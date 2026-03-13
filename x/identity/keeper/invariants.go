package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/types"
)

// RegisterInvariants registers the identity module invariants
func (k Keeper) RegisterInvariants(ir sdk.InvariantRegistry) {
	ir.RegisterRoute(types.ModuleName, "identity-count",
		k.IdentityCountInvariant())
	ir.RegisterRoute(types.ModuleName, "verified-identity",
		k.VerifiedIdentityInvariant())
	ir.RegisterRoute(types.ModuleName, "unique-did",
		k.UniqueDIDInvariant())
	ir.RegisterRoute(types.ModuleName, "provider-validity",
		k.ProviderValidityInvariant())
}

// IdentityCountInvariant checks that the number of registered identities is non-negative
// and that the count matches the actual number of identities in store
func (k Keeper) IdentityCountInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		identities := k.GetAllIdentities(ctx)

		count := len(identities)
		if count < 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"identity-count",
				"identity count is negative (impossible state)",
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"identity-count",
			fmt.Sprintf("registered identities count: %d", count),
		), false
	}
}

// VerifiedIdentityInvariant checks that verified identities have valid verification data
func (k Keeper) VerifiedIdentityInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidVerified []string

		identities := k.GetAllIdentities(ctx)
		for _, identity := range identities {
			if identity.IsVerified {
				// Check that verification provider is set and valid
				if identity.VerificationProvider == "" {
					invalidVerified = append(invalidVerified, fmt.Sprintf("%s:missing provider", identity.Address))
					continue
				}

				// Check that verification provider is in allowed list
				if !types.IsValidProvider(identity.VerificationProvider) {
					invalidVerified = append(invalidVerified, fmt.Sprintf("%s:invalid provider %s", identity.Address, identity.VerificationProvider))
				}

				// Check that verification hash is set
				if identity.VerificationHash == "" {
					invalidVerified = append(invalidVerified, fmt.Sprintf("%s:missing verification hash", identity.Address))
				}
			}
		}

		if len(invalidVerified) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"verified-identity",
				fmt.Sprintf("found %d verified identities with invalid data: %v", len(invalidVerified), invalidVerified),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"verified-identity",
			"all verified identities have valid data",
		), false
	}
}

// UniqueDIDInvariant checks that all DIDs in the system are unique
func (k Keeper) UniqueDIDInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		didMap := make(map[string]string) // DID -> address
		var duplicates []string

		identities := k.GetAllIdentities(ctx)
		for _, identity := range identities {
			if identity.Did != "" {
				if existingAddr, exists := didMap[identity.Did]; exists {
					duplicates = append(duplicates, fmt.Sprintf("%s (used by %s and %s)", identity.Did, existingAddr, identity.Address))
				} else {
					didMap[identity.Did] = identity.Address
				}
			}
		}

		if len(duplicates) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"unique-did",
				fmt.Sprintf("found %d duplicate DIDs: %v", len(duplicates), duplicates),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"unique-did",
			"all DIDs are unique",
		), false
	}
}

// ProviderValidityInvariant checks that all providers used are valid
func (k Keeper) ProviderValidityInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidProviders []string

		identities := k.GetAllIdentities(ctx)
		for _, identity := range identities {
			if identity.VerificationProvider != "" {
				if !types.IsValidProvider(identity.VerificationProvider) {
					invalidProviders = append(invalidProviders, fmt.Sprintf("%s:%s", identity.Address, identity.VerificationProvider))
				}
			}
		}

		if len(invalidProviders) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"provider-validity",
				fmt.Sprintf("found %d invalid providers: %v", len(invalidProviders), invalidProviders),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"provider-validity",
			"all providers are valid",
		), false
	}
}

// AllInvariants runs all identity invariants
func (k Keeper) AllInvariants(ctx sdk.Context) (string, bool) {
	res, stop := k.IdentityCountInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.VerifiedIdentityInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.UniqueDIDInvariant()(ctx)
	if stop {
		return res, stop
	}

	return k.ProviderValidityInvariant()(ctx)
}
