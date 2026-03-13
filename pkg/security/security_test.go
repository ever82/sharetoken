package security

import (
	"testing"
	"time"
)

func TestSanitizeAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "short key",
			input:    "sk-abc",
			expected: "s*****",
		},
		{
			name:     "normal key",
			input:    "sk-abcdefghijklmnopqrstuvwxyz123456789",
			expected: "sk-a******************************6789",
		},
		{
			name:     "exact 8 chars",
			input:    "sk-12345",
			expected: "s*******",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeAPIKey(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeAPIKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizePrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "private key",
			input:    "privatekey123456789",
			expected: "[REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizePrivateKey(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizePrivateKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeMnemonic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "12 word mnemonic",
			input:    "abandon ability able about above absent absorb abstract absurd abuse access accident",
			expected: "[12 words mnemonic - REDACTED]",
		},
		{
			name:     "24 word mnemonic",
			input:    "abandon ability able about above absent absorb abstract absurd abuse access accident acid acquire across act action actor actress actual adapt add address admit",
			expected: "[24 words mnemonic - REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeMnemonic(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeMnemonic() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "short address",
			input:    "abc",
			expected: "***",
		},
		{
			name:     "normal address",
			input:    "sharetoken1xyz789abcdef",
			expected: "sharet...cdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeAddress(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "short hex",
			input:    "0xabc",
			expected: "0x***",
		},
		{
			name:     "64 char hash",
			input:    "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: "0x1234...cdef",
		},
		{
			name:     "without 0x prefix",
			input:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: "0x1234...cdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeHex(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeHex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsSensitiveField(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected bool
	}{
		{
			name:     "api_key",
			field:    "api_key",
			expected: true,
		},
		{
			name:     "password",
			field:    "password",
			expected: true,
		},
		{
			name:     "normal field",
			field:    "name",
			expected: false,
		},
		{
			name:     "token",
			field:    "access_token",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSensitiveField(tt.field)
			if result != tt.expected {
				t.Errorf("IsSensitiveField(%s) = %v, want %v", tt.field, result, tt.expected)
			}
		})
	}
}

func TestNonceTracker(t *testing.T) {
	nt := NewNonceTracker(1 * time.Minute)

	// Test first nonce
	if !nt.CheckAndRecord("addr1", 1) {
		t.Error("First nonce should be valid")
	}

	// Test duplicate nonce
	if nt.CheckAndRecord("addr1", 1) {
		t.Error("Duplicate nonce should be invalid")
	}

	// Test new nonce for same address
	if !nt.CheckAndRecord("addr1", 2) {
		t.Error("New nonce should be valid")
	}

	// Test same nonce for different address
	if !nt.CheckAndRecord("addr2", 1) {
		t.Error("Same nonce for different address should be valid")
	}

	// Test IsUsed
	if !nt.IsUsed("addr1", 1) {
		t.Error("IsUsed should return true for used nonce")
	}
	if nt.IsUsed("addr1", 99) {
		t.Error("IsUsed should return false for unused nonce")
	}
}

func TestSequenceValidator(t *testing.T) {
	sv := NewSequenceValidator()

	// Test first sequence
	if err := sv.ValidateAndUpdate("addr1", 1); err != nil {
		t.Errorf("First sequence should be valid: %v", err)
	}

	// Test old sequence
	if err := sv.ValidateAndUpdate("addr1", 0); err == nil {
		t.Error("Old sequence should be invalid")
	}

	// Test next sequence
	if err := sv.ValidateAndUpdate("addr1", 2); err != nil {
		t.Errorf("Next sequence should be valid: %v", err)
	}

	// Test future sequence
	if err := sv.ValidateAndUpdate("addr1", 5); err != nil {
		t.Errorf("Future sequence should be valid: %v", err)
	}

	// Test expected sequence
	expected := sv.GetExpectedSequence("addr1")
	if expected != 6 {
		t.Errorf("Expected sequence should be 6, got %d", expected)
	}
}

func TestTimestampValidator(t *testing.T) {
	tv := NewTimestampValidator(5 * time.Minute)

	// Test current time
	if err := tv.ValidateTimestamp(time.Now()); err != nil {
		t.Errorf("Current time should be valid: %v", err)
	}

	// Test past time (within window)
	pastTime := time.Now().Add(-1 * time.Minute)
	if err := tv.ValidateTimestamp(pastTime); err != nil {
		t.Errorf("Past time within window should be valid: %v", err)
	}

	// Test past time (outside window)
	oldTime := time.Now().Add(-10 * time.Minute)
	if err := tv.ValidateTimestamp(oldTime); err == nil {
		t.Error("Old time should be invalid")
	}

	// Test future time (within window)
	futureTime := time.Now().Add(1 * time.Minute)
	if err := tv.ValidateTimestamp(futureTime); err != nil {
		t.Errorf("Future time within window should be valid: %v", err)
	}

	// Test future time (outside window)
	farFuture := time.Now().Add(10 * time.Minute)
	if err := tv.ValidateTimestamp(farFuture); err == nil {
		t.Error("Far future time should be invalid")
	}
}

