package types

import (
	"fmt"
)

// GenesisState defines the taskmarket module's genesis state.
type GenesisState struct {
	Tasks        []Task        `json:"tasks"`
	Applications []Application `json:"applications"`
	Auctions     []Auction     `json:"auctions"`
	Bids         []Bid         `json:"bids"`
	Ratings      []Rating      `json:"ratings"`
	Reputations  []Reputation  `json:"reputations"`
}

// DefaultGenesis returns default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Tasks:        []Task{},
		Applications: []Application{},
		Auctions:     []Auction{},
		Bids:         []Bid{},
		Ratings:      []Rating{},
		Reputations:  []Reputation{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	for _, task := range data.Tasks {
		if err := task.Validate(); err != nil {
			return fmt.Errorf("invalid task %s: %w", task.ID, err)
		}
	}

	for _, app := range data.Applications {
		if err := app.Validate(); err != nil {
			return fmt.Errorf("invalid application %s: %w", app.ID, err)
		}
	}

	for _, bid := range data.Bids {
		if err := bid.Validate(); err != nil {
			return fmt.Errorf("invalid bid %s: %w", bid.ID, err)
		}
	}

	for _, rating := range data.Ratings {
		if err := rating.Validate(); err != nil {
			return fmt.Errorf("invalid rating %s: %w", rating.ID, err)
		}
	}

	return nil
}
