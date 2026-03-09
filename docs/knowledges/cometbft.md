# CometBFT 知识文档

## 概述

**CometBFT**（前身为 Tendermint Core）是一个高性能、拜占庭容错的共识引擎，用于构建区块链应用。它将共识层和应用层完全解耦，开发者只需关注业务逻辑，无需从头实现共识和网络层。

```
┌─────────────────────────────────────┐
│         Application Layer           │
│        (Cosmos SDK App)             │
├─────────────────────────────────────┤
│         ABCI (ABCI++)               │  ← 应用区块链接口
├─────────────────────────────────────┤
│         CometBFT Core               │
│  ┌─────────────┬───────────────┐    │
│  │  Consensus  │   P2P Network │    │
│  │   (BFT)     │   (Mempool)   │    │
│  └─────────────┴───────────────┘    │
└─────────────────────────────────────┘
```

---

## 1. ABCI（Application Blockchain Interface）

ABCI 是 CometBFT 与应用层通信的桥梁，定义了一组标准化的消息类型。

### 核心方法

| 方法 | 触发时机 | 用途 |
|------|----------|------|
| `InitChain` | 链初始化 | 初始状态设置（如初始验证者、genesis 参数） |
| `Info` | CometBFT 启动/重连 | 获取应用最新状态（高度、AppHash） |
| `PrepareProposal` | 作为 proposer 准备区块 | 构建区块交易列表（ABCI++） |
| `ProcessProposal` | 收到区块提案 | 验证区块是否应接受 |
| `FinalizeBlock` | 区块最终确定 | 执行交易、更新状态 |
| `Commit` | 区块提交完成 | 持久化状态、返回新的 AppHash |
| `CheckTx` | 交易进入内存池 | 基础验证（签名、格式、基础费用检查） |
| `Query` | 客户端查询 | 状态查询接口 |

### ABCI++ 增强（CometBFT 0.38+）

- **PrepareProposal**: 允许应用参与区块构建，可重排序、过滤或插入交易
- **ProcessProposal**: 应用可验证区块有效性，实现自定义共识规则
- **ExtendVote/VerifyVoteExtension**: 支持预共识数据交换

### 请求/响应流程

```go
// 示例：FinalizeBlock 请求结构
type RequestFinalizeBlock struct {
    Block   *Block
    DecidedLastCommit CommitInfo
    Misbehavior      []Misbehavior
    ByzantineValidators []Evidence
}

// 响应
type ResponseFinalizeBlock struct {
    Events         []Event
    TxResults      []ResponseDeliverTx
    ValidatorUpdates []ValidatorUpdate
    ConsensusParamUpdates *ConsensusParams
    AppHash        []byte
}
```

---

## 2. 共识机制

CometBFT 实现了经典的 **BFT 权益证明**共识算法，可容忍 ≤1/3 验证者作恶。

### 验证者（Validators）

- 按质押权益权重投票
- 每个 height 轮换 proposer（加权轮询）
- 验证者集可通过 ABCI 的 `ValidatorUpdates` 动态更新

### 共识轮次（Round）

每个高度可能有多个轮次，每个轮次包含三个阶段：

```
┌────────────────────────────────────────────────────┐
│                    Height N                         │
├────────────────────────────────────────────────────┤
│  Round 0:                                          │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐          │
│  │ Propose │ → │ Prevote │ → │ Precommit│ → Commit │
│  └─────────┘   └─────────┘   └─────────┘          │
│       ↓             ↓             ↓                │
│    Proposer     2/3 Prevote   2/3 Precommit       │
│    broadcasts    (POL)          required          │
└────────────────────────────────────────────────────┘
```

### 阶段详解

1. **Propose**: Proposer 广播区块提案
2. **Prevote**: 验证者对提案投票（prevote）
3. **Precommit**: 收到 2/3+ prevotes 后发送 precommit
4. **Commit**: 收到 2/3+ precommits 后提交区块

### 锁定机制（Locked Value）

- 验证者一旦 precommit 某区块，就锁定该值
- 后续轮次必须 prevote 锁定值，除非看到更高轮次的 2/3+ prevotes（POL）
- 防止双签和分叉

