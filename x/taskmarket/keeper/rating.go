package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// SubmitRating submits a rating
func (k Keeper) SubmitRating(ctx sdk.Context, rating types.Rating) error {
	if err := rating.Validate(); err != nil {
		return fmt.Errorf("invalid rating: %w", err)
	}
	task, found := k.GetTask(ctx, rating.TaskId)
	if !found {
		return fmt.Errorf("task not found: %s", rating.TaskId)
	}
	if task.Status != types.TaskStatusCompleted {
		return fmt.Errorf("task is not completed")
	}

	// Check if rater has already rated this task
	existingRatings := k.GetRatingsByTask(ctx, rating.TaskId)
	for _, existing := range existingRatings {
		if existing.RaterId == rating.RaterId {
			return fmt.Errorf("rater has already submitted a rating for this task")
		}
	}

	k.SetRating(ctx, rating)

	// Update reputation
	rep, found := k.GetReputation(ctx, rating.RatedId)
	if !found {
		rep = *types.NewReputation(rating.RatedId)
	}
	rep.AddRating(&rating)
	k.SetReputation(ctx, rep)

	return nil
}

// GetReputationByUserID gets reputation for a user (alias for GetReputation for consistency)
func (k Keeper) GetReputationByUserID(ctx sdk.Context, userID string) *types.Reputation {
	rep, found := k.GetReputation(ctx, userID)
	if !found {
		return types.NewReputation(userID)
	}
	return &rep
}

// UpdateReputation updates a user's reputation
func (k Keeper) UpdateReputation(ctx sdk.Context, rep types.Reputation) {
	k.SetReputation(ctx, rep)
}
