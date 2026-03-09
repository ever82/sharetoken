# Feature Research

**Domain:** ShareTokens - Decentralized AI Service Marketplace
**Researched:** 2026-03-02
**Confidence:** MEDIUM-HIGH

## Architecture-Based Feature Organization

Features are organized by the Core + Plugin architecture:

- **Core Module Features**: Required for every node
- **Service Provider Plugin Features**: For nodes providing services
- **User Plugin Features**: For nodes consuming services

---

## Core Module Features

### 1. P2P Communication Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| Node Discovery | Find and connect to peers in the network | MEDIUM | P1 |
| Message Broadcasting | Gossip protocol for network-wide messages | MEDIUM | P1 |
| NAT Traversal | Connect nodes behind firewalls | HIGH | P1 |
| Peer Authentication | Verify peer identity before connection | MEDIUM | P1 |
| Connection Management | Handle peer connections, reconnections | MEDIUM | P1 |

### 2. Identity/Account Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| Account Creation | Create new blockchain account | LOW | P1 |
| Real-Name Verification | KYC for anti-Sybil protection | HIGH | P1 |
| Identity Hash Registry | Store verification hash on-chain | MEDIUM | P1 |
| ZK-DID Support | Privacy-preserving identity proofs | HIGH | P2 |
| Account Recovery | Recover lost account access | MEDIUM | P2 |

### 3. Wallet Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| Balance Query | View STT token balance | LOW | P1 |
| Transaction Signing | Sign transactions securely | MEDIUM | P1 |
| STT Transfer | Send STT to other addresses | LOW | P1 |
| Transaction History | View past transactions | LOW | P1 |
| Keplr Integration | Connect via Keplr wallet | MEDIUM | P1 |
| Multi-Signature | Require multiple approvals | HIGH | P3 |

### 4. Service Market Features (Core Business)

**Three-Layer Service Support:**

| Level | Feature | Description | Complexity | Priority |
|-------|---------|-------------|------------|----------|
| All | Service Registration | Register service offering | MEDIUM | P1 |
| All | Service Discovery | Find available services | MEDIUM | P1 |
| All | Pricing Display | Show costs before commitment | LOW | P1 |
| All | Smart Routing | Match requests to optimal providers | HIGH | P1 |
| All | Service Categories | Organize services by type | LOW | P1 |
| Level 1 | LLM API Routing | Route to LLM providers | MEDIUM | P1 |
| Level 2 | Agent Matching | Match to specialized agents | HIGH | P2 |
| Level 3 | Workflow Orchestration | Coordinate multi-agent workflows | HIGH | P2 |

### 5. Escrow Payment Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| Token Locking | Lock STT before service execution | MEDIUM | P1 |
| Conditional Release | Release on success condition | MEDIUM | P1 |
| Dispute Freeze | Freeze funds during disputes | MEDIUM | P1 |
| Timeout Handling | Auto-release on timeout | MEDIUM | P1 |
| Partial Payments | Support partial releases | MEDIUM | P2 |
| Multi-Party Escrow | Escrow for complex workflows | HIGH | P3 |

### 6. Trust System Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| MQ Scoring | Calculate and store MQ scores | HIGH | P1 |
| Zero-Sum Redistribution | Winners gain, losers lose | HIGH | P1 |
| Score Decay | Time-based reputation decay | MEDIUM | P1 |
| Level Classification | Rank users by MQ level | LOW | P1 |
| History Tracking | Record MQ change events | MEDIUM | P1 |
| Anti-Gaming | Prevent score manipulation | HIGH | P2 |
| Dispute Filing | Submit dispute with evidence | MEDIUM | P1 |
| Jury Selection | Weighted random selection | HIGH | P1 |
| Evidence System | Submit and review evidence | MEDIUM | P1 |
| Voting Mechanism | Jurors vote on outcome | MEDIUM | P1 |
| Resolution Enforcement | Execute resolution decision | MEDIUM | P1 |
| Appeal Process | Request review of decision | HIGH | P2 |

---

## Service Provider Plugin Features

