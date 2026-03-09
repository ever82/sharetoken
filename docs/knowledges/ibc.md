# IBC (Inter-Blockchain Communication) 知识文档

## 1. IBC 概述

### 什么是 IBC？

IBC（Inter-Blockchain Communication）是一种标准化的跨链通信协议，允许独立的区块链之间进行可靠、安全的数据传输。它最初由 Cosmos 生态系统开发，现已成为跨链通信的事实标准。

### 为什么 IBC 重要？

- **互操作性**：打破区块链孤岛，实现链与链之间的资产和数据自由流动
- **模块化**：支持任意数据类型的传输（代币、NFT、消息等）
- **安全性**：基于密码学验证，无需信任第三方中继者
- **可扩展性**：支持任意数量的链之间的连接

---

## 2. 核心概念

### 2.1 连接（Connections）

连接是两条链之间的信任基础，建立在轻客户端验证之上。

```
Chain A ←→ Connection ←→ Chain B
```

**关键属性：**
- 每条链维护对方链的轻客户端状态
- 连接握手：`CONN_OPEN_INIT` → `CONN_OPEN_TRY` → `CONN_OPEN_ACK` → `CONN_OPEN_CONFIRM`
- 连接标识符：如 `connection-0`, `connection-1`

### 2.2 通道（Channels）

通道建立在连接之上，用于特定应用的数据传输。

```
Connection
    ├── Channel 1 (transfer)
    ├── Channel 2 (ica)
    └── Channel 3 (custom)
```

**关键属性：**
- 通道类型：ORDERED（有序）或 UNORDERED（无序）
- 端口标识符：如 `transfer`, `ica`, `custom`
- 通道标识符：如 `channel-0`, `channel-1`
- 通道握手：`CHAN_OPEN_INIT` → `CHAN_OPEN_TRY` → `CHAN_OPEN_ACK` → `CHAN_OPEN_CONFIRM`

### 2.3 数据包（Packets）

数据包是跨链传输的基本单位。

```go
type Packet struct {
    Sequence           uint64          // 包序号
    SourcePort         string          // 源端口
    SourceChannel      string          // 源通道
    DestinationPort    string          // 目标端口
    DestinationChannel string          // 目标通道
    Data               []byte          // 应用数据
    TimeoutHeight      Height          // 超时区块高度
    TimeoutTimestamp   uint64          // 超时时间戳
}
```

**数据包生命周期：**
1. 发送链创建数据包
2. 中继者转发数据包到接收链
3. 接收链处理数据包并返回确认（Acknowledgement）
4. 中继者将确认传回发送链

---

## 3. ibc-go 模块结构

### 3.1 核心模块

```
ibc-go/
├── modules/
│   ├── core/                    # 核心模块
│   │   ├── 01-connection/       # 连接管理
│   │   ├── 02-channel/          # 通道管理
│   │   ├── 03-connection/       # 连接实现
│   │   ├── 04-channel/          # 通道实现
│   │   ├── 05-port/             # 端口管理
│   │   ├── 23-commitment/       # Merkle 证明
│   │   └── keeper/              # 核心 Keeper
│   ├── apps/                    # 应用层
│   │   ├── transfer/            # ICS-20 代币转账
│   │   ├── ica/                 # ICS-27 跨链账户
│   │   └── 27-interchain-accounts/
│   └── light-clients/           # 轻客户端
│       ├── 07-tendermint/       # Tendermint 客户端
│       ├── 06-solomachine/      # 单机客户端
│       └── 09-localhost/        # 本地客户端
```

### 3.2 核心 Keeper 接口

```go
// 核心 Keeper 结构
type Keeper struct {
    ClientKeeper     ClientKeeper
    ConnectionKeeper ConnectionKeeper
    ChannelKeeper    ChannelKeeper
    PortKeeper       PortKeeper
    Router           *porttypes.Router
}

// 基本操作
func (k Keeper) SendPacket(ctx sdk.Context, packet Packet) error
func (k Keeper) RecvPacket(ctx sdk.Context, packet Packet, proof []byte) error
func (k Keeper) AcknowledgePacket(ctx sdk.Context, packet Packet, ack []byte, proof []byte) error
func (k Keeper) TimeoutPacket(ctx sdk.Context, packet Packet, proof []byte) error
```

