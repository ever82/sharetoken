# 用户成就测试覆盖对照表

> 本文件记录 `docs/achievements/for-user.md` 中所有验收条目的测试覆盖情况
> 更新时间: 2026-03-11

---

## 图例

- ✅ **已覆盖** - 已有自动化测试或人工验收清单
- 🔄 **需补充** - 需要创建自动化测试
- 👤 **人工验收** - 需要人工验收（UI/体验/第三方集成）
- ⏳ **待开发** - 功能尚未实现

---

## P0 - 核心路径

### ACH-USER-001: Secure Digital Wallet

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 创建/导入钱包（从打开网页到完成 < 60s） | 👤 人工 | ✅ | `manual/checklist-P0-wallet.md` |
| 查看 STT 余额和交易历史 | 🔄 自动 | ✅ | `e2e/wallet_test.go` |
| 安全地转账和收款 | 🔄 自动 | ✅ | `e2e/wallet_test.go` |
| 使用 Keplr 或手机钱包管理资产 | 👤 人工 | ✅ | `manual/checklist-P0-wallet.md` |
| 导出私钥 | 👤 人工 | ✅ | `manual/checklist-P0-wallet.md` |

**缺失测试:** 无

---

### ACH-USER-002: One-Click AI Access

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 通过 GenieBot 开始对话，无需安装任何软件 | 👤 人工 | ✅ | `manual/checklist-P0-geniebot.md` |
| 从打开网页到发送第一条消息 < 60s | 👤 人工 | ✅ | `manual/checklist-P0-geniebot.md` |
| 自然语言描述需求，系统自动理解并推荐服务 | 👤 人工 | ✅ | `manual/checklist-P0-geniebot.md` |
| 意图识别准确率 >= 85% | 👤 人工 | ✅ | `manual/checklist-P0-geniebot.md` |
| 响应时间 < 3s | 👤 人工 | ✅ | `manual/checklist-P0-geniebot.md` |
| 一键调用至少 5 种主流 AI 模型 | 🔄 自动 | ✅ | `e2e/geniebot_test.go` |
| 查看每次调用的费用明细 | 🔄 自动 | ✅ | `e2e/geniebot_test.go` |

**缺失测试:** 无

---

### ACH-USER-003: Fund Security Guarantee

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 在服务开始前付款，资金进入托管 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 在确认满意后资金才释放给服务者 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 如果有问题，可以申请冻结资金 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 争议解决后，资金按裁决分配 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 可以随时查看托管资金状态 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |

**缺失测试:** 无

---

### ACH-USER-004: Transparent Service Pricing

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 在服务市场浏览所有可用服务 | 🔄 自动 | ✅ | `e2e/marketplace_pricing_test.go` |
| 查看每个服务的定价模式 | 🔄 自动 | ✅ | `e2e/marketplace_pricing_test.go` |
| 比较不同提供者的价格和评分 | 👤 人工 | ⚠️ | 需补充人工验收清单 |
| 在消费前看到预估费用 | 🔄 自动 | ✅ | `e2e/marketplace_pricing_test.go` |
| 查看详细的消费账单（支持导出 CSV） | 👤 人工 | ⚠️ | 需补充人工验收清单 |

**缺失测试:**
- `manual/checklist-P0-pricing.md` - UI 展示、导出功能

---

### ACH-USER-005: First-Time Onboarding

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 用微信/GitHub/Google 一键登录 | 🔄 自动 | ✅ | `e2e/onboarding_test.go` |
| 钱包自动创建（无需记助记词，可后续导出） | 👤 人工 | ✅ | `manual/checklist-P0-wallet.md` |
| 获得初始测试代币（水龙头） | 🔄 自动 | ✅ | `e2e/onboarding_test.go` |
| 首次使用有交互式引导教程（3-5 分钟） | 👤 人工 | ⏳ | 待开发 |
| "What can I do here" 常见用例列表 | 👤 人工 | ⏳ | 待开发 |

**缺失测试:**
- 引导教程和用例列表功能待开发

---

