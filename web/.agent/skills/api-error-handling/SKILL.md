---
name: api-error-handling
description: Standardized error format, error code clusters, and API client usage for consistent error handling.
---

# api-error-handling

## Overview

This skill ensures consistent error handling between the frontend and backend. It defines a standardized error response format and a system of error code clusters to categorize and handle issues effectively.

## Guidelines

### 1. Error Response Format
All error responses from the backend (BFF or mock) MUST follow this format:

```json
{
  "error": "Human-readable error message",
  "code": "CATEGORY_DESCRIPTION"
}
```

### 2. Error Code Clusters
Error codes follow the format `[CATEGORY]_[DESCRIPTION]`. Use the `ErrorCode` constant from `@/http/codes` for comparison.

| Category | Prefix | Description | Examples |
|----------|--------|-------------|----------|
| **System** | `SYS_` | Infrastructural or unexpected errors | `SYS_SERVER_ERROR`, `SYS_TIMEOUT` |
| **Auth** | `AUTH_` | Authentication and authorization issues | `AUTH_TOKEN_EXPIRED`, `AUTH_FORBIDDEN` |
| **Validation** | `VAL_` | Input parameter or schema validation | `VAL_MISSING_FIELD`, `VAL_INVALID_PARAMS` |
| **Business** | `BIZ_` | Business logic constraints | `BIZ_RESOURCE_NOT_FOUND`, `BIZ_ALREADY_EXISTS` |

### 3. Usage in Code
- **Constants**: Always use `ErrorCode` from `@/http/codes` to check for specific error types.
- **Auto-handling**: Errors are automatically caught by `HttpClient` and passed to `handleError` (toast notification, etc.).
- **Manual Handling**: Use `skipErrorHandler: true` in the request configuration to handle errors manually within components or hooks.

```typescript
try {
  await request.post('/resource', data, { skipErrorHandler: true });
} catch (error) {
  if (error.code === ErrorCode.BIZ_ALREADY_EXISTS) {
    // Custom logic
  }
}
```

> [!IMPORTANT]
> **Consistency**: Ensure that backend error codes align with the frontend constants to maintain a robust UX.
