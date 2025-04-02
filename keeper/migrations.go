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
	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/types"
	"dollar.noble.xyz/v2/types/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper  *Keeper
	v1Stats collections.Item[types.Stats]
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper *Keeper) Migrator {
	builder := collections.NewSchemaBuilder(keeper.store)

	return Migrator{
		keeper:  keeper,
		v1Stats: collections.NewItem(builder, types.StatsKey, "stats", codec.CollValue[types.Stats](keeper.cdc)),
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	v1Stats, err := m.v1Stats.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get noble dollar v1 stats")
	}

	stats := v2.Stats{
		TotalHolders:       v1Stats.TotalHolders,
		TotalPrincipal:     v1Stats.TotalPrincipal,
		TotalYieldAccrued:  v1Stats.TotalYieldAccrued,
		TotalExternalYield: make(map[string]string),
	}

	err = m.keeper.Stats.Set(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "failed to set noble dollar v2 stats")
	}

	return nil
}
