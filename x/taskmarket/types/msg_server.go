package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgServer is the message server interface for taskmarket module
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
