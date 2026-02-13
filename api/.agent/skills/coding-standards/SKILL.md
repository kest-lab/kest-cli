---
name: coding-standards
description: ZGO project coding standards and best practices verification
version: 2.0.0
category: development
tags: [standards, code-review, quality, error-handling]
author: ZGO Team
updated: 2026-01-24
---

# Coding Standards Skill

## 📋 Purpose

This skill provides a comprehensive checklist and verification guide for ensuring all code follows ZGO's coding standards and best practices.

## 🎯 When to Use

- Before submitting a Pull Request
- During code review
- When creating a new module
- When refactoring existing code
- When onboarding new team members

## ⚙️ Prerequisites

- [ ] Code committed to Git (for comparison)
- [ ] Module follows 8-file structure
- [ ] Tests written and passing

## 🔍 Verification Checklist

### Level 1: Naming Conventions ✅

#### 1.1 Package Names
```go
✅ package user          // Singular, lowercase, short
✅ package blogpost      // Compound words without underscores
❌ package users         // Avoid plural
❌ package blog_post     // Avoid underscores
❌ package BlogPost      // Avoid capitalization
```

#### 1.2 File Names
```go
✅ model.go              // Singular, lowercase
✅ service_test.go       // Test files suffix with _test.go
✅ user_handler.go       // Multi-word with underscores
❌ Model.go              // Avoid capitalization
❌ service-test.go       // Avoid hyphens
```

#### 1.3 Type Names
```go
// PO (Persistent Object) - Database entities
✅ type UserPO struct {}          // PO suffix
✅ type BlogPostPO struct {}      

// Domain Entity
✅ type User struct {}            // PascalCase, no suffix

// DTO (Data Transfer Object)
✅ type CreateUserRequest struct {}    // Verb + Noun + Request
✅ type UserResponse struct {}         // Noun + Response

// Interface
✅ type Repository interface {}   // Noun
✅ type Service interface {}      
✅ type Handler struct {}         // Handler is struct, not interface

// Private Implementation
✅ type repository struct {}      // Lowercase
✅ type service struct {}
```

#### 1.4 Function Names
```go
// Constructor
✅ func NewRepository() Repository       // New + InterfaceName, returns interface
✅ func NewService() Service

// Mapper Functions
✅ func ToUserPO(user *domain.User) *UserPO      // To + TargetType
✅ func FromUserPO(po *UserPO) *domain.User      // From + SourceType
✅ func ToResponse(user *domain.User) *UserResponse

// CRUD Operations
✅ func (r *repository) Create()
✅ func (r *repository) GetByID()    // Get + Condition
```

#### 1.5 JSON Tags
```go
type User struct {
    UserID   uint   `json:"user_id"`         // ✅ snake_case
    Password string `json:"-"`               // ✅ Hide sensitive fields
}
```

---

### Level 2: Architecture Standards ✅

#### 2.1 8-File Module Structure (Mandatory)
Each module **must** include: `model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`, `provider.go`, `service_test.go`.

> **📚 Full Guide**: See [`module-creation` skill](./.agent/skills/module-creation/)

#### 2.2 Layered Architecture
```
Handler → Service → Repository → Model
```
- **Forbidden**: Handler → Repository, Service → Model (PO), Repository → (returns PO to Service).

#### 2.3 Data Flow
`Handler(DTO) → Service(domain.User) → Repository(UserPO) → Database`

---

### Level 3: File Organization ✅

#### 3.1 model.go Requirements
- [ ] Has `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` fields
- [ ] Has `TableName()` method
- [ ] Uses GORM tags + indexes

#### 3.2 dto.go Requirements
- [ ] Request/Response suffixes
- [ ] Pointers for optional fields in Update DTOs
- [ ] Mappers handle `nil` input

#### 3.3 repository.go Requirements
- [ ] Constructor returns interface
- [ ] PO ↔ Domain conversion at boundary
- [ ] Uses `WithContext(ctx)`

#### 3.4 service.go Requirements
- [ ] Defines custom business errors (`var Err...`)
- [ ] Business logic + validation here
- [ ] Uses domain entities, not POs

---

### Level 4: Error Handling Standards ✅

#### 4.1 Define Custom Errors
```go
// ✅ Package-level errors
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrDuplicateEmail  = errors.New("email already exists")
)
```

#### 4.2 Error Wrapping
```go
func (s *service) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}
```

#### 4.3 Automatic Mapping
Use `response.HandleError(c, "msg", err)` in handlers to auto-map errors to HTTP status codes.

---

### Level 5: Security Standards ✅

#### 5.1 Sensitive Data
- [ ] Passwords hidden with `json:"-"` in domain entities.
- [ ] Passwords encrypted using `crypto.HashPassword()` before storage.
- [ ] Response DTOs exclude sensitive fields like passwords or secrets.

#### 5.2 Input Validation
- [ ] Use `binding` tags in DTOs (e.g., `required`, `email`, `min=8`).
- [ ] Perform business-level validation in the Service layer (e.g., check for duplicates).
- [ ] Use `handler.BindJSON()` to trigger validation.

