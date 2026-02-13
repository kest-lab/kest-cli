---
name: logging-standards
description: Structured logging best practices including levels, context propagation, and monitoring
version: 1.0.0
category: development
tags: [logging, monitoring, debugging, observability]
author: ZGO Team
updated: 2026-01-24
---

# Logging Standards

## 📋 Purpose

This skill provides comprehensive logging standards for the ZGO Go backend project, ensuring consistent, structured, and production-ready logging across all modules.

## 🎯 When to Use

- Setting up logging for a new module
- Debugging production issues
- Implementing monitoring and observability
- Code review for logging practices
- Troubleshooting application behavior

## ⚙️ Prerequisites

- [ ] Understanding of Go and structured logging
- [ ] Familiarity with `logrus` or `zap` libraries
- [ ] Knowledge of log levels and their meanings

## 📚 Core Principles

### 1. Wide Events (CRITICAL)

**Principle**: Emit comprehensive log entries with all relevant context in a single event.

**Why**: Simplifies log querying and analysis. One log event contains everything needed.

```go
// ✅ CORRECT - Wide event with full context
logger.WithFields(logrus.Fields{
    "user_id":     userID,
    "request_id":  requestID,
    "action":      "user_login",
    "ip_address":  ipAddr,
    "user_agent":  userAgent,
    "login_time":  time.Now().Unix(),
    "status":      "success",
}).Info("User logged in successfully")

// ❌ WRONG - Multiple narrow events
logger.Info("User login")
logger.Infof("User ID: %d", userID)
logger.Infof("IP: %s", ipAddr)
logger.Info("Login successful")
```

**Benefits**:
- One query finds all related information
- Easier correlation in log aggregation tools
- Better performance (fewer log calls)
- Cleaner log streams

---

### 2. High Cardinality & Dimensionality (CRITICAL)

**Principle**: Include high-cardinality identifiers (request_id, user_id, session_id) in every log.

**Why**: Enables precise filtering and correlation across distributed systems.

```go
// ✅ CORRECT - High cardinality fields
logger.WithFields(logrus.Fields{
    "request_id":    ctx.Value("request_id"),      // Unique per request
    "user_id":       user.ID,                      // Unique per user
    "tenant_id":     tenant.ID,                    // Unique per tenant
    "correlation_id": ctx.Value("correlation_id"), // Trace across services
    "session_id":    sessionID,                    // Unique per session
}).Info("Processing payment")

// ❌ WRONG - Low cardinality only
logger.WithFields(logrus.Fields{
    "module": "payment",
    "status": "processing",
}).Info("Processing payment")
```

**Key Identifiers to Always Include**:
- `request_id` - Trace a single HTTP request
- `user_id` - Track user actions
- `tenant_id` - Multi-tenancy support
- `correlation_id` - Distributed tracing
- `session_id` - User session tracking

---

### 3. Business Context (CRITICAL)

**Principle**: Log business-relevant information, not just technical details.

**Why**: Helps non-technical stakeholders understand system behavior.

```go
// ✅ CORRECT - Business context
logger.WithFields(logrus.Fields{
    "user_id":        user.ID,
    "order_id":       order.ID,
    "order_total":    order.Total,
    "payment_method": "credit_card",
    "subscription":   "premium",
    "conversion":     "trial_to_paid",
}).Info("User completed purchase")

// ❌ WRONG - Only technical details
logger.WithFields(logrus.Fields{
    "table":  "orders",
    "action": "insert",
    "rows":   1,
}).Info("Database insert")
```

**Business Events to Log**:
- User registration, login, logout
- Purchases, subscriptions, upgrades
- Important state changes
- Feature usage
- Errors affecting business operations

---

### 4. Environment Characteristics (CRITICAL)

**Principle**: Include environment information in every log (via global fields).

**Why**: Distinguishes logs from different environments and deployment contexts.

```go
// ✅ CORRECT - Set global fields at startup
logger := logrus.New()
logger.SetFormatter(&logrus.JSONFormatter{})

// Global context for all logs
logger = logger.WithFields(logrus.Fields{
    "environment": os.Getenv("ENV"),              // dev/staging/production
    "service":     "zgo-api",
    "version":     version.Version,
    "hostname":    hostname,
    "region":      os.Getenv("AWS_REGION"),
    "pod_name":    os.Getenv("POD_NAME"),        // Kubernetes
}).Logger
```

