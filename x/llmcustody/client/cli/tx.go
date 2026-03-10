package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CmdRegisterAPIKey implements the register-api-key command
func CmdRegisterAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-api-key [provider] [encrypted-key]",
		Short: "Register a new API key for LLM provider",
		Long: `Register a new encrypted API key for an LLM provider (openai or anthropic).
The encrypted key should be the API key encrypted with the platform's public key.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := args[0]
			encryptedKey := args[1]

			fmt.Printf("Registering API key for provider: %s\n", provider)
			fmt.Printf("Encrypted key: %s...\n", encryptedKey[:min(20, len(encryptedKey))])
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	return cmd
}

// CmdUpdateAPIKey implements the update-api-key command
func CmdUpdateAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-api-key [api-key-id]",
		Short: "Update an existing API key",
		Long:  `Update an API key's access rules or active status.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKeyID := args[0]

			// Get flags
			active, _ := cmd.Flags().GetBool("active")

			fmt.Printf("Updating API key: %s\n", apiKeyID)
			fmt.Printf("Active status: %v\n", active)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().Bool("active", true, "Set the API key active status")

	return cmd
}

// CmdRevokeAPIKey implements the revoke-api-key command
func CmdRevokeAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-api-key [api-key-id]",
		Short: "Revoke (delete) an API key",
		Long:  `Permanently revoke and delete an API key. This action cannot be undone.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKeyID := args[0]

			fmt.Printf("Revoking API key: %s\n", apiKeyID)
			fmt.Println("Note: Full implementation requires proto-generated message types")

			return nil
		},
	}

	return cmd
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
