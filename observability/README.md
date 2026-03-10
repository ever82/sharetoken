# ShareToken Observability Stack

Complete monitoring, logging, and alerting solution for ShareToken blockchain nodes.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      Observability Stack                        │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Prometheus  │  │    Loki      │  │ AlertManager │          │
│  │  (Metrics)   │  │   (Logs)     │  │  (Alerts)    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                   │
│         └─────────────────┼─────────────────┘                   │
│                           │                                     │
│                   ┌───────┴───────┐                            │
│                   │    Grafana    │                            │
│                   │(Visualization)│                            │
│                   └───────────────┘                            │
├─────────────────────────────────────────────────────────────────┤
│  Exporters: Node Exporter | Cadvisor | Promtail | ShareToken   │
└─────────────────────────────────────────────────────────────────┘
```

## Components

### Metrics (Prometheus)
- **Prometheus v2.48.0**: Time-series metrics collection and storage
- **Node Exporter v1.7.0**: System-level metrics (CPU, memory, disk, network)
- **Cadvisor v0.47.2**: Container metrics
- **ShareToken Exporter**: Blockchain-specific metrics (block height, TPS, peers)

### Logs (Loki)
- **Loki 2.9.0**: Log aggregation and storage
- **Promtail 2.9.0**: Log collection and shipping
- Collects system logs, Docker container logs, and application logs

### Alerting (AlertManager)
- **AlertManager v0.26.0**: Alert routing and notification management
- Multi-channel notifications (email, webhook)
- Alert grouping and inhibition rules
- Severity-based routing (critical, warning, info)

### Visualization (Grafana)
- **Grafana 10.2.3**: Dashboards and visualization
- Pre-configured dashboards:
  - **System Metrics**: CPU, memory, disk usage
  - **Blockchain Metrics**: Block height, block time, TPS, peers
  - **Logs**: Log volume, error analysis, search

## Quick Start

### Start the Stack
```bash
cd observability
docker-compose up -d
```

### Access Services
| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | admin/admin |
| Prometheus | http://localhost:9090 | - |
| AlertManager | http://localhost:9093 | - |
| Loki | http://localhost:3100 | - |

### Stop the Stack
```bash
docker-compose down
```

To remove data volumes:
```bash
docker-compose down -v
```

## Dashboards

### 1. System Metrics (`sharetoken-system`)
- CPU Usage gauge with thresholds (green <60%, yellow <80%, red >80%)
- Memory Usage gauge with thresholds (green <70%, yellow <90%, red >90%)
- Disk Usage gauge with thresholds (green <70%, yellow <90%, red >90%)
- Data from Node Exporter

### 2. Blockchain Metrics (`sharetoken-blockchain`)
- Block Height (current value)
- Block Time P99 (99th percentile)
- TPS (Transactions Per Second)
- Connected Peers count
- Auto-refresh every 5 seconds

### 3. Logs Dashboard (`sharetoken-logs`)
- Log volume by level (error, warn, info, debug)
- Real-time log stream with search
- Error rate statistics
- Top error messages table

## Alerts

### Blockchain Alerts
| Alert | Severity | Condition | Duration |
|-------|----------|-----------|----------|
| HighBlockTime | warning | Block time > 10s | 5m |
| LowTPS | warning | TPS < 10 | 10m |
| ConsensusFailure | critical | Any consensus failure | 1m |
| HighMemoryUsage | warning | Memory > 90% | 5m |
| NodeDown | critical | Node unavailable | 1m |

### System Alerts
| Alert | Severity | Condition | Duration |
|-------|----------|-----------|----------|
| HighCPUUsage | warning | CPU > 80% | 5m |
| SystemHighMemoryUsage | warning | Memory > 90% | 5m |
| LowDiskSpace | critical | Disk < 10% | 5m |
| NetworkErrors | warning | Network errors detected | 5m |

### Alert Routing
- **Critical alerts**: Immediate notification (10s wait), repeat every 30m
- **Warning alerts**: Standard notification (1m wait), repeat every 2h
- **Blockchain alerts**: Route to blockchain-ops team
- **System alerts**: Route to system-ops team

## Configuration

### Prometheus
- **Config**: `prometheus/prometheus.yml`
- **Rules**: `prometheus/recording_rules.yml`, `prometheus/alerting_rules.yml`
- **Retention**: 30 days
- **Scrape interval**: 15s (5s for blockchain metrics)

### AlertManager
- **Config**: `alertmanager/alertmanager.yml`
- **Templates**: `alertmanager/templates/*.tmpl`
- Update email settings for production use

### Loki
- **Config**: `loki/loki-config.yml`
- **Retention**: 14 days
- **Max log age**: 12 hours for journal

### Promtail
- **Config**: `promtail/promtail-config.yml`
- **Jobs**: system-journal, docker, sharetoken-node, application

### Grafana
- **Datasources**: `grafana/provisioning/datasources/`
- **Dashboards**: `grafana/provisioning/dashboards/`
- **Dashboard JSON**: `grafana/dashboards/`

## Production Considerations

### Security
1. Change default Grafana password
2. Configure SMTP for email alerts
3. Enable HTTPS/TLS for all endpoints
4. Use authentication for Prometheus and AlertManager
5. Restrict network access to monitoring ports

### Scaling
1. Use external Prometheus storage (Thanos, Cortex)
2. Deploy Loki in distributed mode
3. Use external AlertManager cluster
4. Consider dedicated Grafana instance

### High Availability
1. Run multiple Prometheus instances with federation
2. Deploy AlertManager in cluster mode
3. Use object storage for Loki (S3, GCS)
4. Enable Grafana alerting HA

## Troubleshooting

### Check Service Status
```bash
docker-compose ps
docker-compose logs <service>
```

### Verify Prometheus Targets
Visit http://localhost:9090/targets

### Test Alerts
```bash
curl -X POST http://localhost:9093/-/reload
```

### Check Log Collection
```bash
curl http://localhost:3100/loki/api/v1/label/job/values
```

## Integration with ShareToken Node

To integrate with a running ShareToken node:

1. **Enable metrics endpoint** in node configuration:
```toml
[telemetry]
enabled = true
prometheus-retention-time = 600
```

2. **Update Prometheus targets** in `prometheus/prometheus.yml`:
```yaml
- targets: ['host.docker.internal:26660']
```

3. **Add node logs** to Promtail configuration:
```yaml
- job_name: sharetoken-node
  static_configs:
    - targets:
        - localhost
      labels:
        job: sharetoken-node
        __path__: /path/to/node/logs/*.log
```

## API Endpoints

### Prometheus
- Query: `GET /api/v1/query?query=<promql>`
- Range query: `GET /api/v1/query_range?query=<promql>&start=<ts>&end=<ts>&step=<duration>`
- Targets: `GET /api/v1/targets`
- Alerts: `GET /api/v1/alerts`

### Loki
- Query: `GET /loki/api/v1/query?query=<logql>`
- Range query: `GET /loki/api/v1/query_range?query=<logql>&start=<ts>&end=<ts>`
- Labels: `GET /loki/api/v1/label/<name>/values`

### AlertManager
- Status: `GET /api/v1/status`
- Alerts: `GET /api/v1/alerts`
- Silences: `GET /api/v1/silences`
