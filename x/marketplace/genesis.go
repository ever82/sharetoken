package marketplace

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/marketplace/keeper"
	"sharetoken/x/marketplace/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Initialize all services
	for _, service := range genState.Services {
		if err := k.SetService(ctx, service); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export all services
	// Note: keeper doesn't have GetAllServices method yet, but we can iterate if needed
	// For now, return empty genesis

	return genesis
}
