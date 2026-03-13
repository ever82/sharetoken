package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/llmcustody/types"
)

// setupMsgServer sets up the test environment and returns a MsgServer and context
func setupMsgServer(t testing.TB) (types.MsgServer, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := NewKeeper(cdc, storeKey)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return NewMsgServerImpl(k), ctx
}

// Test helpers
func createTestOwner() string {
	return "sharetoken10d07y265gmmuvt4z0w9aw880jnsr700jfw7aah"
}

func createSecondOwner() string {
	return "sharetoken1v9l7k6xxm8dqq6rhy5u5q2xjwg5u5q2xjwg5u5q"
}

func createEncryptedKey() []byte {
	return []byte("encrypted-api-key-data-12345")
}

func createAccessRules() []types.AccessRule {
	return []types.AccessRule{
		{
			ServiceID:   "service-1",
			RateLimit:   100,
			MaxRequests: 1000,
			PricePerReq: 10,
		},
		{
			ServiceID:   "service-2",
			RateLimit:   200,
			MaxRequests: 2000,
			PricePerReq: 20,
		},
	}
}

// Test MsgServer initialization
func TestMsgServer_Initialization(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	require.NotNil(t, msgServer)
	require.NotNil(t, ctx)
}

// RegisterAPIKey Tests
func TestMsgServer_RegisterAPIKey_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()
	accessRules := createAccessRules()

	msg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  accessRules,
	}

	resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.APIKeyID)
}

func TestMsgServer_RegisterAPIKey_MultipleProviders(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	providers := []string{"openai", "anthropic"}

	for _, provider := range providers {
		msg := &types.MsgRegisterAPIKey{
			Owner:        owner,
			Provider:     provider,
			EncryptedKey: encryptedKey,
			AccessRules:  createAccessRules(),
		}

		resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.APIKeyID)
	}
}

func TestMsgServer_RegisterAPIKey_InvalidProvider(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	msg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "invalid_provider",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}

	resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestMsgServer_RegisterAPIKey_EmptyOwner(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	encryptedKey := createEncryptedKey()

	msg := &types.MsgRegisterAPIKey{
		Owner:        "",
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}

	resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestMsgServer_RegisterAPIKey_EmptyEncryptedKey(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()

	msg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: []byte{},
		AccessRules:  createAccessRules(),
	}

	resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestMsgServer_RegisterAPIKey_NoAccessRules(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	msg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  []types.AccessRule{},
	}

	resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.APIKeyID)
}

func TestMsgServer_RegisterAPIKey_GeneratesUniqueIDs(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	ids := make(map[string]bool)

	for i := 0; i < 10; i++ {
		msg := &types.MsgRegisterAPIKey{
			Owner:        owner,
			Provider:     "openai",
			EncryptedKey: encryptedKey,
			AccessRules:  createAccessRules(),
		}

		resp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.APIKeyID)

		// Check ID is unique
		require.False(t, ids[resp.APIKeyID], "API Key ID should be unique")
		ids[resp.APIKeyID] = true
	}
}

// UpdateAPIKey Tests
func TestMsgServer_UpdateAPIKey_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()
	accessRules := createAccessRules()

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  accessRules,
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Now update the API key
	newAccessRules := []types.AccessRule{
		{
			ServiceID:   "service-3",
			RateLimit:   300,
			MaxRequests: 3000,
			PricePerReq: 30,
		},
	}

	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    registerResp.APIKeyID,
		AccessRules: newAccessRules,
		Active:      false, // Deactivate
	}

	resp, err := msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestMsgServer_UpdateAPIKey_NotFound(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()

	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    "non-existent-key",
		AccessRules: createAccessRules(),
		Active:      true,
	}

	resp, err := msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, types.ErrAPIKeyNotFound, err)
}

func TestMsgServer_UpdateAPIKey_Unauthorized(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	otherOwner := createSecondOwner()
	encryptedKey := createEncryptedKey()

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Try to update with different owner
	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       otherOwner,
		APIKeyID:    registerResp.APIKeyID,
		AccessRules: createAccessRules(),
		Active:      true,
	}

	resp, err := msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, types.ErrUnauthorized, err)
}

