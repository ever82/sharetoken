package telemetry

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCIHooks wraps ABCI calls with telemetry
// This implements the cosmos-sdk/types/module.Hooks pattern
type ABCIHooks struct {
	module string
}

// NewABCIHooks creates new ABCI hooks for a module
func NewABCIHooks(module string) *ABCIHooks {
	return &ABCIHooks{module: module}
}

// BeforeBeginBlock is called before BeginBlock
func (h *ABCIHooks) BeforeBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) sdk.Context {
	// Record the start time in context
	return ctx.WithContext(context.WithValue(ctx.Context(), "begin_block_start", time.Now()))
}

// AfterBeginBlock is called after BeginBlock
func (h *ABCIHooks) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, resp abci.ResponseBeginBlock) {
	// Calculate and record duration
	if start, ok := ctx.Context().Value("begin_block_start").(time.Time); ok {
		duration := time.Since(start)
		GetMetrics().ObserveBeginBlockDuration(duration)
		AddTelemetryEvent(ctx.Context(), "begin_block_complete",
			StringAttr("module", h.module),
			Int64Attr("duration_ms", duration.Milliseconds()),
		)
	}

	// Record block height
	GetMetrics().SetBlockHeight(ctx.BlockHeight())

	// Log the event
	GetLogger().LogABCI("begin_block", map[string]interface{}{
		"module":       h.module,
		"block_height": ctx.BlockHeight(),
		"block_time":   ctx.BlockTime(),
	})
}

// BeforeEndBlock is called before EndBlock
func (h *ABCIHooks) BeforeEndBlock(ctx sdk.Context, req abci.RequestEndBlock) sdk.Context {
	return ctx.WithContext(context.WithValue(ctx.Context(), "end_block_start", time.Now()))
}

// AfterEndBlock is called after EndBlock
func (h *ABCIHooks) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, resp abci.ResponseEndBlock) {
	if start, ok := ctx.Context().Value("end_block_start").(time.Time); ok {
		duration := time.Since(start)
		GetMetrics().ObserveEndBlockDuration(duration)
		AddTelemetryEvent(ctx.Context(), "end_block_complete",
			StringAttr("module", h.module),
			Int64Attr("duration_ms", duration.Milliseconds()),
		)
	}

	GetLogger().LogABCI("end_block", map[string]interface{}{
		"module":       h.module,
		"block_height": ctx.BlockHeight(),
	})
}

// BeforeDeliverTx is called before DeliverTx
func (h *ABCIHooks) BeforeDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx) sdk.Context {
	return ctx.WithContext(context.WithValue(ctx.Context(), "deliver_tx_start", time.Now()))
}

// AfterDeliverTx is called after DeliverTx
func (h *ABCIHooks) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, resp abci.ResponseDeliverTx) {
	if start, ok := ctx.Context().Value("deliver_tx_start").(time.Time); ok {
		duration := time.Since(start)
		GetMetrics().ObserveDeliverTxDuration(duration)

		status := "success"
		if resp.Code != 0 {
			status = "failure"
		}

		GetMetrics().IncTransaction(status, h.module)
		GetMetrics().ObserveTransactionGasUsed(h.module, uint64(resp.GasUsed))

		AddTelemetryEvent(ctx.Context(), "deliver_tx_complete",
			StringAttr("module", h.module),
			StringAttr("status", status),
			Int64Attr("gas_used", resp.GasUsed),
			Int64Attr("duration_ms", duration.Milliseconds()),
		)
	}

	GetLogger().LogABCI("deliver_tx", map[string]interface{}{
		"module":    h.module,
		"gas_used":  resp.GasUsed,
		"gas_wanted": resp.GasWanted,
		"code":      resp.Code,
	})
}

// TelemetryContext wraps SDK context with telemetry
type TelemetryContext struct {
	sdk.Context
	startTime time.Time
	module    string
	span      interface{} // trace.Span
}

// NewTelemetryContext creates a new telemetry context
func NewTelemetryContext(ctx sdk.Context, module string) *TelemetryContext {
	// Start span for tracing
	spanCtx, span := StartModuleSpan(ctx.Context(), module, "operation")

	return &TelemetryContext{
		Context:   ctx.WithContext(spanCtx),
		startTime: time.Now(),
		module:    module,
		span:      span,
	}
}

// Finish records metrics and ends the span
func (tc *TelemetryContext) Finish(operation string) {
	duration := time.Since(tc.startTime)

	// Record metrics
	GetMetrics().ObserveModuleQueryDuration(tc.module, operation, duration)

	// End span
	if span, ok := tc.span.(interface{ End() }); ok {
		span.End()
	}
}

