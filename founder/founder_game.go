package founder

import (
	"fmt"
	"math"
	"math/rand"
)


func (fs *FounderState) IsGameOver() bool {
	return fs.Cash <= 0 || fs.Turn > fs.MaxTurns || fs.HasExited
}


func (fs *FounderState) GetAvailableExits() []ExitOption {
	var exits []ExitOption
	founderEquity := (100.0 - fs.EquityPool - fs.EquityGivenAway) / 100.0

	// Calculate current valuation (simplified: ARR * multiple based on growth/profitability)
	arr := fs.MRR * 12
	multiple := 10.0                 // Base multiple
	if fs.MonthlyGrowthRate > 0.05 { // >5% monthly growth
		multiple += 5.0
	}
	if fs.MRR > fs.MonthlyTeamCost { // Profitable
		multiple += 3.0
	}
	currentValuation := int64(float64(arr) * multiple)

	// 1. IPO Option
	ipoReqs := []string{}
	canIPO := true

	if arr < 20000000 { // $20M ARR minimum
		ipoReqs = append(ipoReqs, fmt.Sprintf("❌ Need $20M ARR (currently $%s)", formatCurrency(arr)))
		canIPO = false
	} else {
		ipoReqs = append(ipoReqs, "✓ $20M+ ARR")
	}

	if len(fs.FundingRounds) < 2 {
		ipoReqs = append(ipoReqs, "❌ Need at least Series A funding")
		canIPO = false
	} else {
		ipoReqs = append(ipoReqs, "✓ Multiple funding rounds")
	}

	if fs.MonthlyGrowthRate < 0.03 { // <3% monthly = <40% YoY
		ipoReqs = append(ipoReqs, "❌ Need 40%+ YoY growth")
		canIPO = false
	} else {
		ipoReqs = append(ipoReqs, "✓ Strong growth rate")
	}

	ipoValuation := int64(float64(currentValuation) * 1.3)                 // 30% IPO premium
	ipoFounderPayout := int64(float64(ipoValuation) * founderEquity * 0.2) // Can sell 20% at IPO

	exits = append(exits, ExitOption{
		Type:          "ipo",
		Valuation:     ipoValuation,
		FounderPayout: ipoFounderPayout,
		Description:   "Take company public on NASDAQ. Provides liquidity but you remain CEO.",
		Requirements:  ipoReqs,
		CanExit:       canIPO,
	})

	// 2. Strategic Acquisition
	acqReqs := []string{}
	canAcquire := true

	if arr < 5000000 { // $5M ARR minimum
		acqReqs = append(acqReqs, fmt.Sprintf("❌ Need $5M ARR (currently $%s)", formatCurrency(arr)))
		canAcquire = false
	} else {
		acqReqs = append(acqReqs, "✓ $5M+ ARR")
	}

	if fs.Customers < 50 {
		acqReqs = append(acqReqs, fmt.Sprintf("❌ Need 50+ customers (currently %d)", fs.Customers))
		canAcquire = false
	} else {
		acqReqs = append(acqReqs, "✓ Significant customer base")
	}

	acqValuation := int64(float64(currentValuation) * 1.1) // 10% acquisition premium
	acqFounderPayout := int64(float64(acqValuation) * founderEquity)

	exits = append(exits, ExitOption{
		Type:          "acquisition",
		Valuation:     acqValuation,
		FounderPayout: acqFounderPayout,
		Description:   "Sell company to strategic acquirer. Full liquidity, but you're done.",
		Requirements:  acqReqs,
		CanExit:       canAcquire,
	})

	// 3. Secondary Sale (Private Equity)
	secondaryReqs := []string{}
	canSecondary := true

	if arr < 10000000 {
		secondaryReqs = append(secondaryReqs, fmt.Sprintf("❌ Need $10M ARR (currently $%s)", formatCurrency(arr)))
		canSecondary = false
	} else {
		secondaryReqs = append(secondaryReqs, "✓ $10M+ ARR")
	}

	netIncome := fs.MRR - fs.MonthlyTeamCost - fs.MonthlyComputeCost - fs.MonthlyODCCost
	if netIncome < 0 {
		secondaryReqs = append(secondaryReqs, "❌ Must be profitable")
		canSecondary = false
	} else {
		secondaryReqs = append(secondaryReqs, "✓ Profitable")
	}

	secondaryValuation := currentValuation
	secondaryFounderPayout := int64(float64(secondaryValuation) * founderEquity * 0.5) // Sell 50% of your stake

	exits = append(exits, ExitOption{
		Type:          "secondary",
		Valuation:     secondaryValuation,
		FounderPayout: secondaryFounderPayout,
		Description:   "Sell 50% of your shares to private equity. Partial liquidity, stay as CEO.",
		Requirements:  secondaryReqs,
		CanExit:       canSecondary,
	})

	// 4. Continue building
	exits = append(exits, ExitOption{
		Type:          "continue",
		Valuation:     currentValuation,
		FounderPayout: 0,
		Description:   "Keep building. Your current company value.",
		Requirements:  []string{"No requirements - always available"},
		CanExit:       true,
	})

	return exits
}


