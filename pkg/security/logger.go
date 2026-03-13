package security

import (
	"fmt"
	"os"
	"time"

	"github.com/cometbft/cometbft/libs/log"
)

// SecurityEventType defines types of security events
type SecurityEventType string

const (
	// EventAuthFailure authentication failure
	EventAuthFailure SecurityEventType = "AUTH_FAILURE"
	// EventAuthSuccess successful authentication
	EventAuthSuccess SecurityEventType = "AUTH_SUCCESS"
	// EventAccessDenied access denied
	EventAccessDenied SecurityEventType = "ACCESS_DENIED"
	// EventAccessGranted access granted
	EventAccessGranted SecurityEventType = "ACCESS_GRANTED"
	// EventSuspiciousActivity suspicious activity detected
	EventSuspiciousActivity SecurityEventType = "SUSPICIOUS_ACTIVITY"
	// EventRateLimitExceeded rate limit exceeded
	EventRateLimitExceeded SecurityEventType = "RATE_LIMIT_EXCEEDED"
	// EventReplayAttempt replay attack attempt
	EventReplayAttempt SecurityEventType = "REPLAY_ATTEMPT"
	// EventInvalidInput invalid input detected
	EventInvalidInput SecurityEventType = "INVALID_INPUT"
	// EventUnauthorizedAccess unauthorized access attempt
	EventUnauthorizedAccess SecurityEventType = "UNAUTHORIZED_ACCESS"
	// EventConfigChange configuration change
	EventConfigChange SecurityEventType = "CONFIG_CHANGE"
	// EventSecurityAlert general security alert
	EventSecurityAlert SecurityEventType = "SECURITY_ALERT"
)

// SeverityLevel defines the severity of a security event
type SeverityLevel string

const (
	// SeverityCritical critical security event
	SeverityCritical SeverityLevel = "CRITICAL"
	// SeverityHigh high severity security event
	SeverityHigh SeverityLevel = "HIGH"
	// SeverityMedium medium severity security event
	SeverityMedium SeverityLevel = "MEDIUM"
	// SeverityLow low severity security event
	SeverityLow SeverityLevel = "LOW"
	// SeverityInfo informational security event
	SeverityInfo SeverityLevel = "INFO"
)

