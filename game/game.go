package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// InvestmentRound represents a funding round
type InvestmentRound struct {
	RoundType        string  // "Seed", "Series A", "Series B", etc.
	PreMoneyValuation int64  // Valuation before this round
	PostMoneyValuation int64  // Valuation after this round
	AmountRaised     int64   // Total amount raised in this round
	InvestorName     string  // Who invested
	AmountInvested   int64   // Amount this investor put in
	EquityPercent    float64 // Percentage owned after this round
	MonthOccurred    int     // Which month this round happened
}

// CompanyInvestment tracks all investments in a company across all rounds
type CompanyInvestment struct {
	CompanyName      string
	Category         string
	InitialValuation int64   // Seed/pre-money valuation
	CurrentValuation int64   // Current post-money valuation
	Rounds           []InvestmentRound
	TotalRaised      int64   // Total raised across all rounds
	MonthsOld        int
	NegativeNewsSent bool
}

// Investment represents a player's investment in a startup
type Investment struct {
	CompanyName      string
	AmountInvested   int64
	EquityPercent    float64 // Current ownership after dilution
	InitialEquityPercent float64 // Original ownership before dilution
	InitialValuation int64
	CurrentValuation int64
	PostMoneyValuation int64 // Post-money valuation at time of investment
	MonthsHeld       int
	Category         string
	NegativeNewsSent bool
	RoundType        string // Which round this investment was made in
}

// Portfolio tracks all player investments
type Portfolio struct {
	Cash            int64
	FundSize        int64 // Total fund size (for management fee calculation)
	ManagementFee   float64 // Annual management fee percentage (e.g., 2.5%)
	ManagementFeesPaid int64 // Total management fees paid
	Investments     []Investment
	NetWorth        int64
	Turn            int
	MaxTurns        int
}

// ComputerPlayer represents an AI VC player
type ComputerPlayer struct {
	Name            string
	Cash            int64
	FundSize        int64
	ManagementFee   float64
	ManagementFeesPaid int64
	Investments     []Investment
	NetWorth        int64
	Personality     string // "Conservative", "Aggressive", "Balanced"
}

// CompanyState tracks the state of each company with all investors
type CompanyState struct {
	CompanyName      string
	InitialValuation int64
	CurrentValuation int64
	PostMoneyValuation int64
	TotalInvestments []InvestmentRound // All investments across all rounds
	Ownership        map[string]float64 // Map of investor name to ownership %
	Category         string
	MonthsOld        int
	RoundNumber      int // Current round number (0=Seed, 1=Series A, etc.)
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
}

// GameEvent represents something that happens to a startup
type GameEvent struct {
	Event       string  `json:"event"`
	Change      float64 `json:"change"` // multiplier (1.5 = +50%, 0.8 = -20%)
	Description string  `json:"description"`
}

// Difficulty represents game difficulty level
type Difficulty struct {
	Name              string
	StartingCash      int64   // Fund size
	ManagementFee     float64 // Annual management fee as percentage (e.g., 2.5 for 2.5%)
	EventFrequency    float64 // 0-1, chance of event per turn
	Volatility        float64 // 0-1, market volatility
	MaxTurns          int
	Description       string
}

// GameState holds the entire game state
type GameState struct {
	PlayerName        string
	Portfolio         Portfolio
	AvailableStartups []Startup
	EventPool         []GameEvent
	Difficulty        Difficulty
	ComputerPlayers   []ComputerPlayer
	CompanyStates     map[string]*CompanyState // Track state of each company
}

// Predefined difficulty levels
var (
	EasyDifficulty = Difficulty{
		Name:          "Easy",
		StartingCash:  1000000, // $1M fund
		ManagementFee: 2.0,     // 2% annual management fee
		EventFrequency: 0.20,   // 20% chance
		Volatility:     0.03,    // 3% volatility
		MaxTurns:       120,
		Description:    "More cash, lower volatility, fewer bad events",
	}
	
	MediumDifficulty = Difficulty{
		Name:          "Medium",
		StartingCash:  500000,  // $500k fund
		ManagementFee: 2.5,     // 2.5% annual management fee
		EventFrequency: 0.30,   // 30% chance
		Volatility:     0.05,   // 5% volatility
		MaxTurns:       120,
		Description:    "Standard experience - balanced challenge",
	}
	
	HardDifficulty = Difficulty{
		Name:          "Hard",
		StartingCash:  500000,  // $500k fund
		ManagementFee: 2.5,     // 2.5% annual management fee
		EventFrequency: 0.40,   // 40% chance
		Volatility:     0.07,   // 7% volatility
		MaxTurns:       120,
		Description:    "Less cash, higher volatility, more events",
	}
	
	ExpertDifficulty = Difficulty{
		Name:          "Expert",
		StartingCash:  500000,  // $500k fund
		ManagementFee: 2.5,     // 2.5% annual management fee
		EventFrequency: 0.50,   // 50% chance
		Volatility:     0.10,   // 10% volatility
		MaxTurns:       90,     // Only 7.5 years!
		Description:    "Brutal - minimal cash, extreme volatility, shorter time",
	}
)