### ACH-USER-021: Desktop App (开箱即用)

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 从 Releases 下载对应系统的安装包 | 👤 人工 | ✅ | `manual/checklist-P0-desktop.md` |
| 不需要安装 Node.js、npm、Go 等任何依赖 | 👤 人工 | ✅ | `manual/checklist-P0-desktop.md` |
| 不需要运行任何命令行命令 | 👤 人工 | ✅ | `manual/checklist-P0-desktop.md` |
| 双击打开后，图形界面自动启动 | 👤 人工 | ✅ | `manual/checklist-P0-desktop.md` |
| 应用自动检测或选择本地节点/远程节点 | 🔄 自动 | ✅ | `e2e/desktop_app_test.go` |
| 图形界面中完成核心功能 | 👤 人工 | ✅ | `manual/checklist-P0-desktop.md` |

**缺失测试:** 无

---

## P1 - 完整体验

### ACH-USER-006: Task Progress Tracking

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 查看所有进行中的任务列表 | 🔄 自动 | ✅ | `e2e/task_tracking_test.go` |
| 支持分页，每页 >= 10 条 | 🔄 自动 | ⚠️ | 需补充测试 |
| 支持按时间/优先级/状态排序 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看每个任务的当前状态和进度百分比 | 🔄 自动 | ✅ | `e2e/task_tracking_test.go` |
| 状态变更后 10 秒内更新 | 👤 人工 | ⚠️ | 需补充人工验收 |
| 百分比显示，精确到 1% | 👤 人工 | ⚠️ | 需补充人工验收 |
| 查看里程碑完成情况 | 🔄 自动 | ⚠️ | 需补充测试 |
| 收到关键节点的通知 | 👤 人工 | ⚠️ | 需补充人工验收 |
| 与服务者沟通（留言/补充需求） | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/task_tracking_advanced_test.go` - 分页、排序、里程碑、留言
- `manual/checklist-P1-task.md` - 实时更新、通知

---

### ACH-USER-007: Verified Identity Benefits

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 通过微信/GitHub/Google 完成实名认证 | 🔄 自动 | ✅ | `e2e/onboarding_test.go` |
| 认证状态显示在资料中（验证通过标记） | 👤 人工 | ⚠️ | 需补充人工验收 |
| 获得更高的交易限额 | 🔄 自动 | ⚠️ | 需补充测试 |
| 获得更高的提现限额 | 🔄 自动 | ⚠️ | 需补充测试 |
| 被选为陪审员参与争议裁决 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/identity_limits_test.go` - 限额系统
- `e2e/juror_selection_test.go` - 陪审员资格
- `manual/checklist-P1-identity.md` - UI 展示

---

### ACH-USER-008: Fair Dispute Resolution

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 对任何交易发起争议 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 通过 AI 调解员陈述观点 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| 提交证据（截图、聊天记录等） | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 选择陪审团裁决 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 看到裁决结果和理由 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 查看 MQ 分数变化 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |

**缺失测试:** 无

---

### ACH-USER-009: Service Failure & Refund

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 在任务详情页发起退款申请 | 🔄 自动 | ✅ | `e2e/refund_flow_test.go` |
| 选择退款原因 | 🔄 自动 | ✅ | `e2e/refund_flow_test.go` |
| 查看退款处理进度 | 🔄 自动 | ✅ | `e2e/refund_flow_test.go` |
| 退款审核在 48 小时内完成 | 🔄 自动 | ⚠️ | 需补充时间验证 |
| 退款成功后 STT 在 10 分钟内返还到账 | 🔄 自动 | ⚠️ | 需补充时间验证 |

**缺失测试:**
- `e2e/refund_flow_test.go` 补充时间验证

---

### ACH-USER-010: Reputation Dashboard

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 查看 MQ 分数和历史变化 | 🔄 自动 | ⚠️ | 需补充测试 |
| 历史变化图表 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| 查看信誉等级和对应权益 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看参与争议的记录和结果 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看被选为陪审员的记录 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看贡献统计 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/reputation_dashboard_test.go` - MQ 分数、等级、历史记录

---

### ACH-USER-011: Trustworthy Community

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 实名认证（一账户一真人） | 🔄 自动 | ✅ | `e2e/onboarding_test.go` |
| 查看其他用户的信誉评分 | 🔄 自动 | ⚠️ | 需补充测试 |
| 不良行为者会被惩罚或封禁 | 🔄 自动 | ⚠️ | 需补充测试 |
| 争议由社区陪审团公平裁决 | 🔄 自动 | ✅ | `e2e/escrow_security_test.go` |
| 系统去中心化运行 | 🔄 自动 | ⏳ | 需设计测试方案 |

**缺失测试:**
- `e2e/community_trust_test.go` - 信誉查看、惩罚机制

---

## P2 - 增长与差异化

### ACH-USER-012: Service Discovery & Comparison

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 按类别浏览服务（LLM/Agent/Workflow） | 🔄 自动 | ⚠️ | 需补充测试 |
| 按关键词搜索服务 | 🔄 自动 | ⚠️ | 需补充测试 |
| 按价格/评分/响应时间排序 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看服务详情和案例 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看其他用户的评价 | 🔄 自动 | ⚠️ | 需补充测试 |
| 收藏常用服务 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/service_discovery_test.go` - 搜索、排序、收藏

