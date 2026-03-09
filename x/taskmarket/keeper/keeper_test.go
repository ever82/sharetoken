package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/taskmarket/types"
)

func TestTaskCreation(t *testing.T) {
	k := NewKeeper()

	task := types.NewTask("task-1", "Build Website", "Create a responsive website", "requester-1", types.TaskTypeOpen, 1000)
	task.AddMilestone(types.Milestone{
		ID:    "ms-1",
		Title: "Design",
		Amount: 500,
		Order: 1,
		Status: types.MilestoneStatusPending,
	})
	task.AddMilestone(types.Milestone{
		ID:    "ms-2",
		Title: "Development",
		Amount: 500,
		Order: 2,
		Status: types.MilestoneStatusPending,
	})

	err := k.CreateTask(task)
	require.NoError(t, err)

	// Retrieve task
	retrieved := k.GetTask("task-1")
	require.NotNil(t, retrieved)
	require.Equal(t, "Build Website", retrieved.Title)
	require.Equal(t, uint64(1000), retrieved.Budget)
	require.Len(t, retrieved.Milestones, 2)
}

func TestTaskValidation(t *testing.T) {
	k := NewKeeper()

	// Invalid - no title
	invalidTask := types.NewTask("task-1", "", "Description", "requester-1", types.TaskTypeOpen, 1000)
	err := k.CreateTask(invalidTask)
	require.Error(t, err)

	// Invalid - no budget
	invalidTask2 := types.NewTask("task-2", "Title", "Description", "requester-1", types.TaskTypeOpen, 0)
	err = k.CreateTask(invalidTask2)
	require.Error(t, err)
}

func TestTaskLifecycle(t *testing.T) {
	k := NewKeeper()

	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	err := k.CreateTask(task)
	require.NoError(t, err)

	// Publish
	task.Publish()
	require.Equal(t, types.TaskStatusOpen, task.Status)

	// Assign
	task.Assign("worker-1")
	require.Equal(t, types.TaskStatusAssigned, task.Status)
	require.Equal(t, "worker-1", task.WorkerID)

	// Start
	err = k.StartTask("task-1")
	require.NoError(t, err)
	require.Equal(t, types.TaskStatusInProgress, task.Status)

	// Complete
	task.Complete()
	require.Equal(t, types.TaskStatusCompleted, task.Status)
	require.Greater(t, task.CompletedAt, int64(0))
}

func TestMilestoneWorkflow(t *testing.T) {
	k := NewKeeper()

	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task.AddMilestone(types.Milestone{
		ID:     "ms-1",
		Title:  "Design",
		Amount: 500,
		Order:  1,
		Status: types.MilestoneStatusPending,
	})
	task.AddMilestone(types.Milestone{
		ID:     "ms-2",
		Title:  "Development",
		Amount: 500,
		Order:  2,
		Status: types.MilestoneStatusPending,
	})

	err := k.CreateTask(task)
	require.NoError(t, err)

	task.Publish()
	task.Assign("worker-1")
	k.StartTask("task-1")

	// Check first milestone is active
	require.Equal(t, types.MilestoneStatusActive, task.Milestones[0].Status)

	// Submit first milestone
	err = k.SubmitMilestone("task-1", "ms-1", "Design files attached")
	require.NoError(t, err)
	require.Equal(t, types.MilestoneStatusSubmitted, task.Milestones[0].Status)

	// Approve first milestone
	err = k.ApproveMilestone("task-1", "ms-1")
	require.NoError(t, err)
	require.Equal(t, types.MilestoneStatusApproved, task.Milestones[0].Status)

	// Check second milestone is now active
	require.Equal(t, types.MilestoneStatusActive, task.Milestones[1].Status)

	// Task should not be complete yet
	require.False(t, task.AllMilestonesCompleted())

	// Submit and approve second milestone
	err = k.SubmitMilestone("task-1", "ms-2", "Code committed")
	require.NoError(t, err)
	err = k.ApproveMilestone("task-1", "ms-2")
	require.NoError(t, err)

	// All milestones completed
	require.True(t, task.AllMilestonesCompleted())
	require.Equal(t, types.TaskStatusCompleted, task.Status)
}

func TestMilestoneValidation(t *testing.T) {
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)

	// Valid milestones
	task.AddMilestone(types.Milestone{
		ID:     "ms-1",
		Amount: 400,
	})
	task.AddMilestone(types.Milestone{
		ID:     "ms-2",
		Amount: 600,
	})

	err := task.ValidateMilestones()
	require.NoError(t, err)

	// Invalid - doesn't match budget
	task2 := types.NewTask("task-2", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task2.AddMilestone(types.Milestone{
		ID:     "ms-1",
		Amount: 500,
	})

	err = task2.ValidateMilestones()
	require.Error(t, err)
}

