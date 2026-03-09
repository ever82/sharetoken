# x/identity - 身份账号与钱包

> **模块类型:** 核心模块
> **技术栈:** Go (Cosmos SDK)
> **位置:** `src/chain/x/identity`
> **依赖:** 基础类型 (01-base)、Cosmos SDK Auth 模块

---

## 概述

x/identity 是 ShareTokens 的核心模块，包含两个核心功能：
1. **身份账号** - 严格实名制 + 隐私保护的身份验证系统
2. **钱包** - 基于 Cosmos SDK Auth 模块的账户管理

作为 Cosmos SDK 模块实现，是所有其他模块的身份验证基础。

---

## 架构位置

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ShareTokens 核心模块                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────┐                                                               │
│  │ P2P通信  │                                                               │
│  └────┬─────┘                                                               │
│       │                                                                     │
│       ▼                                                                     │
│  ┌──────────────────────────────────────────┐                              │
│  │        身份账号与钱包 (10-identity)        │ ← 本章                       │
│  │  ┌──────────────┐   ┌──────────────┐     │                              │
│  │  │  身份账号     │   │   钱包       │     │                              │
│  │  │  (实名制)    │   │ (Keplr集成)  │     │                              │
│  │  └──────────────┘   └──────────────┘     │                              │
│  └───────────────────┬──────────────────────┘                              │
│                      │                                                      │
│                      ▼                                                      │
│              ┌──────────────┐                                               │
│              │  服务市场    │                                               │
│              └──────┬───────┘                                               │
│                     │                                                       │
│       ┌─────────────┴─────────────┐                                         │
│       ▼                           ▼                                         │
│  ┌──────────────┐         ┌──────────────────┐                             │
│  │  托管支付    │         │   Trust System   │                             │
│  └──────────────┘         └──────────────────┘                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 第一部分：身份账号

### 设计原则

```
1. 严格实名：必须通过微信/GitHub等验证才能注册
2. 隐私保护：链上只存哈希，不存明文
3. 零和唯一：每个身份只能绑定一个地址
4. 本地验证：通过 Merkle 证明验证，不需要全网查询
```

### Cosmos SDK 集成

```
x/identity/
├── module.go              # 模块定义
├── keeper/
│   ├── keeper.go          # 状态管理
│   ├── grpc_query.go      # gRPC 查询
│   └── msg_server.go      # 交易处理
├── types/
│   ├── keys.go            # 存储 key
│   ├── types.go           # 类型定义
│   ├── genesis.go         # 创世状态
│   └── msgs.go            # 消息类型
└── client/
    └── cli/               # CLI 命令
```

### IdentityType - 身份类型

```go
// types/types.go

type IdentityType string

const (
    IdentityTypeWechat IdentityType = "wechat"
    IdentityTypeGitHub IdentityType = "github"
    IdentityTypePhone  IdentityType = "phone"
    IdentityTypeEmail  IdentityType = "email"
    IdentityTypeGoogle IdentityType = "google"
)
```

### IdentityRegistry - 身份注册表

```go
// 身份哈希 → 绑定的地址
// key = sha256(identityType + ":" + identityId)

type IdentityRegistry struct {
    // 身份哈希 → 地址
    Identities Map<Hash, sdk.AccAddress>

    // 地址 → 身份哈希列表
    AddressToIdentities Map<sdk.AccAddress, []Hash>

    // 统计
    Stats RegistryStats
}

type RegistryStats struct {
    TotalIdentities uint64
    ByType          map[IdentityType]uint64
}
```

### IdentityRecord - 身份验证记录

```go
type IdentityRecord struct {
    IdentityHash Hash            // 不暴露真实身份
    IdentityType IdentityType
    Address      sdk.AccAddress

    // 验证证明
    Proof IdentityProof

    Status      string  // "active" | "revoked"
    TxHash      Hash
    BlockNumber uint64
}

type IdentityProof struct {
    ThirdPartySignature []byte
    ThirdPartyPublicKey string
    VerifiedAt          time.Time
}
```

### LocalIdentityStore - 本地存储（用户自己保管）

