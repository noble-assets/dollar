package keeper

import (
	"context"

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
