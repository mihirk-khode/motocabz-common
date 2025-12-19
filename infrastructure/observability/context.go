package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// GetTraceID returns the trace ID from context
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()

	// Check if span context is valid
	if spanCtx.IsValid() {
		return spanCtx.TraceID().String()
	}

	// If span context is invalid, return empty string
	// This happens when tracing is disabled or span is not initialized
	return ""
}

// GetSpanID returns the span ID from context
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// IsSampled checks if the current span is sampled
func IsSampled(ctx context.Context) bool {
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().IsSampled()
}
