package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// StoreKey is a type alias for byte slices
type StoreKey []byte

// codec is a simple interface for marshaling/unmarshaling
type codec interface {
	Marshal(o interface{}) ([]byte, error)
	Unmarshal(b []byte, o interface{}) error
}

// JSONCodec wraps encoding/json for SDK-like interface
type JSONCodec struct{}

func (c JSONCodec) Marshal(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func (c JSONCodec) Unmarshal(b []byte, o interface{}) error {
	return json.Unmarshal(b, o)
}

// Keeper manages task marketplace state using Cosmos SDK KVStore
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec
	paramSpace types.ParamSubspace
}

// NewKeeper creates a new task marketplace keeper
func NewKeeper(storeKey storetypes.StoreKey, cdc codec, paramSpace types.ParamSubspace) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramSpace: paramSpace,
	}
}

// NewKeeperWithCodec creates a keeper with the provided codec
func NewKeeperWithCodec(storeKey storetypes.StoreKey, cdc codec) *Keeper {
	return &Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// NewKeeperSimple creates a new task marketplace keeper without dependencies (for tests)
func NewKeeperSimple() *Keeper {
	return &Keeper{
		storeKey: storetypes.NewKVStoreKey(types.StoreKey),
		cdc:      JSONCodec{},
	}
}

// NewKeeperWithDefaultCodec creates a keeper with default JSON codec (backward compatible)
func NewKeeperWithDefaultCodec(storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{
		storeKey: storeKey,
		cdc:      JSONCodec{},
	}
}

// GetStoreKey returns the store key
func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

// getStore returns the KVStore for this module
func (k Keeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

// Task CRUD Operations

// SetTask stores a task in the KVStore
func (k Keeper) SetTask(ctx sdk.Context, task types.Task) {
	store := k.getStore(ctx)
	key := types.GetTaskKey(task.Id)
	value, err := k.cdc.Marshal(task)
	if err != nil {
		panic(fmt.Errorf("failed to marshal task: %w", err))
	}
	store.Set(key, value)

	// Update indexes
	k.setTaskByRequester(ctx, task)
	k.setTaskByWorker(ctx, task)
	k.setTaskByStatus(ctx, task)
}

// GetTask retrieves a task by ID from KVStore
func (k Keeper) GetTask(ctx sdk.Context, id string) (types.Task, bool) {
	store := k.getStore(ctx)
	key := types.GetTaskKey(id)
	value := store.Get(key)
	if value == nil {
		return types.Task{}, false
	}

	var task types.Task
	if err := k.cdc.Unmarshal(value, &task); err != nil {
		panic(fmt.Errorf("failed to unmarshal task: %w", err))
	}
	return task, true
}

// DeleteTask removes a task from the KVStore
func (k Keeper) DeleteTask(ctx sdk.Context, id string) {
	store := k.getStore(ctx)

	// Get task first to clean up indexes
	task, found := k.GetTask(ctx, id)
	if found {
		k.deleteTaskIndexes(ctx, task)
	}

	key := types.GetTaskKey(id)
	store.Delete(key)
}

// HasTask checks if a task exists
func (k Keeper) HasTask(ctx sdk.Context, id string) bool {
	store := k.getStore(ctx)
	key := types.GetTaskKey(id)
	return store.Has(key)
}

// GetAllTasks returns all tasks from KVStore
func (k Keeper) GetAllTasks(ctx sdk.Context) []types.Task {
	var tasks []types.Task
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var task types.Task
		if err := k.cdc.Unmarshal(iterator.Value(), &task); err != nil {
			panic(fmt.Errorf("failed to unmarshal task: %w", err))
		}
		tasks = append(tasks, task)
	}
	return tasks
}

// CreateTask creates a new task (SDK version)
func (k Keeper) CreateTask(ctx sdk.Context, task types.Task) error {
	if k.HasTask(ctx, task.Id) {
		return fmt.Errorf("task already exists: %s", task.Id)
	}

	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}

	if len(task.Milestones) > 0 {
		if err := task.ValidateMilestones(); err != nil {
			return err
		}
	}

	k.SetTask(ctx, task)
	return nil
}

