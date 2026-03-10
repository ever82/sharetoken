package types

import (
	"cosmossdk.io/errors"
)

const (
	// ModuleName is the name of the trust module
	ModuleName = "trust"
	StoreKey   = ModuleName
)

// MQ errors
var (
	ErrMQNotFound       = errors.Register(ModuleName, 1, "MQ not found")
	ErrInvalidMQ        = errors.Register(ModuleName, 2, "invalid MQ value")
	ErrMQAlreadyExists  = errors.Register(ModuleName, 3, "MQ already exists")
	ErrMQOverflow       = errors.Register(ModuleName, 4, "MQ overflow")
	ErrInvalidDisputeID = errors.Register(ModuleName, 5, "invalid dispute ID")
	ErrDisputeNotFound  = errors.Register(ModuleName, 6, "dispute not found")
)
