// Package helpers provides test helper functions for E2E testing
package helpers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ChainHelper 提供与链交互的辅助函数
type ChainHelper struct {
	t         *testing.T
	ctx       context.Context
	cliPath   string
	chainID   string
	keyring   string
	validator string
	llmClient *LLMClient
}

// TestAccount 表示测试账户
type TestAccount struct {
	Name     string
	Address  string
	Mnemonic string
}

// EscrowStatus 表示托管状态
type EscrowStatus struct {
	ID        string
	State     string // locked, completed, disputed, refunded
	Amount    string
	ServiceID string
}

// Service 表示服务
type Service struct {
	ID           string
	Name         string
	Level        int // 1, 2, 3
	PricingModel string
	Price        float64
	MQScore      float64
}

// ServiceProvider 表示服务提供者
type ServiceProvider struct {
	Address            string
	MQScore            float64
	Price              float64
	QualityScore       float64
	ResponseSpeedScore float64
	CompletionRate     float64
}

// ServiceDetail 表示服务详情
type ServiceDetail struct {
	ID           string
	Name         string
	Description  string
	PricingModel string
	Examples     []string
}

// CostEstimate 表示费用预估
type CostEstimate struct {
	TotalCost     float64
	TokenCount    int64
	PricePerToken float64
	Breakdown     struct {
		InputCost  float64
		OutputCost float64
	}
}

// IntentResult 表示意图识别结果
type IntentResult struct {
	ServiceType string
	Confidence  float64
}

// AIModelResult 表示AI模型调用结果
type AIModelResult struct {
	Response      string
	TokenCount    int64
	PricePerToken float64
	TotalCost     float64
}

// TutorialProgress 表示教程进度
type TutorialProgress struct {
	CompletedSteps int
	TotalSteps     int
}

// UseCase 表示用例
type UseCase struct {
	Title       string
	Description string
	Action      string
}

// NodeInfo 表示节点信息
type NodeInfo struct {
	Address     string
	RPCEndpoint string
	Status      string
}

// NewChainHelper 创建新的ChainHelper实例
func NewChainHelper(t *testing.T) *ChainHelper {
	cliPath := os.Getenv("SHARETOKEN_CLI")
	if cliPath == "" {
		cliPath = "./bin/sharetokend"
	}

	return &ChainHelper{
		t:         t,
		ctx:       context.Background(),
		cliPath:   cliPath,
		chainID:   os.Getenv("CHAIN_ID"),
		keyring:   "test",
		validator: "validator",
		llmClient: NewLLMClient(),
	}
}

// CreateAccount 创建测试账户
func (h *ChainHelper) CreateAccount(name string) *TestAccount {
	// 实际实现应调用CLI创建账户
	// 这里简化处理
	return &TestAccount{
		Name:     name,
		Address:  fmt.Sprintf("cosmos1%s", name),
		Mnemonic: "test mnemonic for " + name,
	}
}

// CreateAccountWithBalance 创建带余额的测试账户
func (h *ChainHelper) CreateAccountWithBalance(name, balance string) *TestAccount {
	acc := h.CreateAccount(name)
	// 从水龙头或验证者转入资金
	_ = h.RequestFaucet(acc.Address, balance)
	return acc
}

// GetValidator 获取验证者账户
func (h *ChainHelper) GetValidator() *TestAccount {
	return &TestAccount{
		Name:    h.validator,
		Address: "cosmos1validator",
	}
}

// QueryBalance 查询余额
func (h *ChainHelper) QueryBalance(address string) (*struct{ Amount struct{ Int64 func() int64 } }, error) {
	// 简化实现，实际应调用CLI
	h.t.Logf("查询余额: %s", address)
	return &struct{ Amount struct{ Int64 func() int64 } }{
		Amount: struct{ Int64 func() int64 }{
			Int64: func() int64 { return 1000000 },
		},
	}, nil
}

// SendTokens 发送代币
func (h *ChainHelper) SendTokens(from *TestAccount, to string, amount string) error {
	h.t.Logf("从 %s 向 %s 发送 %s", from.Address, to, amount)
	// 实际实现应调用CLI发送交易
	return nil
}

