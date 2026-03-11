package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// EscrowSecurityTestSuite 测试 ACH-USER-003: Fund Security Guarantee
type EscrowSecurityTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *EscrowSecurityTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestFundEscrowLocking 测试服务开始前资金进入托管
func (s *EscrowSecurityTestSuite) TestFundEscrowLocking() {
	// 创建用户和服务提供者
	user := s.chain.CreateAccountWithBalance("escrow_user", "5000stt")
	provider := s.chain.CreateAccountWithBalance("escrow_provider", "1000stt")

	// 创建托管订单
	escrowID, err := s.chain.CreateEscrow(
		user,
		provider.Address,
		"1000stt",
		"service-123",
	)
	require.NoError(s.T(), err, "创建托管订单应成功")
	require.NotEmpty(s.T(), escrowID, "托管ID不应为空")

	// 验证资金状态
	status, err := s.chain.QueryEscrowStatus(escrowID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "locked", status.State, "资金应处于锁定状态")

	// 验证用户余额减少
	userBalance, _ := s.chain.QueryBalance(user.Address)
	require.True(s.T(), userBalance.Amount.Int64() < 5000000,
		"用户余额应减少")

	s.T().Logf("托管订单 %s 创建成功，资金已锁定", escrowID)
}

// TestFundReleaseOnCompletion 测试确认满意后资金释放
func (s *EscrowSecurityTestSuite) TestFundReleaseOnCompletion() {
	user := s.chain.CreateAccountWithBalance("release_user", "5000stt")
	provider := s.chain.CreateAccountWithBalance("release_provider", "1000stt")

	// 创建托管订单
	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-456")

	// 记录释放前余额
	providerBalanceBefore, _ := s.chain.QueryBalance(provider.Address)

	// 用户确认完成，释放资金
	err := s.chain.ReleaseEscrow(user, escrowID)
	require.NoError(s.T(), err, "释放资金应成功")

	// 验证资金状态
	status, err := s.chain.QueryEscrowStatus(escrowID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "completed", status.State, "资金应处于已完成状态")

	// 验证提供者余额增加
	providerBalanceAfter, _ := s.chain.QueryBalance(provider.Address)
	require.True(s.T(), providerBalanceAfter.Amount.Int64() > providerBalanceBefore.Amount.Int64(),
		"提供者余额应增加")

	s.T().Log("资金释放成功，提供者已收到款项")
}

// TestFundFreezing 测试争议发起时冻结资金
func (s *EscrowSecurityTestSuite) TestFundFreezing() {
	user := s.chain.CreateAccountWithBalance("freeze_user", "5000stt")
	provider := s.chain.CreateAccountWithBalance("freeze_provider", "1000stt")

	// 创建托管订单
	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-789")

	// 用户发起争议，冻结资金
	err := s.chain.FreezeEscrow(user, escrowID, "服务质量不符合预期")
	require.NoError(s.T(), err, "冻结资金应成功")

	// 验证资金状态
	status, err := s.chain.QueryEscrowStatus(escrowID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "disputed", status.State, "资金应处于争议状态")

	s.T().Log("资金冻结成功，进入争议处理流程")
}

// TestFundDistributionOnDisputeResolution 测试争议解决后资金分配
func (s *EscrowSecurityTestSuite) TestFundDistributionOnDisputeResolution() {
	user := s.chain.CreateAccountWithBalance("dist_user", "5000stt")
	provider := s.chain.CreateAccountWithBalance("dist_provider", "1000stt")

	// 创建托管订单并进入争议状态
	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-101")
	s.chain.FreezeEscrow(user, escrowID, "争议原因")

	// 记录裁决前余额
	userBalanceBefore, _ := s.chain.QueryBalance(user.Address)
	providerBalanceBefore, _ := s.chain.QueryBalance(provider.Address)

	// 模拟争议裁决：70%给用户，30%给提供者
	err := s.chain.ResolveEscrowDispute(
		s.chain.GetValidator(),
		escrowID,
		70, // 用户获得70%
		30, // 提供者获得30%
	)
	require.NoError(s.T(), err, "争议裁决应成功")

	// 验证资金分配
	userBalanceAfter, _ := s.chain.QueryBalance(user.Address)
	providerBalanceAfter, _ := s.chain.QueryBalance(provider.Address)

	require.True(s.T(), userBalanceAfter.Amount.Int64() > userBalanceBefore.Amount.Int64(),
		"用户应获得退款")
	require.True(s.T(), providerBalanceAfter.Amount.Int64() > providerBalanceBefore.Amount.Int64(),
		"提供者应获得部分款项")

	s.T().Log("争议解决，资金按比例分配成功")
}

// TestEscrowStatusVisibility 测试随时查看托管资金状态
func (s *EscrowSecurityTestSuite) TestEscrowStatusVisibility() {
	user := s.chain.CreateAccountWithBalance("status_user", "5000stt")
	provider := s.chain.CreateAccountWithBalance("status_provider", "1000stt")

	// 创建多个托管订单
	for i := 0; i < 3; i++ {
		s.chain.CreateEscrow(user, provider.Address, "500stt", fmt.Sprintf("service-%d", i))
	}

	// 查询用户的所有托管订单
	escrows, err := s.chain.QueryUserEscrows(user.Address)
	require.NoError(s.T(), err, "查询托管订单应成功")
	require.GreaterOrEqual(s.T(), len(escrows), 3, "应至少有3个托管订单")

	// 验证每个订单都有完整的状态信息
	for _, escrow := range escrows {
		require.NotEmpty(s.T(), escrow.ID, "托管ID不应为空")
		require.NotEmpty(s.T(), escrow.State, "状态不应为空")
		require.NotEmpty(s.T(), escrow.Amount, "金额不应为空")
		require.NotEmpty(s.T(), escrow.ServiceID, "服务ID不应为空")
	}

	s.T().Logf("用户有 %d 个托管订单", len(escrows))
}

// TestEscrowE2EFlow 完整托管流程E2E测试
func (s *EscrowSecurityTestSuite) TestEscrowE2EFlow() {
	// 完整流程：创建托管 -> 服务执行 -> 确认完成 -> 资金释放
	user := s.chain.CreateAccountWithBalance("e2e_user", "10000stt")
	provider := s.chain.CreateAccountWithBalance("e2e_provider", "1000stt")

	// 1. 用户创建托管订单
	escrowID, err := s.chain.CreateEscrow(user, provider.Address, "3000stt", "e2e-service")
	require.NoError(s.T(), err)

	// 2. 验证资金已锁定
	status, _ := s.chain.QueryEscrowStatus(escrowID)
	require.Equal(s.T(), "locked", status.State)

	// 3. 模拟服务完成
	providerBalanceBefore, _ := s.chain.QueryBalance(provider.Address)

	// 4. 用户确认完成
	err = s.chain.ReleaseEscrow(user, escrowID)
	require.NoError(s.T(), err)

	// 5. 验证资金已释放
	providerBalanceAfter, _ := s.chain.QueryBalance(provider.Address)
	require.True(s.T(), providerBalanceAfter.Amount.Int64() > providerBalanceBefore.Amount.Int64(),
		"提供者余额应增加")

	s.T().Log("完整托管流程E2E测试通过")
}

func TestEscrowSecuritySuite(t *testing.T) {
	suite.Run(t, new(EscrowSecurityTestSuite))
}
