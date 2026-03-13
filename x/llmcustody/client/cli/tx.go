package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"sharetoken/x/llmcustody/types"
)

// CmdRegisterAPIKey implements the register-api-key command
func CmdRegisterAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-api-key [provider] [encrypted-key-file]",
		Short: "Register a new API key for LLM provider",
		Long: `Register a new encrypted API key for an LLM provider (openai or anthropic).
The encrypted key should be read from a file containing the API key encrypted with the platform's public key.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			provider := args[0]
			encryptedKeyFile := args[1]

			// Read encrypted key from file
			encryptedKey, err := os.ReadFile(encryptedKeyFile)
			if err != nil {
				return fmt.Errorf("failed to read encrypted key file: %w", err)
			}

			// Parse access rules from flags
			accessRules, err := parseAccessRules(cmd)
			if err != nil {
				return err
			}

			// Create message
			msg := types.NewMsgRegisterAPIKey(
				clientCtx.GetFromAddress().String(),
				provider,
				encryptedKey,
				accessRules,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Output the message as JSON for now (until proto is fully integrated)
			output, _ := json.MarshalIndent(msg, "", "  ")
			fmt.Println("Message created successfully:")
			fmt.Println(string(output))
			fmt.Println("\nNote: Transaction broadcasting requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().StringSlice("access-rules", []string{}, "Access rules in format: serviceID:rateLimit:maxRequests:pricePerReq")
	flags.AddTxFlagsToCmd(cmd)

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
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			apiKeyID := args[0]

			// Get flags
			active, err := cmd.Flags().GetBool("active")
			if err != nil {
				return err
			}

			// Parse access rules from flags
			accessRules, err := parseAccessRules(cmd)
			if err != nil {
				return err
			}

			// Create message
			msg := types.NewMsgUpdateAPIKey(
				clientCtx.GetFromAddress().String(),
				apiKeyID,
				accessRules,
				active,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Output the message as JSON for now
			output, _ := json.MarshalIndent(msg, "", "  ")
			fmt.Println("Message created successfully:")
			fmt.Println(string(output))
			fmt.Println("\nNote: Transaction broadcasting requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().Bool("active", true, "Set the API key active status")
	cmd.Flags().StringSlice("access-rules", []string{}, "Access rules in format: serviceID:rateLimit:maxRequests:pricePerReq")
	flags.AddTxFlagsToCmd(cmd)

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
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			apiKeyID := args[0]

			// Create message
			msg := types.NewMsgRevokeAPIKey(
				clientCtx.GetFromAddress().String(),
				apiKeyID,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Output the message as JSON for now
			output, _ := json.MarshalIndent(msg, "", "  ")
			fmt.Println("Message created successfully:")
			fmt.Println(string(output))
			fmt.Println("\nNote: Transaction broadcasting requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdRotateAPIKey implements the rotate-api-key command
func CmdRotateAPIKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate-api-key [api-key-id] [new-encrypted-key-file]",
		Short: "Rotate an API key",
		Long: `Rotate an API key with a new encrypted key. The old key will be deactivated
and a new key ID will be generated. Usage statistics will be preserved.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			apiKeyID := args[0]
			newEncryptedKeyFile := args[1]

			// Read new encrypted key from file
			newEncryptedKey, err := os.ReadFile(newEncryptedKeyFile)
			if err != nil {
				return fmt.Errorf("failed to read encrypted key file: %w", err)
			}

			// Get reason flag
			reason, _ := cmd.Flags().GetString("reason")

			// Create message
			msg := types.NewMsgRotateAPIKey(
				clientCtx.GetFromAddress().String(),
				apiKeyID,
				newEncryptedKey,
				reason,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Output the message as JSON for now
			output, _ := json.MarshalIndent(msg, "", "  ")
			fmt.Println("Message created successfully:")
			fmt.Println(string(output))
			fmt.Println("\nNote: Transaction broadcasting requires proto-generated message types")

			return nil
		},
	}

	cmd.Flags().String("reason", "", "Reason for key rotation")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdRecordUsage implements the record-usage command (typically called by oracle/service)
