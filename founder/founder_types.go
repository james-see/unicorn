package founder

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
	Name           string
	Role           EmployeeRole
	MonthlyCost    int64
	Impact         float64 // Productivity/effectiveness multiplier
	IsExecutive    bool    // C-level executives have 3x impact, $300k/year salary
	Equity         float64 // Equity percentage owned by this employee
	VestingMonths  int     // Total vesting period (typically 48 months)
	CliffMonths    int     // Cliff period (typically 12 months)
	VestedMonths   int     // Months vested so far
	HasCliff       bool    // Has cliff been reached
	MonthHired     int     // Month when hired
	AssignedMarket string  // Market assignment: "USA", "Europe", "Asia", "All", etc.
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
	IsChairman        bool    // Whether this member is the chairman of the board
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
	ReferralProgram    *ReferralProgram
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
	HasExited                 bool
	ExitType                  string // "ipo", "acquisition", "secondary", "time_limit"
	ExitValuation             int64
	ExitMonth                 int
	MonthReachedProfitability int      // -1 if never profitable, otherwise the month when profitability was reached
	PlayerUpgrades            []string // Player's purchased upgrades
	HiresCount                int      // Track number of hires for Quick Hire upgrade
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
	IsCompetitor bool   // true if offer is from a competitor AI
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

// ReferralProgram represents a customer referral program
type ReferralProgram struct {
	LaunchedMonth      int
	RewardPerReferral  int64   // Cash reward per successful referral
	RewardType         string  // "cash", "credit", "equity"
	MonthlyBudget       int64   // Monthly budget for rewards
	ReferralsThisMonth int
	TotalReferrals     int
	CustomersAcquired  int
	MonthlyCost         int64   // Total monthly cost (rewards + platform fees)
	PlatformFee         int64   // Monthly platform/management fee
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

// MonthlyHighlight represents a monthly achievement or concern
type MonthlyHighlight struct {
	Type    string // "win" or "concern"
	Message string
	Icon    string
}

