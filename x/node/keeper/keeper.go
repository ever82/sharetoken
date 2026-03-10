package keeper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"sharetoken/x/node/types"
)

// NodeKeeper manages node role and configuration
type NodeKeeper struct {
	mu sync.RWMutex

	// Current role
	currentRole types.NodeRole

	// Current configuration
	config *types.RoleConfig

	// Config file path
	configPath string

	// State tracking
	state NodeState

	// Plugin manager reference (would be actual plugin manager in production)
	pluginManager PluginManager

	// State callbacks
	roleChangeCallbacks []RoleChangeCallback
}

// NodeState represents the current state of the node
type NodeState int32

const (
	StateStopped NodeState = iota
	StateStarting
	StateRunning
	StateSwitchingRole
	StateStopping
)

func (s NodeState) String() string {
	switch s {
	case StateStopped:
		return "stopped"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateSwitchingRole:
		return "switching_role"
	case StateStopping:
		return "stopping"
	default:
		return "unknown"
	}
}

// RoleChangeCallback is called when role changes
type RoleChangeCallback func(oldRole, newRole types.NodeRole) error

// PluginManager interface for plugin operations
type PluginManager interface {
	LoadPlugin(name string, config map[string]interface{}) error
	UnloadPlugin(name string) error
	IsPluginLoaded(name string) bool
	ListLoadedPlugins() []string
	StopAll() error
}

