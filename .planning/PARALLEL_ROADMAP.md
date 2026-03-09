# ShareTokens Roadmap

> **This is the single source of truth for ShareTokens development roadmap.**
>
> Previous ROADMAP.md has been archived. This document (PARALLEL_ROADMAP.md v3.0) replaces it entirely.

**Created:** 2026-03-01
**Updated:** 2026-03-02
**Approach:** Cosmos SDK AppChain + TypeScript Frontend
**Total Phases:** 4 (Phase 0-3) + Integration + E2E

---

## Architecture: Core + Plugins

```
+---------------------------------------------------------------+
|                    CORE MODULES (Required)                     |
|  Every node must have these modules for basic operation       |
|                                                                |
|  P2P | Identity | Wallet | Service Market | Escrow | MQ | Dispute |
|                           (LLM/Agent/Workflow)                 |
+---------------------------------------------------------------+
                               |
                               v
+---------------------------------------------------------------+
|                    OPTIONAL PLUGINS                            |
|  Install based on node role                                    |
|                                                                |
|  Provider Plugins:     | User Plugins:                        |
|  - LLM API Key Host    | - GenieBot Interface                  |
|  - Agent Executor      |                                       |
|  - Workflow Executor   |                                       |
+---------------------------------------------------------------+
```

### Core Principle: Don't Reinvent the Wheel

| Feature | Open Source Solution | Notes |
|---------|---------------------|-------|
| P2P Network | **CometBFT Built-in** | No need for libp2p |
| Consensus (BFT) | **Tendermint** | No need to implement consensus |
| Account/Signature | **Cosmos SDK Auth** | No need to develop Wallet module |
| Transfer | **Cosmos SDK Bank** | No need to implement |
| Staking | **Cosmos SDK Staking** | No need to implement |
| Price Oracle | **Chainlink Network** | Read existing data, no node deployment |
| IBC Cross-chain | **Cosmos IBC** | v2 feature, use directly |

### What We Need to Develop

1. **Cosmos SDK Business Modules (Go)** - Identity, MQ, Service Market, Escrow, Dispute
2. **Provider Plugins (TypeScript/Rust)** - LLM Key Host, Agent Executor, Workflow Executor
3. **User Plugins (TypeScript)** - GenieBot Interface

---

## Architecture Overview

```
+-------------------------------------------------------------------------+
|                     User Plugins (TypeScript/React)                      |
|  +-----------------------------------------------------------------+    |
|  |  GenieBot Interface - AI Chat, Service Discovery, Result Display |    |
|  +-----------------------------------------------------------------+    |
|                            | REST/gRPC                                  |
+----------------------------+--------------------------------------------+
|                     Provider Plugins (TypeScript/Rust)                   |
|  +----------------+  +----------------+  +----------------+              |
|  | LLM Key Host   |  | Agent Executor |  | Workflow Exec  |              |
|  | (TypeScript)   |  | (OpenFang/Rust)|  | (TypeScript)   |              |
|  +----------------+  +----------------+  +----------------+              |
|                            | ABCI / RPC                                 |
+----------------------------+--------------------------------------------+
|                 CORE MODULES - ShareTokens Chain                         |
|                     (Cosmos SDK + CometBFT)                              |
|  +-----------------------------------------------------------------+    |
|  |  Core Business Modules (Go - What we develop)                    |    |
|  |  +-------------+  +-------------+  +-------------+               |    |
|  |  | x/identity  |  | x/wallet    |  | x/market    |               |    |
|  |  | (ID-01~05)  |  | (WALL-01~05)|  | (SVC,LLM,   |               |    |
|  |  |             |  |             |  |  AGENT,WF)  |               |    |
|  |  +-------------+  +-------------+  +-------------+               |    |
|  |  +-------------+  +-------------+  +-------------+               |    |
|  |  | x/escrow    |  | x/mq        |  | x/dispute   |               |    |
|  |  | (ESC-01~06) |  | (MQ-01~06)  |  | (DISP-01~08)|               |    |
|  |  +-------------+  +-------------+  +-------------+               |    |
|  +-----------------------------------------------------------------+    |
|  +-----------------------------------------------------------------+    |
|  |  Cosmos SDK Base Modules (Use directly)                          |    |
|  |  +-------------------------------------------------------------+  |    |
|  |  | auth (Account) | bank (Transfer) | staking | params         |  |    |
|  |  +-------------------------------------------------------------+  |    |
|  +-----------------------------------------------------------------+    |
|  +-----------------------------------------------------------------+    |
|  |  CometBFT (Use directly)                                         |    |
|  |  +-------------------------------------------------------------+  |    |
|  |  | P2P Network | BFT Consensus | ABCI                          |  |    |
|  |  +-------------------------------------------------------------+  |    |
|  +-----------------------------------------------------------------+    |
+-------------------------------------------------------------------------+
```

