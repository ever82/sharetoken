package validation

import (
	"testing"
)

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		minLen  int
		maxLen  int
		wantErr bool
	}{
		{"valid length", "hello", 1, 10, false},
		{"too short", "hi", 5, 10, true},
		{"too long", "hello world", 1, 5, true},
		{"exact min", "a", 1, 10, false},
		{"exact max", "abcde", 1, 5, false},
		{"empty with min 0", "", 0, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringLength(tt.value, "field", tt.minLen, tt.maxLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStringLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid name", "test_name", false},
		{"with hyphen", "test-name", false},
		{"with dot", "test.name", false},
		{"empty", "", true},
		{"with space", "test name", true},
		{"with special", "test@name", true},
		{"too long", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.value, "name")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid description", "This is a test description", false},
		{"empty", "", false},
		{"with HTML", "<script>alert(1)</script>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.value, "description")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com", false},
		{"empty", "", false},
		{"no scheme", "example.com", true},
		{"ftp scheme", "ftp://example.com", true},
		{"javascript", "javascript:alert(1)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.value, "url")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSafeString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"safe string", "Hello World", false},
		{"script tag", "<script>alert(1)</script>", true},
		{"javascript", "javascript:void(0)", true},
		{"onload", "onload=alert(1)", true},
		{"SQL injection", "SELECT * FROM users", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSafeString(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSafeString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAlphanumeric(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"alphanumeric", "abc123", false},
		{"with underscore", "abc_123", false},
		{"with hyphen", "abc-123", false},
		{"with dot", "abc.123", false},
		{"with space", "abc 123", true},
		{"with special", "abc@123", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAlphanumeric(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAlphanumeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid ID", "abc123", false},
		{"with underscore", "abc_123", false},
		{"with hyphen", "abc-123", false},
		{"empty", "", true},
		{"too long", string(make([]byte, 129)), true},
		{"with special", "abc@123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.value, "id")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"  hello  ", "hello"},
		{"hello\x00world", "helloworld"},
		{"hello\x01world", "helloworld"},
		{"hello\nworld", "hello\nworld"},
	}

	for _, tt := range tests {
		result := SanitizeString(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeLogValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"short", "short"},
		{string(make([]byte, 300)), string(make([]byte, 200)) + "..."},
		{"with\nnewline", "with\\nnewline"},
		{"with\rcarriage", "with\\rcarriage"},
	}

	for _, tt := range tests {
		result := SanitizeLogValue(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeLogValue(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid password", "password123", false},
		{"too short", "pass", true},
		{"too long", string(make([]byte, 129)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.value, "password")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHexString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid hex", "abcdef1234567890", false},
		{"with 0x", "0xabcdef1234567890", false},
		{"empty", "", true},
		{"invalid chars", "ghijkl", true},
		{"mixed case", "ABCDEF1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHexString(tt.value, "hex")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHexString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"with subdomain", "test@mail.example.com", false},
		{"empty", "", true},
		{"no @", "testexample.com", true},
		{"no domain", "test@", true},
		{"no local", "@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.value, "email")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", false},
		{"empty", "", true},
		{"invalid format", "not-a-uuid", true},
		{"too short", "550e8400-e29b-41d4-a716", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.value, "uuid")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompositeValidator(t *testing.T) {
	cv := NewCompositeValidator()
	cv.Add(func() error {
		return nil
	}).Add(func() error {
		return nil
	})

	err := cv.Validate()
	if err != nil {
		t.Errorf("CompositeValidator should not error: %v", err)
	}

	cv2 := NewCompositeValidator()
	cv2.Add(func() error {
		return nil
	}).Add(func() error {
		return ValidateNonEmpty("", "field")
	})

	err = cv2.Validate()
	if err == nil {
		t.Error("CompositeValidator should error when one validation fails")
	}
}

func TestValidationErrors(t *testing.T) {
	ve := &ValidationErrors{}
	ve.Add("field1", "ERROR", "error message 1")
	ve.Add("field2", "ERROR", "error message 2")

	if !ve.HasErrors() {
		t.Error("HasErrors should return true")
	}

	err := ve.Error()
	if err == "" {
		t.Error("Error should not be empty")
	}
}
