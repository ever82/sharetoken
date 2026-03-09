# OpenFang 知识文档

> 适用场景：ShareTokens 项目中 OpenFang 作为 Service Provider Plugin 运行 AI Agents

---

## 一、OpenFang 概述

OpenFang 是一个开源的 **AI Agent 操作系统**，由 RightNow-AI 团队使用 Rust 语言从零构建。与传统的 LLM 封装工具不同，OpenFang 将 AI Agent 从"聊天机器人"升级为"操作系统级守护进程"，支持 7x24 小时自主运行。

### 核心定位

- **Agent OS**：将 AI 能力封装为可管理、可复制的标准化流程
- **自主运行**：Hands 机制让 Agent 按计划自动执行完整工作流
- **生产级安全**：16 层纵深防御体系，企业级安全防护

### 技术指标

| 指标 | 数值 |
|------|------|
| 代码规模 | 14 个 crate，137K 行 Rust 代码 |
| 安装包 | ~32MB 单二进制文件 |
| 冷启动 | ~180ms |
| 空闲内存 | ~40MB |
| 安全机制 | 16 层独立防护 |
| 消息平台 | 40+ 频道适配器 |
| LLM 支持 | 27 家服务提供商 |

### 项目地址

- 官网：https://www.openfang.sh/
- GitHub：https://github.com/RightNow-AI/openfang

---

## 二、Agent 架构与模板

### 2.1 架构设计

OpenFang 采用分层架构：

```
┌─────────────────────────────────────┐
│          Gateway Layer              │  消息平台接入、控制界面
├─────────────────────────────────────┤
│       Agent Runtime                 │  Agent 调度与执行
├─────────────────────────────────────┤
│       Tool System                   │  38 种内置工具 + MCP 协议
├─────────────────────────────────────┤
│       Memory System                 │  SQLite + 向量嵌入持久化
├─────────────────────────────────────┤
│       Channel Adapters              │  40+ 消息平台适配
├─────────────────────────────────────┤
│       Protocol Layer                │  MCP、A2A、OFP 协议
└─────────────────────────────────────┘
```

### 2.2 Agent 模板

OpenFang 提供 30+ 预构建 Agent 模板：

```bash
# 查看可用模板
openfang agent list --templates

# 创建 Agent
openfang agent create --template researcher my-researcher
openfang agent create --template assistant my-assistant
```

常见模板类型：
- **Assistant**：通用助手
- **Researcher**：深度研究
- **Code Reviewer**：代码审查
- **Customer Support**：客户支持
- **Orchestrator**：多 Agent 编排

### 2.3 Agent 配置

Agent 配置文件 `agent.toml`：

```toml
# ~/.openfang/agents/my-agent/agent.toml

[model]
provider = "openai"
model = "gpt-4o"
api_key_env = "OPENAI_API_KEY"
max_tokens = 8192

[capabilities]
tools = ["web_search", "browser", "file_ops"]
memory = true

[guardrails]
require_approval = ["payment", "delete_file"]
```

---

## 三、WASM 沙箱与安全体系

### 3.1 WASM 双计量沙箱

所有工具代码在 WebAssembly 沙箱中隔离运行：

**双重计量机制：**

1. **燃料计量（Fuel Metering）**
   - 统计 WASM 指令执行数量
   - 默认限制：100 万指令
   - 防止 CPU 密集型攻击

2. **时间中断（Epoch Interruption）**
   - 独立看门狗线程监控
   - 默认超时：30 秒
   - 防止 I/O 阻塞攻击

**资源限制：**
- 内存上限：16MB
- 文件操作限制在工作空间内
- 子进程环境清理

### 3.2 16 层安全机制

| # | 安全层 | 作用 |
|---|--------|------|
| 1 | 基于能力的权限模型 | 最小权限原则，权限不可变 |
| 2 | WASM 双计量沙箱 | 资源限制与隔离 |
| 3 | Merkle 哈希链审计 | 操作不可篡改 |
| 4 | 信息流污点追踪 | 防止数据泄露 |
| 5 | Ed25519 清单签名 | 防供应链攻击 |
| 6 | SSRF 防护 | 阻止私有网络访问 |
| 7 | 密钥零化 | 内存中密钥自动擦除 |
| 8 | OFP 双向认证 | HMAC-SHA256 认证 |
| 9 | 安全 HTTP 头 | CSP、XSS 防护 |
| 10 | GCRA 速率限制 | 成本感知限流 |
| 11 | 子进程隔离 | 环境清理 |
| 12 | 提示注入扫描 | 输入净化 |
| 13 | 路径遍历防护 | 文件系统保护 |
| 14 | 递归调用熔断 | 防止无限循环 |
| 15 | 强制人工确认 | 敏感操作审批 |
| 16 | 行为监控 | 异常检测 |

