package types

import (
	"fmt"
	"time"
)

// ApplicationStatus represents application status
type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusAccepted  ApplicationStatus = "accepted"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn ApplicationStatus = "withdrawn"
)

// Application represents a worker's application for an open task
type Application struct {
	ID          string            `json:"id"`
	TaskID      string            `json:"task_id"`
	WorkerID    string            `json:"worker_id"`
	Status      ApplicationStatus `json:"status"`
	ProposedPrice uint64          `json:"proposed_price"`
	CoverLetter string            `json:"cover_letter"`
	RelevantExperience []string   `json:"relevant_experience"`
	PortfolioLinks []string       `json:"portfolio_links"`
	EstimatedDuration int64       `json:"estimated_duration"` // days
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}

// NewApplication creates a new application
func NewApplication(id, taskID, workerID string, price uint64) *Application {
	now := time.Now().Unix()
	return &Application{
		ID:               id,
		TaskID:           taskID,
		WorkerID:         workerID,
		Status:           ApplicationStatusPending,
		ProposedPrice:    price,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Accept accepts the application
func (a *Application) Accept() {
	a.Status = ApplicationStatusAccepted
	a.UpdatedAt = time.Now().Unix()
}

// Reject rejects the application
func (a *Application) Reject() {
	a.Status = ApplicationStatusRejected
	a.UpdatedAt = time.Now().Unix()
}

// Withdraw withdraws the application
func (a *Application) Withdraw() {
	a.Status = ApplicationStatusWithdrawn
	a.UpdatedAt = time.Now().Unix()
}

// Validate validates the application
func (a *Application) Validate() error {
	if a.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if a.WorkerID == "" {
		return fmt.Errorf("worker ID cannot be empty")
	}
	if a.ProposedPrice == 0 {
		return fmt.Errorf("proposed price must be greater than 0")
	}
	return nil
}

// BidStatus represents bid status
type BidStatus string

const (
	BidStatusPending   BidStatus = "pending"
	BidStatusAccepted  BidStatus = "accepted"
	BidStatusRejected  BidStatus = "rejected"
	BidStatusWithdrawn BidStatus = "withdrawn"
	BidStatusOutbid    BidStatus = "outbid"
)

// Bid represents a bid for an auction task
type Bid struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	WorkerID  string    `json:"worker_id"`
	Amount    uint64    `json:"amount"`     // Bid amount (lower is better)
	Status    BidStatus `json:"status"`
	Message   string    `json:"message"`
	Portfolio string    `json:"portfolio"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

// NewBid creates a new bid
func NewBid(id, taskID, workerID string, amount uint64) *Bid {
	now := time.Now().Unix()
	return &Bid{
		ID:        id,
		TaskID:    taskID,
		WorkerID:  workerID,
		Amount:    amount,
		Status:    BidStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Accept accepts the bid
func (b *Bid) Accept() {
	b.Status = BidStatusAccepted
	b.UpdatedAt = time.Now().Unix()
}

// Reject rejects the bid
func (b *Bid) Reject() {
	b.Status = BidStatusRejected
	b.UpdatedAt = time.Now().Unix()
}

// Withdraw withdraws the bid
func (b *Bid) Withdraw() {
	b.Status = BidStatusWithdrawn
	b.UpdatedAt = time.Now().Unix()
}

// MarkOutbid marks the bid as outbid
func (b *Bid) MarkOutbid() {
	b.Status = BidStatusOutbid
	b.UpdatedAt = time.Now().Unix()
}

// Validate validates the bid
func (b *Bid) Validate() error {
	if b.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if b.WorkerID == "" {
		return fmt.Errorf("worker ID cannot be empty")
	}
	if b.Amount == 0 {
		return fmt.Errorf("bid amount must be greater than 0")
	}
	return nil
}

// IsLowerThan checks if this bid is lower than another
func (b *Bid) IsLowerThan(other *Bid) bool {
	return b.Amount < other.Amount
}

// Auction represents an auction with bids
type Auction struct {
	TaskID        string    `json:"task_id"`
	StartingPrice uint64    `json:"starting_price"`
	ReservePrice  uint64    `json:"reserve_price"`   // Minimum acceptable price
	EndTime       int64     `json:"end_time"`        // Auction end time
	Bids          []Bid     `json:"bids"`
	WinningBidID  string    `json:"winning_bid_id"`
	IsActive      bool      `json:"is_active"`
}

// NewAuction creates a new auction
func NewAuction(taskID string, startingPrice, reservePrice uint64, duration int64) *Auction {
	return &Auction{
		TaskID:        taskID,
		StartingPrice: startingPrice,
		ReservePrice:  reservePrice,
		EndTime:       time.Now().Unix() + duration,
		Bids:          []Bid{},
		IsActive:      true,
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

	a.Bids = append(a.Bids, bid)
	a.updateWinningBid()
	return nil
}

// updateWinningBid updates the winning bid
func (a *Auction) updateWinningBid() {
	if len(a.Bids) == 0 {
		return
	}

	// Find lowest bid
	var lowest *Bid
	for i := range a.Bids {
		if a.Bids[i].Status == BidStatusPending {
			if lowest == nil || a.Bids[i].Amount < lowest.Amount {
				lowest = &a.Bids[i]
			}
		}
	}

	if lowest != nil {
		// Mark previous winner as outbid
		if a.WinningBidID != "" {
			for i := range a.Bids {
				if a.Bids[i].ID == a.WinningBidID {
					a.Bids[i].MarkOutbid()
					break
				}
			}
		}
		a.WinningBidID = lowest.ID
	}
}

// GetWinningBid returns the current winning bid
func (a *Auction) GetWinningBid() *Bid {
	for i := range a.Bids {
		if a.Bids[i].ID == a.WinningBidID {
			return &a.Bids[i]
		}
	}
	return nil
}

// CloseAuction closes the auction and selects winner
func (a *Auction) CloseAuction() (*Bid, error) {
	a.IsActive = false

	winner := a.GetWinningBid()
	if winner == nil {
		return nil, fmt.Errorf("no valid bids")
	}

	if winner.Amount > a.ReservePrice {
		return nil, fmt.Errorf("winning bid does not meet reserve price")
	}

	winner.Accept()
	return winner, nil
}

// GetValidBids returns all valid (pending) bids
func (a *Auction) GetValidBids() []Bid {
	var valid []Bid
	for _, bid := range a.Bids {
		if bid.Status == BidStatusPending {
			valid = append(valid, bid)
		}
	}
	return valid
}

// GetBidCount returns the number of valid bids
func (a *Auction) GetBidCount() int {
	return len(a.GetValidBids())
}

// GetLowestBidAmount returns the lowest bid amount
func (a *Auction) GetLowestBidAmount() uint64 {
	winner := a.GetWinningBid()
	if winner != nil {
		return winner.Amount
	}
	return a.StartingPrice
}

// IsEnded checks if auction has ended
func (a *Auction) IsEnded() bool {
	return time.Now().Unix() > a.EndTime || !a.IsActive
}

// TimeRemaining returns time remaining in seconds
func (a *Auction) TimeRemaining() int64 {
	remaining := a.EndTime - time.Now().Unix()
	if remaining < 0 {
		return 0
	}
	return remaining
}