// NewGame initializes a new game with specified difficulty
func NewGame(playerName string, difficulty Difficulty) *GameState {
	rand.Seed(time.Now().UnixNano())
	
	gs := &GameState{
		PlayerName: playerName,
		Difficulty: difficulty,
		Portfolio: Portfolio{
			Cash:          difficulty.StartingCash,
			FundSize:      difficulty.StartingCash,
			ManagementFee: difficulty.ManagementFee,
			NetWorth:      difficulty.StartingCash,
			Turn:          1,
			MaxTurns:      difficulty.MaxTurns,
		},
		CompanyStates: make(map[string]*CompanyState),
	}
	
	gs.LoadStartups()
	gs.LoadEvents()
	gs.InitializeComputerPlayers()
	
	return gs
}

// InitializeComputerPlayers creates AI VC players
func (gs *GameState) InitializeComputerPlayers() {
	// CARL - Default computer player
	carl := ComputerPlayer{
		Name:          "CARL (LP at Sterling and Cooper)",
		Cash:          gs.Difficulty.StartingCash,
		FundSize:      gs.Difficulty.StartingCash,
		ManagementFee: gs.Difficulty.ManagementFee,
		Personality:   "Balanced",
		Investments:   []Investment{},
		NetWorth:      gs.Difficulty.StartingCash,
	}
	
	// Add 1-2 more computer players with different personalities
	players := []ComputerPlayer{carl}
	
	// Add aggressive player
	if rand.Float64() < 0.7 {
		aggressive := ComputerPlayer{
			Name:          "TechVentures Fund",
			Cash:          gs.Difficulty.StartingCash,
			FundSize:      gs.Difficulty.StartingCash,
			ManagementFee: gs.Difficulty.ManagementFee,
			Personality:   "Aggressive",
			Investments:   []Investment{},
			NetWorth:      gs.Difficulty.StartingCash,
		}
		players = append(players, aggressive)
	}
	
	// Add conservative player
	if rand.Float64() < 0.7 {
		conservative := ComputerPlayer{
			Name:          "Stable Capital Partners",
			Cash:          gs.Difficulty.StartingCash,
			FundSize:      gs.Difficulty.StartingCash,
			ManagementFee: gs.Difficulty.ManagementFee,
			Personality:   "Conservative",
			Investments:   []Investment{},
			NetWorth:      gs.Difficulty.StartingCash,
		}
		players = append(players, conservative)
	}
	
	gs.ComputerPlayers = players
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
		
		// Convert valuation to proper scale (assume values are in millions)
		if startup.Valuation < 1000 {
			startup.Valuation *= 1000000
		}
		
		// Calculate risk and growth scores based on metrics
		startup.RiskScore = gs.calculateRiskScore(&startup)
		startup.GrowthPotential = gs.calculateGrowthPotential(&startup)
		
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

// MakeInvestment allows player to invest in a startup
func (gs *GameState) MakeInvestment(startupIndex int, amount int64) error {
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
	
	// Get or create company state
	companyState, exists := gs.CompanyStates[startup.Name]
	if !exists {
		companyState = &CompanyState{
			CompanyName:      startup.Name,
			InitialValuation: startup.Valuation,
			CurrentValuation: startup.Valuation,
			PostMoneyValuation: startup.Valuation,
			Category:         startup.Category,
			Ownership:        make(map[string]float64),
			RoundNumber:      0, // Start at Seed round
		}
		gs.CompanyStates[startup.Name] = companyState
	}
	
	// Determine round type
	roundType := "Seed"
	if companyState.RoundNumber == 1 {
		roundType = "Series A"
	} else if companyState.RoundNumber == 2 {
		roundType = "Series B"
	} else if companyState.RoundNumber >= 3 {
		roundType = fmt.Sprintf("Series %c", 'A'+rune(companyState.RoundNumber))
	}
	
	// Pre-money valuation is current valuation
	preMoneyValuation := companyState.CurrentValuation
	
	// Calculate post-money valuation (pre-money + amount invested)
	postMoneyValuation := preMoneyValuation + amount
	
	// Calculate equity percentage based on post-money valuation
	equityPercent := (float64(amount) / float64(postMoneyValuation)) * 100.0
	
	// Create investment round record
	round := InvestmentRound{
		RoundType:         roundType,
		PreMoneyValuation: preMoneyValuation,
		PostMoneyValuation: postMoneyValuation,
		AmountRaised:     amount,
		InvestorName:     gs.PlayerName,
		AmountInvested:   amount,
		EquityPercent:    equityPercent,
		MonthOccurred:    gs.Portfolio.Turn,
	}
	
	// Apply dilution to existing investors
	gs.applyDilution(startup.Name, postMoneyValuation)
	
	// Update company state
	companyState.PostMoneyValuation = postMoneyValuation
	companyState.CurrentValuation = postMoneyValuation
	companyState.TotalInvestments = append(companyState.TotalInvestments, round)
	if companyState.Ownership[gs.PlayerName] == 0 {
		companyState.Ownership[gs.PlayerName] = equityPercent
	} else {
		companyState.Ownership[gs.PlayerName] += equityPercent
	}
	
	// Create investment record
	investment := Investment{
		CompanyName:        startup.Name,
		AmountInvested:     amount,
		EquityPercent:      equityPercent,
		InitialEquityPercent: equityPercent,
		InitialValuation:   preMoneyValuation,
		CurrentValuation:   postMoneyValuation,
		PostMoneyValuation: postMoneyValuation,
		MonthsHeld:         0,
		Category:           startup.Category,
		NegativeNewsSent:   false,
		RoundType:          roundType,
	}
	
	gs.Portfolio.Investments = append(gs.Portfolio.Investments, investment)
	gs.Portfolio.Cash -= amount
	gs.updateNetWorth()
	
	return nil
}

// applyDilution applies dilution to all existing investors when a new round occurs
func (gs *GameState) applyDilution(companyName string, newPostMoneyValuation int64) {
	// Dilute player's investments
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			inv := &gs.Portfolio.Investments[i]
			// Dilution: new ownership = old ownership * (old post-money / new post-money)
			oldPostMoney := inv.PostMoneyValuation
			if oldPostMoney > 0 {
				inv.EquityPercent = inv.EquityPercent * (float64(oldPostMoney) / float64(newPostMoneyValuation))
				inv.PostMoneyValuation = newPostMoneyValuation
			}
		}
	}
	
	// Dilute computer players' investments
	for p := range gs.ComputerPlayers {
		for i := range gs.ComputerPlayers[p].Investments {
			if gs.ComputerPlayers[p].Investments[i].CompanyName == companyName {
				inv := &gs.ComputerPlayers[p].Investments[i]
				oldPostMoney := inv.PostMoneyValuation
				if oldPostMoney > 0 {
					inv.EquityPercent = inv.EquityPercent * (float64(oldPostMoney) / float64(newPostMoneyValuation))
					inv.PostMoneyValuation = newPostMoneyValuation
				}
			}
		}
	}
}

