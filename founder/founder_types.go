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

	// PHASE 1: Strategic features
	ProductRoadmap      *ProductRoadmap
	SalesPipeline       *SalesPipeline
	PricingStrategy     *PricingStrategy
	ActiveExperiment    *PricingExperiment
	CustomerSegments    []CustomerSegment
	VerticalFocuses     []VerticalFocus
	SelectedICP         string // Selected segment focus
	SelectedVertical    string // Selected industry vertical

	// PHASE 2: Growth & intelligence features
	ContentProgram    *ContentProgram
	CSPlaybooks       []CSPlaybook
	CustomerHealthMap map[int]CustomerHealth // CustomerID -> Health
	CompetitiveIntel  *CompetitiveIntel

	// PHASE 3: Polish & realism features
	TechnicalDebt      *TechnicalDebt
	PRProgram          *PRProgram
	InvestorUpdates    []InvestorUpdate
	BoardRequests      []BoardRequest
	PendingBoardRequest *BoardRequest

	// ADVANCED GROWTH MECHANICS
	AcquisitionTargets []AcquisitionTarget
	Acquisitions        []Acquisition
	PlatformMetrics     *PlatformMetrics
	NetworkEffects      []NetworkEffect
	PartnershipIntegrations []PartnershipIntegration

	// CRISIS MANAGEMENT
	SecurityPosture     *SecurityPosture
	SecurityIncidents   []SecurityIncident
	ActiveSecurityIncident *SecurityIncident
	PRCrises            []PRCrisis
	ActivePRCrisis      *PRCrisis
	CrisisResponses     []CrisisResponse
	EconomicEvent       *EconomicEvent
	SurvivalStrategies  []SurvivalStrategy
	KeyPersonRisks      []KeyPersonRisk
	KeyPersonEvents     []KeyPersonEvent
	SuccessionPlans     []SuccessionPlan
	
	// Roadmap tracking for achievements
	CustomersLostDuringRoadmap int // Track customers churned while features were in progress
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
	Name            string
	Threat          string // "low", "medium", "high", "critical"
	MarketShare     float64
	Strategy        string // "ignore", "monitor", "compete", "partner"
	MonthAppeared   int
	Active          bool
	LastActionMonth int // Prevent multiple actions in same turn
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

// ============================================================================
// PHASE 1 FEATURES: Product Roadmap, Sales Pipeline, Pricing, Segments
// ============================================================================

// ProductFeature represents a feature that can be built
type ProductFeature struct {
	Name                string
	Category            string // "Integration", "Security", "Analytics", etc.
	EngineerMonths      int
	Cost                int64
	ChurnReduction      float64
	CloseRateIncrease   float64
	DealSizeIncrease    float64
	MarketAppealScore   int
	Status              string // "backlog", "in_progress", "completed"
	MonthStarted        int
	MonthCompleted      int
	DevelopmentProgress int // 0-100%
	AllocatedEngineers  int // Engineers currently working on this
}

// ProductRoadmap represents the product development pipeline
type ProductRoadmap struct {
	Features          []ProductFeature
	AvailableFeatures []ProductFeature // Template of possible features
	CompletedCount    int
	InProgressCount   int
	CompetitorLaunches []CompetitorFeatureLaunch // Track competitor feature launches
}

// CompetitorFeatureLaunch tracks when a competitor launched a feature
type CompetitorFeatureLaunch struct {
	FeatureName string // Feature name (e.g., "API", "Mobile App", "SSO")
	CompetitorName string
	MonthLaunched int
}

// Deal represents a sales opportunity
type Deal struct {
	ID              int
	CompanyName     string
	DealSize        int64
	Stage           string // "lead", "qualified", "demo", "negotiation", "closed_won", "closed_lost"
	CloseProbability float64
	DaysInStage     int
	RequiredActions []string
	AssignedSalesRep string
	MonthCreated    int
	LostReason      string // if closed_lost
	Segment         string // "Enterprise", "Mid-Market", "SMB", "Startup"
	Vertical        string // Industry vertical
}

