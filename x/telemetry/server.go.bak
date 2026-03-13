package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer is a Prometheus metrics HTTP server
type MetricsServer struct {
	server   *http.Server
	config   MetricsServerConfig
	registry *prometheus.Registry
}

// MetricsServerConfig holds configuration for the metrics server
type MetricsServerConfig struct {
	// Enabled enables the metrics server
	Enabled bool
	// Address is the server address (default: ":26660")
	Address string
	// Path is the metrics endpoint path (default: "/metrics")
	Path string
	// ReadTimeout is the read timeout
	ReadTimeout time.Duration
	// WriteTimeout is the write timeout
	WriteTimeout time.Duration
}

// DefaultMetricsServerConfig returns default metrics server configuration
func DefaultMetricsServerConfig() MetricsServerConfig {
	return MetricsServerConfig{
		Enabled:      true,
		Address:      ":26660",
		Path:         "/metrics",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(config MetricsServerConfig) *MetricsServer {
	if !config.Enabled {
		return nil
	}

	// Create a custom registry
	registry := prometheus.NewRegistry()

	// Register the default Go metrics
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Register our custom metrics
	metrics := GetMetrics()
	if metrics != nil {
		registry.MustRegister(metrics.BlockHeight)
		registry.MustRegister(metrics.BlockTime)
		registry.MustRegister(metrics.TPS)
		registry.MustRegister(metrics.ConnectedPeers)
		registry.MustRegister(metrics.TransactionCounter)
		registry.MustRegister(metrics.TransactionDuration)
		registry.MustRegister(metrics.TransactionGasUsed)
		registry.MustRegister(metrics.ModuleTxCounter)
		registry.MustRegister(metrics.ModuleQueryCounter)
		registry.MustRegister(metrics.ModuleQueryDuration)
		registry.MustRegister(metrics.BeginBlockDuration)
		registry.MustRegister(metrics.EndBlockDuration)
		registry.MustRegister(metrics.DeliverTxDuration)
		registry.MustRegister(metrics.ValidatorSetSize)
		registry.MustRegister(metrics.BondedTokens)
		registry.MustRegister(metrics.AppMemoryUsage)
		registry.MustRegister(metrics.AppGoroutines)
	}

	// Create HTTP server
	mux := http.NewServeMux()
	mux.Handle(config.Path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	}))

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add ready check endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})

	server := &http.Server{
		Addr:         config.Address,
		Handler:      mux,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return &MetricsServer{
		server:   server,
		config:   config,
		registry: registry,
	}
}

// Start starts the metrics server
func (s *MetricsServer) Start() error {
	if s.server == nil {
		return nil
	}

	Infof("Starting metrics server on %s%s", s.config.Address, s.config.Path)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Errorf("Metrics server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the metrics server
func (s *MetricsServer) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	Info("Stopping metrics server")

	return s.server.Shutdown(ctx)
}

// Registry returns the Prometheus registry
func (s *MetricsServer) Registry() *prometheus.Registry {
	return s.registry
}

// RegisterCollector registers a custom collector
func (s *MetricsServer) RegisterCollector(collector prometheus.Collector) error {
	return s.registry.Register(collector)
}

// TelemetryHandler is a middleware for HTTP telemetry
func TelemetryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start)
		metrics := GetMetrics()
		if metrics != nil {
			metrics.ModuleQueryCounter.WithLabelValues(
				"api",
				r.URL.Path,
			).Inc()
			metrics.ModuleQueryDuration.WithLabelValues(
				"api",
				r.URL.Path,
			).Observe(duration.Seconds())
		}

		// Log
		LogEvent("http_request", map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": wrapped.statusCode,
			"duration":    duration.Seconds(),
			"client_ip":   r.RemoteAddr,
		})
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// GRPCUnaryInterceptor is a gRPC interceptor for telemetry
// This is a placeholder - in production, use proper gRPC interceptors
func GRPCUnaryInterceptor(ctx context.Context, req interface{}, info string, handler func(context.Context, interface{}) (interface{}, error)) (interface{}, error) {
	start := time.Now()
	spanCtx, span := StartSpan(ctx, info)
	defer span.End()

	resp, err := handler(spanCtx, req)

	duration := time.Since(start)
	status := "success"
	if err != nil {
		status = "failure"
		RecordSpanError(spanCtx, err)
	}

	metrics := GetMetrics()
	if metrics != nil {
		metrics.ModuleQueryCounter.WithLabelValues("grpc", info).Inc()
		metrics.ModuleQueryDuration.WithLabelValues("grpc", info).Observe(duration.Seconds())
	}

	LogEvent("grpc_request", map[string]interface{}{
		"method":   info,
		"status":   status,
		"duration": duration.Seconds(),
	})

	return resp, err
}

// StartMetricsServer starts a metrics server with default configuration
func StartMetricsServer() (*MetricsServer, error) {
	config := DefaultMetricsServerConfig()
	server := NewMetricsServer(config)
	if server == nil {
		return nil, fmt.Errorf("metrics server is disabled")
	}

	if err := server.Start(); err != nil {
		return nil, err
	}

	return server, nil
}
