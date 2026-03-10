package types

import (
	"fmt"
	"strings"
)

// NodeRole represents the role of a ShareToken node
type NodeRole int32

const (
	// RoleUndefined is an invalid/undefined role
	RoleUndefined NodeRole = 0
	// RoleLight is a light node - minimal storage, runs core + GenieBot
	RoleLight NodeRole = 1
	// RoleFull is a full node - complete state and current history
	RoleFull NodeRole = 2
	// RoleService is a service node - core + plugins for service provision
	RoleService NodeRole = 3
	// RoleArchive is an archive node - full history with indexes
	RoleArchive NodeRole = 4
)

// String returns the string representation of the role
func (r NodeRole) String() string {
	switch r {
	case RoleLight:
		return "light"
	case RoleFull:
		return "full"
	case RoleService:
		return "service"
	case RoleArchive:
		return "archive"
	default:
		return "undefined"
	}
}

// FromString parses a role from string
func NodeRoleFromString(s string) (NodeRole, error) {
	switch strings.ToLower(s) {
	case "light":
		return RoleLight, nil
	case "full":
		return RoleFull, nil
	case "service":
		return RoleService, nil
	case "archive":
		return RoleArchive, nil
	default:
		return RoleUndefined, fmt.Errorf("unknown node role: %s", s)
	}
}

// IsValid checks if the role is valid
func (r NodeRole) IsValid() bool {
	return r >= RoleLight && r <= RoleArchive
}

// CanSwitchTo checks if current role can switch to target role
// Some role transitions require restart, others can be done hot
func (r NodeRole) CanSwitchTo(target NodeRole) (canHotSwitch bool, requiresRestart bool) {
	// Same role - no switch needed
	if r == target {
		return true, false
	}

	// Role upgrade/downgrade rules:
	// Light <-> Full: Requires restart (different pruning strategy)
	// Light <-> Service: Can hot switch (just load/unload plugins)
	// Full <-> Archive: Requires restart (different storage backend)
	// Service <-> Full: Requires restart
	// Archive <-> any: Requires restart

	switch r {
	case RoleLight:
		switch target {
		case RoleService:
			return true, false // Can hot switch to service
		default:
			return false, true // Others require restart
		}
	case RoleService:
		switch target {
		case RoleLight:
			return true, false // Can hot switch to light
		default:
			return false, true // Others require restart
		}
	case RoleFull, RoleArchive:
		// Full and Archive require restart for any change
		return false, true
	default:
		return false, true
	}
}

// RoleConfig contains configuration specific to a node role
type RoleConfig struct {
	Role NodeRole `json:"role" yaml:"role"`

	// Light node specific
	LightConfig *LightNodeConfig `json:"light_config,omitempty" yaml:"light_config,omitempty"`

	// Full node specific
	FullConfig *FullNodeConfig `json:"full_config,omitempty" yaml:"full_config,omitempty"`

	// Service node specific
	ServiceConfig *ServiceNodeConfig `json:"service_config,omitempty" yaml:"service_config,omitempty"`

	// Archive node specific
	ArchiveConfig *ArchiveNodeConfig `json:"archive_config,omitempty" yaml:"archive_config,omitempty"`
}

// LightNodeConfig configuration for light nodes
type LightNodeConfig struct {
	// MaxPeers maximum number of P2P peers
	MaxPeers int `json:"max_peers" yaml:"max_peers"`
	// PruningStrategy for state storage ("none", "everything", "custom")
	PruningStrategy string `json:"pruning_strategy" yaml:"pruning_strategy"`
	// EnableGenieBot enable GenieBot plugin
	EnableGenieBot bool `json:"enable_genie_bot" yaml:"enable_genie_bot"`
	// TrustedFullNodes list of trusted full nodes for queries
	TrustedFullNodes []string `json:"trusted_full_nodes" yaml:"trusted_full_nodes"`
}

