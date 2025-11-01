package founder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
)

// EmployeeRole represents different types of employees
type EmployeeRole string

const (
	RoleEngineer        EmployeeRole = "engineer"
	RoleSales           EmployeeRole = "sales"
	RoleCustomerSuccess EmployeeRole = "customer_success"
	RoleMarketing       EmployeeRole = "marketing"
	RoleCTO             EmployeeRole = "cto"
	RoleCFO             EmployeeRole = "cfo"
	RoleCOO             EmployeeRole = "coo"
	RoleCGO             EmployeeRole = "cgo" // Chief Growth Officer (sales/marketing)
)

// Employee represents a team member
type Employee struct {
	Role        EmployeeRole
	MonthlyCost int64
	Impact      float64 // Productivity/effectiveness multiplier
	IsExecutive bool    // C-level executives have 3x impact, $300k/year salary
	Equity      float64 // Equity percentage owned by this employee
}

// CapTableEntry tracks individual equity ownership
type CapTableEntry struct {
	Name         string  // Employee name or investor round name
	Type         string  // "employee", "executive", "investor"
	Equity       float64 // Equity percentage
	MonthGranted int     // Month when equity was granted
}

// Team tracks all employees
type Team struct {
	Engineers        []Employee
	Sales            []Employee
	CustomerSuccess  []Employee
	Marketing        []Employee
	Executives       []Employee // C-level: CTO, CFO, COO, CGO
	TotalMonthlyCost int64
	TotalEmployees   int
}

// StartupTemplate represents a startup idea from JSON
type StartupTemplate struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Tagline          string         `json:"tagline"`
	Type             string         `json:"type"`
	Description      string         `json:"description"`
	InitialCash      int64          `json:"initial_cash"`
	MonthlyBurn      int64          `json:"monthly_burn"`
	InitialCustomers int            `json:"initial_customers"`
	InitialMRR       int64          `json:"initial_mrr"`
	AvgDealSize      int64          `json:"avg_deal_size"`
	BaseChurnRate    float64        `json:"base_churn_rate"`
	BaseCAC          int64          `json:"base_cac"`
	TargetMarketSize int            `json:"target_market_size"`
	CompetitionLevel string         `json:"competition_level"`
	InitialTeam      map[string]int `json:"initial_team"`
}

// FounderState represents the current state of your startup
type FounderState struct {
	FounderName       string
	CompanyName       string
	Category          string
	StartupType       string
	Description       string
	Cash              int64
	MRR               int64 // Monthly Recurring Revenue
	Customers         int
	AvgDealSize       int64
	ChurnRate         float64
	CustomerChurnRate float64 // Alias for ChurnRate
	BaseCAC           int64   // Base customer acquisition cost for this business
	Team              Team
	Turn              int
	MaxTurns          int
	ProductMaturity   float64 // 0-1, affects sales velocity
	MarketPenetration float64 // 0-1, % of target market captured
	TargetMarketSize  int
	CompetitionLevel  string
	FundingRounds     []FundingRound
	EquityGivenAway   float64 // Total % equity given to investors
	BoardSeats        int     // Board seats given to investors
	AcquisitionOffers []AcquisitionOffer
	CashRunwayMonths  int
	MonthlyTeamCost   int64 // Cached monthly team cost
	FounderSalary     int64 // $150k/year = $12,500/month

	// Growth metrics
	MonthlyGrowthRate       float64
	CustomerAcquisitionCost int64 // Current effective CAC (changes based on maturity)
	LifetimeValue           int64

	// Advanced features
	Partnerships       []Partnership
	AffiliateProgram   *AffiliateProgram
	Competitors        []Competitor
	GlobalMarkets      []Market
	PivotHistory       []Pivot
	EquityPool         float64 // Employee equity pool %
	InvestorBuybacks   []Buyback
	RandomEvents       []RandomEvent
	ActiveEventEffects map[string]EventImpact // Events currently affecting the business
	CapTable           []CapTableEntry        // Individual equity ownership tracking
}

// FundingRound represents a completed fundraise
type FundingRound struct {
	RoundName   string
	Amount      int64
	Valuation   int64
	EquityGiven float64
	Month       int
	Terms       string // "Founder-friendly", "Standard", "Investor-heavy"
}

// TermSheetOption represents different fundraising options to choose from
type TermSheetOption struct {
	Amount        int64
	PostValuation int64
	PreValuation  int64
	Equity        float64
	Terms         string
	Description   string
}

// AcquisitionOffer represents an offer to buy the company
type AcquisitionOffer struct {
	Acquirer     string
	OfferAmount  int64
	Month        int
	DueDiligence string // "bad", "normal", "good"
	TermsQuality string // "poor", "good", "excellent"
}

// Partnership represents a strategic partnership
type Partnership struct {
	Partner        string
	Type           string // "distribution", "technology", "co-marketing", "data"
	MonthStarted   int
	Duration       int // Months
	Cost           int64
	MRRBoost       int64
	ChurnReduction float64
	Status         string // "active", "expired"
}

// AffiliateProgram represents an affiliate sales program
type AffiliateProgram struct {
	LaunchedMonth      int
	Commission         float64 // % of deal
	Affiliates         int
	SetupCost          int64
	MonthlyPlatformFee int64
	MonthlyRevenue     int64
	CustomersAcquired  int
}

// Competitor represents a competing company
type Competitor struct {
	Name          string
	Threat        string // "low", "medium", "high", "critical"
	MarketShare   float64
	Strategy      string // "ignore", "monitor", "compete", "partner"
	MonthAppeared int
	Active        bool
}

// Market represents a geographic expansion
type Market struct {
	Region           string // "North America", "Europe", "Asia", "LATAM", etc.
	LaunchMonth      int
	SetupCost        int64
	MonthlyCost      int64
	CustomerCount    int
	MRR              int64
	MarketSize       int
	Penetration      float64
	LocalCompetition string
}

// Pivot represents a strategy or market change
type Pivot struct {
	Month         int
	FromStrategy  string
	ToStrategy    string
	Reason        string
	Cost          int64
	CustomersLost int
	Success       bool
}

// Buyback represents buying back equity from investors
type Buyback struct {
	Month        int
	Investor     string // Which round (Seed, Series A, etc)
	EquityBought float64
	PricePaid    int64
	Valuation    int64
}

// RandomEvent represents a random occurrence that affects the business
type RandomEvent struct {
	Month       int
	Type        string // "economy", "regulation", "competition", "talent", "customer", "product", "legal", "press"
	Severity    string // "minor", "moderate", "major", "critical"
	IsPositive  bool
	Title       string
	Description string
	Impact      EventImpact
}

// EventImpact describes the effects of an event
type EventImpact struct {
	CACChange          float64 // Multiplier (1.2 = +20%, 0.8 = -20%)
	ChurnChange        float64 // Additive (0.05 = +5%, -0.02 = -2%)
	GrowthChange       float64 // Multiplier
	CashCost           int64   // One-time cost
	MRRChange          float64 // Multiplier
	EmployeesLost      int     // Number of employees who quit
	ProductivityChange float64 // Team productivity multiplier
	DurationMonths     int     // How long the effect lasts
}

// Decision represents a choice the founder can make
type Decision struct {
	Type        string
	Description string
	Cost        int64
	Impact      string
}

// LoadFounderStartups loads all startup ideas from JSON
func LoadFounderStartups(filename string) ([]StartupTemplate, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open startups.json: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read startups.json: %v", err)
	}

	var templates []StartupTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, fmt.Errorf("failed to parse startups.json: %v", err)
	}

	return templates, nil
}

