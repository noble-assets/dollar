package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
)

// GetPrincipalAmount returns the rounded up principal given a present amount.
//
// https://github.com/m0-foundation/protocol/blob/b1c6e624ed09a9e28f4ae45cd87fda610fafe446/src/abstract/ContinuousIndexing.sol#L106-L114
func (k *Keeper) GetPrincipalAmount(ctx context.Context, presentAmount math.Int) (principalAmount math.Int, err error) {
	index, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "unable to get index from state")
	}

	principalAmount = presentAmount.MulRaw(1e12).AddRaw(index).SubRaw(1).QuoRaw(index)

	return
}

// GetPresentAmount returns the rounded down present amount given a principal.
//
// https://github.com/m0-foundation/protocol/blob/b1c6e624ed09a9e28f4ae45cd87fda610fafe446/src/abstract/ContinuousIndexing.sol#L76-L84
func (k *Keeper) GetPresentAmount(ctx context.Context, principalAmount math.Int) (presentAmount math.Int, err error) {
	index, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "unable to get index from state")
	}

	presentAmount = principalAmount.MulRaw(index).QuoRaw(1e12)

	return
}
