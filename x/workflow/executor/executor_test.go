package executor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sharetoken/x/workflow/types"
)

func TestCapabilityTypes(t *testing.T) {
	caps := types.GetAllCapabilities()
	require.Len(t, caps, 7)

	// Check key capabilities exist
	capMap := make(map[string]bool)
	for _, c := range caps {
		capMap[string(c)] = true
	}

	require.True(t, capMap["collector"])
	require.True(t, capMap["lead"])
	require.True(t, capMap["researcher"])
	require.True(t, capMap["writer"])
	require.True(t, capMap["analyst"])
	require.True(t, capMap["tester"])
	require.True(t, capMap["reviewer"])
}

func TestCapabilityConfig(t *testing.T) {
	config := types.DefaultCapabilityConfig(types.CapabilityWriter)
	require.Equal(t, "writer", config.Name)
	require.Greater(t, config.Timeout, int64(0))
	require.NotNil(t, config.Params)
}

func TestWorkflowCreation(t *testing.T) {
	workflow := types.NewWorkflow("wf-1", "Test Workflow", "owner1")

	require.Equal(t, "wf-1", workflow.ID)
	require.Equal(t, "Test Workflow", workflow.Name)
	require.Equal(t, "owner1", workflow.Owner)
	require.Equal(t, types.WorkflowStatusPending, workflow.Status)
	require.NotNil(t, workflow.Inputs)
	require.NotNil(t, workflow.Outputs)
	require.Greater(t, workflow.CreatedAt, int64(0))
}

func TestWorkflowLifecycle(t *testing.T) {
	workflow := types.NewWorkflow("wf-1", "Test", "owner1")

	// Start
	workflow.Start()
	require.Equal(t, types.WorkflowStatusRunning, workflow.Status)
	require.Greater(t, workflow.StartedAt, int64(0))

	// Complete
	workflow.Complete()
	require.Equal(t, types.WorkflowStatusCompleted, workflow.Status)
	require.Greater(t, workflow.CompletedAt, int64(0))

	// Fail
	workflow2 := types.NewWorkflow("wf-2", "Test", "owner1")
	workflow2.Fail("test error")
	require.Equal(t, types.WorkflowStatusFailed, workflow2.Status)
	require.Greater(t, workflow2.CompletedAt, int64(0))

	// Cancel
	workflow3 := types.NewWorkflow("wf-3", "Test", "owner1")
	workflow3.Cancel()
	require.Equal(t, types.WorkflowStatusCancelled, workflow3.Status)
	require.Greater(t, workflow3.CompletedAt, int64(0))
}

func TestWorkflowValidation(t *testing.T) {
	// Valid workflow
	validWorkflow := types.NewWorkflow("wf-1", "Test", "owner1")
	require.NoError(t, validWorkflow.Validate())

	// Invalid - no ID
	invalidWorkflow := types.NewWorkflow("", "Test", "owner1")
	require.Error(t, invalidWorkflow.Validate())

	// Invalid - no owner
	invalidWorkflow2 := types.NewWorkflow("wf-2", "Test", "")
	require.Error(t, invalidWorkflow2.Validate())
}

func TestCircularDependencyDetection(t *testing.T) {
	workflow := types.NewWorkflow("wf-1", "Test", "owner1")

	// Create steps with circular dependency: A -> B -> C -> A
	stepA := types.WorkflowStep{ID: "step-a", Name: "Step A", DependsOn: []string{"step-c"}}
	stepB := types.WorkflowStep{ID: "step-b", Name: "Step B", DependsOn: []string{"step-a"}}
	stepC := types.WorkflowStep{ID: "step-c", Name: "Step C", DependsOn: []string{"step-b"}}

	workflow.AddStep(stepA)
	workflow.AddStep(stepB)
	workflow.AddStep(stepC)

	err := workflow.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "circular dependency")
}

