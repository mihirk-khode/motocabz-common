# Base Response Structures

This package provides standardized response structures for all microservices in the MotoCabz platform. It ensures consistency across all APIs and provides helper functions for common response patterns.

## Features

- **Standardized Response Format**: Consistent structure across all services
- **Error Handling**: Comprehensive error response structures
- **Pagination Support**: Built-in pagination metadata
- **Validation Errors**: Structured validation error responses
- **Metadata Support**: Request tracking, versioning, and environment info
- **Helper Functions**: Pre-built functions for common response patterns

## Response Structure

### RsBase

The main response structure that all APIs should use:

```go
type RsBase struct {
    ApiVersion string      `json:"apiVersion,omitempty"`
    Data       interface{} `json:"data,omitempty"`
    Error      *ErrorInfo  `json:"error,omitempty"`
    Meta       *MetaInfo   `json:"meta,omitempty"`
}
```

### ErrorInfo

Error information structure:

```go
type ErrorInfo struct {
    Code     int         `json:"code"`
    CodeText string      `json:"codeText"`
    Message  string      `json:"message"`
    ErrorMsg interface{} `json:"errorMsg,omitempty"`
    Details  interface{} `json:"details,omitempty"`
}
```

### MetaInfo

Metadata structure for additional information:

```go
type MetaInfo struct {
    Timestamp   time.Time   `json:"timestamp"`
    RequestID   string      `json:"requestId,omitempty"`
    Version     string      `json:"version,omitempty"`
    Environment string      `json:"environment,omitempty"`
    Pagination  *Pagination `json:"pagination,omitempty"`
}
```

### Pagination

Pagination information:

```go
type Pagination struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"totalPages"`
    HasNext    bool  `json:"hasNext"`
    HasPrev    bool  `json:"hasPrev"`
}
```

## Usage Examples

### Success Responses

#### Simple Success Response

```go
func GetUserHandler(c *gin.Context) {
    user := User{ID: 1, Name: "John Doe", Email: "john@example.com"}
    
    response := common.CreateSuccessResponse(user, "User retrieved successfully")
    c.JSON(http.StatusOK, response)
}
```

#### Success Response with Custom Metadata

```go
func GetUserWithMetaHandler(c *gin.Context) {
    user := User{ID: 1, Name: "John Doe", Email: "john@example.com"}
    
    meta := &common.MetaInfo{
        RequestID:   "req-12345",
        Version:     "v1.2.3",
        Environment: "production",
    }
    
    response := common.CreateSuccessResponseWithMeta(user, "User retrieved successfully", meta)
    c.JSON(http.StatusOK, response)
}
```

#### Paginated Response

```go
func GetUsersHandler(c *gin.Context) {
    users := []User{...} // Your data
    page := 1
    limit := 10
    total := int64(100)
    
    response := common.CreatePaginatedResponse(users, page, limit, total)
    c.JSON(http.StatusOK, response)
}
```

### Error Responses

#### Not Found Error

```go
func GetUserHandler(c *gin.Context) {
    userID := c.Param("id")
    user, err := userService.GetByID(userID)
    if err != nil {
        response := common.CreateNotFoundResponse("User")
        c.JSON(http.StatusNotFound, response)
        return
    }
    
    response := common.CreateSuccessResponse(user, "User retrieved successfully")
    c.JSON(http.StatusOK, response)
}
```

#### Validation Error

```go
func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := []common.ValidationError{
            {
                Field:   "email",
                Message: "Email is required",
                Value:   "",
            },
            {
                Field:   "password",
                Message: "Password must be at least 8 characters",
                Value:   req.Password,
            },
        }
        
        response := common.CreateValidationErrorResponse(validationErrors)
        c.JSON(http.StatusBadRequest, response)
        return
    }
    
    // Process request...
}
```

#### Unauthorized Error

```go
func ProtectedHandler(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if !isValidToken(token) {
        response := common.CreateUnauthorizedResponse("Invalid or expired token")
        c.JSON(http.StatusUnauthorized, response)
        return
    }
    
    // Process request...
}
```

#### Forbidden Error

```go
func AdminHandler(c *gin.Context) {
    user := getUserFromContext(c)
    if !user.IsAdmin {
        response := common.CreateForbiddenResponse("Admin access required")
        c.JSON(http.StatusForbidden, response)
        return
    }
    
    // Process request...
}
```

#### Conflict Error

```go
func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    if userService.EmailExists(req.Email) {
        response := common.CreateConflictResponse("Email already exists", gin.H{
            "field": "email",
            "value": req.Email,
        })
        c.JSON(http.StatusConflict, response)
        return
    }
    
    // Process request...
}
```

#### Bad Request Error

```go
func UpdateUserHandler(c *gin.Context) {
    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response := common.CreateBadRequestResponse("Invalid request format", err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
    }
    
    // Process request...
}
```

#### Internal Server Error

```go
func GetUserHandler(c *gin.Context) {
    user, err := userService.GetByID(userID)
    if err != nil {
        response := common.CreateInternalServerErrorResponse("Database error", gin.H{
            "error": err.Error(),
            "retry": true,
        })
        c.JSON(http.StatusInternalServerError, response)
        return
    }
    
    response := common.CreateSuccessResponse(user, "User retrieved successfully")
    c.JSON(http.StatusOK, response)
}
```

## Migration from Legacy Structures

### Before (Legacy)

```go
// Old way
c.JSON(http.StatusOK, gin.H{
    "status":  "success",
    "message": "User retrieved successfully",
    "data":    user,
})
```

### After (New)

```go
// New way
response := common.CreateSuccessResponse(user, "User retrieved successfully")
c.JSON(http.StatusOK, response)
```

## Response Examples

### Success Response

```json
{
  "apiVersion": "v1",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### Paginated Response

```json
{
  "apiVersion": "v1",
  "data": [
    {"id": 1, "name": "John Doe"},
    {"id": 2, "name": "Jane Smith"}
  ],
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "totalPages": 3,
      "hasNext": true,
      "hasPrev": false
    }
  }
}
```

### Error Response

```json
{
  "apiVersion": "v1",
  "error": {
    "code": 404,
    "codeText": "Not Found",
    "message": "User not found"
  },
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### Validation Error Response

```json
{
  "apiVersion": "v1",
  "error": {
    "code": 400,
    "codeText": "Bad Request",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email is required",
        "value": ""
      },
      {
        "field": "password",
        "message": "Password must be at least 8 characters",
        "value": "123"
      }
    ]
  },
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

## Best Practices

1. **Always use the helper functions** instead of manually constructing responses
2. **Include meaningful error messages** that help clients understand what went wrong
3. **Use appropriate HTTP status codes** with the response helpers
4. **Add request IDs** for better debugging and tracing
5. **Include pagination metadata** for list endpoints
6. **Validate input** and return structured validation errors
7. **Don't expose internal errors** in production - use generic messages

## Integration with Services

### Identity Service

The Identity service has been updated to use these response structures. All handlers should use the common response helpers.

### Trip Service

The Trip service has been updated to use these response structures. All handlers should use the common response helpers.

### Adding to New Services

1. Add the common dependency to your `go.mod`:
   ```go
   require github.com/iamarpitzala/motocabz/common v0.0.0
   replace github.com/iamarpitzala/motocabz/common => C:\mihir\motocabz\Common
   ```

2. Import the common package in your handlers:
   ```go
   import "github.com/iamarpitzala/motocabz/common"
   ```

3. Use the response helpers in your handlers as shown in the examples above.
#   m o t o c a b z - c o m m o n  
 