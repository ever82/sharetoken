package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTxCmd(t *testing.T) {
	cmd := GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "sharetoken", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}

func TestGetQueryCmdSharetokenTx(t *testing.T) {
	cmd := GetQueryCmd("")
	require.NotNil(t, cmd)
	require.Equal(t, "sharetoken", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}

func TestCmdQueryParamsSharetokenTx(t *testing.T) {
	cmd := CmdQueryParams()
	require.NotNil(t, cmd)
	require.Equal(t, "params", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
