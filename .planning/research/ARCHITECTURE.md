# Architecture Research

**Domain:** ShareTokens - Decentralized AI Service Marketplace
**Researched:** 2026-03-02
**Confidence:** MEDIUM-HIGH

## Architecture Overview

ShareTokens follows a **Core + Plugin** architecture pattern. The system is divided into mandatory core modules that every node must have, and optional plugins that nodes can install based on their role.

### Core + Plugin Architecture

```
+-----------------------------------------------------------------------------+
|                        Core Modules (Every Node Must Have)                   |
+-----------------------------------------------------------------------------+
|  +--------+  +--------+  +--------+  +-------------+  +--------+  +------------+ |
|  |  P2P   |  |Identity|  | Wallet |  |   Service   |  | Escrow |  |  Trust     | |
|  |  Comm  |  | Account|  |        |  |   Market    |  |  Pay   |  |  System    | |
|  +---+----+  +---+----+  +---+----+  +------+------+  +---+----+  +---+--------+ |
|      |           |           |              |             |              |        |
|      +-----------+-----------+--------------+-------------+--------------+        |
|                              |                |                                  |
|                              |    (MQ + Dispute Arbitration)            |
|                              +--------------------------------------------------+        |
+-----------------------------------------------------------------------------+
                               |
                               v
+-----------------------------------------------------------------------------+
|                        Optional Plugins (Install By Role)                    |
+-----------------------------------------------------------------------------+
|  Service Provider Plugins:                                                   |
|  +----------------+  +----------------+  +----------------+                  |
|  | LLM API Key    |  | Agent Executor |  | Workflow       |                  |
|  | Hosting        |  | (OpenFang)     |  | Executor       |                  |
|  +----------------+  +----------------+  +----------------+                  |
|                                                                              |
|  User/Demand Side Plugins:                                                   |
|  +----------------+                                                          |
|  | GenieBot       |  AI conversation, idea collection, service invocation    |
|  | Interface      |  result display                                          |
|  +----------------+                                                          |
+-----------------------------------------------------------------------------+
```

## Core Modules (Layer 0)

> These modules are mandatory for every node. Missing any of these prevents the node from functioning.

### 1. P2P Communication

**Purpose:** Node discovery, message broadcasting, NAT traversal

**Implementation:**
- Built on CometBFT's built-in P2P layer
- PEX Reactor + Address Book for peer discovery
- CometBFT Reactor Gossip for broadcast messaging
- SecretConnection (TLS 1.3) for encrypted channels
- NAT traversal via UPnP + external_address configuration

**Protocol Stack:**
```
+------------------+
|   Application    |
+--------+---------+
         |
+--------+---------+
|   CometBFT P2P   |  <- Built-in P2P for consensus
+--------+---------+
         |
+--------+---------+
|  Transport Layer |  TCP/QUIC
+--------+---------+
         |
+--------+---------+
|   Network Layer  |  IP routing
+------------------+
```

### 2. Identity/Account

**Purpose:** Unique identification, real-name verification, anti-Sybil

**Implementation:**
- Cosmos SDK Auth module (built-in)
- Optional ZK-DID for privacy-preserving KYC
- Identity hash registry on-chain
- Real-name verification via attestation oracles

**Data Model:**
```go
type Identity struct {
    Address       cosmos.AccAddress
    VerifiedHash  []byte        // Hash of real identity (not stored raw)
    MQScore       int64         // Moral Quotient score
    CreatedAt     time.Time
    Status        IdentityStatus  // Active, Suspended, Banned
}
```

### 3. Wallet

**Purpose:** Balance queries, transaction signing, STT payments

**Implementation:**
- Cosmos SDK Bank module (built-in)
- Keplr Wallet integration (standard Cosmos wallet)
- WalletConnect for mobile
- STT token as native denom

**Key Operations:**
- View balance
- Sign transactions
- Send/receive STT
- View transaction history

### 4. Service Market (Core Business)

**Purpose:** Service registration, discovery, pricing, routing

**Three-Layer Service Structure:**

| Level | Service Type | Billing | Description |
|-------|-------------|---------|-------------|
| Level 1 | LLM API | Per-token | Direct API calls (GPT-4, Claude, etc.) |
| Level 2 | Agent | Per-skill | AI Agents (coder, researcher, writer, etc.) |
| Level 3 | Workflow | Per-workflow | Multi-agent workflows (dev pipeline, content creation) |

**Marketplace Functions:**
- Service registration and discovery
- Provider capability advertising
- Pricing strategies (fixed, dynamic, auction)
- Smart routing to optimal providers

