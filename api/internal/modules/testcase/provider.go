package testcase

import "github.com/google/wire"

// ProviderSet is the Wire provider set for test case module
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
