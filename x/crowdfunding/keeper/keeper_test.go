package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/crowdfunding/types"
)

func TestIdeaCreation(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Decentralized Social Network", "A social network built on blockchain", "creator-1")
	idea.AddTag("social")
	idea.AddTag("blockchain")
	idea.AddCategory("technology")

	err := k.CreateIdea(idea)
	require.NoError(t, err)

	retrieved := k.GetIdea("idea-1")
	require.NotNil(t, retrieved)
	require.Equal(t, "Decentralized Social Network", retrieved.Title)
	require.Equal(t, "creator-1", retrieved.CreatorID)
	require.Equal(t, types.IdeaStatusDraft, retrieved.Status)
	require.Contains(t, retrieved.Tags, "social")
	require.Contains(t, retrieved.Categories, "technology")
}

func TestIdeaValidation(t *testing.T) {
	k := NewKeeper()

	// Invalid - no title
	invalidIdea := types.NewIdea("idea-1", "", "Description", "creator-1")
	err := k.CreateIdea(invalidIdea)
	require.Error(t, err)

	// Invalid - no creator
	invalidIdea2 := types.NewIdea("idea-2", "Title", "Description", "")
	err = k.CreateIdea(invalidIdea2)
	require.Error(t, err)
}

func TestIdeaLifecycle(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Publish
	err := k.PublishIdea("idea-1")
	require.NoError(t, err)
	require.Equal(t, types.IdeaStatusActive, idea.Status)
	require.Greater(t, idea.PublishedAt, int64(0))

	// Archive
	idea.Archive()
	require.Equal(t, types.IdeaStatusArchived, idea.Status)
}

func TestIdeaVersioning(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Create initial version
	version1, err := k.UpdateIdea("idea-1", "Title", "Description", "Initial version", "creator-1")
	require.NoError(t, err)
	require.Equal(t, 1, version1.Version)
	require.Equal(t, "Initial version", version1.Changes)

	// Update again
	version2, err := k.UpdateIdea("idea-1", "New Title", "New Description", "Updated content", "creator-1")
	require.NoError(t, err)
	require.Equal(t, 2, version2.Version)

	// Check idea updated
	idea = k.GetIdea("idea-1")
	require.Equal(t, "New Title", idea.Title)
	require.Equal(t, "New Description", idea.Description)
	require.Equal(t, 3, idea.CurrentVersion) // 1 -> 2 -> 3

	// Get versions
	versions := k.GetIdeaVersions("idea-1")
	require.Len(t, versions, 2)
}

func TestContributionSubmission(t *testing.T) {
	k := NewKeeper()

	// Create idea
	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Submit contribution
	contribution := types.NewContribution("contrib-1", "idea-1", "contributor-1", types.ContributionCode, "Implemented core feature", 100)
	err := k.SubmitContribution(contribution)
	require.NoError(t, err)

	// Check weight calculated
	expectedWeight := 100.0 * types.CategoryWeights[types.ContributionCode] // 100 * 1.5 = 150
	require.Equal(t, expectedWeight, contribution.Weight)

	// Check idea stats updated
	idea = k.GetIdea("idea-1")
	require.Equal(t, 1, idea.ContributionCount)
	require.Equal(t, expectedWeight, idea.TotalWeight)

	// Check contributor stats
	stats := k.GetContributorStats("idea-1", "contributor-1")
	require.Equal(t, 1, stats.PendingCount)
}

func TestContributionApproval(t *testing.T) {
	k := NewKeeper()

	// Create idea and contribution
	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	contribution := types.NewContribution("contrib-1", "idea-1", "contributor-1", types.ContributionCode, "Code", 100)
	k.SubmitContribution(contribution)

	// Approve contribution
	err := k.ApproveContribution("contrib-1", "creator-1")
	require.NoError(t, err)
	require.Equal(t, types.ContributionStatusApproved, contribution.Status)
	require.Equal(t, "creator-1", contribution.ReviewedBy)
	require.Greater(t, contribution.ReviewedAt, int64(0))

	// Check contributor stats updated
	stats := k.GetContributorStats("idea-1", "contributor-1")
	require.Equal(t, 1, stats.ApprovedCount)
	require.Equal(t, 0, stats.PendingCount)
	require.Equal(t, contribution.Weight, stats.TotalWeight)
}

func TestContributionCategories(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Test different categories
	testCases := []struct {
		category     types.ContributionCategory
		rawScore     float64
		expectedMult float64
	}{
		{types.ContributionCode, 100, 1.5},
		{types.ContributionDesign, 100, 1.2},
		{types.ContributionDocs, 100, 1.0},
		{types.ContributionResearch, 100, 1.3},
	}

	for i, tc := range testCases {
		contrib := types.NewContribution(
			fmt.Sprintf("contrib-%d", i),
			"idea-1",
			"contributor-1",
			tc.category,
			"Test",
			tc.rawScore,
		)
		k.SubmitContribution(contrib)
		k.ApproveContribution(contrib.ID, "creator-1")

		expectedWeight := tc.rawScore * tc.expectedMult
		require.Equal(t, expectedWeight, contrib.Weight)
	}
}

