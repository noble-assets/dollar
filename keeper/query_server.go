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

	"dollar.noble.xyz/v2/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Index(ctx context.Context, req *types.QueryIndex) (*types.QueryIndexResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	rawIndex, err := k.Keeper.Index.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get index from state")
	}

	index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)

	return &types.QueryIndexResponse{Index: index}, nil
}

func (k queryServer) Paused(ctx context.Context, req *types.QueryPaused) (*types.QueryPausedResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	return &types.QueryPausedResponse{
		Paused: k.GetPaused(ctx),
	}, nil
}

func (k queryServer) Principal(ctx context.Context, req *types.QueryPrincipal) (*types.QueryPrincipalResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	account, err := k.address.StringToBytes(req.Account)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode account %s", req.Account)
	}

	principal, _ := k.Keeper.Principal.Get(ctx, account)

	return &types.QueryPrincipalResponse{Principal: principal}, nil
}

func (k queryServer) Yield(ctx context.Context, req *types.QueryYield) (*types.QueryYieldResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	yield, _, err := k.GetYield(ctx, req.Account)

	return &types.QueryYieldResponse{ClaimableAmount: yield}, err
}

func (k queryServer) Stats(ctx context.Context, req *types.QueryStats) (*types.QueryStatsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stats, err := k.Keeper.Stats.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get stats from state")
	}

	return &types.QueryStatsResponse{
		TotalHolders:      stats.TotalHolders,
		TotalPrincipal:    stats.TotalPrincipal,
		TotalYieldAccrued: stats.TotalYieldAccrued,
	}, nil
}
