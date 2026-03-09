#!/usr/bin/env bash
# 多节点本地开发网络启动脚本
# 启动 4 个节点的 Cosmos SDK + CometBFT (cometbft) 开发网络

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
CONFIG_DIR="${PROJECT_ROOT}/config"
DATA_DIR="${PROJECT_ROOT}/.devnet"
BINARY_NAME="sharetokend"
CHAIN_ID="sharetoken-devnet"

# 节点配置
NODES=("node0" "node1" "node2" "node3")
RPC_PORTS=(26657 26667 26677 26687)
P2P_PORTS=(26656 26666 26676 26686)
GRPC_PORTS=(9090 9091 9092 9093)
API_PORTS=(1317 1318 1319 1320)
PPROF_PORTS=(6060 6061 6062 6063)

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

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go 1.19+"
        exit 1
    fi

    # 检查二进制文件
    if [[ ! -f "${PROJECT_ROOT}/bin/${BINARY_NAME}" ]]; then
        log_warn "二进制文件不存在，尝试构建..."
        if [[ -f "${PROJECT_ROOT}/Makefile" ]]; then
            (cd "${PROJECT_ROOT}" && make build)
        else
            log_error "Makefile 不存在，无法自动构建"
            exit 1
        fi
    fi
}

# 清理旧数据
cleanup() {
    log_info "清理旧数据..."
    if [[ -d "${DATA_DIR}" ]]; then
        rm -rf "${DATA_DIR}"
    fi
    mkdir -p "${DATA_DIR}"
}

# 初始化节点配置
init_nodes() {
    log_info "初始化节点配置..."

    local first_node="${NODES[0]}"
    local first_node_dir="${DATA_DIR}/${first_node}"

    # 初始化第一个节点（创世节点）
    "${PROJECT_ROOT}/bin/${BINARY_NAME}" init "${first_node}" --chain-id "${CHAIN_ID}" --home "${first_node_dir}" 2>/dev/null || true

    # 创建其他节点目录
    for i in "${!NODES[@]}"; do
        local node="${NODES[$i]}"
        local node_dir="${DATA_DIR}/${node}"

        if [[ "${node}" != "${first_node}" ]]; then
            mkdir -p "${node_dir}"
            # 初始化节点（创建独立的node_key和priv_validator_key）
            "${PROJECT_ROOT}/bin/${BINARY_NAME}" init "${node}" --chain-id "${CHAIN_ID}" --home "${node_dir}" 2>/dev/null || true
            # 复制node0的创世配置（除了节点密钥）
            cp "${first_node_dir}/config/genesis.json" "${node_dir}/config/"
            cp "${first_node_dir}/config/app.toml" "${node_dir}/config/"
            # 保留此节点自己的config.toml中的节点ID，只修改端口
        fi

        # 配置节点端口
        configure_node "${node}" "${i}"
    done
}

# 配置单个节点
configure_node() {
    local node="$1"
    local index="$2"
    local node_dir="${DATA_DIR}/${node}"
    local config_file="${node_dir}/config/config.toml"
    local app_file="${node_dir}/config/app.toml"

    log_info "配置 ${node}..."

    # 配置 CometBFT (config.toml)
    if [[ -f "${config_file}" ]]; then
        # 设置 P2P 端口
        sed -i.bak "s/laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/0.0.0.0:${P2P_PORTS[$index]}\"/" "${config_file}"
        # 设置 RPC 端口
        sed -i.bak "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/127.0.0.1:${RPC_PORTS[$index]}\"/" "${config_file}"
        # 设置 pprof 端口
        sed -i.bak "s/pprof_laddr = \"localhost:6060\"/pprof_laddr = \"localhost:${PPROF_PORTS[$index]}\"/" "${config_file}"
        # 设置节点名称
        sed -i.bak "s/moniker = \".*\"/moniker = \"${node}\"/" "${config_file}"
        # 禁用 UPnP（开发环境）
        sed -i.bak "s/upnp = true/upnp = false/" "${config_file}"
        # 设置出块时间
        sed -i.bak "s/timeout_commit = \"5s\"/timeout_commit = \"2s\"/" "${config_file}"
        # 设置空块间隔
        sed -i.bak "s/create_empty_blocks_interval = \"0s\"/create_empty_blocks_interval = \"2s\"/" "${config_file}"

        rm -f "${config_file}.bak"
    fi

    # 配置 App (app.toml)
    if [[ -f "${app_file}" ]]; then
        # 设置 gRPC 端口
        sed -i.bak "s/address = \"0.0.0.0:9090\"/address = \"0.0.0.0:${GRPC_PORTS[$index]}\"/" "${app_file}"
        # 设置 API 端口 (注意格式是 tcp://localhost:1317)
        sed -i.bak "s|address = \"tcp://localhost:1317\"|address = \"tcp://localhost:${API_PORTS[$index]}\"|" "${app_file}"
        # 启用 API (修改 [api] 部分下的 enable)
        sed -i.bak '/\[api\]/,/^\[/ s/enable = false/enable = true/' "${app_file}"

        rm -f "${app_file}.bak"
    fi
}

