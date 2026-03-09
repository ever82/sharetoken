# x/idea - 想法模块

> **模块类型:** 用户插件（GenieBot插件的一部分）
> **技术栈:** Go (Cosmos SDK)
> **位置:** `src/chain/x/idea`
> **依赖:** 服务市场 (11-service)、身份账号 (10-identity)、Trust System (09-dispute)

---

## 概述

x/idea 是 ShareTokens 链上的想法管理模块，作为 GenieBot 插件的一部分，负责想法的版本控制、协作、众筹与贡献管理。作为 Cosmos SDK 模块实现，与 x/task、x/escrow 等模块紧密协作。

---

## 架构位置

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ShareTokens 架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  核心模块                                                                    │
│  ├── 服务市场 (11-service) ────────────────────────────────────────────┐   │
│  ├── 托管支付 (11-service)                                              │   │
│  └── Trust System (09-dispute)                                          │   │
│                                                                         │   │
│  用户插件                                                                ▼   │
│  └── GenieBot界面 (12-misc) ──┬── 想法系统 (07-idea) ← 本章              │
│                           │                                               │
│                           └── 任务市场 (08-task)                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 想法生命周期

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          想法生命周期                                         │
└─────────────────────────────────────────────────────────────────────────────┘

                            ┌──────────┐
                            │  draft   │  草稿阶段
                            └────┬─────┘
                                 │
                           Publish
                                 │
                                 ▼
                            ┌──────────┐
                            │published │  已发布 (可见但未众筹)
                            └────┬─────┘
                                 │
                        StartCampaign
                                 │
                                 ▼
                            ┌──────────┐
                            │ funding  │  众筹中
                            └────┬─────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
         Funded│            FundingFailed    CancelCampaign
              │                  │                  │
              ▼                  ▼                  ▼
       ┌───────────┐      ┌──────────┐       ┌──────────┐
       │in_progress│      │ archived │       │ archived │
       └─────┬─────┘      └──────────┘       └──────────┘
             │
    ┌────────┴────────┐
    │                 │
Complete          Abandon
    │                 │
    ▼                 ▼
┌──────────┐    ┌──────────┐
│completed │    │ archived │
└──────────┘    └──────────┘

状态转换表：
┌─────────────────┬───────────────────────────┬──────────────────────────────┐
│ 当前状态        │ 触发事件                   │ 目标状态                     │
├─────────────────┼───────────────────────────┼──────────────────────────────┤
│ (初始)          │ CreateIdea                │ draft                        │
│ draft           │ Publish                   │ published                    │
│ published       │ StartCampaign             │ funding                      │
│ published       │ Archive                   │ archived                     │
│ funding         │ Funded                    │ in_progress                  │
│ funding         │ FundingFailed             │ archived                     │
│ funding         │ CancelCampaign            │ archived                     │
│ in_progress     │ Complete                  │ completed                    │
│ in_progress     │ Abandon                   │ archived                     │
└─────────────────┴───────────────────────────┴──────────────────────────────┘
```

---

## 收益分配机制

```go
// types/revenue.go

// 贡献类型及权重配置
type ContributionType string

const (
    ContributionTypeCode     ContributionType = "code"      // 代码贡献
    ContributionTypeDesign   ContributionType = "design"    // 设计贡献
    ContributionTypeBugFix   ContributionType = "bugfix"    // Bug修复
    ContributionTypeDocs     ContributionType = "docs"      // 文档贡献
    ContributionTypeTesting  ContributionType = "testing"   // 测试贡献
)

// 默认贡献权重配置
// ContributionWeight = Code(30%) + Design(50%) + BugFixes(20%)
var DefaultContributionWeights = map[ContributionType]sdk.Dec{
    ContributionTypeCode:    sdk.NewDecWithPrec(30, 2),   // 30%
    ContributionTypeDesign:  sdk.NewDecWithPrec(50, 2),   // 50%
    ContributionTypeBugFix:  sdk.NewDecWithPrec(20, 2),   // 20%
}

// 贡献记录
type ContributionRecord struct {
    Id          uint64
    IdeaId      uint64
    Contributor sdk.AccAddress

    // 贡献详情
    Type        ContributionType
    Description string
    Weight      sdk.Dec         // 该贡献的权重

    // 验证
    Verified    bool
    VerifiedBy  []sdk.AccAddress

    // 时间
    CreatedAt   time.Time
    VerifiedAt  *time.Time
}

// 收益分配计算
type RevenueDistribution struct {
    IdeaId          uint64
    TotalRevenue    sdk.Coins

    // 各贡献者分配
    Distributions   []ContributorShare
}

type ContributorShare struct {
    Address         sdk.AccAddress
    TotalWeight     sdk.Dec        // 累计贡献权重
    SharePercent    sdk.Dec        // 占比百分比
    Amount          sdk.Coins      // 分配金额
}

