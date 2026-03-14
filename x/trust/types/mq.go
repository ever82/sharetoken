package types

import (
	"time"
)

// InitialMQ is the initial MQ score for new users
const InitialMQ = 100

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
	// Simple reward calculation
	return int32(5)
}

// CalculatePenalty calculates the penalty for incorrect voting
func (m *MQScore) CalculatePenalty(votedWithConsensus bool) int32 {
	if votedWithConsensus {
		return 0
	}
	// Simple penalty calculation
	return int32(10)
}

// ApplyReward applies a reward to the score
func (m *MQScore) ApplyReward(reward int32) {
	m.Score += reward
	if m.Score > 1000 {
		m.Score = 1000 // Cap at 1000
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
