package error

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents a custom error code
type ErrorCode int

const (
	// General errors
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeValidation
	ErrorCodeNotFound
	ErrorCodeUnauthorized
	ErrorCodeForbidden
	ErrorCodeConflict
	ErrorCodeInternal
	ErrorCodeTimeout
	ErrorCodeRateLimit

	// Service-specific errors
	ErrorCodeServiceUnavailable
	ErrorCodeDatabaseError
	ErrorCodeNetworkError
	ErrorCodeConfigurationError

	// Trip-specific errors
	ErrorCodeTripNotFound
	ErrorCodeTripAlreadyExists
	ErrorCodeInvalidTripStatus
	ErrorCodeTripCancelled
	ErrorCodeTripExpired

	// Driver-specific errors
	ErrorCodeDriverNotFound
	ErrorCodeDriverOffline
	ErrorCodeDriverBusy
	ErrorCodeInvalidDriverStatus

	// Rider-specific errors
	ErrorCodeRiderNotFound
	ErrorCodeRiderInactive
	ErrorCodeInvalidRiderStatus

	// Bidding-specific errors
	ErrorCodeBiddingSessionNotFound
	ErrorCodeBiddingSessionExpired
	ErrorCodeInvalidBidAmount
	ErrorCodeBiddingNotAllowed

	// Location-specific errors
	ErrorCodeInvalidLocation
	ErrorCodeLocationNotFound
	ErrorCodeLocationOutOfRange

	// Payment-specific errors
	ErrorCodePaymentFailed
	ErrorCodePaymentNotFound
	ErrorCodeInvalidPaymentMethod
	ErrorCodeInsufficientFunds
)

// ErrorCodeNames maps error codes to their string representations
var ErrorCodeNames = map[ErrorCode]string{
	ErrorCodeUnknown:                "UNKNOWN_ERROR",
	ErrorCodeValidation:             "VALIDATION_ERROR",
	ErrorCodeNotFound:               "NOT_FOUND",
	ErrorCodeUnauthorized:           "UNAUTHORIZED",
	ErrorCodeForbidden:              "FORBIDDEN",
	ErrorCodeConflict:               "CONFLICT",
	ErrorCodeInternal:               "INTERNAL_ERROR",
	ErrorCodeTimeout:                "TIMEOUT",
	ErrorCodeRateLimit:              "RATE_LIMIT_EXCEEDED",
	ErrorCodeServiceUnavailable:     "SERVICE_UNAVAILABLE",
	ErrorCodeDatabaseError:          "DATABASE_ERROR",
	ErrorCodeNetworkError:           "NETWORK_ERROR",
	ErrorCodeConfigurationError:     "CONFIGURATION_ERROR",
	ErrorCodeTripNotFound:           "TRIP_NOT_FOUND",
	ErrorCodeTripAlreadyExists:      "TRIP_ALREADY_EXISTS",
	ErrorCodeInvalidTripStatus:      "INVALID_TRIP_STATUS",
	ErrorCodeTripCancelled:          "TRIP_CANCELLED",
	ErrorCodeTripExpired:            "TRIP_EXPIRED",
	ErrorCodeDriverNotFound:         "DRIVER_NOT_FOUND",
	ErrorCodeDriverOffline:          "DRIVER_OFFLINE",
	ErrorCodeDriverBusy:             "DRIVER_BUSY",
	ErrorCodeInvalidDriverStatus:    "INVALID_DRIVER_STATUS",
	ErrorCodeRiderNotFound:          "RIDER_NOT_FOUND",
	ErrorCodeRiderInactive:          "RIDER_INACTIVE",
	ErrorCodeInvalidRiderStatus:     "INVALID_RIDER_STATUS",
	ErrorCodeBiddingSessionNotFound: "BIDDING_SESSION_NOT_FOUND",
	ErrorCodeBiddingSessionExpired:  "BIDDING_SESSION_EXPIRED",
	ErrorCodeInvalidBidAmount:       "INVALID_BID_AMOUNT",
	ErrorCodeBiddingNotAllowed:      "BIDDING_NOT_ALLOWED",
	ErrorCodeInvalidLocation:        "INVALID_LOCATION",
	ErrorCodeLocationNotFound:       "LOCATION_NOT_FOUND",
	ErrorCodeLocationOutOfRange:     "LOCATION_OUT_OF_RANGE",
	ErrorCodePaymentFailed:          "PAYMENT_FAILED",
	ErrorCodePaymentNotFound:        "PAYMENT_NOT_FOUND",
	ErrorCodeInvalidPaymentMethod:   "INVALID_PAYMENT_METHOD",
	ErrorCodeInsufficientFunds:      "INSUFFICIENT_FUNDS",
}

