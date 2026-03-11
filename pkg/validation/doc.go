// Package validation provides common validation functions for Cosmos SDK modules.
//
// This package contains reusable validation functions for common data types:
//
//   - Address validation (Bech32 format)
//   - String validation (non-empty)
//   - Numeric validation (positive, non-negative, range checks)
//   - Signer extraction
//
// All validation functions return detailed error messages that include the field name.
//
// Example usage:
//
//	import "sharetoken/pkg/validation"
//
//	func (msg MsgMyMessage) ValidateBasic() error {
//	    if err := validation.ValidateAddress(msg.Sender, "sender"); err != nil {
//	        return err
//	    }
//	    if err := validation.ValidateNonEmpty(msg.Name, "name"); err != nil {
//	        return err
//	    }
//	    if err := validation.ValidatePositiveUint64(msg.Amount, "amount"); err != nil {
//	        return err
//	    }
//	    return nil
//	}
//
//	func (msg MsgMyMessage) GetSigners() []sdk.AccAddress {
//	    return validation.MustGetSigners(msg.Sender)
//	}
//
package validation
