package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// AccountI is the interface for accounts
type AccountI interface {
	GetAddress() []byte
	GetPubKey() []byte
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx Context, addr AccAddress) AccountI
	SetAccount(ctx Context, acc AccountI)
	NewAccountWithAddress(ctx Context, addr AccAddress) AccountI
}

// Context is a type alias for sdk.Context
type Context struct{}

// AccAddress is a type alias for sdk.AccAddress
type AccAddress []byte

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	GetBalance(ctx Context, addr AccAddress, denom string) Coin
	SendCoins(ctx Context, fromAddr AccAddress, toAddr AccAddress, amt Coins) error
	SendCoinsFromAccountToModule(ctx Context, senderAddr AccAddress, recipientModule string, amt Coins) error
	SendCoinsFromModuleToAccount(ctx Context, senderModule string, recipientAddr AccAddress, amt Coins) error
}

// Coin represents a coin
type Coin struct {
	Denom  string
	Amount int64
}

// Coins represents multiple coins
type Coins []Coin

// ParamSubspace defines the expected Subspace interface for module parameters
type ParamSubspace interface {
	Get(ctx Context, key []byte, ptr interface{})
	Set(ctx Context, key []byte, param interface{})
	HasKeyTable() bool
	WithKeyTable(table KeyTable) ParamSubspace
}

// KeyTable is a type alias for the paramspace KeyTable
type KeyTable struct{}

// PageRequest is a type alias for query PageRequest
type PageRequest = query.PageRequest
