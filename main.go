package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	achievements "github.com/jamesacampbell/unicorn/achievements"
	ascii "github.com/jamesacampbell/unicorn/ascii"

	// analytics "github.com/jamesacampbell/unicorn/analytics"
	clear "github.com/jamesacampbell/unicorn/clear"
	db "github.com/jamesacampbell/unicorn/database"
	game "github.com/jamesacampbell/unicorn/game"
	leaderboard "github.com/jamesacampbell/unicorn/leaderboard"
	logo "github.com/jamesacampbell/unicorn/logo"
	yaml "gopkg.in/yaml.v2"
)

type gameData struct {
	Pot        int64  `yaml:"starting-cash"`
	BadThings  int64  `yaml:"number-of-bad-things-per-year"`
	Foreground string `yaml:"foreground-color"`
}

func initMenu() (username string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Name: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	fmt.Printf("\nWelcome %s!\n", text)
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

func askForAutomatedMode() bool {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("                 GAME MODE SELECTION")
	cyan.Println(strings.Repeat("=", 60))

	yellow.Println("\n1. Manual Mode (Press Enter each turn)")
	yellow.Println("2. Automated Mode (1 second per turn)")

	fmt.Print("\nEnter your choice (1-2, default 1): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	return choice == "2"
}

func displayWelcome(username string, difficulty game.Difficulty) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta, color.Bold)

	cyan.Printf("\n%s, welcome to your investment journey!\n", username)
	fmt.Printf("\nDifficulty: ")
	yellow.Printf("%s\n", difficulty.Name)
	fmt.Printf("Fund Size: $%s\n", formatMoney(difficulty.StartingCash))
	fmt.Printf("Management Fee: 2%% annually ($%s/year)\n", formatMoney(int64(float64(difficulty.StartingCash)*0.02)))
	fmt.Printf("Game Duration: %d turns (%d years)\n", difficulty.MaxTurns, difficulty.MaxTurns/12)

	fmt.Println("\nEach turn = 1 month. Choose your investments wisely!")
	fmt.Println("Random events and funding rounds will affect valuations.")
	fmt.Println("Watch out for dilution when companies raise new rounds!")
	
	magenta.Println("\n?? COMPETING AGAINST:")
	fmt.Println("   ? CARL (Sterling & Cooper) - Conservative")
	fmt.Println("   ? Sarah Chen (Accel Partners) - Aggressive")
	fmt.Println("   ? Marcus Williams (Sequoia Capital) - Balanced")

	fmt.Print("\nPress 'Enter' to see available startups...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func displayStartup(s game.Startup, index int) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Printf("\n[%d] %s\n", index+1, s.Name)
	fmt.Printf("    %s\n", s.Description)
	yellow.Printf("    Category: %s\n", s.Category)
	fmt.Printf("    Valuation: $%s\n", formatMoney(s.Valuation))
	fmt.Printf("    Monthly Sales: %d units\n", s.MonthlySales)
	fmt.Printf("    Margin: %d%%\n", s.PercentMargin)
	fmt.Printf("    Website Visitors: %s/month\n", formatNumber(s.MonthlyWebsiteVisitors))

	// Risk indicator
	riskColor := color.New(color.FgGreen)
	riskLabel := "Low"
	if s.RiskScore > 0.85 {
		riskColor = color.New(color.FgRed, color.Bold)
		riskLabel = "VERY HIGH"
	} else if s.RiskScore > 0.6 {
		riskColor = color.New(color.FgRed)
		riskLabel = "High"
	} else if s.RiskScore > 0.4 {
		riskColor = color.New(color.FgYellow)
		riskLabel = "Medium"
	} else if s.RiskScore < 0.3 {
		riskColor = color.New(color.FgGreen, color.Bold)
		riskLabel = "LOW"
	}
	riskColor.Printf("    Risk: %s", riskLabel)

	// Growth indicator
	growthColor := color.New(color.FgGreen)
	growthLabel := "High"
	if s.GrowthPotential > 0.85 {
		growthColor = color.New(color.FgGreen, color.Bold)
		growthLabel = "VERY HIGH"
	} else if s.GrowthPotential < 0.4 {
		growthColor = color.New(color.FgRed)
		growthLabel = "Low"
	} else if s.GrowthPotential < 0.6 {
		growthColor = color.New(color.FgYellow)
		growthLabel = "Medium"
	}
	growthColor.Printf(" | Growth Potential: %s\n", growthLabel)
}

func investmentPhase(gs *game.GameState) {
	clear.ClearIt()
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Print(ascii.InvestmentHeader)
	green.Printf("\nTurn %d/%d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	fmt.Printf("Fund Size: $%s\n", formatMoney(gs.Portfolio.InitialFundSize))
	fmt.Printf("Cash Available: $%s\n", formatMoney(gs.Portfolio.Cash))
	fmt.Printf("Portfolio Value: $%s\n", formatMoney(gs.GetPortfolioValue()))
	fmt.Printf("Net Worth: $%s\n", formatMoney(gs.Portfolio.NetWorth))

	// Calculate and display total company valuation
	totalValuation := int64(0)
	for _, startup := range gs.AvailableStartups {
		totalValuation += startup.Valuation
	}
	yellow := color.New(color.FgYellow)
	yellow.Printf("Total Company Valuation: $%s\n", formatMoney(totalValuation))

	// Show available startups
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("AVAILABLE STARTUPS:")
	for i, startup := range gs.AvailableStartups {
		displayStartup(startup, i)
	}
	fmt.Println(strings.Repeat("=", 50))

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\nEnter company number (1-%d) to invest, or 'done' to continue: ", len(gs.AvailableStartups))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "done" || input == "" {
			// Have AI players make their investments too
			gs.AIPlayerMakeInvestments()
			color.Green("\n? AI players have made their investments!")
			break
		}

		companyNum, err := strconv.Atoi(input)
		if err != nil || companyNum < 1 || companyNum > len(gs.AvailableStartups) {
			color.Red("Invalid company number!")
			continue
		}

		fmt.Printf("Enter investment amount ($): ")
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)
		amount, err := strconv.ParseInt(amountStr, 10, 64)

		if err != nil {
			color.Red("Invalid amount!")
			continue
		}

		err = gs.MakeInvestment(companyNum-1, amount)
		if err != nil {
			color.Red("Error: %v", err)
		} else {
			color.Green("%s Investment successful!", ascii.Check)
			fmt.Printf("Cash remaining: $%s\n", formatMoney(gs.Portfolio.Cash))
		}
	}
}

