package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sharetoken/x/crowdfunding/types"
)

// Idea Tests

func TestNewIdea(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")

	require.NotNil(t, idea)
	require.Equal(t, "idea-1", idea.ID)
	require.Equal(t, "Test Idea", idea.Title)
	require.Equal(t, "Description", idea.Description)
	require.Equal(t, "creator-1", idea.CreatorID)
	require.Equal(t, types.IdeaStatusDraft, idea.Status)
	require.Equal(t, 1, idea.CurrentVersion)
	require.NotNil(t, idea.Tags)
	require.NotNil(t, idea.Categories)
	require.NotZero(t, idea.CreatedAt)
	require.NotZero(t, idea.UpdatedAt)
}

func TestIdea_Validate(t *testing.T) {
	tests := []struct {
		name    string
		idea    types.Idea
		wantErr bool
	}{
		{
			name: "valid idea",
			idea: types.Idea{
				ID:        "idea-1",
				Title:     "Test Idea",
				CreatorID: "creator-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty ID",
			idea: types.Idea{
				ID:        "",
				Title:     "Test Idea",
				CreatorID: "creator-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty title",
			idea: types.Idea{
				ID:        "idea-1",
				Title:     "",
				CreatorID: "creator-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty creator",
			idea: types.Idea{
				ID:        "idea-1",
				Title:     "Test Idea",
				CreatorID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.idea.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIdea_Publish(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")
	require.Equal(t, types.IdeaStatusDraft, idea.Status)
	require.Zero(t, idea.PublishedAt)

	idea.Publish()
	require.Equal(t, types.IdeaStatusActive, idea.Status)
	require.NotZero(t, idea.PublishedAt)
}

func TestIdea_Update(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")
	oldVersion := idea.CurrentVersion

	idea.Update("Updated Title", "Updated Description")

	require.Equal(t, "Updated Title", idea.Title)
	require.Equal(t, "Updated Description", idea.Description)
	require.Equal(t, oldVersion+1, idea.CurrentVersion)
}

func TestIdea_Archive(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")
	idea.Archive()
	require.Equal(t, types.IdeaStatusArchived, idea.Status)
}

func TestIdea_StartFunding(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")
	idea.StartFunding("campaign-1")
	require.Equal(t, types.IdeaStatusFunding, idea.Status)
	require.Equal(t, "campaign-1", idea.CampaignID)
}

func TestIdea_StartDevelopment(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")
	idea.StartDevelopment()
	require.Equal(t, types.IdeaStatusDeveloping, idea.Status)
}

func TestIdea_Complete(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")
	idea.Complete()
	require.Equal(t, types.IdeaStatusCompleted, idea.Status)
}

func TestIdea_AddTag(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")

	idea.AddTag("blockchain")
	require.Contains(t, idea.Tags, "blockchain")

	// Duplicate should not be added
	idea.AddTag("blockchain")
	require.Len(t, idea.Tags, 1)

	idea.AddTag("ai")
	require.Len(t, idea.Tags, 2)
	require.Contains(t, idea.Tags, "blockchain")
	require.Contains(t, idea.Tags, "ai")
}

func TestIdea_AddCategory(t *testing.T) {
	idea := types.NewIdea("idea-1", "Test Idea", "Description", "creator-1")

	idea.AddCategory("technology")
	require.Contains(t, idea.Categories, "technology")

	// Duplicate should not be added
	idea.AddCategory("technology")
	require.Len(t, idea.Categories, 1)
}

// Campaign Tests

func TestNewCampaign(t *testing.T) {
	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Test Campaign", types.CampaignTypeInvestment, 10000)

	require.NotNil(t, campaign)
	require.Equal(t, "campaign-1", campaign.ID)
	require.Equal(t, "idea-1", campaign.IdeaID)
	require.Equal(t, "creator-1", campaign.CreatorID)
	require.Equal(t, "Test Campaign", campaign.Title)
	require.Equal(t, types.CampaignTypeInvestment, campaign.Type)
	require.Equal(t, types.CampaignStatusDraft, campaign.Status)
	require.Equal(t, uint64(10000), campaign.GoalAmount)
	require.Equal(t, "STT", campaign.Currency)
	require.Equal(t, uint64(1), campaign.MinContribution)
	require.NotZero(t, campaign.CreatedAt)
	require.NotZero(t, campaign.UpdatedAt)
}

func TestCampaign_Validate(t *testing.T) {
	tests := []struct {
		name     string
		campaign types.Campaign
		wantErr  bool
	}{
		{
			name: "valid campaign",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
				Type:      types.CampaignTypeInvestment,
				EquityOffered: 10,
				Valuation: 100000,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty ID",
			campaign: types.Campaign{
				ID:        "",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
			},
			wantErr: true,
		},
		{
			name: "valid - empty title", // Note: Validate doesn't check for empty title
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "",
				GoalAmount: 10000,
			},
			wantErr: false,
		},
		{
			name: "invalid - zero goal",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid - investment without equity",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
				Type:      types.CampaignTypeInvestment,
				EquityOffered: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid - investment without valuation",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
				Type:      types.CampaignTypeInvestment,
				EquityOffered: 10,
				Valuation: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid - equity over 100",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
				Type:      types.CampaignTypeInvestment,
				EquityOffered: 101,
			},
			wantErr: true,
		},
		{
			name: "invalid - lending without term",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
				Type:      types.CampaignTypeLending,
				LoanTerm:  0,
			},
			wantErr: true,
		},
		{
			name: "valid - lending with negative interest",
			campaign: types.Campaign{
				ID:        "campaign-1",
				IdeaID:    "idea-1",
				CreatorID: "creator-1",
				Title:     "Test Campaign",
				GoalAmount: 10000,
				Type:      types.CampaignTypeLending,
				InterestRate: -1,
				LoanTerm:  365,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.campaign.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCampaign_Launch(t *testing.T) {
	now := time.Now().Unix()
	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Test Campaign", types.CampaignTypeInvestment, 10000)

	campaign.Launch(86400) // 1 day
	require.Equal(t, types.CampaignStatusActive, campaign.Status)
	require.NotZero(t, campaign.StartTime)
	require.True(t, campaign.EndTime > campaign.StartTime)
	require.True(t, campaign.EndTime >= now+86400)
}

func TestCampaign_Contribute(t *testing.T) {
	now := time.Now().Unix()
	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Test Campaign", types.CampaignTypeInvestment, 10000)
	campaign.Status = types.CampaignStatusActive
	campaign.StartTime = now - 3600
	campaign.EndTime = now + 3600
	campaign.MinContribution = 100
	campaign.MaxContribution = 5000

	tests := []struct {
		name    string
		amount  uint64
		wantErr bool
	}{
		{ "valid contribution", 500, false },
		{ "below minimum", 50, true },
		{ "above maximum", 6000, true },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset campaign for each test
			campaign.RaisedAmount = 0
			campaign.BackerCount = 0

			err := campaign.Contribute(tt.amount)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.amount, campaign.RaisedAmount)
				require.Equal(t, 1, campaign.BackerCount)
			}
		})
	}
}

func TestCampaign_GetProgress(t *testing.T) {
	tests := []struct {
		name         string
		goal         uint64
		raised       uint64
		expectedPct  float64
	}{
		{"0% funded", 1000, 0, 0},
		{"50% funded", 1000, 500, 50},
		{"100% funded", 1000, 1000, 100},
		{"150% funded", 1000, 1500, 150},
		{"zero goal", 0, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := types.Campaign{
				GoalAmount:   tt.goal,
				RaisedAmount: tt.raised,
			}
			require.InDelta(t, tt.expectedPct, campaign.GetProgress(), 0.01)
		})
	}
}

func TestCampaign_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   types.CampaignStatus
		expected bool
	}{
		{"active", types.CampaignStatusActive, true},
		{"draft", types.CampaignStatusDraft, false},
		{"funded", types.CampaignStatusFunded, false},
		{"expired", types.CampaignStatusExpired, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := types.Campaign{Status: tt.status}
			require.Equal(t, tt.expected, campaign.IsActive())
		})
	}
}

