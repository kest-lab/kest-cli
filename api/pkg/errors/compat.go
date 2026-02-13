package errors

import "net/http"

// ==========================================
// Backward Compatibility Layer
// ==========================================
// This file provides backward compatibility with the old error types.
// New code should use Error and ErrorCode from error.go and codes.go.

// Code is the legacy string-based error code type
type Code string

// Legacy string error codes
const (
	CodeUnknown            Code = "UNKNOWN"
	CodeValidation         Code = "VALIDATION_ERROR"
	CodeNotFound           Code = "NOT_FOUND"
	CodeUnauthorized       Code = "UNAUTHORIZED"
	CodeForbidden          Code = "FORBIDDEN"
	CodeConflict           Code = "CONFLICT"
	CodeTooManyRequests    Code = "TOO_MANY_REQUESTS"
	CodeInternal           Code = "INTERNAL_ERROR"
	CodeBadRequest         Code = "BAD_REQUEST"
	CodeUnprocessable      Code = "UNPROCESSABLE_ENTITY"
	CodeServiceUnavailable Code = "SERVICE_UNAVAILABLE"
	CodeTimeout            Code = "TIMEOUT"
)

// AppError is the legacy error type for backward compatibility
// New code should use Error instead
type AppError struct {
	// HTTP status code
	Status int `json:"-"`

	// Error code for client identification
	Code Code `json:"code"`

	// Human-readable error message
	Message string `json:"message"`

	// Detailed error description (shown in debug mode)
	Detail string `json:"detail,omitempty"`

	// Validation errors (field -> messages)
	Errors map[string][]string `json:"errors,omitempty"`

	// Additional metadata
	Meta map[string]interface{} `json:"meta,omitempty"`

	// Internal error (not exposed to client)
	Internal error `json:"-"`

	// Stack trace (for debugging)
	Stack string `json:"-"`

	// Source file and line
	File string `json:"-"`
	Line int    `json:"-"`
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return e.Message + ": " + e.Internal.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Internal
}

// WithDetail adds detail to the error
func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

// WithMeta adds metadata to the error
func (e *AppError) WithMeta(key string, value interface{}) *AppError {
	if e.Meta == nil {
		e.Meta = make(map[string]interface{})
	}
	e.Meta[key] = value
	return e
}

// WithInternal sets the internal error
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}

// WithErrors sets validation errors
func (e *AppError) WithErrors(errors map[string][]string) *AppError {
	e.Errors = errors
	return e
}

// Legacy constructors for backward compatibility

// LegacyInternal creates a legacy 500 error
func LegacyInternal(message string) *AppError {
	if message == "" {
		message = "An internal error occurred"
	}
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    CodeInternal,
		Message: message,
	}
}

// LegacyValidationWithErrors creates a legacy validation error
func LegacyValidationWithErrors(errors map[string][]string) *AppError {
	return &AppError{
		Status:  http.StatusUnprocessableEntity,
		Code:    CodeValidation,
		Message: "The given data was invalid",
		Errors:  errors,
	}
}
