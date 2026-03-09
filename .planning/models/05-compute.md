# x/compute - Level 1 LLM API 服务模块

> **模块类型:** 链上模块 (Cosmos SDK)
> **技术栈:** Go
> **位置:** `src/chain/x/compute`
> **所属:** 服务市场 - Level 1 服务层

---

## 模块定位

### 服务市场架构

x/compute 是 **服务市场** 的 **Level 1 服务层**，负责最基础的 LLM API 服务。

```
+-----------------------------------------------------------------------------+
|                          服务市场 (Service Market)                           |
+-----------------------------------------------------------------------------+
|                                                                             |
|  +-----------------------------------------------------------------------+  |
|  | Level 3: Workflow 服务                                                 |  |
|  | - 多步骤工作流编排                                                      |  |
|  | - 跨服务协调                                                           |  |
|  | - 复杂业务流程自动化                                                    |  |
|  +-----------------------------------------------------------------------+  |
|                                    |                                        |
|  +-----------------------------------------------------------------------+  |
|  | Level 2: Agent 服务                                                    |  |
|  | - AI Agent 托管与执行                                                  |  |
|  | - 自主决策与任务执行                                                    |  |
|  | - 多 Agent 协作                                                        |  |
|  +-----------------------------------------------------------------------+  |
|                                    |                                        |
|  +-----------------------------------------------------------------------+  |
|  | Level 1: LLM API 服务 (x/compute)  <-- 本章                            |  |
|  | - API Key 托管                                                         |  |
|  | - LLM 请求/响应处理                                                    |  |
|  | - 基础计费与托管                                                        |  |
|  +-----------------------------------------------------------------------+  |
|                                                                             |
+-----------------------------------------------------------------------------+
```

### 三层服务对比

| 层级 | 服务类型 | 核心能力 | 实现模块 | 典型场景 |
|------|----------|----------|----------|----------|
| Level 1 | LLM API | 基础模型调用、Token计费 | x/compute | GPT-4调用、Claude调用 |
| Level 2 | Agent | 自主任务执行、工具调用 | x/compute (扩展) | 代码生成、数据分析 |
| Level 3 | Workflow | 多步骤编排、状态管理 | x/compute (扩展) | 软件开发流程、内容创作流程 |

### x/compute 在 Level 1 的职责

1. **API Key 托管** - 安全存储和管理 LLM Provider 的 API Key
2. **请求路由** - 将用户请求路由到合适的 Provider
3. **响应处理** - 处理 LLM 响应并验证结果
4. **基础计费** - 基于 Token 使用量进行计费
5. **资金托管** - 通过 x/escrow 实现交易安全保障

---

## 模块概述

x/compute 是 ShareTokens 链上的 Level 1 LLM API 服务模块，作为服务市场的基础层，负责 API Key 托管、LLM 请求/响应处理和资金托管。作为 Cosmos SDK 模块实现，与其他链上模块通过 ABCI 接口交互。

---

## 服务交易状态机

```
+-----------------------------------------------------------------------------+
|                     Level 1 服务交易状态机                                    |
+-----------------------------------------------------------------------------+

                            +----------+
                            | pending  |  等待匹配
                            +----+-----+
                                 |
                    SubmitRequest|
                                 |
                                 v
                            +----------+
                     MatchProvider   | matched  |  已匹配提供者
                      +------------->+----+-----+
                      |                   |
                      |        StartExecution
                      |                   |
                      |                   v
                      |              +----------+
                      |              |executing |  执行中
                      |              +----+-----+
                      |                   |
                      |    +--------------+--------------+
                      |    |              |              |
                      |    |SubmitResponse|   Timeout    |CreateDispute
                      |    |              |              |
                      |    v              v              v
                      | +----------+ +----------+ +----------+
                      | |verifying | |  failed  | | disputed |
                      | +----+-----+ +----------+ +----------+
                      |      |
                      |      +--------------+
                      |      |              |
                      |VerifySuccess  VerifyFail
                      |      |              |
                      |      v              v
                      | +----------+ +----------+
                      | |completed | |  failed  |
                      | +----------+ +----------+
                      |
                      +----------------------------------------->

状态转换表：
+-----------------+-------------------+----------------------------------+
| 当前状态        | 触发事件           | 目标状态                         |
+-----------------+-------------------+----------------------------------+
| (初始)          | SubmitRequest     | pending                          |
| pending         | MatchProvider     | matched                          |
| matched         | StartExecution    | executing                        |
| executing       | SubmitResponse    | verifying                        |
| executing       | Timeout           | failed                           |
| verifying       | VerifySuccess     | completed                        |
| verifying       | VerifyFail        | failed                           |
| verifying       | CreateDispute     | disputed                         |
+-----------------+-------------------+----------------------------------+
```

---

## Cosmos SDK 集成

```
x/compute/
+-- module.go              # 模块定义
+-- keeper/
|   +-- keeper.go          # 状态管理
|   +-- grpc_query.go      # gRPC 查询
|   +-- msg_server.go      # 交易处理
+-- types/
|   +-- keys.go            # 存储 key
|   +-- types.go           # 类型定义
|   +-- genesis.go         # 创世状态
|   +-- msgs.go            # 消息类型
+-- client/
    +-- cli/               # CLI 命令
```

---

## API Key 托管

