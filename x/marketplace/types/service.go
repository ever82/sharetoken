package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "marketplace"
	StoreKey   = ModuleName
)

// ServiceLevel represents the level of service
type ServiceLevel int

const (
	ServiceLevelLLM      ServiceLevel = 1 // Level 1: LLM API
	ServiceLevelAgent    ServiceLevel = 2 // Level 2: Agent
	ServiceLevelWorkflow ServiceLevel = 3 // Level 3: Workflow
)

// PricingMode represents the pricing mode
type PricingMode string

const (
	PricingModeFixed   PricingMode = "fixed"
	PricingModeDynamic PricingMode = "dynamic"
	PricingModeAuction PricingMode = "auction"
)

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

// String implements stringer
func (s Service) String() string {
	return fmt.Sprintf("Service{%s: %s (Level %d), Price: %s}",
		s.ID, s.Name, s.Level, s.Price.String())
}
