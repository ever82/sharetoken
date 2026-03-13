package keeper

import (
	"encoding/json"
	"fmt"

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
	storeKey   StoreKey
	cdc        codec
	paramSpace types.ParamSubspace
}

// NewKeeper creates a new task marketplace keeper
func NewKeeper(storeKey StoreKey, cdc codec, paramSpace types.ParamSubspace) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramSpace: paramSpace,
	}
}

// NewKeeperWithDefaultCodec creates a keeper with default JSON codec
func NewKeeperWithDefaultCodec(storeKey StoreKey) *Keeper {
	return &Keeper{
		storeKey: storeKey,
		cdc:      JSONCodec{},
	}
}

// NewKeeperSimple creates a new task marketplace keeper without dependencies (for tests)
func NewKeeperSimple() *Keeper {
	return &Keeper{
		storeKey: []byte(types.StoreKey),
		cdc:      JSONCodec{},
	}
}

// GetStoreKey returns the store key
func (k Keeper) GetStoreKey() StoreKey {
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
	key := types.GetTaskKey(task.ID)
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
	if k.HasTask(ctx, task.ID) {
		return fmt.Errorf("task already exists: %s", task.ID)
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
	if !k.HasTask(ctx, task.ID) {
		return fmt.Errorf("task not found: %s", task.ID)
	}

	// Delete old indexes
	oldTask, _ := k.GetTask(ctx, task.ID)
	k.deleteTaskIndexes(ctx, oldTask)

	// Set new task with updated indexes
	k.SetTask(ctx, task)
	return nil
}

// Index management functions

func (k Keeper) setTaskByRequester(ctx sdk.Context, task types.Task) {
	if task.RequesterID == "" {
		return
	}
	store := k.getStore(ctx)
	key := types.GetTaskByRequesterKey(task.RequesterID, task.ID)
	store.Set(key, []byte(task.ID))
}

func (k Keeper) setTaskByWorker(ctx sdk.Context, task types.Task) {
	if task.WorkerID == "" {
		return
	}
	store := k.getStore(ctx)
	key := types.GetTaskByWorkerKey(task.WorkerID, task.ID)
	store.Set(key, []byte(task.ID))
}

func (k Keeper) setTaskByStatus(ctx sdk.Context, task types.Task) {
	store := k.getStore(ctx)
	key := types.GetTaskByStatusKey(string(task.Status), task.ID)
	store.Set(key, []byte(task.ID))
}

func (k Keeper) deleteTaskIndexes(ctx sdk.Context, task types.Task) {
	store := k.getStore(ctx)

	// Delete requester index
	if task.RequesterID != "" {
		key := types.GetTaskByRequesterKey(task.RequesterID, task.ID)
		store.Delete(key)
	}

	// Delete worker index
	if task.WorkerID != "" {
		key := types.GetTaskByWorkerKey(task.WorkerID, task.ID)
		store.Delete(key)
	}

	// Delete status index
	key := types.GetTaskByStatusKey(string(task.Status), task.ID)
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
