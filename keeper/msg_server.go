package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) ClaimYield(ctx context.Context, msg *types.MsgClaimYield) (*types.MsgClaimYieldResponse, error) {
	yield, account, err := k.GetYield(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, account, sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
	if err != nil {
		// TODO(@john): Wrap error for developer friendliness!
		return nil, err
	}

	return &types.MsgClaimYieldResponse{}, nil
}

func (k *Keeper) Mint(ctx context.Context, recipient []byte, amount math.Int) error {
	coins := sdk.NewCoins(sdk.NewCoin(k.denom, amount))
	err := k.bank.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		// TODO(@john): Wrap error for developer friendliness!
		return err
	}
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		// TODO(@john): Wrap error for developer friendliness!
		return err
	}

	return nil
}

func (k *Keeper) UpdateIndex(ctx context.Context, rawIndex int64) error {
	oldIndex, err := k.Index.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get index from state")
	}
	if rawIndex <= oldIndex {
		return types.ErrDecreasingIndex
	}

	err = k.Index.Set(ctx, rawIndex)
	if err != nil {
		return errors.Wrap(err, "unable to set index in state")
	}

	totalPrincipal, err := k.GetTotalPrincipal(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get total principal from state")
	}

	index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)
	currentSupply := k.bank.GetSupply(ctx, k.denom).Amount
	// TODO(@john): Ensure that we're always rounding down here, to avoid minting more $USDN than underlying M.
	expectedSupply := index.MulInt(totalPrincipal).TruncateInt()

	err = k.bank.MintCoins(ctx, types.YieldName, sdk.NewCoins(sdk.NewCoin(k.denom, expectedSupply.Sub(currentSupply))))
	if err != nil {
		// TODO(@john): Wrap error for developer friendliness!
		return err
	}

	return nil
}