func TestGetReadySteps(t *testing.T) {
	workflow := types.NewWorkflow("wf-1", "Test", "owner1")

	// Add steps: A (no deps), B (depends on A), C (depends on B)
	stepA := types.WorkflowStep{ID: "step-a", Name: "Step A", Status: types.StepStatusPending}
	stepB := types.WorkflowStep{ID: "step-b", Name: "Step B", DependsOn: []string{"step-a"}, Status: types.StepStatusPending}
	stepC := types.WorkflowStep{ID: "step-c", Name: "Step C", DependsOn: []string{"step-b"}, Status: types.StepStatusPending}

	workflow.AddStep(stepA)
	workflow.AddStep(stepB)
	workflow.AddStep(stepC)

	// Initially, only step A should be ready
	ready := workflow.GetReadySteps()
	require.Len(t, ready, 1)
	require.Equal(t, "step-a", ready[0].ID)

	// Mark A as completed
	stepA.Status = types.StepStatusCompleted
	// Update step in workflow
	for i := range workflow.Steps {
		if workflow.Steps[i].ID == "step-a" {
			workflow.Steps[i] = stepA
		}
	}

	// Now B should be ready
	ready = workflow.GetReadySteps()
	require.Len(t, ready, 1)
	require.Equal(t, "step-b", ready[0].ID)
}

func TestMilestoneCreation(t *testing.T) {
	ms := types.NewMilestone("ms-1", "wf-1", "First Milestone", 1000, 1)

	require.Equal(t, "ms-1", ms.ID)
	require.Equal(t, "wf-1", ms.WorkflowID)
	require.Equal(t, "First Milestone", ms.Name)
	require.Equal(t, uint64(1000), ms.Amount)
	require.Equal(t, 1, ms.Order)
	require.Equal(t, types.MilestoneStatusPending, ms.Status)
}

func TestMilestoneLifecycle(t *testing.T) {
	ms := types.NewMilestone("ms-1", "wf-1", "Test", 1000, 1)

	// Activate
	ms.Activate()
	require.Equal(t, types.MilestoneStatusActive, ms.Status)

	// Complete
	ms.Complete()
	require.Equal(t, types.MilestoneStatusCompleted, ms.Status)
	require.Greater(t, ms.CompletedAt, int64(0))

	// Mark paid
	ms.MarkPaid()
	require.Equal(t, types.MilestoneStatusPaid, ms.Status)
	require.Greater(t, ms.PaidAt, int64(0))
}

func TestMilestonePlanValidation(t *testing.T) {
	// Valid plan
	plan := types.NewMilestonePlan("mp-1", "wf-1", 5000)
	ms1 := types.NewMilestone("ms-1", "wf-1", "First", 2000, 1)
	ms2 := types.NewMilestone("ms-2", "wf-1", "Second", 3000, 2)
	plan.AddMilestone(*ms1)
	plan.AddMilestone(*ms2)
	require.NoError(t, plan.Validate())

	// Invalid - amounts don't sum to total
	plan2 := types.NewMilestonePlan("mp-2", "wf-1", 10000)
	ms3 := types.NewMilestone("ms-3", "wf-1", "First", 2000, 1)
	plan2.AddMilestone(*ms3)
	require.Error(t, plan2.Validate())
}

func TestCapabilityExecution(t *testing.T) {
	cm := NewCapabilityManager()

	ctx := context.Background()
	config := types.DefaultCapabilityConfig(types.CapabilityResearcher)
	params := map[string]string{"topic": "blockchain technology"}

	result, err := cm.Execute(ctx, types.CapabilityResearcher, config, params)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Contains(t, result, "blockchain technology")
}

func TestWorkflowExecution(t *testing.T) {
	executor := NewExecutor()

	// Create a simple workflow
	workflow := types.NewWorkflow("wf-test", "Test Workflow", "owner1")
	workflow.Inputs["input"] = "test input"

	// Add steps
	step1 := types.WorkflowStep{
		ID:          "step-1",
		Name:        "Research",
		Type:        types.StepTypeAgent,
		Capability:  types.CapabilityResearcher,
		Status:      types.StepStatusPending,
		Params:      map[string]string{"topic": "test topic"},
	}
	step2 := types.WorkflowStep{
		ID:          "step-2",
		Name:        "Write",
		Type:        types.StepTypeAgent,
		Capability:  types.CapabilityWriter,
		DependsOn:   []string{"step-1"},
		Status:      types.StepStatusPending,
		Params:      map[string]string{"type": "report"},
	}

	workflow.AddStep(step1)
	workflow.AddStep(step2)

	// Register workflow
	err := executor.RegisterWorkflow(workflow)
	require.NoError(t, err)

	// Execute workflow
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = executor.ExecuteWorkflow(ctx, "wf-test")
	require.NoError(t, err)

	// Verify workflow completed
	completedWorkflow := executor.GetWorkflow("wf-test")
	require.Equal(t, types.WorkflowStatusCompleted, completedWorkflow.Status)
}

