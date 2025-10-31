package achievements

import (
	"time"
)

// Achievement represents an unlockable achievement
type Achievement struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Category    string
	Condition   func(*GameProgress) bool
	Unlocked    bool
	UnlockedAt  time.Time
	Points      int
	Rarity      string // Common, Rare, Epic, Legendary
}

// GameProgress tracks all player progress for achievement checking
type GameProgress struct {
	// Game results
	FinalNetWorth   int64
	ROI             float64
	SuccessfulExits int
	TurnsPlayed     int
	Difficulty      string
	
	// Portfolio details
	InvestmentCount int
	SectorsInvested []string
	
	// Career stats
	TotalGames      int
	TotalWins       int
	BestNetWorth    int64
	BestROI         float64
	TotalExits      int
	
	// Special conditions
	WonWithoutLosses bool
	AllSameCategory  bool
	OnlyOneInvestment bool
	WonInUnder60Turns bool
}

// AchievementCategory types
const (
	CategoryWealth      = "Wealth"
	CategoryPerformance = "Performance"
	CategoryStrategy    = "Strategy"
	CategoryCareer      = "Career"
	CategoryChallenge   = "Challenge"
	CategorySpecial     = "Special"
)

// Rarity levels
const (
	RarityCommon    = "Common"
	RarityRare      = "Rare"
	RarityEpic      = "Epic"
	RarityLegendary = "Legendary"
)

