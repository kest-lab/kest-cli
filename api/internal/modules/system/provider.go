package system

import "github.com/google/wire"

// ProviderSet is the Wire provider set for system module
var ProviderSet = wire.NewSet(
	NewHandler,
)