// RecordTx records transaction metrics
func (tc *TelemetryContext) RecordTx(msgType string, gasUsed uint64, success bool) {
	duration := time.Since(tc.startTime)

	status := "success"
	if !success {
		status = "failure"
	}

	// Record metrics
	GetMetrics().IncModuleTx(tc.module, msgType)
	GetMetrics().IncTransaction(status, tc.module)
	GetMetrics().ObserveTransactionDuration(tc.module, duration)
	GetMetrics().ObserveTransactionGasUsed(tc.module, gasUsed)

	// Log
	GetLogger().LogTx(tc.module, msgType, map[string]interface{}{
		"gas_used": gasUsed,
		"success":  success,
		"duration": duration.Seconds(),
	})

	// End span
	if span, ok := tc.span.(interface{ End() }); ok {
		span.End()
	}
}

// RecordQuery records query metrics
func (tc *TelemetryContext) RecordQuery(queryType string) {
	duration := time.Since(tc.startTime)

	// Record metrics
	GetMetrics().IncModuleQuery(tc.module, queryType)
	GetMetrics().ObserveModuleQueryDuration(tc.module, queryType, duration)

	// Log
	GetLogger().LogQuery(tc.module, queryType, map[string]interface{}{
		"duration": duration.Seconds(),
	})

	// End span
	if span, ok := tc.span.(interface{ End() }); ok {
		span.End()
	}
}

// Attribute helpers for tracing

// StringAttr creates a string attribute
func StringAttr(key, value string) struct {
	Key   string
	Value string
} {
	return struct {
		Key   string
		Value string
	}{Key: key, Value: value}
}

// Int64Attr creates an int64 attribute
func Int64Attr(key string, value int64) struct {
	Key   string
	Value int64
} {
	return struct {
		Key   string
		Value int64
	}{Key: key, Value: value}
}

// TelemetryEvent represents a telemetry event
type TelemetryEvent struct {
	Name  string
	Attrs map[string]interface{}
}

// AddTelemetryEvent adds an event to the span (simplified)
func AddTelemetryEvent(ctx context.Context, name string, attrs ...interface{}) {
	// This is a simplified version - in production, use proper OpenTelemetry attributes
	GetLogger().LogEvent(name, map[string]interface{}{
		"attrs": attrs,
	})
}

// ModuleKeeper is an interface for module keepers that support telemetry
type ModuleKeeper interface {
	SetTelemetry(telemetry *TelemetryKeeper)
}

// TelemetryKeeper wraps a module keeper with telemetry
type TelemetryKeeper struct {
	Module   string
	Metrics  *Metrics
	Logger   *Logger
	tracer   interface{} // trace.Tracer
}

// NewTelemetryKeeper creates a new telemetry keeper
func NewTelemetryKeeper(module string) *TelemetryKeeper {
	return &TelemetryKeeper{
		Module:  module,
		Metrics: GetMetrics(),
		Logger:  GetLogger().WithModule(module),
	}
}

// RecordTxDuration records transaction duration
func (tk *TelemetryKeeper) RecordTxDuration(start time.Time) {
	tk.Metrics.ObserveTransactionDuration(tk.Module, time.Since(start))
}

// RecordQueryDuration records query duration
func (tk *TelemetryKeeper) RecordQueryDuration(queryType string, start time.Time) {
	tk.Metrics.ObserveModuleQueryDuration(tk.Module, queryType, time.Since(start))
}

// LogInfo logs an info message
func (tk *TelemetryKeeper) LogInfo(msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["module"] = tk.Module
	tk.Logger.LogEvent("info", fields)
}

// LogError logs an error message
func (tk *TelemetryKeeper) LogError(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["module"] = tk.Module
	fields["error"] = err.Error()
	tk.Logger.WithError(err).LogEvent("error", fields)
}

// StartSpan starts a new span for the module
func (tk *TelemetryKeeper) StartSpan(ctx context.Context, operation string) (context.Context, interface{}) {
	return StartModuleSpan(ctx, tk.Module, operation)
}

// TelemetryMiddleware wraps keeper methods with telemetry
func TelemetryMiddleware(module string, next func(ctx sdk.Context) error) func(ctx sdk.Context) error {
	return func(ctx sdk.Context) error {
		start := time.Now()
		spanCtx, span := StartModuleSpan(ctx.Context(), module, "operation")
		defer span.End()

		err := next(ctx.WithContext(spanCtx))

		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "failure"
			RecordSpanError(spanCtx, err)
		}

		GetMetrics().ObserveTransactionDuration(module, duration)
		GetMetrics().IncTransaction(status, module)

		return err
	}
}
