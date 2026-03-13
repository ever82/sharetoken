package security

import (
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultSequenceWindow is the default window for sequence number validation
	DefaultSequenceWindow = 100

	// DefaultTimestampWindow is the default time window for transaction validity (5 minutes)
	DefaultTimestampWindow = 5 * time.Minute
)

// NonceTracker tracks used nonces to prevent replay attacks
type NonceTracker struct {
	mu        sync.RWMutex
	usedNonces map[string]map[uint64]time.Time // address -> nonce -> timestamp
	window    time.Duration
}

// NewNonceTracker creates a new nonce tracker
func NewNonceTracker(window time.Duration) *NonceTracker {
	if window == 0 {
		window = DefaultTimestampWindow
	}
	return &NonceTracker{
		usedNonces: make(map[string]map[uint64]time.Time),
		window:     window,
	}
}

// CheckAndRecord checks if a nonce has been used and records it if not
// Returns true if the nonce is valid (not used before)
func (nt *NonceTracker) CheckAndRecord(address string, nonce uint64) bool {
	nt.mu.Lock()
	defer nt.mu.Unlock()

	// Clean old entries periodically
	nt.cleanOldEntries()

	if _, exists := nt.usedNonces[address]; !exists {
		nt.usedNonces[address] = make(map[uint64]time.Time)
	}

	// Check if nonce has been used
	if _, used := nt.usedNonces[address][nonce]; used {
		return false
	}

	// Record nonce
	nt.usedNonces[address][nonce] = time.Now()
	return true
}

// IsUsed checks if a nonce has been used without recording it
func (nt *NonceTracker) IsUsed(address string, nonce uint64) bool {
	nt.mu.RLock()
	defer nt.mu.RUnlock()

	if nonces, exists := nt.usedNonces[address]; exists {
		_, used := nonces[nonce]
		return used
	}
	return false
}

// cleanOldEntries removes entries older than the window
func (nt *NonceTracker) cleanOldEntries() {
	cutoff := time.Now().Add(-nt.window)
	for addr, nonces := range nt.usedNonces {
		for nonce, timestamp := range nonces {
			if timestamp.Before(cutoff) {
				delete(nonces, nonce)
			}
		}
		if len(nonces) == 0 {
			delete(nt.usedNonces, addr)
		}
	}
}

// GetUsedNonceCount returns the number of used nonces for an address
func (nt *NonceTracker) GetUsedNonceCount(address string) int {
	nt.mu.RLock()
	defer nt.mu.RUnlock()

	if nonces, exists := nt.usedNonces[address]; exists {
		return len(nonces)
	}
	return 0
}

// SequenceValidator validates account sequences to prevent replay attacks
type SequenceValidator struct {
	mu        sync.RWMutex
	sequences map[string]uint64 // address -> expected next sequence
}

// NewSequenceValidator creates a new sequence validator
func NewSequenceValidator() *SequenceValidator {
	return &SequenceValidator{
		sequences: make(map[string]uint64),
	}
}

// ValidateAndUpdate validates a sequence number and updates the expected sequence
// Returns an error if the sequence is invalid
func (sv *SequenceValidator) ValidateAndUpdate(address string, sequence uint64) error {
	sv.mu.Lock()
	defer sv.mu.Unlock()

	expected, exists := sv.sequences[address]
	if !exists {
		// First time seeing this address, accept any sequence
		sv.sequences[address] = sequence + 1
		return nil
	}

	if sequence < expected {
		return fmt.Errorf("sequence number too old: got %d, expected >= %d", sequence, expected)
	}

	// Accept the sequence and update expected
	sv.sequences[address] = sequence + 1
	return nil
}

// GetExpectedSequence returns the expected sequence for an address
func (sv *SequenceValidator) GetExpectedSequence(address string) uint64 {
	sv.mu.RLock()
	defer sv.mu.RUnlock()

	if seq, exists := sv.sequences[address]; exists {
		return seq
	}
	return 0
}

// TimestampValidator validates transaction timestamps
type TimestampValidator struct {
	window time.Duration
}

