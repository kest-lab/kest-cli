package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger returns a gin.HandlerFunc that logs requests using the platform logger
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the log record
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Use the platform logger
		fields := map[string]any{
			"status":    statusCode,
			"latency":   latency.String(),
			"client_ip": clientIP,
			"method":    method,
			"path":      path,
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
			Error("HTTP Request Error", fields)
		} else {
			Info("HTTP Request", fields)
		}
	}
}
