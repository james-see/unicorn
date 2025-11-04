package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// FundingRound represents a funding round for a startup
type FundingRound struct {
	RoundName        string // "Seed", "Series A", "Series B", etc.
	PreMoneyVal      int64  // Pre-money valuation
	InvestmentAmount int64  // Total raised in this round
	PostMoneyVal     int64  // Post-money valuation
	Month            int    // Game turn when this happened
}

// InvestmentTerms represents the terms of an investment deal
type InvestmentTerms struct {
	Type               string  // "Common", "Preferred", "SAFE", "Convertible"
	HasProRataRights   bool    // Right to participate in future rounds
	HasInfoRights      bool    // Right to company information
	HasBoardSeat       bool    // Board seat (for larger investments)
	LiquidationPref    float64 // Liquidation preference (1x, 2x, etc.)
	HasAntiDilution    bool    // Anti-dilution protection
	ConversionDiscount float64 // Discount on conversion (for SAFE/Convertible)
}

// Investment represents a player's investment in a startup
type Investment struct {
	CompanyName      string
	AmountInvested   int64   // Total invested across all rounds
	EquityPercent    float64 // Current equity after dilution
	InitialEquity    float64 // Original equity from first investment
	InitialValuation int64
	CurrentValuation int64
	MonthsHeld       int
	Category         string
	NegativeNewsSent bool            // Track if we've already sent negative news for this investment
	Rounds           []FundingRound  // Track all funding rounds
	Terms            InvestmentTerms // Investment terms
	FollowOnThisTurn bool            // Track if follow-on investment was made this turn (prevents double dilution)
}

// Portfolio tracks all player investments
type Portfolio struct {
	Cash                  int64
	Investments           []Investment
	NetWorth              int64
	Turn                  int
	MaxTurns              int
	InitialFundSize       int64   // Original fund size
	ManagementFeesCharged int64   // Total management fees paid
	AnnualManagementFee   float64 // Annual management fee rate (e.g., 0.02 for 2%)
	FollowOnReserve       int64   // Reserve fund for follow-on investments
}

// Startup represents a company available for investment
type Startup struct {
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	Category               string  `json:"category"`
	Valuation              int64   `json:"valuation"`
	GrossBurnRate          int     `json:"grossburnrate"`
	MonthlyActivationRate  int     `json:"Monthly Activation Rate"`
	MonthlyWebsiteVisitors int     `json:"Monthly Active Visitors"`
	MonthlySales           int     `json:"Monthly Sales"`
	Cost                   int     `json:"Cost"`
	SalePrice              int     `json:"Sale Price"`
	PercentMargin          int     `json:"Percent Margin Per Unit"`
	RiskScore              float64 // 0-1, higher is riskier
	GrowthPotential        float64 // 0-1, higher is better

	// Financial tracking
	MonthlyRevenue          int64   // Revenue this month
	MonthlyCosts            int64   // Costs this month (burn rate)
	NetIncome               int64   // Profit/Loss this month
	CumulativeRevenue       int64   // Total revenue to date
	CumulativeCosts         int64   // Total costs to date
	Last409AValuation       int64   // Last 409A valuation
	Last409AMonth           int     // When was last 409A done
	RevenueGrowthRate       float64 // Month-over-month growth
	CustomerCount           int     // Current customers
	MonthlyRecurringRevenue int64   // MRR for SaaS companies
}

// GameEvent represents something that happens to a startup
type GameEvent struct {
	Event       string  `json:"event"`
	Change      float64 `json:"change"` // multiplier (1.5 = +50%, 0.8 = -20%)
	Description string  `json:"description"`
}

// Difficulty represents game difficulty level
type Difficulty struct {
	Name           string
	StartingCash   int64
	EventFrequency float64 // 0-1, chance of event per turn
	Volatility     float64 // 0-1, market volatility
	MaxTurns       int
	Description    string
}

// AIPlayer represents a computer-controlled VC
type AIPlayer struct {
	Name          string
	Firm          string
	Portfolio     Portfolio
	Strategy      string  // "aggressive", "balanced", "conservative"
	RiskTolerance float64 // 0-1
}

// GameState holds the entire game state
type GameState struct {
	PlayerName         string
	Portfolio          Portfolio
	AvailableStartups  []Startup
	EventPool          []GameEvent
	Difficulty         Difficulty
	AIPlayers          []AIPlayer          // Computer opponents
	FundingRoundQueue  []FundingRoundEvent // Scheduled future funding rounds
	AcquisitionQueue   []AcquisitionEvent  // Scheduled acquisition offers
	DramaticEventQueue []DramaticEvent     // Scheduled dramatic events (scandals, splits, etc.)
	PendingBoardVotes  []BoardVote         // Board votes requiring player input
}

// FundingRoundEvent represents a scheduled funding round
type FundingRoundEvent struct {
	CompanyName   string
	RoundName     string
	ScheduledTurn int
	RaiseAmount   int64
	IsDownRound   bool // True if this is a down round
}

// AcquisitionEvent represents a potential acquisition offer
type AcquisitionEvent struct {
	CompanyName   string
	ScheduledTurn int
	OfferMultiple float64 // Multiple of EBITDA or revenue
	DueDiligence  string  // "good", "bad", "normal"
}

// DramaticEvent represents scandals, co-founder splits, etc.
type DramaticEvent struct {
	CompanyName   string
	ScheduledTurn int
	EventType     string  // "cofounder_split", "scandal", "lawsuit", "pivot_fail", "fraud", "data_breach", "key_hire_quit"
	Severity      string  // "minor", "moderate", "severe"
	ImpactPercent float64 // Valuation impact as a multiplier (0.3 = 70% drop)
}

// BoardVote represents a voting opportunity for board members
type BoardVote struct {
	CompanyName  string
	VoteType     string // "acquisition", "down_round", "strategic_pivot", "ceo_removal"
	Title        string
	Description  string
	OptionA      string                 // "Accept" / "Approve" / "Yes"
	OptionB      string                 // "Reject" / "Disapprove" / "No"
	ConsequenceA string                 // What happens if A wins
	ConsequenceB string                 // What happens if B wins
	RequiresVote bool                   // Whether player vote is required
	Turn         int                    // Turn when vote occurs
	Metadata     map[string]interface{} // Additional data (offer value, round details, etc.)
}

// FollowOnOpportunity represents a chance to invest more in a company raising a round
type FollowOnOpportunity struct {
	CompanyName   string
	RoundName     string
	PreMoneyVal   int64
	PostMoneyVal  int64
	CurrentEquity float64
	MinInvestment int64
	MaxInvestment int64
}

// Predefined difficulty levels
var (
	EasyDifficulty = Difficulty{
		Name:           "Easy",
		StartingCash:   1000000, // $1M fund
		EventFrequency: 0.20,    // 20% chance
		Volatility:     0.03,    // 3% volatility
		MaxTurns:       60,      // 5 years
		Description:    "$1M fund, lower volatility, 5 years",
	}

	MediumDifficulty = Difficulty{
		Name:           "Medium",
		StartingCash:   1500000, // $1.5M fund
		EventFrequency: 0.30,    // 30% chance
		Volatility:     0.05,    // 5% volatility
		MaxTurns:       60,      // 5 years
		Description:    "$1.5M fund - balanced challenge, 5 years",
	}

	HardDifficulty = Difficulty{
		Name:           "Hard",
		StartingCash:   2000000, // $2M fund
		EventFrequency: 0.40,    // 40% chance
		Volatility:     0.07,    // 7% volatility
		MaxTurns:       60,      // 5 years
		Description:    "$2M fund, higher volatility, 5 years",
	}

	ExpertDifficulty = Difficulty{
		Name:           "Expert",
		StartingCash:   2500000, // $2.5M fund
		EventFrequency: 0.50,    // 50% chance
		Volatility:     0.10,    // 10% volatility
		MaxTurns:       60,      // 5 years
		Description:    "$2.5M fund, extreme volatility, 5 years",
	}
)

// NewGame initializes a new game with specified difficulty
func NewGame(playerName string, difficulty Difficulty) *GameState {
	rand.Seed(time.Now().UnixNano())

	// Calculate follow-on reserve: $100k base + $50k per potential funding round
	// Assume ~60% of companies will have at least one round we can participate in
	expectedRounds := int64(15 * 0.6 * 2) // 15 companies, 60% raise, avg 2 rounds
	followOnReserve := int64(100000) + (expectedRounds * 50000)

	gs := &GameState{
		PlayerName: playerName,
		Difficulty: difficulty,
		Portfolio: Portfolio{
			Cash:                difficulty.StartingCash,
			NetWorth:            difficulty.StartingCash,
			Turn:                1,
			MaxTurns:            difficulty.MaxTurns,
			InitialFundSize:     difficulty.StartingCash,
			AnnualManagementFee: 0.02,            // 2% annual management fee
			FollowOnReserve:     followOnReserve, // Dynamic based on expected rounds
		},
	}

	gs.LoadStartups()
	gs.LoadEvents()
	gs.InitializeAIPlayers()
	gs.ScheduleFundingRounds()
	gs.ScheduleAcquisitions()
	gs.ScheduleDramaticEvents()

	return gs
}

