package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewEscrow creates a new escrow with the given parameters
func NewEscrow(id, requester, provider string, amount sdk.Coins, duration time.Duration) *Escrow {
	now := time.Now().Unix()
	return &Escrow{
		Id:              id,
		Requester:       requester,
		Provider:        provider,
		Amount:          amount,
		Status:          EscrowStatus_ESCROW_STATUS_PENDING,
		CreatedAt:       now,
		ExpiresAt:       now + int64(duration.Seconds()),
		CompletedAt:     0,
		CompletionProof: "",
		DisputeId:       "",
		RefundAddress:   requester,
	}
}

// CanComplete returns true if the escrow can be completed
func (e Escrow) CanComplete() bool {
	return e.Status == EscrowStatus_ESCROW_STATUS_PENDING && !e.IsExpired()
}

// CanRefund returns true if the escrow can be refunded
func (e Escrow) CanRefund() bool {
	return e.Status == EscrowStatus_ESCROW_STATUS_PENDING && e.IsExpired()
}

// CanDispute returns true if the escrow can be disputed
func (e Escrow) CanDispute() bool {
	return e.Status == EscrowStatus_ESCROW_STATUS_PENDING || e.Status == EscrowStatus_ESCROW_STATUS_COMPLETED
}

// IsExpired returns true if the escrow has expired
func (e Escrow) IsExpired() bool {
	return time.Now().Unix() > e.ExpiresAt
}

// ValidateBasic performs basic validation of the escrow
func (e Escrow) ValidateBasic() error {
	if e.Id == "" {
		return ErrInvalidEscrowID
	}
	if e.Requester == "" {
		return ErrInvalidRequester
	}
	if e.Provider == "" {
		return ErrInvalidProvider
	}
	if !sdk.Coins(e.Amount).IsZero() {
		if err := sdk.Coins(e.Amount).Validate(); err != nil {
			return ErrInvalidAmount.Wrap(err.Error())
		}
	} else {
		return ErrInvalidAmount.Wrap("amount cannot be zero")
	}
	return nil
}

// Validate validates that the allocation sums up to the expected total
func (fa FundAllocation) Validate(total sdk.Coins) error {
	requesterCoins := sdk.Coins(fa.RequesterAmount)
	providerCoins := sdk.Coins(fa.ProviderAmount)
	sum := requesterCoins.Add(providerCoins...)
	if !sum.IsEqual(total) {
		return ErrInvalidAmount.Wrap("allocation does not sum to total")
	}
	return nil
}
