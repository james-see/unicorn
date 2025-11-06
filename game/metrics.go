package game

import (
	"math/rand"
)


func (gs *GameState) UpdateCompanyFinancials(startup *Startup) {
	// Apply growth rate to revenue (with some randomness)
	growthVariance := (rand.Float64()*0.4 - 0.2) // -20% to +20% variance
	actualGrowth := startup.RevenueGrowthRate + growthVariance

	// Update revenue based on growth
	startup.MonthlyRevenue = int64(float64(startup.MonthlyRevenue) * (1 + actualGrowth))

	// Track revenue history (keep last 6 months for trends)
	startup.RevenueHistory = append(startup.RevenueHistory, startup.MonthlyRevenue)
	if len(startup.RevenueHistory) > 6 {
		startup.RevenueHistory = startup.RevenueHistory[len(startup.RevenueHistory)-6:]
	}

	// Costs grow slower than revenue (economies of scale)
	costGrowth := actualGrowth * 0.6 // Costs grow at 60% of revenue growth rate
	startup.MonthlyCosts = int64(float64(startup.MonthlyCosts) * (1 + costGrowth))

	// Calculate net income
	startup.NetIncome = startup.MonthlyRevenue - startup.MonthlyCosts

	// Update cumulative totals
	startup.CumulativeRevenue += startup.MonthlyRevenue
	startup.CumulativeCosts += startup.MonthlyCosts

	// Update customer count based on revenue
	if startup.SalePrice > 0 {
		startup.CustomerCount = int(startup.MonthlyRevenue / int64(startup.SalePrice))
	}

	// Update MRR
	startup.MonthlyRecurringRevenue = startup.MonthlyRevenue

	// Adjust growth rate based on performance
	if startup.NetIncome > 0 {
		startup.RevenueGrowthRate *= 1.02 // Profitable companies grow faster
	} else {
		startup.RevenueGrowthRate *= 0.98 // Unprofitable slow down
	}

	// Cap growth rate
	if startup.RevenueGrowthRate > 0.30 {
		startup.RevenueGrowthRate = 0.30 // Max 30% monthly growth
	}
	if startup.RevenueGrowthRate < -0.15 {
		startup.RevenueGrowthRate = -0.15 // Max 15% monthly decline
	}

	// Update valuation based on financial performance
	annualRevenue := startup.MonthlyRevenue * 12

	// Revenue multiple varies by profitability
	revenueMultiple := 10.0
	if startup.NetIncome > 0 {
		revenueMultiple = 15.0 // Profitable companies get premium
	}

	newValuation := int64(float64(annualRevenue) * revenueMultiple)

	// Smooth valuation changes (max 20% per month)
	maxChange := float64(startup.Valuation) * 0.20
	valuationChange := newValuation - startup.Valuation
	if valuationChange > int64(maxChange) {
		newValuation = startup.Valuation + int64(maxChange)
	} else if valuationChange < -int64(maxChange) {
		newValuation = startup.Valuation - int64(maxChange)
	}

	// Minimum valuation
	if newValuation < 100000 {
		newValuation = 100000
	}

	startup.Valuation = newValuation
}

func (gs *GameState) Calculate409AValuation(startup *Startup) int64 {
	// 409A considers multiple factors
	annualRevenue := startup.MonthlyRevenue * 12

	// Revenue multiple (conservative for 409A)
	revenueMultiple := 8.0
	if startup.NetIncome > 0 {
		revenueMultiple = 12.0
	}
	revenueValue := int64(float64(annualRevenue) * revenueMultiple)

	// Cost to duplicate
	costValue := startup.CumulativeCosts

	// Market value
	marketValue := startup.Valuation

	// Weighted average
	val409A := (revenueValue*4 + costValue*2 + marketValue*4) / 10

	// 409A is typically 20-30% discount to FMV
	val409A = int64(float64(val409A) * 0.75)

	startup.Last409AValuation = val409A
	startup.Last409AMonth = gs.Portfolio.Turn

	return val409A
}

func (gs *GameState) GetLeaderboard() []PlayerScore {
	scores := []PlayerScore{}

	// Add player - ROI based on total starting capital (cash + follow-on reserve)
	totalStartingCapital := gs.Portfolio.InitialFundSize + gs.Portfolio.FollowOnReserve
	playerROI := ((float64(gs.Portfolio.NetWorth) - float64(totalStartingCapital)) / float64(totalStartingCapital)) * 100.0
	scores = append(scores, PlayerScore{
		Name:     gs.PlayerName,
		Firm:     gs.PlayerFirmName,
		NetWorth: gs.Portfolio.NetWorth,
		ROI:      playerROI,
		IsPlayer: true,
	})

	// Add AI players - same calculation
	for _, ai := range gs.AIPlayers {
		aiTotalCapital := ai.Portfolio.InitialFundSize + ai.Portfolio.FollowOnReserve
		aiROI := ((float64(ai.Portfolio.NetWorth) - float64(aiTotalCapital)) / float64(aiTotalCapital)) * 100.0
		scores = append(scores, PlayerScore{
			Name:     ai.Name,
			Firm:     ai.Firm,
			NetWorth: ai.Portfolio.NetWorth,
			ROI:      aiROI,
			IsPlayer: false,
		})
	}

	// Sort by net worth
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].NetWorth > scores[i].NetWorth {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	return scores
}

func (gs *GameState) GetSectorTrends() map[string]string {
	sectorValuations := make(map[string][]int64)

	// Group startups by sector
	for _, startup := range gs.AvailableStartups {
		sectorValuations[startup.Category] = append(sectorValuations[startup.Category], startup.Valuation)
	}

	trends := make(map[string]string)

	// Calculate average valuation per sector
	sectorAverages := make(map[string]float64)
	for sector, valuations := range sectorValuations {
		if len(valuations) > 0 {
			sum := int64(0)
			for _, val := range valuations {
				sum += val
			}
			sectorAverages[sector] = float64(sum) / float64(len(valuations))
		}
	}

	// Find overall average
	overallSum := 0.0
	count := 0
	for _, avg := range sectorAverages {
		overallSum += avg
		count++
	}
	overallAvg := overallSum / float64(count)

	// Categorize sectors
	for sector, avg := range sectorAverages {
		if avg > overallAvg*1.15 {
			trends[sector] = "📈 Hot"
		} else if avg > overallAvg*1.05 {
			trends[sector] = "📊 Strong"
		} else if avg < overallAvg*0.85 {
			trends[sector] = "📉 Cold"
		} else if avg < overallAvg*0.95 {
			trends[sector] = "📊 Weak"
		} else {
			trends[sector] = "➡️ Stable"
		}
	}

	return trends
}