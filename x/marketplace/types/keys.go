package types

const (
	// ModuleName is the module name
	ModuleName = "marketplace"

	// RouterKey is the message route for the marketplace module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the marketplace module
	QuerierRoute = ModuleName

	// StoreKey is the string store key for the marketplace module
	StoreKey = ModuleName
)

var (
	ServiceKey = []byte{0x01}
)

func GetServiceKey(id string) []byte {
	return append(ServiceKey, []byte(id)...)
}
