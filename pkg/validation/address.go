package validation

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateAddress validates that the given string is a valid Bech32 address.
// Returns an error if the address is empty or invalid.
func ValidateAddress(address string, fieldName string) error {
	if address == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	return nil
}

// ValidateAddressOptional validates an address that may be empty.
// Only returns an error if the address is non-empty and invalid.
func ValidateAddressOptional(address string, fieldName string) error {
	if address == "" {
		return nil
	}
	return ValidateAddress(address, fieldName)
}

// MustGetSigners extracts signers from an address string.
// Returns nil if the address is invalid (for use in GetSigners methods).
func MustGetSigners(address string) []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateNonEmpty checks that a string field is not empty.
func ValidateNonEmpty(value string, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidatePositive checks that an int64 value is positive (greater than 0).
func ValidatePositive(value int64, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than 0", fieldName)
	}
	return nil
}

// ValidatePositiveUint64 checks that a uint64 value is positive (greater than 0).
func ValidatePositiveUint64(value uint64, fieldName string) error {
	if value == 0 {
		return fmt.Errorf("%s must be greater than 0", fieldName)
	}
	return nil
}

// ValidateNonNegative checks that an int64 value is non-negative.
func ValidateNonNegative(value int64, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative", fieldName)
	}
	return nil
}

// ValidateRange checks that an int64 value is within the given range (inclusive).
func ValidateRange(value int64, min, max int64, fieldName string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d", fieldName, min, max)
	}
	return nil
}

// ValidateRangeInt32 checks that an int32 value is within the given range (inclusive).
func ValidateRangeInt32(value int32, min, max int32, fieldName string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d", fieldName, min, max)
	}
	return nil
}
