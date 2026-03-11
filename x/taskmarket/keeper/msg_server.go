package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/taskmarket/types"
)

// msgServer implements the MsgServer interface
type msgServer struct {
	K *Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{K: keeper}
}

// CreateTask handles MsgCreateTask
func (m msgServer) CreateTask(ctx sdk.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	// Generate task ID (in production, use proper ID generation)
	taskID := fmt.Sprintf("task-%d", ctx.BlockHeight())

	task := types.NewTask(
		taskID,
		msg.Title,
		msg.Description,
		msg.Creator,
		msg.TaskTypeVal,
		msg.Budget,
	)

	task.Category = msg.Category
	task.Currency = msg.Currency
	task.Deadline = msg.Deadline
	task.Skills = msg.Skills
	task.Subtasks = msg.Subtasks
	task.Milestones = msg.Milestones

	if err := task.Validate(); err != nil {
		return nil, err
	}

	// Use legacy keeper for now
	if err := m.K.legacyKeeper.CreateTask(task); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTask,
			sdk.NewAttribute(types.AttributeKeyTaskID, taskID),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyTaskType, string(msg.TaskTypeVal)),
		),
	)

	return &types.MsgCreateTaskResponse{TaskID: taskID}, nil
}

// UpdateTask handles MsgUpdateTask
func (m msgServer) UpdateTask(ctx sdk.Context, msg *types.MsgUpdateTask) (*types.MsgUpdateTaskResponse, error) {
	task := m.K.legacyKeeper.GetTask(msg.TaskID)
	if task == nil {
		return nil, types.ErrTaskNotFound
	}

	// Verify ownership
	if task.RequesterID != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	// Only allow updates for draft or open tasks
	if task.Status != types.TaskStatusDraft && task.Status != types.TaskStatusOpen {
		return nil, fmt.Errorf("cannot update task in status: %s", task.Status)
	}

	if msg.Title != "" {
		task.Title = msg.Title
	}
	if msg.Description != "" {
		task.Description = msg.Description
	}
	if msg.Budget > 0 {
		task.Budget = msg.Budget
	}
	if msg.Deadline > 0 {
		task.Deadline = msg.Deadline
	}
	if len(msg.Skills) > 0 {
		task.Skills = msg.Skills
	}
	if len(msg.Subtasks) > 0 {
		task.Subtasks = msg.Subtasks
	}
	if len(msg.Milestones) > 0 {
		task.Milestones = msg.Milestones
	}

	if err := m.K.legacyKeeper.UpdateTask(task); err != nil {
		return nil, err
	}

	return &types.MsgUpdateTaskResponse{}, nil
}

// PublishTask handles MsgPublishTask
func (m msgServer) PublishTask(ctx sdk.Context, msg *types.MsgPublishTask) (*types.MsgPublishTaskResponse, error) {
	task := m.K.legacyKeeper.GetTask(msg.TaskID)
	if task == nil {
		return nil, types.ErrTaskNotFound
	}

	// Verify ownership
	if task.RequesterID != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if task.Status != types.TaskStatusDraft {
		return nil, fmt.Errorf("task must be in draft status to publish")
	}

	task.Publish()

	// Create auction if auction type
	if task.Type == types.TaskTypeAuction {
		_, err := m.K.legacyKeeper.CreateAuction(task.ID, task.Budget, task.Budget/2, 7*24*60*60) // 7 days
		if err != nil {
			return nil, err
		}
	}

	if err := m.K.legacyKeeper.UpdateTask(task); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePublishTask,
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskID),
		),
	)

	return &types.MsgPublishTaskResponse{}, nil
}

// CancelTask handles MsgCancelTask
func (m msgServer) CancelTask(ctx sdk.Context, msg *types.MsgCancelTask) (*types.MsgCancelTaskResponse, error) {
	task := m.K.legacyKeeper.GetTask(msg.TaskID)
	if task == nil {
		return nil, types.ErrTaskNotFound
	}

	// Verify ownership
	if task.RequesterID != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if task.Status == types.TaskStatusCompleted {
		return nil, fmt.Errorf("cannot cancel completed task")
	}

	task.Cancel()

	if err := m.K.legacyKeeper.UpdateTask(task); err != nil {
		return nil, err
	}

	return &types.MsgCancelTaskResponse{}, nil
}

// SubmitApplication handles MsgSubmitApplication
func (m msgServer) SubmitApplication(ctx sdk.Context, msg *types.MsgSubmitApplication) (*types.MsgSubmitApplicationResponse, error) {
	appID := fmt.Sprintf("app-%d", ctx.BlockHeight())

	app := types.NewApplication(
		appID,
		msg.TaskID,
		msg.WorkerID,
		msg.ProposedPrice,
	)

	app.CoverLetter = msg.CoverLetter
	app.RelevantExperience = msg.RelevantExperience
	app.PortfolioLinks = msg.PortfolioLinks
	app.EstimatedDuration = msg.EstimatedDuration

	if err := m.K.legacyKeeper.SubmitApplication(app); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitApplication,
			sdk.NewAttribute(types.AttributeKeyApplicationID, appID),
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskID),
			sdk.NewAttribute(types.AttributeKeyWorkerID, msg.WorkerID),
		),
	)

	return &types.MsgSubmitApplicationResponse{ApplicationID: appID}, nil
}

