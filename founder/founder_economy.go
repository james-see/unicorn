package founder

import (
	"fmt"
	"math/rand"
)

// SpawnEconomicEvent generates an economic downturn
func (fs *FounderState) SpawnEconomicEvent() *EconomicEvent {
	// Unlock: Month 12+ (after initial growth phase)
	if fs.Turn < 12 {
		return nil
	}

	// Probability: 5% per month after month 12
	if rand.Float64() > 0.05 {
		return nil
	}

	// Event types
	eventTypes := []string{"recession", "market_crash", "funding_winter", "sector_crash"}
	eventType := eventTypes[rand.Intn(len(eventTypes))]

	// Severity
	severityRoll := rand.Float64()
	severity := "mild"
	if severityRoll < 0.15 {
		severity = "extreme"
	} else if severityRoll < 0.35 {
		severity = "severe"
	} else if severityRoll < 0.65 {
		severity = "moderate"
	}

	// Impact based on severity
	var growthImpact, cacImpact, churnImpact, fundingImpact, customerBudgetCut float64
	var durationMonths int

	switch severity {
	case "extreme":
		growthImpact = 0.3 + rand.Float64()*0.2  // 30-50% growth reduction
		cacImpact = 1.4 + rand.Float64()*0.3    // 1.4-1.7x CAC
		churnImpact = 0.10 + rand.Float64()*0.05 // 10-15% churn
		fundingImpact = 0.2                       // 20% funding availability
		customerBudgetCut = 0.4 + rand.Float64()*0.2 // 40-60% budget cuts
		durationMonths = 18 + rand.Intn(6)       // 18-24 months
	case "severe":
		growthImpact = 0.4 + rand.Float64()*0.2  // 40-60% growth reduction
		cacImpact = 1.3 + rand.Float64()*0.2    // 1.3-1.5x CAC
		churnImpact = 0.07 + rand.Float64()*0.03 // 7-10% churn
		fundingImpact = 0.3                       // 30% funding availability
		customerBudgetCut = 0.3 + rand.Float64()*0.2 // 30-50% budget cuts
		durationMonths = 12 + rand.Intn(6)       // 12-18 months
	case "moderate":
		growthImpact = 0.5 + rand.Float64()*0.2  // 50-70% growth reduction
		cacImpact = 1.2 + rand.Float64()*0.2    // 1.2-1.4x CAC
		churnImpact = 0.04 + rand.Float64()*0.03 // 4-7% churn
		fundingImpact = 0.5                       // 50% funding availability
		customerBudgetCut = 0.2 + rand.Float64()*0.2 // 20-40% budget cuts
		durationMonths = 6 + rand.Intn(6)        // 6-12 months
	case "mild":
		growthImpact = 0.7 + rand.Float64()*0.2  // 70-90% growth reduction
		cacImpact = 1.1 + rand.Float64()*0.1    // 1.1-1.2x CAC
		churnImpact = 0.02 + rand.Float64()*0.02 // 2-4% churn
		fundingImpact = 0.7                       // 70% funding availability
		customerBudgetCut = 0.1 + rand.Float64()*0.1 // 10-20% budget cuts
		durationMonths = 3 + rand.Intn(4)        // 3-6 months
	}

	event := EconomicEvent{
		Type:             eventType,
		Severity:         severity,
		Month:            fs.Turn,
		DurationMonths:   durationMonths,
		GrowthImpact:     growthImpact,
		CACImpact:        cacImpact,
		ChurnImpact:      churnImpact,
		FundingImpact:    fundingImpact,
		CustomerBudgetCut: customerBudgetCut,
		Active:           true,
	}

	fs.EconomicEvent = &event
	return &event
}

