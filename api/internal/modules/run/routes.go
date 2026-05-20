package run

import (
	"github.com/kest-labs/kest/api/internal/infra/router"
)

func (h *Handler) RegisterRoutes(r *router.Router) {
	r.Group("/workspaces/:id/runs", func(runs *router.Router) {
		runs.WithMiddleware("auth")

		runs.POST("", h.Create).
			Name("runs.create").
			WhereUUIDOrNumber("id")
		runs.GET("", h.List).
			Name("runs.list").
			WhereUUIDOrNumber("id")
		runs.GET("/:runid", h.Get).
			Name("runs.show").
			WhereUUIDOrNumber("id", "runid")
	})

	r.Group("/workspaces/:id/collections/:cid/requests/:rid", func(run *router.Router) {
		run.WithMiddleware("auth")

		run.POST("/run", h.Run).
			Name("run.execute").
			WhereUUIDOrNumber("id", "cid", "rid")
	})
}