func TestOpenTaskApplication(t *testing.T) {
	k := NewKeeper()

	// Create open task
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	err := k.CreateTask(task)
	require.NoError(t, err)
	task.Publish()

	// Submit application
	app := types.NewApplication("app-1", "task-1", "worker-1", 900)
	app.CoverLetter = "I have 5 years experience"
	err = k.SubmitApplication(app)
	require.NoError(t, err)

	// Check application count
	require.Equal(t, 1, task.ApplicationCount)

	// Get applications
	apps := k.GetApplicationsByTask("task-1")
	require.Len(t, apps, 1)
	require.Equal(t, "worker-1", apps[0].WorkerID)

	// Accept application
	err = k.AcceptApplication("app-1")
	require.NoError(t, err)
	require.Equal(t, types.ApplicationStatusAccepted, app.Status)
	require.Equal(t, "worker-1", task.WorkerID)
}

func TestApplicationValidation(t *testing.T) {
	k := NewKeeper()

	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	k.CreateTask(task)
	task.Publish()

	// Invalid - auction task
	auctionTask := types.NewTask("task-2", "Build App", "Description", "requester-1", types.TaskTypeAuction, 2000)
	k.CreateTask(auctionTask)
	auctionTask.Publish()

	app := types.NewApplication("app-1", "task-2", "worker-1", 1800)
	err := k.SubmitApplication(app)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not open type")

	// Invalid - no price
	app2 := types.NewApplication("app-2", "task-1", "worker-1", 0)
	err = k.SubmitApplication(app2)
	require.Error(t, err)
}

func TestAuctionBidding(t *testing.T) {
	k := NewKeeper()

	// Create auction task
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeAuction, 1000)
	err := k.CreateTask(task)
	require.NoError(t, err)
	task.Publish()

	// Create auction
	auction, err := k.CreateAuction("task-1", 1000, 800, 86400)
	require.NoError(t, err)
	require.NotNil(t, auction)
	require.Equal(t, uint64(1000), auction.StartingPrice)
	require.Equal(t, uint64(800), auction.ReservePrice)

	// Submit bids
	bid1 := types.NewBid("bid-1", "task-1", "worker-1", 900)
	err = k.SubmitBid(bid1)
	require.NoError(t, err)
	require.Equal(t, types.BidStatusPending, bid1.Status)

	bid2 := types.NewBid("bid-2", "task-1", "worker-2", 850)
	err = k.SubmitBid(bid2)
	require.NoError(t, err)

	// Refresh auction to get updated bids
	auction = k.GetAuction("task-1")

	// Check that first bid is now outbid
	for _, b := range auction.Bids {
		if b.ID == "bid-1" {
			require.Equal(t, types.BidStatusOutbid, b.Status)
			break
		}
	}

	// Check winning bid
	winner := auction.GetWinningBid()
	require.NotNil(t, winner)
	require.Equal(t, "worker-2", winner.WorkerID)
	require.Equal(t, uint64(850), winner.Amount)
}

func TestBidValidation(t *testing.T) {
	k := NewKeeper()

	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeAuction, 1000)
	k.CreateTask(task)
	task.Publish()

	auction, _ := k.CreateAuction("task-1", 1000, 800, 86400)
	require.NotNil(t, auction)

	// Invalid - exceeds starting price
	bid1 := types.NewBid("bid-1", "task-1", "worker-1", 1100)
	err := k.SubmitBid(bid1)
	require.Error(t, err)

	// Valid bid
	bid2 := types.NewBid("bid-2", "task-1", "worker-1", 900)
	err = k.SubmitBid(bid2)
	require.NoError(t, err)
}

