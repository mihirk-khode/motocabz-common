package http

import (
	"fmt"
	"log"

	"github.com/motocabz/common"
	"github.com/motocabz/common/domain"
	"github.com/motocabz/common/infrastructure/observability"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HandleError handles errors and returns appropriate HTTP response
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Get span from context for tracing
	span := trace.SpanFromContext(c.Request.Context())

	var appErr *domain.AppError
	if e, ok := err.(*domain.AppError); ok {
		appErr = e
	} else {
		// Auto-convert unknown errors
		converter := domain.NewErrorConverter()
		appErr = converter.Convert(err)
	}

	// Record error in span
	if span.IsRecording() {
		span.SetStatus(codes.Error, appErr.Message)
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.code", string(appErr.Code)),
			attribute.String("error.message", appErr.Message),
			attribute.Int("http.status", appErr.Status),
		)
	}

	// Extract request ID from context if available
	if requestID := c.GetString("requestId"); requestID != "" {
		appErr.WithDetails("requestId", requestID)
		if span.IsRecording() {
			span.SetAttributes(attribute.String("request.id", requestID))
		}
	}

	// Get trace ID for response - extract directly from span
	traceID := ""
	spanCtx := span.SpanContext()
	if spanCtx.IsValid() {
		traceID = spanCtx.TraceID().String()
	}

	// If still empty, try to get from context (fallback)
	if traceID == "" {
		traceID = observability.GetTraceID(c.Request.Context())
	}

	// Log error with trace ID if available
	if traceID != "" {
		log.Printf("Error [%s] [traceId: %s]: %s - Details: %+v", appErr.Code, traceID, appErr.Message, appErr.Details)
	} else {
		log.Printf("Error [%s]: %s - Details: %+v", appErr.Code, appErr.Message, appErr.Details)
	}

	// Return error response with trace ID using standardized format
	var response common.RsBase
	if appErr.Details != nil {
		// Use details version if details exist
		var errMsg interface{}
		if appErr.Err != nil {
			errMsg = appErr.Err.Error()
		}
		response = common.RsErrDetailsWithTraceID(
			appErr.Status,
			appErr.Message,
			errMsg,
			appErr.Details,
			traceID,
		)
	} else {
		// Use simple error version
		var errMsg interface{}
		if appErr.Err != nil {
			errMsg = appErr.Err.Error()
		}
		response = common.RsErrWithTraceID(
			appErr.Status,
			appErr.Message,
			errMsg,
			traceID,
		)
	}

	c.JSON(appErr.Status, response)
}

// ErrorMiddleware is a Gin middleware that recovers from panics and handles errors
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					HandleError(c, e)
				} else {
					HandleError(c, domain.ErrInternalf("Panic occurred", fmt.Errorf("%v", err)))
				}
				c.Abort()
			}
		}()
		c.Next()

		// Check for errors set by handlers
		if len(c.Errors) > 0 {
			HandleError(c, c.Errors.Last())
			c.Abort()
		}
	}
}
