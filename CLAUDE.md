# CLAUDE.md - ShareToken 项目开发规范

## 项目概述
这是一个基于 Cosmos SDK 的区块链项目，使用 Ignite CLI v29.9.0 创建。

## 开发方式
**TDD（测试驱动开发）**：先写测试，再写实现代码。每个功能都要有对应的单元测试。

## 开发优先级
按照 achievements 中的优先级顺序开发：
1. **P0** - 核心基础（ACH-DEV-001 到 ACH-DEV-003）
2. **P1** - 核心功能（ACH-DEV-004 到 ACH-DEV-010）
3. **P2** - 高级功能（ACH-DEV-011 到 ACH-DEV-018）
4. **P3** - 优化完善（ACH-DEV-019 到 ACH-DEV-022）

## 沟通规则

### 需要通知我的情况
1. **需要用户确认时** - 执行可能影响现有代码的操作前
2. **完成任务时** - 完成一个ACH-DEV任务后
3. **遇到阻塞问题时** - 无法继续开发的技术问题
4. **需要决策时** - 设计选择或实现方案不确定时

### 通知方式
**重要：在以下情况必须使用Bash工具执行echo命令通知CC：**

```bash
Bash(echo "CC-NOTIFY: [类型] 消息内容")
```

通知类型：
- `CC-NOTIFY: [完成] ACH-DEV-XXX - 描述` - 完成任务时
- `CC-NOTIFY: [阻塞] 问题描述` - 遇到技术阻塞时
- `CC-NOTIFY: [确认] 需要确认的操作描述` - 需要用户确认时
- `CC-NOTIFY: [决策] 需要决策的问题描述` - 需要做出设计决策时

示例：
```bash
Bash(echo "CC-NOTIFY: [完成] ACH-DEV-001 CI/CD Pipeline配置完成")
Bash(echo "CC-NOTIFY: [阻塞] proto-gen命令失败，缺少buf工具")
```

## 当前任务
查看 `docs/achievements/for-dev.md` 了解详细开发计划。
下一个任务：**ACH-DEV-004 Identity Module** (P1)

## 已完成任务
查看 `docs/achievements/done/` 目录了解已完成的P0任务。
查看 `docs/achievements/postponed.md` 了解延后实现的功能。

## 开发进度

### ✅ P0 - 核心基础（已完成）

#### ACH-DEV-001: Development Infrastructure ✅
- ✅ CI/CD Pipeline 配置完成
- ✅ Release工作流配置完成
- ✅ 本地开发网络启动脚本 (devnet_multi.sh)
- ✅ 代码规范与 Lint 配置
- 📄 详见 `docs/achievements/done/issue-001.md`

#### ACH-DEV-002: Blockchain Network Foundation ✅
- ✅ 4节点开发网络运行正常
- ✅ P2P网络配置与发现
- ✅ 共识机制配置
- ✅ UPnP/NAT端口映射（实际测试通过）
- ✅ Noise Protocol加密通信
- ⏭️ 区块浏览器（延后到部署阶段）
- 📄 详见 `docs/achievements/done/issue-002.md`

#### ACH-DEV-003: Wallet & Token System ✅
- ✅ STT代币定义与发行
- ✅ 余额查询接口
- ✅ 转账交易签名与广播
- ✅ Keplr钱包集成（代码完成）
- ✅ WalletConnect支持（代码完成）
- ✅ 交易历史查询
- ⏭️ 前端运行时测试（延后到部署阶段）
- 📄 详见 `docs/achievements/done/issue-003.md`

### 🔄 P1 - 核心功能（进行中）

等待开始：
- ⏳ ACH-DEV-004: Identity Module
- ⏳ ACH-DEV-005: Escrow Payment System
- ⏳ ACH-DEV-006: Oracle Service
- ⏳ ACH-DEV-007: MQ Scoring
- ⏳ ACH-DEV-008: Dispute Arbitration
- ⏳ ACH-DEV-009: Service Marketplace Core
- ⏳ ACH-DEV-010: Testnet Launch

## 技术栈
- Cosmos SDK v0.47.3
- CometBFT v0.37.x
- Go 1.21+
- Ignite CLI v29.9.0

## 环境检查清单
开始新任务前请确认：
- [ ] Go版本 >= 1.21 (`go version`)
- [ ] Ignite CLI已安装 (`ignite version`)
- [ ] GitHub CLI已配置 (`gh auth status`)
- [ ] Node.js已安装（前端相关任务）

## 常用命令
```bash
# 开发网络管理
./scripts/devnet_multi.sh        # 启动4节点开发网络
./scripts/devnet_multi.sh status # 查看网络状态
./scripts/devnet_stop.sh         # 停止开发网络

# 构建和测试
make build                       # 构建项目
make test                        # 运行测试
make lint                        # 运行Lint检查
make proto-gen                   # 生成protobuf代码

# 链操作
./bin/sharetokend query bank balances <address>
./bin/sharetokend tx bank send <from> <to> <amount> --chain-id sharetoken

# GitHub Issue管理
gh issue list --state open
gh issue close <number> --comment "完成评论"
```

## 项目文档
- `docs/achievements/for-dev.md` - 开发任务清单
- `docs/achievements/done/` - 已完成的任务
- `docs/achievements/postponed.md` - 延后的功能
- `docs/knowledges/standard-dev-process.md` - 标准开发流程
- `docs/knowledges/lessons-learned.md` - 经验教训总结
- `docs/knowledges/issue-template.md` - Issue文档模板
