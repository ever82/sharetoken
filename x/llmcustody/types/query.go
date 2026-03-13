package types

import (
	"fmt"
)

// QueryAPIKey defines a query function for querying an API key - placeholder for protobuf compatibility
func (req *QueryAPIKeyRequest) ValidateBasic() error {
	if req.Id == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation
func (req *QueryAPIKeysByOwnerRequest) ValidateBasic() error {
	if req.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation
func (req *QueryAllAPIKeysRequest) ValidateBasic() error {
	return nil
}

// ValidateBasic performs basic validation
func (req *QueryUsageStatsRequest) ValidateBasic() error {
	if req.ApiKeyId == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation
func (req *QueryDailyUsageRequest) ValidateBasic() error {
	if req.ApiKeyId == "" {
		return fmt.Errorf("API key ID cannot be empty")
	}
	if req.Date == "" {
		return fmt.Errorf("date cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation
func (req *QueryServiceUsageRequest) ValidateBasic() error {
	if req.ServiceId == "" {
		return fmt.Errorf("service ID cannot be empty")
	}
	return nil
}
