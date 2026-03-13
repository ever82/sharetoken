package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/oracle/keeper"
	"sharetoken/x/oracle/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Initialize all prices
	for _, price := range genState.Prices {
		if err := k.SetPrice(ctx, price); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export all prices
	allPrices := k.GetAllPrices(ctx)
	genesis.Prices = allPrices

	return genesis
}
