package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"

	identitytypes "sharetoken/x/identity/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) identitytypes.Params {
	return identitytypes.NewParams(
		k.VerificationRequired(ctx),
		k.AllowedProviders(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, ps identitytypes.Params) {
	k.paramStore.SetParamSet(ctx, &ps)
}

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable().RegisterParamSet(&identitytypes.Params{})
}

// VerificationRequired returns if verification is required
func (k Keeper) VerificationRequired(ctx sdk.Context) (res bool) {
	k.paramStore.Get(ctx, identitytypes.KeyVerificationRequired, &res)
	return
}

// AllowedProviders returns the list of allowed providers
func (k Keeper) AllowedProviders(ctx sdk.Context) (res []string) {
	k.paramStore.Get(ctx, identitytypes.KeyAllowedProviders, &res)
	return
}

// SetVerificationRequired sets if verification is required
func (k Keeper) SetVerificationRequired(ctx sdk.Context, required bool) {
	k.paramStore.Set(ctx, identitytypes.KeyVerificationRequired, required)
}

// SetAllowedProviders sets the allowed providers
func (k Keeper) SetAllowedProviders(ctx sdk.Context, providers []string) {
	k.paramStore.Set(ctx, identitytypes.KeyAllowedProviders, providers)
}

// IsProviderAllowed checks if a provider is allowed
func (k Keeper) IsProviderAllowed(ctx sdk.Context, provider string) bool {
	allowed := k.AllowedProviders(ctx)
	for _, p := range allowed {
		if p == provider {
			return true
		}
	}
	return false
}