// NewFounderGame initializes a new founder mode game
func NewFounderGame(founderName string, template StartupTemplate) *FounderState {
	fs := &FounderState{
		FounderName:       founderName,
		CompanyName:       template.Name,
		Category:          template.Type,
		StartupType:       template.Type,
		Description:       template.Description,
		Cash:              template.InitialCash,
		MRR:               template.InitialMRR,
		Customers:         template.InitialCustomers,
		AvgDealSize:       template.AvgDealSize,
		ChurnRate:         template.BaseChurnRate,
		CustomerChurnRate: template.BaseChurnRate,
		BaseCAC:           template.BaseCAC,
		Turn:              1,
		MaxTurns:          60,  // 5 years
		ProductMaturity:   0.3, // Start at 30% product maturity
		MarketPenetration: float64(template.InitialCustomers) / float64(template.TargetMarketSize),
		TargetMarketSize:  template.TargetMarketSize,
		CompetitionLevel:  template.CompetitionLevel,

		// Initialize advanced features
		Partnerships:       []Partnership{},
		AffiliateProgram:   nil,
		Competitors:        []Competitor{},
		GlobalMarkets:      []Market{},
		PivotHistory:       []Pivot{},
		EquityPool:         10.0, // Start with 10% equity pool for employees
		InvestorBuybacks:   []Buyback{},
		RandomEvents:       []RandomEvent{},
		ActiveEventEffects: make(map[string]EventImpact),
		EquityGivenAway:    0.0,
		BoardSeats:         1,     // Founder starts with 1 board seat
		MonthlyGrowthRate:  0.10,  // Start with 10% monthly growth
		FounderSalary:      12500, // $150k/year
		CapTable:           []CapTableEntry{},
	}

	// Add randomness to initial cash (±20%)
	cashVariance := 0.20
	cashMultiplier := 1.0 + (rand.Float64()*cashVariance*2 - cashVariance) // 0.8 to 1.2
	fs.Cash = int64(float64(fs.Cash) * cashMultiplier)

	// Add randomness to competition level
	competitionLevels := []string{"low", "medium", "high", "very_high"}
	if rand.Float64() < 0.3 { // 30% chance to randomize competition level
		fs.CompetitionLevel = competitionLevels[rand.Intn(len(competitionLevels))]
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

	// Calculate total equity for initial employees (1-2% each)
	equityPerEmployee := 1.0 + rand.Float64()*1.0 // 1-2% per employee
	totalEmployeeEquity := float64(totalInitialEmployees) * equityPerEmployee

	// Ensure we don't exceed equity pool
	if totalEmployeeEquity > fs.EquityPool {
		equityPerEmployee = fs.EquityPool / float64(totalInitialEmployees)
		totalEmployeeEquity = fs.EquityPool
	}

	fs.EquityPool -= totalEmployeeEquity

	employeeIdx := 0
	for i := range fs.Team.Engineers {
		fs.Team.Engineers[i] = Employee{
			Role:        RoleEngineer,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4, // 0.8-1.2x impact
			IsExecutive: false,
			Equity:      equityPerEmployee,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         fmt.Sprintf("Engineer %d", employeeIdx+1),
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}
	for i := range fs.Team.Sales {
		fs.Team.Sales[i] = Employee{
			Role:        RoleSales,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
			IsExecutive: false,
			Equity:      equityPerEmployee,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         fmt.Sprintf("Sales Rep %d", employeeIdx+1),
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}
	for i := range fs.Team.CustomerSuccess {
		fs.Team.CustomerSuccess[i] = Employee{
			Role:        RoleCustomerSuccess,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
			IsExecutive: false,
			Equity:      equityPerEmployee,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         fmt.Sprintf("CS Rep %d", employeeIdx+1),
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}
	for i := range fs.Team.Marketing {
		fs.Team.Marketing[i] = Employee{
			Role:        RoleMarketing,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
			IsExecutive: false,
			Equity:      equityPerEmployee,
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         fmt.Sprintf("Marketer %d", employeeIdx+1),
			Type:         "employee",
			Equity:       equityPerEmployee,
			MonthGranted: 1,
		})
		employeeIdx++
	}

	fs.CalculateTeamCost()
	fs.CalculateRunway()

	return fs
}

// CalculateTeamCost calculates total monthly team cost
func (fs *FounderState) CalculateTeamCost() {
	total := fs.FounderSalary // Start with founder salary
	count := 1                // Founder counts as 1 employee

	for _, e := range fs.Team.Engineers {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.Sales {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.CustomerSuccess {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.Marketing {
		total += e.MonthlyCost
		count++
	}
	for _, e := range fs.Team.Executives {
		total += e.MonthlyCost
		count++
	}

	fs.Team.TotalMonthlyCost = total
	fs.Team.TotalEmployees = count
	fs.MonthlyTeamCost = total
}

// CalculateRunway calculates months of runway remaining
func (fs *FounderState) CalculateRunway() {
	monthlyBurn := fs.Team.TotalMonthlyCost + 20000 // Team + $20k ops costs
	monthlyRevenue := fs.MRR
	netBurn := monthlyBurn - monthlyRevenue

	if netBurn <= 0 {
		// Cash positive! Runway is infinite
		fs.CashRunwayMonths = -1
	} else {
		fs.CashRunwayMonths = int(fs.Cash / netBurn)
	}
}

// IsGameOver checks if the game has ended
func (fs *FounderState) IsGameOver() bool {
	return fs.Cash <= 0 || fs.Turn > fs.MaxTurns
}

// GetFinalScore calculates the final outcome
func (fs *FounderState) GetFinalScore() (outcome string, valuation int64, founderEquity float64) {
	founderEquity = 100.0 - fs.EquityGivenAway

	// Calculate final valuation based on MRR
	if fs.MRR > 0 {
		multiple := 10.0 // 10x ARR default
		if fs.MonthlyGrowthRate > 0.20 {
			multiple = 15.0
		}
		if fs.MonthlyGrowthRate < 0 {
			multiple = 5.0
		}
		valuation = int64(float64(fs.MRR) * 12 * multiple)
	}

	// Determine outcome
	if fs.Cash <= 0 {
		outcome = "SHUT DOWN - Ran out of cash"
	} else if fs.Turn > fs.MaxTurns {
		if fs.MRR > 1000000 { // $1M+ MRR
			outcome = "UNICORN PATH - Scaled to $1M+ MRR!"
		} else if fs.MRR > 100000 { // $100K+ MRR
			outcome = "SUCCESS - Built a real business"
		} else if fs.MRR > 0 {
			outcome = "SURVIVING - Making revenue but still growing"
		} else {
			outcome = "STRUGGLING - No meaningful traction"
		}
	}

	return outcome, valuation, founderEquity
}

// HireEmployee adds a new team member
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

		// Check if we have enough equity pool
		if executiveEquity > fs.EquityPool {
			return fmt.Errorf("insufficient equity pool (need %.1f%%, have %.1f%%)", executiveEquity, fs.EquityPool)
		}

		fs.EquityPool -= executiveEquity

		employee = Employee{
			Role:        role,
			MonthlyCost: 25000,                            // $300k/year
			Impact:      3.0 * (0.8 + rand.Float64()*0.4), // 3x impact (2.4-3.6x)
			IsExecutive: true,
			Equity:      executiveEquity,
		}
		fs.Team.Executives = append(fs.Team.Executives, employee)

		// Add to cap table
		execTitle := map[EmployeeRole]string{
			RoleCTO: "CTO",
			RoleCFO: "CFO",
			RoleCOO: "COO",
			RoleCGO: "CGO",
		}
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         execTitle[role],
			Type:         "executive",
			Equity:       executiveEquity,
			MonthGranted: fs.Turn,
		})
	} else {
		employee = Employee{
			Role:        role,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
			IsExecutive: false,
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
	return nil
}

// GenerateTermSheetOptions creates multiple term sheet options for a funding round
func (fs *FounderState) GenerateTermSheetOptions(roundName string) []TermSheetOption {
	// Calculate base valuation based on metrics
	baseValuation := int64(float64(fs.MRR) * 12 * 10) // 10x ARR

	// Adjust based on growth and metrics
	if fs.MonthlyGrowthRate > 0.20 {
		baseValuation = int64(float64(baseValuation) * 1.5)
	}
	if fs.ProductMaturity > 0.7 {
		baseValuation = int64(float64(baseValuation) * 1.2)
	}
	if fs.Customers < 10 {
		baseValuation = int64(float64(baseValuation) * 0.5)
	}

	// Minimum valuations by round
	var minValuation, baseRaise int64
	switch roundName {
	case "Seed":
		minValuation = 3000000
		baseRaise = 2000000
	case "Series A":
		minValuation = 15000000
		baseRaise = 10000000
	case "Series B":
		minValuation = 50000000
		baseRaise = 30000000
	default:
		return []TermSheetOption{}
	}

	if baseValuation < minValuation {
		baseValuation = minValuation
	}

	options := []TermSheetOption{}

	// Option 1: Less money, founder-friendly (lower dilution)
	option1Amount := int64(float64(baseRaise) * 0.7)
	option1PreVal := int64(float64(baseValuation) * 1.1) // 10% higher pre-money
	option1PostVal := option1PreVal + option1Amount
	option1Equity := (float64(option1Amount) / float64(option1PostVal)) * 100
	options = append(options, TermSheetOption{
		Amount:        option1Amount,
		PostValuation: option1PostVal,
		PreValuation:  option1PreVal,
		Equity:        option1Equity,
		Terms:         "Founder-friendly",
		Description:   "Lower dilution, founder-friendly terms, but less capital",
	})

	// Option 2: Standard terms (balanced)
	option2Amount := baseRaise
	option2PreVal := baseValuation
	option2PostVal := option2PreVal + option2Amount
	option2Equity := (float64(option2Amount) / float64(option2PostVal)) * 100
	options = append(options, TermSheetOption{
		Amount:        option2Amount,
		PostValuation: option2PostVal,
		PreValuation:  option2PreVal,
		Equity:        option2Equity,
		Terms:         "Standard",
		Description:   "Fair terms, balanced approach",
	})

	// Option 3: More money, higher dilution
	option3Amount := int64(float64(baseRaise) * 1.4)
	option3PreVal := int64(float64(baseValuation) * 0.9) // 10% lower pre-money
	option3PostVal := option3PreVal + option3Amount
	option3Equity := (float64(option3Amount) / float64(option3PostVal)) * 100
	options = append(options, TermSheetOption{
		Amount:        option3Amount,
		PostValuation: option3PostVal,
		PreValuation:  option3PreVal,
		Equity:        option3Equity,
		Terms:         "Growth-focused",
		Description:   "More capital to scale faster, but higher dilution",
	})

	// Option 4: Maximum money, investor-heavy terms
	option4Amount := int64(float64(baseRaise) * 1.8)
	option4PreVal := int64(float64(baseValuation) * 0.75) // 25% lower pre-money
	option4PostVal := option4PreVal + option4Amount
	option4Equity := (float64(option4Amount) / float64(option4PostVal)) * 100
	options = append(options, TermSheetOption{
		Amount:        option4Amount,
		PostValuation: option4PostVal,
		PreValuation:  option4PreVal,
		Equity:        option4Equity,
		Terms:         "Investor-heavy",
		Description:   "Maximum capital, but significant dilution and investor control",
	})

	return options
}

// RaiseFunding attempts to raise a funding round with a chosen term sheet
func (fs *FounderState) RaiseFundingWithTerms(roundName string, option TermSheetOption) (success bool) {
	fs.Cash += option.Amount
	fs.EquityGivenAway += option.Equity

	round := FundingRound{
		RoundName:   roundName,
		Amount:      option.Amount,
		Valuation:   option.PreValuation,
		EquityGiven: option.Equity,
		Month:       fs.Turn,
		Terms:       option.Terms,
	}
	fs.FundingRounds = append(fs.FundingRounds, round)

	// Add investor to cap table
	fs.CapTable = append(fs.CapTable, CapTableEntry{
		Name:         roundName + " Investors",
		Type:         "investor",
		Equity:       option.Equity,
		MonthGranted: fs.Turn,
	})

	fs.CalculateRunway()

	return true
}

// RaiseFunding is the legacy function kept for backwards compatibility
func (fs *FounderState) RaiseFunding(roundName string) (success bool, amount int64, terms string, equityGiven float64) {
	options := fs.GenerateTermSheetOptions(roundName)
	if len(options) == 0 {
		return false, 0, "", 0
	}

	// Use standard (middle) option
	option := options[1]
	success = fs.RaiseFundingWithTerms(roundName, option)
	return success, option.Amount, option.Terms, option.Equity
}

// UpdateCAC recalculates current effective CAC based on product maturity
func (fs *FounderState) UpdateCAC() {
	// Start with business-specific base CAC
	effectiveCAC := float64(fs.BaseCAC)

	// Product maturity reduces CAC (better product = better conversion)
	// At 100% maturity, CAC is 60% of base (40% reduction)
	maturityDiscount := fs.ProductMaturity * 0.4
	effectiveCAC *= (1.0 - maturityDiscount)

	// Competition increases CAC
	switch fs.CompetitionLevel {
	case "very_high":
		effectiveCAC *= 1.5 // +50%
	case "high":
		effectiveCAC *= 1.3 // +30%
	case "medium":
		effectiveCAC *= 1.1 // +10%
		// low = no change
	}

	fs.CustomerAcquisitionCost = int64(effectiveCAC)
}

// SpendOnMarketing allocates budget to customer acquisition
func (fs *FounderState) SpendOnMarketing(amount int64) int {
	if amount > fs.Cash {
		return 0
	}

	fs.Cash -= amount

	// Use current effective CAC (which accounts for product maturity and competition)
	fs.UpdateCAC()

	newCustomers := int(amount / fs.CustomerAcquisitionCost)
	fs.Customers += newCustomers
	newMRR := int64(newCustomers) * fs.AvgDealSize
	fs.MRR += newMRR

	fs.CalculateRunway()

	return newCustomers
}

// CheckForAcquisition checks if an acquisition offer comes in
func (fs *FounderState) CheckForAcquisition() *AcquisitionOffer {
	// Only after Series A and if metrics are good
	hasSeriesA := false
	for _, round := range fs.FundingRounds {
		if round.RoundName == "Series A" || round.RoundName == "Series B" {
			hasSeriesA = true
			break
		}
	}

	if !hasSeriesA {
		return nil
	}

	// 5% chance per month after Series A
	if rand.Float64() > 0.05 {
		return nil
	}

	// Calculate offer
	multiple := 3.0 + rand.Float64()*3.0 // 3-6x revenue
	annualRevenue := fs.MRR * 12
	offerAmount := int64(float64(annualRevenue) * multiple)

	dueDiligence := "normal"
	termsQuality := "good"

	roll := rand.Float64()
	if roll < 0.15 {
		dueDiligence = "bad"
		termsQuality = "poor"
		offerAmount = int64(float64(offerAmount) * 0.6)
	} else if roll > 0.85 {
		dueDiligence = "good"
		termsQuality = "excellent"
		offerAmount = int64(float64(offerAmount) * 1.3)
	}

	acquirers := []string{
		"Google", "Microsoft", "Amazon", "Salesforce", "Oracle",
		"Meta", "Apple", "Adobe", "SAP", "IBM",
	}

	offer := AcquisitionOffer{
		Acquirer:     acquirers[rand.Intn(len(acquirers))],
		OfferAmount:  offerAmount,
		Month:        fs.Turn,
		DueDiligence: dueDiligence,
		TermsQuality: termsQuality,
	}

	return &offer
}

// NeedsLowCashWarning checks if cash is dangerously low
func (fs *FounderState) NeedsLowCashWarning() bool {
	return fs.Cash <= 200000 && fs.CashRunwayMonths < 6
}

// ProcessMonth runs all monthly calculations
func (fs *FounderState) ProcessMonth() []string {
	var messages []string
	fs.Turn++

	// 1. Process revenue growth
	oldMRR := fs.MRR

	// Engineer impact on product (reduces churn, increases sales)
	// CTO counts as 3x engineers
	engImpact := 1.0
	for _, eng := range fs.Team.Engineers {
		engImpact += (eng.Impact * 0.05) // Each engineer adds ~5% product improvement
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCTO {
			engImpact += (exec.Impact * 0.05) // CTO has 3x impact already built into Impact field
		}
	}
	fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+(0.02*engImpact))

	// Sales team impact on growth
	// CGO counts as 3x sales reps
	salesImpact := 1.0
	for _, sales := range fs.Team.Sales {
		salesImpact += (sales.Impact * 0.1) // Each sales rep adds ~10% to close rate
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCGO {
			salesImpact += (exec.Impact * 0.1) // CGO has 3x impact already built into Impact field
		}
	}

	// Marketing impact (residual from spend)
	baseGrowth := fs.MonthlyGrowthRate
	actualGrowth := baseGrowth * salesImpact * engImpact

	// Apply growth
	newRevenue := int64(float64(fs.MRR) * actualGrowth)
	fs.MRR += newRevenue

	// 2. Process churn (only if we have customers)
	var lostCustomers int
	var actualChurn float64
	if fs.Customers > 0 {
		baseChurn := fs.CustomerChurnRate

		// CS team reduces churn
		// COO counts as 3x CS reps
		csImpact := 0.0
		for _, cs := range fs.Team.CustomerSuccess {
			csImpact += (cs.Impact * 0.02) // Each CS rep reduces churn by ~2%
		}
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCOO {
				csImpact += (exec.Impact * 0.02) // COO has 3x impact already built into Impact field
			}
		}
		actualChurn = math.Max(0.01, baseChurn-csImpact)

		churnLoss := int64(float64(fs.MRR) * actualChurn)
		fs.MRR -= churnLoss
		lostCustomers = int(float64(fs.Customers) * actualChurn)
		fs.Customers -= lostCustomers

		if fs.MRR < 0 {
			fs.MRR = 0
		}
		if fs.Customers < 0 {
			fs.Customers = 0
		}
	}

	// 3. Process MRR cash flow
	// MRR converts to cash after deductions:
	// - Taxes: ~15-25% depending on jurisdiction (use 20% average)
	// - Processing fees: ~3-5% for payment processing (use 3%)
	// - Company overhead: ~5-10% for operations (use 5%)
	// - Savings/buffer: ~5% for reserves
	taxRate := 0.20           // 20% taxes
	processingFeeRate := 0.03 // 3% payment processing
	overheadRate := 0.05      // 5% company overhead
	savingsRate := 0.05       // 5% savings/reserves

	totalDeductionRate := taxRate + processingFeeRate + overheadRate + savingsRate // 33% total deductions
	netMRRToCash := int64(float64(fs.MRR) * (1.0 - totalDeductionRate))
	fs.Cash += netMRRToCash

	// 4. Calculate costs
	totalCost := fs.MonthlyTeamCost + (int64(fs.Team.TotalEmployees) * 2000) // +$2k overhead per employee
	fs.Cash -= totalCost

	netIncome := netMRRToCash - totalCost

	// 5. Update runway
	fs.CalculateRunway()

	// 6. Update growth rate for next month
	if oldMRR > 0 {
		fs.MonthlyGrowthRate = float64(fs.MRR-oldMRR) / float64(oldMRR)
	} else if fs.MRR > 0 {
		// First customers! Set initial growth rate
		fs.MonthlyGrowthRate = 0.10 // Start with 10% base growth
	}

	// 7. Generate messages
	if fs.MRR > 0 && oldMRR == 0 {
		messages = append(messages, fmt.Sprintf("🎉 FIRST REVENUE! MRR: $%s", formatCurrency(fs.MRR)))
	} else if fs.MRR > oldMRR && oldMRR > 0 {
		pctGrowth := ((float64(fs.MRR) - float64(oldMRR)) / float64(oldMRR)) * 100
		messages = append(messages, fmt.Sprintf("💰 MRR grew %.1f%% to $%s", pctGrowth, formatCurrency(fs.MRR)))
	} else if fs.MRR < oldMRR && oldMRR > 0 {
		pctDecline := ((float64(oldMRR) - float64(fs.MRR)) / float64(oldMRR)) * 100
		messages = append(messages, fmt.Sprintf("⚠️  MRR declined %.1f%% to $%s", pctDecline, formatCurrency(fs.MRR)))
	} else if fs.MRR == 0 && fs.Turn > 3 {
		messages = append(messages, "⚠️  Still no revenue! Hire sales or spend on marketing!")
	}

	if lostCustomers > 0 {
		messages = append(messages, fmt.Sprintf("📉 Lost %d customers to churn (%.1f%% churn rate)", lostCustomers, actualChurn*100))
	}

	if netIncome > 0 {
		messages = append(messages, fmt.Sprintf("✅ Positive cash flow: $%s/month", formatCurrency(netIncome)))
	} else {
		messages = append(messages, fmt.Sprintf("💸 Burn rate: $%s/month", formatCurrency(-netIncome)))
	}

	if fs.ProductMaturity >= 1.0 {
		messages = append(messages, "🎉 Product has reached full maturity!")
	}

	// 8. Process advanced features
	partnershipMsgs := fs.UpdatePartnerships()
	messages = append(messages, partnershipMsgs...)

	affiliateMsgs := fs.UpdateAffiliateProgram()
	messages = append(messages, affiliateMsgs...)

	competitorMsgs := fs.UpdateCompetitors()
	messages = append(messages, competitorMsgs...)

	marketMsgs := fs.UpdateGlobalMarkets()
	messages = append(messages, marketMsgs...)

	// 9. Spawn new competitors randomly
	if newComp := fs.SpawnCompetitor(); newComp != nil {
		messages = append(messages, fmt.Sprintf("🚨 NEW COMPETITOR: %s entered the market! Threat: %s, Market Share: %.1f%%",
			newComp.Name, newComp.Threat, newComp.MarketShare*100))
	}

	// 10. Process random events
	eventMsgs := fs.ProcessRandomEvents()
	messages = append(messages, eventMsgs...)

	// 11. Spawn new random events (5% chance each month)
	if rand.Float64() < 0.05 {
		if event := fs.SpawnRandomEvent(); event != nil {
			messages = append(messages, fmt.Sprintf("⚡ EVENT: %s - %s", event.Title, event.Description))
		}
	}

	return messages
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

// ============================================================================
// PARTNERSHIPS
// ============================================================================

// StartPartnership creates a new strategic partnership
func (fs *FounderState) StartPartnership(partnerType string) (*Partnership, error) {
	partners := map[string][]string{
		"distribution": {"Salesforce", "HubSpot", "Oracle", "SAP", "Adobe"},
		"technology":   {"AWS", "Google Cloud", "Microsoft Azure", "IBM", "MongoDB"},
		"co-marketing": {"Shopify", "Stripe", "Zendesk", "Slack", "Atlassian"},
		"data":         {"Snowflake", "Databricks", "Tableau", "Segment", "Amplitude"},
	}

	partnerList, ok := partners[partnerType]
	if !ok {
		return nil, fmt.Errorf("unknown partnership type: %s", partnerType)
	}

	partner := partnerList[rand.Intn(len(partnerList))]

	// Calculate costs and benefits
	var cost, mrrBoost int64
	var churnReduction float64
	var duration int

	switch partnerType {
	case "distribution":
		cost = 50000 + rand.Int63n(100000)                             // $50-150k
		mrrBoost = int64(float64(fs.MRR) * (0.1 + rand.Float64()*0.2)) // 10-30% MRR boost
		churnReduction = 0.01 + rand.Float64()*0.02                    // 1-3% churn reduction
		duration = 12 + rand.Intn(12)                                  // 12-24 months
	case "technology":
		cost = 30000 + rand.Int63n(70000)                                // $30-100k
		mrrBoost = int64(float64(fs.MRR) * (0.05 + rand.Float64()*0.15)) // 5-20% MRR boost
		churnReduction = 0.02 + rand.Float64()*0.03                      // 2-5% churn reduction
		duration = 12 + rand.Intn(24)                                    // 12-36 months
	case "co-marketing":
		cost = 25000 + rand.Int63n(50000)                                // $25-75k
		mrrBoost = int64(float64(fs.MRR) * (0.15 + rand.Float64()*0.25)) // 15-40% MRR boost
		churnReduction = 0.005 + rand.Float64()*0.015                    // 0.5-2% churn reduction
		duration = 6 + rand.Intn(12)                                     // 6-18 months
	case "data":
		cost = 40000 + rand.Int63n(60000)                                // $40-100k
		mrrBoost = int64(float64(fs.MRR) * (0.08 + rand.Float64()*0.12)) // 8-20% MRR boost
		churnReduction = 0.01 + rand.Float64()*0.02                      // 1-3% churn reduction
		duration = 12 + rand.Intn(24)                                    // 12-36 months
	}

	if cost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash for partnership (need $%s)", formatCurrency(cost))
	}

	partnership := Partnership{
		Partner:        partner,
		Type:           partnerType,
		MonthStarted:   fs.Turn,
		Duration:       duration,
		Cost:           cost,
		MRRBoost:       mrrBoost,
		ChurnReduction: churnReduction,
		Status:         "active",
	}

	fs.Cash -= cost
	fs.Partnerships = append(fs.Partnerships, partnership)

	return &partnership, nil
}

// UpdatePartnerships processes active partnerships
func (fs *FounderState) UpdatePartnerships() []string {
	var messages []string

	for i := range fs.Partnerships {
		p := &fs.Partnerships[i]
		if p.Status != "active" {
			continue
		}

		// Check if partnership has expired
		monthsActive := fs.Turn - p.MonthStarted
		if monthsActive >= p.Duration {
			p.Status = "expired"
			messages = append(messages, fmt.Sprintf("⏰ Partnership with %s has expired", p.Partner))

			// Remove benefits
			fs.MRR -= p.MRRBoost
			fs.CustomerChurnRate += p.ChurnReduction
			continue
		}

		// Apply ongoing benefits (already included in calculations)
	}

	return messages
}

// ============================================================================
// AFFILIATE PROGRAM
// ============================================================================

// LaunchAffiliateProgram starts an affiliate program
func (fs *FounderState) LaunchAffiliateProgram(commission float64) error {
	if fs.AffiliateProgram != nil {
		return fmt.Errorf("affiliate program already running")
	}

	setupCost := int64(20000 + rand.Int63n(30000)) // $20-50k setup
	if setupCost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(setupCost))
	}

	fs.Cash -= setupCost

	fs.AffiliateProgram = &AffiliateProgram{
		LaunchedMonth:      fs.Turn,
		Commission:         commission / 100,  // Convert % to decimal
		Affiliates:         5 + rand.Intn(10), // Start with 5-15 affiliates
		SetupCost:          setupCost,
		MonthlyPlatformFee: 5000 + rand.Int63n(5000), // $5-10k/month
		MonthlyRevenue:     0,
		CustomersAcquired:  0,
	}

	return nil
}

