package security

import "cosmossdk.io/errors"

// ModuleName is the security module name
const ModuleName = "security"

// Security sentinel errors
var (
	ErrRateLimitExceeded = errors.Register(ModuleName, 1, "rate limit exceeded")
	ErrInvalidSequence   = errors.Register(ModuleName, 2, "invalid sequence")
	ErrReplayDetected    = errors.Register(ModuleName, 3, "replay attack detected")
	ErrInvalidTimestamp  = errors.Register(ModuleName, 4, "invalid timestamp")
	ErrUnauthorized      = errors.Register(ModuleName, 5, "unauthorized")
	ErrInvalidInput      = errors.Register(ModuleName, 6, "invalid input")
	ErrSuspicious        = errors.Register(ModuleName, 7, "suspicious activity detected")
)
