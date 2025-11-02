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
	Name            string
	Role            EmployeeRole
	MonthlyCost     int64
	Impact          float64 // Productivity/effectiveness multiplier
	IsExecutive     bool    // C-level executives have 3x impact, $300k/year salary
	Equity          float64 // Equity percentage owned by this employee
	VestingMonths   int     // Total vesting period (typically 48 months)
	CliffMonths     int     // Cliff period (typically 12 months)
	VestedMonths    int     // Months vested so far
	HasCliff        bool    // Has cliff been reached
	MonthHired      int     // Month when hired
	AssignedMarket  string  // Market assignment: "USA", "Europe", "Asia", "All", etc.
}

// CapTableEntry tracks individual equity ownership
type CapTableEntry struct {
	Name         string  // Employee name or investor round name
	Type         string  // "employee", "executive", "investor", "advisor"
	Equity       float64 // Equity percentage
	MonthGranted int     // Month when equity was granted
}

// BoardMember represents an advisor or board member
type BoardMember struct {
	Name              string
	Type              string // "advisor", "investor", "independent"
	Expertise         string // "sales", "product", "fundraising", "operations", "strategy"
	MonthAdded        int
	EquityCost        float64 // Equity given for this seat
	IsActive          bool
	ContributionScore float64 // 0-1, how valuable their advice has been
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
	FounderName        string
	CompanyName        string
	Category           string
	StartupType        string
	Description        string
	Cash               int64
	MRR                int64 // Monthly Recurring Revenue
	DirectMRR          int64 // MRR from direct customers (excludes affiliate)
	AffiliateMRR       int64 // MRR from affiliate customers
	Customers          int
	DirectCustomers    int // Customers acquired directly (excludes affiliate)
	AffiliateCustomers int // Customers acquired via affiliate program
	AvgDealSize        int64
	MinDealSize        int64 // Minimum deal size (for display)
	MaxDealSize        int64 // Maximum deal size (for display)
	ChurnRate          float64
	CustomerChurnRate  float64 // Alias for ChurnRate
	BaseCAC            int64   // Base customer acquisition cost for this business
	Team               Team
	Turn               int
	MaxTurns           int
	ProductMaturity    float64 // 0-1, affects sales velocity
	MarketPenetration  float64 // 0-1, % of target market captured
	TargetMarketSize   int
	CompetitionLevel   string
	FundingRounds      []FundingRound
	EquityGivenAway    float64       // Total % equity given to investors
	BoardSeats         int           // Board seats given to investors
	BoardMembers       []BoardMember // All board members/advisors
	AcquisitionOffers  []AcquisitionOffer
	CashRunwayMonths   int
	MonthlyTeamCost    int64 // Cached monthly team cost
	FounderSalary      int64 // $150k/year = $12,500/month

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
	EquityPool         float64 // Employee equity pool % (total allocated for employees)
	EquityAllocated    float64 // % of equity pool already allocated to employees
	InvestorBuybacks   []Buyback
	RandomEvents       []RandomEvent
	ActiveEventEffects map[string]EventImpact // Events currently affecting the business
	CapTable           []CapTableEntry        // Individual equity ownership tracking

	// Infrastructure costs
	MonthlyComputeCost int64 // Cloud compute costs (scales with customers)
	MonthlyODCCost     int64 // Other Direct Costs (scales with customers)

	// Customer tracking
	CustomerList       []Customer // Individual customer records
	TotalCustomersEver int        // Total customers acquired (including churned)
	TotalChurned       int        // Total customers that have churned
	NextCustomerID     int        // Next customer ID to assign

	// Investor/Board tracking
	BoardSentiment   string // "happy", "neutral", "concerned", "angry"
	BoardPressure    int    // 0-100, pressure to perform
	LastBoardMeeting int    // Turn of last board meeting

	// Strategic opportunities
	PendingOpportunity *StrategicOpportunity // Current opportunity awaiting decision

	// Exit tracking
	HasExited     bool
	ExitType      string // "ipo", "acquisition", "secondary", "time_limit"
	ExitValuation int64
	ExitMonth     int
}

// Customer represents an individual customer deal
type Customer struct {
	ID           int     // Unique customer ID
	Source       string  // "direct", "affiliate", "partnership", "market"
	DealSize     int64   // Monthly recurring revenue for this customer
	TermMonths   int     // Contract term in months (0 = perpetual/auto-renew)
	MonthAdded   int     // Turn when customer was acquired
	MonthChurned int     // Turn when customer churned (0 if still active)
	IsActive     bool    // Whether customer is currently active
	HealthScore  float64 // 0-1, likelihood to churn (1=healthy, 0=churning soon)
}

// StrategicOpportunity represents a one-time strategic choice
type StrategicOpportunity struct {
	Type        string // "press", "enterprise_pilot", "bridge_round", "acquisition_offer", "conference"
	Title       string
	Description string
	Cost        int64
	Benefit     string
	Risk        string
	ExpiresIn   int // Months until opportunity expires
}