```go
// 注意：这是客户端本地存储，不在链上

type LocalIdentityStore struct {
    Identities []LocalIdentity
}

type LocalIdentity struct {
    Type        IdentityType
    Id          string           // 真实 ID（如微信 openId）
    Hash        Hash             // 计算出的哈希
    AccessToken *EncryptedData   // 访问令牌（加密存储）
}
```

### 验证等级

```go
type VerificationLevel string

const (
    LevelNone     VerificationLevel = "none"      // 未验证
    LevelBasic    VerificationLevel = "basic"     // 单一身份验证
    LevelVerified VerificationLevel = "verified"  // 多重身份验证
    LevelPremium  VerificationLevel = "premium"   // 高级验证 (KYC)
)

func DetermineLevel(identities []IdentityRecord) VerificationLevel {
    count := len(identities)

    switch {
    case count == 0:
        return LevelNone
    case count == 1:
        return LevelBasic
    case count >= 3:
        // 检查是否包含不同类型
        types := make(map[IdentityType]bool)
        for _, id := range identities {
            types[id.IdentityType] = true
        }
        if len(types) >= 3:
            return LevelPremium
        return LevelVerified
    default:
        return LevelVerified
    }
}
```

### 注册流程

```
1. 用户点击"微信验证"
2. 微信返回 openId（只在用户本地）
3. 本地计算：identityHash = sha256("wechat:" + openId)
4. 提交注册请求到链上（只发哈希）
5. 链上检查哈希是否已存在（本地查注册表）
6. 如果不存在，写入注册表
```

### Keeper 接口 - 身份相关

```go
// keeper/keeper.go

type Keeper struct {
    storeKey sdk.StoreKey
    cdc      codec.BinaryCodec
}

// 注册身份
func (k Keeper) RegisterIdentity(ctx sdk.Context, msg types.MsgRegisterIdentity) error

// 检查身份是否已注册
func (k Keeper) IsRegistered(ctx sdk.Context, address sdk.AccAddress) bool

// 获取身份哈希
func (k Keeper) GetIdentityHash(ctx sdk.Context, address sdk.AccAddress) (Hash, error)

// 检查身份哈希是否已存在（防重复）
func (k Keeper) HasIdentityHash(ctx sdk.Context, identityHash Hash) bool

// 获取地址的所有身份
func (k Keeper) GetIdentitiesByAddress(ctx sdk.Context, address sdk.AccAddress) ([]Hash, error)

// 撤销身份
func (k Keeper) RevokeIdentity(ctx sdk.Context, msg types.MsgRevokeIdentity) error

// 生成 Merkle 证明
func (k Keeper) GetIdentityProof(ctx sdk.Context, address sdk.AccAddress) (*MerkleProof, error)

// 验证 Merkle 证明
func (k Keeper) VerifyIdentityProof(ctx sdk.Context, proof MerkleProof) bool
```

### 消息类型

```go
// types/msgs.go

type MsgRegisterIdentity struct {
    Creator      sdk.AccAddress
    IdentityType IdentityType
    IdentityHash Hash
    Proof        IdentityProof
    Signature    []byte
}

type MsgRevokeIdentity struct {
    Creator      sdk.AccAddress
    IdentityHash Hash
    Signature    []byte
}

type MsgVerifyIdentity struct {
    Creator      sdk.AccAddress
    Address      sdk.AccAddress
    MerkleProof  MerkleProof
}
```

### gRPC 查询 - 身份相关

```protobuf
// query.proto

service Query {
    // 查询地址的身份状态
    rpc Identity(QueryIdentityRequest) returns (QueryIdentityResponse);

    // 查询地址的所有身份哈希
    rpc IdentitiesByAddress(QueryIdentitiesByAddressRequest) returns (QueryIdentitiesByAddressResponse);

    // 检查身份哈希是否存在
    rpc HasIdentityHash(QueryHasIdentityHashRequest) returns (QueryHasIdentityHashResponse);

    // 获取身份 Merkle 证明
    rpc IdentityProof(QueryIdentityProofRequest) returns (QueryProofResponse);

    // 获取注册表统计
    rpc Stats(QueryStatsRequest) returns (QueryStatsResponse);
}

message QueryIdentityRequest {
    string address = 1;
}

message QueryIdentityResponse {
    bool registered = 1;
    repeated string identity_hashes = 2;
    string status = 3;
}

message QueryHasIdentityHashRequest {
    bytes identity_hash = 1;
}

message QueryHasIdentityHashResponse {
    bool exists = 1;
    string bound_address = 2;  // 如果存在，返回绑定的地址
}
```

