# Operations & Deployment Guide

**Version:** 2.0.0
**Updated:** 2026-03-02
**Network:** ShareTokens (STT)

---

## Overview

本文档提供 ShareTokens 区块链基础设施的部署、监控和维护指南。

**架构分层:**

```
+-----------------------------------------------------------------------------+
|                          可选插件层 (按需部署)                                 |
|  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐               │
│  │ 用户插件        │ │ 服务提供者插件   │ │ 服务提供者插件   │               │
│  │ GenieBot界面     │ │ LLM API托管      │ │ Agent/Workflow  │               │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘               │
+-----------------------------------------------------------------------------+
                              │
+-----------------------------------------------------------------------------+
|                          核心模块层 (必须有)                                  │
│  P2P通信 | 身份账号 | 钱包 | 服务市场 | 托管支付 | Trust System (MQ + Dispute)  │
+-----------------------------------------------------------------------------+
+-----------------------------------------------------------------------------+
```
```

**运维职责划分:**
- **核心节点运维**: 部署和维护区块链节点
- **插件运维**: 根据角色选择部署相应的插件

---

## Part 1: Core Node Operations (核心节点运维)

> 每个节点必须部署的核心模块

---

## 1. Node Deployment

### 1.1 Hardware Requirements

#### Minimum Requirements (Testnet)

```yaml
minimum:
  cpu: 4 cores
  memory: 16GB
  storage: 500GB SSD
  bandwidth: 100Mbps
  os: Ubuntu 20.04 LTS or later
```

#### Recommended Requirements (Mainnet)

```yaml
recommended:
  cpu: 8 cores (AMD EPYC or Intel Xeon)
  memory: 32GB DDR4
  storage: 2TB NVMe SSD
  bandwidth: 1Gbps (unmetered)
  os: Ubuntu 22.04 LTS
  redundancy:
    - RAID 1 or higher for storage
    - Dual power supplies
    - UPS backup
```

#### Cloud Instance Recommendations

| Provider  | Instance Type      | Use Case        | Cost (Est.) |
| --------- | ------------------ | --------------- | ----------- |
| AWS       | c6i.2xlarge        | Validator       | ~$350/month |
| GCP       | n2-highmem-8       | Validator       | ~$380/month |
| DigitalOcean | CPU-Optimized Droplet | Sentry Node  | ~$160/month |

### 1.2 Software Prerequisites

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install essential packages
sudo apt install -y \
  build-essential \
  git \
  curl \
  wget \
  jq \
  unzip \
  software-properties-common \
  lz4

# Install Go 1.21+
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile

# Install Ignite CLI
curl https://get.ignite.com/cli! | bash
sudo mv ignite /usr/local/bin/

# Verify installations
go version
ignite version
```

### 1.3 Core Module Configuration

核心模块是 Cosmos SDK 链的内置模块，通过 genesis 和配置文件启用。

```yaml
# ~/.sharetokens/config/app.toml
# Core Modules Configuration

# x/identity - 身份账号
identity:
  enable: true

# x/marketplace - 服务市场 (核心业务)
marketplace:
  enable: true
  service_levels: [1, 2, 3]  # LLM, Agent, Workflow

# x/escrow - 托管支付
escrow:
  enable: true
  default_timeout: 72h

# x/mq - 德商系统
mq:
  enable: true
  initial_score: 100

# x/dispute - 争议仲裁
dispute:
  enable: true
  jury_size: 5
  voting_period: 72h
```

### 1.4 Testnet Deployment

#### Step 1: Initialize Node

```bash
# Clone repository
git clone https://github.com/sharetokens/sharetokens-chain.git
cd sharetokens-chain

# Initialize chain
ignite chain init --network testnet

# Or manually initialize
sharetokensd init <moniker> --chain-id sharetokens-testnet-1
```

#### Step 2: Create Validator Key

```bash
# Create new key
sharetokensd keys add validator --keyring-backend file

# Export key backup (STORE SECURELY)
sharetokensd keys export validator --keyring-backend file > validator-key.backup
```

#### Step 3: Add Genesis Account

```bash
# Add initial tokens to validator account
sharetokensd genesis add-genesis-account validator 1000000000000stt \
  --keyring-backend file

# Add other genesis accounts (if applicable)
sharetokensd genesis add-genesis-account <address> 1000000000000stt
```

#### Step 4: Create Gentx

