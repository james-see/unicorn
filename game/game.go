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
	Category         string
	NegativeNewsSent bool // Track if we've already sent negative news for this investment
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

// Difficulty represents game difficulty level
type Difficulty struct {
	Name              string
	StartingCash      int64
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
}

// Predefined difficulty levels
var (
	EasyDifficulty = Difficulty{
		Name:           "Easy",
		StartingCash:   500000,
		EventFrequency: 0.20, // 20% chance
		Volatility:     0.03, // 3% volatility
		MaxTurns:       120,
		Description:    "More cash, lower volatility, fewer bad events",
	}
	
	MediumDifficulty = Difficulty{
		Name:           "Medium",
		StartingCash:   250000,
		EventFrequency: 0.30, // 30% chance
		Volatility:     0.05, // 5% volatility
		MaxTurns:       120,
		Description:    "Standard experience - balanced challenge",
	}
	
	HardDifficulty = Difficulty{
		Name:           "Hard",
		StartingCash:   150000,
		EventFrequency: 0.40, // 40% chance
		Volatility:     0.07, // 7% volatility
		MaxTurns:       120,
		Description:    "Less cash, higher volatility, more events",
	}
	
	ExpertDifficulty = Difficulty{
		Name:           "Expert",
		StartingCash:   100000,
		EventFrequency: 0.50, // 50% chance
		Volatility:     0.10, // 10% volatility
		MaxTurns:       90,   // Only 7.5 years!
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
			Cash:     difficulty.StartingCash,
			NetWorth: difficulty.StartingCash,
			Turn:     1,
			MaxTurns: difficulty.MaxTurns,
		},
	}
	
	gs.LoadStartups()
	gs.LoadEvents()
	
	return gs
}

// LoadStartups loads all startup companies from JSON files
func (gs *GameState) LoadStartups() {
	gs.AvailableStartups = []Startup{}
	
	for i := 1; i <= 20; i++ {
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
		Category:         startup.Category,
		NegativeNewsSent: false,
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
				messages = append(messages, fmt.Sprintf("? %s: %s (+$%d)", inv.CompanyName, event.Event, change))
			} else {
				messages = append(messages, fmt.Sprintf("? %s: %s ($%d)", inv.CompanyName, event.Event, change))
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