// UpdateAffiliateProgram processes monthly affiliate sales
func (fs *FounderState) UpdateAffiliateProgram() []string {
	var messages []string

	if fs.AffiliateProgram == nil {
		return messages
	}

	prog := fs.AffiliateProgram

	// Pay platform fee
	fs.Cash -= prog.MonthlyPlatformFee

	// Calculate affiliate sales (each affiliate brings 0-2 customers/month)
	newCustomers := 0
	for i := 0; i < prog.Affiliates; i++ {
		if rand.Float64() < 0.3 { // 30% chance per affiliate
			newCustomers += 1 + rand.Intn(2)
		}
	}

	if newCustomers > 0 {
		newMRR := int64(newCustomers) * fs.AvgDealSize
		commissionPaid := int64(float64(newMRR) * prog.Commission)

		fs.Customers += newCustomers
		fs.MRR += newMRR
		fs.Cash -= commissionPaid

		prog.CustomersAcquired += newCustomers
		prog.MonthlyRevenue += newMRR

		messages = append(messages, fmt.Sprintf("🤝 Affiliates brought %d customers ($%s MRR, $%s commission)",
			newCustomers, formatCurrency(newMRR), formatCurrency(commissionPaid)))

		// Affiliates grow over time if successful
		if rand.Float64() < 0.2 {
			prog.Affiliates += 1 + rand.Intn(3)
		}
	}

	return messages
}

