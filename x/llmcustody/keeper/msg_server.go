package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/types"
)

// msgServer implements the MsgServer interface
type msgServer struct {
	*Keeper
}

// NewMsgServerImpl creates a new message server
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterAPIKey implements the RegisterAPIKey gRPC method
func (k msgServer) RegisterAPIKey(ctx context.Context, msg *types.MsgRegisterAPIKey) (*types.MsgRegisterAPIKeyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Generate API key ID
	apiKeyID := types.GenerateAPIKeyID()
	provider := types.ProviderFromString(msg.Provider)

	// Create API key
	apiKey := types.NewAPIKey(apiKeyID, provider, msg.EncryptedKey, msg.Owner)
	apiKey.AccessRules = msg.AccessRules

	// Validate
	if err := apiKey.ValidateBasic(); err != nil {
		return nil, err
	}

	// Store
	k.SetAPIKey(sdkCtx, *apiKey)

	return &types.MsgRegisterAPIKeyResponse{
		ApiKeyId: apiKeyID,
	}, nil
}

// UpdateAPIKey implements the UpdateAPIKey gRPC method
func (k msgServer) UpdateAPIKey(ctx context.Context, msg *types.MsgUpdateAPIKey) (*types.MsgUpdateAPIKeyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing API key
	apiKey, found := k.GetAPIKey(sdkCtx, msg.ApiKeyId)
	if !found {
		return nil, types.ErrAPIKeyNotFound
	}

	// Verify ownership
	if apiKey.Owner != msg.Owner {
		return nil, types.ErrUnauthorized
	}

	// Update fields
	if len(msg.AccessRules) > 0 {
		apiKey.AccessRules = msg.AccessRules
	}
	apiKey.Active = msg.Active

	// Store
	k.SetAPIKey(sdkCtx, apiKey)

	return &types.MsgUpdateAPIKeyResponse{}, nil
}

// RevokeAPIKey implements the RevokeAPIKey gRPC method
func (k msgServer) RevokeAPIKey(ctx context.Context, msg *types.MsgRevokeAPIKey) (*types.MsgRevokeAPIKeyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing API key
	apiKey, found := k.GetAPIKey(sdkCtx, msg.ApiKeyId)
	if !found {
		return nil, types.ErrAPIKeyNotFound
	}

	// Verify ownership
	if apiKey.Owner != msg.Owner {
		return nil, types.ErrUnauthorized
	}

	// Delete
	k.DeleteAPIKey(sdkCtx, msg.ApiKeyId)

	return &types.MsgRevokeAPIKeyResponse{}, nil
}

// RecordUsage implements the RecordUsage gRPC method
func (k msgServer) RecordUsage(ctx context.Context, msg *types.MsgRecordUsage) (*types.MsgRecordUsageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Verify access
	apiKey, err := k.VerifyAPIKeyAccess(sdkCtx, msg.ApiKeyId, msg.ServiceId)
	if err != nil {
		return nil, err
	}

	// Record usage
	apiKey.RecordUsage()
	apiKey.LastUsedAt = sdkCtx.BlockTime().Unix()
	k.SetAPIKey(sdkCtx, apiKey)

	return &types.MsgRecordUsageResponse{}, nil
}

// RotateAPIKey implements the RotateAPIKey gRPC method
func (k msgServer) RotateAPIKey(ctx context.Context, msg *types.MsgRotateAPIKey) (*types.MsgRotateAPIKeyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Rotate the API key
	newAPIKeyID, err := k.Keeper.RotateAPIKey(sdkCtx, msg.Owner, msg.ApiKeyId, msg.NewEncryptedKey, msg.Reason)
	if err != nil {
		return nil, err
	}

	return &types.MsgRotateAPIKeyResponse{
		NewApiKeyId: newAPIKeyID,
	}, nil
}