### 隐私保护机制

```
┌──────────────────────────────────────────────────────────────┐
│                    隐私保护设计                               │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  用户本地:                                                   │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  真实身份 (openId, email, phone)                    │    │
│  │  ↓                                                   │    │
│  │  identityHash = sha256("wechat:" + openId)          │    │
│  └─────────────────────────────────────────────────────┘    │
│                           │                                  │
│                           ▼                                  │
│  链上存储:                                                   │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  identityHash → address                             │    │
│  │                                                     │    │
│  │  只存储哈希，无法反推真实身份                        │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  第三方验证:                                                 │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  验证服务签名确认身份真实性                          │    │
│  │  不泄露具体身份信息                                  │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Merkle 证明验证

```go
// 本地验证身份，不需要全网查询

type MerkleProof struct {
    Root  Hash
    Leaf  Hash
    Proof []Hash
    Index uint64
}

func VerifyProof(proof MerkleProof) bool {
    current := proof.Leaf

    for i, sibling := range proof.Proof {
        if (proof.Index>>i)&1 == 0 {
            // 左子节点
            current = sha256(append(current, sibling...))
        } else {
            // 右子节点
            current = sha256(append(sibling, current...))
        }
    }

    return bytes.Equal(current, proof.Root)
}
```

---

## 第二部分：钱包

### 概述

钱包功能基于 Cosmos SDK Auth 模块实现，通过 Keplr 钱包插件提供用户界面。

### 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                       用户界面层                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                 Keplr 钱包插件                       │   │
│  │  - 账户管理                                         │   │
│  │  - 交易签名                                         │   │
│  │  - 资产显示                                         │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Cosmos SDK Auth 模块                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  - 账户 (BaseAccount)                               │   │
│  │  - 序列号 (nonce)                                   │   │
│  │  - 余额 (通过 x/bank)                               │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     x/identity 集成                          │
│  - 身份绑定                                                 │
│  - 验证等级                                                 │
│  - Reputation 关联                                          │
└─────────────────────────────────────────────────────────────┘
```

### 账户类型

```go
// 基于 Cosmos SDK BaseAccount

type Account struct {
    // 基础账户信息（由 Cosmos SDK Auth 提供）
    BaseAccount

    // 身份扩展（由 x/identity 提供）
    IdentityLevel VerificationLevel
    IdentityHashes []Hash

    // Reputation（由 x/dispute 提供）
    ReputationLevel ReputationLevel
    Reputation      uint64
}

// 标准 Cosmos SDK BaseAccount
type BaseAccount struct {
    Address       sdk.AccAddress
    PubKey        crypto.PubKey
    AccountNumber uint64
    Sequence      uint64
}
```

### 钱包集成配置

```typescript
// Keplr 链配置

export const shareTokensChainConfig = {
  chainId: 'sharetokens-1',
  chainName: 'ShareTokens',
  rpc: 'https://rpc.sharetokens.io',
  rest: 'https://api.sharetokens.io',
  bip44: {
    coinType: 118,  // Cosmos ATOM coin type
  },
  bech32Config: {
    bech32PrefixAccAddr: 'share',
    bech32PrefixAccPub: 'sharepub',
    bech32PrefixValAddr: 'sharevaloper',
    bech32PrefixValPub: 'sharevaloperpub',
    bech32PrefixConsAddr: 'sharevalcons',
    bech32PrefixConsPub: 'sharevalconspub',
  },
  currencies: [
    {
      coinDenom: 'STT',
      coinMinimalDenom: 'ustt',
      coinDecimals: 6,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: 'STT',
      coinMinimalDenom: 'ustt',
      coinDecimals: 6,
    },
  ],
  stakeCurrency: {
    coinDenom: 'STT',
    coinMinimalDenom: 'ustt',
    coinDecimals: 6,
  },
  gasPriceStep: {
    low: 0.01,
    average: 0.025,
    high: 0.04,
  },
}
```

### 余额与交易

