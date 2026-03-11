package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/pkg/validation"
)

// GetSigners implements sdk.Msg
func (msg *MsgRegisterIdentity) GetSigners() []sdk.AccAddress {
	return validation.MustGetSigners(msg.Address)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterIdentity) ValidateBasic() error {
	return validation.ValidateAddress(msg.Address, "address")
}

// Type implements sdk.Msg
func (msg *MsgRegisterIdentity) Type() string {
	return "RegisterIdentity"
}

// GetSignBytes implements sdk.Msg
func (msg *MsgRegisterIdentity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements sdk.Msg
func (msg *MsgRegisterIdentity) Route() string {
	return RouterKey
}

// GetSigners implements sdk.Msg
func (msg *MsgVerifyIdentity) GetSigners() []sdk.AccAddress {
	return validation.MustGetSigners(msg.Address)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgVerifyIdentity) ValidateBasic() error {
	if err := validation.ValidateAddress(msg.Address, "address"); err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	if err := validation.ValidateNonEmpty(msg.Provider, "provider"); err != nil {
		return ErrInvalidProvider.Wrap(err.Error())
	}
	return nil
}

// Type implements sdk.Msg
func (msg *MsgVerifyIdentity) Type() string {
	return "VerifyIdentity"
}

// GetSignBytes implements sdk.Msg
func (msg *MsgVerifyIdentity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements sdk.Msg
func (msg *MsgVerifyIdentity) Route() string {
	return RouterKey
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateLimitConfig) GetSigners() []sdk.AccAddress {
	return validation.MustGetSigners(msg.TargetAddress)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateLimitConfig) ValidateBasic() error {
	return validation.ValidateAddress(msg.TargetAddress, "target address")
}

// Type implements sdk.Msg
func (msg *MsgUpdateLimitConfig) Type() string {
	return "UpdateLimitConfig"
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateLimitConfig) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements sdk.Msg
func (msg *MsgUpdateLimitConfig) Route() string {
	return RouterKey
}

// GetSigners implements sdk.Msg
func (msg *MsgResetDailyLimits) GetSigners() []sdk.AccAddress {
	return validation.MustGetSigners(msg.Authority)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgResetDailyLimits) ValidateBasic() error {
	return validation.ValidateAddress(msg.Authority, "authority")
}

// Type implements sdk.Msg
func (msg *MsgResetDailyLimits) Type() string {
	return "ResetDailyLimits"
}

// GetSignBytes implements sdk.Msg
func (msg *MsgResetDailyLimits) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements sdk.Msg
func (msg *MsgResetDailyLimits) Route() string {
	return RouterKey
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return validation.MustGetSigners(msg.Authority)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	return validation.ValidateAddress(msg.Authority, "authority")
}

// Type implements sdk.Msg
func (msg *MsgUpdateParams) Type() string {
	return "UpdateParams"
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements sdk.Msg
func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}
