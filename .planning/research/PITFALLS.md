# Pitfalls Research

**Domain:** ShareTokens - Decentralized AI Service Marketplace
**Researched:** 2026-03-02
**Confidence:** MEDIUM-HIGH

---

## Architecture-Specific Pitfalls

### Pitfall 0: Plugin Architecture Complexity

**What goes wrong:**
The Core + Plugin architecture introduces complexity in versioning, dependencies, and plugin compatibility. Plugins may become incompatible with core updates, or core changes may break plugins.

**Why it happens:**
- Loose coupling requires careful interface design
- Plugin API versioning is often overlooked
- Core modules may change without considering plugin impact

**How to avoid:**
- Define stable Core APIs with semantic versioning
- Use capability negotiation between Core and Plugins
- Implement plugin isolation to prevent cascading failures
- Maintain plugin compatibility test suite

**Phase to address:** Phase 1 (Foundation) - API contracts must be defined early

---

## Core Module Pitfalls

### Pitfall 1: Sybil Attack Vulnerability (Identity Module)

**What goes wrong:**
Malicious actors create multiple fake identities to manipulate reputation, distort rewards, or gain disproportionate influence. This breaks the Trust System and marketplace fairness.

**Why it happens:**
- No cost to create new blockchain addresses
- Pseudonymous by design
- Token-based governance without identity binding

**Consequences:**
- Attackers drain reward pools
- MQ system becomes meaningless
- Legitimate providers get squeezed out

**How to avoid:**
- Real-name identity verification (KYC) for reward eligibility
- Stake-based participation (economic barrier)
- Identity hash stored on-chain (not raw identity data)
- ZK-DID for privacy-preserving verification

**Warning signs:**
- Rapid spike in new account registrations
- Clusters of accounts with similar patterns
- Multiple accounts from same IP ranges

**Phase to address:** Phase 1 (Core Modules) - Identity system before any rewards

---

### Pitfall 2: Escrow Fund Lock (Escrow Module)

**What goes wrong:**
Funds locked in escrow become inaccessible due to bugs, missing conditions, or dispute resolution failures. Users lose trust and funds.

**Why it happens:**
- Incomplete condition logic
- Missing timeout handling
- Dispute resolution deadlocks

**Consequences:**
- User funds permanently locked
- Loss of platform trust
- Regulatory issues

**How to avoid:**
- Always implement timeout auto-release
- Clear escalation path for disputes
- Emergency governance intervention capability
- Comprehensive test coverage for all escrow states

**Warning signs:**
- Escrow accounts accumulating funds
- Long-running unresolved disputes
- User complaints about fund access

**Phase to address:** Phase 2 (Core Business) - Escrow must be bulletproof

---

### Pitfall 3: MQ Gaming (Trust System)

**What goes wrong:**
Users find ways to artificially inflate their MQ scores or manipulate the zero-sum redistribution in their favor.

**Why it happens:**
- Zero-sum mechanism is novel and untested at scale
- Collusion between parties
- Sybil accounts transferring reputation

**Consequences:**
- Trust system becomes meaningless
- High MQ users dominate jury selection
- Platform fairness collapses

**How to avoid:**
- Minimum transaction value for MQ changes
- Time-based decay prevents accumulation
- Behavioral analysis for collusion detection
- Game-theoretic modeling before implementation

**Warning signs:**
- Rapid MQ score changes
- Circular transaction patterns
- Coordinated voting in disputes

**Phase to address:** Phase 2 (Core Business) - Validate model before deployment

---

### Pitfall 4: Jury Collusion (Trust System)

**What goes wrong:**
Jurors collude to always vote in favor of one party, or accept bribes to influence outcomes. This undermines the entire dispute resolution system.

**Why it happens:**
- Juror identities may become known
- Economic incentives for bribery
- No penalty for dishonest voting

**How to avoid:**
- Anonymous juror selection where possible
- MQ-weighted selection (higher MQ = more to lose)
- Stake slashing for proven collusion
- Appeal process with fresh jury

**Warning signs:**
- Same jurors repeatedly selected
- Predictable voting patterns
- Correlation between juror and outcome

**Phase to address:** Phase 2 (Core Business) - Anti-collusion from day one

---

### Pitfall 5: P2P Network Partition (P2P Module)

**What goes wrong:**
Network splits into isolated partitions that cannot communicate. Each partition may have different chain states, causing confusion when reconnected.

**Why it happens:**
- NAT traversal failures
- Bootstrap node failures
- Geographic network issues

**Consequences:**
- Double-spending attempts
- Inconsistent state across partitions
- Service unavailability

