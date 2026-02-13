package permission

import (
	"github.com/google/wire"
)

// ProviderSet is the provider set for the permission module
var ProviderSet = wire.NewSet(
	NewRepository,
	wire.Bind(new(Repository), new(*repository)),
	NewService,
	wire.Bind(new(Service), new(*service)),
	NewHandler,
)
