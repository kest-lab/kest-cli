package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = errors.New("cache: key not found")

// Store defines the cache store interface
type Store interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (interface{}, error)

	// Put stores a value in the cache with expiration
	Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Forever stores a value in the cache indefinitely
	Forever(ctx context.Context, key string, value interface{}) error

	// Forget removes a value from the cache
	Forget(ctx context.Context, key string) error

	// Flush removes all values from the cache
	Flush(ctx context.Context) error

	// Has checks if a key exists in the cache
	Has(ctx context.Context, key string) bool

	// Increment increments a numeric value
	Increment(ctx context.Context, key string, value int64) (int64, error)

	// Decrement decrements a numeric value
	Decrement(ctx context.Context, key string, value int64) (int64, error)
}

// Manager manages multiple cache stores
type Manager struct {
	mu       sync.RWMutex
	stores   map[string]Store
	default_ string
}

var (
	manager *Manager
	once    sync.Once
)

// Global returns the global cache manager
func Global() *Manager {
	once.Do(func() {
		manager = &Manager{
			stores:   make(map[string]Store),
			default_: "memory",
		}
		// Register default memory store
		manager.stores["memory"] = NewMemoryStore()
	})
	return manager
}

// SetDefault sets the default store name
func (m *Manager) SetDefault(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.default_ = name
}

// Register registers a cache store
func (m *Manager) Register(name string, store Store) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[name] = store
}

// Store returns a cache store by name
func (m *Manager) Store(name string) Store {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if store, ok := m.stores[name]; ok {
		return store
	}
	return m.stores[m.default_]
}

// Default returns the default cache store
func (m *Manager) Default() Store {
	return m.Store(m.default_)
}

// --- Convenience functions using default store ---

// Get retrieves a value from the default cache
func Get(ctx context.Context, key string) (interface{}, error) {
	return Global().Default().Get(ctx, key)
}

// GetString retrieves a string value from the cache
func GetString(ctx context.Context, key string) (string, error) {
	val, err := Get(ctx, key)
	if err != nil {
		return "", err
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", errors.New("cache: value is not a string")
}

// GetInt retrieves an int value from the cache
func GetInt(ctx context.Context, key string) (int, error) {
	val, err := Get(ctx, key)
	if err != nil {
		return 0, err
	}
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	}
	return 0, errors.New("cache: value is not an int")
}

// Put stores a value in the default cache
func Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return Global().Default().Put(ctx, key, value, ttl)
}

// Forever stores a value indefinitely in the default cache
func Forever(ctx context.Context, key string, value interface{}) error {
	return Global().Default().Forever(ctx, key, value)
}

// Forget removes a value from the default cache
func Forget(ctx context.Context, key string) error {
	return Global().Default().Forget(ctx, key)
}

// Flush removes all values from the default cache
func Flush(ctx context.Context) error {
	return Global().Default().Flush(ctx)
}

// Has checks if a key exists in the default cache
func Has(ctx context.Context, key string) bool {
	return Global().Default().Has(ctx, key)
}

// Remember gets a value from cache or stores the result of callback
func Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return RememberStore(ctx, Global().Default(), key, ttl, callback)
}

// RememberStore gets a value from a specific store or stores the result of callback
func RememberStore(ctx context.Context, store Store, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if val, err := store.Get(ctx, key); err == nil {
		return val, nil
	}

	// Execute callback
	val, err := callback()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if ttl > 0 {
		_ = store.Put(ctx, key, val, ttl)
	} else {
		_ = store.Forever(ctx, key, val)
	}

	return val, nil
}

// RememberForever gets a value from cache or stores the result indefinitely
func RememberForever(ctx context.Context, key string, callback func() (interface{}, error)) (interface{}, error) {
	return Remember(ctx, key, 0, callback)
}

// Pull retrieves a value and removes it from cache
func Pull(ctx context.Context, key string) (interface{}, error) {
	val, err := Get(ctx, key)
	if err != nil {
		return nil, err
	}
	_ = Forget(ctx, key)
	return val, nil
}

// Add stores a value only if the key doesn't exist
func Add(ctx context.Context, key string, value interface{}, ttl time.Duration) bool {
	if Has(ctx, key) {
		return false
	}
	return Put(ctx, key, value, ttl) == nil
}

// Increment increments a numeric value
func Increment(ctx context.Context, key string, value int64) (int64, error) {
	return Global().Default().Increment(ctx, key, value)
}

// Decrement decrements a numeric value
func Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return Global().Default().Decrement(ctx, key, value)
}
