package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	c := NewCollector()

	c.Start()

	// Record some metrics
	c.Record(10*time.Millisecond, true, nil)
	c.Record(20*time.Millisecond, true, nil)
	c.Record(30*time.Millisecond, true, nil)
	c.Record(40*time.Millisecond, false, nil)

	c.Stop()

	// Check totals
	require.Equal(t, int64(4), c.totalRequests)
	require.Equal(t, int64(3), c.successCount)
	require.Equal(t, int64(1), c.failureCount)

	// Check TPS
	tps := c.GetTPS()
	require.Greater(t, tps, float64(0))
}

func TestPercentiles(t *testing.T) {
	c := NewCollector()

	// Record latencies: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ms
	for i := 1; i <= 10; i++ {
		c.Record(time.Duration(i)*time.Millisecond, true, nil)
	}

	// P50 should be ~5ms
	p50 := c.GetPercentile(50)
	require.GreaterOrEqual(t, p50, 5*time.Millisecond)
	require.Less(t, p50, 7*time.Millisecond)

	// P90 should be ~9ms
	p90 := c.GetPercentile(90)
	require.GreaterOrEqual(t, p90, 9*time.Millisecond)

	// P99 should be ~10ms
	p99 := c.GetPercentile(99)
	require.Equal(t, 10*time.Millisecond, p99)
}

func TestStats(t *testing.T) {
	c := NewCollector()

	c.Start()
	c.Record(10*time.Millisecond, true, nil)
	c.Record(20*time.Millisecond, true, nil)
	c.Record(30*time.Millisecond, true, nil)
	c.Record(100*time.Millisecond, true, nil)
	c.Record(50*time.Millisecond, false, nil)
	c.Stop()

	stats := c.GetStats()

	require.Equal(t, int64(5), stats.TotalRequests)
	require.Equal(t, int64(4), stats.SuccessCount)
	require.Equal(t, int64(1), stats.FailureCount)
	require.Equal(t, 80.0, stats.SuccessRate)

	// Latency stats
	require.Equal(t, 10*time.Millisecond, stats.MinLatency)
	require.Equal(t, 100*time.Millisecond, stats.MaxLatency)
	require.Equal(t, 40*time.Millisecond, stats.AvgLatency)
}

func TestThresholdValidation(t *testing.T) {
	th := Thresholds{
		MinTPS:         100,
		MaxP99Latency:  3 * time.Second,
		MinSuccessRate: 95,
	}

	// Passing stats
	passingStats := Stats{
		TPS:         150,
		P99Latency:  1 * time.Second,
		SuccessRate: 99,
	}
	failures := passingStats.Validate(th)
	require.Empty(t, failures)

	// Failing stats
	failingStats := Stats{
		TPS:         50,
		P99Latency:  5 * time.Second,
		SuccessRate: 90,
	}
	failures = failingStats.Validate(th)
	require.Len(t, failures, 3)
}