func (fs *FounderState) ExecuteExit(exitType string) {
	fs.HasExited = true
	fs.ExitType = exitType
	fs.ExitMonth = fs.Turn

	exits := fs.GetAvailableExits()
	for _, exit := range exits {
		if exit.Type == exitType {
			fs.ExitValuation = exit.Valuation
			break
		}
	}
}


func (fs *FounderState) GetFinalScore() (outcome string, valuation int64, founderEquity float64) {
	founderEquity = 100.0 - fs.EquityPool - fs.EquityGivenAway

	// Calculate final valuation based on MRR
	if fs.MRR > 0 {
		multiple := 10.0 // 10x ARR default
		if fs.MonthlyGrowthRate > 0.20 {
			multiple = 15.0
		}
		if fs.MonthlyGrowthRate < 0 {
			multiple = 5.0
		}
		valuation = int64(float64(fs.MRR) * 12 * multiple)
	}

	// Determine outcome
	if fs.Cash <= 0 {
		outcome = "SHUT DOWN - Ran out of cash"
	} else if fs.Turn > fs.MaxTurns {
		if fs.MRR > 1000000 { // $1M+ MRR
			outcome = "UNICORN PATH - Scaled to $1M+ MRR!"
		} else if fs.MRR > 100000 { // $100K+ MRR
			outcome = "SUCCESS - Built a real business"
		} else if fs.MRR > 0 {
			outcome = "SURVIVING - Making revenue but still growing"
		} else {
			outcome = "STRUGGLING - No meaningful traction"
		}
	}

	return outcome, valuation, founderEquity
}


func (fs *FounderState) CheckForAcquisition() *AcquisitionOffer {
	// Only after Series A and if metrics are good
	hasSeriesA := false
	for _, round := range fs.FundingRounds {
		if round.RoundName == "Series A" || round.RoundName == "Series B" {
			hasSeriesA = true
			break
		}
	}

	if !hasSeriesA {
		return nil
	}

	// Check for competitor acquisition first (higher priority)
	competitorOffer := fs.CheckForCompetitorAcquisition()
	if competitorOffer != nil {
		return competitorOffer
	}

	// 5% chance per month after Series A for strategic acquirer
	if rand.Float64() > 0.05 {
		return nil
	}

	// Calculate offer
	multiple := 3.0 + rand.Float64()*3.0 // 3-6x revenue
	annualRevenue := fs.MRR * 12
	offerAmount := int64(float64(annualRevenue) * multiple)

	dueDiligence := "normal"
	termsQuality := "good"

	roll := rand.Float64()
	if roll < 0.15 {
		dueDiligence = "bad"
		termsQuality = "poor"
		offerAmount = int64(float64(offerAmount) * 0.6)
	} else if roll > 0.85 {
		dueDiligence = "good"
		termsQuality = "excellent"
		offerAmount = int64(float64(offerAmount) * 1.3)
	}

	acquirers := []string{
		"Google", "Microsoft", "Amazon", "Salesforce", "Oracle",
		"Meta", "Apple", "Adobe", "SAP", "IBM",
	}

	offer := AcquisitionOffer{
		Acquirer:     acquirers[rand.Intn(len(acquirers))],
		OfferAmount:  offerAmount,
		Month:        fs.Turn,
		DueDiligence: dueDiligence,
		TermsQuality: termsQuality,
	}

	return &offer
}