// NewNodeKeeper creates a new node keeper
func NewNodeKeeper(configPath string) (*NodeKeeper, error) {
	k := &NodeKeeper{
		configPath:          configPath,
		currentRole:         types.RoleUndefined,
		roleChangeCallbacks: make([]RoleChangeCallback, 0),
	}

	// Try to load existing config
	if err := k.loadConfig(); err != nil {
		// Config doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return k, nil
}

// loadConfig loads configuration from file
func (k *NodeKeeper) loadConfig() error {
	data, err := os.ReadFile(k.configPath)
	if err != nil {
		return err
	}

	var config types.RoleConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	k.config = &config
	k.currentRole = config.Role
	return nil
}

// saveConfig saves configuration to file
func (k *NodeKeeper) saveConfig() error {
	if k.config == nil {
		return fmt.Errorf("no config to save")
	}

	// Ensure directory exists
	dir := filepath.Dir(k.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil { //nolint:gosec
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := json.MarshalIndent(k.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(k.configPath, data, 0644); err != nil { //nolint:gosec
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetCurrentRole returns the current node role
func (k *NodeKeeper) GetCurrentRole() types.NodeRole {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.currentRole
}

// GetCurrentConfig returns the current configuration
func (k *NodeKeeper) GetCurrentConfig() *types.RoleConfig {
	k.mu.RLock()
	defer k.mu.RUnlock()
	if k.config == nil {
		return nil
	}
	// Return a copy
	configCopy := *k.config
	return &configCopy
}

// GetState returns the current node state
func (k *NodeKeeper) GetState() NodeState {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.state
}

// InitializeRole initializes the node with a specific role
func (k *NodeKeeper) InitializeRole(role types.NodeRole, customConfig *types.RoleConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateStopped {
		return fmt.Errorf("node must be stopped to initialize role, current state: %s", k.state.String())
	}

	if !role.IsValid() {
		return fmt.Errorf("invalid role: %s", role.String())
	}

	// Check if already initialized
	if k.currentRole != types.RoleUndefined {
		return fmt.Errorf("node already initialized with role: %s. Use 'switch-role' to change", k.currentRole.String())
	}

	// Use provided config or create default
	var config *types.RoleConfig
	if customConfig != nil {
		config = customConfig
	} else {
		config = types.DefaultRoleConfig(role)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	k.config = config
	k.currentRole = role

	// Save config
	if err := k.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// CanSwitchRole checks if role switch is possible
func (k *NodeKeeper) CanSwitchRole(targetRole types.NodeRole) (canHotSwitch bool, requiresRestart bool, err error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.currentRole == types.RoleUndefined {
		return false, false, fmt.Errorf("node role not initialized")
	}

	if !targetRole.IsValid() {
		return false, false, fmt.Errorf("invalid target role: %s", targetRole.String())
	}

	canHot, needsRestart := k.currentRole.CanSwitchTo(targetRole)
	return canHot, needsRestart, nil
}

// SwitchRoleHot switches role without restarting (if supported)
func (k *NodeKeeper) SwitchRoleHot(targetRole types.NodeRole, targetConfig *types.RoleConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateRunning {
		return fmt.Errorf("node must be running to hot switch, current state: %s", k.state.String())
	}

	canHot, _ := k.currentRole.CanSwitchTo(targetRole)
	if !canHot {
		return fmt.Errorf("cannot hot switch from %s to %s", k.currentRole.String(), targetRole.String())
	}

	// Set switching state
	oldRole := k.currentRole
	k.state = StateSwitchingRole

	// Execute callbacks before switch
	for _, callback := range k.roleChangeCallbacks {
		if err := callback(oldRole, targetRole); err != nil {
			k.state = StateRunning
			return fmt.Errorf("role change callback failed: %w", err)
		}
	}

	// Perform the switch
	var config *types.RoleConfig
	if targetConfig != nil {
		config = targetConfig
	} else {
		config = types.DefaultRoleConfig(targetRole)
	}

	if err := config.Validate(); err != nil {
		k.state = StateRunning
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Update plugins based on role
	if err := k.updatePluginsForRole(targetRole); err != nil {
		k.state = StateRunning
		return fmt.Errorf("failed to update plugins: %w", err)
	}

	k.config = config
	k.currentRole = targetRole

	// Save config
	if err := k.saveConfig(); err != nil {
		k.state = StateRunning
		return fmt.Errorf("failed to save config: %w", err)
	}

	k.state = StateRunning
	return nil
}

// RequestRoleSwitchForRestart requests a role switch that will take effect on next restart
func (k *NodeKeeper) RequestRoleSwitchForRestart(targetRole types.NodeRole, targetConfig *types.RoleConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if !targetRole.IsValid() {
		return fmt.Errorf("invalid target role: %s", targetRole.String())
	}

	var config *types.RoleConfig
	if targetConfig != nil {
		config = targetConfig
	} else {
		config = types.DefaultRoleConfig(targetRole)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Save the new config (will be applied on restart)
	config.Role = targetRole
	k.config = config

	if err := k.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// updatePluginsForRole updates loaded plugins based on role
func (k *NodeKeeper) updatePluginsForRole(role types.NodeRole) error {
	if k.pluginManager == nil {
		return nil
	}

	switch role {
	case types.RoleLight:
		// Load GenieBot plugin
		if !k.pluginManager.IsPluginLoaded("geniebot") {
			if err := k.pluginManager.LoadPlugin("geniebot", map[string]interface{}{}); err != nil {
				return fmt.Errorf("failed to load geniebot plugin: %w", err)
			}
		}
		// Unload service plugins if loaded
		_ = k.pluginManager.UnloadPlugin("llm")      //nolint:errcheck
		_ = k.pluginManager.UnloadPlugin("agent")    //nolint:errcheck
		_ = k.pluginManager.UnloadPlugin("workflow") //nolint:errcheck

	case types.RoleService:
		// Load service plugins
		if k.config.ServiceConfig != nil {
			if k.config.ServiceConfig.EnableLLMPlugin {
				if !k.pluginManager.IsPluginLoaded("llm") {
					_ = k.pluginManager.LoadPlugin("llm", map[string]interface{}{}) //nolint:errcheck
				}
			}
			if k.config.ServiceConfig.EnableAgentPlugin {
				if !k.pluginManager.IsPluginLoaded("agent") {
					_ = k.pluginManager.LoadPlugin("agent", map[string]interface{}{}) //nolint:errcheck
				}
			}
			if k.config.ServiceConfig.EnableWorkflowPlugin {
				if !k.pluginManager.IsPluginLoaded("workflow") {
					_ = k.pluginManager.LoadPlugin("workflow", map[string]interface{}{}) //nolint:errcheck
				}
			}
		}
		// Unload geniebot if loaded
		_ = k.pluginManager.UnloadPlugin("geniebot") //nolint:errcheck

	default:
		// Other roles don't run plugins
		_ = k.pluginManager.StopAll() //nolint:errcheck
	}

	return nil
}

// Start starts the node with current role
func (k *NodeKeeper) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateStopped {
		return fmt.Errorf("node is not stopped, current state: %s", k.state.String())
	}

	if k.currentRole == types.RoleUndefined {
		return fmt.Errorf("node role not initialized")
	}

	k.state = StateStarting

	// Load appropriate plugins for the role
	if err := k.updatePluginsForRole(k.currentRole); err != nil {
		k.state = StateStopped
		return fmt.Errorf("failed to load plugins: %w", err)
	}

	k.state = StateRunning
	return nil
}

// Stop stops the node
func (k *NodeKeeper) Stop() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateRunning && k.state != StateSwitchingRole {
		return fmt.Errorf("node is not running, current state: %s", k.state.String())
	}

	k.state = StateStopping

	// Stop all plugins
	if k.pluginManager != nil {
		if err := k.pluginManager.StopAll(); err != nil {
			k.state = StateRunning
			return fmt.Errorf("failed to stop plugins: %w", err)
		}
	}

	k.state = StateStopped
	return nil
}

// RegisterRoleChangeCallback registers a callback for role changes
func (k *NodeKeeper) RegisterRoleChangeCallback(callback RoleChangeCallback) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.roleChangeCallbacks = append(k.roleChangeCallbacks, callback)
}

// SetPluginManager sets the plugin manager
func (k *NodeKeeper) SetPluginManager(pm PluginManager) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.pluginManager = pm
}

// GetCapabilities returns capabilities for current role
func (k *NodeKeeper) GetCapabilities() types.NodeCapabilities {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.currentRole.GetCapabilities()
}

// GetRoleInfo returns comprehensive role information
func (k *NodeKeeper) GetRoleInfo() map[string]interface{} {
	k.mu.RLock()
	defer k.mu.RUnlock()

	caps := k.currentRole.GetCapabilities()

	return map[string]interface{}{
		"role":           k.currentRole.String(),
		"state":          k.state.String(),
		"capabilities":   caps,
		"can_hot_switch": k.currentRole.CanSwitchTo,
	}
}
