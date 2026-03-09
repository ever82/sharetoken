package types

import (
	"cosmossdk.io/errors"
)

// x/oracle module sentinel errors
var (
	ErrInvalidSymbol    = errors.Register(ModuleName, 1, "invalid symbol")
	ErrPriceNotFound    = errors.Register(ModuleName, 2, "price not found")
	ErrStalePrice       = errors.Register(ModuleName, 3, "stale price")
	ErrInvalidPrice     = errors.Register(ModuleName, 4, "invalid price")
	ErrLowConfidence    = errors.Register(ModuleName, 5, "low confidence price")
)
