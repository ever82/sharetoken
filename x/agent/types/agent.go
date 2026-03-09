package types

import (
	"fmt"
	"time"
)

// Agent represents an AI agent instance
type Agent struct {
	ID       string        `json:"id"`
	Owner    string        `json:"owner"`
	Template AgentTemplate `json:"template"`
	Runtime  string        `json:"runtime"` // "wasm", "rust", "python"
	Config   AgentConfig   `json:"config"`
	Active   bool          `json:"active"`
	CreatedAt int64        `json:"created_at"`
}

// AgentConfig represents agent configuration
type AgentConfig struct {
	MaxTokens      int64    `json:"max_tokens"`
	Temperature    float64  `json:"temperature"`
	Timeout        int64    `json:"timeout"`        // seconds
	MaxMemoryMB    int64    `json:"max_memory_mb"`
	AllowedTools   []string `json:"allowed_tools"`
	RestrictedDirs []string `json:"restricted_dirs"`
}

// DefaultAgentConfig returns default configuration
func DefaultAgentConfig() AgentConfig {
	return AgentConfig{
		MaxTokens:      4096,
		Temperature:    0.7,
		Timeout:        300, // 5 minutes
		MaxMemoryMB:    512,
		AllowedTools:   []string{"file", "search", "calc"},
		RestrictedDirs: []string{"/etc", "/usr", "/bin"},
	}
}

// NewAgent creates a new agent
func NewAgent(id, owner string, template AgentTemplate) *Agent {
	return &Agent{
		ID:        id,
		Owner:     owner,
		Template:  template,
		Runtime:   "wasm",
		Config:    DefaultAgentConfig(),
		Active:    true,
		CreatedAt: time.Now().Unix(),
	}
}

// Validate validates agent configuration
func (a Agent) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("agent ID cannot be empty")
	}
	if a.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if a.Runtime == "" {
		return fmt.Errorf("runtime cannot be empty")
	}
	return nil
}

// String implements stringer
func (a Agent) String() string {
	return fmt.Sprintf("Agent{%s: %s (%s), runtime: %s, active: %v}",
		a.ID, a.Template, a.Owner, a.Runtime, a.Active)
}

// Task represents an agent task
type Task struct {
	ID          string     `json:"id"`
	AgentID     string     `json:"agent_id"`
	ServiceID   string     `json:"service_id"`
	Input       string     `json:"input"`
	Output      string     `json:"output"`
	Status      TaskStatus `json:"status"`
	Error       string     `json:"error"`
	CreatedAt   int64      `json:"created_at"`
	StartedAt   int64      `json:"started_at"`
	CompletedAt int64      `json:"completed_at"`
	GasUsed     uint64     `json:"gas_used"`
}

// NewTask creates a new task
func NewTask(id, agentID, serviceID, input string) *Task {
	return &Task{
		ID:        id,
		AgentID:   agentID,
		ServiceID: serviceID,
		Input:     input,
		Status:    TaskStatusPending,
		CreatedAt: time.Now().Unix(),
	}
}

// Start marks task as started
func (t *Task) Start() {
	t.Status = TaskStatusRunning
	t.StartedAt = time.Now().Unix()
}

// Complete marks task as completed
func (t *Task) Complete(output string) {
	t.Status = TaskStatusCompleted
	t.Output = output
	t.CompletedAt = time.Now().Unix()
}

// Fail marks task as failed
func (t *Task) Fail(err string) {
	t.Status = TaskStatusFailed
	t.Error = err
	t.CompletedAt = time.Now().Unix()
}