// CheckForCompetitorAcquisition checks if any competitor wants to acquire your startup
func (fs *FounderState) CheckForCompetitorAcquisition() *AcquisitionOffer {
	// Only check active competitors with high market share
	for i := range fs.Competitors {
		comp := &fs.Competitors[i]
		if !comp.Active {
			continue
		}

		// Competitors can only acquire if they have significant market share (15%+)
		if comp.MarketShare < 0.15 {
			continue
		}

		// Higher threat competitors are more likely to acquire
		acquisitionChance := 0.02 // 2% base chance per month
		if comp.Threat == "high" {
			acquisitionChance = 0.04 // 4% for high threat
		} else if comp.Threat == "medium" {
			acquisitionChance = 0.03 // 3% for medium threat
		}

		// Silicon Valley companies have different acquisition behaviors
		// Hooli is very aggressive
		if comp.Name == "Hooli" || comp.Name == "Hooli Search" || comp.Name == "Gavin Belson's New Thing" {
			acquisitionChance *= 1.5 // Hooli is more likely to acquire
		}
		
		// Nucleus is competitive
		if comp.Name == "Nucleus" {
			acquisitionChance *= 1.3
		}

		// Competitors are more likely to acquire if you're struggling
		if fs.CashRunwayMonths < 6 {
			acquisitionChance *= 2.0 // Double chance if you're running low on cash
		}

		// Competitors are more likely if you're growing fast (they want to eliminate competition)
		if fs.MonthlyGrowthRate > 0.20 {
			acquisitionChance *= 1.5
		}

		// Hooli especially wants to acquire fast-growing competitors
		if (comp.Name == "Hooli" || comp.Name == "Gavin Belson's New Thing") && fs.MonthlyGrowthRate > 0.25 {
			acquisitionChance *= 2.0
		}

		if rand.Float64() > acquisitionChance {
			continue
		}

		// Calculate offer - competitors typically offer less than strategic acquirers
		// They're buying to eliminate competition, not for strategic value
		multiple := 2.0 + rand.Float64()*2.5 // 2-4.5x revenue (lower than strategic)
		annualRevenue := fs.MRR * 12
		offerAmount := int64(float64(annualRevenue) * multiple)

		// Hooli sometimes makes aggressive offers (but usually lowballs)
		if comp.Name == "Hooli" && rand.Float64() < 0.2 {
			// 20% chance Hooli makes a "Gavin Belson" style aggressive offer
			offerAmount = int64(float64(offerAmount) * 1.5)
		}

		// Competitor offers are usually less favorable
		dueDiligence := "normal"
		termsQuality := "good"
		if rand.Float64() < 0.3 {
			dueDiligence = "bad"
			termsQuality = "poor"
			offerAmount = int64(float64(offerAmount) * 0.7) // Competitors lowball more often
		}

		// Hooli is known for bad terms
		if comp.Name == "Hooli" || comp.Name == "Gavin Belson's New Thing" {
			if rand.Float64() < 0.5 {
				dueDiligence = "bad"
				termsQuality = "poor"
			}
		}

		offer := AcquisitionOffer{
			Acquirer:     comp.Name,
			OfferAmount:  offerAmount,
			Month:        fs.Turn,
			DueDiligence: dueDiligence,
			TermsQuality: termsQuality,
			IsCompetitor: true, // Mark as competitor acquisition
		}

		// Competitor exits market after acquisition (unless it's Hooli - they keep competing)
		if comp.Name != "Hooli" && comp.Name != "Hooli Search" && comp.Name != "Gavin Belson's New Thing" {
			comp.Active = false
			comp.MarketShare = 0
		}

		return &offer
	}

	return nil
}


func (fs *FounderState) NeedsLowCashWarning() bool {
	return fs.Cash <= 200000 && fs.CashRunwayMonths < 6
}


func (fs *FounderState) GenerateMonthlyHighlights() []MonthlyHighlight {
	var highlights []MonthlyHighlight

	// WINS
	if fs.MRR >= 100000 && fs.Turn > 1 {
		// Check if we just crossed $100k
		prevMRR := int64(float64(fs.MRR) / (1.0 + fs.MonthlyGrowthRate))
		if prevMRR < 100000 {
			highlights = append(highlights, MonthlyHighlight{
				Type:    "win",
				Message: "Broke $100k MRR milestone! 🎉",
				Icon:    "💰",
			})
		}
	}

	if fs.CustomerChurnRate < 0.05 && fs.Customers > 5 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "Churn rate below 5% - excellent retention!",
			Icon:    "🎯",
		})
	}

	if fs.MonthlyGrowthRate > 0.20 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: fmt.Sprintf("Strong growth: %.0f%% MoM!", fs.MonthlyGrowthRate*100),
			Icon:    "📈",
		})
	}

	ltvCac := fs.CalculateLTVToCAC()
	if ltvCac >= 3.0 && fs.Customers > 10 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "LTV:CAC ratio is healthy (>3:1)",
			Icon:    "✨",
		})
	}

	if fs.ProductMaturity >= 0.8 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "Product is highly mature - low churn expected",
			Icon:    "🚀",
		})
	}

	ruleOf40 := fs.CalculateRuleOf40()
	if ruleOf40 >= 40 && fs.MRR > 50000 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "Rule of 40 achieved - excellent balance!",
			Icon:    "💎",
		})
	}

	// CONCERNS
	if fs.CashRunwayMonths <= 3 && fs.CashRunwayMonths > 0 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: fmt.Sprintf("Low runway: %d months left!", fs.CashRunwayMonths),
			Icon:    "⚠️",
		})
	}

	if fs.CustomerChurnRate > 0.20 && fs.Customers > 5 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: fmt.Sprintf("High churn rate: %.0f%%/month", fs.CustomerChurnRate*100),
			Icon:    "📉",
		})
	}

	if fs.MonthlyGrowthRate < 0 && fs.Turn > 3 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "Revenue is declining - need to fix growth!",
			Icon:    "🔴",
		})
	}

	if fs.Customers == 0 && fs.Turn > 3 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "No customers yet - focus on acquisition!",
			Icon:    "⚡",
		})
	}

	burnMultiple := fs.CalculateBurnMultiple()
	if burnMultiple > 2.0 && burnMultiple < 999 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "High burn multiple - spending too much per $ of growth",
			Icon:    "💸",
		})
	}

	if ltvCac < 1.0 && ltvCac > 0 && fs.Customers > 10 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "LTV:CAC ratio < 1 - losing money on customers!",
			Icon:    "⛔",
		})
	}

	// Prioritize: show max 3 of each type
	wins := []MonthlyHighlight{}
	concerns := []MonthlyHighlight{}

	for _, h := range highlights {
		if h.Type == "win" {
			wins = append(wins, h)
		} else {
			concerns = append(concerns, h)
		}
	}

	result := []MonthlyHighlight{}

	// Add top 3 wins
	for i := 0; i < len(wins) && i < 3; i++ {
		result = append(result, wins[i])
	}

	// Add top 3 concerns
	for i := 0; i < len(concerns) && i < 3; i++ {
		result = append(result, concerns[i])
	}

	return result
}


