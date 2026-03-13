package llmcustody

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/keeper"
)

// Event types and attributes
const (
	EventTypeRegisterAPIKey = "register_api_key"
	EventTypeUpdateAPIKey   = "update_api_key"
	EventTypeRevokeAPIKey   = "revoke_api_key"
	EventTypeRecordUsage    = "record_usage"

	AttributeKeyAPIKeyID     = "api_key_id"
	AttributeKeyOwner        = "owner"
	AttributeKeyProvider     = "provider"
	AttributeKeyActive       = "active"
	AttributeKeyServiceID    = "service_id"
	AttributeKeyRequestCount = "request_count"
	AttributeKeyTokenCount   = "token_count"
)

// NewHandler returns a handler for "llmcustody" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		_ = ctx.WithEventManager(sdk.NewEventManager())

		// For now, return an error indicating proto generation is needed
		// Once proto files are generated, the message types will implement sdk.Msg
		return nil, fmt.Errorf("llmcustody handler requires proto-generated message types. Please run 'make proto-gen' first")
	}
}
