package game

import (
	"fmt"
	"math/rand"
)


func (gs *GameState) GenerateTermOptions(startup *Startup, amount int64) []InvestmentTerms {
	options := []InvestmentTerms{}

	// Option 1: Preferred Stock (VC standard)
	hasBoardSeat := amount >= 100000 // Board seat for $100k+ investments

	// Check for double_board_seat upgrade - grants 2x voting power
	boardSeatMultiplier := 1
	for _, upgradeID := range gs.PlayerUpgrades {
		if upgradeID == "double_board_seat" {
			boardSeatMultiplier = 2 // Double voting power with upgrade
			break
		}
	}

	options = append(options, InvestmentTerms{
		Type:                "Preferred Stock",
		HasProRataRights:    true,
		HasInfoRights:       true,
		HasBoardSeat:        hasBoardSeat,
		BoardSeatMultiplier: boardSeatMultiplier,
		LiquidationPref:     1.0,
		HasAntiDilution:     true,
		ConversionDiscount:  0.0,
		ValuationCap:        0,
	})

	// Option 2: SAFE (Simple Agreement for Future Equity) - Standard (discount only)
	safeDiscount := 0.20 // Default 20%
	for _, upgradeID := range gs.PlayerUpgrades {
		if upgradeID == "enhanced_safe_discount" {
			safeDiscount = 0.25 // 25% discount
			break
		}
	}

	options = append(options, InvestmentTerms{
		Type:                "SAFE",
		HasProRataRights:    true,
		HasInfoRights:       false,
		HasBoardSeat:        false,
		BoardSeatMultiplier: 1,
		LiquidationPref:     0.0, // No liquidation preference with SAFE
		HasAntiDilution:     false,
		ConversionDiscount:  safeDiscount,
		ValuationCap:        0, // No cap - uses discount only
	})

	// Option 2b: SAFE with Valuation Cap (more investor-friendly)
	// Valuation cap is typically 50-80% of current valuation
	valuationCap := int64(float64(startup.Valuation) * 0.65) // 65% of current valuation
	if valuationCap < 500000 {
		valuationCap = 500000 // Minimum $500k cap
	}

	options = append(options, InvestmentTerms{
		Type:                "SAFE (Capped)",
		HasProRataRights:    true,
		HasInfoRights:       false,
		HasBoardSeat:        false,
		BoardSeatMultiplier: 1,
		LiquidationPref:     0.0,
		HasAntiDilution:     false,
		ConversionDiscount:  safeDiscount,
		ValuationCap:        valuationCap, // Capped at 65% of current valuation
	})

	// Option 3: Common Stock (founder-friendly)
	options = append(options, InvestmentTerms{
		Type:                "Common Stock",
		HasProRataRights:    false,
		HasInfoRights:       false,
		HasBoardSeat:        false,
		BoardSeatMultiplier: 1,
		LiquidationPref:     0.0,
		HasAntiDilution:     false,
		ConversionDiscount:  0.0,
		ValuationCap:        0,
	})

	// Option 4: Preferred Stock with 2x Liquidation Preference (if upgrade unlocked)
	has2xLiquidationPref := false
	for _, upgradeID := range gs.PlayerUpgrades {
		if upgradeID == "liquidation_preference_2x" {
			has2xLiquidationPref = true
			break
		}
	}

	if has2xLiquidationPref {
		options = append(options, InvestmentTerms{
			Type:                "Preferred Stock (2x Liquidation)",
			HasProRataRights:    true,
			HasInfoRights:       true,
			HasBoardSeat:        hasBoardSeat,
			BoardSeatMultiplier: boardSeatMultiplier,
			LiquidationPref:     2.0, // 2x liquidation preference
			HasAntiDilution:     true,
			ConversionDiscount:  0.0,
			ValuationCap:        0,
		})
	}

	return options
}

func (gs *GameState) MakeInvestment(startupIndex int, amount int64) error {
	return gs.MakeInvestmentWithTerms(startupIndex, amount, InvestmentTerms{
		Type:                "Preferred Stock",
		HasProRataRights:    true,
		HasInfoRights:       true,
		HasBoardSeat:        amount >= 100000,
		BoardSeatMultiplier: 1, // Default multiplier
		LiquidationPref:     1.0,
		HasAntiDilution:     true,
	})
}

