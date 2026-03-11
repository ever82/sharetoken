package keeper

import (
	"encoding/json"

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

// LegacyKeeper provides backward compatibility with existing tests.
// Optimized with composite indexes for O(1) lookups.
type LegacyKeeper struct {
	tasks        map[string]*types.Task
	applications map[string]*types.Application
	auctions     map[string]*types.Auction
	ratings      map[string]*types.Rating
	reputations  map[string]*types.Reputation

	// Composite indexes for O(1) lookup instead of O(n) scan
	tasksByRequester map[string][]string // requesterID -> taskIDs
	tasksByWorker    map[string][]string // workerID -> taskIDs
	tasksByStatus    map[string][]string // status -> taskIDs
}

// NewLegacyKeeper creates a new legacy keeper for in-memory operations
func NewLegacyKeeper() *LegacyKeeper {
	return &LegacyKeeper{
		tasks:            make(map[string]*types.Task),
		applications:     make(map[string]*types.Application),
		auctions:         make(map[string]*types.Auction),
		ratings:          make(map[string]*types.Rating),
		reputations:      make(map[string]*types.Reputation),
		tasksByRequester: make(map[string][]string),
		tasksByWorker:    make(map[string][]string),
		tasksByStatus:    make(map[string][]string),
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

// JSONCodec wraps encoding/json for SDK-like interface
type JSONCodec struct{}

func (c JSONCodec) Marshal(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func (c JSONCodec) Unmarshal(b []byte, o interface{}) error {
	return json.Unmarshal(b, o)
}
