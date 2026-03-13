package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
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

// ProviderFromString converts string to Provider
func ProviderFromString(s string) Provider {
	if s == string(ProviderAnthropic) {
		return ProviderAnthropic
	}
	return ProviderOpenAI
}

// AccessRule represents an access control rule
type AccessRule struct {
	ServiceID   string `json:"service_id"`    // 允许访问的服务ID
	Allowed     bool   `json:"allowed"`       // 是否允许访问
	RateLimit   int64  `json:"rate_limit"`    // 每分钟最大请求数
	MaxRequests int64  `json:"max_requests"`  // 总请求上限
	PricePerReq int64  `json:"price_per_req"` // 每次请求价格 (ustt)
}

// Reset implements proto.Message
func (m *AccessRule) Reset() { *m = AccessRule{} }

// String implements proto.Message
func (m AccessRule) String() string {
	return fmt.Sprintf("AccessRule{service_id: %s, allowed: %v}", m.ServiceID, m.Allowed)
}

// ProtoMessage implements proto.Message
func (*AccessRule) ProtoMessage() {}

// APIKey represents an encrypted API key stored on chain
type APIKey struct {
	ID           string       `json:"id"`
	Provider     Provider     `json:"provider"`
	EncryptedKey []byte       `json:"encrypted_key"` // AES-256-GCM encrypted
	Hash         string       `json:"hash"`          // SHA-256 hash for verification
	Owner        string       `json:"owner"`
	AccessRules  []AccessRule `json:"access_rules"`
	CreatedAt    int64        `json:"created_at"`
	LastUsedAt   int64        `json:"last_used_at"`
	UsageCount   int64        `json:"usage_count"`
	Active       bool         `json:"active"`
	Version      int          `json:"version"` // For key rotation tracking
}

// Reset implements proto.Message
func (m *APIKey) Reset() { *m = APIKey{} }

// String implements proto.Message
func (m APIKey) String() string {
	return fmt.Sprintf("APIKey{%s: %s, owner: %s, active: %v, usage: %d}",
		m.ID, m.Provider, m.Owner, m.Active, m.UsageCount)
}

// ProtoMessage implements proto.Message
func (*APIKey) ProtoMessage() {}

// Marshal implements codec.ProtoMarshaler
func (m APIKey) Marshal() ([]byte, error) {
	return proto.Marshal(m.toProto())
}

// MarshalTo implements codec.ProtoMarshaler
func (m APIKey) MarshalTo(data []byte) (n int, err error) {
	return m.MarshalToSizedBuffer(data)
}

// MarshalToSizedBuffer implements codec.ProtoMarshaler
func (m APIKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	encoded, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(encoded)
	if len(dAtA) < n {
		return 0, fmt.Errorf("buffer too small")
	}
	copy(dAtA[:n], encoded)
	return n, nil
}

// Size implements codec.ProtoMarshaler
func (m APIKey) Size() int {
	data, _ := m.Marshal()
	return len(data)
}

// Unmarshal implements codec.ProtoMarshaler
func (m *APIKey) Unmarshal(data []byte) error {
	pm := &APIKeyProto{}
	if err := proto.Unmarshal(data, pm); err != nil {
		return err
	}
	m.fromProto(pm)
	return nil
}

// toProto converts APIKey to proto message
func (m APIKey) toProto() *APIKeyProto {
	rules := make([]*AccessRuleProto, len(m.AccessRules))
	for i, r := range m.AccessRules {
		rules[i] = &AccessRuleProto{
			ServiceId:   r.ServiceID,
			Allowed:     r.Allowed,
			RateLimit:   r.RateLimit,
			MaxRequests: r.MaxRequests,
			PricePerReq: r.PricePerReq,
		}
	}
	return &APIKeyProto{
		Id:           m.ID,
		Provider:     string(m.Provider),
		EncryptedKey: m.EncryptedKey,
		Hash:         m.Hash,
		Owner:        m.Owner,
		AccessRules:  rules,
		CreatedAt:    m.CreatedAt,
		LastUsedAt:   m.LastUsedAt,
		UsageCount:   m.UsageCount,
		Active:       m.Active,
	}
}

// fromProto converts proto message to APIKey
func (m *APIKey) fromProto(pm *APIKeyProto) {
	m.ID = pm.Id
	m.Provider = ProviderFromString(pm.Provider)
	m.EncryptedKey = pm.EncryptedKey
	m.Hash = pm.Hash
	m.Owner = pm.Owner
	m.AccessRules = make([]AccessRule, len(pm.AccessRules))
	for i, r := range pm.AccessRules {
		m.AccessRules[i] = AccessRule{
			ServiceID:   r.ServiceId,
			Allowed:     r.Allowed,
			RateLimit:   r.RateLimit,
			MaxRequests: r.MaxRequests,
			PricePerReq: r.PricePerReq,
		}
	}
	m.CreatedAt = pm.CreatedAt
	m.LastUsedAt = pm.LastUsedAt
	m.UsageCount = pm.UsageCount
	m.Active = pm.Active
	m.Version = 1
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
		Version:      1,
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
		if rule.ServiceID != serviceID {
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

// GetVersion returns the version of the API key
func (k APIKey) GetVersion() int {
	return k.Version
}

// IncrementVersion increments the version of the API key
func (k *APIKey) IncrementVersion() {
	k.Version++
}
