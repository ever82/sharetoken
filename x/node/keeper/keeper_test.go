package keeper_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/node/keeper"
	"sharetoken/x/node/types"
)

func TestNodeKeeper_RoleInitialization(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)
	require.NotNil(t, k)

	// Initially undefined
	require.Equal(t, types.RoleUndefined, k.GetCurrentRole())

	// Test each role
	roles := []types.NodeRole{
		types.RoleLight,
		types.RoleFull,
		types.RoleService,
		types.RoleArchive,
	}

	for _, role := range roles {
		t.Run(role.String(), func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "node_role.json")

			k, err := keeper.NewNodeKeeper(configPath)
			require.NoError(t, err)

			err = k.InitializeRole(role, nil)
			require.NoError(t, err)
			require.Equal(t, role, k.GetCurrentRole())

			// Config should be saved
			_, err = os.Stat(configPath)
			require.NoError(t, err)
		})
	}
}

func TestNodeKeeper_InvalidRole(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	// Invalid role
	err = k.InitializeRole(types.RoleUndefined, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid role")
}

func TestNodeKeeper_DoubleInitialization(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	// First initialization
	err = k.InitializeRole(types.RoleLight, nil)
	require.NoError(t, err)

	// Try to initialize again (should fail - needs restart)
	err = k.InitializeRole(types.RoleFull, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already initialized")
}

func TestNodeKeeper_Capabilities(t *testing.T) {
	tests := []struct {
		role               types.NodeRole
		canValidate        bool
		canQueryState      bool
		canQueryHistory    bool
		canServeLight      bool
		canRunPlugins      bool
		canIndexBlocks     bool
		storageRequirement int
		memoryRequirement  int
	}{
		{
			role:               types.RoleLight,
			canValidate:        false,
			canQueryState:      false,
			canQueryHistory:    false,
			canServeLight:      false,
			canRunPlugins:      true,
			canIndexBlocks:     false,
			storageRequirement: 10,
			memoryRequirement:  2,
		},
		{
			role:               types.RoleFull,
			canValidate:        true,
			canQueryState:      true,
			canQueryHistory:    false,
			canServeLight:      true,
			canRunPlugins:      false,
			canIndexBlocks:     false,
			storageRequirement: 100,
			memoryRequirement:  4,
		},
		{
			role:               types.RoleService,
			canValidate:        false,
			canQueryState:      false,
			canQueryHistory:    false,
			canServeLight:      false,
			canRunPlugins:      true,
			canIndexBlocks:     false,
			storageRequirement: 50,
			memoryRequirement:  8,
		},
		{
			role:               types.RoleArchive,
			canValidate:        true,
			canQueryState:      true,
			canQueryHistory:    true,
			canServeLight:      true,
			canRunPlugins:      false,
			canIndexBlocks:     true,
			storageRequirement: 1000,
			memoryRequirement:  16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.role.String(), func(t *testing.T) {
			caps := tt.role.GetCapabilities()
			require.Equal(t, tt.canValidate, caps.CanValidate)
			require.Equal(t, tt.canQueryState, caps.CanQueryState)
			require.Equal(t, tt.canQueryHistory, caps.CanQueryHistory)
			require.Equal(t, tt.canServeLight, caps.CanServeLightClients)
			require.Equal(t, tt.canRunPlugins, caps.CanRunPlugins)
			require.Equal(t, tt.canIndexBlocks, caps.CanIndexBlocks)
			require.Equal(t, tt.storageRequirement, caps.StorageRequirementGB)
			require.Equal(t, tt.memoryRequirement, caps.MemoryRequirementGB)
		})
	}
}

func TestNodeKeeper_RoleSwitching(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	// Test switching without initialization
	canHot, requiresRestart, err := k.CanSwitchRole(types.RoleFull)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not initialized")
	require.False(t, canHot)
	require.False(t, requiresRestart)

	// Initialize
	err = k.InitializeRole(types.RoleLight, nil)
	require.NoError(t, err)

	// Test Light -> Service (can hot switch)
	canHot, requiresRestart, err = k.CanSwitchRole(types.RoleService)
	require.NoError(t, err)
	require.True(t, canHot)
	require.False(t, requiresRestart)

	// Test Light -> Full (requires restart)
	canHot, requiresRestart, err = k.CanSwitchRole(types.RoleFull)
	require.NoError(t, err)
	require.False(t, canHot)
	require.True(t, requiresRestart)

	// Test Light -> Archive (requires restart)
	canHot, requiresRestart, err = k.CanSwitchRole(types.RoleArchive)
	require.NoError(t, err)
	require.False(t, canHot)
	require.True(t, requiresRestart)

	// Test invalid role
	_, _, err = k.CanSwitchRole(types.RoleUndefined)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid")
}

func TestNodeKeeper_RoleSwitchForRestart(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	// Initialize as Light
	err = k.InitializeRole(types.RoleLight, nil)
	require.NoError(t, err)
	require.Equal(t, types.RoleLight, k.GetCurrentRole())

	// Request switch to Full
	err = k.RequestRoleSwitchForRestart(types.RoleFull, nil)
	require.NoError(t, err)

	// Role should still be Light (until restart)
	require.Equal(t, types.RoleLight, k.GetCurrentRole())

	// Config should be saved
	_, err = os.Stat(configPath)
	require.NoError(t, err)

	// Create new keeper and verify it loads the new role
	k2, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)
	require.Equal(t, types.RoleFull, k2.GetCurrentRole())
}

func TestNodeKeeper_ConfigPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	// Create custom config
	customConfig := &types.RoleConfig{
		Role: types.RoleService,
		ServiceConfig: &types.ServiceNodeConfig{
			MaxPeers:             30,
			EnableLLMPlugin:      true,
			EnableAgentPlugin:    false,
			EnableWorkflowPlugin: true,
			WASMMaxMemory:        1024,
			WASMMaxCPUTime:       60000,
		},
	}

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	err = k.InitializeRole(types.RoleService, customConfig)
	require.NoError(t, err)

	// Create new keeper and verify config is loaded
	k2, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	loadedConfig := k2.GetCurrentConfig()
	require.NotNil(t, loadedConfig)
	require.Equal(t, types.RoleService, loadedConfig.Role)
	require.NotNil(t, loadedConfig.ServiceConfig)
	require.Equal(t, 30, loadedConfig.ServiceConfig.MaxPeers)
	require.Equal(t, true, loadedConfig.ServiceConfig.EnableLLMPlugin)
	require.Equal(t, false, loadedConfig.ServiceConfig.EnableAgentPlugin)
	require.Equal(t, 1024, loadedConfig.ServiceConfig.WASMMaxMemory)
}

