# ShareTokens Modular Architecture v3.0

> **Core + Plugin Architecture for Service Marketplace**
>
> **Core Principle: Essential modules for every node, optional plugins for specific roles**

---

## Architecture Overview

```
+-----------------------------------------------------------------------------+
|                    ShareTokens Core + Plugin Architecture v3.0               |
+-----------------------------------------------------------------------------+

+-----------------------------------------------------------------------------+
|                    User Plugins (Optional)                                   |
|  +-----------------------------------------------------------------------+  |
|  |  P01: GenieBot Interface (Web UI)                                      |  |
|  |  - AI Chat Interface                                                   |  |
|  |  - Idea Incubation                                                     |  |
|  |  - Resource Matching                                                   |  |
|  |  - Workflow Execution                                                  |  |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
|                    Provider Plugins (Optional)                               |
|  +------------+  +------------------------+  +------------------------+    |
|  | P02: LLM   |  | P03: Agent Provider    |  | P04: Workflow Provider |    |
|  | Provider   |  | (OpenFang)             |  |                        |    |
|  | - API Key  |  | - Agent Execution      |  | - Workflow Execution   |    |
|  |   Hosting  |  | - Task Processing      |  | - Multi-step Tasks     |    |
|  +------------+  +------------------------+  +------------------------+    |
+-----------------------------------------------------------------------------+
|                    Core Modules (Every Node Required)                         |
|  +-----------------------------------------------------------------------+  |
|  |  C01: P2P Communication - CometBFT Built-in                            |  |
|  |  C02: Identity/Account - x/identity + Cosmos SDK Auth                  |  |
|  |  C03: Wallet - Cosmos SDK Auth + Keplr                                 |  |
|  +-----------------------------------------------------------------------+  |
|  |  C04: Service Market (Core Business)                                   |  |
|  |  +-----------------------------------------------------------------+  |  |
|  |  | Level 1: LLM API (x/compute)    - API Key托管、请求路由         |  |  |
|  |  | Level 2: Agent Service          - AI Agent任务执行              |  |  |
|  |  | Level 3: Workflow Service       - 多步骤工作流执行              |  |  |
|  |  | - Service Registration          - Service Discovery             |  |  |
|  |  | - Service Pricing               - Service Routing               |  |  |
|  |  +-----------------------------------------------------------------+  |  |
|  +-----------------------------------------------------------------------+  |
|  |  C05: Escrow Payment - x/escrow                                        |  |
|  |  C06: Trust System - x/trust (MQ + Dispute)                    |  |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
|                    Infrastructure (Not Developed - Use Existing)              |
|  +-----------------------------------------------------------------------+  |
|  |  CometBFT: P2P Network, BFT Consensus, ABCI                            |  |
|  |  Cosmos SDK: Auth, Bank, Staking, Params                               |  |
|  |  External: Keplr Wallet, Chainlink Oracle, GitHub API                  |  |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
```

---

## Design Principles

### 1. Core vs Plugin Separation

| Aspect | Core Modules | Plugin Modules |
|--------|--------------|----------------|
| Installation | Required on every node | Optional, role-specific |
| Dependencies | None (self-contained) | May depend on core modules |
| Chain State | Always present | Conditional on plugin |
| Development Priority | Phase 1 | Phase 2+ |

### 2. Do Not Reinvent the Wheel

| Function | Open Source Solution | Notes |
|----------|---------------------|-------|
| P2P Network | **CometBFT Built-in** | No need to develop libp2p |
| Consensus (BFT) | **CometBFT** | No need to implement consensus |
| Account/Signature | **Cosmos SDK Auth** | No need to develop Wallet module |
| Token Transfer | **Cosmos SDK Bank** | No need to implement |
| Staking | **Cosmos SDK Staking** | No need to implement |
| Price Oracle | **Chainlink Network** | Read existing data, no node deployment |
| IBC Cross-chain | **Cosmos IBC** | v2 feature, use directly |
| Wallet UI | **Keplr Integration** | Browser extension wallet |

### 3. Single Responsibility

Each module has one clear purpose.

### 4. Interface-First (Keeper Pattern)

Cosmos SDK Keeper interfaces defined before implementation.

### 5. Mockable Dependencies

Every module can be tested independently using Mock Keepers.

---

## Core Modules (Every Node Required)

Core modules are the minimum required components for any ShareTokens node to participate in the network.

---

## C01: P2P Communication

**Responsibility:** Peer-to-peer networking for consensus and data propagation.

**Implementation:** CometBFT Built-in (Not Developed)

**Status:** NO DEVELOPMENT REQUIRED

**Notes:**
- CometBFT provides production-ready P2P networking
- Includes peer discovery, connection management, message broadcasting
- Battle-tested in Cosmos ecosystem

---

## C02: Identity/Account - x/identity

**Responsibility:** Real-name verification, identity hash storage, Sybil resistance.

**Requirements:** ID-01 to ID-05

**Dependencies:** auth, bank (Cosmos SDK)

**Language:** Go (Cosmos SDK)

**Files:**
```
chain/x/identity/
  ├── keeper/
  │   ├── keeper.go           # Main keeper
  │   ├── grpc_query.go       # Query server
  │   └── msg_server.go       # Transaction server
  ├── types/
  │   ├── genesis.pb.go       # Genesis state
  │   ├── identity.pb.go      # Identity types
  │   ├── query.pb.go         # Query types
  │   ├── tx.pb.go            # Transaction types
  │   ├── keys.go             # Store keys
  │   └── errors.go           # Error definitions
  ├── genesis.go              # Genesis import/export
  └── module.go               # Module definition
```

