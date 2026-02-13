package errors

import (
	"net/http"
	"sync"
)

// ErrorDefinition holds metadata for an error code
type ErrorDefinition struct {
	Code        ErrorCode // Numeric code
	HTTPStatus  int       // HTTP status code
	Key         string    // String key (e.g., "UNAUTHORIZED")
	Message     string    // Default message
	MessageKey  string    // I18n key (e.g., "error.auth.unauthorized")
	Description string    // Documentation description
	Module      string    // Module name
}

var (
	registry     = make(map[ErrorCode]*ErrorDefinition)
	registryLock sync.RWMutex
)

// Register adds an error definition to the registry
func Register(def *ErrorDefinition) {
	registryLock.Lock()
	defer registryLock.Unlock()
	registry[def.Code] = def
}

// Get retrieves an error definition by code
func Get(code ErrorCode) *ErrorDefinition {
	registryLock.RLock()
	defer registryLock.RUnlock()
	return registry[code]
}

// All returns all registered error definitions
func All() []*ErrorDefinition {
	registryLock.RLock()
	defer registryLock.RUnlock()

	result := make([]*ErrorDefinition, 0, len(registry))
	for _, def := range registry {
		result = append(result, def)
	}
	return result
}

// AllByModule returns definitions grouped by module
func AllByModule() map[string][]*ErrorDefinition {
	registryLock.RLock()
	defer registryLock.RUnlock()

	result := make(map[string][]*ErrorDefinition)
	for _, def := range registry {
		result[def.Module] = append(result[def.Module], def)
	}
	return result
}

// init registers all framework error codes
func init() {
	// System errors
	Register(&ErrorDefinition{
		Code:        ErrUnknown,
		HTTPStatus:  http.StatusInternalServerError,
		Key:         "UNKNOWN",
		Message:     "An unknown error occurred",
		MessageKey:  "error.system.unknown",
		Description: "An unexpected error occurred",
		Module:      "system",
	})
	Register(&ErrorDefinition{
		Code:        ErrInternal,
		HTTPStatus:  http.StatusInternalServerError,
		Key:         "INTERNAL_ERROR",
		Message:     "Internal server error",
		MessageKey:  "error.system.internal",
		Description: "An internal server error occurred",
		Module:      "system",
	})
	Register(&ErrorDefinition{
		Code:        ErrBadRequest,
		HTTPStatus:  http.StatusBadRequest,
		Key:         "BAD_REQUEST",
		Message:     "Bad request",
		MessageKey:  "error.system.bad_request",
		Description: "The request was malformed or invalid",
		Module:      "system",
	})
	Register(&ErrorDefinition{
		Code:        ErrValidation,
		HTTPStatus:  http.StatusUnprocessableEntity,
		Key:         "VALIDATION_ERROR",
		Message:     "Validation failed",
		MessageKey:  "error.system.validation",
		Description: "The provided data failed validation",
		Module:      "system",
	})
	Register(&ErrorDefinition{
		Code:        ErrNotFound,
		HTTPStatus:  http.StatusNotFound,
		Key:         "NOT_FOUND",
		Message:     "Resource not found",
		MessageKey:  "error.system.not_found",
		Description: "The requested resource was not found",
		Module:      "system",
	})
	Register(&ErrorDefinition{
		Code:        ErrForbidden,
		HTTPStatus:  http.StatusForbidden,
		Key:         "FORBIDDEN",
		Message:     "Access forbidden",
		MessageKey:  "error.system.forbidden",
		Description: "You do not have permission to access this resource",
		Module:      "system",
	})
	Register(&ErrorDefinition{
		Code:        ErrTooManyRequests,
		HTTPStatus:  http.StatusTooManyRequests,
		Key:         "TOO_MANY_REQUESTS",
		Message:     "Too many requests",
		MessageKey:  "error.system.rate_limit",
		Description: "Rate limit exceeded, please try again later",
		Module:      "system",
	})

	// Auth errors
	Register(&ErrorDefinition{
		Code:        ErrUnauthorized,
		HTTPStatus:  http.StatusUnauthorized,
		Key:         "UNAUTHORIZED",
		Message:     "Authentication required",
		MessageKey:  "error.auth.unauthorized",
		Description: "Valid authentication credentials are required",
		Module:      "auth",
	})
	Register(&ErrorDefinition{
		Code:        ErrTokenExpired,
		HTTPStatus:  http.StatusUnauthorized,
		Key:         "TOKEN_EXPIRED",
		Message:     "Token has expired",
		MessageKey:  "error.auth.token_expired",
		Description: "The authentication token has expired",
		Module:      "auth",
	})
	Register(&ErrorDefinition{
		Code:        ErrTokenInvalid,
		HTTPStatus:  http.StatusUnauthorized,
		Key:         "TOKEN_INVALID",
		Message:     "Token is invalid",
		MessageKey:  "error.auth.token_invalid",
		Description: "The authentication token is invalid",
		Module:      "auth",
	})
	Register(&ErrorDefinition{
		Code:        ErrAccountDisabled,
		HTTPStatus:  http.StatusForbidden,
		Key:         "ACCOUNT_DISABLED",
		Message:     "Account is disabled",
		MessageKey:  "error.auth.account_disabled",
		Description: "This account has been disabled",
		Module:      "auth",
	})

	// User errors
	Register(&ErrorDefinition{
		Code:        ErrUserNotFound,
		HTTPStatus:  http.StatusNotFound,
		Key:         "USER_NOT_FOUND",
		Message:     "User not found",
		MessageKey:  "error.user.not_found",
		Description: "The specified user was not found",
		Module:      "user",
	})
	Register(&ErrorDefinition{
		Code:        ErrUserExists,
		HTTPStatus:  http.StatusConflict,
		Key:         "USER_EXISTS",
		Message:     "User already exists",
		MessageKey:  "error.user.exists",
		Description: "A user with this identifier already exists",
		Module:      "user",
	})
	Register(&ErrorDefinition{
		Code:        ErrEmailExists,
		HTTPStatus:  http.StatusConflict,
		Key:         "EMAIL_EXISTS",
		Message:     "Email already exists",
		MessageKey:  "error.user.email_exists",
		Description: "This email address is already registered",
		Module:      "user",
	})
	Register(&ErrorDefinition{
		Code:        ErrPasswordIncorrect,
		HTTPStatus:  http.StatusUnauthorized,
		Key:         "PASSWORD_INCORRECT",
		Message:     "Incorrect password",
		MessageKey:  "error.user.password_incorrect",
		Description: "The provided password is incorrect",
		Module:      "user",
	})

	// Permission errors
	Register(&ErrorDefinition{
		Code:        ErrPermissionDenied,
		HTTPStatus:  http.StatusForbidden,
		Key:         "PERMISSION_DENIED",
		Message:     "Permission denied",
		MessageKey:  "error.permission.denied",
		Description: "You do not have the required permissions",
		Module:      "permission",
	})
	Register(&ErrorDefinition{
		Code:        ErrRoleNotFound,
		HTTPStatus:  http.StatusNotFound,
		Key:         "ROLE_NOT_FOUND",
		Message:     "Role not found",
		MessageKey:  "error.permission.role_not_found",
		Description: "The specified role was not found",
		Module:      "permission",
	})
}
