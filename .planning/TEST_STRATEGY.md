# ShareTokens Test Strategy

**Defined:** 2025-03-02
**Core Value:** Every great idea deserves the tokens to make it real.

---

## 架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ShareTokens 测试架构                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  核心模块（每个节点必须有）- 测试优先级 P0                                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  ├── P2P通信（Cosmos SDK + CometBFT）                                       │
│  ├── 身份账号（实名制、身份注册表）                                          │
│  ├── 钱包（Cosmos SDK Auth + Keplr）                                        │
│  ├── 服务市场（三层服务：LLM/Agent/Workflow）← 核心业务                      │
│  ├── 托管支付（争议时冻结）                                                  │
│  ├── 德商系统（零和博弈、加权评分、衰减机制）                                │
│  └── 争议仲裁（多AI评估、裁决算法）                                         │
│                                                                             │
│  服务提供者插件 - 测试优先级 P1                                               │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  ├── LLM API Key托管插件（加密存储、OpenFang 16层安全）                     │
│  ├── Agent执行器插件（28+ Agent模板、7个Hands）                             │
│  └── Workflow执行器插件（编排、执行、监控）                                  │
│                                                                             │
│  用户插件 - 测试优先级 P2                                                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  └── 小灯界面插件（AI对话、想法孵化、资源匹配）                              │
│                                                                             │
│  辅助模块 - 测试优先级 P1                                                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  ├── 汇率层（Chainlink）                                                    │
│  ├── 想法系统                                                               │
│  └── 任务市场                                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. 测试概述

### 1.1 测试目标

| 目标 | 描述 | 验收标准 |
|------|------|----------|
| **功能正确性** | 所有需求功能按预期工作 | 所有测试用例通过 |
| **性能达标** | 系统满足性能指标 | TPS > 1000, API < 200ms |
| **安全可靠** | 保护用户资产和数据 | 无高危漏洞 |
| **用户体验** | 界面流畅，响应及时 | 无阻塞性Bug |

### 1.2 测试覆盖率目标

| 测试类型 | 覆盖率目标 | 优先级 |
|----------|------------|--------|
| 单元测试 | 80% | 高 |
| 集成测试 | 60% | 高 |
| E2E 测试 | 核心流程 100% | 高 |
| 性能测试 | 关键路径 100% | 中 |
| 安全测试 | 全量审计 | 高 |

### 1.3 测试范围

**测试优先级:** 核心模块 > 服务提供者插件 > 用户插件

#### 核心模块（每个节点必须有）- 测试优先级 P0

| 模块 | 说明 | 链上模块 | 测试重点 |
|------|------|----------|----------|
| **P2P通信** | 共识、区块、验证者、网络 | CometBFT 内置 | 网络连通性、共识正确性 |
| **身份账号** | 实名制、身份注册表、账户 | x/identity | 身份验证、防重复注册 |
| **钱包** | 账户管理、余额、交易 | Cosmos SDK Auth + Keplr | 交易签名、余额正确性 |
| **服务市场** | 三层服务（LLM/Agent/Workflow）| x/service | 服务注册、发现、定价、路由 |
| **托管支付** | 托管账户、支付、结算 | x/escrow | 托管创建、释放、退款 |
| **德商系统** | 德商评分、零和博弈、衰减 | x/mq | 评分计算、衰减机制 |
| **争议仲裁** | 争议、裁决算法、证据 | x/dispute | 评审团组建、投票、裁决 |

#### 服务提供者插件 - 测试优先级 P1

| 插件 | 说明 | 测试重点 |
|------|------|----------|
| **LLM API Key托管** | API Key 加密存储、请求代理 | Key 安全存储、请求转发 |
| **Agent执行器** | Agent 运行时、资源管理 (OpenFang) | Agent 执行、工具调用 |
| **Workflow执行器** | Workflow 编排、执行监控 | 流程执行、人工介入点 |

#### 用户插件 - 测试优先级 P2

| 插件 | 说明 | 测试重点 |
|------|------|----------|
| **小灯界面** | AI 对话界面、想法孵化 | 用户交互、意图识别、服务调用 |

#### 辅助模块 - 测试优先级 P1

| 模块 | 说明 | 链上模块 | 测试重点 |
|------|------|----------|----------|
| **汇率层** | 汇率快照、价格映射 | x/exchange + Chainlink | 价格更新、汇率转换 |
| **想法系统** | 想法、众筹、贡献 | x/idea | 众筹流程、贡献权重 |
| **任务市场** | 任务、申请、里程碑 | x/task | 任务分配、里程碑验收 |

---

## 2. 单元测试

### 2.1 核心模块 (Go - Cosmos SDK)

**测试框架:** Go testing + testify

**目录结构:**
```
x/
├── identity/              # 身份账号模块
│   ├── keeper/
│   │   ├── keeper_test.go
│   │   ├── identity_test.go
│   │   └── msg_server_test.go
│   └── types/
│       └── identity_test.go
├── service/               # 服务市场模块 (核心)
│   ├── keeper/
│   │   ├── service_test.go
│   │   ├── registration_test.go
│   │   ├── discovery_test.go
│   │   ├── pricing_test.go
│   │   └── router_test.go
│   └── types/
│       └── service_test.go
├── escrow/                # 托管支付模块
│   └── keeper/
│       └── escrow_test.go
├── mq/                    # 德商系统模块
│   └── keeper/
│       └── mq_test.go
├── dispute/               # 争议仲裁模块
│   └── keeper/
│       └── dispute_test.go
├── exchange/              # 汇率层模块
│   └── keeper/
│       └── exchange_test.go
├── idea/                  # 想法系统模块
│   └── keeper/
│       └── idea_test.go
└── task/                  # 任务市场模块
    └── keeper/
        └── task_test.go
```

**Keeper 测试示例:**
```go
// x/identity/keeper/identity_test.go
package keeper_test

import (
    "testing"
    "github.com/stretchr/testify/require"
    "sharetokens/x/identity/types"
)

func TestKeeper_RegisterIdentity(t *testing.T) {
    // Setup test context
    keeper, ctx := setupKeeper(t)

    // Test case: Register new identity
    msg := &types.MsgRegisterIdentity{
        Creator:  "cosmos1abc...",
        Platform: "github",
        Handle:   "testuser",
    }

    err := keeper.RegisterIdentity(ctx, msg)
    require.NoError(t, err)

    // Verify identity was stored
    identity, found := keeper.GetIdentity(ctx, "cosmos1abc...")
    require.True(t, found)
    require.Equal(t, "github", identity.Platform)
}

func TestKeeper_DuplicateIdentity(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    // Register first identity
    msg := &types.MsgRegisterIdentity{
        Creator:  "cosmos1abc...",
        Platform: "github",
        Handle:   "testuser",
    }
    keeper.RegisterIdentity(ctx, msg)

    // Try to register duplicate
    err := keeper.RegisterIdentity(ctx, msg)
    require.Error(t, err)
    require.Contains(t, err.Error(), "already registered")
}

func TestKeeper_VerifyMerkleProof(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    tests := []struct {
        name    string
        proof   []byte
        root    []byte
        isValid bool
    }{
        {"valid proof", validProof, validRoot, true},
        {"invalid proof", invalidProof, validRoot, false},
        {"empty proof", []byte{}, validRoot, false},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := keeper.VerifyMerkleProof(ctx, tc.proof, tc.root)
            require.Equal(t, tc.isValid, result)
        })
    }
}
```