func (gs *GameState) MakeInvestmentWithTerms(startupIndex int, amount int64, terms InvestmentTerms) error {
	if amount <= 0 {
		return fmt.Errorf("investment amount must be positive")
	}

	if amount > gs.Portfolio.Cash {
		return fmt.Errorf("insufficient funds (have $%d, need $%d)", gs.Portfolio.Cash, amount)
	}

	if startupIndex < 0 || startupIndex >= len(gs.AvailableStartups) {
		return fmt.Errorf("invalid startup index")
	}

	startup := gs.AvailableStartups[startupIndex]

	// Check if already invested in this company
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == startup.Name {
			return fmt.Errorf("you have already invested in %s", startup.Name)
		}
	}

	// Minimum investment is $10,000 (standard VC practice)
	minInvestment := int64(10000)
	if amount < minInvestment {
		return fmt.Errorf("minimum investment is $%d", minInvestment)
	}

	// Maximum investment is 20% of company valuation (standard VC practice)
	// Can be increased to 50% with super_pro_rata upgrade
	maxInvestmentPercent := 0.20
	for _, upgradeID := range gs.PlayerUpgrades {
		if upgradeID == "super_pro_rata" {
			maxInvestmentPercent = 0.50 // 50% max with upgrade
			break
		}
	}
	maxInvestment := int64(float64(startup.Valuation) * maxInvestmentPercent)
	if amount > maxInvestment {
		return fmt.Errorf("maximum investment is $%d (%.0f%% of company valuation: $%d)", maxInvestment, maxInvestmentPercent*100, startup.Valuation)
	}

	// Calculate equity percentage based on investment amount and company valuation
	// Only 20% of company is available for investment in this round
	equityPercent := (float64(amount) / float64(startup.Valuation)) * 100.0

	// Apply Seed Accelerator upgrade - first investment gets 25% equity bonus
	isFirstInvestment := len(gs.Portfolio.Investments) == 0
	hasSeedAccelerator := false
	for _, upgradeID := range gs.PlayerUpgrades {
		if upgradeID == "seed_accelerator" {
			hasSeedAccelerator = true
			break
		}
	}
	if hasSeedAccelerator && isFirstInvestment {
		equityPercent *= 1.25 // 25% bonus
	}

	// Apply SAFE discount if applicable
	if terms.Type == "SAFE" && terms.ConversionDiscount > 0 {
		equityPercent = equityPercent * (1 + terms.ConversionDiscount)
		// Cap at 20% even with discount (since only 20% is available)
		maxEquityPercent := 20.0 * (1 + terms.ConversionDiscount)
		if equityPercent > maxEquityPercent {
			equityPercent = maxEquityPercent
		}
	}

	// Safety cap: equity cannot exceed 20% (or 24% with SAFE discount, or 25% with Seed Accelerator)
	maxEquityPercent := 20.0
	if terms.Type == "SAFE" && terms.ConversionDiscount > 0 {
		maxEquityPercent = 20.0 * (1 + terms.ConversionDiscount)
	}
	// Seed Accelerator can push above 20% (25% bonus on first investment)
	if hasSeedAccelerator && isFirstInvestment {
		maxEquityPercent = 25.0 // Allow up to 25% for first investment with Seed Accelerator
		if terms.Type == "SAFE" && terms.ConversionDiscount > 0 {
			maxEquityPercent = 25.0 * (1 + terms.ConversionDiscount) // SAFE + Seed Accelerator
		}
	}
	if equityPercent > maxEquityPercent {
		equityPercent = maxEquityPercent
	}

	investment := Investment{
		CompanyName:      startup.Name,
		AmountInvested:   amount,
		EquityPercent:    equityPercent,
		InitialEquity:    equityPercent,
		InitialValuation: startup.Valuation,
		CurrentValuation: startup.Valuation,
		MonthsHeld:       0,
		Category:         startup.Category,
		NegativeNewsSent: false,
		Rounds:           []FundingRound{},
		Terms:            terms,
		FollowOnThisTurn: false,
	}

	gs.Portfolio.Investments = append(gs.Portfolio.Investments, investment)
	gs.Portfolio.Cash -= amount
	gs.updateNetWorth()

	return nil
}