// ============================================================================
// COMPETITORS
// ============================================================================

// SpawnCompetitor randomly creates a new competitor
func (fs *FounderState) SpawnCompetitor() *Competitor {
	// 3% chance per month after month 12
	if fs.Turn < 12 || rand.Float64() > 0.03 {
		return nil
	}

	names := []string{
		"TechFlow", "DataCore", "CloudSync", "SmartSuite", "VentureStack",
		"NexusApp", "PulseAI", "ZenithCo", "ApexSoft", "CoreLogic",
		"FusionTech", "QuantumLeap", "SynergyLabs", "InnovateX", "SparkCode",
	}

	threatLevels := []string{"low", "medium", "high"}
	threat := threatLevels[rand.Intn(len(threatLevels))]

	var marketShare float64
	switch threat {
	case "low":
		marketShare = 0.01 + rand.Float64()*0.04 // 1-5%
	case "medium":
		marketShare = 0.05 + rand.Float64()*0.10 // 5-15%
	case "high":
		marketShare = 0.10 + rand.Float64()*0.15 // 10-25%
	}

	comp := Competitor{
		Name:          names[rand.Intn(len(names))],
		Threat:        threat,
		MarketShare:   marketShare,
		Strategy:      "ignore", // Default strategy
		MonthAppeared: fs.Turn,
		Active:        true,
	}

	fs.Competitors = append(fs.Competitors, comp)
	return &comp
}

