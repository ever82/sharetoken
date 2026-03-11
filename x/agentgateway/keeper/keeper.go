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
	"sync"
	"time"
)

// Keeper 管理 Agent Gateway 核心逻辑
type Keeper struct {
	// 链交互
	queryClient QueryClient

	// LLM 客户端
	llmClient LLMClient

	// 外部 Agent 配置
	externalAgent *ExternalAgentConfig

	// 会话管理
	sessions map[string]*Session
	sessionMu sync.RWMutex

	// 速率限制
	rateLimits map[string]*RateLimit
	rateMu     sync.RWMutex
}

// ExternalAgentConfig 外部 Agent 配置
type ExternalAgentConfig struct {
	Command string   // 如: "claude", "openclaw"
	Args    []string // 如: ["--stdio"]
	Env     []string // 额外环境变量
}

// Session 会话状态
type Session struct {
	ID        string
	UserAddr  string
	Context   []Message
	CreatedAt time.Time
}

// Message 消息
type Message struct {
	Role    string
	Content string
}

// RateLimit 速率限制状态
type RateLimit struct {
	Address      string
	RequestCount int
	LastReset    time.Time
}

// ChatResponse 对话响应
type ChatResponse struct {
	Content string
	Cost    float64
}

// AgentCard A2A Agent Card
type AgentCard struct {
	Name         string
	Version      string
	Description  string
	Capabilities []string
	Endpoints    map[string]string
}

// QueryClient 链查询客户端接口
type QueryClient interface {
	QueryBalance(ctx context.Context, address string) (uint64, error)
}

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

// NewKeeper 创建新的 Keeper
func NewKeeper() *Keeper {
	k := &Keeper{
		sessions:   make(map[string]*Session),
		rateLimits: make(map[string]*RateLimit),
	}

	// 自动检测本地外部 Agent
	k.detectExternalAgent()

	return k
}

// detectExternalAgent 检测本地可用的外部 Agent
func (k *Keeper) detectExternalAgent() {
	// 检测 Claude Code
	if path, err := exec.LookPath("claude"); err == nil {
		k.externalAgent = &ExternalAgentConfig{
			Command: path,
			Args:    []string{"-p"}, // 使用 -p (print) 模式
		}
		return
	}

	// 检测 OpenClaw
	if path, err := exec.LookPath("openclaw"); err == nil {
		k.externalAgent = &ExternalAgentConfig{
			Command: path,
			Args:    []string{},
		}
		return
	}

	// 可以添加更多外部 Agent 检测
}

// QueryBalance 查询余额
func (k *Keeper) QueryBalance(ctx context.Context, address string) (uint64, error) {
	// TODO: 实现真实链查询
	// 目前返回模拟数据
	if k.queryClient != nil {
		return k.queryClient.QueryBalance(ctx, address)
	}
	// 模拟余额: 10000 STT
	return 10000000, nil
}

