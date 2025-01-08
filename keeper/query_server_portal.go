package keeper

import (
	"context"

	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/portal"
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

	owner, err := k.Keeper.Owner.Get(ctx)

	return &portal.QueryOwnerResponse{Owner: owner}, err
}

func (k portalQueryServer) Peers(ctx context.Context, req *portal.QueryPeers) (*portal.QueryPeersResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	peers, err := k.GetPeers(ctx)

	return &portal.QueryPeersResponse{Peers: peers}, err
}
