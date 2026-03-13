# ShareToken - 去中心化 AI 服务市场

<div align="center">

**一个基于 Cosmos SDK 的去中心化平台，让用户安全地使用和提供 AI 服务**

[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[快速开始](#快速开始) | [安装指南](#安装指南) | [功能特性](#功能特性) | [使用教程](#使用教程)

</div>

---

## 什么是 ShareToken?

ShareToken 是一个**去中心化的 AI 服务市场**，让您能够：

- 🤖 **零配置使用 AI** - 无需安装 Python、配置 API Key，一键使用主流 AI 模型
- 💰 **安全托管交易** - 资金由智能合约托管，确认满意后才付款
- 🏆 **信誉系统** - 基于社区投票的信誉评分，让优质服务脱颖而出
- 🔐 **完全掌控资产** - 私钥自持，不托管给任何第三方

---

## 系统要求

| 组件 | 最低版本 | 说明 |
|------|---------|------|
| **Go** | 1.22+ | 区块链节点运行环境 |
| **Node.js** | 18+ | 前端界面运行环境 |
| **操作系统** | macOS/Linux | Windows 需使用 WSL2 |
| **内存** | 4GB+ | 推荐 8GB |
| **磁盘** | 10GB+ | 区块链数据存储 |

---

## 安装指南

### 🖥️ 方法一：桌面应用（强烈推荐 - 开箱即用）

**最适合普通用户**，下载解压双击即可使用，无需命令行！

**Windows (绿色版):**
```powershell
# 下载
Invoke-WebRequest -Uri https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-win-x64.zip -OutFile ShareToken.zip

# 解压
Expand-Archive -Path ShareToken.zip -DestinationPath .\ShareToken

# 双击打开！
.\ShareToken\ShareToken.exe
```

**macOS:**
```bash
# Intel Mac
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-mac-x64.dmg
open ShareToken-0.1.0-mac-x64.dmg

# Apple Silicon M1/M2/M3
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-mac-arm64.dmg
open ShareToken-0.1.0-mac-arm64.dmg
```

**Linux (AppImage):**
```bash
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-linux-x86_64.AppImage
chmod +x ShareToken-0.1.0-linux-x86_64.AppImage
./ShareToken-0.1.0-linux-x86_64.AppImage
```

**桌面应用特点：**
- ✅ 开箱即用，无需安装 Node.js 或其他依赖
- ✅ 内置区块链节点（轻节点模式）
- ✅ 图形界面钱包（创建/导入/转账）
- ✅ 内置服务市场和 AI 对话界面
- ✅ 双击打开即可使用

---

### 💻 方法二：命令行工具（开发者/高级用户）

**适合开发者或喜欢命令行的用户**。

根据您的系统选择对应版本下载：

**macOS (Intel):**
```bash
# 下载最新版本
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/sharetoken_darwin_amd64.tar.gz

# 解压
tar -xzf sharetoken_darwin_amd64.tar.gz

# 移动到系统目录（可选）
sudo mv sharetokend /usr/local/bin/

# 验证安装
sharetokend version
```

**macOS (Apple Silicon M1/M2/M3):**
```bash
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/sharetoken_darwin_arm64.tar.gz
tar -xzf sharetoken_darwin_arm64.tar.gz
sudo mv sharetokend /usr/local/bin/
sharetokend version
```

**Linux:**
```bash
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/sharetoken_linux_amd64.tar.gz
tar -xzf sharetoken_linux_amd64.tar.gz
sudo mv sharetokend /usr/local/bin/
sharetokend version
```

**Windows (PowerShell):**
```powershell
# 下载最新版本
Invoke-WebRequest -Uri https://github.com/ever82/sharetoken/releases/latest/download/sharetoken_windows_amd64.zip -OutFile sharetoken_windows_amd64.zip

# 解压
Expand-Archive -Path sharetoken_windows_amd64.zip -DestinationPath .\

# 运行
.\sharetokend.exe version
```

> 📦 所有预编译版本可在 [Releases 页面](https://github.com/ever82/sharetoken/releases) 下载

> **注意：** Windows 版本需要开启开发者模式或使用 WSL2 运行完整功能。部分脚本功能可能需要 Git Bash 或 PowerShell。

---

### 方法二：从源码构建（开发者）

如果需要修改代码或参与开发，可以从源码构建：

```bash
# 1. 克隆项目
git clone https://github.com/ever82/sharetoken.git
cd sharetoken

# 2. 安装依赖并构建
make build

# 3. 验证构建
./bin/sharetokend version
```

---

### 安装前端界面（可选）

如果需要使用图形界面，安装前端依赖：

```bash
cd frontend
npm install
npm run serve
# 访问 http://localhost:8080
```

---

### 启动本地网络

```bash
# 一键启动 4 节点开发网络
./scripts/devnet_multi.sh
```

**启动成功后会显示：**
```
==========================================
开发网络已启动!
==========================================

节点信息:
  Chain ID: sharetoken-devnet

  node0:
    RPC: http://127.0.0.1:26657
    API: http://127.0.0.1:1317

  node1:
    RPC: http://127.0.0.1:26667
    API: http://127.0.0.1:1318

命令示例:
  查看状态: ./scripts/devnet_multi.sh status
  停止网络: ./scripts/devnet_stop.sh
```

### 第四步：启动前端界面

```bash
cd frontend
npm run serve
```

**访问 http://localhost:8080 即可看到界面**

---

## 快速开始

### 1️⃣ 创建您的第一个钱包

**方法 A：使用命令行（推荐新手）**

```bash
# 创建新钱包
./bin/sharetokend keys add mywallet --keyring-backend file

# 系统将提示您设置密码并显示助记词
# ⚠️ 重要：请安全保存助记词！
```

**方法 B：使用 Keplr 浏览器钱包（推荐日常使用）**

1. 安装 [Keplr 浏览器扩展](https://www.keplr.app/)
2. 打开 http://localhost:8080
3. 点击 "Connect Keplr"
4. 按提示添加 ShareToken 链

### 2️⃣ 获取测试代币

```bash
# 查看当前账户列表
./bin/sharetokend keys list --keyring-backend file

# 从水龙头获取测试币（需要配置 faucet）
# 或使用创世账户转账
./bin/sharetokend tx bank send validator0 <您的地址> 1000000stake \
    --chain-id sharetoken-devnet \
    --fees 1000stake \
    --yes
```

### 3️⃣ 查看余额

```bash
# 查询余额
./bin/sharetokend query bank balances <您的地址> \
    --node http://127.0.0.1:26657
```

**预期输出：**
```yaml
balances:
- amount: "1000000"
  denom: stake
pagination:
  next_key: null
  total: "0"
```

### 4️⃣ 发送转账

```bash
# 发送 1000 stake 到另一个地址
./bin/sharetokend tx bank send mywallet sharetoken1xxxxx... 1000stake \
    --chain-id sharetoken-devnet \
    --fees 1000stake \
    --memo "第一笔转账" \
    --yes
```

---

## 功能特性

### ✅ P0 - 核心功能（已实现）

| 功能 | 状态 | 说明 |
|------|------|------|
| 🔐 **安全数字钱包** | ✅ 可用 | Keplr + WalletConnect 双钱包支持 |
| 💰 **资金托管** | ✅ 可用 | 智能合约托管，确认后释放 |
| 💱 **代币转账** | ✅ 可用 | STT 代币秒级转账 |
| 🆔 **身份认证** | ✅ 可用 | 微信/GitHub/Google 一键登录 |

### 🔄 P1 - 完整体验（开发中）

| 功能 | 状态 | 说明 |
|------|------|------|
| 🤖 **AI 服务市场** | 🔄 开发中 | 一键调用 AI 模型 |
| 📊 **任务追踪** | 🔄 开发中 | 实时查看任务进度 |
| ⚖️ **争议仲裁** | 🔄 开发中 | 社区陪审团裁决 |
| 🏆 **信誉系统** | 🔄 开发中 | MQ 评分可视化 |

---

## 使用教程

### 教程 1：完整转账流程

```bash
# 步骤 1：创建两个测试钱包
./bin/sharetokend keys add alice --keyring-backend test
./bin/sharetokend keys add bob --keyring-backend test

# 步骤 2：查看钱包地址
./bin/sharetokend keys list --keyring-backend test
# 输出示例：
# - name: alice
#   address: sharetoken1abc...
# - name: bob
#   address: sharetoken1xyz...

# 步骤 3：从创世账户给 alice 充值
./bin/sharetokend tx bank send validator0 sharetoken1abc... 1000000stake \
    --chain-id sharetoken-devnet \
    --keyring-backend test \
    --keyring-dir ./.devnet/node0 \
    --fees 1000stake \
    --yes

# 步骤 4：查看 alice 余额
./bin/sharetokend query bank balances sharetoken1abc... \
    --node http://127.0.0.1:26657

# 步骤 5：alice 转账给 bob
./bin/sharetokend tx bank send alice sharetoken1xyz... 50000stake \
    --chain-id sharetoken-devnet \
    --keyring-backend test \
    --fees 1000stake \
    --yes

# 步骤 6：查看双方余额确认
./bin/sharetokend query bank balances sharetoken1abc... --node http://127.0.0.1:26657
./bin/sharetokend query bank balances sharetoken1xyz... --node http://127.0.0.1:26657
```

### 教程 2：查看交易历史

```bash
# 查询某地址的所有交易
./bin/sharetokend query txs \
    --events "message.sender='sharetoken1abc...'" \
    --node http://127.0.0.1:26657

# 根据交易哈希查询详情
./bin/sharetokend query tx <交易哈希> \
    --node http://127.0.0.1:26657
```

### 教程 3：使用 REST API

当开发网络运行后，您可以使用 HTTP 请求查询数据：

```bash
# 查询余额
curl http://127.0.0.1:1317/cosmos/bank/v1beta1/balances/sharetoken1abc...

# 查询账户信息
curl http://127.0.0.1:1317/cosmos/auth/v1beta1/accounts/sharetoken1abc...

# 查询最新区块
curl http://127.0.0.1:26657/block
```

---

## 项目结构

```
sharetoken/
├── bin/                    # 编译后的可执行文件
├── scripts/                # 实用脚本
│   ├── devnet_multi.sh    # 启动开发网络
│   ├── devnet_stop.sh     # 停止开发网络
│   └── test_wallet.sh     # 钱包测试
├── frontend/              # Vue.js 前端界面
│   ├── src/
│   │   ├── components/    # UI 组件
│   │   │   └── Wallet.vue # 钱包组件
│   │   └── utils/
│   │       ├── keplr.js   # Keplr 钱包集成
│   │       └── walletconnect.js # WalletConnect 集成
│   └── package.json
├── x/                     # Cosmos SDK 自定义模块
│   ├── sharetoken/        # STT 代币模块
│   ├── identity/          # 身份认证模块
│   ├── escrow/            # 资金托管模块
│   ├── trust/             # 信誉系统模块
│   └── ...
├── config.yml             # 链配置
└── readme.md              # 本文件
```

---

## 常见问题

### Q: 启动网络时提示端口被占用？

```bash
# 查找占用端口的进程
lsof -i :26657

# 停止开发网络
./scripts/devnet_stop.sh

# 或者强制清理
rm -rf .devnet/
./scripts/devnet_multi.sh
```

### Q: 查询余额显示为空？

确保：
1. 开发网络正在运行 (`./scripts/devnet_multi.sh status`)
2. 地址正确（使用 `./bin/sharetokend keys list` 查看）
3. 该地址已有转账记录或创世分配

### Q: 如何导出私钥？

```bash
# 导出为 ASCII 格式
./bin/sharetokend keys export mywallet --keyring-backend file

# 导出为 JSON 格式（含助记词）
./bin/sharetokend keys export mywallet --keyring-backend file --unarmored-hex --unsafe
```

---

## 开发团队

**ShareToken** 是一个开源项目，欢迎贡献！

- 💬 社区讨论：[GitHub Discussions](../../discussions)
- 🐛 问题反馈：[GitHub Issues](../../issues)

---

## 许可证

本项目采用 [MIT 许可证](LICENSE)

---

<div align="center">

**[⬆ 返回顶部](#sharetoken---去中心化-ai-服务市场)**

</div>
