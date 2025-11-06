package game

import (
	"fmt"
	"math/rand"
)

// MarketCycle represents the current economic environment
type MarketCycle struct {
	Name             string  // "Bull Market", "Bear Market", "Recession", "Boom"
	StartTurn        int
	Duration         int
	ValuationMult    float64 // 0.7-1.5x multiplier on all valuations
	FundingEase      float64 // 0.5-2.0x multiplier on funding availability
	ExitLikelihood   float64 // 0.5-1.5x multiplier on exit chances
	Description      string
	Color            string // For UI display
}

// EconomicEvent represents a specific economic occurrence
type EconomicEvent struct {
	Name        string
	Description string
	Turn        int
	Duration    int
	Impact      EventImpact
	Sectors     []string // Affected sectors, empty = all
	Active      bool
}

// EventImpact defines the effect of an economic event
type EventImpact struct {
	ValuationChange  float64 // Percentage change to valuations
	FundingChange    float64 // Percentage change to funding ease
	RevenueChange    float64 // Percentage change to revenue growth
	ChurnChange      float64 // Percentage change to churn rate
}

// Market cycle types
const (
	MarketBull      = "Bull Market"
	MarketBear      = "Bear Market"
	MarketRecession = "Recession"
	MarketBoom      = "Boom"
	MarketNormal    = "Normal Market"
)

// InitializeMarketCycle sets up the starting market conditions
func InitializeMarketCycle(difficulty string) *MarketCycle {
	// Start with normal market conditions
	// Harder difficulties start in worse market conditions
	
	switch difficulty {
	case "easy":
		return &MarketCycle{
			Name:           MarketBull,
			StartTurn:      1,
			Duration:       12, // 12 turns
			ValuationMult:  1.2,
			FundingEase:    1.3,
			ExitLikelihood: 1.2,
			Description:    "Strong economic growth and high investor confidence",
			Color:          "green",
		}
	case "hard", "expert":
		return &MarketCycle{
			Name:           MarketBear,
			StartTurn:      1,
			Duration:       8,
			ValuationMult:  0.8,
			FundingEase:    0.7,
			ExitLikelihood: 0.8,
			Description:    "Economic uncertainty and cautious investors",
			Color:          "red",
		}
	default: // medium
		return &MarketCycle{
			Name:           MarketNormal,
			StartTurn:      1,
			Duration:       10,
			ValuationMult:  1.0,
			FundingEase:    1.0,
			ExitLikelihood: 1.0,
			Description:    "Stable economic conditions",
			Color:          "white",
		}
	}
}

// AdvanceMarketCycle potentially changes the market cycle
func AdvanceMarketCycle(currentCycle *MarketCycle, turn int) *MarketCycle {
	// Check if current cycle has ended
	if turn < currentCycle.StartTurn+currentCycle.Duration {
		return currentCycle // Still in current cycle
	}
	
	// Generate new cycle
	return generateNextCycle(currentCycle, turn)
}

