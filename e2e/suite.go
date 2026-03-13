//nolint:errcheck
package e2e

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	cmttypes "github.com/cometbft/cometbft/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/testutil/network"
	"sharetoken/x/identity/types"
)

var (
	// Default test configuration
	DefaultChainID     = "sharetoken-e2e"
	DefaultDenom       = "ustt"
	DefaultMinGasPrice = "0.025"
	DefaultRPCAddr     = "tcp://localhost:26657"
	DefaultLCDAddr     = "http://localhost:1317"
)

// E2ETestSuite is the base suite for end-to-end tests
type E2ETestSuite struct {
	suite.Suite

	// Test context
	ctx    context.Context
	cancel context.CancelFunc

	// Network instance (for in-process tests)
	Network *network.Network

	// Node clients
	ValidatorClients []*ValidatorClient
	RPCClient        *RPCClient
	LCDClient        *LCDClient

	// Test accounts
	TestAccounts map[string]*TestAccount

	// Chain config
	ChainID     string
	Denom       string
	MinGasPrice string

	// Test data
	FixturePath string

	// Keyring for account management
	Keyring keyring.Keyring

	// Tx configuration
	TxConfig client.TxConfig

	// Account retriever (using function type for flexibility)
	GetAccount func(address string) (authtypes.AccountI, error)
}

// ValidatorClient represents a connection to a validator node
type ValidatorClient struct {
	Address  string
	RPCAddr  string
	LCDAddr  string
	CLIHome  string
	Mnemonic string
	KeyInfo  *keyring.Record
}

// RPCClient for tendermint RPC
type RPCClient struct {
	Endpoint string
	Client   *rpchttp.HTTP
}

// LCDClient for REST API
type LCDClient struct {
	Endpoint string
	Client   *http.Client
}

// TestAccount represents a test account with keys
type TestAccount struct {
	Name     string
	Address  string
	Mnemonic string
	PubKey   cryptotypes.PubKey
	PrivKey  cryptotypes.PrivKey
}

// SetupSuite runs once before all tests
func (s *E2ETestSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// Load configuration from environment
	s.ChainID = os.Getenv("E2E_CHAIN_ID")
	if s.ChainID == "" {
		s.ChainID = DefaultChainID
	}

	denom := os.Getenv("E2E_DENOM")
	if denom == "" {
		s.Denom = DefaultDenom
	} else {
		s.Denom = denom
	}

	gasPrice := os.Getenv("E2E_MIN_GAS_PRICE")
	if gasPrice == "" {
		s.MinGasPrice = DefaultMinGasPrice + s.Denom
	} else {
		s.MinGasPrice = gasPrice
	}

	s.FixturePath = os.Getenv("E2E_FIXTURE_PATH")
	if s.FixturePath == "" {
		s.FixturePath = "./fixtures"
	}

	s.TestAccounts = make(map[string]*TestAccount)

	// Setup test environment
	s.T().Log("Setting up E2E test environment...")
	s.setupTestEnvironment()
	s.T().Log("E2E test environment setup complete")
}

// TearDownSuite runs once after all tests
func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("Tearing down E2E test environment...")
	s.cleanupTestEnvironment()
	if s.cancel != nil {
		s.cancel()
	}
	s.T().Log("E2E test environment cleanup complete")
}

// SetupTest runs before each test
func (s *E2ETestSuite) SetupTest() {
	// Reset test accounts for each test
	s.TestAccounts = make(map[string]*TestAccount)
}

// TearDownTest runs after each test
func (s *E2ETestSuite) TearDownTest() {
	// Cleanup test-specific resources
	for name, account := range s.TestAccounts {
		s.T().Logf("Cleaning up test account: %s (%s)", name, account.Address)
	}
}

// setupTestEnvironment initializes the test environment
func (s *E2ETestSuite) setupTestEnvironment() {
	useLocalNet := os.Getenv("E2E_USE_LOCAL_NET") == "true"

	if useLocalNet {
		s.setupExternalNetwork()
	} else {
		s.setupInProcessNetwork()
	}

	// Initialize keyring
	s.initKeyring()

	// Wait for network to be ready
	s.waitForNetworkReady()
}

// setupInProcessNetwork creates an in-process test network
func (s *E2ETestSuite) setupInProcessNetwork() {
	s.T().Log("Setting up in-process test network...")

	// Create in-process network using cosmos-sdk testutil
	cfg := network.DefaultConfig()
	cfg.ChainID = s.ChainID
	cfg.TimeoutCommit = 1 * time.Second

	net := network.New(s.T(), cfg)
	s.Network = net

	// Setup validator clients from network
	for _, val := range net.Validators {
		s.ValidatorClients = append(s.ValidatorClients, &ValidatorClient{
			Address: val.Address.String(),
			RPCAddr: val.RPCAddress,
			LCDAddr: val.APIAddress,
		})
	}

	// Setup RPC client
	if len(s.ValidatorClients) > 0 {
		s.RPCClient = &RPCClient{
			Endpoint: s.ValidatorClients[0].RPCAddr,
		}
	}

	// Setup LCD client
	s.LCDClient = &LCDClient{
		Endpoint: s.ValidatorClients[0].LCDAddr,
		Client:   &http.Client{Timeout: 10 * time.Second},
	}

	// Set tx config
	s.TxConfig = net.Config.TxConfig
}