// AcceptApplication handles MsgAcceptApplication
func (m msgServer) AcceptApplication(ctx sdk.Context, msg *types.MsgAcceptApplication) (*types.MsgAcceptApplicationResponse, error) {
	if err := m.K.legacyKeeper.AcceptApplication(msg.ApplicationID); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAcceptApplication,
			sdk.NewAttribute(types.AttributeKeyApplicationID, msg.ApplicationID),
		),
	)

	return &types.MsgAcceptApplicationResponse{}, nil
}

// RejectApplication handles MsgRejectApplication
func (m msgServer) RejectApplication(ctx sdk.Context, msg *types.MsgRejectApplication) (*types.MsgRejectApplicationResponse, error) {
	if err := m.K.legacyKeeper.RejectApplication(msg.ApplicationID); err != nil {
		return nil, err
	}

	return &types.MsgRejectApplicationResponse{}, nil
}

// SubmitBid handles MsgSubmitBid
func (m msgServer) SubmitBid(ctx sdk.Context, msg *types.MsgSubmitBid) (*types.MsgSubmitBidResponse, error) {
	bidID := fmt.Sprintf("bid-%d", ctx.BlockHeight())

	bid := types.NewBid(
		bidID,
		msg.TaskID,
		msg.WorkerID,
		msg.Amount,
	)

	bid.Message = msg.Message
	bid.Portfolio = msg.Portfolio

	if err := m.K.legacyKeeper.SubmitBid(bid); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitBid,
			sdk.NewAttribute(types.AttributeKeyBidID, bidID),
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskID),
			sdk.NewAttribute(types.AttributeKeyWorkerID, msg.WorkerID),
		),
	)

	return &types.MsgSubmitBidResponse{BidID: bidID}, nil
}

// CloseAuction handles MsgCloseAuction
func (m msgServer) CloseAuction(ctx sdk.Context, msg *types.MsgCloseAuction) (*types.MsgCloseAuctionResponse, error) {
	if err := m.K.legacyKeeper.CloseAuction(msg.TaskID); err != nil {
		return nil, err
	}

	auction := m.K.legacyKeeper.GetAuction(msg.TaskID)
	var winnerID string
	if auction != nil && auction.WinningBidID != "" {
		for _, bid := range auction.Bids {
			if bid.ID == auction.WinningBidID {
				winnerID = bid.WorkerID
				break
			}
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCloseAuction,
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskID),
			sdk.NewAttribute(types.AttributeKeyWinnerID, winnerID),
		),
	)

	return &types.MsgCloseAuctionResponse{WinnerID: winnerID}, nil
}

// StartTask handles MsgStartTask
func (m msgServer) StartTask(ctx sdk.Context, msg *types.MsgStartTask) (*types.MsgStartTaskResponse, error) {
	if err := m.K.legacyKeeper.StartTask(msg.TaskID); err != nil {
		return nil, err
	}

	return &types.MsgStartTaskResponse{}, nil
}

// SubmitMilestone handles MsgSubmitMilestone
func (m msgServer) SubmitMilestone(ctx sdk.Context, msg *types.MsgSubmitMilestone) (*types.MsgSubmitMilestoneResponse, error) {
	if err := m.K.legacyKeeper.SubmitMilestone(msg.TaskID, msg.MilestoneID, msg.Deliverables); err != nil {
		return nil, err
	}

	return &types.MsgSubmitMilestoneResponse{}, nil
}

// ApproveMilestone handles MsgApproveMilestone
func (m msgServer) ApproveMilestone(ctx sdk.Context, msg *types.MsgApproveMilestone) (*types.MsgApproveMilestoneResponse, error) {
	if err := m.K.legacyKeeper.ApproveMilestone(msg.TaskID, msg.MilestoneID); err != nil {
		return nil, err
	}

	return &types.MsgApproveMilestoneResponse{}, nil
}

// RejectMilestone handles MsgRejectMilestone
func (m msgServer) RejectMilestone(ctx sdk.Context, msg *types.MsgRejectMilestone) (*types.MsgRejectMilestoneResponse, error) {
	if err := m.K.legacyKeeper.RejectMilestone(msg.TaskID, msg.MilestoneID, msg.Reason); err != nil {
		return nil, err
	}

	return &types.MsgRejectMilestoneResponse{}, nil
}

// SubmitRating handles MsgSubmitRating
func (m msgServer) SubmitRating(ctx sdk.Context, msg *types.MsgSubmitRating) (*types.MsgSubmitRatingResponse, error) {
	ratingID := fmt.Sprintf("rating-%d", ctx.BlockHeight())

	rating := types.NewRating(
		ratingID,
		msg.TaskID,
		msg.RaterID,
		msg.RatedID,
	)

	for dimStr, val := range msg.Ratings {
		dim := types.RatingDimension(dimStr)
		if err := rating.SetRating(dim, val); err != nil {
			return nil, err
		}
	}
	rating.Comment = msg.Comment

	if err := m.K.legacyKeeper.SubmitRating(rating); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitRating,
			sdk.NewAttribute(types.AttributeKeyRatingID, ratingID),
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskID),
			sdk.NewAttribute(types.AttributeKeyRatedID, msg.RatedID),
		),
	)

	return &types.MsgSubmitRatingResponse{RatingID: ratingID}, nil
}

// Unused context param versions for interface compatibility
var _ = context.Background()
