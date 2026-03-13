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
	gs.Tasks = k.GetAllTasks(ctx)

	// Export all applications
	gs.Applications = k.GetAllApplications(ctx)

	// Export all auctions
	gs.Auctions = k.GetAllAuctions(ctx)

	// Export all bids
	gs.Bids = k.GetAllBids(ctx)

	// Export all ratings
	gs.Ratings = k.GetAllRatings(ctx)

	// Export all reputations
	gs.Reputations = k.GetAllReputations(ctx)

	return gs
}
