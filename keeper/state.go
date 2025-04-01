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
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"dollar.noble.xyz/v2/types/v2"
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

// IncrementTotalExternalYield is a utility that increments the total external yield stat.
func (k *Keeper) IncrementTotalExternalYield(ctx context.Context, provider v2.Provider, identifier string, amount math.Int) error {
	stats, err := k.Stats.Get(ctx)
	if err != nil {
		return err
	}
	if stats.TotalExternalYield == nil {
		stats.TotalExternalYield = make(map[string]string)
	}

	key := fmt.Sprintf("%s/%s", provider, identifier)

	totalExternalYield := math.ZeroInt()
	rawTotalExternalYield, exists := stats.TotalExternalYield[key]
	if exists {
		totalExternalYield, _ = math.NewIntFromString(rawTotalExternalYield)
	}

	totalExternalYield = totalExternalYield.Add(amount)
	stats.TotalExternalYield[key] = totalExternalYield.String()

	return k.Stats.Set(ctx, stats)
}

// GetYieldRecipients is a utility that returns all yield recipients from state.
func (k *Keeper) GetYieldRecipients(ctx context.Context) (map[string]string, error) {
	yieldRecipients := make(map[string]string)

	err := k.YieldRecipients.Walk(ctx, nil, func(key collections.Pair[int32, string], yieldRecipient string) (stop bool, err error) {
		yieldRecipients[fmt.Sprintf("%s/%s", v2.Provider(key.K1()), key.K2())] = yieldRecipient
		return false, nil
	})

	return yieldRecipients, err
}

// GetYieldRecipientsByProvider is utility that returns yield recipients for a specific provider.
func (k *Keeper) GetYieldRecipientsByProvider(ctx context.Context, provider v2.Provider) (map[string]string, error) {
	yieldRecipients := make(map[string]string)

	err := k.YieldRecipients.Walk(
		ctx,
		collections.NewPrefixedPairRange[int32, string](int32(provider)),
		func(key collections.Pair[int32, string], yieldRecipient string) (stop bool, err error) {
			identifier := key.K2()
			yieldRecipients[identifier] = yieldRecipient

			return false, nil
		},
	)

	return yieldRecipients, err
}

// HasYieldRecipient is a utility that returns if there is a yield recipient for a specific provider and identifier.
func (k *Keeper) HasYieldRecipient(ctx context.Context, provider v2.Provider, identifier string) bool {
	key := collections.Join(int32(provider), identifier)
	has, _ := k.YieldRecipients.Has(ctx, key)

	return has
}
