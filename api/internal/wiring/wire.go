//go:build wireinject
// +build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/zgiai/kest-api/internal/app"
	"github.com/zgiai/kest-api/internal/infra"
	"github.com/zgiai/kest-api/internal/modules/apispec"
	"github.com/zgiai/kest-api/internal/modules/category"
	"github.com/zgiai/kest-api/internal/modules/environment"
	"github.com/zgiai/kest-api/internal/modules/event"
	"github.com/zgiai/kest-api/internal/modules/flow"
	"github.com/zgiai/kest-api/internal/modules/ingest"
	"github.com/zgiai/kest-api/internal/modules/issue"
	"github.com/zgiai/kest-api/internal/modules/member"
	"github.com/zgiai/kest-api/internal/modules/permission"
	"github.com/zgiai/kest-api/internal/modules/project"
	"github.com/zgiai/kest-api/internal/modules/system"
	"github.com/zgiai/kest-api/internal/modules/testcase"
	"github.com/zgiai/kest-api/internal/modules/testrunner"
	"github.com/zgiai/kest-api/internal/modules/user"
)

// InitApplication initializes the entire application with all dependencies.
// This is the single entry point for Wire DI.
func InitApplication() (*app.Application, error) {
	wire.Build(
		// Infrastructure providers
		infra.ProviderSet,

		// Module providers
		user.ProviderSet,
		member.ProviderSet,
		permission.ProviderSet,
		project.ProviderSet,
		apispec.ProviderSet,
		category.ProviderSet,
		environment.ProviderSet,
		flow.ProviderSet,
		event.ProviderSet,
		issue.ProviderSet,
		ingest.ProviderSet,
		testcase.ProviderSet,
		testrunner.ProviderSet,
		system.ProviderSet,

		// Bind event.Service to ingest.EventProcessor
		wire.Bind(new(ingest.EventProcessor), new(event.Service)),

		// Bind testrunner.Executor to testcase.Runner
		wire.Bind(new(testcase.Runner), new(*testrunner.Executor)),

		// Aggregate handlers
		wire.Struct(new(app.Handlers), "*"),

		// Build final application
		wire.Struct(new(app.Application), "*"),
	)
	return nil, nil
}
