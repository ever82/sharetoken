package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ModuleNamer is an interface for types that can return their module name
type ModuleNamer interface {
	// ModuleName returns the name of the module
	ModuleName() string
}

// Logger returns a module-specific logger.
// This is a generic helper that can be used by any keeper that implements ModuleNamer.
func Logger(ctx sdk.Context, mn ModuleNamer) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", mn.ModuleName()))
}

// LoggerFunc is a function type that returns a logger for a given context
type LoggerFunc func(ctx sdk.Context) log.Logger

// NewLoggerFunc creates a LoggerFunc for the given module name
func NewLoggerFunc(moduleName string) LoggerFunc {
	return func(ctx sdk.Context) log.Logger {
		return ctx.Logger().With("module", fmt.Sprintf("x/%s", moduleName))
	}
}
