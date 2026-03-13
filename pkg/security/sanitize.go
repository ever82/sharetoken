package security

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// SanitizationConfig configures sanitization behavior
type SanitizationConfig struct {
	// MaskLength is the number of characters to show at the beginning and end
	MaskLength int
	// MaskChar is the character used for masking
	MaskChar rune
	// EnableFullMask masks the entire value if true
	EnableFullMask bool
}

// DefaultSanitizationConfig returns the default configuration
func DefaultSanitizationConfig() SanitizationConfig {
	return SanitizationConfig{
		MaskLength:     4,
		MaskChar:       '*',
		EnableFullMask: false,
	}
}

// StrictSanitizationConfig returns a strict configuration (full masking)
func StrictSanitizationConfig() SanitizationConfig {
	return SanitizationConfig{
		MaskLength:     0,
		MaskChar:       '*',
		EnableFullMask: true,
	}
}

// SanitizeAPIKey masks an API key for safe logging
// Shows first and last MaskLength characters, masks the rest
func SanitizeAPIKey(apiKey string) string {
	return SanitizeAPIKeyWithConfig(apiKey, DefaultSanitizationConfig())
}

// SanitizeAPIKeyWithConfig masks an API key with custom configuration
func SanitizeAPIKeyWithConfig(apiKey string, config SanitizationConfig) string {
	if apiKey == "" {
		return ""
	}

	if config.EnableFullMask {
		return strings.Repeat(string(config.MaskChar), len(apiKey))
	}

	if len(apiKey) <= config.MaskLength*2 {
		// Too short, mask everything except first char
		if len(apiKey) > 1 {
			return apiKey[:1] + strings.Repeat(string(config.MaskChar), len(apiKey)-1)
		}
		return strings.Repeat(string(config.MaskChar), len(apiKey))
	}

	// Show first and last MaskLength chars
	prefix := apiKey[:config.MaskLength]
	suffix := apiKey[len(apiKey)-config.MaskLength:]
	middleLen := len(apiKey) - config.MaskLength*2
	middle := strings.Repeat(string(config.MaskChar), middleLen)

	return prefix + middle + suffix
}

// SanitizePrivateKey masks a private key for safe logging
func SanitizePrivateKey(key string) string {
	return SanitizePrivateKeyWithConfig(key, StrictSanitizationConfig())
}

// SanitizePrivateKeyWithConfig masks a private key with custom configuration
func SanitizePrivateKeyWithConfig(key string, config SanitizationConfig) string {
	if key == "" {
		return ""
	}

	if config.EnableFullMask {
		return "[REDACTED]"
	}

	return SanitizeAPIKeyWithConfig(key, config)
}

// SanitizeMnemonic masks a mnemonic phrase for safe logging
func SanitizeMnemonic(mnemonic string) string {
	if mnemonic == "" {
		return ""
	}

	words := strings.Fields(mnemonic)
	if len(words) == 0 {
		return "[REDACTED]"
	}

	// Show only the count of words
	return fmt.Sprintf("[%d words mnemonic - REDACTED]", len(words))
}

// SanitizePassword masks a password completely
func SanitizePassword(password string) string {
	if password == "" {
		return ""
	}
	return "[PASSWORD REDACTED]"
}

// SanitizeToken masks an authentication token
func SanitizeToken(token string) string {
	return SanitizeTokenWithConfig(token, DefaultSanitizationConfig())
}

// SanitizeTokenWithConfig masks a token with custom configuration
func SanitizeTokenWithConfig(token string, config SanitizationConfig) string {
	if token == "" {
		return ""
	}

	// Detect JWT tokens (contain two dots)
	if strings.Count(token, ".") == 2 {
		parts := strings.Split(token, ".")
		// Mask the payload and signature
		return parts[0] + "." +
			SanitizeAPIKeyWithConfig(parts[1], config) + "." +
			SanitizeAPIKeyWithConfig(parts[2], config)
	}

	return SanitizeAPIKeyWithConfig(token, config)
}

// SanitizeAddress masks a blockchain address
func SanitizeAddress(address string) string {
	if address == "" {
		return ""
	}

	if len(address) <= 10 {
		return strings.Repeat("*", len(address))
	}

	// Show first 6 and last 4 characters
	return address[:6] + "..." + address[len(address)-4:]
}

// SanitizeHex masks a hex string (like a hash or transaction ID)
func SanitizeHex(hex string) string {
	if hex == "" {
		return ""
	}

	// Remove 0x prefix if present
	cleanHex := hex
	if strings.HasPrefix(hex, "0x") {
		cleanHex = hex[2:]
	}

	if len(cleanHex) <= 8 {
		return "0x" + strings.Repeat("*", len(cleanHex))
	}

	// Show first 4 and last 4 characters
	return "0x" + cleanHex[:4] + "..." + cleanHex[len(cleanHex)-4:]
}

// SensitiveField represents a field that needs sanitization
type SensitiveField struct {
	Name        string
	Value       string
	Type        SensitiveType
	SanitizationConfig
}

// SensitiveType defines the type of sensitive data
type SensitiveType int

const (
	TypeAPIKey SensitiveType = iota
	TypePrivateKey
	TypeMnemonic
	TypePassword
	TypeToken
	TypeAddress
	TypeHex
	TypeGeneric
)