// NewTimestampValidator creates a new timestamp validator
func NewTimestampValidator(window time.Duration) *TimestampValidator {
	if window == 0 {
		window = DefaultTimestampWindow
	}
	return &TimestampValidator{
		window: window,
	}
}

// ValidateTimestamp checks if a timestamp is within the valid window
func (tv *TimestampValidator) ValidateTimestamp(txTime time.Time) error {
	now := time.Now()
	minTime := now.Add(-tv.window)
	maxTime := now.Add(tv.window)

	if txTime.Before(minTime) {
		return fmt.Errorf("transaction timestamp too old: %v (minimum: %v)", txTime, minTime)
	}
	if txTime.After(maxTime) {
		return fmt.Errorf("transaction timestamp too far in the future: %v (maximum: %v)", txTime, maxTime)
	}

	return nil
}

// ValidateBlockTime validates that a transaction is not too old based on block time
func (tv *TimestampValidator) ValidateBlockTime(blockTime time.Time, txTime time.Time) error {
	// Transaction must not be older than the window from block time
	minValidTime := blockTime.Add(-tv.window)

	if txTime.Before(minValidTime) {
		return fmt.Errorf("transaction too old: tx time %v, block time %v, window %v", txTime, blockTime, tv.window)
	}

	return nil
}

// ReplayGuard provides comprehensive replay attack protection
type ReplayGuard struct {
	nonceTracker      *NonceTracker
	sequenceValidator *SequenceValidator
	timestampValidator *TimestampValidator
	mu                sync.RWMutex
}

// NewReplayGuard creates a new replay guard
func NewReplayGuard(timestampWindow time.Duration) *ReplayGuard {
	return &ReplayGuard{
		nonceTracker:       NewNonceTracker(timestampWindow),
		sequenceValidator:  NewSequenceValidator(),
		timestampValidator: NewTimestampValidator(timestampWindow),
	}
}

// ValidateTransaction validates a transaction for replay attacks
// This should be called during transaction processing
func (rg *ReplayGuard) ValidateTransaction(ctx sdk.Context, address string, sequence uint64, nonce uint64, txTime time.Time) error {
	// Validate timestamp
	if err := rg.timestampValidator.ValidateBlockTime(ctx.BlockTime(), txTime); err != nil {
		return fmt.Errorf("timestamp validation failed: %w", err)
	}

	// Validate sequence (for account-based replay protection)
	if err := rg.sequenceValidator.ValidateAndUpdate(address, sequence); err != nil {
		return fmt.Errorf("sequence validation failed: %w", err)
	}

	// Check nonce (for additional replay protection)
	if !rg.nonceTracker.CheckAndRecord(address, nonce) {
		return fmt.Errorf("nonce already used: %d", nonce)
	}

	return nil
}

// ValidateBasicSequence validates that a sequence number is reasonable
// This is a lightweight check that can be used in ValidateBasic
func ValidateBasicSequence(sequence uint64) error {
	// Reject obviously invalid sequences (e.g., max uint64)
	if sequence == ^uint64(0) {
		return fmt.Errorf("invalid sequence number")
	}
	return nil
}

// ValidateTimestampRange validates that a timestamp is within a reasonable range
// This is a lightweight check that can be used in ValidateBasic
func ValidateTimestampRange(timestamp int64) error {
	if timestamp < 0 {
		return fmt.Errorf("timestamp cannot be negative")
	}

	// Convert to time
	txTime := time.Unix(timestamp, 0)
	now := time.Now()

	// Check if timestamp is too far in the future (1 hour)
	if txTime.After(now.Add(time.Hour)) {
		return fmt.Errorf("timestamp too far in the future")
	}

	// Check if timestamp is too old (24 hours ago)
	if txTime.Before(now.Add(-24 * time.Hour)) {
		return fmt.Errorf("timestamp too old")
	}

	return nil
}

// ReplayProtectedTx is an interface for transactions that need replay protection
type ReplayProtectedTx interface {
	GetSigners() []sdk.AccAddress
	GetSequence() uint64
}
