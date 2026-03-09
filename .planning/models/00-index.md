# ShareTokens 数据模型索引

> 基于成熟框架构建的业务层数据模型

---

## 技术栈

```
┌─────────────────────────────────────────────────────────────┐
│                    ShareTokens 技术栈                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  应用层（GenieBot - AI对话界面）                              │
│  ├── AI对话、想法孵化、资源匹配                              │
│  └── Workflow执行                                           │
│                                                             │
│  AI Agent 运行时层（OpenFang 框架）                          │
│  ├── 28+ Agent模板：coder, researcher, writer, architect   │
│  ├── 7个Hands：Collector, Clip, Lead, Content, Trade...    │
│  ├── 16层安全防护、WASM沙箱                                 │
│  └── Dashboard: localhost:4200                              │
│                                                             │
│  基础设施层（使用成熟框架，无需自行开发）                      │
│  ├── CometBFT：共识、出块、验证者管理（内置）                │
│  ├── Cosmos SDK：状态机、账户管理、交易处理                  │
│  ├── Keplr：钱包集成（浏览器插件）                           │
│  └── Chainlink：价格预言机、自动化任务                       │
│                                                             │
│  业务层（本项目的核心创新）                                   │
│  ├── Trust System：德商(MQ)评分、零和博弈、争议仲裁              │
│  ├── 实名制身份：哈希存储、隐私保护、防重复注册              │
│  ├── 算力交易：API Key 托管、托管支付、争议仲裁              │
│  ├── 想法众筹：想法孵化、贡献追踪、Token 激励                │
│  └── 任务市场：任务分配、里程碑、评价体系                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### OpenFang 集成映射

| ShareTokens 组件 | OpenFang 组件 | 说明 |
|------------------|---------------|------|
| AI对话界面(GenieBot) | Agent + Channels | 多渠道对话支持 |
| 想法收集 | Collector Hand | 自动数据收集和监控 |
| 软件开发Workflow | Coder Agent + Debugger Agent | 代码生成和调试 |
| 内容创作Workflow | Writer Agent + Content Hand | 内容创作和发布 |
| 想法评估 | Researcher Agent | 多AI评估和分析 |
| 资源匹配 | Lead Hand | 销售线索和客户管理 |
| API Key安全托管 | OpenFang 16层安全 | WASM沙箱隔离 |

### OpenFang Agent 生态

| Agent/Hand | 用途 | ShareTokens 应用 |
|------------|------|------------------|
| **GenieBot Agent** (自定义) | 主对话Agent | 用户入口，想法孵化 |
| **Collector Hand** | 数据收集 | 想法收集、市场调研 |
| **Coder Agent** | 代码编写 | 软件开发Workflow |
| **Content Hand** | 内容创作 | 内容创作Workflow |
| **Lead Hand** | 销售线索 | 资源匹配 |
| **Researcher Agent** | 研究分析 | 想法评估 |

---

## 模块列表

### 核心模块（每个节点必须有）

| 文件 | 模块 | 说明 | 依赖 | OpenFang集成 |
|------|------|------|------|--------------|
| [01-base.md](./01-base.md) | 基础类型 | 原始类型、Token、加密 | - | Agent类型定义 |
| [04-consensus.md](./04-consensus.md) | P2P通信 | 共识、区块、验证者、网络 | Cosmos SDK + CometBFT | - |
| [10-identity.md](./10-identity.md) | 身份账号 | 实名制、身份注册表、账户 | 01-base | Agent身份绑定 |
| [10-identity.md](./10-identity.md) | 钱包 | 账户管理、余额、交易 | Cosmos SDK Auth + Keplr | - |
| [11-service.md](./11-service.md) | 服务市场 | 三层服务市场（LLM/Agent/Workflow） | 01-base, 04-consensus | OpenFang Services |
| [11-service.md](./11-service.md) | 托管支付 | 托管账户、支付、结算 | 01-base, 04-consensus | - |
| [09-dispute.md](./09-dispute.md) | Trust System | 信誉评分、零和博弈、争议仲裁 | 01-base, 10-identity | - |

### 可选模块（插件）

#### 服务提供者插件

| 文件 | 模块 | 说明 | 依赖 | OpenFang集成 |
|------|------|------|------|--------------|
| [05-compute.md](./05-compute.md) | LLM API Key托管 | API Key加密存储、请求 | 11-service (服务市场) | OpenFang Provider |
| [05-compute.md](./05-compute.md) | Agent执行器 | Agent运行时、资源管理 | 11-service (服务市场) | OpenFang Agent Runtime |
| [05-compute.md](./05-compute.md) | Workflow执行器 | Workflow编排、执行 | 11-service (服务市场) | OpenFang Workflow |

### 辅助模块

| 文件 | 模块 | 说明 | 依赖 | OpenFang集成 |
|------|------|------|------|--------------|
| [06-exchange.md](./06-exchange.md) | 汇率层 | 汇率快照、价格映射 | Chainlink | - |
| [07-idea.md](./07-idea.md) | 想法系统 | 想法、众筹、贡献 | 11-service (服务市场) | Collector Hand |
| [08-task.md](./08-task.md) | 任务市场 | 任务、申请、里程碑 | 11-service (服务市场) | Coder Agent, Content Hand |

---

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                    ShareTokens 架构                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  核心模块（每个节点必须有）                                   │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  ├── P2P通信（Cosmos SDK + CometBFT）                       │
│  ├── 身份账号（实名制、身份注册表）                          │
│  ├── 钱包（Cosmos SDK Auth + Keplr）                        │
│  ├── 服务市场（三层服务：LLM/Agent/Workflow）                │
│  │   ├── 服务注册与发现                                     │
│  │   ├── 定价与汇率（Chainlink）                            │
│  │   └── 计费与结算                                         │
│  ├── 托管支付（争议时冻结）                                  │
│  └── Trust System（信誉评分、争议仲裁、裁决算法）            │
│                                                             │
│  可选模块（插件）                                            │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│                                                             │
│  [服务提供者插件]                                            │
│  ├── LLM API Key托管（加密存储、OpenFang 16层安全）         │
│  ├── Agent执行器（28+ Agent模板、7个Hands）                 │
│  └── Workflow执行器（编排、执行、监控）                      │
│                                                             │
│  [用户插件]                                                  │
│  └── GenieBot界面（AI对话、想法孵化、资源匹配）              │
│      *详见 11-service.md 的 GenieBot 部分*                  │
│                                                             │
│  OpenFang 集成层                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  ├── 28+ Agent模板 (coder, researcher, writer...)          │
│  ├── 7个Hands (Collector, Clip, Lead, Content...)          │
│  ├── 16层安全防护、WASM沙箱                                 │
│  └── Dashboard (localhost:4200)                             │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 模块依赖关系

```
核心模块依赖图：

