package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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
	Type                string  // "Common", "Preferred", "SAFE", "Convertible"
	HasProRataRights    bool    // Right to participate in future rounds
	HasInfoRights       bool    // Right to company information
	HasBoardSeat        bool    // Board seat (for larger investments)
	BoardSeatMultiplier int     // Number of votes per board seat (1 = normal, 2 = double)
	LiquidationPref     float64 // Liquidation preference (1x, 2x, etc.)
	HasAntiDilution     bool    // Anti-dilution protection
	ConversionDiscount  float64 // Discount on conversion (for SAFE/Convertible)
	ValuationCap        int64   // Valuation cap for SAFE conversion (0 = no cap)
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

	// Founder Relationship (VC Reputation System)
	FounderName       string  // Name of the company's founder
	RelationshipScore float64 // 0-100, quality of relationship with founder
	LastInteraction   int     // Turn number of last interaction
	ValueAddProvided  int     // Count of value-add actions taken
	HasDueDiligence   bool    // Whether DD was performed before investment
	DDLevel           string  // "none", "quick", "standard", "deep"
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
	CarryInterestPaid     int64   // Total carry interest paid to LPs (20% of profits above hurdle)

	// LP Commitments & Capital Calls
	LPCommittedCapital  int64 // Total capital committed by LPs
	LPCalledCapital     int64 // Capital already called from LPs
	LastCapitalCallTurn int   // Last turn when capital was called
	CapitalCallSchedule []int // Scheduled turns for capital calls (e.g., [1, 13, 25, 37, 49])
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
	RevenueHistory          []int64 // Track last 6 months of revenue for trends
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
	PlayerName             string
	PlayerFirmName         string // Player's VC firm name
	Portfolio              Portfolio
	AvailableStartups      []Startup
	EventPool              []GameEvent
	Difficulty             Difficulty
	AIPlayers              []AIPlayer             // Computer opponents
	FundingRoundQueue      []FundingRoundEvent    // Scheduled future funding rounds
	AcquisitionQueue       []AcquisitionEvent     // Scheduled acquisition offers
	DramaticEventQueue     []DramaticEvent        // Scheduled dramatic events (scandals, splits, etc.)
	PendingBoardVotes      []BoardVote            // Board votes requiring player input
	PlayerUpgrades         []string               // Player's purchased upgrades
	InsuranceUsed          bool                   // Track if Portfolio Insurance has been used this game
	ProtectedCompany       string                 // Company protected by Portfolio Insurance
	SyndicateOpportunities []SyndicateOpportunity // Available syndicate deals

	// VC Reputation System
	PlayerReputation      *VCReputation    // Player's reputation across games
	ActiveValueAddActions []ValueAddAction // Ongoing value-add actions
	PendingDDDecisions    []DDDecision     // Due diligence decisions waiting for player
	SecondaryMarketOffers []SecondaryOffer // Offers to buy stakes
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

