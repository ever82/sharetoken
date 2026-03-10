# ShareTokens Developer Achievements

> 开发视角下的可验收成果清单。每个 Achievement 代表一个较大的功能里程碑，可独立验收。
>
> **优先级定义:**
> - **P0**: 核心基础，必须最先完成
> - **P1**: 核心功能，MVP 必需
> - **P2**: 重要功能，完整体验
> - **P3**: 增强功能，可后期迭代

---

## P0 - 核心基础 ✅ 已完成

> **状态**: P0全部完成，详见 `docs/achievements/done/` 目录
>
> **完成日期**: 2026-03-10

### ACH-DEV-001. Development Infrastructure ✅

**优先级:** P0
**描述:** 开发基础设施与工具链搭建。

**验收标准:**
- [x] CI/CD Pipeline 配置完成（测试、构建、部署）
- [x] 本地开发网络一键启动脚本
- [x] 代码规范与 Lint 配置
- [~] Protobuf 定义完成并生成 Go/TypeScript 代码（Go✅ TS手动实现）

**关联 Spec:** ARCH-001
**验收文档**: `docs/achievements/done/issue-001.md`

---

### ACH-DEV-002. Blockchain Network Foundation ✅

**优先级:** P0
**描述:** 搭建基于 Cosmos SDK + CometBFT 的区块链网络基础设施。

**验收标准:**
- [x] 单节点可启动并出块（出块时间 <= 2s）
- [x] 多节点组成网络（4节点），共识正常运行
- [x] P2P 节点发现与通信正常
  - [~] 新节点加入后 30s 内完成发现（基础验证通过）
  - [~] 1000 条消息广播无丢失（未压力测试）
- [x] 支持 UPnP 自动端口映射(家庭网络用户可自动连网)
- [~] 支持手动配置端口映射(未提供专门文档)
- [x] Noise Protocol 加密通信验证（默认启用）
- [⏭️] 区块浏览器可查看链上数据（延后到部署阶段）

**关联 Spec:** CORE-001, ARCH-001
**验收文档**: `docs/achievements/done/issue-002.md`

---

### ACH-DEV-003. Wallet & Token System ✅

**优先级:** P0
**描述:** 实现 STT 代币的钱包与转账功能。

**验收标准:**
- [x] STT 代币定义与发行
- [x] 余额查询接口可用
- [x] 转账交易签名与广播正常
- [~] Keplr 钱包集成（桌面端）- 代码完成，运行时测试延后
- [~] WalletConnect 支持（移动端）- 代码完成，运行时测试延后
- [x] 交易历史查询

**关联 Spec:** CORE-003
**验收文档**: `docs/achievements/done/issue-003.md`

---

### ACH-DEV-021. Desktop Application (桌面应用)

**优先级:** P0
**描述:** 实现开箱即用的桌面图形界面应用，让用户无需技术背景即可使用。

**为什么这是 P0:**
- 命令行工具 (`sharetokend`) 只适合开发者，普通用户无法使用
- Web 前端 (`npm run serve`) 需要安装 Node.js，门槛太高
- Keplr 插件只能提供钱包功能，无法承载完整平台功能（服务市场、AI 对话等）
- 现代用户预期：下载 → 解压 → 双击打开，三步完成

**验收标准:**
- [ ] 支持 Windows (`.exe` 或 `.zip`)
- [ ] 支持 macOS (`.app` 或 `.dmg`)
- [ ] 支持 Linux (`.AppImage` 或 `.tar.gz`)
- [ ] 内置区块链节点（轻节点模式），无需单独启动
- [ ] 内置钱包功能（创建/导入/转账）
- [ ] 内置服务市场浏览器
- [ ] 内置 GenieBot AI 对话界面
- [ ] 自动更新机制
- [ ] CI/CD 自动构建和发布

**技术方案选型:**
| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **Electron + Go 后端** | 成熟生态、Vue 代码可复用 | 包体积大 (~150MB) | ⭐⭐⭐⭐⭐ |
| **Wails (Go + WebView)** | 包体积小 (~20MB)、性能高 | 生态较新 | ⭐⭐⭐⭐ |
| **Tauri (Rust + WebView)** | 包体积极小 (~5MB)、安全 | 需要 Rust | ⭐⭐⭐ |

