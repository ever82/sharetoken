package types

import (
	"fmt"
)

// GenesisState is the initial state for the node module
type GenesisState struct {
	// NodeConfigs stores node configurations
	NodeConfigs []NodeConfigEntry `json:"node_configs"`
}

// NodeConfigEntry represents a node configuration entry
type NodeConfigEntry struct {
	Address string     `json:"address"`
	Config  RoleConfig `json:"config"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		NodeConfigs: []NodeConfigEntry{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	for _, config := range data.NodeConfigs {
		if config.Address == "" {
			return fmt.Errorf("node address cannot be empty")
		}
		if err := config.Config.Validate(); err != nil {
			return fmt.Errorf("invalid config for address %s: %w", config.Address, err)
		}
	}
	return nil
}
