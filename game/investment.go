package game

import (
	"fmt"
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
	})

	// Option 2: SAFE (Simple Agreement for Future Equity)
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

			// Calculate effective investment amount (applying SAFE conversion discount if applicable)
			effectiveAmount := float64(inv.AmountInvested)
			if inv.Terms.Type == "SAFE" && inv.Terms.ConversionDiscount > 0 {
				// SAFE converts at a discount - your investment gets more equity
				effectiveAmount = float64(inv.AmountInvested) * (1.0 + inv.Terms.ConversionDiscount)
			}

			// Recalculate total equity based on total invested amount and post-money valuation
			// This ensures equity is always calculated correctly relative to current valuation
			newEquityPercent := (effectiveAmount / float64(postMoneyVal)) * 100.0

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