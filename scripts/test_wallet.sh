#!/usr/bin/env bash
# 钱包与代币系统测试 (ACH-DEV-003)
# TDD: 验证钱包功能是否符合验收标准

set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
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
echo "钱包与代币系统测试 (ACH-DEV-003)"
echo "=========================================="
echo ""

# 测试 1: STT 代币定义
echo "--- STT 代币定义测试 ---"

# 检查 genesis.json 中是否有 stake 代币（默认）
if [[ -f "${PROJECT_ROOT}/config/genesis.json" ]] || [[ -f "${PROJECT_ROOT}/.devnet/node0/config/genesis.json" ]]; then
    GENESIS_FILE="${PROJECT_ROOT}/config/genesis.json"
    if [[ -f "${PROJECT_ROOT}/.devnet/node0/config/genesis.json" ]]; then
        GENESIS_FILE="${PROJECT_ROOT}/.devnet/node0/config/genesis.json"
    fi

    # 检查是否有代币定义
    if grep -q '"denom"' "$GENESIS_FILE" 2>/dev/null; then
        assert_pass "创世文件中包含代币定义"
    else
        assert_fail "代币定义" "创世文件中未找到 denom 定义"
    fi

    # 检查是否支持自定义代币
    if grep -q "stake\|token\|STT\|stt" "$GENESIS_FILE" 2>/dev/null; then
        assert_pass "创世文件包含代币名称"
    else
        assert_fail "代币名称" "未找到代币名称"
    fi
else
    assert_skip "代币定义测试 (需要先启动 devnet)"
fi

# 测试 2: 余额查询接口
echo ""
echo "--- 余额查询接口测试 ---"

# 检查 CLI 查询命令
if [[ -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
    # 检查是否支持 bank 查询
    if "${PROJECT_ROOT}/bin/${BINARY_NAME}" query bank --help 2>&1 | grep -q "balance"; then
        assert_pass "CLI 支持 bank balance 查询"
    else
        assert_fail "bank balance 查询" "CLI 不支持"
    fi
else
    assert_skip "余额查询接口测试 (需要先构建)"
fi

# 测试 3: 转账交易
echo ""
echo "--- 转账交易测试 ---"

if [[ -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
    # 检查是否支持 send 交易
    if "${PROJECT_ROOT}/bin/${BINARY_NAME}" tx bank --help 2>&1 | grep -q "send"; then
        assert_pass "CLI 支持 bank send 转账"
    else
        assert_fail "bank send 转账" "CLI 不支持"
    fi
else
    assert_skip "转账交易测试 (需要先构建)"
fi

# 测试 4: Keplr 钱包集成
echo ""
echo "--- Keplr 钱包集成测试 ---"

# 检查是否有前端目录
if [[ -d "${PROJECT_ROOT}/vue" ]] || [[ -d "${PROJECT_ROOT}/frontend" ]]; then
    FRONTEND_DIR=""
    if [[ -d "${PROJECT_ROOT}/vue" ]]; then
        FRONTEND_DIR="${PROJECT_ROOT}/vue"
    elif [[ -d "${PROJECT_ROOT}/frontend" ]]; then
        FRONTEND_DIR="${PROJECT_ROOT}/frontend"
    fi

    if [[ -n "$FRONTEND_DIR" ]]; then
        assert_pass "前端目录存在"

        # 检查 Keplr 相关代码
        if find "$FRONTEND_DIR" -name "*.ts" -o -name "*.js" -o -name "*.vue" 2>/dev/null | head -1 | grep -q "."; then
            # 检查是否提到 Keplr
            if grep -r "keplr\|Keplr" "$FRONTEND_DIR" --include="*.ts" --include="*.js" --include="*.vue" 2>/dev/null | head -1 | grep -q "keplr"; then
                assert_pass "前端代码包含 Keplr 集成"
            else
                assert_fail "Keplr 集成" "前端代码未找到 Keplr 相关代码"
            fi
        else
            assert_fail "Keplr 集成" "前端目录为空或无代码文件"
        fi
    fi
else
    assert_fail "Keplr 钱包集成" "前端目录不存在 (vue 或 frontend)"
fi

# 测试 5: WalletConnect 支持
echo ""
echo "--- WalletConnect 支持测试 ---"

if [[ -n "$FRONTEND_DIR" ]] && [[ -d "$FRONTEND_DIR" ]]; then
    if grep -r "walletconnect\|WalletConnect" "$FRONTEND_DIR" --include="*.ts" --include="*.js" --include="*.vue" 2>/dev/null | head -1 | grep -q "walletconnect"; then
        assert_pass "前端代码包含 WalletConnect 集成"
    else
        assert_fail "WalletConnect 集成" "前端代码未找到 WalletConnect 相关代码"
    fi
else
    assert_fail "WalletConnect 支持" "前端目录不存在"
fi

# 测试 6: 交易历史查询
echo ""
echo "--- 交易历史查询测试 ---"

if [[ -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
    # 检查是否支持 txs 查询
    if "${PROJECT_ROOT}/bin/${BINARY_NAME}" query txs --help 2>&1 | grep -q "events"; then
        assert_pass "CLI 支持交易历史查询"
    else
        assert_fail "交易历史查询" "CLI 不支持 txs 查询"
    fi
else
    assert_skip "交易历史查询测试 (需要先构建)"
fi

# 测试 7: 地址前缀配置
echo ""
echo "--- 地址前缀配置测试 ---"

if grep -q "sharetoken" "${PROJECT_ROOT}/app/app.go" 2>/dev/null; then
    assert_pass "地址前缀配置为 'sharetoken'"
else
    assert_fail "地址前缀配置" "未找到 sharetoken 前缀配置"
fi

# 测试 8: Bank 模块集成
echo ""
echo "--- Bank 模块集成测试 ---"

if grep -q "bank.AppModule" "${PROJECT_ROOT}/app/app.go" 2>/dev/null; then
    assert_pass "Bank 模块已集成"
else
    assert_fail "Bank 模块" "未在 app.go 中找到 Bank 模块"
fi

# 测试 9: 密钥管理
echo ""
echo "--- 密钥管理测试 ---"

if [[ -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
    # 检查是否支持 keys 命令
    if "${PROJECT_ROOT}/bin/${BINARY_NAME}" keys --help 2>&1 | grep -q "add"; then
        assert_pass "CLI 支持密钥管理 (keys add)"
    else
        assert_fail "密钥管理" "CLI 不支持 keys 命令"
    fi
else
    assert_skip "密钥管理测试 (需要先构建)"
fi

# 测试 10: 助记词支持
echo ""
echo "--- 助记词支持测试 ---"

if [[ -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
    if "${PROJECT_ROOT}/bin/${BINARY_NAME}" keys add --help 2>&1 | grep -q "recover"; then
        assert_pass "CLI 支持助记词恢复"
    else
        assert_fail "助记词恢复" "CLI 不支持 --recover 参数"
    fi
else
    assert_skip "助记词支持测试 (需要先构建)"
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
    echo -e "${YELLOW}部分测试跳过，请先构建项目${NC}"
    exit 0
else
    echo -e "${GREEN}所有测试通过!${NC}"
    exit 0
fi
