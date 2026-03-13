package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTxCmd(t *testing.T) {
	cmd := GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "taskmarket", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestGetQueryCmdTaskmarketTx(t *testing.T) {
	cmd := GetQueryCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "taskmarket", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	// Test subcommands are registered
	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)
}

func TestCmdCreateTask(t *testing.T) {
	cmd := CmdCreateTask()
	require.NotNil(t, cmd)
	require.Equal(t, "create-task [title] [description] [type] [budget]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags are registered
	require.NotNil(t, cmd.Flag("category"))
	require.NotNil(t, cmd.Flag("skills"))
	require.NotNil(t, cmd.Flag("deadline"))
}

func TestCmdUpdateTask(t *testing.T) {
	cmd := CmdUpdateTask()
	require.NotNil(t, cmd)
	require.Equal(t, "update-task [task-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdPublishTask(t *testing.T) {
	cmd := CmdPublishTask()
	require.NotNil(t, cmd)
	require.Equal(t, "publish-task [task-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdCancelTask(t *testing.T) {
	cmd := CmdCancelTask()
	require.NotNil(t, cmd)
	require.Equal(t, "cancel-task [task-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdSubmitApplication(t *testing.T) {
	cmd := CmdSubmitApplication()
	require.NotNil(t, cmd)
	require.Equal(t, "submit-application [task-id] [price]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags
	require.NotNil(t, cmd.Flag("cover-letter"))
	require.NotNil(t, cmd.Flag("duration"))
}

func TestCmdAcceptApplication(t *testing.T) {
	cmd := CmdAcceptApplication()
	require.NotNil(t, cmd)
	require.Equal(t, "accept-application [application-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdRejectApplication(t *testing.T) {
	cmd := CmdRejectApplication()
	require.NotNil(t, cmd)
	require.Equal(t, "reject-application [application-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdSubmitBid(t *testing.T) {
	cmd := CmdSubmitBid()
	require.NotNil(t, cmd)
	require.Equal(t, "submit-bid [task-id] [amount]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags
	require.NotNil(t, cmd.Flag("message"))
}

func TestCmdCloseAuction(t *testing.T) {
	cmd := CmdCloseAuction()
	require.NotNil(t, cmd)
	require.Equal(t, "close-auction [task-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdStartTask(t *testing.T) {
	cmd := CmdStartTask()
	require.NotNil(t, cmd)
	require.Equal(t, "start-task [task-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdSubmitMilestone(t *testing.T) {
	cmd := CmdSubmitMilestone()
	require.NotNil(t, cmd)
	require.Equal(t, "submit-milestone [task-id] [milestone-id] [deliverables]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdApproveMilestone(t *testing.T) {
	cmd := CmdApproveMilestone()
	require.NotNil(t, cmd)
	require.Equal(t, "approve-milestone [task-id] [milestone-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdRejectMilestone(t *testing.T) {
	cmd := CmdRejectMilestone()
	require.NotNil(t, cmd)
	require.Equal(t, "reject-milestone [task-id] [milestone-id] [reason]", cmd.Use)
	require.NotEmpty(t, cmd.Short)
}

func TestCmdSubmitRating(t *testing.T) {
	cmd := CmdSubmitRating()
	require.NotNil(t, cmd)
	require.Equal(t, "submit-rating [task-id] [rated-id]", cmd.Use)
	require.NotEmpty(t, cmd.Short)

	// Check flags
	require.NotNil(t, cmd.Flag("quality"))
	require.NotNil(t, cmd.Flag("communication"))
	require.NotNil(t, cmd.Flag("timeliness"))
	require.NotNil(t, cmd.Flag("professionalism"))
	require.NotNil(t, cmd.Flag("comment"))
}