// FundingRound represents a completed fundraise
type FundingRound struct {
	RoundName   string
	Amount      int64
	Valuation   int64
	EquityGiven float64
	Month       int
	Terms       string   // "Founder-friendly", "Standard", "Investor-heavy"
	Investors   []string // Names of investors in this round
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

// ExitOption represents different ways to exit the company
type ExitOption struct {
	Type          string // "ipo", "acquisition", "secondary", "continue"
	Valuation     int64
	FounderPayout int64 // How much founder gets after dilution
	Description   string
	Requirements  []string
	CanExit       bool
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

// generateDealSize creates a variable deal size based on avg, with realistic distribution
// Returns a deal size between 50% and 200% of avg, weighted toward avg, with category-based caps
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

// addCustomer adds a new customer to the tracking system
func (fs *FounderState) addCustomer(dealSize int64, source string) Customer {
	// Determine if contract is perpetual (80% chance) or fixed term (20% chance)
	termMonths := 0 // Default to perpetual
	if rand.Float64() < 0.2 {
		// Fixed term contracts: 6, 12, 24, or 36 months
		terms := []int{6, 12, 24, 36}
		termMonths = terms[rand.Intn(len(terms))]
	}

	// Initial health score based on product maturity and source
	healthScore := 0.5 + (fs.ProductMaturity * 0.3) // Base: 0.5-0.8
	if source == "affiliate" {
		healthScore -= 0.1 // Affiliates slightly less sticky
	}
	healthScore = math.Min(1.0, math.Max(0.1, healthScore))

	customer := Customer{
		ID:           fs.NextCustomerID,
		Source:       source,
		DealSize:     dealSize,
		TermMonths:   termMonths,
		MonthAdded:   fs.Turn,
		MonthChurned: 0,
		IsActive:     true,
		HealthScore:  healthScore,
	}

	fs.CustomerList = append(fs.CustomerList, customer)
	fs.NextCustomerID++
	fs.TotalCustomersEver++

	return customer
}

// churnCustomer marks a customer as churned
func (fs *FounderState) churnCustomer(customerID int) {
	for i := range fs.CustomerList {
		if fs.CustomerList[i].ID == customerID && fs.CustomerList[i].IsActive {
			fs.CustomerList[i].IsActive = false
			fs.CustomerList[i].MonthChurned = fs.Turn
			fs.TotalChurned++
			break
		}
	}
}

// GetActiveCustomers returns all currently active customers
func (fs *FounderState) GetActiveCustomers() []Customer {
	var active []Customer
	for _, c := range fs.CustomerList {
		if c.IsActive {
			active = append(active, c)
		}
	}
	return active
}

// GetChurnedCustomers returns all churned customers
func (fs *FounderState) GetChurnedCustomers() []Customer {
	var churned []Customer
	for _, c := range fs.CustomerList {
		if !c.IsActive {
			churned = append(churned, c)
		}
	}
	return churned
}

// syncMRR ensures MRR is always the sum of DirectMRR + AffiliateMRR
// syncMRR recalculates MRR from CustomerList (single source of truth)
func (fs *FounderState) syncMRR() {
	var directMRR int64
	var affiliateMRR int64
	var directCount int
	var affiliateCount int

	// Recalculate from actual customer list to prevent accumulation errors
	for _, c := range fs.CustomerList {
		if !c.IsActive {
			continue
		}

		if c.Source == "affiliate" {
			affiliateMRR += c.DealSize
			affiliateCount++
		} else {
			// "direct", "partnership", "market" all count as direct
			directMRR += c.DealSize
			directCount++
		}
	}

	fs.DirectMRR = directMRR
	fs.AffiliateMRR = affiliateMRR
	fs.MRR = directMRR + affiliateMRR
	fs.DirectCustomers = directCount
	fs.AffiliateCustomers = affiliateCount
	fs.Customers = directCount + affiliateCount

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}
}

// updateDealSizeRange updates min/max deal size tracking
func (fs *FounderState) updateDealSizeRange(newDealSize int64) {
	if fs.MinDealSize == 0 || newDealSize < fs.MinDealSize {
		fs.MinDealSize = newDealSize
	}
	if newDealSize > fs.MaxDealSize {
		fs.MaxDealSize = newDealSize
	}
}

// NewFounderGame initializes a new founder mode game
func NewFounderGame(founderName string, template StartupTemplate) *FounderState {
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
	return fs.Cash <= 0 || fs.Turn > fs.MaxTurns || fs.HasExited
}

// GetAvailableExits returns possible exit options based on current state
func (fs *FounderState) GetAvailableExits() []ExitOption {
	var exits []ExitOption
	founderEquity := (100.0 - fs.EquityPool - fs.EquityGivenAway) / 100.0

	// Calculate current valuation (simplified: ARR * multiple based on growth/profitability)
	arr := fs.MRR * 12
	multiple := 10.0                 // Base multiple
	if fs.MonthlyGrowthRate > 0.05 { // >5% monthly growth
		multiple += 5.0
	}
	if fs.MRR > fs.MonthlyTeamCost { // Profitable
		multiple += 3.0
	}
	currentValuation := int64(float64(arr) * multiple)

	// 1. IPO Option
	ipoReqs := []string{}
	canIPO := true

	if arr < 20000000 { // $20M ARR minimum
		ipoReqs = append(ipoReqs, fmt.Sprintf("❌ Need $20M ARR (currently $%s)", formatCurrency(arr)))
		canIPO = false
	} else {
		ipoReqs = append(ipoReqs, "✓ $20M+ ARR")
	}

	if len(fs.FundingRounds) < 2 {
		ipoReqs = append(ipoReqs, "❌ Need at least Series A funding")
		canIPO = false
	} else {
		ipoReqs = append(ipoReqs, "✓ Multiple funding rounds")
	}

	if fs.MonthlyGrowthRate < 0.03 { // <3% monthly = <40% YoY
		ipoReqs = append(ipoReqs, "❌ Need 40%+ YoY growth")
		canIPO = false
	} else {
		ipoReqs = append(ipoReqs, "✓ Strong growth rate")
	}

	ipoValuation := int64(float64(currentValuation) * 1.3)                 // 30% IPO premium
	ipoFounderPayout := int64(float64(ipoValuation) * founderEquity * 0.2) // Can sell 20% at IPO

	exits = append(exits, ExitOption{
		Type:          "ipo",
		Valuation:     ipoValuation,
		FounderPayout: ipoFounderPayout,
		Description:   "Take company public on NASDAQ. Provides liquidity but you remain CEO.",
		Requirements:  ipoReqs,
		CanExit:       canIPO,
	})

	// 2. Strategic Acquisition
	acqReqs := []string{}
	canAcquire := true

	if arr < 5000000 { // $5M ARR minimum
		acqReqs = append(acqReqs, fmt.Sprintf("❌ Need $5M ARR (currently $%s)", formatCurrency(arr)))
		canAcquire = false
	} else {
		acqReqs = append(acqReqs, "✓ $5M+ ARR")
	}

	if fs.Customers < 50 {
		acqReqs = append(acqReqs, fmt.Sprintf("❌ Need 50+ customers (currently %d)", fs.Customers))
		canAcquire = false
	} else {
		acqReqs = append(acqReqs, "✓ Significant customer base")
	}

	acqValuation := int64(float64(currentValuation) * 1.1) // 10% acquisition premium
	acqFounderPayout := int64(float64(acqValuation) * founderEquity)

	exits = append(exits, ExitOption{
		Type:          "acquisition",
		Valuation:     acqValuation,
		FounderPayout: acqFounderPayout,
		Description:   "Sell company to strategic acquirer. Full liquidity, but you're done.",
		Requirements:  acqReqs,
		CanExit:       canAcquire,
	})

	// 3. Secondary Sale (Private Equity)
	secondaryReqs := []string{}
	canSecondary := true

	if arr < 10000000 {
		secondaryReqs = append(secondaryReqs, fmt.Sprintf("❌ Need $10M ARR (currently $%s)", formatCurrency(arr)))
		canSecondary = false
	} else {
		secondaryReqs = append(secondaryReqs, "✓ $10M+ ARR")
	}

	netIncome := fs.MRR - fs.MonthlyTeamCost - fs.MonthlyComputeCost - fs.MonthlyODCCost
	if netIncome < 0 {
		secondaryReqs = append(secondaryReqs, "❌ Must be profitable")
		canSecondary = false
	} else {
		secondaryReqs = append(secondaryReqs, "✓ Profitable")
	}

	secondaryValuation := currentValuation
	secondaryFounderPayout := int64(float64(secondaryValuation) * founderEquity * 0.5) // Sell 50% of your stake

	exits = append(exits, ExitOption{
		Type:          "secondary",
		Valuation:     secondaryValuation,
		FounderPayout: secondaryFounderPayout,
		Description:   "Sell 50% of your shares to private equity. Partial liquidity, stay as CEO.",
		Requirements:  secondaryReqs,
		CanExit:       canSecondary,
	})

	// 4. Continue building
	exits = append(exits, ExitOption{
		Type:          "continue",
		Valuation:     currentValuation,
		FounderPayout: 0,
		Description:   "Keep building. Your current company value.",
		Requirements:  []string{"No requirements - always available"},
		CanExit:       true,
	})

	return exits
}

// ExecuteExit processes an exit decision
func (fs *FounderState) ExecuteExit(exitType string) {
	fs.HasExited = true
	fs.ExitType = exitType
	fs.ExitMonth = fs.Turn

	exits := fs.GetAvailableExits()
	for _, exit := range exits {
		if exit.Type == exitType {
			fs.ExitValuation = exit.Valuation
			break
		}
	}
}

// GetFinalScore calculates the final outcome
func (fs *FounderState) GetFinalScore() (outcome string, valuation int64, founderEquity float64) {
	founderEquity = 100.0 - fs.EquityPool - fs.EquityGivenAway

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

// RecalculateChurnRate updates the displayed churn rate based on product maturity and CS team
func (fs *FounderState) RecalculateChurnRate() {
	// Calculate base churn from product maturity
	baseChurnFromMaturity := (1.0-fs.ProductMaturity)*0.65 + 0.05
	baseChurnFromMaturity = math.Max(0.05, math.Min(0.70, baseChurnFromMaturity))
	baseChurn := baseChurnFromMaturity

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

	// Calculate effective churn rate (after CS team impact)
	actualChurn := math.Max(0.01, baseChurn-csImpact)
	fs.CustomerChurnRate = actualChurn
	fs.ChurnRate = fs.CustomerChurnRate
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
		employee = Employee{
			Role:           role,
			MonthlyCost:    avgSalary / 12,
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

// HireEmployeeWithMarket adds a new team member assigned to a specific market
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
func (fs *FounderState) GenerateTermSheetOptions(roundName string) []TermSheetOption {
	// Use fixed equity percentages per round (industry standard)
	// Then calculate valuations from those percentages
	var baseEquityPercent float64
	var baseRaise int64
	switch roundName {
	case "Seed":
		baseEquityPercent = 12.0 // 12% base equity for seed
		baseRaise = 2000000      // $2M base raise
	case "Series A":
		baseEquityPercent = 22.0 // 22% base equity for Series A
		baseRaise = 10000000     // $10M base raise
	case "Series B":
		baseEquityPercent = 18.0 // 18% base equity for Series B
		baseRaise = 30000000     // $30M base raise
	default:
		return []TermSheetOption{}
	}

	options := []TermSheetOption{}

	// Option 1: Less money, founder-friendly (lower dilution)
	option1Amount := int64(float64(baseRaise) * 0.7)
	option1Equity := baseEquityPercent * 0.85 // 85% of base (lower dilution)
	option1PostVal := int64(float64(option1Amount) / (option1Equity / 100.0))
	option1PreVal := option1PostVal - option1Amount
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
	option2Equity := baseEquityPercent // Base equity percentage
	option2PostVal := int64(float64(option2Amount) / (option2Equity / 100.0))
	option2PreVal := option2PostVal - option2Amount
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
	option3Equity := baseEquityPercent * 1.15 // 15% more equity (higher dilution)
	option3PostVal := int64(float64(option3Amount) / (option3Equity / 100.0))
	option3PreVal := option3PostVal - option3Amount
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
	option4Equity := baseEquityPercent * 1.35 // 35% more equity (investor-heavy)
	option4PostVal := int64(float64(option4Amount) / (option4Equity / 100.0))
	option4PreVal := option4PostVal - option4Amount
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

// GenerateInvestorNames creates realistic investor names based on round type
func GenerateInvestorNames(roundName string, amount int64) []string {
	var investors []string

	// Angel/Pre-Seed: individual angels
	angelInvestors := []string{
		"Naval Ravikant", "Balaji Srinivasan", "Jason Calacanis", "David Sacks",
		"Elad Gil", "Lachy Groom", "Sahil Bloom", "Anthony Pompliano",
		"Cyan Banister", "Alexis Ohanian", "Arlan Hamilton", "Kevin Hale",
	}

	// Seed: Mix of angels and micro VCs
	seedFirms := []string{
		"Y Combinator", "Sequoia Scout", "First Round", "SV Angel",
		"Hustle Fund", "Khosla Ventures", "Initialized Capital", "Haystack",
		"Precursor Ventures", "Boost VC", "Felicis Ventures", "Homebrew",
	}

	// Series A: Traditional VCs
	seriesAFirms := []string{
		"Sequoia Capital", "Andreessen Horowitz", "Accel Partners", "Benchmark",
		"Greylock Partners", "Kleiner Perkins", "Lightspeed Venture Partners",
		"Index Ventures", "General Catalyst", "NEA", "GGV Capital", "Bessemer",
	}

	// Series B+: Growth firms and strategics
	growthFirms := []string{
		"Tiger Global", "Insight Partners", "Coatue Management", "DST Global",
		"SoftBank Vision Fund", "General Atlantic", "Thrive Capital", "IVP",
		"Ribbit Capital", "T. Rowe Price", "Fidelity Investments", "BlackRock",
	}

	// Family offices and strategics
	familyOffices := []string{
		"Founders Fund", "Bezos Expeditions", "Schmidt Futures", "Emerson Collective",
		"Zuckerberg Family Office", "Gates Ventures", "Cuban Companies", "Thiel Capital",
	}

	switch roundName {
	case "Angel", "Pre-Seed":
		// 2-4 angels
		count := 2 + rand.Intn(3)
		for i := 0; i < count && i < len(angelInvestors); i++ {
			investors = append(investors, angelInvestors[rand.Intn(len(angelInvestors))])
		}

	case "Seed":
		// Lead VC + 1-2 angels or micro VCs
		investors = append(investors, seedFirms[rand.Intn(len(seedFirms))])
		if rand.Float64() > 0.5 {
			investors = append(investors, angelInvestors[rand.Intn(len(angelInvestors))])
		}
		if amount > 2000000 { // Larger seed rounds have more investors
			investors = append(investors, seedFirms[rand.Intn(len(seedFirms))])
		}

	case "Series A":
		// Lead VC + co-investors
		investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
		if amount > 10000000 {
			investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
		}
		// Sometimes strategic or family office
		if rand.Float64() > 0.7 {
			investors = append(investors, familyOffices[rand.Intn(len(familyOffices))])
		}

	case "Series B", "Series C", "Series D":
		// Growth firms + existing investors
		investors = append(investors, growthFirms[rand.Intn(len(growthFirms))])
		investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
		if amount > 50000000 {
			investors = append(investors, growthFirms[rand.Intn(len(growthFirms))])
		}
		if rand.Float64() > 0.6 {
			investors = append(investors, familyOffices[rand.Intn(len(familyOffices))])
		}

	default:
		// Generic round - mix it up
		investors = append(investors, seriesAFirms[rand.Intn(len(seriesAFirms))])
	}

	return investors
}

// RaiseFunding attempts to raise a funding round with a chosen term sheet
func (fs *FounderState) RaiseFundingWithTerms(roundName string, option TermSheetOption) (success bool) {
	fs.Cash += option.Amount
	fs.EquityGivenAway += option.Equity

	// Generate investor names for this round
	investors := GenerateInvestorNames(roundName, option.Amount)

	round := FundingRound{
		RoundName:   roundName,
		Amount:      option.Amount,
		Valuation:   option.PreValuation,
		EquityGiven: option.Equity,
		Month:       fs.Turn,
		Terms:       option.Terms,
		Investors:   investors,
	}
	fs.FundingRounds = append(fs.FundingRounds, round)

	// Add investors to cap table (split equity among them)
	equityPerInvestor := option.Equity / float64(len(investors))
	for _, investor := range investors {
		fs.CapTable = append(fs.CapTable, CapTableEntry{
			Name:         investor,
			Type:         "investor",
			Equity:       equityPerInvestor,
			MonthGranted: fs.Turn,
		})
	}

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

	// Calculate MRR with variable deal sizes
	// Use category-based defaults if AvgDealSize is 0 (no customers yet)
	baseDealSize := fs.AvgDealSize
	if baseDealSize == 0 {
		switch fs.Category {
		case "SaaS":
			baseDealSize = 1000 // Default $1k/month for SaaS
		case "DeepTech":
			baseDealSize = 5000 // Default $5k/month for DeepTech
		case "GovTech":
			baseDealSize = 2000 // Default $2k/month for GovTech
		case "Hardware":
			baseDealSize = 3000 // Default $3k/month for Hardware
		default:
			baseDealSize = 1000 // Default $1k/month
		}
	}

	var totalMRR int64
	var dealSizes []int64 // Store deal sizes for customer tracking
	for i := 0; i < newCustomers; i++ {
		dealSize := generateDealSize(baseDealSize, fs.Category)
		fs.updateDealSizeRange(dealSize)
		totalMRR += dealSize
		dealSizes = append(dealSizes, dealSize)
	}

	// These are direct customers (not from affiliate program)
	fs.Customers += newCustomers
	fs.DirectCustomers += newCustomers
	fs.DirectMRR += totalMRR

	// Add customers to tracking system
	for _, dealSize := range dealSizes {
		fs.addCustomer(dealSize, "direct")
	}

	// Sync MRR from DirectMRR + AffiliateMRR
	fs.syncMRR()

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

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

// CalculateInfrastructureCosts computes cloud compute and ODC costs based on customers
func (fs *FounderState) CalculateInfrastructureCosts() {
	// Calculate compute costs per customer (variable, random, but never more than deal size)
	// For SaaS, compute costs are typically 10-30% of deal size
	// For DeepTech/Hardware, compute costs might be higher (20-40%)
	// For GovTech, compute costs might be lower (5-15%)

	// Compute cost percentages vary by category
	var computePercent float64
	switch fs.Category {
	case "SaaS":
		computePercent = 0.10 + rand.Float64()*0.20 // 10-30% of deal size
	case "DeepTech":
		computePercent = 0.20 + rand.Float64()*0.20 // 20-40% of deal size
	case "GovTech":
		computePercent = 0.05 + rand.Float64()*0.10 // 5-15% of deal size
	case "Hardware":
		computePercent = 0.15 + rand.Float64()*0.25 // 15-40% of deal size
	default:
		computePercent = 0.10 + rand.Float64()*0.20 // 10-30% default
	}

	// ODC costs are typically 5-15% of deal size (support, data transfer, etc.)
	odcPercent := 0.05 + rand.Float64()*0.10

	// Calculate costs based on each active customer's deal size
	var totalComputeCost int64
	var totalODCCost int64

	for _, c := range fs.CustomerList {
		if !c.IsActive {
			continue
		}

		dealSize := c.DealSize
		if dealSize == 0 {
			continue
		}

		// Compute cost for this customer
		customerComputeCost := int64(float64(dealSize) * computePercent)
		customerODCCost := int64(float64(dealSize) * odcPercent)

		// Ensure total infrastructure cost never exceeds 80% of deal size
		maxCost := int64(float64(dealSize) * 0.80)
		if customerComputeCost+customerODCCost > maxCost {
			// Scale down proportionally
			totalCost := float64(customerComputeCost + customerODCCost)
			scale := float64(maxCost) / totalCost
			customerComputeCost = int64(float64(customerComputeCost) * scale)
			customerODCCost = int64(float64(customerODCCost) * scale)
		}

		totalComputeCost += customerComputeCost
		totalODCCost += customerODCCost
	}

	fs.MonthlyComputeCost = totalComputeCost
	fs.MonthlyODCCost = totalODCCost
}

// SolicitCustomerFeedback gathers feedback from customers to improve product maturity
func (fs *FounderState) SolicitCustomerFeedback() error {
	if fs.Customers == 0 {
		return fmt.Errorf("no customers to solicit feedback from")
	}

	// Feedback improves product maturity by 1-5% based on customer count
	// More customers = better feedback = more improvement
	improvement := 0.01 + (float64(fs.Customers)/100.0)*0.04 // 1-5% improvement
	if improvement > 0.05 {
		improvement = 0.05 // Cap at 5%
	}

	fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+improvement)

	// Customer feedback also reduces churn by 3-10%
	churnReduction := 0.03 + rand.Float64()*0.07                               // 3-10% reduction
	fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-churnReduction) // Minimum 1% churn
	fs.ChurnRate = fs.CustomerChurnRate

	return nil
}

// ============================================================================
// KEY METRICS & ANALYTICS
// ============================================================================

// CalculateLTVToCAC calculates the lifetime value to customer acquisition cost ratio
func (fs *FounderState) CalculateLTVToCAC() float64 {
	if fs.CustomerAcquisitionCost == 0 {
		return 0
	}

	// LTV = Average Revenue per Customer / Churn Rate
	avgRevenuePerCustomer := float64(fs.AvgDealSize)
	ltv := avgRevenuePerCustomer / math.Max(0.01, fs.CustomerChurnRate)

	return ltv / float64(fs.CustomerAcquisitionCost)
}

// CalculateCACPayback calculates months to recover customer acquisition cost
func (fs *FounderState) CalculateCACPayback() float64 {
	if fs.AvgDealSize == 0 {
		return 0
	}

	// Payback period = CAC / Monthly Revenue per Customer
	return float64(fs.CustomerAcquisitionCost) / float64(fs.AvgDealSize)
}

// CalculateRuleOf40 calculates growth rate + profit margin (should be >40% for healthy SaaS)
func (fs *FounderState) CalculateRuleOf40() float64 {
	// Growth rate (as %)
	growthRate := fs.MonthlyGrowthRate * 100

	// Profit margin = (MRR - Costs) / MRR * 100
	annualizedMRR := fs.MRR * 12
	annualCosts := (fs.MonthlyTeamCost + fs.MonthlyComputeCost + fs.MonthlyODCCost) * 12

	var profitMargin float64
	if annualizedMRR > 0 {
		profitMargin = (float64(annualizedMRR-annualCosts) / float64(annualizedMRR)) * 100
	}

	return growthRate + profitMargin
}

// CalculateBurnMultiple calculates cash burned per dollar of new ARR
func (fs *FounderState) CalculateBurnMultiple() float64 {
	if fs.Turn < 2 {
		return 0 // Need at least 2 months of data
	}

	// Monthly burn (if negative cash flow)
	monthlyRevenue := fs.MRR
	monthlyCosts := fs.MonthlyTeamCost + fs.MonthlyComputeCost + fs.MonthlyODCCost
	monthlyBurn := monthlyCosts - monthlyRevenue

	if monthlyBurn <= 0 {
		return 0 // Profitable, no burn
	}

	// New ARR = Growth * 12
	newMonthlyRevenue := int64(float64(fs.MRR) * fs.MonthlyGrowthRate)
	newARR := newMonthlyRevenue * 12

	if newARR <= 0 {
		return 999 // Burning with no growth
	}

	return float64(monthlyBurn) / float64(newARR)
}

// CalculateMagicNumber calculates sales efficiency (revenue per dollar spent on sales/marketing)
func (fs *FounderState) CalculateMagicNumber() float64 {
	if fs.Turn < 2 {
		return 0
	}

	// Sales & Marketing spend = salaries for sales + marketing + any marketing campaigns
	salesMarketingCost := int64(0)
	for range fs.Team.Sales {
		salesMarketingCost += 100000 / 12 // $100k/year salary
	}
	for range fs.Team.Marketing {
		salesMarketingCost += 100000 / 12
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCGO {
			salesMarketingCost += 300000 / 12 // $300k/year
		}
	}

	if salesMarketingCost == 0 {
		return 0
	}

	// New revenue this quarter / S&M spend this quarter
	newQuarterlyRevenue := int64(float64(fs.MRR) * fs.MonthlyGrowthRate * 3) // 3 months

	return float64(newQuarterlyRevenue) / float64(salesMarketingCost*3)
}

// MonthlyHighlight represents a significant event or metric
type MonthlyHighlight struct {
	Type    string // "win" or "concern"
	Message string
	Icon    string
}

// GenerateMonthlyHighlights creates top wins and concerns for the month
func (fs *FounderState) GenerateMonthlyHighlights() []MonthlyHighlight {
	var highlights []MonthlyHighlight

	// WINS
	if fs.MRR >= 100000 && fs.Turn > 1 {
		// Check if we just crossed $100k
		prevMRR := int64(float64(fs.MRR) / (1.0 + fs.MonthlyGrowthRate))
		if prevMRR < 100000 {
			highlights = append(highlights, MonthlyHighlight{
				Type:    "win",
				Message: "Broke $100k MRR milestone! 🎉",
				Icon:    "💰",
			})
		}
	}

	if fs.CustomerChurnRate < 0.05 && fs.Customers > 5 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "Churn rate below 5% - excellent retention!",
			Icon:    "🎯",
		})
	}

	if fs.MonthlyGrowthRate > 0.20 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: fmt.Sprintf("Strong growth: %.0f%% MoM!", fs.MonthlyGrowthRate*100),
			Icon:    "📈",
		})
	}

	ltvCac := fs.CalculateLTVToCAC()
	if ltvCac >= 3.0 && fs.Customers > 10 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "LTV:CAC ratio is healthy (>3:1)",
			Icon:    "✨",
		})
	}

	if fs.ProductMaturity >= 0.8 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "Product is highly mature - low churn expected",
			Icon:    "🚀",
		})
	}

	ruleOf40 := fs.CalculateRuleOf40()
	if ruleOf40 >= 40 && fs.MRR > 50000 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "win",
			Message: "Rule of 40 achieved - excellent balance!",
			Icon:    "💎",
		})
	}

	// CONCERNS
	if fs.CashRunwayMonths <= 3 && fs.CashRunwayMonths > 0 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: fmt.Sprintf("Low runway: %d months left!", fs.CashRunwayMonths),
			Icon:    "⚠️",
		})
	}

	if fs.CustomerChurnRate > 0.20 && fs.Customers > 5 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: fmt.Sprintf("High churn rate: %.0f%%/month", fs.CustomerChurnRate*100),
			Icon:    "📉",
		})
	}

	if fs.MonthlyGrowthRate < 0 && fs.Turn > 3 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "Revenue is declining - need to fix growth!",
			Icon:    "🔴",
		})
	}

	if fs.Customers == 0 && fs.Turn > 3 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "No customers yet - focus on acquisition!",
			Icon:    "⚡",
		})
	}

	burnMultiple := fs.CalculateBurnMultiple()
	if burnMultiple > 2.0 && burnMultiple < 999 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "High burn multiple - spending too much per $ of growth",
			Icon:    "💸",
		})
	}

	if ltvCac < 1.0 && ltvCac > 0 && fs.Customers > 10 {
		highlights = append(highlights, MonthlyHighlight{
			Type:    "concern",
			Message: "LTV:CAC ratio < 1 - losing money on customers!",
			Icon:    "⛔",
		})
	}

	// Prioritize: show max 3 of each type
	wins := []MonthlyHighlight{}
	concerns := []MonthlyHighlight{}

	for _, h := range highlights {
		if h.Type == "win" {
			wins = append(wins, h)
		} else {
			concerns = append(concerns, h)
		}
	}

	result := []MonthlyHighlight{}

	// Add top 3 wins
	for i := 0; i < len(wins) && i < 3; i++ {
		result = append(result, wins[i])
	}

	// Add top 3 concerns
	for i := 0; i < len(concerns) && i < 3; i++ {
		result = append(result, concerns[i])
	}

	return result
}