**德商系统测试示例:**
```go
// x/mq/keeper/mq_test.go
func TestKeeper_InitializeMQ(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    addr := "cosmos1abc..."
    keeper.InitializeMQ(ctx, addr)

    mq, found := keeper.GetMQ(ctx, addr)
    require.True(t, found)
    require.Equal(t, int64(100), mq.Value) // 初始值 100
}

func TestKeeper_MQDecay(t *testing.T) {
    keeper, ctx := setupKeeper(t)

    addr := "cosmos1abc..."
    keeper.InitializeMQ(ctx, addr)

    // 模拟 30 天不活跃
    ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 86400 * 30)

    keeper.ApplyDecay(ctx, addr)

    mq, _ := keeper.GetMQ(ctx, addr)
    // 每日衰减 0.1%, 30天后约 97%
    require.Less(t, mq.Value, int64(100))
    require.Greater(t, mq.Value, int64(90))
}

func TestKeeper_MQLevel(t *testing.T) {
    tests := []struct {
        mq    int64
        level string
    }{
        {30, "Newcomer"},   // Lv1
        {80, "Member"},     // Lv2
        {150, "Trusted"},   // Lv3
        {300, "Expert"},    // Lv4
        {600, "Guardian"},  // Lv5
    }

    keeper, _ := setupKeeper(t)
    for _, tc := range tests {
        level := keeper.CalculateLevel(tc.mq)
        require.Equal(t, tc.level, level)
    }
}
```

### 2.2 服务提供者插件 (TypeScript)

**测试框架:** Jest + ts-jest

**目录结构:**
```
plugins/
├── llm-provider/              # LLM API Key 托管插件
│   ├── src/
│   └── __tests__/
│       ├── key-management.test.ts
│       ├── proxy.test.ts
│       └── mocks/
│           └── llm.mock.ts
├── agent-provider/            # Agent 执行器插件 (OpenFang)
│   └── __tests__/
│       ├── agent-execution.test.ts
│       ├── tool-calling.test.ts
│       └── resource-management.test.ts
├── workflow-provider/         # Workflow 执行器插件
│   └── __tests__/
│       ├── executor.test.ts
│       ├── orchestrator.test.ts
│       └── human-gate.test.ts
└── oracle/                    # 价格预言机 (Chainlink 集成)
    └── __tests__/
        └── price.test.ts
```

**LLM Provider 测试示例:**
```typescript
// plugins/llm-provider/__tests__/key-management.test.ts
import { LLMProviderPlugin } from '../src/index';
import { mockKeyStorage } from './mocks/storage.mock';

describe('LLMProviderPlugin - Key Management', () => {
  let plugin: LLMProviderPlugin;

  beforeEach(() => {
    plugin = new LLMProviderPlugin(mockKeyStorage);
  });

  describe('registerApiKey', () => {
    it('should encrypt and store API key', async () => {
      const provider = 'openai';
      const apiKey = 'sk-test-xxx';

      await plugin.registerApiKey(provider, apiKey);

      // 验证 Key 被加密存储
      const stored = mockKeyStorage.get(provider);
      expect(stored).toBeDefined();
      expect(stored).not.toBe(apiKey); // 不应存储明文
    });

    it('should rotate API key securely', async () => {
      const provider = 'openai';
      await plugin.registerApiKey(provider, 'sk-old');

      await plugin.rotateApiKey(provider, 'sk-new');

      const stored = mockKeyStorage.get(provider);
      expect(plugin.decrypt(stored)).toBe('sk-new');
    });
  });
});
```

**Agent Provider (OpenFang) 测试示例:**
```typescript
// plugins/agent-provider/__tests__/agent-execution.test.ts
import { AgentProviderPlugin } from '../src/index';
import { OpenFangClient } from '../src/openfang-client';

describe('AgentProviderPlugin', () => {
  let plugin: AgentProviderPlugin;
  let mockOpenFang: jest.Mocked<OpenFangClient>;

  beforeEach(() => {
    mockOpenFang = createMockOpenFangClient();
    plugin = new AgentProviderPlugin(mockOpenFang);
  });

  describe('executeTask', () => {
    it('should execute Coder Agent task', async () => {
      const request = {
        agentId: 'coder-agent',
        task: 'Write a hello world function in Python',
      };

      const result = await plugin.executeTask(request);

      expect(result.status).toBe('completed');
      expect(result.artifacts).toBeDefined();
      expect(mockOpenFang.execute).toHaveBeenCalledWith(
        expect.objectContaining({ agentType: 'coder' })
      );
    });

    it('should handle tool calling correctly', async () => {
      const request = {
        agentId: 'researcher-agent',
        task: 'Search for latest AI news',
        requiredTools: ['web_browser'],
      };

      const result = await plugin.executeTask(request);

      expect(result.steps.some(s => s.tool === 'web_browser')).toBe(true);
    });
  });
});
```

**Workflow Provider 测试示例:**
```typescript
// plugins/workflow-provider/__tests__/executor.test.ts
import { WorkflowProviderPlugin } from '../src/index';

describe('WorkflowProviderPlugin', () => {
  describe('SoftwareWorkflow', () => {
    it('should execute software development workflow', async () => {
      const plugin = new WorkflowProviderPlugin();
      const context = {
        workflowId: 'software-dev',
        idea: 'Build a todo app',
        deliverables: ['source_code', 'tests', 'docs'],
      };

      const result = await plugin.startExecution(context);

      expect(result.status).toBe('running');
      expect(result.currentStep).toBe('requirement_analysis');
    });

    it('should pause at human review step', async () => {
      const plugin = new WorkflowProviderPlugin();
      const context = {
        workflowId: 'software-dev',
        idea: 'Build a todo app',
        requiresReview: true,
      };

      const execution = await plugin.startExecution(context);
      // 模拟执行到审核步骤
      await advanceToStep(execution.id, 'code_review');

      const status = await plugin.getStatus(execution.id);
      expect(status.status).toBe('waiting_input');
      expect(status.currentStep).toBe('code_review');
    });
  });
});
```

