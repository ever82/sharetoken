package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"sharetoken/x/workflow/types"
)

// Executor manages workflow execution
type Executor struct {
	workflows     map[string]*types.Workflow
	milestones    map[string]*types.MilestonePlan
	capabilityMgr *CapabilityManager
	mutex         sync.RWMutex
	eventCh       chan WorkflowEvent
}

// WorkflowEvent represents a workflow event
type WorkflowEvent struct {
	WorkflowID string
	StepID     string
	Type       EventType
	Payload    interface{}
	Timestamp  int64
}

// EventType represents the type of workflow event
type EventType string

const (
	EventWorkflowStarted   EventType = "workflow_started"
	EventWorkflowCompleted EventType = "workflow_completed"
	EventWorkflowFailed    EventType = "workflow_failed"
	EventStepStarted       EventType = "step_started"
	EventStepCompleted     EventType = "step_completed"
	EventStepFailed        EventType = "step_failed"
	EventMilestoneReached  EventType = "milestone_reached"
	EventMilestonePaid     EventType = "milestone_paid"
)

// NewExecutor creates a new workflow executor
func NewExecutor() *Executor {
	return &Executor{
		workflows:     make(map[string]*types.Workflow),
		milestones:    make(map[string]*types.MilestonePlan),
		capabilityMgr: NewCapabilityManager(),
		eventCh:       make(chan WorkflowEvent, 100),
	}
}

// RegisterWorkflow registers a workflow for execution
func (e *Executor) RegisterWorkflow(workflow *types.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return fmt.Errorf("invalid workflow: %w", err)
	}

	e.mutex.Lock()
	e.workflows[workflow.ID] = workflow
	e.mutex.Unlock()

	return nil
}

// RegisterMilestonePlan registers a milestone plan for a workflow
func (e *Executor) RegisterMilestonePlan(plan *types.MilestonePlan) error {
	if err := plan.Validate(); err != nil {
		return fmt.Errorf("invalid milestone plan: %w", err)
	}

	e.mutex.Lock()
	e.milestones[plan.WorkflowID] = plan
	e.mutex.Unlock()

	return nil
}

// ExecuteWorkflow executes a workflow
func (e *Executor) ExecuteWorkflow(ctx context.Context, workflowID string) error {
	e.mutex.RLock()
	workflow, exists := e.workflows[workflowID]
	milestones := e.milestones[workflowID]
	e.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("workflow not found: %s", workflowID)
	}

	// Start workflow
	workflow.Start()
	e.emitEvent(WorkflowEvent{
		WorkflowID: workflowID,
		Type:       EventWorkflowStarted,
		Timestamp:  time.Now().Unix(),
	})

	// Execute steps
	for {
		// Check context
		if ctx.Err() != nil {
			workflow.Fail("execution cancelled")
			e.emitEvent(WorkflowEvent{
				WorkflowID: workflowID,
				Type:       EventWorkflowFailed,
				Payload:    "context cancelled",
				Timestamp:  time.Now().Unix(),
			})
			return ctx.Err()
		}

		// Get ready steps
		readySteps := workflow.GetReadySteps()
		if len(readySteps) == 0 {
			// Check if workflow is complete
			if workflow.IsComplete() {
				workflow.Complete()
				e.emitEvent(WorkflowEvent{
					WorkflowID: workflowID,
					Type:       EventWorkflowCompleted,
					Timestamp:  time.Now().Unix(),
				})
				return nil
			}

			// Check for stuck workflow (no ready steps but not complete)
			if hasFailedSteps(workflow) {
				workflow.Fail("workflow has failed steps")
				e.emitEvent(WorkflowEvent{
					WorkflowID: workflowID,
					Type:       EventWorkflowFailed,
					Timestamp:  time.Now().Unix(),
				})
				return fmt.Errorf("workflow has failed steps")
			}

			// Wait a bit before checking again
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Execute ready steps in parallel
		var wg sync.WaitGroup
		errCh := make(chan error, len(readySteps))

		for _, step := range readySteps {
			wg.Add(1)
			go func(s *types.WorkflowStep) {
				defer wg.Done()
				if err := e.executeStep(ctx, workflow, s, milestones); err != nil {
					errCh <- fmt.Errorf("step %s failed: %w", s.ID, err)
				}
			}(step)
		}

		wg.Wait()
		close(errCh)

		// Check for errors
		for err := range errCh {
			_ = err // Ignore errors, let other steps complete
		}
	}
}

