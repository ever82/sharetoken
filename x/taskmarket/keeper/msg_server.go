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
func (m msgServer) CreateTask(ctx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Generate task ID (in production, use proper ID generation)
	taskID := fmt.Sprintf("task-%d", sdkCtx.BlockHeight())

	task := types.NewTask(
		taskID,
		msg.Title,
		msg.Description,
		msg.Creator,
		msg.TaskType,
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

	if err := m.K.CreateTask(sdkCtx, *task); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTask,
			sdk.NewAttribute(types.AttributeKeyTaskID, taskID),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyTaskType, msg.TaskType.String()),
		),
	)

	return &types.MsgCreateTaskResponse{TaskId: taskID}, nil
}

// UpdateTask handles MsgUpdateTask
func (m msgServer) UpdateTask(ctx context.Context, msg *types.MsgUpdateTask) (*types.MsgUpdateTaskResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	task, found := m.K.GetTask(sdkCtx, msg.TaskId)
	if !found {
		return nil, types.ErrTaskNotFound
	}

	// Verify ownership
	if task.RequesterId != msg.Creator {
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

	if err := m.K.UpdateTask(sdkCtx, task); err != nil {
		return nil, err
	}

	return &types.MsgUpdateTaskResponse{}, nil
}

// PublishTask handles MsgPublishTask
func (m msgServer) PublishTask(ctx context.Context, msg *types.MsgPublishTask) (*types.MsgPublishTaskResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	task, found := m.K.GetTask(sdkCtx, msg.TaskId)
	if !found {
		return nil, types.ErrTaskNotFound
	}

	// Verify ownership
	if task.RequesterId != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if task.Status != types.TaskStatusDraft {
		return nil, fmt.Errorf("task must be in draft status to publish")
	}

	task.Publish()

	// Create auction if auction type
	if task.TaskType == types.TaskTypeAuction {
		_, err := m.K.CreateAuction(sdkCtx, task.Id, task.Budget, task.Budget/2, 7*24*60*60) // 7 days
		if err != nil {
			return nil, err
		}
	}

	if err := m.K.UpdateTask(sdkCtx, task); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePublishTask,
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskId),
		),
	)

	return &types.MsgPublishTaskResponse{}, nil
}

// CancelTask handles MsgCancelTask
func (m msgServer) CancelTask(ctx context.Context, msg *types.MsgCancelTask) (*types.MsgCancelTaskResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	task, found := m.K.GetTask(sdkCtx, msg.TaskId)
	if !found {
		return nil, types.ErrTaskNotFound
	}

	// Verify ownership
	if task.RequesterId != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if task.Status == types.TaskStatusCompleted {
		return nil, fmt.Errorf("cannot cancel completed task")
	}

	task.Cancel()

	if err := m.K.UpdateTask(sdkCtx, task); err != nil {
		return nil, err
	}

	return &types.MsgCancelTaskResponse{}, nil
}

// SubmitApplication handles MsgSubmitApplication
func (m msgServer) SubmitApplication(ctx context.Context, msg *types.MsgSubmitApplication) (*types.MsgSubmitApplicationResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	appID := fmt.Sprintf("app-%d", sdkCtx.BlockHeight())

	app := types.NewApplication(
		appID,
		msg.TaskId,
		msg.WorkerId,
		msg.ProposedPrice,
	)

	app.CoverLetter = msg.CoverLetter
	app.RelevantExperience = msg.RelevantExperience
	app.PortfolioLinks = msg.PortfolioLinks
	app.EstimatedDuration = msg.EstimatedDuration

	if err := m.K.SubmitApplication(sdkCtx, *app); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitApplication,
			sdk.NewAttribute(types.AttributeKeyApplicationID, appID),
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskId),
			sdk.NewAttribute(types.AttributeKeyWorkerID, msg.WorkerId),
		),
	)

	return &types.MsgSubmitApplicationResponse{ApplicationId: appID}, nil
}

// AcceptApplication handles MsgAcceptApplication
func (m msgServer) AcceptApplication(ctx context.Context, msg *types.MsgAcceptApplication) (*types.MsgAcceptApplicationResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.AcceptApplication(sdkCtx, msg.ApplicationId); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAcceptApplication,
			sdk.NewAttribute(types.AttributeKeyApplicationID, msg.ApplicationId),
		),
	)

	return &types.MsgAcceptApplicationResponse{}, nil
}

