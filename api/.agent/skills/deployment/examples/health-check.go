package infra

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHealthCheck adds health check routes to the engine
func RegisterHealthCheck(router *gin.RouterGroup, db *sql.DB) {
	h := &healthHandler{db: db}

	router.GET("/healthz", h.Liveness)
	router.GET("/readyz", h.Readiness)
}

type healthHandler struct {
	db *sql.DB
}

// Liveness check tells the platform if the app is alive
func (h *healthHandler) Liveness(c *gin.Context) {
	// If the server is running and can handle this request, it is alive
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

// Readiness check tells the platform if the app is ready to serve queries
func (h *healthHandler) Readiness(c *gin.Context) {
	// Check DB connectivity
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"error":  "database disconnected",
		})
		return
	}

	// Check other dependencies here (Redis, etc.)

	c.JSON(http.StatusOK, gin.H{"status": "READY"})
}
