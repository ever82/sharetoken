package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"sharetoken/x/identity/types"
)

func TestIdentity_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validDID := "did:sharetoken:abc123"

	tests := []struct {
		name      string
		identity  types.Identity
		wantErr   bool
		errType   error
	}{
		{
			name: "valid identity with address and DID",
			identity: types.Identity{
				Address: validAddress,
				Did:     validDID,
			},
			wantErr: false,
		},
		{
			name: "valid identity with address only",
			identity: types.Identity{
				Address: validAddress,
				Did:     "",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty address",
			identity: types.Identity{
				Address: "",
				Did:     validDID,
			},
			wantErr: true,
			errType: types.ErrInvalidAddress,
		},
		{
			name: "invalid - malformed address",
			identity: types.Identity{
				Address: "invalid_address",
				Did:     validDID,
			},
			wantErr: true,
			errType: types.ErrInvalidAddress,
		},
		{
			name: "invalid - malformed DID",
			identity: types.Identity{
				Address: validAddress,
				Did:     "not_a_did",
			},
			wantErr: true,
			errType: types.ErrInvalidDID,
		},
		{
			name: "invalid - DID too short",
			identity: types.Identity{
				Address: validAddress,
				Did:     "did",
			},
			wantErr: true,
			errType: types.ErrInvalidDID,
		},
		{
			name: "valid - minimal DID format",
			identity: types.Identity{
				Address: validAddress,
				Did:     "did:method:test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.identity.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					require.ErrorIs(t, err, tt.errType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewIdentity(t *testing.T) {
	ctx := sdk.Context{} // Note: We can't easily set block time in unit test
	address := sdk.AccAddress([]byte("test_address_1")).String()
	did := "did:sharetoken:test"

	identity := types.NewIdentity(ctx, address, did)

	require.NotNil(t, identity)
	require.Equal(t, address, identity.Address)
	require.Equal(t, did, identity.Did)
	require.False(t, identity.IsVerified)
}

func TestIsValidProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		want     bool
	}{
		{"valid - wechat", "wechat", true},
		{"valid - github", "github", true},
		{"valid - google", "google", true},
		{"invalid - unknown", "facebook", false},
		{"invalid - empty", "", false},
		{"invalid - case sensitive", "WeChat", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.IsValidProvider(tt.provider)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetVerificationHash(t *testing.T) {
	provider := "wechat"
	providerID := "user123"
	timestamp := "1234567890"

	hash1 := types.GetVerificationHash(provider, providerID, timestamp)
	hash2 := types.GetVerificationHash(provider, providerID, timestamp)

	// Same inputs should produce same hash
	require.Equal(t, hash1, hash2)

	// Different inputs should produce different hash
	hash3 := types.GetVerificationHash(provider, providerID, "different_timestamp")
	require.NotEqual(t, hash1, hash3)
}

func TestLimitConfig_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	invalidAddress := "invalid_address"

	tests := []struct {
		name    string
		config  types.LimitConfig
		wantErr bool
	}{
		{
			name: "valid limit config",
			config: types.LimitConfig{
				Address: validAddress,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty address",
			config: types.LimitConfig{
				Address: "",
			},
			wantErr: true,
		},
		{
			name: "invalid - malformed address",
			config: types.LimitConfig{
				Address: invalidAddress,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewLimitConfig(t *testing.T) {
	address := sdk.AccAddress([]byte("test_address_1")).String()

	config := types.NewLimitConfig(address)

	require.NotNil(t, config)
	require.Equal(t, address, config.Address)
	require.NotZero(t, config.UpdatedAt)
}

func TestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  types.Params
		wantErr bool
	}{
		{
			name: "valid params with providers",
			params: types.NewParams(false, []string{"wechat", "github"}),
			wantErr: false,
		},
		{
			name: "invalid - empty providers",
			params: types.NewParams(false, []string{}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultParams(t *testing.T) {
	params := types.DefaultParams()

	require.NotEmpty(t, params.AllowedProviders)
	require.Contains(t, params.AllowedProviders, "wechat")
	require.Contains(t, params.AllowedProviders, "github")
	require.Contains(t, params.AllowedProviders, "google")
	require.False(t, params.VerificationRequired)
}

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Params)
	require.Empty(t, genesis.Identities)
	require.Empty(t, genesis.LimitConfigs)
}

func TestValidateGenesis(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validAddress2 := sdk.AccAddress([]byte("test_address_2")).String()

	tests := []struct {
		name    string
		data    types.GenesisState
		wantErr bool
	}{
		{
			name: "valid genesis with default",
			data: *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "valid genesis with identities",
			data: types.GenesisState{
				Params: types.DefaultParams(),
				Identities: []types.Identity{
					{Address: validAddress},
					{Address: validAddress2},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate addresses",
			data: types.GenesisState{
				Params: types.DefaultParams(),
				Identities: []types.Identity{
					{Address: validAddress},
					{Address: validAddress},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid identity",
			data: types.GenesisState{
				Params: types.DefaultParams(),
				Identities: []types.Identity{
					{Address: ""},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - params with no providers",
			data: types.GenesisState{
				Params: types.Params{
					AllowedProviders: []string{},
				},
				Identities: []types.Identity{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateGenesis(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLimitError_Error(t *testing.T) {
	err := types.LimitError{
		LimitType:   "transaction_daily",
		Current:     "1000",
		Max:         "500",
		Description: "transaction would exceed daily limit",
	}

	result := err.Error()
	require.Contains(t, result, "transaction_daily")
	require.Contains(t, result, "1000")
	require.Contains(t, result, "500")
}

// Message Tests

func TestMsgRegisterIdentity_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()

	tests := []struct {
		name    string
		msg     types.MsgRegisterIdentity
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgRegisterIdentity{
				Address: validAddress,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty address",
			msg: types.MsgRegisterIdentity{
				Address: "",
			},
			wantErr: true,
		},
		{
			name: "invalid - malformed address",
			msg: types.MsgRegisterIdentity{
				Address: "invalid_address",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgVerifyIdentity_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()

	tests := []struct {
		name    string
		msg     types.MsgVerifyIdentity
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgVerifyIdentity{
				Address:  validAddress,
				Provider: "wechat",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty address",
			msg: types.MsgVerifyIdentity{
				Address:  "",
				Provider: "wechat",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty provider",
			msg: types.MsgVerifyIdentity{
				Address:  validAddress,
				Provider: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateLimitConfig_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()

	tests := []struct {
		name    string
		msg     types.MsgUpdateLimitConfig
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateLimitConfig{
				TargetAddress: validAddress,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty target address",
			msg: types.MsgUpdateLimitConfig{
				TargetAddress: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgResetDailyLimits_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()

	tests := []struct {
		name    string
		msg     types.MsgResetDailyLimits
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgResetDailyLimits{
				Authority: validAddress,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty authority",
			msg: types.MsgResetDailyLimits{
				Authority: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()

	tests := []struct {
		name    string
		msg     types.MsgUpdateParams
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateParams{
				Authority: validAddress,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty authority",
			msg: types.MsgUpdateParams{
				Authority: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLimitConfig_CheckTransactionLimit(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	config := types.NewLimitConfig(validAddress)

	tests := []struct {
		name    string
		amount  sdk.Coin
		wantErr bool
	}{
		{
			name:    "valid amount below limit",
			amount:  sdk.NewCoin("ustt", sdk.NewInt(100)),
			wantErr: false,
		},
		{
			name:    "amount at limit boundary",
			amount:  sdk.NewCoin("ustt", sdk.NewInt(1000000000)),
			wantErr: true, // IsGTE check
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.CheckTransactionLimit(tt.amount)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