// 计算收益分配
// 公式: Share = (个人总权重 / 所有人总权重) * 总收益
func CalculateRevenueDistribution(
    ideaId uint64,
    totalRevenue sdk.Coins,
    contributions []ContributionRecord,
) RevenueDistribution
```

---

## Cosmos SDK 集成

```
x/idea/
├── module.go              # 模块定义
├── keeper/
│   ├── keeper.go          # 状态管理
│   ├── grpc_query.go      # gRPC 查询
│   └── msg_server.go      # 交易处理
├── types/
│   ├── keys.go            # 存储 key
│   ├── types.go           # 类型定义
│   ├── genesis.go         # 创世状态
│   └── msgs.go            # 消息类型
└── client/
    └── cli/               # CLI 命令
```

---

## 核心类型定义

### IdeaStatus - 想法状态

```go
// types/types.go

type IdeaStatus string

const (
    IdeaStatusDraft      IdeaStatus = "draft"
    IdeaStatusPublished  IdeaStatus = "published"
    IdeaStatusFunding    IdeaStatus = "funding"
    IdeaStatusInProgress IdeaStatus = "in_progress"
    IdeaStatusCompleted  IdeaStatus = "completed"
    IdeaStatusArchived   IdeaStatus = "archived"
    IdeaStatusAbandoned  IdeaStatus = "abandoned"
)
```

### Idea - 想法

```go
type Idea struct {
    Id           uint64
    Creator      sdk.AccAddress

    // 基本信息
    Title       string
    Description string  // Markdown
    Summary     string  // 摘要

    // 分类与标签
    Category    string
    Tags        []string

    // 版本控制
    Version       uint64
    VersionHistory []uint64  // IdeaVersion IDs

    // 协作
    Collaborators []IdeaCollaborator
    IsPublic      bool

    // 分解
    Decompositions []uint64  // IdeaDecomposition IDs

    // 众筹
    CampaignId *uint64

    // 状态
    Status     IdeaStatus
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type IdeaCollaborator struct {
    Address     sdk.AccAddress
    Role        string  // "owner" | "editor" | "commenter"
    JoinedAt    time.Time
    Contribution sdk.Dec  // 贡献百分比
}
```

### IdeaVersion - 想法版本

```go
type IdeaVersion struct {
    Id          uint64
    IdeaId      uint64
    Version     uint64

    // 版本内容
    Title       string
    Description string
    Summary     string

    // 变更说明
    ChangeLog   string

    // 元数据
    Creator     sdk.AccAddress
    CreatedAt   time.Time

    // 差异
    Diff        *IdeaDiff
}

type IdeaDiff struct {
    Title       *string  // 标题差异 (nil 表示无变化)
    Description string   // 描述差异 (diff 格式)
}
```

### IdeaDecomposition - 想法分解

```go
type IdeaDecomposition struct {
    Id          uint64
    IdeaId      uint64

    // 分解信息
    Name        string
    Description string

    // 子想法/任务
    Items       []DecompositionItem

    // 依赖关系
    Dependencies []DecompositionDependency

    // 元数据
    Creator     sdk.AccAddress
    CreatedAt   time.Time
}

type DecompositionItem struct {
    Type            string  // "idea" | "task"
    RefId           *uint64 // 已创建的想法/任务ID
    Title           string
    Description     string
    EstimatedValue  *sdk.Coin
    Assignee        *sdk.AccAddress
}

type DecompositionDependency struct {
    From  uint64  // 源项索引
    To    uint64  // 目标项索引
    Type  string  // "requires" | "blocks" | "relates"
}
```

---

## 众筹系统

### Campaign - 众筹

```go
type CampaignType string

const (
    CampaignTypeInvestment CampaignType = "investment"
    CampaignTypeLending    CampaignType = "lending"
    CampaignTypeDonation   CampaignType = "donation"
)

type CampaignStatus string

const (
    CampaignStatusDraft      CampaignStatus = "draft"
    CampaignStatusActive     CampaignStatus = "active"
    CampaignStatusFunded     CampaignStatus = "funded"
    CampaignStatusFailed     CampaignStatus = "failed"
    CampaignStatusInProgress CampaignStatus = "in_progress"
    CampaignStatusCompleted  CampaignStatus = "completed"
    CampaignStatusCancelled  CampaignStatus = "cancelled"
    CampaignStatusDisputed   CampaignStatus = "disputed"
)

type Campaign struct {
    Id          uint64
    IdeaId      uint64
    Creator     sdk.AccAddress

    // 基本信息
    Title       string
    Description string
    Image       *string  // 封面图片 URL

    // 类型与条款
    Type        CampaignType
    Terms       CampaignTerms

    // 目标
    TargetAmount    sdk.Coin
    CurrentAmount   sdk.Coin
    ContributorCount uint64

    // 时间
    StartDate   time.Time
    EndDate     time.Time

    // 状态
    Status      CampaignStatus

    // 资金管理
    EscrowAddress *sdk.AccAddress
    FundsReleased  bool

    // 更新
    Updates     []uint64  // CampaignUpdate IDs

    // 元数据
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### CampaignTerms - 众筹条款

```go
type CampaignTerms struct {
    // 投资型条款
    Investment *InvestmentTerms

    // 借贷型条款
    Lending *LendingTerms

    // 捐赠型条款
    Donation *DonationTerms
}

type InvestmentTerms struct {
    EquityShare    sdk.Dec  // 股权比例 (0-100)
    ProfitShare    sdk.Dec  // 收益分成比例 (0-100)
    MinInvestment  sdk.Coin
    MaxInvestment  sdk.Coin
    VestingPeriod  time.Duration
}

type LendingTerms struct {
    InterestRate        sdk.Dec  // 年利率
    RepaymentPeriod     time.Duration
    RepaymentSchedule   string  // "monthly" | "quarterly" | "yearly" | "lump_sum"
    CollateralRequired  bool
    CollateralDescription *string
}

type DonationTerms struct {
    RewardTiers     []DonationReward
    AllowAnonymous  bool
}

type DonationReward struct {
    MinAmount          sdk.Coin
    Title              string
    Description        string
    EstimatedDelivery  time.Time
}
```

### Contribution - 贡献

```go
type Contribution struct {
    Id          uint64
    CampaignId  uint64
    Contributor sdk.AccAddress

    // 贡献信息
    Amount      sdk.Coin
    Type        CampaignType

    // 条款快照
    TermsSnapshot CampaignTerms

    // 状态
    Status      string  // "pending" | "confirmed" | "refunded" | "claimed"

    // 回报
    RewardTier      *DonationReward
    EquityAllocated *sdk.Dec

    // 交易
    TxHash      *string

    // 元数据
    CreatedAt   time.Time
    ConfirmedAt *time.Time
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
    mqKeeper      mqkeeper.Keeper
}

func (k Keeper) CreateIdea(ctx sdk.Context, msg types.MsgCreateIdea) (uint64, error)
func (k Keeper) UpdateIdea(ctx sdk.Context, msg types.MsgUpdateIdea) error
func (k Keeper) GetIdea(ctx sdk.Context, id uint64) (types.Idea, error)
func (k Keeper) CreateCampaign(ctx sdk.Context, msg types.MsgCreateCampaign) (uint64, error)
func (k Keeper) Contribute(ctx sdk.Context, msg types.MsgContribute) (uint64, error)
func (k Keeper) DecomposeIdea(ctx sdk.Context, msg types.MsgDecomposeIdea) (uint64, error)
func (k Keeper) ConvertToTask(ctx sdk.Context, ideaId uint64, itemIndex int) (uint64, error)
```

---

## gRPC 查询

```protobuf
// query.proto

service Query {
    rpc Idea(QueryIdeaRequest) returns (QueryIdeaResponse);
    rpc Ideas(QueryIdeasRequest) returns (QueryIdeasResponse);
    rpc IdeasByCreator(QueryIdeasByCreatorRequest) returns (QueryIdeasByCreatorResponse);
    rpc Campaign(QueryCampaignRequest) returns (QueryCampaignResponse);
    rpc Campaigns(QueryCampaignsRequest) returns (QueryCampaignsResponse);
    rpc Contributions(QueryContributionsRequest) returns (QueryContributionsResponse);
    rpc Decomposition(QueryDecompositionRequest) returns (QueryDecompositionResponse);
}
```

---

## 模块依赖

```
x/idea (用户插件)
    │
    ├── 核心模块依赖
    │   ├── x/identity  (身份验证)
    │   ├── x/dispute   (Trust System - 众筹信任)
    │   ├── x/escrow    (资金托管)
    │   └── x/bank      (代币转账)
    │
    └── 插件依赖
        └── x/task      (想法→任务转化)
```

---

## 与GenieBot界面交互

```
GenieBot界面 (12-misc)                x/idea (链上)
      │                                      │
      │  1. 用户提交想法到 AI                 │
      │                                      │
      │  2. AI 完善/评估想法                  │
      │                                      │
      │  3. 提交 MsgCreateIdea               │
      │─────────────────────────────────────►│
      │                                      │
      │  4. 监听想法状态变更                  │
      │◄─────────────────────────────────────│
      │                                      │
      │  5. 建议分解方案                      │
      │                                      │
      │  6. 提交 MsgDecomposeIdea            │
      │─────────────────────────────────────►│
```

---

## OpenFang 集成

| OpenFang 组件 | 想法系统集成 |
|--------------|-------------|
| Collector Hand | 想法收集、市场调研 |
| Researcher Agent | 想法评估、可行性分析 |
| GenieBot Agent | 想法孵化、协作管理 |

---

[上一章：Oracle服务](./06-exchange.md) | [返回索引](./00-index.md) | [下一章：x/task →](./08-task.md)
