package keeper

import (
	"fmt"

	"sharetoken/x/crowdfunding/types"
)

// CreateIdea creates a new idea
func (k *Keeper) CreateIdea(idea *types.Idea) error {
	if err := idea.Validate(); err != nil {
		return fmt.Errorf("invalid idea: %w", err)
	}

	k.mutex.Lock()
	k.ideas[idea.ID] = idea
	k.mutex.Unlock()

	return nil
}

// GetIdea retrieves an idea
func (k *Keeper) GetIdea(id string) *types.Idea {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	return k.ideas[id]
}

// UpdateIdea updates an idea and creates a version record
func (k *Keeper) UpdateIdea(ideaID, title, description, changes, updatedBy string) (*types.IdeaVersion, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	idea, exists := k.ideas[ideaID]
	if !exists {
		return nil, fmt.Errorf("idea not found: %s", ideaID)
	}

	// Create version record
	version := types.NewIdeaVersion(
		fmt.Sprintf("%s-v%d", ideaID, idea.CurrentVersion),
		ideaID,
		idea.CurrentVersion,
		idea.Title,
		idea.Description,
		changes,
		updatedBy,
	)

	// Store version
	k.versions[ideaID] = append(k.versions[ideaID], version)

	// Update idea
	idea.Update(title, description)

	return version, nil
}

// GetIdeaVersions retrieves all versions of an idea
func (k *Keeper) GetIdeaVersions(ideaID string) []*types.IdeaVersion {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	return k.versions[ideaID]
}

// PublishIdea publishes an idea
func (k *Keeper) PublishIdea(ideaID string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	idea, exists := k.ideas[ideaID]
	if !exists {
		return fmt.Errorf("idea not found: %s", ideaID)
	}

	idea.Publish()
	return nil
}
