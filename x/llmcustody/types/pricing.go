package types

import (
	"fmt"
	"time"
)

// PricingConfig 定价配置
type PricingConfig struct {
	ID          string      `json:"id"`
	Owner       string      `json:"owner"`
	APIKeyID    string      `json:"api_key_id"`
	Provider    Provider    `json:"provider"`
	Model       string      `json:"model"`
	PricingMode PricingMode `json:"pricing_mode"` // fixed, dynamic, auction

	// 固定价格配置
	FixedPrice FixedPriceConfig `json:"fixed_price,omitempty"`

	// 动态价格配置
	DynamicPrice DynamicPriceConfig `json:"dynamic_price,omitempty"`

	// 竞价配置
	AuctionPrice AuctionPriceConfig `json:"auction_price,omitempty"`

	// 访问控制
	AccessRules []AccessRule `json:"access_rules"`

	// 状态
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PricingMode 定价模式
type PricingMode string

const (
	PricingModeFixed   PricingMode = "fixed"
	PricingModeDynamic PricingMode = "dynamic"
	PricingModeAuction PricingMode = "auction"
)

// FixedPriceConfig 固定价格配置
type FixedPriceConfig struct {
	InputTokenPrice  int64 `json:"input_token_price"`  // 每 1K input tokens 价格 (ustt)
	OutputTokenPrice int64 `json:"output_token_price"` // 每 1K output tokens 价格 (ustt)
	MinCharge        int64 `json:"min_charge"`         // 最小收费 (ustt)
	RequestFee       int64 `json:"request_fee"`        // 每次请求固定费用 (ustt)
}

// DynamicPriceConfig 动态价格配置
type DynamicPriceConfig struct {
	BaseInputPrice   int64   `json:"base_input_price"`  // 基础 input 价格
	BaseOutputPrice  int64   `json:"base_output_price"` // 基础 output 价格
	DemandMultiplier float64 `json:"demand_multiplier"` // 需求倍数
	SupplyMultiplier float64 `json:"supply_multiplier"` // 供应倍数
	MinMultiplier    float64 `json:"min_multiplier"`    // 最小倍数
	MaxMultiplier    float64 `json:"max_multiplier"`    // 最大倍数
	UpdateInterval   int64   `json:"update_interval"`   // 价格更新间隔（秒）
}

// AuctionPriceConfig 竞价价格配置
type AuctionPriceConfig struct {
	ReservePrice    int64     `json:"reserve_price"`     // 保留价
	AuctionDuration int64     `json:"auction_duration"`  // 竞价周期（秒）
	NextAuctionTime time.Time `json:"next_auction_time"` // 下次竞价时间
	CurrentWinner   string    `json:"current_winner"`    // 当前中标者
	CurrentPrice    int64     `json:"current_price"`     // 当前中标价
}

// UsageStats 使用统计
type UsageStats struct {
	APIKeyID          string    `json:"api_key_id"`
	TotalRequests     int64     `json:"total_requests"`
	TotalInputTokens  int64     `json:"total_input_tokens"`
	TotalOutputTokens int64     `json:"total_output_tokens"`
	TotalRevenue      int64     `json:"total_revenue"` // ustt
	LastUsedAt        time.Time `json:"last_used_at"`
}

// CalculateCost 计算调用成本
func (pc *PricingConfig) CalculateCost(inputTokens, outputTokens int64) int64 {
	switch pc.PricingMode {
	case PricingModeFixed:
		return pc.calculateFixedCost(inputTokens, outputTokens)
	case PricingModeDynamic:
		return pc.calculateDynamicCost(inputTokens, outputTokens)
	case PricingModeAuction:
		return pc.calculateAuctionCost(inputTokens, outputTokens)
	default:
		return pc.calculateFixedCost(inputTokens, outputTokens)
	}
}

// calculateFixedCost 计算固定价格成本
func (pc *PricingConfig) calculateFixedCost(inputTokens, outputTokens int64) int64 {
	fp := pc.FixedPrice

	// 计算 token 费用
	inputCost := (inputTokens * fp.InputTokenPrice) / 1000
	outputCost := (outputTokens * fp.OutputTokenPrice) / 1000

	total := inputCost + outputCost + fp.RequestFee

	// 应用最小收费
	if total < fp.MinCharge {
		total = fp.MinCharge
	}

	return total
}

// calculateDynamicCost 计算动态价格成本
func (pc *PricingConfig) calculateDynamicCost(inputTokens, outputTokens int64) int64 {
	dp := pc.DynamicPrice

	// 计算基础价格
	inputCost := (inputTokens * dp.BaseInputPrice) / 1000
	outputCost := (outputTokens * dp.BaseOutputPrice) / 1000
	baseTotal := float64(inputCost + outputCost)

	// 应用倍数
	multiplier := dp.DemandMultiplier * dp.SupplyMultiplier

	// 限制倍数范围
	if multiplier < dp.MinMultiplier {
		multiplier = dp.MinMultiplier
	}
	if multiplier > dp.MaxMultiplier {
		multiplier = dp.MaxMultiplier
	}

	return int64(baseTotal * multiplier)
}

// calculateAuctionCost 计算竞价价格成本
func (pc *PricingConfig) calculateAuctionCost(inputTokens, outputTokens int64) int64 {
	ap := pc.AuctionPrice

	// 使用当前中标价
	return ap.CurrentPrice
}

// UpdateDynamicMultiplier 更新动态价格倍数
func (pc *PricingConfig) UpdateDynamicMultiplier(demand, supply int64) {
	if pc.PricingMode != PricingModeDynamic {
		return
	}

	dp := &pc.DynamicPrice

	// 简单的供需模型
	if supply > 0 {
		ratio := float64(demand) / float64(supply)
		dp.DemandMultiplier = 1.0 + (ratio-1.0)*0.5
		dp.SupplyMultiplier = 1.0 / ratio
	}

	// 限制范围
	if dp.DemandMultiplier < dp.MinMultiplier {
		dp.DemandMultiplier = dp.MinMultiplier
	}
	if dp.DemandMultiplier > dp.MaxMultiplier {
		dp.DemandMultiplier = dp.MaxMultiplier
	}

	pc.UpdatedAt = time.Now()
}

// CanAccess 检查是否允许访问
func (pc *PricingConfig) CanAccess(serviceID string, usageCount int64) bool {
	for _, rule := range pc.AccessRules {
		if rule.ServiceId == serviceID {
			if !rule.Allowed {
				return false
			}
			if rule.MaxRequests > 0 && usageCount >= rule.MaxRequests {
				return false
			}
			return true
		}
	}
	// 默认允许
	return true
}

// GetServicePrice 获取服务特定价格
func (pc *PricingConfig) GetServicePrice(serviceID string) int64 {
	for _, rule := range pc.AccessRules {
		if rule.ServiceId == serviceID && rule.PricePerReq > 0 {
			return rule.PricePerReq
		}
	}
	// 返回默认价格
	return pc.FixedPrice.RequestFee
}

// ValidateBasic 基础验证
func (pc *PricingConfig) ValidateBasic() error {
	if pc.Owner == "" {
		return fmt.Errorf("owner is required")
	}
	if pc.APIKeyID == "" {
		return fmt.Errorf("api_key_id is required")
	}
	if !IsValidProvider(string(pc.Provider)) {
		return fmt.Errorf("invalid provider: %s", pc.Provider)
	}
	if pc.Model == "" {
		return fmt.Errorf("model is required")
	}

	switch pc.PricingMode {
	case PricingModeFixed:
		if pc.FixedPrice.InputTokenPrice == 0 && pc.FixedPrice.OutputTokenPrice == 0 {
			return fmt.Errorf("fixed price must be set")
		}
	case PricingModeDynamic:
		if pc.DynamicPrice.BaseInputPrice == 0 && pc.DynamicPrice.BaseOutputPrice == 0 {
			return fmt.Errorf("dynamic price base must be set")
		}
	case PricingModeAuction:
		if pc.AuctionPrice.ReservePrice == 0 {
			return fmt.Errorf("auction reserve price must be set")
		}
	default:
		return fmt.Errorf("invalid pricing mode: %s", pc.PricingMode)
	}

	return nil
}

// NewFixedPricing 创建固定定价配置
func NewFixedPricing(owner, apiKeyID string, provider Provider, model string, inputPrice, outputPrice int64) *PricingConfig {
	now := time.Now()
	return &PricingConfig{
		ID:          GeneratePricingID(),
		Owner:       owner,
		APIKeyID:    apiKeyID,
		Provider:    provider,
		Model:       model,
		PricingMode: PricingModeFixed,
		FixedPrice: FixedPriceConfig{
			InputTokenPrice:  inputPrice,
			OutputTokenPrice: outputPrice,
			MinCharge:        1000, // 0.001 STT
			RequestFee:       0,
		},
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewDynamicPricing 创建动态定价配置
func NewDynamicPricing(owner, apiKeyID string, provider Provider, model string, baseInputPrice, baseOutputPrice int64) *PricingConfig {
	now := time.Now()
	return &PricingConfig{
		ID:          GeneratePricingID(),
		Owner:       owner,
		APIKeyID:    apiKeyID,
		Provider:    provider,
		Model:       model,
		PricingMode: PricingModeDynamic,
		DynamicPrice: DynamicPriceConfig{
			BaseInputPrice:   baseInputPrice,
			BaseOutputPrice:  baseOutputPrice,
			DemandMultiplier: 1.0,
			SupplyMultiplier: 1.0,
			MinMultiplier:    0.5,
			MaxMultiplier:    5.0,
			UpdateInterval:   300, // 5 分钟
		},
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GeneratePricingID 生成定价配置 ID
func GeneratePricingID() string {
	return fmt.Sprintf("price-%d", time.Now().UnixNano())
}
