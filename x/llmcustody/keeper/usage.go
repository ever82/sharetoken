package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/types"
)

// RecordUsage records API usage with detailed tracking
func (k Keeper) RecordUsageWithStats(ctx sdk.Context, apiKeyID, serviceID string, requests, inputTokens, outputTokens, cost int64) error {
	// Get API key
	apiKey, found := k.GetAPIKey(ctx, apiKeyID)
	if !found {
		return types.ErrAPIKeyNotFound
	}

	// Check access
	if !apiKey.CanAccess(serviceID) {
		return types.ErrAccessDenied
	}

	// Create usage record
	record := types.NewUsageRecord(apiKeyID, serviceID, requests, inputTokens, outputTokens, cost, ctx.BlockHeight())

	// Store usage record
	if err := k.SetUsageRecord(ctx, *record); err != nil {
		return fmt.Errorf("failed to store usage record: %w", err)
	}

	// Update API key stats
	if err := k.updateAPIKeyStats(ctx, apiKeyID, record); err != nil {
		return fmt.Errorf("failed to update API key stats: %w", err)
	}

	// Update daily stats
	if err := k.updateDailyStats(ctx, apiKeyID, record); err != nil {
		return fmt.Errorf("failed to update daily stats: %w", err)
	}

	// Update service stats
	if err := k.updateServiceStats(ctx, serviceID, apiKeyID, record); err != nil {
		return fmt.Errorf("failed to update service stats: %w", err)
	}

	// Update API key usage count and last used time
	apiKey.RecordUsage()
	apiKey.LastUsedAt = ctx.BlockTime().Unix()
	if err := k.SetAPIKey(ctx, apiKey); err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	return nil
}

// SetUsageRecord stores a usage record
func (k Keeper) SetUsageRecord(ctx sdk.Context, record types.UsageRecord) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUsageRecordKey(record.ID)

	value, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal usage record: %w", err)
	}

	store.Set(key, value)
	return nil
}

// GetUsageRecord retrieves a usage record by ID
func (k Keeper) GetUsageRecord(ctx sdk.Context, id string) (types.UsageRecord, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUsageRecordKey(id)
	value := store.Get(key)

	if value == nil {
		return types.UsageRecord{}, false
	}

	var record types.UsageRecord
	if err := json.Unmarshal(value, &record); err != nil {
		return types.UsageRecord{}, false
	}

	return record, true
}

// updateAPIKeyStats updates aggregated statistics for an API key
func (k Keeper) updateAPIKeyStats(ctx sdk.Context, apiKeyID string, record *types.UsageRecord) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAPIKeyStatsKey(apiKeyID)

	var stats types.APIKeyUsageStats
	value := store.Get(key)
	if value != nil {
		if err := json.Unmarshal(value, &stats); err != nil {
			return fmt.Errorf("failed to unmarshal API key stats: %w", err)
		}
	} else {
		stats = *types.NewAPIKeyUsageStats(apiKeyID)
	}

	stats.AddRecord(record)

	value, err := json.Marshal(&stats)
	if err != nil {
		return fmt.Errorf("failed to marshal API key stats: %w", err)
	}

	store.Set(key, value)
	return nil
}

// GetAPIKeyStats retrieves usage statistics for an API key
func (k Keeper) GetAPIKeyStats(ctx sdk.Context, apiKeyID string) (types.APIKeyUsageStats, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAPIKeyStatsKey(apiKeyID)
	value := store.Get(key)

	if value == nil {
		return types.APIKeyUsageStats{}, false
	}

	var stats types.APIKeyUsageStats
	if err := json.Unmarshal(value, &stats); err != nil {
		return types.APIKeyUsageStats{}, false
	}

	return stats, true
}

// updateDailyStats updates daily usage statistics
func (k Keeper) updateDailyStats(ctx sdk.Context, apiKeyID string, record *types.UsageRecord) error {
	store := ctx.KVStore(k.storeKey)
	date := record.Timestamp.Format("2006-01-02")
	key := types.GetDailyStatsKey(date, apiKeyID)

	var stats types.DailyUsageStats
	value := store.Get(key)
	if value != nil {
		if err := json.Unmarshal(value, &stats); err != nil {
			return fmt.Errorf("failed to unmarshal daily stats: %w", err)
		}
	} else {
		stats = *types.NewDailyUsageStats(date, apiKeyID)
	}

	stats.AddRecord(record)

	value, err := json.Marshal(&stats)
	if err != nil {
		return fmt.Errorf("failed to marshal daily stats: %w", err)
	}

	store.Set(key, value)
	return nil
}

// GetDailyStats retrieves daily usage statistics
func (k Keeper) GetDailyStats(ctx sdk.Context, date, apiKeyID string) (types.DailyUsageStats, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDailyStatsKey(date, apiKeyID)
	value := store.Get(key)

	if value == nil {
		return types.DailyUsageStats{}, false
	}

	var stats types.DailyUsageStats
	if err := json.Unmarshal(value, &stats); err != nil {
		return types.DailyUsageStats{}, false
	}

	return stats, true
}

