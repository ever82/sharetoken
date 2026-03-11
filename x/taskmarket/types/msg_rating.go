package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for ratings
const (
	TypeMsgSubmitRating = "submit_rating"
)

// Response types
type MsgSubmitRatingResponse struct {
	RatingID string `json:"rating_id"`
}

// -----------------------------------------------------------------------------
// MsgSubmitRating
// -----------------------------------------------------------------------------

type MsgSubmitRating struct {
	TaskID  string         `json:"task_id"`
	RaterID string         `json:"rater_id"`
	RatedID string         `json:"rated_id"`
	Ratings map[string]int `json:"ratings"` // dimension -> score (1-5)
	Comment string         `json:"comment"`
}

func NewMsgSubmitRating(taskID, raterID, ratedID string) *MsgSubmitRating {
	return &MsgSubmitRating{
		TaskID:  taskID,
		RaterID: raterID,
		RatedID: ratedID,
		Ratings: make(map[string]int),
	}
}

func (msg MsgSubmitRating) Route() string { return RouterKey }
func (msg MsgSubmitRating) Type() string  { return TypeMsgSubmitRating }

func (msg MsgSubmitRating) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RaterID)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSubmitRating) GetSignBytes() []byte {
	b, _ := json.Marshal(&msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSubmitRating) ValidateBasic() error {
	if msg.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if msg.RaterID == "" {
		return fmt.Errorf("rater ID cannot be empty")
	}
	if msg.RatedID == "" {
		return fmt.Errorf("rated ID cannot be empty")
	}
	if len(msg.Ratings) == 0 {
		return fmt.Errorf("at least one rating dimension required")
	}
	return nil
}