// SalesPipeline represents the sales funnel and metrics
type SalesPipeline struct {
	ActiveDeals       []Deal
	ClosedDeals       []Deal
	LeadsPerMonth     int
	ConversionRates   map[string]float64 // stage -> conversion rate
	AverageDealSize   int64
	AverageSalesCycle int // days
	WinRate           float64
	TotalDealsCreated int
	NextDealID        int
}

// PricingStrategy represents the company's pricing model
type PricingStrategy struct {
	Model         string           // "freemium", "trial", "annual_upfront", "usage_based", "tiered"
	CurrentTier   map[string]int64 // "starter" -> $99, "pro" -> $299, "enterprise" -> $999
	IsAnnual      bool
	Discount      float64 // annual discount %
	ChangeHistory []PricingChange
}

// PricingChange represents a historical pricing change
type PricingChange struct {
	Month       int
	FromModel   string
	ToModel     string
	Reason      string
	Impact      string
}

// PricingExperiment represents an A/B pricing test
type PricingExperiment struct {
	Name          string
	Cost          int64
	StartMonth    int
	Duration      int
	TestStrategy  PricingStrategy
	Results       PricingResults
	IsComplete    bool
}

// PricingResults represents experiment outcomes
type PricingResults struct {
	ConversionRateChange float64
	AvgDealSizeChange    int64
	ChurnRateChange      float64
	Confidence           float64 // 0-1
}

// CustomerSegment represents a market segment
type CustomerSegment struct {
	Name                string // "Enterprise", "Mid-Market", "SMB", "Startup"
	AvgDealSize         int64
	ChurnRate           float64
	SalesCycle          int // months
	CAC                 int64
	FeatureRequirements []string
	Volume              int // current customers in this segment
}

// VerticalFocus represents industry specialization
type VerticalFocus struct {
	Industry         string // "FinTech", "HealthTech", "Retail", "Manufacturing"
	ICPMatch         float64 // 0-1, how well targeted
	MarketSize       int
	Competition      string
	SpecializedSales int // sales reps trained for this vertical
	CACReduction     float64 // benefit of specialization
	IsActive         bool
}

// ============================================================================
// PHASE 2 FEATURES: Content Marketing, CS Playbooks, Competitive Intel
// ============================================================================

// ContentProgram represents content marketing efforts
type ContentProgram struct {
	MonthlyBudget   int64
	ContentTypes    map[string]bool // "blog", "seo", "webinars", "ebooks", "case_studies"
	OrganicTraffic  int
	InboundLeads    int
	ContentQuality  float64 // 0-1
	SEOScore        int     // 0-100
	MonthsActive    int
	TotalInvestment int64
	CumulativeLeads int
	LaunchedMonth   int
}

// CSPlaybook represents a customer success program
type CSPlaybook struct {
	Name           string // "Onboarding", "Health Monitoring", "Upsell", "Renewal", "Churn Prevention"
	CSHeadcount    int
	MonthlyBudget  int64
	ToolCosts      int64
	ChurnReduction float64
	UpsellRate     float64
	NPSScore       int // 0-100
	Active         bool
	LaunchedMonth  int
}

// CustomerHealth represents individual customer health tracking
type CustomerHealth struct {
	CustomerID      int
	HealthScore     float64 // 0-100
	Risk            string  // "healthy", "at_risk", "critical"
	Interventions   []string
	UpsellPotential bool
	LastTouchpoint  int // months ago
}

// CompetitiveIntel represents intelligence gathering
type CompetitiveIntel struct {
	HasAnalyst      bool
	MonthlyBudget   int64
	IntelReports    []IntelReport
	BattleCards     []BattleCard
	WinLossInsights map[string]int // reason -> count
	AnalystSalary   int64
	LaunchedMonth   int
}

// IntelReport represents a competitor analysis report
type IntelReport struct {
	CompetitorName string
	Pricing        map[string]int64
	Features       []string
	Funding        string
	TeamSize       int
	RecentMoves    []string
	ThreatLevel    string
	Cost           int64
	Month          int
}

// BattleCard represents competitive positioning
type BattleCard struct {
	CompetitorName  string
	OurAdvantages   []string
	TheirAdvantages []string
	ResponseTactics []string
	WinRateBonus    float64 // bonus to close rate when used
	CreatedMonth    int
}

