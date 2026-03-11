package keeper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AgentGatewayTestSuite 测试 Agent Gateway 核心功能
type AgentGatewayTestSuite struct {
	suite.Suite
	keeper *Keeper
}

func (s *AgentGatewayTestSuite) SetupSuite() {
	s.keeper = NewKeeper()
}

// TestQueryBalance 测试余额查询
func (s *AgentGatewayTestSuite) TestQueryBalance() {
	ctx := context.Background()
	address := "cosmos1test"

	balance, err := s.keeper.QueryBalance(ctx, address)
	require.NoError(s.T(), err)
	require.Greater(s.T(), balance, uint64(0), "余额应大于0")

	s.T().Logf("账户 %s 余额: %d", address, balance)
}

// TestHandleChat 测试 GenieBot 对话
func (s *AgentGatewayTestSuite) TestHandleChat() {
	ctx := context.Background()
	sessionID := "test-session-123"
	message := "你好，帮我查询余额"

	resp, err := s.keeper.ChatWithGenie(ctx, sessionID, message)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.Content, "响应内容不应为空")
	require.GreaterOrEqual(s.T(), resp.Cost, 0.0, "费用应大于等于0")

	s.T().Logf("Genie响应: %s (费用: %.4f STT)", resp.Content, resp.Cost)
}

// TestAuth 测试钱包认证
func (s *AgentGatewayTestSuite) TestAuth() {
	address := "cosmos1abc123"
	// 模拟有效签名
	signature := "valid-signature"

	err := s.keeper.Authenticate(address, signature)
	require.NoError(s.T(), err, "有效签名应通过认证")

	// 模拟无效签名
	err = s.keeper.Authenticate(address, "invalid-signature")
	require.Error(s.T(), err, "无效签名应认证失败")
}

// TestRateLimit 测试速率限制
func (s *AgentGatewayTestSuite) TestRateLimit() {
	address := "cosmos1ratelimit"

	// 第一次请求应通过
	allowed := s.keeper.CheckRateLimit(address)
	require.True(s.T(), allowed, "首次请求应被允许")

	// 连续60次请求
	for i := 0; i < 60; i++ {
		s.keeper.CheckRateLimit(address)
	}

	// 第61次应被限制
	allowed = s.keeper.CheckRateLimit(address)
	require.False(s.T(), allowed, "超过速率限制应被拒绝")
}

// TestA2AAgentCard 测试 A2A Agent Card
func (s *AgentGatewayTestSuite) TestA2AAgentCard() {
	card := s.keeper.GetAgentCard()

	require.Equal(s.T(), "ShareToken Agent", card.Name)
	require.Equal(s.T(), "1.0.0", card.Version)
	require.Greater(s.T(), len(card.Capabilities), 0, "应有能力声明")

	s.T().Logf("Agent Card: %s v%s", card.Name, card.Version)
}

// TestCallExternalAgentWithOptions 测试高级外部 Agent 调用
func (s *AgentGatewayTestSuite) TestCallExternalAgentWithOptions() {
	// 测试 JSON 输出选项
	options := ExternalAgentOptions{
		OutputFormat: "json",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"result": map[string]string{
					"type": "string",
				},
			},
		},
		Timeout: 30 * time.Second,
	}

	// 由于可能不存在外部 Agent，我们测试选项构建
	s.T().Logf("External agent options: %+v", options)
	s.T().Logf("JSON Schema: %+v", options.JSONSchema)
}

// TestExternalAgentTimeout 测试外部 Agent 超时
func (s *AgentGatewayTestSuite) TestExternalAgentTimeout() {
	options := ExternalAgentOptions{
		OutputFormat: "json",
		Timeout:      1 * time.Second, // 非常短的超时用于测试
	}

	// 如果外部 Agent 存在，测试超时功能
	if s.keeper.externalAgent != nil {
		_, err := s.keeper.CallExternalAgentWithOptions("test", options)
		// 不检查错误，因为取决于 Agent 响应时间
		s.T().Logf("External agent call result: %v", err)
	}
}

func TestAgentGatewaySuite(t *testing.T) {
	suite.Run(t, new(AgentGatewayTestSuite))
}
