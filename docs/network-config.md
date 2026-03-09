# 网络配置文档

ShareToken 区块链网络配置指南

## 网络架构

ShareToken 使用 Cosmos SDK + CometBFT 构建，支持以下网络特性：

- **单节点模式**: 本地开发和测试
- **多节点模式**: 4 节点共识网络
- **P2P 通信**: 基于 CometBFT 的 P2P 网络
- **UPnP/NAT**: 自动端口映射支持

## 端口配置

### 默认端口

| 服务 | 端口 | 说明 |
|------|------|------|
| RPC | 26657 | CometBFT RPC 接口 |
| P2P | 26656 | P2P 通信端口 |
| gRPC | 9090 | gRPC API 接口 |
| REST API | 1317 | HTTP API 接口 |
| Pprof | 6060 | 性能分析接口 |

### 多节点端口分配

| 节点 | RPC | P2P | gRPC | API |
|------|-----|-----|------|-----|
| node0 | 26657 | 26656 | 9090 | 1317 |
| node1 | 26667 | 26666 | 9091 | 1318 |
| node2 | 26677 | 26676 | 9092 | 1319 |
| node3 | 26687 | 26686 | 9093 | 1320 |

## UPnP 自动端口映射

### 配置

在 `config/config.toml` 中配置：

```toml
[p2p]
# 启用 UPnP 自动端口映射
upnp = true

# P2P 监听地址
laddr = "tcp://0.0.0.0:26656"

# 外部地址（可选，UPnP 会自动检测）
external_address = ""
```

### 工作原理

1. **自动发现**: 节点启动时自动发现 UPnP 设备
2. **端口映射**: 自动将内部端口映射到公网
3. **心跳维持**: 定期发送 SSDP 广播维持映射

### 系统要求

- 路由器支持 UPnP IGD 协议
- 网络防火墙允许 SSDP 广播
- 路由器 NAT 表有可用条目

## 手动端口映射

如果 UPnP 不可用，可以手动配置端口映射：

### 路由器配置

1. 登录路由器管理界面
2. 找到 "端口转发" 或 "虚拟服务器" 设置
3. 添加端口映射规则：
   - 外部端口: 26656
   - 内部端口: 26656
   - 内部 IP: 节点内网 IP
   - 协议: TCP

### 防火墙配置

```bash
# Linux (iptables)
sudo iptables -A INPUT -p tcp --dport 26656 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 26657 -j ACCEPT

# macOS
sudo pfctl -f /etc/pf.conf
# 添加规则: pass in proto tcp to any port {26656, 26657}
```

## 节点发现

### 种子节点

在 `config/config.toml` 中配置种子节点：

```toml
[p2p]
# 种子节点列表
seeds = "node0-id@seed1.sharetoken.network:26656,node1-id@seed2.sharetoken.network:26656"

# 持久连接节点
persistent_peers = "node0-id@192.168.1.100:26656"
```

### 获取节点 ID

```bash
sharetokend tendermint show-node-id
```

## 网络启动

### 单节点

```bash
make devnet
```

### 多节点

```bash
./scripts/devnet_multi.sh
```

### 检查网络状态

```bash
./scripts/devnet_status.sh
```

## 故障排查

### UPnP 不工作

1. 检查路由器 UPnP 是否启用
2. 检查防火墙是否允许 SSDP (UDP 1900)
3. 查看日志：`grep -i upnp ~/.sharetoken/logs/*.log`

### 节点无法连接

1. 检查端口是否开放：`nc -zv <ip> 26656`
2. 检查节点 ID 是否正确
3. 查看连接状态：`curl http://localhost:26657/net_info`

### 出块问题

1. 检查验证者是否在线
2. 检查共识状态：`curl http://localhost:26657/consensus_state`
3. 查看日志中的错误信息

## 安全配置

### P2P 加密

CometBFT 默认使用 Noise Protocol 加密 P2P 通信：

- 基于 Curve25519 密钥交换
- ChaCha20-Poly1305 加密
- 前向保密

### 节点密钥

节点密钥存储在 `config/node_key.json`：

```json
{
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "..."
  }
}
```

**注意**: 保护好 node_key.json，它是节点身份的唯一标识。

## 性能优化

### 连接数配置

```toml
[p2p]
# 最大入站连接数
max_num_inbound_peers = 40

# 最大出站连接数
max_num_outbound_peers = 10

# 发送缓冲区大小
send_rate = 5120000

# 接收缓冲区大小
recv_rate = 5120000
```

### 区块传播

```toml
[consensus]
# 区块传播超时
peer_gossip_sleep_duration = "100ms"

# 区块部分大小
block_part_size = "65536"
```

## 参考

- [CometBFT 文档](https://docs.cometbft.com/)
- [Cosmos SDK 网络配置](https://docs.cosmos.network/)
- [UPnP IGD 规范](http://upnp.org/specs/gw/UPnP-gw-InternetGatewayDevice-v2-Device.pdf)
