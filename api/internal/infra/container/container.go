package container

import (
	"fmt"
	"sync"
)

// Container provides a global service registry for dependency management.
type Container struct {
	mu        sync.RWMutex
	entries   map[string]any
	factories map[string]func() (any, error)
}

var (
	global     = &Container{entries: make(map[string]any), factories: make(map[string]func() (any, error))}
	onceGlobal sync.Once
)

// App returns the global container instance.
func App() *Container {
	onceGlobal.Do(func() {
		if global.entries == nil {
			global.entries = make(map[string]any)
		}
		if global.factories == nil {
			global.factories = make(map[string]func() (any, error))
		}
	})
	return global
}

// Reset clears all registered services and factories. Intended for tests.
func Reset() {
	c := App()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]any)
	c.factories = make(map[string]func() (any, error))
}

// Set registers a concrete singleton value under the provided key.
func (c *Container) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.factories, key)
	c.entries[key] = value
}

// Bind registers a lazy factory for the given key. The factory will be invoked once on first resolve.
func (c *Container) Bind(key string, factory func() (any, error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
	c.factories[key] = factory
}

// Has returns true if the key is registered.
func (c *Container) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.entries[key]; ok {
		return true
	}
	_, ok := c.factories[key]
	return ok
}

// Resolve retrieves a service by key.
func (c *Container) Resolve(key string) (any, error) {
	c.mu.RLock()
	value, ok := c.entries[key]
	factory, factoryOK := c.factories[key]
	c.mu.RUnlock()

	switch {
	case ok:
		return value, nil
	case factoryOK:
		c.mu.Lock()
		defer c.mu.Unlock()
		// Check again after acquiring write lock in case another goroutine already resolved it.
		if value, ok := c.entries[key]; ok {
			return value, nil
		}
		instance, err := factory()
		if err != nil {
			return nil, err
		}
		c.entries[key] = instance
		delete(c.factories, key)
		return instance, nil
	default:
		return nil, fmt.Errorf("container: service '%s' not registered", key)
	}
}

// MustResolve retrieves a service or panics if not registered.
func (c *Container) MustResolve(key string) any {
	instance, err := c.Resolve(key)
	if err != nil {
		panic(err)
	}
	return instance
}

// ResolveAs retrieves a service and casts it to the requested type.
func ResolveAs[T any](key string) (T, error) {
	var zero T
	instance, err := App().Resolve(key)
	if err != nil {
		return zero, err
	}
	typed, ok := instance.(T)
	if !ok {
		return zero, fmt.Errorf("container: service '%s' has unexpected type", key)
	}
	return typed, nil
}

// MustResolveAs retrieves a typed service or panics if not available.
func MustResolveAs[T any](key string) T {
	val, err := ResolveAs[T](key)
	if err != nil {
		panic(err)
	}
	return val
}

// Keys returns a snapshot of registered service keys (including factories).
func (c *Container) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.entries)+len(c.factories))
	for k := range c.entries {
		keys = append(keys, k)
	}
	for k := range c.factories {
		keys = append(keys, k)
	}
	return keys
}