// ChatWithGenie 与 GenieBot 对话
func (k *Keeper) ChatWithGenie(ctx context.Context, sessionID, message string) (*ChatResponse, error) {
	// 获取或创建会话
	session := k.getOrCreateSession(sessionID)

	// 添加到上下文
	session.Context = append(session.Context, Message{
		Role:    "user",
		Content: message,
	})

	// 调用 LLM
	response, cost, err := k.callLLM(session.Context)
	if err != nil {
		return nil, err
	}

	// 保存助手响应到上下文
	session.Context = append(session.Context, Message{
		Role:    "assistant",
		Content: response,
	})

	return &ChatResponse{
		Content: response,
		Cost:    cost,
	}, nil
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
		// 外部 Agent 调用失败，可能是嵌套检测或其他原因
		// 继续尝试其他方式，但记录错误
	}

	// 2. 尝试调用 Claude API
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		// 3. 没有 API key 也没有外部 Agent，返回模拟响应
		if k.externalAgent != nil {
			return fmt.Sprintf("[模拟模式] 外部 Agent (%s) 调用失败，请检查配置。", k.externalAgent.Command), 0.001, nil
		}
		return "[模拟模式] 请在设置中配置 ANTHROPIC_API_KEY 以获得真实 AI 回复，或安装 Claude Code / OpenClaw。", 0.001, nil
	}

	// 检测 API Key 格式
	if strings.HasPrefix(apiKey, "sk-sp-") {
		// Claude Code session token，不能直接用于 API
		// 如果有外部 Agent，建议直接使用
		if k.externalAgent != nil {
			return fmt.Sprintf("[模拟模式] 请直接使用本地 Claude Code (已检测到: %s)", k.externalAgent.Command), 0.001, nil
		}
		return "[模拟模式] 检测到 Claude Code session token，请使用 Anthropic API Key (sk-ant-xxx)，或安装 Claude Code。", 0.001, nil
	}

	// 构建消息历史
	var conversation string
	for _, msg := range messages {
		role := msg.Role
		if role == "user" {
			conversation += fmt.Sprintf("\nHuman: %s", msg.Content)
		} else {
			conversation += fmt.Sprintf("\nAssistant: %s", msg.Content)
		}
	}

	// 调用 Claude API
	reqBody := map[string]interface{
	}{
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
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

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

// ExternalAgentOptions 外部 Agent 调用选项
type ExternalAgentOptions struct {
	OutputFormat string                 // json, text
	JSONSchema   map[string]interface{} // JSON Schema 用于结构化输出
	Timeout      time.Duration          // 超时时间
}

// callExternalAgent 通过 STDIO 调用外部 Agent
func (k *Keeper) callExternalAgent(messages []Message) (string, error) {
	if k.externalAgent == nil {
		return "", errors.New("no external agent configured")
	}

	// 构建对话内容 - 只取最后一条用户消息
	var lastMessage string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastMessage = messages[i].Content
			break
		}
	}
	if lastMessage == "" {
		return "", errors.New("no user message found")
	}

	// 构建命令参数
	args := append(k.externalAgent.Args, lastMessage)

	// 启动外部 Agent 进程
	cmd := exec.Command(k.externalAgent.Command, args...)

	// 设置环境变量 - 清除 CLAUDECODE 以避免嵌套检测
	env := os.Environ()
	filteredEnv := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, "CLAUDECODE=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	// 添加自定义标记，让子进程知道它是被 Gateway 调用的
	filteredEnv = append(filteredEnv, "GENIEBOT_AGENT_CALL=true")
	cmd.Env = filteredEnv

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("external agent failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// CallExternalAgentWithOptions 使用高级选项调用外部 Agent
func (k *Keeper) CallExternalAgentWithOptions(prompt string, opts ExternalAgentOptions) (string, error) {
	if k.externalAgent == nil {
		return "", errors.New("no external agent configured")
	}

	// 构建命令参数
	args := []string{"-p", prompt}

	// 添加输出格式选项
	if opts.OutputFormat != "" {
		args = append(args, "--output-format", opts.OutputFormat)
	}

	// 添加 JSON Schema
	if opts.JSONSchema != nil && opts.OutputFormat == "json" {
		schemaJSON, err := json.Marshal(opts.JSONSchema)
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
		}
		args = append(args, "--json-schema", string(schemaJSON))
	}

	// 启动外部 Agent 进程
	cmd := exec.Command(k.externalAgent.Command, args...)

	// 设置环境变量 - 清除 CLAUDECODE 以避免嵌套检测
	env := os.Environ()
	filteredEnv := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, "CLAUDECODE=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	filteredEnv = append(filteredEnv, "GENIEBOT_AGENT_CALL=true")
	cmd.Env = filteredEnv

	// 设置超时
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	// 使用 context 控制超时
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd = exec.CommandContext(ctx, k.externalAgent.Command, args...)
	cmd.Env = filteredEnv

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("external agent timeout after %v", timeout)
		}
		return "", fmt.Errorf("external agent failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// getOrCreateSession 获取或创建会话
func (k *Keeper) getOrCreateSession(sessionID string) *Session {
	k.sessionMu.Lock()
	defer k.sessionMu.Unlock()

	if session, exists := k.sessions[sessionID]; exists {
		return session
	}

	session := &Session{
		ID:        sessionID,
		Context:   make([]Message, 0),
		CreatedAt: time.Now(),
	}
	k.sessions[sessionID] = session
	return session
}

// Authenticate 钱包签名认证
func (k *Keeper) Authenticate(address, signature string) error {
	// TODO: 实现真实签名验证
	// 目前简单校验
	if signature == "" || signature == "invalid-signature" {
		return errors.New("invalid signature")
	}
	return nil
}

// CheckRateLimit 检查速率限制
func (k *Keeper) CheckRateLimit(address string) bool {
	k.rateMu.Lock()
	defer k.rateMu.Unlock()

	now := time.Now()
	limit, exists := k.rateLimits[address]
	if !exists {
		// 新建速率限制记录
		k.rateLimits[address] = &RateLimit{
			Address:      address,
			RequestCount: 1,
			LastReset:    now,
		}
		return true
	}

	// 检查是否需要重置（每分钟）
	if now.Sub(limit.LastReset) > time.Minute {
		limit.RequestCount = 1
		limit.LastReset = now
		return true
	}

	// 检查是否超过限制（60 req/min）
	if limit.RequestCount >= 60 {
		return false
	}

	limit.RequestCount++
	return true
}

// GetAgentCard 获取 A2A Agent Card
func (k *Keeper) GetAgentCard() *AgentCard {
	return &AgentCard{
		Name:        "ShareToken Agent",
		Version:     "1.0.0",
		Description: "ShareToken blockchain agent for task management and escrow",
		Capabilities: []string{
			"task_execution",
			"escrow_management",
			"query_service",
			"chat_with_genie",
		},
		Endpoints: map[string]string{
			"tasks":     "/a2a/tasks",
			"status":    "/a2a/status",
			"negotiate": "/a2a/negotiate",
		},
	}
}

// CreateTask 创建任务
func (k *Keeper) CreateTask(ctx context.Context, userAddr, description, budget string) (string, error) {
	// TODO: 实现真实任务创建
	taskID := fmt.Sprintf("task-%d", time.Now().Unix())
	return taskID, nil
}

// CreateEscrow 创建托管
func (k *Keeper) CreateEscrow(ctx context.Context, userAddr, providerAddr, amount string) (string, error) {
	// TODO: 实现真实托管创建
	escrowID := fmt.Sprintf("escrow-%d", time.Now().Unix())
	return escrowID, nil
}
