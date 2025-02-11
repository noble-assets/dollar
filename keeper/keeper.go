// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
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
	"dollar.noble.xyz/types/portal/ntt"
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

	Paused    collections.Item[bool]
	Index     collections.Item[int64]
	Principal collections.Map[[]byte, math.Int]

	PortalOwner  collections.Item[string]
	PortalPaused collections.Item[bool]
	PortalPeers  collections.Map[uint16, portal.Peer]
	PortalNonce  collections.Item[uint32]

	VaultsPaused                 collections.Item[int32]
	VaultsPositions              *collections.IndexedMap[collections.Triple[[]byte, int32, int64], vaults.Position, VaultsPositionsIndexes]
	VaultsTotalFlexiblePrincipal collections.Item[math.Int]
	VaultsRewards                collections.Map[string, vaults.Reward]
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

		Paused:    collections.NewItem(builder, types.PausedKey, "paused", collections.BoolValue),
		Index:     collections.NewItem(builder, types.IndexKey, "index", collections.Int64Value),
		Principal: collections.NewMap(builder, types.PrincipalPrefix, "principal", collections.BytesKey, sdk.IntValue),

		PortalOwner:  collections.NewItem(builder, portal.OwnerKey, "portal_owner", collections.StringValue),
		PortalPaused: collections.NewItem(builder, portal.PausedKey, "portal_paused", collections.BoolValue),
		PortalPeers:  collections.NewMap(builder, portal.PeerPrefix, "portal_peers", collections.Uint16Key, codec.CollValue[portal.Peer](cdc)),
		PortalNonce:  collections.NewItem(builder, portal.NonceKey, "portal_nonce", collections.Uint32Value),

		VaultsPaused:                 collections.NewItem(builder, vaults.PausedKey, "vaults_paused", collections.Int32Value),
		VaultsPositions:              collections.NewIndexedMap(builder, vaults.PositionPrefix, "vaults_positions", collections.TripleKeyCodec(collections.BytesKey, collections.Int32Key, collections.Int64Key), codec.CollValue[vaults.Position](cdc), NewVaultsPositionsIndexes(builder)),
		VaultsTotalFlexiblePrincipal: collections.NewItem(builder, vaults.TotalFlexiblePrincipalKey, "vaults_total_flexible_principal", sdk.IntValue),
		VaultsRewards:                collections.NewMap(builder, vaults.RewardPrefix, "vaults_rewards", collections.StringKey, codec.CollValue[vaults.Reward](cdc)),
	}

	_, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return keeper
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bank = bankKeeper
}

// SendRestrictionFn performs an underlying transfer of principal when executing a $USDN transfer.
func (k *Keeper) SendRestrictionFn(ctx context.Context, sender, recipient sdk.AccAddress, coins sdk.Coins) (newRecipient sdk.AccAddress, err error) {
	if amount := coins.AmountOf(k.denom); !amount.IsZero() {
		// We don't want to perform any principal updates in the case of yield payout.
		// -> Transfer from Yield to User account.
		if sender.Equals(types.YieldAddress) {
			return recipient, nil
		}
		// Handle transfers where the recipient is the yield account.
		if recipient.Equals(types.YieldAddress) {
			if sender.Equals(types.ModuleAddress) {
				// We don't want to perform any principal updates in the case of yield accrual.
				// -> Transfer from Module to Yield account.
				return recipient, nil
			} else {
				// We don't want to allow any other transfers to the yield account.
				// -> Transfer from User to Yield account.
				return recipient, fmt.Errorf("transfers of %s to %s are not allowed", k.denom, recipient.String())
			}
		}

		rawIndex, err := k.Index.Get(ctx)
		if err != nil {
			return recipient, errors.Wrap(err, "unable to get index from state")
		}
		index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)
		principal := amount.ToLegacyDec().Quo(index).TruncateInt()

		// We don't want to update the sender's principal in the case of issuance.
		// -> Transfer from Module to User account.
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

		// We don't want to update the recipient's principal in the case of withdrawal.
		// -> Transfer from User to Module account.
		if !recipient.Equals(types.ModuleAddress) {
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
	}

	return recipient, nil
}

// GetDenom is a utility that returns the configured denomination of $USDN.
func (k *Keeper) GetDenom() string {
	return k.denom
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
	expectedBalance := index.MulInt(principal).TruncateInt()

	yield, _ := expectedBalance.SafeSub(currentBalance)

	// TODO: temporary fix for negative coin amounts
	if yield.Abs().Equal(math.OneInt()) || yield.IsNegative() {
		return math.ZeroInt(), nil, nil
	}

	return yield, bz, nil
}

func (k *Keeper) DeliverInjections(ctx context.Context, msg *portal.MsgDeliverInjection) error {
	vaa, err := k.wormhole.ParseAndVerifyVAA(ctx, msg.Vaa)
	if err != nil {
		return err
	}

	peer, err := k.Peers.Get(ctx, uint16(vaa.EmitterChain))
	if err != nil {
		return errors.Wrapf(portal.ErrInvalidPeer, "chain %d not configured", vaa.EmitterChain)
	}

	if !bytes.Equal(peer.Transceiver, vaa.EmitterAddress.Bytes()) {
		return errors.Wrapf(
			portal.ErrInvalidPeer,
			"expected transceiver %s for chain %d, got %s",
			hex.EncodeToString(peer.Transceiver), vaa.EmitterChain,
			vaa.EmitterAddress.String(),
		)
	}

	transceiverMessage, err := ntt.ParseTransceiverMessage(vaa.Payload)
	if err != nil {
		return err
	}

	if !bytes.Equal(peer.Manager, transceiverMessage.SourceManagerAddress) {
		return errors.Wrapf(
			portal.ErrInvalidPeer,
			"expected manager %s for chain %d, got %s",
			hex.EncodeToString(peer.Manager), vaa.EmitterChain,
			hex.EncodeToString(transceiverMessage.SourceManagerAddress),
		)
	}

	if !bytes.Equal(portal.PaddedManagerAddress, transceiverMessage.RecipientManagerAddress) {
		return errors.Wrapf(
			portal.ErrInvalidMessage,
			"expected recipient manager %s, got %s",
			hex.EncodeToString(portal.PaddedManagerAddress),
			hex.EncodeToString(transceiverMessage.RecipientManagerAddress),
		)
	}

	managerMessage, err := ntt.ParseManagerMessage(transceiverMessage.ManagerPayload)
	if err != nil {
		return err
	}

	return k.HandlePayload(ctx, managerMessage.Payload)
}

// HandlePayload is a utility that handles custom payloads when delivering portal messages.
func (k *Keeper) HandlePayload(ctx context.Context, payload []byte) error {
	chain, _ := k.wormhole.GetChain(ctx)

	switch portal.GetPayloadType(payload) {
	case portal.Unknown:
		return nil
	case portal.Token:
		amount, index, recipient, destination := portal.DecodeTokenPayload(payload)
		if chain != destination {
			return fmt.Errorf("not destination chain: expected %d, got %d", chain, destination)
		}

		return k.Mint(ctx, recipient, amount, &index)
	case portal.Index:
		index, destination := portal.DecodeIndexPayload(payload)
		if chain != destination {
			return fmt.Errorf("not destination chain: expected %d, got %d", chain, destination)
		}

		return k.UpdateIndex(ctx, index)
	}

	return nil
}
