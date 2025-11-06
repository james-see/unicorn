package founder

import (
	"fmt"
	"math/rand"
)


func (fs *FounderState) GenerateTermSheetOptions(roundName string) []TermSheetOption {
	// Use fixed equity percentages per round (industry standard)
	// Then calculate valuations from those percentages
	var baseEquityPercent float64
	var baseRaise int64
	switch roundName {
	case "Seed":
		baseEquityPercent = 12.0 // 12% base equity for seed
		baseRaise = 2000000      // $2M base raise
	case "Series A":
		baseEquityPercent = 22.0 // 22% base equity for Series A
		baseRaise = 10000000     // $10M base raise
	case "Series B":
		baseEquityPercent = 18.0 // 18% base equity for Series B
		baseRaise = 30000000     // $30M base raise
	default:
		return []TermSheetOption{}
	}

	options := []TermSheetOption{}

	// Option 1: Less money, founder-friendly (lower dilution)
	option1Amount := int64(float64(baseRaise) * 0.7)
	option1Equity := baseEquityPercent * 0.85 // 85% of base (lower dilution)
	option1PostVal := int64(float64(option1Amount) / (option1Equity / 100.0))
	option1PreVal := option1PostVal - option1Amount
	options = append(options, TermSheetOption{
		Amount:        option1Amount,
		PostValuation: option1PostVal,
		PreValuation:  option1PreVal,
		Equity:        option1Equity,
		Terms:         "Founder-friendly",
		Description:   "Lower dilution, founder-friendly terms, but less capital",
	})

	// Option 2: Standard terms (balanced)
	option2Amount := baseRaise
	option2Equity := baseEquityPercent // Base equity percentage
	option2PostVal := int64(float64(option2Amount) / (option2Equity / 100.0))
	option2PreVal := option2PostVal - option2Amount
	options = append(options, TermSheetOption{
		Amount:        option2Amount,
		PostValuation: option2PostVal,
		PreValuation:  option2PreVal,
		Equity:        option2Equity,
		Terms:         "Standard",
		Description:   "Fair terms, balanced approach",
	})

	// Option 3: More money, higher dilution
	option3Amount := int64(float64(baseRaise) * 1.4)
	option3Equity := baseEquityPercent * 1.15 // 15% more equity (higher dilution)
	option3PostVal := int64(float64(option3Amount) / (option3Equity / 100.0))
	option3PreVal := option3PostVal - option3Amount
	options = append(options, TermSheetOption{
		Amount:        option3Amount,
		PostValuation: option3PostVal,
		PreValuation:  option3PreVal,
		Equity:        option3Equity,
		Terms:         "Growth-focused",
		Description:   "More capital to scale faster, but higher dilution",
	})

	// Option 4: Maximum money, investor-heavy terms
	option4Amount := int64(float64(baseRaise) * 1.8)
	option4Equity := baseEquityPercent * 1.35 // 35% more equity (investor-heavy)
	option4PostVal := int64(float64(option4Amount) / (option4Equity / 100.0))
	option4PreVal := option4PostVal - option4Amount
	options = append(options, TermSheetOption{
		Amount:        option4Amount,
		PostValuation: option4PostVal,
		PreValuation:  option4PreVal,
		Equity:        option4Equity,
		Terms:         "Investor-heavy",
		Description:   "Maximum capital, but significant dilution and investor control",
	})

	return options
}


func GenerateInvestorNames(roundName string, amount int64) []string {
	var investors []string

	// Angel/Pre-Seed: individual angels
	angelInvestors := []string{
		"Naval Ravikant", "Balaji Srinivasan", "Jason Calacanis", "David Sacks",
		"Elad Gil", "Lachy Groom", "Sahil Bloom", "Anthony Pompliano",
		"Cyan Banister", "Alexis Ohanian", "Arlan Hamilton", "Kevin Hale",
	}

	// Seed: Mix of angels and micro VCs
	seedFirms := []string{
		"Y Combinator", "Sequoia Scout", "First Round", "SV Angel",
		"Hustle Fund", "Khosla Ventures", "Initialized Capital", "Haystack",
		"Precursor Ventures", "Boost VC", "Felicis Ventures", "Homebrew",
	}

	// Series A: Traditional VCs
	seriesAFirms := []string{
		"Sequoia Capital", "Andreessen Horowitz", "Accel Partners", "Benchmark",
		"Greylock Partners", "Kleiner Perkins", "Lightspeed Venture Partners",
		"Index Ventures", "General Catalyst", "NEA", "GGV Capital", "Bessemer",
	}

	// Series B+: Growth firms and strategics
	growthFirms := []string{
		"Tiger Global", "Insight Partners", "Coatue Management", "DST Global",
		"SoftBank Vision Fund", "General Atlantic", "Thrive Capital", "IVP",
		"Ribbit Capital", "T. Rowe Price", "Fidelity Investments", "BlackRock",
	}

	// Family offices and strategics
	familyOffices := []string{
		"Founders Fund", "Bezos Expeditions", "Schmidt Futures", "Emerson Collective",
		"Zuckerberg Family Office", "Gates Ventures", "Cuban Companies", "Thiel Capital",
	}

	switch roundName {
	case "Angel", "Pre-Seed":
		// 2-4 angels
		count := 2 + rand.Intn(3)
		for i := 0; i < count && i < len(angelInvestors); i++ {
			investors = append(investors, angelInvestors[rand.Intn(len(angelInvestors))])
		}

	case "Seed":
		// Lead VC + 1-2 angels or micro VCs
		investors = append(investors, seedFirms[rand.Intn(len(seedFirms))])
		if rand.Float64() > 0.5 {
			investors = append(investors, angelInvestors[rand.Intn(len(angelInvestors))])
		}
		if amount > 2000000 { // Larger seed rounds have more investors
			investors = append(investors, seedFirms[rand.Intn(len(seedFirms))])
		}

	case "Series A":
		// Lead VC + co-investors
		investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
		if amount > 10000000 {
			investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
		}
		// Sometimes strategic or family office
		if rand.Float64() > 0.7 {
			investors = append(investors, familyOffices[rand.Intn(len(familyOffices))])
		}

	case "Series B", "Series C", "Series D":
		// Growth firms + existing investors
		investors = append(investors, growthFirms[rand.Intn(len(growthFirms))])
		investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
		if amount > 50000000 {
			investors = append(investors, growthFirms[rand.Intn(len(growthFirms))])
		}
		if rand.Float64() > 0.6 {
			investors = append(investors, familyOffices[rand.Intn(len(familyOffices))])
		}

	default:
		// Generic round - mix it up
		investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
	}

	return investors
}


