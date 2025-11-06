package founder

import (
	"fmt"
	"math/rand"
)


func (fs *FounderState) HireEmployee(role EmployeeRole) error {
	avgSalary := int64(100000)
	var employee Employee

	// C-level executives cost $300k and have 3x impact
	isExec := (role == RoleCTO || role == RoleCFO || role == RoleCOO || role == RoleCGO)

	if isExec {
		// Check if we already have this executive
		for _, exec := range fs.Team.Executives {
			if exec.Role == role {
				return fmt.Errorf("already have a %s", role)
			}
		}

		// C-suite executives get 3-10% equity from the pool
		executiveEquity := 3.0 + rand.Float64()*7.0 // 3-10% equity

		// Check if we have enough equity pool available
		availableEquity := fs.EquityPool - fs.EquityAllocated
		if executiveEquity > availableEquity {
			return fmt.Errorf("insufficient equity pool (need %.1f%%, have %.1f%% available)", executiveEquity, availableEquity)
		}

		fs.EquityAllocated += executiveEquity

		// Famous C-suite names from Silicon Valley (show & real life)
		execNames := map[EmployeeRole][]string{
			RoleCTO: {"Gilfoyle", "Steve Wozniak", "Sergey Brin", "Marc Andreessen", "Brendan Eich"},
			RoleCFO: {"Jared Dunn", "Ruth Porat", "David Wehner", "Ned Segal", "Luca Maestri"},
			RoleCOO: {"Sheryl Sandberg", "Gwart", "Tim Cook", "Jeff Weiner", "Stephanie McMahon"},
			RoleCGO: {"Richard Hendricks", "Erlich Bachman", "Andrew Chen", "Alex Schultz", "Sean Ellis"},
		}

		employee = Employee{
			Name:          execNames[role][rand.Intn(len(execNames[role]))],
			Role:          role,
			MonthlyCost:   25000,                            // $300k/year
			Impact:        3.0 * (0.8 + rand.Float64()*0.4), // 3x impact (2.4-3.6x)
			IsExecutive:   true,
			Equity:        executiveEquity,
			VestingMonths: 48, // 4 year vesting
			CliffMonths:   12, // 1 year cliff
			VestedMonths:  0,
			HasCliff:      false,
			MonthHired:    fs.Turn,
		}
		fs.Team.Executives = append(fs.Team.Executives, employee)

		// Add to cap table with executive's name
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         employee.Name,
			Type:         "executive",
			Equity:       executiveEquity,
			MonthGranted: fs.Turn,
		})
	} else {
		monthlyCost := avgSalary / 12
		
		// Apply Quick Hire upgrade (first 3 hires cost 50% less)
		hasQuickHire := false
		for _, upgradeID := range fs.PlayerUpgrades {
			if upgradeID == "quick_hire" {
				hasQuickHire = true
				break
			}
		}
		if hasQuickHire && fs.HiresCount < 3 {
			monthlyCost = monthlyCost / 2 // 50% discount
		}
		fs.HiresCount++
		
		employee = Employee{
			Role:           role,
			MonthlyCost:    monthlyCost,
			Impact:         0.8 + rand.Float64()*0.4,
			IsExecutive:    false,
			AssignedMarket: "USA", // Default to USA market
		}

		switch role {
		case RoleEngineer:
			fs.Team.Engineers = append(fs.Team.Engineers, employee)
		case RoleSales:
			fs.Team.Sales = append(fs.Team.Sales, employee)
		case RoleCustomerSuccess:
			fs.Team.CustomerSuccess = append(fs.Team.CustomerSuccess, employee)
		case RoleMarketing:
			fs.Team.Marketing = append(fs.Team.Marketing, employee)
		default:
			return fmt.Errorf("unknown role: %s", role)
		}
	}

	fs.CalculateTeamCost()
	fs.CalculateRunway()

	// Recalculate churn rate if hiring CS or COO (affects churn)
	if role == RoleCustomerSuccess || role == RoleCOO {
		fs.RecalculateChurnRate()
	}

	return nil
}


