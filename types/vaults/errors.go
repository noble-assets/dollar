package vaults

import "cosmossdk.io/errors"

var (
	ErrInvalidAuthority = errors.Register(SubmoduleName, 1, "invalid authority")
	ErrInvalidAmount    = errors.Register(SubmoduleName, 2, "invalid amount")
	ErrInvalidVaultType = errors.Register(SubmoduleName, 3, "invalid vault type")
	ErrInvalidPauseType = errors.Register(SubmoduleName, 4, "invalid pause type")
	ErrActionPaused     = errors.Register(SubmoduleName, 5, "action is paused")
)