// FullNodeConfig configuration for full nodes
type FullNodeConfig struct {
	// MaxPeers maximum number of P2P peers
	MaxPeers int `json:"max_peers" yaml:"max_peers"`
	// PruningStrategy for state storage
	PruningStrategy string `json:"pruning_strategy" yaml:"pruning_strategy"`
	// PruningKeepRecent number of recent blocks to keep
	PruningKeepRecent uint64 `json:"pruning_keep_recent" yaml:"pruning_keep_recent"`
	// PruningKeepEvery interval to keep blocks
	PruningKeepEvery uint64 `json:"pruning_keep_every" yaml:"pruning_keep_every"`
	// EnableIndexer enable transaction indexer
	EnableIndexer bool `json:"enable_indexer" yaml:"enable_indexer"`
}

// ServiceNodeConfig configuration for service nodes
type ServiceNodeConfig struct {
	// MaxPeers maximum number of P2P peers
	MaxPeers int `json:"max_peers" yaml:"max_peers"`
	// EnableLLMPlugin enable LLM API Key custody
	EnableLLMPlugin bool `json:"enable_llm_plugin" yaml:"enable_llm_plugin"`
	// EnableAgentPlugin enable Agent executor
	EnableAgentPlugin bool `json:"enable_agent_plugin" yaml:"enable_agent_plugin"`
	// EnableWorkflowPlugin enable Workflow executor
	EnableWorkflowPlugin bool `json:"enable_workflow_plugin" yaml:"enable_workflow_plugin"`
	// WASMMaxMemory maximum memory for WASM sandboxes (MB)
	WASMMaxMemory int `json:"wasm_max_memory" yaml:"wasm_max_memory"`
	// WASMMaxCPUTime maximum CPU time for WASM execution (ms)
	WASMMaxCPUTime int `json:"wasm_max_cpu_time" yaml:"wasm_max_cpu_time"`
	// TrustedFullNodes list of trusted full nodes for state queries
	TrustedFullNodes []string `json:"trusted_full_nodes" yaml:"trusted_full_nodes"`
}

// ArchiveNodeConfig configuration for archive nodes
type ArchiveNodeConfig struct {
	// MaxPeers maximum number of P2P peers
	MaxPeers int `json:"max_peers" yaml:"max_peers"`
	// EnableBlockExplorer enable block explorer API
	EnableBlockExplorer bool `json:"enable_block_explorer" yaml:"enable_block_explorer"`
	// BlockExplorerPort port for block explorer API
	BlockExplorerPort int `json:"block_explorer_port" yaml:"block_explorer_port"`
	// IndexAllTransactions index all transactions
	IndexAllTransactions bool `json:"index_all_transactions" yaml:"index_all_transactions"`
	// IndexEvents index all events
	IndexEvents bool `json:"index_events" yaml:"index_events"`
	// CompressionEnabled enable block compression
	CompressionEnabled bool `json:"compression_enabled" yaml:"compression_enabled"`
	// ArchiveRetentionYears years of data to retain (0 = forever)
	ArchiveRetentionYears int `json:"archive_retention_years" yaml:"archive_retention_years"`
}

// DefaultLightConfig returns default light node config
func DefaultLightConfig() *LightNodeConfig {
	return &LightNodeConfig{
		MaxPeers:         10,
		PruningStrategy:  "everything",
		EnableGenieBot:   true,
		TrustedFullNodes: []string{},
	}
}

// DefaultFullConfig returns default full node config
func DefaultFullConfig() *FullNodeConfig {
	return &FullNodeConfig{
		MaxPeers:          40,
		PruningStrategy:   "custom",
		PruningKeepRecent: 100,
		PruningKeepEvery:  10000,
		EnableIndexer:     true,
	}
}

// DefaultServiceConfig returns default service node config
func DefaultServiceConfig() *ServiceNodeConfig {
	return &ServiceNodeConfig{
		MaxPeers:             20,
		EnableLLMPlugin:      true,
		EnableAgentPlugin:    true,
		EnableWorkflowPlugin: true,
		WASMMaxMemory:        512,
		WASMMaxCPUTime:       30000,
		TrustedFullNodes:     []string{},
	}
}

// DefaultArchiveConfig returns default archive node config
func DefaultArchiveConfig() *ArchiveNodeConfig {
	return &ArchiveNodeConfig{
		MaxPeers:              50,
		EnableBlockExplorer:   true,
		BlockExplorerPort:     8080,
		IndexAllTransactions:  true,
		IndexEvents:           true,
		CompressionEnabled:    true,
		ArchiveRetentionYears: 0, // Forever
	}
}

