package keeper

import (
	"context"

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
func (k queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// Identity implements QueryServer
func (k queryServer) Identity(ctx context.Context, req *types.QueryIdentityRequest) (*types.QueryIdentityResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	identity, found := k.GetIdentity(sdkCtx, req.Address)
	if !found {
		return nil, types.ErrIdentityNotFound.Wrap(req.Address)
	}

	return &types.QueryIdentityResponse{
		Identity: identity,
	}, nil
}

// Identities implements QueryServer
func (k queryServer) Identities(ctx context.Context, req *types.QueryIdentitiesRequest) (*types.QueryIdentitiesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	identities := k.GetAllIdentities(sdkCtx)

	return &types.QueryIdentitiesResponse{
		Identities: identities,
	}, nil
}

// LimitConfig implements QueryServer
func (k queryServer) LimitConfig(ctx context.Context, req *types.QueryLimitConfigRequest) (*types.QueryLimitConfigResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	limitConfig, found := k.GetLimitConfig(sdkCtx, req.Address)
	if !found {
		return nil, types.ErrInvalidLimitConfig.Wrapf("limit config not found for %s", req.Address)
	}

	return &types.QueryLimitConfigResponse{
		LimitConfig: limitConfig,
	}, nil
}

// IsVerified implements QueryServer
func (k queryServer) IsVerified(ctx context.Context, req *types.QueryIsVerifiedRequest) (*types.QueryIsVerifiedResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	isVerified := k.Keeper.IsVerified(sdkCtx, req.Address)

	return &types.QueryIsVerifiedResponse{
		IsVerified: isVerified,
	}, nil
}
