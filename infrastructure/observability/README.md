# OpenTelemetry Observability Package

This package provides OpenTelemetry integration for distributed tracing and metrics collection across all Motocabz microservices.

## Components

### 1. **Tracer** - Creates spans (units of work)
- Initialized via `InitTracer()`
- Creates spans for operations
- Manages span lifecycle

### 2. **Exporter** - Sends telemetry data
- OTLP gRPC exporter
- Sends to OpenTelemetry Collector
- Configurable endpoint

### 3. **Propagator** - Passes trace context
- W3C Trace Context propagation
- W3C Baggage propagation
- Automatic context passing via HTTP headers and gRPC metadata

### 4. **Resource** - Adds metadata
- Service name
- Service version
- Environment (development/production)

### 5. **BatchSpanProcessor** - Buffers & exports spans
- Automatic batching for efficiency
- Handles retries and backoff
- Minimizes performance overhead

## Usage

### Initialization in Service

```go
package main

import (
    "context"
    "log"
    "time"
    
    "motocabz/common/infrastructure/observability"
)

func main() {
    // Initialize tracing
    tp, err := observability.InitTracer("trip-service", "1.0.0")
    if err != nil {
        log.Printf("⚠️  Failed to initialize tracer: %v", err)
    } else {
        log.Println("✅ OpenTelemetry tracer initialized")
        defer func() {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            tp.Shutdown(ctx)
        }()
    }

    // Initialize metrics
    mp, err := observability.InitMeter("trip-service", "1.0.0")
    if err != nil {
        log.Printf("⚠️  Failed to initialize meter: %v", err)
    } else {
        log.Println("✅ OpenTelemetry metrics initialized")
        defer func() {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            mp.Shutdown(ctx)
        }()
    }

    // Your service code...
}
```

### Environment Variables

```bash
# OpenTelemetry Configuration
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317  # Default
OTEL_SAMPLING_RATE=1.0  # 1.0 = 100%, 0.1 = 10%
ENVIRONMENT=production
```

### Manual Span Creation

```go
import "github.com/mihirk-khode/motocabz-common/infrastructure/observability"

func myFunction(ctx context.Context) error {
    // Start a span
    ctx, span := observability.StartSpan(ctx, "my.operation")
    defer span.End()

    // Your code here
    // Span automatically records errors and timing

    return nil
}
```

### Get Trace ID

```go
import "github.com/mihirk-khode/motocabz-common/infrastructure/observability"

traceID := observability.GetTraceID(ctx)
spanID := observability.GetSpanID(ctx)
```

## Integration Points

This package is automatically integrated into:

1. **HTTP Handlers** - Via `http.TracingMiddleware()`
2. **gRPC Clients** - Via `infrastructure/grpc.Client.Call()`
3. **gRPC Servers** - Via `infrastructure/grpc.WithTracing()`
4. **Database Transactions** - Via `infrastructure/persistence.WithEntTransaction()`
5. **Error Handling** - Via `http.HandleError()`

## Architecture

```
Service Code
    │
    ├─→ HTTP Request → http.TracingMiddleware() → Creates span
    │
    ├─→ gRPC Call → grpc.Client.Call() → Creates span
    │
    ├─→ Database → persistence.WithEntTransaction() → Creates span
    │
    └─→ Error → http.HandleError() → Records error in span
            │
            └─→ BatchSpanProcessor → Exporter → OpenTelemetry Collector
```

## Benefits

- **Automatic Instrumentation** - No code changes needed in most cases
- **Distributed Tracing** - Track requests across service boundaries
- **Error Tracking** - Automatic error recording in spans
- **Performance Monitoring** - Automatic timing and latency tracking
- **Context Propagation** - Automatic trace context passing

