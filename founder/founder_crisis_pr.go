package founder

import (
	"fmt"
	"math/rand"
)

// SpawnPRCrisis generates a PR crisis
func (fs *FounderState) SpawnPRCrisis() *PRCrisis {
	// Unlock: $1M+ ARR OR Series A raised
	arr := fs.MRR * 12
	hasSeriesA := false
	for _, round := range fs.FundingRounds {
		if round.RoundName == "Series A" {
			hasSeriesA = true
			break
		}
	}
	if arr < 1000000 && !hasSeriesA {
		return nil
	}

	// Probability: 3% per month, increases with negative events
	baseProbability := 0.03
	if fs.PRProgram != nil && fs.PRProgram.BrandScore < 50 {
		baseProbability = 0.08 // 8% if brand score low
	}

	if rand.Float64() > baseProbability {
		return nil
	}

	// Crisis types
	crisisTypes := []string{"scandal", "product_failure", "layoffs", "founder_drama", "competitor_attack"}
	crisisType := crisisTypes[rand.Intn(len(crisisTypes))]

	// Severity
	severityRoll := rand.Float64()
	severity := "minor"
	if severityRoll < 0.2 {
		severity = "critical"
	} else if severityRoll < 0.4 {
		severity = "major"
	} else if severityRoll < 0.7 {
		severity = "moderate"
	}

	// Media coverage
	mediaOutlets := []string{"TechCrunch", "WSJ", "Forbes", "Bloomberg", "The Verge", "Trade Publication"}
	numOutlets := 1
	if severity == "critical" {
		numOutlets = 3 + rand.Intn(3) // 3-5 outlets
	} else if severity == "major" {
		numOutlets = 2 + rand.Intn(2) // 2-3 outlets
	}

	mediaCoverage := []string{}
	for i := 0; i < numOutlets && i < len(mediaOutlets); i++ {
		mediaCoverage = append(mediaCoverage, mediaOutlets[rand.Intn(len(mediaOutlets))])
	}

	// Impact based on severity
	var cacImpact, churnImpact, brandDamage float64
	var durationMonths int

	switch severity {
	case "critical":
		cacImpact = 1.5 + rand.Float64()*0.5  // 1.5-2.0x CAC
		churnImpact = 0.08 + rand.Float64()*0.04 // 8-12% churn
		brandDamage = 0.3 + rand.Float64()*0.2   // 30-50% brand damage
		durationMonths = 9 + rand.Intn(4)        // 9-12 months
	case "major":
		cacImpact = 1.3 + rand.Float64()*0.3  // 1.3-1.6x CAC
		churnImpact = 0.05 + rand.Float64()*0.03 // 5-8% churn
		brandDamage = 0.2 + rand.Float64()*0.2   // 20-40% brand damage
		durationMonths = 6 + rand.Intn(4)        // 6-9 months
	case "moderate":
		cacImpact = 1.2 + rand.Float64()*0.2  // 1.2-1.4x CAC
		churnImpact = 0.03 + rand.Float64()*0.02 // 3-5% churn
		brandDamage = 0.1 + rand.Float64()*0.1   // 10-20% brand damage
		durationMonths = 3 + rand.Intn(4)        // 3-6 months
	case "minor":
		cacImpact = 1.1 + rand.Float64()*0.1  // 1.1-1.2x CAC
		churnImpact = 0.01 + rand.Float64()*0.02 // 1-3% churn
		brandDamage = 0.05 + rand.Float64()*0.05 // 5-10% brand damage
		durationMonths = 1 + rand.Intn(3)        // 1-3 months
	}

	crisis := PRCrisis{
		Type:           crisisType,
		Severity:       severity,
		Month:          fs.Turn,
		MediaCoverage:  mediaCoverage,
		Response:       "none",
		ResponseCost:   0,
		DurationMonths: durationMonths,
		CACImpact:      cacImpact,
		ChurnImpact:    churnImpact,
		BrandDamage:    brandDamage,
		Resolved:       false,
		ResolutionMonth: 0,
	}

	fs.PRCrises = append(fs.PRCrises, crisis)
	fs.ActivePRCrisis = &crisis

	return &crisis
}

