package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	taskmarketQueryCmd := &cobra.Command{
		Use:                        "taskmarket",
		Short:                      "Querying commands for the taskmarket module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	taskmarketQueryCmd.AddCommand(CmdQueryTask())
	taskmarketQueryCmd.AddCommand(CmdQueryTasks())
	taskmarketQueryCmd.AddCommand(CmdQueryApplications())
	taskmarketQueryCmd.AddCommand(CmdQueryAuction())
	taskmarketQueryCmd.AddCommand(CmdQueryReputation())
	taskmarketQueryCmd.AddCommand(CmdQueryStatistics())

	return taskmarketQueryCmd
}

// CmdQueryTask implements the query task command
func CmdQueryTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task [id]",
		Short: "Query a task by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Querying task: %s\n", taskID)
			fmt.Println("Note: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryTasks implements the query tasks command
func CmdQueryTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Query tasks with optional filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := cmd.Flags().GetString("status")
			category, _ := cmd.Flags().GetString("category")
			requester, _ := cmd.Flags().GetString("requester")

			fmt.Printf("Querying tasks:\n")
			if status != "" {
				fmt.Printf("  Status: %s\n", status)
			}
			if category != "" {
				fmt.Printf("  Category: %s\n", category)
			}
			if requester != "" {
				fmt.Printf("  Requester: %s\n", requester)
			}
			fmt.Println("\nNote: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	cmd.Flags().String("status", "", "Filter by status (draft, open, assigned, in_progress, completed)")
	cmd.Flags().String("category", "", "Filter by category")
	cmd.Flags().String("requester", "", "Filter by requester address")
	cmd.Flags().String("worker", "", "Filter by worker address")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryApplications implements the query applications command
func CmdQueryApplications() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applications",
		Short: "Query applications for a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, _ := cmd.Flags().GetString("task-id")
			workerID, _ := cmd.Flags().GetString("worker")

			fmt.Printf("Querying applications:\n")
			if taskID != "" {
				fmt.Printf("  Task ID: %s\n", taskID)
			}
			if workerID != "" {
				fmt.Printf("  Worker ID: %s\n", workerID)
			}
			fmt.Println("\nNote: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	cmd.Flags().String("task-id", "", "Filter by task ID")
	cmd.Flags().String("worker", "", "Filter by worker address")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryAuction implements the query auction command
func CmdQueryAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction [task-id]",
		Short: "Query auction details for a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Querying auction for task: %s\n", taskID)
			fmt.Println("Note: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryReputation implements the query reputation command
func CmdQueryReputation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reputation [user-id]",
		Short: "Query a user's reputation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID := args[0]

			fmt.Printf("Querying reputation for user: %s\n", userID)
			fmt.Println("Note: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryStatistics implements the query statistics command
func CmdQueryStatistics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Query marketplace statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Querying marketplace statistics:")
			fmt.Println("  Total Tasks: --")
			fmt.Println("  Open Tasks: --")
			fmt.Println("  Assigned Tasks: --")
			fmt.Println("  In Progress Tasks: --")
			fmt.Println("  Completed Tasks: --")
			fmt.Println("  Total Applications: --")
			fmt.Println("  Total Bids: --")
			fmt.Println("  Total Ratings: --")
			fmt.Println("\nNote: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
