package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// testEvent is a simple event for testing
type testEvent struct {
	BaseEvent
	name    string
	payload string
}

func (e testEvent) EventName() string {
	return e.name
}

func newTestEvent(name, payload string) testEvent {
	return testEvent{
		BaseEvent: NewBaseEvent(),
		name:      name,
		payload:   payload,
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	bus := NewEventBus()

	called := false
	sub := bus.Subscribe("user.created", func(ctx context.Context, event Event) error {
		called = true
		return nil
	})

	if sub == nil {
		t.Fatal("Subscribe should return a subscription")
	}

	if sub.EventName() != "user.created" {
		t.Errorf("EventName() = %q, want %q", sub.EventName(), "user.created")
	}

	if sub.ID() == "" {
		t.Error("Subscription ID should not be empty")
	}

	// Publish event
	err := bus.Publish(context.Background(), newTestEvent("user.created", "test"))
	if err != nil {
		t.Errorf("Publish error: %v", err)
	}

	if !called {
		t.Error("Handler should have been called")
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := NewEventBus()

	callCount := 0
	sub := bus.Subscribe("user.created", func(ctx context.Context, event Event) error {
		callCount++
		return nil
	})

	// First publish
	_ = bus.Publish(context.Background(), newTestEvent("user.created", "test"))
	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	// Unsubscribe
	sub.Unsubscribe()

	// Second publish - handler should not be called
	_ = bus.Publish(context.Background(), newTestEvent("user.created", "test"))
	if callCount != 1 {
		t.Errorf("callCount = %d after unsubscribe, want 1", callCount)
	}
}

func TestEventBus_Priority(t *testing.T) {
	bus := NewEventBus()

	var order []int
	var mu sync.Mutex

	// Subscribe with different priorities
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		mu.Lock()
		order = append(order, 1)
		mu.Unlock()
		return nil
	}, WithPriority(10))

	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
		return nil
	}, WithPriority(100)) // Higher priority

	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		mu.Lock()
		order = append(order, 3)
		mu.Unlock()
		return nil
	}, WithPriority(50))

	_ = bus.Publish(context.Background(), newTestEvent("test.event", "test"))

	// Expected order: 2 (100), 3 (50), 1 (10)
	expected := []int{2, 3, 1}
	if len(order) != len(expected) {
		t.Fatalf("order length = %d, want %d", len(order), len(expected))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %d, want %d", i, order[i], v)
		}
	}
}

func TestEventBus_WildcardSubscription(t *testing.T) {
	bus := NewEventBus()

	var received []string
	var mu sync.Mutex

	bus.Subscribe("user.*", func(ctx context.Context, event Event) error {
		mu.Lock()
		received = append(received, event.EventName())
		mu.Unlock()
		return nil
	})

	// These should match
	_ = bus.Publish(context.Background(), newTestEvent("user.created", ""))
	_ = bus.Publish(context.Background(), newTestEvent("user.deleted", ""))

	// This should not match
	_ = bus.Publish(context.Background(), newTestEvent("order.created", ""))

	if len(received) != 2 {
		t.Errorf("received %d events, want 2", len(received))
	}
}

func TestEventBus_ErrorPropagation(t *testing.T) {
	bus := NewEventBus()

	expectedErr := errors.New("handler error")
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		return expectedErr
	})

	err := bus.Publish(context.Background(), newTestEvent("test.event", ""))
	if !errors.Is(err, expectedErr) {
		t.Errorf("Publish error = %v, want %v", err, expectedErr)
	}
}

func TestEventBus_ContextCancellation(t *testing.T) {
	bus := NewEventBus()

	ctx, cancel := context.WithCancel(context.Background())

	// First handler cancels context
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		cancel()
		return nil
	}, WithPriority(100))

	// Second handler should not be called
	secondCalled := false
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		secondCalled = true
		return nil
	}, WithPriority(10))

	err := bus.Publish(ctx, newTestEvent("test.event", ""))
	if err != context.Canceled {
		t.Errorf("Publish error = %v, want context.Canceled", err)
	}

	if secondCalled {
		t.Error("Second handler should not have been called after context cancellation")
	}
}

