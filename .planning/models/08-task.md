# x/task - 任务模块

> **模块类型:** 用户插件（GenieBot插件的一部分）
> **技术栈:** Go (Cosmos SDK)
> **位置:** `src/chain/x/task`
> **依赖:** 服务市场 (11-service)、身份账号 (10-identity)、Trust System (09-dispute)

---

## 概述

x/task 是 ShareTokens 链上的任务市场模块，作为 GenieBot 插件的一部分，负责任务发布、申请、执行、评价的全流程管理。作为 Cosmos SDK 模块实现，与 x/idea、x/escrow、x/dispute 等模块紧密协作。

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
│  └── GenieBot界面 (12-misc) ──┬── 想法系统 (07-idea)                      │
│                               │                                           │
│                               └── 任务市场 (08-task) ← 本章               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 任务状态机

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          任务状态机                                           │
└─────────────────────────────────────────────────────────────────────────────┘

                            ┌──────────┐
                            │  draft   │  草稿
                            └────┬─────┘
                                 │
                           Publish
                                 │
                                 ▼
                            ┌──────────┐
                            │   open   │  开放申请 (7天无人申请自动关闭)
                            └────┬─────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
         AssignTask       ApplicationTimeout  CancelTask
              │                  │                  │
              ▼                  ▼                  ▼
        ┌──────────┐      ┌──────────┐       ┌──────────┐
        │ assigned │      │ expired  │       │cancelled │
        └────┬─────┘      └──────────┘       └──────────┘
             │
        StartWork
             │
             ▼
       ┌───────────┐
       │in_progress│  执行中 (14天未提交自动失败)
       └─────┬─────┘
             │
    ┌────────┼────────┐
    │        │        │
SubmitWork│   │   ExecutionTimeout
    │        │        │
    ▼        │        ▼
┌──────────┐ │  ┌──────────┐
│under_review  │  │  failed  │
└────┬─────┘ │  └──────────┘
     │       │
     │ ReviewTimeout (3天自动通过)
     │       │
     ├───────┼──────────────┐
     │       │              │
 Approve  RequestRevision  │
     │       │              │
     ▼       ▼              ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│completed │ │revision_requested│ │disputed │
└──────────┘ └────┬─────┘ └──────────┘
                  │
             Resubmit
                  │
                  ▼
           ┌───────────┐
           │under_review│
           └───────────┘

状态转换表：
┌─────────────────┬───────────────────────────┬──────────────────────────────┐
│ 当前状态        │ 触发事件                   │ 目标状态                     │
├─────────────────┼───────────────────────────┼──────────────────────────────┤
│ (初始)          │ CreateTask                │ draft                        │
│ draft           │ Publish                   │ open                         │
│ open            │ AssignTask                │ assigned                     │
│ open            │ ApplicationTimeout (7天)  │ expired                      │
│ open            │ CancelTask                │ cancelled                    │
│ assigned        │ StartWork                 │ in_progress                  │
│ in_progress     │ SubmitWork                │ under_review                 │
│ in_progress     │ ExecutionTimeout (14天)   │ failed                       │
│ under_review    │ Approve                   │ completed                    │
│ under_review    │ RequestRevision           │ revision_requested           │
│ under_review    │ ReviewTimeout (3天)       │ completed (自动通过)         │
│ under_review    │ CreateDispute             │ disputed                     │
│ revision_requested │ Resubmit              │ under_review                 │
└─────────────────┴───────────────────────────┴──────────────────────────────┘
```

---

## 自动超时配置

```go
// types/timeout.go

type TaskTimeoutConfig struct {
    // 申请超时: 7天无人申请自动关闭
    ApplicationTimeout time.Duration  // 7 * 24 * time.Hour

    // 执行超时: 14天未提交自动失败
    ExecutionTimeout time.Duration    // 14 * 24 * time.Hour

    // 评审超时: 3天未评审自动通过
    ReviewTimeout time.Duration       // 3 * 24 * time.Hour

    // 延期配置
    MaxExtensions     uint64          // 最大延期次数
    ExtensionDuration time.Duration   // 每次延期时长
}

// 默认配置
var DefaultTaskTimeoutConfig = TaskTimeoutConfig{
    ApplicationTimeout: 7 * 24 * time.Hour,   // 7天
    ExecutionTimeout:   14 * 24 * time.Hour,  // 14天
    ReviewTimeout:      3 * 24 * time.Hour,   // 3天
    MaxExtensions:      2,                    // 最多2次
    ExtensionDuration:  7 * 24 * time.Hour,   // 每次7天
}
```

---

## Cosmos SDK 集成

```
x/task/
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

### TaskStatus - 任务状态

