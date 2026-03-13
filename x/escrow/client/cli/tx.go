package cli

import (
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "escrow",
		Short:                      "Escrow transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       nil,
	}

	return cmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "escrow",
		Short:                      "Querying commands for the escrow module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       nil,
	}

	return cmd
}