// LoadStartups loads 15 randomly selected startup companies from 30 available JSON files
func (gs *GameState) LoadStartups() {
	gs.AvailableStartups = []Startup{}
	allStartups := []Startup{}

	// Load all 30 startups
	for i := 1; i <= 30; i++ {
		var startup Startup
		jsonFile, err := os.Open(fmt.Sprintf("startups/%d.json", i))
		if err != nil {
			fmt.Printf("Warning: Could not load startup %d: %v\n", i, err)
			continue
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)
		jsonFile.Close()

		json.Unmarshal(byteValue, &startup)

		// Cap all initial valuations at $1M or less (pre-seed stage)
		// Generate realistic pre-seed valuations between $250k - $1M
		startup.Valuation = int64(250000 + rand.Intn(750000))

		// Calculate risk and growth scores based on metrics
		startup.RiskScore = gs.calculateRiskScore(&startup)
		startup.GrowthPotential = gs.calculateGrowthPotential(&startup)

		// Initialize financial metrics
		startup.MonthlyRevenue = int64(startup.MonthlySales * startup.SalePrice)
		startup.MonthlyCosts = int64(startup.GrossBurnRate * 1000) // Convert to actual dollars
		startup.NetIncome = startup.MonthlyRevenue - startup.MonthlyCosts
		startup.CustomerCount = startup.MonthlySales // Approximate
		startup.MonthlyRecurringRevenue = startup.MonthlyRevenue
		startup.RevenueGrowthRate = 0.05 // Default 5% growth
		startup.Last409AValuation = startup.Valuation
		startup.Last409AMonth = 0

		allStartups = append(allStartups, startup)
	}

	// Randomly select 15 from the 30 startups
	if len(allStartups) > 15 {
		// Shuffle and take first 15
		rand.Shuffle(len(allStartups), func(i, j int) {
			allStartups[i], allStartups[j] = allStartups[j], allStartups[i]
		})
		gs.AvailableStartups = allStartups[:15]
	} else {
		gs.AvailableStartups = allStartups
	}
}

// LoadEvents loads all possible game events
func (gs *GameState) LoadEvents() {
	gs.EventPool = []GameEvent{}

	jsonFile, err := os.Open("rounds/round-options.json")
	if err != nil {
		fmt.Printf("Warning: Could not load events: %v\n", err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var events [][]GameEvent
	json.Unmarshal(byteValue, &events)

	if len(events) > 0 {
		gs.EventPool = events[0]
	}
}

// GenerateTermOptions generates 3 term sheet options for an investment
func (gs *GameState) GenerateTermOptions(startup *Startup, amount int64) []InvestmentTerms {
	options := []InvestmentTerms{}

	// Option 1: Preferred Stock (VC standard)
	options = append(options, InvestmentTerms{
		Type:               "Preferred Stock",
		HasProRataRights:   true,
		HasInfoRights:      true,
		HasBoardSeat:       amount >= 100000, // Board seat for $100k+ investments
		LiquidationPref:    1.0,
		HasAntiDilution:    true,
		ConversionDiscount: 0.0,
	})

	// Option 2: SAFE (Simple Agreement for Future Equity)
	options = append(options, InvestmentTerms{
		Type:               "SAFE",
		HasProRataRights:   true,
		HasInfoRights:      false,
		HasBoardSeat:       false,
		LiquidationPref:    0.0, // No liquidation preference with SAFE
		HasAntiDilution:    false,
		ConversionDiscount: 0.20, // 20% discount on conversion
	})

	// Option 3: Common Stock (founder-friendly)
	options = append(options, InvestmentTerms{
		Type:               "Common Stock",
		HasProRataRights:   false,
		HasInfoRights:      false,
		HasBoardSeat:       false,
		LiquidationPref:    0.0,
		HasAntiDilution:    false,
		ConversionDiscount: 0.0,
	})

	return options
}

// MakeInvestment allows player to invest in a startup
func (gs *GameState) MakeInvestment(startupIndex int, amount int64) error {
	return gs.MakeInvestmentWithTerms(startupIndex, amount, InvestmentTerms{
		Type:             "Preferred Stock",
		HasProRataRights: true,
		HasInfoRights:    true,
		HasBoardSeat:     amount >= 100000,
		LiquidationPref:  1.0,
		HasAntiDilution:  true,
	})
}

// MakeInvestmentWithTerms allows player to invest with specific terms
func (gs *GameState) MakeInvestmentWithTerms(startupIndex int, amount int64, terms InvestmentTerms) error {
	if amount <= 0 {
		return fmt.Errorf("investment amount must be positive")
	}

	if amount > gs.Portfolio.Cash {
		return fmt.Errorf("insufficient funds (have $%d, need $%d)", gs.Portfolio.Cash, amount)
	}

	if startupIndex < 0 || startupIndex >= len(gs.AvailableStartups) {
		return fmt.Errorf("invalid startup index")
	}

	startup := gs.AvailableStartups[startupIndex]

	// Check if already invested in this company
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == startup.Name {
			return fmt.Errorf("you have already invested in %s", startup.Name)
		}
	}

	// Minimum investment is $10,000 (standard VC practice)
	minInvestment := int64(10000)
	if amount < minInvestment {
		return fmt.Errorf("minimum investment is $%d", minInvestment)
	}

	// Maximum investment is 20% of company valuation (standard VC practice)
	maxInvestment := int64(float64(startup.Valuation) * 0.20)
	if amount > maxInvestment {
		return fmt.Errorf("maximum investment is $%d (20%% of company valuation: $%d)", maxInvestment, startup.Valuation)
	}

	// Calculate equity percentage based on investment amount and company valuation
	// Only 20% of company is available for investment in this round
	equityPercent := (float64(amount) / float64(startup.Valuation)) * 100.0

	// Apply SAFE discount if applicable
	if terms.Type == "SAFE" && terms.ConversionDiscount > 0 {
		equityPercent = equityPercent * (1 + terms.ConversionDiscount)
		// Cap at 20% even with discount (since only 20% is available)
		maxEquityPercent := 20.0 * (1 + terms.ConversionDiscount)
		if equityPercent > maxEquityPercent {
			equityPercent = maxEquityPercent
		}
	}

	// Safety cap: equity cannot exceed 20% (or 24% with SAFE discount)
	maxEquityPercent := 20.0
	if terms.Type == "SAFE" && terms.ConversionDiscount > 0 {
		maxEquityPercent = 20.0 * (1 + terms.ConversionDiscount)
	}
	if equityPercent > maxEquityPercent {
		equityPercent = maxEquityPercent
	}

	investment := Investment{
		CompanyName:      startup.Name,
		AmountInvested:   amount,
		EquityPercent:    equityPercent,
		InitialEquity:    equityPercent,
		InitialValuation: startup.Valuation,
		CurrentValuation: startup.Valuation,
		MonthsHeld:       0,
		Category:         startup.Category,
		NegativeNewsSent: false,
		Rounds:           []FundingRound{},
		Terms:            terms,
		FollowOnThisTurn: false,
	}

	gs.Portfolio.Investments = append(gs.Portfolio.Investments, investment)
	gs.Portfolio.Cash -= amount
	gs.updateNetWorth()

	return nil
}

// GetFollowOnOpportunities returns any follow-on investment opportunities for this turn
func (gs *GameState) GetFollowOnOpportunities() []FollowOnOpportunity {
	opportunities := []FollowOnOpportunity{}

	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Check if player has invested in this company
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == event.CompanyName {
					// Find the startup
					for _, startup := range gs.AvailableStartups {
						if startup.Name == event.CompanyName {
							preMoneyVal := startup.Valuation
							postMoneyVal := preMoneyVal + event.RaiseAmount

							// Calculate min/max investment amounts
							minInvestment := int64(10000) // $10k minimum
							// Maximum investment is 20% of pre-money valuation (standard VC practice)
							maxInvestmentByValuation := int64(float64(preMoneyVal) * 0.20)
							// Use available cash (uninvested money from beginning) + follow-on reserve
							availableCash := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve
							// Maximum is the lower of: 20% of valuation, available cash, or 50% of raise amount
							maxInvestment := maxInvestmentByValuation
							if maxInvestment > availableCash {
								maxInvestment = availableCash
							}
							if maxInvestment > event.RaiseAmount/2 {
								maxInvestment = event.RaiseAmount / 2 // Can't invest more than half the round
							}

							opportunities = append(opportunities, FollowOnOpportunity{
								CompanyName:   event.CompanyName,
								RoundName:     event.RoundName,
								PreMoneyVal:   preMoneyVal,
								PostMoneyVal:  postMoneyVal,
								CurrentEquity: inv.EquityPercent,
								MinInvestment: minInvestment,
								MaxInvestment: maxInvestment,
							})
							break
						}
					}
					break
				}
			}
		}
	}

	return opportunities
}