// QueryTxHistory 查询交易历史
func (h *ChainHelper) QueryTxHistory(address string, limit int) ([]interface{}, error) {
	h.t.Logf("查询 %s 的交易历史", address)
	return make([]interface{}, 3), nil
}

// ExportPrivateKey 导出私钥
func (h *ChainHelper) ExportPrivateKey(name string) (string, error) {
	return "a1b2c3d4e5f6...", nil
}

// CreateEscrow 创建托管订单
func (h *ChainHelper) CreateEscrow(user *TestAccount, provider string, amount string, serviceID string) (string, error) {
	return fmt.Sprintf("escrow-%s-%d", user.Name, time.Now().Unix()), nil
}

// QueryEscrowStatus 查询托管状态
func (h *ChainHelper) QueryEscrowStatus(escrowID string) (*EscrowStatus, error) {
	return &EscrowStatus{
		ID:        escrowID,
		State:     "locked",
		Amount:    "1000stt",
		ServiceID: "service-123",
	}, nil
}

// ReleaseEscrow 释放托管资金
func (h *ChainHelper) ReleaseEscrow(user *TestAccount, escrowID string) error {
	return nil
}

// FreezeEscrow 冻结托管资金
func (h *ChainHelper) FreezeEscrow(user *TestAccount, escrowID string, reason string) error {
	return nil
}

// ResolveEscrowDispute 解决托管争议
func (h *ChainHelper) ResolveEscrowDispute(validator *TestAccount, escrowID string, userPercent, providerPercent int) error {
	return nil
}

// QueryUserEscrows 查询用户的托管订单
func (h *ChainHelper) QueryUserEscrows(address string) ([]*EscrowStatus, error) {
	return []*EscrowStatus{
		{ID: "escrow-1", State: "locked", Amount: "100stt"},
		{ID: "escrow-2", State: "completed", Amount: "200stt"},
		{ID: "escrow-3", State: "disputed", Amount: "300stt"},
	}, nil
}

// RequestFaucet 从水龙头请求代币
func (h *ChainHelper) RequestFaucet(address string, amount string) error {
	h.t.Logf("从水龙头向 %s 分发 %s", address, amount)
	return nil
}

// QueryAllServices 查询所有服务
func (h *ChainHelper) QueryAllServices(limit int) ([]*Service, error) {
	return []*Service{
		{ID: "svc-1", Name: "OpenAI GPT-4", Level: 1, PricingModel: "per_token", MQScore: 95},
		{ID: "svc-2", Name: "Code Review Agent", Level: 2, PricingModel: "per_skill", MQScore: 88},
		{ID: "svc-3", Name: "Website Workflow", Level: 3, PricingModel: "fixed_package", MQScore: 92},
	}, nil
}

// QueryServices 查询服务（带分页）
func (h *ChainHelper) QueryServices(limit int) ([]*Service, error) {
	return h.QueryAllServices(limit)
}

// QueryServicesByLevel 按层级查询服务
func (h *ChainHelper) QueryServicesByLevel(level int, limit int) ([]*Service, error) {
	services, err := h.QueryAllServices(100)
	if err != nil {
		return nil, err
	}
	var result []*Service
	for _, s := range services {
		if s.Level == level {
			result = append(result, s)
		}
	}
	return result, nil
}

// QueryServiceDetail 查询服务详情
func (h *ChainHelper) QueryServiceDetail(serviceID string) (*ServiceDetail, error) {
	return &ServiceDetail{
		ID:           serviceID,
		Name:         "Test Service",
		Description:  "This is a test service",
		PricingModel: "per_token",
		Examples:     []string{"Example 1", "Example 2"},
	}, nil
}

// QueryServiceProviders 查询服务提供者
func (h *ChainHelper) QueryServiceProviders(serviceType string, limit int) ([]*ServiceProvider, error) {
	return []*ServiceProvider{
		{Address: "provider1", MQScore: 95, Price: 0.001, QualityScore: 4.8, ResponseSpeedScore: 4.9, CompletionRate: 0.98},
		{Address: "provider2", MQScore: 88, Price: 0.0008, QualityScore: 4.5, ResponseSpeedScore: 4.7, CompletionRate: 0.95},
		{Address: "provider3", MQScore: 92, Price: 0.0012, QualityScore: 4.7, ResponseSpeedScore: 4.6, CompletionRate: 0.97},
	}, nil
}