### 3.3 模块集成

```go
// app.go 中的集成
type App struct {
    // ...其他模块
    IBCKeeper        *ibckeeper.Keeper
    TransferKeeper   transferkeeper.Keeper
    ICAControllerKeeper icactrlkeeper.Keeper
    ICAHostKeeper    icahostkeeper.Keeper
}

// 路由配置
func (app *App) setupIBCModules() {
    // 注册 IBC 路由
    ibcRouter := porttypes.NewRouter()
    ibcRouter.AddRoute(transfer.ModuleName, transfer.NewIBCModule(app.TransferKeeper))
    ibcRouter.AddRoute(icacontroller.ModuleName, icacontroller.NewIBCModule(app.ICAControllerKeeper))
    app.IBCKeeper.SetRouter(ibcRouter)
}
```

---

## 4. ICS-20 代币跨链转账

### 4.1 转账流程

```
┌─────────┐    1.SendPacket    ┌─────────┐
│ Chain A │ ─────────────────→ │ Relayer │
│ (Source)│                    └─────────┘
└─────────┘                         │
     ↑                               │ 2.RelayPacket
     │                               ↓
     │ 4.Ack    ┌─────────────────────────┐
     │          │                         │
     └──────────│ Chain B (Destination)   │
                │                         │
                │ 3.OnRecvPacket          │
                │ (Mint/Burn)             │
                └─────────────────────────┘
```

### 4.2 代币处理逻辑

```go
// 发送端 (Chain A)
func (k Keeper) SendTransfer(ctx sdk.Context, sender sdk.AccAddress, token sdk.Coin, receiver string, destChain string) error {
    if isSourceChain {
        // 本地代币：锁定（escrow）
        k.LockToken(ctx, sender, token)
    } else {
        // 外部代币：销毁（burn）
        k.BurnToken(ctx, sender, token)
    }

    // 创建 IBC 数据包
    packet := FungibleTokenPacketData{
        Denom:    token.Denom,
        Amount:   token.Amount,
        Sender:   sender.String(),
        Receiver: receiver,
    }

    return k.SendPacket(ctx, packet)
}

// 接收端 (Chain B)
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet Packet) error {
    var data FungibleTokenPacketData
    proto.Unmarshal(packet.Data, &data)

    // 计算接收代币的 denom
    // 格式: ibc/{hash(packet.destPort + packet.destChannel + data.denom)}
    denom := k.DeriveDenom(packet.DestPort, packet.DestChannel, data.Denom)

    if isSourceChain(packet) {
        // 来源链：从 escrow 解锁
        k.UnlockToken(ctx, denom, data.Amount, receiver)
    } else {
        // 目标链：铸造新代币
        k.MintToken(ctx, denom, data.Amount, receiver)
    }

    return nil
}
```

### 4.3 Denom 追踪

```go
// IBC 代币 denom 格式
// 原始: uatom (Cosmos Hub)
// 跨链到 Osmosis: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2
//                                                    ↑
//                                        SHA256(port + channel + denom)

func (k Keeper) DenomHash(port, channel, denom string) string {
    hash := sha256.Sum256([]byte(fmt.Sprintf("%s/%s/%s", port, channel, denom)))
    return fmt.Sprintf("ibc/%X", hash)
}
```

### 4.4 超时处理

```go
// 如果数据包在超时前未被接收，可以回滚
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet Packet) error {
    var data FungibleTokenPacketData
    proto.Unmarshal(packet.Data, &data)

    if isSourceChain {
        // 返还锁定的代币
        k.UnlockToken(ctx, data.Denom, data.Amount, data.Sender)
    } else {
        // 重新铸造销毁的代币
        k.MintToken(ctx, data.Denom, data.Amount, data.Sender)
    }

    return nil
}
```

---

## 5. 轻客户端验证

### 5.1 轻客户端类型

| 客户端类型 | 用途 | 状态 |
|-----------|------|------|
| 07-tendermint | Tendermint/CometBFT 链 | 稳定 |
| 06-solomachine | 单机客户端 | 稳定 |
| 08-wasm | WASM 客户端（Substrate 等） | 稳定 |
| 09-localhost | 本地客户端 | 实验性 |

### 5.2 验证流程