func (gs *GameState) GetFollowOnOpportunities() []FollowOnOpportunity {
	opportunities := []FollowOnOpportunity{}

	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Check if player has invested in this company
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == event.CompanyName {
					// Find the startup
					for _, startup := range gs.AvailableStartups {
						if startup.Name == event.CompanyName {
							preMoneyVal := startup.Valuation
							postMoneyVal := preMoneyVal + event.RaiseAmount

							// Calculate min/max investment amounts
							minInvestment := int64(10000) // $10k minimum
							// Maximum investment is 20% of pre-money valuation (standard VC practice)
							maxInvestmentByValuation := int64(float64(preMoneyVal) * 0.20)
							// Use available cash (uninvested money from beginning) + follow-on reserve
							availableCash := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve
							// Maximum is the lower of: 20% of valuation, available cash, or 50% of raise amount
							maxInvestment := maxInvestmentByValuation
							if maxInvestment > availableCash {
								maxInvestment = availableCash
							}
							if maxInvestment > event.RaiseAmount/2 {
								maxInvestment = event.RaiseAmount / 2 // Can't invest more than half the round
							}

							opportunities = append(opportunities, FollowOnOpportunity{
								CompanyName:   event.CompanyName,
								RoundName:     event.RoundName,
								PreMoneyVal:   preMoneyVal,
								PostMoneyVal:  postMoneyVal,
								CurrentEquity: inv.EquityPercent,
								MinInvestment: minInvestment,
								MaxInvestment: maxInvestment,
							})
							break
						}
					}
					break
				}
			}
		}
	}

	return opportunities
}

func (gs *GameState) MakeFollowOnInvestment(companyName string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("investment amount must be positive")
	}

	if amount > gs.Portfolio.Cash+gs.Portfolio.FollowOnReserve {
		return fmt.Errorf("insufficient follow-on funds (have $%d, need $%d)", gs.Portfolio.Cash+gs.Portfolio.FollowOnReserve, amount)
	}

	// Use cash first, then follow-on reserve
	drawnFromCash := amount
	if drawnFromCash > gs.Portfolio.Cash {
		drawnFromCash = gs.Portfolio.Cash
	}
	drawnFromReserve := amount - drawnFromCash

	// Find the funding round event for this turn to get the post-money valuation
	var postMoneyVal int64
	foundRound := false

	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn && event.CompanyName == companyName {
			// Find the company
			for _, startup := range gs.AvailableStartups {
				if startup.Name == companyName {
					preMoneyVal := startup.Valuation
					postMoneyVal = preMoneyVal + event.RaiseAmount
					foundRound = true
					break
				}
			}
			break
		}
	}

	if !foundRound {
		return fmt.Errorf("no funding round happening for %s this turn", companyName)
	}

	// Find the company and calculate max investment (20% of pre-money valuation)
	var preMoneyVal int64
	var foundCompany bool
	for _, startup := range gs.AvailableStartups {
		if startup.Name == companyName {
			preMoneyVal = startup.Valuation
			foundCompany = true
			break
		}
	}

	if !foundCompany {
		return fmt.Errorf("company %s not found", companyName)
	}

	// Maximum follow-on investment is 20% of current pre-money valuation for THIS round
	maxInvestment := int64(float64(preMoneyVal) * 0.20)

	// Find the investment
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			inv := &gs.Portfolio.Investments[i]

			// Check if this follow-on investment exceeds 20% limit for THIS round
			// The 20% limit applies to each round separately, not cumulatively
			if amount > maxInvestment {
				return fmt.Errorf("follow-on investment of $%d exceeds maximum of $%d (20%% of current pre-money valuation: $%d)", amount, maxInvestment, preMoneyVal)
			}

			// Update total amount invested
			inv.AmountInvested += amount

			// Calculate effective investment amount (applying SAFE conversion discount and valuation cap if applicable)
			effectiveAmount := float64(inv.AmountInvested)
			if inv.Terms.Type == "SAFE" || inv.Terms.Type == "SAFE (Capped)" {
				// Apply discount to get effective amount
				if inv.Terms.ConversionDiscount > 0 {
					effectiveAmount = float64(inv.AmountInvested) * (1.0 + inv.Terms.ConversionDiscount)
				}
			} else if inv.Terms.ConversionDiscount > 0 {
				// Non-SAFE with discount (shouldn't happen, but handle it)
				effectiveAmount = float64(inv.AmountInvested) * (1.0 + inv.Terms.ConversionDiscount)
			}

			// Recalculate total equity based on total invested amount and post-money valuation
			// For SAFEs with valuation caps, use the cap if round valuation exceeds cap
			conversionValuation := float64(postMoneyVal)
			if (inv.Terms.Type == "SAFE" || inv.Terms.Type == "SAFE (Capped)") && inv.Terms.ValuationCap > 0 {
				// For SAFEs, conversion happens at the pre-money valuation (before new money comes in)
				// But we calculate equity based on post-money
				// If pre-money exceeds cap, use cap for conversion calculation
				if preMoneyVal > inv.Terms.ValuationCap {
					// Convert at cap valuation, then calculate post-money equity
					// Effective post-money = cap + (postMoneyVal - preMoneyVal) = cap + raiseAmount
					raiseAmount := postMoneyVal - preMoneyVal
					effectivePostMoney := float64(inv.Terms.ValuationCap) + float64(raiseAmount)
					conversionValuation = effectivePostMoney
				}
			}
			
			newEquityPercent := (effectiveAmount / conversionValuation) * 100.0

			// Cap equity at 100% (should never happen, but safety check)
			if newEquityPercent > 100.0 {
				newEquityPercent = 100.0
			}

			inv.EquityPercent = newEquityPercent
			inv.FollowOnThisTurn = true // Mark that follow-on was made this turn

			// Deduct from cash first, then follow-on reserve
			gs.Portfolio.Cash -= drawnFromCash
			gs.Portfolio.FollowOnReserve -= drawnFromReserve
			gs.updateNetWorth()

			return nil
		}
	}

	return fmt.Errorf("you have not invested in %s", companyName)
}

