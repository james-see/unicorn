package founder

import (
	"fmt"
	"math/rand"
)


func (fs *FounderState) SpawnCompetitor() *Competitor {
	// 8% chance per month after month 6 (increased from 3% after month 12)
	// This ensures competitors spawn more reliably
	if fs.Turn < 6 || rand.Float64() > 0.08 {
		return nil
	}

	// Silicon Valley TV show startup names
	siliconValleyStartups := []struct {
		Name        string
		Description string
		ThreatBias  string // "low", "medium", "high" - bias towards this threat level
	}{
		{"Hooli", "The tech giant - always a high threat", "high"},
		{"Nucleus", "Compression technology competitor", "high"},
		{"End Frame", "Video streaming platform", "medium"},
		{"Aviato", "Erlich's previous startup", "low"},
		{"Bream-Hall", "VC-backed competitor", "medium"},
		{"Raviga", "Another VC-backed startup", "medium"},
		{"Action Jack's Company", "Enterprise-focused competitor", "high"},
		{"SeeFood", "AI/ML startup", "medium"},
		{"Optimal Tip-to-Tip", "Efficiency-focused startup", "low"},
		{"Pied Piper", "Wait, that's you! (or parallel universe)", "high"},
		{"Hooli XYZ", "Hooli's experimental division", "medium"},
		{"Hooli Chat", "Hooli's messaging platform", "medium"},
		{"Hooli Search", "Hooli's search engine", "high"},
		{"Hooli Box", "Hooli's hardware division", "low"},
		{"Hooli Connect", "Hooli's social network", "medium"},
		{"Gavin Belson's New Thing", "Gavin's latest venture", "high"},
		{"Bachmanity", "Erlich's other venture", "low"},
		{"Fiber", "Infrastructure startup", "medium"},
		{"Homicide", "Gaming platform", "low"},
		{"Bro", "Social app", "low"},
	}

	selected := siliconValleyStartups[rand.Intn(len(siliconValleyStartups))]
	
	// Determine threat level based on bias and randomness
	threat := selected.ThreatBias
	
	// 30% chance to deviate from bias
	if rand.Float64() < 0.3 {
		// Can go up or down one level
		if selected.ThreatBias == "high" && rand.Float64() < 0.5 {
			threat = "medium"
		} else if selected.ThreatBias == "low" && rand.Float64() < 0.5 {
			threat = "medium"
		} else if selected.ThreatBias == "medium" {
			if rand.Float64() < 0.5 {
				threat = "high"
			} else {
				threat = "low"
			}
		}
	}

	var marketShare float64
	switch threat {
	case "low":
		marketShare = 0.01 + rand.Float64()*0.04 // 1-5%
	case "medium":
		marketShare = 0.05 + rand.Float64()*0.10 // 5-15%
	case "high":
		marketShare = 0.10 + rand.Float64()*0.15 // 10-25%
	}

	comp := Competitor{
		Name:          selected.Name,
		Threat:        threat,
		MarketShare:   marketShare,
		Strategy:      "ignore", // Default strategy
		MonthAppeared: fs.Turn,
		Active:        true,
	}

	fs.Competitors = append(fs.Competitors, comp)
	return &comp
}


func (fs *FounderState) HandleCompetitor(compIndex int, strategy string) (string, error) {
	if compIndex < 0 || compIndex >= len(fs.Competitors) {
		return "", fmt.Errorf("invalid competitor index")
	}

	comp := &fs.Competitors[compIndex]
	if !comp.Active {
		return "", fmt.Errorf("competitor is no longer active")
	}

	switch strategy {
	case "ignore":
		comp.Strategy = "ignore"
		return fmt.Sprintf("Ignoring %s. They may gain market share.", comp.Name), nil

	case "compete":
		cost := int64(50000 + rand.Int63n(100000)) // $50-150k
		if cost > fs.Cash {
			return "", fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}

		fs.Cash -= cost
		comp.Strategy = "compete"

		// Reduce their threat
		if comp.Threat == "high" {
			comp.Threat = "medium"
		} else if comp.Threat == "medium" {
			comp.Threat = "low"
		}
		comp.MarketShare *= 0.7 // Reduce their market share

		return fmt.Sprintf("Competing aggressively! Cost: $%s. Reduced %s threat to %s",
			formatCurrency(cost), comp.Name, comp.Threat), nil

	case "partner":
		cost := int64(100000 + rand.Int63n(150000)) // $100-250k
		if cost > fs.Cash {
			return "", fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}

		fs.Cash -= cost
		comp.Strategy = "partner"
		comp.Active = false

		// Merge customer bases
		newCustomers := int(float64(fs.Customers) * comp.MarketShare)

		// Calculate MRR with variable deal sizes
		var totalMRR int64
		var dealSizes []int64 // Store deal sizes for customer tracking
		for i := 0; i < newCustomers; i++ {
			dealSize := generateDealSize(fs.AvgDealSize, fs.Category)
			fs.updateDealSizeRange(dealSize)
			totalMRR += dealSize
			dealSizes = append(dealSizes, dealSize)
		}

		// These are direct customers (acquired via partnership)
		fs.Customers += newCustomers
		fs.DirectCustomers += newCustomers
		fs.DirectMRR += totalMRR

		// Add customers to tracking system
		for _, dealSize := range dealSizes {
			fs.addCustomer(dealSize, "partnership")
		}

		// Sync MRR from DirectMRR + AffiliateMRR
		fs.syncMRR()

		// Recalculate average deal size
		if fs.Customers > 0 {
			fs.AvgDealSize = fs.MRR / int64(fs.Customers)
		}

		return fmt.Sprintf("Partnered with %s! Cost: $%s. Gained %d customers (+$%s MRR)",
			comp.Name, formatCurrency(cost), newCustomers, formatCurrency(totalMRR)), nil

	default:
		return "", fmt.Errorf("unknown strategy: %s", strategy)
	}
}


