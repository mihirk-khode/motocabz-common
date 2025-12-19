package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/mihirk-khode/motocabz-common"
)

// SuccessResponse sends a successful response using RsBase
// Deprecated: Use http.HandleError for errors and direct c.JSON for success responses
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := common.RsOK(data, message)
	c.JSON(http.StatusOK, response)
}

// ErrorResponse sends an error response using RsBase
// Deprecated: Use http.HandleError instead
func ErrorResponse(c *gin.Context, code int, message string, errMsg interface{}) {
	response := common.RsErr(code, message, errMsg)
	c.JSON(code, response)
}

// JSONResponse sends a custom JSON response
func JSONResponse(c *gin.Context, response common.RsBase) {
	statusCode := http.StatusOK
	if response.Error != nil {
		statusCode = response.Error.Code
	}
	c.JSON(statusCode, response)
}