**选定方案:** Electron + Go 嵌入式节点
- 前端：复用现有的 Vue 代码
- 后端：嵌入式轻节点，使用 `sharetokend` 作为库调用
- 通信：HTTP API 或 WebSocket

**依赖:** ACH-DEV-003 (Wallet), ACH-DEV-004 (Identity)

**关联 Spec:** CORE-003, REQ-001, NODE-002

---

## P1 - 核心功能

### ACH-DEV-004. Identity Module

**优先级:** P1
**描述:** 实现用户身份注册与实名验证系统。

**验收标准:**
- [ ] 用户可注册链上账户（地址生成）
- [ ] 支持第三方实名验证（WeChat/GitHub/Google）
- [ ] 验证结果以 Hash 形式存储，无明文上链
- [ ] 全局身份注册表防止重复注册
- [ ] 本地 Merkle 证明验证可用
- [ ] 用户限额系统（交易/提现/争议/服务限制）
  - 交易限额：单笔/日/月
  - 提现限额：日限额、冷却期
  - 争议限额：最大活跃争议数
  - 服务限额：并发调用、速率限制

**关联 Spec:** CORE-002, IDENTITY-001, IDENTITY-002

---

### ACH-DEV-005. Escrow Payment System

**优先级:** P1
**描述:** 实现交易资金托管与释放机制。

**验收标准:**
- [ ] 任务开始前锁定 STT 到托管账户
- [ ] 任务完成后自动释放资金给提供者
- [ ] 争议发起时冻结资金
- [ ] 争议解决后按比例分配
- [ ] 多签托管账户安全验证（明确签名配置）

**关联 Spec:** CORE-005, TRUST-004

---

### ACH-DEV-006. Oracle Service

**优先级:** P1
**描述:** 实现去中心化价格数据服务，为服务定价提供汇率转换。

**验收标准:**
- [ ] Chainlink 集成获取汇率数据（明确 Cosmos 集成方案）
- [ ] 各 LLM 官方价格标准化为 STT
- [ ] 价格订阅与缓存机制
- [ ] 价格数据上链可验证
- [ ] 支持价格更新频率配置（默认 5 分钟）

**关联 Spec:** SERVICE-002

---

### ACH-DEV-007. Trust System - MQ Scoring

**优先级:** P1
**描述:** 实现信誉评分系统（MQ）。

**验收标准:**
- [ ] MQ 评分系统：初始100，零和博弈
- [ ] 单次争议 MQ 最大损失 3%，不低于0
- [ ] MQ 加权投票权重计算
- [ ] **收敛机制**：高 MQ 用户偏离共识时损失更多
- [ ] 模拟 100 次争议，验证 MQ 变更符合预期

**关联 Spec:** CORE-006, TRUST-001

---

### ACH-DEV-008. Trust System - Dispute Arbitration

**优先级:** P1
**描述:** 实现去中心化争议仲裁系统。

**验收标准:**
- [ ] AI 调解：对话、证据收集、提案评分（明确链上/链下实现）
- [ ] 陪审团投票：MQ 加权随机抽取
- [ ] MQ 再分配：偏离共识者惩罚，接近共识者奖励
- [ ] 争议全流程可追溯
- [ ] 争议状态变更通知

**依赖:** ACH-DEV-007 (MQ Scoring)

**关联 Spec:** CORE-006, TRUST-002, TRUST-003

---

### ACH-DEV-009. Service Marketplace Core

**优先级:** P1
**描述:** 实现三层服务市场的核心交易逻辑。

**验收标准:**
- [ ] 服务注册与发现接口
- [ ] Level 1 服务：LLM API 按token计费
- [ ] Level 2 服务：Agent 按skill计费
- [ ] Level 3 服务：Workflow 按里程碑打包计费
- [ ] **三种定价模式**: Fixed（固定价）、Dynamic（动态价）、Auction（竞价）
- [ ] 智能路由：自动匹配最优服务提供者
  - 明确路由策略：MQ优先/价格优先/能力匹配
- [ ] 服务提供者管理（注册/下线/状态）

**依赖:** ACH-DEV-004 (Identity), ACH-DEV-005 (Escrow), ACH-DEV-006 (Oracle)

**关联 Spec:** CORE-004, MARKET-001~004

**边界说明:**
- 本成就仅包含链上服务注册、发现、定价、路由逻辑
- 实际的服务执行由 Provider Plugins (ACH-DEV-011/012/013) 完成

