package telemetry

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Namespace = "sharetoken"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	// Chain metrics
	BlockHeight    prometheus.Gauge
	BlockTime      prometheus.Histogram
	TPS            prometheus.Gauge
	ConnectedPeers prometheus.Gauge

	// Transaction metrics
	TransactionCounter  *prometheus.CounterVec
	TransactionDuration *prometheus.HistogramVec
	TransactionGasUsed  *prometheus.HistogramVec

	// Module metrics
	ModuleTxCounter     *prometheus.CounterVec
	ModuleQueryCounter  *prometheus.CounterVec
	ModuleQueryDuration *prometheus.HistogramVec

	// ABCI metrics
	BeginBlockDuration prometheus.Histogram
	EndBlockDuration   prometheus.Histogram
	DeliverTxDuration  prometheus.Histogram

	// Validator metrics
	ValidatorSetSize prometheus.Gauge
	BondedTokens     prometheus.Gauge

	// System metrics (if exposed)
	AppMemoryUsage prometheus.Gauge
	AppGoroutines  prometheus.Gauge
}

var (
	// globalMetrics holds the singleton metrics instance
	globalMetrics *Metrics
)

// NewMetrics creates a new Metrics instance with all collectors registered
func NewMetrics() *Metrics {
	m := &Metrics{
		// Chain metrics
		BlockHeight: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "chain",
			Name:      "block_height",
			Help:      "Current block height",
		}),
		BlockTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "chain",
			Name:      "block_time_seconds",
			Help:      "Time between blocks in seconds",
			Buckets:   []float64{0.1, 0.5, 1, 2, 3, 5, 10, 15, 30},
		}),
		TPS: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "chain",
			Name:      "tps",
			Help:      "Transactions per second",
		}),
		ConnectedPeers: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "chain",
			Name:      "connected_peers",
			Help:      "Number of connected peers",
		}),

		// Transaction metrics
		TransactionCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "tx",
			Name:      "total",
			Help:      "Total number of transactions",
		}, []string{"status", "module"}),
		TransactionDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "tx",
			Name:      "duration_seconds",
			Help:      "Transaction processing duration in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		}, []string{"module"}),
		TransactionGasUsed: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "tx",
			Name:      "gas_used",
			Help:      "Gas used by transactions",
			Buckets:   []float64{1000, 10000, 50000, 100000, 500000, 1000000, 2000000},
		}, []string{"module"}),

		// Module metrics
		ModuleTxCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "module",
			Name:      "tx_total",
			Help:      "Total number of transactions per module",
		}, []string{"module", "msg_type"}),
		ModuleQueryCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "module",
			Name:      "query_total",
			Help:      "Total number of queries per module",
		}, []string{"module", "query_type"}),
		ModuleQueryDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "module",
			Name:      "query_duration_seconds",
			Help:      "Query processing duration per module",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
		}, []string{"module", "query_type"}),

		// ABCI metrics
		BeginBlockDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "abci",
			Name:      "begin_block_duration_seconds",
			Help:      "Duration of BeginBlock in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		}),
		EndBlockDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "abci",
			Name:      "end_block_duration_seconds",
			Help:      "Duration of EndBlock in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		}),
		DeliverTxDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "abci",
			Name:      "deliver_tx_duration_seconds",
			Help:      "Duration of DeliverTx in seconds",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		}),

		// Validator metrics
		ValidatorSetSize: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "validator",
			Name:      "set_size",
			Help:      "Current validator set size",
		}),
		BondedTokens: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "validator",
			Name:      "bonded_tokens",
			Help:      "Total bonded tokens",
		}),

		// System metrics
		AppMemoryUsage: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "app",
			Name:      "memory_usage_bytes",
			Help:      "Application memory usage in bytes",
		}),
		AppGoroutines: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: "app",
			Name:      "goroutines",
			Help:      "Number of goroutines",
		}),
	}

	return m
}

// InitMetrics initializes the global metrics singleton
func InitMetrics() {
	if globalMetrics == nil {
		globalMetrics = NewMetrics()
	}
}

// GetMetrics returns the global metrics instance
func GetMetrics() *Metrics {
	if globalMetrics == nil {
		InitMetrics()
	}
	return globalMetrics
}

// SetBlockHeight sets the current block height
func (m *Metrics) SetBlockHeight(height int64) {
	m.BlockHeight.Set(float64(height))
}

// RecordBlockTime records the time between blocks
func (m *Metrics) RecordBlockTime(duration time.Duration) {
	m.BlockTime.Observe(duration.Seconds())
}