func (fs *FounderState) UpdateCompetitors() []string {
	var messages []string

	for i := range fs.Competitors {
		comp := &fs.Competitors[i]
		if !comp.Active {
			continue
		}

		// Ignored competitors grow stronger
		if comp.Strategy == "ignore" {
			comp.MarketShare *= 1.05 // Grow 5% per month

			// May increase threat level
			if comp.MarketShare > 0.20 && comp.Threat != "high" {
				comp.Threat = "high"
				messages = append(messages, fmt.Sprintf("⚠️  %s is now a HIGH threat!", comp.Name))
			} else if comp.MarketShare > 0.10 && comp.Threat == "low" {
				comp.Threat = "medium"
				messages = append(messages, fmt.Sprintf("⚠️  %s is now a MEDIUM threat", comp.Name))
			}
		}

		// Competing with them slows their growth
		if comp.Strategy == "compete" {
			comp.MarketShare *= 0.95 // Shrink 5% per month
			if comp.MarketShare < 0.02 {
				comp.Active = false
				messages = append(messages, fmt.Sprintf("✅ %s has exited the market!", comp.Name))
			}
		}

		// AI Competitor Actions - Silicon Valley companies are strategic
		if rand.Float64() < 0.15 { // 15% chance per month of competitor action
			actionType := rand.Float64()
			
			if actionType < 0.4 {
				// Steal customers (40% chance)
				customersStolen := int(float64(fs.Customers) * comp.MarketShare * (0.05 + rand.Float64()*0.10)) // 5-15% of their market share
				if customersStolen > 0 && fs.Customers > 0 {
					// Cap at reasonable amount
					if customersStolen > fs.Customers/10 {
						customersStolen = fs.Customers / 10
					}
					
					mrrLost := int64(customersStolen) * fs.AvgDealSize
					fs.Customers -= customersStolen
					fs.DirectCustomers -= customersStolen
					fs.MRR -= mrrLost
					fs.DirectMRR -= mrrLost
					
					// Hooli is especially aggressive
					if comp.Name == "Hooli" || comp.Name == "Hooli Search" {
						customersStolen = int(float64(customersStolen) * 1.5)
						mrrLost = int64(float64(mrrLost) * 1.5)
					}
					
					messages = append(messages, fmt.Sprintf("🔥 %s stole %d customers from you! (-$%s MRR)",
						comp.Name, customersStolen, formatCurrency(mrrLost)))
					comp.MarketShare *= 1.1 // Competitor gains market share
				}
			} else if actionType < 0.7 {
				// Launch competing product/feature (30% chance)
				// This increases their threat and market share
				if comp.Threat != "high" {
					comp.Threat = "high"
				}
				comp.MarketShare *= 1.15
				
				// Your CAC increases due to competition
				fs.BaseCAC = int64(float64(fs.BaseCAC) * 1.1)
				
				// Hooli launches "Box" style competing products
				if comp.Name == "Hooli" || comp.Name == "Hooli Box" {
					messages = append(messages, fmt.Sprintf("🚨 %s launched a competing product! Your CAC increased by 10%%",
						comp.Name))
				} else {
					messages = append(messages, fmt.Sprintf("🚨 %s launched a competing feature! Market competition increased",
						comp.Name))
				}
			} else if actionType < 0.85 {
				// Aggressive pricing/promotion (15% chance)
				// Competitor undercuts you, reducing your growth
				growthReduction := 0.05 + rand.Float64()*0.10 // 5-15% growth reduction
				fs.MonthlyGrowthRate *= (1.0 - growthReduction)
				
				if comp.Name == "Hooli" {
					messages = append(messages, fmt.Sprintf("💰 %s launched aggressive pricing! Your growth rate reduced by %.0f%%",
						comp.Name, growthReduction*100))
				} else {
					messages = append(messages, fmt.Sprintf("💰 %s is undercutting your prices! Growth slowed",
						comp.Name))
				}
			} else {
				// Poach talent (15% chance)
				// Competitor hires away employees
				if fs.Team.TotalEmployees > 0 {
					employeesPoached := 1
					if fs.Team.TotalEmployees > 10 {
						employeesPoached = 1 + rand.Intn(2)
					}
					
					// Remove employees (simplified - just reduce team cost)
					costReduction := int64(employeesPoached) * 8333 // ~$100k/year per employee
					fs.MonthlyTeamCost -= costReduction
					if fs.MonthlyTeamCost < 0 {
						fs.MonthlyTeamCost = 0
					}
					
					// Reduce productivity
					fs.ProductMaturity *= 0.98 // Slight reduction
					
					if comp.Name == "Hooli" || comp.Name == "Action Jack's Company" {
						messages = append(messages, fmt.Sprintf("👔 %s poached %d employee(s) from you! Team productivity reduced",
							comp.Name, employeesPoached))
					} else {
						messages = append(messages, fmt.Sprintf("👔 %s hired away %d employee(s)!",
							comp.Name, employeesPoached))
					}
				}
			}
		}

		// Hooli-specific behaviors - they're always up to something
		if (comp.Name == "Hooli" || comp.Name == "Gavin Belson's New Thing") && rand.Float64() < 0.1 {
			// 10% chance Hooli does something dramatic
			hooliActions := []string{
				"launched a massive marketing campaign",
				"announced a competing product at a major conference",
				"hired away a key executive",
				"filed a patent lawsuit",
			}
			action := hooliActions[rand.Intn(len(hooliActions))]
			
			switch action {
			case "launched a massive marketing campaign":
				comp.MarketShare *= 1.2
				fs.BaseCAC = int64(float64(fs.BaseCAC) * 1.15)
				messages = append(messages, fmt.Sprintf("📢 Hooli %s! Your CAC increased significantly", action))
			case "announced a competing product at a major conference":
				comp.Threat = "high"
				comp.MarketShare *= 1.25
				messages = append(messages, fmt.Sprintf("🎤 Hooli %s! Market competition intensified", action))
			case "hired away a key executive":
				fs.MonthlyTeamCost -= 25000 // Executive salary
				fs.ProductMaturity *= 0.95
				messages = append(messages, fmt.Sprintf("💼 Hooli %s! Product development slowed", action))
			case "filed a patent lawsuit":
				legalCost := int64(50000 + rand.Int63n(100000))
				fs.Cash -= legalCost
				messages = append(messages, fmt.Sprintf("⚖️  Hooli %s! Legal costs: $%s", action, formatCurrency(legalCost)))
			}
		}
	}

	return messages
}

