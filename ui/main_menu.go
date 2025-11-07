package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/achievements"
	"github.com/jamesacampbell/unicorn/animations"
	"github.com/jamesacampbell/unicorn/ascii"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/logo"
	"github.com/jamesacampbell/unicorn/upgrades"
)


func DisplayMainMenu() string {
	clear.ClearIt()

	// Display unicorn logo
	cyan := color.New(color.FgCyan, color.Bold)
	logo.InitLogo(cyan)

	yellow := color.New(color.FgYellow)

	cyan.Println("\n" + strings.Repeat("=", 50))
	cyan.Println("           ? UNICORN - MAIN MENU ?")
	cyan.Println(strings.Repeat("=", 50))

	yellow.Println("\n1. New Game")
	yellow.Println("2. Leaderboards")
	yellow.Println("3. Player Statistics")
	yellow.Println("4. Achievements")
	yellow.Println("5. Upgrades")
	yellow.Println("6. Progression & Levels")
	yellow.Println("7. Analytics Dashboard")
	yellow.Println("8. Help & Info")
	yellow.Println("9. Quit")

	fmt.Print("\nEnter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Show spinner after input
	spinner, _ := animations.StartSpinner("Processing...")
	time.Sleep(500 * time.Millisecond)
	spinner.Stop()

	return choice
}

func DisplayLeaderboards() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)

	cyan.Print(ascii.LeaderboardHeader)

	fmt.Println("\n1. By Net Worth (All Difficulties)")
	fmt.Println("2. By ROI (All Difficulties)")
	fmt.Println("3. Easy Difficulty")
	fmt.Println("4. Medium Difficulty")
	fmt.Println("5. Hard Difficulty")
	fmt.Println("6. Expert Difficulty")
	fmt.Println("7. Recent Games")
	fmt.Println("8. Back to Main Menu")

	fmt.Print("\nEnter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	clear.ClearIt()

	switch choice {
	case "1":
		ShowTopScores("net_worth", "all")
	case "2":
		ShowTopScores("roi", "all")
	case "3":
		ShowTopScores("net_worth", "Easy")
	case "4":
		ShowTopScores("net_worth", "Medium")
	case "5":
		ShowTopScores("net_worth", "Hard")
	case "6":
		ShowTopScores("net_worth", "Expert")
	case "7":
		ShowRecentGames()
	case "8":
		return
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	DisplayLeaderboards()
}

func DisplayAchievementsMenu() {
	// Use enhanced achievement menu with new features
	EnhancedAchievementMenu()
}

func DisplayUpgradeMenu() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	fmt.Print("\nEnter player name: ")
	reader := bufio.NewReader(os.Stdin)
	playerName, _ := reader.ReadString('\n')
	playerName = strings.TrimSpace(playerName)

	if playerName == "" {
		color.Red("Invalid player name!")
		return
	}

	// Get player's points
	allUnlocked, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		allUnlocked = []string{}
	}

	totalLifetimePoints := 0
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalLifetimePoints += ach.Points
		}
	}
	
	// Add level-up points
	profile, _ := database.GetPlayerProfile(playerName)
	if profile != nil {
		totalLifetimePoints += profile.LevelUpPoints
	}

	// Get owned upgrades
	ownedUpgrades, err := database.GetPlayerUpgrades(playerName)
	if err != nil {
		ownedUpgrades = []string{}
	}

	// Calculate available balance (total points - spent on upgrades)
	availableBalance := totalLifetimePoints
	spentOnUpgrades := 0
	for _, upgradeID := range ownedUpgrades {
		if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
			availableBalance -= upgrade.Cost
			spentOnUpgrades += upgrade.Cost
		}
	}

	clear.ClearIt()
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Printf("     🎁 UPGRADE STORE 🎁\n")
	cyan.Println(strings.Repeat("=", 70))

	yellow.Printf("\nPlayer: %s\n", playerName)
	green.Printf("Available Balance: %d pts\n", availableBalance)
	fmt.Printf("Total Lifetime Points: %d pts", totalLifetimePoints)
	if spentOnUpgrades > 0 {
		fmt.Printf(" (Spent: %d pts)", spentOnUpgrades)
	}
	fmt.Println()

	level, title, _ := achievements.CalculateCareerLevel(totalLifetimePoints)
	fmt.Printf("Career Level: ")
	green.Printf("%d - %s\n", level, title)

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("1. Browse All Upgrades")
	fmt.Println("2. View My Upgrades")
	fmt.Println("3. Purchase Upgrades")
	fmt.Println("4. Back to Main Menu")
	fmt.Print("\nEnter your choice: ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	clear.ClearIt()

	switch choice {
	case "1":
		BrowseAllUpgrades(playerName, availableBalance, ownedUpgrades)
	case "2":
		ViewPlayerUpgrades(playerName, ownedUpgrades)
	case "3":
		PurchaseUpgrades(playerName, availableBalance, ownedUpgrades)
	case "4":
		return
	default:
		color.Red("Invalid choice!")
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	DisplayUpgradeMenu()
}

func BrowseAllUpgrades(playerName string, availableBalance int, ownedUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)

	categories := upgrades.GetAllCategories()

	for _, category := range categories {
		cyan.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
		cyan.Printf("║  %s\n", category)
		cyan.Printf("╚══════════════════════════════════════════════════════════════╝\n")

		categoryUpgrades := upgrades.GetUpgradesByCategory(category)
		for i, upgrade := range categoryUpgrades {
			owned := upgrades.IsOwned(upgrade.ID, ownedUpgrades)
			canAfford := availableBalance >= upgrade.Cost

			status := ""
			if owned {
				status = green.Sprintf("[✓ OWNED]")
			} else if canAfford {
				status = yellow.Sprintf("[AVAILABLE]")
			} else {
				status = magenta.Sprintf("[Need %d more pts]", upgrade.Cost-availableBalance)
			}

			fmt.Printf("\n%d. %s %s\n", i+1, upgrade.Icon, upgrade.Name)
			fmt.Printf("   %s\n", upgrade.Description)
			fmt.Printf("   Cost: %d points %s\n", upgrade.Cost, status)
		}
	}
}

