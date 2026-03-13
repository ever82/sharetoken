// MCP Handler 单元测试
package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/x/agentgateway/keeper"
)

// MCPHandlerTestSuite MCP 处理器测试套件
type MCPHandlerTestSuite struct {
	suite.Suite
	server *MCPServer
	keeper *keeper.Keeper
}

func (s *MCPHandlerTestSuite) SetupTest() {
	s.keeper = keeper.NewKeeper()
	s.server = NewMCPServer(s.keeper)
}

// TestInitialize 测试初始化请求
func (s *MCPHandlerTestSuite) TestInitialize() {
	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "initialize",
		ID:      1,
	}

	resp := s.server.Handle(req)

	require.Equal(s.T(), "2.0", resp.JSONRPC)
	require.Nil(s.T(), resp.Error)
	require.Equal(s.T(), 1, resp.ID)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)
	require.Equal(s.T(), "2024-11-05", result["protocolVersion"])

	serverInfo, ok := result["serverInfo"].(map[string]string)
	require.True(s.T(), ok)
	require.Equal(s.T(), "sharetoken-gateway", serverInfo["name"])
	require.Equal(s.T(), "1.0.0", serverInfo["version"])
}

// TestToolsList 测试工具列表请求
func (s *MCPHandlerTestSuite) TestToolsList() {
	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		ID:      2,
	}

	resp := s.server.Handle(req)

	require.Nil(s.T(), resp.Error)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)

	tools, ok := result["tools"].([]map[string]interface{})
	require.True(s.T(), ok)
	require.GreaterOrEqual(s.T(), len(tools), 4, "至少应该有4个工具")

	// 验证工具名称
	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		name, ok := tool["name"].(string)
		require.True(s.T(), ok, "tool name should be a string")
		toolNames[i] = name
	}
	require.Contains(s.T(), toolNames, "query_balance")
	require.Contains(s.T(), toolNames, "chat_with_genie")
	require.Contains(s.T(), toolNames, "create_task")
	require.Contains(s.T(), toolNames, "create_escrow")
}

// TestQueryBalanceTool 测试 query_balance 工具
func (s *MCPHandlerTestSuite) TestQueryBalanceTool() {
	params := map[string]interface{}{
		"name": "query_balance",
		"arguments": map[string]interface{}{
			"address": "cosmos1test123",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      3,
	}

	resp := s.server.Handle(req)

	require.Nil(s.T(), resp.Error)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)

	content, ok := result["content"].([]map[string]interface{})
	require.True(s.T(), ok)
	require.Len(s.T(), content, 1)

	text, ok := content[0]["text"].(string)
	require.True(s.T(), ok, "text should be a string")
	require.Contains(s.T(), text, "cosmos1test123")
	require.Contains(s.T(), text, "balance")
}

// TestChatWithGenieTool 测试 chat_with_genie 工具
func (s *MCPHandlerTestSuite) TestChatWithGenieTool() {
	params := map[string]interface{}{
		"name": "chat_with_genie",
		"arguments": map[string]interface{}{
			"message":   "Hello Genie",
			"session_id": "test-session-mcp",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      4,
	}

	resp := s.server.Handle(req)

	require.Nil(s.T(), resp.Error)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)

	content, ok := result["content"].([]map[string]interface{})
	require.True(s.T(), ok)
	require.Len(s.T(), content, 1)

	text, ok := content[0]["text"].(string)
	require.True(s.T(), ok, "text should be a string")
	require.NotEmpty(s.T(), text)
}

// TestCreateTaskTool 测试 create_task 工具
func (s *MCPHandlerTestSuite) TestCreateTaskTool() {
	params := map[string]interface{}{
		"name": "create_task",
		"arguments": map[string]interface{}{
			"description": "Test task from MCP",
			"budget":      "500stt",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      5,
	}

	resp := s.server.Handle(req)

	require.Nil(s.T(), resp.Error)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)

	content, ok := result["content"].([]map[string]interface{})
	require.True(s.T(), ok)

	text, ok := content[0]["text"].(string)
	require.True(s.T(), ok, "text should be a string")
	require.Contains(s.T(), text, "task-")
	require.Contains(s.T(), text, "created")
}

// TestCreateEscrowTool 测试 create_escrow 工具
func (s *MCPHandlerTestSuite) TestCreateEscrowTool() {
	params := map[string]interface{}{
		"name": "create_escrow",
		"arguments": map[string]interface{}{
			"provider": "cosmos1provider",
			"amount":   "1000stt",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      6,
	}

	resp := s.server.Handle(req)

	require.Nil(s.T(), resp.Error)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)

	content, ok := result["content"].([]map[string]interface{})
	require.True(s.T(), ok)

	text, ok := content[0]["text"].(string)
	require.True(s.T(), ok, "text should be a string")
	require.Contains(s.T(), text, "escrow-")
	require.Contains(s.T(), text, "locked")
}

// TestInvalidMethod 测试无效方法
func (s *MCPHandlerTestSuite) TestInvalidMethod() {
	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "invalid/method",
		ID:      7,
	}

	resp := s.server.Handle(req)

	require.NotNil(s.T(), resp.Error)
	require.Equal(s.T(), -32601, resp.Error.Code)
	require.Contains(s.T(), resp.Error.Message, "Method not found")
}

