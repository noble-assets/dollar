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

func (k vaultsMsgServer) Lock(ctx context.Context, msg *vaults.MsgLock) (*vaults.MsgLockResponse, error) {
	if paused := k.GetVaultsPaused(ctx); paused == vaults.ALL || paused == vaults.LOCK {
		return nil, errors.Wrapf(vaults.ErrActionPaused, "lock is paused")
	}

	// Ensure that the signer is a valid address.
	addr, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("unable to decode user address: %s", msg.Signer)
	}

	// Ensure that the amount is at least the `vaultsMinimumLock`.
	if msg.Amount.LT(math.NewInt(k.vaultsMinimumLock)) {
		return nil, errors.Wrapf(
			vaults.ErrInvalidAmount,
			"must lock at least %d%s",
			k.vaultsMinimumLock,
			k.denom,
		)
	}

	// Ensure that the Vault type does exist.
	_, vaultTypeExists := vaults.VaultType_value[msg.Vault.String()]
	if !vaultTypeExists || msg.Vault == vaults.UNSPECIFIED {
		return nil, errors.Wrapf(vaults.ErrInvalidVaultType, "vault type %s does not exist", msg.Vault)
	}

	currentTime := k.header.GetHeaderInfo(ctx).Time.Unix()

	// Verify that no position from the same user and vault exists within the current block.
	key := collections.Join3(addr, int32(msg.Vault), currentTime)
	if has, _ := k.VaultsPositions.Has(ctx, key); has {
		return nil, errors.Wrapf(vaults.ErrInvalidVaultType, "cannot create multiple user positions in the same block")
	}

	// Verify that the user has sufficient balance.
	if k.bank.GetBalance(ctx, addr, k.denom).Amount.LT(msg.Amount) {
		return nil, errors.Wrapf(vaults.ErrInvalidAmount, "insufficient balance")
	}

	vaultUserAccount := authtypes.NewEmptyModuleAccount(k.ToUserVaultPositionModuleAccount(msg.Signer, msg.Vault, currentTime))
	vaultAccount := k.account.NewAccount(ctx, vaultUserAccount).(*authtypes.ModuleAccount)
	k.account.SetModuleAccount(ctx, vaultAccount)

	// Transfer the specified amount from the user to the submodule Vault account.
	if err = k.bank.SendCoins(ctx,
		addr,
		vaultAccount.GetAddress(),
		sdk.NewCoins(sdk.NewCoin(k.denom, msg.Amount)),
	); err != nil {
		return nil, err
	}

	// Get the current Index.
	index, err := k.Index.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get index from state")
	}
	amountPrincipal := k.GetPrincipalAmountRoundedDown(msg.Amount, index)

	// Create the user Vault Position.
	if err = k.VaultsPositions.Set(ctx, key, vaults.Position{
		Index:     index,
		Principal: amountPrincipal,
		Amount:    msg.Amount,
		Time:      k.header.GetHeaderInfo(ctx).Time,
	}); err != nil {
		return nil, errors.Wrapf(err, "unable to set position")
	}

	// If the Vault type is Flexible, handle the additional login.
	if msg.Vault == vaults.FLEXIBLE {
		// Increase the Total Flexible Principal
		total := math.ZeroInt()
		if has, _ := k.VaultsTotalFlexiblePrincipal.Has(ctx); has {
			current, err := k.VaultsTotalFlexiblePrincipal.Get(ctx)
			if err != nil {
				return nil, err
			}
			total = total.Add(current)
		}
		if err = k.VaultsTotalFlexiblePrincipal.Set(ctx, total.Add(amountPrincipal)); err != nil {
			return nil, err
		}
	}

	// Update Vaults stats.
	if positions, _ := k.GetVaultsPositionsByProviderAndVault(ctx, addr, vaults.VaultType_value[msg.Vault.String()]); len(positions) == 1 {
		if err = k.IncrementVaultUsers(ctx, msg.Vault); err != nil {
			return nil, errors.Wrap(err, "unable to increment vault total users")
		}
	}
	if err = k.IncrementVaultTotalPrincipal(ctx, msg.Vault, amountPrincipal); err != nil {
		return nil, errors.Wrap(err, "unable to increment vault total principal")
	}

	return &vaults.MsgLockResponse{}, k.event.EventManager(ctx).Emit(ctx, &vaults.PositionLocked{
		Account:   msg.Signer,
		VaultType: msg.Vault.String(),
		Index:     index,
		Amount:    msg.Amount,
		Principal: amountPrincipal,
	})
}

