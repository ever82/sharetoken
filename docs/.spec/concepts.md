# ShareTokens Core Business Concepts

> Extracted from `.planning/models` directory. Defines the most core business concepts of the system.

---

## MARKET. Marketplace Concepts

### MARKET-001. Service Marketplace

Three-tier service trading infrastructure connecting service providers and consumers.

- **Level 1 - LLM API Service**: Basic model calls (GPT-4, Claude, etc.), priced per token
- **Level 2 - Agent Service**: Autonomous task execution (code generation, data analysis), priced by complexity
- **Level 3 - Workflow Service**: Multi-step orchestration (software development, content creation), priced by milestone

### MARKET-002. Task Marketplace

Full task lifecycle management.

- **Task**: Decomposable, assignable work units
- **Application/Bidding**: Open application or competitive bidding mode
- **Milestones**: Phased delivery and payment
- **Review**: Multi-dimensional scoring (quality/communication/timeliness/expertise)

### MARKET-003. Idea System

Idea incubation and crowdfunding platform.

- **Idea**: Versionable, collaborative creative documents
- **Crowdfunding**: Supports investment, lending, and donation types
- **Contribution Tracking**: Records weights by contribution type (code/design/docs)
- **Revenue Distribution**: Distributed by cumulative contribution weight

---

## TRUST. Trust & Security Concepts

### TRUST-001. MQ (Moral Quotient)

Reputation scoring system representing a user's credibility on the platform.

- Initial value 100, zero-sum game (total supply constant)
- Voting weight = MQ value itself
- Maximum 3% loss per dispute, never goes below 0
- Higher MQ users lose more when deviating from consensus (convergence mechanism)

### TRUST-002. Trust System

Platform trust infrastructure with two core functions:

- **MQ Scoring**: Zero-sum reputation system
- **Dispute Arbitration**: AI mediation first, jury voting as fallback

### TRUST-003. Dispute Arbitration

Decentralized dispute resolution mechanism.

- **AI Mediation**: Free conversation, evidence submission, proposal scoring
- **Jury Voting**: MQ-weighted random jury selection when AI mediation fails
- **MQ Redistribution**: Those deviating from consensus are penalized, those close to consensus are rewarded

### TRUST-004. Escrow

Transaction fund security mechanism.

- Funds locked in escrow account before transaction
- Automatically frozen when dispute arises
- Released to provider after service completion
- Distributed proportionally after dispute resolution

---

## IDENTITY. Identity & Access Concepts

### IDENTITY-001. Real-name Identity

Strict real-name verification with privacy protection.

- Supports WeChat, GitHub, Google and other third-party verification
- Only hashes stored on-chain, no plaintext
- Global identity registry prevents duplicate registration
- Local Merkle proof verification

### IDENTITY-002. User Limits

Risk control mechanism based on identity verification level and MQ.

- Trading limits: Per-transaction/daily/monthly
- Withdrawal limits: Daily limit, cooldown period
- Dispute limits: Maximum active disputes
- Service limits: Concurrent calls, rate limiting

---

## SERVICE. Service Infrastructure Concepts

### SERVICE-001. API Key Custody

Secure storage and management of LLM Provider API Keys.

- Encrypted storage on-chain
- Decrypted and used within WASM sandbox
- Immediately erased after use (Secret Zeroization)
- Supports access control and pricing configuration

### SERVICE-002. Oracle Service (Exchange Rate Layer)

Decentralized price data service.

- Obtains exchange rate data via Chainlink
- Standardizes official prices from various LLMs to STT
- Supports price subscription and caching

### SERVICE-003. OpenFang Integration

Integration architecture with open-source AI Agent OS.

- **Agent Templates**: 28+ templates including coder, researcher, writer
- **Hands**: 7 autonomous capability packs including Collector, Lead, Researcher
- **Security**: 16-layer security mechanism, WASM sandbox isolation
- **Sidecar Deployment**: Tightly coupled provider plugin with chain node

---

## NODE. Node & Interface Concepts

### NODE-001. Node Types

Network nodes with different roles.

- **Light Node**: Regular users, core modules + GenieBot plugin
- **Full Node**: Validators, complete state and history
- **Service Node**: Service providers, core modules + service plugins
- **Archive Node**: Block explorers, complete history index

### NODE-002. GenieBot

User interaction interface, client for the service marketplace.

- Natural language dialogue entry point
- Intent recognition and service recommendation
- One-click access to LLM/Agent/Workflow services
- Task management and progress tracking
