package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"sharetoken/x/marketplace/types"
)

func TestService_NewService(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	price := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))

	service := types.NewService("svc-1", validAddress, "Test Service", types.ServiceLevelLLM, price)

	require.NotNil(t, service)
	require.Equal(t, "svc-1", service.ID)
	require.Equal(t, validAddress, service.Provider)
	require.Equal(t, "Test Service", service.Name)
	require.Equal(t, types.ServiceLevelLLM, service.Level)
	require.Equal(t, types.PricingModeFixed, service.PricingMode)
	require.True(t, service.Price.IsEqual(price))
	require.True(t, service.Active)
}

func TestPricingModeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.PricingMode
	}{
		{"dynamic", "dynamic", types.PricingModeDynamic},
		{"auction", "auction", types.PricingModeAuction},
		{"fixed", "fixed", types.PricingModeFixed},
		{"empty", "", types.PricingModeFixed},
		{"unknown", "unknown", types.PricingModeFixed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.PricingModeFromString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestService_String(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	price := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))
	service := types.NewService("svc-1", validAddress, "Test Service", types.ServiceLevelLLM, price)

	result := service.String()
	require.Contains(t, result, "svc-1")
	require.Contains(t, result, "Test Service")
	require.Contains(t, result, "1000")
}

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Services)
	require.Empty(t, genesis.Services)
}

func TestValidateGenesis(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validAddress2 := sdk.AccAddress([]byte("test_address_2")).String()
	price := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))

	tests := []struct {
		name    string
		data    types.GenesisState
		wantErr bool
	}{
		{
			name:    "valid genesis with default",
			data:    *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "valid genesis with services",
			data: types.GenesisState{
				Services: []types.Service{
					{ID: "svc-1", Provider: validAddress, Name: "Service 1", Level: types.ServiceLevelLLM, Price: price, Active: true},
					{ID: "svc-2", Provider: validAddress2, Name: "Service 2", Level: types.ServiceLevelAgent, Price: price, Active: true},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate service IDs",
			data: types.GenesisState{
				Services: []types.Service{
					{ID: "svc-1", Provider: validAddress, Name: "Service 1", Level: types.ServiceLevelLLM, Price: price, Active: true},
					{ID: "svc-1", Provider: validAddress2, Name: "Service 2", Level: types.ServiceLevelAgent, Price: price, Active: true},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty service ID",
			data: types.GenesisState{
				Services: []types.Service{
					{ID: "", Provider: validAddress, Name: "Service 1", Level: types.ServiceLevelLLM, Price: price, Active: true},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty provider",
			data: types.GenesisState{
				Services: []types.Service{
					{ID: "svc-1", Provider: "", Name: "Service 1", Level: types.ServiceLevelLLM, Price: price, Active: true},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty name",
			data: types.GenesisState{
				Services: []types.Service{
					{ID: "svc-1", Provider: validAddress, Name: "", Level: types.ServiceLevelLLM, Price: price, Active: true},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateGenesis(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Service Marshal/Unmarshal Tests

func TestService_MarshalUnmarshal(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	price := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))
	original := types.NewService("svc-1", validAddress, "Test Service", types.ServiceLevelLLM, price)
	original.Description = "Test Description"

	// Marshal
	data, err := original.Marshal()
	require.NoError(t, err)
	require.NotNil(t, data)
	require.True(t, len(data) > 0)

	// Unmarshal
	var restored types.Service
	err = restored.Unmarshal(data)
	require.NoError(t, err)

	require.Equal(t, original.ID, restored.ID)
	require.Equal(t, original.Provider, restored.Provider)
	require.Equal(t, original.Name, restored.Name)
	require.Equal(t, original.Description, restored.Description)
	require.Equal(t, original.Level, restored.Level)
	require.Equal(t, original.PricingMode, restored.PricingMode)
	require.True(t, original.Price.IsEqual(restored.Price))
	require.Equal(t, original.Active, restored.Active)
	require.Equal(t, original.CreatedAt, restored.CreatedAt)
}

func TestService_Size(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	price := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))
	service := types.NewService("svc-1", validAddress, "Test Service", types.ServiceLevelLLM, price)

	size := service.Size()
	data, _ := service.Marshal()
	require.Equal(t, len(data), size)
}

// Service Level Tests

func TestServiceLevel_Values(t *testing.T) {
	require.Equal(t, types.ServiceLevel(0), types.ServiceLevelUnspecified)
	require.Equal(t, types.ServiceLevel(1), types.ServiceLevelLLM)
	require.Equal(t, types.ServiceLevel(2), types.ServiceLevelAgent)
	require.Equal(t, types.ServiceLevel(3), types.ServiceLevelWorkflow)
}

// Pricing Mode Tests

func TestPricingMode_Values(t *testing.T) {
	require.Equal(t, types.PricingMode(0), types.PricingModeUnspecified)
	require.Equal(t, types.PricingMode(1), types.PricingModeFixed)
	require.Equal(t, types.PricingMode(2), types.PricingModeDynamic)
	require.Equal(t, types.PricingMode(3), types.PricingModeAuction)
}
