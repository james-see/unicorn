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
	Points      int
	Rarity      string
	Hidden      bool
}

// PlayerAchievement tracks when a player unlocked an achievement
type PlayerAchievement struct {
	AchievementID string
	UnlockedAt    time.Time
}

// GameStats contains all stats needed for achievement checking
type GameStats struct {
	// Game mode
	GameMode string // "vc" or "founder"
	
	// Game results
	FinalNetWorth   int64
	ROI             float64
	SuccessfulExits int
	TurnsPlayed     int
	Difficulty      string
	
	// Portfolio details (VC mode)
	InvestmentCount int
	SectorsInvested []string
	TotalInvested   int64
	RiskScores      []float64 // Risk scores of all investments (for achievement tracking)
	
	// Performance (VC mode)
	PositiveInvestments int
	NegativeInvestments int
	BestROI             float64
	WorstROI            float64
	
	// Founder mode stats
	FinalMRR              int64
	FinalValuation        int64
	FinalEquity           float64
	Customers             int
	FundingRoundsRaised   int
	TotalFundingRaised    int64
	HasExited             bool
	ExitType              string // "ipo", "acquisition", "secondary"
	ExitValuation         int64
	MonthsToProfitability int
	
	// Career stats
	TotalGames      int
	TotalWins       int
	WinStreak       int
	BestNetWorth    int64
	TotalExits      int
}