### 安全性保证

| 属性 | 说明 |
|------|------|
| **安全性（Safety）** | 不会有两个区块在同一高度被提交 |
| **活跃性（Liveness）** | 只要 >2/3 验证者诚实，链就会继续出块 |
| **最终性（Finality）** | 一旦区块提交，即不可回滚（即时最终性） |

---

## 3. P2P 网络层

### 节点类型

- **Full Node**: 完整节点，同步完整区块链数据
- **Validator Node**: 验证者节点，参与共识
- **Seed Node**: 种子节点，提供地址发现服务

### 传输协议

- **TCP + Multiplex**: 默认传输层
- **WebSocket**: 支持 RPC 订阅
- **P2P Encryption**: 节点间通信加密

### 核心组件

```
┌─────────────────────────────────────────┐
│              P2P Layer                   │
├─────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌────────┐│
│  │ Switch   │  │  Peer    │  │ Addr   ││
│  │ (Router) │  │ Manager  │  │ Book   ││
│  └──────────┘  └──────────┘  └────────┘│
├─────────────────────────────────────────┤
│  Channels (Reactor-based)               │
│  - Block Channel                        │
│  - Consensus Channel                    │
│  - Mempool Channel                      │
│  - Evidence Channel                     │
└─────────────────────────────────────────┘
```

### Mempool（内存池）

- 存储待打包交易
- 支持 `CheckTx` 验证
- 交易广播（gossip）
- 配置项：
  - `size`: 最大交易数
  - `cache_size`: 缓存大小（防止重放）
  - `keep-invalid-txs-in-cache`: 是否保留无效交易

### 区块同步

- **State Sync**: 从快照快速同步（新节点）
- **Block Sync**: 从其他节点拉取历史区块
- **Fast Sync**: 结合两种方式

---

## 4. 区块结构与生命周期

### 区块结构

```
┌─────────────────────────────────────────┐
│               Block                      │
├─────────────────────────────────────────┤
│  Header (11 fields)                     │
│  ├── Version {Block, App}               │
│  ├── ChainID                            │
│  ├── Height                             │
│  ├── Time                               │
│  ├── LastBlockID                        │
│  ├── LastCommitHash                     │
│  ├── DataHash (txs merkle)              │
│  ├── ValidatorsHash                     │
│  ├── NextValidatorsHash                 │
│  ├── ConsensusHash                      │
│  ├── AppHash (应用状态根)                │
│  └── LastResultsHash                    │
├─────────────────────────────────────────┤
│  Data (Transactions)                    │
│  └── [][]byte                           │
├─────────────────────────────────────────┤
│  Evidence (Byzantine evidence)          │
│  └── []Evidence                         │
├─────────────────────────────────────────┤
│  LastCommit                              │
│  ├── BlockID                            │
│  └── []Signatures                       │
└─────────────────────────────────────────┘
```

### 区块生命周期

```
    ┌──────────────┐
    │   Mempool    │  交易进入内存池
    │   (CheckTx)  │
    └──────┬───────┘
           ↓
    ┌──────────────┐
    │   Proposed   │  Proposer 选择交易打包
    │   (Prepare)  │
    └──────┬───────┘
           ↓
    ┌──────────────┐
    │   Prevote    │  验证者投票
    │    Phase     │
    └──────┬───────┘
           ↓
    ┌──────────────┐
    │  Precommit   │  2/3+ Prevotes
    │    Phase     │
    └──────┬───────┘
           ↓
    ┌──────────────┐
    │   Commit     │  2/3+ Precommits
    │ (Finalize)   │
    └──────┬───────┘
           ↓
    ┌──────────────┐
    │   Committed  │  区块写入磁盘
    │              │  应用状态更新
    └──────────────┘
```

### 关键哈希计算

```go
// Block ID
type BlockID struct {
    Hash         []byte
    PartSetHeader PartSetHeader
}

// AppHash 更新流程
AppHash_n = Hash(AppState_n-1 + BlockTxs_n)
```

---

## 5. 与 Cosmos SDK 的关系

