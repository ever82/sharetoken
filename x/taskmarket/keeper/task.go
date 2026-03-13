package keeper

import (
	"fmt"

	"sharetoken/x/taskmarket/types"
)

// CreateTask creates a new task
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

	// Update composite indexes
	lk.tasksByRequester[task.RequesterID] = append(lk.tasksByRequester[task.RequesterID], task.ID)
	lk.tasksByStatus[string(task.Status)] = append(lk.tasksByStatus[string(task.Status)], task.ID)

	return nil
}

// GetTask retrieves a task by ID
func (lk *LegacyKeeper) GetTask(id string) *types.Task {
	return lk.tasks[id]
}

// UpdateTask updates a task
func (lk *LegacyKeeper) UpdateTask(task *types.Task) error {
	// Check if status changed to update status index
	if oldTask, exists := lk.tasks[task.ID]; exists && oldTask.Status != task.Status {
		// Remove from old status index
		lk.tasksByStatus[string(oldTask.Status)] = removeFromSlice(lk.tasksByStatus[string(oldTask.Status)], task.ID)
		// Add to new status index
		lk.tasksByStatus[string(task.Status)] = append(lk.tasksByStatus[string(task.Status)], task.ID)
	}

	// Check if worker changed to update worker index
	if oldTask, exists := lk.tasks[task.ID]; exists && oldTask.WorkerID != task.WorkerID {
		// Remove from old worker index if old worker existed
		if oldTask.WorkerID != "" {
			lk.tasksByWorker[oldTask.WorkerID] = removeFromSlice(lk.tasksByWorker[oldTask.WorkerID], task.ID)
		}
		// Add to new worker index if new worker is set
		if task.WorkerID != "" {
			lk.tasksByWorker[task.WorkerID] = append(lk.tasksByWorker[task.WorkerID], task.ID)
		}
	}

	lk.tasks[task.ID] = task
	return nil
}

// GetAllTasks returns all tasks
func (lk *LegacyKeeper) GetAllTasks() []*types.Task {
	result := make([]*types.Task, 0, len(lk.tasks))
	for _, task := range lk.tasks {
		result = append(result, task)
	}
	return result
}

// GetTasksByRequester returns tasks by requester using composite index
func (lk *LegacyKeeper) GetTasksByRequester(requesterID string) []*types.Task {
	tasks := make([]*types.Task, 0)
	if taskIDs, exists := lk.tasksByRequester[requesterID]; exists {
		tasks = make([]*types.Task, 0, len(taskIDs))
		for _, taskID := range taskIDs {
			if task, ok := lk.tasks[taskID]; ok {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

// GetTasksByWorker returns tasks by worker using composite index
func (lk *LegacyKeeper) GetTasksByWorker(workerID string) []*types.Task {
	tasks := make([]*types.Task, 0)
	if taskIDs, exists := lk.tasksByWorker[workerID]; exists {
		tasks = make([]*types.Task, 0, len(taskIDs))
		for _, taskID := range taskIDs {
			if task, ok := lk.tasks[taskID]; ok {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

// GetOpenTasks returns open tasks
func (lk *LegacyKeeper) GetOpenTasks() []*types.Task {
	var open []*types.Task
	for _, task := range lk.tasks {
		if task.IsOpen() {
			open = append(open, task)
		}
	}
	return open
}

// StartTask starts a task
func (lk *LegacyKeeper) StartTask(taskID string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Update status index
	lk.tasksByStatus[string(task.Status)] = removeFromSlice(lk.tasksByStatus[string(task.Status)], task.ID)

	task.Start()
	if len(task.Milestones) > 0 {
		for i := range task.Milestones {
			if task.Milestones[i].Status == types.MilestoneStatusPending {
				task.Milestones[i].Status = types.MilestoneStatusActive
				break
			}
		}
	}

	// Update new status index
	lk.tasksByStatus[string(task.Status)] = append(lk.tasksByStatus[string(task.Status)], task.ID)
	return nil
}

// SubmitMilestone submits a milestone
func (lk *LegacyKeeper) SubmitMilestone(taskID, milestoneID, deliverables string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	return task.SubmitMilestone(milestoneID, deliverables)
}

// ApproveMilestone approves a milestone
func (lk *LegacyKeeper) ApproveMilestone(taskID, milestoneID string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Update status index
	lk.tasksByStatus[string(task.Status)] = removeFromSlice(lk.tasksByStatus[string(task.Status)], task.ID)

	err := task.ApproveMilestone(milestoneID)
	if task.AllMilestonesCompleted() {
		task.Complete()
	}

	// Update new status index
	lk.tasksByStatus[string(task.Status)] = append(lk.tasksByStatus[string(task.Status)], task.ID)
	return err
}

// RejectMilestone rejects a milestone
func (lk *LegacyKeeper) RejectMilestone(taskID, milestoneID, reason string) error {
	task := lk.GetTask(taskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	return task.RejectMilestone(milestoneID, reason)
}

// GetTaskStatistics returns task statistics
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

// removeFromSlice removes a string from a slice (helper function)
func removeFromSlice(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
