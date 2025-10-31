package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// Investment represents a player's investment in a startup
type Investment struct {
	CompanyName      string
	AmountInvested   int64
	EquityPercent    float64
	InitialValuation int64
	CurrentValuation int64
	MonthsHeld       int
}

// Portfolio tracks all player investments
type Portfolio struct {
	Cash        int64
	Investments []Investment
	NetWorth    int64
	Turn        int
	MaxTurns    int
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

// GameState holds the entire game state
type GameState struct {
	PlayerName string
	Portfolio  Portfolio
	AvailableStartups []Startup
	EventPool  []GameEvent
}

// NewGame initializes a new game
func NewGame(playerName string, startingCash int64) *GameState {
	rand.Seed(time.Now().UnixNano())
	
	gs := &GameState{
		PlayerName: playerName,
		Portfolio: Portfolio{
			Cash:     startingCash,
			NetWorth: startingCash,
			Turn:     1,
			MaxTurns: 120, // 10 years
		},
	}
	
	gs.LoadStartups()
	gs.LoadEvents()
	
	return gs
}

// LoadStartups loads all startup companies from JSON files
func (gs *GameState) LoadStartups() {
	gs.AvailableStartups = []Startup{}
	
	for i := 1; i <= 10; i++ {
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
		
		gs.AvailableStartups = append(gs.AvailableStartups, startup)
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
	
	// Calculate equity percentage (simple: investment / valuation)
	equityPercent := (float64(amount) / float64(startup.Valuation)) * 100.0
	
	investment := Investment{
		CompanyName:      startup.Name,
		AmountInvested:   amount,
		EquityPercent:    equityPercent,
		InitialValuation: startup.Valuation,
		CurrentValuation: startup.Valuation,
		MonthsHeld:       0,
	}
	
	gs.Portfolio.Investments = append(gs.Portfolio.Investments, investment)
	gs.Portfolio.Cash -= amount
	gs.updateNetWorth()
	
	return nil
}

// ProcessTurn simulates one month of game time
func (gs *GameState) ProcessTurn() []string {
	messages := []string{}
	
	// Apply random events to each investment
	for i := range gs.Portfolio.Investments {
		inv := &gs.Portfolio.Investments[i]
		inv.MonthsHeld++
		
		// Random chance of an event happening
		if rand.Float64() < 0.3 && len(gs.EventPool) > 0 { // 30% chance per turn
			event := gs.EventPool[rand.Intn(len(gs.EventPool))]
			
			oldVal := inv.CurrentValuation
			inv.CurrentValuation = int64(float64(inv.CurrentValuation) * event.Change)
			
			// Prevent negative valuations
			if inv.CurrentValuation < 0 {
				inv.CurrentValuation = 0
			}
			
			change := inv.CurrentValuation - oldVal
			if change > 0 {
				messages = append(messages, fmt.Sprintf("? %s: %s (+$%d)", inv.CompanyName, event.Event, change))
			} else {
				messages = append(messages, fmt.Sprintf("? %s: %s ($%d)", inv.CompanyName, event.Event, change))
			}
		} else {
			// Natural growth/decline (random walk)
			volatility := 0.05 // 5% volatility
			change := (rand.Float64()*2 - 1) * volatility
			inv.CurrentValuation = int64(float64(inv.CurrentValuation) * (1 + change))
		}
	}
	
	gs.Portfolio.Turn++
	gs.updateNetWorth()
	
	return messages
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

// IsGameOver checks if game has ended
func (gs *GameState) IsGameOver() bool {
	return gs.Portfolio.Turn > gs.Portfolio.MaxTurns
}

// GetFinalScore calculates the final score
func (gs *GameState) GetFinalScore() (netWorth int64, roi float64, successfulExits int) {
	netWorth = gs.Portfolio.NetWorth
	
	// Calculate ROI
	startingCash := int64(250000) // from config
	roi = ((float64(netWorth) - float64(startingCash)) / float64(startingCash)) * 100.0
	
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
	
	// High burn rate = higher risk
	if s.GrossBurnRate > 10 {
		risk += 0.2
	}
	
	// Low sales = higher risk
	if s.MonthlySales < 50 {
		risk += 0.2
	}
	
	// Ensure 0-1 range
	if risk > 1.0 {
		risk = 1.0
	}
	
	return risk
}

func (gs *GameState) calculateGrowthPotential(s *Startup) float64 {
	growth := 0.5
	
	// High activation rate = good growth
	if s.MonthlyActivationRate > 100 {
		growth += 0.2
	}
	
	// Good margins = good growth
	if s.PercentMargin > 25 {
		growth += 0.2
	}
	
	// Ensure 0-1 range
	if growth > 1.0 {
		growth = 1.0
	}
	
	return growth
}