// ============================================================================
// GLOBAL MARKETS
// ============================================================================


func (fs *FounderState) ExpandToMarket(region string) (*Market, error) {
	// Check if already in this market
	for _, m := range fs.GlobalMarkets {
		if m.Region == region {
			return nil, fmt.Errorf("already operating in %s", region)
		}
	}

	// Define market parameters
	var setupCost, monthlyCost int64
	var marketSize int
	var competition string

	switch region {
	case "Europe":
		setupCost = 200000
		monthlyCost = 30000
		marketSize = 50000
		competition = "high"
	case "Asia":
		setupCost = 250000
		monthlyCost = 40000
		marketSize = 100000
		competition = "very_high"
	case "LATAM":
		setupCost = 150000
		monthlyCost = 20000
		marketSize = 30000
		competition = "medium"
	case "Middle East":
		setupCost = 180000
		monthlyCost = 25000
		marketSize = 20000
		competition = "low"
	case "Africa":
		setupCost = 120000
		monthlyCost = 15000
		marketSize = 25000
		competition = "low"
	case "Australia":
		setupCost = 100000
		monthlyCost = 18000
		marketSize = 15000
		competition = "medium"
	default:
		return nil, fmt.Errorf("unknown market: %s", region)
	}

	if setupCost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(setupCost))
	}

	fs.Cash -= setupCost

	// Calculate initial customers based on setup cost / effective local CAC
	fs.UpdateCAC()
	localCAC := fs.CustomerAcquisitionCost

	// Adjust CAC for local competition
	switch competition {
	case "very_high":
		localCAC = int64(float64(localCAC) * 1.8)
	case "high":
		localCAC = int64(float64(localCAC) * 1.5)
	case "medium":
		localCAC = int64(float64(localCAC) * 1.2)
	case "low":
		localCAC = int64(float64(localCAC) * 0.8)
	}

	// Further adjust for product maturity (immature product = harder to sell internationally)
	if fs.ProductMaturity < 0.7 {
		localCAC = int64(float64(localCAC) / fs.ProductMaturity)
	}

	initialCustomers := int(setupCost / localCAC)
	if initialCustomers < 1 {
		initialCustomers = 1 // At least 1 customer
	}

	// Calculate MRR with variable deal sizes
	var initialMRR int64
	var dealSizes []int64 // Store deal sizes for customer tracking
	for i := 0; i < initialCustomers; i++ {
		dealSize := generateDealSize(fs.AvgDealSize, fs.Category)
		fs.updateDealSizeRange(dealSize)
		initialMRR += dealSize
		dealSizes = append(dealSizes, dealSize)
	}

	market := Market{
		Region:           region,
		LaunchMonth:      fs.Turn,
		SetupCost:        setupCost,
		MonthlyCost:      monthlyCost,
		CustomerCount:    initialCustomers,
		MRR:              initialMRR,
		MarketSize:       marketSize,
		Penetration:      float64(initialCustomers) / float64(marketSize),
		LocalCompetition: competition,
	}

	fs.GlobalMarkets = append(fs.GlobalMarkets, market)
	// These are direct customers (from market expansion)
	fs.Customers += initialCustomers
	fs.DirectCustomers += initialCustomers
	fs.DirectMRR += initialMRR

	// Add customers to tracking system
	for _, dealSize := range dealSizes {
		fs.addCustomer(dealSize, "market")
	}

	// Sync MRR from DirectMRR + AffiliateMRR
	fs.syncMRR()

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

	// Increase global churn rate due to operational complexity
	fs.CustomerChurnRate += 0.01 + (rand.Float64() * 0.01)

	// Add regional competitors based on competition level
	numCompetitors := 0
	switch competition {
	case "very_high":
		numCompetitors = 2 + rand.Intn(2) // 2-3 competitors
	case "high":
		numCompetitors = 1 + rand.Intn(2) // 1-2 competitors
	case "medium":
		numCompetitors = rand.Intn(2) // 0-1 competitors
	case "low":
		numCompetitors = 0 // No competitors
	}

	// Regional competitor names by market
	regionalCompetitors := map[string][]string{
		"Europe":      {"Zalando", "Klarna", "N26", "TransferWise", "BlaBlaCar", "Deliveroo EU"},
		"Asia":        {"Grab", "GoJek", "Alibaba Local", "Meituan", "Tokopedia", "Paytm"},
		"LATAM":       {"Nubank", "Mercado Libre", "Rappi", "Kavak", "Creditas", "QuintoAndar"},
		"Middle East": {"Careem", "Souq", "Fetchr", "Talabat", "Noon", "Swvl"},
		"Africa":      {"Jumia", "Flutterwave", "Andela", "Paystack", "Konga", "M-Pesa"},
		"Australia":   {"Afterpay", "Canva AU", "Atlassian", "WiseTech", "Xero", "SEEK"},
	}

	if names, ok := regionalCompetitors[region]; ok && numCompetitors > 0 {
		for i := 0; i < numCompetitors && i < len(names); i++ {
			compName := names[rand.Intn(len(names))]

			// Check if competitor already exists
			exists := false
			for _, existing := range fs.Competitors {
				if existing.Name == compName {
					exists = true
					break
				}
			}

			if !exists {
				threatLevel := "medium"
				marketShare := 0.05 + rand.Float64()*0.15 // 5-20% market share

				switch competition {
				case "very_high":
					threatLevel = "high"
					marketShare = 0.10 + rand.Float64()*0.20 // 10-30%
				case "high":
					threatLevel = "high"
					marketShare = 0.08 + rand.Float64()*0.15 // 8-23%
				case "medium":
					threatLevel = "medium"
					marketShare = 0.05 + rand.Float64()*0.10 // 5-15%
				}

				competitor := Competitor{
					Name:          compName + " (" + region + ")",
					Threat:        threatLevel,
					MarketShare:   marketShare,
					Strategy:      "ignore", // Default strategy
					MonthAppeared: fs.Turn,
					Active:        true,
				}
				fs.Competitors = append(fs.Competitors, competitor)
			}
		}
	}

	return &market, nil
}


