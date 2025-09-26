package grpc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

// GRPCErrorHandler provides centralized error handling for gRPC services
type GRPCErrorHandler struct {
	serviceName string
}

// NewGRPCErrorHandler creates a new gRPC error handler
func NewGRPCErrorHandler(serviceName string) *GRPCErrorHandler {
	return &GRPCErrorHandler{
		serviceName: serviceName,
	}
}

// HandleError converts internal errors to appropriate gRPC status codes
func (eh *GRPCErrorHandler) HandleError(err error) error {
	if err == nil {
		return nil
	}

	log.Printf("%s gRPC: Error occurred: %v", eh.serviceName, err)

	// Check if it's already a gRPC status error
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Convert common errors to gRPC status codes
	switch {
	case strings.Contains(err.Error(), "not found"):
		return status.Error(codes.NotFound, err.Error())
	case strings.Contains(err.Error(), "already exists"):
		return status.Error(codes.AlreadyExists, err.Error())
	case strings.Contains(err.Error(), "permission denied"):
		return status.Error(codes.PermissionDenied, err.Error())
	case strings.Contains(err.Error(), "invalid argument"):
		return status.Error(codes.InvalidArgument, err.Error())
	case strings.Contains(err.Error(), "timeout"):
		return status.Error(codes.DeadlineExceeded, err.Error())
	case strings.Contains(err.Error(), "unavailable"):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, fmt.Sprintf("Internal server error: %v", err))
	}
}

// ValidateRequest performs common request validations
func (eh *GRPCErrorHandler) ValidateRequest(req interface{}) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "Request cannot be nil")
	}
	return nil
}

// ValidateID validates that an ID is not empty
func (eh *GRPCErrorHandler) ValidateID(id, fieldName string) error {
	if strings.TrimSpace(id) == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%s cannot be empty", fieldName))
	}
	return nil
}

// ValidateLocation validates latitude and longitude coordinates
func (eh *GRPCErrorHandler) ValidateLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return status.Error(codes.InvalidArgument, "Latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return status.Error(codes.InvalidArgument, "Longitude must be between -180 and 180")
	}
	return nil
}

// ValidateEmail validates email format
func (eh *GRPCErrorHandler) ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return status.Error(codes.InvalidArgument, "Email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return status.Error(codes.InvalidArgument, "Invalid email format")
	}
	return nil
}

// ValidatePhone validates phone number format
func (eh *GRPCErrorHandler) ValidatePhone(phone string) error {
	if strings.TrimSpace(phone) == "" {
		return status.Error(codes.InvalidArgument, "Phone number cannot be empty")
	}
	if len(phone) < 10 {
		return status.Error(codes.InvalidArgument, "Phone number must be at least 10 digits")
	}
	return nil
}

// LogRequest logs incoming gRPC requests
func (eh *GRPCErrorHandler) LogRequest(method string, req interface{}) {
	log.Printf("%s gRPC: %s called with request: %+v", eh.serviceName, method, req)
}

// LogResponse logs outgoing gRPC responses
func (eh *GRPCErrorHandler) LogResponse(method string, resp interface{}, err error) {
	if err != nil {
		log.Printf("%s gRPC: %s failed with error: %v", eh.serviceName, method, err)
	} else {
		log.Printf("%s gRPC: %s completed successfully", eh.serviceName, method)
	}
}

// ContextWithTimeout creates a context with timeout
func (eh *GRPCErrorHandler) ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// RetryOperation retries an operation with exponential backoff
func (eh *GRPCErrorHandler) RetryOperation(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if i < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(i)) // Exponential backoff
			log.Printf("%s gRPC: Operation failed (attempt %d/%d), retrying in %v: %v",
				eh.serviceName, i+1, maxRetries, delay, err)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}

// CircuitBreaker provides circuit breaker functionality
type CircuitBreaker struct {
	failureCount    int
	successCount    int
	lastFailureTime time.Time
	state           CircuitState
	threshold       int
	timeout         time.Duration
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		state:     CircuitClosed,
	}
}

// Execute executes an operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(operation func() error) error {
	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = CircuitHalfOpen
		} else {
			return status.Error(codes.Unavailable, "Circuit breaker is open")
		}
	}

	err := operation()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if cb.failureCount >= cb.threshold {
			cb.state = CircuitOpen
		}

		return err
	}

	cb.successCount++
	cb.failureCount = 0

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}

	return nil
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// MetricsCollector collects gRPC service metrics
type MetricsCollector struct {
	serviceName     string
	requestCount    int64
	errorCount      int64
	successCount    int64
	avgResponseTime time.Duration
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(serviceName string) *MetricsCollector {
	return &MetricsCollector{
		serviceName: serviceName,
	}
}

// RecordRequest records a request metric
func (mc *MetricsCollector) RecordRequest(duration time.Duration, err error) {
	mc.requestCount++

	if err != nil {
		mc.errorCount++
	} else {
		mc.successCount++
	}

	// Simple moving average for response time
	if mc.avgResponseTime == 0 {
		mc.avgResponseTime = duration
	} else {
		mc.avgResponseTime = (mc.avgResponseTime + duration) / 2
	}
}

// GetMetrics returns current metrics
func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
	successRate := float64(0)
	if mc.requestCount > 0 {
		successRate = float64(mc.successCount) / float64(mc.requestCount) * 100
	}

	return map[string]interface{}{
		"service_name":      mc.serviceName,
		"request_count":     mc.requestCount,
		"error_count":       mc.errorCount,
		"success_count":     mc.successCount,
		"success_rate":      successRate,
		"avg_response_time": mc.avgResponseTime.String(),
	}
}