// ExecuteSurvivalStrategy executes a survival strategy
func (fs *FounderState) ExecuteSurvivalStrategy(strategy string) error {
	validStrategies := map[string]bool{
		"cut_costs":    true,
		"pivot":        true,
		"downround":    true,
		"extend_runway": true,
		"acquire":      true,
	}
	if !validStrategies[strategy] {
		return fmt.Errorf("invalid strategy: %s", strategy)
	}

	var cost int64
	var effectiveness float64
	var tradeoffs []string

	switch strategy {
	case "cut_costs":
		// Layoffs: reduce team by 20-40%
		layoffPercent := 0.20 + rand.Float64()*0.20
		engineersToLayoff := int(float64(len(fs.Team.Engineers)) * layoffPercent)
		salesToLayoff := int(float64(len(fs.Team.Sales)) * layoffPercent)
		csToLayoff := int(float64(len(fs.Team.CustomerSuccess)) * layoffPercent)
		marketingToLayoff := int(float64(len(fs.Team.Marketing)) * layoffPercent)

		// Remove employees
		if engineersToLayoff > 0 && len(fs.Team.Engineers) > engineersToLayoff {
			fs.Team.Engineers = fs.Team.Engineers[:len(fs.Team.Engineers)-engineersToLayoff]
		}
		if salesToLayoff > 0 && len(fs.Team.Sales) > salesToLayoff {
			fs.Team.Sales = fs.Team.Sales[:len(fs.Team.Sales)-salesToLayoff]
		}
		if csToLayoff > 0 && len(fs.Team.CustomerSuccess) > csToLayoff {
			fs.Team.CustomerSuccess = fs.Team.CustomerSuccess[:len(fs.Team.CustomerSuccess)-csToLayoff]
		}
		if marketingToLayoff > 0 && len(fs.Team.Marketing) > marketingToLayoff {
			fs.Team.Marketing = fs.Team.Marketing[:len(fs.Team.Marketing)-marketingToLayoff]
		}

		fs.CalculateTeamCost()
		cost = 0 // No cash cost, but severance is implicit
		effectiveness = 0.6
		tradeoffs = []string{"Reduced team productivity", "Lower growth capacity", "Morale impact"}

	case "pivot":
		cost = 100000 + rand.Int63n(200000) // $100-300k
		if cost > fs.Cash {
			return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}
		fs.Cash -= cost
		effectiveness = 0.5
		tradeoffs = []string{"Lose 20-50% of customers", "Product development reset", "Market confusion"}

	case "downround":
		// Raise at lower valuation (handled in funding system)
		cost = 0
		effectiveness = 0.7
		tradeoffs = []string{"Higher dilution", "Investor pressure", "Valuation reset"}

	case "extend_runway":
		// Reduce growth to extend runway
		fs.MonthlyGrowthRate *= 0.5 // Cut growth in half
		cost = 0
		effectiveness = 0.4
		tradeoffs = []string{"Slower growth", "Market share loss", "Competitor advantage"}

	case "acquire":
		// Acquire struggling competitors cheap
		cost = 200000 + rand.Int63n(300000) // $200-500k (cheap during downturn)
		if cost > fs.Cash {
			return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}
		fs.Cash -= cost
		// Gain customers and MRR (simplified)
		customersGained := 10 + rand.Intn(20)
		mrrGained := int64(customersGained) * fs.AvgDealSize
		fs.Customers += customersGained
		fs.DirectCustomers += customersGained
		fs.DirectMRR += mrrGained
		fs.syncMRR()
		effectiveness = 0.8
		tradeoffs = []string{"Integration complexity", "Cultural challenges"}
	}

	strategyRecord := SurvivalStrategy{
		Strategy:     strategy,
		Cost:         cost,
		Effectiveness: effectiveness,
		Tradeoffs:    tradeoffs,
		MonthStarted: fs.Turn,
		Active:       true,
	}

	fs.SurvivalStrategies = append(fs.SurvivalStrategies, strategyRecord)

	// Apply effectiveness to economic event
	if fs.EconomicEvent != nil && fs.EconomicEvent.Active {
		fs.EconomicEvent.GrowthImpact = 1.0 - ((1.0 - fs.EconomicEvent.GrowthImpact) * (1.0 - effectiveness))
		fs.EconomicEvent.CACImpact = 1.0 + ((fs.EconomicEvent.CACImpact - 1.0) * (1.0 - effectiveness))
		fs.EconomicEvent.ChurnImpact *= (1.0 - effectiveness)
	}

	return nil
}

// ProcessEconomicEvent processes active economic event
func (fs *FounderState) ProcessEconomicEvent() []string {
	var messages []string

	if fs.EconomicEvent == nil || !fs.EconomicEvent.Active {
		return messages
	}

	event := fs.EconomicEvent

	// Apply growth impact
	fs.MonthlyGrowthRate *= event.GrowthImpact

	// Apply CAC impact
	fs.BaseCAC = int64(float64(fs.BaseCAC) * event.CACImpact)

	// Apply churn impact
	fs.CustomerChurnRate += event.ChurnImpact

	// Check if duration expired
	monthsActive := fs.Turn - event.Month
	if monthsActive >= event.DurationMonths {
		event.Active = false
		messages = append(messages, fmt.Sprintf("📈 Economic event ended: %s", event.Type))
	}

	return messages
}