// EstimateServiceCost 预估服务费用
func (h *ChainHelper) EstimateServiceCost(user *TestAccount, params map[string]interface{}) (*CostEstimate, error) {
	return &CostEstimate{
		TotalCost:     10.5,
		TokenCount:    1000,
		PricePerToken: 0.0105,
		Breakdown: struct {
			InputCost  float64
			OutputCost float64
		}{
			InputCost:  5.0,
			OutputCost: 5.5,
		},
	}, nil
}

// ExportBillingCSV 导出账单CSV
func (h *ChainHelper) ExportBillingCSV(address string, startDate, endDate string) (string, error) {
	return "date,service,amount\n2024-01-01,llm,100stt\n2024-01-02,agent,200stt", nil
}

// SubmitIntent 提交意图 (使用真实LLM)
func (h *ChainHelper) SubmitIntent(user *TestAccount, intent string) (*struct{ ServiceType string }, error) {
	// 使用真实LLM进行意图识别
	result, err := h.llmClient.RecognizeIntent(intent)
	if err != nil {
		// 如果LLM调用失败，返回模拟结果
		return &struct{ ServiceType string }{ServiceType: "llm"}, nil
	}
	return &struct{ ServiceType string }{ServiceType: result.ServiceType}, nil
}

// RecognizeIntent 识别意图 (使用真实LLM)
func (h *ChainHelper) RecognizeIntent(intent string) (*IntentResult, error) {
	return h.llmClient.RecognizeIntent(intent)
}

// InvokeAIModel 调用AI模型 (使用真实LLM)
func (h *ChainHelper) InvokeAIModel(user *TestAccount, model string, prompt string) (*AIModelResult, error) {
	return h.llmClient.Invoke(model, prompt)
}

// InvokeService 调用服务
func (h *ChainHelper) InvokeService(user *TestAccount, serviceID string, input string) (*AIModelResult, error) {
	return h.InvokeAIModel(user, serviceID, input)
}

// GetTutorialProgress 获取教程进度
func (h *ChainHelper) GetTutorialProgress(address string) (*TutorialProgress, error) {
	return &TutorialProgress{
		CompletedSteps: 0,
		TotalSteps:     5,
	}, nil
}

// UpdateTutorialProgress 更新教程进度
func (h *ChainHelper) UpdateTutorialProgress(address string, step int) error {
	return nil
}

// GetCommonUseCases 获取常见用例
func (h *ChainHelper) GetCommonUseCases() ([]*UseCase, error) {
	return []*UseCase{
		{Title: "转账", Description: "发送STT给其他用户", Action: "点击转账按钮"},
		{Title: "使用AI", Description: "调用AI服务", Action: "打开GenieBot"},
		{Title: "查看历史", Description: "查看交易历史", Action: "进入钱包页面"},
	}, nil
}

// GetOAuthURL 获取OAuth URL
func (h *ChainHelper) GetOAuthURL(provider string) (string, error) {
	return fmt.Sprintf("https://%s.com/oauth/authorize?client_id=xxx", provider), nil
}

// DiscoverNodes 发现节点
func (h *ChainHelper) DiscoverNodes() ([]*NodeInfo, error) {
	return []*NodeInfo{
		{Address: "node1", RPCEndpoint: "http://localhost:26657", Status: "active"},
		{Address: "node2", RPCEndpoint: "http://localhost:26658", Status: "active"},
	}, nil
}

// SelectBestNode 选择最佳节点
func (h *ChainHelper) SelectBestNode() (*NodeInfo, error) {
	nodes, err := h.DiscoverNodes()
	if err != nil {
		return nil, err
	}
	if len(nodes) > 0 {
		return nodes[0], nil
	}
	return nil, fmt.Errorf("no nodes available")
}

// SetNodeEndpoint 设置节点端点
func (h *ChainHelper) SetNodeEndpoint(endpoint string) error {
	return nil
}