```go
// types/types.go

type TaskStatus string

const (
    TaskStatusDraft           TaskStatus = "draft"
    TaskStatusOpen            TaskStatus = "open"
    TaskStatusAssigned        TaskStatus = "assigned"
    TaskStatusInProgress      TaskStatus = "in_progress"
    TaskStatusUnderReview     TaskStatus = "under_review"
    TaskStatusRevisionRequested TaskStatus = "revision_requested"
    TaskStatusCompleted       TaskStatus = "completed"
    TaskStatusCancelled       TaskStatus = "cancelled"
    TaskStatusDisputed        TaskStatus = "disputed"
    TaskStatusExpired         TaskStatus = "expired"
)
```

### Task - 任务

```go
type Task struct {
    Id       uint64
    Creator  sdk.AccAddress

    // 基本信息
    Title       string
    Description string  // Markdown
    Summary     string

    // 分类与技能
    Category   string
    Skills     []string
    Difficulty string  // "beginner" | "intermediate" | "advanced" | "expert"

    // 价值与奖励
    Budget TokenAmount
    Reward TokenAmount

    // 人员配置
    MaxAssignees    uint64
    CurrentAssignees []sdk.AccAddress

    // 指派配置
    AssignmentConfig TaskAssignmentConfig

    // 里程碑
    Milestones      []TaskMilestone
    CurrentMilestone uint64  // 当前里程碑索引

    // 托管
    Escrow *EscrowInfo

    // 时间
    Deadline          *time.Time
    EstimatedDuration time.Duration
    TimeoutConfig     TimeoutConfig

    // 状态
    Status       TaskStatus
    StatusHistory []StatusChange

    // 申请与竞标
    Applications []TaskApplication
    Bids         []TaskBid

    // 提交与评审
    Submissions []TaskSubmission
    Reviews     []TaskReview

    // 协作任务
    Collaborative *CollaborativeTask

    // 来源
    IdeaId *uint64  // 关联的想法 (可选)

    // 附件
    Attachments []Attachment

    // 元数据
    CreatedAt  time.Time
    UpdatedAt  time.Time
    CompletedAt *time.Time
}
```

### TaskApplication - 任务申请

```go
type TaskApplication struct {
    Id      uint64
    TaskId  uint64
    Applicant sdk.AccAddress

    // 申请内容
    Message         string
    ProposedTimeline time.Duration
    ProposedBudget  *TokenAmount  // 可协商

    // 附件
    Portfolio []PortfolioItem

    // 状态
    Status string  // "pending" | "accepted" | "rejected" | "withdrawn"

    // 元数据
    CreatedAt   time.Time
    RespondedAt *time.Time
    RespondedBy *sdk.AccAddress
}
```

### TaskMilestone - 任务里程碑

```go
type TaskMilestone struct {
    Id     uint64
    TaskId uint64

    // 里程碑信息
    Title       string
    Description string
    Order       uint64

    // 价值
    Amount    TokenAmount
    Percentage sdk.Dec  // 占总预算百分比

    // 交付物
    Deliverables []Deliverable

    // 状态
    Status string  // "pending" | "in_progress" | "submitted" | "reviewing" | "approved" | "rejected"

    // 时间
    Deadline   *time.Time
    StartedAt  *time.Time
    SubmittedAt *time.Time
    ApprovedAt *time.Time

    // 评审
    Review *TaskReview

    // 交付
    Submission *TaskSubmission
}
```

### TaskSubmission - 任务提交

```go
type TaskSubmission struct {
    Id         uint64
    TaskId     uint64
    MilestoneId *uint64
    Submitter  sdk.AccAddress

    // 提交内容
    Message      string
    Deliverables []SubmittedDeliverable

    // 版本
    Version           uint64
    PreviousSubmission *uint64

    // 状态
    Status string  // "pending" | "reviewing" | "approved" | "revision_requested" | "rejected"

    // 元数据
    CreatedAt  time.Time
    ReviewedAt *time.Time
}
```

### TaskReview - 任务评审

```go
type TaskReview struct {
    Id          uint64
    TaskId      uint64
    SubmissionId uint64
    MilestoneId *uint64
    Reviewer    sdk.AccAddress

    // 评审结果
    Result string  // "approved" | "revision_requested" | "rejected"

    // 评审意见
    Comment       string
    RevisionNotes *string

    // 评分 (任务完成后)
    Ratings *TaskRating

    // 状态
    Status string  // "pending" | "completed" | "disputed"

    // 元数据
    CreatedAt time.Time
}

type TaskRating struct {
    // 各维度评分 (1-5)
    Quality       sdk.Dec
    Communication sdk.Dec
    Timeliness    sdk.Dec
    Expertise     sdk.Dec

    // 综合评分
    Overall sdk.Dec

    // 评价文字
    Review string

    // 是否公开
    IsPublic bool
}
```