**Keeper Interface:**
```go
// Keeper defines the identity module keeper
type Keeper interface {
    // Registration
    RegisterIdentity(ctx sdk.Context, owner sdk.AccAddress, identityProof IdentityProof) error
    VerifyIdentity(ctx sdk.Context, address sdk.AccAddress) (VerificationResult, error)
    RevokeIdentity(ctx sdk.Context, owner sdk.AccAddress) error

    // Queries
    GetIdentity(ctx sdk.Context, address sdk.AccAddress) (Identity, bool)
    GetIdentityByHash(ctx sdk.Context, hash []byte) (Identity, bool)
    HasDuplicateIdentity(ctx sdk.Context, hash []byte) bool

    // Merkle Proofs
    GetIdentityProof(ctx sdk.Context, address sdk.AccAddress) (MerkleProof, bool)
    VerifyIdentityProof(ctx sdk.Context, proof MerkleProof) bool

    // Dependencies
    GetAuthKeeper() authkeeper.AccountKeeper
    GetBankKeeper() bankkeeper.Keeper
}

// IdentityProof defines the proof for identity registration
type IdentityProof struct {
    Type       string        // "wechat" | "github" | "email" | "phone"
    Hash       []byte        // Hash of identity data (NOT the data itself)
    Attestation *Attestation // Optional oracle attestation
    Timestamp  time.Time
    Signature  []byte
}

// Identity defines the on-chain identity record
type Identity struct {
    Owner       sdk.AccAddress
    Hash        []byte
    Level       VerificationLevel // "none" | "basic" | "verified" | "premium"
    RegisteredAt time.Time
    Status      IdentityStatus // "active" | "revoked"
}
```

**Can be Mocked:** Yes (use MockIdentityKeeper)

**Complexity:** M (Medium)

---

## C03: Wallet

**Responsibility:** Account management, transaction signing, balance queries.

**Implementation:** Cosmos SDK Auth + Keplr (Not Developed)

**Status:** NO DEVELOPMENT REQUIRED

**Notes:**
- Cosmos SDK Auth module handles account state
- Keplr provides browser extension wallet UI
- Focus on integration, not development

---

## C04: Service Market (Core Business)

**Responsibility:** Service registration, discovery, pricing, routing across three service levels.

**Requirements:** COMP-01 to COMP-10, SVC-01 to SVC-06

**Dependencies:** auth, bank, escrow, identity

**Language:** Go (Cosmos SDK)

**Service Levels:**

| Level | Type | Description | Module |
|-------|------|-------------|--------|
| 1 | LLM API | API Key托管、请求路由 | x/compute |
| 2 | Agent Service | AI Agent任务执行 | OpenFang Integration |
| 3 | Workflow Service | 多步骤工作流执行 | Workflow Engine |

**Files:**
```
chain/x/compute/
  ├── keeper/
  │   ├── keeper.go           # Main keeper
  │   ├── provider.go         # Provider registration
  │   ├── request.go          # Compute request handling
  │   ├── response.go         # Response validation
  │   ├── market.go           # Compute offer discovery
  │   ├── service.go          # Service registry (Level 2 & 3)
  │   ├── pricing.go          # Service pricing
  │   ├── routing.go          # Service routing
  │   ├── grpc_query.go       # Query server
  │   └── msg_server.go       # Transaction server
  ├── types/
  │   ├── genesis.pb.go
  │   ├── compute.pb.go
  │   ├── service.pb.go       # Service types for all levels
  │   ├── query.pb.go
  │   ├── tx.pb.go
  │   ├── keys.go
  │   └── errors.go
  ├── genesis.go
  └── module.go
```

**Keeper Interface:**
```go
// Keeper defines the compute/service market module keeper
type Keeper interface {
    // === Provider Operations ===
    RegisterProvider(ctx sdk.Context, provider ProviderRegistration) error
    UpdateProviderOffer(ctx sdk.Context, providerAddr sdk.AccAddress, offer ComputeOffer) error
    PauseProvider(ctx sdk.Context, providerAddr sdk.AccAddress) error
    RevokeProvider(ctx sdk.Context, providerAddr sdk.AccAddress) error
    GetProvider(ctx sdk.Context, addr sdk.AccAddress) (Provider, bool)

    // === Compute Requests (Level 1 - LLM API) ===
    SubmitRequest(ctx sdk.Context, request ComputeRequest) (RequestResult, error)
    CancelRequest(ctx sdk.Context, requestId string, requester sdk.AccAddress) error
    GetRequest(ctx sdk.Context, requestId string) (ComputeRequest, bool)
    GetRequestsByRequester(ctx sdk.Context, requester sdk.AccAddress) []ComputeRequest
    SubmitResponse(ctx sdk.Context, response ComputeResponse) error
    VerifyResponse(ctx sdk.Context, response ComputeResponse) (VerificationResult, error)

    // === Service Registry (Level 2 & 3) ===
    RegisterService(ctx sdk.Context, service ServiceRegistration) error
    UpdateService(ctx sdk.Context, serviceId string, updates ServiceUpdates) error
    DeregisterService(ctx sdk.Context, serviceId string) error
    GetService(ctx sdk.Context, serviceId string) (Service, bool)

    // === Service Discovery ===
    GetServices(ctx sdk.Context, filter ServiceFilter) []Service
    FindService(ctx sdk.Context, criteria ServiceCriteria) ([]Service, error)

    // === Pricing ===
    SetPricing(ctx sdk.Context, serviceId string, pricing Pricing) error
    GetPricing(ctx sdk.Context, serviceId string) (Pricing, bool)
    CalculatePrice(ctx sdk.Context, request ServiceRequest) (sdk.Coins, error)

    // === Routing ===
    RouteRequest(ctx sdk.Context, request ServiceRequest) (RouteResult, error)
    GetAvailableProviders(ctx sdk.Context, serviceId string) ([]Provider, error)

    // === Market Discovery ===
    GetOffers(ctx sdk.Context, criteria OfferCriteria) []ComputeOffer
    MatchRequest(ctx sdk.Context, request ComputeRequest) ([]ComputeOffer, error)

    // Dependencies
    GetEscrowKeeper() escrowkeeper.Keeper
    GetBankKeeper() bankkeeper.Keeper
    GetIdentityKeeper() identitykeeper.Keeper
}

// ServiceRegistration for all service levels
type ServiceRegistration struct {
    Owner       sdk.AccAddress
    Name        string
    Level       ServiceLevel   // 1 = LLM, 2 = Agent, 3 = Workflow
    Category    string         // Service category
    Description string
    Pricing     Pricing
    Endpoint    string         // Service endpoint
    Metadata    map[string]string
}

// Service defines a registered service
type Service struct {
    Id          string
    Owner       sdk.AccAddress
    Name        string
    Level       ServiceLevel
    Category    string
    Description string
    Pricing     Pricing
    Endpoint    string
    Status      ServiceStatus  // "active" | "paused" | "deprecated"
    Rating      float64
    TotalJobs   int64
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Pricing defines service pricing model
type Pricing struct {
    Model       PricingModel   // "per_request" | "per_token" | "per_hour" | "fixed"
    BasePrice   sdk.Coins
    Unit        string         // "request" | "1k_tokens" | "hour"
    MinPrice    sdk.Coins
    MaxPrice    sdk.Coins
}

// ServiceLevel enum
type ServiceLevel int32

const (
    ServiceLevelLLM      ServiceLevel = 1  // Level 1: LLM API
    ServiceLevelAgent    ServiceLevel = 2  // Level 2: Agent Service
    ServiceLevelWorkflow ServiceLevel = 3  // Level 3: Workflow Service
)
```

