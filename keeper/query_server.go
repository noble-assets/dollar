package keeper

import (
	"context"

	"cosmossdk.io/errors"

	"dollar.noble.xyz/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
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
