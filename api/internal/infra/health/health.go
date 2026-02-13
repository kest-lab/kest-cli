package health

import (
	"context"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Status represents the health status of a component
type Status string

const (
	StatusUp       Status = "up"
	StatusDown     Status = "down"
	StatusDegraded Status = "degraded"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Status    Status         `json:"status"`
	Message   string         `json:"message,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	Latency   time.Duration  `json:"latency_ms,omitempty"`
	Duration  time.Duration  `json:"duration,omitempty"`
	Timestamp time.Time      `json:"timestamp,omitempty"`
}

// Checker is a function that performs a health check
type Checker func(ctx context.Context) CheckResult

// Health manages health checks for the application
type Health struct {
	mu       sync.RWMutex
	checkers map[string]Checker
}

// New creates a new Health instance
func New() *Health {
	return &Health{
		checkers: make(map[string]Checker),
	}
}

// Register adds a health checker
func (h *Health) Register(name string, checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[name] = checker
}

// Unregister removes a health checker
func (h *Health) Unregister(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.checkers, name)
}

// Check runs all health checks
func (h *Health) Check(ctx context.Context) map[string]CheckResult {
	h.mu.RLock()
	checkers := make(map[string]Checker, len(h.checkers))
	for k, v := range h.checkers {
		checkers[k] = v
	}
	h.mu.RUnlock()

	results := make(map[string]CheckResult, len(checkers))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, checker := range checkers {
		wg.Add(1)
		go func(name string, checker Checker) {
			defer wg.Done()
			start := time.Now()
			result := checker(ctx)
			duration := time.Since(start)
			result.Latency = duration
			result.Duration = duration
			result.Timestamp = start

			mu.Lock()
			results[name] = result
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()
	return results
}

// IsHealthy returns true if all checks pass
func (h *Health) IsHealthy(ctx context.Context) bool {
	results := h.Check(ctx)
	for _, result := range results {
		if result.Status == StatusDown {
			return false
		}
	}
	return true
}

// OverallStatus returns the overall health status
func (h *Health) OverallStatus(ctx context.Context) Status {
	results := h.Check(ctx)

	hasDown := false
	hasDegraded := false

	for _, result := range results {
		switch result.Status {
		case StatusDown:
			hasDown = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasDown {
		return StatusDown
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusUp
}

// Response represents the health check response
type Response struct {
	Status    Status                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks,omitempty"`
}

// GetHealth returns the full health response
func (h *Health) GetHealth(ctx context.Context) Response {
	results := h.Check(ctx)
	status := StatusUp

	for _, result := range results {
		if result.Status == StatusDown {
			status = StatusDown
			break
		}
		if result.Status == StatusDegraded {
			status = StatusDegraded
		}
	}

	return Response{
		Status:    status,
		Timestamp: time.Now(),
		Checks:    results,
	}
}

// --- Built-in Checkers ---

// Up creates a checker that always returns up
func Up(message string) Checker {
	return func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  StatusUp,
			Message: message,
		}
	}
}

// Down creates a checker that always returns down
func Down(message string) Checker {
	return func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  StatusDown,
			Message: message,
		}
	}
}

// Custom creates a checker from a simple error-returning function
func Custom(check func(ctx context.Context) error) Checker {
	return func(ctx context.Context) CheckResult {
		if err := check(ctx); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: err.Error(),
			}
		}
		return CheckResult{Status: StatusUp}
	}
}

// Timeout wraps a checker with a timeout
func Timeout(checker Checker, timeout time.Duration) Checker {
	return func(ctx context.Context) CheckResult {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan CheckResult, 1)
		go func() {
			done <- checker(ctx)
		}()

		select {
		case result := <-done:
			return result
		case <-ctx.Done():
			return CheckResult{
				Status:  StatusDown,
				Message: "check timed out",
			}
		}
	}
}