**Can be Mocked:** Yes (use MockComputeKeeper with MockEscrowKeeper)

**Complexity:** XL (Extra Large - core business logic)

---

## C05: Escrow Payment - x/escrow

**Responsibility:** Payment escrow, release, partial release, dispute locking, ruling distribution.

**Requirements:** ESC-01 to ESC-06

**Dependencies:** auth, bank

**Language:** Go (Cosmos SDK)

**Files:**
```
chain/x/escrow/
  ├── keeper/
  │   ├── keeper.go           # Main keeper
  │   ├── escrow.go           # Escrow CRUD
  │   ├── release.go          # Release operations
  │   ├── dispute.go          # Dispute locking
  │   ├── grpc_query.go       # Query server
  │   └── msg_server.go       # Transaction server
  ├── types/
  │   ├── genesis.pb.go
  │   ├── escrow.pb.go
  │   ├── query.pb.go
  │   ├── tx.pb.go
  │   ├── keys.go
  │   └── errors.go
  ├── genesis.go
  └── module.go
```

**Keeper Interface:**
```go
// Keeper defines the escrow module keeper
type Keeper interface {
    // Escrow Operations
    CreateEscrow(ctx sdk.Context, params CreateEscrowParams) (Escrow, error)
    ReleaseEscrow(ctx sdk.Context, escrowId string, releaser sdk.AccAddress) error
    PartialRelease(ctx sdk.Context, escrowId string, amount sdk.Coins, releaser sdk.AccAddress) error
    CancelEscrow(ctx sdk.Context, escrowId string, creator sdk.AccAddress) error

    // Dispute Operations
    LockForDispute(ctx sdk.Context, escrowId string) error
    ResolveByRuling(ctx sdk.Context, escrowId string, ruling Ruling) error

    // Queries
    GetEscrow(ctx sdk.Context, escrowId string) (Escrow, bool)
    GetEscrowsByCreator(ctx sdk.Context, creator sdk.AccAddress) []Escrow
    GetEscrowsByBeneficiary(ctx sdk.Context, beneficiary sdk.AccAddress) []Escrow

    // Dependencies
    GetBankKeeper() bankkeeper.Keeper
}

// Escrow defines an escrow record
type Escrow struct {
    Id          string
    Creator     sdk.AccAddress
    Beneficiary sdk.AccAddress
    Amount      sdk.Coins
    Released    sdk.Coins
    Status      EscrowStatus    // "active" | "partial" | "locked" | "released" | "cancelled"
    LockedAt    *time.Time      // Set when locked for dispute
    CreatedAt   time.Time
    ExpiresAt   time.Time
}

// Ruling defines dispute resolution ruling
type Ruling struct {
    PlaintiffShare   sdk.Dec   // Percentage to plaintiff
    DefendantShare   sdk.Dec   // Percentage to defendant
    ArbitratorFee    sdk.Coins // Fee for arbitrators
}

// CreateEscrowParams for creating new escrow
type CreateEscrowParams struct {
    Creator     sdk.AccAddress
    Beneficiary sdk.AccAddress
    Amount      sdk.Coins
    Duration    time.Duration
    Metadata    string         // Reference to order/task
}
```

**Can be Mocked:** Yes (use MockEscrowKeeper)

**Complexity:** M (Medium)

---

## C06: Trust System - x/trust

**Responsibility:** Moral Quotient (MQ) scoring, dispute arbitration, jury selection, zero-sum redistribution, decay mechanism.

**Requirements:** MQ-01 to MQ-06, DISP-01 to DISP-08

**Dependencies:** auth, identity, escrow

**Language:** Go (Cosmos SDK)