// GetCustomerHealthSegments categorizes customers by health score
func (fs *FounderState) GetCustomerHealthSegments() (healthy, atRisk, critical int, atRiskMRR, criticalMRR int64) {
	for i := range fs.CustomerList {
		c := &fs.CustomerList[i]
		if !c.IsActive {
			continue
		}

		// Update health score based on product maturity and CS team
		csImpact := float64(len(fs.Team.CustomerSuccess)) * 0.05
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCOO {
				csImpact += 0.15 // COO = 3x CS rep
			}
		}

		// Health improves with product maturity and CS team
		c.HealthScore = 0.3 + (fs.ProductMaturity * 0.5) + csImpact
		c.HealthScore = math.Min(1.0, math.Max(0.0, c.HealthScore))

		// Categorize
		if c.HealthScore >= 0.7 {
			healthy++
		} else if c.HealthScore >= 0.4 {
			atRisk++
			atRiskMRR += c.DealSize
		} else {
			critical++
			criticalMRR += c.DealSize
		}
	}

	return
}

// GetBoardGuidance generates monthly advice/opportunities from board members
func (fs *FounderState) GetBoardGuidance() []string {
	var guidance []string

	if len(fs.BoardMembers) == 0 {
		return guidance
	}

	// Board members provide value based on their expertise
	for _, member := range fs.BoardMembers {
		if !member.IsActive {
			continue
		}

		// 30% chance per month a board member provides useful guidance
		if rand.Float64() < 0.3 {
			switch member.Expertise {
			case "sales":
				// Sales expertise helps with customer acquisition
				boost := int64(float64(fs.MRR) * (0.02 + rand.Float64()*0.03)) // 2-5% boost
				if boost > 0 {
					guidance = append(guidance, fmt.Sprintf("📊 %s (Sales Advisor) introduced you to potential customers (+$%s MRR opportunity)",
						member.Name, formatCurrency(boost)))
					// Could apply boost here or make it an opportunity
				}
			case "product":
				// Product expertise improves product maturity
				if fs.ProductMaturity < 1.0 {
					improvement := 0.02 + rand.Float64()*0.03 // 2-5% improvement
					fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+improvement)
					guidance = append(guidance, fmt.Sprintf("🎯 %s (Product Advisor) helped improve product (%.0f%% maturity gained)",
						member.Name, improvement*100))
				}
			case "fundraising":
				// Fundraising expertise improves future round terms
				if len(fs.FundingRounds) < 3 {
					guidance = append(guidance, fmt.Sprintf("💰 %s (Fundraising Advisor) is warming up investors for your next round",
						member.Name))
				}
			case "operations":
				// Operations expertise reduces costs
				if fs.MonthlyTeamCost > 50000 {
					savings := int64(float64(fs.MonthlyTeamCost) * (0.01 + rand.Float64()*0.02)) // 1-3% savings
					fs.Cash += savings
					guidance = append(guidance, fmt.Sprintf("⚙️  %s (Operations Advisor) identified cost savings (+$%s this month)",
						member.Name, formatCurrency(savings)))
				}
			case "strategy":
				// Strategy expertise helps avoid bad decisions
				if fs.CustomerChurnRate > 0.15 {
					reduction := 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
					fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-reduction)
					guidance = append(guidance, fmt.Sprintf("🎓 %s (Strategy Advisor) helped reduce churn (%.0f%% improvement)",
						member.Name, reduction*100))
				}
			}
		}
	}

	return guidance
}