┌──────────────────────────────────────────────────────────┐
│                    基础类型（01-base）                     │
└────────────────┬─────────────────────────────────────────┘
                 │
        ┌────────┴────────┐
        ▼                 ▼
┌───────────────┐  ┌──────────────┐
│  P2P通信      │  │  身份账号     │
│  (04-consensus)│  │  (10-identity)│
└───────┬───────┘  └──────┬───────┘
        │                  │
        └────────┬─────────┘
                 ▼
         ┌───────────────┐
         │   钱包         │
         │  (Keplr集成)   │
         └───────┬───────┘
                 │
                 ▼
         ┌───────────────┐
         │  服务市场      │
         │  (11-service)  │
         │  - LLM服务     │
         │  - Agent服务   │
         │  - Workflow服务│
         └───────┬───────┘
                 │
        ┌────────┴────────┐
        ▼                 ▼
┌───────────────┐  ┌──────────────┐
│  托管支付      │  │ Trust System  │
│  (11-service)  │  │  (09-dispute) │
└───────┬───────┘  └──────┬───────┘
        │                  │
        └────────┬─────────┘
                 ▼
         ┌───────────────┐
         │ Trust System   │
         │  (09-dispute)  │
         └───────────────┘

可选插件依赖：

服务提供者插件：
  LLM API Key托管 ──────┐
  Agent执行器      ──────┼──> 依赖服务市场（11-service）
  Workflow执行器   ──────┘

