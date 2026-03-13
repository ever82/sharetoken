package escrow

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/escrow/keeper"
	"sharetoken/x/escrow/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Initialize all escrows
	for _, escrow := range genState.Escrows {
		if err := k.SetEscrow(ctx, escrow); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export all escrows
	allEscrows := k.GetAllEscrows(ctx)
	genesis.Escrows = allEscrows

	return genesis
}
