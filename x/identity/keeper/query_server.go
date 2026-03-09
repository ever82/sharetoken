package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/types"
)

// queryServer is the query server
type queryServer struct {
	*Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

// RegisterQueryServer registers the query server functions
func RegisterQueryServer(server interface{}, srv types.QueryServer) {
	// This is a placeholder - actual implementation depends on the framework
}

var _ types.QueryServer = &queryServer{}

// Params implements QueryServer
func (k queryServer) Params(ctx sdk.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// Identity implements QueryServer
func (k queryServer) Identity(ctx sdk.Context, req *types.QueryIdentityRequest) (*types.QueryIdentityResponse, error) {
	identity, found := k.GetIdentity(ctx, req.Address)
	if !found {
		return nil, types.ErrIdentityNotFound.Wrap(req.Address)
	}

	return &types.QueryIdentityResponse{
		Identity: identity,
	}, nil
}

// Identities implements QueryServer
func (k queryServer) Identities(ctx sdk.Context, req *types.QueryIdentitiesRequest) (*types.QueryIdentitiesResponse, error) {
	identities := k.GetAllIdentities(ctx)

	return &types.QueryIdentitiesResponse{
		Identities: identities,
	}, nil
}

// LimitConfig implements QueryServer
func (k queryServer) LimitConfig(ctx sdk.Context, req *types.QueryLimitConfigRequest) (*types.QueryLimitConfigResponse, error) {
	limitConfig, found := k.GetLimitConfig(ctx, req.Address)
	if !found {
		return nil, types.ErrInvalidLimitConfig.Wrapf("limit config not found for %s", req.Address)
	}

	return &types.QueryLimitConfigResponse{
		LimitConfig: limitConfig,
	}, nil
}

// IsVerified implements QueryServer
func (k queryServer) IsVerified(ctx sdk.Context, req *types.QueryIsVerifiedRequest) (*types.QueryIsVerifiedResponse, error) {
	isVerified := k.Keeper.IsVerified(ctx, req.Address)

	return &types.QueryIsVerifiedResponse{
		IsVerified: isVerified,
	}, nil
}