**Data Model:**
```go
type Service struct {
    ID            string
    Provider      cosmos.AccAddress
    Level         ServiceLevel    // 1=LLM, 2=Agent, 3=Workflow
    Name          string
    Description   string
    Pricing       PricingModel
    Capabilities  []string
    MQRequired    int64          // Minimum MQ to use
    Status        ServiceStatus
}
```

### 5. Escrow Payment

**Purpose:** Lock STT before tasks, release on completion, freeze on disputes

**Implementation:**
- Custom Cosmos SDK module (x/escrow)
- Multi-sig holding accounts
- Time-locked releases
- Dispute triggers

**Escrow Flow:**
```
[Consumer] --lock STT--> [Escrow]
                           |
                    [Task Execution]
                           |
          +----------------+----------------+
          v                                 v
    [Success]                         [Dispute]
          |                                 |
    Release to Provider              Freeze for Arbitration
```

### 6. Trust System

**Purpose:** MQ scoring, level classification, history tracking, dispute arbitration

**Combines:** MQ + Dispute Arbitration

**Key Features:**
- Zero-sum MQ (winners gain = losers lose)
- Time-based decay
- MQ-weighted privileges
- Anti-gaming mechanisms

**Scoring Model:**
```go
type MQScore struct {
    Address       cosmos.AccAddress
    Score         int64         // Can be negative
    Level         MQLevel       // Based on score thresholds
    LastActive    time.Time
    History       []MQEvent
}

func (m *MQScore) ApplyDecay() {
    daysSinceActive := time.Since(m.LastActive).Hours() / 24
    decay := math.Pow(0.99, daysSinceActive)  // 1% daily decay
    m.Score = int64(float64(m.Score) * decay)
}
```

**Dispute Resolution (within Trust System):**
- Jury mechanism, voting decisions, MQ-weighted voting

**Implementation (Trust System):**
- Custom Cosmos SDK module (x/trust)
- Weighted random jury selection (higher MQ = higher probability)
- Evidence submission system
- MQ redistribution on resolution
```go
func SelectJury(eligible []Juror, size int, exclude map[string]bool) []string {
    candidates := filter(eligible, func(j Juror) bool {
        return !exclude[j.Address] && j.MQ > MinJuryMQ
    })
    totalMQ := sum(candidates, func(j Juror) int64 { return j.MQ })

    selected := []string{}
    for len(selected) < size {
        r := randomInt(totalMQ)
        cumulative := int64(0)
        for _, c := range candidates {
            if contains(selected, c.Address) { continue }
            cumulative += c.MQ
            if cumulative >= r {
                selected = append(selected, c.Address)
                break
            }
        }
    }
    return selected
}
```

## Optional Plugins (Layer 1)

> These plugins are installed based on the node's role. Not required for basic operation.

### Service Provider Plugins

#### Plugin A: LLM API Key Hosting

**Purpose:** Host and manage LLM API keys (OpenAI, Anthropic, etc.)

**Key Features:**
- Encrypted key storage (never stored on-chain)
- Usage monitoring and rate limiting
- Automatic shutoff on anomaly detection
- Key rotation support

**Security Model:**
```
[Provider] stores [Encrypted Key] locally
                     |
                [Proxy Layer]
                     |
            [Service Request] -> [Key never exposed to workers]
```

#### Plugin B: Agent Executor (OpenFang)

**Purpose:** Run AI Agents, provide Level 2 services

**Implementation:**
- OpenFang runtime (Rust)
- 28+ pre-built Agent templates
- WASM sandbox for isolation
- 16-layer security protection

**OpenFang Agent Templates:**
| Type | Capabilities |
|------|-------------|
| coder | Code generation, debugging, review |
| researcher | Information gathering, analysis |
| writer | Content creation, editing |
| architect | System design, planning |
| analyst | Data analysis, reporting |

#### Plugin C: Workflow Executor

**Purpose:** Orchestrate multi-Agent workflows, provide Level 3 services

**Implementation:**
- OpenFang Hands system
- 7 pre-built Hands: Collector, Clip, Lead, Content, Trade, Browser, Twitter
- Workflow definition language
- State management

**Workflow Examples:**
| Type | Flow |
|------|------|
| Software Dev | Requirements -> Architecture -> Coding -> Testing -> Deploy |
| Content Creation | Topic -> Research -> Writing -> Review -> Publish |
| Business Plan | Market Analysis -> Strategy -> Review -> Tracking |

### User/Demand Side Plugins

#### Plugin D: GenieBot Interface

**Purpose:** AI conversation, idea collection, service invocation, result display

**Key Features:**
- Natural language interface
- Idea refinement assistance
- Service discovery and recommendation
- Cost estimation
- Progress tracking