func (fs *FounderState) GetBoardGuidance() []string {
	var guidance []string

	if len(fs.BoardMembers) == 0 {
		return guidance
	}

	// Check for chairman first - chairman provides enhanced benefits
	chairman := fs.GetChairman()
	if chairman != nil && chairman.IsActive {
		// Chairman provides guidance 60% of the time (vs 30% for regular advisors)
		if rand.Float64() < 0.6 {
			impactMultiplier := 2.0 // Chairman has 2x impact

			switch chairman.Expertise {
			case "sales":
				// Sales expertise helps with customer acquisition
				boost := int64(float64(fs.MRR) * (0.02 + rand.Float64()*0.03) * impactMultiplier) // 4-10% boost (2x)
				if boost > 0 {
					fs.MRR += boost
					fs.DirectMRR += boost
					guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman - Sales) introduced you to major customers (+$%s MRR)",
						chairman.Name, formatCurrency(boost)))
				}
			case "product":
				// Product expertise improves product maturity
				if fs.ProductMaturity < 1.0 {
					improvement := (0.02 + rand.Float64()*0.03) * impactMultiplier // 4-10% improvement (2x)
					fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+improvement)
					guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman - Product) provided strategic product guidance (%.0f%% maturity gained)",
						chairman.Name, improvement*100))
				}
			case "fundraising":
				// Fundraising expertise improves future round terms
				if len(fs.FundingRounds) < 3 {
					guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman - Fundraising) is actively opening doors to top-tier investors",
						chairman.Name))
				}
			case "operations":
				// Operations expertise reduces costs
				if fs.MonthlyTeamCost > 50000 {
					savings := int64(float64(fs.MonthlyTeamCost) * (0.01 + rand.Float64()*0.02) * impactMultiplier) // 2-6% savings (2x)
					fs.Cash += savings
					guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman - Operations) identified significant cost savings (+$%s this month)",
						chairman.Name, formatCurrency(savings)))
				}
			case "strategy":
				// Strategy expertise helps avoid bad decisions
				if fs.CustomerChurnRate > 0.15 {
					reduction := (0.01 + rand.Float64()*0.02) * impactMultiplier // 2-6% churn reduction (2x)
					fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-reduction)
					guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman - Strategy) provided strategic guidance to reduce churn (%.0f%% improvement)",
						chairman.Name, reduction*100))
				}
			}

			// Chairman also provides investor relations benefit
			if fs.BoardPressure > 0 {
				pressureReduction := 5 + rand.Intn(11) // 5-15 point reduction
				fs.BoardPressure -= pressureReduction
				if fs.BoardPressure < 0 {
					fs.BoardPressure = 0
				}
				guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman) improved investor relations (board pressure reduced by %d points)",
					chairman.Name, pressureReduction))
			}
		}

		// Chairman represents company at events (saves founder time, unlocks opportunities)
		if rand.Float64() < 0.3 {
			// 30% chance chairman attends event on your behalf
			opportunityTypes := []string{"partnership", "customer", "fundraising"}
			opportunityType := opportunityTypes[rand.Intn(len(opportunityTypes))]
			
			switch opportunityType {
			case "partnership":
				guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman) represented company at industry conference - opened partnership discussions",
					chairman.Name))
			case "customer":
				guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman) spoke at conference on your behalf - generated customer leads",
					chairman.Name))
			case "fundraising":
				guidance = append(guidance, fmt.Sprintf("👔 %s (Chairman) networked at investor event - warming up potential investors",
					chairman.Name))
			}
		}
	}

	// Regular board members provide value based on their expertise (30% chance)
	for _, member := range fs.BoardMembers {
		if !member.IsActive || member.IsChairman {
			continue // Skip chairman (already handled above)
		}

		// 30% chance per month a board member provides useful guidance
		if rand.Float64() < 0.3 {
			switch member.Expertise {
			case "sales":
				// Sales expertise helps with customer acquisition
				boost := int64(float64(fs.MRR) * (0.02 + rand.Float64()*0.03)) // 2-5% boost
				if boost > 0 {
					guidance = append(guidance, fmt.Sprintf("📊 %s (Sales Advisor) introduced you to potential customers (+$%s MRR opportunity)",
						member.Name, formatCurrency(boost)))
					// Could apply boost here or make it an opportunity
				}
			case "product":
				// Product expertise improves product maturity
				if fs.ProductMaturity < 1.0 {
					improvement := 0.02 + rand.Float64()*0.03 // 2-5% improvement
					fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+improvement)
					guidance = append(guidance, fmt.Sprintf("🎯 %s (Product Advisor) helped improve product (%.0f%% maturity gained)",
						member.Name, improvement*100))
				}
			case "fundraising":
				// Fundraising expertise improves future round terms
				if len(fs.FundingRounds) < 3 {
					guidance = append(guidance, fmt.Sprintf("💰 %s (Fundraising Advisor) is warming up investors for your next round",
						member.Name))
				}
			case "operations":
				// Operations expertise reduces costs
				if fs.MonthlyTeamCost > 50000 {
					savings := int64(float64(fs.MonthlyTeamCost) * (0.01 + rand.Float64()*0.02)) // 1-3% savings
					fs.Cash += savings
					guidance = append(guidance, fmt.Sprintf("⚙️  %s (Operations Advisor) identified cost savings (+$%s this month)",
						member.Name, formatCurrency(savings)))
				}
			case "strategy":
				// Strategy expertise helps avoid bad decisions
				if fs.CustomerChurnRate > 0.15 {
					reduction := 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
					fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-reduction)
					guidance = append(guidance, fmt.Sprintf("🎓 %s (Strategy Advisor) helped reduce churn (%.0f%% improvement)",
						member.Name, reduction*100))
				}
			}
		}
	}

	return guidance
}


