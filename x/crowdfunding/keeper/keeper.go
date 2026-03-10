package keeper

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"sharetoken/x/crowdfunding/types"
)

// Keeper manages crowdfunding state
type Keeper struct {
	ideas            map[string]*types.Idea
	versions         map[string][]*types.IdeaVersion
	contributions    map[string]*types.Contribution
	campaigns        map[string]*types.Campaign
	backers          map[string]*types.Backer
	updates          map[string][]*types.CampaignUpdate
	contributorStats map[string]*types.ContributorStats
	mutex            sync.RWMutex
}

// NewKeeper creates a new crowdfunding keeper
func NewKeeper() *Keeper {
	return &Keeper{
		ideas:            make(map[string]*types.Idea),
		versions:         make(map[string][]*types.IdeaVersion),
		contributions:    make(map[string]*types.Contribution),
		campaigns:        make(map[string]*types.Campaign),
		backers:          make(map[string]*types.Backer),
		updates:          make(map[string][]*types.CampaignUpdate),
		contributorStats: make(map[string]*types.ContributorStats),
	}
}

// CreateIdea creates a new idea
func (k *Keeper) CreateIdea(idea *types.Idea) error {
	if err := idea.Validate(); err != nil {
		return fmt.Errorf("invalid idea: %w", err)
	}

	k.mutex.Lock()
	k.ideas[idea.ID] = idea
	k.mutex.Unlock()

	return nil
}

// GetIdea retrieves an idea
func (k *Keeper) GetIdea(id string) *types.Idea {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	return k.ideas[id]
}

// UpdateIdea updates an idea and creates a version
func (k *Keeper) UpdateIdea(ideaID, title, description, changes, updatedBy string) (*types.IdeaVersion, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	idea, exists := k.ideas[ideaID]
	if !exists {
		return nil, fmt.Errorf("idea not found: %s", ideaID)
	}

	// Create version record
	version := types.NewIdeaVersion(
		fmt.Sprintf("%s-v%d", ideaID, idea.CurrentVersion),
		ideaID,
		idea.CurrentVersion,
		idea.Title,
		idea.Description,
		changes,
		updatedBy,
	)

	// Store version
	k.versions[ideaID] = append(k.versions[ideaID], version)

	// Update idea
	idea.Update(title, description)

	return version, nil
}

// GetIdeaVersions retrieves all versions of an idea
func (k *Keeper) GetIdeaVersions(ideaID string) []*types.IdeaVersion {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	return k.versions[ideaID]
}

// PublishIdea publishes an idea
func (k *Keeper) PublishIdea(ideaID string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	idea, exists := k.ideas[ideaID]
	if !exists {
		return fmt.Errorf("idea not found: %s", ideaID)
	}

	idea.Publish()
	return nil
}

// SubmitContribution submits a contribution
func (k *Keeper) SubmitContribution(contribution *types.Contribution) error {
	if err := contribution.Validate(); err != nil {
		return fmt.Errorf("invalid contribution: %w", err)
	}

	k.mutex.Lock()
	defer k.mutex.Unlock()

	idea, exists := k.ideas[contribution.IdeaID]
	if !exists {
		return fmt.Errorf("idea not found: %s", contribution.IdeaID)
	}

	// Calculate weight
	contribution.CalculateWeight()

	// Store contribution
	k.contributions[contribution.ID] = contribution

	// Update idea stats
	idea.ContributionCount++
	idea.TotalWeight += contribution.Weight

	// Update contributor stats
	statsKey := fmt.Sprintf("%s:%s", contribution.IdeaID, contribution.ContributorID)
	stats, exists := k.contributorStats[statsKey]
	if !exists {
		stats = types.NewContributorStats(contribution.ContributorID, contribution.IdeaID)
		k.contributorStats[statsKey] = stats
	}
	stats.PendingCount++

	return nil
}

// ApproveContribution approves a contribution
func (k *Keeper) ApproveContribution(contributionID, reviewerID string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	contribution, exists := k.contributions[contributionID]
	if !exists {
		return fmt.Errorf("contribution not found: %s", contributionID)
	}

	contribution.Approve(reviewerID)

	// Update contributor stats
	statsKey := fmt.Sprintf("%s:%s", contribution.IdeaID, contribution.ContributorID)
	stats, exists := k.contributorStats[statsKey]
	if exists {
		stats.AddContribution(contribution)
		stats.PendingCount--
	}

	return nil
}

