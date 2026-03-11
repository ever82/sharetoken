package llmcustody

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/keeper"
	"sharetoken/x/llmcustody/types"
)

// InitGenesis initializes the module's state from a genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Store all API keys
	for _, apiKey := range genState.APIKeys {
		// Don't restore keys with empty encrypted data (genesis export clears them)
		if len(apiKey.EncryptedKey) > 0 {
			if err := k.SetAPIKey(ctx, apiKey); err != nil {
				panic(err)
			}
		}
	}
}

// ExportGenesis exports the module's state to a genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export all API keys (with sensitive data cleared)
	allKeys := k.GetAllAPIKeys(ctx)
	for _, key := range allKeys {
		genesis.APIKeys = append(genesis.APIKeys, types.ExportAPIKeyForGenesis(key))
	}

	return genesis
}
