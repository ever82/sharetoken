package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

// PriceSource represents the source of price data
type PriceSource int32

const (
	PriceSourceUnspecified PriceSource = 0
	PriceSourceChainlink   PriceSource = 1 // Chainlink oracle
	PriceSourceManual      PriceSource = 2 // Manual price

	// PriceSourceConfidenceThreshold is the minimum confidence threshold
	PriceSourceConfidenceThreshold = 50
	// MaxConfidence is the maximum confidence value
	MaxConfidence = 100
	// TokenUnitDivisor is the divisor for token calculations
	TokenUnitDivisor = 1000
)

// PriceSourceFromString converts string to PriceSource
func PriceSourceFromString(s string) PriceSource {
	switch s {
	case "chainlink":
		return PriceSourceChainlink
	default:
		return PriceSourceManual
	}
}

// PriceSourceToString converts PriceSource to string
func PriceSourceToString(s PriceSource) string {
	switch s {
	case PriceSourceChainlink:
		return "chainlink"
	default:
		return "manual"
	}
}

// Price represents a price data point
type Price struct {
	Symbol     string      `json:"symbol"`
	Price      sdk.Dec     `json:"price"`
	Timestamp  int64       `json:"timestamp"`
	Source     PriceSource `json:"source"`
	Confidence int32       `json:"confidence"`
}

// Reset implements proto.Message
func (m *Price) Reset() { *m = Price{} }

// String implements proto.Message
func (m Price) String() string {
	return fmt.Sprintf("%s: %s (confidence: %d%%, source: %s, time: %d)",
		m.Symbol, m.Price.String(), m.Confidence, PriceSourceToString(m.Source), m.Timestamp)
}

// ProtoMessage implements proto.Message
func (*Price) ProtoMessage() {}

// Marshal implements codec.ProtoMarshaler
func (m Price) Marshal() ([]byte, error) {
	return proto.Marshal(m.toProto())
}

// MarshalTo implements codec.ProtoMarshaler
func (m Price) MarshalTo(data []byte) (n int, err error) {
	return m.MarshalToSizedBuffer(data)
}

// MarshalToSizedBuffer implements codec.ProtoMarshaler
func (m Price) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	encoded, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(encoded)
	if len(dAtA) < n {
		return 0, fmt.Errorf("buffer too small")
	}
	copy(dAtA[:n], encoded)
	return n, nil
}

// Size implements codec.ProtoMarshaler
func (m Price) Size() int {
	data, _ := m.Marshal()
	return len(data)
}

// Unmarshal implements codec.ProtoMarshaler
func (m *Price) Unmarshal(data []byte) error {
	pm := &PriceProto{}
	if err := proto.Unmarshal(data, pm); err != nil {
		return err
	}
	m.fromProto(pm)
	return nil
}

// toProto converts Price to proto message
func (m Price) toProto() *PriceProto {
	return &PriceProto{
		Symbol:     m.Symbol,
		Price:      m.Price.String(),
		Timestamp:  m.Timestamp,
		Source:     int32(m.Source),
		Confidence: m.Confidence,
	}
}

// fromProto converts proto message to Price
func (m *Price) fromProto(pm *PriceProto) {
	m.Symbol = pm.Symbol
	m.Price, _ = sdk.NewDecFromStr(pm.Price)
	m.Timestamp = pm.Timestamp
	m.Source = PriceSource(pm.Source)
	m.Confidence = pm.Confidence
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