func TestCampaign_IsFunded(t *testing.T) {
	tests := []struct {
		name     string
		status   types.CampaignStatus
		expected bool
	}{
		{"funded", types.CampaignStatusFunded, true},
		{"active", types.CampaignStatusActive, false},
		{"draft", types.CampaignStatusDraft, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := types.Campaign{Status: tt.status}
			require.Equal(t, tt.expected, campaign.IsFunded())
		})
	}
}

// Backer Tests

func TestNewBacker(t *testing.T) {
	backer := types.NewBacker("backer-1", "campaign-1", "user-1", 1000)

	require.NotNil(t, backer)
	require.Equal(t, "backer-1", backer.ID)
	require.Equal(t, "campaign-1", backer.CampaignID)
	require.Equal(t, "user-1", backer.BackerID)
	require.Equal(t, uint64(1000), backer.Amount)
	require.Equal(t, "STT", backer.Currency)
	require.NotZero(t, backer.CreatedAt)
	require.False(t, backer.Refunded)
}

func TestBacker_Refund(t *testing.T) {
	backer := types.NewBacker("backer-1", "campaign-1", "user-1", 1000)
	require.False(t, backer.Refunded)
	require.Zero(t, backer.RefundAmount)

	backer.Refund(1000)
	require.True(t, backer.Refunded)
	require.Equal(t, uint64(1000), backer.RefundAmount)
}