// CustomError represents a custom error with additional context
type CustomError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Stack     string                 `json:"stack,omitempty"`
	Cause     error                  `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *CustomError) Unwrap() error {
	return e.Cause
}

// NewCustomError creates a new custom error
func NewCustomError(code ErrorCode, message string) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// NewCustomErrorWithCause creates a new custom error with a cause
func NewCustomErrorWithCause(code ErrorCode, message string, cause error) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
		Cause:     cause,
	}
}

// NewCustomErrorWithDetails creates a new custom error with details
func NewCustomErrorWithDetails(code ErrorCode, message string, details map[string]interface{}) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// WrapToCustomError wraps an existing error into a custom error
func WrapToCustomError(domain string, err error) *CustomError {
	if err == nil {
		return nil
	}

	// Check if it's already a custom error
	if customErr, ok := err.(*CustomError); ok {
		return customErr
	}

	// Determine error code based on error type
	code := ErrorCodeInternal
	message := fmt.Sprintf("%s: %v", domain, err)

	// Map common error patterns to specific codes
	errStr := strings.ToLower(err.Error())
	switch {
	case strings.Contains(errStr, "not found"):
		code = ErrorCodeNotFound
	case strings.Contains(errStr, "unauthorized"):
		code = ErrorCodeUnauthorized
	case strings.Contains(errStr, "forbidden"):
		code = ErrorCodeForbidden
	case strings.Contains(errStr, "conflict"):
		code = ErrorCodeConflict
	case strings.Contains(errStr, "timeout"):
		code = ErrorCodeTimeout
	case strings.Contains(errStr, "validation"):
		code = ErrorCodeValidation
	case strings.Contains(errStr, "database"):
		code = ErrorCodeDatabaseError
	case strings.Contains(errStr, "network"):
		code = ErrorCodeNetworkError
	}

	return &CustomError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
		Cause:     err,
	}
}

// getStackTrace returns the current stack trace
func getStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// ErrorCodeToHTTPStatus maps error codes to HTTP status codes
func ErrorCodeToHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrorCodeValidation:
		return http.StatusBadRequest
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrorCodeForbidden:
		return http.StatusForbidden
	case ErrorCodeConflict:
		return http.StatusConflict
	case ErrorCodeTimeout:
		return http.StatusRequestTimeout
	case ErrorCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrorCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrorCodeInternal, ErrorCodeDatabaseError, ErrorCodeNetworkError, ErrorCodeConfigurationError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ErrorCodeToGRPCStatus maps error codes to gRPC status codes
func ErrorCodeToGRPCStatus(code ErrorCode) codes.Code {
	switch code {
	case ErrorCodeValidation:
		return codes.InvalidArgument
	case ErrorCodeNotFound:
		return codes.NotFound
	case ErrorCodeUnauthorized:
		return codes.Unauthenticated
	case ErrorCodeForbidden:
		return codes.PermissionDenied
	case ErrorCodeConflict:
		return codes.AlreadyExists
	case ErrorCodeTimeout:
		return codes.DeadlineExceeded
	case ErrorCodeRateLimit:
		return codes.ResourceExhausted
	case ErrorCodeServiceUnavailable:
		return codes.Unavailable
	case ErrorCodeInternal, ErrorCodeDatabaseError, ErrorCodeNetworkError, ErrorCodeConfigurationError:
		return codes.Internal
	default:
		return codes.Internal
	}
}

// ConvertToGRPCError converts a custom error to a gRPC error
func ConvertToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	if customErr, ok := err.(*CustomError); ok {
		grpcCode := ErrorCodeToGRPCStatus(customErr.Code)
		return status.Error(grpcCode, customErr.Message)
	}

	return status.Error(codes.Internal, err.Error())
}

// ConvertToHTTPError converts a custom error to an HTTP error response
func ConvertToHTTPError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	if customErr, ok := err.(*CustomError); ok {
		httpStatus := ErrorCodeToHTTPStatus(customErr.Code)
		return httpStatus, customErr.Message
	}

	return http.StatusInternalServerError, err.Error()
}

// IsCustomError checks if an error is a custom error
func IsCustomError(err error) bool {
	_, ok := err.(*CustomError)
	return ok
}

// GetErrorCode returns the error code from a custom error
func GetErrorCode(err error) ErrorCode {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code
	}
	return ErrorCodeUnknown
}

// GetErrorMessage returns the error message
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// GetErrorDetails returns the error details from a custom error
func GetErrorDetails(err error) map[string]interface{} {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Details
	}
	return nil
}

// Common error constructors

// NewValidationError creates a validation error
func NewValidationError(message string) *CustomError {
	return NewCustomError(ErrorCodeValidation, message)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *CustomError {
	return NewCustomError(ErrorCodeNotFound, resource+" not found")
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *CustomError {
	return NewCustomError(ErrorCodeUnauthorized, message)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *CustomError {
	return NewCustomError(ErrorCodeForbidden, message)
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *CustomError {
	return NewCustomError(ErrorCodeConflict, message)
}

// NewInternalError creates an internal error
func NewInternalError(message string) *CustomError {
	return NewCustomError(ErrorCodeInternal, message)
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(message string) *CustomError {
	return NewCustomError(ErrorCodeTimeout, message)
}

// NewServiceUnavailableError creates a service unavailable error
func NewServiceUnavailableError(message string) *CustomError {
	return NewCustomError(ErrorCodeServiceUnavailable, message)
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string) *CustomError {
	return NewCustomError(ErrorCodeDatabaseError, message)
}

// NewNetworkError creates a network error
func NewNetworkError(message string) *CustomError {
	return NewCustomError(ErrorCodeNetworkError, message)
}

// Service-specific error constructors

// NewTripNotFoundError creates a trip not found error
func NewTripNotFoundError(tripID string) *CustomError {
	return NewCustomError(ErrorCodeTripNotFound, "Trip not found: "+tripID)
}

// NewDriverNotFoundError creates a driver not found error
func NewDriverNotFoundError(driverID string) *CustomError {
	return NewCustomError(ErrorCodeDriverNotFound, "Driver not found: "+driverID)
}

// NewRiderNotFoundError creates a rider not found error
func NewRiderNotFoundError(riderID string) *CustomError {
	return NewCustomError(ErrorCodeRiderNotFound, "Rider not found: "+riderID)
}

// NewBiddingSessionNotFoundError creates a bidding session not found error
func NewBiddingSessionNotFoundError(sessionID string) *CustomError {
	return NewCustomError(ErrorCodeBiddingSessionNotFound, "Bidding session not found: "+sessionID)
}

// NewInvalidLocationError creates an invalid location error
func NewInvalidLocationError(message string) *CustomError {
	return NewCustomError(ErrorCodeInvalidLocation, message)
}

// NewPaymentFailedError creates a payment failed error
func NewPaymentFailedError(message string) *CustomError {
	return NewCustomError(ErrorCodePaymentFailed, message)
}

// Error aggregation and handling

// ErrorList represents a list of errors
type ErrorList struct {
	Errors []*CustomError `json:"errors"`
}

// Add adds an error to the list
func (el *ErrorList) Add(err *CustomError) {
	if err != nil {
		el.Errors = append(el.Errors, err)
	}
}

// HasErrors returns true if there are any errors
func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

// GetFirstError returns the first error in the list
func (el *ErrorList) GetFirstError() *CustomError {
	if len(el.Errors) > 0 {
		return el.Errors[0]
	}
	return nil
}

// GetAllMessages returns all error messages
func (el *ErrorList) GetAllMessages() []string {
	var messages []string
	for _, err := range el.Errors {
		messages = append(messages, err.Message)
	}
	return messages
}

// Error implements the error interface
func (el *ErrorList) Error() string {
	if len(el.Errors) == 0 {
		return ""
	}
	if len(el.Errors) == 1 {
		return el.Errors[0].Error()
	}
	return fmt.Sprintf("Multiple errors: %s", strings.Join(el.GetAllMessages(), "; "))
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		Errors: make([]*CustomError, 0),
	}
}

// Error recovery and handling

// RecoverFromPanic recovers from a panic and returns a custom error
func RecoverFromPanic(domain string) *CustomError {
	if r := recover(); r != nil {
		message := fmt.Sprintf("%s: panic recovered: %v", domain, r)
		return NewCustomError(ErrorCodeInternal, message)
	}
	return nil
}

// SafeExecute executes a function safely and returns any panic as an error
func SafeExecute(domain string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewCustomError(ErrorCodeInternal, fmt.Sprintf("%s: panic recovered: %v", domain, r))
		}
	}()

	return fn()
}

// Error context and tracing

// ErrorContext represents additional context for an error
type ErrorContext struct {
	RequestID   string                 `json:"requestId,omitempty"`
	UserID      string                 `json:"userId,omitempty"`
	ServiceName string                 `json:"serviceName,omitempty"`
	Operation   string                 `json:"operation,omitempty"`
	Additional  map[string]interface{} `json:"additional,omitempty"`
}

// AddContext adds context to a custom error
func (e *CustomError) AddContext(ctx ErrorContext) {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}

	if ctx.RequestID != "" {
		e.Details["requestId"] = ctx.RequestID
	}
	if ctx.UserID != "" {
		e.Details["userId"] = ctx.UserID
	}
	if ctx.ServiceName != "" {
		e.Details["serviceName"] = ctx.ServiceName
	}
	if ctx.Operation != "" {
		e.Details["operation"] = ctx.Operation
	}

	for k, v := range ctx.Additional {
		e.Details[k] = v
	}
}
