# ShareTokens 详细设计规范

> **目的:** 汇总所有需要补充的详细设计细节，作为开发实施的具体指导
> **版本:** 1.0
> **创建日期:** 2026-03-02

---

## 目录

1. [架构设计](#1-架构设计)
2. [接口定义补充](#2-接口定义补充)
3. [业务逻辑补充](#3-业务逻辑补充)
4. [技术实现补充](#4-技术实现补充)
5. [状态市场与插件系统](#5-服务市场与插件系统)
6. [状态机补充](#6-状态机补充)
7. [事件系统](#7-事件系统)
8. [数据流设计](#8-数据流设计)
9. [边界条件](#9-边界条件)
10. [安全规范](#10-安全规范)
11. [补充细节](#11-补充细节)

---

## 1. 架构设计

### 1.1 系统架构总览

ShareTokens 采用模块化架构，分为**核心模块**和**可选模块（插件）**两大类。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ShareTokens 节点架构                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         核心模块 (Core Modules)                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐ │   │
│  │  │ P2P通信  │  │ 身份账号 │  │   钱包   │  │     服务市场         │ │   │
│  │  │(Libp2p)  │  │(Identity)│  │ (Wallet) │  │ (三层服务架构)       │ │   │
│  │  └──────────┘  └──────────┘  └──────────┘  │  LLM/Agent/Workflow  │ │   │
│  │                                            └──────────────────────┘ │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐                           │   │
│  │  │ 托管支付 │  │ 德商系统 │  │ 争议仲裁 │                           │   │
│  │  │(Escrow)  │  │  (MQ)    │  │(Dispute) │                           │   │
│  │  └──────────┘  └──────────┘  └──────────┘                           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         插件系统 (Plugin System)                     │   │
│  │  ┌─────────────────────────────┐  ┌─────────────────────────────┐   │   │
│  │  │    服务提供者插件            │  │      用户插件               │   │   │
│  │  │  ┌─────────────────────────┐│  │  ┌─────────────────────────┐│   │   │
│  │  │  │ LLM API Key 托管插件    ││  │  │     GenieBot界面插件    ││   │   │
│  │  │  └─────────────────────────┘│  │  │   (GenieBot UI)       ││   │   │
│  │  │  ┌─────────────────────────┐│  │  └─────────────────────────┘│   │   │
│  │  │  │ Agent 执行器 (OpenFang) ││  │                              │   │   │
│  │  │  └─────────────────────────┘│  │                              │   │   │
│  │  │  ┌─────────────────────────┐│  │                              │   │   │
│  │  │  │ Workflow 执行器         ││  │                              │   │   │
│  │  │  └─────────────────────────┘│  │                              │   │   │
│  │  └─────────────────────────────┘  └─────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 核心模块定义

**每个节点必须具备的核心模块：**

| 模块名称 | 说明 | 职责 |
|---------|------|------|
| **P2P通信** | 基于 Libp2p 的点对点网络 | 节点发现、消息广播、数据同步 |
| **身份账号** | 身份验证与管理 | OAuth 验证、身份哈希、等级管理 |
| **钱包** | 资产管理 | STT 代币存储、转账、签名 |
| **服务市场** | 核心业务模块 | 三层服务交易 (LLM/Agent/Workflow) |
| **托管支付** | 资金托管 | 创建托管、释放资金、争议锁定 |
| **德商系统** | 信誉评分 | 零和配权、多方案评分、评审义务 |
| **争议仲裁** | 纠纷处理 | 协商、投票、裁决、上诉 |

### 1.3 可选模块（插件系统）

#### 1.3.1 服务提供者插件

| 插件名称 | 功能 | 接口 |
|---------|------|------|
| **LLM API Key 托管插件** | 托管 OpenAI/Claude 等 API Key，提供 LLM 服务 | `ILLMProvider` |
| **Agent 执行器 (OpenFang)** | 执行 AI Agent 任务 | `IAgentExecutor` |
| **Workflow 执行器** | 执行工作流编排任务 | `IWorkflowExecutor` |

#### 1.3.2 用户插件

| 插件名称 | 功能 | 接口 |
|---------|------|------|
| **GenieBot界面插件** | 提供用户交互界面，支持服务发现、任务管理 | `IUserInterface` |

### 1.4 服务市场三层架构

服务市场是核心业务模块，提供三层服务类型：

```
┌────────────────────────────────────────────────────────────┐
│                     服务市场架构                            │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Layer 1: LLM 服务 (大语言模型推理)                         │
│  ├── 输入: Prompt (文本/图片)                              │
│  ├── 输出: Completion (文本)                               │
│  ├── 计费: Token 用量                                      │
│  └── 示例: GPT-4, Claude, LLaMA                            │
│                                                            │
│  Layer 2: Agent 服务 (智能代理执行)                         │
│  ├── 输入: 任务描述                                        │
│  ├── 输出: 任务结果                                        │
│  ├── 计费: 任务复杂度 + 执行时间                           │
│  └── 示例: 代码生成、数据分析、文档编写                    │
│                                                            │
│  Layer 3: Workflow 服务 (工作流编排)                        │
│  ├── 输入: 工作流定义                                      │
│  ├── 输出: 工作流执行结果                                  │
│  ├── 计费: 节点数量 + 执行资源                             │
│  └── 示例: 多步骤自动化、数据处理流水线                    │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

**关键设计决策：**
- Workflow 是服务类型，由服务市场统一管理
- GenieBot界面是用户插件，用于消费服务市场的服务
- 每层服务都有独立的定价模型和匹配机制

### 1.5 模块依赖关系

```
                    ┌──────────────┐
                    │   P2P通信    │
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐      ┌─────▼─────┐     ┌─────▼─────┐
    │ 身份账号 │      │   钱包    │     │ 服务市场  │
    └────┬────┘      └─────┬─────┘     └─────┬─────┘
         │                 │                 │
         │           ┌─────▼─────┐           │
         │           │ 托管支付  │◄──────────┤
         │           └─────┬─────┘           │
         │                 │                 │
         └────────►┌───────▼───────┐◄────────┘
                   │   德商系统    │
                   └───────┬───────┘
                           │
                    ┌──────▼──────┐
                    │  争议仲裁   │
                    └─────────────┘
```

---

## 2. 接口定义补充

### 2.1 统一错误处理规范 (ShareTokensError)

```typescript
// 链下服务统一错误类型
interface ShareTokensError {
  code: ErrorCode;
  message: string;
  details?: Record<string, unknown>;
  traceId: string;
  timestamp: number;
}

enum ErrorCode {
  // 通用错误 (1xxx)
  UNKNOWN = 1000,
  INVALID_PARAMETER = 1001,
  UNAUTHORIZED = 1002,
  FORBIDDEN = 1003,
  NOT_FOUND = 1004,
  RATE_LIMITED = 1005,

  // 身份错误 (2xxx)
  IDENTITY_NOT_VERIFIED = 2001,
  IDENTITY_ALREADY_REGISTERED = 2002,
  IDENTITY_REVOKED = 2003,
  IDENTITY_PROOF_INVALID = 2004,

  // 德商错误 (3xxx)
  MQ_INSUFFICIENT = 3001,
  MQ_LOCKED = 3002,
  MQ_ZERO = 3003,  // 德商已归零

  // 评审义务错误 (3x1x)
  JURY_INELIGIBLE = 3101,  // 无评审资格
  JURY_SUSPENDED = 3102,   // 评审资格暂停
  JURY_ABSENCE_LIMIT = 3103, // 缺席次数过多

  // 算力交易错误 (4xxx)
  COMPUTE_REQUEST_NOT_FOUND = 4001,
  COMPUTE_PROVIDER_NOT_AVAILABLE = 4002,
  COMPUTE_EXECUTION_TIMEOUT = 4003,
  COMPUTE_VERIFICATION_FAILED = 4004,
  API_KEY_ENCRYPTION_FAILED = 4005,

  // 托管错误 (5xxx)
  ESCROW_INSUFFICIENT_BALANCE = 5001,
  ESCROW_LOCKED_BY_DISPUTE = 5002,
  ESCROW_ALREADY_RELEASED = 5003,
  ESCROW_EXPIRED = 5004,

  // 争议错误 (6xxx)
  DISPUTE_NOT_FOUND = 6001,
  DISPUTE_ALREADY_RESOLVED = 6002,
  DISPUTE_APPEAL_LIMIT_REACHED = 6003,
  JURY_DUTY_DECLINED = 6004,

  // 任务错误 (7xxx)
  TASK_NOT_FOUND = 7001,
  TASK_ALREADY_ASSIGNED = 7002,
  TASK_APPLICATION_REJECTED = 7003,
  TASK_EXECUTION_TIMEOUT = 7004,
  TASK_REVISION_LIMIT_REACHED = 7005,

  // 想法错误 (8xxx)
  IDEA_NOT_FOUND = 8001,
  IDEA_FUNDING_FAILED = 8002,
  IDEA_CONTRIBUTION_WEIGHT_INVALID = 8003,
}

// 错误构造器
class ErrorBuilder {
  static fromCode(code: ErrorCode, details?: Record<string, unknown>): ShareTokensError {
    return {
      code,
      message: ErrorMessageMap[code] || 'Unknown error',
      details,
      traceId: generateTraceId(),
      timestamp: Date.now(),
    };
  }
}

const ErrorMessageMap: Record<ErrorCode, string> = {
  [ErrorCode.UNKNOWN]: 'An unknown error occurred',
  [ErrorCode.INVALID_PARAMETER]: 'Invalid parameter provided',
  // ... 其他映射
};
```

```go
// 链上模块统一错误定义 (Cosmos SDK 风格)
// 各模块 types/errors.go

// x/compute/types/errors.go
var (
  ErrComputeRequestNotFound = sdkerrors.Register(ModuleName, 4001, "compute request not found")
  ErrProviderNotAvailable   = sdkerrors.Register(ModuleName, 4002, "compute provider not available")
  ErrExecutionTimeout       = sdkerrors.Register(ModuleName, 4003, "compute execution timeout")
  ErrVerificationFailed     = sdkerrors.Register(ModuleName, 4004, "result verification failed")
)
```

### 2.2 参数验证规则

```typescript
// 通用验证规则
const ValidationRules = {
  // 地址验证
  address: {
    pattern: /^share[a-z0-9]{38}$/,  // Cosmos SDK Bech32 格式
    message: 'Invalid ShareTokens address format',
  },

  // 金额验证
  amount: {
    min: 1,           // 最小 1 micro-STT
    max: 1e15,        // 最大 10^9 STT
    precision: 0,     // 整数 (micro-STT)
  },

  // 德商验证
  mq: {
    min: 10,          // 最低德商
    initial: 100,     // 初始德商
    max: 1000,        // 最大德商
  },

  // 文本验证
  title: {
    minLength: 3,
    maxLength: 200,
  },
  description: {
    minLength: 10,
    maxLength: 10000,
  },

  // 时间验证
  duration: {
    minSeconds: 60,           // 最小 1 分钟
    maxSeconds: 365 * 24 * 3600,  // 最大 1 年
  },
};

// 验证器
class Validator {
  static validateAddress(value: string): boolean {
    return ValidationRules.address.pattern.test(value);
  }

  static validateAmount(value: bigint, rules = ValidationRules.amount): boolean {
    return value >= rules.min && value <= rules.max;
  }

  static validateMQ(value: number): boolean {
    return value >= ValidationRules.mq.min && value <= ValidationRules.mq.max;
  }

  static validateTitle(value: string): boolean {
    const len = value.trim().length;
    return len >= ValidationRules.title.minLength && len <= ValidationRules.title.maxLength;
  }

  static validateDescription(value: string): boolean {
    const len = value.trim().length;
    return len >= ValidationRules.description.minLength && len <= ValidationRules.description.maxLength;
  }
}
```

### 2.3 链上链下通信事件类型

```typescript
// 事件类型定义
enum SystemEventType {
  // ===== 核心模块事件 =====

  // P2P 通信事件
  P2P_NODE_CONNECTED = 'p2p.node_connected',
  P2P_NODE_DISCONNECTED = 'p2p.node_disconnected',
  P2P_MESSAGE_RECEIVED = 'p2p.message_received',

  // 身份事件
  IDENTITY_REGISTERED = 'identity.registered',
  IDENTITY_REVOKED = 'identity.revoked',

  // 钱包事件
  WALLET_TRANSFER = 'wallet.transfer',
  WALLET_BALANCE_CHANGED = 'wallet.balance_changed',

  // 服务市场事件 (核心业务)
  SERVICE_REGISTERED = 'service.registered',
  SERVICE_DEREGISTERED = 'service.deregistered',
  SERVICE_REQUEST_CREATED = 'service.request_created',
  SERVICE_REQUEST_MATCHED = 'service.request_matched',
  SERVICE_EXECUTION_STARTED = 'service.execution_started',
  SERVICE_RESPONSE_SUBMITTED = 'service.response_submitted',
  SERVICE_COMPLETED = 'service.completed',
  SERVICE_FAILED = 'service.failed',

  // LLM 服务事件
  LLM_SERVICE_REGISTERED = 'llm.service_registered',
  LLM_REQUEST_CREATED = 'llm.request_created',
  LLM_TOKENS_USED = 'llm.tokens_used',
  LLM_RESPONSE_READY = 'llm.response_ready',

  // Agent 服务事件
  AGENT_SERVICE_REGISTERED = 'agent.service_registered',
  AGENT_TASK_STARTED = 'agent.task_started',
  AGENT_TASK_PROGRESS = 'agent.task_progress',
  AGENT_TASK_COMPLETED = 'agent.task_completed',

  // Workflow 服务事件
  WORKFLOW_REGISTERED = 'workflow.registered',
  WORKFLOW_EXECUTION_STARTED = 'workflow.execution_started',
  WORKFLOW_NODE_STARTED = 'workflow.node_started',
  WORKFLOW_NODE_COMPLETED = 'workflow.node_completed',
  WORKFLOW_EXECUTION_COMPLETED = 'workflow.execution_completed',

  // 托管事件
  ESCROW_CREATED = 'escrow.created',
  ESCROW_RELEASED = 'escrow.released',
  ESCROW_PARTIAL_RELEASE = 'escrow.partial_release',
  ESCROW_LOCKED = 'escrow.locked',

  // 德商事件
  MQ_INITIALIZED = 'mq.initialized',
  MQ_REDISTRIBUTED = 'mq.redistributed',

  // 评审事件
  JURY_ABSENCE = 'jury.absence',

  // 争议事件
  DISPUTE_CREATED = 'dispute.created',
  DISPUTE_NEGOTIATION_STARTED = 'dispute.negotiation_started',
  DISPUTE_VOTING_STARTED = 'dispute.voting_started',
  DISPUTE_RESOLVED = 'dispute.resolved',
  DISPUTE_APPEALED = 'dispute.appealed',

  // ===== 插件系统事件 =====

  // 插件生命周期事件
  PLUGIN_LOADED = 'plugin.loaded',
  PLUGIN_INITIALIZED = 'plugin.initialized',
  PLUGIN_STARTED = 'plugin.started',
  PLUGIN_STOPPED = 'plugin.stopped',
  PLUGIN_ERROR = 'plugin.error',
}

// 事件结构
interface SystemEvent {
  id: string;
  type: SystemEventType;
  module: string;          // 模块名: p2p, identity, wallet, service, escrow, mq, dispute, plugin
  height: number;          // 区块高度
  txHash: string;          // 交易哈希
  timestamp: number;
  data: Record<string, unknown>;
  metadata?: {
    emitter: string;       // 事件发射者地址
    relatedEntities?: string[];  // 相关实体 ID
  };
}
```

### 2.4 API 版本控制策略

```yaml
# API 版本规范
api_versioning:
  current_version: "v1"
  header_name: "X-API-Version"
  url_prefix: "/api/v1"

  # 版本策略
  strategy:
    - type: url_prefix       # URL 前缀 /api/v1/
    - type: header           # Header X-API-Version: 1

  # 兼容性保证
  compatibility:
    breaking_change_policy: "major version bump"  # 破坏性变更需要大版本升级
    deprecation_period: "6 months"                # 弃用期 6 个月
    sunset_header: true                           # 使用 Sunset Header 通知

  # 版本生命周期
  lifecycle:
    alpha:
      duration: "3 months"
      stability: "unstable"
    beta:
      duration: "3 months"
      stability: "mostly stable"
    stable:
      duration: "12 months"
      stability: "stable"
    deprecated:
      duration: "6 months"
      stability: "security fixes only"
    retired:
      status: "removed"

# 响应格式
response_format:
  success:
    code: 200
    body:
      data: {}
      meta:
        api_version: "v1.2.3"
        timestamp: 1709251200000
        request_id: "req_abc123"

  error:
    code: 4xx/5xx
    body:
      error:
        code: "INVALID_PARAMETER"
        message: "Invalid parameter provided"
        details: {}
        trace_id: "trace_xyz789"
```

### 2.5 缓存策略接口

```typescript
// 缓存策略配置
interface CacheConfig {
  enabled: boolean;
  ttl: number;              // 秒
  strategy: CacheStrategy;
  invalidation: InvalidationPolicy;
}

enum CacheStrategy {
  CACHE_ASIDE = 'cache_aside',       // 应用层管理
  WRITE_THROUGH = 'write_through',   // 写入时同时更新缓存
  WRITE_BEHIND = 'write_behind',     // 异步写入缓存
  REFRESH_AHEAD = 'refresh_ahead',   // 主动刷新
}

interface InvalidationPolicy {
  method: 'ttl' | 'event' | 'manual';
  events?: SystemEventType[];  // 事件触发失效
}

// 各模块缓存配置
const ModuleCacheConfigs: Record<string, CacheConfig> = {
  // 身份缓存
  'identity:status': {
    enabled: true,
    ttl: 3600,  // 1 小时
    strategy: CacheStrategy.CACHE_ASIDE,
    invalidation: {
      method: 'event',
      events: [SystemEventType.IDENTITY_REGISTERED, SystemEventType.IDENTITY_REVOKED],
    },
  },

  // 德商缓存
  'mq:score': {
    enabled: true,
    ttl: 300,   // 5 分钟
    strategy: CacheStrategy.CACHE_ASIDE,
    invalidation: {
      method: 'event',
      events: [SystemEventType.MQ_REDISTRIBUTED],
    },
  },

  // 汇率缓存
  'price:stt_usd': {
    enabled: true,
    ttl: 60,    // 1 分钟
    strategy: CacheStrategy.REFRESH_AHEAD,
    invalidation: { method: 'ttl' },
  },

  // 任务列表缓存
  'task:list': {
    enabled: true,
    ttl: 30,    // 30 秒
    strategy: CacheStrategy.CACHE_ASIDE,
    invalidation: { method: 'ttl' },
  },
};

// 缓存接口
interface ICacheService {
  get<T>(key: string): Promise<T | null>;
  set<T>(key: string, value: T, ttl?: number): Promise<void>;
  delete(key: string): Promise<void>;
  invalidatePattern(pattern: string): Promise<void>;
}
```

### 2.6 监控指标接口

```typescript
// 监控指标定义
interface MetricDefinition {
  name: string;
  type: MetricType;
  description: string;
  labels: string[];
  unit?: string;
}

enum MetricType {
  COUNTER = 'counter',       // 只增计数器
  GAUGE = 'gauge',           // 可增可减
  HISTOGRAM = 'histogram',   // 直方图
  SUMMARY = 'summary',       // 摘要
}

// 系统指标
const SystemMetrics: MetricDefinition[] = [
  // 交易指标
  {
    name: 'sharetokens_tx_total',
    type: MetricType.COUNTER,
    description: 'Total number of transactions',
    labels: ['module', 'type', 'status'],
  },
  {
    name: 'sharetokens_tx_duration_seconds',
    type: MetricType.HISTOGRAM,
    description: 'Transaction processing duration',
    labels: ['module', 'type'],
    unit: 'seconds',
  },

  // 区块指标
  {
    name: 'sharetokens_block_height',
    type: MetricType.GAUGE,
    description: 'Current block height',
    labels: [],
  },
  {
    name: 'sharetokens_block_tx_count',
    type: MetricType.HISTOGRAM,
    description: 'Transactions per block',
    labels: [],
  },

  // 业务指标
  {
    name: 'sharetokens_compute_requests_total',
    type: MetricType.COUNTER,
    description: 'Total compute requests',
    labels: ['model', 'status'],
  },
  {
    name: 'sharetokens_disputes_active',
    type: MetricType.GAUGE,
    description: 'Number of active disputes',
    labels: ['type'],
  },
  {
    name: 'sharetokens_tasks_open',
    type: MetricType.GAUGE,
    description: 'Number of open tasks',
    labels: ['category'],
  },

  // 德商指标
  {
    name: 'sharetokens_mq_average',
    type: MetricType.GAUGE,
    description: 'Average MQ score',
    labels: [],
  },
  {
    name: 'sharetokens_mq_distribution',
    type: MetricType.HISTOGRAM,
    description: 'MQ score distribution',
    labels: [],
  },
];

// 指标收集接口
interface IMetricsCollector {
  incrementCounter(name: string, labels: Record<string, string>, value?: number): void;
  setGauge(name: string, labels: Record<string, string>, value: number): void;
  observeHistogram(name: string, labels: Record<string, string>, value: number): void;
}
```

---

## 3. 业务逻辑补充

### 3.1 德商系统详细设计

#### 2.1.1 德商零和配权原理

```
核心原则:
1. 零和博弈：惩罚池 = 奖励池，德商总量不变
2. 收敛性：德商越高增加越难（对数增长），高德商者付出>收获
3. 风险可控：每次最多损失德商的3%，永远不为负
4. 公正义务：参与裁决是义务，缺席扣代币/德商

初始值: 100
下限: 0 (无限趋近，永不触底)
上限: 无硬性上限，但对数增长自然收敛
```

#### 2.1.2 德商配权参数

```go
// 德商配权配置
type MQConfig struct {
    // 基础参数
    InitialMQ       uint64  // 100 - 初始德商

    // 配权参数
    RiskRate        sdk.Dec // 0.03 (3%) - 每次最大风险比例
    Lambda          sdk.Dec // 1.5 - 基准偏差
    MaxDeviation    sdk.Dec // 6.0 - 最大合理偏差（评分-10~10的范围）
    ChangeFactor    sdk.Dec // 0.01 (1%) - 基础变化系数
}

var DefaultMQConfig = MQConfig{
    InitialMQ:     100,
    RiskRate:      sdk.NewDecWithPrec(3, 2),   // 3%
    Lambda:        sdk.MustNewDecFromStr("1.5"),
    MaxDeviation:  sdk.MustNewDecFromStr("6.0"),
    ChangeFactor:  sdk.NewDecWithPrec(1, 2),   // 1%
}
```

#### 2.1.3 多方案评分流程

```
1. 方案提出:
   - 当事人A提出解决方案
   - 当事人B提出解决方案

2. 双方评分:
   - 双方对所有方案评分 (-10 ~ 10)
   - 计算分歧程度

3. 分歧判断:
   - 分歧小 → 自动和解
   - 分歧大 → AI补充中间方案

4. AI撮合:
   - AI分析双方评分
   - 提出折中方案
   - 再次撮合协商

5. 评审团评分 (如仍无法和解):
   - 随机选择评审团
   - 每人对每个方案评分
   - 计算德商加权平均分（群体共识）
   - 计算每个人与共识的偏差

6. 德商配权:
   - 偏差 < λ → 获得奖励
   - 偏差 = λ → 不奖不罚
   - 偏差 > λ → 受到惩罚
```

#### 2.1.4 德商配权算法

```go
// 第一步：计算惩罚池
for each person i where d_i > λ:
    // 惩罚比例 = 3% × (d - λ) / (max_d - λ)
    penaltyRate = RiskRate × (d_i - λ) / (max_d - λ)
    penaltyRate = min(penaltyRate, RiskRate)  // 不超过3%

    loss_i = D_i × penaltyRate
    penaltyPool += loss_i

// 第二步：计算奖励分配
for each person i where d_i < λ:
    // 贡献分 = 公正程度 × 对数抑制（收敛）
    // contrib = (λ - d) × log(D + 1)
    contrib_i = (λ - d_i) × log(D_i + 1)
    totalContrib += contrib_i

// 第三步：分配奖励池
allocFactor = penaltyPool / totalContrib

for each person i where d_i < λ:
    gain_i = contrib_i × allocFactor

// 确保零和: Σ(gain) = Σ(loss) = penaltyPool
```

#### 2.1.5 算法示例

```
参数: λ=1.5, max_d=6, risk=3%

参与者:
┌──────┬────────┬────────┬────────┬─────────────────────────────────┐
│ 人员 │ 德商 D │ 偏差 d │ 类型   │ 计算                            │
├──────┼────────┼────────┼────────┼─────────────────────────────────┤
│ A    │ 100    │ 0      │ 奖励   │ contrib = 1.5 × log(101) = 6.9  │
│ B    │ 200    │ 0.5    │ 奖励   │ contrib = 1.0 × log(201) = 5.3  │
│ C    │ 100    │ 1.5    │ 不奖罚 │ -                               │
│ D    │ 100    │ 3.0    │ 惩罚   │ rate = 3% × 1.5/4.5 = 1%        │
│      │        │        │        │ loss = 100 × 1% = 1             │
│ E    │ 200    │ 6.0    │ 惩罚   │ rate = 3% × 4.5/4.5 = 3%        │
│      │        │        │        │ loss = 200 × 3% = 6             │
└──────┴────────┴────────┴────────┴─────────────────────────────────┘

惩罚池 = 1 + 6 = 7
总贡献 = 6.9 + 5.3 = 12.2
分配系数 = 7 / 12.2 = 0.57

结果:
┌──────┬─────────┬──────────┬───────────┐
│ 人员 │ 变化    │ 新德商   │ 比例变化  │
├──────┼─────────┼──────────┼───────────┤
│ A    │ +3.9    │ 103.9    │ +3.9%     │
│ B    │ +3.0    │ 203.0    │ +1.5%     │ ← 高德商比例增长慢
│ C    │ 0       │ 100      │ 0%        │
│ D    │ -1      │ 99       │ -1%       │
│ E    │ -6      │ 194      │ -3%       │ ← 高德商风险大
└──────┴─────────┴──────────┴───────────┘

验证:
- 零和: 3.9 + 3.0 = 6.9 ≈ 7 (惩罚池) ✓
- 收敛: 低德商A (+3.9%) > 高德商B (+1.5%) ✓
- 安全: 最多损失3%，不为负 ✓
```

#### 2.1.6 评审义务与缺席惩罚

```go
// 评审义务配置
type JuryDutyConfig struct {
    // 缺席惩罚
    AbsenceTokenPenalty  sdk.Coin  // 扣代币
    AbsenceMQPenalty     uint64    // 扣德商

    // 多次缺席
    MaxAbsences         uint64    // 最大缺席次数
    MultipleAbsencePenalty uint64 // 多次缺席额外惩罚
}

var DefaultJuryDutyConfig = JuryDutyConfig{
    AbsenceTokenPenalty:     sdk.NewCoin("stt", sdk.NewInt(10_000000)), // 10 STT
    AbsenceMQPenalty:        5,  // 扣5点德商
    MaxAbsences:             3,  // 最多3次缺席
    MultipleAbsencePenalty:  10, // 多次缺席额外扣10点
}
```

### 3.2 评审团选择规则

```go
// 评审团规模配置
type JurySizeConfig struct {
    Small  uint64  // 小额争议评审人数
    Medium uint64  // 中等争议评审人数
    Large  uint64  // 大额争议评审人数
    Huge   uint64  // 巨额争议评审人数
}

// 金额阈值
type AmountThresholdConfig struct {
    Small  sdk.Int  // 小额上限 (100 STT)
    Medium sdk.Int  // 中额上限 (1000 STT)
    Large  sdk.Int  // 大额上限 (10000 STT)
}

var DefaultJuryConfig = struct {
    Size       JurySizeConfig
    Thresholds AmountThresholdConfig
}{
    Size: JurySizeConfig{
        Small:  3,   // 小额(≤100 STT): 3人
        Medium: 5,   // 中等(100-1000 STT): 5人
        Large:  7,   // 大额(1000-10000 STT): 7人
        Huge:   11,  // 巨额(>10000 STT): 11人
    },
    Thresholds: AmountThresholdConfig{
        Small:  sdk.NewInt(100_000000),   // 100 STT (100 × 10^6 micro-STT)
        Medium: sdk.NewInt(1000_000000),  // 1000 STT
        Large:  sdk.NewInt(10000_000000), // 10000 STT
    },
}

// 评审员资格要求
type JurorEligibility struct {
    MinMQ           uint64  // 最低德商要求 (默认 50)
    MinActiveDays   uint64  // 最少活跃天数 (默认 30)
    MaxActiveDisputes uint64  // 最大同时参与争议数 (默认 3)
    MinIdentityLevel string // 最低身份等级 (默认 "basic")
}

var DefaultJurorEligibility = JurorEligibility{
    MinMQ:             50,
    MinActiveDays:     30,
    MaxActiveDisputes: 3,
    MinIdentityLevel:  "basic",
}

// 评审团选择算法
func SelectJury(
    amount sdk.Int,
    eligibleJurors []JurorCandidate,
    excludeAddresses []sdk.AccAddress,
) ([]sdk.AccAddress, error) {

    // 1. 确定评审团规模
    size := DetermineJurySize(amount)

    // 2. 过滤不合格的候选人
    candidates := FilterEligibleJurors(eligibleJurors, excludeAddresses)

    if uint64(len(candidates)) < size {
        return nil, errors.New("insufficient eligible jurors")
    }

    // 3. 加权随机选择 (权重 = MQ^1.2)
    selected := WeightedRandomSelect(candidates, size)

    return selected, nil
}
```

### 3.3 托管释放规则

```go
// 托管释放配置
type EscrowReleaseConfig struct {
    // 自动释放条件
    AutoReleaseOnCompletion bool          // 任务完成时自动释放
    AutoReleaseDelay        time.Duration // 自动释放延迟 (冷静期)

    // 部分释放
    AllowPartialRelease     bool          // 是否允许部分释放
    PartialReleaseThreshold sdk.Dec       // 部分释放阈值 (如 50%)

    // 争议锁定
    DisputeLockPeriod       time.Duration // 争议锁定期间
}

var DefaultEscrowReleaseConfig = EscrowReleaseConfig{
    AutoReleaseOnCompletion: true,
    AutoReleaseDelay:        24 * time.Hour,  // 24 小时冷静期
    AllowPartialRelease:     true,
    PartialReleaseThreshold: sdk.NewDecWithPrec(50, 2),  // 50%
    DisputeLockPeriod:       7 * 24 * time.Hour,  // 7 天
}

// 释放类型
type ReleaseType string

const (
    ReleaseTypeFull        ReleaseType = "full"         // 全额释放
    ReleaseTypePartial     ReleaseType = "partial"      // 部分释放
    ReleaseTypeAuto        ReleaseType = "auto"         // 自动释放
    ReleaseTypeDispute     ReleaseType = "dispute"      // 争议裁决释放
    ReleaseTypeRefund      ReleaseType = "refund"       // 退款
)

// 释放条件检查
func CheckReleaseConditions(escrow Escrow, releaseType ReleaseType) error {
    // 任务完成：自动释放
    if releaseType == ReleaseTypeAuto && escrow.TaskCompleted {
        if time.Since(escrow.CompletedAt) > DefaultEscrowReleaseConfig.AutoReleaseDelay {
            return nil
        }
        return errors.New("release delay not met")
    }

    // 部分释放：需双方同意
    if releaseType == ReleaseTypePartial {
        if !escrow.CreatorApproved || !escrow.BeneficiaryApproved {
            return errors.New("both parties must agree for partial release")
        }
    }

    // 争议锁定：争议期间不可动
    if escrow.DisputeLocked {
        return errors.New("escrow is locked due to active dispute")
    }

    return nil
}
```

### 3.4 争议处理时间线

```go
// 争议处理配置
type DisputeTimeoutConfig struct {
    NegotiationPeriod  time.Duration  // 协商期
    VotingPeriod       time.Duration  // 投票期
    AppealPeriod       time.Duration  // 上诉期
    EvidenceDeadline   time.Duration  // 证据提交截止 (投票结束前)
}

var DefaultDisputeTimeoutConfig = DisputeTimeoutConfig{
    NegotiationPeriod:  7 * 24 * time.Hour,  // 协商期: 7天
    VotingPeriod:       3 * 24 * time.Hour,  // 投票期: 3天
    AppealPeriod:       7 * 24 * time.Hour,  // 上诉期: 7天
    EvidenceDeadline:   24 * time.Hour,      // 证据截止: 24小时
}

// 投票通过阈值
type VotingThresholdConfig struct {
    MajorityThreshold  sdk.Dec  // 通过阈值 (默认 60%)
    QuorumRequired     sdk.Dec  // 法定人数 (默认 50%)
}

var DefaultVotingThreshold = VotingThresholdConfig{
    MajorityThreshold: sdk.NewDecWithPrec(60, 2),  // 60%
    QuorumRequired:    sdk.NewDecWithPrec(50, 2),  // 50%
}

// 投票结果判定
func DetermineVerdict(votes []DisputeVote, totalWeight sdk.Dec) Verdict {
    var plaintiffWeight, defendantWeight, neutralWeight sdk.Dec

    for _, vote := range votes {
        weight := CalculateVotingWeight(vote.JurorMQ, vote.ExpertBonus)
        switch vote.Verdict {
        case VerdictPlaintiff:
            plaintiffWeight = plaintiffWeight.Add(weight)
        case VerdictDefendant:
            defendantWeight = defendantWeight.Add(weight)
        case VerdictNeutral:
            neutralWeight = neutralWeight.Add(weight)
        }
    }

    // 检查法定人数
    votedWeight := plaintiffWeight.Add(defendantWeight).Add(neutralWeight)
    if votedWeight.Quo(totalWeight).LT(DefaultVotingThreshold.QuorumRequired) {
        return VerdictNeutral  // 法定人数不足，中立判定
    }

    // 检查通过阈值
    if plaintiffWeight.Quo(votedWeight).GTE(DefaultVotingThreshold.MajorityThreshold) {
        return VerdictPlaintiff
    }
    if defendantWeight.Quo(votedWeight).GTE(DefaultVotingThreshold.MajorityThreshold) {
        return VerdictDefendant
    }

    return VerdictNeutral  // 未达阈值，中立判定
}
```

### 3.5 收益分配规则

```go
// 收益分配配置
type RevenueDistributionConfig struct {
    // 贡献权重
    ContributionWeights map[ContributionType]sdk.Dec

    // 平台手续费
    PlatformFeeRate     sdk.Dec  // 平台手续费率
    PlatformFeeAddress  sdk.AccAddress
}

var DefaultRevenueConfig = RevenueDistributionConfig{
    ContributionWeights: map[ContributionType]sdk.Dec{
        ContributionTypeCode:   sdk.NewDecWithPrec(30, 2),  // 代码: 30%
        ContributionTypeDesign: sdk.NewDecWithPrec(50, 2),  // 设计: 50%
        ContributionTypeBugFix: sdk.NewDecWithPrec(20, 2),  // Bug修复: 20%
    },
    PlatformFeeRate: sdk.NewDecWithPrec(2, 2),  // 平台手续费: 2%
}

// 收益分配计算
func CalculateRevenueDistribution(
    totalRevenue sdk.Coins,
    contributions []ContributionRecord,
    config RevenueDistributionConfig,
) (*RevenueDistributionResult, error) {

    // 1. 扣除平台手续费
    platformFee := calculatePlatformFee(totalRevenue, config.PlatformFeeRate)
    distributable := totalRevenue.Sub(platformFee)

    // 2. 计算各贡献者权重
    contributorWeights := make(map[string]sdk.Dec)
    var totalWeight sdk.Dec

    for _, c := range contributions {
        if !c.Verified {
            continue
        }

        // 基础权重 × 贡献类型权重
        typeWeight := config.ContributionWeights[c.Type]
        effectiveWeight := c.Weight.Mul(typeWeight)

        addr := c.Contributor.String()
        contributorWeights[addr] = contributorWeights[addr].Add(effectiveWeight)
        totalWeight = totalWeight.Add(effectiveWeight)
    }

    // 3. 计算每人分配金额
    distributions := make([]ContributorShare, 0)
    for addrStr, weight := range contributorWeights {
        sharePercent := weight.Quo(totalWeight)

        var amount sdk.Coins
        for _, coin := range distributable {
            allocAmt := coin.Amount.ToDec().Mul(sharePercent).TruncateInt()
            amount = append(amount, sdk.NewCoin(coin.Denom, allocAmt))
        }

        distributions = append(distributions, ContributorShare{
            Address:      sdk.MustAccAddressFromBech32(addrStr),
            TotalWeight:  weight,
            SharePercent: sharePercent,
            Amount:       amount.Sort(),
        })
    }

    return &RevenueDistributionResult{
        TotalRevenue:   totalRevenue,
        PlatformFee:    platformFee,
        Distributable:  distributable,
        Distributions:  distributions,
    }, nil
}
```

### 3.6 任务里程碑规则

```go
// 任务里程碑配置
type MilestoneConfig struct {
    MaxRevisions       uint64  // 最多修改次数
    ReviewTimeout      time.Duration  // 评审超时
    AutoApproveOnTimeout bool  // 超时自动通过
}

var DefaultMilestoneConfig = MilestoneConfig{
    MaxRevisions:         3,                    // 最多 3 次修改
    ReviewTimeout:        3 * 24 * time.Hour,   // 评审超时: 3 天
    AutoApproveOnTimeout: true,                 // 超时自动通过
}

// 里程碑提交验证
func ValidateMilestoneSubmission(milestone TaskMilestone, revisionCount uint64) error {
    // 检查修改次数
    if revisionCount > DefaultMilestoneConfig.MaxRevisions {
        return errors.New("maximum revision count exceeded")
    }

    // 检查状态
    if milestone.Status != "submitted" && milestone.Status != "rejected" {
        return errors.New("milestone not ready for submission")
    }

    return nil
}

// 评审超时处理
func CheckReviewTimeout(milestone TaskMilestone) (bool, error) {
    if milestone.SubmittedAt == nil {
        return false, nil
    }

    timeSinceSubmission := time.Since(*milestone.SubmittedAt)
    if timeSinceSubmission > DefaultMilestoneConfig.ReviewTimeout {
        if DefaultMilestoneConfig.AutoApproveOnTimeout {
            return true, nil  // 自动通过
        }
        return false, errors.New("review timeout")
    }

    return false, nil
}
```

---

## 4. 技术实现补充

### 4.1 Cosmos SDK 配置

```yaml
# chain/config.yml (Ignite CLI 配置)

version: 1
chain:
  name: sharetokens
  binary: sharetokensd

# 共识配置
consensus:
  block_time: 10s           # 出块时间
  max_tx_size: 1048576      # 最大交易大小 1MB
  max_gas: 10000000         # 最大 Gas

# 验证者配置
validators:
  min_stake: 1000000stt     # 最小质押 100万 STT
  max_validators: 21        # 最大验证者数量
  unbonding_period: 1209600s  # 解绑期 14 天

# 模块配置
modules:
  # ===== 核心模块 (每个节点必须有) =====
  - name: auth
  - name: bank
  - name: staking
  - name: params
  - name: p2p               # 自定义: P2P通信 (基于Libp2p)
  - name: identity          # 自定义: 身份账号
  - name: wallet            # 自定义: 钱包
  - name: service           # 自定义: 服务市场 (核心业务)
  - name: escrow            # 自定义: 托管支付
  - name: mq                # 自定义: 德商系统
  - name: dispute           # 自定义: 争议仲裁

  # ===== 可选模块 (按需加载) =====
  # - name: plugin_manager   # 插件管理器

# Genesis 配置
genesis:
  chain_id: sharetokens-1
  app_state:
    staking:
      params:
        bond_denom: stt
        unbonding_time: 1209600s
        max_validators: 21
    bank:
      supply:
        - denom: stt
          amount: "1000000000000000"  # 10亿 STT
```

### 4.2 Ignite CLI 命令

```bash
# ===== 项目初始化 =====
ignite scaffold chain sharetokens --module-dir=./x

# ===== x/identity 模块 =====
ignite scaffold module identity --dep auth,bank

ignite scaffold message register-identity \
  --module identity \
  --signer creator \
  identityType:string \
  identityHash:bytes \
  proof:object

ignite scaffold message revoke-identity \
  --module identity \
  --signer creator \
  identityHash:bytes

ignite scaffold query get-identity \
  --module identity \
  --req address:string

ignite scaffold query has-identity-hash \
  --module identity \
  --req identityHash:bytes

# ===== x/mq 模块 =====
ignite scaffold module mq --dep auth,identity

ignite scaffold message initialize-mq \
  --module mq \
  --signer creator

ignite scaffold message redistribute-mq \
  --module mq \
  --signer creator \
  disputeId:uint64 \
  plaintiff:string \
  defendant:string \
  verdict:string

ignite scaffold query get-mq \
  --module mq \
  --req address:string

ignite scaffold query select-jury \
  --module mq \
  --req disputeId:uint64 \
  --req size:uint64

# ===== x/compute 模块 =====
ignite scaffold module compute --dep auth,bank,escrow,identity

ignite scaffold message register-provider \
  --module compute \
  --signer creator \
  provider:string \
  models:list \
  encryptedKey:bytes \
  offer:object

ignite scaffold message submit-request \
  --module compute \
  --signer creator \
  model:string \
  promptHash:bytes \
  priceOffer:coin \
  timeout:duration

ignite scaffold message submit-response \
  --module compute \
  --signer creator \
  requestId:uint64 \
  resultHash:bytes \
  tokensUsed:object

ignite scaffold query get-request \
  --module compute \
  --req requestId:uint64

ignite scaffold query list-offers \
  --module compute \
  --req model:string

# ===== x/escrow 模块 =====
ignite scaffold module escrow --dep auth,bank

ignite scaffold message create-escrow \
  --module escrow \
  --signer creator \
  beneficiary:string \
  amount:coin \
  duration:duration

ignite scaffold message release-escrow \
  --module escrow \
  --signer creator \
  escrowId:uint64

ignite scaffold message partial-release \
  --module escrow \
  --signer creator \
  escrowId:uint64 \
  amount:coin

ignite scaffold query get-escrow \
  --module escrow \
  --req escrowId:uint64

# ===== x/dispute 模块 =====
ignite scaffold module dispute --dep auth,escrow,mq

ignite scaffold message create-dispute \
  --module dispute \
  --signer creator \
  orderId:uint64 \
  defendant:string \
  disputeType:string \
  title:string \
  description:string \
  amount:coin

ignite scaffold message submit-evidence \
  --module dispute \
  --signer creator \
  disputeId:uint64 \
  evidenceType:string \
  evidenceHash:bytes

ignite scaffold message cast-vote \
  --module dispute \
  --signer creator \
  disputeId:uint64 \
  verdict:string \
  reasoning:string

ignite scaffold query get-dispute \
  --module dispute \
  --req disputeId:uint64

# ===== x/idea 模块 =====
ignite scaffold module idea --dep auth,bank,escrow

ignite scaffold message create-idea \
  --module idea \
  --signer creator \
  title:string \
  description:string \
  category:string

ignite scaffold message create-campaign \
  --module idea \
  --signer creator \
  ideaId:uint64 \
  targetAmount:coin \
  endDate:time

ignite scaffold message contribute \
  --module idea \
  --signer creator \
  campaignId:uint64 \
  amount:coin

ignite scaffold query get-idea \
  --module idea \
  --req ideaId:uint64

# ===== x/task 模块 =====
ignite scaffold module task --dep auth,bank,escrow,idea

ignite scaffold message create-task \
  --module task \
  --signer creator \
  title:string \
  description:string \
  budget:coin \
  deadline:time

ignite scaffold message apply-for-task \
  --module task \
  --signer creator \
  taskId:uint64 \
  message:string

ignite scaffold message submit-work \
  --module task \
  --signer creator \
  taskId:uint64 \
  submissionHash:bytes

ignite scaffold message review-submission \
  --module task \
  --signer creator \
  taskId:uint64 \
  result:string \
  comment:string

ignite scaffold query get-task \
  --module task \
  --req taskId:uint64
```

### 4.3 链下服务架构

```yaml
# services/docker-compose.yml

services:
  # API 网关
  api-gateway:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - ai-service
      - oracle-service
      - matching-service
      - workflow-service

  # Redis 缓存
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  # IPFS 证据存储
  ipfs:
    image: ipfs/kubo:latest
    ports:
      - "4001:4001"
      - "5001:5001"
      - "8080:8080"
    volumes:
      - ipfs_data:/data/ipfs

  # InfluxDB 指标存储
  influxdb:
    image: influxdb:2.7
    ports:
      - "8086:8086"
    volumes:
      - influxdb_data:/var/lib/influxdb2
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_ORG=sharetokens
      - DOCKER_INFLUXDB_INIT_BUCKET=metrics

  # Elasticsearch 日志存储
  elasticsearch:
    image: elasticsearch:8.11.0
    ports:
      - "9200:9200"
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    volumes:
      - es_data:/usr/share/elasticsearch/data

  # Kibana 日志可视化
  kibana:
    image: kibana:8.11.0
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch

volumes:
  redis_data:
  ipfs_data:
  influxdb_data:
  es_data:
```

```typescript
// services/cache/redis.ts - Redis 缓存配置

import { createClient } from 'redis';

export interface CacheServiceConfig {
  url: string;
  keyPrefix: string;
  defaultTTL: number;
}

export class RedisCacheService implements ICacheService {
  private client: RedisClient;
  private config: CacheServiceConfig;

  constructor(config: CacheServiceConfig) {
    this.config = config;
    this.client = createClient({ url: config.url });
  }

  async get<T>(key: string): Promise<T | null> {
    const fullKey = `${this.config.keyPrefix}:${key}`;
    const data = await this.client.get(fullKey);
    if (!data) return null;
    return JSON.parse(data) as T;
  }

  async set<T>(key: string, value: T, ttl?: number): Promise<void> {
    const fullKey = `${this.config.keyPrefix}:${key}`;
    const data = JSON.stringify(value);
    const effectiveTTL = ttl ?? this.config.defaultTTL;
    await this.client.setEx(fullKey, effectiveTTL, data);
  }

  async delete(key: string): Promise<void> {
    const fullKey = `${this.config.keyPrefix}:${key}`;
    await this.client.del(fullKey);
  }

  async invalidatePattern(pattern: string): Promise<void> {
    const fullPattern = `${this.config.keyPrefix}:${pattern}`;
    const keys = await this.client.keys(fullPattern);
    if (keys.length > 0) {
      await this.client.del(keys);
    }
  }
}
```

### 4.4 安全机制

```yaml
# 安全配置

security:
  # API Key 加密
  api_key_encryption:
    algorithm: "AES-256-GCM"
    key_derivation: "PBKDF2"
    iterations: 100000
    salt_length: 32
    nonce_length: 12
    key_rotation_days: 30

  # 签名算法
  signature:
    algorithm: "secp256k1"
    hash_function: "SHA-256"
    # Cosmos SDK 默认使用 secp256k1

  # 密钥轮换
  key_rotation:
    enabled: true
    period_days: 30
    grace_period_days: 7
    notification_days_before: 14

  # 访问控制
  access_control:
    rate_limiting:
      enabled: true
      requests_per_minute: 60
      burst: 100

    ip_whitelist:
      enabled: false
      allowed_ips: []

    cors:
      enabled: true
      allowed_origins:
        - "https://sharetokens.io"
        - "https://*.sharetokens.io"
      allowed_methods:
        - GET
        - POST
        - PUT
        - DELETE
      allowed_headers:
        - Authorization
        - Content-Type
        - X-API-Version
```

```go
// x/compute/keeper/encryption.go - API Key 加密实现

package keeper

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"

    "golang.org/x/crypto/pbkdf2"
)

type EncryptionConfig struct {
    Iterations int
    SaltLen    int
    NonceLen   int
}

var DefaultEncryptionConfig = EncryptionConfig{
    Iterations: 100000,
    SaltLen:    32,
    NonceLen:   12,
}

// EncryptAPIKey 使用 AES-256-GCM 加密 API Key
func EncryptAPIKey(plaintext []byte, masterKey []byte, config EncryptionConfig) (string, error) {
    // 生成随机盐
    salt := make([]byte, config.SaltLen)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    // 使用 PBKDF2 派生密钥
    key := pbkdf2.Key(masterKey, salt, config.Iterations, 32, sha256.New)

    // 创建 AES-GCM 加密器
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    // 生成随机 nonce
    nonce := make([]byte, config.NonceLen)
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    // 加密
    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

    // 组合: salt + nonce + ciphertext
    result := append(salt, nonce...)
    result = append(result, ciphertext...)

    return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptAPIKey 解密 API Key
func DecryptAPIKey(encoded string, masterKey []byte, config EncryptionConfig) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return nil, err
    }

    // 提取 salt, nonce, ciphertext
    salt := data[:config.SaltLen]
    nonce := data[config.SaltLen : config.SaltLen+config.NonceLen]
    ciphertext := data[config.SaltLen+config.NonceLen:]

    // 派生密钥
    key := pbkdf2.Key(masterKey, salt, config.Iterations, 32, sha256.New)

    // 解密
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}
```

---

## 5. 服务市场与插件系统

### 5.1 服务市场核心接口

```typescript
// 服务市场核心接口定义

/**
 * 服务类型枚举
 */
enum ServiceType {
  LLM = 'llm',           // 大语言模型服务
  AGENT = 'agent',       // 智能代理服务
  WORKFLOW = 'workflow', // 工作流服务
}

/**
 * 服务基础接口
 */
interface IService {
  id: string;
  type: ServiceType;
  provider: string;        // 提供者地址
  name: string;
  description: string;
  pricing: PricingModel;
  status: ServiceStatus;
  createdAt: number;
  updatedAt: number;
}

/**
 * LLM 服务定义
 */
interface LLMService extends IService {
  type: ServiceType.LLM;
  model: string;           // 模型名称: gpt-4, claude-3, etc.
  maxTokens: number;
  supportedFormats: ('text' | 'image')[];
  latency: number;         // 平均响应延迟 ms
}

/**
 * Agent 服务定义
 */
interface AgentService extends IService {
  type: ServiceType.AGENT;
  capabilities: string[];  // 能力标签
  executorType: string;    // 执行器类型: openfang, custom
  maxDuration: number;     // 最大执行时间 ms
  sandboxConfig?: SandboxConfig;
}

/**
 * Workflow 服务定义
 */
interface WorkflowService extends IService {
  type: ServiceType.WORKFLOW;
  nodes: WorkflowNode[];   // 工作流节点定义
  edges: WorkflowEdge[];   // 节点连接关系
  triggers: TriggerConfig[]; // 触发条件
}

/**
 * 定价模型
 */
interface PricingModel {
  type: PricingType;
  basePrice: bigint;       // 基础价格 (micro-STT)
  unitPrice?: bigint;      // 单价
  tiers?: PricingTier[];   // 阶梯定价
}

enum PricingType {
  FIXED = 'fixed',         // 固定价格
  PER_TOKEN = 'per_token', // 按Token计费
  PER_HOUR = 'per_hour',   // 按小时计费
  PER_EXECUTION = 'per_execution', // 按执行次数
  TIERED = 'tiered',       // 阶梯定价
}

/**
 * 服务请求
 */
interface ServiceRequest {
  id: string;
  serviceType: ServiceType;
  requester: string;
  provider?: string;       // 可选，指定提供者
  input: ServiceInput;
  budget: bigint;
  timeout: number;
  status: RequestStatus;
}

/**
 * 服务响应
 */
interface ServiceResponse {
  requestId: string;
  provider: string;
  output: ServiceOutput;
  usage: UsageMetrics;
  price: bigint;
  signature: string;       // 提供者签名
  timestamp: number;
}

/**
 * 使用量指标
 */
interface UsageMetrics {
  tokensUsed?: number;     // LLM Token 用量
  executionTime?: number;  // 执行时间 ms
  nodeCount?: number;      // Workflow 节点数
  computeUnits?: number;   // 计算单元
}
```

### 5.2 插件系统接口

```typescript
/**
 * 插件基础接口
 */
interface IPlugin {
  id: string;
  name: string;
  version: string;
  type: PluginType;
  initialize(config: PluginConfig): Promise<void>;
  start(): Promise<void>;
  stop(): Promise<void>;
  getStatus(): PluginStatus;
}

enum PluginType {
  LLM_PROVIDER = 'llm_provider',
  AGENT_EXECUTOR = 'agent_executor',
  WORKFLOW_EXECUTOR = 'workflow_executor',
  USER_INTERFACE = 'user_interface',
}

/**
 * LLM 提供者插件接口
 */
interface ILLMProviderPlugin extends IPlugin {
  type: PluginType.LLM_PROVIDER;

  // 注册服务到市场
  registerService(service: LLMService): Promise<string>;

  // 处理请求
  handleRequest(request: ServiceRequest): Promise<ServiceResponse>;

  // 验证 API Key
  validateApiKey(encryptedKey: string): Promise<boolean>;

  // 获取支持的模型列表
  getSupportedModels(): Promise<ModelInfo[]>;
}

/**
 * Agent 执行器插件接口
 */
interface IAgentExecutorPlugin extends IPlugin {
  type: PluginType.AGENT_EXECUTOR;

  // 注册 Agent 服务
  registerAgent(service: AgentService): Promise<string>;

  // 执行 Agent 任务
  execute(request: ServiceRequest): Promise<ServiceResponse>;

  // 获取执行状态
  getExecutionStatus(executionId: string): Promise<ExecutionStatus>;
}

/**
 * Workflow 执行器插件接口
 */
interface IWorkflowExecutorPlugin extends IPlugin {
  type: PluginType.WORKFLOW_EXECUTOR;

  // 注册 Workflow 模板
  registerWorkflow(service: WorkflowService): Promise<string>;

  // 执行工作流
  execute(request: ServiceRequest): Promise<ServiceResponse>;

  // 暂停/恢复工作流
  pause(executionId: string): Promise<void>;
  resume(executionId: string): Promise<void>;
}

/**
 * 用户界面插件接口 (小灯)
 */
interface IUserInterfacePlugin extends IPlugin {
  type: PluginType.USER_INTERFACE;

  // 服务发现
  discoverServices(filter: ServiceFilter): Promise<IService[]>;

  // 创建服务请求
  createRequest(request: ServiceRequest): Promise<string>;

  // 查询请求状态
  getRequestStatus(requestId: string): Promise<RequestStatus>;

  // 构建和提交 Workflow
  buildWorkflow(definition: WorkflowDefinition): Promise<string>;

  // 管理用户偏好
  setUserPreferences(preferences: UserPreferences): Promise<void>;
}

/**
 * 插件管理器
 */
interface IPluginManager {
  // 加载插件
  loadPlugin(pluginPath: string): Promise<IPlugin>;

  // 卸载插件
  unloadPlugin(pluginId: string): Promise<void>;

  // 获取插件
  getPlugin(pluginId: string): IPlugin | undefined;

  // 获取所有插件
  getPlugins(type?: PluginType): IPlugin[];

  // 启用/禁用插件
  enablePlugin(pluginId: string): Promise<void>;
  disablePlugin(pluginId: string): Promise<void>;
}
```

### 5.3 插件与核心模块交互

```go
// x/service/keeper/plugin_manager.go - 插件管理器实现

package keeper

import (
    "context"
    "sync"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

// PluginManager 插件管理器
type PluginManager struct {
    mu      sync.RWMutex
    plugins map[string]Plugin
    router  *ServiceRouter
}

// Plugin 插件接口 (Go 侧)
type Plugin interface {
    ID() string
    Type() PluginType
    Initialize(ctx sdk.Context, config []byte) error
    HandleRequest(ctx sdk.Context, request ServiceRequest) (*ServiceResponse, error)
}

// ServiceRouter 服务路由器
type ServiceRouter struct {
    llmProviders      map[string]Plugin
    agentExecutors    map[string]Plugin
    workflowExecutors map[string]Plugin
}

// RegisterPlugin 注册插件
func (pm *PluginManager) RegisterPlugin(ctx sdk.Context, plugin Plugin) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // 验证插件
    if err := pm.validatePlugin(ctx, plugin); err != nil {
        return err
    }

    // 注册到对应路由
    switch plugin.Type() {
    case PluginTypeLLMProvider:
        pm.router.llmProviders[plugin.ID()] = plugin
    case PluginTypeAgentExecutor:
        pm.router.agentExecutors[plugin.ID()] = plugin
    case PluginTypeWorkflowExecutor:
        pm.router.workflowExecutors[plugin.ID()] = plugin
    }

    pm.plugins[plugin.ID()] = plugin
    return nil
}

// RouteRequest 路由请求到合适的插件
func (pm *PluginManager) RouteRequest(ctx sdk.Context, request ServiceRequest) (*ServiceResponse, error) {
    switch request.ServiceType {
    case ServiceTypeLLM:
        return pm.routeLLMRequest(ctx, request)
    case ServiceTypeAgent:
        return pm.routeAgentRequest(ctx, request)
    case ServiceTypeWorkflow:
        return pm.routeWorkflowRequest(ctx, request)
    default:
        return nil, ErrUnsupportedServiceType
    }
}

// routeLLMRequest 路由 LLM 请求
func (pm *PluginManager) routeLLMRequest(ctx sdk.Context, request ServiceRequest) (*ServiceResponse, error) {
    // 1. 查找匹配的服务
    service, err := pm.FindService(ctx, request.ServiceType, request.Provider)
    if err != nil {
        return nil, err
    }

    // 2. 获取对应的插件
    plugin, exists := pm.router.llmProviders[service.Provider]
    if !exists {
        return nil, ErrProviderNotAvailable
    }

    // 3. 创建托管
    escrowID, err := pm.escrowKeeper.CreateEscrow(ctx, request.Requester, service.Provider, request.Budget)
    if err != nil {
        return nil, err
    }

    // 4. 调用插件处理请求
    response, err := plugin.HandleRequest(ctx, request)
    if err != nil {
        // 退款
        pm.escrowKeeper.Refund(ctx, escrowID)
        return nil, err
    }

    // 5. 结算并释放托管
    actualCost := pm.calculateCost(service.Pricing, response.Usage)
    pm.escrowKeeper.Release(ctx, escrowID, actualCost)

    return response, nil
}
```

### 5.4 小灯界面插件设计

```typescript
// plugins/xiaodeng-ui/index.ts - 小灯界面插件实现

import {
  IUserInterfacePlugin,
  ServiceFilter,
  IService,
  ServiceRequest,
  RequestStatus,
  WorkflowDefinition,
  UserPreferences,
} from '@sharetokens/plugin-sdk';

export class XiaoDengUIPlugin implements IUserInterfacePlugin {
  id = 'xiaodeng-ui-v1';
  name = 'XiaoDeng UI Plugin';
  version = '1.0.0';
  type = PluginType.USER_INTERFACE;

  private coreClient: CoreClient;
  private preferences: UserPreferences;

  async initialize(config: PluginConfig): Promise<void> {
    this.coreClient = new CoreClient(config.nodeEndpoint);
  }

  async start(): Promise<void> {
    // 启动 UI 服务
    await this.startUIServer();
  }

  async stop(): Promise<void> {
    await this.stopUIServer();
  }

  getStatus(): PluginStatus {
    return {
      running: true,
      health: 'healthy',
    };
  }

  // 服务发现 - 浏览市场上的服务
  async discoverServices(filter: ServiceFilter): Promise<IService[]> {
    const services = await this.coreClient.serviceMarket.queryServices({
      type: filter.type,
      minRating: filter.minRating,
      maxPrice: filter.maxPrice,
      tags: filter.tags,
    });

    return services.map(this.transformService);
  }

  // 创建服务请求
  async createRequest(request: ServiceRequest): Promise<string> {
    // 1. 验证预算
    const balance = await this.coreClient.wallet.getBalance();
    if (balance < request.budget) {
      throw new Error('Insufficient balance');
    }

    // 2. 提交请求到服务市场
    const requestId = await this.coreClient.serviceMarket.createRequest({
      serviceType: request.serviceType,
      input: request.input,
      budget: request.budget,
      timeout: request.timeout,
      provider: request.provider, // 可选指定
    });

    // 3. 订阅状态更新
    this.subscribeToUpdates(requestId);

    return requestId;
  }

  // 查询请求状态
  async getRequestStatus(requestId: string): Promise<RequestStatus> {
    return await this.coreClient.serviceMarket.getRequestStatus(requestId);
  }

  // 构建 Workflow - 小灯的核心功能
  async buildWorkflow(definition: WorkflowDefinition): Promise<string> {
    // 1. 验证工作流定义
    this.validateWorkflowDefinition(definition);

    // 2. 将工作流注册为 Workflow 服务
    const serviceId = await this.coreClient.serviceMarket.registerWorkflow({
      name: definition.name,
      description: definition.description,
      nodes: definition.nodes,
      edges: definition.edges,
      triggers: definition.triggers,
      pricing: {
        type: PricingType.PER_EXECUTION,
        basePrice: definition.estimatedCost,
      },
    });

    return serviceId;
  }

  // 执行工作流
  async executeWorkflow(workflowId: string, input: any): Promise<string> {
    return await this.createRequest({
      serviceType: ServiceType.WORKFLOW,
      input: { workflowId, params: input },
      budget: await this.estimateWorkflowCost(workflowId),
      timeout: 3600000, // 1 小时
    });
  }

  // 设置用户偏好
  async setUserPreferences(preferences: UserPreferences): Promise<void> {
    this.preferences = preferences;
    await this.coreClient.storage.save('preferences', preferences);
  }

  // 私有方法
  private async subscribeToUpdates(requestId: string): Promise<void> {
    this.coreClient.events.subscribe(`request.${requestId}`, (event) => {
      this.notifyUser(event);
    });
  }

  private notifyUser(event: ServiceEvent): void {
    // 通知用户 (桌面通知、邮件等)
    console.log(`Request update: ${event.type}`, event.data);
  }
}

// 插件导出
export default XiaoDengUIPlugin;
```

### 5.5 服务匹配算法

```go
// x/service/keeper/matching.go - 服务匹配算法

package keeper

import (
    "sort"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

// MatchConfig 匹配配置
type MatchConfig struct {
    // 评分权重
    PriceWeight        sdk.Dec  // 价格权重 (默认 0.4)
    RatingWeight       sdk.Dec  // 评分权重 (默认 0.3)
    LatencyWeight      sdk.Dec  // 延迟权重 (默认 0.2)
    AvailabilityWeight sdk.Dec  // 可用性权重 (默认 0.1)

    // 匹配策略
    Strategy           MatchStrategy
}

type MatchStrategy string

const (
    MatchStrategyBestScore  MatchStrategy = "best_score"   // 最高评分
    MatchStrategyLowestPrice MatchStrategy = "lowest_price" // 最低价格
    MatchStrategyRoundRobin  MatchStrategy = "round_robin"  // 轮询
)

// ServiceProvider 服务提供者信息
type ServiceProvider struct {
    Address      string
    Service      IService
    Rating       sdk.Dec      // 服务评分 (基于德商和历史)
    Latency      int64        // 平均延迟 ms
    Availability sdk.Dec      // 可用性 (0-1)
    Price        sdk.Int      // 价格
}

// MatchProvider 匹配最佳服务提供者
func (k Keeper) MatchProvider(
    ctx sdk.Context,
    request ServiceRequest,
    config MatchConfig,
) (*ServiceProvider, error) {

    // 1. 获取所有符合条件的提供者
    providers := k.GetEligibleProviders(ctx, request)
    if len(providers) == 0 {
        return nil, ErrNoProviderAvailable
    }

    // 2. 计算综合评分
    scoredProviders := make([]ScoredProvider, 0, len(providers))
    for _, p := range providers {
        score := k.calculateScore(p, request, config)
        scoredProviders = append(scoredProviders, ScoredProvider{
            Provider: p,
            Score:    score,
        })
    }

    // 3. 根据策略选择
    switch config.Strategy {
    case MatchStrategyBestScore:
        return k.selectByScore(scoredProviders), nil
    case MatchStrategyLowestPrice:
        return k.selectLowestPrice(providers), nil
    case MatchStrategyRoundRobin:
        return k.selectRoundRobin(ctx, request.ServiceType, providers), nil
    default:
        return k.selectByScore(scoredProviders), nil
    }
}

// calculateScore 计算综合评分
func (k Keeper) calculateScore(
    provider ServiceProvider,
    request ServiceRequest,
    config MatchConfig,
) sdk.Dec {
    // 价格评分 (价格越低越好)
    priceScore := k.calculatePriceScore(provider.Price, request.Budget)

    // 评分评分
    ratingScore := provider.Rating

    // 延迟评分 (延迟越低越好)
    latencyScore := k.calculateLatencyScore(provider.Latency)

    // 可用性评分
    availabilityScore := provider.Availability

    // 加权综合评分
    totalScore := priceScore.Mul(config.PriceWeight).
        Add(ratingScore.Mul(config.RatingWeight)).
        Add(latencyScore.Mul(config.LatencyWeight)).
        Add(availabilityScore.Mul(config.AvailabilityWeight))

    return totalScore
}

// calculatePriceScore 计算价格评分
func (k Keeper) calculatePriceScore(price, budget sdk.Int) sdk.Dec {
    // 价格 / 预算，越接近 0 越好
    if budget.IsZero() {
        return sdk.ZeroDec()
    }
    ratio := price.ToDec().Quo(budget.ToDec())
    // 1 - ratio，价格越低分数越高
    score := sdk.OneDec().Sub(ratio)
    if score.IsNegative() {
        return sdk.ZeroDec()
    }
    return score
}

// calculateLatencyScore 计算延迟评分
func (k Keeper) calculateLatencyScore(latencyMs int64) sdk.Dec {
    // 延迟阈值: 100ms = 1.0, 1000ms = 0.5, 5000ms = 0.0
    if latencyMs <= 100 {
        return sdk.OneDec()
    }
    if latencyMs >= 5000 {
        return sdk.ZeroDec()
    }
    // 线性插值
    score := sdk.OneDec().Sub(sdk.NewDec(latencyMs - 100).Quo(sdk.NewDec(4900)))
    return score
}

// ScoredProvider 带评分的提供者
type ScoredProvider struct {
    Provider ServiceProvider
    Score    sdk.Dec
}

// selectByScore 选择评分最高的提供者
func (k Keeper) selectByScore(providers []ScoredProvider) *ServiceProvider {
    sort.Slice(providers, func(i, j int) bool {
        return providers[i].Score.GT(providers[j].Score)
    })
    return &providers[0].Provider
}
```

---

## 6. 状态机补充

### 6.1 服务请求状态

```
服务请求状态 (统一):

pending → matched → executing → completed
    ↓         ↓          ↓
cancelled cancelled  failed
                         ↓
                     disputed

完整状态:
- pending: 等待匹配
- matched: 已匹配提供者
- executing: 执行中
- completed: 已完成
- failed: 失败
- cancelled: 已取消
- disputed: 争议中
```

### 6.2 LLM 服务状态

```
LLM 服务状态:

idle → processing → responding → completed
  ↓         ↓
disabled  rate_limited

完整状态:
- idle: 空闲
- processing: 处理中
- responding: 响应中
- completed: 完成
- disabled: 已禁用
- rate_limited: 限流中
```

### 6.3 Agent 执行状态

```
Agent 执行状态:

initialized → planning → executing → verifying → completed
      ↓          ↓           ↓           ↓
   cancelled   cancelled   failed      failed
                               ↓
                           retrying
                               ↓
                            failed

完整状态:
- initialized: 初始化
- planning: 规划中
- executing: 执行中
- verifying: 验证中
- completed: 完成
- cancelled: 已取消
- failed: 失败
- retrying: 重试中
```

### 6.4 Workflow 执行状态

```
Workflow 执行状态:

initialized → running → paused → running → completed
      ↓          ↓        ↓
   cancelled  failed   cancelled
                 ↓
             partial_completed

节点状态:
- pending: 等待执行
- running: 执行中
- completed: 完成
- failed: 失败
- skipped: 跳过

完整状态:
- initialized: 初始化
- running: 运行中
- paused: 已暂停
- completed: 完成
- cancelled: 已取消
- failed: 失败
- partial_completed: 部分完成
```

### 6.5 算力交易状态

```
算力交易状态:

pending → matched → executing → verifying → completed
    ↓         ↓          ↓           ↓
cancelled  cancelled  failed      failed
                            ↓
                        disputed

完整状态:
- pending: 等待匹配
- matched: 已匹配提供者
- executing: 执行中
- verifying: 验证中
- completed: 已完成
- failed: 失败
- cancelled: 已取消
- disputed: 争议中
```

### 6.6 争议状态

```
争议状态:

pending → negotiating → voting → resolved → final
    ↓         ↓           ↓         ↓
cancelled  settled    appealed   (end)
                           ↓
                      voting → resolved → final

完整状态:
- pending: 争议创建
- negotiating: 协商阶段 (7天)
- voting: 投票阶段 (3天)
- resolved: 已裁决 (可上诉期 7天)
- appealed: 上诉中
- final: 最终裁决
- cancelled: 已取消
```

### 6.7 任务状态

```
任务状态:

draft → open → assigned → in_progress → under_review → completed
  ↓      ↓        ↓            ↓              ↓
cancelled expired cancelled   failed    revision_requested
                                        ↓         ↓
                                      disputed  under_review

完整状态:
- draft: 草稿
- open: 开放申请 (7天无人申请自动关闭)
- assigned: 已指派
- in_progress: 执行中 (14天未提交自动失败)
- under_review: 评审中 (3天自动通过)
- completed: 已完成
- revision_requested: 需修改 (最多3次)
- cancelled: 已取消
- expired: 已过期
- failed: 失败
- disputed: 争议中
```

### 6.8 想法状态

```
想法状态:

draft → published → funding → in_progress → completed
  ↓        ↓          ↓            ↓
archived archived  archived     archived

完整状态:
- draft: 草稿
- published: 已发布 (可见但未众筹)
- funding: 众筹中
- in_progress: 执行中
- completed: 已完成
- archived: 已归档
- abandoned: 已放弃
```

---

## 7. 事件系统

### 7.1 统一事件结构

```protobuf
// proto/sharetokens/base/event.proto

syntax = "proto3";

package sharetokens.base;

import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";

// 系统事件
message SystemEvent {
  string id = 1;                          // 事件唯一 ID
  string type = 2;                        // 事件类型
  string module = 3;                      // 模块名
  uint64 height = 4;                      // 区块高度
  string tx_hash = 5;                     // 交易哈希
  google.protobuf.Timestamp timestamp = 6;
  google.protobuf.Any data = 7;           // 事件数据
  EventMetadata metadata = 8;             // 元数据
}

message EventMetadata {
  string emitter = 1;                     // 事件发射者地址
  repeated string related_entities = 2;   // 相关实体 ID
  map<string, string> tags = 3;           // 标签
}
```

### 7.2 各模块事件定义

```protobuf
// proto/sharetokens/identity/events.proto

message EventIdentityRegistered {
  string address = 1;
  string identity_type = 2;
  bytes identity_hash = 3;
}

message EventIdentityRevoked {
  string address = 1;
  bytes identity_hash = 2;
}

// proto/sharetokens/mq/events.proto

message EventMQInitialized {
  string address = 1;
  uint64 initial_mq = 2;
}

message EventMQRedistributed {
  string dispute_id = 1;
  repeated ParticipantChange participants = 2;
  int64 total_penalty = 3;
  int64 total_reward = 4;
}

message ParticipantChange {
  string address = 1;
  string role = 2;  // "plaintiff" | "defendant" | "juror"
  uint64 before = 3;
  uint64 after = 4;
  int64 change = 5;
  double deviation = 6;
  bool is_penalized = 7;
  bool is_rewarded = 8;
}

message EventJuryAbsence {
  string address = 1;
  uint64 absence_count = 2;
  uint64 mq_penalty = 3;
  string token_penalty = 4;
}

// proto/sharetokens/compute/events.proto

message EventComputeRequestCreated {
  uint64 request_id = 1;
  string requester = 2;
  string model = 3;
  string price_offer = 4;
}

message EventComputeRequestMatched {
  uint64 request_id = 1;
  string provider = 2;
}

message EventComputeExecutionStarted {
  uint64 request_id = 1;
  string provider = 2;
}

message EventComputeResponseSubmitted {
  uint64 request_id = 1;
  uint64 response_id = 2;
  string provider = 3;
  uint64 tokens_used = 4;
}

message EventComputeCompleted {
  uint64 request_id = 1;
  string requester = 2;
  string provider = 3;
  string actual_cost = 4;
}

message EventComputeFailed {
  uint64 request_id = 1;
  string reason = 2;
}

// proto/sharetokens/escrow/events.proto

message EventEscrowCreated {
  uint64 escrow_id = 1;
  string creator = 2;
  string beneficiary = 3;
  string amount = 4;
}

message EventEscrowReleased {
  uint64 escrow_id = 1;
  string beneficiary = 2;
  string amount = 3;
}

message EventEscrowPartialRelease {
  uint64 escrow_id = 1;
  string beneficiary = 2;
  string amount = 3;
  string remaining = 4;
}

message EventEscrowLocked {
  uint64 escrow_id = 1;
  uint64 dispute_id = 2;
}

// proto/sharetokens/dispute/events.proto

message EventDisputeCreated {
  uint64 dispute_id = 1;
  uint64 order_id = 2;
  string plaintiff = 3;
  string defendant = 4;
  string amount = 5;
}

message EventDisputeNegotiationStarted {
  uint64 dispute_id = 1;
  google.protobuf.Timestamp deadline = 2;
}

message EventDisputeVotingStarted {
  uint64 dispute_id = 1;
  repeated string jurors = 2;
  google.protobuf.Timestamp deadline = 3;
}

message EventDisputeResolved {
  uint64 dispute_id = 1;
  string verdict = 2;
  string plaintiff_share = 3;
  string defendant_share = 4;
}

message EventDisputeAppealed {
  uint64 dispute_id = 1;
  string appellant = 2;
  string reason = 3;
}

// proto/sharetokens/task/events.proto

message EventTaskCreated {
  uint64 task_id = 1;
  string creator = 2;
  string title = 3;
  string budget = 4;
}

message EventTaskAssigned {
  uint64 task_id = 1;
  string assignee = 2;
}

message EventTaskSubmitted {
  uint64 task_id = 1;
  uint64 submission_id = 2;
  string assignee = 3;
}

message EventTaskReviewed {
  uint64 task_id = 1;
  string result = 2;  // approved, revision_requested, rejected
  string reviewer = 3;
}

message EventTaskCompleted {
  uint64 task_id = 1;
  string assignee = 2;
  string reward = 3;
}

// proto/sharetokens/idea/events.proto

message EventIdeaCreated {
  uint64 idea_id = 1;
  string creator = 2;
  string title = 3;
  string category = 4;
}

message EventIdeaCampaignStarted {
  uint64 idea_id = 1;
  uint64 campaign_id = 2;
  string target_amount = 3;
  google.protobuf.Timestamp end_date = 4;
}

message EventIdeaFunded {
  uint64 idea_id = 1;
  uint64 campaign_id = 2;
  string total_raised = 3;
  uint64 contributor_count = 4;
}

message EventIdeaContributionAdded {
  uint64 idea_id = 1;
  uint64 contribution_id = 2;
  string contributor = 3;
  string contribution_type = 4;
  string weight = 5;
}
```

### 7.3 事件发射示例

```go
// x/compute/keeper/msg_server.go

func (k msgServer) SubmitRequest(goCtx context.Context, msg *types.MsgSubmitRequest) (*types.MsgSubmitRequestResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // ... 业务逻辑 ...

    // 发射事件
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeComputeRequestCreated,
            sdk.NewAttribute(types.AttributeKeyRequestId, fmt.Sprintf("%d", requestId)),
            sdk.NewAttribute(types.AttributeKeyRequester, msg.Creator),
            sdk.NewAttribute(types.AttributeKeyModel, msg.Model),
            sdk.NewAttribute(types.AttributeKeyPriceOffer, msg.PriceOffer.String()),
        ),
    )

    return &types.MsgSubmitRequestResponse{RequestId: requestId}, nil
}
```

---

## 8. 数据流设计

### 8.1 服务请求数据流

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        服务请求完整数据流                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. 请求发起                                                                 │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────────┐                  │
│  │ 用户插件 │───►│ 服务市场 │───►│  托管支付 (锁定资金) │                  │
│  │ (小灯)   │    │ (匹配)   │    │                      │                  │
│  └──────────┘    └──────────┘    └──────────────────────┘                  │
│       │               │                       │                             │
│       │               ▼                       │                             │
│       │        ┌──────────┐                   │                             │
│       │        │ 德商检查 │                   │                             │
│       │        └──────────┘                   │                             │
│       │               │                       │                             │
│       └───────────────┼───────────────────────┘                             │
│                       ▼                                                     │
│  2. 服务执行                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                       服务提供者插件                                   │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐                     │  │
│  │  │ LLM 执行器 │  │Agent 执行器│  │Workflow执行│                     │  │
│  │  └────────────┘  └────────────┘  └────────────┘                     │  │
│  │         │              │               │                             │  │
│  │         └──────────────┼───────────────┘                             │  │
│  │                        ▼                                              │  │
│  │                 ┌────────────┐                                        │  │
│  │                 │ 结果签名   │                                        │  │
│  │                 └────────────┘                                        │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                       │                                                     │
│                       ▼                                                     │
│  3. 结果验证与结算                                                           │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────────┐                  │
│  │ 服务市场 │───►│ 结果验证 │───►│  托管支付 (释放资金) │                  │
│  │          │    │          │    │                      │                  │
│  └──────────┘    └──────────┘    └──────────────────────┘                  │
│       │               │                       │                             │
│       │               │                       ▼                             │
│       │               │              ┌──────────────┐                       │
│       │               │              │ 德商更新     │                       │
│       │               │              │ (交易完成)   │                       │
│       │               │              └──────────────┘                       │
│       │               │                                                     │
│       └───────────────┼─────────────────────────────────────────────────►  │
│                       ▼                                                     │
│  4. 争议处理 (可选)                                                          │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐             │
│  │ 发起争议 │───►│ 协商阶段 │───►│ 投票阶段 │───►│ 裁决执行 │             │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘             │
│       │                                                   │                │
│       │                                                   ▼                │
│       │                                          ┌──────────────┐          │
│       │                                          │ 德商重分配   │          │
│       │                                          └──────────────┘          │
│       │                                                                     │
│       └─────────────────────────────────────────────────────────────────►  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 插件与核心模块交互流

```typescript
// 数据流接口定义

/**
 * 核心模块到插件的数据流
 */
interface CoreToPluginFlow {
  // 1. 服务市场 -> LLM 插件
  llmRequest: {
    requestId: string;
    prompt: string | Buffer;  // 文本或图片
    model: string;
    maxTokens: number;
    temperature?: number;
    budget: bigint;
    callbackUrl: string;      // 结果回调地址
  };

  // 2. 服务市场 -> Agent 插件
  agentRequest: {
    requestId: string;
    taskDescription: string;
    context: Record<string, unknown>;
    constraints: AgentConstraints;
    budget: bigint;
    timeout: number;
    callbackUrl: string;
  };

  // 3. 服务市场 -> Workflow 插件
  workflowRequest: {
    requestId: string;
    workflowId: string;
    inputs: Record<string, unknown>;
    budget: bigint;
    timeout: number;
    callbackUrl: string;
  };
}

/**
 * 插件到核心模块的数据流
 */
interface PluginToCoreFlow {
  // 1. 插件 -> 服务市场 (注册)
  serviceRegistration: {
    serviceType: ServiceType;
    provider: string;
    serviceDefinition: IService;
    pricing: PricingModel;
    signature: string;
  };

  // 2. 插件 -> 服务市场 (结果)
  serviceResult: {
    requestId: string;
    output: ServiceOutput;
    usage: UsageMetrics;
    actualCost: bigint;
    signature: string;
    timestamp: number;
  };

  // 3. 插件 -> 身份账号 (验证)
  identityVerification: {
    provider: string;
    identityProof: IdentityProof;
    callbackUrl: string;
  };
}

/**
 * 用户插件 (小灯) 数据流
 */
interface UserPluginFlow {
  // 服务发现
  discoverServices: {
    filter: ServiceFilter;
    pagination: Pagination;
  } => {
    services: IService[];
    total: number;
  };

  // 创建请求
  createRequest: {
    request: ServiceRequest;
  } => {
    requestId: string;
    estimatedWaitTime: number;
  };

  // 监控请求
  monitorRequest: {
    requestId: string;
  } => Observable<RequestUpdate>;

  // 构建工作流
  buildWorkflow: {
    definition: WorkflowDefinition;
  } => {
    workflowId: string;
    estimatedCost: bigint;
    previewUrl: string;
  };
}
```

### 8.3 跨模块数据流

```go
// x/service/keeper/cross_module_flow.go - 跨模块数据流管理

package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

// CrossModuleFlowManager 跨模块数据流管理器
type CrossModuleFlowManager struct {
    serviceKeeper *Keeper
    escrowKeeper  EscrowKeeper
    mqKeeper      MQKeeper
    identityKeeper IdentityKeeper
    disputeKeeper DisputeKeeper
}

// ExecuteServiceRequest 执行服务请求 (跨模块流程)
func (fm *CrossModuleFlowManager) ExecuteServiceRequest(
    ctx sdk.Context,
    request ServiceRequest,
) (*ServiceResponse, error) {

    // ========== 第一步: 验证 ==========
    // 1.1 身份验证
    if err := fm.verifyIdentity(ctx, request.Requester); err != nil {
        return nil, err
    }

    // 1.2 德商检查
    requesterMQ := fm.mqKeeper.GetMQ(ctx, request.Requester)
    if requesterMQ < MinMQForRequest {
        return nil, ErrInsufficientMQ
    }

    // 1.3 余额检查
    balance := fm.escrowKeeper.GetBalance(ctx, request.Requester)
    if balance.LT(request.Budget) {
        return nil, ErrInsufficientBalance
    }

    // ========== 第二步: 托管 ==========
    escrowID, err := fm.escrowKeeper.CreateEscrow(
        ctx,
        request.Requester,
        "", // provider 待定
        request.Budget,
        EscrowTypeService,
    )
    if err != nil {
        return nil, err
    }

    // ========== 第三步: 匹配 ==========
    provider, err := fm.serviceKeeper.MatchProvider(ctx, request, DefaultMatchConfig)
    if err != nil {
        // 匹配失败，退款
        fm.escrowKeeper.Refund(ctx, escrowID)
        return nil, err
    }

    // 更新托管受益人
    fm.escrowKeeper.UpdateBeneficiary(ctx, escrowID, provider.Address)

    // ========== 第四步: 执行 ==========
    response, err := fm.serviceKeeper.ExecuteRequest(ctx, request, provider)
    if err != nil {
        // 执行失败，退款
        fm.escrowKeeper.Refund(ctx, escrowID)
        return nil, err
    }

    // ========== 第五步: 结算 ==========
    actualCost := fm.calculateActualCost(request, response)
    if err := fm.escrowKeeper.Release(ctx, escrowID, actualCost); err != nil {
        return nil, err
    }

    // ========== 第六步: 更新德商 ==========
    // 交易成功，双方获得小额德商奖励
    fm.mqKeeper.AwardTransactionMQ(ctx, request.Requester, provider.Address)

    // 记录交易历史
    fm.serviceKeeper.RecordTransaction(ctx, request, response, actualCost)

    return response, nil
}

// HandleDispute 处理争议 (跨模块流程)
func (fm *CrossModuleFlowManager) HandleDispute(
    ctx sdk.Context,
    dispute Dispute,
) error {

    // ========== 第一步: 锁定托管 ==========
    if err := fm.escrowKeeper.LockByDispute(ctx, dispute.OrderID, dispute.ID); err != nil {
        return err
    }

    // ========== 第二步: 协商阶段 ==========
    if err := fm.disputeKeeper.StartNegotiation(ctx, dispute.ID); err != nil {
        return err
    }

    // ========== 第三步: 投票阶段 (如果协商失败) ==========
    if dispute.Status == DisputeStatusNegotiationFailed {
        // 选择评审团
        jurors, err := fm.mqKeeper.SelectJury(ctx, dispute.Amount, DefaultJuryConfig)
        if err != nil {
            return err
        }

        if err := fm.disputeKeeper.StartVoting(ctx, dispute.ID, jurors); err != nil {
            return err
        }
    }

    // ========== 第四步: 裁决 ==========
    verdict := fm.disputeKeeper.GetVerdict(ctx, dispute.ID)

    // ========== 第五步: 执行裁决 ==========
    switch verdict {
    case VerdictPlaintiff:
        // 原告胜，资金退还
        fm.escrowKeeper.RefundToCreator(ctx, dispute.OrderID)
    case VerdictDefendant:
        // 被告胜，资金释放给提供者
        fm.escrowKeeper.ReleaseToBeneficiary(ctx, dispute.OrderID)
    case VerdictNeutral:
        // 中立，平分
        fm.escrowKeeper.Split(ctx, dispute.OrderID)
    }

    // ========== 第六步: 德商重分配 ==========
    if err := fm.mqKeeper.RedistributeMQ(ctx, dispute.ID, verdict); err != nil {
        return err
    }

    return nil
}
```

### 8.4 事件驱动数据流

```yaml
# 事件驱动数据流配置

event_flows:
  # 服务请求事件流
  service_request_flow:
    trigger: "service.request_created"
    steps:
      - event: "escrow.created"
        module: "escrow"
        action: "create_escrow"
      - event: "mq.checked"
        module: "mq"
        action: "verify_mq"
      - event: "service.matched"
        module: "service"
        action: "match_provider"
      - event: "service.executing"
        module: "service"
        action: "execute_request"
      - event: "service.completed"
        module: "service"
        action: "complete_request"
      - event: "escrow.released"
        module: "escrow"
        action: "release_funds"
      - event: "mq.updated"
        module: "mq"
        action: "update_mq"

  # 争议处理事件流
  dispute_flow:
    trigger: "dispute.created"
    steps:
      - event: "escrow.locked"
        module: "escrow"
        action: "lock_by_dispute"
      - event: "dispute.negotiating"
        module: "dispute"
        action: "start_negotiation"
      - event: "dispute.voting"
        module: "dispute"
        action: "start_voting"
        condition: "negotiation_failed"
      - event: "mq.jury_selected"
        module: "mq"
        action: "select_jury"
        condition: "voting_started"
      - event: "dispute.resolved"
        module: "dispute"
        action: "resolve"
      - event: "escrow.distributed"
        module: "escrow"
        action: "distribute_by_verdict"
      - event: "mq.redistributed"
        module: "mq"
        action: "redistribute"

  # 插件生命周期事件流
  plugin_lifecycle_flow:
    trigger: "plugin.loaded"
    steps:
      - event: "plugin.initialized"
        module: "plugin_manager"
        action: "initialize"
      - event: "plugin.registered"
        module: "service"
        action: "register_services"
        condition: "provider_plugin"
      - event: "plugin.started"
        module: "plugin_manager"
        action: "start"
      - event: "plugin.ready"
        module: "plugin_manager"
        action: "notify_ready"
```

### 8.5 数据一致性保证

```go
// x/base/keeper/consistency.go - 数据一致性保证

package keeper

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConsistencyManager 数据一致性管理器
type ConsistencyManager struct {
    eventBus *EventBus
    cache    *StateCache
}

// TransactionalOperation 事务性操作
type TransactionalOperation func(ctx sdk.Context) error

// ExecuteAtomically 原子执行操作
func (cm *ConsistencyManager) ExecuteAtomically(
    ctx sdk.Context,
    operations []TransactionalOperation,
) error {
    // 创建缓存上下文
    cacheCtx, writeCache := ctx.CacheContext()

    // 保存初始状态用于回滚
    snapshot := cm.cache.CreateSnapshot()

    // 执行所有操作
    for i, op := range operations {
        if err := op(cacheCtx); err != nil {
            // 回滚到初始状态
            cm.cache.RestoreSnapshot(snapshot)

            // 发射失败事件
            cm.eventBus.EmitEvent(Event{
                Type: "transaction.failed",
                Data: map[string]interface{}{
                    "operation_index": i,
                    "error":           err.Error(),
                },
            })

            return fmt.Errorf("operation %d failed: %w", i, err)
        }
    }

    // 所有操作成功，原子写入
    writeCache()

    // 发射成功事件
    cm.eventBus.EmitEvent(Event{
        Type: "transaction.completed",
        Data: map[string]interface{}{
            "operation_count": len(operations),
        },
    })

    return nil
}

// InvariantCheck 不变量检查
func (cm *ConsistencyManager) InvariantCheck(ctx sdk.Context) error {
    // 1. 德商总量恒定
    if err := cm.checkMQConservation(ctx); err != nil {
        return err
    }

    // 2. 资金平衡
    if err := cm.checkFundBalance(ctx); err != nil {
        return err
    }

    // 3. 服务状态一致性
    if err := cm.checkServiceStateConsistency(ctx); err != nil {
        return err
    }

    return nil
}

// checkMQConservation 检查德商守恒
func (cm *ConsistencyManager) checkMQConservation(ctx sdk.Context) error {
    expectedTotal := cm.mqKeeper.GetExpectedTotalMQ(ctx)
    actualTotal := cm.mqKeeper.GetActualTotalMQ(ctx)

    if !expectedTotal.Equal(actualTotal) {
        return fmt.Errorf("MQ conservation violated: expected %s, got %s",
            expectedTotal.String(), actualTotal.String())
    }

    return nil
}

// checkFundBalance 检查资金平衡
func (cm *ConsistencyManager) checkFundBalance(ctx sdk.Context) error {
    // 总供应 = 流通中 + 托管中 + 质押中
    totalSupply := cm.bankKeeper.GetSupply(ctx, "stt")
    circulating := cm.bankKeeper.GetCirculatingSupply(ctx, "stt")
    escrowed := cm.escrowKeeper.GetTotalEscrowed(ctx)
    staked := cm.stakingKeeper.GetTotalStaked(ctx)

    expectedTotal := circulating.Add(escrowed).Add(staked)
    if !totalSupply.Amount.Equal(expectedTotal) {
        return fmt.Errorf("fund balance mismatch: supply %s != circ %s + escrow %s + stake %s",
            totalSupply.String(), circulating.String(),
            escrowed.String(), staked.String())
    }

    return nil
}
```

---

## 9. 边界条件

### 9.1 数值边界

```yaml
token:
  min_unit: 0.00000001 STT  # 8位小数
  max_single_tx: 1,000,000 STT
  max_balance_per_user: 100,000,000 STT

mq:
  min_value: 0   # 无限趋近，永不触底
  max_value: 无硬性上限  # 对数增长自然收敛
  max_single_change: ±3%  # 每次最多变化3%
  initial_value: 100  # 初始德商

timeout:
  compute_request_max: 3600s  # 1小时
  dispute_total_max: 259200s  # 3天
  task_review_max: 259200s    # 3天
```

### 9.2 并发处理

- 任务分配：乐观锁 + 先到先得
- 出价锁定：5秒锁定窗口
- 投票同步：实时广播 + 最终一致性

### 9.3 资源配额

```yaml
user_quotas:
  storage: 100MB
  api_rate_limit: 100/minute
  max_active_tasks: 5
  max_active_ideas: 10
  max_concurrent_compute: 3
```

### 9.4 异常处理

- 交易失败：自动回滚
- 系统故障：自动重试3次
- 数据备份：每小时增量，每天全量
- 灾难恢复：RPO < 1小时，RTO < 4小时

---

## 10. 安全规范

### 10.1 密钥管理

```yaml
key_management:
  # API Key 加密
  api_key_encryption:
    algorithm: "AES-256-GCM"
    key_derivation: "PBKDF2"
    iterations: 100000
    salt_length: 32
    nonce_length: 12

  # 主密钥存储
  master_key_storage:
    type: "HSM"  # Hardware Security Module
    provider: "AWS CloudHSM"  # 或 Azure Dedicated HSM
    backup_enabled: true

  # 密钥轮换
  key_rotation:
    enabled: true
    period_days: 30
    grace_period_days: 7
    notification_days_before: 14

  # 密钥备份
  key_backup:
    method: "Shamir's Secret Sharing"
    threshold: "3/5"  # 5份中需要3份恢复
    storage_locations:
      - "HSM Primary"
      - "HSM Secondary (DR)"
      - "Offline Vault 1"
      - "Offline Vault 2"
      - "Trusted Custodian"
```

```go
// x/compute/keeper/key_management.go - 密钥管理实现

package keeper

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "time"

    "golang.org/x/crypto/pbkdf2"
)

// KeyMetadata 密钥元数据
type KeyMetadata struct {
    KeyID       string    // 密钥标识
    Version     uint64    // 版本号
    CreatedAt   time.Time // 创建时间
    ExpiresAt   time.Time // 过期时间
    Status      string    // active, rotating, deprecated
}

// MasterKeyManager 主密钥管理器
type MasterKeyManager struct {
    hsmClient    HSMClient
    currentKeyID string
    keyCache     map[string][]byte
    metadata     map[string]KeyMetadata
}

// RotateMasterKey 密钥轮换
func (m *MasterKeyManager) RotateMasterKey() error {
    // 1. 在 HSM 中生成新密钥
    newKeyID, err := m.hsmClient.GenerateKey()
    if err != nil {
        return err
    }

    // 2. 设置旧密钥为 deprecated (宽限期)
    if m.currentKeyID != "" {
        oldMeta := m.metadata[m.currentKeyID]
        oldMeta.Status = "deprecated"
        oldMeta.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
        m.metadata[m.currentKeyID] = oldMeta
    }

    // 3. 激活新密钥
    m.currentKeyID = newKeyID
    m.metadata[newKeyID] = KeyMetadata{
        KeyID:     newKeyID,
        Version:   uint64(len(m.metadata) + 1),
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
        Status:    "active",
    }

    return nil
}

// BackupKeyWithShamir 使用 Shamir's Secret Sharing 备份密钥
func (m *MasterKeyManager) BackupKeyWithShamir(keyID string, totalShares int, threshold int) ([][]byte, error) {
    key, err := m.hsmClient.ExportKey(keyID)
    if err != nil {
        return nil, err
    }

    // 使用 Shamir's Secret Sharing 分割密钥
    shares, err := shamir.Split(key, totalShares, threshold)
    if err != nil {
        return nil, err
    }

    return shares, nil
}

// RestoreKeyFromShares 从分片恢复密钥
func (m *MasterKeyManager) RestoreKeyFromShares(shares [][]byte) ([]byte, error) {
    return shamir.Combine(shares)
}
```

### 10.2 身份验证安全

```yaml
authentication_security:
  # OAuth 验证
  oauth:
    flow: "PKCE"  # Proof Key for Code Exchange
    code_challenge_method: "S256"  # SHA256
    state_parameter: true
    nonce_enabled: true

  # Token 配置
  tokens:
    access_token:
      algorithm: "RS256"
      validity: "1h"
      refresh_enabled: true
    refresh_token:
      validity: "7d"
      rotation_on_use: true
      max_lifetime: "30d"

  # 防伪造措施
  anti_forgery:
    state_parameter:
      enabled: true
      validity: "10m"
      single_use: true
    nonce:
      enabled: true
      min_length: 16
      storage: "redis"
      ttl: "1h"
```

```typescript
// services/auth/oauth-pkce.ts - OAuth PKCE 流程实现

import crypto from 'crypto';
import { v4 as uuidv4 } from 'uuid';

interface PKCEChallenge {
  code_verifier: string;
  code_challenge: string;
  code_challenge_method: 'S256' | 'plain';
}

interface OAuthState {
  state: string;
  nonce: string;
  redirect_uri: string;
  created_at: number;
}

// 生成 PKCE Code Verifier
export function generateCodeVerifier(): string {
  return crypto.randomBytes(32)
    .toString('base64url')
    .slice(0, 128);  // 43-128 字符
}

// 生成 PKCE Code Challenge
export function generateCodeChallenge(verifier: string): PKCEChallenge {
  const hash = crypto.createHash('sha256')
    .update(verifier)
    .digest('base64url');

  return {
    code_verifier: verifier,
    code_challenge: hash,
    code_challenge_method: 'S256',
  };
}

// 生成防伪造 State 和 Nonce
export function generateOAuthState(redirectUri: string): OAuthState {
  return {
    state: crypto.randomBytes(16).toString('hex'),
    nonce: uuidv4(),
    redirect_uri: redirectUri,
    created_at: Date.now(),
  };
}

// 验证 OAuth 回调
export function validateOAuthCallback(
  receivedState: string,
  receivedNonce: string,
  storedState: OAuthState,
  stateValidityMs: number = 10 * 60 * 1000,  // 10 分钟
): { valid: boolean; error?: string } {
  // 检查 state 是否匹配
  if (receivedState !== storedState.state) {
    return { valid: false, error: 'State mismatch - possible CSRF attack' };
  }

  // 检查 nonce 是否匹配
  if (receivedNonce !== storedState.nonce) {
    return { valid: false, error: 'Nonce mismatch - possible replay attack' };
  }

  // 检查 state 是否过期
  if (Date.now() - storedState.created_at > stateValidityMs) {
    return { valid: false, error: 'State expired' };
  }

  return { valid: true };
}

// Token 配置
export const TokenConfig = {
  accessToken: {
    expiresIn: '1h',
    algorithm: 'RS256',
  },
  refreshToken: {
    expiresIn: '7d',
    rotateOnUse: true,
    maxLifetime: '30d',
  },
};
```

### 10.3 交易安全

```yaml
transaction_security:
  # 重放防护
  replay_protection:
    method: "sequence_nonce"
    nonce_cache_ttl: "24h"
    max_nonce_gap: 1000

  # 双花防护
  double_spend_protection:
    mechanism: "CometBFT_consensus"
    finality: "instant"  # 即时确定性
    block_confirmation: 1

  # 原子性保证
  atomicity:
    transaction_type: "ACID"
    rollback_enabled: true
    state_commit: "atomic"

  # MEV 防护
  mev_protection:
    enabled: true
    method: "private_transaction_pool"
    flashbots_enabled: false  # Cosmos 不需要
    transaction_ordering: "fee_priority"
```

```go
// x/base/keeper/transaction_security.go - 交易安全实现

package keeper

import (
    "errors"
    "fmt"
    "time"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

// NonceManager 重放防护 - Nonce 管理
type NonceManager struct {
    store      sdk.KVStore
    cacheTTL   time.Duration
    maxGap     uint64
}

// CheckAndIncrementNonce 检查并递增 Nonce
func (n *NonceManager) CheckAndIncrementNonce(ctx sdk.Context, addr sdk.AccAddress, nonce uint64) error {
    currentNonce := n.getCurrentNonce(ctx, addr)

    // 检查 Nonce 是否有效
    if nonce < currentNonce {
        return errors.New("nonce too low - possible replay attack")
    }

    if nonce > currentNonce+n.maxGap {
        return fmt.Errorf("nonce too high - max gap is %d", n.maxGap)
    }

    // 更新 Nonce
    n.setCurrentNonce(ctx, addr, nonce+1)
    return nil
}

// DoubleSpendProtection 双花防护 - 由 CometBFT 共识保证
// Cosmos SDK 通过以下机制防止双花:
// 1. 每个账户的 sequence nonce
// 2. CometBFT 的即时确定性
// 3. UTXO 模型的余额检查

// AtomicTransaction 原子事务包装器
type AtomicTransaction struct {
    keeper     *Keeper
    operations []Operation
    rollback   []RollbackFunc
}

type Operation func(ctx sdk.Context) error
type RollbackFunc func(ctx sdk.Context)

// Execute 原子执行所有操作
func (at *AtomicTransaction) Execute(ctx sdk.Context) error {
    // 在缓存上下文中执行
    cacheCtx, writeCache := ctx.CacheContext()

    for i, op := range at.operations {
        if err := op(cacheCtx); err != nil {
            // 执行回滚
            for j := i - 1; j >= 0; j-- {
                if at.rollback[j] != nil {
                    at.rollback[j](ctx)
                }
            }
            return fmt.Errorf("operation %d failed: %w", i, err)
        }
    }

    // 原子写入
    writeCache()
    return nil
}

// MEVProtection MEV 防护配置
type MEVProtectionConfig struct {
    Enabled              bool
    PrivateTxPool        bool
    FeePriorityOrdering  bool
}

var DefaultMEVProtection = MEVProtectionConfig{
    Enabled:              true,
    PrivateTxPool:        true,   // 私有交易池
    FeePriorityOrdering:  true,   // 按手续费排序
}

// ValidateTransaction 验证交易安全性
func (k Keeper) ValidateTransaction(ctx sdk.Context, tx sdk.Tx) error {
    // 1. 重放防护检查
    for _, msg := range tx.GetMsgs() {
        if sigMsg, ok := msg.(sdk.Msg); ok {
            signers := sigMsg.GetSigners()
            for _, signer := range signers {
                // Nonce 检查在 AnteHandler 中完成
                _ = signer
            }
        }
    }

    // 2. 双花检查 (余额检查)
    // 由 bank module 的 AnteHandler 完成

    // 3. MEV 检查
    if k.mevConfig.Enabled {
        // 检查是否为私有交易
        // 检查手续费是否合理
    }

    return nil
}
```

### 10.4 模块安全审计要点

```go
// 各模块安全审计清单

// =============================================================================
// x/identity - 身份模块
// =============================================================================
// 安全要点:
// 1. Merkle 证明验证
// 2. OAuth 回调验证
// 3. 身份哈希唯一性检查
// 4. 防止身份盗用

func (k Keeper) VerifyMerkleProof(proof MerkleProof, root []byte, leaf []byte) bool {
    // Merkle 证明验证实现
    // 防止伪造证明
    computedHash := leaf
    for _, sibling := range proof.Siblings {
        if proof.IsLeftSibling {
            computedHash = hash(append(sibling, computedHash...))
        } else {
            computedHash = hash(append(computedHash, sibling...))
        }
    }
    return bytes.Equal(computedHash, root)
}

func (k Keeper) ValidateOAuthCallback(ctx sdk.Context, callback OAuthCallback) error {
    // 验证 state 参数
    if !k.ValidateState(ctx, callback.State) {
        return errors.New("invalid state parameter")
    }
    // 验证 nonce
    if !k.ValidateNonce(ctx, callback.Nonce) {
        return errors.New("invalid nonce")
    }
    // 验证签名
    if !k.VerifyOAuthSignature(callback) {
        return errors.New("invalid OAuth signature")
    }
    return nil
}

// =============================================================================
// x/mq - 德商模块
// =============================================================================
// 安全要点:
// 1. 零和校验 (德商总量恒定)
// 2. 溢出检查
// 3. 下界检查 (德商不为负)

func (k Keeper) RedistributeMQ(ctx sdk.Context, disputeID uint64) error {
    // 记录重分配前总量
    totalBefore := k.GetTotalMQ(ctx)

    // 执行德商转移
    if err := k.transferMQ(ctx, plaintiff, defendant, amount); err != nil {
        return err
    }

    // 零和校验
    totalAfter := k.GetTotalMQ(ctx)
    if totalBefore != totalAfter {
        return errors.New("MQ conservation violated - zero-sum check failed")
    }

    return nil
}

func (k Keeper) safeMQAdd(current uint64, delta int64) (uint64, error) {
    // 溢出检查
    if delta > 0 {
        if current > math.MaxUint64-uint64(delta) {
            return 0, errors.New("overflow detected")
        }
        return current + uint64(delta), nil
    } else {
        if current < uint64(-delta) {
            return 0, errors.New("underflow detected")
        }
        return current - uint64(-delta), nil
    }
}

// =============================================================================
// x/compute - 算力模块
// =============================================================================
// 安全要点:
// 1. API Key 隔离存储
// 2. 执行沙箱
// 3. 资源限制

func (k Keeper) StoreAPIKey(ctx sdk.Context, provider string, encryptedKey []byte) error {
    // API Key 按提供者隔离存储
    key := types.APIKeyStoreKey(provider)

    // 使用最新的主密钥加密
    encrypted, err := k.encryptWithCurrentKey(encryptedKey)
    if err != nil {
        return err
    }

    // 存储到隔离的存储空间
    k.store.Set(key, encrypted)
    return nil
}

type ExecutionSandbox struct {
    MaxMemoryMB    uint64  // 最大内存
    MaxCPUMillis   uint64  // 最大 CPU 时间
    MaxDuration    time.Duration  // 最大执行时间
    NetworkAccess  bool    // 网络访问
}

func (s *ExecutionSandbox) Execute(fn func() error) error {
    // 1. 设置资源限制
    // 2. 设置超时
    // 3. 在隔离环境中执行
    // 4. 清理资源
    return fn()
}

// =============================================================================
// x/escrow - 托管模块
// =============================================================================
// 安全要点:
// 1. 余额校验
// 2. 权限检查
// 3. 争议锁定

func (k Keeper) CreateEscrow(ctx sdk.Context, creator sdk.AccAddress, beneficiary sdk.AccAddress, amount sdk.Coins) error {
    // 1. 余额校验
    balance := k.bankKeeper.GetBalance(ctx, creator, amount[0].Denom)
    if balance.IsLT(amount[0]) {
        return errors.New("insufficient balance")
    }

    // 2. 锁定资金到托管账户
    if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, amount); err != nil {
        return err
    }

    // 3. 创建托管记录
    escrow := types.Escrow{
        Creator:     creator.String(),
        Beneficiary: beneficiary.String(),
        Amount:      amount,
        Status:      types.EscrowStatus_ACTIVE,
        CreatedAt:   ctx.BlockTime(),
    }

    k.SetEscrow(ctx, escrow)
    return nil
}

func (k Keeper) ReleaseEscrow(ctx sdk.Context, escrowID uint64, releaser sdk.AccAddress) error {
    escrow, _ := k.GetEscrow(ctx, escrowID)

    // 权限检查 - 只有创建者或受益人可以释放
    if releaser.String() != escrow.Creator && releaser.String() != escrow.Beneficiary {
        return errors.New("unauthorized releaser")
    }

    // 争议锁定检查
    if escrow.Status == types.EscrowStatus_DISPUTE_LOCKED {
        return errors.New("escrow is locked due to active dispute")
    }

    // 释放资金
    return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, releaser, escrow.Amount)
}

// =============================================================================
// x/dispute - 争议模块
// =============================================================================
// 安全要点:
// 1. 投票权重验证
// 2. 时间锁
// 3. 评审员资格验证

func (k Keeper) CastVote(ctx sdk.Context, disputeID uint64, juror sdk.AccAddress, verdict string) error {
    // 1. 验证评审员资格
    if !k.IsEligibleJuror(ctx, juror) {
        return errors.New("juror not eligible")
    }

    // 2. 计算投票权重 (基于德商)
    jurorMQ := k.mqKeeper.GetMQ(ctx, juror)
    weight := k.calculateVotingWeight(jurorMQ)

    // 3. 权重验证 - 防止权重篡改
    expectedWeight := k.calculateVotingWeight(jurorMQ)
    if weight != expectedWeight {
        return errors.New("vote weight manipulation detected")
    }

    // 4. 记录投票
    vote := types.Vote{
        DisputeID: disputeID,
        Juror:     juror.String(),
        Verdict:   verdict,
        Weight:    weight,
        Timestamp: ctx.BlockTime(),
    }

    k.SetVote(ctx, vote)
    return nil
}

// 时间锁 - 上诉期
func (k Keeper) CanAppeal(ctx sdk.Context, disputeID uint64) bool {
    dispute, _ := k.GetDispute(ctx, disputeID)

    // 检查是否在上诉期内 (裁决后 7 天)
    appealDeadline := dispute.ResolvedAt.Add(7 * 24 * time.Hour)
    return ctx.BlockTime().Before(appealDeadline)
}

// =============================================================================
// x/idea - 想法模块
// =============================================================================
// 安全要点:
// 1. 贡献权重算法验证
// 2. 资金目标验证
// 3. 贡献者权益保护

func (k Keeper) AddContribution(ctx sdk.Context, ideaID uint64, contributor sdk.AccAddress, contributionType string, weight sdk.Dec) error {
    // 1. 验证贡献权重算法
    if err := k.validateContributionWeight(contributionType, weight); err != nil {
        return err
    }

    // 2. 检查权重总和不超过 100%
    currentTotal := k.GetTotalContributionWeight(ctx, ideaID)
    newTotal := currentTotal.Add(weight)
    if newTotal.GT(sdk.OneDec()) {
        return errors.New("total contribution weight exceeds 100%")
    }

    // 3. 记录贡献
    contribution := types.Contribution{
        IdeaID:      ideaID,
        Contributor: contributor.String(),
        Type:        contributionType,
        Weight:      weight,
        Timestamp:   ctx.BlockTime(),
    }

    k.SetContribution(ctx, contribution)
    return nil
}

// =============================================================================
// x/task - 任务模块
// =============================================================================
// 安全要点:
// 1. 里程碑状态机完整性
// 2. 预算检查
// 3. 截止时间验证

func (k Keeper) TransitionMilestone(ctx sdk.Context, taskID uint64, milestoneID uint64, newStatus string) error {
    milestone, _ := k.GetMilestone(ctx, taskID, milestoneID)

    // 状态机转换验证
    validTransitions := map[string][]string{
        "pending":    {"submitted"},
        "submitted":  {"approved", "revision_requested"},
        "revision_requested": {"submitted"},
        "approved":   {"completed"},
    }

    allowed, exists := validTransitions[milestone.Status]
    if !exists {
        return errors.New("invalid current status")
    }

    isValid := false
    for _, s := range allowed {
        if s == newStatus {
            isValid = true
            break
        }
    }

    if !isValid {
        return fmt.Errorf("invalid transition from %s to %s", milestone.Status, newStatus)
    }

    // 更新状态
    milestone.Status = newStatus
    milestone.UpdatedAt = ctx.BlockTime()
    k.SetMilestone(ctx, milestone)

    return nil
}
```

### 10.5 数据隐私

```yaml
data_privacy:
  # 法规遵循
  compliance:
    - "GDPR"    # 欧盟通用数据保护条例
    - "CCPA"    # 加州消费者隐私法

  # 数据最小化
  data_minimization:
    enabled: true
    principles:
      - "只收集必要数据"
      - "数据用途明确"
      - "存储期限限制"
      - "定期清理过期数据"

  # 删除权
  right_to_deletion:
    enabled: true
    types:
      - "soft_delete"   # 软删除 (可恢复)
      - "hard_delete"   # 硬删除 (不可恢复)
      - "anonymization" # 匿名化
    verification_required: true

  # 日志保留
  log_retention:
    duration: "90d"
    types:
      - name: "access_logs"
        retention: "90d"
        include_pii: false
      - name: "audit_logs"
        retention: "365d"
        include_pii: true
        encrypted: true
      - name: "error_logs"
        retention: "30d"
        include_pii: false
```

```typescript
// services/privacy/data-privacy.ts - 数据隐私实现

import { v4 as uuidv4 } from 'uuid';

interface DeletionRequest {
  id: string;
  userId: string;
  type: 'soft_delete' | 'hard_delete' | 'anonymization';
  status: 'pending' | 'processing' | 'completed' | 'failed';
  createdAt: Date;
  completedAt?: Date;
  verificationCode: string;
}

interface DataCategory {
  name: string;
  retentionDays: number;
  includesPII: boolean;
  anonymizable: boolean;
}

// 数据分类配置
const DataCategories: DataCategory[] = [
  { name: 'identity_proofs', retentionDays: 365, includesPII: true, anonymizable: false },
  { name: 'transaction_history', retentionDays: 365, includesPII: true, anonymizable: true },
  { name: 'mq_scores', retentionDays: 90, includesPII: false, anonymizable: true },
  { name: 'dispute_records', retentionDays: 730, includesPII: true, anonymizable: true },
  { name: 'task_submissions', retentionDays: 365, includesPII: true, anonymizable: true },
];

// 数据最小化检查
export function validateDataCollection(dataType: string, fields: string[]): { valid: boolean; reason?: string } {
  const requiredFields: Record<string, string[]> = {
    'identity_verification': ['identity_hash', 'identity_type'],
    'transaction': ['sender', 'receiver', 'amount'],
    'task_creation': ['title', 'description', 'budget'],
  };

  const required = requiredFields[dataType] || [];
  const extraFields = fields.filter(f => !required.includes(f));

  if (extraFields.length > 0) {
    console.warn(`Extra fields collected for ${dataType}: ${extraFields.join(', ')}`);
    // 根据配置决定是否允许
  }

  return { valid: true };
}

// 处理删除请求
export class DataDeletionService {
  private requests: Map<string, DeletionRequest> = new Map();

  // 创建删除请求
  async createDeletionRequest(
    userId: string,
    type: 'soft_delete' | 'hard_delete' | 'anonymization',
  ): Promise<DeletionRequest> {
    const request: DeletionRequest = {
      id: uuidv4(),
      userId,
      type,
      status: 'pending',
      createdAt: new Date(),
      verificationCode: this.generateVerificationCode(),
    };

    this.requests.set(request.id, request);
    return request;
  }

  // 验证删除请求
  async verifyAndProcess(requestId: string, verificationCode: string): Promise<boolean> {
    const request = this.requests.get(requestId);
    if (!request || request.verificationCode !== verificationCode) {
      return false;
    }

    request.status = 'processing';

    try {
      await this.executeDeletion(request);
      request.status = 'completed';
      request.completedAt = new Date();
      return true;
    } catch (error) {
      request.status = 'failed';
      return false;
    }
  }

  // 执行删除
  private async executeDeletion(request: DeletionRequest): Promise<void> {
    const { userId, type } = request;

    switch (type) {
      case 'soft_delete':
        // 软删除 - 标记为已删除
        await this.softDeleteUserData(userId);
        break;

      case 'hard_delete':
        // 硬删除 - 永久删除
        await this.hardDeleteUserData(userId);
        break;

      case 'anonymization':
        // 匿名化 - 保留数据但移除关联
        await this.anonymizeUserData(userId);
        break;
    }
  }

  private async softDeleteUserData(userId: string): Promise<void> {
    // 标记所有用户数据为已删除
    console.log(`Soft deleting data for user ${userId}`);
  }

  private async hardDeleteUserData(userId: string): Promise<void> {
    // 永久删除所有用户数据
    // 注意: 链上数据无法删除，只能弃用
    console.log(`Hard deleting data for user ${userId}`);
  }

  private async anonymizeUserData(userId: string): Promise<void> {
    // 匿名化用户数据
    console.log(`Anonymizing data for user ${userId}`);
  }

  private generateVerificationCode(): string {
    return Math.random().toString(36).substring(2, 8).toUpperCase();
  }
}

// 日志管理
export class LogRetentionManager {
  private retentionPolicies: Map<string, number> = new Map([
    ['access_logs', 90],
    ['audit_logs', 365],
    ['error_logs', 30],
  ]);

  // 清理过期日志
  async cleanupExpiredLogs(): Promise<void> {
    const now = new Date();

    for (const [logType, retentionDays] of this.retentionPolicies) {
      const cutoffDate = new Date(now.getTime() - retentionDays * 24 * 60 * 60 * 1000);

      // 删除过期日志
      console.log(`Cleaning up ${logType} older than ${cutoffDate.toISOString()}`);
      // await this.deleteLogsOlderThan(logType, cutoffDate);
    }
  }

  // 脱敏日志中的 PII
  sanitizeLogEntry(entry: string): string {
    // 移除或脱敏 PII
    const patterns = [
      { pattern: /\b[\w.-]+@[\w.-]+\.\w+\b/g, replacement: '[EMAIL]' },
      { pattern: /\b\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}\b/g, replacement: '[CARD]' },
      { pattern: /\bshare[a-z0-9]{38}\b/g, replacement: '[ADDRESS]' },
    ];

    let sanitized = entry;
    for (const { pattern, replacement } of patterns) {
      sanitized = sanitized.replace(pattern, replacement);
    }

    return sanitized;
  }
}
```

### 10.6 共识安全

```yaml
consensus_security:
  # Slashing 条件
  slashing:
    # 双签惩罚
    double_sign:
      slash_percentage: "5%"      # 质押罚没比例
      jail_duration: "forever"    # 永久监禁
      evidence_age: "336h"        # 证据有效期 (14天)

    # 下线惩罚
    downtime:
      slash_percentage: "0.01%"   # 质押罚没比例
      jail_duration: "24h"        # 监禁时长
      min_signed_per_window: 50   # 窗口内最少签名比例 %
      signed_blocks_window: 10000 # 签名窗口大小

    # 恶意投票
    malicious_voting:
      slash_percentage: "2%"
      jail_duration: "72h"

  # 治理安全
  governance_security:
    proposal_minimum_deposit: "1000stt"  # 提案最低押金
    quorum: "20%"                        # 参与率要求
    threshold: "50%"                     # 通过阈值
    veto_threshold: "33.4%"              # 否决阈值
    voting_period: "72h"                 # 投票期

    # 紧急提案
    expedited_proposals:
      voting_period: "24h"
      threshold: "66.7%"
```

```go
// x/slashing/keeper/slashing.go - Slashing 实现

package keeper

import (
    "fmt"
    "time"

    sdk "github.com/cosmos/cosmos-sdk/types"
    slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// Slashing 配置
type SlashingConfig struct {
    DoubleSignSlashFraction   sdk.Dec  // 双签罚没比例
    DoubleSignJailDuration    time.Duration
    DowntimeSlashFraction     sdk.Dec  // 下线罚没比例
    DowntimeJailDuration      time.Duration
    SignedBlocksWindow        int64
    MinSignedPerWindow        sdk.Dec
}

var DefaultSlashingConfig = SlashingConfig{
    DoubleSignSlashFraction:   sdk.NewDecWithPrec(5, 2),    // 5%
    DoubleSignJailDuration:    time.Duration(1<<63 - 1),    // forever
    DowntimeSlashFraction:     sdk.NewDecWithPrec(1, 4),    // 0.01%
    DowntimeJailDuration:      24 * time.Hour,
    SignedBlocksWindow:        10000,
    MinSignedPerWindow:        sdk.NewDecWithPrec(50, 2),   // 50%
}

// HandleDoubleSign 处理双签
func (k Keeper) HandleDoubleSign(ctx sdk.Context, evidence *slashingtypes.Equivocation) error {
    // 1. 验证证据
    if err := k.ValidateDoubleSignEvidence(ctx, evidence); err != nil {
        return err
    }

    // 2. 获取验证者信息
    validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, evidence.GetConsensusAddress())
    if !found {
        return fmt.Errorf("validator not found")
    }

    // 3. 计算罚没金额
    slashAmount := validator.Tokens.ToDec().Mul(k.config.DoubleSignSlashFraction).TruncateInt()

    // 4. 执行罚没
    if err := k.stakingKeeper.Slash(ctx, evidence.GetConsensusAddress(), evidence.Height, slashAmount); err != nil {
        return err
    }

    // 5. 监禁验证者 (永久)
    if err := k.stakingKeeper.Jail(ctx, evidence.GetConsensusAddress()); err != nil {
        return err
    }

    // 6. 记录事件
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            slashingtypes.EventTypeSlash,
            sdk.NewAttribute("type", "double_sign"),
            sdk.NewAttribute("validator", validator.OperatorAddress),
            sdk.NewAttribute("slash_amount", slashAmount.String()),
            sdk.NewAttribute("jailed", "forever"),
        ),
    )

    return nil
}

// HandleDowntime 处理下线
func (k Keeper) HandleDowntime(ctx sdk.Context, consensusAddr sdk.ConsAddress) error {
    // 1. 检查下线窗口
    signedBlocksWindow := k.SignedBlocksWindow(ctx)
    minSignedPerWindow := k.MinSignedPerWindow(ctx)

    // 获取签名统计
    signingInfo, found := k.GetValidatorSigningInfo(ctx, consensusAddr)
    if !found {
        return fmt.Errorf("signing info not found")
    }

    // 计算签名率
    signedRatio := sdk.NewDec(signingInfo.SignedBlocksCounter).Quo(sdk.NewDec(signedBlocksWindow))

    // 2. 判断是否需要罚没
    if signedRatio.LT(minSignedPerWindow) {
        // 获取验证者
        validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, consensusAddr)
        if !found {
            return fmt.Errorf("validator not found")
        }

        // 3. 计算罚没金额
        slashAmount := validator.Tokens.ToDec().Mul(k.config.DowntimeSlashFraction).TruncateInt()

        // 4. 执行罚没
        if err := k.stakingKeeper.Slash(ctx, consensusAddr, ctx.BlockHeight(), slashAmount); err != nil {
            return err
        }

        // 5. 监禁验证者 (24小时)
        if err := k.stakingKeeper.JailUntil(ctx, consensusAddr, ctx.BlockTime().Add(k.config.DowntimeJailDuration)); err != nil {
            return err
        }

        // 6. 记录事件
        ctx.EventManager().EmitEvent(
            sdk.NewEvent(
                slashingtypes.EventTypeSlash,
                sdk.NewAttribute("type", "downtime"),
                sdk.NewAttribute("validator", validator.OperatorAddress),
                sdk.NewAttribute("slash_amount", slashAmount.String()),
                sdk.NewAttribute("signed_ratio", signedRatio.String()),
                sdk.NewAttribute("jailed_until", ctx.BlockTime().Add(k.config.DowntimeJailDuration).String()),
            ),
        )
    }

    return nil
}

// =============================================================================
// x/gov/keeper/governance.go - 治理安全
// ==================================================================

// GovernanceConfig 治理配置
type GovernanceConfig struct {
    MinDeposit       sdk.Coins    // 最低押金
    Quorum           sdk.Dec      // 参与率要求
    Threshold        sdk.Dec      // 通过阈值
    VetoThreshold    sdk.Dec      // 否决阈值
    VotingPeriod     time.Duration
}

var DefaultGovernanceConfig = GovernanceConfig{
    MinDeposit:    sdk.NewCoins(sdk.NewCoin("stt", sdk.NewInt(1000_000000))),  // 1000 STT
    Quorum:        sdk.NewDecWithPrec(20, 2),    // 20% 参与率
    Threshold:     sdk.NewDecWithPrec(50, 2),    // 50% 通过
    VetoThreshold: sdk.NewDecWithPrec(334, 3),   // 33.4% 否决
    VotingPeriod:  72 * time.Hour,
}

// TallyVotes 计票并验证
func (k Keeper) TallyVotes(ctx sdk.Context, proposalID uint64) (passes bool, tallyResult TallyResult, err error) {
    proposal, found := k.GetProposal(ctx, proposalID)
    if !found {
        return false, tallyResult, fmt.Errorf("proposal not found")
    }

    // 1. 收集所有投票
    votes := k.GetAllVotes(ctx, proposalID)

    // 2. 计算投票权重
    var yesWeight, noWeight, noWithVetoWeight, abstainWeight sdk.Dec
    var totalVotingPower sdk.Dec

    for _, vote := range votes {
        // 获取投票者的质押权重
        validator, found := k.stakingKeeper.GetDelegatorValidator(ctx, vote.Voter)
        if !found {
            continue
        }

        votingPower := validator.Tokens.ToDec()
        totalVotingPower = totalVotingPower.Add(votingPower)

        switch vote.Option {
        case OptionYes:
            yesWeight = yesWeight.Add(votingPower)
        case OptionNo:
            noWeight = noWeight.Add(votingPower)
        case OptionNoWithVeto:
            noWithVetoWeight = noWithVetoWeight.Add(votingPower)
        case OptionAbstain:
            abstainWeight = abstainWeight.Add(votingPower)
        }
    }

    // 3. 计算参与率
    totalBondedTokens := k.stakingKeeper.TotalBondedTokens(ctx).ToDec()
    participationRate := totalVotingPower.Quo(totalBondedTokens)

    // 4. 检查是否达到参与率要求
    if participationRate.LT(k.config.Quorum) {
        return false, tallyResult, nil  // 参与率不足，提案失败
    }

    // 5. 计算投票结果
    nonAbstainTotal := yesWeight.Add(noWeight).Add(noWithVetoWeight)

    // 检查否决
    vetoRatio := noWithVetoWeight.Quo(nonAbstainTotal)
    if vetoRatio.GTE(k.config.VetoThreshold) {
        tallyResult = TallyResult{
            Yes:        yesWeight,
            No:         noWeight,
            NoWithVeto: noWithVetoWeight,
            Abstain:    abstainWeight,
        }
        return false, tallyResult, nil  // 被否决
    }

    // 检查通过
    yesRatio := yesWeight.Quo(nonAbstainTotal)
    passes = yesRatio.GTE(k.config.Threshold)

    tallyResult = TallyResult{
        Yes:        yesWeight,
        No:         noWeight,
        NoWithVeto: noWithVetoWeight,
        Abstain:    abstainWeight,
    }

    return passes, tallyResult, nil
}
```

---

## 附录

### A. 配置参数汇总

#### A.1 核心模块参数

| 参数 | 值 | 说明 |
|------|-----|------|
| 出块时间 | 10s | CometBFT 共识 |
| 最大交易大小 | 1MB | 单笔交易上限 |
| 最大 Gas | 10,000,000 | 区块 Gas 上限 |
| 最小验证者质押 | 1,000,000 STT | 100万 STT |
| 最大验证者数 | 21 | 出块验证者 |

#### A.2 服务市场参数

| 参数 | 值 | 说明 |
|------|-----|------|
| LLM 最大超时 | 3600s | 1小时 |
| Agent 最大执行时间 | 86400s | 24小时 |
| Workflow 最大执行时间 | 604800s | 7天 |
| 最小服务价格 | 0.01 STT | 服务最低定价 |
| 最大并发请求 | 100 | 每提供者最大并发 |
| 匹配超时 | 30s | 服务匹配等待时间 |

#### A.3 德商系统参数

| 参数 | 值 | 说明 |
|------|-----|------|
| 初始德商 | 100 | 新用户初始 MQ |
| 最低德商 | 0 | MQ 下限（无限趋近） |
| 最大风险比例 | 3% | 每次评分最大损失 |
| 基准偏差 λ | 1.5 | 奖惩分界线 |
| 最大偏差 max_d | 6.0 | 评分范围(-10~10) |

#### A.4 评审义务参数

| 参数 | 值 | 说明 |
|------|-----|------|
| 缺席代币惩罚 | 10 STT | 缺席扣代币 |
| 缺席德商惩罚 | 5 | 缺席扣德商 |
| 最大缺席次数 | 3 | 触发额外惩罚 |
| 多次缺席惩罚 | 10 MQ | 额外扣德商 |
| 暂停阈值 | 5次 | 暂停评审资格 |

#### A.4 争议仲裁参数

| 参数 | 值 | 说明 |
|------|-----|------|
| 协商期 | 7天 | 争议协商 |
| 投票期 | 3天 | 评审投票 |
| 上诉期 | 7天 | 裁决上诉 |
| 通过阈值 | 60% | 投票多数 |
| 评审超时 | 3天 | 自动通过 |

### B. 技术栈版本

#### B.1 核心技术栈

| 组件 | 版本 | 说明 |
|------|------|------|
| Cosmos SDK | v0.47+ | 区块链框架 |
| CometBFT | v0.37+ | 共识引擎 |
| Go | 1.21+ | 链上开发语言 |
| Libp2p | v0.30+ | P2P 通信 |

#### B.2 服务端技术栈

| 组件 | 版本 | 说明 |
|------|------|------|
| TypeScript | 5.0+ | 插件开发语言 |
| Node.js | 20+ | 运行时 |
| Redis | 7+ | 缓存服务 |
| InfluxDB | 2.7+ | 指标存储 |
| Elasticsearch | 8.11+ | 日志存储 |

#### B.3 插件技术栈

| 组件 | 版本 | 说明 |
|------|------|------|
| React | 18+ | UI 框架 (小灯) |
| OpenFang | latest | Agent 执行器 |

---

## 11. 补充细节

### 11.1 用户等级体系映射

```
用户等级体系：
- 身份验证等级：none → email → phone → KYC
- 德商等级：Newcomer → Member → Trusted → Expert → Guardian
- 信誉等级：0-100分

关系：
1. 身份验证等级影响交易限额（未验证 100 STT/天，已验证 10000 STT/天）
2. 德商等级影响投票权重和评审资格
3. 信誉等级影响风控策略（低于 50 触发人工审核）
```

```go
// x/identity/types/user_levels.go - 用户等级定义

package types

// IdentityLevel 身份验证等级
type IdentityLevel int

const (
    IdentityLevelNone   IdentityLevel = iota  // 未验证
    IdentityLevelEmail                         // 邮箱验证
    IdentityLevelPhone                         // 手机验证
    IdentityLevelKYC                           // KYC 验证
)

// MQTier 德商等级
type MQTier string

const (
    MQTierNewcomer MQTier = "Newcomer"   // 新手: MQ 10-49
    MQTierMember   MQTier = "Member"     // 成员: MQ 50-99
    MQTierTrusted  MQTier = "Trusted"    // 信任: MQ 100-299
    MQTierExpert   MQTier = "Expert"     // 专家: MQ 300-499
    MQTierGuardian MQTier = "Guardian"   // 守护者: MQ 500+
)

// ReputationScore 信誉分 (0-100)
type ReputationScore uint8

// UserLevelProfile 用户等级档案
type UserLevelProfile struct {
    IdentityLevel    IdentityLevel    `json:"identity_level"`
    MQTier           MQTier           `json:"mq_tier"`
    ReputationScore  ReputationScore  `json:"reputation_score"`
}

// GetMQTier 根据 MQ 值获取德商等级
func GetMQTier(mq uint64) MQTier {
    switch {
    case mq >= 500:
        return MQTierGuardian
    case mq >= 300:
        return MQTierExpert
    case mq >= 100:
        return MQTierTrusted
    case mq >= 50:
        return MQTierMember
    default:
        return MQTierNewcomer
    }
}

// TransactionLimits 交易限额配置
type TransactionLimits struct {
    DailyLimit  sdk.Int  // 每日限额
    TxLimit     sdk.Int  // 单笔限额
}

// 身份验证等级对应的交易限额
var IdentityTransactionLimits = map[IdentityLevel]TransactionLimits{
    IdentityLevelNone: {
        DailyLimit: sdk.NewInt(100_000000),     // 100 STT/天
        TxLimit:    sdk.NewInt(50_000000),      // 50 STT/笔
    },
    IdentityLevelEmail: {
        DailyLimit: sdk.NewInt(500_000000),     // 500 STT/天
        TxLimit:    sdk.NewInt(200_000000),     // 200 STT/笔
    },
    IdentityLevelPhone: {
        DailyLimit: sdk.NewInt(2000_000000),    // 2000 STT/天
        TxLimit:    sdk.NewInt(1000_000000),    // 1000 STT/笔
    },
    IdentityLevelKYC: {
        DailyLimit: sdk.NewInt(10000_000000),   // 10000 STT/天
        TxLimit:    sdk.NewInt(5000_000000),    // 5000 STT/笔
    },
}

// 德商等级对应的特权
var MQTierPrivileges = map[MQTier][]string{
    MQTierNewcomer: {
        "basic_tasks",        // 可接基础任务
    },
    MQTierMember: {
        "basic_tasks",
        "jury_duty",          // 可担任评审
        "dispute_creation",   // 可发起争议
    },
    MQTierTrusted: {
        "basic_tasks",
        "jury_duty",
        "dispute_creation",
        "premium_tasks",      // 可接高级任务
        "reduced_fees",       // 手续费减免 10%
    },
    MQTierExpert: {
        "basic_tasks",
        "jury_duty",
        "dispute_creation",
        "premium_tasks",
        "reduced_fees",       // 手续费减免 20%
        "lead_juror",         // 可担任首席评审
        "api_access",         // API 高级访问
    },
    MQTierGuardian: {
        "basic_tasks",
        "jury_duty",
        "dispute_creation",
        "premium_tasks",
        "reduced_fees",       // 手续费减免 30%
        "lead_juror",
        "api_access",
        "governance_vote",    // 治理投票权重加成
        "priority_support",   // 优先支持
    },
}

// 风控策略
type RiskControlStrategy struct {
    Threshold      ReputationScore  // 触发阈值
    Action         string           // 采取行动
    ReviewRequired bool             // 是否需要人工审核
}

var RiskControlStrategies = []RiskControlStrategy{
    {Threshold: 50, Action: "manual_review", ReviewRequired: true},   // < 50: 人工审核
    {Threshold: 70, Action: "enhanced_monitoring", ReviewRequired: false}, // 50-70: 加强监控
    {Threshold: 85, Action: "standard", ReviewRequired: false},       // 70-85: 标准流程
    {Threshold: 100, Action: "fast_track", ReviewRequired: false},    // 85+: 快速通道
}

// CheckTransactionLimit 检查交易限额
func CheckTransactionLimit(identityLevel IdentityLevel, dailyUsed sdk.Int, txAmount sdk.Int) error {
    limits, exists := IdentityTransactionLimits[identityLevel]
    if !exists {
        limits = IdentityTransactionLimits[IdentityLevelNone]
    }

    // 检查单笔限额
    if txAmount.GT(limits.TxLimit) {
        return fmt.Errorf("transaction amount %s exceeds single tx limit %s",
            txAmount.String(), limits.TxLimit.String())
    }

    // 检查每日限额
    totalDaily := dailyUsed.Add(txAmount)
    if totalDaily.GT(limits.DailyLimit) {
        return fmt.Errorf("daily limit %s would be exceeded (used: %s, tx: %s)",
            limits.DailyLimit.String(), dailyUsed.String(), txAmount.String())
    }

    return nil
}

// CheckRiskControl 检查风控策略
func CheckRiskControl(reputation ReputationScore) RiskControlStrategy {
    for _, strategy := range RiskControlStrategies {
        if reputation < strategy.Threshold {
            return strategy
        }
    }
    return RiskControlStrategies[len(RiskControlStrategies)-1]
}
```

### 11.2 网络分区处理

```
网络分区处理策略：
1. 检测机制：
   - 心跳超时 > 10s
   - 区块高度差 > 100

2. 分区行为：
   - 小分区（< 1/3 验证者）：暂停出块
   - 大分区（≥ 2/3 验证者）：继续运行

3. 恢复流程：
   - 同步缺失区块
   - 验证分区期间交易
   - 更新世界状态
```

```go
// x/consensus/keeper/partition.go - 网络分区处理

package keeper

import (
    "context"
    "fmt"
    "time"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

// PartitionConfig 网络分区配置
type PartitionConfig struct {
    // 检测参数
    HeartbeatTimeout    time.Duration  // 心跳超时 (默认 10s)
    MaxHeightDiff       int64          // 最大区块高度差 (默认 100)

    // 分区阈值
    MinorPartitionRatio sdk.Dec        // 小分区阈值 (< 1/3)
    MajorPartitionRatio sdk.Dec        // 大分区阈值 (≥ 2/3)

    // 恢复参数
    SyncBatchSize       int64          // 同步批次大小
    MaxPendingBlocks    int64          // 最大待处理区块数
}

var DefaultPartitionConfig = PartitionConfig{
    HeartbeatTimeout:    10 * time.Second,
    MaxHeightDiff:       100,
    MinorPartitionRatio: sdk.NewDecWithPrec(33, 2),  // 33%
    MajorPartitionRatio: sdk.NewDecWithPrec(67, 2),  // 67%
    SyncBatchSize:       50,
    MaxPendingBlocks:    1000,
}

// PartitionState 分区状态
type PartitionState struct {
    IsPartitioned       bool            `json:"is_partitioned"`
    DetectedAt          time.Time       `json:"detected_at"`
    PartitionType       PartitionType   `json:"partition_type"`
    LocalVotingPower    sdk.Dec         `json:"local_voting_power"`
    TotalVotingPower    sdk.Dec         `json:"total_voting_power"`
    HeightAtDetection   int64           `json:"height_at_detection"`
}

type PartitionType string

const (
    PartitionTypeNone   PartitionType = "none"
    PartitionTypeMinor  PartitionType = "minor"   // 小分区 (< 1/3)
    PartitionTypeMajor  PartitionType = "major"   // 大分区 (≥ 2/3)
)

// PartitionDetector 分区检测器
type PartitionDetector struct {
    config           PartitionConfig
    lastHeartbeat    map[string]time.Time  // 验证者 -> 最后心跳时间
    knownHeights     map[string]int64      // 验证者 -> 已知高度
    localState       PartitionState
}

// DetectPartition 检测网络分区
func (d *PartitionDetector) DetectPartition(
    ctx sdk.Context,
    currentHeight int64,
    validators []ValidatorInfo,
) (PartitionState, error) {
    state := PartitionState{
        IsPartitioned: false,
        PartitionType: PartitionTypeNone,
    }

    now := time.Now()
    var totalPower, localPower sdk.Dec
    totalPower = sdk.ZeroDec()
    localPower = sdk.ZeroDec()

    for _, val := range validators {
        totalPower = totalPower.Add(val.VotingPower)

        // 检查心跳超时
        lastHeartbeat, exists := d.lastHeartbeat[val.Address]
        if !exists || now.Sub(lastHeartbeat) > d.config.HeartbeatTimeout {
            // 心跳超时 - 可能分区
            continue
        }

        // 检查区块高度差
        knownHeight, exists := d.knownHeights[val.Address]
        if exists && abs(currentHeight-knownHeight) > d.config.MaxHeightDiff {
            // 高度差过大 - 可能分区
            continue
        }

        // 该验证者正常通信
        localPower = localPower.Add(val.VotingPower)
    }

    // 计算本地分区占比
    if totalPower.IsZero() {
        return state, fmt.Errorf("no voting power available")
    }

    state.LocalVotingPower = localPower
    state.TotalVotingPower = totalPower
    localRatio := localPower.Quo(totalPower)

    // 判断分区类型
    if localRatio.LT(d.config.MinorPartitionRatio) {
        // 小分区 (< 1/3) - 暂停出块
        state.IsPartitioned = true
        state.PartitionType = PartitionTypeMinor
        state.DetectedAt = now
        state.HeightAtDetection = currentHeight
    } else if localRatio.GTE(d.config.MajorPartitionRatio) {
        // 大分区 (≥ 2/3) - 继续运行
        state.IsPartitioned = false
        state.PartitionType = PartitionTypeNone
    } else {
        // 中间状态 (1/3 - 2/3) - 警告但不暂停
        state.IsPartitioned = false
        state.PartitionType = PartitionTypeNone
        // 可以触发警告事件
    }

    d.localState = state
    return state, nil
}

// HandlePartitionBehavior 处理分区行为
func (k Keeper) HandlePartitionBehavior(ctx sdk.Context, state PartitionState) error {
    switch state.PartitionType {
    case PartitionTypeMinor:
        // 小分区行为: 暂停出块
        return k.pauseBlockProduction(ctx, state)

    case PartitionTypeMajor:
        // 大分区行为: 继续运行
        return k.continueOperation(ctx)

    default:
        return nil
    }
}

// pauseBlockProduction 暂停出块
func (k Keeper) pauseBlockProduction(ctx sdk.Context, state PartitionState) error {
    // 1. 停止共识参与
    k.SetConsensusPaused(ctx, true)

    // 2. 记录分区状态
    k.SetPartitionState(ctx, state)

    // 3. 发出分区事件
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "partition_detected",
            sdk.NewAttribute("type", string(state.PartitionType)),
            sdk.NewAttribute("local_power", state.LocalVotingPower.String()),
            sdk.NewAttribute("total_power", state.TotalVotingPower.String()),
            sdk.NewAttribute("height", fmt.Sprintf("%d", state.HeightAtDetection)),
        ),
    )

    return nil
}

// PartitionRecovery 分区恢复处理器
type PartitionRecovery struct {
    config     PartitionConfig
    syncBuffer []*BlockData
}

// RecoverFromPartition 从分区恢复
func (r *PartitionRecovery) RecoverFromPartition(
    ctx sdk.Context,
    fromHeight int64,
    toHeight int64,
    peerProvider PeerProvider,
) error {
    // 1. 同步缺失区块
    if err := r.syncMissingBlocks(ctx, fromHeight, toHeight, peerProvider); err != nil {
        return fmt.Errorf("failed to sync blocks: %w", err)
    }

    // 2. 验证分区期间交易
    if err := r.validatePartitionTransactions(ctx, fromHeight, toHeight); err != nil {
        return fmt.Errorf("transaction validation failed: %w", err)
    }

    // 3. 更新世界状态
    if err := r.updateWorldState(ctx, toHeight); err != nil {
        return fmt.Errorf("state update failed: %w", err)
    }

    // 4. 恢复共识参与
    // 由上层调用 SetConsensusPaused(false)

    return nil
}

// syncMissingBlocks 同步缺失区块
func (r *PartitionRecovery) syncMissingBlocks(
    ctx sdk.Context,
    fromHeight int64,
    toHeight int64,
    peerProvider PeerProvider,
) error {
    currentHeight := fromHeight

    for currentHeight <= toHeight {
        // 批量获取区块
        endHeight := min(currentHeight+r.config.SyncBatchSize-1, toHeight)
        blocks, err := peerProvider.GetBlocks(ctx, currentHeight, endHeight)
        if err != nil {
            return err
        }

        for _, block := range blocks {
            // 验证区块签名
            if err := r.validateBlockSignatures(block); err != nil {
                return fmt.Errorf("invalid block at height %d: %w", block.Height, err)
            }

            r.syncBuffer = append(r.syncBuffer, block)
        }

        currentHeight = endHeight + 1
    }

    return nil
}

// validatePartitionTransactions 验证分区期间交易
func (r *PartitionRecovery) validatePartitionTransactions(
    ctx sdk.Context,
    fromHeight int64,
    toHeight int64,
) error {
    for _, block := range r.syncBuffer {
        if block.Height < fromHeight || block.Height > toHeight {
            continue
        }

        for _, tx := range block.Transactions {
            // 验证交易签名
            if err := r.validateTransactionSignature(tx); err != nil {
                // 记录无效交易，但不中断恢复
                ctx.EventManager().EmitEvent(
                    sdk.NewEvent(
                        "invalid_partition_tx",
                        sdk.NewAttribute("tx_hash", tx.Hash),
                        sdk.NewAttribute("reason", err.Error()),
                    ),
                )
                continue
            }

            // 验证交易 nonce (防重放)
            if err := r.validateTransactionNonce(ctx, tx); err != nil {
                continue
            }
        }
    }

    return nil
}

// updateWorldState 更新世界状态
func (r *PartitionRecovery) updateWorldState(ctx sdk.Context, targetHeight int64) error {
    // 按顺序应用区块状态变更
    for _, block := range r.syncBuffer {
        if err := r.applyBlockStateChanges(ctx, block); err != nil {
            return fmt.Errorf("failed to apply state at height %d: %w", block.Height, err)
        }
    }

    // 清理缓冲区
    r.syncBuffer = nil

    return nil
}

// PeerProvider 对端提供者接口
type PeerProvider interface {
    GetBlocks(ctx context.Context, fromHeight, toHeight int64) ([]*BlockData, error)
}

// BlockData 区块数据
type BlockData struct {
    Height       int64
    Hash         string
    Transactions []TransactionData
    Signatures   []ValidatorSignature
    AppStateHash []byte
}

// ValidatorInfo 验证者信息
type ValidatorInfo struct {
    Address      string
    VotingPower  sdk.Dec
    IsOnline     bool
    LastSeen     time.Time
}

// Helper functions

func abs(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}
```

### 11.3 评审义务与缺席惩罚

```
评审义务规则：
- 参与公正裁决是每位用户的义务
- 被抽中后必须在规定时间内完成评分
- 缺席将受到代币和德商的双重惩罚
- 多次缺席将受到更严厉的惩罚
```

```go
// x/dispute/keeper/jury_duty.go - 评审义务实现

package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// JuryDutyConfig 评审义务配置
type JuryDutyConfig struct {
    // 缺席惩罚
    AbsenceTokenPenalty     sdk.Coin  // 扣代币
    AbsenceMQPenalty        uint64    // 扣德商

    // 多次缺席
    MaxAbsences             uint64    // 最大缺席次数
    MultipleAbsencePenalty  uint64    // 多次缺席额外惩罚
    SuspensionThreshold     uint64    // 暂停评审资格阈值
}

var DefaultJuryDutyConfig = JuryDutyConfig{
    AbsenceTokenPenalty:     sdk.NewCoin("stt", sdk.NewInt(10_000000)), // 10 STT
    AbsenceMQPenalty:        5,  // 扣5点德商
    MaxAbsences:             3,  // 最多3次缺席
    MultipleAbsencePenalty:  10, // 多次缺席额外扣10点
    SuspensionThreshold:     5,  // 5次缺席后暂停评审资格
}

// ApplyAbsencePenalty 处理缺席惩罚
func (k Keeper) ApplyAbsencePenalty(ctx sdk.Context, juror sdk.AccAddress) error {
    config := k.GetJuryDutyConfig(ctx)

    // 记录缺席
    absenceCount := k.RecordAbsence(ctx, juror)

    // 扣代币
    if err := k.bankKeeper.SendCoinsFromAccountToModule(
        ctx, juror, types.ModuleName, sdk.NewCoins(config.AbsenceTokenPenalty),
    ); err != nil {
        // 代币不足时只扣德商
        ctx.Logger().Info("juror has insufficient tokens for absence penalty", "address", juror)
    }

    // 计算德商惩罚
    mqPenalty := config.AbsenceMQPenalty

    // 多次缺席额外惩罚
    if absenceCount >= config.MaxAbsences {
        mqPenalty += config.MultipleAbsencePenalty
    }

    // 扣德商
    currentMQ := k.GetMQ(ctx, juror)
    newMQ := uint64(max(0, int64(currentMQ)-int64(mqPenalty)))
    k.SetMQ(ctx, juror, newMQ)

    // 检查是否暂停评审资格
    if absenceCount >= config.SuspensionThreshold {
        k.SuspendJuryEligibility(ctx, juror, 30*24*time.Hour) // 暂停30天
    }

    // 发出事件
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "jury_absence",
            sdk.NewAttribute("address", juror.String()),
            sdk.NewAttribute("absence_count", fmt.Sprintf("%d", absenceCount)),
            sdk.NewAttribute("mq_penalty", fmt.Sprintf("%d", mqPenalty)),
            sdk.NewAttribute("token_penalty", config.AbsenceTokenPenalty.String()),
        ),
    )

    return nil
}

