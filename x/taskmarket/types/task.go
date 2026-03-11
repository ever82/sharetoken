package types

import (
	"fmt"
	"time"
)

// TaskType represents the type of marketplace task
type TaskType string

const (
	TaskTypeOpen    TaskType = "open"    // Open application
	TaskTypeAuction TaskType = "auction" // Bidding auction
)

// TaskStatus represents the status of a marketplace task
type TaskStatus string

const (
	TaskStatusDraft      TaskStatus = "draft"
	TaskStatusOpen       TaskStatus = "open"     // Accepting applications/bids
	TaskStatusAssigned   TaskStatus = "assigned" // Worker selected
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReview     TaskStatus = "review" // Under review
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusDisputed   TaskStatus = "disputed"
)

// TaskCategory represents task categories
type TaskCategory string

const (
	CategoryDevelopment TaskCategory = "development"
	CategoryDesign      TaskCategory = "design"
	CategoryWriting     TaskCategory = "writing"
	CategoryResearch    TaskCategory = "research"
	CategoryMarketing   TaskCategory = "marketing"
	CategoryConsulting  TaskCategory = "consulting"
	CategoryOther       TaskCategory = "other"
)

// Task represents a marketplace task
type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	RequesterID string       `json:"requester_id"`
	WorkerID    string       `json:"worker_id"`
	Type        TaskType     `json:"type"`
	Category    TaskCategory `json:"category"`
	Status      TaskStatus   `json:"status"`
	Budget      uint64       `json:"budget"`   // In STT
	Currency    string       `json:"currency"` // "STT", "USDC"
	Deadline    int64        `json:"deadline"` // Unix timestamp
	Skills      []string     `json:"skills"`   // Required skills

	// Subtasks for decomposition
	Subtasks []Subtask `json:"subtasks"`

	// Milestones
	Milestones []Milestone `json:"milestones"`

	// Metadata
	CreatedAt   int64 `json:"created_at"`
	UpdatedAt   int64 `json:"updated_at"`
	AssignedAt  int64 `json:"assigned_at"`
	CompletedAt int64 `json:"completed_at"`

	// Statistics
	ViewCount        int `json:"view_count"`
	ApplicationCount int `json:"application_count"`
	BidCount         int `json:"bid_count"`
}

// Subtask represents a decomposed subtask
type Subtask struct {
	ID          string `json:"id"`
	TaskID      string `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Budget      uint64 `json:"budget"`
	Completed   bool   `json:"completed"`
}

// Milestone represents a task milestone
type Milestone struct {
	ID           string          `json:"id"`
	TaskID       string          `json:"task_id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Amount       uint64          `json:"amount"`
	Order        int             `json:"order"`
	Status       MilestoneStatus `json:"status"`
	Deliverables string          `json:"deliverables"`
	CreatedAt    int64           `json:"created_at"`
	CompletedAt  int64           `json:"completed_at"`
}

// MilestoneStatus represents milestone status
type MilestoneStatus string

const (
	MilestoneStatusPending   MilestoneStatus = "pending"
	MilestoneStatusActive    MilestoneStatus = "active"
	MilestoneStatusSubmitted MilestoneStatus = "submitted"
	MilestoneStatusApproved  MilestoneStatus = "approved"
	MilestoneStatusRejected  MilestoneStatus = "rejected"
	MilestoneStatusPaid      MilestoneStatus = "paid"
)

