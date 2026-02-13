package singleflight

import (
	"context"
	"sync"
)

// SingleFlight lets concurrent calls with the same key share the call result.
// This prevents cache stampede and reduces load on downstream services.
//
// Example:
//
//	sf := singleflight.New()
//	result, err := sf.Do("user:123", func() (any, error) {
//	    return db.GetUser(123)
//	})
type SingleFlight interface {
	// Do executes the function for the given key, sharing the result with concurrent callers
	Do(key string, fn func() (any, error)) (any, error)

	// DoEx is like Do but also returns whether the result was freshly computed
	DoEx(key string, fn func() (any, error)) (val any, fresh bool, err error)

	// DoCtx is like Do but respects context cancellation
	DoCtx(ctx context.Context, key string, fn func() (any, error)) (any, error)

	// Forget removes a key from the group, allowing the next call to execute
	Forget(key string)
}

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type flightGroup struct {
	mu    sync.Mutex
	calls map[string]*call
}

// New creates a new SingleFlight instance
func New() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
}

func (g *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	c, done := g.createCall(key)
	if done {
		return c.val, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *flightGroup) DoEx(key string, fn func() (any, error)) (val any, fresh bool, err error) {
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, true, c.err
}

func (g *flightGroup) DoCtx(ctx context.Context, key string, fn func() (any, error)) (any, error) {
	c, done := g.createCall(key)
	if done {
		// Wait for existing call or context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return c.val, c.err
		}
	}

	// Execute the call
	done2 := make(chan struct{})
	go func() {
		g.makeCall(c, key, fn)
		close(done2)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done2:
		return c.val, c.err
	}
}

func (g *flightGroup) Forget(key string) {
	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()
}

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.mu.Lock()
	if c, ok := g.calls[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.mu.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		g.mu.Lock()
		delete(g.calls, key)
		g.mu.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}

// Typed provides type-safe singleflight operations
type Typed[T any] struct {
	sf SingleFlight
}

// NewTyped creates a new type-safe SingleFlight instance
func NewTyped[T any]() *Typed[T] {
	return &Typed[T]{sf: New()}
}

// Do executes the function for the given key with type safety
func (t *Typed[T]) Do(key string, fn func() (T, error)) (T, error) {
	val, err := t.sf.Do(key, func() (any, error) {
		return fn()
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return val.(T), nil
}

// DoEx is like Do but also returns whether the result was freshly computed
func (t *Typed[T]) DoEx(key string, fn func() (T, error)) (val T, fresh bool, err error) {
	v, fresh, err := t.sf.DoEx(key, func() (any, error) {
		return fn()
	})
	if err != nil {
		var zero T
		return zero, fresh, err
	}
	return v.(T), fresh, nil
}

// DoCtx is like Do but respects context cancellation
func (t *Typed[T]) DoCtx(ctx context.Context, key string, fn func() (T, error)) (T, error) {
	val, err := t.sf.DoCtx(ctx, key, func() (any, error) {
		return fn()
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return val.(T), nil
}

// Forget removes a key from the group
func (t *Typed[T]) Forget(key string) {
	t.sf.Forget(key)
}