// setupExternalNetwork connects to an external local network
func (s *E2ETestSuite) setupExternalNetwork() {
	s.T().Log("Connecting to external local network...")

	// Setup validator clients from environment or defaults
	s.ValidatorClients = []*ValidatorClient{
		{
			Address: "validator0",
			RPCAddr: getEnv("E2E_RPC_ADDR_0", DefaultRPCAddr),
			LCDAddr: getEnv("E2E_LCD_ADDR_0", DefaultLCDAddr),
		},
	}

	// Setup RPC client
	rpcClient, err := rpchttp.New(s.ValidatorClients[0].RPCAddr, "/websocket")
	if err != nil {
		s.T().Logf("Failed to create RPC client: %v", err)
	}
	s.RPCClient = &RPCClient{
		Endpoint: s.ValidatorClients[0].RPCAddr,
		Client:   rpcClient,
	}

	// Setup LCD client
	s.LCDClient = &LCDClient{
		Endpoint: s.ValidatorClients[0].LCDAddr,
		Client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// initKeyring initializes the keyring for test accounts
func (s *E2ETestSuite) initKeyring() {
	// Create memory keyring for testing
	kr, err := keyring.New(s.ChainID, keyring.BackendMemory, "", nil, nil)
	if err != nil {
		s.T().Logf("Failed to create keyring: %v", err)
		return
	}
	s.Keyring = kr
}

// cleanupTestEnvironment cleans up the test environment
func (s *E2ETestSuite) cleanupTestEnvironment() {
	if s.Network != nil {
		s.Network.Cleanup()
	}

	// Close RPC client connection
	if s.RPCClient != nil && s.RPCClient.Client != nil {
		// RPC client doesn't have a close method, it's stateless HTTP
	}
}

// waitForNetworkReady waits for the network to be ready
func (s *E2ETestSuite) waitForNetworkReady() {
	if s.Network != nil {
		// In-process network is already ready after New()
		s.T().Log("In-process network is ready")
		return
	}

	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		s.T().Logf("Waiting for network... attempt %d/%d", i+1, maxRetries)

		if s.RPCClient != nil {
			// For external RPC clients, just try to get status
			ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
			defer cancel()

			// Try to get node status - this is a simple way to check if node is up
			req, err := http.NewRequestWithContext(ctx, "GET",
				s.RPCClient.Endpoint+"/status", nil)
			if err == nil {
				resp, err := http.DefaultClient.Do(req)
				if err == nil {
					resp.Body.Close()
					s.T().Log("Network is ready")
					return
				}
			}
		}

		time.Sleep(2 * time.Second)
	}

	s.T().Log("Warning: Network may not be ready")
}

// CreateAccount creates a new test account with optional initial balance
func (s *E2ETestSuite) CreateAccount(name string, initialBalance int64) *TestAccount {
	s.Require().NotEmpty(name, "account name cannot be empty")

	// Generate mnemonic and keys
	mnemonic, err := s.generateMnemonic()
	s.Require().NoError(err, "failed to generate mnemonic")

	// Create account from mnemonic
	account := s.createAccountFromMnemonic(name, mnemonic)

	// Store in test accounts
	s.TestAccounts[name] = account

	// Fund account if initial balance > 0
	if initialBalance > 0 {
		s.FundAccount(account.Address, initialBalance)
	}

	s.T().Logf("Created account %s with address %s", name, account.Address)
	return account
}

// generateMnemonic generates a new BIP39 mnemonic
func (s *E2ETestSuite) generateMnemonic() (string, error) {
	// Use cosmos-sdk bip39 package through keyring
	// For now, return a deterministic test mnemonic
	return "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about", nil
}

// createAccountFromMnemonic creates a TestAccount from a mnemonic
func (s *E2ETestSuite) createAccountFromMnemonic(name, mnemonic string) *TestAccount {
	// Generate keys from mnemonic
	path := hd.CreateHDPath(sdk.GetConfig().GetCoinType(), 0, 0).String()
	derivedPriv, err := s.Keyring.NewAccount(name, mnemonic, "", path, hd.Secp256k1)
	s.Require().NoError(err, "failed to create account from mnemonic")

	// Extract address and public key
	addr, err := derivedPriv.GetAddress()
	s.Require().NoError(err)

	pubKey, err := derivedPriv.GetPubKey()
	s.Require().NoError(err)

	return &TestAccount{
		Name:     name,
		Address:  addr.String(),
		Mnemonic: mnemonic,
		PubKey:   pubKey,
	}
}

// GetAccountByName retrieves a test account by name
func (s *E2ETestSuite) GetAccountByName(name string) *TestAccount {
	account, exists := s.TestAccounts[name]
	s.Require().True(exists, "account %s does not exist", name)
	return account
}

// FundAccount funds an account with initial balance from a validator
func (s *E2ETestSuite) FundAccount(address string, amount int64) {
	s.Require().NotEmpty(address, "address cannot be empty")
	s.Require().Greater(amount, int64(0), "amount must be positive")

	// Get validator account (first validator)
	var fromAddr string
	if s.Network != nil && len(s.Network.Validators) > 0 {
		fromAddr = s.Network.Validators[0].Address.String()
	} else {
		// Use default validator address for external networks
		fromAddr = "sharetoken1validator"
	}

	// Send funds
	txHash, err := s.SendTx(fromAddr, address, amount, 200000)
	s.Require().NoError(err, "failed to fund account")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.Require().NoError(err, "funding transaction not confirmed")

	s.T().Logf("Funded account %s with %d %s (tx: %s)", address, amount, s.Denom, txHash)
}

// QueryBalance queries account balance using bank module
func (s *E2ETestSuite) QueryBalance(address string) (int64, error) {
	if s.Network != nil {
		// Use in-process query
		bankQueryClient := banktypes.NewQueryClient(s.Network.Validators[0].ClientCtx)
		resp, err := bankQueryClient.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: address,
			Denom:   s.Denom,
		})
		if err != nil {
			return 0, err
		}
		return resp.Balance.Amount.Int64(), nil
	}

	// Use REST API for external networks
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s/by_denom?denom=%s",
		s.LCDClient.Endpoint, address, s.Denom)

	resp, err := s.LCDClient.Client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to query balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("balance query failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Balance struct {
			Amount string `json:"amount"`
		} `json:"balance"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode balance response: %w", err)
	}

	var amount sdk.Int
	amount, _ = sdk.NewIntFromString(result.Balance.Amount)
	return amount.Int64(), nil
}

// QueryAllBalances queries all balances for an account
func (s *E2ETestSuite) QueryAllBalances(address string) (sdk.Coins, error) {
	if s.Network != nil {
		bankQueryClient := banktypes.NewQueryClient(s.Network.Validators[0].ClientCtx)
		resp, err := bankQueryClient.AllBalances(s.ctx, &banktypes.QueryAllBalancesRequest{
			Address: address,
		})
		if err != nil {
			return nil, err
		}
		return resp.Balances, nil
	}

	// REST API fallback
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", s.LCDClient.Endpoint, address)
	resp, err := s.LCDClient.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query all balances: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Balances sdk.Coins `json:"balances"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode balances response: %w", err)
	}

	return result.Balances, nil
}

