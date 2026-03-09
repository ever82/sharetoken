package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Provider represents LLM provider type
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
)

// IsValidProvider checks if provider is valid
func IsValidProvider(p string) bool {
	return p == string(ProviderOpenAI) || p == string(ProviderAnthropic)
}

// AccessRule represents an access control rule
type AccessRule struct {
	ServiceID   string `json:"service_id"`   // 允许访问的服务ID
	RateLimit   int64  `json:"rate_limit"`   // 每分钟最大请求数
	MaxRequests int64  `json:"max_requests"` // 总请求上限
	PricePerReq int64  `json:"price_per_req"`// 每次请求价格 (ustt)
}

// APIKey represents an encrypted API key stored on chain
type APIKey struct {
	ID            string       `json:"id"`
	Provider      Provider     `json:"provider"`
	EncryptedKey  []byte       `json:"encrypted_key"` // AES-256-GCM encrypted
	Hash          string       `json:"hash"`          // SHA-256 hash for verification
	Owner         string       `json:"owner"`
	AccessRules   []AccessRule `json:"access_rules"`
	CreatedAt     int64        `json:"created_at"`
	LastUsedAt    int64        `json:"last_used_at"`
	UsageCount    int64        `json:"usage_count"`
	Active        bool         `json:"active"`
}

// NewAPIKey creates a new API key record
func NewAPIKey(id string, provider Provider, encryptedKey []byte, owner string) *APIKey {
	// Calculate hash of the encrypted key for integrity verification
	hash := sha256.Sum256(encryptedKey)

	return &APIKey{
		ID:           id,
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
	if k.ID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if k.Provider != ProviderOpenAI && k.Provider != ProviderAnthropic {
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
		if rule.ServiceID == serviceID {
			return k.UsageCount < rule.MaxRequests || rule.MaxRequests == 0
		}
	}
	return false
}

// RecordUsage records a usage of the API key
func (k *APIKey) RecordUsage() {
	k.UsageCount++
	k.LastUsedAt = 0 // Will be set by keeper with block time
}

// String implements stringer
func (k APIKey) String() string {
	return fmt.Sprintf("APIKey{%s: %s, owner: %s, active: %v, usage: %d}",
		k.ID, k.Provider, k.Owner, k.Active, k.UsageCount)
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
