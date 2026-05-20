package run

import (
	"github.com/google/wire"

	"github.com/kest-labs/kest/api/internal/contracts"
)

var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
	wire.Bind(new(contracts.Module), new(*Handler)),
)