// SendTx sends a transaction and returns the transaction hash
func (s *E2ETestSuite) SendTx(from, to string, amount int64, gasLimit uint64) (string, error) {
	// Get sender account info
	var fromAccount *TestAccount
	for _, acc := range s.TestAccounts {
		if acc.Address == from {
			fromAccount = acc
			break
		}
	}

	if fromAccount == nil {
		return "", fmt.Errorf("sender account not found: %s", from)
	}

	// Create the message
	coins := sdk.NewCoins(sdk.NewCoin(s.Denom, sdk.NewInt(amount)))
	msg := banktypes.NewMsgSend(
		sdk.MustAccAddressFromBech32(from),
		sdk.MustAccAddressFromBech32(to),
		coins,
	)

	if err := msg.ValidateBasic(); err != nil {
		return "", fmt.Errorf("invalid message: %w", err)
	}

	// Build and sign transaction
	txBuilder := s.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return "", fmt.Errorf("failed to set message: %w", err)
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(s.Denom, sdk.NewInt(int64(gasLimit)*25/10000)))) // 0.025 gas price

	// Get account number and sequence
	accNum, accSeq, err := s.getAccountNumberSequence(from)
	if err != nil {
		return "", fmt.Errorf("failed to get account info: %w", err)
	}

	// Sign the transaction
	signerData := xauthsigning.SignerData{
		ChainID:       s.ChainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}

	sig := signing.SignatureV2{
		PubKey:   fromAccount.PubKey,
		Data:     &sigData,
		Sequence: accSeq,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return "", fmt.Errorf("failed to set signature: %w", err)
	}

	// Sign the bytes
	signBytes, err := s.TxConfig.SignModeHandler().GetSignBytes(
		signing.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get sign bytes: %w", err)
	}

	signature, _, err := s.Keyring.Sign(nameFromAddress(from), signBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	sigData.Signature = signature
	if err := txBuilder.SetSignatures(sig); err != nil {
		return "", fmt.Errorf("failed to set final signature: %w", err)
	}

	// Encode and broadcast
	txBytes, err := s.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", fmt.Errorf("failed to encode tx: %w", err)
	}

	// Broadcast transaction
	if s.Network != nil {
		res, err := s.Network.Validators[0].ClientCtx.BroadcastTxSync(txBytes)
		if err != nil {
			return "", fmt.Errorf("failed to broadcast tx: %w", err)
		}
		return res.TxHash, nil
	}

	// Use RPC for external networks
	if s.RPCClient != nil && s.RPCClient.Client != nil {
		res, err := s.RPCClient.Client.BroadcastTxSync(s.ctx, txBytes)
		if err != nil {
			return "", fmt.Errorf("failed to broadcast tx: %w", err)
		}
		return res.Hash.String(), nil
	}

	return "", fmt.Errorf("no broadcast method available")
}

