# ShareTokens 并行开发计划

> 基于 TDD + Mock First 策略，最大化并行开发效率。

---

## 开发策略

1. **接口契约优先** - 所有模块先定义接口，放在 `docs/contracts/`
2. **Mock First** - 依赖方用 Mock 先行开发，完成后联调替换
3. **TDD** - 每个模块自测 + Consumer Tests 验证接口

---

## Wave 0 - 核心基础设施（3 并行）

**无相互依赖，可立即开工**

| Agent | 任务 | 交付物 |
|-------|------|--------|
| Agent-1 | ACH-DEV-001 Dev Infrastructure | proto/, CI/CD, Makefile, Lint |
| Agent-2 | ACH-DEV-002 Blockchain Network | app/, cmd/, config/, 启动脚本 |
| Agent-3 | ACH-DEV-003 Wallet & Token | x/bank/, x/token/, 钱包集成 |

---

## Wave 1 - 核心功能（5 并行 + 2 后续）

**依赖 Wave 0 的接口契约**

### 第一批（可立即并行）

| Agent | 任务 | Mock 依赖 | 提供接口 |
|-------|------|-----------|----------|
| Agent-4 | ACH-DEV-004 Identity | Trust, Escrow | Identity |
| Agent-5 | ACH-DEV-005 Escrow | Trust, Market | Escrow |
| Agent-6 | ACH-DEV-006 Oracle | Chainlink | Oracle |
| Agent-7 | ACH-DEV-007 MQ Scoring | Escrow, Market | Trust (部分) |
| Agent-8 | ACH-DEV-010 Testnet | Core modules | - |

### 第二批（等待依赖完成）

| Agent | 任务 | 依赖 |
|-------|------|------|
| Agent-9 | ACH-DEV-008 Dispute | ACH-007 (MQ Scoring) |
| Agent-10 | ACH-DEV-009 Service Market | ACH-004, 005, 006 |

---

## Wave 2 - 插件与扩展（8 并行）

**依赖 Wave 1 的 Market 模块**

| Agent | 任务 | Mock 依赖 |
|-------|------|-----------|
| Agent-11 | ACH-DEV-011 LLM Plugin | Market, Escrow |
| Agent-12 | ACH-DEV-012 Agent Plugin | Market, Escrow |
| Agent-13 | ACH-DEV-013 Workflow Plugin | Market, Escrow |
| Agent-14 | ACH-DEV-014 GenieBot UI | Market, Wallet |
| Agent-15 | ACH-DEV-015 Task Market | Escrow |
| Agent-16 | ACH-DEV-016 Idea/Crowdfunding | Escrow |
| Agent-17 | ACH-DEV-017 Performance | Wave 1 完成 |
| Agent-18 | ACH-DEV-018 Observability | Core modules |

---

## 联调计划

```
Wave 0 完成后:
├── Identity + Wallet (注册流程)
├── Escrow + Wallet (托管流程)
└── Network + All (链上交互)

Wave 1 完成后:
├── Identity + Escrow + Trust (限额+争议)
├── Oracle + Market (定价)
└── Market + Plugins (服务调用)

Wave 2 完成后:
└── 端到端集成测试
```

---

## 接口契约清单

启动 Wave 1 前需定义：

| 契约 | 文件 | 使用方 |
|------|------|--------|
| Identity | `docs/contracts/identity.go` | Market, Escrow, Trust |
| Escrow | `docs/contracts/escrow.go` | Market, Trust, Plugins |
| Oracle | `docs/contracts/oracle.go` | Market |
| Trust | `docs/contracts/trust.go` | Escrow, Identity, Market |
| Market | `docs/contracts/market.go` | Plugins, GenieBot |

---

## 立即启动

```bash
# 第一波 3 个任务可现在开工
git worktree add ../sharetokens-agent1 -b feature/dev-infra
git worktree add ../sharetokens-agent2 -b feature/blockchain
git worktree add ../sharetokens-agent3 -b feature/wallet
```

---

## 总结

| 波次 | 并行数 | 任务 |
|------|--------|------|
| Wave 0 | 3 | P0 基础设施 |
| Wave 1 | 5+2 | P1 核心功能 |
| Wave 2 | 8 | P2 插件扩展 |

**第一波可同时开工: 3 个任务** (ACH-001, 002, 003)
