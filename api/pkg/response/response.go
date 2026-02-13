package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// Success Responses (Zero-Intrusion API)
// ============================================================================

// Success sends a successful response with any data.
// Automatically handles pagination if data is a Paginator.
//
// Example:
//
//	// Simple data
//	response.Success(c, user)
//	// Output: {"code": 0, "message": "success", "data": {...}}
//
//	// Paginated data
//	users, paginator, _ := pagination.New[User](c, db)
//	response.Success(c, paginator)
//	// Output: {"code": 0, "message": "success", "data": [...], "meta": {...}, "links": {...}}
func Success(c *gin.Context, data any) {
	// Check if data is a Paginator (implements Paginatable with Items)
	if p, ok := data.(PaginatableWithItems); ok {
		c.JSON(http.StatusOK, PaginatedResponse{
			Code:    0,
			Message: "success",
			Data:    p.GetItems(),
			Meta:    p.GetMeta(),
			Links:   p.GetLinks(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// OK is an alias for Success.
func OK(c *gin.Context, data any) {
	Success(c, data)
}

// Created sends a 201 Created response.
// Use this after successfully creating a new resource.
//
// Example:
//
//	response.Created(c, newUser)
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
// Use this for successful operations that don't return data (e.g., DELETE).
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Accepted sends a 202 Accepted response.
// Use this for async operations that have been queued.
func Accepted(c *gin.Context, data any) {
	c.JSON(http.StatusAccepted, Response{
		Code:    0,
		Message: "accepted",
		Data:    data,
	})
}

// ============================================================================
// Paginated Responses (Zero-Intrusion)
// ============================================================================

// Paginated sends a paginated response.
// Deprecated: Use Success(c, paginator) instead - it auto-detects pagination.
//
// Example:
//
//	users, paginator, _ := pagination.New[User](c, db)
//	response.Success(c, paginator)  // Preferred
//	response.Paginated(c, users, paginator)  // Still works
func Paginated(c *gin.Context, data any, paginator Paginatable) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Meta:    paginator.GetMeta(),
		Links:   paginator.GetLinks(),
	})
}

// List sends a simple list response without pagination.
// Use this for small collections that don't need pagination.
//
// Example:
//
//	response.List(c, categories)
func List(c *gin.Context, data any) {
	Success(c, data)
}

// ============================================================================
// Transform Responses (Optional - for complex transformations)
// ============================================================================

// Transform sends data after applying a transformation function.
// Use this when you need to transform data without defining a Resource type.
//
// Example:
//
//	response.Transform(c, user, func(u *User) any {
//	    return map[string]any{
//	        "id":       u.ID,
//	        "username": u.Username,
//	    }
//	})
func Transform[T any](c *gin.Context, data T, fn func(T) any) {
	Success(c, fn(data))
}

// TransformList sends a list after applying transformation to each item.
//
// Example:
//
//	response.TransformList(c, users, func(u *User) any {
//	    return map[string]any{"id": u.ID, "name": u.Username}
//	})
func TransformList[T any](c *gin.Context, items []T, fn func(T) any) {
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	Success(c, result)
}

// TransformPaginated sends paginated data with transformation.
//
// Example:
//
//	users, paginator, _ := pagination.New[*User](c, db)
//	response.TransformPaginated(c, users, paginator, func(u *User) any {
//	    return map[string]any{"id": u.ID, "name": u.Username}
//	})
func TransformPaginated[T any](c *gin.Context, items []T, paginator Paginatable, fn func(T) any) {
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	Paginated(c, result, paginator)
}

// ============================================================================
// Resource Responses (Optional - for reusable transformations)
// ============================================================================

// Resource sends a response using a Resourceable transformer.
// Use this when you have defined a reusable Resource type.
//
// Example:
//
//	response.Resource(c, NewUserResource(user))
func Resource(c *gin.Context, resource Resourceable) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    resource.ToArray(),
	})
}

