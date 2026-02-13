# Gin Error Reporting Example

This example demonstrates how to integrate a Gin web application with trac error reporting.

## Features

- Automatic panic recovery with stack trace
- Error capture and reporting to trac API
- Sentry-compatible envelope format
- Multiple test endpoints to simulate different error scenarios

## Setup

1. Install dependencies:
```bash
cd examples/gin-error-reporting
go mod tidy
```

2. Set the trac DSN (optional, defaults to test project):
```bash
export TRAC_DSN="http://44da05e841a422acd8b414dbb452e88a@localhost:8025/2"
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Test Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Test Panic (will be reported to trac)
```bash
curl http://localhost:8080/test/panic
```

### Test Error (will be reported to trac)
```bash
curl http://localhost:8080/test/error
```

### Test Division by Zero (will be reported to trac)
```bash
curl http://localhost:8080/test/divide
```

### Test Validation Error
```bash
curl -X POST http://localhost:8080/test/process \
  -H "Content-Type: application/json" \
  -d '{"name": "", "age": -1}'
```

### Test Processing Error
```bash
curl -X POST http://localhost:8080/test/process \
  -H "Content-Type: application/json" \
  -d '{"name": "error", "age": 25}'
```

### Valid Request
```bash
curl -X POST http://localhost:8080/test/process \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "age": 30}'
```

## How It Works

1. **Middleware**: The application uses Gin middleware to:
   - Recover from panics
   - Capture errors
   - Send them to trac API

2. **Error Reporting**: Errors are sent as Sentry-compatible envelopes to:
   ```
   POST /api/{project_id}/envelope/
   ```

3. **Event Data**: Each error event includes:
   - Stack trace
   - Request information
   - User context
   - Custom tags

## Integration in Your Project

To integrate this in your own Gin project:

1. Copy the `TracClient` struct and related code
2. Add the middleware to your router:
   ```go
   r.Use(tracMiddleware)
   ```
3. Set your project DSN in environment variables

## Verify Errors in Trac

Check the trac API to see reported errors:

```bash
# Login to get token
curl -X POST http://localhost:8025/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'

# List issues for the project
curl -X GET http://localhost:8025/v1/issues \
  -H "Authorization: Bearer YOUR_TOKEN"
```
