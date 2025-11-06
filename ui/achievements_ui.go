package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/achievements"
	"github.com/jamesacampbell/unicorn/animations"
	"github.com/jamesacampbell/unicorn/ascii"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/game"
	"github.com/jamesacampbell/unicorn/progression"
	"github.com/jamesacampbell/unicorn/upgrades"
)


func CheckAndUnlockAchievements(gs *game.GameState) {
	// Get player's previously unlocked achievements
	previouslyUnlocked, err := database.GetPlayerAchievements(gs.PlayerName)
	if err != nil {
		previouslyUnlocked = []string{}
	}

	// Get player stats
	stats, _ := database.GetPlayerStats(gs.PlayerName)
	winStreak, _ := database.GetWinStreak(gs.PlayerName)

	// Count sectors and get investment details
	sectors := make(map[string]bool)
	positiveCount := 0
	negativeCount := 0
	totalInvested := int64(0)
	riskScores := []float64{}

	for _, inv := range gs.Portfolio.Investments {
		totalInvested += inv.AmountInvested
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		if value > inv.AmountInvested {
			positiveCount++
		} else if value < inv.AmountInvested {
			negativeCount++
		}

		// Find sector and risk score
		for _, startup := range gs.AvailableStartups {
			if startup.Name == inv.CompanyName {
				sectors[startup.Category] = true
				riskScores = append(riskScores, startup.RiskScore)
				break
			}
		}
	}

	sectorsInvested := []string{}
	for sector := range sectors {
		sectorsInvested = append(sectorsInvested, sector)
	}

	netWorth, roi, successfulExits := gs.GetFinalScore()

	// Build game stats for achievement checking
	gameStats := achievements.GameStats{
		GameMode:            "vc",
		FinalNetWorth:       netWorth,
		ROI:                 roi,
		SuccessfulExits:     successfulExits,
		TurnsPlayed:         gs.Portfolio.Turn - 1,
		Difficulty:          gs.Difficulty.Name,
		InvestmentCount:     len(gs.Portfolio.Investments),
		SectorsInvested:     sectorsInvested,
		TotalInvested:       totalInvested,
		RiskScores:          riskScores,
		PositiveInvestments: positiveCount,
		NegativeInvestments: negativeCount,
		TotalGames:          stats.TotalGames,
		TotalWins:           int(stats.WinRate * float64(stats.TotalGames) / 100.0),
		WinStreak:           winStreak,
		BestNetWorth:        stats.BestNetWorth,
		TotalExits:          stats.TotalExits,
	}

	// Check for new achievements
	newAchievements := achievements.CheckAchievements(gameStats, previouslyUnlocked)
	
	// Get list of newly unlocked achievement IDs
	newAchievementIDs := []string{}
	for _, ach := range newAchievements {
		newAchievementIDs = append(newAchievementIDs, ach.ID)
	}

	// Always show achievement section
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println("\n" + strings.Repeat("=", 70))
	cyan.Printf("                    ACHIEVEMENT CHECK\n")
	fmt.Println(strings.Repeat("=", 70))

	// Save and display new achievements
	if len(newAchievements) > 0 {
		fmt.Println("\n" + strings.Repeat("?", 60))
		cyan.Printf("     %s NEW ACHIEVEMENTS UNLOCKED! %s\n", ascii.Star, ascii.Star)
		fmt.Println(strings.Repeat("?", 60))

		for _, ach := range newAchievements {
			// Save to database
			database.UnlockAchievement(gs.PlayerName, ach.ID)

			// Display with animation
			achievementText := fmt.Sprintf("%s %s [%s]\n+%d points", ach.Icon, ach.Name, ach.Rarity, ach.Points)
			animations.ShowAchievementUnlock(achievementText, ach.Description)
		}
	} else {
		yellow.Println("\nNo new achievements unlocked this game.")
		yellow.Println("Keep playing to unlock more achievements!")
		fmt.Println("\nTips to unlock achievements:")
		fmt.Println("  • Wealth: Reach net worth milestones ($1M, $5M, $10M, $50M)")
		fmt.Println("  • Performance: Achieve positive ROI (break even, 2x, 5x, 10x returns)")
		fmt.Println("  • Strategy: Diversify investments, master sectors, get successful exits")
		fmt.Println("  • Career: Play more games, build win streaks")
	}
	
	// Calculate and award XP
	xpEarned := progression.CalculateXPReward(&gameStats, newAchievementIDs)
	
	// Get profile before adding XP to compare levels
	profileBefore, _ := database.GetPlayerProfile(gs.PlayerName)
	oldLevel := profileBefore.Level
	
	// Add XP to player profile
	leveledUp, newLevel, err := database.AddExperience(gs.PlayerName, xpEarned)
	if err != nil {
		color.Yellow("\nWarning: Could not add XP: %v", err)
	} else {
		// Display XP breakdown
		xpBreakdown := make(map[string]int)
		xpBreakdown["Game Completion"] = progression.XPGameComplete
		
		if gameStats.ROI > 0 {
			xpBreakdown["Positive ROI"] = progression.XPPositiveROI
		}
		
		if gameStats.SuccessfulExits > 0 {
			xpBreakdown[fmt.Sprintf("Successful Exits (%d)", gameStats.SuccessfulExits)] = progression.XPSuccessfulExit * gameStats.SuccessfulExits
		}
		
		switch strings.ToLower(gameStats.Difficulty) {
		case "medium":
			xpBreakdown["Medium Difficulty"] = progression.XPDifficultyMedium
		case "hard":
			xpBreakdown["Hard Difficulty"] = progression.XPDifficultyHard
		case "expert":
			xpBreakdown["Expert Difficulty"] = progression.XPDifficultyExpert
		}
		
		if len(newAchievementIDs) > 0 {
			achXP := 0
			for _, achvID := range newAchievementIDs {
				if achv, exists := achievements.AllAchievements[achvID]; exists {
					achXP += achv.Points * progression.XPAchievementBase
				}
			}
			if achXP > 0 {
				xpBreakdown[fmt.Sprintf("New Achievements (%d)", len(newAchievementIDs))] = achXP
			}
		}
		
		DisplayXPGained(xpBreakdown, xpEarned)
		
		// Show level up screen if leveled up
		if leveledUp {
			levelInfo := progression.GetLevelInfo(newLevel)
			DisplayLevelUp(gs.PlayerName, oldLevel, newLevel, levelInfo.Unlocks)
		} else {
			// Show progress towards next level
			profileAfter, _ := database.GetPlayerProfile(gs.PlayerName)
			fmt.Println()
			yellow.Printf("   Level %d Progress: ", profileAfter.Level)
			progressBar := progression.FormatXPBar(profileAfter.ExperiencePoints, profileAfter.NextLevelXP, 20)
			fmt.Printf("%s %d/%d XP\n", progressBar, profileAfter.ExperiencePoints, profileAfter.NextLevelXP)
		}
	}

	// Calculate and display career level and points (always show)
	totalLifetimePoints := 0
	allUnlocked, _ := database.GetPlayerAchievements(gs.PlayerName)
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalLifetimePoints += ach.Points
		}
	}

	// Get owned upgrades to calculate available balance
	ownedUpgrades, _ := database.GetPlayerUpgrades(gs.PlayerName)
	availableBalance := totalLifetimePoints
	spentOnUpgrades := 0
	for _, upgradeID := range ownedUpgrades {
		if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
			availableBalance -= upgrade.Cost
			spentOnUpgrades += upgrade.Cost
		}
	}

	level, title, _ := achievements.CalculateCareerLevel(totalLifetimePoints)
	green := color.New(color.FgGreen)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("Career Level: ")
	yellow.Printf("%d - %s", level, title)
	fmt.Printf("\nAvailable Balance: ")
	green.Printf("%d pts", availableBalance)
	fmt.Printf("\nTotal Lifetime Points: %d pts", totalLifetimePoints)
	if spentOnUpgrades > 0 {
		fmt.Printf(" (Spent: %d pts)", spentOnUpgrades)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
}

