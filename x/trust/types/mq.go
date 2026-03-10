package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MQScore represents a user's Moral Quotient score
type MQScore struct {
	Address   string `json:"address"`
	Score     int32  `json:"score"`     // Current MQ score (0-100)
	Disputes  uint64 `json:"disputes"`  // Number of disputes participated
	Consensus uint64 `json:"consensus"` // Number of times voted with consensus
	UpdatedAt int64  `json:"updated_at"`
}

// InitialMQ is the starting MQ score for new users
const InitialMQ int32 = 100

// MaxMQLossPercent is the maximum percentage of MQ that can be lost in a single dispute
const MaxMQLossPercent = 3

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
	if mq.Score < 0 || mq.Score > 100 {
		return ErrInvalidMQ.Wrapf("score must be between 0 and 100, got %d", mq.Score)
	}
	return nil
}

// CalculateVoteWeight calculates voting weight based on MQ score
// Higher MQ = higher weight
func (mq MQScore) CalculateVoteWeight() sdk.Dec {
	// Weight = MQ / 100
	return sdk.NewDec(int64(mq.Score)).Quo(sdk.NewDec(100))
}

// CalculatePenalty calculates MQ penalty for voting against consensus
// Higher MQ users lose more when they vote against consensus (convergence mechanism)
func (mq MQScore) CalculatePenalty(isConsensus bool) int32 {
	if isConsensus {
		return 0 // No penalty for voting with consensus
	}

	// Base penalty: 1%
	penalty := int32(1)

	// Additional penalty based on MQ (higher MQ = more penalty)
	// This is the convergence mechanism - high MQ users lose more for being wrong
	if mq.Score > 80 {
		penalty += 2 // +2% for high MQ users
	} else if mq.Score > 50 {
		penalty += 1 // +1% for medium MQ users
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

	// Base reward: 0.5%
	reward := int32(1)

	// Higher reward for low MQ users to help them recover
	if mq.Score < 30 {
		reward += 1 // +1% for low MQ users
	}

	return reward
}

// ApplyPenalty applies penalty to MQ score
func (mq *MQScore) ApplyPenalty(penalty int32) {
	mq.Score -= penalty
	if mq.Score < 0 {
		mq.Score = 0
	}
}

// ApplyReward applies reward to MQ score
func (mq *MQScore) ApplyReward(reward int32) {
	mq.Score += reward
	if mq.Score > 100 {
		mq.Score = 100
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
