package keeper

import (
	"fmt"

	"sharetoken/x/taskmarket/types"
)

// SubmitRating submits a rating
func (lk *LegacyKeeper) SubmitRating(rating *types.Rating) error {
	if err := rating.Validate(); err != nil {
		return fmt.Errorf("invalid rating: %w", err)
	}
	task := lk.GetTask(rating.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", rating.TaskID)
	}
	if task.Status != types.TaskStatusCompleted {
		return fmt.Errorf("task is not completed")
	}
	lk.ratings[rating.ID] = rating
	rep, exists := lk.reputations[rating.RatedID]
	if !exists {
		rep = types.NewReputation(rating.RatedID)
	}
	rep.AddRating(rating)
	lk.reputations[rating.RatedID] = rep
	return nil
}

// GetRating gets a rating by ID
func (lk *LegacyKeeper) GetRating(id string) *types.Rating {
	return lk.ratings[id]
}

// GetReputation gets reputation for a user
func (lk *LegacyKeeper) GetReputation(userID string) *types.Reputation {
	rep, exists := lk.reputations[userID]
	if !exists {
		return types.NewReputation(userID)
	}
	return rep
}
