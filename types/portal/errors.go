package portal

import "cosmossdk.io/errors"

var (
	ErrNoOwner          = errors.Register(SubmoduleName, 1, "there is no owner")
	ErrNotOwner         = errors.Register(SubmoduleName, 2, "signer is not owner")
	ErrSameOwner        = errors.Register(SubmoduleName, 3, "provided owner is the current owner")
	ErrInvalidPeer      = errors.Register(SubmoduleName, 4, "invalid peer")
	ErrInvalidMessage   = errors.Register(SubmoduleName, 5, "invalid message")
	ErrInvalidRecipient = errors.Register(SubmoduleName, 6, "invalid recipient")
)
