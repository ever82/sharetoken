# 四、共识层

> 基于 Cosmos SDK + CometBFT 的共识机制

---

> **重要说明：共识机制由 CometBFT（原 Tendermint）内置提供**
>
> 本项目使用 Cosmos SDK 框架，共识相关的核心功能（出块、验证、投票、状态同步等）
> 全部由 CometBFT 内置实现，**无需自行开发**。
>
> 我们只需要定义业务相关的交易类型和状态转换逻辑。

---

## 技术栈

```
共识框架：Cosmos SDK + CometBFT (原 Tendermint)
├── BFT 共识：CometBFT Core（内置，无需开发）
├── 状态机：Cosmos SDK 模块（业务逻辑）
├── ABCI++：应用层接口
└── IBC：跨链通信
```

---

## 4.1 交易

```typescript
type Transaction = {
  hash: Hash                        // 交易哈希

  // 基础字段
  nonce: Nonce                      // 发送方序列号
  chainId: number                   // 链 ID
  from: Address                     // 发送方地址
  to: Address                       // 接收方地址

  // 交易内容
  value: TokenAmount                // 转账金额
  data?: string                     // 附加数据
  gasLimit: bigint                  // Gas 限制
  gasPrice: TokenAmount             // Gas 价格

  // 签名
  signature: Signature
  publicKey: string

  timestamp: Timestamp
}
```

---

## 4.2 业务交易类型

```typescript
type BusinessTxType =
  // 代币交易
  | 'transfer'                      // 转账
  | 'stake'                         // 质押
  | 'unstake'                       // 解质押

  // 算力交易
  | 'api_key_create'                // 创建 API Key
  | 'compute_request'               // 算力请求
  | 'escrow_create'                 // 创建托管
  | 'escrow_release'                // 释放托管
  | 'escrow_dispute'                // 争议托管

  // 想法系统
  | 'idea_create'                   // 创建想法
  | 'idea_support'                  // 支持想法
  | 'contribution_submit'           // 提交贡献

  // 任务市场
  | 'task_create'                   // 创建任务
  | 'task_apply'                    // 申请任务
  | 'task_assign'                   // 分配任务
  | 'milestone_submit'              // 提交里程碑
  | 'milestone_approve'             // 批准里程碑

  // 争议系统
  | 'dispute_create'                // 创建争议
  | 'dispute_vote'                  // 争议投票

  // 身份验证
  | 'identity_register'             // 注册身份
  | 'identity_revoke'               // 撤销身份

  // 服务市场
  | 'service_register'              // 注册服务
  | 'service_call'                  // 调用服务

type BusinessTx = {
  base: Transaction
  type: BusinessTxType
  payload: object                   // 根据类型不同
}
```

---

## 4.3 交易状态

```typescript
type TxStatus =
  | 'pending'                       // 待处理
  | 'included'                      // 已入块
  | 'confirmed'                     // 已确认
  | 'failed'                        // 执行失败
  | 'rejected'                      // 被拒绝

type TxReceipt = {
  txHash: Hash
  status: TxStatus
  blockHash?: Hash
  blockNumber?: bigint
  gasUsed: bigint
  logs: TxLog[]
  timestamp: Timestamp
}

type TxLog = {
  address: Address
  topics: Hash[]
  data: string
  logIndex: number
}
```

---

## 4.4 区块

```typescript
type Block = {
  header: BlockHeader
  transactions: Transaction[]
  signatures: BlockSignature[]      // Tendermint 验证者签名
}

type BlockHeader = {
  hash: Hash
  number: bigint
  parentHash: Hash

  // 状态根
  stateRoot: Hash                   // 状态 Merkle 根
  txRoot: Hash                      // 交易 Merkle 根

  // 时间
  timestamp: Timestamp
  proposer: Address

  // Gas
  gasUsed: bigint
  gasLimit: bigint
}
```

---

## 4.5 Epoch

```typescript
type Epoch = {
  number: number
  startBlock: bigint
  endBlock: bigint

  // 验证者集合
  validators: ValidatorInfo[]

  // 统计
  stats: {
    totalBlocks: number
    totalTransactions: number
    totalGasUsed: bigint
    avgBlockTime: number
  }

  startTimestamp: Timestamp
  endTimestamp?: Timestamp
  status: 'active' | 'completed' | 'finalized'
}
```

---

## 4.6 验证者

```typescript
type ValidatorInfo = {
  address: Address
  publicKey: string

  // 质押
  stake: TokenAmount
  delegatedStake: TokenAmount

  // 状态
  status: 'active' | 'inactive' | 'jailed' | 'unbonding'

  // 性能
  performance: {
    blocksProduced: number
    blocksSigned: number
    missedBlocks: number
    uptime: number                   // 0-100
  }

  joinedAt: Timestamp
  lastActiveAt: Timestamp
}
```

---

## 4.7 Gas 配置建议

```typescript
type GasConfig = {
  // 基础交易
  transfer: 21_000
  stake: 50_000
  unstake: 50_000

  // 算力交易
  compute_request: 100_000
  escrow_create: 80_000
  escrow_release: 50_000

  // 想法/任务
  idea_create: 100_000
  task_create: 150_000
  milestone_submit: 80_000

  // 争议（复杂计算）
  dispute_create: 200_000
  dispute_vote: 100_000
}
```

---

## Tendermint 提供的能力（无需自行实现）

| 功能 | Tendermint 模块 | 说明 |
|------|----------------|------|
| 共识 | Tendermint Core | BFT 共识算法 |
| 出块 | Block Protocol | 确定性出块 |
| 验证者 | Validator Set | 验证者管理 |
| 投票 | Vote Protocol | Prevote/Precommit |
| 锁定 | Lock Mechanism | 防双花 |
| 状态同步 | State Sync | 快速同步 |
| 轻客户端 | Light Client | SPV 验证 |

---

[← 上一章：基础类型](./01-base.md) | [返回索引](./00-index.md) | [下一章：算力层 →](./05-compute.md)