---

## Development Priority

```
Phase 0: Project Init & Interface Definition (1 week)
           |
           v
Phase 1: Core Modules (4-5 weeks) <-- PRIORITY
           |
           +-- 1a: Foundation (Identity, Wallet, Service Market Core)
           |
           +-- 1b: Trust (Escrow, MQ, Dispute)
           |
           v
Phase 2: Provider Plugins (2-3 weeks)
           |
           +-- LLM Key Host Plugin
           +-- Agent Executor Plugin (OpenFang)
           +-- Workflow Executor Plugin
           |
           v
Phase 3: User Plugins (2-3 weeks)
           |
           +-- GenieBot Interface
           |
           v
Integration & E2E (3-4 weeks)
```

---

## Phase 0: Project Init & Interface Definition

**Duration:** 1 week
**Goal:** Set up Cosmos SDK chain skeleton, define all module interfaces

### 0.1 Cosmos SDK Chain Init

```bash
# Use Ignite CLI to create chain skeleton
ignite scaffold chain sharetokens --no-module

# Directory structure
chain/
+-- app/
|   +-- app.go              # Application config
|   +-- encoding.go         # Encoding config
+-- cmd/
|   +-- sharetokensd/       # Chain binary
+-- x/                      # Custom modules directory
+-- proto/                  # Protobuf definitions
+-- testnets/               # Test network config
```

### 0.2 Interface Definition (Protobuf)

```
proto/
+-- sharetokens/
    +-- identity/
    |   +-- v1/
    |       +-- genesis.proto
    |       +-- query.proto
    |       +-- tx.proto
    +-- market/
    |   +-- v1/
    |       +-- service.proto    # Service registration
    |       +-- llm.proto        # LLM service layer
    |       +-- agent.proto      # Agent service layer
    |       +-- workflow.proto   # Workflow service layer
    +-- escrow/
    |   +-- v1/...
    +-- mq/
    |   +-- v1/...
    +-- dispute/
        +-- v1/...
```

### 0.3 TypeScript Client Interfaces

```typescript
// src/interfaces/
interfaces/
+-- chain/
|   +-- IIdentityClient.ts
|   +-- IMarketClient.ts
|   +-- IEscrowClient.ts
|   +-- IMQClient.ts
|   +-- IDisputeClient.ts
+-- plugin/
|   +-- ILLMProviderPlugin.ts
|   +-- IAgentExecutorPlugin.ts
|   +-- IWorkflowExecutorPlugin.ts
|   +-- IGenieBotPlugin.ts
+-- mocks/
    +-- MockChainClient.ts
    +-- MockPluginServices.ts
```

### Success Criteria

- [ ] `sharetokensd` can start single-node dev network
- [ ] All Protobuf definitions complete
- [ ] TypeScript interface definitions complete
- [ ] Mock implementations can return test data

---

## Phase 1: Core Modules (Go - Cosmos SDK)

**Duration:** 4-5 weeks (Parallel)
**Goal:** Develop Cosmos SDK business modules - THE FOUNDATION

### Module Dependency Graph

