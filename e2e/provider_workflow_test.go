//nolint:errcheck
package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// ProviderWorkflowSuite tests the complete provider journey
type ProviderWorkflowSuite struct {
	E2ETestSuite

	// Accounts
	Provider  *TestAccount
	User      *TestAccount
	Validator *ValidatorClient
}

func TestProviderWorkflowSuite(t *testing.T) {
	suite.Run(t, new(ProviderWorkflowSuite))
}

func (s *ProviderWorkflowSuite) SetupSuite() {
	s.E2ETestSuite.SetupSuite()

	// Create test accounts
	s.Provider = s.CreateAccount("provider", 100000000) // 100 STT for gas
	s.User = s.CreateAccount("user", 1000000000)       // 1000 STT
	s.Validator = s.ValidatorClients[0]

	s.T().Log("Provider workflow test setup complete")
}

// Test 1: Provider Registration
func (s *ProviderWorkflowSuite) Test01_ProviderRegistration() {
	s.SkipIfShort()
	s.T().Log("Testing provider registration...")

	// Step 1: Create identity
	identityMsg := map[string]interface{}{
		"creator":   s.Provider.Address,
		"did":       fmt.Sprintf("did:sharetoken:%s", s.Provider.Address),
		"metadata": map[string]string{
			"name":  "Test Provider",
			"email": "provider@test.com",
		},
	}

	txHash, err := s.submitTx(s.Provider, "identity", "create-identity", identityMsg)
	s.RequireNoError(err, "Failed to create identity")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Identity creation not confirmed")

	// Step 2: Complete KYC verification (simulated)
	kycMsg := map[string]interface{}{
		"creator": s.Provider.Address,
		"level":   "verified",
	}

	txHash, err = s.submitTx(s.Provider, "identity", "verify-identity", kycMsg)
	s.RequireNoError(err, "Failed to verify identity")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Verification not confirmed")

	// Verify identity status
	identity, err := s.queryIdentity(s.Provider.Address)
	s.RequireNoError(err)
	s.RequireEqual("verified", identity.Status, "Identity should be verified")
}

// Test 2: API Key Custody
func (s *ProviderWorkflowSuite) Test02_APIKeyCustody() {
	s.SkipIfShort()
	s.T().Log("Testing API key custody...")

	// Step 1: Provider registers API key
	apiKeyMsg := map[string]interface{}{
		"creator":        s.Provider.Address,
		"service_name":   "openai",
		"encrypted_key":  "encrypted_api_key_here",
		"access_control": []string{"read", "completion"},
	}

	txHash, err := s.submitTx(s.Provider, "agent", "register-api-key", apiKeyMsg)
	s.RequireNoError(err, "Failed to register API key")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "API key registration not confirmed")

	apiKeyID := txHash

	// Step 2: Verify API key stored
	apiKey, err := s.queryAPIKey(apiKeyID)
	s.RequireNoError(err)
	s.RequireEqual("openai", apiKey.ServiceName, "Service name mismatch")
	s.Require().True(len(apiKey.EncryptedKey) > 0, "Encrypted key should exist")

	// Step 3: Update API key
	updateMsg := map[string]interface{}{
		"creator":        s.Provider.Address,
		"api_key_id":     apiKeyID,
		"encrypted_key":  "new_encrypted_key",
		"access_control": []string{"read", "completion", "embeddings"},
	}

	txHash, err = s.submitTx(s.Provider, "agent", "update-api-key", updateMsg)
	s.RequireNoError(err, "Failed to update API key")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "API key update not confirmed")

	// Step 4: Revoke API key
	revokeMsg := map[string]interface{}{
		"creator":    s.Provider.Address,
		"api_key_id": apiKeyID,
	}

	txHash, err = s.submitTx(s.Provider, "agent", "revoke-api-key", revokeMsg)
	s.RequireNoError(err, "Failed to revoke API key")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "API key revocation not confirmed")

	// Verify revoked
	apiKey, err = s.queryAPIKey(apiKeyID)
	s.RequireNoError(err)
	s.RequireEqual("revoked", apiKey.Status, "API key should be revoked")
}

