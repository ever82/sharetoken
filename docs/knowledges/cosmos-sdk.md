# Cosmos SDK 知识文档

## 概述

**Cosmos SDK** 是一个模块化的区块链开发框架，用于构建高性能、可互操作的权益证明（PoS）区块链。它提供了核心区块链功能的开箱即用组件，开发者只需关注业务逻辑。

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌─────────────────────────────────────────────────────────┐│
│  │              Custom Modules (x/)                         ││
│  │   x/compute  x/trust  x/escrow  x/identity              ││
│  └─────────────────────────────────────────────────────────┘│
│  ┌─────────────────────────────────────────────────────────┐│
│  │              Cosmos SDK Modules                          ││
│  │   auth  bank  staking  distribution  gov  params        ││
│  └─────────────────────────────────────────────────────────┘│
│  ┌─────────────────────────────────────────────────────────┐│
│  │              BaseApp (ABCI Handler)                      ││
│  │   Router  AnteHandler  GasMeter  Store                  ││
│  └─────────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│                    ABCI / ABCI++                            │
├─────────────────────────────────────────────────────────────┤
│                    CometBFT Core                            │
│        Consensus (BFT) + P2P Network + Mempool             │
└─────────────────────────────────────────────────────────────┘
```

---

## 1. 核心概念

### 1.1 模块（Modules）

模块是 Cosmos SDK 的核心构建块，每个模块封装了特定的业务逻辑。

```
模块职责划分:
- 状态定义（State）
- 消息处理（Message Handler）
- 查询处理（Query Handler）
- 状态变更触发（Event）
```

### 1.2 Keeper

Keeper 是模块的核心控制器，负责：
- 管理模块的状态存储
- 提供状态读写方法
- 跨模块交互（通过依赖其他模块的 Keeper）

```go
// Keeper 基本结构
type Keeper struct {
    storeKey     storetypes.StoreKey    // 状态存储键
    cdc          codec.BinaryCodec      // 编解码器
    paramSpace   paramstypes.Subspace   // 参数存储

    // 跨模块依赖
    bankKeeper   bankkeeper.Keeper
    authKeeper   authkeeper.AccountKeeper
}
```

### 1.3 Handler（消息处理器）

Handler 处理交易消息，执行状态变更。

```go
// 消息处理流程
func NewHandler(k Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        switch msg := msg.(type) {
        case *MsgCreateService:
            return handleMsgCreateService(ctx, k, msg)
        case *MsgUpdateService:
            return handleMsgUpdateService(ctx, k, msg)
        default:
            return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type")
        }
    }
}
```

### 1.4 Querier（查询处理器）

Querier 处理状态查询请求。

```go
// 查询处理（gRPC方式 - 推荐）
func (k Keeper) GetService(c context.Context, req *QueryGetServiceRequest) (*QueryGetServiceResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    ctx := sdk.UnwrapSDKContext(c)
    service, found := k.GetService(ctx, req.Id)
    if !found {
        return nil, status.Error(codes.NotFound, "service not found")
    }

    return &QueryGetServiceResponse{Service: service}, nil
}
```

---

## 2. 模块结构（x/ 目录模式）

```
chain/x/<module_name>/
├── keeper/
│   ├── keeper.go           # 主 Keeper 实现
│   ├── msg_server.go       # 交易消息服务器（gRPC）
│   ├── grpc_query.go       # 查询服务器（gRPC）
│   └── genesis.go          # Genesis 状态导入/导出
├── types/
│   ├── genesis.pb.go       # Genesis 状态类型
│   ├── module.pb.go        # 模块类型
│   ├── query.pb.go         # 查询类型
│   ├── tx.pb.go            # 交易类型
│   ├── keys.go             # 存储键常量
│   ├── errors.go           # 错误定义
│   └── expected_keepers.go # 依赖接口定义（用于 mock）
├── module.go               # 模块定义（实现 AppModuleBasic）
├── genesis.go              # Genesis 初始化
└── handler.go              # 消息路由（可选，gRPC 可替代）
```

---

## 3. 状态管理（Store）

### 3.1 存储键（Store Key）

每个模块需要自己的存储键来隔离状态。

```go
// types/keys.go
const (
    StoreKey = "store_key_name"
)

func KeyPrefix(p string) []byte {
    return []byte(p)
}

// 常用前缀
var (
    ServiceKey     = []byte{0x01}
    ProviderKey    = []byte{0x02}
    RequestKey     = []byte{0x03}
)
```

### 3.2 状态读写

```go
// 设置状态
func (k Keeper) SetService(ctx sdk.Context, service Service) {
    store := ctx.KVStore(k.storeKey)
    bz := k.cdc.MustMarshal(&service)
    store.Set(ServiceKey(service.Id), bz)
}

