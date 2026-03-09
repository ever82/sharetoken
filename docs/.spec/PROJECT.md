# ShareTokens

---

## BASE. Project Overview

### BASE-001. What This Is

ShareTokens is a decentralized AI network. One great idea, the world's AI united to make it real.

### BASE-002. Core Value

**Every great idea deserves the tokens to make it real.**

**Every token deserves the chance to create greater returns.**

**Human minds create and choose. AI executes.**

### BASE-003. Who It's For

**The AI Enthusiast with Ideas but No Technical Skills**

AI agents are booming, but using them requires technical setup — installing software, configuring environments, managing API keys. Many people can only chat with AI, not put it to work. ShareTokens lets you submit ideas and the network's AI executes them. No setup required.

**The Technical Expert Who Doesn't Know What to Build**

You've mastered AI agents, but what do you build? Toy projects waste tokens with no return. ShareTokens connects you with great ideas that need execution. Your skills earn tokens.

**The Visionary Who Needs Trust and Resources**

You have bold ideas and technical skills, but launching alone is too slow. Fundraising is out of reach. Sharing your idea feels risky — what if someone copies it? ShareTokens provides a trust system to protect ideas and match resources, so great projects can take off.

### BASE-004. Why Choose Us

We are an **open-source initiative**, committed to remaining non-commercial and purely dedicated to public welfare.

Our platform is **decentralized**, as we are convinced that decentralization lays the foundation for the highest levels of fairness and justice.

But decentralization does not mean a lack of governance:
- **Trust-based arbitration system** ensures proper oversight
- **Strict real-name verification** — one account per individual
- **Accountability** — those who engage in disruptive activities face fines or permanent bans

---

## CORE. Core Modules

> Essential components — the network cannot run without any of these

### CORE-001. P2P Communication
- Node discovery, message broadcasting
- NAT traversal
- Built on CometBFT P2P

### CORE-002. Identity
- Unique identifiers
- Real-name verification
- Cosmos SDK Auth module

### CORE-003. Wallet
- Balance queries
- Transaction signing
- STT payments

### CORE-004. Service Marketplace (Core Business)

**Three-tier service structure:**

| Level | Service Type | Pricing | Examples |
|-------|-------------|---------|----------|
| Level 1 | LLM API | Per-token billing | GPT-4, Claude API calls |
| Level 2 | Agent | Per-skill billing | coder, researcher, writer |
| Level 3 | Workflow | Package pricing | Software development, content creation |

**Marketplace features:**
- Service registration and discovery
- Pricing strategies
- Intelligent routing

### CORE-005. Escrow
- Lock STT before tasks
- Release on completion
- Freeze during disputes

### CORE-006. Trust System
- Moral Quotient (MQ) scoring
- Tier classification
- Transaction history
- Dispute arbitration
- Jury mechanism
- MQ-weighted voting

---

## PLUGIN. Optional Plugins

> Install based on node role

### PLUGIN-001. Service Provider Plugins

| Plugin | Function | Service Provided |
|--------|----------|------------------|
| LLM API Key Hosting | Host OpenAI/Anthropic API keys | Level 1 service |
| Agent Executor (OpenFang) | Run AI agents | Level 2 service |
| Workflow Executor | Orchestrate multi-agent task flows | Level 3 service |

### PLUGIN-002. User Plugins

| Plugin | Function |
|--------|----------|
| GenieBot Interface | AI chat, service calls, result display |

---

## ARCH. Technical Architecture

### ARCH-001. Core Module Stack

| Module | Technology |
|--------|------------|
| Blockchain Framework | Cosmos SDK + CometBFT |
| P2P Network | CometBFT built-in |
| Consensus | Tendermint |
| Identity | Cosmos SDK Auth |
| Wallet | Keplr integration |

### ARCH-002. Plugin Stack

| Plugin | Technology |
|--------|------------|
| Agent Executor | OpenFang (Rust) |
| Workflow Executor | OpenFang Hands |
| GenieBot Interface | React + TypeScript |

### ARCH-003. Language Distribution

| Layer | Language | Components |
|-------|----------|------------|
| On-chain Core | Go | x/identity, x/market, x/escrow, x/dispute, x/reputation |
| Agent Executor | Rust | OpenFang |
| Off-chain Services | TypeScript | API gateway, service routing |
| Frontend | TypeScript + React | GenieBot Interface |

---

## MARKET. Service Marketplace Details

### MARKET-001. Level 1: LLM API Service

```
Requester → Service Marketplace → LLM Provider
              ↓
         Per-token billing
         Auto-route to optimal provider
```

### MARKET-002. Level 2: Agent Service

```
Requester → Service Marketplace → Agent Executor
              ↓
         Per-skill billing
         Match specialized agents (coder/researcher/writer)
```

### MARKET-003. Level 3: Workflow Service

```
Requester → Service Marketplace → Workflow Executor
              ↓
         Package pricing
         Multi-agent collaboration for complex tasks
```

### MARKET-004. Workflow Examples

| Type | Flow |
|------|------|
| Software Development | Requirements → Architecture → Code Generation → Testing → Deployment |
| Content Creation | Topic Selection → Research → Creation → Review → Publish |
| Business Planning | Market Analysis → Solution Design → Review → Execution Tracking |

---

## GENIE. GenieBot Interface

> GenieBot is the user-facing plugin for the service marketplace, not the workflow executor itself

### GENIE-001. Features
- AI chat entry point
- Idea collection and refinement
- Service discovery and invocation
- Result display

### GENIE-002. Interaction Example

```
User: "I want to build an AI writing assistant"
GenieBot: "This idea requires:
       1. Software Development Workflow (Level 3)
       2. Estimated cost: 500 STT
       Start now?"
User: "Start"
GenieBot: [Calls marketplace → Matches workflow executor → Executes task]
```

---

## REQ. Requirements

### REQ-001. Active

**Core Modules**
- [ ] P2P node discovery and communication
- [ ] Identity registration and real-name verification
- [ ] Wallet balance and transfers
- [ ] Service marketplace three-tier structure
- [ ] Escrow lock/release
- [ ] Moral Quotient (MQ) scoring system
- [ ] Dispute arbitration mechanism

**Service Provider Plugins**
- [ ] LLM API Key hosting
- [ ] OpenFang Agent executor integration
- [ ] Workflow orchestration engine

**User Plugins**
- [ ] GenieBot AI chat interface
- [ ] Service marketplace browsing and invocation

### REQ-002. Out of Scope

- Mobile App (desktop web first)
- Centralized exchange
- NFT marketplace

---

## DECISION. Key Decisions

| Decision | Rationale |
|----------|-----------|
| Core + Plugin architecture | Core modules ensure basic operation, plugins extend as needed |
| Three-tier service marketplace | Serves different needs from simple API to complex workflows |
| Workflow as marketplace service | Workflow is a service type, not a GenieBot feature |
| OpenFang as Agent executor | Provides ready-made agent templates and security isolation |

---

*Last updated: 2026-03-02*
