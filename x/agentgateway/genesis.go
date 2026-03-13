package agentgateway

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/agentgateway/keeper"
	"sharetoken/x/agentgateway/types"
)

// InitGenesis initializes the agentgateway module's state from a genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Initialize any state from genState if needed
	// For now, sessions are stored in memory so nothing to persist
}

// ExportGenesis exports the agentgateway module's state to a genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export any state if needed
	// For now, sessions are stored in memory so nothing to export

	return genesis
}
