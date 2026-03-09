package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"sharetoken/x/identity/types"
)

type (
	// Keeper of this module
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramStore paramtypes.Subspace

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
	}
)

// NewKeeper creates a new identity Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramStore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the module's authority address (x/gov by default)
func (k Keeper) GetAuthority() string {
	return "sharetoken10d07y265gmmuvt4z0w9aw880jnsr700jfw7aah" // gov module address placeholder
}

// IsAuthority checks if the given address is the authority
func (k Keeper) IsAuthority(ctx sdk.Context, addr sdk.AccAddress) bool {
	authorityAddr, err := sdk.AccAddressFromBech32(k.GetAuthority())
	if err != nil {
		return false
	}
	return addr.Equals(authorityAddr)
}