// UpdateTask updates a task in the store
func (k Keeper) UpdateTask(ctx sdk.Context, task types.Task) error {
	if !k.HasTask(ctx, task.Id) {
		return fmt.Errorf("task not found: %s", task.Id)
	}

	// Delete old indexes
	oldTask, _ := k.GetTask(ctx, task.Id)
	k.deleteTaskIndexes(ctx, oldTask)

	// Set new task with updated indexes
	k.SetTask(ctx, task)
	return nil
}

// Index management functions

func (k Keeper) setTaskByRequester(ctx sdk.Context, task types.Task) {
	if task.RequesterId == "" {
		return
	}
	store := k.getStore(ctx)
	key := types.GetTaskByRequesterKey(task.RequesterId, task.Id)
	store.Set(key, []byte(task.Id))
}

func (k Keeper) setTaskByWorker(ctx sdk.Context, task types.Task) {
	if task.WorkerId == "" {
		return
	}
	store := k.getStore(ctx)
	key := types.GetTaskByWorkerKey(task.WorkerId, task.Id)
	store.Set(key, []byte(task.Id))
}

func (k Keeper) setTaskByStatus(ctx sdk.Context, task types.Task) {
	store := k.getStore(ctx)
	key := types.GetTaskByStatusKey(string(task.Status), task.Id)
	store.Set(key, []byte(task.Id))
}

func (k Keeper) deleteTaskIndexes(ctx sdk.Context, task types.Task) {
	store := k.getStore(ctx)

	// Delete requester index
	if task.RequesterId != "" {
		key := types.GetTaskByRequesterKey(task.RequesterId, task.Id)
		store.Delete(key)
	}

	// Delete worker index
	if task.WorkerId != "" {
		key := types.GetTaskByWorkerKey(task.WorkerId, task.Id)
		store.Delete(key)
	}

	// Delete status index
	key := types.GetTaskByStatusKey(string(task.Status), task.Id)
	store.Delete(key)
}