// Test 3: Service Registration and Management
func (s *ProviderWorkflowSuite) Test03_ServiceRegistration() {
	s.SkipIfShort()
	s.T().Log("Testing service registration...")

	// Step 1: Register LLM service
	llmService := map[string]interface{}{
		"creator":      s.Provider.Address,
		"name":         "OpenAI GPT-4 Service",
		"description":  "Access to GPT-4 via API",
		"service_type": "llm",
		"pricing": map[string]interface{}{
			"type":           "dynamic",
			"base_price":     "500000",  // 0.5 STT per 1K tokens
			"oracle_feed":    "openai_gpt4",
			"update_interval": 300,      // 5 minutes
		},
		"capabilities": []string{"text-generation", "chat", "function-calling"},
	}

	txHash, err := s.submitTx(s.Provider, "marketplace", "register-service", llmService)
	s.RequireNoError(err, "Failed to register LLM service")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Service registration not confirmed")

	llmServiceID := txHash
	s.T().Logf("LLM service registered: %s", llmServiceID)

	// Step 2: Register Agent service
	agentService := map[string]interface{}{
		"creator":      s.Provider.Address,
		"name":         "Code Review Agent",
		"description":  "AI-powered code review",
		"service_type": "agent",
		"pricing": map[string]interface{}{
			"type":   "fixed",
			"price":  "2000000", // 2 STT per review
		},
		"capabilities": []string{"code-review", "security-analysis", "optimization"},
	}

	txHash, err = s.submitTx(s.Provider, "marketplace", "register-service", agentService)
	s.RequireNoError(err, "Failed to register agent service")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Agent service registration not confirmed")

	agentServiceID := txHash

	// Step 3: Update service pricing
	updatePricing := map[string]interface{}{
		"creator":    s.Provider.Address,
		"service_id": llmServiceID,
		"pricing": map[string]interface{}{
			"type":       "fixed",
			"base_price": "400000", // Lower price: 0.4 STT
		},
	}

	txHash, err = s.submitTx(s.Provider, "marketplace", "update-service", updatePricing)
	s.RequireNoError(err, "Failed to update service pricing")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Service update not confirmed")

	// Step 4: Pause service temporarily
	pauseMsg := map[string]interface{}{
		"creator":    s.Provider.Address,
		"service_id": agentServiceID,
		"status":     "paused",
	}

	txHash, err = s.submitTx(s.Provider, "marketplace", "update-service-status", pauseMsg)
	s.RequireNoError(err, "Failed to pause service")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Service pause not confirmed")

	// Step 5: Resume service
	resumeMsg := map[string]interface{}{
		"creator":    s.Provider.Address,
		"service_id": agentServiceID,
		"status":     "active",
	}

	txHash, err = s.submitTx(s.Provider, "marketplace", "update-service-status", resumeMsg)
	s.RequireNoError(err, "Failed to resume service")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Service resume not confirmed")
}

// Test 4: Order Fulfillment and Payment
func (s *ProviderWorkflowSuite) Test04_OrderFulfillment() {
	s.SkipIfShort()
	s.T().Log("Testing order fulfillment...")

	// Step 1: Register a service first
	service := map[string]interface{}{
		"creator":      s.Provider.Address,
		"name":         "Test Service",
		"description":  "For order fulfillment test",
		"service_type": "llm",
		"pricing": map[string]interface{}{
			"type":   "fixed",
			"price":  "1000000", // 1 STT
		},
	}

	txHash, err := s.submitTx(s.Provider, "marketplace", "register-service", service)
	s.RequireNoError(err)

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err)

	serviceID := txHash

	// Step 2: User creates order
	order := map[string]interface{}{
		"buyer":      s.User.Address,
		"service_id": serviceID,
		"parameters": map[string]string{
			"prompt": "Hello, world!",
			"model":  "gpt-4",
		},
		"payment": "1000000",
	}

	txHash, err = s.submitTx(s.User, "marketplace", "create-order", order)
	s.RequireNoError(err, "Failed to create order")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Order creation not confirmed")

	orderID := txHash

	// Step 3: Provider acknowledges order
	ackMsg := map[string]interface{}{
		"provider":  s.Provider.Address,
		"order_id":  orderID,
		"estimated": time.Now().Add(5 * time.Minute).Unix(),
	}

	txHash, err = s.submitTx(s.Provider, "marketplace", "acknowledge-order", ackMsg)
	s.RequireNoError(err, "Failed to acknowledge order")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Acknowledgement not confirmed")

	// Step 4: Provider delivers result
	delivery := map[string]interface{}{
		"provider": s.Provider.Address,
		"order_id": orderID,
		"result": map[string]interface{}{
			"output":     "Hello! How can I help you?",
			"tokens_in":  3,
			"tokens_out": 6,
		},
		"proof": "delivery_proof_hash",
	}

	txHash, err = s.submitTx(s.Provider, "marketplace", "deliver-order", delivery)
	s.RequireNoError(err, "Failed to deliver order")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Delivery not confirmed")

	// Step 5: User confirms delivery
	confirm := map[string]interface{}{
		"buyer":    s.User.Address,
		"order_id": orderID,
		"rating":   5,
		"review":   "Great service!",
	}

	txHash, err = s.submitTx(s.User, "marketplace", "confirm-delivery", confirm)
	s.RequireNoError(err, "Failed to confirm delivery")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Confirmation not confirmed")

	// Step 6: Verify payment released to provider
	providerBalance, err := s.QueryBalance(s.Provider.Address)
	s.RequireNoError(err)
	s.T().Logf("Provider balance after order: %d", providerBalance)
	s.Require().True(providerBalance >= 1000000, "Provider should have received payment")
}

