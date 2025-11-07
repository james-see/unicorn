package founder

import (
	"fmt"
	"math"
	"math/rand"
)

// InitializePricingStrategy sets up the default pricing model
func (fs *FounderState) InitializePricingStrategy() {
	fs.PricingStrategy = &PricingStrategy{
		Model: "tiered",
		CurrentTier: map[string]int64{
			"starter":    2000,  // $2k/month
			"pro":        5000,  // $5k/month
			"enterprise": 15000, // $15k/month
		},
		IsAnnual:      false,
		Discount:      0.0,
		ChangeHistory: []PricingChange{},
	}
}

// ChangePricingModel switches to a new pricing model
func (fs *FounderState) ChangePricingModel(newModel string, isAnnual bool, discount float64) error {
	validModels := map[string]bool{
		"freemium":        true,
		"trial":           true,
		"annual_upfront":  true,
		"usage_based":     true,
		"tiered":          true,
	}

	if !validModels[newModel] {
		return fmt.Errorf("invalid pricing model: %s", newModel)
	}

	if fs.PricingStrategy == nil {
		fs.InitializePricingStrategy()
	}

	oldModel := fs.PricingStrategy.Model
	fs.PricingStrategy.Model = newModel
	fs.PricingStrategy.IsAnnual = isAnnual
	fs.PricingStrategy.Discount = discount

	// Record the change
	change := PricingChange{
		Month:     fs.Turn,
		FromModel: oldModel,
		ToModel:   newModel,
		Reason:    "Manual change",
		Impact:    "Analyzing impact...",
	}
	fs.PricingStrategy.ChangeHistory = append(fs.PricingStrategy.ChangeHistory, change)

	// Apply immediate impacts based on model
	switch newModel {
	case "freemium":
		// Freemium: Much higher lead volume but lower conversion
		// Increases churn as well
		fs.CustomerChurnRate = math.Min(0.30, fs.CustomerChurnRate+0.03)
	case "annual_upfront":
		// Annual upfront: Lower churn, higher commitment
		fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-0.05)
	case "usage_based":
		// Usage-based: Variable, can lead to higher churn if not managed well
		fs.CustomerChurnRate = math.Min(0.30, fs.CustomerChurnRate+0.02)
	}

	return nil
}

// IncreasePrices raises prices by a percentage (causes churn risk)
func (fs *FounderState) IncreasePrices(percentage float64) error {
	if fs.PricingStrategy == nil {
		fs.InitializePricingStrategy()
	}

	if percentage < 0 || percentage > 0.50 {
		return fmt.Errorf("price increase must be between 0%% and 50%%")
	}

	// Increase all tier prices
	for tier, price := range fs.PricingStrategy.CurrentTier {
		newPrice := int64(float64(price) * (1.0 + percentage))
		fs.PricingStrategy.CurrentTier[tier] = newPrice
	}

	// Increase avg deal size
	fs.AvgDealSize = int64(float64(fs.AvgDealSize) * (1.0 + percentage))
	fs.MinDealSize = int64(float64(fs.MinDealSize) * (1.0 + percentage))
	fs.MaxDealSize = int64(float64(fs.MaxDealSize) * (1.0 + percentage))

	// Increase MRR from existing customers
	fs.MRR = int64(float64(fs.MRR) * (1.0 + percentage))
	fs.DirectMRR = int64(float64(fs.DirectMRR) * (1.0 + percentage))

	// Update individual customer deal sizes
	for i := range fs.CustomerList {
		if fs.CustomerList[i].IsActive {
			fs.CustomerList[i].DealSize = int64(float64(fs.CustomerList[i].DealSize) * (1.0 + percentage))
		}
	}

	// Cause churn risk (3% base + percentage * 5%)
	churnRisk := 0.03 + (percentage * 0.05)
	fs.CustomerChurnRate = math.Min(0.30, fs.CustomerChurnRate+churnRisk)

	// Record the change
	change := PricingChange{
		Month:     fs.Turn,
		FromModel: fs.PricingStrategy.Model,
		ToModel:   fs.PricingStrategy.Model,
		Reason:    fmt.Sprintf("Price increase +%.0f%%", percentage*100),
		Impact:    fmt.Sprintf("MRR +%.0f%%, Churn Risk +%.1f%%", percentage*100, churnRisk*100),
	}
	fs.PricingStrategy.ChangeHistory = append(fs.PricingStrategy.ChangeHistory, change)

	return nil
}

