package keeper

import (
	"fmt"

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

func (k Keeper) SetService(ctx sdk.Context, service types.Service) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceKey(service.ID)
	value, err := k.cdc.Marshal(&service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}
	store.Set(key, value)
	return nil
}

func (k Keeper) GetService(ctx sdk.Context, id string) (types.Service, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceKey(id)
	value := store.Get(key)
	if value == nil {
		return types.Service{}, false
	}
	var service types.Service
	if err := k.cdc.Unmarshal(value, &service); err != nil {
		return types.Service{}, false
	}
	return service, true
}

// GetAllServices returns all services from the store
func (k Keeper) GetAllServices(ctx sdk.Context) []types.Service {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ServiceKey)
	defer iterator.Close()

	var services []types.Service
	for ; iterator.Valid(); iterator.Next() {
		var service types.Service
		if err := k.cdc.Unmarshal(iterator.Value(), &service); err != nil {
			ctx.Logger().Error("failed to unmarshal service", "error", err)
			continue
		}
		services = append(services, service)
	}

	return services
}
