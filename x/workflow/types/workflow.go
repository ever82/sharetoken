package types

import (
	"fmt"
	"time"
)

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusPending    WorkflowStatus = "pending"
	WorkflowStatusRunning    WorkflowStatus = "running"
	WorkflowStatusPaused     WorkflowStatus = "paused"
	WorkflowStatusCompleted  WorkflowStatus = "completed"
	WorkflowStatusFailed     WorkflowStatus = "failed"
	WorkflowStatusCancelled  WorkflowStatus = "cancelled"
)

// StepStatus represents the status of a workflow step
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusReady     StepStatus = "ready"     // Dependencies satisfied
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
)

// StepType represents the type of workflow step
type StepType string

const (
	StepTypeAgent      StepType = "agent"      // Execute agent task
	StepTypeParallel   StepType = "parallel"   // Run steps in parallel
	StepTypeSequence   StepType = "sequence"   // Run steps sequentially
	StepTypeCondition  StepType = "condition"  // Conditional execution
	StepTypeMilestone  StepType = "milestone"  // Payment milestone
)

// Workflow represents a multi-agent workflow
type Workflow struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Owner       string         `json:"owner"`
	Status      WorkflowStatus `json:"status"`
	Steps       []WorkflowStep `json:"steps"`
	Inputs      map[string]interface{} `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs"`
	CreatedAt   int64          `json:"created_at"`
	StartedAt   int64          `json:"started_at"`
	CompletedAt int64          `json:"completed_at"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Type         StepType          `json:"type"`
	Capability   Capability        `json:"capability"`
	DependsOn    []string          `json:"depends_on"`    // IDs of dependent steps
	Params       map[string]string `json:"params"`
	Status       StepStatus        `json:"status"`
	Output       string            `json:"output"`
	Error        string            `json:"error"`
	StartedAt    int64             `json:"started_at"`
	CompletedAt  int64             `json:"completed_at"`
	SubSteps     []WorkflowStep    `json:"sub_steps"`     // For parallel/sequence
	Condition    string            `json:"condition"`     // For condition step
	MilestoneID  string            `json:"milestone_id"`  // For milestone step
}

// NewWorkflow creates a new workflow
func NewWorkflow(id, name, owner string) *Workflow {
	return &Workflow{
		ID:        id,
		Name:      name,
		Owner:     owner,
		Status:    WorkflowStatusPending,
		Steps:     []WorkflowStep{},
		Inputs:    make(map[string]interface{}),
		Outputs:   make(map[string]interface{}),
		CreatedAt: time.Now().Unix(),
	}
}

// Start marks the workflow as started
func (w *Workflow) Start() {
	w.Status = WorkflowStatusRunning
	w.StartedAt = time.Now().Unix()
}

// Complete marks the workflow as completed
func (w *Workflow) Complete() {
	w.Status = WorkflowStatusCompleted
	w.CompletedAt = time.Now().Unix()
}

// Fail marks the workflow as failed
func (w *Workflow) Fail(err string) {
	w.Status = WorkflowStatusFailed
	w.CompletedAt = time.Now().Unix()
}

// Cancel marks the workflow as cancelled
func (w *Workflow) Cancel() {
	w.Status = WorkflowStatusCancelled
	w.CompletedAt = time.Now().Unix()
}

// AddStep adds a step to the workflow
func (w *Workflow) AddStep(step WorkflowStep) {
	w.Steps = append(w.Steps, step)
}

// GetStep returns a step by ID
func (w *Workflow) GetStep(id string) *WorkflowStep {
	for i := range w.Steps {
		if w.Steps[i].ID == id {
			return &w.Steps[i]
		}
	}
	return nil
}

// GetReadySteps returns steps that are ready to execute
func (w *Workflow) GetReadySteps() []*WorkflowStep {
	var ready []*WorkflowStep

	for i := range w.Steps {
		if w.Steps[i].Status != StepStatusPending {
			continue
		}

		// Check if all dependencies are completed
		allDepsComplete := true
		for _, depID := range w.Steps[i].DependsOn {
			dep := w.GetStep(depID)
			if dep == nil || dep.Status != StepStatusCompleted {
				allDepsComplete = false
				break
			}
		}

		if allDepsComplete {
			ready = append(ready, &w.Steps[i])
		}
	}

	return ready
}

// IsComplete checks if all steps are completed
func (w *Workflow) IsComplete() bool {
	for _, step := range w.Steps {
		if step.Status != StepStatusCompleted && step.Status != StepStatusSkipped {
			return false
		}
	}
	return true
}

// Validate validates the workflow structure
func (w *Workflow) Validate() error {
	if w.ID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}
	if w.Owner == "" {
		return fmt.Errorf("workflow owner cannot be empty")
	}

	// Validate no circular dependencies
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, step := range w.Steps {
		if step.ID == "" {
			return fmt.Errorf("step ID cannot be empty")
		}
		if hasCycle(w, step.ID, visited, recStack) {
			return fmt.Errorf("circular dependency detected at step %s", step.ID)
		}
	}

	return nil
}

// hasCycle detects circular dependencies using DFS
func hasCycle(w *Workflow, stepID string, visited, recStack map[string]bool) bool {
	visited[stepID] = true
	recStack[stepID] = true

	step := w.GetStep(stepID)
	if step == nil {
		return false
	}

	for _, depID := range step.DependsOn {
		if !visited[depID] {
			if hasCycle(w, depID, visited, recStack) {
				return true
			}
		} else if recStack[depID] {
			return true
		}
	}

	recStack[stepID] = false
	return false
}
