package crowdfunding

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/crowdfunding/keeper"
	"sharetoken/x/crowdfunding/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Initialize ideas
	for _, idea := range genState.Ideas {
		ideaCopy := idea // Create a copy to avoid pointer issues
		if err := k.CreateIdea(&ideaCopy); err != nil {
			panic(err)
		}
	}

	// Initialize campaigns
	for _, campaign := range genState.Campaigns {
		campaignCopy := campaign // Create a copy to avoid pointer issues
		if err := k.CreateCampaign(&campaignCopy); err != nil {
			panic(err)
		}
	}

	// Initialize backers
	for _, backer := range genState.Backers {
		backerCopy := backer // Create a copy to avoid pointer issues
		if err := k.ContributeToCampaign(&backerCopy); err != nil {
			panic(err)
		}
	}

	// Initialize contributions
	for _, contribution := range genState.Contributions {
		// Note: keeper doesn't have a direct method to add contributions
		// This would need to be implemented in the keeper
		_ = contribution
	}

	// Initialize updates
	for _, update := range genState.Updates {
		updateCopy := update // Create a copy to avoid pointer issues
		if err := k.AddCampaignUpdate(&updateCopy); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Note: The keeper uses memory storage (maps)
	// In a real implementation with KVStore, we would iterate through all items
	// For now, return empty genesis

	return genesis
}
