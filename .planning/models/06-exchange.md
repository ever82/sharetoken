# Oracle Service - 预言机服务（汇率层）

> **模块类型:** 辅助模块
> **技术栈:** TypeScript + Chainlink
> **位置:** `src/services/oracle`
> **依赖:** Chainlink 网络

---

## 概述

Oracle Service 是 ShareTokens 的辅助预言机服务，通过 Chainlink 网络获取去中心化的价格数据，为服务市场提供可靠的汇率转换。这是一个链下服务，不作为核心共识模块。

---

## 架构位置

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ShareTokens 架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  核心模块                                                                    │
│  ├── 服务市场 (11-service) ────────────────┐                               │
│  │   └── 定价需要汇率                       │                               │
│  └──────────────────────────────────────────┤                               │
│                                             ▼                               │
│  辅助模块: Oracle Service (06-exchange)                                     │
│  └── 通过 Chainlink 获取汇率数据                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 技术架构

```
Chainlink Network (外部)
       │
       ▼
Oracle Service (TypeScript)
├── price-feed.ts      # 价格订阅
├── aggregator.ts      # 数据聚合
├── signer.ts          # 消息签名
├── submitter.ts       # 链上提交
└── cache.ts           # 本地缓存
       │
       ▼
x/compute (链上模块)
```

---

## Chainlink 集成

```typescript
// Chainlink 提供的能力（直接使用，无需自行实现）

interface ChainlinkCapabilities {
  // Price Feeds: 去中心化价格数据
  priceFeeds: {
    getPrice(pair: string): Promise<PriceData>
    subscribe(pair: string, callback: PriceCallback): Promise<Subscription>
  }

  // VRF: 可验证随机数
  vrf: {
    requestRandomness(seed: bigint): Promise<bigint>
  }

  // Automation: 自动化任务
  automation: {
    registerUpkeep(task: UpkeepTask): Promise<string>
    performUpkeep(id: string): Promise<void>
  }

  // Functions: 链下计算
  functions: {
    execute(request: FunctionRequest): Promise<FunctionResult>
  }

  // CCIP: 跨链数据传输
  ccip: {
    sendMessage(destChain: string, data: bytes): Promise<TxHash>
  }
}
```

---

## 汇率快照

```typescript
interface ExchangeRate {
  id: string
  pair: {
    base: 'STT'
    quote: string  // 'USD' | 'ETH' | 'BTC' | 'CNY'
  }
  rate: number  // 1 STT = ? quote

  // Chainlink 数据
  chainlinkRoundId: bigint
  chainlinkAnswer: bigint
  startedAt: Timestamp
  updatedAt: Timestamp

  // 置信度
  confidence: number  // 0-1

  // 有效性
  validFrom: Timestamp
  validUntil: Timestamp
}
```

---

## Token 价格映射

> 标准化各 LLM 官方价格到 STT

```typescript
interface TokenPriceMapping {
  provider: 'openai' | 'anthropic' | 'google' | 'azure'
  model: string

  // 官方定价（USD）
  officialPrice: {
    inputPerToken: number   // USD per 1K tokens
    outputPerToken: number
  }

  // STT 定价
  sttPrice: {
    inputPerToken: TokenAmount  // micro-STT per token
    outputPerToken: TokenAmount
  }

  // 汇率快照
  exchangeRateId: string

  updatedAt: Timestamp
}

// 示例
const exampleMapping: TokenPriceMapping = {
  provider: 'openai',
  model: 'gpt-4-turbo',
  officialPrice: {
    inputPerToken: 0.01,   // $0.01 / 1K tokens
    outputPerToken: 0.03   // $0.03 / 1K tokens
  },
  sttPrice: {
    inputPerToken: { amount: 100n, symbol: 'STT' },   // 100 micro-STT
    outputPerToken: { amount: 300n, symbol: 'STT' }   // 300 micro-STT
  },
  exchangeRateId: 'rate-001',
  updatedAt: Date.now()
}
```

---

## 价格订阅服务

```typescript
// src/services/oracle/price-feed.ts

import { ChainlinkPriceFeed } from '@chainlink/contracts'

export class PriceFeedService {
  private feeds: Map<string, ChainlinkPriceFeed>
  private cache: Map<string, CachedPrice>

  async initialize(pairs: string[]): Promise<void> {
    for (const pair of pairs) {
      const feed = await ChainlinkPriceFeed.create(pair)
      this.feeds.set(pair, feed)

      // 订阅价格更新
      feed.on('AnswerUpdated', (answer, roundId, timestamp) => {
        this.handlePriceUpdate(pair, answer, roundId, timestamp)
      })
    }
  }

  async getPrice(pair: string): Promise<ExchangeRate> {
    // 先检查缓存
    const cached = this.cache.get(pair)
    if (cached && !this.isStale(cached)) {
      return cached.rate
    }

    // 从 Chainlink 获取最新价格
    const feed = this.feeds.get(pair)
    const roundData = await feed.getLatestRoundData()

    return this.toExchangeRate(pair, roundData)
  }

  async subscribe(pair: string, callback: (rate: ExchangeRate) => void): Promise<Subscription> {
    const feed = this.feeds.get(pair)
    if (!feed) throw new Error(`Unknown pair: ${pair}`)

    const handler = (answer, roundId, timestamp) => {
      const rate = this.toExchangeRate(pair, { answer, roundId, timestamp })
      callback(rate)
    }

    feed.on('AnswerUpdated', handler)

    return {
      unsubscribe: () => feed.off('AnswerUpdated', handler)
    }
  }

  private isStale(cached: CachedPrice): boolean {
    const heartbeat = 60 * 60 * 1000 // 1 hour
    return Date.now() - cached.timestamp > heartbeat
  }
}
```

