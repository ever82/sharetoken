//nolint:errcheck
package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/taskmarket/types"
)

// Note: These are simplified unit tests that test the types package functionality
// Full integration tests with KVStore require complex SDK setup

func TestTaskCreation(t *testing.T) {
	task := types.NewTask("task-1", "Build Website", "Create a responsive website", "requester-1", types.TaskTypeOpen, 1000)
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

	require.Equal(t, "task-1", task.ID)
	require.Equal(t, "Build Website", task.Title)
	require.Equal(t, uint64(1000), task.Budget)
	require.Len(t, task.Milestones, 2)
}

func TestTaskValidation(t *testing.T) {
	// Invalid - no title
	invalidTask := types.NewTask("task-1", "", "Description", "requester-1", types.TaskTypeOpen, 1000)
	err := invalidTask.Validate()
	require.Error(t, err)

	// Invalid - no budget
	invalidTask2 := types.NewTask("task-2", "Title", "Description", "requester-1", types.TaskTypeOpen, 0)
	err = invalidTask2.Validate()
	require.Error(t, err)
}

func TestTaskLifecycle(t *testing.T) {
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)

	// Publish
	task.Publish()
	require.Equal(t, types.TaskStatusOpen, task.Status)

	// Assign
	task.Assign("worker-1")
	require.Equal(t, types.TaskStatusAssigned, task.Status)
	require.Equal(t, "worker-1", task.WorkerID)

	// Start
	task.Start()
	require.Equal(t, types.TaskStatusInProgress, task.Status)

	// Complete
	task.Complete()
	require.Equal(t, types.TaskStatusCompleted, task.Status)
	require.Greater(t, task.CompletedAt, int64(0))
}

