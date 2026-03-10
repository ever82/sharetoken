# ACH-DEV-018: Observability Stack

## Summary
Complete observability infrastructure for ShareToken blockchain including metrics collection, log aggregation, alerting, and visualization.

## Components Implemented

### 1. Metrics Collection (Prometheus)
- **Prometheus v2.48.0**: Time-series metrics collection
- **Node Exporter v1.7.0**: System metrics (CPU, memory, disk, network)
- **Cadvisor v0.47.2**: Container metrics
- **ShareToken Exporter**: Mock metrics endpoint for blockchain data
- **Scrape targets**: Prometheus self, Node Exporter, Cadvisor, ShareToken node
- **Recording rules**: TPS, block time, validator participation pre-aggregation
- **Alerting rules**: 9 alerts covering blockchain and system metrics

### 2. Log Aggregation (Loki + Promtail)
- **Loki 2.9.0**: Log aggregation with 14-day retention
- **Promtail 2.9.0**: Log collection agent
- **Log sources**:
  - System journal logs
  - Docker container logs
  - ShareToken node logs (structured JSON)
  - Application logs
- **Log pipeline**: JSON parsing, label extraction, timestamp normalization

### 3. Alerting (AlertManager)
- **AlertManager v0.26.0**: Alert routing and notifications
- **Alert groups**: sharetoken_alerts, system_alerts
- **Severity levels**: critical, warning
- **Routing rules**: Severity-based and category-based (blockchain/system)
- **Notification channels**: Email, webhook
- **Templates**: HTML and text templates for critical, warning, blockchain, system alerts
- **Inhibition rules**: Prevent alert storms

### 4. Visualization (Grafana)
- **Grafana 10.2.3**: Dashboards and visualization
- **Datasources**: Prometheus, Loki
- **Dashboards**:
  - System Metrics (CPU, Memory, Disk gauges)
  - Blockchain Metrics (Block Height, Block Time, TPS, Peers)
  - Logs (Log volume, error analysis, search)
- **Auto-provisioning**: Datasources and dashboards auto-configured

## Files Created

```
observability/
├── docker-compose.yml              # Complete stack orchestration
├── README.md                       # Documentation
├── prometheus/
│   ├── prometheus.yml              # Prometheus configuration
│   ├── recording_rules.yml         # Recording rules
│   ├── alerting_rules.yml          # Alerting rules
│   └── mock-metrics.conf           # Mock metrics for testing
├── alertmanager/
│   ├── alertmanager.yml            # Alert routing configuration
│   └── templates/
│       ├── critical.tmpl           # Critical alert templates
│       ├── warning.tmpl            # Warning alert templates
│       ├── blockchain.tmpl         # Blockchain alert templates
│       └── system.tmpl             # System alert templates
├── loki/
│   └── loki-config.yml             # Loki configuration
├── promtail/
│   └── promtail-config.yml         # Promtail configuration
└── grafana/
    ├── provisioning/
    │   ├── dashboards/
    │   │   └── default.yml         # Dashboard provisioning
    │   └── datasources/
    │       └── prometheus.yml      # Datasource provisioning
    └── dashboards/
        ├── blockchain.json         # Blockchain dashboard
        ├── system.json             # System metrics dashboard
        └── logs.json               # Logs dashboard
```

## Configuration Details

### Prometheus Configuration
- Scrape interval: 15s (5s for blockchain)
- Retention: 30 days
- AlertManager endpoint: alertmanager:9093
- Recording rules: tps, block_time_99, validator_participation

### AlertManager Configuration
- Resolve timeout: 5m
- Group wait: 30s (10s critical, 1m warning)
- Group interval: 5m
- Repeat interval: 4h (30m critical, 2h warning)
- Inhibition: Critical suppresses warnings

### Loki Configuration
- Retention: 14 days
- Max log age: 12h
- Ingestion rate: 16 MB/s
- Ingestion burst: 32 MB

### Grafana Configuration
- Admin: admin/admin (change in production)
- Plugins: clock-panel, simple-json-datasource
- Auto-refresh: dashboards reload every 10s

## Alerts Configured

### Blockchain Alerts
| Alert | Severity | Threshold | For |
|-------|----------|-----------|-----|
| HighBlockTime | warning | > 10s | 5m |
| LowTPS | warning | < 10 | 10m |
| ConsensusFailure | critical | > 0 | 1m |
| HighMemoryUsage | warning | > 90% | 5m |
| NodeDown | critical | = 0 | 1m |

### System Alerts
| Alert | Severity | Threshold | For |
|-------|----------|-----------|-----|
| HighCPUUsage | warning | > 80% | 5m |
| SystemHighMemoryUsage | warning | > 90% | 5m |
| LowDiskSpace | critical | < 10% | 5m |
| NetworkErrors | warning | > 0 | 5m |

## Usage

### Start the Stack
```bash
cd observability
docker-compose up -d
```

### Access Points
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- AlertManager: http://localhost:9093
- Loki: http://localhost:3100

### Verify Installation
```bash
# Check all services are running
docker-compose ps

# View logs
docker-compose logs -f

# Test Prometheus targets
curl http://localhost:9090/api/v1/targets

# Test Loki labels
curl http://localhost:3100/loki/api/v1/label/job/values
```

## Integration

To integrate with actual ShareToken node:

1. Enable metrics in node config:
```toml
[telemetry]
enabled = true
prometheus-retention-time = 600
```

2. Update Prometheus targets:
```yaml
- targets: ['host.docker.internal:26660']
```

3. Configure Promtail for node logs:
```yaml
- job_name: sharetoken-node
  static_configs:
    - targets:
        - localhost
      labels:
        job: sharetoken-node
        __path__: /path/to/logs/*.log
```

## Production Considerations

1. **Security**
   - Change default Grafana password
   - Configure SMTP credentials in AlertManager
   - Enable HTTPS/TLS
   - Use authentication

2. **Scaling**
   - External Prometheus storage (Thanos/Cortex)
   - Distributed Loki deployment
   - AlertManager clustering

3. **High Availability**
   - Multiple Prometheus instances
   - Object storage for Loki
   - Grafana HA setup

## Testing

All components validated:
- [x] Prometheus scrapes all targets
- [x] Recording rules execute correctly
- [x] Alerting rules evaluate properly
- [x] AlertManager routes alerts by severity
- [x] Grafana displays all dashboards
- [x] Loki receives and indexes logs
- [x] Promtail ships logs successfully

## References

- Prometheus: https://prometheus.io/
- Grafana: https://grafana.com/
- Loki: https://grafana.com/oss/loki/
- AlertManager: https://prometheus.io/docs/alerting/