// SetTPS sets the current transactions per second
func (m *Metrics) SetTPS(tps float64) {
	m.TPS.Set(tps)
}

// SetConnectedPeers sets the number of connected peers
func (m *Metrics) SetConnectedPeers(peers int) {
	m.ConnectedPeers.Set(float64(peers))
}

// IncTransaction increments the transaction counter
func (m *Metrics) IncTransaction(status, module string) {
	m.TransactionCounter.WithLabelValues(status, module).Inc()
}

// ObserveTransactionDuration records transaction duration
func (m *Metrics) ObserveTransactionDuration(module string, duration time.Duration) {
	m.TransactionDuration.WithLabelValues(module).Observe(duration.Seconds())
}

// ObserveTransactionGasUsed records gas used by a transaction
func (m *Metrics) ObserveTransactionGasUsed(module string, gasUsed uint64) {
	m.TransactionGasUsed.WithLabelValues(module).Observe(float64(gasUsed))
}

// IncModuleTx increments the module transaction counter
func (m *Metrics) IncModuleTx(module, msgType string) {
	m.ModuleTxCounter.WithLabelValues(module, msgType).Inc()
}

// IncModuleQuery increments the module query counter
func (m *Metrics) IncModuleQuery(module, queryType string) {
	m.ModuleQueryCounter.WithLabelValues(module, queryType).Inc()
}

// ObserveModuleQueryDuration records query duration per module
func (m *Metrics) ObserveModuleQueryDuration(module, queryType string, duration time.Duration) {
	m.ModuleQueryDuration.WithLabelValues(module, queryType).Observe(duration.Seconds())
}

// ObserveBeginBlockDuration records BeginBlock duration
func (m *Metrics) ObserveBeginBlockDuration(duration time.Duration) {
	m.BeginBlockDuration.Observe(duration.Seconds())
}

// ObserveEndBlockDuration records EndBlock duration
func (m *Metrics) ObserveEndBlockDuration(duration time.Duration) {
	m.EndBlockDuration.Observe(duration.Seconds())
}

// ObserveDeliverTxDuration records DeliverTx duration
func (m *Metrics) ObserveDeliverTxDuration(duration time.Duration) {
	m.DeliverTxDuration.Observe(duration.Seconds())
}

// SetValidatorSetSize sets the validator set size
func (m *Metrics) SetValidatorSetSize(size int) {
	m.ValidatorSetSize.Set(float64(size))
}

// SetBondedTokens sets the total bonded tokens
func (m *Metrics) SetBondedTokens(tokens int64) {
	m.BondedTokens.Set(float64(tokens))
}

// SetAppMemoryUsage sets the application memory usage
func (m *Metrics) SetAppMemoryUsage(bytes uint64) {
	m.AppMemoryUsage.Set(float64(bytes))
}

// SetAppGoroutines sets the number of goroutines
func (m *Metrics) SetAppGoroutines(count int) {
	m.AppGoroutines.Set(float64(count))
}

// RecordTxDuration is a helper to record transaction duration with defer
func RecordTxDuration(module string, start time.Time) {
	GetMetrics().ObserveTransactionDuration(module, time.Since(start))
}

// RecordQueryDuration is a helper to record query duration with defer
func RecordQueryDuration(module, queryType string, start time.Time) {
	GetMetrics().ObserveModuleQueryDuration(module, queryType, time.Since(start))
}

// ContextKey is the type for context keys
type ContextKey string

const (
	// MetricsContextKey is used to store metrics in context
	MetricsContextKey ContextKey = "telemetry_metrics"
)

// ContextWithMetrics adds metrics to context
func ContextWithMetrics(ctx context.Context, m *Metrics) context.Context {
	return context.WithValue(ctx, MetricsContextKey, m)
}

// MetricsFromContext extracts metrics from context
func MetricsFromContext(ctx context.Context) *Metrics {
	if m, ok := ctx.Value(MetricsContextKey).(*Metrics); ok {
		return m
	}
	return GetMetrics()
}

// SdkContextWithMetrics adds metrics to SDK context
func SdkContextWithMetrics(ctx sdk.Context, m *Metrics) sdk.Context {
	return ctx.WithContext(ContextWithMetrics(ctx.Context(), m))
}

// SdkMetricsFromContext extracts metrics from SDK context
func SdkMetricsFromContext(ctx sdk.Context) *Metrics {
	return MetricsFromContext(ctx.Context())
}
