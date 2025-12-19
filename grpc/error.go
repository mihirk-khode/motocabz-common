package grpc

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", ve.Field, ve.Message)
}

// ToGRPCStatus converts ValidationError to gRPC status
func (ve *ValidationError) ToGRPCStatus() error {
	return status.Error(codes.InvalidArgument, ve.Error())
}

// ValidateRequest validates that request is not nil
func ValidateRequest(req interface{}) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "Request cannot be nil")
	}
	return nil
}

// ValidateID validates that an ID is not empty
func ValidateID(id, fieldName string) error {
	if strings.TrimSpace(id) == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%s cannot be empty", fieldName))
	}
	return nil
}

// ValidateEmail validates email format (basic check)
func ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return status.Error(codes.InvalidArgument, "Email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return status.Error(codes.InvalidArgument, "Invalid email format")
	}
	return nil
}

// ValidatePhone validates phone number format (basic check)
func ValidatePhone(phone string) error {
	if strings.TrimSpace(phone) == "" {
		return status.Error(codes.InvalidArgument, "Phone number cannot be empty")
	}
	if len(phone) < 10 {
		return status.Error(codes.InvalidArgument, "Phone number must be at least 10 digits")
	}
	return nil
}