### 2.3 用户插件 - 小灯界面 (TypeScript)

**测试框架:** Jest + ts-jest

**目录结构:**
```
plugins/
└── xiaodeng-client/           # 小灯界面插件
    └── __tests__/
        ├── intent-recognition.test.ts
        ├── service-routing.test.ts
        └── marketplace-integration.test.ts
```

**小灯意图识别测试示例:**
```typescript
// plugins/xiaodeng-client/__tests__/intent-recognition.test.ts
import { XiaodengClient } from '../src/index';

describe('XiaodengClient - Intent Recognition', () => {
  let client: XiaodengClient;

  beforeEach(() => {
    client = new XiaodengClient(mockMarketplace);
  });

  describe('determineServiceLevel', () => {
    it('should route simple chat to Level 1 (LLM)', async () => {
      const intent = await client.analyzeIntent('Hello, how are you?');

      expect(intent.type).toBe('chat');
      expect(intent.serviceLevel).toBe(1);
    });

    it('should route specific task to Level 2 (Agent)', async () => {
      const intent = await client.analyzeIntent('Write a Python function to sort a list');

      expect(intent.type).toBe('task');
      expect(intent.serviceLevel).toBe(2);
      expect(intent.requiredCapabilities).toContain('code_generation');
    });

    it('should route complex idea to Level 3 (Workflow)', async () => {
      const intent = await client.analyzeIntent('I want to build a mobile app for task management');

      expect(intent.type).toBe('idea');
      expect(intent.serviceLevel).toBe(3);
    });
  });
});
```

### 2.4 前端组件 (小灯 UI)

**测试框架:** Vitest + React Testing Library

**目录结构:**
```
web/
├── src/
│   ├── components/
│   └── __tests__/
│       ├── Chat.test.tsx
│       ├── IdeaCard.test.tsx
│       ├── WalletConnect.test.tsx
│       └── ServiceBrowser.test.tsx
└── vitest.config.ts
```

**组件测试示例:**
```typescript
// web/src/__tests__/Chat.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Chat } from '../components/Chat';

describe('Chat Component', () => {
  it('should render chat input', () => {
    render(<Chat />);
    expect(screen.getByPlaceholderText(/ask xiaodeng/i)).toBeInTheDocument();
  });

  it('should send message and display response', async () => {
    const mockSend = vi.fn().mockResolvedValue('AI response');
    render(<Chat onSend={mockSend} />);

    const input = screen.getByPlaceholderText(/ask xiaodeng/i);
    fireEvent.change(input, { target: { value: 'Hello' } });
    fireEvent.submit(input.closest('form')!);

    await waitFor(() => {
      expect(mockSend).toHaveBeenCalledWith('Hello');
    });
  });

  it('should display typing indicator while waiting', async () => {
    const mockSend = vi.fn().mockImplementation(() => new Promise(r => setTimeout(r, 100)));
    render(<Chat onSend={mockSend} />);

    const input = screen.getByPlaceholderText(/ask xiaodeng/i);
    fireEvent.change(input, { target: { value: 'Hello' } });
    fireEvent.submit(input.closest('form')!);

    expect(screen.getByTestId('typing-indicator')).toBeInTheDocument();
  });
});
```

**钱包连接测试:**
```typescript
// web/src/__tests__/WalletConnect.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { WalletConnect } from '../components/WalletConnect';

describe('WalletConnect', () => {
  it('should show connect button when not connected', () => {
    render(<WalletConnect connected={false} />);
    expect(screen.getByText(/connect wallet/i)).toBeInTheDocument();
  });

  it('should show address when connected', () => {
    render(<WalletConnect connected={true} address="cosmos1abc...xyz" />);
    expect(screen.getByText(/cosmos1abc...xyz/i)).toBeInTheDocument();
  });

  it('should call onConnect when button clicked', () => {
    const onConnect = vi.fn();
    render(<WalletConnect connected={false} onConnect={onConnect} />);
    fireEvent.click(screen.getByText(/connect wallet/i));
    expect(onConnect).toHaveBeenCalled();
  });
});
```

---

## 3. 集成测试

### 3.1 核心模块间集成

**测试矩阵:**

| 模块组合 | 测试场景 | 验证点 |
|----------|----------|--------|
| identity + mq | 用户注册后初始化德商 | MQ = 100 |
| service + escrow | 服务交易托管 | 托管状态正确 |
| dispute + escrow + mq | 争议仲裁后结算 | 德商+托管正确分配 |
| service + exchange | 服务定价汇率转换 | 价格计算正确 |
| idea + escrow | 想法众筹托管 | 托管创建成功 |
| task + escrow | 任务支付托管 | 里程碑释放正确 |

**服务市场集成测试示例:**
```go
// x/test/integration/service_escrow_test.go
package integration_test

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestServiceToEscrow(t *testing.T) {
    // Setup 集成环境
    app := setupTestApp()
    ctx := app.NewContext(false)

    // 1. 注册服务提供者 (Level 1: LLM)
    provider := "cosmos1provider..."
    app.ServiceKeeper.RegisterService(ctx, &types.ServiceRegistration{
        Provider:    provider,
        Level:       1,  // LLM 服务
        Name:        "GPT-4 Service",
        Category:    "llm_chat",
        Pricing:     types.PricingModel{Type: "per_token", PricePerToken: 100},
    })

    // 2. 用户发起服务请求
    requester := "cosmos1requester..."
    requestID := app.ServiceKeeper.CreateRequest(ctx, &types.ServiceRequest{
        Consumer:  requester,
        Level:     1,
        Prompt:    "Hello, world!",
        Model:     "gpt-4",
        MaxTokens: 100,
    })

    // 3. 验证托管创建
    escrow, found := app.EscrowKeeper.GetEscrow(ctx, requestID)
    require.True(t, found)
    require.Equal(t, requester, escrow.Buyer)
    require.Equal(t, provider, escrow.Seller)
}

func TestDisputeResolution(t *testing.T) {
    app := setupTestApp()
    ctx := app.NewContext(false)

    // 初始化双方德商
    buyer := "cosmos1buyer..."
    seller := "cosmos1seller..."
    app.MQKeeper.InitializeMQ(ctx, buyer)
    app.MQKeeper.InitializeMQ(ctx, seller)

    initialBuyerMQ, _ := app.MQKeeper.GetMQ(ctx, buyer)
    initialSellerMQ, _ := app.MQKeeper.GetMQ(ctx, seller)

    // 创建争议
    disputeID := app.DisputeKeeper.CreateDispute(ctx, buyer, seller, 500)

    // 组建评审团并投票
    app.DisputeKeeper.FormJury(ctx, disputeID)
    app.DisputeKeeper.Vote(ctx, disputeID, buyer, true) // 买方胜

    // 验证德商变化 (零和博弈)
    finalBuyerMQ, _ := app.MQKeeper.GetMQ(ctx, buyer)
    finalSellerMQ, _ := app.MQKeeper.GetMQ(ctx, seller)

    delta := finalBuyerMQ.Value - initialBuyerMQ.Value
    require.Equal(t, delta, initialSellerMQ.Value-finalSellerMQ.Value)
}
```

