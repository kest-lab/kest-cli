package middleware

import (
	"github.com/gin-gonic/gin"
)

// HelmetConfig holds security headers configuration
type HelmetConfig struct {
	// XSSProtection sets X-XSS-Protection header
	// Default: "1; mode=block"
	XSSProtection string

	// ContentTypeNosniff sets X-Content-Type-Options header
	// Default: "nosniff"
	ContentTypeNosniff string

	// XFrameOptions sets X-Frame-Options header
	// Default: "SAMEORIGIN"
	XFrameOptions string

	// HSTSMaxAge sets Strict-Transport-Security max-age
	// Default: 31536000 (1 year)
	// Set to 0 to disable HSTS
	HSTSMaxAge int

	// HSTSIncludeSubdomains adds includeSubDomains to HSTS
	// Default: true
	HSTSIncludeSubdomains bool

	// HSTSPreload adds preload to HSTS
	// Default: false
	HSTSPreload bool

	// ContentSecurityPolicy sets Content-Security-Policy header
	// Default: "" (disabled)
	ContentSecurityPolicy string

	// ReferrerPolicy sets Referrer-Policy header
	// Default: "strict-origin-when-cross-origin"
	ReferrerPolicy string

	// PermissionsPolicy sets Permissions-Policy header
	// Default: "" (disabled)
	PermissionsPolicy string

	// CrossOriginEmbedderPolicy sets Cross-Origin-Embedder-Policy header
	// Default: "" (disabled)
	CrossOriginEmbedderPolicy string

	// CrossOriginOpenerPolicy sets Cross-Origin-Opener-Policy header
	// Default: "same-origin"
	CrossOriginOpenerPolicy string

	// CrossOriginResourcePolicy sets Cross-Origin-Resource-Policy header
	// Default: "same-origin"
	CrossOriginResourcePolicy string
}

// DefaultHelmetConfig returns default security headers configuration
func DefaultHelmetConfig() HelmetConfig {
	return HelmetConfig{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		HSTSMaxAge:                31536000,
		HSTSIncludeSubdomains:     true,
		HSTSPreload:               false,
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
	}
}

// Helmet returns security headers middleware with default config
func Helmet() gin.HandlerFunc {
	return HelmetWithConfig(DefaultHelmetConfig())
}

// HelmetWithConfig returns security headers middleware with custom config
func HelmetWithConfig(cfg HelmetConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-XSS-Protection
		if cfg.XSSProtection != "" {
			c.Header("X-XSS-Protection", cfg.XSSProtection)
		}

		// X-Content-Type-Options
		if cfg.ContentTypeNosniff != "" {
			c.Header("X-Content-Type-Options", cfg.ContentTypeNosniff)
		}

		// X-Frame-Options
		if cfg.XFrameOptions != "" {
			c.Header("X-Frame-Options", cfg.XFrameOptions)
		}

		// Strict-Transport-Security (HSTS)
		if cfg.HSTSMaxAge > 0 {
			hsts := "max-age=" + itoa(cfg.HSTSMaxAge)
			if cfg.HSTSIncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			if cfg.HSTSPreload {
				hsts += "; preload"
			}
			c.Header("Strict-Transport-Security", hsts)
		}

		// Content-Security-Policy
		if cfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}

		// Referrer-Policy
		if cfg.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", cfg.ReferrerPolicy)
		}

		// Permissions-Policy
		if cfg.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", cfg.PermissionsPolicy)
		}

		// Cross-Origin-Embedder-Policy
		if cfg.CrossOriginEmbedderPolicy != "" {
			c.Header("Cross-Origin-Embedder-Policy", cfg.CrossOriginEmbedderPolicy)
		}

		// Cross-Origin-Opener-Policy
		if cfg.CrossOriginOpenerPolicy != "" {
			c.Header("Cross-Origin-Opener-Policy", cfg.CrossOriginOpenerPolicy)
		}

		// Cross-Origin-Resource-Policy
		if cfg.CrossOriginResourcePolicy != "" {
			c.Header("Cross-Origin-Resource-Policy", cfg.CrossOriginResourcePolicy)
		}

		c.Next()
	}
}

// itoa converts int to string (simple implementation)
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte
	negative := n < 0
	if negative {
		n = -n
	}

	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}
