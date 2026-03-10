package types

const (
	// ModuleName is the name of the escrow module
	ModuleName = "escrow"

	// StoreKey is the string store key for the escrow module
	StoreKey = ModuleName

	// RouterKey is the message route for the escrow module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the escrow module
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// EscrowKey is the prefix for escrow store
	EscrowKey = []byte{0x01}

	// EscrowByRequesterKey is the prefix for escrow by requester index
	EscrowByRequesterKey = []byte{0x02}

	// EscrowByProviderKey is the prefix for escrow by provider index
	EscrowByProviderKey = []byte{0x03}
)

// GetEscrowKey returns the key for an escrow by ID
func GetEscrowKey(id string) []byte {
	return append(EscrowKey, []byte(id)...)
}

// GetEscrowByRequesterKey returns the key for escrow by requester index
func GetEscrowByRequesterKey(requester, id string) []byte {
	return append(append(EscrowByRequesterKey, []byte(requester)...), []byte(id)...)
}

// GetEscrowByProviderKey returns the key for escrow by provider index
func GetEscrowByProviderKey(provider, id string) []byte {
	return append(append(EscrowByProviderKey, []byte(provider)...), []byte(id)...)
}

// Event types
const (
	EventTypeCreateEscrow = "create_escrow"
	EventTypeRelease      = "release_escrow"
	EventTypeRefund       = "refund_escrow"
	EventTypeDispute      = "dispute_escrow"
	EventTypeResolve      = "resolve_escrow"
)

// Attribute keys
const (
	AttributeKeyEscrowID        = "escrow_id"
	AttributeKeyRequester       = "requester"
	AttributeKeyProvider        = "provider"
	AttributeKeyAmount          = "amount"
	AttributeKeyRequesterAmount = "requester_amount"
	AttributeKeyProviderAmount  = "provider_amount"
)
