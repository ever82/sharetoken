package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/escrow/types"
)

// Keeper of the escrow store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey

	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

// NewKeeper creates a new escrow Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