// ProcessTurn simulates one month of game time
func (gs *GameState) ProcessTurn() []string {
	messages := []string{}
	
	// Pay management fees annually (every 12 months)
	if gs.Portfolio.Turn%12 == 0 {
		monthlyFee := (gs.Portfolio.ManagementFee / 100.0) * float64(gs.Portfolio.FundSize) / 12.0
		feeAmount := int64(monthlyFee)
		if gs.Portfolio.Cash >= feeAmount {
			gs.Portfolio.Cash -= feeAmount
			gs.Portfolio.ManagementFeesPaid += feeAmount
			messages = append(messages, fmt.Sprintf("?? Management fee: -$%d (%.1f%% annually)", feeAmount, gs.Portfolio.ManagementFee))
		}
		
		// Pay fees for computer players
		for p := range gs.ComputerPlayers {
			cp := &gs.ComputerPlayers[p]
			cpMonthlyFee := (cp.ManagementFee / 100.0) * float64(cp.FundSize) / 12.0
			cpFeeAmount := int64(cpMonthlyFee)
			if cp.Cash >= cpFeeAmount {
				cp.Cash -= cpFeeAmount
				cp.ManagementFeesPaid += cpFeeAmount
			}
		}
	}
	
	// Let computer players make investments (at start of game and periodically)
	if gs.Portfolio.Turn == 1 || (gs.Portfolio.Turn > 1 && rand.Float64() < 0.15) {
		gs.processComputerPlayerInvestments(messages)
	}
	
	// Process each company's valuation changes
	for companyName, companyState := range gs.CompanyStates {
		companyState.MonthsOld++
		
		oldValuation := companyState.CurrentValuation
		
		// Random chance of an event happening (based on difficulty)
		if rand.Float64() < gs.Difficulty.EventFrequency && len(gs.EventPool) > 0 {
			event := gs.EventPool[rand.Intn(len(gs.EventPool))]
			
			companyState.CurrentValuation = int64(float64(companyState.CurrentValuation) * event.Change)
			
			// Prevent negative valuations
			if companyState.CurrentValuation < 0 {
				companyState.CurrentValuation = 0
			}
			
			change := companyState.CurrentValuation - oldValuation
			if change > 0 {
				messages = append(messages, fmt.Sprintf("?? %s: %s (+$%d)", companyName, event.Event, change))
			} else {
				messages = append(messages, fmt.Sprintf("?? %s: %s ($%d)", companyName, event.Event, change))
			}
		} else {
			// Natural growth/decline (random walk) - volatility based on difficulty
			change := (rand.Float64()*2 - 1) * gs.Difficulty.Volatility
			companyState.CurrentValuation = int64(float64(companyState.CurrentValuation) * (1 + change))
		}
		
		// Update all investments in this company
		gs.updateCompanyInvestments(companyName, companyState.CurrentValuation)
		
		// Chance for new funding round (Series A, B, etc.)
		if companyState.MonthsOld >= 6 && rand.Float64() < 0.05 {
			gs.processNewFundingRound(companyName, messages)
		}
	}
	
	gs.Portfolio.Turn++
	gs.updateNetWorth()
	gs.updateComputerPlayerNetWorth()
	
	return messages
}

