package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/oracle/types"
)

func TestNewPrice(t *testing.T) {
	price := types.NewPrice("LLM-API/USD", sdk.NewDec(10), types.PriceSourceManual, 95)

	require.Equal(t, "LLM-API/USD", price.Symbol)
	require.True(t, price.Price.Equal(sdk.NewDec(10)))
	require.Equal(t, types.PriceSourceManual, price.Source)
	require.Equal(t, int32(95), price.Confidence)
	require.False(t, price.IsStale(time.Hour))
}

func TestPriceValidation(t *testing.T) {
	tests := []struct {
		name    string
		price   types.Price
		wantErr bool
	}{
		{
			name: "valid price",
			price: types.Price{
				Symbol:     "ETH/USD",
				Price:      sdk.NewDec(2000),
				Timestamp:  time.Now().Unix(),
				Source:     types.PriceSourceChainlink,
				Confidence: 95,
			},
			wantErr: false,
		},
		{
			name: "missing symbol",
			price: types.Price{
				Symbol:     "",
				Price:      sdk.NewDec(2000),
				Timestamp:  time.Now().Unix(),
				Source:     types.PriceSourceChainlink,
				Confidence: 95,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			price: types.Price{
				Symbol:     "ETH/USD",
				Price:      sdk.NewDec(-100),
				Timestamp:  time.Now().Unix(),
				Source:     types.PriceSourceChainlink,
				Confidence: 95,
			},
			wantErr: true,
		},
		{
			name: "low confidence",
			price: types.Price{
				Symbol:     "ETH/USD",
				Price:      sdk.NewDec(2000),
				Timestamp:  time.Now().Unix(),
				Source:     types.PriceSourceChainlink,
				Confidence: 150,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.price.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPriceStale(t *testing.T) {
	price := types.NewPrice("ETH/USD", sdk.NewDec(2000), types.PriceSourceChainlink, 95)

	// Should not be stale immediately
	require.False(t, price.IsStale(time.Hour))

	// Simulate old price
	price.Timestamp = time.Now().Add(-2 * time.Hour).Unix()
	require.True(t, price.IsStale(time.Hour))
}

func TestLLMPriceConvertToSTT(t *testing.T) {
	llmPrice := types.LLMPrice{
		Provider:    "openai",
		Model:       "gpt-4",
		InputPrice:  sdk.NewDecWithPrec(3, 2), // $0.03 per 1K tokens = 0.03
		OutputPrice: sdk.NewDecWithPrec(6, 2), // $0.06 per 1K tokens = 0.06
		Currency:    "USD",
	}

	// 1 USD = 10 STT
	usdToSTTRate := sdk.NewDec(10)
	inputSTT, outputSTT := llmPrice.ConvertToSTT(usdToSTTRate)

	// $0.03 * 10 = 0.3 STT
	require.True(t, inputSTT.Equal(sdk.NewDecWithPrec(3, 1)))
	// $0.06 * 10 = 0.6 STT
	require.True(t, outputSTT.Equal(sdk.NewDecWithPrec(6, 1)))
}

func TestPriceString(t *testing.T) {
	price := types.NewPrice("ETH/USD", sdk.NewDec(2000), types.PriceSourceChainlink, 95)
	str := price.String()

	require.Contains(t, str, "ETH/USD")
	require.Contains(t, str, "2000")
	require.Contains(t, str, "95")
	require.Contains(t, str, "PRICE_SOURCE_CHAINLINK")
}