func CmdRecordUsage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-usage [api-key-id] [service-id] [request-count] [token-count] [cost]",
		Short: "Record API usage",
		Long:  `Record API usage for billing and tracking purposes.`,
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			apiKeyID := args[0]
			serviceID := args[1]

			// Parse numeric arguments
			requestCount, err := parseInt64(args[2])
			if err != nil {
				return fmt.Errorf("invalid request count: %w", err)
			}

			tokenCount, err := parseInt64(args[3])
			if err != nil {
				return fmt.Errorf("invalid token count: %w", err)
			}

			cost, err := parseInt64(args[4])
			if err != nil {
				return fmt.Errorf("invalid cost: %w", err)
			}

			// Create message
			msg := types.NewMsgRecordUsage(
				apiKeyID,
				serviceID,
				requestCount,
				tokenCount,
				cost,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Output the message as JSON for now
			output, _ := json.MarshalIndent(msg, "", "  ")
			fmt.Println("Message created successfully:")
			fmt.Println(string(output))
			fmt.Println("\nNote: Transaction broadcasting requires proto-generated message types")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Helper functions

func parseAccessRules(cmd *cobra.Command) ([]types.AccessRule, error) {
	rulesStr, err := cmd.Flags().GetStringSlice("access-rules")
	if err != nil {
		return nil, err
	}

	if len(rulesStr) == 0 {
		return []types.AccessRule{}, nil
	}

	rules := make([]types.AccessRule, len(rulesStr))
	for i, ruleStr := range rulesStr {
		// Parse format: serviceID:rateLimit:maxRequests:pricePerReq
		parts := splitRule(ruleStr)
		if len(parts) < 1 {
			return nil, fmt.Errorf("invalid access rule format: %s", ruleStr)
		}

		rule := types.AccessRule{
			ServiceID: parts[0],
			Allowed:   true,
		}

		if len(parts) > 1 {
			rateLimit, err := parseInt64(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid rate limit: %w", err)
			}
			rule.RateLimit = rateLimit
		}

		if len(parts) > 2 {
			maxRequests, err := parseInt64(parts[2])
			if err != nil {
				return nil, fmt.Errorf("invalid max requests: %w", err)
			}
			rule.MaxRequests = maxRequests
		}

		if len(parts) > 3 {
			pricePerReq, err := parseInt64(parts[3])
			if err != nil {
				return nil, fmt.Errorf("invalid price per request: %w", err)
			}
			rule.PricePerReq = pricePerReq
		}

		rules[i] = rule
	}

	return rules, nil
}

func splitRule(rule string) []string {
	// Simple split by colon
	parts := []string{}
	start := 0
	for i := 0; i < len(rule); i++ {
		if rule[i] == ':' {
			parts = append(parts, rule[start:i])
			start = i + 1
		}
	}
	parts = append(parts, rule[start:])
	return parts
}

func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// EncryptAPIKeyCmd is a utility command to encrypt an API key before registration
func EncryptAPIKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt-api-key [api-key]",
		Short: "Encrypt an API key for secure storage",
		Long: `Encrypt an API key using the platform's encryption key.
This command is for testing and demonstration purposes only.
In production, API keys should be encrypted client-side.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := args[0]

			// Create a new encryption key
			encKey, err := types.NewEncryptionKey()
			if err != nil {
				return fmt.Errorf("failed to create encryption key: %w", err)
			}

			// Encrypt the API key
			encrypted, err := encKey.Encrypt([]byte(apiKey))
			if err != nil {
				return fmt.Errorf("failed to encrypt API key: %w", err)
			}

			// Output the encrypted key (base64 encoded)
			fmt.Println("Encrypted API Key (base64):")
			fmt.Println(base64.StdEncoding.EncodeToString(encrypted))

			// Securely wipe the key
			apiKey = ""
			types.Zeroize(encKey.Key)

			return nil
		},
	}

	return cmd
}
