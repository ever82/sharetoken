package keeper

import (
	"fmt"

	"sharetoken/x/taskmarket/types"
)

// SubmitApplication submits an application
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

// GetApplication gets an application by ID
func (lk *LegacyKeeper) GetApplication(id string) *types.Application {
	return lk.applications[id]
}

// GetApplicationsByTask returns applications for a task
func (lk *LegacyKeeper) GetApplicationsByTask(taskID string) []*types.Application {
	var apps []*types.Application
	for _, app := range lk.applications {
		if app.TaskID == taskID {
			apps = append(apps, app)
		}
	}
	return apps
}

// AcceptApplication accepts an application
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

	// Update status indexes
	lk.tasksByStatus[string(task.Status)] = removeFromSlice(lk.tasksByStatus[string(task.Status)], task.ID)

	app.Accept()
	task.Assign(app.WorkerID)

	// Update new status
	lk.tasksByStatus[string(task.Status)] = append(lk.tasksByStatus[string(task.Status)], task.ID)
	// Update worker index
	lk.tasksByWorker[task.WorkerID] = append(lk.tasksByWorker[task.WorkerID], task.ID)
	return nil
}

// RejectApplication rejects an application
func (lk *LegacyKeeper) RejectApplication(appID string) error {
	app := lk.GetApplication(appID)
	if app == nil {
		return fmt.Errorf("application not found: %s", appID)
	}
	app.Reject()
	return nil
}
