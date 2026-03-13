package llmcustody

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"sharetoken/testutil/sample"
	"sharetoken/x/llmcustody/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = rand.Rand{}
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	llmcustodyGenesis := types.GenesisState{
		APIKeys:       []types.APIKey{},
		EncryptionKey: nil,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	// Use standard JSON marshaling since GenesisState is not a protobuf type
	genesisJSON, err := json.Marshal(&llmcustodyGenesis)
	if err != nil {
		panic(err)
	}
	simState.GenState[types.ModuleName] = genesisJSON
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
// nolint:staticcheck // WeightedProposalContent is deprecated but kept for compatibility
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
