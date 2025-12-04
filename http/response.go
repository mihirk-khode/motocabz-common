package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/mihirk-khode/motocabz-common"
)

// SuccessResponse sends a successful response using RsBase
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := common.RsOK(data, message)
	c.JSON(http.StatusOK, response)
}

// ErrorResponse sends an error response using RsBase
func ErrorResponse(c *gin.Context, code int, message string, errMsg interface{}) {
	response := common.RsErr(code, message, errMsg)
	c.JSON(code, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	response := common.RsBadRequest("Validation failed", err.Error())
	c.JSON(http.StatusBadRequest, response)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, err error) {
	response := common.RsInternalErr("Internal server error", err.Error())
	c.JSON(http.StatusInternalServerError, response)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, message string) {
	response := common.RsNotFound(message)
	c.JSON(http.StatusNotFound, response)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	response := common.RsUnauthorized(message)
	c.JSON(http.StatusUnauthorized, response)
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context, message string) {
	response := common.RsForbidden(message)
	c.JSON(http.StatusForbidden, response)
}

// BadRequestResponse sends a bad request error response
func BadRequestResponse(c *gin.Context, message string, errMsg interface{}) {
	response := common.RsBadRequest(message, errMsg)
	c.JSON(http.StatusBadRequest, response)
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *gin.Context, message string, errMsg interface{}) {
	response := common.RsConflict(message, errMsg)
	c.JSON(http.StatusConflict, response)
}

// JSONResponse sends a custom JSON response
func JSONResponse(c *gin.Context, response common.RsBase) {
	c.JSON(response.Error.Code, response)
}
