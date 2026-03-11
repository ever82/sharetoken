package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// LLMClient LLM客户端接口
type LLMClient interface {
	Invoke(model, prompt string) (*LLMResult, error)
}

// LLMResult LLM调用结果
type LLMResult struct {
	Response      string
	TokenCount    int64
	PricePerToken float64
	TotalCost     float64
}

// ExternalAgentOptions 外部 Agent 调用选项
type ExternalAgentOptions struct {
	OutputFormat string                 // json, text
	JSONSchema   map[string]interface{} // JSON Schema 用于结构化输出
	Timeout      time.Duration          // 超时时间
}

// callLLM 调用真实 LLM (Claude)
func (k *Keeper) callLLM(messages []Message) (string, float64, error) {
	// 1. 优先尝试使用本地外部 Agent (Claude Code / OpenClaw)
	if k.externalAgent != nil {
		response, err := k.callExternalAgent(messages)
		if err == nil {
			// 使用外部 Agent 不消耗 STT
			return response, 0, nil
		}
		// 外部 Agent 调用失败，继续尝试其他方式
	}

	// 2. 检查 API Key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return k.buildMockResponse(), 0.001, nil
	}

	// 3. 检测 API Key 格式
	if strings.HasPrefix(apiKey, "sk-sp-") {
		return k.handleClaudeCodeToken(), 0.001, nil
	}

	// 4. 调用 Claude API
	return k.callClaudeAPI(messages, apiKey)
}

// buildMockResponse 构建模拟响应
func (k *Keeper) buildMockResponse() string {
	if k.externalAgent != nil {
		return fmt.Sprintf("[模拟模式] 外部 Agent (%s) 调用失败，请检查配置。", k.externalAgent.Command)
	}
	return "[模拟模式] 请在设置中配置 ANTHROPIC_API_KEY 以获得真实 AI 回复，或安装 Claude Code / OpenClaw。"
}

// handleClaudeCodeToken 处理 Claude Code session token
func (k *Keeper) handleClaudeCodeToken() string {
	if k.externalAgent != nil {
		return fmt.Sprintf("[模拟模式] 请直接使用本地 Claude Code (已检测到: %s)", k.externalAgent.Command)
	}
	return "[模拟模式] 检测到 Claude Code session token，请使用 Anthropic API Key (sk-ant-xxx)，或安装 Claude Code。"
}

// buildPrompt 构建对话提示
func (k *Keeper) buildPrompt(messages []Message) string {
	var conversation string
	for _, msg := range messages {
		role := msg.Role
		if role == "user" {
			conversation += fmt.Sprintf("\nHuman: %s", msg.Content)
		} else {
			conversation += fmt.Sprintf("\nAssistant: %s", msg.Content)
		}
	}
	return conversation
}

// createLLMRequest 创建 LLM API 请求
func (k *Keeper) createLLMRequest(messages []Message, apiKey string) (*http.Request, error) {
	conversation := k.buildPrompt(messages)

	reqBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307",
		"max_tokens": 1024,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": conversation,
			},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	return req, nil
}

// parseLLMResponse 解析 LLM API 响应
func (k *Keeper) parseLLMResponse(resp *http.Response) (string, float64, error) {
	if resp.StatusCode == 403 {
		return "", 0, errors.New("LLM API 403: 请检查 API Key 是否有效，或在 https://console.anthropic.com/ 生成新的 API Key")
	}
	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("LLM API error: %d", resp.StatusCode)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	if len(result.Content) == 0 {
		return "", 0, errors.New("empty response from LLM")
	}

	// 计算费用 (Claude Haiku: $0.25/1M input, $1.25/1M output)
	inputCost := float64(result.Usage.InputTokens) * 0.25 / 1000000
	outputCost := float64(result.Usage.OutputTokens) * 1.25 / 1000000
	totalCost := inputCost + outputCost

	return result.Content[0].Text, totalCost, nil
}

// callClaudeAPI 调用 Claude API
func (k *Keeper) callClaudeAPI(messages []Message, apiKey string) (string, float64, error) {
	req, err := k.createLLMRequest(messages, apiKey)
	if err != nil {
		return "", 0, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	return k.parseLLMResponse(resp)
}

// callExternalAgent 通过 STDIO 调用外部 Agent
func (k *Keeper) callExternalAgent(messages []Message) (string, error) {
	if k.externalAgent == nil {
		return "", errors.New("no external agent configured")
	}

	// 构建对话内容 - 只取最后一条用户消息
	lastMessage := k.getLastUserMessage(messages)
	if lastMessage == "" {
		return "", errors.New("no user message found")
	}

	// 执行外部 Agent
	return k.executeExternalAgent(lastMessage, k.externalAgent.Args)
}

// getLastUserMessage 获取最后一条用户消息
func (k *Keeper) getLastUserMessage(messages []Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return messages[i].Content
		}
	}
	return ""
}

// executeExternalAgent 执行外部 Agent 命令
func (k *Keeper) executeExternalAgent(message string, extraArgs []string) (string, error) {
	args := append(extraArgs, message)

	cmd := exec.Command(k.externalAgent.Command, args...)

	// 设置环境变量 - 清除 CLAUDECODE 以避免嵌套检测
	cmd.Env = k.buildAgentEnv()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("external agent failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// buildAgentEnv 构建 Agent 环境变量
func (k *Keeper) buildAgentEnv() []string {
	env := os.Environ()
	filteredEnv := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, "CLAUDECODE=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	filteredEnv = append(filteredEnv, "GENIEBOT_AGENT_CALL=true")
	return filteredEnv
}

// CallExternalAgentWithOptions 使用高级选项调用外部 Agent
func (k *Keeper) CallExternalAgentWithOptions(prompt string, opts ExternalAgentOptions) (string, error) {
	if k.externalAgent == nil {
		return "", errors.New("no external agent configured")
	}

	args := k.buildAgentArgs(prompt, opts)

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, k.externalAgent.Command, args...)
	cmd.Env = k.buildAgentEnv()

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("external agent timeout after %v", timeout)
		}
		return "", fmt.Errorf("external agent failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// buildAgentArgs 构建 Agent 参数
func (k *Keeper) buildAgentArgs(prompt string, opts ExternalAgentOptions) []string {
	args := []string{"-p", prompt}

	if opts.OutputFormat != "" {
		args = append(args, "--output-format", opts.OutputFormat)
	}

	if opts.JSONSchema != nil && opts.OutputFormat == "json" {
		schemaJSON, err := json.Marshal(opts.JSONSchema)
		if err == nil {
			args = append(args, "--json-schema", string(schemaJSON))
		}
	}

	return args
}
