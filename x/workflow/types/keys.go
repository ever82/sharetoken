package types

const (
	// ModuleName is the name of the workflow module
	ModuleName = "workflow"

	// StoreKey is the string store key for the workflow module
	StoreKey = ModuleName

	// RouterKey is the message route for the workflow module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the workflow module
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// WorkflowKey is the prefix for workflow store
	WorkflowKey = []byte{0x01}

	// WorkflowStepKey is the prefix for workflow step store
	WorkflowStepKey = []byte{0x02}

	// CapabilityKey is the prefix for capability store
	CapabilityKey = []byte{0x03}
)

// GetWorkflowKey returns the key for a workflow by ID
func GetWorkflowKey(id string) []byte {
	return append(WorkflowKey, []byte(id)...)
}

// GetWorkflowStepKey returns the key for a workflow step
func GetWorkflowStepKey(workflowID, stepID string) []byte {
	return append(append(WorkflowStepKey, []byte(workflowID)...), []byte(stepID)...)
}

// GetCapabilityKey returns the key for a capability by ID
func GetCapabilityKey(id string) []byte {
	return append(CapabilityKey, []byte(id)...)
}
