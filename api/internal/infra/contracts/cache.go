package contracts

import (
	"context"
	"time"
)

// Cache defines the caching interface
type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Has(ctx context.Context, key string) bool
	Flush(ctx context.Context) error
}

// CacheStore defines a cache store that can be tagged
type CacheStore interface {
	Cache
	Tags(tags ...string) Cache
	Forever(ctx context.Context, key string, value any) error
	Remember(ctx context.Context, key string, ttl time.Duration, callback func() (any, error)) (any, error)
}