// MakeFollowOnInvestment allows investing more in a company during a funding round
func (gs *GameState) MakeFollowOnInvestment(companyName string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("investment amount must be positive")
	}

	if amount > gs.Portfolio.Cash+gs.Portfolio.FollowOnReserve {
		return fmt.Errorf("insufficient follow-on funds (have $%d, need $%d)", gs.Portfolio.Cash+gs.Portfolio.FollowOnReserve, amount)
	}

	// Use cash first, then follow-on reserve
	drawnFromCash := amount
	if drawnFromCash > gs.Portfolio.Cash {
		drawnFromCash = gs.Portfolio.Cash
	}
	drawnFromReserve := amount - drawnFromCash

	// Find the funding round event for this turn to get the post-money valuation
	var postMoneyVal int64
	foundRound := false

	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn && event.CompanyName == companyName {
			// Find the company
			for _, startup := range gs.AvailableStartups {
				if startup.Name == companyName {
					preMoneyVal := startup.Valuation
					postMoneyVal = preMoneyVal + event.RaiseAmount
					foundRound = true
					break
				}
			}
			break
		}
	}

	if !foundRound {
		return fmt.Errorf("no funding round happening for %s this turn", companyName)
	}

	// Find the company and calculate max investment (20% of pre-money valuation)
	var preMoneyVal int64
	var foundCompany bool
	for _, startup := range gs.AvailableStartups {
		if startup.Name == companyName {
			preMoneyVal = startup.Valuation
			foundCompany = true
			break
		}
	}

	if !foundCompany {
		return fmt.Errorf("company %s not found", companyName)
	}

	// Maximum follow-on investment is 20% of pre-money valuation
	maxInvestment := int64(float64(preMoneyVal) * 0.20)

	// Find the investment
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			inv := &gs.Portfolio.Investments[i]

			// Check if this follow-on investment would exceed 20% limit
			// We need to check total investment (existing + new) against 20% of valuation
			totalInvestmentAfterFollowOn := inv.AmountInvested + amount
			if totalInvestmentAfterFollowOn > maxInvestment {
				return fmt.Errorf("total investment would exceed maximum of $%d (20%% of company valuation: $%d)", maxInvestment, preMoneyVal)
			}

			// Update total amount invested
			inv.AmountInvested += amount

			// Calculate effective investment amount (applying SAFE conversion discount if applicable)
			effectiveAmount := float64(inv.AmountInvested)
			if inv.Terms.Type == "SAFE" && inv.Terms.ConversionDiscount > 0 {
				// SAFE converts at a discount - your investment gets more equity
				effectiveAmount = float64(inv.AmountInvested) * (1.0 + inv.Terms.ConversionDiscount)
			}

			// Recalculate total equity based on total invested amount and post-money valuation
			// This ensures equity is always calculated correctly relative to current valuation
			newEquityPercent := (effectiveAmount / float64(postMoneyVal)) * 100.0

			// Cap equity at 100% (should never happen, but safety check)
			if newEquityPercent > 100.0 {
				newEquityPercent = 100.0
			}

			inv.EquityPercent = newEquityPercent
			inv.FollowOnThisTurn = true // Mark that follow-on was made this turn

			// Deduct from cash first, then follow-on reserve
			gs.Portfolio.Cash -= drawnFromCash
			gs.Portfolio.FollowOnReserve -= drawnFromReserve
			gs.updateNetWorth()

			return nil
		}
	}

	return fmt.Errorf("you have not invested in %s", companyName)
}

// HasFollowOnOpportunities checks if there are any follow-on opportunities this turn
func (gs *GameState) HasFollowOnOpportunities() bool {
	opportunities := gs.GetFollowOnOpportunities()
	return len(opportunities) > 0
}

// UpdateCompanyFinancials updates monthly financials for a company
func (gs *GameState) UpdateCompanyFinancials(startup *Startup) {
	// Apply growth rate to revenue (with some randomness)
	growthVariance := (rand.Float64()*0.4 - 0.2) // -20% to +20% variance
	actualGrowth := startup.RevenueGrowthRate + growthVariance

	// Update revenue based on growth
	startup.MonthlyRevenue = int64(float64(startup.MonthlyRevenue) * (1 + actualGrowth))

	// Costs grow slower than revenue (economies of scale)
	costGrowth := actualGrowth * 0.6 // Costs grow at 60% of revenue growth rate
	startup.MonthlyCosts = int64(float64(startup.MonthlyCosts) * (1 + costGrowth))

	// Calculate net income
	startup.NetIncome = startup.MonthlyRevenue - startup.MonthlyCosts

	// Update cumulative totals
	startup.CumulativeRevenue += startup.MonthlyRevenue
	startup.CumulativeCosts += startup.MonthlyCosts

	// Update customer count based on revenue
	if startup.SalePrice > 0 {
		startup.CustomerCount = int(startup.MonthlyRevenue / int64(startup.SalePrice))
	}

	// Update MRR
	startup.MonthlyRecurringRevenue = startup.MonthlyRevenue

	// Adjust growth rate based on performance
	if startup.NetIncome > 0 {
		startup.RevenueGrowthRate *= 1.02 // Profitable companies grow faster
	} else {
		startup.RevenueGrowthRate *= 0.98 // Unprofitable slow down
	}

	// Cap growth rate
	if startup.RevenueGrowthRate > 0.30 {
		startup.RevenueGrowthRate = 0.30 // Max 30% monthly growth
	}
	if startup.RevenueGrowthRate < -0.15 {
		startup.RevenueGrowthRate = -0.15 // Max 15% monthly decline
	}

	// Update valuation based on financial performance
	annualRevenue := startup.MonthlyRevenue * 12

	// Revenue multiple varies by profitability
	revenueMultiple := 10.0
	if startup.NetIncome > 0 {
		revenueMultiple = 15.0 // Profitable companies get premium
	}

	newValuation := int64(float64(annualRevenue) * revenueMultiple)

	// Smooth valuation changes (max 20% per month)
	maxChange := float64(startup.Valuation) * 0.20
	valuationChange := newValuation - startup.Valuation
	if valuationChange > int64(maxChange) {
		newValuation = startup.Valuation + int64(maxChange)
	} else if valuationChange < -int64(maxChange) {
		newValuation = startup.Valuation - int64(maxChange)
	}

	// Minimum valuation
	if newValuation < 100000 {
		newValuation = 100000
	}

	startup.Valuation = newValuation
}

// Calculate409AValuation performs quarterly 409A valuation
func (gs *GameState) Calculate409AValuation(startup *Startup) int64 {
	// 409A considers multiple factors
	annualRevenue := startup.MonthlyRevenue * 12

	// Revenue multiple (conservative for 409A)
	revenueMultiple := 8.0
	if startup.NetIncome > 0 {
		revenueMultiple = 12.0
	}
	revenueValue := int64(float64(annualRevenue) * revenueMultiple)

	// Cost to duplicate
	costValue := startup.CumulativeCosts

	// Market value
	marketValue := startup.Valuation

	// Weighted average
	val409A := (revenueValue*4 + costValue*2 + marketValue*4) / 10

	// 409A is typically 20-30% discount to FMV
	val409A = int64(float64(val409A) * 0.75)

	startup.Last409AValuation = val409A
	startup.Last409AMonth = gs.Portfolio.Turn

	return val409A
}

// formatCurrency formats a number as currency
func formatCurrency(amount int64) string {
	abs := amount
	if abs < 0 {
		abs = -abs
	}

	s := fmt.Sprintf("%d", abs)
	n := len(s)
	if n <= 3 {
		if amount < 0 {
			return "-" + s
		}
		return s
	}

	result := ""
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	if amount < 0 {
		return "-" + result
	}
	return result
}

// ProcessTurn simulates one month of game time
func (gs *GameState) ProcessTurn() []string {
	messages := []string{}

	// Process management fees
	feeMessages := gs.ProcessManagementFees()
	messages = append(messages, feeMessages...)

	// NOTE: Follow-on investments should be handled BEFORE this function is called
	// Process funding rounds
	roundMessages := gs.ProcessFundingRounds()
	messages = append(messages, roundMessages...)

	// Process dramatic events (scandals, co-founder splits, etc.)
	dramaMessages := gs.ProcessDramaticEvents()
	messages = append(messages, dramaMessages...)

	// Process acquisitions
	acqMessages := gs.ProcessAcquisitions()
	messages = append(messages, acqMessages...)

	// Old random event code removed - now using financial-based valuation below

	// Update financials for all companies
	for i := range gs.AvailableStartups {
		startup := &gs.AvailableStartups[i]
		gs.UpdateCompanyFinancials(startup)

		// Do 409A valuation quarterly (every 3 months, starting at month 4)
		// First 409A should be at month 4 (not month 3), then every 3 months: 4, 7, 10, 13...
		if gs.Portfolio.Turn >= 4 && (gs.Portfolio.Turn-1)%3 == 0 {
			val409A := gs.Calculate409AValuation(startup)

			// Show 409A for companies we're invested in
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == startup.Name {
					profitLossStr := ""
					if startup.NetIncome >= 0 {
						profitLossStr = fmt.Sprintf("Profit: $%s", formatCurrency(startup.NetIncome))
					} else {
						profitLossStr = fmt.Sprintf("Loss: $%s", formatCurrency(-startup.NetIncome))
					}

					messages = append(messages, fmt.Sprintf(
						"?? %s 409A: $%s (FMV: $%s, Revenue: $%s/mo, %s)",
						startup.Name,
						formatCurrency(val409A),
						formatCurrency(startup.Valuation),
						formatCurrency(startup.MonthlyRevenue),
						profitLossStr,
					))
					break
				}
			}
		}
	}

	// Update player investments based on company valuations
	for i := range gs.Portfolio.Investments {
		inv := &gs.Portfolio.Investments[i]

		wasAboveInitial := inv.CurrentValuation >= inv.InitialValuation

		// Find the company and update valuation
		for _, startup := range gs.AvailableStartups {
			if startup.Name == inv.CompanyName {
				oldVal := inv.CurrentValuation
				inv.CurrentValuation = startup.Valuation

				// Show significant monthly changes (>15%)
				change := inv.CurrentValuation - oldVal
				if oldVal > 0 {
					percentChange := float64(change) / float64(oldVal) * 100.0

					if percentChange > 15.0 {
						messages = append(messages, fmt.Sprintf(
							"?? %s: Strong growth! Revenue $%s/mo (+%.1f%%)",
							startup.Name,
							formatCurrency(startup.MonthlyRevenue),
							percentChange,
						))
					} else if percentChange < -15.0 {
						messages = append(messages, fmt.Sprintf(
							"?? %s: Declining. Revenue $%s/mo (%.1f%%)",
							startup.Name,
							formatCurrency(startup.MonthlyRevenue),
							percentChange,
						))
					}
				}
				break
			}
		}

		// Check if investment just went negative and generate news
		nowBelowInitial := inv.CurrentValuation < inv.InitialValuation
		if wasAboveInitial && nowBelowInitial && !inv.NegativeNewsSent {
			inv.NegativeNewsSent = true
			news := gs.generateNegativeNews(inv)
			messages = append(messages, news)
		}
	}

	gs.Portfolio.Turn++
	gs.updateNetWorth()

	// Process AI player turns
	gs.ProcessAITurns()

	return messages
}

