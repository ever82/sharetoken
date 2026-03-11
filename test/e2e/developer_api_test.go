package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DeveloperAPITest 测试开发者 API 功能
type DeveloperAPITest struct {
	E2ETestSuite
}

func TestDeveloperAPI(t *testing.T) {
	suite.Run(t, new(DeveloperAPITest))
}

// TestGenerateAPIKey 测试 API Key 生成
func (s *DeveloperAPITest) TestGenerateAPIKey() {
	// 创建开发者用户
	developer := s.CreateVerifiedUser("github")

	// 生成 API Key
	apiKey, err := s.GenerateAPIKey(developer.Address, "Test Key")
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), apiKey.ID)
	assert.NotEmpty(s.T(), apiKey.Key)

	// 验证 API Key 格式
	assert.True(s.T(), len(apiKey.Key) >= 32, "API Key 应至少 32 位")
	assert.NotEmpty(s.T(), apiKey.Prefix)
}

// TestAPIKeyPermissions 测试 API Key 权限
func (s *DeveloperAPITest) TestAPIKeyPermissions() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成带权限的 API Key
	permissions := APIKeyPermissions{
		ReadBalance:  true,
		SendTx:       true,
		CallServices: true,
		Admin:        false,
	}

	apiKey, err := s.GenerateAPIKeyWithPermissions(developer.Address, "Scoped Key", permissions)
	s.Require().NoError(err)

	// 验证权限
	storedPermissions, err := s.GetAPIKeyPermissions(apiKey.ID)
	s.Require().NoError(err)
	assert.True(s.T(), storedPermissions.ReadBalance)
	assert.True(s.T(), storedPermissions.SendTx)
	assert.False(s.T(), storedPermissions.Admin)
}

// TestAPIKeyRateLimit 测试 API Key 速率限制
func (s *DeveloperAPITest) TestAPIKeyRateLimit() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成带限额的 API Key
	limits := APIKeyLimits{
		RequestsPerMinute: 60,
		RequestsPerHour:   1000,
		RequestsPerDay:    10000,
		MaxCostPerDay:     10000000, // 10 STT
	}

	apiKey, err := s.GenerateAPIKeyWithLimits(developer.Address, "Limited Key", limits)
	s.Require().NoError(err)

	// 验证限额
	storedLimits, err := s.GetAPIKeyLimits(apiKey.ID)
	s.Require().NoError(err)
	assert.Equal(s.T(), limits.RequestsPerMinute, storedLimits.RequestsPerMinute)
	assert.Equal(s.T(), limits.MaxCostPerDay, storedLimits.MaxCostPerDay)
}

// TestAPIKeyWhitelist 测试 API Key IP 白名单
func (s *DeveloperAPITest) TestAPIKeyWhitelist() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成带白名单的 API Key
	whitelist := []string{
		"192.168.1.0/24",
		"10.0.0.0/8",
	}

	apiKey, err := s.GenerateAPIKeyWithWhitelist(developer.Address, "Whitelisted Key", whitelist)
	s.Require().NoError(err)

	// 验证白名单
	storedWhitelist, err := s.GetAPIKeyWhitelist(apiKey.ID)
	s.Require().NoError(err)
	assert.Equal(s.T(), whitelist, storedWhitelist)
}

// TestAPIKeyRevocation 测试 API Key 撤销
func (s *DeveloperAPITest) TestAPIKeyRevocation() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成 API Key
	apiKey, _ := s.GenerateAPIKey(developer.Address, "Key to Revoke")

	// 验证 Key 有效
	isValid, err := s.ValidateAPIKey(apiKey.Key)
	s.Require().NoError(err)
	assert.True(s.T(), isValid)

	// 撤销 Key
	err = s.RevokeAPIKey(developer.Address, apiKey.ID)
	s.Require().NoError(err)

	// 验证 Key 已失效
	isValid, err = s.ValidateAPIKey(apiKey.Key)
	s.Require().NoError(err)
	assert.False(s.T(), isValid)
}