// DecreasePrices lowers prices by a percentage (improves acquisition)
func (fs *FounderState) DecreasePrices(percentage float64) error {
	if fs.PricingStrategy == nil {
		fs.InitializePricingStrategy()
	}

	if percentage < 0 || percentage > 0.50 {
		return fmt.Errorf("price decrease must be between 0%% and 50%%")
	}

	// Decrease all tier prices
	for tier, price := range fs.PricingStrategy.CurrentTier {
		newPrice := int64(float64(price) * (1.0 - percentage))
		if newPrice < 100 {
			newPrice = 100 // Minimum $100/month
		}
		fs.PricingStrategy.CurrentTier[tier] = newPrice
	}

	// Decrease avg deal size
	fs.AvgDealSize = int64(float64(fs.AvgDealSize) * (1.0 - percentage))
	if fs.AvgDealSize < 100 {
		fs.AvgDealSize = 100
	}
	fs.MinDealSize = int64(float64(fs.MinDealSize) * (1.0 - percentage))
	fs.MaxDealSize = int64(float64(fs.MaxDealSize) * (1.0 - percentage))

	// Decrease MRR from existing customers
	fs.MRR = int64(float64(fs.MRR) * (1.0 - percentage))
	fs.DirectMRR = int64(float64(fs.DirectMRR) * (1.0 - percentage))

	// Update individual customer deal sizes
	for i := range fs.CustomerList {
		if fs.CustomerList[i].IsActive {
			fs.CustomerList[i].DealSize = int64(float64(fs.CustomerList[i].DealSize) * (1.0 - percentage))
			if fs.CustomerList[i].DealSize < 100 {
				fs.CustomerList[i].DealSize = 100
			}
		}
	}

	// Reduce churn risk (percentage * 3%)
	churnReduction := percentage * 0.03
	fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-churnReduction)

	// Record the change
	change := PricingChange{
		Month:     fs.Turn,
		FromModel: fs.PricingStrategy.Model,
		ToModel:   fs.PricingStrategy.Model,
		Reason:    fmt.Sprintf("Price decrease -%.0f%%", percentage*100),
		Impact:    fmt.Sprintf("MRR -%.0f%%, Churn Risk -%.1f%%", percentage*100, churnReduction*100),
	}
	fs.PricingStrategy.ChangeHistory = append(fs.PricingStrategy.ChangeHistory, change)

	return nil
}

