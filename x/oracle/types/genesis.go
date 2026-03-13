package types

import (
	"fmt"
)

// GenesisState defines the oracle module's genesis state.
type GenesisState struct {
	Prices []Price `json:"prices"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Prices: []Price{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	seenSymbols := make(map[string]bool)
	for _, price := range data.Prices {
		if seenSymbols[price.Symbol] {
			return fmt.Errorf("duplicate price symbol: %s", price.Symbol)
		}
		seenSymbols[price.Symbol] = true

		if err := price.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid price for symbol %s: %w", price.Symbol, err)
		}
	}
	return nil
}