```go
// types/types.go

type APIKeyEntry struct {
    Id        uint64
    Owner     sdk.AccAddress
    Provider  string  // "openai" | "anthropic" | "google" | "azure" | "local"
    Models    []string
    EncryptedKey []byte

    AccessControl AccessControl `yaml:"access_control"`
    Pricing       PricingConfig
    Status        string  // "active" | "paused" | "revoked"
    Stats         APIKeyStats
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type AccessControl struct {
    MaxDailySpend     *sdk.Coins
    MaxRequestsPerHour *uint64
    AllowedCallers    []sdk.AccAddress
    BlockedCallers    []sdk.AccAddress
}

type PricingConfig struct {
    BasePricePerToken sdk.Coin
    ModelMultipliers  map[string]sdk.Dec
    DynamicPricing    *DynamicPricingConfig `yaml:"dynamic_pricing"`
}

type DynamicPricingConfig struct {
    Enabled         bool
    DemandMultiplier sdk.Dec
    MinPrice        sdk.Coin
    MaxPrice        sdk.Coin
}

type APIKeyStats struct {
    TotalRequests   uint64
    TotalTokensUsed uint64
    TotalRevenue    sdk.Coins
    LastUsedAt      *time.Time
}
```

---

## LLM 服务请求

```go
// types/msgs.go

type MsgSubmitComputeRequest struct {
    Creator    sdk.AccAddress
    Model      string
    Prompt     string
    PromptHash []byte
    Nonce      uint64
    Signature  []byte
    PriceOffer sdk.Coin
    Timeout    time.Duration
}

type ComputeRequest struct {
    Id         uint64
    Requester  sdk.AccAddress
    Model      string
    Prompt     string
    PromptHash []byte
    Nonce      uint64
    Signature  []byte
    PriceOffer sdk.Coin
    Timeout    time.Duration
    Status     string  // "pending" | "matched" | "executing" | "completed" | "failed"
    CreatedAt  time.Time
}
```

---

## LLM 服务响应

```go
type MsgSubmitComputeResponse struct {
    Creator    sdk.AccAddress
    RequestId  uint64
    Provider   sdk.AccAddress
    Result     string
    TokensUsed TokenUsage
    ActualCost sdk.Coin
    ResultHash []byte
    Signature  []byte
}

type ComputeResponse struct {
    Id          uint64
    RequestId   uint64
    Provider    sdk.AccAddress
    Result      string
    TokensUsed  TokenUsage
    ActualCost  sdk.Coin
    ResultHash  []byte
    Signature   []byte
    CreatedAt   time.Time
}

type TokenUsage struct {
    Input  uint64
    Output uint64
    Total  uint64
}
```

---

## 资金托管 (Escrow)

```go
// 与 x/escrow 模块集成

type EscrowStatus string

const (
    EscrowStatusPending   EscrowStatus = "pending"
    EscrowStatusLocked    EscrowStatus = "locked"
    EscrowStatusPartial   EscrowStatus = "partial"
    EscrowStatusReleased  EscrowStatus = "released"
    EscrowStatusRefunded  EscrowStatus = "refunded"
    EscrowStatusDisputed  EscrowStatus = "disputed"
)

type Escrow struct {
    Id               uint64
    TaskId           *uint64  // 关联任务（可选）
    ComputeRequestId *uint64  // 关联 LLM 请求（可选）
    Depositor        sdk.AccAddress
    Beneficiary      sdk.AccAddress
    TotalAmount      sdk.Coins
    ReleasedAmount   sdk.Coins
    Status           EscrowStatus
    ReleaseConditions []ReleaseCondition
    Transactions     []EscrowTransaction
    Arbitrator       *sdk.AccAddress
    DisputeReason    *string
    CreatedAt        time.Time
    UpdatedAt        time.Time
    ReleasedAt       *time.Time
}

type ReleaseCondition struct {
    Type      string  // "task_completion" | "compute_delivery" | "time_lock" | "multi_sig"
    Params    map[string]interface{}
    Fulfilled bool
}

type EscrowTransaction struct {
    Id        uint64
    EscrowId  uint64
    Type      string  // "deposit" | "release" | "refund" | "partial_release"
    Amount    sdk.Coins
    From      sdk.AccAddress
    To        sdk.AccAddress
    Signature []byte
    Timestamp time.Time
}
```

---

## Level 1 服务市场

```go
type ComputeOffer struct {
    Id           uint64
    Provider     sdk.AccAddress
    APIKeyEntryId uint64
    Models       []ModelOffer
    SLA          ServiceLevelAgreement
    Rating       ProviderRating
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type ModelOffer struct {
    Name            string
    Available       bool
    QueueLength     uint64
    AvgResponseTime time.Duration
    PricePerToken   sdk.Coin
}

type ServiceLevelAgreement struct {
    Uptime     sdk.Dec  // 百分比，如 0.999
    AvgLatency time.Duration
    MaxLatency time.Duration
}

type ProviderRating struct {
    Score       sdk.Dec  // 1-5
    TotalReviews uint64
}
```

---

## Keeper 接口

