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

### 1.4 Mainnet Deployment

#### Pre-Deployment Checklist

- [ ] Hardware meets recommended requirements
- [ ] Validator key generated and backed up securely
- [ ] Genesis file downloaded from official source
- [ ] Network configuration verified
- [ ] Monitoring stack ready
- [ ] Backup strategy in place

#### Step 1: Initialize Node

```bash
# Clone repository
git clone https://github.com/sharetokens/sharetokens-chain.git
cd sharetokens-chain

# Checkout mainnet release tag
git checkout v2.0.0

# Build binary
make build

# Initialize chain for mainnet
sharetokensd init <moniker> --chain-id sharetokens-1
```

#### Step 2: Download Mainnet Genesis

```bash
# Download official genesis
curl -L https://raw.githubusercontent.com/sharetokens/mainnet/main/genesis.json \
  -o ~/.sharetokens/config/genesis.json

# Verify genesis hash
sha256sum ~/.sharetokens/config/genesis.json
# Expected: <official-genesis-hash>

# Download mainnet config
curl -L https://raw.githubusercontent.com/sharetokens/mainnet/main/config.toml \
  -o ~/.sharetokens/config/config.toml

curl -L https://raw.githubusercontent.com/sharetokens/mainnet/main/app.toml \
  -o ~/.sharetokens/config/app.toml
```

#### Step 3: Configure Node for Mainnet

```bash
# Set persistent peers from mainnet seed list
SEEDS="seed1@seed1.sharetokens.io:26656,seed2@seed2.sharetokens.io:26656"
sed -i "s/^seeds = .*/seeds = \"$SEEDS\"/" ~/.sharetokens/config/config.toml

# Configure minimum gas prices
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.025stt"/' ~/.sharetokens/config/app.toml

# Enable Prometheus metrics
sed -i 's/^prometheus = .*/prometheus = true/' ~/.sharetokens/config/config.toml

# Configure pruning for mainnet (aggressive)
sed -i 's/^pruning = .*/pruning = "everything"/' ~/.sharetokens/config/app.toml
```

#### Step 4: Start Syncing

```bash
# Option A: Start from genesis (slow, historical data)
sharetokensd start

# Option B: Start from snapshot (fast, recommended)
# Download snapshot
wget https://snapshots.sharetokens.io/mainnet/latest.tar.lz4

# Stop node if running
sudo systemctl stop sharetokensd

# Reset and restore from snapshot
sharetokensd unsafe-reset-all
lz4 -d latest.tar.lz4 | tar -x -C ~/.sharetokens/

# Start node
sudo systemctl start sharetokensd
```

#### Step 5: Create Validator (Optional)

```bash
# Create validator key
sharetokensd keys add validator --keyring-backend file

# Fund validator address (from exchange or faucet)
# Wait for funds to arrive

# Create validator
current_height=$(sharetokensd status | jq -r '.SyncInfo.latest_block_height')
pubkey=$(sharetokensd tendermint show-validator)

sharetokensd tx staking create-validator \
  --amount=1000000stt \
  --pubkey="$pubkey" \
  --moniker="<your-validator-name>" \
  --website="https://your-website.com" \
  --details="Your validator description" \
  --commission-rate=0.10 \
  --commission-max-rate=0.20 \
  --commission-max-change-rate=0.01 \
  --min-self-delegation=1 \
  --from=validator \
  --chain-id=sharetokens-1 \
  --gas=auto \
  --gas-adjustment=1.2 \
  --gas-prices=0.025stt \
  --keyring-backend=file
```

---

### 1.5 Testnet Deployment

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

---

## 2. Mainnet Genesis

### Initial Token Distribution

| Allocation          | Percentage | Amount (STT)    | Vesting Period |
| ------------------- | ---------- | --------------- | -------------- |
| Community Pool      | 20%        | 200,000,000     | Unlocked       |
| Team & Advisors     | 15%        | 150,000,000     | 2 years        |
| Early Investors     | 10%        | 100,000,000     | 1 year         |
| Ecosystem Fund      | 25%        | 250,000,000     | 4 years        |
| Airdrop             | 10%        | 100,000,000     | Unlocked       |
| Liquidity Mining    | 20%        | 200,000,000     | 4 years        |

