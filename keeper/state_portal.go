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