// SecurityEvent represents a security event log entry
type SecurityEvent struct {
	Timestamp   time.Time         `json:"timestamp"`
	Type        SecurityEventType `json:"type"`
	Severity    SeverityLevel     `json:"severity"`
	Message     string            `json:"message"`
	UserID      string            `json:"user_id,omitempty"`
	IPAddress   string            `json:"ip_address,omitempty"`
	Resource    string            `json:"resource,omitempty"`
	Action      string            `json:"action,omitempty"`
	Result      string            `json:"result,omitempty"`
	Details     map[string]string `json:"details,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
	SessionID   string            `json:"session_id,omitempty"`
}

// SecurityLogger provides security event logging capabilities
type SecurityLogger struct {
	logger    log.Logger
	sanitizer *Sanitizer
	logFile   *os.File
}

// NewSecurityLogger creates a new security logger
func NewSecurityLogger(baseLogger log.Logger) *SecurityLogger {
	return &SecurityLogger{
		logger:    baseLogger.With("module", "security"),
		sanitizer: NewSanitizer(),
	}
}

// NewSecurityLoggerWithFile creates a security logger with file output
func NewSecurityLoggerWithFile(baseLogger log.Logger, filePath string) (*SecurityLogger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open security log file: %w", err)
	}

	return &SecurityLogger{
		logger:    baseLogger.With("module", "security"),
		sanitizer: NewSanitizer(),
		logFile:   file,
	}, nil
}

// Close closes the security logger
func (sl *SecurityLogger) Close() error {
	if sl.logFile != nil {
		return sl.logFile.Close()
	}
	return nil
}

// sanitizeDetails sanitizes sensitive information from event details
func (sl *SecurityLogger) sanitizeDetails(details map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for key, value := range details {
		if IsSensitiveField(key) {
			sanitized[key] = SanitizeString(value)
		} else {
			sanitized[key] = value
		}
	}
	return sanitized
}

// LogEvent logs a security event
func (sl *SecurityLogger) LogEvent(event SecurityEvent) {
	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Sanitize sensitive details
	event.Details = sl.sanitizeDetails(event.Details)

	// Build log message
	msg := fmt.Sprintf("[SECURITY] %s | %s | %s",
		event.Severity,
		event.Type,
		event.Message,
	)

	// Add context fields
	fields := []interface{}{
		"timestamp", event.Timestamp.Format(time.RFC3339),
		"type", string(event.Type),
		"severity", string(event.Severity),
	}

	if event.UserID != "" {
		fields = append(fields, "user_id", SanitizeAddress(event.UserID))
	}
	if event.IPAddress != "" {
		fields = append(fields, "ip_address", event.IPAddress)
	}
	if event.Resource != "" {
		fields = append(fields, "resource", event.Resource)
	}
	if event.Action != "" {
		fields = append(fields, "action", event.Action)
	}
	if event.Result != "" {
		fields = append(fields, "result", event.Result)
	}
	if event.RequestID != "" {
		fields = append(fields, "request_id", event.RequestID)
	}
	if event.SessionID != "" {
		fields = append(fields, "session_id", event.SessionID)
	}
	if len(event.Details) > 0 {
		fields = append(fields, "details", event.Details)
	}

	// Log based on severity
	switch event.Severity {
	case SeverityCritical:
		sl.logger.Error(msg, fields...)
	case SeverityHigh:
		sl.logger.Error(msg, fields...)
	case SeverityMedium:
		sl.logger.Info(msg, fields...)
	case SeverityLow:
		sl.logger.Debug(msg, fields...)
	default:
		sl.logger.Info(msg, fields...)
	}

	// Write to dedicated security log file if configured
	if sl.logFile != nil {
		logLine := fmt.Sprintf("%s | %s | %s | %s | %s\n",
			event.Timestamp.Format(time.RFC3339),
			event.Severity,
			event.Type,
			event.UserID,
			event.Message,
		)
		sl.logFile.WriteString(logLine)
	}
}

// LogAuthFailure logs an authentication failure event
func (sl *SecurityLogger) LogAuthFailure(userID, ipAddress, reason string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventAuthFailure,
		Severity:  SeverityHigh,
		Message:   fmt.Sprintf("Authentication failed: %s", reason),
		UserID:    userID,
		IPAddress: ipAddress,
		Action:    "authenticate",
		Result:    "failure",
		Details:   details,
	})
}

// LogAuthSuccess logs a successful authentication event
func (sl *SecurityLogger) LogAuthSuccess(userID, ipAddress string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventAuthSuccess,
		Severity:  SeverityInfo,
		Message:   "Authentication successful",
		UserID:    userID,
		IPAddress: ipAddress,
		Action:    "authenticate",
		Result:    "success",
		Details:   details,
	})
}

// LogAccessDenied logs an access denied event
func (sl *SecurityLogger) LogAccessDenied(userID, ipAddress, resource, reason string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventAccessDenied,
		Severity:  SeverityMedium,
		Message:   fmt.Sprintf("Access denied to %s: %s", resource, reason),
		UserID:    userID,
		IPAddress: ipAddress,
		Resource:  resource,
		Action:    "access",
		Result:    "denied",
	})
}

// LogRateLimitExceeded logs a rate limit exceeded event
func (sl *SecurityLogger) LogRateLimitExceeded(userID, ipAddress, resource string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventRateLimitExceeded,
		Severity:  SeverityMedium,
		Message:   fmt.Sprintf("Rate limit exceeded for resource: %s", resource),
		UserID:    userID,
		IPAddress: ipAddress,
		Resource:  resource,
		Action:    "request",
		Result:    "rate_limited",
		Details:   details,
	})
}

// LogReplayAttempt logs a potential replay attack attempt
func (sl *SecurityLogger) LogReplayAttempt(userID, ipAddress string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventReplayAttempt,
		Severity:  SeverityHigh,
		Message:   "Potential replay attack detected",
		UserID:    userID,
		IPAddress: ipAddress,
		Action:    "transaction",
		Result:    "blocked",
		Details:   details,
	})
}

// LogSuspiciousActivity logs suspicious activity
func (sl *SecurityLogger) LogSuspiciousActivity(userID, ipAddress, description string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventSuspiciousActivity,
		Severity:  SeverityHigh,
		Message:   description,
		UserID:    userID,
		IPAddress: ipAddress,
		Details:   details,
	})
}

// LogInvalidInput logs invalid input attempts
func (sl *SecurityLogger) LogInvalidInput(userID, ipAddress, field, reason string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventInvalidInput,
		Severity:  SeverityLow,
		Message:   fmt.Sprintf("Invalid input in field '%s': %s", field, reason),
		UserID:    userID,
		IPAddress: ipAddress,
		Action:    "input_validation",
		Result:    "invalid",
		Details: map[string]string{
			"field": field,
		},
	})
}

// LogUnauthorizedAccess logs unauthorized access attempts
func (sl *SecurityLogger) LogUnauthorizedAccess(userID, ipAddress, resource string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:      EventUnauthorizedAccess,
		Severity:  SeverityHigh,
		Message:   fmt.Sprintf("Unauthorized access attempt to: %s", resource),
		UserID:    userID,
		IPAddress: ipAddress,
		Resource:  resource,
		Action:    "access",
		Result:    "unauthorized",
		Details:   details,
	})
}

// LogConfigChange logs configuration changes
func (sl *SecurityLogger) LogConfigChange(userID, configKey, oldValue, newValue string) {
	sl.LogEvent(SecurityEvent{
		Type:     EventConfigChange,
		Severity: SeverityMedium,
		Message:  fmt.Sprintf("Configuration changed: %s", configKey),
		UserID:   userID,
		Action:   "config_change",
		Result:   "success",
		Details: map[string]string{
			"config_key": configKey,
			"old_value":  SanitizeString(oldValue),
			"new_value":  SanitizeString(newValue),
		},
	})
}

// LogSecurityAlert logs a general security alert
func (sl *SecurityLogger) LogSecurityAlert(severity SeverityLevel, alertType SecurityEventType, message string, details map[string]string) {
	sl.LogEvent(SecurityEvent{
		Type:     alertType,
		Severity: severity,
		Message:  message,
		Details:  details,
	})
}

// SecurityAuditLogger provides audit logging capabilities
type SecurityAuditLogger struct {
	logger *SecurityLogger
}

// NewSecurityAuditLogger creates a new security audit logger
func NewSecurityAuditLogger(baseLogger log.Logger) *SecurityAuditLogger {
	return &SecurityAuditLogger{
		logger: NewSecurityLogger(baseLogger),
	}
}

// LogTransaction logs a transaction for audit purposes
func (sal *SecurityAuditLogger) LogTransaction(txHash, sender, receiver string, amount string, success bool) {
	severity := SeverityInfo
	result := "success"
	if !success {
		severity = SeverityMedium
		result = "failure"
	}

	sal.logger.LogEvent(SecurityEvent{
		Type:     EventSecurityAlert,
		Severity: severity,
		Message:  fmt.Sprintf("Transaction %s: %s from %s to %s", result, amount, sender, receiver),
		UserID:   sender,
		Details: map[string]string{
			"tx_hash":  SanitizeHex(txHash),
			"sender":   SanitizeAddress(sender),
			"receiver": SanitizeAddress(receiver),
			"amount":   amount,
			"result":   result,
		},
	})
}

// LogPermissionChange logs permission changes
func (sal *SecurityAuditLogger) LogPermissionChange(granter, grantee, permission string, granted bool) {
	action := "revoked"
	if granted {
		action = "granted"
	}

	sal.logger.LogEvent(SecurityEvent{
		Type:     EventConfigChange,
		Severity: SeverityMedium,
		Message:  fmt.Sprintf("Permission %s: %s to %s", action, permission, grantee),
		UserID:   granter,
		Details: map[string]string{
			"granter":    granter,
			"grantee":    grantee,
			"permission": permission,
			"action":     action,
		},
	})
}

// LogDataAccess logs data access events
func (sal *SecurityAuditLogger) LogDataAccess(userID, resource, dataType, action string) {
	sal.logger.LogEvent(SecurityEvent{
		Type:     EventAccessGranted,
		Severity: SeverityInfo,
		Message:  fmt.Sprintf("Data accessed: %s (%s)", resource, dataType),
		UserID:   userID,
		Resource: resource,
		Action:   action,
		Result:   "success",
		Details: map[string]string{
			"data_type": dataType,
		},
	})
}
