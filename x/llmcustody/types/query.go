package types

import (
	"context"
	"fmt"
)

// QueryServer is the query server interface
type QueryServer interface {
	// APIKey queries an API key by ID
	APIKey(ctx context.Context, req *QueryAPIKeyRequest) (*QueryAPIKeyResponse, error)

	// APIKeysByOwner queries all API keys owned by an address
	APIKeysByOwner(ctx context.Context, req *QueryAPIKeysByOwnerRequest) (*QueryAPIKeysByOwnerResponse, error)

	// AllAPIKeys queries all registered API keys
	AllAPIKeys(ctx context.Context, req *QueryAllAPIKeysRequest) (*QueryAllAPIKeysResponse, error)

	// UsageStats queries usage statistics for an API key
	UsageStats(ctx context.Context, req *QueryUsageStatsRequest) (*QueryUsageStatsResponse, error)

	// DailyUsage queries daily usage statistics
	DailyUsage(ctx context.Context, req *QueryDailyUsageRequest) (*QueryDailyUsageResponse, error)

	// ServiceUsage queries service usage statistics
	ServiceUsage(ctx context.Context, req *QueryServiceUsageRequest) (*QueryServiceUsageResponse, error)
}

// -----------------------------------------------------------------------------
// QueryAPIKey
// -----------------------------------------------------------------------------

// QueryAPIKeyRequest is the request for APIKey query
type QueryAPIKeyRequest struct {
	ID string `json:"id"`
}

// QueryAPIKeyResponse is the response for APIKey query
type QueryAPIKeyResponse struct {
	APIKey APIKey `json:"api_key"`
}

// ValidateBasic performs basic validation
func (req QueryAPIKeyRequest) ValidateBasic() error {
	if req.ID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// QueryAPIKeysByOwner
// -----------------------------------------------------------------------------

// QueryAPIKeysByOwnerRequest is the request for APIKeysByOwner query
type QueryAPIKeysByOwnerRequest struct {
	Owner  string `json:"owner"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// QueryAPIKeysByOwnerResponse is the response for APIKeysByOwner query
type QueryAPIKeysByOwnerResponse struct {
	APIKeys []APIKey `json:"api_keys"`
	Total   int      `json:"total"`
}

// ValidateBasic performs basic validation
func (req QueryAPIKeysByOwnerRequest) ValidateBasic() error {
	if req.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if req.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	return nil
}

// -----------------------------------------------------------------------------
// QueryAllAPIKeys
// -----------------------------------------------------------------------------

// QueryAllAPIKeysRequest is the request for AllAPIKeys query
type QueryAllAPIKeysRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// QueryAllAPIKeysResponse is the response for AllAPIKeys query
type QueryAllAPIKeysResponse struct {
	APIKeys []APIKey `json:"api_keys"`
	Total   int      `json:"total"`
}

// ValidateBasic performs basic validation
func (req QueryAllAPIKeysRequest) ValidateBasic() error {
	if req.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	return nil
}

// -----------------------------------------------------------------------------
// QueryUsageStats
// -----------------------------------------------------------------------------

// QueryUsageStatsRequest is the request for UsageStats query
type QueryUsageStatsRequest struct {
	APIKeyID string `json:"api_key_id"`
}

// QueryUsageStatsResponse is the response for UsageStats query
type QueryUsageStatsResponse struct {
	Stats APIKeyUsageStats `json:"stats"`
}

// ValidateBasic performs basic validation
func (req QueryUsageStatsRequest) ValidateBasic() error {
	if req.APIKeyID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// QueryDailyUsage
// -----------------------------------------------------------------------------

// QueryDailyUsageRequest is the request for DailyUsage query
type QueryDailyUsageRequest struct {
	APIKeyID string `json:"api_key_id"`
	Date     string `json:"date"` // YYYY-MM-DD format
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

// QueryDailyUsageResponse is the response for DailyUsage query
type QueryDailyUsageResponse struct {
	DailyStats []DailyUsageStats `json:"daily_stats"`
	Total      int               `json:"total"`
}

// ValidateBasic performs basic validation
func (req QueryDailyUsageRequest) ValidateBasic() error {
	if req.APIKeyID == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if req.Date == "" {
		return fmt.Errorf("date cannot be empty")
	}
	if req.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	return nil
}

// -----------------------------------------------------------------------------
// QueryServiceUsage
// -----------------------------------------------------------------------------

// QueryServiceUsageRequest is the request for ServiceUsage query
type QueryServiceUsageRequest struct {
	ServiceID string `json:"service_id"`
	APIKeyID  string `json:"api_key_id"`
}

// QueryServiceUsageResponse is the response for ServiceUsage query
type QueryServiceUsageResponse struct {
	Stats ServiceUsageStats `json:"stats"`
}

// ValidateBasic performs basic validation
func (req QueryServiceUsageRequest) ValidateBasic() error {
	if req.ServiceID == "" {
		return fmt.Errorf("service ID cannot be empty")
	}
	return nil
}
