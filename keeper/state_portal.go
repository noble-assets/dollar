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

	"dollar.noble.xyz/types/portal"
)

// GetPeers is a utility that returns all peers from state.
func (k *Keeper) GetPeers(ctx context.Context) (map[uint16]portal.Peer, error) {
	peers := make(map[uint16]portal.Peer)

	err := k.Peers.Walk(ctx, nil, func(chain uint16, peer portal.Peer) (stop bool, err error) {
		peers[chain] = peer
		return false, nil
	})

	return peers, err
}

// IncrementNonce is a utility that returns the next nonce and increments.
func (k *Keeper) IncrementNonce(ctx context.Context) (uint32, error) {
	nonce, err := k.Nonce.Get(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get nonce from state")
	}

	err = k.Nonce.Set(ctx, nonce+1)
	if err != nil {
		return 0, errors.Wrap(err, "unable to set nonce in state")
	}

	return nonce, nil
}
