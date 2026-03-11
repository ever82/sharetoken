package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/types"
)

// SetLimitConfig sets a limit config in the store
func (k Keeper) SetLimitConfig(ctx sdk.Context, limitConfig types.LimitConfig) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLimitConfigKey(limitConfig.Address)
	value, err := json.Marshal(limitConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal limit config: %w", err)
	}
	store.Set(key, value)
	return nil
}

// GetLimitConfig retrieves a limit config by address
func (k Keeper) GetLimitConfig(ctx sdk.Context, address string) (types.LimitConfig, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLimitConfigKey(address)
	value := store.Get(key)

	if value == nil {
		return types.LimitConfig{}, false
	}

	var limitConfig types.LimitConfig
	if err := json.Unmarshal(value, &limitConfig); err != nil {
		return types.LimitConfig{}, false
	}
	return limitConfig, true
}

// HasLimitConfig checks if a limit config exists
func (k Keeper) HasLimitConfig(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLimitConfigKey(address)
	return store.Has(key)
}

// DeleteLimitConfig deletes a limit config from the store
func (k Keeper) DeleteLimitConfig(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLimitConfigKey(address)
	store.Delete(key)
}

// GetOrCreateLimitConfig gets or creates a limit config for an address
func (k Keeper) GetOrCreateLimitConfig(ctx sdk.Context, address string) types.LimitConfig {
	limitConfig, found := k.GetLimitConfig(ctx, address)
	if found {
		return limitConfig
	}
	return types.NewLimitConfig(address)
}

// UpdateLimitConfig updates a user's limit configuration
func (k Keeper) UpdateLimitConfig(ctx sdk.Context, targetAddress string, newConfig types.LimitConfig) error {
	// Ensure the address matches
	if newConfig.Address != targetAddress {
		return types.ErrInvalidLimitConfig.Wrap("address mismatch")
	}

	// Validate the config
	if err := newConfig.ValidateBasic(); err != nil {
		return err
	}

	// Update timestamp
	newConfig.UpdatedAt = uint64(ctx.BlockTime().Unix())

	// Store the config
	if err := k.SetLimitConfig(ctx, newConfig); err != nil {
		return err
	}

	return nil
}

// CheckTransactionLimit checks if a transaction is within limits
func (k Keeper) CheckTransactionLimit(ctx sdk.Context, address string, amount sdk.Coin) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	return limitConfig.CheckTransactionLimit(amount)
}

// RecordTransaction records a transaction
func (k Keeper) RecordTransaction(ctx sdk.Context, address string, amount sdk.Coin) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	limitConfig.RecordTransaction(amount)
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}
	return nil
}

// CheckWithdrawalLimit checks if a withdrawal is within limits
func (k Keeper) CheckWithdrawalLimit(ctx sdk.Context, address string, amount sdk.Coin) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	return limitConfig.CheckWithdrawalLimit(amount)
}

// RecordWithdrawal records a withdrawal
func (k Keeper) RecordWithdrawal(ctx sdk.Context, address string, amount sdk.Coin) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	limitConfig.RecordWithdrawal(amount)
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}
	return nil
}

// CheckDisputeLimit checks if a new dispute can be created
func (k Keeper) CheckDisputeLimit(ctx sdk.Context, address string) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	return limitConfig.CheckDisputeLimit()
}

// IncrementActiveDisputes increments the active dispute count
func (k Keeper) IncrementActiveDisputes(ctx sdk.Context, address string) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	limitConfig.IncrementActiveDisputes()
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}
	return nil
}

// DecrementActiveDisputes decrements the active dispute count
func (k Keeper) DecrementActiveDisputes(ctx sdk.Context, address string) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	limitConfig.DecrementActiveDisputes()
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}
	return nil
}

// CheckServiceLimit checks if a service call is within limits
func (k Keeper) CheckServiceLimit(ctx sdk.Context, address string) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	return limitConfig.CheckServiceLimit()
}

// RecordServiceCall records a service call
func (k Keeper) RecordServiceCall(ctx sdk.Context, address string) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	limitConfig.RecordServiceCall()
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}
	return nil
}

// ReleaseServiceCall releases a service call slot
func (k Keeper) ReleaseServiceCall(ctx sdk.Context, address string) error {
	limitConfig := k.GetOrCreateLimitConfig(ctx, address)
	limitConfig.ReleaseServiceCall()
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}
	return nil
}

// ResetDailyLimits resets daily limits for all users
func (k Keeper) ResetDailyLimits(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LimitConfigKey)
	defer iterator.Close()

	var resetCount uint64
	for ; iterator.Valid(); iterator.Next() {
		var limitConfig types.LimitConfig
		if err := json.Unmarshal(iterator.Value(), &limitConfig); err != nil {
			continue
		}

		// Reset daily counters
		limitConfig.TxLimit.DailyTxCount = 0
		limitConfig.TxLimit.DailySpent = types.DefaultCoin()
		limitConfig.WithdrawalLimit.DailyWithdrawn = types.DefaultCoin()

		if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
			ctx.Logger().Error("failed to reset limit config", "address", limitConfig.Address, "error", err)
			continue
		}
		resetCount++
	}

	return resetCount
}

// GetAllLimitConfigs returns all limit configs
func (k Keeper) GetAllLimitConfigs(ctx sdk.Context) []types.LimitConfig {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LimitConfigKey)
	defer iterator.Close()

	var configs []types.LimitConfig
	for ; iterator.Valid(); iterator.Next() {
		var limitConfig types.LimitConfig
		if err := json.Unmarshal(iterator.Value(), &limitConfig); err != nil {
			ctx.Logger().Error("failed to unmarshal limit config", "error", err)
			continue
		}
		configs = append(configs, limitConfig)
	}

	return configs
}
