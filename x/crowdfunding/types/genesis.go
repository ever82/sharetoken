package types

import (
	"fmt"
)

// GenesisState defines the crowdfunding module's genesis state.
type GenesisState struct {
	Ideas         []Idea               `json:"ideas"`
	Campaigns     []Campaign           `json:"campaigns"`
	Backers       []Backer             `json:"backers"`
	Contributions []Contribution       `json:"contributions"`
	Updates       []CampaignUpdate     `json:"updates"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Ideas:         []Idea{},
		Campaigns:     []Campaign{},
		Backers:       []Backer{},
		Contributions: []Contribution{},
		Updates:       []CampaignUpdate{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	// Validate ideas
	seenIdeaIDs := make(map[string]bool)
	for _, idea := range data.Ideas {
		if seenIdeaIDs[idea.ID] {
			return fmt.Errorf("duplicate idea ID: %s", idea.ID)
		}
		seenIdeaIDs[idea.ID] = true

		if err := idea.Validate(); err != nil {
			return fmt.Errorf("invalid idea %s: %w", idea.ID, err)
		}
	}

	// Validate campaigns
	seenCampaignIDs := make(map[string]bool)
	for _, campaign := range data.Campaigns {
		if seenCampaignIDs[campaign.ID] {
			return fmt.Errorf("duplicate campaign ID: %s", campaign.ID)
		}
		seenCampaignIDs[campaign.ID] = true

		if err := campaign.Validate(); err != nil {
			return fmt.Errorf("invalid campaign %s: %w", campaign.ID, err)
		}

		// Check that referenced idea exists
		if campaign.IdeaID != "" && !seenIdeaIDs[campaign.IdeaID] {
			return fmt.Errorf("campaign %s references non-existent idea %s", campaign.ID, campaign.IdeaID)
		}
	}

	// Validate backers
	seenBackerIDs := make(map[string]bool)
	for _, backer := range data.Backers {
		if seenBackerIDs[backer.ID] {
			return fmt.Errorf("duplicate backer ID: %s", backer.ID)
		}
		seenBackerIDs[backer.ID] = true

		// Check that referenced campaign exists
		if backer.CampaignID != "" && !seenCampaignIDs[backer.CampaignID] {
			return fmt.Errorf("backer %s references non-existent campaign %s", backer.ID, backer.CampaignID)
		}
	}

	// Validate contributions
	seenContributionIDs := make(map[string]bool)
	for _, contribution := range data.Contributions {
		if seenContributionIDs[contribution.ID] {
			return fmt.Errorf("duplicate contribution ID: %s", contribution.ID)
		}
		seenContributionIDs[contribution.ID] = true

		if err := contribution.Validate(); err != nil {
			return fmt.Errorf("invalid contribution %s: %w", contribution.ID, err)
		}

		// Check that referenced idea exists
		if contribution.IdeaID != "" && !seenIdeaIDs[contribution.IdeaID] {
			return fmt.Errorf("contribution %s references non-existent idea %s", contribution.ID, contribution.IdeaID)
		}
	}

	return nil
}
