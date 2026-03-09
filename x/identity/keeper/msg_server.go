package keeper

import (
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
func (k msgServer) RegisterIdentity(ctx sdk.Context, msg *types.MsgRegisterIdentity) (*types.MsgRegisterIdentityResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err := k.Keeper.RegisterIdentity(ctx, msg.Address, msg.DID, msg.MetadataHash)
	if err != nil {
		return nil, err
	}

	// Get the created identity to return merkle root
	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		return nil, types.ErrIdentityNotFound.Wrap(msg.Address)
	}

	return &types.MsgRegisterIdentityResponse{
		MerkleRoot: identity.MerkleRoot,
	}, nil
}

// VerifyIdentity implements MsgServer
func (k msgServer) VerifyIdentity(ctx sdk.Context, msg *types.MsgVerifyIdentity) (*types.MsgVerifyIdentityResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err := k.Keeper.VerifyIdentity(ctx, msg.Address, msg.Provider, msg.VerificationHash, msg.Proof)
	if err != nil {
		return nil, err
	}

	// Get the updated identity to return merkle root
	identity, found := k.GetIdentity(ctx, msg.Address)
	if !found {
		return nil, types.ErrIdentityNotFound.Wrap(msg.Address)
	}

	return &types.MsgVerifyIdentityResponse{
		IsVerified:        identity.IsVerified,
		UpdatedMerkleRoot: identity.MerkleRoot,
	}, nil
}

// UpdateLimitConfig implements MsgServer
func (k msgServer) UpdateLimitConfig(ctx sdk.Context, msg *types.MsgUpdateLimitConfig) (*types.MsgUpdateLimitConfigResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if caller is authority
	caller := sdk.AccAddress(ctx.BlockHeader().ProposerAddress)
	if !k.IsAuthority(ctx, caller) {
		return nil, types.ErrUnauthorized.Wrap("only authority can update limits")
	}

	err := k.Keeper.UpdateLimitConfig(ctx, msg.TargetAddress, msg.NewConfig)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateLimitConfigResponse{}, nil
}

// ResetDailyLimits implements MsgServer
func (k msgServer) ResetDailyLimits(ctx sdk.Context, msg *types.MsgResetDailyLimits) (*types.MsgResetDailyLimitsResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if caller is authority
	caller := sdk.AccAddress(ctx.BlockHeader().ProposerAddress)
	if !k.IsAuthority(ctx, caller) {
		return nil, types.ErrUnauthorized.Wrap("only authority can reset limits")
	}

	resetCount := k.Keeper.ResetDailyLimits(ctx)

	return &types.MsgResetDailyLimitsResponse{
		ResetCount: resetCount,
	}, nil
}

// UpdateParams implements MsgServer (governance)
func (k msgServer) UpdateParams(ctx sdk.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Only governance can update params
	if msg.Authority != k.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrap("only governance can update params")
	}

	k.Keeper.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
