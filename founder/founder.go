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
)

// Employee represents a team member
type Employee struct {
	Role        EmployeeRole
	MonthlyCost int64
	Impact      float64 // Productivity/effectiveness multiplier
}

// Team tracks all employees
type Team struct {
	Engineers        []Employee
	Sales            []Employee
	CustomerSuccess  []Employee
	Marketing        []Employee
	TotalMonthlyCost int64
	TotalEmployees   int
}

// StartupTemplate represents a startup idea from JSON
type StartupTemplate struct {
	ID                 string                    `json:"id"`
	Name               string                    `json:"name"`
	Tagline            string                    `json:"tagline"`
	Type               string                    `json:"type"`
	Description        string                    `json:"description"`
	InitialCash        int64                     `json:"initial_cash"`
	MonthlyBurn        int64                     `json:"monthly_burn"`
	InitialCustomers   int                       `json:"initial_customers"`
	InitialMRR         int64                     `json:"initial_mrr"`
	AvgDealSize        int64                     `json:"avg_deal_size"`
	BaseChurnRate      float64                   `json:"base_churn_rate"`
	TargetMarketSize   int                       `json:"target_market_size"`
	CompetitionLevel   string                    `json:"competition_level"`
	InitialTeam        map[string]int            `json:"initial_team"`
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
	
	// Growth metrics
	MonthlyGrowthRate       float64
	CustomerAcquisitionCost int64
	LifetimeValue           int64
	
	// Advanced features
	Partnerships      []Partnership
	AffiliateProgram  *AffiliateProgram
	Competitors       []Competitor
	GlobalMarkets     []Market
	PivotHistory      []Pivot
	EquityPool        float64 // Employee equity pool %
	InvestorBuybacks  []Buyback
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

// AcquisitionOffer represents an offer to buy the company
type AcquisitionOffer struct {
	Acquirer      string
	OfferAmount   int64
	Month         int
	DueDiligence  string // "good", "normal", "bad"
	TermsQuality  string // "excellent", "good", "poor"
}

// Partnership represents a strategic partnership
type Partnership struct {
	Partner       string
	Type          string // "distribution", "technology", "co-marketing", "data"
	MonthStarted  int
	Duration      int // months
	Cost          int64
	MRRBoost      int64
	ChurnReduction float64
	Status        string // "active", "expired", "terminated"
}

// AffiliateProgram represents an affiliate marketing program
type AffiliateProgram struct {
	Active            bool
	MonthStarted      int
	Commission        float64 // % of deal
	MonthlyCost       int64   // Platform + management
	Affiliates        int
	MonthlyRevenue    int64
	CustomersAcquired int
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
	Region          string // "North America", "Europe", "Asia", "LATAM", etc.
	LaunchMonth     int
	SetupCost       int64
	MonthlyCost     int64
	CustomerCount   int
	MRR             int64
	MarketSize      int
	Penetration     float64
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
	Month         int
	Investor      string // Which round (Seed, Series A, etc)
	EquityBought  float64
	PricePaid     int64
	Valuation     int64
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
		Turn:              1,
		MaxTurns:          60, // 5 years
		ProductMaturity:   0.3, // Start at 30% product maturity
		MarketPenetration: float64(template.InitialCustomers) / float64(template.TargetMarketSize),
		TargetMarketSize:  template.TargetMarketSize,
		CompetitionLevel:  template.CompetitionLevel,
		
		// Initialize advanced features
		Partnerships:     []Partnership{},
		AffiliateProgram: nil,
		Competitors:      []Competitor{},
		GlobalMarkets:    []Market{},
		PivotHistory:     []Pivot{},
		EquityPool:       10.0, // Start with 10% equity pool for employees
		InvestorBuybacks: []Buyback{},
		EquityGivenAway:   0.0,
		BoardSeats:        1, // Founder starts with 1 board seat
		MonthlyGrowthRate: 0.10, // Start with 10% monthly growth
	}
	
	// Initialize team from template
	fs.Team = Team{
		Engineers:        make([]Employee, template.InitialTeam["engineers"]),
		Sales:            make([]Employee, template.InitialTeam["sales"]),
		CustomerSuccess:  make([]Employee, template.InitialTeam["customer_success"]),
		Marketing:        make([]Employee, template.InitialTeam["marketing"]),
	}
	
	// Set up initial employees
	avgSalary := int64(100000)
	for i := range fs.Team.Engineers {
		fs.Team.Engineers[i] = Employee{
			Role:        RoleEngineer,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4, // 0.8-1.2x impact
		}
	}
	for i := range fs.Team.Sales {
		fs.Team.Sales[i] = Employee{
			Role:        RoleSales,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
		}
	}
	for i := range fs.Team.CustomerSuccess {
		fs.Team.CustomerSuccess[i] = Employee{
			Role:        RoleCustomerSuccess,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
		}
	}
	for i := range fs.Team.Marketing {
		fs.Team.Marketing[i] = Employee{
			Role:        RoleMarketing,
			MonthlyCost: avgSalary / 12,
			Impact:      0.8 + rand.Float64()*0.4,
		}
	}
	
	fs.CalculateTeamCost()
	fs.CalculateRunway()
	
	return fs
}

