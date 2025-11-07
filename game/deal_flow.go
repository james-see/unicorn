package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
)

// GenerateStartupsWithReputation generates startups based on player reputation
func GenerateStartupsWithReputation(reputation *VCReputation, count int, filename string) ([]Startup, error) {
	// Load base startups from file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read startups file: %v", err)
	}

	var allStartups []Startup
	if err := json.Unmarshal(file, &allStartups); err != nil {
		return nil, fmt.Errorf("failed to parse startups file: %v", err)
	}

	if len(allStartups) < count {
		return allStartups, nil // Return all if not enough
	}

	aggregateRep := reputation.GetAggregateReputation()

	// Determine deal quality distribution based on reputation
	tier1Percent := 0.0 // Hot deals
	tier2Percent := 0.0 // Standard deals
	// tier3 is implicit: 1.0 - (tier1Percent + tier2Percent)

	if aggregateRep >= 70 {
		// Tier 1 (Hot Deals) access
		tier1Percent = 0.25 // 25% hot deals
		tier2Percent = 0.60 // 60% standard (15% struggling implicit)
	} else if aggregateRep >= 40 {
		// Tier 2 (Standard Deals)
		tier1Percent = 0.05 // 5% hot deals (some lucky access)
		tier2Percent = 0.75 // 75% standard (20% struggling implicit)
	} else {
		// Tier 3 (Struggling Deals)
		tier1Percent = 0.0  // No hot deals
		tier2Percent = 0.40 // 40% standard (60% struggling implicit)
	}

	// Shuffle startups
	rand.Shuffle(len(allStartups), func(i, j int) {
		allStartups[i], allStartups[j] = allStartups[j], allStartups[i]
	})

	// Adjust startups based on deal quality tier
	adjustedStartups := make([]Startup, 0, count)

	for i := 0; i < count && i < len(allStartups); i++ {
		startup := allStartups[i]

		// Determine which tier this startup slot should be
		roll := rand.Float64()
		var tier string

		if roll < tier1Percent {
			tier = "tier1" // Hot deal
		} else if roll < (tier1Percent + tier2Percent) {
			tier = "tier2" // Standard deal
		} else {
			tier = "tier3" // Struggling deal
		}

		// Adjust startup characteristics based on tier
		startup = adjustStartupForTier(startup, tier)

		adjustedStartups = append(adjustedStartups, startup)
	}

	return adjustedStartups, nil
}

// adjustStartupForTier modifies startup characteristics based on deal quality tier
func adjustStartupForTier(startup Startup, tier string) Startup {
	switch tier {
	case "tier1": // Hot deals - lower risk, higher growth
		// Reduce risk (to 0.2-0.4 range)
		if startup.RiskScore > 0.4 {
			startup.RiskScore = 0.2 + rand.Float64()*0.2 // 0.2-0.4
		}

		// Increase growth potential (to 0.7-0.9 range)
		if startup.GrowthPotential < 0.7 {
			startup.GrowthPotential = 0.7 + rand.Float64()*0.2 // 0.7-0.9
		}

		// Slightly higher initial valuation (hot deal premium)
		startup.Valuation = int64(float64(startup.Valuation) * (1.1 + rand.Float64()*0.2)) // 1.1-1.3x

	case "tier2": // Standard deals - balanced
		// Keep risk in 0.4-0.6 range
		if startup.RiskScore < 0.4 {
			startup.RiskScore = 0.4 + rand.Float64()*0.1
		} else if startup.RiskScore > 0.6 {
			startup.RiskScore = 0.5 + rand.Float64()*0.1
		}

		// Keep growth in 0.5-0.7 range
		if startup.GrowthPotential < 0.5 {
			startup.GrowthPotential = 0.5 + rand.Float64()*0.1
		} else if startup.GrowthPotential > 0.7 {
			startup.GrowthPotential = 0.6 + rand.Float64()*0.1
		}

		// Standard valuation (no adjustment)

	case "tier3": // Struggling deals - higher risk, lower growth
		// Increase risk (to 0.6-0.8 range)
		if startup.RiskScore < 0.6 {
			startup.RiskScore = 0.6 + rand.Float64()*0.2 // 0.6-0.8
		}

		// Decrease growth potential (to 0.3-0.5 range)
		if startup.GrowthPotential > 0.5 {
			startup.GrowthPotential = 0.3 + rand.Float64()*0.2 // 0.3-0.5
		}

		// Lower valuation (struggling company discount)
		startup.Valuation = int64(float64(startup.Valuation) * (0.7 + rand.Float64()*0.2)) // 0.7-0.9x
	}

	return startup
}

