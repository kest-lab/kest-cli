package cache

import (
	"context"
	"sync"
	"time"
)

// item represents a cached item with expiration
type item struct {
	value      interface{}
	expiration int64 // Unix nano timestamp, 0 means no expiration
}

// isExpired checks if the item has expired
func (i *item) isExpired() bool {
	if i.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.expiration
}

// MemoryStore implements an in-memory cache store
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]*item

	// Cleanup settings
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// MemoryOption configures the memory store
type MemoryOption func(*MemoryStore)

// WithCleanupInterval sets the cleanup interval for expired items
func WithCleanupInterval(d time.Duration) MemoryOption {
	return func(s *MemoryStore) {
		s.cleanupInterval = d
	}
}

// NewMemoryStore creates a new in-memory cache store
func NewMemoryStore(opts ...MemoryOption) *MemoryStore {
	s := &MemoryStore{
		items:           make(map[string]*item),
		cleanupInterval: 5 * time.Minute,
		stopCleanup:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(s)
	}

	// Start cleanup goroutine
	go s.cleanup()

	return s
}

// cleanup periodically removes expired items
func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.deleteExpired()
		case <-s.stopCleanup:
			return
		}
	}
}

// deleteExpired removes all expired items
func (s *MemoryStore) deleteExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range s.items {
		if item.expiration > 0 && now > item.expiration {
			delete(s.items, key)
		}
	}
}

// Close stops the cleanup goroutine
func (s *MemoryStore) Close() {
	close(s.stopCleanup)
}

// Get retrieves a value from the cache
func (s *MemoryStore) Get(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	item, ok := s.items[key]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrCacheMiss
	}

	if item.isExpired() {
		s.Forget(ctx, key)
		return nil, ErrCacheMiss
	}

	return item.value, nil
}

// Put stores a value in the cache with expiration
func (s *MemoryStore) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	s.mu.Lock()
	s.items[key] = &item{
		value:      value,
		expiration: expiration,
	}
	s.mu.Unlock()

	return nil
}

// Forever stores a value in the cache indefinitely
func (s *MemoryStore) Forever(ctx context.Context, key string, value interface{}) error {
	return s.Put(ctx, key, value, 0)
}

// Forget removes a value from the cache
func (s *MemoryStore) Forget(ctx context.Context, key string) error {
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
	return nil
}

// Flush removes all values from the cache
func (s *MemoryStore) Flush(ctx context.Context) error {
	s.mu.Lock()
	s.items = make(map[string]*item)
	s.mu.Unlock()
	return nil
}

// Has checks if a key exists in the cache
func (s *MemoryStore) Has(ctx context.Context, key string) bool {
	s.mu.RLock()
	item, ok := s.items[key]
	s.mu.RUnlock()

	if !ok {
		return false
	}

	if item.isExpired() {
		s.Forget(ctx, key)
		return false
	}

	return true
}

// Increment increments a numeric value
func (s *MemoryStore) Increment(ctx context.Context, key string, value int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itm, ok := s.items[key]
	if !ok || itm.isExpired() {
		s.items[key] = &item{value: value, expiration: 0}
		return value, nil
	}

	var current int64
	switch v := itm.value.(type) {
	case int:
		current = int64(v)
	case int64:
		current = v
	case float64:
		current = int64(v)
	default:
		return 0, ErrCacheMiss
	}

	newValue := current + value
	itm.value = newValue
	return newValue, nil
}

// Decrement decrements a numeric value
func (s *MemoryStore) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return s.Increment(ctx, key, -value)
}

// Keys returns all keys in the cache (for debugging)
func (s *MemoryStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.items))
	for k, item := range s.items {
		if !item.isExpired() {
			keys = append(keys, k)
		}
	}
	return keys
}

// Len returns the number of items in the cache
func (s *MemoryStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
