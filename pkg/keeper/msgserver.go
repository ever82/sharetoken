package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgServer is a generic message server interface that can be implemented by any module.
// This provides a common pattern for message server implementations.
type MsgServer interface {
	// UnwrapSDKContext extracts the SDK context from the Go context
	// and performs any necessary validation
}

// UnwrapContext extracts the SDK context from the Go context.
// This is a helper function that centralizes the context unwrapping logic.
func UnwrapContext(ctx context.Context) sdk.Context {
	return sdk.UnwrapSDKContext(ctx)
}

// MsgHandler is a function type that handles a message and returns a response and error
type MsgHandler[T any, R any] func(sdk.Context, T) (R, error)

// HandleMsg is a generic helper for handling messages.
// It unwraps the context, calls the handler, and returns the result.
func HandleMsg[T any, R any](
	ctx context.Context,
	msg T,
	handler MsgHandler[T, R],
) (R, error) {
	sdkCtx := UnwrapContext(ctx)
	return handler(sdkCtx, msg)
}