// RejectContribution rejects a contribution
func (k *Keeper) RejectContribution(contributionID, reviewerID string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	contribution, exists := k.contributions[contributionID]
	if !exists {
		return fmt.Errorf("contribution not found: %s", contributionID)
	}

	contribution.Reject(reviewerID)

	// Update contributor stats
	statsKey := fmt.Sprintf("%s:%s", contribution.IdeaID, contribution.ContributorID)
	stats, exists := k.contributorStats[statsKey]
	if exists {
		stats.RejectedCount++
		stats.PendingCount--
	}

	return nil
}

// GetContributionsByIdea returns contributions for an idea
func (k *Keeper) GetContributionsByIdea(ideaID string) []*types.Contribution {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	var contributions []*types.Contribution
	for _, c := range k.contributions {
		if c.IdeaID == ideaID {
			contributions = append(contributions, c)
		}
	}

	return contributions
}

// GetContributionsByContributor returns contributions by a contributor
func (k *Keeper) GetContributionsByContributor(contributorID string) []*types.Contribution {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	var contributions []*types.Contribution
	for _, c := range k.contributions {
		if c.ContributorID == contributorID {
			contributions = append(contributions, c)
		}
	}

	return contributions
}

// GetContributorStats returns stats for a contributor on an idea
func (k *Keeper) GetContributorStats(ideaID, contributorID string) *types.ContributorStats {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	statsKey := fmt.Sprintf("%s:%s", ideaID, contributorID)
	stats, exists := k.contributorStats[statsKey]
	if !exists {
		return types.NewContributorStats(contributorID, ideaID)
	}
	return stats
}

// CreateCampaign creates a new crowdfunding campaign
func (k *Keeper) CreateCampaign(campaign *types.Campaign) error {
	if err := campaign.Validate(); err != nil {
		return fmt.Errorf("invalid campaign: %w", err)
	}

	k.mutex.Lock()
	k.campaigns[campaign.ID] = campaign
	k.mutex.Unlock()

	return nil
}

// GetCampaign retrieves a campaign
func (k *Keeper) GetCampaign(id string) *types.Campaign {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	return k.campaigns[id]
}

// LaunchCampaign launches a campaign
func (k *Keeper) LaunchCampaign(campaignID string, duration int64) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	campaign, exists := k.campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found: %s", campaignID)
	}

	campaign.Launch(duration)

	// Update idea status
	if idea, exists := k.ideas[campaign.IdeaID]; exists {
		idea.StartFunding(campaignID)
	}

	return nil
}

// ContributeToCampaign adds a contribution to a campaign
func (k *Keeper) ContributeToCampaign(backer *types.Backer) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	campaign, exists := k.campaigns[backer.CampaignID]
	if !exists {
		return fmt.Errorf("campaign not found: %s", backer.CampaignID)
	}

	if err := campaign.Contribute(backer.Amount); err != nil {
		return err
	}

	k.backers[backer.ID] = backer
	return nil
}

// GetBackersByCampaign returns backers for a campaign
func (k *Keeper) GetBackersByCampaign(campaignID string) []*types.Backer {
	k.mutex.RLock()
	defer k.mutex.Unlock()

	var backers []*types.Backer
	for _, b := range k.backers {
		if b.CampaignID == campaignID {
			backers = append(backers, b)
		}
	}

	return backers
}

// AddCampaignUpdate adds an update to a campaign
func (k *Keeper) AddCampaignUpdate(update *types.CampaignUpdate) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	campaign, exists := k.campaigns[update.CampaignID]
	if !exists {
		return fmt.Errorf("campaign not found: %s", update.CampaignID)
	}

	k.updates[update.CampaignID] = append(k.updates[update.CampaignID], update)
	campaign.UpdateCount++
	campaign.UpdatedAt = time.Now().Unix()

	return nil
}

// GetCampaignUpdates returns updates for a campaign
func (k *Keeper) GetCampaignUpdates(campaignID string) []*types.CampaignUpdate {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	return k.updates[campaignID]
}