### 3.3 能力权限模型

```rust
// openfang-types/src/capability.rs
pub enum Capability {
    FileRead,      // 文件读取
    FileWrite,     // 文件写入
    NetConnect,    // 网络连接
    NetListen,     // 网络监听
    ToolInvoke,    // 工具调用
    ToolAll,       // 所有工具
    LlmQuery,      // LLM 查询
    LlmMaxTokens,  // Token 限制
    AgentSpawn,    // 创建子 Agent
    AgentMessage,  // Agent 通信
}
```

权限验证流程：
1. `check_capability()` 在 WASM 沙箱中执行
2. `validate_capability_inheritance()` 确保子 Agent 权限小于父 Agent
3. 敏感操作记录到 Merkle 审计链

---

## 四、Hands 自主能力包

Hands 是 OpenFang 的核心创新，每个 Hand 是一个 **预打包的自主能力包**，包含：

- `HAND.toml`：工具声明、配置、仪表盘指标
- 系统提示词：500+ 字多阶段专家级操作手册
- `SKILL.md`：领域知识注入
- Guardrails：敏感操作审批规则

### 4.1 内置 Hands 列表

| Hand | 功能 | 运行模式 |
|------|------|----------|
| **Collector** | 情报收集与监控 | 持续监控 |
| **Clip** | 视频内容剪辑与发布 | 按计划 |
| **Lead** | 潜在客户生成与管理 | 每日执行 |
| **Content** | 内容创作与运营 | 按计划 |
| **Trade** | 交易执行（强制审批） | 按触发 |
| **Browser** | 网页自动化操作 | 按需 |
| **Twitter** | 社交媒体管理 | 持续运行 |

### 4.2 Collector - 情报收集

**功能：**
- OSINT 风格持续监控目标
- 变更检测与情感分析
- 自动构建知识图谱
- 关键事件即时告警

**使用场景：**
- 竞品动态监控
- 舆情预警
- 行业趋势追踪

**配置示例：**

```toml
# HAND.toml - Collector
[hand]
name = "collector"
schedule = "*/15 * * * *"  # 每 15 分钟

[targets]
companies = ["competitor-a", "competitor-b"]
keywords = ["AI Agent", "LLM", "automation"]
sources = ["news", "twitter", "github"]

[output]
format = "markdown"
channel = "telegram"
```

### 4.3 Clip - 视频剪辑

**8 阶段自动化流水线：**

1. 下载视频
2. 语音转文字（5 种 STT 后端）
3. 高光片段识别
4. 竖屏剪辑适配
5. 字幕嵌入
6. 封面生成
7. AI 配音（可选）
8. 跨平台发布

**配置示例：**

```toml
[hand]
name = "clip"

[pipeline]
stt_backend = "whisper"
output_format = "vertical"  # 9:16
add_captions = true
add_thumbnail = true

[platforms]
telegram = true
whatsapp = true
```

### 4.4 Lead - 线索生成

**功能：**
- 根据 ICP 画像自动发现潜在客户
- 网络调研与数据富化
- 0-100 多维评分
- 自动去重
- 输出 CSV/JSON/Markdown

**每日执行流程：**
```
搜索 → 调研 → 评分 → 去重 → 输出报告
```

### 4.5 Browser - 浏览器自动化

**基于 Playwright 驱动：**
- 多步表单填写
- 页面导航与点击
- 会话持久化
- **强制购买审批门**：涉及金额自动暂停

**安全机制：**
```toml
[guardrails]
require_approval = [
    "payment",
    "checkout",
    "transfer",
    "subscribe"
]
```

### 4.6 Twitter - 社媒管理

**功能：**
- 7 种内容格式轮换
- 最优时间自动发帖
- 回复提及互动
- 绩效追踪
- **审批队列**：所有内容人工确认后发布

### 4.7 Hands CLI 操作

```bash
# 激活 Hand
openfang hand activate collector --config ./collector.toml

# 查看状态
openfang hand status collector

# 暂停（保留状态）
openfang hand pause lead

# 恢复
openfang hand resume lead

# 列出所有 Hands
openfang hand list

# 查看日志
openfang hand logs researcher --tail 50
```

### 4.8 自定义 Hand