// GetTasksByRequester returns tasks by requester using composite index
func (k Keeper) GetTasksByRequester(ctx sdk.Context, requesterID string) []types.Task {
	var tasks []types.Task
	store := k.getStore(ctx)
	prefix := types.GetTaskByRequesterPrefix(requesterID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		taskID := string(iterator.Value())
		task, found := k.GetTask(ctx, taskID)
		if found {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// GetTasksByWorker returns tasks by worker using composite index
func (k Keeper) GetTasksByWorker(ctx sdk.Context, workerID string) []types.Task {
	var tasks []types.Task
	store := k.getStore(ctx)
	prefix := types.GetTaskByWorkerPrefix(workerID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		taskID := string(iterator.Value())
		task, found := k.GetTask(ctx, taskID)
		if found {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// GetTasksByStatus returns tasks by status
func (k Keeper) GetTasksByStatus(ctx sdk.Context, status string) []types.Task {
	var tasks []types.Task
	store := k.getStore(ctx)
	prefix := types.GetTaskByStatusPrefix(status)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		taskID := string(iterator.Value())
		task, found := k.GetTask(ctx, taskID)
		if found {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// GetOpenTasks returns all open tasks
func (k Keeper) GetOpenTasks(ctx sdk.Context) []types.Task {
	return k.GetTasksByStatus(ctx, string(types.TaskStatusOpen))
}

// Task Lifecycle Operations

// StartTask starts a task
func (k Keeper) StartTask(ctx sdk.Context, taskID string) error {
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Delete old status index
	k.deleteTaskIndexes(ctx, task)

	task.Start()

	if len(task.Milestones) > 0 {
		for i := range task.Milestones {
			if task.Milestones[i].Status == types.MilestoneStatusPending {
				task.Milestones[i].Status = types.MilestoneStatusActive
				break
			}
		}
	}

	k.SetTask(ctx, task)
	return nil
}

// SubmitMilestone submits a milestone
func (k Keeper) SubmitMilestone(ctx sdk.Context, taskID, milestoneID, deliverables string) error {
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}
	if err := task.SubmitMilestone(milestoneID, deliverables); err != nil {
		return err
	}
	return k.UpdateTask(ctx, task)
}

// ApproveMilestone approves a milestone
func (k Keeper) ApproveMilestone(ctx sdk.Context, taskID, milestoneID string) error {
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Delete old indexes before status change
	k.deleteTaskIndexes(ctx, task)

	err := task.ApproveMilestone(milestoneID)
	if task.AllMilestonesCompleted() {
		task.Complete()
	}

	k.SetTask(ctx, task)
	return err
}

// RejectMilestone rejects a milestone
func (k Keeper) RejectMilestone(ctx sdk.Context, taskID, milestoneID, reason string) error {
	task, found := k.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}
	if err := task.RejectMilestone(milestoneID, reason); err != nil {
		return err
	}
	return k.UpdateTask(ctx, task)
}

// Application CRUD Operations

// SetApplication stores an application in the KVStore
func (k Keeper) SetApplication(ctx sdk.Context, app types.Application) {
	store := k.getStore(ctx)
	key := types.GetApplicationKey(app.Id)
	value, err := k.cdc.Marshal(app)
	if err != nil {
		panic(fmt.Errorf("failed to marshal application: %w", err))
	}
	store.Set(key, value)

	// Update index
	k.setApplicationByTask(ctx, app)
}

// GetApplication retrieves an application by ID from KVStore
func (k Keeper) GetApplication(ctx sdk.Context, id string) (types.Application, bool) {
	store := k.getStore(ctx)
	key := types.GetApplicationKey(id)
	value := store.Get(key)
	if value == nil {
		return types.Application{}, false
	}

	var app types.Application
	if err := k.cdc.Unmarshal(value, &app); err != nil {
		panic(fmt.Errorf("failed to unmarshal application: %w", err))
	}
	return app, true
}

// GetAllApplications returns all applications from KVStore
func (k Keeper) GetAllApplications(ctx sdk.Context) []types.Application {
	var apps []types.Application
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.ApplicationKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var app types.Application
		if err := k.cdc.Unmarshal(iterator.Value(), &app); err != nil {
			panic(fmt.Errorf("failed to unmarshal application: %w", err))
		}
		apps = append(apps, app)
	}
	return apps
}

func (k Keeper) setApplicationByTask(ctx sdk.Context, app types.Application) {
	store := k.getStore(ctx)
	key := types.GetApplicationByTaskKey(app.TaskId, app.Id)
	store.Set(key, []byte(app.Id))
}

// GetApplicationsByTask returns applications by task using composite index
func (k Keeper) GetApplicationsByTask(ctx sdk.Context, taskID string) []types.Application {
	var apps []types.Application
	store := k.getStore(ctx)
	prefix := types.GetApplicationByTaskPrefix(taskID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		appID := string(iterator.Value())
		app, found := k.GetApplication(ctx, appID)
		if found {
			apps = append(apps, app)
		}
	}
	return apps
}

// Auction CRUD Operations

// SetAuction stores an auction in the KVStore
func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) {
	store := k.getStore(ctx)
	key := types.GetAuctionKey(auction.TaskId)
	value, err := k.cdc.Marshal(auction)
	if err != nil {
		panic(fmt.Errorf("failed to marshal auction: %w", err))
	}
	store.Set(key, value)
}

// GetAuction retrieves an auction by task ID from KVStore
func (k Keeper) GetAuction(ctx sdk.Context, taskID string) (types.Auction, bool) {
	store := k.getStore(ctx)
	key := types.GetAuctionKey(taskID)
	value := store.Get(key)
	if value == nil {
		return types.Auction{}, false
	}

	var auction types.Auction
	if err := k.cdc.Unmarshal(value, &auction); err != nil {
		panic(fmt.Errorf("failed to unmarshal auction: %w", err))
	}
	return auction, true
}

// GetAllAuctions returns all auctions from KVStore
func (k Keeper) GetAllAuctions(ctx sdk.Context) []types.Auction {
	var auctions []types.Auction
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.AuctionKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var auction types.Auction
		if err := k.cdc.Unmarshal(iterator.Value(), &auction); err != nil {
			panic(fmt.Errorf("failed to unmarshal auction: %w", err))
		}
		auctions = append(auctions, auction)
	}
	return auctions
}

// Bid CRUD Operations

// SetBid stores a bid in the KVStore
func (k Keeper) SetBid(ctx sdk.Context, bid types.Bid) {
	store := k.getStore(ctx)
	key := types.GetBidKey(bid.Id)
	value, err := k.cdc.Marshal(bid)
	if err != nil {
		panic(fmt.Errorf("failed to marshal bid: %w", err))
	}
	store.Set(key, value)

	// Update index
	k.setBidByTask(ctx, bid)
}

// GetBid retrieves a bid by ID from KVStore
func (k Keeper) GetBid(ctx sdk.Context, id string) (types.Bid, bool) {
	store := k.getStore(ctx)
	key := types.GetBidKey(id)
	value := store.Get(key)
	if value == nil {
		return types.Bid{}, false
	}

	var bid types.Bid
	if err := k.cdc.Unmarshal(value, &bid); err != nil {
		panic(fmt.Errorf("failed to unmarshal bid: %w", err))
	}
	return bid, true
}

// GetAllBids returns all bids from KVStore
func (k Keeper) GetAllBids(ctx sdk.Context) []types.Bid {
	var bids []types.Bid
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.BidKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var bid types.Bid
		if err := k.cdc.Unmarshal(iterator.Value(), &bid); err != nil {
			panic(fmt.Errorf("failed to unmarshal bid: %w", err))
		}
		bids = append(bids, bid)
	}
	return bids
}

func (k Keeper) setBidByTask(ctx sdk.Context, bid types.Bid) {
	store := k.getStore(ctx)
	key := types.GetBidByTaskKey(bid.TaskId, bid.Id)
	store.Set(key, []byte(bid.Id))
}

// GetBidsByTask returns bids by task using composite index
func (k Keeper) GetBidsByTask(ctx sdk.Context, taskID string) []types.Bid {
	var bids []types.Bid
	store := k.getStore(ctx)
	prefix := types.GetBidByTaskPrefix(taskID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		bidID := string(iterator.Value())
		bid, found := k.GetBid(ctx, bidID)
		if found {
			bids = append(bids, bid)
		}
	}
	return bids
}

// Rating CRUD Operations

// SetRating stores a rating in the KVStore
func (k Keeper) SetRating(ctx sdk.Context, rating types.Rating) {
	store := k.getStore(ctx)
	key := types.GetRatingKey(rating.Id)
	value, err := k.cdc.Marshal(rating)
	if err != nil {
		panic(fmt.Errorf("failed to marshal rating: %w", err))
	}
	store.Set(key, value)

	// Update indexes
	k.setRatingByTask(ctx, rating)
	k.setRatingByRatedUser(ctx, rating)
}

// GetRating retrieves a rating by ID from KVStore
func (k Keeper) GetRating(ctx sdk.Context, id string) (types.Rating, bool) {
	store := k.getStore(ctx)
	key := types.GetRatingKey(id)
	value := store.Get(key)
	if value == nil {
		return types.Rating{}, false
	}

	var rating types.Rating
	if err := k.cdc.Unmarshal(value, &rating); err != nil {
		panic(fmt.Errorf("failed to unmarshal rating: %w", err))
	}
	return rating, true
}

// GetAllRatings returns all ratings from KVStore
func (k Keeper) GetAllRatings(ctx sdk.Context) []types.Rating {
	var ratings []types.Rating
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.RatingKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rating types.Rating
		if err := k.cdc.Unmarshal(iterator.Value(), &rating); err != nil {
			panic(fmt.Errorf("failed to unmarshal rating: %w", err))
		}
		ratings = append(ratings, rating)
	}
	return ratings
}

func (k Keeper) setRatingByTask(ctx sdk.Context, rating types.Rating) {
	store := k.getStore(ctx)
	key := types.GetRatingByTaskKey(rating.TaskId, rating.Id)
	store.Set(key, []byte(rating.Id))
}

func (k Keeper) setRatingByRatedUser(ctx sdk.Context, rating types.Rating) {
	store := k.getStore(ctx)
	key := types.GetRatingByRatedUserKey(rating.RatedId, rating.Id)
	store.Set(key, []byte(rating.Id))
}

// GetRatingsByTask returns ratings by task using composite index
func (k Keeper) GetRatingsByTask(ctx sdk.Context, taskID string) []types.Rating {
	var ratings []types.Rating
	store := k.getStore(ctx)
	prefix := types.GetRatingByTaskPrefix(taskID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ratingID := string(iterator.Value())
		rating, found := k.GetRating(ctx, ratingID)
		if found {
			ratings = append(ratings, rating)
		}
	}
	return ratings
}

// GetRatingsByRatedUser returns ratings by rated user using composite index
func (k Keeper) GetRatingsByRatedUser(ctx sdk.Context, ratedID string) []types.Rating {
	var ratings []types.Rating
	store := k.getStore(ctx)
	prefix := types.GetRatingByRatedUserPrefix(ratedID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ratingID := string(iterator.Value())
		rating, found := k.GetRating(ctx, ratingID)
		if found {
			ratings = append(ratings, rating)
		}
	}
	return ratings
}

// Reputation CRUD Operations

// SetReputation stores a reputation in the KVStore
func (k Keeper) SetReputation(ctx sdk.Context, rep types.Reputation) {
	store := k.getStore(ctx)
	key := types.GetReputationKey(rep.UserId)
	value, err := k.cdc.Marshal(rep)
	if err != nil {
		panic(fmt.Errorf("failed to marshal reputation: %w", err))
	}
	store.Set(key, value)
}

// GetReputation retrieves a reputation by user ID from KVStore
func (k Keeper) GetReputation(ctx sdk.Context, userID string) (types.Reputation, bool) {
	store := k.getStore(ctx)
	key := types.GetReputationKey(userID)
	value := store.Get(key)
	if value == nil {
		return types.Reputation{}, false
	}

	var rep types.Reputation
	if err := k.cdc.Unmarshal(value, &rep); err != nil {
		panic(fmt.Errorf("failed to unmarshal reputation: %w", err))
	}
	return rep, true
}

// GetAllReputations returns all reputations from KVStore
func (k Keeper) GetAllReputations(ctx sdk.Context) []types.Reputation {
	var reps []types.Reputation
	store := k.getStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.ReputationKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rep types.Reputation
		if err := k.cdc.Unmarshal(iterator.Value(), &rep); err != nil {
			panic(fmt.Errorf("failed to unmarshal reputation: %w", err))
		}
		reps = append(reps, rep)
	}
	return reps
}

// GetTaskStatistics returns task statistics
func (k Keeper) GetTaskStatistics(ctx sdk.Context) map[string]interface{} {
	stats := map[string]interface{}{
		"total_tasks":        0,
		"open_tasks":         0,
		"assigned_tasks":     0,
		"in_progress_tasks":  0,
		"completed_tasks":    0,
		"cancelled_tasks":    0,
		"total_applications": 0,
		"total_bids":         0,
		"total_ratings":      0,
	}

	tasks := k.GetAllTasks(ctx)
	stats["total_tasks"] = len(tasks)

	for _, task := range tasks {
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

	stats["total_applications"] = len(k.GetAllApplications(ctx))
	stats["total_ratings"] = len(k.GetAllRatings(ctx))

	// Count bids from auctions
	auctions := k.GetAllAuctions(ctx)
	for _, auction := range auctions {
		stats["total_bids"] = stats["total_bids"].(int) + len(auction.Bids)
	}

	return stats
}
