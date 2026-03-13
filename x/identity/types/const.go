package types

// User limit constants for verified and unverified users
// All amounts are in micro-STT (ustt), where 1 STT = 1,000,000 ustt
const (
	// STTDecimals is the number of decimal places for STT token
	STTDecimals = 6
	// STTBaseUnit is the base unit conversion factor (1 STT = 10^6 ustt)
	STTBaseUnit int64 = 1_000_000

	// Unverified user limits
	// UnverifiedUserTransactionLimit is the daily transaction limit for unverified users (1000 STT)
	UnverifiedUserTransactionLimit int64 = 1_000 * STTBaseUnit
	// UnverifiedUserWithdrawalLimit is the daily withdrawal limit for unverified users (500 STT)
	UnverifiedUserWithdrawalLimit int64 = 500 * STTBaseUnit
	// UnverifiedUserDisputeLimit is the dispute limit for unverified users (100 STT)
	UnverifiedUserDisputeLimit int64 = 100 * STTBaseUnit
	// UnverifiedUserServiceLimit is the service usage limit for unverified users (500 STT)
	UnverifiedUserServiceLimit int64 = 500 * STTBaseUnit

	// Verified user limits (10x unverified limits)
	// VerifiedUserTransactionLimit is the daily transaction limit for verified users (10000 STT)
	VerifiedUserTransactionLimit int64 = 10_000 * STTBaseUnit
	// VerifiedUserWithdrawalLimit is the daily withdrawal limit for verified users (5000 STT)
	VerifiedUserWithdrawalLimit int64 = 5_000 * STTBaseUnit
	// VerifiedUserDisputeLimit is the dispute limit for verified users (1000 STT)
	VerifiedUserDisputeLimit int64 = 1_000 * STTBaseUnit
	// VerifiedUserServiceLimit is the service usage limit for verified users (5000 STT)
	VerifiedUserServiceLimit int64 = 5_000 * STTBaseUnit
)

// Reputation level thresholds based on MQ score
const (
	// ReputationLevelOutstandingMin is the minimum MQ score for outstanding reputation level
	ReputationLevelOutstandingMin int64 = 500
	// ReputationLevelExcellentMin is the minimum MQ score for excellent reputation level
	ReputationLevelExcellentMin int64 = 200
	// ReputationLevelGoodMin is the minimum MQ score for good reputation level
	ReputationLevelGoodMin int64 = 100
	// ReputationLevelNormalMin is the minimum MQ score for normal reputation level
	ReputationLevelNormalMin int64 = 50
	// ReputationLevelNoviceMax is the maximum MQ score for novice reputation level
	ReputationLevelNoviceMax int64 = 49
)

// Juror eligibility constants
const (
	// JurorMinRequiredScore is the minimum MQ score required to be a juror
	JurorMinRequiredScore int = 100
)

// Pagination defaults
const (
	// DefaultPageSize is the default number of items per page
	DefaultPageSize = 10
	// MaxPageSize is the maximum number of items per page
	MaxPageSize = 100
)

// Dispute resolution constants
const (
	// DefaultUserSharePercent is the default percentage for user in dispute resolution
	DefaultUserSharePercent int64 = 70
	// DefaultProviderSharePercent is the default percentage for provider in dispute resolution
	DefaultProviderSharePercent int64 = 30
	// MaxSharePercent is the maximum percentage (100%)
	MaxSharePercent int64 = 100
)

// Timeout constants
const (
	// DefaultPacketTimeoutMinutes is the default IBC packet timeout in minutes
	DefaultPacketTimeoutMinutes = 10
	// NATRefreshIntervalMinutes is the NAT refresh interval in minutes
	NATRefreshIntervalMinutes = 30
	// CrowdfundingCheckIntervalMinutes is the crowdfunding campaign check interval
	CrowdfundingCheckIntervalMinutes = 1
)

// Gas and fee constants
const (
	// DefaultMinGasPrice is the default minimum gas price
	DefaultMinGasPrice = "0.025"
	// DefaultGasPriceMultiplier is the multiplier for calculating fees from gas limit
	DefaultGasPriceMultiplier = 25
	// DefaultGasPriceDivisor is the divisor for calculating fees (1/10000 = 0.0001)
	DefaultGasPriceDivisor = 10000
)

// Oracle price constants
const (
	// PriceSourceConfidenceThreshold is the minimum confidence level for valid prices
	PriceSourceConfidenceThreshold int32 = 90
	// MaxConfidence is the maximum confidence value (100%)
	MaxConfidence int32 = 100
	// TokenUnitDivisor is the divisor for token-based calculations (1000 tokens)
	TokenUnitDivisor int64 = 1000
	// DefaultSTTPrice is the default STT price in USD
	DefaultSTTPrice int64 = 100
)

// LLM custody constants
const (
	// DefaultRateLimit is the default API rate limit (requests per minute)
	DefaultRateLimit = 100
	// DefaultMaxRequests is the default maximum requests per day
	DefaultMaxRequests = 1000
	// MaxRateLimit is the maximum allowed rate limit
	MaxRateLimit = 1000
	// MaxMaxRequests is the maximum allowed max requests
	MaxMaxRequests = 10000
)

// Escrow constants
const (
	// DefaultEscrowDurationHours is the default escrow duration in hours
	DefaultEscrowDurationHours = 24
)

// API retry constants
const (
	// MaxAPIRetries is the maximum number of API retry attempts
	MaxAPIRetries = 5
	// InitialBackoffSeconds is the initial backoff duration in seconds for retries
	InitialBackoffSeconds = 1
)

// Network constants
const (
	// ChainIDRandLength is the length of random string appended to chain ID
	ChainIDRandLength = 6
)

// NAT/UPnP constants
const (
	// NATPortMappingLeaseSeconds is the UPnP port mapping lease duration in seconds
	NATPortMappingLeaseSeconds uint32 = 3600 // 1 hour lease, will be renewed by keepalive
)
