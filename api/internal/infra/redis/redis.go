package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client with a Laravel-style API
type Client struct {
	rdb *redis.Client
	ctx context.Context
}

// Config holds Redis connection configuration
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
	Prefix   string
}

// DefaultConfig returns default Redis configuration
func DefaultConfig() Config {
	return Config{
		Host:   "localhost",
		Port:   "6379",
		DB:     0,
		Prefix: "",
	}
}

var (
	defaultClient *Client
	defaultPrefix string
)

// Connect creates a new Redis connection
func Connect(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	client := &Client{
		rdb: rdb,
		ctx: ctx,
	}

	defaultPrefix = cfg.Prefix
	defaultClient = client

	return client, nil
}

// Default returns the default client
func Default() *Client {
	return defaultClient
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}

// --- Key Operations ---

// key adds prefix to key
func (c *Client) key(k string) string {
	if defaultPrefix != "" {
		return defaultPrefix + k
	}
	return k
}

// Set stores a value with optional expiration
func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(c.ctx, c.key(key), value, expiration).Err()
}

// Get retrieves a value
func (c *Client) Get(key string) (string, error) {
	return c.rdb.Get(c.ctx, c.key(key)).Result()
}

// GetBytes retrieves a value as bytes
func (c *Client) GetBytes(key string) ([]byte, error) {
	return c.rdb.Get(c.ctx, c.key(key)).Bytes()
}

// Del deletes one or more keys
func (c *Client) Del(keys ...string) (int64, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.key(k)
	}
	return c.rdb.Del(c.ctx, prefixedKeys...).Result()
}

// Exists checks if keys exist
func (c *Client) Exists(keys ...string) (int64, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.key(k)
	}
	return c.rdb.Exists(c.ctx, prefixedKeys...).Result()
}

// Expire sets a key's TTL
func (c *Client) Expire(key string, expiration time.Duration) (bool, error) {
	return c.rdb.Expire(c.ctx, c.key(key), expiration).Result()
}

// TTL returns a key's TTL
func (c *Client) TTL(key string) (time.Duration, error) {
	return c.rdb.TTL(c.ctx, c.key(key)).Result()
}

// Keys returns keys matching pattern
func (c *Client) Keys(pattern string) ([]string, error) {
	return c.rdb.Keys(c.ctx, c.key(pattern)).Result()
}

// --- String Operations ---

// Incr increments a key
func (c *Client) Incr(key string) (int64, error) {
	return c.rdb.Incr(c.ctx, c.key(key)).Result()
}

// IncrBy increments a key by value
func (c *Client) IncrBy(key string, value int64) (int64, error) {
	return c.rdb.IncrBy(c.ctx, c.key(key), value).Result()
}

// Decr decrements a key
func (c *Client) Decr(key string) (int64, error) {
	return c.rdb.Decr(c.ctx, c.key(key)).Result()
}

// DecrBy decrements a key by value
func (c *Client) DecrBy(key string, value int64) (int64, error) {
	return c.rdb.DecrBy(c.ctx, c.key(key), value).Result()
}

// SetNX sets a key only if it doesn't exist
func (c *Client) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(c.ctx, c.key(key), value, expiration).Result()
}

// --- Hash Operations ---

// HSet sets hash fields
func (c *Client) HSet(key string, values ...interface{}) (int64, error) {
	return c.rdb.HSet(c.ctx, c.key(key), values...).Result()
}

// HGet gets a hash field
func (c *Client) HGet(key, field string) (string, error) {
	return c.rdb.HGet(c.ctx, c.key(key), field).Result()
}

// HGetAll gets all hash fields
func (c *Client) HGetAll(key string) (map[string]string, error) {
	return c.rdb.HGetAll(c.ctx, c.key(key)).Result()
}

// HDel deletes hash fields
func (c *Client) HDel(key string, fields ...string) (int64, error) {
	return c.rdb.HDel(c.ctx, c.key(key), fields...).Result()
}

// HExists checks if hash field exists
func (c *Client) HExists(key, field string) (bool, error) {
	return c.rdb.HExists(c.ctx, c.key(key), field).Result()
}

// --- List Operations ---

// LPush prepends values to a list
func (c *Client) LPush(key string, values ...interface{}) (int64, error) {
	return c.rdb.LPush(c.ctx, c.key(key), values...).Result()
}

// RPush appends values to a list
func (c *Client) RPush(key string, values ...interface{}) (int64, error) {
	return c.rdb.RPush(c.ctx, c.key(key), values...).Result()
}

// LPop removes and returns the first element
func (c *Client) LPop(key string) (string, error) {
	return c.rdb.LPop(c.ctx, c.key(key)).Result()
}

