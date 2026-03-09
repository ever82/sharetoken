package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EscrowStatus represents the status of an escrow
type EscrowStatus string

const (
	// EscrowStatusPending means the escrow is active and waiting for completion
	EscrowStatusPending EscrowStatus = "pending"
	// EscrowStatusCompleted means the escrow has been completed and funds released
	EscrowStatusCompleted EscrowStatus = "completed"
	// EscrowStatusDisputed means the escrow is under dispute
	EscrowStatusDisputed EscrowStatus = "disputed"
	// EscrowStatusRefunded means the escrow has been refunded to the requester
	EscrowStatusRefunded EscrowStatus = "refunded"
)

// Escrow represents an escrow agreement
type Escrow struct {
	ID              string       `json:"id"`
	Requester       string       `json:"requester"`
	Provider        string       `json:"provider"`
	Amount          sdk.Coins    `json:"amount"`
	Status          EscrowStatus `json:"status"`
	CreatedAt       int64        `json:"created_at"`
	ExpiresAt       int64        `json:"expires_at"`
	CompletedAt     int64        `json:"completed_at"`
	CompletionProof string       `json:"completion_proof"`
	DisputeID       string       `json:"dispute_id"`
	RefundAddress   string       `json:"refund_address"`
}

// NewEscrow creates a new escrow
func NewEscrow(id, requester, provider string, amount sdk.Coins, duration time.Duration) *Escrow {
	now := time.Now().Unix()
	return &Escrow{
		ID:            id,
		Requester:     requester,
		Provider:      provider,
		Amount:        amount,
		Status:        EscrowStatusPending,
		CreatedAt:     now,
		ExpiresAt:     now + int64(duration.Seconds()),
		RefundAddress: requester,
	}
}

// ValidateBasic performs basic validation of escrow fields
func (e Escrow) ValidateBasic() error {
	if e.ID == "" {
		return ErrInvalidEscrowID
	}
	if e.Requester == "" {
		return ErrUnauthorized.Wrap("requester address required")
	}
	if _, err := sdk.AccAddressFromBech32(e.Requester); err != nil {
		return ErrUnauthorized.Wrap("invalid requester address")
	}
	if e.Provider == "" {
		return ErrUnauthorized.Wrap("provider address required")
	}
	if _, err := sdk.AccAddressFromBech32(e.Provider); err != nil {
		return ErrUnauthorized.Wrap("invalid provider address")
	}
	if !e.Amount.IsValid() || e.Amount.IsZero() {
		return ErrInvalidAmount
	}
	if e.Status == "" {
		return ErrInvalidStatus
	}
	return nil
}

// IsExpired checks if the escrow has expired
func (e Escrow) IsExpired() bool {
	return time.Now().Unix() > e.ExpiresAt
}

// CanComplete checks if the escrow can be completed
func (e Escrow) CanComplete() bool {
	return e.Status == EscrowStatusPending && !e.IsExpired()
}

// CanDispute checks if the escrow can be disputed
func (e Escrow) CanDispute() bool {
	return e.Status == EscrowStatusPending || e.Status == EscrowStatusCompleted
}

// CanRefund checks if the escrow can be refunded
func (e Escrow) CanRefund() bool {
	return e.Status == EscrowStatusPending && e.IsExpired()
}

// String implements stringer interface
func (e Escrow) String() string {
	return fmt.Sprintf(`
Escrow %s:
  Requester: %s
  Provider: %s
  Amount: %s
  Status: %s
  Created At: %d
  Expires At: %d
  Dispute ID: %s
`, e.ID, e.Requester, e.Provider, e.Amount.String(), e.Status, e.CreatedAt, e.ExpiresAt, e.DisputeID)
}

// FundAllocation represents how funds should be allocated after a dispute
type FundAllocation struct {
	RequesterAmount sdk.Coins `json:"requester_amount"`
	ProviderAmount  sdk.Coins `json:"provider_amount"`
}

// Validate validates the fund allocation
func (fa FundAllocation) Validate(total sdk.Coins) error {
	totalAllocated := fa.RequesterAmount.Add(fa.ProviderAmount...)
	if !totalAllocated.IsEqual(total) {
		return ErrInvalidAllocation.Wrapf("total allocated %s does not match escrow amount %s", totalAllocated.String(), total.String())
	}
	return nil
}

// GenesisState defines the escrow module's genesis state
type GenesisState struct {
	Escrows []Escrow `json:"escrows"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Escrows: []Escrow{},
	}
}

// ValidateGenesis validates genesis state
func ValidateGenesis(data GenesisState) error {
	seenIDs := make(map[string]bool)
	for _, escrow := range data.Escrows {
		if seenIDs[escrow.ID] {
			return fmt.Errorf("duplicate escrow ID: %s", escrow.ID)
		}
		seenIDs[escrow.ID] = true

		if err := escrow.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
