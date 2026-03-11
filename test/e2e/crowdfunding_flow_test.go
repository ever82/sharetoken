package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CrowdfundingFlowTest 测试众筹完整流程
type CrowdfundingFlowTest struct {
	E2ETestSuite
}

func TestCrowdfundingFlow(t *testing.T) {
	suite.Run(t, new(CrowdfundingFlowTest))
}

// TestIdeaCreation 测试 Idea 创建
func (s *CrowdfundingFlowTest) TestIdeaCreation() {
	// 创建创作者
	creator := s.CreateVerifiedUser("github")

	// 创建 Idea
	idea := IdeaDraft{
		Title:       "Decentralized Social Network",
		Description: "A social network built on ShareToken",
		Category:    "social",
		Tags:        []string{"social", "decentralized", "web3"},
	}

	ideaID, err := s.CreateIdea(creator.Address, idea)
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), ideaID)

	// 验证 Idea
	storedIdea, err := s.GetIdea(ideaID)
	s.Require().NoError(err)
	assert.Equal(s.T(), idea.Title, storedIdea.Title)
	assert.Equal(s.T(), creator.Address, storedIdea.Creator)
	assert.Equal(s.T(), "draft", storedIdea.Status)
}

// TestIdeaVersioning 测试 Idea 版本管理
func (s *CrowdfundingFlowTest) TestIdeaVersioning() {
	// 创建 Idea
	creator := s.CreateVerifiedUser("github")
	idea := IdeaDraft{
		Title:       "Test Idea",
		Description: "Initial version",
	}

	ideaID, _ := s.CreateIdea(creator.Address, idea)

	// 更新 Idea（创建新版本）
	updatedIdea := IdeaDraft{
		Title:       "Test Idea",
		Description: "Updated version with more details",
	}

	err := s.UpdateIdea(creator.Address, ideaID, updatedIdea)
	s.Require().NoError(err)

	// 查询版本历史
	versions, err := s.GetIdeaVersions(ideaID)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(versions), 2)
}

// TestCrowdfundingLaunch 测试众筹发布
func (s *CrowdfundingFlowTest) TestCrowdfundingLaunch() {
	// 创建 Idea
	creator := s.CreateVerifiedUser("github")
	idea := IdeaDraft{
		Title:       "AI Marketplace",
		Description: "AI services marketplace",
	}

	ideaID, _ := s.CreateIdea(creator.Address, idea)

	// 启动投资型众筹
	crowdfunding := CrowdfundingConfig{
		Type:        "investment",
		Goal:        1000000000, // 1000 STT
		MinPledge:   10000000,   // 10 STT
		MaxPledge:   100000000,  // 100 STT
		Duration:    30,         // 30 days
		EquityShare: 20,         // 20% equity
	}

	campaignID, err := s.LaunchCrowdfunding(creator.Address, ideaID, crowdfunding)
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), campaignID)

	// 验证众筹状态
	campaign, err := s.GetCrowdfundingCampaign(campaignID)
	s.Require().NoError(err)
	assert.Equal(s.T(), "active", campaign.Status)
	assert.Equal(s.T(), crowdfunding.Goal, campaign.Goal)
}

// TestInvestmentCrowdfunding 测试投资型众筹
func (s *CrowdfundingFlowTest) TestInvestmentCrowdfunding() {
	// 创建众筹
	creator := s.CreateVerifiedUser("github")
	ideaID, _ := s.CreateIdea(creator.Address, IdeaDraft{Title: "Test"})

	campaignID, _ := s.LaunchCrowdfunding(creator.Address, ideaID, CrowdfundingConfig{
		Type:        "investment",
		Goal:        100000000,
		EquityShare: 10,
	})

	// 多个投资者参与
	investors := []struct {
		amount int64
	}{
		{20000000}, // 20 STT
		{30000000}, // 30 STT
		{50000000}, // 50 STT
	}

	for _, inv := range investors {
		investor := s.CreateVerifiedUser("github")
		s.FundAccount(investor.Address, inv.amount+1000000)

		err := s.PledgeToCrowdfunding(investor.Address, campaignID, inv.amount)
		s.Require().NoError(err)
	}

	// 验证众筹进度
	campaign, err := s.GetCrowdfundingCampaign(campaignID)
	s.Require().NoError(err)
	assert.Equal(s.T(), int64(100000000), campaign.Raised)
	assert.Equal(s.T(), "funded", campaign.Status)
}

// TestDonationCrowdfunding 测试捐赠型众筹
func (s *CrowdfundingFlowTest) TestDonationCrowdfunding() {
	// 创建捐赠型众筹
	creator := s.CreateVerifiedUser("github")
	ideaID, _ := s.CreateIdea(creator.Address, IdeaDraft{Title: "Open Source Tool"})

	campaignID, _ := s.LaunchCrowdfunding(creator.Address, ideaID, CrowdfundingConfig{
		Type:     "donation",
		Goal:     50000000, // 50 STT
		MinPledge: 1000000, // 1 STT
	})

	// 捐赠者参与
	donor := s.CreateVerifiedUser("github")
	s.FundAccount(donor.Address, 20000000)

	err := s.PledgeToCrowdfunding(donor.Address, campaignID, 10000000) // 10 STT
	s.Require().NoError(err)

	// 验证捐赠记录
	donations, err := s.GetCrowdfundingPledges(campaignID)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(donations), 1)
	assert.Equal(s.T(), donor.Address, donations[0].Contributor)
}

