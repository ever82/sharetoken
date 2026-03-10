# ShareToken Benchmark Tool

Performance benchmarking tool for ShareToken blockchain.

## Features

- Load generation with configurable concurrency
- TPS (Transactions Per Second) measurement
- Latency distribution tracking (P50, P90, P99, P999)
- Multiple benchmark scenarios
- Threshold validation
- Multiple output formats

## Installation

```bash
cd benchmark
go build -o bin/benchmark ./cmd/benchmark
```

## Usage

### Basic Benchmark

```bash
# Run transfer benchmark with 100 workers for 30 seconds
./bin/benchmark -scenario=transfer -workers=100 -duration=30s

# Run query benchmark
./bin/benchmark -scenario=query -workers=200 -duration=60s

# Run mixed workload
./bin/benchmark -scenario=mixed -workers=150 -duration=30s

# Run stress test
./bin/benchmark -scenario=stress -workers=1000 -duration=5m
```

### Output Formats

```bash
# JSON output
./bin/benchmark -output=json > results.json

# Markdown report
./bin/benchmark -output=markdown > report.md

# Text output (default)
./bin/benchmark -output=text
```

### Ramp-up Testing

```bash
./bin/benchmark -rampup=true -scenario=transfer
```

### Without Threshold Validation

```bash
./bin/benchmark -thresholds=false
```

## Benchmark Scenarios

### transfer
Simulates token transfer transactions.
- Target latency: ~5-15ms
- Failure rate: ~1%

### query
Simulates balance query operations.
- Target latency: ~1-5ms
- Failure rate: ~0%

### mixed
Simulates mixed workload (70% queries, 30% transfers).
- Target latency: ~2-10ms
- Failure rate: ~0.3%

### stress
Simulates high load with latency spikes.
- Target latency: ~5-150ms
- Failure rate: ~5%

## Performance Thresholds

Default thresholds (can be customized in code):

- **TPS**: >= 100
- **P99 Latency**: < 3s
- **Success Rate**: >= 95%

## Sample Output

```
============================================================
Benchmark Report: transfer
============================================================

Summary:
  Duration:        30s
  Total Requests:  3000
  Successful:      2970
  Failed:          30
  Success Rate:    99.00%
  TPS:             99.00 ✓ (target: 100.00)

Latency Distribution:
  Min:             5.00ms
  Avg:             9.50ms
  Max:             15.00ms
  P50:             9.00ms
  P90:             12.00ms
  P99:             14.50ms ✓ (target: 3.00s)
  P99.9:           14.95ms

✓ All thresholds passed

============================================================
```

## Running All Benchmarks

```bash
# Run all benchmark scenarios
make benchmark

# Run specific scenario
make benchmark-transfer
make benchmark-query
make benchmark-mixed
make benchmark-stress
```

## Integration Tests

```bash
go test ./benchmark/... -v
```

## Architecture

```
benchmark/
├── cmd/benchmark/
│   └── main.go           # CLI entry point
├── internal/
│   ├── generator/
│   │   └── load.go       # Load generator
│   ├── metrics/
│   │   └── collector.go  # Metrics collection
│   └── reporter/
│       └── reporter.go   # Report generation
└── README.md
```

## Extending

To add a new benchmark scenario:

1. Create a new operation type in `main.go`:
```go
type MyOperation struct{}

func (o *MyOperation) Name() string { return "My Operation" }

func (o *MyOperation) Execute(ctx context.Context) error {
    // Your benchmark logic here
    return nil
}
```

2. Add to `createOperation` function:
```go
case "myscenario":
    return &MyOperation{}
```

3. Run:
```bash
./bin/benchmark -scenario=myscenario
```
