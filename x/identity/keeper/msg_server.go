package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/types"
)

// msgServer is the message server
type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = &msgServer{}

// RegisterIdentity implements MsgServer
func (k msgServer) RegisterIdentity(ctx context.Context, msg *types.MsgRegisterIdentity) (*types.MsgRegisterIdentityResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.Keeper.RegisterIdentity(sdkCtx, msg.Address, msg.Did, msg.MetadataHash)
	if err != nil {
		return nil, err
	}

	// Get the created identity to return merkle root
	identity, found := k.GetIdentity(sdkCtx, msg.Address)
	if !found {
		return nil, types.ErrIdentityNotFound.Wrap(msg.Address)
	}

	return &types.MsgRegisterIdentityResponse{
		MerkleRoot: identity.MerkleRoot,
	}, nil
}

// VerifyIdentity implements MsgServer
func (k msgServer) VerifyIdentity(ctx context.Context, msg *types.MsgVerifyIdentity) (*types.MsgVerifyIdentityResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.Keeper.VerifyIdentity(sdkCtx, msg.Address, msg.Provider, msg.VerificationHash, msg.Proof)
	if err != nil {
		return nil, err
	}

	// Get the updated identity to return merkle root
	identity, found := k.GetIdentity(sdkCtx, msg.Address)
	if !found {
		return nil, types.ErrIdentityNotFound.Wrap(msg.Address)
	}

	return &types.MsgVerifyIdentityResponse{
		IsVerified:        identity.IsVerified,
		UpdatedMerkleRoot: identity.MerkleRoot,
	}, nil
}

// UpdateLimitConfig implements MsgServer
func (k msgServer) UpdateLimitConfig(ctx context.Context, msg *types.MsgUpdateLimitConfig) (*types.MsgUpdateLimitConfigResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if caller is authority
	caller := sdk.AccAddress(sdkCtx.BlockHeader().ProposerAddress)
	if !k.IsAuthority(sdkCtx, caller) {
		return nil, types.ErrUnauthorized.Wrap("only authority can update limits")
	}

	err := k.Keeper.UpdateLimitConfig(sdkCtx, msg.TargetAddress, msg.NewConfig)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateLimitConfigResponse{}, nil
}

// ResetDailyLimits implements MsgServer
func (k msgServer) ResetDailyLimits(ctx context.Context, msg *types.MsgResetDailyLimits) (*types.MsgResetDailyLimitsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if caller is authority
	caller := sdk.AccAddress(sdkCtx.BlockHeader().ProposerAddress)
	if !k.IsAuthority(sdkCtx, caller) {
		return nil, types.ErrUnauthorized.Wrap("only authority can reset limits")
	}

	resetCount := k.Keeper.ResetDailyLimits(sdkCtx)

	return &types.MsgResetDailyLimitsResponse{
		ResetCount: resetCount,
	}, nil
}

// UpdateParams implements MsgServer (governance)
func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Only governance can update params
	if msg.Authority != k.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrap("only governance can update params")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.Keeper.SetParams(sdkCtx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
