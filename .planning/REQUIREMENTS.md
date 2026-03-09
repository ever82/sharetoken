# Requirements: ShareTokens

**Defined:** 2025-03-01
**Core Value:** Every great idea deserves the tokens to make it real.

---

## Architecture Overview

### Core Modules (Required for every node)
1. P2P Communication
2. Identity
3. Wallet
4. Service Marketplace (Three-tier: LLM/Agent/Workflow)
5. Escrow
6. Trust System (MQ + Dispute Arbitration)

### Optional Modules (Plugins)
- **Service Provider Plugins**: LLM API Key Hosting, Agent Executor, Workflow Executor
- **User Plugins**: GenieBot Interface

---

## v1 Requirements

---

## Core Modules

### P2P Communication (Core Module - Foundation)

- [ ] **P2P-01**: Users can start a node and connect to the network
- [ ] **P2P-02**: Nodes can discover other nodes (DHT)
- [ ] **P2P-03**: Nodes can broadcast messages (PubSub)
- [ ] **P2P-04**: Nodes can traverse NAT to establish connections
- [ ] **P2P-05**: Nodes can verify other nodes' identity

### Identity (Core Module - Security)

- [ ] **ID-01**: Users can register real-name identity (WeChat/GitHub/Phone)
- [ ] **ID-02**: System can verify identity uniqueness (prevent duplicate registration)
- [ ] **ID-03**: Identity information stores only hash values (privacy protection)
- [ ] **ID-04**: Users can revoke registered identity
- [ ] **ID-05**: System can verify Merkle proofs

### Wallet (Core Module - Foundation)

- [ ] **WALL-01**: Users can create new wallets
- [ ] **WALL-02**: Users can import existing wallets (mnemonic)
- [ ] **WALL-03**: Users can export wallet backup
- [ ] **WALL-04**: Users can view balance
- [ ] **WALL-05**: Wallet can sign transactions

### Service Marketplace (Core Module - Three-tier Service Architecture)

**Service Registration and Discovery**

- [ ] **SVC-01**: Users can register services to marketplace (LLM/Agent/Workflow)
- [ ] **SVC-02**: Users can discover available services (filter by type, capability)
- [ ] **SVC-03**: Users can view service details (capability description, price, provider MQ)
- [ ] **SVC-04**: Service providers can update service information
- [ ] **SVC-05**: Service providers can delist services

**Service Pricing**

- [ ] **SVC-06**: Service providers can set pricing strategies (per-use/subscription)
- [ ] **SVC-07**: System can convert service price to STT (based on exchange rate)
- [ ] **SVC-08**: Users can view service price history
- [ ] **SVC-09**: System supports dynamic pricing (based on supply/demand)

**Request Routing**

- [ ] **SVC-10**: System can route requests to appropriate service providers
- [ ] **SVC-11**: System can match best provider based on MQ, price, capability
- [ ] **SVC-12**: System supports request retry and failover
- [ ] **SVC-13**: System can load balance requests

**LLM Service Layer**

- [ ] **LLM-01**: Users can host LLM API Keys (encrypted storage)
- [ ] **LLM-02**: Users can query available LLM model list
- [ ] **LLM-03**: Users can initiate LLM requests
- [ ] **LLM-04**: Providers can execute LLM requests and return results
- [ ] **LLM-05**: System can verify execution results (token usage)
- [ ] **LLM-06**: System can auto-settle payments (Token -> Provider)
- [ ] **LLM-07**: Users can set API Key pricing strategy
- [ ] **LLM-08**: Users can pause/revoke API Key hosting

**Agent Service Layer**

- [ ] **AGENT-01**: Users can register Agent services (based on OpenFang)
- [ ] **AGENT-02**: System supports multiple Agent types (Coder/Researcher/Writer)
- [ ] **AGENT-03**: Users can initiate Agent task requests
- [ ] **AGENT-04**: Agent executors can execute tasks and return results
- [ ] **AGENT-05**: System can verify Agent execution results
- [ ] **AGENT-06**: System can auto-settle Agent service fees
- [ ] **AGENT-07**: Users can view Agent execution status
- [ ] **AGENT-08**: Agent services support custom templates

