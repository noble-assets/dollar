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

// DecrementTotalHolders is a utility that decrements the total holders stat.
func (k *Keeper) DecrementTotalHolders(ctx context.Context) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}

	if stats.TotalHolders > 0 {
		stats.TotalHolders -= 1
	}

	return k.Stats.Set(ctx, stats)
}

// IncrementTotalHolders is a utility that increments the total holders stat.
func (k *Keeper) IncrementTotalHolders(ctx context.Context) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}

	stats.TotalHolders += 1

	return k.Stats.Set(ctx, stats)
}

// GetTotalPrincipal is a utility that returns the total principal stat.
func (k *Keeper) GetTotalPrincipal(ctx context.Context) (math.Int, error) {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	return stats.TotalPrincipal, nil
}

// DecrementTotalPrincipal is a utility that decrements the total principal stat.
func (k *Keeper) DecrementTotalPrincipal(ctx context.Context, amount math.Int) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}

	stats.TotalPrincipal, err = stats.TotalPrincipal.SafeSub(amount)
	if err != nil {
		return err
	}

	return k.Stats.Set(ctx, stats)
}

// IncrementTotalPrincipal is a utility that increments the total principal stat.
func (k *Keeper) IncrementTotalPrincipal(ctx context.Context, amount math.Int) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}

	stats.TotalPrincipal, err = stats.TotalPrincipal.SafeAdd(amount)
	if err != nil {
		return err
	}

	return k.Stats.Set(ctx, stats)
}

// IncrementTotalYieldAccrued is a utility that increments the total yield accrued stat.
func (k *Keeper) IncrementTotalYieldAccrued(ctx context.Context, amount math.Int) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}

	stats.TotalYieldAccrued, err = stats.TotalYieldAccrued.SafeAdd(amount)
	if err != nil {
		return err
	}

	return k.Stats.Set(ctx, stats)
}

// IncrementTotalChannelYield is a utility that increments the total channel yield stat.
func (k *Keeper) IncrementTotalChannelYield(ctx context.Context, channelId string, amount math.Int) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}
	if stats.TotalChannelYield == nil {
		stats.TotalChannelYield = make(map[string]string)
	}

	totalChannelYield := math.ZeroInt()
	rawTotalChannelYield, exists := stats.TotalChannelYield[channelId]
	if exists {
		totalChannelYield, _ = math.NewIntFromString(rawTotalChannelYield)
	}

	totalChannelYield = totalChannelYield.Add(amount)
	stats.TotalChannelYield[channelId] = totalChannelYield.String()

	return k.Stats.Set(ctx, stats)
}

// GetYieldRecipients is a utility that returns all yield recipients from state.
func (k *Keeper) GetYieldRecipients(ctx context.Context) (map[string]string, error) {
	yieldRecipients := make(map[string]string)

	err := k.YieldRecipients.Walk(ctx, nil, func(channelId string, yieldRecipient string) (stop bool, err error) {
		yieldRecipients[channelId] = yieldRecipient
		return false, nil
	})

	return yieldRecipients, err
}