```bash
# Generate genesis transaction
sharetokensd genesis gentx validator 100000000stt \
  --chain-id sharetokens-testnet-1 \
  --moniker "<your-moniker>" \
  --commission-rate 0.10 \
  --commission-max-rate 0.20 \
  --commission-max-change-rate 0.01 \
  --keyring-backend file

# Collect genesis transactions
sharetokensd genesis collect-gentxs
```

#### Step 5: Configure Node

```bash
# Set persistent peers
PEERS="<peer-id-1>@<host-1>:26656,<peer-id-2>@<host-2>:26656"
sed -i "s/^persistent_peers = .*/persistent_peers = \"$PEERS\"/" ~/.sharetokens/config/config.toml

# Set seeds
SEEDS="<seed-id-1>@<seed-host-1>:26656"
sed -i "s/^seeds = .*/seeds = \"$SEEDS\"/" ~/.sharetokens/config/config.toml

# Configure pruning (recommended for storage efficiency)
sed -i 's/^pruning = .*/pruning = "custom"/' ~/.sharetokens/config/app.toml
sed -i 's/^pruning-keep-recent = .*/pruning-keep-recent = "100"/' ~/.sharetokens/config/app.toml
sed -i 's/^pruning-keep-every = .*/pruning-keep-every = "2000"/' ~/.sharetokens/config/app.toml
sed -i 's/^pruning-interval = .*/pruning-interval = "10"/' ~/.sharetokens/config/app.toml
```

#### Step 6: Start Node

```bash
# Start node (foreground)
sharetokensd start

# Or start as systemd service (recommended)
sudo tee /etc/systemd/system/sharetokensd.service << EOF
[Unit]
Description=ShareTokens Node
After=network.target

[Service]
Type=simple
User=<your-user>
ExecStart=/usr/local/bin/sharetokensd start
Restart=on-failure
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable sharetokensd
sudo systemctl start sharetokensd
```

### 1.5 Mainnet Genesis

#### Initial Token Distribution

| Allocation          | Percentage | Amount (STT)    | Vesting Period |
| ------------------- | ---------- | --------------- | -------------- |
| Community Pool      | 20%        | 200,000,000     | Unlocked       |
| Team & Advisors     | 15%        | 150,000,000     | 2 years        |
| Early Investors     | 10%        | 100,000,000     | 1 year         |
| Ecosystem Fund      | 25%        | 250,000,000     | 4 years        |
| Airdrop             | 10%        | 100,000,000     | Unlocked       |
| Liquidity Mining    | 20%        | 200,000,000     | 4 years        |

#### Governance Parameters

```yaml
governance:
  voting_period: 120h  # 5 days
  quorum: 0.334        # 33.4%
  threshold: 0.5       # 50%
  veto_threshold: 0.334 # 33.4%
  min_deposit: 100000stt

staking:
  unbonding_time: 1814400s  # 21 days
  max_validators: 100
  min_self_delegation: 1stt

slashing:
  signed_blocks_window: 10000
  min_signed_per_window: 0.05
  downtime_jail_duration: 600s
  slash_fraction_double_sign: 0.05
  slash_fraction_downtime: 0.0001
```

---

## 2. Core Module Operations

### 2.1 Service Marketplace Operations

服务市场是核心业务模块，节点运行后自动参与。

```bash
# Query registered services
sharetokensd query marketplace services --level 1  # LLM services
sharetokensd query marketplace services --level 2  # Agent services
sharetokensd query marketplace services --level 3  # Workflow services

# Query service details
sharetokensd query marketplace service <service-id>

# Query service requests
sharetokensd query marketplace requests --consumer <address>
```

### 2.2 Escrow Operations

托管支付由链自动执行，但也支持手动操作。

```bash
# Query escrow status
sharetokensd query escrow <escrow-id>

# Create escrow (usually done by marketplace)
sharetokensd tx escrow create \
  --provider <provider-address> \
  --amount 100stt \
  --service-id <service-id> \
  --from <consumer-address>

# Release escrow (after service completion)
sharetokensd tx escrow release <escrow-id> \
  --from <consumer-address>

# Create dispute
sharetokensd tx escrow dispute <escrow-id> \
  --reason "Service not delivered" \
  --from <consumer-address>
```

### 2.3 MQ Operations

德商系统自动更新，支持查询。

```bash
# Query MQ score
sharetokensd query mq score <address>

# Query MQ history
sharetokensd query mq history <address> --from 2024-01-01 --to 2024-12-31

# Query leaderboard
sharetokensd query mq leaderboard --limit 100
```

### 2.4 Dispute Operations

