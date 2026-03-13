package types

import (
	"fmt"
)

// GenesisState defines the marketplace module's genesis state.
type GenesisState struct {
	Services []Service `json:"services"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Services: []Service{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	seenIDs := make(map[string]bool)
	for _, service := range data.Services {
		if seenIDs[service.ID] {
			return fmt.Errorf("duplicate service ID: %s", service.ID)
		}
		seenIDs[service.ID] = true

		if service.ID == "" {
			return fmt.Errorf("service ID cannot be empty")
		}
		if service.Provider == "" {
			return fmt.Errorf("service provider cannot be empty")
		}
		if service.Name == "" {
			return fmt.Errorf("service name cannot be empty")
		}
	}
	return nil
}
