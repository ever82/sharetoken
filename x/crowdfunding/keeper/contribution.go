package keeper

import (
	"context"
	"fmt"
	"sort"
	"time"

	"sharetoken/x/crowdfunding/types"
	identitytypes "sharetoken/x/identity/types"
)

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
	ticker := time.NewTicker(time.Duration(identitytypes.CrowdfundingCheckIntervalMinutes) * time.Minute)
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