### Governance Parameters

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

## 3. Core Module Operations

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

## 4. Upgrade Mechanism

### 4.1 Chain Upgrade Types

#### Planned Upgrades (via Governance)
- Feature additions
- Parameter changes
- Consensus updates
- Security patches

#### Emergency Upgrades
- Critical security fixes
- Consensus bugs
- Requires coordination via validators channel

### 4.2 Software Upgrade (Governance Proposal)

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

# Vote on proposal
sharetokensd tx gov vote <proposal-id> yes \
  --from validator \
  --chain-id sharetokens-1 \
  --keyring-backend file

# Monitor voting progress
sharetokensd query gov proposal <proposal-id>
```

### 4.3 Cosmovisor Setup (Recommended)

Cosmovisor is a process manager for Cosmos SDK applications that automates chain upgrades.

#### Installation

```bash
# Install Cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

# Setup directories
mkdir -p ~/.sharetokens/cosmovisor/genesis/bin
mkdir -p ~/.sharetokens/cosmovisor/upgrades

# Copy current binary to genesis
cp $(which sharetokensd) ~/.sharetokens/cosmovisor/genesis/bin/

# Set environment variables
echo 'export DAEMON_NAME=sharetokensd' >> ~/.profile
echo 'export DAEMON_HOME=$HOME/.sharetokens' >> ~/.profile
echo 'export DAEMON_ALLOW_DOWNLOAD_BINARIES=false' >> ~/.profile
echo 'export DAEMON_LOG_BUFFER_SIZE=512' >> ~/.profile
echo 'export DAEMON_RESTART_AFTER_UPGRADE=true' >> ~/.profile
source ~/.profile

# Verify installation
cosmovisor version
```

#### Systemd Service Configuration

```bash
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
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable cosmovisor
sudo systemctl start cosmovisor

# Check status
sudo systemctl status cosmovisor
journalctl -u cosmovisor -f
```

### 4.4 Upgrade Process with Cosmovisor

#### Phase 1: Pre-Upgrade (Before Target Height)

```bash
# Download new binary
wget https://github.com/sharetokens/sharetokens-chain/releases/download/v2.0.0/sharetokensd-v2.0.0-linux-amd64.tar.gz
tar -xzf sharetokensd-v2.0.0-linux-amd64.tar.gz

# Verify checksum
sha256sum sharetokensd
cat checksum.txt | grep sharetokensd

# Create upgrade directory
mkdir -p ~/.sharetokens/cosmovisor/upgrades/v2.0.0/bin

# Copy new binary
cp sharetokensd ~/.sharetokens/cosmovisor/upgrades/v2.0.0/bin/

# Verify binary version
~/.sharetokens/cosmovisor/upgrades/v2.0.0/bin/sharetokensd version

# Set upgrade info (optional - for automated downloads)
echo '{"binaries":{"linux/amd64":"https://github.com/sharetokens/sharetokens-chain/releases/download/v2.0.0/sharetokensd-v2.0.0-linux-amd64.tar.gz"}}' > ~/.sharetokens/cosmovisor/upgrades/v2.0.0/upgrade-info.json
```

#### Phase 2: Monitor Upgrade Height

```bash
# Check current height
watch -n 5 'sharetokensd status | jq .SyncInfo.latest_block_height'

# Check upgrade plan
sharetokensd query upgrade plan

# Monitor logs
journalctl -u cosmovisor -f
```

#### Phase 3: Post-Upgrade Verification

```bash
# Check new version
sharetokensd version

# Check node is syncing
sharetokensd status | jq .SyncInfo

# Check validator is signing (if validator)
curl -s http://localhost:26657/block | jq .result.block.header.height

# Verify chain continues
watch -n 2 'curl -s http://localhost:26657/block | jq .result.block.header.height'
```

### 4.5 Manual Upgrade (Without Cosmovisor)

```bash
# Monitor for upgrade height
# When halt-height reached:

# Stop node
sudo systemctl stop sharetokensd

