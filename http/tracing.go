package http

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware returns Gin middleware for OpenTelemetry tracing
// This automatically creates spans for HTTP requests
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}

// AddSpanAttributes adds attributes to the current span from Gin context
func AddSpanAttributes(c *gin.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(c.Request.Context())
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// SetSpanStatus sets the status of the current span
func SetSpanStatus(c *gin.Context, code codes.Code, description string) {
	span := trace.SpanFromContext(c.Request.Context())
	if span.IsRecording() {
		span.SetStatus(code, description)
	}
}

// GetTraceID returns the trace ID from the current context
func GetTraceID(c *gin.Context) string {
	span := trace.SpanFromContext(c.Request.Context())
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// GetSpanID returns the span ID from the current context
func GetSpanID(c *gin.Context) string {
	span := trace.SpanFromContext(c.Request.Context())
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}
