package founder

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"

	"github.com/jamesacampbell/unicorn/assets"
)

// LoadFounderStartups loads startup templates from embedded assets
// The filename parameter is kept for backwards compatibility but is ignored
func LoadFounderStartups(filename string) ([]StartupTemplate, error) {
	data, err := assets.ReadFounderStartups()
	if err != nil {
		return nil, fmt.Errorf("failed to read startups.json: %v", err)
	}

	var templates []StartupTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, fmt.Errorf("failed to parse startups.json: %v", err)
	}

	return templates, nil
}

func generateDealSize(avgDealSize int64, category string) int64 {
	// If avgDealSize is 0 (no customers yet), use category-based defaults
	if avgDealSize == 0 {
		switch category {
		case "SaaS":
			avgDealSize = 1000 // Default $1k/month for SaaS
		case "DeepTech":
			avgDealSize = 5000 // Default $5k/month for DeepTech
		case "GovTech":
			avgDealSize = 2000 // Default $2k/month for GovTech
		case "Hardware":
			avgDealSize = 3000 // Default $3k/month for Hardware
		default:
			avgDealSize = 1000 // Default $1k/month
		}
	}

	// Use normal distribution approximation: 70% chance within ±30%, 30% chance outside
	// Most deals cluster around avg, but some are much larger/smaller
	var dealSize int64

	if rand.Float64() < 0.7 {
		// 70% of deals: within ±30% of avg
		variation := -0.3 + rand.Float64()*0.6 // -30% to +30%
		dealSize = int64(float64(avgDealSize) * (1.0 + variation))
	} else {
		// 30% of deals: wider range
		if rand.Float64() < 0.5 {
			// Smaller deals: 50% to 70% of avg
			variation := 0.5 + rand.Float64()*0.2
			dealSize = int64(float64(avgDealSize) * variation)
		} else {
			// Larger deals: 130% to 200% of avg (enterprise deals)
			variation := 1.3 + rand.Float64()*0.7
			dealSize = int64(float64(avgDealSize) * variation)
		}
	}

	// Ensure minimum deal size (at least 10% of avg)
	if dealSize < avgDealSize/10 {
		dealSize = avgDealSize / 10
	}

	// Apply category-based maximum caps (per month)
	var maxDealSize int64
	switch category {
	case "SaaS":
		// Software: Max $83k/month (million dollar annual license)
		maxDealSize = 83000
	case "DeepTech":
		// Hardware/Physical: Max $100k/month (could be per unit for expensive items)
		maxDealSize = 100000
	case "GovTech":
		// Government contracts: Can be larger, but cap at $200k/month
		maxDealSize = 200000
	case "Hardware":
		// Physical hardware: Max $100k/month per unit
		maxDealSize = 100000
	default:
		// Default cap: $100k/month for unknown types
		maxDealSize = 100000
	}

	// Cap the deal size at the maximum
	if dealSize > maxDealSize {
		dealSize = maxDealSize
	}

	return dealSize
}

