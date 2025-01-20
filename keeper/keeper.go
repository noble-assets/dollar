package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/portal"
	"dollar.noble.xyz/types/vaults"
)

type Keeper struct {
	denom     string
	authority string

	header   header.Service
	event    event.Service
	address  address.Codec
	bank     types.BankKeeper
	account  types.AccountKeeper
	wormhole portal.WormholeKeeper

	Index     collections.Item[int64]
	Principal collections.Map[[]byte, math.Int]

	Owner collections.Item[string]
	Peers collections.Map[uint16, portal.Peer]
	Nonce collections.Item[uint32]

	Paused                 collections.Item[int32]
	Positions              *collections.IndexedMap[collections.Triple[[]byte, int32, int64], vaults.Position, PositionsIndexes]
	TotalFlexiblePrincipal collections.Item[math.Int]
	Rewards                collections.Map[string, vaults.Reward]
}

// Positions Indexes

type PositionsIndexes struct {
	ByProvider *indexes.Multi[[]byte, collections.Triple[[]byte, int32, int64], vaults.Position]
}

func (i PositionsIndexes) IndexesList() []collections.Index[collections.Triple[[]byte, int32, int64], vaults.Position] {
	return []collections.Index[collections.Triple[[]byte, int32, int64], vaults.Position]{
		i.ByProvider,
	}
}

func NewPositionsIndexes(builder *collections.SchemaBuilder) PositionsIndexes {
	return PositionsIndexes{
		ByProvider: indexes.NewMulti(
			builder, []byte("positions_by_provider"), "positions_by_provider",
			collections.BytesKey,
			collections.TripleKeyCodec(collections.BytesKey, collections.Int32Key, collections.Int64Key),
			func(key collections.Triple[[]byte, int32, int64], value vaults.Position) ([]byte, error) {
				return key.K1(), nil
			},
		),
	}
}

//

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bank = bankKeeper
}

func NewKeeper(denom string, authority string, cdc codec.Codec, store store.KVStoreService, header header.Service, event event.Service, address address.Codec, bank types.BankKeeper, account types.AccountKeeper, wormhole portal.WormholeKeeper) *Keeper {
	transceiverAddress := authtypes.NewModuleAddress(fmt.Sprintf("%s/transceiver", portal.SubmoduleName))
	copy(portal.PaddedTransceiverAddress[12:], transceiverAddress)
	portal.TransceiverAddress, _ = address.BytesToString(transceiverAddress)

	managerAddress := authtypes.NewModuleAddress(fmt.Sprintf("%s/manager", portal.SubmoduleName))
	copy(portal.PaddedManagerAddress[12:], managerAddress)
	portal.ManagerAddress, _ = address.BytesToString(managerAddress)

	bz := []byte(denom)
	copy(portal.RawToken[32-len(bz):], bz)

	builder := collections.NewSchemaBuilder(store)

	keeper := &Keeper{
		denom:     denom,
		authority: authority,
		header:    header,
		event:     event,
		address:   address,
		bank:      bank,
		wormhole:  wormhole,
		account:   account,

		Index:     collections.NewItem(builder, types.IndexKey, "index", collections.Int64Value),
		Principal: collections.NewMap(builder, types.PrincipalPrefix, "principal", collections.BytesKey, sdk.IntValue),

		Owner: collections.NewItem(builder, portal.OwnerKey, "owner", collections.StringValue),
		Peers: collections.NewMap(builder, portal.PeerPrefix, "peers", collections.Uint16Key, codec.CollValue[portal.Peer](cdc)),
		Nonce: collections.NewItem(builder, portal.NonceKey, "nonce", collections.Uint32Value),

		Paused:                 collections.NewItem(builder, vaults.PausedKey, "paused", collections.Int32Value),
		Positions:              collections.NewIndexedMap(builder, vaults.PositionPrefix, "positions", collections.TripleKeyCodec(collections.BytesKey, collections.Int32Key, collections.Int64Key), codec.CollValue[vaults.Position](cdc), NewPositionsIndexes(builder)),
		TotalFlexiblePrincipal: collections.NewItem(builder, vaults.TotalFlexiblePrincipalKey, "total_flexible_principal", sdk.IntValue),
		Rewards:                collections.NewMap(builder, vaults.RewardPrefix, "rewards", collections.StringKey, codec.CollValue[vaults.Reward](cdc)),
	}

	_, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return keeper
}

// SendRestrictionFn performs an underlying transfer of principal when executing a $USDN transfer.
func (k *Keeper) SendRestrictionFn(ctx context.Context, sender, recipient sdk.AccAddress, coins sdk.Coins) (newRecipient sdk.AccAddress, err error) {
	if amount := coins.AmountOf(k.denom); !amount.IsZero() {
		if sender.Equals(types.YieldAddress) {
			return recipient, nil
		}

		rawIndex, err := k.Index.Get(ctx)
		if err != nil {
			return recipient, errors.Wrap(err, "unable to get index from state")
		}
		index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)
		principal := amount.ToLegacyDec().Quo(index).TruncateInt()

		if !sender.Equals(types.ModuleAddress) {
			senderPrincipal, err := k.Principal.Get(ctx, sender)
			if err != nil {
				if errors.IsOf(err, collections.ErrNotFound) {
					senderPrincipal = math.ZeroInt()
				} else {
					return recipient, errors.Wrap(err, "unable to get sender principal from state")
				}
			}
			err = k.Principal.Set(ctx, sender, senderPrincipal.Sub(principal))
			if err != nil {
				return recipient, errors.Wrap(err, "unable to set sender principal to state")
			}
		}

		recipientPrincipal, err := k.Principal.Get(ctx, recipient)
		if err != nil {
			if errors.IsOf(err, collections.ErrNotFound) {
				recipientPrincipal = math.ZeroInt()
			} else {
				return recipient, errors.Wrap(err, "unable to get recipient principal from state")
			}
		}
		err = k.Principal.Set(ctx, recipient, recipientPrincipal.Add(principal))
		if err != nil {
			return recipient, errors.Wrap(err, "unable to set recipient principal to state")
		}
	}

	return recipient, nil
}

// GetYield is a utility that returns the user's current amount of claimable $USDN yield.
func (k *Keeper) GetYield(ctx context.Context, account string) (math.Int, []byte, error) {
	bz, err := k.address.StringToBytes(account)
	if err != nil {
		return math.ZeroInt(), nil, errors.Wrapf(err, "unable to decode account %s", account)
	}

	principal, err := k.Principal.Get(ctx, bz)
	if err != nil {
		return math.ZeroInt(), nil, errors.Wrapf(err, "unable to get principal for account %s from state", account)
	}

	rawIndex, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), nil, errors.Wrap(err, "unable to get index from state")
	}
	index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)

	currentBalance := k.bank.GetBalance(ctx, bz, k.denom).Amount
	// TODO(@john): Ensure that we're always rounding down here, to avoid giving users more $USDN than underlying M.
	expectedBalance := index.MulInt(principal).TruncateInt()

	yield, _ := expectedBalance.SafeSub(currentBalance)

	// TODO: temporary fix for negative coin amounts
	if yield.Abs().Equal(math.OneInt()) || yield.IsNegative() {
		return math.ZeroInt(), nil, nil
	}

	return yield, bz, nil
}
