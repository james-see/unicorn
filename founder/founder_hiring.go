package founder

import (
	"fmt"
	"math/rand"
)

// ExecOffer represents a compensation package for an executive candidate
type ExecOffer struct {
	Name        string
	Role        EmployeeRole
	Equity      float64 // Equity percentage
	MonthlyCost int64   // Monthly salary
	AnnualCost  int64   // Annual salary (for display)
	Impact      float64 // Effectiveness multiplier
	Label       string  // "Standard", "Equity-Heavy", "Cash-Heavy"
	Description string  // Short description of the tradeoff
}

// GenerateExecOffers creates 3 compensation packages for an executive hire
func (fs *FounderState) GenerateExecOffers(role EmployeeRole) ([]ExecOffer, error) {
	// Check if we already have this executive
	for _, exec := range fs.Team.Executives {
		if exec.Role == role {
			return nil, fmt.Errorf("already have a %s", role)
		}
	}

	availableEquity := fs.EquityPool - fs.EquityAllocated
	if availableEquity < 0 {
		availableEquity = 0
	}

	if availableEquity < 0.5 {
		return nil, fmt.Errorf("insufficient equity pool (%.1f%% available) — expand pool via Board & Equity", availableEquity)
	}

	// Famous C-suite names from Silicon Valley (show & real life)
	execNames := map[EmployeeRole][]string{
		RoleCTO: {"Gilfoyle", "Steve Wozniak", "Sergey Brin", "Marc Andreessen", "Brendan Eich"},
		RoleCFO: {"Jared Dunn", "Ruth Porat", "David Wehner", "Ned Segal", "Luca Maestri"},
		RoleCOO: {"Sheryl Sandberg", "Gwart", "Tim Cook", "Jeff Weiner", "Stephanie McMahon"},
		RoleCGO: {"Richard Hendricks", "Erlich Bachman", "Andrew Chen", "Alex Schultz", "Sean Ellis"},
	}

	name := execNames[role][rand.Intn(len(execNames[role]))]
	baseImpact := 3.0 * (0.8 + rand.Float64()*0.4) // 2.4-3.6x

	// Generate 3 offers with different equity/cash tradeoffs
	// Standard: balanced equity and salary
	standardEquity := 2.0 + rand.Float64()*2.0 // 2-4%
	if standardEquity > availableEquity {
		standardEquity = availableEquity
	}
	standardSalary := int64(300000) // $300k/year

	// Equity-heavy: more equity, lower salary — exec is betting on the company
	equityHeavyEquity := standardEquity * 1.5
	if equityHeavyEquity > availableEquity {
		equityHeavyEquity = availableEquity
	}
	equityHeavySalary := int64(200000) // $200k/year

	// Cash-heavy: minimal equity, higher salary — exec wants guaranteed comp
	cashHeavyEquity := standardEquity * 0.4
	if cashHeavyEquity < 0.5 {
		cashHeavyEquity = 0.5
	}
	cashHeavySalary := int64(450000) // $450k/year

	offers := []ExecOffer{
		{
			Name:        name,
			Role:        role,
			Equity:      standardEquity,
			MonthlyCost: standardSalary / 12,
			AnnualCost:  standardSalary,
			Impact:      baseImpact,
			Label:       "Standard",
			Description: "Balanced equity and salary",
		},
		{
			Name:        name,
			Role:        role,
			Equity:      equityHeavyEquity,
			MonthlyCost: equityHeavySalary / 12,
			AnnualCost:  equityHeavySalary,
			Impact:      baseImpact * 1.1, // Slightly more motivated — skin in the game
			Label:       "Equity-Heavy",
			Description: "More equity, lower salary — believes in the mission",
		},
		{
			Name:        name,
			Role:        role,
			Equity:      cashHeavyEquity,
			MonthlyCost: cashHeavySalary / 12,
			AnnualCost:  cashHeavySalary,
			Impact:      baseImpact * 0.9, // Slightly less invested
			Label:       "Cash-Heavy",
			Description: "Minimal equity, higher salary — wants guaranteed comp",
		},
	}

	return offers, nil
}

// HireExecWithOffer hires an executive using a selected offer package
func (fs *FounderState) HireExecWithOffer(offer ExecOffer) error {
	// Validate equity is still available
	availableEquity := fs.EquityPool - fs.EquityAllocated
	if availableEquity < 0 {
		availableEquity = 0
	}
	if offer.Equity > availableEquity {
		return fmt.Errorf("insufficient equity pool (need %.1f%%, have %.1f%% available)", offer.Equity, availableEquity)
	}

	// Check duplicate
	for _, exec := range fs.Team.Executives {
		if exec.Role == offer.Role {
			return fmt.Errorf("already have a %s", offer.Role)
		}
	}

	fs.EquityAllocated += offer.Equity

	employee := Employee{
		Name:          offer.Name,
		Role:          offer.Role,
		MonthlyCost:   offer.MonthlyCost,
		Impact:        offer.Impact,
		IsExecutive:   true,
		Equity:        offer.Equity,
		VestingMonths: 48, // 4 year vesting
		CliffMonths:   12, // 1 year cliff
		VestedMonths:  0,
		HasCliff:      false,
		MonthHired:    fs.Turn,
	}
	fs.Team.Executives = append(fs.Team.Executives, employee)

	fs.CapTable = append(fs.CapTable, CapTableEntry{
		Name:         employee.Name,
		Type:         "executive",
		Equity:       offer.Equity,
		MonthGranted: fs.Turn,
	})

	fs.CalculateTeamCost()
	fs.CalculateRunway()

	if offer.Role == RoleCustomerSuccess || offer.Role == RoleCOO {
		fs.RecalculateChurnRate()
	}

	return nil
}

func (fs *FounderState) HireEmployee(role EmployeeRole) error {
	avgSalary := int64(100000)
	var employee Employee

	// C-level executives should use GenerateExecOffers + HireExecWithOffer
	isExec := (role == RoleCTO || role == RoleCFO || role == RoleCOO || role == RoleCGO)

	if isExec {
		// Fallback for non-TUI callers: auto-select standard offer
		offers, err := fs.GenerateExecOffers(role)
		if err != nil {
			return err
		}
		return fs.HireExecWithOffer(offers[0]) // Use standard offer
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