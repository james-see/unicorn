package founder

import (
	"math"
	"math/rand"
)

func (fs *FounderState) CalculateTeamCost() {
	total := fs.FounderSalary // Start with founder salary
	count := 1                // Founder counts as 1 employee

	for _, e := range fs.Team.Engineers {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.Sales {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.CustomerSuccess {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.Marketing {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.Executives {
		total += e.MonthlyCost
		count++
	}

	fs.Team.TotalMonthlyCost = total
	fs.Team.TotalEmployees = count
	
	// Apply Lower Burn upgrade (-10% team costs)
	hasLowerBurn := false
	for _, upgradeID := range fs.PlayerUpgrades {
		if upgradeID == "lower_burn" {
			hasLowerBurn = true
			break
		}
	}
	if hasLowerBurn {
		total = int64(float64(total) * 0.9)
		fs.Team.TotalMonthlyCost = total
	}
	
	fs.MonthlyTeamCost = total
}


func (fs *FounderState) CalculateRunway() {
	monthlyBurn := fs.Team.TotalMonthlyCost + 20000 // Team + $20k ops costs
	monthlyRevenue := fs.MRR
	netBurn := monthlyBurn - monthlyRevenue

	if netBurn <= 0 {
		// Cash positive! Runway is infinite
		fs.CashRunwayMonths = -1
		// Track when profitability was first reached
		if fs.MonthReachedProfitability == -1 {
			fs.MonthReachedProfitability = fs.Turn
		}
	} else {
		fs.CashRunwayMonths = int(fs.Cash / netBurn)
	}
}


func (fs *FounderState) CalculateInfrastructureCosts() {
	// Calculate compute costs per customer (variable, random, but never more than deal size)
	// For SaaS, compute costs are typically 10-30% of deal size
	// For DeepTech/Hardware, compute costs might be higher (20-40%)
	// For GovTech, compute costs might be lower (5-15%)

	// Compute cost percentages vary by category
	var computePercent float64
	switch fs.Category {
	case "SaaS":
		computePercent = 0.10 + rand.Float64()*0.20 // 10-30% of deal size
	case "DeepTech":
		computePercent = 0.20 + rand.Float64()*0.20 // 20-40% of deal size
	case "GovTech":
		computePercent = 0.05 + rand.Float64()*0.10 // 5-15% of deal size
	case "Hardware":
		computePercent = 0.15 + rand.Float64()*0.25 // 15-40% of deal size
	default:
		computePercent = 0.10 + rand.Float64()*0.20 // 10-30% default
	}

	// ODC costs are typically 5-15% of deal size (support, data transfer, etc.)
	odcPercent := 0.05 + rand.Float64()*0.10

	// Calculate costs based on each active customer's deal size
	var totalComputeCost int64
	var totalODCCost int64

	for _, c := range fs.CustomerList {
		if !c.IsActive {
			continue
		}

		dealSize := c.DealSize
		if dealSize == 0 {
			continue
		}

		// Compute cost for this customer
		customerComputeCost := int64(float64(dealSize) * computePercent)
		customerODCCost := int64(float64(dealSize) * odcPercent)

		// Ensure total infrastructure cost never exceeds 80% of deal size
		maxCost := int64(float64(dealSize) * 0.80)
		if customerComputeCost+customerODCCost > maxCost {
			// Scale down proportionally
			totalCost := float64(customerComputeCost + customerODCCost)
			scale := float64(maxCost) / totalCost
			customerComputeCost = int64(float64(customerComputeCost) * scale)
			customerODCCost = int64(float64(customerODCCost) * scale)
		}

		totalComputeCost += customerComputeCost
		totalODCCost += customerODCCost
	}

	fs.MonthlyComputeCost = totalComputeCost
	fs.MonthlyODCCost = totalODCCost
	
	// Apply Cloud Free First Year upgrade (no compute costs for first 12 months)
	hasCloudFree := false
	for _, upgradeID := range fs.PlayerUpgrades {
		if upgradeID == "cloud_free_first_year" {
			hasCloudFree = true
			break
		}
	}
	if hasCloudFree && fs.Turn <= 12 {
		fs.MonthlyComputeCost = 0 // Free compute for first 12 months
	}
}


func (fs *FounderState) CalculateLTVToCAC() float64 {
	if fs.CustomerAcquisitionCost == 0 {
		return 0
	}

	// LTV = Average Revenue per Customer / Churn Rate
	avgRevenuePerCustomer := float64(fs.AvgDealSize)
	ltv := avgRevenuePerCustomer / math.Max(0.01, fs.CustomerChurnRate)

	return ltv / float64(fs.CustomerAcquisitionCost)
}


func (fs *FounderState) CalculateCACPayback() float64 {
	if fs.AvgDealSize == 0 {
		return 0
	}

	// Payback period = CAC / Monthly Revenue per Customer
	return float64(fs.CustomerAcquisitionCost) / float64(fs.AvgDealSize)
}


func (fs *FounderState) CalculateRuleOf40() float64 {
	// Growth rate (as %)
	growthRate := fs.MonthlyGrowthRate * 100

	// Profit margin = (MRR - Costs) / MRR * 100
	annualizedMRR := fs.MRR * 12
	annualCosts := (fs.MonthlyTeamCost + fs.MonthlyComputeCost + fs.MonthlyODCCost) * 12

	var profitMargin float64
	if annualizedMRR > 0 {
		profitMargin = (float64(annualizedMRR-annualCosts) / float64(annualizedMRR)) * 100
	}

	return growthRate + profitMargin
}


func (fs *FounderState) CalculateBurnMultiple() float64 {
	if fs.Turn < 2 {
		return 0 // Need at least 2 months of data
	}

	// Monthly burn (if negative cash flow)
	monthlyRevenue := fs.MRR
	monthlyCosts := fs.MonthlyTeamCost + fs.MonthlyComputeCost + fs.MonthlyODCCost
	monthlyBurn := monthlyCosts - monthlyRevenue

	if monthlyBurn <= 0 {
		return 0 // Profitable, no burn
	}

	// New ARR = Growth * 12
	newMonthlyRevenue := int64(float64(fs.MRR) * fs.MonthlyGrowthRate)
	newARR := newMonthlyRevenue * 12

	if newARR <= 0 {
		return 999 // Burning with no growth
	}

	return float64(monthlyBurn) / float64(newARR)
}


func (fs *FounderState) CalculateMagicNumber() float64 {
	if fs.Turn < 2 {
		return 0
	}

	// Sales & Marketing spend = salaries for sales + marketing + any marketing campaigns
	salesMarketingCost := int64(0)
	for range fs.Team.Sales {
		salesMarketingCost += 100000 / 12 // $100k/year salary
	}
	for range fs.Team.Marketing {
		salesMarketingCost += 100000 / 12
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCGO {
			salesMarketingCost += 300000 / 12 // $300k/year
		}
	}

	if salesMarketingCost == 0 {
		return 0
	}

	// New revenue this quarter / S&M spend this quarter
	newQuarterlyRevenue := int64(float64(fs.MRR) * fs.MonthlyGrowthRate * 3) // 3 months

	return float64(newQuarterlyRevenue) / float64(salesMarketingCost*3)
}

func (fs *FounderState) GetCustomerHealthSegments() (healthy, atRisk, critical int, atRiskMRR, criticalMRR int64) {
	for i := range fs.CustomerList {
		c := &fs.CustomerList[i]
		if !c.IsActive {
			continue
		}

		// Update health score based on product maturity and CS team
		csImpact := float64(len(fs.Team.CustomerSuccess)) * 0.05
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCOO {
				csImpact += 0.15 // COO = 3x CS rep
			}
		}

		// Health improves with product maturity and CS team
		c.HealthScore = 0.3 + (fs.ProductMaturity * 0.5) + csImpact
		c.HealthScore = math.Min(1.0, math.Max(0.0, c.HealthScore))

		// Categorize
		if c.HealthScore >= 0.7 {
			healthy++
		} else if c.HealthScore >= 0.4 {
			atRisk++
			atRiskMRR += c.DealSize
		} else {
			critical++
			criticalMRR += c.DealSize
		}
	}

	return
}
