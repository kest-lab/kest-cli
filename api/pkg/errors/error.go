package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// Error represents an application error with rich context
type Error struct {
	// Numeric error code
	Code ErrorCode `json:"code"`

	// String error key (for backward compatibility)
	Key string `json:"key,omitempty"`

	// HTTP status code
	HTTPStatus int `json:"-"`

	// Human-readable error message
	Message string `json:"message"`

	// Detailed error description (debug mode only)
	Detail string `json:"detail,omitempty"`

	// Validation errors (field -> messages)
	Errors map[string][]string `json:"errors,omitempty"`

	// Additional metadata
	Meta map[string]interface{} `json:"meta,omitempty"`

	// Request ID for tracing
	RequestID string `json:"request_id,omitempty"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`

	// Internal error (not exposed to client)
	internal error

	// Stack trace (for debugging)
	stack string

	// Source file and line
	file string
	line int
}

// Error implements error interface
func (e *Error) Error() string {
	if e.internal != nil {
		return fmt.Sprintf("[%d] %s: %s (internal: %v)", e.Code, e.Key, e.Message, e.internal)
	}
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Key, e.Message)
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.internal
}

// --- Builder Methods ---

// WithDetail adds detail to the error
func (e *Error) WithDetail(detail string) *Error {
	e.Detail = detail
	return e
}

// WithMeta adds metadata to the error
func (e *Error) WithMeta(key string, value interface{}) *Error {
	if e.Meta == nil {
		e.Meta = make(map[string]interface{})
	}
	e.Meta[key] = value
	return e
}

// WithInternal sets the internal error
func (e *Error) WithInternal(err error) *Error {
	e.internal = err
	return e
}

// WithRequestID sets the request ID
func (e *Error) WithRequestID(id string) *Error {
	e.RequestID = id
	return e
}

// WithErrors sets validation errors
func (e *Error) WithErrors(errors map[string][]string) *Error {
	e.Errors = errors
	return e
}

// AddError adds a validation error for a field
func (e *Error) AddError(field, message string) *Error {
	if e.Errors == nil {
		e.Errors = make(map[string][]string)
	}
	e.Errors[field] = append(e.Errors[field], message)
	return e
}

// --- Constructors ---

// newError creates a new Error with stack trace
func newError(code ErrorCode) *Error {
	def := Get(code)

	e := &Error{
		Code:      code,
		Timestamp: time.Now(),
	}

	if def != nil {
		e.HTTPStatus = def.HTTPStatus
		e.Key = def.Key
		e.Message = def.Message
	} else {
		e.HTTPStatus = http.StatusInternalServerError
		e.Key = "UNKNOWN"
		e.Message = "An unknown error occurred"
	}

	// Capture caller info
	if _, file, line, ok := runtime.Caller(2); ok {
		e.file = file
		e.line = line
	}

	// Capture stack trace
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	e.stack = string(buf[:n])

	return e
}

// New creates an error from error code
func New(code ErrorCode) *Error {
	return newError(code)
}

// NewWithMessage creates an error with custom message
func NewWithMessage(code ErrorCode, message string) *Error {
	e := newError(code)
	e.Message = message
	return e
}

// Wrap wraps an existing error
func Wrap(code ErrorCode, err error) *Error {
	if err == nil {
		return nil
	}

	// If already our Error type, enhance it
	if appErr, ok := err.(*Error); ok {
		return appErr
	}

	e := newError(code)
	e.internal = err
	return e
}

// --- Convenience Constructors ---

// BadRequest creates a 400 error
func BadRequest(message string) *Error {
	return NewWithMessage(ErrBadRequest, message)
}

// Unauthorized creates a 401 error
func Unauthorized(message ...string) *Error {
	e := New(ErrUnauthorized)
	if len(message) > 0 && message[0] != "" {
		e.Message = message[0]
	}
	return e
}

// Forbidden creates a 403 error
func Forbidden(message ...string) *Error {
	e := New(ErrForbidden)
	if len(message) > 0 && message[0] != "" {
		e.Message = message[0]
	}
	return e
}

// NotFound creates a 404 error
func NotFound(message ...string) *Error {
	e := New(ErrNotFound)
	if len(message) > 0 && message[0] != "" {
		e.Message = message[0]
	}
	return e
}

// Validation creates a 422 error
func Validation(message ...string) *Error {
	e := New(ErrValidation)
	if len(message) > 0 && message[0] != "" {
		e.Message = message[0]
	}
	return e
}

// ValidationWithErrors creates a validation error with field errors
func ValidationWithErrors(errors map[string][]string) *Error {
	e := New(ErrValidation)
	e.Errors = errors
	return e
}

// Internal creates a 500 error
func Internal(message ...string) *Error {
	e := New(ErrInternal)
	if len(message) > 0 && message[0] != "" {
		e.Message = message[0]
	}
	return e
}

// TooManyRequests creates a 429 error
func TooManyRequests(message ...string) *Error {
	e := New(ErrTooManyRequests)
	if len(message) > 0 && message[0] != "" {
		e.Message = message[0]
	}
	return e
}

// --- Error Checking ---

// Is checks if an error has a specific code
func Is(err error, code ErrorCode) bool {
	if appErr, ok := err.(*Error); ok {
		return appErr.Code == code
	}
	return false
}

// GetCode returns the error code from an error
func GetCode(err error) ErrorCode {
	if appErr, ok := err.(*Error); ok {
		return appErr.Code
	}
	return ErrUnknown
}

// GetHTTPStatus returns the HTTP status code from an error
func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*Error); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
