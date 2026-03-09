package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/llmcustody/types"
)

func TestAPIKeyCreation(t *testing.T) {
	// Generate encryption key
	kek, err := types.NewEncryptionKey()
	require.NoError(t, err)

	// Encrypt an API key
	plainKey := "sk-test123456789"
	encrypted, err := kek.Encrypt([]byte(plainKey))
	require.NoError(t, err)

	// Create API key record
	apiKey := types.NewAPIKey("key-1", types.ProviderOpenAI, encrypted, "owner1")
	require.Equal(t, "key-1", apiKey.ID)
	require.Equal(t, types.ProviderOpenAI, apiKey.Provider)
	require.Equal(t, "owner1", apiKey.Owner)
	require.True(t, apiKey.Active)
	require.NotEmpty(t, apiKey.Hash)

	// Verify hash
	require.True(t, apiKey.VerifyHash(encrypted))

	// Decrypt and verify
	decrypted, err := kek.Decrypt(apiKey.EncryptedKey)
	require.NoError(t, err)
	require.Equal(t, plainKey, string(decrypted))

	// Zeroize for security
	types.Zeroize(decrypted)
}

func TestAPIKeyValidation(t *testing.T) {
	kek, _ := types.NewEncryptionKey()
	encrypted, _ := kek.Encrypt([]byte("test"))

	// Valid key
	apiKey := types.NewAPIKey("key-1", types.ProviderOpenAI, encrypted, "owner1")
	require.NoError(t, apiKey.ValidateBasic())

	// Invalid - missing ID
	invalidKey := types.APIKey{
		ID:           "",
		Provider:     types.ProviderOpenAI,
		EncryptedKey: encrypted,
		Owner:        "owner1",
	}
	require.Error(t, invalidKey.ValidateBasic())

	// Invalid - empty encrypted key
	invalidKey2 := types.NewAPIKey("key-2", types.ProviderOpenAI, []byte{}, "owner1")
	require.Error(t, invalidKey2.ValidateBasic())

	// Invalid - invalid provider
	invalidKey3 := types.NewAPIKey("key-3", types.Provider("invalid"), encrypted, "owner1")
	require.Error(t, invalidKey3.ValidateBasic())
}

func TestAccessRules(t *testing.T) {
	kek, _ := types.NewEncryptionKey()
	encrypted, _ := kek.Encrypt([]byte("test"))

	apiKey := types.NewAPIKey("key-1", types.ProviderOpenAI, encrypted, "owner1")

	// Add access rule
	apiKey.AccessRules = []types.AccessRule{
		{
			ServiceID:   "service-1",
			RateLimit:   100,
			MaxRequests: 1000,
			PricePerReq: 100,
		},
	}

	// Should be able to access
	require.True(t, apiKey.CanAccess("service-1"))

	// Should not be able to access unknown service
	require.False(t, apiKey.CanAccess("service-2"))

	// Should not be able to access when max requests reached
	apiKey.UsageCount = 1000
	require.False(t, apiKey.CanAccess("service-1"))

	// Should not be able to access when inactive
	apiKey.Active = false
	require.False(t, apiKey.CanAccess("service-1"))
}

func TestProviderValidation(t *testing.T) {
	require.True(t, types.IsValidProvider("openai"))
	require.True(t, types.IsValidProvider("anthropic"))
	require.False(t, types.IsValidProvider("invalid"))
	require.False(t, types.IsValidProvider(""))
}

func TestSecureString(t *testing.T) {
	sensitive := "my-secret-api-key"
	ss := types.NewSecureString(sensitive)

	require.Equal(t, sensitive, ss.String())
	require.Equal(t, []byte(sensitive), ss.Bytes())

	// Zeroize
	ss.Zeroize()
	require.True(t, ss.IsZeroized())
}

func TestEncryption(t *testing.T) {
	kek, err := types.NewEncryptionKey()
	require.NoError(t, err)

	plaintext := []byte("secret-api-key-12345")

	// Encrypt
	ciphertext, err := kek.Encrypt(plaintext)
	require.NoError(t, err)
	require.NotEqual(t, plaintext, ciphertext)

	// Decrypt
	decrypted, err := kek.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)

	// Zeroize
	types.Zeroize(decrypted)
}

func TestHashKey(t *testing.T) {
	apiKey := "sk-abc123"
	hash1 := types.HashKey(apiKey)
	hash2 := types.HashKey(apiKey)

	// Same input should produce same hash
	require.Equal(t, hash1, hash2)

	// Different input should produce different hash
	hash3 := types.HashKey("sk-different")
	require.NotEqual(t, hash1, hash3)
}

func TestAPIKeyWipe(t *testing.T) {
	kek, _ := types.NewEncryptionKey()
	encrypted, _ := kek.Encrypt([]byte("test"))

	apiKey := types.NewAPIKey("key-1", types.ProviderOpenAI, encrypted, "owner1")

	// Verify key exists
	require.NotEmpty(t, apiKey.EncryptedKey)

	// Wipe
	apiKey.SecureWipe()

	// Key should be zeroed
	require.Empty(t, apiKey.EncryptedKey)
}