// updateNetWorth calculates current net worth
func (gs *GameState) updateNetWorth() {
	netWorth := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve

	for _, inv := range gs.Portfolio.Investments {
		// Value of investment = (equity % / 100) * current valuation
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		netWorth += value
	}

	gs.Portfolio.NetWorth = netWorth
}

// GetPortfolioValue returns the current value of all investments
func (gs *GameState) GetPortfolioValue() int64 {
	total := int64(0)
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		total += value
	}
	return total
}

// IsGameOver checks if game has ended
func (gs *GameState) IsGameOver() bool {
	return gs.Portfolio.Turn > gs.Portfolio.MaxTurns
}

// GetFinalScore calculates the final score
func (gs *GameState) GetFinalScore() (netWorth int64, roi float64, successfulExits int) {
	netWorth = gs.Portfolio.NetWorth

	// Calculate ROI based on TOTAL starting capital (cash + follow-on reserve)
	totalStartingCapital := gs.Portfolio.InitialFundSize + gs.Portfolio.FollowOnReserve
	roi = ((float64(netWorth) - float64(totalStartingCapital)) / float64(totalStartingCapital)) * 100.0

	// Count successful exits (investments that 5x'd or more)
	successfulExits = 0
	for _, inv := range gs.Portfolio.Investments {
		if inv.CurrentValuation >= inv.InitialValuation*5 {
			successfulExits++
		}
	}

	return netWorth, roi, successfulExits
}

// InitializeAIPlayers creates computer VC opponents
func (gs *GameState) InitializeAIPlayers() {
	gs.AIPlayers = []AIPlayer{
		{
			Name:          "CARL",
			Firm:          "Sterling & Cooper",
			Strategy:      "conservative",
			RiskTolerance: 0.3,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
			},
		},
		{
			Name:          "Sarah Chen",
			Firm:          "Accel Partners",
			Strategy:      "aggressive",
			RiskTolerance: 0.8,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
			},
		},
		{
			Name:          "Marcus Williams",
			Firm:          "Sequoia Capital",
			Strategy:      "balanced",
			RiskTolerance: 0.5,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     gs.Portfolio.FollowOnReserve,
			},
		},
	}
}

// ScheduleFundingRounds schedules future funding rounds for companies
func (gs *GameState) ScheduleFundingRounds() {
	gs.FundingRoundQueue = []FundingRoundEvent{}

	// Schedule funding rounds with realistic amounts
	for _, startup := range gs.AvailableStartups {
		// Seed round (3-9 months) - raise $2M-$5M
		seedTurn := 3 + rand.Intn(7)
		if seedTurn < gs.Portfolio.MaxTurns {
			seedAmount := int64(2000000 + rand.Intn(3000000)) // $2M-$5M
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Seed",
				ScheduledTurn: seedTurn,
				RaiseAmount:   seedAmount,
			})
		}

		// Series A (12-24 months) - raise $10M-$20M
		seriesATurn := 12 + rand.Intn(13)
		if seriesATurn < gs.Portfolio.MaxTurns {
			seriesAAmount := int64(10000000 + rand.Intn(10000000)) // $10M-$20M
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Series A",
				ScheduledTurn: seriesATurn,
				RaiseAmount:   seriesAAmount,
			})
		}

		// Series B (30-48 months) - raise $30M-$50M
		seriesBTurn := 30 + rand.Intn(19)
		if seriesBTurn < gs.Portfolio.MaxTurns {
			seriesBAmount := int64(30000000 + rand.Intn(20000000)) // $30M-$50M
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Series B",
				ScheduledTurn: seriesBTurn,
				RaiseAmount:   seriesBAmount,
			})
		}

		// Series C (48-60 months) - raise $50M-$100M, only for top performers
		if rand.Float64() < 0.3 { // 30% of companies
			seriesCTurn := 48 + rand.Intn(13)
			if seriesCTurn < gs.Portfolio.MaxTurns {
				seriesCAmount := int64(50000000 + rand.Intn(50000000)) // $50M-$100M
				gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
					CompanyName:   startup.Name,
					RoundName:     "Series C",
					ScheduledTurn: seriesCTurn,
					RaiseAmount:   seriesCAmount,
				})
			}
		}

		// 20% chance of a down round occurring (usually Series A or B)
		if rand.Float64() < 0.2 {
			downRoundTurn := 20 + rand.Intn(30) // Months 20-50
			if downRoundTurn < gs.Portfolio.MaxTurns {
				downRoundName := "Series A (Down)"
				if rand.Float64() < 0.5 {
					downRoundName = "Series B (Down)"
				}
				downAmount := int64(5000000 + rand.Intn(15000000)) // $5M-$20M
				gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
					CompanyName:   startup.Name,
					RoundName:     downRoundName,
					ScheduledTurn: downRoundTurn,
					RaiseAmount:   downAmount,
					IsDownRound:   true,
				})
			}
		}
	}
}

// ScheduleAcquisitions schedules potential acquisition offers
func (gs *GameState) ScheduleAcquisitions() {
	gs.AcquisitionQueue = []AcquisitionEvent{}

	// 40% of companies get acquisition offers
	for _, startup := range gs.AvailableStartups {
		if rand.Float64() < 0.4 {
			// Acquisitions happen between months 24-60
			acqTurn := 24 + rand.Intn(37)
			if acqTurn < gs.Portfolio.MaxTurns {
				// Multiple ranges from 3x to 6x EBITDA (4x average)
				multiple := 3.0 + rand.Float64()*3.0

				// Due diligence quality
				dueDiligence := "normal"
				roll := rand.Float64()
				if roll < 0.15 { // 15% bad due diligence
					dueDiligence = "bad"
					multiple *= 0.6 // Offer falls through or gets cut 40%
				} else if roll > 0.85 { // 15% great due diligence
					dueDiligence = "good"
					multiple *= 1.2 // Offer increases 20%
				}

				gs.AcquisitionQueue = append(gs.AcquisitionQueue, AcquisitionEvent{
					CompanyName:   startup.Name,
					ScheduledTurn: acqTurn,
					OfferMultiple: multiple,
					DueDiligence:  dueDiligence,
				})
			}
		}
	}
}

// ScheduleDramaticEvents schedules scandals, co-founder splits, etc.
func (gs *GameState) ScheduleDramaticEvents() {
	gs.DramaticEventQueue = []DramaticEvent{}

	// Event frequency based on difficulty
	// Easy: 10%, Medium: 20%, Hard: 30%, Expert: 40%
	eventChance := 0.10
	if gs.Difficulty.Name == "Medium" {
		eventChance = 0.20
	} else if gs.Difficulty.Name == "Hard" {
		eventChance = 0.30
	} else if gs.Difficulty.Name == "Expert" {
		eventChance = 0.40
	}

	eventTypes := []string{
		"cofounder_split", "scandal", "lawsuit", "pivot_fail",
		"fraud", "data_breach", "key_hire_quit", "regulatory_issue",
		"competitor_attack", "product_failure",
	}

	for _, startup := range gs.AvailableStartups {
		if rand.Float64() < eventChance {
			// Events happen between months 6-55
			eventTurn := 6 + rand.Intn(50)
			if eventTurn < gs.Portfolio.MaxTurns {
				eventType := eventTypes[rand.Intn(len(eventTypes))]

				// Determine severity (difficulty affects this)
				severityRoll := rand.Float64()
				severity := "minor"
				impactPercent := 0.85 // 15% drop

				if gs.Difficulty.Name == "Hard" || gs.Difficulty.Name == "Expert" {
					// Harder difficulties have worse outcomes
					if severityRoll < 0.25 {
						severity = "severe"
						impactPercent = 0.40 // 60% drop
					} else if severityRoll < 0.55 {
						severity = "moderate"
						impactPercent = 0.65 // 35% drop
					}
				} else if gs.Difficulty.Name == "Medium" {
					if severityRoll < 0.15 {
						severity = "severe"
						impactPercent = 0.50 // 50% drop
					} else if severityRoll < 0.40 {
						severity = "moderate"
						impactPercent = 0.70 // 30% drop
					}
				} else {
					// Easy mode
					if severityRoll < 0.10 {
						severity = "moderate"
						impactPercent = 0.75 // 25% drop
					}
				}

				gs.DramaticEventQueue = append(gs.DramaticEventQueue, DramaticEvent{
					CompanyName:   startup.Name,
					ScheduledTurn: eventTurn,
					EventType:     eventType,
					Severity:      severity,
					ImpactPercent: impactPercent,
				})
			}
		}
	}
}

