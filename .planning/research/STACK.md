# Stack Research

**Domain:** ShareTokens - Decentralized AI Service Marketplace
**Researched:** 2026-03-02
**Confidence:** HIGH (for core components), MEDIUM (for supporting libraries)

## Architecture-Aligned Stack Selection

The stack is organized by the Core + Plugin architecture:

- **Core Module Stack**: Technologies for mandatory node components
- **Service Provider Plugin Stack**: Technologies for provider plugins
- **User Plugin Stack**: Technologies for user-facing plugins

---

## Core Module Stack

### 1. P2P Communication Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **CometBFT P2P** | v0.38.x | Built-in P2P networking | Included with Cosmos SDK; handles consensus P2P; no additional setup |

**Protocol Stack:**
```
+------------------+
|   Application    |  Service discovery, messaging
+--------+---------+
         |
+--------+---------+
|   CometBFT P2P   |  PEX Reactor, Address Book, Reactor Gossip
+--------+---------+
         |
+--------+---------+
|  Transport Layer |  TCP, WebSocket, QUIC
+--------+---------+
         |
+--------+---------+
|   Encryption     |  SecretConnection (TLS 1.3)
+------------------+
```

### 2. Identity/Account Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Cosmos SDK Auth** | v0.50.x | Account management | Built-in module; signature verification; address derivation |
| **Cosmos SDK Authz** | v0.50.x | Authorization grants | Grant permissions to other accounts |
| **ZK-DID (optional)** | circom/snarkjs | Privacy-preserving KYC | Zero-knowledge identity proofs |

### 3. Wallet Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Cosmos SDK Bank** | v0.50.x | Token transfers | Built-in module; STT token operations |
| **Keplr Wallet** | latest | User wallet | Standard Cosmos ecosystem wallet; browser extension + mobile |
| **WalletConnect** | v2.x | Mobile wallet connection | Universal wallet connection protocol |
| **cosmjs** | v0.32.x | JavaScript SDK | Client-side transaction building and signing |

### 4. Service Market Stack (Core Business)

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Cosmos SDK (Custom Module)** | v0.50.x | Service registry | x/market module; on-chain service listings |
| **CometBFT** | v0.38.x | Consensus for transactions | BFT consensus; instant finality |
| **gRPC** | built-in | Service queries | High-performance queries; Cosmos SDK native |
| **GraphQL (optional)** | v1.x | Flexible queries | Alternative query interface |

### 5. Escrow Payment Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Cosmos SDK (Custom Module)** | v0.50.x | x/escrow module | Custom escrow logic; time-locked releases |
| **Cosmos SDK Bank** | v0.50.x | Token operations | Lock, release, transfer STT |
| **Chainlink Automation** | 2025 | Time-based triggers | Automated timeout releases (optional) |

### 6. Trust System Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Cosmos SDK (Custom Module)** | v0.50.x | x/trust module | MQ scoring; dispute resolution |
| **Cosmos SDK State** | v0.50.x | Persistent storage | On-chain MQ scores and history |
| **IPFS (optional)** | latest | Evidence storage | Large evidence files off-chain |
| **Chainlink VRF** | v2.x | Randomness | Unpredictable jury selection (optional) |

---

## Core Stack Summary

| Component | Technology | Version | Language |
|-----------|------------|---------|----------|
| **Blockchain Framework** | Cosmos SDK | v0.50.x | Go |
| **Consensus Engine** | CometBFT | v0.38.x | Go |
| **Account Module** | Cosmos SDK Auth | v0.50.x | Go |
| **Token Module** | Cosmos SDK Bank | v0.50.x | Go |
| **Custom Modules** | x/market, x/escrow, x/trust, x/identity | custom | Go |
| **Cross-Chain** | IBC-Go | v8.x | Go |
| **Client SDK** | cosmjs | v0.32.x | TypeScript |

---

## Service Provider Plugin Stack

### Plugin A: LLM API Key Hosting Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Rust** | 1.75+ | Plugin implementation | Memory safety; high performance; no GC |
| **AEAD Encryption** | ring crate | Key encryption | Authenticated encryption for key storage |
| **HTTP Proxy** | hyper/reqwest | LLM API proxy | High-performance async HTTP |
| **SQLite/ RocksDB** | latest | Local encrypted storage | Encrypted key storage; offline capability |
| **Prometheus** | latest | Usage metrics | Standard monitoring; alerts on anomalies |

