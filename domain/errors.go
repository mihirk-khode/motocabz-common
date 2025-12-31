package domain

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// General errors
	ErrValidation         ErrorCode = "VALIDATION_ERROR"
	ErrNotFound           ErrorCode = "NOT_FOUND"
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrForbidden          ErrorCode = "FORBIDDEN"
	ErrConflict           ErrorCode = "CONFLICT"
	ErrInternal           ErrorCode = "INTERNAL_ERROR"
	ErrTimeout            ErrorCode = "TIMEOUT"
	ErrRateLimit          ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrDatabaseError      ErrorCode = "DATABASE_ERROR"
	ErrNetworkError       ErrorCode = "NETWORK_ERROR"
	ErrConfigurationError ErrorCode = "CONFIGURATION_ERROR"
)

// AppError is a simple, structured error
type AppError struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
	Err     error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// Simple constructors - easy to use
func ErrValidationf(format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    ErrValidation,
		Message: fmt.Sprintf(format, args...),
		Status:  http.StatusBadRequest,
	}
}

func ErrNotFoundf(resource, id string) *AppError {
	return &AppError{
		Code:    ErrNotFound,
		Message: fmt.Sprintf("%s not found: %s", resource, id),
		Status:  http.StatusNotFound,
		Details: map[string]interface{}{"resource": resource, "id": id},
	}
}

func ErrUnauthorizedf(format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    ErrUnauthorized,
		Message: fmt.Sprintf(format, args...),
		Status:  http.StatusUnauthorized,
	}
}

func ErrForbiddenf(format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    ErrForbidden,
		Message: fmt.Sprintf(format, args...),
		Status:  http.StatusForbidden,
	}
}

func ErrConflictf(format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    ErrConflict,
		Message: fmt.Sprintf(format, args...),
		Status:  http.StatusConflict,
	}
}

func ErrInternalf(msg string, err error) *AppError {
	return &AppError{
		Code:    ErrInternal,
		Message: msg,
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

func ErrTimeoutf(operation string, timeout time.Duration) *AppError {
	return &AppError{
		Code:    ErrTimeout,
		Message: fmt.Sprintf("Operation '%s' timed out after %v", operation, timeout),
		Status:  http.StatusRequestTimeout,
		Details: map[string]interface{}{"operation": operation, "timeout": timeout.String()},
	}
}

func ErrServiceUnavailablef(service string, err error) *AppError {
	return &AppError{
		Code:    ErrServiceUnavailable,
		Message: fmt.Sprintf("Service %s unavailable", service),
		Status:  http.StatusServiceUnavailable,
		Details: map[string]interface{}{"service": service},
		Err:     err,
	}
}

// ErrorConverter converts various error types to AppError
type ErrorConverter struct{}

func NewErrorConverter() *ErrorConverter {
	return &ErrorConverter{}
}

// Convert converts any error to AppError
func (ec *ErrorConverter) Convert(err error) *AppError {
	if err == nil {
		return nil
	}

	// If already an AppError, return as-is
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Check for common error patterns
	errStr := err.Error()

	switch {
	case contains(errStr, "not found"):
		return ErrNotFoundf("resource", "unknown").WithDetails("originalError", errStr)
	case contains(errStr, "already exists"):
		return ErrConflictf("Resource already exists").WithDetails("originalError", errStr)
	case contains(errStr, "permission denied") || contains(errStr, "unauthorized"):
		return ErrUnauthorizedf("Access denied").WithDetails("originalError", errStr)
	case contains(errStr, "timeout") || contains(errStr, "deadline exceeded"):
		return ErrTimeoutf("operation", 0).WithDetails("originalError", errStr)
	case contains(errStr, "connection refused") || contains(errStr, "unavailable"):
		return ErrServiceUnavailablef("external service", err)
	default:
		return ErrInternalf("An unexpected error occurred", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		strings.Contains(s, substr))
}
