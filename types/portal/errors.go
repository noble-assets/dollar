package portal

import "cosmossdk.io/errors"

var (
	ErrNoOwner          = errors.Register(SubmoduleName, 1, "there is no owner")
	ErrNotOwner         = errors.Register(SubmoduleName, 2, "signer is not owner")
	ErrInvalidPeer      = errors.Register(SubmoduleName, 3, "invalid peer")
	ErrInvalidMessage   = errors.Register(SubmoduleName, 4, "invalid message")
	ErrInvalidRecipient = errors.Register(SubmoduleName, 5, "invalid recipient")
)
