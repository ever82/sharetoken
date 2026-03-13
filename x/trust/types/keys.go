package types

const (
	// RouterKey is the message route for the trust module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the trust module
	QuerierRoute = ModuleName

	// StoreKey is the string store key for the trust module
	StoreKey = ModuleName
)

// Key prefixes for store
var (
	// MQScoreKey is the prefix for MQ score store
	MQScoreKey = []byte{0x01}

	// ParticipationKey is the prefix for participation history store
	ParticipationKey = []byte{0x02}
)

// GetMQScoreKey returns the key for an MQ score by address
func GetMQScoreKey(address string) []byte {
	return append(MQScoreKey, []byte(address)...)
}

// GetParticipationKey returns the key for participation history by address
func GetParticipationKey(address string) []byte {
	return append(ParticipationKey, []byte(address)...)
}