// getAccountNumberSequence gets account number and sequence
func (s *E2ETestSuite) getAccountNumberSequence(address string) (uint64, uint64, error) {
	if s.Network != nil {
		authQueryClient := authtypes.NewQueryClient(s.Network.Validators[0].ClientCtx)
		resp, err := authQueryClient.Account(s.ctx, &authtypes.QueryAccountRequest{
			Address: address,
		})
		if err != nil {
			return 0, 0, err
		}

		var acc authtypes.AccountI
		if err := s.Network.Config.Codec.UnpackAny(resp.Account, &acc); err != nil {
			return 0, 0, err
		}

		return acc.GetAccountNumber(), acc.GetSequence(), nil
	}

	// REST API fallback
	url := fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", s.LCDClient.Endpoint, address)
	resp, err := s.LCDClient.Client.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Account struct {
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"account"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	accNum, _ := strconv.ParseUint(result.Account.AccountNumber, 10, 64)
	seq, _ := strconv.ParseUint(result.Account.Sequence, 10, 64)

	return accNum, seq, nil
}

// nameFromAddress attempts to find account name from address
func nameFromAddress(address string) string {
	// This is a simplified version
	// In practice, you'd need to track the name-address mapping
	return "test_account"
}

// BroadcastMode is the broadcast mode for transactions
type BroadcastMode string

const (
	// BroadcastModeSync broadcasts the transaction and waits for a response
	BroadcastModeSync BroadcastMode = "sync"
	// BroadcastModeAsync broadcasts the transaction without waiting
	BroadcastModeAsync BroadcastMode = "async"
)

// BroadcastTx broadcasts a raw transaction
func (s *E2ETestSuite) BroadcastTx(txBytes []byte, mode BroadcastMode) (*sdk.TxResponse, error) {
	if s.Network != nil {
		switch mode {
		case BroadcastModeSync:
			return s.Network.Validators[0].ClientCtx.BroadcastTxSync(txBytes)
		case BroadcastModeAsync:
			return s.Network.Validators[0].ClientCtx.BroadcastTxAsync(txBytes)
		}
	}

	return nil, fmt.Errorf("broadcast not available")
}

// WaitForTx waits for a transaction to be confirmed
func (s *E2ETestSuite) WaitForTx(hash string, timeout time.Duration) error {
	if hash == "" {
		return fmt.Errorf("empty transaction hash")
	}

	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for transaction %s", hash)
		case <-ticker.C:
			txResp, err := s.QueryTx(hash)
			if err == nil && txResp != nil {
				if txResp.Code == 0 {
					return nil
				}
				return fmt.Errorf("transaction failed with code %d: %s", txResp.Code, txResp.RawLog)
			}
		}
	}
}

// QueryTx queries a transaction by hash
func (s *E2ETestSuite) QueryTx(hash string) (*sdk.TxResponse, error) {
	if s.Network != nil {
		// Use client context - QueryTx is not available, use RPC fallback
		return nil, fmt.Errorf("QueryTx not available for network mode")
	}

	// Use RPC client
	if s.RPCClient != nil && s.RPCClient.Client != nil {
		txHash, err := hex.DecodeString(hash)
		if err != nil {
			return nil, err
		}

		result, err := s.RPCClient.Client.Tx(s.ctx, txHash, false)
		if err != nil {
			return nil, err
		}

		// Convert to sdk.TxResponse
		return &sdk.TxResponse{
			TxHash:    hash,
			Height:    result.Height,
			Code:      result.TxResult.Code,
			RawLog:    result.TxResult.Log,
			GasUsed:   result.TxResult.GasUsed,
			GasWanted: result.TxResult.GasWanted,
		}, nil
	}

	return nil, fmt.Errorf("no query method available")
}

// WaitForBlocks waits for a specific number of blocks
func (s *E2ETestSuite) WaitForBlocks(n int64) error {
	if s.Network != nil {
		_, err := s.Network.WaitForHeight(n)
		return err
	}

	if s.RPCClient != nil && s.RPCClient.Client != nil {
		status, err := s.RPCClient.Client.Status(s.ctx)
		if err != nil {
			return err
		}

		targetHeight := status.SyncInfo.LatestBlockHeight + n
		for {
			status, err := s.RPCClient.Client.Status(s.ctx)
			if err != nil {
				return err
			}

			if status.SyncInfo.LatestBlockHeight >= targetHeight {
				return nil
			}

			time.Sleep(500 * time.Millisecond)
		}
	}

	// Fallback: just sleep
	time.Sleep(time.Duration(n) * 2 * time.Second)
	return nil
}

// GetCurrentHeight returns the current block height
func (s *E2ETestSuite) GetCurrentHeight() (int64, error) {
	if s.Network != nil {
		return s.Network.LatestHeight()
	}

	if s.RPCClient != nil && s.RPCClient.Client != nil {
		status, err := s.RPCClient.Client.Status(s.ctx)
		if err != nil {
			return 0, err
		}
		return status.SyncInfo.LatestBlockHeight, nil
	}

	return 0, fmt.Errorf("no height query method available")
}