func (fs *FounderState) UpdateGlobalMarkets() []string {
	var messages []string

	for i := range fs.GlobalMarkets {
		m := &fs.GlobalMarkets[i]

		// Pay monthly costs
		fs.Cash -= m.MonthlyCost

		// Calculate market-specific churn rate
		marketChurn := fs.CustomerChurnRate

		// Count CS team assigned to this market or "All" markets
		csAssignedToMarket := 0
		for _, cs := range fs.Team.CustomerSuccess {
			if cs.AssignedMarket == m.Region || cs.AssignedMarket == "All" {
				csAssignedToMarket++
			}
		}

		// COO also helps with churn across all markets
		hasCOO := false
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCOO {
				hasCOO = true
				break
			}
		}

		// No CS team assigned to this market? Much higher churn (up to 50% in new markets)
		if csAssignedToMarket == 0 && !hasCOO {
			marketChurn += 0.30 // +30% base churn without CS in this market
		} else {
			// Each CS rep assigned to market reduces churn by ~2%
			churnReduction := float64(csAssignedToMarket) * 0.02
			if hasCOO {
				churnReduction += 0.06 // COO equivalent to 3 CS reps
			}
			marketChurn -= churnReduction
			if marketChurn < 0.01 {
				marketChurn = 0.01 // Min 1% churn
			}
		}

		// Product not mature? Higher churn
		if fs.ProductMaturity < 1.0 {
			marketChurn += (1.0 - fs.ProductMaturity) * 0.20 // Up to +20% if product immature
		}

		// No engineers? Product degrades, churn increases
		if len(fs.Team.Engineers) == 0 {
			marketChurn += 0.15 // +15% without engineers
		}

		// Local competition increases churn
		switch m.LocalCompetition {
		case "very_high":
			marketChurn += 0.10
		case "high":
			marketChurn += 0.07
		case "medium":
			marketChurn += 0.04
		}

		// Process churn first
		customersLost := int(float64(m.CustomerCount) * marketChurn)
		if customersLost > 0 {
			// Estimate MRR lost using average deal size
			mrrLost := int64(customersLost) * fs.AvgDealSize

			m.CustomerCount -= customersLost
			m.MRR -= mrrLost
			// These are direct customers (from market expansion)
			fs.MRR -= mrrLost
			fs.DirectMRR -= mrrLost
			fs.Customers -= customersLost
			fs.DirectCustomers -= customersLost

			// Recalculate average deal size
			if fs.Customers > 0 {
				fs.AvgDealSize = fs.MRR / int64(fs.Customers)
			}

			messages = append(messages, fmt.Sprintf("📉 %s: Lost %d customers (%.1f%% churn)",
				m.Region, customersLost, marketChurn*100))
		}

		// Competitors actively take customers
		for _, comp := range fs.Competitors {
			if !comp.Active {
				continue
			}
			// Competitors are more effective in markets where you're weak
			if comp.Strategy == "ignore" && len(fs.Team.Sales) < 3 {
				competitorSteal := int(float64(m.CustomerCount) * comp.MarketShare * 0.15) // 15% of their market share
				if competitorSteal > 0 {
					// Estimate MRR stolen using average deal size
					stolenMRR := int64(competitorSteal) * fs.AvgDealSize

					m.CustomerCount -= competitorSteal
					m.MRR -= stolenMRR
					// These are direct customers (from market expansion)
					fs.MRR -= stolenMRR
					fs.DirectMRR -= stolenMRR
					fs.Customers -= competitorSteal
					fs.DirectCustomers -= competitorSteal

					// Recalculate average deal size
					if fs.Customers > 0 {
						fs.AvgDealSize = fs.MRR / int64(fs.Customers)
					}

					messages = append(messages, fmt.Sprintf("⚠️  %s took %d customers in %s",
						comp.Name, competitorSteal, m.Region))
				}
			}
		}

		// Now attempt growth
		// Base growth rate - you need sales/marketing to grow in new markets
		baseGrowth := 0.03 + (rand.Float64() * 0.02) // 3-5% base monthly growth

		// Count employees assigned to this market or "All" markets
		salesInMarket := 0
		for _, sales := range fs.Team.Sales {
			if sales.AssignedMarket == m.Region || sales.AssignedMarket == "All" {
				salesInMarket++
			}
		}

		marketingInMarket := 0
		for _, marketing := range fs.Team.Marketing {
			if marketing.AssignedMarket == m.Region || marketing.AssignedMarket == "All" {
				marketingInMarket++
			}
		}

		csInMarket := 0
		for _, cs := range fs.Team.CustomerSuccess {
			if cs.AssignedMarket == m.Region || cs.AssignedMarket == "All" {
				csInMarket++
			}
		}

		// Sales team impact (critical for growing in new markets)
		salesImpact := float64(salesInMarket) * 0.05 // Each sales rep adds 5%

		// CGO amplifies sales impact (CGO works across all markets)
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCGO {
				salesImpact += (exec.Impact * 0.05) // CGO adds significant growth boost
			}
		}

		// Marketing team impact (brand awareness in new markets)
		marketingImpact := 0.03 * float64(marketingInMarket) // Each marketer adds 3%

		// Adjust for competition
		competitionMultiplier := 1.0
		switch m.LocalCompetition {
		case "very_high":
			competitionMultiplier = 0.6 // Harder growth in very competitive markets
		case "high":
			competitionMultiplier = 0.75
		case "medium":
			competitionMultiplier = 0.9
		case "low":
			competitionMultiplier = 1.2 // Much easier growth in low competition
		}

		// Product maturity affects conversion
		productMultiplier := fs.ProductMaturity
		if productMultiplier < 0.5 {
			productMultiplier = 0.5 // Can't grow much with immature product
		}

		totalGrowthRate := (baseGrowth + salesImpact + marketingImpact) * competitionMultiplier * productMultiplier

		// Calculate new customers
		// Percentage growth based on current base (compounds over time)
		percentageGrowth := int(float64(m.CustomerCount) * totalGrowthRate)

		// Plus absolute growth (helps new/small markets grow)
		// Sales/marketing teams directly acquire customers even in new markets
		// Only count employees assigned to this market or "All"
		absoluteGrowth := (salesInMarket * 2) + (marketingInMarket * 1)

		// CGO contributes to absolute growth (works across all markets)
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCGO {
				absoluteGrowth += int(exec.Impact * 3) // CGO brings in customers directly
			}
		}

		newCustomers := percentageGrowth + absoluteGrowth

		// But cap at remaining market opportunity
		remainingMarket := m.MarketSize - m.CustomerCount
		if newCustomers > remainingMarket {
			newCustomers = remainingMarket
		}

		if newCustomers > 0 {
			// Calculate MRR with variable deal sizes
			var totalMRR int64
			for i := 0; i < newCustomers; i++ {
				dealSize := generateDealSize(fs.AvgDealSize, fs.Category)
				fs.updateDealSizeRange(dealSize)
				totalMRR += dealSize
			}

			m.CustomerCount += newCustomers
			m.MRR += totalMRR
			// These are direct customers (from market expansion)
			fs.DirectMRR += totalMRR
			fs.Customers += newCustomers
			fs.DirectCustomers += newCustomers

			// Sync MRR from DirectMRR + AffiliateMRR
			fs.syncMRR()

			// Recalculate average deal size
			if fs.Customers > 0 {
				fs.AvgDealSize = fs.MRR / int64(fs.Customers)
			}

			m.Penetration = float64(m.CustomerCount) / float64(m.MarketSize)

			messages = append(messages, fmt.Sprintf("🌍 %s: +%d customers, $%s MRR (%.1f%% penetration)",
				m.Region, newCustomers, formatCurrency(m.MRR), m.Penetration*100))
		}

		// Market can shrink if churn > growth
		if m.CustomerCount < 0 {
			m.CustomerCount = 0
			m.MRR = 0
		}
	}

	return messages
}

