/**
 * Standardized API Error Codes
 * 
 * Format: [CATEGORY]_[DESCRIPTION]
 */

export const ErrorCode = {
    // --- System Errors (SYS) ---
    UNKNOWN_ERROR: 'SYS_001',
    NETWORK_ERROR: 'SYS_002',
    TIMEOUT: 'SYS_003',
    SERVER_ERROR: 'SYS_500',
    MAINTENANCE: 'SYS_503',

    // --- Auth Errors (AUTH) ---
    UNAUTHORIZED: 'AUTH_401',
    FORBIDDEN: 'AUTH_403',
    TOKEN_EXPIRED: 'AUTH_001',
    INVALID_CREDENTIALS: 'AUTH_002',
    USER_NOT_FOUND: 'AUTH_003',
    SESSION_EXPIRED: 'AUTH_004',

    // --- Validation Errors (VAL) ---
    INVALID_PARAMS: 'VAL_400',
    MISSING_FIELD: 'VAL_001',
    SCHEMA_MISMATCH: 'VAL_002',

    // --- Business Errors (BIZ) ---
    RESOURCE_NOT_FOUND: 'BIZ_404',
    ALREADY_EXISTS: 'BIZ_001',
    OPERATION_FAILED: 'BIZ_002',
    QUOTA_EXCEEDED: 'BIZ_003',
} as const;

export type ErrorCodeValue = typeof ErrorCode[keyof typeof ErrorCode];

/**
 * Maps HTTP Status codes to default Error Codes
 */
export const HttpStatusMap: Record<number, ErrorCodeValue> = {
    400: ErrorCode.INVALID_PARAMS,
    401: ErrorCode.UNAUTHORIZED,
    403: ErrorCode.FORBIDDEN,
    404: ErrorCode.RESOURCE_NOT_FOUND,
    500: ErrorCode.SERVER_ERROR,
    503: ErrorCode.MAINTENANCE,
};
