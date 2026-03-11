package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/identity/types"
)

// SetIdentity sets an identity in the store
func (k Keeper) SetIdentity(ctx sdk.Context, identity types.Identity) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIdentityKey(identity.Address)
	value, err := json.Marshal(identity)
	if err != nil {
		return fmt.Errorf("failed to marshal identity: %w", err)
	}
	store.Set(key, value)
	return nil
}

// GetIdentity retrieves an identity by address
func (k Keeper) GetIdentity(ctx sdk.Context, address string) (types.Identity, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIdentityKey(address)
	value := store.Get(key)

	if value == nil {
		return types.Identity{}, false
	}

	var identity types.Identity
	if err := json.Unmarshal(value, &identity); err != nil {
		return types.Identity{}, false
	}
	return identity, true
}

// HasIdentity checks if an identity exists
func (k Keeper) HasIdentity(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIdentityKey(address)
	return store.Has(key)
}

// DeleteIdentity deletes an identity from the store
func (k Keeper) DeleteIdentity(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIdentityKey(address)
	store.Delete(key)
}

// RegisterIdentity registers a new identity
func (k Keeper) RegisterIdentity(ctx sdk.Context, address, did, metadataHash string) error {
	// Check if identity already exists
	if k.HasIdentity(ctx, address) {
		return types.ErrIdentityAlreadyExists.Wrap(address)
	}

	// Check if DID is already registered (if provided)
	if did != "" && k.IsDIDRegistered(ctx, did) {
		return types.ErrDIDAlreadyRegistered.Wrap(did)
	}

	// Create new identity
	identity := types.NewIdentity(ctx, address, did)
	identity.MetadataHash = metadataHash

	// Generate merkle root
	identity.MerkleRoot = identity.GenerateMerkleRoot()

	// Validate identity
	if err := identity.ValidateBasic(); err != nil {
		return err
	}

	// Store identity
	if err := k.SetIdentity(ctx, *identity); err != nil {
		return err
	}

	// Register DID if provided
	if did != "" {
		k.RegisterDID(ctx, did)
	}

	// Initialize default limit config
	limitConfig := types.NewLimitConfig(address)
	if err := k.SetLimitConfig(ctx, limitConfig); err != nil {
		return err
	}

	return nil
}

// VerifyIdentity verifies an identity with a third-party provider
func (k Keeper) VerifyIdentity(ctx sdk.Context, address, provider, verificationHash, proof string) error {
	// Check if identity exists
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound.Wrap(address)
	}

	// Check if provider is valid
	if !types.IsValidProvider(provider) {
		return types.ErrInvalidProvider.Wrap(provider)
	}

	// Check if provider is already used for this account
	if k.IsProviderUsed(ctx, provider, address) {
		return types.ErrProviderAlreadyUsed.Wrapf("provider %s already used for %s", provider, address)
	}

	// Verify the proof (simplified - in production, this would validate OAuth tokens)
	// For now, we accept any non-empty proof
	if proof == "" {
		return types.ErrInvalidProof.Wrap("verification proof is required")
	}

	// Update identity
	identity.VerificationProvider = provider
	identity.VerificationHash = verificationHash
	identity.IsVerified = true

	// Generate new merkle root
	identity.MerkleRoot = identity.GenerateMerkleRoot()

	// Store updated identity
	if err := k.SetIdentity(ctx, identity); err != nil {
		return err
	}

	// Mark provider as used
	k.MarkProviderUsed(ctx, provider, address)

	return nil
}

// GetAllIdentities returns all identities
func (k Keeper) GetAllIdentities(ctx sdk.Context) []types.Identity {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IdentityKey)
	defer iterator.Close()

	var identities []types.Identity
	for ; iterator.Valid(); iterator.Next() {
		var identity types.Identity
		if err := json.Unmarshal(iterator.Value(), &identity); err != nil {
			ctx.Logger().Error("failed to unmarshal identity", "error", err)
			continue
		}
		identities = append(identities, identity)
	}

	return identities
}

// IsDIDRegistered checks if a DID is already registered
func (k Keeper) IsDIDRegistered(ctx sdk.Context, did string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRegisteredDIDKey(did)
	return store.Has(key)
}

// RegisterDID registers a DID
func (k Keeper) RegisterDID(ctx sdk.Context, did string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRegisteredDIDKey(did)
	store.Set(key, []byte{1})
}

// IsProviderUsed checks if a provider is already used for an account
func (k Keeper) IsProviderUsed(ctx sdk.Context, provider, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVerificationProviderKey(provider, address)
	return store.Has(key)
}

// MarkProviderUsed marks a provider as used for an account
func (k Keeper) MarkProviderUsed(ctx sdk.Context, provider, address string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVerificationProviderKey(provider, address)
	store.Set(key, []byte{1})
}

// IsVerified checks if an address is verified
func (k Keeper) IsVerified(ctx sdk.Context, address string) bool {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return false
	}
	return identity.IsVerified
}

// RequireVerification checks if verification is required and address is verified
func (k Keeper) RequireVerification(ctx sdk.Context, address string) error {
	params := k.GetParams(ctx)
	if !params.VerificationRequired {
		return nil
	}

	if !k.IsVerified(ctx, address) {
		return types.ErrInvalidVerification.Wrapf("address %s is not verified", address)
	}

	return nil
}
