package game

import (
	"fmt"
	"math/rand"
)

// SecondaryOffer represents an offer to buy a stake from the player
type SecondaryOffer struct {
	CompanyName   string
	BuyerName     string // AI investor name
	BuyerFirm     string
	OfferAmount   int64   // Total offer for the stake
	OfferPercent  float64 // Percent of current valuation (0.7-0.9)
	EquityOffered float64 // Percentage of equity being sold
	CurrentValue  int64   // Current value of the stake
	Turn          int     // When offer was made
	ExpiresIn     int     // Turns until expiration
}

// GenerateSecondaryOffers creates offers for qualifying investments
func (gs *GameState) GenerateSecondaryOffers() []SecondaryOffer {
	offers := []SecondaryOffer{}

	for _, inv := range gs.Portfolio.Investments {
		// Requirements for secondary sale:
		// 1. Held for 12+ months
		// 2. Company not failed (valuation > 0)
		// 3. 10% chance per eligible investment per turn

		if inv.MonthsHeld < 12 {
			continue
		}

		if inv.CurrentValuation <= 0 {
			continue
		}

		// Check if already has an active offer
		hasActiveOffer := false
		for _, offer := range gs.SecondaryMarketOffers {
			if offer.CompanyName == inv.CompanyName && offer.ExpiresIn > 0 {
				hasActiveOffer = true
				break
			}
		}

		if hasActiveOffer {
			continue
		}

		// 10% chance to get an offer
		if rand.Float64() > 0.10 {
			continue
		}

		// Select a buyer from AI players
		if len(gs.AIPlayers) == 0 {
			continue
		}

		buyer := gs.AIPlayers[rand.Intn(len(gs.AIPlayers))]

		// Calculate offer
		currentStakeValue := int64(float64(inv.CurrentValuation) * inv.EquityPercent / 100.0)

		// Offer is 70-90% of current value (secondary market discount)
		// Better companies get better offers
		offerPercent := 0.70

		// Adjust based on company performance
		if inv.CurrentValuation > inv.InitialValuation*2 {
			offerPercent = 0.85 // Strong performer
		} else if inv.CurrentValuation > inv.InitialValuation {
			offerPercent = 0.80 // Good performer
		} else if inv.CurrentValuation < int64(float64(inv.InitialValuation)*0.7) {
			offerPercent = 0.70 // Struggling
		}

		// Adjust based on buyer strategy
		if buyer.Strategy == "aggressive" {
			offerPercent += 0.05 // Aggressive buyers pay more
		} else if buyer.Strategy == "conservative" {
			offerPercent -= 0.03 // Conservative buyers pay less
		}

		// Ensure bounds
		if offerPercent < 0.70 {
			offerPercent = 0.70
		}
		if offerPercent > 0.90 {
			offerPercent = 0.90
		}

		offerAmount := int64(float64(currentStakeValue) * offerPercent)

		// Transaction fee (20% of proceeds) built into offer
		// This is already factored into the discount

		offer := SecondaryOffer{
			CompanyName:   inv.CompanyName,
			BuyerName:     buyer.Name,
			BuyerFirm:     buyer.Firm,
			OfferAmount:   offerAmount,
			OfferPercent:  offerPercent,
			EquityOffered: inv.EquityPercent,
			CurrentValue:  currentStakeValue,
			Turn:          gs.Portfolio.Turn,
			ExpiresIn:     3, // Offer valid for 3 turns
		}

		offers = append(offers, offer)
	}

	return offers
}

