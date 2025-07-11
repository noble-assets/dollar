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

	"dollar.noble.xyz/v2/types"
	"dollar.noble.xyz/v2/types/portal"
)

var _ portal.QueryServer = &portalQueryServer{}

type portalQueryServer struct {
	*Keeper
}

func NewPortalQueryServer(keeper *Keeper) portal.QueryServer {
	return &portalQueryServer{Keeper: keeper}
}

func (k portalQueryServer) Owner(ctx context.Context, req *portal.QueryOwner) (*portal.QueryOwnerResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	owner, err := k.PortalOwner.Get(ctx)

	return &portal.QueryOwnerResponse{Owner: owner}, err
}

func (k portalQueryServer) Paused(ctx context.Context, req *portal.QueryPaused) (*portal.QueryPausedResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	return &portal.QueryPausedResponse{
		Paused: k.GetPortalPaused(ctx),
	}, nil
}

func (k portalQueryServer) Peers(ctx context.Context, req *portal.QueryPeers) (*portal.QueryPeersResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	peers, err := k.GetPortalPeers(ctx)

	return &portal.QueryPeersResponse{Peers: peers}, err
}

func (k portalQueryServer) DestinationTokens(ctx context.Context, req *portal.QueryDestinationTokens) (*portal.QueryDestinationTokensResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	// NOTE: We have to typecast here because protoc-gen-grpc-gateway doesn't support it via gogoproto.
	chainID := uint16(req.ChainId)

	var destinationTokens [][]byte
	err := k.PortalBridgingPaths.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint16, []byte](chainID),
		func(key collections.Pair[uint16, []byte], supported bool) (stop bool, err error) {
			if supported {
				destinationTokens = append(destinationTokens, key.K2())
			}

			return false, nil
		},
	)

	return &portal.QueryDestinationTokensResponse{DestinationTokens: destinationTokens}, err
}

func (k portalQueryServer) Nonce(ctx context.Context, req *portal.QueryNonce) (*portal.QueryNonceResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	nonce, err := k.PortalNonce.Get(ctx)

	return &portal.QueryNonceResponse{Nonce: nonce}, err
}
