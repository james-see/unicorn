package founder

import (
	"fmt"
	"math/rand"
)

// InitializePlatform initializes the platform system
func InitializePlatform(fs *FounderState) {
	if fs.PlatformMetrics == nil {
		fs.PlatformMetrics = &PlatformMetrics{
			IsPlatform:         false,
			ThirdPartyApps:     0,
			DeveloperCount:      0,
			APIUsage:           0,
			MarketplaceRevenue: 0,
			NetworkEffectScore: 0.0,
			PlatformType:       "",
			LaunchedMonth:      0,
		}
	}
	if fs.NetworkEffects == nil {
		fs.NetworkEffects = []NetworkEffect{}
	}
}

// CanLaunchPlatform checks if founder can launch platform
func (fs *FounderState) CanLaunchPlatform() bool {
	// Unlock: $1M+ ARR OR 500+ customers
	arr := fs.MRR * 12
	return arr >= 1000000 || fs.Customers >= 500
}

// LaunchPlatform launches a platform business model
func (fs *FounderState) LaunchPlatform(platformType string) error {
	if !fs.CanLaunchPlatform() {
		return fmt.Errorf("platform requires $1M+ ARR or 500+ customers")
	}

	if fs.PlatformMetrics.IsPlatform {
		return fmt.Errorf("platform already launched")
	}

	validTypes := map[string]bool{
		"marketplace":     true,
		"social":         true,
		"data":           true,
		"infrastructure": true,
	}
	if !validTypes[platformType] {
		return fmt.Errorf("invalid platform type: %s", platformType)
	}

	// Setup cost: $100-300k
	setupCost := int64(100000 + rand.Int63n(200000))
	if setupCost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(setupCost))
	}

	fs.Cash -= setupCost

	fs.PlatformMetrics.IsPlatform = true
	fs.PlatformMetrics.PlatformType = platformType
	fs.PlatformMetrics.LaunchedMonth = fs.Turn

	// Initialize network effect based on platform type
	networkEffect := NetworkEffect{
		Type:          platformType,
		Strength:      0.0,
		CustomerValue: 0.0,
		Threshold:     0,
		Active:        false,
		ActivatedMonth: 0,
	}

	switch platformType {
	case "marketplace":
		networkEffect.Threshold = 500 // Need 500 customers for marketplace network effects
		networkEffect.CustomerValue = 0.001 // Each customer adds 0.1% value
	case "social":
		networkEffect.Threshold = 1000 // Need 1000 customers for social network effects
		networkEffect.CustomerValue = 0.002 // Each customer adds 0.2% value
	case "data":
		networkEffect.Threshold = 300 // Need 300 customers for data network effects
		networkEffect.CustomerValue = 0.0015 // Each customer adds 0.15% value
	case "infrastructure":
		networkEffect.Threshold = 200 // Need 200 customers for infrastructure network effects
		networkEffect.CustomerValue = 0.0005 // Each customer adds 0.05% value
	}

	fs.NetworkEffects = append(fs.NetworkEffects, networkEffect)

	return nil
}

