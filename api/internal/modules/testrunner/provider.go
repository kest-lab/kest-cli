package testrunner

import "github.com/google/wire"

// ProviderSet is the Wire provider set for test runner module
var ProviderSet = wire.NewSet(
	NewExecutor,
)
