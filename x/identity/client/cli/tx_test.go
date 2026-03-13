package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTxCmd(t *testing.T) {
	cmd := GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "identity", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}

func TestGetQueryCmdIdentityTx(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "identity", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}

func TestCmdRegisterIdentity(t *testing.T) {
	cmd := CmdRegisterIdentity()
	require.NotNil(t, cmd)
	require.Equal(t, "register-identity [address] [did] [metadata-hash]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdVerifyIdentity(t *testing.T) {
	cmd := CmdVerifyIdentity()
	require.NotNil(t, cmd)
	require.Equal(t, "verify-identity [address] [provider] [verification-hash] [proof]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
