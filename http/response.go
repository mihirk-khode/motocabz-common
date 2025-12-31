package http

import (
	"net/http"

	common "github.com/motocabz/common"

	"github.com/gin-gonic/gin"
)

// JSONResponse sends a custom JSON response
func JSONResponse(c *gin.Context, response common.RsBase) {
	statusCode := http.StatusOK
	if response.Error != nil {
		statusCode = response.Error.Code
	}
	c.JSON(statusCode, response)
}
