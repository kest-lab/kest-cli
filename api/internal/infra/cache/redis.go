package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements a Redis-backed cache store
type RedisStore struct {
	client *redis.Client
	prefix string
}

// RedisOption configures the Redis store
type RedisOption func(*RedisStore)

// WithPrefix sets a key prefix for all cache keys
func WithPrefix(prefix string) RedisOption {
	return func(s *RedisStore) {
		s.prefix = prefix
	}
}

// NewRedisStore creates a new Redis cache store
func NewRedisStore(client *redis.Client, opts ...RedisOption) *RedisStore {
	s := &RedisStore{
		client: client,
		prefix: "cache:",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// NewRedisStoreFromConfig creates a Redis store from connection config
func NewRedisStoreFromConfig(addr, password string, db int, opts ...RedisOption) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return NewRedisStore(client, opts...)
}

// prefixKey adds the prefix to a key
func (s *RedisStore) prefixKey(key string) string {
	return s.prefix + key
}

// Get retrieves a value from the cache
func (s *RedisStore) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := s.client.Get(ctx, s.prefixKey(key)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	// Try to unmarshal as JSON
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// Return raw string if not JSON
		return val, nil
	}

	return result, nil
}

// Put stores a value in the cache with expiration
func (s *RedisStore) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var val string

	switch v := value.(type) {
	case string:
		val = v
	case []byte:
		val = string(v)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		val = string(data)
	}

	return s.client.Set(ctx, s.prefixKey(key), val, ttl).Err()
}

// Forever stores a value in the cache indefinitely
func (s *RedisStore) Forever(ctx context.Context, key string, value interface{}) error {
	return s.Put(ctx, key, value, 0)
}

// Forget removes a value from the cache
func (s *RedisStore) Forget(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.prefixKey(key)).Err()
}

// Flush removes all values with the prefix from the cache
func (s *RedisStore) Flush(ctx context.Context) error {
	iter := s.client.Scan(ctx, 0, s.prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Has checks if a key exists in the cache
func (s *RedisStore) Has(ctx context.Context, key string) bool {
	exists, err := s.client.Exists(ctx, s.prefixKey(key)).Result()
	return err == nil && exists > 0
}

// Increment increments a numeric value
func (s *RedisStore) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return s.client.IncrBy(ctx, s.prefixKey(key), value).Result()
}

// Decrement decrements a numeric value
func (s *RedisStore) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return s.client.DecrBy(ctx, s.prefixKey(key), value).Result()
}

// GetClient returns the underlying Redis client
func (s *RedisStore) GetClient() *redis.Client {
	return s.client
}

// Close closes the Redis connection
func (s *RedisStore) Close() error {
	return s.client.Close()
}

// TTL returns the remaining time to live of a key
func (s *RedisStore) TTL(ctx context.Context, key string) (time.Duration, error) {
	return s.client.TTL(ctx, s.prefixKey(key)).Result()
}

// Expire sets a new expiration time on a key
func (s *RedisStore) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return s.client.Expire(ctx, s.prefixKey(key), ttl).Err()
}