# Backup (just in case)
cp -r ~/.sharetokens/data ~/.sharetokens/data-backup-$(date +%Y%m%d)

# Download and install new binary
wget <new-binary-url>
tar -xzf <new-binary.tar.gz>
sudo cp sharetokensd /usr/local/bin/
sudo chmod +x /usr/local/bin/sharetokensd

# Verify
sharetokensd version

# Start node
sudo systemctl start sharetokensd

# Check status
sharetokensd status
```

### 4.6 Rollback Procedure

If upgrade fails:

```bash
# Stop node
sudo systemctl stop sharetokensd

# Restore from backup (if using manual upgrade)
cp -r ~/.sharetokens/data-backup-<date> ~/.sharetokens/data

# Restore old binary
sudo cp /path/to/backup/sharetokensd /usr/local/bin/

# Start with old version
sudo systemctl start sharetokensd
```

### 4.7 Emergency Upgrade Coordination

When emergency upgrade is required:

1. **Immediate Actions:**
   - Stop node if instructed
   - Join emergency validator call
   - Monitor official channels

2. **Download Patch:**
   ```bash
   # Download emergency patch
   wget <emergency-binary-url>

   # Verify signature
   gpg --verify <binary>.sig <binary>

   # Install
   sudo systemctl stop sharetokensd
   sudo cp <new-binary> /usr/local/bin/sharetokensd
   sudo systemctl start sharetokensd
   ```

3. **Communication Channels:**
   - Validator Discord: #emergency
   - Telegram: @STT_Validators
   - Email: validators@sharetokens.io

---

## 5. Monitoring & Alerting

### 5.1 Key Metrics

#### Core Node Metrics

| Metric                  | Target   | Alert Threshold | Severity |
| ----------------------- | -------- | --------------- | -------- |
| Block Production Delay  | < 10s    | > 15s           | Critical |
| Transaction Confirmation| < 30s    | > 60s           | Warning  |
| Peer Count              | > 25     | < 10            | Warning  |
| CPU Usage               | < 70%    | > 90%           | Warning  |
| Memory Usage            | < 80%    | > 95%           | Critical |
| Disk Usage              | < 80%    | > 90%           | Critical |
| Block Height Lag        | < 5      | > 100           | Warning  |
| Missed Blocks (24h)     | < 50     | > 200           | Critical |

#### Service Marketplace Metrics

| Metric                  | Target   | Alert Threshold | Severity |
| ----------------------- | -------- | --------------- | -------- |
| Service Registration    | Working  | Failed          | Warning  |
| Request Routing         | < 5s     | > 30s           | Warning  |
| Escrow Operations       | Working  | Failed          | Critical |
| MQ Updates              | Working  | Stale > 1h      | Warning  |

### 5.2 Prometheus Setup

#### Installation

```bash
# Download Prometheus
cd /tmp
wget https://github.com/prometheus/prometheus/releases/download/v2.47.0/prometheus-2.47.0.linux-amd64.tar.gz
tar -xzf prometheus-2.47.0.linux-amd64.tar.gz
cd prometheus-2.47.0.linux-amd64

# Move binaries
sudo mv prometheus /usr/local/bin/
sudo mv promtool /usr/local/bin/

# Create directories
sudo mkdir -p /etc/prometheus
sudo mkdir -p /var/lib/prometheus

# Copy config files
sudo cp -r consoles /etc/prometheus/
sudo cp -r console_libraries /etc/prometheus/
```

#### Configuration

```yaml
# /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    monitor: 'sharetokens-monitor'
    chain: 'sharetokens-1'

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - localhost:9093
      timeout: 10s
      api_version: v2

