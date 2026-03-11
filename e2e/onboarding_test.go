package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// OnboardingTestSuite 测试 ACH-USER-005: First-Time Onboarding
type OnboardingTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *OnboardingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestWalletAutoCreation 测试钱包自动创建
func (s *OnboardingTestSuite) TestWalletAutoCreation() {
	// 新用户注册时自动创建钱包
	user := s.chain.CreateAccount("onboarding_user")

	// 验证钱包已创建
	require.NotEmpty(s.T(), user.Address, "钱包地址不应为空")
	require.NotEmpty(s.T(), user.Mnemonic, "助记词不应为空")

	// 验证钱包可用（可以查询）
	balance, err := s.chain.QueryBalance(user.Address)
	require.NoError(s.T(), err, "应能查询新创建的钱包")
	require.NotNil(s.T(), balance, "余额对象不应为nil")

	s.T().Logf("新用户钱包自动创建成功: %s", user.Address)
}

// TestFaucetDistribution 测试获得初始测试代币
func (s *OnboardingTestSuite) TestFaucetDistribution() {
	// 创建新用户
	user := s.chain.CreateAccount("faucet_user")

	// 查询初始余额（应为0或最小值）
	initialBalance, _ := s.chain.QueryBalance(user.Address)
	initialAmount := initialBalance.Amount.Int64()

	// 从水龙头获取代币
	err := s.chain.RequestFaucet(user.Address, "1000stt")
	require.NoError(s.T(), err, "水龙头请求应成功")

	// 验证余额增加
	finalBalance, err := s.chain.QueryBalance(user.Address)
	require.NoError(s.T(), err)
	finalAmount := finalBalance.Amount.Int64()

	require.True(s.T(), finalAmount > initialAmount,
		"获取水龙头代币后余额应增加，初始: %d, 最终: %d", initialAmount, finalAmount)

	s.T().Logf("水龙头分发成功，用户获得 %d STT", finalAmount-initialAmount)
}

// TestOnboardingWithin3Minutes 测试3分钟内完成注册体验
func (s *OnboardingTestSuite) TestOnboardingWithin3Minutes() {
	start := time.Now()

	// 1. 自动创建钱包
	user := s.chain.CreateAccount("fast_onboard_user")

	// 2. 从水龙头获取代币
	s.chain.RequestFaucet(user.Address, "1000stt")

	// 3. 验证可以执行操作
	recipient := s.chain.CreateAccount("test_recipient")
	err := s.chain.SendTokens(user, recipient.Address, "100stt")
	require.NoError(s.T(), err, "应能执行转账操作")

	elapsed := time.Since(start)
	require.Less(s.T(), elapsed.Minutes(), 3.0,
		"完整注册流程应在3分钟内完成，实际用时: %.2f分钟", elapsed.Minutes())

	s.T().Logf("注册体验完成，用时: %.2f秒", elapsed.Seconds())
}

// TestTutorialProgress 测试交互式引导教程
func (s *OnboardingTestSuite) TestTutorialProgress() {
	user := s.chain.CreateAccount("tutorial_user")

	// 查询教程进度
	progress, err := s.chain.GetTutorialProgress(user.Address)
	require.NoError(s.T(), err, "应能查询教程进度")

	// 新用户应有默认进度
	require.NotNil(s.T(), progress, "教程进度不应为nil")

	// 更新教程进度（模拟完成步骤）
	err = s.chain.UpdateTutorialProgress(user.Address, 1)
	require.NoError(s.T(), err, "更新教程进度应成功")

	// 验证进度已更新
	progress, _ = s.chain.GetTutorialProgress(user.Address)
	require.Equal(s.T(), 1, progress.CompletedSteps, "应完成1步教程")

	s.T().Log("交互式引导教程功能测试通过")
}

// TestUseCaseList 测试 "What can I do" 常见用例列表
func (s *OnboardingTestSuite) TestUseCaseList() {
	// 查询常见用例
	useCases, err := s.chain.GetCommonUseCases()
	require.NoError(s.T(), err, "应能查询常见用例")
	require.Greater(s.T(), len(useCases), 0, "应至少有一个用例")

	// 验证用例结构
	for _, uc := range useCases {
		require.NotEmpty(s.T(), uc.Title, "用例应有标题")
		require.NotEmpty(s.T(), uc.Description, "用例应有描述")
		require.NotEmpty(s.T(), uc.Action, "用例应有操作指引")
	}

	s.T().Logf("找到 %d 个常见用例", len(useCases))
}

// TestOAuthLoginIntegration 测试微信/GitHub/Google一键登录集成
func (s *OnboardingTestSuite) TestOAuthLoginIntegration() {
	// 模拟OAuth登录流程
	testCases := []string{"wechat", "github", "google"}

	for _, provider := range testCases {
		// 获取OAuth URL
		oauthURL, err := s.chain.GetOAuthURL(provider)
		require.NoError(s.T(), err, "%s OAuth应可获取URL", provider)
		require.Contains(s.T(), oauthURL, "http", "OAuth URL应包含http")

		s.T().Logf("%s OAuth集成可用: %s...", provider, oauthURL[:30])
	}
}

// TestCompleteOnboardingFlow 完整注册流程E2E测试
func (s *OnboardingTestSuite) TestCompleteOnboardingFlow() {
	start := time.Now()

	// 1. 新用户访问，钱包自动创建
	user := s.chain.CreateAccount("complete_flow_user")

	// 2. 获取测试代币
	s.chain.RequestFaucet(user.Address, "1000stt")

	// 3. 查看常见用例
	useCases, _ := s.chain.GetCommonUseCases()
	require.Greater(s.T(), len(useCases), 0)

	// 4. 尝试第一个用例：转账
	recipient := s.chain.CreateAccount("flow_recipient")
	err := s.chain.SendTokens(user, recipient.Address, "100stt")
	require.NoError(s.T(), err)

	elapsed := time.Since(start)
	require.Less(s.T(), elapsed.Seconds(), 180.0, "完整流程应在3分钟内")

	s.T().Logf("完整注册流程E2E测试通过，用时: %.2f秒", elapsed.Seconds())
}

func TestOnboardingSuite(t *testing.T) {
	suite.Run(t, new(OnboardingTestSuite))
}
