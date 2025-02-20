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
