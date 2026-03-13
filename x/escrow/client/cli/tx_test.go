package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTxCmd(t *testing.T) {
	cmd := GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "escrow", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestGetQueryCmdEscrow(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "escrow", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
