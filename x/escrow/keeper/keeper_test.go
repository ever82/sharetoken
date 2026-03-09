package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/escrow/types"
)

const (
	testRequester = "sharetoken10d07y265gmmuvt4z0w9aw880jnsr700jfw7aah"
	testProvider  = "sharetoken1xy50nxn6jxplnp8zr9f6ph2h4h5j2q8w7z9abc"
)

func TestNewEscrow(t *testing.T) {
	id := "escrow-1"
	amount := sdk.NewCoins(sdk.NewInt64Coin("ustt", 1000000))
	duration := 24 * time.Hour

	escrow := types.NewEscrow(id, testRequester, testProvider, amount, duration)

	require.Equal(t, id, escrow.ID)
	require.Equal(t, testRequester, escrow.Requester)
	require.Equal(t, testProvider, escrow.Provider)
	require.True(t, escrow.Amount.IsEqual(amount))
	require.Equal(t, types.EscrowStatusPending, escrow.Status)
	require.False(t, escrow.IsExpired())
	require.True(t, escrow.CanComplete())
}

func TestEscrowValidation(t *testing.T) {
	tests := []struct {
		name    string
		escrow  types.Escrow
		wantErr bool
	}{
		{
			name: "missing ID",
			escrow: types.Escrow{
				Requester: testRequester,
				Provider:  testProvider,
				Amount:    sdk.NewCoins(sdk.NewInt64Coin("ustt", 1000000)),
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			escrow: types.Escrow{
				ID:        "escrow-1",
				Requester: testRequester,
				Provider:  testProvider,
				Amount:    sdk.Coins{},
				Status:    types.EscrowStatusPending,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.escrow.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEscrowStatus(t *testing.T) {
	escrow := types.NewEscrow("1", testRequester, testProvider, sdk.NewCoins(sdk.NewInt64Coin("ustt", 100)), time.Hour)

	// Test initial status
	require.Equal(t, types.EscrowStatusPending, escrow.Status)
	require.True(t, escrow.CanComplete())
	require.True(t, escrow.CanDispute())
	require.False(t, escrow.CanRefund())

	// Test disputed status
	escrow.Status = types.EscrowStatusDisputed
	require.False(t, escrow.CanComplete())
	require.False(t, escrow.CanDispute())

	// Test completed status
	escrow.Status = types.EscrowStatusCompleted
	require.False(t, escrow.CanComplete())
	require.True(t, escrow.CanDispute()) // Can dispute after completion
}

func TestEscrowExpiration(t *testing.T) {
	// Create an escrow that expired 1 hour ago
	escrow := types.Escrow{
		ID:        "expired-escrow",
		Requester: testRequester,
		Provider:  testProvider,
		Amount:    sdk.NewCoins(sdk.NewInt64Coin("ustt", 1000000)),
		Status:    types.EscrowStatusPending,
		CreatedAt: time.Now().Add(-2 * time.Hour).Unix(),
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
	}

	require.True(t, escrow.IsExpired())
	require.False(t, escrow.CanComplete())
	require.True(t, escrow.CanRefund())
}

func TestFundAllocation(t *testing.T) {
	total := sdk.NewCoins(sdk.NewInt64Coin("ustt", 1000000))

	allocation := types.FundAllocation{
		RequesterAmount: sdk.NewCoins(sdk.NewInt64Coin("ustt", 300000)),
		ProviderAmount:  sdk.NewCoins(sdk.NewInt64Coin("ustt", 700000)),
	}

	err := allocation.Validate(total)
	require.NoError(t, err)

	// Test invalid allocation
	invalidAllocation := types.FundAllocation{
		RequesterAmount: sdk.NewCoins(sdk.NewInt64Coin("ustt", 500000)),
		ProviderAmount:  sdk.NewCoins(sdk.NewInt64Coin("ustt", 600000)),
	}

	err = invalidAllocation.Validate(total)
	require.Error(t, err)
}

func TestGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()
	require.NotNil(t, genesis)
	require.Empty(t, genesis.Escrows)

	// Test validation
	err := types.ValidateGenesis(*genesis)
	require.NoError(t, err)

	// Test with duplicate IDs
	genesisWithDups := types.GenesisState{
		Escrows: []types.Escrow{
			{ID: "1", Requester: testRequester, Provider: testProvider, Amount: sdk.NewCoins(sdk.NewInt64Coin("ustt", 100)), Status: types.EscrowStatusPending},
			{ID: "1", Requester: testRequester, Provider: testProvider, Amount: sdk.NewCoins(sdk.NewInt64Coin("ustt", 200)), Status: types.EscrowStatusPending},
		},
	}
	err = types.ValidateGenesis(genesisWithDups)
	require.Error(t, err)
}
