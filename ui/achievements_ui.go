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