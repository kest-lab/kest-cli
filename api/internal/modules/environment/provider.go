package environment

import "github.com/google/wire"

// ProviderSet is the Wire provider set for environment module
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