### 3.2 核心模块与插件集成

**测试场景:**

| 场景 | 链上操作 | 插件服务 | 验证点 |
|------|----------|----------|--------|
| LLM 服务请求 | 创建服务请求 | LLM Provider 执行 | 结果验证+结算 |
| Agent 任务执行 | 创建任务请求 | Agent Provider (OpenFang) | 任务完成+结算 |
| Workflow 流程执行 | 创建流程请求 | Workflow Provider | 交付物验证+结算 |
| 汇率更新 | 定时触发 | Chainlink 预言机 | 价格更新 |
| 想法评估 | 提交想法 | 多 AI 评估 | 结果上链 |

**事件监听测试:**
```typescript
// plugins/test/event-listener.test.ts
import { EventListener } from '../src/event-listener';
import { createTestClient } from './test-utils';

describe('EventListener Integration', () => {
  let client: TestClient;
  let listener: EventListener;

  beforeAll(async () => {
    client = await createTestClient();
    listener = new EventListener(client);
  });

  it('should listen to ServiceRequestCreated events', async () => {
    const events: any[] = [];
    listener.on('ServiceRequestCreated', (e) => events.push(e));

    // 在链上创建服务请求
    const tx = await client.service.createRequest({
      level: 1,
      model: 'gpt-4',
      prompt: 'test',
    });
    await tx.waitForConfirmation();

    // 验证事件被捕获
    await new Promise(r => setTimeout(r, 1000));
    expect(events.length).toBeGreaterThan(0);
    expect(events[0].model).toBe('gpt-4');
  });
});
```

**API 集成测试:**
```typescript
// plugins/test/api-integration.test.ts
import request from 'supertest';
import { createApp } from '../src/app';

describe('API Integration', () => {
  const app = createApp();

  describe('POST /api/service/request', () => {
    it('should create service request and return escrow info', async () => {
      const response = await request(app)
        .post('/api/service/request')
        .set('Authorization', 'Bearer test-token')
        .send({
          level: 1,  // LLM 服务
          model: 'gpt-4',
          prompt: 'Hello, world!',
          maxTokens: 500,
        });

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('requestId');
      expect(response.body).toHaveProperty('escrowAddress');
    });
  });
});
```

---

## 4. E2E 测试

### 4.1 核心场景

**测试框架:** Playwright / Cypress

**测试环境:** Docker Compose (本地测试网)

#### 场景 1: 服务交易完整流程 (LLM 服务)

```typescript
// e2e/service-trading.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Service Trading Flow - LLM', () => {
  test('complete LLM service request flow', async ({ page }) => {
    // 1. 连接钱包
    await page.goto('/');
    await page.click('[data-testid="connect-wallet"]');
    await page.fill('[data-testid="mnemonic-input"]', process.env.TEST_MNEMONIC);
    await page.click('[data-testid="confirm-connect"]');

    // 2. 充值 STT
    await page.click('[data-testid="faucet-button"]');
    await expect(page.locator('[data-testid="balance"]')).toContainText('1000 STT');

    // 3. 注册为服务提供者 (用户A) - LLM Provider
    await page.click('[data-testid="become-provider"]');
    await page.selectOption('[data-testid="service-level"]', '1');  // Level 1: LLM
    await page.fill('[data-testid="api-key-input"]', 'sk-test-xxx');
    await page.fill('[data-testid="model-select"]', 'gpt-4');
    await page.click('[data-testid="register-provider"]');

    // 4. 发起服务请求 (用户B - 通过小灯界面)
    const context2 = await browser.newContext();
    const page2 = await context2.newPage();
    await page2.goto('/');
    await connectWallet(page2);

    await page2.click('[data-testid="service-market"]');
    await page2.click('[data-testid="service-llm-gpt-4"]');
    await page2.fill('[data-testid="prompt-input"]', 'Write a hello world program');
    await page2.click('[data-testid="submit-request"]');

    // 5. 验证请求被处理
    await expect(page2.locator('[data-testid="response"]')).toBeVisible({ timeout: 30000 });

    // 6. 验证支付完成
    await expect(page2.locator('[data-testid="tx-status"]')).toContainText('completed');
  });
});
```

#### 场景 2: 争议处理完整流程

```typescript
// e2e/dispute-resolution.spec.ts
test('complete dispute resolution flow', async ({ page, browser }) => {
  // 前置: 创建一笔有争议的交易
  const { buyerPage, sellerPage, requestId } = await setupDisputedTransaction(browser);

  // 1. 买方发起争议
  await buyerPage.goto(`/requests/${requestId}`);
  await buyerPage.click('[data-testid="raise-dispute"]');
  await buyerPage.fill('[data-testid="dispute-reason"]', 'Provider did not deliver as promised');
  await buyerPage.click('[data-testid="submit-dispute"]');

  // 2. 验证托管锁定
  await expect(buyerPage.locator('[data-testid="escrow-status"]')).toContainText('locked');

  // 3. 双方提交证据
  await buyerPage.click('[data-testid="add-evidence"]');
  await buyerPage.setInputFiles('[data-testid="evidence-file"]', './test-assets/evidence1.png');
  await buyerPage.click('[data-testid="submit-evidence"]');

  await sellerPage.goto(`/disputes/${requestId}`);
  await sellerPage.click('[data-testid="add-evidence"]');
  await sellerPage.setInputFiles('[data-testid="evidence-file"]', './test-assets/evidence2.png');
  await sellerPage.click('[data-testid="submit-evidence"]');

  // 4. 评审团投票 (模拟3个评审员)
  for (let i = 0; i < 3; i++) {
    const juryPage = await createJuryMemberPage(browser, i);
    await juryPage.goto(`/disputes/${requestId}`);
    await juryPage.click('[data-testid="vote-for-buyer"]');
    await juryPage.click('[data-testid="confirm-vote"]');
  }

  // 5. 验证裁决结果
  await buyerPage.goto(`/disputes/${requestId}/result`);
  await expect(buyerPage.locator('[data-testid="verdict"]')).toContainText('buyer wins');

  // 6. 验证德商变化
  const buyerMQ = await getMQValue(buyerPage);
  const sellerMQ = await getMQValue(sellerPage);
  expect(buyerMQ).toBeGreaterThan(100); // 买方德商增加
  expect(sellerMQ).toBeLessThan(100);   // 卖方德商减少
});
```