// CalculateTeamCost calculates total monthly team cost
func (fs *FounderState) CalculateTeamCost() {
	total := int64(0)
	count := 0
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
		fs.CashRunwayMonths = 999 // Effectively infinite if cashflow positive
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
	// Calculate valuation based on MRR
	multiple := 12.0 // Base 12x MRR multiple
	if fs.MRR > 1000000 {
		multiple = 15.0 // Higher multiple for scale
	}
	
	valuation = int64(float64(fs.MRR) * 12 * multiple)
	founderEquity = 100.0 - fs.EquityGivenAway
	
	if fs.Cash <= 0 {
		outcome = "Shut Down - Ran Out of Cash"
	} else if fs.MRR > 5000000 {
		outcome = "Unicorn! - Exceptional Success"
	} else if fs.MRR > 1000000 {
		outcome = "Major Success - Strong Growth"
	} else if fs.MRR > 100000 {
		outcome = "Growing Business - Solid Progress"
	} else {
		outcome = "Surviving - Keep Pushing"
	}
	
	return outcome, valuation, founderEquity
}

// HireEmployee adds a new team member
func (fs *FounderState) HireEmployee(role EmployeeRole) error {
	avgSalary := int64(100000)
	employee := Employee{
		Role:        role,
		MonthlyCost: avgSalary / 12,
		Impact:      0.8 + rand.Float64()*0.4,
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
	return nil
}

// FireEmployee removes a team member
func (fs *FounderState) FireEmployee(role EmployeeRole) error {
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

// RaiseFunding attempts to raise a funding round
func (fs *FounderState) RaiseFunding(roundName string) (success bool, amount int64, terms string, equityGiven float64) {
	// Calculate valuation based on metrics
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
	var minValuation, targetRaise int64
	switch roundName {
	case "Seed":
		minValuation = 3000000
		targetRaise = 2000000 + rand.Int63n(3000000) // $2-5M
	case "Series A":
		minValuation = 15000000
		targetRaise = 10000000 + rand.Int63n(10000000) // $10-20M
	case "Series B":
		minValuation = 50000000
		targetRaise = 30000000 + rand.Int63n(20000000) // $30-50M
	default:
		return false, 0, "", 0
	}
	
	if baseValuation < minValuation {
		baseValuation = minValuation
	}
	
	// Due diligence affects terms
	dueDiligence := rand.Float64()
	if dueDiligence < 0.15 {
		// Bad due diligence - harsh terms
		terms = "Investor-heavy"
		equityGiven = (float64(targetRaise) / float64(baseValuation)) * 100 * 1.3
		targetRaise = int64(float64(targetRaise) * 0.7) // Get less money
	} else if dueDiligence > 0.85 {
		// Great due diligence - founder friendly
		terms = "Founder-friendly"
		equityGiven = (float64(targetRaise) / float64(baseValuation)) * 100 * 0.8
		targetRaise = int64(float64(targetRaise) * 1.2) // Get more money
	} else {
		// Normal terms
		terms = "Standard"
		equityGiven = (float64(targetRaise) / float64(baseValuation)) * 100
	}
	
	// Cap equity dilution
	if equityGiven > 30 {
		equityGiven = 30
	}
	
	fs.Cash += targetRaise
	fs.EquityGivenAway += equityGiven
	
	round := FundingRound{
		RoundName:   roundName,
		Amount:      targetRaise,
		Valuation:   baseValuation,
		EquityGiven: equityGiven,
		Month:       fs.Turn,
		Terms:       terms,
	}
	fs.FundingRounds = append(fs.FundingRounds, round)
	
	fs.CalculateRunway()
	
	return true, targetRaise, terms, equityGiven
}

// SpendOnMarketing allocates budget to customer acquisition
func (fs *FounderState) SpendOnMarketing(amount int64) int {
	if amount > fs.Cash {
		return 0
	}
	
	fs.Cash -= amount
	
	// Calculate customers acquired
	// CAC varies by competition level
	cac := int64(5000) // Base $5k CAC
	switch fs.CompetitionLevel {
	case "very_high":
		cac = 10000
	case "high":
		cac = 7500
	case "medium":
		cac = 5000
	case "low":
		cac = 3000
	}
	
	newCustomers := int(amount / cac)
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
	engImpact := 1.0
	for _, eng := range fs.Team.Engineers {
		engImpact += (eng.Impact * 0.05) // Each engineer adds ~5% product improvement
	}
	fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+(0.02*engImpact))
	
	// Sales team impact on growth
	salesImpact := 1.0
	for _, sales := range fs.Team.Sales {
		salesImpact += (sales.Impact * 0.1) // Each sales rep adds ~10% to close rate
	}
	
	// Marketing impact (residual from spend)
	baseGrowth := fs.MonthlyGrowthRate
	actualGrowth := baseGrowth * salesImpact * engImpact
	
	// Apply growth
	newRevenue := int64(float64(fs.MRR) * actualGrowth)
	fs.MRR += newRevenue
	
	// 2. Process churn
	baseChurn := fs.CustomerChurnRate
	
	// CS team reduces churn
	csImpact := 0.0
	for _, cs := range fs.Team.CustomerSuccess {
		csImpact += (cs.Impact * 0.02) // Each CS rep reduces churn by ~2%
	}
	actualChurn := math.Max(0.01, baseChurn-csImpact)
	
	churnLoss := int64(float64(fs.MRR) * actualChurn)
	fs.MRR -= churnLoss
	lostCustomers := int(float64(fs.Customers) * actualChurn)
	fs.Customers -= lostCustomers
	
	if fs.MRR < 0 {
		fs.MRR = 0
	}
	if fs.Customers < 0 {
		fs.Customers = 0
	}
	
	// 3. Calculate costs
	totalCost := fs.MonthlyTeamCost + (int64(fs.Team.TotalEmployees) * 2000) // +$2k overhead per employee
	fs.Cash -= totalCost
	
	netIncome := fs.MRR - totalCost
	
	// 4. Update runway
	fs.CalculateRunway()
	
	// 5. Update growth rate for next month
	if oldMRR > 0 {
		fs.MonthlyGrowthRate = float64(fs.MRR-oldMRR) / float64(oldMRR)
	}
	
	// 6. Generate messages
	if fs.MRR > oldMRR {
		pctGrowth := ((float64(fs.MRR) - float64(oldMRR)) / float64(oldMRR)) * 100
		messages = append(messages, fmt.Sprintf("💰 MRR grew %.1f%% to $%s", pctGrowth, formatCurrency(fs.MRR)))
	} else if fs.MRR < oldMRR {
		pctDecline := ((float64(oldMRR) - float64(fs.MRR)) / float64(oldMRR)) * 100
		messages = append(messages, fmt.Sprintf("⚠️  MRR declined %.1f%% to $%s", pctDecline, formatCurrency(fs.MRR)))
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
	
	// 7. Process advanced features
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
		"technology": {"AWS", "Google Cloud", "Microsoft Azure", "IBM", "MongoDB"},
		"co-marketing": {"Shopify", "Stripe", "Zendesk", "Slack", "Atlassian"},
		"data": {"Snowflake", "Databricks", "Tableau", "Segment", "Amplitude"},
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
		cost = 50000 + rand.Int63n(100000) // $50-150k
		mrrBoost = int64(float64(fs.MRR) * (0.1 + rand.Float64()*0.2)) // 10-30% MRR boost
		churnReduction = 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
		duration = 12 + rand.Intn(12) // 12-24 months
	case "technology":
		cost = 30000 + rand.Int63n(70000) // $30-100k
		mrrBoost = int64(float64(fs.MRR) * (0.05 + rand.Float64()*0.15)) // 5-20% MRR boost
		churnReduction = 0.02 + rand.Float64()*0.03 // 2-5% churn reduction
		duration = 12 + rand.Intn(24) // 12-36 months
	case "co-marketing":
		cost = 25000 + rand.Int63n(50000) // $25-75k
		mrrBoost = int64(float64(fs.MRR) * (0.15 + rand.Float64()*0.25)) // 15-40% MRR boost
		churnReduction = 0.005 + rand.Float64()*0.015 // 0.5-2% churn reduction
		duration = 6 + rand.Intn(12) // 6-18 months
	case "data":
		cost = 40000 + rand.Int63n(60000) // $40-100k
		mrrBoost = int64(float64(fs.MRR) * (0.08 + rand.Float64()*0.12)) // 8-20% MRR boost
		churnReduction = 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
		duration = 12 + rand.Intn(24) // 12-36 months
	}
	
	if cost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash for partnership (need $%s)", formatCurrency(cost))
	}
	
	partnership := Partnership{
		Partner:       partner,
		Type:          partnerType,
		MonthStarted:  fs.Turn,
		Duration:      duration,
		Cost:          cost,
		MRRBoost:      mrrBoost,
		ChurnReduction: churnReduction,
		Status:        "active",
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
		
		monthsActive := fs.Turn - p.MonthStarted
		if monthsActive >= p.Duration {
			p.Status = "expired"
			messages = append(messages, fmt.Sprintf("🤝 Partnership with %s has expired", p.Partner))
			continue
		}
		
		// Apply partnership benefits
		fs.MRR += p.MRRBoost / int64(p.Duration) // Distribute boost over duration
		fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate - (p.ChurnReduction / float64(p.Duration)))
	}
	
	return messages
}

