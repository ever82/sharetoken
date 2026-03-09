#!/usr/bin/env bash
# 停止本地开发网络脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DATA_DIR="${PROJECT_ROOT}/.devnet"
BINARY_NAME="sharetokend"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 停止节点
stop_nodes() {
    log_info "停止开发网络节点..."

    local nodes=("node0" "node1" "node2" "node3")
    local stopped=0

    for node in "${nodes[@]}"; do
        local pid_file="${DATA_DIR}/${node}.pid"

        if [[ -f "${pid_file}" ]]; then
            local pid
            pid=$(cat "${pid_file}")

            if kill -0 "${pid}" 2>/dev/null; then
                log_info "停止 ${node} (PID: ${pid})..."
                kill "${pid}" 2>/dev/null || true
                sleep 1

                # 强制终止如果还在运行
                if kill -0 "${pid}" 2>/dev/null; then
                    log_warn "强制终止 ${node}..."
                    kill -9 "${pid}" 2>/dev/null || true
                fi

                ((stopped++))
            else
                log_warn "${node} 已经不在运行"
            fi

            rm -f "${pid_file}"
        else
            log_warn "${node} 的 PID 文件不存在"
        fi
    done

    # 清理可能残留的进程
    local remaining
    remaining=$(pgrep -f "${BINARY_NAME}.*${DATA_DIR}" 2>/dev/null || true)

    if [[ -n "${remaining}" ]]; then
        log_warn "发现残留的 ${BINARY_NAME} 进程，正在清理..."
        echo "${remaining}" | xargs kill -9 2>/dev/null || true
    fi

    log_info "已停止 ${stopped} 个节点"
}

# 清理数据（可选）
cleanup_data() {
    if [[ "$1" == "--clean" ]]; then
        log_warn "清理数据目录: ${DATA_DIR}"
        if [[ -d "${DATA_DIR}" ]]; then
            rm -rf "${DATA_DIR}"
            log_info "数据已清理"
        fi
    fi
}

# 主函数
main() {
    echo "=========================================="
    echo "停止 ShareToken 开发网络"
    echo "=========================================="
    echo ""

    stop_nodes
    cleanup_data "$1"

    echo ""
    log_info "开发网络已停止"
}

main "$@"
