package keeper

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"sharetoken/x/taskmarket/types"
)

// Keeper manages task marketplace state
type Keeper struct {
	tasks        map[string]*types.Task
	applications map[string]*types.Application
	auctions     map[string]*types.Auction
	ratings      map[string]*types.Rating
	reputations  map[string]*types.Reputation
	mutex        sync.RWMutex
}

// NewKeeper creates a new task marketplace keeper
func NewKeeper() *Keeper {
	return &Keeper{
		tasks:        make(map[string]*types.Task),
		applications: make(map[string]*types.Application),
		auctions:     make(map[string]*types.Auction),
		ratings:      make(map[string]*types.Rating),
		reputations:  make(map[string]*types.Reputation),
	}
}

// CreateTask creates a new task
func (k *Keeper) CreateTask(task *types.Task) error {
	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}

	if len(task.Milestones) > 0 {
		if err := task.ValidateMilestones(); err != nil {
			return err
		}
	}

	k.mutex.Lock()
	k.tasks[task.ID] = task
	k.mutex.Unlock()

	return nil
}

// GetTask retrieves a task
func (k *Keeper) GetTask(id string) *types.Task {
	k.mutex.RLock()
	defer k.mutex.RUnlock()
	return k.tasks[id]
}

// UpdateTask updates a task
func (k *Keeper) UpdateTask(task *types.Task) error {
	k.mutex.Lock()
	k.tasks[task.ID] = task
	k.mutex.Unlock()
	return nil
}

// GetOpenTasks returns all open tasks
func (k *Keeper) GetOpenTasks() []*types.Task {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	var open []*types.Task
	for _, task := range k.tasks {
		if task.IsOpen() {
			open = append(open, task)
		}
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i].CreatedAt > open[j].CreatedAt
	})

	return open
}

