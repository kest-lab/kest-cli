package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ============================================================================
// Logger Setup
// ============================================================================

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// JSON formatter for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	if level, err := logrus.ParseLevel(logLevel); err == nil {
		Log.SetLevel(level)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}

	// Add global fields (environment characteristics)
	Log = Log.WithFields(logrus.Fields{
		"service":     "zgo-api",
		"environment": getEnv("ENV", "development"),
		"version":     getEnv("VERSION", "1.0.0"),
		"hostname":    getHostname(),
	}).Logger

	Log.Info("Logger initialized")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

// ============================================================================
// Context Logger - Pattern for Propagation
// ============================================================================

type loggerKey struct{}

// WithLogger adds logger to context
func WithLogger(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) *logrus.Entry {
	if log, ok := ctx.Value(loggerKey{}).(*logrus.Entry); ok {
		return log
	}
	return Log.WithField("source", "unknown")
}

// ============================================================================
// HTTP Middleware
// ============================================================================

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Create request-scoped logger
		requestLogger := Log.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// Store in context for handlers to use
		c.Set("logger", requestLogger)

		// Log request start
		requestLogger.Info("Request started")

		// Process request
		c.Next()

		// Log request completion
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logFields := logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     statusCode,
			"duration":   duration.Milliseconds(),
			"size":       c.Writer.Size(),
		}

		// Log at appropriate level based on status code
		if statusCode >= 500 {
			requestLogger.WithFields(logFields).Error("Request failed")
		} else if statusCode >= 400 {
			requestLogger.WithFields(logFields).Warn("Request error")
		} else {
			requestLogger.WithFields(logFields).Info("Request completed")
		}
	}
}

// ============================================================================
// Example: Service Layer with Logging
// ============================================================================

type User struct {
	ID       uint
	Username string
	Email    string
	Status   string
}

type UserService struct {
	// dependencies
}

// Create demonstrates proper logging in service layer
func (s *UserService) Create(ctx context.Context, username, email string) (*User, error) {
	// ✅ Get logger from context
	log := FromContext(ctx)

	// ✅ Log with business context
	log.WithFields(logrus.Fields{
		"username": username,
		"email":    email,
		"action":   "create_user",
	}).Info("Creating new user")

	// Simulate validation
	if username == "" {
		log.Warn("User creation failed: username is empty")
		return nil, fmt.Errorf("username is required")
	}

	// Simulate database operation
	user := &User{
		ID:       123,
		Username: username,
		Email:    email,
		Status:   "active",
	}

	// ✅ Log successful creation with user_id
	log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}).Info("User created successfully")

	return user, nil
}

// Update demonstrates error logging
func (s *UserService) Update(ctx context.Context, userID uint, updates map[string]interface{}) error {
	log := FromContext(ctx)

	log.WithFields(logrus.Fields{
		"user_id": userID,
		"updates": updates,
	}).Info("Updating user")

	// Simulate error
	if userID == 0 {
		// ✅ Log error with full context
		log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   "invalid user ID",
		}).Error("User update failed")

		return fmt.Errorf("invalid user ID")
	}

	log.WithField("user_id", userID).Info("User updated successfully")
	return nil
}

// ProcessPayment demonstrates wide events
func (s *UserService) ProcessPayment(ctx context.Context, userID uint, amount float64, method string) error {
	log := FromContext(ctx)

	// ✅ WIDE EVENT - All context in one log
	log.WithFields(logrus.Fields{
		"user_id":        userID,
		"amount":         amount,
		"currency":       "USD",
		"payment_method": method,
		"action":         "process_payment",
		"timestamp":      time.Now().Unix(),
	}).Info("Processing payment")

	// Simulate payment gateway call
	if err := callPaymentGateway(amount); err != nil {
		// ✅ Error with full context
		log.WithFields(logrus.Fields{
			"user_id":        userID,
			"amount":         amount,
			"payment_method": method,
			"error":          err.Error(),
			"gateway":        "stripe",
		}).Error("Payment gateway failed")

		return err
	}

	// ✅ Success with business outcome
	log.WithFields(logrus.Fields{
		"user_id":        userID,
		"amount":         amount,
		"payment_method": method,
		"status":         "completed",
		"transaction_id": "txn_123456",
	}).Info("Payment completed successfully")

	return nil
}

func callPaymentGateway(amount float64) error {
	// Simulate payment gateway
	return nil
}

// ============================================================================
// Example: Handler with Logging
// ============================================================================

type UserHandler struct {
	service *UserService
}

// CreateUser demonstrates handler logging
func (h *UserHandler) CreateUser(c *gin.Context) {
	// ✅ Get request-scoped logger
	log := c.MustGet("logger").(*logrus.Entry)

	// Parse request
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("error", err.Error()).Warn("Invalid request body")
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// Create context with logger
	ctx := WithLogger(c.Request.Context(), log)

	// Call service
	user, err := h.service.Create(ctx, req.Username, req.Email)
	if err != nil {
		log.WithError(err).Error("Failed to create user")
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	// Success
	log.WithField("user_id", user.ID).Info("User creation successful")
	c.JSON(201, user)
}

// ============================================================================
// Example: Different Log Levels
// ============================================================================

func DemonstrateLogLevels(ctx context.Context) {
	log := FromContext(ctx)

	// DEBUG - Detailed debugging information
	log.WithFields(logrus.Fields{
		"query":        "SELECT * FROM users WHERE id = ?",
		"params":       []interface{}{123},
		"result_count": 1,
	}).Debug("Database query executed")

	// INFO - General information
	log.WithFields(logrus.Fields{
		"user_id": 123,
		"action":  "login",
	}).Info("User logged in")

	// WARN - Warning conditions
	log.WithFields(logrus.Fields{
		"user_id": 123,
		"attempt": 3,
		"max":     5,
	}).Warn("Multiple failed login attempts")

	// ERROR - Error conditions
	log.WithFields(logrus.Fields{
		"user_id": 123,
		"error":   "database connection timeout",
	}).Error("Failed to fetch user")

	// FATAL - Critical errors (terminates app)
	// log.WithField("error", "database unreachable").Fatal("Cannot start application")
}

// ============================================================================
// Example: Sensitive Data Handling
// ============================================================================

func LogUserCreation(user *User) {
	// ❌ WRONG - Logs password
	// Log.WithFields(logrus.Fields{
	//     "username": user.Username,
	//     "password": user.Password,  // NEVER LOG PASSWORDS
	//     "email":    user.Email,
	// }).Info("User created")

	// ✅ CORRECT - Omits sensitive data
	Log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
		// Password NOT included
	}).Info("User created")
}

// Helper to sanitize user data for logging
func sanitizeUser(user *User) logrus.Fields {
	return logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	}
}

// ============================================================================
// Main - Example Setup
// ============================================================================

func main() {
	// Initialize logger
	InitLogger()

	// Setup Gin
	r := gin.New()

	// Add logging middleware
	r.Use(LoggingMiddleware())

	// Setup routes
	service := &UserService{}
	handler := &UserHandler{service: service}

	r.POST("/api/users", handler.CreateUser)

	// Start server
	Log.Info("Server starting on :8080")
	r.Run(":8080")
}