// generateNextCycle creates the next market cycle
func generateNextCycle(previousCycle *MarketCycle, turn int) *MarketCycle {
	// Cycles tend to revert to normal, with some randomness
	cycles := []struct {
		name       string
		prob       float64
		valuation  float64
		funding    float64
		exitChance float64
		duration   int
		desc       string
		color      string
	}{
		{MarketBull, 0.20, 1.3, 1.4, 1.3, 10, "Strong economic growth and high investor confidence", "green"},
		{MarketBoom, 0.10, 1.5, 1.8, 1.5, 6, "Exceptional market conditions and abundant capital", "green"},
		{MarketNormal, 0.40, 1.0, 1.0, 1.0, 12, "Stable economic conditions", "white"},
		{MarketBear, 0.20, 0.75, 0.7, 0.8, 8, "Economic uncertainty and cautious investors", "red"},
		{MarketRecession, 0.10, 0.6, 0.5, 0.6, 8, "Economic downturn and scarce capital", "red"},
	}
	
	// Adjust probabilities based on previous cycle (reversion to mean)
	roll := rand.Float64()
	cumProb := 0.0
	
	// If coming from extreme, more likely to normalize
	if previousCycle.Name == MarketBoom || previousCycle.Name == MarketRecession {
		// Force more normal conditions
		if roll < 0.6 {
			return &MarketCycle{
				Name:           MarketNormal,
				StartTurn:      turn,
				Duration:       12,
				ValuationMult:  1.0,
				FundingEase:    1.0,
				ExitLikelihood: 1.0,
				Description:    "Stable economic conditions",
				Color:          "white",
			}
		}
	}
	
	for _, cycle := range cycles {
		cumProb += cycle.prob
		if roll < cumProb {
			return &MarketCycle{
				Name:           cycle.name,
				StartTurn:      turn,
				Duration:       cycle.duration,
				ValuationMult:  cycle.valuation,
				FundingEase:    cycle.funding,
				ExitLikelihood: cycle.exitChance,
				Description:    cycle.desc,
				Color:          cycle.color,
			}
		}
	}
	
	// Fallback to normal
	return &MarketCycle{
		Name:           MarketNormal,
		StartTurn:      turn,
		Duration:       12,
		ValuationMult:  1.0,
		FundingEase:    1.0,
		ExitLikelihood: 1.0,
		Description:    "Stable economic conditions",
		Color:          "white",
	}
}

// GenerateEconomicEvent creates a random economic event
func GenerateEconomicEvent(turn int) *EconomicEvent {
	// 15% chance per turn to generate an event
	if rand.Float64() > 0.15 {
		return nil
	}
	
	events := []struct {
		name        string
		desc        string
		duration    int
		impact      EventImpact
		sectors     []string
		probability float64
	}{
		{
			name:     "Interest Rate Hike",
			desc:     "Central bank raises interest rates, making capital more expensive",
			duration: 6,
			impact: EventImpact{
				ValuationChange: -15.0,
				FundingChange:   -25.0,
				RevenueChange:   -5.0,
			},
			sectors:     []string{}, // Affects all
			probability: 0.15,
		},
		{
			name:     "Tech Boom",
			desc:     "Surge in technology sector investment and valuations",
			duration: 8,
			impact: EventImpact{
				ValuationChange: 25.0,
				FundingChange:   30.0,
				RevenueChange:   15.0,
			},
			sectors:     []string{"AI/ML", "SaaS", "Enterprise Software", "Cybersecurity"},
			probability: 0.12,
		},
		{
			name:     "Credit Crunch",
			desc:     "Banks tighten lending standards, reducing available capital",
			duration: 10,
			impact: EventImpact{
				ValuationChange: -20.0,
				FundingChange:   -40.0,
				RevenueChange:   -10.0,
			},
			sectors:     []string{},
			probability: 0.10,
		},
		{
			name:     "IPO Window Opens",
			desc:     "Public markets become receptive to new offerings",
			duration: 6,
			impact: EventImpact{
				ValuationChange: 20.0,
				FundingChange:   25.0,
			},
			sectors:     []string{},
			probability: 0.12,
		},
		{
			name:     "Healthcare Innovation Wave",
			desc:     "Breakthrough in healthcare drives sector enthusiasm",
			duration: 8,
			impact: EventImpact{
				ValuationChange: 30.0,
				FundingChange:   35.0,
				RevenueChange:   20.0,
			},
			sectors:     []string{"Healthcare", "Biotech"},
			probability: 0.10,
		},
		{
			name:     "Consumer Spending Surge",
			desc:     "Strong consumer confidence boosts retail and e-commerce",
			duration: 6,
			impact: EventImpact{
				ValuationChange: 20.0,
				RevenueChange:   25.0,
			},
			sectors:     []string{"E-commerce", "Consumer Products"},
			probability: 0.12,
		},
		{
			name:     "Regulatory Crackdown",
			desc:     "New regulations increase compliance costs and uncertainty",
			duration: 12,
			impact: EventImpact{
				ValuationChange: -15.0,
				FundingChange:   -20.0,
				RevenueChange:   -8.0,
			},
			sectors:     []string{"Fintech", "Crypto"},
			probability: 0.10,
		},
		{
			name:     "AI Investment Frenzy",
			desc:     "Artificial intelligence dominates investor attention",
			duration: 10,
			impact: EventImpact{
				ValuationChange: 40.0,
				FundingChange:   50.0,
				RevenueChange:   30.0,
			},
			sectors:     []string{"AI/ML"},
			probability: 0.10,
		},
		{
			name:     "Supply Chain Crisis",
			desc:     "Global supply chain disruptions impact operations",
			duration: 8,
			impact: EventImpact{
				RevenueChange: -15.0,
				ChurnChange:   10.0,
			},
			sectors:     []string{"E-commerce", "Hardware"},
			probability: 0.09,
		},
	}
	
	// Select event based on probabilities
	roll := rand.Float64()
	cumProb := 0.0
	
	for _, evt := range events {
		cumProb += evt.probability
		if roll < cumProb {
			return &EconomicEvent{
				Name:        evt.name,
				Description: evt.desc,
				Turn:        turn,
				Duration:    evt.duration,
				Impact:      evt.impact,
				Sectors:     evt.sectors,
				Active:      true,
			}
		}
	}
	
	return nil
}

