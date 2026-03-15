package types

import "cosmossdk.io/errors"

var (
	ErrInvalidService = errors.Register(ModuleName, 1, "invalid service")
	ErrServiceNotFound = errors.Register(ModuleName, 2, "service not found")
	ErrInvalidParams   = errors.Register(ModuleName, 3, "invalid params")
)
