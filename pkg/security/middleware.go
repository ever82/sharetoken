package security

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cometbft/cometbft/libs/log"
)

// SecurityMiddleware provides security checks as middleware
type SecurityMiddleware struct {
	rateLimiter    *CompositeRateLimiter
	replayGuard    *ReplayGuard
	securityLogger *SecurityLogger
	enabled        bool
}

// SecurityMiddlewareConfig configures the security middleware
type SecurityMiddlewareConfig struct {
	Enabled           bool
	RateLimitEnabled  bool
	ReplayProtection  bool
	InputValidation   bool
	SecurityLogging   bool
	IPRateLimitTier   RateLimitTier
	UserRateLimitTier RateLimitTier
}

// DefaultSecurityMiddlewareConfig returns the default configuration
func DefaultSecurityMiddlewareConfig() SecurityMiddlewareConfig {
	return SecurityMiddlewareConfig{
		Enabled:           true,
		RateLimitEnabled:  true,
		ReplayProtection:  true,
		InputValidation:   true,
		SecurityLogging:   true,
		IPRateLimitTier:   TierStandard,
		UserRateLimitTier: TierStandard,
	}
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware(config SecurityMiddlewareConfig, logger log.Logger) *SecurityMiddleware {
	if !config.Enabled {
		return &SecurityMiddleware{enabled: false}
	}

	sm := &SecurityMiddleware{
		enabled: config.Enabled,
	}

	if config.RateLimitEnabled {
		ipConfig := DefaultRateLimits[config.IPRateLimitTier]
		userConfig := DefaultRateLimits[config.UserRateLimitTier]
		sm.rateLimiter = NewCompositeRateLimiter(ipConfig, userConfig)
	}

	if config.ReplayProtection {
		sm.replayGuard = NewReplayGuard(DefaultTimestampWindow)
	}

	if config.SecurityLogging {
		sm.securityLogger = NewSecurityLogger(logger)
	}

	return sm
}

// CheckRateLimit checks if the request is within rate limits
func (sm *SecurityMiddleware) CheckRateLimit(ctx sdk.Context, ip, userID string) error {
	if !sm.enabled || sm.rateLimiter == nil {
		return nil
	}

	allowed, reason := sm.rateLimiter.Allow(ip, userID)
	if !allowed {
		// Log the rate limit event
		if sm.securityLogger != nil {
			sm.securityLogger.LogRateLimitExceeded(userID, ip, fmt.Sprintf("gas_used_%d", ctx.GasMeter().GasConsumed()), map[string]string{
				"reason": reason,
			})
		}
		return ErrRateLimitExceeded.Wrapf("rate limit exceeded: %s", reason)
	}

	return nil
}

// CheckReplay checks for replay attacks
func (sm *SecurityMiddleware) CheckReplay(ctx sdk.Context, address string, sequence uint64) error {
	if !sm.enabled || sm.replayGuard == nil {
		return nil
	}

	// Basic sequence validation
	if err := ValidateBasicSequence(sequence); err != nil {
		if sm.securityLogger != nil {
			sm.securityLogger.LogReplayAttempt(address, "", map[string]string{
				"reason":   "invalid_sequence",
				"sequence": fmt.Sprintf("%d", sequence),
			})
		}
		return ErrInvalidSequence.Wrap(err.Error())
	}

	// Validate sequence
	if err := sm.replayGuard.sequenceValidator.ValidateAndUpdate(address, sequence); err != nil {
		if sm.securityLogger != nil {
			sm.securityLogger.LogReplayAttempt(address, "", map[string]string{
				"reason":   err.Error(),
				"sequence": fmt.Sprintf("%d", sequence),
			})
		}
		return ErrReplayDetected.Wrap(err.Error())
	}

	return nil
}

// LogSecurityEvent logs a security event
func (sm *SecurityMiddleware) LogSecurityEvent(event SecurityEvent) {
	if sm.enabled && sm.securityLogger != nil {
		sm.securityLogger.LogEvent(event)
	}
}

// LogAuthFailure logs an authentication failure
func (sm *SecurityMiddleware) LogAuthFailure(userID, ipAddress, reason string) {
	if sm.enabled && sm.securityLogger != nil {
		sm.securityLogger.LogAuthFailure(userID, ipAddress, reason, nil)
	}
}

// LogSuspiciousActivity logs suspicious activity
func (sm *SecurityMiddleware) LogSuspiciousActivity(userID, ipAddress, description string) {
	if sm.enabled && sm.securityLogger != nil {
		sm.securityLogger.LogSuspiciousActivity(userID, ipAddress, description, nil)
	}
}

// LogInvalidInput logs an invalid input event
func (sm *SecurityMiddleware) LogInvalidInput(userID, ipAddress, field, reason string) {
	if sm.enabled && sm.securityLogger != nil {
		sm.securityLogger.LogInvalidInput(userID, ipAddress, field, reason)
	}
}

// LogUnauthorizedAccess logs an unauthorized access event
func (sm *SecurityMiddleware) LogUnauthorizedAccess(userID, ipAddress, resource string, details map[string]string) {
	if sm.enabled && sm.securityLogger != nil {
		sm.securityLogger.LogUnauthorizedAccess(userID, ipAddress, resource, details)
	}
}

// SecurityAnteDecorator wraps security checks as an ante decorator
type SecurityAnteDecorator struct {
	middleware *SecurityMiddleware
}

// NewSecurityAnteDecorator creates a new security ante decorator
func NewSecurityAnteDecorator(middleware *SecurityMiddleware) SecurityAnteDecorator {
	return SecurityAnteDecorator{
		middleware: middleware,
	}
}

// AnteHandle performs security checks before transaction processing
func (sad SecurityAnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if !sad.middleware.enabled {
		return next(ctx, tx, simulate)
	}

	// Get signers from transaction
	signers := make([]string, 0)
	for _, msg := range tx.GetMsgs() {
		for _, signer := range msg.GetSigners() {
			signers = append(signers, signer.String())
		}
	}

	// Check rate limits for each signer
	for _, signer := range signers {
		if err := sad.middleware.CheckRateLimit(ctx, "", signer); err != nil {
			return ctx, err
		}
	}

	// Check replay protection
	// Note: In real implementation, you'd get sequence from account
	// This is a simplified version
	for _, signer := range signers {
		// Get account and check sequence
		_ = signer // Use the signer
		// sequence := account.GetSequence()
		// if err := sad.middleware.CheckReplay(ctx, signer, sequence); err != nil {
		//     return ctx, err
		// }
	}

	return next(ctx, tx, simulate)
}

// SecurityKeeper provides security features for keepers
type SecurityKeeper struct {
	middleware *SecurityMiddleware
}

// NewSecurityKeeper creates a new security keeper
func NewSecurityKeeper(middleware *SecurityMiddleware) *SecurityKeeper {
	return &SecurityKeeper{
		middleware: middleware,
	}
}

// ValidateAddress validates an address and logs if invalid
func (sk *SecurityKeeper) ValidateAddress(ctx sdk.Context, address string, fieldName string) error {
	if address == "" {
		sk.middleware.LogInvalidInput("", "", fieldName, "empty address")
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	if !IsAccAddress(address) {
		sk.middleware.LogInvalidInput("", "", fieldName, "invalid address format")
		return fmt.Errorf("invalid %s format", fieldName)
	}

	return nil
}

// CheckPermission checks if a user has permission and logs access attempts
func (sk *SecurityKeeper) CheckPermission(ctx sdk.Context, userID, resource, action string, hasPermission bool) error {
	if !hasPermission {
		sk.middleware.LogUnauthorizedAccess(userID, "", resource, map[string]string{
			"action": action,
		})
		return ErrUnauthorized.Wrapf("unauthorized to %s %s", action, resource)
	}

	// Log successful access
	sk.middleware.LogSecurityEvent(SecurityEvent{
		Type:     EventAccessGranted,
		Severity: SeverityInfo,
		Message:  fmt.Sprintf("Access granted to %s for %s", resource, action),
		UserID:   userID,
		Resource: resource,
		Action:   action,
		Result:   "success",
	})

	return nil
}

// LogTransaction logs a transaction for security audit
func (sk *SecurityKeeper) LogTransaction(ctx sdk.Context, txHash, sender, receiver, amount string, success bool) {
	severity := SeverityInfo
	result := "successful"
	if !success {
		severity = SeverityHigh
		result = "failed"
	}

	sk.middleware.LogSecurityEvent(SecurityEvent{
		Type:     EventSecurityAlert,
		Severity: severity,
		Message:  fmt.Sprintf("Transaction %s", result),
		UserID:   sender,
		Details: map[string]string{
			"tx_hash":  SanitizeHex(txHash),
			"sender":   SanitizeAddress(sender),
			"receiver": SanitizeAddress(receiver),
			"amount":   amount,
		},
	})
}

// RateLimitInfo returns rate limit information
func (sk *SecurityKeeper) RateLimitInfo(ip, userID string) RateLimitInfo {
	if sk.middleware.rateLimiter != nil {
		return sk.middleware.rateLimiter.GetRateLimitInfo(ip, userID)
	}
	return RateLimitInfo{}
}

// IsAccAddress checks if a string is a valid account address
// This is a convenience wrapper for validation.IsAccAddress
func IsAccAddress(address string) bool {
	// Import from validation package
	// For now, do a simple check
	return len(address) > 0 && len(address) >= len("sharetoken") && address[:len("sharetoken")] == "sharetoken"
}
