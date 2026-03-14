package types

import (
	"fmt"
)

// DefaultGenesis returns default genesis state for the taskmarket module
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Tasks:       []Task{},
		Applications: []Application{},
		Auctions:    []Auction{},
		Bids:        []Bid{},
		Ratings:     []Rating{},
		Reputations: []Reputation{},
	}
}

// Validate validates the genesis state
func (gs *GenesisState) Validate() error {
	// Validate tasks
	taskIDs := make(map[string]bool)
	for _, task := range gs.Tasks {
		if task.Id == "" {
			return fmt.Errorf("task ID cannot be empty")
		}
		if taskIDs[task.Id] {
			return fmt.Errorf("duplicate task ID: %s", task.Id)
		}
		taskIDs[task.Id] = true
	}

	// Validate applications
	appIDs := make(map[string]bool)
	for _, app := range gs.Applications {
		if app.Id == "" {
			return fmt.Errorf("application ID cannot be empty")
		}
		if appIDs[app.Id] {
			return fmt.Errorf("duplicate application ID: %s", app.Id)
		}
		appIDs[app.Id] = true
	}

	// Validate bids
	bidIDs := make(map[string]bool)
	for _, bid := range gs.Bids {
		if bid.Id == "" {
			return fmt.Errorf("bid ID cannot be empty")
		}
		if bidIDs[bid.Id] {
			return fmt.Errorf("duplicate bid ID: %s", bid.Id)
		}
		bidIDs[bid.Id] = true
	}

	// Validate ratings
	ratingIDs := make(map[string]bool)
	for _, rating := range gs.Ratings {
		if rating.Id == "" {
			return fmt.Errorf("rating ID cannot be empty")
		}
		if ratingIDs[rating.Id] {
			return fmt.Errorf("duplicate rating ID: %s", rating.Id)
		}
		ratingIDs[rating.Id] = true
	}

	return nil
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(gs GenesisState) error {
	return gs.Validate()
}
