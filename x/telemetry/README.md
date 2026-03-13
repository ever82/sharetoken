# ShareToken Telemetry Package

A comprehensive observability solution for ShareToken blockchain, providing Prometheus metrics, OpenTelemetry tracing, and structured logging.

## Features

- **Prometheus Metrics**: Chain, module, and transaction-level metrics
- **OpenTelemetry Tracing**: Distributed tracing with Jaeger/Zipkin support
- **Structured Logging**: JSON/text logging with configurable levels
- **Module Wrappers**: Automatic instrumentation of Cosmos SDK modules
- **ABCI Hooks**: Telemetry for BeginBlock, EndBlock, and DeliverTx

## Installation

The telemetry package is located at `x/telemetry` and is automatically included in the project.

## Quick Start

### 1. Initialize Telemetry

```go
import "sharetoken/x/telemetry"

func main() {
    // Configure telemetry
    config := telemetry.AppTelemetryConfig{
        MetricsEnabled:       true,
        MetricsServerAddress: ":26660",
        TracingEnabled:       true,
        TracingEndpoint:      "http://localhost:14268/api/traces",
        LoggingFormat:        "json",
        LoggingLevel:         telemetry.InfoLevel,
    }

    // Initialize
    telem, err := telemetry.InitializeAppTelemetry(config)
    if err != nil {
        log.Fatal(err)
    }
    defer telem.Shutdown()
}
```

### 2. Use in Keepers

```go
type MyKeeper struct {
    telemetry *telemetry.TelemetryKeeper
}

func NewMyKeeper() *MyKeeper {
    return &MyKeeper{
        telemetry: telemetry.NewTelemetryKeeper("mymodule"),
    }
}

func (k *MyKeeper) DoSomething(ctx sdk.Context) error {
    // Start timing
    start := time.Now()

    // Your logic here
    err := doWork()

    // Record metrics
    telemetry.RecordKeeperQuery(ctx, "mymodule", "DoSomething", start)

    return err
}
```

### 3. Record Transactions

```go
func (k *MyKeeper) ProcessTx(ctx sdk.Context, msg *types.MsgDoSomething, gasUsed uint64) error {
    // Your logic
    err := processMessage(msg)
    success := err == nil

    // Record transaction telemetry
    telemetry.RecordKeeperTx(ctx, "mymodule", "MsgDoSomething", gasUsed, success)

    return err
}
```

## Available Metrics

### Chain Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `sharetoken_block_height` | Gauge | Current block height |
| `sharetoken_block_time_seconds` | Histogram | Time between blocks |
| `sharetoken_tps` | Gauge | Transactions per second |
| `sharetoken_connected_peers` | Gauge | Number of connected peers |

### Transaction Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `sharetoken_tx_total` | Counter | Total transactions by status and module |
| `sharetoken_tx_duration_seconds` | Histogram | Transaction processing duration |
| `sharetoken_tx_gas_used` | Histogram | Gas used by transactions |

### Module Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `sharetoken_module_tx_total` | Counter | Module transactions by type |
| `sharetoken_module_query_total` | Counter | Module queries by type |
| `sharetoken_module_query_duration_seconds` | Histogram | Query processing duration |

### ABCI Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `sharetoken_abci_begin_block_duration_seconds` | Histogram | BeginBlock duration |
| `sharetoken_abci_end_block_duration_seconds` | Histogram | EndBlock duration |
| `sharetoken_abci_deliver_tx_duration_seconds` | Histogram | DeliverTx duration |

### System Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `sharetoken_app_memory_usage_bytes` | Gauge | Application memory usage |
| `sharetoken_app_goroutines` | Gauge | Number of goroutines |

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEMETRY_ENABLED` | Enable telemetry | `true` |
| `METRICS_ENABLED` | Enable Prometheus metrics | `true` |
| `METRICS_SERVER_ADDRESS` | Metrics server address | `:26660` |
| `TRACING_ENABLED` | Enable distributed tracing | `true` |
| `TRACING_ENDPOINT` | Jaeger/Zipkin endpoint | `http://localhost:14268/api/traces` |
| `TRACING_SAMPLE_RATE` | Tracing sample rate (0-1) | `1.0` |
| `LOG_FORMAT` | Log format (json/text) | `json` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |

### Config File (app.toml)