---

### ACH-DEV-010. Testnet Launch

**优先级:** P1
**描述:** 部署公共测试网络，供早期用户和开发者测试。

**验收标准:**
- [ ] 至少 4 个验证者节点运行
- [ ] 公开 RPC/LCD 端点可用
- [ ] 区块浏览器可访问
- [ ] 测试代币水龙头可用
- [ ] 网络稳定性 7 天无重大故障

**关联 Spec:** CORE-001

---

## P2 - 重要功能

### ACH-DEV-011. LLM API Key Custody Plugin

**优先级:** P2
**描述:** 实现 LLM Provider API Key 的安全托管与代理服务。

**验收标准:**
- [ ] API Key 加密存储（明确：链上存加密后hash）
- [ ] WASM 沙箱内解密使用
- [ ] 使用后立即清除（Secret Zeroization）
- [ ] 访问控制与定价配置
- [ ] 支持 OpenAI / Anthropic API 代理
- [ ] 密钥管理方案文档化（KMS/HSM 集成说明）

**关联 Spec:** PLUGIN-001, SERVICE-001

---

### ACH-DEV-012. Agent Executor Plugin

**优先级:** P2
**描述:** 集成 OpenFang 作为 Level 2 Agent 执行器。

**验收标准:**
- [ ] OpenFang Rust 运行时集成
- [ ] 28+ Agent 模板可用（coder/researcher/writer...）
- [ ] 16 层安全机制生效（列出关键安全层）
- [ ] WASM 沙箱隔离运行
- [ ] 与 Service Marketplace 对接：接收任务、返回结果
- [ ] Sidecar 部署模式与链节点通信

**关联 Spec:** PLUGIN-001, SERVICE-003

---

### ACH-DEV-013. Workflow Executor Plugin

**优先级:** P2
**描述:** 实现多 Agent 协作的 Workflow 执行引擎。

**验收标准:**
- [ ] OpenFang Hands 集成
- [ ] 7 个自主能力包可用（Collector/Lead/Researcher...）
- [ ] 软件开发 Workflow 可执行
- [ ] 内容创作 Workflow 可执行
- [ ] 里程碑追踪与进度报告
- [ ] 与 Escrow 系统联动（里程碑付款）

**关联 Spec:** PLUGIN-001, MARKET-003~004

---

### ACH-DEV-014. GenieBot User Interface

**优先级:** P2
**描述:** 实现面向用户的 AI 对话与服务调用界面。

**验收标准:**
- [ ] React + TypeScript 前端可运行
- [ ] 自然语言对话入口
- [ ] 意图识别与服务推荐（准确率 >= 85%）
- [ ] 一键调用 LLM/Agent/Workflow 服务
- [ ] 任务管理与进度追踪
- [ ] 结果展示与下载

**关联 Spec:** PLUGIN-002, GENIE-001~002, NODE-002

---

### ACH-DEV-015. Task Marketplace Module

**优先级:** P2
**描述:** 实现任务全生命周期管理（人工任务市场）。

**验收标准:**
- [ ] 任务创建与分解
- [ ] 开放申请/竞价两种模式
- [ ] 里程碑定义与阶段性交付
- [ ] 多维评分（质量/沟通/时效/专业度）
- [ ] 任务历史与统计

**关联 Spec:** MARKET-002

**边界说明:**
- 本成就针对人工任务市场
- AI Agent 执行由 ACH-DEV-012 负责

---

### ACH-DEV-016. Idea & Crowdfunding System

**优先级:** P2
**描述:** 实现创意孵化与众筹平台。

**验收标准:**
- [ ] Idea 创建与版本管理
- [ ] 协作编辑与贡献追踪
- [ ] 众筹：投资/借贷/捐赠三种类型
- [ ] 贡献权重记录（代码/设计/文档）
- [ ] 收益按累计贡献权重分配

**关联 Spec:** MARKET-003

---

### ACH-DEV-017. Performance Benchmark

**优先级:** P2
**描述:** 建立性能基准并验证达标。

**验收标准:**
- [ ] TPS >= 100
- [ ] 交易确认延迟 P99 < 3s
- [ ] 支持 1000 并发用户
- [ ] 性能测试报告发布

**关联 Spec:** CORE-001

---

