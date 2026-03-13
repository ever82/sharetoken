package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MQScore represents a user's Moral Quotient score
type MQScore struct {
	Address   string `json:"address"`
	Score     int32  `json:"score"`     // Current MQ score (MinMQScore-MaxMQScore)
	Disputes  uint64 `json:"disputes"`  // Number of disputes participated
	Consensus uint64 `json:"consensus"` // Number of times voted with consensus
	UpdatedAt int64  `json:"updated_at"`
}

// MQ score constants
const (
	// MinMQScore is the minimum possible MQ score
	MinMQScore int32 = 0
	// MaxMQScore is the maximum possible MQ score
	MaxMQScore int32 = 100
	// InitialMQ is the starting MQ score for new users
	InitialMQ int32 = 100

	// MQ score thresholds for penalty calculation
	// HighMQThreshold is the threshold for high MQ users (affected more by penalties)
	HighMQThreshold int32 = 80
	// MediumMQThreshold is the threshold for medium MQ users
	MediumMQThreshold int32 = 50
	// LowMQThreshold is the threshold for low MQ users (receive more rewards)
	LowMQThreshold int32 = 30

	// Penalty and reward constants
	// BasePenalty is the base MQ penalty for voting against consensus
	BasePenalty int32 = 1
	// HighMQAdditionalPenalty is additional penalty for high MQ users
	HighMQAdditionalPenalty int32 = 2
	// MediumMQAdditionalPenalty is additional penalty for medium MQ users
	MediumMQAdditionalPenalty int32 = 1
	// MaxMQLossPercent is the maximum percentage of MQ that can be lost in a single dispute
	MaxMQLossPercent int32 = 3

	// BaseReward is the base MQ reward for voting with consensus
	BaseReward int32 = 1
	// LowMQAdditionalReward is additional reward for low MQ users to help recovery
	LowMQAdditionalReward int32 = 1

	// VoteWeightDivisor is the divisor for calculating vote weight from MQ score
	VoteWeightDivisor int64 = 100
)

// NewMQScore creates a new MQ score for a user
func NewMQScore(address string) *MQScore {
	return &MQScore{
		Address:   address,
		Score:     InitialMQ,
		Disputes:  0,
		Consensus: 0,
		UpdatedAt: 0,
	}
}

// ValidateBasic performs basic validation
func (mq MQScore) ValidateBasic() error {
	if mq.Address == "" {
		return ErrInvalidMQ.Wrap("address required")
	}
	if mq.Score < MinMQScore || mq.Score > MaxMQScore {
		return ErrInvalidMQ.Wrapf("score must be between %d and %d, got %d", MinMQScore, MaxMQScore, mq.Score)
	}
	return nil
}

// CalculateVoteWeight calculates voting weight based on MQ score
// Higher MQ = higher weight
func (mq MQScore) CalculateVoteWeight() sdk.Dec {
	// Weight = MQ / VoteWeightDivisor
	return sdk.NewDec(int64(mq.Score)).Quo(sdk.NewDec(VoteWeightDivisor))
}

// CalculatePenalty calculates MQ penalty for voting against consensus
// Higher MQ users lose more when they vote against consensus (convergence mechanism)
func (mq MQScore) CalculatePenalty(isConsensus bool) int32 {
	if isConsensus {
		return 0 // No penalty for voting with consensus
	}

	// Base penalty
	penalty := BasePenalty

	// Additional penalty based on MQ (higher MQ = more penalty)
	// This is the convergence mechanism - high MQ users lose more for being wrong
	if mq.Score > HighMQThreshold {
		penalty += HighMQAdditionalPenalty
	} else if mq.Score > MediumMQThreshold {
		penalty += MediumMQAdditionalPenalty
	}

	// Cap at MaxMQLossPercent
	if penalty > MaxMQLossPercent {
		penalty = MaxMQLossPercent
	}

	return penalty
}

// CalculateReward calculates MQ reward for voting with consensus
func (mq MQScore) CalculateReward(votersCount int64, isConsensus bool) int32 {
	if !isConsensus {
		return 0 // No reward for voting against consensus
	}

	// Base reward
	reward := BaseReward

	// Higher reward for low MQ users to help them recover
	if mq.Score < LowMQThreshold {
		reward += LowMQAdditionalReward
	}

	return reward
}

// ApplyPenalty applies penalty to MQ score
func (mq *MQScore) ApplyPenalty(penalty int32) {
	mq.Score -= penalty
	if mq.Score < MinMQScore {
		mq.Score = MinMQScore
	}
}

// ApplyReward applies reward to MQ score
func (mq *MQScore) ApplyReward(reward int32) {
	mq.Score += reward
	if mq.Score > MaxMQScore {
		mq.Score = MaxMQScore
	}
}

// RecordDispute records a dispute participation
func (mq *MQScore) RecordDispute(votedWithConsensus bool) {
	mq.Disputes++
	if votedWithConsensus {
		mq.Consensus++
	}
}

// ConsensusRate returns the rate of voting with consensus
func (mq MQScore) ConsensusRate() sdk.Dec {
	if mq.Disputes == 0 {
		return sdk.NewDec(0)
	}
	return sdk.NewDec(int64(mq.Consensus)).Quo(sdk.NewDec(int64(mq.Disputes)))
}

// String implements stringer
func (mq MQScore) String() string {
	consensusRate := mq.ConsensusRate()
	return fmt.Sprintf("MQScore{%s: %d, disputes: %d, consensus rate: %s}",
		mq.Address, mq.Score, mq.Disputes, consensusRate.String())
}

// GenesisState represents the genesis state
type GenesisState struct {
	MQScores []MQScore `json:"mq_scores"`
}

// DefaultGenesis returns default genesis
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		MQScores: []MQScore{},
	}
}

// ValidateGenesis validates genesis
func ValidateGenesis(data GenesisState) error {
	seen := make(map[string]bool)
	for _, mq := range data.MQScores {
		if seen[mq.Address] {
			return fmt.Errorf("duplicate MQ score for address: %s", mq.Address)
		}
		seen[mq.Address] = true

		if err := mq.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
