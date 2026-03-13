package types

const (
	// ModuleName is the name of the dispute module
	ModuleName = "dispute"

	// StoreKey is the string store key for the dispute module
	StoreKey = ModuleName

	// RouterKey is the message route for the dispute module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the dispute module
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// DisputeKey is the prefix for dispute store
	DisputeKey = []byte{0x01}

	// EvidenceKey is the prefix for evidence store
	EvidenceKey = []byte{0x02}

	// VoteKey is the prefix for vote store
	VoteKey = []byte{0x03}
)

// GetDisputeKey returns the key for a dispute by ID
func GetDisputeKey(id string) []byte {
	return append(DisputeKey, []byte(id)...)
}

// GetEvidenceKey returns the key for evidence by dispute ID
func GetEvidenceKey(disputeID string) []byte {
	return append(EvidenceKey, []byte(disputeID)...)
}

// GetVoteKey returns the key for votes by dispute ID
func GetVoteKey(disputeID string) []byte {
	return append(VoteKey, []byte(disputeID)...)
}