**Security Architecture:**
```
+------------------+
|   API Request    |
+--------+---------+
         |
+--------+---------+
|   Proxy Layer    |  Key never exposed
+--------+---------+
         |
+--------+---------+
|  Encrypted Store |  Keys encrypted at rest
+--------+---------+
         |
+--------+---------+
|   LLM Provider   |  OpenAI, Anthropic, etc.
+------------------+
```

### Plugin B: Agent Executor (OpenFang) Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **OpenFang** | latest | AI Agent OS | 28+ pre-built agents; WASM sandbox; 16-layer security |
| **Rust** | 1.75+ | OpenFang runtime | Memory safety; performance; WASM support |
| **WASM Runtime** | wasmtime | Agent sandbox | Isolation; capability-based security |
| **LLM Clients** | async-openai, anthropic | LLM integration | API clients for major providers |

**OpenFang Components:**
| Component | Description |
|-----------|-------------|
| OpenFang Kernel | Agent runtime, WASM sandbox |
| OpenFang Agents | 28+ templates (coder, researcher, writer, etc.) |
| OpenFang CLI | `openfang init`, `openfang start` |
| OpenFang Dashboard | localhost:4200 |
| OpenFang SDK | JavaScript, Python |

**OpenFang Security (16 Layers):**
1. WASM Sandbox
2. Taint Tracking
3. Audit Logs
4. SSRF Protection
5. Least Privilege
6. Input Validation
7. Output Encoding
8. Key Management
9. Network Isolation
10. Resource Limits
11. Behavior Monitoring
12. Rollback Mechanism
13. Encrypted Transport
14. Access Control
15. Compliance Check
16. Threat Intelligence

### Plugin C: Workflow Executor Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **OpenFang Hands** | latest | Workflow capabilities | 7 pre-built hands; proven automation |
| **Rust** | 1.75+ | Workflow engine | High performance; async support |
| **Temporal (optional)** | latest | Workflow orchestration | Enterprise-grade workflow engine |
| **SQLite** | latest | State persistence | Checkpoint storage; recovery support |

**OpenFang Hands:**
| Hand | Function | Use Case |
|------|----------|----------|
| Collector | Data collection | Idea gathering, market research |
| Clip | Video editing | Content workflows |
| Lead | Sales leads | Resource matching |
| Content | Content creation | Article writing |
| Trade | Trade monitoring | Service trading |
| Browser | Browser automation | Task automation |
| Twitter | Social media | Social matching |

---

## User Plugin Stack

### Plugin D: GenieBot Interface Stack

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **React** | 18.x | UI framework | Component-based; large ecosystem |
| **TypeScript** | 5.x | Type safety | Type-safe development; better DX |
| **TailwindCSS** | 3.x | Styling | Utility-first; rapid prototyping |
| **cosmjs** | v0.32.x | Blockchain client | Transaction building; wallet connection |
| **Keplr Integration** | latest | Wallet connection | Standard Cosmos wallet |
| **Zustand/Jotai** | latest | State management | Lightweight; React integration |

**Frontend Architecture:**
```
+------------------+
|  GenieBot UI     |  React + TypeScript
+--------+---------+
         |
+--------+---------+
|   State Layer    |  Zustand/Jotai
+--------+---------+
         |
+--------+---------+
|  cosmjs Client   |  Blockchain queries
+--------+---------+
         |
+--------+---------+
|  Keplr Wallet    |  Transaction signing
+--------+---------+
         |
+--------+---------+
|  Service Market  |  Core module via RPC
+------------------+
```

---

## Infrastructure Stack

| Technology | Purpose | Notes |
|------------|---------|-------|
| **Docker** | Containerization | Node deployment; reproducibility |
| **Kubernetes** | Orchestration | Provider node clusters (optional) |
| **Prometheus/Grafana** | Monitoring | Standard Cosmos monitoring |
| **CometBFT Remote Signer** | Key management | HSM integration for validators |
| **Ignite CLI** | Development | Cosmos chain scaffolding |
| **cosmovisor** | Upgrades | Automated binary swap |

---

## Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| **Go 1.21+** | Chain development | Primary language for Cosmos SDK |
| **Rust 1.75+** | Plugin development | Provider plugins; OpenFang |
| **TypeScript 5.x** | Frontend/SDK | Type-safe client development |
| **golangci-lint** | Go linting | Static analysis |
| **cargo/clippy** | Rust linting | Static analysis |
| **Protocol Buffers** | Type definitions | Cross-language serialization |