func TestRateLimiter(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 3,
		WindowSize:        1 * time.Minute,
		BurstSize:         1,
	}

	rl := NewRateLimiter(config)

	// Test first 3 requests should be allowed
	for i := 0; i < 3; i++ {
		if !rl.Allow("key1") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Test 4th request should be denied
	if rl.Allow("key1") {
		t.Error("4th request should be denied")
	}

	// Test Check returns false when rate limited
	if rl.Check("key1") {
		t.Error("Check should return false when rate limited")
	}

	// Test different key should be allowed
	if !rl.Allow("key2") {
		t.Error("Different key should be allowed")
	}

	// Test GetRemaining
	remaining := rl.GetRemaining("key1")
	if remaining != 0 {
		t.Errorf("Remaining should be 0, got %d", remaining)
	}

	remaining = rl.GetRemaining("key2")
	if remaining != 2 {
		t.Errorf("Remaining should be 2, got %d", remaining)
	}
}

func TestRateLimiterTier(t *testing.T) {
	tests := []struct {
		tier     RateLimitTier
		expected int
	}{
		{TierStrict, 10},
		{TierStandard, 60},
		{TierRelaxed, 300},
	}

	for _, tt := range tests {
		rl := NewRateLimiterWithTier(tt.tier)
		if rl.config.RequestsPerWindow != tt.expected {
			t.Errorf("Tier %d: expected %d requests per window, got %d",
				tt.tier, tt.expected, rl.config.RequestsPerWindow)
		}
	}
}

func TestIPRateLimiter(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerWindow: 2,
		WindowSize:        1 * time.Minute,
		BurstSize:         1,
	}

	iprl := NewIPRateLimiter(config)

	// Test first 2 requests from IP
	for i := 0; i < 2; i++ {
		if !iprl.Allow("192.168.1.1") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Test 3rd request from same IP
	if iprl.Allow("192.168.1.1") {
		t.Error("3rd request should be denied")
	}

	// Test different IP
	if !iprl.Allow("192.168.1.2") {
		t.Error("Different IP should be allowed")
	}
}

func TestCompositeRateLimiter(t *testing.T) {
	ipConfig := RateLimitConfig{
		RequestsPerWindow: 10,
		WindowSize:        1 * time.Minute,
	}
	userConfig := RateLimitConfig{
		RequestsPerWindow: 5,
		WindowSize:        1 * time.Minute,
	}

	crl := NewCompositeRateLimiter(ipConfig, userConfig)

	// Test allowed request
	allowed, reason := crl.Allow("192.168.1.1", "user1")
	if !allowed {
		t.Errorf("First request should be allowed, got reason: %s", reason)
	}

	// Test rate limit info
	info := crl.GetRateLimitInfo("192.168.1.1", "user1")
	if info.IPRemaining != 9 {
		t.Errorf("IP remaining should be 9, got %d", info.IPRemaining)
	}
	if info.UserRemaining != 4 {
		t.Errorf("User remaining should be 4, got %d", info.UserRemaining)
	}
}

func TestParseIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.1", "192.168.1.1"},
		{"192.168.1.1:8080", "192.168.1.1"},
		{"[::1]:8080", "::1"},
		{"::1", "::1"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		result := ParseIP(tt.input)
		if result != tt.expected {
			t.Errorf("ParseIP(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
	}

	for _, tt := range tests {
		result := IsPrivateIP(tt.ip)
		if result != tt.expected {
			t.Errorf("IsPrivateIP(%s) = %v, want %v", tt.ip, result, tt.expected)
		}
	}
}

func TestValidateBasicSequence(t *testing.T) {
	// Test valid sequence
	if err := ValidateBasicSequence(1); err != nil {
		t.Errorf("Valid sequence should not error: %v", err)
	}

	// Test max uint64 (should fail)
	if err := ValidateBasicSequence(^uint64(0)); err == nil {
		t.Error("Max uint64 should be invalid")
	}
}

func TestSecurityEventCreation(t *testing.T) {
	event := SecurityEvent{
		Type:      EventAuthFailure,
		Severity:  SeverityHigh,
		Message:   "Test message",
		UserID:    "user123",
		IPAddress: "192.168.1.1",
		Resource:  "/api/test",
		Action:    "authenticate",
		Result:    "failure",
		Details: map[string]string{
			"reason": "invalid_credentials",
		},
	}

	if event.Type != EventAuthFailure {
		t.Error("Event type should be AUTH_FAILURE")
	}
	if event.Severity != SeverityHigh {
		t.Error("Event severity should be HIGH")
	}
}
