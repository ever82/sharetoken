package types

import (
	"fmt"
)

// GenesisState is the initial state for the agentgateway module
type GenesisState struct {
	// Sessions stores active sessions
	Sessions []Session `json:"sessions"`
}

// Session represents an agent gateway session
type Session struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	CreatedAt int64  `json:"created_at"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Sessions: []Session{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	for _, session := range data.Sessions {
		if session.ID == "" {
			return fmt.Errorf("session ID cannot be empty")
		}
		if session.Address == "" {
			return fmt.Errorf("session address cannot be empty")
		}
	}
	return nil
}
