package game

import (
	"math/rand"
)

// Founder name pools inspired by real tech founders
var founderFirstNames = []string{
	// Male names
	"Mark", "Jeff", "Elon", "Larry", "Sergey", "Bill", "Steve", "Jack", "Reed",
	"Brian", "Travis", "Drew", "Nathan", "Patrick", "Daniel", "Kevin", "Aaron",
	"Ben", "Paul", "Peter", "Marc", "Reid", "Sam", "Andrew", "Chris", "David",
	"Eric", "Max", "Ryan", "Sean", "Tony", "Mike", "Tom", "Travis", "Garrett",
	// Female names
	"Sara", "Whitney", "Jessica", "Jennifer", "Sarah", "Emily", "Julia", "Anne",
	"Katrina", "Diane", "Marissa", "Sheryl", "Susan", "Ginni", "Meg", "Safra",
	"Leila", "Reshma", "Melanie", "Stacy", "Aileen", "Ann", "Megan", "Whitney",
}

var founderLastNames = []string{
	"Zuckerberg", "Bezos", "Musk", "Page", "Brin", "Gates", "Jobs", "Dorsey",
	"Hastings", "Chesky", "Kalanick", "Houston", "Blecharczyk", "Collison",
	"Zuckerberg", "Systrom", "Chen", "Pincus", "Hoffman", "Thiel", "Andreessen",
	"Chen", "Altman", "Mason", "Dixon", "Kim", "Lin", "Patel", "Gupta", "Lee",
	"Wang", "Chen", "Li", "Kumar", "Singh", "Shah", "Nguyen", "Garcia", "Lopez",
	"Martinez", "Rodriguez", "Hernandez", "Brown", "Johnson", "Williams", "Jones",
	"Taylor", "Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Thompson",
	"Garcia", "Martinez", "Robinson", "Clark", "Rodriguez", "Lewis", "Walker",
}

// GenerateFounderName creates a realistic founder name
func GenerateFounderName() string {
	firstName := founderFirstNames[rand.Intn(len(founderFirstNames))]
	lastName := founderLastNames[rand.Intn(len(founderLastNames))]
	return firstName + " " + lastName
}

// CalculateInitialRelationship determines starting relationship score based on investment terms
func CalculateInitialRelationship(terms InvestmentTerms, hasDueDiligence bool, amount int64) float64 {
	// Base relationship: 50-70 range
	relationship := 50.0 + float64(rand.Intn(21)) // 50-70

	// Founder-friendly terms improve relationship
	if terms.Type == "Common Stock" {
		relationship += 10.0 // Founders love this
	} else if terms.Type == "SAFE" || terms.Type == "SAFE (Capped)" {
		relationship += 5.0 // SAFE is founder-friendly
	} else if terms.Type == "Preferred Stock (2x Liquidation)" {
		relationship -= 5.0 // Aggressive terms hurt relationship
	}

	// Due diligence shows seriousness and professionalism
	if hasDueDiligence {
		relationship += 8.0
	}

	// Larger investments show confidence
	if amount >= 100000 {
		relationship += 5.0
	} else if amount >= 50000 {
		relationship += 3.0
	}

	// Cap at reasonable range
	if relationship > 85.0 {
		relationship = 85.0
	}
	if relationship < 40.0 {
		relationship = 40.0
	}

	return relationship
}

// GetRelationshipLevel returns a descriptive level for a relationship score
func GetRelationshipLevel(score float64) string {
	if score >= 90 {
		return "Exceptional"
	} else if score >= 80 {
		return "Strong"
	} else if score >= 70 {
		return "Good"
	} else if score >= 60 {
		return "Cordial"
	} else if score >= 50 {
		return "Professional"
	} else if score >= 40 {
		return "Strained"
	} else if score >= 30 {
		return "Tense"
	}
	return "Hostile"
}

