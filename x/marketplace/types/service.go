package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

const (
	ModuleName = "marketplace"
	StoreKey   = ModuleName
)

// ServiceLevel represents the level of service
type ServiceLevel int32

const (
	ServiceLevelUnspecified ServiceLevel = 0
	ServiceLevelLLM         ServiceLevel = 1 // Level 1: LLM API
	ServiceLevelAgent       ServiceLevel = 2 // Level 2: Agent
	ServiceLevelWorkflow    ServiceLevel = 3 // Level 3: Workflow
)

// PricingMode represents the pricing mode
type PricingMode int32

const (
	PricingModeUnspecified PricingMode = 0
	PricingModeFixed       PricingMode = 1
	PricingModeDynamic     PricingMode = 2
	PricingModeAuction     PricingMode = 3
)

// PricingModeFromString converts string to PricingMode
func PricingModeFromString(s string) PricingMode {
	switch s {
	case "dynamic":
		return PricingModeDynamic
	case "auction":
		return PricingModeAuction
	default:
		return PricingModeFixed
	}
}

// Service represents a service offering
type Service struct {
	ID          string       `json:"id"`
	Provider    string       `json:"provider"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Level       ServiceLevel `json:"level"`
	PricingMode PricingMode  `json:"pricing_mode"`
	Price       sdk.Coins    `json:"price"`
	Active      bool         `json:"active"`
	CreatedAt   int64        `json:"created_at"`
}

// Reset implements proto.Message
func (m *Service) Reset() { *m = Service{} }

// String implements proto.Message
func (m Service) String() string {
	return fmt.Sprintf("Service{%s: %s (Level %d), Price: %s}",
		m.ID, m.Name, m.Level, m.Price.String())
}

// ProtoMessage implements proto.Message
func (*Service) ProtoMessage() {}

// Marshal implements codec.ProtoMarshaler
func (m Service) Marshal() ([]byte, error) {
	return proto.Marshal(m.toProto())
}

// MarshalTo implements codec.ProtoMarshaler
func (m Service) MarshalTo(data []byte) (n int, err error) {
	return m.MarshalToSizedBuffer(data)
}

// MarshalToSizedBuffer implements codec.ProtoMarshaler
func (m Service) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m Service) Size() int {
	data, _ := m.Marshal()
	return len(data)
}

// Unmarshal implements codec.ProtoMarshaler
func (m *Service) Unmarshal(data []byte) error {
	pm := &ServiceProto{}
	if err := proto.Unmarshal(data, pm); err != nil {
		return err
	}
	m.fromProto(pm)
	return nil
}

// toProto converts Service to proto message
func (m Service) toProto() *ServiceProto {
	coins := make([]*CoinProto, len(m.Price))
	for i, c := range m.Price {
		coins[i] = &CoinProto{
			Denom:  c.Denom,
			Amount: c.Amount.String(),
		}
	}
	return &ServiceProto{
		Id:          m.ID,
		Provider:    m.Provider,
		Name:        m.Name,
		Description: m.Description,
		Level:       int32(m.Level),
		PricingMode: int32(m.PricingMode),
		Price:       coins,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
	}
}

// fromProto converts proto message to Service
func (m *Service) fromProto(pm *ServiceProto) {
	m.ID = pm.Id
	m.Provider = pm.Provider
	m.Name = pm.Name
	m.Description = pm.Description
	m.Level = ServiceLevel(pm.Level)
	m.PricingMode = PricingMode(pm.PricingMode)
	m.Price = make(sdk.Coins, len(pm.Price))
	for i, c := range pm.Price {
		amount, _ := sdk.NewIntFromString(c.Amount)
		m.Price[i] = sdk.NewCoin(c.Denom, amount)
	}
	m.Active = pm.Active
	m.CreatedAt = pm.CreatedAt
}

// NewService creates a new service
func NewService(id, provider, name string, level ServiceLevel, price sdk.Coins) *Service {
	return &Service{
		ID:          id,
		Provider:    provider,
		Name:        name,
		Level:       level,
		PricingMode: PricingModeFixed,
		Price:       price,
		Active:      true,
	}
}