**How to avoid:**
- Multiple bootstrap nodes in different regions
- Aggressive peer discovery
- Partition detection and recovery protocols
- CometBFT's built-in consensus handles chain consistency

**Warning signs:**
- Declining peer count
- Stalled block production
- Unable to reach known peers

**Phase to address:** Phase 1 (Core Modules) - P2P stability is foundational

---

## Service Provider Plugin Pitfalls

### Pitfall 6: API Key Theft (LLM Hosting Plugin)

**What goes wrong:**
LLM API keys stored by providers are stolen or leaked, causing financial loss and security breaches.

**Why it happens:**
- Keys stored in plaintext
- Keys transmitted to workers
- Poor access control

**Consequences:**
- Financial loss for key owners
- Platform liability
- Loss of provider trust

**How to avoid:**
- NEVER store keys on-chain
- Encrypt keys at rest (provider's local storage)
- Proxy layer prevents key exposure to workers
- Real-time usage monitoring with auto-shutoff
- Key rotation support

**Warning signs:**
- API usage spikes
- Keys used from unexpected locations
- Cost anomalies

**Phase to address:** Phase 3 (Provider Plugins) - Security before launch

---

### Pitfall 7: Agent Escape (OpenFang Plugin)

**What goes wrong:**
Malicious or buggy Agent escapes its sandbox and accesses host system resources or data.

**Why it happens:**
- WASM sandbox vulnerabilities
- Misconfigured permissions
- Resource limit bypass

**Consequences:**
- Host system compromise
- Data leakage
- Service disruption

**How to avoid:**
- Leverage OpenFang's 16-layer security
- WASM sandbox with capability-based security
- Strict resource limits (CPU, memory, network)
- Regular security audits
- Isolated execution environment

**Warning signs:**
- Unexpected resource usage
- Failed capability requests
- Unusual network connections

**Phase to address:** Phase 3 (Provider Plugins) - Security hardening required

---

### Pitfall 8: Workflow State Corruption (Workflow Plugin)

**What goes wrong:**
Long-running workflows lose state, fail to recover from errors, or produce inconsistent results.

**Why it happens:**
- State not persisted properly
- No error recovery logic
- Partial failures not handled

**Consequences:**
- Lost work and wasted resources
- User frustration
- Payment disputes

**How to avoid:**
- Checkpoint workflow state regularly
- Implement idempotent operations
- Clear error recovery paths
- Timeout and retry mechanisms
- Compensating transactions for rollback

**Warning signs:**
- Workflows stuck in intermediate states
- Inconsistent outputs for same inputs
- Frequent manual interventions

**Phase to address:** Phase 3 (Provider Plugins) - State management is critical

---

## User Plugin Pitfalls

### Pitfall 9: GenieBot Misinterpretation (User Plugin)

**What goes wrong:**
GenieBot misinterprets user intent and recommends wrong services, leading to wasted money and frustration.

**Why it happens:**
- Natural language ambiguity
- Incomplete context understanding
- Limited service matching logic

**Consequences:**
- Wasted escrow funds
- User frustration
- Disputes over wrong service execution

**How to avoid:**
- Confirmation before service invocation
- Clear cost estimation
- Service matching transparency
- Easy cancellation/refund path
- Continuous improvement based on feedback

**Warning signs:**
- High cancellation rate
- Frequent service changes after recommendation
- User complaints about relevance

**Phase to address:** Phase 4 (User Plugins) - UX validation required

---

## Integration Pitfalls

### Pitfall 10: Plugin-Core API Mismatch

**What goes wrong:**
Plugins expect certain Core APIs that have changed or don't exist, causing runtime failures.

**Why it happens:**
- Independent development of Core and Plugins
- API versioning not enforced
- Breaking changes without migration path

**How to avoid:**
- Semantic versioning for Core APIs
- Capability negotiation at plugin load
- Backward compatibility requirements
- Comprehensive integration tests

**Phase to address:** Phase 1 (Foundation) - Define API contracts

---

### Pitfall 11: Cross-Plugin Data Leakage

**What goes wrong:**
One plugin accidentally accesses or modifies data belonging to another plugin, causing security issues or data corruption.

**Why it happens:**
- Shared memory space
- Insufficient isolation
- Global state pollution

**How to avoid:**
- Strict plugin isolation
- Well-defined inter-plugin communication
- No shared mutable state
- Security sandboxing

**Phase to address:** Phase 3 (Provider Plugins) - Isolation architecture

---

## Token Economics Pitfalls

### Pitfall 12: Unsustainable Emissions

**What goes wrong:**
Token emissions exceed actual network revenue, causing token value collapse and provider exodus.

**Why it happens:**
- Over-optimistic growth projections
- Speculation-driven design
- No relationship between emissions and usage

**Consequences:**
- Token price collapse
- Provider flight
- Death spiral

**How to avoid:**
- Design around real utility (service payment)
- Match emissions to network usage
- Treasury reserves for stability
- Stake requirements for participation

**Phase to address:** Phase 1 (Foundation) - Economic model validation

---

## Reliability Pitfalls

### Pitfall 13: Decentralized Reliability Gap

**What goes wrong:**
Decentralized network reliability (3.2% failure rate) is significantly worse than centralized alternatives (0.07%), making the platform unusable for serious work.

**Why it happens:**
- Consumer-grade hardware
- Node churn
- No enforceable SLAs

**Consequences:**
- Users abandon platform
- Cannot compete with centralized providers
- Enterprise customers won't adopt

**How to avoid:**
- Redundancy (multiple providers per request)
- Quality tiers (enterprise-grade vs consumer-grade)
- MQ-weighted routing
- SLA simulation via escrow refunds
- Graceful degradation

**Phase to address:** Phase 2 (Core Business) - Reliability targets from start

---

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Skip identity verification | Easy onboarding | Sybil attacks | Never in production |
| Trust-based escrow | Simple implementation | Fraud, disputes | Never |
| Single-provider routing | Simpler code | Single point of failure | MVP only |
| No plugin isolation | Easy integration | Security issues | Never |
| Fixed MQ scoring | Simple logic | Gaming, manipulation | Never |
| Centralized discovery | Faster lookups | Censorship, failure | Early testing only |

---

## Security Checklist

### Core Modules Security

- [ ] **Identity:** Sybil resistance enforced
- [ ] **Wallet:** Key management secure
- [ ] **Escrow:** Timeout release implemented
- [ ] **Trust System:** Gaming prevention active
- [ ] **Trust System:** Collusion detection working
- [ ] **P2P:** Peer authentication required

### Plugin Security

- [ ] **LLM Hosting:** Keys encrypted, never on-chain
- [ ] **Agent Executor:** WASM sandbox verified
- [ ] **Workflow:** State isolation confirmed
- [ ] **GenieBot:** Input sanitization active

---

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification Method |
|---------|------------------|---------------------|
| Plugin Architecture | Phase 1 | API compatibility tests |
| Sybil Attack | Phase 1 | Simulated attack testing |
| Escrow Lock | Phase 2 | State transition tests |
| MQ Gaming | Phase 2 | Game-theoretic modeling |
| Jury Collusion | Phase 2 | Anonymity verification |
| P2P Partition | Phase 1 | Network partition testing |
| API Key Theft | Phase 3 | Penetration testing |
| Agent Escape | Phase 3 | Sandbox escape testing |
| Workflow Corruption | Phase 3 | State recovery testing |
| GenieBot Misinterpretation | Phase 4 | User acceptance testing |
| API Mismatch | Phase 1 | Integration tests |
| Cross-Plugin Leakage | Phase 3 | Isolation verification |
| Unsustainable Emissions | Phase 1 | Economic modeling |
| Reliability Gap | Phase 2 | Load testing |

---

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Sybil attack detected | VERY HIGH | Freeze rewards, audit accounts, implement verification |
| Escrow funds locked | HIGH | Emergency governance release, bug fix |
| MQ gaming detected | HIGH | Reset affected scores, implement anti-gaming |
| Jury collusion proven | MEDIUM | Revote with new jury, slash colluders |
| P2P partition | MEDIUM | Wait for reconnection, consensus handles state |
| API key leaked | MEDIUM | Rotate affected keys, audit usage |
| Agent escaped sandbox | HIGH | Kill agent, audit system, patch vulnerability |
| Workflow corrupted | MEDIUM | Restore from checkpoint, compensate users |
| Plugin incompatible | LOW | Update plugin or Core API version |

---

## Sources

- [DePIN Analysis - ChainCatcher](https://www.chaincatcher.com/article/2178002) - DePIN challenges
- [Decentralized Cloud Compute Analysis](https://m.btcbaike.com/choice/2628r3.html) - Reliability data
- [Sybil Attack Prevention Research](https://link.springer.com/article/10.1007/s44443-025-00080-9) - Academic analysis
- [OpenFang Security Documentation](https://openfang.sh/) - Agent security
- [Cosmos SDK Security](https://docs.cosmos.network/) - Blockchain security patterns
- [Kleros Documentation](https://kleros.io/) - Dispute resolution patterns
- [OpenAI API Key Security](https://help.openai.com/) - Key security best practices

---
*Pitfalls research for: ShareTokens Decentralized AI Service Marketplace*
*Updated: 2026-03-02*
*Architecture Version: Core + Plugin*
