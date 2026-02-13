package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

// ErrMaxRetriesExceeded is returned when max retries are exceeded
var ErrMaxRetriesExceeded = errors.New("max retries exceeded")

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// RetryableFuncWithResult is a function that returns a result and can be retried
type RetryableFuncWithResult[T any] func(ctx context.Context) (T, error)

// ShouldRetry determines if an error should trigger a retry
type ShouldRetry func(err error) bool

// Config holds retry configuration
type Config struct {
	// MaxAttempts is the maximum number of attempts (including the first one)
	MaxAttempts int

	// InitialDelay is the initial delay before the first retry
	InitialDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// Multiplier is the factor by which the delay increases
	Multiplier float64

	// Jitter adds randomness to the delay (0.0 to 1.0)
	Jitter float64

	// ShouldRetry determines if an error should trigger a retry
	ShouldRetry ShouldRetry
}

// DefaultConfig returns default retry configuration
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       0.1,
		ShouldRetry:  func(err error) bool { return err != nil },
	}
}

// Option is a function that modifies Config
type Option func(*Config)

// WithMaxAttempts sets the maximum number of attempts
func WithMaxAttempts(n int) Option {
	return func(c *Config) {
		c.MaxAttempts = n
	}
}

// WithInitialDelay sets the initial delay
func WithInitialDelay(d time.Duration) Option {
	return func(c *Config) {
		c.InitialDelay = d
	}
}

// WithMaxDelay sets the maximum delay
func WithMaxDelay(d time.Duration) Option {
	return func(c *Config) {
		c.MaxDelay = d
	}
}

// WithMultiplier sets the backoff multiplier
func WithMultiplier(m float64) Option {
	return func(c *Config) {
		c.Multiplier = m
	}
}

// WithJitter sets the jitter factor
func WithJitter(j float64) Option {
	return func(c *Config) {
		c.Jitter = j
	}
}

// WithShouldRetry sets the retry condition function
func WithShouldRetry(fn ShouldRetry) Option {
	return func(c *Config) {
		c.ShouldRetry = fn
	}
}

// Do executes a function with retry logic
func Do(ctx context.Context, fn RetryableFunc, opts ...Option) error {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		// Check context before attempting
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute the function
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry
		if !cfg.ShouldRetry(err) {
			return err
		}

		// Don't sleep after the last attempt
		if attempt < cfg.MaxAttempts-1 {
			delay := calculateDelay(attempt, cfg)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return errors.Join(ErrMaxRetriesExceeded, lastErr)
}

// DoWithResult executes a function that returns a result with retry logic
func DoWithResult[T any](ctx context.Context, fn RetryableFuncWithResult[T], opts ...Option) (T, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	var zero T
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		default:
		}

		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}

		lastErr = err

		if !cfg.ShouldRetry(err) {
			return zero, err
		}

		if attempt < cfg.MaxAttempts-1 {
			delay := calculateDelay(attempt, cfg)
			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return zero, errors.Join(ErrMaxRetriesExceeded, lastErr)
}

func calculateDelay(attempt int, cfg Config) time.Duration {
	// Calculate exponential backoff
	delay := float64(cfg.InitialDelay) * math.Pow(cfg.Multiplier, float64(attempt))

	// Apply max delay cap
	if delay > float64(cfg.MaxDelay) {
		delay = float64(cfg.MaxDelay)
	}

	// Apply jitter
	if cfg.Jitter > 0 {
		jitter := delay * cfg.Jitter * (rand.Float64()*2 - 1) // -jitter to +jitter
		delay += jitter
	}

	return time.Duration(delay)
}

// Retrier provides a reusable retry configuration
type Retrier struct {
	cfg Config
}

// NewRetrier creates a new Retrier with the given options
func NewRetrier(opts ...Option) *Retrier {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Retrier{cfg: cfg}
}

// Do executes a function with the retrier's configuration
func (r *Retrier) Do(ctx context.Context, fn RetryableFunc) error {
	return Do(ctx, fn,
		WithMaxAttempts(r.cfg.MaxAttempts),
		WithInitialDelay(r.cfg.InitialDelay),
		WithMaxDelay(r.cfg.MaxDelay),
		WithMultiplier(r.cfg.Multiplier),
		WithJitter(r.cfg.Jitter),
		WithShouldRetry(r.cfg.ShouldRetry),
	)
}

// DoWithResult executes a function that returns a result
func (r *Retrier) DoWithResult(ctx context.Context, fn RetryableFuncWithResult[any]) (any, error) {
	return DoWithResult(ctx, fn,
		WithMaxAttempts(r.cfg.MaxAttempts),
		WithInitialDelay(r.cfg.InitialDelay),
		WithMaxDelay(r.cfg.MaxDelay),
		WithMultiplier(r.cfg.Multiplier),
		WithJitter(r.cfg.Jitter),
		WithShouldRetry(r.cfg.ShouldRetry),
	)
}

// --- Convenience Functions ---

// Times retries a function n times with no delay
func Times(ctx context.Context, n int, fn RetryableFunc) error {
	return Do(ctx, fn, WithMaxAttempts(n), WithInitialDelay(0), WithMaxDelay(0))
}

// Forever retries a function indefinitely until success or context cancellation
func Forever(ctx context.Context, fn RetryableFunc, delay time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

// ExponentialBackoff retries with exponential backoff
func ExponentialBackoff(ctx context.Context, fn RetryableFunc, maxAttempts int) error {
	return Do(ctx, fn,
		WithMaxAttempts(maxAttempts),
		WithInitialDelay(100*time.Millisecond),
		WithMaxDelay(30*time.Second),
		WithMultiplier(2.0),
		WithJitter(0.1),
	)
}