func ViewPlayerUpgrades(playerName string, ownedUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)

	if len(ownedUpgrades) == 0 {
		color.Yellow("\nYou haven't purchased any upgrades yet!")
		return
	}

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Printf("     YOUR UPGRADES\n")
	cyan.Println(strings.Repeat("=", 70))

	categories := upgrades.GetAllCategories()
	for _, category := range categories {
		categoryUpgrades := upgrades.GetUpgradesByCategory(category)
		hasOwned := false
		for _, upgrade := range categoryUpgrades {
			if upgrades.IsOwned(upgrade.ID, ownedUpgrades) {
				if !hasOwned {
					fmt.Printf("\n%s:\n", category)
					hasOwned = true
				}
				green.Printf("  ✓ %s %s - %s\n", upgrade.Icon, upgrade.Name, upgrade.Description)
			}
		}
	}
}

func PurchaseUpgrades(playerName string, totalPoints int, ownedUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Refresh points and owned upgrades from database
	allUnlocked, err := database.GetPlayerAchievements(playerName)
	if err != nil {
		allUnlocked = []string{}
	}

	currentPoints := 0
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			currentPoints += ach.Points
		}
	}
	
	// Add level-up points
	profilePurchase, _ := database.GetPlayerProfile(playerName)
	if profilePurchase != nil {
		currentPoints += profilePurchase.LevelUpPoints
	}

	// Deduct cost of owned upgrades
	currentOwnedUpgrades, err := database.GetPlayerUpgrades(playerName)
	if err != nil {
		currentOwnedUpgrades = []string{}
	}

	// Recalculate points after subtracting purchased upgrade costs
	for _, upgradeID := range currentOwnedUpgrades {
		if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
			currentPoints -= upgrade.Cost
		}
	}

	// Use refreshed values
	totalPoints = currentPoints
	ownedUpgrades = currentOwnedUpgrades

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Printf("     PURCHASE UPGRADES\n")
	cyan.Println(strings.Repeat("=", 70))

	yellow.Printf("\nYour Points: %d\n\n", totalPoints)

	// Show available upgrades
	availableUpgrades := []upgrades.Upgrade{}
	for _, upgrade := range upgrades.AllUpgrades {
		if !upgrades.IsOwned(upgrade.ID, ownedUpgrades) && totalPoints >= upgrade.Cost {
			availableUpgrades = append(availableUpgrades, upgrade)
		}
	}

	if len(availableUpgrades) == 0 {
		color.Yellow("\nNo upgrades available for purchase!")
		color.Yellow("Earn more achievement points to unlock upgrades.")
		return
	}

	fmt.Println("Available Upgrades:")
	for i, upgrade := range availableUpgrades {
		fmt.Printf("%d. %s %s - %d pts\n", i+1, upgrade.Icon, upgrade.Name, upgrade.Cost)
		fmt.Printf("   %s\n", upgrade.Description)
	}

	fmt.Print("\nEnter upgrade number to purchase (or 0 to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)

	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 0 || choice > len(availableUpgrades) {
		color.Red("Invalid choice!")
		return
	}

	if choice == 0 {
		return
	}

	upgrade := availableUpgrades[choice-1]

	if totalPoints < upgrade.Cost {
		color.Red("Insufficient points! Need %d, have %d", upgrade.Cost, totalPoints)
		return
	}

	// Purchase upgrade
	err = database.PurchaseUpgrade(playerName, upgrade.ID)
	if err != nil {
		color.Red("Error purchasing upgrade: %v", err)
		return
	}

	green.Printf("\n✓ Successfully purchased: %s %s!\n", upgrade.Icon, upgrade.Name)

	// Refresh points and owned upgrades from database
	allUnlockedRefresh, errRefresh := database.GetPlayerAchievements(playerName)
	if errRefresh != nil {
		allUnlockedRefresh = []string{}
	}

	newTotalPoints := 0
	for _, id := range allUnlockedRefresh {
		if ach, exists := achievements.AllAchievements[id]; exists {
			newTotalPoints += ach.Points
		}
	}
	
	// Add level-up points
	profileRefresh, _ := database.GetPlayerProfile(playerName)
	if profileRefresh != nil {
		newTotalPoints += profileRefresh.LevelUpPoints
	}

	newOwnedUpgrades, errUpgrades := database.GetPlayerUpgrades(playerName)
	if errUpgrades != nil {
		newOwnedUpgrades = []string{}
	}

	// Deduct cost of all owned upgrades
	for _, upgradeID := range newOwnedUpgrades {
		if up, exists := upgrades.AllUpgrades[upgradeID]; exists {
			newTotalPoints -= up.Cost
		}
	}

	green.Printf("Points remaining: %d\n", newTotalPoints)

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Return to purchase menu with refreshed values
	PurchaseUpgrades(playerName, newTotalPoints, newOwnedUpgrades)
}