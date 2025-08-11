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
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/math"

	"dollar.noble.xyz/v3/types/vaults"
)

type VaultsPositionsIndexes struct {
	ByProvider         *indexes.Multi[[]byte, collections.Triple[[]byte, int32, int64], vaults.Position]
	ByProviderAndVault *indexes.Multi[collections.Pair[[]byte, int32], collections.Triple[[]byte, int32, int64], vaults.Position]
}

func (i VaultsPositionsIndexes) IndexesList() []collections.Index[collections.Triple[[]byte, int32, int64], vaults.Position] {
	return []collections.Index[collections.Triple[[]byte, int32, int64], vaults.Position]{
		i.ByProvider,
		i.ByProviderAndVault,
	}
}

func NewVaultsPositionsIndexes(builder *collections.SchemaBuilder) VaultsPositionsIndexes {
	return VaultsPositionsIndexes{
		ByProvider: indexes.NewMulti(
			builder, []byte("vaults_positions_by_provider"), "vaults_positions_by_provider",
			collections.BytesKey,
			collections.TripleKeyCodec(collections.BytesKey, collections.Int32Key, collections.Int64Key),
			func(key collections.Triple[[]byte, int32, int64], value vaults.Position) ([]byte, error) {
				return key.K1(), nil
			},
		),
		ByProviderAndVault: indexes.NewMulti(
			builder, []byte("vaults_positions_by_pair_provider_and_vault"), "vaults_positions_by_pair_provider_and_vault",
			collections.PairKeyCodec(collections.BytesKey, collections.Int32Key),
			collections.TripleKeyCodec(collections.BytesKey, collections.Int32Key, collections.Int64Key),
			func(key collections.Triple[[]byte, int32, int64], value vaults.Position) (collections.Pair[[]byte, int32], error) {
				return collections.Join(key.K1(), key.K2()), nil
			},
		),
	}
}

//

// GetVaultsTotalFlexiblePrincipal is a utility that returns the total flexible vault principal from state.
func (k *Keeper) GetVaultsTotalFlexiblePrincipal(ctx context.Context) (math.Int, error) {
	value, err := k.VaultsTotalFlexiblePrincipal.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	return value, nil
}

// GetVaultsPaused is a utility that returns the current paused state.
func (k *Keeper) GetVaultsPaused(ctx context.Context) vaults.PausedType {
	value, err := k.VaultsPaused.Get(ctx)
	if err != nil {
		return vaults.NONE
	}
	return vaults.PausedType(value)
}

// IsVaultsSeasonOneEnded is a utility that returns if Vaults Season One has ended.
func (k *Keeper) IsVaultsSeasonOneEnded(ctx context.Context) bool {
	ended, _ := k.VaultsSeasonOneEnded.Get(ctx)
	return ended
}

// GetVaultsPositions is a utility that returns all vaults positions from state.
func (k *Keeper) GetVaultsPositions(ctx context.Context) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.VaultsPositions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		position, _ := k.VaultsPositions.Get(ctx, key)
		positions = append(positions, vaults.PositionEntry{
			Address:   key.K1(),
			Vault:     vaults.VaultType(key.K2()),
			Index:     position.Index,
			Principal: position.Principal,
			Amount:    position.Amount,
			Time:      position.Time,
		})
	}

	return positions, err
}

// GetVaultsPositionsByProvider is a utility that returns all vaults positions from state by a given provider.
func (k *Keeper) GetVaultsPositionsByProvider(ctx context.Context, provider []byte) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.VaultsPositions.Indexes.ByProvider.MatchExact(ctx, provider)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		position, _ := k.VaultsPositions.Get(ctx, key)
		positions = append(positions, vaults.PositionEntry{
			Address:   key.K1(),
			Vault:     vaults.VaultType(key.K2()),
			Index:     position.Index,
			Principal: position.Principal,
			Amount:    position.Amount,
			Time:      position.Time,
		})
	}

	return positions, err
}

