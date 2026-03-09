# State: ShareTokens

## Project Reference

See: .planning/PROJECT.md (updated 2026-03-02)

**Core value:** Every great idea deserves the tokens to make it real.

**Current focus:** Phase 0 - Interface Definitions

---

## Architecture

```
+-------------------------------------------------------------+
|                    Core Modules (Required)                   |
|  P2P | Identity | Wallet | Market | Escrow | Trust System   |
+-------------------------------------------------------------+
                              |
+-------------------------------------------------------------+
|                    Optional Plugins                          |
|  Providers: LLM Host | Agent Executor | Workflow Executor    |
|  Users: GenieBot Interface                                   |
+-------------------------------------------------------------+
```

---

## Module Progress

### Core Modules (0%)

| Module | Status | Description | Tech |
|--------|--------|-------------|------|
| P2P Communication | Pending | Node discovery, messaging, NAT traversal | CometBFT Built-in |
| Identity | Pending | Unique ID, verification, privacy protection | Cosmos SDK Auth |
| Wallet | Pending | Balance, signing, STT transactions | Cosmos SDK + Keplr |
| Service Market | Pending | 3-layer services (LLM/Agent/Workflow) | Cosmos SDK x/service |
| Escrow Payment | Pending | Lock, release, dispute freeze | Cosmos SDK x/escrow |
| Trust System | In Progress | MQ scoring, dispute arbitration | Cosmos SDK x/trust |

### Provider Plugins (0%)

| Plugin | Status | Description |
|--------|--------|-------------|
| LLM API Key Host | Pending | Secure storage, multi-provider support |
| Agent Executor (OpenFang) | Pending | 28+ templates, 7 Hands, 16 security layers |
| Workflow Executor | Pending | Multi-agent orchestration, GitHub integration |

### User Plugins (0%)

| Plugin | Status | Description |
|--------|--------|-------------|
| GenieBot Interface | Pending | Natural language, idea processing, service calls |

---

## Statistics

### Requirements Coverage

| Category | Count | Percentage |
|----------|-------|------------|
| **Core Modules** | 76 | 64% |
| - P2P Communication | 5 | - |
| - Identity | 5 | - |
| - Wallet | 5 | - |
| - Service Market | 47 | - |
| - Escrow Payment | 6 | - |
| - Trust System | 14 | - |
| **Plugin Modules** | 42 | 36% |
| - LLM Provider Plugin | 3 | - |
| - Agent Executor Plugin | 10 | - |
| - Workflow Executor Plugin | 6 | - |
| - GenieBot Plugin | 23 | - |
| **Total** | 118 | 100% |

### Development Estimates

| Category | Modules | Est. LOC |
|----------|---------|----------|
| Core (Go) | 5 Cosmos SDK modules | ~15,000 |
| Plugins (Rust/TS) | 4 plugins | ~8,000 |
| Frontend (React) | 1 UI | ~5,000 |
| **Total** | 10 | ~28,000 |

---

## Blockers

(None)

---

## Quick Tasks Completed

| # | Description | Date | Commit | Directory |
|---|-------------|------|--------|-----------|
| 1 | Trust System Module Foundation | 2026-03-03 | 199ebe5 | [1-trust-system-module-development](./quick/1-trust-system-module-development/) |

---

## Recent Activity

- 2026-03-03: **Trust System Foundation** - Proto definitions, types, keeper, MQ scoring (Quick Task 1)
- 2026-03-02: **Architecture reorganized** - Core modules + optional plugins
- 2026-03-02: **Cosmos SDK adopted** - Unified blockchain framework
- 2026-03-02: **Trust System consolidated** - MQ + Dispute merged
- 2025-03-02: OpenFang adopted as AI Agent OS
- 2025-03-01: Project initialized
- 2025-03-01: Research completed (4 dimensions)
- 2025-03-01: Requirements defined (118 requirements)

---

## Next Actions

1. ~~Trust System foundation~~ - DONE (Quick Task 1)
2. Trust System: MsgServer & QueryServer implementation
3. Trust System: AI Mediation events and Evidence storage
4. Define Cosmos SDK module structure (x/identity, x/service, x/escrow)

---

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-03-02 | Core + Plugin architecture | Core ensures basic operation, plugins extend by role |
| 2026-03-02 | Service market 3 layers | Meets needs from simple API to complex workflows |
| 2026-03-02 | Workflow as service type | Workflow is a service type, not GenieBot feature |
| 2026-03-02 | OpenFang as Agent executor | Ready-made templates and security isolation |
| 2026-03-02 | Cosmos SDK framework | Battle-tested, IBC support, modular design |
| 2026-03-02 | Trust System consolidated | MQ and dispute are tightly coupled |

---

*Last updated: 2026-03-03*
