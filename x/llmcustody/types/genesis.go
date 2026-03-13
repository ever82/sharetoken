package types

import (
	"fmt"
)

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ApiKeys:       []APIKey{},
		EncryptionKey: nil,
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	for _, key := range data.ApiKeys {
		if err := key.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid API key %s: %w", key.Id, err)
		}
	}
	return nil
}

// ExportAPIKeyForGenesis exports an API key for genesis (with sensitive data cleared)
func ExportAPIKeyForGenesis(key APIKey) APIKey {
	// Clear sensitive data
	key.EncryptedKey = nil
	key.Hash = ""
	return key
}
