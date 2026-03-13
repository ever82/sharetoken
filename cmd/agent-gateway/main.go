// Agent Gateway 主入口
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"sharetoken/x/agentgateway/a2a"
	"sharetoken/x/agentgateway/keeper"
)

func main() {
	var (
		transport   = flag.String("transport", "stdio", "Transport type: stdio or http")
		httpPort    = flag.String("port", "8080", "HTTP server port")
		chainEndpoint = flag.String("chain-endpoint", "http://localhost:26657", "Chain RPC endpoint")
	)
	flag.Parse()

	// 创建 keeper
	k := keeper.NewKeeper()

	// TODO: 使用 chainEndpoint 配置 keeper 的链客户端
	_ = *chainEndpoint

	// 启动对应 transport
	switch *transport {
	case "stdio":
		log.Println("Starting Agent Gateway in stdio mode...")
		runStdio(k)
	case "http":
		log.Printf("Starting Agent Gateway HTTP server on port %s...", *httpPort)
		runHTTP(k, *httpPort)
	default:
		log.Fatalf("Unknown transport: %s", *transport)
	}
}

// runStdio 运行 stdio 模式
func runStdio(k *keeper.Keeper) {
	mcpServer := NewMCPServer(k)
	ServeStdio(mcpServer)
}

// runHTTP 运行 HTTP 模式
func runHTTP(k *keeper.Keeper, port string) {
	mux := http.NewServeMux()

	// MCP HTTP endpoint
	mcpServer := NewMCPServer(k)
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req MCPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		resp := mcpServer.Handle(req)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// MCP SSE endpoint (流式响应)
	mux.HandleFunc("/mcp/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// TODO: 实现 SSE 流式响应
		fmt.Fprintf(w, "data: %s\n\n", `{"type": "connected"}`)
		w.(http.Flusher).Flush()
	})

	// A2A endpoints
	mux.Handle("/.well-known/", a2a.Routes(k))
	mux.Handle("/a2a/", a2a.Routes(k))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws", handleWebSocket(k))

	addr := ":" + port
	log.Printf("Agent Gateway listening on %s", addr)
	log.Printf("MCP endpoint: http://localhost%s/mcp", addr)
	log.Printf("A2A endpoint: http://localhost%s/a2a/", addr)
	log.Printf("Agent Card: http://localhost%s/.well-known/agent.json", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// handleWebSocket WebSocket 处理
func handleWebSocket(k *keeper.Keeper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: 实现 WebSocket 支持
		// 需要引入 github.com/gorilla/websocket
		http.Error(w, "WebSocket not implemented yet", http.StatusNotImplemented)
	}
}
