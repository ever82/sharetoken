package types

import (
	"errors"
)

// Constants for task types and statuses
const (
	// Task types
	TaskTypeOpen    = TaskType_TASK_TYPE_OPEN
	TaskTypeAuction = TaskType_TASK_TYPE_AUCTION

	// Application statuses
	ApplicationStatusPending   = ApplicationStatus_APPLICATION_STATUS_PENDING
	ApplicationStatusAccepted  = ApplicationStatus_APPLICATION_STATUS_ACCEPTED
	ApplicationStatusRejected  = ApplicationStatus_APPLICATION_STATUS_REJECTED
	ApplicationStatusWithdrawn = ApplicationStatus_APPLICATION_STATUS_WITHDRAWN

	// Rating dimensions
	DimensionQuality         = "quality"
	DimensionCommunication   = "communication"
	DimensionTimeliness      = "timeliness"
	DimensionProfessionalism = "professionalism"
)

// Validate validates an application
func (a Application) Validate() error {
	if a.Id == "" {
		return errors.New("application ID cannot be empty")
	}
	if a.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	if a.WorkerId == "" {
		return errors.New("worker ID cannot be empty")
	}
	if a.ProposedPrice == 0 {
		return errors.New("proposed price must be greater than 0")
	}
	return nil
}

// IsOpen checks if a task is open for applications
func (t Task) IsOpen() bool {
	return t.Status == TaskStatus_TASK_STATUS_OPEN
}

// Accept marks the application as accepted
func (a *Application) Accept() {
	a.Status = ApplicationStatus_APPLICATION_STATUS_ACCEPTED
}

// Reject marks the application as rejected
func (a *Application) Reject() {
	a.Status = ApplicationStatus_APPLICATION_STATUS_REJECTED
}

// Withdraw marks the application as withdrawn
func (a *Application) Withdraw() {
	a.Status = ApplicationStatus_APPLICATION_STATUS_WITHDRAWN
}

// Assign assigns a worker to a task
func (t *Task) Assign(workerID string) {
	t.WorkerId = workerID
	t.Status = TaskStatus_TASK_STATUS_ASSIGNED
}
