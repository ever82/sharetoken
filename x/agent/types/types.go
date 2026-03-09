package types

const (
	ModuleName = "agent"
	StoreKey   = ModuleName
)

// AgentTemplate represents agent template type
type AgentTemplate string

const (
	TemplateCoder      AgentTemplate = "coder"
	TemplateResearcher AgentTemplate = "researcher"
	TemplateWriter     AgentTemplate = "writer"
	TemplateAnalyst    AgentTemplate = "analyst"
	TemplateDesigner   AgentTemplate = "designer"
	TemplateTester     AgentTemplate = "tester"
	TemplateReviewer   AgentTemplate = "reviewer"
	TemplateArchitect  AgentTemplate = "architect"
)

// GetAllTemplates returns all available agent templates
func GetAllTemplates() []AgentTemplate {
	return []AgentTemplate{
		TemplateCoder,
		TemplateResearcher,
		TemplateWriter,
		TemplateAnalyst,
		TemplateDesigner,
		TemplateTester,
		TemplateReviewer,
		TemplateArchitect,
	}
}

// TaskStatus represents task execution status
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)
