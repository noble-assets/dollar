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
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	sdkerrors "cosmossdk.io/errors"
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
	Stats     collections.Item[types.Stats]

	PortalOwner  collections.Item[string]
	PortalPaused collections.Item[bool]
	PortalPeers  collections.Map[uint16, portal.Peer]
	PortalNonce  collections.Item[uint32]

	VaultsPaused                 collections.Item[int32]
	VaultsPositions              *collections.IndexedMap[collections.Triple[[]byte, int32, int64], vaults.Position, VaultsPositionsIndexes]
	VaultsTotalFlexiblePrincipal collections.Item[math.Int]
	VaultsRewards                collections.Map[int64, vaults.Reward]
	VaultsStats                  collections.Item[vaults.Stats]
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
		Stats:     collections.NewItem(builder, types.StatsKey, "stats", codec.CollValue[types.Stats](cdc)),

		PortalOwner:  collections.NewItem(builder, portal.OwnerKey, "portal_owner", collections.StringValue),
		PortalPaused: collections.NewItem(builder, portal.PausedKey, "portal_paused", collections.BoolValue),
		PortalPeers:  collections.NewMap(builder, portal.PeerPrefix, "portal_peers", collections.Uint16Key, codec.CollValue[portal.Peer](cdc)),
		PortalNonce:  collections.NewItem(builder, portal.NonceKey, "portal_nonce", collections.Uint32Value),

		VaultsPaused:                 collections.NewItem(builder, vaults.PausedKey, "vaults_paused", collections.Int32Value),
		VaultsPositions:              collections.NewIndexedMap(builder, vaults.PositionPrefix, "vaults_positions", collections.TripleKeyCodec(collections.BytesKey, collections.Int32Key, collections.Int64Key), codec.CollValue[vaults.Position](cdc), NewVaultsPositionsIndexes(builder)),
		VaultsTotalFlexiblePrincipal: collections.NewItem(builder, vaults.TotalFlexiblePrincipalKey, "vaults_total_flexible_principal", sdk.IntValue),
		VaultsRewards:                collections.NewMap(builder, vaults.RewardPrefix, "vaults_rewards", collections.Int64Key, codec.CollValue[vaults.Reward](cdc)),
		VaultsStats:                  collections.NewItem(builder, vaults.StatsKey, "vaults_stats", codec.CollValue[vaults.Stats](cdc)),
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
	coin := coins.AmountOf(k.denom)
	if amount := coin; !amount.IsZero() {
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

		// When burning and transferring, the $M token executes the
		// `_getPrincipalAmountRoundedUp` function. When minting, the $M token
		// executes the `_getPrincipalAmountRoundedDown` function. As $USDN
		// inherits the yielding properties of $M, we mimic that functionality
		// here.
		index, err := k.Index.Get(ctx)
		if err != nil {
			return recipient, sdkerrors.Wrap(err, "unable to get index from state")
		}
		principal := k.GetPrincipalAmountRoundedUp(amount, index)

		// We don't want to update the sender's principal in the case of issuance.
		// -> Transfer from Module to User account.
		if !sender.Equals(types.ModuleAddress) {
			senderPrincipal, err := k.Principal.Get(ctx, sender)
			if err != nil {
				if sdkerrors.IsOf(err, collections.ErrNotFound) {
					senderPrincipal = math.ZeroInt()
				} else {
					return recipient, sdkerrors.Wrap(err, "unable to get sender principal from state")
				}
			}
			err = k.Principal.Set(ctx, sender, senderPrincipal.Sub(principal))
			if err != nil {
				return recipient, sdkerrors.Wrap(err, "unable to set sender principal to state")
			}

			balance := k.bank.GetBalance(ctx, sender, k.denom)
			if balance.IsZero() {
				// If the sender's $USDN balance is zero, this indicates that
				// they are no longer a holder, and we should decrement the
				// statistic.
				err = k.DecrementTotalHolders(ctx)
				if err != nil {
					return recipient, sdkerrors.Wrap(err, "unable to decrement total holders")
				}
			}
		} else {
			err = k.IncrementTotalPrincipal(ctx, principal)
			if err != nil {
				return recipient, sdkerrors.Wrap(err, "unable to increment total principal")
			}
		}

		// We don't want to update the recipient's principal in the case of withdrawal.
		// -> Transfer from User to Module account.
		if !recipient.Equals(types.ModuleAddress) {
			if sender.Equals(types.ModuleAddress) {
				principal = k.GetPrincipalAmountRoundedDown(amount, index)
			}

			recipientPrincipal, err := k.Principal.Get(ctx, recipient)
			if err != nil {
				if sdkerrors.IsOf(err, collections.ErrNotFound) {
					recipientPrincipal = math.ZeroInt()
				} else {
					return recipient, sdkerrors.Wrap(err, "unable to get recipient principal from state")
				}
			}
			err = k.Principal.Set(ctx, recipient, recipientPrincipal.Add(principal))
			if err != nil {
				return recipient, sdkerrors.Wrap(err, "unable to set recipient principal to state")
			}

			balance := k.bank.GetBalance(ctx, recipient, k.denom)
			if balance.IsZero() {
				// If the recipient's $USDN balance is zero, this indicates
				// that they are a new holder, and we should increment the
				// statistic.
				err = k.IncrementTotalHolders(ctx)
				if err != nil {
					return recipient, sdkerrors.Wrap(err, "unable to increment total holders")
				}
			}
		} else {
			err = k.DecrementTotalPrincipal(ctx, principal)
			if err != nil {
				return recipient, sdkerrors.Wrap(err, "unable to decrement total principal")
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
		return math.ZeroInt(), nil, sdkerrors.Wrapf(err, "unable to decode account %s", account)
	}

	principal, err := k.Principal.Get(ctx, bz)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return math.ZeroInt(), nil, sdkerrors.Wrapf(err, "unable to get principal for account %s from state", account)
		}

		principal = math.ZeroInt()
	}

	index, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), nil, sdkerrors.Wrap(err, "unable to get index from state")
	}

	currentBalance := k.bank.GetBalance(ctx, bz, k.denom).Amount
	expectedBalance := k.GetPresentAmount(principal, index)

	// We need to make sure that the yield value is valid and > 1.
	yield, err := expectedBalance.SafeSub(currentBalance)
	if err != nil || yield.Equal(math.OneInt()) || yield.IsNegative() {
		yield = math.ZeroInt()
	}

	return yield, bz, nil
}

// Deliver is internal logic executed when delivering portal messages.
func (k *Keeper) Deliver(ctx context.Context, bz []byte) error {
	if k.GetPortalPaused(ctx) {
		return portal.ErrPaused
	}

	vaa, err := k.wormhole.ParseAndVerifyVAA(ctx, bz)
	if err != nil {
		return err
	}

	peer, err := k.PortalPeers.Get(ctx, uint16(vaa.EmitterChain))
	if err != nil {
		return sdkerrors.Wrapf(portal.ErrInvalidPeer, "chain %d not configured", vaa.EmitterChain)
	}

	if !bytes.Equal(peer.Transceiver, vaa.EmitterAddress.Bytes()) {
		return sdkerrors.Wrapf(
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
		return sdkerrors.Wrapf(
			portal.ErrInvalidPeer,
			"expected manager %s for chain %d, got %s",
			hex.EncodeToString(peer.Manager), vaa.EmitterChain,
			hex.EncodeToString(transceiverMessage.SourceManagerAddress),
		)
	}

	if !bytes.Equal(portal.PaddedManagerAddress, transceiverMessage.RecipientManagerAddress) {
		return sdkerrors.Wrapf(
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
