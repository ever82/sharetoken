# Service Marketplace - 服务市场

> **模块类型:** 核心模块
> **技术栈:** TypeScript + Cosmos SDK
> **位置:** `src/modules/marketplace/`

---

## 概述

服务市场是 ShareTokens 的核心模块，提供 AI 服务的交易基础设施。它实现三层服务结构（LLM / Agent / Workflow），支持服务的注册、发现、定价和路由。

服务市场本身不直接提供 AI 能力，而是作为连接服务提供者和服务消费者的中介层。

---

## 三层服务结构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Service Marketplace (核心模块)                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Level 3: Workflow 服务                            │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  输入: 想法/需求 (自然语言描述)                               │    │   │
│  │  │  输出: 完整交付物 (代码仓库、文档、部署应用等)                │    │   │
│  │  │  定价: 按流程打包计费 (固定价格或里程碑计费)                  │    │   │
│  │  │  示例: 软件开发、内容创作、商业策划、生活服务                 │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Level 2: Agent 服务                               │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  输入: 任务描述 (结构化或半结构化)                            │    │   │
│  │  │  输出: 执行结果 (代码、报告、分析等)                          │    │   │
│  │  │  能力: tools + skills (工具调用 + 技能模块)                   │    │   │
│  │  │  定价: 按能力/复杂度计费                                      │    │   │
│  │  │  示例: Coder Agent, Researcher Agent, Writer Agent           │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Level 1: LLM API 服务                             │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  输入: prompt (文本提示)                                     │    │   │
│  │  │  输出: completion (文本补全)                                 │    │   │
│  │  │  定价: 按 Token 计费                                         │    │   │
│  │  │  示例: GPT-4, Claude, Llama, Qwen 等                         │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ LLM Provider    │ │ Agent Provider  │ │ Workflow        │
│ Plugin          │ │ Plugin          │ │ Provider Plugin │
│ (托管API Key)   │ │ (OpenFang)      │ │                 │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

---

## 核心功能

### 1. 服务注册

服务提供者在市场上注册其服务，声明能力、定价和可用性。

```typescript
// src/modules/marketplace/types.ts

interface ServiceRegistration {
  // 服务标识
  id: string
  provider: Address
  level: ServiceLevel  // 1=LLM, 2=Agent, 3=Workflow

  // 服务信息
  name: string
  description: string
  category: ServiceCategory

  // 能力声明
  capabilities: string[]
  tools?: string[]        // Level 2+: 支持的工具
  skills?: string[]       // Level 2+: 支持的技能
  workflow?: WorkflowDef  // Level 3: 工作流定义

  // 定价
  pricing: PricingModel

  // SLA
  sla: ServiceLevelAgreement

  // 状态
  status: ServiceStatus
  registeredAt: Timestamp
  updatedAt: Timestamp
}

type ServiceLevel = 1 | 2 | 3

type ServiceCategory =
  | 'llm_chat'           // LLM 对话
  | 'llm_completion'     // LLM 补全
  | 'llm_embedding'      // LLM 嵌入
  | 'agent_coder'        // 编程 Agent
  | 'agent_researcher'   // 研究 Agent
  | 'agent_writer'       // 写作 Agent
  | 'workflow_software'  // 软件开发流程
  | 'workflow_content'   // 内容创作流程
  | 'workflow_business'  // 商业策划流程
  | 'workflow_service'   // 生活服务流程

type PricingModel =
  | { type: 'per_token'; pricePerToken: TokenAmount }
  | { type: 'per_task'; pricePerTask: TokenAmount }
  | { type: 'complexity'; basePrice: TokenAmount; complexityMultiplier: Record<string, number> }
  | { type: 'package'; totalPrice: TokenAmount; deliverables: string[] }
  | { type: 'milestone'; milestones: { name: string; price: TokenAmount }[] }

interface ServiceLevelAgreement {
  responseTime: Duration      // 响应时间承诺
  uptime: number              // 可用性百分比
  refundPolicy: RefundPolicy
  supportLevel: 'basic' | 'standard' | 'premium'
}
```

### 2. 服务发现

用户可以搜索、筛选和比较市场上的服务。

