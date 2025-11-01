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

