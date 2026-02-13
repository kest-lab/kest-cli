# Sentry Go SDK Integration Example

This example demonstrates how to use the official Sentry Go SDK with Trac API.

## Prerequisites

- Go 1.21 or later
- A running Trac API instance
- A project created in Trac with a valid DSN

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Update the DSN in `main.go`:
```go
Dsn: "http://YOUR_PUBLIC_KEY@YOUR_SERVER:8025/YOUR_PROJECT_ID",
```

3. Run the example:
```bash
go run main.go
```

## What This Example Does

The example demonstrates four types of event capture:

1. **Simple Message**: Basic message logging
2. **Exception**: Error/exception capture
3. **Rich Context**: Message with user info, tags, and custom context
4. **Panic Recovery**: Automatic panic capture and recovery

## Expected Output

```
Test 1: Capturing a message...
Test 2: Capturing an exception...
Test 3: Capturing with context...
Test 4: Recovering from panic...

All events sent! Check your Trac dashboard.
Waiting for events to be flushed...
```

## Verifying Events

After running the example, verify events were captured:

### Using Trac API

```bash
# Get your access token
TOKEN=$(curl -s -X POST http://YOUR_SERVER:8025/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"YOUR_USERNAME","password":"YOUR_PASSWORD"}' \
  | jq -r '.data.access_token')

# List issues for your project
curl -s http://YOUR_SERVER:8025/v1/projects/YOUR_PROJECT_ID/issues/ \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Using ClickHouse

```bash
docker exec zgo-clickhouse clickhouse-client -u trac_user --password trac_pass -d trac -q \
  "SELECT event_id, level, message, timestamp FROM events WHERE project_id = YOUR_PROJECT_ID ORDER BY timestamp DESC LIMIT 10"
```

## Configuration Options

The Sentry SDK supports many configuration options:

```go
sentry.Init(sentry.ClientOptions{
    Dsn:              "http://...",
    Environment:      "production",     // Environment name
    Release:          "my-app@1.0.0",   // Release version
    Debug:            true,              // Enable debug logging
    SampleRate:       1.0,               // Sample rate (0.0 to 1.0)
    TracesSampleRate: 1.0,               // Traces sample rate
    BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
        // Modify or filter events before sending
        return event
    },
})
```

## Troubleshooting

### Events not appearing?

1. Check if the SDK is sending events:
   - Enable `Debug: true` in ClientOptions
   - Look for HTTP requests in the debug output

2. Verify the DSN is correct:
   - Format: `http://PUBLIC_KEY@HOST:PORT/PROJECT_ID`
   - Check if the public key matches your project

3. Check server logs:
   ```bash
   docker logs zgo-api --tail 50
   ```

### Connection refused?

- Ensure Trac API is running: `curl http://YOUR_SERVER:8025/health`
- Check firewall rules if using a remote server
- Verify the port (default: 8025)

## Next Steps

- Integrate Sentry SDK into your application
- Configure error filtering and sampling
- Set up custom tags and contexts
- Implement breadcrumbs for better debugging
- Configure release tracking

## Resources

- [Sentry Go SDK Documentation](https://docs.sentry.io/platforms/go/)
- [Trac API Documentation](../../docs/api.md)
- [Sentry Protocol Specification](https://develop.sentry.dev/sdk/event-payloads/)
