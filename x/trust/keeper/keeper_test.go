package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/trust/types"
)

func TestNewMQScore(t *testing.T) {
	address := "sharetoken1xyz"
	mq := types.NewMQScore(address)

	require.Equal(t, address, mq.Address)
	require.Equal(t, int32(100), mq.Score)
	require.Equal(t, uint64(0), mq.Disputes)
	require.Equal(t, uint64(0), mq.Consensus)
}

func TestMQValidation(t *testing.T) {
	mq := types.NewMQScore("test")
	require.NoError(t, mq.ValidateBasic())

	// Invalid score
	mq.Score = 150
	require.Error(t, mq.ValidateBasic())

	mq.Score = -10
	require.Error(t, mq.ValidateBasic())

	// Missing address
	mq.Address = ""
	require.Error(t, mq.ValidateBasic())
}

func TestCalculateVoteWeight(t *testing.T) {
	mq := types.NewMQScore("test")

	// Initial score is 100, weight should be 1.0
	weight := mq.CalculateVoteWeight()
	require.True(t, weight.Equal(sdk.NewDec(1)))

	// Score 50, weight should be 0.5
	mq.Score = 50
	weight = mq.CalculateVoteWeight()
	require.True(t, weight.Equal(sdk.NewDecWithPrec(5, 1)))
}

func TestCalculatePenalty(t *testing.T) {
	// High MQ user (100) - should get max penalty
	highMQ := types.NewMQScore("high")
	highMQ.Score = 100
	penalty := highMQ.CalculatePenalty(false)
	require.Equal(t, int32(3), penalty) // Max 3%

	// Medium MQ user (60) - should get medium penalty
	medMQ := types.NewMQScore("med")
	medMQ.Score = 60
	penalty = medMQ.CalculatePenalty(false)
	require.Equal(t, int32(2), penalty)

	// Low MQ user (30) - should get min penalty
	lowMQ := types.NewMQScore("low")
	lowMQ.Score = 30
	penalty = lowMQ.CalculatePenalty(false)
	require.Equal(t, int32(1), penalty)

	// Consensus vote - no penalty
	penalty = highMQ.CalculatePenalty(true)
	require.Equal(t, int32(0), penalty)
}

func TestCalculateReward(t *testing.T) {
	mq := types.NewMQScore("test")

	// Normal reward
	reward := mq.CalculateReward(10, true)
	require.Equal(t, int32(1), reward)

	// Low MQ user gets higher reward
	mq.Score = 20
	reward = mq.CalculateReward(10, true)
	require.Equal(t, int32(2), reward)

	// No reward for voting against consensus
	reward = mq.CalculateReward(10, false)
	require.Equal(t, int32(0), reward)
}

func TestApplyPenaltyAndReward(t *testing.T) {
	mq := types.NewMQScore("test")
	mq.Score = 50

	// Apply penalty
	mq.ApplyPenalty(5)
	require.Equal(t, int32(45), mq.Score)

	// Apply reward
	mq.ApplyReward(10)
	require.Equal(t, int32(55), mq.Score)

	// Cap at 100
	mq.Score = 95
	mq.ApplyReward(10)
	require.Equal(t, int32(100), mq.Score)

	// Floor at 0
	mq.Score = 5
	mq.ApplyPenalty(10)
	require.Equal(t, int32(0), mq.Score)
}

func TestRecordDispute(t *testing.T) {
	mq := types.NewMQScore("test")

	mq.RecordDispute(true)
	require.Equal(t, uint64(1), mq.Disputes)
	require.Equal(t, uint64(1), mq.Consensus)

	mq.RecordDispute(false)
	require.Equal(t, uint64(2), mq.Disputes)
	require.Equal(t, uint64(1), mq.Consensus)
}

func TestConsensusRate(t *testing.T) {
	mq := types.NewMQScore("test")

	// No disputes yet
	rate := mq.ConsensusRate()
	require.True(t, rate.IsZero())

	// 3 disputes, 2 with consensus
	mq.Disputes = 3
	mq.Consensus = 2
	rate = mq.ConsensusRate()
	// 2/3 = 0.666666...
	require.True(t, rate.GT(sdk.NewDecWithPrec(66, 2)) && rate.LT(sdk.NewDecWithPrec(67, 2)))
}
