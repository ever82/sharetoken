# x/dispute - Trust System (Reputation & Dispute Arbitration)

> **Module Type:** Core Module
> **Tech Stack:** Go (Cosmos SDK)
> **Location:** `src/chain/x/dispute`
> **Dependencies:** Base Types (01-base), Identity (10-identity)

---

## Overview

x/dispute is a core module of ShareTokens, containing two core functions:
1. **Moral Quotient (MQ) Scoring** - Zero-sum reputation scoring system
2. **Dispute Arbitration** - Decentralized dispute resolution based on MQ weights

Implemented as a Cosmos SDK module, it is a core component of the platform's trust foundation.

---

## Architecture Position

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ShareTokens Core Modules                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐                               │
│  │    P2P   │   │ Identity │   │  Wallet  │                               │
│  └────┬─────┘   └────┬─────┘   └────┬─────┘                               │
│       │              │              │                                       │
│       └──────────────┼──────────────┘                                       │
│                      ▼                                                      │
│              ┌──────────────┐                                               │
│              │   Service    │                                               │
│              │   Marketplace│                                               │
│              └──────┬───────┘                                               │
│                     │                                                       │
│       ┌─────────────┴─────────────┐                                         │
│       ▼                           ▼                                         │
│  ┌──────────────┐         ┌──────────────────┐                             │
│  │    Escrow    │────────►│   Trust System   │ ← This Chapter              │
│  └──────────────┘         │  (09-dispute)    │                             │
│                           └──────────────────┘                             │
│                                   │                                         │
│                                   ▼                                         │
│                      Trust foundation for all transactions                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Part 1: Moral Quotient (MQ) Scoring

### Design Principles

```
1. Zero-Sum: Total MQ is constant, loser's loss = winner's gain
2. Convergence: Higher MQ = harder to increase, high MQ users give more than receive
3. Duty of Justice: Participation in arbitration is mandatory, absence results in token/MQ penalty
4. Controlled Risk: Maximum 3% MQ loss per transaction, never negative
5. Wisdom of Crowds: Random jury selection, multi-proposal scoring for consensus
```

### System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Trust System Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Chain (x/dispute)                             │   │
│  │                                                                     │   │
│  │  • Dispute state machine                                           │   │
│  │  • MediationSession event log                                      │   │
│  │  • MQ records                                                      │   │
│  │  • Jury selection & voting                                         │   │
│  │  • Zero-sum redistribution algorithm                               │   │
│  │  • Evidence hashes (not files)                                     │   │
│  │  • Debt tracking (for unpaid fines)                                │   │
│  │                                                                     │   │
│  └────────────────────────────────┬────────────────────────────────────┘   │
│                                   │                                         │
│                                   │ Transactions & Queries                  │
│                                   │                                         │
│  ┌────────────────────────────────┴────────────────────────────────────┐   │
│  │                     Off-chain Services                               │   │
│  │                                                                     │   │
│  │  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐ │   │
│  │  │   AI Mediator   │    │ Evidence Store  │    │   Oracle/Relayer │   │
│  │  │   (LLM Service) │    │   (IPFS/S3)     │    │                 │ │   │
│  │  │                 │    │                 │    │                 │ │   │
│  │  │ • Conversations │    │ • File storage │    │ • AI signature  │ │   │
│  │  │ • Evidence AI   │    │ • URL → hash   │    │   verification  │ │   │
│  │  │ • Proposals     │    │ • Content hash │    │ • Submit AI     │ │   │
│  │  │ • Decisions     │    │   on chain     │    │   verdicts      │ │   │
│  │  └────────┬────────┘    └────────┬────────┘    └────────┬────────┘ │   │
│  │           │                      │                      │          │   │
│  └───────────┼──────────────────────┼──────────────────────┼──────────┘   │
│              │                      │                      │              │
│              └──────────────────────┼──────────────────────┘              │
│                                     │                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### AI Service Architecture

```go
// AI is an off-chain service that submits signed verdicts to chain

type AIServiceConfig struct {
    // AI service identity
    ServiceAddress   sdk.AccAddress   // Authorized AI service address
    PublicKey        crypto.PubKey    // For signature verification

    // Rate limiting
    MaxResponseTime  time.Duration    // 30 seconds per message
    MediationTimeout time.Duration    // 7 days max
}

// AI verdict submission (called by oracle/relayer)
type MsgSubmitAIVerdict struct {
    DisputeId    uint64
    Verdict      VerdictContent
    Timestamp    time.Time
    AISignature  []byte            // AI service signature
    Submitter    sdk.AccAddress    // Oracle/relayer address
}

// Chain verifies AI signature before accepting verdict
func (k Keeper) VerifyAIVerdictSignature(verdict MsgSubmitAIVerdict) bool {
    // 1. Check submitter is authorized oracle
    // 2. Verify AI signature on verdict content
    // 3. Check timestamp is recent
    return true
}
```

**Why AI is an Off-chain Service:**
- LLM computation is heavy, cannot run on-chain
- Needs access to external knowledge (market prices, technical docs)
- Response time is unpredictable
- Can upgrade models without affecting on-chain logic

### Evidence Storage

```go
// Evidence stored off-chain, only hash on chain

type Evidence struct {
    Id          uint64
    DisputeId   uint64
    Submitter   sdk.AccAddress

    // On-chain (minimal)
    Type        EvidenceType     // "screenshot" | "chat_log" | "document" | "link"
    ContentHash []byte           // SHA-256 of file content
    URL         string           // IPFS/S3 URL (optional)

    // Off-chain metadata (submitted but not stored on chain)
    // - Actual file content
    // - File size, mime type
    // - AI analysis result

    Timestamp   time.Time
}

type EvidenceType string

const (
    EvidenceTypeScreenshot EvidenceType = "screenshot"
    EvidenceTypeChatLog    EvidenceType = "chat_log"
    EvidenceTypeDocument   EvidenceType = "document"
    EvidenceTypeLink       EvidenceType = "link"
    EvidenceTypeVideo      EvidenceType = "video"
)
```

**Storage Strategy:**
```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│     User         │     │  Evidence Store  │     │   Chain          │
│                  │     │    (IPFS/S3)     │     │                  │
│  1. Upload file  │────►│  2. Store file   │     │                  │
│                  │     │  3. Return URL   │     │                  │
│                  │     │     + hash       │     │                  │
│                  │     │                  │     │                  │
│  4. Submit hash  │────►│                  │────►│  5. Store hash   │
│     + URL        │     │                  │     │     only         │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

### No-Stake Enforcement Mechanism

```go
// No upfront stake required. System enforces payment after ruling.

