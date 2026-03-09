package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"sharetoken/x/identity/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryIdentity())
	cmd.AddCommand(CmdQueryIdentities())
	cmd.AddCommand(CmdQueryLimitConfig())
	cmd.AddCommand(CmdQueryIsVerified())
	cmd.AddCommand(CmdQueryParams())

	return cmd
}

func CmdQueryIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity [address]",
		Short: "Query an identity by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// For now, return a placeholder
			// In production, this would query the chain
			identity := types.Identity{
				Address: args[0],
			}

			data, err := json.MarshalIndent(identity, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryIdentities() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identities",
		Short: "Query all identities",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// For now, return a placeholder
			// In production, this would query the chain
			identities := []types.Identity{}

			data, err := json.MarshalIndent(identities, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryLimitConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limit-config [address]",
		Short: "Query a user's limit configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// For now, return a placeholder
			// In production, this would query the chain
			limitConfig := types.NewLimitConfig(args[0])

			data, err := json.MarshalIndent(limitConfig, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryIsVerified() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-verified [address]",
		Short: "Check if an address is verified",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// For now, return a placeholder
			// In production, this would query the chain
			result := map[string]bool{
				"is_verified": false,
			}

			data, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query module parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// For now, return a placeholder
			// In production, this would query the chain
			params := types.DefaultParams()

			data, err := json.MarshalIndent(params, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
