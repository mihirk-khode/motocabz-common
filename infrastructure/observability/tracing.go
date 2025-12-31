package observability

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TracerProvider wraps OpenTelemetry tracer provider
type TracerProvider struct {
	tp *sdktrace.TracerProvider
}

// InitTracer initializes OpenTelemetry tracer for a service
// This sets up:
// 1. Tracer - Creates spans
// 2. Exporter - Sends telemetry data (OTLP gRPC)
// 3. Propagator - Passes trace context across boundaries
// 4. Resource - Adds metadata (service name, version, env)
// 5. BatchSpanProcessor - Buffers & exports spans efficiently
// Returns a no-op tracer provider if OTEL_EXPORTER_OTLP_ENDPOINT is not configured
func InitTracer(serviceName, serviceVersion string) (*TracerProvider, error) {
	ctx := context.Background()

	// Get endpoint from env or use default
	endpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	// If endpoint is disabled or empty, return no-op tracer provider
	if endpoint == "" || endpoint == "none" || endpoint == "disabled" {
		// Create a no-op tracer provider that still allows spans to be created
		// but doesn't export them anywhere
		res, err := resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceName(serviceName),
				semconv.ServiceVersion(serviceVersion),
				semconv.DeploymentEnvironment(getEnv("ENVIRONMENT", "development")),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource: %w", err)
		}

		// No-op tracer provider (no exporter, spans are created but not exported)
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
			sdktrace.WithSampler(sdktrace.NeverSample()), // Don't sample when disabled
		)

		// Set global tracer provider
		otel.SetTracerProvider(tp)

		// Set global propagator for context propagation (still useful even without export)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		return &TracerProvider{tp: tp}, nil
	}

	// Strip http:// or https:// prefix if present (gRPC endpoint should be host:port only)
	endpoint = normalizeEndpoint(endpoint)

	// 2. EXPORTER - Create OTLP gRPC exporter
	// Note: Connection is non-blocking - BatchSpanProcessor handles retries
	// WithTimeout reduces connection timeout to fail faster and reduce noisy errors
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(2*time.Second), // Shorter timeout to fail faster
		// Removed WithBlock() to allow non-blocking connection
		// The BatchSpanProcessor will retry connections automatically
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// 4. RESOURCE - Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.DeploymentEnvironment(getEnv("ENVIRONMENT", "development")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure sampler
	sampler := sdktrace.AlwaysSample()
	if getEnv("OTEL_SAMPLING_RATE", "1.0") != "1.0" {
		// Can implement probabilistic sampling here if needed
		sampler = sdktrace.AlwaysSample()
	}

	// Wrap exporter to suppress connection errors gracefully
	wrappedExporter := newErrorSuppressingExporter(exporter)

	// 5. BATCHSPANPROCESSOR - Created automatically by WithBatcher()
	// WithBatcher creates a BatchSpanProcessor that:
	// - Buffers spans in memory
	// - Batches them for efficient export
	// - Handles retries and backoff
	// BatchSpanProcessorOptions reduce export timeout to fail faster
	batchProcessor := sdktrace.NewBatchSpanProcessor(wrappedExporter,
		sdktrace.WithBatchTimeout(5*time.Second),  // Export batch every 5s
		sdktrace.WithExportTimeout(2*time.Second), // Fail export after 2s
		sdktrace.WithMaxExportBatchSize(512),      // Batch size
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batchProcessor),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// 3. PROPAGATOR - Set global propagator for context propagation
	// This enables trace context to be passed via HTTP headers and gRPC metadata
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C Trace Context (traceparent header)
		propagation.Baggage{},      // W3C Baggage (baggage header)
	))

	return &TracerProvider{tp: tp}, nil
}

// Shutdown gracefully shuts down the tracer provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return tp.tp.Shutdown(ctx)
}

// GetTracer returns a tracer for the given name
// 1. TRACER - Creates spans (units of work)
func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a new span in the current context
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := GetTracer("motocabz-common")
	return tracer.Start(ctx, name, opts...)
}

// SpanFromContext extracts span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// normalizeEndpoint removes http:// or https:// prefix from endpoint
// OTLP gRPC endpoints should be in format host:port (no scheme)
func normalizeEndpoint(endpoint string) string {
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		return endpoint[7:]
	}
	if len(endpoint) > 8 && endpoint[:8] == "https://" {
		return endpoint[8:]
	}
	return endpoint
}

// errorSuppressingExporter wraps a span exporter to suppress connection-related errors
// This prevents noisy logs when the OTLP endpoint is unreachable
type errorSuppressingExporter struct {
	exporter sdktrace.SpanExporter
}

// newErrorSuppressingExporter creates a new error-suppressing exporter wrapper
func newErrorSuppressingExporter(exporter sdktrace.SpanExporter) sdktrace.SpanExporter {
	return &errorSuppressingExporter{exporter: exporter}
}

// ExportSpans suppresses connection errors but allows other errors through
func (e *errorSuppressingExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	err := e.exporter.ExportSpans(ctx, spans)
	if err == nil {
		return nil
	}

	// Suppress connection-related errors
	errStr := err.Error()
	if isConnectionError(err, errStr) {
		// Silently ignore connection errors - the BatchSpanProcessor will retry
		return nil
	}

	// Allow other errors through
	return err
}

// Shutdown forwards shutdown to the underlying exporter
func (e *errorSuppressingExporter) Shutdown(ctx context.Context) error {
	return e.exporter.Shutdown(ctx)
}

// isConnectionError checks if an error is connection-related
func isConnectionError(err error, errStr string) bool {
	// Check for gRPC connection errors
	if grpcStatus, ok := status.FromError(err); ok {
		code := grpcStatus.Code()
		if code == codes.Unavailable || code == codes.DeadlineExceeded {
			return true
		}
	}

	// Check error message for connection-related keywords
	lowerErr := strings.ToLower(errStr)
	connectionKeywords := []string{
		"connection",
		"connectex",
		"connection refused",
		"no connection could be made",
		"dial tcp",
		"transport: error while dialing",
		"connection error",
	}

	for _, keyword := range connectionKeywords {
		if strings.Contains(lowerErr, keyword) {
			return true
		}
	}

	return false
}