type DebtRecord struct {
    Address     sdk.AccAddress
    Amount      sdk.Coins        // Amount owed
    Reason      string           // "dispute_fine" | "jury_penalty" | "absence_penalty"
    DisputeId   uint64           // Related dispute
    CreatedAt   time.Time
}

// Restrictions while in debt
type DebtRestrictions struct {
    CanTransfer     bool  // false - cannot send tokens
    CanWithdraw     bool  // false - cannot withdraw from escrow
    CanCreateOrder  bool  // false - cannot create new orders
    CanApplyForTask bool  // false - cannot apply for tasks
    CanUseAI        bool  // false - premium AI features disabled
}

// Automatic deduction when funds available
func (k Keeper) TryCollectDebt(ctx sdk.Context, debtor sdk.AccAddress) error {
    debt := k.GetDebt(ctx, debtor)
    if debt.Amount.IsZero() {
        return nil
    }

    balance := k.bankKeeper.GetBalance(ctx, debtor)
    if balance.IsGTE(debt.Amount) {
        // Full payment
        k.bankKeeper.SendCoinsFromAccountToModule(ctx, debtor, types.ModuleName, debt.Amount)
        k.ClearDebt(ctx, debtor)
    } else {
        // Partial payment
        k.bankKeeper.SendCoinsFromAccountToModule(ctx, debtor, types.ModuleName, balance)
        k.ReduceDebt(ctx, debtor, balance)
    }
    return nil
}

// Called on every transaction from this address
func (k Keeper) CheckDebtAndRestrict(ctx sdk.Context, addr sdk.AccAddress) error {
    debt := k.GetDebt(ctx, addr)
    if debt.Amount.IsZero() {
        return nil  // No restrictions
    }

    // Apply restrictions based on debt age and amount
    restrictions := k.GetDebtRestrictions(ctx, debt)

    // Log warning for user
    ctx.EventManager().EmitEvent(sdk.NewEvent(
        "debt_warning",
        sdk.NewAttribute("address", addr.String()),
        sdk.NewAttribute("amount", debt.Amount.String()),
    ))

    return nil  // Allow transaction but restrictions apply elsewhere
}
```

**Enforcement Flow:**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      No-Stake Enforcement Flow                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Dispute Resolved                                                           │
│       │                                                                     │
│       ▼                                                                     │
│  Calculate Fine (e.g., loser pays winner 100 STT)                          │
│       │                                                                     │
│       ▼                                                                     │
│  Check Loser's Wallet                                                       │
│       │                                                                     │
│       ├─────────────────┬─────────────────┐                                │
│       │                 │                 │                                │
│   Balance >= 100    Balance < 100     Balance = 0                           │
│       │                 │                 │                                │
│       ▼                 ▼                 ▼                                │
│  Deduct 100 STT    Deduct all         Create debt                          │
│  Transfer to       Create debt        = 100 STT                            │
│  winner            = 100 - balance    Restrictions:                         │
│       │                 │             - No transfers                        │
│       │                 │             - No new orders                       │
│       │                 │             - No withdrawals                      │
│       │                 │             - Premium features blocked            │
│       │                 │                                                 │
│       ▼                 ▼                                                 │
│  Done              On next deposit:                                        │
│                    Auto-deduct until debt cleared                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### MQ Initialization

```go
// New users get initial MQ upon identity verification

type InitMQConfig struct {
    InitialMQ    uint64  // 100 - starting MQ
    RequireVerification  bool    // true - must verify identity first
}

// Called when user completes identity verification
func (k Keeper) InitializeMQ(ctx sdk.Context, address sdk.AccAddress) error {
    // Check if already has MQ
    if k.HasMQ(ctx, address) {
        return errors.New("MQ already initialized")
    }

    config := k.GetInitMQConfig(ctx)

    // Create initial MQ record
    record := MQRecord{
        Address:    address,
        MQ: config.InitialMQ,  // 100
        History: []MQHistoryEntry{{
            Change:    int64(config.InitialMQ),
            Reason:    "initial_verification",
            Timestamp: ctx.BlockTime(),
        }},
        CreatedAt: ctx.BlockTime(),
        UpdatedAt: ctx.BlockTime(),
    }

    k.SetMQRecord(ctx, record)
    return nil
}
```

**Initialization Flow:**
```
User Registration → Identity Verification → Automatically get 100 MQ
```

---

### MQ Scoring

```go
// types/mq.go

// MQ Configuration
type MQConfig struct {
    // Basic parameters
    InitialMQ uint64  // 100 - initial MQ

    // Weighting parameters
    RiskRate        sdk.Dec // 0.03 (3%) - max risk rate per transaction
    Lambda          sdk.Dec // 1.5 - baseline deviation
    MaxDeviation    sdk.Dec // 6.0 - max reasonable deviation (score range -10~10)
    ChangeFactor    sdk.Dec // 0.01 (1%) - base change factor

    // Jury configuration
    JurySize        JurySizeConfig
    AmountThresholds AmountThresholdConfig
}

type JurySizeConfig struct {
    Small   uint64  // 5 (small amount)
    Medium  uint64  // 11 (medium amount)
    Large   uint64  // 21 (large amount)
    Huge    uint64  // 31 (huge amount)
}

type AmountThresholdConfig struct {
    Small  AmountThreshold
    Medium AmountThreshold
    Large AmountThreshold
}

type AmountThreshold struct {
    Amount sdk.Int
    Symbol string
}

