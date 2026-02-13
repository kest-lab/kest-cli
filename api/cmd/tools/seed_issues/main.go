package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	// Initialize Sentry SDK pointing to our local trac server
	dsn := "http://test_public_key@localhost:8025/v1/1"

	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(5 * time.Second)

	// Group 1: "database-connection-error" - Send 5 events
	log.Println("Sending 5 events for fingerprint: database-connection-error")
	for i := 1; i <= 5; i++ {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			scope.SetTag("error_type", "database")
			scope.SetFingerprint([]string{"database-connection-error"})
			sentry.CaptureException(fmt.Errorf("failed to connect to database (occurrence %d)", i))
		})
		time.Sleep(100 * time.Millisecond)
	}

	// Group 2: "api-timeout" - Send 3 events
	log.Println("Sending 3 events for fingerprint: api-timeout")
	for i := 1; i <= 3; i++ {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelWarning)
			scope.SetTag("error_type", "timeout")
			scope.SetFingerprint([]string{"api-timeout"})
			sentry.CaptureException(fmt.Errorf("API request timeout (occurrence %d)", i))
		})
		time.Sleep(100 * time.Millisecond)
	}

	// Group 3: "null-pointer" - Send 7 events
	log.Println("Sending 7 events for fingerprint: null-pointer")
	for i := 1; i <= 7; i++ {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelFatal)
			scope.SetTag("error_type", "panic")
			scope.SetFingerprint([]string{"null-pointer"})
			sentry.CaptureException(fmt.Errorf("null pointer dereference (occurrence %d)", i))
		})
		time.Sleep(100 * time.Millisecond)
	}

	// Group 4: "auth-failed" - Send 2 events
	log.Println("Sending 2 events for fingerprint: auth-failed")
	for i := 1; i <= 2; i++ {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			scope.SetTag("error_type", "auth")
			scope.SetFingerprint([]string{"auth-failed"})
			sentry.CaptureException(fmt.Errorf("authentication failed (occurrence %d)", i))
		})
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("All events sent! Flushing...")
	sentry.Flush(5 * time.Second)

	log.Println("\nEvent Summary:")
	log.Println("  - database-connection-error: 5 events")
	log.Println("  - api-timeout: 3 events")
	log.Println("  - null-pointer: 7 events")
	log.Println("  - auth-failed: 2 events")
	log.Println("\nTotal: 17 events in 4 issue groups")
	log.Println("\nNext steps:")
	log.Println("1. Wait ~10 seconds for materialized view to aggregate")
	log.Println("2. Query issues: curl http://localhost:8025/v1/projects/1/issues")
	log.Println("3. Check ClickHouse: docker exec zgo-clickhouse clickhouse-client -u trac_user --password trac_pass -d trac -q \"SELECT fingerprint, event_count FROM issues WHERE project_id = 1\"")
}