// ============================================================================
// AFFILIATE PROGRAM
// ============================================================================

// LaunchAffiliateProgram starts an affiliate marketing program
func (fs *FounderState) LaunchAffiliateProgram(commission float64) error {
	if fs.AffiliateProgram != nil && fs.AffiliateProgram.Active {
		return fmt.Errorf("affiliate program already active")
	}
	
	if commission < 5 || commission > 30 {
		return fmt.Errorf("commission must be between 5%% and 30%%")
	}
	
	setupCost := int64(20000 + rand.Intn(30000)) // $20-50k setup
	if setupCost > fs.Cash {
		return fmt.Errorf("insufficient cash for affiliate program setup")
	}
	
	fs.Cash -= setupCost
	fs.AffiliateProgram = &AffiliateProgram{
		Active:       true,
		MonthStarted: fs.Turn,
		Commission:   commission,
		MonthlyCost:  5000 + int64(rand.Intn(5000)), // $5-10k/month platform fees
		Affiliates:   5 + rand.Intn(10), // Start with 5-15 affiliates
	}
	
	return nil
}

// UpdateAffiliateProgram processes monthly affiliate sales
func (fs *FounderState) UpdateAffiliateProgram() []string {
	var messages []string
	
	if fs.AffiliateProgram == nil || !fs.AffiliateProgram.Active {
		return messages
	}
	
	ap := fs.AffiliateProgram
	
	// Pay monthly platform costs
	fs.Cash -= ap.MonthlyCost
	
	// Affiliates grow over time
	if rand.Float64() < 0.3 {
		newAffiliates := 1 + rand.Intn(3)
		ap.Affiliates += newAffiliates
		messages = append(messages, fmt.Sprintf("🤝 +%d new affiliates joined (total: %d)", newAffiliates, ap.Affiliates))
	}
	
	// Generate affiliate sales (each affiliate brings ~0.5-2 customers/month)
	customersPerAffiliate := 0.5 + rand.Float64()*1.5
	newCustomers := int(float64(ap.Affiliates) * customersPerAffiliate)
	
	if newCustomers > 0 {
		revenue := int64(newCustomers) * fs.AvgDealSize
		commissionPaid := int64(float64(revenue) * ap.Commission / 100.0)
		
		fs.Customers += newCustomers
		fs.MRR += revenue
		fs.Cash -= commissionPaid
		
		ap.CustomersAcquired += newCustomers
		ap.MonthlyRevenue = revenue
		
		messages = append(messages, fmt.Sprintf("💰 Affiliates brought %d customers ($%s MRR, $%s commission)", 
			newCustomers, formatCurrency(revenue), formatCurrency(commissionPaid)))
	}
	
	return messages
}

