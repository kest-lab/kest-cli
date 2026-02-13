package ingest

import (
	"github.com/google/wire"
)

// ProviderSet is the provider set for ingest module
var ProviderSet = wire.NewSet(
	NewHandler,
)