```go
// 轻客户端验证核心逻辑
func (k ClientKeeper) VerifyMembership(
    ctx sdk.Context,
    clientID string,
    height Height,
    delayTimePeriod uint64,
    delayBlockPeriod uint64,
    proof []byte,
    path Path,
    value []byte,
) error {
    // 获取客户端状态
    clientState := k.GetClientState(ctx, clientID)

    // 检查是否经过足够的延迟（防止攻击）
    if err := k.VerifyDelayPeriod(ctx, clientID, height, delayTimePeriod, delayBlockPeriod); err != nil {
        return err
    }

    // 验证 Merkle 证明
    return clientState.VerifyMembership(
        k.cdc,
        k.GetClientConsensusState(ctx, clientID, height),
        proof,
        path,
        value,
    )
}
```

### 5.3 客户端更新

```go
// 更新轻客户端状态
func (k ClientKeeper) UpdateClient(ctx sdk.Context, clientID string, header Header) error {
    clientState := k.GetClientState(ctx, clientID)

    // 验证 header
    if err := clientState.VerifyHeader(k.cdc, k.GetClientConsensusState(ctx, clientID, height), header); err != nil {
        return err
    }

    // 更新客户端状态和共识状态
    newClientState, newConsensusState := clientState.UpdateState(header)
    k.SetClientState(ctx, clientID, newClientState)
    k.SetClientConsensusState(ctx, clientID, header.GetHeight(), newConsensusState)

    return nil
}
```

### 5.4 Merkle 证明结构

```
Root Hash
    │
    ├── [0] ────┬── Key: "ibc" ──── Value Hash
    │           │
    │           └── Sub-tree
    │                ├── "connections/connection-0"
    │                ├── "channels/channel-0"
    │                └── "packetCommitments/..."
    │
    └── [1] ──── Other data
```

---

## 6. Relayer 操作

### 6.1 Relayer 职责

- 监听链上的 IBC 事件
- 将数据包从源链中继到目标链
- 转发确认和超时
- 更新轻客户端状态

### 6.2 主流 Relayer 实现

| Relayer | 语言 | 特点 |
|---------|------|------|
| Hermes | Rust | 功能完整，官方推荐 |
| Go Relayer | Go | Cosmos 生态原生 |
| Ts-Relayer | TypeScript | 易于扩展 |

### 6.3 Relayer 核心逻辑

```go
// 简化的 Relayer 工作流程
func (r *Relayer) RelayPackets(srcChain, dstChain Chain) error {
    for {
        // 1. 查询源链未中继的数据包
        packets, err := srcChain.QueryUnrelayedPackets(dstChain.ClientID)
        if err != nil {
            return err
        }

        // 2. 为每个数据包获取证明
        for _, packet := range packets {
            proof, err := srcChain.QueryPacketProof(packet)
            if err != nil {
                continue
            }

            // 3. 更新目标链的轻客户端（如果需要）
            if err := r.UpdateClientIfNeeded(dstChain, srcChain); err != nil {
                continue
            }

            // 4. 提交数据包到目标链
            if err := dstChain.RecvPacket(packet, proof); err != nil {
                log.Printf("Failed to relay packet: %v", err)
                continue
            }

            // 5. 中继确认（可选）
            ack, err := dstChain.QueryPacketAcknowledgement(packet)
            if err != nil {
                continue
            }

            ackProof, err := dstChain.QueryAckProof(packet, ack)
            if err != nil {
                continue
            }

            srcChain.AcknowledgePacket(packet, ack, ackProof)
        }

        time.Sleep(r.PollInterval)
    }
}
```

### 6.4 Relayer 配置示例 (Hermes)

```toml
# hermes config.toml
[global]
log_level = "info"

[mode]
[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = true

[mode.channels]
enabled = true

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true

[[chains]]
id = "cosmoshub-4"
rpc_addr = "https://cosmos-rpc.example.com"
grpc_addr = "https://cosmos-grpc.example.com:443"
event_source = { mode = "push", url = "wss://cosmos-rpc.example.com/websocket" }
account_prefix = "cosmos"
key_name = "relayer"
store_prefix = "ibc"
gas_price = { price = 0.01, denom = "uatom" }

[[chains]]
id = "osmosis-1"
rpc_addr = "https://osmosis-rpc.example.com"
grpc_addr = "https://osmosis-grpc.example.com:443"
event_source = { mode = "push", url = "wss://osmosis-rpc.example.com/websocket" }
account_prefix = "osmo"
key_name = "relayer"
store_prefix = "ibc"
gas_price = { price = 0.01, denom = "uosmo" }
```

