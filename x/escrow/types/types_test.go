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
	validCoins := []sdk.Coin{sdk.NewCoin("ustt", sdk.NewInt(1000))}

	tests := []struct {
		name    string
		escrow  types.Escrow
		wantErr bool
		errType error
	}{
		{
			name: "valid escrow",
			escrow: types.Escrow{
				Id:        "escrow-1",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty ID",
			escrow: types.Escrow{
				Id:        "",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
			},
			wantErr: true,
			errType: types.ErrInvalidEscrowID,
		},
		{
			name: "invalid - empty requester",
			escrow: types.Escrow{
				Id:        "escrow-1",
				Requester: "",
				Provider:  validAddress2,
				Amount:    validCoins,
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
			},
			wantErr: true,
			errType: types.ErrInvalidRequester,
		},
		{
			name: "invalid - empty provider",
			escrow: types.Escrow{
				Id:        "escrow-1",
				Requester: validAddress,
				Provider:  "",
				Amount:    validCoins,
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
			},
			wantErr: true,
			errType: types.ErrInvalidProvider,
		},
		{
			name: "invalid - zero amount",
			escrow: types.Escrow{
				Id:        "escrow-1",
				Requester: validAddress,
				Provider:  validAddress2,
				Amount:    []sdk.Coin{},
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
			},
			wantErr: true,
			errType: types.ErrInvalidAmount,
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
	require.Equal(t, "escrow-1", escrow.Id)
	require.Equal(t, validAddress, escrow.Requester)
	require.Equal(t, validAddress2, escrow.Provider)
	require.True(t, sdk.Coins(escrow.Amount).IsEqual(validCoins))
	require.Equal(t, types.EscrowStatus_ESCROW_STATUS_PENDING, escrow.Status)
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
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
				ExpiresAt: now + 3600,
			},
			expected: true,
		},
		{
			name: "cannot complete - completed status",
			escrow: types.Escrow{
				Status:    types.EscrowStatus_ESCROW_STATUS_COMPLETED,
				ExpiresAt: now + 3600,
			},
			expected: false,
		},
		{
			name: "cannot complete - expired",
			escrow: types.Escrow{
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
				ExpiresAt: now - 3600,
			},
			expected: false,
		},
		{
			name: "cannot complete - disputed",
			escrow: types.Escrow{
				Status:    types.EscrowStatus_ESCROW_STATUS_DISPUTED,
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
		{"can dispute - pending", types.EscrowStatus_ESCROW_STATUS_PENDING, true},
		{"can dispute - completed", types.EscrowStatus_ESCROW_STATUS_COMPLETED, true},
		{"cannot dispute - disputed", types.EscrowStatus_ESCROW_STATUS_DISPUTED, false},
		{"cannot dispute - refunded", types.EscrowStatus_ESCROW_STATUS_REFUNDED, false},
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
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
				ExpiresAt: now - 3600,
			},
			expected: true,
		},
		{
			name: "cannot refund - not expired",
			escrow: types.Escrow{
				Status:    types.EscrowStatus_ESCROW_STATUS_PENDING,
				ExpiresAt: now + 3600,
			},
			expected: false,
		},
		{
			name: "cannot refund - not pending",
			escrow: types.Escrow{
				Status:    types.EscrowStatus_ESCROW_STATUS_COMPLETED,
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

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Escrows)
	require.Empty(t, genesis.Escrows)
}

func TestValidateGenesis(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	validAddress2 := sdk.AccAddress([]byte("test_address_2")).String()
	validCoins := []sdk.Coin{sdk.NewCoin("ustt", sdk.NewInt(1000))}

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
					{Id: "escrow-1", Requester: validAddress, Provider: validAddress2, Amount: validCoins, Status: types.EscrowStatus_ESCROW_STATUS_PENDING},
					{Id: "escrow-2", Requester: validAddress2, Provider: validAddress, Amount: validCoins, Status: types.EscrowStatus_ESCROW_STATUS_COMPLETED},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - duplicate escrow IDs",
			data: types.GenesisState{
				Escrows: []types.Escrow{
					{Id: "escrow-1", Requester: validAddress, Provider: validAddress2, Amount: validCoins, Status: types.EscrowStatus_ESCROW_STATUS_PENDING},
					{Id: "escrow-1", Requester: validAddress2, Provider: validAddress, Amount: validCoins, Status: types.EscrowStatus_ESCROW_STATUS_PENDING},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid escrow",
			data: types.GenesisState{
				Escrows: []types.Escrow{
					{Id: "", Requester: validAddress, Provider: validAddress2, Amount: validCoins, Status: types.EscrowStatus_ESCROW_STATUS_PENDING},
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
