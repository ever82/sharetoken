package types

import (
	"fmt"
)

// GenesisState defines the dispute module's genesis state.
type GenesisState struct {
	Disputes []Dispute `json:"disputes"`
	JurorPool []string `json:"juror_pool"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Disputes:  []Dispute{},
		JurorPool: []string{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	seenIDs := make(map[string]bool)
	for _, dispute := range data.Disputes {
		if seenIDs[dispute.ID] {
			return fmt.Errorf("duplicate dispute ID: %s", dispute.ID)
		}
		seenIDs[dispute.ID] = true

		if dispute.ID == "" {
			return fmt.Errorf("dispute ID cannot be empty")
		}
		if dispute.EscrowID == "" {
			return fmt.Errorf("escrow ID cannot be empty")
		}
		if dispute.Requester == "" {
			return fmt.Errorf("requester cannot be empty")
		}
		if dispute.Provider == "" {
			return fmt.Errorf("provider cannot be empty")
		}
	}

	// Validate juror pool (no duplicates)
	seenJurors := make(map[string]bool)
	for _, juror := range data.JurorPool {
		if seenJurors[juror] {
			return fmt.Errorf("duplicate juror in pool: %s", juror)
		}
		seenJurors[juror] = true
	}

	return nil
}
