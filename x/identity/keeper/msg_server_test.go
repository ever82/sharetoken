package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"sharetoken/testutil/sample"
	"sharetoken/x/identity/types"
)

// setupMsgServer sets up the test environment and returns a MsgServer and context
func setupMsgServer(t testing.TB) (types.MsgServer, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"IdentityParams",
	)

	// Create keeper without account and bank keepers for these tests
	k := NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		nil, // accountKeeper
		nil, // bankKeeper
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return NewMsgServerImpl(k), ctx
}

// Test helpers
func createTestAddress() string {
	return sample.AccAddress()
}

func createSecondAddress() string {
	return sample.AccAddress()
}

func createThirdAddress() string {
	return sample.AccAddress()
}

// Test MsgServer initialization
func TestMsgServer_Initialization(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	require.NotNil(t, msgServer)
	require.NotNil(t, ctx)
}

// RegisterIdentity Tests
func TestMsgServer_RegisterIdentity_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	msg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}

	resp, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.MerkleRoot)
}

func TestMsgServer_RegisterIdentity_AlreadyExists(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// First registration
	msg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// Second registration should fail
	resp, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "already exists")
}

func TestMsgServer_RegisterIdentity_DuplicateDID(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)

	// First registration with DID
	msg1 := &types.MsgRegisterIdentity{
		Address:      createTestAddress(),
		Did:          "did:sharetoken:shared",
		MetadataHash: "QmHash1",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), msg1)
	require.NoError(t, err)

	// Second registration with same DID should fail
	msg2 := &types.MsgRegisterIdentity{
		Address:      createSecondAddress(),
		Did:          "did:sharetoken:shared",
		MetadataHash: "QmHash2",
	}
	resp, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), msg2)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "already registered")
}

// VerifyIdentity Tests
func TestMsgServer_VerifyIdentity_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// First register identity
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Now verify
	verifyMsg := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "github",
		VerificationHash: "verify_hash_123",
		Proof:            "proof_token_123",
	}

	resp, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.IsVerified)
	require.NotEmpty(t, resp.UpdatedMerkleRoot)
}

func TestMsgServer_VerifyIdentity_NotFound(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	verifyMsg := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "github",
		VerificationHash: "verify_hash_123",
		Proof:            "proof_token_123",
	}

	resp, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "not found")
}

func TestMsgServer_VerifyIdentity_InvalidProvider(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// First register identity
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Verify with invalid provider
	verifyMsg := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "invalid_provider",
		VerificationHash: "verify_hash_123",
		Proof:            "proof_token_123",
	}

	resp, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "invalid verification provider")
}

func TestMsgServer_VerifyIdentity_DuplicateProvider(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// First register identity
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// First verification
	verifyMsg1 := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "github",
		VerificationHash: "verify_hash_123",
		Proof:            "proof_token_123",
	}
	_, err = msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg1)
	require.NoError(t, err)

	// Second verification with same provider should succeed (as provider is marked as used)
	verifyMsg2 := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "wechat",
		VerificationHash: "verify_hash_456",
		Proof:            "proof_token_456",
	}
	resp, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg2)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestMsgServer_VerifyIdentity_EmptyProof(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// First register identity
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Verify with empty proof
	verifyMsg := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "github",
		VerificationHash: "verify_hash_123",
		Proof:            "",
	}

	resp, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "proof is required")
}

// UpdateLimitConfig Tests
func TestMsgServer_UpdateLimitConfig_Unauthorized(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// First register identity
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:abc123",
		MetadataHash: "QmHash123",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Update limit config
	limitConfig := types.NewLimitConfig(address)
	msg := &types.MsgUpdateLimitConfig{
		TargetAddress: address,
		NewConfig:     limitConfig,
	}

	resp, err := msgServer.UpdateLimitConfig(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "only authority can update limits")
}

// ResetDailyLimits Tests
func TestMsgServer_ResetDailyLimits_Unauthorized(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)

	msg := &types.MsgResetDailyLimits{
		Authority: createTestAddress(),
	}

	resp, err := msgServer.ResetDailyLimits(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "only authority can reset limits")
}

// UpdateParams Tests
func TestMsgServer_UpdateParams_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)

	params := types.DefaultParams()
	params.VerificationRequired = true

	msg := &types.MsgUpdateParams{
		Authority: "sharetoken10d07y265gmmuvt4z0w9aw880jnsr700jfw7aah", // gov address
		Params:    params,
	}

	resp, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestMsgServer_UpdateParams_Unauthorized(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)

	params := types.DefaultParams()

	msg := &types.MsgUpdateParams{
		Authority: createSecondAddress(), // not gov address
		Params:    params,
	}

	resp, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "only governance can update params")
}

// Integration Tests
func TestMsgServer_FullIdentityLifecycle(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// 1. Register identity
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:lifecycle123",
		MetadataHash: "QmLifecycleHash",
	}
	resp1, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.NotEmpty(t, resp1.MerkleRoot)
	merkleRoot1 := resp1.MerkleRoot

	// 2. Verify identity
	verifyMsg := &types.MsgVerifyIdentity{
		Address:          address,
		Provider:         "google",
		VerificationHash: "lifecycle_verify_hash",
		Proof:            "lifecycle_proof_token",
	}
	resp2, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg)
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.True(t, resp2.IsVerified)
	require.NotEmpty(t, resp2.UpdatedMerkleRoot)
	require.NotEqual(t, merkleRoot1, resp2.UpdatedMerkleRoot)
}

func TestMsgServer_MultipleIdentities(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)

	for i := 0; i < 3; i++ {
		addr := createTestAddress()
		msg := &types.MsgRegisterIdentity{
			Address:      addr,
			Did:          "did:sharetoken:user" + string(rune('0'+i)),
			MetadataHash: "QmHash" + string(rune('0'+i)),
		}
		resp, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.MerkleRoot)
	}
}

func TestMsgServer_VerifyWithDifferentProviders(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	address := createTestAddress()

	// Register
	registerMsg := &types.MsgRegisterIdentity{
		Address:      address,
		Did:          "did:sharetoken:multi",
		MetadataHash: "QmMultiHash",
	}
	_, err := msgServer.RegisterIdentity(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Verify with first provider
	providers := []string{"github", "wechat", "google"}
	for _, provider := range providers {
		verifyMsg := &types.MsgVerifyIdentity{
			Address:          address,
			Provider:         provider,
			VerificationHash: "verify_" + provider,
			Proof:            "proof_" + provider,
		}
		resp, err := msgServer.VerifyIdentity(sdk.WrapSDKContext(ctx), verifyMsg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.True(t, resp.IsVerified)
	}
}