func (fs *FounderState) UpdateBoardSentiment() {
	if len(fs.FundingRounds) == 0 {
		fs.BoardSentiment = ""
		fs.BoardPressure = 0
		return
	}

	// Calculate performance score (0-100)
	score := 50.0 // Start neutral

	// Growth is good
	if fs.MonthlyGrowthRate > 0.15 {
		score += 20
	} else if fs.MonthlyGrowthRate > 0.05 {
		score += 10
	} else if fs.MonthlyGrowthRate < 0 {
		score -= 20
	}

	// Runway matters
	if fs.CashRunwayMonths <= 3 {
		score -= 25
	} else if fs.CashRunwayMonths > 12 {
		score += 10
	}

	// Profitability is good
	netIncome := fs.MRR - fs.MonthlyTeamCost - fs.MonthlyComputeCost - fs.MonthlyODCCost
	if netIncome > 0 {
		score += 15
	}

	// Churn matters
	if fs.CustomerChurnRate < 0.05 {
		score += 10
	} else if fs.CustomerChurnRate > 0.20 {
		score -= 15
	}

	// Set sentiment
	if score >= 75 {
		fs.BoardSentiment = "happy"
		fs.BoardPressure = 10
	} else if score >= 60 {
		fs.BoardSentiment = "pleased"
		fs.BoardPressure = 25
	} else if score >= 40 {
		fs.BoardSentiment = "neutral"
		fs.BoardPressure = 50
	} else if score >= 25 {
		fs.BoardSentiment = "concerned"
		fs.BoardPressure = 75
	} else {
		fs.BoardSentiment = "angry"
		fs.BoardPressure = 95
	}
}


