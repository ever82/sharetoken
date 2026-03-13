package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Register messages
	cdc.RegisterConcrete(&MsgRegisterAPIKey{}, "llmcustody/RegisterAPIKey", nil)
	cdc.RegisterConcrete(&MsgUpdateAPIKey{}, "llmcustody/UpdateAPIKey", nil)
	cdc.RegisterConcrete(&MsgRevokeAPIKey{}, "llmcustody/RevokeAPIKey", nil)
	cdc.RegisterConcrete(&MsgRecordUsage{}, "llmcustody/RecordUsage", nil)
	cdc.RegisterConcrete(&MsgRotateAPIKey{}, "llmcustody/RotateAPIKey", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Register messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterAPIKey{},
		&MsgUpdateAPIKey{},
		&MsgRevokeAPIKey{},
		&MsgRecordUsage{},
		&MsgRotateAPIKey{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino = codec.NewLegacyAmino()
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}
