// MCP Handler 实现
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"sharetoken/x/agentgateway/keeper"
)

// MCPRequest MCP 请求
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

// MCPResponse MCP 响应
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// MCPError MCP 错误
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToolHandler 工具处理函数类型
type ToolHandler func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// MCPServer MCP 服务器
type MCPServer struct {
	keeper  *keeper.Keeper
	tools   map[string]ToolHandler
	version string
	name    string
}

// NewMCPServer 创建 MCP 服务器
func NewMCPServer(keeper *keeper.Keeper) *MCPServer {
	s := &MCPServer{
		keeper:  keeper,
		tools:   make(map[string]ToolHandler),
		version: "1.0.0",
		name:    "sharetoken-gateway",
	}

	// 注册工具
	s.registerTools()

	return s
}

// registerTools 注册所有工具
func (s *MCPServer) registerTools() {
	// query_balance 工具
	s.tools["query_balance"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		address, ok := args["address"].(string)
		if !ok || address == "" {
			return nil, fmt.Errorf("address is required")
		}

		balance, err := s.keeper.QueryBalance(ctx, address)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"address": address,
			"balance": balance,
			"denom":   "stt",
		}, nil
	}

	// chat_with_genie 工具
	s.tools["chat_with_genie"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		message, ok := args["message"].(string)
		if !ok || message == "" {
			return nil, fmt.Errorf("message is required")
		}

		sessionID := "default-session"
		if sid, ok := args["session_id"].(string); ok && sid != "" {
			sessionID = sid
		}

		resp, err := s.keeper.ChatWithGenie(ctx, sessionID, message)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"content": resp.Content,
			"cost":    resp.Cost,
		}, nil
	}

	// create_task 工具
	s.tools["create_task"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		description, ok := args["description"].(string)
		if !ok || description == "" {
			return nil, fmt.Errorf("description is required")
		}

		budget := "100stt"
		if b, ok := args["budget"].(string); ok && b != "" {
			budget = b
		}

		userAddr := "cosmos1user"
		taskID, err := s.keeper.CreateTask(ctx, userAddr, description, budget)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"task_id":     taskID,
			"description": description,
			"budget":      budget,
			"status":      "created",
		}, nil
	}

	// create_escrow 工具
	s.tools["create_escrow"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		provider, ok := args["provider"].(string)
		if !ok || provider == "" {
			return nil, fmt.Errorf("provider is required")
		}

		amount, ok := args["amount"].(string)
		if !ok || amount == "" {
			return nil, fmt.Errorf("amount is required")
		}

		userAddr := "cosmos1user"
		escrowID, err := s.keeper.CreateEscrow(ctx, userAddr, provider, amount)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"escrow_id": escrowID,
			"provider":  provider,
			"amount":    amount,
			"status":    "locked",
		}, nil
	}

	// extract_functions 工具 - 从代码中提取函数名
	s.tools["extract_functions"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		code, ok := args["code"].(string)
		if !ok || code == "" {
			return nil, fmt.Errorf("code is required")
		}

		language := "python"
		if lang, ok := args["language"].(string); ok && lang != "" {
			language = lang
		}

		// 使用外部 Agent 调用 Claude Code
		prompt := fmt.Sprintf("Extract all function names from the following %s code and return them as a JSON array:\n\n```%s\n%s\n```", language, language, code)

		response, err := s.keeper.CallExternalAgentWithOptions(prompt, keeper.ExternalAgentOptions{
			OutputFormat: "json",
			JSONSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"functions": map[string]interface{}{
						"type": "array",
						"items": map[string]string{
							"type": "string",
						},
					},
				},
				"required": []string{"functions"},
			},
		})
		if err != nil {
			return nil, err
		}

		// 尝试解析 JSON 响应
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			// 如果不是 JSON，直接返回文本
			return map[string]interface{}{
				"functions": []string{},
				"raw":       response,
			}, nil
		}

		return result, nil
	}

	// analyze_code 工具 - 分析代码结构和复杂度
	s.tools["analyze_code"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		code, ok := args["code"].(string)
		if !ok || code == "" {
			return nil, fmt.Errorf("code is required")
		}

		language := "python"
		if lang, ok := args["language"].(string); ok && lang != "" {
			language = lang
		}

		prompt := fmt.Sprintf("Analyze the following %s code and provide:\n1. Function count\n2. Cyclomatic complexity estimate\n3. Dependencies used\n4. Code structure summary\n\nReturn as JSON:\n\n```%s\n%s\n```", language, language, code)

		response, err := s.keeper.CallExternalAgentWithOptions(prompt, keeper.ExternalAgentOptions{
			OutputFormat: "json",
			JSONSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"function_count": map[string]string{
						"type": "integer",
					},
					"complexity": map[string]string{
						"type": "string",
					},
					"dependencies": map[string]interface{}{
						"type": "array",
						"items": map[string]string{
							"type": "string",
						},
					},
					"summary": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"function_count", "summary"},
			},
		})
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			return map[string]interface{}{
				"function_count": 0,
				"summary":        response,
			}, nil
		}

		return result, nil
	}

	// search_code 工具 - 在代码库中搜索特定模式
	s.tools["search_code"] = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		query, ok := args["query"].(string)
		if !ok || query == "" {
			return nil, fmt.Errorf("query is required")
		}

		language := ""
		if lang, ok := args["language"].(string); ok {
			language = lang
		}

		// 这里可以实现真正的代码搜索
		// 目前返回模拟结果
		return map[string]interface{}{
			"query":    query,
			"language": language,
			"results": []map[string]string{
				{
					"file":    "example.py",
					"line":    "10",
					"content": "def example_function():",
				},
			},
			"total": 1,
		}, nil
	}
}