// HandleCompetitor allows player to respond to a competitor
func (fs *FounderState) HandleCompetitor(compIndex int, strategy string) (string, error) {
	if compIndex < 0 || compIndex >= len(fs.Competitors) {
		return "", fmt.Errorf("invalid competitor index")
	}

	comp := &fs.Competitors[compIndex]
	if !comp.Active {
		return "", fmt.Errorf("competitor is no longer active")
	}

	switch strategy {
	case "ignore":
		comp.Strategy = "ignore"
		return fmt.Sprintf("Ignoring %s. They may gain market share.", comp.Name), nil

	case "compete":
		cost := int64(50000 + rand.Int63n(100000)) // $50-150k
		if cost > fs.Cash {
			return "", fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}

		fs.Cash -= cost
		comp.Strategy = "compete"

		// Reduce their threat
		if comp.Threat == "high" {
			comp.Threat = "medium"
		} else if comp.Threat == "medium" {
			comp.Threat = "low"
		}
		comp.MarketShare *= 0.7 // Reduce their market share

		return fmt.Sprintf("Competing aggressively! Cost: $%s. Reduced %s threat to %s",
			formatCurrency(cost), comp.Name, comp.Threat), nil

	case "partner":
		cost := int64(100000 + rand.Int63n(150000)) // $100-250k
		if cost > fs.Cash {
			return "", fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}

		fs.Cash -= cost
		comp.Strategy = "partner"
		comp.Active = false

		// Merge customer bases
		newCustomers := int(float64(fs.Customers) * comp.MarketShare)
		fs.Customers += newCustomers
		newMRR := int64(newCustomers) * fs.AvgDealSize
		fs.MRR += newMRR

		return fmt.Sprintf("Partnered with %s! Cost: $%s. Gained %d customers (+$%s MRR)",
			comp.Name, formatCurrency(cost), newCustomers, formatCurrency(newMRR)), nil

	default:
		return "", fmt.Errorf("unknown strategy: %s", strategy)
	}
}