func (gs *GameState) HasFollowOnOpportunities() bool {
	opportunities := gs.GetFollowOnOpportunities()
	return len(opportunities) > 0
}

// GenerateSyndicateOpportunities creates co-investment opportunities with AI investors
// Only generates if player has unlocked syndicate feature (level 2+)
func (gs *GameState) GenerateSyndicateOpportunities(playerLevel int) {
	gs.SyndicateOpportunities = []SyndicateOpportunity{}
	
	// Check if syndicates are unlocked (level 2+)
	if playerLevel < 2 {
		return
	}
	
	// Generate 2-4 syndicate opportunities from available startups
	numOpportunities := 2 + rand.Intn(3) // 2-4 opportunities
	
	// Select random startups that aren't already in player's portfolio
	availableForSyndicate := []int{}
	for i, startup := range gs.AvailableStartups {
		// Check if player already invested
		alreadyInvested := false
		for _, inv := range gs.Portfolio.Investments {
			if inv.CompanyName == startup.Name {
				alreadyInvested = true
				break
			}
		}
		if !alreadyInvested {
			availableForSyndicate = append(availableForSyndicate, i)
		}
	}
	
	// Shuffle and take first N
	if len(availableForSyndicate) > numOpportunities {
		rand.Shuffle(len(availableForSyndicate), func(i, j int) {
			availableForSyndicate[i], availableForSyndicate[j] = availableForSyndicate[j], availableForSyndicate[i]
		})
		availableForSyndicate = availableForSyndicate[:numOpportunities]
	}
	
	// Generate syndicate opportunities
	for _, startupIdx := range availableForSyndicate {
		startup := gs.AvailableStartups[startupIdx]
		
		// Pick a random AI investor to lead
		leadInvestorIdx := rand.Intn(len(gs.AIPlayers))
		leadInvestor := gs.AIPlayers[leadInvestorIdx]
		
		// Calculate round size (typically 1.5-3x company valuation for seed rounds)
		roundMultiplier := 1.5 + rand.Float64()*1.5 // 1.5x to 3x
		totalRoundSize := int64(float64(startup.Valuation) * roundMultiplier)
		
		// Player can invest 20-40% of the round
		playerSharePercent := 0.20 + rand.Float64()*0.20 // 20-40%
		yourMaxShare := int64(float64(totalRoundSize) * playerSharePercent)
		yourMinShare := int64(25000) // $25k minimum
		
		// Cap max share at available cash
		if yourMaxShare > gs.Portfolio.Cash {
			yourMaxShare = gs.Portfolio.Cash
		}
		
		// Generate benefits based on startup characteristics
		benefits := []string{}
		if startup.GrowthPotential > 0.7 {
			benefits = append(benefits, "High-growth opportunity")
		}
		if startup.RiskScore < 0.4 {
			benefits = append(benefits, "Lower risk profile")
		}
		if leadInvestor.Strategy == "aggressive" {
			benefits = append(benefits, "Access to hot deal")
		} else if leadInvestor.Strategy == "conservative" {
			benefits = append(benefits, "Vetted by conservative investor")
		}
		benefits = append(benefits, "Shared due diligence costs")
		benefits = append(benefits, "Network access through lead investor")
		
		// Generate description
		descriptions := []string{
			fmt.Sprintf("%s is leading a syndicate round for %s", leadInvestor.Name, startup.Name),
			fmt.Sprintf("Co-investment opportunity with %s on %s", leadInvestor.Firm, startup.Name),
			fmt.Sprintf("%s invites you to join their deal on %s", leadInvestor.Name, startup.Name),
		}
		description := descriptions[rand.Intn(len(descriptions))]
		
		opportunity := SyndicateOpportunity{
			CompanyName:      startup.Name,
			StartupIndex:      startupIdx,
			LeadInvestor:     leadInvestor.Name,
			LeadInvestorFirm: leadInvestor.Firm,
			TotalRoundSize:   totalRoundSize,
			YourMaxShare:     yourMaxShare,
			YourMinShare:     yourMinShare,
			Valuation:        startup.Valuation,
			Description:      description,
			Benefits:         benefits,
		}
		
		gs.SyndicateOpportunities = append(gs.SyndicateOpportunities, opportunity)
	}
}

