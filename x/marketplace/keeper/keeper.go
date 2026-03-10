package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/marketplace/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{cdc: cdc, storeKey: storeKey}
}

func (k Keeper) SetService(ctx sdk.Context, service types.Service) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceKey(service.ID)
	value, err := json.Marshal(service)
	if err != nil {
		panic(err)
	}
	store.Set(key, value)
}

func (k Keeper) GetService(ctx sdk.Context, id string) (types.Service, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceKey(id)
	value := store.Get(key)
	if value == nil {
		return types.Service{}, false
	}
	var service types.Service
	if err := json.Unmarshal(value, &service); err != nil {
		return types.Service{}, false
	}
	return service, true
}