```typescript
// src/modules/marketplace/discovery.ts

export class ServiceDiscovery {
  // 搜索服务
  async searchServices(query: ServiceQuery): Promise<ServiceResult[]>

  // 按类别浏览
  async listByCategory(category: ServiceCategory): Promise<ServiceResult[]>

  // 按能力查找
  async findByCapability(capability: string): Promise<ServiceResult[]>

  // 推荐服务
  async recommendServices(userNeed: UserNeed): Promise<ServiceRecommendation[]>

  // 比较服务
  async compareServices(serviceIds: string[]): Promise<ServiceComparison>
}

interface ServiceQuery {
  keyword?: string
  level?: ServiceLevel
  category?: ServiceCategory
  capabilities?: string[]
  maxPrice?: TokenAmount
  minRating?: number
  provider?: Address
  sortBy?: 'price' | 'rating' | 'popularity' | 'response_time'
  sortOrder?: 'asc' | 'desc'
}

interface ServiceResult {
  service: ServiceRegistration
  rating: number
  reviewCount: number
  completedTasks: number
  avgResponseTime: Duration
  priceExample: string
}
```

### 3. 服务定价

市场提供透明的定价机制和费用结算。

```typescript
// src/modules/marketplace/pricing.ts

export class PricingEngine {
  // 计算服务费用
  async calculatePrice(
    serviceId: string,
    request: ServiceRequest
  ): Promise<PriceEstimate>

  // 创建支付托管
  async createEscrow(
    serviceId: string,
    request: ServiceRequest,
    price: TokenAmount
  ): Promise<EscrowInfo>

  // 释放支付
  async releasePayment(escrowId: string, result: ServiceResult): Promise<void>

  // 退款处理
  async processRefund(escrowId: string, reason: RefundReason): Promise<void>
}

interface PriceEstimate {
  basePrice: TokenAmount
  platformFee: TokenAmount      // 平台手续费
  totalPrice: TokenAmount
  breakdown: PriceBreakdown
  validUntil: Timestamp
}

interface PriceBreakdown {
  items: {
    name: string
    amount: TokenAmount
    description?: string
  }[]
}
```

### 4. 请求路由

市场将用户请求路由到合适的服务提供者。

```typescript
// src/modules/marketplace/router.ts

export class RequestRouter {
  // 路由请求
  async routeRequest(request: ServiceRequest): Promise<RouteResult>

  // 选择提供者
  async selectProvider(
    request: ServiceRequest,
    candidates: Address[]
  ): Promise<Address>

  // 负载均衡
  async balanceLoad(providers: ProviderLoad[]): Promise<Address>

  // 故障转移
  async failover(request: ServiceRequest, failedProvider: Address): Promise<RouteResult>
}

interface ServiceRequest {
  id: string
  consumer: Address
  level: ServiceLevel

  // Level 1: LLM 请求
  prompt?: string
  model?: string
  maxTokens?: number

  // Level 2: Agent 请求
  taskDescription?: string
  requiredTools?: string[]
  requiredSkills?: string[]

  // Level 3: Workflow 请求
  idea?: string
  deliverables?: string[]
  milestones?: string[]

  // 通用参数
  budget?: TokenAmount
  deadline?: Timestamp
  preferences?: Record<string, any>
}

interface RouteResult {
  provider: Address
  serviceId: string
  estimatedTime: Duration
  estimatedPrice: TokenAmount
  escrowAddress: Address
}
```

### 5. 标准化接口

市场定义统一的服务接口规范，确保互操作性。

