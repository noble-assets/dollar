package types

import "cosmossdk.io/errors"

var (
	ErrInvalidRequest = errors.Register(ModuleName, 0, "invalid request")

	ErrDecreasingIndex = errors.Register(ModuleName, 1, "decreasing index")
)
