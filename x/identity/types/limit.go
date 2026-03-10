package types

import (
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LimitConfig stores user limits configuration
type LimitConfig struct {
	Address         string           `json:"address"`
	TxLimit         TransactionLimit `json:"tx_limit"`
	WithdrawalLimit WithdrawalLimit  `json:"withdrawal_limit"`
	DisputeLimit    DisputeLimit     `json:"dispute_limit"`
	ServiceLimit    ServiceLimit     `json:"service_limit"`
	UpdatedAt       int64            `json:"updated_at"`
}

// TransactionLimit defines transaction related limits
type TransactionLimit struct {
	MaxSingle      sdk.Coin `json:"max_single"`
	MaxDaily       sdk.Coin `json:"max_daily"`
	MaxMonthly     sdk.Coin `json:"max_monthly"`
	DailyTxCount   uint64   `json:"daily_tx_count"`
	MonthlyTxCount uint64   `json:"monthly_tx_count"`
	DailySpent     sdk.Coin `json:"daily_spent"`
	MonthlySpent   sdk.Coin `json:"monthly_spent"`
	LastResetDay   int64    `json:"last_reset_day"`
	LastResetMonth int64    `json:"last_reset_month"`
}

// WithdrawalLimit defines withdrawal related limits
type WithdrawalLimit struct {
	MaxDaily           sdk.Coin `json:"max_daily"`
	CooldownHours      uint64   `json:"cooldown_hours"`
	LastWithdrawalTime int64    `json:"last_withdrawal_time"`
	DailyWithdrawn     sdk.Coin `json:"daily_withdrawn"`
	LastResetDay       int64    `json:"last_reset_day"`
}

// DisputeLimit defines dispute related limits
type DisputeLimit struct {
	MaxActiveDisputes uint64 `json:"max_active_disputes"`
	CurrentActive     uint64 `json:"current_active"`
}

// ServiceLimit defines service usage limits
type ServiceLimit struct {
	MaxConcurrent      uint64 `json:"max_concurrent"`
	RateLimitPerMinute uint64 `json:"rate_limit_per_minute"`
	CurrentConcurrent  uint64 `json:"current_concurrent"`
	RequestsLastMinute uint64 `json:"requests_last_minute"`
	LastRateReset      int64  `json:"last_rate_reset"`
}

// DefaultLimitConfig stores default limit values
type DefaultLimitConfig struct {
	DefaultTxLimit         TransactionLimit `json:"default_tx_limit"`
	DefaultWithdrawalLimit WithdrawalLimit  `json:"default_withdrawal_limit"`
	DefaultDisputeLimit    DisputeLimit     `json:"default_dispute_limit"`
	DefaultServiceLimit    ServiceLimit     `json:"default_service_limit"`
}

// DefaultCoin returns a default coin for limit initialization
func DefaultCoin() sdk.Coin {
	return sdk.NewCoin("ustt", sdk.ZeroInt())
}

// DefaultDefaultLimitConfig returns default limit configuration
func DefaultDefaultLimitConfig() DefaultLimitConfig {
	return DefaultLimitConfig{
		DefaultTxLimit: TransactionLimit{
			MaxSingle:      sdk.NewCoin("ustt", sdk.NewInt(1000000000)),   // 1000 STT
			MaxDaily:       sdk.NewCoin("ustt", sdk.NewInt(10000000000)),  // 10000 STT
			MaxMonthly:     sdk.NewCoin("ustt", sdk.NewInt(100000000000)), // 100000 STT
			DailySpent:     DefaultCoin(),
			MonthlySpent:   DefaultCoin(),
			LastResetDay:   time.Now().Unix(),
			LastResetMonth: time.Now().Unix(),
		},
		DefaultWithdrawalLimit: WithdrawalLimit{
			MaxDaily:       sdk.NewCoin("ustt", sdk.NewInt(5000000000)), // 5000 STT
			CooldownHours:  24,
			DailyWithdrawn: DefaultCoin(),
			LastResetDay:   time.Now().Unix(),
		},
		DefaultDisputeLimit: DisputeLimit{
			MaxActiveDisputes: 5,
			CurrentActive:     0,
		},
		DefaultServiceLimit: ServiceLimit{
			MaxConcurrent:      10,
			RateLimitPerMinute: 60,
			LastRateReset:      time.Now().Unix(),
		},
	}
}

// NewLimitConfig creates a new limit config with default values
func NewLimitConfig(address string) LimitConfig {
	defaults := DefaultDefaultLimitConfig()
	now := time.Now().Unix()

	return LimitConfig{
		Address: address,
		TxLimit: TransactionLimit{
			MaxSingle:      defaults.DefaultTxLimit.MaxSingle,
			MaxDaily:       defaults.DefaultTxLimit.MaxDaily,
			MaxMonthly:     defaults.DefaultTxLimit.MaxMonthly,
			DailySpent:     DefaultCoin(),
			MonthlySpent:   DefaultCoin(),
			LastResetDay:   now,
			LastResetMonth: now,
		},
		WithdrawalLimit: WithdrawalLimit{
			MaxDaily:       defaults.DefaultWithdrawalLimit.MaxDaily,
			CooldownHours:  defaults.DefaultWithdrawalLimit.CooldownHours,
			DailyWithdrawn: DefaultCoin(),
			LastResetDay:   now,
		},
		DisputeLimit: DisputeLimit{
			MaxActiveDisputes: defaults.DefaultDisputeLimit.MaxActiveDisputes,
			CurrentActive:     0,
		},
		ServiceLimit: ServiceLimit{
			MaxConcurrent:      defaults.DefaultServiceLimit.MaxConcurrent,
			RateLimitPerMinute: defaults.DefaultServiceLimit.RateLimitPerMinute,
			LastRateReset:      now,
		},
		UpdatedAt: now,
	}
}

// ValidateBasic performs basic validation
func (lc LimitConfig) ValidateBasic() error {
	if lc.Address == "" {
		return ErrInvalidAddress
	}

	_, err := sdk.AccAddressFromBech32(lc.Address)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}

	return nil
}

