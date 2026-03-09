package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/keeper"
	"sharetoken/x/identity/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set params
	k.SetParams(ctx, genState.Params)

	// Set all identities
	for _, identity := range genState.Identities {
		k.SetIdentity(ctx, identity)
	}

	// Set all limit configs
	for _, limitConfig := range genState.LimitConfigs {
		k.SetLimitConfig(ctx, limitConfig)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Params = k.GetParams(ctx)
	genesis.Identities = k.GetAllIdentities(ctx)
	genesis.LimitConfigs = k.GetAllLimitConfigs(ctx)

	return genesis
}