**Workflow Service Layer**

- [ ] **WF-01**: System can identify Workflow type for requests
- [ ] **WF-02**: Users can register Workflow services (software/content/business/life)
- [ ] **WF-03**: Workflow executors can execute software development workflows
- [ ] **WF-04**: Workflow executors can execute content creation workflows
- [ ] **WF-05**: Workflow executors can execute business planning workflows
- [ ] **WF-06**: Workflow executors can execute life service workflows
- [ ] **WF-07**: Users can view Workflow execution status
- [ ] **WF-08**: Workflows support manual review nodes
- [ ] **WF-09**: System can verify Workflow execution results
- [ ] **WF-10**: System can auto-settle Workflow service fees

**Exchange Rate Service**

- [ ] **PRICE-01**: System can get STT/USD exchange rate (Chainlink)
- [ ] **PRICE-02**: System can convert service price to STT
- [ ] **PRICE-03**: System can cache exchange rate data
- [ ] **PRICE-04**: System can notify users when exchange rate changes

### Escrow (Core Module - Trust)

- [ ] **ESC-01**: Users can create escrow for transactions
- [ ] **ESC-02**: System can release escrow funds when task completes
- [ ] **ESC-03**: Users can request partial release
- [ ] **ESC-04**: System can lock escrow during disputes
- [ ] **ESC-05**: System can distribute escrow after dispute resolution
- [ ] **ESC-06**: Users can view escrow status

### Trust System (Core Module - MQ + Dispute)

**Moral Quotient (MQ) Scoring**

- [ ] **MQ-01**: System can initialize user MQ (initial value 100)
- [ ] **MQ-02**: System can redistribute MQ after disputes (zero-sum game)
- [ ] **MQ-03**: System can apply MQ decay (inactive users)
- [ ] **MQ-04**: System can select jury based on MQ
- [ ] **MQ-05**: User voting weight depends on MQ value
- [ ] **MQ-06**: System can calculate MQ tiers

**Dispute Arbitration**

- [ ] **DISP-01**: Users can initiate disputes
- [ ] **DISP-02**: Users can submit evidence
- [ ] **DISP-03**: Dispute parties can negotiate
- [ ] **DISP-04**: System can form a jury
- [ ] **DISP-05**: Jury members can vote
- [ ] **DISP-06**: System can rule based on voting results
- [ ] **DISP-07**: Ruling results can trigger MQ redistribution
- [ ] **DISP-08**: Users can view dispute history

---

## Optional Modules - Service Provider Plugins

### LLM API Key Hosting Plugin

*Inherits LLM-01~08 requirements*

- [ ] **LLM-PLUGIN-01**: Plugin can securely store and manage API Keys
- [ ] **LLM-PLUGIN-02**: Plugin can monitor API Key usage
- [ ] **LLM-PLUGIN-03**: Plugin supports multiple provider API Keys (OpenAI, Anthropic, etc.)

### Agent Executor Plugin (OpenFang Integration)

*Inherits AGENT-01~08 requirements*

- [ ] **AGENT-PLUGIN-01**: Plugin can install and configure OpenFang
- [ ] **AGENT-PLUGIN-02**: Plugin can use Coder Agent for software development tasks
- [ ] **AGENT-PLUGIN-03**: Plugin can use Researcher Agent for research and analysis
- [ ] **AGENT-PLUGIN-04**: Plugin can use Writer Agent for content creation
- [ ] **AGENT-PLUGIN-05**: Plugin can use Collector Hand for data collection and monitoring
- [ ] **AGENT-PLUGIN-06**: Plugin can use Content Hand for content creation workflows
- [ ] **AGENT-PLUGIN-07**: Plugin can use Lead Hand for resource matching
- [ ] **AGENT-PLUGIN-08**: Plugin can monitor Agent status via Dashboard
- [ ] **AGENT-PLUGIN-09**: Plugin can use secure sandbox to protect API Keys (16-layer security)
- [ ] **AGENT-PLUGIN-10**: Plugin can create custom Agent templates and Hand workflows

