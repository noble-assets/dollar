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

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/vaults"
)

var _ vaults.QueryServer = &vaultsQueryServer{}

type vaultsQueryServer struct {
	*Keeper
}

func NewVaultsQueryServer(keeper *Keeper) vaults.QueryServer {
	return &vaultsQueryServer{Keeper: keeper}
}

func (k vaultsQueryServer) Paused(ctx context.Context, req *vaults.QueryPaused) (*vaults.QueryPausedResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	return &vaults.QueryPausedResponse{
		Paused: k.GetVaultsPaused(ctx),
	}, nil
}

func (k vaultsQueryServer) PositionsByProvider(ctx context.Context, req *vaults.QueryPositionsByProvider) (*vaults.QueryPositionsByProviderResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	addr, err := k.address.StringToBytes(req.Provider)
	if err != nil {
		return nil, types.ErrInvalidRequest
	}

	positions, err := k.GetVaultsPositionsByProvider(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &vaults.QueryPositionsByProviderResponse{
		Positions: positions,
	}, nil
}

func (k vaultsQueryServer) PendingRewards(ctx context.Context, req *vaults.QueryPendingRewards) (*vaults.QueryPendingRewardsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	rewards, err := k.GetVaultsRewards(ctx)
	if err != nil {
		return &vaults.QueryPendingRewardsResponse{
			PendingRewards: math.ZeroInt(),
		}, nil
	}

	totalPendingRewards := math.ZeroInt()
	for _, reward := range rewards {
		totalPendingRewards = totalPendingRewards.Add(reward.Rewards)
	}

	return &vaults.QueryPendingRewardsResponse{
		PendingRewards: totalPendingRewards,
	}, nil
}

func (k vaultsQueryServer) PendingRewardsByProvider(ctx context.Context, req *vaults.QueryPendingRewardsByProvider) (*vaults.QueryPendingRewardsByProviderResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	addr, err := k.address.StringToBytes(req.Provider)
	if err != nil {
		return nil, types.ErrInvalidRequest
	}

	// Retrieve all the user positions in the Flexible Vault.
	positions, err := k.GetVaultsPositionsByProviderAndVault(ctx, addr, vaults.VaultType_value[vaults.FLEXIBLE.String()])
	if err != nil {
		return nil, err
	}

	// Retrieve the current Index.
	currentIndex, err := k.Index.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get index from state")
	}

	totalUserPendingRewards := math.ZeroInt()
	for _, position := range positions {
		amountPrincipal := k.GetPrincipalAmountRoundedDown(position.Amount, position.Index)

		// Iterate through the rewards to calculate the amount owed to the user, proportional to their position.
		// NOTE: For the user to be eligible, they must have joined before and exited after a complete `UpdateIndex` cycle.
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

				totalUserPendingRewards = totalUserPendingRewards.Add(userReward)
				return false, nil
			}); err != nil {
			return nil, err
		}
	}

	return &vaults.QueryPendingRewardsByProviderResponse{
		PendingRewards: totalUserPendingRewards,
	}, nil
}

func (k vaultsQueryServer) Stats(ctx context.Context, req *vaults.QueryStats) (*vaults.QueryStatsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get vaults stats from state")
	}

	return &vaults.QueryStatsResponse{
		FlexibleTotalPrincipal:                   stats.FlexibleTotalPrincipal,
		FlexibleTotalUsers:                       stats.FlexibleTotalUsers,
		FlexibleTotalDistributedRewardsPrincipal: stats.FlexibleTotalDistributedRewardsPrincipal,
		StakedTotalPrincipal:                     stats.StakedTotalPrincipal,
		StakedTotalUsers:                         stats.StakedTotalUsers,
	}, nil
}