```
                    +-------------+
                    | auth        |  (Cosmos SDK Built-in)
                    | bank        |
                    +------+------+
                           |
         +-----------------+-----------------+
         |                 |                 |
         v                 v                 v
    +---------+      +------------+    +---------+
    | identity|      |   market   |    | escrow  |  <-- Phase 1a (Parallel)
    +---------+      +------------+    +---------+
         |                 |                 |
         +-----------------+-----------------+
                           |
         +-----------------+-----------------+
         |                 |                 |
         v                 v                 v
    +---------+      +------------+    +---------+
    |   mq    |      |  dispute   |    | (task)  |  <-- Phase 1b (Parallel)
    +---------+      +------------+    +---------+
```

### Phase 1a: Foundation Modules (Week 1-2)

---

### Module 1.1: Identity (x/identity)

| Attribute | Value |
|-----------|-------|
| **Requirements** | ID-01 to ID-05 |
| **Dependencies** | auth, bank (Cosmos SDK) |
| **Language** | Go (Cosmos SDK) |
| **Parallel safe** | YES (Use MockKeeper) |

**Development Tasks:**
- [ ] Use Ignite scaffold to create module skeleton
- [ ] Define State: Identity, IdentityRegistry
- [ ] Implement Keeper: RegisterIdentity, VerifyIdentity, RevokeIdentity
- [ ] Implement Msg Server: Transaction handling
- [ ] Implement Query Server: Query interface
- [ ] Implement Genesis import/export

**Ignite Commands:**
```bash
ignite scaffold module identity --dep auth,bank
ignite scaffold map identity owner hash:type=string level:type=string --module identity
```

---

### Module 1.2: Service Market (x/market) - CORE BUSINESS

| Attribute | Value |
|-----------|-------|
| **Requirements** | SVC-01~13, LLM-01~08, AGENT-01~08, WF-01~10, PRICE-01~04 |
| **Dependencies** | auth, bank |
| **Language** | Go (Cosmos SDK) |
| **Parallel safe** | YES |
| **Priority** | HIGHEST - Core Business Logic |

**Sub-modules:**

```
x/market/
+-- keeper/
|   +-- keeper.go           # Main keeper
|   +-- service.go          # Service registration/discovery
|   +-- pricing.go          # Pricing strategy
|   +-- routing.go          # Request routing
|   +-- llm.go              # LLM service layer
|   +-- agent.go            # Agent service layer
|   +-- workflow.go         # Workflow service layer
+-- types/
    +-- service.go          # Service types
    +-- llm.go              # LLM types
    +-- agent.go            # Agent types
    +-- workflow.go         # Workflow types
```

**Development Tasks:**

**Service Registration & Discovery:**
- [ ] Define State: Service, ServiceType, ServiceStatus
- [ ] Implement RegisterService (LLM/Agent/Workflow)
- [ ] Implement UpdateService / DecommissionService
- [ ] Implement QueryServices (filter by type, capability, price)
- [ ] Implement GetServiceDetails

**Pricing:**
- [ ] Define State: PricingStrategy, PriceHistory
- [ ] Implement SetPricingStrategy (per-use/per-quantity/subscription)
- [ ] Implement ConvertToSTT (based on oracle rate)
- [ ] Implement DynamicPricing (based on supply/demand)

**Request Routing:**
- [ ] Implement RouteRequest (find best provider)
- [ ] Implement MatchByMQ (reputation-based)
- [ ] Implement MatchByPrice (price-based)
- [ ] Implement MatchByCapability (ability-based)
- [ ] Implement LoadBalancing
- [ ] Implement RetryAndFailover

**LLM Service Layer:**
- [ ] Define State: LLMProvider, LLMRequest, LLMResponse
- [ ] Implement RegisterLLMProvider
- [ ] Implement SubmitLLMRequest
- [ ] Implement RecordLLMResponse
- [ ] Implement VerifyTokenUsage
- [ ] Implement AutoSettlement

**Agent Service Layer:**
- [ ] Define State: AgentProvider, AgentRequest, AgentResponse
- [ ] Implement RegisterAgentProvider
- [ ] Implement SubmitAgentRequest
- [ ] Implement RecordAgentResponse
- [ ] Implement VerifyAgentResult
- [ ] Implement AutoSettlement

