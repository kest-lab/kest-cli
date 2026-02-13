package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	// 1. Initialize Sentry SDK
	// Replace with your actual project ID and public key if different
	// Format: http://{public_key}@{host}:{port}/{project_id}
	// Note: /v1 prefix is currently required due to our route registration
	dsn := "http://test_public_key@localhost:8025/v1/1"

	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: true,
		// Performance monitoring
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// 2. Capture a message
	sentry.CaptureMessage("Hello Trac! This is a test message from Sentry Go SDK.")
	log.Println("Message captured")

	// 3. Capture an exception
	errToCapture := fmt.Errorf("something went wrong in the test")
	sentry.CaptureException(errToCapture)
	log.Println("Exception captured")

	// 4. Flush events to the server
	sentry.Flush(5 * time.Second)
	log.Println("Events flushed. Check ClickHouse 'events' table for results.")
}
