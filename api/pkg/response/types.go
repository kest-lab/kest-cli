// Package response provides unified API response handling for the Eogo framework.
//
// This package implements Laravel-style API responses with support for:
//   - Unified response structure with code, message, and data
//   - Automatic pagination metadata (links, meta)
//   - Resource transformation for clean API output
//   - Error handling with proper HTTP status codes
//
// Basic Usage:
//
//	// Success response
//	response.Success(c, user)
//
//	// Paginated response
//	response.Paginated(c, users, paginator)
//
//	// Error response
//	response.NotFound(c, "User not found")
//
// With Resources:
//
//	// Single resource
//	response.Resource(c, NewUserResource(user))
//
//	// Collection with pagination
//	response.Collection(c, NewUserCollection(users, paginator))
package response

import "github.com/gin-gonic/gin"

// Response is the standard API response structure.
// All API endpoints should return this structure for consistency.
//
// Example JSON output:
//
//	{
//	    "code": 0,
//	    "message": "success",
//	    "data": {...}
//	}
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PaginatedResponse extends Response with pagination metadata.
// Used for list endpoints that return paginated data.
//
// Example JSON output:
//
//	{
//	    "code": 0,
//	    "message": "success",
//	    "data": [...],
//	    "meta": {
//	        "current_page": 1,
//	        "per_page": 15,
//	        "total": 100,
//	        "last_page": 7,
//	        "from": 1,
//	        "to": 15
//	    },
//	    "links": {
//	        "first": "http://api.example.com/users?page=1",
//	        "last": "http://api.example.com/users?page=7",
//	        "prev": null,
//	        "next": "http://api.example.com/users?page=2"
//	    }
//	}
type PaginatedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta"`
	Links   *Links `json:"links"`
}

// Meta contains pagination metadata.
type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

// Links contains pagination URLs.
type Links struct {
	First string  `json:"first"`
	Last  string  `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

// ErrorResponse is returned for error cases.
// Includes optional error details for debugging.
//
// Example JSON output:
//
//	{
//	    "code": 404,
//	    "message": "User not found",
//	    "error": "record not found"
//	}
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ValidationErrorResponse is returned for validation failures.
// Includes field-level error details.
//
// Example JSON output:
//
//	{
//	    "code": 422,
//	    "message": "Validation failed",
//	    "errors": {
//	        "email": ["The email field is required", "The email must be valid"],
//	        "password": ["The password must be at least 8 characters"]
//	    }
//	}
type ValidationErrorResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

// Responder interface for custom response types.
// Implement this interface to create custom response handlers.
type Responder interface {
	Respond(c *gin.Context)
}

// Resourceable interface for transforming models to API output.
// Implement this interface on your resource structs.
//
// Example:
//
//	type UserResource struct {
//	    user *domain.User
//	}
//
//	func (r *UserResource) ToArray() map[string]any {
//	    return map[string]any{
//	        "id":       r.user.ID,
//	        "username": r.user.Username,
//	        "email":    r.user.Email,
//	    }
//	}
type Resourceable interface {
	ToArray() map[string]any
}

// Collectable interface for resource collections.
// Implement this for paginated list responses.
type Collectable interface {
	ToArray() []map[string]any
	GetPaginator() Paginatable
}

// Paginatable interface for pagination data.
// The pagination package implements this interface.
type Paginatable interface {
	GetMeta() *Meta
	GetLinks() *Links
}

// PaginatableWithItems extends Paginatable with items access.
// Used by Success() to auto-detect paginated responses.
type PaginatableWithItems interface {
	Paginatable
	GetItems() any
}