**Workflow Service Layer:**
- [ ] Define State: WorkflowProvider, WorkflowRequest, WorkflowResponse
- [ ] Implement RegisterWorkflowProvider
- [ ] Implement IdentifyWorkflowType
- [ ] Implement SubmitWorkflowRequest
- [ ] Implement RecordWorkflowProgress
- [ ] Implement VerifyWorkflowResult
- [ ] Implement AutoSettlement

---

### Module 1.3: Escrow (x/escrow)

| Attribute | Value |
|-----------|-------|
| **Requirements** | ESC-01 to ESC-06 |
| **Dependencies** | auth, bank |
| **Language** | Go (Cosmos SDK) |
| **Parallel safe** | YES |

**Development Tasks:**
- [ ] Define State: Escrow, EscrowStatus
- [ ] Implement CreateEscrow (lock funds)
- [ ] Implement ReleaseEscrow (release to payee)
- [ ] Implement PartialRelease (partial release)
- [ ] Implement LockForDispute (dispute lock)
- [ ] Implement ResolveByRuling (ruling distribution)

---

### Phase 1b: Trust Modules (Week 3-4)

---

### Module 1.4: MQ - Moral Quotient (x/mq)

| Attribute | Value |
|-----------|-------|
| **Requirements** | MQ-01 to MQ-06 |
| **Dependencies** | auth, identity |
| **Language** | Go (Cosmos SDK) |
| **Parallel safe** | YES (Use MockIdentityKeeper) |

**Development Tasks:**
- [ ] Define State: MoralQuotient, MQConfig
- [ ] Implement MQ initialization (new user 100 MQ)
- [ ] Implement RedistributeMQ (zero-sum game)
- [ ] Implement ApplyDecay (decay mechanism)
- [ ] Implement SelectJury (based on MQ)
- [ ] Implement GetVoteWeight (vote weight calculation)
- [ ] Implement GetMQLevel (level calculation)

**MQ Levels:**
| Level | Name | MQ Range | Permissions |
|-------|------|----------|-------------|
| Lv1 | Newcomer | 0-50 | Basic trading, no jury |
| Lv2 | Member | 50-100 | Standard trading, small dispute jury |
| Lv3 | Trusted | 100-200 | Full trading, medium dispute jury |
| Lv4 | Expert | 200-500 | Advanced trading, large dispute jury |
| Lv5 | Guardian | 500+ | Highest trading, jury chair, governance |

---

### Module 1.5: Dispute (x/dispute)

| Attribute | Value |
|-----------|-------|
| **Requirements** | DISP-01 to DISP-08 |
| **Dependencies** | auth, escrow, mq |
| **Language** | Go (Cosmos SDK) |
| **Parallel safe** | YES (Use MockKeeper) |

**Development Tasks:**
- [ ] Define State: Dispute, Evidence, Vote, Resolution
- [ ] Implement CreateDispute
- [ ] Implement SubmitEvidence
- [ ] Implement AssembleJury (call MQ module)
- [ ] Implement CastVote
- [ ] Implement ResolveDispute (trigger MQ redistribution)
- [ ] Implement AppealDispute

---

### Phase 1 Parallel Execution Plan

```
Week 1-2 (Phase 1a - 3 modules parallel):
+----------------+  +----------------+  +----------------+
|    Identity    |  | Service Market |  |     Escrow     |
|     (Go)       |  |      (Go)      |  |      (Go)      |
+----------------+  +----------------+  +----------------+

Week 3-4 (Phase 1b - 2 modules parallel):
+----------------+  +----------------+
|       MQ       |  |    Dispute     |
|      (Go)      |  |      (Go)      |
+----------------+  +----------------+

Week 5: Integration & Testing
+------------------------------------------------+
|          Core Module Integration                |
|  - Keeper dependency injection                  |
|  - Cross-module testing                         |
|  - Genesis state testing                        |
+------------------------------------------------+
```

---

## Phase 2: Provider Plugins (TypeScript/Rust)

**Duration:** 2-3 weeks (Parallel)
**Goal:** Develop plugins for service providers

