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

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/vaults"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) ClaimYield(ctx context.Context, msg *types.MsgClaimYield) (*types.MsgClaimYieldResponse, error) {
	if k.GetPaused(ctx) {
		return nil, types.ErrPaused
	}

	yield, account, err := k.GetYield(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, account, sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to distribute yield to user")
	}

	return &types.MsgClaimYieldResponse{}, nil
}

func (k msgServer) SetPausedState(ctx context.Context, msg *types.MsgSetPausedState) (*types.MsgSetPausedStateResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != k.authority {
		return nil, errors.Wrapf(vaults.ErrInvalidAuthority, "expected %s, got %s", k.authority, msg.Signer)
	}

	if err := k.Paused.Set(ctx, msg.Paused); err != nil {
		return nil, err
	}

	return &types.MsgSetPausedStateResponse{}, nil
}

func (k *Keeper) Burn(ctx context.Context, sender []byte, amount math.Int) error {
	coins := sdk.NewCoins(sdk.NewCoin(k.denom, amount))
	err := k.bank.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	if err != nil {
		return err
	}
	err = k.bank.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Mint(ctx context.Context, recipient []byte, amount math.Int, index *int64) error {
	if index != nil {
		_ = k.UpdateIndex(ctx, *index)
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.denom, amount))
	err := k.bank.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) UpdateIndex(ctx context.Context, rawIndex int64) error {
	oldIndex, err := k.Index.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get index from state")
	}
	if rawIndex <= oldIndex {
		return types.ErrDecreasingIndex
	}

	err = k.Index.Set(ctx, rawIndex)
	if err != nil {
		return errors.Wrap(err, "unable to set index in state")
	}

	totalPrincipal, err := k.GetTotalPrincipal(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get total principal from state")
	}

	index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)
	currentSupply := k.bank.GetSupply(ctx, k.denom).Amount
	expectedSupply := index.MulInt(totalPrincipal).TruncateInt()

	coins := sdk.NewCoins(sdk.NewCoin(k.denom, expectedSupply.Sub(currentSupply)))
	if coins.IsAllPositive() {
		err = k.bank.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return errors.Wrap(err, "unable to mint coins")
		}
		err = k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.YieldName, coins)
		if err != nil {
			return errors.Wrap(err, "unable to send coins")
		}
		err = k.IncrementTotalYieldAccrued(ctx, coins.AmountOf(k.denom))
		if err != nil {
			return errors.Wrap(err, "unable to increment total yield accrued")
		}
	}

	// Claim the yield of the Flexible vault.
	flexibleYield, err := k.claimModuleYield(ctx, vaults.FlexibleVaultAddress)
	if err != nil {
		return err
	}

	// Claim the yield of the Staked vault and redirect it to the Flexible vault.
	stakedYield, err := k.claimStakedVaultYield(ctx)
	if err != nil {
		return err
	}

	// get the current Flexible total principal.
	totalFlexiblePrincipal := math.ZeroInt()
	if has, _ := k.VaultsTotalFlexiblePrincipal.Has(ctx); has {
		current, err := k.VaultsTotalFlexiblePrincipal.Get(ctx)
		if err != nil {
			return err
		}
		totalFlexiblePrincipal = totalFlexiblePrincipal.Add(current)
	}

	// Register the new Rewards record.
	rewards := stakedYield.Add(flexibleYield)
	if err = k.VaultsRewards.Set(ctx, index.String(), vaults.Reward{
		Index:   index,
		Total:   totalFlexiblePrincipal,
		Rewards: rewards,
	}); err != nil {
		return err
	}
	return nil
}

func (k *Keeper) claimModuleYield(ctx context.Context, addr sdk.Address) (math.Int, error) {
	// Get the Yield of the module address.
	yield, _, err := k.GetYield(ctx, addr.String())
	if err != nil {
		return math.ZeroInt(), nil
	}

	// If the Yield exists, claim it.
	if yield.IsPositive() {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, addr.Bytes(), sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	return yield, nil
}

func (k *Keeper) claimStakedVaultYield(ctx context.Context) (math.Int, error) {
	// Get the Yield of the Staked Vault.
	yield, _, err := k.GetYield(ctx, vaults.StakedVaultAddress.String())
	if err != nil {
		return math.ZeroInt(), nil
	}
	if !yield.IsPositive() {
		return math.ZeroInt(), nil
	}

	// Redirect the Yield from the Yield Module to the Flexible Vault instead of the Staked Vault.
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, vaults.FlexibleVaultAddress, sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
	if err != nil {
		return math.ZeroInt(), err
	}

	// Get the current Index.
	rawIndex, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)
	// Calculate the Yield Principal.
	yieldPrincipal := yield.ToLegacyDec().Quo(index).TruncateInt()

	// Reduce the Staked Vault Principal by the Yield Principal.
	stakedPrincipal, err := k.Principal.Get(ctx, vaults.StakedVaultAddress)
	if err != nil {
		stakedPrincipal = math.ZeroInt()
	}
	if err = k.Principal.Set(ctx, vaults.StakedVaultAddress, stakedPrincipal.Sub(yieldPrincipal)); err != nil {
		return math.ZeroInt(), err
	}

	// Add the Yield Principal to the Flexible Vault Principal.
	flexiblePrincipal, err := k.Principal.Get(ctx, vaults.FlexibleVaultAddress)
	if err != nil {
		flexiblePrincipal = math.ZeroInt()
	}
	if err = k.Principal.Set(ctx, vaults.FlexibleVaultAddress, flexiblePrincipal.Add(yieldPrincipal)); err != nil {
		return math.ZeroInt(), err
	}
	return yield, nil
}
