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
	cdc.RegisterConcrete(&MsgRotateAPIKey{}, "llmcustody/RotateAPIKey", nil)
}

// MsgServer registration helpers
func RegisterMsgServer(server interface{}, srv MsgServer) {
	// This is a placeholder for when proto-generated code is available
	// The actual implementation would call the generated RegisterMsgServer function
}

// RegisterQueryServer registers the query server
func RegisterQueryServer(server interface{}, srv QueryServer) {
	// This is a placeholder for when proto-generated code is available
	// The actual implementation would call the generated RegisterQueryServer function
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
