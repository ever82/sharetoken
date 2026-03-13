package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/escrow/types"
)

// RegisterInvariants registers the escrow module invariants
func (k Keeper) RegisterInvariants(ir sdk.InvariantRegistry) {
	ir.RegisterRoute(types.ModuleName, "escrow-amounts",
		k.EscrowAmountsInvariant())
	ir.RegisterRoute(types.ModuleName, "escrow-status",
		k.EscrowStatusInvariant())
	ir.RegisterRoute(types.ModuleName, "escrow-expiration",
		k.EscrowExpirationInvariant())
}

// EscrowAmountsInvariant checks that all escrow amounts are valid and positive
func (k Keeper) EscrowAmountsInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidEscrows []string

		escrows := k.GetAllEscrows(ctx)
		for _, escrow := range escrows {
			if !escrow.Amount.IsValid() || escrow.Amount.IsZero() {
				invalidEscrows = append(invalidEscrows, escrow.ID)
			}
		}

		if len(invalidEscrows) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"escrow-amounts",
				fmt.Sprintf("found %d escrows with invalid amounts: %v", len(invalidEscrows), invalidEscrows),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"escrow-amounts",
			"all escrow amounts are valid",
		), false
	}
}

// EscrowStatusInvariant checks that escrow statuses are valid
func (k Keeper) EscrowStatusInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidStatuses []string

		escrows := k.GetAllEscrows(ctx)
		for _, escrow := range escrows {
			switch escrow.Status {
			case types.EscrowStatusPending,
				types.EscrowStatusCompleted,
				types.EscrowStatusDisputed,
				types.EscrowStatusRefunded:
				// Valid status
			default:
				invalidStatuses = append(invalidStatuses, fmt.Sprintf("%s:%s", escrow.ID, escrow.Status))
			}
		}

		if len(invalidStatuses) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"escrow-status",
				fmt.Sprintf("found %d escrows with invalid statuses: %v", len(invalidStatuses), invalidStatuses),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"escrow-status",
			"all escrow statuses are valid",
		), false
	}
}

// EscrowExpirationInvariant checks that escrow expiration times are valid
func (k Keeper) EscrowExpirationInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidExpirations []string

		escrows := k.GetAllEscrows(ctx)
		for _, escrow := range escrows {
			// CreatedAt must be before ExpiresAt
			if escrow.ExpiresAt <= escrow.CreatedAt {
				invalidExpirations = append(invalidExpirations, escrow.ID)
			}
		}

		if len(invalidExpirations) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"escrow-expiration",
				fmt.Sprintf("found %d escrows with invalid expiration times: %v", len(invalidExpirations), invalidExpirations),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"escrow-expiration",
			"all escrow expiration times are valid",
		), false
	}
}

// AllInvariants runs all escrow invariants
func (k Keeper) AllInvariants(ctx sdk.Context) (string, bool) {
	res, stop := k.EscrowAmountsInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.EscrowStatusInvariant()(ctx)
	if stop {
		return res, stop
	}

	return k.EscrowExpirationInvariant()(ctx)
}