func TestCampaignCreation(t *testing.T) {
	k := NewKeeper()

	// Create idea first
	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Create investment campaign
	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Social Network Funding", types.CampaignTypeInvestment, 100000)
	campaign.EquityOffered = 20.0
	campaign.Valuation = 500000

	err := k.CreateCampaign(campaign)
	require.NoError(t, err)

	retrieved := k.GetCampaign("campaign-1")
	require.NotNil(t, retrieved)
	require.Equal(t, types.CampaignTypeInvestment, retrieved.Type)
	require.Equal(t, uint64(100000), retrieved.GoalAmount)
}

func TestCampaignTypes(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Investment campaign
	investment := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Investment", types.CampaignTypeInvestment, 100000)
	investment.EquityOffered = 10.0
	investment.Valuation = 1000000
	err := k.CreateCampaign(investment)
	require.NoError(t, err)

	// Lending campaign
	lending := types.NewCampaign("campaign-2", "idea-1", "creator-1", "Lending", types.CampaignTypeLending, 50000)
	lending.InterestRate = 10.0
	lending.LoanTerm = 365
	err = k.CreateCampaign(lending)
	require.NoError(t, err)

	// Donation campaign
	donation := types.NewCampaign("campaign-3", "idea-1", "creator-1", "Donation", types.CampaignTypeDonation, 10000)
	err = k.CreateCampaign(donation)
	require.NoError(t, err)
}

func TestCampaignLifecycle(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Funding", types.CampaignTypeInvestment, 10000)
	campaign.EquityOffered = 10.0
	campaign.Valuation = 100000
	k.CreateCampaign(campaign)

	// Launch
	err := k.LaunchCampaign("campaign-1", 86400) // 1 day
	require.NoError(t, err)
	require.Equal(t, types.CampaignStatusActive, campaign.Status)
	require.Greater(t, campaign.StartTime, int64(0))
	require.Greater(t, campaign.EndTime, campaign.StartTime)

	// Idea status updated
	idea = k.GetIdea("idea-1")
	require.Equal(t, types.IdeaStatusFunding, idea.Status)

	// Contribute
	backer1 := types.NewBacker("backer-1", "campaign-1", "user-1", 5000)
	err = k.ContributeToCampaign(backer1)
	require.NoError(t, err)
	require.Equal(t, uint64(5000), campaign.RaisedAmount)
	require.Equal(t, 1, campaign.BackerCount)

	// Contribute more to reach goal
	backer2 := types.NewBacker("backer-2", "campaign-1", "user-2", 5000)
	err = k.ContributeToCampaign(backer2)
	require.NoError(t, err)
	require.Equal(t, uint64(10000), campaign.RaisedAmount)
	require.Equal(t, types.CampaignStatusFunded, campaign.Status)
	require.Equal(t, 100.0, campaign.GetProgress())

	// Close campaign
	err = k.CloseCampaign("campaign-1")
	require.NoError(t, err)
	require.Equal(t, types.CampaignStatusClosed, campaign.Status)

	// Idea status updated
	idea = k.GetIdea("idea-1")
	require.Equal(t, types.IdeaStatusDeveloping, idea.Status)
}

func TestContributionLimits(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Funding", types.CampaignTypeDonation, 10000)
	campaign.MinContribution = 100
	campaign.MaxContribution = 5000
	k.CreateCampaign(campaign)
	k.LaunchCampaign("campaign-1", 86400)

	// Below minimum
	backer1 := types.NewBacker("backer-1", "campaign-1", "user-1", 50)
	err := k.ContributeToCampaign(backer1)
	require.Error(t, err)

	// Above maximum
	backer2 := types.NewBacker("backer-2", "campaign-1", "user-2", 6000)
	err = k.ContributeToCampaign(backer2)
	require.Error(t, err)

	// Within limits
	backer3 := types.NewBacker("backer-3", "campaign-1", "user-3", 1000)
	err = k.ContributeToCampaign(backer3)
	require.NoError(t, err)
}

func TestExpectedReturns(t *testing.T) {
	// Investment
	investment := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Investment", types.CampaignTypeInvestment, 100000)
	investment.EquityOffered = 20.0
	investment.Valuation = 500000

	// 10000 investment gets 2% equity (200 basis points)
	returned, err := investment.GetExpectedReturn(10000)
	require.NoError(t, err)
	require.Equal(t, uint64(200), returned) // 2% = 200 basis points

	// Lending
	lending := types.NewCampaign("campaign-2", "idea-1", "creator-1", "Lending", types.CampaignTypeLending, 50000)
	lending.InterestRate = 10.0
	lending.LoanTerm = 365

	// 10000 loan at 10% for 1 year = 11000
	returned, err = lending.GetExpectedReturn(10000)
	require.NoError(t, err)
	require.Equal(t, uint64(11000), returned)

	// Donation
	donation := types.NewCampaign("campaign-3", "idea-1", "creator-1", "Donation", types.CampaignTypeDonation, 10000)
	returned, err = donation.GetExpectedReturn(10000)
	require.NoError(t, err)
	require.Equal(t, uint64(0), returned)
}

