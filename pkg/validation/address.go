package validation

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ShareTokenPrefix is the expected Bech32 prefix for ShareToken addresses
const ShareTokenPrefix = "sharetoken"

var (
	// addressRegex matches ShareToken Bech32 addresses
	// sharetoken + 1 (separator) + 39 (data) + 6 (checksum) = ~44-51 chars
	accAddressRegex  = regexp.MustCompile(fmt.Sprintf(`^%s1[a-z0-9]{38,40}$`, ShareTokenPrefix))
	valAddressRegex  = regexp.MustCompile(fmt.Sprintf(`^%svaloper1[a-z0-9]{38,40}$`, ShareTokenPrefix))
	consAddressRegex = regexp.MustCompile(fmt.Sprintf(`^%svalcons1[a-z0-9]{38,40}$`, ShareTokenPrefix))
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

// ValidateNonEmpty checks that a string field is not empty.
func ValidateNonEmpty(value string, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
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

// ValidateAccAddress validates a ShareToken account address
func ValidateAccAddress(address string, fieldName string) error {
	if address == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	if !accAddressRegex.MatchString(address) {
		return fmt.Errorf("invalid %s format: expected %s1... format", fieldName, ShareTokenPrefix)
	}
	return nil
}

// ValidateAccAddressOptional validates an account address that may be empty
func ValidateAccAddressOptional(address string, fieldName string) error {
	if address == "" {
		return nil
	}
	return ValidateAccAddress(address, fieldName)
}

// ValidateValAddress validates a ShareToken validator address
func ValidateValAddress(address string, fieldName string) error {
	if address == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	if !valAddressRegex.MatchString(address) {
		return fmt.Errorf("invalid %s format: expected %svaloper1... format", fieldName, ShareTokenPrefix)
	}
	return nil
}

// ValidateValAddressOptional validates a validator address that may be empty
func ValidateValAddressOptional(address string, fieldName string) error {
	if address == "" {
		return nil
	}
	return ValidateValAddress(address, fieldName)
}

// ValidateConsAddress validates a ShareToken consensus address
func ValidateConsAddress(address string, fieldName string) error {
	if address == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	if !consAddressRegex.MatchString(address) {
		return fmt.Errorf("invalid %s format: expected %svalcons1... format", fieldName, ShareTokenPrefix)
	}
	return nil
}

// ValidateConsAddressOptional validates a consensus address that may be empty
func ValidateConsAddressOptional(address string, fieldName string) error {
	if address == "" {
		return nil
	}
	return ValidateConsAddress(address, fieldName)
}

// ValidateShareTokenAddress validates any type of ShareToken address
func ValidateShareTokenAddress(address string, fieldName string) error {
	if address == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	// Check if it matches any ShareToken address format
	if accAddressRegex.MatchString(address) ||
		valAddressRegex.MatchString(address) ||
		consAddressRegex.MatchString(address) {
		return nil
	}

	return fmt.Errorf("invalid %s format: expected ShareToken address format", fieldName)
}

// IsAccAddress checks if a string is a valid account address format
func IsAccAddress(address string) bool {
	return accAddressRegex.MatchString(address)
}

// IsValAddress checks if a string is a valid validator address format
func IsValAddress(address string) bool {
	return valAddressRegex.MatchString(address)
}

// IsConsAddress checks if a string is a valid consensus address format
func IsConsAddress(address string) bool {
	return consAddressRegex.MatchString(address)
}

// ValidateAddressList validates a list of addresses
func ValidateAddressList(addresses []string, fieldName string) error {
	if len(addresses) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	seen := make(map[string]bool)
	for i, addr := range addresses {
		if err := ValidateAccAddress(addr, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
		if seen[addr] {
			return fmt.Errorf("%s contains duplicate address: %s", fieldName, addr)
		}
		seen[addr] = true
	}

	return nil
}
