# Project Research Summary

**Project:** ShareTokens
**Domain:** Decentralized AI Service Marketplace
**Researched:** 2026-03-02
**Confidence:** MEDIUM-HIGH

## Executive Summary

ShareTokens is a **decentralized AI service marketplace** built on a **Core + Plugin architecture**. The system enables trading of three levels of AI services: LLM API access (Level 1), Agent execution (Level 2), and Workflow orchestration (Level 3).

### Architecture Overview

```
+-----------------------------------------------------------------------------+
|                        Core Modules (Every Node Must Have)                   |
+-----------------------------------------------------------------------------+
|  P2P Communication | Identity/Account | Wallet | Service Market (Core)       |
|  Escrow Payment | Trust System (MQ + Dispute)                                   |
+-----------------------------------------------------------------------------+
                               |
                               v
+-----------------------------------------------------------------------------+
|                        Optional Plugins (Install By Role)                    |
+-----------------------------------------------------------------------------+
|  Service Provider Plugins:                                                   |
|  - LLM API Key Hosting (Level 1)                                            |
|  - Agent Executor / OpenFang (Level 2)                                      |
|  - Workflow Executor (Level 3)                                              |
|                                                                              |
|  User Plugin:                                                                |
|  - GenieBot Interface - AI conversation, service invocation                 |
+-----------------------------------------------------------------------------+
```

### Core Value Proposition

**Every great idea deserves the tokens to make it real.**

| Role | Need | ShareTokens Solution |
|------|------|----------------------|
| API Key Holders | Monetize idle resources | Provide Level 1 LLM services, earn STT |
| Agent Developers | Monetize skills | Provide Level 2 Agent services, earn STT |
| Workflow Builders | Monetize processes | Provide Level 3 Workflow services, earn STT |
| Idea Owners | Access AI capabilities | Purchase services on-demand via GenieBot |

---

## Key Findings

### 1. Core Module Stack (Mandatory)

Every node must run these modules. Built primarily in **Go** using Cosmos SDK.

| Module | Technology | Purpose |
|--------|------------|---------|
| P2P Communication | CometBFT built-in | Node discovery, messaging |
| Identity/Account | Cosmos SDK Auth | Account management, anti-Sybil |
| Wallet | Cosmos SDK Bank + Keplr | STT token operations |
| Service Market | Custom x/market | Three-level service registry |
| Escrow Payment | Custom x/escrow | Lock/release STT for services |
| Trust System | Custom x/trust | MQ + Dispute resolution |

### 2. Plugin Stack (Optional)

Installed based on node role. Built in **Rust** (providers) or **TypeScript** (users).

| Plugin | Technology | Provides |
|--------|------------|----------|
| LLM API Key Hosting | Rust + AEAD encryption | Level 1 services |
| Agent Executor | OpenFang (Rust) | Level 2 services |
| Workflow Executor | OpenFang Hands | Level 3 services |
| GenieBot Interface | React + TypeScript | User-facing UI |

### 3. Three-Level Service Model

The Service Market supports three service levels, each with distinct characteristics:

| Level | Service Type | Billing | Provider Plugin |
|-------|-------------|---------|-----------------|
| Level 1 | LLM API | Per-token | LLM API Key Hosting |
| Level 2 | Agent | Per-skill | Agent Executor (OpenFang) |
| Level 3 | Workflow | Per-workflow | Workflow Executor |

### 4. Critical Pitfalls Identified

| Pitfall | Affected Component | Prevention Phase |
|---------|-------------------|------------------|
| Plugin Architecture Complexity | Core-Plugin Interface | Phase 1 |
| Sybil Attack | Identity Module | Phase 1 |
| Escrow Fund Lock | Escrow Module | Phase 2 |
| MQ Gaming | Trust System | Phase 2 |
| Jury Collusion | Trust System | Phase 2 |
| API Key Theft | LLM Hosting Plugin | Phase 3 |
| Agent Sandbox Escape | Agent Executor Plugin | Phase 3 |
| GenieBot Misinterpretation | GenieBot Plugin | Phase 4 |

---

## Architecture Decisions

### Core + Plugin Pattern

**Decision:** Separate mandatory core modules from optional plugins.

**Rationale:**
- Core modules provide essential infrastructure (P2P, identity, payments, reputation)
- Plugins enable flexible service provision based on node capabilities
- Allows lightweight nodes (users) vs. full-featured nodes (providers)
- Supports independent development and versioning

**Trade-offs:**
- Pros: Modularity, flexibility, independent scaling
- Cons: API versioning complexity, integration testing overhead

### Zero-Sum MQ (Trust System)

**Decision:** Use zero-sum MQ where winners gain exactly what losers lose.

**Rationale:**
- Prevents MQ inflation
- Creates strong incentives for honest behavior
- Naturally balances over time
- Novel approach requiring careful game-theoretic validation

**Trade-offs:**
- Pros: No inflation, game-theoretic fairness, automatic balancing
- Cons: Can feel punitive to newcomers, requires active participation

### OpenFang for Agent/Workflow Execution

**Decision:** Use OpenFang as the Agent OS for Level 2 and Level 3 services.

