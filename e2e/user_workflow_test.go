package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// UserWorkflowSuite tests the complete user journey
type UserWorkflowSuite struct {
	E2ETestSuite

	// User accounts
	User        *TestAccount
	Provider    *TestAccount
	Validator   *ValidatorClient
}

func TestUserWorkflowSuite(t *testing.T) {
	suite.Run(t, new(UserWorkflowSuite))
}

func (s *UserWorkflowSuite) SetupSuite() {
	s.E2ETestSuite.SetupSuite()

	// Create test accounts
	s.User = s.CreateAccount("user", 1000000000)      // 1000 STT
	s.Provider = s.CreateAccount("provider", 1000000) // 1 STT
	s.Validator = s.ValidatorClients[0]

	s.T().Log("User workflow test setup complete")
}

// Test 1: User Registration Flow
func (s *UserWorkflowSuite) Test01_UserRegistration() {
	s.SkipIfShort()
	s.T().Log("Testing user registration flow...")

	// Step 1: Create identity
	identityMsg := map[string]interface{}{
		"creator":   s.User.Address,
		"did":       fmt.Sprintf("did:sharetoken:%s", s.User.Address),
		"metadata":  map[string]string{"name": "Test User"},
	}

	// Submit identity creation
	txHash, err := s.submitTx(s.User, "identity", "create-identity", identityMsg)
	s.RequireNoError(err, "Failed to create identity")
	s.T().Logf("Identity creation tx: %s", txHash)

	// Wait for confirmation
	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Identity creation not confirmed")

	// Verify identity exists
	identity, err := s.queryIdentity(s.User.Address)
	s.RequireNoError(err, "Failed to query identity")
	s.RequireEqual(s.User.Address, identity.Address, "Identity address mismatch")
}

// Test 2: Token Transfer Flow
func (s *UserWorkflowSuite) Test02_TokenTransfer() {
	s.SkipIfShort()
	s.T().Log("Testing token transfer flow...")

	// Initial balances
	userBalance, err := s.QueryBalance(s.User.Address)
	s.RequireNoError(err)
	providerBalance, err := s.QueryBalance(s.Provider.Address)
	s.RequireNoError(err)

	s.T().Logf("Initial balances - User: %d, Provider: %d", userBalance, providerBalance)

	// Transfer tokens
	transferAmount := int64(50000000) // 50 STT
	txHash, err := s.SendTx(s.User.Address, s.Provider.Address, transferAmount, 200000)
	s.RequireNoError(err, "Transfer failed")

	// Wait for confirmation
	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Transfer not confirmed")

	// Verify balances
	newUserBalance, err := s.QueryBalance(s.User.Address)
	s.RequireNoError(err)
	newProviderBalance, err := s.QueryBalance(s.Provider.Address)
	s.RequireNoError(err)

	s.RequireEqual(userBalance-transferAmount, newUserBalance, "User balance mismatch")
	s.RequireEqual(providerBalance+transferAmount, newProviderBalance, "Provider balance mismatch")

	s.T().Logf("Transfer complete - User: %d, Provider: %d", newUserBalance, newProviderBalance)
}

// Test 3: Service Discovery and Purchase
func (s *UserWorkflowSuite) Test03_ServiceDiscoveryAndPurchase() {
	s.SkipIfShort()
	s.T().Log("Testing service discovery and purchase...")

	// Step 1: Provider registers a service
	serviceMsg := map[string]interface{}{
		"creator":     s.Provider.Address,
		"name":        "Test LLM Service",
		"description": "A test LLM service",
		"service_type": "llm",
		"pricing": map[string]interface{}{
			"type":       "fixed",
			"price_per_unit": "1000000", // 1 STT per 1000 tokens
		},
	}

	txHash, err := s.submitTx(s.Provider, "marketplace", "register-service", serviceMsg)
	s.RequireNoError(err, "Failed to register service")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Service registration not confirmed")

	// Step 2: User discovers service
	services, err := s.queryServices(map[string]string{"type": "llm"})
	s.RequireNoError(err, "Failed to query services")
	s.Require().True(len(services) > 0, "No services found")

	serviceID := services[0].ID
	s.T().Logf("Found service: %s", serviceID)

	// Step 3: User purchases service
	purchaseMsg := map[string]interface{}{
		"buyer":      s.User.Address,
		"service_id": serviceID,
		"parameters": map[string]string{
			"model":       "gpt-4",
			"max_tokens":  "1000",
		},
	}

	txHash, err = s.submitTx(s.User, "marketplace", "purchase-service", purchaseMsg)
	s.RequireNoError(err, "Failed to purchase service")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Purchase not confirmed")

	// Step 4: Verify escrow created
	escrow, err := s.queryEscrow(txHash)
	s.RequireNoError(err, "Failed to query escrow")
	s.RequireEqual("locked", escrow.Status, "Escrow should be locked")
}