// ============================================================================
// COMPETITORS
// ============================================================================

// SpawnCompetitor creates a new competitor
func (fs *FounderState) SpawnCompetitor() *Competitor {
	// 10% chance each month after month 6
	if fs.Turn < 6 || rand.Float64() > 0.1 {
		return nil
	}
	
	names := []string{
		"RivalTech", "FastGrowth Inc", "MarketLeader", "DisruptCo", 
		"NextGen Solutions", "AgileStartup", "InnovateLabs", "ScaleUp",
		"CompeteCorp", "ChallengerTech",
	}
	
	threats := []string{"low", "medium", "high"}
	weights := []float64{0.5, 0.35, 0.15}
	
	// Weighted random threat selection
	r := rand.Float64()
	var threat string
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			threat = threats[i]
			break
		}
	}
	
	marketShare := 0.05 + rand.Float64()*0.15 // 5-20% market share
	
	competitor := Competitor{
		Name:          names[rand.Intn(len(names))],
		Threat:        threat,
		MarketShare:   marketShare,
		Strategy:      "monitor", // Default strategy
		MonthAppeared: fs.Turn,
		Active:        true,
	}
	
	fs.Competitors = append(fs.Competitors, competitor)
	return &competitor
}

// HandleCompetitor responds to a competitor
func (fs *FounderState) HandleCompetitor(compIndex int, strategy string) (string, error) {
	if compIndex < 0 || compIndex >= len(fs.Competitors) {
		return "", fmt.Errorf("invalid competitor index")
	}
	
	comp := &fs.Competitors[compIndex]
	if !comp.Active {
		return "", fmt.Errorf("competitor no longer active")
	}
	
	var message string
	var cost int64
	
	switch strategy {
	case "ignore":
		comp.Strategy = "ignore"
		// Competitor may take market share
		lostCustomers := int(float64(fs.Customers) * comp.MarketShare * 0.5)
		fs.Customers -= lostCustomers
		fs.MRR -= int64(lostCustomers) * fs.AvgDealSize
		message = fmt.Sprintf("Ignored %s. Lost %d customers to competition", comp.Name, lostCustomers)
		
	case "compete":
		comp.Strategy = "compete"
		// Aggressive competition - costs money but reduces their threat
		cost = 50000 + int64(rand.Intn(100000)) // $50-150k
		if cost > fs.Cash {
			return "", fmt.Errorf("insufficient cash to compete aggressively")
		}
		fs.Cash -= cost
		
		// Reduce their market share
		comp.MarketShare *= 0.7
		if comp.MarketShare < 0.02 {
			comp.Active = false
			message = fmt.Sprintf("Successfully competed against %s! They shut down. Cost: $%s", 
				comp.Name, formatCurrency(cost))
		} else {
			// Lower threat level
			if comp.Threat == "high" {
				comp.Threat = "medium"
			} else if comp.Threat == "medium" {
				comp.Threat = "low"
			}
			message = fmt.Sprintf("Competed aggressively with %s. Reduced their threat to '%s'. Cost: $%s", 
				comp.Name, comp.Threat, formatCurrency(cost))
		}
		
	case "partner":
		comp.Strategy = "partner"
		// Partner with competitor - costs money but creates synergy
		cost = 100000 + int64(rand.Intn(150000)) // $100-250k
		if cost > fs.Cash {
			return "", fmt.Errorf("insufficient cash for partnership")
		}
		fs.Cash -= cost
		
		// Merge customer bases (partial)
		gainedCustomers := int(float64(comp.MarketShare) * float64(fs.TargetMarketSize) * 0.3)
		fs.Customers += gainedCustomers
		fs.MRR += int64(gainedCustomers) * fs.AvgDealSize
		
		comp.Active = false
		message = fmt.Sprintf("Partnered with %s! Gained %d customers. Cost: $%s", 
			comp.Name, gainedCustomers, formatCurrency(cost))
		
	default:
		return "", fmt.Errorf("unknown strategy: %s", strategy)
	}
	
	return message, nil
}

