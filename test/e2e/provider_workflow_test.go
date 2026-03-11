package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ProviderWorkflowTest 测试服务提供者完整流程
type ProviderWorkflowTest struct {
	E2ETestSuite
}

func TestProviderWorkflow(t *testing.T) {
	suite.Run(t, new(ProviderWorkflowTest))
}

// TestProviderRegistration 测试服务提供者注册
func (s *ProviderWorkflowTest) TestProviderRegistration() {
	// 创建提供者用户
	provider := s.CreateVerifiedUser("github")

	// 注册为服务提供者
	providerID, err := s.RegisterAsProvider(provider.Address, "Test Provider", "A test service provider")
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), providerID)

	// 查询提供者信息
	info, err := s.GetProviderInfo(providerID)
	s.Require().NoError(err)
	assert.Equal(s.T(), provider.Address, info.Owner)
	assert.Equal(s.T(), "Test Provider", info.Name)
	assert.True(s.T(), info.IsActive)
}

// TestLLMServiceHosting 测试 LLM 服务托管
func (s *ProviderWorkflowTest) TestLLMServiceHosting() {
	// 创建提供者
	provider := s.CreateVerifiedUser("github")
	providerID, _ := s.RegisterAsProvider(provider.Address, "LLM Provider", "")

	// 托管 OpenAI API Key
	apiKeyConfig := LLMAPIKeyConfig{
		Provider:    "openai",
		APIKey:      "sk-test123456",
		Model:       "gpt-4",
		PricePerToken: 100, // 100 ustt per 1k tokens
	}

	keyID, err := s.HostLLMAPIKey(provider.Address, apiKeyConfig)
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), keyID)

	// 验证 API Key 托管成功
	hostedKey, err := s.GetHostedAPIKey(keyID)
	s.Require().NoError(err)
	assert.Equal(s.T(), "openai", hostedKey.Provider)
	assert.Equal(s.T(), "gpt-4", hostedKey.Model)

	// 验证加密存储（不能获取明文 key）
	assert.Empty(s.T(), hostedKey.PlainKey)
	assert.NotEmpty(s.T(), hostedKey.EncryptedKey)
	assert.NotEmpty(s.T(), hostedKey.Hash)
}

// TestAgentSkillRegistration 测试 Agent 技能注册
func (s *ProviderWorkflowTest) TestAgentSkillRegistration() {
	// 创建提供者
	provider := s.CreateVerifiedUser("github")
	providerID, _ := s.RegisterAsProvider(provider.Address, "Agent Provider", "")

	// 注册 Agent 技能
	skill := AgentSkill{
		Name:        "Code Review",
		Description: "Automated code review and suggestions",
		Category:    "development",
		PricePerUse: 500000, // 0.5 STT
		InputSchema: map[string]interface{}{
			"code":     "string",
			"language": "string",
		},
		OutputSchema: map[string]interface{}{
			"suggestions": []string{},
			"score":       "number",
		},
	}

	skillID, err := s.RegisterAgentSkill(provider.Address, skill)
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), skillID)

	// 验证技能注册
	registeredSkill, err := s.GetAgentSkill(skillID)
	s.Require().NoError(err)
	assert.Equal(s.T(), "Code Review", registeredSkill.Name)
	assert.Equal(s.T(), "development", registeredSkill.Category)
}

// TestWorkflowTaskAcceptance 测试 Workflow 任务承接
func (s *ProviderWorkflowTest) TestWorkflowTaskAcceptance() {
	// 创建用户和提供者
	user := s.CreateVerifiedUser("github")
	provider := s.CreateVerifiedUser("github")
	_, _ = s.RegisterAsProvider(provider.Address, "Workflow Provider", "")

	// 创建 workflow 任务
	task := WorkflowTask{
		Title:       "Build a website",
		Description: "Create a simple portfolio website",
		Budget:      10000000, // 10 STT
		Milestones: []Milestone{
			{Description: "Design mockup", Percentage: 20},
			{Description: "Frontend development", Percentage: 40},
			{Description: "Backend integration", Percentage: 30},
			{Description: "Deployment", Percentage: 10},
		},
	}

	taskID, err := s.CreateWorkflowTask(user.Address, task)
	s.Require().NoError(err)

	// 提供者承接任务
	err = s.AcceptWorkflowTask(provider.Address, taskID)
	s.Require().NoError(err)

	// 验证任务状态
	taskInfo, err := s.GetWorkflowTask(taskID)
	s.Require().NoError(err)
	assert.Equal(s.T(), provider.Address, taskInfo.AssignedTo)
	assert.Equal(s.T(), "accepted", taskInfo.Status)
}

// TestTaskDeliveryAndPayment 测试任务交付和收款
func (s *ProviderWorkflowTest) TestTaskDeliveryAndPayment() {
	// 创建并承接任务
	user := s.CreateVerifiedUser("github")
	provider := s.CreateVerifiedUser("github")
	_, _ = s.RegisterAsProvider(provider.Address, "Provider", "")

	task := WorkflowTask{
		Title:       "Simple task",
		Description: "A simple test task",
		Budget:      1000000, // 1 STT
		Milestones: []Milestone{
			{Description: "Complete", Percentage: 100},
		},
	}

	taskID, _ := s.CreateWorkflowTask(user.Address, task)
	s.FundEscrowForTask(user.Address, taskID, 1000000)
	s.AcceptWorkflowTask(provider.Address, taskID)

	// 完成里程碑
	err := s.CompleteMilestone(provider.Address, taskID, 0)
	s.Require().NoError(err)

	// 用户确认交付
	err = s.ConfirmDelivery(user.Address, taskID)
	s.Require().NoError(err)

	// 验证资金释放给提供者
	providerBalance, err := s.QueryBalance(provider.Address)
	s.Require().NoError(err)
	// 应该收到 1 STT 减去平台费用
	assert.Greater(s.T(), providerBalance, int64(0))
}