### Plugin Architecture

```
+-----------------------------------------------------------------+
|                    Provider Plugin System                         |
|                                                                   |
|  +-------------+    +-------------+    +-------------+            |
|  | LLM Key     |    | Agent       |    | Workflow    |            |
|  | Host        |    | Executor    |    | Executor    |            |
|  | (TypeScript)|    | (OpenFang)  |    | (TypeScript)|            |
|  +------+------+    +------+------+    +------+------+            |
|         |                  |                  |                   |
|         v                  v                  v                   |
|  +--------------------------------------------------------+       |
|  |              Plugin SDK (Common Interface)              |       |
|  |  - Service Registration                                 |       |
|  |  - Health Reporting                                     |       |
|  |  - Payment Settlement                                   |       |
|  |  - Error Handling                                       |       |
|  +--------------------------------------------------------+       |
+-----------------------------------------------------------------+
                            |
                            v
                  +-------------------+
                  |  ShareTokens Chain |
                  |  (Core Modules)    |
                  +-------------------+
```

---

### Plugin 2.1: LLM API Key Host

| Attribute | Value |
|-----------|-------|
| **Requirements** | LLM-PLUGIN-01~03 |
| **Dependencies** | x/market (Chain) |
| **Language** | TypeScript |
| **Parallel safe** | YES |

**Development Tasks:**
- [ ] Secure API Key storage (encrypted)
- [ ] Multi-provider support (OpenAI, Anthropic, etc.)
- [ ] Usage monitoring
- [ ] Request proxying
- [ ] Response caching
- [ ] Auto-pricing based on cost

**Directory Structure:**
```
plugins/llm-provider/
+-- src/
|   +-- index.ts           # Plugin entry
|   +-- keyManager.ts      # API key management
|   +-- proxy.ts           # Request proxy
|   +-- pricing.ts         # Auto pricing
|   +-- monitor.ts         # Usage monitoring
+-- config/
    +-- providers.yaml     # Provider configurations
```

---

### Plugin 2.2: Agent Executor (OpenFang Integration)

| Attribute | Value |
|-----------|-------|
| **Requirements** | AGENT-PLUGIN-01~10 |
| **Dependencies** | x/market (Chain), OpenFang |
| **Language** | TypeScript + Rust (OpenFang) |
| **Parallel safe** | YES |

**Development Tasks:**
- [ ] OpenFang installation and configuration
- [ ] Agent type implementations:
  - [ ] Coder Agent (software development)
  - [ ] Researcher Agent (research and analysis)
  - [ ] Writer Agent (content creation)
  - [ ] Collector Hand (data collection)
  - [ ] Content Hand (content workflows)
  - [ ] Lead Hand (resource matching)
- [ ] Dashboard integration
- [ ] Security sandbox (16-layer security)
- [ ] Custom agent template support

**Directory Structure:**
```
plugins/agent-executor/
+-- src/
|   +-- index.ts           # Plugin entry
|   +-- openfang.ts        # OpenFang bridge
|   +-- agents/
|   |   +-- coder.ts
|   |   +-- researcher.ts
|   |   +-- writer.ts
|   +-- hands/
|   |   +-- collector.ts
|   |   +-- content.ts
|   |   +-- lead.ts
|   +-- sandbox.ts         # Security sandbox
+-- openfang/              # OpenFang config
```

---

### Plugin 2.3: Workflow Executor

| Attribute | Value |
|-----------|-------|
| **Requirements** | WF-PLUGIN-01~06 |
| **Dependencies** | x/market (Chain), Agent Executor |
| **Language** | TypeScript |
| **Parallel safe** | YES |

**Development Tasks:**
- [ ] GitHub integration
- [ ] Workflow type implementations:
  - [ ] Software Development Workflow
  - [ ] Content Creation Workflow
  - [ ] Business Planning Workflow
  - [ ] Life Services Workflow
- [ ] Human review node support
- [ ] Progress tracking
- [ ] Result verification

**Workflow Definitions:**

