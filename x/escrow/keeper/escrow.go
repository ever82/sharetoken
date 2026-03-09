package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/escrow/types"
)

// SetEscrow sets an escrow in the store
func (k Keeper) SetEscrow(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(escrow.ID)
	value, err := json.Marshal(escrow)
	if err != nil {
		panic(err)
	}
	store.Set(key, value)

	// Update indexes
	k.setEscrowByRequester(ctx, escrow)
	k.setEscrowByProvider(ctx, escrow)
}

// GetEscrow retrieves an escrow by ID
func (k Keeper) GetEscrow(ctx sdk.Context, id string) (types.Escrow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(id)
	value := store.Get(key)

	if value == nil {
		return types.Escrow{}, false
	}

	var escrow types.Escrow
	if err := json.Unmarshal(value, &escrow); err != nil {
		return types.Escrow{}, false
	}
	return escrow, true
}

// HasEscrow checks if an escrow exists
func (k Keeper) HasEscrow(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(id)
	return store.Has(key)
}

// DeleteEscrow deletes an escrow from the store
func (k Keeper) DeleteEscrow(ctx sdk.Context, id string) {
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(id)
	store.Delete(key)

	// Delete indexes
	k.deleteEscrowByRequester(ctx, escrow)
	k.deleteEscrowByProvider(ctx, escrow)
}

// setEscrowByRequester sets the escrow by requester index
func (k Keeper) setEscrowByRequester(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByRequesterKey(escrow.Requester, escrow.ID)
	store.Set(key, []byte{1})
}

// deleteEscrowByRequester deletes the escrow by requester index
func (k Keeper) deleteEscrowByRequester(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByRequesterKey(escrow.Requester, escrow.ID)
	store.Delete(key)
}

// setEscrowByProvider sets the escrow by provider index
func (k Keeper) setEscrowByProvider(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByProviderKey(escrow.Provider, escrow.ID)
	store.Set(key, []byte{1})
}

// deleteEscrowByProvider deletes the escrow by provider index
func (k Keeper) deleteEscrowByProvider(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByProviderKey(escrow.Provider, escrow.ID)
	store.Delete(key)
}

// GetEscrowsByRequester returns all escrows by requester
func (k Keeper) GetEscrowsByRequester(ctx sdk.Context, requester string) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.EscrowByRequesterKey, []byte(requester)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var escrows []types.Escrow
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract escrow ID from key
		escrowID := string(key[len(prefix):])
		escrow, found := k.GetEscrow(ctx, escrowID)
		if found {
			escrows = append(escrows, escrow)
		}
	}

	return escrows
}

// GetEscrowsByProvider returns all escrows by provider
func (k Keeper) GetEscrowsByProvider(ctx sdk.Context, provider string) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.EscrowByProviderKey, []byte(provider)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var escrows []types.Escrow
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract escrow ID from key
		escrowID := string(key[len(prefix):])
		escrow, found := k.GetEscrow(ctx, escrowID)
		if found {
			escrows = append(escrows, escrow)
		}
	}

	return escrows
}

// GetAllEscrows returns all escrows
func (k Keeper) GetAllEscrows(ctx sdk.Context) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.EscrowKey)
	defer iterator.Close()

	var escrows []types.Escrow
	for ; iterator.Valid(); iterator.Next() {
		var escrow types.Escrow
		if err := json.Unmarshal(iterator.Value(), &escrow); err != nil {
			continue
		}
		escrows = append(escrows, escrow)
	}

	return escrows
}