创建目录结构：

```
my-custom-hand/
├── HAND.toml      # 配置文件
├── SKILL.md       # 领域知识
└── prompts/
    └── system.md  # 系统提示词
```

HAND.toml 示例：

```toml
[hand]
name = "blockchain-monitor"
version = "1.0.0"

[tools]
allowed = ["web_search", "http_request", "file_write"]

[schedule]
cron = "*/5 * * * *"

[metrics]
dashboard = ["transactions", "alerts", "gas_price"]

[guardrails]
require_approval = []
```

---

## 五、SDK 与 CLI 使用

### 5.1 安装

**macOS / Linux：**
```bash
curl -fsSL https://openfang.sh/install | sh
```

**Windows：**
```powershell
irm https://openfang.sh/install.ps1 | iex
```

**Cargo：**
```bash
cargo install --git https://github.com/RightNow-AI/openfang openfang-cli
```

**Docker：**
```bash
docker pull ghcr.io/rightnow-ai/openfang:latest
```

### 5.2 初始化与启动

```bash
# 初始化配置
openfang init

# 快速初始化（非交互）
openfang init --quick

# 启动守护进程
openfang start

# 启动 Dashboard
openfang dashboard --port 4200
# 访问 http://localhost:4200
```

### 5.3 LLM 配置

编辑 `~/.openfang/config.toml`：

```toml
# OpenAI
[default_model]
provider = "openai"
model = "gpt-4o"
api_key_env = "OPENAI_API_KEY"

# 通义千问（国内推荐）
[default_model]
provider = "openai"
model = "qwen-plus"
api_key_env = "DASHSCOPE_API_KEY"
base_url = "https://dashscope.aliyuncs.com/compatible-mode/v1"

# Anthropic
[providers.anthropic]
api_key_env = "ANTHROPIC_API_KEY"
```

环境变量文件 `~/.openfang/.env`：
```bash
OPENAI_API_KEY=sk-xxx
DASHSCOPE_API_KEY=sk-xxx
ANTHROPIC_API_KEY=sk-ant-xxx
```

### 5.4 CLI 常用命令

```bash
# 诊断
openfang doctor

# 版本
openfang --version

# Agent 管理
openfang agent list
openfang agent create --template researcher my-agent
openfang agent start my-agent
openfang agent stop my-agent

# Hand 管理
openfang hand list
openfang hand activate collector
openfang hand status

# 频道配置
openfang channel list
openfang channel enable telegram

# 从 OpenClaw 迁移
openfang migrate --from openclaw
```

### 5.5 SDK 集成（Rust）

```rust
use openfang::{Agent, Hand, Config};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 加载配置
    let config = Config::from_file("~/.openfang/config.toml")?;

    // 创建 Agent
    let agent = Agent::new("my-agent")
        .with_model("gpt-4o")
        .with_tools(vec!["web_search", "browser"])
        .build()?;

    // 激活 Hand
    let hand = Hand::activate("collector", "./collector.toml").await?;

    // 获取状态
    let status = hand.status().await?;
    println!("Hand status: {:?}", status);

    Ok(())
}
```

### 5.6 REST API

```bash
# 启动 API 服务
openfang api --port 8080

# 调用示例
curl http://localhost:8080/v1/agents
curl http://localhost:8080/v1/hands
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "agent": "assistant"}'
```

---

## 六、区块链服务集成

### 6.1 适用场景

在 ShareTokens 项目中，OpenFang 可作为 **Service Provider Plugin** 实现：

- 链上数据监控与采集
- 智能合约交互自动化
- 交易审批与执行
- 多链状态追踪

### 6.2 自定义区块链 Hand

创建 `blockchain-monitor/HAND.toml`：

```toml
[hand]
name = "blockchain-monitor"
version = "1.0.0"
description = "Multi-chain monitor for ShareTokens"

[tools]
allowed = [
    "http_request",
    "web_socket",
    "file_write",
    "json_parse"
]

[schedule]
cron = "*/1 * * * *"  # 每分钟检查

[metrics]
dashboard = [
    "block_height",
    "pending_transactions",
    "contract_events",
    "gas_price"
]

[guardrails]
# 交易操作必须人工确认
require_approval = [
    "sign_transaction",
    "send_transaction",
    "contract_write"
]

[notifications]
channels = ["telegram", "discord"]
alert_on = ["large_transfer", "contract_event", "price_change"]
```

领域知识 `SKILL.md`：