争议仲裁需要陪审员参与。

```bash
# Query open disputes
sharetokensd query dispute list --status open

# Query dispute details
sharetokensd query dispute <dispute-id>

# Vote on dispute (jurors only)
sharetokensd tx dispute vote <dispute-id> \
  --vote for \
  --reason "Evidence supports plaintiff" \
  --from <juror-address>
```

---

## 3. Upgrade Mechanism

### 3.1 Software Upgrade (Governance Proposal)

```bash
# Submit software upgrade proposal
sharetokensd tx gov submit-proposal software-upgrade v2.0.0 \
  --title "Upgrade to v2.0.0" \
  --description "Adding new marketplace features" \
  --upgrade-height 1000000 \
  --upgrade-info '{"binaries":{"linux/amd64":"https://github.com/sharetokens/sharetokens-chain/releases/download/v2.0.0/sharetokensd-v2.0.0-linux-amd64.tar.gz"}}' \
  --from validator \
  --deposit 10000000stt \
  --chain-id sharetokens-1 \
  --keyring-backend file \
  --fees 5000stt
```

### 3.2 Auto-Upgrade with Cosmovisor

```bash
# Install Cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

# Setup directories
mkdir -p ~/.sharetokens/cosmovisor/genesis/bin
mkdir -p ~/.sharetokens/cosmovisor/upgrades/v2.0.0/bin

# Copy current binary
cp $(which sharetokensd) ~/.sharetokens/cosmovisor/genesis/bin/

# Set environment variables
echo 'export DAEMON_NAME=sharetokensd' >> ~/.profile
echo 'export DAEMON_HOME=~/.sharetokens' >> ~/.profile
source ~/.profile

# Create systemd service for Cosmovisor
sudo tee /etc/systemd/system/cosmovisor.service << EOF
[Unit]
Description=Cosmovisor (ShareTokens)
After=network.target

[Service]
Type=simple
User=<your-user>
ExecStart=/home/<your-user>/go/bin/cosmovisor run start
Restart=on-failure
RestartSec=5
LimitNOFILE=65535
Environment="DAEMON_NAME=sharetokensd"
Environment="DAEMON_HOME=/home/<your-user>/.sharetokens"

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable cosmovisor
sudo systemctl start cosmovisor
```

---

## 4. Monitoring & Alerting

### 4.1 Key Metrics

#### Core Node Metrics

| Metric                  | Target   | Alert Threshold |
| ----------------------- | -------- | --------------- |
| Block Production Delay  | < 10s    | > 15s           |
| Transaction Confirmation| < 30s    | > 60s           |
| Peer Count              | > 25     | < 10            |
| CPU Usage               | < 70%    | > 90%           |
| Memory Usage            | < 80%    | > 95%           |
| Disk Usage              | < 80%    | > 90%           |

#### Service Marketplace Metrics

| Metric                  | Target   | Alert Threshold |
| ----------------------- | -------- | --------------- |
| Service Registration    | Working  | Failed          |
| Request Routing         | < 5s     | > 30s           |
| Escrow Operations       | Working  | Failed          |
| MQ Updates              | Working  | Stale > 1h      |

### 4.2 Prometheus Configuration

```yaml
# /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - localhost:9093

rule_files:
  - /etc/prometheus/rules/*.yml

scrape_configs:
  # ShareTokens Core Node
  - job_name: 'sharetokens-core'
    static_configs:
      - targets: ['localhost:26660']

  # Node Exporter (System Metrics)
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']

  # Prometheus self-monitoring
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### 4.3 Alert Rules

```yaml
# /etc/prometheus/rules/sharetokens-alerts.yml
groups:
  - name: sharetokens_core_alerts
    interval: 30s
    rules:
      # Block production halted
      - alert: BlockProductionHalted
        expr: increase(tendermint_block_height[1m]) == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Block production halted"
          description: "No new blocks produced for > 60 seconds"

      # Node offline
      - alert: NodeOffline
        expr: up{job="sharetokens-core"} == 0
        for: 30s
        labels:
          severity: critical
        annotations:
          summary: "Node offline"
          description: "ShareTokens core node is not responding"

      # Low peer count
      - alert: LowPeerCount
        expr: tendermint_p2p_peers < 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Low peer count"
          description: "Connected peers < 10"

      # Disk space low
      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) < 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Disk space critically low"
          description: "Available disk space < 10%"
