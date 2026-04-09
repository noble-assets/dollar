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
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"dollar.noble.xyz/v2/types/vaults"
)

var _ vaults.MsgServer = &vaultsMsgServer{}

type vaultsMsgServer struct {
	*Keeper
}

func NewVaultsMsgServer(keeper *Keeper) vaults.MsgServer {
	return &vaultsMsgServer{Keeper: keeper}
}

func (k vaultsMsgServer) Lock(_ context.Context, _ *vaults.MsgLock) (*vaults.MsgLockResponse, error) {
	return nil, errors.Wrapf(vaults.ErrActionPaused, "locking is unsupported")
}

func (k vaultsMsgServer) Unlock(ctx context.Context, msg *vaults.MsgUnlock) (*vaults.MsgUnlockResponse, error) {
	if paused := k.GetVaultsPaused(ctx); paused == vaults.ALL || paused == vaults.UNLOCK {
		return nil, errors.Wrapf(vaults.ErrActionPaused, "unlock is paused")
	}

	return k.unlock(ctx, msg)
}

func (k *Keeper) unlock(ctx context.Context, msg *vaults.MsgUnlock) (*vaults.MsgUnlockResponse, error) {
	// Ensure that the signer is a valid address.
	addr, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("unable to decode user address: %s", msg.Signer)
	}

	// Ensure that the amount is valid.
	if msg.Amount.IsNil() || !msg.Amount.IsPositive() {
		return nil, errors.Wrapf(vaults.ErrInvalidAmount, "amount is zero")
	}

	// Ensure that the Vault type does exist.
	_, vaultTypeExists := vaults.VaultType_value[msg.Vault.String()]
	if !vaultTypeExists || msg.Vault == vaults.UNSPECIFIED {
		return nil, errors.Wrapf(vaults.ErrInvalidVaultType, "vault type %s does not exist", msg.Vault)
	}

	// Retrieve all positions associated with the user.
	positions, err := k.GetVaultsPositionsByProviderAndVault(ctx, addr, vaults.VaultType_value[msg.Vault.String()])
	if err != nil {
		return nil, err
	}

	// Calculate the total user positions amount.
	totalPositions := math.ZeroInt()
	for _, position := range positions {
		totalPositions = totalPositions.Add(position.Amount)
	}

	// Early check to ensure that the user has a sufficient locked amount.
	if msg.Amount.GT(totalPositions) {
		return nil, errors.Wrapf(
			vaults.ErrInvalidAmount,
			"%s%s is greater than the total amount left of %s%s",
			msg.Amount,
			k.denom,
			totalPositions.String(),
			k.denom,
		)
	}

	// Ensure that the amount to unlock is at least `vaultsMinimumUnlock`
	// or the total remaining position when the remaining amount is less than `vaultsMinimumUnlock`.
	if msg.Amount.LT(math.NewInt(k.vaultsMinimumUnlock)) && !totalPositions.Equal(msg.Amount) {
		if !msg.Amount.Equal(totalPositions) && totalPositions.LT(math.NewInt(k.vaultsMinimumUnlock)) {
			return nil, errors.Wrapf(
				vaults.ErrInvalidAmount,
				"must unlock the total amount left of %s%s",
				totalPositions.String(),
				k.denom,
			)
		}
		return nil, errors.Wrapf(
			vaults.ErrInvalidAmount,
			"must unlock at least %d%s",
			k.vaultsMinimumUnlock,
			k.denom,
		)
	}

	// Iterate through the user's positions until the required principal amount for removal is reached.
	remainingAmountToRemove := msg.Amount
	removedPrincipal := math.ZeroInt()
	for _, position := range positions {
		// Exit when the amount to remove is zero.
		if remainingAmountToRemove.IsZero() {
			break
		}

		// Determine the amount and principal to remove from the current position, either partially or in full, and to send to the user.
		positionAmountToRemove := position.Amount
		if position.Amount.GT(remainingAmountToRemove) {
			positionAmountToRemove = remainingAmountToRemove
		}
		amountToSend := positionAmountToRemove
		positionPrincipalToRemove := k.GetPrincipalAmountRoundedDown(positionAmountToRemove, position.Index)

		// Transfer the specified amount from submodule's vault account to the user.
		err = k.bank.SendCoins(ctx,
			authtypes.NewModuleAddress(k.ToUserVaultPositionModuleAccount(msg.Signer, position.Vault, position.Time.Unix())),
			addr,
			sdk.NewCoins(sdk.NewCoin(k.denom, amountToSend)),
		)
		if err != nil {
			return nil, err
		}

		// Remove or update the user's position.
		if positionAmountToRemove.GTE(position.Amount) {
			if err = k.VaultsPositions.Remove(ctx, collections.Join3(position.Address, int32(position.Vault), position.Time.Unix())); err != nil {
				return nil, errors.Wrapf(err, "unable to remove position")
			}
		} else {
			updatedPrincipal := position.Principal.Sub(positionPrincipalToRemove)
			if err = k.VaultsPositions.Set(ctx, collections.Join3(position.Address, int32(position.Vault), position.Time.Unix()), vaults.Position{
				Principal: updatedPrincipal,
				Index:     position.Index,
				Amount:    position.Amount.Sub(positionAmountToRemove),
				Time:      position.Time,
			}); err != nil {
				return nil, errors.Wrapf(err, "unable to update position")
			}
		}

		if err = k.event.EventManager(ctx).Emit(ctx, &vaults.PositionUnlocked{
			Account:   msg.Signer,
			VaultType: msg.Vault.String(),
			Index:     position.Index,
			Amount:    amountToSend,
			Principal: positionPrincipalToRemove,
		}); err != nil {
			return nil, errors.Wrap(err, "unable to emit position unlocked event")
		}

		removedPrincipal = removedPrincipal.Add(positionPrincipalToRemove)

		// Update the remaining amount to be removed.
		remainingAmountToRemove = remainingAmountToRemove.Sub(positionAmountToRemove)
	}

	if !remainingAmountToRemove.IsZero() || !remainingAmountToRemove.Abs().Equal(math.ZeroInt()) {
		return nil, errors.Wrapf(vaults.ErrInvalidAmount, "invalid amount left: %s", remainingAmountToRemove.String())
	}

	// Update Vaults stats.
	if positions, _ = k.GetVaultsPositionsByProviderAndVault(ctx, addr, vaults.VaultType_value[msg.Vault.String()]); len(positions) == 0 {
		if err = k.DecrementVaultUsers(ctx, msg.Vault); err != nil {
			return nil, errors.Wrap(err, "unable to decrement vault total users")
		}
	}
	if err = k.DecrementVaultTotalPrincipal(ctx, msg.Vault, removedPrincipal); err != nil {
		return nil, errors.Wrap(err, "unable to decrement vault total principal")
	}

	return &vaults.MsgUnlockResponse{}, nil
}

