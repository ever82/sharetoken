package types

import (
	"fmt"
	"time"
)

// NewTask creates a new task
func NewTask(id, title, description, requesterId string, taskType TaskType, budget uint64) *Task {
	now := time.Now().Unix()
	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
		RequesterId: requesterId,
		TaskType:    taskType,
		Status:      TaskStatus_TASK_STATUS_DRAFT,
		Budget:      budget,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Publish marks the task as open/published
func (t *Task) Publish() {
	t.Status = TaskStatus_TASK_STATUS_OPEN
	t.UpdatedAt = time.Now().Unix()
}

// ID returns the task ID
func (t *Task) ID() string {
	return t.Id
}

// RequesterID returns the requester ID
func (t *Task) RequesterID() string {
	return t.RequesterId
}

// WorkerID returns the worker ID
func (t *Task) WorkerID() string {
	return t.WorkerId
}

// AddMilestone adds a milestone to the task
func (t *Task) AddMilestone(milestone Milestone) {
	t.Milestones = append(t.Milestones, milestone)
}

// Validate validates the task
func (t *Task) Validate() error {
	if t.Id == "" {
		return fmt.Errorf("task ID is required")
	}
	if t.Title == "" {
		return fmt.Errorf("task title is required")
	}
	if t.RequesterId == "" {
		return fmt.Errorf("requester ID is required")
	}
	if t.Budget == 0 {
		return fmt.Errorf("budget must be greater than 0")
	}
	if t.Deadline <= time.Now().Unix() {
		return fmt.Errorf("deadline must be in the future")
	}
	return nil
}

// ValidateMilestones validates task milestones
func (t *Task) ValidateMilestones() error {
	if len(t.Milestones) == 0 {
		return nil
	}

	totalAmount := uint64(0)
	for _, m := range t.Milestones {
		if m.Title == "" {
			return fmt.Errorf("milestone title is required")
		}
		totalAmount += m.Amount
	}

	if totalAmount != t.Budget {
		return fmt.Errorf("milestone amounts must sum to task budget")
	}

	return nil
}

// Start marks the task as in progress
func (t *Task) Start() {
	t.Status = TaskStatus_TASK_STATUS_IN_PROGRESS
	t.UpdatedAt = time.Now().Unix()
}

// SubmitMilestone submits a milestone for review
func (t *Task) SubmitMilestone(milestoneID, deliverables string) error {
	for i := range t.Milestones {
		if t.Milestones[i].Id == milestoneID {
			if t.Milestones[i].Status != MilestoneStatus_MILESTONE_STATUS_ACTIVE {
				return fmt.Errorf("milestone is not active")
			}
			t.Milestones[i].Status = MilestoneStatus_MILESTONE_STATUS_SUBMITTED
			t.Milestones[i].Deliverables = deliverables
			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// ApproveMilestone approves a milestone
func (t *Task) ApproveMilestone(milestoneID string) error {
	for i := range t.Milestones {
		if t.Milestones[i].Id == milestoneID {
			if t.Milestones[i].Status != MilestoneStatus_MILESTONE_STATUS_SUBMITTED {
				return fmt.Errorf("milestone is not submitted")
			}
			t.Milestones[i].Status = MilestoneStatus_MILESTONE_STATUS_APPROVED
			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// RejectMilestone rejects a milestone
func (t *Task) RejectMilestone(milestoneID, reason string) error {
	for i := range t.Milestones {
		if t.Milestones[i].Id == milestoneID {
			if t.Milestones[i].Status != MilestoneStatus_MILESTONE_STATUS_SUBMITTED {
				return fmt.Errorf("milestone is not submitted")
			}
			t.Milestones[i].Status = MilestoneStatus_MILESTONE_STATUS_REJECTED
			return nil
		}
	}
	return fmt.Errorf("milestone not found: %s", milestoneID)
}

// AllMilestonesCompleted checks if all milestones are completed
func (t *Task) AllMilestonesCompleted() bool {
	for _, m := range t.Milestones {
		if m.Status != MilestoneStatus_MILESTONE_STATUS_APPROVED {
			return false
		}
	}
	return true
}

// Complete marks the task as completed
func (t *Task) Complete() {
	t.Status = TaskStatus_TASK_STATUS_COMPLETED
	t.CompletedAt = time.Now().Unix()
	t.UpdatedAt = time.Now().Unix()
}

// TaskStatusOpen is a shorthand for open task status
var TaskStatusOpen = TaskStatus_TASK_STATUS_OPEN

// TaskStatusInProgress is a shorthand for in-progress task status
var TaskStatusInProgress = TaskStatus_TASK_STATUS_IN_PROGRESS

// TaskStatusCompleted is a shorthand for completed task status
var TaskStatusCompleted = TaskStatus_TASK_STATUS_COMPLETED

// TaskStatusDraft is a shorthand for draft task status
var TaskStatusDraft = TaskStatus_TASK_STATUS_DRAFT

// TaskStatusAssigned is a shorthand for assigned task status
var TaskStatusAssigned = TaskStatus_TASK_STATUS_ASSIGNED

// TaskStatusCancelled is a shorthand for cancelled task status
var TaskStatusCancelled = TaskStatus_TASK_STATUS_CANCELLED

// Cancel marks the task as cancelled
func (t *Task) Cancel() {
	t.Status = TaskStatus_TASK_STATUS_CANCELLED
	t.UpdatedAt = time.Now().Unix()
}

// MilestoneStatusPending is a shorthand for pending milestone status
var MilestoneStatusPending = MilestoneStatus_MILESTONE_STATUS_PENDING

// MilestoneStatusActive is a shorthand for active milestone status
var MilestoneStatusActive = MilestoneStatus_MILESTONE_STATUS_ACTIVE

// MilestoneStatusSubmitted is a shorthand for submitted milestone status
var MilestoneStatusSubmitted = MilestoneStatus_MILESTONE_STATUS_SUBMITTED

// MilestoneStatusApproved is a shorthand for approved milestone status
var MilestoneStatusApproved = MilestoneStatus_MILESTONE_STATUS_APPROVED

// MilestoneStatusPaid is a shorthand for paid milestone status
var MilestoneStatusPaid = MilestoneStatus_MILESTONE_STATUS_PAID

// NewApplication creates a new application
func NewApplication(id, taskId, workerId string, proposedPrice uint64) *Application {
	now := time.Now().Unix()
	return &Application{
		Id:            id,
		TaskId:        taskId,
		WorkerId:      workerId,
		Status:        ApplicationStatus_APPLICATION_STATUS_PENDING,
		ProposedPrice: proposedPrice,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewBid creates a new bid
func NewBid(id, taskId, workerId string, amount uint64) *Bid {
	now := time.Now().Unix()
	return &Bid{
		Id:        id,
		TaskId:    taskId,
		WorkerId:  workerId,
		Amount:    amount,
		Status:    BidStatus_BID_STATUS_PENDING,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// BidStatusPending is a shorthand for pending bid status
var BidStatusPending = BidStatus_BID_STATUS_PENDING

// BidStatusAccepted is a shorthand for accepted bid status
var BidStatusAccepted = BidStatus_BID_STATUS_ACCEPTED

// BidStatusRejected is a shorthand for rejected bid status
var BidStatusRejected = BidStatus_BID_STATUS_REJECTED

// BidStatusOutbid is a shorthand for outbid bid status
var BidStatusOutbid = BidStatus_BID_STATUS_OUTBID

// BidStatusWithdrawn is a shorthand for withdrawn bid status
var BidStatusWithdrawn = BidStatus_BID_STATUS_WITHDRAWN

// Validate validates a bid
func (b Bid) Validate() error {
	if b.Id == "" {
		return fmt.Errorf("bid ID is required")
	}
	if b.TaskId == "" {
		return fmt.Errorf("task ID is required")
	}
	if b.WorkerId == "" {
		return fmt.Errorf("worker ID is required")
	}
	if b.Amount == 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

// Accept marks the bid as accepted
func (b *Bid) Accept() {
	b.Status = BidStatus_BID_STATUS_ACCEPTED
}

// Reject marks the bid as rejected
func (b *Bid) Reject() {
	b.Status = BidStatus_BID_STATUS_REJECTED
}

// Withdraw marks the bid as withdrawn
func (b *Bid) Withdraw() {
	b.Status = BidStatus_BID_STATUS_WITHDRAWN
}

// IsLowerThan checks if this bid has a lower amount than another bid
func (b *Bid) IsLowerThan(other *Bid) bool {
	return b.Amount < other.Amount
}

// NewAuction creates a new auction
func NewAuction(taskId string, startingPrice, reservePrice uint64, durationSec int64) *Auction {
	now := time.Now().Unix()
	return &Auction{
		TaskId:        taskId,
		StartingPrice: startingPrice,
		ReservePrice:  reservePrice,
		IsActive:      true,
		Bids:          []Bid{},
		EndTime:       now + durationSec,
	}
}

// AddBid adds a bid to the auction
func (a *Auction) AddBid(bid Bid) error {
	if !a.IsActive {
		return fmt.Errorf("auction is not active")
	}
	if time.Now().Unix() > a.EndTime {
		return fmt.Errorf("auction has ended")
	}
	if bid.Amount > a.StartingPrice {
		return fmt.Errorf("bid exceeds starting price")
	}

	// Mark existing higher bids as outbid
	for i := range a.Bids {
		if a.Bids[i].Amount > bid.Amount && a.Bids[i].Status == BidStatus_BID_STATUS_PENDING {
			a.Bids[i].Status = BidStatus_BID_STATUS_OUTBID
		}
	}

	a.Bids = append(a.Bids, bid)
	return nil
}

// GetWinningBid returns the current winning bid (lowest amount)
func (a *Auction) GetWinningBid() *Bid {
	if len(a.Bids) == 0 {
		return nil
	}

	var winner *Bid
	for i := range a.Bids {
		if a.Bids[i].Status == BidStatus_BID_STATUS_PENDING {
			if winner == nil || a.Bids[i].Amount < winner.Amount {
				winner = &a.Bids[i]
			}
		}
	}
	return winner
}

// CloseAuction closes the auction and returns the winning bid
func (a *Auction) CloseAuction() (*Bid, error) {
	if !a.IsActive {
		return nil, fmt.Errorf("auction is not active")
	}

	winner := a.GetWinningBid()
	if winner == nil {
		return nil, fmt.Errorf("no valid bids")
	}

	if winner.Amount > a.ReservePrice {
		return nil, fmt.Errorf("winning bid does not meet reserve price")
	}

	winner.Accept()
	a.IsActive = false
	a.WinningBidId = winner.Id

	return winner, nil
}

// IsEnded checks if the auction has ended
func (a *Auction) IsEnded() bool {
	return !a.IsActive || time.Now().Unix() > a.EndTime
}

// GetTotalMilestoneAmount returns the total amount of all milestones
func (t *Task) GetTotalMilestoneAmount() uint64 {
	total := uint64(0)
	for _, m := range t.Milestones {
		total += m.Amount
	}
	return total
}

// GetCompletionPercentage returns the completion percentage based on milestones
func (t *Task) GetCompletionPercentage() float64 {
	if t.Status == TaskStatus_TASK_STATUS_COMPLETED {
		return 100.0
	}
	if len(t.Milestones) == 0 {
		return 0.0
	}

	completed := 0
	for _, m := range t.Milestones {
		if m.Status == MilestoneStatus_MILESTONE_STATUS_APPROVED ||
			m.Status == MilestoneStatus_MILESTONE_STATUS_PAID {
			completed++
		}
	}

	return float64(completed) / float64(len(t.Milestones)) * 100.0
}

// NewRating creates a new rating
func NewRating(id, taskId, raterId, ratedId string) *Rating {
	return &Rating{
		Id:        id,
		TaskId:    taskId,
		RaterId:   raterId,
		RatedId:   ratedId,
		Ratings:   make(map[string]int32),
		CreatedAt: time.Now().Unix(),
	}
}

// SetRating sets a rating value for a dimension
func (r *Rating) SetRating(dim string, val int32) error {
	if val < 1 || val > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	r.Ratings[dim] = val
	return nil
}

// Validate validates the rating
func (r *Rating) Validate() error {
	if r.Id == "" {
		return fmt.Errorf("rating ID is required")
	}
	if r.TaskId == "" {
		return fmt.Errorf("task ID is required")
	}
	if r.RaterId == "" {
		return fmt.Errorf("rater ID is required")
	}
	if r.RatedId == "" {
		return fmt.Errorf("rated ID is required")
	}
	if len(r.Ratings) == 0 {
		return fmt.Errorf("at least one rating is required")
	}
	return nil
}

// NewReputation creates a new reputation for a user
func NewReputation(userId string) *Reputation {
	now := time.Now().Unix()
	return &Reputation{
		UserId:             userId,
		RatingsByDimension: make(map[string]float64),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// AddRating adds a rating to the reputation
func (r *Reputation) AddRating(rating *Rating) {
	r.TotalRatings++

	// Update ratings by dimension
	for dim, val := range rating.Ratings {
		if existing, ok := r.RatingsByDimension[dim]; ok {
			// Weighted average
			r.RatingsByDimension[dim] = (existing*float64(r.TotalRatings-1) + float64(val)) / float64(r.TotalRatings)
		} else {
			r.RatingsByDimension[dim] = float64(val)
		}
	}

	// Recalculate average rating
	var total float64
	for _, val := range r.RatingsByDimension {
		total += val
	}
	r.AverageRating = total / float64(len(r.RatingsByDimension))

	r.UpdatedAt = time.Now().Unix()
}
