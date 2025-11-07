package founder

import "math"

// InitializeTechnicalDebt sets up tech debt tracking
func (fs *FounderState) InitializeTechnicalDebt() {
	if fs.TechnicalDebt == nil {
		fs.TechnicalDebt = &TechnicalDebt{
			CurrentLevel:       20, // Start at 20/100
			VelocityImpact:     1.0,
			BugFrequency:       0.02,
			SecurityRisks:      0,
			ScalingProblems:    false,
			EngineerMorale:     0.80,
			RefactoringCosts:   0,
			MonthsSinceRefactor: 0,
		}
	}
}

// AccumulateTechnicalDebt increases debt each month based on factors
func (fs *FounderState) AccumulateTechnicalDebt() []string {
	var messages []string

	if fs.TechnicalDebt == nil {
		fs.InitializeTechnicalDebt()
	}

	// Debt accumulates when shipping fast without senior engineers
	debtIncrease := 2
	
	// Check for CTO in Executives
	hasCTO := false
	for _, exec := range fs.Team.Executives {
		if exec.Role == "CTO" {
			hasCTO = true
			break
		}
	}
	if hasCTO {
		debtIncrease = 1 // CTO reduces accumulation by 50%
	}

	// Check for senior engineers (simplified - all engineers count as experience)
	seniorEngineers := len(fs.Team.Engineers) / 3 // Assume 1/3 are senior
	if seniorEngineers > 0 {
		// Having senior engineers helps
	}

	if seniorEngineers == 0 && len(fs.Team.Engineers) > 3 {
		debtIncrease += 2 // No seniors = more debt
	}

	fs.TechnicalDebt.CurrentLevel += debtIncrease
	if fs.TechnicalDebt.CurrentLevel > 100 {
		fs.TechnicalDebt.CurrentLevel = 100
	}
	fs.TechnicalDebt.MonthsSinceRefactor++

	// Apply impacts
	if fs.TechnicalDebt.CurrentLevel > 40 {
		fs.TechnicalDebt.VelocityImpact = 0.90 // -10% velocity
	}
	if fs.TechnicalDebt.CurrentLevel > 60 {
		fs.TechnicalDebt.VelocityImpact = 0.75 // -25% velocity
		fs.TechnicalDebt.BugFrequency = 0.05
		fs.CustomerChurnRate = math.Min(0.30, fs.CustomerChurnRate+0.01) // +1% churn from bugs
	}
	if fs.TechnicalDebt.CurrentLevel > 80 {
		fs.TechnicalDebt.ScalingProblems = true
		messages = append(messages, "🚨 HIGH TECH DEBT: Scaling problems, bugs increasing churn!")
	}

	return messages
}

// RefactorTechDebt reduces debt through investment
func (fs *FounderState) RefactorTechDebt(cost int64, engineersAllocated int) error {
	if fs.TechnicalDebt == nil {
		fs.InitializeTechnicalDebt()
	}

	if fs.Cash < cost {
		return nil
	}

	fs.Cash -= cost
	fs.TechnicalDebt.RefactoringCosts = cost
	
	// Reduce debt based on investment
	debtReduction := 10 + (engineersAllocated * 5)
	fs.TechnicalDebt.CurrentLevel -= debtReduction
	if fs.TechnicalDebt.CurrentLevel < 10 {
		fs.TechnicalDebt.CurrentLevel = 10 // Minimum 10 debt
	}
	
	fs.TechnicalDebt.MonthsSinceRefactor = 0
	fs.TechnicalDebt.VelocityImpact = 1.0
	fs.TechnicalDebt.BugFrequency = 0.02
	fs.TechnicalDebt.ScalingProblems = false

	return nil
}

