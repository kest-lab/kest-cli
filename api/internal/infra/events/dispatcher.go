package events

import (
	"context"
	"reflect"
	"sync"
)

// SimpleEvent is a marker interface for simple events (legacy dispatcher)
type SimpleEvent interface {
	EventName() string
}

// SimpleListener handles a simple event
type SimpleListener interface {
	Handle(ctx context.Context, event SimpleEvent) error
}

// SimpleListenerFunc is a function adapter for SimpleListener
type SimpleListenerFunc func(ctx context.Context, event SimpleEvent) error

func (f SimpleListenerFunc) Handle(ctx context.Context, event SimpleEvent) error {
	return f(ctx, event)
}

// SimpleDispatcher manages simple event dispatching (legacy)
type SimpleDispatcher struct {
	mu        sync.RWMutex
	listeners map[string][]SimpleListener
	async     bool
}

// simpleDispatcher is the global dispatcher instance
var (
	globalSimpleDispatcher *SimpleDispatcher
	simpleOnce             sync.Once
)

// GlobalSimpleDispatcher returns the global simple dispatcher instance
func GlobalSimpleDispatcher() *SimpleDispatcher {
	simpleOnce.Do(func() {
		globalSimpleDispatcher = NewSimpleDispatcher()
	})
	return globalSimpleDispatcher
}

// NewSimpleDispatcher creates a new simple event dispatcher
func NewSimpleDispatcher() *SimpleDispatcher {
	return &SimpleDispatcher{
		listeners: make(map[string][]SimpleListener),
	}
}

// Listen registers a listener for an event
func (d *SimpleDispatcher) Listen(eventName string, listener SimpleListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

// ListenFunc registers a function listener for an event
func (d *SimpleDispatcher) ListenFunc(eventName string, fn func(ctx context.Context, event SimpleEvent) error) {
	d.Listen(eventName, SimpleListenerFunc(fn))
}

// Subscribe registers a listener for a typed event using reflection
func (d *SimpleDispatcher) Subscribe(eventType SimpleEvent, listener SimpleListener) {
	d.Listen(eventType.EventName(), listener)
}

// Dispatch fires an event to all registered listeners
func (d *SimpleDispatcher) Dispatch(ctx context.Context, event SimpleEvent) error {
	d.mu.RLock()
	listeners := d.listeners[event.EventName()]
	d.mu.RUnlock()

	for _, listener := range listeners {
		if err := listener.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// DispatchAsync fires an event asynchronously to all registered listeners
func (d *SimpleDispatcher) DispatchAsync(ctx context.Context, event SimpleEvent) {
	d.mu.RLock()
	listeners := d.listeners[event.EventName()]
	d.mu.RUnlock()

	for _, listener := range listeners {
		go func(l SimpleListener) {
			_ = l.Handle(ctx, event)
		}(listener)
	}
}

// DispatchAll fires multiple events
func (d *SimpleDispatcher) DispatchAll(ctx context.Context, events ...SimpleEvent) error {
	for _, event := range events {
		if err := d.Dispatch(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// HasListeners checks if an event has any listeners
func (d *SimpleDispatcher) HasListeners(eventName string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.listeners[eventName]) > 0
}

// Forget removes all listeners for an event
func (d *SimpleDispatcher) Forget(eventName string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.listeners, eventName)
}

// Flush removes all listeners
func (d *SimpleDispatcher) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners = make(map[string][]SimpleListener)
}

// GetSimpleEventName extracts event name from type (helper for struct events)
func GetSimpleEventName(event SimpleEvent) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}

// --- Convenience functions using global simple dispatcher ---

// ListenSimple registers a listener on the global simple dispatcher
func ListenSimple(eventName string, listener SimpleListener) {
	GlobalSimpleDispatcher().Listen(eventName, listener)
}

// ListenSimpleFunc registers a function listener on the global simple dispatcher
func ListenSimpleFunc(eventName string, fn func(ctx context.Context, event SimpleEvent) error) {
	GlobalSimpleDispatcher().ListenFunc(eventName, fn)
}

// DispatchSimple fires an event on the global simple dispatcher
func DispatchSimple(ctx context.Context, event SimpleEvent) error {
	return GlobalSimpleDispatcher().Dispatch(ctx, event)
}

// DispatchSimpleAsync fires an event asynchronously on the global simple dispatcher
func DispatchSimpleAsync(ctx context.Context, event SimpleEvent) {
	GlobalSimpleDispatcher().DispatchAsync(ctx, event)
}
