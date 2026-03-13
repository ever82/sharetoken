package types

const (
	// ModuleName is the name of the llmcustody module
	ModuleName = "llmcustody"
	StoreKey   = ModuleName
	// RouterKey is the message route for llmcustody module
	RouterKey = ModuleName
)

// Store key prefixes
var (
	// APIKeyPrefix is the prefix for API key store
	APIKeyPrefix = []byte{0x01}

	// UsageRecordPrefix is the prefix for usage records
	UsageRecordPrefix = []byte{0x02}

	// APIKeyStatsPrefix is the prefix for API key usage statistics
	APIKeyStatsPrefix = []byte{0x03}

	// DailyStatsPrefix is the prefix for daily usage statistics
	DailyStatsPrefix = []byte{0x04}

	// ServiceStatsPrefix is the prefix for service usage statistics
	ServiceStatsPrefix = []byte{0x05}

	// KeyRotationPrefix is the prefix for key rotation history
	KeyRotationPrefix = []byte{0x06}
)

// GetAPIKeyKey returns the key for an API key
func GetAPIKeyKey(id string) []byte {
	return append(APIKeyPrefix, []byte(id)...)
}

// GetUsageRecordKey returns the key for a usage record
func GetUsageRecordKey(id string) []byte {
	return append(UsageRecordPrefix, []byte(id)...)
}

// GetAPIKeyStatsKey returns the key for API key usage statistics
func GetAPIKeyStatsKey(apiKeyID string) []byte {
	return append(APIKeyStatsPrefix, []byte(apiKeyID)...)
}

// GetDailyStatsKey returns the key for daily usage statistics
func GetDailyStatsKey(date, apiKeyID string) []byte {
	key := append(DailyStatsPrefix, []byte(date)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(apiKeyID)...)
}

// GetServiceStatsKey returns the key for service usage statistics
func GetServiceStatsKey(serviceID, apiKeyID string) []byte {
	key := append(ServiceStatsPrefix, []byte(serviceID)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(apiKeyID)...)
}

// GetKeyRotationKey returns the key for key rotation history
func GetKeyRotationKey(apiKeyID string) []byte {
	return append(KeyRotationPrefix, []byte(apiKeyID)...)
}
