package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// RefundFlowTestSuite 测试 ACH-USER-009: Service Failure & Refund
type RefundFlowTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *RefundFlowTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestRefundApplication 测试发起退款申请
func (s *RefundFlowTestSuite) TestRefundApplication() {
	user := s.chain.CreateAccountWithBalance("refund_user", "5000stt")
	provider := s.chain.CreateAccount("refund_provider")

	// 创建托管订单
	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-123")

	// 发起退款申请
	refundID, err := s.chain.CreateRefundRequest(user, escrowID, "service_not_delivered")
	require.NoError(s.T(), err, "退款申请应成功")
	require.NotEmpty(s.T(), refundID)

	s.T().Logf("退款申请 %s 创建成功", refundID)
}

// TestRefundReasonSelection 测试选择退款原因
func (s *RefundFlowTestSuite) TestRefundReasonSelection() {
	user := s.chain.CreateAccountWithBalance("reason_user", "5000stt")
	provider := s.chain.CreateAccount("reason_provider")

	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-456")

	// 测试不同退款原因
	reasons := []string{
		"service_not_delivered",
		"quality_not_met",
		"timeout",
	}

	for _, reason := range reasons {
		refundID, err := s.chain.CreateRefundRequest(user, escrowID, reason)
		require.NoError(s.T(), err, "退款原因 '%s' 应被接受", reason)
		require.NotEmpty(s.T(), refundID)
	}

	s.T().Log("退款原因选择测试通过")
}

// TestRefundProgressQuery 测试查看退款处理进度
func (s *RefundFlowTestSuite) TestRefundProgressQuery() {
	user := s.chain.CreateAccountWithBalance("progress_user", "5000stf")
	provider := s.chain.CreateAccount("progress_provider")

	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-789")
	refundID, _ := s.chain.CreateRefundRequest(user, escrowID, "quality_not_met")

	// 查询退款进度
	progress, err := s.chain.QueryRefundProgress(refundID)
	require.NoError(s.T(), err, "应能查询退款进度")
	require.NotEmpty(s.T(), progress.Status)
	require.NotNil(s.T(), progress.CreatedAt)
	require.NotNil(s.T(), progress.EstimatedCompletion)

	s.T().Logf("退款进度: %s", progress.Status)
}

// TestRefundReviewSLA 测试退款审核48小时内完成
func (s *RefundFlowTestSuite) TestRefundReviewSLA() {
	user := s.chain.CreateAccountWithBalance("sla_user", "5000stt")
	provider := s.chain.CreateAccount("sla_provider")

	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-101")
	refundID, _ := s.chain.CreateRefundRequest(user, escrowID, "service_not_delivered")

	// 记录申请时间
	createdAt := time.Now()

	// 模拟审核完成
	err := s.chain.ProcessRefund(refundID, "approved")
	require.NoError(s.T(), err)

	// 查询完成时间
	progress, _ := s.chain.QueryRefundProgress(refundID)
	completedAt := progress.CompletedAt

	// 验证在48小时内
	duration := completedAt.Sub(createdAt)
	require.Less(s.T(), duration.Hours(), 48.0,
		"退款审核应在48小时内完成，实际用时: %.2f小时", duration.Hours())

	s.T().Logf("退款审核用时: %.2f小时", duration.Hours())
}

// TestRefundReturnTime 测试退款10分钟内到账
func (s *RefundFlowTestSuite) TestRefundReturnTime() {
	user := s.chain.CreateAccountWithBalance("return_user", "5000stt")
	provider := s.chain.CreateAccount("return_provider")

	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "1000stt", "service-202")
	refundID, _ := s.chain.CreateRefundRequest(user, escrowID, "service_not_delivered")

	// 记录退款前余额
	balanceBefore, _ := s.chain.QueryBalance(user.Address)

	// 批准退款
	approvedAt := time.Now()
	s.chain.ProcessRefund(refundID, "approved")

	// 验证余额增加
	balanceAfter, _ := s.chain.QueryBalance(user.Address)
	elapsed := time.Since(approvedAt)

	require.True(s.T(), balanceAfter.Amount.Int64() > balanceBefore.Amount.Int64(),
		"退款后余额应增加")
	require.Less(s.T(), elapsed.Minutes(), 10.0,
		"退款应在10分钟内到账，实际用时: %.2f分钟", elapsed.Minutes())

	s.T().Logf("退款到账用时: %.2f分钟", elapsed.Minutes())
}

// TestCompleteRefundFlow 完整退款流程E2E测试
func (s *RefundFlowTestSuite) TestCompleteRefundFlow() {
	// 1. 用户创建订单
	user := s.chain.CreateAccountWithBalance("e2e_refund_user", "10000stt")
	provider := s.chain.CreateAccount("e2e_refund_provider")

	escrowID, _ := s.chain.CreateEscrow(user, provider.Address, "3000stt", "e2e-service")

	// 2. 发起退款
	refundID, err := s.chain.CreateRefundRequest(user, escrowID, "quality_not_met")
	require.NoError(s.T(), err)

	// 3. 查询进度
	progress, _ := s.chain.QueryRefundProgress(refundID)
	require.Equal(s.T(), "pending_review", progress.Status)

	// 4. 处理退款
	err = s.chain.ProcessRefund(refundID, "approved")
	require.NoError(s.T(), err)

	// 5. 验证退款完成
	progress, _ = s.chain.QueryRefundProgress(refundID)
	require.Equal(s.T(), "completed", progress.Status)

	s.T().Log("完整退款流程E2E测试通过")
}

func TestRefundFlowSuite(t *testing.T) {
	suite.Run(t, new(RefundFlowTestSuite))
}