```go
// keeper/keeper.go

type Keeper struct {
    storeKey      sdk.StoreKey
    cdc           codec.BinaryCodec
    bankKeeper    bankkeeper.Keeper
    escrowKeeper  escrowkeeper.Keeper
    identityKeeper identitykeeper.Keeper
}

func (k Keeper) RegisterAPIKey(ctx sdk.Context, entry types.APIKeyEntry) error
func (k Keeper) SubmitRequest(ctx sdk.Context, req types.MsgSubmitComputeRequest) (uint64, error)
func (k Keeper) SubmitResponse(ctx sdk.Context, res types.MsgSubmitComputeResponse) error
func (k Keeper) GetComputeOffer(ctx sdk.Context, id uint64) (types.ComputeOffer, error)
func (k Keeper) MatchRequest(ctx sdk.Context, reqId uint64) ([]types.ComputeOffer, error)
func (k Keeper) CreateEscrow(ctx sdk.Context, escrow types.Escrow) error
func (k Keeper) ReleaseEscrow(ctx sdk.Context, escrowId uint64, to sdk.AccAddress) error
```

---

## gRPC 查询

```protobuf
// query.proto

service Query {
    rpc APIKey(QueryAPIKeyRequest) returns (QueryAPIKeyResponse);
    rpc APIKeysByOwner(QueryAPIKeysByOwnerRequest) returns (QueryAPIKeysByOwnerResponse);
    rpc ComputeRequest(QueryComputeRequestRequest) returns (QueryComputeRequestResponse);
    rpc ComputeOffers(QueryComputeOffersRequest) returns (QueryComputeOffersResponse);
    rpc Escrow(QueryEscrowRequest) returns (QueryEscrowResponse);
}
```

---

## 模块依赖

```
x/compute (Level 1 LLM API 服务)
+-- x/identity  (身份验证 - Provider/Consumer 身份)
+-- x/escrow    (资金托管 - 交易安全保障)
+-- x/dispute   (Trust System - 信任加权匹配)
+-- x/bank      (代币转账 - STT 支付)
```

---

## 与链下服务交互

```
链下 AI Service                        链上 x/compute (Level 1)
      |                                      |
      |  1. 监听 ComputeRequest 事件          |
      |<-------------------------------------|
      |                                      |
      |  2. 执行 LLM 推理                     |
      |                                      |
      |  3. 提交 ComputeResponse             |
      |------------------------------------->|
      |                                      |
      |  4. 触发 Escrow 释放                  |
      |<-------------------------------------|
```

---

## OpenFang 集成设计

> **重要说明:** OpenFang 是一个真实存在的开源项目，本节详细描述其与 ShareTokens 的集成架构。

### OpenFang 简介

**OpenFang** 是由 RightNow-AI 开发的开源 AI Agent 操作系统，用 Rust 编写，提供企业级安全防护和自主 Agent 能力。

| 属性 | 值 |
|------|-----|
| **项目类型** | 开源 AI Agent 操作系统 |
| **语言** | Rust (137K+ LOC, 14 crates) |
| **二进制大小** | ~32MB 单文件 |
| **许可证** | MIT |
| **版本** | v0.1.0 (2026年2月发布) |
| **官网** | https://www.openfang.sh/ |
| **GitHub** | https://github.com/RightNow-AI/openfang |

#### 核心能力

| 能力 | 说明 |
|------|------|
| **7个 Hands** | 自主运行的 Agent 能力包 (Clip, Lead, Collector, Predictor, Researcher, Twitter, Browser) |
| **16层安全防护** | WASM沙箱、Merkle审计链、污点追踪等纵深防御体系 |
| **40个消息渠道** | Telegram, Discord, Slack, WhatsApp 等全平台支持 |
| **27个 LLM Provider** | OpenAI, Anthropic, Gemini, DeepSeek 等全覆盖 |
| **53+ 内置工具** | + MCP + A2A 协议支持 |

#### 关键 Hands 与 ShareTokens 映射

| OpenFang Hand | 功能 | ShareTokens 应用 |
|---------------|------|------------------|
| **Collector** | 持续监控、知识图谱构建 | 想法收集、市场调研、竞品监控 |
| **Lead** | 每日自动发现潜在客户 | 资源匹配、销售线索挖掘 |
| **Researcher** | 深度研究、交叉验证、引用报告 | 想法评估、多AI评估系统 |
| **Browser** | 网页自动化、强制审批门 | 安全的链上交互操作 |
| **Clip** | 视频自动剪辑、字幕、发布 | 内容创作工作流 |
| **Predictor** | 多源信号预测、Brier评分 | 市场预测、风险评估 |
| **Twitter** | 社交账号自主运营 | 营销推广、社区管理 |

---

### 集成架构图