// processComputerPlayerInvestments allows AI players to make investments
func (gs *GameState) processComputerPlayerInvestments(messages []string) {
	for p := range gs.ComputerPlayers {
		cp := &gs.ComputerPlayers[p]
		
		// Skip if no cash
		if cp.Cash < 10000 {
			continue
		}
		
		// AI decision making based on personality
		for i, startup := range gs.AvailableStartups {
			// Check if already invested in this company
			alreadyInvested := false
			for _, inv := range cp.Investments {
				if inv.CompanyName == startup.Name {
					alreadyInvested = true
					break
				}
			}
			
			if alreadyInvested {
				continue
			}
			
			shouldInvest := false
			investmentAmount := int64(0)
			
			switch cp.Personality {
			case "Aggressive":
				// Aggressive: Invest in high-growth, high-risk companies
				if startup.GrowthPotential > 0.6 && cp.Cash >= 50000 {
					shouldInvest = true
					investmentAmount = int64(float64(cp.Cash) * 0.15) // 15% of cash
					if investmentAmount < 50000 {
						investmentAmount = 50000
					}
				}
			case "Conservative":
				// Conservative: Invest in low-risk, stable companies
				if startup.RiskScore < 0.4 && cp.Cash >= 30000 {
					shouldInvest = true
					investmentAmount = int64(float64(cp.Cash) * 0.10) // 10% of cash
					if investmentAmount < 30000 {
						investmentAmount = 30000
					}
				}
			case "Balanced":
				// Balanced: Mix of risk and growth
				if startup.GrowthPotential > 0.5 && startup.RiskScore < 0.7 && cp.Cash >= 40000 {
					shouldInvest = true
					investmentAmount = int64(float64(cp.Cash) * 0.12) // 12% of cash
					if investmentAmount < 40000 {
						investmentAmount = 40000
					}
				}
			}
			
			if shouldInvest && investmentAmount <= cp.Cash {
				// Make the investment similar to player investment
				startup := gs.AvailableStartups[i]
				companyState, exists := gs.CompanyStates[startup.Name]
				if !exists {
					companyState = &CompanyState{
						CompanyName:      startup.Name,
						InitialValuation: startup.Valuation,
						CurrentValuation: startup.Valuation,
						PostMoneyValuation: startup.Valuation,
						Category:         startup.Category,
						Ownership:        make(map[string]float64),
						RoundNumber:      0,
					}
					gs.CompanyStates[startup.Name] = companyState
				}
				
				roundType := "Seed"
				if companyState.RoundNumber == 1 {
					roundType = "Series A"
				} else if companyState.RoundNumber >= 2 {
					roundType = fmt.Sprintf("Series %c", 'A'+rune(companyState.RoundNumber))
				}
				
				preMoneyValuation := companyState.CurrentValuation
				postMoneyValuation := preMoneyValuation + investmentAmount
				equityPercent := (float64(investmentAmount) / float64(postMoneyValuation)) * 100.0
				
				// Apply dilution
				gs.applyDilution(startup.Name, postMoneyValuation)
				
				// Update company state
				companyState.PostMoneyValuation = postMoneyValuation
				companyState.CurrentValuation = postMoneyValuation
				companyState.TotalInvestments = append(companyState.TotalInvestments, InvestmentRound{
					RoundType:         roundType,
					PreMoneyValuation: preMoneyValuation,
					PostMoneyValuation: postMoneyValuation,
					AmountRaised:     investmentAmount,
					InvestorName:     cp.Name,
					AmountInvested:   investmentAmount,
					EquityPercent:    equityPercent,
					MonthOccurred:    gs.Portfolio.Turn,
				})
				if companyState.Ownership[cp.Name] == 0 {
					companyState.Ownership[cp.Name] = equityPercent
				} else {
					companyState.Ownership[cp.Name] += equityPercent
				}
				
				// Create investment record
				investment := Investment{
					CompanyName:        startup.Name,
					AmountInvested:     investmentAmount,
					EquityPercent:      equityPercent,
					InitialEquityPercent: equityPercent,
					InitialValuation:   preMoneyValuation,
					CurrentValuation:   postMoneyValuation,
					PostMoneyValuation: postMoneyValuation,
					MonthsHeld:         0,
					Category:           startup.Category,
					NegativeNewsSent:   false,
					RoundType:          roundType,
				}
				
				cp.Investments = append(cp.Investments, investment)
				cp.Cash -= investmentAmount
				
				messages = append(messages, fmt.Sprintf("?? %s invested $%d in %s (%s round)", cp.Name, investmentAmount, startup.Name, roundType))
			}
		}
	}
}

