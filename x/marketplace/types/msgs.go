package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// marketplace message types
const (
	TypeMsgRegisterService   = "register_service"
	TypeMsgUpdateService     = "update_service"
	TypeMsgActivateService   = "activate_service"
	TypeMsgDeactivateService = "deactivate_service"
	TypeMsgPurchaseService   = "purchase_service"
)

var (
	_ sdk.Msg = &MsgRegisterService{}
	_ sdk.Msg = &MsgUpdateService{}
	_ sdk.Msg = &MsgActivateService{}
	_ sdk.Msg = &MsgDeactivateService{}
	_ sdk.Msg = &MsgPurchaseService{}
)

// -- MsgRegisterService --

// Route Implements Msg
func (msg MsgRegisterService) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgRegisterService) Type() string { return TypeMsgRegisterService }

// GetSignBytes Implements Msg
func (msg MsgRegisterService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic Implements Msg
func (msg MsgRegisterService) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid provider address")
	}
	if msg.Name == "" {
		return sdkerrors.Wrap(ErrInvalidService, "name cannot be empty")
	}
	return nil
}

// GetSigners Implements Msg
func (msg MsgRegisterService) GetSigners() []sdk.AccAddress {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{provider}
}

// -- MsgUpdateService --

// Route Implements Msg
func (msg MsgUpdateService) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgUpdateService) Type() string { return TypeMsgUpdateService }

// GetSignBytes Implements Msg
func (msg MsgUpdateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic Implements Msg
func (msg MsgUpdateService) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid provider address")
	}
	if msg.ServiceId == "" {
		return sdkerrors.Wrap(ErrInvalidService, "service ID cannot be empty")
	}
	return nil
}

// GetSigners Implements Msg
func (msg MsgUpdateService) GetSigners() []sdk.AccAddress {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{provider}
}

// -- MsgActivateService --

// Route Implements Msg
func (msg MsgActivateService) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgActivateService) Type() string { return TypeMsgActivateService }

// GetSignBytes Implements Msg
func (msg MsgActivateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic Implements Msg
func (msg MsgActivateService) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid provider address")
	}
	if msg.ServiceId == "" {
		return sdkerrors.Wrap(ErrInvalidService, "service ID cannot be empty")
	}
	return nil
}

// GetSigners Implements Msg
func (msg MsgActivateService) GetSigners() []sdk.AccAddress {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{provider}
}

// -- MsgDeactivateService --

// Route Implements Msg
func (msg MsgDeactivateService) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDeactivateService) Type() string { return TypeMsgDeactivateService }

// GetSignBytes Implements Msg
func (msg MsgDeactivateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic Implements Msg
func (msg MsgDeactivateService) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid provider address")
	}
	if msg.ServiceId == "" {
		return sdkerrors.Wrap(ErrInvalidService, "service ID cannot be empty")
	}
	return nil
}

// GetSigners Implements Msg
func (msg MsgDeactivateService) GetSigners() []sdk.AccAddress {
	provider, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{provider}
}

// -- MsgPurchaseService --

// Route Implements Msg
func (msg MsgPurchaseService) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgPurchaseService) Type() string { return TypeMsgPurchaseService }

// GetSignBytes Implements Msg
func (msg MsgPurchaseService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic Implements Msg
func (msg MsgPurchaseService) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid buyer address")
	}
	if msg.ServiceId == "" {
		return sdkerrors.Wrap(ErrInvalidService, "service ID cannot be empty")
	}
	return nil
}

// GetSigners Implements Msg
func (msg MsgPurchaseService) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{buyer}
}
