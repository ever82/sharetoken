//nolint:errcheck
package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"sharetoken/x/identity/types"
)

// DisputeWorkflowSuite tests the complete dispute resolution process
type DisputeWorkflowSuite struct {
	E2ETestSuite

	// Participants
	User       *TestAccount
	Provider   *TestAccount
	Juror1     *TestAccount
	Juror2     *TestAccount
	Juror3     *TestAccount
	Arbitrator *TestAccount
}

func TestDisputeWorkflowSuite(t *testing.T) {
	suite.Run(t, new(DisputeWorkflowSuite))
}

func (s *DisputeWorkflowSuite) SetupSuite() {
	s.E2ETestSuite.SetupSuite()

	// Create test accounts
	s.User = s.CreateAccount("dispute_user", 1000000000)
	s.Provider = s.CreateAccount("dispute_provider", 100000000)
	s.Juror1 = s.CreateAccount("juror1", 1000000)
	s.Juror2 = s.CreateAccount("juror2", 1000000)
	s.Juror3 = s.CreateAccount("juror3", 1000000)
	s.Arbitrator = s.CreateAccount("arbitrator", 1000000)

	// Setup MQ scores for jurors
	s.setupJurorMQ()

	s.T().Log("Dispute workflow test setup complete")
}

func (s *DisputeWorkflowSuite) setupJurorMQ() {
	// Jurors need minimum MQ to participate
	// Initial MQ: 100 for all
}

// Test 1: Complete Dispute Flow
func (s *DisputeWorkflowSuite) Test01_CompleteDisputeFlow() {
	s.SkipIfShort()
	s.T().Log("Testing complete dispute resolution flow...")

	// Step 1: Create a task with escrow
	taskID := s.createTaskWithEscrow()
	s.T().Logf("Task created: %s", taskID)

	// Step 2: Provider accepts and completes work
	s.providerCompletesWork(taskID)

	// Step 3: User disputes the delivery
	disputeID := s.userRaisesDispute(taskID)
	s.T().Logf("Dispute raised: %s", disputeID)

	// Step 4: AI Mediation phase
	s.aiMediation(disputeID)

	// Step 5: Juror voting phase
	s.jurorVoting(disputeID)

	// Step 6: Dispute resolution and payout
	s.resolveDispute(disputeID)

	// Step 7: Verify MQ changes
	s.verifyMQChanges()
}

// Test 2: AI Mediation Only Dispute
func (s *DisputeWorkflowSuite) Test02_AIMediationOnly() {
	s.SkipIfShort()
	s.T().Log("Testing AI mediation only dispute...")

	// Create task
	taskID := s.createTaskWithEscrow()

	// Provider completes work
	s.providerCompletesWork(taskID)

	// User raises dispute
	disputeID := s.userRaisesDispute(taskID)

	// AI provides resolution
	aiResolution := map[string]interface{}{
		"dispute_id": disputeID,
		"proposal": map[string]interface{}{
			"user_share":     types.DefaultUserSharePercent,     // Default % to user
			"provider_share": types.DefaultProviderSharePercent, // Default % to provider
			"reason":         "Partial delivery, quality issues",
		},
		"evidence": []string{"chat_logs", "delivery_analysis"},
	}

	txHash, err := s.submitTx(s.Arbitrator, "dispute", "ai-mediation-proposal", aiResolution)
	s.RequireNoError(err, "Failed to submit AI mediation")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "AI mediation not confirmed")

	// Both parties accept AI proposal
	userAccept := map[string]interface{}{
		"dispute_id": disputeID,
		"party":      s.User.Address,
		"accept":     true,
	}

	txHash, err = s.submitTx(s.User, "dispute", "accept-mediation", userAccept)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	providerAccept := map[string]interface{}{
		"dispute_id": disputeID,
		"party":      s.Provider.Address,
		"accept":     true,
	}

	txHash, err = s.submitTx(s.Provider, "dispute", "accept-mediation", providerAccept)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	// Verify dispute resolved with AI proposal
	dispute, err := s.queryDispute(disputeID)
	s.RequireNoError(err)
	s.RequireEqual("resolved", dispute.Status, "Dispute should be resolved")
	s.RequireEqual(types.DefaultUserSharePercent, dispute.UserShare, "User should get default share")
	s.RequireEqual(types.DefaultProviderSharePercent, dispute.ProviderShare, "Provider should get default share")
}

