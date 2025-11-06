package progression

// Unlock represents something that becomes available at a certain level
type Unlock struct {
	Level       int
	Feature     string
	Name        string
	Description string
	Icon        string
}

// AllUnlocks defines everything that unlocks at each level
var AllUnlocks = []Unlock{
	// Level 1 - Starting content
	{
		Level:       1,
		Feature:     "easy",
		Name:        "Easy Difficulty",
		Description: "Perfect for learning the game mechanics",
		Icon:        "🎮",
	},
	{
		Level:       1,
		Feature:     "medium",
		Name:        "Medium Difficulty",
		Description: "Balanced challenge for most players",
		Icon:        "🎮",
	},
	{
		Level:       1,
		Feature:     "founder_mode",
		Name:        "Founder Mode",
		Description: "Build your own startup from scratch",
		Icon:        "🚀",
	},
	
	// Level 5 - Hard mode
	{
		Level:       5,
		Feature:     "hard",
		Name:        "Hard Difficulty",
		Description: "For experienced investors seeking a challenge",
		Icon:        "🔥",
	},
	
	// Level 10 - Expert and analytics
	{
		Level:       10,
		Feature:     "expert",
		Name:        "Expert Difficulty",
		Description: "Only for the most skilled players",
		Icon:        "💎",
	},
	{
		Level:       10,
		Feature:     "analytics",
		Name:        "Advanced Analytics",
		Description: "Detailed performance tracking and insights",
		Icon:        "📊",
	},
	
	// Level 15 - Nightmare and chains
	{
		Level:       15,
		Feature:     "nightmare",
		Name:        "Nightmare Mode",
		Description: "Extreme difficulty with harsh penalties",
		Icon:        "👹",
	},
	{
		Level:       15,
		Feature:     "achievement_chains",
		Name:        "Achievement Chains",
		Description: "See connected achievements and progression paths",
		Icon:        "🔗",
	},
	{
		Level:       15,
		Feature:     "market_cycles",
		Name:        "Market Cycles",
		Description: "Economic cycles affect investments",
		Icon:        "📈",
	},
	
	// Level 20 - Game modifiers
	{
		Level:       20,
		Feature:     "game_modifiers",
		Name:        "Game Modifiers",
		Description: "Custom rules to change gameplay",
		Icon:        "⚙️",
	},
	{
		Level:       20,
		Feature:     "secondary_market",
		Name:        "Secondary Markets",
		Description: "Buy and sell stakes before exits",
		Icon:        "💱",
	},
	
	// Level 25 - Prestige
	{
		Level:       25,
		Feature:     "prestige",
		Name:        "Prestige System",
		Description: "Reset progress for permanent bonuses",
		Icon:        "✨",
	},
	{
		Level:       25,
		Feature:     "legendary_achievements",
		Name:        "Legendary Achievements",
		Description: "The rarest and most challenging achievements",
		Icon:        "🏆",
	},
	
	// Level 30 - Master difficulty
	{
		Level:       30,
		Feature:     "master",
		Name:        "Master Difficulty",
		Description: "The ultimate test of skill",
		Icon:        "👑",
	},
	
	// Level 50 - Titan status
	{
		Level:       50,
		Feature:     "titan",
		Name:        "Titan Status",
		Description: "You've mastered everything",
		Icon:        "🌟",
	},
}

// GetUnlocksForLevel returns all unlocks that become available at a specific level
func GetUnlocksForLevel(level int) []Unlock {
	unlocks := []Unlock{}
	for _, unlock := range AllUnlocks {
		if unlock.Level == level {
			unlocks = append(unlocks, unlock)
		}
	}
	return unlocks
}

// GetNextUnlock returns the next unlock after the current level
func GetNextUnlock(currentLevel int) *Unlock {
	for _, unlock := range AllUnlocks {
		if unlock.Level > currentLevel {
			return &unlock
		}
	}
	return nil
}

// GetAllUnlockedFeatures returns all features unlocked at or below a level
func GetAllUnlockedFeatures(level int) []string {
	features := []string{}
	seen := make(map[string]bool)
	
	for _, unlock := range AllUnlocks {
		if unlock.Level <= level {
			if !seen[unlock.Feature] {
				features = append(features, unlock.Feature)
				seen[unlock.Feature] = true
			}
		}
	}
	
	return features
}

// GetLevelForFeature returns what level is required to unlock a feature
func GetLevelForFeature(featureName string) int {
	for _, unlock := range AllUnlocks {
		if unlock.Feature == featureName {
			return unlock.Level
		}
	}
	return 1 // Default to level 1 if not found
}

// IsFeatureUnlocked checks if a feature is unlocked at a given level
func IsFeatureUnlocked(level int, featureName string) bool {
	requiredLevel := GetLevelForFeature(featureName)
	return level >= requiredLevel
}