// TestCrowdfundingProgress 测试众筹进度跟踪
func (s *CrowdfundingFlowTest) TestCrowdfundingProgress() {
	// 创建众筹
	creator := s.CreateVerifiedUser("github")
	ideaID, _ := s.CreateIdea(creator.Address, IdeaDraft{Title: "Test"})

	campaignID, _ := s.LaunchCrowdfunding(creator.Address, ideaID, CrowdfundingConfig{
		Type: "investment",
		Goal: 100000000, // 100 STT
	})

	// 部分投资
	investor := s.CreateVerifiedUser("github")
	s.FundAccount(investor.Address, 30000000)
	s.PledgeToCrowdfunding(investor.Address, campaignID, 25000000) // 25 STT

	// 验证进度
	progress, err := s.GetCrowdfundingProgress(campaignID)
	s.Require().NoError(err)
	assert.Equal(s.T(), int64(25000000), progress.Raised)
	assert.Equal(s.T(), float64(25), progress.Percentage) // 25%
}

// TestContributionTracking 测试贡献追踪
func (s *CrowdfundingFlowTest) TestContributionTracking() {
	// 创建协作式 Idea
	creator := s.CreateVerifiedUser("github")
	ideaID, _ := s.CreateIdea(creator.Address, IdeaDraft{Title: "Collaborative Project"})

	// 多个贡献者协作
	contributors := []struct {
		address     string
		contribType string
		weight      int64
	}{
		{s.CreateVerifiedUser("github").Address, "code", 40},
		{s.CreateVerifiedUser("github").Address, "design", 30},
		{s.CreateVerifiedUser("github").Address, "docs", 20},
	}

	for _, c := range contributors {
		err := s.RecordContribution(ideaID, c.address, c.contribType, c.weight)
		s.Require().NoError(err)
	}

	// 验证贡献权重
	weights, err := s.GetContributionWeights(ideaID)
	s.Require().NoError(err)
	assert.Equal(s.T(), int64(90), weights.Total) // 40 + 30 + 20
}

// TestRevenueDistribution 测试收益分配
func (s *CrowdfundingFlowTest) TestRevenueDistribution() {
	// 创建并完成众筹
	creator := s.CreateVerifiedUser("github")
	ideaID, _ := s.CreateIdea(creator.Address, IdeaDraft{Title: "Revenue Project"})

	// 记录贡献
	contributor := s.CreateVerifiedUser("github")
	s.RecordContribution(ideaID, contributor.Address, "code", 50)
	s.RecordContribution(ideaID, creator.Address, "idea", 50)

	// 众筹成功并产生收益
	campaignID, _ := s.LaunchCrowdfunding(creator.Address, ideaID, CrowdfundingConfig{
		Type:        "investment",
		Goal:        100000000,
		EquityShare: 20,
	})

	investor := s.CreateVerifiedUser("github")
	s.FundAccount(investor.Address, 100000000)
	s.PledgeToCrowdfunding(investor.Address, campaignID, 100000000)

	// 项目产生收益
	revenue := int64(50000000) // 50 STT
	err := s.DistributeRevenue(creator.Address, ideaID, revenue)
	s.Require().NoError(err)

	// 验证收益分配
	distributions, err := s.GetRevenueDistributions(ideaID)
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), distributions)
}

// TestFailedCrowdfundingRefund 测试失败众筹退款
func (s *CrowdfundingFlowTest) TestFailedCrowdfundingRefund() {
	// 创建众筹
	creator := s.CreateVerifiedUser("github")
	ideaID, _ := s.CreateIdea(creator.Address, IdeaDraft{Title: "Test"})

	campaignID, _ := s.LaunchCrowdfunding(creator.Address, ideaID, CrowdfundingConfig{
		Type:     "investment",
		Goal:     100000000,
		Duration: 7, // 7 days
	})

	// 投资者参与但未达目标
	investor := s.CreateVerifiedUser("github")
	initialBalance := int64(50000000)
	s.FundAccount(investor.Address, initialBalance)
	pledgeAmount := int64(20000000)
	s.PledgeToCrowdfunding(investor.Address, campaignID, pledgeAmount)

	// 模拟众筹失败（到期未达标）
	err := s.ExpireCrowdfunding(campaignID)
	s.Require().NoError(err)

	// 申请退款
	err = s.RefundCrowdfundingPledge(investor.Address, campaignID)
	s.Require().NoError(err)

	// 验证退款
	finalBalance, _ := s.QueryBalance(investor.Address)
	assert.GreaterOrEqual(s.T(), finalBalance, initialBalance-1000000) // 考虑 gas 费
}

// Helper types

type IdeaDraft struct {
	Title       string
	Description string
	Category    string
	Tags        []string
}

type CrowdfundingConfig struct {
	Type        string // investment, loan, donation
	Goal        int64
	MinPledge   int64
	MaxPledge   int64
	Duration    int    // days
	EquityShare int    // for investment type
}

type CampaignProgress struct {
	Raised     int64
	Percentage float64
}

type ContributionRecord struct {
	Contributor string
	Type        string
	Weight      int64
}