// UpdateCompetitors processes competitor actions
func (fs *FounderState) UpdateCompetitors() []string {
	var messages []string
	
	for i := range fs.Competitors {
		comp := &fs.Competitors[i]
		if !comp.Active {
			continue
		}
		
		// Competitors take actions based on threat level
		if comp.Strategy == "ignore" {
			// They keep taking market share
			if rand.Float64() < 0.3 {
				lostCustomers := int(float64(fs.Customers) * comp.MarketShare * 0.1)
				if lostCustomers > 0 {
					fs.Customers -= lostCustomers
					fs.MRR -= int64(lostCustomers) * fs.AvgDealSize
					messages = append(messages, fmt.Sprintf("⚠️  %s took %d customers", comp.Name, lostCustomers))
				}
			}
		}
		
		// Competitors may escalate threat
		if rand.Float64() < 0.05 {
			if comp.Threat == "low" {
				comp.Threat = "medium"
				messages = append(messages, fmt.Sprintf("⚠️  %s is growing! Threat increased to medium", comp.Name))
			} else if comp.Threat == "medium" {
				comp.Threat = "high"
				messages = append(messages, fmt.Sprintf("🚨 %s is now a high threat!", comp.Name))
			}
		}
	}
	
	return messages
}

// ============================================================================
// GLOBAL EXPANSION
// ============================================================================

