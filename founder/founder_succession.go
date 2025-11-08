package founder

import (
	"fmt"
	"math/rand"
)

// InitializeKeyPersonRisks initializes key person risk assessment
func InitializeKeyPersonRisks(fs *FounderState) {
	if fs.KeyPersonRisks == nil {
		fs.KeyPersonRisks = []KeyPersonRisk{}
	}
	if fs.KeyPersonEvents == nil {
		fs.KeyPersonEvents = []KeyPersonEvent{}
	}
	if fs.SuccessionPlans == nil {
		fs.SuccessionPlans = []SuccessionPlan{}
	}

	// Assess key persons
	fs.AssessKeyPersonRisks()
}

// CanHaveKeyPersonRisk checks if company is large enough
func (fs *FounderState) CanHaveKeyPersonRisk() bool {
	// Unlock: $2M+ ARR OR 50+ employees
	arr := fs.MRR * 12
	totalEmployees := fs.Team.TotalEmployees
	return arr >= 2000000 || totalEmployees >= 50
}

// AssessKeyPersonRisks assesses risk for all key persons
func (fs *FounderState) AssessKeyPersonRisks() {
	if !fs.CanHaveKeyPersonRisk() {
		return
	}

	// Clear existing risks
	fs.KeyPersonRisks = []KeyPersonRisk{}

	// Founder
	founderRisk := KeyPersonRisk{
		PersonName:      "Founder",
		Role:            "CEO",
		RiskLevel:       "medium",
		Dependencies:    []string{"Company vision", "Investor relations", "Team culture"},
		SuccessionReady: false,
		RetentionScore:  0.8, // 80% retention
	}
	fs.KeyPersonRisks = append(fs.KeyPersonRisks, founderRisk)

	// CTO
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCTO {
			ctoRisk := KeyPersonRisk{
				PersonName:      exec.Name,
				Role:            "CTO",
				RiskLevel:       "high",
				Dependencies:    []string{"Product architecture", "Engineering team", "Technical decisions"},
				SuccessionReady: false,
				RetentionScore:  0.7, // 70% retention
			}
			fs.KeyPersonRisks = append(fs.KeyPersonRisks, ctoRisk)
			break
		}
	}

	// CFO
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCFO {
			cfoRisk := KeyPersonRisk{
				PersonName:      exec.Name,
				Role:            "CFO",
				RiskLevel:       "medium",
				Dependencies:    []string{"Financial management", "Fundraising", "Compliance"},
				SuccessionReady: false,
				RetentionScore:  0.75, // 75% retention
			}
			fs.KeyPersonRisks = append(fs.KeyPersonRisks, cfoRisk)
			break
		}
	}

	// Head of Sales
	if len(fs.Team.Sales) > 5 {
		headOfSalesRisk := KeyPersonRisk{
			PersonName:      "Head of Sales",
			Role:            "Head of Sales",
			RiskLevel:       "high",
			Dependencies:    []string{"Sales strategy", "Customer relationships", "Revenue targets"},
			SuccessionReady: false,
			RetentionScore:  0.65, // 65% retention (sales people move around)
		}
		fs.KeyPersonRisks = append(fs.KeyPersonRisks, headOfSalesRisk)
	}
}

// CreateSuccessionPlan creates a succession plan for a key person
func (fs *FounderState) CreateSuccessionPlan(personName string, backupPerson string) error {
	// Check if person exists in risks
	personExists := false
	for _, kpr := range fs.KeyPersonRisks {
		if kpr.PersonName == personName {
			personExists = true
			break
		}
	}
	if !personExists {
		return fmt.Errorf("key person %s not found", personName)
	}

	// Check if plan already exists
	for _, sp := range fs.SuccessionPlans {
		if sp.PersonName == personName {
			return fmt.Errorf("succession plan for %s already exists", personName)
		}
	}

	// Training takes 3-6 months
	trainingMonths := 3 + rand.Intn(4)

	plan := SuccessionPlan{
		PersonName:     personName,
		BackupPerson:   backupPerson,
		TrainingMonths: trainingMonths,
		Ready:          false,
		MonthCreated:   fs.Turn,
	}

	fs.SuccessionPlans = append(fs.SuccessionPlans, plan)

	// Update risk assessment
	for i := range fs.KeyPersonRisks {
		if fs.KeyPersonRisks[i].PersonName == personName {
			fs.KeyPersonRisks[i].SuccessionReady = true
			fs.KeyPersonRisks[i].RetentionScore += 0.1 // +10% retention with plan
			break
		}
	}

	return nil
}

