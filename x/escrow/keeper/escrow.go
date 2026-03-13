package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/escrow/types"
)

// SetEscrow sets an escrow in the store
func (k Keeper) SetEscrow(ctx sdk.Context, escrow types.Escrow) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(escrow.ID)
	value, err := json.Marshal(escrow)
	if err != nil {
		return fmt.Errorf("failed to marshal escrow: %w", err)
	}
	store.Set(key, value)

	// Update indexes
	k.setEscrowByRequester(ctx, escrow)
	k.setEscrowByProvider(ctx, escrow)
	return nil
}

// GetEscrow retrieves an escrow by ID
func (k Keeper) GetEscrow(ctx sdk.Context, id string) (types.Escrow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(id)
	value := store.Get(key)

	if value == nil {
		return types.Escrow{}, false
	}

	var escrow types.Escrow
	if err := json.Unmarshal(value, &escrow); err != nil {
		return types.Escrow{}, false
	}
	return escrow, true
}

// HasEscrow checks if an escrow exists
func (k Keeper) HasEscrow(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(id)
	return store.Has(key)
}

// DeleteEscrow deletes an escrow from the store
func (k Keeper) DeleteEscrow(ctx sdk.Context, id string) {
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(id)
	store.Delete(key)

	// Delete indexes
	k.deleteEscrowByRequester(ctx, escrow)
	k.deleteEscrowByProvider(ctx, escrow)
}

// setEscrowByRequester sets the escrow by requester index
//
//go:noinline
func (k Keeper) setEscrowByRequester(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByRequesterKey(escrow.Requester, escrow.ID)
	store.Set(key, []byte{1})
}

// deleteEscrowByRequester deletes the escrow by requester index
//
//go:noinline
func (k Keeper) deleteEscrowByRequester(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByRequesterKey(escrow.Requester, escrow.ID)
	store.Delete(key)
}

// setEscrowByProvider sets the escrow by provider index
//
//go:noinline
func (k Keeper) setEscrowByProvider(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByProviderKey(escrow.Provider, escrow.ID)
	store.Set(key, []byte{1})
}

// deleteEscrowByProvider deletes the escrow by provider index
//
//go:noinline
func (k Keeper) deleteEscrowByProvider(ctx sdk.Context, escrow types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowByProviderKey(escrow.Provider, escrow.ID)
	store.Delete(key)
}

// GetEscrowsByRequester returns all escrows by requester
func (k Keeper) GetEscrowsByRequester(ctx sdk.Context, requester string) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.EscrowByRequesterKey, []byte(requester)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close() //nolint:errcheck

	var escrows []types.Escrow
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract escrow ID from key
		escrowID := string(key[len(prefix):])
		escrow, found := k.GetEscrow(ctx, escrowID)
		if found {
			escrows = append(escrows, escrow)
		}
	}

	return escrows
}

// GetEscrowsByProvider returns all escrows by provider
func (k Keeper) GetEscrowsByProvider(ctx sdk.Context, provider string) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.EscrowByProviderKey, []byte(provider)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close() //nolint:errcheck

	var escrows []types.Escrow
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract escrow ID from key
		escrowID := string(key[len(prefix):])
		escrow, found := k.GetEscrow(ctx, escrowID)
		if found {
			escrows = append(escrows, escrow)
		}
	}

	return escrows
}

// GetAllEscrows returns all escrows
func (k Keeper) GetAllEscrows(ctx sdk.Context) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.EscrowKey)
	defer iterator.Close() //nolint:errcheck

	var escrows []types.Escrow
	for ; iterator.Valid(); iterator.Next() {
		var escrow types.Escrow
		if err := json.Unmarshal(iterator.Value(), &escrow); err != nil {
			ctx.Logger().Error("failed to unmarshal escrow", "error", err)
			continue
		}
		escrows = append(escrows, escrow)
	}

	return escrows
}

// CreateEscrow creates a new escrow
func (k Keeper) CreateEscrow(ctx sdk.Context, requester, provider string, amount sdk.Coins, duration time.Duration) (types.Escrow, error) {
	// Generate unique ID
	id := fmt.Sprintf("escrow-%d-%d", ctx.BlockHeight(), ctx.BlockTime().Unix())

	escrow := types.NewEscrow(id, requester, provider, amount, duration)

	// Validate
	if err := escrow.ValidateBasic(); err != nil {
		return types.Escrow{}, err
	}

	// Store escrow
	if err := k.SetEscrow(ctx, *escrow); err != nil {
		return types.Escrow{}, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateEscrow,
			sdk.NewAttribute(types.AttributeKeyEscrowID, id),
			sdk.NewAttribute(types.AttributeKeyRequester, requester),
			sdk.NewAttribute(types.AttributeKeyProvider, provider),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)

	return *escrow, nil
}

// Release releases funds to provider
func (k Keeper) Release(ctx sdk.Context, escrowID string) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}

	if !escrow.CanComplete() {
		return types.ErrInvalidStatus
	}

	// Update status
	escrow.Status = types.EscrowStatusCompleted
	escrow.CompletedAt = ctx.BlockTime().Unix()
	if err := k.SetEscrow(ctx, escrow); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRelease,
			sdk.NewAttribute(types.AttributeKeyEscrowID, escrowID),
			sdk.NewAttribute(types.AttributeKeyProvider, escrow.Provider),
			sdk.NewAttribute(types.AttributeKeyAmount, escrow.Amount.String()),
		),
	)

	return nil
}

// Refund refunds funds to requester (only if expired)
func (k Keeper) Refund(ctx sdk.Context, escrowID string) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}

	if !escrow.CanRefund() {
		return types.ErrInvalidStatus
	}

	// Update status
	escrow.Status = types.EscrowStatusRefunded
	if err := k.SetEscrow(ctx, escrow); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRefund,
			sdk.NewAttribute(types.AttributeKeyEscrowID, escrowID),
			sdk.NewAttribute(types.AttributeKeyRequester, escrow.Requester),
			sdk.NewAttribute(types.AttributeKeyAmount, escrow.Amount.String()),
		),
	)

	return nil
}

// Dispute marks escrow as disputed
func (k Keeper) Dispute(ctx sdk.Context, escrowID string) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}

	if !escrow.CanDispute() {
		return types.ErrInvalidStatus
	}

	// Update status
	escrow.Status = types.EscrowStatusDisputed
	if err := k.SetEscrow(ctx, escrow); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDispute,
			sdk.NewAttribute(types.AttributeKeyEscrowID, escrowID),
			sdk.NewAttribute(types.AttributeKeyRequester, escrow.Requester),
			sdk.NewAttribute(types.AttributeKeyProvider, escrow.Provider),
		),
	)

	return nil
}

// ResolveDispute resolves a disputed escrow with fund allocation
func (k Keeper) ResolveDispute(ctx sdk.Context, escrowID string, allocation types.FundAllocation) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}

	if escrow.Status != types.EscrowStatusDisputed {
		return types.ErrInvalidStatus
	}

	// Validate allocation
	if err := allocation.Validate(escrow.Amount); err != nil {
		return err
	}

	// Update escrow
	escrow.Status = types.EscrowStatusCompleted
	escrow.CompletedAt = ctx.BlockTime().Unix()
	if err := k.SetEscrow(ctx, escrow); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeResolve,
			sdk.NewAttribute(types.AttributeKeyEscrowID, escrowID),
			sdk.NewAttribute(types.AttributeKeyRequesterAmount, allocation.RequesterAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProviderAmount, allocation.ProviderAmount.String()),
		),
	)

	return nil
}