func (fs *FounderState) RaiseFundingWithTerms(roundName string, option TermSheetOption) (success bool) {
	// Apply Better Terms upgrade (-5% equity given away)
	equityToGive := option.Equity
	hasBetterTerms := false
	for _, upgradeID := range fs.PlayerUpgrades {
		if upgradeID == "better_terms" {
			hasBetterTerms = true
			break
		}
	}
	if hasBetterTerms {
		equityToGive = option.Equity * 0.95 // 5% reduction
	}
	
	fs.Cash += option.Amount
	fs.EquityGivenAway += equityToGive

	// Generate investor names for this round
	investors := GenerateInvestorNames(roundName, option.Amount)

	round := FundingRound{
		RoundName:   roundName,
		Amount:      option.Amount,
		Valuation:   option.PreValuation,
		EquityGiven: equityToGive, // Use reduced equity if Better Terms upgrade is active
		Month:       fs.Turn,
		Terms:       option.Terms,
		Investors:   investors,
	}
	fs.FundingRounds = append(fs.FundingRounds, round)

	// Add investors to cap table (split equity among them)
	equityPerInvestor := equityToGive / float64(len(investors))
	for _, investor := range investors {
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         investor,
			Type:         "investor",
			Equity:       equityPerInvestor,
			MonthGranted: fs.Turn,
		})
	}

	fs.CalculateRunway()

	return true
}


func (fs *FounderState) RaiseFunding(roundName string) (success bool, amount int64, terms string, equityGiven float64) {
	options := fs.GenerateTermSheetOptions(roundName)
	if len(options) == 0 {
		return false, 0, "", 0
	}

	// Use standard (middle) option
	option := options[1]
	success = fs.RaiseFundingWithTerms(roundName, option)
	return success, option.Amount, option.Terms, option.Equity
}


func (fs *FounderState) BuybackEquity(roundName string, equityPercent float64) (*Buyback, error) {
	// Must be profitable
	if fs.MRR <= fs.MonthlyTeamCost {
		return nil, fmt.Errorf("must be profitable to buy back equity")
	}

	// Find the round
	var round *FundingRound
	for i := range fs.FundingRounds {
		if fs.FundingRounds[i].RoundName == roundName {
			round = &fs.FundingRounds[i]
			break
		}
	}

	if round == nil {
		return nil, fmt.Errorf("round not found: %s", roundName)
	}

	// Calculate current valuation
	currentValuation := int64(float64(fs.MRR) * 12 * 10) // 10x ARR

	// Price to buy back equity (at current valuation)
	price := int64(float64(currentValuation) * equityPercent / 100)

	if price > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(price))
	}

	fs.Cash -= price
	fs.EquityGivenAway -= equityPercent

	buyback := Buyback{
		Month:        fs.Turn,
		Investor:     roundName,
		EquityBought: equityPercent,
		PricePaid:    price,
		Valuation:    currentValuation,
	}

	fs.InvestorBuybacks = append(fs.InvestorBuybacks, buyback)
	fs.CalculateRunway()

	return &buyback, nil
}


func (fs *FounderState) AddBoardSeat(reason string) {
	fs.BoardSeats++
	fs.EquityPool -= 2.0 // Each board seat costs 2% from equity pool
	if fs.EquityPool < 0 {
		fs.EquityPool = 0
	}
}


func (fs *FounderState) ExpandEquityPool(percentToAdd float64) {
	fs.EquityPool += percentToAdd
	// Dilution happens automatically because founder equity = 100 - EquityGivenAway - EquityPool
}

// ============================================================================
// RANDOM EVENTS
// ============================================================================