```

---

## 5. Backup & Recovery

### 5.1 Backup Strategy

| Type       | Frequency | Retention | Size (Est.) |
| ---------- | --------- | --------- | ----------- |
| Full       | Daily     | 30 days   | ~50GB       |
| Incremental| Hourly    | 7 days    | ~1GB        |
| State Snapshots | Every 1000 blocks | 14 days | ~5GB |
| Key Backup | Once (secure storage) | Forever | <1MB |

### 5.2 Recovery Procedures

#### Scenario 1: Data Corruption

```bash
# Stop node
sudo systemctl stop sharetokensd

# Backup corrupted data
mv ~/.sharetokens/data ~/.sharetokens/data-corrupted-$(date +%Y%m%d)

# Reset node state
sharetokensd unsafe-reset-all

# Download recent snapshot
wget https://snapshots.sharetokens.io/mainnet/latest.tar.gz -O snapshot.tar.gz

# Extract snapshot
tar -xzf snapshot.tar.gz -C ~/.sharetokens/data

# Start node
sudo systemctl start sharetokensd
```

#### Scenario 2: Complete Node Failure

```bash
# On new server, install prerequisites

# Download and extract backup
wget https://backups.sharetokens.io/sharetokens-full-latest.tar.gz
tar -xzf sharetokens-full-latest.tar.gz -C ~/

# Restore validator key
sharetokensd keys import validator validator-key.backup --keyring-backend file

# Start node
sudo systemctl start sharetokensd
```

---

## 6. Security Checklist

### Core Node Security

- [ ] Firewall configured (UFW or cloud security groups)
- [ ] SSH key-based authentication only
- [ ] Fail2ban installed and configured
- [ ] Automatic security updates enabled
- [ ] Validator key encrypted and backed up offline
- [ ] Sentry node architecture for validators
- [ ] TLS enabled for RPC and API endpoints

### Network Security

```bash
# UFW Firewall Configuration
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 26656/tcp  # P2P
sudo ufw allow 26657/tcp  # RPC (restrict to localhost if possible)
sudo ufw allow 26660/tcp  # Prometheus metrics
sudo ufw enable
```

---

## Part 2: Plugin Operations (插件运维)

> 根据节点角色选择部署的可选模块

---

## 7. User Plugin: GenieBot

GenieBot是服务市场的用户端插件，提供 AI 对话界面。

### 7.1 Installation

```bash
# Using npm
npm install @sharetokens/plugin-geniebot

# Or using Docker
docker pull sharetokens/geniebot:latest
```

### 7.2 Configuration

```yaml
# geniebot-config.yml
geniebot:
  # Core node connection
  rpc_url: http://localhost:26657
  rest_url: http://localhost:1317

  # Marketplace settings
  marketplace:
    default_level: 1  # Default to LLM
    auto_recommend: true

  # UI settings
  ui:
    language: zh-CN
    theme: light

  # Wallet integration
  wallet:
    keplr_chain_id: sharetokens-1
```

### 7.3 Deployment

#### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  xiaodeng:
    image: sharetokens/xiaodeng:latest
    container_name: sharetokens-xiaodeng
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - REST_URL=http://sharetokensd:1317
      - NODE_ENV=production
    depends_on:
      - sharetokensd
    networks:
      - sharetokens-net

networks:
  sharetokens-net:
    external: true
```

#### Start Service

```bash
# Start with Docker
docker-compose up -d geniebot

# Check logs
docker-compose logs -f geniebot

# Health check
curl http://localhost:3000/health
```

### 7.4 Monitoring

```yaml
# Prometheus scrape config for Xiaodeng
scrape_configs:
  - job_name: 'geniebot'
    static_configs:
      - targets: ['localhost:3000/metrics']
```

---

## 8. Provider Plugin: LLM API Key Hosting

托管 API Key，提供 Level 1 LLM 服务。

### 8.1 Installation

```bash
# Using npm
npm install @sharetokens/plugin-llm-provider

# Or using Docker
docker pull sharetokens/llm-provider:latest
```

### 8.2 Configuration

```yaml
# llm-provider-config.yml
llm_provider:
  # Core node connection
  rpc_url: http://localhost:26657

  # Supported providers
  providers:
    - name: openai
      enabled: true
      models: [gpt-4, gpt-3.5-turbo]
    - name: anthropic
      enabled: true
      models: [claude-3-opus, claude-3-sonnet]
    - name: google
      enabled: true
      models: [gemini-pro]

  # Rate limiting
  rate_limit:
    requests_per_minute: 100
    tokens_per_minute: 100000

  # Pricing (STT per 1K tokens)
  pricing:
    gpt-4: 0.01
    gpt-3.5-turbo: 0.001
    claude-3-opus: 0.015
    gemini-pro: 0.001
```