### Workflow Executor Plugin

*Inherits WF-01~10 requirements*

- [ ] **WF-PLUGIN-01**: Plugin can execute GitHub integration workflows
- [ ] **WF-PLUGIN-02**: Users can connect GitHub account
- [ ] **WF-PLUGIN-03**: Plugin can create code repository based on request
- [ ] **WF-PLUGIN-04**: Plugin can generate code and create PRs
- [ ] **WF-PLUGIN-05**: Users can review PRs via GenieBot
- [ ] **WF-PLUGIN-06**: Plugin can auto-deploy merged code

---

## Optional Modules - User Plugins

### GenieBot AI Chat Plugin (Application Layer - User Interface)

**Basic Chat**

- [ ] **LAMP-01**: Users can have natural language conversations with GenieBot
- [ ] **LAMP-02**: GenieBot can identify user intent (idea/question/command)
- [ ] **LAMP-03**: GenieBot can maintain conversation context (multi-turn dialogue)
- [ ] **LAMP-04**: Users can view conversation history
- [ ] **LAMP-05**: GenieBot can display task processing progress

**Idea Processing**

- [ ] **IDEA-PROC-01**: GenieBot can auto-classify ideas (software/content/business/life)
- [ ] **IDEA-PROC-02**: GenieBot can refine user ideas through dialogue
- [ ] **IDEA-PROC-03**: System can call multiple AIs to evaluate token requirements
- [ ] **IDEA-PROC-04**: GenieBot can display idea evaluation report
- [ ] **IDEA-PROC-05**: Users can save and manage idea drafts

**Multi-AI Evaluation**

- [ ] **EVAL-01**: System can call multiple AI models to evaluate idea value
- [ ] **EVAL-02**: System can aggregate multi-AI evaluation results
- [ ] **EVAL-03**: System can detect AI evaluation discrepancies (anomaly detection)
- [ ] **EVAL-04**: Users can view detailed evaluation report
- [ ] **EVAL-05**: System can estimate tokens needed to realize idea

**Resource Matching**

- [ ] **MATCH-01**: System can match collaborators with complementary skills
- [ ] **MATCH-02**: System can match potential investors
- [ ] **MATCH-03**: System can match like-minded users (social)
- [ ] **MATCH-04**: System can match suitable job positions
- [ ] **MATCH-05**: System can match buyers and sellers
- [ ] **MATCH-06**: Users can set matching preferences

**Idea Crowdfunding**

- [ ] **IDEA-01**: Users can create ideas
- [ ] **IDEA-02**: Users can support ideas (Tokens)
- [ ] **IDEA-03**: Users can contribute code/design to ideas
- [ ] **IDEA-04**: System can track contribution weights
- [ ] **IDEA-05**: Idea creators can allocate revenue shares

**Task Marketplace**

- [ ] **TASK-01**: Users can publish tasks
- [ ] **TASK-02**: Users can apply for tasks
- [ ] **TASK-03**: Task publishers can assign tasks
- [ ] **TASK-04**: Users can submit milestones
- [ ] **TASK-05**: Task publishers can accept milestones
- [ ] **TASK-06**: System can auto-pay on completion
- [ ] **TASK-07**: Users can rate task completers

---

## Consensus Chain (Infrastructure - Foundation)

*Note: Consensus chain is underlying infrastructure supporting all core modules*

- [ ] **CONS-01**: System can process transfer transactions
- [ ] **CONS-02**: System can process staking transactions
- [ ] **CONS-03**: System can verify transaction signatures
- [ ] **CONS-04**: System can prevent double-spending
- [ ] **CONS-05**: System can query transaction status

---

## v2 Requirements (Future)

### Advanced Identity Verification

- **ID-A1**: OAuth login (Google, GitHub)
- **ID-A2**: Two-factor authentication
- **ID-A3**: Enterprise identity verification

