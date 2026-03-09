package types

import (
	"fmt"
	"time"
)

// IdeaStatus represents the status of an idea
type IdeaStatus string

const (
	IdeaStatusDraft      IdeaStatus = "draft"
	IdeaStatusActive     IdeaStatus = "active"
	IdeaStatusFunding    IdeaStatus = "funding"    // In crowdfunding
	IdeaStatusDeveloping IdeaStatus = "developing" // Being developed
	IdeaStatusCompleted  IdeaStatus = "completed"
	IdeaStatusArchived   IdeaStatus = "archived"
)

// ContributionCategory represents types of contributions
type ContributionCategory string

const (
	ContributionCode    ContributionCategory = "code"
	ContributionDesign  ContributionCategory = "design"
	ContributionDocs    ContributionCategory = "docs"
	ContributionResearch ContributionCategory = "research"
	ContributionMarketing ContributionCategory = "marketing"
	ContributionTesting  ContributionCategory = "testing"
)

// GetAllContributionCategories returns all contribution categories
func GetAllContributionCategories() []ContributionCategory {
	return []ContributionCategory{
		ContributionCode,
		ContributionDesign,
		ContributionDocs,
		ContributionResearch,
		ContributionMarketing,
		ContributionTesting,
	}
}

// CategoryWeights defines default weights for contribution categories
var CategoryWeights = map[ContributionCategory]float64{
	ContributionCode:      1.5,
	ContributionDesign:    1.2,
	ContributionDocs:      1.0,
	ContributionResearch:  1.3,
	ContributionMarketing: 1.0,
	ContributionTesting:   1.1,
}

// Idea represents a creative idea/project
type Idea struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	CreatorID    string     `json:"creator_id"`
	Status       IdeaStatus `json:"status"`
	CurrentVersion int      `json:"current_version"`

	// Tags and categories
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`

	// Statistics
	ViewCount        int     `json:"view_count"`
	ContributionCount int    `json:"contribution_count"`
	TotalWeight      float64 `json:"total_weight"`

	// Crowdfunding
	CampaignID string `json:"campaign_id"` // Associated crowdfunding campaign

	// Timestamps
	CreatedAt   int64 `json:"created_at"`
	UpdatedAt   int64 `json:"updated_at"`
	PublishedAt int64 `json:"published_at"`
}

// NewIdea creates a new idea
func NewIdea(id, title, description, creatorID string) *Idea {
	now := time.Now().Unix()
	return &Idea{
		ID:              id,
		Title:           title,
		Description:     description,
		CreatorID:       creatorID,
		Status:          IdeaStatusDraft,
		CurrentVersion:  1,
		Tags:            []string{},
		Categories:      []string{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Publish publishes the idea
func (i *Idea) Publish() {
	i.Status = IdeaStatusActive
	i.PublishedAt = time.Now().Unix()
	i.UpdatedAt = i.PublishedAt
}

// Update updates the idea
func (i *Idea) Update(title, description string) {
	if title != "" {
		i.Title = title
	}
	if description != "" {
		i.Description = description
	}
	i.UpdatedAt = time.Now().Unix()
	i.CurrentVersion++
}

// Archive archives the idea
func (i *Idea) Archive() {
	i.Status = IdeaStatusArchived
	i.UpdatedAt = time.Now().Unix()
}

// StartFunding starts crowdfunding
func (i *Idea) StartFunding(campaignID string) {
	i.Status = IdeaStatusFunding
	i.CampaignID = campaignID
	i.UpdatedAt = time.Now().Unix()
}

// StartDevelopment marks idea as being developed
func (i *Idea) StartDevelopment() {
	i.Status = IdeaStatusDeveloping
	i.UpdatedAt = time.Now().Unix()
}

// Complete marks idea as completed
func (i *Idea) Complete() {
	i.Status = IdeaStatusCompleted
	i.UpdatedAt = time.Now().Unix()
}

// AddTag adds a tag
func (i *Idea) AddTag(tag string) {
	for _, t := range i.Tags {
		if t == tag {
			return
		}
	}
	i.Tags = append(i.Tags, tag)
	i.UpdatedAt = time.Now().Unix()
}

// AddCategory adds a category
func (i *Idea) AddCategory(category string) {
	for _, c := range i.Categories {
		if c == category {
			return
		}
	}
	i.Categories = append(i.Categories, category)
	i.UpdatedAt = time.Now().Unix()
}

// Validate validates the idea
func (i *Idea) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("idea ID cannot be empty")
	}
	if i.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if i.CreatorID == "" {
		return fmt.Errorf("creator ID cannot be empty")
	}
	return nil
}

// IdeaVersion represents a version of an idea
type IdeaVersion struct {
	ID          string `json:"id"`
	IdeaID      string `json:"idea_id"`
	Version     int    `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Changes     string `json:"changes"` // Description of changes
	CreatedBy   string `json:"created_by"`
	CreatedAt   int64  `json:"created_at"`
}