用户插件：
  GenieBot界面 ────────────> 依赖服务市场（11-service）
```

---

## 核心创新

### 1. 三层服务市场（核心模块）

```
服务市场 = {
  Layer 1: LLM服务
    - API Key托管
    - 标准化定价
    - 请求转发

  Layer 2: Agent服务
    - Agent执行器
    - 资源管理
    - 结果返回

  Layer 3: Workflow服务
    - Workflow编排
    - 多Agent协作
    - 复杂任务执行
}

统一功能：
  - 服务注册与发现
  - 定价与汇率（Chainlink）
  - 托管支付与结算
  - 争议仲裁
```

### 2. Trust System

```
- 每个用户初始 MQ = 100
- 总 MQ 恒定（零和博弈）
- 评分权重 = 参与者 MQ / 所有参与者 MQ 之和
- 每次评分重新分配 3%
- 不活跃账户 MQ 衰减
- 争议仲裁与 MQ 系统整合
```

### 3. 实名制身份

```
- 严格实名制（微信、GitHub 等）
- 只存哈希，不存明文
- 全局身份注册表防止重复注册
- 本地 Merkle 证明验证
```

### 3. 托管支付

```
- 托管账户管理
- 争议时自动冻结
- MQ 权重分配
- 多币种支持
```

### 4. 模块化插件架构

```
核心节点 = {
  P2P通信 + 身份账号 + 钱包 + 服务市场 + 托管支付 + Trust System
}

服务提供者节点 = 核心节点 + {
  LLM API Key托管插件
  Agent执行器插件
  Workflow执行器插件
}

用户节点 = 核心节点 + {
  GenieBot界面插件
}
```

### 5. 想法众筹

```
- 想法孵化（Token 支持）
- 贡献追踪（权重分配）
- 收益分享（智能合约执行）
```

---

## 节点类型

| 类型 | 核心模块 | 可选插件 | 适用场景 |
|------|----------|----------|----------|
| 轻节点 | ✓ 6个核心模块 | GenieBot界面插件 | 普通用户 |
| 全节点 | ✓ 6个核心模块 | - | 验证者 |
| 服务节点 | ✓ 6个核心模块 | LLM API Key托管插件<br>Agent执行器插件<br>Workflow执行器插件 | 服务提供者 |
| 归档节点 | ✓ 6个核心模块 | - | 区块浏览器 |

**核心模块（每个节点必须有）：**
1. P2P通信（Cosmos SDK + CometBFT）
2. 身份账号（实名制、身份注册表）
3. 钱包（Cosmos SDK Auth + Keplr）
4. 服务市场（三层服务：LLM/Agent/Workflow）
5. 托管支付（争议时冻结）
6. Trust System（德商MQ评分、争议仲裁）

**可选插件（按需安装）：**
- 服务提供者插件：LLM API Key托管、Agent执行器、Workflow执行器（详见 05-compute.md）
- 用户插件：GenieBot界面（详见 11-service.md）

---

## 文件统计

| 模块 | 简化前行数 | 简化后行数 | 减少 |
|------|-----------|-----------|------|
| 02-network | ~736 | ~287 | 61% |
| 04-consensus | ~254 | ~254 | - |
| 06-exchange | ~126 | ~123 | 2% |

---

## 变更记录

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| 1.0 | 2024-03-01 | 初始版本 |
| 1.1 | 2024-03-01 | 拆分为多模块 |
| 1.2 | 2024-03-01 | 基于 libp2p/Cosmos SDK/Chainlink 简化 |
| 1.3 | 2024-03-02 | 移除过时模块（P2P网络、钱包），更新为 Cosmos SDK + CometBFT 架构 |
| 1.4 | 2025-03-02 | **集成 OpenFang Agent OS**：添加Agent层、28+模板、7个Hands、16层安全 |
| 1.5 | 2025-03-02 | **模块化架构重组**：划分为7个核心模块+可选插件，服务市场统一为三层服务（LLM/Agent/Workflow） |
| 1.6 | 2025-03-03 | **清理重复定义**：删除 12-misc.md，MQ系统统一到 09-dispute.md，FraudIndicators 迁移到 05-compute.md，UserLimits 迁移到 10-identity.md |
