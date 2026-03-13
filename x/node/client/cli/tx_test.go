package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTxCmd(t *testing.T) {
	cmd := GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "node", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestGetQueryCmdNodeTx(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "node", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestGetCmdInitializeRole(t *testing.T) {
	cmd := GetCmdInitializeRole()
	require.NotNil(t, cmd)
	require.Equal(t, "init-role [role]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check config-file flag exists
	flag := cmd.Flag("config-file")
	require.NotNil(t, flag)
}

func TestGetCmdSwitchRole(t *testing.T) {
	cmd := GetCmdSwitchRole()
	require.NotNil(t, cmd)
	require.Equal(t, "switch-role [role]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags exist
	flagConfig := cmd.Flag("config-file")
	require.NotNil(t, flagConfig)
	flagForce := cmd.Flag("force")
	require.NotNil(t, flagForce)
}

func TestGetCmdUpdateConfig(t *testing.T) {
	cmd := GetCmdUpdateConfig()
	require.NotNil(t, cmd)
	require.Equal(t, "update-config", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