// GetVaultsPositionsByProviderAndVault is a utility that returns all vaults positions from state by a given provider and vault.
func (k *Keeper) GetVaultsPositionsByProviderAndVault(ctx context.Context, provider []byte, vault int32) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.VaultsPositions.Indexes.ByProviderAndVault.MatchExact(ctx, collections.Join(provider, vault))
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		position, _ := k.VaultsPositions.Get(ctx, key)
		positions = append(positions, vaults.PositionEntry{
			Address:   key.K1(),
			Vault:     vaults.VaultType(key.K2()),
			Index:     position.Index,
			Principal: position.Principal,
			Amount:    position.Amount,
			Time:      position.Time,
		})
	}

	return positions, err
}

// GetVaultsRewards is a utility that returns all rewards positions from state.
func (k *Keeper) GetVaultsRewards(ctx context.Context) ([]vaults.Reward, error) {
	var rewards []vaults.Reward

	itr, err := k.VaultsRewards.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		reward, _ := k.VaultsRewards.Get(ctx, key)
		rewards = append(rewards, vaults.Reward{
			Index:   reward.Index,
			Total:   reward.Total,
			Rewards: reward.Rewards,
		})
	}

	return rewards, err
}

// GetVaultsStats is a utility that returns all vaults stats from state.
func (k *Keeper) GetVaultsStats(ctx context.Context) (vaults.Stats, error) {
	stats, err := k.VaultsStats.Get(ctx)
	if err != nil {
		return vaults.Stats{}, err
	}
	return stats, nil
}

// IncrementVaultUsers is a utility that increments the total vault users stat.
func (k *Keeper) IncrementVaultUsers(ctx context.Context, vault vaults.VaultType) error {
	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return err
	}

	switch vault {
	case vaults.STAKED:
		stats.StakedTotalUsers += 1
	case vaults.FLEXIBLE:
		stats.FlexibleTotalUsers += 1
	}

	return k.VaultsStats.Set(ctx, stats)
}

// DecrementVaultUsers is a utility that decrements the total vault users stat.
func (k *Keeper) DecrementVaultUsers(ctx context.Context, vault vaults.VaultType) error {
	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return err
	}

	switch vault {
	case vaults.STAKED:
		stats.StakedTotalUsers -= 1
	case vaults.FLEXIBLE:
		stats.FlexibleTotalUsers -= 1
	}

	return k.VaultsStats.Set(ctx, stats)
}

// IncrementVaultTotalPrincipal is a utility that increments the total vault principal stat.
func (k *Keeper) IncrementVaultTotalPrincipal(ctx context.Context, vault vaults.VaultType, amount math.Int) error {
	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return err
	}

	switch vault {
	case vaults.STAKED:
		stats.StakedTotalPrincipal = stats.StakedTotalPrincipal.Add(amount)
	case vaults.FLEXIBLE:
		stats.FlexibleTotalPrincipal = stats.FlexibleTotalPrincipal.Add(amount)
	}

	return k.VaultsStats.Set(ctx, stats)
}

// DecrementVaultTotalPrincipal is a utility that decrements the total vault principal stat.
func (k *Keeper) DecrementVaultTotalPrincipal(ctx context.Context, vault vaults.VaultType, amount math.Int) error {
	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return err
	}

	switch vault {
	case vaults.STAKED:
		stats.StakedTotalPrincipal = stats.StakedTotalPrincipal.Sub(amount)
	case vaults.FLEXIBLE:
		stats.FlexibleTotalPrincipal = stats.FlexibleTotalPrincipal.Sub(amount)
	}

	return k.VaultsStats.Set(ctx, stats)
}

// IncrementFlexibleTotalDistributedRewardsPrincipal is a utility that increments the total flexible vault distributed principal rewards stat.
func (k *Keeper) IncrementFlexibleTotalDistributedRewardsPrincipal(ctx context.Context, amount math.Int) error {
	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return err
	}

	if !stats.FlexibleTotalDistributedRewardsPrincipal.IsNil() {
		stats.FlexibleTotalDistributedRewardsPrincipal = stats.FlexibleTotalDistributedRewardsPrincipal.Add(amount)
	} else {
		stats.FlexibleTotalDistributedRewardsPrincipal = amount
	}

	return k.VaultsStats.Set(ctx, stats)
}