// ============================================================================
// PHASE 3 FEATURES: Tech Debt, PR, Enhanced Investor Relations
// ============================================================================

// TechnicalDebt represents accumulated technical debt
type TechnicalDebt struct {
	CurrentLevel       int     // 0-100
	VelocityImpact     float64 // multiplier on engineer productivity
	BugFrequency       float64 // increases churn
	SecurityRisks      int
	ScalingProblems    bool
	EngineerMorale     float64 // affects attrition
	RefactoringCosts   int64
	MonthsSinceRefactor int
}

// PRProgram represents PR and media relations
type PRProgram struct {
	HasPRFirm       bool
	MonthlyRetainer int64
	Campaigns       []PRCampaign
	MediaCoverage   []MediaCoverage
	BrandScore      int // 0-100
	LaunchedMonth   int
}

// PRCampaign represents a PR initiative
type PRCampaign struct {
	Type         string // "product_launch", "funding", "thought_leadership", "crisis_response"
	Cost         int64
	Duration     int
	TargetMedia  []string
	StartMonth   int
	Success      bool
	Impact       PRImpact
}

// MediaCoverage represents a press mention
type MediaCoverage struct {
	Outlet    string // "TechCrunch", "WSJ", "Trade Publication", "Podcast"
	Type      string // "positive", "negative", "neutral"
	Reach     int
	CACImpact float64
	Month     int
}

// PRImpact represents effects of PR
type PRImpact struct {
	CACReduction   float64
	BrandBoost     int
	InboundLeads   int
	DurationMonths int
}

// InvestorUpdate represents monthly investor communication
type InvestorUpdate struct {
	Month        int
	Metrics      map[string]string // what you chose to share
	Transparency string            // "full", "optimistic", "selective"
	BoardResponse string
}

// BoardRequest represents asking the board for help
type BoardRequest struct {
	Type        string // "customer_intro", "recruiting_help", "strategic_advice", "fundraising_prep"
	Description string
	Month       int
	Result      BoardHelp
	Completed   bool
}

// BoardHelp represents value provided by the board
type BoardHelp struct {
	Type           string
	Benefit        string
	CACImpact      float64
	LeadsGenerated int
	IntrosMade     []string
}

// ============================================================================
// ADVANCED GROWTH MECHANICS: Acquisitions, Platform Effects, Partnerships
// ============================================================================

// AcquisitionTarget represents a potential acquisition target
type AcquisitionTarget struct {
	Name            string
	Category        string
	MRR             int64
	Customers       int
	TeamSize        int
	Technology      []string // Features/IP you gain
	AcquisitionCost int64
	IntegrationCost int64
	SynergyBonus    float64 // Revenue boost from integration
	Risk            string  // "low", "medium", "high"
	MonthAppeared   int
	ExpiresIn       int // Months until opportunity expires
}

// Acquisition represents a completed acquisition
type Acquisition struct {
	TargetName        string
	Month             int
	Cost              int64
	CustomersGained   int
	MRRGained         int64
	TeamGained        int
	Success           bool
	IntegrationMonths int
	IntegrationProgress int // 0-100%
	SynergyRealized   float64 // Actual synergy achieved
}

// PlatformMetrics represents platform business metrics
type PlatformMetrics struct {
	IsPlatform         bool
	ThirdPartyApps     int
	DeveloperCount     int
	APIUsage           int64 // API calls per month
	MarketplaceRevenue int64 // Revenue from marketplace fees
	NetworkEffectScore float64 // 0-1, how strong network effects are
	PlatformType       string // "marketplace", "social", "data", "infrastructure"
	LaunchedMonth      int
}

// NetworkEffect represents network effect mechanics
type NetworkEffect struct {
	Type          string  // "marketplace", "social", "data", "infrastructure"
	Strength      float64 // 0-1
	CustomerValue float64 // How much value each new customer adds
	Threshold     int     // Customers needed to activate
	Active        bool
	ActivatedMonth int
}

// Enhanced Partnership (extends existing Partnership struct)
// Note: Partnership struct already exists, we'll add fields via composition in FounderState