func TestParallelStepExecution(t *testing.T) {
	executor := NewExecutor()

	// Create workflow with parallel step
	workflow := types.NewWorkflow("wf-parallel", "Parallel Test", "owner1")

	// Create parallel step with sub-steps
	subStep1 := types.WorkflowStep{
		ID:         "sub-1",
		Name:       "Analysis",
		Type:       types.StepTypeAgent,
		Capability: types.CapabilityAnalyst,
		Status:     types.StepStatusPending,
	}
	subStep2 := types.WorkflowStep{
		ID:         "sub-2",
		Name:       "Review",
		Type:       types.StepTypeAgent,
		Capability: types.CapabilityReviewer,
		Status:     types.StepStatusPending,
	}

	parallelStep := types.WorkflowStep{
		ID:        "parallel-1",
		Name:      "Parallel Execution",
		Type:      types.StepTypeParallel,
		Status:    types.StepStatusPending,
		SubSteps:  []types.WorkflowStep{subStep1, subStep2},
	}

	workflow.AddStep(parallelStep)

	// Register and execute
	err := executor.RegisterWorkflow(workflow)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = executor.ExecuteWorkflow(ctx, "wf-parallel")
	require.NoError(t, err)

	// Verify
	completedWorkflow := executor.GetWorkflow("wf-parallel")
	require.Equal(t, types.WorkflowStatusCompleted, completedWorkflow.Status)
}

func TestProgressReport(t *testing.T) {
	executor := NewExecutor()

	// Create workflow with milestones
	workflow := types.NewWorkflow("wf-progress", "Progress Test", "owner1")

	step1 := types.WorkflowStep{
		ID:   "step-1",
		Name: "Step 1",
		Type: types.StepTypeAgent,
		Capability: types.CapabilityResearcher,
		Status: types.StepStatusCompleted,
	}
	step2 := types.WorkflowStep{
		ID:   "step-2",
		Name: "Step 2",
		Type: types.StepTypeAgent,
		Capability: types.CapabilityWriter,
		Status: types.StepStatusRunning,
	}

	workflow.AddStep(step1)
	workflow.AddStep(step2)

	// Create milestone plan
	plan := types.NewMilestonePlan("mp-1", "wf-progress", 3000)
	ms1 := types.NewMilestone("ms-1", "wf-progress", "First", 1000, 1)
	ms1.Complete()
	ms1.MarkPaid()
	ms2 := types.NewMilestone("ms-2", "wf-progress", "Second", 2000, 2)
	plan.AddMilestone(*ms1)
	plan.AddMilestone(*ms2)

	// Register
	err := executor.RegisterWorkflow(workflow)
	require.NoError(t, err)
	err = executor.RegisterMilestonePlan(plan)
	require.NoError(t, err)

	// Generate report
	report, err := executor.GenerateProgressReport("wf-progress")
	require.NoError(t, err)

	require.Equal(t, "wf-progress", report.WorkflowID)
	require.Equal(t, 2, report.TotalSteps)
	require.Equal(t, 1, report.CompletedSteps)
	require.Equal(t, 2, report.TotalMilestones)
	require.Equal(t, 1, report.CompletedMs)
	require.Equal(t, 1, report.PaidMs)
	require.Greater(t, report.ProgressPct, float64(0))
}

func TestWorkflowEvents(t *testing.T) {
	executor := NewExecutor()

	// Subscribe to events
	eventCh := executor.GetEventChannel()

	// Create and execute workflow
	workflow := types.NewWorkflow("wf-events", "Event Test", "owner1")
	step := types.WorkflowStep{
		ID:         "step-1",
		Name:       "Test Step",
		Type:       types.StepTypeAgent,
		Capability: types.CapabilityResearcher,
		Status:     types.StepStatusPending,
	}
	workflow.AddStep(step)

	err := executor.RegisterWorkflow(workflow)
	require.NoError(t, err)

	// Execute in background
	go func() {
		ctx := context.Background()
		err := executor.ExecuteWorkflow(ctx, "wf-events")
		_ = err // ignore error in goroutine
	}()

	// Collect events
	var events []WorkflowEvent
	timeout := time.After(5 * time.Second)

	for len(events) < 3 { // workflow_started, step_started, step_completed
		select {
		case event := <-eventCh:
			events = append(events, event)
		case <-timeout:
			return
		}
	}

	// Verify we got events
	require.GreaterOrEqual(t, len(events), 1)

	// First event should be workflow started
	require.Equal(t, "wf-events", events[0].WorkflowID)
	require.Equal(t, EventWorkflowStarted, events[0].Type)
}