// MQ record
type MQRecord struct {
    Address     sdk.AccAddress
    MQ  uint64          // Current MQ score
    History     []MQHistoryEntry
    Stats       MQStats
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type MQStats struct {
    TotalDisputes    uint64
    DisputesWon      uint64
    DisputesLost     uint64
    JuryDuties       uint64
    CorrectVotes     uint64
    WrongVotes       uint64
    TotalGained      int64
    TotalLost        int64
}

type MQHistoryEntry struct {
    Change     int64
    Reason     string
    RelatedId  uint64  // Dispute ID or other reference
    Timestamp  time.Time
}
```

### MQ Permission Configuration

```go
// MQ-based permissions (configurable, no fixed levels)
type MQConfig struct {
    // Minimum MQ to create dispute
    MinMQToCreateDispute uint64  // e.g., 50

    // Minimum MQ to be eligible for jury duty
    MinMQForJury uint64  // e.g., 50

    // Jury voting weight = MQ directly (no multipliers)
    // Higher MQ = higher weight naturally
}
```

**Design Principles:**
- No preset levels, MQ is a continuous value
- Voting weight = MQ value itself (natural weighting)
- Minimum thresholds are configurable (e.g., 50 MQ to create dispute)
- Let the system evolve naturally, no artificial limits

---

## Part 2: Dispute Arbitration

### Design Philosophy

```
1. AI-First: AI mediates the entire dispute resolution process
2. Free Conversation: Natural dialogue flow, not rigid stages
3. Intelligent Decisions: AI decides when to request evidence, propose solutions, or escalate
4. Early Settlement: Any moment can lead to settlement if both parties agree
5. Last Resort: Jury voting only when AI mediation truly fails
```

### Simplified State Machine

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Dispute Resolution State Machine                        │
└─────────────────────────────────────────────────────────────────────────────┘

                            ┌──────────┐
                            │  filed   │  Dispute filed
                            └────┬─────┘
                                 │
                            AI Accepts
                                 │
                                 ▼
                     ┌────────────────────┐
                     │    mediating       │  AI Mediation (unified)
                     │                    │  - Free conversation
                     │                    │  - Evidence submission
                     │                    │  - Proposals & scoring
                     │                    │  - AI-guided negotiation
                     └─────────┬──────────┘
                               │
             ┌─────────────────┼─────────────────┐
             │                 │                 │
        Settled           CannotSettle      OnePartyForfeit
             │                 │                 │
             ▼                 ▼                 ▼
      ┌───────────┐    ┌────────────┐    ┌───────────┐
      │  settled  │    │   juried   │    │  forfeit  │
      │ (settled) │    │(jury voting)│    │(default)  │
      └─────┬─────┘    └─────┬──────┘    └─────┬─────┘
            │                │                 │
            │          Jury Scores            │
            │                │                 │
            │                ▼                 │
            │         ┌──────────┐            │
            │         │ resolved │            │
            │         └────┬─────┘            │
            │              │                  │
            └──────────────┼──────────────────┘
                           │
                    7-day Appeal Window
                           │
              ┌────────────┼────────────┐
              │            │            │
          No Appeal    AppealFiled   AppealRejected
              │            │            │
              ▼            ▼            │
        ┌──────────┐ ┌──────────┐      │
        │  final   │ │ appealed │      │
        └──────────┘ └────┬─────┘      │
                           │           │
                      Expired           │
                           │           │
                           ▼           │
                     ┌──────────┐      │
                     │  final   │◄─────┘
                     └──────────┘

State Transitions:
┌─────────────────┬───────────────────────────┬──────────────────────────────┐
│ Current State   │ Event                     │ Next State                   │
├─────────────────┼───────────────────────────┼──────────────────────────────┤
│ (initial)       │ CreateDispute             │ filed                        │
│ filed           │ AIAccepts                 │ mediating                    │
│ mediating       │ Settled                   │ settled                      │
│ mediating       │ CannotSettle              │ juried                       │
│ mediating       │ OnePartyForfeit           │ forfeit                      │
│ juried          │ AllScored                 │ resolved                     │
│ settled         │ 7-day pass                │ final                        │
│ resolved        │ 7-day pass                │ final                        │
│ resolved        │ AppealFiled               │ appealed                     │
│ appealed        │ AppealExpired             │ final                        │
│ forfeit         │ 7-day pass                │ final                        │
└─────────────────┴───────────────────────────┴──────────────────────────────┘
```

### AI Mediation Session

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AI Mediation Session (Free Conversation)                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  AI: "I'm the mediation AI. Party A, please state your claim."             │
│                                                                             │
│  A: "I paid 100 STT for code service, but the code doesn't run."           │
│                                                                             │
│  AI: "Party B, please respond. Both parties please submit evidence."       │
│                                                                             │
│  B: "The code runs fine, it's their environment issue."                     │
│  [B uploads: success screenshot, README doc]                                │
│  [A uploads: error screenshot, terminal logs]                               │
│                                                                             │
│  AI: "I see B's screenshot is Mac, A's error is Windows. A, did you        │
│       follow the README configuration?"                                     │
│                                                                             │
│  A: "Yes, still doesn't work. Also, delivery was 2 days late."             │
│                                                                             │
│  B: "Late because A changed requirements 3 times!"                          │
│                                                                             │
│  AI: "I understand both sides. Here's my proposal: refund 40 STT,          │
│       B keeps 60 STT. Please score this from +10 to -10."                  │
│                                                                             │
│  A scores: +3                                                               │
│  B scores: -2                                                               │
│                                                                             │
│  AI: "Not far apart. Let's compromise: refund 35 STT. Agree?"              │
│                                                                             │
│  A: "Agreed"                                                                │
│  B: "OK"                                                                    │
│                                                                             │
│  → Settlement reached                                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### AI Decision Logic

```go
// AI evaluates continuously and decides next action

type AIDecision struct {
    // Conversation analysis
    ConflictLevel    float64   // 0-1, how intense is the conflict
    EmotionA         float64   // A's emotional state (angry/calm)
    EmotionB         float64   // B's emotional state

    // Evidence analysis
    EvidenceA_Quality float64  // Quality of A's evidence
    EvidenceB_Quality float64  // Quality of B's evidence

    // Settlement probability
    SettlementProbability float64  // Likelihood of reaching settlement

    // Next action recommendation
    NextAction       AINextAction
}

type AINextAction string

const (
    ActionChat        AINextAction = "chat"         // Continue conversation
    ActionAskEvidence AINextAction = "ask_evidence" // Request evidence
    ActionPropose     AINextAction = "propose"      // Propose solution
    ActionScore       AINextAction = "score"        // Request scoring
    ActionSettle      AINextAction = "settle"       // Declare settlement
    ActionJury        AINextAction = "jury"         // Escalate to jury
)

// AI Decision Logic (simplified)
//
// After each message, AI updates its assessment:
//
// 1. Conversation rounds > 10 with no progress?
//    → May need to escalate to jury
//
// 2. Both parties calming down?
//    → Try proposing a solution
//
// 3. Evidence insufficient?
//    → Request more evidence
//
// 4. Just scored, divergence < 3?
//    → Declare settlement
//
// 5. Tried 2-3 proposals, divergence still > 6?
//    → Escalate to jury voting
//
// 6. Default
//    → Continue conversation
```

### Timeout Configuration

```go
// types/timeout.go

type DisputeTimeoutConfig struct {
    // Mediation timeout
    MediationTimeout   time.Duration  // 7 days - max time for AI mediation

    // Jury voting timeout
    JuryVotingTimeout  time.Duration  // 3 days - timeout for jury scoring

    // Appeal period
    AppealPeriod       time.Duration  // 7 days - window for appeals

    // Inactivity timeout
    InactivityTimeout  time.Duration  // 48 hours - if no response, escalate/forfeit
}

var DefaultDisputeTimeoutConfig = DisputeTimeoutConfig{
    MediationTimeout:   7 * 24 * time.Hour,   // 7 days
    JuryVotingTimeout:  3 * 24 * time.Hour,   // 3 days
    AppealPeriod:       7 * 24 * time.Hour,   // 7 days
    InactivityTimeout:  48 * time.Hour,       // 48 hours
}
```

---

## Core Type Definitions

### Dispute

```go
type Dispute struct {
    Id uint64

    // Associated order
    OrderId uint64

    // Parties
    Plaintiff  DisputeParty
    Defendant  DisputeParty

    // Dispute content
    Type        DisputeType
    Title       string
    Description string
    Amount      TokenAmount

    // AI Mediation session (unified)
    Mediation   MediationSession

    // Jury (if mediation fails)
    Jury        *JuryInfo

    // Resolution
    Resolution  *Resolution

    // Status
    Status      DisputeStatus
    AppealCount uint64

    // Timestamps
    CreatedAt   time.Time
    UpdatedAt   time.Time
    ResolvedAt  *time.Time
    FinalAt     *time.Time
}

type DisputeParty struct {
    Address         sdk.AccAddress
    MQStake uint64  // Staked MQ
}

type DisputeType string

const (
    DisputeTypeServiceNotDelivered DisputeType = "service_not_delivered"
    DisputeTypePaymentNotReceived  DisputeType = "payment_not_received"
    DisputeTypeWrongAmount         DisputeType = "wrong_amount"
    DisputeTypeQualityIssue        DisputeType = "quality_issue"
    DisputeTypeServiceFailure      DisputeType = "service_failure"
    DisputeTypeFraud               DisputeType = "fraud"
    DisputeTypeBreachOfContract    DisputeType = "breach_of_contract"
    DisputeTypeOther               DisputeType = "other"
)

type DisputeStatus string

const (
    DisputeStatusFiled     DisputeStatus = "filed"      // Just filed
    DisputeStatusMediating DisputeStatus = "mediating"  // AI mediation in progress
    DisputeStatusSettled   DisputeStatus = "settled"    // Settled via mediation
    DisputeStatusJuried    DisputeStatus = "juried"     // Escalated to jury
    DisputeStatusForfeit   DisputeStatus = "forfeit"    // One party forfeited
    DisputeStatusResolved  DisputeStatus = "resolved"   // Resolved (by jury)
    DisputeStatusAppealed  DisputeStatus = "appealed"   // Under appeal
    DisputeStatusFinal     DisputeStatus = "final"      // Final, no more appeals
)
```

### MediationSession - AI Mediation Session

```go
// MediationSession represents the entire AI-led mediation process
type MediationSession struct {
    DisputeId   uint64

    // Event timeline (unified)
    Events      []MediationEvent

    // Current state
    Status      MediationStatus  // "active" | "settled" | "juried" | "forfeit"

    // Outcome
    Outcome     *MediationOutcome

    // Metadata
    StartedAt   time.Time
    EndedAt     *time.Time
}

// MediationEvent - unified event type for the timeline
type MediationEvent struct {
    Id        uint64
    Timestamp time.Time
    Actor     string       // "A" | "B" | "AI"
    Type      EventType

    // Content varies by type
    Content   *MessageContent
    Evidence  *EvidenceContent
    Proposal  *ProposalContent
    Score     *ScoreContent
    Verdict   *VerdictContent
}

type EventType string

const (
    EventMessage   EventType = "message"    // Chat message
    EventEvidence  EventType = "evidence"   // Evidence submission
    EventProposal  EventType = "proposal"   // Solution proposal
    EventScore     EventType = "score"      // Scoring
    EventVerdict   EventType = "verdict"    // AI judgment
)

type MessageContent struct {
    Text      string
    Private   bool        // Private message to AI only
}

type EvidenceContent struct {
    EvidenceId   uint64
    Description  string
    Attachment   *string   // URL or hash
    Verified     bool      // AI verified
}

type ProposalContent struct {
    ProposalId  uint64
    Content     string
    TokenSplit  *TokenSplit
    Proposer    string      // "A" | "B" | "AI"
}

type ScoreContent struct {
    ProposalId  uint64
    Score       int8        // -10 to +10
}

type VerdictContent struct {
    Action      string      // "continue" | "settle" | "jury"
    Reasoning   string
    Divergence  float64     // Score divergence if applicable
}

type MediationStatus string

const (
    MediationStatusActive  MediationStatus = "active"
    MediationStatusSettled MediationStatus = "settled"
    MediationStatusJuried  MediationStatus = "juried"
    MediationStatusForfeit MediationStatus = "forfeit"
)

type MediationOutcome struct {
    Type          string          // "settlement" | "jury" | "forfeit"
    Settlement    *Settlement     // If settled
    JurySummary   *JurySummary    // If escalated to jury
}

type Settlement struct {
    FinalProposal  uint64
    TokenSplit     TokenSplit
    AgreedAt       time.Time
}

type JurySummary struct {
    // Summary generated by AI for jury review
    PartyA_Claim     string
    PartyA_Evidence  []string
    PartyB_Claim     string
    PartyB_Evidence  []string
    Proposals        []ProposalSummary
    AI_Assessment    string
}

type ProposalSummary struct {
    Id          uint64
    Content     string
    Proposer    string
    AvgScore    float64
}
```

### Proposal - Resolution Proposal

```go
// Proposal for dispute resolution
type Proposal struct {
    Id          uint64
    DisputeId   uint64
    Proposer    sdk.AccAddress  // Proposer
    Type        ProposalType    // Proposal type
    Content     string          // Proposal content
    TokenSplit  *TokenSplit     // Token allocation (optional)
    CreatedAt   time.Time
}

type ProposalType string

const (
    ProposalTypePlaintiff  ProposalType = "plaintiff"   // Plaintiff's proposal
    ProposalTypeDefendant  ProposalType = "defendant"   // Defendant's proposal
    ProposalTypeAI         ProposalType = "ai"          // AI supplement proposal
    ProposalTypeMediation  ProposalType = "mediation"   // Mediation proposal
)

type TokenSplit struct {
    PlaintiffPercent sdk.Dec  // Plaintiff's share
    DefendantPercent sdk.Dec  // Defendant's share
}
```

### ProposalScore - Proposal Scoring

```go
// ProposalScore - proposal scoring record
type ProposalScore struct {
    ProposalId uint64
    DisputeId  uint64
    Scorer     sdk.AccAddress  // Scorer
    Score      int8            // Score -10 ~ 10
    Comment    string          // Scoring reason (optional)
    CreatedAt  time.Time
}

// ScoringResult - scoring result summary
type ScoringResult struct {
    DisputeId       uint64
    ProposalResults []ProposalResult
    Deviations      []ScoringDeviation
    FinalProposal   uint64  // Final selected proposal ID
    CalculatedAt    time.Time
}

type ProposalResult struct {
    ProposalId       uint64
    WeightedAverage  sdk.Dec  // MQ weighted average score
    TotalScores      uint64
    TotalWeight      uint64
}

type ScoringDeviation struct {
    Address        sdk.AccAddress
    MQ     uint64
    Deviations     []float64  // Deviation per proposal
    StdDeviation   float64    // Combined standard deviation
    IsPenalized    bool       // d > λ
    IsRewarded     bool       // d < λ
}
```

### JuryInfo - Jury Information

```go
type JuryInfo struct {
    Members         []JuryMember
    Size            uint64
    SelectionMethod string  // "random" | "weighted_random"
    SelectedAt      *time.Time
}

type JuryMember struct {
    Address     sdk.AccAddress
    MQ  uint64  // MQ at time of selection
    Scores      map[uint64]int8  // proposalId -> score
    VotedAt     *time.Time
    Reward      int64   // Reward received (can be positive or negative)
    Absent      bool    // Whether absent
}

type Verdict string

const (
    VerdictPlaintiff Verdict = "plaintiff"
    VerdictDefendant Verdict = "defendant"
    VerdictNeutral   Verdict = "neutral"
)
```

### Resolution - Ruling Result

```go
type Resolution struct {
    DisputeId uint64

    // Ruling result
    FinalProposal   uint64         // Final selected proposal
    ProposalResults []ProposalResult

    // MQ redistribution
    MQRedistribution MQRedistribution

    // Fund handling
    FundResolution FundResolution

    // Timestamps
    ResolvedAt time.Time
    IsFinal    bool
}

type MQRedistribution struct {
    Participants []ParticipantMQChange  // All participants (parties + jurors)
    TotalPenalty int64                  // Total penalty pool
    TotalReward  int64                  // Total reward pool
}

type ParticipantMQChange struct {
    Address     sdk.AccAddress
    Role        string    // "plaintiff" | "defendant" | "juror"
    Before      uint64
    After       uint64
    Change      int64     // Positive = gain, negative = loss
    Deviation   float64   // Scoring deviation
    IsPenalized bool
    IsRewarded  bool
}

type FundResolution struct {
    TokenReleasedTo   *sdk.AccAddress
    PaymentReleasedTo *sdk.AccAddress
    PartialSplit      *PartialSplit
}

type PartialSplit struct {
    TokenPercent   sdk.Dec
    PaymentPercent sdk.Dec
}
```

---

## Voting Weight Calculation

```go
// keeper/voting.go

// Voting weight formula: Weight = MQ
//
// Simplified design: Use MQ directly as weight
// MQ weighting algorithm already achieves convergence, no additional weight calculation needed

type VotingWeight struct {
    MQ   uint64
    Weight      sdk.Dec
}

func CalculateVotingWeight(mq uint64) sdk.Dec {
    return sdk.NewDec(int64(mq))
}

// Examples:
// MQ=100: Weight = 100
// MQ=200: Weight = 200
// MQ=500: Weight = 500
```

---

## MQ Weighting Algorithm (Zero-Sum Convergence)

### Core Principles

```
1. Zero-Sum: Penalty pool = Reward pool, total MQ unchanged
2. Convergence: Rewards use logarithm (diminishing returns), penalties use linear (high risk)
3. Fairness: Same deviation → Same proportional change
4. Safety: Maximum 3% loss per transaction, MQ never negative
```

### Parameter Definitions

```go
// MQ weighting parameters
const (
    InitialMQ = 100   // Initial MQ
    RiskRate          = 0.03  // Max risk rate 3%
    Lambda            = 1.5   // Baseline deviation
    MaxDeviation      = 6.0   // Max reasonable deviation
)
```

### Unified Mediation Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Unified AI Mediation Flow                                 │
└─────────────────────────────────────────────────────────────────────────────┘

  Dispute Filed
       │
       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     AI Mediation (Unified Process)                           │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │ Free Conversation Flow (AI controls the pace):                         │ │
│  │                                                                       │ │
│  │  Party A states claim ──────────────────────────────────────────────┐ │ │
│  │                                                                      │ │
│  │  Party B responds ──────────────────────────────────────────────────┤ │ │
│  │                                                                      │ │
│  │  [Evidence submission] ◄── AI requests when needed                  │ │ │
│  │                                                                      │ │
│  │  [Dialogue continues] ◄── AI guides, asks questions                 │ │ │
│  │                                                                      │ │
│  │  [AI proposes solution] ◄── When timing is right                    │ │ │
│  │                                                                      │ │
│  │  [Scoring] ◄── Both parties score proposals                         │ │ │
│  │                                                                      │ │
│  │  ┌───────────────────┴───────────────────┐                          │ │
│  │  │                                       │                          │ │
│  │  ▼                                       ▼                          │ │
│  │ Divergence small                    Divergence large                │ │
│  │  │                                       │                          │ │
│  │  ▼                                       ▼                          │ │
│  │ Settlement                        More proposals                   │ │
│  │ (Done!)                                │                            │ │
│  │                                        ▼                            │ │
│  │                              Still can't settle?                    │ │
│  │                                        │                            │ │
│  │                                        ▼                            │ │
│  │                              Escalate to Jury                       │ │
│  │                                                                      │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  Any moment → Settlement possible → End                                     │
│  Any moment → One party absent → Forfeit                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
       │
       │ Cannot settle
       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Jury Voting Phase                                    │
│                                                                             │
│  1. AI generates dispute summary (claims + evidence + proposals)            │
│  2. Random jury selection (MQ-weighted)                            │
│  3. Each juror scores all proposals (-10 to +10)                           │
│  4. Calculate weighted consensus (MQ-weighted average)              │
│  5. Calculate deviation from consensus for each juror                       │
│  6. Redistribute MQ (zero-sum algorithm)                           │
│  7. Final verdict based on highest-scoring proposal                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
       │
       ▼
   Resolved
```

### Scoring Data Structure

```go
// Proposal scoring
type ProposalScore struct {
    ProposalId   uint64
    Scores       map[string]int8  // address -> score (-10 ~ 10)
}

// Scoring deviation calculation
type ScoringDeviation struct {
    Address        sdk.AccAddress
    MQ     uint64
    Deviations     []float64  // Deviation per proposal
    StdDeviation   float64    // Combined standard deviation
    IsPenalized    bool       // d > λ
    IsRewarded     bool       // d < λ
}
```

### Algorithm Implementation

```go
// keeper/mq_redistribution.go

// CalculateScoringDeviations calculates each person's scoring deviation
func CalculateScoringDeviations(
    proposals []ProposalScore,
    jurors map[string]uint64,  // address -> MQ
    config MQConfig,
) []ScoringDeviation {
    results := make([]ScoringDeviation, 0)

    // Calculate MQ weighted average score for each proposal (crowd consensus)
    consensusScores := make([]float64, len(proposals))
    for i, proposal := range proposals {
        var weightedSum, totalWeight float64
        for addr, score := range proposal.Scores {
            mq := float64(jurors[addr])
            weightedSum += float64(score) * mq
            totalWeight += mq
        }
        if totalWeight > 0 {
            consensusScores[i] = weightedSum / totalWeight
        }
    }

    // Calculate each person's deviation from consensus
    for addr, mq := range jurors {
        deviation := ScoringDeviation{
            Address:    sdk.AccAddress(addr),
            MQ: mq,
            Deviations: make([]float64, len(proposals)),
        }

        // Calculate deviation for each proposal
        var sumSquaredDiff float64
        for i, proposal := range proposals {
            score := float64(proposal.Scores[addr])
            diff := score - consensusScores[i]
            deviation.Deviations[i] = diff
            sumSquaredDiff += diff * diff
        }

        // Combined standard deviation
        deviation.StdDeviation = math.Sqrt(sumSquaredDiff / float64(len(proposals)))

        // Determine reward/penalty type
        if deviation.StdDeviation > config.Lambda.MustFloat64() {
            deviation.IsPenalized = true
        } else if deviation.StdDeviation < config.Lambda.MustFloat64() {
            deviation.IsRewarded = true
        }

        results = append(results, deviation)
    }

    return results
}

// CalculateMQRedistribution zero-sum MQ weighting
func CalculateMQRedistribution(
    deviations []ScoringDeviation,
    config MQConfig,
) map[string]int64 {
    changes := make(map[string]int64)
    lambda := config.Lambda.MustFloat64()
    maxD := config.MaxDeviation.MustFloat64()
    riskRate := config.RiskRate.MustFloat64()

    // ========== Step 1: Calculate penalty pool ==========
    var penaltyPool float64
    for _, d := range deviations {
        if d.IsPenalized {
            // Penalty rate = 3% × (d - λ) / (max_d - λ)
            penaltyRate := riskRate * (d.StdDeviation - lambda) / (maxD - lambda)
            penaltyRate = math.Min(penaltyRate, riskRate) // No more than 3%

            loss := float64(d.MQ) * penaltyRate
            penaltyPool += loss
            changes[d.Address.String()] = -int64(math.Round(loss))
        }
    }

    // ========== Step 2: Calculate reward distribution ==========
    type rewardCandidate struct {
        addr      string
        mq       uint64
        contrib   float64
    }

    var candidates []rewardCandidate
    var totalContrib float64

    for _, d := range deviations {
        if d.IsRewarded {
            // Contribution score = fairness × logarithmic suppression (convergence)
            // contrib = (λ - d) × log(D + 1)
            contrib := (lambda - d.StdDeviation) * math.Log(float64(d.MQ)+1)
            candidates = append(candidates, rewardCandidate{
                addr:    d.Address.String(),
                mq:     d.MQ,
                contrib: contrib,
            })
            totalContrib += contrib
        }
    }

    // ========== Step 3: Distribute reward pool ==========
    if totalContrib > 0 && penaltyPool > 0 {
        allocFactor := penaltyPool / totalContrib

        for _, c := range candidates {
            gain := c.contrib * allocFactor
            changes[c.addr] = changes[c.addr] + int64(math.Round(gain))
        }
    }

    return changes
}

// ApplyMQChanges applies MQ changes (ensures non-negative)
func ApplyMQChanges(
    ctx sdk.Context,
    k Keeper,
    changes map[string]int64,
) error {
    for addrStr, change := range changes {
        addr, _ := sdk.AccAddressFromBech32(addrStr)
        currentMQ := k.GetMQ(ctx, addr)

        newMQ := int64(currentMQ) + change
        if newMQ < 0 {
            newMQ = 0  // MQ never negative, asymptotically approaches 0
        }

        k.SetMQ(ctx, addr, uint64(newMQ))
    }
    return nil
}
```

### Algorithm Example

```
Parameters: λ=1.5, max_d=6, risk=3%

Participants:
┌──────┬─────────────┬────────┬────────┬─────────────────────────────────┐
│ User │ MQ  │ Dev d │ Type   │ Calculation                      │
├──────┼─────────────┼────────┼────────┼─────────────────────────────────┤
│ A    │ 100         │ 0      │ Reward │ contrib = 1.5 × log(101) = 6.9  │
│ B    │ 200         │ 0.5    │ Reward │ contrib = 1.0 × log(201) = 5.3  │
│ C    │ 100         │ 1.5    │ Neutral│ -                               │
│ D    │ 100         │ 3.0    │ Penalty│ rate = 3% × 1.5/4.5 = 1%        │
│      │             │        │        │ loss = 100 × 1% = 1             │
│ E    │ 200         │ 6.0    │ Penalty│ rate = 3% × 4.5/4.5 = 3%        │
│      │             │        │        │ loss = 200 × 3% = 6             │
└──────┴─────────────┴────────┴────────┴─────────────────────────────────┘

Penalty pool = 1 + 6 = 7
Total contribution = 6.9 + 5.3 = 12.2
Allocation factor = 7 / 12.2 = 0.57

Results:
┌──────┬─────────┬──────────────┬───────────┐
│ User │ Change  │ New MQ │ % Change  │
├──────┼─────────┼──────────────┼───────────┤
│ A    │ +3.9    │ 103.9        │ +3.9%     │
│ B    │ +3.0    │ 203.0        │ +1.5%     │ ← High MQ, slower % growth
│ C    │ 0       │ 100          │ 0%        │
│ D    │ -1      │ 99           │ -1%       │
│ E    │ -6      │ 194          │ -3%       │ ← High MQ, higher risk
└──────┴─────────┴──────────────┴───────────┘

Verify zero-sum: 3.9 + 3.0 = 6.9 ≈ 7 (penalty pool) ✓
Verify convergence: Low MQ A (+3.9%) > High MQ B (+1.5%) ✓
```

### Jury Duty and Absence Penalty

```go
// Jury duty configuration
type JuryDutyConfig struct {
    // Absence penalty
    AbsenceTokenPenalty       sdk.Coin  // Token deduction
    AbsenceMQPenalty  uint64    // MQ deduction

    // Multiple absences
    MaxAbsences             uint64    // Max absence count
    MultipleAbsencePenalty  uint64    // Extra penalty for multiple absences
}

// ApplyAbsencePenalty handles absence penalty
func (k Keeper) ApplyAbsencePenalty(ctx sdk.Context, juror sdk.AccAddress) error {
    config := k.GetJuryDutyConfig(ctx)

    // Deduct tokens
    if err := k.bankKeeper.SendCoinsFromAccountToModule(
        ctx, juror, types.ModuleName, sdk.NewCoins(config.AbsenceTokenPenalty),
    ); err != nil {
        return err
    }

    // Deduct MQ
    currentMQ := k.GetMQ(ctx, juror)
    newMQ := uint64(math.Max(0, float64(currentMQ)-float64(config.AbsenceMQPenalty)))
    k.SetMQ(ctx, juror, newMQ)

    // Record absence
    k.RecordAbsence(ctx, juror)

    return nil
}
```

---

## Jury Selection Algorithm

```go
// keeper/jury.go

func DetermineJurySize(amount TokenAmount, config MQConfig) uint64 {
    amt := amount.Amount

    if amt < config.AmountThresholds.Small.Amount {
        return config.JurySize.Small
    } else if amt < config.AmountThresholds.Medium.Amount {
        return config.JurySize.Medium
    } else if amt < config.AmountThresholds.Large.Amount {
        return config.JurySize.Large
    }
    return config.JurySize.Huge
}

func SelectJury(
    eligibleJurors []JurorCandidate,
    size uint64,
    excludeAddresses map[string]bool,
) []sdk.AccAddress {
    // Filter excluded addresses
    var candidates []JurorCandidate
    for _, j := range eligibleJurors {
        if !excludeAddresses[j.Address.String()] {
            candidates = append(candidates, j)
        }
    }

    // Calculate weights
    var totalMQ uint64
    for _, c := range candidates {
        totalMQ += c.MQ
    }

    selected := make([]sdk.AccAddress, 0, size)
    selectedSet := make(map[string]bool)

    for uint64(len(selected)) < size && uint64(len(selected)) < uint64(len(candidates)) {
        // Weighted random selection
        random := rand.Uint64() % totalMQ
        var cumulative uint64

        for _, candidate := range candidates {
            if selectedSet[candidate.Address.String()] {
                continue
            }

            cumulative += candidate.MQ
            if random < cumulative {
                selected = append(selected, candidate.Address)
                selectedSet[candidate.Address.String()] = true
                break
            }
        }
    }

    return selected
}
```

---

## Keeper Interface

```go
// keeper/keeper.go

type Keeper struct {
    storeKey     sdk.StoreKey
    cdc          codec.BinaryCodec
    bankKeeper   bankkeeper.Keeper
    escrowKeeper escrowkeeper.Keeper
}

// MQ
func (k Keeper) GetMQ(ctx sdk.Context, address sdk.AccAddress) uint64
func (k Keeper) SetMQ(ctx sdk.Context, address sdk.AccAddress, mq uint64)
func (k Keeper) InitializeMQ(ctx sdk.Context, address sdk.AccAddress) error

// Dispute
func (k Keeper) CreateDispute(ctx sdk.Context, msg types.MsgCreateDispute) (uint64, error)
func (k Keeper) GetDispute(ctx sdk.Context, id uint64) (types.Dispute, error)

// Mediation Events (unified)
func (k Keeper) AppendMediationEvent(ctx sdk.Context, disputeId uint64, event types.MediationEvent) error
func (k Keeper) GetMediationEvents(ctx sdk.Context, disputeId uint64) ([]types.MediationEvent, error)

// AI Mediation Actions
func (k Keeper) SendMessage(ctx sdk.Context, msg types.MsgSendMessage) error
func (k Keeper) SubmitEvidence(ctx sdk.Context, msg types.MsgSubmitEvidence) error
func (k Keeper) SubmitProposal(ctx sdk.Context, msg types.MsgSubmitProposal) error
func (k Keeper) SubmitScore(ctx sdk.Context, msg types.MsgSubmitScore) error
func (k Keeper) AcceptSettlement(ctx sdk.Context, msg types.MsgAcceptSettlement) error

// AI Decision (called by off-chain AI service)
func (k Keeper) RecordAIVerdict(ctx sdk.Context, disputeId uint64, verdict types.VerdictContent) error
func (k Keeper) EscalateToJury(ctx sdk.Context, disputeId uint64) error
func (k Keeper) DeclareForfeit(ctx sdk.Context, disputeId uint64, absentParty sdk.AccAddress) error

// Jury
func (k Keeper) SelectJury(ctx sdk.Context, disputeId uint64) error
func (k Keeper) CastJuryVote(ctx sdk.Context, msg types.MsgCastJuryVote) error
func (k Keeper) RecordAbsence(ctx sdk.Context, juror sdk.AccAddress) error

// MQ Redistribution
func (k Keeper) CalculateScoringDeviations(ctx sdk.Context, disputeId uint64) ([]types.ScoringDeviation, error)
func (k Keeper) RedistributeMQ(ctx sdk.Context, disputeId uint64) error

// Resolution
func (k Keeper) ResolveDispute(ctx sdk.Context, disputeId uint64) (*types.Resolution, error)
func (k Keeper) AppealDispute(ctx sdk.Context, msg types.MsgAppealDispute) error
```

---

## gRPC Queries

```protobuf
// query.proto

service Query {
    // MQ queries
    rpc MQ(QueryMQRequest) returns (QueryMQResponse);
    rpc MQByLevel(QueryMQByLevelRequest) returns (QueryMQByLevelResponse);
    rpc MQHistory(QueryMQHistoryRequest) returns (QueryMQHistoryResponse);
    rpc MQLeaderboard(QueryMQLeaderboardRequest) returns (QueryMQLeaderboardResponse);

    // Dispute queries
    rpc Dispute(QueryDisputeRequest) returns (QueryDisputeResponse);
    rpc Disputes(QueryDisputesRequest) returns (QueryDisputesResponse);
    rpc DisputesByParty(QueryDisputesByPartyRequest) returns (QueryDisputesByPartyResponse);

    // Mediation session queries
    rpc MediationEvents(QueryMediationEventsRequest) returns (QueryMediationEventsResponse);
    rpc MediationStatus(QueryMediationStatusRequest) returns (QueryMediationStatusResponse);

    // Evidence queries
    rpc Evidence(QueryEvidenceRequest) returns (QueryEvidenceResponse);

    // Proposal & scoring queries
    rpc Proposals(QueryProposalsRequest) returns (QueryProposalsResponse);
    rpc ProposalScores(QueryProposalScoresRequest) returns (QueryProposalScoresResponse);

    // Jury queries
    rpc JuryInfo(QueryJuryInfoRequest) returns (QueryJuryInfoResponse);
    rpc JuryDuties(QueryJuryDutiesRequest) returns (QueryJuryDutiesResponse);
}
```

---

## Module Dependencies

```
x/dispute (Core Module)
    │
    ├── Core Dependencies
    │   ├── x/identity  (Identity verification)
    │   ├── x/bank      (Token transfers)
    │   └── Base Types (01-base)
    │
    └── Depended By
        ├── x/compute   (Compute trading - dispute association)
        ├── x/task      (Task marketplace - dispute association)
        └── x/escrow    (Fund escrow - dispute lock)
```

---

## Complete Dispute Resolution Flow

```
┌──────────────┐
│ Order Complete│
└──────┬───────┘
       │
       ▼
┌──────────────┐     Dissatisfied     ┌──────────────────────┐
│ Rate Order   │────────────────────► │    File Dispute      │
└──────┬───────┘                      │  (stake MQ)  │
       │                              └───────────┬──────────┘
       │                                          │
    Satisfied│                                          │
       │                                          ▼
       │                              ┌──────────────────────────────┐
       │                              │   AI Mediation Session       │
       │                              │   (Unified Free Conversation)│
       │                              │                              │
       │                              │   ┌────────────────────────┐│
       │                              │   │ • Both parties chat    ││
       │                              │   │ • Submit evidence      ││
       │                              │   │ • AI guides discussion││
       │                              │   │ • Proposals emerge     ││
       │                              │   │ • Scoring when ready   ││
       │                              │   └────────────────────────┘│
       │                              │                              │
       │                              │         ┌────────────────┐   │
       │                              │         │ Any settlement │   │
       │                              │         │ at any point?  │   │
       │                              │         └───────┬────────┘   │
       │                              │                 │            │
       │                              │    ┌────────────┼───────────┐│
       │                              │    │            │           ││
       │                              │    ▼            ▼           ▼│
       │                              │ Settled    Cannot       Forfeit│
       │                              │    │      Settle          │   │
       │                              │    │         │           │   │
       │                              └────┼─────────┼───────────┼───┘
       │                                   │         │           │
       │                                   ▼         ▼           ▼
       │                            ┌─────────┐ ┌─────────┐ ┌─────────┐
       │                            │ Settled │ │  Juried │ │ Forfeit │
       │                            │ (done!) │ └────┬────┘ │(default)│
       │                            └─────────┘      │      └─────────┘
       │                                             │
       │                                             ▼
       │                                    ┌────────────────┐
       │                                    │   Jury Voting  │
       │                                    │                │
       │                                    │ 1. AI generates│
       │                                    │    summary     │
       │                                    │ 2. Select jury │
       │                                    │    (weighted)  │
       │                                    │ 3. Jurors score│
       │                                    │    proposals   │
       │                                    │ 4. Calculate   │
       │                                    │    consensus   │
       │                                    │ 5. Redistribute│
       │                                    │    MQ  │
       │                                    └───────┬────────┘
       │                                            │
       │                                            ▼
       │                                    ┌────────────────┐
       │                                    │   Resolved     │
       │                                    └───────┬────────┘
       │                                            │
       │                          ┌─────────────────┼─────────────────┐
       │                          │                 │                 │
       │                     No Appeal         Appeal          Rejected
       │                          │                 │                 │
       │                          ▼                 ▼                 │
       │                   ┌────────────┐   ┌────────────┐           │
       │                   │   Final    │   │  Appealed  │           │
       │                   └────────────┘   └─────┬──────┘           │
       │                                          │                  │
       │                                     Expired                  │
       │                                          │                  │
       │                                          ▼                  │
       │                                   ┌────────────┐            │
       │                                   │   Final    │◄───────────┘
       │                                   └─────┬──────┘
       │                                         │
       └─────────────────────────────────────────┘
                                                 │
                                                 ▼
                                        ┌────────────────┐
                                        │  Dispute End   │
                                        │ (funds + Rep)  │
                                        └────────────────┘
```

---

## Summary: Key Design Decisions

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Mediation | Unified, AI-led free conversation | Natural flow, better UX |
| States | 5 main states (filed→mediating→settled/juried/forfeit→resolved→final) | Simplified from 7 states |
| AI Role | Full control of mediation process | Intelligent timing for proposals/scoring |
| Settlement | Can happen at any moment | Maximizes early resolution |
| Jury | Last resort only | Reduces jury burden, respects party autonomy |
| MQ | Zero-sum redistribution | Maintains total MQ, incentivizes fairness |

---

[Previous Chapter: x/task](./08-task.md) | [Back to Index](./00-index.md) | [Next Chapter: x/identity →](./10-identity.md)
