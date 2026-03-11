# STT 经济模型修复方案

## 当前问题

1. **没有初始 STT 分配** - 创世时所有账户余额为0
2. **mint_denom 是 stake 不是 stt** - STT 没有增发机制
3. **任务完成没有支付集成** - 任务模块没有连接到银行模块

## 修复方案

### 1. 创世配置添加初始 STT 分配

需要在创世时给以下账户分配初始 STT：

```json
{
  "bank": {
    "balances": [
      {
        "address": "validator_address",
        "coins": [
          { "denom": "stt", "amount": "1000000000000" },  // 1,000,000 STT
          { "denom": "stake", "amount": "1000000000" }
        ]
      }
    ]
  }
}
```

### 2. 任务奖励银行集成

任务完成时需要：
1. 从托管释放资金到工作者
2. 从请求者扣除 STT
3. 平台收取手续费（比如 5%）

### 3. 本地 Agent 赚 STT 机制

用户可以通过以下方式赚 STT：
1. 完成任务获得奖励
2. 参与争议仲裁获得奖励
3. 提供验证服务获得奖励
4. 质押 stake 获得 stt 奖励

## 快速修复

### 方式1: 修改创世配置（立即生效）

```bash
# 在 config/genesis.json 中添加初始账户
./bin/sharetokend add-genesis-account cosmos1... 1000000000000stt --home config/
```

### 方式2: 空投机制

创建 faucet 服务，新用户注册时自动获得少量 STT。

### 方式3: 任务即挖矿

完成任务自动获得系统奖励的 STT（类似区块奖励）。

## 实施计划

1. ✅ 前端 GenieBot 连接 Agent Gateway
2. ⏳ 添加创世账户分配
3. ⏳ 集成银行模块到任务市场
4. ⏳ 实现 faucet 空投
