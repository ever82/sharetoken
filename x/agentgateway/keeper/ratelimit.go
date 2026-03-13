package keeper

import (
	"time"

	"sharetoken/x/agentgateway/types"
)

// RateLimit 速率限制状态
type RateLimit struct {
	Address      string
	RequestCount int
	LastReset    time.Time
}

// CheckRateLimit 检查速率限制
func (k *Keeper) CheckRateLimit(address string) bool {
	k.rateMu.Lock()
	defer k.rateMu.Unlock()

	now := time.Now()
	limit, exists := k.rateLimits[address]
	if !exists {
		// 新建速率限制记录
		k.rateLimits[address] = &RateLimit{
			Address:      address,
			RequestCount: 1,
			LastReset:    now,
		}
		return true
	}

	// 检查是否需要重置（每分钟）
	if now.Sub(limit.LastReset) > time.Minute {
		limit.RequestCount = 1
		limit.LastReset = now
		return true
	}

	// 检查是否超过限制（60 req/min）
	if limit.RequestCount >= types.DefaultRateLimitPerMinute {
		return false
	}

	limit.RequestCount++
	return true
}