func TestRatingSystem(t *testing.T) {
	k := NewKeeper()

	// Create and complete task
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	k.CreateTask(task)
	task.Publish()
	task.Assign("worker-1")
	task.Complete()

	// Submit rating
	rating := types.NewRating("rating-1", "task-1", "requester-1", "worker-1")
	rating.SetRating(types.DimensionQuality, 5)
	rating.SetRating(types.DimensionCommunication, 4)
	rating.SetRating(types.DimensionTimeliness, 5)
	rating.SetRating(types.DimensionProfessionalism, 4)
	rating.Comment = "Great work!"

	err := k.SubmitRating(rating)
	require.NoError(t, err)

	// Check rating
	require.Equal(t, 4.5, rating.GetAverage())
	require.True(t, rating.IsComplete())

	// Check reputation
	rep := k.GetReputation("worker-1")
	require.Equal(t, 1, rep.TotalRatings)
	require.Equal(t, 4.5, rep.AverageRating)
	require.Equal(t, float64(5), rep.GetRatingForDimension(types.DimensionQuality))
	require.Equal(t, float64(4), rep.GetRatingForDimension(types.DimensionCommunication))
}

func TestRatingValidation(t *testing.T) {
	k := NewKeeper()

	// Try to rate incomplete task
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	k.CreateTask(task)
	task.Publish()
	task.Assign("worker-1")
	// Not completed

	rating := types.NewRating("rating-1", "task-1", "requester-1", "worker-1")
	rating.SetRating(types.DimensionQuality, 5)

	err := k.SubmitRating(rating)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not completed")
}

func TestReputationTrustLevel(t *testing.T) {
	rep := types.NewReputation("user-1")

	// New user
	require.True(t, rep.IsNew())
	require.Equal(t, "new", rep.GetTrustLevel())

	// Add ratings to become established
	for i := 0; i < 3; i++ {
		rating := types.NewRating(fmt.Sprintf("rating-%d", i), fmt.Sprintf("task-%d", i), "other", "user-1")
		rating.SetRating(types.DimensionQuality, 4)
		rating.SetRating(types.DimensionCommunication, 4)
		rating.SetRating(types.DimensionTimeliness, 4)
		rating.SetRating(types.DimensionProfessionalism, 4)
		rep.AddRating(rating)
	}

	require.False(t, rep.IsNew())
	require.Equal(t, "established", rep.GetTrustLevel())

	// Add more high ratings to become trusted
	for i := 3; i < 6; i++ {
		rating := types.NewRating(fmt.Sprintf("rating-%d", i), fmt.Sprintf("task-%d", i), "other", "user-1")
		rating.SetRating(types.DimensionQuality, 5)
		rating.SetRating(types.DimensionCommunication, 5)
		rating.SetRating(types.DimensionTimeliness, 5)
		rating.SetRating(types.DimensionProfessionalism, 5)
		rep.AddRating(rating)
	}

	require.True(t, rep.IsTrusted())
	require.Equal(t, "trusted", rep.GetTrustLevel())
}

func TestTaskStatistics(t *testing.T) {
	k := NewKeeper()

	// Create tasks
	task1 := types.NewTask("task-1", "Task 1", "Description", "requester-1", types.TaskTypeOpen, 1000)
	k.CreateTask(task1)
	task1.Publish()

	task2 := types.NewTask("task-2", "Task 2", "Description", "requester-1", types.TaskTypeOpen, 2000)
	k.CreateTask(task2)
	task2.Publish()
	task2.Assign("worker-1")

	task3 := types.NewTask("task-3", "Task 3", "Description", "requester-1", types.TaskTypeAuction, 1500)
	k.CreateTask(task3)
	task3.Publish()
	k.CreateAuction("task-3", 1500, 1200, 86400)

	// Submit bid
	bid := types.NewBid("bid-1", "task-3", "worker-1", 1300)
	k.SubmitBid(bid)

	// Get statistics
	stats := k.GetTaskStatistics()
	require.Equal(t, 3, stats["total_tasks"])
	require.Equal(t, 2, stats["open_tasks"]) // task-1 and task-3
	require.Equal(t, 1, stats["assigned_tasks"]) // task-2
	require.GreaterOrEqual(t, stats["total_bids"].(int), 1)
}

func TestGetOpenTasks(t *testing.T) {
	k := NewKeeper()

	// Create open tasks
	task1 := types.NewTask("task-1", "Task 1", "Description", "requester-1", types.TaskTypeOpen, 1000)
	k.CreateTask(task1)
	task1.Publish()

	task2 := types.NewTask("task-2", "Task 2", "Description", "requester-1", types.TaskTypeOpen, 2000)
	k.CreateTask(task2)
	task2.Publish()

	// Create assigned task
	task3 := types.NewTask("task-3", "Task 3", "Description", "requester-1", types.TaskTypeOpen, 1500)
	k.CreateTask(task3)
	task3.Publish()
	task3.Assign("worker-1")

	openTasks := k.GetOpenTasks()
	require.Len(t, openTasks, 2)
}
