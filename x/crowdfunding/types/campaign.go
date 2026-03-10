package types

import (
	"fmt"
	"time"
)

// CampaignType represents the type of crowdfunding campaign
type CampaignType string

const (
	CampaignTypeInvestment CampaignType = "investment" // Equity investment
	CampaignTypeLending    CampaignType = "lending"    // Loan with interest
	CampaignTypeDonation   CampaignType = "donation"   // No return
)

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusFunded    CampaignStatus = "funded"  // Reached goal
	CampaignStatusExpired   CampaignStatus = "expired" // Time ran out
	CampaignStatusCancelled CampaignStatus = "cancelled"
	CampaignStatusClosed    CampaignStatus = "closed" // Funds distributed
)

// Campaign represents a crowdfunding campaign
type Campaign struct {
	ID          string         `json:"id"`
	IdeaID      string         `json:"idea_id"`
	CreatorID   string         `json:"creator_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        CampaignType   `json:"type"`
	Status      CampaignStatus `json:"status"`

	// Funding goal
	GoalAmount   uint64 `json:"goal_amount"`   // Target amount
	RaisedAmount uint64 `json:"raised_amount"` // Currently raised
	Currency     string `json:"currency"`      // "STT", "USDC"

	// Time limits
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`

	// Campaign specific
	MinContribution uint64 `json:"min_contribution"`
	MaxContribution uint64 `json:"max_contribution"`

	// For investment type
	EquityOffered float64 `json:"equity_offered"` // Percentage (0-100)
	Valuation     uint64  `json:"valuation"`      // Company valuation

	// For lending type
	InterestRate      float64 `json:"interest_rate"`      // Annual interest rate
	LoanTerm          int64   `json:"loan_term"`          // Term in days
	RepaymentSchedule string  `json:"repayment_schedule"` // "monthly", "quarterly", "lump"

	// Statistics
	BackerCount int `json:"backer_count"`
	UpdateCount int `json:"update_count"`

	// Timestamps
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// NewCampaign creates a new campaign
func NewCampaign(id, ideaID, creatorID, title string, campaignType CampaignType, goal uint64) *Campaign {
	now := time.Now().Unix()
	return &Campaign{
		ID:              id,
		IdeaID:          ideaID,
		CreatorID:       creatorID,
		Title:           title,
		Type:            campaignType,
		Status:          CampaignStatusDraft,
		GoalAmount:      goal,
		Currency:        "STT",
		MinContribution: 1,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Launch launches the campaign
func (c *Campaign) Launch(duration int64) {
	c.Status = CampaignStatusActive
	c.StartTime = time.Now().Unix()
	c.EndTime = c.StartTime + duration
	c.UpdatedAt = c.StartTime
}

// Contribute adds a contribution to the campaign
func (c *Campaign) Contribute(amount uint64) error {
	if c.Status != CampaignStatusActive {
		return fmt.Errorf("campaign is not active")
	}

	if time.Now().Unix() > c.EndTime {
		return fmt.Errorf("campaign has ended")
	}

	if amount < c.MinContribution {
		return fmt.Errorf("contribution below minimum")
	}

	if c.MaxContribution > 0 && amount > c.MaxContribution {
		return fmt.Errorf("contribution exceeds maximum")
	}

	c.RaisedAmount += amount
	c.BackerCount++
	c.UpdatedAt = time.Now().Unix()

	// Check if goal reached
	if c.RaisedAmount >= c.GoalAmount {
		c.Status = CampaignStatusFunded
	}

	return nil
}

// CheckExpired checks if campaign has expired
func (c *Campaign) CheckExpired() {
	if c.Status == CampaignStatusActive && time.Now().Unix() > c.EndTime {
		if c.RaisedAmount >= c.GoalAmount {
			c.Status = CampaignStatusFunded
		} else {
			c.Status = CampaignStatusExpired
		}
		c.UpdatedAt = time.Now().Unix()
	}
}

// Cancel cancels the campaign
func (c *Campaign) Cancel() {
	c.Status = CampaignStatusCancelled
	c.UpdatedAt = time.Now().Unix()
}

// Close closes the campaign after distribution
func (c *Campaign) Close() {
	c.Status = CampaignStatusClosed
	c.UpdatedAt = time.Now().Unix()
}

// GetProgress returns funding progress percentage
func (c *Campaign) GetProgress() float64 {
	if c.GoalAmount == 0 {
		return 0
	}
	return float64(c.RaisedAmount) / float64(c.GoalAmount) * 100
}

// GetTimeRemaining returns time remaining in seconds
func (c *Campaign) GetTimeRemaining() int64 {
	remaining := c.EndTime - time.Now().Unix()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsActive checks if campaign is active
func (c *Campaign) IsActive() bool {
	return c.Status == CampaignStatusActive
}

// IsFunded checks if campaign is funded
func (c *Campaign) IsFunded() bool {
	return c.Status == CampaignStatusFunded
}

// Validate validates the campaign
func (c *Campaign) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if c.IdeaID == "" {
		return fmt.Errorf("idea ID cannot be empty")
	}
	if c.CreatorID == "" {
		return fmt.Errorf("creator ID cannot be empty")
	}
	if c.GoalAmount == 0 {
		return fmt.Errorf("goal amount must be greater than 0")
	}

	// Validate campaign type specific fields
	switch c.Type {
	case CampaignTypeInvestment:
		if c.EquityOffered <= 0 || c.EquityOffered > 100 {
			return fmt.Errorf("equity offered must be between 0 and 100")
		}
		if c.Valuation == 0 {
			return fmt.Errorf("valuation must be specified for investment")
		}
	case CampaignTypeLending:
		if c.InterestRate < 0 {
			return fmt.Errorf("interest rate cannot be negative")
		}
		if c.LoanTerm <= 0 {
			return fmt.Errorf("loan term must be specified for lending")
		}
	}

	return nil
}

// GetExpectedReturn calculates expected return for backers
func (c *Campaign) GetExpectedReturn(amount uint64) (uint64, error) {
	switch c.Type {
	case CampaignTypeInvestment:
		// Return is based on equity
		// e.g., if equity offered is 10% for 1000 STT goal,
		// then 100 STT investment gets 1% equity
		equity := float64(amount) / float64(c.GoalAmount) * c.EquityOffered
		return uint64(equity * 100), nil // Return basis points

	case CampaignTypeLending:
		// Simple interest calculation
		interest := float64(amount) * c.InterestRate * float64(c.LoanTerm) / 365 / 100
		return amount + uint64(interest), nil

	case CampaignTypeDonation:
		return 0, nil // No return

	default:
		return 0, fmt.Errorf("unknown campaign type")
	}
}

// Backer represents a backer/contribution to a campaign
type Backer struct {
	ID           string `json:"id"`
	CampaignID   string `json:"campaign_id"`
	BackerID     string `json:"backer_id"`
	Amount       uint64 `json:"amount"`
	Currency     string `json:"currency"`
	Message      string `json:"message"`
	Refunded     bool   `json:"refunded"`
	RefundAmount uint64 `json:"refund_amount"`
	CreatedAt    int64  `json:"created_at"`
}

// NewBacker creates a new backer
func NewBacker(id, campaignID, backerID string, amount uint64) *Backer {
	return &Backer{
		ID:         id,
		CampaignID: campaignID,
		BackerID:   backerID,
		Amount:     amount,
		Currency:   "STT",
		CreatedAt:  time.Now().Unix(),
	}
}

// Refund refunds the backer
func (b *Backer) Refund(amount uint64) {
	b.Refunded = true
	b.RefundAmount = amount
}

// CampaignUpdate represents an update from the campaign creator
type CampaignUpdate struct {
	ID         string `json:"id"`
	CampaignID string `json:"campaign_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  int64  `json:"created_at"`
}

// NewCampaignUpdate creates a new campaign update
func NewCampaignUpdate(id, campaignID, title, content, createdBy string) *CampaignUpdate {
	return &CampaignUpdate{
		ID:         id,
		CampaignID: campaignID,
		Title:      title,
		Content:    content,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now().Unix(),
	}
}

// RevenueDistribution represents revenue distribution to contributors
type RevenueDistribution struct {
	ID            string              `json:"id"`
	IdeaID        string              `json:"idea_id"`
	CampaignID    string              `json:"campaign_id"`
	TotalAmount   uint64              `json:"total_amount"`
	DistributedAt int64               `json:"distributed_at"`
	Distributions []ContributorPayout `json:"distributions"`
	Status        DistributionStatus  `json:"status"`
}

// DistributionStatus represents distribution status
type DistributionStatus string

const (
	DistributionStatusPending   DistributionStatus = "pending"
	DistributionStatusProcessed DistributionStatus = "processed"
	DistributionStatusFailed    DistributionStatus = "failed"
)

// ContributorPayout represents a payout to a contributor
type ContributorPayout struct {
	ContributorID string  `json:"contributor_id"`
	Weight        float64 `json:"weight"`
	Share         float64 `json:"share"` // Percentage
	Amount        uint64  `json:"amount"`
	Status        string  `json:"status"` // "pending", "sent", "failed"
}

// NewRevenueDistribution creates a new revenue distribution
func NewRevenueDistribution(id, ideaID string, totalAmount uint64) *RevenueDistribution {
	return &RevenueDistribution{
		ID:            id,
		IdeaID:        ideaID,
		TotalAmount:   totalAmount,
		Distributions: []ContributorPayout{},
		Status:        DistributionStatusPending,
	}
}

// CalculateDistributions calculates payouts based on contributor weights
func (rd *RevenueDistribution) CalculateDistributions(contributors map[string]float64) {
	var totalWeight float64
	for _, weight := range contributors {
		totalWeight += weight
	}

	if totalWeight == 0 {
		return
	}

	for contributorID, weight := range contributors {
		share := weight / totalWeight
		amount := uint64(float64(rd.TotalAmount) * share)

		payout := ContributorPayout{
			ContributorID: contributorID,
			Weight:        weight,
			Share:         share * 100, // Convert to percentage
			Amount:        amount,
			Status:        "pending",
		}

		rd.Distributions = append(rd.Distributions, payout)
	}

	rd.Status = DistributionStatusProcessed
	rd.DistributedAt = time.Now().Unix()
}

// CampaignStats provides statistics for a campaign
type CampaignStats struct {
	CampaignID           string  `json:"campaign_id"`
	AverageContribution  float64 `json:"average_contribution"`
	LargestContribution  uint64  `json:"largest_contribution"`
	SmallestContribution uint64  `json:"smallest_contribution"`
	BackerRetention      float64 `json:"backer_retention"` // Repeat backers
}