#### 场景 3: 想法众筹完整流程

```typescript
// e2e/idea-crowdfunding.spec.ts
test('idea crowdfunding flow', async ({ page, browser }) => {
  await connectWallet(page);

  // 1. 创建想法
  await page.goto('/ideas/new');
  await page.fill('[data-testid="idea-title"]', 'AI Code Reviewer');
  await page.fill('[data-testid="idea-description"]', 'An AI that reviews pull requests');
  await page.click('[data-testid="submit-idea"]');

  // 2. 等待 AI 评估
  await expect(page.locator('[data-testid="evaluation-result"]')).toBeVisible({ timeout: 60000 });
  const tokenEstimate = await page.locator('[data-testid="token-estimate"]').textContent();

  // 3. 支持想法
  const supporterPage = await createNewUserPage(browser);
  await supporterPage.goto('/ideas/ai-code-reviewer');
  await supporterPage.fill('[data-testid="support-amount"]', '100');
  await supporterPage.click('[data-testid="support-idea"]');

  // 4. 验证贡献权重
  await expect(supporterPage.locator('[data-testid="contribution-weight"]')).toBeVisible();

  // 5. 验证收益份额分配
  const creatorPage = page;
  await creatorPage.goto('/ideas/ai-code-reviewer/members');
  await expect(creatorPage.locator('[data-testid="revenue-share"]')).toContainText('90%'); // 创建者
});
```

#### 场景 4: 任务市场完整流程

```typescript
// e2e/task-market.spec.ts
test('task marketplace flow', async ({ page, browser }) => {
  await connectWallet(page);

  // 1. 发布任务
  await page.goto('/tasks/new');
  await page.fill('[data-testid="task-title"]', 'Design a logo');
  await page.fill('[data-testid="task-budget"]', '500');
  await page.fill('[data-testid="task-deadline"]', '2025-04-01');
  await page.click('[data-testid="publish-task"]');

  // 2. 托管预算
  await expect(page.locator('[data-testid="escrow-status"]')).toContainText('escrowed');

  // 3. 申请任务 (自由职业者)
  const freelancerPage = await createNewUserPage(browser);
  await freelancerPage.goto('/tasks/design-a-logo');
  await freelancerPage.click('[data-testid="apply-task"]');
  await freelancerPage.fill('[data-testid="proposal"]', 'I have 5 years of design experience');
  await freelancerPage.click('[data-testid="submit-proposal"]');

  // 4. 分配任务
  await page.goto('/tasks/design-a-logo/applicants');
  await page.click('[data-testid="assign-freelancer"]');

  // 5. 提交里程碑
  await freelancerPage.goto('/tasks/design-a-logo');
  await freelancerPage.click('[data-testid="submit-milestone"]');
  await freelancerPage.fill('[data-testid="milestone-description"]', 'Initial design concepts');
  await freelancerPage.setInputFiles('[data-testid="deliverable"]', './test-assets/design.png');
  await freelancerPage.click('[data-testid="submit"]');

  // 6. 验收里程碑
  await page.goto('/tasks/design-a-logo');
  await page.click('[data-testid="approve-milestone"]');

  // 7. 验证支付释放
  await expect(page.locator('[data-testid="payment-status"]')).toContainText('released');
});
```

#### 场景 5: Agent 服务完整流程 (Level 2)

```typescript
// e2e/agent-service.spec.ts
test('Agent service flow - Coder Agent', async ({ page, browser }) => {
  await connectWallet(page);

  // 1. 注册 Agent 服务提供者
  const providerPage = await createNewUserPage(browser);
  await providerPage.goto('/');
  await connectWallet(providerPage);
  await providerPage.click('[data-testid="become-provider"]');
  await providerPage.selectOption('[data-testid="service-level"]', '2');  // Level 2: Agent
  await providerPage.selectOption('[data-testid="agent-type"]', 'coder');
  await providerPage.click('[data-testid="register-provider"]');

  // 2. 通过小灯发起 Agent 任务
  await page.goto('/');
  await page.fill('[data-testid="chat-input"]', 'Write a Python function to sort a list');
  await page.click('[data-testid="send-message"]');

  // 3. 验证意图识别正确路由到 Agent 服务
  await expect(page.locator('[data-testid="service-level"]')).toContainText('Agent');

  // 4. 等待 Agent 执行完成
  await expect(page.locator('[data-testid="agent-result"]')).toBeVisible({ timeout: 60000 });

  // 5. 验证交付物
  const codeArtifact = await page.locator('[data-testid="code-artifact"]').textContent();
  expect(codeArtifact).toContain('def sort');

  // 6. 验证支付完成
  await expect(page.locator('[data-testid="tx-status"]')).toContainText('completed');
});
```

#### 场景 6: Workflow 服务完整流程 (Level 3)

```typescript
// e2e/workflow-service.spec.ts
test('Workflow service flow - Software Development', async ({ page, browser }) => {
  await connectWallet(page);

  // 1. 注册 Workflow 服务提供者
  const providerPage = await createNewUserPage(browser);
  await providerPage.goto('/');
  await connectWallet(providerPage);
  await providerPage.click('[data-testid="become-provider"]');
  await providerPage.selectOption('[data-testid="service-level"]', '3');  // Level 3: Workflow
  await providerPage.selectOption('[data-testid="workflow-type"]', 'software');
  await providerPage.click('[data-testid="register-provider"]');

  // 2. 通过小灯发起 Workflow
  await page.goto('/');
  await page.fill('[data-testid="chat-input"]', 'I want to build a todo app with React');
  await page.click('[data-testid="send-message"]');

  // 3. 验证意图识别正确路由到 Workflow 服务
  await expect(page.locator('[data-testid="service-level"]')).toContainText('Workflow');

  // 4. 等待 Workflow 启动
  await expect(page.locator('[data-testid="workflow-status"]')).toContainText('running');

  // 5. 验证当前步骤
  await expect(page.locator('[data-testid="current-step"]')).toContainText('requirement_analysis');

  // 6. 模拟人工审核步骤
  await advanceWorkflowToHumanGate(page, 'code_review');
  await page.click('[data-testid="approve-step"]');

  // 7. 等待 Workflow 完成
  await expect(page.locator('[data-testid="workflow-status"]')).toContainText('completed', { timeout: 120000 });

  // 8. 验证交付物
  const deliverables = await page.locator('[data-testid="deliverables"]').all();
  expect(deliverables.length).toBeGreaterThan(0);
});
```

