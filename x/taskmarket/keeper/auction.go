package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// CreateAuction creates an auction
func (k Keeper) CreateAuction(ctx sdk.Context, taskID string, startingPrice, reservePrice uint64, duration int64) (*types.Auction, error) {
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	if task.Type != types.TaskTypeAuction {
		return nil, fmt.Errorf("task is not auction type")
	}
	auction := types.NewAuction(taskID, startingPrice, reservePrice, duration)
	k.SetAuction(ctx, *auction)
	return auction, nil
}

// GetAuction gets an auction by task ID
func (k Keeper) GetAuctionByTaskID(ctx sdk.Context, taskID string) (types.Auction, bool) {
	return k.GetAuction(ctx, taskID)
}

// SubmitBid submits a bid
func (k Keeper) SubmitBid(ctx sdk.Context, bid types.Bid) error {
	if err := bid.Validate(); err != nil {
		return fmt.Errorf("invalid bid: %w", err)
	}
	auction, found := k.GetAuction(ctx, bid.TaskID)
	if !found {
		return fmt.Errorf("auction not found: %s", bid.TaskID)
	}
	if err := auction.AddBid(bid); err != nil {
		return err
	}

	// Update task bid count
	task, found := k.GetTask(ctx, bid.TaskID)
	if found {
		task.BidCount = auction.GetBidCount()
		k.SetTask(ctx, task)
	}

	k.SetAuction(ctx, auction)
	k.SetBid(ctx, bid)
	return nil
}

// CloseAuction closes an auction
func (k Keeper) CloseAuction(ctx sdk.Context, taskID string) error {
	auction, found := k.GetAuction(ctx, taskID)
	if !found {
		return fmt.Errorf("auction not found: %s", taskID)
	}
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Delete old indexes
	k.deleteTaskIndexes(ctx, task)

	winner, err := auction.CloseAuction()
	if err != nil {
		return err
	}
	auction.IsActive = false
	task.Assign(winner.WorkerID)

	// Update new indexes
	k.setTaskByWorker(ctx, task)
	k.setTaskByStatus(ctx, task)
	k.SetTask(ctx, task)
	k.SetAuction(ctx, auction)

	return nil
}

// CancelBid cancels a bid
func (k Keeper) CancelBid(ctx sdk.Context, bidID string) error {
	bid, found := k.GetBid(ctx, bidID)
	if !found {
		return fmt.Errorf("bid not found: %s", bidID)
	}
	bid.Withdraw()
	k.SetBid(ctx, bid)
	return nil
}