**Files:**
```
chain/x/trust/
  ├── keeper/
  │   ├── keeper.go           # Main keeper
  │   ├── scoring.go          # MQ calculations
  │   ├── redistribution.go   # Zero-sum redistribution algorithm
  │   ├── decay.go            # Inactivity decay
  │   ├── jury.go             # Jury selection
  │   ├── dispute.go          # Dispute lifecycle
  │   ├── evidence.go         # Evidence submission
  │   ├── voting.go           # Jury voting
  │   ├── resolution.go       # Resolution and redistribution
  │   ├── appeal.go           # Appeal handling
  │   ├── grpc_query.go       # Query server
  │   └── msg_server.go       # Transaction server
  ├── types/
  │   ├── genesis.pb.go
  │   ├── trust.pb.go         # Trust types
  │   ├── mq.pb.go            # MQ types
  │   ├── dispute.pb.go       # Dispute types
  │   ├── query.pb.go
  │   ├── tx.pb.go
  │   ├── keys.go
  │   └── errors.go
  ├── genesis.go
  └── module.go
```

**Keeper Interface:**
```go
// Keeper defines the Trust System module keeper
type Keeper interface {
    // === MQ Operations ===
    GetMQ(ctx sdk.Context, address sdk.AccAddress) (MQ, bool)
    GetMQLevel(ctx sdk.Context, address sdk.AccAddress) (MQLevel, error)
    GetTotalMQ(ctx sdk.Context) sdk.Int
    InitializeMQ(ctx sdk.Context, address sdk.AccAddress) error
    RedistributeMQ(ctx sdk.Context, input RedistributionInput) (RedistributionResult, error)
    ApplyDecay(ctx sdk.Context, config DecayConfig) (DecayResult, error)

    // === Jury Operations ===
    SelectJury(ctx sdk.Context, disputeId string, size int, exclude []sdk.AccAddress) ([]sdk.AccAddress, error)
    GetVoteWeight(ctx sdk.Context, juror sdk.AccAddress) (sdk.Dec, error)

    // === Dispute Lifecycle ===
    CreateDispute(ctx sdk.Context, params DisputeParams) (Dispute, error)
    CancelDispute(ctx sdk.Context, disputeId string, creator sdk.AccAddress) error
    GetDispute(ctx sdk.Context, disputeId string) (Dispute, bool)
    GetDisputes(ctx sdk.Context, filter DisputeFilter) []Dispute

    // === Evidence ===
    SubmitEvidence(ctx sdk.Context, disputeId string, evidence Evidence) error
    GetEvidence(ctx sdk.Context, disputeId string) []Evidence

    // === Voting ===
    CastVote(ctx sdk.Context, disputeId string, juror sdk.AccAddress, vote Vote) error
    GetVotes(ctx sdk.Context, disputeId string) []DisputeVote

    // === Resolution ===
    ResolveDispute(ctx sdk.Context, disputeId string) (Resolution, error)

    // === Appeal ===
    AppealDispute(ctx sdk.Context, disputeId string, appellant sdk.AccAddress, reason string) error

    // Configuration
    GetConfig(ctx sdk.Context) TrustConfig
    SetConfig(ctx sdk.Context, config TrustConfig) error

    // Dependencies
    GetIdentityKeeper() identitykeeper.Keeper
    GetEscrowKeeper() escrowkeeper.Keeper
}

// MQ defines the Moral Quotient record for a user
type MQ struct {
    Address     sdk.AccAddress
    Current     sdk.Int       // Current MQ balance
    Locked      sdk.Int       // MQ locked in disputes
    Available   sdk.Int       // Available MQ
    Level       MQLevel // Calculated level
    LastUpdated time.Time
}

// MQLevel defines user levels based on MQ
type MQLevel int

const (
    LevelNewcomer MQLevel = iota  // 0-50
    LevelMember                            // 50-100
    LevelTrusted                           // 100-200
    LevelExpert                            // 200-500
    LevelGuardian                          // 500+
)

// Dispute defines a dispute record
type Dispute struct {
    Id              string
    OrderId         string
    Plaintiff       Party
    Defendant       Party
    Type            DisputeType   // "quality" | "delivery" | "payment" | "fraud"
    Title           string
    Description     string
    Amount          sdk.Coins
    Evidence        []Evidence
    Jury            []JurorInfo
    Votes           []DisputeVote
    Resolution      *Resolution
    Status          DisputeStatus // "open" | "negotiating" | "voting" | "resolved" | "appealed"
    AppealCount     int
    CreatedAt       time.Time
    ResolvedAt      *time.Time
}

// Vote defines a juror's vote
type Vote struct {
    Verdict   Verdict   // "plaintiff" | "defendant" | "neutral"
    Reasoning string
}

// Resolution defines dispute resolution
type Resolution struct {
    Verdict              Verdict
    PlaintiffShare       sdk.Dec
    DefendantShare       sdk.Dec
    MQRedistribution RedistributionResult
    ResolvedAt           time.Time
}
```

**Can be Mocked:** Yes (use MockTrustKeeper)

**Complexity:** L (Large - combined module)

---

## Provider Plugins (Optional)

Provider plugins are optional modules that extend node functionality to provide specific services.

---

## P02: LLM Provider

**Responsibility:** Host and manage API keys for LLM services (OpenAI, Anthropic, etc.).

**Requirements:** COMP-01 to COMP-04

**Dependencies:** x/compute (core)

**Language:** TypeScript (Offchain Service)