### 4.2 边界场景

```typescript
// e2e/edge-cases.spec.ts

test.describe('Edge Cases', () => {
  test('concurrent requests to same service provider', async ({ browser }) => {
    const provider = await createServiceProviderPage(browser, { level: 1 });
    const pages = await Promise.all([
      createConsumerPage(browser),
      createConsumerPage(browser),
      createConsumerPage(browser),
    ]);

    // 并发发送服务请求
    const requests = pages.map(p =>
      p.evaluate(() => submitServiceRequest({ level: 1, model: 'gpt-4', prompt: 'test' }))
    );

    const results = await Promise.all(requests);

    // 所有请求都应成功排队
    results.forEach(r => expect(r.status).toBe('queued'));
  });

  test('timeout handling for service request', async ({ page }) => {
    await connectWallet(page);
    await submitServiceRequest(page, { level: 1, model: 'gpt-4', prompt: 'test', timeout: 5000 });

    // 验证超时后自动退款
    await expect(page.locator('[data-testid="refund-status"]')).toContainText('refunded', { timeout: 10000 });
  });

  test('error recovery from network interruption', async ({ page, context }) => {
    await connectWallet(page);
    const requestId = await submitServiceRequest(page, { level: 1, model: 'gpt-4', prompt: 'test' });

    // 模拟网络中断
    await context.setOffline(true);
    await page.waitForTimeout(2000);
    await context.setOffline(false);

    // 验证自动重连和状态恢复
    await page.reload();
    await expect(page.locator(`[data-testid="request-${requestId}"]`)).toBeVisible();
  });

  test('insufficient balance handling', async ({ page }) => {
    await connectWallet(page, { balance: 10 }); // 低余额

    await page.goto('/services');
    await page.fill('[data-testid="tokens-input"]', '1000');
    await page.click('[data-testid="submit-request"]');

    await expect(page.locator('[data-testid="error-message"]')).toContainText('Insufficient balance');
  });
});
```

---

## 5. 性能测试

### 5.1 链性能

**工具:** Cosmos SDK 内置 benchmark + 自定义脚本

**测试指标:**

| 指标 | 目标值 | 测试方法 |
|------|--------|----------|
| TPS (Transactions Per Second) | > 1000 | 压力测试 |
| 交易延迟 (Latency) | < 5s (P99) | 延迟监控 |
| 区块时间 | ~1s | 节点日志 |
| 最终确认时间 | < 10s | 端到端测量 |

**TPS 测试脚本:**
```bash
#!/bin/bash
# scripts/benchmark-tps.sh

# 参数
DURATION=60        # 测试持续时间 (秒)
CONCURRENCY=100    # 并发数
ENDPOINT="http://localhost:26657"

# 生成测试交易
echo "Starting TPS benchmark..."
for i in $(seq 1 $CONCURRENCY); do
  (
    start_time=$(date +%s)
    while [ $(($(date +%s) - start_time)) -lt $DURATION ]; do
      # 发送测试交易
      shared tx send $SENDER $RECEIVER 1stt --fees 10stt -y
    done
  ) &
done

# 等待完成
wait

# 分析结果
echo "Analyzing results..."
curl -s "$ENDPOINT/status" | jq '.result.sync_info'
```

**Go benchmark 示例:**
```go
// x/service/keeper/benchmark_test.go
func BenchmarkServiceRequest(b *testing.B) {
    app := setupBenchmarkApp()
    ctx := app.NewContext(false)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        requester := fmt.Sprintf("cosmos1requester%d", i)
        app.ServiceKeeper.CreateRequest(ctx, requester, &types.ServiceRequest{
            Level:     1,
            Model:     "gpt-4",
            MaxTokens: 100,
        })
    }
}
```

### 5.2 服务性能

**工具:** k6 / Artillery / Locust

**API 性能目标:**

| 端点 | 目标响应时间 | 目标吞吐量 |
|------|-------------|------------|
| GET /api/balance | < 50ms | 5000 req/s |
| POST /api/service/request | < 200ms | 1000 req/s |
| GET /api/services | < 100ms | 2000 req/s |
| WS /api/chat | < 50ms (首字节) | 1000 conn |

**k6 测试脚本:**
```javascript
// scripts/load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 },   // 预热
    { duration: '1m', target: 1000 },   // 峰值
    { duration: '30s', target: 100 },   // 降温
  ],
  thresholds: {
    http_req_duration: ['p(99)<200'], // 99% 请求 < 200ms
    http_req_failed: ['rate<0.01'],   // 错误率 < 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // 测试获取余额
  const balanceRes = http.get(`${BASE_URL}/api/balance`, {
    headers: { Authorization: `Bearer ${__ENV.TEST_TOKEN}` },
  });
  check(balanceRes, {
    'balance status 200': (r) => r.status === 200,
    'balance < 50ms': (r) => r.timings.duration < 50,
  });

  sleep(1);

  // 测试创建服务请求
  const createRes = http.post(
    `${BASE_URL}/api/service/request`,
    JSON.stringify({ level: 1, model: 'gpt-4', prompt: 'test' }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(createRes, {
    'create status 200': (r) => r.status === 200,
    'create < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(1);
}
```

**WebSocket 负载测试:**
```javascript
// scripts/ws-load-test.js
import { check } from 'k6';
import { WebSocket } from 'k6/experimental/websockets';

export const options = {
  vus: 1000,
  duration: '1m',
};

export default function () {
  const url = 'ws://localhost:8080/api/chat';
  const ws = new WebSocket(url);

  ws.on('open', () => {
    ws.send(JSON.stringify({ type: 'message', content: 'Hello' }));
  });

  ws.on('message', (msg) => {
    check(msg, {
      'received response': (m) => JSON.parse(m).type === 'response',
    });
    ws.close();
  });
}
```

### 5.3 数据库性能

```sql
-- PostgreSQL 查询性能测试
EXPLAIN ANALYZE SELECT * FROM ideas WHERE status = 'active' ORDER BY created_at DESC LIMIT 20;

-- 索引优化
CREATE INDEX idx_ideas_status_created ON ideas(status, created_at DESC);
```

---

## 6. 安全测试

### 6.1 智能合约审计清单

**Cosmos SDK 模块审计要点:**