// TestEarningsOverview 测试收入概览
func (s *ProviderWorkflowTest) TestEarningsOverview() {
	// 创建提供者并执行一些任务
	provider := s.CreateVerifiedUser("github")
	_, _ = s.RegisterAsProvider(provider.Address, "Provider", "")

	// 查询收入概览
	earnings, err := s.GetProviderEarnings(provider.Address)
	s.Require().NoError(err)

	// 验证字段
	assert.GreaterOrEqual(s.T(), earnings.TotalEarned, int64(0))
	assert.GreaterOrEqual(s.T(), earnings.TotalTasks, int64(0))
	assert.GreaterOrEqual(s.T(), earnings.PendingPayment, int64(0))
}

// TestCustomerReviews 测试客户评价
func (s *ProviderWorkflowTest) TestCustomerReviews() {
	// 创建并完成任务
	user := s.CreateVerifiedUser("github")
	provider := s.CreateVerifiedUser("github")
	_, _ = s.RegisterAsProvider(provider.Address, "Provider", "")

	task := WorkflowTask{
		Title:       "Review task",
		Description: "Test review",
		Budget:      1000000,
	}

	taskID, _ := s.CreateWorkflowTask(user.Address, task)
	s.FundEscrowForTask(user.Address, taskID, 1000000)
	s.AcceptWorkflowTask(provider.Address, taskID)
	s.CompleteMilestone(provider.Address, taskID, 0)
	s.ConfirmDelivery(user.Address, taskID)

	// 提交评价
	review := TaskReview{
		TaskID:  taskID,
		Rating:  5,
		Comment: "Excellent work!",
		Aspects: ReviewAspects{
			Quality:   5,
			Communication: 5,
			Timeliness: 5,
			Professionalism: 5,
		},
	}

	err := s.SubmitTaskReview(user.Address, review)
	s.Require().NoError(err)

	// 验证评价
	reviews, err := s.GetProviderReviews(provider.Address)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(reviews), 1)
	assert.Equal(s.T(), 5, reviews[0].Rating)
}

// TestProviderWithdrawal 测试提供者提现
func (s *ProviderWorkflowTest) TestProviderWithdrawal() {
	// 创建有收入的提供者
	provider := s.CreateVerifiedUser("github")
	s.FundAccount(provider.Address, 10000000) // 10 STT

	// 查询可提现金额
	withdrawable, err := s.GetWithdrawableAmount(provider.Address)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), withdrawable, int64(10000000))

	// 执行提现
	withdrawAmount := int64(5000000) // 5 STT
	txHash, err := s.WithdrawFunds(provider.Address, withdrawAmount, "external_wallet_address")
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), txHash)

	// 验证余额减少
	newBalance, err := s.QueryBalance(provider.Address)
	s.Require().NoError(err)
	assert.Less(s.T(), newBalance, int64(10000000))
}

// TestServicePricingUpdate 测试服务定价更新
func (s *ProviderWorkflowTest) TestServicePricingUpdate() {
	// 创建提供者并托管服务
	provider := s.CreateVerifiedUser("github")
	providerID, _ := s.RegisterAsProvider(provider.Address, "Provider", "")

	apiKeyConfig := LLMAPIKeyConfig{
		Provider:      "openai",
		APIKey:        "sk-test123",
		Model:         "gpt-4",
		PricePerToken: 100,
	}

	keyID, _ := s.HostLLMAPIKey(provider.Address, apiKeyConfig)

	// 更新定价
	newPrice := int64(150) // 涨价到 150 ustt
	err := s.UpdateServicePricing(provider.Address, keyID, newPrice)
	s.Require().NoError(err)

	// 验证新定价
	service, err := s.GetServiceDetails(keyID)
	s.Require().NoError(err)
	assert.Equal(s.T(), newPrice, service.Price)
}

// TestProviderDeactivation 测试提供者下线
func (s *ProviderWorkflowTest) TestProviderDeactivation() {
	// 创建提供者
	provider := s.CreateVerifiedUser("github")
	providerID, _ := s.RegisterAsProvider(provider.Address, "Provider", "")

	// 下线提供者
	err := s.DeactivateProvider(provider.Address, providerID)
	s.Require().NoError(err)

	// 验证状态
	info, err := s.GetProviderInfo(providerID)
	s.Require().NoError(err)
	assert.False(s.T(), info.IsActive)
}

// Helper types

type LLMAPIKeyConfig struct {
	Provider      string
	APIKey        string
	Model         string
	PricePerToken int64
}

type AgentSkill struct {
	Name         string
	Description  string
	Category     string
	PricePerUse  int64
	InputSchema  map[string]interface{}
	OutputSchema map[string]interface{}
}

type WorkflowTask struct {
	Title       string
	Description string
	Budget      int64
	Milestones  []Milestone
}

type Milestone struct {
	Description string
	Percentage  int
}

type TaskReview struct {
	TaskID  string
	Rating  int
	Comment string
	Aspects ReviewAspects
}

type ReviewAspects struct {
	Quality         int
	Communication   int
	Timeliness      int
	Professionalism int
}

type ProviderEarnings struct {
	TotalEarned    int64
	TotalTasks     int64
	PendingPayment int64
}
