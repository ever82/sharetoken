# ShareTokens Technology Stack

> Technology decisions for the ShareTokens decentralized AI service marketplace.

---

## ARCH. Architecture

### ARCH-001. Core + Plugin Pattern

**Decision:** Modular architecture with mandatory core modules and optional plugins.

**Rationale:**
- Core modules: every node must have these, operational
- Plugins: install based on node role (provider vs user)

---

## CORE. Core Modules

> Every node must have these modules. Missing any prevents the node from functioning.

### CORE-001. P2P Communication
| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | CometBFT built-in P2P | Integrated with consensus |
| Peer Discovery | PEX Reactor + Address Book | Built-in, battle-tested in Cosmos ecosystem |
| Broadcast | CometBFT Reactor Gossip | Native consensus messaging, channel-based routing |
| Encryption | SecretConnection (TLS 1.3) | Built-in authenticated encryption with PFS |
| NAT Traversal | UPnP + external_address config | Standard Cosmos approach |

**Notes:**
- PEX (Peer Exchange) automatically discovers and exchanges peer addresses
- Address Book persists known peers across restarts
- For production validators: manually configure `external_address`
- For home users: enable UPnP or use cloud VPS

### CORE-002. Identity
| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | Cosmos SDK Auth | Built-in account management |
| Privacy KYC | ZK-DID (optional) | Privacy-preserving verification |

### CORE-003. Wallet & Authentication

| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | Cosmos SDK Bank | Built-in token management |
| Desktop Wallet | Keplr | Standard Cosmos ecosystem wallet |
| Mobile Wallet | WalletConnect | Cross-platform mobile support |
| Native Token | STT | Platform denomination |
| Social Login | GitHub OAuth (primary), Zero-friction onboarding |
| Auto Wallet Creation | Backend automation | New users don't need to understand mnemonics |

#### Wallet Connection Flow

**Web Desktop:**
1. **Existing Keplr Users**: Click "Connect Wallet" → Select Keplr → Approve connection
2. **New Users**: Click "Quick Start with GitHub" → OAuth → Auto-create Keplr wallet → Show address

   - Backend automatically creates a Keplr wallet (encrypted mnemonic storage)
   - User receives wallet address immediately usable
   - User can export mnemonic to full Keplr anytime

**Mobile:**
1. Scan WalletConnect QR code
2. Open Keplr mobile app
3. Approve connection
4. Complete login

**Backend Auto-Creation Flow:**
```go
// 1. User authenticates via GitHub OAuth
// 2. Check if wallet exists for GitHub ID
// 3. If not: Create new Keplr wallet
//    - Generate mnemonic (24 words)
//    - Create account on chain
//    - Encrypt and store mnemonic in database
// 4. Return wallet address to frontend
// 5. User can later export mnemonic to gain full control
```

**User Experience Benefits:**
- **Zero Friction**: No need to understand mnemonics for new users
- **Familiar Login**: GitHub OAuth is a known, trusted method
- **Progressive Ownership**: Start custustodial, export to full Keplr when ready
- **Ecosystem Compatible**: All wallets are standard Keplr wallets

### CORE-004. Service Marketplace
| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | Custom Cosmos SDK module (x/market) | Domain-specific logic |
| Service Levels | 3 tiers (LLM/Agent/Workflow) | Cover different complexity needs |
| Pricing | Fixed, Dynamic, Auction | Flexible monetization |

### CORE-005. Escrow Payment
| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | Custom Cosmos SDK module (x/escrow) | Domain-specific logic |
| Lock Mechanism | Module-controlled escrow accounts | State-machine based fund custody |
| Release Triggers | Completion, Dispute resolution, Timeout | Event-driven automated settlement |
| Dispute Integration | EscrowHooks interface | Jury verdict triggers fund distribution |

**Notes:**
- Funds held in module account (not multi-sig)
- Release controlled by escrow state machine
- Dispute resolution triggers automatic fund distribution via hooks

### CORE-006. Trust System
| Decision | Choice | Reason |
|----------|--------|--------|
| Framework | Custom Cosmos SDK module (x/trust) | Domain-specific logic |
| MQ Scoring | Zero-sum game | Anti-inflationary reputation |
| Dispute Resolution | Jury voting with MQ weights | Fair, decentralized arbitration |

---

## PLUGIN. Optional Plugins

> Install based on node role. Not required for basic operation.

### PLUGIN-001. Service Provider Plugins

| Plugin | Technology | Purpose |
|--------|------------|---------|
| LLM API Key Hosting | Custom (encrypted storage + proxy) | Host OpenAI/Anthropic keys securely |
| Agent Executor | OpenFang (Rust) | Run AI Agents, Level 2 service |
| Workflow Executor | OpenFang Hands | Orchestrate multi-agent workflows, Level 3 service |

**OpenFang Details:**
- 28+ pre-built Agent templates (coder, researcher, writer, etc.)
- 16-layer security protection
- WASM sandbox isolation

### PLUGIN-002. User Plugins

| Plugin | Technology | Purpose |
|--------|------------|---------|
| GenieBot Interface | React + TypeScript | User-facing AI chat interface |

---