// SyndicateOpportunity represents a co-investment opportunity with other VCs
type SyndicateOpportunity struct {
	CompanyName      string
	StartupIndex     int
	LeadInvestor     string // Which AI investor is leading (e.g., "Sarah Chen - Accel Partners")
	LeadInvestorFirm string
	TotalRoundSize   int64    // Total amount being raised in this syndicate
	YourMaxShare     int64    // Maximum you can invest (typically 20-40% of round)
	YourMinShare     int64    // Minimum investment to join (typically $25k)
	Valuation        int64    // Company valuation
	Description      string   // Why this is a good deal
	Benefits         []string // Benefits of joining (e.g., "Access to hot deal", "Lower risk")
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

// NewGame initializes a new game with specified difficulty and player upgrades

// PlayerScore represents a player's score in the leaderboard
type PlayerScore struct {
	Name     string
	Firm     string
	NetWorth int64
	ROI      float64
	IsPlayer bool
}

// initializeLPCommitments sets up LP commitments and capital call schedule
func initializeLPCommitments(startingCash int64, maxTurns int) (int64, []int) {
	// Typical VC fund: LPs commit 2-3x the initial fund size
	// Capital calls happen quarterly (every 3 months) or semi-annually
	lpCommittedCapital := startingCash * 2 // 2x committed capital
	capitalCallFrequency := 3              // Quarterly capital calls (every 3 months)
	capitalCallSchedule := []int{}

	// Schedule capital calls: first call at turn 1, then every 3 months
	// Each call is 25% of committed capital (4 calls total over 12 months)
	for turn := 1; turn <= maxTurns; turn += capitalCallFrequency {
		capitalCallSchedule = append(capitalCallSchedule, turn)
	}

	return lpCommittedCapital, capitalCallSchedule
}

// GenerateDefaultFirmName creates a default firm name from player name
// Takes last name and adds "Capital" (e.g., "James Campbell" -> "Campbell Capital")
func GenerateDefaultFirmName(playerName string) string {
	parts := strings.Fields(playerName)
	if len(parts) == 0 {
		return "Your Capital"
	}
	// Use last name if multiple names, otherwise use the whole name
	lastName := parts[len(parts)-1]
	return lastName + " Capital"
}

func NewGame(playerName string, firmName string, difficulty Difficulty, playerUpgrades []string) *GameState {
	rand.Seed(time.Now().UnixNano())

	// Calculate follow-on reserve: $100k base + $50k per potential funding round
	// Assume ~60% of companies will have at least one round we can participate in
	expectedRounds := int64(15 * 0.6 * 2) // 15 companies, 60% raise, avg 2 rounds
	followOnReserve := int64(100000) + (expectedRounds * 50000)

	// Apply upgrades
	startingCash := difficulty.StartingCash
	managementFee := 0.02
	maxTurns := difficulty.MaxTurns

	// Check for upgrades
	for _, upgradeID := range playerUpgrades {
		switch upgradeID {
		case "fund_booster":
			startingCash = int64(float64(startingCash) * 1.1) // +10% cash
		case "management_fee_reduction":
			managementFee = 0.015 // 1.5% instead of 2%
		case "follow_on_reserve_boost":
			followOnReserve += 200000 // +$200k
		case "speed_mode":
			maxTurns = 30 // 30 turns instead of 60
		case "endurance_mode":
			maxTurns = 120 // 120 turns instead of 60
		case "angel_investor":
			startingCash += 100000 // +$100k bonus cash (stacks with Fund Booster)
		}
	}

	// Initialize LP Commitments
	lpCommittedCapital, capitalCallSchedule := initializeLPCommitments(startingCash, maxTurns)

	gs := &GameState{
		PlayerName:     playerName,
		PlayerFirmName: firmName,
		Difficulty:     difficulty,
		PlayerUpgrades: playerUpgrades,
		Portfolio: Portfolio{
			Cash:                startingCash,
			NetWorth:            startingCash,
			Turn:                1,
			MaxTurns:            maxTurns,
			InitialFundSize:     startingCash,
			AnnualManagementFee: managementFee,
			FollowOnReserve:     followOnReserve,
			CarryInterestPaid:   0,
			LPCommittedCapital:  lpCommittedCapital,
			LPCalledCapital:     0,
			LastCapitalCallTurn: 0,
			CapitalCallSchedule: capitalCallSchedule,
		},
	}

	gs.LoadStartups(playerUpgrades)
	gs.LoadEvents()
	gs.InitializeAIPlayers()
	gs.ScheduleFundingRounds()
	gs.ScheduleAcquisitions()
	gs.ScheduleDramaticEvents()

	// Initialize syndicate opportunities (empty for now, generated during investment phase if unlocked)
	gs.SyndicateOpportunities = []SyndicateOpportunity{}

	return gs
}

func (gs *GameState) LoadStartups(playerUpgrades []string) {
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
		startup.RevenueHistory = []int64{startup.MonthlyRevenue} // Initialize with first month

		allStartups = append(allStartups, startup)
	}

	// Randomly select 15 from the 30 startups
	if len(allStartups) > 15 {
		// Apply early_access upgrade - show extra startups
		extraStartups := 0
		for _, upgradeID := range playerUpgrades {
			if upgradeID == "early_access" {
				extraStartups += 2
			}
			if upgradeID == "founder_network" {
				extraStartups += 1 // Add one more startup
			}
		}

		// Shuffle and take first 15+extra
		rand.Shuffle(len(allStartups), func(i, j int) {
			allStartups[i], allStartups[j] = allStartups[j], allStartups[i]
		})
		gs.AvailableStartups = allStartups[:15+extraStartups]
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

func (gs *GameState) ProcessTurn() []string {
	messages := []string{}

	// Process capital calls (before management fees, so fees are charged on larger fund)
	capitalCallMessages := gs.ProcessCapitalCalls()
	messages = append(messages, capitalCallMessages...)

	// Process capital calls for AI players
	for i := range gs.AIPlayers {
		gs.processAICapitalCalls(i)
	}

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

func (gs *GameState) updateNetWorth() {
	netWorth := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve

	for _, inv := range gs.Portfolio.Investments {
		// Value of investment = (equity % / 100) * current valuation
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		netWorth += value
	}

	gs.Portfolio.NetWorth = netWorth
}

func (gs *GameState) GetPortfolioValue() int64 {
	total := int64(0)
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		total += value
	}
	return total
}

func (gs *GameState) GetTotalInvested() int64 {
	total := int64(0)
	for _, inv := range gs.Portfolio.Investments {
		total += inv.AmountInvested
	}
	return total
}

// SectorBreakdownData represents sector performance metrics
type SectorBreakdownData struct {
	Count         int
	TotalInvested int64
	CurrentValue  int64
}

func (gs *GameState) GetSectorBreakdown() map[string]SectorBreakdownData {
	breakdown := make(map[string]SectorBreakdownData)

	for _, inv := range gs.Portfolio.Investments {
		data := breakdown[inv.Category]
		data.Count++
		data.TotalInvested += inv.AmountInvested
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		data.CurrentValue += value
		breakdown[inv.Category] = data
	}

	return breakdown
}

func (gs *GameState) calculateInvestmentROI(inv Investment) float64 {
	if inv.AmountInvested == 0 {
		return 0.0
	}
	currentValue := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
	profit := float64(currentValue - inv.AmountInvested)
	return (profit / float64(inv.AmountInvested)) * 100.0
}

// ROIProjection represents projected ROI data for an investment
type ROIProjection struct {
	CurrentROI      float64 // Current ROI %
	ProjectedROI    float64 // Projected ROI % at game end
	BestCaseROI     float64 // Best case scenario ROI %
	WorstCaseROI    float64 // Worst case scenario ROI %
	ConfidenceLevel float64 // 0-1, confidence in projection
	MonthsRemaining int     // Months until game end
}

// PredictROI calculates projected ROI for an investment based on:
// - Current growth trajectory
// - Risk factors
// - Market conditions
// - Historical performance
func (gs *GameState) PredictROI(inv Investment) ROIProjection {
	monthsRemaining := gs.Portfolio.MaxTurns - gs.Portfolio.Turn
	if monthsRemaining <= 0 {
		monthsRemaining = 1
	}

	currentROI := gs.calculateInvestmentROI(inv)

	// Find the startup
	var startup *Startup
	for i := range gs.AvailableStartups {
		if gs.AvailableStartups[i].Name == inv.CompanyName {
			startup = &gs.AvailableStartups[i]
			break
		}
	}

	if startup == nil {
		return ROIProjection{
			CurrentROI:      currentROI,
			ProjectedROI:    currentROI,
			BestCaseROI:     currentROI,
			WorstCaseROI:    currentROI,
			ConfidenceLevel: 0.0,
			MonthsRemaining: monthsRemaining,
		}
	}

	// Calculate growth rate from rounds (if any)
	growthRate := 1.0 // Base growth rate
	if len(inv.Rounds) > 0 {
		// Calculate average valuation increase per round
		lastRound := inv.Rounds[len(inv.Rounds)-1]
		if inv.InitialValuation > 0 {
			roundGrowth := float64(lastRound.PostMoneyVal) / float64(inv.InitialValuation)
			// Annualize growth (assume rounds happen every 12-18 months)
			monthsSinceFirstRound := float64(lastRound.Month)
			if monthsSinceFirstRound > 0 {
				annualGrowth := math.Pow(roundGrowth, 12.0/monthsSinceFirstRound)
				growthRate = annualGrowth
			}
		}
	} else {
		// No rounds yet, use growth potential
		growthRate = 1.0 + startup.GrowthPotential*0.5 // 0-50% annual growth potential
	}

	// Adjust for risk
	riskAdjustment := 1.0 - (startup.RiskScore * 0.3) // Risk reduces growth by up to 30%
	adjustedGrowthRate := growthRate * riskAdjustment

	// Project forward
	monthsToProject := float64(monthsRemaining)
	projectedMultiplier := math.Pow(adjustedGrowthRate, monthsToProject/12.0)

	// Calculate projected value
	currentValue := float64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
	projectedValue := currentValue * projectedMultiplier
	projectedROI := ((projectedValue - float64(inv.AmountInvested)) / float64(inv.AmountInvested)) * 100.0

	// Best case: 1.5x growth rate
	bestCaseMultiplier := math.Pow(adjustedGrowthRate*1.5, monthsToProject/12.0)
	bestCaseValue := currentValue * bestCaseMultiplier
	bestCaseROI := ((bestCaseValue - float64(inv.AmountInvested)) / float64(inv.AmountInvested)) * 100.0

	// Worst case: 0.5x growth rate (or down round)
	worstCaseMultiplier := math.Pow(adjustedGrowthRate*0.5, monthsToProject/12.0)
	worstCaseValue := currentValue * worstCaseMultiplier
	worstCaseROI := ((worstCaseValue - float64(inv.AmountInvested)) / float64(inv.AmountInvested)) * 100.0

	// Confidence based on:
	// - Number of rounds (more rounds = more data = higher confidence)
	// - Risk level (lower risk = higher confidence)
	// - Time remaining (more time = less confidence)
	confidence := 0.5 // Base confidence
	if len(inv.Rounds) > 0 {
		confidence += 0.2 // More rounds = more confidence
	}
	if startup.RiskScore < 0.3 {
		confidence += 0.2 // Low risk = more confidence
	}
	if monthsRemaining < 12 {
		confidence += 0.1 // Less time = more predictable
	}
	if confidence > 1.0 {
		confidence = 1.0
	}

	return ROIProjection{
		CurrentROI:      currentROI,
		ProjectedROI:    projectedROI,
		BestCaseROI:     bestCaseROI,
		WorstCaseROI:    worstCaseROI,
		ConfidenceLevel: confidence,
		MonthsRemaining: monthsRemaining,
	}
}

func (gs *GameState) GetBestPerformers(count int) []Investment {
	investments := make([]Investment, len(gs.Portfolio.Investments))
	copy(investments, gs.Portfolio.Investments)

	// Sort by ROI (descending)
	for i := 0; i < len(investments)-1; i++ {
		for j := i + 1; j < len(investments); j++ {
			roiI := gs.calculateInvestmentROI(investments[i])
			roiJ := gs.calculateInvestmentROI(investments[j])
			if roiJ > roiI {
				investments[i], investments[j] = investments[j], investments[i]
			}
		}
	}

	if count > len(investments) {
		count = len(investments)
	}
	return investments[:count]
}

func (gs *GameState) GetWorstPerformers(count int) []Investment {
	investments := make([]Investment, len(gs.Portfolio.Investments))
	copy(investments, gs.Portfolio.Investments)

	// Sort by ROI (ascending)
	for i := 0; i < len(investments)-1; i++ {
		for j := i + 1; j < len(investments); j++ {
			roiI := gs.calculateInvestmentROI(investments[i])
			roiJ := gs.calculateInvestmentROI(investments[j])
			if roiJ < roiI {
				investments[i], investments[j] = investments[j], investments[i]
			}
		}
	}

	if count > len(investments) {
		count = len(investments)
	}
	return investments[:count]
}

func (gs *GameState) IsGameOver() bool {
	return gs.Portfolio.Turn > gs.Portfolio.MaxTurns
}

// CalculateCarryInterest calculates projected carry interest based on current portfolio value
// Returns: projected carry, hurdle return, excess profit, and whether carry applies
func (gs *GameState) CalculateCarryInterest() (projectedCarry int64, hurdleReturn float64, excessProfit float64, applies bool) {
	totalStartingCapital := float64(gs.Portfolio.InitialFundSize + gs.Portfolio.FollowOnReserve)

	// Hurdle rate: 8% annual = ~0.67% monthly over 60 months = ~40% total return
	hurdleReturn = totalStartingCapital * 0.40 // 40% hurdle (8% annual over 5 years)

	currentNetWorth := float64(gs.Portfolio.NetWorth)
	profit := currentNetWorth - totalStartingCapital

	if profit > hurdleReturn {
		excessProfit = profit - hurdleReturn
		projectedCarry = int64(excessProfit * 0.20) // 20% carry on excess profit
		applies = true
	} else {
		projectedCarry = 0
		excessProfit = 0
		applies = false
	}

	return projectedCarry, hurdleReturn, excessProfit, applies
}

func (gs *GameState) GetFinalScore() (netWorth int64, roi float64, successfulExits int) {
	netWorth = gs.Portfolio.NetWorth

	// Calculate ROI based on TOTAL starting capital (cash + follow-on reserve)
	totalStartingCapital := gs.Portfolio.InitialFundSize + gs.Portfolio.FollowOnReserve
	roi = ((float64(netWorth) - float64(totalStartingCapital)) / float64(totalStartingCapital)) * 100.0

	// Calculate carry interest (20% of profits above 8% annual hurdle rate)
	carryInterest, _, _, applies := gs.CalculateCarryInterest()

	if applies {
		gs.Portfolio.CarryInterestPaid = carryInterest
		// Net worth after carry goes to LPs
		netWorth -= carryInterest
	} else {
		gs.Portfolio.CarryInterestPaid = 0
	}

	// Recalculate ROI after carry
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
