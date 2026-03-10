# 用户验收指南

> 本指南帮助您像普通用户一样验证 ShareToken 的各项功能。
> 每项验证都有明确的步骤和预期结果。

---

## 📋 验收清单

### P0 - 核心功能（必须完成）

- [ ] **ACH-USER-001**: 安全数字钱包
  - [ ] 创建钱包 < 60 秒
  - [ ] 查看余额
  - [ ] 转账收款
  - [ ] 导出私钥

- [ ] **ACH-USER-002**: 一键 AI 访问
  - [ ] 无需配置使用 AI
  - [ ] 自然语言描述需求
  - [ ] 费用明细展示

- [ ] **ACH-USER-003**: 资金安全保证
  - [ ] 托管资金
  - [ ] 确认后释放
  - [ ] 争议冻结

- [ ] **ACH-USER-004**: 透明服务定价
  - [ ] 浏览服务市场
  - [ ] 查看定价模式
  - [ ] 预估费用

- [ ] **ACH-USER-005**: 首次入门
  - [ ] 一键登录
  - [ ] 自动创建钱包
  - [ ] 获得测试代币

### P1 - 完整体验（建议完成）

- [ ] **ACH-USER-006**: 任务进度追踪
- [ ] **ACH-USER-007**: 实名认证权益
- [ ] **ACH-USER-008**: 公平争议解决
- [ ] **ACH-USER-009**: 服务失败退款
- [ ] **ACH-USER-010**: 信誉仪表板
- [ ] **ACH-USER-011**: 可信社区

---

## 🚀 开始验收

### 准备环境

```bash
# 1. 进入项目目录
cd sharetoken

# 2. 确保已构建
make build

# 3. 启动开发网络
./scripts/devnet_multi.sh

# 4. 等待网络启动（约 10 秒）
./scripts/devnet_multi.sh status
```

---

## ✅ ACH-USER-001: 安全数字钱包

### 测试步骤

#### 步骤 1：创建钱包（计时测试）

```bash
# 开始计时，执行以下命令
./bin/sharetokend keys add testuser --keyring-backend test
```

**预期结果：**
- ✅ 命令在 60 秒内完成
- ✅ 显示助记词（24个单词）
- ✅ 显示地址（sharetoken1... 开头）

#### 步骤 2：查看余额

```bash
# 替换为您刚创建的地址
./bin/sharetokend query bank balances <您的地址> \
    --node http://127.0.0.1:26657
```

**预期结果：**
- ✅ 显示余额列表（初始可能为空）
- ✅ 无错误信息

#### 步骤 3：获取测试代币

```bash
# 从创世账户转账
./bin/sharetokend tx bank send validator0 <您的地址> 1000000stake \
    --chain-id sharetoken-devnet \
    --keyring-backend test \
    --keyring-dir ./.devnet/node0 \
    --fees 1000stake \
    --yes

# 再次查询余额
./bin/sharetokend query bank balances <您的地址> \
    --node http://127.0.0.1:26657
```

**预期结果：**
- ✅ 交易成功提交（显示 txhash）
- ✅ 余额显示 1000000 stake

#### 步骤 4：转账测试

```bash
# 创建第二个钱包
./bin/sharetokend keys add testuser2 --keyring-backend test

# 查看地址
./bin/sharetokend keys list --keyring-backend test

# 从第一个钱包转账到第二个
./bin/sharetokend tx bank send testuser <testuser2地址> 100000stake \
    --chain-id sharetoken-devnet \
    --keyring-backend test \
    --fees 1000stake \
    --yes

# 验证第二个钱包收到款项
./bin/sharetokend query bank balances <testuser2地址> \
    --node http://127.0.0.1:26657
```

**预期结果：**
- ✅ 转账成功
- ✅ testuser2 显示余额 100000 stake

#### 步骤 5：导出私钥

```bash
# 导出私钥
./bin/sharetokend keys export testuser --keyring-backend test

# 输入密码后，将显示私钥
```

**预期结果：**
- ✅ 成功提示输入密码
- ✅ 显示 ASCII 格式的私钥

### 验收结论

- [ ] 全部步骤通过
- [ ] 部分步骤失败（请记录）

**用户反馈：** _________________

---

## ✅ ACH-USER-002: 一键 AI 访问

### 测试步骤

#### 步骤 1：启动前端

```bash
cd frontend
npm install
npm run serve
```

#### 步骤 2：访问界面

打开浏览器访问 http://localhost:8080

**预期结果：**
- ✅ 页面正常加载
- ✅ 显示 ShareToken 钱包界面

#### 步骤 3：连接 Keplr

1. 安装 Keplr 浏览器扩展
2. 点击 "Connect Keplr"
3. 按提示添加 ShareToken 链

**预期结果：**
- ✅ Keplr 弹窗提示添加链
- ✅ 连接成功后显示钱包地址

### 验收结论

- [ ] 全部步骤通过
- [ ] 部分步骤失败（请记录）

---

## 📊 验收结果汇总

| 功能项 | 状态 | 备注 |
|--------|------|------|
| ACH-USER-001 | ⏳ 待验收 | |
| ACH-USER-002 | ⏳ 待验收 | |
| ACH-USER-003 | ⏳ 待验收 | |
| ACH-USER-004 | ⏳ 待验收 | |
| ACH-USER-005 | ⏳ 待验收 | |

**总体完成度：** __%

**是否通过 MVP 验收：** [ ] 是 [ ] 否

---

## 🐛 问题反馈

如果在验收过程中遇到问题，请记录：

1. **问题描述：**
2. **复现步骤：**
3. **预期结果：**
4. **实际结果：**
5. **错误信息：**

---

*最后更新：2026-03-10*
