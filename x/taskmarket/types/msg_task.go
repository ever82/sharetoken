package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for task tasks
const (
	TypeMsgCreateTask  = "create_task"
	TypeMsgUpdateTask  = "update_task"
	TypeMsgPublishTask = "publish_task"
	TypeMsgCancelTask  = "cancel_task"
	TypeMsgStartTask   = "start_task"
)

// Response types
type MsgCreateTaskResponse struct {
	TaskID string `json:"task_id"`
}

type MsgUpdateTaskResponse struct{}
type MsgPublishTaskResponse struct{}
type MsgCancelTaskResponse struct{}
type MsgStartTaskResponse struct{}

// -----------------------------------------------------------------------------
// MsgCreateTask
// -----------------------------------------------------------------------------

type MsgCreateTask struct {
	Creator     string       `json:"creator"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	TaskTypeVal TaskType     `json:"task_type"`
	Category    TaskCategory `json:"category"`
	Budget      uint64       `json:"budget"`
	Currency    string       `json:"currency"`
	Deadline    int64        `json:"deadline"`
	Skills      []string     `json:"skills"`
	Subtasks    []Subtask    `json:"subtasks"`
	Milestones  []Milestone  `json:"milestones"`
}

func NewMsgCreateTask(creator, title, description string, taskType TaskType, category TaskCategory, budget uint64) *MsgCreateTask {
	return &MsgCreateTask{
		Creator:     creator,
		Title:       title,
		Description: description,
		TaskTypeVal: taskType,
		Category:    category,
		Budget:      budget,
		Currency:    "STT",
		Skills:      []string{},
		Subtasks:    []Subtask{},
		Milestones:  []Milestone{},
	}
}

func (msg MsgCreateTask) Route() string { return RouterKey }
func (msg MsgCreateTask) Type() string  { return TypeMsgCreateTask }

func (msg MsgCreateTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateTask) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgCreateTask) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}
	if msg.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if msg.Budget == 0 {
		return fmt.Errorf("budget must be greater than 0")
	}
	if msg.TaskTypeVal != TaskTypeOpen && msg.TaskTypeVal != TaskTypeAuction {
		return fmt.Errorf("invalid task type")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgUpdateTask
// -----------------------------------------------------------------------------

type MsgUpdateTask struct {
	Creator     string      `json:"creator"`
	TaskID      string      `json:"task_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Budget      uint64      `json:"budget"`
	Deadline    int64       `json:"deadline"`
	Skills      []string    `json:"skills"`
	Subtasks    []Subtask   `json:"subtasks"`
	Milestones  []Milestone `json:"milestones"`
}

func NewMsgUpdateTask(creator, taskID string) *MsgUpdateTask {
	return &MsgUpdateTask{
		Creator: creator,
		TaskID:  taskID,
		Skills:  []string{},
	}
}

func (msg MsgUpdateTask) Route() string { return RouterKey }
func (msg MsgUpdateTask) Type() string  { return TypeMsgUpdateTask }

func (msg MsgUpdateTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUpdateTask) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgUpdateTask) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgPublishTask
// -----------------------------------------------------------------------------

type MsgPublishTask struct {
	Creator string `json:"creator"`
	TaskID  string `json:"task_id"`
}

func NewMsgPublishTask(creator, taskID string) *MsgPublishTask {
	return &MsgPublishTask{
		Creator: creator,
		TaskID:  taskID,
	}
}

func (msg MsgPublishTask) Route() string { return RouterKey }
func (msg MsgPublishTask) Type() string  { return TypeMsgPublishTask }

func (msg MsgPublishTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgPublishTask) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgPublishTask) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgCancelTask
// -----------------------------------------------------------------------------

type MsgCancelTask struct {
	Creator string `json:"creator"`
	TaskID  string `json:"task_id"`
}

func NewMsgCancelTask(creator, taskID string) *MsgCancelTask {
	return &MsgCancelTask{
		Creator: creator,
		TaskID:  taskID,
	}
}

func (msg MsgCancelTask) Route() string { return RouterKey }
func (msg MsgCancelTask) Type() string  { return TypeMsgCancelTask }

func (msg MsgCancelTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelTask) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgCancelTask) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgStartTask
// -----------------------------------------------------------------------------

type MsgStartTask struct {
	WorkerID string `json:"worker_id"`
	TaskID   string `json:"task_id"`
}

func NewMsgStartTask(workerID, taskID string) *MsgStartTask {
	return &MsgStartTask{
		WorkerID: workerID,
		TaskID:   taskID,
	}
}

func (msg MsgStartTask) Route() string { return RouterKey }
func (msg MsgStartTask) Type() string  { return TypeMsgStartTask }

func (msg MsgStartTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WorkerID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgStartTask) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgStartTask) ValidateBasic() error {
	if msg.WorkerID == "" {
		return fmt.Errorf("worker ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	return nil
}