// RecordAbsence 记录缺席并返回累计次数
func (k Keeper) RecordAbsence(ctx sdk.Context, juror sdk.AccAddress) uint64 {
    store := ctx.KVStore(k.storeKey)
    key := types.AbsenceCountKey(juror)

    var count uint64 = 1
    if bz := store.Get(key); bz != nil {
        count = sdk.BigEndianToUint64(bz) + 1
    }

    store.Set(key, sdk.Uint64ToBigEndian(count))
    return count
}

// SuspendJuryEligibility 暂停评审资格
func (k Keeper) SuspendJuryEligibility(ctx sdk.Context, juror sdk.AccAddress, duration time.Duration) {
    store := ctx.KVStore(k.storeKey)
    suspendedUntil := ctx.BlockTime().Add(duration)
    key := types.JurySuspensionKey(juror)
    store.Set(key, sdk.FormatTimeBytes(suspendedUntil))
}

// IsJuryEligible 检查是否有评审资格
func (k Keeper) IsJuryEligible(ctx sdk.Context, address sdk.AccAddress) bool {
    // 检查德商是否足够
    mq := k.GetMQ(ctx, address)
    if mq < MinMQForJury { // 50
        return false
    }

    // 检查是否被暂停
    store := ctx.KVStore(k.storeKey)
    key := types.JurySuspensionKey(address)
    if bz := store.Get(key); bz != nil {
        suspendedUntil, _ := sdk.ParseTimeBytes(bz)
        if ctx.BlockTime().Before(suspendedUntil) {
            return false
        }
        // 暂停期满，清除记录
        store.Delete(key)
    }

    return true
}
```

---

*文档版本: 1.2*
*更新日期: 2026-03-02*
*维护者: ShareTokens 开发团队*
