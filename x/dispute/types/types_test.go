package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sharetoken/x/dispute/types"
)

func TestNewDispute(t *testing.T) {
	dispute := types.NewDispute("dispute-1", "escrow-1", "requester-1", "provider-1", "service not delivered")

	require.NotNil(t, dispute)
	require.Equal(t, "dispute-1", dispute.Id)
	require.Equal(t, "escrow-1", dispute.EscrowId)
	require.Equal(t, "requester-1", dispute.Requester)
	require.Equal(t, "provider-1", dispute.Provider)
	require.Equal(t, types.DisputeStatusOpen, dispute.Status)
	require.Equal(t, "service not delivered", dispute.Reason)
	require.NotNil(t, dispute.Evidence)
	require.Empty(t, dispute.Evidence)
	require.NotNil(t, dispute.Votes)
	require.Empty(t, dispute.Votes)
}

func TestDispute_AddEvidence(t *testing.T) {
	dispute := types.NewDispute("dispute-1", "escrow-1", "requester-1", "provider-1", "reason")

	require.Empty(t, dispute.Evidence)

	evidence := types.Evidence{
		SubmittedBy: "requester-1",
		EvidenceType: "document",
		Content:     "proof of payment",
		Timestamp:   1234567890,
	}

	dispute.AddEvidence(evidence)
	require.Len(t, dispute.Evidence, 1)
	require.Equal(t, evidence, dispute.Evidence[0])

	// Add more evidence
	evidence2 := types.Evidence{
		SubmittedBy: "provider-1",
		EvidenceType: "image",
		Content:     "delivery confirmation",
		Timestamp:   1234567891,
	}
	dispute.AddEvidence(evidence2)
	require.Len(t, dispute.Evidence, 2)
}

func TestDispute_AddVote(t *testing.T) {
	dispute := types.NewDispute("dispute-1", "escrow-1", "requester-1", "provider-1", "reason")

	require.Empty(t, dispute.Votes)

	vote := types.Vote{
		Voter:     "voter-1",
		Weight:    "1",
		Decision:  "requester",
		Timestamp: 1234567890,
	}

	dispute.AddVote(vote)
	require.Len(t, dispute.Votes, 1)
	require.Equal(t, vote, dispute.Votes[0])

	// Add more votes
	vote2 := types.Vote{
		Voter:     "voter-2",
		Weight:    "2",
		Decision:  "provider",
		Timestamp: 1234567891,
	}
	dispute.AddVote(vote2)
	require.Len(t, dispute.Votes, 2)
}

func TestDispute_CalculateResult(t *testing.T) {
	tests := []struct {
		name           string
		votes          []types.Vote
		expectedWinner string
	}{
		{
			name:           "no votes",
			votes:          []types.Vote{},
			expectedWinner: "split",
		},
		{
			name: "requester wins",
			votes: []types.Vote{
				{Voter: "v1", Weight: "3", Decision: "requester"},
				{Voter: "v2", Weight: "2", Decision: "provider"},
				{Voter: "v3", Weight: "1", Decision: "split"},
			},
			expectedWinner: "requester",
		},
		{
			name: "provider wins",
			votes: []types.Vote{
				{Voter: "v1", Weight: "2", Decision: "requester"},
				{Voter: "v2", Weight: "4", Decision: "provider"},
				{Voter: "v3", Weight: "1", Decision: "split"},
			},
			expectedWinner: "provider",
		},
		{
			name: "split wins",
			votes: []types.Vote{
				{Voter: "v1", Weight: "2", Decision: "requester"},
				{Voter: "v2", Weight: "2", Decision: "provider"},
				{Voter: "v3", Weight: "5", Decision: "split"},
			},
			expectedWinner: "split",
		},
		{
			name: "tie goes to split",
			votes: []types.Vote{
				{Voter: "v1", Weight: "2", Decision: "requester"},
				{Voter: "v2", Weight: "2", Decision: "provider"},
			},
			expectedWinner: "split",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispute := types.NewDispute("dispute-1", "escrow-1", "requester-1", "provider-1", "reason")
			dispute.Votes = tt.votes

			result := dispute.CalculateResult()
			require.Equal(t, tt.expectedWinner, result.Decision)
		})
	}
}

func TestDispute_String(t *testing.T) {
	dispute := types.NewDispute("dispute-1", "escrow-1", "requester-1", "provider-1", "reason")
	result := dispute.String()

	require.Contains(t, result, "dispute-1")
	require.Contains(t, result, "requester-1")
	require.Contains(t, result, "provider-1")
	require.Contains(t, result, "open")
}

// Genesis Tests

func TestDefaultDisputeGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Disputes)
	require.Empty(t, genesis.Disputes)
}

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Disputes)
	require.NotNil(t, genesis.JurorPool)
	require.Empty(t, genesis.Disputes)
	require.Empty(t, genesis.JurorPool)
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
			name: "valid genesis with disputes",
			data: types.GenesisState{
				Disputes: []types.Dispute{
					{Id: "dispute-1", EscrowId: "escrow-1", Requester: "requester-1", Provider: "provider-1"},
					{Id: "dispute-2", EscrowId: "escrow-2", Requester: "requester-2", Provider: "provider-2"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate dispute IDs",
			data: types.GenesisState{
				Disputes: []types.Dispute{
					{Id: "dispute-1", EscrowId: "escrow-1", Requester: "requester-1", Provider: "provider-1"},
					{Id: "dispute-1", EscrowId: "escrow-2", Requester: "requester-2", Provider: "provider-2"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty dispute ID",
			data: types.GenesisState{
				Disputes: []types.Dispute{
					{Id: "", EscrowId: "escrow-1", Requester: "requester-1", Provider: "provider-1"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty escrow ID",
			data: types.GenesisState{
				Disputes: []types.Dispute{
					{Id: "dispute-1", EscrowId: "", Requester: "requester-1", Provider: "provider-1"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty requester",
			data: types.GenesisState{
				Disputes: []types.Dispute{
					{Id: "dispute-1", EscrowId: "escrow-1", Requester: "", Provider: "provider-1"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty provider",
			data: types.GenesisState{
				Disputes: []types.Dispute{
					{Id: "dispute-1", EscrowId: "escrow-1", Requester: "requester-1", Provider: ""},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - duplicate jurors",
			data: types.GenesisState{
				Disputes: []types.Dispute{},
				JurorPool: []string{"juror-1", "juror-1"},
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
