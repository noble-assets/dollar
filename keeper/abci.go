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

	"github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of each block.
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			k.logger.Error("recovered panic from x/dollar BeginBlocker while ending the vaults program", "err", r)
			return
		}
	}()

	// If the current time exceeds the Season One end time, and it hasn't ended, handle it.
	if k.header.GetHeaderInfo(ctx).Time.Unix() > k.vaultsSeasonOneEndTimestamp && !k.IsVaultsSeasonOneEnded(ctx) {
		defer func() {
			// No matter the result of the execution, Season One must be marked
			// as completed.
			k.VaultsSeasonOneEnded.Set(ctx, true)
		}()

		// Create a cached context for the execution.
		cachedCtx, commit := types.UnwrapSDKContext(ctx).CacheContext()

		if err := k.endVaultsSeasonOne(cachedCtx); err != nil {
			return err
		}

		// Commit the results.
		commit()
	}

	return nil
}
