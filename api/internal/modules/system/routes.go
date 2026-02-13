package system

import (
	"github.com/zgiai/kest-api/internal/infra/router"
)

// RegisterRoutes registers system routes
func RegisterRoutes(r *router.Router, handler *Handler) {
	// Public routes - no authentication required
	r.GET("/system-features", handler.GetSystemFeatures)
	r.GET("/setup-status", handler.GetSetupStatus)
}