# 生成密钥和地址
setup_keys() {
    log_info "设置节点密钥..."

    for i in "${!NODES[@]}"; do
        local node="${NODES[$i]}"
        local node_dir="${DATA_DIR}/${node}"

        # 生成密钥（使用测试助记词）
        local mnemonic="mnemonic"
        echo "${mnemonic}" | "${PROJECT_ROOT}/bin/${BINARY_NAME}" keys add "validator${i}" --recover --home "${node_dir}" --keyring-backend test 2>/dev/null || \
        "${PROJECT_ROOT}/bin/${BINARY_NAME}" keys add "validator${i}" --home "${node_dir}" --keyring-backend test 2>/dev/null || true
    done
}

# 配置创世文件
setup_genesis() {
    log_info "配置创世文件..."

    local first_node="${NODES[0]}"
    local first_node_dir="${DATA_DIR}/${first_node}"
    local genesis_file="${DATA_DIR}/${first_node}/config/genesis.json"

    # 修改创世参数
    if [[ -f "${genesis_file}" ]]; then
        # 设置出块时间参数（如果jq存在）
        if command -v jq &> /dev/null; then
            jq '.consensus.params.block.max_bytes = "22020096"' "${genesis_file}" > "${genesis_file}.tmp" && mv "${genesis_file}.tmp" "${genesis_file}"
            jq '.consensus.params.block.max_gas = "-1"' "${genesis_file}" > "${genesis_file}.tmp" && mv "${genesis_file}.tmp" "${genesis_file}"
        else
            log_warn "jq未安装，跳过创世参数修改"
        fi
    fi

    # 为所有节点添加创世账户和创建验证人
    for i in "${!NODES[@]}"; do
        local node="${NODES[$i]}"
        local node_dir="${DATA_DIR}/${node}"

        # 添加账户到创世文件（使用第一个节点的home，因为创世文件在那里）
        "${PROJECT_ROOT}/bin/${BINARY_NAME}" add-genesis-account "validator${i}" 1000000000stake --home "${DATA_DIR}/${first_node}" --keyring-backend test 2>/dev/null || true
    done

    # 为每个节点创建创世交易（gentx）
    for i in "${!NODES[@]}"; do
        "${PROJECT_ROOT}/bin/${BINARY_NAME}" gentx "validator${i}" 100000000stake \
            --chain-id "${CHAIN_ID}" \
            --home "${DATA_DIR}/${first_node}" \
            --keyring-backend test 2>/dev/null || true
    done

    # 收集创世交易
    "${PROJECT_ROOT}/bin/${BINARY_NAME}" collect-gentxs \
        --home "${DATA_DIR}/${first_node}" 2>/dev/null || true

    # 验证创世文件
    "${PROJECT_ROOT}/bin/${BINARY_NAME}" validate-genesis \
        --home "${DATA_DIR}/${first_node}" 2>/dev/null || true

    # 将创世文件复制到其他节点
    for i in "${!NODES[@]}"; do
        local node="${NODES[$i]}"
        if [[ "${node}" != "${first_node}" ]]; then
            local node_dir="${DATA_DIR}/${node}"
            cp "${DATA_DIR}/${first_node}/config/genesis.json" "${node_dir}/config/genesis.json"
        fi
    done
}