### 6.5 Relayer 命令示例

```bash
# 创建连接
hermes create connection cosmoshub-4 osmosis-1

# 创建通道
hermes create channel cosmoshub-4 osmosis-1 --port-a transfer --port-b transfer

# 启动中继
hermes start

# 手动中继数据包
hermes tx packet-recv cosmoshub-4 osmosis-1 --port transfer --channel channel-0
hermes tx packet-ack osmosis-1 cosmoshub-4 --port transfer --channel channel-0
```

---

## 7. ShareTokens 集成建议

### 7.1 模块集成清单

```go
// 在 app.go 中添加 IBC 支持
type App struct {
    // IBC 核心模块
    IBCKeeper      *ibckeeper.Keeper
    TransferKeeper transferkeeper.Keeper

    // 自定义 IBC 应用（如需要）
    // CustomIBCKeeper customkeeper.Keeper
}

// 初始化
func NewApp(...) *App {
    // 1. 初始化 IBC Keeper
    app.IBCKeeper = ibckeeper.NewKeeper(
        appCodec, keys[ibcexported.StoreKey],
        app.GetSubspace(ibcexported.ModuleName),
        stakingKeeper, upgradeKeeper,
    )

    // 2. 初始化 Transfer Keeper
    app.TransferKeeper = transferkeeper.NewKeeper(
        appCodec, keys[transfertypes.StoreKey],
        app.GetSubspace(transfertypes.ModuleName),
        app.IBCKeeper.ChannelKeeper,
        app.IBCKeeper.ChannelKeeper,
        app.BankKeeper,
        authtypes.FeeCollectorName,
    )

    // 3. 注册路由
    ibcRouter := porttypes.NewRouter()
    ibcRouter.AddRoute(transfertypes.ModuleName, transfer.NewIBCModule(app.TransferKeeper))
    app.IBCKeeper.SetRouter(ibcRouter)
}
```

### 7.2 自定义 IBC 应用

如果需要传输自定义数据（如 ShareTokens 的 idea/task 数据）：

```go
// 自定义 IBC 模块
type IBCModule struct {
    keeper Keeper
}

// 实现 IBCModule 接口
func (im IBCModule) OnChanOpenInit(...) (string, error)
func (im IBCModule) OnChanOpenTry(...) (string, error)
func (im IBCModule) OnChanOpenAck(...) error
func (im IBCModule) OnChanOpenConfirm(...) error
func (im IBCModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) error
func (im IBCModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, ack []byte) error
func (im IBCModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet) error
```

### 7.3 安全考虑

1. **超时设置**：合理设置数据包超时时间
2. **客户端信任等级**：配置合适的 trusting period
3. **中继者激励**：确保有足够的中继者参与
4. **监控**：监控 IBC 连接和数据包状态

---

## 8. 常用命令速查

```bash
# 查询连接
gaiad query ibc connection connections
gaiad query ibc connection connection-0

# 查询通道
gaiad query ibc channel channels
gaiad query ibc channel channel-0 --port transfer

# 查询数据包承诺
gaiad query ibc channel packet-commitments transfer channel-0
gaiad query ibc channel packet-receipt transfer channel-0 1

# 查询未中继数据包
gaiad query ibc channel unreceived-packets transfer channel-0

# 发送 IBC 转账
gaiad tx ibc-transfer transfer transfer channel-0 osmo1xxx 1000uatom --from user
```

---

## 参考资源

- [IBC Protocol Specification](https://github.com/cosmos/ibc)
- [ibc-go Documentation](https://ibc.cosmos.network/)
- [ICS-20 Fungible Token Transfer](https://github.com/cosmos/ibc/tree/main/spec/app/ics-020-fungible-token-transfer)
- [Hermes Relayer](https://hermes.informal.systems/)
- [Cosmos SDK IBC Module](https://docs.cosmos.network/main/modules/ibc)