### Advanced Service Marketplace Features

- **SVC-A1**: Batch requests
- **SVC-A2**: Streaming responses
- **SVC-A3**: Model fine-tuning configuration
- **SVC-A4**: Auto-retry and failover
- **SVC-A5**: Service quality monitoring

### Advanced Reputation Features

- **MQ-A1**: Reputation history charts
- **MQ-A2**: Reputation prediction model
- **MQ-A3**: Reputation staking (lock for higher privileges)

### Cross-chain Features

- **CHAIN-01**: IBC cross-chain asset transfer
- **CHAIN-02**: Cross-chain identity verification
- **CHAIN-03**: Cross-chain transactions

---

## Out of Scope

| Feature | Reason |
|---------|--------|
| Mobile App | Desktop web first, mobile UX is complex |
| Centralized Exchange | Not needed, fiat on-ramp violates decentralization principles |
| NFT Marketplace | Phase 1 focuses on utility transactions, not collectibles |
| Enterprise Deployment | Start with light nodes, validate P2P model first |
| Layer 2 Scaling | Validate if Layer 1 is sufficient first |
| DAO Governance | Build core features first, governance needs community |

---

## Traceability

### Core Modules

| Requirement | Module | Type | Status |
|-------------|--------|------|--------|
| **P2P-01~05** | CometBFT (built-in) | Core | N/A - Using existing solution |
| **ID-01~05** | x/identity (Cosmos SDK) | Core | Pending |
| **WALL-01~05** | Cosmos SDK Auth + Keplr | Core | N/A - Using existing solution |
| **SVC-01~13** | x/service (Cosmos SDK) | Core | Pending |
| **LLM-01~08** | x/service/llm (Cosmos SDK) | Core | Pending |
| **AGENT-01~08** | x/service/agent (Cosmos SDK) | Core | Pending |
| **WF-01~10** | x/service/workflow (Cosmos SDK) | Core | Pending |
| **PRICE-01~04** | Oracle Service (off-chain) | Core | Pending |
| **ESC-01~06** | x/escrow (Cosmos SDK) | Core | Pending |
| **MQ-01~06** | x/mq (Cosmos SDK) | Core | Pending |
| **DISP-01~08** | x/dispute (Cosmos SDK) | Core | Pending |
| **CONS-01~05** | CometBFT (built-in) | Core | N/A - Using existing solution |

### Optional Modules - Service Provider Plugins

| Requirement | Module | Type | Status |
|-------------|--------|------|--------|
| **LLM-PLUGIN-01~03** | LLM Provider Plugin | Plugin | Pending |
| **AGENT-PLUGIN-01~10** | OpenFang Executor Plugin | Plugin | Pending |
| **WF-PLUGIN-01~06** | Workflow Executor Plugin | Plugin | Pending |

### Optional Modules - User Plugins

| Requirement | Module | Type | Status |
|-------------|--------|------|--------|
| **LAMP-01~05** | AI Chat UI Plugin (GenieBot) | Plugin | Pending |
| **IDEA-PROC-01~05** | AI Chat UI Plugin (GenieBot) | Plugin | Pending |
| **EVAL-01~05** | AI Service Plugin (off-chain) | Plugin | Pending |
| **MATCH-01~06** | Resource Matching Plugin (off-chain) | Plugin | Pending |
| **IDEA-01~05** | x/idea (Cosmos SDK) | Plugin | Pending |
| **TASK-01~07** | x/task (Cosmos SDK) | Plugin | Pending |

### Legacy Requirement Mapping

| Old ID | New ID | Change Notes |
|--------|--------|--------------|
| COMP-01~10 | LLM-01~08, SVC-01~13 | Compute trading renamed to LLM service layer, core routing logic moved to service marketplace |
| WF-01~07 | WF-01~10 | Workflow upgraded to service type, added service registration and settlement requirements |
| GH-01~05 | WF-PLUGIN-01~06 | GitHub integration became Workflow executor plugin |
| OPENFANG-01~10 | AGENT-PLUGIN-01~10 | OpenFang became Agent executor plugin |