| Type | Steps |
|------|-------|
| Software Dev | Requirements -> Architecture -> Code -> Test -> Deploy |
| Content | Topic -> Research -> Create -> Review -> Publish |
| Business | Market Analysis -> Plan Design -> Review -> Track |
| Life | Request -> Match -> Execute -> Verify -> Complete |

**Directory Structure:**
```
plugins/workflow-executor/
+-- src/
|   +-- index.ts           # Plugin entry
|   +-- github.ts          # GitHub integration
|   +-- workflows/
|   |   +-- software.ts
|   |   +-- content.ts
|   |   +-- business.ts
|   |   +-- life.ts
|   +-- review.ts          # Human review
|   +-- tracker.ts         # Progress tracking
```

---

### Phase 2 Parallel Execution Plan

```
Week 1-2 (Phase 2 - 3 plugins parallel):
+----------------+  +----------------+  +----------------+
|   LLM Host     |  | Agent Executor |  | Workflow Exec  |
|  (TypeScript)  |  | (OpenFang/TS)  |  |  (TypeScript)  |
+----------------+  +----------------+  +----------------+

Week 3: Plugin Integration
+------------------------------------------------+
|          Plugin System Testing                  |
|  - Plugin-to-chain communication                |
|  - Service registration testing                 |
|  - Settlement testing                           |
+------------------------------------------------+
```

---

## Phase 3: User Plugins (TypeScript/React)

**Duration:** 2-3 weeks
**Goal:** Develop user-facing plugins

---

### Plugin 3.1: GenieBot Interface

| Attribute | Value |
|-----------|-------|
| **Requirements** | LAMP-01~05, IDEA-PROC-01~05, EVAL-01~05, MATCH-01~06 |
| **Dependencies** | Core Modules, Provider Plugins |
| **Language** | TypeScript + React |
| **Parallel safe** | YES |

**Development Tasks:**

**Core UI Components:**
- [ ] Wallet connection (Keplr integration)
- [ ] Navigation layout
- [ ] Theme system

**AI Chat UI:**
- [ ] Conversation interface
- [ ] Message rendering
- [ ] Idea cards
- [ ] Progress display

**Service Market UI:**
- [ ] Service discovery
- [ ] Service details
- [ ] Pricing display
- [ ] Request submission

**Idea Processing:**
- [ ] Idea classification (software/content/business/life)
- [ ] Idea refinement through conversation
- [ ] Multi-AI evaluation
- [ ] Token estimation

**Resource Matching:**
- [ ] Collaborator matching
- [ ] Investor matching
- [ ] Social matching
- [ ] Job matching
- [ ] Trading matching

**Idea Crowdfunding:**
- [ ] Create idea
- [ ] Support idea (Token)
- [ ] Contribute code/design
- [ ] Revenue distribution

**Task Market:**
- [ ] Post task
- [ ] Apply for task
- [ ] Milestone submission
- [ ] Review and payment

**Directory Structure:**
```
plugins/xiaodeng/
+-- src/
|   +-- App.tsx
|   +-- components/
|   |   +-- chat/
|   |   +-- market/
|   |   +-- idea/
|   |   +-- task/
|   +-- hooks/
|   +-- services/
|   |   +-- chain.ts
|   |   +-- ai.ts
|   +-- store/
+-- public/
```

---

## Integration Phase

**Duration:** 2 weeks
**Goal:** Integrate all modules

### Integration Order

```
1. Chain Integration
   Cosmos SDK modules connection (Keeper dependency)

2. Chain-Plugin Integration
   Provider plugins connect to chain (RPC/WebSocket)

3. User Plugin Integration
   GenieBot connects to chain and services

4. End-to-End Testing
```

---

## E2E Phase

**Duration:** 1-2 weeks
**Goal:** Complete flow verification

### Test Scenarios

#### Scenario 1: LLM Service Trading
```
1. User A registers identity
2. User A installs LLM Provider Plugin, hosts API Key
3. Service Market registers LLM service
4. User B submits LLM request
5. System routes to User A's service
6. User A executes request
7. System verifies and settles payment
```

