package types

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for llmcustody module
const (
	TypeMsgRegisterAPIKey = "register_api_key"
	TypeMsgUpdateAPIKey   = "update_api_key"
	TypeMsgRevokeAPIKey   = "revoke_api_key"
	TypeMsgRecordUsage    = "record_usage"
)

// MsgServer is the message server interface
type MsgServer interface {
	RegisterAPIKey(ctx context.Context, msg *MsgRegisterAPIKey) (*MsgRegisterAPIKeyResponse, error)
	UpdateAPIKey(ctx context.Context, msg *MsgUpdateAPIKey) (*MsgUpdateAPIKeyResponse, error)
	RevokeAPIKey(ctx context.Context, msg *MsgRevokeAPIKey) (*MsgRevokeAPIKeyResponse, error)
	RecordUsage(ctx context.Context, msg *MsgRecordUsage) (*MsgRecordUsageResponse, error)
}

// Response types
type MsgRegisterAPIKeyResponse struct {
	APIKeyID string `json:"api_key_id"`
}

type MsgUpdateAPIKeyResponse struct{}
type MsgRevokeAPIKeyResponse struct{}
type MsgRecordUsageResponse struct{}

// GenerateAPIKeyID generates a unique API key ID
func GenerateAPIKeyID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// -----------------------------------------------------------------------------
// MsgRegisterAPIKey
// -----------------------------------------------------------------------------

// MsgRegisterAPIKey is the message for registering a new API key
type MsgRegisterAPIKey struct {
	Owner       string       `json:"owner"`
	Provider    string       `json:"provider"`
	EncryptedKey []byte      `json:"encrypted_key"`
	AccessRules []AccessRule `json:"access_rules"`
}

// NewMsgRegisterAPIKey creates a new MsgRegisterAPIKey
func NewMsgRegisterAPIKey(owner string, provider string, encryptedKey []byte, rules []AccessRule) *MsgRegisterAPIKey {
	return &MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     provider,
		EncryptedKey: encryptedKey,
		AccessRules:  rules,
	}
}

// Route returns the module route
func (msg MsgRegisterAPIKey) Route() string { return RouterKey }

// Type returns the message type
func (msg MsgRegisterAPIKey) Type() string { return TypeMsgRegisterAPIKey }

// GetSigners returns the signers
func (msg MsgRegisterAPIKey) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the bytes to sign
func (msg MsgRegisterAPIKey) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic performs basic validation
func (msg MsgRegisterAPIKey) ValidateBasic() error {
	if msg.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if !IsValidProvider(msg.Provider) {
		return fmt.Errorf("invalid provider: %s", msg.Provider)
	}
	if len(msg.EncryptedKey) == 0 {
		return fmt.Errorf("encrypted key cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgUpdateAPIKey
// -----------------------------------------------------------------------------

// MsgUpdateAPIKey is the message for updating an API key
type MsgUpdateAPIKey struct {
	Owner       string       `json:"owner"`
	APIKeyID    string       `json:"api_key_id"`
	AccessRules []AccessRule `json:"access_rules"`
	Active      bool         `json:"active"`
}

// NewMsgUpdateAPIKey creates a new MsgUpdateAPIKey
func NewMsgUpdateAPIKey(owner string, apiKeyID string, rules []AccessRule, active bool) *MsgUpdateAPIKey {
	return &MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    apiKeyID,
		AccessRules: rules,
		Active:      active,
	}
}

// Route returns the module route
func (msg MsgUpdateAPIKey) Route() string { return RouterKey }

// Type returns the message type
func (msg MsgUpdateAPIKey) Type() string { return TypeMsgUpdateAPIKey }

// GetSigners returns the signers
func (msg MsgUpdateAPIKey) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the bytes to sign
func (msg MsgUpdateAPIKey) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic performs basic validation
func (msg MsgUpdateAPIKey) ValidateBasic() error {
	if msg.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if msg.APIKeyID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgRevokeAPIKey
// -----------------------------------------------------------------------------

// MsgRevokeAPIKey is the message for revoking an API key
type MsgRevokeAPIKey struct {
	Owner    string `json:"owner"`
	APIKeyID string `json:"api_key_id"`
}

// NewMsgRevokeAPIKey creates a new MsgRevokeAPIKey
func NewMsgRevokeAPIKey(owner string, apiKeyID string) *MsgRevokeAPIKey {
	return &MsgRevokeAPIKey{
		Owner:    owner,
		APIKeyID: apiKeyID,
	}
}

// Route returns the module route
func (msg MsgRevokeAPIKey) Route() string { return RouterKey }

// Type returns the message type
func (msg MsgRevokeAPIKey) Type() string { return TypeMsgRevokeAPIKey }

// GetSigners returns the signers
func (msg MsgRevokeAPIKey) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the bytes to sign
func (msg MsgRevokeAPIKey) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic performs basic validation
func (msg MsgRevokeAPIKey) ValidateBasic() error {
	if msg.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if msg.APIKeyID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgRecordUsage
// -----------------------------------------------------------------------------

// MsgRecordUsage is the message for recording API usage
type MsgRecordUsage struct {
	APIKeyID  string `json:"api_key_id"`
	ServiceID string `json:"service_id"`
	RequestCount int64 `json:"request_count"`
	TokenCount   int64 `json:"token_count"`
	Cost         int64 `json:"cost"` // in ustt
}

// NewMsgRecordUsage creates a new MsgRecordUsage
func NewMsgRecordUsage(apiKeyID string, serviceID string, requests int64, tokens int64, cost int64) *MsgRecordUsage {
	return &MsgRecordUsage{
		APIKeyID:     apiKeyID,
		ServiceID:    serviceID,
		RequestCount: requests,
		TokenCount:   tokens,
		Cost:         cost,
	}
}

// Route returns the module route
func (msg MsgRecordUsage) Route() string { return RouterKey }

// Type returns the message type
func (msg MsgRecordUsage) Type() string { return TypeMsgRecordUsage }

// GetSigners returns the signers
func (msg MsgRecordUsage) GetSigners() []sdk.AccAddress {
	// This message is typically signed by the service provider or oracle
	return []sdk.AccAddress{}
}

// GetSignBytes returns the bytes to sign
func (msg MsgRecordUsage) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic performs basic validation
func (msg MsgRecordUsage) ValidateBasic() error {
	if msg.APIKeyID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if msg.ServiceID == "" {
		return fmt.Errorf("service ID cannot be empty")
	}
	if msg.RequestCount < 0 {
		return fmt.Errorf("request count cannot be negative")
	}
	if msg.TokenCount < 0 {
		return fmt.Errorf("token count cannot be negative")
	}
	return nil
}