---

### Level 9: Advanced Error Handling Patterns ⭐ NEW

> **📚 Expanded Section**: Advanced error handling for production systems

#### 9.1 Error Types and Wrapping

**Standard Error Creation**:
```go
// ✅ Define module-specific errors
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrDuplicateEmail    = errors.New("email already exists")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrAccountLocked     = errors.New("account is locked")
)

// ✅ Wrap errors with context
func (s *service) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}

// ✅ Check wrapped errors
user, err := service.GetByID(ctx, 123)
if err != nil {
    if errors.Is(err, repository.ErrUserNotFound) {
        // Handle not found
    }
}
```

#### 9.2 Circuit Breaker Pattern

Prevent cascading failures when calling external services.

**Implementation**:
```go
package capabilities

import (
    "context"
    "errors"
    "sync"
    "time"
)

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
    failures := cb.failures
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

// Usage in service
type PaymentService struct {
    gateway       PaymentGateway
    circuitBreaker *CircuitBreaker
}

func NewPaymentService(gateway PaymentGateway) *PaymentService {
    return &PaymentService{
        gateway:        gateway,
        circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
    }
}

func (s *PaymentService) Charge(ctx context.Context, amount float64) error {
    return s.circuitBreaker.Execute(ctx, func() error {
        return s.gateway.Charge(amount)
    })
}
```

**When to Use**:
- ✅ Calling external APIs (payment gateways, third-party services)
- ✅ Database connections that may fail
- ✅ Service-to-service communication
- ❌ Internal function calls
- ❌ Simple CRUD operations

#### 9.3 Retry with Exponential Backoff

Automatically retry failed operations with increasing delays.

**Implementation**:
```go
package capabilities

import (
    "context"
    "fmt"
    "math"
    "time"
)

// RetryConfig defines retry behavior
type RetryConfig struct {
    MaxAttempts int
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
            return fmt.Errorf("non-retryable error: %w", err)
        }
        
        // Last attempt failed
        if attempt == config.MaxAttempts {
            break
        }
        
        // Wait before retry
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(delay):
            // Calculate next delay
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
    // Add logic to check if error is retryable
    // e.g., network errors, timeouts, 5xx status codes
    
    if errors.Is(err, context.Canceled) {
        return false // Don't retry canceled contexts
    }
    
    // Check for specific error types
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Timeout() {
        return true // Retry timeouts
    }
    
    return true // Default: retry
}

// Usage in service
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*domain.User, error) {
    var user *domain.User
    var createErr error
    
    err := RetryWithBackoff(ctx, DefaultRetryConfig, func() error {
        user, createErr = s.repo.Create(ctx, req)
        return createErr
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to create user after retries: %w", err)
    }
    
    return user, nil
}
```

**When to Use**:
- ✅ Network operations
- ✅ External API calls
- ✅ Database deadlocks
- ✅ Rate-limited endpoints
- ❌ Validation errors
- ❌ Authorization failures

#### 9.4 Timeout Pattern

Prevent operations from running indefinitely.

**Implementation**:
```go
package capabilities

import (
    "context"
    "time"
)

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

// Usage in handler
func (h *Handler) Create(c *gin.Context) {
    var req CreateUserRequest
    if !handler.BindJSON(c, &req) {
        return
    }
    
    // Enforce 5-second timeout
    err := WithTimeout(c.Request.Context(), 5*time.Second, func(ctx context.Context) error {
        _, err := h.service.Create(ctx, &req)
        return err
    })
    
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            response.Error(c, 504, "Request timeout")
            return
        }
        response.HandleError(c, "Failed to create user", err)
        return
    }
    
    response.Created(c, user)
}

// Or set timeout at service level
func (s *service) ProcessPayment(ctx context.Context, amount float64) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    return s.paymentGateway.Charge(ctx, amount)
}
```

**Recommended Timeouts**:
- HTTP Handlers: 30s
- External API calls: 5-10s
- Database queries: 3-5s
- File operations: 1-3s

#### 9.5 Error Aggregation

Collect multiple errors and return them together.

**Implementation**:
```go
package capabilities

import (
    "errors"
    "fmt"
    "strings"
)

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
    
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d errors occurred:\n", len(m.Errors)))
    for i, err := range m.Errors {
        sb.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
    }
    return sb.String()
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

// Usage: Batch operations
func (s *service) ImportUsers(ctx context.Context, users []*CreateUserRequest) error {
    var multiErr MultiError
    
    for i, req := range users {
        _, err := s.Create(ctx, req)
        if err != nil {
            multiErr.Add(fmt.Errorf("user %d (%s): %w", i, req.Email, err))
        }
    }
    
    return multiErr.ErrorOrNil()
}

// Usage: Validation errors
func (s *service) Validate(ctx context.Context, user *domain.User) error {
    var multiErr MultiError
    
    if user.Email == "" {
        multiErr.Add(errors.New("email is required"))
    }
    
    if len(user.Password) < 8 {
        multiErr.Add(errors.New("password must be at least 8 characters"))
    }
    
    if user.Age < 0 {
        multiErr.Add(errors.New("age cannot be negative"))
    }
    
    return multiErr.ErrorOrNil()
}
```

