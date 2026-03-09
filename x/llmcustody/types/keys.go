package types

const (
	// ModuleName is the name of the llmcustody module
	ModuleName = "llmcustody"
	StoreKey   = ModuleName
)

// APIKeyPrefix is the prefix for API key store
var APIKeyPrefix = []byte{0x01}

// GetAPIKeyKey returns the key for an API key
func GetAPIKeyKey(id string) []byte {
	return append(APIKeyPrefix, []byte(id)...)
}
