package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	taskmarketTxCmd := &cobra.Command{
		Use:                        "taskmarket",
		Short:                      "Task marketplace transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	taskmarketTxCmd.AddCommand(CmdCreateTask())
	taskmarketTxCmd.AddCommand(CmdUpdateTask())
	taskmarketTxCmd.AddCommand(CmdPublishTask())
	taskmarketTxCmd.AddCommand(CmdCancelTask())
	taskmarketTxCmd.AddCommand(CmdSubmitApplication())
	taskmarketTxCmd.AddCommand(CmdAcceptApplication())
	taskmarketTxCmd.AddCommand(CmdRejectApplication())
	taskmarketTxCmd.AddCommand(CmdSubmitBid())
	taskmarketTxCmd.AddCommand(CmdCloseAuction())
	taskmarketTxCmd.AddCommand(CmdStartTask())
	taskmarketTxCmd.AddCommand(CmdSubmitMilestone())
	taskmarketTxCmd.AddCommand(CmdApproveMilestone())
	taskmarketTxCmd.AddCommand(CmdRejectMilestone())
	taskmarketTxCmd.AddCommand(CmdSubmitRating())

	return taskmarketTxCmd
}

// CmdCreateTask implements the create-task command
func CmdCreateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-task [title] [description] [type] [budget]",
		Short: "Create a new marketplace task",
		Long: `Create a new task in the marketplace.
Types: open (for applications) or auction (for bidding).`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			description := args[1]
			taskType := args[2]
			budget, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid budget: %w", err)
			}

			// Get flags
			category, _ := cmd.Flags().GetString("category")
			skills, _ := cmd.Flags().GetStringSlice("skills")
			deadline, _ := cmd.Flags().GetInt64("deadline")

			fmt.Printf("Creating task:\n")
			fmt.Printf("  Title: %s\n", title)
			fmt.Printf("  Description: %s\n", description)
			fmt.Printf("  Type: %s\n", taskType)
			fmt.Printf("  Budget: %d STT\n", budget)
			fmt.Printf("  Category: %s\n", category)
			fmt.Printf("  Skills: %s\n", strings.Join(skills, ", "))
			fmt.Printf("  Deadline: %d\n", deadline)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().String("category", "other", "Task category (development, design, writing, etc.)")
	cmd.Flags().StringSlice("skills", []string{}, "Required skills (comma-separated)")
	cmd.Flags().Int64("deadline", 0, "Task deadline (Unix timestamp)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdateTask implements the update-task command
func CmdUpdateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-task [task-id]",
		Short: "Update an existing task",
		Long:  `Update a task that is still in draft or open status.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Updating task: %s\n", taskID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().String("title", "", "New title")
	cmd.Flags().String("description", "", "New description")
	cmd.Flags().Uint64("budget", 0, "New budget")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdPublishTask implements the publish-task command
func CmdPublishTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish-task [task-id]",
		Short: "Publish a draft task to the marketplace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Publishing task: %s\n", taskID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCancelTask implements the cancel-task command
func CmdCancelTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-task [task-id]",
		Short: "Cancel a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Cancelling task: %s\n", taskID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitApplication implements the submit-application command
func CmdSubmitApplication() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-application [task-id] [price]",
		Short: "Submit an application for an open task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			price, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid price: %w", err)
			}

			coverLetter, _ := cmd.Flags().GetString("cover-letter")
			duration, _ := cmd.Flags().GetInt64("duration")

			fmt.Printf("Submitting application for task: %s\n", taskID)
			fmt.Printf("  Proposed Price: %d STT\n", price)
			fmt.Printf("  Cover Letter: %s\n", coverLetter)
			fmt.Printf("  Estimated Duration: %d days\n", duration)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().String("cover-letter", "", "Cover letter/message")
	cmd.Flags().Int64("duration", 0, "Estimated duration in days")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdAcceptApplication implements the accept-application command
func CmdAcceptApplication() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-application [application-id]",
		Short: "Accept an application for your task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			applicationID := args[0]

			fmt.Printf("Accepting application: %s\n", applicationID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRejectApplication implements the reject-application command
func CmdRejectApplication() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-application [application-id]",
		Short: "Reject an application for your task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			applicationID := args[0]

			fmt.Printf("Rejecting application: %s\n", applicationID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitBid implements the submit-bid command
func CmdSubmitBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-bid [task-id] [amount]",
		Short: "Submit a bid for an auction task",
		Long:  `Submit a bid amount (lower is better) for an auction-based task.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			amount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			message, _ := cmd.Flags().GetString("message")

			fmt.Printf("Submitting bid for task: %s\n", taskID)
			fmt.Printf("  Amount: %d STT\n", amount)
			fmt.Printf("  Message: %s\n", message)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().String("message", "", "Bid message")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdCloseAuction implements the close-auction command
func CmdCloseAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-auction [task-id]",
		Short: "Close an auction and select the winner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Closing auction for task: %s\n", taskID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdStartTask implements the start-task command
func CmdStartTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-task [task-id]",
		Short: "Start working on an assigned task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			fmt.Printf("Starting task: %s\n", taskID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitMilestone implements the submit-milestone command
func CmdSubmitMilestone() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-milestone [task-id] [milestone-id] [deliverables]",
		Short: "Submit deliverables for a milestone",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			milestoneID := args[1]
			deliverables := args[2]

			fmt.Printf("Submitting milestone for task: %s\n", taskID)
			fmt.Printf("  Milestone ID: %s\n", milestoneID)
			fmt.Printf("  Deliverables: %s\n", deliverables)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdApproveMilestone implements the approve-milestone command
func CmdApproveMilestone() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-milestone [task-id] [milestone-id]",
		Short: "Approve a submitted milestone",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			milestoneID := args[1]

			fmt.Printf("Approving milestone for task: %s\n", taskID)
			fmt.Printf("  Milestone ID: %s\n", milestoneID)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRejectMilestone implements the reject-milestone command
func CmdRejectMilestone() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-milestone [task-id] [milestone-id] [reason]",
		Short: "Reject a submitted milestone",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			milestoneID := args[1]
			reason := args[2]

			fmt.Printf("Rejecting milestone for task: %s\n", taskID)
			fmt.Printf("  Milestone ID: %s\n", milestoneID)
			fmt.Printf("  Reason: %s\n", reason)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitRating implements the submit-rating command
func CmdSubmitRating() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-rating [task-id] [rated-id]",
		Short: "Submit a rating for a completed task",
		Long:  `Submit a multi-dimensional rating (1-5) for the other party after task completion.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			ratedID := args[1]

			quality, _ := cmd.Flags().GetInt("quality")
			communication, _ := cmd.Flags().GetInt("communication")
			timeliness, _ := cmd.Flags().GetInt("timeliness")
			professionalism, _ := cmd.Flags().GetInt("professionalism")
			comment, _ := cmd.Flags().GetString("comment")

			fmt.Printf("Submitting rating for task: %s\n", taskID)
			fmt.Printf("  Rated User: %s\n", ratedID)
			fmt.Printf("  Quality: %d\n", quality)
			fmt.Printf("  Communication: %d\n", communication)
			fmt.Printf("  Timeliness: %d\n", timeliness)
			fmt.Printf("  Professionalism: %d\n", professionalism)
			fmt.Printf("  Comment: %s\n", comment)
			fmt.Println("\nNote: Full implementation requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().Int("quality", 0, "Quality rating (1-5)")
	cmd.Flags().Int("communication", 0, "Communication rating (1-5)")
	cmd.Flags().Int("timeliness", 0, "Timeliness rating (1-5)")
	cmd.Flags().Int("professionalism", 0, "Professionalism rating (1-5)")
	cmd.Flags().String("comment", "", "Rating comment")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
