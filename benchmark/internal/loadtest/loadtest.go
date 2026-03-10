package loadtest

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"sharetoken/benchmark/internal/metrics"
)

// Config holds load test configuration
type Config struct {
	Endpoint    string
	Duration    time.Duration
	Concurrency int
	TxRate      int
	TestType    string
}

// LoadTester runs load tests
type LoadTester struct {
	config    Config
	metrics   *metrics.Collector
	client    *Client
	stopCh    chan struct{}
}

// Client simulates blockchain client
type Client struct {
	endpoint string
	accounts []string
}

// NewLoadTester creates a new load tester
func NewLoadTester(config Config, metrics *metrics.Collector) *LoadTester {
	return &LoadTester{
		config:  config,
		metrics: metrics,
		client: &Client{
			endpoint: config.Endpoint,
			accounts: generateTestAccounts(config.Concurrency),
		},
		stopCh: make(chan struct{}),
	}
}

// Run executes the load test
func (lt *LoadTester) Run(ctx context.Context) error {
	fmt.Printf("Starting load test with %d concurrent users...\n", lt.config.Concurrency)

	// Calculate TPS per worker
	tpsPerWorker := float64(lt.config.TxRate) / float64(lt.config.Concurrency)
	if tpsPerWorker < 1 {
		tpsPerWorker = 1
	}
	interval := time.Duration(float64(time.Second) / tpsPerWorker)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < lt.config.Concurrency; i++ {
		wg.Add(1)
		go lt.worker(ctx, &wg, i, interval)
	}

	// Run for specified duration
	time.Sleep(lt.config.Duration)
	
	// Signal workers to stop
	close(lt.stopCh)
	
	// Wait for all workers to finish
	wg.Wait()
	
	lt.metrics.Stop()
	
	return nil
}

// worker runs a single concurrent user
func (lt *LoadTester) worker(ctx context.Context, wg *sync.WaitGroup, id int, interval time.Duration) {
	defer wg.Done()
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	account := lt.client.accounts[id%len(lt.client.accounts)]
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-lt.stopCh:
			return
		case <-ticker.C:
			lt.executeTransaction(account)
		}
	}
}

// executeTransaction simulates executing a transaction
func (lt *LoadTester) executeTransaction(account string) {
	start := time.Now()
	
	// Simulate transaction based on test type
	var success bool
	switch lt.config.TestType {
	case "transfer":
		success = lt.simulateTransfer(account)
	case "marketplace":
		success = lt.simulateMarketplace(account)
	case "mixed":
		if rand.Float32() < 0.5 { //nolint:gosec
			success = lt.simulateTransfer(account)
		} else {
			success = lt.simulateMarketplace(account)
		}
	default:
		success = lt.simulateTransfer(account)
	}
	
	latency := time.Since(start)

	lt.metrics.Record(latency, success, nil)
}

// simulateTransfer simulates a token transfer
func (lt *LoadTester) simulateTransfer(from string) bool {
	// Simulate transfer logic
	// In real implementation, this would call the blockchain
	
	// Simulate occasional failures (2% failure rate)
	if rand.Float32() < 0.02 { //nolint:gosec
		return false
	}
	
	// Simulate processing time
	time.Sleep(time.Duration(rand.Int63n(50)+10) * time.Millisecond) //nolint:gosec
	
	return true
}

// simulateMarketplace simulates a marketplace operation
func (lt *LoadTester) simulateMarketplace(account string) bool {
	// Simulate marketplace operations (service purchase, etc.)
	
	// Simulate occasional failures (3% failure rate)
	if rand.Float32() < 0.03 {
		return false
	}
	
	// Simulate processing time (marketplace ops take longer)
	time.Sleep(time.Duration(rand.Int63n(100)+20) * time.Millisecond)
	
	return true
}

// generateTestAccounts generates test account addresses
func generateTestAccounts(count int) []string {
	accounts := make([]string, count)
	for i := 0; i < count; i++ {
		accounts[i] = fmt.Sprintf("sharetoken1test%04d", i)
	}
	return accounts
}