// GetTasksByRequester returns tasks by requester
func (k *Keeper) GetTasksByRequester(requesterID string) []*types.Task {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	var tasks []*types.Task
	for _, task := range k.tasks {
		if task.RequesterID == requesterID {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// GetTasksByWorker returns tasks by worker
func (k *Keeper) GetTasksByWorker(workerID string) []*types.Task {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	var tasks []*types.Task
	for _, task := range k.tasks {
		if task.WorkerID == workerID {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// SubmitApplication submits an application
func (k *Keeper) SubmitApplication(app *types.Application) error {
	if err := app.Validate(); err != nil {
		return fmt.Errorf("invalid application: %w", err)
	}

	task := k.GetTask(app.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", app.TaskID)
	}

	if task.Type != types.TaskTypeOpen {
		return fmt.Errorf("task is not open type")
	}

	if !task.IsOpen() {
		return fmt.Errorf("task is not accepting applications")
	}

	k.mutex.Lock()
	k.applications[app.ID] = app
	task.ApplicationCount++
	k.mutex.Unlock()

	return nil
}

// GetApplication retrieves an application
func (k *Keeper) GetApplication(id string) *types.Application {
	k.mutex.RLock()
	defer k.mutex.RUnlock()
	return k.applications[id]
}

// GetApplicationsByTask returns applications for a task
func (k *Keeper) GetApplicationsByTask(taskID string) []*types.Application {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	var apps []*types.Application
	for _, app := range k.applications {
		if app.TaskID == taskID {
			apps = append(apps, app)
		}
	}

	return apps
}

// AcceptApplication accepts an application
func (k *Keeper) AcceptApplication(appID string) error {
	app := k.GetApplication(appID)
	if app == nil {
		return fmt.Errorf("application not found: %s", appID)
	}

	if app.Status != types.ApplicationStatusPending {
		return fmt.Errorf("application is not pending")
	}

	task := k.GetTask(app.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", app.TaskID)
	}

	k.mutex.Lock()
	app.Accept()
	task.Assign(app.WorkerID)
	k.mutex.Unlock()

	return nil
}

// RejectApplication rejects an application
func (k *Keeper) RejectApplication(appID string) error {
	app := k.GetApplication(appID)
	if app == nil {
		return fmt.Errorf("application not found: %s", appID)
	}

	k.mutex.Lock()
	app.Reject()
	k.mutex.Unlock()

	return nil
}

// CreateAuction creates an auction
func (k *Keeper) CreateAuction(taskID string, startingPrice, reservePrice uint64, duration int64) (*types.Auction, error) {
	task := k.GetTask(taskID)
	if task == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	if task.Type != types.TaskTypeAuction {
		return nil, fmt.Errorf("task is not auction type")
	}

	auction := types.NewAuction(taskID, startingPrice, reservePrice, duration)

	k.mutex.Lock()
	k.auctions[taskID] = auction
	k.mutex.Unlock()

	return auction, nil
}

// GetAuction retrieves an auction
func (k *Keeper) GetAuction(taskID string) *types.Auction {
	k.mutex.RLock()
	defer k.mutex.RUnlock()
	return k.auctions[taskID]
}

// SubmitBid submits a bid
func (k *Keeper) SubmitBid(bid *types.Bid) error {
	if err := bid.Validate(); err != nil {
		return fmt.Errorf("invalid bid: %w", err)
	}

	auction := k.GetAuction(bid.TaskID)
	if auction == nil {
		return fmt.Errorf("auction not found: %s", bid.TaskID)
	}

	if err := auction.AddBid(*bid); err != nil {
		return err
	}

	task := k.GetTask(bid.TaskID)

	k.mutex.Lock()
	k.auctions[bid.TaskID] = auction
	if task != nil {
		task.BidCount = auction.GetBidCount()
	}
	k.mutex.Unlock()

	return nil
}

// CloseAuction closes an auction
func (k *Keeper) CloseAuction(taskID string) error {
	auction := k.GetAuction(taskID)
	if auction == nil {
		return fmt.Errorf("auction not found: %s", taskID)
	}

	task := k.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	winner, err := auction.CloseAuction()
	if err != nil {
		return err
	}

	k.mutex.Lock()
	auction.IsActive = false
	task.Assign(winner.WorkerID)
	k.mutex.Unlock()

	return nil
}

// SubmitRating submits a rating
func (k *Keeper) SubmitRating(rating *types.Rating) error {
	if err := rating.Validate(); err != nil {
		return fmt.Errorf("invalid rating: %w", err)
	}

	task := k.GetTask(rating.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", rating.TaskID)
	}

	if task.Status != types.TaskStatusCompleted {
		return fmt.Errorf("task is not completed")
	}

	k.mutex.Lock()
	k.ratings[rating.ID] = rating

	rep, exists := k.reputations[rating.RatedID]
	if !exists {
		rep = types.NewReputation(rating.RatedID)
	}
	rep.AddRating(rating)
	k.reputations[rating.RatedID] = rep

	k.mutex.Unlock()

	return nil
}

// GetReputation retrieves reputation
func (k *Keeper) GetReputation(userID string) *types.Reputation {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	rep, exists := k.reputations[userID]
	if !exists {
		return types.NewReputation(userID)
	}
	return rep
}

// StartTask starts a task
func (k *Keeper) StartTask(taskID string) error {
	task := k.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	k.mutex.Lock()
	task.Start()

	if len(task.Milestones) > 0 {
		for i := range task.Milestones {
			if task.Milestones[i].Status == types.MilestoneStatusPending {
				task.Milestones[i].Status = types.MilestoneStatusActive
				break
			}
		}
	}

	k.mutex.Unlock()

	return nil
}

// SubmitMilestone submits a milestone
func (k *Keeper) SubmitMilestone(taskID, milestoneID, deliverables string) error {
	task := k.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	k.mutex.Lock()
	err := task.SubmitMilestone(milestoneID, deliverables)
	k.mutex.Unlock()

	return err
}

// ApproveMilestone approves a milestone
func (k *Keeper) ApproveMilestone(taskID, milestoneID string) error {
	task := k.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	k.mutex.Lock()
	err := task.ApproveMilestone(milestoneID)

	if task.AllMilestonesCompleted() {
		task.Complete()
	}

	k.mutex.Unlock()

	return err
}

// RejectMilestone rejects a milestone
func (k *Keeper) RejectMilestone(taskID, milestoneID string, reason string) error {
	task := k.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	k.mutex.Lock()
	err := task.RejectMilestone(milestoneID, reason)
	k.mutex.Unlock()

	return err
}

// GetTaskStatistics returns statistics
func (k *Keeper) GetTaskStatistics() map[string]interface{} {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_tasks":        len(k.tasks),
		"open_tasks":         0,
		"assigned_tasks":     0,
		"in_progress_tasks":  0,
		"completed_tasks":    0,
		"cancelled_tasks":    0,
		"total_applications": len(k.applications),
		"total_bids":         0,
		"total_ratings":      len(k.ratings),
	}

	for _, task := range k.tasks {
		switch task.Status {
		case types.TaskStatusOpen:
			stats["open_tasks"] = stats["open_tasks"].(int) + 1
		case types.TaskStatusAssigned:
			stats["assigned_tasks"] = stats["assigned_tasks"].(int) + 1
		case types.TaskStatusInProgress:
			stats["in_progress_tasks"] = stats["in_progress_tasks"].(int) + 1
		case types.TaskStatusCompleted:
			stats["completed_tasks"] = stats["completed_tasks"].(int) + 1
		case types.TaskStatusCancelled:
			stats["cancelled_tasks"] = stats["cancelled_tasks"].(int) + 1
		}
	}

	for _, auction := range k.auctions {
		stats["total_bids"] = stats["total_bids"].(int) + len(auction.Bids)
	}

	return stats
}

// AutoCloseAuctions closes expired auctions
func (k *Keeper) AutoCloseAuctions(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			k.mutex.Lock()
			for taskID, auction := range k.auctions {
				if auction.IsActive && auction.IsEnded() {
					if winner, err := auction.CloseAuction(); err == nil && winner != nil {
						if task := k.tasks[taskID]; task != nil {
							task.Assign(winner.WorkerID)
						}
					}
				}
			}
			k.mutex.Unlock()
		}
	}
}