### Plugin A: LLM API Key Hosting Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| Key Storage | Securely store API keys | HIGH | P1 |
| Key Encryption | Encrypt keys at rest | HIGH | P1 |
| Usage Proxy | Proxy requests to LLM APIs | MEDIUM | P1 |
| Rate Limiting | Prevent key abuse | MEDIUM | P1 |
| Usage Monitoring | Track key usage | MEDIUM | P1 |
| Auto Shutoff | Disable on anomaly | MEDIUM | P2 |
| Key Rotation | Rotate keys periodically | MEDIUM | P2 |
| Multi-Provider | Support OpenAI, Anthropic, etc. | MEDIUM | P2 |

### Plugin B: Agent Executor (OpenFang) Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| OpenFang Integration | Run OpenFang runtime | HIGH | P1 |
| Agent Templates | 28+ pre-built agents | LOW | P1 |
| WASM Sandbox | Isolated execution | HIGH | P1 |
| Resource Limits | CPU/memory constraints | MEDIUM | P1 |
| Output Capture | Capture agent outputs | MEDIUM | P1 |
| Custom Agents | User-defined agents | HIGH | P2 |
| Agent Marketplace | Share/sell agents | MEDIUM | P3 |

**OpenFang Agent Templates (Level 2 Services):**

| Agent Type | Capabilities |
|------------|-------------|
| coder | Code generation, debugging, review |
| researcher | Information gathering, analysis |
| writer | Content creation, editing |
| architect | System design, planning |
| analyst | Data analysis, reporting |
| tester | Test generation, execution |
| reviewer | Code review, quality check |

### Plugin C: Workflow Executor Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| Workflow Definition | Define multi-step workflows | HIGH | P1 |
| OpenFang Hands | 7 pre-built hands | MEDIUM | P1 |
| State Management | Track workflow state | HIGH | P1 |
| Error Recovery | Handle failures gracefully | HIGH | P1 |
| Parallel Execution | Run steps concurrently | HIGH | P2 |
| Workflow Templates | Pre-built workflows | MEDIUM | P2 |

**OpenFang Hands (Level 3 Services):**

| Hand | Function | Use Case |
|------|----------|----------|
| Collector | Data collection | Idea gathering, market research |
| Clip | Video editing | Content creation workflows |
| Lead | Sales leads | Resource matching |
| Content | Content creation | Article writing, publishing |
| Trade | Trade monitoring | Service trading |
| Browser | Browser automation | Task automation |
| Twitter | Social media | Social matching |

**Workflow Examples:**

| Workflow | Steps |
|----------|-------|
| Software Dev | Requirements -> Architecture -> Coding -> Testing -> Deploy |
| Content Creation | Topic -> Research -> Writing -> Review -> Publish |
| Business Plan | Market Analysis -> Strategy -> Review -> Tracking |

---

## User Plugin Features

### Plugin D: GenieBot Interface Features

| Feature | Description | Complexity | Priority |
|---------|-------------|------------|----------|
| AI Conversation | Natural language dialogue | MEDIUM | P1 |
| Idea Collection | Capture user ideas | MEDIUM | P1 |
| Idea Refinement | Help develop ideas | HIGH | P1 |
| Service Discovery | Find relevant services | MEDIUM | P1 |
| Cost Estimation | Estimate service costs | MEDIUM | P1 |
| Progress Tracking | Track task progress | MEDIUM | P2 |
| Result Display | Show service results | MEDIUM | P1 |
| History View | View past interactions | LOW | P2 |

**GenieBot Interaction Flow:**

```
User Input --> Idea Analysis --> Service Recommendation --> Cost Display
                                     |
                               [User Confirms]
                                     |
                               Service Invocation
                                     |
                               Result Display
```

---

## Feature Dependencies

### Core Module Dependencies

```
[P2P Communication] - foundation for all modules

[Identity/Account]
    └──requires──> [P2P Communication]

[Wallet]
    └──requires──> [Identity/Account]
    └──requires──> [P2P Communication]

[Service Market]
    └──requires──> [Identity/Account]
    └──requires──> [Wallet]
    └──requires──> [P2P Communication]

[Escrow Payment]
    └──requires──> [Wallet]
    └──requires──> [Service Market]

[Trust System]
    └──requires──> [Identity/Account]
    └──requires──> [Escrow Payment]
```

