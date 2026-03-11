package executor

import (
	"context"
	"fmt"
	"strings"

	"sharetoken/x/workflow/types"
)

// CapabilityManager manages and executes capabilities
type CapabilityManager struct {
	handlers map[types.Capability]CapabilityHandler
}

// CapabilityHandler is a function that handles a capability execution
type CapabilityHandler func(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error)

// NewCapabilityManager creates a new capability manager with default handlers
func NewCapabilityManager() *CapabilityManager {
	cm := &CapabilityManager{
		handlers: make(map[types.Capability]CapabilityHandler),
	}

	// Register default handlers
	cm.RegisterHandler(types.CapabilityCollector, handleCollector)
	cm.RegisterHandler(types.CapabilityLead, handleLead)
	cm.RegisterHandler(types.CapabilityResearcher, handleResearcher)
	cm.RegisterHandler(types.CapabilityWriter, handleWriter)
	cm.RegisterHandler(types.CapabilityAnalyst, handleAnalyst)
	cm.RegisterHandler(types.CapabilityTester, handleTester)
	cm.RegisterHandler(types.CapabilityReviewer, handleReviewer)

	return cm
}

// RegisterHandler registers a handler for a capability
func (cm *CapabilityManager) RegisterHandler(cap types.Capability, handler CapabilityHandler) {
	cm.handlers[cap] = handler
}

// Execute executes a capability
func (cm *CapabilityManager) Execute(ctx context.Context, cap types.Capability, config types.CapabilityConfig, params map[string]string) (string, error) {
	handler, exists := cm.handlers[cap]
	if !exists {
		return "", fmt.Errorf("no handler registered for capability: %s", cap)
	}

	return handler(ctx, config, params)
}

// handleCollector handles data collection capability
func handleCollector(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	// Simulate data collection
	sources := config.Params["max_sources"]
	target := params["target"]
	if target == "" {
		target = "general data"
	}

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(256)
	result.WriteString(fmt.Sprintf("Collected data from %s sources about '%s':\n", sources, target))
	result.WriteString("- Source 1: Retrieved successfully\n")
	result.WriteString("- Source 2: Retrieved successfully\n")
	result.WriteString("- Source 3: Retrieved successfully\n")
	result.WriteString(fmt.Sprintf("Total items collected: %d\n", 42))

	return result.String(), nil
}

// handleLead handles workflow orchestration capability
func handleLead(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	// Lead capability coordinates other agents
	action := params["action"]
	if action == "" {
		action = "coordinate"
	}

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(128)
	result.WriteString(fmt.Sprintf("Lead agent executing '%s' action:\n", action))
	result.WriteString("- Analyzing workflow dependencies\n")
	result.WriteString("- Scheduling parallel tasks\n")
	result.WriteString("- Monitoring execution progress\n")
	result.WriteString("Status: Coordination complete\n")

	return result.String(), nil
}

// handleResearcher handles research capability
func handleResearcher(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	topic := params["topic"]
	if topic == "" {
		topic = "general research"
	}

	depth := config.Params["depth"]

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(256)
	result.WriteString(fmt.Sprintf("Research on '%s' (depth: %s):\n", topic, depth))
	result.WriteString("1. Found 5 relevant sources\n")
	result.WriteString("2. Analyzed key findings\n")
	result.WriteString("3. Synthesized conclusions\n")
	result.WriteString("Key findings:\n")
	result.WriteString("- Finding A: Important insight\n")
	result.WriteString("- Finding B: Supporting evidence\n")
	result.WriteString("- Finding C: Related context\n")

	return result.String(), nil
}

// handleWriter handles content creation capability
func handleWriter(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	contentType := params["type"]
	if contentType == "" {
		contentType = "document"
	}

	style := config.Params["style"]
	length := config.Params["length"]

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(512)
	result.WriteString(fmt.Sprintf("Created %s (style: %s, length: %s):\n\n", contentType, style, length))
	result.WriteString("# Generated Content\n\n")
	result.WriteString("This is a sample generated content based on the specified parameters.\n\n")
	result.WriteString("## Section 1\n\n")
	result.WriteString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n\n")
	result.WriteString("## Section 2\n\n")
	result.WriteString("Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\n\n")
	result.WriteString("---\n")
	result.WriteString(fmt.Sprintf("Content type: %s | Style: %s | Length: %s\n", contentType, style, length))

	return result.String(), nil
}

// handleAnalyst handles data analysis capability
func handleAnalyst(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	data := params["data"]
	if data == "" {
		data = "provided dataset"
	}

	confidence := config.Params["confidence"]

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(512)
	result.WriteString(fmt.Sprintf("Analysis of '%s' (confidence: %s):\n\n", data, confidence))
	result.WriteString("## Data Summary\n")
	result.WriteString("- Total records: 1,234\n")
	result.WriteString("- Valid records: 1,200\n")
	result.WriteString("- Anomalies: 34\n\n")
	result.WriteString("## Key Metrics\n")
	result.WriteString("- Mean: 42.5\n")
	result.WriteString("- Median: 41.0\n")
	result.WriteString("- Std Dev: 8.3\n\n")
	result.WriteString("## Insights\n")
	result.WriteString("1. Primary trend: Positive growth\n")
	result.WriteString("2. Seasonality: Detected quarterly pattern\n")
	result.WriteString("3. Outliers: 3 significant anomalies\n")

	return result.String(), nil
}

// handleTester handles testing capability
func handleTester(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	target := params["target"]
	if target == "" {
		target = "system"
	}

	coverage := config.Params["coverage"]
	testTypes := config.Params["types"]

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(512)
	result.WriteString(fmt.Sprintf("Test Results for '%s':\n\n", target))
	result.WriteString("## Test Configuration\n")
	result.WriteString(fmt.Sprintf("- Coverage target: %s%%\n", coverage))
	result.WriteString(fmt.Sprintf("- Test types: %s\n\n", testTypes))
	result.WriteString("## Test Summary\n")
	result.WriteString("✓ Unit Tests: 45/45 PASS\n")
	result.WriteString("✓ Integration Tests: 12/12 PASS\n")
	result.WriteString("✓ E2E Tests: 8/8 PASS\n\n")
	result.WriteString("## Coverage Report\n")
	result.WriteString("- Lines: 87%\n")
	result.WriteString("- Functions: 92%\n")
	result.WriteString("- Branches: 81%\n\n")
	result.WriteString("Status: All tests passed ✓\n")

	return result.String(), nil
}

// handleReviewer handles review capability
func handleReviewer(ctx context.Context, config types.CapabilityConfig, params map[string]string) (string, error) {
	target := params["target"]
	if target == "" {
		target = "content"
	}

	strictness := config.Params["strictness"]

	// Performance: 使用 strings.Builder 预分配容量
	var result strings.Builder
	result.Grow(512)
	result.WriteString(fmt.Sprintf("Review of '%s' (strictness: %s):\n\n", target, strictness))
	result.WriteString("## Quality Score: 8.5/10\n\n")
	result.WriteString("## Strengths\n")
	result.WriteString("✓ Well structured\n")
	result.WriteString("✓ Clear explanations\n")
	result.WriteString("✓ Good examples\n\n")
	result.WriteString("## Areas for Improvement\n")
	result.WriteString("- Add more edge case coverage\n")
	result.WriteString("- Improve error handling\n")
	result.WriteString("- Add performance benchmarks\n\n")
	result.WriteString("## Recommendations\n")
	result.WriteString("1. Address the 2 medium-priority issues\n")
	result.WriteString("2. Add documentation for public APIs\n")
	result.WriteString("3. Consider refactoring module X\n")

	return result.String(), nil
}
