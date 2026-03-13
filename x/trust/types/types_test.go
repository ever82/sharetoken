package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"sharetoken/x/trust/types"
)

func TestNewMQScore(t *testing.T) {
	mq := types.NewMQScore("address-1")

	require.NotNil(t, mq)
	require.Equal(t, "address-1", mq.Address)
	require.Equal(t, types.InitialMQ, mq.Score)
	require.Equal(t, uint64(0), mq.Disputes)
	require.Equal(t, uint64(0), mq.Consensus)
	require.Equal(t, int64(0), mq.UpdatedAt)
}

func TestMQScore_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		mq      types.MQScore
		wantErr bool
		errType error
	}{
		{
			name: "valid MQ score",
			mq: types.MQScore{
				Address: "address-1",
				Score:   100,
			},
			wantErr: false,
		},
		{
			name: "valid - minimum score",
			mq: types.MQScore{
				Address: "address-1",
				Score:   0,
			},
			wantErr: false,
		},
		{
			name: "valid - maximum score",
			mq: types.MQScore{
				Address: "address-1",
				Score:   100,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty address",
			mq: types.MQScore{
				Address: "",
				Score:   100,
			},
			wantErr: true,
			errType: types.ErrInvalidMQ,
		},
		{
			name: "invalid - negative score",
			mq: types.MQScore{
				Address: "address-1",
				Score:   -1,
			},
			wantErr: true,
			errType: types.ErrInvalidMQ,
		},
		{
			name: "invalid - score over 100",
			mq: types.MQScore{
				Address: "address-1",
				Score:   101,
			},
			wantErr: true,
			errType: types.ErrInvalidMQ,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mq.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					require.ErrorIs(t, err, tt.errType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMQScore_CalculateVoteWeight(t *testing.T) {
	tests := []struct {
		name     string
		score    int32
		expected sdk.Dec
	}{
		{"score 100", 100, sdk.NewDec(1)},     // 100/100 = 1
		{"score 50", 50, sdk.NewDecWithPrec(5, 1)},   // 50/100 = 0.5
		{"score 0", 0, sdk.NewDec(0)},         // 0/100 = 0
		{"score 75", 75, sdk.NewDecWithPrec(75, 2)}, // 75/100 = 0.75
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := types.MQScore{Score: tt.score}
			weight := mq.CalculateVoteWeight()
			require.True(t, weight.Equal(tt.expected))
		})
	}
}

func TestMQScore_CalculatePenalty(t *testing.T) {
	tests := []struct {
		name       string
		score      int32
		isConsensus bool
		expected   int32
	}{
		{"consensus - no penalty", 100, true, 0},
		{"low MQ penalty", 30, false, 1},
		{"medium MQ penalty", 60, false, 2},
		{"high MQ penalty", 85, false, 3},
		{"max MQ capped", 100, false, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := types.MQScore{Score: tt.score}
			penalty := mq.CalculatePenalty(tt.isConsensus)
			require.Equal(t, tt.expected, penalty)
		})
	}
}

func TestMQScore_CalculateReward(t *testing.T) {
	tests := []struct {
		name        string
		score       int32
		votersCount int64
		isConsensus bool
		expected    int32
	}{
		{"not consensus - no reward", 100, 10, false, 0},
		{"consensus - base reward", 50, 10, true, 1},
		{"consensus - low MQ bonus", 25, 10, true, 2},
		{"consensus - very low MQ", 20, 10, true, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := types.MQScore{Score: tt.score}
			reward := mq.CalculateReward(tt.votersCount, tt.isConsensus)
			require.Equal(t, tt.expected, reward)
		})
	}
}

func TestMQScore_ApplyPenalty(t *testing.T) {
	tests := []struct {
		name           string
		initialScore   int32
		penalty        int32
		expectedScore  int32
	}{
		{"normal penalty", 100, 10, 90},
		{"penalty to zero", 50, 50, 0},
		{"penalty exceeds score", 50, 100, 0},
		{"no penalty", 80, 0, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := types.MQScore{Score: tt.initialScore}
			mq.ApplyPenalty(tt.penalty)
			require.Equal(t, tt.expectedScore, mq.Score)
		})
	}
}

func TestMQScore_ApplyReward(t *testing.T) {
	tests := []struct {
		name          string
		initialScore  int32
		reward        int32
		expectedScore int32
	}{
		{"normal reward", 50, 10, 60},
		{"reward to max", 95, 10, 100},
		{"reward exceeds max", 100, 10, 100},
		{"no reward", 80, 0, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := types.MQScore{Score: tt.initialScore}
			mq.ApplyReward(tt.reward)
			require.Equal(t, tt.expectedScore, mq.Score)
		})
	}
}

func TestMQScore_RecordDispute(t *testing.T) {
	mq := types.NewMQScore("address-1")

	require.Equal(t, uint64(0), mq.Disputes)
	require.Equal(t, uint64(0), mq.Consensus)

	mq.RecordDispute(true)
	require.Equal(t, uint64(1), mq.Disputes)
	require.Equal(t, uint64(1), mq.Consensus)

	mq.RecordDispute(false)
	require.Equal(t, uint64(2), mq.Disputes)
	require.Equal(t, uint64(1), mq.Consensus)
}

func TestMQScore_ConsensusRate(t *testing.T) {
	tests := []struct {
		name     string
		disputes uint64
		consensus uint64
		expected sdk.Dec
	}{
		{"no disputes", 0, 0, sdk.NewDec(0)},
		{"100% consensus", 10, 10, sdk.NewDec(1)},
		{"50% consensus", 10, 5, sdk.NewDecWithPrec(5, 1)},
		{"0% consensus", 10, 0, sdk.NewDec(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := types.MQScore{
				Disputes:  tt.disputes,
				Consensus: tt.consensus,
			}
			rate := mq.ConsensusRate()
			require.True(t, rate.Equal(tt.expected))
		})
	}
}

func TestMQScore_String(t *testing.T) {
	mq := types.NewMQScore("address-1")
	mq.Disputes = 10
	mq.Consensus = 8

	result := mq.String()
	require.Contains(t, result, "address-1")
	require.Contains(t, result, "100")
	require.Contains(t, result, "10")
}

// Genesis Tests

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.MQScores)
	require.Empty(t, genesis.MQScores)
}

func TestValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		data    types.GenesisState
		wantErr bool
	}{
		{
			name:    "valid genesis with default",
			data:    *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "valid genesis with MQ scores",
			data: types.GenesisState{
				MQScores: []types.MQScore{
					{Address: "address-1", Score: 100},
					{Address: "address-2", Score: 80},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate addresses",
			data: types.GenesisState{
				MQScores: []types.MQScore{
					{Address: "address-1", Score: 100},
					{Address: "address-1", Score: 80},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid MQ score",
			data: types.GenesisState{
				MQScores: []types.MQScore{
					{Address: "", Score: 100},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - score out of range",
			data: types.GenesisState{
				MQScores: []types.MQScore{
					{Address: "address-1", Score: 150},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateGenesis(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
