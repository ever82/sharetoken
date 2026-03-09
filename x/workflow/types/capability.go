package types

// Capability represents an autonomous capability package
type Capability string

const (
	CapabilityCollector Capability = "collector" // Data collection and aggregation
	CapabilityLead      Capability = "lead"      // Workflow orchestration and coordination
	CapabilityResearcher Capability = "researcher" // Research and information gathering
	CapabilityWriter    Capability = "writer"    // Content creation and documentation
	CapabilityAnalyst   Capability = "analyst"   // Data analysis and insights
	CapabilityTester    Capability = "tester"    // Testing and validation
	CapabilityReviewer  Capability = "reviewer"  // Code/content review and feedback
)

// GetAllCapabilities returns all available capabilities
func GetAllCapabilities() []Capability {
	return []Capability{
		CapabilityCollector,
		CapabilityLead,
		CapabilityResearcher,
		CapabilityWriter,
		CapabilityAnalyst,
		CapabilityTester,
		CapabilityReviewer,
	}
}

// CapabilityConfig represents configuration for a capability
type CapabilityConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Params      map[string]string `json:"params"`
	Timeout     int64             `json:"timeout"` // seconds
}

// DefaultCapabilityConfig returns default config for a capability
func DefaultCapabilityConfig(cap Capability) CapabilityConfig {
	configs := map[Capability]CapabilityConfig{
		CapabilityCollector: {
			Name:        "collector",
			Description: "Collect and aggregate data from multiple sources",
			Params:      map[string]string{"max_sources": "10", "timeout": "60"},
			Timeout:     300,
		},
		CapabilityLead: {
			Name:        "lead",
			Description: "Orchestrate workflow execution and coordinate agents",
			Params:      map[string]string{"max_parallel": "5", "retry_count": "3"},
			Timeout:     600,
		},
		CapabilityResearcher: {
			Name:        "researcher",
			Description: "Research topics and gather information",
			Params:      map[string]string{"depth": "medium", "sources": "web"},
			Timeout:     300,
		},
		CapabilityWriter: {
			Name:        "writer",
			Description: "Create content and documentation",
			Params:      map[string]string{"style": "professional", "length": "medium"},
			Timeout:     300,
		},
		CapabilityAnalyst: {
			Name:        "analyst",
			Description: "Analyze data and generate insights",
			Params:      map[string]string{"metrics": "all", "confidence": "0.95"},
			Timeout:     300,
		},
		CapabilityTester: {
			Name:        "tester",
			Description: "Test implementations and validate results",
			Params:      map[string]string{"coverage": "80", "types": "unit,integration"},
			Timeout:     600,
		},
		CapabilityReviewer: {
			Name:        "reviewer",
			Description: "Review code and content for quality",
			Params:      map[string]string{"strictness": "medium", "focus": "all"},
			Timeout:     300,
		},
	}
	return configs[cap]
}
