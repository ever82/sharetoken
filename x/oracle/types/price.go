package types

import (
	"errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Constants for price sources and thresholds
const (
	// PriceSourceConfidenceThreshold is the minimum confidence level required
	PriceSourceConfidenceThreshold = 50

	// TokenUnitDivisor is used for LLM price calculations (1000 tokens)
	TokenUnitDivisor int64 = 1000

	// MaxConfidence is the maximum confidence level (100)
	MaxConfidence = 100
)

// Aliases for protobuf enum values
const (
	// PriceSourceChainlink is an alias for PriceSource_PRICE_SOURCE_CHAINLINK
	PriceSourceChainlink = PriceSource_PRICE_SOURCE_CHAINLINK

	// PriceSourceManual is an alias for PriceSource_PRICE_SOURCE_MANUAL
	PriceSourceManual = PriceSource_PRICE_SOURCE_MANUAL
)

// NewPrice creates a new price instance
func NewPrice(symbol string, price sdk.Dec, source PriceSource, confidence int32) *Price {
	return &Price{
		Symbol:     symbol,
		Price:      price,
		Source:     source,
		Confidence: confidence,
		Timestamp:  0, // Should be set by caller if needed
	}
}

// IsStale checks if the price is older than the given duration
func (p *Price) IsStale(maxAge time.Duration) bool {
	if p.Timestamp == 0 {
		return false
	}
	return time.Now().Unix()-p.Timestamp > int64(maxAge.Seconds())
}

// ValidateBasic performs basic validation of the price
func (p Price) ValidateBasic() error {
	if p.Symbol == "" {
		return errors.New("price symbol cannot be empty")
	}
	if p.Price.IsNil() || p.Price.IsNegative() {
		return errors.New("price must be positive")
	}
	if p.Confidence < 0 || p.Confidence > MaxConfidence {
		return errors.New("confidence must be between 0 and 100")
	}
	if p.Source == PriceSource_PRICE_SOURCE_UNSPECIFIED {
		return errors.New("price source must be specified")
	}
	return nil
}

// ConvertToSTT converts the LLM price to STT Dec values based on USD price rate
// Returns (inputPriceInSTT, outputPriceInSTT) per 1000 tokens
func (p LLMPrice) ConvertToSTT(usdPrice sdk.Dec) (sdk.Dec, sdk.Dec) {
	// Calculate STT amount for input: (InputPrice * usdPrice) / TokenUnitDivisor
	inputSTT := p.InputPrice.Mul(usdPrice)

	// Calculate STT amount for output: (OutputPrice * usdPrice) / TokenUnitDivisor
	outputSTT := p.OutputPrice.Mul(usdPrice)

	return inputSTT, outputSTT
}
