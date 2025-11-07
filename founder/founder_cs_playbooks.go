package founder

import "fmt"

// LaunchCSPlaybook starts a customer success program
func (fs *FounderState) LaunchCSPlaybook(name string, budget int64) error {
	if fs.Customers < 100 {
		return fmt.Errorf("need at least 100 customers for CS playbooks")
	}

	csPlaybook := CSPlaybook{
		Name:           name,
		CSHeadcount:    len(fs.Team.CustomerSuccess),
		MonthlyBudget:  budget,
		ToolCosts:      5000, // $5k/month for CS tools
		ChurnReduction: 0.03,
		UpsellRate:     0.05,
		NPSScore:       50, // Start at 50
		Active:         true,
		LaunchedMonth:  fs.Turn,
	}

	fs.CSPlaybooks = append(fs.CSPlaybooks, csPlaybook)

	// Apply immediate churn reduction
	if fs.CustomerChurnRate > 0.05 {
		fs.CustomerChurnRate -= 0.03
	}

	return nil
}

// UpdateCSPlaybooks processes customer success monthly
func (fs *FounderState) UpdateCSPlaybooks() []string {
	var messages []string

	for i := range fs.CSPlaybooks {
		playbook := &fs.CSPlaybooks[i]
		if !playbook.Active {
			continue
		}

		// Pay costs
		totalCost := playbook.MonthlyBudget + playbook.ToolCosts
		if fs.Cash < totalCost {
			continue
		}
		fs.Cash -= totalCost

		// Improve NPS over time
		if playbook.NPSScore < 80 {
			playbook.NPSScore += 2
		}

		// Track at-risk customers
		atRisk := 0
		for j := range fs.CustomerList {
			if fs.CustomerList[j].IsActive && fs.CustomerList[j].HealthScore < 0.40 {
				atRisk++
				// CS intervention
				fs.CustomerList[j].HealthScore += 0.15
			}
		}

		if atRisk > 0 {
			messages = append(messages, fmt.Sprintf("✅ CS: Saved %d at-risk customers | NPS: %d", atRisk, playbook.NPSScore))
		}
	}

	return messages
}

