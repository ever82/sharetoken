package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Escrows: []Escrow{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(gs GenesisState) error {
	seenIDs := make(map[string]bool)
	for _, escrow := range gs.Escrows {
		if escrow.Id == "" {
			return fmt.Errorf("escrow ID cannot be empty")
		}
		if seenIDs[escrow.Id] {
			return fmt.Errorf("duplicate escrow ID: %s", escrow.Id)
		}
		seenIDs[escrow.Id] = true

		if escrow.Requester == "" {
			return fmt.Errorf("requester cannot be empty for escrow %s", escrow.Id)
		}
		if escrow.Provider == "" {
			return fmt.Errorf("provider cannot be empty for escrow %s", escrow.Id)
		}
	}
	return nil
}
