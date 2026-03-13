# ShareToken 故障排查手册

> 基于项目开发经验和常见问题整理的故障排查指南

---

## 目录

- [节点启动问题](#节点启动问题)
- [网络连接问题](#网络连接问题)
- [交易失败问题](#交易失败问题)
- [密钥管理问题](#密钥管理问题)
- [常见错误码解释](#常见错误码解释)

---

## 节点启动问题

### 问题1: 节点启动失败，端口已被占用

**症状**:
```
ERROR: failed to listen: listen tcp 0.0.0.0:26656: bind: address already in use
```

**原因**:
- 另一个节点进程正在运行
- 之前的节点未正确停止
- 其他应用程序占用了该端口

**解决方案**:
```bash
# 1. 检查端口占用
lsof -i :26656
netstat -tlnp | grep 26656

# 2. 停止现有节点进程
./scripts/devnet_stop.sh

# 3. 或手动终止进程
pkill -f sharetokend

# 4. 更换端口测试
# 编辑 config/config.toml，修改 p2p 端口
sed -i 's/26656/26666/g' config/config.toml
```

---

### 问题2: devnet 目录数据残留导致启动失败

**症状**:
```
ERROR: failed to initialize node: genesis file mismatch
```
或节点启动后立即崩溃

**原因**:
- 多次启动开发网络后，旧数据导致新启动失败
- 创世文件配置不一致

**解决方案**:
```bash
# 1. 清理旧数据（开发网络启动脚本会自动执行）
rm -rf .devnet/

# 2. 重新初始化节点
./bin/sharetokend init <node_name> --chain-id sharetoken-devnet

# 3. 重新启动开发网络
./scripts/devnet_multi.sh
```

**预防措施**:
> 启动脚本 `devnet_multi.sh` 已包含自动清理逻辑，建议使用脚本启动而非手动启动。

---

### 问题3: 二进制文件不存在

**症状**:
```
ERROR: Binary not found: ./bin/sharetokend
```

**解决方案**:
```bash
# 1. 构建项目
make build

# 2. 或手动构建
go build -o bin/sharetokend ./cmd/sharetokend

# 3. 验证二进制文件
ls -la bin/sharetokend
./bin/sharetokend version
```

---

### 问题4: pprof 端口冲突

**症状**:
多节点网络启动时，第二个节点报错端口冲突

**原因**:
- 多个节点默认使用相同的 pprof 端口 6060

**解决方案**:
```bash
# 配置文件中为每个节点设置不同的 pprof 端口
# node0: 6060
# node1: 6061
# node2: 6062
# node3: 6063

# 修改 config.toml
sed -i 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:6061"/' config/config.toml
```

---

## 网络连接问题

### 问题1: P2P 节点无法相互发现

**症状**:
- 节点日志显示 `numPeers=0`
- 区块高度不增长
- 日志中出现 `No addresses to dial`

**解决方案**:
```bash
# 1. 检查节点间连接配置
curl http://localhost:26657/net_info

# 2. 确认 persistent_peers 配置正确
# 在 config/config.toml 中：
[p2p]
persistent_peers = "<node_id>@<ip>:<port>"

# 3. 获取节点 ID
./bin/sharetokend tendermint show-node-id --home .devnet/node0

# 4. 检查防火墙设置
# Linux
sudo iptables -L | grep 26656

# macOS
sudo pfctl -sr | grep 26656
```

---

### 问题2: UPnP 自动端口映射失败

**症状**:
```
WARN: no UPnP IGD device found
ERROR: failed to add port mapping
```

**原因**:
- 路由器 UPnP 功能未启用
- 多层 NAT（光猫+路由器）
- 防火墙阻止 SSDP 广播

**解决方案**:
```bash
# 1. 确认路由器 UPnP 已开启
# 登录路由器管理界面，找到 UPnP 设置并启用

# 2. 确认本机内网 IP
ifconfig | grep "inet 192.168\|inet 10."

# 3. 检查 SSDP 端口（UDP 1900）
# 确保防火墙允许 UDP 1900

# 4. 手动端口映射（如果 UPnP 不可用）
# 在路由器中手动添加端口映射规则：
# - 外部端口: 26656
# - 内部端口: 26656
# - 内部 IP: 本机内网 IP
# - 协议: TCP

# 5. 禁用 UPnP 使用手动配置
# 在 config.toml 中：
[p2p]
upnp = false
external_address = "<公网IP>:26656"
```

**参考**: 详细 UPnP 测试指南见 `docs/upnp-testing-guide.md`

---

### 问题3: 端口映射成功但外部无法连接

**症状**:
- 日志显示 UPnP 配置成功
- 但外部无法连接节点

**原因**:
- 运营商封锁端口
- 多层 NAT（光猫+路由器）
- 防火墙拦截

**解决方案**:
```bash
# 1. 检查是否多层 NAT
# 登录光猫查看 WAN IP，如果是 10.x.x.x 或 192.168.x.x，说明是内网 IP

# 2. 尝试更换端口
cat >> ~/.sharetoken/config/config.toml << 'EOF'
[p2p]
laddr = "tcp://0.0.0.0:26666"
EOF

# 3. 关闭系统防火墙测试
# Ubuntu
sudo ufw disable

# CentOS
sudo systemctl stop firewalld

# macOS
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate off
```

---

### 问题4: RPC 接口无法访问

**症状**:
```
curl: (7) Failed to connect to localhost port 26657: Connection refused
```

**解决方案**:
```bash
# 1. 检查节点是否在运行
./scripts/devnet_status.sh

# 2. 确认 RPC 配置正确
# config.toml:
[rpc]
laddr = "tcp://127.0.0.1:26657"

# 3. 如需外部访问，修改为
laddr = "tcp://0.0.0.0:26657"

# 4. 检查 CORS 设置（如果从前端访问）
cors_allowed_origins = ["*"]
```

---

## 交易失败问题

### 问题1: Gas 不足 (out of gas)

**症状**:
```
raw_log: "out of gas in location: WritePerByte; gasWanted: 200000, gasUsed: 201234"
code: 11
```

**解决方案**:
```bash
# 1. 增加 gas limit
./bin/sharetokend tx bank send <from> <to> 100000stake \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025stake

# 2. 或指定固定 gas
./bin/sharetokend tx bank send <from> <to> 100000stake \
  --gas 300000 \
  --gas-prices 0.025stake
```

**预防措施**:
- 使用 `--gas auto` 让节点自动估算
- 配合 `--gas-adjustment` 增加安全边际

---

### 问题2: Nonce/Sequence 错误 (sequence mismatch)

**症状**:
```
raw_log: "account sequence mismatch, expected 10, got 9"
code: 32
```

**原因**:
- 账户序列号不匹配
- 多客户端同时发送交易
- 交易未确认就发送下一笔

**解决方案**:
```bash
# 1. 查询正确的 sequence
./bin/sharetokend query account <address>

# 2. 使用正确的 sequence 发送交易
./bin/sharetokend tx bank send ... --sequence 10

# 3. 或使用 async 模式（不等待结果）
./bin/sharetokend tx bank send ... --broadcast-mode async

# 4. 查询当前 sequence
curl "http://localhost:1317/cosmos/auth/v1beta1/accounts/<address>"
```

**CosmJS 处理示例**:
```javascript
// 获取账户信息（包含 sequence）
const account = await client.getAccount(address);
if (account) {
  console.log("Sequence:", account.sequence);
}

// 错误处理
try {
  const result = await client.signAndBroadcast(...);
} catch (error) {
  if (error.message.includes("sequence mismatch")) {
    // 重试交易
    console.error("Sequence mismatch - retrying transaction");
  }
}
```

---

### 问题3: 余额不足 (insufficient funds)

**症状**:
```
raw_log: "insufficient funds: insufficient funds to pay for fees"
code: 5
```

**解决方案**:
```bash
# 1. 查询当前余额
./bin/sharetokend query bank balances <address>

# 2. 确认费用代币正确
# 检查 --fees 或 --gas-prices 使用的 denom 是否存在于账户

# 3. 检查是否有足够的代币支付费用
# 费用 = gas * gas-prices
# 例如: 200000 * 0.025stake = 5000stake
```

---

### 问题4: 交易签名失败

**症状**:
```
ERROR: key not found
ERROR: signing failed: unable to sign
```

**原因**:
- 密钥不存在于 keyring
- 使用了错误的 keyring-backend
- 密钥被锁定

**解决方案**:
```bash
# 1. 确认密钥存在
./bin/sharetokend keys list --keyring-backend test

# 2. 使用正确的 keyring-backend
./bin/sharetokend tx bank send ... --keyring-backend test
./bin/sharetokend tx bank send ... --keyring-backend file
./bin/sharetokend tx bank send ... --keyring-backend os

# 3. 如果使用 file backend，确保已解锁
# 系统会提示输入密码
```

---

### 问题5: 交易一直处于 pending 状态

**症状**:
- 交易提交后长时间未确认
- 查询不到交易哈希

**解决方案**:
```bash
# 1. 检查节点状态
./bin/sharetokend status

# 2. 检查网络是否出块
curl http://localhost:26657/status | jq .result.sync_info.latest_block_height

# 3. 查询交易状态
./bin/sharetokend query tx <txhash>

# 4. 如果网络卡住，可能需要重启
./scripts/devnet_stop.sh
./scripts/devnet_multi.sh
```

---

## 密钥管理问题

### 问题1: 密钥恢复失败

**症状**:
```
ERROR: invalid mnemonic
ERROR: checksum error
```

**原因**:
- 助记词拼写错误
- 单词顺序错误
- 使用了不支持的助记词语言

**解决方案**:
```bash
# 1. 确认助记词格式（BIP39）
# 应为 12 或 24 个英文单词

# 2. 重新恢复密钥
./bin/sharetokend keys add <name> --recover --keyring-backend test
# 然后输入正确的助记词

# 3. 检查恢复的地址是否正确
./bin/sharetokend keys show <name> --keyring-backend test
```

---

### 问题2: 密钥文件权限问题

**症状**:
```
ERROR: failed to read key file: permission denied
```

**解决方案**:
```bash
# 1. 检查密钥文件权限
ls -la ~/.sharetoken/keyring-file/

# 2. 修复权限
chmod 600 ~/.sharetoken/keyring-file/*
chmod 700 ~/.sharetoken/keyring-file/

# 3. 如果使用 OS keyring，检查系统密钥链访问权限
```

---

### 问题3: 创世账户未正确配置

**症状**:
- 新初始化节点没有初始余额
- 验证人无法开始出块

**解决方案**:
```bash
# 1. 添加创世账户
./bin/sharetokend add-genesis-account <address> 1000000000stake

# 2. 创建创世交易
gentx <name> 100000000stake --keyring-backend test

# 3. 收集创世交易
collect-gentxs

# 4. 验证创世文件
validate-genesis
```

---

### 问题4: node_key.json 损坏

**症状**:
```
ERROR: failed to load node key: invalid character
```

**解决方案**:
```bash
# 1. 备份并删除损坏的密钥文件
mv config/node_key.json config/node_key.json.bak

# 2. 重新初始化（会生成新密钥）
./bin/sharetokend init <moniker> --chain-id sharetoken-devnet

# 3. 注意：这将改变节点 ID，需要更新其他节点的 persistent_peers
```

**重要提示**:
> `node_key.json` 是节点身份标识，P2P 网络通过它识别节点。更换后其他节点需要更新配置才能连接。

---

## 常见错误码解释

### Cosmos SDK 标准错误码

| Code | 名称 | 说明 | 解决方案 |
|------|------|------|----------|
| 1 | ErrUnauthorized | 未授权操作 | 检查签名和权限 |
| 2 | ErrInsufficientFunds | 余额不足 | 充值账户余额 |
| 3 | ErrUnknownRequest | 未知请求 | 检查消息类型 |
| 4 | ErrInvalidAddress | 无效地址 | 验证地址格式 |
| 5 | ErrInvalidCoins | 无效代币 | 检查 denom 格式 |
| 6 | ErrOutOfGas | Gas 不足 | 增加 gas limit |
| 7 | ErrMemoTooLarge | Memo 过长 | 缩短 memo |
| 8 | ErrInsufficientFee | 费用不足 | 增加 fees |
| 9 | ErrTooManySignatures | 签名过多 | 减少签名人数 |
| 10 | ErrNoSignatures | 无签名 | 添加签名 |
| 11 | ErrOutOfGas | Gas 已用完 | 增加 gas limit |
| 12 | ErrInvalidSequence | 无效序列号 | 查询正确 sequence |
| 13 | ErrInvalidPubKey | 无效公钥 | 检查密钥配置 |

### CometBFT 错误码

| Code | 名称 | 说明 | 解决方案 |
|------|------|------|----------|
| -1 | ErrTxDecode | 交易解码失败 | 检查交易格式 |
| -2 | ErrInvalidSequence | 序列号错误 | 同步 sequence |
| -3 | ErrUnauthorized | 未授权 | 检查签名 |
| -4 | ErrInsufficientFunds | 资金不足 | 检查余额 |
| -5 | ErrUnknownRequest | 未知请求 | 检查 API 路径 |
| -6 | ErrInvalidAddress | 无效地址 | 检查地址格式 |

### 交易结果状态码

| Code | 状态 | 说明 |
|------|------|------|
| 0 | Success | 交易成功 |
| 1-99 | 系统错误 | SDK/CometBFT 错误 |
| 100+ | 应用错误 | 应用层自定义错误 |

### 查询交易结果

```bash
# 使用 CLI
./bin/sharetokend query tx <txhash>

# 使用 REST API
curl "http://localhost:1317/cosmos/tx/v1beta1/txs/<txhash>"

# 响应示例
{
  "tx_response": {
    "code": 0,
    "height": "65",
    "txhash": "5EDBBADD...",
    "gas_wanted": "200000",
    "gas_used": "180234",
    "raw_log": "[]"
  }
}
```

---

## 通用调试命令

### 查看节点状态
```bash
# 节点状态
./bin/sharetokend status

# 网络信息
curl http://localhost:26657/net_info | jq .

# 共识状态
curl http://localhost:26657/consensus_state | jq .

# 最新区块
curl http://localhost:26657/block | jq .result.block.header
```

### 查看日志
```bash
# 开发网络日志
tail -f .devnet/node0.log

# 使用 grep 过滤特定信息
grep -i "error\|failed\|panic" .devnet/node0.log

# 查看 UPnP 相关日志
grep -i "upnp\|nat\|external" .devnet/node0.log
```

### 验证配置
```bash
# 验证创世文件
./bin/sharetokend validate-genesis

# 查看节点配置
cat .devnet/node0/config/config.toml | grep -A5 "\[p2p\]"

# 查看应用配置
cat .devnet/node0/config/app.toml | grep -A3 "\[api\]"
```

---

## 参考资料

- [网络配置文档](./network-config.md)
- [UPnP 测试指南](./upnp-testing-guide.md)
- [标准开发流程](./standard-dev-process.md)
- [经验教训总结](./lessons-learned.md)

---

*最后更新: 2026-03-13*
*基于 Issue #1/2/3 开发经验及后续问题整理*