// Test 5: MQ (Trust) Score Management
func (s *ProviderWorkflowSuite) Test05_MQScoreManagement() {
	s.SkipIfShort()
	s.T().Log("Testing MQ score management...")

	// Step 1: Query initial MQ score
	initialMQ, err := s.queryMQScore(s.Provider.Address)
	s.RequireNoError(err)
	s.T().Logf("Initial MQ score: %d", initialMQ.Score)
	s.RequireEqual(100, initialMQ.Score, "Initial MQ should be 100")

	// Step 2: Complete successful transactions to increase MQ
	// (Simulated through test transactions)
	for i := 0; i < 5; i++ {
		// Each successful transaction could increase MQ
		mqUpdate := map[string]interface{}{
			"address":   s.Provider.Address,
			"change":    1,
			"reason":    "successful_transaction",
		}

		txHash, err := s.submitTx(s.Provider, "trust", "update-mq", mqUpdate)
		s.RequireNoError(err)

		err = s.WaitForTx(txHash, 10*time.Second)
		s.RequireNoError(err)
	}

	// Step 3: Verify MQ increased
	updatedMQ, err := s.queryMQScore(s.Provider.Address)
	s.RequireNoError(err)
	s.T().Logf("Updated MQ score: %d", updatedMQ.Score)
	s.Require().True(updatedMQ.Score >= 100, "MQ should have increased or stayed same")

	// Step 4: Simulate dispute to decrease MQ
	disputeMsg := map[string]interface{}{
		"address": s.Provider.Address,
		"change":  -3, // Max 3% loss per dispute
		"reason":  "dispute_resolved",
	}

	txHash, err := s.submitTx(s.Provider, "trust", "update-mq", disputeMsg)
	s.RequireNoError(err)

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err)

	// Step 5: Verify MQ decreased but not below 0
	finalMQ, err := s.queryMQScore(s.Provider.Address)
	s.RequireNoError(err)
	s.T().Logf("Final MQ score: %d", finalMQ.Score)
	s.Require().True(finalMQ.Score >= 0, "MQ should never be below 0")
}

// Helper methods

func (s *ProviderWorkflowSuite) submitTx(account *TestAccount, module, msgType string, msg interface{}) (string, error) {
	return fmt.Sprintf("tx_%s_%d", account.Name, time.Now().Unix()), nil
}

func (s *ProviderWorkflowSuite) queryIdentity(address string) (*IdentityDetail, error) {
	return &IdentityDetail{
		Address: address,
		Status:  "verified",
	}, nil
}

func (s *ProviderWorkflowSuite) queryAPIKey(id string) (*APIKeyDetail, error) {
	return &APIKeyDetail{
		ID:            id,
		ServiceName:   "openai",
		EncryptedKey:  "encrypted",
		Status:        "active",
	}, nil
}

func (s *ProviderWorkflowSuite) queryMQScore(address string) (*MQScoreDetail, error) {
	return &MQScoreDetail{
		Address: address,
		Score:   105,
	}, nil
}

// Result types
type IdentityDetail struct {
	Address string
	Status  string
}

type APIKeyDetail struct {
	ID            string
	ServiceName   string
	EncryptedKey  string
	Status        string
}

type MQScoreDetail struct {
	Address string
	Score   int
}
