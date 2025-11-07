package founder

import "fmt"

// LaunchCompetitiveIntel starts intel gathering
func (fs *FounderState) LaunchCompetitiveIntel(analystSalary int64) error {
	if len(fs.Competitors) < 5 {
		return fmt.Errorf("need at least 5 competitors to justify analyst")
	}

	fs.CompetitiveIntel = &CompetitiveIntel{
		HasAnalyst:      true,
		MonthlyBudget:   analystSalary,
		IntelReports:    []IntelReport{},
		BattleCards:     []BattleCard{},
		WinLossInsights: make(map[string]int),
		AnalystSalary:   analystSalary,
		LaunchedMonth:   fs.Turn,
	}

	return nil
}

// CommissionIntelReport creates a competitor analysis report
func (fs *FounderState) CommissionIntelReport(competitorName string, cost int64) error {
	if fs.CompetitiveIntel == nil {
		return fmt.Errorf("no competitive intel program active")
	}

	if fs.Cash < cost {
		return fmt.Errorf("insufficient cash")
	}

	fs.Cash -= cost

	report := IntelReport{
		CompetitorName: competitorName,
		Pricing:        make(map[string]int64),
		Features:       []string{"API", "Mobile", "Analytics"},
		Funding:        "Series B",
		TeamSize:       50,
		RecentMoves:    []string{"Launched new feature", "Hired VP Sales"},
		ThreatLevel:    "high",
		Cost:           cost,
		Month:          fs.Turn,
	}

	fs.CompetitiveIntel.IntelReports = append(fs.CompetitiveIntel.IntelReports, report)

	return nil
}

