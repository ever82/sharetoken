#!/usr/bin/env bash
# 本地开发网络状态检查脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DATA_DIR="${PROJECT_ROOT}/.devnet"

# 节点配置
NODES=("node0" "node1" "node2" "node3")
RPC_PORTS=(26657 26667 26677 26687)
P2P_PORTS=(26656 26666 26676 26686)

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检查节点状态
check_node_status() {
    local node="$1"
    local rpc_port="$2"
    local index="$3"

    local pid_file="${DATA_DIR}/${node}.pid"
    local pid=""
    local status="stopped"
    local height="N/A"
    local peers="N/A"

    # 检查 PID
    if [[ -f "${pid_file}" ]]; then
        pid=$(cat "${pid_file}" 2>/dev/null)
        if kill -0 "${pid}" 2>/dev/null; then
            status="running"
        else
            status="dead"
        fi
    fi

    # 通过 RPC 获取详细信息
    if [[ "${status}" == "running" ]]; then
        local response
        response=$(curl -s "http://127.0.0.1:${rpc_port}/status" 2>/dev/null || true)

        if [[ -n "${response}" ]]; then
            height=$(echo "${response}" | grep -o '"latest_block_height":"[0-9]*"' | cut -d'"' -f4)
            peers=$(echo "${response}" | grep -o '"n_peers":"[0-9]*"' | cut -d'"' -f4)

            [[ -z "${height}" ]] && height="syncing"
            [[ -z "${peers}" ]] && peers="0"
        fi
    fi

    # 输出状态
    local status_color="${RED}"
    if [[ "${status}" == "running" ]]; then
        status_color="${GREEN}"
    elif [[ "${status}" == "dead" ]]; then
        status_color="${YELLOW}"
    fi

    printf "  %-8s | PID: %-6s | Status: ${status_color}%-8s${NC} | Height: %-10s | Peers: %-4s | RPC: %-20s\n" \
        "${node}" "${pid:-N/A}" "${status}" "${height}" "${peers}" "127.0.0.1:${rpc_port}"
}

# 检查网络整体状态
check_network_status() {
    echo ""
    echo "=========================================="
    echo "网络整体状态"
    echo "=========================================="
    echo ""

    local running_count=0
    local total_count=${#NODES[@]}

    for i in "${!NODES[@]}"; do
        local pid_file="${DATA_DIR}/${NODES[$i]}.pid"
        if [[ -f "${pid_file}" ]]; then
            local pid
            pid=$(cat "${pid_file}" 2>/dev/null)
            if kill -0 "${pid}" 2>/dev/null; then
                ((running_count++))
            fi
        fi
    done

    echo "  运行节点: ${running_count}/${total_count}"

    if [[ ${running_count} -eq ${total_count} ]]; then
        echo -e "  网络状态: ${GREEN}健康${NC}"
    elif [[ ${running_count} -gt 0 ]]; then
        echo -e "  网络状态: ${YELLOW}部分运行${NC}"
    else
        echo -e "  网络状态: ${RED}已停止${NC}"
    fi
}

# 检查端口占用
check_ports() {
    echo ""
    echo "=========================================="
    echo "端口状态"
    echo "=========================================="
    echo ""

    for i in "${!NODES[@]}"; do
        local rpc_port="${RPC_PORTS[$i]}"
        local p2p_port="${P2P_PORTS[$i]}"

        local rpc_status="${RED}closed${NC}"
        local p2p_status="${RED}closed${NC}"

        if nc -z 127.0.0.1 "${rpc_port}" 2>/dev/null; then
            rpc_status="${GREEN}open${NC}"
        fi

        if nc -z 127.0.0.1 "${p2p_port}" 2>/dev/null; then
            p2p_status="${GREEN}open${NC}"
        fi

        printf "  %-8s | RPC: %-4s (%b) | P2P: %-4s (%b)\n" \
            "${NODES[$i]}" "${rpc_port}" "${rpc_status}" "${p2p_port}" "${p2p_status}"
    done
}

# 显示最近的区块信息
show_recent_blocks() {
    echo ""
    echo "=========================================="
    echo "最近区块信息"
    echo "=========================================="
    echo ""

    local response
    response=$(curl -s "http://127.0.0.1:26657/blockchain?minHeight=1&maxHeight=10" 2>/dev/null || true)

    if [[ -n "${response}" ]]; then
        local latest_height
        latest_height=$(echo "${response}" | grep -o '"last_height":"[0-9]*"' | head -1 | cut -d'"' -f4)
        echo "  最新区块高度: ${latest_height:-N/A}"
    else
        echo "  无法获取区块信息（网络可能未启动）"
    fi
}

# 显示日志提示
show_log_hints() {
    echo ""
    echo "=========================================="
    echo "日志查看"
    echo "=========================================="
    echo ""
    echo "  查看所有节点日志:"
    echo "    tail -f ${DATA_DIR}/*.log"
    echo ""
    echo "  查看特定节点日志:"
    for node in "${NODES[@]}"; do
        echo "    tail -f ${DATA_DIR}/${node}.log"
    done
}

# 主函数
main() {
    echo "=========================================="
    echo "ShareToken 开发网络状态"
    echo "=========================================="
    echo ""

    echo "节点详情:"
    echo "  名称     | PID      | Status   | Height       | Peers | RPC"
    echo "  ---------+----------+----------+--------------+-------+----------------------"

    for i in "${!NODES[@]}"; do
        check_node_status "${NODES[$i]}" "${RPC_PORTS[$i]}" "${i}"
    done

    check_network_status
    check_ports
    show_recent_blocks

    if [[ "$1" == "--verbose" || "$1" == "-v" ]]; then
        show_log_hints
    fi

    echo ""
    echo "提示: 使用 --verbose 或 -v 查看更多信息"
}

main "$@"
