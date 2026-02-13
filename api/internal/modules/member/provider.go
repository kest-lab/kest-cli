package member

import (
	"github.com/google/wire"
)

// ProviderSet is the provider set for member module
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