```
+-----------------------------------------------------------------------------+
|                    ShareTokens <-> OpenFang 集成架构                          |
+-----------------------------------------------------------------------------+

+-----------------------------------------------------------------------------+
|  用户层 (GenieBot UI)                                                        |
|  +-----------------------------------------------------------------------+  |
|  |  React Web Application                                                |  |
|  |  +-- 想法孵化界面          +-- 服务市场浏览                            |  |
|  |  +-- Workflow 执行监控     +-- 争议仲裁投票                            |  |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
                                    |
                                    | REST/gRPC
                                    v
+-----------------------------------------------------------------------------+
|  ShareTokens Chain (Cosmos SDK)                                             |
|  +-----------------------------------------------------------------------+  |
|  |  x/compute (Level 1 LLM API)                                          |  |
|  |  +-- API Key 托管         +-- 请求路由         +-- 计费结算            |  |
|  +-----------------------------------------------------------------------+  |
|  |  x/escrow | x/identity | x/trust | x/bank                             |  |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
                                    |
                                    | ABCI Events / gRPC
                                    v
+-----------------------------------------------------------------------------+
|  OpenFang Provider Plugin (Sidecar 部署)                                     |
|  +-----------------------------------------------------------------------+  |
|  |  ShareTokens Bridge                                                    |  |
|  |  +-- ChainEventSubscriber   +-- ComputeRequestHandler                 |  |
|  |  +-- ResponseSubmitter      +-- EscrowManager                          |  |
|  +-----------------------------------------------------------------------+  |
|                                    |                                         |
|  +---------------------------------+------------------------------------+    |
|  |  OpenFang Kernel (Rust)                                              |    |
|  |  +----------------------------------------------------------------+  |    |
|  |  | 16层安全防护                                                       |  |    |
|  |  | +-- 1. WASM 双计量沙箱        +-- 2. Merkle 哈希链审计            |  |    |
|  |  | +-- 3. 信息流污点追踪        +-- 4. Ed25519 签名清单              |  |    |
|  |  | +-- 5. SSRF 防护             +-- 6. Secret Zeroization           |  |    |
|  |  | +-- 7. OFP 双向认证          +-- 8. 能力门控                     |  |    |
|  |  | +-- 9-16. 其他安全层...                                          |  |    |
|  |  +----------------------------------------------------------------+  |    |
|  |                                    |                                  |    |
|  |  +---------------------------------+--------------------------------+|    |
|  |  | OpenFang Runtime                                                   ||    |
|  |  | +-- 3 Native LLM Drivers    +-- 53 Built-in Tools                 ||    |
|  |  | +-- MCP Protocol            +-- A2A Protocol                      ||    |
|  |  | +-- SQLite + Vector Memory  +-- Session Management                ||    |
|  |  +----------------------------------------------------------------+  |    |
|  |                                    |                                  |    |
|  |  +---------------------------------+--------------------------------+|    |
|  |  | OpenFang Hands (7 Bundled)                                         ||    |
|  |  | +-- Collector (监控/采集)   +-- Lead (线索挖掘)                    ||    |
|  |  | +-- Researcher (深度研究)   +-- Browser (网页自动化)               ||    |
|  |  | +-- Clip (视频剪辑)         +-- Predictor (预测)                   ||    |
|  |  | +-- Twitter (社交运营)                                             ||    |
|  |  +----------------------------------------------------------------+  |    |
|  +-----------------------------------------------------------------------+  |
|                                    |                                         |
|  +---------------------------------+------------------------------------+    |
|  |  LLM Provider Gateway                                                  |    |
|  |  +-- OpenAI Provider        +-- Anthropic Provider                    |    |
|  |  +-- Gemini Provider         +-- OpenRouter Provider                  |    |
|  |  +-- DeepSeek Provider       +-- Other 27+ Providers                  |    |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
```

---

### 具体的集成接口定义

#### 1. Agent 调用接口

ShareTokens 通过 OpenFang 的 OpenAI-Compatible API 调用 Agent 能力：

```typescript
// OpenFang OpenAI-Compatible API Endpoint
// POST http://localhost:4200/v1/chat/completions

interface OpenFangAgentRequest {
  model: string;           // Agent ID, e.g., "researcher", "collector"
  messages: Message[];
  stream?: boolean;
  temperature?: number;
  max_tokens?: number;
  metadata?: {
    sharetokens_request_id: string;  // 链上请求 ID
    escrow_id: string;               // 托管账户 ID
    requester_address: string;       // 请求者链上地址
  };
}

interface OpenFangAgentResponse {
  id: string;
  object: string;
  choices: [{
    index: number;
    message: {
      role: string;
      content: string;
    };
    finish_reason: string;
  }];
  usage: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
  // OpenFang 扩展字段
  sharetokens_proof: {
    merkle_root: string;     // Merkle 审计链根哈希
    execution_hash: string;  // 执行结果哈希
    signature: string;       // Ed25519 签名
  };
}
```

#### 2. Hand 工具接口

OpenFang Hands 作为 ShareTokens 的 Level 2/3 服务：

```typescript
// ShareTokens Bridge - Hand 调用接口
interface ShareTokensHandInterface {
  // Collector Hand - 数据收集
  collector: {
    startMonitoring(params: {
      target: string;           // 监控目标 (公司/人/话题)
      filters: string[];        // 过滤条件
      alert_threshold: number;  // 告警阈值
      callback_webhook: string; // 结果回调
    }): Promise<MonitoringJob>;

    getKnowledgeGraph(jobId: string): Promise<KnowledgeGraph>;

    getAlerts(jobId: string): Promise<Alert[]>;
  };

  // Lead Hand - 线索挖掘
  lead: {
    startDiscovery(params: {
      icp: ICPProfile;          // 理想客户画像
      daily_limit: number;
      min_score: number;
    }): Promise<DiscoveryJob>;

    getLeads(jobId: string): Promise<Lead[]>;
  };

  // Researcher Hand - 深度研究
  researcher: {
    startResearch(params: {
      topic: string;
      depth: "quick" | "standard" | "comprehensive";
      sources: string[];
      citation_format: "APA" | "MLA" | "Chicago";
    }): Promise<ResearchJob>;

    getReport(jobId: string): Promise<ResearchReport>;
  };

  // Browser Hand - 网页自动化 (带强制审批门)
  browser: {
    executeTask(params: {
      actions: BrowserAction[];
      require_approval: boolean;  // 涉及交易必须 true
    }): Promise<BrowserTask>;

    approveAction(taskId: string, actionId: string): Promise<void>;
  };
}

// Hand 任务状态同步到链上
interface HandTaskOnChain {
  task_id: string;
  hand_type: string;
  status: "pending" | "running" | "completed" | "failed";
  result_hash?: string;       // IPFS CID
  proof: {
    merkle_root: string;
    signature: string;
  };
  escrow_release_requested: boolean;
}
```

