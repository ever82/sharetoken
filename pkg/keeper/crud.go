package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KVStoreItem is an interface that must be implemented by types
// that can be stored in the KV store using the generic CRUD operations.
type KVStoreItem interface {
	// GetID returns the unique identifier for this item
	GetID() string
}

// StoreAccessor is a function type that returns a KV store for the given context
type StoreAccessor func(ctx sdk.Context) sdk.KVStore

// CRUDKeeper provides generic CRUD operations for any type that implements KVStoreItem.
// This reduces boilerplate code in keepers that perform similar operations.
type CRUDKeeper[T KVStoreItem] struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	prefix   []byte
}

// NewCRUDKeeper creates a new CRUD keeper for the given type
func NewCRUDKeeper[T KVStoreItem](cdc codec.BinaryCodec, storeKey storetypes.StoreKey, prefix []byte) *CRUDKeeper[T] {
	return &CRUDKeeper[T]{
		cdc:      cdc,
		storeKey: storeKey,
		prefix:   prefix,
	}
}

// Set stores an item in the KV store
func (k *CRUDKeeper[T]) Set(ctx sdk.Context, item T) {
	store := ctx.KVStore(k.storeKey)
	key := BuildKey(k.prefix, item.GetID())
	value, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	store.Set(key, value)
}

// Get retrieves an item by ID from the KV store
func (k *CRUDKeeper[T]) Get(ctx sdk.Context, id string) (T, bool) {
	var result T
	store := ctx.KVStore(k.storeKey)
	key := BuildKey(k.prefix, id)
	value := store.Get(key)

	if value == nil {
		return result, false
	}

	if err := json.Unmarshal(value, &result); err != nil {
		return result, false
	}
	return result, true
}

// Has checks if an item exists by ID
func (k *CRUDKeeper[T]) Has(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	key := BuildKey(k.prefix, id)
	return store.Has(key)
}

// Delete removes an item by ID from the KV store
func (k *CRUDKeeper[T]) Delete(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	key := BuildKey(k.prefix, id)
	store.Delete(key)
}

// GetAll retrieves all items with the given prefix
func (k *CRUDKeeper[T]) GetAll(ctx sdk.Context) []T {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, k.prefix)
	defer iterator.Close() //nolint:errcheck

	var items []T
	for ; iterator.Valid(); iterator.Next() {
		var item T
		if err := json.Unmarshal(iterator.Value(), &item); err != nil {
			continue
		}
		items = append(items, item)
	}

	return items
}

// BuildKey builds a store key from prefix and ID
func BuildKey(prefix []byte, id string) []byte {
	return append(prefix, []byte(id)...)
}

// BuildCompositeKey builds a store key from prefix and multiple components
func BuildCompositeKey(prefix []byte, components ...string) []byte {
	key := prefix
	for _, component := range components {
		key = append(key, []byte(component)...)
	}
	return key
}