// ResourceCreated sends a created response using a Resource.
func ResourceCreated(c *gin.Context, resource Resourceable) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    resource.ToArray(),
	})
}

// Collection sends a collection response using a Collectable.
// Use this when you have defined a reusable Collection type.
//
// Example:
//
//	response.Collection(c, NewUserCollection(users, paginator))
func Collection(c *gin.Context, collection Collectable) {
	paginator := collection.GetPaginator()
	if paginator != nil {
		c.JSON(http.StatusOK, PaginatedResponse{
			Code:    0,
			Message: "success",
			Data:    collection.ToArray(),
			Meta:    paginator.GetMeta(),
			Links:   paginator.GetLinks(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    collection.ToArray(),
	})
}

// ============================================================================
// Error Responses
// ============================================================================

// Error sends a generic error response.
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Code:    statusCode,
		Message: message,
	})
}

// ErrorWithDetails sends an error response with error details.
func ErrorWithDetails(c *gin.Context, statusCode int, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	c.JSON(statusCode, ErrorResponse{
		Code:    statusCode,
		Message: message,
		Error:   errMsg,
	})
}

// BadRequest sends a 400 Bad Request response.
// Use this for invalid request data.
//
// Example:
//
//	response.BadRequest(c, "Invalid email format")
func BadRequest(c *gin.Context, message string, err ...error) {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	ErrorWithDetails(c, http.StatusBadRequest, message, e)
}

// Unauthorized sends a 401 Unauthorized response.
// Use this when authentication is required but not provided or invalid.
func Unauthorized(c *gin.Context, message ...string) {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(c, http.StatusUnauthorized, msg)
}

// Forbidden sends a 403 Forbidden response.
// Use this when the user is authenticated but lacks permission.
func Forbidden(c *gin.Context, message ...string) {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(c, http.StatusForbidden, msg)
}

// NotFound sends a 404 Not Found response.
// Use this when the requested resource doesn't exist.
//
// Example:
//
//	response.NotFound(c, "User not found")
func NotFound(c *gin.Context, message string, err ...error) {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	ErrorWithDetails(c, http.StatusNotFound, message, e)
}

// MethodNotAllowed sends a 405 Method Not Allowed response.
func MethodNotAllowed(c *gin.Context) {
	Error(c, http.StatusMethodNotAllowed, "Method not allowed")
}

// Conflict sends a 409 Conflict response.
// Use this for duplicate entries or conflicting operations.
func Conflict(c *gin.Context, message string, err ...error) {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	ErrorWithDetails(c, http.StatusConflict, message, e)
}

// UnprocessableEntity sends a 422 Unprocessable Entity response.
// Use this for semantic validation errors.
func UnprocessableEntity(c *gin.Context, message string, err ...error) {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	ErrorWithDetails(c, http.StatusUnprocessableEntity, message, e)
}

// ValidationFailed sends a 422 response with field-level errors.
// Use this for form validation failures.
//
// Example:
//
//	response.ValidationFailed(c, map[string][]string{
//	    "email": {"The email field is required"},
//	    "password": {"The password must be at least 8 characters"},
//	})
func ValidationFailed(c *gin.Context, errors map[string][]string) {
	c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: "Validation failed",
		Errors:  errors,
	})
}

// TooManyRequests sends a 429 Too Many Requests response.
// Use this for rate limiting.
func TooManyRequests(c *gin.Context, message ...string) {
	msg := "Too many requests"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(c, http.StatusTooManyRequests, msg)
}

// InternalServerError sends a 500 Internal Server Error response.
// Use this for unexpected server errors.
func InternalServerError(c *gin.Context, message string, err ...error) {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	ErrorWithDetails(c, http.StatusInternalServerError, message, e)
}

// ServiceUnavailable sends a 503 Service Unavailable response.
// Use this when a dependent service is down.
func ServiceUnavailable(c *gin.Context, message ...string) {
	msg := "Service unavailable"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(c, http.StatusServiceUnavailable, msg)
}
