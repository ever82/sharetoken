package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for applications
const (
	TypeMsgSubmitApplication = "submit_application"
	TypeMsgAcceptApplication = "accept_application"
	TypeMsgRejectApplication = "reject_application"
)

// Response types
type MsgSubmitApplicationResponse struct {
	ApplicationID string `json:"application_id"`
}
type MsgAcceptApplicationResponse struct{}
type MsgRejectApplicationResponse struct{}

// -----------------------------------------------------------------------------
// MsgSubmitApplication
// -----------------------------------------------------------------------------

type MsgSubmitApplication struct {
	WorkerID           string   `json:"worker_id"`
	TaskID             string   `json:"task_id"`
	ProposedPrice      uint64   `json:"proposed_price"`
	CoverLetter        string   `json:"cover_letter"`
	RelevantExperience []string `json:"relevant_experience"`
	PortfolioLinks     []string `json:"portfolio_links"`
	EstimatedDuration  int64    `json:"estimated_duration"`
}

func NewMsgSubmitApplication(workerID, taskID string, price uint64) *MsgSubmitApplication {
	return &MsgSubmitApplication{
		WorkerID:      workerID,
		TaskID:        taskID,
		ProposedPrice: price,
	}
}

func (msg MsgSubmitApplication) Route() string { return RouterKey }
func (msg MsgSubmitApplication) Type() string  { return TypeMsgSubmitApplication }

func (msg MsgSubmitApplication) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WorkerID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSubmitApplication) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSubmitApplication) ValidateBasic() error {
	if msg.WorkerID == "" {
		return fmt.Errorf("worker ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if msg.ProposedPrice == 0 {
		return fmt.Errorf("proposed price must be greater than 0")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgAcceptApplication
// -----------------------------------------------------------------------------

type MsgAcceptApplication struct {
	RequesterID   string `json:"requester_id"`
	ApplicationID string `json:"application_id"`
}

func NewMsgAcceptApplication(requesterID, applicationID string) *MsgAcceptApplication {
	return &MsgAcceptApplication{
		RequesterID:   requesterID,
		ApplicationID: applicationID,
	}
}

func (msg MsgAcceptApplication) Route() string { return RouterKey }
func (msg MsgAcceptApplication) Type() string  { return TypeMsgAcceptApplication }

func (msg MsgAcceptApplication) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgAcceptApplication) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgAcceptApplication) ValidateBasic() error {
	if msg.RequesterID == "" {
		return fmt.Errorf("requester ID cannot be empty")
	}
	if msg.ApplicationID == "" {
		return fmt.Errorf("application ID cannot be empty")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgRejectApplication
// -----------------------------------------------------------------------------

type MsgRejectApplication struct {
	RequesterID   string `json:"requester_id"`
	ApplicationID string `json:"application_id"`
	Reason        string `json:"reason"`
}

func NewMsgRejectApplication(requesterID, applicationID string) *MsgRejectApplication {
	return &MsgRejectApplication{
		RequesterID:   requesterID,
		ApplicationID: applicationID,
	}
}

func (msg MsgRejectApplication) Route() string { return RouterKey }
func (msg MsgRejectApplication) Type() string  { return TypeMsgRejectApplication }

func (msg MsgRejectApplication) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgRejectApplication) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgRejectApplication) ValidateBasic() error {
	if msg.RequesterID == "" {
		return fmt.Errorf("requester ID cannot be empty")
	}
	if msg.ApplicationID == "" {
		return fmt.Errorf("application ID cannot be empty")
	}
	return nil
}
