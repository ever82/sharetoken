package taskmarket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/keeper"
	"sharetoken/x/taskmarket/types"
)

// InitGenesis initializes the taskmarket module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	// Initialize tasks
	for _, task := range gs.Tasks {
		k.SetTask(ctx, task)
	}

	// Initialize applications
	for _, app := range gs.Applications {
		k.SetApplication(ctx, app)
	}

	// Initialize auctions
	for _, auction := range gs.Auctions {
		k.SetAuction(ctx, auction)
	}

	// Initialize bids
	for _, bid := range gs.Bids {
		k.SetBid(ctx, bid)
	}

	// Initialize ratings
	for _, rating := range gs.Ratings {
		k.SetRating(ctx, rating)
	}

	// Initialize reputations
	for _, rep := range gs.Reputations {
		k.SetReputation(ctx, rep)
	}
}

// ExportGenesis returns the taskmarket module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	var gs types.GenesisState

	// Export all tasks
	allTasks := k.GetAllTasks()
	for _, task := range allTasks {
		gs.Tasks = append(gs.Tasks, *task)
	}

	// Export all applications
	gs.Applications = []types.Application{}
	// Note: In full implementation, iterate through all applications

	// Export all auctions
	gs.Auctions = []types.Auction{}
	// Note: In full implementation, iterate through all auctions

	// Export all bids
	gs.Bids = []types.Bid{}
	// Note: In full implementation, iterate through all bids

	// Export all ratings
	gs.Ratings = []types.Rating{}
	// Note: In full implementation, iterate through all ratings

	// Export all reputations
	gs.Reputations = []types.Reputation{}
	// Note: In full implementation, iterate through all reputations

	return gs
}
