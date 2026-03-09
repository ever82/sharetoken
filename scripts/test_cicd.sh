#!/usr/bin/env bash
# CI/CD Pipeline 配置测试
# TDD: 先写测试，验证 CI/CD 配置是否符合要求

# 不要设置 set -e，让测试完成所有检查后再退出
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
WORKFLOW_DIR="${PROJECT_ROOT}/.github/workflows"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

# 测试函数
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

assert_yaml_contains() {
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

assert_command_exists() {
    local cmd="$1"
    local desc="$2"
    if command -v "$cmd" &> /dev/null; then
        echo -e "${GREEN}✓ PASS${NC}: $desc"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $desc - 命令不存在: $cmd"
        ((FAILED++))
    fi
}

echo "=========================================="
echo "CI/CD Pipeline 配置测试"
echo "=========================================="
echo ""

# 测试 1: 工作流目录存在
echo "--- 基础结构测试 ---"
assert_file_exists "${WORKFLOW_DIR}/ci.yml" "存在 CI 工作流文件"
assert_file_exists "${WORKFLOW_DIR}/release.yml" "存在 Release 工作流文件"

echo ""
echo "--- CI 工作流内容测试 ---"

# 测试 2: CI workflow 包含必要的 job（测试实现时验证）
if [[ -f "${WORKFLOW_DIR}/ci.yml" ]]; then
    assert_yaml_contains "${WORKFLOW_DIR}/ci.yml" "test" "CI 包含 test job"
    assert_yaml_contains "${WORKFLOW_DIR}/ci.yml" "build" "CI 包含 build job"
    assert_yaml_contains "${WORKFLOW_DIR}/ci.yml" "lint" "CI 包含 lint job"
    assert_yaml_contains "${WORKFLOW_DIR}/ci.yml" "go test" "CI 运行 Go 测试"
    assert_yaml_contains "${WORKFLOW_DIR}/ci.yml" "go build" "CI 运行 Go 构建"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: CI 工作流不存在，跳过内容测试"
fi

echo ""
echo "--- 工具链测试 ---"

# 测试 3: 必要的命令行工具
assert_command_exists "go" "Go 已安装"
assert_command_exists "git" "Git 已安装"

# 测试 4: Go 版本检查
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | head -1)
    echo -e "${GREEN}✓ INFO${NC}: Go 版本: $GO_VERSION"
fi

echo ""
echo "--- Makefile 测试 ---"

# 测试 5: Makefile 存在并包含必要目标
assert_file_exists "${PROJECT_ROOT}/Makefile" "存在 Makefile"
if [[ -f "${PROJECT_ROOT}/Makefile" ]]; then
    assert_yaml_contains "${PROJECT_ROOT}/Makefile" "test:" "Makefile 包含 test 目标"
    assert_yaml_contains "${PROJECT_ROOT}/Makefile" "build:" "Makefile 包含 build 目标"
    assert_yaml_contains "${PROJECT_ROOT}/Makefile" "lint:" "Makefile 包含 lint 目标"
fi

echo ""
echo "=========================================="
echo "测试结果: $PASSED 通过, $FAILED 失败"
echo "=========================================="

if [[ $FAILED -gt 0 ]]; then
    echo -e "${RED}测试失败! 需要先实现缺失的配置${NC}"
    exit 1
else
    echo -e "${GREEN}所有测试通过!${NC}"
    exit 0
fi