func playTurn(gs *game.GameState, autoMode bool) {
	yellow := color.New(color.FgYellow, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)

	// Print separator line instead of clearing screen
	fmt.Println(strings.Repeat("=", 70))
	yellow.Printf("\n%s MONTH %d of %d\n", ascii.Calendar, gs.Portfolio.Turn, gs.Portfolio.MaxTurns)

	messages := gs.ProcessTurn()

	if len(messages) > 0 {
		cyan.Print(ascii.NewsHeader)
		for _, msg := range messages {
			fmt.Println(msg)
		}
		fmt.Println() // Add spacing after news
	}

	// Show portfolio status
	cyan.Print(ascii.PortfolioHeader)
	if len(gs.Portfolio.Investments) == 0 {
		fmt.Println("   No investments yet")
	} else {
		// Calculate total company valuation
		totalCompanyValuation := int64(0)
		for _, inv := range gs.Portfolio.Investments {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			profit := value - inv.AmountInvested
			profitColor := color.New(color.FgGreen)
			profitSign := "+"
			if profit < 0 {
				profitColor = color.New(color.FgRed)
				profitSign = ""
			}

			// Show dilution if applicable
			dilutionInfo := ""
			if len(inv.Rounds) > 0 {
				dilutionInfo = fmt.Sprintf(" (was %.2f%%, %d rounds)", inv.InitialEquity, len(inv.Rounds))
			}
			
			fmt.Printf("   %s: $%s invested, %.2f%% equity%s\n",
				inv.CompanyName, formatMoney(inv.AmountInvested), inv.EquityPercent, dilutionInfo)
			fmt.Printf("      Current Value: $%s ", formatMoney(value))
			profitColor.Printf("(%s$%s)\n", profitSign, formatMoney(abs(profit)))

			totalCompanyValuation += inv.CurrentValuation
		}

		// Display total company valuation
		yellow.Printf("\n   Total Company Valuation: $%s\n", formatMoney(totalCompanyValuation))
	}

	fmt.Printf("\n%s Net Worth: $%s", ascii.Money, formatMoney(gs.Portfolio.NetWorth))
	fmt.Printf(" | Management Fees Paid: $%s\n", formatMoney(gs.Portfolio.ManagementFeesCharged))
	
	// Show competitive leaderboard every quarter
	if gs.Portfolio.Turn%3 == 0 {
		displayMiniLeaderboard(gs)
	}

	// Always wait for user input when there are news messages, otherwise use auto mode logic
	if len(messages) > 0 {
		fmt.Print("\nPress 'Enter' to continue to next month...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	} else if autoMode {
		time.Sleep(1 * time.Second)
	} else {
		fmt.Print("\nPress 'Enter' to continue to next month...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func displayFinalScore(gs *game.GameState) {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)

	cyan.Print(ascii.GameOverHeader)

	netWorth, roi, successfulExits := gs.GetFinalScore()

	fmt.Printf("\n%s Player: %s\n", ascii.Star, gs.PlayerName)
	fmt.Printf("%s Turns Played: %d\n", ascii.Calendar, gs.Portfolio.Turn-1)
	fmt.Printf("%s Management Fees Paid: $%s\n\n", ascii.Money, formatMoney(gs.Portfolio.ManagementFeesCharged))

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("%s Final Net Worth: $%s\n", ascii.Money, formatMoney(netWorth))

	roiColor := color.New(color.FgGreen)
	if roi < 0 {
		roiColor = color.New(color.FgRed)
	}
	roiColor.Printf("%s Return on Investment: %.2f%%\n", ascii.Chart, roi)
	fmt.Printf("%s Successful Exits (5x+): %d\n", ascii.Rocket, successfulExits)
	
	// Show competitive results
	magenta.Println("\n" + strings.Repeat("=", 70))
	magenta.Println("                   FINAL LEADERBOARD")
	magenta.Println(strings.Repeat("=", 70))
	
	leaderboard := gs.GetLeaderboard()
	fmt.Printf("\n%-5s %-25s %-25s %-15s %-10s\n", "RANK", "INVESTOR", "FIRM", "NET WORTH", "ROI")
	fmt.Println(strings.Repeat("-", 90))
	
	for i, entry := range leaderboard {
		rankColor := color.New(color.FgWhite)
		if i == 0 {
			rankColor = color.New(color.FgYellow, color.Bold)
		} else if i == 1 {
			rankColor = color.New(color.FgCyan)
		} else if i == 2 {
			rankColor = color.New(color.FgGreen)
		}
		
		playerMarker := ""
		if entry.IsPlayer {
			playerMarker = " ? YOU"
		}
		
		roiColorEntry := color.New(color.FgGreen)
		if entry.ROI < 0 {
			roiColorEntry = color.New(color.FgRed)
		}
		
		rankColor.Printf("%-5d ", i+1)
		fmt.Printf("%-25s ", entry.Name+playerMarker)
		fmt.Printf("%-25s ", entry.Firm)
		fmt.Printf("$%-14s ", formatMoney(entry.NetWorth))
		roiColorEntry.Printf("%.1f%%\n", entry.ROI)
	}
	
	if leaderboard[0].IsPlayer {
		magenta.Println("\n?? CONGRATULATIONS! You beat all the AI investors!")
	} else {
		magenta.Printf("\nYou finished in position #%d. Better luck next time!\n", findPlayerRank(leaderboard))
	}

	fmt.Println("\n" + strings.Repeat("?", 50))
	fmt.Println("FINAL PORTFOLIO:")
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		fmt.Printf("   %s: $%s ? $%s\n",
			inv.CompanyName, formatMoney(inv.AmountInvested), formatMoney(value))
	}

	// Performance rating
	fmt.Println("\n" + strings.Repeat("?", 50))
	var rating string
	var icon string
	if roi >= 1000 {
		rating = "UNICORN HUNTER - Legendary!"
		icon = ascii.Crown
	} else if roi >= 500 {
		rating = "Elite VC - Outstanding!"
		icon = ascii.Trophy
	} else if roi >= 200 {
		rating = "Great Investor - Excellent!"
		icon = ascii.Star
	} else if roi >= 50 {
		rating = "Solid Performance - Good!"
		icon = ascii.Check
	} else if roi >= 0 {
		rating = "Break Even - Not Bad"
		icon = "="
	} else {
		rating = "Lost Money - Better Luck Next Time"
		icon = ascii.Warning
	}

	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("Rating: %s %s\n", icon, rating)
	fmt.Println(strings.Repeat("=", 70) + "\n")
}

func displayMiniLeaderboard(gs *game.GameState) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Println("\n?? Current Standings:")
	leaderboard := gs.GetLeaderboard()
	
	for i, entry := range leaderboard {
		marker := "  "
		if entry.IsPlayer {
			marker = "? "
			yellow.Printf("%s%d. %s (%s): $%s (ROI: %.1f%%)\n",
				marker, i+1, entry.Name, entry.Firm, formatMoney(entry.NetWorth), entry.ROI)
		} else {
			fmt.Printf("%s%d. %s (%s): $%s (ROI: %.1f%%)\n",
				marker, i+1, entry.Name, entry.Firm, formatMoney(entry.NetWorth), entry.ROI)
		}
	}
}

func findPlayerRank(leaderboard []game.PlayerScore) int {
	for i, entry := range leaderboard {
		if entry.IsPlayer {
			return i + 1
		}
	}
	return len(leaderboard)
}

func main() {
	// Initialize database
	err := db.InitDB("unicorn_scores.db")
	if err != nil {
		fmt.Printf("Warning: Could not initialize database: %v\n", err)
		fmt.Println("Scores will not be saved.")
		time.Sleep(2 * time.Second)
	}
	defer db.CloseDB()

	// Main menu loop
	for {
		choice := displayMainMenu()
		clear.ClearIt()

		switch choice {
		case "1":
			playNewGame()
		case "2":
			displayLeaderboards()
		case "3":
			displayPlayerStats()
		case "4":
			displayAchievementsMenu()
		case "5":
			displayHelpGuide()
		case "6":
			fmt.Println("\nThanks for playing! " + ascii.Star2)
			return
		default:
			color.Red("Invalid choice!")
			time.Sleep(1 * time.Second)
		}
	}
}

func displayMainMenu() string {
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
	yellow.Println("5. Help & Info")
	yellow.Println("6. Quit")

	fmt.Print("\nEnter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

func selectDifficulty() game.Difficulty {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	magenta := color.New(color.FgMagenta)

	cyan.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("                 SELECT DIFFICULTY")
	cyan.Println(strings.Repeat("=", 60))

	green.Printf("\n1. Easy")
	fmt.Printf(" - %s\n", game.EasyDifficulty.Description)
	fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
		formatMoney(game.EasyDifficulty.StartingCash), game.EasyDifficulty.MaxTurns)

	yellow.Printf("\n2. Medium")
	fmt.Printf(" - %s\n", game.MediumDifficulty.Description)
	fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
		formatMoney(game.MediumDifficulty.StartingCash), game.MediumDifficulty.MaxTurns)

	red.Printf("\n3. Hard")
	fmt.Printf(" - %s\n", game.HardDifficulty.Description)
	fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
		formatMoney(game.HardDifficulty.StartingCash), game.HardDifficulty.MaxTurns)

	magenta.Printf("\n4. Expert")
	fmt.Printf(" - %s\n", game.ExpertDifficulty.Description)
	fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
		formatMoney(game.ExpertDifficulty.StartingCash), game.ExpertDifficulty.MaxTurns)

	fmt.Print("\nEnter your choice (1-4): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return game.EasyDifficulty
	case "3":
		return game.HardDifficulty
	case "4":
		return game.ExpertDifficulty
	default:
		return game.MediumDifficulty
	}
}

