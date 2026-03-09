package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DisputeStatus represents the status of a dispute
type DisputeStatus string

const (
	DisputeStatusOpen       DisputeStatus = "open"
	DisputeStatusMediating  DisputeStatus = "mediating"
	DisputeStatusVoting     DisputeStatus = "voting"
	DisputeStatusResolved   DisputeStatus = "resolved"
	DisputeStatusCancelled  DisputeStatus = "cancelled"
)

// Dispute represents a dispute between parties
type Dispute struct {
	ID            string        `json:"id"`
	EscrowID      string        `json:"escrow_id"`
	Requester     string        `json:"requester"`
	Provider      string        `json:"provider"`
	Status        DisputeStatus `json:"status"`
	Reason        string        `json:"reason"`
	Evidence      []Evidence    `json:"evidence"`
	Votes         []Vote        `json:"votes"`
	Result        VoteResult    `json:"result"`
	CreatedAt     int64         `json:"created_at"`
	CompletedAt   int64         `json:"completed_at"`
}

// Evidence represents evidence submitted in a dispute
type Evidence struct {
	SubmittedBy string `json:"submitted_by"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Timestamp   int64  `json:"timestamp"`
}

// Vote represents a vote in a dispute
type Vote struct {
	Voter     string  `json:"voter"`
	Weight    sdk.Dec `json:"weight"`
	Decision  string  `json:"decision"` // "requester", "provider", or "split"
	Timestamp int64   `json:"timestamp"`
}

// VoteResult represents the result of voting
type VoteResult struct {
	Decision      string    `json:"decision"`
	RequesterVotes sdk.Dec  `json:"requester_votes"`
	ProviderVotes  sdk.Dec  `json:"provider_votes"`
	SplitVotes     sdk.Dec  `json:"split_votes"`
	TotalWeight    sdk.Dec  `json:"total_weight"`
}

// NewDispute creates a new dispute
func NewDispute(id, escrowID, requester, provider, reason string) *Dispute {
	return &Dispute{
		ID:        id,
		EscrowID:  escrowID,
		Requester: requester,
		Provider:  provider,
		Status:    DisputeStatusOpen,
		Reason:    reason,
		Evidence:  []Evidence{},
		Votes:     []Vote{},
	}
}

// AddEvidence adds evidence to the dispute
func (d *Dispute) AddEvidence(evidence Evidence) {
	d.Evidence = append(d.Evidence, evidence)
}

// AddVote adds a vote to the dispute
func (d *Dispute) AddVote(vote Vote) {
	d.Votes = append(d.Votes, vote)
}

// CalculateResult calculates the voting result
func (d *Dispute) CalculateResult() VoteResult {
	result := VoteResult{
		RequesterVotes: sdk.ZeroDec(),
		ProviderVotes:  sdk.ZeroDec(),
		SplitVotes:     sdk.ZeroDec(),
		TotalWeight:    sdk.ZeroDec(),
	}

	for _, vote := range d.Votes {
		switch vote.Decision {
		case "requester":
			result.RequesterVotes = result.RequesterVotes.Add(vote.Weight)
		case "provider":
			result.ProviderVotes = result.ProviderVotes.Add(vote.Weight)
		case "split":
			result.SplitVotes = result.SplitVotes.Add(vote.Weight)
		}
		result.TotalWeight = result.TotalWeight.Add(vote.Weight)
	}

	// Determine the decision
	if result.RequesterVotes.GT(result.ProviderVotes) && result.RequesterVotes.GT(result.SplitVotes) {
		result.Decision = "requester"
	} else if result.ProviderVotes.GT(result.RequesterVotes) && result.ProviderVotes.GT(result.SplitVotes) {
		result.Decision = "provider"
	} else {
		result.Decision = "split"
	}

	return result
}

// String implements stringer
func (d Dispute) String() string {
	return fmt.Sprintf("Dispute{%s: %s vs %s, status: %s, votes: %d}",
		d.ID, d.Requester, d.Provider, d.Status, len(d.Votes))
}

// DisputeGenesis represents genesis state
type DisputeGenesis struct {
	Disputes []Dispute `json:"disputes"`
}

// DefaultDisputeGenesis returns default genesis
func DefaultDisputeGenesis() *DisputeGenesis {
	return &DisputeGenesis{
		Disputes: []Dispute{},
	}
}
