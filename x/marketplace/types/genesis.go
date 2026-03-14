package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Services: []Service{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(gs GenesisState) error {
	seenIDs := make(map[string]bool)
	for _, service := range gs.Services {
		if service.Id == "" {
			return ErrInvalidService.Wrap("service ID cannot be empty")
		}
		if seenIDs[service.Id] {
			return ErrInvalidService.Wrapf("duplicate service ID: %s", service.Id)
		}
		seenIDs[service.Id] = true
		if service.Provider == "" {
			return ErrInvalidService.Wrapf("service provider cannot be empty for ID: %s", service.Id)
		}
	}
	return nil
}

