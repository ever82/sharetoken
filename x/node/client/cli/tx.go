package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"sharetoken/x/node/keeper"
	"sharetoken/x/node/types"
)

// GetQueryCmd returns the query commands for node module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "node",
		Short:                      "Querying commands for the node module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryRole(),
		GetCmdQueryCapabilities(),
		GetCmdQueryState(),
	)

	return cmd
}

// GetTxCmd returns the transaction commands for node module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "node",
		Short:                      "Node role management transactions",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdInitializeRole(),
		GetCmdSwitchRole(),
		GetCmdUpdateConfig(),
	)

	return cmd
}

// GetCmdQueryRole returns the query role command
func GetCmdQueryRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Query current node role",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// In a real implementation, this would query the node via gRPC/REST
			// For now, we'll read from the local config
			k, err := keeper.NewNodeKeeper(getConfigPath(clientCtx))
			if err != nil {
				return fmt.Errorf("failed to create node keeper: %w", err)
			}

			role := k.GetCurrentRole()
			if role == types.RoleUndefined {
				return fmt.Errorf("node role not initialized")
			}

			fmt.Printf("Current Role: %s\n", role.String())
			fmt.Printf("Role Value: %d\n", role)

			return nil
		},
	}

	return cmd
}

// GetCmdQueryCapabilities returns the query capabilities command
func GetCmdQueryCapabilities() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capabilities [role]",
		Short: "Query capabilities for a role (or current role if not specified)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var role types.NodeRole
			var err error

			if len(args) > 0 {
				role, err = types.NodeRoleFromString(args[0])
				if err != nil {
					return err
				}
			} else {
				clientCtx, err := client.GetClientQueryContext(cmd)
				if err != nil {
					return err
				}

				k, err := keeper.NewNodeKeeper(getConfigPath(clientCtx))
				if err != nil {
					return fmt.Errorf("failed to create node keeper: %w", err)
				}

				role = k.GetCurrentRole()
				if role == types.RoleUndefined {
					return fmt.Errorf("node role not initialized")
				}
			}

			caps := role.GetCapabilities()

			fmt.Printf("Capabilities for role '%s':\n", role.String())
			fmt.Printf("  Can Validate:          %v\n", caps.CanValidate)
			fmt.Printf("  Can Query State:       %v\n", caps.CanQueryState)
			fmt.Printf("  Can Query History:     %v\n", caps.CanQueryHistory)
			fmt.Printf("  Can Serve Light:       %v\n", caps.CanServeLightClients)
			fmt.Printf("  Can Run Plugins:       %v\n", caps.CanRunPlugins)
			fmt.Printf("  Can Index Blocks:      %v\n", caps.CanIndexBlocks)
			fmt.Printf("  Storage Required:      %d GB\n", caps.StorageRequirementGB)
			fmt.Printf("  Memory Required:       %d GB\n", caps.MemoryRequirementGB)

			return nil
		},
	}

	return cmd
}

// GetCmdQueryState returns the query state command
func GetCmdQueryState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Query current node state",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			k, err := keeper.NewNodeKeeper(getConfigPath(clientCtx))
			if err != nil {
				return fmt.Errorf("failed to create node keeper: %w", err)
			}

			state := k.GetState()
			role := k.GetCurrentRole()
			caps := k.GetCapabilities()

			info := map[string]interface{}{
				"role":        role.String(),
				"state":       state.String(),
				"config_path": getConfigPath(clientCtx),
				"capabilities": map[string]interface{}{
					"can_validate":     caps.CanValidate,
					"can_query_state":  caps.CanQueryState,
					"can_run_plugins":  caps.CanRunPlugins,
					"storage_gb":       caps.StorageRequirementGB,
					"memory_gb":        caps.MemoryRequirementGB,
				},
			}

			jsonData, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(jsonData))
			return nil
		},
	}

	return cmd
}

// GetCmdInitializeRole returns the initialize role command
func GetCmdInitializeRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-role [role]",
		Short: "Initialize node with a specific role",
		Long: `Initialize the node with a specific role. Valid roles are:
  - light: Light node (minimal storage, runs core + GenieBot)
  - full: Full node (complete state and current history)
  - service: Service node (core + plugins for service provision)
  - archive: Archive node (full history with indexes)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			role, err := types.NodeRoleFromString(args[0])
			if err != nil {
				return err
			}

			k, err := keeper.NewNodeKeeper(getConfigPath(clientCtx))
			if err != nil {
				return fmt.Errorf("failed to create node keeper: %w", err)
			}

			// Check if already initialized
			if k.GetCurrentRole() != types.RoleUndefined {
				return fmt.Errorf("node already initialized with role: %s. Use 'switch-role' to change", k.GetCurrentRole().String())
			}

			// Get custom config if provided
			configPath, _ := cmd.Flags().GetString("config-file")
			var customConfig *types.RoleConfig
			if configPath != "" {
				data, err := os.ReadFile(configPath)
				if err != nil {
					return fmt.Errorf("failed to read config file: %w", err)
				}
				customConfig = &types.RoleConfig{}
				if err := json.Unmarshal(data, customConfig); err != nil {
					return fmt.Errorf("failed to parse config file: %w", err)
				}
			}

			if err := k.InitializeRole(role, customConfig); err != nil {
				return fmt.Errorf("failed to initialize role: %w", err)
			}

			fmt.Printf("Successfully initialized node with role: %s\n", role.String())
			fmt.Printf("Config saved to: %s\n", getConfigPath(clientCtx))
			fmt.Println("\nTo start the node, run: sharetokend start")

			return nil
		},
	}

	cmd.Flags().String("config-file", "", "Path to custom role configuration JSON file")

	return cmd
}

// GetCmdSwitchRole returns the switch role command
func GetCmdSwitchRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch-role [role]",
		Short: "Switch node to a different role",
		Long: `Switch the node to a different role. Some role changes can be done hot (without restart),