func ViewPlayerAchievements() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Print("\nEnter player name: ")
	reader := bufio.NewReader(os.Stdin)
	playerName, _ := reader.ReadString('\n')
	playerName = strings.TrimSpace(playerName)

	if playerName == "" {
		color.Red("Invalid player name!")
		return
	}

	unlocked, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		color.Red("Error loading achievements: %v", err)
		return
	}

	clear.ClearIt()
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Printf("     ACHIEVEMENTS FOR: %s\n", strings.ToUpper(playerName))
	cyan.Println(strings.Repeat("=", 70))

	// Calculate stats
	totalAchievements := len(achievements.AllAchievements)
	unlockedCount := len(unlocked)
	progress := float64(unlockedCount) / float64(totalAchievements) * 100

	// Calculate total points
	totalPoints := 0
	for _, id := range unlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalPoints += ach.Points
		}
	}

	// Get career level
	level, title, nextLevelPoints := achievements.CalculateCareerLevel(totalPoints)

	green := color.New(color.FgGreen, color.Bold)
	fmt.Printf("\n%s Progress: %d/%d (%.1f%%)\n", ascii.Chart, unlockedCount, totalAchievements, progress)
	fmt.Printf("%s Total Points: ", ascii.Coin)
	green.Printf("%d\n", totalPoints)
	fmt.Printf("%s Career Level: ", ascii.Level)
	yellow.Printf("%d - %s\n", level, title)
	if level < 10 {
		fmt.Printf("%s Next Level: %d points needed\n", ascii.Target, nextLevelPoints-totalPoints)
	}

	if unlockedCount == 0 {
		yellow.Println("\nNo achievements unlocked yet. Keep playing!")
		return
	}

	// Group by category
	categories := map[string][]achievements.Achievement{
		achievements.CategoryWealth:      {},
		achievements.CategoryPerformance: {},
		achievements.CategoryStrategy:    {},
		achievements.CategoryCareer:      {},
		achievements.CategoryChallenge:   {},
		achievements.CategorySpecial:     {},
	}

	unlockedMap := make(map[string]bool)
	for _, id := range unlocked {
		unlockedMap[id] = true
	}

	for id, ach := range achievements.AllAchievements {
		if unlockedMap[id] {
			categories[ach.Category] = append(categories[ach.Category], ach)
		}
	}

	// Display by category
	for _, category := range []string{
		achievements.CategoryWealth,
		achievements.CategoryPerformance,
		achievements.CategoryStrategy,
		achievements.CategoryCareer,
		achievements.CategoryChallenge,
		achievements.CategorySpecial,
	} {
		achs := categories[category]
		if len(achs) == 0 {
			continue
		}

		fmt.Printf("\n%s:\n", category)
		for _, ach := range achs {
			rarityColor := color.New(color.Attribute(achievements.GetRarityColor(ach.Rarity)))
			fmt.Printf("  %s ", ach.Icon)
			rarityColor.Printf("%s", ach.Name)
			fmt.Printf(" - %s (+%d pts)\n", ach.Description, ach.Points)
		}
	}
}