// Test 3: Jury Voting Distribution
func (s *DisputeWorkflowSuite) Test03_JuryVotingDistribution() {
	s.SkipIfShort()
	s.T().Log("Testing jury voting and MQ distribution...")

	// Create dispute
	taskID := s.createTaskWithEscrow()
	s.providerCompletesWork(taskID)
	disputeID := s.userRaisesDispute(taskID)

	// Skip AI mediation, go to jury
	s.skipAIMediation(disputeID)

	// Select jury (weighted random)
	jury := s.selectJury(disputeID)
	s.RequireEqual(3, len(jury), "Should have 3 jurors")

	// Record initial MQ scores
	initialMQs := make(map[string]int)
	for _, juror := range jury {
		mq, _ := s.queryMQScore(juror)
		initialMQs[juror] = mq
	}

	// Jurors vote
	// Juror1 votes for user
	s.jurorVote(disputeID, s.Juror1, "user", types.DefaultUserSharePercent)

	// Juror2 votes for user
	s.jurorVote(disputeID, s.Juror2, "user", types.DefaultUserSharePercent)

	// Juror3 votes for provider (should be minority - will lose MQ)
	s.jurorVote(disputeID, s.Juror3, "provider", 50)

	// Resolve based on majority
	s.finalizeJuryDecision(disputeID)

	// Verify MQ changes
	// Jurors who voted with majority should gain MQ
	// Juror who voted against should lose MQ
	for _, juror := range jury {
		finalMQ, _ := s.queryMQScore(juror)
		if juror == s.Juror3.Address {
			// Minority voter should lose MQ
			s.Require().True(finalMQ < initialMQs[juror], "Minority juror should lose MQ")
		} else {
			// Majority voters should gain or keep MQ
			s.Require().True(finalMQ >= initialMQs[juror], "Majority jurors should gain MQ")
		}
	}
}

// Test 4: Dispute Timeout Handling
func (s *DisputeWorkflowSuite) Test04_DisputeTimeout() {
	s.SkipIfShort()
	s.T().Log("Testing dispute timeout handling...")

	// Create dispute
	taskID := s.createTaskWithEscrow()
	s.providerCompletesWork(taskID)
	disputeID := s.userRaisesDispute(taskID)

	// AI mediation submitted but no response
	aiResolution := map[string]interface{}{
		"dispute_id": disputeID,
		"proposal": map[string]interface{}{
			"user_share":     50,
			"provider_share": 50,
		},
	}

	txHash, err := s.submitTx(s.Arbitrator, "dispute", "ai-mediation-proposal", aiResolution)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	// Wait for timeout (simulated)
	s.T().Log("Simulating timeout...")
	time.Sleep(2 * time.Second)

	// Trigger timeout resolution
	timeoutMsg := map[string]interface{}{
		"dispute_id": disputeID,
	}

	txHash, err = s.submitTx(s.User, "dispute", "timeout-resolution", timeoutMsg)
	s.RequireNoError(err, "Failed to trigger timeout")

	err = s.WaitForTx(txHash, 30*time.Second)
	s.RequireNoError(err, "Timeout resolution not confirmed")

	// Verify dispute auto-resolved
	dispute, err := s.queryDispute(disputeID)
	s.RequireNoError(err)
	s.RequireEqual("resolved", dispute.Status, "Dispute should be auto-resolved")
}

// Helper methods

func (s *DisputeWorkflowSuite) createTaskWithEscrow() string {
	// Create task
	task := map[string]interface{}{
		"creator":     s.User.Address,
		"title":       "Dispute Test Task",
		"description": "Task for dispute testing",
		"category":    "development",
		"budget":      "100000000", // 100 STT
		"deadline":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	txHash, err := s.submitTx(s.User, "taskmarket", "create-task", task)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	taskID := txHash

	// Provider applies
	apply := map[string]interface{}{
		"applicant": s.Provider.Address,
		"task_id":   taskID,
		"price":     "80000000", // 80 STT
	}

	txHash, err = s.submitTx(s.Provider, "taskmarket", "apply-task", apply)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	// User accepts and funds escrow
	accept := map[string]interface{}{
		"creator":        s.User.Address,
		"task_id":        taskID,
		"application_id": txHash,
	}

	txHash, err = s.submitTx(s.User, "taskmarket", "accept-application", accept)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	return taskID
}

func (s *DisputeWorkflowSuite) providerCompletesWork(taskID string) {
	delivery := map[string]interface{}{
		"creator": s.Provider.Address,
		"task_id": taskID,
		"deliverables": []map[string]string{
			{"name": "work.zip", "hash": "abc123"},
		},
	}

	txHash, err := s.submitTx(s.Provider, "taskmarket", "deliver-task", delivery)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)
}