// AcceptSecondaryOffer processes a secondary sale
func (gs *GameState) AcceptSecondaryOffer(offerIndex int) error {
	if offerIndex < 0 || offerIndex >= len(gs.SecondaryMarketOffers) {
		return fmt.Errorf("invalid offer index")
	}

	offer := gs.SecondaryMarketOffers[offerIndex]

	// Check if offer has expired
	if offer.ExpiresIn <= 0 {
		return fmt.Errorf("offer has expired")
	}

	// Find the investment
	invIdx := -1
	for i, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == offer.CompanyName {
			invIdx = i
			break
		}
	}

	if invIdx == -1 {
		return fmt.Errorf("investment not found")
	}

	// Remove the investment from portfolio
	gs.Portfolio.Investments = append(
		gs.Portfolio.Investments[:invIdx],
		gs.Portfolio.Investments[invIdx+1:]...,
	)

	// Add cash (already includes transaction fee discount)
	gs.Portfolio.Cash += offer.OfferAmount

	// Remove the offer
	gs.SecondaryMarketOffers = append(
		gs.SecondaryMarketOffers[:offerIndex],
		gs.SecondaryMarketOffers[offerIndex+1:]...,
	)

	gs.updateNetWorth()

	return nil
}

// DeclineSecondaryOffer removes an offer
func (gs *GameState) DeclineSecondaryOffer(offerIndex int) error {
	if offerIndex < 0 || offerIndex >= len(gs.SecondaryMarketOffers) {
		return fmt.Errorf("invalid offer index")
	}

	// Remove the offer
	gs.SecondaryMarketOffers = append(
		gs.SecondaryMarketOffers[:offerIndex],
		gs.SecondaryMarketOffers[offerIndex+1:]...,
	)

	return nil
}

// ProcessSecondaryOfferExpirations decrements expiration timers
func (gs *GameState) ProcessSecondaryOfferExpirations() []string {
	messages := []string{}

	for i := range gs.SecondaryMarketOffers {
		gs.SecondaryMarketOffers[i].ExpiresIn--

		if gs.SecondaryMarketOffers[i].ExpiresIn == 0 {
			messages = append(messages, fmt.Sprintf(
				"Secondary market offer for %s has expired",
				gs.SecondaryMarketOffers[i].CompanyName))
		}
	}

	// Remove expired offers
	activeOffers := []SecondaryOffer{}
	for _, offer := range gs.SecondaryMarketOffers {
		if offer.ExpiresIn > 0 {
			activeOffers = append(activeOffers, offer)
		}
	}
	gs.SecondaryMarketOffers = activeOffers

	return messages
}

// GetSecondaryMarketSummary returns stats about secondary sales
func (gs *GameState) GetSecondaryMarketSummary() (totalSales int, totalProceeds int64) {
	// This would need to track secondary sales history
	// For now, returning zeros as placeholder
	return 0, 0
}

// CalculateSecondaryROI calculates ROI on a potential secondary sale
func CalculateSecondaryROI(offer SecondaryOffer, amountInvested int64) float64 {
	if amountInvested == 0 {
		return 0.0
	}

	return (float64(offer.OfferAmount) - float64(amountInvested)) / float64(amountInvested) * 100.0
}

// ShouldAcceptSecondaryOffer provides AI recommendation
func ShouldAcceptSecondaryOffer(offer SecondaryOffer, inv Investment) (bool, string) {
	// Calculate current ROI
	currentROI := (float64(offer.CurrentValue) - float64(inv.AmountInvested)) / float64(inv.AmountInvested) * 100.0

	// Recommendation logic
	if currentROI >= 300 { // 3x return
		return true, "Strong return achieved (3x+). Consider taking profits."
	}

	if currentROI >= 200 { // 2x return
		return true, "Good return (2x+). De-risking portfolio makes sense."
	}

	if currentROI >= 100 && inv.MonthsHeld >= 24 { // 2x return after 2 years
		return true, "Solid return after long hold. Opportunity to redeploy capital."
	}

	if currentROI < 0 { // Underwater
		return true, "Opportunity to cut losses and redeploy capital."
	}

	if inv.CurrentValuation < int64(float64(inv.InitialValuation)*0.5) { // Down 50%+
		return true, "Significant underperformance. Consider cutting losses."
	}

	// Hold recommendations
	if currentROI >= 50 && currentROI < 150 {
		return false, "Moderate gains. Consider holding for higher returns."
	}

	return false, "Company showing promise. Hold for potential upside."
}

// GetSecondaryMarketFee returns the transaction fee percentage
func GetSecondaryMarketFee() float64 {
	return 0.20 // 20% fee (already built into offer discounts)
}
