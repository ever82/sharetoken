package keeper

import (
	"sync"

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
