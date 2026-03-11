package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// StoreKey is a type alias for byte slices
type StoreKey []byte

// Keeper manages task marketplace state using Cosmos SDK store
type Keeper struct {
	storeKey   StoreKey
	cdc        interface{} // Simplified - not using protobuf codec for now
	paramSpace types.ParamSubspace

	// Legacy in-memory keeper for backward compatibility
	legacyKeeper *LegacyKeeper
}

// LegacyKeeper provides backward compatibility with existing tests
type LegacyKeeper struct {
	tasks        map[string]*types.Task
	applications map[string]*types.Application
	auctions     map[string]*types.Auction
	ratings      map[string]*types.Rating
	reputations  map[string]*types.Reputation
}

// NewLegacyKeeper creates a new legacy keeper for in-memory operations
func NewLegacyKeeper() *LegacyKeeper {
	return &LegacyKeeper{
		tasks:        make(map[string]*types.Task),
		applications: make(map[string]*types.Application),
		auctions:     make(map[string]*types.Auction),
		ratings:      make(map[string]*types.Rating),
		reputations:  make(map[string]*types.Reputation),
	}
}

// NewKeeper creates a new task marketplace keeper
func NewKeeper() *Keeper {
	return &Keeper{
		legacyKeeper: NewLegacyKeeper(),
	}
}

// GetLegacyKeeper returns the legacy keeper for tests and migration
func (k Keeper) GetLegacyKeeper() *LegacyKeeper {
	return k.legacyKeeper
}

// Backward compatible methods (no context)

// CreateTask creates a new task
func (k Keeper) CreateTask(task *types.Task) error {
	return k.legacyKeeper.CreateTask(task)
}

// GetTask retrieves a task by ID
func (k Keeper) GetTask(id string) *types.Task {
	return k.legacyKeeper.GetTask(id)
}

// UpdateTask updates a task
func (k Keeper) UpdateTask(task *types.Task) error {
	return k.legacyKeeper.UpdateTask(task)
}

// GetAllTasks returns all tasks
func (k Keeper) GetAllTasks() []*types.Task {
	return k.legacyKeeper.GetAllTasks()
}

// GetTasksByRequester returns tasks by requester
func (k Keeper) GetTasksByRequester(requesterID string) []*types.Task {
	return k.legacyKeeper.GetTasksByRequester(requesterID)
}

// GetTasksByWorker returns tasks by worker
func (k Keeper) GetTasksByWorker(workerID string) []*types.Task {
	return k.legacyKeeper.GetTasksByWorker(workerID)
}

// GetOpenTasks returns open tasks
func (k Keeper) GetOpenTasks() []*types.Task {
	return k.legacyKeeper.GetOpenTasks()
}

// StartTask starts a task
func (k Keeper) StartTask(taskID string) error {
	return k.legacyKeeper.StartTask(taskID)
}

// SubmitMilestone submits a milestone
func (k Keeper) SubmitMilestone(taskID, milestoneID, deliverables string) error {
	return k.legacyKeeper.SubmitMilestone(taskID, milestoneID, deliverables)
}

// ApproveMilestone approves a milestone
func (k Keeper) ApproveMilestone(taskID, milestoneID string) error {
	return k.legacyKeeper.ApproveMilestone(taskID, milestoneID)
}

// SubmitApplication submits an application
func (k Keeper) SubmitApplication(app *types.Application) error {
	return k.legacyKeeper.SubmitApplication(app)
}

// GetApplication gets an application
func (k Keeper) GetApplication(id string) *types.Application {
	return k.legacyKeeper.GetApplication(id)
}

// GetApplicationsByTask gets applications for a task
func (k Keeper) GetApplicationsByTask(taskID string) []*types.Application {
	return k.legacyKeeper.GetApplicationsByTask(taskID)
}

// AcceptApplication accepts an application
func (k Keeper) AcceptApplication(appID string) error {
	return k.legacyKeeper.AcceptApplication(appID)
}

// RejectApplication rejects an application
func (k Keeper) RejectApplication(appID string) error {
	return k.legacyKeeper.RejectApplication(appID)
}

