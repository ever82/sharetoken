// A2A Handler 实现
package a2a

import (
	"encoding/json"
	"net/http"

	"sharetoken/x/agentgateway/keeper"
)

// Handler A2A HTTP Handler
type Handler struct {
	keeper *keeper.Keeper
}

// NewHandler 创建 A2A Handler
func NewHandler(keeper *keeper.Keeper) *Handler {
	return &Handler{keeper: keeper}
}

// ServeHTTP 实现 http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/.well-known/agent.json":
		h.handleAgentCard(w, r)
	case "/a2a/tasks":
		h.handleTasks(w, r)
	case "/a2a/status":
		h.handleStatus(w, r)
	case "/a2a/negotiate":
		h.handleNegotiate(w, r)
	default:
		h.writeError(w, http.StatusNotFound, "Not found")
	}
}

// handleAgentCard 处理 Agent Card 请求
func (h *Handler) handleAgentCard(w http.ResponseWriter, r *http.Request) {
	card := h.keeper.GetAgentCard()

	response := map[string]interface{}{
		"name":         card.Name,
		"version":      card.Version,
		"description":  card.Description,
		"capabilities": card.Capabilities,
		"endpoints":    card.Endpoints,
		"authentication": map[string]string{
			"type":  "wallet_signature",
			"chain": "cosmos",
		},
	}

	json.NewEncoder(w).Encode(response)
}

// handleTasks 处理任务请求
func (h *Handler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 列出任务
		h.listTasks(w, r)
	case http.MethodPost:
		// 创建任务
		h.createTask(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// listTasks 列出任务
func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现真实任务列表
	tasks := []map[string]interface{}{
		{
			"id":          "task-1",
			"description": "开发DeFi仪表盘",
			"status":      "open",
			"budget":      "500stt",
		},
		{
			"id":          "task-2",
			"description": "智能合约审计",
			"status":      "in_progress",
			"budget":      "1000stt",
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
	})
}

// createTask 创建任务
func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Description string `json:"description"`
		Budget      string `json:"budget"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: 从请求中提取用户地址
	userAddr := "cosmos1user"

	taskID, err := h.keeper.CreateTask(r.Context(), userAddr, req.Description, req.Budget)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id": taskID,
		"status":  "created",
	})
}

// handleStatus 处理状态查询
func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":    "online",
		"version":   "1.0.0",
		"uptime":    "99.9%",
		"capabilities": []string{
			"task_execution",
			"escrow_management",
			"query_service",
		},
	}

	json.NewEncoder(w).Encode(status)
}

// handleNegotiate 处理协商请求
func (h *Handler) handleNegotiate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		TaskID   string `json:"task_id"`
		Provider string `json:"provider"`
		Bid      string `json:"bid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: 实现真实协商逻辑
	response := map[string]interface{}{
		"task_id":   req.TaskID,
		"provider":  req.Provider,
		"bid":       req.Bid,
		"status":    "accepted",
		"message":   "Bid accepted",
	}

	json.NewEncoder(w).Encode(response)
}

// writeError 写入错误响应
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})
}

// Routes 返回 A2A 路由
func Routes(keeper *keeper.Keeper) http.Handler {
	return NewHandler(keeper)
}
