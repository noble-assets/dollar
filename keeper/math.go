package keeper

import "cosmossdk.io/math"

// GetPrincipalAmountRoundedUp returns the rounded up principal given a present amount.
//
// https://github.com/m0-foundation/protocol/blob/b1c6e624ed09a9e28f4ae45cd87fda610fafe446/src/abstract/ContinuousIndexing.sol#L106-L114
func (k *Keeper) GetPrincipalAmountRoundedUp(presentAmount math.Int, index int64) (principalAmount math.Int) {
	return presentAmount.MulRaw(1e12).AddRaw(index).SubRaw(1).QuoRaw(index)
}

// GetPrincipalAmountRoundedDown returns the rounded down principal given a present amount.
//
// https://github.com/m0-foundation/protocol/blob/b1c6e624ed09a9e28f4ae45cd87fda610fafe446/src/abstract/ContinuousIndexing.sol#L96-L104
func (k *Keeper) GetPrincipalAmountRoundedDown(presentAmount math.Int, index int64) (principalAmount math.Int) {
	return presentAmount.MulRaw(1e12).QuoRaw(index)
}

// GetPresentAmount returns the rounded down present amount given a principal.
//
// https://github.com/m0-foundation/protocol/blob/b1c6e624ed09a9e28f4ae45cd87fda610fafe446/src/abstract/ContinuousIndexing.sol#L76-L84
func (k *Keeper) GetPresentAmount(principalAmount math.Int, index int64) (presentAmount math.Int) {
	return principalAmount.MulRaw(index).QuoRaw(1e12)
}
