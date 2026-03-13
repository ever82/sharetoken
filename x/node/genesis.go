package node

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/node/keeper"
	"sharetoken/x/node/types"
)

// InitGenesis initializes the node module's state from a genesis state.
func InitGenesis(ctx sdk.Context, k keeper.NodeKeeper, genState types.GenesisState) {
	// Initialize any state from genState if needed
	// For now, node uses file-based config so nothing to persist in KVStore
}

// ExportGenesis exports the node module's state to a genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.NodeKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export any state if needed
	// For now, node uses file-based config so nothing to export

	return genesis
}
