package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"sharetoken/x/identity/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdRegisterIdentity())
	cmd.AddCommand(CmdVerifyIdentity())

	return cmd
}

func CmdRegisterIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-identity [address] [did] [metadata-hash]",
		Short: "Register a new identity",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Registering identity for address: %s\n", args[0])
			fmt.Printf("DID: %s\n", args[1])
			fmt.Printf("Metadata Hash: %s\n", args[2])
			fmt.Println("Note: This is a placeholder command. Integration with chain pending.")
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdVerifyIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-identity [address] [provider] [verification-hash] [proof]",
		Short: "Verify an identity with a third-party provider",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Verifying identity for address: %s\n", args[0])
			fmt.Printf("Provider: %s\n", args[1])
			fmt.Printf("Verification Hash: %s\n", args[2])
			fmt.Println("Note: This is a placeholder command. Integration with chain pending.")
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
