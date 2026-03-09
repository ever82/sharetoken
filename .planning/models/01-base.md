# 基础类型

> **模块类型:** 核心基础
> **被依赖:** 所有核心模块和插件
> **位置:** `src/types/`

---

## 概述

基础类型是 ShareTokens 所有模块共享的类型定义，包括原始类型、Token金额、加密数据、网络端点等。作为整个系统的类型基础，被所有核心模块和插件依赖。

---

## 架构位置

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ShareTokens 架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  核心模块                                                                    │
│  ├── P2P通信 ──────────────────┐                                           │
│  ├── 身份账号 ─────────────────┤                                           │
│  ├── 钱包 ─────────────────────┤                                           │
│  ├── 服务市场 ─────────────────┤  都依赖基础类型 (01-base)                  │
│  ├── 托管支付 ─────────────────┤                                           │
│  └── Trust System ─────────────────┘                                       │
│                                                                             │
│  可选插件                                                                    │
│  ├── LLM API Key托管 ──────────┐                                           │
│  ├── Agent执行器 ──────────────┤                                           │
│  ├── Workflow执行器 ───────────┤  都依赖基础类型 (01-base)                  │
│  └── GenieBot界面 ─────────────────┘                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 原始类型

```typescript
type ID = string                    // 唯一标识符，UUID 或哈希
type Address = string               // 钱包地址，cosmos1... 格式
type Hash = string                  // 32 字节哈希，hex 编码
type Signature = string             // 签名，hex 编码
type Timestamp = number             // Unix 时间戳 (毫秒)
type Duration = number              // 时长 (秒)
type Nonce = bigint                 // 序列号，防重放
```

---

## Token 金额

```typescript
type TokenAmount = {
  amount: bigint                    // 最小单位数量 (1 STT = 1_000_000 micro-STT)
  symbol: 'STT'                     // 代币符号
}

// 币种定义
type Denom = 'stt' | 'ustt'         // stt = 主币, ustt = micro-STT

// 金额工具函数
function parseCoin(amount: string): TokenAmount
function formatCoin(coin: TokenAmount): string
function addCoins(a: TokenAmount, b: TokenAmount): TokenAmount
function subtractCoins(a: TokenAmount, b: TokenAmount): TokenAmount
```

---

## 加密数据

```typescript
type EncryptedData = {
  ciphertext: string                // Base64 编码的密文
  nonce: string                     // 加密随机数
  algorithm: 'x25519-xsalsa20-poly1305' | 'aes-256-gcm'
  version: number                   // 算法版本，支持升级
}

// 加密工具
interface CryptoService {
  encrypt(data: Uint8Array, publicKey: string): EncryptedData
  decrypt(encrypted: EncryptedData, privateKey: string): Uint8Array
  sign(data: Uint8Array, privateKey: string): Signature
  verify(data: Uint8Array, signature: Signature, publicKey: string): boolean
}
```

---

## 网络端点

```typescript
type Endpoint = {
  type: 'ipv4' | 'ipv6' | 'dns'
  host: string
  port: number
  // NAT 穿透信息
  natType?: 'none' | 'full_cone' | 'restricted' | 'port_restricted' | 'symmetric'
  publicEndpoint?: Endpoint         // 公网端点
  relayEndpoint?: Endpoint          // 中继端点 (无法直连时使用)
}
```

---

## 状态根

```typescript
type StateRoot = {
  root: Hash                        // Merkle 根
  version: bigint                   // 状态版本
  timestamp: Timestamp
}
```

---

## Merkle 证明

```typescript
type MerkleProof = {
  root: Hash                        // Merkle 根
  leaf: Hash                        // 叶子节点哈希
  proof: Hash[]                     // 兄弟节点哈希路径
  index: number                     // 叶子节点索引
}

// 验证函数
function verifyProof(proof: MerkleProof): boolean
```

---

## 节点类型

```typescript
type NodeType = 'light' | 'full' | 'archive' | 'service'

type NodeConfig = {
  type: NodeType

  // 存储配置
  storage: {
    // 轻节点
    lightStorage?: {
      ownAccounts: boolean
      blockHeaders: number
      cacheSize: number
    }

    // 全节点
    fullStorage?: {
      fullState: boolean
      fullHistory: boolean
      pruneAfter: number
    }

    // 归档节点
    archiveStorage?: {
      fullHistory: boolean
      indexAll: boolean
    }
  }

  // 服务配置（服务节点）
  services?: {
    types: ServiceType[]
    maxConcurrent: number
  }
}
```

---

## 服务类型定义