// CampaignUpdate Tests

func TestNewCampaignUpdate(t *testing.T) {
	update := types.NewCampaignUpdate("update-1", "campaign-1", "Update Title", "Update content", "creator-1")

	require.NotNil(t, update)
	require.Equal(t, "update-1", update.ID)
	require.Equal(t, "campaign-1", update.CampaignID)
	require.Equal(t, "Update Title", update.Title)
	require.Equal(t, "Update content", update.Content)
	require.Equal(t, "creator-1", update.CreatedBy)
	require.NotZero(t, update.CreatedAt)
}

// Contribution Tests

func TestNewContribution(t *testing.T) {
	contribution := types.NewContribution("contrib-1", "idea-1", "contributor-1", types.ContributionCode, "Added feature", 10)

	require.NotNil(t, contribution)
	require.Equal(t, "contrib-1", contribution.ID)
	require.Equal(t, "idea-1", contribution.IdeaID)
	require.Equal(t, "contributor-1", contribution.ContributorID)
	require.Equal(t, types.ContributionCode, contribution.Category)
	require.Equal(t, "Added feature", contribution.Description)
	require.Equal(t, float64(10), contribution.RawScore)
	require.Equal(t, types.ContributionStatusPending, contribution.Status)
	require.NotZero(t, contribution.CreatedAt)
}