// UpdateBoardSentiment updates board/investor sentiment based on performance
func (fs *FounderState) UpdateBoardSentiment() {
	if len(fs.FundingRounds) == 0 {
		fs.BoardSentiment = ""
		fs.BoardPressure = 0
		return
	}

	// Calculate performance score (0-100)
	score := 50.0 // Start neutral

	// Growth is good
	if fs.MonthlyGrowthRate > 0.15 {
		score += 20
	} else if fs.MonthlyGrowthRate > 0.05 {
		score += 10
	} else if fs.MonthlyGrowthRate < 0 {
		score -= 20
	}

	// Runway matters
	if fs.CashRunwayMonths <= 3 {
		score -= 25
	} else if fs.CashRunwayMonths > 12 {
		score += 10
	}

	// Profitability is good
	netIncome := fs.MRR - fs.MonthlyTeamCost - fs.MonthlyComputeCost - fs.MonthlyODCCost
	if netIncome > 0 {
		score += 15
	}

	// Churn matters
	if fs.CustomerChurnRate < 0.05 {
		score += 10
	} else if fs.CustomerChurnRate > 0.20 {
		score -= 15
	}

	// Set sentiment
	if score >= 75 {
		fs.BoardSentiment = "happy"
		fs.BoardPressure = 10
	} else if score >= 60 {
		fs.BoardSentiment = "pleased"
		fs.BoardPressure = 25
	} else if score >= 40 {
		fs.BoardSentiment = "neutral"
		fs.BoardPressure = 50
	} else if score >= 25 {
		fs.BoardSentiment = "concerned"
		fs.BoardPressure = 75
	} else {
		fs.BoardSentiment = "angry"
		fs.BoardPressure = 95
	}
}