// UpdateCompetitors processes competitor actions
func (fs *FounderState) UpdateCompetitors() []string {
	var messages []string

	for i := range fs.Competitors {
		comp := &fs.Competitors[i]
		if !comp.Active {
			continue
		}

		// Ignored competitors grow stronger
		if comp.Strategy == "ignore" {
			comp.MarketShare *= 1.05 // Grow 5% per month

			// May increase threat level
			if comp.MarketShare > 0.20 && comp.Threat != "high" {
				comp.Threat = "high"
				messages = append(messages, fmt.Sprintf("⚠️  %s is now a HIGH threat!", comp.Name))
			} else if comp.MarketShare > 0.10 && comp.Threat == "low" {
				comp.Threat = "medium"
				messages = append(messages, fmt.Sprintf("⚠️  %s is now a MEDIUM threat", comp.Name))
			}
		}

		// Competing with them slows their growth
		if comp.Strategy == "compete" {
			comp.MarketShare *= 0.95 // Shrink 5% per month
			if comp.MarketShare < 0.02 {
				comp.Active = false
				messages = append(messages, fmt.Sprintf("✅ %s has exited the market!", comp.Name))
			}
		}
	}

	return messages
}

// ============================================================================
// GLOBAL MARKETS
// ============================================================================

// ExpandToMarket launches in a new geographic region
func (fs *FounderState) ExpandToMarket(region string) (*Market, error) {
	// Check if already in this market
	for _, m := range fs.GlobalMarkets {
		if m.Region == region {
			return nil, fmt.Errorf("already operating in %s", region)
		}
	}

	// Define market parameters
	var setupCost, monthlyCost int64
	var marketSize int
	var competition string

	switch region {
	case "Europe":
		setupCost = 200000
		monthlyCost = 30000
		marketSize = 50000
		competition = "high"
	case "Asia":
		setupCost = 250000
		monthlyCost = 40000
		marketSize = 100000
		competition = "very_high"
	case "LATAM":
		setupCost = 150000
		monthlyCost = 20000
		marketSize = 30000
		competition = "medium"
	case "Middle East":
		setupCost = 180000
		monthlyCost = 25000
		marketSize = 20000
		competition = "low"
	case "Africa":
		setupCost = 120000
		monthlyCost = 15000
		marketSize = 25000
		competition = "low"
	case "Australia":
		setupCost = 100000
		monthlyCost = 18000
		marketSize = 15000
		competition = "medium"
	default:
		return nil, fmt.Errorf("unknown market: %s", region)
	}

	if setupCost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(setupCost))
	}

	fs.Cash -= setupCost

	// Calculate initial customers based on setup cost / effective local CAC
	fs.UpdateCAC()
	localCAC := fs.CustomerAcquisitionCost

	// Adjust CAC for local competition
	switch competition {
	case "very_high":
		localCAC = int64(float64(localCAC) * 1.8)
	case "high":
		localCAC = int64(float64(localCAC) * 1.5)
	case "medium":
		localCAC = int64(float64(localCAC) * 1.2)
	case "low":
		localCAC = int64(float64(localCAC) * 0.8)
	}

	// Further adjust for product maturity (immature product = harder to sell internationally)
	if fs.ProductMaturity < 0.7 {
		localCAC = int64(float64(localCAC) / fs.ProductMaturity)
	}

	initialCustomers := int(setupCost / localCAC)
	if initialCustomers < 1 {
		initialCustomers = 1 // At least 1 customer
	}

	initialMRR := int64(initialCustomers) * fs.AvgDealSize

	market := Market{
		Region:           region,
		LaunchMonth:      fs.Turn,
		SetupCost:        setupCost,
		MonthlyCost:      monthlyCost,
		CustomerCount:    initialCustomers,
		MRR:              initialMRR,
		MarketSize:       marketSize,
		Penetration:      float64(initialCustomers) / float64(marketSize),
		LocalCompetition: competition,
	}

	fs.GlobalMarkets = append(fs.GlobalMarkets, market)
	fs.Customers += initialCustomers
	fs.MRR += initialMRR

	// Increase global churn rate due to operational complexity

	fs.CustomerChurnRate += 0.01 + (rand.Float64() * 0.01)

	return &market, nil
}

// UpdateGlobalMarkets processes international operations
func (fs *FounderState) UpdateGlobalMarkets() []string {
	var messages []string

	for i := range fs.GlobalMarkets {
		m := &fs.GlobalMarkets[i]

		// Pay monthly costs
		fs.Cash -= m.MonthlyCost

		// Calculate market-specific churn rate
		marketChurn := fs.CustomerChurnRate

		// No CS team? Much higher churn (up to 50% in new markets)
		if len(fs.Team.CustomerSuccess) == 0 {
			marketChurn += 0.30 // +30% base churn without CS
		}

		// Product not mature? Higher churn
		if fs.ProductMaturity < 1.0 {
			marketChurn += (1.0 - fs.ProductMaturity) * 0.20 // Up to +20% if product immature
		}

		// No engineers? Product degrades, churn increases
		if len(fs.Team.Engineers) == 0 {
			marketChurn += 0.15 // +15% without engineers
		}

		// Local competition increases churn
		switch m.LocalCompetition {
		case "very_high":
			marketChurn += 0.10
		case "high":
			marketChurn += 0.07
		case "medium":
			marketChurn += 0.04
		}

		// Process churn first
		customersLost := int(float64(m.CustomerCount) * marketChurn)
		if customersLost > 0 {
			m.CustomerCount -= customersLost
			mrrLost := int64(customersLost) * fs.AvgDealSize
			m.MRR -= mrrLost
			fs.MRR -= mrrLost
			fs.Customers -= customersLost

			messages = append(messages, fmt.Sprintf("📉 %s: Lost %d customers (%.1f%% churn)",
				m.Region, customersLost, marketChurn*100))
		}

		// Competitors actively take customers
		for _, comp := range fs.Competitors {
			if !comp.Active {
				continue
			}
			// Competitors are more effective in markets where you're weak
			if comp.Strategy == "ignore" && len(fs.Team.Sales) < 3 {
				competitorSteal := int(float64(m.CustomerCount) * comp.MarketShare * 0.15) // 15% of their market share
				if competitorSteal > 0 {
					m.CustomerCount -= competitorSteal
					stolenMRR := int64(competitorSteal) * fs.AvgDealSize
					m.MRR -= stolenMRR
					fs.MRR -= stolenMRR
					fs.Customers -= competitorSteal

					messages = append(messages, fmt.Sprintf("⚠️  %s took %d customers in %s",
						comp.Name, competitorSteal, m.Region))
				}
			}
		}

		// Now attempt growth
		// Base growth rate much lower - you need to work for it
		baseGrowth := 0.02 + (rand.Float64() * 0.03) // 2-5% base monthly growth

		// Sales team impact (need sales to grow in new markets)
		salesImpact := float64(len(fs.Team.Sales)) * 0.02 // Each sales rep adds 2%

		// Marketing spend helps (if they spent on marketing this turn, residual effect)
		marketingImpact := 0.01 * float64(len(fs.Team.Marketing)) // Each marketer adds 1%

		// Adjust for competition
		competitionMultiplier := 1.0
		switch m.LocalCompetition {
		case "very_high":
			competitionMultiplier = 0.5 // Half growth in very competitive markets
		case "high":
			competitionMultiplier = 0.7
		case "medium":
			competitionMultiplier = 0.85
		case "low":
			competitionMultiplier = 1.1 // Easier growth in low competition
		}

		// Product maturity affects conversion
		productMultiplier := fs.ProductMaturity
		if productMultiplier < 0.5 {
			productMultiplier = 0.5 // Can't grow much with immature product
		}

		totalGrowth := (baseGrowth + salesImpact + marketingImpact) * competitionMultiplier * productMultiplier

		// Calculate new customers (as % of CURRENT customer base in this market)
		// This is much more realistic - you grow based on what you already have
		newCustomers := int(float64(m.CustomerCount) * totalGrowth)

		// But cap at remaining market opportunity
		remainingMarket := m.MarketSize - m.CustomerCount
		if newCustomers > remainingMarket {
			newCustomers = remainingMarket
		}

		if newCustomers > 0 {
			m.CustomerCount += newCustomers
			newMRR := int64(newCustomers) * fs.AvgDealSize
			m.MRR += newMRR
			fs.MRR += newMRR
			fs.Customers += newCustomers

			m.Penetration = float64(m.CustomerCount) / float64(m.MarketSize)

			messages = append(messages, fmt.Sprintf("🌍 %s: +%d customers, $%s MRR (%.1f%% penetration)",
				m.Region, newCustomers, formatCurrency(m.MRR), m.Penetration*100))
		}

		// Market can shrink if churn > growth
		if m.CustomerCount < 0 {
			m.CustomerCount = 0
			m.MRR = 0
		}
	}

	return messages
}

