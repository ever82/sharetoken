package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/marketplace/types"
)

// RegisterInvariants registers the marketplace module invariants
func (k Keeper) RegisterInvariants(ir sdk.InvariantRegistry) {
	ir.RegisterRoute(types.ModuleName, "service-status",
		k.ServiceStatusInvariant())
	ir.RegisterRoute(types.ModuleName, "service-pricing",
		k.ServicePricingInvariant())
	ir.RegisterRoute(types.ModuleName, "service-level-validity",
		k.ServiceLevelValidityInvariant())
}

// ServiceStatusInvariant checks that service statuses are consistent
func (k Keeper) ServiceStatusInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidStatuses []string

		services := k.GetAllServices(ctx)
		for _, service := range services {
			// Active services should have valid pricing
			if service.Active {
				if !service.Price.IsValid() {
					invalidStatuses = append(invalidStatuses, fmt.Sprintf("%s:active with invalid price", service.Id))
				}
				if service.Provider == "" {
					invalidStatuses = append(invalidStatuses, fmt.Sprintf("%s:active with no provider", service.Id))
				}
			}

			// Service should have a valid name
			if service.Name == "" {
				invalidStatuses = append(invalidStatuses, fmt.Sprintf("%s:empty name", service.Id))
			}

			// Service should have a valid ID
			if service.Id == "" {
				invalidStatuses = append(invalidStatuses, "service with empty ID")
			}
		}

		if len(invalidStatuses) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"service-status",
				fmt.Sprintf("found %d services with invalid status: %v", len(invalidStatuses), invalidStatuses),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"service-status",
			"all services have valid status",
		), false
	}
}

// ServicePricingInvariant checks that service pricing is valid
func (k Keeper) ServicePricingInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidPricing []string

		services := k.GetAllServices(ctx)
		for _, service := range services {
			// Check pricing mode validity
			switch service.PricingMode {
			case types.PricingMode_PRICING_MODE_FIXED,
				types.PricingMode_PRICING_MODE_DYNAMIC,
				types.PricingMode_PRICING_MODE_AUCTION:
				// Valid pricing mode
			default:
				invalidPricing = append(invalidPricing, fmt.Sprintf("%s:invalid pricing mode %d", service.Id, service.PricingMode))
				continue
			}

			// Price should be valid (not empty for active services)
			if service.Active && (service.Price.IsZero() || !service.Price.IsValid()) {
				invalidPricing = append(invalidPricing, fmt.Sprintf("%s:invalid price %s", service.Id, service.Price.String()))
			}
		}

		if len(invalidPricing) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"service-pricing",
				fmt.Sprintf("found %d services with invalid pricing: %v", len(invalidPricing), invalidPricing),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"service-pricing",
			"all services have valid pricing",
		), false
	}
}

// ServiceLevelValidityInvariant checks that service levels are valid
func (k Keeper) ServiceLevelValidityInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidLevels []string

		services := k.GetAllServices(ctx)
		for _, service := range services {
			// Service level should be one of the valid levels
			switch service.Level {
			case types.ServiceLevel_SERVICE_LEVEL_LLM,
				types.ServiceLevel_SERVICE_LEVEL_AGENT,
				types.ServiceLevel_SERVICE_LEVEL_WORKFLOW:
				// Valid level
			default:
				invalidLevels = append(invalidLevels, fmt.Sprintf("%s:invalid level %d", service.Id, service.Level))
			}
		}

		if len(invalidLevels) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"service-level-validity",
				fmt.Sprintf("found %d services with invalid levels: %v", len(invalidLevels), invalidLevels),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"service-level-validity",
			"all services have valid levels",
		), false
	}
}

// AllInvariants runs all marketplace invariants
func (k Keeper) AllInvariants(ctx sdk.Context) (string, bool) {
	res, stop := k.ServiceStatusInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.ServicePricingInvariant()(ctx)
	if stop {
		return res, stop
	}

	return k.ServiceLevelValidityInvariant()(ctx)
}
