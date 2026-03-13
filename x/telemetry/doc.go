// Package telemetry provides observability infrastructure for ShareToken blockchain.
//
// It includes:
// - Prometheus metrics collection
// - OpenTelemetry tracing (Jaeger/Zipkin)
// - Structured logging with JSON output
// - Module wrappers for automatic telemetry collection
//
// Basic Usage:
//
// Initialize telemetry at application startup:
//
//	// Initialize logger
//	telemetry.InitLogger(telemetry.LoggerConfig{
//	    Format: "json",
//	    Level:  telemetry.InfoLevel,
//	})
//
//	// Initialize metrics
//	telemetry.InitMetrics()
//
//	// Initialize tracing
//	telemetry.InitTracer(telemetry.TracerConfig{
//	    Enabled:        true,
//	    ExporterType:   "jaeger",
//	    JaegerEndpoint: "http://localhost:14268/api/traces",
//	})
//
//	// Start metrics server
//	server, err := telemetry.StartMetricsServer()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer server.Stop(context.Background())
//
// Using Metrics:
//
//	// Record block height
//	telemetry.GetMetrics().SetBlockHeight(ctx.BlockHeight())
//
//	// Record transaction
//	telemetry.GetMetrics().IncTransaction("success", "bank")
//
//	// Record query duration
//	defer telemetry.RecordQueryDuration("my-module", "my-query", time.Now())()
//
// Using Tracing:
//
//	// Start a span
//	ctx, span := telemetry.StartModuleSpan(ctx, "my-module", "my-operation")
//	defer span.End()
//
//	// Add attributes
//	telemetry.SetSpanAttributes(ctx, telemetry.BlockHeightAttribute(ctx.BlockHeight()))
//
//	// Record error
//	if err != nil {
//	    telemetry.RecordSpanError(ctx, err)
//	}
//
// Using Logging:
//
//	// Simple logging
//	telemetry.Info("Processing transaction")
//
//	// Structured logging
//	telemetry.LogEvent("transaction_processed", map[string]interface{}{
//	    "tx_hash": hash,
//	    "sender":  sender,
//	    "amount":  amount,
//	})
//
//	// Module-specific logging
//	logger := telemetry.WithModule("bank")
//	logger.Info("Transfer completed")
//
// Module Wrapper:
//
// Wrap existing modules with telemetry:
//
//	wrappedBank := telemetry.WrapModule("bank", bankModule)
//
// This automatically instruments BeginBlock, EndBlock, queries, and transactions.
//
package telemetry