// NewTask creates a new marketplace task
func NewTask(id, title, description, requesterID string, taskType TaskType, budget uint64) *Task {
	now := time.Now().Unix()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		RequesterID: requesterID,
		Type:        taskType,
		Status:      TaskStatusDraft,
		Budget:      budget,
		Currency:    "STT",
		Skills:      []string{},
		Subtasks:    []Subtask{},
		Milestones:  []Milestone{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddSubtask adds a subtask
func (t *Task) AddSubtask(subtask Subtask) {
	t.Subtasks = append(t.Subtasks, subtask)
	t.UpdatedAt = time.Now().Unix()
}

// AddMilestone adds a milestone
func (t *Task) AddMilestone(milestone Milestone) {
	t.Milestones = append(t.Milestones, milestone)
	t.UpdatedAt = time.Now().Unix()
}

// Publish publishes the task
func (t *Task) Publish() {
	t.Status = TaskStatusOpen
	t.UpdatedAt = time.Now().Unix()
}

// Assign assigns the task to a worker
func (t *Task) Assign(workerID string) {
	t.WorkerID = workerID
	t.Status = TaskStatusAssigned
	t.AssignedAt = time.Now().Unix()
	t.UpdatedAt = t.AssignedAt
}

// Start marks the task as in progress
func (t *Task) Start() {
	t.Status = TaskStatusInProgress
	t.UpdatedAt = time.Now().Unix()
}

// SubmitForReview marks the task for review
func (t *Task) SubmitForReview() {
	t.Status = TaskStatusReview
	t.UpdatedAt = time.Now().Unix()
}

// Complete marks the task as completed
func (t *Task) Complete() {
	t.Status = TaskStatusCompleted
	t.CompletedAt = time.Now().Unix()
	t.UpdatedAt = t.CompletedAt
}

// Cancel cancels the task
func (t *Task) Cancel() {
	t.Status = TaskStatusCancelled
	t.UpdatedAt = time.Now().Unix()
}

// IsOpen checks if task is accepting applications/bids
func (t *Task) IsOpen() bool {
	return t.Status == TaskStatusOpen
}

// Validate validates the task
func (t *Task) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if t.Title == "" {
		return fmt.Errorf("task title cannot be empty")
	}
	if t.RequesterID == "" {
		return fmt.Errorf("requester ID cannot be empty")
	}
	if t.Budget == 0 {
		return fmt.Errorf("task budget must be greater than 0")
	}
	if t.Deadline > 0 && t.Deadline < time.Now().Unix() {
		return fmt.Errorf("deadline must be in the future")
	}
	return nil
}

// GetTotalMilestoneAmount returns sum of milestone amounts
func (t *Task) GetTotalMilestoneAmount() uint64 {
	var total uint64
	for _, m := range t.Milestones {
		total += m.Amount
	}
	return total
}

// ValidateMilestones checks if milestones total matches budget
func (t *Task) ValidateMilestones() error {
	total := t.GetTotalMilestoneAmount()
	if total != t.Budget {
		return fmt.Errorf("milestone total (%d) does not match budget (%d)", total, t.Budget)
	}
	return nil
}

// GetActiveMilestone returns the currently active milestone
func (t *Task) GetActiveMilestone() *Milestone {
	for i := range t.Milestones {
		if t.Milestones[i].Status == MilestoneStatusActive {
			return &t.Milestones[i]
		}
	}
	return nil
}

// GetNextPendingMilestone returns next pending milestone
func (t *Task) GetNextPendingMilestone() *Milestone {
	for i := range t.Milestones {
		if t.Milestones[i].Status == MilestoneStatusPending {
			return &t.Milestones[i]
		}
	}
	return nil
}

// ActivateMilestone activates a milestone
func (t *Task) ActivateMilestone(milestoneID string) error {
	for i := range t.Milestones {
		if t.Milestones[i].ID == milestoneID {
			t.Milestones[i].Status = MilestoneStatusActive
			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// SubmitMilestone submits a milestone for approval
func (t *Task) SubmitMilestone(milestoneID string, deliverables string) error {
	for i := range t.Milestones {
		if t.Milestones[i].ID == milestoneID {
			if t.Milestones[i].Status != MilestoneStatusActive {
				return fmt.Errorf("milestone is not active")
			}
			t.Milestones[i].Status = MilestoneStatusSubmitted
			t.Milestones[i].Deliverables = deliverables
			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// ApproveMilestone approves a milestone
func (t *Task) ApproveMilestone(milestoneID string) error {
	for i := range t.Milestones {
		if t.Milestones[i].ID == milestoneID {
			if t.Milestones[i].Status != MilestoneStatusSubmitted {
				return fmt.Errorf("milestone is not submitted")
			}
			t.Milestones[i].Status = MilestoneStatusApproved
			t.Milestones[i].CompletedAt = time.Now().Unix()

			// Activate next milestone if exists
			next := t.GetNextPendingMilestone()
			if next != nil {
				next.Status = MilestoneStatusActive
			}

			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// RejectMilestone rejects a milestone
func (t *Task) RejectMilestone(milestoneID string, reason string) error {
	for i := range t.Milestones {
		if t.Milestones[i].ID == milestoneID {
			if t.Milestones[i].Status != MilestoneStatusSubmitted {
				return fmt.Errorf("milestone is not submitted")
			}
			t.Milestones[i].Status = MilestoneStatusRejected
			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// AllMilestonesCompleted checks if all milestones are approved
func (t *Task) AllMilestonesCompleted() bool {
	if len(t.Milestones) == 0 {
		return false
	}
	for _, m := range t.Milestones {
		if m.Status != MilestoneStatusApproved && m.Status != MilestoneStatusPaid {
			return false
		}
	}
	return true
}

// GetCompletionPercentage returns task completion percentage
func (t *Task) GetCompletionPercentage() float64 {
	if len(t.Milestones) == 0 {
		if t.Status == TaskStatusCompleted {
			return 100.0
		}
		return 0.0
	}

	completed := 0
	for _, m := range t.Milestones {
		if m.Status == MilestoneStatusApproved || m.Status == MilestoneStatusPaid {
			completed++
		}
	}

	return float64(completed) / float64(len(t.Milestones)) * 100
}