// updateServiceStats updates service usage statistics
func (k Keeper) updateServiceStats(ctx sdk.Context, serviceID, apiKeyID string, record *types.UsageRecord) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceStatsKey(serviceID, apiKeyID)

	var stats types.ServiceUsageStats
	value := store.Get(key)
	if value != nil {
		if err := json.Unmarshal(value, &stats); err != nil {
			return fmt.Errorf("failed to unmarshal service stats: %w", err)
		}
	} else {
		stats = types.ServiceUsageStats{
			ServiceID: serviceID,
		}
	}

	stats.TotalRequests += record.RequestCount
	stats.TotalInputTokens += record.InputTokens
	stats.TotalOutputTokens += record.OutputTokens
	stats.TotalTokens += record.TotalTokens
	stats.TotalCost += record.Cost

	value, err := json.Marshal(&stats)
	if err != nil {
		return fmt.Errorf("failed to marshal service stats: %w", err)
	}

	store.Set(key, value)
	return nil
}

// GetServiceStats retrieves service usage statistics
func (k Keeper) GetServiceStats(ctx sdk.Context, serviceID, apiKeyID string) (types.ServiceUsageStats, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceStatsKey(serviceID, apiKeyID)
	value := store.Get(key)

	if value == nil {
		return types.ServiceUsageStats{}, false
	}

	var stats types.ServiceUsageStats
	if err := json.Unmarshal(value, &stats); err != nil {
		return types.ServiceUsageStats{}, false
	}

	return stats, true
}

// RotateAPIKey rotates an API key with a new encrypted key
func (k Keeper) RotateAPIKey(ctx sdk.Context, owner, apiKeyID string, newEncryptedKey []byte, reason string) (string, error) {
	// Get existing API key
	apiKey, found := k.GetAPIKey(ctx, apiKeyID)
	if !found {
		return "", types.ErrAPIKeyNotFound
	}

	// Verify ownership
	if apiKey.Owner != owner {
		return "", types.ErrUnauthorized
	}

	// Store old key ID and encrypted key
	oldKeyID := apiKey.ID
	oldEncryptedKey := make([]byte, len(apiKey.EncryptedKey))
	copy(oldEncryptedKey, apiKey.EncryptedKey)

	// Generate new API key ID
	newAPIKeyID := types.GenerateAPIKeyID()

	// Create new API key with same configuration but new encrypted key
	newAPIKey := types.NewAPIKey(
		newAPIKeyID,
		apiKey.Provider,
		newEncryptedKey,
		owner,
	)
	newAPIKey.AccessRules = apiKey.AccessRules
	newAPIKey.Active = apiKey.Active
	newAPIKey.CreatedAt = ctx.BlockTime().Unix()
	newAPIKey.LastUsedAt = apiKey.LastUsedAt
	newAPIKey.UsageCount = apiKey.UsageCount
	newAPIKey.Version = apiKey.Version + 1

	// Store new API key
	if err := k.SetAPIKey(ctx, *newAPIKey); err != nil {
		return "", fmt.Errorf("failed to store new API key: %w", err)
	}

	// Mark old key as inactive (do not delete to preserve history)
	apiKey.Active = false
	if err := k.SetAPIKey(ctx, apiKey); err != nil {
		return "", fmt.Errorf("failed to update old API key: %w", err)
	}

	// Store rotation history
	rotation := types.NewKeyRotationHistory(apiKeyID, oldKeyID, newAPIKeyID, owner, reason)
	if err := k.SetKeyRotationHistory(ctx, *rotation); err != nil {
		return "", fmt.Errorf("failed to store rotation history: %w", err)
	}

	// Securely wipe old encrypted key
	apiKey.SecureWipe()
	types.Zeroize(oldEncryptedKey)

	return newAPIKeyID, nil
}

// SetKeyRotationHistory stores key rotation history
func (k Keeper) SetKeyRotationHistory(ctx sdk.Context, rotation types.KeyRotationHistory) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKeyRotationKey(rotation.APIKeyID)

	value, err := json.Marshal(&rotation)
	if err != nil {
		return fmt.Errorf("failed to marshal rotation history: %w", err)
	}

	store.Set(key, value)
	return nil
}

// GetKeyRotationHistory retrieves key rotation history
func (k Keeper) GetKeyRotationHistory(ctx sdk.Context, apiKeyID string) (types.KeyRotationHistory, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKeyRotationKey(apiKeyID)
	value := store.Get(key)

	if value == nil {
		return types.KeyRotationHistory{}, false
	}

	var rotation types.KeyRotationHistory
	if err := json.Unmarshal(value, &rotation); err != nil {
		return types.KeyRotationHistory{}, false
	}

	return rotation, true
}