// CloseCampaign closes a funded campaign
func (k *Keeper) CloseCampaign(campaignID string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	campaign, exists := k.campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found: %s", campaignID)
	}

	if !campaign.IsFunded() {
		return fmt.Errorf("campaign is not funded")
	}

	campaign.Close()

	// Update idea status
	if idea, exists := k.ideas[campaign.IdeaID]; exists {
		idea.StartDevelopment()
	}

	return nil
}

// GetActiveCampaigns returns active campaigns
func (k *Keeper) GetActiveCampaigns() []*types.Campaign {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	var active []*types.Campaign
	for _, campaign := range k.campaigns {
		if campaign.IsActive() {
			campaign.CheckExpired()
			if campaign.IsActive() {
				active = append(active, campaign)
			}
		}
	}

	// Sort by end time (ending soonest first)
	sort.Slice(active, func(i, j int) bool {
		return active[i].EndTime < active[j].EndTime
	})

	return active
}

// GetContributionSummary generates a contribution summary for an idea
func (k *Keeper) GetContributionSummary(ideaID string) *types.ContributionSummary {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	summary := types.NewContributionSummary(ideaID)
	contributorWeights := make(map[string]float64)

	for _, c := range k.contributions {
		if c.IdeaID == ideaID && c.Status == types.ContributionStatusApproved {
			summary.TotalWeight += c.Weight
			summary.ByCategory[c.Category] += c.Weight
			contributorWeights[c.ContributorID] += c.Weight
		}
	}

	summary.TotalContributors = len(contributorWeights)

	// Get top contributors
	for contributorID, weight := range contributorWeights {
		stats := k.contributorStats[fmt.Sprintf("%s:%s", ideaID, contributorID)]
		if stats != nil {
			stats.TotalWeight = weight
			summary.TopContributors = append(summary.TopContributors, *stats)
		}
	}

	// Sort top contributors by weight
	sort.Slice(summary.TopContributors, func(i, j int) bool {
		return summary.TopContributors[i].TotalWeight > summary.TopContributors[j].TotalWeight
	})

	return summary
}

// DistributeRevenue distributes revenue to contributors
func (k *Keeper) DistributeRevenue(ideaID string, totalAmount uint64) (*types.RevenueDistribution, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, exists := k.ideas[ideaID]
	if !exists {
		return nil, fmt.Errorf("idea not found: %s", ideaID)
	}

	// Get contributor weights
	contributorWeights := make(map[string]float64)
	for _, c := range k.contributions {
		if c.IdeaID == ideaID && c.Status == types.ContributionStatusApproved {
			contributorWeights[c.ContributorID] += c.Weight
		}
	}

	distribution := types.NewRevenueDistribution(
		fmt.Sprintf("dist-%d", time.Now().Unix()),
		ideaID,
		totalAmount,
	)

	distribution.CalculateDistributions(contributorWeights)

	return distribution, nil
}

// AutoCloseExpiredCampaigns auto-closes expired campaigns
func (k *Keeper) AutoCloseExpiredCampaigns(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			k.mutex.Lock()
			for _, campaign := range k.campaigns {
				if campaign.IsActive() {
					campaign.CheckExpired()
				}
			}
			k.mutex.Unlock()
		}
	}
}

// GetCampaignStats returns campaign statistics
func (k *Keeper) GetCampaignStats(campaignID string) *types.CampaignStats {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	// Access backers directly without calling GetBackersByCampaign to avoid deadlock
	var backers []*types.Backer
	for _, b := range k.backers {
		if b.CampaignID == campaignID {
			backers = append(backers, b)
		}
	}
	if len(backers) == 0 {
		return &types.CampaignStats{CampaignID: campaignID}
	}

	var total, largest, smallest uint64
	smallest = backers[0].Amount

	for _, b := range backers {
		total += b.Amount
		if b.Amount > largest {
			largest = b.Amount
		}
		if b.Amount < smallest {
			smallest = b.Amount
		}
	}

	return &types.CampaignStats{
		CampaignID:           campaignID,
		AverageContribution:  float64(total) / float64(len(backers)),
		LargestContribution:  largest,
		SmallestContribution: smallest,
	}
}