**Files:**
```
plugins/llm-provider/
  ├── src/
  │   ├── index.ts
  │   ├── key-management/
  │   │   ├── store.ts          # Encrypted key storage
  │   │   ├── rotate.ts         # Key rotation
  │   │   └── audit.ts          # Key usage audit
  │   ├── providers/
  │   │   ├── base.ts           # Base provider interface
  │   │   ├── openai.ts         # OpenAI integration
  │   │   ├── anthropic.ts      # Anthropic integration
  │   │   ├── google.ts         # Google AI integration
  │   │   └── azure.ts          # Azure OpenAI integration
  │   ├── routing/
  │   │   ├── router.ts         # Request routing
  │   │   ├── load-balancer.ts  # Load balancing
  │   │   └── fallback.ts       # Fallback handling
  │   ├── billing/
  │   │   ├── meter.ts          # Token metering
  │   │   └── report.ts         # Usage reporting
  │   └── security/
  │       ├── encryption.ts     # Key encryption
  │       └── access-control.ts # Access control
  ├── package.json
  └── tsconfig.json
```

**Interface:**
```typescript
// LLM Provider Interface
export interface ILLMProviderPlugin {
  // Key Management
  registerAPIKey(provider: string, encryptedKey: EncryptedKey): Promise<void>
  rotateKey(provider: string): Promise<void>
  revokeKey(provider: string): Promise<void>

  // Request Handling
  processRequest(request: LLMRequest): Promise<LLMResponse>
  streamRequest(request: LLMRequest): AsyncIterable<LLMChunk>

  // Billing
  getTokenUsage(requestId: string): Promise<TokenUsage>
  getBillingReport(period: DateRange): Promise<BillingReport>

  // Health
  healthCheck(): Promise<HealthStatus>
  getProviderStatus(): Promise<ProviderStatus[]>
}

export interface LLMRequest {
  id: string
  model: string
  prompt: string
  maxTokens?: number
  temperature?: number
  metadata: Record<string, any>
}

export interface LLMResponse {
  id: string
  content: string
  usage: TokenUsage
  latency: number
  provider: string
}
```

**Can be Mocked:** Yes (echo responses)

**Complexity:** L (Large)

**Installation:** Optional for LLM API providers

---

## P03: Agent Provider (OpenFang)

**Responsibility:** Execute AI Agent tasks using OpenFang framework integration.

**Requirements:** AGENT-01 to AGENT-06

**Dependencies:** x/compute (core), OpenFang Runtime

**Language:** TypeScript + OpenFang

**Files:**
```
plugins/agent-provider/
  ├── src/
  │   ├── index.ts
  │   ├── openfang/
  │   │   ├── client.ts         # OpenFang SDK client
  │   │   ├── agents/           # Agent configurations
  │   │   │   ├── coder.ts      # Coder Agent
  │   │   │   ├── researcher.ts # Researcher Agent
  │   │   │   ├── writer.ts     # Writer Agent
  │   │   │   └── custom.ts     # Custom Agents
  │   │   └── hands/            # OpenFang Hands
  │   │       ├── collector.ts  # Collector Hand
  │   │       ├── content.ts    # Content Hand
  │   │       └── lead.ts       # Lead Hand
  │   ├── task-executor/
  │   │   ├── executor.ts       # Task execution engine
  │   │   ├── scheduler.ts      # Task scheduling
  │   │   └── monitor.ts        # Execution monitoring
  │   ├── security/
  │   │   ├── sandbox.ts        # WASM sandbox
  │   │   └── limits.ts         # Resource limits
  │   └── reporting/
  │       ├── status.ts         # Status reporting
  │       └── metrics.ts        # Performance metrics
  ├── package.json
  └── tsconfig.json
```

**Interface:**
```typescript
// Agent Provider Interface
export interface IAgentProviderPlugin {
  // Agent Management
  registerAgent(config: AgentConfig): Promise<Agent>
  updateAgent(agentId: string, config: AgentConfig): Promise<Agent>
  deregisterAgent(agentId: string): Promise<void>
  getAgent(agentId: string): Promise<Agent>

  // Task Execution
  submitTask(task: AgentTask): Promise<TaskResult>
  getTaskStatus(taskId: string): Promise<TaskStatus>
  cancelTask(taskId: string): Promise<void>

  // OpenFang Integration
  getAvailableAgents(): Promise<AgentTemplate[]>
  getAgentCapabilities(agentId: string): Promise<Capability[]>

  // Monitoring
  getExecutionMetrics(agentId: string): Promise<ExecutionMetrics>
  getHealthStatus(): Promise<HealthStatus>
}

export interface AgentTask {
  id: string
  agentId: string
  input: TaskInput
  constraints: TaskConstraints
  callback?: string
}

export interface TaskResult {
  taskId: string
  status: 'completed' | 'failed' | 'timeout'
  output: any
  executionTime: number
  tokensUsed: number
}
```

**OpenFang Agent Templates (28+):**

| Agent | Purpose | ShareTokens Use Case |
|-------|---------|---------------------|
| Coder | Code generation | Software development workflow |
| Researcher | Research analysis | Idea evaluation |
| Writer | Content creation | Content workflow |
| Architect | System design | Technical architecture |
| Debugger | Code debugging | Software development |
| Reviewer | Code review | Quality assurance |

**OpenFang Hands (7):**

| Hand | Purpose | ShareTokens Use Case |
|------|---------|---------------------|
| Collector | Data collection | Market research |
| Clip | Content clipping | Content curation |
| Lead | Sales leads | Resource matching |
| Content | Content management | Publishing |
| Trade | Trading | Token trading |
| Calendar | Scheduling | Task scheduling |
| Notify | Notifications | User alerts |

**Can be Mocked:** Yes (simulated agent responses)

**Complexity:** XL (Extra Large)

**Installation:** Optional for Agent service providers

---

## P04: Workflow Provider

