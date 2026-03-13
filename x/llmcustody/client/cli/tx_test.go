package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTxCmd(t *testing.T) {
	cmd := GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "llmcustody", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestGetQueryCmdLlmcustodyTx(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "llmcustody", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestCmdRegisterAPIKey(t *testing.T) {
	cmd := CmdRegisterAPIKey()
	require.NotNil(t, cmd)
	require.Equal(t, "register-api-key [provider] [encrypted-key-file]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdUpdateAPIKey(t *testing.T) {
	cmd := CmdUpdateAPIKey()
	require.NotNil(t, cmd)
	require.Equal(t, "update-api-key [api-key-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags are registered
	flag := cmd.Flag("active")
	require.NotNil(t, flag)
}

func TestCmdRevokeAPIKey(t *testing.T) {
	cmd := CmdRevokeAPIKey()
	require.NotNil(t, cmd)
	require.Equal(t, "revoke-api-key [api-key-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
