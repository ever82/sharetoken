package security

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// RateLimitTier defines different rate limit tiers
type RateLimitTier int

const (
	// TierStrict is for strict rate limiting (e.g., authentication endpoints)
	TierStrict RateLimitTier = iota
	// TierStandard is for standard API endpoints
	TierStandard
	// TierRelaxed is for read-only or less sensitive endpoints
	TierRelaxed
)

// RateLimitConfig defines rate limit configuration
type RateLimitConfig struct {
	RequestsPerWindow int
	WindowSize        time.Duration
	BurstSize         int
}

// Default rate limit configurations per tier
var DefaultRateLimits = map[RateLimitTier]RateLimitConfig{
	TierStrict: {
		RequestsPerWindow: 10,
		WindowSize:        time.Minute,
		BurstSize:         5,
	},
	TierStandard: {
		RequestsPerWindow: 60,
		WindowSize:        time.Minute,
		BurstSize:         10,
	},
	TierRelaxed: {
		RequestsPerWindow: 300,
		WindowSize:        time.Minute,
		BurstSize:         30,
	},
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	mu           sync.RWMutex
	limits       map[string]*RateLimitEntry // key -> entry
	config       RateLimitConfig
	keyExtractor func(interface{}) string
}

// RateLimitEntry tracks rate limit state for a key
type RateLimitEntry struct {
	Count       int
	WindowStart time.Time
	LastAccess  time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		limits: make(map[string]*RateLimitEntry),
		config: config,
	}
}

// NewRateLimiterWithTier creates a rate limiter with a predefined tier
func NewRateLimiterWithTier(tier RateLimitTier) *RateLimiter {
	config, ok := DefaultRateLimits[tier]
	if !ok {
		config = DefaultRateLimits[TierStandard]
	}
	return NewRateLimiter(config)
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.limits[key]

	if !exists || now.Sub(entry.WindowStart) > rl.config.WindowSize {
		// New window
		rl.limits[key] = &RateLimitEntry{
			Count:       1,
			WindowStart: now,
			LastAccess:  now,
		}
		return true
	}

	// Check burst limit
	if entry.Count >= rl.config.RequestsPerWindow {
		return false
	}

	// Allow request
	entry.Count++
	entry.LastAccess = now
	return true
}

// Check checks if a request would be allowed without incrementing the counter
func (rl *RateLimiter) Check(key string) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	entry, exists := rl.limits[key]

	if !exists || now.Sub(entry.WindowStart) > rl.config.WindowSize {
		return true // New window, would be allowed
	}

	return entry.Count < rl.config.RequestsPerWindow
}

// GetRemaining returns the remaining requests for a key
func (rl *RateLimiter) GetRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	entry, exists := rl.limits[key]

	if !exists || now.Sub(entry.WindowStart) > rl.config.WindowSize {
		return rl.config.RequestsPerWindow
	}

	remaining := rl.config.RequestsPerWindow - entry.Count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetResetTime returns when the rate limit will reset for a key
func (rl *RateLimiter) GetResetTime(key string) time.Time {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	entry, exists := rl.limits[key]
	if !exists {
		return time.Now()
	}

	return entry.WindowStart.Add(rl.config.WindowSize)
}

// Cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, entry := range rl.limits {
		if now.Sub(entry.LastAccess) > rl.config.WindowSize*2 {
			delete(rl.limits, key)
		}
	}
}

// IPRateLimiter provides rate limiting by IP address
type IPRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*RateLimiter // IP -> limiter
	config   RateLimitConfig
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(config RateLimitConfig) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		config:   config,
	}
}

// NewIPRateLimiterWithTier creates an IP rate limiter with a predefined tier
func NewIPRateLimiterWithTier(tier RateLimitTier) *IPRateLimiter {
	config, ok := DefaultRateLimits[tier]
	if !ok {
		config = DefaultRateLimits[TierStandard]
	}
	return NewIPRateLimiter(config)
}

// Allow checks if a request from an IP should be allowed
func (iprl *IPRateLimiter) Allow(ip string) bool {
	iprl.mu.Lock()
	if _, exists := iprl.limiters[ip]; !exists {
		iprl.limiters[ip] = NewRateLimiter(iprl.config)
	}
	limiter := iprl.limiters[ip]
	iprl.mu.Unlock()

	return limiter.Allow(ip)
}

// Check checks if a request would be allowed without incrementing
func (iprl *IPRateLimiter) Check(ip string) bool {
	iprl.mu.RLock()
	limiter, exists := iprl.limiters[ip]
	iprl.mu.RUnlock()

	if !exists {
		return true
	}

	return limiter.Check(ip)
}