**Responsibility:** Execute multi-step workflow tasks.

**Requirements:** WF-01 to WF-07

**Dependencies:** x/compute (core), AI Service, GitHub Integration

**Language:** TypeScript

**Files:**
```
plugins/workflow-provider/
  ├── src/
  │   ├── index.ts
  │   ├── workflows/
  │   │   ├── base.ts           # Base workflow class
  │   │   ├── software.ts       # Software development workflow
  │   │   ├── content.ts        # Content creation workflow
  │   │   ├── business.ts       # Business planning workflow
  │   │   └── service.ts        # Life service workflow
  │   ├── executor/
  │   │   ├── runner.ts         # Workflow runner
  │   │   ├── state.ts          # State management
  │   │   └── recovery.ts       # Failure recovery
  │   ├── nodes/
  │   │   ├── auto.ts           # Automated nodes
  │   │   ├── human.ts          # Human approval nodes
  │   │   └── external.ts       # External service nodes
  │   └── monitoring/
  │       ├── progress.ts       # Progress tracking
  │       └── alerts.ts         # Alert system
  ├── package.json
  └── tsconfig.json
```

**Interface:**
```typescript
// Workflow Provider Interface
export interface IWorkflowProviderPlugin {
  // Workflow Registration
  registerWorkflow(config: WorkflowConfig): Promise<Workflow>

  // Execution
  startWorkflow(workflowId: string, input: WorkflowInput): Promise<WorkflowExecution>
  getExecutionStatus(executionId: string): Promise<ExecutionStatus>
  cancelExecution(executionId: string): Promise<void>

  // Node Operations
  submitNodeResult(executionId: string, nodeId: string, result: NodeResult): Promise<void>
  approveNode(executionId: string, nodeId: string, approver: string): Promise<void>
  rejectNode(executionId: string, nodeId: string, reason: string): Promise<void>

  // Monitoring
  getActiveExecutions(): Promise<WorkflowExecution[]>
  getExecutionHistory(limit: number): Promise<WorkflowExecution[]>
}

export interface WorkflowExecution {
  id: string
  workflowId: string
  status: 'pending' | 'running' | 'waiting_approval' | 'completed' | 'failed' | 'cancelled'
  nodes: ExecutionNode[]
  currentNode: number
  input: WorkflowInput
  output?: any
  startedAt: Date
  completedAt?: Date
}

export interface ExecutionNode {
  id: string
  type: 'auto' | 'human_approval' | 'external'
  name: string
  status: 'pending' | 'running' | 'waiting' | 'completed' | 'failed'
  input: any
  output?: any
  error?: string
}
```

**Can be Mocked:** Yes (simple state machine)

**Complexity:** L (Large)

**Installation:** Optional for Workflow service providers

---

## User Plugins (Optional)

User plugins provide enhanced user interfaces and experiences.

---

## P01: GenieBot Interface (Web UI)

**Responsibility:** AI chat interface, idea incubation, resource matching, workflow execution UI.

**Requirements:** LAMP-01 to LAMP-05, UI-01 to UI-06

**Dependencies:** All core modules, offchain services

**Language:** React + TypeScript

**Files:**
```
frontend/
  ├── src/
  │   ├── core/
  │   │   ├── wallet/
  │   │   │   ├── KeplrProvider.tsx
  │   │   │   ├── useWallet.ts
  │   │   │   └── ConnectButton.tsx
  │   │   ├── layout/
  │   │   │   ├── AppLayout.tsx
  │   │   │   ├── Sidebar.tsx
  │   │   │   └── Header.tsx
  │   │   └── theme/
  │   │       ├── ThemeProvider.tsx
  │   │       └── theme.ts
  │   ├── chat/
  │   │   ├── ChatWindow.tsx
  │   │   ├── MessageList.tsx
  │   │   ├── MessageItem.tsx
  │   │   ├── InputArea.tsx
  │   │   └── IdeaCard/
  │   │       ├── IdeaCard.tsx
  │   │       ├── EvalReport.tsx
  │   │       └── TokenEstimate.tsx
  │   ├── market/
  │   │   ├── compute/
  │   │   │   ├── ComputeMarket.tsx
  │   │   │   └── ProviderCard.tsx
  │   │   ├── agent/
  │   │   │   ├── AgentMarket.tsx
  │   │   │   └── AgentCard.tsx
  │   │   └── workflow/
  │   │       ├── WorkflowMarket.tsx
  │   │       └── WorkflowCard.tsx
  │   ├── dispute/
  │   │   ├── DisputeForm.tsx
  │   │   ├── DisputeList.tsx
  │   │   ├── EvidenceUpload.tsx
  │   │   └── VotingPanel.tsx
  │   └── profile/
  │       ├── ProfilePage.tsx
  │       ├── MQDisplay.tsx
  │       └── ServiceHistory.tsx
  ├── package.json
  └── tsconfig.json
```

**Interface:**
```typescript
// GenieBot UI Interface
export interface IXiaoDengUI {
  // Chat
  startConversation(): Promise<Conversation>
  sendMessage(conversationId: string, content: string): Promise<Message>

  // Idea
  submitIdea(idea: IdeaInput): Promise<Idea>
  evaluateIdea(ideaId: string): Promise<Evaluation>

  // Service Market
  browseServices(filter: ServiceFilter): Promise<Service[]>
  requestService(serviceId: string, request: ServiceRequest): Promise<ServiceOrder>

  // Dispute
  createDispute(orderId: string, reason: string): Promise<Dispute>
  voteOnDispute(disputeId: string, vote: Vote): Promise<void>
}

export interface Conversation {
  id: string
  messages: Message[]
  status: 'active' | 'completed'
  idea?: Idea
  workflow?: Workflow
}
```