// GetDealFlowQualityMessage returns a message about deal flow based on reputation
func GetDealFlowQualityMessage(reputation *VCReputation) string {
	tier := reputation.GetDealQualityTier()
	aggregate := reputation.GetAggregateReputation()

	switch tier {
	case "Tier 1 (Hot Deals)":
		return fmt.Sprintf("🔥 EXCELLENT DEAL FLOW (Reputation: %.1f/100)\nYour strong reputation grants access to high-quality startups with lower risk and higher growth potential. Top founders want you on their cap table!", aggregate)

	case "Tier 2 (Standard Deals)":
		return fmt.Sprintf("✓ STANDARD DEAL FLOW (Reputation: %.1f/100)\nYou have access to solid startup opportunities with balanced risk/reward profiles. Continue building your reputation for access to better deals.", aggregate)

	default: // Tier 3
		return fmt.Sprintf("⚠️  LIMITED DEAL FLOW (Reputation: %.1f/100)\nYour deal flow consists mainly of higher-risk opportunities. Focus on building relationships and delivering strong returns to access better deals.", aggregate)
	}
}

// CalculateFounderReferralEffect applies founder referral bonus to deal quality
func CalculateFounderReferralEffect(avgFounderRelationship float64, startup *Startup) {
	bonus := GetFounderReferralBonus(avgFounderRelationship)

	if bonus > 0 {
		// Improve risk and growth slightly
		startup.RiskScore -= bonus / 2
		if startup.RiskScore < 0.1 {
			startup.RiskScore = 0.1
		}

		startup.GrowthPotential += bonus
		if startup.GrowthPotential > 1.0 {
			startup.GrowthPotential = 1.0
		}
	}
}

// GetReputationBonus returns reputation impact on various game mechanics
type ReputationBonus struct {
	DealQualityBonus   float64 // 0-0.15 improvement to deal quality
	FounderTrustBonus  float64 // 0-10 points to initial relationships
	ExitMultiplier     float64 // 1.0-1.1 multiplier on exit values
	BoardInfluence     float64 // 0-0.2 better board vote outcomes
	NetworkEffectBonus float64 // 0-0.1 bonus to company valuations from network
}

func GetReputationBonus(reputation *VCReputation) ReputationBonus {
	aggregate := reputation.GetAggregateReputation()

	bonus := ReputationBonus{
		DealQualityBonus:   0.0,
		FounderTrustBonus:  0.0,
		ExitMultiplier:     1.0,
		BoardInfluence:     0.0,
		NetworkEffectBonus: 0.0,
	}

	// Scale bonuses based on aggregate reputation
	if aggregate >= 80 {
		bonus.DealQualityBonus = 0.15
		bonus.FounderTrustBonus = 10.0
		bonus.ExitMultiplier = 1.10
		bonus.BoardInfluence = 0.20
		bonus.NetworkEffectBonus = 0.10
	} else if aggregate >= 70 {
		bonus.DealQualityBonus = 0.12
		bonus.FounderTrustBonus = 8.0
		bonus.ExitMultiplier = 1.08
		bonus.BoardInfluence = 0.15
		bonus.NetworkEffectBonus = 0.08
	} else if aggregate >= 60 {
		bonus.DealQualityBonus = 0.08
		bonus.FounderTrustBonus = 5.0
		bonus.ExitMultiplier = 1.05
		bonus.BoardInfluence = 0.10
		bonus.NetworkEffectBonus = 0.05
	} else if aggregate >= 50 {
		bonus.DealQualityBonus = 0.04
		bonus.FounderTrustBonus = 3.0
		bonus.ExitMultiplier = 1.02
		bonus.BoardInfluence = 0.05
		bonus.NetworkEffectBonus = 0.02
	}

	return bonus
}

// ApplyReputationBonusToInvestment applies reputation bonus when making investment
func ApplyReputationBonusToInvestment(inv *Investment, bonus ReputationBonus) {
	// Apply founder trust bonus to initial relationship
	inv.RelationshipScore += bonus.FounderTrustBonus

	if inv.RelationshipScore > 100 {
		inv.RelationshipScore = 100
	}
}