// ============================================================================
// PIVOTS
// ============================================================================


func (fs *FounderState) ExecutePivot(toStrategy string, reason string) (*Pivot, error) {
	cost := int64(100000 + rand.Int63n(200000)) // $100-300k
	if cost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	fs.Cash -= cost

	// Lose customers during pivot (20-50%)
	churnRate := 0.20 + rand.Float64()*0.30
	customersLost := int(float64(fs.Customers) * churnRate)
	mrrLost := int64(customersLost) * fs.AvgDealSize

	// Apply churn proportionally to direct and affiliate customers
	directCustomersLost := int(float64(fs.DirectCustomers) * churnRate)
	affiliateCustomersLost := int(float64(fs.AffiliateCustomers) * churnRate)
	directMRRLost := int64(float64(fs.DirectMRR) * churnRate)
	affiliateMRRLost := int64(float64(fs.AffiliateMRR) * churnRate)

	fs.Customers -= customersLost
	fs.DirectCustomers -= directCustomersLost
	fs.AffiliateCustomers -= affiliateCustomersLost
	fs.MRR -= mrrLost
	fs.DirectMRR -= directMRRLost
	fs.AffiliateMRR -= affiliateMRRLost

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

	// Success rate depends on product maturity and timing
	successChance := 0.30 // Base 30%
	if fs.ProductMaturity > 0.7 {
		successChance += 0.20 // +20% if product mature
	}
	if fs.Turn < 24 {
		successChance += 0.20 // +20% if early (under 2 years)
	}

	success := rand.Float64() < successChance

	pivot := Pivot{
		Month:         fs.Turn,
		FromStrategy:  fs.StartupType,
		ToStrategy:    toStrategy,
		Reason:        reason,
		Cost:          cost,
		CustomersLost: customersLost,
		Success:       success,
	}

	if success {
		fs.StartupType = toStrategy
		// Expand market size on successful pivot
		fs.TargetMarketSize = int(float64(fs.TargetMarketSize) * (1.5 + rand.Float64()))
	}

	fs.PivotHistory = append(fs.PivotHistory, pivot)
	fs.CalculateRunway()

	return &pivot, nil
}

