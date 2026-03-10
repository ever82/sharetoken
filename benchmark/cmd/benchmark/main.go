package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"sharetoken/benchmark/internal/generator"
	"sharetoken/benchmark/internal/metrics"
	"sharetoken/benchmark/internal/reporter"
)

// Benchmark scenarios
type Scenario struct {
	Name        string
	Description string
	Workers     int
	Duration    time.Duration
	Operation   generator.Operation
}

func main() {
	var (
		scenario    = flag.String("scenario", "transfer", "Benchmark scenario (transfer, query, mixed, stress)")
		workers     = flag.Int("workers", 100, "Number of concurrent workers")
		duration    = flag.Duration("duration", 30*time.Second, "Benchmark duration")
		output      = flag.String("output", "text", "Output format (text, json, markdown)")
		thresholds  = flag.Bool("thresholds", true, "Validate against performance thresholds")
		rampUp      = flag.Bool("rampup", false, "Run ramp-up test")
	)
	flag.Parse()

	fmt.Printf("ShareToken Benchmark Tool\n")
	fmt.Printf("========================\n\n")

	// Create collector
	collector := metrics.NewCollector()

	// Create operation based on scenario
	op := createOperation(*scenario)

	// Run benchmark
	if *rampUp {
		runRampUpTest(collector, op)
	} else {
		runBenchmark(collector, *workers, *duration, op)
	}

	// Generate report
	stats := collector.GetStats()
	rep := reporter.NewReporter(os.Stdout)

	th := metrics.DefaultThresholds()
	if !*thresholds {
		th = metrics.Thresholds{} // Empty thresholds
	}

	switch *output {
	case "json":
		if err := rep.ReportJSON(*scenario, stats); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JSON report: %v\n", err)
		}
	case "markdown":
		rep.ReportMarkdown(*scenario, stats, th)
	default:
		rep.Report(*scenario, stats, th)
	}

	// Exit with error code if thresholds not met
	if *thresholds {
		failures := stats.Validate(th)
		if len(failures) > 0 {
			os.Exit(1)
		}
	}
}

// createOperation creates the benchmark operation
func createOperation(scenario string) generator.Operation {
	switch scenario {
	case "transfer":
		return &TransferOperation{}
	case "query":
		return &QueryOperation{}
	case "mixed":
		return &MixedOperation{}
	case "stress":
		return &StressOperation{}
	default:
		return &TransferOperation{}
	}
}

// runBenchmark runs a single benchmark
func runBenchmark(collector *metrics.Collector, workers int, duration time.Duration, op generator.Operation) {
	fmt.Printf("Running benchmark:\n")
	fmt.Printf("  Workers:  %d\n", workers)
	fmt.Printf("  Duration: %s\n", duration)
	fmt.Printf("  Operation: %s\n\n", op.Name())

	config := generator.Config{
		Workers:  workers,
		Duration: duration,
	}

	gen := generator.NewLoadGenerator(collector, config, op)
	if err := gen.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Benchmark failed: %v\n", err)
		os.Exit(1)
	}
}

// runRampUpTest runs a ramp-up test
func runRampUpTest(collector *metrics.Collector, op generator.Operation) {
	fmt.Printf("Running ramp-up test:\n\n")

	config := generator.RampUpConfig{
		StartWorkers: 10,
		EndWorkers:   1000,
		StepDuration: 30 * time.Second,
		StepSize:     50,
	}

	rampGen := generator.NewRampUpGenerator(collector, config, 5*time.Minute, op)
	results, err := rampGen.Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ramp-up test failed: %v\n", err)
		os.Exit(1)
	}

	// Print ramp-up results
	fmt.Printf("\nRamp-up Results:\n")
	fmt.Printf("%8s %12s %12s %12s %12s\n", "Workers", "TPS", "P50", "P99", "Success%")
	fmt.Printf("%s\n", "---------------------------------------------------------")
	for _, stats := range results {
		fmt.Printf("%8d %12.2f %12s %12s %11.2f%%\n",
			stats.TotalRequests, // This would need to be tracked per step
			stats.TPS,
			formatLatency(stats.P50Latency),
			formatLatency(stats.P99Latency),
			stats.SuccessRate,
		)
	}
}

// TransferOperation simulates token transfers
type TransferOperation struct {
	counter int
}

func (o *TransferOperation) Name() string { return "Token Transfer" }

func (o *TransferOperation) Execute(ctx context.Context) error {
	// Simulate transfer latency
	time.Sleep(time.Millisecond * time.Duration(5+o.counter%10))
	o.counter++

	// Simulate occasional failures (1%)
	if o.counter%100 == 0 {
		return fmt.Errorf("simulated transfer failure")
	}

	return nil
}

// QueryOperation simulates balance queries
type QueryOperation struct {
	counter int
}

func (o *QueryOperation) Name() string { return "Balance Query" }

func (o *QueryOperation) Execute(ctx context.Context) error {
	// Simulate query latency (faster than transfer)
	time.Sleep(time.Millisecond * time.Duration(1+o.counter%5))
	o.counter++
	return nil
}

// MixedOperation simulates mixed workload
type MixedOperation struct {
	counter int
}

func (o *MixedOperation) Name() string { return "Mixed Workload" }

func (o *MixedOperation) Execute(ctx context.Context) error {
	// 70% queries, 30% transfers
	if o.counter%10 < 7 {
		// Query
		time.Sleep(time.Millisecond * time.Duration(1+o.counter%5))
	} else {
		// Transfer
		time.Sleep(time.Millisecond * time.Duration(5+o.counter%10))
	}
	o.counter++
	return nil
}

// StressOperation simulates high load
type StressOperation struct {
	counter int
}

func (o *StressOperation) Name() string { return "Stress Test" }

func (o *StressOperation) Execute(ctx context.Context) error {
	// Simulate variable latency with occasional spikes
	baseLatency := 5 + o.counter%10
	if o.counter%50 == 0 {
		baseLatency *= 10 // Spike
	}
	time.Sleep(time.Millisecond * time.Duration(baseLatency))
	o.counter++

	// Simulate higher failure rate under stress (5%)
	if o.counter%20 == 0 {
		return fmt.Errorf("simulated stress failure")
	}

	return nil
}

// formatLatency formats latency for display
func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}
