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

	// Index prefixes
	TaskByRequesterPrefix   = []byte{0x11}
	TaskByWorkerPrefix      = []byte{0x12}
	TaskByStatusPrefix      = []byte{0x13}
	ApplicationByTaskPrefix = []byte{0x21}
	BidByTaskPrefix         = []byte{0x31}
	RatingByTaskPrefix      = []byte{0x41}
	RatingByRatedUserPrefix = []byte{0x42}
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

// Index keys

// GetTaskByRequesterKey returns the key for task by requester index
func GetTaskByRequesterKey(requesterID, taskID string) []byte {
	key := append(TaskByRequesterPrefix, []byte(requesterID)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(taskID)...)
}

// GetTaskByRequesterPrefix returns the prefix for iterating tasks by requester
func GetTaskByRequesterPrefix(requesterID string) []byte {
	key := append(TaskByRequesterPrefix, []byte(requesterID)...)
	return append(key, 0x00)
}

// GetTaskByWorkerKey returns the key for task by worker index
func GetTaskByWorkerKey(workerID, taskID string) []byte {
	key := append(TaskByWorkerPrefix, []byte(workerID)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(taskID)...)
}

// GetTaskByWorkerPrefix returns the prefix for iterating tasks by worker
func GetTaskByWorkerPrefix(workerID string) []byte {
	key := append(TaskByWorkerPrefix, []byte(workerID)...)
	return append(key, 0x00)
}

// GetTaskByStatusKey returns the key for task by status index
func GetTaskByStatusKey(status, taskID string) []byte {
	key := append(TaskByStatusPrefix, []byte(status)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(taskID)...)
}

// GetTaskByStatusPrefix returns the prefix for iterating tasks by status
func GetTaskByStatusPrefix(status string) []byte {
	key := append(TaskByStatusPrefix, []byte(status)...)
	return append(key, 0x00)
}

// GetApplicationByTaskKey returns the key for application by task index
func GetApplicationByTaskKey(taskID, appID string) []byte {
	key := append(ApplicationByTaskPrefix, []byte(taskID)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(appID)...)
}

// GetApplicationByTaskPrefix returns the prefix for iterating applications by task
func GetApplicationByTaskPrefix(taskID string) []byte {
	key := append(ApplicationByTaskPrefix, []byte(taskID)...)
	return append(key, 0x00)
}

// GetBidByTaskKey returns the key for bid by task index
func GetBidByTaskKey(taskID, bidID string) []byte {
	key := append(BidByTaskPrefix, []byte(taskID)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(bidID)...)
}

// GetBidByTaskPrefix returns the prefix for iterating bids by task
func GetBidByTaskPrefix(taskID string) []byte {
	key := append(BidByTaskPrefix, []byte(taskID)...)
	return append(key, 0x00)
}

// GetRatingByTaskKey returns the key for rating by task index
func GetRatingByTaskKey(taskID, ratingID string) []byte {
	key := append(RatingByTaskPrefix, []byte(taskID)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(ratingID)...)
}

// GetRatingByTaskPrefix returns the prefix for iterating ratings by task
func GetRatingByTaskPrefix(taskID string) []byte {
	key := append(RatingByTaskPrefix, []byte(taskID)...)
	return append(key, 0x00)
}

// GetRatingByRatedUserKey returns the key for rating by rated user index
func GetRatingByRatedUserKey(ratedID, ratingID string) []byte {
	key := append(RatingByRatedUserPrefix, []byte(ratedID)...)
	key = append(key, 0x00) // separator
	return append(key, []byte(ratingID)...)
}

// GetRatingByRatedUserPrefix returns the prefix for iterating ratings by rated user
func GetRatingByRatedUserPrefix(ratedID string) []byte {
	key := append(RatingByRatedUserPrefix, []byte(ratedID)...)
	return append(key, 0x00)
}
