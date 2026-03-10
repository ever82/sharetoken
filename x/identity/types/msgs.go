package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSigners implements sdk.Msg
func (msg *MsgRegisterIdentity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterIdentity) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAddress.Wrap("address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	return nil
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
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgVerifyIdentity) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAddress.Wrap("address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	if msg.Provider == "" {
		return ErrInvalidProvider.Wrap("provider cannot be empty")
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
	addr, err := sdk.AccAddressFromBech32(msg.TargetAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateLimitConfig) ValidateBasic() error {
	if msg.TargetAddress == "" {
		return ErrInvalidAddress.Wrap("target address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.TargetAddress)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	return nil
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
	addr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgResetDailyLimits) ValidateBasic() error {
	if msg.Authority == "" {
		return ErrInvalidAddress.Wrap("authority cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	return nil
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
	addr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	if msg.Authority == "" {
		return ErrInvalidAddress.Wrap("authority cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}
	return nil
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
