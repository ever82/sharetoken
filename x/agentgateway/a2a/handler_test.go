// A2A Handler 单元测试
package a2a

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"sharetoken/x/agentgateway/keeper"
)

// A2AHandlerTestSuite A2A 处理器测试套件
type A2AHandlerTestSuite struct {
	suite.Suite
	handler *Handler
	keeper  *keeper.Keeper
}

func (s *A2AHandlerTestSuite) SetupTest() {
	s.keeper = keeper.NewKeeper()
	s.handler = NewHandler(s.keeper)
}

// TestAgentCard 测试 Agent Card 端点
func (s *A2AHandlerTestSuite) TestAgentCard() {
	req := httptest.NewRequest(http.MethodGet, "/.well-known/agent.json", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)

	// 验证字段
	require.Equal(s.T(), "ShareToken Agent", response["name"])
	require.Equal(s.T(), "1.0.0", response["version"])
	require.NotEmpty(s.T(), response["description"])
	require.NotNil(s.T(), response["capabilities"])
	require.NotNil(s.T(), response["endpoints"])
	require.NotNil(s.T(), response["authentication"])

	// 验证认证信息
	auth, ok := response["authentication"].(map[string]interface{})
	require.True(s.T(), ok)
	require.Equal(s.T(), "wallet_signature", auth["type"])
	require.Equal(s.T(), "cosmos", auth["chain"])
}

// TestListTasks 测试任务列表端点
func (s *A2AHandlerTestSuite) TestListTasks() {
	req := httptest.NewRequest(http.MethodGet, "/a2a/tasks", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)

	tasks, ok := response["tasks"].([]interface{})
	require.True(s.T(), ok)
	require.GreaterOrEqual(s.T(), len(tasks), 2)

	// 验证第一个任务结构
	if len(tasks) > 0 {
		task, ok := tasks[0].(map[string]interface{})
		require.True(s.T(), ok)
		require.NotEmpty(s.T(), task["id"])
		require.NotEmpty(s.T(), task["description"])
		require.NotEmpty(s.T(), task["status"])
		require.NotEmpty(s.T(), task["budget"])
	}
}

// TestCreateTask 测试创建任务端点
func (s *A2AHandlerTestSuite) TestCreateTask() {
	payload := map[string]string{
		"description": "Test task from A2A",
		"budget":      "100stt",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), response["task_id"])
	require.Equal(s.T(), "created", response["status"])
}

// TestCreateTaskInvalidBody 测试创建任务无效请求体
func (s *A2AHandlerTestSuite) TestCreateTaskInvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/a2a/tasks", bytes.NewReader([]byte("invalid json")))
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)
	require.Contains(s.T(), response["error"], "Invalid request")
}

// TestCreateTaskMethodNotAllowed 测试创建任务方法不允许
func (s *A2AHandlerTestSuite) TestCreateTaskMethodNotAllowed() {
	req := httptest.NewRequest(http.MethodPut, "/a2a/tasks", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rec.Code)
}

// TestStatus 测试状态查询端点
func (s *A2AHandlerTestSuite) TestStatus() {
	req := httptest.NewRequest(http.MethodGet, "/a2a/status", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)

	require.Equal(s.T(), "online", response["status"])
	require.Equal(s.T(), "1.0.0", response["version"])
	require.NotEmpty(s.T(), response["uptime"])
	require.NotNil(s.T(), response["capabilities"])

	caps, ok := response["capabilities"].([]interface{})
	require.True(s.T(), ok)
	require.GreaterOrEqual(s.T(), len(caps), 3)
}

// TestNegotiate 测试协商端点
func (s *A2AHandlerTestSuite) TestNegotiate() {
	payload := map[string]string{
		"task_id":   "task-123",
		"provider":  "cosmos1provider",
		"bid":       "450stt",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/negotiate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)

	require.Equal(s.T(), "task-123", response["task_id"])
	require.Equal(s.T(), "cosmos1provider", response["provider"])
	require.Equal(s.T(), "450stt", response["bid"])
	require.Equal(s.T(), "accepted", response["status"])
	require.NotEmpty(s.T(), response["message"])
}

// TestNegotiateInvalidBody 测试协商无效请求体
func (s *A2AHandlerTestSuite) TestNegotiateInvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/a2a/negotiate", bytes.NewReader([]byte("invalid json")))
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

// TestNegotiateMethodNotAllowed 测试协商方法不允许
func (s *A2AHandlerTestSuite) TestNegotiateMethodNotAllowed() {
	req := httptest.NewRequest(http.MethodGet, "/a2a/negotiate", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rec.Code)
}

// TestNotFound 测试404端点
func (s *A2AHandlerTestSuite) TestNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/a2a/unknown", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusNotFound, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)
	require.Contains(s.T(), response["error"], "Not found")
}

// TestRoutes 测试 Routes 函数
func (s *A2AHandlerTestSuite) TestRoutes() {
	handler := Routes(s.keeper)
	require.NotNil(s.T(), handler)

	// 验证 handler 可以处理请求
	req := httptest.NewRequest(http.MethodGet, "/.well-known/agent.json", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	require.Equal(s.T(), http.StatusOK, rec.Code)
}

// TestNewHandler 测试 Handler 创建
func (s *A2AHandlerTestSuite) TestNewHandler() {
	h := NewHandler(s.keeper)
	require.NotNil(s.T(), h)
	require.Equal(s.T(), s.keeper, h.keeper)
}

// TestWriteError 测试错误写入
func (s *A2AHandlerTestSuite) TestWriteError() {
	rec := httptest.NewRecorder()

	// 直接调用 writeError 方法
	s.handler.writeError(rec, http.StatusInternalServerError, "Test error message")

	require.Equal(s.T(), http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Test error message", response["error"])
}

// TestContentType 测试响应内容类型
func (s *A2AHandlerTestSuite) TestContentType() {
	req := httptest.NewRequest(http.MethodGet, "/.well-known/agent.json", nil)
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	require.Equal(s.T(), "application/json", contentType)
}

// TestCreateTaskEmptyDescription 测试创建任务空描述
func (s *A2AHandlerTestSuite) TestCreateTaskEmptyDescription() {
	payload := map[string]string{
		"description": "",
		"budget":      "100stt",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	// 空描述应该被接受（keeper.CreateTask 处理验证）
	require.Equal(s.T(), http.StatusCreated, rec.Code)
}

// TestNegotiateEmptyFields 测试协商空字段
func (s *A2AHandlerTestSuite) TestNegotiateEmptyFields() {
	payload := map[string]string{
		"task_id":   "",
		"provider":  "",
		"bid":       "",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/negotiate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	// 空字段应该被接受
	require.Equal(s.T(), http.StatusOK, rec.Code)
}

func TestA2AHandlerSuite(t *testing.T) {
	suite.Run(t, new(A2AHandlerTestSuite))
}
