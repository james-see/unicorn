package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/animations"
	"github.com/jamesacampbell/unicorn/ascii"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/upgrades"
)


func DisplayPlayerStats() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)

	cyan.Println("\n" + strings.Repeat("?", 50))
	cyan.Println("           PLAYER STATISTICS")
	cyan.Println(strings.Repeat("?", 50))

	fmt.Print("\nEnter player name: ")
	reader := bufio.NewReader(os.Stdin)
	playerName, _ := reader.ReadString('\n')
	playerName = strings.TrimSpace(playerName)

	if playerName == "" {
		color.Red("Invalid player name!")
		time.Sleep(1 * time.Second)
		return
	}

	// Show spinner while loading
	spinner, _ := animations.StartSpinner("Loading player stats...")
	defer spinner.Stop()

	stats, err := database.GetPlayerStats(playerName)
	if err != nil {
		color.Red("Error loading stats: %v", err)
		time.Sleep(2 * time.Second)
		return
	}

	if stats.TotalGames == 0 {
		color.Yellow("\nNo games found for player: %s", playerName)
		fmt.Print("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	clear.ClearIt()
	cyan.Println("\n" + strings.Repeat("=", 50))
	cyan.Printf("    STATS FOR: %s\n", strings.ToUpper(playerName))
	cyan.Println(strings.Repeat("=", 50))

	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	// Show spinner while loading detailed stats
	spinner2, _ := animations.StartSpinner("Loading detailed stats...")
	defer spinner2.Stop()

	// Get stats for VC mode
	vcStats, err := database.GetPlayerStatsByMode(playerName, "vc")
	if err != nil {
		color.Red("Error loading VC stats: %v", err)
		time.Sleep(2 * time.Second)
		return
	}

	// Get stats for Founder mode
	founderStats, err := database.GetPlayerStatsByMode(playerName, "founder")
	if err != nil {
		color.Red("Error loading Founder stats: %v", err)
		time.Sleep(2 * time.Second)
		return
	}

	if vcStats.TotalGames == 0 && founderStats.TotalGames == 0 {
		color.Yellow("\nNo games found for player: %s", playerName)
		fmt.Print("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// Display VC Mode Stats
	if vcStats.TotalGames > 0 {
		cyan.Println("\n🎩 VC INVESTOR MODE STATS:")
		cyan.Println(strings.Repeat("─", 50))

		fmt.Printf("\n%s Total Games Played: ", ascii.Chart)
		green.Printf("%d\n", vcStats.TotalGames)

		fmt.Printf("%s Best Net Worth: ", ascii.Money)
		green.Printf("$%s\n", FormatMoney(vcStats.BestNetWorth))

		fmt.Printf("%s Best ROI: ", ascii.Chart)
		green.Printf("%.2f%%\n", vcStats.BestROI)

		fmt.Printf("%s Total Successful Exits: ", ascii.Rocket)
		green.Printf("%d\n", vcStats.TotalExits)

		fmt.Printf("%s Average Net Worth: ", ascii.Coin)
		green.Printf("$%.0f\n", vcStats.AverageNetWorth)

		fmt.Printf("%s Win Rate (Positive ROI): ", ascii.Trophy)
		if vcStats.WinRate >= 50 {
			green.Printf("%.1f%%\n", vcStats.WinRate)
		} else {
			color.New(color.FgYellow).Printf("%.1f%%\n", vcStats.WinRate)
		}
	}

	// Display Founder Mode Stats
	if founderStats.TotalGames > 0 {
		cyan.Println("\n🚀 FOUNDER MODE STATS:")
		cyan.Println(strings.Repeat("─", 50))

		fmt.Printf("\n%s Total Games Played: ", ascii.Chart)
		green.Printf("%d\n", founderStats.TotalGames)

		fmt.Printf("%s Best Net Worth: ", ascii.Money)
		green.Printf("$%s\n", FormatMoney(founderStats.BestNetWorth))

		fmt.Printf("%s Best ROI: ", ascii.Chart)
		green.Printf("%.2f%%\n", founderStats.BestROI)

		fmt.Printf("%s Total Successful Exits: ", ascii.Rocket)
		green.Printf("%d\n", founderStats.TotalExits)

		fmt.Printf("%s Average Net Worth: ", ascii.Coin)
		green.Printf("$%.0f\n", founderStats.AverageNetWorth)

		fmt.Printf("%s Win Rate (Positive ROI): ", ascii.Trophy)
		if founderStats.WinRate >= 50 {
			green.Printf("%.1f%%\n", founderStats.WinRate)
		} else {
			color.New(color.FgYellow).Printf("%.1f%%\n", founderStats.WinRate)
		}
	}

	// Get and display active upgrades
	playerUpgrades, err := database.GetPlayerUpgrades(playerName)
	if err == nil && len(playerUpgrades) > 0 {
		fmt.Println()
		green.Println("✨ ACTIVE UPGRADES:")
		for _, upgradeID := range playerUpgrades {
			if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
				fmt.Printf("   %s %s - %s\n", upgrade.Icon, upgrade.Name, upgrade.Description)
			}
		}
	} else if err == nil && len(playerUpgrades) == 0 {
		fmt.Println()
		yellow.Printf("%s No upgrades purchased yet\n", ascii.Star)
		fmt.Println("   Purchase upgrades in the Upgrades menu to gain permanent advantages!")
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}