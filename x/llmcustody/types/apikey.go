package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Provider constants for backward compatibility
const (
	ProviderOpenAI    = Provider_PROVIDER_OPENAI
	ProviderAnthropic = Provider_PROVIDER_ANTHROPIC
)

// IsValidProvider checks if provider is valid
func IsValidProvider(p string) bool {
	return p == "openai" || p == "anthropic"
}

// ProviderFromString converts string to Provider
func ProviderFromString(s string) Provider {
	if s == "anthropic" || s == "PROVIDER_ANTHROPIC" {
		return Provider_PROVIDER_ANTHROPIC
	}
	if s == "openai" || s == "PROVIDER_OPENAI" {
		return Provider_PROVIDER_OPENAI
	}
	return Provider_PROVIDER_UNSPECIFIED
}

// NewAPIKey creates a new API key record
func NewAPIKey(id string, provider Provider, encryptedKey []byte, owner string) *APIKey {
	// Calculate hash of the encrypted key for integrity verification
	hash := sha256.Sum256(encryptedKey)

	return &APIKey{
		Id:           id,
		Provider:     provider,
		EncryptedKey: encryptedKey,
		Hash:         hex.EncodeToString(hash[:]),
		Owner:        owner,
		AccessRules:  []AccessRule{},
		Active:       true,
	}
}

// ValidateBasic performs basic validation
func (k APIKey) ValidateBasic() error {
	if k.Id == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if k.Provider != Provider_PROVIDER_OPENAI && k.Provider != Provider_PROVIDER_ANTHROPIC {
		return fmt.Errorf("invalid provider: %s", k.Provider)
	}
	if len(k.EncryptedKey) == 0 {
		return fmt.Errorf("encrypted key cannot be empty")
	}
	if k.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	return nil
}

// VerifyHash verifies the encrypted key hash matches
func (k APIKey) VerifyHash(encryptedKey []byte) bool {
	hash := sha256.Sum256(encryptedKey)
	return hex.EncodeToString(hash[:]) == k.Hash
}

// CanAccess checks if the API key can be used for a service
func (k APIKey) CanAccess(serviceID string) bool {
	if !k.Active {
		return false
	}

	// Check if there's a rule allowing access to this service
	for _, rule := range k.AccessRules {
		if rule.ServiceId != serviceID {
			continue
		}
		if rule.MaxRequests == 0 {
			return true
		}
		return k.UsageCount < rule.MaxRequests
	}
	return false
}

// RecordUsage records a usage of the API key
func (k *APIKey) RecordUsage() {
	k.UsageCount++
	k.LastUsedAt = 0 // Will be set by keeper with block time
}

// SecureWipe securely wipes sensitive data from memory
// Note: This is a best-effort approach in Go due to garbage collection
func (k *APIKey) SecureWipe() {
	// Overwrite encrypted key with zeros
	for i := range k.EncryptedKey {
		k.EncryptedKey[i] = 0
	}
	k.EncryptedKey = nil
}

// SetActive sets the active status of the API key
func (k *APIKey) SetActive(active bool) {
	k.Active = active
}
