package types

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for llmcustody module
const (
	TypeMsgRegisterAPIKey = "register_api_key"
	TypeMsgUpdateAPIKey   = "update_api_key"
	TypeMsgRevokeAPIKey   = "revoke_api_key"
	TypeMsgRecordUsage    = "record_usage"
	TypeMsgRotateAPIKey   = "rotate_api_key"
)

// MsgServer is the message server interface
// Note: Using protobuf-generated interface from tx.pb.go

// GenerateAPIKeyID generates a unique API key ID
func GenerateAPIKeyID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
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

// ValidateBasic performs basic validation
func (msg MsgUpdateAPIKey) ValidateBasic() error {
	if msg.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if msg.ApiKeyId == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
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

// ValidateBasic performs basic validation
func (msg MsgRevokeAPIKey) ValidateBasic() error {
	if msg.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if msg.ApiKeyId == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
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

// ValidateBasic performs basic validation
func (msg MsgRecordUsage) ValidateBasic() error {
	if msg.ApiKeyId == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if msg.ServiceId == "" {
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

// Route returns the module route
func (msg MsgRotateAPIKey) Route() string { return RouterKey }

// Type returns the message type
func (msg MsgRotateAPIKey) Type() string { return TypeMsgRotateAPIKey }

// GetSigners returns the signers
func (msg MsgRotateAPIKey) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic performs basic validation
func (msg MsgRotateAPIKey) ValidateBasic() error {
	if msg.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if msg.ApiKeyId == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if len(msg.NewEncryptedKey) == 0 {
		return fmt.Errorf("new encrypted key cannot be empty")
	}
	return nil
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

// NewMsgUpdateAPIKey creates a new MsgUpdateAPIKey
func NewMsgUpdateAPIKey(owner string, apiKeyID string, rules []AccessRule, active bool) *MsgUpdateAPIKey {
	return &MsgUpdateAPIKey{
		Owner:       owner,
		ApiKeyId:    apiKeyID,
		AccessRules: rules,
		Active:      active,
	}
}

// NewMsgRevokeAPIKey creates a new MsgRevokeAPIKey
func NewMsgRevokeAPIKey(owner string, apiKeyID string) *MsgRevokeAPIKey {
	return &MsgRevokeAPIKey{
		Owner:    owner,
		ApiKeyId: apiKeyID,
	}
}

// NewMsgRecordUsage creates a new MsgRecordUsage
func NewMsgRecordUsage(apiKeyID string, serviceID string, requests int64, tokens int64, cost int64) *MsgRecordUsage {
	return &MsgRecordUsage{
		ApiKeyId:     apiKeyID,
		ServiceId:    serviceID,
		RequestCount: requests,
		TokenCount:   tokens,
		Cost:         cost,
	}
}

// NewMsgRotateAPIKey creates a new MsgRotateAPIKey
func NewMsgRotateAPIKey(owner, apiKeyID string, newEncryptedKey []byte, reason string) *MsgRotateAPIKey {
	return &MsgRotateAPIKey{
		Owner:           owner,
		ApiKeyId:        apiKeyID,
		NewEncryptedKey: newEncryptedKey,
		Reason:          reason,
	}
}
