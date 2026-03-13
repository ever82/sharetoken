package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	identitytypes "sharetoken/x/identity/types"
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
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceKey(price.Symbol)
	value, err := k.cdc.Marshal(&price)
	if err != nil {
		return fmt.Errorf("failed to marshal price: %w", err)
	}
	store.Set(key, value)
	return nil
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
	if err := k.cdc.Unmarshal(value, &price); err != nil {
		return types.Price{}, false
	}
	return price, true
}

// GetAllPrices returns all prices
func (k Keeper) GetAllPrices(ctx sdk.Context) []types.Price {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PriceKey)
	defer func() {
		if err := iterator.Close(); err != nil {
			ctx.Logger().Error("failed to close iterator", "error", err)
		}
	}()

	var prices []types.Price
	for ; iterator.Valid(); iterator.Next() {
		var price types.Price
		if err := k.cdc.Unmarshal(iterator.Value(), &price); err != nil {
			ctx.Logger().Error("failed to unmarshal price", "error", err)
			continue
		}
		prices = append(prices, price)
	}

	return prices
}

// PriceFeedClient interface for external price feeds
type PriceFeedClient interface {
	FetchPrice(symbol string) (types.Price, error)
}

// ChainlinkClient implements Chainlink price feed client
type ChainlinkClient struct {
	endpoint string
	apiKey   string
}

// NewChainlinkClient creates a new Chainlink client
func NewChainlinkClient(endpoint, apiKey string) *ChainlinkClient {
	return &ChainlinkClient{
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

// FetchPrice fetches price from Chainlink
func (c *ChainlinkClient) FetchPrice(symbol string) (types.Price, error) {
	// This would call Chainlink API in production
	// For now, return mock data
	return types.Price{
		Symbol:     symbol,
		Price:      sdk.NewDec(identitytypes.DefaultSTTPrice),
		Source:     types.PriceSourceChainlink,
		Timestamp:  time.Now().Unix(),
		Confidence: types.PriceSourceConfidenceThreshold + 5,
	}, nil
}

// CalculateLLMPrice calculates the STT price for LLM service
func (k Keeper) CalculateLLMPrice(ctx sdk.Context, model string, inputTokens, outputTokens int64) (sdk.Coins, error) {
	// Get USD price for the model
	usdPricePer1K, err := k.GetLLMUSDPrice(model)
	if err != nil {
		return sdk.Coins{}, err
	}

	// Get STT/USD exchange rate
	sttPrice, found := k.GetPrice(ctx, "STT/USD")
	if !found {
		return sdk.Coins{}, types.ErrPriceNotFound
	}

	// Calculate USD cost
	totalTokens := inputTokens + outputTokens

	// Convert to STT using precise decimal arithmetic
	// usdPricePer1K is float64 (e.g., 0.03 for gpt-4)
	// sttPrice.Price is sdk.Dec (e.g., 0.001 for 1 STT = $0.001)
	// Formula: STT cost = (totalTokens / 1000) * usdPricePer1K / sttPriceInUSD

	// Create sdk.Dec from token count
	totalTokensDec := sdk.NewDec(totalTokens)
	usdPriceDec := sdk.NewDecFromIntWithPrec(sdk.NewInt(int64(usdPricePer1K*1e6)), 6)
	sttRate := sttPrice.Price

	// Calculate: (totalTokens / TokenUnitDivisor) * usdPricePer1K / sttRate
	// = totalTokens * usdPricePer1K / (TokenUnitDivisor * sttRate)
	sttCost := totalTokensDec.Mul(usdPriceDec).Quo(sttRate.MulInt64(types.TokenUnitDivisor))

	return sdk.NewCoins(sdk.NewCoin("ustt", sttCost.TruncateInt())), nil
}

// GetLLMUSDPrice gets the USD price per 1K tokens for a model
func (k Keeper) GetLLMUSDPrice(model string) (float64, error) {
	// Model pricing table (USD per 1K tokens)
	prices := map[string]float64{
		"gpt-4":           0.03,
		"gpt-4-turbo":     0.01,
		"gpt-3.5-turbo":   0.0015,
		"claude-3-opus":   0.015,
		"claude-3-sonnet": 0.003,
		"claude-3-haiku":  0.00025,
	}

	price, ok := prices[model]
	if !ok {
		return 0, fmt.Errorf("unknown model: %s", model)
	}
	return price, nil
}

// PriceAggregator aggregates prices from multiple sources
type PriceAggregator struct {
	sources []PriceFeedClient
}

// AggregatePrice aggregates prices and returns median
func (pa *PriceAggregator) AggregatePrice(symbol string) (types.Price, error) {
	var prices []sdk.Dec
	var priceSources []types.PriceSource

	for _, source := range pa.sources {
		price, err := source.FetchPrice(symbol)
		if err != nil {
			continue
		}
		prices = append(prices, price.Price) //nolint:staticcheck
		_ = priceSources
	}

	if len(prices) == 0 {
		return types.Price{}, fmt.Errorf("no price data available")
	}

	// Sort and take median
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].LT(prices[j])
	})
	median := prices[len(prices)/2]

	return types.Price{
		Symbol:     symbol,
		Price:      median,
		Source:     types.PriceSourceManual, // aggregated
		Timestamp:  time.Now().Unix(),
		Confidence: types.PriceSourceConfidenceThreshold,
	}, nil
}

// UpdatePrices updates all prices from external sources
func (k Keeper) UpdatePrices(ctx sdk.Context) error {
	// List of symbols to update
	symbols := []string{"STT/USD", "BTC/USD", "ETH/USD"}

	for _, symbol := range symbols {
		// In production, this would fetch from actual sources
		// For now, use mock data
		price := types.NewPrice(symbol, sdk.NewDec(1), types.PriceSourceManual, types.MaxConfidence)
		if err := k.SetPrice(ctx, *price); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePrice validates a price update
func (k Keeper) ValidatePrice(ctx sdk.Context, price types.Price) error {
	if price.Symbol == "" {
		return types.ErrInvalidPrice
	}

	if price.Price.IsNil() || price.Price.IsNegative() {
		return types.ErrInvalidPrice
	}

	// Check price is not too old (using escrow default duration as max age)
	if ctx.BlockTime().Unix()-price.Timestamp > int64(identitytypes.DefaultEscrowDurationHours)*identitytypes.SecondsPerHour {
		return types.ErrStalePrice
	}

	return nil
}

// GetPriceWithCache gets price with caching
func (k Keeper) GetPriceWithCache(ctx sdk.Context, symbol string, maxAge time.Duration) (types.Price, bool) {
	price, found := k.GetPrice(ctx, symbol)
	if !found {
		return types.Price{}, false
	}

	// Check if price is fresh
	if ctx.BlockTime().Unix()-price.Timestamp > int64(maxAge.Seconds()) {
		return types.Price{}, false
	}

	return price, true
}
