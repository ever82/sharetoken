package keeper

import (
	"context"
	"time"
)

// Session 会话状态
type Session struct {
	ID        string
	UserAddr  string
	Context   []Message
	CreatedAt time.Time
}

// Message 消息
type Message struct {
	Role    string
	Content string
}

// getOrCreateSession 获取或创建会话
func (k *Keeper) getOrCreateSession(sessionID string) *Session {
	k.sessionMu.Lock()
	defer k.sessionMu.Unlock()

	if session, exists := k.sessions[sessionID]; exists {
		return session
	}

	session := &Session{
		ID:        sessionID,
		Context:   make([]Message, 0),
		CreatedAt: time.Now(),
	}
	k.sessions[sessionID] = session
	return session
}

// ChatWithGenie 与 GenieBot 对话
func (k *Keeper) ChatWithGenie(ctx context.Context, sessionID, message string) (*ChatResponse, error) {
	// 获取或创建会话
	session := k.getOrCreateSession(sessionID)

	// 添加到上下文
	session.Context = append(session.Context, Message{
		Role:    "user",
		Content: message,
	})

	// 调用 LLM
	response, cost, err := k.callLLM(session.Context)
	if err != nil {
		return nil, err
	}

	// 保存助手响应到上下文
	session.Context = append(session.Context, Message{
		Role:    "assistant",
		Content: response,
	})

	return &ChatResponse{
		Content: response,
		Cost:    cost,
	}, nil
}