---

## 里程碑状态机

```
       ┌─────────┐
       │ pending │  等待开始
       └────┬────┘
            │
       StartMilestone
            │
            ▼
      ┌───────────┐
      │in_progress│  进行中
      └─────┬─────┘
            │
      SubmitMilestone
            │
            ▼
       ┌──────────┐
       │submitted │  已提交
       └────┬─────┘
            │
       StartReview
            │
            ▼
       ┌──────────┐
       │reviewing │  评审中
       └────┬─────┘
            │
    ┌───────┴───────┐
    │               │
Approve         Reject
    │               │
    ▼               ▼
┌─────────┐   ┌──────────┐
│approved │   │ rejected │
└─────────┘   └──────────┘
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
    ideaKeeper    ideakeeper.Keeper
}

func (k Keeper) CreateTask(ctx sdk.Context, msg types.MsgCreateTask) (uint64, error)
func (k Keeper) UpdateTask(ctx sdk.Context, msg types.MsgUpdateTask) error
func (k Keeper) GetTask(ctx sdk.Context, id uint64) (types.Task, error)
func (k Keeper) ApplyForTask(ctx sdk.Context, msg types.MsgApplyForTask) (uint64, error)
func (k Keeper) SubmitBid(ctx sdk.Context, msg types.MsgSubmitBid) (uint64, error)
func (k Keeper) AcceptApplication(ctx sdk.Context, taskId, appId uint64) error
func (k Keeper) SubmitWork(ctx sdk.Context, msg types.MsgSubmitWork) (uint64, error)
func (k Keeper) ReviewSubmission(ctx sdk.Context, msg types.MsgReviewSubmission) error
func (k Keeper) CompleteTask(ctx sdk.Context, taskId uint64) error
func (k Keeper) CancelTask(ctx sdk.Context, taskId uint64, reason string) error
```

---

## gRPC 查询

```protobuf
// query.proto

service Query {
    rpc Task(QueryTaskRequest) returns (QueryTaskResponse);
    rpc Tasks(QueryTasksRequest) returns (QueryTasksResponse);
    rpc TasksByCreator(QueryTasksByCreatorRequest) returns (QueryTasksByCreatorResponse);
    rpc TasksByAssignee(QueryTasksByAssigneeRequest) returns (QueryTasksByAssigneeResponse);
    rpc Applications(QueryApplicationsRequest) returns (QueryApplicationsResponse);
    rpc Bids(QueryBidsRequest) returns (QueryBidsResponse);
    rpc Submissions(QuerySubmissionsRequest) returns (QuerySubmissionsResponse);
    rpc Reviews(QueryReviewsRequest) returns (QueryReviewsResponse);
    rpc Categories(QueryCategoriesRequest) returns (QueryCategoriesResponse);
    rpc Skills(QuerySkillsRequest) returns (QuerySkillsResponse);
}
```

---

## 模块依赖

```
x/task (用户插件)
    │
    ├── 核心模块依赖
    │   ├── x/identity  (身份验证)
    │   ├── x/dispute   (Trust System - 信任加权指派)
    │   ├── x/escrow    (资金托管)
    │   ├── x/dispute   (Trust System - 争议处理)
    │   └── x/bank      (代币转账)
    │
    └── 插件依赖
        └── x/idea      (想法→任务转化)
```

---

## 与GenieBot界面交互

```
GenieBot界面 (12-misc)                x/task (链上)
      │                                      │
      │  1. 用户通过 AI 创建任务              │
      │                                      │
      │  2. 提交 MsgCreateTask               │
      │─────────────────────────────────────►│
      │                                      │
      │  3. 监听 TaskCreated 事件             │
      │◄─────────────────────────────────────│
      │                                      │
      │  4. AI 匹配合适的执行者               │
      │                                      │
      │  5. 通知匹配到的用户                  │
      │                                      │
      │  6. 监听任务状态变更                  │
      │◄─────────────────────────────────────│
      │                                      │
      │  7. 自动化工作流执行                  │
      │                                      │
      │  8. 提交工作成果                      │
      │─────────────────────────────────────►│
```

---

## OpenFang 集成

| OpenFang 组件 | 任务市场集成 |
|--------------|-------------|
| Coder Agent | 软件开发任务执行 |
| Content Hand | 内容创作任务执行 |
| Lead Hand | 任务匹配、资源调度 |
| GenieBot Agent | 任务创建、进度跟踪 |

---

[上一章：x/idea](./07-idea.md) | [返回索引](./00-index.md) | [下一章：x/dispute →](./09-dispute.md)