// processNewFundingRound handles a new funding round for a company
func (gs *GameState) processNewFundingRound(companyName string, messages []string) {
	companyState, exists := gs.CompanyStates[companyName]
	if !exists {
		return
	}
	
	// Advance to next round
	companyState.RoundNumber++
	roundType := "Series A"
	if companyState.RoundNumber == 2 {
		roundType = "Series B"
	} else if companyState.RoundNumber >= 3 {
		roundType = fmt.Sprintf("Series %c", 'A'+rune(companyState.RoundNumber))
	}
	
	// New round typically raises 2-5x the current valuation
	preMoneyValuation := companyState.CurrentValuation
	roundMultiplier := 1.5 + rand.Float64()*3.5 // 1.5x to 5x
	amountRaised := int64(float64(preMoneyValuation) * roundMultiplier)
	postMoneyValuation := preMoneyValuation + amountRaised
	
	// Apply dilution
	gs.applyDilution(companyName, postMoneyValuation)
	
	// Update company state
	companyState.PostMoneyValuation = postMoneyValuation
	companyState.CurrentValuation = postMoneyValuation
	
	// Randomly assign new investors (could be computer players or external)
	investors := []string{"External VC", "Corporate Investor", "Strategic Partner"}
	if len(gs.ComputerPlayers) > 0 {
		investors = append(investors, gs.ComputerPlayers[rand.Intn(len(gs.ComputerPlayers))].Name)
	}
	
	investorName := investors[rand.Intn(len(investors))]
	equityPercent := (float64(amountRaised) / float64(postMoneyValuation)) * 100.0
	
	companyState.TotalInvestments = append(companyState.TotalInvestments, InvestmentRound{
		RoundType:         roundType,
		PreMoneyValuation: preMoneyValuation,
		PostMoneyValuation: postMoneyValuation,
		AmountRaised:     amountRaised,
		InvestorName:     investorName,
		AmountInvested:   amountRaised,
		EquityPercent:    equityPercent,
		MonthOccurred:    gs.Portfolio.Turn,
	})
	
	messages = append(messages, fmt.Sprintf("?? %s raised $%d in %s (pre-money: $%d, post-money: $%d) - Your ownership diluted!", companyName, amountRaised, roundType, preMoneyValuation, postMoneyValuation))
}