```go
// 余额查询通过 x/bank 模块

type Balance struct {
    Address sdk.AccAddress
    Coins   sdk.Coins
}

// 交易类型
type Transaction struct {
    Hash        Hash
    From        sdk.AccAddress
    To          sdk.AccAddress
    Amount      sdk.Coins
    Gas         uint64
    GasPrice    sdk.DecCoin
    Memo        string
    Status      string  // "pending" | "success" | "failed"
    BlockHeight uint64
    Timestamp   time.Time
}
```

### 用户限额 (UserLimits)

> **说明:** 用户限额基于身份验证等级和 MQ 评分，提供风控保护。

```go
// types/limits.go

type UserLimits struct {
    Address sdk.AccAddress

    // 交易限制
    Trading TradingLimits

    // 提现限制
    Withdrawal WithdrawalLimits

    // 争议限制
    Disputes DisputeLimits

    // 服务限制（提供者）
    Service *ServiceLimits

    // 限制原因
    Restrictions []Restriction

    // 时间戳
    LastUpdated time.Time
}

type TradingLimits struct {
    // 单笔限额
    MaxOrderAmount sdk.Coin
    MinOrderAmount sdk.Coin

    // 日限额
    DailyLimit sdk.Coin
    DailyUsed  sdk.Coin

    // 月限额
    MonthlyLimit sdk.Coin
    MonthlyUsed  sdk.Coin

    // 订单数量
    MaxActiveOrders uint64
    MaxDailyOrders  uint64
}

type WithdrawalLimits struct {
    DailyLimit          sdk.Coin
    DailyUsed           sdk.Coin
    SingleWithdrawalMin sdk.Coin
    SingleWithdrawalMax sdk.Coin
    WithdrawalCooldown  time.Duration
}

type DisputeLimits struct {
    MaxActiveDisputes     uint64
    MinReputationToCreate uint64
    CooldownAfterLoss     time.Duration
}

type ServiceLimits struct {
    MaxConcurrentCalls uint64
    RateLimitPerMinute uint64
    MaxServices        uint64
}

type Restriction struct {
    Id        uint64
    Type      RestrictionType
    Reason    string
    ImposedAt time.Time
    ExpiresAt *time.Time
    ImposedBy *sdk.AccAddress  // 手动限制时
}

type RestrictionType string

const (
    RestrictionTypeNewAccount       RestrictionType = "new_account"        // 新账户限制
    RestrictionTypeLowVerification  RestrictionType = "low_verification"   // 低验证级别
    RestrictionTypeLowReputation    RestrictionType = "low_reputation"     // 低信誉
    RestrictionTypeDisputeLoss      RestrictionType = "dispute_loss"       // 争议败诉
    RestrictionTypePolicyViolation  RestrictionType = "policy_violation"   // 政策违规
    RestrictionTypeTempSuspension   RestrictionType = "temporary_suspension" // 临时暂停
    RestrictionTypePermanentBan     RestrictionType = "permanent_ban"      // 永久封禁
)

// 根据身份验证等级获取限额
func GetLimitsByVerificationLevel(level VerificationLevel) UserLimits {
    switch level {
    case LevelNone:
        return UserLimits{
            Trading: TradingLimits{
                MaxOrderAmount: sdk.NewCoin("stt", sdk.NewInt(100)),
                DailyLimit:     sdk.NewCoin("stt", sdk.NewInt(500)),
                MonthlyLimit:   sdk.NewCoin("stt", sdk.NewInt(2000)),
                MaxActiveOrders: 3,
                MaxDailyOrders:  5,
            },
            Withdrawal: WithdrawalLimits{
                DailyLimit:          sdk.NewCoin("stt", sdk.NewInt(200)),
                WithdrawalCooldown:  24 * time.Hour,
            },
            Disputes: DisputeLimits{
                MaxActiveDisputes: 1,
            },
        }
    case LevelBasic:
        return UserLimits{
            Trading: TradingLimits{
                MaxOrderAmount: sdk.NewCoin("stt", sdk.NewInt(500)),
                DailyLimit:     sdk.NewCoin("stt", sdk.NewInt(2000)),
                MonthlyLimit:   sdk.NewCoin("stt", sdk.NewInt(10000)),
                MaxActiveOrders: 10,
                MaxDailyOrders:  20,
            },
            Withdrawal: WithdrawalLimits{
                DailyLimit:         sdk.NewCoin("stt", sdk.NewInt(1000)),
                WithdrawalCooldown: 12 * time.Hour,
            },
            Disputes: DisputeLimits{
                MaxActiveDisputes: 3,
            },
        }
    case LevelVerified:
        return UserLimits{
            Trading: TradingLimits{
                MaxOrderAmount: sdk.NewCoin("stt", sdk.NewInt(2000)),
                DailyLimit:     sdk.NewCoin("stt", sdk.NewInt(10000)),
                MonthlyLimit:   sdk.NewCoin("stt", sdk.NewInt(50000)),
                MaxActiveOrders: 30,
                MaxDailyOrders:  50,
            },
            Withdrawal: WithdrawalLimits{
                DailyLimit:         sdk.NewCoin("stt", sdk.NewInt(5000)),
                WithdrawalCooldown: 6 * time.Hour,
            },
            Disputes: DisputeLimits{
                MaxActiveDisputes: 5,
            },
        }
    case LevelPremium:
        return UserLimits{
            Trading: TradingLimits{
                MaxOrderAmount: sdk.NewCoin("stt", sdk.NewInt(100000)),
                DailyLimit:     sdk.NewCoin("stt", sdk.NewInt(500000)),
                MonthlyLimit:   sdk.NewCoin("stt", sdk.NewInt(2000000)),
                MaxActiveOrders: 100,
                MaxDailyOrders:  200,
            },
            Withdrawal: WithdrawalLimits{
                DailyLimit:         sdk.NewCoin("stt", sdk.NewInt(100000)),
                WithdrawalCooldown: 1 * time.Hour,
            },
            Disputes: DisputeLimits{
                MaxActiveDisputes: 10,
            },
        }
    default:
        return UserLimits{}
    }
}
```