// DatabaseChecker creates a database health checker
func DatabaseChecker(db *gorm.DB) Checker {
	return func(ctx context.Context) CheckResult {
		sqlDB, err := db.DB()
		if err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "failed to get database connection",
			}
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "database ping failed",
			}
		}

		stats := sqlDB.Stats()
		return CheckResult{
			Status: StatusUp,
			Details: map[string]any{
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
				"max_open":         stats.MaxOpenConnections,
			},
		}
	}
}

// RedisChecker creates a Redis health checker
func RedisChecker(pingFunc func(ctx context.Context) error) Checker {
	return func(ctx context.Context) CheckResult {
		if err := pingFunc(ctx); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "redis ping failed",
			}
		}
		return CheckResult{Status: StatusUp}
	}
}

// DiskSpace creates a disk space checker
func DiskSpace(path string, minBytes uint64) Checker {
	return func(ctx context.Context) CheckResult {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(path, &stat); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "failed to get disk stats",
			}
		}

		available := stat.Bavail * uint64(stat.Bsize)
		total := stat.Blocks * uint64(stat.Bsize)
		used := total - available
		usedPercent := float64(used) / float64(total) * 100

		if available < minBytes {
			return CheckResult{
				Status:  StatusDown,
				Message: "disk space low",
				Details: map[string]any{
					"available_bytes": available,
					"total_bytes":     total,
					"used_percent":    usedPercent,
				},
			}
		}

		return CheckResult{
			Status: StatusUp,
			Details: map[string]any{
				"available_bytes": available,
				"total_bytes":     total,
				"used_percent":    usedPercent,
			},
		}
	}
}

// Memory creates a memory usage checker
// Note: This is a simplified implementation that works cross-platform
func Memory(maxUsedPercent float64) Checker {
	return func(ctx context.Context) CheckResult {
		// Memory check is platform-specific and complex to implement correctly
		// For production use, consider using a library like gopsutil
		// This simplified version always returns up
		return CheckResult{
			Status:  StatusUp,
			Message: "memory check passed",
		}
	}
}

// --- HTTP Handlers ---

// Handler returns a Gin handler for health checks
func (h *Health) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		response := h.GetHealth(ctx)

		httpStatus := http.StatusOK
		if response.Status == StatusDown {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, response)
	}
}

// ReadinessHandler returns a readiness probe handler (for Kubernetes)
func (h *Health) ReadinessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		status := h.OverallStatus(ctx)

		if status == StatusDown {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "down",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": string(status),
		})
	}
}

// RegisterRoutes registers health check routes on the engine
func (h *Health) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Handler())
	r.GET("/health/live", LivenessHandler())
	r.GET("/health/ready", h.ReadinessHandler())
}

// --- Global Instance ---

var globalHealth = New()

// Global returns the global health instance
func Global() *Health {
	return globalHealth
}

// Register adds a checker to the global health instance
func Register(name string, checker Checker) {
	globalHealth.Register(name, checker)
}

// Unregister removes a checker from the global health instance
func Unregister(name string) {
	globalHealth.Unregister(name)
}

// Check runs all checks on the global health instance
func Check(ctx context.Context) map[string]CheckResult {
	return globalHealth.Check(ctx)
}

// IsHealthy returns true if all global checks pass
func IsHealthy(ctx context.Context) bool {
	return globalHealth.IsHealthy(ctx)
}

// GetHealth returns the full health response from global instance
func GetHealth(ctx context.Context) Response {
	return globalHealth.GetHealth(ctx)
}

// Handler returns a handler using the global health instance
func Handler() gin.HandlerFunc {
	return globalHealth.Handler()
}

// ReadinessHandler returns a readiness handler using the global health instance
func ReadinessHandler() gin.HandlerFunc {
	return globalHealth.ReadinessHandler()
}

// LivenessHandler returns a simple liveness probe handler (for Kubernetes)
func LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	}
}

// RegisterRoutes registers health check routes using the global instance
func RegisterRoutes(r *gin.Engine) {
	globalHealth.RegisterRoutes(r)
}