#### Scenario 2: Agent Service Trading
```
1. Provider installs Agent Executor Plugin
2. Registers Agent service (coder/researcher/writer)
3. User submits Agent task request
4. System routes to appropriate Agent
5. Agent executes and returns result
6. System verifies and settles payment
```

#### Scenario 3: Workflow Service Trading
```
1. Provider installs Workflow Executor Plugin
2. Registers Workflow service
3. User submits idea via GenieBot
4. System identifies Workflow type
5. Workflow executes multi-step process
6. Human review at key nodes
7. System verifies and settles payment
```

#### Scenario 4: Dispute Resolution
```
1. Dispute occurs
2. System locks Escrow
3. MQ module selects Jury
4. Jury votes
5. Ruling executed
6. MQ redistributed
```

---

## Module Summary

### Core Modules (Go - Cosmos SDK)

| Module | Requirements | Dependencies | Complexity | Parallel |
|--------|--------------|--------------|------------|----------|
| x/identity | ID-01~05 | auth, bank | M | YES |
| x/market | SVC,LLM,AGENT,WF,PRICE | auth, bank | L | YES |
| x/escrow | ESC-01~06 | auth, bank | M | YES |
| x/mq | MQ-01~06 | auth, identity | M | YES |
| x/dispute | DISP-01~08 | auth, escrow, mq | M | YES |

### Provider Plugins (TypeScript/Rust)

| Plugin | Requirements | Dependencies | Complexity | Parallel |
|--------|--------------|--------------|------------|----------|
| LLM Key Host | LLM-PLUGIN-01~03 | x/market | M | YES |
| Agent Executor | AGENT-PLUGIN-01~10 | x/market, OpenFang | L | YES |
| Workflow Executor | WF-PLUGIN-01~06 | x/market, GitHub | M | YES |

### User Plugins (React + TypeScript)

| Plugin | Requirements | Complexity | Parallel |
|--------|--------------|------------|----------|
| GenieBot Interface | LAMP,IDEA-PROC,EVAL,MATCH,IDEA,TASK | L | YES |

---

## Timeline Summary

| Phase | Duration | Cumulative | Description |
|-------|----------|------------|-------------|
| Phase 0: Init | 1 week | 1 week | Project setup & interfaces |
| Phase 1: Core Modules | 4-5 weeks | 5-6 weeks | Chain business modules |
| Phase 2: Provider Plugins | 2-3 weeks | 7-9 weeks | Service provider plugins |
| Phase 3: User Plugins | 2-3 weeks | 9-12 weeks | User interface plugins |
| Integration | 2 weeks | 11-14 weeks | All modules integration |
| E2E Testing | 1-2 weeks | 12-16 weeks | Complete flow testing |

**Total:** 12-16 weeks (3-4 months)

---

## Key Changes from v2.0

| Aspect | v2.0 | v3.0 (This Version) |
|--------|------|---------------------|
| Architecture | All modules equal | Core + Plugins separation |
| Service Market | compute module | Dedicated x/market with 3 layers |
| Plugin System | Not defined | Clear plugin architecture |
| Development Priority | Mixed | Core -> Provider Plugins -> User Plugins |
| OpenFang | Separate module | Agent Executor Plugin |
| Workflow | Chain module | Workflow Executor Plugin |
| GenieBot | Frontend | User Plugin (clear boundary) |

---

## Comparison: Old Plan vs New Plan

| Aspect | Old Plan | New Plan |
|--------|----------|----------|
| P2P Network | Develop with libp2p | Use CometBFT built-in |
| Consensus | Develop | Use Tendermint |
| Wallet | Develop module | Use Cosmos SDK Auth + Keplr |
| Architecture | Monolithic | Core + Plugins |
| Service Market | compute module | Dedicated market with 3 layers |
| Development Language | TypeScript + Go mixed | Chain Go, Plugins TypeScript |
| Development Time | 14-20 weeks | 12-16 weeks |
| Reinventing Wheel | Much | Little |

---

*Roadmap Version: 3.0*
*Updated: 2026-03-02 - Refactored to Core + Plugins architecture, Service Market as core business*
