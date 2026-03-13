package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountI is the interface for accounts
type AccountI interface {
	GetAddress() []byte
	GetPubKey() []byte
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) AccountI
	SetAccount(ctx sdk.Context, acc AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) AccountI
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for module parameters
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	HasKeyTable() bool
	WithKeyTable(table paramtypes.KeyTable) ParamSubspace
}

// PageRequest is a type alias for query PageRequest
type PageRequest = query.PageRequest