// RejectApplication handles MsgRejectApplication
func (m msgServer) RejectApplication(ctx context.Context, msg *types.MsgRejectApplication) (*types.MsgRejectApplicationResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.RejectApplication(sdkCtx, msg.ApplicationId); err != nil {
		return nil, err
	}

	return &types.MsgRejectApplicationResponse{}, nil
}

// SubmitBid handles MsgSubmitBid
func (m msgServer) SubmitBid(ctx context.Context, msg *types.MsgSubmitBid) (*types.MsgSubmitBidResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	bidID := fmt.Sprintf("bid-%d", sdkCtx.BlockHeight())

	bid := types.NewBid(
		bidID,
		msg.TaskId,
		msg.WorkerId,
		msg.Amount,
	)

	bid.Message = msg.Message
	bid.Portfolio = msg.Portfolio

	if err := m.K.SubmitBid(sdkCtx, *bid); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitBid,
			sdk.NewAttribute(types.AttributeKeyBidID, bidID),
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskId),
			sdk.NewAttribute(types.AttributeKeyWorkerID, msg.WorkerId),
		),
	)

	return &types.MsgSubmitBidResponse{BidId: bidID}, nil
}

// CloseAuction handles MsgCloseAuction
func (m msgServer) CloseAuction(ctx context.Context, msg *types.MsgCloseAuction) (*types.MsgCloseAuctionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.CloseAuction(sdkCtx, msg.TaskId); err != nil {
		return nil, err
	}

	auction, found := m.K.GetAuction(sdkCtx, msg.TaskId)
	var winnerID string
	if found && auction.WinningBidId != "" {
		for _, bid := range auction.Bids {
			if bid.Id == auction.WinningBidId {
				winnerID = bid.WorkerId
				break
			}
		}
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCloseAuction,
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskId),
			sdk.NewAttribute(types.AttributeKeyWinnerID, winnerID),
		),
	)

	return &types.MsgCloseAuctionResponse{WinnerId: winnerID}, nil
}

// StartTask handles MsgStartTask
func (m msgServer) StartTask(ctx context.Context, msg *types.MsgStartTask) (*types.MsgStartTaskResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.StartTask(sdkCtx, msg.TaskId); err != nil {
		return nil, err
	}

	return &types.MsgStartTaskResponse{}, nil
}

// SubmitMilestone handles MsgSubmitMilestone
func (m msgServer) SubmitMilestone(ctx context.Context, msg *types.MsgSubmitMilestone) (*types.MsgSubmitMilestoneResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.SubmitMilestone(sdkCtx, msg.TaskId, msg.MilestoneId, msg.Deliverables); err != nil {
		return nil, err
	}

	return &types.MsgSubmitMilestoneResponse{}, nil
}

// ApproveMilestone handles MsgApproveMilestone
func (m msgServer) ApproveMilestone(ctx context.Context, msg *types.MsgApproveMilestone) (*types.MsgApproveMilestoneResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.ApproveMilestone(sdkCtx, msg.TaskId, msg.MilestoneId); err != nil {
		return nil, err
	}

	return &types.MsgApproveMilestoneResponse{}, nil
}

// RejectMilestone handles MsgRejectMilestone
func (m msgServer) RejectMilestone(ctx context.Context, msg *types.MsgRejectMilestone) (*types.MsgRejectMilestoneResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := m.K.RejectMilestone(sdkCtx, msg.TaskId, msg.MilestoneId, msg.Reason); err != nil {
		return nil, err
	}

	return &types.MsgRejectMilestoneResponse{}, nil
}

// SubmitRating handles MsgSubmitRating
func (m msgServer) SubmitRating(ctx context.Context, msg *types.MsgSubmitRating) (*types.MsgSubmitRatingResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	ratingID := fmt.Sprintf("rating-%d", sdkCtx.BlockHeight())

	rating := types.NewRating(
		ratingID,
		msg.TaskId,
		msg.RaterId,
		msg.RatedId,
	)

	for dimStr, val := range msg.Ratings {
		if err := rating.SetRating(dimStr, val); err != nil {
			return nil, err
		}
	}
	rating.Comment = msg.Comment

	if err := m.K.SubmitRating(sdkCtx, *rating); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitRating,
			sdk.NewAttribute(types.AttributeKeyRatingID, ratingID),
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskId),
			sdk.NewAttribute(types.AttributeKeyRatedID, msg.RatedId),
		),
	)

	return &types.MsgSubmitRatingResponse{RatingId: ratingID}, nil
}