// executeStep executes a single workflow step
func (e *Executor) executeStep(ctx context.Context, workflow *types.Workflow, step *types.WorkflowStep, milestones *types.MilestonePlan) error {
	// Update step status
	step.Status = types.StepStatusRunning
	step.StartedAt = time.Now().Unix()

	e.emitEvent(WorkflowEvent{
		WorkflowID: workflow.ID,
		StepID:     step.ID,
		Type:       EventStepStarted,
		Timestamp:  step.StartedAt,
	})

	// Execute based on step type
	switch step.Type {
	case types.StepTypeAgent:
		return e.executeAgentStep(ctx, workflow, step, milestones)
	case types.StepTypeParallel:
		return e.executeParallelStep(ctx, workflow, step, milestones)
	case types.StepTypeSequence:
		return e.executeSequenceStep(ctx, workflow, step, milestones)
	case types.StepTypeCondition:
		return e.executeConditionStep(ctx, workflow, step, milestones)
	case types.StepTypeMilestone:
		return e.executeMilestoneStep(ctx, workflow, step, milestones)
	default:
		step.Status = types.StepStatusFailed
		step.Error = fmt.Sprintf("unknown step type: %s", step.Type)
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
}

// executeAgentStep executes an agent-based step
func (e *Executor) executeAgentStep(ctx context.Context, workflow *types.Workflow, step *types.WorkflowStep, milestones *types.MilestonePlan) error {
	// Get capability config
	config := types.DefaultCapabilityConfig(step.Capability)

	// Apply step-specific params
	for k, v := range step.Params {
		config.Params[k] = v
	}

	// Execute using capability manager
	result, err := e.capabilityMgr.Execute(ctx, step.Capability, config, step.Params)
	if err != nil {
		step.Status = types.StepStatusFailed
		step.Error = err.Error()
		e.emitEvent(WorkflowEvent{
			WorkflowID: workflow.ID,
			StepID:     step.ID,
			Type:       EventStepFailed,
			Payload:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return err
	}

	step.Status = types.StepStatusCompleted
	step.Output = result
	step.CompletedAt = time.Now().Unix()

	e.emitEvent(WorkflowEvent{
		WorkflowID: workflow.ID,
		StepID:     step.ID,
		Type:       EventStepCompleted,
		Payload:    result,
		Timestamp:  step.CompletedAt,
	})

	return nil
}

// executeParallelStep executes sub-steps in parallel
func (e *Executor) executeParallelStep(ctx context.Context, workflow *types.Workflow, step *types.WorkflowStep, milestones *types.MilestonePlan) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(step.SubSteps))

	for i := range step.SubSteps {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := e.executeStep(ctx, workflow, &step.SubSteps[idx], milestones); err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	// Check for errors
	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		step.Status = types.StepStatusFailed
		step.Error = fmt.Sprintf("%d sub-steps failed", len(errs))
		return fmt.Errorf("parallel step failed")
	}

	step.Status = types.StepStatusCompleted
	step.CompletedAt = time.Now().Unix()

	return nil
}

// executeSequenceStep executes sub-steps sequentially
func (e *Executor) executeSequenceStep(ctx context.Context, workflow *types.Workflow, step *types.WorkflowStep, milestones *types.MilestonePlan) error {
	for i := range step.SubSteps {
		if err := e.executeStep(ctx, workflow, &step.SubSteps[i], milestones); err != nil {
			step.Status = types.StepStatusFailed
			step.Error = fmt.Sprintf("sub-step %d failed: %v", i, err)
			return err
		}
	}

	step.Status = types.StepStatusCompleted
	step.CompletedAt = time.Now().Unix()

	return nil
}

// executeConditionStep executes a conditional step
func (e *Executor) executeConditionStep(ctx context.Context, workflow *types.Workflow, step *types.WorkflowStep, milestones *types.MilestonePlan) error {
	// Evaluate condition (simplified - in real implementation would use expression evaluator)
	// For now, always take first sub-step if condition not empty
	if step.Condition == "" || len(step.SubSteps) == 0 {
		step.Status = types.StepStatusSkipped
		step.CompletedAt = time.Now().Unix()
		return nil
	}

	// Execute the appropriate branch
	// Simplified: execute first sub-step
	if err := e.executeStep(ctx, workflow, &step.SubSteps[0], milestones); err != nil {
		step.Status = types.StepStatusFailed
		step.Error = err.Error()
		return err
	}

	step.Status = types.StepStatusCompleted
	step.CompletedAt = time.Now().Unix()

	return nil
}

// executeMilestoneStep handles milestone completion
func (e *Executor) executeMilestoneStep(ctx context.Context, workflow *types.Workflow, step *types.WorkflowStep, milestones *types.MilestonePlan) error {
	if milestones == nil {
		step.Status = types.StepStatusFailed
		step.Error = "no milestone plan registered"
		return fmt.Errorf("no milestone plan for workflow %s", workflow.ID)
	}

	// Mark milestone as completed
	for i := range milestones.Milestones {
		if milestones.Milestones[i].ID == step.MilestoneID {
			milestones.Milestones[i].Complete()

			e.emitEvent(WorkflowEvent{
				WorkflowID: workflow.ID,
				StepID:     step.ID,
				Type:       EventMilestoneReached,
				Payload:    step.MilestoneID,
				Timestamp:  time.Now().Unix(),
			})

			// TODO: Trigger escrow payment
			break
		}
	}

	step.Status = types.StepStatusCompleted
	step.CompletedAt = time.Now().Unix()

	return nil
}

// GetWorkflow returns a workflow by ID
func (e *Executor) GetWorkflow(id string) *types.Workflow {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.workflows[id]
}

// GetMilestonePlan returns a milestone plan by workflow ID
func (e *Executor) GetMilestonePlan(workflowID string) *types.MilestonePlan {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.milestones[workflowID]
}

// GetEventChannel returns the event channel for subscribing to events
func (e *Executor) GetEventChannel() <-chan WorkflowEvent {
	return e.eventCh
}

// emitEvent emits a workflow event
func (e *Executor) emitEvent(event WorkflowEvent) {
	select {
	case e.eventCh <- event:
	default:
		// Channel full, drop event
	}
}

// hasFailedSteps checks if workflow has any failed steps
func hasFailedSteps(w *types.Workflow) bool {
	for _, step := range w.Steps {
		if step.Status == types.StepStatusFailed {
			return true
		}
	}
	return false
}

// GenerateProgressReport generates a progress report for a workflow
func (e *Executor) GenerateProgressReport(workflowID string) (*types.ProgressReport, error) {
	e.mutex.RLock()
	workflow := e.workflows[workflowID]
	milestones := e.milestones[workflowID]
	e.mutex.RUnlock()

	if workflow == nil {
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	report := types.NewProgressReport(workflowID)
	report.TotalSteps = len(workflow.Steps)

	// Count steps
	for _, step := range workflow.Steps {
		switch step.Status {
		case types.StepStatusCompleted:
			report.CompletedSteps++
		case types.StepStatusFailed:
			report.FailedSteps++
		case types.StepStatusRunning:
			report.CurrentStep = step.Name
		}
	}

	// Add milestone info
	if milestones != nil {
		report.TotalMilestones = len(milestones.Milestones)
		for _, ms := range milestones.Milestones {
			switch ms.Status {
			case types.MilestoneStatusCompleted, types.MilestoneStatusPaid:
				report.CompletedMs++
			}
			if ms.Status == types.MilestoneStatusPaid {
				report.PaidMs++
			}
		}
	}

	report.CalculateProgress()

	return report, nil
}
