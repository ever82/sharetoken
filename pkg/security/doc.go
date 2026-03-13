// Package security provides security-related utilities for the ShareToken blockchain.
//
// This package implements the security measures required by GitHub Issue #58.
//
// Security Measures Implemented:
//
// 1. Input Validation and Sanitization
//   - See pkg/validation/ for comprehensive input validation
//   - Includes address validation, string sanitization, safe string checking
//   - URL validation, coins validation, and password strength validation
//
// 2. Replay Attack Prevention
//   - NonceTracker: Tracks used nonces to prevent replay attacks
//   - SequenceValidator: Validates account sequences
//   - TimestampValidator: Validates transaction timestamps
//   - ReplayGuard: Comprehensive replay attack protection combining all mechanisms
//   - See replay.go for implementation
//
// 3. Rate Limiting
//   - RateLimiter: Token bucket rate limiter
//   - IPRateLimiter: Rate limiting by IP address
//   - UserRateLimiter: Rate limiting by user ID
//   - CompositeRateLimiter: Combines IP and user rate limiting
//   - See ratelimit.go for implementation
//
// 4. Sensitive Data Sanitization
//   - SanitizeAPIKey: Masks API keys for safe logging
//   - SanitizePrivateKey: Masks private keys
//   - SanitizeMnemonic: Masks mnemonic phrases
//   - SanitizeToken: Masks authentication tokens
//   - SanitizeAddress: Masks blockchain addresses
//   - SanitizeHex: Masks hex strings (hashes, tx IDs)
//   - See sanitize.go for implementation
//
// 5. Security Logging
//   - SecurityLogger: Security event logging with severity levels
//   - SecurityAuditLogger: Audit logging for transactions and data access
//   - Event types: AuthFailure, AuthSuccess, AccessDenied, RateLimitExceeded, etc.
//   - See logger.go for implementation
//
// All security functions are designed to be composable and reusable across modules.
package security