// ============================================================================
// EQUITY BUYBACKS
// ============================================================================


func (fs *FounderState) SpawnRandomEvent() *RandomEvent {
	// Don't spawn events in first 6 months
	if fs.Turn < 6 {
		return nil
	}

	// Choose event type
	eventTypes := []string{"economy", "regulation", "competition", "talent", "customer", "product", "legal", "press"}
	eventType := eventTypes[rand.Intn(len(eventTypes))]

	// 40% chance of positive event
	isPositive := rand.Float64() < 0.4

	var event *RandomEvent

	switch eventType {
	case "economy":
		event = fs.generateEconomyEvent(isPositive)
	case "regulation":
		event = fs.generateRegulationEvent(isPositive)
	case "competition":
		event = fs.generateCompetitionEvent(isPositive)
	case "talent":
		event = fs.generateTalentEvent(isPositive)
	case "customer":
		event = fs.generateCustomerEvent(isPositive)
	case "product":
		event = fs.generateProductEvent(isPositive)
	case "legal":
		event = fs.generateLegalEvent(isPositive)
	case "press":
		event = fs.generatePressEvent(isPositive)
	}

	if event != nil {
		event.Month = fs.Turn
		event.Type = eventType
		
		// Check if chairman can mitigate negative events
		if !event.IsPositive {
			if fs.MitigateCrisis(event) {
				// Event was mitigated - add message about chairman's help
				chairman := fs.GetChairman()
				if chairman != nil {
					event.Description += fmt.Sprintf(" (Chairman %s helped mitigate impact)", chairman.Name)
				}
			}
		}

		fs.RandomEvents = append(fs.RandomEvents, *event)

		// Apply the effect (may have been mitigated by chairman)
		if event.Impact.DurationMonths > 0 {
			fs.ActiveEventEffects[event.Title] = event.Impact
		}

		// Apply immediate cash cost (may have been reduced by chairman)
		if event.Impact.CashCost > 0 {
			fs.Cash -= event.Impact.CashCost
		}

		// Handle employee losses
		if event.Impact.EmployeesLost > 0 {
			fs.handleEmployeeLoss(event.Impact.EmployeesLost)
		}
	}

	return event
}

