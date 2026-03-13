package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"sharetoken/x/escrow/types"
	identitytypes "sharetoken/x/identity/types"
)

func TestEscrow_ValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validAddress2 := sdk.AccAddress([]byte("test_address_2")).String()
	validCoins := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))

	tests := []struct {
		name    string
		escrow  types.Escrow
		wantErr bool
		errType error
	}{
		{
			name: "valid escrow",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatusPending,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty ID",
			escrow: types.Escrow{
				ID:        "",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrInvalidEscrowID,
		},
		{
			name: "invalid - empty requester",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: "",
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrUnauthorized,
		},
		{
			name: "invalid - malformed requester address",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: "invalid_address",
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrUnauthorized,
		},
		{
			name: "invalid - empty provider",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: validAddress,
				Provider:  "",
				Amount:    validCoins,
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrUnauthorized,
		},
		{
			name: "invalid - malformed provider address",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: validAddress,
				Provider:  "invalid_address",
				Amount:    validCoins,
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrUnauthorized,
		},
		{
			name: "invalid - zero amount",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    sdk.NewCoins(),
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrInvalidAmount,
		},
		{
			name: "invalid - invalid amount",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    sdk.Coins{{Denom: "ustt", Amount: sdk.NewInt(-1)}},
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
			errType: types.ErrInvalidAmount,
		},
		{
			name: "invalid - empty status",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    "",
			},
			wantErr: true,
			errType: types.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.escrow.ValidateBasic()
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

func TestNewEscrow(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validAddress2 := sdk.AccAddress([]byte("test_address_2")).String()
	validCoins := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))
	duration := time.Duration(identitytypes.DefaultEscrowDurationHours) * time.Hour

	escrow := types.NewEscrow("escrow-1", validAddress, validAddress2, validCoins, duration)

	require.NotNil(t, escrow)
	require.Equal(t, "escrow-1", escrow.ID)
	require.Equal(t, validAddress, escrow.Requester)
	require.Equal(t, validAddress2, escrow.Provider)
	require.True(t, escrow.Amount.IsEqual(validCoins))
	require.Equal(t, types.EscrowStatusPending, escrow.Status)
	require.Equal(t, validAddress, escrow.RefundAddress)
	require.NotZero(t, escrow.CreatedAt)
	require.NotZero(t, escrow.ExpiresAt)
	require.True(t, escrow.ExpiresAt > escrow.CreatedAt)
}

func TestEscrow_IsExpired(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name     string
		escrow   types.Escrow
		expected bool
	}{
		{
			name: "not expired",
			escrow: types.Escrow{
				ExpiresAt: now + 3600, // 1 hour from now
			},
			expected: false,
		},
		{
			name: "expired",
			escrow: types.Escrow{
				ExpiresAt: now - 3600, // 1 hour ago
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.escrow.IsExpired()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEscrow_CanComplete(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name     string
		escrow   types.Escrow
		expected bool
	}{
		{
			name: "can complete - pending and not expired",
			escrow: types.Escrow{
				Status:    types.EscrowStatusPending,
				ExpiresAt: now + 3600,
			},
			expected: true,
		},
		{
			name: "cannot complete - completed status",
			escrow: types.Escrow{
				Status:    types.EscrowStatusCompleted,
				ExpiresAt: now + 3600,
			},
			expected: false,
		},
		{
			name: "cannot complete - expired",
			escrow: types.Escrow{
				Status:    types.EscrowStatusPending,
				ExpiresAt: now - 3600,
			},
			expected: false,
		},
		{
			name: "cannot complete - disputed",
			escrow: types.Escrow{
				Status:    types.EscrowStatusDisputed,
				ExpiresAt: now + 3600,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.escrow.CanComplete()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEscrow_CanDispute(t *testing.T) {
	tests := []struct {
		name     string
		status   types.EscrowStatus
		expected bool
	}{
		{"can dispute - pending", types.EscrowStatusPending, true},
		{"can dispute - completed", types.EscrowStatusCompleted, true},
		{"cannot dispute - disputed", types.EscrowStatusDisputed, false},
		{"cannot dispute - refunded", types.EscrowStatusRefunded, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			escrow := types.Escrow{Status: tt.status}
			result := escrow.CanDispute()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEscrow_CanRefund(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name     string
		escrow   types.Escrow
		expected bool
	}{
		{
			name: "can refund - pending and expired",
			escrow: types.Escrow{
				Status:    types.EscrowStatusPending,
				ExpiresAt: now - 3600,
			},
			expected: true,
		},
		{
			name: "cannot refund - not expired",
			escrow: types.Escrow{
				Status:    types.EscrowStatusPending,
				ExpiresAt: now + 3600,
			},
			expected: false,
		},
		{
			name: "cannot refund - not pending",
			escrow: types.Escrow{
				Status:    types.EscrowStatusCompleted,
				ExpiresAt: now - 3600,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.escrow.CanRefund()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFundAllocation_Validate(t *testing.T) {
	total := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))

	tests := []struct {
		name      string
		allocation types.FundAllocation
		wantErr   bool
	}{
		{
			name: "valid allocation",
			allocation: types.FundAllocation{
				RequesterAmount: sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(400))),
				ProviderAmount:  sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(600))),
			},
			wantErr: false,
		},
		{
			name: "invalid - sum exceeds total",
			allocation: types.FundAllocation{
				RequesterAmount: sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(600))),
				ProviderAmount:  sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(600))),
			},
			wantErr: true,
		},
		{
			name: "invalid - sum less than total",
			allocation: types.FundAllocation{
				RequesterAmount: sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(300))),
				ProviderAmount:  sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(300))),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.allocation.Validate(total)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Escrows)
	require.Empty(t, genesis.Escrows)
}

func TestValidateGenesis(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validAddress2 := sdk.AccAddress([]byte("test_address_2")).String()
	validCoins := sdk.NewCoins(sdk.NewCoin("ustt", sdk.NewInt(1000)))

	tests := []struct {
		name    string
		data    types.GenesisState
		wantErr bool
	}{
		{
			name:    "valid genesis with default",
			data:    *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "valid genesis with escrows",
			data: types.GenesisState{
				Escrows: []types.Escrow{
					{ID: "escrow-1", Requester: validAddress, Provider: validAddress2, Amount: validCoins, Status: types.EscrowStatusPending},
					{ID: "escrow-2", Requester: validAddress2, Provider: validAddress, Amount: validCoins, Status: types.EscrowStatusCompleted},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate escrow IDs",
			data: types.GenesisState{
				Escrows: []types.Escrow{
					{ID: "escrow-1", Requester: validAddress, Provider: validAddress2, Amount: validCoins, Status: types.EscrowStatusPending},
					{ID: "escrow-1", Requester: validAddress2, Provider: validAddress, Amount: validCoins, Status: types.EscrowStatusPending},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid escrow",
			data: types.GenesisState{
				Escrows: []types.Escrow{
					{ID: "", Requester: validAddress, Provider: validAddress2, Amount: validCoins, Status: types.EscrowStatusPending},
				},
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