---

### ACH-USER-013: Idea Submission & Refinement

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 提交粗糙的想法（自然语言描述） | 🔄 自动 | ⏳ | 待开发 |
| AI 分析和细化想法 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| 评估所需服务和成本 | 🔄 自动 | ⏳ | 待开发 |
| 查看类似想法的参考案例 | 🔄 自动 | ⏳ | 待开发 |
| 选择开始执行或保存想法 | 🔄 自动 | ⏳ | 待开发 |

**缺失测试:** 功能待开发

---

### ACH-USER-014: Service Provider Earnings

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 注册成为服务提供者 | 🔄 自动 | ⚠️ | 需补充测试 |
| 托管 LLM API Key 并设置定价 | 🔄 自动 | ✅ | 已有 llmcustody 模块测试 |
| 注册 Agent 技能 | 🔄 自动 | ⚠️ | 需补充测试 |
| 承接 Workflow 任务 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看收入明细和提现 | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看客户评价和评分 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/provider_workflow_test.go` - 完整的提供者流程

---

### ACH-USER-015: Crowdfunding for Ideas

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 发布 Idea 并设置众筹目标 | 🔄 自动 | ⚠️ | 需补充测试 |
| 选择众筹类型（投资/借贷/捐赠） | 🔄 自动 | ⚠️ | 需补充测试 |
| 更新 Idea 进展 | 👤 人工 | ⚠️ | 需补充人工验收 |
| 查看众筹进度 | 🔄 自动 | ⚠️ | 需补充测试 |
| 项目成功后收益分配 | 🔄 自动 | ⚠️ | 需补充测试 |
| 支持者查看回报分配 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/crowdfunding_flow_test.go` - 众筹完整流程

---

### ACH-USER-016: Developer API Access

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 在控制台生成 API Key | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看完整的 API 文档和 SDK 示例 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| 设置 API 调用限额和白名单 IP | 🔄 自动 | ⚠️ | 需补充测试 |
| 查看 API 调用日志和错误统计 | 🔄 自动 | ⚠️ | 需补充测试 |
| 支持 Python, TypeScript SDK | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |

**缺失测试:**
- `e2e/developer_api_test.go` - API Key、限额、日志

---

## P3 - 规模化

### ACH-USER-017: Multi-Agent Task Execution

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 描述复杂需求 | 🔄 自动 | ⏳ | 待开发 |
| 系统自动分解为子任务 | 🔄 自动 | ⏳ | 待开发 |
| 看到多个 Agent 协作执行 | 👤 人工 | ⏳ | 待开发 |
| 查看子任务的负责 Agent 和进度 | 🔄 自动 | ⏳ | 待开发 |
| 验收阶段性成果 | 🔄 自动 | ⏳ | 待开发 |
| 最终交付完整成果 | 👤 人工 | ⏳ | 待开发 |

**缺失测试:** 功能待开发

---

