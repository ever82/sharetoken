package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sharetoken/x/identity/types"
)

// OpenAIProxy OpenAI API 代理
type OpenAIProxy struct {
	client     *http.Client
	baseURL    string
	maxRetries int
	timeout    time.Duration
}

// NewOpenAIProxy 创建新的 OpenAI 代理
func NewOpenAIProxy() *OpenAIProxy {
	return &OpenAIProxy{
		client: &http.Client{
			Timeout: types.DefaultEscrowDurationHours * time.Hour / 24 * 60, // 60 seconds
		},
		baseURL:    "https://api.openai.com/v1",
		maxRetries: types.MaxAPIRetries,
		timeout:    types.DefaultEscrowDurationHours * time.Hour / 24 * 60, // 60 seconds
	}
}

// ChatCompletionRequest 聊天完成请求
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	N           int           `json:"n,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   TokenUsage             `json:"usage"`
}

// ChatCompletionChoice 完成选项
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// TokenUsage Token 使用情况
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// EmbeddingsRequest 嵌入向量请求
type EmbeddingsRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingsResponse 嵌入向量响应
type EmbeddingsResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  TokenUsage      `json:"usage"`
}

// EmbeddingData 嵌入数据
type EmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

// ModelsResponse 模型列表响应
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// CallChatCompletion 调用聊天完成 API
func (p *OpenAIProxy) CallChatCompletion(apiKey string, req ChatCompletionRequest) (*ChatCompletionResponse, float64, error) {
	// 构建请求
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

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
			time.Sleep(time.Duration(types.InitialBackoffSeconds) * time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("OpenAI API error %d: failed to read error body: %w", resp.StatusCode, err)
		}
		return nil, 0, fmt.Errorf("OpenAI API error %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 计算费用
	cost := p.calculateChatCost(req.Model, result.Usage)

	return &result, cost, nil
}

// CallEmbeddings 调用嵌入向量 API
func (p *OpenAIProxy) CallEmbeddings(apiKey string, req EmbeddingsRequest) (*EmbeddingsResponse, float64, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", p.baseURL+"/embeddings", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("OpenAI API error %d: failed to read error body: %w", resp.StatusCode, err)
		}
		return nil, 0, fmt.Errorf("OpenAI API error %d: %s", resp.StatusCode, string(body))
	}

	var result EmbeddingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 计算费用
	cost := p.calculateEmbeddingsCost(req.Model, result.Usage)

	return &result, cost, nil
}

// ListModels 获取可用模型列表
func (p *OpenAIProxy) ListModels(apiKey string) ([]ModelInfo, error) {
	httpReq, err := http.NewRequest("GET", p.baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("OpenAI API error %d: failed to read error body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("OpenAI API error %d: %s", resp.StatusCode, string(body))
	}

	var result ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// calculateChatCost 计算聊天 API 费用（以 USD 为单位）
func (p *OpenAIProxy) calculateChatCost(model string, usage TokenUsage) float64 {
	// 价格表（每 1K tokens）
	prices := map[string]struct {
		Input  float64
		Output float64
	}{
		"gpt-4":             {Input: 0.03, Output: 0.06},
		"gpt-4-32k":         {Input: 0.06, Output: 0.12},
		"gpt-4-turbo":       {Input: 0.01, Output: 0.03},
		"gpt-3.5-turbo":     {Input: 0.0005, Output: 0.0015},
		"gpt-3.5-turbo-16k": {Input: 0.001, Output: 0.002},
	}

	price, ok := prices[model]
	if !ok {
		// 默认使用 gpt-3.5-turbo 价格
		price = prices["gpt-3.5-turbo"]
	}

	inputCost := float64(usage.PromptTokens) * price.Input / 1000
	outputCost := float64(usage.CompletionTokens) * price.Output / 1000

	return inputCost + outputCost
}

// calculateEmbeddingsCost 计算 Embeddings API 费用
func (p *OpenAIProxy) calculateEmbeddingsCost(model string, usage TokenUsage) float64 {
	// 价格表（每 1K tokens）
	prices := map[string]float64{
		"text-embedding-ada-002": 0.0001,
		"text-embedding-3-small": 0.00002,
		"text-embedding-3-large": 0.00013,
	}

	price, ok := prices[model]
	if !ok {
		price = prices["text-embedding-ada-002"]
	}

	return float64(usage.TotalTokens) * price / 1000
}

// ValidateAPIKey 验证 API Key 是否有效
func (p *OpenAIProxy) ValidateAPIKey(apiKey string) error {
	models, err := p.ListModels(apiKey)
	if err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}
	if len(models) == 0 {
		return fmt.Errorf("API key valid but no models available")
	}
	return nil
}
