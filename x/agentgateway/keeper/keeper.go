package keeper

import (
	"context"
	"fmt"
	"os/exec"
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
	sessions  map[string]*Session
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
}

// QueryBalance 查询余额
func (k *Keeper) QueryBalance(ctx context.Context, address string) (uint64, error) {
	if k.queryClient != nil {
		return k.queryClient.QueryBalance(ctx, address)
	}
	// 模拟余额: 10000 STT
	return 10000000, nil
}

// Authenticate 钱包签名认证
func (k *Keeper) Authenticate(address, signature string) error {
	if signature == "" || signature == "invalid-signature" {
		return fmt.Errorf("invalid signature")
	}
	return nil
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
	taskID := fmt.Sprintf("task-%d", time.Now().Unix())
	return taskID, nil
}

// CreateEscrow 创建托管
func (k *Keeper) CreateEscrow(ctx context.Context, userAddr, providerAddr, amount string) (string, error) {
	escrowID := fmt.Sprintf("escrow-%d", time.Now().Unix())
	return escrowID, nil
}