// ProcessPlatformMetrics processes platform growth and network effects
func (fs *FounderState) ProcessPlatformMetrics() []string {
	var messages []string

	if !fs.PlatformMetrics.IsPlatform {
		return messages
	}

	// Developer growth: depends on API usage and marketing
	baseDeveloperGrowth := 2 // 2 developers per month base
	if fs.PlatformMetrics.APIUsage > 1000000 {
		baseDeveloperGrowth += 5 // More API usage = more developers
	}
	if len(fs.Team.Marketing) > 0 {
		baseDeveloperGrowth += len(fs.Team.Marketing) * 2
	}

	fs.PlatformMetrics.DeveloperCount += baseDeveloperGrowth

	// API usage growth: depends on developers and customers
	apiGrowthRate := 0.10 // 10% growth per month
	if fs.PlatformMetrics.DeveloperCount > 0 {
		apiGrowthRate += float64(fs.PlatformMetrics.DeveloperCount) * 0.01 // +1% per developer
	}
	if fs.Customers > 0 {
		apiGrowthRate += float64(fs.Customers) * 0.0001 // +0.01% per customer
	}

	fs.PlatformMetrics.APIUsage = int64(float64(fs.PlatformMetrics.APIUsage) * (1.0 + apiGrowthRate))
	if fs.PlatformMetrics.APIUsage < 10000 {
		fs.PlatformMetrics.APIUsage = 10000 // Minimum baseline
	}

	// Third-party apps: grows with developers
	if fs.PlatformMetrics.DeveloperCount > 50 && rand.Float64() < 0.1 {
		fs.PlatformMetrics.ThirdPartyApps++
		messages = append(messages, fmt.Sprintf("📱 New third-party app launched on platform (Total: %d)", fs.PlatformMetrics.ThirdPartyApps))
	}

	// Marketplace revenue: % of transactions
	if fs.PlatformMetrics.PlatformType == "marketplace" {
		// Marketplace takes 5-15% of transaction value
		commissionRate := 0.05 + rand.Float64()*0.10
		// Estimate transaction volume as 2x MRR
		transactionVolume := fs.MRR * 2
		fs.PlatformMetrics.MarketplaceRevenue = int64(float64(transactionVolume) * commissionRate)
	}

	// Process network effects
	for i := range fs.NetworkEffects {
		ne := &fs.NetworkEffects[i]

		if !ne.Active && fs.Customers >= ne.Threshold {
			// Activate network effect
			ne.Active = true
			ne.ActivatedMonth = fs.Turn
			ne.Strength = 0.3 // Start at 30% strength
			messages = append(messages, fmt.Sprintf("🌐 Network effect activated! (%s platform reached %d customers)", ne.Type, ne.Threshold))
		}

		if ne.Active {
			// Network effect strength grows with customers
			customerGrowth := float64(fs.Customers) * ne.CustomerValue
			ne.Strength = customerGrowth
			if ne.Strength > 1.0 {
				ne.Strength = 1.0 // Cap at 100%
			}

			// Update platform network effect score
			fs.PlatformMetrics.NetworkEffectScore = ne.Strength
		}
	}

	return messages
}

// ApplyNetworkEffectBonuses applies network effect benefits to metrics
func (fs *FounderState) ApplyNetworkEffectBonuses() (cacReduction float64, retentionBonus float64, growthBonus float64) {
	if !fs.PlatformMetrics.IsPlatform {
		return 1.0, 0.0, 1.0
	}

	cacReduction = 1.0
	retentionBonus = 0.0
	growthBonus = 1.0

	// Find active network effect
	for _, ne := range fs.NetworkEffects {
		if ne.Active {
			// CAC reduction: customers refer others (5-15% reduction)
			cacReduction = 1.0 - (ne.Strength * 0.15)

			// Retention bonus: switching costs increase (2-8% churn reduction)
			retentionBonus = ne.Strength * 0.08

			// Growth bonus: network effects accelerate growth (5-20% boost)
			growthBonus = 1.0 + (ne.Strength * 0.20)
			break
		}
	}

	return cacReduction, retentionBonus, growthBonus
}

// InvestInDeveloperProgram invests in developer relations
func (fs *FounderState) InvestInDeveloperProgram(monthlyBudget int64) error {
	if monthlyBudget > fs.Cash {
		return fmt.Errorf("insufficient cash")
	}

	if !fs.PlatformMetrics.IsPlatform {
		return fmt.Errorf("platform not launched")
	}

	fs.Cash -= monthlyBudget
	fs.PlatformMetrics.DeveloperCount += int(monthlyBudget / 5000) // $5k per developer attracted

	return nil
}

// GetPlatformSummary returns platform metrics summary
func (fs *FounderState) GetPlatformSummary() (isPlatform bool, developers int, apps int, apiUsage int64, marketplaceRevenue int64, networkScore float64) {
	if fs.PlatformMetrics == nil {
		return false, 0, 0, 0, 0, 0.0
	}
	return fs.PlatformMetrics.IsPlatform,
		fs.PlatformMetrics.DeveloperCount,
		fs.PlatformMetrics.ThirdPartyApps,
		fs.PlatformMetrics.APIUsage,
		fs.PlatformMetrics.MarketplaceRevenue,
		fs.PlatformMetrics.NetworkEffectScore
}