// CreateAuction creates an auction
func (k Keeper) CreateAuction(taskID string, startingPrice, reservePrice uint64, duration int64) (*types.Auction, error) {
	return k.legacyKeeper.CreateAuction(taskID, startingPrice, reservePrice, duration)
}

// GetAuction gets an auction
func (k Keeper) GetAuction(taskID string) *types.Auction {
	return k.legacyKeeper.GetAuction(taskID)
}

// SubmitBid submits a bid
func (k Keeper) SubmitBid(bid *types.Bid) error {
	return k.legacyKeeper.SubmitBid(bid)
}

// CloseAuction closes an auction
func (k Keeper) CloseAuction(taskID string) error {
	return k.legacyKeeper.CloseAuction(taskID)
}

// SubmitRating submits a rating
func (k Keeper) SubmitRating(rating *types.Rating) error {
	return k.legacyKeeper.SubmitRating(rating)
}

// GetRating gets a rating
func (k Keeper) GetRating(id string) *types.Rating {
	return k.legacyKeeper.GetRating(id)
}

// GetReputation gets reputation
func (k Keeper) GetReputation(userID string) *types.Reputation {
	return k.legacyKeeper.GetReputation(userID)
}

// GetTaskStatistics gets statistics
func (k Keeper) GetTaskStatistics() map[string]interface{} {
	return k.legacyKeeper.GetTaskStatistics()
}

// SDK-compatible methods (with context)

// SetTask sets a task in the store
func (k Keeper) SetTask(ctx sdk.Context, task types.Task) {
	k.legacyKeeper.tasks[task.ID] = &task
}

// SetApplication sets an application in the store
func (k Keeper) SetApplication(ctx sdk.Context, app types.Application) {
	k.legacyKeeper.applications[app.ID] = &app
}

// SetAuction sets an auction in the store
func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) {
	k.legacyKeeper.auctions[auction.TaskID] = &auction
}

// SetBid sets a bid in the store
func (k Keeper) SetBid(ctx sdk.Context, bid types.Bid) {
	// Bids are stored within auctions in the legacy keeper
}

// SetRating sets a rating in the store
func (k Keeper) SetRating(ctx sdk.Context, rating types.Rating) {
	k.legacyKeeper.ratings[rating.ID] = &rating
}

// SetReputation sets a reputation in the store
func (k Keeper) SetReputation(ctx sdk.Context, rep types.Reputation) {
	k.legacyKeeper.reputations[rep.UserID] = &rep
}

// LegacyKeeper methods

func (lk *LegacyKeeper) CreateTask(task *types.Task) error {
	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}
	if len(task.Milestones) > 0 {
		if err := task.ValidateMilestones(); err != nil {
			return err
		}
	}
	lk.tasks[task.ID] = task
	return nil
}

func (lk *LegacyKeeper) GetTask(id string) *types.Task {
	return lk.tasks[id]
}

func (lk *LegacyKeeper) UpdateTask(task *types.Task) error {
	lk.tasks[task.ID] = task
	return nil
}

func (lk *LegacyKeeper) GetAllTasks() []*types.Task {
	var result []*types.Task
	for _, task := range lk.tasks {
		result = append(result, task)
	}
	return result
}

