#!/usr/bin/env bash
# 代码规范与 Lint 配置测试
# TDD: 先写测试，验证代码规范配置是否符合要求

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

assert_command_exists() {
    local cmd="$1"
    local desc="$2"
    if command -v "$cmd" &> /dev/null; then
        echo -e "${GREEN}✓ PASS${NC}: $desc"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠ WARN${NC}: $desc - 命令不存在: $cmd"
    fi
}

echo "=========================================="
echo "代码规范与 Lint 配置测试"
echo "=========================================="
echo ""

# 测试 1: 配置文件存在
echo "--- 配置文件测试 ---"
assert_file_exists "${PROJECT_ROOT}/.golangci.yml" "golangci-lint 配置文件"
assert_file_exists "${PROJECT_ROOT}/.editorconfig" "EditorConfig 配置文件"
assert_file_exists "${PROJECT_ROOT}/.gitattributes" "Git 属性配置文件"

echo ""
echo "--- golangci-lint 配置测试 ---"

# 测试 2: golangci-lint 配置包含必要的检查器
if [[ -f "${PROJECT_ROOT}/.golangci.yml" ]]; then
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "linters:" "配置 linters 部分"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "gofmt" "启用 gofmt 检查"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "govet" "启用 govet 检查"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "errcheck" "启用 errcheck 检查"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "staticcheck" "启用 staticcheck 检查"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "gosimple" "启用 gosimple 检查"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "ineffassign" "启用 ineffassign 检查"
    assert_contains "${PROJECT_ROOT}/.golangci.yml" "deadcode" "启用 deadcode 检查"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: golangci-lint 配置不存在，跳过内容测试"
fi

echo ""
echo "--- EditorConfig 配置测试 ---"

# 测试 3: EditorConfig 配置
if [[ -f "${PROJECT_ROOT}/.editorconfig" ]]; then
    assert_contains "${PROJECT_ROOT}/.editorconfig" "root = true" "配置 root 选项"
    assert_contains "${PROJECT_ROOT}/.editorconfig" "*.go" "配置 Go 文件规则"
    assert_contains "${PROJECT_ROOT}/.editorconfig" "indent_style" "配置缩进风格"
    assert_contains "${PROJECT_ROOT}/.editorconfig" "indent_size" "配置缩进大小"
    assert_contains "${PROJECT_ROOT}/.editorconfig" "end_of_line" "配置换行符"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: EditorConfig 配置不存在，跳过内容测试"
fi

echo ""
echo "--- Git 配置测试 ---"

# 测试 4: Git 属性配置
if [[ -f "${PROJECT_ROOT}/.gitattributes" ]]; then
    assert_contains "${PROJECT_ROOT}/.gitattributes" "*.go" "配置 Go 文件属性"
    assert_contains "${PROJECT_ROOT}/.gitattributes" "text" "配置文本属性"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: Git 属性配置不存在，跳过内容测试"
fi

echo ""
echo "--- Makefile Lint 目标测试 ---"

# 测试 5: Makefile 包含 lint 相关目标
if [[ -f "${PROJECT_ROOT}/Makefile" ]]; then
    assert_contains "${PROJECT_ROOT}/Makefile" "lint:" "Makefile 包含 lint 目标"
    assert_contains "${PROJECT_ROOT}/Makefile" "fmt:" "Makefile 包含 fmt 目标"
    assert_contains "${PROJECT_ROOT}/Makefile" "golangci-lint" "Makefile 调用 golangci-lint"
else
    echo -e "${YELLOW}⚠ SKIP${NC}: Makefile 不存在，跳过内容测试"
fi

echo ""
echo "--- 工具链测试 ---"

# 测试 6: 可选工具
assert_command_exists "golangci-lint" "golangci-lint 已安装"
assert_command_exists "goimports" "goimports 已安装"

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