// NewIdeaVersion creates a new idea version
func NewIdeaVersion(id, ideaID string, version int, title, description, changes, createdBy string) *IdeaVersion {
	return &IdeaVersion{
		ID:          id,
		IdeaID:      ideaID,
		Version:     version,
		Title:       title,
		Description: description,
		Changes:     changes,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now().Unix(),
	}
}

// Contribution represents a contribution to an idea
type Contribution struct {
	ID          string               `json:"id"`
	IdeaID      string               `json:"idea_id"`
	ContributorID string             `json:"contributor_id"`
	Category    ContributionCategory `json:"category"`
	Description string               `json:"description"`
	Weight      float64              `json:"weight"`       // Calculated weight
	RawScore    float64              `json:"raw_score"`    // Base contribution score
	Evidence    string               `json:"evidence"`     // Link to work/evidence
	Status      ContributionStatus   `json:"status"`
	ReviewedBy  string               `json:"reviewed_by"`
	ReviewedAt  int64                `json:"reviewed_at"`
	CreatedAt   int64                `json:"created_at"`
}

// ContributionStatus represents contribution status
type ContributionStatus string

const (
	ContributionStatusPending   ContributionStatus = "pending"
	ContributionStatusApproved  ContributionStatus = "approved"
	ContributionStatusRejected  ContributionStatus = "rejected"
)

// NewContribution creates a new contribution
func NewContribution(id, ideaID, contributorID string, category ContributionCategory, description string, rawScore float64) *Contribution {
	return &Contribution{
		ID:            id,
		IdeaID:        ideaID,
		ContributorID: contributorID,
		Category:      category,
		Description:   description,
		RawScore:      rawScore,
		Status:        ContributionStatusPending,
		CreatedAt:     time.Now().Unix(),
	}
}

// CalculateWeight calculates the weighted score
func (c *Contribution) CalculateWeight() {
	weight := CategoryWeights[c.Category]
	c.Weight = c.RawScore * weight
}

// Approve approves the contribution
func (c *Contribution) Approve(reviewerID string) {
	c.Status = ContributionStatusApproved
	c.ReviewedBy = reviewerID
	c.ReviewedAt = time.Now().Unix()
}

// Reject rejects the contribution
func (c *Contribution) Reject(reviewerID string) {
	c.Status = ContributionStatusRejected
	c.ReviewedBy = reviewerID
	c.ReviewedAt = time.Now().Unix()
}

// Validate validates the contribution
func (c *Contribution) Validate() error {
	if c.IdeaID == "" {
		return fmt.Errorf("idea ID cannot be empty")
	}
	if c.ContributorID == "" {
		return fmt.Errorf("contributor ID cannot be empty")
	}
	if c.RawScore <= 0 {
		return fmt.Errorf("raw score must be positive")
	}
	return nil
}

// ContributorStats represents statistics for a contributor
type ContributorStats struct {
	ContributorID    string                 `json:"contributor_id"`
	IdeaID           string                 `json:"idea_id"`
	TotalWeight      float64                `json:"total_weight"`
	ContributionCount int                   `json:"contribution_count"`
	ByCategory       map[ContributionCategory]float64 `json:"by_category"`
	ApprovedCount    int                    `json:"approved_count"`
	PendingCount     int                    `json:"pending_count"`
	RejectedCount    int                    `json:"rejected_count"`
}

// NewContributorStats creates new contributor stats
func NewContributorStats(contributorID, ideaID string) *ContributorStats {
	return &ContributorStats{
		ContributorID: contributorID,
		IdeaID:        ideaID,
		ByCategory:    make(map[ContributionCategory]float64),
	}
}

// AddContribution adds a contribution to stats
func (cs *ContributorStats) AddContribution(contribution *Contribution) {
	if contribution.Status != ContributionStatusApproved {
		return
	}

	cs.TotalWeight += contribution.Weight
	cs.ContributionCount++
	cs.ByCategory[contribution.Category] += contribution.Weight
	cs.ApprovedCount++
}

// ContributionSummary provides a summary of contributions
type ContributionSummary struct {
	IdeaID           string                 `json:"idea_id"`
	TotalContributors int                   `json:"total_contributors"`
	TotalWeight      float64                `json:"total_weight"`
	ByCategory       map[ContributionCategory]float64 `json:"by_category"`
	TopContributors  []ContributorStats     `json:"top_contributors"`
}

// NewContributionSummary creates a contribution summary
func NewContributionSummary(ideaID string) *ContributionSummary {
	return &ContributionSummary{
		IdeaID:     ideaID,
		ByCategory: make(map[ContributionCategory]float64),
	}
}
