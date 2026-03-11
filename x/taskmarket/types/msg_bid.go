package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for bids and auctions
const (
	TypeMsgSubmitBid    = "submit_bid"
	TypeMsgCloseAuction = "close_auction"
)

// Response types
type MsgSubmitBidResponse struct {
	BidID string `json:"bid_id"`
}
type MsgCloseAuctionResponse struct {
	WinnerID string `json:"winner_id"`
}

// -----------------------------------------------------------------------------
// MsgSubmitBid
// -----------------------------------------------------------------------------

type MsgSubmitBid struct {
	WorkerID  string `json:"worker_id"`
	TaskID    string `json:"task_id"`
	Amount    uint64 `json:"amount"`
	Message   string `json:"message"`
	Portfolio string `json:"portfolio"`
}

func NewMsgSubmitBid(workerID, taskID string, amount uint64) *MsgSubmitBid {
	return &MsgSubmitBid{
		WorkerID: workerID,
		TaskID:   taskID,
		Amount:   amount,
	}
}

func (msg MsgSubmitBid) Route() string { return RouterKey }
func (msg MsgSubmitBid) Type() string  { return TypeMsgSubmitBid }

func (msg MsgSubmitBid) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WorkerID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSubmitBid) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSubmitBid) ValidateBasic() error {
	if msg.WorkerID == "" {
		return fmt.Errorf("worker ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if msg.Amount == 0 {
		return fmt.Errorf("bid amount must be greater than 0")
	}
	return nil
}

// -----------------------------------------------------------------------------
// MsgCloseAuction
// -----------------------------------------------------------------------------

type MsgCloseAuction struct {
	RequesterID string `json:"requester_id"`
	TaskID      string `json:"task_id"`
}

func NewMsgCloseAuction(requesterID, taskID string) *MsgCloseAuction {
	return &MsgCloseAuction{
		RequesterID: requesterID,
		TaskID:      taskID,
	}
}

func (msg MsgCloseAuction) Route() string { return RouterKey }
func (msg MsgCloseAuction) Type() string  { return TypeMsgCloseAuction }

func (msg MsgCloseAuction) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCloseAuction) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgCloseAuction) ValidateBasic() error {
	if msg.RequesterID == "" {
		return fmt.Errorf("requester ID cannot be empty")
	}
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	return nil
}