# 配置节点间连接
setup_peers() {
    log_info "配置节点间连接..."

    local first_node="${NODES[0]}"
    local first_node_dir="${DATA_DIR}/${first_node}"

    # 获取第一个节点的 node ID
    local node0_id
    node0_id=$("${PROJECT_ROOT}/bin/${BINARY_NAME}" tendermint show-node-id --home "${first_node_dir}" 2>/dev/null)

    if [[ -z "${node0_id}" ]]; then
        log_warn "无法获取 node0 ID，使用默认配置"
        return
    fi

    # 配置其他节点连接到 node0
    for i in "${!NODES[@]}"; do
        if [[ $i -eq 0 ]]; then
            continue
        fi

        local node="${NODES[$i]}"
        local node_dir="${DATA_DIR}/${node}"
        local config_file="${node_dir}/config/config.toml"

        # 设置 persistent_peers
        local peer="${node0_id}@127.0.0.1:${P2P_PORTS[0]}"
        sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"${peer}\"/" "${config_file}"
        rm -f "${config_file}.bak"
    done
}

# 启动节点
start_nodes() {
    log_info "启动节点..."

    for i in "${!NODES[@]}"; do
        local node="${NODES[$i]}"
        local node_dir="${DATA_DIR}/${node}"
        local log_file="${DATA_DIR}/${node}.log"

        log_info "启动 ${node} (RPC: ${RPC_PORTS[$i]}, P2P: ${P2P_PORTS[$i]})..."

        # 启动节点
        nohup "${PROJECT_ROOT}/bin/${BINARY_NAME}" start \
            --home "${node_dir}" \
            > "${log_file}" 2>&1 &

        # 保存 PID
        echo $! > "${DATA_DIR}/${node}.pid"
    done
}

# 等待节点启动
wait_for_startup() {
    log_info "等待节点启动..."

    local max_attempts=30
    local attempt=0

    while [[ $attempt -lt $max_attempts ]]; do
        local all_ready=true

        for i in "${!NODES[@]}"; do
            local rpc_port="${RPC_PORTS[$i]}"

            # 检查 RPC 端口
            if ! curl -s "http://127.0.0.1:${rpc_port}/status" > /dev/null 2>&1; then
                all_ready=false
                break
            fi
        done

        if [[ "$all_ready" == "true" ]]; then
            log_info "所有节点已启动!"
            return 0
        fi

        sleep 1
        ((attempt++))
        echo -n "."
    done

    log_error "节点启动超时"
    return 1
}

# 显示状态信息
show_status() {
    echo ""
    echo "=========================================="
    echo "开发网络已启动!"
    echo "=========================================="
    echo ""
    echo "节点信息:"
    echo "  Chain ID: ${CHAIN_ID}"
    echo ""

    for i in "${!NODES[@]}"; do
        local node="${NODES[$i]}"
        local pid_file="${DATA_DIR}/${node}.pid"
        local pid=""

        if [[ -f "${pid_file}" ]]; then
            pid=$(cat "${pid_file}")
        fi

        echo "  ${node}:"
        echo "    PID: ${pid}"
        echo "    RPC: http://127.0.0.1:${RPC_PORTS[$i]}"
        echo "    P2P: 127.0.0.1:${P2P_PORTS[$i]}"
        echo "    gRPC: 127.0.0.1:${GRPC_PORTS[$i]}"
        echo "    API: http://127.0.0.1:${API_PORTS[$i]}"
        echo "    Log: ${DATA_DIR}/${node}.log"
        echo ""
    done

    echo "命令示例:"
    echo "  查看状态: ./scripts/devnet_status.sh"
    echo "  停止网络: ./scripts/devnet_stop.sh"
    echo "  查看日志: tail -f ${DATA_DIR}/node0.log"
    echo ""
}

# 主函数
main() {
    echo "=========================================="
    echo "ShareToken 多节点开发网络启动脚本"
    echo "=========================================="
    echo ""

    check_dependencies
    cleanup
    init_nodes
    setup_keys
    setup_genesis
    setup_peers
    start_nodes
    wait_for_startup
    show_status
}

# 处理信号
cleanup_on_exit() {
    log_warn "接收到中断信号，正在停止节点..."
    if [[ -f "${SCRIPT_DIR}/devnet_stop.sh" ]]; then
        bash "${SCRIPT_DIR}/devnet_stop.sh"
    fi
    exit 0
}

trap cleanup_on_exit SIGINT SIGTERM

main "$@"
