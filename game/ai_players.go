package game

import (
	"math/rand"
)


func (gs *GameState) InitializeAIPlayers() {
	gs.AIPlayers = []AIPlayer{
		{
			Name:          "CARL",
			Firm:          "Sterling & Cooper",
			Strategy:      "conservative",
			RiskTolerance: 0.3,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
			},
		},
		{
			Name:          "Sarah Chen",
			Firm:          "Accel Partners",
			Strategy:      "aggressive",
			RiskTolerance: 0.8,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
			},
		},
		{
			Name:          "Marcus Williams",
			Firm:          "Sequoia Capital",
			Strategy:      "balanced",
			RiskTolerance: 0.5,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
			},
		},
	}
}

func (gs *GameState) AIPlayerMakeInvestments() {
	for i := range gs.AIPlayers {
		ai := &gs.AIPlayers[i]

		// Only invest on turn 1 (initial investment phase)
		if gs.Portfolio.Turn != 1 {
			continue
		}

		// AI investment strategy based on risk tolerance
		targetInvestmentCount := 3 + rand.Intn(4) // Invest in 3-6 companies
		availableCash := ai.Portfolio.Cash

		// Shuffle startups for variety
		startups := make([]Startup, len(gs.AvailableStartups))
		copy(startups, gs.AvailableStartups)
		rand.Shuffle(len(startups), func(i, j int) {
			startups[i], startups[j] = startups[j], startups[i]
		})

		investmentsMade := 0
		for _, startup := range startups {
			if investmentsMade >= targetInvestmentCount {
				break
			}

			// Decision based on risk tolerance and startup metrics
			shouldInvest := false
			if ai.Strategy == "conservative" {
				shouldInvest = startup.RiskScore < 0.4 && startup.GrowthPotential > 0.5
			} else if ai.Strategy == "aggressive" {
				shouldInvest = startup.GrowthPotential > 0.7 || (startup.RiskScore > 0.7 && startup.GrowthPotential > 0.6)
			} else { // balanced
				shouldInvest = startup.GrowthPotential > 0.5 && startup.RiskScore < 0.7
			}

			if shouldInvest {
				// Invest portion of available cash
				investmentAmount := availableCash / int64(targetInvestmentCount-investmentsMade)
				if investmentAmount > availableCash {
					investmentAmount = availableCash
				}

				// Maximum investment is 20% of company valuation (standard VC practice)
				maxInvestment := int64(float64(startup.Valuation) * 0.20)
				if investmentAmount > maxInvestment {
					investmentAmount = maxInvestment
				}

				if investmentAmount > 10000 { // Minimum investment
					// Calculate equity percentage (only 20% of company is available)
					equityPercent := (float64(investmentAmount) / float64(startup.Valuation)) * 100.0

					// Safety cap at 20%
					if equityPercent > 20.0 {
						equityPercent = 20.0
					}

					// AI players get Preferred Stock terms (like the player)
					terms := InvestmentTerms{
						Type:                "Preferred Stock",
						HasProRataRights:    true,
						HasInfoRights:       true,
						HasBoardSeat:        investmentAmount >= 100000, // Board seat for $100k+ investments
						BoardSeatMultiplier: 1,                          // AI players don't get upgrades
						LiquidationPref:     1.0,
						HasAntiDilution:     true,
						ConversionDiscount:  0.0,
					}

					investment := Investment{
						CompanyName:      startup.Name,
						AmountInvested:   investmentAmount,
						EquityPercent:    equityPercent,
						InitialEquity:    equityPercent,
						InitialValuation: startup.Valuation,
						CurrentValuation: startup.Valuation,
						MonthsHeld:       0,
						Category:         startup.Category,
						Rounds:           []FundingRound{},
						Terms:            terms,
						FollowOnThisTurn: false,
					}

					ai.Portfolio.Investments = append(ai.Portfolio.Investments, investment)
					ai.Portfolio.Cash -= investmentAmount
					availableCash -= investmentAmount
					investmentsMade++
				}
			}
		}

		// Update AI net worth
		gs.updateAINetWorth(i)
	}
}

func (gs *GameState) updateAINetWorth(aiIndex int) {
	ai := &gs.AIPlayers[aiIndex]
	netWorth := ai.Portfolio.Cash + ai.Portfolio.FollowOnReserve

	for _, inv := range ai.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		netWorth += value
	}

	ai.Portfolio.NetWorth = netWorth
}

func (gs *GameState) ProcessAITurns() {
	// Update all AI investments with same events as player
	for i := range gs.AIPlayers {
		for j := range gs.AIPlayers[i].Portfolio.Investments {
			inv := &gs.AIPlayers[i].Portfolio.Investments[j]
			inv.MonthsHeld++

			wasAboveInitial := inv.CurrentValuation >= inv.InitialValuation

			// Apply same random events and volatility as player investments
			// Random chance of an event happening (based on difficulty)
			if rand.Float64() < gs.Difficulty.EventFrequency && len(gs.EventPool) > 0 {
				event := gs.EventPool[rand.Intn(len(gs.EventPool))]

				inv.CurrentValuation = int64(float64(inv.CurrentValuation) * event.Change)

				// Prevent negative valuations
				if inv.CurrentValuation < 0 {
					inv.CurrentValuation = 0
				}
			} else {
				// Natural growth/decline (random walk) - volatility based on difficulty
				change := (rand.Float64()*2 - 1) * gs.Difficulty.Volatility
				inv.CurrentValuation = int64(float64(inv.CurrentValuation) * (1 + change))
			}

			// Check if investment just went negative (for consistency, but don't generate news for AI)
			nowBelowInitial := inv.CurrentValuation < inv.InitialValuation
			if wasAboveInitial && nowBelowInitial && !inv.NegativeNewsSent {
				inv.NegativeNewsSent = true
			}
		}

		gs.AIPlayers[i].Portfolio.Turn++
		gs.updateAINetWorth(i)
	}
}