package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/types"
)

// Keeper of the llmcustody store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new llmcustody Keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// SetAPIKey sets an API key in the store
func (k Keeper) SetAPIKey(ctx sdk.Context, apiKey types.APIKey) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAPIKeyKey(apiKey.ID)
	value, err := k.cdc.Marshal(&apiKey)
	if err != nil {
		return fmt.Errorf("failed to marshal API key: %w", err)
	}
	store.Set(key, value)
	return nil
}

// GetAPIKey retrieves an API key by ID
func (k Keeper) GetAPIKey(ctx sdk.Context, id string) (types.APIKey, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAPIKeyKey(id)
	value := store.Get(key)

	if value == nil {
		return types.APIKey{}, false
	}

	var apiKey types.APIKey
	if err := k.cdc.Unmarshal(value, &apiKey); err != nil {
		return types.APIKey{}, false
	}
	return apiKey, true
}

// HasAPIKey checks if an API key exists
func (k Keeper) HasAPIKey(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAPIKeyKey(id)
	return store.Has(key)
}

// DeleteAPIKey deletes an API key
func (k Keeper) DeleteAPIKey(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAPIKeyKey(id)
	store.Delete(key)
}

// GetAllAPIKeys returns all API keys
func (k Keeper) GetAllAPIKeys(ctx sdk.Context) []types.APIKey {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.APIKeyPrefix)
	defer func() {
		if err := iterator.Close(); err != nil {
			ctx.Logger().Error("failed to close iterator", "error", err)
		}
	}()

	var apiKeys []types.APIKey
	for ; iterator.Valid(); iterator.Next() {
		var apiKey types.APIKey
		if err := k.cdc.Unmarshal(iterator.Value(), &apiKey); err != nil {
			ctx.Logger().Error("failed to unmarshal API key", "error", err)
			continue
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys
}

// GetAPIKeysByOwner returns all API keys owned by an address
func (k Keeper) GetAPIKeysByOwner(ctx sdk.Context, owner string) []types.APIKey {
	var result []types.APIKey
	allKeys := k.GetAllAPIKeys(ctx)
	for _, key := range allKeys {
		if key.Owner == owner {
			result = append(result, key)
		}
	}
	return result
}

// VerifyAPIKeyAccess verifies if an API key can be used for a service
func (k Keeper) VerifyAPIKeyAccess(ctx sdk.Context, keyID, serviceID string) (types.APIKey, error) {
	apiKey, found := k.GetAPIKey(ctx, keyID)
	if !found {
		return types.APIKey{}, fmt.Errorf("API key not found: %s", keyID)
	}

	if !apiKey.CanAccess(serviceID) {
		return types.APIKey{}, fmt.Errorf("API key %s cannot access service %s", keyID, serviceID)
	}

	return apiKey, nil
}
