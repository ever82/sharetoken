package keeper

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"sharetoken/x/trust/types"
)

// MQKeeper manages MQ (Meta Quality) scores for users
type MQKeeper struct {
	mu sync.RWMutex

	// MQ scores storage: address -> MQScore
	scores map[string]*types.MQScore

	// Dispute participation history: address -> list of booleans (true = voted with consensus)
	participation map[string][]bool

	// Cosmos SDK storage
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewMQKeeper creates a new MQ keeper (legacy in-memory version)
func NewMQKeeper() *MQKeeper {
	return &MQKeeper{
		scores:        make(map[string]*types.MQScore),
		participation: make(map[string][]bool),
	}
}

// NewKeeper creates a new MQ keeper for Cosmos SDK module (uses KVStore)
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) MQKeeper {
	return MQKeeper{
		scores:        make(map[string]*types.MQScore),
		participation: make(map[string][]bool),
		cdc:           cdc,
		storeKey:      storeKey,
	}
}

// GetScore retrieves the MQ score for an address
func (k *MQKeeper) GetScore(address string) *types.MQScore {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if score, exists := k.scores[address]; exists {
		return score
	}

	// Return default score (100) if not exists
	return types.NewMQScore(address)
}

// InitializeScore initializes MQ score for a new user
func (k *MQKeeper) InitializeScore(address string) *types.MQScore {
	k.mu.Lock()
	defer k.mu.Unlock()

	if _, exists := k.scores[address]; exists {
		return k.scores[address]
	}

	score := types.NewMQScore(address)
	k.scores[address] = score
	return score
}

// CalculateVotingWeight calculates voting weight based on MQ score
func (k *MQKeeper) CalculateVotingWeight(address string) int64 {
	score := k.GetScore(address)
	// Weight = MQ score (for integer calculations)
	return int64(score.Score)
}

// RecordDisputeParticipation records a user's participation in a dispute
func (k *MQKeeper) RecordDisputeParticipation(address, disputeID string, votedWithConsensus bool) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Store participation history
	k.participation[address] = append(k.participation[address], votedWithConsensus)

	// Get or initialize score
	score, exists := k.scores[address]
	if !exists {
		score = types.NewMQScore(address)
		k.scores[address] = score
	}

	// Record the dispute participation
	score.RecordDispute(votedWithConsensus)

	// Apply penalty or reward
	if votedWithConsensus {
		// Calculate reward
		votersCount := int64(len(k.participation[address]))
		reward := score.CalculateReward(votersCount, true)
		score.ApplyReward(reward)
	} else {
		// Calculate penalty
		penalty := score.CalculatePenalty(false)
		score.ApplyPenalty(penalty)
	}

	return nil
}

// GetParticipationRate returns the participation rate for an address
func (k *MQKeeper) GetParticipationRate(address string) int64 {
	k.mu.RLock()
	defer k.mu.RUnlock()

	participations := k.participation[address]
	if len(participations) == 0 {
		return 100
	}

	consensus := 0
	for _, p := range participations {
		if p {
			consensus++
		}
	}

	return int64(consensus * 100 / len(participations))
}

// GetAllScores returns all MQ scores
func (k *MQKeeper) GetAllScores() map[string]*types.MQScore {
	k.mu.RLock()
	defer k.mu.RUnlock()

	result := make(map[string]*types.MQScore)
	for addr, score := range k.scores {
		result[addr] = score
	}
	return result
}

// GetTopScorers returns top N users by MQ score
func (k *MQKeeper) GetTopScorers(n int) []*types.MQScore {
	k.mu.RLock()
	defer k.mu.RUnlock()

	allScores := make([]*types.MQScore, 0, len(k.scores))
	for _, score := range k.scores {
		allScores = append(allScores, score)
	}

	// Sort by score descending
	sort.Slice(allScores, func(i, j int) bool {
		return allScores[i].Score > allScores[j].Score
	})

	if n > len(allScores) {
		n = len(allScores)
	}

	if n == 0 {
		return []*types.MQScore{}
	}

	return allScores[:n]
}

// ValidateScore validates if a score is within valid range
func (k *MQKeeper) ValidateScore(score int32) error {
	if score < 0 || score > 100 {
		return fmt.Errorf("score must be between 0 and 100, got %d", score)
	}
	return nil
}

// SelectJurorsWeighted selects jurors weighted by their MQ score using weighted reservoir sampling.
//
// Algorithm: Weighted Reservoir Sampling (O(n))
// Instead of O(n²) repeated selection, we use a single pass algorithm:
// 1. Build prefix sum array of weights (O(n))
// 2. For each selection, do binary search on prefix sum (O(log n))
// Overall complexity: O(n log n) for selection, O(n) for preprocessing
//
// Alternative: If count is small relative to n, we can use reservoir sampling
// which achieves O(n) total complexity.
func (k *MQKeeper) SelectJurorsWeighted(candidates []string, count int) []string {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if len(candidates) == 0 || count <= 0 {
		return []string{}
	}

	// Calculate weights and build prefix sum array
	// This allows O(log n) selection via binary search
	type weightedCandidate struct {
		address string
		weight  int64
	}

	weighted := make([]weightedCandidate, 0, len(candidates))
	prefixSum := make([]int64, 0, len(candidates))
	var totalWeight int64

	for _, candidate := range candidates {
		score := k.GetScore(candidate)
		weight := int64(score.Score)
		if weight <= 0 {
			weight = 1 // Minimum weight
		}
		weighted = append(weighted, weightedCandidate{candidate, weight})
		totalWeight += weight
		prefixSum = append(prefixSum, totalWeight)
	}

	if totalWeight == 0 {
		if count > len(candidates) {
			count = len(candidates)
		}
		return candidates[:count]
	}

	if count > len(candidates) {
		count = len(candidates)
	}

	// Weighted random selection using binary search on prefix sum
	// Complexity: O(count * log n) instead of O(count * n)
	selected := make([]string, 0, count)
	used := make(map[int]bool) // Track indices to avoid duplicates

	for len(selected) < count {
		// Generate random target in [1, totalWeight]
		target := rand.Int63n(totalWeight) + 1 //nolint:gosec

		// Binary search to find the selected candidate
		// prefixSum[i] represents cumulative weight up to index i
		idx := binarySearchPrefixSum(prefixSum, target)

		// Skip if already selected
		if !used[idx] {
			used[idx] = true
			selected = append(selected, weighted[idx].address)
		}
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
