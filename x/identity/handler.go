package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/keeper"
	"sharetoken/x/identity/types"
)

// NewHandler creates an sdk.Handler for the identity module
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(&k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRegisterIdentity:
			res, err := msgServer.RegisterIdentity(ctx.Context(), msg)
			if err != nil {
				return nil, err
			}
			return sdk.WrapServiceResult(ctx, res, nil)

		case *types.MsgVerifyIdentity:
			res, err := msgServer.VerifyIdentity(ctx.Context(), msg)
			if err != nil {
				return nil, err
			}
			return sdk.WrapServiceResult(ctx, res, nil)

		case *types.MsgUpdateLimitConfig:
			res, err := msgServer.UpdateLimitConfig(ctx.Context(), msg)
			if err != nil {
				return nil, err
			}
			return sdk.WrapServiceResult(ctx, res, nil)

		case *types.MsgResetDailyLimits:
			res, err := msgServer.ResetDailyLimits(ctx.Context(), msg)
			if err != nil {
				return nil, err
			}
			return sdk.WrapServiceResult(ctx, res, nil)

		default:
			return nil, fmt.Errorf("unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