| 类别 | 检查项 | 风险等级 |
|------|--------|----------|
| **输入验证** | 所有消息字段验证 | 高 |
| **权限控制** | 只有授权者可执行操作 | 高 |
| **重入攻击** | 状态更新顺序正确 | 高 |
| **整数溢出** | 使用 safe math | 中 |
| **状态一致性** | 跨模块调用正确 | 高 |
| **Gas 消耗** | 避免无限循环 | 中 |

**审计清单:**
```markdown
## x/identity 审计清单

- [ ] RegisterIdentity
  - [ ] 验证 handle 格式
  - [ ] 检查重复注册
  - [ ] 验证签名正确性
  - [ ] 检查 Merkle proof 验证逻辑

- [ ] RevokeIdentity
  - [ ] 只有 owner 可撤销
  - [ ] 撤销后状态正确
  - [ ] 事件正确发出

## x/service 审计清单 (服务市场核心模块)

- [ ] RegisterService
  - [ ] 服务级别验证
  - [ ] 定价模型验证
  - [ ] 防止重复注册

- [ ] CreateRequest
  - [ ] 余额检查
  - [ ] 托管创建正确
  - [ ] 事件发出

- [ ] ExecuteRequest
  - [ ] 只有分配的 provider 可执行
  - [ ] 服务使用量验证
  - [ ] 支付计算正确

## x/escrow 审计清单 (托管支付核心模块)

- [ ] CreateEscrow
  - [ ] 资金锁定正确
  - [ ] 状态初始化正确

- [ ] ReleaseEscrow
  - [ ] 只有授权方可释放
  - [ ] 金额分配正确

- [ ] RefundEscrow
  - [ ] 退款条件验证
  - [ ] 争议状态检查
```

**自动化安全扫描:**
```yaml
# .github/workflows/security.yml
name: Security Scan

on: [push, pull_request]

jobs:
  gosec:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: ./...

  semgrep:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: returntocorp/semgrep-action@v1
        with:
          config: >-
            p/security-audit
            p/secrets
            p/golang
```

### 6.2 渗透测试

**API 安全测试:**

| 测试项 | 描述 | 工具 |
|--------|------|------|
| SQL 注入 | 验证输入过滤 | SQLMap |
| XSS | 检查输出转义 | OWASP ZAP |
| CSRF | 验证 token 机制 | Burp Suite |
| 认证绕过 | 测试权限边界 | 自定义脚本 |
| 速率限制 | 验证 API 限流 | k6 |

**渗透测试脚本:**
```typescript
// security/api-pentest.ts
describe('API Security Tests', () => {
  it('should prevent SQL injection', async () => {
    const maliciousInput = "'; DROP TABLE users; --";
    const response = await request(app)
      .get(`/api/users?handle=${maliciousInput}`);

    expect(response.status).not.toBe(500);
    expect(response.body).not.toContain('error');
  });

  it('should require authentication for protected routes', async () => {
    const response = await request(app)
      .post('/api/service/request')
      .send({ level: 1, model: 'gpt-4', prompt: 'test' });

    expect(response.status).toBe(401);
  });

  it('should validate JWT token expiration', async () => {
    const expiredToken = generateExpiredToken();
    const response = await request(app)
      .get('/api/balance')
      .set('Authorization', `Bearer ${expiredToken}`);

    expect(response.status).toBe(401);
    expect(response.body.error).toContain('expired');
  });

  it('should enforce rate limiting', async () => {
    const requests = Array(100).fill(null).map(() =>
      request(app).get('/api/balance')
    );

    const responses = await Promise.all(requests);
    const rateLimited = responses.filter(r => r.status === 429);

    expect(rateLimited.length).toBeGreaterThan(0);
  });
});
```

**DDoS 测试:**
```bash
# 使用 hping3 进行 DDoS 测试 (仅测试环境)
# 注意: 需要在隔离环境中进行
hping3 -S -p 8080 --flood --rand-source target-server
```

---

## 7. 测试环境

### 7.1 本地开发

**Ignite 链启动:**
```bash
# 启动本地开发链
ignite chain serve -r

# 运行单元测试
ignite chain test
```

**Docker Compose 配置:**
```yaml
# docker-compose.test.yml
version: '3.8'

services:
  # 核心模块 - 共识链节点
  shared-node:
    build:
      context: ./chain
      dockerfile: Dockerfile
    ports:
      - "26657:26657"  # RPC
      - "26656:26656"  # P2P
    environment:
      - CHAIN_ID=shared-test

  # PostgreSQL (插件服务)
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: sharetokens_test
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
    ports:
      - "5432:5432"

  # Redis (缓存)
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  # 服务提供者插件 - LLM Provider
  llm-provider-plugin:
    build:
      context: ./plugins/llm-provider
    environment:
      - DATABASE_URL=postgresql://test:test@postgres:5432/sharetokens_test
      - CHAIN_RPC=http://shared-node:26657
    depends_on:
      - postgres
      - shared-node

  # 服务提供者插件 - Agent Provider (OpenFang)
  agent-provider-plugin:
    build:
      context: ./plugins/agent-provider
    environment:
      - OPENFANG_DASHBOARD=http://localhost:4200
      - CHAIN_RPC=http://shared-node:26657
    depends_on:
      - shared-node

  # 服务提供者插件 - Workflow Provider
  workflow-provider-plugin:
    build:
      context: ./plugins/workflow-provider
    environment:
      - DATABASE_URL=postgresql://test:test@postgres:5432/sharetokens_test
      - CHAIN_RPC=http://shared-node:26657
    depends_on:
      - postgres
      - shared-node

  # 用户插件 - 小灯界面
  xiaodeng-client:
    build:
      context: ./plugins/xiaodeng-client
    ports:
      - "3000:3000"
    environment:
      - CHAIN_RPC=http://shared-node:26657
      - MARKETPLACE_URL=http://marketplace:8080
    depends_on:
      - shared-node

  # 服务市场 API
  marketplace:
    build:
      context: ./modules/marketplace
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://test:test@postgres:5432/sharetokens_test
      - CHAIN_RPC=http://shared-node:26657
    depends_on:
      - postgres
      - shared-node

volumes:
  postgres_data:
```

**启动测试环境:**
```bash
# 启动所有服务
docker-compose -f docker-compose.test.yml up -d

# 等待服务就绪
./scripts/wait-for-services.sh

# 运行测试
npm run test:e2e

# 清理
docker-compose -f docker-compose.test.yml down -v
```

### 7.2 测试网

**多节点网络配置:**
```yaml
# testnet/config.yaml
chain_id: shared-testnet
validators:
  - name: validator-1
    ip: 10.0.0.1
    stake: 1000000stt
  - name: validator-2
    ip: 10.0.0.2
    stake: 1000000stt
  - name: validator-3
    ip: 10.0.0.3
    stake: 1000000stt
  - name: validator-4
    ip: 10.0.0.4
    stake: 1000000stt

genesis:
  app_state:
    auth:
      accounts:
        - name: faucet
          coins: ["1000000000stt"]
```

