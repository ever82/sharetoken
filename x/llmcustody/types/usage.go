package types

import (
	"fmt"
	"time"
)

// UsageRecord represents a single API usage record
type UsageRecord struct {
	ID            string    `json:"id"`
	APIKeyID      string    `json:"api_key_id"`
	ServiceID     string    `json:"service_id"`
	RequestCount  int64     `json:"request_count"`
	InputTokens   int64     `json:"input_tokens"`
	OutputTokens  int64     `json:"output_tokens"`
	TotalTokens   int64     `json:"total_tokens"`
	Cost          int64     `json:"cost"` // in ustt
	Timestamp     time.Time `json:"timestamp"`
	BlockHeight   int64     `json:"block_height"`
}

// NewUsageRecord creates a new usage record
func NewUsageRecord(apiKeyID, serviceID string, requests, inputTokens, outputTokens, cost int64, blockHeight int64) *UsageRecord {
	now := time.Now()
	return &UsageRecord{
		ID:           fmt.Sprintf("usage-%d-%s", now.UnixNano(), apiKeyID[:8]),
		APIKeyID:     apiKeyID,
		ServiceID:    serviceID,
		RequestCount: requests,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalTokens:  inputTokens + outputTokens,
		Cost:         cost,
		Timestamp:    now,
		BlockHeight:  blockHeight,
	}
}

// ValidateBasic performs basic validation
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

// APIKeyUsageStats represents aggregated usage statistics for an API key
type APIKeyUsageStats struct {
	APIKeyID          string    `json:"api_key_id"`
	TotalRequests     int64     `json:"total_requests"`
	TotalInputTokens  int64     `json:"total_input_tokens"`
	TotalOutputTokens int64     `json:"total_output_tokens"`
	TotalTokens       int64     `json:"total_tokens"`
	TotalCost         int64     `json:"total_cost"` // in ustt
	FirstUsedAt       time.Time `json:"first_used_at"`
	LastUsedAt        time.Time `json:"last_used_at"`
}

// NewAPIKeyUsageStats creates new API key usage statistics
func NewAPIKeyUsageStats(apiKeyID string) *APIKeyUsageStats {
	now := time.Now()
	return &APIKeyUsageStats{
		APIKeyID:    apiKeyID,
		FirstUsedAt: now,
		LastUsedAt:  now,
	}
}

// AddRecord adds a usage record to the statistics
func (stats *APIKeyUsageStats) AddRecord(record *UsageRecord) {
	stats.TotalRequests += record.RequestCount
	stats.TotalInputTokens += record.InputTokens
	stats.TotalOutputTokens += record.OutputTokens
	stats.TotalTokens += record.TotalTokens
	stats.TotalCost += record.Cost
	stats.LastUsedAt = record.Timestamp
}

// DailyUsageStats represents daily aggregated usage statistics
type DailyUsageStats struct {
	Date              string `json:"date"` // YYYY-MM-DD format
	APIKeyID          string `json:"api_key_id"`
	TotalRequests     int64  `json:"total_requests"`
	TotalInputTokens  int64  `json:"total_input_tokens"`
	TotalOutputTokens int64  `json:"total_output_tokens"`
	TotalTokens       int64  `json:"total_tokens"`
	TotalCost         int64  `json:"total_cost"`
}

// NewDailyUsageStats creates new daily usage statistics
func NewDailyUsageStats(date, apiKeyID string) *DailyUsageStats {
	return &DailyUsageStats{
		Date:     date,
		APIKeyID: apiKeyID,
	}
}

// AddRecord adds a usage record to daily statistics
func (ds *DailyUsageStats) AddRecord(record *UsageRecord) {
	ds.TotalRequests += record.RequestCount
	ds.TotalInputTokens += record.InputTokens
	ds.TotalOutputTokens += record.OutputTokens
	ds.TotalTokens += record.TotalTokens
	ds.TotalCost += record.Cost
}

// ServiceUsageStats represents usage statistics for a specific service
type ServiceUsageStats struct {
	ServiceID         string `json:"service_id"`
	TotalRequests     int64  `json:"total_requests"`
	TotalInputTokens  int64  `json:"total_input_tokens"`
	TotalOutputTokens int64  `json:"total_output_tokens"`
	TotalTokens       int64  `json:"total_tokens"`
	TotalCost         int64  `json:"total_cost"`
}

// KeyRotationHistory tracks API key rotation history
type KeyRotationHistory struct {
	APIKeyID      string    `json:"api_key_id"`
	PreviousKeyID string    `json:"previous_key_id"`
	NewKeyID      string    `json:"new_key_id"`
	RotatedAt     time.Time `json:"rotated_at"`
	RotatedBy     string    `json:"rotated_by"`
	Reason        string    `json:"reason"`
}

// NewKeyRotationHistory creates a new key rotation history record
func NewKeyRotationHistory(apiKeyID, previousKeyID, newKeyID, rotatedBy, reason string) *KeyRotationHistory {
	return &KeyRotationHistory{
		APIKeyID:      apiKeyID,
		PreviousKeyID: previousKeyID,
		NewKeyID:      newKeyID,
		RotatedAt:     time.Now(),
		RotatedBy:     rotatedBy,
		Reason:        reason,
	}
}