func playNewGame() {
	// Get username
	username := initMenu()
	clear.ClearIt()

	// Select difficulty
	difficulty := selectDifficulty()
	clear.ClearIt()

	// Ask for automated mode
	autoMode := askForAutomatedMode()
	clear.ClearIt()

	// Display welcome and rules
	displayWelcome(username, difficulty)

	// Initialize game
	gs := game.NewGame(username, difficulty)

	// Investment phase at start
	investmentPhase(gs)

	// Main game loop
	for !gs.IsGameOver() {
		playTurn(gs, autoMode)
	}

	// Show final score
	displayFinalScore(gs)

	// Save score to database
	netWorth, roi, successfulExits := gs.GetFinalScore()
	score := db.GameScore{
		PlayerName:      gs.PlayerName,
		FinalNetWorth:   netWorth,
		ROI:             roi,
		SuccessfulExits: successfulExits,
		TurnsPlayed:     gs.Portfolio.Turn - 1,
		Difficulty:      gs.Difficulty.Name,
		PlayedAt:        time.Now(),
	}

	err := db.SaveGameScore(score)
	if err != nil {
		color.Yellow("\nWarning: Could not save score: %v", err)
	} else {
		color.Green("\n%s Score saved to local leaderboard!", ascii.Check)
	}

	// Ask if player wants to submit to global leaderboard
	askToSubmitToGlobalLeaderboard(score)

	// Check for achievements
	checkAndUnlockAchievements(gs)

	fmt.Print("\nPress 'Enter' to return to main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func checkAndUnlockAchievements(gs *game.GameState) {
	// Get player's previously unlocked achievements
	previouslyUnlocked, err := db.GetPlayerAchievements(gs.PlayerName)
	if err != nil {
		previouslyUnlocked = []string{}
	}

	// Get player stats
	stats, _ := db.GetPlayerStats(gs.PlayerName)
	winStreak, _ := db.GetWinStreak(gs.PlayerName)

	// Count sectors and get investment details
	sectors := make(map[string]bool)
	positiveCount := 0
	negativeCount := 0
	totalInvested := int64(0)

	for _, inv := range gs.Portfolio.Investments {
		totalInvested += inv.AmountInvested
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		if value > inv.AmountInvested {
			positiveCount++
		} else if value < inv.AmountInvested {
			negativeCount++
		}

		// Find sector
		for _, startup := range gs.AvailableStartups {
			if startup.Name == inv.CompanyName {
				sectors[startup.Category] = true
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

	// Save and display new achievements
	if len(newAchievements) > 0 {
		cyan := color.New(color.FgCyan, color.Bold)
		yellow := color.New(color.FgYellow)

		fmt.Println("\n" + strings.Repeat("?", 60))
		cyan.Printf("     %s NEW ACHIEVEMENTS UNLOCKED! %s\n", ascii.Star, ascii.Star)
		fmt.Println(strings.Repeat("?", 60))

		for _, ach := range newAchievements {
			// Save to database
			db.UnlockAchievement(gs.PlayerName, ach.ID)

			// Display
			rarityColor := color.New(color.Attribute(achievements.GetRarityColor(ach.Rarity)))
			fmt.Printf("\n%s  ", ach.Icon)
			rarityColor.Printf("%s", ach.Name)
			fmt.Printf(" [%s]\n", ach.Rarity)
			yellow.Printf("   %s\n", ach.Description)
			fmt.Printf("   +%d points\n", ach.Points)
		}

		// Calculate new career level
		totalPoints := 0
		allUnlocked, _ := db.GetPlayerAchievements(gs.PlayerName)
		for _, id := range allUnlocked {
			if ach, exists := achievements.AllAchievements[id]; exists {
				totalPoints += ach.Points
			}
		}

		level, title, _ := achievements.CalculateCareerLevel(totalPoints)

		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Printf("Career Level: ")
		yellow.Printf("%d - %s", level, title)
		fmt.Printf(" | Total Points: ")
		yellow.Printf("%d\n", totalPoints)
		fmt.Println(strings.Repeat("=", 60))
	}
}

func displayLeaderboards() {
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
		showTopScores("net_worth", "all")
	case "2":
		showTopScores("roi", "all")
	case "3":
		showTopScores("net_worth", "Easy")
	case "4":
		showTopScores("net_worth", "Medium")
	case "5":
		showTopScores("net_worth", "Hard")
	case "6":
		showTopScores("net_worth", "Expert")
	case "7":
		showRecentGames()
	case "8":
		return
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	displayLeaderboards()
}

func showTopScores(sortBy string, difficulty string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	var scores []db.GameScore
	var err error
	var title string

	if sortBy == "roi" {
		scores, err = db.GetTopScoresByROI(10, difficulty)
		title = "TOP 10 BY ROI"
	} else {
		scores, err = db.GetTopScoresByNetWorth(10, difficulty)
		title = "TOP 10 BY NET WORTH"
	}

	if difficulty != "all" && difficulty != "" {
		title += fmt.Sprintf(" (%s)", strings.ToUpper(difficulty))
	}

	if err != nil {
		color.Red("Error loading leaderboard: %v", err)
		return
	}

	cyan.Println("\n" + strings.Repeat("=", 90))
	cyan.Printf("%-40s\n", title)
	cyan.Println(strings.Repeat("=", 90))

	if len(scores) == 0 {
		yellow.Println("\nNo games played yet! Be the first!")
		return
	}

	fmt.Printf("\n%-5s %-20s %-15s %-15s %-10s %-12s\n",
		"RANK", "PLAYER", "NET WORTH", "ROI", "EXITS", "DIFFICULTY")
	fmt.Println(strings.Repeat("-", 90))

	for i, score := range scores {
		rankColor := color.New(color.FgWhite)
		if i == 0 {
			rankColor = color.New(color.FgYellow, color.Bold)
		} else if i == 1 {
			rankColor = color.New(color.FgCyan)
		} else if i == 2 {
			rankColor = color.New(color.FgGreen)
		}

		roiColor := color.New(color.FgGreen)
		if score.ROI < 0 {
			roiColor = color.New(color.FgRed)
		}

		rankColor.Printf("%-5d ", i+1)
		fmt.Printf("%-20s ", truncateString(score.PlayerName, 20))
		fmt.Printf("$%-14s ", formatMoney(score.FinalNetWorth))
		roiColor.Printf("%-15s ", fmt.Sprintf("%.1f%%", score.ROI))
		fmt.Printf("%-10d ", score.SuccessfulExits)
		fmt.Printf("%-12s\n", score.Difficulty)
	}
}

func showRecentGames() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	scores, err := db.GetRecentGames(10)
	if err != nil {
		color.Red("Error loading recent games: %v", err)
		return
	}

	cyan.Println("\n" + strings.Repeat("=", 90))
	cyan.Println("                           RECENT GAMES")
	cyan.Println(strings.Repeat("=", 90))

	if len(scores) == 0 {
		yellow.Println("\nNo games played yet!")
		return
	}

	fmt.Printf("\n%-20s %-15s %-15s %-12s %-20s\n",
		"PLAYER", "NET WORTH", "ROI", "DIFFICULTY", "DATE")
	fmt.Println(strings.Repeat("-", 90))

	for _, score := range scores {
		roiColor := color.New(color.FgGreen)
		if score.ROI < 0 {
			roiColor = color.New(color.FgRed)
		}

		fmt.Printf("%-20s ", truncateString(score.PlayerName, 20))
		fmt.Printf("$%-14s ", formatMoney(score.FinalNetWorth))
		roiColor.Printf("%-15s ", fmt.Sprintf("%.1f%%", score.ROI))
		fmt.Printf("%-12s ", score.Difficulty)
		fmt.Printf("%-20s\n", score.PlayedAt.Format("2006-01-02 15:04"))
	}
}

func displayPlayerStats() {
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

	stats, err := db.GetPlayerStats(playerName)
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

	fmt.Printf("\n%s Total Games Played: ", ascii.Chart)
	green.Printf("%d\n", stats.TotalGames)

	fmt.Printf("%s Best Net Worth: ", ascii.Money)
	green.Printf("$%s\n", formatMoney(stats.BestNetWorth))

	fmt.Printf("%s Best ROI: ", ascii.Chart)
	green.Printf("%.2f%%\n", stats.BestROI)

	fmt.Printf("%s Total Successful Exits: ", ascii.Rocket)
	green.Printf("%d\n", stats.TotalExits)

	fmt.Printf("%s Average Net Worth: ", ascii.Coin)
	green.Printf("$%.0f\n", stats.AverageNetWorth)

	fmt.Printf("%s Win Rate (Positive ROI): ", ascii.Trophy)
	if stats.WinRate >= 50 {
		green.Printf("%.1f%%\n", stats.WinRate)
	} else {
		color.New(color.FgYellow).Printf("%.1f%%\n", stats.WinRate)
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Helper functions
func formatMoney(amount int64) string {
	abs := amount
	if abs < 0 {
		abs = -abs
	}

	s := strconv.FormatInt(abs, 10)

	// Add commas
	n := len(s)
	if n <= 3 {
		if amount < 0 {
			return "-" + s
		}
		return s
	}

	result := ""
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	if amount < 0 {
		return "-" + result
	}
	return result
}

func formatNumber(n int) string {
	return formatMoney(int64(n))
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func displayAchievementsMenu() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)

	cyan.Print(ascii.AchievementsHeader)

	fmt.Println("\n1. View My Achievements")
	fmt.Println("2. Browse All Achievements")
	fmt.Println("3. Leaderboard (Most Achievements)")
	fmt.Println("4. Back to Main Menu")

	fmt.Print("\nEnter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	clear.ClearIt()

	switch choice {
	case "1":
		viewPlayerAchievements()
	case "2":
		browseAllAchievements()
	case "3":
		displayAchievementLeaderboard()
	case "4":
		return
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	displayAchievementsMenu()
}

func viewPlayerAchievements() {
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

	unlocked, err := db.GetPlayerAchievements(playerName)
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

func browseAllAchievements() {
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

func displayAchievementLeaderboard() {
	cyan := color.New(color.FgCyan, color.Bold)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("            ACHIEVEMENT LEADERBOARD (Coming Soon)")
	cyan.Println(strings.Repeat("=", 70))

	color.Yellow("\nThis feature will show players with the most achievements!")
}

func askToSubmitToGlobalLeaderboard(score db.GameScore) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("           ?? GLOBAL LEADERBOARD")
	fmt.Println(strings.Repeat("=", 60))

	yellow.Println("\nWould you like to submit your score to the global leaderboard?")
	fmt.Println("Your score will be visible to all players worldwide!")
	fmt.Print("\nSubmit to global leaderboard? (y/n, default n): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice != "y" && choice != "yes" {
		fmt.Println("Okay, score saved locally only.")
		return
	}

	// Check API availability
	fmt.Print("\nChecking global leaderboard service...")
	if !leaderboard.IsAPIAvailable("") {
		color.Yellow("\n??  Global leaderboard service is not available right now.")
		color.Yellow("Your score has been saved locally.")
		return
	}
	color.Green(" ?")

	// Submit score
	fmt.Print("Submitting your score...")
	submission := leaderboard.ScoreSubmission{
		PlayerName:      score.PlayerName,
		FinalNetWorth:   score.FinalNetWorth,
		ROI:             score.ROI,
		SuccessfulExits: score.SuccessfulExits,
		TurnsPlayed:     score.TurnsPlayed,
		Difficulty:      score.Difficulty,
	}

	err := leaderboard.SubmitScore(submission, "")
	if err != nil {
		color.Red("\n? Failed to submit score: %v", err)
		color.Yellow("Your score has been saved locally.")
		return
	}

	color.Green(" ?")
	cyan.Println("\n?? Success! Your score has been submitted to the global leaderboard!")
	yellow.Println("\nView the global leaderboard at:")
	yellow.Println("https://jamesacampbell.github.io/unicorn")
}

func displayHelpGuide() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n" + strings.Repeat("?", 70))
	cyan.Println("              HELP & INFORMATION")
	cyan.Println(strings.Repeat("?", 70))

	yellow.Printf("\n%s GAME OVERVIEW\n", ascii.Lightbulb)
	fmt.Println("You're a VC investor with limited capital. Invest in 15 randomly")
	fmt.Println("selected startups from a pool of 30 and watch your portfolio")
	fmt.Println("grow (or shrink) over 10 years.")

	yellow.Printf("\n%s HOW TO PLAY\n", ascii.Target)
	fmt.Println("1. Select difficulty (Easy/Medium/Hard/Expert)")
	fmt.Println("2. Review 15 randomly selected companies (from 30 total)")
	fmt.Println("3. Invest your capital across multiple startups")
	fmt.Println("4. Watch events unfold each turn (1 turn = 1 month)")
	fmt.Println("5. After 90-120 turns, see your final score")

	yellow.Printf("\n%s COMPANY METRICS\n", ascii.Building)
	fmt.Printf("%s Risk Score: Low/Medium/High - chance of failure\n", ascii.Warning)
	fmt.Printf("%s Growth Potential: Projected growth trajectory\n", ascii.Chart)
	fmt.Printf("%s Valuation: Current company worth (in millions)\n", ascii.Money)
	fmt.Printf("%s Category: Industry sector (FinTech, BioTech, etc.)\n", ascii.Star)

	yellow.Printf("\n%s SCORING\n", ascii.Trophy)
	fmt.Printf("%s Net Worth: Cash + Portfolio Value\n", ascii.Money)
	fmt.Printf("%s ROI: Return on Investment percentage\n", ascii.Chart)
	fmt.Printf("%s Successful Exits: Companies that 5x or more\n", ascii.Rocket)
	fmt.Printf("%s Rating: Based on ROI (Unicorn Hunter = 1000%%+)\n", ascii.Crown)

	yellow.Printf("\n%s DIFFICULTY LEVELS\n", ascii.Shield)
	fmt.Printf("%s Easy: $1M fund, 20%% events, 3%% volatility\n", ascii.Check)
	fmt.Printf("%s Medium: $750k fund, 30%% events, 5%% volatility\n", ascii.Star)
	fmt.Printf("%s Hard: $500k fund, 40%% events, 7%% volatility\n", ascii.Warning)
	fmt.Printf("%s Expert: $500k fund, 50%% events, 10%% volatility, 90 turns\n", ascii.Zap)

	yellow.Printf("\n%s ANALYTICS\n", ascii.Chart)
	fmt.Println("After each game, view detailed portfolio analytics:")
	fmt.Printf("%s Best/Worst performers\n", ascii.Medal)
	fmt.Printf("%s Sector breakdown\n", ascii.Building)
	fmt.Printf("%s Win/loss ratio\n", ascii.Trophy)
	fmt.Printf("%s Investment distribution\n", ascii.Portfolio)

	yellow.Printf("\n%s AVAILABLE COMPANIES\n", ascii.Building)
	fmt.Println("30 diverse startups across 12+ sectors:")
	fmt.Println("Each game randomly selects 15 companies from the pool.")
	fmt.Println("FinTech ? BioTech ? CleanTech ? HealthTech ? EdTech")
	fmt.Println("Robotics ? Security ? Gaming ? LegalTech ? AgriTech")
	fmt.Println("Logistics ? IoT ? Creative ? CloudTech ? and more!")
	fmt.Println("Includes LOW risk stable companies and VERY HIGH risk moonshots!")

	yellow.Printf("\n%s RANDOM EVENTS\n", ascii.News)
	fmt.Println("60+ possible events can affect your companies:")
	fmt.Printf("%s Funding rounds (Series A/B, IPO)\n", ascii.Money)
	fmt.Printf("%s Product launches (success/failure)\n", ascii.Rocket)
	fmt.Printf("%s Partnerships & acquisitions\n", ascii.Building)
	fmt.Printf("%s Scandals & regulatory issues\n", ascii.Warning)
	fmt.Printf("%s Market conditions & competition\n", ascii.Chart)

	yellow.Printf("\n%s REALISTIC VC MECHANICS\n", ascii.Money)
	fmt.Printf("%s Management Fees: 2%% annual fee charged monthly\n", ascii.Coin)
	fmt.Printf("%s Multiple Rounds: Companies raise Seed, Series A/B/C\n", ascii.Rocket)
	fmt.Printf("%s Dilution: Your equity %% decreases with new rounds\n", ascii.Warning)
	fmt.Printf("%s Post-Money Valuation: Pre-money + investment\n", ascii.Chart)
	fmt.Printf("%s AI Competition: Play against 3 AI VCs\n", ascii.Trophy)
	fmt.Println("   \u2022 CARL (Sterling & Cooper) - Conservative")
	fmt.Println("   \u2022 Sarah Chen (Accel) - Aggressive")
	fmt.Println("   \u2022 Marcus Williams (Sequoia) - Balanced")
	
	yellow.Printf("\n%s STRATEGY TIPS\n", ascii.Lightbulb)
	fmt.Printf("%s Diversify: Don't put everything in one company\n", ascii.Portfolio)
	fmt.Printf("%s Balance: Mix high-risk and low-risk investments\n", ascii.Shield)
	fmt.Printf("%s Sectors: Different industries perform differently\n", ascii.Building)
	fmt.Printf("%s Research: Read company metrics carefully\n", ascii.Star)
	fmt.Printf("%s Early Entry: Invest before dilution hits\n", ascii.Target)
	fmt.Printf("%s Management Fees: Factor in 2%% annual costs\n", ascii.Coin)

	cyan.Println("\n" + strings.Repeat("=", 70))
	fmt.Print("\nPress 'Enter' to return to menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
