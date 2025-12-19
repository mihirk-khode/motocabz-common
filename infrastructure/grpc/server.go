package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

// WithTracing adds OpenTelemetry tracing to gRPC server (unary)
// This automatically creates spans for gRPC unary calls
func WithTracing() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
	))
}

// WithTracingStream adds OpenTelemetry tracing to gRPC server (stream)
// This automatically creates spans for gRPC streaming calls
func WithTracingStream() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
	))
}

// WithClientTracing adds OpenTelemetry tracing to gRPC client
// This automatically creates spans for gRPC calls
func WithClientTracing() grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}

// WithClientTracingStream adds OpenTelemetry tracing to gRPC client (stream)
// This is now the same as WithClientTracing since stats handler covers both unary and stream
func WithClientTracingStream() grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}