func TestEventBus_AsyncHandler(t *testing.T) {
	bus := NewEventBus()

	var called int32
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		atomic.AddInt32(&called, 1)
		wg.Done()
		return nil
	}, WithAsync())

	err := bus.Publish(context.Background(), newTestEvent("test.event", ""))
	if err != nil {
		t.Errorf("Publish error: %v", err)
	}

	// Wait for async handler
	wg.Wait()

	if atomic.LoadInt32(&called) != 1 {
		t.Error("Async handler should have been called")
	}
}

func TestEventBus_Middleware(t *testing.T) {
	bus := NewEventBus()

	var order []string
	var mu sync.Mutex

	// Add middleware
	bus.Use(func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			mu.Lock()
			order = append(order, "middleware1-before")
			mu.Unlock()
			err := next(ctx, event)
			mu.Lock()
			order = append(order, "middleware1-after")
			mu.Unlock()
			return err
		}
	})

	bus.Use(func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			mu.Lock()
			order = append(order, "middleware2-before")
			mu.Unlock()
			err := next(ctx, event)
			mu.Lock()
			order = append(order, "middleware2-after")
			mu.Unlock()
			return err
		}
	})

	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		mu.Lock()
		order = append(order, "handler")
		mu.Unlock()
		return nil
	})

	_ = bus.Publish(context.Background(), newTestEvent("test.event", ""))

	expected := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(order) != len(expected) {
		t.Fatalf("order length = %d, want %d", len(order), len(expected))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

func TestEventBus_ClosedBus(t *testing.T) {
	bus := NewEventBus()
	bus.Close()

	err := bus.Publish(context.Background(), newTestEvent("test.event", ""))
	if !errors.Is(err, ErrEventBusClosed) {
		t.Errorf("Publish error = %v, want ErrEventBusClosed", err)
	}
}

func TestEventBus_HasSubscribers(t *testing.T) {
	bus := NewEventBus()

	if bus.HasSubscribers("user.created") {
		t.Error("HasSubscribers should return false for no subscribers")
	}

	bus.Subscribe("user.*", func(ctx context.Context, event Event) error {
		return nil
	})

	if !bus.HasSubscribers("user.created") {
		t.Error("HasSubscribers should return true for matching wildcard")
	}

	if bus.HasSubscribers("order.created") {
		t.Error("HasSubscribers should return false for non-matching pattern")
	}
}

func TestBaseEvent_Metadata(t *testing.T) {
	event := NewBaseEvent()

	meta := event.Metadata()
	if meta.ID == "" {
		t.Error("Event ID should not be empty")
	}

	if meta.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	// Test with source
	eventWithSource := NewBaseEventWithSource("user-service")
	if eventWithSource.Metadata().Source != "user-service" {
		t.Errorf("Source = %q, want %q", eventWithSource.Metadata().Source, "user-service")
	}

	// Test with correlation
	eventWithCorr := NewBaseEventWithCorrelation("corr-123", "cause-456", "order-service")
	meta = eventWithCorr.Metadata()
	if meta.CorrelationID != "corr-123" {
		t.Errorf("CorrelationID = %q, want %q", meta.CorrelationID, "corr-123")
	}
	if meta.CausationID != "cause-456" {
		t.Errorf("CausationID = %q, want %q", meta.CausationID, "cause-456")
	}
	if meta.Source != "order-service" {
		t.Errorf("Source = %q, want %q", meta.Source, "order-service")
	}
}

func TestEventBus_PublishAsync(t *testing.T) {
	bus := NewEventBus()

	var called int32
	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	bus.PublishAsync(context.Background(), newTestEvent("test.event", ""))

	// Give async handler time to execute
	time.Sleep(50 * time.Millisecond)

	if atomic.LoadInt32(&called) != 1 {
		t.Error("Handler should have been called asynchronously")
	}
}

func TestEventBus_Clear(t *testing.T) {
	bus := NewEventBus()

	bus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		return nil
	})

	if bus.HandlerCount() != 1 {
		t.Errorf("HandlerCount = %d, want 1", bus.HandlerCount())
	}

	bus.Clear()

	if bus.HandlerCount() != 0 {
		t.Errorf("HandlerCount after Clear = %d, want 0", bus.HandlerCount())
	}
}