// CheckTransactionLimit checks if a transaction amount is within limits
func (lc *LimitConfig) CheckTransactionLimit(amount sdk.Coin) error {
	// Reset counters if needed
	lc.resetTxCountersIfNeeded()

	// Check single transaction limit
	if amount.IsGTE(lc.TxLimit.MaxSingle) {
		return LimitError{
			LimitType:   "transaction_single",
			Current:     amount.String(),
			Max:         lc.TxLimit.MaxSingle.String(),
			Description: "transaction amount exceeds single transaction limit",
		}
	}

	// Check daily limit
	newDailySpent := lc.TxLimit.DailySpent.Add(amount)
	if newDailySpent.IsGTE(lc.TxLimit.MaxDaily) {
		return LimitError{
			LimitType:   "transaction_daily",
			Current:     newDailySpent.String(),
			Max:         lc.TxLimit.MaxDaily.String(),
			Description: "transaction would exceed daily limit",
		}
	}

	// Check monthly limit
	newMonthlySpent := lc.TxLimit.MonthlySpent.Add(amount)
	if newMonthlySpent.IsGTE(lc.TxLimit.MaxMonthly) {
		return LimitError{
			LimitType:   "transaction_monthly",
			Current:     newMonthlySpent.String(),
			Max:         lc.TxLimit.MaxMonthly.String(),
			Description: "transaction would exceed monthly limit",
		}
	}

	return nil
}