#### 3. 事件回调接口

OpenFang 通过回调向 ShareTokens 链报告执行结果：

```typescript
// ShareTokens Bridge 事件回调
interface ShareTokensEventCallbacks {
  // 计算请求完成回调
  onComputeComplete: (event: {
    request_id: string;
    provider_address: string;
    result: {
      content: string;
      tokens_used: TokenUsage;
      actual_cost: bigint;
    };
    proof: {
      merkle_root: string;
      execution_hash: string;
      signature: string;
    };
  }) => Promise<void>;

  // Hand 任务状态更新回调
  onHandTaskUpdate: (event: {
    task_id: string;
    hand_type: string;
    status: string;
    progress: number;
    intermediate_result?: string;
  }) => Promise<void>;

  // 托管释放请求回调
  onEscrowReleaseRequest: (event: {
    escrow_id: string;
    request_id: string;
    amount: bigint;
    recipient: string;
    proof: MerkleProof;
  }) => Promise<void>;

  // 错误/超时回调
  onError: (event: {
    request_id: string;
    error_code: string;
    error_message: string;
    requires_dispute: boolean;
  }) => Promise<void>;
}

// gRPC 服务定义 (proto)
syntax = "proto3";

package sharetokens.openfang.v1;

service OpenFangBridge {
  // 计算请求处理
  rpc HandleComputeRequest(ComputeRequest) returns (ComputeResponse);

  // 任务状态查询
  rpc GetTaskStatus(TaskStatusRequest) returns (TaskStatusResponse);

  // 任务取消
  rpc CancelTask(CancelTaskRequest) returns (CancelTaskResponse);

  // 健康检查
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}

message ComputeRequest {
  string request_id = 1;
  string model = 2;
  string prompt = 3;
  bytes prompt_hash = 4;
  uint64 nonce = 5;
  bytes signature = 6;
  string price_offer = 7;  // sdk.Coin
  uint32 timeout_seconds = 8;
  string requester_address = 9;
  string escrow_id = 10;
}

message ComputeResponse {
  string response_id = 1;
  string request_id = 2;
  string result = 3;
  TokenUsage tokens_used = 4;
  string actual_cost = 5;
  bytes result_hash = 6;
  bytes signature = 7;
  MerkleProof proof = 8;
}

message MerkleProof {
  string merkle_root = 1;
  repeated string siblings = 2;
  uint32 leaf_index = 3;
}
```

---

### 安全模型

#### 1. API Key 安全传递

```
+-----------------------------------------------------------------------------+
|                    API Key 安全传递流程                                       |
+-----------------------------------------------------------------------------+

  ShareTokens Chain                      OpenFang Provider
       |                                       |
       |  1. Provider 注册 API Key             |
       |     (链上存储加密密文)                  |
       |-------------------------------------->|
       |                                       |
       |  2. 用户发起 ComputeRequest           |
       |     (指定 Provider)                   |
       |-------------------------------------->|
       |                                       |
       |                          3. Provider 监听到请求
       |                                       |
       |                          4. 在 WASM 沙箱内解密 API Key
       |                             +-------------------+
       |                             |  WASM Sandbox     |
       |                             |  +-- 解密 Key     |
       |                             |  +-- 调用 LLM API  |
       |                             |  +-- Zeroize Key  |
       |                             +-------------------+
       |                                       |
       |  5. 返回结果 + Merkle 证明             |
       |<--------------------------------------|
       |                                       |
       |  6. 验证证明，释放托管                  |
       |-------------------------------------->|
       |                                       |

关键安全特性:
- API Key 永不离开 WASM 沙箱
- Secret Zeroization: 使用后立即擦除内存
- 所有操作记录到 Merkle 审计链
- 污点追踪确保敏感数据不泄露
```

#### 2. WASM 沙箱边界

```typescript
// OpenFang WASM 沙箱安全配置
interface WASMSandboxConfig {
  // 双计量机制
  fuel_metering: {
    initial_fuel: 1000000;      // 初始燃料
    fuel_per_instruction: 1;    // 每条指令消耗
    out_of_fuel_action: "terminate";  // 燃料耗尽行为
  };

  epoch_interruption: {
    epoch_duration_ms: 10;      // 纪元时长
    max_epochs_per_call: 1000;  // 单次调用最大纪元
  };

  // 看门狗线程
  watchdog: {
    check_interval_ms: 100;
    max_execution_time_ms: 30000;  // 最大执行时间
    terminate_on_timeout: true;
  };

  // 资源限制
  resource_limits: {
    max_memory_pages: 256;      // 最大内存页 (16MB)
    max_table_size: 65536;
    max_br_table_size: 16384;
  };

  // 能力门控
  capability_gates: {
    network_access: "whitelist";  // 只允许白名单网络
    file_system_access: "none";   // 无文件系统访问
    env_access: "restricted";     // 限制环境变量访问
  };
}

// ShareTokens 特定限制
interface ShareTokensSandboxPolicy extends WASMSandboxConfig {
  // 只允许访问的 LLM Provider 端点
  allowed_endpoints: [
    "api.openai.com",
    "api.anthropic.com",
    "generativelanguage.googleapis.com",
    "openrouter.ai",
    "api.deepseek.com"
  ];

  // 禁止的操作
  forbidden_operations: [
    "external_file_write",
    "arbitrary_code_exec",
    "network_listen",
    "process_fork"
  ];
}
```

