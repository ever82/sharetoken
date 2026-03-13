package types

const (
	// ModuleName is the name of the agentgateway module
	ModuleName = "agentgateway"

	// StoreKey is the string store key for the agentgateway module
	StoreKey = ModuleName

	// RouterKey is the message route for the agentgateway module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the agentgateway module
	QuerierRoute = ModuleName
)

// Rate limit constants
const (
	// DefaultRateLimitPerMinute is the default maximum requests per minute per address
	DefaultRateLimitPerMinute = 60
)

// Key prefixes for store
var (
	// SessionKey is the prefix for session store
	SessionKey = []byte{0x01}

	// RateLimitKey is the prefix for rate limit store
	RateLimitKey = []byte{0x02}
)

// GetSessionKey returns the key for a session by ID
func GetSessionKey(id string) []byte {
	return append(SessionKey, []byte(id)...)
}

// GetRateLimitKey returns the key for rate limit by address
func GetRateLimitKey(address string) []byte {
	return append(RateLimitKey, []byte(address)...)
}
