package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// I18nTest 测试多语言支持
type I18nTest struct {
	E2ETestSuite
}

func TestI18n(t *testing.T) {
	suite.Run(t, new(I18nTest))
}

// TestSupportedLanguages 测试支持的语言
func (s *I18nTest) TestSupportedLanguages() {
	// 查询支持的语言
	languages, err := s.GetSupportedLanguages()
	s.Require().NoError(err)

	// 验证至少支持 3 种语言
	assert.GreaterOrEqual(s.T(), len(languages), 3)

	// 验证包含中文、英文、日文
	codes := make([]string, len(languages))
	for i, lang := range languages {
		codes[i] = lang.Code
	}
	assert.Contains(s.T(), codes, "zh-CN")
	assert.Contains(s.T(), codes, "en")
	assert.Contains(s.T(), codes, "ja")
}

// TestLanguageDetails 测试语言详情
func (s *I18nTest) TestLanguageDetails() {
	// 查询中文详情
	zh, err := s.GetLanguageDetails("zh-CN")
	s.Require().NoError(err)
	assert.Equal(s.T(), "zh-CN", zh.Code)
	assert.Equal(s.T(), "简体中文", zh.Name)
	assert.Equal(s.T(), "Chinese (Simplified)", zh.NameEn)

	// 查询英文详情
	en, err := s.GetLanguageDetails("en")
	s.Require().NoError(err)
	assert.Equal(s.T(), "en", en.Code)
	assert.Equal(s.T(), "English", en.Name)

	// 查询日文详情
	ja, err := s.GetLanguageDetails("ja")
	s.Require().NoError(err)
	assert.Equal(s.T(), "ja", ja.Code)
	assert.Equal(s.T(), "日本語", ja.Name)
}

// TestSetUserLanguage 测试设置用户语言
func (s *I18nTest) TestSetUserLanguage() {
	// 创建用户
	user := s.CreateVerifiedUser("github")

	// 设置语言为中文
	err := s.SetUserLanguage(user.Address, "zh-CN")
	s.Require().NoError(err)

	// 验证设置
	lang, err := s.GetUserLanguage(user.Address)
	s.Require().NoError(err)
	assert.Equal(s.T(), "zh-CN", lang)

	// 切换为英文
	err = s.SetUserLanguage(user.Address, "en")
	s.Require().NoError(err)

	lang, err = s.GetUserLanguage(user.Address)
	s.Require().NoError(err)
	assert.Equal(s.T(), "en", lang)
}

// TestTranslationKeys 测试翻译键
func (s *I18nTest) TestTranslationKeys() {
	// 查询所有翻译键
	keys, err := s.GetTranslationKeys()
	s.Require().NoError(err)

	// 验证包含核心键
	coreKeys := []string{
		"common.welcome",
		"common.confirm",
		"common.cancel",
		"wallet.balance",
		"wallet.send",
		"marketplace.services",
		"task.my_tasks",
	}

	for _, key := range coreKeys {
		assert.Contains(s.T(), keys, key, "应包含翻译键: %s", key)
	}
}

// TestGetTranslation 测试获取翻译
func (s *I18nTest) TestGetTranslation() {
	testCases := []struct {
		key      string
		lang     string
		expected string
	}{
		{"common.welcome", "zh-CN", "欢迎"},
		{"common.welcome", "en", "Welcome"},
		{"common.welcome", "ja", "ようこそ"},
		{"wallet.balance", "zh-CN", "余额"},
		{"wallet.balance", "en", "Balance"},
	}

	for _, tc := range testCases {
		s.Run(tc.lang+"_"+tc.key, func() {
			translation, err := s.GetTranslation(tc.key, tc.lang)
			s.Require().NoError(err)
			assert.Equal(s.T(), tc.expected, translation)
		})
	}
}

// TestFallbackTranslation 测试翻译回退
func (s *I18nTest) TestFallbackTranslation() {
	// 测试不存在的语言回退到英文
	translation, err := s.GetTranslationWithFallback("common.welcome", "fr", "en")
	s.Require().NoError(err)
	assert.Equal(s.T(), "Welcome", translation) // 回退到英文

	// 测试不存在的键
	translation, err = s.GetTranslationWithFallback("nonexistent.key", "zh-CN", "en")
	s.Require().NoError(err)
	assert.Equal(s.T(), "nonexistent.key", translation) // 返回键名
}

