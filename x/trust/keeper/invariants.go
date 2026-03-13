package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/trust/types"
)

// RegisterInvariants registers the trust module invariants
func (k *MQKeeper) RegisterInvariants(ir sdk.InvariantRegistry) {
	ir.RegisterRoute(types.ModuleName, "mq-score-validity",
		k.MQScoreValidityInvariant())
	ir.RegisterRoute(types.ModuleName, "mq-consensus-rate",
		k.MQConsensusRateInvariant())
	ir.RegisterRoute(types.ModuleName, "mq-participation-consistency",
		k.MQParticipationConsistencyInvariant())
}

// MQScoreValidityInvariant checks that all MQ scores are within valid range (0-100)
func (k *MQKeeper) MQScoreValidityInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidScores []string

		scores := k.GetAllScores()
		for addr, score := range scores {
			if score.Score < 0 || score.Score > 100 {
				invalidScores = append(invalidScores, fmt.Sprintf("%s:%d", addr, score.Score))
			}
		}

		if len(invalidScores) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"mq-score-validity",
				fmt.Sprintf("found %d MQ scores outside valid range (0-100): %v", len(invalidScores), invalidScores),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"mq-score-validity",
			"all MQ scores are within valid range",
		), false
	}
}

// MQConsensusRateInvariant checks that consensus rates are valid
// Consensus rate = consensus votes / total disputes
func (k *MQKeeper) MQConsensusRateInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidRates []string

		scores := k.GetAllScores()
		for addr, score := range scores {
			if score.Disputes > 0 {
				// Consensus should not exceed disputes
				if score.Consensus > score.Disputes {
					invalidRates = append(invalidRates, fmt.Sprintf("%s:consensus(%d)>disputes(%d)", addr, score.Consensus, score.Disputes))
				}

				// Calculate expected consensus rate
				expectedRate := float64(score.Consensus) / float64(score.Disputes)
				if expectedRate < 0 || expectedRate > 1 {
					invalidRates = append(invalidRates, fmt.Sprintf("%s:rate(%f)", addr, expectedRate))
				}
			}
		}

		if len(invalidRates) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"mq-consensus-rate",
				fmt.Sprintf("found %d invalid consensus rates: %v", len(invalidRates), invalidRates),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"mq-consensus-rate",
			"all consensus rates are valid",
		), false
	}
}

// MQParticipationConsistencyInvariant checks that participation history is consistent with scores
func (k *MQKeeper) MQParticipationConsistencyInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var inconsistencies []string

		scores := k.GetAllScores()
		for addr, score := range scores {
			// Get participation history
			participation := k.participation[addr]

			// If there are disputes, there should be participation history
			if score.Disputes > 0 && len(participation) == 0 {
				inconsistencies = append(inconsistencies, fmt.Sprintf("%s:disputes(%d) but no participation", addr, score.Disputes))
				continue
			}

			// Count actual participations
			actualParticipations := uint64(len(participation))
			if score.Disputes != actualParticipations {
				inconsistencies = append(inconsistencies,
					fmt.Sprintf("%s:disputes(%d)!=participations(%d)", addr, score.Disputes, actualParticipations))
			}

			// Count actual consensus votes
			var actualConsensus uint64
			for _, votedWithConsensus := range participation {
				if votedWithConsensus {
					actualConsensus++
				}
			}
			if score.Consensus != actualConsensus {
				inconsistencies = append(inconsistencies,
					fmt.Sprintf("%s:consensus(%d)!=actual(%d)", addr, score.Consensus, actualConsensus))
			}
		}

		if len(inconsistencies) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"mq-participation-consistency",
				fmt.Sprintf("found %d participation inconsistencies: %v", len(inconsistencies), inconsistencies),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"mq-participation-consistency",
			"all participation data is consistent",
		), false
	}
}

// AllInvariants runs all trust invariants
func (k *MQKeeper) AllInvariants(ctx sdk.Context) (string, bool) {
	res, stop := k.MQScoreValidityInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.MQConsensusRateInvariant()(ctx)
	if stop {
		return res, stop
	}

	return k.MQParticipationConsistencyInvariant()(ctx)
}
