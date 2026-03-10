package types

import (
	"fmt"
)

// GenesisState defines the llmcustody module's genesis state.
type GenesisState struct {
	APIKeys []APIKey `json:"api_keys"`
	// EncryptionKey is the key encryption key (KEK) used to encrypt API keys
	// In production, this should be managed by a secure key management service
	EncryptionKey []byte `json:"encryption_key,omitempty"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		APIKeys:       []APIKey{},
		EncryptionKey: nil,
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	for _, key := range data.APIKeys {
		if err := key.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid API key %s: %w", key.ID, err)
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
