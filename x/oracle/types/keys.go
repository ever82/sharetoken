package types

const (
	// ModuleName is the name of the oracle module
	ModuleName = "oracle"

	// StoreKey is the string store key
	StoreKey = ModuleName
)

// PriceKey is the prefix for price store
var PriceKey = []byte{0x01}

// GetPriceKey returns the key for a price
func GetPriceKey(symbol string) []byte {
	return append(PriceKey, []byte(symbol)...)
}