func (fs *FounderState) generateEconomyEvent(isPositive bool) *RandomEvent {
	if isPositive {
		// Economic boom
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Economic Boom",
			Description: "Strong economy increases customer budgets and spending",
			Impact: EventImpact{
				GrowthChange:   1.15, // +15% growth
				DurationMonths: 6,
			},
		}
	} else {
		// Economic downturn
		severity := "moderate"
		impact := EventImpact{
			GrowthChange:   0.85, // -15% growth
			ChurnChange:    0.02, // +2% churn
			DurationMonths: 6,
		}

		// Tariffs specifically affect hardware companies
		if fs.Category == "Hardware" || fs.Category == "Deep Tech" {
			severity = "major"
			impact.CACChange = 1.3 // +30% CAC due to tariffs
			impact.CashCost = 20000 + rand.Int63n(30000)
			return &RandomEvent{
				Severity:    severity,
				IsPositive:  false,
				Title:       "Trade Tariffs Imposed",
				Description: "New tariffs on hardware components increase costs significantly",
				Impact:      impact,
			}
		}

		return &RandomEvent{
			Severity:    severity,
			IsPositive:  false,
			Title:       "Economic Downturn",
			Description: "Recession fears cause customers to cut budgets",
			Impact:      impact,
		}
	}
}

func (fs *FounderState) generateRegulationEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "minor",
			IsPositive:  true,
			Title:       "Favorable Regulations Passed",
			Description: "New laws favor your business model and reduce compliance costs",
			Impact: EventImpact{
				GrowthChange:   1.10,
				CashCost:       -10000, // Negative cost = gain
				DurationMonths: 12,
			},
		}
	} else {
		impact := EventImpact{
			CashCost:       50000 + rand.Int63n(100000), // $50-150k compliance
			GrowthChange:   0.90,                        // -10% growth
			DurationMonths: 12,
		}

		// Data privacy regulations
		if fs.Category == "SaaS" {
			impact.CashCost += 50000
			return &RandomEvent{
				Severity:    "major",
				IsPositive:  false,
				Title:       "New Data Privacy Regulations",
				Description: "GDPR/CCPA-style laws require expensive compliance changes",
				Impact:      impact,
			}
		}

		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "New Industry Regulations",
			Description: "New compliance requirements slow growth and increase costs",
			Impact:      impact,
		}
	}
}

func (fs *FounderState) generateCompetitionEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Major Competitor Exits Market",
			Description: "A key competitor shuts down, opening up opportunities",
			Impact: EventImpact{
				GrowthChange:   1.25, // +25% growth
				CACChange:      0.80, // -20% CAC
				DurationMonths: 6,
			},
		}
	} else {
		// Check for open source threat
		if rand.Float64() < 0.3 {
			return &RandomEvent{
				Severity:    "major",
				IsPositive:  false,
				Title:       "Open Source Alternative Launched",
				Description: "Free alternative threatens paid product - must differentiate!",
				Impact: EventImpact{
					ChurnChange:    0.05, // +5% churn
					GrowthChange:   0.70, // -30% growth
					CACChange:      1.40, // +40% CAC
					DurationMonths: 12,
				},
			}
		}

		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "Well-Funded Competitor Emerges",
			Description: "New competitor raises $50M Series B, plans aggressive expansion",
			Impact: EventImpact{
				CACChange:      1.25, // +25% CAC
				ChurnChange:    0.03, // +3% churn
				DurationMonths: 9,
			},
		}
	}
}