**Complexity:** XL (Extra Large)

**Installation:** Optional for end users

---

## Dependency Graph

```
+-----------------------------------------------------------------------------+
|                           Dependency Graph                                   |
+-----------------------------------------------------------------------------+

NOT DEVELOPED (Existing):
+-----------------------------------------------------------------------------+
|  Cosmos SDK Modules        |  CometBFT           |  External               |
|  +---------------------+   |  +---------------+  |  +-------------------+  |
|  | auth (Account)      |   |  | P2P Network   |  |  | Keplr Wallet      |  |
|  | bank (Transfer)     |   |  | BFT Consensus |  |  | Chainlink Oracle  |  |
|  | staking (Staking)   |   |  | ABCI          |  |  | GitHub API        |  |
|  | params (Config)     |   |  +---------------+  |  | LLM APIs          |  |
|  +---------------------+   |                     |  +-------------------+  |
+-----------------------------------------------------------------------------+

CORE MODULES (Go - Every Node):
+-----------------------------------------------------------------------------+
|                                                                             |
|  +---------+     +---------+     +---------+                               |
|  | auth    |     | bank    |     | staking |  (Cosmos SDK - Use Directly)  |
|  +----+----+     +----+----+     +---------+                               |
|       |               |                                                     |
|       +-------+-------+                                                     |
|               |                                                             |
|    +----------+----------+                                                  |
|    |                     |                                                  |
|    v                     v                                                  |
|  +------+            +-------+                                              |
|  |x/iden|            |x/escro|     Wave 1a (Parallel)                      |
|  |tity  |            |w      |                                              |
|  +--+---+            +---+---+                                              |
|     |                    |                                                  |
|     v                    |                                                  |
|  +------+                |                                                  |
|  |x/trust|<---------------+     Wave 1b (Trust System)                         |
|  +--+---+                                                                    |
|     |                                                                        |
|     v                                                                        |
|  +----------+                                                                |
|  |x/compute |   Service Market (Core Business) - Wave 1a                    |
|  |(Service  |   - Level 1: LLM API                                          |
|  | Market)  |   - Level 2: Agent Service                                    |
|  +----------+   - Level 3: Workflow Service                                 |
|                                                                            |
|                                                                             |
|  +----------+                                                                |
|  |x/compute |   Service Market (Core Business) - Wave 1a                    |
|  |(Service  |   - Level 1: LLM API                                          |
|  | Market)  |   - Level 2: Agent Service                                    |
|  +----------+   - Level 3: Workflow Service                                 |
|                                                                             |
+-----------------------------------------------------------------------------+

PROVIDER PLUGINS (TypeScript - Optional):
+-----------------------------------------------------------------------------+
|                                                                             |
|  +------------+        +----------------+        +------------------+       |
|  | P02: LLM   |        | P03: Agent     |        | P04: Workflow    |       |
|  | Provider   |        | Provider       |        | Provider         |       |
|  +-----+------+        +-------+--------+        +--------+---------+       |
|        |                       |                          |                 |
|        +-----------------------+--------------------------+                 |
|                                |                                           |
|                                v                                           |
|                        +-------------+                                     |
|                        | x/compute   |  (Core Module)                      |
|                        +-------------+                                     |
|                                                                             |
+-----------------------------------------------------------------------------+

USER PLUGINS (React - Optional):
+-----------------------------------------------------------------------------+
|                                                                             |
|  +--------------------------+                                               |
|  | P01: GenieBot Interface (Web UI)   |                                               |
|  +------------+-------------+                                               |
|               |                                                             |
|               v                                                             |
|  +------------+-------------+                                               |
|  | All Core Modules          |  (Uses all core chain modules)              |
|  +--------------------------+                                               |
|                                                                             |
+-----------------------------------------------------------------------------+
```

---

## Module Summary

### Core Modules (Go - Cosmos SDK)

| Module | ID | Dependencies | Complexity | Required | Requirements |
|--------|-----|--------------|------------|----------|--------------|
| P2P Network | C01 | None | - (CometBFT) | YES | - |
| Identity | C02 | auth, bank | M | YES | ID-01 to ID-05 |
| Wallet | C03 | None | - (Keplr) | YES | - |
| Service Market | C04 | auth, bank, escrow, identity | XL | YES | COMP, SVC |
| Escrow | C05 | auth, bank | M | YES | ESC-01 to ESC-06 |
| Trust System | C06 | auth, identity, escrow | L | YES | MQ-01 to MQ-06, DISP-01 to DISP-08 |

### Provider Plugins (TypeScript)

| Module | ID | Dependencies | Complexity | Required | Requirements |
|--------|-----|--------------|------------|----------|--------------|
| LLM Provider | P02 | x/compute | L | NO | COMP-01 to 04 |
| Agent Provider | P03 | x/compute, OpenFang | XL | NO | AGENT-01 to 06 |
| Workflow Provider | P04 | x/compute | L | NO | WF-01 to 07 |

### User Plugins (React + TypeScript)

| Module | ID | Dependencies | Complexity | Required |
|--------|-----|--------------|------------|----------|
| GenieBot Interface (Web UI) | P01 | All core | XL | NO |

---

## Parallel Development Strategy

### Core Modules Development

```
+-----------------------------------------------------------------------------+
| Wave 1a (Week 1-2): 4 Chain Modules in Parallel                             |
|   [C02: x/identity] [C04: x/compute] [C05: x/escrow]                        |
+-----------------------------------------------------------------------------+
| Wave 1b (Week 2-3): 3 Chain Modules in Parallel                             |
|   [C06: x/mq] [C07: x/dispute]                                              |
+-----------------------------------------------------------------------------+
```

