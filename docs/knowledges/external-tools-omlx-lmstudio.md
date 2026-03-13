# 外部工具研究：OMLX & LM Studio

> 研究日期: 2026-03-12
> 状态: 已评估，暂未集成，供未来参考

---

## 1. OMLX

### 基本信息

| 属性 | 内容 |
|------|------|
| **项目** | omlx |
| **GitHub** | https://github.com/jundot/omlx |
| **License** | Apache-2.0 |
| **Stars** | 3,001 |
| **Forks** | 226 |
| **主要语言** | Python |
| **定位** | Apple Silicon 优化的本地 LLM 推理服务器 |
| **更新时间** | 2026-03-11 (活跃开发中) |

### 核心特性

1. **OpenAI 兼容 API** - `http://localhost:8000/v1`
2. **连续批处理 (Continuous Batching)** - 并发请求高效处理
3. **分层 KV 缓存 (Tiered KV Cache)**
   - Hot tier (RAM): 频繁访问块驻留内存
   - Cold tier (SSD): 缓存持久化，safetensors 格式
4. **多模型同时服务** - LLM/VLM/Embedding/Reranker
5. **MCP 支持** - Model Context Protocol
6. **Web 管理后台** - 实时监控、模型管理、聊天界面
7. **Claude Code 优化** - 上下文缩放、SSE keep-alive

### 安装方式

```bash
# macOS App
# 从 Releases 下载 .dmg，拖拽安装

# Homebrew (推荐)
brew tap jundot/omlx https://github.com/jundot/omlx
brew install omlx
brew services start omlx  # 后台服务

# 源码安装
git clone https://github.com/jundot/omlx.git
cd omlx
pip install -e .
# 或带 MCP 支持: pip install -e ".[mcp]"
```

### 系统要求

- macOS 15.0+ (Sequoia)
- Python 3.10+
- **Apple Silicon (M1/M2/M3/M4)** - 仅支持 Apple Silicon

### MCP 配置示例

```json
{
  "servers": {
    "filesystem": {
      "transport": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "enabled": true,
      "timeout": 30
    },
    "sqlite": {
      "transport": "stdio",
      "command": "uvx",
      "args": ["mcp-server-sqlite", "--db-path", "data.db"],
      "enabled": false,
      "timeout": 30
    },
    "web-search": {
      "transport": "sse",
      "url": "http://localhost:3001/sse",
      "enabled": false,
      "timeout": 60
    }
  },
  "max_tool_calls": 10,
  "default_timeout": 30.0
}
```

---

## 2. LM Studio

### 基本信息

| 属性 | 内容 |
|------|------|
| **官网** | https://lmstudio.ai |
| **定位** | 跨平台本地 LLM 管理桌面应用 |
| **支持平台** | Windows / macOS / Linux |
| **API** | OpenAI 兼容本地服务器 |
| **特点** | 图形化管理、一键下载模型、跨平台 |

### 核心特性

- 本地运行开源 LLM (Llama, Mistral, Qwen 等)
- OpenAI 兼容的 HTTP API
- 图形化模型管理界面
- 内置聊天界面
- 隐私优先 (数据不上云)

---

## 3. 对 ShareToken 项目的潜在价值

### 3.1 帮助 ACH-DEV-011 (LLM API Key Custody Plugin)

**问题场景**:
- 开发需要调用 LLM，但 OpenAI API 需要付费
- 测试 API Key 托管功能时，真实 Key 有泄露风险
- CI/CD 自动化测试无法依赖外部 API

**解决方案**:

```
┌─────────────┐      ┌──────────────────┐      ┌──────────┐
│ OMLX /      │      │ ShareToken       │      │ GenieBot │
│ LM Studio   │─────▶│ LLM Key Plugin   │─────▶│  UI      │
│             │      │                  │      │          │
│ • 完全本地  │      │ • 本地模式配置    │      │ • 测试    │
│ • 零成本    │      │ • API 转发       │      │ • 演示    │
│ • 隐私安全  │      │ • 计费统计       │      │ • 开发    │
└─────────────┘      └──────────────────┘      └──────────┘
```

**代码示例**:
```go
// x/agentgateway/llm/client.go
if config.LocalMode {
    baseURL = "http://localhost:8000/v1"  // omlx/LM Studio
    apiKey = "local"  // 无需真实 API Key
} else {
    baseURL = provider.BaseURL
    apiKey = decryptFromChain(config.EncryptedAPIKey)
}
```

**推荐用法**:

| 场景 | 工具选择 | 原因 |
|------|----------|------|
| Mac 开发者 | **OMLX** | 性能更好 (Apple Silicon 优化)、MCP 支持、Apache 协议可二次开发 |
| 跨平台团队 | **LM Studio** | 支持 Windows/Linux，图形化易用 |
| CI/CD 测试 | **OMLX** | 可命令行安装、可脚本化控制 |

### 3.2 帮助 ACH-DEV-014 (GenieBot UI)

OMLX 的分层 KV 缓存技术对 GenieBot 的长对话场景有参考价值:

- **Hot tier (RAM)**: 活跃对话缓存
- **Cold tier (SSD)**: 历史对话持久化

可借鉴用于:
- 本地对话历史缓存
- 多会话上下文管理
- Token 使用量优化

### 3.3 帮助 ACH-DEV-012/013 (Agent/Workflow)

**MCP (Model Context Protocol)** 支持:

MCP 是标准化的 LLM 工具调用接口，让 LLM 可以调用外部工具 (文件、数据库、搜索等)。

对 ShareToken Agent 系统的价值:
- 可复用 MCP 生态的工具服务器
- Agent Executor 可对接 MCP 协议
- 避免重复开发工具调用基础设施

---

## 4. 工具对比

| 维度 | OMLX | LM Studio |
|------|------|-----------|
| **平台支持** | 仅 macOS + Apple Silicon | Windows/macOS/Linux |
| **协议** | Apache-2.0 (可二次开发) | 商业软件 (免费使用) |
| **性能优化** | Apple Silicon 深度优化 | 通用优化 |
| **MCP 支持** | ✅ 原生支持 | ❌ 暂不支持 |
| **管理方式** | CLI + Web UI | 桌面 GUI |
| **KV 缓存** | Hot+Cold 分层缓存 | 标准缓存 |
| **适用场景** | 开发/测试/生产部署 | 个人用户/快速体验 |

---

## 5. 集成建议

### 何时考虑集成

- [ ] ACH-DEV-011 开发测试阶段 (本地 LLM 替代 OpenAI)
- [ ] 需要 MCP 协议支持 Agent 工具调用
- [ ] 桌面应用用户要求本地运行 LLM 降低成本

### 快速集成步骤 (供参考)

```bash
# 1. 安装 OMLX
brew tap jundot/omlx
brew install omlx
brew services start omlx

# 2. 下载测试模型 (通过 Web UI: http://localhost:8000/admin)
# 推荐: Qwen2.5-7B-Instruct (中文支持好)

# 3. 测试 API
curl http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen2.5-7b-instruct",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

---

## 6. 相关资源

- **OMLX GitHub**: https://github.com/jundot/omlx
- **OMLX 官网**: https://omlx.ai
- **LM Studio**: https://lmstudio.ai
- **MCP 协议**: https://modelcontextprotocol.io

---

*文档创建: 2026-03-12*
*状态: 待集成*
