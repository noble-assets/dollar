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

package dollar

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"dollar.noble.xyz/types/vaults"

	"dollar.noble.xyz/keeper"
	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/portal"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, address address.Codec, genesis types.GenesisState) {
	var err error

	err = k.Index.Set(ctx, genesis.Index)
	if err != nil {
		panic(errors.Wrap(err, "unable to set genesis index"))
	}

	for rawAccount, rawPrincipal := range genesis.Principal {
		account, err := address.StringToBytes(rawAccount)
		if err != nil {
			panic(errors.Wrapf(err, "unable to decode account %s", rawAccount))
		}

		principal, ok := math.NewIntFromString(rawPrincipal)
		if !ok {
			panic(fmt.Errorf("unable to parse principal %s", rawPrincipal))
		}

		err = k.Principal.Set(ctx, account, principal)
		if err != nil {
			panic(errors.Wrapf(err, "unable to set genesis principal (%s:%s)", rawAccount, rawPrincipal))
		}
	}

	err = k.Stats.Set(ctx, genesis.Stats)
	if err != nil {
		panic(errors.Wrap(err, "unable to set genesis stats"))
	}

	if err = k.Owner.Set(ctx, genesis.Portal.Owner); err != nil {
		panic(errors.Wrap(err, "unable to set genesis owner"))
	}

	for chain, peer := range genesis.Portal.Peers {
		err = k.Peers.Set(ctx, chain, peer)
		if err != nil {
			panic(errors.Wrapf(err, "unable to set genesis peer (%d:%s)", chain, peer))
		}
	}

	if err = k.Nonce.Set(ctx, genesis.Portal.Nonce); err != nil {
		panic(errors.Wrap(err, "unable to set genesis nonce"))
	}

	for _, position := range genesis.Vaults.Positions {
		if err = k.Positions.Set(ctx, collections.Join3(position.Address, int32(position.Vault), position.Time.Unix()), vaults.Position{
			Principal: position.Principal,
			Index:     position.Index,
			Amount:    position.Amount,
			Time:      position.Time,
		}); err != nil {
			panic(errors.Wrapf(err, "unable to set position (%s:%s)", position.Address, position.Vault))
		}
	}

	for _, reward := range genesis.Vaults.Rewards {
		if err = k.Rewards.Set(ctx, reward.Index.String(), vaults.Reward{
			Index:   reward.Index,
			Total:   reward.Total,
			Rewards: reward.Rewards,
		}); err != nil {
			panic(errors.Wrapf(err, "unable to set reward (index:%s)", reward.Index))
		}
	}

	if err = k.VaultsStats.Set(ctx, genesis.Vaults.Stats); err != nil {
		panic(errors.Wrapf(err, "unable to set genesis vaults stats"))
	}

	if err = k.Paused.Set(ctx, int32(genesis.Vaults.Paused)); err != nil {
		panic(errors.Wrap(err, "unable to set genesis vaults paused"))
	}

	if err = k.TotalFlexiblePrincipal.Set(ctx, genesis.Vaults.TotalFlexiblePrincipal); err != nil {
		panic(errors.Wrap(err, "unable to set total flexible principal"))
	}
}

func ExportGenesis(ctx context.Context, k *keeper.Keeper) *types.GenesisState {
	index, _ := k.Index.Get(ctx)
	principal, _ := k.GetPrincipal(ctx)
	stats, _ := k.Stats.Get(ctx)

	owner, _ := k.Owner.Get(ctx)
	peers, _ := k.GetPeers(ctx)
	nonce, _ := k.Nonce.Get(ctx)

	rewards, _ := k.GetRewards(ctx)
	positions, _ := k.GetPositions(ctx)
	totalFlexiblePrincipal, _ := k.GetTotalFlexiblePrincipal(ctx)
	paused := k.GetPaused(ctx)
	vaultsStats, _ := k.GetVaultsStats(ctx)

	return &types.GenesisState{
		Portal: portal.GenesisState{
			Owner: owner,
			Peers: peers,
			Nonce: nonce,
		},
		Vaults: vaults.GenesisState{
			Positions:              positions,
			Rewards:                rewards,
			TotalFlexiblePrincipal: totalFlexiblePrincipal,
			Paused:                 paused,
			Stats:                  vaultsStats,
		},
		Index:     index,
		Principal: principal,
		Stats:     stats,
	}
}
