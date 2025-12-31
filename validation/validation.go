package validation

import (
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool              `json:"isValid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// ValidateRequired validates that a string field is not empty
func ValidateRequired(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
			Value:   value,
		}
	}
	return nil
}

// ValidateUUID validates that a string is a valid UUID
func ValidateUUID(value, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
			Value:   value,
		}
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(strings.ToLower(value)) {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a valid UUID",
			Value:   value,
		}
	}
	return nil
}

// ValidateEmail validates that a string is a valid email address
func ValidateEmail(value, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
			Value:   value,
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a valid email address",
			Value:   value,
		}
	}
	return nil
}

// ValidatePhone validates that a string is a valid phone number
func ValidatePhone(value, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
			Value:   value,
		}
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(value) {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a valid phone number",
			Value:   value,
		}
	}
	return nil
}

// ValidateLength validates that a string has a specific length range
func ValidateLength(value, fieldName string, min, max int) *ValidationError {
	length := len(strings.TrimSpace(value))
	if length < min || length > max {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be between " + string(rune(min)) + " and " + string(rune(max)) + " characters",
			Value:   value,
		}
	}
	return nil
}

// ValidateNumeric validates that a string contains only numeric characters
func ValidateNumeric(value, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
			Value:   value,
		}
	}

	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numericRegex.MatchString(value) {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must contain only numeric characters",
			Value:   value,
		}
	}
	return nil
}

// ValidateTripRequest validates common trip request fields
func ValidateTripRequest(tripID, userID string) error {
	if tripID == "" {
		return status.Error(codes.InvalidArgument, "trip ID is required")
	}
	if userID == "" {
		return status.Error(codes.InvalidArgument, "user ID is required")
	}
	return nil
}

// ValidateLocation validates location coordinates
func ValidateLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return status.Error(codes.InvalidArgument, "invalid latitude: must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return status.Error(codes.InvalidArgument, "invalid longitude: must be between -180 and 180")
	}
	return nil
}

// ValidatePrice validates that a price is positive
func ValidatePrice(price float64, fieldName string) *ValidationError {
	if price <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be greater than 0",
			Value:   string(rune(int(price))),
		}
	}
	return nil
}

// ValidateTime validates that a time is not in the past
func ValidateTime(t time.Time, fieldName string) *ValidationError {
	if t.IsZero() {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}

	if t.Before(time.Now()) {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " cannot be in the past",
		}
	}
	return nil
}

// ValidateEnum validates that a value is one of the allowed enum values
func ValidateEnum(value, fieldName string, allowedValues []string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
			Value:   value,
		}
	}

	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}

	return &ValidationError{
		Field:   fieldName,
		Message: fieldName + " must be one of: " + strings.Join(allowedValues, ", "),
		Value:   value,
	}
}

// ValidateMultiple validates multiple fields and returns all errors
func ValidateMultiple(validators ...func() *ValidationError) []ValidationError {
	var errors []ValidationError

	for _, validator := range validators {
		if err := validator(); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// IsValidUUID checks if a string is a valid UUID without returning an error
func IsValidUUID(value string) bool {
	if value == "" {
		return false
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(value))
}

// IsValidEmail checks if a string is a valid email without returning an error
func IsValidEmail(value string) bool {
	if value == "" {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(value)
}

// IsValidPhone checks if a string is a valid phone number without returning an error
func IsValidPhone(value string) bool {
	if value == "" {
		return false
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(value)
}

// SanitizeString removes leading/trailing whitespace and converts to lowercase
func SanitizeString(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// SanitizeEmail sanitizes an email address
func SanitizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// Custom validation functions for specific business logic

// ValidateTripStatus, ValidatePaymentStatus, ValidatePriceModel, and ValidateBiddingStatus
// have been moved to Common/validation/constants.go to use standardized constants.
// These functions are kept here for backward compatibility but delegate to the new implementations.
// Deprecated: Use the functions in constants.go which use standardized constants from Common/constants.go

// ValidateUserType validates user type
func ValidateUserType(userType string) *ValidationError {
	allowedTypes := []string{
		"driver",
		"rider",
		"admin",
	}
	return ValidateEnum(userType, "userType", allowedTypes)
}

// Helper function to convert validation errors to gRPC status
func ValidationErrorsToStatus(errors []ValidationError) error {
	if len(errors) == 0 {
		return nil
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Field+": "+err.Message)
	}

	return status.Error(codes.InvalidArgument, strings.Join(messages, "; "))
}

// Helper function to check if a validation result is valid
func IsValid(result ValidationResult) bool {
	return result.IsValid && len(result.Errors) == 0
}

// Helper function to get first validation error message
func GetFirstError(result ValidationResult) string {
	if len(result.Errors) > 0 {
		return result.Errors[0].Message
	}
	return ""
}

// Helper function to get all validation error messages
func GetAllErrors(result ValidationResult) []string {
	var messages []string
	for _, err := range result.Errors {
		messages = append(messages, err.Message)
	}
	return messages
}
