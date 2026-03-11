package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for milestones
const (
	TypeMsgSubmitMilestone  = "submit_milestone"
	TypeMsgApproveMilestone = "approve_milestone"
	TypeMsgRejectMilestone  = "reject_milestone"
)

// Response types
type MsgSubmitMilestoneResponse struct{}
type MsgApproveMilestoneResponse struct{}
type MsgRejectMilestoneResponse struct{}

// -----------------------------------------------------------------------------
// MsgSubmitMilestone
// -----------------------------------------------------------------------------

type MsgSubmitMilestone struct {
	WorkerID     string `json:"worker_id"`
	TaskID       string `json:"task_id"`
	MilestoneID  string `json:"milestone_id"`
	Deliverables string `json:"deliverables"`
}

func NewMsgSubmitMilestone(workerID, taskID, milestoneID, deliverables string) *MsgSubmitMilestone {
	return &MsgSubmitMilestone{
		WorkerID:     workerID,
		TaskID:       taskID,
		MilestoneID:  milestoneID,
		Deliverables: deliverables,
	}
}

func (msg MsgSubmitMilestone) Route() string { return RouterKey }
func (msg MsgSubmitMilestone) Type() string  { return TypeMsgSubmitMilestone }

func (msg MsgSubmitMilestone) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WorkerID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSubmitMilestone) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSubmitMilestone) ValidateBasic() error {
	if msg.WorkerID == "" {
		return fmt.Errorf("worker ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if msg.MilestoneID == "" {
		return fmt.Errorf("milestone ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgApproveMilestone
// -----------------------------------------------------------------------------

type MsgApproveMilestone struct {
	RequesterID string `json:"requester_id"`
	TaskID      string `json:"task_id"`
	MilestoneID string `json:"milestone_id"`
}

func NewMsgApproveMilestone(requesterID, taskID, milestoneID string) *MsgApproveMilestone {
	return &MsgApproveMilestone{
		RequesterID: requesterID,
		TaskID:      taskID,
		MilestoneID: milestoneID,
	}
}

func (msg MsgApproveMilestone) Route() string { return RouterKey }
func (msg MsgApproveMilestone) Type() string  { return TypeMsgApproveMilestone }

func (msg MsgApproveMilestone) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgApproveMilestone) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgApproveMilestone) ValidateBasic() error {
	if msg.RequesterID == "" {
		return fmt.Errorf("requester ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if msg.MilestoneID == "" {
		return fmt.Errorf("milestone ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgRejectMilestone
// -----------------------------------------------------------------------------

type MsgRejectMilestone struct {
	RequesterID string `json:"requester_id"`
	TaskID      string `json:"task_id"`
	MilestoneID string `json:"milestone_id"`
	Reason      string `json:"reason"`
}

func NewMsgRejectMilestone(requesterID, taskID, milestoneID, reason string) *MsgRejectMilestone {
	return &MsgRejectMilestone{
		RequesterID: requesterID,
		TaskID:      taskID,
		MilestoneID: milestoneID,
		Reason:      reason,
	}
}

func (msg MsgRejectMilestone) Route() string { return RouterKey }
func (msg MsgRejectMilestone) Type() string  { return TypeMsgRejectMilestone }

func (msg MsgRejectMilestone) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgRejectMilestone) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgRejectMilestone) ValidateBasic() error {
	if msg.RequesterID == "" {
		return fmt.Errorf("requester ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if msg.MilestoneID == "" {
		return fmt.Errorf("milestone ID cannot be empty")
	}
	return nil
}
