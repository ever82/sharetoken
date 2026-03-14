package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		MqScores: []MQScore{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(gs GenesisState) error {
	seenAddresses := make(map[string]bool)
	for _, score := range gs.MqScores {
		if score.Address == "" {
			return fmt.Errorf("MQ score address cannot be empty")
		}
		if seenAddresses[score.Address] {
			return fmt.Errorf("duplicate MQ score for address: %s", score.Address)
		}
		seenAddresses[score.Address] = true
		if score.Score < 0 || score.Score > 1000 {
			return fmt.Errorf("MQ score out of range for address %s: %d", score.Address, score.Score)
		}
	}
	return nil
}
