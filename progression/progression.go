package progression

import (
	"github.com/jamesacampbell/unicorn/achievements"
	"github.com/jamesacampbell/unicorn/database"
)

// XP reward constants
const (
	XPGameComplete      = 100
	XPPositiveROI       = 50
	XPSuccessfulExit    = 200
	XPAchievementBase   = 10  // Base XP, multiplied by achievement points
	XPDifficultyEasy    = 0
	XPDifficultyMedium  = 50
	XPDifficultyHard    = 100
	XPDifficultyExpert  = 200
)

// LevelInfo contains information about a specific level
type LevelInfo struct {
	Level       int
	XPRequired  int
	TotalXP     int
	Unlocks     []string
	Title       string
}

// CalculateXPReward calculates total XP earned from a game session
func CalculateXPReward(gameStats *achievements.GameStats, unlockedAchievements []string) int {
	totalXP := XPGameComplete
	
	// Bonus for positive ROI
	if gameStats.ROI > 0 {
		totalXP += XPPositiveROI
	}
	
	// Bonus for successful exits
	if gameStats.SuccessfulExits > 0 {
		totalXP += XPSuccessfulExit * gameStats.SuccessfulExits
	}
	
	// Difficulty bonus
	switch gameStats.Difficulty {
	case "easy":
		totalXP += XPDifficultyEasy
	case "medium":
		totalXP += XPDifficultyMedium
	case "hard":
		totalXP += XPDifficultyHard
	case "expert":
		totalXP += XPDifficultyExpert
	}
	
	// Achievement bonuses
	for _, achvID := range unlockedAchievements {
		if achv, exists := achievements.AllAchievements[achvID]; exists {
			totalXP += achv.Points * XPAchievementBase
		}
	}
	
	// Founder mode specific bonuses
	if gameStats.GameMode == "founder" {
		if gameStats.HasExited && gameStats.ExitType == "ipo" {
			totalXP += 500 // Big bonus for IPO
		} else if gameStats.HasExited && gameStats.ExitType == "acquisition" {
			totalXP += 300 // Good bonus for acquisition
		}
		
		// Bonus for reaching profitability
		if gameStats.MonthsToProfitability > 0 && gameStats.MonthsToProfitability <= 24 {
			totalXP += 100
		}
	}
	
	return totalXP
}

// GetLevelInfo returns detailed information about a level
func GetLevelInfo(level int) LevelInfo {
	return LevelInfo{
		Level:      level,
		XPRequired: database.GetLevelRequirement(level),
		TotalXP:    getTotalXPForLevel(level),
		Unlocks:    GetUnlockablesAtLevel(level),
		Title:      getLevelTitle(level),
	}
}

// getTotalXPForLevel calculates total XP needed to reach a level (sum of all previous levels)
func getTotalXPForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	
	total := 0
	for i := 2; i <= level; i++ {
		total += database.GetLevelRequirement(i)
	}
	return total
}

// getLevelTitle returns a title/rank for a level
func getLevelTitle(level int) string {
	switch {
	case level < 5:
		return "Novice Investor"
	case level < 10:
		return "Angel Investor"
	case level < 15:
		return "Seed Investor"
	case level < 20:
		return "Venture Capitalist"
	case level < 25:
		return "Senior VC"
	case level < 30:
		return "Managing Partner"
	case level < 40:
		return "Super Angel"
	case level < 50:
		return "Legendary Investor"
	default:
		return "Titan of Industry"
	}
}

// GetUnlockablesAtLevel returns what unlocks at a specific level
func GetUnlockablesAtLevel(level int) []string {
	unlocks := []string{}
	
	switch level {
	case 1:
		unlocks = append(unlocks, "Easy & Medium difficulty available")
	case 2:
		unlocks = append(unlocks, "🤝 Investor Syndicates unlocked")
	case 5:
		unlocks = append(unlocks, "🔓 Hard difficulty unlocked")
	case 10:
		unlocks = append(unlocks, "🔓 Expert difficulty unlocked")
		unlocks = append(unlocks, "📊 Advanced analytics unlocked")
	case 15:
		unlocks = append(unlocks, "🔓 Nightmare mode unlocked")
		unlocks = append(unlocks, "🎯 Achievement chains visible")
	case 20:
		unlocks = append(unlocks, "🔓 Special game modifiers unlocked")
		unlocks = append(unlocks, "💎 Prestige system preview")
	case 25:
		unlocks = append(unlocks, "🔓 Prestige system unlocked")
		unlocks = append(unlocks, "🏆 Legendary achievements available")
	case 30:
		unlocks = append(unlocks, "🔓 Master difficulty unlocked")
	case 50:
		unlocks = append(unlocks, "👑 Titan title unlocked")
		unlocks = append(unlocks, "🌟 All content unlocked")
	}
	
	return unlocks
}

// IsLevelLocked checks if a difficulty/feature is locked based on level
func IsLevelLocked(playerLevel int, feature string) bool {
	requiredLevel := GetRequiredLevelFor(feature)
	return playerLevel < requiredLevel
}

// GetRequiredLevelFor returns the required level for a feature
func GetRequiredLevelFor(feature string) int {
	levelRequirements := map[string]int{
		"easy":           1,
		"medium":         1,
		"hard":           5,
		"expert":         10,
		"nightmare":      15,
		"master":         30,
		"game_modifiers": 20,
		"prestige":       25,
		"analytics":      10,
		"achievement_chains": 15,
		"syndicate_deals": 2,
	}
	
	if level, exists := levelRequirements[feature]; exists {
		return level
	}
	return 1 // Default to level 1 (unlocked)
}

// FormatXPBar creates a visual XP progress bar
func FormatXPBar(current, required int, width int) string {
	if required == 0 {
		return "[MAX LEVEL]"
	}
	
	progress := float64(current) / float64(required)
	filled := int(progress * float64(width))
	
	if filled > width {
		filled = width
	}
	
	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "■"
		} else {
			bar += "□"
		}
	}
	bar += "]"
	
	return bar
}