**Interaction Example:**
```
User: "I want to build an AI writing assistant"
GenieBot: "This idea needs:
          1. Software Development Workflow (Level 3)
          2. Estimated cost: 500 STT
          Start now?"

User: "Start"
GenieBot: [Invokes Service Market -> Matches Workflow Executor -> Executes Task]
```

## Data Flow

### Service Request Flow

```
[User via GenieBot]
        |
        | 1. Submit ServiceRequest (type, requirements, priceOffer)
        v
[Service Market] --> [Provider Discovery via DHT]
        |
        | 2. Match with provider, create Escrow
        v
[Escrow Module] --> [Lock consumer STT]
        |
        | 3. Notify provider via PubSub
        v
[Provider Node]
        |
        | 4a. Level 1: LLM API call via proxy
        | 4b. Level 2: Agent execution via OpenFang
        | 4c. Level 3: Workflow orchestration
        v
[Service Delivery]
        |
        | 5. Return result with proof
        v
[Consumer] --> [Verify response]
        |
        | 6. Accept/Dispute
        v
[Settlement]
        |
        | 7a. Accept --> Release escrow to provider
        | 7b. Dispute --> Jury resolution --> MQ redistribution
        v
[Complete]
```

### State Management

```
[On-Chain State]                     [Off-Chain State]
+------------------+                 +------------------+
| Account balances |                 | API keys (local) |
| MQ scores        |                 | Service offers   |
| Escrow status    |                 | Response cache   |
| Dispute records  |                 | Metrics/logs     |
| Service registry |                 | Workflow state   |
+--------+---------+                 +--------+---------+
         |                                    |
         +-------- RPC/ABCI queries ---------+
```

## Module Architecture

### Cosmos SDK Custom Modules

```
sharetokens-chain/
+-- app/
|   +-- app.go              # Module wiring
|   +-- encoding.go         # Encoding config
+-- x/
|   +-- market/             # Service marketplace
|   |   +-- keeper/
|   |   +-- types/
|   |   +-- handler.go
|   +-- escrow/             # Payment escrow
|   +-- trust/              # Trust System (reputation + dispute)
|   +-- identity/           # Identity verification
+-- cmd/
    +-- sharetokensd/
```

### Plugin Architecture

```
sharetokens-plugins/
+-- llm-hosting/            # Level 1: LLM API Key Hosting
|   +-- proxy/              # Secure proxy layer
|   +-- storage/            # Encrypted key storage
|   +-- monitor/            # Usage monitoring
+-- agent-executor/         # Level 2: OpenFang integration
|   +-- runtime/            # OpenFang kernel
|   +-- agents/             # Pre-built agents
+-- workflow-executor/      # Level 3: Workflow orchestration
|   +-- hands/              # OpenFang Hands
|   +-- engine/             # Workflow engine
+-- geniebot-ui/            # User interface plugin
    +-- web/                # React frontend
    +-- api/                # API client
```

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|--------------------------|
| 0-1k users | Single full node + light clients; basic P2P |
| 1k-100k users | Multiple full nodes with load balancing; provider clusters |
| 100k+ users | Layer 2 scaling; sharded provider network; hierarchical juries |

## Build Order Implications

### Phase 1: Core Modules Foundation
1. P2P Communication (CometBFT built-in)
2. Identity/Account (Cosmos SDK Auth)
3. Wallet (Cosmos SDK Bank + Keplr)
4. Service Market (basic registration)

**Rationale:** These are the foundation. Without them, nothing works.

### Phase 2: Core Business Modules
1. Escrow Payment (custom module)
2. Trust System (custom module - combines reputation + dispute)
3. Service Market (full functionality)

**Rationale:** These enable the core value proposition - trusted service trading.

### Phase 3: Service Provider Plugins
1. LLM API Key Hosting
2. Agent Executor (OpenFang)
3. Workflow Executor

**Rationale:** These provide actual services. Core must work first.

### Phase 4: User Plugins
1. GenieBot Interface

**Rationale:** User-facing features come after backend is stable.

## Sources

- [Cosmos SDK Documentation](https://docs.cosmos.network/) - Modular blockchain framework
- [CometBFT Documentation](https://cometbft.com/) - Consensus engine
- [OpenFang Documentation](https://openfang.sh/) - AI Agent OS
- [libp2p Documentation](https://docs.libp2p.io) - P2P networking
- [Keplr Wallet](https://keplr.app/) - Cosmos wallet

---
*Architecture research for: ShareTokens Decentralized AI Service Marketplace*
*Updated: 2026-03-02*
*Architecture Version: Core + Plugin*
