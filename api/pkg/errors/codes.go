package errors

// ErrorCode represents a numeric error code
// Format: MMMMYYYY where MMMM is module code (0-9999) and YYYY is error code (0-9999)
type ErrorCode int

// Module codes
const (
	ModuleSystem     = 0    // System/Framework errors (0000-0999)
	ModuleAuth       = 1000 // Authentication errors (1000-1999)
	ModuleUser       = 2000 // User errors (2000-2999)
	ModulePermission = 3000 // Permission/RBAC errors (3000-3999)
	// Business modules start from 4000
)

// ===== System Errors (0000-0999) =====
const (
	ErrUnknown         ErrorCode = 0  // Unknown error
	ErrInternal        ErrorCode = 1  // Internal server error
	ErrBadRequest      ErrorCode = 2  // Bad request
	ErrValidation      ErrorCode = 3  // Validation failed
	ErrNotFound        ErrorCode = 4  // Resource not found
	ErrForbidden       ErrorCode = 5  // Access forbidden
	ErrConflict        ErrorCode = 6  // Resource conflict
	ErrTooManyRequests ErrorCode = 7  // Rate limit exceeded
	ErrTimeout         ErrorCode = 8  // Request timeout
	ErrServiceUnavail  ErrorCode = 9  // Service unavailable
	ErrDatabase        ErrorCode = 10 // Database error
	ErrCache           ErrorCode = 11 // Cache error
	ErrQueue           ErrorCode = 12 // Queue error
	ErrThirdParty      ErrorCode = 13 // Third-party service error
	ErrConfiguration   ErrorCode = 14 // Configuration error
	ErrIO              ErrorCode = 15 // I/O error
	ErrSerialization   ErrorCode = 16 // Serialization error
	ErrNetwork         ErrorCode = 17 // Network error
)

// ===== Authentication Errors (1000-1999) =====
const (
	ErrUnauthorized     ErrorCode = 1000 // Not authenticated
	ErrTokenExpired     ErrorCode = 1001 // Token expired
	ErrTokenInvalid     ErrorCode = 1002 // Token invalid
	ErrTokenMalformed   ErrorCode = 1003 // Token malformed
	ErrTokenMissing     ErrorCode = 1004 // Token missing
	ErrLoginFailed      ErrorCode = 1005 // Login failed
	ErrLogoutFailed     ErrorCode = 1006 // Logout failed
	ErrRefreshFailed    ErrorCode = 1007 // Token refresh failed
	ErrAccountDisabled  ErrorCode = 1008 // Account disabled
	ErrAccountLocked    ErrorCode = 1009 // Account locked
	ErrAccountNotVerify ErrorCode = 1010 // Account not verified
	ErrSessionExpired   ErrorCode = 1011 // Session expired
	ErrSessionInvalid   ErrorCode = 1012 // Session invalid
)

// ===== User Errors (2000-2999) =====
const (
	ErrUserNotFound       ErrorCode = 2000 // User not found
	ErrUserExists         ErrorCode = 2001 // User already exists
	ErrEmailExists        ErrorCode = 2002 // Email already exists
	ErrUsernameExists     ErrorCode = 2003 // Username already exists
	ErrPasswordIncorrect  ErrorCode = 2004 // Incorrect password
	ErrPasswordTooWeak    ErrorCode = 2005 // Password too weak
	ErrPasswordMismatch   ErrorCode = 2006 // Password mismatch
	ErrEmailInvalid       ErrorCode = 2007 // Invalid email format
	ErrProfileIncomplete  ErrorCode = 2008 // Profile incomplete
	ErrAvatarUploadFailed ErrorCode = 2009 // Avatar upload failed
)

// ===== Permission/RBAC Errors (3000-3999) =====
const (
	ErrPermissionDenied   ErrorCode = 3000 // Permission denied
	ErrRoleNotFound       ErrorCode = 3001 // Role not found
	ErrRoleExists         ErrorCode = 3002 // Role already exists
	ErrRoleInUse          ErrorCode = 3003 // Role in use
	ErrPermissionNotFound ErrorCode = 3004 // Permission not found
	ErrNoRoleAssigned     ErrorCode = 3005 // No role assigned
	ErrCannotRemoveAdmin  ErrorCode = 3006 // Cannot remove admin role
)