```typescript
// 服务层级
type ServiceLevel = 1 | 2 | 3

// Level 1: LLM API 服务
// Level 2: Agent 服务
// Level 3: Workflow 服务

type ServiceType =
  // Level 1: LLM 服务
  | 'llm_chat'           // LLM 对话
  | 'llm_completion'     // LLM 补全
  | 'llm_embedding'      // LLM 嵌入

  // Level 2: Agent 服务
  | 'agent_coder'        // 编程 Agent
  | 'agent_researcher'   // 研究 Agent
  | 'agent_writer'       // 写作 Agent
  | 'agent_geniebot'     // GenieBot Agent (自定义)

  // Level 3: Workflow 服务
  | 'workflow_software'  // 软件开发流程
  | 'workflow_content'   // 内容创作流程
  | 'workflow_business'  // 商业策划流程
  | 'workflow_service'   // 生活服务流程
```

---

## OpenFang 集成类型

> 这些类型用于与 OpenFang Agent OS 集成

```typescript
// OpenFang Agent 模板类型
type AgentType =
  | 'coder'           // 代码编写
  | 'researcher'      // 研究分析
  | 'writer'          // 内容写作
  | 'architect'       // 架构设计
  | 'debugger'        // 调试修复
  | 'analyst'         // 数据分析
  | 'geniebot'        // GenieBot - 自定义主Agent

// OpenFang Hand 类型
type HandType =
  | 'collector'       // 数据收集
  | 'clip'            // 视频剪辑
  | 'lead'            // 销售线索
  | 'content'         // 内容创作
  | 'trade'           // 交易监控
  | 'browser'         // 浏览器自动化
  | 'twitter'         // 社媒管理

// OpenFang 通信频道
type ChannelType =
  | 'cli'             // 命令行
  | 'telegram'        // Telegram
  | 'discord'         // Discord
  | 'slack'           // Slack
  | 'web'             // Web界面（GenieBot）

// Agent 配置
type AgentConfig = {
  id: ID
  type: AgentType
  provider: string    // LLM Provider (openai, anthropic, openrouter)
  model: string       // 模型名称
  channels: ChannelType[]
  hands?: HandType[]  // 关联的Hands
  securityPolicy: {
    sandbox: boolean  // WASM沙箱
    maxTokens?: number
    allowedTools?: string[]
  }
}

// Hand 配置
type HandConfig = {
  id: ID
  type: HandType
  schedule?: string   // Cron表达式
  enabled: boolean
  config: Record<string, any>  // Hand特定配置
}
```

---

## OpenFang 消息类型

```typescript
// Agent 请求
type AgentRequest = {
  agentId: ID
  conversationId: ID
  message: string
  context?: {
    history: Message[]
    metadata: Record<string, any>
  }
}

// Agent 响应
type AgentResponse = {
  agentId: ID
  conversationId: ID
  response: string
  toolCalls?: ToolCall[]
  metadata?: {
    tokensUsed: number
    model: string
    provider: string
  }
}

// 工具调用
type ToolCall = {
  tool: string
  args: Record<string, any>
  result?: any
}
```

---

## 依赖关系

```
基础类型 (01-base)
    │
    ├── 被核心模块依赖
    │   ├── P2P通信 (04-consensus)
    │   ├── 身份账号 (10-identity)
    │   ├── 服务市场 (11-service)
    │   ├── 托管支付 (11-service)
    │   └── Trust System (09-dispute)
    │
    └── 被可选插件依赖
        ├── LLM API Key托管 (05-compute)
        ├── Agent执行器 (05-compute)
        ├── Workflow执行器 (05-compute)
        ├── 想法系统 (07-idea)
        ├── 任务市场 (08-task)
        └── GenieBot界面 (12-misc)
```

---

## 类型使用示例

```typescript
// Token 金额操作
const price: TokenAmount = { amount: 1_000_000n, symbol: 'STT' }  // 1 STT
const fee: TokenAmount = { amount: 50_000n, symbol: 'STT' }       // 0.05 STT
const total = addCoins(price, fee)                                // 1.05 STT

// 加密数据
const apiKey: EncryptedData = {
  ciphertext: 'base64-encoded-ciphertext',
  nonce: 'base64-encoded-nonce',
  algorithm: 'x25519-xsalsa20-poly1305',
  version: 1
}

// Merkle 证明验证
const proof: MerkleProof = {
  root: '0x1234...',
  leaf: '0xabcd...',
  proof: ['0x1111...', '0x2222...'],
  index: 3
}
if (verifyProof(proof)) {
  // 证明有效
}
```

---

[返回索引](./00-index.md) | [下一章：共识层 →](./04-consensus.md)