// StartPricingExperiment launches an A/B test
func (fs *FounderState) StartPricingExperiment(name string, testModel string, testAnnual bool, testDiscount float64, cost int64, duration int) error {
	if fs.ActiveExperiment != nil && !fs.ActiveExperiment.IsComplete {
		return fmt.Errorf("already running an experiment: %s", fs.ActiveExperiment.Name)
	}

	if fs.Cash < cost {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	// Pay for the experiment
	fs.Cash -= cost

	// Create test strategy
	testStrategy := PricingStrategy{
		Model:         testModel,
		CurrentTier:   make(map[string]int64),
		IsAnnual:      testAnnual,
		Discount:      testDiscount,
		ChangeHistory: []PricingChange{},
	}

	// Copy current tier prices with potential modifications
	for tier, price := range fs.PricingStrategy.CurrentTier {
		testStrategy.CurrentTier[tier] = price
	}

	fs.ActiveExperiment = &PricingExperiment{
		Name:         name,
		Cost:         cost,
		StartMonth:   fs.Turn,
		Duration:     duration,
		TestStrategy: testStrategy,
		Results:      PricingResults{},
		IsComplete:   false,
	}

	return nil
}

// ProcessPricingExperiment updates running experiments
func (fs *FounderState) ProcessPricingExperiment() []string {
	var messages []string

	if fs.ActiveExperiment == nil || fs.ActiveExperiment.IsComplete {
		return messages
	}

	monthsRunning := fs.Turn - fs.ActiveExperiment.StartMonth

	if monthsRunning >= fs.ActiveExperiment.Duration {
		// Experiment complete - generate results
		fs.ActiveExperiment.IsComplete = true

		// Generate results based on the test model
		results := PricingResults{
			ConversionRateChange: (rand.Float64()*0.20 - 0.10), // -10% to +10%
			AvgDealSizeChange:    int64(rand.Float64()*4000 - 2000), // -$2k to +$2k
			ChurnRateChange:      (rand.Float64()*0.06 - 0.03), // -3% to +3%
			Confidence:           0.70 + rand.Float64()*0.25, // 70-95% confidence
		}

		// Adjust based on model characteristics
		switch fs.ActiveExperiment.TestStrategy.Model {
		case "freemium":
			results.ConversionRateChange = -0.90 // 90% lower conversion (but 10x more leads)
			results.AvgDealSizeChange = -500
			results.ChurnRateChange = 0.05
		case "annual_upfront":
			results.ConversionRateChange = -0.15 // 15% lower conversion (commitment)
			results.AvgDealSizeChange = 3000
			results.ChurnRateChange = -0.05 // Lower churn
		case "usage_based":
			results.ConversionRateChange = 0.10 // 10% higher conversion (no commitment)
			results.AvgDealSizeChange = -1000 // Lower initial value
			results.ChurnRateChange = 0.03
		}

		fs.ActiveExperiment.Results = results

		messages = append(messages, fmt.Sprintf("🧪 Pricing experiment '%s' complete!", fs.ActiveExperiment.Name))
		messages = append(messages, fmt.Sprintf("   Conversion Rate: %+.1f%% | Deal Size: %+$%s | Churn: %+.1f%%",
			results.ConversionRateChange*100,
			formatCurrency(int64(math.Abs(float64(results.AvgDealSizeChange)))),
			results.ChurnRateChange*100))
		messages = append(messages, fmt.Sprintf("   Confidence: %.0f%%", results.Confidence*100))
	} else {
		// Still running
		messages = append(messages, fmt.Sprintf("🧪 Pricing experiment running: %d/%d months complete",
			monthsRunning, fs.ActiveExperiment.Duration))
	}

	return messages
}

// ApplyExperimentResults applies the winning experiment strategy
func (fs *FounderState) ApplyExperimentResults() error {
	if fs.ActiveExperiment == nil || !fs.ActiveExperiment.IsComplete {
		return fmt.Errorf("no completed experiment to apply")
	}

	// Apply the test strategy
	oldModel := fs.PricingStrategy.Model
	fs.PricingStrategy = &fs.ActiveExperiment.TestStrategy

	// Apply the results to actual metrics
	results := fs.ActiveExperiment.Results
	
	// Update avg deal size
	fs.AvgDealSize += results.AvgDealSizeChange
	if fs.AvgDealSize < 100 {
		fs.AvgDealSize = 100
	}

	// Update churn rate
	fs.CustomerChurnRate += results.ChurnRateChange
	if fs.CustomerChurnRate < 0.01 {
		fs.CustomerChurnRate = 0.01
	}
	if fs.CustomerChurnRate > 0.30 {
		fs.CustomerChurnRate = 0.30
	}

	// Record the change
	change := PricingChange{
		Month:     fs.Turn,
		FromModel: oldModel,
		ToModel:   fs.PricingStrategy.Model,
		Reason:    fmt.Sprintf("Applied experiment: %s", fs.ActiveExperiment.Name),
		Impact:    fmt.Sprintf("Conversion %+.1f%%, Deal Size %+$%s, Churn %+.1f%%",
			results.ConversionRateChange*100,
			formatCurrency(int64(math.Abs(float64(results.AvgDealSizeChange)))),
			results.ChurnRateChange*100),
	}
	fs.PricingStrategy.ChangeHistory = append(fs.PricingStrategy.ChangeHistory, change)

	// Clear active experiment
	fs.ActiveExperiment = nil

	return nil
}

// CheckCompetitorPricing simulates competitor pricing pressure
func (fs *FounderState) CheckCompetitorPricing() []string {
	var messages []string

	if fs.PricingStrategy == nil {
		return messages
	}

	// Get average competitor pricing (simulate)
	var avgCompetitorPrice int64 = fs.AvgDealSize // Default to same as ours

	// Check if we have active competitors
	if len(fs.Competitors) > 0 {
		activeCompetitors := 0
		for _, comp := range fs.Competitors {
			if comp.Active {
				activeCompetitors++
			}
		}

		if activeCompetitors > 0 {
			// Competitors typically price ±20% of market average
			variance := 0.20
			avgCompetitorPrice = int64(float64(fs.AvgDealSize) * (1.0 + (rand.Float64()*variance*2 - variance)))
		}
	}

	// Check if we're significantly more expensive (>30% over market)
	priceDiff := float64(fs.AvgDealSize-avgCompetitorPrice) / float64(avgCompetitorPrice)
	
	if priceDiff > 0.30 {
		messages = append(messages, fmt.Sprintf("⚠️  Your pricing is %.0f%% above market average - losing deals to competitors!", priceDiff*100))
		// Reduce growth slightly
		fs.MonthlyGrowthRate = math.Max(0.01, fs.MonthlyGrowthRate*0.95)
	} else if priceDiff < -0.30 {
		messages = append(messages, fmt.Sprintf("💡 Your pricing is %.0f%% below market - consider increasing prices!", math.Abs(priceDiff)*100))
	}

	return messages
}

// GetPricingModelDescription returns a description of the pricing model
func GetPricingModelDescription(model string) string {
	descriptions := map[string]string{
		"freemium":       "Free tier with paid upgrades. 10x leads, 2% conversion, higher churn",
		"trial":          "Free trial then paid. Standard conversion, moderate churn",
		"annual_upfront": "Annual contracts paid upfront. Lower conversion, lower churn, higher deal size",
		"usage_based":    "Pay for what you use. Higher conversion, variable revenue, moderate churn",
		"tiered":         "Multiple pricing tiers. Balanced approach, standard metrics",
	}
	return descriptions[model]
}

