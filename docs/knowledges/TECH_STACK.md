# ShareToken 技术栈

本文档详细说明 ShareToken 区块链项目的完整技术栈，包括各组件版本、依赖关系和兼容性要求。

---

## 1. 概述

ShareToken 是基于 Cosmos SDK 构建的区块链应用，采用分层架构设计：

| 层级 | 技术 | 主要用途 |
|------|------|----------|
| 区块链核心 | Go + Cosmos SDK | 共识、状态机、IBC |
| 前端应用 | Vue.js 3 + TypeScript | 用户界面、钱包交互 |
| 智能合约 | CosmWasm (可选) | 复杂业务逻辑 |

---

## 2. 后端技术栈 (区块链核心)

### 2.1 核心框架

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.19+ | 开发语言 |
| Cosmos SDK | v0.47.3 | 区块链应用框架 |
| CometBFT | v0.37.1 | 共识引擎 (原 Tendermint) |
| IBC-Go | v7.1.0 | 跨链通信协议 |

### 2.2 Cosmos SDK 生态组件

| 组件 | 版本 | 说明 |
|------|------|------|
| cosmossdk.io/api | v0.3.1 | Cosmos SDK API 定义 |
| cosmossdk.io/errors | v1.0.0-beta.7 | 错误处理库 |
| cosmossdk.io/core | v0.5.1 | 核心接口定义 |
| cosmossdk.io/depinject | v1.0.0-alpha.3 | 依赖注入框架 |
| cosmossdk.io/log | v1.1.0 | 日志库 |
| cosmossdk.io/math | v1.0.1 | 精确实算库 |
| cosmossdk.io/tools/rosetta | v0.2.1 | Rosetta API 支持 |
| cosmos/gogoproto | v1.4.10 | Protobuf 生成工具 |
| cosmos/iavl | v0.20.0 | 默克尔树存储 |
| cosmos/ics23/go | v0.10.0 | ICS23 证明格式 |

### 2.3 存储技术

| 组件 | 版本 | 说明 |
|------|------|------|
| cometbft-db | v0.7.0 | CometBFT 数据库接口 |
| goleveldb | v1.0.1 | LevelDB Go 绑定 |
| go.etcd.io/bbolt | v1.3.7 | BoltDB 嵌入式数据库 |
| dgraph-io/badger/v2 | v2.2007.4 | Badger 键值存储 |
| dgraph-io/ristretto | v0.1.1 | 内存缓存 |
| syndtr/goleveldb | v1.0.1-0.20210819022825-2ae1ddf74ef7 | LevelDB (替换版本) |

### 2.4 通信与 API

| 组件 | 版本 | 说明 |
|------|------|------|
| google.golang.org/grpc | v1.55.0 | gRPC 通信框架 |
| google.golang.org/protobuf | v1.30.0 | Protobuf Go 运行时 |
| github.com/golang/protobuf | v1.5.3 | Protobuf Go 生成器 |
| google.golang.org/genproto | v0.0.0-20230306155012-7f2fa6fef1f4 | 生成代码 |
| grpc-gateway | v1.16.0 / v2.15.2 | RESTful API 网关 |
| gorilla/mux | v1.8.0 | HTTP 路由器 |
| gorilla/websocket | v1.5.0 | WebSocket 支持 |
| rs/cors | v1.8.3 | CORS 中间件 |

### 2.5 加密与安全

| 组件 | 版本 | 说明 |
|------|------|------|
| golang.org/x/crypto | v0.8.0 | Go 加密库 |
| filippo.io/edwards25519 | v1.0.0 | Ed25519 曲线实现 |
| decred/dcrd/dcrec/secp256k1/v4 | v4.1.0 | Secp256k1 曲线 |
| btcsuite/btcd/btcec/v2 | v2.3.2 | 比特币 ECDSA |

### 2.6 CLI 与工具

| 组件 | 版本 | 说明 |
|------|------|------|
| spf13/cobra | v1.6.1 | CLI 框架 |
| spf13/pflag | v1.0.5 | POSIX 风格命令行标志 |
| spf13/viper | v1.15.0 | 配置管理 |
| rs/zerolog | v1.29.1 | 结构化日志 |

---

## 3. 前端技术栈

### 3.1 核心框架

| 组件 | 版本 | 说明 |
|------|------|------|
| Vue.js | ^3.3.4 | 前端框架 |
| TypeScript | ^5.x | 类型安全的 JavaScript |
| Vue Router | ^4.2.4 | 路由管理 |
| Vuex | ^4.1.0 | 状态管理 |
| Core-js | ^3.32.0 | JavaScript  polyfill |

### 3.2 Cosmos 生态 JavaScript 库

| 组件 | 版本 | 说明 |
|------|------|------|
| @cosmjs/amino | ^0.31.0 | Amino 编码支持 |
| @cosmjs/cosmwasm-stargate | ^0.31.0 | CosmWasm 客户端 |
| @cosmjs/crypto | ^0.31.0 | 加密工具 |
| @cosmjs/encoding | ^0.31.0 | 编码工具 |
| @cosmjs/proto-signing | ^0.31.0 | Protobuf 签名 |
| @cosmjs/stargate | ^0.31.0 | Cosmos SDK 客户端 |
| @keplr-wallet/types | ^0.12.0 | Keplr 钱包类型 |

