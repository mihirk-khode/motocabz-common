package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// BindAndValidate binds request data and validates it
// Supports both JSON binding and query/path parameter binding
func BindAndValidate(c *gin.Context, rq interface{}) error {
	// Handle map[string]string for path/query params
	if m, ok := rq.(*map[string]string); ok {
		params := make(map[string]string)
		for _, p := range c.Params {
			if p.Value == "" {
				return errors.New("missing required path parameter: " + p.Key)
			}
			params[p.Key] = p.Value
		}
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
		*m = params
		return nil
	}

	// Try ShouldBindJSON first, then fall back to ShouldBind
	if err := c.ShouldBindJSON(rq); err != nil {
		// If JSON binding fails, try ShouldBind (for form data, query params, etc.)
		if err := c.ShouldBind(rq); err != nil {
			return err
		}
	}

	// Validate struct
	if err := validate.Struct(rq); err != nil {
		return err
	}

	return nil
}

// BindJSON binds JSON request data only
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}

	if err := validate.Struct(obj); err != nil {
		return err
	}

	return nil
}
