package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Disputes:  []Dispute{},
		JurorPool: []string{},
	}
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		VotingPeriod: 86400, // 1 day in seconds
		MinJurors:    3,
		MinMqScore:   100,
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(state GenesisState) error {
	// Validate disputes
	seenIDs := make(map[string]bool)
	for _, dispute := range state.Disputes {
		if dispute.Id == "" {
			return fmt.Errorf("dispute ID cannot be empty")
		}
		if seenIDs[dispute.Id] {
			return fmt.Errorf("duplicate dispute ID: %s", dispute.Id)
		}
		seenIDs[dispute.Id] = true

		if dispute.EscrowId == "" {
			return fmt.Errorf("escrow ID cannot be empty for dispute %s", dispute.Id)
		}
		if dispute.Requester == "" {
			return fmt.Errorf("requester cannot be empty for dispute %s", dispute.Id)
		}
		if dispute.Provider == "" {
			return fmt.Errorf("provider cannot be empty for dispute %s", dispute.Id)
		}
	}

	// Validate juror pool
	seenJurors := make(map[string]bool)
	for _, juror := range state.JurorPool {
		if seenJurors[juror] {
			return fmt.Errorf("duplicate juror: %s", juror)
		}
		seenJurors[juror] = true
	}

	return nil
}

// Validate validates the params
func (p Params) Validate() error {
	if p.VotingPeriod == 0 {
		return fmt.Errorf("voting period cannot be zero")
	}
	if p.MinJurors == 0 {
		return fmt.Errorf("min jurors cannot be zero")
	}
	return nil
}