**Environment Fields**:
- `environment` - dev, staging, production
- `service` - Service name
- `version` - Application version
- `hostname` - Server hostname
- `region` - Cloud region
- `pod_name` - Container/pod identifier

---

### 5. Single Logger Instance (HIGH)

**Principle**: Use a single, globally configured logger throughout the application.

**Why**: Ensures consistent formatting, levels, and configuration.

```go
// ✅ CORRECT - Single logger (pkg/logger/logger.go)
package logger

import (
    "os"
    "github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
    Log = logrus.New()
    Log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02T15:04:05.000Z",
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "level",
            logrus.FieldKeyMsg:   "message",
        },
    })
    
    // Set log level from environment
    if level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
        Log.SetLevel(level)
    } else {
        Log.SetLevel(logrus.InfoLevel)
    }
    
    // Add global fields
    Log = Log.WithFields(logrus.Fields{
        "service":     "zgo-api",
        "environment": os.Getenv("ENV"),
        "version":     os.Getenv("VERSION"),
    }).Logger
}

// Usage in modules
import "github.com/zgiai/zgo/pkg/logger"

func someFunction() {
    logger.Log.WithFields(logrus.Fields{
        "user_id": 123,
    }).Info("User action")
}
```

---

### 6. Middleware Pattern (HIGH)

**Principle**: Use middleware to log all HTTP requests automatically.

**Why**: Consistent request logging, automatic request_id propagation.

```go
// ✅ CORRECT - Logging middleware
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        requestID := uuid.New().String()
        
        // Store request_id in context
        c.Set("request_id", requestID)
        
        // Log request start
        logger.Log.WithFields(logrus.Fields{
            "request_id": requestID,
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "ip":         c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
        }).Info("Request started")
        
        // Process request
        c.Next()
        
        // Log request completion
        duration := time.Since(start)
        logger.Log.WithFields(logrus.Fields{
            "request_id": requestID,
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     c.Writer.Status(),
            "duration":   duration.Milliseconds(),
            "size":       c.Writer.Size(),
        }).Info("Request completed")
    }
}

// Register in main.go
r := gin.New()
r.Use(LoggingMiddleware())
```

---

### 7. Structure & Consistency (HIGH)

**Principle**: Use structured logging (JSON) with consistent field names.

**Why**: Machine-readable, queryable, works with log aggregation tools.

```go
// ✅ CORRECT - Structured JSON logging
logger.SetFormatter(&logrus.JSONFormatter{
    TimestampFormat: "2006-01-02T15:04:05.000Z",
})

logger.WithFields(logrus.Fields{
    "user_id":    123,
    "action":     "login",
    "ip_address": "192.168.1.1",
}).Info("User logged in")

// Output:
// {"timestamp":"2026-01-24T18:20:01.000Z","level":"info","message":"User logged in","user_id":123,"action":"login","ip_address":"192.168.1.1"}

// ❌ WRONG - Unstructured text logging
logger.Infof("User %d logged in from %s", userID, ipAddr)

// Output:
// 2026-01-24 18:20:01 INFO User 123 logged in from 192.168.1.1
```

**Standard Field Names**:
- Use **snake_case** for consistency
- Use **descriptive names**: `user_id` not `uid`
- Use **common names**: `timestamp` not `time` or `t`

---

## 📊 Log Levels

### Level Definitions

| Level | Usage | Example | Production |
|-------|-------|---------|------------|
| **DEBUG** | Detailed info for debugging | Variable values, flow trace | ❌ Disabled |
| **INFO** | General information | Request/response, state changes | ✅ Enabled |
| **WARN** | Warning conditions | Deprecated API usage, retries | ✅ Enabled |
| **ERROR** | Error conditions | Failed operations, exceptions | ✅ Enabled |
| **FATAL** | Critical errors | Startup failures, panics | ✅ Enabled |

### When to Use Each Level

#### DEBUG
```go
// ✅ Use for development debugging
logger.WithFields(logrus.Fields{
    "user_id":      user.ID,
    "query":        query,
    "result_count": len(results),
}).Debug("Database query executed")

// ✅ Trace execution flow
logger.Debug("Entering function ProcessPayment")
logger.Debug("Validation passed")
logger.Debug("Exiting function ProcessPayment")
```