// ExpandToMarket launches in a new geographic market
func (fs *FounderState) ExpandToMarket(region string) (*Market, error) {
	// Check if already in this market
	for _, m := range fs.GlobalMarkets {
		if m.Region == region {
			return nil, fmt.Errorf("already operating in %s", region)
		}
	}
	
	marketData := map[string]struct {
		setupCost        int64
		monthlyCost      int64
		marketSize       int
		competition      string
	}{
		"Europe": {200000, 30000, 50000, "high"},
		"Asia": {250000, 40000, 100000, "very_high"},
		"LATAM": {150000, 20000, 30000, "medium"},
		"Middle East": {180000, 25000, 20000, "low"},
		"Africa": {120000, 15000, 15000, "low"},
		"Australia": {100000, 18000, 10000, "medium"},
	}
	
	data, ok := marketData[region]
	if !ok {
		return nil, fmt.Errorf("unknown region: %s", region)
	}
	
	if data.setupCost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash for expansion (need $%s)", formatCurrency(data.setupCost))
	}
	
	fs.Cash -= data.setupCost
	
	market := Market{
		Region:          region,
		LaunchMonth:     fs.Turn,
		SetupCost:       data.setupCost,
		MonthlyCost:     data.monthlyCost,
		CustomerCount:   0,
		MRR:             0,
		MarketSize:      data.marketSize,
		Penetration:     0,
		LocalCompetition: data.competition,
	}
	
	fs.GlobalMarkets = append(fs.GlobalMarkets, market)
	
	return &market, nil
}

// UpdateGlobalMarkets processes all international markets
func (fs *FounderState) UpdateGlobalMarkets() []string {
	var messages []string
	
	for i := range fs.GlobalMarkets {
		m := &fs.GlobalMarkets[i]
		
		// Pay monthly costs
		fs.Cash -= m.MonthlyCost
		
		// Grow customers in market
		growthRate := 0.05 + rand.Float64()*0.1 // 5-15% monthly growth
		
		// Adjust for competition
		switch m.LocalCompetition {
		case "very_high":
			growthRate *= 0.6
		case "high":
			growthRate *= 0.8
		case "medium":
			growthRate *= 0.9
		}
		
		// Adjust for product maturity
		growthRate *= fs.ProductMaturity
		
		newCustomers := int(float64(m.MarketSize) * growthRate)
		if m.CustomerCount + newCustomers > m.MarketSize {
			newCustomers = m.MarketSize - m.CustomerCount
		}
		
		if newCustomers > 0 {
			m.CustomerCount += newCustomers
			newMRR := int64(newCustomers) * fs.AvgDealSize
			m.MRR += newMRR
			fs.MRR += newMRR // Add to global MRR
			fs.Customers += newCustomers
			
			m.Penetration = float64(m.CustomerCount) / float64(m.MarketSize)
			
			messages = append(messages, fmt.Sprintf("🌍 %s: +%d customers, $%s MRR (%.1f%% penetration)", 
				m.Region, newCustomers, formatCurrency(m.MRR), m.Penetration*100))
		}
	}
	
	return messages
}

// ============================================================================
// PIVOTS
// ============================================================================

