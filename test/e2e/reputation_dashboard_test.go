package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ReputationDashboardTest 测试信誉仪表板功能
type ReputationDashboardTest struct {
	E2ETestSuite
}

func TestReputationDashboard(t *testing.T) {
	suite.Run(t, new(ReputationDashboardTest))
}

// TestViewMQScore 测试查看 MQ 分数
func (s *ReputationDashboardTest) TestViewMQScore() {
	// 创建测试用户
	user := s.CreateTestUser()

	// 查询 MQ 分数
	mqScore, err := s.QueryMQScore(user.Address)
	s.Require().NoError(err)

	// 验证初始分数为 100
	assert.Equal(s.T(), int64(100), mqScore.Score, "初始 MQ 分数应为 100")
	assert.Equal(s.T(), "normal", mqScore.Level, "初始等级应为 normal")
}

// TestMQScoreHistory 测试 MQ 分数历史
func (s *ReputationDashboardTest) TestMQScoreHistory() {
	// 创建测试用户
	user := s.CreateTestUser()

	// 执行一些操作影响 MQ 分数
	// 模拟一次成功的交易
	s.SimulateSuccessfulTransaction(user.Address)

	// 查询 MQ 分数历史
	history, err := s.QueryMQScoreHistory(user.Address)
	s.Require().NoError(err)

	// 验证历史记录存在
	assert.NotEmpty(s.T(), history, "应有 MQ 分数历史记录")
	assert.GreaterOrEqual(s.T(), len(history), 1, "至少应有 1 条历史记录")
}

// TestReputationLevel 测试信誉等级
func (s *ReputationDashboardTest) TestReputationLevel() {
	testCases := []struct {
		name          string
		mqScore       int64
		expectedLevel string
	}{
		{"新手", 0, "novice"},
		{"普通", 50, "normal"},
		{"良好", 100, "good"},
		{"优秀", 200, "excellent"},
		{"杰出", 500, "outstanding"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// 查询等级
			level, err := s.GetReputationLevel(tc.mqScore)
			s.Require().NoError(err)

			// 验证等级
			assert.Equal(s.T(), tc.expectedLevel, level)
		})
	}
}

// TestReputationBenefits 测试等级权益
func (s *ReputationDashboardTest) TestReputationBenefits() {
	// 创建测试用户
	user := s.CreateTestUser()

	// 查询用户权益
	benefits, err := s.QueryReputationBenefits(user.Address)
	s.Require().NoError(err)

	// 验证权益存在
	assert.NotNil(s.T(), benefits, "应有权益信息")
	assert.GreaterOrEqual(s.T(), benefits.TransactionLimit, int64(0), "交易限额应 >= 0")
	assert.GreaterOrEqual(s.T(), benefits.WithdrawalLimit, int64(0), "提现限额应 >= 0")
}

// TestDisputeParticipationRecord 测试争议参与记录
func (s *ReputationDashboardTest) TestDisputeParticipationRecord() {
	// 创建测试用户
	user := s.CreateTestUser()

	// 查询争议参与记录
	records, err := s.QueryDisputeParticipation(user.Address)
	s.Require().NoError(err)

	// 新用户应无记录
	assert.Empty(s.T(), records, "新用户应无争议记录")
}

// TestJurorSelectionRecord 测试陪审员参与记录
func (s *ReputationDashboardTest) TestJurorSelectionRecord() {
	// 创建高信誉用户
	user := s.CreateTestUserWithMQScore(200) // excellent 等级

	// 查询陪审员资格
	eligibility, err := s.QueryJurorEligibility(user.Address)
	s.Require().NoError(err)

	// 验证有资格
	assert.True(s.T(), eligibility.IsEligible, "MQ 200 应有陪审员资格")

	// 查询陪审员参与记录
	records, err := s.QueryJurorParticipation(user.Address)
	s.Require().NoError(err)

	// 验证记录存在（即使为空）
	assert.NotNil(s.T(), records, "应有陪审员参与记录")
}

// TestContributionStats 测试贡献统计
func (s *ReputationDashboardTest) TestContributionStats() {
	// 创建测试用户
	user := s.CreateTestUser()

	// 查询贡献统计
	stats, err := s.QueryContributionStats(user.Address)
	s.Require().NoError(err)

	// 验证统计字段
	assert.GreaterOrEqual(s.T(), stats.TransactionCount, int64(0), "交易量应 >= 0")
	assert.GreaterOrEqual(s.T(), stats.ReviewCount, int64(0), "评价数应 >= 0")
	assert.GreaterOrEqual(s.T(), stats.ServiceCount, int64(0), "服务次数应 >= 0")
}

// TestMQScoreAfterDispute 测试争议后 MQ 分数变化
func (s *ReputationDashboardTest) TestMQScoreAfterDispute() {
	// 创建测试用户
	user := s.CreateTestUser()

	// 获取初始分数
	initialScore, err := s.QueryMQScore(user.Address)
	s.Require().NoError(err)

	// 模拟参与争议并获得奖励
	s.SimulateDisputeParticipation(user.Address, true)

	// 获取新分数
	newScore, err := s.QueryMQScore(user.Address)
	s.Require().NoError(err)

	// 验证分数增加
	assert.Greater(s.T(), newScore.Score, initialScore.Score, "参与争议获胜后分数应增加")
}
