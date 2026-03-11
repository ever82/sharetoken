package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper interface
// used by various modules for account operations.
type AccountKeeper interface {
	// GetAccount returns the account for the given address
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// HasAccount checks if an account exists for the given address
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
}

// BankKeeper defines the expected bank keeper interface
// used by various modules for balance and coin operations.
type BankKeeper interface {
	// SpendableCoins returns the spendable coins for the given address
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// GetBalance returns the balance of a specific denom for the given address
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// ParamSubspace defines the expected subspace interface for parameters
type ParamSubspace interface {
	// HasKeyTable checks if the subspace has a key table
	HasKeyTable() bool
	// WithKeyTable returns the subspace with the given key table
	WithKeyTable(table interface{}) interface{}
}
