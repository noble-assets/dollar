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

	"dollar.noble.xyz/v3/types/portal"
)

// GetPortalPaused is a utility that returns the current paused state.
func (k *Keeper) GetPortalPaused(ctx context.Context) bool {
	paused, _ := k.PortalPaused.Get(ctx)
	return paused
}

// GetPortalPeers is a utility that returns all peers from state.
func (k *Keeper) GetPortalPeers(ctx context.Context) (map[uint16]portal.Peer, error) {
	peers := make(map[uint16]portal.Peer)

	err := k.PortalPeers.Walk(ctx, nil, func(chain uint16, peer portal.Peer) (stop bool, err error) {
		peers[chain] = peer
		return false, nil
	})

	return peers, err
}

// GetPortalBridgingPaths is a utility that returns all supported bridging paths from state.
func (k *Keeper) GetPortalBridgingPaths(ctx context.Context) ([]portal.BridgingPath, error) {
	var supportedBridgingPaths []portal.BridgingPath

	err := k.PortalBridgingPaths.Walk(
		ctx,
		nil,
		func(key collections.Pair[uint16, []byte], supported bool) (stop bool, err error) {
			if supported {
				supportedBridgingPaths = append(supportedBridgingPaths, portal.BridgingPath{
					DestinationChainId: key.K1(),
					DestinationToken:   key.K2(),
				})
			}
			return false, nil
		},
	)

	return supportedBridgingPaths, err
}

// IncrementPortalNonce is a utility that returns the next nonce and increments.
func (k *Keeper) IncrementPortalNonce(ctx context.Context) (uint32, error) {
	nonce, err := k.PortalNonce.Get(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get portal nonce from state")
	}

	err = k.PortalNonce.Set(ctx, nonce+1)
	if err != nil {
		return 0, errors.Wrap(err, "unable to set portal nonce in state")
	}

	return nonce, nil
}