func TestNodeKeeper_StartStop(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "node_role.json")

	k, err := keeper.NewNodeKeeper(configPath)
	require.NoError(t, err)

	// Try to start without initialization
	err = k.Start()
	require.Error(t, err)
	require.Contains(t, err.Error(), "not initialized")

	// Initialize
	err = k.InitializeRole(types.RoleLight, nil)
	require.NoError(t, err)

	// Start
	err = k.Start()
	require.NoError(t, err)
	require.Equal(t, keeper.StateRunning, k.GetState())

	// Stop
	err = k.Stop()
	require.NoError(t, err)
	require.Equal(t, keeper.StateStopped, k.GetState())
}

func TestRoleConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  *types.RoleConfig
		wantErr bool
	}{
		{
			name: "valid light config",
			config: &types.RoleConfig{
				Role:        types.RoleLight,
				LightConfig: types.DefaultLightConfig(),
			},
			wantErr: false,
		},
		{
			name: "light config missing",
			config: &types.RoleConfig{
				Role: types.RoleLight,
			},
			wantErr: true,
		},
		{
			name: "light invalid max_peers",
			config: &types.RoleConfig{
				Role: types.RoleLight,
				LightConfig: &types.LightNodeConfig{
					MaxPeers: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "valid service config",
			config: &types.RoleConfig{
				Role:          types.RoleService,
				ServiceConfig: types.DefaultServiceConfig(),
			},
			wantErr: false,
		},
		{
			name: "service invalid wasm memory",
			config: &types.RoleConfig{
				Role: types.RoleService,
				ServiceConfig: &types.ServiceNodeConfig{
					MaxPeers:      10,
					WASMMaxMemory: 32,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			config: &types.RoleConfig{
				Role: types.RoleUndefined,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNodeRole_String(t *testing.T) {
	tests := []struct {
		role types.NodeRole
		want string
	}{
		{types.RoleLight, "light"},
		{types.RoleFull, "full"},
		{types.RoleService, "service"},
		{types.RoleArchive, "archive"},
		{types.RoleUndefined, "undefined"},
		{99, "undefined"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.role.String()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNodeRole_FromString(t *testing.T) {
	tests := []struct {
		input    string
		expected types.NodeRole
		wantErr  bool
	}{
		{"light", types.RoleLight, false},
		{"LIGHT", types.RoleLight, false},
		{"Light", types.RoleLight, false},
		{"full", types.RoleFull, false},
		{"service", types.RoleService, false},
		{"archive", types.RoleArchive, false},
		{"invalid", types.RoleUndefined, true},
		{"", types.RoleUndefined, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := types.NodeRoleFromString(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestNodeRole_IsValid(t *testing.T) {
	require.True(t, types.RoleLight.IsValid())
	require.True(t, types.RoleFull.IsValid())
	require.True(t, types.RoleService.IsValid())
	require.True(t, types.RoleArchive.IsValid())
	require.False(t, types.RoleUndefined.IsValid())
	require.False(t, types.NodeRole(99).IsValid())
	require.False(t, types.NodeRole(-1).IsValid())
}
