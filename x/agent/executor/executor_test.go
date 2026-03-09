package executor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/agent/security"
	"sharetoken/x/agent/types"
)

func TestSecurityLayers(t *testing.T) {
	layers := security.GetSecurityLayers()
	require.Len(t, layers, 16)

	// Check key layers exist
	layerNames := []string{}
	for _, l := range layers {
		layerNames = append(layerNames, l.Name)
	}

	require.Contains(t, layerNames, "WASM Sandbox")
	require.Contains(t, layerNames, "Memory Limit")
	require.Contains(t, layerNames, "Network Isolation")
	require.Contains(t, layerNames, "Emergency Kill")
}

func TestSecurityConfig(t *testing.T) {
	// Test high security
	highSec := security.NewSecurityConfig(security.LevelHigh)
	require.Equal(t, security.LevelHigh, highSec.Level)
	require.Len(t, highSec.Layers, 16)

	enabledHigh := highSec.GetEnabledLayers()
	require.GreaterOrEqual(t, len(enabledHigh), 9)

	// Test minimal security
	minSec := security.NewSecurityConfig(security.LevelMinimal)
	require.Equal(t, security.LevelMinimal, minSec.Level)

	enabledMin := minSec.GetEnabledLayers()
	require.LessOrEqual(t, len(enabledMin), 13)

	// Validate
	require.NoError(t, highSec.Validate())
	require.NoError(t, minSec.Validate())
}

func TestAgentCreation(t *testing.T) {
	agent := types.NewAgent("agent-1", "owner1", types.TemplateCoder)

	require.Equal(t, "agent-1", agent.ID)
	require.Equal(t, "owner1", agent.Owner)
	require.Equal(t, types.TemplateCoder, agent.Template)
	require.Equal(t, "wasm", agent.Runtime)
	require.True(t, agent.Active)
}

func TestAgentValidation(t *testing.T) {
	validAgent := types.NewAgent("agent-1", "owner1", types.TemplateCoder)
	require.NoError(t, validAgent.Validate())

	invalidAgent := types.Agent{ID: "", Owner: "owner1"}
	require.Error(t, invalidAgent.Validate())
}

func TestTaskExecution(t *testing.T) {
	executor := NewExecutor(security.LevelStandard)

	tests := []struct {
		name     string
		template types.AgentTemplate
		input    string
	}{
		{"Coder", types.TemplateCoder, "sort array"},
		{"Researcher", types.TemplateResearcher, "blockchain history"},
		{"Writer", types.TemplateWriter, "blog post"},
		{"Analyst", types.TemplateAnalyst, "market data"},
		{"Tester", types.TemplateTester, "login flow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := types.NewAgent("agent-1", "owner1", tt.template)
			task := types.NewTask("task-1", agent.ID, "service-1", tt.input)

			result, err := executor.ExecuteTask(agent, task)
			require.NoError(t, err)
			require.True(t, result.Success)
			require.NotEmpty(t, result.Output)
			require.NotEmpty(t, result.Layers)
			require.Greater(t, result.GasUsed, uint64(0))
			require.Equal(t, types.TaskStatusCompleted, task.Status)
		})
	}
}

func TestTaskLifecycle(t *testing.T) {
	task := types.NewTask("task-1", "agent-1", "service-1", "test input")

	require.Equal(t, types.TaskStatusPending, task.Status)

	task.Start()
	require.Equal(t, types.TaskStatusRunning, task.Status)
	require.Greater(t, task.StartedAt, int64(0))

	task.Complete("output")
	require.Equal(t, types.TaskStatusCompleted, task.Status)
	require.Equal(t, "output", task.Output)
	require.Greater(t, task.CompletedAt, int64(0))

	// Test fail
	task2 := types.NewTask("task-2", "agent-1", "service-1", "test")
	task2.Fail("error message")
	require.Equal(t, types.TaskStatusFailed, task2.Status)
	require.Equal(t, "error message", task2.Error)
}

func TestAgentTemplates(t *testing.T) {
	templates := types.GetAllTemplates()
	require.GreaterOrEqual(t, len(templates), 8)

	// Check that all expected templates exist
	templateMap := make(map[string]bool)
	for _, t := range templates {
		templateMap[string(t)] = true
	}

	require.True(t, templateMap["coder"])
	require.True(t, templateMap["researcher"])
	require.True(t, templateMap["writer"])
	require.True(t, templateMap["analyst"])
}

func TestAgentConfig(t *testing.T) {
	config := types.DefaultAgentConfig()

	require.Greater(t, config.MaxTokens, int64(0))
	require.Greater(t, config.Timeout, int64(0))
	require.Greater(t, config.MaxMemoryMB, int64(0))
	require.NotEmpty(t, config.AllowedTools)
	require.NotEmpty(t, config.RestrictedDirs)
}