**Rationale:**
- 28+ pre-built agent templates
- 7 pre-built Hands for workflows
- 16-layer security model
- WASM sandbox for isolation
- Production-ready Rust implementation

**Trade-offs:**
- Pros: Immediate agent capability, proven security, active development
- Cons: Dependency on external project, learning curve for customization

---

## Implications for Roadmap

### Phase 1: Core Modules Foundation

**Delivers:** P2P Communication, Identity/Account, Wallet, Basic Service Market

**Key Technologies:**
- Cosmos SDK v0.50.x (Go)
- CometBFT v0.38.x
- Keplr Wallet integration

**Addresses Pitfalls:**
- Plugin Architecture Complexity (API contracts)
- Sybil Attack (identity verification)
- Unsustainable Emissions (economic model)

### Phase 2: Core Business Modules

**Delivers:** Escrow, Trust System, Full Service Market

**Key Technologies:**
- Custom Cosmos SDK modules (x/escrow, x/trust)
- Chainlink (optional for automation)

**Addresses Pitfalls:**
- Escrow Fund Lock (timeout handling)
- MQ Gaming (anti-gaming mechanisms)
- Jury Collusion (anonymous selection)
- Reliability Gap (redundancy)

### Phase 3: Service Provider Plugins

**Delivers:** LLM Hosting, Agent Executor, Workflow Executor

**Key Technologies:**
- Rust for plugin implementation
- OpenFang for Agent/Workflow runtime
- AEAD encryption for key storage

**Addresses Pitfalls:**
- API Key Theft (encryption, proxy layer)
- Agent Escape (WASM sandbox, 16-layer security)
- Workflow Corruption (checkpointing, recovery)

### Phase 4: User Plugins

**Delivers:** GenieBot Interface

**Key Technologies:**
- React 18.x + TypeScript 5.x
- cosmjs for blockchain client
- Keplr for wallet connection

**Addresses Pitfalls:**
- GenieBot Misinterpretation (confirmation, transparency)

---

## Feature Priorities

### Core Module Features (P1)

| Module | P1 Features |
|--------|-------------|
| P2P | Node discovery, message broadcasting, NAT traversal |
| Identity | Account creation, real-name verification |
| Wallet | Balance query, STT transfer, Keplr integration |
| Service Market | Service registration, discovery, pricing display |
| Escrow | Token locking, conditional release |
| Trust System | MQ scoring, zero-sum redistribution, dispute resolution |

### Plugin Features (P1-P2)

| Plugin | P1 Features | P2 Features |
|--------|-------------|-------------|
| LLM Hosting | Key storage, encryption, proxy | Auto shutoff, key rotation |
| Agent Executor | OpenFang integration, templates | Custom agents |
| Workflow Executor | Workflow definition, Hands | Parallel execution |
| GenieBot | AI conversation, service discovery | Progress tracking |

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack Selection | HIGH | Cosmos SDK, CometBFT, OpenFang are mature, well-documented |
| Architecture Pattern | HIGH | Core + Plugin is proven pattern; clear separation of concerns |
| Core Modules | MEDIUM-HIGH | Standard Cosmos SDK modules; some novel components (Trust System) |
| Plugin Integration | MEDIUM | OpenFang integration needs validation; API contracts need testing |
| Feature Priorities | MEDIUM | Based on competitor analysis; needs user validation |
| Pitfall Prevention | MEDIUM | Based on research and domain expertise; needs testing |

**Overall confidence:** MEDIUM-HIGH

### Gaps to Address

1. **Zero-sum MQ validation:** Novel mechanism needs game-theoretic modeling and simulation
2. **OpenFang integration specifics:** API compatibility, performance characteristics
3. **Plugin API versioning:** Need explicit versioning strategy before Phase 3
4. **GenieBot NLU accuracy:** Natural language understanding for service matching

---

## Build Order Summary

```
Phase 1: Core Modules Foundation (Go)
    |
    v
Phase 2: Core Business Modules (Go)
    |
    v
Phase 3: Service Provider Plugins (Rust)
    |
    v
Phase 4: User Plugins (TypeScript)
```

**Language Distribution:**
- **Go**: Core modules (chain layer)
- **Rust**: Provider plugins (performance, security)
- **TypeScript**: User plugin (frontend, SDK)

---

## Sources

### Primary (HIGH confidence)
- [Cosmos SDK Documentation](https://docs.cosmos.network/) - Core framework
- [CometBFT Documentation](https://cometbft.com/) - Consensus engine
- [OpenFang Documentation](https://openfang.sh/) - Agent OS
- [Keplr Wallet](https://keplr.app/) - Cosmos wallet

### Secondary (MEDIUM confidence)
- [Kleros Documentation](https://kleros.io/) - Dispute resolution patterns
- [OpenRouter](https://openrouter.ai/) - LLM marketplace comparison
- [Golem Network](https://golem.network/) - Compute marketplace patterns

### Tertiary (LOW confidence)
- DePIN Analysis articles - Market analysis
- Various blog posts on decentralized compute

---
*Research completed: 2026-03-02*
*Architecture Version: Core + Plugin*
*Ready for roadmap: yes*
