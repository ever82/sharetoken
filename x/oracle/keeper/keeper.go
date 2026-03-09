package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/oracle/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new oracle Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// SetPrice sets a price in the store
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceKey(price.Symbol)
	value, err := json.Marshal(price)
	if err != nil {
		panic(err)
	}
	store.Set(key, value)
}

// GetPrice retrieves a price by symbol
func (k Keeper) GetPrice(ctx sdk.Context, symbol string) (types.Price, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceKey(symbol)
	value := store.Get(key)

	if value == nil {
		return types.Price{}, false
	}

	var price types.Price
	if err := json.Unmarshal(value, &price); err != nil {
		return types.Price{}, false
	}
	return price, true
}

// GetAllPrices returns all prices
func (k Keeper) GetAllPrices(ctx sdk.Context) []types.Price {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PriceKey)
	defer iterator.Close()

	var prices []types.Price
	for ; iterator.Valid(); iterator.Next() {
		var price types.Price
		if err := json.Unmarshal(iterator.Value(), &price); err != nil {
			continue
		}
		prices = append(prices, price)
	}

	return prices
}