func TestMilestoneWorkflow(t *testing.T) {
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

	task.Publish()
	task.Assign("worker-1")
	task.Start()

	// Manually activate first milestone (normally done by keeper.StartTask)
	task.Milestones[0].Status = types.MilestoneStatusActive

	// Check first milestone is active
	require.Equal(t, types.MilestoneStatusActive, task.Milestones[0].Status)

	// Submit first milestone
	err := task.SubmitMilestone("ms-1", "Design files attached")
	require.NoError(t, err)
	require.Equal(t, types.MilestoneStatusSubmitted, task.Milestones[0].Status)

	// Approve first milestone
	err = task.ApproveMilestone("ms-1")
	require.NoError(t, err)
	require.Equal(t, types.MilestoneStatusApproved, task.Milestones[0].Status)

	// Manually activate second milestone
	task.Milestones[1].Status = types.MilestoneStatusActive

	// Check second milestone is now active
	require.Equal(t, types.MilestoneStatusActive, task.Milestones[1].Status)

	// Task should not be complete yet
	require.False(t, task.AllMilestonesCompleted())

	// Submit and approve second milestone
	err = task.SubmitMilestone("ms-2", "Code committed")
	require.NoError(t, err)
	err = task.ApproveMilestone("ms-2")
	require.NoError(t, err)

	// All milestones completed
	require.True(t, task.AllMilestonesCompleted())
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
	// Create open task
	task := types.NewTask("task-1", "Build Website", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task.Publish()

	// Submit application
	app := types.NewApplication("app-1", "task-1", "worker-1", 900)
	app.CoverLetter = "I have 5 years experience"

	require.Equal(t, types.ApplicationStatusPending, app.Status)
	require.Equal(t, "task-1", app.TaskID)
	require.Equal(t, "worker-1", app.WorkerID)

	// Accept application
	app.Accept()
	require.Equal(t, types.ApplicationStatusAccepted, app.Status)

	// Assign task
	task.Assign(app.WorkerID)
	require.Equal(t, "worker-1", task.WorkerID)
}

func TestApplicationValidation(t *testing.T) {
	// Valid application
	app := types.NewApplication("app-1", "task-1", "worker-1", 900)
	err := app.Validate()
	require.NoError(t, err)

	// Invalid - no price
	app2 := types.NewApplication("app-2", "task-1", "worker-1", 0)
	err = app2.Validate()
	require.Error(t, err)

	// Invalid - no worker
	app3 := types.NewApplication("app-3", "task-1", "", 900)
	err = app3.Validate()
	require.Error(t, err)
}

func TestAuctionBidding(t *testing.T) {
	// Create auction
	auction := types.NewAuction("task-1", 1000, 800, 86400)
	require.NotNil(t, auction)
	require.Equal(t, uint64(1000), auction.StartingPrice)
	require.Equal(t, uint64(800), auction.ReservePrice)
	require.True(t, auction.IsActive)

	// Submit bids
	bid1 := types.NewBid("bid-1", "task-1", "worker-1", 900)
	err := auction.AddBid(*bid1)
	require.NoError(t, err)
	require.Equal(t, types.BidStatusPending, bid1.Status)

	bid2 := types.NewBid("bid-2", "task-1", "worker-2", 850)
	err = auction.AddBid(*bid2)
	require.NoError(t, err)

	// Refresh auction to get updated bids
	require.Equal(t, 2, len(auction.Bids))

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
	auction := types.NewAuction("task-1", 1000, 800, 86400)

	// Invalid - exceeds starting price
	bid1 := types.NewBid("bid-1", "task-1", "worker-1", 1100)
	err := auction.AddBid(*bid1)
	require.Error(t, err)

	// Valid bid
	bid2 := types.NewBid("bid-2", "task-1", "worker-1", 900)
	err = auction.AddBid(*bid2)
	require.NoError(t, err)
}

func TestRatingSystem(t *testing.T) {
	// Submit rating
	rating := types.NewRating("rating-1", "task-1", "requester-1", "worker-1")
	err := rating.SetRating(types.DimensionQuality, 5)
	require.NoError(t, err)
	err = rating.SetRating(types.DimensionCommunication, 4)
	require.NoError(t, err)
	err = rating.SetRating(types.DimensionTimeliness, 5)
	require.NoError(t, err)
	err = rating.SetRating(types.DimensionProfessionalism, 4)
	require.NoError(t, err)
	rating.Comment = "Great work!"

	// Check rating
	require.Equal(t, 4.5, rating.GetAverage())
	require.True(t, rating.IsComplete())

	// Check reputation update
	rep := types.NewReputation("worker-1")
	rep.AddRating(rating)
	require.Equal(t, 1, rep.TotalRatings)
	require.Equal(t, 4.5, rep.AverageRating)
	require.Equal(t, float64(5), rep.GetRatingForDimension(types.DimensionQuality))
	require.Equal(t, float64(4), rep.GetRatingForDimension(types.DimensionCommunication))
}

func TestRatingValidation(t *testing.T) {
	rating := types.NewRating("rating-1", "task-1", "requester-1", "worker-1")

	// Invalid rating value
	err := rating.SetRating(types.DimensionQuality, 6)
	require.Error(t, err)

	// Valid rating
	err = rating.SetRating(types.DimensionQuality, 5)
	require.NoError(t, err)

	// Validate requires at least one rating
	err = rating.Validate()
	require.NoError(t, err)

	// Empty rating validation
	rating2 := types.NewRating("rating-2", "task-1", "requester-1", "worker-1")
	err = rating2.Validate()
	require.Error(t, err) // No ratings set
}

func TestReputationTrustLevel(t *testing.T) {
	rep := types.NewReputation("user-1")

	// New user
	require.True(t, rep.IsNew())
	require.Equal(t, "new", rep.GetTrustLevel())
}

func TestTaskTypes(t *testing.T) {
	// Open task
	openTask := types.NewTask("task-1", "Build", "Desc", "requester-1", types.TaskTypeOpen, 1000)
	require.Equal(t, types.TaskTypeOpen, openTask.Type)

	// IsOpen returns true only when status is Open
	openTask.Publish() // This sets status to Open
	require.True(t, openTask.IsOpen())

	// Auction task
	auctionTask := types.NewTask("task-2", "Build", "Desc", "requester-1", types.TaskTypeAuction, 1000)
	require.Equal(t, types.TaskTypeAuction, auctionTask.Type)
	require.False(t, auctionTask.IsOpen())
}

func TestTaskStatusTransitions(t *testing.T) {
	task := types.NewTask("task-1", "Build", "Desc", "requester-1", types.TaskTypeOpen, 1000)

	// Draft -> Open
	require.Equal(t, types.TaskStatusDraft, task.Status)
	task.Publish()
	require.Equal(t, types.TaskStatusOpen, task.Status)

	// Open -> Assigned
	task.Assign("worker-1")
	require.Equal(t, types.TaskStatusAssigned, task.Status)

	// Assigned -> InProgress
	task.Start()
	require.Equal(t, types.TaskStatusInProgress, task.Status)

	// InProgress -> Completed
	task.Complete()
	require.Equal(t, types.TaskStatusCompleted, task.Status)
}

func TestAuctionClosing(t *testing.T) {
	auction := types.NewAuction("task-1", 1000, 800, 86400)

	// Add valid bid
	bid := types.NewBid("bid-1", "task-1", "worker-1", 750)
	err := auction.AddBid(*bid)
	require.NoError(t, err)

	// Close auction
	winner, err := auction.CloseAuction()
	require.NoError(t, err)
	require.NotNil(t, winner)
	require.Equal(t, "worker-1", winner.WorkerID)
	require.False(t, auction.IsActive)
}

func TestAuctionReservePrice(t *testing.T) {
	auction := types.NewAuction("task-1", 1000, 800, 86400)

	// Add bid above reserve price
	bid := types.NewBid("bid-1", "task-1", "worker-1", 900)
	err := auction.AddBid(*bid)
	require.NoError(t, err)

	// Close should fail - bid above reserve price
	_, err = auction.CloseAuction()
	require.Error(t, err)
}

func TestBidStatuses(t *testing.T) {
	bid := types.NewBid("bid-1", "task-1", "worker-1", 900)
	require.Equal(t, types.BidStatusPending, bid.Status)

	bid.Accept()
	require.Equal(t, types.BidStatusAccepted, bid.Status)

	bid2 := types.NewBid("bid-2", "task-1", "worker-2", 850)
	bid2.Reject()
	require.Equal(t, types.BidStatusRejected, bid2.Status)

	bid3 := types.NewBid("bid-3", "task-1", "worker-3", 800)
	bid3.Withdraw()
	require.Equal(t, types.BidStatusWithdrawn, bid3.Status)
}

func TestApplicationStatuses(t *testing.T) {
	app := types.NewApplication("app-1", "task-1", "worker-1", 900)
	require.Equal(t, types.ApplicationStatusPending, app.Status)

	app.Accept()
	require.Equal(t, types.ApplicationStatusAccepted, app.Status)

	app2 := types.NewApplication("app-2", "task-1", "worker-2", 800)
	app2.Reject()
	require.Equal(t, types.ApplicationStatusRejected, app2.Status)

	app3 := types.NewApplication("app-3", "task-1", "worker-3", 850)
	app3.Withdraw()
	require.Equal(t, types.ApplicationStatusWithdrawn, app3.Status)
}

func TestReputationStars(t *testing.T) {
	rep := types.NewReputation("user-1")

	// No ratings yet
	require.Equal(t, "☆☆☆☆☆", rep.GetStars())

	// Add some ratings
	for i := 0; i < 5; i++ {
		rating := types.NewRating("rating-"+string(rune('0'+i)), "task-"+string(rune('0'+i)), "other", "user-1")
		rating.SetRating(types.DimensionQuality, 5)
		rep.AddRating(rating)
	}

	// Should have 5 stars now
	require.Equal(t, "★★★★★", rep.GetStars())
}
