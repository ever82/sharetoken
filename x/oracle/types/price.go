package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceSource represents the source of price data
type PriceSource string

const (
	// PriceSourceChainlink represents Chainlink oracle
	PriceSourceChainlink PriceSource = "chainlink"
	// PriceSourceManual represents manually set price
	PriceSourceManual PriceSource = "manual"
)

// Price represents a price data point
type Price struct {
	Symbol     string      `json:"symbol"`
	Price      sdk.Dec     `json:"price"`
	Timestamp  int64       `json:"timestamp"`
	Source     PriceSource `json:"source"`
	Confidence int32       `json:"confidence"`
}

// NewPrice creates a new price
func NewPrice(symbol string, price sdk.Dec, source PriceSource, confidence int32) *Price {
	return &Price{
		Symbol:     symbol,
		Price:      price,
		Timestamp:  time.Now().Unix(),
		Source:     source,
		Confidence: confidence,
	}
}

// ValidateBasic performs basic validation
func (p Price) ValidateBasic() error {
	if p.Symbol == "" {
		return ErrInvalidSymbol
	}
	if p.Price.IsNil() || p.Price.IsNegative() {
		return ErrInvalidPrice
	}
	if p.Confidence < 0 || p.Confidence > 100 {
		return ErrLowConfidence.Wrap("confidence must be between 0 and 100")
	}
	return nil
}

// IsStale checks if the price is stale (older than maxAge)
func (p Price) IsStale(maxAge time.Duration) bool {
	return time.Now().Unix()-p.Timestamp > int64(maxAge.Seconds())
}

// String implements stringer
func (p Price) String() string {
	return fmt.Sprintf("%s: %s (confidence: %d%%, source: %s, time: %d)",
		p.Symbol, p.Price.String(), p.Confidence, p.Source, p.Timestamp)
}

// LLMPrice represents LLM API pricing
type LLMPrice struct {
	Provider    string  `json:"provider"`     // e.g., "openai", "anthropic"
	Model       string  `json:"model"`        // e.g., "gpt-4", "claude-3"
	InputPrice  sdk.Dec `json:"input_price"`  // per 1K tokens
	OutputPrice sdk.Dec `json:"output_price"` // per 1K tokens
	Currency    string  `json:"currency"`     // e.g., "USD"
}

// ConvertToSTT converts USD price to STT
func (lp LLMPrice) ConvertToSTT(usdToSTTRate sdk.Dec) (inputSTT, outputSTT sdk.Dec) {
	inputSTT = lp.InputPrice.Mul(usdToSTTRate)
	outputSTT = lp.OutputPrice.Mul(usdToSTTRate)
	return
}