**测试数据管理:**
```typescript
// scripts/seed-test-data.ts
import { TestClient } from './test-utils';

async function seedTestData() {
  const client = await TestClient.create();

  // 创建测试用户
  const users = await Promise.all([
    client.createAccount('alice', { balance: 10000 }),
    client.createAccount('bob', { balance: 10000 }),
    client.createAccount('charlie', { balance: 10000 }),
  ]);

  // 注册测试服务 (三层服务)
  // Level 1: LLM 服务
  await client.registerService({
    provider: users[0].address,
    level: 1,
    name: 'GPT-4 Test Service',
    category: 'llm_chat',
    pricing: { type: 'per_token', pricePerToken: 100 },
  });

  // Level 2: Agent 服务
  await client.registerService({
    provider: users[1].address,
    level: 2,
    name: 'Coder Agent Test',
    category: 'agent_coder',
    capabilities: ['code_generation', 'debugging'],
    pricing: { type: 'per_task', pricePerTask: 1000 },
  });

  // Level 3: Workflow 服务
  await client.registerService({
    provider: users[2].address,
    level: 3,
    name: 'Software Development Workflow',
    category: 'workflow_software',
    pricing: { type: 'package', totalPrice: 10000 },
  });

  // 创建测试想法
  await client.createIdea({
    title: 'Test Idea 1',
    creator: users[0].address,
  });

  // 创建测试任务
  await client.createTask({
    title: 'Test Task 1',
    budget: 100,
    creator: users[1].address,
  });

  console.log('Test data seeded successfully');
}

seedTestData();
```

---

## 8. CI/CD 集成

### 8.1 GitHub Actions 工作流

```yaml
# .github/workflows/test.yml
name: Test Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.21'
  NODE_VERSION: '20'

jobs:
  # 单元测试
  unit-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Run Go unit tests
        working-directory: ./chain
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Run TypeScript unit tests
        run: |
          npm ci
          npm run test:unit -- --coverage

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./chain/coverage.out, ./coverage/lcov.info

  # 集成测试
  integration-test:
    runs-on: ubuntu-latest
    needs: unit-test
    steps:
      - uses: actions/checkout@v4

      - name: Start test environment
        run: docker-compose -f docker-compose.test.yml up -d

      - name: Wait for services
        run: ./scripts/wait-for-services.sh

      - name: Run integration tests
        run: npm run test:integration

      - name: Cleanup
        if: always()
        run: docker-compose -f docker-compose.test.yml down -v

  # E2E 测试
  e2e-test:
    runs-on: ubuntu-latest
    needs: integration-test
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install Playwright
        run: |
          npm ci
          npx playwright install --with-deps

      - name: Start test environment
        run: docker-compose -f docker-compose.test.yml up -d

      - name: Run E2E tests
        run: npm run test:e2e

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-report
          path: playwright-report/

  # 性能测试
  performance-test:
    runs-on: ubuntu-latest
    needs: integration-test
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4

      - name: Run k6 load test
        uses: grafana/k6-action@v0.3.0
        with:
          filename: scripts/load-test.js
        env:
          K6_CLOUD_TOKEN: ${{ secrets.K6_CLOUD_TOKEN }}
          BASE_URL: http://localhost:8080

  # 安全扫描
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run GoSec
        uses: securego/gosec@master
        with:
          args: ./chain/...

      - name: Run npm audit
        run: npm audit --audit-level=high

      - name: Run Snyk
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
```

### 8.2 测试命令

```json
// package.json
{
  "scripts": {
    "test": "npm run test:unit && npm run test:integration",
    "test:unit": "vitest run --coverage",
    "test:unit:watch": "vitest watch",
    "test:integration": "vitest run --config vitest.integration.config.ts",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:coverage": "vitest run --coverage --reporter=json --outputFile=coverage.json",
    "test:benchmark": "vitest bench"
  }
}
```

---

## 9. 测试报告

### 9.1 覆盖率报告

```bash
# 生成覆盖率报告
npm run test:coverage

# Go 覆盖率
cd chain && go tool cover -html=coverage.out -o coverage.html
```

### 9.2 测试仪表板

| 指标 | 目标 | 当前状态 |
|------|------|----------|
| 单元测试覆盖率 | 80% | - |
| 集成测试覆盖率 | 60% | - |
| E2E 测试通过率 | 100% | - |
| 性能测试 TPS | > 1000 | - |
| 安全漏洞 | 0 高危 | - |

---

## 10. 附录

### 10.1 测试工具清单

#### 核心模块测试工具

| 工具 | 用途 | 链接 |
|------|------|------|
| Go testing | Go 单元测试 (核心模块) | 内置 |
| testify | Go 断言库 | github.com/stretchr/testify |
| Cosmos SDK testutil | Cosmos SDK 模块测试 | docs.cosmos.network |

#### 插件模块测试工具

| 工具 | 用途 | 链接 |
|------|------|------|
| Jest | TypeScript 单元测试 (插件) | jestjs.io |
| ts-jest | TypeScript Jest 支持 | kulshekhar.github.io/ts-jest |
| Vitest | 前端单元测试 | vitest.dev |
| React Testing Library | React 组件测试 | testing-library.com |

#### E2E 和性能测试工具

| 工具 | 用途 | 链接 |
|------|------|------|
| Playwright | E2E 测试 | playwright.dev |
| k6 | 性能测试 | k6.io |
| Artillery | API 负载测试 | artillery.io |

#### 安全测试工具

| 工具 | 用途 | 链接 |
|------|------|------|
| OWASP ZAP | 安全扫描 | zaproxy.org |
| SQLMap | SQL 注入检测 | sqlmap.org |
| GoSec | Go 安全扫描 | securego.io |
| Semgrep | 代码安全分析 | semgrep.dev |

#### OpenFang 相关工具

| 工具 | 用途 | 链接 |
|------|------|------|
| OpenFang Dashboard | Agent 监控和调试 | localhost:4200 |
| OpenFang CLI | Agent 执行和管理 | OpenFang CLI |

### 10.2 参考资料

- [Cosmos SDK Testing Guide](https://docs.cosmos.network/main/building-modules/testing)
- [Testing Blockchain Applications](https://ethereum.org/en/developers/docs/smart-contracts/testing/)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)
- [OpenFang Documentation](https://openfang.ai/docs)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)

---

*Test Strategy defined: 2025-03-02*
*Last updated: 2025-03-02*