```typescript
// src/modules/marketplace/interface.ts

// Level 1: LLM 接口
interface LLMApiInterface {
  complete(request: LLMRequest): Promise<LLMResponse>
  stream(request: LLMRequest): AsyncIterable<LLMChunk>
  embed(request: EmbedRequest): Promise<EmbedResponse>
}

interface LLMRequest {
  prompt: string
  model: string
  maxTokens?: number
  temperature?: number
  stopSequences?: string[]
}

interface LLMResponse {
  completion: string
  usage: TokenUsage
  finishReason: 'stop' | 'length' | 'error'
}

// Level 2: Agent 接口
interface AgentApiInterface {
  execute(request: AgentRequest): Promise<AgentResponse>
  getStatus(executionId: string): Promise<AgentStatus>
  cancel(executionId: string): Promise<void>
}

interface AgentRequest {
  task: string
  context?: Record<string, any>
  tools?: string[]
  constraints?: AgentConstraints
}

interface AgentResponse {
  result: any
  steps: AgentStep[]
  usage: ResourceUsage
  artifacts?: Artifact[]
}

// Level 3: Workflow 接口
interface WorkflowApiInterface {
  start(request: WorkflowRequest): Promise<WorkflowExecution>
  getStatus(executionId: string): Promise<WorkflowStatus>
  provideInput(executionId: string, input: HumanInput): Promise<void>
  approve(executionId: string, stepId: string): Promise<void>
  cancel(executionId: string): Promise<void>
}

interface WorkflowRequest {
  idea: string
  requirements?: string[]
  deliverables?: string[]
  preferences?: Record<string, any>
}

interface WorkflowExecution {
  id: string
  status: 'pending' | 'running' | 'waiting_input' | 'completed' | 'failed'
  currentStep?: string
  progress: number
  deliverables: Deliverable[]
  startedAt: Timestamp
  estimatedCompletion?: Timestamp
}
```

---

## 服务提供者插件

服务提供者通过插件接入市场，实现具体的服务能力。

### LLM Provider Plugin

托管 API Key，提供 Level 1 LLM 服务。

```typescript
// src/plugins/llm-provider/index.ts

export class LLMProviderPlugin {
  private apiKeys: Map<string, EncryptedKey>
  private rateLimiter: RateLimiter

  // API Key 管理
  async registerApiKey(provider: string, apiKey: string): Promise<void>
  async rotateApiKey(provider: string): Promise<void>
  async revokeApiKey(provider: string): Promise<void>

  // 请求代理
  async proxyRequest(request: LLMRequest): Promise<LLMResponse>

  // 使用量统计
  async getUsage(provider: string, period: TimePeriod): Promise<UsageStats>
}

// 支持的 LLM 提供者
type LLMProvider =
  | 'openai'      // GPT-4, GPT-3.5
  | 'anthropic'   // Claude
  | 'google'      // Gemini
  | 'meta'        // Llama
  | 'alibaba'     // Qwen
  | 'deepseek'    // DeepSeek
  | 'local'       // 本地部署模型
```

### Agent Provider Plugin (OpenFang)

提供 Level 2 Agent 服务。

```typescript
// src/plugins/agent-provider/index.ts

export class AgentProviderPlugin {
  private openfangClient: OpenFangClient

  // Agent 管理
  async registerAgent(definition: AgentDefinition): Promise<string>
  async updateAgent(agentId: string, updates: Partial<AgentDefinition>): Promise<void>
  async unregisterAgent(agentId: string): Promise<void>

  // 任务执行
  async executeTask(agentId: string, request: AgentRequest): Promise<AgentResponse>

  // 状态监控
  async getAgentStatus(agentId: string): Promise<AgentStatus>
  async getActiveExecutions(agentId: string): Promise<Execution[]>
}

interface AgentDefinition {
  name: string
  description: string
  capabilities: string[]
  tools: ToolDefinition[]
  skills: SkillDefinition[]
  llmConfig: LLMConfig
  pricing: AgentPricing
}

// 预置 Agent 类型
const BUILTIN_AGENTS = {
  coder: {
    name: 'Coder Agent',
    capabilities: ['code_generation', 'code_review', 'debugging'],
    tools: ['github', 'terminal', 'file_system'],
  },
  researcher: {
    name: 'Researcher Agent',
    capabilities: ['web_search', 'data_analysis', 'report_writing'],
    tools: ['web_browser', 'calculator', 'document_editor'],
  },
  writer: {
    name: 'Writer Agent',
    capabilities: ['content_creation', 'editing', 'translation'],
    tools: ['document_editor', 'grammar_checker'],
  },
}
```

### Workflow Provider Plugin

提供 Level 3 Workflow 服务。

