package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdentityQueryCommands(t *testing.T) {
	// Test that all query commands can be created without errors
	tests := []struct {
		name string
		cmd  func()
	}{
		{
			name: "CmdQueryIdentity",
			cmd:  func() { CmdQueryIdentity() },
		},
		{
			name: "CmdQueryIdentities",
			cmd:  func() { CmdQueryIdentities() },
		},
		{
			name: "CmdQueryLimitConfig",
			cmd:  func() { CmdQueryLimitConfig() },
		},
		{
			name: "CmdQueryIsVerified",
			cmd:  func() { CmdQueryIsVerified() },
		},
		{
			name: "CmdQueryParams",
			cmd:  func() { CmdQueryParams() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			tt.cmd()
		})
	}
}

func TestCmdQueryIdentityQuery(t *testing.T) {
	cmd := CmdQueryIdentity()
	require.NotNil(t, cmd)
	require.Equal(t, "identity [address]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryIdentitiesQuery(t *testing.T) {
	cmd := CmdQueryIdentities()
	require.NotNil(t, cmd)
	require.Equal(t, "identities", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryLimitConfigQuery(t *testing.T) {
	cmd := CmdQueryLimitConfig()
	require.NotNil(t, cmd)
	require.Equal(t, "limit-config [address]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryIsVerifiedQuery(t *testing.T) {
	cmd := CmdQueryIsVerified()
	require.NotNil(t, cmd)
	require.Equal(t, "is-verified [address]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryParamsIdentity(t *testing.T) {
	cmd := CmdQueryParams()
	require.NotNil(t, cmd)
	require.Equal(t, "params", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestGetQueryCmdIdentityQuery(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "identity", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}