#### INFO
```go
// ✅ Normal application events
logger.WithFields(logrus.Fields{
    "user_id": user.ID,
    "action":  "login",
}).Info("User logged in")

// ✅ Important state changes
logger.WithFields(logrus.Fields{
    "order_id": order.ID,
    "status":   "completed",
}).Info("Order completed")
```

#### WARN
```go
// ✅ Recoverable issues
logger.WithFields(logrus.Fields{
    "user_id": user.ID,
    "attempt": 3,
}).Warn("Password attempt failed, user locked")

// ✅ Deprecated features
logger.Warn("Using deprecated API endpoint: /api/v1/users")

// ✅ Resource limits approaching
logger.WithFields(logrus.Fields{
    "disk_usage": "85%",
}).Warn("Disk usage approaching limit")
```

#### ERROR
```go
// ✅ Operation failures
logger.WithFields(logrus.Fields{
    "user_id": user.ID,
    "error":   err.Error(),
}).Error("Failed to create user")

// ✅ External service failures
logger.WithFields(logrus.Fields{
    "service": "payment_gateway",
    "error":   err.Error(),
}).Error("Payment gateway unavailable")
```

#### FATAL
```go
// ✅ Unrecoverable errors (terminates app)
logger.WithFields(logrus.Fields{
    "error": err.Error(),
}).Fatal("Failed to connect to database")

// ⚠️ Use sparingly - only for startup failures
```

---

## 🔧 Implementation Patterns

### Pattern 1: Request-Scoped Logger

Create a logger with request context for each HTTP request.

```go
// middleware/context_logger.go
func ContextLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("request_id", requestID)
        
        // Create request-scoped logger
        requestLogger := logger.Log.WithFields(logrus.Fields{
            "request_id": requestID,
            "path":       c.Request.URL.Path,
            "method":     c.Request.Method,
        })
        
        c.Set("logger", requestLogger)
        c.Next()
    }
}

// Usage in handler
func (h *Handler) Get(c *gin.Context) {
    log := c.MustGet("logger").(*logrus.Entry)
    
    log.WithField("user_id", userID).Info("Fetching user")
    
    user, err := h.service.GetByID(ctx, userID)
    if err != nil {
        log.WithError(err).Error("Failed to fetch user")
        return
    }
    
    log.Info("User fetched successfully")
}
```

### Pattern 2: Context Propagation

Pass logger through context for service/repository layers.

```go
// pkg/logger/context.go
type loggerKey struct{}

func WithLogger(ctx context.Context, log *logrus.Entry) context.Context {
    return context.WithValue(ctx, loggerKey{}, log)
}

func FromContext(ctx context.Context) *logrus.Entry {
    if log, ok := ctx.Value(loggerKey{}).(*logrus.Entry); ok {
        return log
    }
    return logger.Log.WithField("source", "unknown")
}

// Usage in service
func (s *service) Create(ctx context.Context, req *CreateRequest) error {
    log := logger.FromContext(ctx)
    
    log.WithField("username", req.Username).Info("Creating user")
    
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        log.WithError(err).Error("Failed to create user in database")
        return err
    }
    
    log.WithField("user_id", user.ID).Info("User created successfully")
    return nil
}
```

### Pattern 3: Error Logging

Always log errors with full context.

```go
// ✅ CORRECT - Log with context and error
func (s *service) ProcessPayment(ctx context.Context, orderID uint, amount float64) error {
    log := logger.FromContext(ctx)
    
    log.WithFields(logrus.Fields{
        "order_id": orderID,
        "amount":   amount,
    }).Info("Processing payment")
    
    if err := s.paymentGateway.Charge(amount); err != nil {
        log.WithFields(logrus.Fields{
            "order_id": orderID,
            "amount":   amount,
            "error":    err.Error(),
            "gateway":  "stripe",
        }).Error("Payment gateway charge failed")
        return fmt.Errorf("charge failed: %w", err)
    }
    
    log.WithField("order_id", orderID).Info("Payment processed successfully")
    return nil
}
```

### Pattern 4: Sensitive Data Handling

Never log sensitive information.

