package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// GetBidsByAuction returns bids by auction
func (k Keeper) GetBidsByAuction(ctx sdk.Context, auctionID string) []types.Bid {
	var bids []types.Bid
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.BidKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var bid types.Bid
		if err := k.cdc.Unmarshal(iterator.Value(), &bid); err != nil {
			continue
		}
		if bid.TaskId == auctionID {
			bids = append(bids, bid)
		}
	}
	return bids
}

// CreateAuction creates an auction
func (k Keeper) CreateAuction(ctx sdk.Context, taskID string, startingPrice, reservePrice uint64, duration int64) (*types.Auction, error) {
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	if task.TaskType != types.TaskType_TASK_TYPE_AUCTION {
		return nil, fmt.Errorf("task is not auction type")
	}

	auction := &types.Auction{
		TaskId:        taskID,
		StartingPrice: startingPrice,
		ReservePrice:  reservePrice,
		EndTime:       time.Now().Unix() + duration,
		IsActive:      true,
	}
	k.SetAuction(ctx, *auction)
	return auction, nil
}

// GetAuctionByTaskID gets an auction by task ID
func (k Keeper) GetAuctionByTaskID(ctx sdk.Context, taskID string) (types.Auction, bool) {
	return k.GetAuction(ctx, taskID)
}

// SubmitBid submits a bid
func (k Keeper) SubmitBid(ctx sdk.Context, bid types.Bid) error {
	if bid.TaskId == "" {
		return fmt.Errorf("task ID is required")
	}
	if bid.WorkerId == "" {
		return fmt.Errorf("worker ID is required")
	}
	if bid.Amount == 0 {
		return fmt.Errorf("bid amount must be greater than 0")
	}

	auction, found := k.GetAuction(ctx, bid.TaskId)
	if !found {
		return fmt.Errorf("auction not found: %s", bid.TaskId)
	}
	if !auction.IsActive {
		return fmt.Errorf("auction is not active")
	}
	if time.Now().Unix() > auction.EndTime {
		return fmt.Errorf("auction has ended")
	}

	bid.Status = types.BidStatus_BID_STATUS_PENDING
	k.SetBid(ctx, bid)

	// Update task bid count
	task, found := k.GetTask(ctx, bid.TaskId)
	if found {
		task.BidCount++
		k.SetTask(ctx, task)
	}

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

	auction.IsActive = false
	// Find winner (simplified - just get the last bid)
	bids := k.GetBidsByAuction(ctx, taskID)
	if len(bids) > 0 {
		winner := bids[len(bids)-1]
		task.WorkerId = winner.WorkerId
		task.Status = types.TaskStatus_TASK_STATUS_ASSIGNED
	}

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
	bid.Status = types.BidStatus_BID_STATUS_WITHDRAWN
	k.SetBid(ctx, bid)
	return nil
}