// GenerateStrategicOpportunity creates a random strategic choice
func (fs *FounderState) GenerateStrategicOpportunity() *StrategicOpportunity {
	// Don't generate if one already pending
	if fs.PendingOpportunity != nil {
		return nil
	}

	// 15% chance per month (after month 3)
	if fs.Turn < 3 || rand.Float64() > 0.15 {
		return nil
	}

	opportunityTypes := []string{"press", "enterprise_pilot", "conference", "talent", "competitor_distress"}

	// Add bridge round opportunity if running low on cash and have raised before
	if fs.CashRunwayMonths <= 6 && len(fs.FundingRounds) > 0 {
		opportunityTypes = append(opportunityTypes, "bridge_round")
	}

	oppType := opportunityTypes[rand.Intn(len(opportunityTypes))]

	var opp StrategicOpportunity

	switch oppType {
	case "press":
		opp = StrategicOpportunity{
			Type:        "press",
			Title:       "📰 TechCrunch Feature Opportunity",
			Description: "TechCrunch wants to write a feature story about your company. This could significantly boost brand awareness and inbound leads.",
			Cost:        10000 + rand.Int63n(15000), // $10-25k PR prep
			Benefit:     fmt.Sprintf("+%d potential customers over next 3 months, +15%% brand awareness", 5+rand.Intn(10)),
			Risk:        "Requires founder time (1 week) and PR prep costs",
			ExpiresIn:   1,
		}

	case "enterprise_pilot":
		dealSize := 50000 + rand.Int63n(150000)
		opp = StrategicOpportunity{
			Type:        "enterprise_pilot",
			Title:       "🏢 Enterprise Pilot Program",
			Description: "Fortune 500 company wants to pilot your product. High revenue potential but requires dedicated engineering resources.",
			Cost:        0, // No upfront cost, but requires team time
			Benefit:     fmt.Sprintf("$%s annual contract if successful (80%% chance), reference customer for enterprise sales", formatCurrency(dealSize)),
			Risk:        "Requires 2 engineers for 3 months, may slow product development",
			ExpiresIn:   1,
		}

	case "bridge_round":
		amount := 200000 + rand.Int63n(500000)
		equity := 3.0 + rand.Float64()*5.0
		opp = StrategicOpportunity{
			Type:        "bridge_round",
			Title:       "💰 Bridge Round Opportunity",
			Description: "Existing investor offers bridge financing at favorable terms. Quick capital to extend runway.",
			Cost:        0,
			Benefit:     fmt.Sprintf("$%s at %.1f%% equity (better terms than raising a full round)", formatCurrency(amount), equity),
			Risk:        "Additional dilution, may signal to market that you're struggling",
			ExpiresIn:   2,
		}

	case "conference":
		opp = StrategicOpportunity{
			Type:        "conference",
			Title:       "🎤 Conference Speaking Slot",
			Description: "Invited to speak at major industry conference. Great for leads and recruiting, but takes founder time.",
			Cost:        5000 + rand.Int63n(10000), // Travel + booth costs
			Benefit:     fmt.Sprintf("+%d qualified leads, improved recruiting pipeline, industry credibility", 10+rand.Intn(20)),
			Risk:        "Founder unavailable for 1 week, may not convert leads immediately",
			ExpiresIn:   2,
		}

	case "talent":
		opp = StrategicOpportunity{
			Type:        "talent",
			Title:       "⭐ Star Engineer Available",
			Description: "Senior engineer from Google/Meta is interested in joining. Exceptional talent but expensive and expects senior role.",
			Cost:        200000, // $200k/year salary
			Benefit:     "Accelerates product development 2x, attracts other top talent, improved technical credibility",
			Risk:        "High salary, may create team dynamics issues if not managed well",
			ExpiresIn:   1,
		}

	case "competitor_distress":
		opp = StrategicOpportunity{
			Type:        "competitor_distress",
			Title:       "🎯 Competitor in Distress",
			Description: "Main competitor is struggling (layoffs, negative press). Perfect time to steal their customers or acquire them cheaply.",
			Cost:        50000 + rand.Int63n(150000),
			Benefit:     fmt.Sprintf("+%d customers (from their base), eliminate key competitor", 15+rand.Intn(25)),
			Risk:        "May inherit technical debt or unhappy customers",
			ExpiresIn:   2,
		}
	}

	return &opp
}

