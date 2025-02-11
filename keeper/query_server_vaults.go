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
		Paused: k.GetPaused(ctx),
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

	positions, err := k.GetPositionsByProvider(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &vaults.QueryPositionsByProviderResponse{
		Positions: positions,
	}, nil
}

func (k vaultsQueryServer) Stats(ctx context.Context, req *vaults.QueryStats) (*vaults.QueryStatsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stats, err := k.GetVaultsStats(ctx)
	if err != nil {
		return nil, err
	}

	return &vaults.QueryStatsResponse{
		FlexibleTotalPrincipal:                   stats.FlexibleTotalPrincipal,
		FlexibleTotalUsers:                       stats.FlexibleTotalUsers,
		FlexibleTotalDistributedRewardsPrincipal: stats.FlexibleTotalDistributedRewardsPrincipal,
		StakedTotalPrincipal:                     stats.StakedTotalPrincipal,
		StakedTotalUsers:                         stats.StakedTotalUsers,
	}, nil
}
