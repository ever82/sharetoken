package types

import (
	"errors"
	"fmt"
)

// ValidateBasic performs basic validation of MsgCreateTask
func (m MsgCreateTask) ValidateBasic() error {
	if m.Creator == "" {
		return errors.New("creator cannot be empty")
	}
	if m.Title == "" {
		return errors.New("title cannot be empty")
	}
	if m.Budget == 0 {
		return errors.New("budget must be greater than 0")
	}
	if m.TaskType == TaskType_TASK_TYPE_UNSPECIFIED {
		return errors.New("task type must be specified")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgCreateTask) GetSigners() []string {
	return []string{m.Creator}
}

// Type returns the message type
func (m MsgCreateTask) Type() string {
	return "create_task"
}

// ValidateBasic performs basic validation of MsgUpdateTask
func (m MsgUpdateTask) ValidateBasic() error {
	if m.Creator == "" {
		return errors.New("creator cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgUpdateTask) GetSigners() []string {
	return []string{m.Creator}
}

// Type returns the message type
func (m MsgUpdateTask) Type() string {
	return "update_task"
}

// ValidateBasic performs basic validation of MsgPublishTask
func (m MsgPublishTask) ValidateBasic() error {
	if m.Creator == "" {
		return errors.New("creator cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgPublishTask) GetSigners() []string {
	return []string{m.Creator}
}

// ValidateBasic performs basic validation of MsgCancelTask
func (m MsgCancelTask) ValidateBasic() error {
	if m.Creator == "" {
		return errors.New("creator cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgCancelTask) GetSigners() []string {
	return []string{m.Creator}
}

// ValidateBasic performs basic validation of MsgStartTask
func (m MsgStartTask) ValidateBasic() error {
	if m.WorkerId == "" {
		return errors.New("worker ID cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgStartTask) GetSigners() []string {
	return []string{m.WorkerId}
}

// ValidateBasic performs basic validation of MsgSubmitApplication
func (m MsgSubmitApplication) ValidateBasic() error {
	if m.WorkerId == "" {
		return errors.New("worker ID cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgSubmitApplication) GetSigners() []string {
	return []string{m.WorkerId}
}

// ValidateBasic performs basic validation of MsgAcceptApplication
func (m MsgAcceptApplication) ValidateBasic() error {
	if m.RequesterId == "" {
		return errors.New("requester ID cannot be empty")
	}
	if m.ApplicationId == "" {
		return errors.New("application ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgAcceptApplication) GetSigners() []string {
	return []string{m.RequesterId}
}

// ValidateBasic performs basic validation of MsgRejectApplication
func (m MsgRejectApplication) ValidateBasic() error {
	if m.RequesterId == "" {
		return errors.New("requester ID cannot be empty")
	}
	if m.ApplicationId == "" {
		return errors.New("application ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgRejectApplication) GetSigners() []string {
	return []string{m.RequesterId}
}

// ValidateBasic performs basic validation of MsgSubmitBid
func (m MsgSubmitBid) ValidateBasic() error {
	if m.WorkerId == "" {
		return errors.New("worker ID cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	if m.Amount == 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgSubmitBid) GetSigners() []string {
	return []string{m.WorkerId}
}

// ValidateBasic performs basic validation of MsgCloseAuction
func (m MsgCloseAuction) ValidateBasic() error {
	if m.RequesterId == "" {
		return errors.New("requester ID cannot be empty")
	}
	if m.TaskId == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	return nil
}

// GetSigners returns the signers of the message
func (m MsgCloseAuction) GetSigners() []string {
	return []string{m.RequesterId}
}