func TestMsgServer_UpdateAPIKey_EmptyAccessRules(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Update with empty access rules (should still work, just keeps existing)
	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    registerResp.APIKeyID,
		AccessRules: []types.AccessRule{},
		Active:      true,
	}

	resp, err := msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// RevokeAPIKey Tests
func TestMsgServer_RevokeAPIKey_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Revoke the API key
	revokeMsg := &types.MsgRevokeAPIKey{
		Owner:    owner,
		APIKeyID: registerResp.APIKeyID,
	}

	resp, err := msgServer.RevokeAPIKey(sdk.WrapSDKContext(ctx), revokeMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestMsgServer_RevokeAPIKey_NotFound(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()

	revokeMsg := &types.MsgRevokeAPIKey{
		Owner:    owner,
		APIKeyID: "non-existent-key",
	}

	resp, err := msgServer.RevokeAPIKey(sdk.WrapSDKContext(ctx), revokeMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, types.ErrAPIKeyNotFound, err)
}

func TestMsgServer_RevokeAPIKey_Unauthorized(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	otherOwner := createSecondOwner()
	encryptedKey := createEncryptedKey()

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Try to revoke with different owner
	revokeMsg := &types.MsgRevokeAPIKey{
		Owner:    otherOwner,
		APIKeyID: registerResp.APIKeyID,
	}

	resp, err := msgServer.RevokeAPIKey(sdk.WrapSDKContext(ctx), revokeMsg)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, types.ErrUnauthorized, err)
}

// RecordUsage Tests
func TestMsgServer_RecordUsage_Success(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()
	accessRules := []types.AccessRule{
		{
			ServiceID:   "service-1",
			RateLimit:   100,
			MaxRequests: 1000,
			PricePerReq: 10,
		},
	}

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  accessRules,
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Record usage
	usageMsg := &types.MsgRecordUsage{
		APIKeyID:     registerResp.APIKeyID,
		ServiceID:    "service-1",
		RequestCount: 1,
		TokenCount:   100,
		Cost:         10,
	}

	resp, err := msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestMsgServer_RecordUsage_InvalidService(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()
	accessRules := []types.AccessRule{
		{
			ServiceID:   "service-1",
			RateLimit:   100,
			MaxRequests: 1000,
			PricePerReq: 10,
		},
	}

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  accessRules,
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Try to record usage for invalid service
	usageMsg := &types.MsgRecordUsage{
		APIKeyID:     registerResp.APIKeyID,
		ServiceID:    "invalid-service",
		RequestCount: 1,
		TokenCount:   100,
		Cost:         10,
	}

	resp, err := msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestMsgServer_RecordUsage_KeyNotFound(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)

	usageMsg := &types.MsgRecordUsage{
		APIKeyID:     "non-existent-key",
		ServiceID:    "service-1",
		RequestCount: 1,
		TokenCount:   100,
		Cost:         10,
	}

	resp, err := msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestMsgServer_RecordUsage_InactiveKey(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()
	accessRules := []types.AccessRule{
		{
			ServiceID:   "service-1",
			RateLimit:   100,
			MaxRequests: 1000,
			PricePerReq: 10,
		},
	}

	// First register an API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  accessRules,
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)

	// Deactivate the key
	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    registerResp.APIKeyID,
		AccessRules: accessRules,
		Active:      false,
	}
	_, err = msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.NoError(t, err)

	// Try to record usage
	usageMsg := &types.MsgRecordUsage{
		APIKeyID:     registerResp.APIKeyID,
		ServiceID:    "service-1",
		RequestCount: 1,
		TokenCount:   100,
		Cost:         10,
	}

	resp, err := msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg)
	require.Error(t, err)
	require.Nil(t, resp)
}