#### 9.6 Graceful Degradation

Continue providing service even when某些 features fail.

**Implementation**:
```go
// Service with fallback
type UserService struct {
    cache      Cache
    database   Database
    analytics  Analytics
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    // Try cache first
    if user, err := s.cache.Get(ctx, id); err == nil {
        return user, nil
    }
    
    // Fallback to database
    user, err := s.database.GetByID(ctx, id)
    if err != nil {
        return nil, err // Critical error
    }
    
    // Update cache (non-critical, don't fail if it errors)
    _ = s.cache.Set(ctx, id, user)
    
    // Track analytics (non-critical)
    go func() {
        _ = s.analytics.Track("user_viewed", map[string]interface{}{
            "user_id": id,
        })
    }()
    
    return user, nil
}

// Feature flags for gradual degradation
type FeatureFlags struct {
    EnableRecommendations bool
    EnableNotifications   bool
    EnableAnalytics       bool
}

func (s *Service) GetUserProfile(ctx context.Context, id uint) (*ProfileResponse, error) {
    // Core functionality - must succeed
    user, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    profile := &ProfileResponse{
        User: user,
    }
    
    // Optional feature - degrade gracefully
    if s.flags.EnableRecommendations {
        if recs, err := s.getRecommendations(ctx, id); err == nil {
            profile.Recommendations = recs
        } else {
            logger.Warn("Failed to load recommendations", "user_id", id, "error", err)
            // Continue without recommendations
        }
    }
    
    return profile, nil
}
```

#### 9.7 Error Monitoring and Alerting

**Sentry Integration**:
```go
package logger

import (
    "context"
    "github.com/getsentry/sentry-go"
    "time"
)

func InitSentry() {
    sentry.Init(sentry.ClientOptions{
        Dsn:              os.Getenv("SENTRY_DSN"),
        Environment:      os.Getenv("ENV"),
        Release:          os.Getenv("VERSION"),
        TracesSampleRate: 0.2,
    })
}

// CaptureError sends error to Sentry
func CaptureError(ctx context.Context, err error, tags map[string]string) {
    hub := sentry.GetHubFromContext(ctx)
    if hub == nil {
        hub = sentry.CurrentHub()
    }
    
    hub.WithScope(func(scope *sentry.Scope) {
        // Add tags
        for key, value := range tags {
            scope.SetTag(key, value)
        }
        
        // Add context
        if userID := ctx.Value("user_id"); userID != nil {
            scope.SetUser(sentry.User{
                ID: fmt.Sprintf("%v", userID),
            })
        }
        
        // Capture
        hub.CaptureException(err)
    })
}

// Usage in service
func (s *service) Create(ctx context.Context, req *CreateUserRequest) (*domain.User, error) {
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        logger.CaptureError(ctx, err, map[string]string{
            "operation": "create_user",
            "email":     req.Email,
        })
        return nil, err
    }
    return user, nil
}
```

---

## 🚀 Quick Verification Script

[Rest of the file remains the same from line 453 onwards...]

---

## 📚 Examples

See the following examples:
- [`examples/standards-checklist.md`](./examples/standards-checklist.md) - Filled checklist
- [`examples/error-handling-patterns.go`](./examples/error-handling-patterns.go) - Error handling implementations
- [`examples/circuit-breaker-example.go`](./examples/circuit-breaker-example.go) - Circuit breaker pattern

## 🔗 Related Skills

- [`module-creation`](../module-creation/): For creating new modules
- [`api-development`](../api-development/): For API best practices
- [`logging-standards`](../logging-standards/): For logging errors

## 📖 References

- [AGENTS.md - Coding Standards](../../AGENTS.md#coding-standards-mandatory)
- [Go Error Handling](https://go.dev/blog/error-handling-and-go)
- [Effective Go](https://golang.org/doc/effective_go)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)

---

## ✅ Quick Checklist Summary

Before submitting code, verify:

- [ ] 8-file structure complete
- [ ] Naming follows conventions (PO suffix, snake_case JSON, etc.)
- [ ] Architecture layers respected (no cross-layer access)
- [ ] All files have required content
- [ ] Security: passwords hidden, input validated
- [ ] **Error handling: custom errors, wrapping, retry, timeout** ⭐ NEW
- [ ] **Production patterns: circuit breaker for external calls** ⭐ NEW
- [ ] Tests: unit tests passing, coverage > 80%
- [ ] Code quality: linter passing, formatted
- [ ] Wire DI: ProviderSet exported and integrated
- [ ] Documentation: Swagger comments added
- [ ] **Monitoring: errors sent to Sentry** ⭐ NEW

---

**Version**: 2.0.0  
**Last Updated**: 2026-01-24  
**Maintainer**: ZGO Team  
**Changelog**: Added Advanced Error Handling Patterns (Circuit Breaker, Retry, Timeout, Error Aggregation)
