package types

import (
	"cosmossdk.io/errors"
)

// x/escrow module sentinel errors
var (
	ErrInvalidEscrowID      = errors.Register(ModuleName, 1, "invalid escrow ID")
	ErrEscrowNotFound       = errors.Register(ModuleName, 2, "escrow not found")
	ErrEscrowAlreadyExists  = errors.Register(ModuleName, 3, "escrow already exists")
	ErrInvalidAmount        = errors.Register(ModuleName, 4, "invalid amount")
	ErrInvalidStatus        = errors.Register(ModuleName, 5, "invalid escrow status")
	ErrUnauthorized         = errors.Register(ModuleName, 6, "unauthorized")
	ErrEscrowExpired        = errors.Register(ModuleName, 7, "escrow expired")
	ErrEscrowNotExpired     = errors.Register(ModuleName, 8, "escrow not expired")
	ErrInsufficientFunds    = errors.Register(ModuleName, 9, "insufficient funds")
	ErrInvalidCompletion    = errors.Register(ModuleName, 10, "invalid completion proof")
	ErrDisputeAlreadyExists = errors.Register(ModuleName, 11, "dispute already exists")
	ErrEscrowNotDisputed    = errors.Register(ModuleName, 12, "escrow not disputed")
	ErrInvalidAllocation    = errors.Register(ModuleName, 13, "invalid fund allocation")
	ErrInvalidRequester     = errors.Register(ModuleName, 14, "invalid requester")
	ErrInvalidProvider      = errors.Register(ModuleName, 15, "invalid provider")
)
