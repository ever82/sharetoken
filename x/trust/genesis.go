package trust

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/trust/keeper"
	"sharetoken/x/trust/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.MQKeeper, genState types.GenesisState) {
	// Note: The keeper currently uses memory storage (MQKeeper)
	// In the future, this should be migrated to use KVStore

	// Initialize MQ scores
	for _, mqScore := range genState.MqScores {
		// Initialize score for each address
		k.InitializeScore(mqScore.Address)
		// Note: Full restoration of MQ scores would require additional methods in the keeper
		_ = mqScore
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *keeper.MQKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export all MQ scores
	allScores := k.GetAllScores()
	for _, score := range allScores {
		if score != nil {
			genesis.MqScores = append(genesis.MqScores, *score)
		}
	}

	return genesis
}
