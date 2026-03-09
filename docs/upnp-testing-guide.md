# UPnP 实际测试指南

## 测试环境要求

### 硬件/网络要求
1. **家用路由器**（支持UPnP功能）
   - 主流品牌：TP-Link、小米、华为、华硕、网件等
   - 需要能登录路由器管理界面开启UPnP

2. **测试机器**
   - 连接到该路由器的电脑（有线或无线）
   - 能获取到内网IP（如192.168.x.x 或 10.x.x.x）

3. **网络环境**
   - 路由器有公网IP（或至少能访问互联网）
   - 不是公司/企业网络（通常有防火墙限制）

## 测试前准备

### 1. 开启路由器UPnP功能
```
以TP-Link为例：
1. 登录路由器管理界面（通常是 http://192.168.1.1 或 http://192.168.0.1）
2. 找到 "转发规则" 或 "NAT" 或 "UPnP" 菜单
3. 开启 "UPnP" 功能
4. 保存设置
```

### 2. 确认本机内网IP
```bash
# macOS/Linux
ifconfig | grep "inet " | grep -v 127.0.0.1

# 或
ip addr show | grep "inet " | grep -v 127.0.0.1

# 应该显示类似：192.168.1.100
```

### 3. 查看当前公网IP
```bash
# 方法1：通过命令行
curl ifconfig.me
curl ip.sb
curl cip.cc

# 方法2：浏览器访问
# https://www.whatismyip.com/
```

## 测试步骤

### 步骤1：编译带UPnP功能的节点
```bash
cd /Users/apple/projects/sharetoken

# 确保代码已更新
git pull origin main

# 编译
make build
# 或
go build -o bin/sharetokend ./cmd/sharetokend
```

### 步骤2：准备测试配置
```bash
# 创建测试目录
mkdir -p ~/.sharetoken-test

# 初始化节点（启用UPnP）
./bin/sharetokend init test-node --chain-id sharetoken-test --home ~/.sharetoken-test

# 修改配置启用UPnP
cat >> ~/.sharetoken-test/config/config.toml << 'EOF'

[p2p]
upnp = true
laddr = "tcp://0.0.0.0:26656"
EOF
```

### 步骤3：启动节点
```bash
# 方式1：直接启动（带日志）
./bin/sharetokend start --home ~/.sharetoken-test 2>&1 | grep -i "upnp\|nat\|external"

# 方式2：使用脚本启动
./scripts/devnet_multi.sh
```

### 步骤4：验证UPnP日志
预期看到的日志输出：
```
INFO Starting NAT manager upnp_enabled=true
INFO Attempting UPnP port mapping port=26656
INFO Found UPnP IGD v2  # 或 v1
INFO UPnP device discovered external_ip=203.x.x.x local_ip=192.168.1.100
INFO UPnP port mapping configured protocol=TCP external_port=26656 internal_port=26656 external_ip=203.x.x.x
```

## 验证方法

### 方法1：查看路由器UPnP映射表
```
登录路由器管理界面 -> 转发规则/NAT/UPnP -> 查看映射表

应该看到类似条目：
- 外部端口: 26656
- 内部IP: 192.168.1.100 (你的内网IP)
- 内部端口: 26656
- 协议: TCP
- 状态: Enabled
```

### 方法2：外部端口扫描测试
```bash
# 在另一台机器（或手机4G网络）上测试
# 使用 nc 或 telnet 测试公网IP+端口

telnet <你的公网IP> 26656
# 或
nc -vz <你的公网IP> 26656

# 预期结果：连接成功
```

### 方法3：在线端口检测工具
```
1. 访问 https://www.yougetsignal.com/tools/open-ports/
2. 输入你的公网IP
3. 输入端口 26656
4. 点击 Check

预期结果：端口显示为 OPEN
```

### 方法4：查看节点外部地址
```bash
# 使用节点CLI查询（如果有实现查询命令）
# 或通过API

curl http://localhost:26657/net_info
# 检查是否有正确的公网IP显示
```

## 常见问题排查

### 问题1："no UPnP IGD device found"
**原因**：
- 路由器UPnP未开启
- 不在同一个局域网
- 路由器不支持UPnP

**解决**：
```bash
# 1. 确认路由器UPnP已开启
# 2. 确认本机获取的是路由器分配的内网IP
ifconfig | grep "inet 192.168\|inet 10."

# 3. 重启路由器UPnP功能
```

### 问题2：端口映射成功但外部无法连接
**原因**：
- 运营商封锁端口
- 多层NAT（光猫+路由器）
- 防火墙拦截

**解决**：
```bash
# 1. 检查是否多层NAT
# 登录光猫查看WAN IP，如果是10.x.x.x或192.168.x.x，说明是内网IP

# 2. 尝试更换端口
# 修改 config.toml 使用其他端口如 26666

# 3. 关闭系统防火墙测试
sudo ufw disable  # Ubuntu
# 或
sudo systemctl stop firewalld  # CentOS
```

### 问题3："failed to add port mapping"
**原因**：
- 端口已被其他应用占用
- 权限不足
- 路由器拒绝

**解决**：
```bash
# 1. 检查端口占用
lsof -i :26656
netstat -tlnp | grep 26656

# 2. 更换端口测试
sed -i 's/26656/26666/g' ~/.sharetoken-test/config/config.toml
```

## 测试脚本

创建一个自动化测试脚本：

```bash
#!/bin/bash
# test_upnp.sh - UPnP测试脚本

echo "=== UPnP 测试开始 ==="

# 1. 检查环境
echo "1. 检查内网IP..."
LOCAL_IP=$(ip addr show | grep "inet 192.168\|inet 10." | head -1 | awk '{print $2}' | cut -d/ -f1)
echo "   内网IP: $LOCAL_IP"

echo "2. 检查公网IP..."
PUBLIC_IP=$(curl -s ifconfig.me)
echo "   公网IP: $PUBLIC_IP"

echo "3. 启动节点（10秒后检查）..."
./bin/sharetokend start --home ~/.sharetoken-test &
NODE_PID=$!
sleep 10

echo "4. 检查日志..."
grep -i "upnp\|external_ip" ~/.sharetoken-test/logs/*.log 2>/dev/null || echo "   未找到日志"

echo "5. 外部端口测试..."
echo "   请在其他网络执行: telnet $PUBLIC_IP 26656"

echo "6. 清理..."
kill $NODE_PID 2>/dev/null

echo "=== 测试完成 ==="
```

## 预期结果

### 成功标志
- [x] 日志显示"Found UPnP IGD"
- [x] 日志显示"UPnP port mapping configured"
- [x] 路由器管理界面显示映射条目
- [x] 外部端口扫描显示端口开放
- [x] 其他节点能通过公网IP+端口连接

### 失败处理
如果测试失败：
1. 记录错误日志
2. 检查路由器型号和固件版本
3. 尝试更新路由器固件
4. 考虑使用替代方案（内网穿透/云服务器）

## 下一步

UPnP测试成功后：
1. 更新 issue-002.md 标记UPnP为"✅ 测试通过"
2. 进行P2P连接测试（多节点通过公网互联）
3. 更新整体完成度至100%