// ProcessDramaticEvents handles scandals, co-founder splits, etc.
func (gs *GameState) ProcessDramaticEvents() []string {
	messages := []string{}

	for _, event := range gs.DramaticEventQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Find the company
			for i := range gs.AvailableStartups {
				if gs.AvailableStartups[i].Name == event.CompanyName {
					startup := &gs.AvailableStartups[i]

					oldValuation := startup.Valuation
					startup.Valuation = int64(float64(startup.Valuation) * event.ImpactPercent)

					// Generate message based on event type
					var eventMsg string
					var emoji string

					switch event.EventType {
					case "cofounder_split":
						emoji = "💔"
						if event.Severity == "severe" {
							eventMsg = "Co-founders had MAJOR falling out! CEO resigned. Board in chaos."
						} else if event.Severity == "moderate" {
							eventMsg = "Co-founder conflict! One founder left with equity dispute."
						} else {
							eventMsg = "Minor co-founder disagreement resolved, but caused delays."
						}
					case "scandal":
						emoji = "🔥"
						if event.Severity == "severe" {
							eventMsg = "MAJOR SCANDAL! CEO involved in workplace harassment allegations."
						} else if event.Severity == "moderate" {
							eventMsg = "PR scandal! Questionable business practices exposed."
						} else {
							eventMsg = "Minor controversy in the press, manageable."
						}
					case "lawsuit":
						emoji = "⚖️"
						if event.Severity == "severe" {
							eventMsg = "Class-action lawsuit filed! Facing $50M+ in liabilities."
						} else if event.Severity == "moderate" {
							eventMsg = "Patent infringement lawsuit. Legal costs mounting."
						} else {
							eventMsg = "Small legal dispute, expected to settle."
						}
					case "fraud":
						emoji = "🚨"
						if event.Severity == "severe" {
							eventMsg = "FRAUD DISCOVERED! CFO cooking books. SEC investigation."
						} else {
							eventMsg = "Financial irregularities found. Auditors called in."
						}
					case "data_breach":
						emoji = "🔓"
						if event.Severity == "severe" {
							eventMsg = "MASSIVE DATA BREACH! Customer data leaked. GDPR fines incoming."
						} else if event.Severity == "moderate" {
							eventMsg = "Security breach! Customer trust damaged."
						} else {
							eventMsg = "Minor security incident, quickly patched."
						}
					case "key_hire_quit":
						emoji = "👋"
						if event.Severity == "severe" {
							eventMsg = "CTO quit and joined competitor! Taking team with them."
						} else if event.Severity == "moderate" {
							eventMsg = "VP Engineering resigned. Product roadmap delayed."
						} else {
							eventMsg = "Senior engineer left. Minor setback."
						}
					case "regulatory_issue":
						emoji = "📋"
						if event.Severity == "severe" {
							eventMsg = "Regulatory crackdown! Business model under threat."
						} else {
							eventMsg = "New compliance requirements. Extra costs."
						}
					case "pivot_fail":
						emoji = "🔄"
						if event.Severity == "severe" {
							eventMsg = "Pivot FAILED! Lost key customers and burning cash fast."
						} else {
							eventMsg = "Pivot struggling. Market not responding well."
						}
					case "competitor_attack":
						emoji = "⚔️"
						if event.Severity == "severe" {
							eventMsg = "Competitor launched predatory pricing! Market share plummeting."
						} else {
							eventMsg = "New competitor with better product. Losing customers."
						}
					case "product_failure":
						emoji = "💥"
						if event.Severity == "severe" {
							eventMsg = "Major product launch FLOPPED! Customers demanding refunds."
						} else {
							eventMsg = "Product update buggy. Customer complaints rising."
						}
					default:
						emoji = "⚠️"
						eventMsg = "Unexpected crisis hit the company."
					}

					// Check if player invested
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]
							inv.CurrentValuation = startup.Valuation

							valuationDrop := oldValuation - startup.Valuation
							dropPercent := float64(valuationDrop) / float64(oldValuation) * 100

							messages = append(messages, fmt.Sprintf(
								"%s %s: %s (Valuation: $%s → $%s, -%.0f%%)",
								emoji,
								event.CompanyName,
								eventMsg,
								formatCurrency(oldValuation),
								formatCurrency(startup.Valuation),
								dropPercent,
							))
							break
						}
					}

					// Update AI investments
					for k := range gs.AIPlayers {
						for j := range gs.AIPlayers[k].Portfolio.Investments {
							if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
								gs.AIPlayers[k].Portfolio.Investments[j].CurrentValuation = startup.Valuation
							}
						}
					}
				}
			}
		}
	}

	return messages
}

// ProcessFundingRounds handles any scheduled funding rounds this turn
func (gs *GameState) ProcessFundingRounds() []string {
	messages := []string{}

	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Find the company
			for i := range gs.AvailableStartups {
				if gs.AvailableStartups[i].Name == event.CompanyName {
					startup := &gs.AvailableStartups[i]

					var preMoneyVal int64
					var postMoneyVal int64
					var dilutionFactor float64

					if event.IsDownRound {
						// Down round: pre-money is 60-90% of current valuation
						downFactor := 0.6 + rand.Float64()*0.3 // 60%-90%
						preMoneyVal = int64(float64(startup.Valuation) * downFactor)
						postMoneyVal = preMoneyVal + event.RaiseAmount
						dilutionFactor = float64(preMoneyVal) / float64(postMoneyVal)

						// Check if any investor has board seat - down rounds require board approval
						if gs.HasAnyBoardSeat(event.CompanyName) {
							// Only create vote if player has board seat (player votes, AI votes are simulated)
							if gs.HasBoardSeat(event.CompanyName) {
								// Create board vote for down round
								vote := BoardVote{
									CompanyName:  event.CompanyName,
									VoteType:     "down_round",
									Title:        fmt.Sprintf("Down Round: $%s at $%s pre-money", formatCurrency(event.RaiseAmount), formatCurrency(preMoneyVal)),
									Description:  fmt.Sprintf("%s proposes raising $%s in a DOWN ROUND at $%s pre-money (down from $%s). This will significantly dilute your equity.", event.CompanyName, formatCurrency(event.RaiseAmount), formatCurrency(preMoneyVal), formatCurrency(startup.Valuation)),
									OptionA:      "Approve",
									OptionB:      "Reject",
									ConsequenceA: fmt.Sprintf("Down round approved. Company raises $%s at reduced valuation.", formatCurrency(event.RaiseAmount)),
									ConsequenceB: "Down round rejected. Company must find alternative funding or accept worse terms.",
									RequiresVote: true,
									Turn:         gs.Portfolio.Turn,
									Metadata: map[string]interface{}{
										"raiseAmount":      event.RaiseAmount,
										"preMoneyVal":      preMoneyVal,
										"postMoneyVal":     postMoneyVal,
										"currentValuation": startup.Valuation,
									},
								}
								gs.PendingBoardVotes = append(gs.PendingBoardVotes, vote)
								messages = append(messages, fmt.Sprintf(
									"🏛️  BOARD VOTE REQUIRED: %s proposes a DOWN ROUND. Vote will be required.",
									event.CompanyName,
								))
								continue // Skip processing this round until vote is complete
							}
						}
					} else {
						// Normal round
						preMoneyVal = startup.Valuation
						postMoneyVal = preMoneyVal + event.RaiseAmount
						dilutionFactor = float64(preMoneyVal) / float64(postMoneyVal)
					}

					// Update player's investment if they invested in this company
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]

							// If follow-on investment was made this turn, equity was already recalculated
							// in MakeFollowOnInvestment based on post-money valuation, so we don't dilute again
							if !inv.FollowOnThisTurn {
								// Normal case: dilute existing equity
								oldEquity := inv.EquityPercent
								inv.EquityPercent *= dilutionFactor

								// Record the round
								inv.Rounds = append(inv.Rounds, FundingRound{
									RoundName:        event.RoundName,
									PreMoneyVal:      preMoneyVal,
									InvestmentAmount: event.RaiseAmount,
									PostMoneyVal:     postMoneyVal,
									Month:            gs.Portfolio.Turn,
								})

								if event.IsDownRound {
									messages = append(messages, fmt.Sprintf(
										"⚠️  %s raised $%s in DOWN ROUND (%s)! Valuation dropped. Equity: %.2f%% → %.2f%%",
										event.CompanyName,
										formatCurrency(event.RaiseAmount),
										event.RoundName,
										oldEquity,
										inv.EquityPercent,
									))
								} else {
									messages = append(messages, fmt.Sprintf(
										"🚀 %s raised $%s in %s round! Your equity diluted from %.2f%% to %.2f%%",
										event.CompanyName,
										formatCurrency(event.RaiseAmount),
										event.RoundName,
										oldEquity,
										inv.EquityPercent,
									))
								}
							} else {
								// Follow-on investment case: equity already calculated correctly, just record the round
								oldEquity := inv.EquityPercent
								inv.Rounds = append(inv.Rounds, FundingRound{
									RoundName:        event.RoundName,
									PreMoneyVal:      preMoneyVal,
									InvestmentAmount: event.RaiseAmount,
									PostMoneyVal:     postMoneyVal,
									Month:            gs.Portfolio.Turn,
								})

								// Reset flag for next turn
								inv.FollowOnThisTurn = false

								messages = append(messages, fmt.Sprintf(
									"🚀 %s raised $%s in %s round! Your equity: %.2f%% (includes your follow-on investment)",
									event.CompanyName,
									formatCurrency(event.RaiseAmount),
									event.RoundName,
									oldEquity,
								))
							}
						}
					}

					// Update AI players' investments
					for k := range gs.AIPlayers {
						for j := range gs.AIPlayers[k].Portfolio.Investments {
							if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
								inv := &gs.AIPlayers[k].Portfolio.Investments[j]
								inv.EquityPercent *= dilutionFactor
								inv.Rounds = append(inv.Rounds, FundingRound{
									RoundName:        event.RoundName,
									PreMoneyVal:      preMoneyVal,
									InvestmentAmount: event.RaiseAmount,
									PostMoneyVal:     postMoneyVal,
									Month:            gs.Portfolio.Turn,
								})
							}
						}
					}

					// Update company valuation
					startup.Valuation = postMoneyVal

					// Also update current valuation for all investors
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							gs.Portfolio.Investments[j].CurrentValuation = postMoneyVal
						}
					}
					for k := range gs.AIPlayers {
						for j := range gs.AIPlayers[k].Portfolio.Investments {
							if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
								gs.AIPlayers[k].Portfolio.Investments[j].CurrentValuation = postMoneyVal
							}
						}
					}
				}
			}
		}
	}

	return messages
}

