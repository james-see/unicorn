package game

import (
	"math/rand"
)

func (gs *GameState) InitializeAIPlayers() {
	// Initialize LP commitments for AI players
	lpCommittedCapital, capitalCallSchedule := initializeLPCommitments(gs.Difficulty.StartingCash, gs.Difficulty.MaxTurns)

	// Randomly select 3-5 AI players from the pool
	allAIPlayers := []AIPlayer{
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
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
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
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
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
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Alex Rodriguez",
			Firm:          "Tiger Global",
			Strategy:      "aggressive",
			RiskTolerance: 0.85, // Very aggressive growth investor
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Jessica Park",
			Firm:          "Y Combinator",
			Strategy:      "early_stage",
			RiskTolerance: 0.6, // Early-stage specialist, moderate risk
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Raj Patel",
			Firm:          "SoftBank Vision Fund",
			Strategy:      "mega_fund",
			RiskTolerance: 0.75, // Mega-fund, high valuations
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "David Kim",
			Firm:          "Andreessen Horowitz",
			Strategy:      "balanced",
			RiskTolerance: 0.55, // Balanced approach, tech-focused
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Lisa Thompson",
			Firm:          "Benchmark Capital",
			Strategy:      "conservative",
			RiskTolerance: 0.35, // Conservative, quality-focused
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Michael Chen",
			Firm:          "First Round Capital",
			Strategy:      "seed_focused",
			RiskTolerance: 0.65, // Seed-stage specialist, higher risk tolerance
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Emily Rodriguez",
			Firm:          "Greylock Partners",
			Strategy:      "enterprise_focused",
			RiskTolerance: 0.45, // Enterprise SaaS focus, moderate risk
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "James Wilson",
			Firm:          "Index Ventures",
			Strategy:      "deep_tech",
			RiskTolerance: 0.7, // Deep tech focus, high risk/high reward
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
		{
			Name:          "Sophie Martin",
			Firm:          "Lightspeed Venture Partners",
			Strategy:      "consumer_focused",
			RiskTolerance: 0.6, // Consumer products focus, moderate-high risk
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
				CarryInterestPaid:   0,
				LPCommittedCapital:  lpCommittedCapital,
				LPCalledCapital:     0,
				LastCapitalCallTurn: 0,
				CapitalCallSchedule: capitalCallSchedule,
			},
		},
	}

	// Shuffle and select 3-5 players
	rand.Shuffle(len(allAIPlayers), func(i, j int) {
		allAIPlayers[i], allAIPlayers[j] = allAIPlayers[j], allAIPlayers[i]
	})

	numPlayers := 3 + rand.Intn(3) // 3-5 AI players
	if numPlayers > len(allAIPlayers) {
		numPlayers = len(allAIPlayers)
	}

	gs.AIPlayers = allAIPlayers[:numPlayers]
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
			// Note: RiskScore ranges from 0.5-1.0, GrowthPotential from 0.5-1.0
			shouldInvest := false
			if ai.Strategy == "conservative" {
				// Conservative: prefer lower risk (0.5-0.65) and decent growth
				shouldInvest = startup.RiskScore <= 0.65 && startup.GrowthPotential >= 0.55
			} else if ai.Strategy == "aggressive" {
				// Aggressive: chase high growth, accept higher risk
				shouldInvest = startup.GrowthPotential >= 0.7 || (startup.RiskScore >= 0.7 && startup.GrowthPotential >= 0.6)
			} else if ai.Strategy == "early_stage" {
				// Y Combinator: Focus on early-stage companies with high growth potential
				shouldInvest = startup.Valuation < 800000 && startup.GrowthPotential >= 0.6
			} else if ai.Strategy == "mega_fund" {
				// SoftBank: Invest in larger rounds, focus on scale
				shouldInvest = startup.Valuation >= 500000 && startup.GrowthPotential >= 0.6
			} else if ai.Strategy == "seed_focused" {
				// First Round: Seed-stage focus, lower valuations, high growth potential
				shouldInvest = startup.Valuation < 600000 && startup.GrowthPotential >= 0.6 && startup.RiskScore <= 0.8
			} else if ai.Strategy == "enterprise_focused" {
				// Greylock: Enterprise SaaS focus, prefer SaaS category
				shouldInvest = (startup.Category == "SaaS" || startup.Category == "GovTech") && startup.GrowthPotential >= 0.55
			} else if ai.Strategy == "deep_tech" {
				// Index Ventures: Deep tech focus, high risk tolerance
				shouldInvest = (startup.Category == "DeepTech" || startup.Category == "Hardware" || startup.Category == "CleanTech") && startup.GrowthPotential >= 0.55
			} else if ai.Strategy == "consumer_focused" {
				// Lightspeed: Consumer products focus
				shouldInvest = (startup.Category == "SaaS" || startup.Category == "Advertising" || startup.Category == "Consumer") && startup.GrowthPotential >= 0.55
			} else { // balanced
				// Balanced: moderate approach - will invest in most decent startups
				shouldInvest = startup.GrowthPotential >= 0.5 && startup.RiskScore <= 0.85
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
						ValuationCap:        0,
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
