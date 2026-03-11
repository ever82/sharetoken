package types

import (
	"errors"
)

// Event types and attributes
const (
	EventTypeCreateTask        = "create_task"
	EventTypeUpdateTask        = "update_task"
	EventTypePublishTask       = "publish_task"
	EventTypeCancelTask        = "cancel_task"
	EventTypeSubmitApplication = "submit_application"
	EventTypeAcceptApplication = "accept_application"
	EventTypeRejectApplication = "reject_application"
	EventTypeSubmitBid         = "submit_bid"
	EventTypeCloseAuction      = "close_auction"
	EventTypeStartTask         = "start_task"
	EventTypeSubmitMilestone   = "submit_milestone"
	EventTypeApproveMilestone  = "approve_milestone"
	EventTypeRejectMilestone   = "reject_milestone"
	EventTypeSubmitRating      = "submit_rating"

	AttributeKeyTaskID        = "task_id"
	AttributeKeyCreator       = "creator"
	AttributeKeyTaskType      = "task_type"
	AttributeKeyApplicationID = "application_id"
	AttributeKeyBidID         = "bid_id"
	AttributeKeyRatingID      = "rating_id"
	AttributeKeyWorkerID      = "worker_id"
	AttributeKeyWinnerID      = "winner_id"
	AttributeKeyRatedID       = "rated_id"
)

// Common errors for taskmarket module
var (
	ErrTaskNotFound         = errors.New("task not found")
	ErrApplicationNotFound  = errors.New("application not found")
	ErrAuctionNotFound      = errors.New("auction not found")
	ErrBidNotFound          = errors.New("bid not found")
	ErrRatingNotFound       = errors.New("rating not found")
	ErrReputationNotFound   = errors.New("reputation not found")
	ErrInvalidTaskStatus    = errors.New("invalid task status")
	ErrInvalidBidAmount     = errors.New("invalid bid amount")
	ErrTaskNotOpen          = errors.New("task is not open")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidMilestone     = errors.New("invalid milestone")
	ErrMilestoneNotFound    = errors.New("milestone not found")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrDuplicateApplication = errors.New("duplicate application")
	ErrDuplicateBid         = errors.New("duplicate bid")
	ErrAuctionClosed        = errors.New("auction is closed")
	ErrAuctionNotEnded      = errors.New("auction has not ended")
)