// ApplyMarketEffects applies market cycle effects to a startup's valuation
func ApplyMarketEffects(startup *Startup, marketCycle *MarketCycle, events []EconomicEvent) *Startup {
	if startup == nil || marketCycle == nil {
		return startup
	}
	
	// Apply market cycle valuation multiplier
	startup.Valuation = int64(float64(startup.Valuation) * marketCycle.ValuationMult)
	
	// Apply economic event effects
	for _, event := range events {
		if !event.Active {
			continue
		}
		
		// Check if event affects this startup's sector
		affectsStartup := len(event.Sectors) == 0 // Empty means affects all
		if !affectsStartup {
			for _, sector := range event.Sectors {
				if sector == startup.Category {
					affectsStartup = true
					break
				}
			}
		}
		
		if affectsStartup {
			// Apply valuation change
			valuationChange := 1.0 + (event.Impact.ValuationChange / 100.0)
			startup.Valuation = int64(float64(startup.Valuation) * valuationChange)
		}
	}
	
	// Ensure valuation doesn't go below zero
	if startup.Valuation < 0 {
		startup.Valuation = 100000 // Minimum $100k valuation
	}
	
	return startup
}

// GetMarketSentiment returns a string describing current market conditions
func GetMarketSentiment(marketCycle *MarketCycle) string {
	if marketCycle == nil {
		return "Unknown"
	}
	
	switch marketCycle.Name {
	case MarketBoom:
		return "🚀 Extremely Bullish"
	case MarketBull:
		return "📈 Bullish"
	case MarketNormal:
		return "➡️ Neutral"
	case MarketBear:
		return "📉 Bearish"
	case MarketRecession:
		return "💔 Very Bearish"
	default:
		return "➡️ Neutral"
	}
}

// UpdateEconomicEvents checks which events are still active
func UpdateEconomicEvents(events []EconomicEvent, currentTurn int) []EconomicEvent {
	activeEvents := []EconomicEvent{}
	
	for _, event := range events {
		if event.Active && currentTurn < event.Turn+event.Duration {
			activeEvents = append(activeEvents, event)
		}
	}
	
	return activeEvents
}

// GetMarketConditionsSummary returns a formatted string of market conditions
func GetMarketConditionsSummary(marketCycle *MarketCycle, events []EconomicEvent) string {
	summary := fmt.Sprintf("%s: %s", marketCycle.Name, marketCycle.Description)
	
	if len(events) > 0 {
		summary += "\n\nActive Economic Events:"
		for _, event := range events {
			if event.Active {
				summary += fmt.Sprintf("\n  • %s: %s", event.Name, event.Description)
			}
		}
	}
	
	return summary
}