// ExecutePivot changes company strategy or market focus
func (fs *FounderState) ExecutePivot(toStrategy string, reason string) (*Pivot, error) {
	strategies := []string{
		"Enterprise B2B", "SMB B2B", "B2C", "Marketplace", "Platform",
		"Vertical SaaS", "Horizontal SaaS", "Deep Tech", "Consumer Apps",
	}
	
	valid := false
	for _, s := range strategies {
		if s == toStrategy {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("invalid strategy: %s", toStrategy)
	}
	
	// Pivots are expensive and risky
	cost := 100000 + int64(rand.Intn(200000)) // $100-300k
	if cost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash for pivot (need $%s)", formatCurrency(cost))
	}
	
	// Lose some customers in transition
	customersLost := int(float64(fs.Customers) * (0.2 + rand.Float64()*0.3)) // Lose 20-50%
	
	// Success rate depends on product maturity and timing
	successChance := 0.3 + (fs.ProductMaturity * 0.4) // 30-70% success
	if fs.Turn > 36 {
		successChance *= 0.7 // Harder to pivot late
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
	
	fs.Cash -= cost
	fs.Customers -= customersLost
	fs.MRR -= int64(customersLost) * fs.AvgDealSize
	
	if success {
		// Successful pivot opens new market
		fs.TargetMarketSize = int(float64(fs.TargetMarketSize) * (1.5 + rand.Float64()))
		fs.StartupType = toStrategy
		fs.MonthlyGrowthRate *= 1.5 // Higher growth potential
	} else {
		// Failed pivot
		fs.MonthlyGrowthRate *= 0.7 // Slower growth
	}
	
	fs.PivotHistory = append(fs.PivotHistory, pivot)
	
	return &pivot, nil
}

// ============================================================================
// INVESTOR BUYBACKS
// ============================================================================

// BuybackEquity buys back equity from investors
func (fs *FounderState) BuybackEquity(roundName string, equityPercent float64) (*Buyback, error) {
	// Must be profitable
	monthlyProfit := fs.MRR - fs.MonthlyTeamCost
	if monthlyProfit <= 0 {
		return nil, fmt.Errorf("must be profitable to buy back equity")
	}
	
	// Find the round
	var foundRound *FundingRound
	for i := range fs.FundingRounds {
		if fs.FundingRounds[i].RoundName == roundName {
			foundRound = &fs.FundingRounds[i]
			break
		}
	}
	if foundRound == nil {
		return nil, fmt.Errorf("funding round not found: %s", roundName)
	}
	
	// Can't buy back more than they own
	if equityPercent > foundRound.EquityGiven {
		return nil, fmt.Errorf("can't buy back more equity than investors own")
	}
	
	// Calculate current valuation
	currentValuation := int64(float64(fs.MRR) * 12 * 12) // 12x ARR
	priceToPay := int64(float64(currentValuation) * equityPercent / 100.0)
	
	if priceToPay > fs.Cash {
		return nil, fmt.Errorf("insufficient cash for buyback (need $%s)", formatCurrency(priceToPay))
	}
	
	buyback := Buyback{
		Month:        fs.Turn,
		Investor:     roundName,
		EquityBought: equityPercent,
		PricePaid:    priceToPay,
		Valuation:    currentValuation,
	}
	
	fs.Cash -= priceToPay
	fs.EquityGivenAway -= equityPercent
	foundRound.EquityGiven -= equityPercent
	
	fs.InvestorBuybacks = append(fs.InvestorBuybacks, buyback)
	
	return &buyback, nil
}

// ============================================================================
// BOARD MANAGEMENT
// ============================================================================

// AddBoardSeat adds a new board seat (dilutes equity pool)
func (fs *FounderState) AddBoardSeat(reason string) error {
	// Each board seat costs ~2% from employee equity pool
	equityCost := 1.5 + rand.Float64()
	
	if fs.EquityPool < equityCost {
		return fmt.Errorf("insufficient equity pool (need %.1f%%)", equityCost)
	}
	
	fs.BoardSeats++
	fs.EquityPool -= equityCost
	
	return nil
}

// ExpandEquityPool increases employee equity pool (dilutes founders)
func (fs *FounderState) ExpandEquityPool(percentToAdd float64) error {
	if percentToAdd < 1 || percentToAdd > 10 {
		return fmt.Errorf("can only add 1-10%% to equity pool at once")
	}
	
	// This dilutes the founder
	fs.EquityGivenAway += percentToAdd
	fs.EquityPool += percentToAdd
	
	return nil
}