func BrowseAllAchievements() {
	cyan := color.New(color.FgCyan, color.Bold)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                  ALL ACHIEVEMENTS")
	cyan.Println(strings.Repeat("=", 70))

	// Group by category
	for _, category := range []string{
		achievements.CategoryWealth,
		achievements.CategoryPerformance,
		achievements.CategoryStrategy,
		achievements.CategoryCareer,
		achievements.CategoryChallenge,
		achievements.CategorySpecial,
	} {
		achs := achievements.GetAchievementsByCategory(category)
		if len(achs) == 0 {
			continue
		}

		yellow := color.New(color.FgYellow, color.Bold)
		yellow.Printf("\n%s:\n", category)

		for _, ach := range achs {
			rarityColor := color.New(color.Attribute(achievements.GetRarityColor(ach.Rarity)))
			fmt.Printf("  %s ", ach.Icon)
			rarityColor.Printf("%s", ach.Name)
			fmt.Printf(" [%s] - %s (+%d pts)\n", ach.Rarity, ach.Description, ach.Points)
		}
	}

	fmt.Printf("\n\nTotal Achievements: %d\n", len(achievements.AllAchievements))
}

func DisplayAchievementLeaderboard() {
	cyan := color.New(color.FgCyan, color.Bold)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("            ACHIEVEMENT LEADERBOARD (Coming Soon)")
	cyan.Println(strings.Repeat("=", 70))

	color.Yellow("\nThis feature will show players with the most achievements!")
}