// RPop removes and returns the last element
func (c *Client) RPop(key string) (string, error) {
	return c.rdb.RPop(c.ctx, c.key(key)).Result()
}

// LRange returns a range of elements
func (c *Client) LRange(key string, start, stop int64) ([]string, error) {
	return c.rdb.LRange(c.ctx, c.key(key), start, stop).Result()
}

// LLen returns the length of a list
func (c *Client) LLen(key string) (int64, error) {
	return c.rdb.LLen(c.ctx, c.key(key)).Result()
}

// --- Set Operations ---

// SAdd adds members to a set
func (c *Client) SAdd(key string, members ...interface{}) (int64, error) {
	return c.rdb.SAdd(c.ctx, c.key(key), members...).Result()
}

// SMembers returns all members
func (c *Client) SMembers(key string) ([]string, error) {
	return c.rdb.SMembers(c.ctx, c.key(key)).Result()
}

// SIsMember checks if member exists
func (c *Client) SIsMember(key string, member interface{}) (bool, error) {
	return c.rdb.SIsMember(c.ctx, c.key(key), member).Result()
}

// SRem removes members
func (c *Client) SRem(key string, members ...interface{}) (int64, error) {
	return c.rdb.SRem(c.ctx, c.key(key), members...).Result()
}

// SCard returns the set size
func (c *Client) SCard(key string) (int64, error) {
	return c.rdb.SCard(c.ctx, c.key(key)).Result()
}

// --- Sorted Set Operations ---

// ZAdd adds members with scores
func (c *Client) ZAdd(key string, members ...redis.Z) (int64, error) {
	return c.rdb.ZAdd(c.ctx, c.key(key), members...).Result()
}

// ZRange returns members by rank
func (c *Client) ZRange(key string, start, stop int64) ([]string, error) {
	return c.rdb.ZRange(c.ctx, c.key(key), start, stop).Result()
}

// ZRangeWithScores returns members with scores by rank
func (c *Client) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return c.rdb.ZRangeWithScores(c.ctx, c.key(key), start, stop).Result()
}

// ZRem removes members
func (c *Client) ZRem(key string, members ...interface{}) (int64, error) {
	return c.rdb.ZRem(c.ctx, c.key(key), members...).Result()
}

// ZCard returns the sorted set size
func (c *Client) ZCard(key string) (int64, error) {
	return c.rdb.ZCard(c.ctx, c.key(key)).Result()
}

// ZScore returns a member's score
func (c *Client) ZScore(key string, member string) (float64, error) {
	return c.rdb.ZScore(c.ctx, c.key(key), member).Result()
}

// --- Pub/Sub ---

// Publish publishes a message to a channel
func (c *Client) Publish(channel string, message interface{}) (int64, error) {
	return c.rdb.Publish(c.ctx, channel, message).Result()
}

// Subscribe subscribes to channels
func (c *Client) Subscribe(channels ...string) *redis.PubSub {
	return c.rdb.Subscribe(c.ctx, channels...)
}

// --- Pipeline ---

// Pipeline returns a pipeline
func (c *Client) Pipeline() redis.Pipeliner {
	return c.rdb.Pipeline()
}

// Pipelined executes commands in a pipeline
func (c *Client) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.rdb.Pipelined(c.ctx, fn)
}

// TxPipeline returns a transactional pipeline
func (c *Client) TxPipeline() redis.Pipeliner {
	return c.rdb.TxPipeline()
}

// TxPipelined executes commands in a transaction
func (c *Client) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return c.rdb.TxPipelined(c.ctx, fn)
}

// --- Utility ---

// FlushDB flushes the current database
func (c *Client) FlushDB() error {
	return c.rdb.FlushDB(c.ctx).Err()
}

// Ping pings the Redis server
func (c *Client) Ping() error {
	return c.rdb.Ping(c.ctx).Err()
}

// Raw returns the underlying redis client
func (c *Client) Raw() *redis.Client {
	return c.rdb
}

// --- Package-level convenience functions ---

// Set stores a value (uses default client)
func Set(key string, value interface{}, expiration time.Duration) error {
	return defaultClient.Set(key, value, expiration)
}

// Get retrieves a value (uses default client)
func Get(key string) (string, error) {
	return defaultClient.Get(key)
}

// Del deletes keys (uses default client)
func Del(keys ...string) (int64, error) {
	return defaultClient.Del(keys...)
}

// Exists checks if keys exist (uses default client)
func Exists(keys ...string) (int64, error) {
	return defaultClient.Exists(keys...)
}
