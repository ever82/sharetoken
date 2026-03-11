package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// GenieBotTestSuite 测试 ACH-USER-002: One-Click AI Access
type GenieBotTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *GenieBotTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestServiceDiscovery 测试服务发现接口
func (s *GenieBotTestSuite) TestServiceDiscovery() {
	// 查询可用服务列表
	services, err := s.chain.QueryServices(5)
	require.NoError(s.T(), err, "服务查询应成功")
	require.GreaterOrEqual(s.T(), len(services), 5, "应至少有5种服务")

	s.T().Logf("发现 %d 个服务", len(services))
}

// TestIntentRecognitionLatency 测试意图识别响应时间 (< 3s)
func (s *GenieBotTestSuite) TestIntentRecognitionLatency() {
	user := s.chain.CreateAccountWithBalance("intent_user", "1000stt")

	intent := "帮我写一个Python快速排序算法"

	start := time.Now()
	result, err := s.chain.SubmitIntent(user, intent)
	elapsed := time.Since(start)

	require.NoError(s.T(), err, "意图提交应成功")
	require.Less(s.T(), elapsed.Seconds(), 3.0,
		"意图识别响应应在3秒内，实际用时: %.2f秒", elapsed.Seconds())
	require.NotNil(s.T(), result, "结果不应为nil")

	s.T().Logf("意图识别响应时间: %.2f秒", elapsed.Seconds())
}

// TestIntentRecognitionAccuracy 测试意图识别准确率
func (s *GenieBotTestSuite) TestIntentRecognitionAccuracy() {
	// 准备测试用例
	testCases := []struct {
		intent   string
		expected string // 期望识别的服务类型
	}{
		{"帮我翻译这段英文", "translation"},
		{"写一段Python代码", "coding"},
		{"总结一下这篇文章", "summarization"},
		{"帮我画一张图", "image_generation"},
		{"分析一下市场数据", "data_analysis"},
		{"写一篇文章", "writing"},
		{"帮我修bug", "debugging"},
		{"生成测试数据", "test_generation"},
		{"解释一下这个算法", "explanation"},
		{"优化这段代码", "optimization"},
	}

	correct := 0
	for _, tc := range testCases {
		result, err := s.chain.RecognizeIntent(tc.intent)
		if err == nil && result.ServiceType == tc.expected {
			correct++
		}
	}

	accuracy := float64(correct) / float64(len(testCases))
	require.GreaterOrEqual(s.T(), accuracy, 0.85,
		"意图识别准确率应>=85%%，实际: %.2f%%", accuracy*100)

	s.T().Logf("意图识别准确率: %.2f%% (%d/%d)", accuracy*100, correct, len(testCases))
}

// TestAIModelInvocation 测试一键调用AI模型
func (s *GenieBotTestSuite) TestAIModelInvocation() {
	user := s.chain.CreateAccountWithBalance("ai_user", "5000stt")

	// 测试调用5种主流AI模型
	models := []string{
		"openai/gpt-4",
		"anthropic/claude-3",
		"openai/gpt-3.5-turbo",
		"stability/stable-diffusion",
		"cohere/command",
	}

	for _, model := range models {
		result, err := s.chain.InvokeAIModel(user, model, "Hello, test prompt")
		require.NoError(s.T(), err, "模型 %s 调用应成功", model)
		require.NotNil(s.T(), result, "结果不应为nil")
	}

	s.T().Logf("成功调用 %d 个AI模型", len(models))
}

// TestCostBreakdown 测试费用明细
func (s *GenieBotTestSuite) TestCostBreakdown() {
	user := s.chain.CreateAccountWithBalance("cost_user", "1000stt")

	// 调用服务并获取费用明细
	result, err := s.chain.InvokeService(user, "openai/gpt-4", "测试消息")
	require.NoError(s.T(), err)

	// 验证费用明细字段
	require.NotZero(s.T(), result.TokenCount, "token数量应不为0")
	require.NotZero(s.T(), result.PricePerToken, "单价应不为0")
	require.NotZero(s.T(), result.TotalCost, "总费用应不为0")
	require.Equal(s.T(), result.TotalCost, float64(result.TokenCount)*result.PricePerToken,
		"总费用 = token数 * 单价")

	s.T().Logf("费用明细 - Tokens: %d, 单价: %f, 总费用: %f",
		result.TokenCount, result.PricePerToken, result.TotalCost)
}

// TestFirstMessageLatency 测试从打开到发送第一条消息的用时 (< 60s)
func (s *GenieBotTestSuite) TestFirstMessageLatency() {
	// 模拟首次使用流程
	start := time.Now()

	// 1. 创建账户 (模拟注册)
	user := s.chain.CreateAccount("first_time_user")

	// 2. 获取测试代币
	err := s.chain.RequestFaucet(user.Address, "1000stt")
	require.NoError(s.T(), err)

	// 3. 发送第一条消息
	_, err = s.chain.SubmitIntent(user, "你好")
	require.NoError(s.T(), err)

	elapsed := time.Since(start)
	require.Less(s.T(), elapsed.Seconds(), 60.0,
		"首次对话应在60秒内完成，实际用时: %.2f秒", elapsed.Seconds())

	s.T().Logf("首次对话总用时: %.2f秒", elapsed.Seconds())
}

func TestGenieBotSuite(t *testing.T) {
	suite.Run(t, new(GenieBotTestSuite))
}
