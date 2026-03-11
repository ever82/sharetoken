package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for taskmarket module
const (
	TypeMsgCreateTask          = "create_task"
	TypeMsgUpdateTask          = "update_task"
	TypeMsgPublishTask         = "publish_task"
	TypeMsgCancelTask          = "cancel_task"
	TypeMsgSubmitApplication   = "submit_application"
	TypeMsgAcceptApplication   = "accept_application"
	TypeMsgRejectApplication   = "reject_application"
	TypeMsgSubmitBid           = "submit_bid"
	TypeMsgCloseAuction        = "close_auction"
	TypeMsgStartTask           = "start_task"
	TypeMsgSubmitMilestone     = "submit_milestone"
	TypeMsgApproveMilestone    = "approve_milestone"
	TypeMsgRejectMilestone     = "reject_milestone"
	TypeMsgSubmitRating        = "submit_rating"
)

// MsgServer is the message server interface
type MsgServer interface {
	CreateTask(ctx sdk.Context, msg *MsgCreateTask) (*MsgCreateTaskResponse, error)
	UpdateTask(ctx sdk.Context, msg *MsgUpdateTask) (*MsgUpdateTaskResponse, error)
	PublishTask(ctx sdk.Context, msg *MsgPublishTask) (*MsgPublishTaskResponse, error)
	CancelTask(ctx sdk.Context, msg *MsgCancelTask) (*MsgCancelTaskResponse, error)
	SubmitApplication(ctx sdk.Context, msg *MsgSubmitApplication) (*MsgSubmitApplicationResponse, error)
	AcceptApplication(ctx sdk.Context, msg *MsgAcceptApplication) (*MsgAcceptApplicationResponse, error)
	RejectApplication(ctx sdk.Context, msg *MsgRejectApplication) (*MsgRejectApplicationResponse, error)
	SubmitBid(ctx sdk.Context, msg *MsgSubmitBid) (*MsgSubmitBidResponse, error)
	CloseAuction(ctx sdk.Context, msg *MsgCloseAuction) (*MsgCloseAuctionResponse, error)
	StartTask(ctx sdk.Context, msg *MsgStartTask) (*MsgStartTaskResponse, error)
	SubmitMilestone(ctx sdk.Context, msg *MsgSubmitMilestone) (*MsgSubmitMilestoneResponse, error)
	ApproveMilestone(ctx sdk.Context, msg *MsgApproveMilestone) (*MsgApproveMilestoneResponse, error)
	RejectMilestone(ctx sdk.Context, msg *MsgRejectMilestone) (*MsgRejectMilestoneResponse, error)
	SubmitRating(ctx sdk.Context, msg *MsgSubmitRating) (*MsgSubmitRatingResponse, error)
}

// Response types
type MsgCreateTaskResponse struct {
	TaskID string `json:"task_id"`
}

type MsgUpdateTaskResponse struct{}
type MsgPublishTaskResponse struct{}
type MsgCancelTaskResponse struct{}
type MsgSubmitApplicationResponse struct {
	ApplicationID string `json:"application_id"`
}
type MsgAcceptApplicationResponse struct{}
type MsgRejectApplicationResponse struct{}
type MsgSubmitBidResponse struct {
	BidID string `json:"bid_id"`
}
type MsgCloseAuctionResponse struct {
	WinnerID string `json:"winner_id"`
}
type MsgStartTaskResponse struct{}
type MsgSubmitMilestoneResponse struct{}
type MsgApproveMilestoneResponse struct{}
type MsgRejectMilestoneResponse struct{}
type MsgSubmitRatingResponse struct {
	RatingID string `json:"rating_id"`
}

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

// -----------------------------------------------------------------------------
// MsgSubmitRating
// -----------------------------------------------------------------------------

type MsgSubmitRating struct {
	TaskID  string          `json:"task_id"`
	RaterID string          `json:"rater_id"`
	RatedID string          `json:"rated_id"`
	Ratings map[string]int  `json:"ratings"` // dimension -> score (1-5)
	Comment string          `json:"comment"`
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
