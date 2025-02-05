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

	"dollar.noble.xyz/types/vaults"
)

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

func (k *Keeper) GetTotalFlexiblePrincipal(ctx context.Context) (math.Int, error) {
	value, err := k.TotalFlexiblePrincipal.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	return value, nil
}

func (k *Keeper) GetPaused(ctx context.Context) vaults.PausedType {
	value, err := k.Paused.Get(ctx)
	if err != nil {
		return vaults.NONE
	}
	return vaults.PausedType(value)
}

func (k *Keeper) GetPositions(ctx context.Context) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.Positions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		position, _ := k.Positions.Get(ctx, key)
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

func (k *Keeper) GetPositionsByProvider(ctx context.Context, provider []byte) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.Positions.Indexes.ByProvider.MatchExact(ctx, provider)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		position, _ := k.Positions.Get(ctx, key)
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

func (k *Keeper) GetRewards(ctx context.Context) ([]vaults.Reward, error) {
	var rewards []vaults.Reward

	itr, err := k.Rewards.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		reward, _ := k.Rewards.Get(ctx, key)
		rewards = append(rewards, vaults.Reward{
			Index:   reward.Index,
			Total:   reward.Total,
			Rewards: reward.Rewards,
		})
	}

	return rewards, err
}