// TestAPICallLogging 测试 API 调用日志
func (s *DeveloperAPITest) TestAPICallLogging() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成 API Key
	apiKey, _ := s.GenerateAPIKey(developer.Address, "Logged Key")

	// 模拟 API 调用
	call := APICall{
		Method:   "GET",
		Endpoint: "/v1/balance",
		Status:   200,
		Cost:     0,
	}

	err := s.LogAPICall(apiKey.ID, call)
	s.Require().NoError(err)

	// 查询日志
	logs, err := s.GetAPIKeyLogs(apiKey.ID, 10)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(logs), 1)
	assert.Equal(s.T(), "GET", logs[0].Method)
	assert.Equal(s.T(), "/v1/balance", logs[0].Endpoint)
}

// TestAPIKeyUsageStats 测试 API Key 使用统计
func (s *DeveloperAPITest) TestAPIKeyUsageStats() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成 API Key
	apiKey, _ := s.GenerateAPIKey(developer.Address, "Stats Key")

	// 模拟多个 API 调用
	for i := 0; i < 5; i++ {
		s.LogAPICall(apiKey.ID, APICall{
			Method:   "POST",
			Endpoint: "/v1/llm/call",
			Status:   200,
			Cost:     100000, // 0.1 STT
		})
	}

	// 查询使用统计
	stats, err := s.GetAPIKeyUsageStats(apiKey.ID, "24h")
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), stats.TotalCalls, int64(5))
	assert.GreaterOrEqual(s.T(), stats.TotalCost, int64(500000))
	assert.Equal(s.T(), int64(0), stats.ErrorCount) // 全部成功
}

// TestAPICallWithInvalidKey 测试无效 Key 调用
func (s *DeveloperAPITest) TestAPICallWithInvalidKey() {
	// 尝试使用无效 Key 调用
	isValid, err := s.ValidateAPIKey("invalid-key-12345")
	s.Require().NoError(err)
	assert.False(s.T(), isValid)
}

// TestRateLimitExceeded 测试超出速率限制
func (s *DeveloperAPITest) TestRateLimitExceeded() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成严格限流的 API Key
	limits := APIKeyLimits{
		RequestsPerMinute: 1, // 每分钟 1 次
	}

	apiKey, _ := s.GenerateAPIKeyWithLimits(developer.Address, "Strict Key", limits)

	// 第一次调用应成功
	canCall, err := s.CheckRateLimit(apiKey.ID)
	s.Require().NoError(err)
	assert.True(s.T(), canCall)

	// 立即第二次调用应被拒绝
	canCall, err = s.CheckRateLimit(apiKey.ID)
	s.Require().NoError(err)
	assert.False(s.T(), canCall)
}

// TestDailyCostLimit 测试每日费用限制
func (s *DeveloperAPITest) TestDailyCostLimit() {
	// 创建开发者
	developer := s.CreateVerifiedUser("github")

	// 生成带每日费用限制的 API Key
	limits := APIKeyLimits{
		MaxCostPerDay: 1000000, // 1 STT
	}

	apiKey, _ := s.GenerateAPIKeyWithLimits(developer.Address, "Cost Limited Key", limits)

	// 模拟消耗费用
	err := s.LogAPICall(apiKey.ID, APICall{
		Method:   "POST",
		Endpoint: "/v1/llm/call",
		Status:   200,
		Cost:     900000, // 0.9 STT
	})
	s.Require().NoError(err)

	// 检查剩余额度
	remaining, err := s.GetAPIKeyDailyCostRemaining(apiKey.ID)
	s.Require().NoError(err)
	assert.Less(s.T(), remaining, int64(200000)) // 剩余不到 0.2 STT
}

// Helper types

type APIKey struct {
	ID        string
	Key       string
	Prefix    string
	Name      string
	CreatedAt int64
}

type APIKeyPermissions struct {
	ReadBalance  bool
	SendTx       bool
	CallServices bool
	Admin        bool
}

type APIKeyLimits struct {
	RequestsPerMinute int64
	RequestsPerHour   int64
	RequestsPerDay    int64
	MaxCostPerDay     int64
}

type APICall struct {
	Method    string
	Endpoint  string
	Status    int
	Cost      int64
	Error     string
	Timestamp int64
}

type APIKeyUsageStats struct {
	TotalCalls  int64
	TotalCost   int64
	ErrorCount  int64
	AvgLatency  float64
}
