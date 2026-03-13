package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetQueryCmdTaskmarketQuery(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "taskmarket", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestCmdQueryTaskQuery(t *testing.T) {
	cmd := CmdQueryTask()
	require.NotNil(t, cmd)
	require.Equal(t, "task [id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryTasksQuery(t *testing.T) {
	cmd := CmdQueryTasks()
	require.NotNil(t, cmd)
	require.Equal(t, "tasks", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags
	require.NotNil(t, cmd.Flag("status"))
	require.NotNil(t, cmd.Flag("category"))
	require.NotNil(t, cmd.Flag("requester"))
	require.NotNil(t, cmd.Flag("worker"))
}

func TestCmdQueryApplicationsQuery(t *testing.T) {
	cmd := CmdQueryApplications()
	require.NotNil(t, cmd)
	require.Equal(t, "applications", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags
	require.NotNil(t, cmd.Flag("task-id"))
	require.NotNil(t, cmd.Flag("worker"))
}

func TestCmdQueryAuctionQuery(t *testing.T) {
	cmd := CmdQueryAuction()
	require.NotNil(t, cmd)
	require.Equal(t, "auction [task-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryReputationQuery(t *testing.T) {
	cmd := CmdQueryReputation()
	require.NotNil(t, cmd)
	require.Equal(t, "reputation [user-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdQueryStatisticsQuery(t *testing.T) {
	cmd := CmdQueryStatistics()
	require.NotNil(t, cmd)
	require.Equal(t, "statistics", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}
