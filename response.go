package common

import (
	"net/http"
	"time"
)

// RsBase represents the standard API response structure
type RsBase struct {
	ApiVersion string      `json:"apiVersion,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      *ErrorInfo  `json:"error,omitempty"`
	Meta       *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo represents error information in API responses
type ErrorInfo struct {
	Code     int         `json:"code"`
	CodeText string      `json:"codeText"`
	Message  string      `json:"message"`
	ErrorMsg interface{} `json:"errorMsg,omitempty"`
	TraceID  string      `json:"traceId,omitempty"`
}

// MetaInfo represents metadata for API responses
type MetaInfo struct {
	Timestamp   time.Time   `json:"timestamp"`
	RequestID   string      `json:"requestId,omitempty"`
	TraceID     string      `json:"traceId,omitempty"`
	Version     string      `json:"version,omitempty"`
	Environment string      `json:"environment,omitempty"`
	Pagination  *Pagination `json:"pagination,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// RsSuccess represents a successful API response
type RsSuccess struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func RsOK(data interface{}, msg string) RsBase {
	return RsBase{
		ApiVersion: "v1",
		Data:       data,
		Meta: &MetaInfo{
			Timestamp: time.Now(),
		},
	}
}

func RsOKMeta(data interface{}, msg string, meta *MetaInfo) RsBase {
	if meta == nil {
		meta = &MetaInfo{
			Timestamp: time.Now(),
		}
	}
	if meta.Timestamp.IsZero() {
		meta.Timestamp = time.Now()
	}

	return RsBase{
		ApiVersion: "v1",
		Data:       data,
		Meta:       meta,
	}
}

func RsErr(code int, msg string, errMsg interface{}) RsBase {
	return RsErrWithTraceID(code, msg, errMsg, "")
}

// RsErrWithTraceID creates an error response with trace ID
func RsErrWithTraceID(code int, msg string, errMsg interface{}, traceID string) RsBase {
	meta := &MetaInfo{
		Timestamp: time.Now(),
	}
	if traceID != "" {
		meta.TraceID = traceID
	}
	return RsBase{
		ApiVersion: "v1",
		Error: &ErrorInfo{
			Code:     code,
			CodeText: http.StatusText(code),
			Message:  msg,
			ErrorMsg: errMsg,
		},
		Meta: meta,
	}
}

func RsErrDetails(code int, msg string, errMsg interface{}, details interface{}) RsBase {
	return RsErrDetailsWithTraceID(code, msg, errMsg, details, "")
}

// RsErrDetailsWithTraceID creates an error response with details and trace ID
func RsErrDetailsWithTraceID(code int, msg string, errMsg interface{}, details interface{}, traceID string) RsBase {
	meta := &MetaInfo{
		Timestamp: time.Now(),
	}
	if traceID != "" {
		meta.TraceID = traceID
	}
	return RsBase{
		ApiVersion: "v1",
		Error: &ErrorInfo{
			Code:     code,
			CodeText: http.StatusText(code),
			Message:  msg,
			ErrorMsg: errMsg,
		},
		Meta: meta,
	}
}

func RsValidationErr(validationErrors []ValidationError) RsBase {
	return RsValidationErrWithTraceID(validationErrors, "")
}

// RsValidationErrWithTraceID creates a validation error response with trace ID
func RsValidationErrWithTraceID(validationErrors []ValidationError, traceID string) RsBase {
	meta := &MetaInfo{
		Timestamp: time.Now(),
	}
	if traceID != "" {
		meta.TraceID = traceID
	}
	return RsBase{
		ApiVersion: "v1",
		Error: &ErrorInfo{
			Code:     http.StatusBadRequest,
			CodeText: http.StatusText(http.StatusBadRequest),
			Message:  "Validation failed",
			ErrorMsg: validationErrors,
		},
		Meta: meta,
	}
}

func RsPaginated(data interface{}, page, limit int, total int64) RsBase {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return RsBase{
		ApiVersion: "v1",
		Data:       data,
		Meta: &MetaInfo{
			Timestamp: time.Now(),
			Pagination: &Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
				HasNext:    page < totalPages,
				HasPrev:    page > 1,
			},
		},
	}
}

func RsNotFound(resource string) RsBase {
	return RsErr(
		http.StatusNotFound,
		resource+" not found",
		nil,
	)
}

func RsUnauthorized(msg string) RsBase {
	if msg == "" {
		msg = "Unauthorized access"
	}
	return RsErr(
		http.StatusUnauthorized,
		msg,
		nil,
	)
}

func RsForbidden(msg string) RsBase {
	if msg == "" {
		msg = "Access forbidden"
	}
	return RsErr(
		http.StatusForbidden,
		msg,
		nil,
	)
}

func RsInternalErr(msg string, errMsg interface{}) RsBase {
	if msg == "" {
		msg = "An internal server error occurred"
	}
	return RsErr(
		http.StatusInternalServerError,
		msg,
		errMsg,
	)
}

func RsBadRequest(msg string, errMsg interface{}) RsBase {
	if msg == "" {
		msg = "Bad request"
	}
	return RsErr(
		http.StatusBadRequest,
		msg,
		errMsg,
	)
}

func RsConflict(msg string, errMsg interface{}) RsBase {
	if msg == "" {
		msg = "Resource conflict"
	}
	return RsErr(
		http.StatusConflict,
		msg,
		errMsg,
	)
}
