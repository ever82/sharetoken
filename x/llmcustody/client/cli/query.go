package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	llmcustodyTxCmd := &cobra.Command{
		Use:                        "llmcustody",
		Short:                      "LLM Custody transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	llmcustodyTxCmd.AddCommand(CmdRegisterAPIKey())
	llmcustodyTxCmd.AddCommand(CmdUpdateAPIKey())
	llmcustodyTxCmd.AddCommand(CmdRevokeAPIKey())
	llmcustodyTxCmd.AddCommand(CmdRotateAPIKey())
	llmcustodyTxCmd.AddCommand(CmdRecordUsage())
	llmcustodyTxCmd.AddCommand(EncryptAPIKeyCmd())

	return llmcustodyTxCmd
}

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	llmcustodyQueryCmd := &cobra.Command{
		Use:                        "llmcustody",
		Short:                      "Querying commands for the llmcustody module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	llmcustodyQueryCmd.AddCommand(CmdQueryAPIKey())
	llmcustodyQueryCmd.AddCommand(CmdQueryAPIKeysByOwner())
	llmcustodyQueryCmd.AddCommand(CmdQueryAllAPIKeys())
	llmcustodyQueryCmd.AddCommand(CmdQueryUsageStats())
	llmcustodyQueryCmd.AddCommand(CmdQueryDailyUsage())
	llmcustodyQueryCmd.AddCommand(CmdQueryServiceUsage())

	return llmcustodyQueryCmd
}

// CmdQueryAPIKey implements the query api-key command
func CmdQueryAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-key [id]",
		Short: "Query an API key by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id := args[0]

			// Create query client
			queryClient := NewQueryClient(clientCtx)

			res, err := queryClient.APIKey(context.Background(), id)
			if err != nil {
				return err
			}

			return printOutput(clientCtx, res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryAPIKeysByOwner implements the query api-keys-by-owner command
func CmdQueryAPIKeysByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-keys-by-owner [owner-address]",
		Short: "Query all API keys owned by an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner := args[0]

			// Get pagination flags
			page, _ := cmd.Flags().GetInt(flags.FlagPage)
			limit, _ := cmd.Flags().GetInt(flags.FlagLimit)

			// Create query client
			queryClient := NewQueryClient(clientCtx)

			res, err := queryClient.APIKeysByOwner(context.Background(), owner, page, limit)
			if err != nil {
				return err
			}

			return printOutput(clientCtx, res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "api-keys")

	return cmd
}

// CmdQueryAllAPIKeys implements the query all-api-keys command
func CmdQueryAllAPIKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-api-keys",
		Short: "Query all registered API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Get pagination flags
			page, _ := cmd.Flags().GetInt(flags.FlagPage)
			limit, _ := cmd.Flags().GetInt(flags.FlagLimit)

			// Create query client
			queryClient := NewQueryClient(clientCtx)

			res, err := queryClient.AllAPIKeys(context.Background(), page, limit)
			if err != nil {
				return err
			}

			return printOutput(clientCtx, res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "api-keys")

	return cmd
}

// CmdQueryUsageStats implements the query usage-stats command
func CmdQueryUsageStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "usage-stats [api-key-id]",
		Short: "Query usage statistics for an API key",
		Long:  `Query aggregated usage statistics including total requests, tokens, and cost for an API key.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			apiKeyID := args[0]

			// Create query client
			queryClient := NewQueryClient(clientCtx)

			res, err := queryClient.UsageStats(context.Background(), apiKeyID)
			if err != nil {
				return err
			}

			return printOutput(clientCtx, res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryDailyUsage implements the query daily-usage command
func CmdQueryDailyUsage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daily-usage [api-key-id] [date]",
		Short: "Query daily usage statistics",
		Long: `Query daily usage statistics for an API key.
Date format: YYYY-MM-DD`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			apiKeyID := args[0]
			date := args[1]

			// Create query client
			queryClient := NewQueryClient(clientCtx)

			res, err := queryClient.DailyUsage(context.Background(), apiKeyID, date)
			if err != nil {
				return err
			}

			return printOutput(clientCtx, res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryServiceUsage implements the query service-usage command
func CmdQueryServiceUsage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-usage [service-id] [api-key-id]",
		Short: "Query service usage statistics",
		Long: `Query usage statistics for a specific service.
api-key-id is optional - if not provided, returns stats for all keys.`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serviceID := args[0]

			apiKeyID := ""
			if len(args) > 1 {
				apiKeyID = args[1]
			}

			// Create query client
			queryClient := NewQueryClient(clientCtx)

			res, err := queryClient.ServiceUsage(context.Background(), serviceID, apiKeyID)
			if err != nil {
				return err
			}

			return printOutput(clientCtx, res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// printOutput prints the output in JSON format
func printOutput(clientCtx client.Context, v interface{}) error {
	output, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
