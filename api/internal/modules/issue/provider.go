package issue

import "github.com/google/wire"

// ProviderSet is the Wire provider set for issue module
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