func (fs *FounderState) generateTalentEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "minor",
			IsPositive:  true,
			Title:       "Rock Star Candidate Interested",
			Description: "Top talent from FAANG wants to join - hire them!",
			Impact: EventImpact{
				ProductivityChange: 1.15, // +15% team productivity
				DurationMonths:     12,
			},
		}
	} else {
		// Key employees quit
		numQuitting := 1
		if fs.Team.TotalEmployees > 10 {
			numQuitting = 1 + rand.Intn(2)
		}

		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "Key Employees Resign",
			Description: fmt.Sprintf("%d employee(s) leave for higher-paying competitors", numQuitting),
			Impact: EventImpact{
				EmployeesLost:      numQuitting,
				ProductivityChange: 0.90, // -10% productivity
				DurationMonths:     3,    // Takes 3 months to recover
			},
		}
	}
}

func (fs *FounderState) generateCustomerEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  true,
			Title:       "Enterprise Customer Win!",
			Description: "Fortune 500 company signs major contract",
			Impact: EventImpact{
				CashCost:       -(50000 + rand.Int63n(150000)), // $50-200k deal
				DurationMonths: 0,                              // One-time
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  false,
			Title:       "Major Customer Churns",
			Description: "Top 3 customer cancels unexpectedly",
			Impact: EventImpact{
				MRRChange:      0.90, // -10% MRR
				DurationMonths: 0,
			},
		}
	}
}

func (fs *FounderState) generateProductEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Product Breakthrough!",
			Description: "Engineering team ships game-changing feature",
			Impact: EventImpact{
				ChurnChange:    -0.02, // -2% churn
				GrowthChange:   1.20,  // +20% growth
				DurationMonths: 6,
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  false,
			Title:       "Critical Bug in Production",
			Description: "Major outage affects all customers for 48 hours",
			Impact: EventImpact{
				ChurnChange:    0.05,                       // +5% churn
				CashCost:       10000 + rand.Int63n(20000), // Emergency fixes
				DurationMonths: 1,
			},
		}
	}
}

func (fs *FounderState) generateLegalEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "minor",
			IsPositive:  true,
			Title:       "Patent Approved",
			Description: "Key technology patent granted, provides competitive moat",
			Impact: EventImpact{
				GrowthChange:   1.10,
				DurationMonths: 24,
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  false,
			Title:       "Patent Infringement Claim",
			Description: "Large company alleges patent violation, requires legal defense",
			Impact: EventImpact{
				CashCost:       100000 + rand.Int63n(200000), // $100-300k legal fees
				DurationMonths: 0,
			},
		}
	}
}

func (fs *FounderState) generatePressEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Major Press Coverage",
			Description: "Featured in TechCrunch/WSJ - huge brand boost!",
			Impact: EventImpact{
				CACChange:      0.75, // -25% CAC
				GrowthChange:   1.30, // +30% growth
				DurationMonths: 3,
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "PR Crisis",
			Description: "Negative press about company culture or product issues",
			Impact: EventImpact{
				ChurnChange:    0.04, // +4% churn
				CACChange:      1.35, // +35% CAC
				GrowthChange:   0.80, // -20% growth
				DurationMonths: 6,
			},
		}
	}
}


func (fs *FounderState) ProcessRandomEvents() []string {
	var messages []string

	// Expire old events
	for key, impact := range fs.ActiveEventEffects {
		// Find the original event
		var originalEvent *RandomEvent
		for i := range fs.RandomEvents {
			if fs.RandomEvents[i].Title == key {
				originalEvent = &fs.RandomEvents[i]
				break
			}
		}

		if originalEvent != nil {
			monthsActive := fs.Turn - originalEvent.Month
			if monthsActive >= impact.DurationMonths {
				delete(fs.ActiveEventEffects, key)
				messages = append(messages, fmt.Sprintf("⏰ Event effect expired: %s", key))
			}
		}
	}

	// Apply active event effects (these are already factored into calculations)
	// The effects modify CAC, churn, growth rates which are used in ProcessMonth

	return messages
}


func (fs *FounderState) handleEmployeeLoss(count int) {
	for i := 0; i < count; i++ {
		// Randomly pick a team to lose from
		roll := rand.Intn(4)
		switch roll {
		case 0:
			if len(fs.Team.Engineers) > 0 {
				fs.Team.Engineers = fs.Team.Engineers[:len(fs.Team.Engineers)-1]
			}
		case 1:
			if len(fs.Team.Sales) > 0 {
				fs.Team.Sales = fs.Team.Sales[:len(fs.Team.Sales)-1]
			}
		case 2:
			if len(fs.Team.CustomerSuccess) > 0 {
				fs.Team.CustomerSuccess = fs.Team.CustomerSuccess[:len(fs.Team.CustomerSuccess)-1]
			}
		case 3:
			if len(fs.Team.Marketing) > 0 {
				fs.Team.Marketing = fs.Team.Marketing[:len(fs.Team.Marketing)-1]
			}
		}
	}
	fs.CalculateTeamCost()
	fs.CalculateRunway()
}