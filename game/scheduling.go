package game

import (
	"math/rand"
)


func (gs *GameState) ScheduleFundingRounds() {
	gs.FundingRoundQueue = []FundingRoundEvent{}

	// Schedule funding rounds with realistic amounts
	for _, startup := range gs.AvailableStartups {
		// Seed round (3-9 months) - raise $2M-$5M
		seedTurn := 3 + rand.Intn(7)
		if seedTurn < gs.Portfolio.MaxTurns {
			seedAmount := int64(2000000 + rand.Intn(3000000)) // $2M-$5M
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Seed",
				ScheduledTurn: seedTurn,
				RaiseAmount:   seedAmount,
			})
		}

		// Series A (12-24 months) - raise $10M-$20M
		seriesATurn := 12 + rand.Intn(13)
		if seriesATurn < gs.Portfolio.MaxTurns {
			seriesAAmount := int64(10000000 + rand.Intn(10000000)) // $10M-$20M
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Series A",
				ScheduledTurn: seriesATurn,
				RaiseAmount:   seriesAAmount,
			})
		}

		// Series B (30-48 months) - raise $30M-$50M
		seriesBTurn := 30 + rand.Intn(19)
		if seriesBTurn < gs.Portfolio.MaxTurns {
			seriesBAmount := int64(30000000 + rand.Intn(20000000)) // $30M-$50M
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Series B",
				ScheduledTurn: seriesBTurn,
				RaiseAmount:   seriesBAmount,
			})
		}

		// Series C (48-60 months) - raise $50M-$100M, only for top performers
		if rand.Float64() < 0.3 { // 30% of companies
			seriesCTurn := 48 + rand.Intn(13)
			if seriesCTurn < gs.Portfolio.MaxTurns {
				seriesCAmount := int64(50000000 + rand.Intn(50000000)) // $50M-$100M
				gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
					CompanyName:   startup.Name,
					RoundName:     "Series C",
					ScheduledTurn: seriesCTurn,
					RaiseAmount:   seriesCAmount,
				})
			}
		}

		// 20% chance of a down round occurring (usually Series A or B)
		if rand.Float64() < 0.2 {
			downRoundTurn := 20 + rand.Intn(30) // Months 20-50
			if downRoundTurn < gs.Portfolio.MaxTurns {
				downRoundName := "Series A (Down)"
				if rand.Float64() < 0.5 {
					downRoundName = "Series B (Down)"
				}
				downAmount := int64(5000000 + rand.Intn(15000000)) // $5M-$20M
				gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
					CompanyName:   startup.Name,
					RoundName:     downRoundName,
					ScheduledTurn: downRoundTurn,
					RaiseAmount:   downAmount,
					IsDownRound:   true,
				})
			}
		}
	}
}


func (gs *GameState) ScheduleAcquisitions() {
	gs.AcquisitionQueue = []AcquisitionEvent{}

	// 40% of companies get acquisition offers
	for _, startup := range gs.AvailableStartups {
		if rand.Float64() < 0.4 {
			// Acquisitions happen between months 24-60
			acqTurn := 24 + rand.Intn(37)
			if acqTurn < gs.Portfolio.MaxTurns {
				// Multiple ranges from 3x to 6x EBITDA (4x average)
				multiple := 3.0 + rand.Float64()*3.0

				// Due diligence quality
				dueDiligence := "normal"
				roll := rand.Float64()
				if roll < 0.15 { // 15% bad due diligence
					dueDiligence = "bad"
					multiple *= 0.6 // Offer falls through or gets cut 40%
				} else if roll > 0.85 { // 15% great due diligence
					dueDiligence = "good"
					multiple *= 1.2 // Offer increases 20%
				}

				gs.AcquisitionQueue = append(gs.AcquisitionQueue, AcquisitionEvent{
					CompanyName:   startup.Name,
					ScheduledTurn: acqTurn,
					OfferMultiple: multiple,
					DueDiligence:  dueDiligence,
				})
			}
		}
	}
}

func (gs *GameState) ScheduleDramaticEvents() {
	gs.DramaticEventQueue = []DramaticEvent{}

	// Event frequency based on difficulty
	// Easy: 10%, Medium: 20%, Hard: 30%, Expert: 40%
	eventChance := 0.10
	if gs.Difficulty.Name == "Medium" {
		eventChance = 0.20
	} else if gs.Difficulty.Name == "Hard" {
		eventChance = 0.30
	} else if gs.Difficulty.Name == "Expert" {
		eventChance = 0.40
	}

	eventTypes := []string{
		"cofounder_split", "scandal", "lawsuit", "pivot_fail",
		"fraud", "data_breach", "key_hire_quit", "regulatory_issue",
		"competitor_attack", "product_failure",
	}

	for _, startup := range gs.AvailableStartups {
		if rand.Float64() < eventChance {
			// Events happen between months 6-55
			eventTurn := 6 + rand.Intn(50)
			if eventTurn < gs.Portfolio.MaxTurns {
				eventType := eventTypes[rand.Intn(len(eventTypes))]

				// Determine severity (difficulty affects this)
				severityRoll := rand.Float64()
				severity := "minor"
				impactPercent := 0.85 // 15% drop

				if gs.Difficulty.Name == "Hard" || gs.Difficulty.Name == "Expert" {
					// Harder difficulties have worse outcomes
					if severityRoll < 0.25 {
						severity = "severe"
						impactPercent = 0.40 // 60% drop
					} else if severityRoll < 0.55 {
						severity = "moderate"
						impactPercent = 0.65 // 35% drop
					}
				} else if gs.Difficulty.Name == "Medium" {
					if severityRoll < 0.15 {
						severity = "severe"
						impactPercent = 0.50 // 50% drop
					} else if severityRoll < 0.40 {
						severity = "moderate"
						impactPercent = 0.70 // 30% drop
					}
				} else {
					// Easy mode
					if severityRoll < 0.10 {
						severity = "moderate"
						impactPercent = 0.75 // 25% drop
					}
				}

				gs.DramaticEventQueue = append(gs.DramaticEventQueue, DramaticEvent{
					CompanyName:   startup.Name,
					ScheduledTurn: eventTurn,
					EventType:     eventType,
					Severity:      severity,
					ImpactPercent: impactPercent,
				})
			}
		}
	}
}