// Achievement categories
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
var AllAchievements = map[string]Achievement{
	// Wealth Achievements
	"first_profit": {
		ID:          "first_profit",
		Name:        "First Profit",
		Description: "Make your first dollar of profit",
		Icon:        "$",
		Category:    CategoryWealth,
		Points:      5,
		Rarity:      RarityCommon,
	},
	"millionaire": {
		ID:          "millionaire",
		Name:        "Millionaire",
		Description: "Reach $1,000,000 net worth",
		Icon:        "💰",
		Category:    CategoryWealth,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"multi_millionaire": {
		ID:          "multi_millionaire",
		Name:        "Multi-Millionaire",
		Description: "Reach $5,000,000 net worth",
		Icon:        "💵",
		Category:    CategoryWealth,
		Points:      25,
		Rarity:      RarityRare,
	},
	"deca_millionaire": {
		ID:          "deca_millionaire",
		Name:        "Deca-Millionaire",
		Description: "Reach $10,000,000 net worth",
		Icon:        "🏦",
		Category:    CategoryWealth,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"mega_rich": {
		ID:          "mega_rich",
		Name:        "Mega Rich",
		Description: "Reach $50,000,000 net worth",
		Icon:        "👑",
		Category:    CategoryWealth,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	
	// Performance Achievements
	"break_even": {
		ID:          "break_even",
		Name:        "Break Even",
		Description: "Achieve 0% or better ROI",
		Icon:        "=",
		Category:    CategoryPerformance,
		Points:      5,
		Rarity:      RarityCommon,
	},
	"double_up": {
		ID:          "double_up",
		Name:        "Double Up",
		Description: "Achieve 100%+ ROI",
		Icon:        "📈",
		Category:    CategoryPerformance,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"great_investor": {
		ID:          "great_investor",
		Name:        "Great Investor",
		Description: "Achieve 200%+ ROI",
		Icon:        "⭐",
		Category:    CategoryPerformance,
		Points:      25,
		Rarity:      RarityRare,
	},
	"elite_vc": {
		ID:          "elite_vc",
		Name:        "Elite VC",
		Description: "Achieve 500%+ ROI",
		Icon:        "🏆",
		Category:    CategoryPerformance,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"unicorn_hunter": {
		ID:          "unicorn_hunter",
		Name:        "Unicorn Hunter",
		Description: "Achieve 1000%+ ROI",
		Icon:        "🦄",
		Category:    CategoryPerformance,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	
	// Strategy Achievements
	"diversified": {
		ID:          "diversified",
		Name:        "Diversified",
		Description: "Invest in 5+ companies",
		Icon:        "📊",
		Category:    CategoryStrategy,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"sector_master": {
		ID:          "sector_master",
		Name:        "Sector Master",
		Description: "Invest in 5+ different sectors",
		Icon:        "🏢",
		Category:    CategoryStrategy,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"all_in": {
		ID:          "all_in",
		Name:        "All In",
		Description: "Win with only 1 investment",
		Icon:        "🎲",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityEpic,
	},
	"sector_specialist": {
		ID:          "sector_specialist",
		Name:        "Sector Specialist",
		Description: "Win with all investments in same sector",
		Icon:        "🎯",
		Category:    CategoryStrategy,
		Points:      20,
		Rarity:      RarityRare,
	},
	"exit_master": {
		ID:          "exit_master",
		Name:        "Exit Master",
		Description: "3+ successful exits (5x) in one game",
		Icon:        "🚀",
		Category:    CategoryStrategy,
		Points:      25,
		Rarity:      RarityRare,
	},
	"perfect_portfolio": {
		ID:          "perfect_portfolio",
		Name:        "Perfect Portfolio",
		Description: "Win without any losing investments",
		Icon:        "✨",
		Category:    CategoryStrategy,
		Points:      50,
		Rarity:      RarityEpic,
	},
	
	// Career Achievements
	"first_game": {
		ID:          "first_game",
		Name:        "First Steps",
		Description: "Complete your first game",
		Icon:        "👣",
		Category:    CategoryCareer,
		Points:      5,
		Rarity:      RarityCommon,
	},
	"persistent": {
		ID:          "persistent",
		Name:        "Persistent",
		Description: "Play 10 games",
		Icon:        "💪",
		Category:    CategoryCareer,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"veteran": {
		ID:          "veteran",
		Name:        "Veteran",
		Description: "Play 25 games",
		Icon:        "🎖️",
		Category:    CategoryCareer,
		Points:      25,
		Rarity:      RarityRare,
	},
	"master_investor": {
		ID:          "master_investor",
		Name:        "Master Investor",
		Description: "Play 50 games",
		Icon:        "👑",
		Category:    CategoryCareer,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"win_streak_3": {
		ID:          "win_streak_3",
		Name:        "Hot Streak",
		Description: "Win 3 games in a row",
		Icon:        "🔥",
		Category:    CategoryCareer,
		Points:      20,
		Rarity:      RarityRare,
	},
	"win_streak_5": {
		ID:          "win_streak_5",
		Name:        "On Fire",
		Description: "Win 5 games in a row",
		Icon:        "⚡",
		Category:    CategoryCareer,
		Points:      40,
		Rarity:      RarityEpic,
	},
	
	// Challenge Achievements
	"easy_win": {
		ID:          "easy_win",
		Name:        "Easy Money",
		Description: "Win on Easy difficulty",
		Icon:        "✅",
		Category:    CategoryChallenge,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"medium_win": {
		ID:          "medium_win",
		Name:        "Rising Star",
		Description: "Win on Medium difficulty",
		Icon:        "⭐",
		Category:    CategoryChallenge,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"hard_win": {
		ID:          "hard_win",
		Name:        "Battle Tested",
		Description: "Win on Hard difficulty",
		Icon:        "🛡️",
		Category:    CategoryChallenge,
		Points:      25,
		Rarity:      RarityRare,
	},
	"expert_win": {
		ID:          "expert_win",
		Name:        "Expert Survivor",
		Description: "Win on Expert difficulty",
		Icon:        "💀",
		Category:    CategoryChallenge,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"easy_master": {
		ID:          "easy_master",
		Name:        "Easy Domination",
		Description: "500%+ ROI on Easy",
		Icon:        "🥇",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityRare,
	},
	"expert_master": {
		ID:          "expert_master",
		Name:        "Expert Legend",
		Description: "500%+ ROI on Expert",
		Icon:        "🌟",
		Category:    CategoryChallenge,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	"speed_runner": {
		ID:          "speed_runner",
		Name:        "Speed Runner",
		Description: "Win in under 60 turns",
		Icon:        "🏃",
		Category:    CategoryChallenge,
		Points:      30,
		Rarity:      RarityRare,
	},
	
	// Special Achievements
	"lucky_seven": {
		ID:          "lucky_seven",
		Name:        "Lucky Seven",
		Description: "Invest in exactly 7 companies and win",
		Icon:        "🍀",
		Category:    CategorySpecial,
		Points:      15,
		Rarity:      RarityRare,
	},
	"minimalist": {
		ID:          "minimalist",
		Name:        "Minimalist",
		Description: "Win with exactly 2 investments",
		Icon:        "🎯",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
	},
	"tech_enthusiast": {
		ID:          "tech_enthusiast",
		Name:        "Tech Enthusiast",
		Description: "Only invest in tech sectors and win",
		Icon:        "💻",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
	},
	"clean_investor": {
		ID:          "clean_investor",
		Name:        "Clean Investor",
		Description: "Only invest in CleanTech/AgriTech and win",
		Icon:        "🌱",
		Category:    CategorySpecial,
		Points:      20,
		Rarity:      RarityRare,
	},
	"risk_taker": {
		ID:          "risk_taker",
		Name:        "Risk Taker",
		Description: "Win with only high-risk companies",
		Icon:        "🎲",
		Category:    CategorySpecial,
		Points:      35,
		Rarity:      RarityEpic,
		Hidden:      true,
	},
	"cautious_investor": {
		ID:          "cautious_investor",
		Name:        "Cautious Investor",
		Description: "Win with only low-risk companies",
		Icon:        "🛡️",
		Category:    CategorySpecial,
		Points:      25,
		Rarity:      RarityRare,
	},
	
	// Founder Mode Achievements
	"first_revenue": {
		ID:          "first_revenue",
		Name:        "First Revenue",
		Description: "Generate your first $1,000 MRR",
		Icon:        "💵",
		Category:    CategoryWealth,
		Points:      10,
		Rarity:      RarityCommon,
	},
	"profitable": {
		ID:          "profitable",
		Name:        "Profitable",
		Description: "Reach profitability (positive cash flow)",
		Icon:        "📈",
		Category:    CategoryPerformance,
		Points:      25,
		Rarity:      RarityRare,
	},
	"100k_mrr": {
		ID:          "100k_mrr",
		Name:        "$100K MRR Club",
		Description: "Reach $100,000 monthly recurring revenue",
		Icon:        "🎯",
		Category:    CategoryWealth,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"1m_mrr": {
		ID:          "1m_mrr",
		Name:        "Unicorn MRR",
		Description: "Reach $1,000,000 monthly recurring revenue",
		Icon:        "🦄",
		Category:    CategoryWealth,
		Points:      100,
		Rarity:      RarityLegendary,
	},
	"seed_raised": {
		ID:          "seed_raised",
		Name:        "Seed Raiser",
		Description: "Raise your first funding round",
		Icon:        "🌱",
		Category:    CategoryStrategy,
		Points:      15,
		Rarity:      RarityCommon,
	},
	"series_a": {
		ID:          "series_a",
		Name:        "Series A Graduate",
		Description: "Raise Series A funding",
		Icon:        "🚀",
		Category:    CategoryStrategy,
		Points:      30,
		Rarity:      RarityRare,
	},
	"ipo_exit": {
		ID:          "ipo_exit",
		Name:        "Public Debut",
		Description: "Take your company public via IPO",
		Icon:        "🏛️",
		Category:    CategorySpecial,
		Points:      75,
		Rarity:      RarityLegendary,
	},
	"acquired": {
		ID:          "acquired",
		Name:        "Acquired",
		Description: "Get acquired by another company",
		Icon:        "🤝",
		Category:    CategorySpecial,
		Points:      50,
		Rarity:      RarityEpic,
	},
	"10000_customers": {
		ID:          "10000_customers",
		Name:        "10K Customers",
		Description: "Reach 10,000 customers",
		Icon:        "👥",
		Category:    CategoryPerformance,
		Points:      40,
		Rarity:      RarityEpic,
	},
	"bootstrapped": {
		ID:          "bootstrapped",
		Name:        "Bootstrapped",
		Description: "Reach $100K MRR without raising funding",
		Icon:        "💪",
		Category:    CategoryChallenge,
		Points:      60,
		Rarity:      RarityLegendary,
		Hidden:      true,
	},
}

// CheckAchievements checks which achievements were earned this game
func CheckAchievements(stats GameStats, previouslyUnlocked []string) []Achievement {
	unlocked := make(map[string]bool)
	for _, id := range previouslyUnlocked {
		unlocked[id] = true
	}
	
	var newAchievements []Achievement
	
	for id, achievement := range AllAchievements {
		if unlocked[id] {
			continue
		}
		
		if checkAchievement(id, stats) {
			newAchievements = append(newAchievements, achievement)
		}
	}
	
	return newAchievements
}

func checkAchievement(id string, stats GameStats) bool {
	won := stats.ROI > 0
	
	switch id {
	// Wealth
	case "first_profit":
		return stats.FinalNetWorth > stats.TotalInvested
	case "millionaire":
		return stats.FinalNetWorth >= 1000000
	case "multi_millionaire":
		return stats.FinalNetWorth >= 5000000
	case "deca_millionaire":
		return stats.FinalNetWorth >= 10000000
	case "mega_rich":
		return stats.FinalNetWorth >= 50000000
		
	// Performance
	case "break_even":
		return stats.ROI >= 0
	case "double_up":
		return stats.ROI >= 100
	case "great_investor":
		return stats.ROI >= 200
	case "elite_vc":
		return stats.ROI >= 500
	case "unicorn_hunter":
		return stats.ROI >= 1000
		
	// Strategy
	case "diversified":
		return stats.InvestmentCount >= 5
	case "sector_master":
		return len(stats.SectorsInvested) >= 5
	case "all_in":
		return stats.InvestmentCount == 1 && won
	case "sector_specialist":
		return len(stats.SectorsInvested) == 1 && stats.InvestmentCount > 1 && won
	case "exit_master":
		return stats.SuccessfulExits >= 3
	case "perfect_portfolio":
		return stats.NegativeInvestments == 0 && stats.InvestmentCount > 0 && won
		
	// Career
	case "first_game":
		return stats.TotalGames >= 1
	case "persistent":
		return stats.TotalGames >= 10
	case "veteran":
		return stats.TotalGames >= 25
	case "master_investor":
		return stats.TotalGames >= 50
	case "win_streak_3":
		return stats.WinStreak >= 3
	case "win_streak_5":
		return stats.WinStreak >= 5
		
	// Challenge
	case "easy_win":
		return stats.Difficulty == "Easy" && won
	case "medium_win":
		return stats.Difficulty == "Medium" && won
	case "hard_win":
		return stats.Difficulty == "Hard" && won
	case "expert_win":
		return stats.Difficulty == "Expert" && won
	case "easy_master":
		return stats.Difficulty == "Easy" && stats.ROI >= 500
	case "expert_master":
		return stats.Difficulty == "Expert" && stats.ROI >= 500
	case "speed_runner":
		return stats.TurnsPlayed < 60 && won
		
	// Special
	case "lucky_seven":
		return stats.InvestmentCount == 7 && won
	case "minimalist":
		return stats.InvestmentCount == 2 && won
	case "tech_enthusiast":
		// Check if all sectors are tech-related
		techSectors := map[string]bool{
			"CloudTech": true, "SaaS": true, "DeepTech": true,
			"FinTech": true, "HealthTech": true, "EdTech": true,
			"LegalTech": true, "Gaming": true, "Security": true,
		}
		for _, sector := range stats.SectorsInvested {
			if !techSectors[sector] {
				return false
			}
		}
		return len(stats.SectorsInvested) > 0 && won
	case "clean_investor":
		// Check if only CleanTech/AgriTech
		for _, sector := range stats.SectorsInvested {
			if sector != "CleanTech" && sector != "AgriTech" {
				return false
			}
		}
		return len(stats.SectorsInvested) > 0 && won
	case "risk_taker":
		// Win with only high-risk companies (risk score > 0.6)
		if stats.InvestmentCount == 0 || !won {
			return false
		}
		for _, risk := range stats.RiskScores {
			if risk <= 0.6 {
				return false // Not all high-risk
			}
		}
		return true
	case "cautious_investor":
		// Win with only low-risk companies (risk score < 0.3)
		if stats.InvestmentCount == 0 || !won {
			return false
		}
		for _, risk := range stats.RiskScores {
			if risk >= 0.3 {
				return false // Not all low-risk
			}
		}
		return true
		
	// Founder Mode Achievements
	case "first_revenue":
		return stats.GameMode == "founder" && stats.FinalMRR >= 1000
	case "profitable":
		return stats.GameMode == "founder" && stats.MonthsToProfitability > 0
	case "100k_mrr":
		return stats.GameMode == "founder" && stats.FinalMRR >= 100000
	case "1m_mrr":
		return stats.GameMode == "founder" && stats.FinalMRR >= 1000000
	case "seed_raised":
		return stats.GameMode == "founder" && stats.FundingRoundsRaised >= 1
	case "series_a":
		return stats.GameMode == "founder" && stats.FundingRoundsRaised >= 2 // Assuming Seed + Series A
	case "ipo_exit":
		return stats.GameMode == "founder" && stats.HasExited && stats.ExitType == "ipo"
	case "acquired":
		return stats.GameMode == "founder" && stats.HasExited && stats.ExitType == "acquisition"
	case "10000_customers":
		return stats.GameMode == "founder" && stats.Customers >= 10000
	case "bootstrapped":
		return stats.GameMode == "founder" && stats.FinalMRR >= 100000 && stats.FundingRoundsRaised == 0
	}
	
	return false
}

// CalculateCareerLevel calculates player level based on achievement points
func CalculateCareerLevel(totalPoints int) (level int, title string, nextLevelPoints int) {
	levels := []struct {
		points int
		level  int
		title  string
	}{
		{0, 0, "Intern"},
		{25, 1, "Analyst"},
		{75, 2, "Associate"},
		{150, 3, "Senior Associate"},
		{250, 4, "Principal"},
		{400, 5, "Partner"},
		{600, 6, "Senior Partner"},
		{850, 7, "Managing Partner"},
		{1150, 8, "Elite VC"},
		{1500, 9, "Master Investor"},
		{2000, 10, "Legendary Investor"},
	}
	
	for i := len(levels) - 1; i >= 0; i-- {
		if totalPoints >= levels[i].points {
			nextLevel := 2001 // Max
			if i < len(levels)-1 {
				nextLevel = levels[i+1].points
			}
			return levels[i].level, levels[i].title, nextLevel
		}
	}
	
	return 0, "Intern", 25
}

// GetAchievementsByCategory returns achievements for a category
func GetAchievementsByCategory(category string) []Achievement {
	var achievements []Achievement
	for _, ach := range AllAchievements {
		if ach.Category == category && !ach.Hidden {
			achievements = append(achievements, ach)
		}
	}
	return achievements
}

// GetRarityColor returns color code for rarity
func GetRarityColor(rarity string) int {
	switch rarity {
	case RarityCommon:
		return 37 // White
	case RarityRare:
		return 36 // Cyan
	case RarityEpic:
		return 35 // Magenta
	case RarityLegendary:
		return 33 // Yellow
	default:
		return 37
	}
}