---

## 链上提交器

```typescript
// src/services/oracle/submitter.ts

import { SigningStargateClient } from '@cosmjs/stargate'

export class OracleSubmitter {
  private client: SigningStargateClient
  private signer: OfflineSigner
  private oracleAddress: string

  async submitPriceUpdate(rate: ExchangeRate): Promise<TxResult> {
    const msg = {
      typeUrl: '/sharetokens.oracle.MsgUpdatePrice',
      value: {
        oracle: this.oracleAddress,
        pair: `${rate.pair.base}/${rate.pair.quote}`,
        rate: Math.floor(rate.rate * 1e6), // 使用整数精度
        roundId: rate.chainlinkRoundId.toString(),
        timestamp: BigInt(rate.updatedAt),
        signature: await this.signRate(rate)
      }
    }

    const fee = {
      amount: coins(5000, 'stake'),
      gas: '200000'
    }

    return this.client.signAndBroadcast(
      this.oracleAddress,
      [msg],
      fee
    )
  }

  private async signRate(rate: ExchangeRate): Promise<Uint8Array> {
    const data = JSON.stringify({
      pair: rate.pair,
      rate: rate.rate,
      roundId: rate.chainlinkRoundId.toString(),
      timestamp: rate.updatedAt
    })

    const accounts = await this.signer.getAccounts()
    const signature = await this.signer.signMessage(
      accounts[0].address,
      data
    )

    return signature
  }
}
```

---

## 定价计算器

```typescript
// src/services/oracle/pricing.ts

export class PricingCalculator {
  private priceFeed: PriceFeedService
  private mappings: Map<string, TokenPriceMapping>

  async calculateSttPrice(
    provider: string,
    model: string,
    tokensUsed: { input: number; output: number }
  ): Promise<TokenAmount> {
    // 获取价格映射
    const key = `${provider}/${model}`
    const mapping = this.mappings.get(key)
    if (!mapping) throw new Error(`Unknown model: ${key}`)

    // 计算基础价格 (USD)
    const inputCost = (tokensUsed.input / 1000) * mapping.officialPrice.inputPerToken
    const outputCost = (tokensUsed.output / 1000) * mapping.officialPrice.outputPerToken
    const totalUsd = inputCost + outputCost

    // 获取 STT/USD 汇率
    const rate = await this.priceFeed.getPrice('STT/USD')

    // 转换为 STT
    const sttAmount = totalUsd / rate.rate

    return {
      amount: BigInt(Math.floor(sttAmount * 1e6)), // micro-STT
      symbol: 'STT'
    }
  }

  async updateMappings(): Promise<void> {
    // 从链上或配置中获取最新价格映射
    // 根据官方 API 更新价格
  }
}
```

---

## 服务配置

```typescript
// src/services/oracle/config.ts

export interface OracleConfig {
  // Chainlink 配置
  chainlink: {
    network: 'mainnet' | 'testnet'
    rpcUrl: string
    feedAddresses: Record<string, string>  // pair -> address
  }

  // 链上提交配置
  submission: {
    enabled: boolean
    interval: number  // 提交间隔 (ms)
    threshold: number // 价格变动阈值 (百分比)
  }

  // 缓存配置
  cache: {
    ttl: number       // 缓存过期时间 (ms)
    maxSize: number
  }

  // 监控配置
  monitoring: {
    alertThreshold: number  // 价格偏离阈值
    heartbeatTimeout: number
  }
}
```

---

## 与服务市场交互

```
Oracle Service                          服务市场 (11-service)
      │                                      │
      │  1. 服务市场请求价格转换              │
      │◄─────────────────────────────────────│
      │                                      │
      │  2. 从 Chainlink 获取汇率            │
      │                                      │
      │  3. 计算并返回 STT 价格              │
      │─────────────────────────────────────►│
      │                                      │
      │  4. 服务市场完成定价                  │
      │                                      │
```

---

## Chainlink 能力映射

| 功能 | Chainlink 模块 | 本服务使用方式 |
|------|---------------|----------------|
| 价格数据 | Price Feeds | 订阅并转发到服务市场 |
| 数据签名 | On-chain Validation | 验证价格真实性 |
| 心跳更新 | Automation | 定期更新价格 |
| 跨链 | CCIP | 多链价格同步 |
| 随机数 | VRF | 争议系统随机抽选 |

---

## 依赖关系

```
Oracle Service (06-exchange)
    │
    ├── 外部依赖
    │   └── Chainlink Network
    │
    └── 被依赖
        └── 服务市场 (11-service) - 汇率转换
```

---

[上一章：x/compute](./05-compute.md) | [返回索引](./00-index.md) | [下一章：x/idea →](./07-idea.md)