rule_files:
  - /etc/prometheus/rules/*.yml

scrape_configs:
  # ShareTokens Core Node
  - job_name: 'sharetokens-core'
    static_configs:
      - targets: ['localhost:26660']
    metrics_path: /metrics
    scrape_interval: 5s

  # Node Exporter (System Metrics)
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
    scrape_interval: 15s

  # Prometheus self-monitoring
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  # Validator metrics (if applicable)
  - job_name: 'validator-metrics'
    static_configs:
      - targets: ['localhost:26661']
    scrape_interval: 5s
```

### 5.3 Grafana Setup

#### Installation

```bash
# Add Grafana repository
sudo apt-get install -y software-properties-common
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee /etc/apt/sources.list.d/grafana.list

# Install Grafana
sudo apt-get update
sudo apt-get install -y grafana

# Start Grafana
sudo systemctl daemon-reload
sudo systemctl enable grafana-server
sudo systemctl start grafana-server

# Default credentials: admin/admin
# Access: http://<server-ip>:3000
```

#### Dashboard Configuration

```json
{
  "dashboard": {
    "title": "ShareTokens Node Monitor",
    "tags": ["sharetokens", "blockchain"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Block Height",
        "type": "stat",
        "targets": [
          {
            "expr": "tendermint_block_height",
            "legendFormat": "Current Height"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "thresholds": {
              "steps": [
                {"value": null, "color": "green"}
              ]
            }
          }
        }
      },
      {
        "id": 2,
        "title": "Block Time",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(tendermint_block_height[1m])",
            "legendFormat": "Blocks/sec"
          }
        ]
      },
      {
        "id": 3,
        "title": "Peer Count",
        "type": "stat",
        "targets": [
          {
            "expr": "tendermint_p2p_peers",
            "legendFormat": "Peers"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "thresholds": {
              "steps": [
                {"value": null, "color": "red"},
                {"value": 5, "color": "yellow"},
                {"value": 25, "color": "green"}
              ]
            }
          }
        }
      },
      {
        "id": 4,
        "title": "CPU Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)",
            "legendFormat": "CPU %"
          }
        ]
      },
      {
        "id": 5,
        "title": "Memory Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100",
            "legendFormat": "Memory %"
          }
        ]
      },
      {
        "id": 6,
        "title": "Disk Usage",
        "type": "stat",
        "targets": [
          {
            "expr": "(node_filesystem_size_bytes{mountpoint=\"/\"} - node_filesystem_avail_bytes{mountpoint=\"/\"}) / node_filesystem_size_bytes{mountpoint=\"/\"} * 100",
            "legendFormat": "Disk Usage %"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "thresholds": {
              "steps": [
                {"value": null, "color": "green"},
                {"value": 80, "color": "yellow"},
                {"value": 90, "color": "red"}
              ]
            }
          }
        }
      },
      {
        "id": 7,
        "title": "Transaction Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(tendermint_consensus_total_txs[5m])",
            "legendFormat": "TX/s"
          }
        ]
      },
      {
        "id": 8,
        "title": "Validator Status",
        "type": "stat",
        "targets": [
          {
            "expr": "tendermint_consensus_validator_power",
            "legendFormat": "Voting Power"
          }
        ]
      }
    ]
  }
}
```

#### Dashboard Import

1. Access Grafana at `http://<server-ip>:3000`
2. Login with default credentials (admin/admin)
3. Change default password
4. Navigate to Dashboards → Import
5. Upload JSON or paste JSON content
6. Select Prometheus data source
7. Click Import

### 5.4 Alert Rules

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
          description: "No new blocks produced for > 60 seconds on {{ $labels.instance }}"

      # Node offline
      - alert: NodeOffline
        expr: up{job="sharetokens-core"} == 0
        for: 30s
        labels:
          severity: critical
        annotations:
          summary: "Node offline"
          description: "ShareTokens core node {{ $labels.instance }} is not responding"

      # Low peer count
      - alert: LowPeerCount
        expr: tendermint_p2p_peers < 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Low peer count"
          description: "Connected peers < 10 on {{ $labels.instance }}"

      # Disk space low
      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) < 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Disk space critically low"
          description: "Available disk space < 10% on {{ $labels.instance }}"

      # Memory usage high
      - alert: MemoryUsageHigh
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes > 0.95
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Memory usage critically high"
          description: "Memory usage > 95% on {{ $labels.instance }}"

      # CPU usage high
      - alert: CPUUsageHigh
        expr: 100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 90
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "CPU usage high"
          description: "CPU usage > 90% for > 10 minutes on {{ $labels.instance }}"

      # Validator missed blocks
      - alert: ValidatorMissedBlocks
        expr: increase(tendermint_consensus_validator_missed_blocks[1h]) > 50
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Validator missing blocks"
          description: "Validator missed > 50 blocks in last hour"

      # Block sync lagging
      - alert: BlockSyncLag
        expr: tendermint_consensus_latest_block_height - tendermint_consensus_caught_up_height > 100
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Block sync lagging"
          description: "Node is > 100 blocks behind"
```

### 5.5 Alertmanager Configuration

```yaml
# /etc/prometheus/alertmanager.yml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alerts@sharetokens.io'
  smtp_auth_username: 'alerts@sharetokens.io'
  smtp_auth_password: '<app-password>'

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'default'
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty'
    - match:
        severity: warning
      receiver: 'default'

receivers:
  - name: 'default'
    email_configs:
      - to: 'ops@sharetokens.io'
        subject: '[STT Alert] {{ .GroupLabels.alertname }}'
        body: |
          {{ range .Alerts }}
          Alert: {{ .Annotations.summary }}
          Description: {{ .Annotations.description }}
          Instance: {{ .Labels.instance }}
          Severity: {{ .Labels.severity }}
          Time: {{ .StartsAt }}
          {{ end }}

  - name: 'pagerduty'
    pagerduty_configs:
      - service_key: '<pagerduty-key>'
        severity: '{{ .GroupLabels.severity }}'
        description: '{{ .GroupLabels.alertname }}'
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

## 7. Backup & Recovery

### 6.1 Backup Strategy

#### Backup Types

| Type       | Frequency | Retention | Size (Est.) | Priority |
| ---------- | --------- | --------- | ----------- | -------- |
| Full       | Daily     | 30 days   | ~50GB       | High     |
| Incremental| Hourly    | 7 days    | ~1GB        | Medium   |
| State Snapshots | Every 1000 blocks | 14 days | ~5GB | High |
| Key Backup | Once (secure storage) | Forever | <1MB | Critical |
| Config Backup | On change | 90 days | ~1MB | Medium |

#### Automated Backup Script

```bash
#!/bin/bash
# /usr/local/bin/sharetokens-backup.sh

set -e

# Configuration
BACKUP_DIR="/backups/sharetokens"
DATA_DIR="$HOME/.sharetokens/data"
CONFIG_DIR="$HOME/.sharetokens/config"
RETENTION_DAYS=30
DATE=$(date +%Y%m%d_%H%M%S)
LOG_FILE="$BACKUP_DIR/backup_$DATE.log"

# Ensure backup directory exists
mkdir -p "$BACKUP_DIR"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Stop node temporarily (optional - can backup live with --snapshot)
log "Stopping node for backup..."
sudo systemctl stop sharetokensd || true

# Create backup
log "Creating backup..."
tar -czf "$BACKUP_DIR/sharetokens-data-$DATE.tar.gz" -C "$DATA_DIR" . 2>&1 | tee -a "$LOG_FILE"
tar -czf "$BACKUP_DIR/sharetokens-config-$DATE.tar.gz" -C "$CONFIG_DIR" . 2>&1 | tee -a "$LOG_FILE"

# Restart node
log "Restarting node..."
sudo systemctl start sharetokensd

# Upload to cloud storage (S3 example)
log "Uploading to S3..."
aws s3 cp "$BACKUP_DIR/sharetokens-data-$DATE.tar.gz" s3://sharetokens-backups/mainnet/ --storage-class STANDARD_IA
aws s3 cp "$BACKUP_DIR/sharetokens-config-$DATE.tar.gz" s3://sharetokens-backups/mainnet/ --storage-class STANDARD_IA

# Clean up old backups locally
log "Cleaning up old backups..."
find "$BACKUP_DIR" -name "sharetokens-*.tar.gz" -mtime +$RETENTION_DAYS -delete
find "$BACKUP_DIR" -name "backup_*.log" -mtime +$RETENTION_DAYS -delete

log "Backup completed successfully"
```

#### Cron Configuration

```bash
# /etc/cron.d/sharetokens-backup
# Daily full backup at 2 AM
0 2 * * * root /usr/local/bin/sharetokens-backup.sh

# Hourly incremental backup (state only)
0 * * * * root /usr/local/bin/sharetokens-backup-incremental.sh
```

### 6.2 State Snapshots

```bash
# Create state snapshot
sharetokensd snapshots create --home ~/.sharetokens

# List snapshots
sharetokensd snapshots list --home ~/.sharetokens

# Restore from snapshot
sharetokensd snapshots restore <snapshot-height> --home ~/.sharetokens

# Prune old snapshots
sharetokensd snapshots delete-older-than 1000 --home ~/.sharetokens
```

### 6.3 Validator Key Backup

```bash
# Export validator key
sharetokensd keys export validator --keyring-backend file > validator-key.backup

# Encrypt backup
gpg --symmetric --cipher-algo AES256 validator-key.backup

# Store securely (multiple locations)
# 1. Hardware wallet
# 2. Offline USB (encrypted)
# 3. Password manager
# 4. Paper backup (mnemonic)

# Test restoration
gpg --decrypt validator-key.backup.gpg > validator-key.restore
cp validator-key.restore /tmp/validator-key.backup
sharetokensd keys import validator-test /tmp/validator-key.backup --keyring-backend file
sharetokensd keys show validator-test --keyring-backend file
rm /tmp/validator-key.backup
```

### 6.4 Recovery Procedures

#### Scenario 1: Data Corruption

```bash
# Stop node
sudo systemctl stop sharetokensd

# Backup corrupted data (for forensic analysis)
mv ~/.sharetokens/data ~/.sharetokens/data-corrupted-$(date +%Y%m%d)

# Reset node state
sharetokensd unsafe-reset-all

# Download recent snapshot from official source
wget https://snapshots.sharetokens.io/mainnet/latest.tar.lz4

# Extract snapshot
lz4 -d latest.tar.lz4 | tar -x -C ~/.sharetokens/data

# Or restore from local backup
tar -xzf /backups/sharetokens/sharetokens-data-<date>.tar.gz -C ~/.sharetokens/data

# Start node
sudo systemctl start sharetokensd

# Verify sync
sharetokensd status
```

#### Scenario 2: Complete Node Failure

```bash
# On new server:

# 1. Install prerequisites (see Section 1.2)

# 2. Download and restore backup
mkdir -p ~/.sharetokens
tar -xzf sharetokens-config-<date>.tar.gz -C ~/.sharetokens/config
tar -xzf sharetokens-data-<date>.tar.gz -C ~/.sharetokens/data

# 3. Restore validator key (if validator)
sharetokensd keys import validator validator-key.backup --keyring-backend file

# 4. Verify configuration
cat ~/.sharetokens/config/config.toml | grep moniker
cat ~/.sharetokens/config/config.toml | grep seeds

# 5. Start node
sudo systemctl start sharetokensd

# 6. Verify
sharetokensd status
sharetokensd query staking validator <validator-address>
```

#### Scenario 3: Chain Rollback (Consensus Issue)

```bash
# Stop node
sudo systemctl stop sharetokensd

# Backup current state
cp -r ~/.sharetokens/data ~/.sharetokens/data-pre-rollback-$(date +%Y%m%d)

# Rollback to specific height
sharetokensd rollback --height <target-height>

# Or manually reset and sync from snapshot
sharetokensd unsafe-reset-all
# Restore from snapshot (see Scenario 1)

# Start node
sudo systemctl start sharetokensd

# Monitor for consensus
watch -n 5 'sharetokensd status | jq .SyncInfo'
```

### 6.5 Disaster Recovery Plan

#### RTO/RPO Targets

| Component | RTO | RPO |
|-----------|-----|-----|
| Validator Node | 1 hour | 24 hours |
| Full Node | 4 hours | 24 hours |
| Sentry Node | 2 hours | 24 hours |

#### Recovery Checklist

- [ ] Identify failure type (hardware/software/network)
- [ ] Notify team via emergency channels
- [ ] Assess backup availability
- [ ] Provision new server (if hardware failure)
- [ ] Restore from backup or snapshot
- [ ] Verify node syncs to network
- [ ] Verify validator is signing (if applicable)
- [ ] Update monitoring
- [ ] Post-incident report

---

## 8. Emergency Response

### 7.1 Incident Severity Levels

| Level | Description | Response Time | Example |
|-------|-------------|---------------|---------|
| P0 (Critical) | Chain halt, security breach | 15 minutes | Consensus failure, double-sign detected |
| P1 (High) | Validator jailed, sync issues | 1 hour | Validator downtime, sync lag |
| P2 (Medium) | Performance degradation | 4 hours | High latency, peer count low |
| P3 (Low) | Non-critical issues | 24 hours | Minor config issues |

### 7.2 Emergency Contacts & Communication

#### Primary Channels

| Channel | Purpose | Access |
|---------|---------|--------|
| Validator Discord | Real-time coordination | #validators-emergency |
| Telegram | Quick alerts | @STT_Validators |
| Email | Formal communications | validators@sharetokens.io |
| GitHub | Security advisories | github.com/sharetokens/security |

#### Emergency Response Team

| Role | Contact | Responsibility |
|------|---------|----------------|
| On-call Engineer | oncall@sharetokens.io | Initial response |
| Security Lead | security@sharetokens.io | Security incidents |
| DevOps Lead | devops@sharetokens.io | Infrastructure issues |
| Community Lead | community@sharetokens.io | Public communications |

### 7.3 Critical Incident Response

#### P0: Chain Halt Response

```bash
# 1. IMMEDIATE (0-15 min)
# Check node status
sharetokensd status 2>&1 | jq .

# Check logs for errors
journalctl -u sharetokensd -n 1000 --no-pager | grep -i error

# Check if halt is network-wide
# - Check official Discord/Telegram
# - Check block explorers
# - Query other validators

# 2. ASSESS (15-30 min)
# Determine if local or network issue
# If local: attempt restart with backup
# If network: await official guidance

# 3. RECOVERY
# Stop node
sudo systemctl stop sharetokensd

# Restore from last known good state
tar -xzf /backups/sharetokens/sharetokens-data-<last-known-good>.tar.gz -C ~/.sharetokens/data

# Or reset and sync from snapshot
sharetokensd unsafe-reset-all
# Download and restore snapshot

# Restart
sudo systemctl start sharetokensd

# Monitor
watch -n 5 'sharetokensd status | jq .SyncInfo.catching_up'
```

#### P0: Double-Sign Detected

```bash
# If validator accidentally double-signed:

# 1. IMMEDIATELY stop validator
sudo systemctl stop sharetokensd

# 2. Check validator status
sharetokensd query staking validator <validator-address>

# 3. Validator will be tombstoned (cannot unjail)
# Contact team immediately
# Rotate to backup validator if available

# 4. Post-incident
# Review signing process
# Implement multi-sig or HSM
# Document lessons learned
```

#### P1: Validator Jailed

```bash
# Check jail status
sharetokensd query staking validator <validator-address>

# Check missed blocks
sharetokensd query slashing signing-info <validator-pubkey>

# Unjail (if not tombstoned)
sharetokensd tx slashing unjail \
  --from validator \
  --chain-id sharetokens-1 \
  --keyring-backend file \
  --fees 5000stt

# Monitor for recovery
watch -n 30 'sharetokensd query staking validator <validator-address>'
```

### 7.4 Security Incident Response

#### Suspicious Activity Detection

```bash
# Check for unusual transactions
sharetokensd query tx --events message.sender=<suspicious-address>

# Check validator power changes
sharetokensd query staking validators --limit 1000 | jq '.validators[] | select(.tokens != <expected>)'

# Check for unauthorized access attempts
sudo grep "Failed password" /var/log/auth.log
sudo grep "Invalid user" /var/log/auth.log

# Check firewall logs
sudo ufw status verbose
sudo tail -f /var/log/ufw.log
```

#### Security Incident Response Steps

1. **Contain**
   - Isolate affected systems
   - Revoke compromised credentials
   - Block suspicious IPs

2. **Investigate**
   - Preserve logs
   - Identify attack vector
   - Assess impact scope

3. **Recover**
   - Restore from clean backup
   - Apply security patches
   - Update credentials

4. **Report**
   - Document incident
   - Notify affected parties
   - Publish post-mortem

### 7.5 Emergency Runbooks

#### Runbook: Node Won't Start

```bash
# Check 1: Logs
journalctl -u sharetokensd -n 500 --no-pager

# Check 2: Disk space
df -h

# Check 3: Configuration
sharetokensd validate-genesis

# Check 4: Data integrity
# If corrupted:
sudo systemctl stop sharetokensd
sharetokensd unsafe-reset-all
# Restore from backup/snapshot

# Check 5: Binary version
sharetokensd version
# Ensure correct version for network
```

#### Runbook: Sync Stuck

```bash
# Check peer count
sharetokensd query consensus peers

# Check for errors
journalctl -u sharetokensd -f

# Restart with peer refresh
sudo systemctl restart sharetokensd

# If still stuck:
# 1. Update persistent_peers in config.toml
# 2. Reset and restore from snapshot
```

#### Runbook: High Resource Usage

```bash
# CPU high
htop
# Check for:
# - Too many connections
# - Compaction running
# - Query overload

# Memory high
free -h
# Check for:
# - Memory leak
# - Cache not clearing

# Disk high
du -sh ~/.sharetokens/data/*
# Solutions:
# - Enable pruning
# - Compact database
# - Archive old data
```

### 7.6 Post-Incident Review

After every P0/P1 incident:

1. **Timeline Documentation**
   - Detection time
   - Response start
   - Resolution time
   - Root cause identified

2. **Lessons Learned**
   - What went well
   - What could improve
   - Action items

3. **Process Updates**
   - Update runbooks
   - Improve monitoring
   - Revise playbooks

---

## Part 2: Plugin Operations (插件运维)

> 根据节点角色选择部署的可选模块

---

## 8. User Plugin: GenieBot

GenieBot是服务市场的用户端插件，提供 AI 对话界面。

### 8.1 Installation

```bash
# Using npm
npm install @sharetokens/plugin-geniebot

# Or using Docker
docker pull sharetokens/geniebot:latest
```

### 8.2 Configuration

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

### 8.3 Deployment

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

### 8.4 Monitoring

```yaml
# Prometheus scrape config for Xiaodeng
scrape_configs:
  - job_name: 'geniebot'
    static_configs:
      - targets: ['localhost:3000/metrics']
```

---

## 9. Provider Plugin: LLM API Key Hosting

托管 API Key，提供 Level 1 LLM 服务。

### 9.1 Installation

```bash
# Using npm
npm install @sharetokens/plugin-llm-provider

# Or using Docker
docker pull sharetokens/llm-provider:latest
```

### 9.2 Configuration

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

### 9.3 API Key Management

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

### 9.4 Deployment

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

### 9.5 Security Considerations

- API keys are encrypted at rest
- Never log API keys
- Use hardware security modules (HSM) for production
- Rotate keys regularly
- Monitor for unauthorized access

---

## 10. Provider Plugin: Agent Executor (OpenFang)

运行 AI Agent，提供 Level 2 Agent 服务。

### 10.1 Prerequisites

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

### 10.2 Configuration

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

### 10.3 Deployment

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

### 10.4 Monitoring

```bash
# Check agent status
openfang status

# View active tasks
openfang tasks list

# View agent logs
openfang logs --agent coder --follow
```

---

## 11. Provider Plugin: Workflow Executor

编排多 Agent 任务流程，提供 Level 3 Workflow 服务。

### 11.1 Installation

```bash
# Using npm
npm install @sharetokens/plugin-workflow-provider

# Or using Docker
docker pull sharetokens/workflow-provider:latest
```

### 11.2 Configuration

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

### 11.3 Deployment

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

## 12. Complete Plugin Stack Deployment

### 12.1 Full Stack Docker Compose

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

### 12.2 Role-Based Deployment

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

## 13. Troubleshooting

### 13.1 Core Node Issues

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

### 13.2 Plugin Issues

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

## 14. Support & Resources

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

*Document Version: 2.2.0*
*Last Updated: 2026-03-14*
*Maintained by: ShareTokens Operations Team*