### 架构层次

```
┌────────────────────────────────────────────┐
│           Cosmos SDK Application            │
│  ┌──────────────────────────────────────┐  │
│  │           BaseApp                     │  │
│  │  ┌────────────────────────────────┐  │  │
│  │  │      Modules (x/auth, x/bank)  │  │  │
│  │  └────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────┐  │  │
│  │  │      Keeper / State Machine    │  │  │
│  │  └────────────────────────────────┘  │  │
│  └──────────────────────────────────────┘  │
├────────────────────────────────────────────┤
│              ABCI / ABCI++                 │
├────────────────────────────────────────────┤
│              CometBFT                      │
│  ┌──────────┬──────────┬───────────────┐  │
│  │ Consensus│   P2P    │    Mempool    │  │
│  └──────────┴──────────┴───────────────┘  │
└────────────────────────────────────────────┘
```

### BaseApp 与 ABCI 映射

| ABCI 方法 | BaseApp 处理 |
|-----------|-------------|
| `InitChain` | 初始化模块，设置初始状态 |
| `Info` | 返回最新高度和 AppHash |
| `PrepareProposal` | 调用模块的 `PrepareProposal` 钩子 |
| `ProcessProposal` | 验证区块交易 |
| `FinalizeBlock` | 执行 AnteHandler → Msg Handlers → Events |
| `Commit` | 写入状态存储，计算新 AppHash |
| `CheckTx` | AnteHandler 验证（签名、nonce、gas） |
| `Query` | 路由到模块的 Querier |

### 交易处理流程

```
Tx Submitted
     ↓
┌─────────────────┐
│    CheckTx      │  CometBFT
│  (AnteHandler)  │  验证基础有效性
└────────┬────────┘
         ↓
    [进入 Mempool]
         ↓
┌─────────────────┐
│ FinalizeBlock   │  区块执行
│  ├─ AnteHandler │  再次验证
│  ├─ MsgRouter   │  路由到模块
│  └─ Keeper      │  状态更新
└────────┬────────┘
         ↓
┌─────────────────┐
│     Commit      │  持久化
│  (AppHash)      │  状态根
└─────────────────┘
```

### ShareTokens 项目中的应用

在 ShareTokens 项目中：

```
CometBFT 负责:
- 验证者选择和轮换
- 区块共识和最终性
- 交易广播和排序
- P2P 网络通信

Cosmos SDK (BaseApp) 负责:
- 交易解析和路由
- 模块状态机逻辑 (idea, task, trust, ...)
- 状态存储 (IAVL Store)
- 事件发布和查询

ABCI 交互点:
- CheckTx: 验证交易签名、格式
- FinalizeBlock: 执行 Idea 创建、Task 分配、Trust 评分更新
- Commit: 提交新的状态根
```

---

## 6. 实用配置参考

### config.toml 关键配置

```toml
[consensus]
# 区块时间
timeout_propose = "3s"
timeout_prevote = "1s"
timeout_precommit = "1s"

[mempool]
# 内存池大小
size = 5000
cache_size = 10000
# 广播策略
broadcast = true

[p2p]
# 监听地址
laddr = "tcp://0.0.0.0:26656"
# 种子节点
seeds = ""
# 最大连接数
max_num_inbound_peers = 40
max_num_outbound_peers = 10

[instrumentation]
# Prometheus 指标
prometheus = true
prometheus_listen_addr = ":26660"
```

### 常用 RPC 端点

| 端点 | 用途 |
|------|------|
| `/status` | 节点状态信息 |
| `/block` | 获取区块 |
| `/block_results` | 区块执行结果 |
| `/tx` | 查询交易 |
| `/broadcast_tx_sync` | 同步广播交易 |
| `/abci_query` | 查询应用状态 |

---

## 参考资料

- [CometBFT 官方文档](https://docs.cometbft.com/)
- [Cosmos SDK 文档](https://docs.cosmos.network/)
- [ABCI++ 规范](https://github.com/cometbft/cometbft/blob/main/spec/abci/README.md)
- [Tendermint 论文](https://arxiv.org/abs/1807.04938)
