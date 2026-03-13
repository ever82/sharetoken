package dispute

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/dispute/keeper"
	"sharetoken/x/dispute/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.DisputeKeeper, genState types.GenesisState) {
	// Note: keeper currently uses memory storage (DisputeKeeper)
	// In the future, this should be migrated to use KVStore

	// Initialize disputes
	for _, dispute := range genState.Disputes {
		// Since the keeper uses memory storage, we can't directly restore state
		// This is a placeholder for when the keeper is migrated to KVStore
		_ = dispute
	}

	// Initialize juror pool
	for _, juror := range genState.JurorPool {
		// Placeholder for when the keeper is migrated to KVStore
		_ = juror
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *keeper.DisputeKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export all disputes from the keeper
	// Note: GetAllDisputes returns []*types.Dispute
	allDisputes := k.GetAllDisputes()
	for _, dispute := range allDisputes {
		if dispute != nil {
			genesis.Disputes = append(genesis.Disputes, *dispute)
		}
	}

	// Export juror pool
	// Note: juror pool access would need to be added to the keeper interface

	return genesis
}
