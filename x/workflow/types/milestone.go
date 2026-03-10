package types

import (
	"fmt"
	"time"
)

// MilestoneStatus represents the status of a milestone
type MilestoneStatus string

const (
	MilestoneStatusPending   MilestoneStatus = "pending"
	MilestoneStatusActive    MilestoneStatus = "active"
	MilestoneStatusCompleted MilestoneStatus = "completed"
	MilestoneStatusFailed    MilestoneStatus = "failed"
	MilestoneStatusPaid      MilestoneStatus = "paid"
)

// Milestone represents a payment milestone in a workflow
type Milestone struct {
	ID          string          `json:"id"`
	WorkflowID  string          `json:"workflow_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Amount      uint64          `json:"amount"` // Amount in STT
	Status      MilestoneStatus `json:"status"`
	Order       int             `json:"order"`     // Execution order
	Criteria    string          `json:"criteria"`  // Completion criteria
	EscrowID    string          `json:"escrow_id"` // Associated escrow
	CreatedAt   int64           `json:"created_at"`
	CompletedAt int64           `json:"completed_at"`
	PaidAt      int64           `json:"paid_at"`
}

// NewMilestone creates a new milestone
func NewMilestone(id, workflowID, name string, amount uint64, order int) *Milestone {
	return &Milestone{
		ID:         id,
		WorkflowID: workflowID,
		Name:       name,
		Amount:     amount,
		Status:     MilestoneStatusPending,
		Order:      order,
		CreatedAt:  time.Now().Unix(),
	}
}

// Activate marks the milestone as active
func (m *Milestone) Activate() {
	m.Status = MilestoneStatusActive
}

// Complete marks the milestone as completed
func (m *Milestone) Complete() {
	m.Status = MilestoneStatusCompleted
	m.CompletedAt = time.Now().Unix()
}

// MarkPaid marks the milestone as paid
func (m *Milestone) MarkPaid() {
	m.Status = MilestoneStatusPaid
	m.PaidAt = time.Now().Unix()
}

// Fail marks the milestone as failed
func (m *Milestone) Fail() {
	m.Status = MilestoneStatusFailed
}

// Validate validates milestone configuration
func (m *Milestone) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("milestone ID cannot be empty")
	}
	if m.WorkflowID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}
	if m.Amount == 0 {
		return fmt.Errorf("milestone amount must be greater than 0")
	}
	if m.Order < 0 {
		return fmt.Errorf("milestone order must be non-negative")
	}
	return nil
}

// MilestonePlan represents a complete milestone plan for a workflow
type MilestonePlan struct {
	ID          string      `json:"id"`
	WorkflowID  string      `json:"workflow_id"`
	TotalAmount uint64      `json:"total_amount"`
	Milestones  []Milestone `json:"milestones"`
	CreatedAt   int64       `json:"created_at"`
}

// NewMilestonePlan creates a new milestone plan
func NewMilestonePlan(id, workflowID string, totalAmount uint64) *MilestonePlan {
	return &MilestonePlan{
		ID:          id,
		WorkflowID:  workflowID,
		TotalAmount: totalAmount,
		Milestones:  []Milestone{},
		CreatedAt:   time.Now().Unix(),
	}
}

// AddMilestone adds a milestone to the plan
func (mp *MilestonePlan) AddMilestone(milestone Milestone) {
	mp.Milestones = append(mp.Milestones, milestone)
}

// GetNextPending returns the next pending milestone
func (mp *MilestonePlan) GetNextPending() *Milestone {
	for i := range mp.Milestones {
		if mp.Milestones[i].Status == MilestoneStatusPending {
			return &mp.Milestones[i]
		}
	}
	return nil
}

// GetActive returns the currently active milestone
func (mp *MilestonePlan) GetActive() *Milestone {
	for i := range mp.Milestones {
		if mp.Milestones[i].Status == MilestoneStatusActive {
			return &mp.Milestones[i]
		}
	}
	return nil
}

// Validate validates the milestone plan
func (mp *MilestonePlan) Validate() error {
	if mp.ID == "" {
		return fmt.Errorf("milestone plan ID cannot be empty")
	}

	// Check total amount matches sum of milestones
	var sum uint64
	for _, m := range mp.Milestones {
		sum += m.Amount
	}
	if sum != mp.TotalAmount {
		return fmt.Errorf("milestone amounts sum (%d) does not equal total (%d)", sum, mp.TotalAmount)
	}

	// Check for duplicate orders
	orders := make(map[int]bool)
	for _, m := range mp.Milestones {
		if orders[m.Order] {
			return fmt.Errorf("duplicate milestone order: %d", m.Order)
		}
		orders[m.Order] = true
	}

	return nil
}

// GetCompletedAmount returns total amount of completed/paid milestones
func (mp *MilestonePlan) GetCompletedAmount() uint64 {
	var sum uint64
	for _, m := range mp.Milestones {
		if m.Status == MilestoneStatusCompleted || m.Status == MilestoneStatusPaid {
			sum += m.Amount
		}
	}
	return sum
}

// GetPendingAmount returns total amount of pending/active milestones
func (mp *MilestonePlan) GetPendingAmount() uint64 {
	var sum uint64
	for _, m := range mp.Milestones {
		if m.Status == MilestoneStatusPending || m.Status == MilestoneStatusActive {
			sum += m.Amount
		}
	}
	return sum
}

// ProgressReport represents a progress report for a workflow
type ProgressReport struct {
	WorkflowID      string  `json:"workflow_id"`
	TotalSteps      int     `json:"total_steps"`
	CompletedSteps  int     `json:"completed_steps"`
	FailedSteps     int     `json:"failed_steps"`
	CurrentStep     string  `json:"current_step"`
	TotalMilestones int     `json:"total_milestones"`
	CompletedMs     int     `json:"completed_milestones"`
	PaidMs          int     `json:"paid_milestones"`
	ProgressPct     float64 `json:"progress_percentage"`
	EstimatedTime   int64   `json:"estimated_time_remaining"` // seconds
	GeneratedAt     int64   `json:"generated_at"`
}

// NewProgressReport creates a new progress report
func NewProgressReport(workflowID string) *ProgressReport {
	return &ProgressReport{
		WorkflowID:  workflowID,
		GeneratedAt: time.Now().Unix(),
	}
}

// CalculateProgress calculates progress percentage
func (pr *ProgressReport) CalculateProgress() {
	if pr.TotalSteps > 0 {
		pr.ProgressPct = float64(pr.CompletedSteps) / float64(pr.TotalSteps) * 100
	}
}