```markdown
# Blockchain Monitoring Domain Knowledge

## Supported Chains
- Ethereum
- Polygon
- BSC
- Arbitrum
- Optimism

## Key Concepts
- RPC Endpoints
- Event Logs
- Transaction Receipts
- Gas Estimation

## Security Considerations
- Never expose private keys
- Validate all addresses
- Rate limit RPC calls
```

### 6.3 集成配置

**环境变量：**
```bash
# RPC 端点
ETHEREUM_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/xxx
POLYGON_RPC_URL=https://polygon-mainnet.g.alchemy.com/v2/xxx

# 钱包（仅读取）
WATCH_ADDRESS=0x...

# API Keys
ETHERSCAN_API_KEY=xxx
COINMARKETCAP_API_KEY=xxx
```

**Agent 配置：**

```toml
# ~/.openfang/agents/blockchain-watcher/agent.toml

[model]
provider = "openai"
model = "gpt-4o"

[capabilities]
tools = ["http_request", "web_socket", "json_parse"]
memory = true

[custom]
rpc_endpoints = ["ethereum", "polygon"]
watch_addresses = ["0x..."]
alert_threshold_eth = 10.0
```

### 6.4 与 ShareTokens 集成架构

```
┌─────────────────────────────────────────────────────────┐
│                    ShareTokens Core                      │
│  (Service Provider Registry, Task Management, Dispute)  │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                  OpenFang Plugin Layer                   │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  Collector  │  │   Browser   │  │   Custom    │     │
│  │    Hand     │  │    Hand     │  │ Blockchain  │     │
│  └─────────────┘  └─────────────┘  │    Hand     │     │
│         │                │         └─────────────┘     │
│         ▼                ▼                  │           │
│  ┌─────────────────────────────────────────┐│           │
│  │            Tool Layer                    ││           │
│  │  HTTP │ WebSocket │ File │ Browser      ││           │
│  └─────────────────────────────────────────┘│           │
│                                             ▼           │
│  ┌─────────────────────────────────────────────────────┐│
│  │              WASM Sandbox (16 Layers Security)      ││
│  └─────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
              ┌───────────────────────────┐
              │   Blockchain Networks     │
              │  Ethereum │ Polygon │ ... │
              └───────────────────────────┘
```

### 6.5 审计与合规

OpenFang 的 Merkle 哈希链审计天然适合区块链场景：

```bash
# 导出审计日志
openfang audit export --format json --output ./audit-log.json

# 验证完整性
openfang audit verify

# 查看特定 Agent 操作记录
openfang audit show --agent blockchain-watcher
```

---

## 七、最佳实践

### 7.1 安全建议

1. **最小权限原则**：为每个 Agent 分配最小必要权限
2. **敏感操作审批**：始终为资金相关操作配置 `require_approval`
3. **密钥管理**：使用环境变量，启用密钥零化
4. **定期审计**：定期验证 Merkle 审计链完整性

### 7.2 性能优化

1. **模型选择**：简单任务用 qwen-turbo，复杂推理用 qwen-max
2. **资源限制**：合理设置 WASM 燃料限制
3. **调度优化**：避免多个 Hands 同时高负载运行

### 7.3 国内部署

推荐使用阿里云 DashScope：

```toml
[default_model]
provider = "openai"
model = "qwen-plus"
api_key_env = "DASHSCOPE_API_KEY"
base_url = "https://dashscope.aliyuncs.com/compatible-mode/v1"
```

支持的国内渠道：
- 飞书（Feishu/Lark）
- 钉钉（通过 Webhook）
- 企业微信（通过 Webhook）

---

## 八、故障排查

```bash
# 诊断命令
openfang doctor

# 查看日志
openfang logs --tail 100

# 检查 LLM 连接
openfang test llm

# 检查频道连接
openfang test channel telegram

# 重置 Agent
openfang agent reset my-agent
```

常见问题：

| 问题 | 解决方案 |
|------|----------|
| LLM 认证失败 | 检查 API Key 环境变量 |
| WASM 超时 | 增加 fuel 限制或检查循环 |
| 频道连接失败 | 验证 token 和 chat_id |
| 审计链损坏 | 检查磁盘空间，恢复备份 |

---

## 参考资料

- 官方文档：https://docs.openfang.sh/
- GitHub 仓库：https://github.com/RightNow-AI/openfang
- FangHub 市场：https://hub.openfang.sh/
- 社区 Discord：https://discord.gg/openfang

---

*文档版本：2026-03-05*
*适用于 ShareTokens Service Provider Plugin 集成*
