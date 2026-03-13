package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// IdentityLimitsTest 测试身份认证相关的限额功能
type IdentityLimitsTest struct {
	E2ETestSuite
}

func TestIdentityLimits(t *testing.T) {
	suite.Run(t, new(IdentityLimitsTest))
}

// TestUnverifiedUserLimits 测试未认证用户限额
func (s *IdentityLimitsTest) TestUnverifiedUserLimits() {
	// 创建未认证用户
	user := s.CreateUnverifiedUser()

	// 查询限额
	limits, err := s.QueryUserLimits(user.Address)
	s.Require().NoError(err)

	// 验证未认证用户限额较低
	assert.Equal(s.T(), int64(1000000000), limits.TransactionLimit, "未认证用户交易限额应为 1000 STT")
	assert.Equal(s.T(), int64(500000000), limits.WithdrawalLimit, "未认证用户提现限额应为 500 STT")
	assert.Equal(s.T(), int64(100000000), limits.DisputeLimit, "未认证用户争议限额应为 100 STT")
}

// TestVerifiedUserLimits 测试已认证用户限额
func (s *IdentityLimitsTest) TestVerifiedUserLimits() {
	// 创建已认证用户
	user := s.CreateVerifiedUser("github")

	// 查询限额
	limits, err := s.QueryUserLimits(user.Address)
	s.Require().NoError(err)

	// 验证已认证用户限额更高
	assert.Equal(s.T(), int64(10000000000), limits.TransactionLimit, "认证用户交易限额应为 10000 STT")
	assert.Equal(s.T(), int64(5000000000), limits.WithdrawalLimit, "认证用户提现限额应为 5000 STT")
	assert.Equal(s.T(), int64(1000000000), limits.DisputeLimit, "认证用户争议限额应为 1000 STT")
}

// TestMultipleVerificationMethods 测试多种认证方式
func (s *IdentityLimitsTest) TestMultipleVerificationMethods() {
	methods := []string{"wechat", "github", "google"}

	for _, method := range methods {
		s.Run(method, func() {
			// 使用不同方式认证
			user := s.CreateVerifiedUser(method)

			// 查询认证状态
			status, err := s.QueryIdentityStatus(user.Address)
			s.Require().NoError(err)

			// 验证已认证
			assert.True(s.T(), status.IsVerified, "%s 认证应成功", method)
			assert.Equal(s.T(), method, status.Provider, "认证方式应为 %s", method)
		})
	}
}

// TestLimitEnforcement 测试限额执行
func (s *IdentityLimitsTest) TestLimitEnforcement() {
	// 创建未认证用户
	user := s.CreateUnverifiedUser()

	// 给用户提供足够的资金
	s.FundAccount(user.Address, 2000000000) // 2000 STT

	// 尝试超过交易限额
	_, err := s.CreateEscrow(user.Address, 1500000000) // 1500 STT，超过 1000 限额

	// 应被拒绝
	assert.Error(s.T(), err, "超过限额的交易应被拒绝")
}

// TestLimitUpgradeAfterVerification 测试认证后限额提升
func (s *IdentityLimitsTest) TestLimitUpgradeAfterVerification() {
	// 创建未认证用户
	user := s.CreateUnverifiedUser()

	// 获取认证前限额
	limitsBefore, err := s.QueryUserLimits(user.Address)
	s.Require().NoError(err)

	// 执行认证
	s.VerifyIdentity(user.Address, "github")

	// 获取认证后限额
	limitsAfter, err := s.QueryUserLimits(user.Address)
	s.Require().NoError(err)

	// 验证限额提升
	assert.Greater(s.T(), limitsAfter.TransactionLimit, limitsBefore.TransactionLimit, "认证后交易限额应提升")
	assert.Greater(s.T(), limitsAfter.WithdrawalLimit, limitsBefore.WithdrawalLimit, "认证后提现限额应提升")
}

// TestJurorEligibilityByLevel 测试等级与陪审员资格
func (s *IdentityLimitsTest) TestJurorEligibilityByLevel() {
	testCases := []struct {
		name             string
		mqScore          int64
		isEligible       bool
		minRequiredScore int64
	}{
		{"新手不可", 50, false, 100},
		{"普通刚好", 100, true, 100},
		{"良好可以", 150, true, 100},
		{"优秀可以", 200, true, 100},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// 创建特定 MQ 分数的用户
			user := s.CreateTestUserWithMQScore(tc.mqScore)

			// 查询陪审员资格
			eligibility, err := s.QueryJurorEligibility(user.Address)
			s.Require().NoError(err)

			// 验证资格
			assert.Equal(s.T(), tc.isEligible, eligibility.IsEligible)
			assert.Equal(s.T(), tc.minRequiredScore, eligibility.MinRequiredScore)
		})
	}
}
