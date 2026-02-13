package events

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// TestLoggingMiddleware tests the logging middleware
func TestLoggingMiddleware(t *testing.T) {
	logger := NewSlogAdapter(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	middleware := LoggingMiddleware(logger, LogLevelDebug)

	called := false
	handler := middleware(func(ctx context.Context, event Event) error {
		called = true
		return nil
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

// TestLoggingMiddleware_Error tests logging middleware with error
func TestLoggingMiddleware_Error(t *testing.T) {
	logger := NewSlogAdapter(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	middleware := LoggingMiddleware(logger, LogLevelInfo)

	expectedErr := errors.New("test error")
	handler := middleware(func(ctx context.Context, event Event) error {
		return expectedErr
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err != expectedErr {
		t.Errorf("Expected error %v, got: %v", expectedErr, err)
	}
}

// TestRecoveryMiddleware tests panic recovery
func TestRecoveryMiddleware(t *testing.T) {
	var recovered interface{}
	var recoveredStack []byte

	middleware := RecoveryMiddleware(func(event Event, r interface{}, stack []byte) {
		recovered = r
		recoveredStack = stack
	})

	handler := middleware(func(ctx context.Context, event Event) error {
		panic("test panic")
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err == nil {
		t.Error("Expected error from panic recovery")
	}
	if recovered != "test panic" {
		t.Errorf("Expected recovered value 'test panic', got: %v", recovered)
	}
	if len(recoveredStack) == 0 {
		t.Error("Expected stack trace")
	}
}

// TestRetryMiddleware tests retry behavior
func TestRetryMiddleware(t *testing.T) {
	attempts := 0
	middleware := RetryMiddleware(2, 10*time.Millisecond)

	handler := middleware(func(ctx context.Context, event Event) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err != nil {
		t.Errorf("Expected success after retries, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

// TestRetryMiddleware_AllFail tests retry when all attempts fail
func TestRetryMiddleware_AllFail(t *testing.T) {
	attempts := 0
	middleware := RetryMiddleware(2, 1*time.Millisecond)

	handler := middleware(func(ctx context.Context, event Event) error {
		attempts++
		return errors.New("persistent error")
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err == nil {
		t.Error("Expected error after all retries failed")
	}
	if attempts != 3 { // 1 initial + 2 retries
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

// TestTimeoutMiddleware tests timeout behavior
func TestTimeoutMiddleware(t *testing.T) {
	middleware := TimeoutMiddleware(50 * time.Millisecond)

	handler := middleware(func(ctx context.Context, event Event) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
			return nil
		}
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestTimeoutMiddleware_Timeout tests timeout when handler is slow
func TestTimeoutMiddleware_Timeout(t *testing.T) {
	middleware := TimeoutMiddleware(10 * time.Millisecond)

	handler := middleware(func(ctx context.Context, event Event) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	event := newTestEvent("test.event", "payload")
	err := handler(context.Background(), event)

	if err == nil {
		t.Error("Expected timeout error")
	}
}

// TestFilterMiddleware tests event filtering
func TestFilterMiddleware(t *testing.T) {
	var called int32
	middleware := FilterMiddleware(func(event Event) bool {
		return event.EventName() == "allowed.event"
	})

	handler := middleware(func(ctx context.Context, event Event) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	// Should be filtered out
	event1 := newTestEvent("blocked.event", "payload")
	_ = handler(context.Background(), event1)

	// Should pass through
	event2 := newTestEvent("allowed.event", "payload")
	_ = handler(context.Background(), event2)

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected handler to be called once, got: %d", called)
	}
}

// TestMetricsMiddleware tests metrics recording
func TestMetricsMiddleware(t *testing.T) {
	recorder := &mockMetricsRecorder{}
	middleware := MetricsMiddleware(recorder)

	handler := middleware(func(ctx context.Context, event Event) error {
		return nil
	})

	event := newTestEvent("test.event", "payload")
	_ = handler(context.Background(), event)

	if recorder.eventName != "test.event" {
		t.Errorf("Expected event name 'test.event', got: %s", recorder.eventName)
	}
	if !recorder.success {
		t.Error("Expected success to be true")
	}
}

// mockMetricsRecorder implements MetricsRecorder for testing
type mockMetricsRecorder struct {
	eventName string
	duration  time.Duration
	success   bool
}

func (m *mockMetricsRecorder) RecordEventHandled(eventName string, duration time.Duration, success bool) {
	m.eventName = eventName
	m.duration = duration
	m.success = success
}

// TestMiddlewareChain tests multiple middleware in chain
func TestMiddlewareChain(t *testing.T) {
	var order []string

	m1 := func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			order = append(order, "m1-before")
			err := next(ctx, event)
			order = append(order, "m1-after")
			return err
		}
	}

	m2 := func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			order = append(order, "m2-before")
			err := next(ctx, event)
			order = append(order, "m2-after")
			return err
		}
	}

	bus := NewEventBus()
	bus.Use(m1, m2)
	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		order = append(order, "handler")
		return nil
	})

	event := newTestEvent("test", "payload")
	_ = bus.Publish(context.Background(), event)

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Errorf("Expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if i < len(order) && order[i] != v {
			t.Errorf("Expected order[%d] = %s, got %s", i, v, order[i])
		}
	}
}
