package keeper

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/dispute/types"
)

// DisputeKeeper manages disputes and arbitration
type DisputeKeeper struct {
	mu sync.RWMutex

	// Disputes storage
	disputes map[string]*types.Dispute

	// Juror pool
	jurorPool []string

	// MQ keeper reference (for juror selection)
	mqKeeper MQKeeperInterface
}

// MQKeeperInterface defines the interface for MQ operations
type MQKeeperInterface interface {
	GetScore(address string) interface{}
	CalculateVotingWeight(address string) int64
}

// NewDisputeKeeper creates a new dispute keeper
func NewDisputeKeeper(mqKeeper MQKeeperInterface) *DisputeKeeper {
	return &DisputeKeeper{
		disputes:  make(map[string]*types.Dispute),
		jurorPool: make([]string, 0),
		mqKeeper:  mqKeeper,
	}
}

// CreateDispute creates a new dispute
func (k *DisputeKeeper) CreateDispute(escrowID, requester, provider, reason string) (*types.Dispute, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	disputeID := generateDisputeID()

	dispute := types.NewDispute(disputeID, escrowID, requester, provider, reason)

	k.disputes[disputeID] = dispute
	return dispute, nil
}

// AddEvidence adds evidence to a dispute
func (k *DisputeKeeper) AddEvidence(disputeID, submitter, evidenceType, content string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	dispute, exists := k.disputes[disputeID]
	if !exists {
		return fmt.Errorf("dispute %s not found", disputeID)
	}

	if dispute.Status != types.DisputeStatusOpen {
		return fmt.Errorf("cannot add evidence to dispute in status %s", dispute.Status)
	}

	evidence := types.Evidence{
		Type:        evidenceType,
		Content:     content,
		SubmittedBy: submitter,
		Timestamp:   time.Now().Unix(),
	}

	dispute.AddEvidence(evidence)
	return nil
}

// StartMediation transitions dispute to mediation phase
func (k *DisputeKeeper) StartMediation(disputeID string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	dispute, exists := k.disputes[disputeID]
	if !exists {
		return fmt.Errorf("dispute %s not found", disputeID)
	}

	if dispute.Status != types.DisputeStatusOpen {
		return fmt.Errorf("cannot start mediation for dispute in status %s", dispute.Status)
	}

	dispute.Status = types.DisputeStatusMediating
	return nil
}

// StartJuryVoting transitions dispute to jury voting phase
func (k *DisputeKeeper) StartJuryVoting(disputeID string, jurySize int) ([]string, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	dispute, exists := k.disputes[disputeID]
	if !exists {
		return nil, fmt.Errorf("dispute %s not found", disputeID)
	}

	// Select jury based on MQ-weighted random selection
	jury := k.selectWeightedRandomJurors(jurySize)
	dispute.Status = types.DisputeStatusVoting

	return jury, nil
}

// selectWeightedRandomJurors selects jurors with MQ-weighted random algorithm
func (k *DisputeKeeper) selectWeightedRandomJurors(count int) []string {
	if len(k.jurorPool) == 0 || count <= 0 {
		return []string{}
	}

	// Calculate total weight
	type weightedJuror struct {
		address string
		weight  int64
	}

	weighted := make([]weightedJuror, 0, len(k.jurorPool))
	var totalWeight int64

	for _, juror := range k.jurorPool {
		weight := k.mqKeeper.CalculateVotingWeight(juror)
		if weight <= 0 {
			weight = 1 // Minimum weight
		}
		weighted = append(weighted, weightedJuror{juror, weight})
		totalWeight += weight
	}

	if totalWeight == 0 {
		return []string{}
	}

	// Select jurors based on weights
	selected := make([]string, 0, count)
	used := make(map[string]bool)

	for len(selected) < count && len(used) < len(weighted) {
		target := rand.Int63n(totalWeight) //nolint:gosec
		var current int64

		for _, wj := range weighted {
			if used[wj.address] {
				continue
			}
			current += wj.weight
			if current >= target {
				selected = append(selected, wj.address)
				used[wj.address] = true
				totalWeight -= wj.weight
				break
			}
		}
	}

	return selected
}

// CastVote casts a vote in a dispute
func (k *DisputeKeeper) CastVote(disputeID, voter, decision string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	dispute, exists := k.disputes[disputeID]
	if !exists {
		return fmt.Errorf("dispute %s not found", disputeID)
	}

	if dispute.Status != types.DisputeStatusVoting {
		return fmt.Errorf("dispute is not in voting phase")
	}

	weight := k.mqKeeper.CalculateVotingWeight(voter)
	vote := types.Vote{
		Voter:     voter,
		Decision:  decision,
		Weight:    sdk.NewDec(weight),
		Timestamp: time.Now().Unix(),
	}

	dispute.AddVote(vote)
	return nil
}

// FinalizeVoting calculates voting results and resolves dispute
func (k *DisputeKeeper) FinalizeVoting(disputeID string) (*types.VoteResult, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	dispute, exists := k.disputes[disputeID]
	if !exists {
		return nil, fmt.Errorf("dispute %s not found", disputeID)
	}

	if dispute.Status != types.DisputeStatusVoting {
		return nil, fmt.Errorf("dispute is not in voting phase")
	}

	// Calculate result
	result := dispute.CalculateResult()
	dispute.Result = result
	dispute.Status = types.DisputeStatusResolved
	dispute.CompletedAt = time.Now().Unix()

	return &result, nil
}

// GetDispute retrieves a dispute by ID
func (k *DisputeKeeper) GetDispute(disputeID string) (*types.Dispute, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	dispute, exists := k.disputes[disputeID]
	if !exists {
		return nil, fmt.Errorf("dispute %s not found", disputeID)
	}
	return dispute, nil
}

// GetAllDisputes returns all disputes
func (k *DisputeKeeper) GetAllDisputes() []*types.Dispute {
	k.mu.RLock()
	defer k.mu.RUnlock()

	result := make([]*types.Dispute, 0, len(k.disputes))
	for _, dispute := range k.disputes {
		result = append(result, dispute)
	}

	// Sort by creation time (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})

	return result
}

// RegisterJuror adds a juror to the pool
func (k *DisputeKeeper) RegisterJuror(address string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Check if already registered
	for _, j := range k.jurorPool {
		if j == address {
			return
		}
	}

	k.jurorPool = append(k.jurorPool, address)
}

// generateDisputeID generates a unique dispute ID
func generateDisputeID() string {
	return fmt.Sprintf("DISPUTE-%d-%d", time.Now().Unix(), rand.Intn(10000)) //nolint:gosec
}