```go
// ❌ WRONG - Logging sensitive data
logger.WithFields(logrus.Fields{
    "password":     user.Password,
    "credit_card":  payment.CardNumber,
    "ssn":          user.SSN,
}).Info("User created")

// ✅ CORRECT - Mask or omit sensitive data
logger.WithFields(logrus.Fields{
    "user_id":      user.ID,
    "email":        user.Email,
    "card_last4":   payment.CardLast4,  // Only last 4 digits
}).Info("User created")

// ✅ CORRECT - Helper function
func sanitizeUser(user *User) logrus.Fields {
    return logrus.Fields{
        "user_id": user.ID,
        "email":   user.Email,
        "username": user.Username,
        // Password excluded
    }
}

logger.WithFields(sanitizeUser(user)).Info("User logged in")
```

---

## 🚫 Anti-Patterns to Avoid

### 1. String Interpolation in Messages

```go
// ❌ WRONG - Variables in message string
logger.Infof("User %d logged in from %s", userID, ipAddr)

// ✅ CORRECT - Variables in structured fields
logger.WithFields(logrus.Fields{
    "user_id":    userID,
    "ip_address": ipAddr,
}).Info("User logged in")
```

### 2. Multiple Logs for Single Event

```go
// ❌ WRONG - Multiple narrow logs
logger.Info("Starting payment")
logger.Infof("Amount: %.2f", amount)
logger.Infof("Order: %d", orderID)
logger.Info("Payment completed")

// ✅ CORRECT - Single wide log
logger.WithFields(logrus.Fields{
    "order_id": orderID,
    "amount":   amount,
    "status":   "completed",
}).Info("Payment processed")
```

### 3. Logging Without Context

```go
// ❌ WRONG - No context
logger.Info("Error occurred")

// ✅ CORRECT - With full context
logger.WithFields(logrus.Fields{
    "user_id":    userID,
    "operation":  "create_order",
    "error":      err.Error(),
}).Error("Failed to create order")
```

### 4. Excessive DEBUG Logging

```go
// ❌ WRONG - Too verbose
logger.Debug("Entering function")
logger.Debug("x = 1")
logger.Debug("y = 2")
logger.Debug("Calling helper")
logger.Debug("Helper returned")
logger.Debug("Exiting function")

// ✅ CORRECT - Meaningful DEBUG logs
logger.WithField("input", input).Debug("Processing input")
logger.WithField("result", result).Debug("Function completed")
```

---

## ✅ Verification Checklist

- [ ] Using structured logging (JSON) not text
- [ ] Single logger instance configured at startup
- [ ] Log levels correctly set (INFO in production)
- [ ] Global fields included (service, environment, version)
- [ ] Request middleware logs all HTTP requests
- [ ] Request_id propagated through context
- [ ] High-cardinality fields included (user_id, request_id)
- [ ] Business context logged for important events
- [ ] Errors logged with full context
- [ ] No sensitive data in logs
- [ ] Consistent field naming (snake_case)
- [ ] Wide events (not multiple narrow logs)

---

## 📚 Complete Example

See [`.agent/skills/logging-standards/examples/logging-setup.go`](./examples/logging-setup.go)

---

## 🔧 Validation Script

Run the logging standards validation:

```bash
.agent/skills/logging-standards/scripts/validate-logging.sh <module_name>
```

---

## 🔗 Related Skills

- [`api-development`](../api-development/): API logging patterns
- [`coding-standards`](../coding-standards/): General code quality
- [`module-creation`](../module-creation/): Module structure

---

## 📝 Quick Reference

```go
// Setup (main.go or init)
logger.Init()

// HTTP middleware
r.Use(middleware.ContextLogger())

// In handlers
log := c.MustGet("logger").(*logrus.Entry)
log.WithField("user_id", userID).Info("Action performed")

// In services
log := logger.FromContext(ctx)
log.WithError(err).Error("Operation failed")

// Common patterns
log.Info("Normal event")
log.Warn("Warning condition")
log.WithError(err).Error("Error occurred")
log.WithFields(logrus.Fields{
    "user_id": 123,
    "action":  "login",
}).Info("User action")
```

---

**Version**: 1.0.0  
**Last Updated**: 2026-01-24  
**Maintainer**: ZGO Team