// Full Workflow Integration Tests
func TestMsgServer_FullAPIKeyLifecycle(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner := createTestOwner()
	encryptedKey := createEncryptedKey()

	// 1. Register API key
	registerMsg := &types.MsgRegisterAPIKey{
		Owner:        owner,
		Provider:     "openai",
		EncryptedKey: encryptedKey,
		AccessRules:  createAccessRules(),
	}
	registerResp, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg)
	require.NoError(t, err)
	require.NotNil(t, registerResp)
	apiKeyID := registerResp.APIKeyID

	// 2. Record some usage
	usageMsg := &types.MsgRecordUsage{
		APIKeyID:     apiKeyID,
		ServiceID:    "service-1",
		RequestCount: 5,
		TokenCount:   500,
		Cost:         50,
	}
	_, err = msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg)
	require.NoError(t, err)

	// 3. Update API key (change access rules and deactivate)
	newAccessRules := []types.AccessRule{
		{
			ServiceID:   "service-3",
			RateLimit:   500,
			MaxRequests: 5000,
			PricePerReq: 50,
		},
	}
	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    apiKeyID,
		AccessRules: newAccessRules,
		Active:      false,
	}
	_, err = msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.NoError(t, err)

	// 4. Reactivate the key
	updateMsg2 := &types.MsgUpdateAPIKey{
		Owner:       owner,
		APIKeyID:    apiKeyID,
		AccessRules: newAccessRules,
		Active:      true,
	}
	_, err = msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg2)
	require.NoError(t, err)

	// 5. Record usage for new service
	usageMsg2 := &types.MsgRecordUsage{
		APIKeyID:     apiKeyID,
		ServiceID:    "service-3",
		RequestCount: 10,
		TokenCount:   1000,
		Cost:         100,
	}
	_, err = msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg2)
	require.NoError(t, err)

	// 6. Revoke the API key
	revokeMsg := &types.MsgRevokeAPIKey{
		Owner:    owner,
		APIKeyID: apiKeyID,
	}
	_, err = msgServer.RevokeAPIKey(sdk.WrapSDKContext(ctx), revokeMsg)
	require.NoError(t, err)

	// 7. Try to record usage after revocation (should fail)
	usageMsg3 := &types.MsgRecordUsage{
		APIKeyID:     apiKeyID,
		ServiceID:    "service-3",
		RequestCount: 1,
		TokenCount:   100,
		Cost:         10,
	}
	resp, err := msgServer.RecordUsage(sdk.WrapSDKContext(ctx), usageMsg3)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestMsgServer_MultipleOwners(t *testing.T) {
	msgServer, ctx := setupMsgServer(t)
	owner1 := createTestOwner()
	owner2 := createSecondOwner()

	// Owner 1 registers a key
	registerMsg1 := &types.MsgRegisterAPIKey{
		Owner:        owner1,
		Provider:     "openai",
		EncryptedKey: createEncryptedKey(),
		AccessRules:  createAccessRules(),
	}
	resp1, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg1)
	require.NoError(t, err)

	// Owner 2 registers a key
	registerMsg2 := &types.MsgRegisterAPIKey{
		Owner:        owner2,
		Provider:     "anthropic",
		EncryptedKey: createEncryptedKey(),
		AccessRules:  createAccessRules(),
	}
	resp2, err := msgServer.RegisterAPIKey(sdk.WrapSDKContext(ctx), registerMsg2)
	require.NoError(t, err)

	// Verify IDs are different
	require.NotEqual(t, resp1.APIKeyID, resp2.APIKeyID)

	// Owner 1 cannot update owner 2's key
	updateMsg := &types.MsgUpdateAPIKey{
		Owner:       owner1,
		APIKeyID:    resp2.APIKeyID,
		AccessRules: createAccessRules(),
		Active:      false,
	}
	_, err = msgServer.UpdateAPIKey(sdk.WrapSDKContext(ctx), updateMsg)
	require.Error(t, err)
	require.Equal(t, types.ErrUnauthorized, err)

	// Owner 2 cannot revoke owner 1's key
	revokeMsg := &types.MsgRevokeAPIKey{
		Owner:    owner2,
		APIKeyID: resp1.APIKeyID,
	}
	_, err = msgServer.RevokeAPIKey(sdk.WrapSDKContext(ctx), revokeMsg)
	require.Error(t, err)
	require.Equal(t, types.ErrUnauthorized, err)
}