### 8.3 API Key Management

```bash
# Register API key (encrypted)
llm-provider keys add openai --key <your-openai-api-key>

# List registered keys
llm-provider keys list

# Rotate key
llm-provider keys rotate openai --key <new-api-key>

# Revoke key
llm-provider keys revoke openai
```

### 8.4 Deployment

```yaml
# docker-compose.yml for LLM Provider
version: '3.8'

services:
  llm-provider:
    image: sharetokens/llm-provider:latest
    container_name: sharetokens-llm-provider
    restart: unless-stopped
    ports:
      - "3001:3001"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    volumes:
      - ./llm-provider-keys:/app/keys:ro
    depends_on:
      - sharetokensd
    networks:
      - sharetokens-net
```

### 8.5 Security Considerations

- API keys are encrypted at rest
- Never log API keys
- Use hardware security modules (HSM) for production
- Rotate keys regularly
- Monitor for unauthorized access

---

## 9. Provider Plugin: Agent Executor (OpenFang)

运行 AI Agent，提供 Level 2 Agent 服务。

### 9.1 Prerequisites

```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# Install OpenFang
git clone https://github.com/sharetokens/openfang.git
cd openfang
cargo build --release
cargo install --path .
```

### 9.2 Configuration

```yaml
# openfang-config.yml
openfang:
  # Core node connection
  rpc_url: http://localhost:26657

  # Agent definitions
  agents:
    - name: coder
      enabled: true
      capabilities: [code_generation, code_review, debugging]
      tools: [github, terminal, file_system]
      llm: gpt-4

    - name: researcher
      enabled: true
      capabilities: [web_search, data_analysis, report_writing]
      tools: [web_browser, calculator, document_editor]
      llm: claude-3-opus

    - name: writer
      enabled: true
      capabilities: [content_creation, editing, translation]
      tools: [document_editor, grammar_checker]
      llm: claude-3-sonnet

  # Resource limits
  limits:
    max_concurrent_tasks: 10
    max_task_duration: 1h
    memory_per_task: 512MB

  # Pricing (STT per task)
  pricing:
    coder: 30
    researcher: 20
    writer: 15
```

### 9.3 Deployment

```yaml
# docker-compose.yml for Agent Provider
version: '3.8'

services:
  agent-provider:
    image: sharetokens/agent-provider:latest
    container_name: sharetokens-agent-provider
    restart: unless-stopped
    ports:
      - "3002:3002"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
    volumes:
      - ./openfang-config.yml:/app/config.yml:ro
    depends_on:
      - sharetokensd
    networks:
      - sharetokens-net
```

### 9.4 Monitoring

```bash
# Check agent status
openfang status

# View active tasks
openfang tasks list

# View agent logs
openfang logs --agent coder --follow
```

---

## 10. Provider Plugin: Workflow Executor

编排多 Agent 任务流程，提供 Level 3 Workflow 服务。

### 10.1 Installation

```bash
# Using npm
npm install @sharetokens/plugin-workflow-provider

# Or using Docker
docker pull sharetokens/workflow-provider:latest
```

### 10.2 Configuration

```yaml
# workflow-provider-config.yml
workflow_provider:
  # Core node connection
  rpc_url: http://localhost:26657

  # Agent provider connection
  agent_provider_url: http://localhost:3002

  # Workflow definitions
  workflows:
    - name: software_development
      enabled: true
      steps:
        - name: requirements_analysis
          agent: researcher
          human_gate: true
        - name: architecture_design
          agent: coder
          human_gate: true
        - name: implementation
          agent: coder
        - name: testing
          agent: coder
        - name: deployment
          agent: coder
          human_gate: true
        - name: documentation
          agent: writer
      deliverables:
        - code_repository
        - test_report
        - deployed_application
        - technical_documentation

    - name: content_creation
      enabled: true
      steps:
        - name: topic_research
          agent: researcher
        - name: outline_design
          agent: writer
          human_gate: true
        - name: content_writing
          agent: writer
        - name: review_edit
          agent: writer
          human_gate: true
        - name: formatting
          agent: writer
      deliverables:
        - content_document
        - images
        - publish_link

  # Resource limits
  limits:
    max_concurrent_workflows: 5
    max_workflow_duration: 168h  # 7 days

  # Pricing (STT per workflow)
  pricing:
    software_development: 500
    content_creation: 100
```

