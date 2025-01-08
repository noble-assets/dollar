package keeper

import (
	"context"

	"cosmossdk.io/math"
)

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
