package types

import (
	"fmt"

	"cosmossdk.io/errors"
)

// x/identity module sentinel errors
var (
	ErrInvalidAddress        = errors.Register(ModuleName, 1, "invalid address")
	ErrIdentityNotFound      = errors.Register(ModuleName, 2, "identity not found")
	ErrIdentityAlreadyExists = errors.Register(ModuleName, 3, "identity already exists")
	ErrInvalidVerification   = errors.Register(ModuleName, 4, "invalid verification")
	ErrDIDAlreadyRegistered  = errors.Register(ModuleName, 5, "DID already registered")
	ErrProviderAlreadyUsed   = errors.Register(ModuleName, 6, "provider already used for this account")
	ErrInvalidProvider       = errors.Register(ModuleName, 7, "invalid verification provider")
	ErrLimitExceeded         = errors.Register(ModuleName, 8, "limit exceeded")
	ErrInvalidLimitConfig    = errors.Register(ModuleName, 9, "invalid limit configuration")
	ErrUnauthorized          = errors.Register(ModuleName, 10, "unauthorized")
	ErrInvalidProof          = errors.Register(ModuleName, 11, "invalid merkle proof")
	ErrCooldownNotMet        = errors.Register(ModuleName, 12, "cooldown period not met")
	ErrMaxDisputesReached    = errors.Register(ModuleName, 13, "maximum active disputes reached")
	ErrMaxConcurrentReached  = errors.Register(ModuleName, 14, "maximum concurrent services reached")
	ErrRateLimitExceeded     = errors.Register(ModuleName, 15, "rate limit exceeded")
	ErrInvalidDID            = errors.Register(ModuleName, 16, "invalid DID format")
)

// LimitError represents a detailed limit error
type LimitError struct {
	LimitType   string
	Current     string
	Max         string
	Description string
}

func (e LimitError) Error() string {
	return fmt.Sprintf("%s limit exceeded: current %s, max %s - %s",
		e.LimitType, e.Current, e.Max, e.Description)
}