func (lk *LegacyKeeper) GetTasksByRequester(requesterID string) []*types.Task {
	var tasks []*types.Task
	for _, task := range lk.tasks {
		if task.RequesterID == requesterID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (lk *LegacyKeeper) GetTasksByWorker(workerID string) []*types.Task {
	var tasks []*types.Task
	for _, task := range lk.tasks {
		if task.WorkerID == workerID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (lk *LegacyKeeper) GetOpenTasks() []*types.Task {
	var open []*types.Task
	for _, task := range lk.tasks {
		if task.IsOpen() {
			open = append(open, task)
		}
	}
	return open
}

func (lk *LegacyKeeper) SubmitApplication(app *types.Application) error {
	if err := app.Validate(); err != nil {
		return fmt.Errorf("invalid application: %w", err)
	}
	task := lk.GetTask(app.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", app.TaskID)
	}
	if task.Type != types.TaskTypeOpen {
		return fmt.Errorf("task is not open type")
	}
	if !task.IsOpen() {
		return fmt.Errorf("task is not accepting applications")
	}
	lk.applications[app.ID] = app
	task.ApplicationCount++
	return nil
}

func (lk *LegacyKeeper) GetApplication(id string) *types.Application {
	return lk.applications[id]
}

func (lk *LegacyKeeper) GetApplicationsByTask(taskID string) []*types.Application {
	var apps []*types.Application
	for _, app := range lk.applications {
		if app.TaskID == taskID {
			apps = append(apps, app)
		}
	}
	return apps
}

func (lk *LegacyKeeper) AcceptApplication(appID string) error {
	app := lk.GetApplication(appID)
	if app == nil {
		return fmt.Errorf("application not found: %s", appID)
	}
	if app.Status != types.ApplicationStatusPending {
		return fmt.Errorf("application is not pending")
	}
	task := lk.GetTask(app.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", app.TaskID)
	}
	app.Accept()
	task.Assign(app.WorkerID)
	return nil
}

func (lk *LegacyKeeper) RejectApplication(appID string) error {
	app := lk.GetApplication(appID)
	if app == nil {
		return fmt.Errorf("application not found: %s", appID)
	}
	app.Reject()
	return nil
}

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

func (lk *LegacyKeeper) GetAuction(taskID string) *types.Auction {
	return lk.auctions[taskID]
}

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

func (lk *LegacyKeeper) CloseAuction(taskID string) error {
	auction := lk.GetAuction(taskID)
	if auction == nil {
		return fmt.Errorf("auction not found: %s", taskID)
	}
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	winner, err := auction.CloseAuction()
	if err != nil {
		return err
	}
	auction.IsActive = false
	task.Assign(winner.WorkerID)
	return nil
}

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

func (lk *LegacyKeeper) GetRating(id string) *types.Rating {
	return lk.ratings[id]
}

func (lk *LegacyKeeper) GetReputation(userID string) *types.Reputation {
	rep, exists := lk.reputations[userID]
	if !exists {
		return types.NewReputation(userID)
	}
	return rep
}

func (lk *LegacyKeeper) StartTask(taskID string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	task.Start()
	if len(task.Milestones) > 0 {
		for i := range task.Milestones {
			if task.Milestones[i].Status == types.MilestoneStatusPending {
				task.Milestones[i].Status = types.MilestoneStatusActive
				break
			}
		}
	}
	return nil
}

func (lk *LegacyKeeper) SubmitMilestone(taskID, milestoneID, deliverables string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	return task.SubmitMilestone(milestoneID, deliverables)
}

func (lk *LegacyKeeper) ApproveMilestone(taskID, milestoneID string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	err := task.ApproveMilestone(milestoneID)
	if task.AllMilestonesCompleted() {
		task.Complete()
	}
	return err
}

func (lk *LegacyKeeper) RejectMilestone(taskID, milestoneID, reason string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	return task.RejectMilestone(milestoneID, reason)
}

func (lk *LegacyKeeper) GetTaskStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"total_tasks":        len(lk.tasks),
		"open_tasks":         0,
		"assigned_tasks":     0,
		"in_progress_tasks":  0,
		"completed_tasks":    0,
		"cancelled_tasks":    0,
		"total_applications": len(lk.applications),
		"total_bids":         0,
		"total_ratings":      len(lk.ratings),
	}

	for _, task := range lk.tasks {
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

	for _, auction := range lk.auctions {
		stats["total_bids"] = stats["total_bids"].(int) + len(auction.Bids)
	}

	return stats
}

// JSONCodec wraps encoding/json for SDK-like interface
type JSONCodec struct{}

func (c JSONCodec) Marshal(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func (c JSONCodec) Unmarshal(b []byte, o interface{}) error {
	return json.Unmarshal(b, o)
}