// All available achievements
var AllAchievements = []Achievement{
	// Wealth Achievements
	{
		ID:          "millionaire",
		Name:        "Millionaire",
		Description: "Reach a net worth of $1,000,000",
		Icon:        "??",
		Category:    CategoryWealth,
		Points:      10,
		Rarity:      RarityCommon,
		Condition: func(p *GameProgress) bool {
			return p.FinalNetWorth >= 1000000
		},
	},
	{
		ID:          "multimillionaire",
		Name:        "Multi-Millionaire",
		Description: "Reach a net worth of $5,000,000",
		Icon:        "??",
		Category:    CategoryWealth,
		Points:      25,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.FinalNetWorth >= 5000000
		},
	},
	{
		ID:          "decamillionaire",
		Name:        "Deca-Millionaire",
		Description: "Reach a net worth of $10,000,000",
		Icon:        "??",
		Category:    CategoryWealth,
		Points:      50,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.FinalNetWorth >= 10000000
		},
	},
	{
		ID:          "king_midas",
		Name:        "King Midas",
		Description: "Reach a net worth of $50,000,000",
		Icon:        "??",
		Category:    CategoryWealth,
		Points:      100,
		Rarity:      RarityLegendary,
		Condition: func(p *GameProgress) bool {
			return p.FinalNetWorth >= 50000000
		},
	},
	
	// Performance Achievements
	{
		ID:          "profitable",
		Name:        "In the Green",
		Description: "Achieve positive ROI",
		Icon:        "??",
		Category:    CategoryPerformance,
		Points:      5,
		Rarity:      RarityCommon,
		Condition: func(p *GameProgress) bool {
			return p.ROI > 0
		},
	},
	{
		ID:          "double_up",
		Name:        "Double Up",
		Description: "Achieve 100% ROI or better",
		Icon:        "??",
		Category:    CategoryPerformance,
		Points:      15,
		Rarity:      RarityCommon,
		Condition: func(p *GameProgress) bool {
			return p.ROI >= 100
		},
	},
	{
		ID:          "great_investor",
		Name:        "Great Investor",
		Description: "Achieve 200% ROI or better",
		Icon:        "?",
		Category:    CategoryPerformance,
		Points:      25,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.ROI >= 200
		},
	},
	{
		ID:          "elite_vc",
		Name:        "Elite VC",
		Description: "Achieve 500% ROI or better",
		Icon:        "??",
		Category:    CategoryPerformance,
		Points:      50,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.ROI >= 500
		},
	},
	{
		ID:          "unicorn_hunter",
		Name:        "Unicorn Hunter",
		Description: "Achieve 1000% ROI or better",
		Icon:        "??",
		Category:    CategoryPerformance,
		Points:      100,
		Rarity:      RarityLegendary,
		Condition: func(p *GameProgress) bool {
			return p.ROI >= 1000
		},
	},
	
	// Strategy Achievements
	{
		ID:          "diversified",
		Name:        "Diversified Portfolio",
		Description: "Invest in 5 or more companies",
		Icon:        "??",
		Category:    CategoryStrategy,
		Points:      10,
		Rarity:      RarityCommon,
		Condition: func(p *GameProgress) bool {
			return p.InvestmentCount >= 5
		},
	},
	{
		ID:          "all_in",
		Name:        "All In",
		Description: "Win with only one investment",
		Icon:        "??",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.OnlyOneInvestment && p.ROI > 0
		},
	},
	{
		ID:          "sector_specialist",
		Name:        "Sector Specialist",
		Description: "Win with all investments in the same category",
		Icon:        "??",
		Category:    CategoryStrategy,
		Points:      20,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.AllSameCategory && p.ROI > 0
		},
	},
	{
		ID:          "exit_master",
		Name:        "Exit Master",
		Description: "Achieve 3 or more successful exits (5x+) in one game",
		Icon:        "??",
		Category:    CategoryStrategy,
		Points:      25,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.SuccessfulExits >= 3
		},
	},
	{
		ID:          "perfect_portfolio",
		Name:        "Perfect Portfolio",
		Description: "Win without a single losing investment",
		Icon:        "?",
		Category:    CategoryStrategy,
		Points:      50,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.WonWithoutLosses && p.ROI > 0
		},
	},
	
	// Career Achievements
	{
		ID:          "getting_started",
		Name:        "Getting Started",
		Description: "Complete your first game",
		Icon:        "??",
		Category:    CategoryCareer,
		Points:      5,
		Rarity:      RarityCommon,
		Condition: func(p *GameProgress) bool {
			return p.TotalGames >= 1
		},
	},
	{
		ID:          "persistent",
		Name:        "Persistent",
		Description: "Play 10 games",
		Icon:        "??",
		Category:    CategoryCareer,
		Points:      15,
		Rarity:      RarityCommon,
		Condition: func(p *GameProgress) bool {
			return p.TotalGames >= 10
		},
	},
	{
		ID:          "veteran",
		Name:        "Veteran",
		Description: "Play 25 games",
		Icon:        "???",
		Category:    CategoryCareer,
		Points:      25,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.TotalGames >= 25
		},
	},
	{
		ID:          "master",
		Name:        "Master Investor",
		Description: "Play 50 games",
		Icon:        "??",
		Category:    CategoryCareer,
		Points:      50,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.TotalGames >= 50
		},
	},
	{
		ID:          "consistent",
		Name:        "Consistent Winner",
		Description: "Win 10 games with positive ROI",
		Icon:        "??",
		Category:    CategoryCareer,
		Points:      30,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.TotalWins >= 10
		},
	},
	
	// Challenge Achievements
	{
		ID:          "speed_runner",
		Name:        "Speed Runner",
		Description: "Win in under 60 turns",
		Icon:        "?",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.WonInUnder60Turns && p.ROI > 0
		},
	},
	{
		ID:          "easy_master",
		Name:        "Easy Master",
		Description: "Achieve 500%+ ROI on Easy difficulty",
		Icon:        "??",
		Category:    CategoryChallenge,
		Points:      20,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.Difficulty == "Easy" && p.ROI >= 500
		},
	},
	{
		ID:          "medium_master",
		Name:        "Medium Master",
		Description: "Achieve 500%+ ROI on Medium difficulty",
		Icon:        "?",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.Difficulty == "Medium" && p.ROI >= 500
		},
	},
	{
		ID:          "hard_master",
		Name:        "Hard Master",
		Description: "Achieve 500%+ ROI on Hard difficulty",
		Icon:        "??",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
		Condition: func(p *GameProgress) bool {
			return p.Difficulty == "Hard" && p.ROI >= 500
		},
	},
	{
		ID:          "expert_master",
		Name:        "Expert Master",
		Description: "Achieve 500%+ ROI on Expert difficulty",
		Icon:        "??",
		Category:    CategoryChallenge,
		Points:      100,
		Rarity:      RarityLegendary,
		Condition: func(p *GameProgress) bool {
			return p.Difficulty == "Expert" && p.ROI >= 500
		},
	},
	{
		ID:          "expert_survivor",
		Name:        "Expert Survivor",
		Description: "Simply complete an Expert game with positive ROI",
		Icon:        "???",
		Category:    CategoryChallenge,
		Points:      25,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.Difficulty == "Expert" && p.ROI > 0
		},
	},
	
	// Special Achievements
	{
		ID:          "lucky_seven",
		Name:        "Lucky Seven",
		Description: "Invest in exactly 7 companies and win",
		Icon:        "??",
		Category:    CategorySpecial,
		Points:      15,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.InvestmentCount == 7 && p.ROI > 0
		},
	},
	{
		ID:          "minimalist",
		Name:        "Minimalist",
		Description: "Win with exactly 2 investments",
		Icon:        "??",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return p.InvestmentCount == 2 && p.ROI > 0
		},
	},
	{
		ID:          "sector_rotation",
		Name:        "Sector Rotation",
		Description: "Invest in 5+ different sectors",
		Icon:        "??",
		Category:    CategorySpecial,
		Points:      15,
		Rarity:      RarityRare,
		Condition: func(p *GameProgress) bool {
			return len(p.SectorsInvested) >= 5
		},
	},
}

