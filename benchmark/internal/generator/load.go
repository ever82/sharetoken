package generator

import (
	"context"
	"sync"
	"time"

	"sharetoken/benchmark/internal/metrics"
)

// LoadGenerator generates load for benchmarking
type LoadGenerator struct {
	collector  *metrics.Collector
	workers    int
	rate       int           // Requests per second
	duration   time.Duration
	operation  Operation
}

// Operation represents a benchmark operation
type Operation interface {
	Name() string
	Execute(ctx context.Context) error
}

// OperationFunc allows using functions as operations
type OperationFunc func(ctx context.Context) error

func (f OperationFunc) Name() string { return "custom" }
func (f OperationFunc) Execute(ctx context.Context) error { return f(ctx) }

// Config represents load generator configuration
type Config struct {
	Workers  int
	Rate     int           // Requests per second (0 = unlimited)
	Duration time.Duration
}

// NewLoadGenerator creates a new load generator
func NewLoadGenerator(collector *metrics.Collector, config Config, op Operation) *LoadGenerator {
	return &LoadGenerator{
		collector: collector,
		workers:   config.Workers,
		rate:      config.Rate,
		duration:  config.Duration,
		operation: op,
	}
}

// Run runs the load generator
func (g *LoadGenerator) Run(ctx context.Context) error {
	g.collector.Start()
	defer g.collector.Stop()

	ctx, cancel := context.WithTimeout(ctx, g.duration)
	defer cancel()

	var wg sync.WaitGroup
	workCh := make(chan struct{}, g.workers)

	// Start rate limiter if rate is set
	var rateLimiter <-chan time.Time
	if g.rate > 0 {
		interval := time.Second / time.Duration(g.rate)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		rateLimiter = ticker.C
	}

	// Start workers
	for i := 0; i < g.workers; i++ {
		wg.Add(1)
		go g.worker(ctx, &wg, workCh)
	}

	// Generate work
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(workCh)
				return
			case <-rateLimiter:
				select {
				case workCh <- struct{}{}:
				case <-ctx.Done():
					close(workCh)
					return
				}
			case workCh <- struct{}{}:
				// Unlimited rate
			}
		}
	}()

	wg.Wait()
	return nil
}

// worker processes work items
func (g *LoadGenerator) worker(ctx context.Context, wg *sync.WaitGroup, workCh <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-workCh:
			if !ok {
				return
			}

			start := time.Now()
			err := g.operation.Execute(ctx)
			duration := time.Since(start)

			g.collector.Record(duration, err == nil, err)
		}
	}
}

// RampUpConfig represents ramp-up configuration
type RampUpConfig struct {
	StartWorkers int
	EndWorkers   int
	StepDuration time.Duration
	StepSize     int
}

// RampUpGenerator performs ramp-up testing
type RampUpGenerator struct {
	collector *metrics.Collector
	config    RampUpConfig
	duration  time.Duration
	operation Operation
}

// NewRampUpGenerator creates a new ramp-up generator
func NewRampUpGenerator(collector *metrics.Collector, config RampUpConfig, duration time.Duration, op Operation) *RampUpGenerator {
	return &RampUpGenerator{
		collector: collector,
		config:    config,
		duration:  duration,
		operation: op,
	}
}

// Run runs the ramp-up test
func (g *RampUpGenerator) Run(ctx context.Context) ([]metrics.Stats, error) {
	var results []metrics.Stats

	for workers := g.config.StartWorkers; workers <= g.config.EndWorkers; workers += g.config.StepSize {
		// Create collector for this step
		stepCollector := metrics.NewCollector()

		config := Config{
			Workers:  workers,
			Duration: g.config.StepDuration,
		}

		generator := NewLoadGenerator(stepCollector, config, g.operation)

		if err := generator.Run(ctx); err != nil {
			return results, err
		}

		stats := stepCollector.GetStats()
		results = append(results, stats)

		// Check if performance degraded significantly
		if len(results) > 1 {
			prev := results[len(results)-2]
			if stats.TPS < prev.TPS*0.8 { // TPS dropped by 20%
				break // Stop ramping
			}
		}
	}

	return results, nil
}
