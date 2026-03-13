package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetQueryCmdNodeQuery(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "node", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}

func TestGetCmdQueryRoleQuery(t *testing.T) {
	cmd := GetCmdQueryRole()
	require.NotNil(t, cmd)
	require.Equal(t, "role", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestGetCmdQueryCapabilitiesQuery(t *testing.T) {
	cmd := GetCmdQueryCapabilities()
	require.NotNil(t, cmd)
	require.Equal(t, "capabilities [role]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestGetCmdQueryStateQuery(t *testing.T) {
	cmd := GetCmdQueryState()
	require.NotNil(t, cmd)
	require.Equal(t, "state", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
