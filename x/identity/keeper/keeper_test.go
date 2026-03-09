package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/identity/types"
)

func TestNewLimitConfig(t *testing.T) {
	address := "sharetoken1xyz"
	limitConfig := types.NewLimitConfig(address)

	require.Equal(t, address, limitConfig.Address)
	require.NotNil(t, limitConfig.TxLimit)
	require.NotNil(t, limitConfig.WithdrawalLimit)
	require.NotNil(t, limitConfig.DisputeLimit)
	require.NotNil(t, limitConfig.ServiceLimit)
}

func TestTransactionLimit(t *testing.T) {
	address := "sharetoken1xyz"
	limitConfig := types.NewLimitConfig(address)

	// Test daily limit check
	// This is a simplified test - in production, you would need to set up proper test fixtures
	err := limitConfig.CheckTransactionLimit(types.DefaultCoin())
	require.NoError(t, err)
}

func TestWithdrawalLimit(t *testing.T) {
	address := "sharetoken1xyz"
	limitConfig := types.NewLimitConfig(address)

	// Test withdrawal limit check
	err := limitConfig.CheckWithdrawalLimit(types.DefaultCoin())
	require.NoError(t, err)
}

func TestDisputeLimit(t *testing.T) {
	address := "sharetoken1xyz"
	limitConfig := types.NewLimitConfig(address)

	// Test dispute limit check - should pass when no active disputes
	err := limitConfig.CheckDisputeLimit()
	require.NoError(t, err)

	// Increment active disputes to max
	for i := uint64(0); i < limitConfig.DisputeLimit.MaxActiveDisputes; i++ {
		limitConfig.IncrementActiveDisputes()
	}

	// Should fail when at max
	err = limitConfig.CheckDisputeLimit()
	require.Error(t, err)
}

func TestServiceLimit(t *testing.T) {
	address := "sharetoken1xyz"
	limitConfig := types.NewLimitConfig(address)

	// Test service limit check - should pass when no concurrent services
	err := limitConfig.CheckServiceLimit()
	require.NoError(t, err)

	// Record multiple service calls
	for i := uint64(0); i < limitConfig.ServiceLimit.MaxConcurrent; i++ {
		limitConfig.RecordServiceCall()
	}

	// Should fail when at max concurrent
	err = limitConfig.CheckServiceLimit()
	require.Error(t, err)
}

func TestIdentityValidation(t *testing.T) {
	identity := types.Identity{
		Address: "sharetoken1xyz",
		DID:     "did:sharetoken:abc123",
	}

	// Test basic validation - this will fail because address is not a valid bech32 address
	// In production, you would use a valid address
	err := identity.ValidateBasic()
	require.Error(t, err) // Expected to fail with invalid address
}

func TestIsValidProvider(t *testing.T) {
	require.True(t, types.IsValidProvider("wechat"))
	require.True(t, types.IsValidProvider("github"))
	require.True(t, types.IsValidProvider("google"))
	require.False(t, types.IsValidProvider("invalid"))
}

func TestDefaultParams(t *testing.T) {
	params := types.DefaultParams()
	require.False(t, params.VerificationRequired)
	require.NotEmpty(t, params.AllowedProviders)
	require.Contains(t, params.AllowedProviders, "wechat")
	require.Contains(t, params.AllowedProviders, "github")
	require.Contains(t, params.AllowedProviders, "google")
}

func TestGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()
	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Params)
	require.Empty(t, genesis.Identities)
	require.Empty(t, genesis.LimitConfigs)

	// Test validation
	err := types.ValidateGenesis(*genesis)
	require.NoError(t, err)
}
