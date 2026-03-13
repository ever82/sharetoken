package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sharetoken/x/taskmarket/types"
)

// Task Tests

func TestNewTask(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)

	require.NotNil(t, task)
	require.Equal(t, "task-1", task.ID)
	require.Equal(t, "Test Task", task.Title)
	require.Equal(t, "Description", task.Description)
	require.Equal(t, "requester-1", task.RequesterID)
	require.Equal(t, types.TaskTypeOpen, task.Type)
	require.Equal(t, types.TaskStatusDraft, task.Status)
	require.Equal(t, uint64(1000), task.Budget)
	require.Equal(t, "STT", task.Currency)
	require.NotNil(t, task.Skills)
	require.NotNil(t, task.Subtasks)
	require.NotNil(t, task.Milestones)
	require.NotZero(t, task.CreatedAt)
	require.NotZero(t, task.UpdatedAt)
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    types.Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: types.Task{
				ID:          "task-1",
				Title:       "Test Task",
				RequesterID: "requester-1",
				Budget:      1000,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty ID",
			task: types.Task{
				ID:          "",
				Title:       "Test Task",
				RequesterID: "requester-1",
				Budget:      1000,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty title",
			task: types.Task{
				ID:          "task-1",
				Title:       "",
				RequesterID: "requester-1",
				Budget:      1000,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty requester",
			task: types.Task{
				ID:          "task-1",
				Title:       "Test Task",
				RequesterID: "",
				Budget:      1000,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero budget",
			task: types.Task{
				ID:          "task-1",
				Title:       "Test Task",
				RequesterID: "requester-1",
				Budget:      0,
			},
			wantErr: true,
		},
		{
			name: "invalid - past deadline",
			task: types.Task{
				ID:          "task-1",
				Title:       "Test Task",
				RequesterID: "requester-1",
				Budget:      1000,
				Deadline:    time.Now().Unix() - 3600,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTask_Publish(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)
	require.Equal(t, types.TaskStatusDraft, task.Status)

	task.Publish()
	require.Equal(t, types.TaskStatusOpen, task.Status)
}

func TestTask_Assign(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task.Publish()

	task.Assign("worker-1")
	require.Equal(t, types.TaskStatusAssigned, task.Status)
	require.Equal(t, "worker-1", task.WorkerID)
	require.NotZero(t, task.AssignedAt)
}

func TestTask_Start(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task.Publish()
	task.Assign("worker-1")

	task.Start()
	require.Equal(t, types.TaskStatusInProgress, task.Status)
}

func TestTask_Complete(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task.Complete()
	require.Equal(t, types.TaskStatusCompleted, task.Status)
	require.NotZero(t, task.CompletedAt)
}

func TestTask_Cancel(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)
	task.Cancel()
	require.Equal(t, types.TaskStatusCancelled, task.Status)
}

func TestTask_IsOpen(t *testing.T) {
	tests := []struct {
		name     string
		status   types.TaskStatus
		expected bool
	}{
		{"open", types.TaskStatusOpen, true},
		{"draft", types.TaskStatusDraft, false},
		{"assigned", types.TaskStatusAssigned, false},
		{"completed", types.TaskStatusCompleted, false},
		{"cancelled", types.TaskStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := types.Task{Status: tt.status}
			require.Equal(t, tt.expected, task.IsOpen())
		})
	}
}

func TestTask_GetTotalMilestoneAmount(t *testing.T) {
	task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, 1000)

	// No milestones
	require.Equal(t, uint64(0), task.GetTotalMilestoneAmount())

	// Add milestones
	task.Milestones = []types.Milestone{
		{Amount: 300},
		{Amount: 400},
		{Amount: 300},
	}
	require.Equal(t, uint64(1000), task.GetTotalMilestoneAmount())
}

func TestTask_ValidateMilestones(t *testing.T) {
	tests := []struct {
		name    string
		budget  uint64
		amounts []uint64
		wantErr bool
	}{
		{
			name:    "valid - milestones match budget",
			budget:  1000,
			amounts: []uint64{400, 600},
			wantErr: false,
		},
		{
			name:    "invalid - milestones exceed budget",
			budget:  1000,
			amounts: []uint64{600, 600},
			wantErr: true,
		},
		{
			name:    "invalid - milestones less than budget",
			budget:  1000,
			amounts: []uint64{400, 400},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := types.NewTask("task-1", "Test Task", "Description", "requester-1", types.TaskTypeOpen, tt.budget)
			for i, amount := range tt.amounts {
				task.Milestones = append(task.Milestones, types.Milestone{
					ID:     string(rune('a' + i)),
					Amount: amount,
				})
			}
			err := task.ValidateMilestones()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTask_GetCompletionPercentage(t *testing.T) {
	tests := []struct {
		name     string
		status   types.TaskStatus
		milestoneStatuses []types.MilestoneStatus
		expected float64
	}{
		{
			name:     "no milestones - not completed",
			status:   types.TaskStatusInProgress,
			milestoneStatuses: []types.MilestoneStatus{},
			expected: 0.0,
		},
		{
			name:     "no milestones - completed",
			status:   types.TaskStatusCompleted,
			milestoneStatuses: []types.MilestoneStatus{},
			expected: 100.0,
		},
		{
			name:     "50% complete",
			status:   types.TaskStatusInProgress,
			milestoneStatuses: []types.MilestoneStatus{
				types.MilestoneStatusApproved,
				types.MilestoneStatusPending,
			},
			expected: 50.0,
		},
		{
			name:     "100% complete",
			status:   types.TaskStatusCompleted,
			milestoneStatuses: []types.MilestoneStatus{
				types.MilestoneStatusPaid,
				types.MilestoneStatusApproved,
			},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := types.Task{Status: tt.status}
			for _, status := range tt.milestoneStatuses {
				task.Milestones = append(task.Milestones, types.Milestone{Status: status})
			}
			require.InDelta(t, tt.expected, task.GetCompletionPercentage(), 0.01)
		})
	}
}

// Application Tests

func TestNewApplication(t *testing.T) {
	app := types.NewApplication("app-1", "task-1", "worker-1", 800)

	require.NotNil(t, app)
	require.Equal(t, "app-1", app.ID)
	require.Equal(t, "task-1", app.TaskID)
	require.Equal(t, "worker-1", app.WorkerID)
	require.Equal(t, types.ApplicationStatusPending, app.Status)
	require.Equal(t, uint64(800), app.ProposedPrice)
	require.NotZero(t, app.CreatedAt)
	require.NotZero(t, app.UpdatedAt)
}

func TestApplication_Validate(t *testing.T) {
	tests := []struct {
		name    string
		app     types.Application
		wantErr bool
	}{
		{
			name: "valid application",
			app: types.Application{
				TaskID:        "task-1",
				WorkerID:      "worker-1",
				ProposedPrice: 800,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty task ID",
			app: types.Application{
				TaskID:        "",
				WorkerID:      "worker-1",
				ProposedPrice: 800,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty worker ID",
			app: types.Application{
				TaskID:        "task-1",
				WorkerID:      "",
				ProposedPrice: 800,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero price",
			app: types.Application{
				TaskID:        "task-1",
				WorkerID:      "worker-1",
				ProposedPrice: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.app.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestApplication_Accept(t *testing.T) {
	app := types.NewApplication("app-1", "task-1", "worker-1", 800)
	require.Equal(t, types.ApplicationStatusPending, app.Status)

	app.Accept()
	require.Equal(t, types.ApplicationStatusAccepted, app.Status)
}

func TestApplication_Reject(t *testing.T) {
	app := types.NewApplication("app-1", "task-1", "worker-1", 800)
	app.Reject()
	require.Equal(t, types.ApplicationStatusRejected, app.Status)
}

func TestApplication_Withdraw(t *testing.T) {
	app := types.NewApplication("app-1", "task-1", "worker-1", 800)
	app.Withdraw()
	require.Equal(t, types.ApplicationStatusWithdrawn, app.Status)
}

// Bid Tests

func TestNewBid(t *testing.T) {
	bid := types.NewBid("bid-1", "task-1", "worker-1", 500)

	require.NotNil(t, bid)
	require.Equal(t, "bid-1", bid.ID)
	require.Equal(t, "task-1", bid.TaskID)
	require.Equal(t, "worker-1", bid.WorkerID)
	require.Equal(t, types.BidStatusPending, bid.Status)
	require.Equal(t, uint64(500), bid.Amount)
	require.NotZero(t, bid.CreatedAt)
	require.NotZero(t, bid.UpdatedAt)
}

func TestBid_Validate(t *testing.T) {
	tests := []struct {
		name    string
		bid     types.Bid
		wantErr bool
	}{
		{
			name: "valid bid",
			bid: types.Bid{
				TaskID:   "task-1",
				WorkerID: "worker-1",
				Amount:   500,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty task ID",
			bid: types.Bid{
				TaskID:   "",
				WorkerID: "worker-1",
				Amount:   500,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty worker ID",
			bid: types.Bid{
				TaskID:   "task-1",
				WorkerID: "",
				Amount:   500,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero amount",
			bid: types.Bid{
				TaskID:   "task-1",
				WorkerID: "worker-1",
				Amount:   0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bid.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBid_IsLowerThan(t *testing.T) {
	bid1 := &types.Bid{Amount: 400}
	bid2 := &types.Bid{Amount: 500}

	require.True(t, bid1.IsLowerThan(bid2))
	require.False(t, bid2.IsLowerThan(bid1))
}

// Auction Tests

func TestNewAuction(t *testing.T) {
	auction := types.NewAuction("task-1", 1000, 800, 86400)

	require.NotNil(t, auction)
	require.Equal(t, "task-1", auction.TaskID)
	require.Equal(t, uint64(1000), auction.StartingPrice)
	require.Equal(t, uint64(800), auction.ReservePrice)
	require.True(t, auction.IsActive)
	require.NotNil(t, auction.Bids)
	require.Empty(t, auction.Bids)
	require.NotZero(t, auction.EndTime)
}

func TestAuction_AddBid(t *testing.T) {
	now := time.Now().Unix()
	auction := types.NewAuction("task-1", 1000, 800, 3600)

	tests := []struct {
		name    string
		bid     types.Bid
		wantErr bool
	}{
		{
			name:    "valid bid",
			bid:     types.Bid{ID: "bid-1", Amount: 900},
			wantErr: false,
		},
		{
			name:    "invalid - bid exceeds starting price",
			bid:     types.Bid{ID: "bid-2", Amount: 1100},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset auction for each test
			auction = types.NewAuction("task-1", 1000, 800, 3600)
			err := auction.AddBid(tt.bid)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, 1, len(auction.Bids))
			}
		})
	}

	// Test expired auction
	expiredAuction := types.Auction{
		TaskID:        "task-1",
		StartingPrice: 1000,
		EndTime:       now - 3600,
		IsActive:      true,
	}
	err := expiredAuction.AddBid(types.Bid{ID: "bid-1", Amount: 900})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ended")
}

func TestAuction_GetWinningBid(t *testing.T) {
	auction := types.NewAuction("task-1", 1000, 800, 3600)

	// No bids yet
	require.Nil(t, auction.GetWinningBid())

	// Add bids
	_ = auction.AddBid(types.Bid{ID: "bid-1", Amount: 900, Status: types.BidStatusPending})
	_ = auction.AddBid(types.Bid{ID: "bid-2", Amount: 800, Status: types.BidStatusPending})

	winner := auction.GetWinningBid()
	require.NotNil(t, winner)
	require.Equal(t, "bid-2", winner.ID) // Lower amount wins
	require.Equal(t, uint64(800), winner.Amount)
}

func TestAuction_CloseAuction(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name        string
		auction     types.Auction
		expectError bool
	}{
		{
			name: "successful close",
			auction: types.Auction{
				TaskID:        "task-1",
				ReservePrice:  800,
				Bids:          []types.Bid{{ID: "bid-1", Amount: 700, Status: types.BidStatusPending}},
				WinningBidID:  "bid-1",
				EndTime:       now + 3600,
				IsActive:      true,
			},
			expectError: false,
		},
		{
			name: "no valid bids",
			auction: types.Auction{
				TaskID:        "task-1",
				ReservePrice:  800,
				Bids:          []types.Bid{},
				IsActive:      true,
			},
			expectError: true,
		},
		{
			name: "bid below reserve",
			auction: types.Auction{
				TaskID:        "task-1",
				ReservePrice:  800,
				Bids:          []types.Bid{{ID: "bid-1", Amount: 900, Status: types.BidStatusPending}},
				WinningBidID:  "bid-1",
				IsActive:      true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winner, err := tt.auction.CloseAuction()
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, winner)
			} else {
				require.NoError(t, err)
				require.NotNil(t, winner)
				require.Equal(t, types.BidStatusAccepted, winner.Status)
			}
		})
	}
}

func TestAuction_IsEnded(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name     string
		auction  types.Auction
		expected bool
	}{
		{
			name:     "not ended - active and time remaining",
			auction:  types.Auction{IsActive: true, EndTime: now + 3600},
			expected: false,
		},
		{
			name:     "ended - expired",
			auction:  types.Auction{IsActive: true, EndTime: now - 3600},
			expected: true,
		},
		{
			name:     "ended - not active",
			auction:  types.Auction{IsActive: false, EndTime: now + 3600},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.auction.IsEnded())
		})
	}
}

// Genesis Tests

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Tasks)
	require.NotNil(t, genesis.Applications)
	require.NotNil(t, genesis.Auctions)
	require.NotNil(t, genesis.Bids)
	require.NotNil(t, genesis.Ratings)
	require.NotNil(t, genesis.Reputations)
	require.Empty(t, genesis.Tasks)
	require.Empty(t, genesis.Applications)
}

func TestValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		data    types.GenesisState
		wantErr bool
	}{
		{
			name:    "valid genesis with default",
			data:    *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "valid genesis with data",
			data: types.GenesisState{
				Tasks: []types.Task{
					{ID: "task-1", Title: "Task 1", RequesterID: "req-1", Budget: 1000},
				},
				Applications: []types.Application{
					{ID: "app-1", TaskID: "task-1", WorkerID: "worker-1", ProposedPrice: 800},
				},
				Bids: []types.Bid{
					{ID: "bid-1", TaskID: "task-1", WorkerID: "worker-1", Amount: 500},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - invalid task",
			data: types.GenesisState{
				Tasks: []types.Task{
					{ID: "", Title: "Task 1", RequesterID: "req-1", Budget: 1000},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid application",
			data: types.GenesisState{
				Applications: []types.Application{
					{ID: "app-1", TaskID: "", WorkerID: "worker-1", ProposedPrice: 800},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid bid",
			data: types.GenesisState{
				Bids: []types.Bid{
					{ID: "bid-1", TaskID: "task-1", WorkerID: "", Amount: 500},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateGenesis(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