func (fs *FounderState) GenerateStrategicOpportunity() *StrategicOpportunity {
	// Don't generate if one already pending
	if fs.PendingOpportunity != nil {
		return nil
	}

	// 15% chance per month (after month 3)
	if fs.Turn < 3 || rand.Float64() > 0.15 {
		return nil
	}

	opportunityTypes := []string{"press", "enterprise_pilot", "conference", "talent", "competitor_distress"}

	// Add bridge round opportunity if running low on cash and have raised before
	if fs.CashRunwayMonths <= 6 && len(fs.FundingRounds) > 0 {
		opportunityTypes = append(opportunityTypes, "bridge_round")
	}

	oppType := opportunityTypes[rand.Intn(len(opportunityTypes))]

	var opp StrategicOpportunity

	switch oppType {
	case "press":
		opp = StrategicOpportunity{
			Type:        "press",
			Title:       "📰 TechCrunch Feature Opportunity",
			Description: "TechCrunch wants to write a feature story about your company. This could significantly boost brand awareness and inbound leads.",
			Cost:        10000 + rand.Int63n(15000), // $10-25k PR prep
			Benefit:     fmt.Sprintf("+%d potential customers over next 3 months, +15%% brand awareness", 5+rand.Intn(10)),
			Risk:        "Requires founder time (1 week) and PR prep costs",
			ExpiresIn:   1,
		}

	case "enterprise_pilot":
		dealSize := 50000 + rand.Int63n(150000)
		opp = StrategicOpportunity{
			Type:        "enterprise_pilot",
			Title:       "🏢 Enterprise Pilot Program",
			Description: "Fortune 500 company wants to pilot your product. High revenue potential but requires dedicated engineering resources.",
			Cost:        0, // No upfront cost, but requires team time
			Benefit:     fmt.Sprintf("$%s annual contract if successful (80%% chance), reference customer for enterprise sales", formatCurrency(dealSize)),
			Risk:        "Requires 2 engineers for 3 months, may slow product development",
			ExpiresIn:   1,
		}

	case "bridge_round":
		amount := 200000 + rand.Int63n(500000)
		equity := 3.0 + rand.Float64()*5.0
		opp = StrategicOpportunity{
			Type:        "bridge_round",
			Title:       "💰 Bridge Round Opportunity",
			Description: "Existing investor offers bridge financing at favorable terms. Quick capital to extend runway.",
			Cost:        0,
			Benefit:     fmt.Sprintf("$%s at %.1f%% equity (better terms than raising a full round)", formatCurrency(amount), equity),
			Risk:        "Additional dilution, may signal to market that you're struggling",
			ExpiresIn:   2,
		}

	case "conference":
		opp = StrategicOpportunity{
			Type:        "conference",
			Title:       "🎤 Conference Speaking Slot",
			Description: "Invited to speak at major industry conference. Great for leads and recruiting, but takes founder time.",
			Cost:        5000 + rand.Int63n(10000), // Travel + booth costs
			Benefit:     fmt.Sprintf("+%d qualified leads, improved recruiting pipeline, industry credibility", 10+rand.Intn(20)),
			Risk:        "Founder unavailable for 1 week, may not convert leads immediately",
			ExpiresIn:   2,
		}

	case "talent":
		opp = StrategicOpportunity{
			Type:        "talent",
			Title:       "⭐ Star Engineer Available",
			Description: "Senior engineer from Google/Meta is interested in joining. Exceptional talent but expensive and expects senior role.",
			Cost:        200000, // $200k/year salary
			Benefit:     "Accelerates product development 2x, attracts other top talent, improved technical credibility",
			Risk:        "High salary, may create team dynamics issues if not managed well",
			ExpiresIn:   1,
		}

	case "competitor_distress":
		opp = StrategicOpportunity{
			Type:        "competitor_distress",
			Title:       "🎯 Competitor in Distress",
			Description: "Main competitor is struggling (layoffs, negative press). Perfect time to steal their customers or acquire them cheaply.",
			Cost:        50000 + rand.Int63n(150000),
			Benefit:     fmt.Sprintf("+%d customers (from their base), eliminate key competitor", 15+rand.Intn(25)),
			Risk:        "May inherit technical debt or unhappy customers",
			ExpiresIn:   2,
		}
	}

	return &opp
}


// ProcessMonth runs all monthly calculations
func (fs *FounderState) ProcessMonth() []string {
	return fs.ProcessMonthWithBaseline(fs.MRR)
}