```toml
[telemetry]
enabled = true
metrics-enabled = true
metrics-server-address = ":26660"
tracing-enabled = true
tracing-endpoint = "http://localhost:14268/api/traces"
tracing-sample-rate = 1.0
log-format = "json"
log-level = "info"
```

## Grafana Dashboards

Pre-configured dashboards are available in `observability/grafana/dashboards/`:

1. **System Metrics** (`system.json`): CPU, memory, disk usage
2. **Blockchain Metrics** (`blockchain.json`): Block height, TPS, peers
3. **Module Metrics** (`modules.json`): Module-level transaction and query metrics
4. **Logs** (`logs.json`): Log volume and error analysis

To import:
```bash
cd observability
docker-compose up -d
# Access Grafana at http://localhost:3000 (admin/admin)
```

## Tracing

The telemetry package supports multiple tracing backends:

### Jaeger

```go
config := telemetry.TracerConfig{
    Enabled:        true,
    ExporterType:   "jaeger",
    JaegerEndpoint: "http://localhost:14268/api/traces",
    SampleRate:     1.0,
}
telemetry.InitTracer(config)
```

### Zipkin

```go
config := telemetry.TracerConfig{
    Enabled:      true,
    ExporterType: "zipkin",
    // Configure Zipkin endpoint
}
```

### Stdout (for debugging)

```go
config := telemetry.TracerConfig{
    Enabled:      true,
    ExporterType: "stdout",
}
```

## Logging

### Structured Logging

```go
// Simple log
telemetry.Info("Processing block")

// With fields
telemetry.WithFields(map[string]interface{}{
    "height": 100,
    "hash": "0xabc...",
}).Info("Block committed")

// Module-specific
logger := telemetry.WithModule("bank")
logger.Info("Transfer completed")

// Events
telemetry.LogEvent("transaction_processed", map[string]interface{}{
    "tx_hash": hash,
    "sender":  sender,
    "amount":  amount,
})
```

### Log Levels

- `DebugLevel` - Detailed debugging information
- `InfoLevel` - General operational information
- `WarnLevel` - Warning events
- `ErrorLevel` - Error events
- `FatalLevel` - Fatal errors (exits application)

## Testing

Run telemetry tests:

```bash
go test ./x/telemetry/...
```

## Integration with Existing Code

### Wrap Existing Modules

```go
// In app.go, wrap modules with telemetry
bankModule := bank.NewAppModule(...)
wrappedBank := telemetry.WrapModule("bank", bankModule)
```

### Add to Message Handlers

```go
func (k msgServer) DoSomething(goCtx context.Context, msg *types.MsgDoSomething) (*types.MsgDoSomethingResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Start telemetry
    spanCtx, span := telemetry.StartTxSpan(ctx.Context(), "mymodule", "MsgDoSomething")
    defer span.End()
    ctx = ctx.WithContext(spanCtx)

    // Your handler logic
    // ...

    return &types.MsgDoSomethingResponse{}, nil
}
```

### Add to Query Handlers

```go
func (k queryServer) GetSomething(goCtx context.Context, req *types.QueryGetSomethingRequest) (*types.QueryGetSomethingResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    start := time.Now()

    // Your query logic
    result, err := k.GetSomething(ctx, req.Id)

    // Record telemetry
    telemetry.RecordKeeperQuery(ctx, "mymodule", "GetSomething", start)

    return result, err
}
```

## Performance Considerations

- Metrics collection has minimal overhead (< 1ms per operation)
- Tracing is sampled based on configuration (default: 100%)
- Logging is asynchronous and non-blocking
- Disable telemetry in production if not needed:
  ```go
  config := telemetry.AppTelemetryConfig{
      MetricsEnabled: false,
      TracingEnabled: false,
  }
  ```

## Troubleshooting

### Metrics not appearing

1. Check metrics server is running: `curl http://localhost:26660/metrics`
2. Verify Prometheus is configured to scrape the endpoint
3. Check application logs for errors

### Tracing not working

1. Verify Jaeger/Zipkin is running
2. Check tracing configuration
3. Look for errors in application logs

### High memory usage

1. Reduce log verbosity: `LOG_LEVEL=warn`
2. Lower tracing sample rate: `TRACING_SAMPLE_RATE=0.1`
3. Disable unnecessary metrics

## Contributing

When adding new telemetry features:

1. Follow naming convention: `sharetoken_<subsystem>_<metric_name>`
2. Add appropriate labels for dimensionality
3. Update this documentation
4. Add corresponding Grafana panels
