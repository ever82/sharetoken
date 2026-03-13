package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Common validation constants
const (
	// MaxStringLength is the maximum allowed string length for general fields
	MaxStringLength = 1000

	// MaxNameLength is the maximum allowed length for name fields
	MaxNameLength = 100

	// MaxDescriptionLength is the maximum allowed length for description fields
	MaxDescriptionLength = 10000

	// MaxURLLength is the maximum allowed length for URL fields
	MaxURLLength = 2048

	// MinPasswordLength is the minimum required password length
	MinPasswordLength = 8
)

// Common patterns
var (
	// alphanumericPattern matches alphanumeric characters and common safe symbols
	alphanumericPattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

	// safeStringPattern allows common characters but blocks dangerous ones
	safeStringPattern = regexp.MustCompile(`^[\p{L}\p{N}\s_\-\.\,\!\?\(\)\[\]\{\}\/\:\;\"\'\@\#\$\%\&\*\+\=\<\>\|\\]+$`)

	// dangerousPattern matches potentially dangerous characters/patterns
	dangerousPattern = regexp.MustCompile(`[<>{})](?:script|iframe|object|embed|form)|javascript:|data:text/html|on\w+\s*=|<\s*\?|<%|` +
		`select\s+.*\s+from|insert\s+into|delete\s+from|drop\s+table|union\s+select`)
)

// ValidateStringLength validates that a string's length is within acceptable bounds
func ValidateStringLength(value string, fieldName string, minLen, maxLen int) error {
	length := utf8.RuneCountInString(value)
	if length < minLen {
		return fmt.Errorf("%s is too short: minimum %d characters, got %d", fieldName, minLen, length)
	}
	if length > maxLen {
		return fmt.Errorf("%s is too long: maximum %d characters, got %d", fieldName, maxLen, length)
	}
	return nil
}

// ValidateName validates a name field
func ValidateName(name string, fieldName string) error {
	if err := ValidateNonEmpty(name, fieldName); err != nil {
		return err
	}
	if err := ValidateStringLength(name, fieldName, 1, MaxNameLength); err != nil {
		return err
	}
	if !alphanumericPattern.MatchString(name) {
		return fmt.Errorf("%s contains invalid characters: only alphanumeric, underscore, hyphen, and period are allowed", fieldName)
	}
	return nil
}

// ValidateDescription validates a description field
func ValidateDescription(description string, fieldName string) error {
	if description == "" {
		return nil // Description is optional
	}
	if err := ValidateStringLength(description, fieldName, 0, MaxDescriptionLength); err != nil {
		return err
	}
	if err := ValidateSafeString(description, fieldName); err != nil {
		return err
	}
	return nil
}

// ValidateURL validates a URL field
func ValidateURL(urlStr string, fieldName string) error {
	if urlStr == "" {
		return nil // URL is optional
	}
	if err := ValidateStringLength(urlStr, fieldName, 1, MaxURLLength); err != nil {
		return err
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%s must use http or https scheme", fieldName)
	}
	if u.Host == "" {
		return fmt.Errorf("%s must have a valid host", fieldName)
	}
	return nil
}

// ValidateSafeString checks that a string doesn't contain dangerous characters or patterns
func ValidateSafeString(value string, fieldName string) error {
	if value == "" {
		return nil
	}
	if dangerousPattern.MatchString(strings.ToLower(value)) {
		return fmt.Errorf("%s contains potentially dangerous content", fieldName)
	}
	return nil
}

// ValidateAlphanumeric checks that a string contains only alphanumeric characters and safe symbols
func ValidateAlphanumeric(value string, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	if !alphanumericPattern.MatchString(value) {
		return fmt.Errorf("%s must contain only alphanumeric characters, underscore, hyphen, or period", fieldName)
	}
	return nil
}

// ValidateID validates an identifier (e.g., UUID, custom ID)
func ValidateID(id string, fieldName string) error {
	if err := ValidateNonEmpty(id, fieldName); err != nil {
		return err
	}
	if err := ValidateStringLength(id, fieldName, 1, 128); err != nil {
		return err
	}
	// IDs should be alphanumeric with limited safe characters
	idPattern := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !idPattern.MatchString(id) {
		return fmt.Errorf("%s contains invalid characters", fieldName)
	}
	return nil
}

// ValidateCoins validates a coins amount
func ValidateCoins(coins sdk.Coins, fieldName string) error {
	if coins.Empty() {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	if !coins.IsValid() {
		return fmt.Errorf("%s is not valid", fieldName)
	}
	if coins.IsAnyNegative() {
		return fmt.Errorf("%s cannot contain negative amounts", fieldName)
	}
	return nil
}

// ValidatePositiveCoins validates that coins are positive (greater than 0)
func ValidatePositiveCoins(coins sdk.Coins, fieldName string) error {
	if err := ValidateCoins(coins, fieldName); err != nil {
		return err
	}
	if coins.IsZero() {
		return fmt.Errorf("%s must be greater than 0", fieldName)
	}
	return nil
}

// ValidateDenoms validates that the provided denoms are valid
func ValidateDenoms(denoms []string, fieldName string) error {
	if len(denoms) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	seen := make(map[string]bool)
	for _, denom := range denoms {
		if denom == "" {
			return fmt.Errorf("%s contains empty denomination", fieldName)
		}
		if seen[denom] {
			return fmt.Errorf("%s contains duplicate denomination: %s", fieldName, denom)
		}
		seen[denom] = true
		// Validate denom format (basic check)
		if len(denom) > 128 {
			return fmt.Errorf("%s contains denomination that is too long: %s", fieldName, denom)
		}
	}
	return nil
}

// SanitizeString removes potentially dangerous characters from a string
func SanitizeString(input string) string {
	// Remove null bytes
	sanitized := strings.ReplaceAll(input, "\x00", "")
	// Remove control characters except common whitespace
	sanitized = regexp.MustCompile(`[\x01-\x08\x0B\x0C\x0E-\x1F\x7F]`).ReplaceAllString(sanitized, "")
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

// SanitizeLogValue prepares a value for safe logging
func SanitizeLogValue(input string) string {
	if len(input) > 200 {
		input = input[:200] + "..."
	}
	// Remove newlines to prevent log injection
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\r", "\\r")
	return input
}

// ValidatePassword validates a password meets minimum requirements
func ValidatePassword(password string, fieldName string) error {
	if err := ValidateStringLength(password, fieldName, MinPasswordLength, 128); err != nil {
		return err
	}
	return nil
}

// ValidateHexString validates a hexadecimal string
func ValidateHexString(hex string, fieldName string) error {
	if hex == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	hexPattern := regexp.MustCompile(`^(0x)?[0-9a-fA-F]+$`)
	if !hexPattern.MatchString(hex) {
		return fmt.Errorf("%s is not a valid hexadecimal string", fieldName)
	}
	return nil
}

// ValidateJSON validates that a string is valid JSON (basic check)
func ValidateJSON(jsonStr string, fieldName string) error {
	if jsonStr == "" {
		return nil // Empty is considered valid
	}
	if !strings.HasPrefix(jsonStr, "{") && !strings.HasPrefix(jsonStr, "[") {
		return fmt.Errorf("%s is not valid JSON", fieldName)
	}
	return nil
}
