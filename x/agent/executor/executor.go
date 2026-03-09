package executor

import (
	"fmt"
	"time"

	"sharetoken/x/agent/security"
	"sharetoken/x/agent/types"
)

// Executor represents the agent executor
type Executor struct {
	securityConfig *security.SecurityConfig
	timeout        time.Duration
}

// NewExecutor creates a new agent executor
func NewExecutor(secLevel security.SecurityLevel) *Executor {
	return &Executor{
		securityConfig: security.NewSecurityConfig(secLevel),
		timeout:        5 * time.Minute,
	}
}

// ExecutionResult represents task execution result
type ExecutionResult struct {
	Success   bool
	Output    string
	Error     string
	GasUsed   uint64
	Duration  time.Duration
	Layers    []string // Enabled security layers
}

// ExecuteTask executes an agent task
func (e *Executor) ExecuteTask(agent *types.Agent, task *types.Task) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Success: false,
		Layers:  []string{},
	}

	startTime := time.Now()

	// Validate security config
	if err := e.securityConfig.Validate(); err != nil {
		result.Error = fmt.Sprintf("security validation failed: %v", err)
		return result, err
	}

	// Record enabled security layers
	for _, layer := range e.securityConfig.GetEnabledLayers() {
		result.Layers = append(result.Layers, layer.Name)
	}

	// Start task	task.Start()

	// Simulate execution based on agent template
	output, err := e.simulateExecution(agent, task)
	if err != nil {
		task.Fail(err.Error())
		result.Error = err.Error()
		return result, err
	}

	// Complete task
	task.Complete(output)
	result.Success = true
	result.Output = output
	result.GasUsed = 1000 // Simulated gas
	result.Duration = time.Since(startTime)

	return result, nil
}

// simulateExecution simulates agent execution based on template
func (e *Executor) simulateExecution(agent *types.Agent, task *types.Task) (string, error) {
	// Simulate different behaviors based on template
	switch agent.Template {
	case types.TemplateCoder:
		return fmt.Sprintf("// Generated code for: %s\nfunction solution() {\n  return 'implemented';\n}", task.Input), nil

	case types.TemplateResearcher:
		return fmt.Sprintf("Research findings for '%s':\n1. Found relevant sources\n2. Analyzed data\n3. Conclusion generated", task.Input), nil

	case types.TemplateWriter:
		return fmt.Sprintf("Content for '%s':\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit...", task.Input), nil

	case types.TemplateAnalyst:
		return fmt.Sprintf("Analysis of '%s':\n- Data points: 100\n- Trend: positive\n- Confidence: 95%%", task.Input), nil

	case types.TemplateTester:
		return fmt.Sprintf("Test results for '%s':\n✓ Test 1: PASS\n✓ Test 2: PASS\n✓ Test 3: PASS\nCoverage: 95%%", task.Input), nil

	default:
		return fmt.Sprintf("Processed '%s' using %s template", task.Input, agent.Template), nil
	}
}

// GetSecurityConfig returns the executor's security configuration
func (e *Executor) GetSecurityConfig() *security.SecurityConfig {
	return e.securityConfig
}
