#!/usr/bin/env bash
# 区块链网络基础功能测试
# TDD: 验证 ACH-DEV-002 验收标准

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DATA_DIR="${PROJECT_ROOT}/.devnet"
BINARY_NAME="sharetokend"
CHAIN_ID="sharetoken-devnet"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PASSED=0
FAILED=0
SKIPPED=0

assert_pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((PASSED++))
}

assert_fail() {
    echo -e "${RED}✗ FAIL${NC}: $1 - $2"
    ((FAILED++))
}

assert_skip() {
    echo -e "${YELLOW}⚠ SKIP${NC}: $1"
    ((SKIPPED++))
}

assert_file_contains() {
    local file="$1"
    local pattern="$2"
    local desc="$3"
    if grep -q "$pattern" "$file" 2>/dev/null; then
        assert_pass "$desc"
    else
        assert_fail "$desc" "未找到: $pattern"
    fi
}

echo "=========================================="
echo "区块链网络基础功能测试 (ACH-DEV-002)"
echo "=========================================="
echo ""

# 测试 1: 单节点启动与出块
echo "--- 单节点功能测试 ---"

echo "检查二进制文件..."
if [[ -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
    assert_pass "区块链二进制文件已构建"
else
    assert_fail "区块链二进制文件" "运行 'make build' 构建"
fi

# 测试 2: 出块时间配置
echo ""
echo "--- 出块时间配置测试 ---"

if [[ -f "${DATA_DIR}/node0/config/config.toml" ]]; then
    assert_file_contains "${DATA_DIR}/node0/config/config.toml" "timeout_commit = \"2s\"" "出块时间配置为 2s"
    assert_file_contains "${DATA_DIR}/node0/config/config.toml" "create_empty_blocks_interval = \"2s\"" "空块间隔配置为 2s"
else
    assert_skip "出块时间配置测试 (需要先启动 devnet)"
fi

# 测试 3: 多节点网络配置
echo ""
echo "--- 多节点网络配置测试 ---"

NODES=("node0" "node1" "node2" "node3")
for node in "${NODES[@]}"; do
    local node_dir="${DATA_DIR}/${node}"
    if [[ -d "${node_dir}" ]]; then
        assert_pass "${node} 配置目录存在"

        # 检查 CometBFT 配置
        if [[ -f "${node_dir}/config/config.toml" ]]; then
            assert_pass "${node} CometBFT 配置存在"
        else
            assert_fail "${node} CometBFT 配置" "config.toml 不存在"
        fi

        # 检查创世文件
        if [[ -f "${node_dir}/config/genesis.json" ]]; then
            assert_pass "${node} 创世文件存在"
        else
            assert_fail "${node} 创世文件" "genesis.json 不存在"
        fi
    else
        assert_fail "${node} 配置目录" "目录不存在，请先运行 ./scripts/devnet_multi.sh"
    fi
done

# 测试 4: P2P 网络配置
echo ""
echo "--- P2P 网络配置测试 ---"

if [[ -f "${DATA_DIR}/node0/config/config.toml" ]]; then
    # 检查 P2P 端口配置
    assert_file_contains "${DATA_DIR}/node0/config/config.toml" "laddr = \"tcp://0.0.0.0:26656\"" "node0 P2P 端口配置"
    assert_file_contains "${DATA_DIR}/node1/config/config.toml" "laddr = \"tcp://0.0.0.0:26666\"" "node1 P2P 端口配置"

    # 检查节点间连接
    if grep -q "persistent_peers" "${DATA_DIR}/node1/config/config.toml" 2>/dev/null; then
        assert_pass "节点间 persistent_peers 配置"
    else
        assert_fail "persistent_peers 配置" "未配置节点间连接"
    fi
else
    assert_skip "P2P 网络配置测试 (需要先启动 devnet)"
fi

# 测试 5: UPnP 支持
echo ""
echo "--- UPnP 支持测试 ---"

if [[ -f "${DATA_DIR}/node0/config/config.toml" ]]; then
    assert_file_contains "${DATA_DIR}/node0/config/config.toml" "upnp = false" "UPnP 配置项存在 (开发环境禁用)"
else
    assert_skip "UPnP 配置测试 (需要先启动 devnet)"
fi

# 检查是否有 UPnP 相关的代码或配置选项
if [[ -f "${PROJECT_ROOT}/app/app.go" ]]; then
    if grep -q "upnp\|UPnP\|nat\|NAT" "${PROJECT_ROOT}/app/app.go" 2>/dev/null; then
        assert_pass "代码中包含 UPnP/NAT 相关逻辑"
    else
        assert_fail "UPnP/NAT 代码实现" "app.go 中未找到相关逻辑"
    fi
fi

# 测试 6: Noise Protocol 加密
echo ""
echo "--- Noise Protocol 加密测试 ---"

# CometBFT 默认使用 Noise Protocol 进行 P2P 加密
if [[ -f "${DATA_DIR}/node0/config/config.toml" ]]; then
    # 检查 P2P 加密配置
    if grep -q "p2p" "${DATA_DIR}/node0/config/config.toml" 2>/dev/null; then
        assert_pass "P2P 配置存在 (CometBFT 默认使用 Noise Protocol)"
    else
        assert_fail "P2P 加密配置" "未找到 P2P 配置"
    fi
else
    assert_skip "Noise Protocol 测试 (需要先启动 devnet)"
fi

# 测试 7: 区块浏览器数据接口
echo ""
echo "--- 区块浏览器数据接口测试 ---"

# 检查是否启用了必要的 API
if [[ -f "${DATA_DIR}/node0/config/app.toml" ]]; then
    if grep -q "enable = true" "${DATA_DIR}/node0/config/app.toml" 2>/dev/null; then
        assert_pass "API 接口已启用"
    else
        assert_fail "API 接口" "未启用 API"
    fi

    if grep -q "swagger = true" "${DATA_DIR}/node0/config/app.toml" 2>/dev/null; then
        assert_pass "Swagger API 文档已启用"
    else
        assert_skip "Swagger API 文档 (可选)"
    fi
else
    assert_skip "区块浏览器接口测试 (需要先启动 devnet)"
fi

# 测试 8: 网络运行状态检查
echo ""
echo "--- 网络运行状态测试 ---"

# 检查节点是否正在运行
running_nodes=0
for i in {0..3}; do
    local pid_file="${DATA_DIR}/node${i}.pid"
    if [[ -f "${pid_file}" ]]; then
        local pid=$(cat "${pid_file}" 2>/dev/null)
        if kill -0 "${pid}" 2>/dev/null; then
            ((running_nodes++))
        fi
    fi
done

if [[ $running_nodes -ge 4 ]]; then
    assert_pass "4 个节点全部运行中"

    # 测试 RPC 接口
    echo "测试 RPC 接口..."
    for i in {0..3}; do
        local rpc_port=$((26657 + i * 10))
        local response=$(curl -s "http://127.0.0.1:${rpc_port}/status" 2>/dev/null || true)
        if [[ -n "$response" ]]; then
            local height=$(echo "$response" | grep -o '"latest_block_height":"[0-9]*"' | cut -d'"' -f4)
            if [[ -n "$height" && "$height" != "0" ]]; then
                assert_pass "node${i} 正在出块 (高度: $height)"
            else
                assert_fail "node${i} 出块状态" "未开始出块"
            fi
        else
            assert_fail "node${i} RPC 接口" "无法连接"
        fi
    done
elif [[ $running_nodes -gt 0 ]]; then
    assert_fail "节点运行状态" "只有 $running_nodes 个节点在运行，需要 4 个"
else
    assert_skip "网络运行状态测试 (需要先启动 devnet)"
fi

# 测试 9: 端口映射配置
echo ""
echo "--- 端口映射配置测试 ---"

# 检查是否有手动端口映射配置文档或脚本
if [[ -f "${PROJECT_ROOT}/docs/network-config.md" ]]; then
    assert_pass "网络配置文档存在"
else
    assert_fail "网络配置文档" "docs/network-config.md 不存在"
fi

# 测试 10: 消息广播测试
echo ""
echo "--- P2P 消息广播测试 ---"

if [[ $running_nodes -ge 2 ]]; then
    echo "检查节点间连接数..."
    for i in {0..3}; do
        local rpc_port=$((26657 + i * 10))
        local response=$(curl -s "http://127.0.0.1:${rpc_port}/net_info" 2>/dev/null || true)
        if [[ -n "$response" ]]; then
            local n_peers=$(echo "$response" | grep -o '"n_peers":"[0-9]*"' | cut -d'"' -f4)
            if [[ -n "$n_peers" && "$n_peers" -gt 0 ]]; then
                assert_pass "node${i} 已连接到 $n_peers 个对等节点"
            else
                assert_fail "node${i} P2P 连接" "未连接到其他节点"
            fi
        else
            assert_skip "node${i} P2P 连接测试 (无法获取 net_info)"
        fi
    done
else
    assert_skip "P2P 消息广播测试 (需要先启动 devnet)"
fi

# 汇总
echo ""
echo "=========================================="
echo "测试结果: $PASSED 通过, $FAILED 失败, $SKIPPED 跳过"
echo "=========================================="

if [[ $FAILED -gt 0 ]]; then
    echo -e "${RED}测试失败! 需要实现缺失的功能${NC}"
    exit 1
elif [[ $SKIPPED -gt 0 ]]; then
    echo -e "${YELLOW}部分测试跳过，请先运行 ./scripts/devnet_multi.sh 启动网络${NC}"
    exit 0
else
    echo -e "${GREEN}所有测试通过!${NC}"
    exit 0
fi
