package types

const (
	// ModuleName is the name of the node module
	ModuleName = "node"

	// StoreKey is the string store key for the node module
	StoreKey = ModuleName

	// RouterKey is the message route for the node module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the node module
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// NodeConfigKey is the prefix for node configuration store
	NodeConfigKey = []byte{0x01}

	// NodeRoleKey is the prefix for node role store
	NodeRoleKey = []byte{0x02}
)

// GetNodeConfigKey returns the key for node configuration
func GetNodeConfigKey(address string) []byte {
	return append(NodeConfigKey, []byte(address)...)
}

// GetNodeRoleKey returns the key for node role
func GetNodeRoleKey(address string) []byte {
	return append(NodeRoleKey, []byte(address)...)
}
