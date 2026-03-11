# ACH-DEV-024: Agent Gateway - External Agent Integration

**优先级:** P2
**类型:** Infrastructure / Protocol
**状态:** 已完成 ✅
**完成日期:** 2026-03-11
**实际时间:** 2小时
**依赖:** ACH-DEV-014 (GenieBot UI), ACH-DEV-012 (Agent Executor)

---

## 目标

实现 Agent Gateway，让 Claude Code、OpenFang/OpenClaw 等外部 Agent 工具可以通过标准协议（MCP + A2A）对接 GenieBot。

---

## 快速实现清单 (2-4小时)

### Phase 1: MCP 核心 (30分钟)

**使用 mcp-go 库，几行代码搞定**

```go
// cmd/agent-gateway/main.go
package main

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    // 创建 MCP Server
    s := server.NewMCPServer(
        "sharetoken-gateway",
        "1.0.0",
    )

    // 注册 Tools
    s.AddTool(mcp.NewTool("query_balance",
        mcp.WithDescription("查询账户STT余额"),
        mcp.WithString("address", mcp.Required(), mcp.Description("账户地址")),
    ), handleQueryBalance)

    s.AddTool(mcp.NewTool("chat_with_genie",
        mcp.WithDescription("与GenieBot对话"),
        mcp.WithString("message", mcp.Required()),
    ), handleChat)

    // 启动 stdio server
    if err := server.ServeStdio(s); err != nil {
        panic(err)
    }
}
```

- [x] 引入 `github.com/mark3labs/mcp-go` (实际使用标准库自定义实现，避免外部依赖)
- [x] 创建 `cmd/agent-gateway/main.go`
- [x] 实现 4个核心 Tool (query_balance, chat_with_genie, create_task, create_escrow)
- [x] stdio transport 启动

### Phase 2: A2A 核心 (30分钟)

**A2A 是简单 HTTP API**

```go
// agentgateway/a2a/handler.go
package a2a

import (
    "encoding/json"
    "net/http"
)

func Handler() http.Handler {
    mux := http.NewServeMux()

    // Agent Card
    mux.HandleFunc("/.well-known/agent.json", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(AgentCard{
            Name: "ShareToken Agent",
            Version: "1.0.0",
            Capabilities: []string{"task_execution", "query"},
        })
    })

    // Tasks
    mux.HandleFunc("/a2a/tasks", handleTasks)

    return mux
}
```

- [x] 创建 `x/agentgateway/a2a/handler.go`
- [x] 实现 Agent Card endpoint
- [x] 实现 Task 基础接口 (list, create, negotiate)

### Phase 3: 链交互 (30分钟)

**复用现有 ChainHelper**

```go
// agentgateway/keeper/chain.go
package keeper

type ChainKeeper struct {
    queryClient types.QueryClient
}

func (k *ChainKeeper) QueryBalance(ctx context.Context, address string) (uint64, error) {
    resp, err := k.queryClient.Balance(ctx, &types.QueryBalanceRequest{Address: address})
    return resp.Balance.Amount, err
}
```

- [x] 创建 keeper 与链交互
- [x] 实现余额查询 (mock模式，真实链集成可配置)
- [x] 实现任务创建
- [x] 实现托管创建

### Phase 4: GenieBot 集成 (30分钟)

**复用 LLMClient**

```go
// agentgateway/keeper/genie.go
func (k *Keeper) ChatWithGenie(ctx context.Context, sessionID, message string) (*ChatResponse, error) {
    // 直接调用现有 LLMClient
    result, err := k.llmClient.Invoke("claude", message)
    if err != nil {
        return nil, err
    }

    // 扣费
    if err := k.ChargeUser(ctx, sessionID, result.TotalCost); err != nil {
        return nil, err
    }

    return &ChatResponse{
        Content: result.Response,
        Cost: result.TotalCost,
    }, nil
}
```

- [x] 实现 chat_with_genie handler
- [x] 会话管理
- [x] 基础计费逻辑