while others require a restart. This command will tell you if a restart is needed.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			targetRole, err := types.NodeRoleFromString(args[0])
			if err != nil {
				return err
			}

			k, err := keeper.NewNodeKeeper(getConfigPath(clientCtx))
			if err != nil {
				return fmt.Errorf("failed to create node keeper: %w", err)
			}

			currentRole := k.GetCurrentRole()
			if currentRole == types.RoleUndefined {
				return fmt.Errorf("node role not initialized. Use 'init-role' first")
			}

			if currentRole == targetRole {
				fmt.Printf("Node is already running as %s\n", targetRole.String())
				return nil
			}

			// Check if hot switch is possible
			canHot, requiresRestart, err := k.CanSwitchRole(targetRole)
			if err != nil {
				return err
			}

			// Get custom config if provided
			configPath, _ := cmd.Flags().GetString("config-file")
			var customConfig *types.RoleConfig
			if configPath != "" {
				data, err := os.ReadFile(configPath)
				if err != nil {
					return fmt.Errorf("failed to read config file: %w", err)
				}
				customConfig = &types.RoleConfig{}
				if err := json.Unmarshal(data, customConfig); err != nil {
					return fmt.Errorf("failed to parse config file: %w", err)
				}
			}

			force, _ := cmd.Flags().GetBool("force")

			if requiresRestart {
				// Request switch for next restart
				if err := k.RequestRoleSwitchForRestart(targetRole, customConfig); err != nil {
					return fmt.Errorf("failed to request role switch: %w", err)
				}

				fmt.Printf("Role switch from '%s' to '%s' scheduled\n", currentRole.String(), targetRole.String())
				fmt.Println("Configuration saved. The new role will take effect on next restart.")
				fmt.Println("\nTo apply the change, restart the node:")
				fmt.Println("  sharetokend stop && sharetokend start")

				return nil
			}

			// Can hot switch
			if !canHot {
				return fmt.Errorf("cannot switch from '%s' to '%s'", currentRole.String(), targetRole.String())
			}

			if !force {
				fmt.Printf("Can hot switch from '%s' to '%s'\n", currentRole.String(), targetRole.String())
				fmt.Println("Use --force to perform the switch immediately")
				return nil
			}

			// Perform hot switch
			if err := k.SwitchRoleHot(targetRole, customConfig); err != nil {
				return fmt.Errorf("failed to switch role: %w", err)
			}

			fmt.Printf("Successfully switched node role from '%s' to '%s'\n", currentRole.String(), targetRole.String())
			return nil
		},
	}

	cmd.Flags().String("config-file", "", "Path to custom role configuration JSON file")
	cmd.Flags().Bool("force", false, "Perform hot switch immediately")

	return cmd
}

// GetCmdUpdateConfig returns the update config command
func GetCmdUpdateConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "Update current role configuration",
		Long:  `Update the configuration for the current role without changing the role itself.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			configPath, _ := cmd.Flags().GetString("config-file")
			if configPath == "" {
				return fmt.Errorf("--config-file is required")
			}

			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			newConfig := &types.RoleConfig{}
			if err := json.Unmarshal(data, newConfig); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}

			k, err := keeper.NewNodeKeeper(getConfigPath(clientCtx))
			if err != nil {
				return fmt.Errorf("failed to create node keeper: %w", err)
			}

			currentRole := k.GetCurrentRole()
			if currentRole == types.RoleUndefined {
				return fmt.Errorf("node role not initialized. Use 'init-role' first")
			}

			// Ensure role matches
			if newConfig.Role != currentRole {
				return fmt.Errorf("config role (%s) does not match current role (%s). Use 'switch-role' instead", newConfig.Role.String(), currentRole.String())
			}

			// Re-initialize with new config
			if err := k.InitializeRole(currentRole, newConfig); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}

			fmt.Printf("Successfully updated configuration for role: %s\n", currentRole.String())
			return nil
		},
	}

	cmd.Flags().String("config-file", "", "Path to new configuration JSON file (required)")
	cmd.MarkFlagRequired("config-file")

	return cmd
}

// getConfigPath returns the config path for the node
func getConfigPath(clientCtx client.Context) string {
	return filepath.Join(clientCtx.HomeDir, "config", "node_role.json")
}
