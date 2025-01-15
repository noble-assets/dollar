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

func (k vaultsQueryServer) Owner(ctx context.Context, req *vaults.QueryOwner) (*vaults.QueryOwnerResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	owner, err := k.Keeper.Owner.Get(ctx)

	return &vaults.QueryOwnerResponse{Owner: owner}, err
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