// TestLanguageDetection 测试语言检测
func (s *I18nTest) TestLanguageDetection() {
	testCases := []struct {
		text     string
		expected string
	}{
		{"这是一段中文文本", "zh-CN"},
		{"This is English text", "en"},
		{"これは日本語のテキストです", "ja"},
	}

	for _, tc := range testCases {
		s.Run(tc.expected, func() {
			detected, err := s.DetectLanguage(tc.text)
			s.Require().NoError(err)
			assert.Equal(s.T(), tc.expected, detected)
		})
	}
}

// TestTranslationCompletion 测试翻译完成度
func (s *I18nTest) TestTranslationCompletion() {
	// 查询各语言的翻译完成度
	completion, err := s.GetTranslationCompletion()
	s.Require().NoError(err)

	// 英文应 100% 完成
	assert.Equal(s.T(), 100, completion["en"])

	// 中文和日文应达到一定比例
	assert.GreaterOrEqual(s.T(), completion["zh-CN"], 80)
	assert.GreaterOrEqual(s.T(), completion["ja"], 80)
}

// TestMultilingualServiceDescription 测试多语言服务描述
func (s *I18nTest) TestMultilingualServiceDescription() {
	// 创建多语言服务
	service := ServiceWithI18n{
		ID:      "service-1",
		DefaultLang: "en",
		Descriptions: map[string]string{
			"en":    "AI Translation Service",
			"zh-CN": "AI 翻译服务",
			"ja":    "AI翻訳サービス",
		},
	}

	// 存储多语言服务
	err := s.CreateMultilingualService(service)
	s.Require().NoError(err)

	// 查询不同语言版本
	descEN, err := s.GetServiceDescription("service-1", "en")
	s.Require().NoError(err)
	assert.Equal(s.T(), "AI Translation Service", descEN)

	descZH, err := s.GetServiceDescription("service-1", "zh-CN")
	s.Require().NoError(err)
	assert.Equal(s.T(), "AI 翻译服务", descZH)

	descJA, err := s.GetServiceDescription("service-1", "ja")
	s.Require().NoError(err)
	assert.Equal(s.T(), "AI翻訳サービス", descJA)
}

// TestDisputeLanguageSelection 测试争议语言选择
func (s *I18nTest) TestDisputeLanguageSelection() {
	// 创建争议
	user := s.CreateVerifiedUser("github")
	dispute := DisputeWithLanguage{
		Title:    "Test Dispute",
		Language: "zh-CN", // 选择中文
	}

	disputeID, err := s.CreateDisputeWithLanguage(user.Address, dispute)
	s.Require().NoError(err)

	// 验证语言设置
	storedDispute, err := s.GetDispute(disputeID)
	s.Require().NoError(err)
	assert.Equal(s.T(), "zh-CN", storedDispute.Language)
}

// TestAutoTranslation 测试自动翻译
func (s *I18nTest) TestAutoTranslation() {
	// 测试自动翻译服务描述
	serviceID := "service-auto-translate"

	// 创建只有英文描述的服务
	err := s.CreateServiceWithDescription(serviceID, "en", "Advanced AI Assistant")
	s.Require().NoError(err)

	// 请求自动翻译
	err = s.AutoTranslateService(serviceID, []string{"zh-CN", "ja"})
	s.Require().NoError(err)

	// 验证翻译生成
	zhDesc, err := s.GetServiceDescription(serviceID, "zh-CN")
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), zhDesc)
	assert.NotEqual(s.T(), "Advanced AI Assistant", zhDesc)
}

// Helper types

type Language struct {
	Code   string
	Name   string
	NameEn string
	Flag   string
}

type ServiceWithI18n struct {
	ID           string
	DefaultLang  string
	Descriptions map[string]string
}

type DisputeWithLanguage struct {
	Title       string
	Description string
	Language    string
}
