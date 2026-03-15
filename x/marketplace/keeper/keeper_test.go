package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/marketplace/types"
)

func TestNewService(t *testing.T) {
	service := types.NewService("svc-1", "provider1", "GPT-4 API", types.ServiceLevel_SERVICE_LEVEL_LLM, sdk.NewCoins(sdk.NewInt64Coin("ustt", 1000)))

	require.Equal(t, "svc-1", service.Id)
	require.Equal(t, "provider1", service.Provider)
	require.Equal(t, "GPT-4 API", service.Name)
	require.Equal(t, types.ServiceLevel_SERVICE_LEVEL_LLM, service.Level)
	require.True(t, service.Active)
}

func TestPricingModes(t *testing.T) {
	require.Equal(t, types.PricingMode_PRICING_MODE_FIXED, types.PricingModeFromString("fixed"))
	require.Equal(t, types.PricingMode_PRICING_MODE_DYNAMIC, types.PricingModeFromString("dynamic"))
	require.Equal(t, types.PricingMode_PRICING_MODE_AUCTION, types.PricingModeFromString("auction"))
}
