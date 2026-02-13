package breaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrServiceUnavailable is returned when the circuit breaker is open
var ErrServiceUnavailable = errors.New("circuit breaker is open")

// State represents the circuit breaker state
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// Acceptable is the func to check if the error can be accepted
type Acceptable func(err error) bool

// Fallback is the func to be called if the request is rejected
type Fallback func(err error) error

// Promise interface defines the callbacks that returned by Breaker.Allow
type Promise interface {
	Accept()
	Reject(reason string)
}

// Breaker represents a circuit breaker
type Breaker interface {
	Name() string
	Allow() (Promise, error)
	AllowCtx(ctx context.Context) (Promise, error)
	Do(req func() error) error
	DoCtx(ctx context.Context, req func() error) error
	DoWithFallback(req func() error, fallback Fallback) error
	DoWithAcceptable(req func() error, acceptable Acceptable) error
	State() State
}

// Config holds circuit breaker configuration
type Config struct {
	Name string
	// Threshold is the number of consecutive failures to open the circuit
	Threshold int
	// Timeout is the duration the circuit stays open before transitioning to half-open
	Timeout time.Duration
	// MaxHalfOpenRequests is the max number of requests allowed in half-open state
	MaxHalfOpenRequests int
}

// DefaultConfig returns default circuit breaker configuration
func DefaultConfig() Config {
	return Config{
		Threshold:           5,
		Timeout:             10 * time.Second,
		MaxHalfOpenRequests: 1,
	}
}

type circuitBreaker struct {
	name                string
	threshold           int
	timeout             time.Duration
	maxHalfOpenRequests int

	mu               sync.Mutex
	state            State
	failures         int
	successes        int
	halfOpenRequests int
	lastFailureTime  time.Time
	lastStateChange  time.Time
	recentErrors     []string
	maxRecentErrors  int
}

// New creates a new circuit breaker
func New(cfg Config) Breaker {
	if cfg.Threshold <= 0 {
		cfg.Threshold = 5
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.MaxHalfOpenRequests <= 0 {
		cfg.MaxHalfOpenRequests = 1
	}

	return &circuitBreaker{
		name:                cfg.Name,
		threshold:           cfg.Threshold,
		timeout:             cfg.Timeout,
		maxHalfOpenRequests: cfg.MaxHalfOpenRequests,
		state:               StateClosed,
		lastStateChange:     time.Now(),
		maxRecentErrors:     5,
		recentErrors:        make([]string, 0, 5),
	}
}

func (cb *circuitBreaker) Name() string {
	return cb.name
}

func (cb *circuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.currentState()
}

func (cb *circuitBreaker) currentState() State {
	now := time.Now()

	switch cb.state {
	case StateOpen:
		if now.Sub(cb.lastStateChange) >= cb.timeout {
			cb.setState(StateHalfOpen)
		}
	}

	return cb.state
}

func (cb *circuitBreaker) setState(state State) {
	if cb.state == state {
		return
	}
	cb.state = state
	cb.lastStateChange = time.Now()

	switch state {
	case StateClosed:
		cb.failures = 0
		cb.successes = 0
	case StateHalfOpen:
		cb.halfOpenRequests = 0
		cb.successes = 0
	case StateOpen:
		// Keep failures count for logging
	}
}

func (cb *circuitBreaker) Allow() (Promise, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	state := cb.currentState()

	switch state {
	case StateClosed:
		return &promise{cb: cb}, nil
	case StateHalfOpen:
		if cb.halfOpenRequests >= cb.maxHalfOpenRequests {
			return nil, ErrServiceUnavailable
		}
		cb.halfOpenRequests++
		return &promise{cb: cb}, nil
	case StateOpen:
		return nil, ErrServiceUnavailable
	default:
		return nil, ErrServiceUnavailable
	}
}

func (cb *circuitBreaker) AllowCtx(ctx context.Context) (Promise, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return cb.Allow()
	}
}

func (cb *circuitBreaker) Do(req func() error) error {
	return cb.DoWithAcceptable(req, func(err error) bool { return err == nil })
}

func (cb *circuitBreaker) DoCtx(ctx context.Context, req func() error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return cb.Do(req)
	}
}

func (cb *circuitBreaker) DoWithFallback(req func() error, fallback Fallback) error {
	err := cb.Do(req)
	if err != nil && fallback != nil {
		if errors.Is(err, ErrServiceUnavailable) {
			return fallback(err)
		}
	}
	return err
}

func (cb *circuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	promise, err := cb.Allow()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			promise.Reject("panic")
			panic(r)
		}
	}()

	err = req()
	if acceptable(err) {
		promise.Accept()
	} else {
		reason := "error"
		if err != nil {
			reason = err.Error()
		}
		promise.Reject(reason)
	}

	return err
}

func (cb *circuitBreaker) onSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.failures = 0
	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.maxHalfOpenRequests {
			cb.setState(StateClosed)
		}
	}
}

func (cb *circuitBreaker) onFailure(reason string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()

	// Track recent errors
	if len(cb.recentErrors) >= cb.maxRecentErrors {
		cb.recentErrors = cb.recentErrors[1:]
	}
	cb.recentErrors = append(cb.recentErrors, reason)

	switch cb.state {
	case StateClosed:
		cb.failures++
		if cb.failures >= cb.threshold {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.setState(StateOpen)
	}
}

type promise struct {
	cb *circuitBreaker
}

func (p *promise) Accept() {
	p.cb.onSuccess()
}

func (p *promise) Reject(reason string) {
	p.cb.onFailure(reason)
}

// Group manages multiple circuit breakers by name
type Group struct {
	mu       sync.RWMutex
	breakers map[string]Breaker
	cfg      Config
}

// NewGroup creates a new circuit breaker group
func NewGroup(cfg Config) *Group {
	return &Group{
		breakers: make(map[string]Breaker),
		cfg:      cfg,
	}
}

// Get returns a circuit breaker by name, creating one if it doesn't exist
func (g *Group) Get(name string) Breaker {
	g.mu.RLock()
	cb, ok := g.breakers[name]
	g.mu.RUnlock()

	if ok {
		return cb
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// Double check after acquiring write lock
	if cb, ok = g.breakers[name]; ok {
		return cb
	}

	cfg := g.cfg
	cfg.Name = name
	cb = New(cfg)
	g.breakers[name] = cb

	return cb
}

// Do executes a function with circuit breaker protection
func (g *Group) Do(name string, req func() error) error {
	return g.Get(name).Do(req)
}

// DoWithFallback executes a function with circuit breaker and fallback
func (g *Group) DoWithFallback(name string, req func() error, fallback Fallback) error {
	return g.Get(name).DoWithFallback(req, fallback)
}