func (fs *FounderState) HireEmployeeWithMarket(role EmployeeRole, market string) error {
	avgSalary := int64(100000)

	employee := Employee{
		Role:           role,
		MonthlyCost:    avgSalary / 12,
		Impact:         0.8 + rand.Float64()*0.4,
		IsExecutive:    false,
		AssignedMarket: market,
		MonthHired:     fs.Turn,
	}

	switch role {
	case RoleEngineer:
		fs.Team.Engineers = append(fs.Team.Engineers, employee)
	case RoleSales:
		fs.Team.Sales = append(fs.Team.Sales, employee)
	case RoleCustomerSuccess:
		fs.Team.CustomerSuccess = append(fs.Team.CustomerSuccess, employee)
	case RoleMarketing:
		fs.Team.Marketing = append(fs.Team.Marketing, employee)
	default:
		return fmt.Errorf("unknown role: %s", role)
	}

	fs.CalculateTeamCost()
	fs.CalculateRunway()

	// Recalculate churn rate if hiring CS or COO (affects churn)
	if role == RoleCustomerSuccess || role == RoleCOO {
		fs.RecalculateChurnRate()
	}

	return nil
}

// FireEmployee removes a team member

func (fs *FounderState) FireEmployee(role EmployeeRole) error {
	// Check if it's an executive role
	isExec := (role == RoleCTO || role == RoleCFO || role == RoleCOO || role == RoleCGO)

	if isExec {
		for i, exec := range fs.Team.Executives {
			if exec.Role == role {
				fs.Team.Executives = append(fs.Team.Executives[:i], fs.Team.Executives[i+1:]...)
				fs.CalculateTeamCost()
				fs.CalculateRunway()

				// Recalculate churn rate if firing COO (affects churn)
				if role == RoleCOO {
					fs.RecalculateChurnRate()
				}

				return nil
			}
		}
		return fmt.Errorf("don't have a %s to let go", role)
	}

	switch role {
	case RoleEngineer:
		if len(fs.Team.Engineers) > 0 {
			fs.Team.Engineers = fs.Team.Engineers[:len(fs.Team.Engineers)-1]
		} else {
			return fmt.Errorf("no engineers to fire")
		}
	case RoleSales:
		if len(fs.Team.Sales) > 0 {
			fs.Team.Sales = fs.Team.Sales[:len(fs.Team.Sales)-1]
		} else {
			return fmt.Errorf("no sales reps to fire")
		}
	case RoleCustomerSuccess:
		if len(fs.Team.CustomerSuccess) > 0 {
			fs.Team.CustomerSuccess = fs.Team.CustomerSuccess[:len(fs.Team.CustomerSuccess)-1]
		} else {
			return fmt.Errorf("no CS reps to fire")
		}
	case RoleMarketing:
		if len(fs.Team.Marketing) > 0 {
			fs.Team.Marketing = fs.Team.Marketing[:len(fs.Team.Marketing)-1]
		} else {
			return fmt.Errorf("no marketers to fire")
		}
	default:
		return fmt.Errorf("unknown role: %s", role)
	}

	fs.CalculateTeamCost()
	fs.CalculateRunway()

	// Recalculate churn rate if firing CS (affects churn)
	if role == RoleCustomerSuccess {
		fs.RecalculateChurnRate()
	}

	return nil
}

// GenerateTermSheetOptions creates multiple term sheet options for a funding round

func (fs *FounderState) UpdateEmployeeVesting() {
	updateVesting := func(employees *[]Employee) {
		for i := range *employees {
			e := &(*employees)[i]
			if e.Equity > 0 {
				e.VestedMonths = fs.Turn - e.MonthHired
				if e.VestedMonths >= e.CliffMonths && !e.HasCliff {
					e.HasCliff = true // Cliff reached!
				}
			}
		}
	}

	updateVesting(&fs.Team.Engineers)
	updateVesting(&fs.Team.Sales)
	updateVesting(&fs.Team.CustomerSuccess)
	updateVesting(&fs.Team.Marketing)
	updateVesting(&fs.Team.Executives)
}

// ProcessMonth runs all monthly calculations