#### 3. 权限控制 (基于 MQ)

```typescript
// ShareTokens MQ 权限控制集成
interface MQBasedPermissions {
  // MQ 等级与权限映射
  mq_levels: {
    newcomer: { range: [0, 50],      max_daily_spend: 100,  max_requests_per_hour: 10 };
    member:   { range: [50, 100],    max_daily_spend: 500,  max_requests_per_hour: 50 };
    trusted:  { range: [100, 200],   max_daily_spend: 2000, max_requests_per_hour: 200 };
    expert:   { range: [200, 500],   max_daily_spend: 10000, max_requests_per_hour: 500 };
    guardian: { range: [500, null],  max_daily_spend: null, max_requests_per_hour: null };
  };

  // Provider MQ 要求
  provider_requirements: {
    min_mq_for_provider: 50;      // 成为 Provider 最低 MQ
    min_mq_for_premium_models: 100; // 提供 GPT-4 等高端模型
    dispute_penalty: -10;          // 争议败诉 MQ 惩罚
  };

  // 服务匹配权重
  matching_weights: {
    mq_weight: 0.4;              // MQ 权重
    price_weight: 0.3;           // 价格权重
    rating_weight: 0.2;          // 评分权重
    latency_weight: 0.1;         // 延迟权重
  };
}
```

---

### 部署模式

#### 推荐方案: Sidecar 模式

```
+-----------------------------------------------------------------------------+
|                    Sidecar 部署架构                                           |
+-----------------------------------------------------------------------------+

+-----------------------------------------------------------------------------+
|  服务提供者节点 (单个物理机或 K8s Pod)                                         |
|                                                                             |
|  +---------------------------+  +----------------------------------------+  |
|  | ShareTokens Chain Node    |  | OpenFang Provider (Sidecar)            |  |
|  | (Cosmos SDK)              |  |                                        |  |
|  | +-- CometBFT P2P          |  | +-- ShareTokens Bridge                 |  |
|  | +-- x/compute Module      |  | +-- OpenFang Kernel (32MB binary)      |  |
|  | +-- x/escrow Module       |  | +-- WASM Runtime                       |  |
|  | +-- x/trust Module        |  | +-- Hands Runtime                      |  |
|  | +-- State DB (RocksDB)    |  | +-- LLM Provider Gateway               |  |
|  +---------------------------+  +----------------------------------------+  |
|         |                                    |                              |
|         |        gRPC / Unix Socket          |                              |
|         +------------------------------------+                              |
|                                                                             |
+-----------------------------------------------------------------------------+

优点:
- 紧密耦合，低延迟通信
- 共享生命周期管理
- 简化运维 (单一部署单元)
- 适合中小规模服务提供者
```

#### 替代方案: 独立服务模式

```
+-----------------------------------------------------------------------------+
|                    独立服务部署架构                                           |
+-----------------------------------------------------------------------------+

+-----------------------+     gRPC/REST      +--------------------------------+
| ShareTokens Chain     |<------------------>| OpenFang Provider Cluster      |
| (多个验证者节点)       |                    | (独立扩缩容)                    |
| +-- Validator 1       |                    | +-- OpenFang Instance 1        |
| +-- Validator 2       |                    | +-- OpenFang Instance 2        |
| +-- Full Node 1..N    |                    | +-- OpenFang Instance N        |
+-----------------------+                    | +-- Load Balancer              |
                                             +--------------------------------+

优点:
- 可独立扩缩容 OpenFang 集群
- 适合大规模服务提供者
- OpenFang 可服务多个链
- 高可用性更好

缺点:
- 网络延迟增加
- 运维复杂度提高
```

#### OpenFang 配置文件

```toml
# /etc/openfang/sharetokens-bridge.toml

[openfang]
# OpenFang 核心配置
data_dir = "/var/lib/openfang"
log_level = "info"
dashboard_enabled = true
dashboard_port = 4200

[sharetokens]
# ShareTokens 链连接配置
chain_id = "sharetokens-1"
rpc_endpoint = "http://localhost:26657"
grpc_endpoint = "http://localhost:9090"

# Provider 身份
provider_address = "share1abc..."  # 链上注册的 Provider 地址
mnemonic = "encrypted:..."         # 加密的助记词

[compute.market]
# 服务市场配置
enabled = true
min_mq = 50                        # 最低 MQ 要求
escrow_required = true

# 支持的模型
[[compute.models]]
name = "gpt-4o"
provider = "openai"
price_per_1k_tokens = "0.05stt"
enabled = true

[[compute.models]]
name = "claude-3-opus"
provider = "anthropic"
price_per_1k_tokens = "0.03stt"
enabled = true

[[compute.models]]
name = "deepseek-chat"
provider = "deepseek"
price_per_1k_tokens = "0.001stt"
enabled = true

[security]
# 安全配置
sandbox_enabled = true
audit_chain_enabled = true
taint_tracking_enabled = true

# API Key 存储 (AES-256-GCM 加密)
[secrets.vault]
type = "encrypted_file"
path = "/var/lib/openfang/vault/secrets.enc"

[secrets.api_keys]
# API Key 从链上获取，不在本地存储
source = "chain"
cache_ttl_seconds = 3600

[bridge]
# Bridge 服务配置
grpc_port = 50051
max_concurrent_requests = 1000
request_timeout_seconds = 60

[hands]
# 启用的 Hands
enabled = ["collector", "lead", "researcher", "browser"]

[hands.browser]
# Browser Hand 安全配置
require_approval_for_purchase = true
allowed_domains = ["*"]
```