// CheckAchievements returns newly unlocked achievements
func CheckAchievements(progress *GameProgress, unlockedIDs []string) []Achievement {
	var newlyUnlocked []Achievement
	
	// Create map of already unlocked achievements
	unlockedMap := make(map[string]bool)
	for _, id := range unlockedIDs {
		unlockedMap[id] = true
	}
	
	// Check each achievement
	for _, achievement := range AllAchievements {
		// Skip if already unlocked
		if unlockedMap[achievement.ID] {
			continue
		}
		
		// Check if condition is met
		if achievement.Condition(progress) {
			achievement.Unlocked = true
			achievement.UnlockedAt = time.Now()
			newlyUnlocked = append(newlyUnlocked, achievement)
		}
	}
	
	return newlyUnlocked
}

// GetAchievementByID retrieves a specific achievement
func GetAchievementByID(id string) *Achievement {
	for _, achievement := range AllAchievements {
		if achievement.ID == id {
			return &achievement
		}
	}
	return nil
}

// CalculateCareerLevel calculates player level based on achievement points
func CalculateCareerLevel(totalPoints int) (level int, title string) {
	if totalPoints >= 1000 {
		return 10, "Legendary Investor"
	} else if totalPoints >= 750 {
		return 9, "Master Investor"
	} else if totalPoints >= 500 {
		return 8, "Elite VC"
	} else if totalPoints >= 350 {
		return 7, "Senior Partner"
	} else if totalPoints >= 250 {
		return 6, "Partner"
	} else if totalPoints >= 150 {
		return 5, "Principal"
	} else if totalPoints >= 100 {
		return 4, "Senior Associate"
	} else if totalPoints >= 50 {
		return 3, "Associate"
	} else if totalPoints >= 25 {
		return 2, "Junior Associate"
	} else if totalPoints >= 10 {
		return 1, "Analyst"
	}
	return 0, "Intern"
}

// GetAchievementsByCategory returns achievements grouped by category
func GetAchievementsByCategory(category string) []Achievement {
	var filtered []Achievement
	for _, achievement := range AllAchievements {
		if achievement.Category == category {
			filtered = append(filtered, achievement)
		}
	}
	return filtered
}

// GetAchievementProgress returns completion percentage
func GetAchievementProgress(unlockedCount int) float64 {
	return (float64(unlockedCount) / float64(len(AllAchievements))) * 100.0
}
