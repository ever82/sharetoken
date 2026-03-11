package taskmarket

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/keeper"
)

// NewHandler returns a handler for "taskmarket" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		// For now, return an error indicating proto generation is needed
		// Once proto files are generated, the message types will implement sdk.Msg
		return nil, fmt.Errorf("taskmarket handler requires proto-generated message types. Please run 'make proto-gen' first")
	}
}
