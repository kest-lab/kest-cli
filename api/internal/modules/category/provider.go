package category

import "github.com/google/wire"

// ProviderSet is the Wire provider set for category module
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
