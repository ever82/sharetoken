package types

import (
	"errors"
)

var (
	// ErrAPIKeyNotFound is returned when an API key is not found
	ErrAPIKeyNotFound = errors.New("API key not found")

	// ErrAPIKeyExists is returned when an API key already exists
	ErrAPIKeyExists = errors.New("API key already exists")

	// ErrInvalidProvider is returned when the provider is invalid
	ErrInvalidProvider = errors.New("invalid provider")

	// ErrInvalidAPIKey is returned when the API key is invalid
	ErrInvalidAPIKey = errors.New("invalid API key")

	// ErrUnauthorized is returned when the caller is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrQuotaExceeded is returned when the usage quota is exceeded
	ErrQuotaExceeded = errors.New("usage quota exceeded")

	// ErrRateLimitExceeded is returned when the rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrEncryptionFailed is returned when encryption fails
	ErrEncryptionFailed = errors.New("encryption failed")

	// ErrDecryptionFailed is returned when decryption fails
	ErrDecryptionFailed = errors.New("decryption failed")

	// ErrAccessDenied is returned when access is denied
	ErrAccessDenied = errors.New("access denied")
)