// ProcessAcquisitions handles acquisition offers this turn
func (gs *GameState) ProcessAcquisitions() []string {
	messages := []string{}

	for _, event := range gs.AcquisitionQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Find the company
			for i := range gs.AvailableStartups {
				if gs.AvailableStartups[i].Name == event.CompanyName {
					startup := &gs.AvailableStartups[i]

					// Calculate EBITDA (approximation: annual net income)
					annualEBITDA := startup.NetIncome * 12
					if annualEBITDA < 0 {
						// For unprofitable companies, use revenue multiple instead
						annualEBITDA = startup.MonthlyRevenue * 12
						event.OfferMultiple *= 0.3 // Lower multiple for revenue-based
					}

					// Calculate acquisition offer
					offerValue := int64(float64(annualEBITDA) * event.OfferMultiple)

					// Ensure minimum offer value
					if offerValue < startup.Valuation/2 {
						offerValue = startup.Valuation / 2
					}

					// Check if player invested in this company
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]

							// Calculate payout
							payout := int64((inv.EquityPercent / 100.0) * float64(offerValue))
							returnMultiple := float64(payout) / float64(inv.AmountInvested)

							// If player has board seat, require board vote for acquisitions (unless bad due diligence)
							if inv.Terms.HasBoardSeat && event.DueDiligence != "bad" {
								// Create board vote
								vote := BoardVote{
									CompanyName:  event.CompanyName,
									VoteType:     "acquisition",
									Title:        fmt.Sprintf("Acquisition Offer: $%s", formatCurrency(offerValue)),
									Description:  fmt.Sprintf("Acquirer offers $%s (%.1fx EBITDA) for %s. Your payout would be $%s (%.1fx return).", formatCurrency(offerValue), event.OfferMultiple, event.CompanyName, formatCurrency(payout), returnMultiple),
									OptionA:      "Accept",
									OptionB:      "Reject",
									ConsequenceA: fmt.Sprintf("Acquisition proceeds. You receive $%s.", formatCurrency(payout)),
									ConsequenceB: "Acquisition rejected. Company continues operating independently.",
									RequiresVote: true,
									Turn:         gs.Portfolio.Turn,
									Metadata: map[string]interface{}{
										"offerValue":       offerValue,
										"currentValuation": startup.Valuation,
										"dueDiligence":     event.DueDiligence,
										"offerMultiple":    event.OfferMultiple,
									},
								}
								gs.PendingBoardVotes = append(gs.PendingBoardVotes, vote)
								messages = append(messages, fmt.Sprintf(
									"🏛️  BOARD VOTE REQUIRED: %s received acquisition offer of $%s. Vote will be required.",
									event.CompanyName,
									formatCurrency(offerValue),
								))
								break // Don't execute acquisition yet - wait for vote
							}

							// Add acquisition message based on due diligence
							switch event.DueDiligence {
							case "bad":
								messages = append(messages, fmt.Sprintf(
									"⚠️  %s acquisition FELL THROUGH! Due diligence issues. Offer was $%s (%.1fx EBITDA)",
									event.CompanyName,
									formatCurrency(offerValue),
									event.OfferMultiple,
								))
							case "good":
								messages = append(messages, fmt.Sprintf(
									"🎉 %s ACQUIRED for $%s (%.1fx EBITDA)! Your %.2f%% = $%s (%.1fx return)",
									event.CompanyName,
									formatCurrency(offerValue),
									event.OfferMultiple,
									inv.EquityPercent,
									formatCurrency(payout),
									returnMultiple,
								))
								// Execute acquisition
								gs.Portfolio.Cash += payout
								// Remove investment from portfolio
								gs.Portfolio.Investments = append(gs.Portfolio.Investments[:j], gs.Portfolio.Investments[j+1:]...)
							default: // normal
								messages = append(messages, fmt.Sprintf(
									"💰 %s ACQUIRED for $%s (%.1fx EBITDA)! Your %.2f%% = $%s (%.1fx return)",
									event.CompanyName,
									formatCurrency(offerValue),
									event.OfferMultiple,
									inv.EquityPercent,
									formatCurrency(payout),
									returnMultiple,
								))
								// Execute acquisition
								gs.Portfolio.Cash += payout
								// Remove investment from portfolio
								gs.Portfolio.Investments = append(gs.Portfolio.Investments[:j], gs.Portfolio.Investments[j+1:]...)
							}
							break
						}
					}

					// Handle AI player acquisitions
					if event.DueDiligence != "bad" {
						for k := range gs.AIPlayers {
							for j := len(gs.AIPlayers[k].Portfolio.Investments) - 1; j >= 0; j-- {
								if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
									inv := &gs.AIPlayers[k].Portfolio.Investments[j]
									payout := int64((inv.EquityPercent / 100.0) * float64(offerValue))
									gs.AIPlayers[k].Portfolio.Cash += payout
									// Remove from AI portfolio
									gs.AIPlayers[k].Portfolio.Investments = append(
										gs.AIPlayers[k].Portfolio.Investments[:j],
										gs.AIPlayers[k].Portfolio.Investments[j+1:]...,
									)
									break
								}
							}
						}
					}
				}
			}
		}
	}

	return messages
}

// ProcessManagementFees charges monthly management fees
func (gs *GameState) ProcessManagementFees() []string {
	messages := []string{}

	// Charge management fee monthly (annual rate / 12)
	monthlyFeeRate := gs.Portfolio.AnnualManagementFee / 12.0
	fee := int64(float64(gs.Portfolio.InitialFundSize) * monthlyFeeRate)

	if fee > 0 && gs.Portfolio.Cash >= fee {
		gs.Portfolio.Cash -= fee
		gs.Portfolio.ManagementFeesCharged += fee

		// Also charge AI players
		for i := range gs.AIPlayers {
			aiFee := int64(float64(gs.AIPlayers[i].Portfolio.InitialFundSize) * monthlyFeeRate)
			if gs.AIPlayers[i].Portfolio.Cash >= aiFee {
				gs.AIPlayers[i].Portfolio.Cash -= aiFee
				gs.AIPlayers[i].Portfolio.ManagementFeesCharged += aiFee
			}
		}

		// Only show message every 12 months (annually)
		if gs.Portfolio.Turn%12 == 0 {
			annualFee := fee * 12
			messages = append(messages, fmt.Sprintf(
				"?? Annual management fee charged: $%d (2%% of fund size)",
				annualFee,
			))
		}
	}

	return messages
}