// Sanitize sanitizes a sensitive field based on its type
func (sf SensitiveField) Sanitize() string {
	switch sf.Type {
	case TypeAPIKey:
		return SanitizeAPIKeyWithConfig(sf.Value, sf.SanitizationConfig)
	case TypePrivateKey:
		return SanitizePrivateKeyWithConfig(sf.Value, sf.SanitizationConfig)
	case TypeMnemonic:
		return SanitizeMnemonic(sf.Value)
	case TypePassword:
		return SanitizePassword(sf.Value)
	case TypeToken:
		return SanitizeTokenWithConfig(sf.Value, sf.SanitizationConfig)
	case TypeAddress:
		return SanitizeAddress(sf.Value)
	case TypeHex:
		return SanitizeHex(sf.Value)
	case TypeGeneric:
		return SanitizeAPIKeyWithConfig(sf.Value, sf.SanitizationConfig)
	default:
		return SanitizeAPIKeyWithConfig(sf.Value, sf.SanitizationConfig)
	}
}

// Sanitizer provides a structured way to sanitize multiple fields
type Sanitizer struct {
	config SanitizationConfig
	fields []SensitiveField
}

// NewSanitizer creates a new sanitizer with default config
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		config: DefaultSanitizationConfig(),
		fields: []SensitiveField{},
	}
}

// NewSanitizerWithConfig creates a new sanitizer with custom config
func NewSanitizerWithConfig(config SanitizationConfig) *Sanitizer {
	return &Sanitizer{
		config: config,
		fields: []SensitiveField{},
	}
}

// AddField adds a field to be sanitized
func (s *Sanitizer) AddField(name, value string, fieldType SensitiveType) *Sanitizer {
	s.fields = append(s.fields, SensitiveField{
		Name:               name,
		Value:              value,
		Type:               fieldType,
		SanitizationConfig: s.config,
	})
	return s
}

// AddFieldWithConfig adds a field with custom sanitization config
func (s *Sanitizer) AddFieldWithConfig(name, value string, fieldType SensitiveType, config SanitizationConfig) *Sanitizer {
	s.fields = append(s.fields, SensitiveField{
		Name:               name,
		Value:              value,
		Type:               fieldType,
		SanitizationConfig: config,
	})
	return s
}

// SanitizeAll returns a map of sanitized field values
func (s *Sanitizer) SanitizeAll() map[string]string {
	result := make(map[string]string)
	for _, field := range s.fields {
		result[field.Name] = field.Sanitize()
	}
	return result
}

// SanitizeString is a generic function that detects and sanitizes sensitive data
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Check for common patterns

	// API Key patterns (sk-*, pk-*, etc.)
	if regexp.MustCompile(`^(sk|pk)_[a-zA-Z0-9]{20,}$`).MatchString(input) {
		return SanitizeAPIKey(input)
	}

	// Hex strings (64 chars = private key, 40 chars = address)
	cleanInput := strings.TrimPrefix(input, "0x")
	if regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(cleanInput) {
		return "[HEX64 REDACTED]"
	}
	if regexp.MustCompile(`^[0-9a-fA-F]{40}$`).MatchString(cleanInput) {
		return SanitizeAddress(input)
	}

	// Base64 encoded data
	if regexp.MustCompile(`^[A-Za-z0-9+/]{40,}={0,2}$`).MatchString(input) {
		return SanitizeAPIKey(input)
	}

	// Mnemonic (12-24 words)
	words := strings.Fields(input)
	if len(words) >= 12 && len(words) <= 24 {
		// Check if all words are lowercase letters
		isPotentialMnemonic := true
		for _, word := range words {
			if !regexp.MustCompile(`^[a-z]+$`).MatchString(word) {
				isPotentialMnemonic = false
				break
			}
		}
		if isPotentialMnemonic {
			return SanitizeMnemonic(input)
		}
	}

	// Default: return as-is if no pattern matches
	return input
}

// SanitizeMap sanitizes all values in a map
func SanitizeMap(input map[string]string, sensitiveKeys []string) map[string]string {
	sensitiveSet := make(map[string]bool)
	for _, key := range sensitiveKeys {
		sensitiveSet[strings.ToLower(key)] = true
	}

	result := make(map[string]string)
	for key, value := range input {
		if sensitiveSet[strings.ToLower(key)] {
			result[key] = SanitizeString(value)
		} else {
			result[key] = value
		}
	}
	return result
}

// DefaultSensitiveKeys is a list of commonly sensitive field names
var DefaultSensitiveKeys = []string{
	"apikey", "api_key", "api-key",
	"secret", "secretkey", "secret_key", "secret-key",
	"privatekey", "private_key", "private-key",
	"password", "passwd", "pwd",
	"token", "accesstoken", "access_token", "access-token",
	"refreshtoken", "refresh_token", "refresh-token",
	"mnemonic", "seed",
	"authorization", "auth",
}

// IsSensitiveField checks if a field name is likely sensitive
func IsSensitiveField(fieldName string) bool {
	lowerName := strings.ToLower(fieldName)
	for _, key := range DefaultSensitiveKeys {
		if strings.Contains(lowerName, key) {
			return true
		}
	}
	return false
}

// SanitizeBase64 masks base64 encoded data
func SanitizeBase64(data string) string {
	if data == "" {
		return ""
	}

	// Try to decode to check if it's valid base64
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// Not valid base64, treat as regular string
		return SanitizeString(data)
	}

	// If it's valid base64, sanitize the decoded content
	sanitized := SanitizeString(string(decoded))

	// Re-encode if needed, or return masked indicator
	if sanitized == string(decoded) {
		// No sensitive data found, return original
		return data
	}

	return "[BASE64 ENCODED SENSITIVE DATA]"
}
