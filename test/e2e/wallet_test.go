package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// WalletTestSuite 测试 ACH-USER-001: Secure Digital Wallet
type WalletTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *WalletTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestQueryBalance 测试查看STT余额和交易历史
func (s *WalletTestSuite) TestQueryBalance() {
	// 创建测试账户
	acc := s.chain.CreateAccount("test_balance_user")

	// 验证余额查询API
	balance, err := s.chain.QueryBalance(acc.Address)
	require.NoError(s.T(), err, "余额查询应成功")
	require.NotNil(s.T(), balance, "余额不应为nil")

	// 验证余额格式 (简化检查)
	require.Greater(s.T(), balance.Amount.Int64(), int64(0), "余额应大于0")

	s.T().Logf("账户 %s 余额查询成功", acc.Address)
}

// TestTransactionHistory 测试交易历史查询
func (s *WalletTestSuite) TestTransactionHistory() {
	acc := s.chain.CreateAccount("test_history_user")

	// 发送几笔测试交易
	for i := 0; i < 3; i++ {
		recipient := s.chain.CreateAccount(fmt.Sprintf("recipient_%d", i))
		err := s.chain.SendTokens(acc, recipient.Address, "100stt")
		require.NoError(s.T(), err, "转账应成功")
	}

	// 查询交易历史
	txs, err := s.chain.QueryTxHistory(acc.Address, 10)
	require.NoError(s.T(), err, "交易历史查询应成功")
	require.GreaterOrEqual(s.T(), len(txs), 3, "应至少有3笔交易")

	s.T().Logf("账户 %s 交易历史记录数: %d", acc.Address, len(txs))
}

// TestSecureTransfer 测试安全转账和收款
func (s *WalletTestSuite) TestSecureTransfer() {
	sender := s.chain.CreateAccountWithBalance("test_sender", "10000stt")
	recipient := s.chain.CreateAccount("test_recipient")

	// 查询转账前余额
	senderBalanceBefore, _ := s.chain.QueryBalance(sender.Address)
	recipientBalanceBefore, _ := s.chain.QueryBalance(recipient.Address)

	// 执行转账
	transferAmount := "1000stt"
	err := s.chain.SendTokens(sender, recipient.Address, transferAmount)
	require.NoError(s.T(), err, "转账应成功")

	// 验证余额更新
	senderBalanceAfter, err := s.chain.QueryBalance(sender.Address)
	require.NoError(s.T(), err)
	recipientBalanceAfter, err := s.chain.QueryBalance(recipient.Address)
	require.NoError(s.T(), err)

	// 验证金额变化
	require.Greater(s.T(), senderBalanceBefore.Amount.Int64(), senderBalanceAfter.Amount.Int64(),
		"发送方余额应减少")
	require.Greater(s.T(), recipientBalanceAfter.Amount.Int64(), recipientBalanceBefore.Amount.Int64(),
		"接收方余额应增加")

	s.T().Logf("转账 %s 成功", transferAmount)
}

// TestExportPrivateKey 测试私钥导出功能
func (s *WalletTestSuite) TestExportPrivateKey() {
	acc := s.chain.CreateAccount("test_export_user")

	// 验证私钥导出
	privKey, err := s.chain.ExportPrivateKey(acc.Name)
	require.NoError(s.T(), err, "私钥导出应成功")
	require.NotEmpty(s.T(), privKey, "私钥不应为空")
	require.Greater(s.T(), len(privKey), 20, "私钥长度应合理")

	// 验证私钥格式 (hex encoded)
	require.Regexp(s.T(), "^[0-9a-fA-F]+$", privKey, "私钥应为hex格式")

	s.T().Log("私钥导出功能验证通过")
}

// TestWalletCreationTime 测试钱包创建时间 (< 60s)
func (s *WalletTestSuite) TestWalletCreationTime() {
	start := time.Now()

	_ = s.chain.CreateAccount("test_creation_time")

	elapsed := time.Since(start)
	require.Less(s.T(), elapsed.Seconds(), 60.0,
		"钱包创建应在60秒内完成，实际用时: %.2f秒", elapsed.Seconds())

	s.T().Logf("钱包创建用时: %.2f秒", elapsed.Seconds())
}

// TestWalletFunctionE2E 完整钱包功能E2E测试
func (s *WalletTestSuite) TestWalletFunctionE2E() {
	// 1. 创建钱包
	user := s.chain.CreateAccountWithBalance("e2e_user", "5000stt")

	// 2. 查询余额
	balance, err := s.chain.QueryBalance(user.Address)
	require.NoError(s.T(), err)
	require.Greater(s.T(), balance.Amount.Int64(), int64(0))

	// 3. 发送交易
	recipient := s.chain.CreateAccount("e2e_recipient")
	err = s.chain.SendTokens(user, recipient.Address, "1000stt")
	require.NoError(s.T(), err)

	// 4. 验证交易历史
	txs, err := s.chain.QueryTxHistory(user.Address, 5)
	require.NoError(s.T(), err)
	require.GreaterOrEqual(s.T(), len(txs), 1)

	// 5. 导出私钥
	privKey, err := s.chain.ExportPrivateKey(user.Name)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), privKey)

	s.T().Log("钱包功能E2E测试通过")
}

func TestWalletSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}