### Plugin Development

```
+-----------------------------------------------------------------------------+
| Wave 2a (Week 3-5): Provider Plugins in Parallel                            |
|   [P02: LLM Provider] [P03: Agent Provider] [P04: Workflow Provider]       |
+-----------------------------------------------------------------------------+
| Wave 2b (Week 5-7): User Plugin                                             |
|   [P01: GenieBot Interface]                                                    |
+-----------------------------------------------------------------------------+
```

### Agent Assignment Recommendation

| Agents | Wave 1a (3 modules) | Wave 1b (2 modules) | Wave 2a (3 plugins) | Wave 2b |
|--------|---------------------|---------------------|---------------------|---------|
| 1 | All sequential | All sequential | All sequential | UI |
| 2 | 2+1 | 1+1 | 2+1 | UI |
| 3 | 1 each | 1+1 | 1 each | UI |

---

## Complexity Summary

### Core Modules (Required)

| Module | Complexity | Est. Lines | Language |
|--------|------------|------------|----------|
| P2P Network | - | 0 (CometBFT) | - |
| Identity | M | ~600 | Go |
| Wallet | - | 0 (Keplr) | - |
| Service Market | XL | ~2,000 | Go |
| Escrow | M | ~600 | Go |
| Trust System | L | ~1,600 | Go |

**Core Total:** ~4,800 lines

### Provider Plugins (Optional)

| Module | Complexity | Est. Lines | Language |
|--------|------------|------------|----------|
| LLM Provider | L | ~1,500 | TypeScript |
| Agent Provider | XL | ~2,500 | TypeScript |
| Workflow Provider | L | ~1,500 | TypeScript |

**Provider Total:** ~5,500 lines

### User Plugins (Optional)

| Module | Complexity | Est. Lines | Language |
|--------|------------|------------|----------|
| GenieBot Interface | XL | ~3,000 | TypeScript |

**User Total:** ~3,000 lines

**Grand Total:** ~13,500 lines (vs 18,500 in v1 - 27% reduction)

---

## What We Don't Develop

| Function | Solution | Why We Don't Develop |
|----------|----------|---------------------|
| P2P Network | CometBFT Built-in | Battle-tested, no need to reinvent |
| BFT Consensus | CometBFT | Industry standard |
| Account Management | Cosmos SDK Auth | Full-featured |
| Token Transfer | Cosmos SDK Bank | Includes all needed functionality |
| Staking | Cosmos SDK Staking | Validator delegation |
| Wallet UI | Keplr Integration | Users already have it |
| Price Oracle | Chainlink Network | Read existing feeds |
| IBC | Cosmos IBC | Built-in cross-chain |

**Estimated Lines Saved:** ~8,000+ lines

---

## Testing Strategy

### Unit Testing (Per Module)

Each module should have:
- `keeper/keeper_test.go` - Keeper unit tests with mocked dependencies
- `types/` - Type validation tests

### Integration Testing

```
chain/tests/
  ├── integration/
  │   ├── identity_mq_test.go     # Identity + MQ integration
  │   ├── service_escrow_test.go  # Service + Escrow integration
  │   ├── dispute_flow_test.go    # Full dispute flow
  │   └── service_market_test.go  # Service marketplace flow
  └── e2e/
      └── full_flow_test.go       # End-to-end scenarios
```

### Plugin Testing

```
plugins/
  ├── llm-provider/tests/
  ├── agent-provider/tests/
  └── workflow-provider/tests/
```

---

## File Structure

```
sharetokens/
+-- chain/                        # Cosmos SDK Chain (Go)
|   +-- app/
|   |   +-- app.go               # Application configuration
|   |   +-- encoding.go          # Encoding configuration
|   +-- cmd/
|   |   +-- sharetokensd/        # Chain binary
|   +-- x/                       # Core modules
|   |   +-- identity/            # C02: Identity
|   |   +-- compute/             # C04: Service Market
|   |   +-- escrow/              # C05: Escrow
|   |   +-- mq/                  # C06: Moral Quotient
|   |   +-- dispute/             # C07: Dispute
|   +-- proto/                   # Protobuf definitions
|   +-- testnets/                # Test network configurations
|   +-- go.mod
|   +-- go.sum
|
+-- plugins/                     # Optional Plugins (TypeScript)
|   +-- llm-provider/            # P02: LLM Provider
|   +-- agent-provider/          # P03: Agent Provider
|   +-- workflow-provider/       # P04: Workflow Provider
|
+-- frontend/                    # User Plugin (React + TypeScript)
|   +-- src/
|   |   +-- core/
|   |   +-- chat/
|   |   +-- market/
|   |   +-- dispute/
|   +-- package.json
|   +-- tsconfig.json
|
+-- proto/                       # Shared protobuf definitions
|
+-- docs/                        # Documentation
|
+-- .planning/                   # Planning documents
    +-- PROJECT.md
    +-- REQUIREMENTS.md
    +-- ROADMAP.md
    +-- MODULES.md
    +-- STATE.md
```

---

## Next Steps

1. **Wave 0:** Initialize Cosmos SDK chain with Ignite CLI
2. **Protobuf Definitions:** Define all message types in proto/
3. **Keeper Interfaces:** Define all Keeper interfaces for mocking
4. **Mock Implementations:** Create mock keepers for parallel development
5. **CI/CD Setup:** Configure parallel builds and tests

---

*Architecture Version: 3.0*
*Created: 2026-03-02*
*Updated: 2026-03-02 - Reorganized into Core + Plugin architecture*
*Key Change: Clear separation between required core modules and optional plugins*
