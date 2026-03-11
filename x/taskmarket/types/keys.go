package types

const (
	// ModuleName is the name of the taskmarket module
	ModuleName = "taskmarket"
	StoreKey   = ModuleName
	// RouterKey is the message route for taskmarket module
	RouterKey = ModuleName
)

// Key prefixes for store
var (
	TaskKeyPrefix        = []byte{0x01}
	ApplicationKeyPrefix = []byte{0x02}
	AuctionKeyPrefix     = []byte{0x03}
	BidKeyPrefix         = []byte{0x04}
	RatingKeyPrefix      = []byte{0x05}
	ReputationKeyPrefix  = []byte{0x06}
)

// GetTaskKey returns the key for a task
func GetTaskKey(id string) []byte {
	return append(TaskKeyPrefix, []byte(id)...)
}

// GetApplicationKey returns the key for an application
func GetApplicationKey(id string) []byte {
	return append(ApplicationKeyPrefix, []byte(id)...)
}

// GetAuctionKey returns the key for an auction
func GetAuctionKey(taskID string) []byte {
	return append(AuctionKeyPrefix, []byte(taskID)...)
}

// GetBidKey returns the key for a bid
func GetBidKey(id string) []byte {
	return append(BidKeyPrefix, []byte(id)...)
}

// GetRatingKey returns the key for a rating
func GetRatingKey(id string) []byte {
	return append(RatingKeyPrefix, []byte(id)...)
}

// GetReputationKey returns the key for a reputation
func GetReputationKey(userID string) []byte {
	return append(ReputationKeyPrefix, []byte(userID)...)
}
