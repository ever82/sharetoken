package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdQueryAPIKeyQuery(t *testing.T) {
	cmd := CmdQueryAPIKey()
	require.NotNil(t, cmd)
	require.Equal(t, "api-key [id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryAPIKeysByOwnerQuery(t *testing.T) {
	cmd := CmdQueryAPIKeysByOwner()
	require.NotNil(t, cmd)
	require.Equal(t, "api-keys-by-owner [owner-address]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryAllAPIKeysQuery(t *testing.T) {
	cmd := CmdQueryAllAPIKeys()
	require.NotNil(t, cmd)
	require.Equal(t, "all-api-keys", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestGetQueryCmdLlmcustodyQuery(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "llmcustody", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}