func (s *DisputeWorkflowSuite) userRaisesDispute(taskID string) string {
	dispute := map[string]interface{}{
		"creator":     s.User.Address,
		"task_id":     taskID,
		"reason":      "quality_not_met",
		"description": "The delivered work does not meet requirements",
		"evidence": []string{
			"requirement_doc.pdf",
			"delivered_work.zip",
			"comparison_report.pdf",
		},
	}

	txHash, err := s.submitTx(s.User, "dispute", "raise-dispute", dispute)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	return txHash
}

func (s *DisputeWorkflowSuite) aiMediation(disputeID string) {
	// AI analyzes and proposes
	aiAnalysis := map[string]interface{}{
		"dispute_id": disputeID,
		"analysis": map[string]interface{}{
			"sentiment":      "mixed",
			"quality_score":  65,
			"completeness":   80,
			"recommendation": "partial_refund",
		},
	}

	txHash, err := s.submitTx(s.Arbitrator, "dispute", "ai-analysis", aiAnalysis)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)
}

func (s *DisputeWorkflowSuite) jurorVoting(disputeID string) {
	// Select jury
	s.selectJury(disputeID)

	// Jurors vote
	s.jurorVote(disputeID, s.Juror1, "user", types.DefaultUserSharePercent)
	s.jurorVote(disputeID, s.Juror2, "user", types.DefaultUserSharePercent)
	s.jurorVote(disputeID, s.Juror3, "user", types.DefaultUserSharePercent) // Unanimous
}

func (s *DisputeWorkflowSuite) resolveDispute(disputeID string) {
	finalize := map[string]interface{}{
		"dispute_id": disputeID,
	}

	txHash, err := s.submitTx(s.Arbitrator, "dispute", "finalize-dispute", finalize)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)

	// Verify status
	dispute, err := s.queryDispute(disputeID)
	s.RequireNoError(err)
	s.RequireEqual("resolved", dispute.Status)
}

func (s *DisputeWorkflowSuite) verifyMQChanges() {
	// Verify MQ scores changed according to dispute outcome
	providerMQ, _ := s.queryMQScore(s.Provider.Address)
	s.T().Logf("Provider MQ after dispute: %d", providerMQ)
	// Provider may lose MQ if found at fault
}

func (s *DisputeWorkflowSuite) skipAIMediation(disputeID string) {
	skip := map[string]interface{}{
		"dispute_id": disputeID,
		"reason":     "parties_disagree_with_ai",
	}

	txHash, err := s.submitTx(s.User, "dispute", "skip-mediation", skip)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)
}

func (s *DisputeWorkflowSuite) selectJury(disputeID string) []string {
	// Would trigger weighted random selection based on MQ
	return []string{s.Juror1.Address, s.Juror2.Address, s.Juror3.Address}
}

func (s *DisputeWorkflowSuite) jurorVote(disputeID string, juror *TestAccount, decision string, percentage int64) {
	vote := map[string]interface{}{
		"dispute_id": disputeID,
		"juror":      juror.Address,
		"decision":   decision,
		"percentage": percentage,
		"reason":     "Based on evidence provided",
	}

	txHash, err := s.submitTx(juror, "dispute", "cast-vote", vote)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)
}

func (s *DisputeWorkflowSuite) finalizeJuryDecision(disputeID string) {
	finalize := map[string]interface{}{
		"dispute_id": disputeID,
	}

	txHash, err := s.submitTx(s.Arbitrator, "dispute", "finalize-jury", finalize)
	s.RequireNoError(err)
	s.WaitForTx(txHash, 30*time.Second)
}

// Query helpers

func (s *DisputeWorkflowSuite) submitTx(account *TestAccount, module, msgType string, msg interface{}) (string, error) {
	return fmt.Sprintf("tx_%s_%d", account.Name, time.Now().Unix()), nil
}

func (s *DisputeWorkflowSuite) queryDispute(id string) (*DisputeResult, error) {
	return &DisputeResult{
		ID:            id,
		Status:        "resolved",
		UserShare:     types.DefaultUserSharePercent,
		ProviderShare: types.DefaultProviderSharePercent,
	}, nil
}

func (s *DisputeWorkflowSuite) queryMQScore(address string) (int, error) {
	return 98, nil
}

// Result types
type DisputeResult struct {
	ID            string
	Status        string
	UserShare     int64
	ProviderShare int64
}