func NewFounderGame(founderName string, template StartupTemplate, playerUpgrades []string) *FounderState {
	fs := &FounderState{
		FounderName:        founderName,
		CompanyName:        template.Name,
		Category:           template.Type,
		StartupType:        template.Type,
		Description:        template.Description,
		Cash:               template.InitialCash,
		MRR:                template.InitialMRR,
		DirectMRR:          template.InitialMRR, // Initial MRR is all direct
		AffiliateMRR:       0,
		Customers:          template.InitialCustomers,
		DirectCustomers:    template.InitialCustomers, // Initial customers are all direct
		AffiliateCustomers: 0,
		AvgDealSize:        template.AvgDealSize,
		MinDealSize:        template.AvgDealSize, // Will be updated as deals are made
		MaxDealSize:        template.AvgDealSize, // Will be updated as deals are made
		ChurnRate:          template.BaseChurnRate,
		CustomerChurnRate:  template.BaseChurnRate,
		BaseCAC:            template.BaseCAC,
		Turn:               1,
		MaxTurns:           60,                         // 5 years
		ProductMaturity:    0.20 + rand.Float64()*0.25, // Randomize between 20-45% product maturity
		MarketPenetration:  float64(template.InitialCustomers) / float64(template.TargetMarketSize),
		TargetMarketSize:   template.TargetMarketSize,
		CompetitionLevel:   template.CompetitionLevel,

		// Initialize advanced features
		Partnerships:       []Partnership{},
		AffiliateProgram:   nil,
		ReferralProgram:    nil,
		Competitors:        []Competitor{},
		GlobalMarkets:      []Market{},
		PivotHistory:       []Pivot{},
		EquityPool:         20.0, // Start with 20% equity pool for employees (industry standard)
		EquityAllocated:    0.0,  // Track allocated equity separately
		InvestorBuybacks:   []Buyback{},
		RandomEvents:       []RandomEvent{},
		ActiveEventEffects: make(map[string]EventImpact),
		EquityGivenAway:    0.0,
		BoardSeats:         1,               // Founder starts with 1 board seat
		BoardMembers:       []BoardMember{}, // Start with no advisors
		MonthlyGrowthRate:  0.10,            // Start with 10% monthly growth
		FounderSalary:      12500,           // $150k/year
		CapTable:           []CapTableEntry{},

		// Initialize infrastructure costs
		MonthlyComputeCost: 0,
		MonthlyODCCost:     0,

		// Initialize customer tracking
		CustomerList:       []Customer{},
		TotalCustomersEver: template.InitialCustomers,
		TotalChurned:       0,
		NextCustomerID:     1,

		// Initialize board tracking
		BoardSentiment:     "",
		BoardPressure:      0,
		LastBoardMeeting:   0,
		PendingOpportunity: nil,

		HasExited:     false,
		ExitType:      "",
		ExitValuation: 0,
		ExitMonth:     0,

		MonthReachedProfitability: -1, // -1 means not yet profitable
		PlayerUpgrades:            playerUpgrades,
		HiresCount:                0,

		// Roadmap tracking for achievements
		CustomersLostDuringRoadmap: 0,
	}

	// Add randomness to initial cash (±20%) - only if not already randomized
	// (If randomized in UI, use that value; otherwise randomize here)
	if fs.Cash == template.InitialCash {
		cashVariance := 0.20
		cashMultiplier := 1.0 + (rand.Float64()*cashVariance*2 - cashVariance) // 0.8 to 1.2
		fs.Cash = int64(float64(fs.Cash) * cashMultiplier)
	}

	// Competition level should already be randomized from UI, but if not, randomize it
	competitionLevels := []string{"low", "medium", "high", "very_high"}
	if fs.CompetitionLevel == template.CompetitionLevel {
		// Only randomize if it's still the original template value
		fs.CompetitionLevel = competitionLevels[rand.Intn(len(competitionLevels))]
	}

	// Calculate initial churn based on product maturity
	// Lower maturity = higher churn
	// Formula: baseChurn = (1.0 - ProductMaturity) * 0.65 + 0.05
	// At 25% maturity: (1.0 - 0.25) * 0.65 + 0.05 = 0.5375 ≈ 54% churn
	// At 45% maturity: (1.0 - 0.45) * 0.65 + 0.05 = 0.4075 ≈ 41% churn
	// At 100% maturity: (1.0 - 1.0) * 0.65 + 0.05 = 0.05 = 5% churn (minimum)
	baseChurnFromMaturity := (1.0-fs.ProductMaturity)*0.65 + 0.05
	// Add some variation (±10%)
	churnVariation := -0.10 + rand.Float64()*0.20 // -10% to +10%
	fs.CustomerChurnRate = baseChurnFromMaturity * (1.0 + churnVariation)
	fs.ChurnRate = fs.CustomerChurnRate

	// Apply upgrades
	for _, upgradeID := range playerUpgrades {
		switch upgradeID {
		case "fast_track":
			// Start with 10% more product maturity
			fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+0.1)
			// Recalculate churn after maturity boost
			baseChurnFromMaturity = (1.0-fs.ProductMaturity)*0.65 + 0.05
			fs.CustomerChurnRate = baseChurnFromMaturity * (1.0 + churnVariation)
		case "sales_boost":
			// +15% to initial MRR
			fs.MRR = int64(float64(fs.MRR) * 1.15)
			fs.DirectMRR = fs.MRR
		case "churn_shield":
			// Reduce churn by 10% permanently
			fs.ChurnRate *= 0.9
			fs.CustomerChurnRate = fs.ChurnRate
		}
	}

	// Apply Churn Shield if active (after all maturity calculations)
	for _, upgradeID := range playerUpgrades {
		if upgradeID == "churn_shield" {
			fs.CustomerChurnRate *= 0.9
			fs.ChurnRate = fs.CustomerChurnRate
			break
		}
	}

	// Calculate initial effective CAC
	fs.UpdateCAC()

	// Initialize team from template
	fs.Team = Team{
		Engineers:       make([]Employee, template.InitialTeam["engineers"]),
		Sales:           make([]Employee, template.InitialTeam["sales"]),
		CustomerSuccess: make([]Employee, template.InitialTeam["customer_success"]),
		Marketing:       make([]Employee, template.InitialTeam["marketing"]),
	}

	// Set up initial employees with equity grants (1-2% each)
	avgSalary := int64(100000)
	totalInitialEmployees := template.InitialTeam["engineers"] + template.InitialTeam["sales"] +
		template.InitialTeam["customer_success"] + template.InitialTeam["marketing"]

	// Calculate total equity for initial employees (0.5-1.5% each)
	equityPerEmployee := 0.5 + rand.Float64()*1.0 // 0.5-1.5% per employee
	totalEmployeeEquity := float64(totalInitialEmployees) * equityPerEmployee

	// Ensure we don't exceed equity pool
	if totalEmployeeEquity > fs.EquityPool {
		equityPerEmployee = fs.EquityPool / float64(totalInitialEmployees)
		totalEmployeeEquity = fs.EquityPool
	}

	// Track allocated equity (deducted from founder's ownership)
	fs.EquityAllocated = totalEmployeeEquity

	// Generate employee names
	engineerNames := []string{"Alex Chen", "Jordan Smith", "Taylor Kim", "Morgan Lee", "Casey Park"}
	salesNames := []string{"Sam Rivera", "Riley Cooper", "Jamie Foster", "Drew Mitchell", "Quinn Baker"}
	csNames := []string{"Avery Thompson", "Dakota Reed", "Skylar Hayes", "River Martinez", "Phoenix Adams"}
	marketingNames := []string{"Sage Wilson", "Rowan Clark", "Eden Wright", "Blake Turner", "Finley Ross"}

	employeeIdx := 0
	for i := range fs.Team.Engineers {
		name := engineerNames[i%len(engineerNames)]
		fs.Team.Engineers[i] = Employee{
			Name:          name,
			Role:          RoleEngineer,
			MonthlyCost:   avgSalary / 12,
			Impact:        0.8 + rand.Float64()*0.4, // 0.8-1.2x impact
			IsExecutive:   false,
			Equity:        equityPerEmployee,
			VestingMonths: 48, // 4 year vesting
			CliffMonths:   12, // 1 year cliff
			VestedMonths:  0,
			HasCliff:      false,
			MonthHired:    1,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         name,
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}
	for i := range fs.Team.Sales {
		name := salesNames[i%len(salesNames)]
		fs.Team.Sales[i] = Employee{
			Name:          name,
			Role:          RoleSales,
			MonthlyCost:   avgSalary / 12,
			Impact:        0.8 + rand.Float64()*0.4,
			IsExecutive:   false,
			Equity:        equityPerEmployee,
			VestingMonths: 48,
			CliffMonths:   12,
			VestedMonths:  0,
			HasCliff:      false,
			MonthHired:    1,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         name,
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}
	for i := range fs.Team.CustomerSuccess {
		name := csNames[i%len(csNames)]
		fs.Team.CustomerSuccess[i] = Employee{
			Name:          name,
			Role:          RoleCustomerSuccess,
			MonthlyCost:   avgSalary / 12,
			Impact:        0.8 + rand.Float64()*0.4,
			IsExecutive:   false,
			Equity:        equityPerEmployee,
			VestingMonths: 48,
			CliffMonths:   12,
			VestedMonths:  0,
			HasCliff:      false,
			MonthHired:    1,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         name,
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}
	for i := range fs.Team.Marketing {
		name := marketingNames[i%len(marketingNames)]
		fs.Team.Marketing[i] = Employee{
			Name:          name,
			Role:          RoleMarketing,
			MonthlyCost:   avgSalary / 12,
			Impact:        0.8 + rand.Float64()*0.4,
			IsExecutive:   false,
			Equity:        equityPerEmployee,
			VestingMonths: 48,
			CliffMonths:   12,
			VestedMonths:  0,
			HasCliff:      false,
			MonthHired:    1,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         name,
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}

	fs.CalculateTeamCost()
	fs.CalculateRunway()

	// Initialize new advanced features
	InitializeAcquisitions(fs)
	InitializePlatform(fs)
	InitializePartnershipIntegrations(fs)
	InitializeSecurity(fs)
	InitializeKeyPersonRisks(fs)

	// Initialize customer segments and verticals so ICP/pricing views work from the start
	fs.InitializeSegments()
	fs.InitializeVerticals()
	fs.UpdateSegmentVolumes()

	return fs
}

func formatCurrency(amount int64) string {
	if amount < 0 {
		return fmt.Sprintf("-$%s", formatCurrency(-amount))
	}

	str := fmt.Sprintf("%d", amount)
	var result string
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}
