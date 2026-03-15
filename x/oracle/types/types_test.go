package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"sharetoken/x/oracle/types"
)

func TestNewPrice(t *testing.T) {
	price := types.NewPrice("STT", sdk.NewDec(10), types.PriceSourceChainlink, 95)

	require.NotNil(t, price)
	require.Equal(t, "STT", price.Symbol)
	require.True(t, price.Price.Equal(sdk.NewDec(10)))
	require.Equal(t, types.PriceSourceChainlink, price.Source)
	require.Equal(t, int32(95), price.Confidence)
	// Timestamp is set to 0 by default (not automatically set)
	require.Equal(t, int64(0), price.Timestamp)
}

func TestPrice_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		price   types.Price
		wantErr bool
		errType error
	}{
		{
			name: "valid price",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(10),
				Confidence: 95,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: false,
		},
		{
			name: "valid - zero price",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(0),
				Confidence: 95,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: false,
		},
		{
			name: "valid - confidence 0",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(10),
				Confidence: 0,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: false,
		},
		{
			name: "valid - confidence 100",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(10),
				Confidence: 100,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty symbol",
			price: types.Price{
				Symbol:     "",
				Price:      sdk.NewDec(10),
				Confidence: 95,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: true,
			errType: types.ErrInvalidSymbol,
		},
		{
			name: "invalid - nil price",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.Dec{},
				Confidence: 95,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: true,
			errType: types.ErrInvalidPrice,
		},
		{
			name: "invalid - negative price",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(-10),
				Confidence: 95,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: true,
			errType: types.ErrInvalidPrice,
		},
		{
			name: "invalid - negative confidence",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(10),
				Confidence: -1,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: true,
			errType: types.ErrLowConfidence,
		},
		{
			name: "invalid - confidence over 100",
			price: types.Price{
				Symbol:     "STT",
				Price:      sdk.NewDec(10),
				Confidence: 101,
				Source:     types.PriceSourceChainlink,
			},
			wantErr: true,
			errType: types.ErrLowConfidence,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.price.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					require.ErrorIs(t, err, tt.errType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPrice_IsStale(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name     string
		price    types.Price
		maxAge   time.Duration
		expected bool
	}{
		{
			name: "fresh price",
			price: types.Price{
				Timestamp: now,
			},
			maxAge:   time.Hour,
			expected: false,
		},
		{
			name: "stale price",
			price: types.Price{
				Timestamp: now - 7200, // 2 hours ago
			},
			maxAge:   time.Hour,
			expected: true,
		},
		{
			name: "exactly at boundary - not stale",
			price: types.Price{
				Timestamp: now - 3600, // 1 hour ago
			},
			maxAge:   time.Hour,
			expected: false, // > not >=
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.price.IsStale(tt.maxAge)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestPrice_String(t *testing.T) {
	price := types.Price{
		Symbol:     "STT",
		Price:      sdk.NewDec(10),
		Confidence: 95,
		Source:     types.PriceSourceChainlink,
		Timestamp:  1234567890,
	}

	result := price.String()
	require.Contains(t, result, "STT")
	require.Contains(t, result, "10")
	require.Contains(t, result, "95")
	require.Contains(t, result, "PRICE_SOURCE_CHAINLINK")
}

func TestPriceSourceFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.PriceSource
	}{
		{"chainlink", "chainlink", types.PriceSourceChainlink},
		{"manual", "manual", types.PriceSourceManual},
		{"empty", "", types.PriceSourceManual},
		{"unknown", "unknown", types.PriceSourceManual},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.PriceSourceFromString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestPriceSourceToString(t *testing.T) {
	tests := []struct {
		name     string
		input    types.PriceSource
		expected string
	}{
		{"chainlink", types.PriceSourceChainlink, "chainlink"},
		{"manual", types.PriceSourceManual, "manual"},
		{"unspecified", types.PriceSourceUnspecified, "manual"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.PriceSourceToString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

// LLMPrice Tests

func TestLLMPrice_ConvertToSTT(t *testing.T) {
	llmPrice := types.LLMPrice{
		Provider:    "openai",
		Model:       "gpt-4",
		InputPrice:  sdk.NewDecWithPrec(3, 0),  // $3 per 1K tokens
		OutputPrice: sdk.NewDecWithPrec(6, 0),  // $6 per 1K tokens
		Currency:    "USD",
	}

	usdToSTTRate := sdk.NewDec(10) // 1 USD = 10 STT
	inputSTT, outputSTT := llmPrice.ConvertToSTT(usdToSTTRate)

	require.True(t, inputSTT.Equal(sdk.NewDec(30)))  // $3 * 10 = 30 STT
	require.True(t, outputSTT.Equal(sdk.NewDec(60))) // $6 * 10 = 60 STT
}

// Marshal/Unmarshal Tests

func TestPrice_MarshalUnmarshal(t *testing.T) {
	original := types.NewPrice("STT", sdk.NewDec(10), types.PriceSourceChainlink, 95)

	// Marshal
	data, err := original.Marshal()
	require.NoError(t, err)
	require.NotNil(t, data)
	require.True(t, len(data) > 0)

	// Unmarshal
	var restored types.Price
	err = restored.Unmarshal(data)
	require.NoError(t, err)

	require.Equal(t, original.Symbol, restored.Symbol)
	require.True(t, original.Price.Equal(restored.Price))
	require.Equal(t, original.Confidence, restored.Confidence)
	require.Equal(t, original.Source, restored.Source)
	require.Equal(t, original.Timestamp, restored.Timestamp)
}

func TestPrice_Size(t *testing.T) {
	price := types.NewPrice("STT", sdk.NewDec(10), types.PriceSourceChainlink, 95)

	size := price.Size()
	data, _ := price.Marshal()
	require.Equal(t, len(data), size)
}

// Genesis Tests

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Prices)
	// Default genesis now includes STT/USD price
	require.Len(t, genesis.Prices, 1)
	require.Equal(t, "STT/USD", genesis.Prices[0].Symbol)
}

func TestValidateGenesis(t *testing.T) {
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
			name: "valid genesis with prices",
			data: types.GenesisState{
				Prices: []types.Price{
					{Symbol: "STT", Price: sdk.NewDec(10), Confidence: 95},
					{Symbol: "BTC", Price: sdk.NewDec(50000), Confidence: 90},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate symbols",
			data: types.GenesisState{
				Prices: []types.Price{
					{Symbol: "STT", Price: sdk.NewDec(10), Confidence: 95},
					{Symbol: "STT", Price: sdk.NewDec(11), Confidence: 90},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty symbol",
			data: types.GenesisState{
				Prices: []types.Price{
					{Symbol: "", Price: sdk.NewDec(10), Confidence: 95},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - negative price",
			data: types.GenesisState{
				Prices: []types.Price{
					{Symbol: "STT", Price: sdk.NewDec(-10), Confidence: 95},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - confidence over 100",
			data: types.GenesisState{
				Prices: []types.Price{
					{Symbol: "STT", Price: sdk.NewDec(10), Confidence: 101},
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