// RespondToPRCrisis handles PR crisis response
func (fs *FounderState) RespondToPRCrisis(responseType string) error {
	if fs.ActivePRCrisis == nil {
		return fmt.Errorf("no active PR crisis")
	}

	crisis := fs.ActivePRCrisis

	validResponses := map[string]bool{
		"deny":       true,
		"apologize":  true,
		"transparent": true,
		"aggressive": true,
	}
	if !validResponses[responseType] {
		return fmt.Errorf("invalid response type: %s", responseType)
	}

	// Response costs and effectiveness
	var cost int64
	var effectiveness float64
	var outcome string

	switch responseType {
	case "deny":
		cost = 10000 // Cheap but risky
		effectiveness = 0.2 // 20% effective
		if rand.Float64() < 0.5 {
			outcome = "escalated" // 50% chance makes it worse
		} else {
			outcome = "contained"
		}

	case "apologize":
		cost = 50000 // Moderate cost
		effectiveness = 0.5 // 50% effective
		if rand.Float64() < 0.3 {
			outcome = "escalated"
		} else {
			outcome = "contained"
		}

	case "transparent":
		cost = 100000 // Higher cost
		effectiveness = 0.8 // 80% effective
		if rand.Float64() < 0.1 {
			outcome = "escalated"
		} else {
			outcome = "resolved"
		}

	case "aggressive":
		cost = 200000 // Expensive (legal)
		effectiveness = 0.4 // 40% effective (can backfire)
		if rand.Float64() < 0.4 {
			outcome = "escalated" // 40% chance backfires
		} else {
			outcome = "contained"
		}
	}

	// PR firm helps
	if fs.PRProgram != nil && fs.PRProgram.HasPRFirm {
		effectiveness += 0.2 // +20% effectiveness
		cost = int64(float64(cost) * 0.7) // 30% cost reduction
	}

	if cost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	fs.Cash -= cost
	crisis.Response = responseType
	crisis.ResponseCost = cost

	// Apply effectiveness
	crisis.CACImpact = 1.0 + (crisis.CACImpact-1.0)*(1.0-effectiveness)
	crisis.ChurnImpact *= (1.0 - effectiveness)
	crisis.BrandDamage *= (1.0 - effectiveness)

	// Create response record
	response := CrisisResponse{
		CrisisType:    crisis.Type,
		ResponseType:  responseType,
		Cost:          cost,
		Effectiveness: effectiveness,
		Outcome:       outcome,
		Month:         fs.Turn,
	}
	fs.CrisisResponses = append(fs.CrisisResponses, response)

	// Check if resolved
	if outcome == "resolved" || (outcome == "contained" && effectiveness > 0.6) {
		crisis.Resolved = true
		crisis.ResolutionMonth = fs.Turn
		fs.ActivePRCrisis = nil
	}

	return nil
}

// ProcessPRCrises processes active PR crises
func (fs *FounderState) ProcessPRCrises() []string {
	var messages []string

	if fs.ActivePRCrisis != nil {
		crisis := fs.ActivePRCrisis

		if !crisis.Resolved {
			// Apply ongoing damage
			// CAC increase
			cacMultiplier := crisis.CACImpact
			fs.BaseCAC = int64(float64(fs.BaseCAC) * cacMultiplier)

			// Churn increase
			fs.CustomerChurnRate += crisis.ChurnImpact

			// Brand damage
			if fs.PRProgram != nil {
				fs.PRProgram.BrandScore -= int(crisis.BrandDamage * 10)
				if fs.PRProgram.BrandScore < 0 {
					fs.PRProgram.BrandScore = 0
				}
			}

			// Check if duration expired
			monthsActive := fs.Turn - crisis.Month
			if monthsActive >= crisis.DurationMonths {
				crisis.Resolved = true
				crisis.ResolutionMonth = fs.Turn
				fs.ActivePRCrisis = nil
				messages = append(messages, "⏰ PR crisis duration expired (damage reduced)")
			}
		}
	}

	return messages
}