// updateCompanyInvestments updates all investments for a company when valuation changes
func (gs *GameState) updateCompanyInvestments(companyName string, newValuation int64) {
	// Update player investments
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			gs.Portfolio.Investments[i].CurrentValuation = newValuation
			gs.Portfolio.Investments[i].MonthsHeld++
		}
	}
	
	// Update computer player investments
	for p := range gs.ComputerPlayers {
		for i := range gs.ComputerPlayers[p].Investments {
			if gs.ComputerPlayers[p].Investments[i].CompanyName == companyName {
				gs.ComputerPlayers[p].Investments[i].CurrentValuation = newValuation
				gs.ComputerPlayers[p].Investments[i].MonthsHeld++
			}
		}
	}
	
	// Update company state
	if companyState, exists := gs.CompanyStates[companyName]; exists {
		companyState.CurrentValuation = newValuation
	}
}

// updateComputerPlayerNetWorth updates net worth for all computer players
func (gs *GameState) updateComputerPlayerNetWorth() {
	for p := range gs.ComputerPlayers {
		cp := &gs.ComputerPlayers[p]
		cp.NetWorth = cp.Cash
		
		for _, inv := range cp.Investments {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			cp.NetWorth += value
		}
	}
}

// updateNetWorth calculates current net worth
func (gs *GameState) updateNetWorth() {
	netWorth := gs.Portfolio.Cash
	
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

// GetCompanyOwnership returns ownership percentages for all investors in a company
func (gs *GameState) GetCompanyOwnership(companyName string) map[string]float64 {
	if companyState, exists := gs.CompanyStates[companyName]; exists {
		return companyState.Ownership
	}
	return make(map[string]float64)
}

// GetVCLeaderboard returns all VCs sorted by net worth
func (gs *GameState) GetVCLeaderboard() []struct {
	Name     string
	NetWorth int64
	IsPlayer bool
} {
	leaderboard := []struct {
		Name     string
		NetWorth int64
		IsPlayer bool
	}{
		{Name: gs.PlayerName, NetWorth: gs.Portfolio.NetWorth, IsPlayer: true},
	}
	
	for _, cp := range gs.ComputerPlayers {
		leaderboard = append(leaderboard, struct {
			Name     string
			NetWorth int64
			IsPlayer bool
		}{Name: cp.Name, NetWorth: cp.NetWorth, IsPlayer: false})
	}
	
	// Sort by net worth (descending)
	for i := 0; i < len(leaderboard)-1; i++ {
		for j := i + 1; j < len(leaderboard); j++ {
			if leaderboard[j].NetWorth > leaderboard[i].NetWorth {
				leaderboard[i], leaderboard[j] = leaderboard[j], leaderboard[i]
			}
		}
	}
	
	return leaderboard
}

// IsGameOver checks if game has ended
func (gs *GameState) IsGameOver() bool {
	return gs.Portfolio.Turn > gs.Portfolio.MaxTurns
}

// GetFinalScore calculates the final score
func (gs *GameState) GetFinalScore() (netWorth int64, roi float64, successfulExits int) {
	netWorth = gs.Portfolio.NetWorth
	
	// Calculate ROI based on starting cash for this difficulty
	roi = ((float64(netWorth) - float64(gs.Difficulty.StartingCash)) / float64(gs.Difficulty.StartingCash)) * 100.0
	
	// Count successful exits (investments that 5x'd or more)
	successfulExits = 0
	for _, inv := range gs.Portfolio.Investments {
		if inv.CurrentValuation >= inv.InitialValuation*5 {
			successfulExits++
		}
	}
	
	return netWorth, roi, successfulExits
}

// Helper functions
func (gs *GameState) calculateRiskScore(s *Startup) float64 {
	risk := 0.5
	
	// Very high burn rate = VERY HIGH risk
	if s.GrossBurnRate > 40 {
		risk += 0.4
	} else if s.GrossBurnRate > 20 {
		risk += 0.3
	} else if s.GrossBurnRate > 10 {
		risk += 0.2
	} else if s.GrossBurnRate <= 3 {
		// Low burn rate = lower risk
		risk -= 0.2
	}
	
	// Very low sales = VERY HIGH risk
	if s.MonthlySales < 5 {
		risk += 0.4
	} else if s.MonthlySales < 20 {
		risk += 0.3
	} else if s.MonthlySales < 50 {
		risk += 0.2
	} else if s.MonthlySales > 300 {
		// High sales = lower risk
		risk -= 0.2
	}
	
	// Ensure 0-1 range
	if risk > 1.0 {
		risk = 1.0
	}
	if risk < 0.0 {
		risk = 0.0
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