// ============================================================================
// PIVOTS
// ============================================================================

// ExecutePivot changes company strategy/market
func (fs *FounderState) ExecutePivot(toStrategy string, reason string) (*Pivot, error) {
	cost := int64(100000 + rand.Int63n(200000)) // $100-300k
	if cost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	fs.Cash -= cost

	// Lose customers during pivot (20-50%)
	churnRate := 0.20 + rand.Float64()*0.30
	customersLost := int(float64(fs.Customers) * churnRate)
	fs.Customers -= customersLost
	mrrLost := int64(customersLost) * fs.AvgDealSize
	fs.MRR -= mrrLost

	// Success rate depends on product maturity and timing
	successChance := 0.30 // Base 30%
	if fs.ProductMaturity > 0.7 {
		successChance += 0.20 // +20% if product mature
	}
	if fs.Turn < 24 {
		successChance += 0.20 // +20% if early (under 2 years)
	}

	success := rand.Float64() < successChance

	pivot := Pivot{
		Month:         fs.Turn,
		FromStrategy:  fs.StartupType,
		ToStrategy:    toStrategy,
		Reason:        reason,
		Cost:          cost,
		CustomersLost: customersLost,
		Success:       success,
	}

	if success {
		fs.StartupType = toStrategy
		// Expand market size on successful pivot
		fs.TargetMarketSize = int(float64(fs.TargetMarketSize) * (1.5 + rand.Float64()))
	}

	fs.PivotHistory = append(fs.PivotHistory, pivot)
	fs.CalculateRunway()

	return &pivot, nil
}

// ============================================================================
// EQUITY BUYBACKS
// ============================================================================

// BuybackEquity buys back equity from investors
func (fs *FounderState) BuybackEquity(roundName string, equityPercent float64) (*Buyback, error) {
	// Must be profitable
	if fs.MRR <= fs.MonthlyTeamCost {
		return nil, fmt.Errorf("must be profitable to buy back equity")
	}

	// Find the round
	var round *FundingRound
	for i := range fs.FundingRounds {
		if fs.FundingRounds[i].RoundName == roundName {
			round = &fs.FundingRounds[i]
			break
		}
	}

	if round == nil {
		return nil, fmt.Errorf("round not found: %s", roundName)
	}

	// Calculate current valuation
	currentValuation := int64(float64(fs.MRR) * 12 * 10) // 10x ARR

	// Price to buy back equity (at current valuation)
	price := int64(float64(currentValuation) * equityPercent / 100)

	if price > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(price))
	}

	fs.Cash -= price
	fs.EquityGivenAway -= equityPercent

	buyback := Buyback{
		Month:        fs.Turn,
		Investor:     roundName,
		EquityBought: equityPercent,
		PricePaid:    price,
		Valuation:    currentValuation,
	}

	fs.InvestorBuybacks = append(fs.InvestorBuybacks, buyback)
	fs.CalculateRunway()

	return &buyback, nil
}

// AddBoardSeat adds a board seat (usually with funding)
func (fs *FounderState) AddBoardSeat(reason string) {
	fs.BoardSeats++
	fs.EquityPool -= 2.0 // Each board seat costs 2% from equity pool
	if fs.EquityPool < 0 {
		fs.EquityPool = 0
	}
}

// ExpandEquityPool increases employee equity pool (dilutes founder)
func (fs *FounderState) ExpandEquityPool(percentToAdd float64) {
	fs.EquityPool += percentToAdd
	// This dilutes everyone proportionally
	fs.EquityGivenAway = fs.EquityGivenAway * (100.0 / (100.0 + percentToAdd))
}

// ============================================================================
// RANDOM EVENTS
// ============================================================================

// SpawnRandomEvent creates a new random event
func (fs *FounderState) SpawnRandomEvent() *RandomEvent {
	// Don't spawn events in first 6 months
	if fs.Turn < 6 {
		return nil
	}

	// Choose event type
	eventTypes := []string{"economy", "regulation", "competition", "talent", "customer", "product", "legal", "press"}
	eventType := eventTypes[rand.Intn(len(eventTypes))]

	// 40% chance of positive event
	isPositive := rand.Float64() < 0.4

	var event *RandomEvent

	switch eventType {
	case "economy":
		event = fs.generateEconomyEvent(isPositive)
	case "regulation":
		event = fs.generateRegulationEvent(isPositive)
	case "competition":
		event = fs.generateCompetitionEvent(isPositive)
	case "talent":
		event = fs.generateTalentEvent(isPositive)
	case "customer":
		event = fs.generateCustomerEvent(isPositive)
	case "product":
		event = fs.generateProductEvent(isPositive)
	case "legal":
		event = fs.generateLegalEvent(isPositive)
	case "press":
		event = fs.generatePressEvent(isPositive)
	}

	if event != nil {
		event.Month = fs.Turn
		event.Type = eventType
		fs.RandomEvents = append(fs.RandomEvents, *event)

		// Apply the effect
		if event.Impact.DurationMonths > 0 {
			fs.ActiveEventEffects[event.Title] = event.Impact
		}

		// Apply immediate cash cost
		if event.Impact.CashCost > 0 {
			fs.Cash -= event.Impact.CashCost
		}

		// Handle employee losses
		if event.Impact.EmployeesLost > 0 {
			fs.handleEmployeeLoss(event.Impact.EmployeesLost)
		}
	}

	return event
}

func (fs *FounderState) generateEconomyEvent(isPositive bool) *RandomEvent {
	if isPositive {
		// Economic boom
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Economic Boom",
			Description: "Strong economy increases customer budgets and spending",
			Impact: EventImpact{
				GrowthChange:   1.15, // +15% growth
				DurationMonths: 6,
			},
		}
	} else {
		// Economic downturn
		severity := "moderate"
		impact := EventImpact{
			GrowthChange:   0.85, // -15% growth
			ChurnChange:    0.02, // +2% churn
			DurationMonths: 6,
		}

		// Tariffs specifically affect hardware companies
		if fs.Category == "Hardware" || fs.Category == "Deep Tech" {
			severity = "major"
			impact.CACChange = 1.3 // +30% CAC due to tariffs
			impact.CashCost = 20000 + rand.Int63n(30000)
			return &RandomEvent{
				Severity:    severity,
				IsPositive:  false,
				Title:       "Trade Tariffs Imposed",
				Description: "New tariffs on hardware components increase costs significantly",
				Impact:      impact,
			}
		}

		return &RandomEvent{
			Severity:    severity,
			IsPositive:  false,
			Title:       "Economic Downturn",
			Description: "Recession fears cause customers to cut budgets",
			Impact:      impact,
		}
	}
}