// RecordTransaction records a completed transaction
func (lc *LimitConfig) RecordTransaction(amount sdk.Coin) {
	lc.resetTxCountersIfNeeded()
	lc.TxLimit.DailyTxCount++
	lc.TxLimit.MonthlyTxCount++
	lc.TxLimit.DailySpent = lc.TxLimit.DailySpent.Add(amount)
	lc.TxLimit.MonthlySpent = lc.TxLimit.MonthlySpent.Add(amount)
}

// resetTxCountersIfNeeded resets daily/monthly counters if needed
func (lc *LimitConfig) resetTxCountersIfNeeded() {
	now := time.Now().Unix()
	nowDay := now / 86400
	nowMonth := now / (86400 * 30)

	lastDay := lc.TxLimit.LastResetDay / 86400
	if nowDay > lastDay {
		lc.TxLimit.DailyTxCount = 0
		lc.TxLimit.DailySpent = DefaultCoin()
		lc.TxLimit.LastResetDay = now
	}

	lastMonth := lc.TxLimit.LastResetMonth / (86400 * 30)
	if nowMonth > lastMonth {
		lc.TxLimit.MonthlyTxCount = 0
		lc.TxLimit.MonthlySpent = DefaultCoin()
		lc.TxLimit.LastResetMonth = now
	}
}

// CheckWithdrawalLimit checks if a withdrawal is within limits
func (lc *LimitConfig) CheckWithdrawalLimit(amount sdk.Coin) error {
	// Reset counters if needed
	lc.resetWithdrawalCountersIfNeeded()

	// Check cooldown
	now := time.Now().Unix()
	cooldownSeconds := int64(lc.WithdrawalLimit.CooldownHours) * 3600
	if now-lc.WithdrawalLimit.LastWithdrawalTime < cooldownSeconds {
		remaining := cooldownSeconds - (now - lc.WithdrawalLimit.LastWithdrawalTime)
		return LimitError{
			LimitType:   "withdrawal_cooldown",
			Current:     fmt.Sprintf("%d seconds since last withdrawal", now-lc.WithdrawalLimit.LastWithdrawalTime),
			Max:         fmt.Sprintf("%d seconds cooldown", cooldownSeconds),
			Description: fmt.Sprintf("cooldown period not met, remaining: %d seconds", remaining),
		}
	}

	// Check daily limit
	newDailyWithdrawn := lc.WithdrawalLimit.DailyWithdrawn.Add(amount)
	if newDailyWithdrawn.IsGTE(lc.WithdrawalLimit.MaxDaily) {
		return LimitError{
			LimitType:   "withdrawal_daily",
			Current:     newDailyWithdrawn.String(),
			Max:         lc.WithdrawalLimit.MaxDaily.String(),
			Description: "withdrawal would exceed daily limit",
		}
	}

	return nil
}

// RecordWithdrawal records a completed withdrawal
func (lc *LimitConfig) RecordWithdrawal(amount sdk.Coin) {
	lc.resetWithdrawalCountersIfNeeded()
	lc.WithdrawalLimit.LastWithdrawalTime = time.Now().Unix()
	lc.WithdrawalLimit.DailyWithdrawn = lc.WithdrawalLimit.DailyWithdrawn.Add(amount)
}

// resetWithdrawalCountersIfNeeded resets daily counters if needed
func (lc *LimitConfig) resetWithdrawalCountersIfNeeded() {
	now := time.Now().Unix()
	nowDay := now / 86400

	lastDay := lc.WithdrawalLimit.LastResetDay / 86400
	if nowDay > lastDay {
		lc.WithdrawalLimit.DailyWithdrawn = DefaultCoin()
		lc.WithdrawalLimit.LastResetDay = now
	}
}

// CheckDisputeLimit checks if a new dispute can be created
func (lc *LimitConfig) CheckDisputeLimit() error {
	if lc.DisputeLimit.CurrentActive >= lc.DisputeLimit.MaxActiveDisputes {
		return LimitError{
			LimitType:   "dispute_active",
			Current:     strconv.FormatUint(lc.DisputeLimit.CurrentActive, 10),
			Max:         strconv.FormatUint(lc.DisputeLimit.MaxActiveDisputes, 10),
			Description: "maximum number of active disputes reached",
		}
	}
	return nil
}

