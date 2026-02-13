package capabilities

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// Circuit Breaker Pattern
// =============================================================================

// CircuitBreaker protects against cascading failures
type CircuitBreaker struct {
	maxFailures  int
	resetTimeout time.Duration

	mu              sync.RWMutex
	failures        int
	lastFailureTime time.Time
	state           string // "closed", "open", "half-open"
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "closed",
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

// Execute runs the function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// Check state
	cb.mu.RLock()
	state := cb.state
	lastFailure := cb.lastFailureTime
	cb.mu.RUnlock()

	// If open, check if we should try again
	if state == "open" {
		if time.Since(lastFailure) > cb.resetTimeout {
			cb.mu.Lock()
			cb.state = "half-open"
			cb.mu.Unlock()
		} else {
			return ErrCircuitOpen
		}
	}

	// Execute function
	err := fn()

	// Update state based on result
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailureTime = time.Now()

		if cb.failures >= cb.maxFailures {
			cb.state = "open"
		}
		return err
	}

	// Success - reset
	if cb.state == "half-open" {
		cb.state = "closed"
	}
	cb.failures = 0

	return nil
}

// GetState returns current circuit breaker state
func (cb *CircuitBreaker) GetState() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// =============================================================================
// Retry with Exponential Backoff
// =============================================================================

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts:  3,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     5 * time.Second,
	Multiplier:   2.0,
}

// RetryWithBackoff retries a function with exponential backoff
func RetryWithBackoff(ctx context.Context, config RetryConfig, fn func() error) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Try execution
		err := fn()
		if err == nil {
			return nil // Success!
		}

		lastErr = err

		// Check if we should retry
		if !IsRetryable(err) {
			return fmt.Errorf("non-retryable error on attempt %d: %w", attempt, err)
		}

		// Last attempt failed
		if attempt == config.MaxAttempts {
			break
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled after %d attempts: %w", attempt, ctx.Err())
		case <-time.After(delay):
			// Calculate next delay with exponential backoff
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}
	}

	return fmt.Errorf("max retry attempts (%d) reached: %w", config.MaxAttempts, lastErr)
}

// IsRetryable determines if an error should be retried
func IsRetryable(err error) bool {
	// Don't retry canceled contexts
	if errors.Is(err, context.Canceled) {
		return false
	}

	// Don't retry deadline exceeded
	if errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Retry network timeouts
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	// Retry temporary network errors
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	// Default: retry
	return true
}

// =============================================================================
// Timeout Pattern
// =============================================================================

// WithTimeout executes function with timeout
func WithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- fn(ctx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// =============================================================================
// Error Aggregation
// =============================================================================

// MultiError holds multiple errors
type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}

	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}

	var result string
	result = fmt.Sprintf("%d errors occurred:\n", len(m.Errors))
	for i, err := range m.Errors {
		result += fmt.Sprintf("  %d. %v\n", i+1, err)
	}
	return result
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}

func (m *MultiError) ErrorOrNil() error {
	if !m.HasErrors() {
		return nil
	}
	return m
}

// =============================================================================
// Usage Examples
// =============================================================================

// Example 1: Service with Circuit Breaker
type PaymentService struct {
	gateway        PaymentGateway
	circuitBreaker *CircuitBreaker
}

type PaymentGateway interface {
	Charge(ctx context.Context, amount float64) error
}

func NewPaymentService(gateway PaymentGateway) *PaymentService {
	return &PaymentService{
		gateway: gateway,
		// Open circuit after 5 failures, reset after 30 seconds
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
	}
}

func (s *PaymentService) Charge(ctx context.Context, amount float64) error {
	return s.circuitBreaker.Execute(ctx, func() error {
		return s.gateway.Charge(ctx, amount)
	})
}

// Example 2: Service with Retry
type ExternalAPIService struct {
	client HTTPClient
}

type HTTPClient interface {
	Post(ctx context.Context, url string, data interface{}) error
}

func (s *ExternalAPIService) SendNotification(ctx context.Context, message string) error {
	return RetryWithBackoff(ctx, DefaultRetryConfig, func() error {
		return s.client.Post(ctx, "/notifications", map[string]string{
			"message": message,
		})
	})
}

// Example 3: Combined Circuit Breaker + Retry
type RobustService struct {
	client         HTTPClient
	circuitBreaker *CircuitBreaker
}

func NewRobustService(client HTTPClient) *RobustService {
	return &RobustService{
		client:         client,
		circuitBreaker: NewCircuitBreaker(3, 10*time.Second),
	}
}

func (s *RobustService) CallAPI(ctx context.Context, endpoint string) error {
	// Circuit breaker wraps retry logic
	return s.circuitBreaker.Execute(ctx, func() error {
		// Retry with backoff inside circuit breaker
		return RetryWithBackoff(ctx, RetryConfig{
			MaxAttempts:  2,
			InitialDelay: 50 * time.Millisecond,
			MaxDelay:     500 * time.Millisecond,
			Multiplier:   2.0,
		}, func() error {
			return s.client.Post(ctx, endpoint, nil)
		})
	})
}

// Example 4: Batch Operations with Error Aggregation
type BatchService struct {
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, entity interface{}) error
}

func (s *BatchService) ImportUsers(ctx context.Context, users []map[string]interface{}) error {
	var multiErr MultiError

	for i, userData := range users {
		err := s.repo.Create(ctx, userData)
		if err != nil {
			// Collect error but continue processing
			multiErr.Add(fmt.Errorf("user %d: %w", i, err))
		}
	}

	return multiErr.ErrorOrNil()
}

// Example 5: Timeout on Handler Level
func ExampleHandler(service *PaymentService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req struct {
			Amount float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		// Enforce 5-second timeout
		err := WithTimeout(c.Request.Context(), 5*time.Second, func(ctx context.Context) error {
			return service.Charge(ctx, req.Amount)
		})

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(504, gin.H{"error": "request timeout"})
				return
			}
			if errors.Is(err, ErrCircuitOpen) {
				c.JSON(503, gin.H{"error": "service temporarily unavailable"})
				return
			}
			c.JSON(500, gin.H{"error": "payment failed"})
			return
		}

		c.JSON(200, gin.H{"status": "success"})
	}
}

// Example 6: Graceful Degradation
type UserService struct {
	cache     Cache
	database  Database
	analytics Analytics
}

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
}

type Database interface {
	GetByID(ctx context.Context, id uint) (interface{}, error)
}

type Analytics interface {
	Track(event string, properties map[string]interface{}) error
}

func (s *UserService) GetUser(ctx context.Context, id uint) (interface{}, error) {
	// Try cache first (fast but optional)
	cached, err := s.cache.Get(ctx, fmt.Sprintf("user:%d", id))
	if err == nil && cached != nil {
		return cached, nil
	}

	// Fallback to database (critical)
	user, err := s.database.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %w", err)
	}

	// Update cache asynchronously (non-critical)
	go func() {
		_ = s.cache.Set(context.Background(), fmt.Sprintf("user:%d", id), user)
	}()

	// Track analytics asynchronously (non-critical)
	go func() {
		_ = s.analytics.Track("user_viewed", map[string]interface{}{
			"user_id": id,
		})
	}()

	return user, nil
}
