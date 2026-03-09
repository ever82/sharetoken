package types

const (
	// ModuleName is the name of the identity module
	ModuleName = "identity"

	// StoreKey is the string store key for the identity module
	StoreKey = ModuleName

	// RouterKey is the message route for the identity module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the identity module
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// IdentityKey is the prefix for identity store
	IdentityKey = []byte{0x01}

	// LimitConfigKey is the prefix for limit config store
	LimitConfigKey = []byte{0x02}

	// RegisteredDIDKey is the prefix for registered DID store
	RegisteredDIDKey = []byte{0x03}

	// VerificationProviderKey is the prefix for verification provider tracking
	VerificationProviderKey = []byte{0x04}
)

// GetIdentityKey returns the key for an identity by address
func GetIdentityKey(address string) []byte {
	return append(IdentityKey, []byte(address)...)
}

// GetLimitConfigKey returns the key for a limit config by address
func GetLimitConfigKey(address string) []byte {
	return append(LimitConfigKey, []byte(address)...)
}

// GetRegisteredDIDKey returns the key for a registered DID
func GetRegisteredDIDKey(did string) []byte {
	return append(RegisteredDIDKey, []byte(did)...)
}

// GetVerificationProviderKey returns the key for verification provider tracking
func GetVerificationProviderKey(provider, providerID string) []byte {
	return append(append(VerificationProviderKey, []byte(provider)...), []byte(providerID)...)
}