// 获取状态
func (k Keeper) GetService(ctx sdk.Context, id string) (Service, bool) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get(ServiceKey(id))
    if bz == nil {
        return Service{}, false
    }
    var service Service
    k.cdc.MustUnmarshal(bz, &service)
    return service, true
}

// 删除状态
func (k Keeper) DeleteService(ctx sdk.Context, id string) {
    store := ctx.KVStore(k.storeKey)
    store.Delete(ServiceKey(id))
}
```

### 3.3 迭代器

```go
// 获取所有服务
func (k Keeper) GetAllServices(ctx sdk.Context) (services []Service) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, ServiceKey)
    defer iterator.Close()

    for ; iterator.Valid(); iterator.Next() {
        var service Service
        k.cdc.MustUnmarshal(iterator.Value(), &service)
        services = append(services, service)
    }
    return services
}
```

---

## 4. 常用内置模块

| 模块 | 用途 | 主要功能 |
|------|------|----------|
| **x/auth** | 账户管理 | 账户创建、余额追踪、签名验证 |
| **x/bank** | 代币转账 | 代币发送、余额查询、供应管理 |
| **x/staking** | 质押 | 验证者管理、委托、解绑 |
| **x/distribution** | 奖励分配 | 质押奖励、手续费分配 |
| **x/gov** | 链上治理 | 提案、投票、执行 |
| **x/params** | 参数管理 | 模块参数存储和查询 |
| **x/slashing** | 惩罚机制 | 验证者惩罚、证据处理 |
| **x/upgrade** | 链升级 | 计划性升级管理 |

### 4.1 auth 模块使用

```go
// 获取账户
account := k.authKeeper.GetAccount(ctx, addr)

// 创建新账户
acc := k.authKeeper.NewAccountWithAddress(ctx, addr)
k.authKeeper.SetAccount(ctx, acc)
```

### 4.2 bank 模块使用

```go
// 发送代币
err := k.bankKeeper.SendCoins(ctx, senderAddr, recipientAddr, coins)

// 查询余额
balance := k.bankKeeper.GetBalance(ctx, addr, denom)

// 铸造代币（需要权限）
err := k.bankKeeper.MintCoins(ctx, moduleName, coins)
```

---

## 5. 创建自定义模块

### 5.1 定义 Proto 消息

```protobuf
// proto/<module>/tx.proto

syntax = "proto3";

package sharetokens.<module>;

option go_package = "github.com/sharetokens/x/<module>/types";

// 消息定义
message MsgCreateService {
    string creator = 1;
    string name = 2;
    string description = 3;
    repeated cosmos.base.v1beta1.Coin pricing = 4;
}

message MsgCreateServiceResponse {
    uint64 id = 1;
}
```

### 5.2 实现 Keeper

```go
// keeper/keeper.go

type Keeper struct {
    storeKey   storetypes.StoreKey
    cdc        codec.BinaryCodec
    paramSpace paramstypes.Subspace

    // 依赖
    bankKeeper   bankkeeper.Keeper
    authKeeper   authkeeper.AccountKeeper
}

func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey storetypes.StoreKey,
    paramSpace paramstypes.Subspace,
    bankKeeper bankkeeper.Keeper,
    authKeeper authkeeper.AccountKeeper,
) Keeper {
    return Keeper{
        storeKey:     storeKey,
        cdc:          cdc,
        paramSpace:   paramSpace,
        bankKeeper:   bankKeeper,
        authKeeper:   authKeeper,
    }
}
```

### 5.3 实现 MsgServer

```go
// keeper/msg_server.go

type msgServer struct {
    Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
    return &msgServer{Keeper: keeper}
}

func (k msgServer) CreateService(goCtx context.Context, msg *types.MsgCreateService) (*types.MsgCreateServiceResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // 验证
    creator, err := sdk.AccAddressFromBech32(msg.Creator)
    if err != nil {
        return nil, err
    }

    // 创建服务
    service := types.Service{
        Creator:     msg.Creator,
        Name:        msg.Name,
        Description: msg.Description,
        Pricing:     msg.Pricing,
        CreatedAt:   ctx.BlockTime(),
    }

    // 存储状态
    id := k.AppendService(ctx, service)

    // 发出事件
    ctx.EventManager().EmitTypedEvent(&types.EventServiceCreated{
        Creator: msg.Creator,
        Id:      id,
    })

    return &types.MsgCreateServiceResponse{Id: id}, nil
}
```

### 5.4 注册模块

```go
// app/app.go

// 添加 StoreKey
keys := sdk.NewKVStoreKeys(
    ..., // 其他模块
    types.StoreKey,
)

// 初始化 Keeper
app.ServiceKeeper = keeper.NewKeeper(
    appCodec,
    keys[types.StoreKey],
    app.GetSubspace(types.ModuleName),
    app.BankKeeper,
    app.AuthKeeper,
)