// Handle 处理 MCP 请求
func (s *MCPServer) Handle(req MCPRequest) MCPResponse {
	ctx := context.Background()

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolCall(ctx, req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
			ID: req.ID,
		}
	}
}

// handleInitialize 处理初始化
func (s *MCPServer) handleInitialize(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]string{
				"name":    s.name,
				"version": s.version,
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		},
		ID: req.ID,
	}
}

// handleToolsList 处理工具列表
func (s *MCPServer) handleToolsList(req MCPRequest) MCPResponse {
	tools := []map[string]interface{}{
		{
			"name":        "query_balance",
			"description": "查询账户STT余额",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"address": map[string]string{
						"type":        "string",
						"description": "账户地址",
					},
				},
				"required": []string{"address"},
			},
		},
		{
			"name":        "chat_with_genie",
			"description": "与GenieBot AI助手对话",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]string{
						"type":        "string",
						"description": "对话消息",
					},
					"session_id": map[string]string{
						"type":        "string",
						"description": "会话ID（可选）",
					},
				},
				"required": []string{"message"},
			},
		},
		{
			"name":        "create_task",
			"description": "创建新任务",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"description": map[string]string{
						"type":        "string",
						"description": "任务描述",
					},
					"budget": map[string]string{
						"type":        "string",
						"description": "预算（如: 100stt）",
					},
				},
				"required": []string{"description"},
			},
		},
		{
			"name":        "create_escrow",
			"description": "创建资金托管",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"provider": map[string]string{
						"type":        "string",
						"description": "服务提供者地址",
					},
					"amount": map[string]string{
						"type":        "string",
						"description": "托管金额（如: 1000stt）",
					},
				},
				"required": []string{"provider", "amount"},
			},
		},
		{
			"name":        "extract_functions",
			"description": "从代码文件中提取函数名",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]string{
						"type":        "string",
						"description": "代码内容",
					},
					"language": map[string]string{
						"type":        "string",
						"description": "编程语言（默认: python）",
					},
				},
				"required": []string{"code"},
			},
		},
		{
			"name":        "analyze_code",
			"description": "分析代码结构和复杂度",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]string{
						"type":        "string",
						"description": "代码内容",
					},
					"language": map[string]string{
						"type":        "string",
						"description": "编程语言（默认: python）",
					},
				},
				"required": []string{"code"},
			},
		},
		{
			"name":        "search_code",
			"description": "在代码库中搜索特定模式",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]string{
						"type":        "string",
						"description": "搜索查询",
					},
					"language": map[string]string{
						"type":        "string",
						"description": "编程语言过滤（可选）",
					},
				},
				"required": []string{"query"},
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"tools": tools,
		},
		ID: req.ID,
	}
}

// handleToolCall 处理工具调用
func (s *MCPServer) handleToolCall(ctx context.Context, req MCPRequest) MCPResponse {
	// 解析参数
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
			ID: req.ID,
		}
	}

	// 查找工具
	handler, exists := s.tools[params.Name]
	if !exists {
		return MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Tool not found: %s", params.Name),
			},
			ID: req.ID,
		}
	}

	// 执行工具
	result, err := handler(ctx, params.Arguments)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32603,
				Message: err.Error(),
			},
			ID: req.ID,
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": toJSON(result),
				},
			},
		},
		ID: req.ID,
	}
}

// toJSON 转换为 JSON 字符串
func toJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return fmt.Sprintf("{\"error\": \"failed to marshal result: %v\"}", err)
	}
	return string(b)
}

// ServeStdio stdio 模式服务
func ServeStdio(server *MCPServer) {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Decode error: %v", err)
			continue
		}

		resp := server.Handle(req)
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Encode error: %v", err)
		}
	}
}