// DisplayAchievementChains shows all achievement chains with progress
func DisplayAchievementChains(playerName string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	gray := color.New(color.FgHiBlack)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    ACHIEVEMENT CHAINS")
	cyan.Println(strings.Repeat("=", 70))
	
	// Get all chains
	chains := achievements.GetAllChains()
	
	if len(chains) == 0 {
		color.Yellow("\nNo achievement chains available yet!")
		return
	}
	
	// Display each chain
	for _, chainID := range chains {
		unlocked, total := achievements.GetChainProgress(playerName, chainID)
		chainAchievements := achievements.GetAchievementsByChain(chainID)
		
		fmt.Println()
		yellow.Printf("▶ %s Chain: ", strings.Title(strings.ReplaceAll(chainID, "_", " ")))
		green.Printf("%d/%d ", unlocked, total)
		
		// Progress bar
		progressBar := formatProgressBar(unlocked, total, 20)
		fmt.Printf("%s\n", progressBar)
		
		// Display achievements in chain
		playerAchievements, _ := database.GetPlayerAchievements(playerName)
		unlockedMap := make(map[string]bool)
		for _, id := range playerAchievements {
			unlockedMap[id] = true
		}
		
		for i, achv := range chainAchievements {
			prefix := "   "
			if i > 0 {
				prefix = "   ↓ "
			}
			
			if unlockedMap[achv.ID] {
				green.Printf("%s✓ %s %s", prefix, achv.Icon, achv.Name)
				fmt.Printf(" (%d pts)\n", achv.Points)
			} else if achievements.CheckAchievementChain(playerName, achv.ID) {
				yellow.Printf("%s○ %s %s", prefix, achv.Icon, achv.Name)
				fmt.Printf(" (%d pts) - Available\n", achv.Points)
			} else {
				gray.Printf("%s● %s %s", prefix, achv.Icon, achv.Name)
				fmt.Printf(" (%d pts) - Locked\n", achv.Points)
			}
			
			// Show progress for progressive achievements
			if achv.ProgressTracking {
				current, max, err := database.GetAchievementProgress(playerName, achv.ID)
				if err == nil && max > 0 {
					progBar := formatProgressBar(current, max, 15)
					fmt.Printf("        Progress: %s %d/%d\n", progBar, current, max)
				}
			}
		}
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayHiddenAchievements shows hidden achievements (only if unlocked)
func DisplayHiddenAchievements(playerName string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta, color.Bold)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    HIDDEN ACHIEVEMENTS")
	cyan.Println(strings.Repeat("=", 70))
	
	hidden := achievements.GetHiddenAchievements(playerName)
	
	if len(hidden) == 0 {
		yellow.Println("\nYou haven't unlocked any hidden achievements yet!")
		yellow.Println("Keep playing to discover secret challenges...")
		
		// Count total hidden achievements
		totalHidden := 0
		for _, achv := range achievements.AllAchievements {
			if achv.Hidden {
				totalHidden++
			}
		}
		
		fmt.Printf("\nHidden Achievements Available: %d\n", totalHidden)
	} else {
		fmt.Printf("\n%s You've discovered %d secret achievement(s)!\n", ascii.Star, len(hidden))
		fmt.Println()
		
		// Display unlocked hidden achievements
		for _, achv := range hidden {
			rarityColor := getRarityColor(achv.Rarity)
			rarityColor.Printf("   %s %s\n", achv.Icon, achv.Name)
			fmt.Printf("      %s\n", achv.Description)
			magenta.Printf("      %s • %d pts\n", achv.Rarity, achv.Points)
			fmt.Println()
		}
	}
	
	cyan.Println(strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayProgressiveAchievements shows achievements with progress tracking
func DisplayProgressiveAchievements(playerName string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                  PROGRESSIVE ACHIEVEMENTS")
	cyan.Println(strings.Repeat("=", 70))
	
	progressiveAchievements := achievements.GetProgressiveAchievements()
	
	if len(progressiveAchievements) == 0 {
		color.Yellow("\nNo progressive achievements available!")
		return
	}
	
	// Get player's unlocked achievements
	playerAchievements, _ := database.GetPlayerAchievements(playerName)
	unlockedMap := make(map[string]bool)
	for _, id := range playerAchievements {
		unlockedMap[id] = true
	}
	
	// Separate into unlocked and in-progress
	fmt.Println("\n📈 IN PROGRESS:")
	inProgressCount := 0
	
	for _, achv := range progressiveAchievements {
		if !unlockedMap[achv.ID] {
			current, max, err := database.GetAchievementProgress(playerName, achv.ID)
			if err == nil && current > 0 {
				yellow.Printf("\n   %s %s\n", achv.Icon, achv.Name)
				fmt.Printf("      %s\n", achv.Description)
				
				progBar := formatProgressBar(current, max, 30)
				fmt.Printf("      %s ", progBar)
				green.Printf("%d/%d ", current, max)
				fmt.Printf("(%.1f%%)\n", float64(current)/float64(max)*100)
				fmt.Printf("      %s • %d pts\n", achv.Rarity, achv.Points)
				
				inProgressCount++
			}
		}
	}
	
	if inProgressCount == 0 {
		fmt.Println("   No achievements in progress yet!")
	}
	
	// Show completed progressive achievements
	fmt.Println("\n✓ COMPLETED:")
	completedCount := 0
	
	for _, achv := range progressiveAchievements {
		if unlockedMap[achv.ID] {
			green.Printf("\n   %s %s\n", achv.Icon, achv.Name)
			fmt.Printf("      %s\n", achv.Description)
			progBar := formatProgressBar(achv.MaxProgress, achv.MaxProgress, 30)
			fmt.Printf("      %s ", progBar)
			green.Printf("%d/%d ", achv.MaxProgress, achv.MaxProgress)
			fmt.Printf("(100%%)\n")
			completedCount++
		}
	}
	
	if completedCount == 0 {
		fmt.Println("   No progressive achievements completed yet!")
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Helper function to format a progress bar
func formatProgressBar(current, max, width int) string {
	if max == 0 {
		return "[" + strings.Repeat("□", width) + "]"
	}
	
	progress := float64(current) / float64(max)
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

// Helper function to get color for rarity
func getRarityColor(rarity string) *color.Color {
	switch rarity {
	case achievements.RarityCommon:
		return color.New(color.FgWhite)
	case achievements.RarityRare:
		return color.New(color.FgCyan)
	case achievements.RarityEpic:
		return color.New(color.FgMagenta)
	case achievements.RarityLegendary:
		return color.New(color.FgYellow, color.Bold)
	default:
		return color.New(color.FgWhite)
	}
}

// EnhancedAchievementMenu provides options for viewing achievements with new features
func EnhancedAchievementMenu() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	for {
		clear.ClearIt()
		
		cyan.Println("\n" + strings.Repeat("=", 60))
		cyan.Println("                ACHIEVEMENTS MENU")
		cyan.Println(strings.Repeat("=", 60))
		
		yellow.Println("\n1. View All Achievements")
		yellow.Println("2. View Achievement Chains")
		yellow.Println("3. View Progressive Achievements")
		yellow.Println("4. View Hidden Achievements")
		yellow.Println("5. Browse All")
		yellow.Println("6. Back to Main Menu")
		
		fmt.Print("\nEnter your choice: ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		if choice == "6" {
			return
		}
		
		if choice == "5" {
			BrowseAllAchievements()
			continue
		}
		
		// Get player name for other options
		fmt.Print("\nEnter your name: ")
		playerName, _ := reader.ReadString('\n')
		playerName = strings.TrimSpace(playerName)
		
		if playerName == "" {
			color.Red("Invalid player name!")
			continue
		}
		
		switch choice {
		case "1":
			ViewPlayerAchievements()
		case "2":
			DisplayAchievementChains(playerName)
		case "3":
			DisplayProgressiveAchievements(playerName)
		case "4":
			DisplayHiddenAchievements(playerName)
		default:
			color.Red("Invalid choice!")
		}
	}
}