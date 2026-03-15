package types

import (
	"fmt"
	"time"
)

// KeyRotationHistory represents the history of API key rotations
type KeyRotationHistory struct {
	APIKeyID        string `json:"api_key_id"`
	PreviousKeyHash string `json:"previous_key_hash"`
	NewKeyID        string `json:"new_key_id"`
	Reason          string `json:"reason"`
	RotatedAt       int64  `json:"rotated_at"`
}

// UsageRecord represents a single API usage record
// Note: This is kept for backward compatibility - query.pb.go has similar types
type UsageRecord struct {
	ID            string `json:"id"`
	APIKeyID      string `json:"api_key_id"`
	ServiceID     string `json:"service_id"`
	RequestCount  int64  `json:"request_count"`
	InputTokens   int64  `json:"input_tokens"`
	OutputTokens  int64  `json:"output_tokens"`
	TotalTokens   int64  `json:"total_tokens"`
	Cost          int64  `json:"cost"` // in ustt
	Timestamp     int64  `json:"timestamp"`
	BlockHeight   int64  `json:"block_height"`
}

// NewAPIKeyUsageStats creates a new API key usage stats
func NewAPIKeyUsageStats(apiKeyID string) *APIKeyUsageStats {
	return &APIKeyUsageStats{
		ApiKeyId:          apiKeyID,
		TotalRequests:     0,
		TotalInputTokens:  0,
		TotalOutputTokens: 0,
		TotalTokens:       0,
		TotalCost:         0,
		LastUsedAt:        0,
	}
}

// NewDailyUsageStats creates a new daily usage stats
func NewDailyUsageStats(date, apiKeyID string) *DailyUsageStats {
	return &DailyUsageStats{
		Date:              date,
		ApiKeyId:          apiKeyID,
		TotalRequests:     0,
		TotalInputTokens:  0,
		TotalOutputTokens: 0,
		TotalTokens:       0,
		TotalCost:         0,
	}
}

// NewUsageRecord creates a new usage record
func NewUsageRecord(apiKeyID, serviceID string, requests, inputTokens, outputTokens, cost, blockHeight int64) *UsageRecord {
	return &UsageRecord{
		ID:            GenerateUsageRecordID(),
		APIKeyID:      apiKeyID,
		ServiceID:     serviceID,
		RequestCount:  requests,
		InputTokens:   inputTokens,
		OutputTokens:  outputTokens,
		TotalTokens:   inputTokens + outputTokens,
		Cost:          cost,
		Timestamp:     0, // Will be set by caller
		BlockHeight:   blockHeight,
	}
}

// GenerateUsageRecordID generates a unique usage record ID
func GenerateUsageRecordID() string {
	return fmt.Sprintf("usage_%d", time.Now().UnixNano())
}

// NewKeyRotationHistory creates a new key rotation history entry
func NewKeyRotationHistory(apiKeyID, previousKeyHash, newKeyID, reason string) *KeyRotationHistory {
	return &KeyRotationHistory{
		APIKeyID:        apiKeyID,
		PreviousKeyHash: previousKeyHash,
		NewKeyID:        newKeyID,
		Reason:          reason,
		RotatedAt:       time.Now().Unix(),
	}
}
func (ur UsageRecord) ValidateBasic() error {
	if ur.APIKeyID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if ur.ServiceID == "" {
		return fmt.Errorf("service ID cannot be empty")
	}
	if ur.RequestCount < 0 {
		return fmt.Errorf("request count cannot be negative")
	}
	if ur.InputTokens < 0 {
		return fmt.Errorf("input tokens cannot be negative")
	}
	if ur.OutputTokens < 0 {
		return fmt.Errorf("output tokens cannot be negative")
	}
	if ur.Cost < 0 {
		return fmt.Errorf("cost cannot be negative")
	}
	return nil
}

// AddRecord adds a usage record to daily statistics
func (ds *DailyUsageStats) AddRecord(record UsageRecord) {
	ds.TotalRequests += record.RequestCount
	ds.TotalInputTokens += record.InputTokens
	ds.TotalOutputTokens += record.OutputTokens
	ds.TotalTokens += record.TotalTokens
	ds.TotalCost += record.Cost
}

// AddRecord adds a usage record to the statistics
func (stats *APIKeyUsageStats) AddRecord(record UsageRecord) {
	stats.TotalRequests += record.RequestCount
	stats.TotalInputTokens += record.InputTokens
	stats.TotalOutputTokens += record.OutputTokens
	stats.TotalTokens += record.TotalTokens
	stats.TotalCost += record.Cost
	stats.LastUsedAt = record.Timestamp
}
