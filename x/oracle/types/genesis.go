package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Prices: []Price{
			{
				Symbol:     "STT/USD",
				Price:      sdk.NewDec(1),
				Timestamp:  0,
				Source:     PriceSource_PRICE_SOURCE_MANUAL,
				Confidence: 100,
			},
		},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(gs GenesisState) error {
	seenSymbols := make(map[string]bool)
	for _, price := range gs.Prices {
		if price.Symbol == "" {
			return fmt.Errorf("price symbol cannot be empty")
		}
		if seenSymbols[price.Symbol] {
			return fmt.Errorf("duplicate price symbol: %s", price.Symbol)
		}
		seenSymbols[price.Symbol] = true
		if price.Price.IsNil() || price.Price.IsNegative() {
			return fmt.Errorf("invalid price for symbol: %s", price.Symbol)
		}
	}
	return nil
}