// 注册模块
app.ModuleBasics = module.NewBasicManager(
    ..., // 其他模块
    types.AppModuleBasic{},
)

// 注册 gRPC
types.RegisterMsgServer(app.MsgServiceRouter(), keeper.NewMsgServerImpl(app.ServiceKeeper))
types.RegisterQueryServer(app.GRPCQueryRouter(), keeper.NewQueryServerImpl(app.ServiceKeeper))
```

---

## 6. 交易生命周期

```
1. 用户签名交易
        ↓
2. 提交到 CometBFT RPC
        ↓
3. 进入 Mempool（CheckTx 验证）
   └── AnteHandler: 验证签名、gas、nonce
        ↓
4. 被选入区块（Proposer 选择）
        ↓
5. 区块共识达成
        ↓
6. FinalizeBlock 执行
   └── AnteHandler（再次验证）
   └── Router 路由到模块 Handler
   └── Handler 执行业务逻辑
   └── 状态变更 + 事件发出
        ↓
7. Commit 持久化状态
   └── 计算新的 AppHash
```

---

## 7. ShareTokens 项目应用

### 7.1 自定义模块清单

| 模块 | 职责 | 依赖 |
|------|------|------|
| **x/identity** | 身份验证、Sybil 防护 | auth, bank |
| **x/compute** | 服务市场（LLM/Agent/Workflow） | auth, bank, escrow, identity |
| **x/escrow** | 支付托管、争议锁定 | auth, bank |
| **x/trust** | MQ 评分、争议仲裁、零和再分配 | auth, identity, escrow |

### 7.2 推荐目录结构

```
chain/
├── app/
│   ├── app.go              # 应用主入口
│   ├── encoding.go         # 编解码配置
│   └── export.go           # 状态导出
├── cmd/
│   └── sharetokensd/
│       ├── main.go         # 主程序入口
│       └── cmd/            # CLI 命令
├── x/
│   ├── identity/           # 身份模块
│   ├── compute/            # 服务市场模块
│   ├── escrow/             # 托管模块
│   └── trust/              # 信任系统模块
├── proto/
│   ├── identity/
│   ├── compute/
│   ├── escrow/
│   └── trust/
└── go.mod
```

---

## 8. 开发工具

### 8.1 Ignite CLI（推荐）

```bash
# 安装
curl https://get.ignite.com/cli | bash

# 创建新链
ignite scaffold chain github.com/sharetokens/sharetokens

# 创建模块
ignite scaffold module trust --dep bank,auth

# 创建消息
ignite scaffold message create-service name description pricing

# 创建查询
ignite scaffold query service id --response name,description

# 启动开发链
ignite chain serve
```

### 8.2 常用命令

```bash
# 编译
go build ./cmd/sharetokensd

# 初始化节点
sharetokensd init node0 --chain-id sharetokens-1

# 创建密钥
sharetokensd keys add validator

# 添加创世账户
sharetokensd add-genesis-account validator 1000000000stake

# 创建验证者
sharetokensd gentx validator 100000000stake --chain-id sharetokens-1

# 收集创世交易
sharetokensd collect-gentxs

# 启动节点
sharetokensd start
```

---

## 9. 最佳实践

### 9.1 模块设计原则

1. **单一职责**：每个模块只负责一个领域
2. **依赖注入**：通过构造函数传入依赖的 Keeper
3. **接口隔离**：使用 `expected_keepers.go` 定义依赖接口
4. **事件驱动**：状态变更时发出事件，便于索引器监听

### 9.2 安全考虑

```go
// 1. 始终验证调用者
creator, err := sdk.AccAddressFromBech32(msg.Creator)
if err != nil {
    return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
}

// 2. 检查权限
if service.Creator != msg.Creator {
    return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not the creator")
}

// 3. 验证输入
if msg.Amount.IsNegative() {
    return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount cannot be negative")
}

// 4. 使用 AnteHandler 验证基础条件（签名、gas）
```

### 9.3 测试策略

```go
// keeper/keeper_test.go

func TestCreateService(t *testing.T) {
    // 创建测试上下文
    keeper, ctx := setupKeeper(t)

    // 测试创建服务
    msg := &types.MsgCreateService{
        Creator:     addr.String(),
        Name:        "Test Service",
        Description: "Test Description",
    }

    _, err := keeper.CreateService(ctx, msg)
    require.NoError(t, err)

    // 验证状态
    services := keeper.GetAllServices(ctx)
    require.Len(t, services, 1)
}
```

---

## 参考资料

- [Cosmos SDK 官方文档](https://docs.cosmos.network/)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network/)
- [Ignite CLI 文档](https://docs.ignite.com/)
- [Cosmos SDK 源码](https://github.com/cosmos/cosmos-sdk)
- [ABCI 规范](https://github.com/cometbft/cometbft/tree/main/spec/abci)
