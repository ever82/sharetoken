package keeper

import (
	"fmt"
	"sort"
	"time"

	"sharetoken/x/crowdfunding/types"
)

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
	defer k.mutex.RUnlock()

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
