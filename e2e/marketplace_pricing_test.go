package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// MarketplacePricingTestSuite 测试 ACH-USER-004: Transparent Service Pricing
type MarketplacePricingTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *MarketplacePricingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestBrowseAllServices 测试浏览所有可用服务
func (s *MarketplacePricingTestSuite) TestBrowseAllServices() {
	services, err := s.chain.QueryAllServices(100)
	require.NoError(s.T(), err, "服务查询应成功")
	require.Greater(s.T(), len(services), 0, "应至少有一个服务")

	s.T().Logf("市场中共有 %d 个服务", len(services))
}

// TestPricingModelDisplay 测试查看定价模式（按token/skill/打包）
func (s *MarketplacePricingTestSuite) TestPricingModelDisplay() {
	// 查询不同类型的服务
	testCases := []struct {
		level    int
		expected string
	}{
		{1, "per_token"},     // Level 1: 按token计费
		{2, "per_skill"},     // Level 2: 按skill计费
		{3, "fixed_package"}, // Level 3: 打包计费
	}

	for _, tc := range testCases {
		services, err := s.chain.QueryServicesByLevel(tc.level, 5)
		require.NoError(s.T(), err)

		if len(services) > 0 {
			require.Equal(s.T(), tc.expected, services[0].PricingModel,
				"Level %d 服务应有 %s 定价模式", tc.level, tc.expected)
		}
	}

	s.T().Log("定价模式显示正确")
}

// TestProviderComparison 测试比较不同提供者价格和评分
func (s *MarketplacePricingTestSuite) TestProviderComparison() {
	// 查询同一类型的多个提供者
	serviceType := "llm"
	providers, err := s.chain.QueryServiceProviders(serviceType, 10)
	require.NoError(s.T(), err)

	// 验证至少显示3个提供者
	require.GreaterOrEqual(s.T(), len(providers), 3,
		"同一服务应至少显示3个提供者")

	// 验证每个提供者都有评分信息
	for _, provider := range providers {
		require.NotZero(s.T(), provider.MQScore, "提供者应有MQ评分")
		require.NotZero(s.T(), provider.Price, "提供者应有价格")
	}

	s.T().Logf("找到 %d 个 %s 服务提供者", len(providers), serviceType)
}

// TestRatingDimensions 测试评分维度（质量/速度/完成率）
func (s *MarketplacePricingTestSuite) TestRatingDimensions() {
	// 查询服务提供者详情
	providers, err := s.chain.QueryServiceProviders("llm", 1)
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(providers), 0)

	provider := providers[0]

	// 验证评分维度
	require.NotNil(s.T(), provider.QualityScore, "应有质量评分")
	require.NotNil(s.T(), provider.ResponseSpeedScore, "应有响应速度评分")
	require.NotNil(s.T(), provider.CompletionRate, "应有完成率评分")

	s.T().Logf("提供者评分 - 质量: %.2f, 速度: %.2f, 完成率: %.2f%%",
		provider.QualityScore, provider.ResponseSpeedScore, provider.CompletionRate*100)
}

// TestCostEstimation 测试消费前预估费用
func (s *MarketplacePricingTestSuite) TestCostEstimation() {
	user := s.chain.CreateAccountWithBalance("estimate_user", "1000stt")

	// 预估费用
	params := map[string]interface{}{
		"service_type": "llm",
		"model":        "gpt-4",
		"input_tokens": 100,
	}

	estimate, err := s.chain.EstimateServiceCost(user, params)
	require.NoError(s.T(), err, "费用预估应成功")
	require.NotNil(s.T(), estimate, "预估结果不应为nil")
	require.NotZero(s.T(), estimate.TotalCost, "预估总费用应不为0")
	require.NotZero(s.T(), estimate.Breakdown.InputCost, "应有输入费用")
	require.NotZero(s.T(), estimate.Breakdown.OutputCost, "应有输出费用")

	s.T().Logf("费用预估: %f STT", estimate.TotalCost)
}

// TestCSVExport 测试导出CSV账单
func (s *MarketplacePricingTestSuite) TestCSVExport() {
	user := s.chain.CreateAccountWithBalance("export_user", "1000stt")

	// 先执行一些消费
	provider := s.chain.CreateAccount("export_provider")
	s.chain.CreateEscrow(user, provider.Address, "100stt", "test-service")

	// 导出账单
	csvData, err := s.chain.ExportBillingCSV(user.Address, "2024-01-01", "2024-12-31")
	require.NoError(s.T(), err, "CSV导出应成功")
	require.NotEmpty(s.T(), csvData, "CSV数据不应为空")

	// 验证CSV格式
	require.Contains(s.T(), csvData, "date", "CSV应包含日期列")
	require.Contains(s.T(), csvData, "service", "CSV应包含服务列")
	require.Contains(s.T(), csvData, "amount", "CSV应包含金额列")

	s.T().Log("CSV账单导出成功")
}

// TestServiceDetails 测试服务详情和案例
func (s *MarketplacePricingTestSuite) TestServiceDetails() {
	services, _ := s.chain.QueryAllServices(1)
	require.Greater(s.T(), len(services), 0)

	serviceID := services[0].ID

	// 查询服务详情
	detail, err := s.chain.QueryServiceDetail(serviceID)
	require.NoError(s.T(), err, "查询服务详情应成功")
	require.NotEmpty(s.T(), detail.Description, "服务应有描述")
	require.NotEmpty(s.T(), detail.Examples, "服务应有示例")

	s.T().Logf("服务 %s 详情查询成功", serviceID)
}

// TestServicePricingE2E 完整定价流程E2E测试
func (s *MarketplacePricingTestSuite) TestServicePricingE2E() {
	// 1. 浏览服务市场
	services, err := s.chain.QueryAllServices(10)
	require.NoError(s.T(), err)

	// 2. 查看定价模式
	if len(services) > 0 {
		detail, _ := s.chain.QueryServiceDetail(services[0].ID)
		require.NotEmpty(s.T(), detail.PricingModel)
	}

	// 3. 比较提供者
	providers, _ := s.chain.QueryServiceProviders("llm", 5)
	require.GreaterOrEqual(s.T(), len(providers), 3)

	// 4. 预估费用
	user := s.chain.CreateAccountWithBalance("pricing_e2e_user", "1000stt")
	estimate, _ := s.chain.EstimateServiceCost(user, map[string]interface{}{"service_type": "llm"})
	require.NotNil(s.T(), estimate)

	s.T().Log("定价流程E2E测试通过")
}

func TestMarketplacePricingSuite(t *testing.T) {
	suite.Run(t, new(MarketplacePricingTestSuite))
}
