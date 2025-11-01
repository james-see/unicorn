package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// FundingRound represents a funding round for a startup
type FundingRound struct {
	RoundName        string  // "Seed", "Series A", "Series B", etc.
	PreMoneyVal      int64   // Pre-money valuation
	InvestmentAmount int64   // Total raised in this round
	PostMoneyVal     int64   // Post-money valuation
	Month            int     // Game turn when this happened
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
	NegativeNewsSent bool // Track if we've already sent negative news for this investment
	Rounds           []FundingRound // Track all funding rounds
}

// Portfolio tracks all player investments
type Portfolio struct {
	Cash                int64
	Investments         []Investment
	NetWorth            int64
	Turn                int
	MaxTurns            int
	InitialFundSize     int64   // Original fund size
	ManagementFeesCharged int64 // Total management fees paid
	AnnualManagementFee float64 // Annual management fee rate (e.g., 0.02 for 2%)
	FollowOnReserve     int64   // Reserve fund for follow-on investments
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
	StartingCash      int64
	EventFrequency    float64 // 0-1, chance of event per turn
	Volatility        float64 // 0-1, market volatility
	MaxTurns          int
	Description       string
}

// AIPlayer represents a computer-controlled VC
type AIPlayer struct {
	Name            string
	Firm            string
	Portfolio       Portfolio
	Strategy        string  // "aggressive", "balanced", "conservative"
	RiskTolerance   float64 // 0-1
}

// GameState holds the entire game state
type GameState struct {
	PlayerName        string
	Portfolio         Portfolio
	AvailableStartups []Startup
	EventPool         []GameEvent
	Difficulty        Difficulty
	AIPlayers         []AIPlayer // Computer opponents
	FundingRoundQueue []FundingRoundEvent // Scheduled future funding rounds
}

// FundingRoundEvent represents a scheduled funding round
type FundingRoundEvent struct {
	CompanyName string
	RoundName   string
	ScheduledTurn int
	RaiseAmount int64
}

// FollowOnOpportunity represents a chance to invest more in a company raising a round
type FollowOnOpportunity struct {
	CompanyName      string
	RoundName        string
	PreMoneyVal      int64
	PostMoneyVal     int64
	CurrentEquity    float64
	MinInvestment    int64
	MaxInvestment    int64
}

// Predefined difficulty levels
var (
	EasyDifficulty = Difficulty{
		Name:           "Easy",
		StartingCash:   1000000, // $1M fund
		EventFrequency: 0.20,    // 20% chance
		Volatility:     0.03,    // 3% volatility
		MaxTurns:       120,
		Description:    "$1M fund, lower volatility, fewer bad events",
	}
	
	MediumDifficulty = Difficulty{
		Name:           "Medium",
		StartingCash:   750000, // $750k fund
		EventFrequency: 0.30,   // 30% chance
		Volatility:     0.05,   // 5% volatility
		MaxTurns:       120,
		Description:    "$750k fund - balanced challenge",
	}
	
	HardDifficulty = Difficulty{
		Name:           "Hard",
		StartingCash:   500000, // $500k fund
		EventFrequency: 0.40,   // 40% chance
		Volatility:     0.07,   // 7% volatility
		MaxTurns:       120,
		Description:    "$500k fund, higher volatility, more events",
	}
	
	ExpertDifficulty = Difficulty{
		Name:           "Expert",
		StartingCash:   500000, // $500k fund
		EventFrequency: 0.50,   // 50% chance
		Volatility:     0.10,   // 10% volatility
		MaxTurns:       90,     // Only 7.5 years!
		Description:    "$500k fund, extreme volatility, shorter time",
	}
)