### Keeper 接口 - 限额相关

```go
// keeper/limits.go

// 获取用户限额
func (k Keeper) GetUserLimits(ctx sdk.Context, address sdk.AccAddress) UserLimits

// 检查是否超出限额
func (k Keeper) CheckLimitExceeded(ctx sdk.Context, address sdk.AccAddress, limitType string, amount sdk.Coin) bool

// 更新使用量
func (k Keeper) UpdateLimitUsage(ctx sdk.Context, address sdk.AccAddress, limitType string, amount sdk.Coin) error

// 添加限制
func (k Keeper) AddRestriction(ctx sdk.Context, address sdk.AccAddress, restriction Restriction) error

// 移除限制
func (k Keeper) RemoveRestriction(ctx sdk.Context, address sdk.AccAddress, restrictionId uint64) error
```

---

## 模块依赖

```
x/identity (核心模块)
    │
    ├── 核心依赖
    │   ├── x/auth     (账户认证 - Cosmos SDK)
    │   ├── x/bank     (代币转账 - Cosmos SDK)
    │   └── 基础类型 (01-base)
    │
    └── 被依赖
        ├── x/dispute   (Trust System - 身份验证)
        ├── x/compute   (算力交易 - 身份验证)
        ├── x/task      (任务市场 - 身份验证)
        ├── x/idea      (想法系统 - 身份验证)
        └── x/escrow    (资金托管 - 身份验证)
```

---

## 与链下服务交互

```
链下验证服务                            链上 x/identity
      │                                      │
      │  1. 用户发起身份验证                  │
      │                                      │
      │  2. 调用第三方 OAuth                  │
      │  (微信/GitHub/Google)                │
      │                                      │
      │  3. 获取用户身份 ID                   │
      │                                      │
      │  4. 计算身份哈希                      │
      │                                      │
      │  5. 签名验证结果                      │
      │                                      │
      │  6. 提交 MsgRegisterIdentity         │
      │─────────────────────────────────────►│
      │                                      │
      │  7. 链上验证签名并记录                │
      │                                      │
      │  8. 返回注册结果                      │
      │◄─────────────────────────────────────│
```

---

## OpenFang 集成

| OpenFang 组件 | 身份系统集成 |
|--------------|-------------|
| Keplr 钱包 | 账户管理、交易签名 |
| GenieBot Agent | 身份绑定、验证状态 |

---

[上一章：x/dispute](./09-dispute.md) | [返回索引](./00-index.md) | [下一章：服务市场 →](./11-service.md)