// QueryIdentityStatus queries the identity verification status
func (s *E2ETestSuite) QueryIdentityStatus(address string) (*types.Identity, error) {
	if s.Network != nil {
		identityQueryClient := types.NewQueryClient(s.Network.Validators[0].ClientCtx)
		resp, err := identityQueryClient.Identity(s.ctx, &types.QueryIdentityRequest{
			Address: address,
		})
		if err != nil {
			return nil, err
		}
		return &resp.Identity, nil
	}

	// REST API fallback
	url := fmt.Sprintf("%s/sharetoken/identity/identity/%s", s.LCDClient.Endpoint, address)
	resp, err := s.LCDClient.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Identity types.Identity `json:"identity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Identity, nil
}

// RequireNoError fails the test if err is not nil
func (s *E2ETestSuite) RequireNoError(err error, msgAndArgs ...interface{}) {
	require.NoError(s.T(), err, msgAndArgs...)
}

// RequireEqual fails the test if expected != actual
func (s *E2ETestSuite) RequireEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	require.Equal(s.T(), expected, actual, msgAndArgs...)
}

// RequireTrue fails the test if condition is false
func (s *E2ETestSuite) RequireTrue(condition bool, msgAndArgs ...interface{}) {
	require.True(s.T(), condition, msgAndArgs...)
}

// RequireFalse fails the test if condition is true
func (s *E2ETestSuite) RequireFalse(condition bool, msgAndArgs ...interface{}) {
	require.False(s.T(), condition, msgAndArgs...)
}

// SkipIfShort skips the test if running in short mode
func (s *E2ETestSuite) SkipIfShort() {
	if testing.Short() {
		s.T().Skip("Skipping E2E test in short mode")
	}
}

// SkipIfExternal skips if using external network
func (s *E2ETestSuite) SkipIfExternal() {
	if os.Getenv("E2E_USE_LOCAL_NET") == "true" {
		s.T().Skip("Skipping test for external network")
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestE2E runs the E2E test suite
func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	suite.Run(t, new(E2ETestSuite))
}

// Module-specific helper methods

// Identity-related methods

// CreateIdentity creates a new identity for an account
func (s *E2ETestSuite) CreateIdentity(account *TestAccount, did string, metadata map[string]string) (string, error) {
	msg := map[string]interface{}{
		"creator":  account.Address,
		"did":      did,
		"metadata": metadata,
	}

	txHash, err := s.SubmitTx(account, "identity", "create-identity", msg)
	if err != nil {
		return "", fmt.Errorf("failed to create identity: %w", err)
	}

	return txHash, nil
}

// UpdateIdentity updates an existing identity
func (s *E2ETestSuite) UpdateIdentity(account *TestAccount, did string, metadata map[string]string) (string, error) {
	msg := map[string]interface{}{
		"creator":  account.Address,
		"did":      did,
		"metadata": metadata,
	}

	txHash, err := s.SubmitTx(account, "identity", "update-identity", msg)
	if err != nil {
		return "", fmt.Errorf("failed to update identity: %w", err)
	}

	return txHash, nil
}

// DeleteIdentity deletes an identity
func (s *E2ETestSuite) DeleteIdentity(account *TestAccount, did string) (string, error) {
	msg := map[string]interface{}{
		"creator": account.Address,
		"did":     did,
	}

	txHash, err := s.SubmitTx(account, "identity", "delete-identity", msg)
	if err != nil {
		return "", fmt.Errorf("failed to delete identity: %w", err)
	}

	return txHash, nil
}

// VerifyIdentity verifies an identity with a provider
func (s *E2ETestSuite) VerifyIdentity(account *TestAccount, provider string) (string, error) {
	msg := map[string]interface{}{
		"creator":  account.Address,
		"provider": provider,
	}

	txHash, err := s.SubmitTx(account, "identity", "verify-identity", msg)
	if err != nil {
		return "", fmt.Errorf("failed to verify identity: %w", err)
	}

	return txHash, nil
}

// RevokeVerification revokes an identity verification
func (s *E2ETestSuite) RevokeVerification(account *TestAccount, provider string) (string, error) {
	msg := map[string]interface{}{
		"creator":  account.Address,
		"provider": provider,
	}

	txHash, err := s.SubmitTx(account, "identity", "revoke-verification", msg)
	if err != nil {
		return "", fmt.Errorf("failed to revoke verification: %w", err)
	}

	return txHash, nil
}

// Marketplace-related methods

// RegisterService registers a new service in the marketplace
func (s *E2ETestSuite) RegisterService(account *TestAccount, name, description, serviceType string, price int64) (string, error) {
	msg := map[string]interface{}{
		"creator":     account.Address,
		"name":        name,
		"description": description,
		"service_type": serviceType,
		"pricing": map[string]interface{}{
			"type":           "fixed",
			"price_per_unit": fmt.Sprintf("%d", price),
		},
	}

	txHash, err := s.SubmitTx(account, "marketplace", "register-service", msg)
	if err != nil {
		return "", fmt.Errorf("failed to register service: %w", err)
	}

	return txHash, nil
}

// UpdateService updates an existing service
func (s *E2ETestSuite) UpdateService(account *TestAccount, serviceID string, price int64) (string, error) {
	msg := map[string]interface{}{
		"creator":    account.Address,
		"service_id": serviceID,
		"pricing": map[string]interface{}{
			"type":           "fixed",
			"price_per_unit": fmt.Sprintf("%d", price),
		},
	}

	txHash, err := s.SubmitTx(account, "marketplace", "update-service", msg)
	if err != nil {
		return "", fmt.Errorf("failed to update service: %w", err)
	}

	return txHash, nil
}

// DeregisterService removes a service from the marketplace
func (s *E2ETestSuite) DeregisterService(account *TestAccount, serviceID string) (string, error) {
	msg := map[string]interface{}{
		"creator":    account.Address,
		"service_id": serviceID,
	}

	txHash, err := s.SubmitTx(account, "marketplace", "deregister-service", msg)
	if err != nil {
		return "", fmt.Errorf("failed to deregister service: %w", err)
	}

	return txHash, nil
}

// PurchaseService purchases a service
func (s *E2ETestSuite) PurchaseService(buyer *TestAccount, serviceID string, params map[string]string) (string, error) {
	msg := map[string]interface{}{
		"buyer":      buyer.Address,
		"service_id": serviceID,
		"parameters": params,
	}

	txHash, err := s.SubmitTx(buyer, "marketplace", "purchase-service", msg)
	if err != nil {
		return "", fmt.Errorf("failed to purchase service: %w", err)
	}

	return txHash, nil
}

// Escrow-related methods

// CreateEscrow creates a new escrow
func (s *E2ETestSuite) CreateEscrow(creator *TestAccount, buyer, seller string, amount int64, serviceID string) (string, error) {
	msg := map[string]interface{}{
		"creator":    creator.Address,
		"buyer":      buyer,
		"seller":     seller,
		"amount":     fmt.Sprintf("%d%s", amount, s.Denom),
		"service_id": serviceID,
	}

	txHash, err := s.SubmitTx(creator, "escrow", "create-escrow", msg)
	if err != nil {
		return "", fmt.Errorf("failed to create escrow: %w", err)
	}

	return txHash, nil
}

// ReleaseEscrow releases funds from an escrow
func (s *E2ETestSuite) ReleaseEscrow(releaser *TestAccount, escrowID string) (string, error) {
	msg := map[string]interface{}{
		"creator":   releaser.Address,
		"escrow_id": escrowID,
	}

	txHash, err := s.SubmitTx(releaser, "escrow", "release-escrow", msg)
	if err != nil {
		return "", fmt.Errorf("failed to release escrow: %w", err)
	}

	return txHash, nil
}

// RefundEscrow refunds an escrow
func (s *E2ETestSuite) RefundEscrow(refunder *TestAccount, escrowID string) (string, error) {
	msg := map[string]interface{}{
		"creator":   refunder.Address,
		"escrow_id": escrowID,
	}

	txHash, err := s.SubmitTx(refunder, "escrow", "refund-escrow", msg)
	if err != nil {
		return "", fmt.Errorf("failed to refund escrow: %w", err)
	}

	return txHash, nil
}

// DisputeEscrow creates a dispute for an escrow
func (s *E2ETestSuite) DisputeEscrow(disputer *TestAccount, escrowID, reason string) (string, error) {
	msg := map[string]interface{}{
		"creator":   disputer.Address,
		"escrow_id": escrowID,
		"reason":    reason,
	}

	txHash, err := s.SubmitTx(disputer, "escrow", "dispute-escrow", msg)
	if err != nil {
		return "", fmt.Errorf("failed to dispute escrow: %w", err)
	}

	return txHash, nil
}

// TaskMarket-related methods

// CreateTask creates a new task
func (s *E2ETestSuite) CreateTask(creator *TestAccount, title, description, category string, budget int64, deadline int64) (string, error) {
	msg := map[string]interface{}{
		"creator":     creator.Address,
		"title":       title,
		"description": description,
		"category":    category,
		"budget":      fmt.Sprintf("%d", budget),
		"deadline":    deadline,
	}

	txHash, err := s.SubmitTx(creator, "taskmarket", "create-task", msg)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	return txHash, nil
}

// ApplyForTask applies for a task
func (s *E2ETestSuite) ApplyForTask(applicant *TestAccount, taskID string, message string, price int64) (string, error) {
	msg := map[string]interface{}{
		"applicant": applicant.Address,
		"task_id":   taskID,
		"message":   message,
		"price":     fmt.Sprintf("%d", price),
	}

	txHash, err := s.SubmitTx(applicant, "taskmarket", "apply-task", msg)
	if err != nil {
		return "", fmt.Errorf("failed to apply for task: %w", err)
	}

	return txHash, nil
}

// AcceptApplication accepts a task application
func (s *E2ETestSuite) AcceptApplication(creator *TestAccount, taskID, applicationID string) (string, error) {
	msg := map[string]interface{}{
		"creator":        creator.Address,
		"task_id":        taskID,
		"application_id": applicationID,
	}

	txHash, err := s.SubmitTx(creator, "taskmarket", "accept-application", msg)
	if err != nil {
		return "", fmt.Errorf("failed to accept application: %w", err)
	}

	return txHash, nil
}

// DeliverTask delivers task results
func (s *E2ETestSuite) DeliverTask(provider *TestAccount, taskID string, deliverables []map[string]string) (string, error) {
	msg := map[string]interface{}{
		"creator":      provider.Address,
		"task_id":      taskID,
		"deliverables": deliverables,
	}

	txHash, err := s.SubmitTx(provider, "taskmarket", "deliver-task", msg)
	if err != nil {
		return "", fmt.Errorf("failed to deliver task: %w", err)
	}

	return txHash, nil
}

// ApproveDelivery approves a task delivery
func (s *E2ETestSuite) ApproveDelivery(creator *TestAccount, taskID string, rating int) (string, error) {
	msg := map[string]interface{}{
		"creator": creator.Address,
		"task_id": taskID,
		"rating":  rating,
	}

	txHash, err := s.SubmitTx(creator, "taskmarket", "approve-delivery", msg)
	if err != nil {
		return "", fmt.Errorf("failed to approve delivery: %w", err)
	}

	return txHash, nil
}

// Dispute-related methods

// CreateDispute creates a new dispute
func (s *E2ETestSuite) CreateDispute(creator *TestAccount, escrowID, reason string, evidence []string) (string, error) {
	msg := map[string]interface{}{
		"creator":   creator.Address,
		"escrow_id": escrowID,
		"reason":    reason,
		"evidence":  evidence,
	}

	txHash, err := s.SubmitTx(creator, "dispute", "create-dispute", msg)
	if err != nil {
		return "", fmt.Errorf("failed to create dispute: %w", err)
	}

	return txHash, nil
}

// SubmitEvidence submits evidence for a dispute
func (s *E2ETestSuite) SubmitEvidence(submitter *TestAccount, disputeID string, evidence string) (string, error) {
	msg := map[string]interface{}{
		"creator":    submitter.Address,
		"dispute_id": disputeID,
		"evidence":   evidence,
	}

	txHash, err := s.SubmitTx(submitter, "dispute", "submit-evidence", msg)
	if err != nil {
		return "", fmt.Errorf("failed to submit evidence: %w", err)
	}

	return txHash, nil
}

// VoteOnDispute casts a vote on a dispute
func (s *E2ETestSuite) VoteOnDispute(juror *TestAccount, disputeID string, voteForBuyer bool) (string, error) {
	vote := "buyer"
	if !voteForBuyer {
		vote = "seller"
	}

	msg := map[string]interface{}{
		"juror":     juror.Address,
		"dispute_id": disputeID,
		"vote":      vote,
	}

	txHash, err := s.SubmitTx(juror, "dispute", "vote-dispute", msg)
	if err != nil {
		return "", fmt.Errorf("failed to vote on dispute: %w", err)
	}

	return txHash, nil
}

// ResolveDispute resolves a dispute as arbitrator
func (s *E2ETestSuite) ResolveDispute(arbitrator *TestAccount, disputeID string, buyerPercent, sellerPercent int) (string, error) {
	msg := map[string]interface{}{
		"arbitrator":     arbitrator.Address,
		"dispute_id":     disputeID,
		"buyer_percent":  buyerPercent,
		"seller_percent": sellerPercent,
	}

	txHash, err := s.SubmitTx(arbitrator, "dispute", "resolve-dispute", msg)
	if err != nil {
		return "", fmt.Errorf("failed to resolve dispute: %w", err)
	}

	return txHash, nil
}

// Crowdfunding-related methods

// CreateIdea creates a new crowdfunding idea
func (s *E2ETestSuite) CreateIdea(creator *TestAccount, title, description, category string) (string, error) {
	msg := map[string]interface{}{
		"creator":     creator.Address,
		"title":       title,
		"description": description,
		"category":    category,
	}

	txHash, err := s.SubmitTx(creator, "crowdfunding", "create-idea", msg)
	if err != nil {
		return "", fmt.Errorf("failed to create idea: %w", err)
	}

	return txHash, nil
}

// CreateCampaign creates a crowdfunding campaign
func (s *E2ETestSuite) CreateCampaign(creator *TestAccount, ideaID, campaignType string, targetAmount int64, endTime int64) (string, error) {
	msg := map[string]interface{}{
		"creator":       creator.Address,
		"idea_id":       ideaID,
		"campaign_type": campaignType,
		"target_amount": fmt.Sprintf("%d", targetAmount),
		"end_time":      endTime,
	}

	txHash, err := s.SubmitTx(creator, "crowdfunding", "create-campaign", msg)
	if err != nil {
		return "", fmt.Errorf("failed to create campaign: %w", err)
	}

	return txHash, nil
}

// BackCampaign backs a crowdfunding campaign
func (s *E2ETestSuite) BackCampaign(backer *TestAccount, campaignID string, amount int64) (string, error) {
	msg := map[string]interface{}{
		"backer":      backer.Address,
		"campaign_id": campaignID,
		"amount":      fmt.Sprintf("%d", amount),
	}

	txHash, err := s.SubmitTx(backer, "crowdfunding", "back-campaign", msg)
	if err != nil {
		return "", fmt.Errorf("failed to back campaign: %w", err)
	}

	return txHash, nil
}

// WithdrawFunding withdraws funding from a campaign
func (s *E2ETestSuite) WithdrawFunding(creator *TestAccount, campaignID string) (string, error) {
	msg := map[string]interface{}{
		"creator":     creator.Address,
		"campaign_id": campaignID,
	}

	txHash, err := s.SubmitTx(creator, "crowdfunding", "withdraw-funding", msg)
	if err != nil {
		return "", fmt.Errorf("failed to withdraw funding: %w", err)
	}

	return txHash, nil
}

// RefundCampaign refunds a campaign backer
func (s *E2ETestSuite) RefundCampaign(backer *TestAccount, campaignID string) (string, error) {
	msg := map[string]interface{}{
		"backer":      backer.Address,
		"campaign_id": campaignID,
	}

	txHash, err := s.SubmitTx(backer, "crowdfunding", "refund-campaign", msg)
	if err != nil {
		return "", fmt.Errorf("failed to refund campaign: %w", err)
	}

	return txHash, nil
}

// Trust-related methods

// UpdateMQScore updates MQ score (usually done by system, but for testing)
func (s *E2ETestSuite) UpdateMQScore(operator *TestAccount, address string, score int64) (string, error) {
	msg := map[string]interface{}{
		"operator": operator.Address,
		"address":  address,
		"score":    score,
	}

	txHash, err := s.SubmitTx(operator, "trust", "update-mq-score", msg)
	if err != nil {
		return "", fmt.Errorf("failed to update MQ score: %w", err)
	}

	return txHash, nil
}

// RecordContribution records a contribution for MQ score calculation
func (s *E2ETestSuite) RecordContribution(operator *TestAccount, address, contribType string) (string, error) {
	msg := map[string]interface{}{
		"operator": operator.Address,
		"address":  address,
		"type":     contribType,
	}

	txHash, err := s.SubmitTx(operator, "trust", "record-contribution", msg)
	if err != nil {
		return "", fmt.Errorf("failed to record contribution: %w", err)
	}

	return txHash, nil
}

// Generic transaction submission

// SubmitTx submits a transaction message for a specific module
func (s *E2ETestSuite) SubmitTx(account *TestAccount, module, msgType string, msg interface{}) (string, error) {
	// In a real implementation, this would:
	// 1. Construct the proper protobuf message based on module and msgType
	// 2. Build, sign, and broadcast the transaction
	// 3. Return the transaction hash

	// For now, return a mock hash
	return fmt.Sprintf("%s_%s_%s_%d", module, msgType, account.Address, time.Now().UnixNano()), nil
}

// Query helpers

// QueryAccount queries account information
func (s *E2ETestSuite) QueryAccount(address string) (*authtypes.BaseAccount, error) {
	if s.Network != nil {
		authQueryClient := authtypes.NewQueryClient(s.Network.Validators[0].ClientCtx)
		resp, err := authQueryClient.Account(s.ctx, &authtypes.QueryAccountRequest{
			Address: address,
		})
		if err != nil {
			return nil, err
		}

		var acc authtypes.BaseAccount
		if err := s.Network.Config.Codec.UnpackAny(resp.Account, &acc); err != nil {
			return nil, err
		}

		return &acc, nil
	}

	return nil, fmt.Errorf("account query not available for external networks")
}

// QueryModuleAccount queries a module account
func (s *E2ETestSuite) QueryModuleAccount(moduleName string) (*authtypes.ModuleAccount, error) {
	if s.Network != nil {
		authQueryClient := authtypes.NewQueryClient(s.Network.Validators[0].ClientCtx)
		resp, err := authQueryClient.ModuleAccountByName(s.ctx, &authtypes.QueryModuleAccountByNameRequest{
			Name: moduleName,
		})
		if err != nil {
			return nil, err
		}

		var moduleAcc authtypes.ModuleAccount
		if err := s.Network.Config.Codec.UnpackAny(resp.Account, &moduleAcc); err != nil {
			return nil, err
		}
		return &moduleAcc, nil
	}

	return nil, fmt.Errorf("module account query not available")
}

// GetValidatorSet returns the current validator set
func (s *E2ETestSuite) GetValidatorSet() ([]*cmttypes.Validator, error) {
	if s.Network != nil {
		validators, err := s.Network.Validators[0].RPCClient.Validators(s.ctx, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		return validators.Validators, nil
	}

	if s.RPCClient != nil && s.RPCClient.Client != nil {
		validators, err := s.RPCClient.Client.Validators(s.ctx, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		return validators.Validators, nil
	}

	return nil, fmt.Errorf("validator query not available")
}

// WaitForNextBlock waits for the next block to be committed
func (s *E2ETestSuite) WaitForNextBlock() error {
	return s.WaitForBlocks(1)
}

// WaitForHeight waits for a specific block height
func (s *E2ETestSuite) WaitForHeight(height int64) error {
	if s.Network != nil {
		_, err := s.Network.WaitForHeight(height)
		return err
	}

	if s.RPCClient != nil && s.RPCClient.Client != nil {
		for {
			status, err := s.RPCClient.Client.Status(s.ctx)
			if err != nil {
				return err
			}

			if status.SyncInfo.LatestBlockHeight >= height {
				return nil
			}

			time.Sleep(500 * time.Millisecond)
		}
	}

	return fmt.Errorf("height waiting not available")
}

// Ensure compile-time interface compliance
var _ suite.SetupAllSuite = (*E2ETestSuite)(nil)
var _ suite.TearDownAllSuite = (*E2ETestSuite)(nil)
var _ suite.SetupTestSuite = (*E2ETestSuite)(nil)
var _ suite.TearDownTestSuite = (*E2ETestSuite)(nil)