// IncrementActiveDisputes increments the active dispute count
func (lc *LimitConfig) IncrementActiveDisputes() {
	lc.DisputeLimit.CurrentActive++
}

// DecrementActiveDisputes decrements the active dispute count
func (lc *LimitConfig) DecrementActiveDisputes() {
	if lc.DisputeLimit.CurrentActive > 0 {
		lc.DisputeLimit.CurrentActive--
	}
}

// CheckServiceLimit checks if a service call is within limits
func (lc *LimitConfig) CheckServiceLimit() error {
	// Reset rate limit if needed
	lc.resetServiceCountersIfNeeded()

	// Check concurrent limit
	if lc.ServiceLimit.CurrentConcurrent >= lc.ServiceLimit.MaxConcurrent {
		return LimitError{
			LimitType:   "service_concurrent",
			Current:     strconv.FormatUint(lc.ServiceLimit.CurrentConcurrent, 10),
			Max:         strconv.FormatUint(lc.ServiceLimit.MaxConcurrent, 10),
			Description: "maximum concurrent service calls reached",
		}
	}

	// Check rate limit
	if lc.ServiceLimit.RequestsLastMinute >= lc.ServiceLimit.RateLimitPerMinute {
		return LimitError{
			LimitType:   "service_rate",
			Current:     strconv.FormatUint(lc.ServiceLimit.RequestsLastMinute, 10),
			Max:         strconv.FormatUint(lc.ServiceLimit.RateLimitPerMinute, 10),
			Description: "rate limit exceeded (requests per minute)",
		}
	}

	return nil
}

// RecordServiceCall records a service call
func (lc *LimitConfig) RecordServiceCall() {
	lc.resetServiceCountersIfNeeded()
	lc.ServiceLimit.CurrentConcurrent++
	lc.ServiceLimit.RequestsLastMinute++
}

// ReleaseServiceCall releases a service call slot
func (lc *LimitConfig) ReleaseServiceCall() {
	if lc.ServiceLimit.CurrentConcurrent > 0 {
		lc.ServiceLimit.CurrentConcurrent--
	}
}

// resetServiceCountersIfNeeded resets rate limit counters if needed
func (lc *LimitConfig) resetServiceCountersIfNeeded() {
	now := time.Now().Unix()
	nowMinute := now / 60

	lastMinute := lc.ServiceLimit.LastRateReset / 60
	if nowMinute > lastMinute {
		lc.ServiceLimit.RequestsLastMinute = 0
		lc.ServiceLimit.LastRateReset = now
	}
}

// String implements stringer interface
func (lc LimitConfig) String() string {
	return fmt.Sprintf(`
LimitConfig for %s:
  Transaction:
    Single: %s
    Daily: %s
    Monthly: %s
    Daily Count: %d
    Monthly Count: %d
  Withdrawal:
    Daily: %s
    Cooldown: %d hours
  Dispute:
    Max Active: %d
    Current Active: %d
  Service:
    Max Concurrent: %d
    Rate Limit: %d/min
    Current Concurrent: %d
`, lc.Address,
		lc.TxLimit.MaxSingle.String(), lc.TxLimit.MaxDaily.String(), lc.TxLimit.MaxMonthly.String(),
		lc.TxLimit.DailyTxCount, lc.TxLimit.MonthlyTxCount,
		lc.WithdrawalLimit.MaxDaily.String(), lc.WithdrawalLimit.CooldownHours,
		lc.DisputeLimit.MaxActiveDisputes, lc.DisputeLimit.CurrentActive,
		lc.ServiceLimit.MaxConcurrent, lc.ServiceLimit.RateLimitPerMinute, lc.ServiceLimit.CurrentConcurrent)
}