**Coverage:**
- v1 requirements: 118 total
  - Core: 76 requirements (64%)
  - Plugin: 42 requirements (36%)
- Mapped to modules: 118
- Unmapped: 0 ✓

**Module Distribution:**
- Core Modules: 76 requirements
  - P2P Communication: 5
  - Identity: 5
  - Wallet: 5
  - Service Marketplace: 47 (Service Management 13 + LLM Layer 8 + Agent Layer 8 + Workflow Layer 10 + Pricing 4)
  - Escrow: 6
  - Trust System: 14 (MQ 6 + Dispute 8)
  - Consensus Chain: 5
- Plugin Modules: 42 requirements
  - LLM Provider Plugin: 3
  - Agent Executor Plugin: 10
  - Workflow Executor Plugin: 6
  - GenieBot User Plugin: 23 (Chat 5 + Idea Processing 5 + Evaluation 5 + Matching 6 + Crowdfunding 5 + Tasks 7)

---

## Business Rules

### Trust System (MQ)
- **Initial Value**: Each new user gets 100 MQ
- **Constant Total**: Total MQ = User Count x 100
- **Zero-Sum Game**: Sum of MQ changes for dispute parties = 0
- **Decay Mechanism**: 0.1% daily decay, minimum value 10
- **Tier Classification** (consistent with x/mq module implementation):
  - Lv1 Newcomer: 0-50 MQ (Restricted User)
    - Basic transaction permissions
    - Cannot serve as jury member
    - Limited service provision
  - Lv2 Member: 50-100 MQ (Regular User)
    - Standard transaction permissions
    - Can participate in small dispute reviews
    - Normal service provision permissions
  - Lv3 Trusted: 100-200 MQ (Active User)
    - Full transaction permissions
    - Can participate in medium dispute reviews
    - Priority service matching
  - Lv4 Expert: 200-500 MQ (Excellent User)
    - Advanced transaction permissions
    - Can participate in large dispute reviews
    - Priority display in service marketplace
    - Voting weight bonus
  - Lv5 Guardian: 500+ MQ (Outstanding User)
    - Highest transaction permissions
    - Can serve as jury chair
    - Highest voting weight
    - Can participate in governance proposals

### Service Marketplace
- **Service Types**: LLM / Agent / Workflow three tiers
- **Pricing Modes**: Per-use / Volume-based / Subscription
- **Routing Strategy**: MQ priority / Price priority / Capability priority / Load balancing
- **Service Levels**:
  - Basic Service: MQ >= 50
  - Advanced Service: MQ >= 100
  - Professional Service: MQ >= 200

### Jury Selection
- **Size Rules**:
  - Small dispute (<=100 STT): 3 people
  - Medium dispute (100-1000 STT): 5 people
  - Large dispute (>1000 STT): 7 people
- **Qualification Requirements**: MQ > 50, Activity > 30 days
- **Recusal Rules**: Auto-recuse if transaction history with either party

### Escrow Release Conditions
- **Auto Release**: 24 hours after task completion confirmation
- **Partial Release**: Requires both parties' signatures
- **Dispute Lock**: Funds frozen during dispute
- **Timeout Handling**: Auto-release to provider after 7 days no response

### Dispute Processing Time Limits
- **Negotiation Period**: 7 days
- **Voting Period**: 3 days
- **Appeal Period**: 7 days
- **Pass Threshold**: 60% weighted votes

### Revenue Distribution
- **Contribution Weights**: Code (30%) + Design (50%) + Bug Fixes (20%)
- **Platform Fee**: 2%
- **Distribution Cycle**: Auto-distribute when project generates revenue

### Task Milestones
- **Revision Limit**: Max 3 revisions
- **Review Timeout**: Auto-pass after 3 days no response
- **Acceptance Criteria**: Defined by task publisher via checklist

---

*Requirements defined: 2025-03-01*
*Last updated: 2026-03-02 - Restructured to service marketplace architecture, separated core modules and plugin modules*
