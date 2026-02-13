package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDConfig holds RequestID middleware configuration
type RequestIDConfig struct {
	// Header is the name of the header to check/set
	// Default: X-Request-ID
	Header string

	// Generator is a function to generate request IDs
	// Default: UUID v4
	Generator func() string

	// ContextKey is the key used to store the request ID in context
	// Default: request_id
	ContextKey string
}

// DefaultRequestIDConfig returns default configuration
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		Header:     "X-Request-ID",
		Generator:  func() string { return uuid.New().String() },
		ContextKey: "request_id",
	}
}

// RequestID returns RequestID middleware with default config
func RequestID() gin.HandlerFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig())
}

// RequestIDWithConfig returns RequestID middleware with custom config
func RequestIDWithConfig(cfg RequestIDConfig) gin.HandlerFunc {
	// Set defaults
	if cfg.Header == "" {
		cfg.Header = "X-Request-ID"
	}
	if cfg.Generator == nil {
		cfg.Generator = func() string { return uuid.New().String() }
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "request_id"
	}

	return func(c *gin.Context) {
		// Check if request already has a request ID
		requestID := c.GetHeader(cfg.Header)

		// Generate new ID if not present
		if requestID == "" {
			requestID = cfg.Generator()
		}

		// Set request ID in context
		c.Set(cfg.ContextKey, requestID)

		// Set request ID in response header
		c.Header(cfg.Header, requestID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		return id.(string)
	}
	return ""
}