---

### 集成实施路线图

| 阶段 | 内容 | 依赖 |
|------|------|------|
| **Phase 1** | OpenFang Sidecar 基础集成 | x/compute 完成 |
| **Phase 2** | ShareTokens Bridge gRPC 服务 | Phase 1 |
| **Phase 3** | 7 Hands 与 ShareTokens 映射 | Phase 2 |
| **Phase 4** | API Key 链上托管 + WASM 沙箱集成 | Phase 2 |
| **Phase 5** | MQ 权限控制集成 | x/trust 完成 |
| **Phase 6** | 生产环境部署 + 监控告警 | Phase 3-5 |

---

## 与 Level 2/3 服务的关系

```
+-----------------------------------------------------------------------------+
|                          服务调用链                                          |
+-----------------------------------------------------------------------------+

  Level 3 Workflow                Level 2 Agent               Level 1 LLM API
       |                              |                             |
       |  1. 用户发起工作流            |                             |
       |----------------------------->|                             |
       |                              |                             |
       |                              |  2. Agent 需要调用 LLM       |
       |                              |---------------------------->|
       |                              |                             |
       |                              |  3. x/compute 处理请求       |
       |                              |<----------------------------|
       |                              |                             |
       |  4. 返回结果                  |                             |
       |<-----------------------------|                             |
       |                              |                             |

说明：
- Level 2 Agent 服务依赖 Level 1 完成实际的 LLM 调用
- Level 3 Workflow 编排多个 Level 2 Agent
- x/compute 作为 Level 1 提供基础的 LLM API 能力
```

---

## Level 1 服务计费模型

```go
// 基于 Token 的计费
type BillingModel struct {
    // 基础价格
    BasePricePerToken sdk.Coin

    // 模型价格乘数
    ModelMultipliers map[string]sdk.Dec  // "gpt-4": 1.5, "gpt-3.5": 0.5

    // 动态定价
    DynamicPricing DynamicPricingConfig
}

// 费用计算
func CalculateCost(model string, usage TokenUsage, pricing PricingConfig) sdk.Coin {
    baseCost := pricing.BasePricePerToken.Amount.Mul(sdk.NewInt(int64(usage.Total)))

    if multiplier, ok := pricing.ModelMultipliers[model]; ok {
        baseCost = baseCost.Mul(multiplier)
    }

    return sdk.NewCoin("stt", baseCost)
}
```

---

## OpenFang 参考链接

| 资源 | 链接 |
|------|------|
| 官方网站 | https://www.openfang.sh/ |
| GitHub 仓库 | https://github.com/RightNow-AI/openfang |
| 文档 | https://docs.openfang.sh/ |
| Discord | https://discord.gg/openfang |
| Twitter/X | https://twitter.com/openfang_ai |

---

## 防欺诈检测 (FraudIndicators)

> **说明:** 本节内容从 12-misc.md 迁移整合，用于检测服务市场中的异常行为。

### 概述

FraudIndicators 是一个综合的欺诈检测系统，用于识别服务市场中的刷单、价格操纵、女巫攻击等恶意行为。

### 数据结构

```go
// types/fraud.go

type FraudIndicators struct {
    // 风险等级
    RiskLevel    RiskLevel    // low | medium | high | critical
    RiskScore    uint8        // 0-100

    // 交易模式异常
    TradingPatterns TradingPatterns

    // 关联账户检测
    AccountLinks AccountLinks

    // 行为异常
    BehavioralAnomalies BehavioralAnomalies

    // 历史标记
    Flags []ReputationFlag

    // 最后检测时间
    LastChecked time.Time
}

type RiskLevel string

const (
    RiskLevelLow      RiskLevel = "low"
    RiskLevelMedium   RiskLevel = "medium"
    RiskLevelHigh     RiskLevel = "high"
    RiskLevelCritical RiskLevel = "critical"
)
```

### 交易模式检测

```go
type TradingPatterns struct {
    // 对手方集中度
    HighCounterpartyConcentration bool
    TopCounterpartyRatio           sdk.Dec  // 最大对手方交易占比

    // 交易时间异常
    UnusualTradingHours bool
    NightTradingRatio   sdk.Dec  // 夜间交易占比

    // 交易金额异常
    RoundNumberPreference bool
    AvgTransactionSize    sdk.Coin

    // 快速连续交易
    RapidTransactions bool
    AvgTimeBetweenTx  time.Duration
}
```

### 关联账户检测

