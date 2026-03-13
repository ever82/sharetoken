package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetQueryCmdEscrowQuery(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "escrow", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
