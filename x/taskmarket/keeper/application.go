package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// SubmitApplication submits an application
func (k Keeper) SubmitApplication(ctx sdk.Context, app types.Application) error {
	if err := app.Validate(); err != nil {
		return fmt.Errorf("invalid application: %w", err)
	}
	task, found := k.GetTask(ctx, app.TaskID)
	if !found {
		return fmt.Errorf("task not found: %s", app.TaskID)
	}
	if task.Type != types.TaskTypeOpen {
		return fmt.Errorf("task is not open type")
	}
	if !task.IsOpen() {
		return fmt.Errorf("task is not accepting applications")
	}

	// Check for existing application from this worker
	existingApps := k.GetApplicationsByTask(ctx, app.TaskID)
	for _, existing := range existingApps {
		if existing.WorkerID == app.WorkerID {
			return fmt.Errorf("worker already applied to this task")
		}
	}

	k.SetApplication(ctx, app)

	// Update task application count
	task.ApplicationCount++
	k.SetTask(ctx, task)

	return nil
}

// AcceptApplication accepts an application
func (k Keeper) AcceptApplication(ctx sdk.Context, appID string) error {
	app, found := k.GetApplication(ctx, appID)
	if !found {
		return fmt.Errorf("application not found: %s", appID)
	}
	if app.Status != types.ApplicationStatusPending {
		return fmt.Errorf("application is not pending")
	}

	task, found := k.GetTask(ctx, app.TaskID)
	if !found {
		return fmt.Errorf("task not found: %s", app.TaskID)
	}

	// Delete old indexes
	k.deleteTaskIndexes(ctx, task)

	app.Accept()
	task.Assign(app.WorkerID)

	// Update new indexes
	k.setTaskByWorker(ctx, task)
	k.setTaskByStatus(ctx, task)
	k.SetTask(ctx, task)
	k.SetApplication(ctx, app)

	return nil
}

// RejectApplication rejects an application
func (k Keeper) RejectApplication(ctx sdk.Context, appID string) error {
	app, found := k.GetApplication(ctx, appID)
	if !found {
		return fmt.Errorf("application not found: %s", appID)
	}
	app.Reject()
	k.SetApplication(ctx, app)
	return nil
}

// WithdrawApplication withdraws an application
func (k Keeper) WithdrawApplication(ctx sdk.Context, appID string) error {
	app, found := k.GetApplication(ctx, appID)
	if !found {
		return fmt.Errorf("application not found: %s", appID)
	}
	app.Withdraw()
	k.SetApplication(ctx, app)
	return nil
}
