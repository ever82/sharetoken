package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"

	"sharetoken/x/llmcustody/types"
)

// QueryClient provides a simple interface for querying the llmcustody module
// This is a temporary implementation until proto-generated query client is available
type QueryClient struct {
	clientCtx client.Context
}

// NewQueryClient creates a new QueryClient
func NewQueryClient(clientCtx client.Context) *QueryClient {
	return &QueryClient{clientCtx: clientCtx}
}

// APIKeyResponse wraps the API key response
type APIKeyResponse struct {
	APIKey types.APIKey `json:"api_key"`
}

// APIKeysResponse wraps the API keys response
type APIKeysResponse struct {
	APIKeys []types.APIKey `json:"api_keys"`
	Total   int            `json:"total"`
}

// UsageStatsResponse wraps the usage stats response
type UsageStatsResponse struct {
	Stats types.APIKeyUsageStats `json:"stats"`
}

// DailyUsageResponse wraps the daily usage response
type DailyUsageResponse struct {
	DailyStats []types.DailyUsageStats `json:"daily_stats"`
	Total      int                     `json:"total"`
}

// ServiceUsageResponse wraps the service usage response
type ServiceUsageResponse struct {
	Stats types.ServiceUsageStats `json:"stats"`
}

// String returns string representation
func (r *APIKeyResponse) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *APIKeysResponse) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *UsageStatsResponse) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *DailyUsageResponse) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *ServiceUsageResponse) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

// APIKey queries an API key by ID
func (qc *QueryClient) APIKey(ctx context.Context, id string) (*APIKeyResponse, error) {
	// For now, return a mock response
	// In full implementation, this would make a gRPC query
	fmt.Printf("Querying API key: %s\n", id)

	return &APIKeyResponse{
		APIKey: types.APIKey{
			ID:       id,
			Provider: types.ProviderOpenAI,
			Owner:    "",
			Active:   true,
		},
	}, nil
}

// APIKeysByOwner queries all API keys owned by an address
func (qc *QueryClient) APIKeysByOwner(ctx context.Context, owner string, page, limit int) (*APIKeysResponse, error) {
	fmt.Printf("Querying API keys for owner: %s (page: %d, limit: %d)\n", owner, page, limit)

	return &APIKeysResponse{
		APIKeys: []types.APIKey{},
		Total:   0,
	}, nil
}

// AllAPIKeys queries all registered API keys
func (qc *QueryClient) AllAPIKeys(ctx context.Context, page, limit int) (*APIKeysResponse, error) {
	fmt.Printf("Querying all API keys (page: %d, limit: %d)\n", page, limit)

	return &APIKeysResponse{
		APIKeys: []types.APIKey{},
		Total:   0,
	}, nil
}

// UsageStats queries usage statistics for an API key
func (qc *QueryClient) UsageStats(ctx context.Context, apiKeyID string) (*UsageStatsResponse, error) {
	fmt.Printf("Querying usage stats for API key: %s\n", apiKeyID)

	return &UsageStatsResponse{
		Stats: *types.NewAPIKeyUsageStats(apiKeyID),
	}, nil
}

// DailyUsage queries daily usage statistics
func (qc *QueryClient) DailyUsage(ctx context.Context, apiKeyID, date string) (*DailyUsageResponse, error) {
	fmt.Printf("Querying daily usage for API key: %s on date: %s\n", apiKeyID, date)

	return &DailyUsageResponse{
		DailyStats: []types.DailyUsageStats{
			*types.NewDailyUsageStats(date, apiKeyID),
		},
		Total: 1,
	}, nil
}

// ServiceUsage queries service usage statistics
func (qc *QueryClient) ServiceUsage(ctx context.Context, serviceID, apiKeyID string) (*ServiceUsageResponse, error) {
	fmt.Printf("Querying service usage for service: %s, API key: %s\n", serviceID, apiKeyID)

	return &ServiceUsageResponse{
		Stats: types.ServiceUsageStats{
			ServiceID: serviceID,
		},
	}, nil
}