### Plugin Dependencies

```
[LLM API Key Hosting]
    └──requires──> [Service Market]
    └──requires──> [Escrow Payment]

[Agent Executor]
    └──requires──> [Service Market]
    └──requires──> [Escrow Payment]
    └──requires──> [LLM API Key Hosting] (for LLM calls)

[Workflow Executor]
    └──requires──> [Service Market]
    └──requires──> [Escrow Payment]
    └──requires──> [Agent Executor]

[GenieBot Interface]
    └──requires──> [Service Market]
    └──requires──> [Wallet]
```

---

## MVP Definition

### Phase 1: Core Modules MVP

**Launch With (v0.1):**

- [ ] P2P node discovery and messaging
- [ ] Basic account creation
- [ ] Wallet balance and transfers
- [ ] Service registration (Level 1 only)
- [ ] Basic escrow (lock/release)

**Rationale:** These create the minimal viable loop for LLM API sharing.

### Phase 2: Complete Core

**Add After Core MVP (v0.2):**

- [ ] Real-name verification (anti-Sybil)
- [ ] Trust System (reputation + dispute)
- [ ] Service Market (all levels)

**Trigger:** When transactions are flowing and disputes arise.

### Phase 3: Service Provider Plugins

**Add for Providers (v0.3):**

- [ ] LLM API Key Hosting
- [ ] Agent Executor (OpenFang)
- [ ] Workflow Executor

**Trigger:** When core modules are stable.

### Phase 4: User Plugins

**Add for Users (v0.4):**

- [ ] GenieBot Interface

**Trigger:** When provider plugins are ready.

---

## Anti-Features

Features that seem good but cause problems:

| Anti-Feature | Why Requested | Why Problematic | Alternative |
|--------------|---------------|-----------------|-------------|
| Anonymous Accounts | Privacy | Breaks reputation, enables Sybil | Pseudonymous with verified identity |
| On-Chain API Keys | Simplicity | Security risk | Encrypted local storage |
| Fixed Pricing | Simplicity | Market inefficiency | Dynamic market pricing |
| Centralized Discovery | Speed | Single point of failure | DHT-based discovery |
| Mobile-First | Reach | Complex crypto UX | Mobile-responsive web first |
| Free Tier | Adoption | Attracts abuse | Task-earned tokens |

---

## Competitor Feature Comparison

| Feature | ShareTokens | OpenRouter | Golem | Akash |
|---------|-------------|------------|-------|-------|
| **Architecture** | Core + Plugin | Monolithic | Monolithic | Monolithic |
| **Service Levels** | 3 (LLM/Agent/Workflow) | 1 (LLM only) | 1 (Compute) | 1 (Cloud) |
| **MQ** | Trust System (zero-sum) | None | Basic | None |
| **Dispute Resolution** | Jury-based (in Trust System) | None | Limited | None |
| **Identity** | ZK-DID optional | None | None | None |
| **Agent Support** | OpenFang | None | None | None |
| **Workflow** | OpenFang Hands | None | Limited | None |
| **User Interface** | GenieBot | Web | CLI | CLI/Web |

### Competitive Positioning

**What makes ShareTokens unique:**
1. **Core + Plugin architecture** - Modular, extensible
2. **Three service levels** - From simple API to complex workflows
3. **Trust System** - Zero-sum reputation, game-theoretic fairness
4. **OpenFang integration** - Production-ready Agent OS

---

## Sources

- [OpenFang Documentation](https://openfang.sh/) - Agent OS features
- [Cosmos SDK Modules](https://docs.cosmos.network/) - Core module patterns
- [Kleros](https://kleros.io/) - Dispute resolution patterns
- [OpenRouter](https://openrouter.ai/) - LLM marketplace comparison
- [Golem Network](https://golem.network/) - Compute marketplace patterns

---
*Feature research for: ShareTokens Decentralized AI Service Marketplace*
*Updated: 2026-03-02*
*Architecture Version: Core + Plugin*