```go
type AccountLinks struct {
    // 共享 IP
    SharedIPs      []string
    IPOverlapRatio sdk.Dec

    // 共享设备
    SharedDevices      []string
    DeviceOverlapRatio sdk.Dec

    // 交易网络重叠
    CounterpartyOverlap []sdk.AccAddress
    OverlapRatio        sdk.Dec

    // 资金流向异常
    CircularFlows bool    // 资金循环
    BackAndForth  bool    // 频繁来回转账
}
```

### 行为异常检测

```go
type BehavioralAnomalies struct {
    // 新账户高风险行为
    NewAccountHighVolume bool
    AccountAge           time.Duration
    VolumeInFirstWeek    sdk.Coin

    // 评价异常
    RatingManipulation bool
    AvgRatingDeviation sdk.Dec  // 评分偏离正常值

    // 争议行为
    ExcessiveDisputes bool
    DisputeWinRate    sdk.Dec  // 争议胜率异常

    // 取消模式
    CancellationPattern bool
    CancelRate          sdk.Dec
    CancelAfterMatchRatio sdk.Dec  // 匹配后取消率
}
```

### 风险标记

```go
type ReputationFlag struct {
    Id          uint64
    Type        FlagType
    Severity    Severity    // info | warning | danger
    Description string
    Source      FlagSource  // algorithm | manual | report
    CreatedAt   time.Time
    ResolvedAt  *time.Time
    ResolvedBy  *sdk.AccAddress
}

type FlagType string

const (
    FlagTypeSuspectedWashTrading  FlagType = "suspected_wash_trading"   // 疑似刷单
    FlagTypePriceManipulation     FlagType = "price_manipulation"       // 价格操纵
    FlagTypeCircularTransactions  FlagType = "circular_transactions"    // 循环交易
    FlagTypeExcessiveDisputes     FlagType = "excessive_disputes"       // 过多争议
    FlagTypeRatingManipulation    FlagType = "rating_manipulation"      // 评分操纵
    FlagTypeSybilSuspected        FlagType = "sybil_suspected"          // 疑似女巫攻击
    FlagTypeUnusualActivity       FlagType = "unusual_activity"         // 异常活动
    FlagTypeConfirmedFraud        FlagType = "confirmed_fraud"          // 确认欺诈
)

type Severity string

const (
    SeverityInfo    Severity = "info"
    SeverityWarning Severity = "warning"
    SeverityDanger  Severity = "danger"
)

type FlagSource string

const (
    FlagSourceAlgorithm FlagSource = "algorithm"
    FlagSourceManual    FlagSource = "manual"
    FlagSourceReport    FlagSource = "report"
)
```

### 风险评分流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          欺诈检测流程                                         │
└─────────────────────────────────────────────────────────────────────────────┘

                    ┌───────────────────────────────────────┐
                    │           交易完成后                   │
                    └───────────────────┬───────────────────┘
                                        │
                                        ▼
                    ┌───────────────────────────────────────┐
                    │          FraudIndicators              │
                    │          欺诈指标检测                  │
                    └───────────────────┬───────────────────┘
                                        │
                ┌───────────────────────┼───────────────────────┐
                │                       │                       │
                ▼                       ▼                       ▼
        ┌───────────────┐       ┌───────────────┐       ┌───────────────┐
        │ 交易模式分析   │       │ 关联账户检测   │       │ 行为异常分析   │
        │               │       │               │       │               │
        │ - 对手方集中度 │       │ - 共享 IP     │       │ - 新账户行为   │
        │ - 时间模式    │       │ - 共享设备    │       │ - 评价异常     │
        │ - 金额模式    │       │ - 交易重叠    │       │ - 争议行为     │
        └───────┬───────┘       └───────┬───────┘       └───────┬───────┘
                │                       │                       │
                └───────────────────────┼───────────────────────┘
                                        │
                                        ▼
                    ┌───────────────────────────────────────┐
                    │            风险评分                    │
                    │         riskScore: 0-100             │
                    └───────────────────┬───────────────────┘
                                        │
                ┌───────────────────────┼───────────────────────┐
                │                       │                       │
                ▼                       ▼                       ▼
        ┌───────────────┐       ┌───────────────┐       ┌───────────────┐
        │  低风险        │       │  中风险        │       │  高风险        │
        │  (0-30)       │       │  (31-70)      │       │  (71-100)     │
        │               │       │               │       │               │
        │  正常交易      │       │  加强监控      │       │  限制交易      │
        │  无限制        │       │  降低限额      │       │  人工审核      │
        └───────────────┘       └───────────────┘       └───────────────┘
```

### Keeper 接口

```go
// keeper/fraud.go

// 检测欺诈指标
func (k Keeper) DetectFraud(ctx sdk.Context, address sdk.AccAddress) (FraudIndicators, error)

// 获取风险评分
func (k Keeper) GetRiskScore(ctx sdk.Context, address sdk.AccAddress) uint8

// 添加风险标记
func (k Keeper) AddFlag(ctx sdk.Context, address sdk.AccAddress, flag ReputationFlag) error

// 解决风险标记
func (k Keeper) ResolveFlag(ctx sdk.Context, flagId uint64, resolver sdk.AccAddress) error

// 根据风险等级获取限制
func (k Keeper) GetRestrictionsByRiskLevel(level RiskLevel) Restrictions
```

---

[<- 上一章：共识层](./04-consensus.md) | [返回索引](./00-index.md) | [下一章：Oracle服务 ->](./06-exchange.md)