// SpawnKeyPersonEvent generates a key person leaving event
func (fs *FounderState) SpawnKeyPersonEvent() *KeyPersonEvent {
	if !fs.CanHaveKeyPersonRisk() {
		return nil
	}

	// Probability: 2% per month per key person
	if len(fs.KeyPersonRisks) == 0 {
		return nil
	}

	totalProbability := float64(len(fs.KeyPersonRisks)) * 0.02
	if rand.Float64() > totalProbability {
		return nil
	}

	// Select random key person
	personIndex := rand.Intn(len(fs.KeyPersonRisks))
	person := fs.KeyPersonRisks[personIndex]

	// Check retention score
	if rand.Float64() < person.RetentionScore {
		return nil // Person stays
	}

	// Event types
	eventTypes := []string{"quit", "poached", "illness", "scandal", "death"}
	eventType := eventTypes[rand.Intn(len(eventTypes))]

	// Death is very rare
	if eventType == "death" && rand.Float64() > 0.05 {
		eventType = eventTypes[rand.Intn(len(eventTypes)-1)] // Re-roll excluding death
	}

	// Impact based on role
	impact := EventImpact{
		CACChange:    1.0,
		ChurnChange:  0.0,
		GrowthChange: 1.0,
		CashCost:     0,
		MRRChange:    1.0,
	}

	switch person.Role {
	case "CEO":
		impact.GrowthChange = 0.7  // -30% growth
		impact.ChurnChange = 0.03  // +3% churn
		impact.CACChange = 1.2     // +20% CAC
	case "CTO":
		impact.GrowthChange = 0.6  // -40% growth (product velocity)
		impact.ChurnChange = 0.05   // +5% churn (product issues)
	case "CFO":
		impact.CashCost = 50000     // $50k immediate cost
		impact.GrowthChange = 0.9   // -10% growth
	case "Head of Sales":
		impact.GrowthChange = 0.5   // -50% growth (revenue impact)
		impact.CACChange = 1.3      // +30% CAC
	}

	// Replacement cost
	replacementCost := int64(50000 + rand.Int63n(150000)) // $50-200k

	// Recovery months
	recoveryMonths := 3
	if person.Role == "CEO" {
		recoveryMonths = 6
	}
	if !person.SuccessionReady {
		recoveryMonths += 3 // +3 months without succession plan
	}

	event := KeyPersonEvent{
		PersonName:      person.PersonName,
		EventType:       eventType,
		Month:           fs.Turn,
		Impact:          impact,
		ReplacementCost: replacementCost,
		RecoveryMonths:  recoveryMonths,
		Resolved:        false,
	}

	fs.KeyPersonEvents = append(fs.KeyPersonEvents, event)

	// Apply immediate impact
	fs.BaseCAC = int64(float64(fs.BaseCAC) * impact.CACChange)
	fs.CustomerChurnRate += impact.ChurnChange
	fs.MonthlyGrowthRate *= impact.GrowthChange
	fs.Cash -= impact.CashCost

	return &event
}

// ProcessKeyPersonEvents processes active key person events
func (fs *FounderState) ProcessKeyPersonEvents() []string {
	var messages []string

	for i := range fs.KeyPersonEvents {
		event := &fs.KeyPersonEvents[i]
		if event.Resolved {
			continue
		}

		// Check if recovery complete
		monthsSinceEvent := fs.Turn - event.Month
		if monthsSinceEvent >= event.RecoveryMonths {
			event.Resolved = true

			// Hire replacement
			fs.Cash -= event.ReplacementCost

			// Restore metrics gradually
			fs.BaseCAC = int64(float64(fs.BaseCAC) / event.Impact.CACChange)
			fs.CustomerChurnRate -= event.Impact.ChurnChange
			fs.MonthlyGrowthRate /= event.Impact.GrowthChange

			messages = append(messages, fmt.Sprintf("✅ Replaced %s after %d months", event.PersonName, event.RecoveryMonths))
		}
	}

	return messages
}

// ProcessSuccessionPlans processes succession plan training
func (fs *FounderState) ProcessSuccessionPlans() []string {
	var messages []string

	for i := range fs.SuccessionPlans {
		plan := &fs.SuccessionPlans[i]
		if plan.Ready {
			continue
		}

		monthsSinceCreated := fs.Turn - plan.MonthCreated
		if monthsSinceCreated >= plan.TrainingMonths {
			plan.Ready = true
			messages = append(messages, fmt.Sprintf("✅ Succession plan ready for %s (backup: %s)", plan.PersonName, plan.BackupPerson))
		}
	}

	return messages
}

