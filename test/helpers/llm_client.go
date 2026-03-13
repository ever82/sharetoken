// Package helpers provides LLM client for real API integration
package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// LLMClient 提供真实LLM API调用能力
type LLMClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewLLMClient 创建新的LLM客户端
func NewLLMClient() *LLMClient {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_AUTH_TOKEN")
	}

	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	model := os.Getenv("ANTHROPIC_MODEL")
	if model == "" {
		model = "claude-3-sonnet-20240229"
	}

	return &LLMClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsAvailable 检查LLM是否可用
func (c *LLMClient) IsAvailable() bool {
	return c.apiKey != ""
}

// RecognizeIntent 使用LLM进行意图识别
func (c *LLMClient) RecognizeIntent(intent string) (*IntentResult, error) {
	if !c.IsAvailable() {
		// 如果没有API key，使用本地规则
		return c.recognizeIntentLocal(intent), nil
	}

	prompt := fmt.Sprintf(`Analyze the following user intent and classify it into one of these service types:
- coding: programming, code writing, debugging, code review
- translation: language translation, convert text between languages
- summarization: summarize content, create summaries
- image_generation: create images, generate pictures, draw
- data_analysis: analyze data, statistics, charts
- writing: write articles, content creation, creative writing
- debugging: fix bugs, troubleshoot errors
- test_generation: generate test cases, unit tests
- explanation: explain concepts, how things work
- optimization: optimize code, improve performance
- formatting: format data, JSON, XML conversion
- grammar_check: check grammar, spelling
- recommendation: recommend solutions, suggest options
- comparison: compare options, pros and cons
- prediction: predict trends, forecasting
- classification: categorize items, classify data
- extraction: extract information, data extraction
- qa: question answering, Q&A
- list_creation: create lists, enumerate items
- calculation: calculate values, mathematical operations

User intent: "%s"

Respond with ONLY the service type name in lowercase, nothing else.`, intent)

	response, err := c.callLLM(prompt)
	if err != nil {
		return c.recognizeIntentLocal(intent), nil
	}

	serviceType := c.cleanResponse(response)
	confidence := 0.95 // 假设LLM有较高置信度

	return &IntentResult{
		ServiceType: serviceType,
		Confidence:  confidence,
	}, nil
}

// Invoke 调用LLM模型
func (c *LLMClient) Invoke(model string, prompt string) (*AIModelResult, error) {
	if !c.IsAvailable() {
		// 返回模拟结果
		return &AIModelResult{
			Response:      "This is a simulated response (no API key available)",
			TokenCount:    100,
			PricePerToken: 0.001,
			TotalCost:     0.1,
		}, nil
	}

	start := time.Now()
	response, err := c.callLLM(prompt)
	_ = time.Since(start) // 可用于日志记录响应时间

	if err != nil {
		return nil, err
	}

	// 估算token数量 (简化计算)
	tokenCount := int64(len(response) / 4) // 粗略估算: 4字符/token
	pricePerToken := c.getPricePerToken(model)
	totalCost := float64(tokenCount) * pricePerToken

	// 模拟费用计算
	if model == "openai/gpt-4" || model == "claude-3" {
		pricePerToken = 0.00003
		totalCost = float64(tokenCount) * pricePerToken
	}

	return &AIModelResult{
		Response:      response,
		TokenCount:    tokenCount,
		PricePerToken: pricePerToken,
		TotalCost:     totalCost,
	}, nil
}

// callLLM 调用LLM API
func (c *LLMClient) callLLM(prompt string) (string, error) {
	// 构建Anthropic API请求
	reqBody := map[string]interface{
	}{
		"model":      c.model,
		"max_tokens": 1024,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := c.baseURL + "/v1/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	// 解析响应
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}

// recognizeIntentLocal 本地意图识别 (当API不可用时)
func (c *LLMClient) recognizeIntentLocal(intent string) *IntentResult {
	keywords := map[string]string{
		"翻译":      "translation",
		"translate": "translation",
		"代码":      "coding",
		"code":      "coding",
		"写":       "writing",
		"write":     "writing",
		"画":       "image_generation",
		"draw":      "image_generation",
		"图":       "image_generation",
		"image":     "image_generation",
		"分析":      "data_analysis",
		"analyze":   "data_analysis",
		"总结":      "summarization",
		"summarize": "summarization",
		"解释":      "explanation",
		"explain":   "explanation",
		"优化":      "optimization",
		"optimize":  "optimization",
		"修复":      "debugging",
		"fix":       "debugging",
		"bug":       "debugging",
		"测试":      "test_generation",
		"test":      "test_generation",
		"推荐":      "recommendation",
		"recommend": "recommendation",
		"比较":      "comparison",
		"compare":   "comparison",
		"预测":      "prediction",
		"predict":   "prediction",
		"分类":      "classification",
		"classify":  "classification",
		"提取":      "extraction",
		"extract":   "extraction",
		"回答":      "qa",
		"answer":    "qa",
		"计算":      "calculation",
		"calculate": "calculation",
	}

	for keyword, service := range keywords {
		if contains(intent, keyword) {
			return &IntentResult{
				ServiceType: service,
				Confidence:  0.8,
			}
		}
	}

	return &IntentResult{
		ServiceType: "unknown",
		Confidence:  0.5,
	}
}

// cleanResponse 清理LLM响应
func (c *LLMClient) cleanResponse(response string) string {
	// 移除多余空格和换行
	response = trimSpace(response)
	// 转为小写
	response = strings.ToLower(response)
	return response
}

// getPricePerToken 获取模型单价
func (c *LLMClient) getPricePerToken(model string) float64 {
	prices := map[string]float64{
		"openai/gpt-4":             0.00003,
		"gpt-4":                    0.00003,
		"anthropic/claude-3":       0.00003,
		"claude-3":                 0.00003,
		"gpt-3.5-turbo":            0.000002,
		"stability/stable-diffusion": 0.0001,
		"cohere/command":           0.000015,
	}

	if price, ok := prices[model]; ok {
		return price
	}
	return 0.00001 // 默认价格
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// trimSpace 移除字符串首尾空格和换行
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
