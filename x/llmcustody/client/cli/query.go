package cli

import (
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

	return llmcustodyQueryCmd
}

// CmdQueryAPIKey implements the query api-key command
func CmdQueryAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-key [id]",
		Short: "Query an API key by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			// For now, query from local keeper directly
			// In full implementation, this would use gRPC
			fmt.Printf("Querying API key: %s\n", id)
			fmt.Println("Note: Full query implementation requires proto-generated query client")

			return nil
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
			owner := args[0]

			fmt.Printf("Querying API keys for owner: %s\n", owner)
			fmt.Println("Note: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryAllAPIKeys implements the query all-api-keys command
func CmdQueryAllAPIKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-api-keys",
		Short: "Query all registered API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Querying all API keys")
			fmt.Println("Note: Full query implementation requires proto-generated query client")

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