// MakeSyndicateInvestment allows player to co-invest via syndicate
func (gs *GameState) MakeSyndicateInvestment(opportunityIndex int, amount int64) error {
	if opportunityIndex < 0 || opportunityIndex >= len(gs.SyndicateOpportunities) {
		return fmt.Errorf("invalid syndicate opportunity index")
	}
	
	opp := gs.SyndicateOpportunities[opportunityIndex]
	
	if amount < opp.YourMinShare {
		return fmt.Errorf("minimum investment is $%d", opp.YourMinShare)
	}
	
	if amount > opp.YourMaxShare {
		return fmt.Errorf("maximum investment is $%d", opp.YourMaxShare)
	}
	
	if amount > gs.Portfolio.Cash {
		return fmt.Errorf("insufficient funds (have $%d, need $%d)", gs.Portfolio.Cash, amount)
	}
	
	// Check if already invested in this company
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == opp.CompanyName {
			return fmt.Errorf("you have already invested in %s", opp.CompanyName)
		}
	}
	
	startup := gs.AvailableStartups[opp.StartupIndex]
	
	// Calculate equity - syndicate rounds typically have better terms
	// Player gets equity based on their share of the round
	equityPercent := (float64(amount) / float64(opp.TotalRoundSize)) * 100.0
	
	// Syndicate bonus: slightly better terms (5% bonus equity)
	equityPercent *= 1.05
	
	// Cap at 20% (standard max)
	if equityPercent > 20.0 {
		equityPercent = 20.0
	}
	
	// Use Preferred Stock terms (standard for syndicates)
	terms := InvestmentTerms{
		Type:                "Preferred Stock",
		HasProRataRights:    true,
		HasInfoRights:       true,
		HasBoardSeat:        amount >= 100000,
		BoardSeatMultiplier: 1,
		LiquidationPref:     1.0,
		HasAntiDilution:     true,
		ConversionDiscount:  0.0,
		ValuationCap:        0,
	}
	
	investment := Investment{
		CompanyName:      startup.Name,
		AmountInvested:   amount,
		EquityPercent:    equityPercent,
		InitialEquity:    equityPercent,
		InitialValuation: startup.Valuation,
		CurrentValuation: startup.Valuation,
		MonthsHeld:       0,
		Category:         startup.Category,
		NegativeNewsSent: false,
		Rounds:           []FundingRound{},
		Terms:            terms,
		FollowOnThisTurn: false,
	}
	
	gs.Portfolio.Investments = append(gs.Portfolio.Investments, investment)
	gs.Portfolio.Cash -= amount
	gs.updateNetWorth()
	
	// Remove this opportunity (can only invest once)
	gs.SyndicateOpportunities = append(
		gs.SyndicateOpportunities[:opportunityIndex],
		gs.SyndicateOpportunities[opportunityIndex+1:]...,
	)
	
	return nil
}