// PartnershipIntegration represents deep integration details
type PartnershipIntegration struct {
	PartnerName     string
	IntegrationDepth string // "surface", "deep", "native"
	RevenueShare    float64 // % of deals from partnership
	CACReduction    float64 // % reduction in CAC
	ChurnReduction  float64 // % reduction in churn
	MRRContribution int64   // MRR directly from partnership
	CoMarketingActive bool
	DataSharingActive bool
}

// ============================================================================
// CRISIS MANAGEMENT: Security, PR Crisis, Economy, Key Person Risk
// ============================================================================

// SecurityIncident represents a security breach or incident
type SecurityIncident struct {
	Type             string  // "data_breach", "ransomware", "ddos", "insider_threat", "vulnerability"
	Severity         string  // "low", "medium", "high", "critical"
	Month             int
	CustomersAffected int
	DataExposed       string  // "PII", "financial", "health", "none"
	ResponseCost      int64
	LegalCosts        int64
	ReputationDamage  float64 // 0-1
	Resolved          bool
	ResolutionMonth   int
	ResponseActions   []string
}

// SecurityPosture represents overall security status
type SecurityPosture struct {
	SecurityScore    int      // 0-100
	ComplianceCerts  []string // "SOC2", "ISO27001", "HIPAA", "GDPR"
	SecurityTeamSize int
	SecurityBudget    int64
	LastAudit         int
	Vulnerabilities   int
	BugBountyActive   bool
	BugBountyBudget   int64
}

// PRCrisis represents a PR crisis event
type PRCrisis struct {
	Type            string   // "scandal", "product_failure", "layoffs", "founder_drama", "competitor_attack"
	Severity        string   // "minor", "moderate", "major", "critical"
	Month           int
	MediaCoverage    []string // Outlets covering the story
	Response        string   // "none", "deny", "apologize", "transparent", "aggressive"
	ResponseCost     int64
	DurationMonths   int
	CACImpact        float64 // Multiplier on CAC
	ChurnImpact      float64 // Additional churn %
	BrandDamage      float64 // 0-1
	Resolved         bool
	ResolutionMonth  int
}

// CrisisResponse represents response to a PR crisis
type CrisisResponse struct {
	CrisisType    string
	ResponseType  string
	Cost          int64
	Effectiveness float64 // 0-1
	Outcome       string  // "contained", "escalated", "resolved"
	Month         int
}

// EconomicEvent represents economic downturn or market crash
type EconomicEvent struct {
	Type             string  // "recession", "market_crash", "funding_winter", "sector_crash"
	Severity         string  // "mild", "moderate", "severe", "extreme"
	Month             int
	DurationMonths    int
	GrowthImpact      float64 // Multiplier on growth (0.5 = 50% reduction)
	CACImpact         float64 // Multiplier on CAC
	ChurnImpact       float64 // Additional churn %
	FundingImpact     float64 // Multiplier on funding availability
	CustomerBudgetCut float64 // % reduction in customer budgets
	Active            bool
}

// SurvivalStrategy represents strategies to survive economic downturn
type SurvivalStrategy struct {
	Strategy      string   // "cut_costs", "pivot", "downround", "extend_runway", "acquire"
	Cost          int64
	Effectiveness float64 // 0-1
	Tradeoffs     []string // What you lose
	MonthStarted  int
	Active        bool
}

// KeyPersonRisk represents risk associated with key personnel
type KeyPersonRisk struct {
	PersonName      string   // "founder", "cto", "cfo", "head_of_sales"
	Role            string
	RiskLevel       string   // "low", "medium", "high", "critical"
	Dependencies    []string // What breaks if they leave
	SuccessionReady bool
	RetentionScore  float64 // 0-1, likelihood of staying
}

// KeyPersonEvent represents a key person leaving or crisis
type KeyPersonEvent struct {
	PersonName      string
	EventType       string // "quit", "poached", "illness", "scandal", "death"
	Month           int
	Impact          EventImpact
	ReplacementCost int64
	RecoveryMonths  int
	Resolved        bool
}

// SuccessionPlan represents succession planning for key persons
type SuccessionPlan struct {
	PersonName     string
	BackupPerson   string
	TrainingMonths int
	Ready          bool
	MonthCreated   int
}

