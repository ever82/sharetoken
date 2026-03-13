package types

import (
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StringCoin helpers for working with string-based coins
func zeroCoin() sdk.Coin {
	return sdk.NewCoin("ustt", sdk.ZeroInt())
}

// DefaultCoin returns a default coin string for limit initialization
func DefaultCoin() string {
	return sdk.NewCoin("ustt", sdk.ZeroInt()).String()
}

// DefaultDefaultLimitConfig returns default limit configuration
func DefaultDefaultLimitConfig() DefaultLimitConfig {
	now := uint64(time.Now().Unix())
	return DefaultLimitConfig{
		DefaultTxLimit: TransactionLimit{
			MaxSingle:      sdk.NewCoin("ustt", sdk.NewInt(1000000000)).String(),   // 1000 STT
			MaxDaily:       sdk.NewCoin("ustt", sdk.NewInt(10000000000)).String(),  // 10000 STT
			MaxMonthly:     sdk.NewCoin("ustt", sdk.NewInt(100000000000)).String(), // 100000 STT
			DailySpent:     DefaultCoin(),
			MonthlySpent:   DefaultCoin(),
			LastResetDay:   now,
			LastResetMonth: now,
		},
		DefaultWithdrawalLimit: WithdrawalLimit{
			MaxDaily:       sdk.NewCoin("ustt", sdk.NewInt(5000000000)).String(), // 5000 STT
			CooldownHours:  24,
			DailyWithdrawn: DefaultCoin(),
			LastResetDay:   now,
		},
		DefaultDisputeLimit: DisputeLimit{
			MaxActiveDisputes: 5,
			CurrentActive:     0,
		},
		DefaultServiceLimit: ServiceLimit{
			MaxConcurrent:      10,
			RateLimitPerMinute: 60,
			LastRateReset:      now,
		},
	}
}

// NewLimitConfig creates a new limit config with default values
func NewLimitConfig(address string) LimitConfig {
	defaults := DefaultDefaultLimitConfig()
	now := uint64(time.Now().Unix())

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

// parseCoin parses a coin string
func parseCoin(s string) sdk.Coin {
	coin, err := sdk.ParseCoinNormalized(s)
	if err != nil {
		return zeroCoin()
	}
	return coin
}

// CheckTransactionLimit checks if a transaction amount is within limits
func (lc *LimitConfig) CheckTransactionLimit(amount sdk.Coin) error {
	// Reset counters if needed
	lc.resetTxCountersIfNeeded()

	maxSingle := parseCoin(lc.TxLimit.MaxSingle)
	maxDaily := parseCoin(lc.TxLimit.MaxDaily)
	maxMonthly := parseCoin(lc.TxLimit.MaxMonthly)
	dailySpent := parseCoin(lc.TxLimit.DailySpent)
	monthlySpent := parseCoin(lc.TxLimit.MonthlySpent)

	// Check single transaction limit
	if amount.IsGTE(maxSingle) {
		return LimitError{
			LimitType:   "transaction_single",
			Current:     amount.String(),
			Max:         maxSingle.String(),
			Description: "transaction amount exceeds single transaction limit",
		}
	}

	// Check daily limit
	newDailySpent := dailySpent.Add(amount)
	if newDailySpent.IsGTE(maxDaily) {
		return LimitError{
			LimitType:   "transaction_daily",
			Current:     newDailySpent.String(),
			Max:         maxDaily.String(),
			Description: "transaction would exceed daily limit",
		}
	}

	// Check monthly limit
	newMonthlySpent := monthlySpent.Add(amount)
	if newMonthlySpent.IsGTE(maxMonthly) {
		return LimitError{
			LimitType:   "transaction_monthly",
			Current:     newMonthlySpent.String(),
			Max:         maxMonthly.String(),
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
	dailySpent := parseCoin(lc.TxLimit.DailySpent).Add(amount)
	monthlySpent := parseCoin(lc.TxLimit.MonthlySpent).Add(amount)
	lc.TxLimit.DailySpent = dailySpent.String()
	lc.TxLimit.MonthlySpent = monthlySpent.String()
}

// resetTxCountersIfNeeded resets daily/monthly counters if needed
func (lc *LimitConfig) resetTxCountersIfNeeded() {
	now := uint64(time.Now().Unix())
	nowDay := now / SecondsPerDay
	nowMonth := now / (SecondsPerDay * DaysPerMonth)

	lastDay := lc.TxLimit.LastResetDay / SecondsPerDay
	if nowDay > lastDay {
		lc.TxLimit.DailyTxCount = 0
		lc.TxLimit.DailySpent = DefaultCoin()
		lc.TxLimit.LastResetDay = now
	}

	lastMonth := lc.TxLimit.LastResetMonth / (SecondsPerDay * DaysPerMonth)
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

	maxDaily := parseCoin(lc.WithdrawalLimit.MaxDaily)
	dailyWithdrawn := parseCoin(lc.WithdrawalLimit.DailyWithdrawn)

	// Check cooldown
	now := time.Now().Unix()
	cooldownSeconds := int64(lc.WithdrawalLimit.CooldownHours) * SecondsPerHour
	if int64(lc.WithdrawalLimit.LastWithdrawalTime) > 0 && now-int64(lc.WithdrawalLimit.LastWithdrawalTime) < cooldownSeconds {
		remaining := cooldownSeconds - (now - int64(lc.WithdrawalLimit.LastWithdrawalTime))
		return LimitError{
			LimitType:   "withdrawal_cooldown",
			Current:     fmt.Sprintf("%d seconds since last withdrawal", now-int64(lc.WithdrawalLimit.LastWithdrawalTime)),
			Max:         fmt.Sprintf("%d seconds cooldown", cooldownSeconds),
			Description: fmt.Sprintf("cooldown period not met, remaining: %d seconds", remaining),
		}
	}

	// Check daily limit
	newDailyWithdrawn := dailyWithdrawn.Add(amount)
	if newDailyWithdrawn.IsGTE(maxDaily) {
		return LimitError{
			LimitType:   "withdrawal_daily",
			Current:     newDailyWithdrawn.String(),
			Max:         maxDaily.String(),
			Description: "withdrawal would exceed daily limit",
		}
	}

	return nil
}

// RecordWithdrawal records a completed withdrawal
func (lc *LimitConfig) RecordWithdrawal(amount sdk.Coin) {
	lc.resetWithdrawalCountersIfNeeded()
	lc.WithdrawalLimit.LastWithdrawalTime = uint64(time.Now().Unix())
	dailyWithdrawn := parseCoin(lc.WithdrawalLimit.DailyWithdrawn).Add(amount)
	lc.WithdrawalLimit.DailyWithdrawn = dailyWithdrawn.String()
}

// resetWithdrawalCountersIfNeeded resets daily counters if needed
func (lc *LimitConfig) resetWithdrawalCountersIfNeeded() {
	now := uint64(time.Now().Unix())
	nowDay := now / SecondsPerDay

	lastDay := lc.WithdrawalLimit.LastResetDay / SecondsPerDay
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
	now := uint64(time.Now().Unix())
	nowMinute := now / SecondsPerMinute

	lastMinute := lc.ServiceLimit.LastRateReset / SecondsPerMinute
	if nowMinute > lastMinute {
		lc.ServiceLimit.RequestsLastMinute = 0
		lc.ServiceLimit.LastRateReset = now
	}
}
