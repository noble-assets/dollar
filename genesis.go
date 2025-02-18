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

	err = k.Paused.Set(ctx, genesis.Paused)
	if err != nil {
		panic(errors.Wrap(err, "unable to set genesis paused state"))
	}

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

	if err = k.PortalOwner.Set(ctx, genesis.Portal.Owner); err != nil {
		panic(errors.Wrap(err, "unable to set genesis portal owner"))
	}

	if err = k.PortalPaused.Set(ctx, genesis.Portal.Paused); err != nil {
		panic(errors.Wrap(err, "unable to set genesis portal paused state"))
	}

	for chain, peer := range genesis.Portal.Peers {
		err = k.PortalPeers.Set(ctx, chain, peer)
		if err != nil {
			panic(errors.Wrapf(err, "unable to set genesis portal peer (%d:%s)", chain, peer))
		}
	}

	if err = k.PortalNonce.Set(ctx, genesis.Portal.Nonce); err != nil {
		panic(errors.Wrap(err, "unable to set genesis portal nonce"))
	}

	for _, position := range genesis.Vaults.Positions {
		if err = k.VaultsPositions.Set(ctx, collections.Join3(position.Address, int32(position.Vault), position.Time.Unix()), vaults.Position{
			Principal: position.Principal,
			Index:     position.Index,
			Amount:    position.Amount,
			Time:      position.Time,
		}); err != nil {
			panic(errors.Wrapf(err, "unable to set vaults position (%s:%s)", position.Address, position.Vault))
		}
	}

	for _, reward := range genesis.Vaults.Rewards {
		if err = k.VaultsRewards.Set(ctx, reward.Index, vaults.Reward{
			Index:   reward.Index,
			Total:   reward.Total,
			Rewards: reward.Rewards,
		}); err != nil {
			panic(errors.Wrapf(err, "unable to set vaults reward (index:%d)", reward.Index))
		}
	}

	if err = k.VaultsPaused.Set(ctx, int32(genesis.Vaults.Paused)); err != nil {
		panic(errors.Wrap(err, "unable to set genesis vaults paused state"))
	}

	if err = k.VaultsStats.Set(ctx, genesis.Vaults.Stats); err != nil {
		panic(errors.Wrapf(err, "unable to set genesis vaults stats"))
	}

	if err = k.VaultsTotalFlexiblePrincipal.Set(ctx, genesis.Vaults.TotalFlexiblePrincipal); err != nil {
		panic(errors.Wrap(err, "unable to set total vaults flexible principal"))
	}
}

func ExportGenesis(ctx context.Context, k *keeper.Keeper) *types.GenesisState {
	paused := k.GetPaused(ctx)
	index, _ := k.Index.Get(ctx)
	principal, _ := k.GetPrincipal(ctx)
	stats, _ := k.Stats.Get(ctx)

	portalOwner, _ := k.PortalOwner.Get(ctx)
	portalPaused := k.GetPortalPaused(ctx)
	portalPeers, _ := k.GetPortalPeers(ctx)
	portalNonce, _ := k.PortalNonce.Get(ctx)

	vaultsRewards, _ := k.GetVaultsRewards(ctx)
	vaultsPositions, _ := k.GetVaultsPositions(ctx)
	vaultsTotalFlexiblePrincipal, _ := k.GetVaultsTotalFlexiblePrincipal(ctx)
	vaultsPaused := k.GetVaultsPaused(ctx)
	vaultsStats, _ := k.GetVaultsStats(ctx)

	return &types.GenesisState{
		Portal: portal.GenesisState{
			Owner:  portalOwner,
			Paused: portalPaused,
			Peers:  portalPeers,
			Nonce:  portalNonce,
		},
		Vaults: vaults.GenesisState{
			Positions:              vaultsPositions,
			Rewards:                vaultsRewards,
			TotalFlexiblePrincipal: vaultsTotalFlexiblePrincipal,
			Paused:                 vaultsPaused,
			Stats:                  vaultsStats,
		},
		Paused:    paused,
		Index:     index,
		Principal: principal,
		Stats:     stats,
	}
}
