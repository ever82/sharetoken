# Phase 1: ShareTokens 核心模块 B - 技术研究

**Researched:** 2026-03-03
**Domain:** 服务市场模块 (Service Marketplace) + 托管支付模块 (Escrow)
**Confidence:** HIGH (Akash/LLM routing), MEDIUM (Cosmos escrow implementations), HIGH (Cosmos SDK patterns)

## Summary

本研究报告聚焦于 ShareTokens 的两个核心模块：**服务市场模块**（三层服务结构、智能路由）和**托管支付模块**（任务锁定、完成释放、争议冻结）。

通过调研 Akash Network、OpenRouter、Ocean Protocol 等类似项目，以及 Cosmos SDK 生态系统中的现有实现，我们得出了以下关键结论：

**Primary recommendation:**
1. **服务市场模块**：参考 Akash 的密封拍卖机制 + OpenRouter 的 LLM 路由策略，基于 Cosmos SDK 自定义 x/service 模块实现
2. **托管支付模块**：基于 Cosmos SDK x/auth/vesting 模式的锁定账户机制 + 自定义 x/escrow 模块实现争议冻结

---

## 类似项目分析

### 1. Akash Network - 去中心化算力市场

| 属性 | 详情 |
|------|------|
| **项目链接** | [akash.network](https://akash.network) / [GitHub Docs](https://github.com/akash-network/docs) |
| **技术栈** | Cosmos SDK + Tendermint (CometBFT) |
| **代币** | AKT |
| **相似之处** | 去中心化资源市场、密封拍卖订单匹配、链上结算 |

**核心机制 - 密封拍卖订单匹配：**

```
用户请求 (CPU/GPU/存储需求 + 最高价格)
          ↓
    多个提供者匿名投标
          ↓
    系统自动匹配最优组合
          ↓
    生成链上租约 (Lease)
```

**可学习的架构决策：**

| Akash 设计 | ShareTokens 可借鉴 |
|------------|-------------------|
| **动态密封拍卖** | 服务请求可采用类似拍卖机制获取最优价格 |
| **反向拍卖** | 消费者竞价租用资源，驱动价格竞争下降 |
| **链上结算** | 资源使用数据触发自动链上结算，无需第三方托管 |
| **无信任环境** | 区块链技术确保安全，降低欺诈风险 |
| **四层架构** | 区块链层、应用层、提供者层、用户层分离 |

**成本优势参考：**
- Akash 比 AWS/GCP/Azure 便宜约 66%
- ShareTokens 可追求类似成本优势

---

### 2. OpenRouter - LLM API 聚合路由

| 属性 | 详情 |
|------|------|
| **项目链接** | [openrouter.ai](https://openrouter.ai) |
| **服务类型** | LLM API 聚合平台 |
| **模型数量** | 80-100+ 种 AI 模型 |
| **相似之处** | 多模型路由、智能选择、统一接口 |

**核心机制 - 智能路由：**

```
用户请求 (prompt + 任务类型)
          ↓
    路由决策引擎
    ├─ 任务复杂度分析
    ├─ 成本优化
    ├─ 性能需求
    └─ 可用性检查
          ↓
    选择最优 LLM 提供者
          ↓
    返回结果 + 统一格式
```

**LLMRouter 项目参考 (2025年12月, 1K+ Stars)：**

| 特性 | 描述 |
|------|------|
| **16+ 路由策略** | KNN、SVM、MLP、图方法、BERT-based、混合概率 |
| **四大路由类型** | 单轮、多轮、Agentic、个性化路由 |
| **智能决策** | 基于任务复杂度、成本、性能自动分配 |
| **插件式扩展** | 自定义 router/task/metric 可无缝接入 |

**可学习的架构决策：**

| OpenRouter 设计 | ShareTokens 可借鉴 |
|-----------------|-------------------|
| **统一接口** | 将不同 LLM 提供者统一为兼容格式 |
| **任务类型路由** | 根据编程/长文本/通用对话自动选择模型 |
| **按需计费** | 灵活的计费模式支持 |
| **提供商抽象** | 用户无需关心底层 API 细节 |

---

### 3. Ocean Protocol - 去中心化数据市场

| 属性 | 详情 |
|------|------|
| **项目链接** | [oceanprotocol.com](https://oceanprotocol.com) |
| **技术栈** | Ethereum + Polygon |
| **核心功能** | 数据交换、Compute-to-Data |
| **相似之处** | AI 服务市场、数据定价、隐私保护 |

**核心机制 - Compute-to-Data：**

```
数据请求 → 算法发送到数据端 → 本地计算 → 返回结果
              ↓
        数据永不离开原始位置
```

**可学习的架构决策：**

| Ocean 设计 | ShareTokens 可借鉴 |
|------------|-------------------|
| **Data NFTs (ERC-721)** | 服务可表示为 NFT，保护知识产权 |
| **Data Tokens** | Token 化访问控制 |
| **Compute-to-Data** | AI Agent 在数据端执行，保护隐私 |
| **单边质押** | 解决 rug pull 问题 |

---

### 4. Golem Network + Render + iExec - 去中心化算力三巨头

| 项目 | GLM | RNDR | RLC |
|------|-----|------|-----|
| **网络** | Golem | Render | iExec |
| **专注领域** | 通用计算 + AI/ML | GPU 渲染 | 企业级 HPC |
| **链** | Ethereum | Ethereum | Ethereum |
| **Market Share** | ~15% (组合) | 渲染领域领先 | 企业专注 |

**2026 市场趋势：**
- 去中心化渲染市场份额：1-5% (2024)
- 年增长率：30%+
- AI 算力市场预测：500 亿美元 (2027)

**可学习的架构决策：**

| 设计模式 | ShareTokens 应用 |
|----------|-----------------|
| **闲置资源利用** | LLM API Key 闲置时共享 |
| **GPU 实例支持** | A100/H100 实例路由 |
| **被动收入机制** | 提供者赚取 STT |

---

## 可复用代码库

### Cosmos SDK 官方模块

#### 1. x/auth/vesting - 锁定账户机制

| 属性 | 详情 |
|------|------|
| **GitHub** | [cosmos/cosmos-sdk/x/auth/vesting](https://github.com/cosmos/cosmos-sdk/tree/main/x/auth/vesting) |
| **许可证** | Apache-2.0 |
| **成熟度** | HIGH (Cosmos Hub 官方使用) |
| **状态** | v0.52 后 deprecated，但向后兼容 |

**账户类型：**

| 类型 | 用途 | ShareTokens 应用 |
|------|------|-----------------|
| `BaseVestingAccount` | 基础锁定 | Escrow 基础结构 |
| `ContinuousVestingAccount` | 线性释放 | 部分释放场景 |
| `DelayedVestingAccount` | 延迟释放 | 任务完成后释放 |
| `PeriodicVestingAccount` | 周期释放 | 里程碑付款 |

**核心代码模式：**

```go
// 锁定硬币计算
func (va VestingAccount) LockedCoins(t Time) Coins {
   return max(va.GetVestingCoins(t) - va.DelegatedVesting, 0)
}

// 可消费余额
func (k Keeper) SpendableCoins(ctx Context, addr AccAddress) Coins {
    bc := k.GetBalances(ctx, addr)
    v := k.LockedCoins(ctx, addr)
    return bc - v
}
```

**直接可用性：** 可作为 Escrow 模块的基础模式，但需要自定义争议冻结逻辑

---

#### 2. Cosmos SDK Bank Module - 代币转移

| 属性 | 详情 |
|------|------|
| **GitHub** | [cosmos/cosmos-sdk/x/bank](https://github.com/cosmos/cosmos-sdk/tree/main/x/bank) |
| **许可证** | Apache-2.0 |
| **成熟度** | HIGH (生产级) |

**核心功能：**
- `SendCoins` - 代币转移
- `LockedCoins` - 锁定查询
- `SpendableCoins` - 可消费余额

**直接可用性：** 直接使用，无需开发

---

#### 3. Ignite CLI - Cosmos SDK 开发脚手架

| 属性 | 详情 |
|------|------|
| **官网** | [ignite.com](https://ignite.com) |
| **GitHub** | [ignite/cli](https://github.com/ignite/cli) |
| **许可证** | Apache-2.0 |
| **成熟度** | HIGH (官方推荐工具) |

**快速启动命令：**

```bash
# 初始化新链
ignite scaffold chain github.com/sharetokens/sharetokens-chain

# 创建自定义模块
ignite scaffold module service --dep bank,auth
ignite scaffold module escrow --dep bank,auth

# 添加消息类型
ignite scaffold message create-escrow amount beneficiary duration
ignite scaffold message release-escrow escrow-id
ignite scaffold message lock-escrow escrow-id
```

**直接可用性：** 强烈推荐用于快速启动开发

---

### 第三方开源实现

#### 4. MeshJS - Cardano 智能合约库

| 属性 | 详情 |
|------|------|
| **GitHub** | [MeshJS/mesh](https://github.com/MeshJS/mesh) |
| **许可证** | Apache-2.0 |
| **Stars** | 2K+ |
| **平台** | Cardano (非 Cosmos) |

**可用合约模板：**

| 合约 | 描述 | 借鉴价值 |
|------|------|---------|
| **Escrow** | 安全资产交换 | 参考接口设计 |
| **Marketplace** | NFT 买卖市场 | 参考市场机制 |
| **Vesting** | 锁定释放 | 参考释放逻辑 |
| **Payment Splitter** | 支付分配 | 争议裁决后分配 |

**注意：** 平台不同，但接口设计可借鉴

---

#### 5. TokenLockup - 代币锁定释放

| 属性 | 详情 |
|------|------|
| **GitHub** | [Cerebellum-Network/TokenLockup](https://github.com/Cerebellum-Network/TokenLockup) |
| **许可证** | MIT |
| **功能** | 计划代币释放、固定锁定 |

**直接可用性：** 可参考实现逻辑

---

### Akash Network 模块

#### 6. Akash Provider Module

| 属性 | 详情 |
|------|------|
| **GitHub** | [akash-network/provider](https://github.com/akash-network/provider) |
| **许可证** | Apache-2.0 |
| **成熟度** | HIGH (主网运行中) |

**核心组件：**
- `x/provider` - 提供者注册
- `x/deployment` - 部署管理
- `x/market` - 订单匹配
- `x/escrow` - 托管结算 (Akash 的实现)

**直接可用性：** 可参考 x/escrow 和 x/market 设计

---

## 快速启动建议

### 可以直接使用

| 组件 | 来源 | 操作 |
|------|------|------|
| **P2P 网络** | CometBFT 内置 | 无需开发 |
| **账户管理** | Cosmos SDK Auth | 无需开发 |
| **代币转移** | Cosmos SDK Bank | 无需开发 |
| **链脚手架** | Ignite CLI | `ignite scaffold chain` |
| **锁定账户模式** | x/auth/vesting | 参考设计模式 |
| **共识引擎** | CometBFT | 无需开发 |

### 需要适配

| 组件 | 来源 | 适配工作 |
|------|------|---------|
| **Escrow 模块** | 参考 Akash x/escrow + Cosmos vesting | 添加争议冻结、部分释放、超时释放 |
| **服务注册** | 参考 Akash x/provider | 添加三层服务结构 (LLM/Agent/Workflow) |
| **定价策略** | 参考 OpenRouter | 添加按次/批量/订阅模式 |
| **智能路由** | 参考 LLMRouter + Akash Market | 实现MQ/价格/能力/负载均衡路由 |

### 必须自研

| 组件 | 原因 |
|------|------|
| **x/service 模块** | 三层服务结构是 ShareTokens 独有需求 |
| **MQ 路由权重** | 基于信任评分的路由是独特需求 |
| **争议冻结逻辑** | 与 Trust System 的集成需要自定义 |
| **超时自动释放** | 业务规则需要自定义实现 |

---

## 架构模式建议

### 服务市场模块 (x/service)

```
chain/x/service/
├── keeper/
│   ├── keeper.go           # 主 Keeper
│   ├── provider.go         # 提供者注册 (参考 Akash)
│   ├── service.go          # 服务注册 (三层结构)
│   ├── pricing.go          # 定价策略
│   ├── routing.go          # 智能路由 (参考 OpenRouter)
│   ├── grpc_query.go       # 查询服务
│   └── msg_server.go       # 交易服务
├── types/
│   ├── service.pb.go       # 服务类型定义
│   ├── pricing.pb.go       # 定价类型
│   ├── routing.pb.go       # 路由类型
│   └── errors.go           # 错误定义
└── module.go
```

### 托管支付模块 (x/escrow)

```
chain/x/escrow/
├── keeper/
│   ├── keeper.go           # 主 Keeper
│   ├── escrow.go           # 托管 CRUD (参考 Cosmos vesting)
│   ├── release.go          # 释放操作
│   ├── dispute.go          # 争议锁定
│   ├── timeout.go          # 超时处理
│   ├── grpc_query.go       # 查询服务
│   └── msg_server.go       # 交易服务
├── types/
│   ├── escrow.pb.go        # 托管类型
│   ├── ruling.pb.go        # 裁决类型
│   └── errors.go           # 错误定义
└── module.go
```

---

## 代码示例

### Escrow 模块核心接口

```go
// Keeper defines the escrow module keeper
// 参考: Cosmos SDK x/auth/vesting + Akash x/escrow
type Keeper interface {
    // 托管操作
    CreateEscrow(ctx sdk.Context, params CreateEscrowParams) (Escrow, error)
    ReleaseEscrow(ctx sdk.Context, escrowId string, releaser sdk.AccAddress) error
    PartialRelease(ctx sdk.Context, escrowId string, amount sdk.Coins, releaser sdk.AccAddress) error
    CancelEscrow(ctx sdk.Context, escrowId string, creator sdk.AccAddress) error

    // 争议操作 (ShareTokens 特有)
    LockForDispute(ctx sdk.Context, escrowId string) error
    ResolveByRuling(ctx sdk.Context, escrowId string, ruling Ruling) error

    // 超时处理 (ShareTokens 特有)
    CheckTimeout(ctx sdk.Context, escrowId string) error

    // 查询
    GetEscrow(ctx sdk.Context, escrowId string) (Escrow, bool)
    GetEscrowsByCreator(ctx sdk.Context, creator sdk.AccAddress) []Escrow
    GetEscrowsByBeneficiary(ctx sdk.Context, beneficiary sdk.AccAddress) []Escrow

    // 依赖
    GetBankKeeper() bankkeeper.Keeper
}

// Escrow 状态
type EscrowStatus string

const (
    EscrowStatusActive    EscrowStatus = "active"     // 活跃
    EscrowStatusPartial   EscrowStatus = "partial"    // 部分释放
    EscrowStatusLocked    EscrowStatus = "locked"     // 争议锁定
    EscrowStatusReleased  EscrowStatus = "released"   // 已释放
    EscrowStatusCancelled EscrowStatus = "cancelled"  // 已取消
)

// Escrow 定义
type Escrow struct {
    Id          string
    Creator     sdk.AccAddress
    Beneficiary sdk.AccAddress
    Amount      sdk.Coins
    Released    sdk.Coins
    Status      EscrowStatus
    LockedAt    *time.Time      // 争议锁定时间
    CreatedAt   time.Time
    ExpiresAt   time.Time       // 超时时间
}
```

### 服务路由接口

```go
// 参考: OpenRouter + LLMRouter + Akash Market
type ServiceRouter interface {
    // 路由请求到最优提供者
    RouteRequest(ctx sdk.Context, request ServiceRequest) (RouteResult, error)

    // 获取可用提供者
    GetAvailableProviders(ctx sdk.Context, serviceId string) ([]Provider, error)

    // 计算路由分数
    CalculateScore(ctx sdk.Context, provider Provider, criteria RouteCriteria) sdk.Dec
}

// 路由标准
type RouteCriteria struct {
    MQWeight       sdk.Dec  // MQ 权重 (ShareTokens 特有)
    PriceWeight    sdk.Dec  // 价格权重
    CapabilityWeight sdk.Dec  // 能力权重
    LoadWeight     sdk.Dec  // 负载权重
}

// 路由策略
type RouteStrategy string

const (
    RouteStrategyMQPriority    RouteStrategy = "mq_priority"     // MQ 优先
    RouteStrategyPricePriority RouteStrategy = "price_priority"  // 价格优先
    RouteStrategyBalanced      RouteStrategy = "balanced"        // 均衡
    RouteStrategyAuction       RouteStrategy = "auction"         // 拍卖 (参考 Akash)
)
```

---

## 常见陷阱

### 1. 服务市场模块陷阱

| 陷阱 | 描述 | 解决方案 |
|------|------|---------|
| **路由中心化** | 单一路由器成为瓶颈 | 使用 DHT 分布式发现 + 本地路由决策 |
| **价格操纵** | 提供者恶意抬高/压低价格 | 密封拍卖机制 + 历史价格参考 |
| **服务描述欺诈** | 提供者虚假描述能力 | MQ 评分系统 + 服务验证机制 |
| **锁定资金滥用** | Escrow 资金被挪用 | 链上锁定 + 多签控制 |

### 2. 托管支付模块陷阱

| 陷阱 | 描述 | 解决方案 |
|------|------|---------|
| **争议永久冻结** | 争议未解决导致资金永远锁定 | 强制超时机制 + 默认裁决 |
| **部分释放复杂** | 多次部分释放状态管理混乱 | 原子操作 + 事件日志 |
| **时区问题** | 超时计算因时区出错 | 使用区块高度而非时间戳 |
| **重放攻击** | 释放交易被重放 | nonce + 幂等性设计 |

---

## 开放问题

### 1. 路由策略选择

**问题：** ShareTokens 应该采用哪种路由策略？

| 选项 | 优点 | 缺点 |
|------|------|------|
| **密封拍卖 (Akash)** | 价格发现、竞争性 | 延迟较高 |
| **固定定价** | 简单、快速 | 缺乏灵活性 |
| **动态定价** | 供需平衡 | 实现复杂 |

**推荐：** 混合策略 - LLM 服务用固定定价，Agent/Workflow 用密封拍卖

### 2. Escrow 超时机制

**问题：** 超时应该基于时间还是区块高度？

| 选项 | 优点 | 缺点 |
|------|------|------|
| **时间戳** | 用户友好 | 区块链时间不准确 |
| **区块高度** | 确定性 | 用户需要换算 |

**推荐：** 区块高度 + 预估时间显示

---

## Sources

### Primary (HIGH confidence)

- [Akash Network GitHub Docs](https://github.com/akash-network/docs) - 订单匹配机制
- [Cosmos SDK x/auth/vesting](https://github.com/cosmos/cosmos-sdk/tree/main/x/auth/vesting) - 锁定账户实现
- [Cosmos SDK Documentation](https://docs.cosmos.network/) - 模块开发指南
- [Ignite CLI](https://ignite.com) - 开发脚手架

### Secondary (MEDIUM confidence)

- [OpenRouter Platform](https://openrouter.ai) - LLM 路由策略参考
- [Ocean Protocol Documentation](https://oceanprotocol.com) - 数据市场设计参考
- [MeshJS GitHub](https://github.com/MeshJS/mesh) - 智能合约模板参考
- [Akash Network Price Analysis](https://coinmarketcap.com/currencies/akash-network/) - 市场数据

### Tertiary (LOW confidence - 需验证)

- Golem/Render/iExec 市场份额数据 - 需要进一步验证

---

## Metadata

**Confidence breakdown:**
- 类似项目分析: HIGH - Akash/OpenRouter/Ocean 文档详尽
- Cosmos SDK 模块: HIGH - 官方文档和源码可验证
- 第三方实现: MEDIUM - 需要评估许可证兼容性
- 市场数据: MEDIUM - 多来源交叉验证

**Research date:** 2026-03-03
**Valid until:** 2026-04-03 (1个月 - 技术栈稳定，但市场数据需更新)