func (fs *FounderState) ProcessMonthWithBaseline(baselineMRR int64) []string {
	var messages []string
	fs.Turn++

	// Update employee vesting
	fs.UpdateEmployeeVesting()

	// Ensure MRR is in sync before processing
	fs.syncMRR()

	// 1. Process revenue growth
	oldMRR := baselineMRR // Use the baseline MRR from start of turn, not current MRR

	// Engineer impact on product (reduces churn, increases sales)
	// CTO counts as 3x engineers
	engImpact := 1.0
	for _, eng := range fs.Team.Engineers {
		engImpact += (eng.Impact * 0.05) // Each engineer adds ~5% product improvement
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCTO {
			engImpact += (exec.Impact * 0.05) // CTO has 3x impact already built into Impact field
		}
	}
	fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+(0.02*engImpact))

	// Sales team impact on growth
	// CGO counts as 3x sales reps
	salesImpact := 1.0
	for _, sales := range fs.Team.Sales {
		salesImpact += (sales.Impact * 0.1) // Each sales rep adds ~10% to close rate
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCGO {
			salesImpact += (exec.Impact * 0.1) // CGO has 3x impact already built into Impact field
		}
	}

	// Marketing impact (residual from spend)
	baseGrowth := fs.MonthlyGrowthRate
	actualGrowth := baseGrowth * salesImpact * engImpact

	// Apply growth
	// Calculate growth first, then apply proportionally
	newRevenue := int64(float64(fs.MRR) * actualGrowth)

	// Apply growth proportionally to direct and affiliate MRR
	if fs.MRR > 0 {
		directRatio := float64(fs.DirectMRR) / float64(fs.MRR)
		affiliateRatio := float64(fs.AffiliateMRR) / float64(fs.MRR)
		fs.DirectMRR += int64(float64(newRevenue) * directRatio)
		fs.AffiliateMRR += int64(float64(newRevenue) * affiliateRatio)
	} else {
		// If no MRR yet, growth goes to direct
		fs.DirectMRR += newRevenue
	}

	// Sync MRR from DirectMRR + AffiliateMRR
	fs.syncMRR()

	// 2. Process churn (only if we have customers)
	var lostCustomers int
	var lostDirectCustomers int
	var lostAffiliateCustomers int
	var actualChurn float64

	// Always recalculate base churn based on current product maturity
	// Lower maturity = higher churn
	// Formula: baseChurn = (1.0 - ProductMaturity) * 0.65 + 0.05
	// This ensures churn decreases as product matures
	baseChurnFromMaturity := (1.0-fs.ProductMaturity)*0.65 + 0.05
	// Cap at reasonable range (5% minimum, 70% maximum)
	baseChurnFromMaturity = math.Max(0.05, math.Min(0.70, baseChurnFromMaturity))
	baseChurn := baseChurnFromMaturity

	// CS team reduces churn
	// COO counts as 3x CS reps
	csImpact := 0.0
	for _, cs := range fs.Team.CustomerSuccess {
		csImpact += (cs.Impact * 0.02) // Each CS rep reduces churn by ~2%
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCOO {
			csImpact += (exec.Impact * 0.02) // COO has 3x impact already built into Impact field
		}
	}

	// Calculate effective churn rate (after CS team impact)
	actualChurn = math.Max(0.01, baseChurn-csImpact)

	// Update displayed churn rate using helper function
	fs.RecalculateChurnRate()

	if fs.Customers > 0 {
		// Get active customers by source for proportional churn
		activeDirectCustomers := []Customer{}
		activeAffiliateCustomers := []Customer{}
		for _, c := range fs.CustomerList {
			if c.IsActive {
				if c.Source == "direct" || c.Source == "partnership" || c.Source == "market" {
					activeDirectCustomers = append(activeDirectCustomers, c)
				} else if c.Source == "affiliate" {
					activeAffiliateCustomers = append(activeAffiliateCustomers, c)
				}
			}
		}

		// Calculate how many to churn from each source
		lostDirectCustomers = int(float64(len(activeDirectCustomers)) * actualChurn)
		lostAffiliateCustomers = int(float64(len(activeAffiliateCustomers)) * actualChurn)

		// Mark customers as churned (randomly select from active customers)
		rand.Shuffle(len(activeDirectCustomers), func(i, j int) {
			activeDirectCustomers[i], activeDirectCustomers[j] = activeDirectCustomers[j], activeDirectCustomers[i]
		})
		rand.Shuffle(len(activeAffiliateCustomers), func(i, j int) {
			activeAffiliateCustomers[i], activeAffiliateCustomers[j] = activeAffiliateCustomers[j], activeAffiliateCustomers[i]
		})

		// Churn direct customers
		for i := 0; i < lostDirectCustomers && i < len(activeDirectCustomers); i++ {
			customer := activeDirectCustomers[i]
			fs.churnCustomer(customer.ID)
			fs.DirectMRR -= customer.DealSize
		}

		// Churn affiliate customers
		for i := 0; i < lostAffiliateCustomers && i < len(activeAffiliateCustomers); i++ {
			customer := activeAffiliateCustomers[i]
			fs.churnCustomer(customer.ID)
			fs.AffiliateMRR -= customer.DealSize
		}

		lostCustomers = lostDirectCustomers + lostAffiliateCustomers
		fs.Customers -= lostCustomers
		fs.DirectCustomers -= lostDirectCustomers
		fs.AffiliateCustomers -= lostAffiliateCustomers
	}

	// Clamp values to prevent negatives
	if fs.DirectMRR < 0 {
		fs.DirectMRR = 0
	}
	if fs.AffiliateMRR < 0 {
		fs.AffiliateMRR = 0
	}
	if fs.Customers < 0 {
		fs.Customers = 0
	}
	if fs.DirectCustomers < 0 {
		fs.DirectCustomers = 0
	}
	if fs.AffiliateCustomers < 0 {
		fs.AffiliateCustomers = 0
	}

	// Ensure MRR is in sync with DirectMRR + AffiliateMRR
	fs.syncMRR()

	// Recalculate average deal size after churn
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}
	// If no customers, keep AvgDealSize from template (don't reset to 0)

	// 3. Process MRR cash flow
	// MRR converts to cash after deductions:
	// - Taxes: ~15-25% depending on jurisdiction (use 20% average)
	// - Processing fees: ~3-5% for payment processing (use 3%)
	// - Company overhead: ~5-10% for operations (use 5%)
	// - Savings/buffer: ~5% for reserves
	taxRate := 0.20           // 20% taxes
	processingFeeRate := 0.03 // 3% payment processing
	overheadRate := 0.05      // 5% company overhead
	savingsRate := 0.05       // 5% savings/reserves

	totalDeductionRate := taxRate + processingFeeRate + overheadRate + savingsRate // 33% total deductions
	netMRRToCash := int64(float64(fs.MRR) * (1.0 - totalDeductionRate))
	fs.Cash += netMRRToCash

	// 4. Calculate costs
	totalCost := fs.MonthlyTeamCost + (int64(fs.Team.TotalEmployees) * 2000) // +$2k overhead per employee

	// Calculate infrastructure costs (compute + ODC)
	fs.CalculateInfrastructureCosts()
	totalInfrastructureCost := fs.MonthlyComputeCost + fs.MonthlyODCCost

	fs.Cash -= totalCost
	fs.Cash -= totalInfrastructureCost

	netIncome := netMRRToCash - totalCost - totalInfrastructureCost

	// 5. Update runway
	fs.CalculateRunway()

	// Store messages for churn, cash flow, etc. (will add MRR comparison later)
	if lostCustomers > 0 {
		messages = append(messages, fmt.Sprintf("📉 Lost %d customers to churn (%.1f%% churn rate)", lostCustomers, actualChurn*100))
	}

	if netIncome > 0 {
		messages = append(messages, fmt.Sprintf("✅ Positive cash flow: $%s/month", formatCurrency(netIncome)))
	} else {
		messages = append(messages, fmt.Sprintf("💸 Burn rate: $%s/month", formatCurrency(-netIncome)))
	}

	if fs.ProductMaturity >= 1.0 {
		messages = append(messages, "🎉 Product has reached full maturity!")
	}

	// 7. Process advanced features (affiliates, partnerships, etc.)
	// These will add more MRR, so we'll compare baseline to final MRR after all processing
	partnershipMsgs := fs.UpdatePartnerships()
	messages = append(messages, partnershipMsgs...)

	affiliateMsgs := fs.UpdateAffiliateProgram()
	messages = append(messages, affiliateMsgs...)

	referralMsgs := fs.UpdateReferralProgram()
	messages = append(messages, referralMsgs...)

	competitorMsgs := fs.UpdateCompetitors()
	messages = append(messages, competitorMsgs...)

	marketMsgs := fs.UpdateGlobalMarkets()
	messages = append(messages, marketMsgs...)

	// 8. Spawn new competitors randomly
	if newComp := fs.SpawnCompetitor(); newComp != nil {
		messages = append(messages, fmt.Sprintf("🚨 NEW COMPETITOR: %s entered the market! Threat: %s, Market Share: %.1f%%",
			newComp.Name, newComp.Threat, newComp.MarketShare*100))
	}

	// 9. Process random events
	eventMsgs := fs.ProcessRandomEvents()
	messages = append(messages, eventMsgs...)

	// 10. Spawn new random events (5% chance each month)
	if rand.Float64() < 0.05 {
		if event := fs.SpawnRandomEvent(); event != nil {
			messages = append(messages, fmt.Sprintf("⚡ EVENT: %s - %s", event.Title, event.Description))
		}
	}

	// Final sync to ensure MRR is always correct (after all processing)
	fs.syncMRR()

	// Recalculate infrastructure costs after all customer changes (affiliates, partnerships, etc.)
	fs.CalculateInfrastructureCosts()

	// Show infrastructure costs if significant
	if fs.MonthlyComputeCost > 0 || fs.MonthlyODCCost > 0 {
		messages = append(messages, fmt.Sprintf("💻 Infrastructure: Compute $%s/mo, ODC $%s/mo",
			formatCurrency(fs.MonthlyComputeCost), formatCurrency(fs.MonthlyODCCost)))
	}

	// 11. Update growth rate for display on next month's dashboard (AFTER all processing)
	if oldMRR > 0 {
		fs.MonthlyGrowthRate = float64(fs.MRR-oldMRR) / float64(oldMRR)
	} else if fs.MRR > 0 {
		// First customers! Set initial growth rate
		fs.MonthlyGrowthRate = 0.10 // Start with 10% base growth
	}

	// Update board sentiment if raised funding
	fs.UpdateBoardSentiment()

	// Get board guidance
	boardGuidance := fs.GetBoardGuidance()
	messages = append(messages, boardGuidance...)

	// Generate strategic opportunity (15% chance, or 25% with good board)
	if fs.PendingOpportunity == nil {
		newOpp := fs.GenerateStrategicOpportunity()
		if newOpp != nil {
			fs.PendingOpportunity = newOpp
		}
	} else {
		// Decrement expiration timer
		fs.PendingOpportunity.ExpiresIn--
		if fs.PendingOpportunity.ExpiresIn <= 0 {
			messages = append(messages, fmt.Sprintf("⏰ Opportunity expired: %s", fs.PendingOpportunity.Title))
			fs.PendingOpportunity = nil
		}
	}

	// 12. Generate MRR comparison message (AFTER all customer additions)
	if fs.MRR > 0 && oldMRR == 0 {
		messages = append(messages, fmt.Sprintf("🎉 FIRST REVENUE! MRR: $%s", formatCurrency(fs.MRR)))
	} else if fs.MRR > oldMRR && oldMRR > 0 {
		pctGrowth := ((float64(fs.MRR) - float64(oldMRR)) / float64(oldMRR)) * 100
		messages = append(messages, fmt.Sprintf("💰 MRR grew %.1f%% to $%s", pctGrowth, formatCurrency(fs.MRR)))
	} else if fs.MRR < oldMRR && oldMRR > 0 {
		pctDecline := ((float64(oldMRR) - float64(fs.MRR)) / float64(oldMRR)) * 100
		messages = append(messages, fmt.Sprintf("⚠️  MRR declined %.1f%% to $%s", pctDecline, formatCurrency(fs.MRR)))
	} else if fs.MRR == 0 && fs.Turn > 3 {
		messages = append(messages, "⚠️  Still no revenue! Hire sales or spend on marketing!")
	}

	// Recalculate average deal size if we have customers
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}
	// If no customers, keep AvgDealSize from template (don't reset to 0)

	return messages
}