// GetRemaining returns remaining requests for an IP
func (iprl *IPRateLimiter) GetRemaining(ip string) int {
	iprl.mu.RLock()
	limiter, exists := iprl.limiters[ip]
	iprl.mu.RUnlock()

	if !exists {
		return iprl.config.RequestsPerWindow
	}

	return limiter.GetRemaining(ip)
}

// Cleanup removes old IP entries
func (iprl *IPRateLimiter) Cleanup() {
	iprl.mu.Lock()
	defer iprl.mu.Unlock()

	for ip, limiter := range iprl.limiters {
		limiter.Cleanup()
		// If limiter has no entries, remove it
		if len(limiter.limits) == 0 {
			delete(iprl.limiters, ip)
		}
	}
}

// UserRateLimiter provides rate limiting by user ID
type UserRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*RateLimiter // userID -> limiter
	config   RateLimitConfig
}

// NewUserRateLimiter creates a new user-based rate limiter
func NewUserRateLimiter(config RateLimitConfig) *UserRateLimiter {
	return &UserRateLimiter{
		limiters: make(map[string]*RateLimiter),
		config:   config,
	}
}

// NewUserRateLimiterWithTier creates a user rate limiter with a predefined tier
func NewUserRateLimiterWithTier(tier RateLimitTier) *UserRateLimiter {
	config, ok := DefaultRateLimits[tier]
	if !ok {
		config = DefaultRateLimits[TierStandard]
	}
	return NewUserRateLimiter(config)
}

// Allow checks if a request from a user should be allowed
func (url *UserRateLimiter) Allow(userID string) bool {
	url.mu.Lock()
	if _, exists := url.limiters[userID]; !exists {
		url.limiters[userID] = NewRateLimiter(url.config)
	}
	limiter := url.limiters[userID]
	url.mu.Unlock()

	return limiter.Allow(userID)
}

// Check checks if a request would be allowed without incrementing
func (url *UserRateLimiter) Check(userID string) bool {
	url.mu.RLock()
	limiter, exists := url.limiters[userID]
	url.mu.RUnlock()

	if !exists {
		return true
	}

	return limiter.Check(userID)
}

// GetRemaining returns remaining requests for a user
func (url *UserRateLimiter) GetRemaining(userID string) int {
	url.mu.RLock()
	limiter, exists := url.limiters[userID]
	url.mu.RUnlock()

	if !exists {
		return url.config.RequestsPerWindow
	}

	return limiter.GetRemaining(userID)
}

// CompositeRateLimiter combines IP and user rate limiting
type CompositeRateLimiter struct {
	ipLimiter   *IPRateLimiter
	userLimiter *UserRateLimiter
}

// NewCompositeRateLimiter creates a composite rate limiter
func NewCompositeRateLimiter(ipConfig, userConfig RateLimitConfig) *CompositeRateLimiter {
	return &CompositeRateLimiter{
		ipLimiter:   NewIPRateLimiter(ipConfig),
		userLimiter: NewUserRateLimiter(userConfig),
	}
}

// Allow checks both IP and user rate limits
func (crl *CompositeRateLimiter) Allow(ip, userID string) (bool, string) {
	// Check IP limit first
	if !crl.ipLimiter.Allow(ip) {
		return false, "ip_rate_limited"
	}

	// Check user limit if userID is provided
	if userID != "" && !crl.userLimiter.Allow(userID) {
		return false, "user_rate_limited"
	}

	return true, ""
}

// RateLimitInfo provides rate limit information
type RateLimitInfo struct {
	IPRemaining    int
	UserRemaining  int
	IPResetTime    time.Time
	UserResetTime  time.Time
	Limited        bool
	Reason         string
}

// GetRateLimitInfo returns rate limit information for an IP and user
func (crl *CompositeRateLimiter) GetRateLimitInfo(ip, userID string) RateLimitInfo {
	info := RateLimitInfo{
		IPRemaining: crl.ipLimiter.GetRemaining(ip),
	}

	if userID != "" {
		info.UserRemaining = crl.userLimiter.GetRemaining(userID)
	}

	return info
}

// ParseIP extracts a clean IP address from a string
func ParseIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// Try to extract IP from host:port format
		host, _, err := net.SplitHostPort(ipStr)
		if err == nil {
			ip = net.ParseIP(host)
			if ip != nil {
				return ip.String()
			}
		}
		return ipStr // Return original if parsing fails
	}
	return ip.String()
}

// ValidateIP validates that a string is a valid IP address
func ValidateIP(ipStr string) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}
	return nil
}

// IsPrivateIP checks if an IP is a private address
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
	}

	for _, cidr := range privateRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