```typescript
// src/plugins/workflow-provider/index.ts

export class WorkflowProviderPlugin {
  private workflowEngine: WorkflowEngine

  // Workflow 注册
  async registerWorkflow(definition: WorkflowDefinition): Promise<string>
  async updateWorkflow(workflowId: string, updates: Partial<WorkflowDefinition>): Promise<void>

  // 执行管理
  async startExecution(workflowId: string, request: WorkflowRequest): Promise<WorkflowExecution>
  async getExecutionStatus(executionId: string): Promise<WorkflowStatus>
  async provideHumanInput(executionId: string, input: HumanInput): Promise<void>

  // 交付物管理
  async getDeliverables(executionId: string): Promise<Deliverable[]>
  async verifyDeliverable(deliverableId: string): Promise<VerificationResult>
}

interface WorkflowDefinition {
  id: string
  name: string
  description: string
  category: WorkflowCategory

  // 流程定义
  steps: WorkflowStep[]
  transitions: Transition[]

  // 人工介入点
  humanGates: HumanGate[]

  // 交付物定义
  deliverables: DeliverableSpec[]

  // 定价
  pricing: WorkflowPricing
}

// 预置 Workflow 类型
const BUILTIN_WORKFLOWS = {
  software: {
    name: '软件开发',
    steps: [
      '需求分析', '架构设计', '代码实现',
      '测试验证', '部署发布', '文档交付'
    ],
    deliverables: ['代码仓库', '测试报告', '部署应用', '技术文档'],
  },
  content: {
    name: '内容创作',
    steps: [
      '主题研究', '大纲设计', '内容撰写',
      '审核修改', '格式排版', '发布交付'
    ],
    deliverables: ['内容文档', '配图素材', '发布链接'],
  },
  business: {
    name: '商业策划',
    steps: [
      '市场调研', '竞品分析', '策略制定',
      '方案设计', '财务预测', '报告撰写'
    ],
    deliverables: ['商业计划书', '财务模型', '演示文稿'],
  },
}
```

---

## 用户客户端：GenieBot 界面

GenieBot 是服务市场的客户端界面，帮助用户便捷地使用各种服务。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            GenieBot 界面 (客户端插件)                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          对话界面                                    │   │
│  │  - 自然语言交互                                                      │   │
│  │  - 意图识别和需求理解                                                │   │
│  │  - 服务推荐和选择                                                    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │  LLM 对话   │ │ Agent 任务  │ │ Workflow    │ │  服务管理   │          │
│  │  (Level 1)  │ │  (Level 2)  │ │  (Level 3)  │ │             │          │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           服务市场 (核心模块)                                 │
│  - 服务注册和发现                                                           │
│  - 请求路由和定价                                                           │
│  - 托管支付和结算                                                           │
└─────────────────────────────────────────────────────────────────────────────┘
```

```typescript
// src/plugins/geniebot-client/index.ts

export class GenieBotClient {
  private marketplace: ServiceMarketplaceClient

  // 对话入口
  async chat(message: string): Promise<ChatResponse> {
    // 1. 理解用户意图
    const intent = await this.analyzeIntent(message)

    // 2. 根据意图选择服务层级
    const serviceLevel = this.determineServiceLevel(intent)

    // 3. 查找合适的服务
    const services = await this.marketplace.discover({
      level: serviceLevel,
      capabilities: intent.requiredCapabilities,
    })

    // 4. 路由请求
    const result = await this.marketplace.routeRequest({
      level: serviceLevel,
      ...intent.requestParams,
    })

    return result
  }

  // 服务级别判断
  private determineServiceLevel(intent: UserIntent): ServiceLevel {
    // 简单对话 → Level 1 (LLM)
    if (intent.type === 'chat' || intent.type === 'question') {
      return 1
    }

    // 明确任务 → Level 2 (Agent)
    if (intent.type === 'task' && intent.scope === 'specific') {
      return 2
    }

    // 复杂需求 → Level 3 (Workflow)
    if (intent.type === 'idea' || intent.type === 'project') {
      return 3
    }

    // 默认 Level 1
    return 1
  }

