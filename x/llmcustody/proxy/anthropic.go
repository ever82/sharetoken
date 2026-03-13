package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicProxy Anthropic API 代理
type AnthropicProxy struct {
	client     *http.Client
	baseURL    string
	maxRetries int
	timeout    time.Duration
	version    string
}

// NewAnthropicProxy 创建新的 Anthropic 代理
func NewAnthropicProxy() *AnthropicProxy {
	return &AnthropicProxy{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL:    "https://api.anthropic.com/v1",
		maxRetries: 3,
		timeout:    60 * time.Second,
		version:    "2023-06-01",
	}
}

// MessagesRequest Claude Messages API 请求
type MessagesRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []ClaudeMessage `json:"messages"`
	System      string          `json:"system,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
	TopK        int             `json:"top_k,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

// ClaudeMessage Claude 消息
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MessagesResponse Claude Messages API 响应
type MessagesResponse struct {
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Role         string           `json:"role"`
	Model        string           `json:"model"`
	Content      []ContentBlock   `json:"content"`
	StopReason   string           `json:"stop_reason"`
	StopSequence string           `json:"stop_sequence,omitempty"`
	Usage        ClaudeTokenUsage `json:"usage"`
}

// ContentBlock 内容块
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ClaudeTokenUsage Token 使用情况
type ClaudeTokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// CompleteRequest Claude Complete API 请求（旧版）
type CompleteRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens_to_sample"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
}

// CompleteResponse Claude Complete API 响应（旧版）
type CompleteResponse struct {
	Completion string           `json:"completion"`
	StopReason string           `json:"stop_reason"`
	Model      string           `json:"model"`
	Usage      ClaudeTokenUsage `json:"usage"`
}

// CallMessages 调用 Claude Messages API（推荐）
func (p *AnthropicProxy) CallMessages(apiKey string, req MessagesRequest) (*MessagesResponse, float64, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", p.baseURL+"/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", p.version)

	// 执行请求（带重试）
	var resp *http.Response
	for i := 0; i < p.maxRetries; i++ {
		resp, err = p.client.Do(httpReq)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		if i < p.maxRetries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to call Anthropic API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("Anthropic API error %d: failed to read error body: %w", resp.StatusCode, err)
		}
		return nil, 0, fmt.Errorf("Anthropic API error %d: %s", resp.StatusCode, string(body))
	}

	var result MessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 计算费用
	cost := p.calculateCost(req.Model, result.Usage)

	return &result, cost, nil
}

// CallComplete 调用 Claude Complete API（旧版，兼容）
func (p *AnthropicProxy) CallComplete(apiKey string, req CompleteRequest) (*CompleteResponse, float64, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", p.baseURL+"/complete", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", p.version)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call Anthropic API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("Anthropic API error %d: failed to read error body: %w", resp.StatusCode, err)
		}
		return nil, 0, fmt.Errorf("Anthropic API error %d: %s", resp.StatusCode, string(body))
	}

	var result CompleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 计算费用
	cost := p.calculateCost(req.Model, result.Usage)

	return &result, cost, nil
}

// ListModels 获取可用模型列表
func (p *AnthropicProxy) ListModels(apiKey string) ([]string, error) {
	// Anthropic 目前没有模型列表 API，返回已知模型
	return []string{
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
		"claude-2.1",
		"claude-2.0",
		"claude-instant-1.2",
	}, nil
}

// calculateCost 计算 API 费用（以 USD 为单位）
func (p *AnthropicProxy) calculateCost(model string, usage ClaudeTokenUsage) float64 {
	// 价格表（每 1K tokens）- 2024 年 3 月价格
	prices := map[string]struct {
		Input  float64
		Output float64
	}{
		"claude-3-opus-20240229":   {Input: 0.015, Output: 0.075},
		"claude-3-sonnet-20240229": {Input: 0.003, Output: 0.015},
		"claude-3-haiku-20240307":  {Input: 0.00025, Output: 0.00125},
		"claude-2.1":               {Input: 0.008, Output: 0.024},
		"claude-2.0":               {Input: 0.008, Output: 0.024},
		"claude-instant-1.2":       {Input: 0.0008, Output: 0.0024},
	}

	price, ok := prices[model]
	if !ok {
		// 默认使用 claude-3-haiku 价格
		price = prices["claude-3-haiku-20240307"]
	}

	inputCost := float64(usage.InputTokens) * price.Input / 1000
	outputCost := float64(usage.OutputTokens) * price.Output / 1000

	return inputCost + outputCost
}

// ValidateAPIKey 验证 API Key 是否有效
func (p *AnthropicProxy) ValidateAPIKey(apiKey string) error {
	// 简单调用验证
	_, err := p.ListModels(apiKey)
	if err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}
	return nil
}

// GetModelInfo 获取模型信息
func (p *AnthropicProxy) GetModelInfo(model string) (map[string]interface{}, error) {
	modelInfo := map[string]interface{}{
		"claude-3-opus-20240229": map[string]interface{}{
			"name":        "Claude 3 Opus",
			"context":     200000,
			"description": "Most powerful model for highly complex tasks",
		},
		"claude-3-sonnet-20240229": map[string]interface{}{
			"name":        "Claude 3 Sonnet",
			"context":     200000,
			"description": "Ideal balance of intelligence and speed",
		},
		"claude-3-haiku-20240307": map[string]interface{}{
			"name":        "Claude 3 Haiku",
			"context":     200000,
			"description": "Fastest model for lightweight actions",
		},
	}

	info, ok := modelInfo[model]
	if !ok {
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	return info.(map[string]interface{}), nil
}
