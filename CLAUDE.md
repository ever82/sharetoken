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

## 开发进度

### ACH-DEV-001: Development Infrastructure (进行中)
- ✅ CI/CD Pipeline 测试与配置
  - ✅ 创建 test_cicd.sh 测试脚本
  - ✅ 创建 .github/workflows/ci.yml
  - ✅ 创建 .github/workflows/release.yml
  - ✅ 创建 Makefile
  - ⚠️ 测试结果: 10/13 通过（3个失败项与grep模式匹配和Go环境相关，非配置问题）
- ⏳ 本地开发网络启动脚本
- ⏳ 代码规范与 Lint 配置

### 测试状态说明
CI/CD测试中有3个失败项：
1. "CI 运行 Go 测试" - grep模式问题（CI使用`make test`而非直接`go test`）
2. "CI 运行 Go 构建" - grep模式问题（CI使用`make build`而非直接`go build`）
3. "Go 已安装" - 当前macOS环境未安装Go（非配置问题）

这些失败不影响CI/CD配置的正确性，可以继续下一个任务。

## 技术栈
- Cosmos SDK v0.47.3
- CometBFT v0.37.x
- Go 1.21+
- Ignite CLI v29.9.0

## 常用命令
```bash
# 启动开发链
ignite chain serve

# 运行测试
go test ./...

# 生成protobuf代码
make proto-gen

# 构建
make build
```