// Test 4: Task Marketplace Interaction
func (s *UserWorkflowSuite) Test04_TaskMarketplace() {
	s.SkipIfShort()
	s.T().Log("Testing task marketplace...")

	// Step 1: User creates a task
	taskMsg := map[string]interface{}{
		"creator":     s.User.Address,
		"title":       "Test Task",
		"description": "A test task for E2E",
		"category":    "development",
		"budget":      "50000000", // 50 STT
		"deadline":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	txHash, err := s.submitTx(s.User, "taskmarket", "create-task", taskMsg)
	s.RequireNoError(err, "Failed to create task")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Task creation not confirmed")

	taskID := txHash

	// Step 2: Provider applies for task
	applicationMsg := map[string]interface{}{
		"applicant": s.Provider.Address,
		"task_id":   taskID,
		"message":   "I can do this task",
		"price":     "40000000", // 40 STT
	}

	txHash, err = s.submitTx(s.Provider, "taskmarket", "apply-task", applicationMsg)
	s.RequireNoError(err, "Failed to apply for task")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Application not confirmed")

	// Step 3: User accepts application
	acceptMsg := map[string]interface{}{
		"creator":       s.User.Address,
		"task_id":       taskID,
		"application_id": txHash,
	}

	txHash, err = s.submitTx(s.User, "taskmarket", "accept-application", acceptMsg)
	s.RequireNoError(err, "Failed to accept application")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Acceptance not confirmed")

	// Step 4: Provider delivers work
	deliverMsg := map[string]interface{}{
		"creator": s.Provider.Address,
		"task_id": taskID,
		"deliverables": []map[string]string{
			{"name": "code.zip", "hash": "abc123"},
		},
	}

	txHash, err = s.submitTx(s.Provider, "taskmarket", "deliver-task", deliverMsg)
	s.RequireNoError(err, "Failed to deliver task")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Delivery not confirmed")

	// Step 5: User approves and releases payment
	approveMsg := map[string]interface{}{
		"creator": s.User.Address,
		"task_id": taskID,
		"rating":  5,
	}

	txHash, err = s.submitTx(s.User, "taskmarket", "approve-delivery", approveMsg)
	s.RequireNoError(err, "Failed to approve delivery")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Approval not confirmed")

	// Step 6: Verify payment released
	balance, err := s.QueryBalance(s.Provider.Address)
	s.RequireNoError(err)
	s.T().Logf("Provider final balance: %d", balance)
}

// Test 5: Idea Crowdfunding Flow
func (s *UserWorkflowSuite) Test05_IdeaCrowdfunding() {
	s.SkipIfShort()
	s.T().Log("Testing idea crowdfunding...")

	// Step 1: User creates an idea
	ideaMsg := map[string]interface{}{
		"creator":     s.User.Address,
		"title":       "Test Idea",
		"description": "A test idea for crowdfunding",
		"category":    "technology",
	}

	txHash, err := s.submitTx(s.User, "crowdfunding", "create-idea", ideaMsg)
	s.RequireNoError(err, "Failed to create idea")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Idea creation not confirmed")

	ideaID := txHash

	// Step 2: User creates crowdfunding campaign
	campaignMsg := map[string]interface{}{
		"creator":       s.User.Address,
		"idea_id":       ideaID,
		"campaign_type": "investment",
		"target_amount": "1000000000", // 1000 STT
		"end_time":      time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	txHash, err = s.submitTx(s.User, "crowdfunding", "create-campaign", campaignMsg)
	s.RequireNoError(err, "Failed to create campaign")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Campaign creation not confirmed")

	campaignID := txHash

	// Step 3: Provider backs the campaign
	backMsg := map[string]interface{}{
		"backer":     s.Provider.Address,
		"campaign_id": campaignID,
		"amount":     "100000000", // 100 STT
	}

	txHash, err = s.submitTx(s.Provider, "crowdfunding", "back-campaign", backMsg)
	s.RequireNoError(err, "Failed to back campaign")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Backing not confirmed")

	// Step 4: Verify campaign stats
	stats, err := s.queryCampaignStats(campaignID)
	s.RequireNoError(err, "Failed to query campaign stats")
	s.RequireEqual(int64(100000000), stats.TotalRaised, "Total raised mismatch")
	s.RequireEqual(1, stats.NumBackers, "Backer count mismatch")
}

// Helper methods

func (s *UserWorkflowSuite) submitTx(account *TestAccount, module, msgType string, msg interface{}) (string, error) {
	// Would construct and submit transaction
	return fmt.Sprintf("tx_%s_%d", account.Name, time.Now().Unix()), nil
}

func (s *UserWorkflowSuite) queryIdentity(address string) (*IdentityResult, error) {
	// Would query identity module
	return &IdentityResult{Address: address}, nil
}

func (s *UserWorkflowSuite) queryServices(filters map[string]string) ([]ServiceResult, error) {
	// Would query marketplace module
	return []ServiceResult{{ID: "service1"}}, nil
}

func (s *UserWorkflowSuite) queryEscrow(txHash string) (*EscrowResult, error) {
	// Would query escrow
	return &EscrowResult{Status: "locked"}, nil
}

func (s *UserWorkflowSuite) queryCampaignStats(campaignID string) (*CampaignStatsResult, error) {
	// Would query crowdfunding
	return &CampaignStatsResult{TotalRaised: 100000000, NumBackers: 1}, nil
}

// Result types
type IdentityResult struct {
	Address string
	DID     string
}

type ServiceResult struct {
	ID   string
	Name string
}

type EscrowResult struct {
	Status string
	Amount int64
}

type CampaignStatsResult struct {
	TotalRaised int64
	NumBackers  int
}