// AIPlayerMakeInvestments has AI players make investment decisions
func (gs *GameState) AIPlayerMakeInvestments() {
	for i := range gs.AIPlayers {
		ai := &gs.AIPlayers[i]

		// Only invest on turn 1 (initial investment phase)
		if gs.Portfolio.Turn != 1 {
			continue
		}

		// AI investment strategy based on risk tolerance
		targetInvestmentCount := 3 + rand.Intn(4) // Invest in 3-6 companies
		availableCash := ai.Portfolio.Cash

		// Shuffle startups for variety
		startups := make([]Startup, len(gs.AvailableStartups))
		copy(startups, gs.AvailableStartups)
		rand.Shuffle(len(startups), func(i, j int) {
			startups[i], startups[j] = startups[j], startups[i]
		})

		investmentsMade := 0
		for _, startup := range startups {
			if investmentsMade >= targetInvestmentCount {
				break
			}

			// Decision based on risk tolerance and startup metrics
			shouldInvest := false
			if ai.Strategy == "conservative" {
				shouldInvest = startup.RiskScore < 0.4 && startup.GrowthPotential > 0.5
			} else if ai.Strategy == "aggressive" {
				shouldInvest = startup.GrowthPotential > 0.7 || (startup.RiskScore > 0.7 && startup.GrowthPotential > 0.6)
			} else { // balanced
				shouldInvest = startup.GrowthPotential > 0.5 && startup.RiskScore < 0.7
			}

			if shouldInvest {
				// Invest portion of available cash
				investmentAmount := availableCash / int64(targetInvestmentCount-investmentsMade)
				if investmentAmount > availableCash {
					investmentAmount = availableCash
				}

				// Maximum investment is 20% of company valuation (standard VC practice)
				maxInvestment := int64(float64(startup.Valuation) * 0.20)
				if investmentAmount > maxInvestment {
					investmentAmount = maxInvestment
				}

				if investmentAmount > 10000 { // Minimum investment
					// Calculate equity percentage (only 20% of company is available)
					equityPercent := (float64(investmentAmount) / float64(startup.Valuation)) * 100.0

					// Safety cap at 20%
					if equityPercent > 20.0 {
						equityPercent = 20.0
					}

					// AI players get Preferred Stock terms (like the player)
					terms := InvestmentTerms{
						Type:               "Preferred Stock",
						HasProRataRights:   true,
						HasInfoRights:      true,
						HasBoardSeat:       investmentAmount >= 100000, // Board seat for $100k+ investments
						LiquidationPref:    1.0,
						HasAntiDilution:    true,
						ConversionDiscount: 0.0,
					}

					investment := Investment{
						CompanyName:      startup.Name,
						AmountInvested:   investmentAmount,
						EquityPercent:    equityPercent,
						InitialEquity:    equityPercent,
						InitialValuation: startup.Valuation,
						CurrentValuation: startup.Valuation,
						MonthsHeld:       0,
						Category:         startup.Category,
						Rounds:           []FundingRound{},
						Terms:            terms,
						FollowOnThisTurn: false,
					}

					ai.Portfolio.Investments = append(ai.Portfolio.Investments, investment)
					ai.Portfolio.Cash -= investmentAmount
					availableCash -= investmentAmount
					investmentsMade++
				}
			}
		}

		// Update AI net worth
		gs.updateAINetWorth(i)
	}
}

// updateAINetWorth calculates AI player net worth
func (gs *GameState) updateAINetWorth(aiIndex int) {
	ai := &gs.AIPlayers[aiIndex]
	netWorth := ai.Portfolio.Cash + ai.Portfolio.FollowOnReserve

	for _, inv := range ai.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		netWorth += value
	}

	ai.Portfolio.NetWorth = netWorth
}

// ProcessAITurns updates AI players' portfolios
func (gs *GameState) ProcessAITurns() {
	// Update all AI investments with same events as player
	for i := range gs.AIPlayers {
		for j := range gs.AIPlayers[i].Portfolio.Investments {
			inv := &gs.AIPlayers[i].Portfolio.Investments[j]
			inv.MonthsHeld++

			wasAboveInitial := inv.CurrentValuation >= inv.InitialValuation

			// Apply same random events and volatility as player investments
			// Random chance of an event happening (based on difficulty)
			if rand.Float64() < gs.Difficulty.EventFrequency && len(gs.EventPool) > 0 {
				event := gs.EventPool[rand.Intn(len(gs.EventPool))]

				inv.CurrentValuation = int64(float64(inv.CurrentValuation) * event.Change)

				// Prevent negative valuations
				if inv.CurrentValuation < 0 {
					inv.CurrentValuation = 0
				}
			} else {
				// Natural growth/decline (random walk) - volatility based on difficulty
				change := (rand.Float64()*2 - 1) * gs.Difficulty.Volatility
				inv.CurrentValuation = int64(float64(inv.CurrentValuation) * (1 + change))
			}

			// Check if investment just went negative (for consistency, but don't generate news for AI)
			nowBelowInitial := inv.CurrentValuation < inv.InitialValuation
			if wasAboveInitial && nowBelowInitial && !inv.NegativeNewsSent {
				inv.NegativeNewsSent = true
			}
		}

		gs.AIPlayers[i].Portfolio.Turn++
		gs.updateAINetWorth(i)
	}
}

// PlayerScore represents a player's score in the leaderboard
type PlayerScore struct {
	Name     string
	Firm     string
	NetWorth int64
	ROI      float64
	IsPlayer bool
}

// GetLeaderboard returns sorted leaderboard of all players
func (gs *GameState) GetLeaderboard() []PlayerScore {
	scores := []PlayerScore{}

	// Add player - ROI based on total starting capital (cash + follow-on reserve)
	totalStartingCapital := gs.Portfolio.InitialFundSize + gs.Portfolio.FollowOnReserve
	playerROI := ((float64(gs.Portfolio.NetWorth) - float64(totalStartingCapital)) / float64(totalStartingCapital)) * 100.0
	scores = append(scores, PlayerScore{
		Name:     gs.PlayerName,
		Firm:     "Your Fund",
		NetWorth: gs.Portfolio.NetWorth,
		ROI:      playerROI,
		IsPlayer: true,
	})

	// Add AI players - same calculation
	for _, ai := range gs.AIPlayers {
		aiTotalCapital := ai.Portfolio.InitialFundSize + ai.Portfolio.FollowOnReserve
		aiROI := ((float64(ai.Portfolio.NetWorth) - float64(aiTotalCapital)) / float64(aiTotalCapital)) * 100.0
		scores = append(scores, PlayerScore{
			Name:     ai.Name,
			Firm:     ai.Firm,
			NetWorth: ai.Portfolio.NetWorth,
			ROI:      aiROI,
			IsPlayer: false,
		})
	}

	// Sort by net worth
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].NetWorth > scores[i].NetWorth {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	return scores
}

// Helper functions
func (gs *GameState) calculateRiskScore(s *Startup) float64 {
	risk := 0.5 // Start at medium risk (minimum is medium)

	// Very high burn rate = VERY HIGH risk
	if s.GrossBurnRate > 40 {
		risk += 0.4
	} else if s.GrossBurnRate > 20 {
		risk += 0.3
	} else if s.GrossBurnRate > 10 {
		risk += 0.2
	} else if s.GrossBurnRate <= 3 {
		// Low burn rate = slightly lower risk (but still medium minimum)
		risk -= 0.1
	}

	// Very low sales = VERY HIGH risk
	if s.MonthlySales < 5 {
		risk += 0.4
	} else if s.MonthlySales < 20 {
		risk += 0.3
	} else if s.MonthlySales < 50 {
		risk += 0.2
	} else if s.MonthlySales > 300 {
		// High sales = slightly lower risk (but still medium minimum)
		risk -= 0.1
	}

	// Ensure 0.5-1.0 range (minimum risk is medium/0.5)
	if risk > 1.0 {
		risk = 1.0
	}
	if risk < 0.5 {
		risk = 0.5
	}

	return risk
}

func (gs *GameState) calculateGrowthPotential(s *Startup) float64 {
	growth := 0.5

	// Very high margins = VERY HIGH growth potential
	if s.PercentMargin > 80 {
		growth += 0.3
	} else if s.PercentMargin > 60 {
		growth += 0.25
	} else if s.PercentMargin > 40 {
		growth += 0.2
	} else if s.PercentMargin > 25 {
		growth += 0.1
	}

	// Very high valuation suggests high growth potential
	if s.Valuation > 100000000 {
		growth += 0.2
	} else if s.Valuation > 50000000 {
		growth += 0.15
	} else if s.Valuation > 30000000 {
		growth += 0.1
	}

	// High activation rate = good growth
	if s.MonthlyActivationRate > 150 {
		growth += 0.15
	} else if s.MonthlyActivationRate > 100 {
		growth += 0.1
	}

	// Ensure 0-1 range
	if growth > 1.0 {
		growth = 1.0
	}
	if growth < 0.0 {
		growth = 0.0
	}

	return growth
}