### 3.3 钱包集成

| 组件 | 版本 | 说明 |
|------|------|------|
| @walletconnect/client | ^1.8.0 | WalletConnect 客户端 |
| @walletconnect/qrcode-modal | ^1.8.0 | WalletConnect 二维码 |

### 3.4 开发工具

| 组件 | 版本 | 说明 |
|------|------|------|
| @vue/cli-service | ^5.0.8 | Vue CLI 服务 |
| @babel/eslint-parser | ^7.22.0 | ESLint Babel 解析器 |
| eslint | ^8.47.0 | 代码检查工具 |
| eslint-plugin-vue | ^9.17.0 | Vue ESLint 插件 |

---

## 4. 开发工具

### 4.1 构建工具

| 工具 | 版本/说明 |
|------|-----------|
| Make | 构建自动化 |
| Ignite CLI | 推荐 (用于脚手架和代码生成) |
| gofmt | Go 代码格式化 |
| golangci-lint | Go 代码检查 |

### 4.2 Protobuf 工具

| 工具 | 说明 |
|------|------|
| cosmos/gogoproto | Cosmos 定制的 Protobuf 生成器 |
| grpc-gateway | RESTful API 生成 |

---

## 5. 依赖关系图

```
┌─────────────────────────────────────────────────────────────────┐
│                        ShareToken Application                    │
├─────────────────────────────────────────────────────────────────┤
│  Cosmos SDK v0.47.3                                             │
│  ├── Core (auth, bank, staking, distribution, gov, mint...)    │
│  ├── IBC-Go v7.1.0 (transfer, interchain-accounts)              │
│  └── Custom Modules (trust, service, etc.)                      │
├─────────────────────────────────────────────────────────────────┤
│  CometBFT v0.37.1 (Consensus + Networking)                      │
│  └── cometbft-db v0.7.0                                         │
├─────────────────────────────────────────────────────────────────┤
│  Storage Layer                                                  │
│  ├── LevelDB (goleveldb)                                        │
│  ├── BoltDB (bbolt)                                            │
│  └── Badger (optional)                                         │
├─────────────────────────────────────────────────────────────────┤
│  gRPC / REST API                                                │
│  ├── grpc-gateway v1.16.0 / v2.15.2                            │
│  └── gorilla/mux v1.8.0                                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Frontend Application                        │
│  ├── Vue.js 3 + TypeScript                                      │
│  ├── @cosmjs/* v0.31.0 (Cosmos client libraries)               │
│  ├── Keplr Wallet Integration                                   │
│  └── WalletConnect Integration                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. 版本兼容性矩阵

### 6.1 Go 版本兼容性

| Go 版本 | 支持状态 | 说明 |
|---------|----------|------|
| 1.19.x | 推荐 | 当前 go.mod 指定版本 |
| 1.20.x | 支持 | 测试通过 |
| 1.21.x | 支持 | 开发环境版本 |

### 6.2 Cosmos SDK 兼容性

| Cosmos SDK | CometBFT | IBC-Go | 状态 |
|------------|----------|--------|------|
| v0.47.x | v0.37.x | v7.x | 当前使用 |
| v0.46.x | v0.34.x | v6.x | 不兼容 |
| v0.45.x | v0.34.x | v3.x | 不兼容 |

### 6.3 前端兼容性

| cosmjs 版本 | Cosmos SDK | 兼容性 |
|-------------|------------|--------|
| v0.31.x | v0.47.x | 完全兼容 |
| v0.30.x | v0.46.x | 不兼容 |

---

## 7. 升级注意事项

### 7.1 Cosmos SDK v0.47 升级要点

1. **CometBFT 替换**: 从 Tendermint v0.34 迁移到 CometBFT v0.37
2. **IAVL 升级**: iavl v0.20 引入新的存储格式
3. **GRPC 升级**: grpc v1.55 可能需要证书更新
4. **Protobuf 更新**: gogoproto v1.4.10 生成代码格式变化

### 7.2 数据库兼容性

- LevelDB 数据库在 v0.47 中保持向后兼容
- 建议升级前备份 `data/` 目录

### 7.3 API 变更

- REST API 路径保持兼容
- gRPC 服务定义无破坏性变更
- 查询参数格式保持一致

---

## 8. 参考资源

| 资源 | URL | 说明 |
|------|-----|------|
| Cosmos SDK 文档 | https://docs.cosmos.network/ | 核心框架文档 |
| Cosmos SDK v0.47 | https://docs.cosmos.network/v0.47 | 特定版本文档 |
| CometBFT | https://cometbft.com/ | 共识引擎 |
| IBC-Go | https://ibc.cosmos.network/ | 跨链协议 |
| Ignite CLI | https://ignite.com/ | 开发脚手架 |
| Keplr Wallet | https://keplr.app/ | 钱包集成 |
| CosmJS | https://cosmos.github.io/cosmjs/ | JavaScript 库 |

---

## 9. 版本历史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2024-XX | v0.1.0 | 初始版本，基于 Cosmos SDK v0.47.3 |

---

*最后更新: 2026-03-13*
