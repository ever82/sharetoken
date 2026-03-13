package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/marketplace/types"
)

func TestNewService(t *testing.T) {
	service := types.NewService("svc-1", "provider1", "GPT-4 API", types.ServiceLevelLLM, sdk.NewCoins(sdk.NewInt64Coin("ustt", 1000)))

	require.Equal(t, "svc-1", service.ID)
	require.Equal(t, "provider1", service.Provider)
	require.Equal(t, "GPT-4 API", service.Name)
	require.Equal(t, types.ServiceLevelLLM, service.Level)
	require.True(t, service.Active)
}

func TestPricingModes(t *testing.T) {
	require.Equal(t, types.PricingModeFixed, types.PricingModeFromString("fixed"))
	require.Equal(t, types.PricingModeDynamic, types.PricingModeFromString("dynamic"))
	require.Equal(t, types.PricingModeAuction, types.PricingModeFromString("auction"))
}