// generateNegativeNews creates contextual news when a startup goes negative
func (gs *GameState) generateNegativeNews(inv *Investment) string {
	category := inv.Category
	reasons := []string{}

	// Category-specific reasons
	switch category {
	case "FinTech", "Financial":
		reasons = []string{
			"Regulatory scrutiny increased compliance costs",
			"Customer trust eroded after security concerns",
			"Competition from established banks intensified",
			"Regulatory changes impacted revenue model",
			"Market saturation slowed customer acquisition",
		}
	case "BioTech", "HealthTech":
		reasons = []string{
			"Clinical trial delays extended timeline to market",
			"Regulatory approval process took longer than expected",
			"Competitor launched similar product first",
			"Funding challenges slowed R&D progress",
			"Partnership negotiations fell through",
		}
	case "CleanTech", "GreenTech":
		reasons = []string{
			"Policy changes reduced government incentives",
			"Raw material costs increased unexpectedly",
			"Market adoption slower than projected",
			"Infrastructure challenges delayed deployment",
			"Competition from cheaper alternatives increased",
		}
	case "EdTech":
		reasons = []string{
			"School budget cuts reduced institutional sales",
			"Market saturation slowed growth",
			"User retention below expectations",
			"Competition from free platforms increased",
			"Content development costs exceeded projections",
		}
	case "Robotics", "Hardware":
		reasons = []string{
			"Supply chain disruptions delayed production",
			"Manufacturing costs exceeded estimates",
			"Technical hurdles extended development timeline",
			"Market demand weaker than anticipated",
			"Component shortages affected scalability",
		}
	case "Security", "Cybersecurity":
		reasons = []string{
			"High-profile breach damaged reputation",
			"Market crowded with well-funded competitors",
			"Enterprise sales cycles longer than expected",
			"Feature gap compared to established players",
			"Integration challenges slowed adoption",
		}
	case "Gaming", "Entertainment":
		reasons = []string{
			"User acquisition costs exceeded revenue",
			"Retention rates below industry benchmarks",
			"Platform changes affected distribution",
			"Competition from AAA studios intensified",
			"Development costs overran budget",
		}
	case "LegalTech":
		reasons = []string{
			"Law firm adoption slower than projected",
			"Integration complexity deterred clients",
			"Regulatory barriers in some jurisdictions",
			"Competition from established legal software",
			"Customer acquisition costs too high",
		}
	case "AgriTech":
		reasons = []string{
			"Farmer adoption slower than expected",
			"Seasonal factors affected sales cycles",
			"Hardware costs challenged unit economics",
			"Regulatory approvals delayed market entry",
			"Distribution challenges in rural markets",
		}
	case "Logistics", "Supply Chain":
		reasons = []string{
			"Fuel costs reduced profit margins",
			"Market volatility affected demand",
			"Competition from established logistics players",
			"Infrastructure investment required more capital",
			"Regulatory compliance costs increased",
		}
	case "IoT", "Internet of Things":
		reasons = []string{
			"Interoperability standards fragmented market",
			"Security concerns slowed enterprise adoption",
			"Hardware costs challenged scalability",
			"Integration complexity deterred customers",
			"Platform competition intensified",
		}
	case "CloudTech", "SaaS":
		reasons = []string{
			"Customer churn exceeded projections",
			"Market saturation slowed growth",
			"Competition from enterprise giants intensified",
			"Sales cycles longer than expected",
			"Feature development lagged competitors",
		}
	case "Advertising", "Marketing":
		reasons = []string{
			"Ad spend budgets decreased",
			"Platform policy changes affected targeting",
			"Market saturation increased competition",
			"Customer acquisition costs rose",
			"ROI metrics failed to meet expectations",
		}
	default:
		reasons = []string{
			"Market conditions deteriorated",
			"Customer acquisition slower than projected",
			"Competition intensified unexpectedly",
			"Operational costs exceeded revenue",
			"Key partnerships failed to materialize",
		}
	}

	// Select random reason
	reason := reasons[rand.Intn(len(reasons))]

	return fmt.Sprintf("?? %s: Valuation dropped below initial investment. %s", inv.CompanyName, reason)
}

// HasBoardSeat checks if player has a board seat for a given company
func (gs *GameState) HasBoardSeat(companyName string) bool {
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == companyName && inv.Terms.HasBoardSeat {
			return true
		}
	}
	return false
}

// HasAnyBoardSeat checks if any investor (player or AI) has a board seat for a company
func (gs *GameState) HasAnyBoardSeat(companyName string) bool {
	// Check player
	if gs.HasBoardSeat(companyName) {
		return true
	}

	// Check AI players
	for _, ai := range gs.AIPlayers {
		for _, inv := range ai.Portfolio.Investments {
			if inv.CompanyName == companyName && inv.Terms.HasBoardSeat {
				return true
			}
		}
	}

	return false
}

// GetPendingBoardVotes returns all pending board votes for companies where player has a board seat
func (gs *GameState) GetPendingBoardVotes() []BoardVote {
	pending := []BoardVote{}
	for _, vote := range gs.PendingBoardVotes {
		if gs.HasBoardSeat(vote.CompanyName) {
			pending = append(pending, vote)
		}
	}
	return pending
}

// ProcessBoardVote processes a player's vote and determines outcome
func (gs *GameState) ProcessBoardVote(voteIndex int, playerVote string) (string, bool, error) {
	if voteIndex < 0 || voteIndex >= len(gs.PendingBoardVotes) {
		return "", false, fmt.Errorf("invalid vote index")
	}

	vote := &gs.PendingBoardVotes[voteIndex]
	if !gs.HasBoardSeat(vote.CompanyName) {
		return "", false, fmt.Errorf("you do not have a board seat for %s", vote.CompanyName)
	}

	// Normalize vote
	playerVote = strings.ToLower(strings.TrimSpace(playerVote))
	var votedForA bool
	if playerVote == "a" || playerVote == "accept" || playerVote == "approve" || playerVote == "yes" || playerVote == "1" {
		votedForA = true
	} else if playerVote == "b" || playerVote == "reject" || playerVote == "disapprove" || playerVote == "no" || playerVote == "2" {
		votedForA = false
	} else {
		return "", false, fmt.Errorf("invalid vote choice")
	}

	// Simulate board vote: player vote + AI board members vote
	// Player vote counts as 1, AI votes randomly based on their strategy
	playerVotes := 1
	aiVotesA := 0
	aiVotesB := 0

	// Count AI board members (simulate other investors with board seats)
	numAIBoardMembers := 2 + rand.Intn(2) // 2-3 AI board members

	for i := 0; i < numAIBoardMembers; i++ {
		// AI votes based on their strategy
		voteChance := 0.5
		if vote.VoteType == "acquisition" {
			// AI more likely to accept acquisitions if good terms
			if offerValue, ok := vote.Metadata["offerValue"].(int64); ok {
				if currentVal, ok := vote.Metadata["currentValuation"].(int64); ok {
					if offerValue >= currentVal {
						voteChance = 0.7 // 70% chance to accept good offers
					}
				}
			}
		} else if vote.VoteType == "down_round" {
			// AI less likely to accept down rounds
			voteChance = 0.3
		}

		if rand.Float64() < voteChance {
			aiVotesA++
		} else {
			aiVotesB++
		}
	}

	// Count votes
	totalVotesA := aiVotesA
	totalVotesB := aiVotesB
	if votedForA {
		totalVotesA += playerVotes
	} else {
		totalVotesB += playerVotes
	}

	// Determine outcome
	votePassed := totalVotesA > totalVotesB

	// Store vote result in metadata for execution
	voteCopy := *vote
	voteCopy.Metadata["votePassed"] = votePassed
	voteCopy.Metadata["playerVotedForA"] = votedForA

	// Remove vote from pending list
	gs.PendingBoardVotes = append(gs.PendingBoardVotes[:voteIndex], gs.PendingBoardVotes[voteIndex+1:]...)

	// Generate result message
	voteOutcome := fmt.Sprintf("Board Vote: %d/%d voted for %s, %d/%d voted for %s. ",
		totalVotesA,
		totalVotesA+totalVotesB,
		vote.OptionA,
		totalVotesB,
		totalVotesA+totalVotesB,
		vote.OptionB)

	if votePassed {
		voteOutcome += vote.ConsequenceA
	} else {
		voteOutcome += vote.ConsequenceB
	}

	return voteOutcome, votePassed, nil
}

// ExecuteBoardVoteOutcome executes the outcome of a board vote
func (gs *GameState) ExecuteBoardVoteOutcome(vote BoardVote, passed bool) []string {
	messages := []string{}

	switch vote.VoteType {
	case "acquisition":
		if passed {
			// Acquisition approved - execute it
			if offerValue, ok := vote.Metadata["offerValue"].(int64); ok {
				companyName := vote.CompanyName
				for j := range gs.Portfolio.Investments {
					if gs.Portfolio.Investments[j].CompanyName == companyName {
						inv := &gs.Portfolio.Investments[j]
						payout := int64((inv.EquityPercent / 100.0) * float64(offerValue))
						returnMultiple := float64(payout) / float64(inv.AmountInvested)

						messages = append(messages, fmt.Sprintf(
							"🎉 %s ACQUIRED (Board Approved)! Your %.2f%% = $%s (%.1fx return)",
							companyName,
							inv.EquityPercent,
							formatCurrency(payout),
							returnMultiple,
						))

						gs.Portfolio.Cash += payout
						gs.Portfolio.Investments = append(gs.Portfolio.Investments[:j], gs.Portfolio.Investments[j+1:]...)
						break
					}
				}
			}
		} else {
			messages = append(messages, fmt.Sprintf(
				"❌ %s acquisition REJECTED by board. Company continues operating.",
				vote.CompanyName,
			))
		}
	case "down_round":
		if passed {
			// Down round approved - it proceeds (already handled in ProcessFundingRounds)
			messages = append(messages, fmt.Sprintf(
				"✅ Board approved down round for %s. Round proceeds.",
				vote.CompanyName,
			))
		} else {
			// Down round rejected - company must find alternative or accept worse terms
			messages = append(messages, fmt.Sprintf(
				"❌ Board rejected down round for %s. Company must find alternative funding.",
				vote.CompanyName,
			))
		}
	}

	return messages
}
