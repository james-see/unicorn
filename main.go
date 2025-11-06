package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/animations"
	"github.com/jamesacampbell/unicorn/ascii"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/ui"
	"github.com/jamesacampbell/unicorn/upgrades"
	yaml "gopkg.in/yaml.v2"
)

type gameData struct {
	Pot        int64  `yaml:"starting-cash"`
	BadThings  int64  `yaml:"number-of-bad-things-per-year"`
	Foreground string `yaml:"foreground-color"`
}

func initMenu() (username string) {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Name: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	// Show spinner while checking player stats
	spinner, _ := animations.StartSpinner("Loading player data...")
	stats, err := database.GetPlayerStats(text)
	spinner.Stop()
	if err == nil && stats.TotalGames > 0 {
		// Returning player - show welcome back message with stats
		fmt.Println()
		green.Printf("🎉 Welcome back, %s!\n\n", text)

		cyan.Println(strings.Repeat("=", 60))
		cyan.Println("                  YOUR PLAYER STATS")
		cyan.Println(strings.Repeat("=", 60))

		yellow.Printf("\n📊 Games Played: %d\n", stats.TotalGames)
		yellow.Printf("💰 Best Net Worth: $%s\n", ui.FormatCurrency(stats.BestNetWorth))
		yellow.Printf("📈 Best ROI: %.1f%%\n", stats.BestROI*100)
		yellow.Printf("🚀 Total Exits: %d\n", stats.TotalExits)
		yellow.Printf("📊 Average Net Worth: $%s\n", ui.FormatCurrency(int64(stats.AverageNetWorth)))
		yellow.Printf("🎯 Win Rate: %.1f%%\n", stats.WinRate)

		// Get achievement count
		spinner2, _ := animations.StartSpinner("Loading achievements...")
		achievementCount, _ := database.GetPlayerAchievementCount(text)
		spinner2.Stop()
		if achievementCount > 0 {
			yellow.Printf("🏆 Achievements Unlocked: %d\n", achievementCount)
		}

		// Get and display active upgrades
		spinner3, _ := animations.StartSpinner("Loading upgrades...")
		playerUpgrades, err := database.GetPlayerUpgrades(text)
		spinner3.Stop()
		if err == nil && len(playerUpgrades) > 0 {
			green := color.New(color.FgGreen)
			fmt.Println()
			green.Println("✨ ACTIVE UPGRADES:")
			for _, upgradeID := range playerUpgrades {
				if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
					fmt.Printf("   %s %s - %s\n", upgrade.Icon, upgrade.Name, upgrade.Description)
				}
			}
		}

		cyan.Println(strings.Repeat("=", 60))
		fmt.Println()

		fmt.Print("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		fmt.Println()
	} else {
		// New player
		fmt.Printf("\nWelcome %s!\n", text)
	}

	return text
}

func loadConfig() gameData {
	var gd gameData
	yamlFile, err := os.Open("config/data.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer yamlFile.Close()
	byteValue, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(byteValue, &gd)
	if err != nil {
		fmt.Println(err)
	}
	return gd
}

func playNewGame() {
	// Get username
	username := initMenu()
	clear.ClearIt()

	// Select game mode
	gameMode := askForGameMode()
	clear.ClearIt()

	if gameMode == "founder" {
		ui.PlayFounderMode(username)
	} else {
		ui.PlayVCMode(username)
	}
}

func askForGameMode() string {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("                 GAME MODE SELECTION")
	cyan.Println(strings.Repeat("=", 60))

	yellow.Println("\n1. VC Investor Mode (Classic)")
	fmt.Println("   Build a portfolio of startups and compete against AI investors")

	yellow.Println("\n2. Startup Founder Mode (New!)")
	fmt.Println("   Build your own startup from the ground up")

	fmt.Print("\nEnter your choice (1-2, default 1): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Show spinner after input
	spinner, _ := animations.StartSpinner("Loading game mode...")
	time.Sleep(500 * time.Millisecond)
	spinner.Stop()

	if choice == "2" {
		return "founder"
	}
	return "vc"
}

func main() {
	// Show animated splash screen on first launch
	animations.ShowGameStartAnimation()

	// Pause to let user enjoy the splash screen
	fmt.Print("\nPress 'Enter' to continue to main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Initialize database
	err := database.InitDB("unicorn_scores.db")
	if err != nil {
		fmt.Printf("Warning: Could not initialize database: %v\n", err)
		fmt.Println("Scores will not be saved.")
		time.Sleep(2 * time.Second)
	}
	defer database.CloseDB()

	// Main menu loop
	for {
		choice := ui.DisplayMainMenu()
		clear.ClearIt()

		switch choice {
		case "1":
			playNewGame()
		case "2":
			ui.DisplayLeaderboards()
		case "3":
			ui.DisplayPlayerStats()
		case "4":
			ui.DisplayAchievementsMenu()
		case "5":
			ui.DisplayUpgradeMenu()
		case "6":
			ui.DisplayHelpGuide()
		case "7":
			animations.ShowInfoMessage("Thanks for playing! " + ascii.Star2)
			return
		default:
			color.Red("Invalid choice!")
			time.Sleep(1 * time.Second)
		}
	}
}