// DefaultRoleConfig returns default config for a role
func DefaultRoleConfig(role NodeRole) *RoleConfig {
	config := &RoleConfig{
		Role: role,
	}

	switch role {
	case RoleLight:
		config.LightConfig = DefaultLightConfig()
	case RoleFull:
		config.FullConfig = DefaultFullConfig()
	case RoleService:
		config.ServiceConfig = DefaultServiceConfig()
	case RoleArchive:
		config.ArchiveConfig = DefaultArchiveConfig()
	}

	return config
}

// Validate checks if the config is valid for the role
func (rc *RoleConfig) Validate() error {
	if !rc.Role.IsValid() {
		return fmt.Errorf("invalid node role: %s", rc.Role.String())
	}

	switch rc.Role {
	case RoleLight:
		if rc.LightConfig == nil {
			return fmt.Errorf("light_config is required for light role")
		}
		if rc.LightConfig.MaxPeers < 1 {
			return fmt.Errorf("max_peers must be >= 1")
		}
	case RoleFull:
		if rc.FullConfig == nil {
			return fmt.Errorf("full_config is required for full role")
		}
		if rc.FullConfig.MaxPeers < 1 {
			return fmt.Errorf("max_peers must be >= 1")
		}
	case RoleService:
		if rc.ServiceConfig == nil {
			return fmt.Errorf("service_config is required for service role")
		}
		if rc.ServiceConfig.MaxPeers < 1 {
			return fmt.Errorf("max_peers must be >= 1")
		}
		if rc.ServiceConfig.WASMMaxMemory < 64 {
			return fmt.Errorf("wasm_max_memory must be >= 64 MB")
		}
	case RoleArchive:
		if rc.ArchiveConfig == nil {
			return fmt.Errorf("archive_config is required for archive role")
		}
		if rc.ArchiveConfig.MaxPeers < 1 {
			return fmt.Errorf("max_peers must be >= 1")
		}
	}

	return nil
}

// NodeCapabilities represents what a node can do
type NodeCapabilities struct {
	// CanValidate can participate in consensus
	CanValidate bool
	// CanQueryState can query full state
	CanQueryState bool
	// CanQueryHistory can query historical data
	CanQueryHistory bool
	// CanServeLightClients can serve light client proofs
	CanServeLightClients bool
	// CanRunPlugins can run service plugins
	CanRunPlugins bool
	// CanIndexBlocks can index blocks for explorer
	CanIndexBlocks bool
	// StorageRequirementGB estimated storage requirement in GB
	StorageRequirementGB int
	// MemoryRequirementGB estimated memory requirement in GB
	MemoryRequirementGB int
}

// GetCapabilities returns capabilities for a role
func (r NodeRole) GetCapabilities() NodeCapabilities {
	switch r {
	case RoleLight:
		return NodeCapabilities{
			CanValidate:          false,
			CanQueryState:        false, // Relies on full nodes
			CanQueryHistory:      false,
			CanServeLightClients: false,
			CanRunPlugins:        true, // GenieBot only
			CanIndexBlocks:       false,
			StorageRequirementGB: 10,
			MemoryRequirementGB:  2,
		}
	case RoleFull:
		return NodeCapabilities{
			CanValidate:          true,
			CanQueryState:        true,
			CanQueryHistory:      false, // Only recent history
			CanServeLightClients: true,
			CanRunPlugins:        false,
			CanIndexBlocks:       false,
			StorageRequirementGB: 100,
			MemoryRequirementGB:  4,
		}
	case RoleService:
		return NodeCapabilities{
			CanValidate:          false, // Doesn't validate, executes services
			CanQueryState:        false, // Relies on full nodes
			CanQueryHistory:      false,
			CanServeLightClients: false,
			CanRunPlugins:        true, // All service plugins
			CanIndexBlocks:       false,
			StorageRequirementGB: 50,
			MemoryRequirementGB:  8,
		}
	case RoleArchive:
		return NodeCapabilities{
			CanValidate:          true,
			CanQueryState:        true,
			CanQueryHistory:      true,
			CanServeLightClients: true,
			CanRunPlugins:        false,
			CanIndexBlocks:       true,
			StorageRequirementGB: 1000,
			MemoryRequirementGB:  16,
		}
	default:
		return NodeCapabilities{}
	}
}
