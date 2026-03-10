package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Register messages
	cdc.RegisterConcrete(&MsgRegisterAPIKey{}, "llmcustody/RegisterAPIKey", nil)
	cdc.RegisterConcrete(&MsgUpdateAPIKey{}, "llmcustody/UpdateAPIKey", nil)
	cdc.RegisterConcrete(&MsgRevokeAPIKey{}, "llmcustody/RevokeAPIKey", nil)
	cdc.RegisterConcrete(&MsgRecordUsage{}, "llmcustody/RecordUsage", nil)
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
