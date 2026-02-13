package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// CompressConfig holds Compress middleware configuration
type CompressConfig struct {
	// Level is the compression level (1-9, where 9 is best compression)
	// Default: gzip.DefaultCompression
	Level int

	// MinLength is the minimum response size to compress
	// Default: 1024 bytes
	MinLength int

	// ExcludedExtensions are file extensions to skip compression
	// Default: [".png", ".gif", ".jpeg", ".jpg", ".webp", ".ico"]
	ExcludedExtensions []string

	// ExcludedPaths are paths to skip compression
	ExcludedPaths []string
}

// DefaultCompressConfig returns default compression configuration
func DefaultCompressConfig() CompressConfig {
	return CompressConfig{
		Level:              gzip.DefaultCompression,
		MinLength:          1024,
		ExcludedExtensions: []string{".png", ".gif", ".jpeg", ".jpg", ".webp", ".ico", ".woff", ".woff2"},
		ExcludedPaths:      []string{},
	}
}

// gzipWriter wraps http.ResponseWriter with gzip compression
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

// gzip writer pool
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
		return w
	},
}

// Compress returns Compress middleware with default config
func Compress() gin.HandlerFunc {
	return CompressWithConfig(DefaultCompressConfig())
}

// CompressWithConfig returns Compress middleware with custom config
func CompressWithConfig(cfg CompressConfig) gin.HandlerFunc {
	// Build excluded extensions map
	excludedExt := make(map[string]bool)
	for _, ext := range cfg.ExcludedExtensions {
		excludedExt[ext] = true
	}

	return func(c *gin.Context) {
		// Check if client accepts gzip
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Check excluded paths
		path := c.Request.URL.Path
		for _, excluded := range cfg.ExcludedPaths {
			if strings.HasPrefix(path, excluded) {
				c.Next()
				return
			}
		}

		// Check excluded extensions
		for ext := range excludedExt {
			if strings.HasSuffix(path, ext) {
				c.Next()
				return
			}
		}

		// Get gzip writer from pool
		gz := gzipWriterPool.Get().(*gzip.Writer)
		gz.Reset(c.Writer)

		// Create wrapped writer
		gzw := &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		// Set headers
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Replace writer
		c.Writer = gzw

		defer func() {
			gz.Close()
			gzipWriterPool.Put(gz)
		}()

		c.Next()

		// Remove Content-Length as it's now different
		c.Header("Content-Length", "")
	}
}

// NoCompress skips compression for a specific handler
func NoCompress(c *gin.Context) {
	c.Header("Content-Encoding", "identity")
}

// compressedResponseWriter implements http.ResponseWriter
var _ http.ResponseWriter = (*gzipWriter)(nil)