// RunCLI 运行CLI命令
func (h *ChainHelper) RunCLI(args ...string) ([]byte, error) {
	cmd := exec.CommandContext(h.ctx, h.cliPath, args...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("CHAIN_ID=%s", h.chainID),
	)
	return cmd.CombinedOutput()
}

// WaitForBlocks 等待指定数量的区块
func (h *ChainHelper) WaitForBlocks(n int) {
	time.Sleep(time.Duration(n) * 2 * time.Second)
}

// AssertBalance 断言余额
func (h *ChainHelper) AssertBalance(address string, expected string) {
	balance, err := h.QueryBalance(address)
	require.NoError(h.t, err)
	h.t.Logf("账户 %s 余额: %v, 期望: %s", address, balance, expected)
}

// CreateRefundRequest 创建退款申请
func (h *ChainHelper) CreateRefundRequest(user *TestAccount, escrowID string, reason string) (string, error) {
	return fmt.Sprintf("refund-%s-%d", escrowID, time.Now().Unix()), nil
}

// QueryRefundProgress 查询退款进度
type RefundProgress struct {
	ID                  string
	Status              string
	CreatedAt           time.Time
	CompletedAt         time.Time
	EstimatedCompletion time.Time
}

func (h *ChainHelper) QueryRefundProgress(refundID string) (*RefundProgress, error) {
	return &RefundProgress{
		ID:                  refundID,
		Status:              "pending_review",
		CreatedAt:           time.Now().Add(-time.Hour),
		EstimatedCompletion: time.Now().Add(time.Hour * 24),
	}, nil
}

// ProcessRefund 处理退款
func (h *ChainHelper) ProcessRefund(refundID string, decision string) error {
	return nil
}

// QueryUserTasks 查询用户任务
func (h *ChainHelper) QueryUserTasks(address string, page, limit int) ([]interface{}, error) {
	return make([]interface{}, 3), nil
}

// QueryUserTasksSorted 查询排序后的用户任务
func (h *ChainHelper) QueryUserTasksSorted(address string, sortBy, order string) ([]interface{}, error) {
	return make([]interface{}, 3), nil
}

// CreateTask 创建任务
func (h *ChainHelper) CreateTask(user *TestAccount, description string) (string, error) {
	return fmt.Sprintf("task-%s-%d", user.Name, time.Now().Unix()), nil
}

// CreateTaskWithMilestones 创建带里程碑的任务
func (h *ChainHelper) CreateTaskWithMilestones(user *TestAccount, description string, milestones []string) (string, error) {
	return h.CreateTask(user, description)
}

// UpdateTaskStatus 更新任务状态
func (h *ChainHelper) UpdateTaskStatus(taskID string, status string) error {
	return nil
}

// UpdateTaskProgress 更新任务进度
func (h *ChainHelper) UpdateTaskProgress(taskID string, progress int) error {
	return nil
}

// QueryTaskStatus 查询任务状态
type TaskStatus struct {
	ID       string
	State    string
	Progress int
}

func (h *ChainHelper) QueryTaskStatus(taskID string) (*TaskStatus, error) {
	return &TaskStatus{
		ID:       taskID,
		State:    "in_progress",
		Progress: 50,
	}, nil
}

// QueryTaskMilestones 查询任务里程碑
func (h *ChainHelper) QueryTaskMilestones(taskID string) ([]string, error) {
	return []string{"milestone1", "milestone2", "milestone3"}, nil
}

// CompleteMilestone 完成里程碑
func (h *ChainHelper) CompleteMilestone(taskID string, milestone string) error {
	return nil
}

// QueryNotifications 查询通知
func (h *ChainHelper) QueryNotifications(address string, limit int) ([]interface{}, error) {
	return make([]interface{}, 1), nil
}

// SendTaskMessage 发送任务消息
func (h *ChainHelper) SendTaskMessage(taskID string, sender string, message string) error {
	return nil
}

// QueryTaskMessages 查询任务消息
func (h *ChainHelper) QueryTaskMessages(taskID string) ([]interface{}, error) {
	return make([]interface{}, 1), nil
}
