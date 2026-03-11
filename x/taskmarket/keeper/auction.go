package keeper

import (
	"fmt"

	"sharetoken/x/taskmarket/types"
)

// CreateAuction creates an auction
func (lk *LegacyKeeper) CreateAuction(taskID string, startingPrice, reservePrice uint64, duration int64) (*types.Auction, error) {
	task := lk.GetTask(taskID)
	if task == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	if task.Type != types.TaskTypeAuction {
		return nil, fmt.Errorf("task is not auction type")
	}
	auction := types.NewAuction(taskID, startingPrice, reservePrice, duration)
	lk.auctions[taskID] = auction
	return auction, nil
}

// GetAuction gets an auction by task ID
func (lk *LegacyKeeper) GetAuction(taskID string) *types.Auction {
	return lk.auctions[taskID]
}

// SubmitBid submits a bid
func (lk *LegacyKeeper) SubmitBid(bid *types.Bid) error {
	if err := bid.Validate(); err != nil {
		return fmt.Errorf("invalid bid: %w", err)
	}
	auction := lk.GetAuction(bid.TaskID)
	if auction == nil {
		return fmt.Errorf("auction not found: %s", bid.TaskID)
	}
	if err := auction.AddBid(*bid); err != nil {
		return err
	}
	task := lk.GetTask(bid.TaskID)
	lk.auctions[bid.TaskID] = auction
	if task != nil {
		task.BidCount = auction.GetBidCount()
	}
	return nil
}

// CloseAuction closes an auction
func (lk *LegacyKeeper) CloseAuction(taskID string) error {
	auction := lk.GetAuction(taskID)
	if auction == nil {
		return fmt.Errorf("auction not found: %s", taskID)
	}
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Update status index
	lk.tasksByStatus[string(task.Status)] = removeFromSlice(lk.tasksByStatus[string(task.Status)], task.ID)

	winner, err := auction.CloseAuction()
	if err != nil {
		return err
	}
	auction.IsActive = false
	task.Assign(winner.WorkerID)

	// Update new status index
	lk.tasksByStatus[string(task.Status)] = append(lk.tasksByStatus[string(task.Status)], task.ID)
	// Update worker index
	lk.tasksByWorker[task.WorkerID] = append(lk.tasksByWorker[task.WorkerID], task.ID)
	return nil
}