### Phase 5: 认证与安全 (30分钟)

**简单钱包签名验证**

```go
// agentgateway/keeper/auth.go
func (k *Keeper) Authenticate(req *http.Request) (string, error) {
    sig := req.Header.Get("X-Wallet-Signature")
    addr := req.Header.Get("X-Wallet-Address")

    if !verifySignature(addr, sig) {
        return "", errors.New("invalid signature")
    }
    return addr, nil
}
```

- [x] 钱包签名验证 (mock实现)
- [x] 基础速率限制 (60 req/min)
- [x] 上下文传递

### Phase 6: HTTP/WebSocket (30分钟)

**标准 Go HTTP Server**

```go
// cmd/agent-gateway/server.go
func startHTTPServer(keeper *keeper.Keeper) {
    mux := http.NewServeMux()

    // MCP over HTTP
    mux.Handle("/mcp/", mcpHandler(keeper))

    // A2A
    mux.Handle("/a2a/", a2a.Handler())

    // WebSocket
    mux.HandleFunc("/ws", wsHandler(keeper))

    http.ListenAndServe(":8080", mux)
}
```

- [x] HTTP Server (端口可配置)
- [x] MCP over HTTP endpoint
- [x] A2A endpoints
- [x] WebSocket handler stub
- [x] SSE 流式响应 stub

---

## 极简配置

### Claude Code (5分钟配置)
```json
{
  "mcpServers": {
    "sharetoken": {
      "command": "go",
      "args": ["run", "./cmd/agent-gateway", "--transport", "stdio"]
    }
  }
}
```

### 启动命令
```bash
# stdio 模式 (Claude Code)
go run ./cmd/agent-gateway --transport stdio

# HTTP 模式
go run ./cmd/agent-gateway --transport http --port 8080
```

---

## 实现顺序 (建议)

| 顺序 | 任务 | 时间 | 验证方式 |
|------|------|------|---------|
| 1 | MCP + 1个Tool | 15min | Claude Code 测试 |
| 2 | 链交互 | 15min | 余额查询测试 |
| 3 | GenieBot | 20min | 对话测试 |
| 4 | A2A | 15min | curl 测试 |
| 5 | HTTP/WebSocket | 20min | 浏览器测试 |
| 6 | 认证 | 15min | 签名验证 |

**总计**: 约 2小时

---

## 文件清单 (最少代码)

```
cmd/agent-gateway/
├── main.go              # 入口 (50行)
├── mcp_handler.go       # MCP处理 (100行)
└── server.go            # HTTP服务 (80行)

x/agentgateway/
├── keeper/
│   ├── keeper.go        # 核心逻辑 (150行)
│   ├── chain.go         # 链交互 (80行)
│   └── auth.go          # 认证 (50行)
└── a2a/
    └── handler.go       # A2A处理 (100行)

# 总计: 约 600 行代码
```

---

## 测试命令

```bash
# 1. 启动 Gateway
go run ./cmd/agent-gateway --transport http --port 8080

# 2. 测试 MCP Tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "query_balance",
      "arguments": {"address": "cosmos1..."}
    }
  }'

# 3. 测试 A2A Agent Card
curl http://localhost:8080/.well-known/agent.json

# 4. WebSocket 测试
websocat ws://localhost:8080/ws
```

---

## 注意事项

1. **快速迭代**: 先跑通1个Tool，再扩展
2. **复用代码**: ChainHelper、LLMClient 直接复用
3. **最小实现**: 不需要完整 Cosmos SDK 模块，独立程序即可
4. **测试驱动**: 每实现一个功能，立即用 Claude Code 测试

---

## 下一步

开始 Phase 1: MCP 核心实现

```bash
# 1. 添加依赖
go get github.com/mark3labs/mcp-go

# 2. 创建入口文件
cmd/agent-gateway/main.go

# 3. 实现第一个 Tool
cmd/agent-gateway/handlers.go

# 4. 测试
# 在 Claude Code 中配置 MCP 并测试
```

需要开始实现吗？
