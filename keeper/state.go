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

	"cosmossdk.io/math"
)

// GetPaused is a utility that returns the current paused state.
func (k *Keeper) GetPaused(ctx context.Context) bool {
	paused, _ := k.Paused.Get(ctx)
	return paused
}

// GetPrincipal is a utility that returns all principal entries from state.
func (k *Keeper) GetPrincipal(ctx context.Context) (map[string]string, error) {
	principal := make(map[string]string)

	err := k.Principal.Walk(ctx, nil, func(key []byte, value math.Int) (stop bool, err error) {
		address, err := k.address.BytesToString(key)
		if err != nil {
			return false, err
		}

		principal[address] = value.String()
		return false, nil
	})

	return principal, err
}

// GetTotalPrincipal is a utility that returns the total principal from state.
func (k *Keeper) GetTotalPrincipal(ctx context.Context) (math.Int, error) {
	totalPrincipal := math.ZeroInt()

	err := k.Principal.Walk(ctx, nil, func(_ []byte, value math.Int) (stop bool, err error) {
		totalPrincipal = totalPrincipal.Add(value)
		return false, nil
	})

	return totalPrincipal, err
}
