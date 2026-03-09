#!/usr/bin/env bash
# 本地开发网络启动脚本测试
# TDD: 先写测试，验证本地开发网络配置是否符合要求

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

assert_file_exists() {
    local file="$1"
    local desc="$2"
    if [[ -f "$file" ]]; then
        echo -e "${GREEN}✓ PASS${NC}: $desc"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $desc - 文件不存在: $file"
        ((FAILED++))
    fi
}

assert_dir_exists() {
    local dir="$1"
    local desc="$2"
    if [[ -d "$dir" ]]; then
        echo -e "${GREEN}✓ PASS${NC}: $desc"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $desc - 目录不存在: $dir"
        ((FAILED++))
    fi
}

assert_script_executable() {
    local file="$1"
    local desc="$2"
    if [[ -f "$file" && -x "$file" ]]; then
        echo -e "${GREEN}✓ PASS${NC}: $desc"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $desc - 脚本不存在或不可执行: $file"
        ((FAILED++))
    fi
}

assert_contains() {
    local file="$1"
    local pattern="$2"
    local desc="$3"
    if grep -q "$pattern" "$file" 2>/dev/null; then
        echo -e "${GREEN}✓ PASS${NC}: $desc"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $desc - 未找到: $pattern"
        ((FAILED++))
    fi
}

echo "=========================================="
echo "本地开发网络启动脚本测试"
echo "=========================================="
echo ""

# 测试 1: 脚本存在
echo "--- 脚本文件测试 ---"
assert_script_executable "${PROJECT_ROOT}/scripts/devnet_multi.sh" "多节点启动脚本存在且可执行"
assert_script_executable "${PROJECT_ROOT}/scripts/devnet_stop.sh" "停止脚本存在且可执行"
assert_script_executable "${PROJECT_ROOT}/scripts/devnet_status.sh" "状态检查脚本存在且可执行"

echo ""
echo "--- 脚本内容测试 ---"

# 测试 2: 多节点脚本包含必要的配置
if [[ -f "${PROJECT_ROOT}/scripts/devnet_multi.sh" ]]; then
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "node0" "配置 node0"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "node1" "配置 node1"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "node2" "配置 node2"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "node3" "配置 node3"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "26657" "配置 RPC 端口"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "26656" "配置 P2P 端口"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "cometbft" "使用 CometBFT"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_multi.sh" "collect-gentxs" "收集创世交易"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: 多节点脚本不存在，跳过内容测试"
fi

# 测试 3: 停止脚本包含必要的命令
if [[ -f "${PROJECT_ROOT}/scripts/devnet_stop.sh" ]]; then
    assert_contains "${PROJECT_ROOT}/scripts/devnet_stop.sh" "kill" "包含 kill 命令"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_stop.sh" "sharetokend" "包含进程名"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: 停止脚本不存在，跳过内容测试"
fi

# 测试 4: 状态检查脚本包含必要的检查
if [[ -f "${PROJECT_ROOT}/scripts/devnet_status.sh" ]]; then
    assert_contains "${PROJECT_ROOT}/scripts/devnet_status.sh" "status" "包含状态检查"
    assert_contains "${PROJECT_ROOT}/scripts/devnet_status.sh" "26657" "检查 RPC 端口"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: 状态检查脚本不存在，跳过内容测试"
fi

echo ""
echo "--- 配置目录测试 ---"

# 测试 5: 配置目录结构
assert_dir_exists "${PROJECT_ROOT}/scripts" "存在 scripts 目录"
assert_dir_exists "${PROJECT_ROOT}/config" "存在 config 目录"

echo ""
echo "--- 网络功能测试 ---"

# 测试 6: 检查是否可以启动（如果 Ignite 已安装）
if command -v ignite &> /dev/null; then
    echo -e "${GREEN}✓ INFO${NC}: Ignite CLI 已安装，可以构建和启动网络"

    # 检查二进制文件是否存在
    if [[ -f "${PROJECT_ROOT}/bin/sharetokend" ]]; then
        echo -e "${GREEN}✓ PASS${NC}: 区块链二进制文件已构建"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠ WARN${NC}: 区块链二进制文件未构建 (运行 'make build')"
    fi
else
    echo -e "${YELLOW}⚠ INFO${NC}: Ignite CLI 未安装，跳过功能测试"
fi

echo ""
echo "=========================================="
echo "测试结果: $PASSED 通过, $FAILED 失败"
echo "=========================================="

if [[ $FAILED -gt 0 ]]; then
    echo -e "${RED}测试失败! 需要先实现缺失的脚本${NC}"
    exit 1
else
    echo -e "${GREEN}所有测试通过!${NC}"
    exit 0
fi
