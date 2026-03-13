package validation

import (
	"fmt"
	"regexp"
)

// ValidationResult represents the result of a validation check
type ValidationResult struct {
	Valid   bool
	Errors  []string
	Field   string
}

// NewValidationResult creates a new validation result
func NewValidationResult(field string) *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Errors: []string{},
		Field:  field,
	}
}

// AddError adds an error to the validation result
func (vr *ValidationResult) AddError(err string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, err)
}

// Error returns the combined error message
func (vr *ValidationResult) Error() error {
	if vr.Valid {
		return nil
	}
	return fmt.Errorf("validation failed for %s: %v", vr.Field, vr.Errors)
}

// Validator interface for custom validation
type Validator interface {
	Validate() error
}

// ValidateAll validates multiple validators and returns all errors
func ValidateAll(validators ...Validator) error {
	var errors []string
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}
	return nil
}

// CompositeValidator chains multiple validation functions
type CompositeValidator struct {
	validations []func() error
}

// NewCompositeValidator creates a new composite validator
func NewCompositeValidator() *CompositeValidator {
	return &CompositeValidator{
		validations: []func() error{},
	}
}

// Add adds a validation function
func (cv *CompositeValidator) Add(fn func() error) *CompositeValidator {
	cv.validations = append(cv.validations, fn)
	return cv
}

// Validate executes all validations and returns the first error
func (cv *CompositeValidator) Validate() error {
	for _, fn := range cv.validations {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateAll executes all validations and returns all errors
func (cv *CompositeValidator) ValidateAll() error {
	var errors []string
	for _, fn := range cv.validations {
		if err := fn(); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("multiple validation errors: %v", errors)
	}
	return nil
}

// Common validation error messages
const (
	ErrMsgEmpty           = "%s cannot be empty"
	ErrMsgInvalidFormat   = "%s has invalid format"
	ErrMsgTooShort        = "%s is too short (minimum %d characters)"
	ErrMsgTooLong         = "%s is too long (maximum %d characters)"
	ErrMsgOutOfRange      = "%s is out of range [%d, %d]"
	ErrMsgInvalidValue    = "%s has invalid value: %v"
	ErrMsgDuplicate       = "%s contains duplicate value: %v"
	ErrMsgNotFound        = "%s not found"
	ErrMsgAlreadyExists   = "%s already exists"
	ErrMsgUnauthorized    = "%s is unauthorized"
	ErrMsgForbidden       = "%s is forbidden"
)

// ValidationError represents a structured validation error
type ValidationError struct {
	Field   string
	Code    string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, code, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "no validation errors"
	}
	if len(ve.Errors) == 1 {
		return ve.Errors[0].Error()
	}
	msgs := make([]string, len(ve.Errors))
	for i, err := range ve.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("%d validation errors: %v", len(ve.Errors), msgs)
}

// Add adds a validation error
func (ve *ValidationErrors) Add(field, code, message string) {
	ve.Errors = append(ve.Errors, NewValidationError(field, code, message))
}

// HasErrors returns true if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// Common regex patterns for validation
var (
	// EmailRegex is a basic email validation pattern
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// UUIDRegex matches UUID format
	UUIDRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	// HashRegex matches common hash formats (64 hex chars)
	HashRegex = regexp.MustCompile(`^(0x)?[0-9a-fA-F]{64}$`)

	// Base64Regex matches base64 encoded strings
	Base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
)

// ValidateEmail validates an email address
func ValidateEmail(email string, fieldName string) error {
	if email == "" {
		return fmt.Errorf(ErrMsgEmpty, fieldName)
	}
	if !EmailRegex.MatchString(email) {
		return fmt.Errorf("%s is not a valid email address", fieldName)
	}
	return nil
}

// ValidateUUID validates a UUID
func ValidateUUID(uuid string, fieldName string) error {
	if uuid == "" {
		return fmt.Errorf(ErrMsgEmpty, fieldName)
	}
	if !UUIDRegex.MatchString(uuid) {
		return fmt.Errorf("%s is not a valid UUID", fieldName)
	}
	return nil
}

// ValidateHash validates a hash string
func ValidateHash(hash string, fieldName string) error {
	if hash == "" {
		return fmt.Errorf(ErrMsgEmpty, fieldName)
	}
	if !HashRegex.MatchString(hash) {
		return fmt.Errorf("%s is not a valid hash", fieldName)
	}
	return nil
}

// ValidateBase64 validates a base64 encoded string
func ValidateBase64(data string, fieldName string) error {
	if data == "" {
		return fmt.Errorf(ErrMsgEmpty, fieldName)
	}
	if !Base64Regex.MatchString(data) {
		return fmt.Errorf("%s is not valid base64", fieldName)
	}
	return nil
}
