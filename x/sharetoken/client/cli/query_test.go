package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdQueryParamsSharetokenQuery(t *testing.T) {
	cmd := CmdQueryParams()
	require.NotNil(t, cmd)
	require.Equal(t, "params", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestGetQueryCmdSharetokenQuery(t *testing.T) {
	cmd := GetQueryCmd("")
	require.NotNil(t, cmd)
	require.Equal(t, "sharetoken", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
