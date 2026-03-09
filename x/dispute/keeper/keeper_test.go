package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/dispute/types"
)

func TestNewDispute(t *testing.T) {
	dispute := types.NewDispute("dispute-1", "escrow-1", "requester", "provider", "service not delivered")

	require.Equal(t, "dispute-1", dispute.ID)
	require.Equal(t, "escrow-1", dispute.EscrowID)
	require.Equal(t, "requester", dispute.Requester)
	require.Equal(t, "provider", dispute.Provider)
	require.Equal(t, types.DisputeStatusOpen, dispute.Status)
	require.Equal(t, "service not delivered", dispute.Reason)
}

func TestAddEvidence(t *testing.T) {
	dispute := types.NewDispute("1", "escrow-1", "req", "prov", "reason")

	evidence := types.Evidence{
		SubmittedBy: "requester",
		Type:        "text",
		Content:     "The service was not delivered on time",
		Timestamp:   1234567890,
	}

	dispute.AddEvidence(evidence)
	require.Len(t, dispute.Evidence, 1)
}

func TestAddVote(t *testing.T) {
	dispute := types.NewDispute("1", "escrow-1", "req", "prov", "reason")

	vote := types.Vote{
		Voter:     "juror1",
		Weight:    sdk.NewDec(1),
		Decision:  "requester",
		Timestamp: 1234567890,
	}

	dispute.AddVote(vote)
	require.Len(t, dispute.Votes, 1)
}

func TestCalculateResult(t *testing.T) {
	dispute := types.NewDispute("1", "escrow-1", "req", "prov", "reason")

	// Add votes
	dispute.AddVote(types.Vote{Voter: "j1", Weight: sdk.NewDec(2), Decision: "requester"})
	dispute.AddVote(types.Vote{Voter: "j2", Weight: sdk.NewDec(1), Decision: "provider"})
	dispute.AddVote(types.Vote{Voter: "j3", Weight: sdk.NewDec(1), Decision: "split"})

	result := dispute.CalculateResult()
	require.Equal(t, "requester", result.Decision)
	require.True(t, result.RequesterVotes.Equal(sdk.NewDec(2)))
	require.True(t, result.ProviderVotes.Equal(sdk.NewDec(1)))
	require.True(t, result.SplitVotes.Equal(sdk.NewDec(1)))
	require.True(t, result.TotalWeight.Equal(sdk.NewDec(4)))
}

func TestCalculateResultSplit(t *testing.T) {
	dispute := types.NewDispute("1", "escrow-1", "req", "prov", "reason")

	// Add votes with split winning
	dispute.AddVote(types.Vote{Voter: "j1", Weight: sdk.NewDec(1), Decision: "requester"})
	dispute.AddVote(types.Vote{Voter: "j2", Weight: sdk.NewDec(2), Decision: "split"})
	dispute.AddVote(types.Vote{Voter: "j3", Weight: sdk.NewDec(1), Decision: "provider"})

	result := dispute.CalculateResult()
	require.Equal(t, "split", result.Decision)
}

func TestDisputeString(t *testing.T) {
	dispute := types.NewDispute("dispute-123", "escrow-1", "req", "prov", "reason")
	str := dispute.String()

	require.Contains(t, str, "dispute-123")
	require.Contains(t, str, "req")
	require.Contains(t, str, "prov")
}
