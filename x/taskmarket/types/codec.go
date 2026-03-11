package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Register messages
	cdc.RegisterConcrete(&MsgCreateTask{}, "taskmarket/CreateTask", nil)
	cdc.RegisterConcrete(&MsgUpdateTask{}, "taskmarket/UpdateTask", nil)
	cdc.RegisterConcrete(&MsgPublishTask{}, "taskmarket/PublishTask", nil)
	cdc.RegisterConcrete(&MsgCancelTask{}, "taskmarket/CancelTask", nil)
	cdc.RegisterConcrete(&MsgSubmitApplication{}, "taskmarket/SubmitApplication", nil)
	cdc.RegisterConcrete(&MsgAcceptApplication{}, "taskmarket/AcceptApplication", nil)
	cdc.RegisterConcrete(&MsgRejectApplication{}, "taskmarket/RejectApplication", nil)
	cdc.RegisterConcrete(&MsgSubmitBid{}, "taskmarket/SubmitBid", nil)
	cdc.RegisterConcrete(&MsgCloseAuction{}, "taskmarket/CloseAuction", nil)
	cdc.RegisterConcrete(&MsgStartTask{}, "taskmarket/StartTask", nil)
	cdc.RegisterConcrete(&MsgSubmitMilestone{}, "taskmarket/SubmitMilestone", nil)
	cdc.RegisterConcrete(&MsgApproveMilestone{}, "taskmarket/ApproveMilestone", nil)
	cdc.RegisterConcrete(&MsgRejectMilestone{}, "taskmarket/RejectMilestone", nil)
	cdc.RegisterConcrete(&MsgSubmitRating{}, "taskmarket/SubmitRating", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// TODO: Register when proto files are generated
	// For now, skip registration
	_ = registry
}

var (
	Amino = codec.NewLegacyAmino()
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}
