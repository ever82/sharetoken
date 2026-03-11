package keeper

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

	disputeID, err := generateDisputeID()
	if err != nil {
		return nil, err
	}

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

// selectWeightedRandomJurors selects jurors with MQ-weighted random algorithm.
//
// Algorithm: Weighted Random Selection with Prefix Sum + Binary Search (O(n log n))
// Instead of O(n²) repeated linear scan, we:
// 1. Calculate all weights and build prefix sum array (O(n))
// 2. Use binary search for each selection (O(log n) per selection)
// Overall complexity: O(n + count * log n) vs O(n²) before
func (k *DisputeKeeper) selectWeightedRandomJurors(count int) []string {
	if len(k.jurorPool) == 0 || count <= 0 {
		return []string{}
	}

	// Calculate weights and build prefix sum array
	// This allows O(log n) selection via binary search
	type weightedJuror struct {
		address string
		weight  int64
	}

	weighted := make([]weightedJuror, 0, len(k.jurorPool))
	prefixSum := make([]int64, 0, len(k.jurorPool))
	var totalWeight int64

	for _, juror := range k.jurorPool {
		weight := k.mqKeeper.CalculateVotingWeight(juror)
		if weight <= 0 {
			weight = 1 // Minimum weight
		}
		weighted = append(weighted, weightedJuror{juror, weight})
		totalWeight += weight
		prefixSum = append(prefixSum, totalWeight)
	}

	if totalWeight == 0 {
		return []string{}
	}

	if count > len(weighted) {
		count = len(weighted)
	}

	// Weighted random selection using binary search on prefix sum
	// Complexity: O(count * log n) instead of O(count * n)
	selected := make([]string, 0, count)
	used := make(map[int]bool) // Track by index to avoid duplicates

	for len(selected) < count {
		// Generate random number in range [1, totalWeight]
		targetBig, err := rand.Int(rand.Reader, big.NewInt(totalWeight))
		if err != nil {
			// Fallback: if crypto/rand fails, break to avoid infinite loop
			break
		}
		target := targetBig.Int64() + 1 // Range [1, totalWeight]

		// Binary search to find the selected juror
		idx := binarySearchPrefixSum(prefixSum, target)

		// Skip if already selected
		if used[idx] {
			continue
		}

		used[idx] = true
		selected = append(selected, weighted[idx].address)
	}

	return selected
}

// binarySearchPrefixSum finds the smallest index i such that prefixSum[i] >= target
// using binary search. Complexity: O(log n)
func binarySearchPrefixSum(prefixSum []int64, target int64) int {
	left, right := 0, len(prefixSum)-1
	for left < right {
		mid := left + (right-left)/2
		if prefixSum[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
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
func generateDisputeID() (string, error) {
	// Generate random number for the ID suffix
	randBig, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}
	return fmt.Sprintf("DISPUTE-%d-%d", time.Now().Unix(), randBig.Int64()), nil
}