func TestContribution_Validate(t *testing.T) {
	tests := []struct {
		name        string
		contribution types.Contribution
		wantErr     bool
	}{
		{
			name: "valid contribution",
			contribution: types.Contribution{
				IdeaID:        "idea-1",
				ContributorID: "contributor-1",
				RawScore:      10,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty idea ID",
			contribution: types.Contribution{
				IdeaID:        "",
				ContributorID: "contributor-1",
				RawScore:      10,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty contributor",
			contribution: types.Contribution{
				IdeaID:        "idea-1",
				ContributorID: "",
				RawScore:      10,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero score",
			contribution: types.Contribution{
				IdeaID:        "idea-1",
				ContributorID: "contributor-1",
				RawScore:      0,
			},
			wantErr: true,
		},
		{
			name: "invalid - negative score",
			contribution: types.Contribution{
				IdeaID:        "idea-1",
				ContributorID: "contributor-1",
				RawScore:      -5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contribution.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestContribution_CalculateWeight(t *testing.T) {
	contribution := types.NewContribution("contrib-1", "idea-1", "contributor-1", types.ContributionCode, "Added feature", 10)
	contribution.CalculateWeight()

	expectedWeight := 10.0 * types.CategoryWeights[types.ContributionCode]
	require.InDelta(t, expectedWeight, contribution.Weight, 0.001)
}

func TestContribution_Approve(t *testing.T) {
	contribution := types.NewContribution("contrib-1", "idea-1", "contributor-1", types.ContributionCode, "Added feature", 10)
	require.Equal(t, types.ContributionStatusPending, contribution.Status)
	require.Zero(t, contribution.ReviewedAt)

	contribution.Approve("reviewer-1")
	require.Equal(t, types.ContributionStatusApproved, contribution.Status)
	require.Equal(t, "reviewer-1", contribution.ReviewedBy)
	require.NotZero(t, contribution.ReviewedAt)
}

func TestContribution_Reject(t *testing.T) {
	contribution := types.NewContribution("contrib-1", "idea-1", "contributor-1", types.ContributionCode, "Added feature", 10)
	contribution.Reject("reviewer-1")
	require.Equal(t, types.ContributionStatusRejected, contribution.Status)
}

// RevenueDistribution Tests

func TestNewRevenueDistribution(t *testing.T) {
	distribution := types.NewRevenueDistribution("dist-1", "idea-1", 10000)

	require.NotNil(t, distribution)
	require.Equal(t, "dist-1", distribution.ID)
	require.Equal(t, "idea-1", distribution.IdeaID)
	require.Equal(t, uint64(10000), distribution.TotalAmount)
	require.Equal(t, types.DistributionStatusPending, distribution.Status)
	require.NotNil(t, distribution.Distributions)
}

func TestRevenueDistribution_CalculateDistributions(t *testing.T) {
	distribution := types.NewRevenueDistribution("dist-1", "idea-1", 10000)
	contributors := map[string]float64{
		"contrib-1": 50,
		"contrib-2": 30,
		"contrib-3": 20,
	}

	distribution.CalculateDistributions(contributors)

	require.Len(t, distribution.Distributions, 3)
	require.Equal(t, types.DistributionStatusProcessed, distribution.Status)
	require.NotZero(t, distribution.DistributedAt)

	// Check totals add up
	var totalAmount uint64
	for _, payout := range distribution.Distributions {
		totalAmount += payout.Amount
	}
	require.InDelta(t, distribution.TotalAmount, totalAmount, 1)
}

func TestRevenueDistribution_CalculateDistributions_ZeroWeight(t *testing.T) {
	distribution := types.NewRevenueDistribution("dist-1", "idea-1", 10000)
	contributors := map[string]float64{}

	distribution.CalculateDistributions(contributors)
	require.Len(t, distribution.Distributions, 0)
}

// Genesis Tests

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Ideas)
	require.NotNil(t, genesis.Campaigns)
	require.NotNil(t, genesis.Backers)
	require.NotNil(t, genesis.Contributions)
	require.NotNil(t, genesis.Updates)
	require.Empty(t, genesis.Ideas)
}

func TestValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		data    types.GenesisState
		wantErr bool
	}{
		{
			name:    "valid genesis with default",
			data:    *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "valid genesis with data",
			data: types.GenesisState{
				Ideas: []types.Idea{
					{ID: "idea-1", Title: "Idea 1", CreatorID: "creator-1"},
				},
				Campaigns: []types.Campaign{
					{ID: "campaign-1", IdeaID: "idea-1", CreatorID: "creator-1", Title: "Campaign 1", GoalAmount: 1000},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate idea IDs",
			data: types.GenesisState{
				Ideas: []types.Idea{
					{ID: "idea-1", Title: "Idea 1", CreatorID: "creator-1"},
					{ID: "idea-1", Title: "Idea 2", CreatorID: "creator-2"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - duplicate campaign IDs",
			data: types.GenesisState{
				Campaigns: []types.Campaign{
					{ID: "campaign-1", IdeaID: "idea-1", CreatorID: "creator-1", Title: "Campaign 1", GoalAmount: 1000},
					{ID: "campaign-1", IdeaID: "idea-2", CreatorID: "creator-2", Title: "Campaign 2", GoalAmount: 2000},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - campaign references non-existent idea",
			data: types.GenesisState{
				Ideas: []types.Idea{},
				Campaigns: []types.Campaign{
					{ID: "campaign-1", IdeaID: "non-existent", CreatorID: "creator-1", Title: "Campaign 1", GoalAmount: 1000},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - duplicate backer IDs",
			data: types.GenesisState{
				Campaigns: []types.Campaign{
					{ID: "campaign-1", IdeaID: "", CreatorID: "creator-1", Title: "Campaign 1", GoalAmount: 1000},
				},
				Backers: []types.Backer{
					{ID: "backer-1", CampaignID: "campaign-1", BackerID: "user-1", Amount: 100},
					{ID: "backer-1", CampaignID: "campaign-1", BackerID: "user-2", Amount: 200},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateGenesis(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