func TestContributionSummary(t *testing.T) {
	k := NewKeeper()

	// Create idea
	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Add contributions from different contributors
	contributions := []struct {
		id          string
		contributor string
		category    types.ContributionCategory
		score       float64
	}{
		{"c1", "user-1", types.ContributionCode, 100},
		{"c2", "user-1", types.ContributionDesign, 50},
		{"c3", "user-2", types.ContributionCode, 80},
		{"c4", "user-3", types.ContributionDocs, 100},
	}

	for _, tc := range contributions {
		c := types.NewContribution(tc.id, "idea-1", tc.contributor, tc.category, "Test", tc.score)
		k.SubmitContribution(c)
		k.ApproveContribution(tc.id, "creator-1")
	}

	// Get summary
	summary := k.GetContributionSummary("idea-1")
	require.Equal(t, 3, summary.TotalContributors)
	require.Greater(t, summary.TotalWeight, float64(0))
	require.Greater(t, summary.ByCategory[types.ContributionCode], float64(0))
	require.Greater(t, summary.ByCategory[types.ContributionDesign], float64(0))
	require.Greater(t, summary.ByCategory[types.ContributionDocs], float64(0))
}

func TestRevenueDistribution(t *testing.T) {
	k := NewKeeper()

	// Create idea and contributions
	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// User 1: Code contribution (high weight)
	c1 := types.NewContribution("c1", "idea-1", "user-1", types.ContributionCode, "Code", 100)
	k.SubmitContribution(c1)
	k.ApproveContribution("c1", "creator-1")

	// User 2: Docs contribution (lower weight)
	c2 := types.NewContribution("c2", "idea-1", "user-2", types.ContributionDocs, "Docs", 100)
	k.SubmitContribution(c2)
	k.ApproveContribution("c2", "creator-1")

	// Distribute 10000 STT
	distribution, err := k.DistributeRevenue("idea-1", 10000)
	require.NoError(t, err)
	require.Len(t, distribution.Distributions, 2)
	require.Equal(t, types.DistributionStatusProcessed, distribution.Status)

	// Check user 1 gets more (code has higher weight)
	var user1Share, user2Share float64
	for _, d := range distribution.Distributions {
		if d.ContributorID == "user-1" {
			user1Share = d.Share
		}
		if d.ContributorID == "user-2" {
			user2Share = d.Share
		}
	}

	require.Greater(t, user1Share, user2Share)
	require.Equal(t, float64(100), user1Share+user2Share)
}

func TestCampaignUpdates(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Funding", types.CampaignTypeDonation, 10000)
	k.CreateCampaign(campaign)
	k.LaunchCampaign("campaign-1", 86400)

	// Add update
	update := types.NewCampaignUpdate("update-1", "campaign-1", "Milestone Reached", "We hit 50% funding!", "creator-1")
	err := k.AddCampaignUpdate(update)
	require.NoError(t, err)
	require.Equal(t, 1, campaign.UpdateCount)

	// Add another update
	update2 := types.NewCampaignUpdate("update-2", "campaign-1", "New Feature", "Added feature X", "creator-1")
	k.AddCampaignUpdate(update2)
	require.Equal(t, 2, campaign.UpdateCount)

	// Get updates
	updates := k.GetCampaignUpdates("campaign-1")
	require.Len(t, updates, 2)
}

func TestActiveCampaigns(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	// Create active campaign
	campaign1 := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Active", types.CampaignTypeDonation, 10000)
	k.CreateCampaign(campaign1)
	k.LaunchCampaign("campaign-1", 86400)

	// Create draft campaign (not active)
	campaign2 := types.NewCampaign("campaign-2", "idea-1", "creator-1", "Draft", types.CampaignTypeDonation, 5000)
	k.CreateCampaign(campaign2)

	// Get active campaigns
	active := k.GetActiveCampaigns()
	require.Len(t, active, 1)
	require.Equal(t, "campaign-1", active[0].ID)
}

func TestCampaignStats(t *testing.T) {
	k := NewKeeper()

	idea := types.NewIdea("idea-1", "Title", "Description", "creator-1")
	k.CreateIdea(idea)

	campaign := types.NewCampaign("campaign-1", "idea-1", "creator-1", "Funding", types.CampaignTypeDonation, 10000)
	k.CreateCampaign(campaign)
	k.LaunchCampaign("campaign-1", 86400)

	// Add backers
	backers := []struct {
		id     string
		amount uint64
	}{
		{"backer-1", 100},
		{"backer-2", 500},
		{"backer-3", 1000},
	}

	for _, b := range backers {
		backer := types.NewBacker(b.id, "campaign-1", b.id, b.amount)
		k.ContributeToCampaign(backer)
	}

	stats := k.GetCampaignStats("campaign-1")
	require.Equal(t, uint64(1000), stats.LargestContribution)
	require.Equal(t, uint64(100), stats.SmallestContribution)
	require.Equal(t, 533.33, float64(int(stats.AverageContribution*100))/100) // 1600/3
}