---

## Version Compatibility Matrix

| Package A | Compatible With | Notes |
|-----------|-----------------|-------|
| cosmos-sdk v0.50.x | cometbft v0.38.x | Officially supported |
| cosmos-sdk v0.52.x | cometbft v0.38.x | Beta; may have breaking changes |
| ibc-go v8.x | cosmos-sdk v0.50.x | Required for IBC |
| cosmjs v0.32.x | cosmos-sdk v0.50.x | Client compatibility |
| OpenFang latest | Rust 1.75+ | Agent runtime |

---

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| **Tendermint Core (legacy)** | Forked to CometBFT; no longer maintained | CometBFT v0.38.x |
| **Cosmos SDK < v0.47** | End of life; no security updates | Cosmos SDK v0.50.x |
| **Web3.js / Ethers.js** | Not an EVM chain | cosmjs |
| **On-Chain Key Storage** | Security risk | Encrypted local storage |
| **Centralized Discovery** | Single point of failure | DHT-based discovery |
| **Ethereum as Base Layer** | Gas costs prohibitive | Cosmos app-chain |
| **GraphQL-only APIs** | Cosmos uses gRPC/REST natively | gRPC for performance |

---

## Stack by Build Phase

### Phase 1: Core Modules Foundation

| Component | Technology |
|-----------|------------|
| Blockchain Framework | Cosmos SDK v0.50.x |
| Consensus | CometBFT v0.38.x |
| Account | Cosmos SDK Auth |
| Wallet | Cosmos SDK Bank + Keplr |
| Service Market (basic) | Custom x/market module |
| Development Tool | Ignite CLI |

**Language Mix:**
- Go: Chain modules (primary)
- TypeScript: Client SDK

### Phase 2: Core Business Modules

| Component | Technology |
|-----------|------------|
| Escrow | Custom x/escrow module |
| Trust System | Custom x/trust module |
| Service Market (full) | Custom x/market module |
| Oracle (optional) | Chainlink Network |

**Language Mix:**
- Go: All modules
- TypeScript: Client updates

### Phase 3: Service Provider Plugins

| Component | Technology |
|-----------|------------|
| LLM Hosting | Rust + AEAD encryption |
| Agent Executor | OpenFang (Rust) |
| Workflow Executor | OpenFang Hands (Rust) |

**Language Mix:**
- Rust: All provider plugins (primary)
- Go: Core integration

### Phase 4: User Plugins

| Component | Technology |
|-----------|------------|
| GenieBot UI | React + TypeScript |
| State Management | Zustand/Jotai |
| Wallet Connection | Keplr + cosmjs |

**Language Mix:**
- TypeScript: All user plugin code

---

## Language Distribution Summary

| Layer | Language | Components |
|-------|----------|------------|
| Core Modules (Chain) | Go | x/market, x/escrow, x/trust, x/identity |
| Provider Plugins | Rust | LLM Hosting, OpenFang Agent/Workflow |
| User Plugin | TypeScript | GenieBot UI |
| Client SDK | TypeScript | cosmjs-based client library |

---

## Installation Guide

```bash
# Core Modules (Go)
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
# In go.mod: require github.com/cosmos/cosmos-sdk v0.50.10

# Ignite CLI (Development)
curl https://get.ignite.com/cli | bash

# Provider Plugins (Rust)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
cargo install openfang

# User Plugin (TypeScript)
npx create-react-app geniebot-ui --template typescript
npm install @cosmjs/stargate @cosmjs/launchpad

# Keplr Wallet (Browser Extension)
# Install from: https://keplr.app/
```

---

## Sources

- [Cosmos SDK Documentation](https://docs.cosmos.network/) - Framework docs (HIGH confidence)
- [CometBFT Documentation](https://cometbft.com/) - Consensus engine (HIGH confidence)
- [OpenFang Documentation](https://openfang.sh/) - Agent OS (HIGH confidence)
- [Keplr Wallet](https://keplr.app/) - Cosmos wallet (HIGH confidence)
- [cosmjs GitHub](https://github.com/cosmos/cosmjs) - Client SDK (HIGH confidence)
- [Chainlink Documentation](https://docs.chain.link/) - Oracle network (HIGH confidence)

---
*Stack research for: ShareTokens Decentralized AI Service Marketplace*
*Updated: 2026-03-02*
*Architecture Version: Core + Plugin*
