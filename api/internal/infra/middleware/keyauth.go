package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// KeyAuthConfig holds API key authentication configuration
type KeyAuthConfig struct {
	// KeyLookup is where to look for the key: "header:X-API-Key", "query:api_key", "cookie:api_key"
	// Default: "header:X-API-Key"
	KeyLookup string

	// Validator is a function to validate the API key
	// Returns true if key is valid
	Validator func(key string) bool

	// ContextKey is the key used to store the API key in context
	// Default: "api_key"
	ContextKey string

	// ErrorMessage is the message returned when authentication fails
	// Default: "Invalid or missing API key"
	ErrorMessage string

	// AuthScheme is the scheme expected before the key (e.g., "Bearer")
	// Default: "" (no scheme)
	AuthScheme string
}

// DefaultKeyAuthConfig returns default key auth configuration
func DefaultKeyAuthConfig() KeyAuthConfig {
	return KeyAuthConfig{
		KeyLookup:    "header:X-API-Key",
		ContextKey:   "api_key",
		ErrorMessage: "Invalid or missing API key",
		AuthScheme:   "",
	}
}

// KeyAuth returns API key authentication middleware
func KeyAuth(validator func(key string) bool) gin.HandlerFunc {
	cfg := DefaultKeyAuthConfig()
	cfg.Validator = validator
	return KeyAuthWithConfig(cfg)
}

// KeyAuthWithConfig returns API key authentication middleware with custom config
func KeyAuthWithConfig(cfg KeyAuthConfig) gin.HandlerFunc {
	if cfg.KeyLookup == "" {
		cfg.KeyLookup = "header:X-API-Key"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "api_key"
	}
	if cfg.ErrorMessage == "" {
		cfg.ErrorMessage = "Invalid or missing API key"
	}
	if cfg.Validator == nil {
		panic("KeyAuth middleware requires a validator function")
	}

	// Parse key lookup
	parts := strings.SplitN(cfg.KeyLookup, ":", 2)
	if len(parts) != 2 {
		panic("KeyAuth KeyLookup must be in format 'source:name'")
	}
	source := strings.ToLower(parts[0])
	name := parts[1]

	return func(c *gin.Context) {
		var key string

		// Extract key based on source
		switch source {
		case "header":
			key = c.GetHeader(name)
			// Handle auth scheme if present
			if cfg.AuthScheme != "" && strings.HasPrefix(key, cfg.AuthScheme+" ") {
				key = strings.TrimPrefix(key, cfg.AuthScheme+" ")
			}
		case "query":
			key = c.Query(name)
		case "cookie":
			key, _ = c.Cookie(name)
		case "form":
			key = c.PostForm(name)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Invalid key auth configuration",
				},
			})
			return
		}

		// Validate key
		if key == "" || !cfg.Validator(key) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": cfg.ErrorMessage,
				},
			})
			return
		}

		// Store key in context
		c.Set(cfg.ContextKey, key)

		c.Next()
	}
}

// GetAPIKey retrieves the API key from context
func GetAPIKey(c *gin.Context) string {
	if key, exists := c.Get("api_key"); exists {
		return key.(string)
	}
	return ""
}
