package achievements

import (
	"github.com/jamesacampbell/unicorn/database"
)

// CheckAchievementChain checks if a player has unlocked all prerequisites for an achievement
func CheckAchievementChain(playerName, achievementID string) bool {
	achievement, exists := AllAchievements[achievementID]
	if !exists {
		return false
	}
	
	// No prerequisites means it's unlocked
	if len(achievement.RequiredAchievements) == 0 {
		return true
	}
	
	// Check if player has all required achievements
	playerAchievements, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		return false
	}
	
	unlockedMap := make(map[string]bool)
	for _, achv := range playerAchievements {
		unlockedMap[achv] = true
	}
	
	// Check all requirements
	for _, requiredID := range achievement.RequiredAchievements {
		if !unlockedMap[requiredID] {
			return false // Missing a required achievement
		}
	}
	
	return true
}

// GetNextInChain returns the next achievement in a chain after the given achievement
func GetNextInChain(playerName, achievementID string) *Achievement {
	achievement, exists := AllAchievements[achievementID]
	if !exists || achievement.ChainID == "" {
		return nil
	}
	
	// Find the next achievement in this chain
	for _, achv := range AllAchievements {
		if achv.ChainID == achievement.ChainID && 
		   len(achv.RequiredAchievements) > 0 &&
		   contains(achv.RequiredAchievements, achievementID) {
			// Check if player can unlock it
			if CheckAchievementChain(playerName, achv.ID) {
				return &achv
			}
		}
	}
	
	return nil
}

// GetProgressiveAchievements returns all achievements with progress tracking
func GetProgressiveAchievements() []Achievement {
	var achievements []Achievement
	for _, achv := range AllAchievements {
		if achv.ProgressTracking {
			achievements = append(achievements, achv)
		}
	}
	return achievements
}

// UpdateProgress updates progress toward a progressive achievement
func UpdateProgress(playerName, achievementID string, currentProgress int) error {
	achievement, exists := AllAchievements[achievementID]
	if !exists || !achievement.ProgressTracking {
		return nil // Not a progressive achievement
	}
	
	// Update progress in database
	err := database.UpdateAchievementProgress(playerName, achievementID, currentProgress, achievement.MaxProgress)
	if err != nil {
		return err
	}
	
	// Check if achievement should be unlocked
	if currentProgress >= achievement.MaxProgress {
		// Check if chain requirements are met
		if CheckAchievementChain(playerName, achievementID) {
			// Check if not already unlocked
			playerAchievements, err := database.GetPlayerAchievements(playerName)
			if err == nil {
				alreadyUnlocked := false
				for _, id := range playerAchievements {
					if id == achievementID {
						alreadyUnlocked = true
						break
					}
				}
				
				if !alreadyUnlocked {
					// Unlock the achievement
					err = database.UnlockAchievement(playerName, achievementID)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	
	return nil
}

// GetAchievementsByChain returns all achievements in a specific chain
func GetAchievementsByChain(chainID string) []Achievement {
	var achievements []Achievement
	for _, achv := range AllAchievements {
		if achv.ChainID == chainID {
			achievements = append(achievements, achv)
		}
	}
	return achievements
}

// GetAllChains returns a list of all unique chain IDs
func GetAllChains() []string {
	chainMap := make(map[string]bool)
	for _, achv := range AllAchievements {
		if achv.ChainID != "" {
			chainMap[achv.ChainID] = true
		}
	}
	
	chains := []string{}
	for chain := range chainMap {
		chains = append(chains, chain)
	}
	return chains
}

// GetChainProgress returns progress through a chain (unlocked count / total count)
func GetChainProgress(playerName, chainID string) (unlocked, total int) {
	chainAchievements := GetAchievementsByChain(chainID)
	total = len(chainAchievements)
	
	if total == 0 {
		return 0, 0
	}
	
	playerAchievements, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		return 0, total
	}
	
	unlockedMap := make(map[string]bool)
	for _, id := range playerAchievements {
		unlockedMap[id] = true
	}
	
	for _, achv := range chainAchievements {
		if unlockedMap[achv.ID] {
			unlocked++
		}
	}
	
	return unlocked, total
}

// IsAchievementLocked checks if an achievement is locked due to chain requirements
func IsAchievementLocked(playerName, achievementID string) bool {
	return !CheckAchievementChain(playerName, achievementID)
}

// GetLockedAchievements returns all achievements that are locked for a player
func GetLockedAchievements(playerName string) []Achievement {
	playerAchievements, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		return []Achievement{}
	}
	
	unlockedMap := make(map[string]bool)
	for _, id := range playerAchievements {
		unlockedMap[id] = true
	}
	
	locked := []Achievement{}
	for _, achv := range AllAchievements {
		if !unlockedMap[achv.ID] && !achv.Hidden {
			locked = append(locked, achv)
		}
	}
	
	return locked
}

// GetUnlockedAchievements returns all achievements that are unlocked for a player
func GetUnlockedAchievements(playerName string) []Achievement {
	playerAchievementIDs, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		return []Achievement{}
	}
	
	unlocked := []Achievement{}
	for _, id := range playerAchievementIDs {
		if achv, exists := AllAchievements[id]; exists {
			unlocked = append(unlocked, achv)
		}
	}
	
	return unlocked
}

// GetHiddenAchievements returns all hidden achievements (only shown after unlock)
func GetHiddenAchievements(playerName string) []Achievement {
	playerAchievements, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		return []Achievement{}
	}
	
	unlockedMap := make(map[string]bool)
	for _, id := range playerAchievements {
		unlockedMap[id] = true
	}
	
	hidden := []Achievement{}
	for _, achv := range AllAchievements {
		if achv.Hidden && unlockedMap[achv.ID] {
			// Only show hidden achievements if player has unlocked them
			hidden = append(hidden, achv)
		}
	}
	
	return hidden
}

// GetAvailableAchievements returns achievements that can be earned right now (not locked by chains)
func GetAvailableAchievements(playerName string) []Achievement {
	playerAchievements, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		return []Achievement{}
	}
	
	unlockedMap := make(map[string]bool)
	for _, id := range playerAchievements {
		unlockedMap[id] = true
	}
	
	available := []Achievement{}
	for _, achv := range AllAchievements {
		if !unlockedMap[achv.ID] && CheckAchievementChain(playerName, achv.ID) {
			available = append(available, achv)
		}
	}
	
	return available
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

