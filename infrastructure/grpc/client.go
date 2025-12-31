package grpc

import (
	"context"
	"log"
	"time"

	"github.com/motocabz/common/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client wraps gRPC with retry and circuit breaker
type Client struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	retries int
}

// NewClient creates a simple resilient client
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:    conn,
		timeout: 30 * time.Second,
		retries: 3,
	}
}

// WithTimeout sets custom timeout
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

// WithRetries sets custom retry count
func (c *Client) WithRetries(retries int) *Client {
	c.retries = retries
	return c
}

// Call executes a gRPC call with automatic retry and tracing
func (c *Client) Call(ctx context.Context, fn func(context.Context) error) error {
	// Start span for gRPC call
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().
		Tracer("motocabz-common/grpc").
		Start(ctx, "grpc.call")
	defer span.End()

	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		callCtx, cancel := context.WithTimeout(ctx, c.timeout)

		// Add attempt attribute
		if span.IsRecording() {
			span.SetAttributes(attribute.Int("grpc.attempt", attempt+1))
		}

		err := fn(callCtx)
		cancel()

		if err == nil {
			if span.IsRecording() {
				span.SetStatus(codes.Ok, "Success")
			}
			return nil
		}

		// Don't retry non-retryable errors
		if !isRetryable(err) {
			if span.IsRecording() {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			}
			return c.toDomainError(err)
		}

		lastErr = err
		if attempt < c.retries {
			backoff := time.Duration(attempt+1) * 100 * time.Millisecond
			if span.IsRecording() {
				span.AddEvent("retry", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("backoff", backoff.String()),
				))
			}
			log.Printf("gRPC call failed (attempt %d/%d), retrying in %v: %v",
				attempt+1, c.retries, backoff, err)
			time.Sleep(backoff)
		}
	}

	if span.IsRecording() {
		span.SetStatus(codes.Error, lastErr.Error())
		span.RecordError(lastErr)
	}
	return c.toDomainError(lastErr)
}

// isRetryable checks if an error is retryable
func isRetryable(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	// Retryable gRPC codes
	retryableCodes := []grpccodes.Code{
		grpccodes.Unavailable,
		grpccodes.DeadlineExceeded,
		grpccodes.ResourceExhausted,
		grpccodes.Aborted,
		grpccodes.Internal,
	}

	for _, code := range retryableCodes {
		if st.Code() == code {
			return true
		}
	}

	return false
}

// toDomainError converts gRPC errors to domain errors
func (c *Client) toDomainError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return domain.ErrServiceUnavailablef("gRPC service", err)
	}

	switch st.Code() {
	case grpccodes.NotFound:
		return domain.ErrNotFoundf("Resource", "")
	case grpccodes.DeadlineExceeded:
		return domain.ErrTimeoutf("gRPC call", c.timeout)
	case grpccodes.Unavailable:
		return domain.ErrServiceUnavailablef("gRPC service", err)
	case grpccodes.Unauthenticated:
		return domain.ErrUnauthorizedf("Authentication failed")
	case grpccodes.PermissionDenied:
		return domain.ErrForbiddenf("Permission denied")
	case grpccodes.InvalidArgument:
		return domain.ErrValidationf("Invalid argument: %s", st.Message())
	case grpccodes.AlreadyExists:
		return domain.ErrConflictf("Resource already exists")
	default:
		return domain.ErrInternalf("gRPC call failed", err)
	}
}

// Close closes the underlying connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