func (k vaultsMsgServer) SetPausedState(ctx context.Context, msg *vaults.MsgSetPausedState) (*vaults.MsgSetPausedStateResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != k.authority {
		return nil, errors.Wrapf(vaults.ErrInvalidAuthority, "expected %s, got %s", k.authority, msg.Signer)
	}

	// Ensure that the Pause type does exist.
	_, pausedTypeExists := vaults.PausedType_value[msg.Paused.String()]
	if !pausedTypeExists {
		return nil, errors.Wrapf(vaults.ErrInvalidPauseType, "vaults pause type %s does not exist", msg.Paused)
	}

	// Set the new Paused status.
	if err := k.VaultsPaused.Set(ctx, int32(msg.Paused)); err != nil {
		return nil, err
	}

	return &vaults.MsgSetPausedStateResponse{}, k.event.EventManager(ctx).Emit(ctx, &vaults.PausedStateUpdated{
		Paused: msg.Paused.String(),
	})
}

func (k *Keeper) ToUserVaultPositionModuleAccount(address string, vaultType vaults.VaultType, timestamp int64) string {
	if vaultType == vaults.FLEXIBLE {
		// Flexible Vaults use individual accounts for each user position.
		return fmt.Sprintf("%s/%s/%s/%d", vaults.SubmoduleName, strings.ToLower(vaultType.String()), strings.ToLower(address), timestamp)
	} else {
		// Staked Vaults use a shared account for all users.
		return vaults.StakedVaultName
	}
}