// NewGame initializes a new game with specified difficulty
func NewGame(playerName string, difficulty Difficulty) *GameState {
	rand.Seed(time.Now().UnixNano())
	
	gs := &GameState{
		PlayerName: playerName,
		Difficulty: difficulty,
		Portfolio: Portfolio{
			Cash:                difficulty.StartingCash,
			NetWorth:            difficulty.StartingCash,
			Turn:                1,
			MaxTurns:            difficulty.MaxTurns,
			InitialFundSize:     difficulty.StartingCash,
			AnnualManagementFee: 0.02, // 2% annual management fee
			FollowOnReserve:     100000, // $100k for follow-on investments
		},
	}
	
	gs.LoadStartups()
	gs.LoadEvents()
	gs.InitializeAIPlayers()
	gs.ScheduleFundingRounds()
	
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
	
	// Check if already invested in this company
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == startup.Name {
			return fmt.Errorf("you have already invested in %s", startup.Name)
		}
	}
	
	// Calculate equity percentage (simple: investment / valuation)
	equityPercent := (float64(amount) / float64(startup.Valuation)) * 100.0
	
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
							maxInvestment := gs.Portfolio.FollowOnReserve
							if maxInvestment > event.RaiseAmount / 2 {
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
	
	if amount > gs.Portfolio.FollowOnReserve {
		return fmt.Errorf("insufficient follow-on funds (have $%d, need $%d)", gs.Portfolio.FollowOnReserve, amount)
	}
	
	// Find the investment
	for i := range gs.Portfolio.Investments {
		if gs.Portfolio.Investments[i].CompanyName == companyName {
			inv := &gs.Portfolio.Investments[i]
			
			// Find the company valuation
			for _, startup := range gs.AvailableStartups {
				if startup.Name == companyName {
					// Calculate additional equity gained
					// New equity = (investment / post-money valuation) * 100
					additionalEquity := (float64(amount) / float64(startup.Valuation)) * 100.0
					
					inv.AmountInvested += amount
					inv.EquityPercent += additionalEquity
					gs.Portfolio.FollowOnReserve -= amount
					gs.updateNetWorth()
					
					return nil
				}
			}
			return fmt.Errorf("company not found")
		}
	}
	
	return fmt.Errorf("you have not invested in %s", companyName)
}