func (fs *FounderState) generateRegulationEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "minor",
			IsPositive:  true,
			Title:       "Favorable Regulations Passed",
			Description: "New laws favor your business model and reduce compliance costs",
			Impact: EventImpact{
				GrowthChange:   1.10,
				CashCost:       -10000, // Negative cost = gain
				DurationMonths: 12,
			},
		}
	} else {
		impact := EventImpact{
			CashCost:       50000 + rand.Int63n(100000), // $50-150k compliance
			GrowthChange:   0.90,                        // -10% growth
			DurationMonths: 12,
		}

		// Data privacy regulations
		if fs.Category == "SaaS" {
			impact.CashCost += 50000
			return &RandomEvent{
				Severity:    "major",
				IsPositive:  false,
				Title:       "New Data Privacy Regulations",
				Description: "GDPR/CCPA-style laws require expensive compliance changes",
				Impact:      impact,
			}
		}

		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "New Industry Regulations",
			Description: "New compliance requirements slow growth and increase costs",
			Impact:      impact,
		}
	}
}

func (fs *FounderState) generateCompetitionEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Major Competitor Exits Market",
			Description: "A key competitor shuts down, opening up opportunities",
			Impact: EventImpact{
				GrowthChange:   1.25, // +25% growth
				CACChange:      0.80, // -20% CAC
				DurationMonths: 6,
			},
		}
	} else {
		// Check for open source threat
		if rand.Float64() < 0.3 {
			return &RandomEvent{
				Severity:    "major",
				IsPositive:  false,
				Title:       "Open Source Alternative Launched",
				Description: "Free alternative threatens paid product - must differentiate!",
				Impact: EventImpact{
					ChurnChange:    0.05, // +5% churn
					GrowthChange:   0.70, // -30% growth
					CACChange:      1.40, // +40% CAC
					DurationMonths: 12,
				},
			}
		}

		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "Well-Funded Competitor Emerges",
			Description: "New competitor raises $50M Series B, plans aggressive expansion",
			Impact: EventImpact{
				CACChange:      1.25, // +25% CAC
				ChurnChange:    0.03, // +3% churn
				DurationMonths: 9,
			},
		}
	}
}

func (fs *FounderState) generateTalentEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "minor",
			IsPositive:  true,
			Title:       "Rock Star Candidate Interested",
			Description: "Top talent from FAANG wants to join - hire them!",
			Impact: EventImpact{
				ProductivityChange: 1.15, // +15% team productivity
				DurationMonths:     12,
			},
		}
	} else {
		// Key employees quit
		numQuitting := 1
		if fs.Team.TotalEmployees > 10 {
			numQuitting = 1 + rand.Intn(2)
		}

		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "Key Employees Resign",
			Description: fmt.Sprintf("%d employee(s) leave for higher-paying competitors", numQuitting),
			Impact: EventImpact{
				EmployeesLost:      numQuitting,
				ProductivityChange: 0.90, // -10% productivity
				DurationMonths:     3,    // Takes 3 months to recover
			},
		}
	}
}

func (fs *FounderState) generateCustomerEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  true,
			Title:       "Enterprise Customer Win!",
			Description: "Fortune 500 company signs major contract",
			Impact: EventImpact{
				CashCost:       -(50000 + rand.Int63n(150000)), // $50-200k deal
				DurationMonths: 0,                              // One-time
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  false,
			Title:       "Major Customer Churns",
			Description: "Top 3 customer cancels unexpectedly",
			Impact: EventImpact{
				MRRChange:      0.90, // -10% MRR
				DurationMonths: 0,
			},
		}
	}
}

func (fs *FounderState) generateProductEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Product Breakthrough!",
			Description: "Engineering team ships game-changing feature",
			Impact: EventImpact{
				ChurnChange:    -0.02, // -2% churn
				GrowthChange:   1.20,  // +20% growth
				DurationMonths: 6,
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  false,
			Title:       "Critical Bug in Production",
			Description: "Major outage affects all customers for 48 hours",
			Impact: EventImpact{
				ChurnChange:    0.05,                       // +5% churn
				CashCost:       10000 + rand.Int63n(20000), // Emergency fixes
				DurationMonths: 1,
			},
		}
	}
}

func (fs *FounderState) generateLegalEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "minor",
			IsPositive:  true,
			Title:       "Patent Approved",
			Description: "Key technology patent granted, provides competitive moat",
			Impact: EventImpact{
				GrowthChange:   1.10,
				DurationMonths: 24,
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "major",
			IsPositive:  false,
			Title:       "Patent Infringement Claim",
			Description: "Large company alleges patent violation, requires legal defense",
			Impact: EventImpact{
				CashCost:       100000 + rand.Int63n(200000), // $100-300k legal fees
				DurationMonths: 0,
			},
		}
	}
}

func (fs *FounderState) generatePressEvent(isPositive bool) *RandomEvent {
	if isPositive {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  true,
			Title:       "Major Press Coverage",
			Description: "Featured in TechCrunch/WSJ - huge brand boost!",
			Impact: EventImpact{
				CACChange:      0.75, // -25% CAC
				GrowthChange:   1.30, // +30% growth
				DurationMonths: 3,
			},
		}
	} else {
		return &RandomEvent{
			Severity:    "moderate",
			IsPositive:  false,
			Title:       "PR Crisis",
			Description: "Negative press about company culture or product issues",
			Impact: EventImpact{
				ChurnChange:    0.04, // +4% churn
				CACChange:      1.35, // +35% CAC
				GrowthChange:   0.80, // -20% growth
				DurationMonths: 6,
			},
		}
	}
}

// ProcessRandomEvents applies effects of active events
func (fs *FounderState) ProcessRandomEvents() []string {
	var messages []string

	// Expire old events
	for key, impact := range fs.ActiveEventEffects {
		// Find the original event
		var originalEvent *RandomEvent
		for i := range fs.RandomEvents {
			if fs.RandomEvents[i].Title == key {
				originalEvent = &fs.RandomEvents[i]
				break
			}
		}

		if originalEvent != nil {
			monthsActive := fs.Turn - originalEvent.Month
			if monthsActive >= impact.DurationMonths {
				delete(fs.ActiveEventEffects, key)
				messages = append(messages, fmt.Sprintf("⏰ Event effect expired: %s", key))
			}
		}
	}

	// Apply active event effects (these are already factored into calculations)
	// The effects modify CAC, churn, growth rates which are used in ProcessMonth

	return messages
}

// handleEmployeeLoss fires random employees
func (fs *FounderState) handleEmployeeLoss(count int) {
	for i := 0; i < count; i++ {
		// Randomly pick a team to lose from
		roll := rand.Intn(4)
		switch roll {
		case 0:
			if len(fs.Team.Engineers) > 0 {
				fs.Team.Engineers = fs.Team.Engineers[:len(fs.Team.Engineers)-1]
			}
		case 1:
			if len(fs.Team.Sales) > 0 {
				fs.Team.Sales = fs.Team.Sales[:len(fs.Team.Sales)-1]
			}
		case 2:
			if len(fs.Team.CustomerSuccess) > 0 {
				fs.Team.CustomerSuccess = fs.Team.CustomerSuccess[:len(fs.Team.CustomerSuccess)-1]
			}
		case 3:
			if len(fs.Team.Marketing) > 0 {
				fs.Team.Marketing = fs.Team.Marketing[:len(fs.Team.Marketing)-1]
			}
		}
	}
	fs.CalculateTeamCost()
	fs.CalculateRunway()
}
