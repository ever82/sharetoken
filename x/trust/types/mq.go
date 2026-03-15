package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitialMQ is the initial MQ score for new users
const InitialMQ = int32(100)

// NewMQScore creates a new MQ score with default values
func NewMQScore(address string) *MQScore {
	return &MQScore{
		Address:   address,
		Score:     100, // Default score
		Disputes:  0,
		Consensus: 0,
		UpdatedAt: time.Now().Unix(),
	}
}

// RecordDispute records dispute participation
func (m *MQScore) RecordDispute(votedWithConsensus bool) {
	m.Disputes++
	if votedWithConsensus {
		m.Consensus++
	}
	m.UpdatedAt = time.Now().Unix()
}

// CalculateReward calculates the reward for correct voting
func (m *MQScore) CalculateReward(votersCount int64, votedWithConsensus bool) int32 {
	if !votedWithConsensus {
		return 0
	}
	// Base reward is 1, low MQ gets bonus
	if m.Score <= 25 {
		return 2
	}
	return 1
}

// CalculatePenalty calculates the penalty for incorrect voting
func (m *MQScore) CalculatePenalty(votedWithConsensus bool) int32 {
	if votedWithConsensus {
		return 0
	}
	// Penalty based on MQ score
	if m.Score <= 50 {
		return 1
	}
	if m.Score <= 80 {
		return 2
	}
	return 3
}

// ApplyReward applies a reward to the score
func (m *MQScore) ApplyReward(reward int32) {
	m.Score += reward
	if m.Score > 100 {
		m.Score = 100 // Cap at 100
	}
	m.UpdatedAt = time.Now().Unix()
}

// ApplyPenalty applies a penalty to the score
func (m *MQScore) ApplyPenalty(penalty int32) {
	m.Score -= penalty
	if m.Score < 0 {
		m.Score = 0 // Floor at 0
	}
	m.UpdatedAt = time.Now().Unix()
}

// ValidateBasic validates the MQScore
func (m MQScore) ValidateBasic() error {
	if m.Address == "" {
		return ErrInvalidMQ.Wrap("address cannot be empty")
	}
	if m.Score < 0 || m.Score > 100 {
		return ErrInvalidMQ.Wrapf("score must be between 0 and 100, got %d", m.Score)
	}
	return nil
}

// CalculateVoteWeight calculates voting weight based on MQ score
func (m *MQScore) CalculateVoteWeight() sdk.Dec {
	return sdk.NewDec(int64(m.Score)).Quo(sdk.NewDec(100))
}

// ConsensusRate calculates the consensus participation rate
func (m *MQScore) ConsensusRate() sdk.Dec {
	if m.Disputes == 0 {
		return sdk.NewDec(0)
	}
	return sdk.NewDec(int64(m.Consensus)).Quo(sdk.NewDec(int64(m.Disputes)))
}