// ProcessTurn simulates one month of game time
func (gs *GameState) ProcessTurn() []string {
	messages := []string{}
	
	// Process management fees
	feeMessages := gs.ProcessManagementFees()
	messages = append(messages, feeMessages...)
	
	// Process funding rounds
	roundMessages := gs.ProcessFundingRounds()
	messages = append(messages, roundMessages...)
	
	// Apply random events to each investment
	for i := range gs.Portfolio.Investments {
		inv := &gs.Portfolio.Investments[i]
		inv.MonthsHeld++
		
		wasAboveInitial := inv.CurrentValuation >= inv.InitialValuation
		
		// Random chance of an event happening (based on difficulty)
		if rand.Float64() < gs.Difficulty.EventFrequency && len(gs.EventPool) > 0 {
			event := gs.EventPool[rand.Intn(len(gs.EventPool))]
			
			oldVal := inv.CurrentValuation
			inv.CurrentValuation = int64(float64(inv.CurrentValuation) * event.Change)
			
			// Prevent negative valuations
			if inv.CurrentValuation < 0 {
				inv.CurrentValuation = 0
			}
			
			change := inv.CurrentValuation - oldVal
			if change > 0 {
				messages = append(messages, fmt.Sprintf("?? %s: %s (+$%d)", inv.CompanyName, event.Event, change))
			} else {
				messages = append(messages, fmt.Sprintf("?? %s: %s ($%d)", inv.CompanyName, event.Event, change))
			}
		} else {
			// Natural growth/decline (random walk) - volatility based on difficulty
			change := (rand.Float64()*2 - 1) * gs.Difficulty.Volatility
			inv.CurrentValuation = int64(float64(inv.CurrentValuation) * (1 + change))
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

// InitializeAIPlayers creates computer VC opponents
func (gs *GameState) InitializeAIPlayers() {
	gs.AIPlayers = []AIPlayer{
		{
			Name:     "CARL",
			Firm:     "Sterling & Cooper",
			Strategy: "conservative",
			RiskTolerance: 0.3,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     100000,
			},
		},
		{
			Name:     "Sarah Chen",
			Firm:     "Accel Partners",
			Strategy: "aggressive",
			RiskTolerance: 0.8,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     100000,
			},
		},
		{
			Name:     "Marcus Williams",
			Firm:     "Sequoia Capital",
			Strategy: "balanced",
			RiskTolerance: 0.5,
			Portfolio: Portfolio{
				Cash:                gs.Difficulty.StartingCash,
				NetWorth:            gs.Difficulty.StartingCash,
				Turn:                1,
				MaxTurns:            gs.Difficulty.MaxTurns,
				InitialFundSize:     gs.Difficulty.StartingCash,
				AnnualManagementFee: 0.02,
				FollowOnReserve:     100000,
			},
		},
	}
}

// ScheduleFundingRounds schedules future funding rounds for companies
func (gs *GameState) ScheduleFundingRounds() {
	gs.FundingRoundQueue = []FundingRoundEvent{}
	
	// Schedule funding rounds for each startup at random intervals
	for _, startup := range gs.AvailableStartups {
		// Seed round (6-12 months)
		seedTurn := 6 + rand.Intn(7)
		if seedTurn < gs.Portfolio.MaxTurns {
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Seed",
				ScheduledTurn: seedTurn,
				RaiseAmount:   startup.Valuation / 10, // Raise 10% of current valuation
			})
		}
		
		// Series A (18-36 months)
		seriesATurn := 18 + rand.Intn(19)
		if seriesATurn < gs.Portfolio.MaxTurns {
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Series A",
				ScheduledTurn: seriesATurn,
				RaiseAmount:   startup.Valuation / 5, // Raise 20% of current valuation
			})
		}
		
		// Series B (36-60 months)
		seriesBTurn := 36 + rand.Intn(25)
		if seriesBTurn < gs.Portfolio.MaxTurns {
			gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
				CompanyName:   startup.Name,
				RoundName:     "Series B",
				ScheduledTurn: seriesBTurn,
				RaiseAmount:   startup.Valuation / 3, // Raise 33% of current valuation
			})
		}
		
		// Series C (60-90 months) - only for some companies
		if rand.Float64() < 0.5 { // 50% of companies get Series C
			seriesCTurn := 60 + rand.Intn(31)
			if seriesCTurn < gs.Portfolio.MaxTurns {
				gs.FundingRoundQueue = append(gs.FundingRoundQueue, FundingRoundEvent{
					CompanyName:   startup.Name,
					RoundName:     "Series C",
					ScheduledTurn: seriesCTurn,
					RaiseAmount:   startup.Valuation / 2, // Raise 50% of current valuation
				})
			}
		}
	}
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
					
					preMoneyVal := startup.Valuation
					postMoneyVal := preMoneyVal + event.RaiseAmount
					
					// Calculate dilution for existing investors
					dilutionFactor := float64(preMoneyVal) / float64(postMoneyVal)
					
					// Update player's investment if they invested in this company
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]
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
							
							messages = append(messages, fmt.Sprintf(
								"?? %s raised $%d in %s round! Your equity diluted from %.2f%% to %.2f%%",
								event.CompanyName,
								event.RaiseAmount,
								event.RoundName,
								oldEquity,
								inv.EquityPercent,
							))
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
				
				if investmentAmount > 10000 { // Minimum investment
					equityPercent := (float64(investmentAmount) / float64(startup.Valuation)) * 100.0
					
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
	
	// Add player
	playerROI := ((float64(gs.Portfolio.NetWorth) - float64(gs.Portfolio.InitialFundSize)) / float64(gs.Portfolio.InitialFundSize)) * 100.0
	scores = append(scores, PlayerScore{
		Name:     gs.PlayerName,
		Firm:     "Your Fund",
		NetWorth: gs.Portfolio.NetWorth,
		ROI:      playerROI,
		IsPlayer: true,
	})
	
	// Add AI players
	for _, ai := range gs.AIPlayers {
		aiROI := ((float64(ai.Portfolio.NetWorth) - float64(ai.Portfolio.InitialFundSize)) / float64(ai.Portfolio.InitialFundSize)) * 100.0
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
