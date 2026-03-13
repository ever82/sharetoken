package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/types"
)

// queryServer implements the QueryServer interface
type queryServer struct {
	*Keeper
}

// NewQueryServerImpl creates a new query server
func NewQueryServerImpl(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

// APIKey implements the APIKey gRPC method
func (k queryServer) APIKey(ctx context.Context, req *types.QueryAPIKeyRequest) (*types.QueryAPIKeyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	apiKey, found := k.GetAPIKey(sdkCtx, req.ID)
	if !found {
		return nil, types.ErrAPIKeyNotFound
	}

	return &types.QueryAPIKeyResponse{
		APIKey: apiKey,
	}, nil
}

// APIKeysByOwner implements the APIKeysByOwner gRPC method
func (k queryServer) APIKeysByOwner(ctx context.Context, req *types.QueryAPIKeysByOwnerRequest) (*types.QueryAPIKeysByOwnerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allKeys := k.GetAPIKeysByOwner(sdkCtx, req.Owner)

	// Apply pagination
	total := len(allKeys)
	start := req.Offset
	end := start + req.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	if req.Limit == 0 {
		end = total
	}

	return &types.QueryAPIKeysByOwnerResponse{
		APIKeys: allKeys[start:end],
		Total:   total,
	}, nil
}

// AllAPIKeys implements the AllAPIKeys gRPC method
func (k queryServer) AllAPIKeys(ctx context.Context, req *types.QueryAllAPIKeysRequest) (*types.QueryAllAPIKeysResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allKeys := k.GetAllAPIKeys(sdkCtx)

	// Apply pagination
	total := len(allKeys)
	start := req.Offset
	end := start + req.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	if req.Limit == 0 {
		end = total
	}

	return &types.QueryAllAPIKeysResponse{
		APIKeys: allKeys[start:end],
		Total:   total,
	}, nil
}

// UsageStats implements the UsageStats gRPC method
func (k queryServer) UsageStats(ctx context.Context, req *types.QueryUsageStatsRequest) (*types.QueryUsageStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats, found := k.GetAPIKeyStats(sdkCtx, req.APIKeyID)
	if !found {
		// Return empty stats if not found
		stats = *types.NewAPIKeyUsageStats(req.APIKeyID)
	}

	return &types.QueryUsageStatsResponse{
		Stats: stats,
	}, nil
}

// DailyUsage implements the DailyUsage gRPC method
func (k queryServer) DailyUsage(ctx context.Context, req *types.QueryDailyUsageRequest) (*types.QueryDailyUsageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats, found := k.GetDailyStats(sdkCtx, req.Date, req.APIKeyID)
	if !found {
		// Return empty stats if not found
		stats = *types.NewDailyUsageStats(req.Date, req.APIKeyID)
	}

	return &types.QueryDailyUsageResponse{
		DailyStats: []types.DailyUsageStats{stats},
		Total:      1,
	}, nil
}

// ServiceUsage implements the ServiceUsage gRPC method
func (k queryServer) ServiceUsage(ctx context.Context, req *types.QueryServiceUsageRequest) (*types.QueryServiceUsageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats, found := k.GetServiceStats(sdkCtx, req.ServiceID, req.APIKeyID)
	if !found {
		// Return empty stats if not found
		stats = types.ServiceUsageStats{
			ServiceID: req.ServiceID,
		}
	}

	return &types.QueryServiceUsageResponse{
		Stats: stats,
	}, nil
}
