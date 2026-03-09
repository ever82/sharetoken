package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for identity module
const (
	TypeMsgRegisterIdentity  = "register_identity"
	TypeMsgVerifyIdentity    = "verify_identity"
	TypeMsgUpdateLimitConfig = "update_limit_config"
	TypeMsgResetDailyLimits  = "reset_daily_limits"
)

// MsgRegisterIdentity defines the message for registering a new identity
type MsgRegisterIdentity struct {
	Address      string `json:"address"`
	DID          string `json:"did"`
	MetadataHash string `json:"metadata_hash"`
}

// NewMsgRegisterIdentity creates a new MsgRegisterIdentity
func NewMsgRegisterIdentity(address, did, metadataHash string) *MsgRegisterIdentity {
	return &MsgRegisterIdentity{
		Address:      address,
		DID:          did,
		MetadataHash: metadataHash,
	}
}

// Route implements sdk.Msg
func (msg MsgRegisterIdentity) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgRegisterIdentity) Type() string {
	return TypeMsgRegisterIdentity
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterIdentity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes implements sdk.Msg
func (msg MsgRegisterIdentity) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterIdentity) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAddress
	}
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	if msg.DID != "" && !isValidDID(msg.DID) {
		return ErrInvalidDID.Wrap(msg.DID)
	}
	return nil
}

// MsgVerifyIdentity defines the message for verifying an identity
type MsgVerifyIdentity struct {
	Address          string `json:"address"`
	Provider         string `json:"provider"`
	VerificationHash string `json:"verification_hash"`
	Proof            string `json:"proof"`
}

// NewMsgVerifyIdentity creates a new MsgVerifyIdentity
func NewMsgVerifyIdentity(address, provider, verificationHash, proof string) *MsgVerifyIdentity {
	return &MsgVerifyIdentity{
		Address:          address,
		Provider:         provider,
		VerificationHash: verificationHash,
		Proof:            proof,
	}
}

// Route implements sdk.Msg
func (msg MsgVerifyIdentity) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgVerifyIdentity) Type() string {
	return TypeMsgVerifyIdentity
}

// GetSigners implements sdk.Msg
func (msg MsgVerifyIdentity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes implements sdk.Msg
func (msg MsgVerifyIdentity) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements sdk.Msg
func (msg MsgVerifyIdentity) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAddress
	}
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	if !IsValidProvider(msg.Provider) {
		return ErrInvalidProvider.Wrap(msg.Provider)
	}
	return nil
}

// MsgUpdateLimitConfig defines the message for updating limit configuration
type MsgUpdateLimitConfig struct {
	Authority     string      `json:"authority"`
	TargetAddress string      `json:"target_address"`
	NewConfig     LimitConfig `json:"new_config"`
}

// NewMsgUpdateLimitConfig creates a new MsgUpdateLimitConfig
func NewMsgUpdateLimitConfig(authority, targetAddress string, newConfig LimitConfig) *MsgUpdateLimitConfig {
	return &MsgUpdateLimitConfig{
		Authority:     authority,
		TargetAddress: targetAddress,
		NewConfig:     newConfig,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateLimitConfig) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateLimitConfig) Type() string {
	return TypeMsgUpdateLimitConfig
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateLimitConfig) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateLimitConfig) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateLimitConfig) ValidateBasic() error {
	if msg.Authority == "" {
		return ErrUnauthorized.Wrap("authority address required")
	}
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	if msg.TargetAddress == "" {
		return ErrInvalidAddress
	}
	_, err = sdk.AccAddressFromBech32(msg.TargetAddress)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	return nil
}

// MsgResetDailyLimits defines the message for resetting daily limits
type MsgResetDailyLimits struct {
	Authority string `json:"authority"`
}

// NewMsgResetDailyLimits creates a new MsgResetDailyLimits
func NewMsgResetDailyLimits(authority string) *MsgResetDailyLimits {
	return &MsgResetDailyLimits{
		Authority: authority,
	}
}

// Route implements sdk.Msg
func (msg MsgResetDailyLimits) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgResetDailyLimits) Type() string {
	return TypeMsgResetDailyLimits
}

// GetSigners implements sdk.Msg
func (msg MsgResetDailyLimits) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes implements sdk.Msg
func (msg MsgResetDailyLimits) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements sdk.Msg
func (msg MsgResetDailyLimits) ValidateBasic() error {
	if msg.Authority == "" {
		return ErrUnauthorized.Wrap("authority address required")
	}
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	return nil
}