// UpdateEmployeeVesting updates vesting progress for all employees
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
func (fs *FounderState) ProcessMonth() []string {
	return fs.ProcessMonthWithBaseline(fs.MRR)
}

// ProcessMonthWithBaseline runs monthly calculations with a baseline MRR for comparison
func (fs *FounderState) ProcessMonthWithBaseline(baselineMRR int64) []string {
	var messages []string
	fs.Turn++

	// Update employee vesting
	fs.UpdateEmployeeVesting()

	// Ensure MRR is in sync before processing
	fs.syncMRR()

	// 1. Process revenue growth
	oldMRR := baselineMRR // Use the baseline MRR from start of turn, not current MRR

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
	// Calculate growth first, then apply proportionally
	newRevenue := int64(float64(fs.MRR) * actualGrowth)

	// Apply growth proportionally to direct and affiliate MRR
	if fs.MRR > 0 {
		directRatio := float64(fs.DirectMRR) / float64(fs.MRR)
		affiliateRatio := float64(fs.AffiliateMRR) / float64(fs.MRR)
		fs.DirectMRR += int64(float64(newRevenue) * directRatio)
		fs.AffiliateMRR += int64(float64(newRevenue) * affiliateRatio)
	} else {
		// If no MRR yet, growth goes to direct
		fs.DirectMRR += newRevenue
	}

	// Sync MRR from DirectMRR + AffiliateMRR
	fs.syncMRR()

	// 2. Process churn (only if we have customers)
	var lostCustomers int
	var lostDirectCustomers int
	var lostAffiliateCustomers int
	var actualChurn float64

	// Always recalculate base churn based on current product maturity
	// Lower maturity = higher churn
	// Formula: baseChurn = (1.0 - ProductMaturity) * 0.65 + 0.05
	// This ensures churn decreases as product matures
	baseChurnFromMaturity := (1.0-fs.ProductMaturity)*0.65 + 0.05
	// Cap at reasonable range (5% minimum, 70% maximum)
	baseChurnFromMaturity = math.Max(0.05, math.Min(0.70, baseChurnFromMaturity))
	baseChurn := baseChurnFromMaturity

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

	// Calculate effective churn rate (after CS team impact)
	actualChurn = math.Max(0.01, baseChurn-csImpact)

	// Update displayed churn rate using helper function
	fs.RecalculateChurnRate()

	if fs.Customers > 0 {
		// Get active customers by source for proportional churn
		activeDirectCustomers := []Customer{}
		activeAffiliateCustomers := []Customer{}
		for _, c := range fs.CustomerList {
			if c.IsActive {
				if c.Source == "direct" || c.Source == "partnership" || c.Source == "market" {
					activeDirectCustomers = append(activeDirectCustomers, c)
				} else if c.Source == "affiliate" {
					activeAffiliateCustomers = append(activeAffiliateCustomers, c)
				}
			}
		}

		// Calculate how many to churn from each source
		lostDirectCustomers = int(float64(len(activeDirectCustomers)) * actualChurn)
		lostAffiliateCustomers = int(float64(len(activeAffiliateCustomers)) * actualChurn)

		// Mark customers as churned (randomly select from active customers)
		rand.Shuffle(len(activeDirectCustomers), func(i, j int) {
			activeDirectCustomers[i], activeDirectCustomers[j] = activeDirectCustomers[j], activeDirectCustomers[i]
		})
		rand.Shuffle(len(activeAffiliateCustomers), func(i, j int) {
			activeAffiliateCustomers[i], activeAffiliateCustomers[j] = activeAffiliateCustomers[j], activeAffiliateCustomers[i]
		})

		// Churn direct customers
		for i := 0; i < lostDirectCustomers && i < len(activeDirectCustomers); i++ {
			customer := activeDirectCustomers[i]
			fs.churnCustomer(customer.ID)
			fs.DirectMRR -= customer.DealSize
		}

		// Churn affiliate customers
		for i := 0; i < lostAffiliateCustomers && i < len(activeAffiliateCustomers); i++ {
			customer := activeAffiliateCustomers[i]
			fs.churnCustomer(customer.ID)
			fs.AffiliateMRR -= customer.DealSize
		}

		lostCustomers = lostDirectCustomers + lostAffiliateCustomers
		fs.Customers -= lostCustomers
		fs.DirectCustomers -= lostDirectCustomers
		fs.AffiliateCustomers -= lostAffiliateCustomers
	}

	// Clamp values to prevent negatives
	if fs.DirectMRR < 0 {
		fs.DirectMRR = 0
	}
	if fs.AffiliateMRR < 0 {
		fs.AffiliateMRR = 0
	}
	if fs.Customers < 0 {
		fs.Customers = 0
	}
	if fs.DirectCustomers < 0 {
		fs.DirectCustomers = 0
	}
	if fs.AffiliateCustomers < 0 {
		fs.AffiliateCustomers = 0
	}

	// Ensure MRR is in sync with DirectMRR + AffiliateMRR
	fs.syncMRR()

	// Recalculate average deal size after churn
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}
	// If no customers, keep AvgDealSize from template (don't reset to 0)

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

	// Calculate infrastructure costs (compute + ODC)
	fs.CalculateInfrastructureCosts()
	totalInfrastructureCost := fs.MonthlyComputeCost + fs.MonthlyODCCost

	fs.Cash -= totalCost
	fs.Cash -= totalInfrastructureCost

	netIncome := netMRRToCash - totalCost - totalInfrastructureCost

	// 5. Update runway
	fs.CalculateRunway()

	// Store messages for churn, cash flow, etc. (will add MRR comparison later)
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

	// 7. Process advanced features (affiliates, partnerships, etc.)
	// These will add more MRR, so we'll compare baseline to final MRR after all processing
	partnershipMsgs := fs.UpdatePartnerships()
	messages = append(messages, partnershipMsgs...)

	affiliateMsgs := fs.UpdateAffiliateProgram()
	messages = append(messages, affiliateMsgs...)

	competitorMsgs := fs.UpdateCompetitors()
	messages = append(messages, competitorMsgs...)

	marketMsgs := fs.UpdateGlobalMarkets()
	messages = append(messages, marketMsgs...)

	// 8. Spawn new competitors randomly
	if newComp := fs.SpawnCompetitor(); newComp != nil {
		messages = append(messages, fmt.Sprintf("🚨 NEW COMPETITOR: %s entered the market! Threat: %s, Market Share: %.1f%%",
			newComp.Name, newComp.Threat, newComp.MarketShare*100))
	}

	// 9. Process random events
	eventMsgs := fs.ProcessRandomEvents()
	messages = append(messages, eventMsgs...)

	// 10. Spawn new random events (5% chance each month)
	if rand.Float64() < 0.05 {
		if event := fs.SpawnRandomEvent(); event != nil {
			messages = append(messages, fmt.Sprintf("⚡ EVENT: %s - %s", event.Title, event.Description))
		}
	}

	// Final sync to ensure MRR is always correct (after all processing)
	fs.syncMRR()

	// Recalculate infrastructure costs after all customer changes (affiliates, partnerships, etc.)
	fs.CalculateInfrastructureCosts()

	// Show infrastructure costs if significant
	if fs.MonthlyComputeCost > 0 || fs.MonthlyODCCost > 0 {
		messages = append(messages, fmt.Sprintf("💻 Infrastructure: Compute $%s/mo, ODC $%s/mo",
			formatCurrency(fs.MonthlyComputeCost), formatCurrency(fs.MonthlyODCCost)))
	}

	// 11. Update growth rate for display on next month's dashboard (AFTER all processing)
	if oldMRR > 0 {
		fs.MonthlyGrowthRate = float64(fs.MRR-oldMRR) / float64(oldMRR)
	} else if fs.MRR > 0 {
		// First customers! Set initial growth rate
		fs.MonthlyGrowthRate = 0.10 // Start with 10% base growth
	}

	// Update board sentiment if raised funding
	fs.UpdateBoardSentiment()

	// Get board guidance
	boardGuidance := fs.GetBoardGuidance()
	messages = append(messages, boardGuidance...)

	// Generate strategic opportunity (15% chance, or 25% with good board)
	if fs.PendingOpportunity == nil {
		newOpp := fs.GenerateStrategicOpportunity()
		if newOpp != nil {
			fs.PendingOpportunity = newOpp
		}
	} else {
		// Decrement expiration timer
		fs.PendingOpportunity.ExpiresIn--
		if fs.PendingOpportunity.ExpiresIn <= 0 {
			messages = append(messages, fmt.Sprintf("⏰ Opportunity expired: %s", fs.PendingOpportunity.Title))
			fs.PendingOpportunity = nil
		}
	}

	// 12. Generate MRR comparison message (AFTER all customer additions)
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

	// Recalculate average deal size if we have customers
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}
	// If no customers, keep AvgDealSize from template (don't reset to 0)

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
	// Ensure MRR is synced before calculating boost
	fs.syncMRR()

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
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - helps acquire first customers
			mrrBoost = 5000 + rand.Int63n(15000) // $5-20k/month minimum
		}
		churnReduction = 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
		duration = 12 + rand.Intn(12)               // 12-24 months
	case "technology":
		cost = 30000 + rand.Int63n(70000)                                // $30-100k
		mrrBoost = int64(float64(fs.MRR) * (0.05 + rand.Float64()*0.15)) // 5-20% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - product integration helps attract customers
			mrrBoost = 3000 + rand.Int63n(7000) // $3-10k/month minimum
		}
		churnReduction = 0.02 + rand.Float64()*0.03 // 2-5% churn reduction
		duration = 12 + rand.Intn(24)               // 12-36 months
	case "co-marketing":
		cost = 25000 + rand.Int63n(50000)                                // $25-75k
		mrrBoost = int64(float64(fs.MRR) * (0.15 + rand.Float64()*0.25)) // 15-40% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - marketing helps acquire customers
			mrrBoost = 8000 + rand.Int63n(12000) // $8-20k/month minimum
		}
		churnReduction = 0.005 + rand.Float64()*0.015 // 0.5-2% churn reduction
		duration = 6 + rand.Intn(12)                  // 6-18 months
	case "data":
		cost = 40000 + rand.Int63n(60000)                                // $40-100k
		mrrBoost = int64(float64(fs.MRR) * (0.08 + rand.Float64()*0.12)) // 8-20% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - analytics help attract customers
			mrrBoost = 4000 + rand.Int63n(8000) // $4-12k/month minimum
		}
		churnReduction = 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
		duration = 12 + rand.Intn(24)               // 12-36 months
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

	// Apply partnership benefits immediately
	fs.MRR += mrrBoost
	fs.DirectMRR += mrrBoost // Partnership boost goes to direct MRR
	fs.CustomerChurnRate -= churnReduction
	if fs.CustomerChurnRate < 0 {
		fs.CustomerChurnRate = 0
	}

	// Sync MRR to ensure consistency
	fs.syncMRR()

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
		// Calculate MRR with variable deal sizes
		// Use template AvgDealSize if current AvgDealSize is 0 (no customers yet)
		baseDealSize := fs.AvgDealSize
		if baseDealSize == 0 {
			// Fallback to category-based defaults
			switch fs.Category {
			case "SaaS":
				baseDealSize = 1000 // Default $1k/month for SaaS
			case "DeepTech":
				baseDealSize = 5000 // Default $5k/month for DeepTech
			case "GovTech":
				baseDealSize = 2000 // Default $2k/month for GovTech
			case "Hardware":
				baseDealSize = 3000 // Default $3k/month for Hardware
			default:
				baseDealSize = 1000 // Default $1k/month
			}
		}

		var totalMRR int64
		var dealSizes []int64 // Store deal sizes for customer tracking
		for i := 0; i < newCustomers; i++ {
			dealSize := generateDealSize(baseDealSize, fs.Category)
			fs.updateDealSizeRange(dealSize)
			totalMRR += dealSize
			dealSizes = append(dealSizes, dealSize)
		}

		commissionPaid := int64(float64(totalMRR) * prog.Commission)

		// These are affiliate customers
		fs.Customers += newCustomers
		fs.AffiliateCustomers += newCustomers
		fs.AffiliateMRR += totalMRR
		fs.Cash -= commissionPaid

		// Add customers to tracking system
		for _, dealSize := range dealSizes {
			fs.addCustomer(dealSize, "affiliate")
		}

		prog.CustomersAcquired += newCustomers
		prog.MonthlyRevenue += totalMRR

		// Sync MRR from DirectMRR + AffiliateMRR
		fs.syncMRR()

		// Recalculate average deal size
		if fs.Customers > 0 {
			fs.AvgDealSize = fs.MRR / int64(fs.Customers)
		}

		messages = append(messages, fmt.Sprintf("🤝 Affiliates brought %d customers ($%s MRR, $%s commission)",
			newCustomers, formatCurrency(totalMRR), formatCurrency(commissionPaid)))

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

		// Calculate MRR with variable deal sizes
		var totalMRR int64
		var dealSizes []int64 // Store deal sizes for customer tracking
		for i := 0; i < newCustomers; i++ {
			dealSize := generateDealSize(fs.AvgDealSize, fs.Category)
			fs.updateDealSizeRange(dealSize)
			totalMRR += dealSize
			dealSizes = append(dealSizes, dealSize)
		}

		// These are direct customers (acquired via partnership)
		fs.Customers += newCustomers
		fs.DirectCustomers += newCustomers
		fs.DirectMRR += totalMRR

		// Add customers to tracking system
		for _, dealSize := range dealSizes {
			fs.addCustomer(dealSize, "partnership")
		}

		// Sync MRR from DirectMRR + AffiliateMRR
		fs.syncMRR()

		// Recalculate average deal size
		if fs.Customers > 0 {
			fs.AvgDealSize = fs.MRR / int64(fs.Customers)
		}

		return fmt.Sprintf("Partnered with %s! Cost: $%s. Gained %d customers (+$%s MRR)",
			comp.Name, formatCurrency(cost), newCustomers, formatCurrency(totalMRR)), nil

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

	// Calculate MRR with variable deal sizes
	var initialMRR int64
	var dealSizes []int64 // Store deal sizes for customer tracking
	for i := 0; i < initialCustomers; i++ {
		dealSize := generateDealSize(fs.AvgDealSize, fs.Category)
		fs.updateDealSizeRange(dealSize)
		initialMRR += dealSize
		dealSizes = append(dealSizes, dealSize)
	}

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
	// These are direct customers (from market expansion)
	fs.Customers += initialCustomers
	fs.DirectCustomers += initialCustomers
	fs.DirectMRR += initialMRR

	// Add customers to tracking system
	for _, dealSize := range dealSizes {
		fs.addCustomer(dealSize, "market")
	}

	// Sync MRR from DirectMRR + AffiliateMRR
	fs.syncMRR()

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

	// Increase global churn rate due to operational complexity
	fs.CustomerChurnRate += 0.01 + (rand.Float64() * 0.01)

	// Add regional competitors based on competition level
	numCompetitors := 0
	switch competition {
	case "very_high":
		numCompetitors = 2 + rand.Intn(2) // 2-3 competitors
	case "high":
		numCompetitors = 1 + rand.Intn(2) // 1-2 competitors
	case "medium":
		numCompetitors = rand.Intn(2) // 0-1 competitors
	case "low":
		numCompetitors = 0 // No competitors
	}

	// Regional competitor names by market
	regionalCompetitors := map[string][]string{
		"Europe":      {"Zalando", "Klarna", "N26", "TransferWise", "BlaBlaCar", "Deliveroo EU"},
		"Asia":        {"Grab", "GoJek", "Alibaba Local", "Meituan", "Tokopedia", "Paytm"},
		"LATAM":       {"Nubank", "Mercado Libre", "Rappi", "Kavak", "Creditas", "QuintoAndar"},
		"Middle East": {"Careem", "Souq", "Fetchr", "Talabat", "Noon", "Swvl"},
		"Africa":      {"Jumia", "Flutterwave", "Andela", "Paystack", "Konga", "M-Pesa"},
		"Australia":   {"Afterpay", "Canva AU", "Atlassian", "WiseTech", "Xero", "SEEK"},
	}

	if names, ok := regionalCompetitors[region]; ok && numCompetitors > 0 {
		for i := 0; i < numCompetitors && i < len(names); i++ {
			compName := names[rand.Intn(len(names))]
			
			// Check if competitor already exists
			exists := false
			for _, existing := range fs.Competitors {
				if existing.Name == compName {
					exists = true
					break
				}
			}
			
			if !exists {
				threatLevel := "medium"
				marketShare := 0.05 + rand.Float64()*0.15 // 5-20% market share
				
				switch competition {
				case "very_high":
					threatLevel = "high"
					marketShare = 0.10 + rand.Float64()*0.20 // 10-30%
				case "high":
					threatLevel = "high"
					marketShare = 0.08 + rand.Float64()*0.15 // 8-23%
				case "medium":
					threatLevel = "medium"
					marketShare = 0.05 + rand.Float64()*0.10 // 5-15%
				}
				
				competitor := Competitor{
					Name:          compName + " (" + region + ")",
					Threat:        threatLevel,
					MarketShare:   marketShare,
					Strategy:      "ignore", // Default strategy
					MonthAppeared: fs.Turn,
					Active:        true,
				}
				fs.Competitors = append(fs.Competitors, competitor)
			}
		}
	}

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

		// Count CS team assigned to this market or "All" markets
		csAssignedToMarket := 0
		for _, cs := range fs.Team.CustomerSuccess {
			if cs.AssignedMarket == m.Region || cs.AssignedMarket == "All" {
				csAssignedToMarket++
			}
		}
		
		// COO also helps with churn across all markets
		hasCOO := false
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCOO {
				hasCOO = true
				break
			}
		}

		// No CS team assigned to this market? Much higher churn (up to 50% in new markets)
		if csAssignedToMarket == 0 && !hasCOO {
			marketChurn += 0.30 // +30% base churn without CS in this market
		} else {
			// Each CS rep assigned to market reduces churn by ~2%
			churnReduction := float64(csAssignedToMarket) * 0.02
			if hasCOO {
				churnReduction += 0.06 // COO equivalent to 3 CS reps
			}
			marketChurn -= churnReduction
			if marketChurn < 0.01 {
				marketChurn = 0.01 // Min 1% churn
			}
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
			// Estimate MRR lost using average deal size
			mrrLost := int64(customersLost) * fs.AvgDealSize

			m.CustomerCount -= customersLost
			m.MRR -= mrrLost
			// These are direct customers (from market expansion)
			fs.MRR -= mrrLost
			fs.DirectMRR -= mrrLost
			fs.Customers -= customersLost
			fs.DirectCustomers -= customersLost

			// Recalculate average deal size
			if fs.Customers > 0 {
				fs.AvgDealSize = fs.MRR / int64(fs.Customers)
			}

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
					// Estimate MRR stolen using average deal size
					stolenMRR := int64(competitorSteal) * fs.AvgDealSize

					m.CustomerCount -= competitorSteal
					m.MRR -= stolenMRR
					// These are direct customers (from market expansion)
					fs.MRR -= stolenMRR
					fs.DirectMRR -= stolenMRR
					fs.Customers -= competitorSteal
					fs.DirectCustomers -= competitorSteal

					// Recalculate average deal size
					if fs.Customers > 0 {
						fs.AvgDealSize = fs.MRR / int64(fs.Customers)
					}

					messages = append(messages, fmt.Sprintf("⚠️  %s took %d customers in %s",
						comp.Name, competitorSteal, m.Region))
				}
			}
		}

		// Now attempt growth
		// Base growth rate - you need sales/marketing to grow in new markets
		baseGrowth := 0.03 + (rand.Float64() * 0.02) // 3-5% base monthly growth

		// Count employees assigned to this market or "All" markets
		salesInMarket := 0
		for _, sales := range fs.Team.Sales {
			if sales.AssignedMarket == m.Region || sales.AssignedMarket == "All" {
				salesInMarket++
			}
		}
		
		marketingInMarket := 0
		for _, marketing := range fs.Team.Marketing {
			if marketing.AssignedMarket == m.Region || marketing.AssignedMarket == "All" {
				marketingInMarket++
			}
		}
		
		csInMarket := 0
		for _, cs := range fs.Team.CustomerSuccess {
			if cs.AssignedMarket == m.Region || cs.AssignedMarket == "All" {
				csInMarket++
			}
		}

		// Sales team impact (critical for growing in new markets)
		salesImpact := float64(salesInMarket) * 0.05 // Each sales rep adds 5%
		
		// CGO amplifies sales impact (CGO works across all markets)
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCGO {
				salesImpact += (exec.Impact * 0.05) // CGO adds significant growth boost
			}
		}

		// Marketing team impact (brand awareness in new markets)
		marketingImpact := 0.03 * float64(marketingInMarket) // Each marketer adds 3%

		// Adjust for competition
		competitionMultiplier := 1.0
		switch m.LocalCompetition {
		case "very_high":
			competitionMultiplier = 0.6 // Harder growth in very competitive markets
		case "high":
			competitionMultiplier = 0.75
		case "medium":
			competitionMultiplier = 0.9
		case "low":
			competitionMultiplier = 1.2 // Much easier growth in low competition
		}

		// Product maturity affects conversion
		productMultiplier := fs.ProductMaturity
		if productMultiplier < 0.5 {
			productMultiplier = 0.5 // Can't grow much with immature product
		}

		totalGrowthRate := (baseGrowth + salesImpact + marketingImpact) * competitionMultiplier * productMultiplier

		// Calculate new customers
		// Percentage growth based on current base (compounds over time)
		percentageGrowth := int(float64(m.CustomerCount) * totalGrowthRate)
		
		// Plus absolute growth (helps new/small markets grow)
		// Sales/marketing teams directly acquire customers even in new markets
		// Only count employees assigned to this market or "All"
		absoluteGrowth := (salesInMarket * 2) + (marketingInMarket * 1)
		
		// CGO contributes to absolute growth (works across all markets)
		for _, exec := range fs.Team.Executives {
			if exec.Role == RoleCGO {
				absoluteGrowth += int(exec.Impact * 3) // CGO brings in customers directly
			}
		}
		
		newCustomers := percentageGrowth + absoluteGrowth

		// But cap at remaining market opportunity
		remainingMarket := m.MarketSize - m.CustomerCount
		if newCustomers > remainingMarket {
			newCustomers = remainingMarket
		}

		if newCustomers > 0 {
			// Calculate MRR with variable deal sizes
			var totalMRR int64
			for i := 0; i < newCustomers; i++ {
				dealSize := generateDealSize(fs.AvgDealSize, fs.Category)
				fs.updateDealSizeRange(dealSize)
				totalMRR += dealSize
			}

			m.CustomerCount += newCustomers
			m.MRR += totalMRR
			// These are direct customers (from market expansion)
			fs.DirectMRR += totalMRR
			fs.Customers += newCustomers
			fs.DirectCustomers += newCustomers

			// Sync MRR from DirectMRR + AffiliateMRR
			fs.syncMRR()

			// Recalculate average deal size
			if fs.Customers > 0 {
				fs.AvgDealSize = fs.MRR / int64(fs.Customers)
			}

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
	mrrLost := int64(customersLost) * fs.AvgDealSize

	// Apply churn proportionally to direct and affiliate customers
	directCustomersLost := int(float64(fs.DirectCustomers) * churnRate)
	affiliateCustomersLost := int(float64(fs.AffiliateCustomers) * churnRate)
	directMRRLost := int64(float64(fs.DirectMRR) * churnRate)
	affiliateMRRLost := int64(float64(fs.AffiliateMRR) * churnRate)

	fs.Customers -= customersLost
	fs.DirectCustomers -= directCustomersLost
	fs.AffiliateCustomers -= affiliateCustomersLost
	fs.MRR -= mrrLost
	fs.DirectMRR -= directMRRLost
	fs.AffiliateMRR -= affiliateMRRLost

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

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
	// Dilution happens automatically because founder equity = 100 - EquityGivenAway - EquityPool
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
