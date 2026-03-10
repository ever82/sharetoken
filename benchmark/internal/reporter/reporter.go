package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"sharetoken/benchmark/internal/metrics"
)

// Reporter generates benchmark reports
type Reporter struct {
	output io.Writer
}

// NewReporter creates a new reporter
func NewReporter(output io.Writer) *Reporter {
	return &Reporter{output: output}
}

// Report generates a report from stats
func (r *Reporter) Report(name string, stats metrics.Stats, thresholds metrics.Thresholds) {
	fmt.Fprintf(r.output, "\n%s\n", strings.Repeat("=", 60))
	fmt.Fprintf(r.output, "Benchmark Report: %s\n", name)
	fmt.Fprintf(r.output, "%s\n\n", strings.Repeat("=", 60))

	// Summary
	fmt.Fprintf(r.output, "Summary:\n")
	fmt.Fprintf(r.output, "  Duration:        %s\n", stats.Duration)
	fmt.Fprintf(r.output, "  Total Requests:  %d\n", stats.TotalRequests)
	fmt.Fprintf(r.output, "  Successful:      %d\n", stats.SuccessCount)
	fmt.Fprintf(r.output, "  Failed:          %d\n", stats.FailureCount)
	fmt.Fprintf(r.output, "  Success Rate:    %.2f%%\n", stats.SuccessRate)

	// TPS
	tpsStatus := "✓"
	if stats.TPS < thresholds.MinTPS {
		tpsStatus = "✗"
	}
	fmt.Fprintf(r.output, "  TPS:             %.2f %s (target: %.2f)\n", stats.TPS, tpsStatus, thresholds.MinTPS)

	// Latency
	fmt.Fprintf(r.output, "\nLatency Distribution:\n")
	fmt.Fprintf(r.output, "  Min:             %s\n", formatDuration(stats.MinLatency))
	fmt.Fprintf(r.output, "  Avg:             %s\n", formatDuration(stats.AvgLatency))
	fmt.Fprintf(r.output, "  Max:             %s\n", formatDuration(stats.MaxLatency))
	fmt.Fprintf(r.output, "  P50:             %s\n", formatDuration(stats.P50Latency))
	fmt.Fprintf(r.output, "  P90:             %s\n", formatDuration(stats.P90Latency))

	p99Status := "✓"
	if stats.P99Latency > thresholds.MaxP99Latency {
		p99Status = "✗"
	}
	fmt.Fprintf(r.output, "  P99:             %s %s (target: %s)\n",
		formatDuration(stats.P99Latency), p99Status, formatDuration(thresholds.MaxP99Latency))
	fmt.Fprintf(r.output, "  P99.9:           %s\n", formatDuration(stats.P999Latency))

	// Threshold validation
	failures := stats.Validate(thresholds)
	if len(failures) > 0 {
		fmt.Fprintf(r.output, "\n✗ FAILED Thresholds:\n")
		for _, f := range failures {
			fmt.Fprintf(r.output, "  - %s\n", f)
		}
	} else {
		fmt.Fprintf(r.output, "\n✓ All thresholds passed\n")
	}

	fmt.Fprintf(r.output, "\n%s\n", strings.Repeat("=", 60))
}

// ReportJSON generates a JSON report
func (r *Reporter) ReportJSON(name string, stats metrics.Stats) error {
	report := map[string]interface {
	}{
		"benchmark": name,
		"timestamp": time.Now().UTC(),
		"summary": map[string]interface{}{
			"duration":       stats.Duration.String(),
			"total_requests": stats.TotalRequests,
			"success_count":  stats.SuccessCount,
			"failure_count":  stats.FailureCount,
			"success_rate":   stats.SuccessRate,
			"tps":            stats.TPS,
		},
		"latency": map[string]string{
			"min":   stats.MinLatency.String(),
			"avg":   stats.AvgLatency.String(),
			"max":   stats.MaxLatency.String(),
			"p50":   stats.P50Latency.String(),
			"p90":   stats.P90Latency.String(),
			"p99":   stats.P99Latency.String(),
			"p99.9": stats.P999Latency.String(),
		},
	}

	encoder := json.NewEncoder(r.output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// ReportMarkdown generates a Markdown report
func (r *Reporter) ReportMarkdown(name string, stats metrics.Stats, thresholds metrics.Thresholds) {
	fmt.Fprintf(r.output, "# Benchmark Report: %s\n\n", name)
	fmt.Fprintf(r.output, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	// Summary table
	fmt.Fprintf(r.output, "## Summary\n\n")
	fmt.Fprintf(r.output, "| Metric | Value | Threshold | Status |\n")
	fmt.Fprintf(r.output, "|--------|-------|-----------|--------|\n")

	tpsStatus := "✓"
	if stats.TPS < thresholds.MinTPS {
		tpsStatus = "✗"
	}
	fmt.Fprintf(r.output, "| TPS | %.2f | %.2f | %s |\n", stats.TPS, thresholds.MinTPS, tpsStatus)

	p99Status := "✓"
	if stats.P99Latency > thresholds.MaxP99Latency {
		p99Status = "✗"
	}
	fmt.Fprintf(r.output, "| P99 Latency | %s | %s | %s |\n",
		formatDuration(stats.P99Latency), formatDuration(thresholds.MaxP99Latency), p99Status)

	successStatus := "✓"
	if stats.SuccessRate < thresholds.MinSuccessRate {
		successStatus = "✗"
	}
	fmt.Fprintf(r.output, "| Success Rate | %.2f%% | %.2f%% | %s |\n", stats.SuccessRate, thresholds.MinSuccessRate, successStatus)

	// Details
	fmt.Fprintf(r.output, "\n## Details\n\n")
	fmt.Fprintf(r.output, "- **Duration**: %s\n", stats.Duration)
	fmt.Fprintf(r.output, "- **Total Requests**: %d\n", stats.TotalRequests)
	fmt.Fprintf(r.output, "- **Successful**: %d\n", stats.SuccessCount)
	fmt.Fprintf(r.output, "- **Failed**: %d\n", stats.FailureCount)

	// Latency table
	fmt.Fprintf(r.output, "\n## Latency Distribution\n\n")
	fmt.Fprintf(r.output, "| Percentile | Latency |\n")
	fmt.Fprintf(r.output, "|------------|---------|\n")
	fmt.Fprintf(r.output, "| Min | %s |\n", formatDuration(stats.MinLatency))
	fmt.Fprintf(r.output, "| Avg | %s |\n", formatDuration(stats.AvgLatency))
	fmt.Fprintf(r.output, "| Max | %s |\n", formatDuration(stats.MaxLatency))
	fmt.Fprintf(r.output, "| P50 | %s |\n", formatDuration(stats.P50Latency))
	fmt.Fprintf(r.output, "| P90 | %s |\n", formatDuration(stats.P90Latency))
	fmt.Fprintf(r.output, "| P99 | %s |\n", formatDuration(stats.P99Latency))
	fmt.Fprintf(r.output, "| P99.9 | %s |\n", formatDuration(stats.P999Latency))

	// Conclusion
	failures := stats.Validate(thresholds)
	if len(failures) > 0 {
		fmt.Fprintf(r.output, "\n## ❌ Failed Thresholds\n\n")
		for _, f := range failures {
			fmt.Fprintf(r.output, "- %s\n", f)
		}
	} else {
		fmt.Fprintf(r.output, "\n## ✅ All Thresholds Passed\n")
	}
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}
