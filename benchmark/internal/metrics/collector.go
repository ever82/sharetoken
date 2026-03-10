package metrics

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// LatencyMetric represents a latency measurement
type LatencyMetric struct {
	Duration  time.Duration
	Timestamp time.Time
	Success   bool
	Error     error
}

// Collector collects and analyzes metrics
type Collector struct {
	latencies     []LatencyMetric
	totalRequests int64
	successCount  int64
	failureCount  int64
	mutex         sync.RWMutex
	startTime     time.Time
	endTime       time.Time
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		latencies: make([]LatencyMetric, 0),
	}
}

// Start marks the start of collection
func (c *Collector) Start() {
	c.mutex.Lock()
	c.startTime = time.Now()
	c.mutex.Unlock()
}

// Stop marks the end of collection
func (c *Collector) Stop() {
	c.mutex.Lock()
	c.endTime = time.Now()
	c.mutex.Unlock()
}

// Record records a latency measurement
func (c *Collector) Record(duration time.Duration, success bool, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.latencies = append(c.latencies, LatencyMetric{
		Duration:  duration,
		Timestamp: time.Now(),
		Success:   success,
		Error:     err,
	})

	c.totalRequests++
	if success {
		c.successCount++
	} else {
		c.failureCount++
	}
}

// GetTPS returns transactions per second
func (c *Collector) GetTPS() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.startTime.IsZero() || c.endTime.IsZero() {
		return 0
	}

	duration := c.endTime.Sub(c.startTime).Seconds()
	if duration == 0 {
		return 0
	}

	return float64(c.successCount) / duration
}

// GetLatencies returns all successful latencies
func (c *Collector) GetLatencies() []time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var latencies []time.Duration
	for _, m := range c.latencies {
		if m.Success {
			latencies = append(latencies, m.Duration)
		}
	}

	return latencies
}

// GetPercentile returns the percentile latency
func (c *Collector) GetPercentile(p float64) time.Duration {
	latencies := c.GetLatencies()
	if len(latencies) == 0 {
		return 0
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	index := int(float64(len(latencies)) * p / 100)
	if index >= len(latencies) {
		index = len(latencies) - 1
	}

	return latencies[index]
}

// GetStats returns comprehensive statistics
func (c *Collector) GetStats() Stats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	latencies := c.GetLatencies()

	return Stats{
		TotalRequests: c.totalRequests,
		SuccessCount:  c.successCount,
		FailureCount:  c.failureCount,
		SuccessRate:   float64(c.successCount) / float64(c.totalRequests) * 100,
		TPS:           c.GetTPS(),
		MinLatency:    c.getMinLatency(latencies),
		MaxLatency:    c.getMaxLatency(latencies),
		AvgLatency:    c.getAvgLatency(latencies),
		P50Latency:    c.GetPercentile(50),
		P90Latency:    c.GetPercentile(90),
		P99Latency:    c.GetPercentile(99),
		P999Latency:   c.GetPercentile(99.9),
		Duration:      c.endTime.Sub(c.startTime),
	}
}

// getMinLatency returns minimum latency
func (c *Collector) getMinLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	min := latencies[0]
	for _, l := range latencies {
		if l < min {
			min = l
		}
	}
	return min
}

// getMaxLatency returns maximum latency
func (c *Collector) getMaxLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	max := latencies[0]
	for _, l := range latencies {
		if l > max {
			max = l
		}
	}
	return max
}

// getAvgLatency returns average latency
func (c *Collector) getAvgLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}

	return sum / time.Duration(len(latencies))
}

// Stats represents benchmark statistics
type Stats struct {
	TotalRequests int64
	SuccessCount  int64
	FailureCount  int64
	SuccessRate   float64
	TPS           float64
	MinLatency    time.Duration
	MaxLatency    time.Duration
	AvgLatency    time.Duration
	P50Latency    time.Duration
	P90Latency    time.Duration
	P99Latency    time.Duration
	P999Latency   time.Duration
	Duration      time.Duration
}

// ValidateThresholds validates against performance thresholds
type Thresholds struct {
	MinTPS         float64
	MaxP99Latency  time.Duration
	MinSuccessRate float64
}

// DefaultThresholds returns default performance thresholds
func DefaultThresholds() Thresholds {
	return Thresholds{
		MinTPS:         100,
		MaxP99Latency:  3 * time.Second,
		MinSuccessRate: 95.0,
	}
}

// Validate checks if stats meet thresholds
func (s Stats) Validate(th Thresholds) []string {
	failures := []string{}

	if s.TPS < th.MinTPS {
		failures = append(failures, fmt.Sprintf("TPS %.2f < %.2f", s.TPS, th.MinTPS))
	}

	if s.P99Latency > th.MaxP99Latency {
		failures = append(failures, fmt.Sprintf("P99 %s > %s", s.P99Latency, th.MaxP99Latency))
	}

	if s.SuccessRate < th.MinSuccessRate {
		failures = append(failures, fmt.Sprintf("Success rate %.2f%% < %.2f%%", s.SuccessRate, th.MinSuccessRate))
	}

	return failures
}
