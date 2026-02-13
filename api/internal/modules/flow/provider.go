package flow

import "github.com/google/wire"

// ProviderSet is the Wire provider set for flow module
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