  // 快捷功能
  async quickChat(prompt: string): Promise<string>       // Level 1 快速对话
  async runAgent(task: string): Promise<AgentResult>     // Level 2 执行任务
  async startWorkflow(idea: string): Promise<Workflow>   // Level 3 启动流程

  // 服务管理
  async browseServices(): Promise<Service[]>             // 浏览可用服务
  async viewServiceHistory(): Promise<ServiceRecord[]>   // 查看服务历史
  async manageSubscriptions(): Promise<Subscription[]>   // 管理订阅
}
```

---

## 链上集成

服务市场与链上模块交互，实现去中心化的服务交易。

```typescript
// src/modules/marketplace/chain-integration.ts

export class MarketplaceChainClient {
  // 服务注册上链
  async registerServiceOnChain(registration: ServiceRegistration): Promise<TxResult>

  // 创建托管支付
  async createEscrowOnChain(params: EscrowParams): Promise<TxResult>

  // 释放支付
  async releaseEscrowOnChain(escrowId: string, result: ServiceResult): Promise<TxResult>

  // 争议处理
  async createDispute(escrowId: string, reason: string): Promise<TxResult>

  // MQ 更新
  async updateMQ(provider: Address, rating: number): Promise<TxResult>
}

// 与链上模块的关系
interface ChainModuleIntegration {
  // x/compute - 算力服务交易
  compute: {
    registerProvider: () => void
    submitRequest: () => void
    submitResponse: () => void
  }

  // x/escrow - 托管支付
  escrow: {
    createEscrow: () => void
    releaseEscrow: () => void
    refundEscrow: () => void
  }

  // x/dispute - 争议处理
  dispute: {
    createDispute: () => void
    voteOnDispute: () => void
    executeResolution: () => void
  }

  // x/mq - MQ 管理
  mq: {
    recordServiceCompletion: () => void
    updateProviderScore: () => void
  }
}
```

---

## 服务目录

| 服务层级 | 类型 | 输入 | 输出 | 定价模式 |
|----------|------|------|------|----------|
| Level 1 | LLM Chat | prompt | completion | 按 Token |
| Level 1 | LLM Embedding | text | vector | 按 Token |
| Level 2 | Coder Agent | 任务描述 | 代码/修复 | 按复杂度 |
| Level 2 | Researcher Agent | 研究问题 | 分析报告 | 按任务 |
| Level 2 | Writer Agent | 内容需求 | 文档内容 | 按任务 |
| Level 3 | Software Workflow | 想法 | 完整软件 | 按里程碑 |
| Level 3 | Content Workflow | 主题 | 发布内容 | 按打包 |
| Level 3 | Business Workflow | 商业目标 | 商业计划 | 按打包 |

---

## 架构总览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              用户层                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  GenieBot 界面 (客户端插件)                                           │   │
│  │  - 自然语言对话界面                                                   │   │
│  │  - 服务浏览和选择                                                     │   │
│  │  - 任务管理和跟踪                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          服务市场 (核心模块)                                  │
│  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐    │
│  │ 服务注册  │ │ 服务发现  │ │ 服务定价  │ │ 请求路由  │ │ 标准接口  │    │
│  └───────────┘ └───────────┘ └───────────┘ └───────────┘ └───────────┘    │
│                                                                             │
│  三层服务结构: LLM (Level 1) → Agent (Level 2) → Workflow (Level 3)        │
└─────────────────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          服务提供者插件层                                     │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐               │
│  │ LLM Provider    │ │ Agent Provider  │ │ Workflow        │               │
│  │ Plugin          │ │ Plugin          │ │ Provider Plugin │               │
│  │ - OpenAI        │ │ (OpenFang)      │ │                 │               │
│  │ - Anthropic     │ │ - Coder Agent   │ │ - Software      │               │
│  │ - Google        │ │ - Researcher    │ │ - Content       │               │
│  │ - Local Models  │ │ - Writer        │ │ - Business      │               │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘               │
└─────────────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          链上模块 (Cosmos SDK)                               │
│  x/compute | x/escrow | x/dispute | x/mq | x/identity | x/task             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

[← 上一章：x/identity](./10-identity.md) | [返回索引](./00-index.md) | [下一章：其他模块 →](./12-misc.md)