// GetRelationshipEmoji returns an emoji representing relationship health
func GetRelationshipEmoji(score float64) string {
	if score >= 80 {
		return "😊" // High relationship
	} else if score >= 50 {
		return "😐" // Medium relationship
	}
	return "😟" // Low relationship
}

// RelationshipEvent represents an event affecting founder relationships
type RelationshipEvent struct {
	CompanyName    string
	FounderName    string
	EventType      string // "positive", "negative", "neutral"
	Description    string
	ScoreChange    float64
	RequiresAction bool // If true, player can take action to mitigate
}

// GenerateRelationshipEvent creates a random relationship event
func GenerateRelationshipEvent(inv *Investment, currentTurn int) *RelationshipEvent {
	// 10% chance per turn for a relationship event
	if rand.Float64() > 0.10 {
		return nil
	}

	events := []struct {
		eventType   string
		description string
		scoreChange float64
	}{
		// Positive events
		{"positive", "reached out to thank you for your support and strategic advice", 5.0},
		{"positive", "mentioned you favorably in a podcast interview", 3.0},
		{"positive", "implemented your suggestion, leading to a successful outcome", 7.0},
		{"positive", "invited you to a company milestone celebration", 4.0},
		{"positive", "publicly credited your firm for their success", 6.0},

		// Negative events
		{"negative", "expressed frustration about lack of engagement from your firm", -5.0},
		{"negative", "disagreed strongly with board direction you supported", -4.0},
		{"negative", "feels micromanaged by board demands", -6.0},
		{"negative", "is upset about recruiting competition from your other portfolio companies", -5.0},
		{"negative", "heard criticism of your firm from other founders", -3.0},

		// Neutral but important events
		{"neutral", "asked for help with a strategic decision", 0.0},
		{"neutral", "requested an introduction to potential customer", 0.0},
		{"neutral", "wants advice on hiring a key executive", 0.0},
	}

	event := events[rand.Intn(len(events))]

	return &RelationshipEvent{
		CompanyName:    inv.CompanyName,
		FounderName:    inv.FounderName,
		EventType:      event.eventType,
		Description:    inv.FounderName + " " + event.description,
		ScoreChange:    event.scoreChange,
		RequiresAction: event.eventType == "neutral", // Neutral events are opportunities
	}
}

// ApplyRelationshipChange updates relationship score with bounds checking
func ApplyRelationshipChange(currentScore float64, change float64) float64 {
	newScore := currentScore + change

	// Bounds: 0-100
	if newScore > 100.0 {
		return 100.0
	}
	if newScore < 0.0 {
		return 0.0
	}

	return newScore
}

// GetRelationshipImpactOnExit calculates exit value multiplier based on relationship
func GetRelationshipImpactOnExit(relationshipScore float64) float64 {
	// High relationships (80+) can boost exit by up to 15%
	// Low relationships (<50) can hurt exit by up to 10%

	if relationshipScore >= 80 {
		// 1.00 to 1.15 multiplier
		bonus := (relationshipScore - 80.0) / 20.0 * 0.15
		return 1.0 + bonus
	} else if relationshipScore < 50 {
		// 0.90 to 1.00 multiplier
		penalty := (50.0 - relationshipScore) / 50.0 * 0.10
		return 1.0 - penalty
	}

	// Neutral impact for medium relationships
	return 1.0
}

// CanBeFiredFromBoard checks if poor relationship leads to board removal
func CanBeFiredFromBoard(relationshipScore float64) bool {
	// Very poor relationships (<30) have chance of board removal
	if relationshipScore < 30 {
		return rand.Float64() < 0.05 // 5% chance per turn when very low
	}
	return false
}

// GetFounderReferralBonus returns deal flow bonus from good founder relationships
func GetFounderReferralBonus(avgRelationshipScore float64) float64 {
	// High average relationships with founders lead to better deal flow
	// 80+ average = 10% better deal quality
	// 70+ average = 5% better deal quality

	if avgRelationshipScore >= 80 {
		return 0.10
	} else if avgRelationshipScore >= 70 {
		return 0.05
	}

	return 0.0
}
