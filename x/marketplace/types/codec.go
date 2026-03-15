package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}

// RegisterCodec registers the necessary x/marketplace interfaces and concrete types
// on the provided LegacyAmino codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	// Legacy Amino registration if needed
}

// RegisterInterfaces registers the x/marketplace interfaces with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Register proto message implementations
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterService{},
		&MsgUpdateService{},
		&MsgActivateService{},
		&MsgDeactivateService{},
		&MsgPurchaseService{},
	)

	// Register the message service descriptor
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
