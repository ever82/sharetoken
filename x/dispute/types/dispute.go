package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DisputeStatusOpen is an alias for the protobuf enum value
const DisputeStatusOpen = DisputeStatus_DISPUTE_STATUS_OPEN

// DisputeStatusMediating is an alias for the protobuf enum value
const DisputeStatusMediating = DisputeStatus_DISPUTE_STATUS_MEDIATING

// DisputeStatusVoting is an alias for the protobuf enum value
const DisputeStatusVoting = DisputeStatus_DISPUTE_STATUS_VOTING

// DisputeStatusResolved is an alias for the protobuf enum value
const DisputeStatusResolved = DisputeStatus_DISPUTE_STATUS_RESOLVED

// DisputeStatusCancelled is an alias for the protobuf enum value
const DisputeStatusCancelled = DisputeStatus_DISPUTE_STATUS_CANCELLED

// NewDispute creates a new dispute
func NewDispute(id, escrowID, requester, provider, reason string) *Dispute {
	return &Dispute{
		Id:          id,
		EscrowId:    escrowID,
		Requester:   requester,
		Provider:    provider,
		Status:      DisputeStatus_DISPUTE_STATUS_OPEN,
		Reason:      reason,
		Evidence:    []Evidence{},
		Votes:       []Vote{},
		CreatedAt:   time.Now().Unix(),
		CompletedAt: 0,
	}
}

// AddEvidence adds evidence to a dispute
func (d *Dispute) AddEvidence(evidence Evidence) {
	d.Evidence = append(d.Evidence, evidence)
}

// AddVote adds a vote to a dispute
func (d *Dispute) AddVote(vote Vote) {
	d.Votes = append(d.Votes, vote)
}

// CalculateResult calculates the voting result
func (d *Dispute) CalculateResult() VoteResult {
	requesterWeight := sdk.ZeroDec()
	providerWeight := sdk.ZeroDec()
	splitWeight := sdk.ZeroDec()

	for _, vote := range d.Votes {
		weight, err := sdk.NewDecFromStr(vote.Weight)
		if err != nil {
			continue
		}
		switch vote.Decision {
		case "requester":
			requesterWeight = requesterWeight.Add(weight)
		case "provider":
			providerWeight = providerWeight.Add(weight)
		case "split":
			splitWeight = splitWeight.Add(weight)
		}
	}

	totalWeight := requesterWeight.Add(providerWeight).Add(splitWeight)

	var decision string
	if requesterWeight.GT(providerWeight) && requesterWeight.GT(splitWeight) {
		decision = "requester"
	} else if providerWeight.GT(requesterWeight) && providerWeight.GT(splitWeight) {
		decision = "provider"
	} else {
		decision = "split"
	}

	return VoteResult{
		Decision:       decision,
		RequesterVotes: requesterWeight.String(),
		ProviderVotes:  providerWeight.String(),
		SplitVotes:     splitWeight.String(),
		TotalWeight:    totalWeight.String(),
	}
}

// NewEvidence creates a new evidence
func NewEvidence(submittedBy, evidenceType, content string) Evidence {
	return Evidence{
		SubmittedBy:  submittedBy,
		EvidenceType: evidenceType,
		Content:      content,
		Timestamp:    time.Now().Unix(),
	}
}
