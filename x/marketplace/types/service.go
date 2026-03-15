package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Service level constants for convenience
const (
	ServiceLevelUnspecified = ServiceLevel_SERVICE_LEVEL_UNSPECIFIED
	ServiceLevelLLM         = ServiceLevel_SERVICE_LEVEL_LLM
	ServiceLevelAgent       = ServiceLevel_SERVICE_LEVEL_AGENT
	ServiceLevelWorkflow    = ServiceLevel_SERVICE_LEVEL_WORKFLOW
)

// Pricing mode constants for convenience
const (
	PricingModeUnspecified = PricingMode_PRICING_MODE_UNSPECIFIED
	PricingModeFixed       = PricingMode_PRICING_MODE_FIXED
	PricingModeDynamic     = PricingMode_PRICING_MODE_DYNAMIC
	PricingModeAuction     = PricingMode_PRICING_MODE_AUCTION
)

// NewService creates a new service with the given parameters
func NewService(id, provider, name string, level ServiceLevel, price sdk.Coins) *Service {
	return &Service{
		Id:          id,
		Provider:    provider,
		Name:        name,
		Level:       level,
		Price:       price,
		PricingMode: PricingModeFixed, // Default pricing mode
		Active:      true,
		CreatedAt:   time.Now().Unix(),
	}
}

// PricingModeFromString converts a string to a PricingMode
func PricingModeFromString(s string) PricingMode {
	switch s {
	case "dynamic":
		return PricingModeDynamic
	case "auction":
		return PricingModeAuction
	case "fixed":
		return PricingModeFixed
	default:
		return PricingModeFixed // Default to fixed
	}
}

// ValidateBasic performs basic validation of the service
func (s Service) ValidateBasic() error {
	if s.Id == "" {
		return ErrInvalidService.Wrap("service ID cannot be empty")
	}
	if s.Provider == "" {
		return ErrInvalidService.Wrap("provider cannot be empty")
	}
	if s.Name == "" {
		return ErrInvalidService.Wrap("name cannot be empty")
	}
	if !sdk.Coins(s.Price).IsZero() {
		if err := sdk.Coins(s.Price).Validate(); err != nil {
			return ErrInvalidService.Wrapf("invalid price: %s", err.Error())
		}
	}
	return nil
}