### ACH-USER-018: Cross-Platform Access

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 在桌面浏览器使用完整功能 | 🔄 自动 | ⚠️ | 需补充浏览器兼容性测试 |
| 在手机浏览器使用核心功能 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| 通过 WalletConnect 连接移动钱包 | 🔄 自动 | ✅ | `e2e/wallet_test.go` |
| 账户和数据在所有设备同步 | 🔄 自动 | ⚠️ | 需补充测试 |
| 安全地在任何设备登出 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/cross_platform_test.go` - 多设备同步、登出

---

### ACH-USER-019: Data Export & Account Deletion

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 一键导出所有个人数据（JSON/CSV） | 🔄 自动 | ⚠️ | 需补充测试 |
| 导出包含交易记录、任务历史、评价、设置 | 🔄 自动 | ⚠️ | 需补充测试 |
| 申请永久删除账户 | 🔄 自动 | ⚠️ | 需补充测试 |
| 30 天冷静期 | 🔄 自动 | ⚠️ | 需补充测试 |
| 删除后链上数据只保留哈希 | 🔄 自动 | ⚠️ | 需补充测试 |

**缺失测试:**
- `e2e/data_export_test.go` - 数据导出、账户删除

---

### ACH-USER-020: Multilingual Support

| 验收标准 | 类型 | 状态 | 测试位置 |
|---------|------|------|---------|
| 平台支持至少 3 种语言 | 🔄 自动 | ⚠️ | 需补充测试 |
| 在设置中切换界面语言 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| AI 对话支持多语言输入 | 👤 人工 | ✅ | `manual/checklist-P1-P3.md` |
| 服务描述自动翻译 | 👤 人工 | ⚠️ | 需补充人工验收 |
| 争议调解支持选择语言 | 👤 人工 | ⚠️ | 需补充人工验收 |

**缺失测试:**
- `e2e/i18n_test.go` - 多语言功能

---

## 统计汇总

### 总体覆盖情况

| 优先级 | 总数 | 已覆盖 | 需补充 | 人工验收 | 待开发 |
|--------|------|--------|--------|----------|--------|
| P0 | 6 | 4 | 1 | 4 | 1 |
| P1 | 6 | 2 | 5 | 2 | 0 |
| P2 | 5 | 0 | 4 | 1 | 2 |
| P3 | 4 | 0 | 3 | 2 | 1 |
| **Total** | **21** | **6** | **13** | **9** | **4** |

### 缺失的自动化测试列表

1. ✅ `e2e/wallet_test.go` - 已存在
2. ✅ `e2e/geniebot_test.go` - 已存在
3. ✅ `e2e/escrow_security_test.go` - 已存在
4. ✅ `e2e/marketplace_pricing_test.go` - 已存在
5. ✅ `e2e/onboarding_test.go` - 已存在
6. ✅ `e2e/desktop_app_test.go` - 已存在
7. ✅ `e2e/refund_flow_test.go` - 已存在
8. ✅ `e2e/task_tracking_test.go` - 已存在
9. 🔄 `e2e/task_tracking_advanced_test.go` - 分页、排序、里程碑
10. 🔄 `e2e/identity_limits_test.go` - 限额系统
11. 🔄 `e2e/juror_selection_test.go` - 陪审员资格
12. 🔄 `e2e/reputation_dashboard_test.go` - MQ 分数、等级
13. 🔄 `e2e/community_trust_test.go` - 信誉查看、惩罚机制
14. 🔄 `e2e/service_discovery_test.go` - 搜索、排序、收藏
15. 🔄 `e2e/provider_workflow_test.go` - 提供者流程
16. 🔄 `e2e/crowdfunding_flow_test.go` - 众筹流程
17. 🔄 `e2e/developer_api_test.go` - API Key、限额
18. 🔄 `e2e/cross_platform_test.go` - 多设备同步
19. 🔄 `e2e/data_export_test.go` - 数据导出
20. 🔄 `e2e/i18n_test.go` - 多语言

### 缺失的人工验收清单

1. ✅ `manual/checklist-P0-wallet.md` - 已存在
2. ✅ `manual/checklist-P0-geniebot.md` - 已存在
3. ✅ `manual/checklist-P0-desktop.md` - 已存在
4. 🔄 `manual/checklist-P0-pricing.md` - 定价展示、导出
5. 🔄 `manual/checklist-P1-task.md` - 任务进度 UI
6. 🔄 `manual/checklist-P1-identity.md` - 认证标识 UI
7. ✅ `manual/checklist-P1-P3.md` - 已存在（综合）

---

## 下一步行动建议

### 高优先级（P0 剩余）

1. 创建 `manual/checklist-P0-pricing.md` - 服务定价 UI 验收

### 中优先级（P1 功能）

1. 创建 `e2e/reputation_dashboard_test.go` - 信誉系统
2. 创建 `e2e/identity_limits_test.go` - 身份限额
3. 创建 `manual/checklist-P1-identity.md` - 认证 UI 验收

### 低优先级（P2/P3 功能）

1. 按开发进度逐步补充相关测试
