package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func init() {
	RegisterCodec(Amino)
}

// RegisterCodec registers the necessary x/identity interfaces and concrete types
// on the provided LegacyAmino codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterIdentity{}, "identity/RegisterIdentity", nil)
	cdc.RegisterConcrete(&MsgVerifyIdentity{}, "identity/VerifyIdentity", nil)
	cdc.RegisterConcrete(&MsgUpdateLimitConfig{}, "identity/UpdateLimitConfig", nil)
	cdc.RegisterConcrete(&MsgResetDailyLimits{}, "identity/ResetDailyLimits", nil)
}

// RegisterInterfaces registers the x/identity interfaces with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// For now, we don't register proto interfaces since we don't have proto generation
	// This would be done by the generated code if we had proper proto generation
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

// MarshalJSON helper for codec
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MustMarshalJSON helper for codec
func MustMarshalJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