### 10.3 Deployment

```yaml
# docker-compose.yml for Workflow Provider
version: '3.8'

services:
  workflow-provider:
    image: sharetokens/workflow-provider:latest
    container_name: sharetokens-workflow-provider
    restart: unless-stopped
    ports:
      - "3003:3003"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - AGENT_PROVIDER_URL=http://agent-provider:3002
    volumes:
      - ./workflow-config.yml:/app/config.yml:ro
    depends_on:
      - sharetokensd
      - agent-provider
    networks:
      - sharetokens-net
```

---

## 11. Complete Plugin Stack Deployment

### 11.1 Full Stack Docker Compose

```yaml
# docker-compose.yml - Full plugin stack
version: '3.8'

services:
  # User Plugin: GenieBot
  geniebot:
    image: sharetokens/geniebot:latest
    container_name: sharetokens-geniebot
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - REST_URL=http://sharetokensd:1317
    networks:
      - sharetokens-net

  # Provider Plugin: LLM
  llm-provider:
    image: sharetokens/llm-provider:latest
    container_name: sharetokens-llm-provider
    restart: unless-stopped
    ports:
      - "3001:3001"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    secrets:
      - openai_api_key
      - anthropic_api_key
    networks:
      - sharetokens-net

  # Provider Plugin: Agent
  agent-provider:
    image: sharetokens/agent-provider:latest
    container_name: sharetokens-agent-provider
    restart: unless-stopped
    ports:
      - "3002:3002"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - LLM_PROVIDER_URL=http://llm-provider:3001
    networks:
      - sharetokens-net

  # Provider Plugin: Workflow
  workflow-provider:
    image: sharetokens/workflow-provider:latest
    container_name: sharetokens-workflow-provider
    restart: unless-stopped
    ports:
      - "3003:3003"
    environment:
      - RPC_URL=http://sharetokensd:26657
      - AGENT_PROVIDER_URL=http://agent-provider:3002
    networks:
      - sharetokens-net

secrets:
  openai_api_key:
    file: ./secrets/openai_api_key.txt
  anthropic_api_key:
    file: ./secrets/anthropic_api_key.txt

networks:
  sharetokens-net:
    external: true
```

### 11.2 Role-Based Deployment

根据节点角色选择部署的插件：

| Node Role | Required Core | Required Plugins |
|-----------|--------------|------------------|
| Validator | All core modules | None |
| Service Provider (LLM) | All core modules | LLM Provider |
| Service Provider (Agent) | All core modules | Agent Provider |
| Service Provider (Workflow) | All core modules | Agent + Workflow Provider |
| User Node | All core modules | GenieBot |
| Full Node | All core modules | All plugins |

---

## 12. Troubleshooting

### 12.1 Core Node Issues

#### Node Not Syncing

```bash
# Check node status
sharetokensd status 2>&1 | jq .

# Check logs
journalctl -u sharetokensd -f

# Verify peers
sharetokensd query consensus peers

# Reset and resync (last resort)
sharetokensd unsafe-reset-all
sharetokensd start
```

#### Validator Jailed

```bash
# Check validator status
sharetokensd query staking validator <validator-address>

# Unjail validator
sharetokensd tx slashing unjail \
  --from validator \
  --chain-id sharetokens-1 \
  --keyring-backend file
```

### 12.2 Plugin Issues

#### GenieBot Not Connecting

```bash
# Check service status
docker-compose ps geniebot

# Check logs
docker-compose logs geniebot

# Verify core node connectivity
curl http://localhost:26657/status
```

#### LLM Provider Errors

```bash
# Check API key status
llm-provider keys list

# Test API connectivity
llm-provider test openai

# Check rate limits
llm-provider stats
```

#### Agent Provider Issues

```bash
# Check agent status
openfang status

# Restart agent
openfang restart coder

# Check resource usage
openfang resources
```

---

## 13. Support & Resources

### Official Channels

- **Documentation**: https://docs.sharetokens.io
- **GitHub**: https://github.com/sharetokens
- **Discord**: https://discord.gg/sharetokens
- **Telegram**: https://t.me/sharetokens

### Emergency Contacts

- **Security Issues**: security@sharetokens.io
- **Critical Bugs**: bugs@sharetokens.io
- **Validator Support**: validators@sharetokens.io
- **Plugin Support**: plugins@sharetokens.io

---

*Document Version: 2.0.0*
*Last Updated: 2026-03-02*
*Maintained by: ShareTokens Operations Team*
