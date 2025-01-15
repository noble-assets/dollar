package vaults

import "cosmossdk.io/errors"

var (
	ErrNoOwner          = errors.Register(SubmoduleName, 1, "there is no owner")
	ErrNotOwner         = errors.Register(SubmoduleName, 2, "signer is not owner")
	ErrInvalidAmount    = errors.Register(SubmoduleName, 3, "invalid amount")
	ErrInvalidVaultType = errors.Register(SubmoduleName, 4, "invalid vault type")
	ErrInvalidPauseType = errors.Register(SubmoduleName, 5, "invalid pause type")
	ErrActionPaused     = errors.Register(SubmoduleName, 6, "action is paused")
)
