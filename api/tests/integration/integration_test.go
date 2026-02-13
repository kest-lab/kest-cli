package integration

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup integration test environment
	// This could include setting up test database, redis, etc.

	// Run tests
	code := m.Run()

	// Cleanup

	os.Exit(code)
}