### ACH-DEV-018. Observability Stack

**优先级:** P2
**描述:** 部署完整的监控和可观测性系统。

**验收标准:**
- [ ] Prometheus + Grafana 监控面板
- [ ] 关键指标采集（区块、交易、共识）
- [ ] 告警规则配置
- [ ] 日志聚合可用
- [ ] 分布式追踪

**关联 Spec:** OPERATIONS

---

## P3 - 增强功能

### ACH-DEV-019. Node Role System

**优先级:** P3
**描述:** 实现不同角色的节点类型。

**验收标准:**
- [ ] Light Node：核心模块 + GenieBot 插件
- [ ] Full Node：完整状态与历史
- [ ] Service Node：核心模块 + 服务插件
- [ ] Archive Node：完整历史索引 + 区块浏览器对接
- [ ] 节点角色配置与切换（明确热切换或重启切换）

**关联 Spec:** NODE-001

---

### ACH-DEV-020. Security Audit

**优先级:** P3
**描述:** 完成智能合约和核心模块的安全审计。

**验收标准:**
- [ ] 至少一家知名审计机构完成审计
- [ ] 所有 Critical/High 漏洞修复
- [ ] 审计报告公开
- [ ] Bug Bounty 计划启动

**关联 Spec:** ALL

---

### ACH-DEV-023. Mainnet Launch

**优先级:** P3
**描述:** 主网正式上线。

**验收标准:**
- [ ] 创世文件确定并公开
- [ ] 初始验证者集合确认
- [ ] 主网区块高度 > 1000
- [ ] 代币流通正常
- [ ] 紧急响应流程就绪

**依赖:** ACH-DEV-020 (Security Audit)

**关联 Spec:** ALL

---

### ACH-DEV-024. End-to-End Integration

**优先级:** P3
**描述:** 端到端集成测试与整体系统验收。

**验收标准:**
- [ ] 用户完整流程：注册 → 充值 → 调用服务 → 支付 → 评价
- [ ] 提供者完整流程：注册 → 托管 API Key → 接单 → 交付 → 收款
- [ ] 争议完整流程：发起 → AI 调解 → 陪审团 → MQ 变更
- [ ] Idea 众筹完整流程：创建 → 众筹 → 执行 → 收益分配
- [ ] 安全测试通过（渗透测试）

**关联 Spec:** ALL

---

## 依赖关系图

```
Wave 0 (P0): Infrastructure & Foundation
├── ACH-DEV-001 (Dev Infrastructure)
├── ACH-DEV-002 (Blockchain Network)
├── ACH-DEV-003 (Wallet & Token) ───────────► ACH-DEV-021 (Desktop App)
│                                                 (依赖钱包功能)
└── ACH-DEV-021 (Desktop App)
    └── 依赖: ACH-DEV-003, ACH-DEV-004

Wave 1 (P1): Core Modules
├── ACH-DEV-004 (Identity) ─────────────┐
├── ACH-DEV-005 (Escrow) ───────────────┼──► ACH-DEV-009 (Service Market)
├── ACH-DEV-006 (Oracle) ───────────────┘
├── ACH-DEV-007 (MQ Scoring) ───────────────► ACH-DEV-008 (Dispute)
└── ACH-DEV-010 (Testnet)

Wave 2 (P2): Plugins & Features
├── ACH-DEV-011 (LLM Plugin)
├── ACH-DEV-012 (Agent Plugin)
├── ACH-DEV-013 (Workflow Plugin)
├── ACH-DEV-014 (GenieBot UI)
├── ACH-DEV-015 (Task Market)
├── ACH-DEV-016 (Idea/Crowdfunding)
├── ACH-DEV-017 (Performance)
└── ACH-DEV-018 (Observability)

Wave 3 (P3): Production Ready
├── ACH-DEV-019 (Node Roles)
├── ACH-DEV-020 (Security Audit) ──────────► ACH-DEV-023 (Mainnet)
└── ACH-DEV-024 (E2E Integration)
```

---

## 统计

| 优先级 | 数量 | 描述 |
|--------|------|------|
| P0 | 4 | 核心基础 (含桌面应用) |
| P1 | 7 | 核心功能 (MVP) |
| P2 | 8 | 重要功能 |
| P3 | 4 | 增强功能 |
| **Total** | **23** | |

---

*Last updated: 2026-03-10*