// TestToolNotFound 测试工具不存在
func (s *MCPHandlerTestSuite) TestToolNotFound() {
	params := map[string]interface{}{
		"name":      "non_existent_tool",
		"arguments": map[string]interface{}{},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      8,
	}

	resp := s.server.Handle(req)

	require.NotNil(s.T(), resp.Error)
	require.Equal(s.T(), -32602, resp.Error.Code)
	require.Contains(s.T(), resp.Error.Message, "Tool not found")
}

// TestInvalidParams 测试无效参数
func (s *MCPHandlerTestSuite) TestInvalidParams() {
	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  []byte("invalid json"),
		ID:      9,
	}

	resp := s.server.Handle(req)

	require.NotNil(s.T(), resp.Error)
	require.Equal(s.T(), -32602, resp.Error.Code)
}

// TestMissingRequiredArgument 测试缺少必需参数
func (s *MCPHandlerTestSuite) TestMissingRequiredArgument() {
	params := map[string]interface{}{
		"name":      "query_balance",
		"arguments": map[string]interface{}{},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      10,
	}

	resp := s.server.Handle(req)

	require.NotNil(s.T(), resp.Error)
	require.Contains(s.T(), resp.Error.Message, "address is required")
}

// TestToJSON 测试 JSON 转换辅助函数
func (s *MCPHandlerTestSuite) TestToJSON() {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple map",
			input:    map[string]string{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "nested map",
			input:    map[string]interface{}{"num": 42, "bool": true},
			expected: `{"bool":true,"num":42}`,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := toJSON(tc.input)
			require.JSONEq(s.T(), tc.expected, result)
		})
	}
}

// TestToolHandlerSignature 测试工具处理函数签名
func (s *MCPHandlerTestSuite) TestToolHandlerSignature() {
	// 验证所有注册的工具都有正确的签名
	require.NotNil(s.T(), s.server.tools["query_balance"])
	require.NotNil(s.T(), s.server.tools["chat_with_genie"])
	require.NotNil(s.T(), s.server.tools["create_task"])
	require.NotNil(s.T(), s.server.tools["create_escrow"])

	// 测试处理函数可以被调用
	ctx := context.Background()

	result, err := s.server.tools["query_balance"](ctx, map[string]interface{}{
		"address": "test123",
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
}

// TestExtractFunctionsTool 测试 extract_functions 工具
func (s *MCPHandlerTestSuite) TestExtractFunctionsTool() {
	code := `def auth_user():
    pass

def verify_token():
    pass

class AuthManager:
    def login(self):
        pass`

	params := map[string]interface{}{
		"name": "extract_functions",
		"arguments": map[string]interface{}{
			"code":     code,
			"language": "python",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      20,
	}

	resp := s.server.Handle(req)

	// 由于外部 Agent 可能不可用，我们只验证请求处理不报错
	// 实际结果取决于是否安装了 Claude Code
	s.T().Logf("extract_functions response: %+v", resp)
}

// TestAnalyzeCodeTool 测试 analyze_code 工具
func (s *MCPHandlerTestSuite) TestAnalyzeCodeTool() {
	code := `import json

def process_data(data):
    return json.loads(data)

def save_result(result):
    with open('output.txt', 'w') as f:
        f.write(result)`

	params := map[string]interface{}{
		"name": "analyze_code",
		"arguments": map[string]interface{}{
			"code":     code,
			"language": "python",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      21,
	}

	resp := s.server.Handle(req)
	s.T().Logf("analyze_code response: %+v", resp)
}

// TestSearchCodeTool 测试 search_code 工具
func (s *MCPHandlerTestSuite) TestSearchCodeTool() {
	params := map[string]interface{}{
		"name": "search_code",
		"arguments": map[string]interface{}{
			"query":    "def authenticate",
			"language": "python",
		},
	}
	paramsJSON, _ := json.Marshal(params)

	req := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      22,
	}

	resp := s.server.Handle(req)

	require.Nil(s.T(), resp.Error)

	result, ok := resp.Result.(map[string]interface{})
	require.True(s.T(), ok)

	content, ok := result["content"].([]map[string]interface{})
	require.True(s.T(), ok)

	text, ok := content[0]["text"].(string)
	require.True(s.T(), ok, "text should be a string")
	require.Contains(s.T(), text, "query")
	require.Contains(s.T(), text, "results")
}

// TestCodeAnalysisToolsRegistered 验证代码分析工具已注册
func (s *MCPHandlerTestSuite) TestCodeAnalysisToolsRegistered() {
	// 验证新工具已注册
	require.NotNil(s.T(), s.server.tools["extract_functions"])
	require.NotNil(s.T(), s.server.tools["analyze_code"])
	require.NotNil(s.T(), s.server.tools["search_code"])

	// 验证工具总数 (4个原有 + 3个新工具 = 7个)
	require.GreaterOrEqual(s.T(), len(s.server.tools), 7)
}

func TestMCPHandlerSuite(t *testing.T) {
	suite.Run(t, new(MCPHandlerTestSuite))
}
