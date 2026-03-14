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

	apiKey, found := k.GetAPIKey(sdkCtx, req.Id)
	if !found {
		return nil, types.ErrAPIKeyNotFound
	}

	return &types.QueryAPIKeyResponse{
		ApiKey: apiKey,
	}, nil
}

// APIKeysByOwner implements the APIKeysByOwner gRPC method
func (k queryServer) APIKeysByOwner(ctx context.Context, req *types.QueryAPIKeysByOwnerRequest) (*types.QueryAPIKeysByOwnerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allKeys := k.GetAPIKeysByOwner(sdkCtx, req.Owner)

	// Apply pagination
	total := len(allKeys)
	start := 0
	end := total

	if req.Pagination != nil {
		start = int(req.Pagination.Offset)
		limit := int(req.Pagination.Limit)
		if limit > 0 {
			end = start + limit
		}
	}

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &types.QueryAPIKeysByOwnerResponse{
		ApiKeys: allKeys[start:end],
		Total:   uint64(total),
	}, nil
}

// AllAPIKeys implements the AllAPIKeys gRPC method
func (k queryServer) AllAPIKeys(ctx context.Context, req *types.QueryAllAPIKeysRequest) (*types.QueryAllAPIKeysResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allKeys := k.GetAllAPIKeys(sdkCtx)

	// Apply pagination
	total := len(allKeys)
	start := 0
	end := total

	if req.Pagination != nil {
		start = int(req.Pagination.Offset)
		limit := int(req.Pagination.Limit)
		if limit > 0 {
			end = start + limit
		}
	}

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &types.QueryAllAPIKeysResponse{
		ApiKeys: allKeys[start:end],
		Total:   uint64(total),
	}, nil
}

// UsageStats implements the UsageStats gRPC method
func (k queryServer) UsageStats(ctx context.Context, req *types.QueryUsageStatsRequest) (*types.QueryUsageStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats, found := k.GetAPIKeyStats(sdkCtx, req.ApiKeyId)
	if !found {
		// Return empty stats if not found
		stats = *types.NewAPIKeyUsageStats(req.ApiKeyId)
	}

	return &types.QueryUsageStatsResponse{
		Stats: stats,
	}, nil
}

// DailyUsage implements the DailyUsage gRPC method
func (k queryServer) DailyUsage(ctx context.Context, req *types.QueryDailyUsageRequest) (*types.QueryDailyUsageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats, found := k.GetDailyStats(sdkCtx, req.Date, req.ApiKeyId)
	if !found {
		// Return empty stats if not found
		stats = *types.NewDailyUsageStats(req.Date, req.ApiKeyId)
	}

	return &types.QueryDailyUsageResponse{
		DailyStats: []types.DailyUsageStats{stats},
		Total:      1,
	}, nil
}

// ServiceUsage implements the ServiceUsage gRPC method
func (k queryServer) ServiceUsage(ctx context.Context, req *types.QueryServiceUsageRequest) (*types.QueryServiceUsageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats, found := k.GetServiceStats(sdkCtx, req.ServiceId, req.ApiKeyId)
	if !found {
		// Return empty stats if not found
		stats = types.ServiceUsageStats{
			ServiceId: req.ServiceId,
		}
	}

	return &types.QueryServiceUsageResponse{
		Stats: stats,
	}, nil
}