func (k vaultsMsgServer) Unlock(ctx context.Context, msg *vaults.MsgUnlock) (*vaults.MsgUnlockResponse, error) {
	return k.unlock(ctx, msg)
}

func (k *Keeper) unlock(ctx context.Context, msg *vaults.MsgUnlock) (*vaults.MsgUnlockResponse, error) {
	if paused := k.GetVaultsPaused(ctx); paused == vaults.ALL || paused == vaults.UNLOCK {
		return nil, errors.Wrapf(vaults.ErrActionPaused, "unlock is paused")
	}

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

		// If the vault is Flexible, handle the additional logic.
		if msg.Vault == vaults.FLEXIBLE {
			// Get total Vault principal.
			totalVaultUsersPrincipal := math.ZeroInt()
			if has, _ := k.VaultsTotalFlexiblePrincipal.Has(ctx); has {
				current, err := k.VaultsTotalFlexiblePrincipal.Get(ctx)
				if err != nil {
					return nil, err
				}
				totalVaultUsersPrincipal = totalVaultUsersPrincipal.Add(current)
			}
			// Deduct the relative position's principal amount from the TotalVaultUsersPrincipal.
			if err = k.VaultsTotalFlexiblePrincipal.Set(ctx, totalVaultUsersPrincipal.Sub(positionPrincipalToRemove)); err != nil {
				return nil, errors.Wrapf(err, "unable to set position for %s", msg.Vault)
			}

			// Claim the rewards associated to the current position.
			rewards, err := k.ClaimRewards(ctx, position, positionAmountToRemove)
			if err != nil {
				return nil, err
			}
			if err = k.IncrementFlexibleTotalDistributedRewardsPrincipal(ctx, rewards); err != nil {
				return nil, errors.Wrap(err, "unable to increment flexible vault total distributed rewards principal")
			}

			// Claim the yield associated to the current position.
			yield, err := k.claimModuleYield(ctx, authtypes.NewModuleAddress(k.ToUserVaultPositionModuleAccount(msg.Signer, msg.Vault, position.Time.Unix())))
			if err != nil {
				return nil, err
			}
			if yield.IsPositive() {
				amountToSend = positionAmountToRemove.Add(yield)
			}

		}

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

func (k *Keeper) ClaimRewards(ctx context.Context, position vaults.PositionEntry, amount math.Int) (math.Int, error) {
	userAddress, err := k.address.BytesToString(position.Address)
	if err != nil {
		return math.Int{}, err
	}

	// Get the total Vault rewards.
	vaultRewardsPrincipal := math.ZeroInt()
	if has, _ := k.Principal.Has(ctx, vaults.FlexibleVaultAddress); has {
		current, err := k.Principal.Get(ctx, vaults.FlexibleVaultAddress)
		if err != nil {
			return math.ZeroInt(), err
		}
		vaultRewardsPrincipal = vaultRewardsPrincipal.Add(current)
	}

	// Exit if there are no rewards.
	if !vaultRewardsPrincipal.IsPositive() {
		return math.ZeroInt(), nil
	}

	// Retrieve the current Index and amount Principal.
	currentIndex, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "unable to get index from state")
	}
	amountPrincipal := k.GetPrincipalAmountRoundedDown(amount, position.Index)

	// Iterate through the rewards to calculate the amount owed to the user, proportional to their position.
	// NOTE: For the user to be eligible, they must have joined before and exited after a complete `UpdateIndex` cycle.
	rewardsAmount := math.ZeroInt()
	if err := k.VaultsRewards.Walk(
		ctx,
		new(collections.Range[int64]).StartExclusive(position.Index), // Exclude the entry point Index.
		func(key int64, record vaults.Reward) (stop bool, err error) {
			if !record.Total.IsPositive() || !record.Rewards.IsPositive() {
				return false, nil
			}

			// Exclude the last Index.
			userReward := math.ZeroInt()
			if record.Index != currentIndex && !record.Rewards.IsNegative() {
				userReward = record.Rewards.ToLegacyDec().Quo(record.Total.ToLegacyDec()).MulInt(amountPrincipal).TruncateInt()
			}

			// Update the Rewards entry.
			if err = k.VaultsRewards.Set(ctx, key, vaults.Reward{
				Index:   record.Index,
				Total:   record.Total.Sub(amountPrincipal),
				Rewards: record.Rewards.Sub(userReward),
			}); err != nil {
				return true, err
			}
			rewardsAmount = rewardsAmount.Add(userReward)
			return false, nil
		}); err != nil {
		return math.ZeroInt(), nil
	}

	// Transfer the specified amount back to the user from the submodule's Vault account.
	err = k.bank.SendCoins(ctx,
		authtypes.NewModuleAddress(vaults.FlexibleVaultName),
		position.Address,
		sdk.NewCoins(sdk.NewCoin(k.denom, rewardsAmount)),
	)
	if err != nil {
		return math.ZeroInt(), err
	}

	return rewardsAmount, k.event.EventManager(ctx).Emit(ctx, &vaults.RewardClaimed{
		Account: userAddress,
		Amount:  rewardsAmount,
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

func (k *Keeper) VaultsEndProgram(ctx context.Context) error {
	k.VaultsProgramEnded.Set(ctx, true)

	k.handleVaultProgramEndUnlocks(ctx)

	// Pause any further interaction with the vaults.
	// This must be executed no matter what
	k.VaultsPaused.Set(ctx, vaults.PausedType_value[vaults.ALL.String()])

	k.logger.Warn("Program Ended!")

	return nil
}

func (k *Keeper) handleVaultProgramEndUnlocks(ctx context.Context) error {
	/* k.VaultsProgramEnded.Set(ctx, true) */

	// Pause any further interaction with the vaults.
	k.logger.Warn("AAAA")

	// Get all the Vaults positions.
	positions, err := k.GetVaultsPositions(ctx)
	if err != nil {
		return err
	}

	k.logger.Info("Collecting Positions...")
	// Create a mapping between the different vault types by the address and the total positions amount
	stakedUsers := map[string]math.Int{}
	flexibleUsers := map[string]math.Int{}

	// Iterate through all the positions
	for _, position := range positions {
		addr, err := k.address.BytesToString(position.Address)
		if err != nil {
			//TODO: what to do with an invalid address? does this case even exists?
			continue
		}
		if position.Vault == vaults.FLEXIBLE {
			if _, exists := flexibleUsers[addr]; !exists {
				flexibleUsers[addr] = position.Amount
			} else {
				flexibleUsers[addr] = flexibleUsers[addr].Add(position.Amount)
			}
		} else {
			if _, exists := stakedUsers[addr]; !exists {
				stakedUsers[addr] = position.Amount
			} else {
				stakedUsers[addr] = stakedUsers[addr].Add(position.Amount)
			}
		}
	}

	// First of all we remove all the Staked Vaults positions.
	k.logger.Info("Unlocking staked vault positions...")
	for stakedUserAddr, stakedUserTotalAmount := range stakedUsers {
		if _, err := k.unlock(ctx, &vaults.MsgUnlock{
			Signer: stakedUserAddr,
			Vault:  vaults.STAKED,
			Amount: stakedUserTotalAmount,
		}); err != nil {
			k.logger.Error(err.Error())
		}
	}

	k.logger.Info("Unlocking flexible vault positions...")
	for flexibleUserAddr, flexibleUserTotalAmount := range flexibleUsers {
		if _, err := k.unlock(ctx, &vaults.MsgUnlock{
			Signer: flexibleUserAddr,
			Vault:  vaults.FLEXIBLE,
			Amount: flexibleUserTotalAmount,
		}); err != nil {
			k.logger.Error(err.Error())
		}
	}

	return nil
}
