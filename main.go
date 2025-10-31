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
	// analytics "github.com/jamesacampbell/unicorn/analytics"
	clear "github.com/jamesacampbell/unicorn/clear"
	db "github.com/jamesacampbell/unicorn/database"
	game "github.com/jamesacampbell/unicorn/game"
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
	
	cyan.Printf("\n%s, welcome to your investment journey!\n", username)
	fmt.Printf("\nDifficulty: ")
	yellow.Printf("%s\n", difficulty.Name)
	fmt.Printf("Starting Cash: $%s\n", formatMoney(difficulty.StartingCash))
	fmt.Printf("Game Duration: %d turns (%d years)\n", difficulty.MaxTurns, difficulty.MaxTurns/12)
	
	fmt.Println("\nEach turn = 1 month. Choose your investments wisely!")
	fmt.Println("Random events will affect valuations throughout the game.")
	
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
	if s.RiskScore > 0.6 {
		riskColor = color.New(color.FgRed)
		riskLabel = "High"
	} else if s.RiskScore > 0.4 {
		riskColor = color.New(color.FgYellow)
		riskLabel = "Medium"
	}
	riskColor.Printf("    Risk: %s", riskLabel)
	
	// Growth indicator
	growthColor := color.New(color.FgGreen)
	growthLabel := "High"
	if s.GrowthPotential < 0.4 {
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
	green.Printf("\n?? INVESTMENT PHASE - Turn %d/%d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	fmt.Printf("Cash Available: $%s\n", formatMoney(gs.Portfolio.Cash))
	fmt.Printf("Portfolio Value: $%s\n", formatMoney(gs.GetPortfolioValue()))
	fmt.Printf("Net Worth: $%s\n", formatMoney(gs.Portfolio.NetWorth))
	
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
			color.Green("? Investment successful!")
			fmt.Printf("Cash remaining: $%s\n", formatMoney(gs.Portfolio.Cash))
		}
	}
}

func playTurn(gs *game.GameState, autoMode bool) {
	clear.ClearIt()
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("\n?? MONTH %d of %d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	
	messages := gs.ProcessTurn()
	
	if len(messages) > 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("COMPANY NEWS:")
		for _, msg := range messages {
			fmt.Println(msg)
		}
		fmt.Println(strings.Repeat("=", 50))
	}
	
	// Show portfolio status
	fmt.Println("\n?? YOUR PORTFOLIO:")
	if len(gs.Portfolio.Investments) == 0 {
		fmt.Println("   No investments yet")
	} else {
		for _, inv := range gs.Portfolio.Investments {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			profit := value - inv.AmountInvested
			profitColor := color.New(color.FgGreen)
			profitSign := "+"
			if profit < 0 {
				profitColor = color.New(color.FgRed)
				profitSign = ""
			}
			
			fmt.Printf("   %s: $%s invested, %.2f%% equity\n", 
				inv.CompanyName, formatMoney(inv.AmountInvested), inv.EquityPercent)
			fmt.Printf("      Current Value: $%s ", formatMoney(value))
			profitColor.Printf("(%s$%s)\n", profitSign, formatMoney(abs(profit)))
		}
	}
	
	fmt.Printf("\n?? Net Worth: $%s\n", formatMoney(gs.Portfolio.NetWorth))
	
	if autoMode {
		time.Sleep(1 * time.Second)
	} else {
		fmt.Print("\nPress 'Enter' to continue to next month...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func displayFinalScore(gs *game.GameState) {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	
	cyan.Println("\n" + strings.Repeat("=", 50))
	cyan.Println("           ?? GAME OVER - FINAL RESULTS ??")
	cyan.Println(strings.Repeat("=", 50))
	
	netWorth, roi, successfulExits := gs.GetFinalScore()
	
	fmt.Printf("\n?? Player: %s\n", gs.PlayerName)
	fmt.Printf("?? Turns Played: %d\n\n", gs.Portfolio.Turn-1)
	
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("?? Final Net Worth: $%s\n", formatMoney(netWorth))
	
	roiColor := color.New(color.FgGreen)
	if roi < 0 {
		roiColor = color.New(color.FgRed)
	}
	roiColor.Printf("?? Return on Investment: %.2f%%\n", roi)
	fmt.Printf("?? Successful Exits (5x+): %d\n", successfulExits)
	
	fmt.Println("\n?? FINAL PORTFOLIO:")
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		fmt.Printf("   %s: $%s ? $%s\n", 
			inv.CompanyName, formatMoney(inv.AmountInvested), formatMoney(value))
	}
	
	// Performance rating
	fmt.Println("\n" + strings.Repeat("=", 50))
	var rating string
	if roi >= 1000 {
		rating = "?? UNICORN HUNTER - Legendary!"
	} else if roi >= 500 {
		rating = "?? Elite VC - Outstanding!"
	} else if roi >= 200 {
		rating = "? Great Investor - Excellent!"
	} else if roi >= 50 {
		rating = "?? Solid Performance - Good!"
	} else if roi >= 0 {
		rating = "?? Break Even - Not Bad"
	} else {
		rating = "?? Lost Money - Better Luck Next Time"
	}
	
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("Rating: %s\n", rating)
	fmt.Println(strings.Repeat("=", 50) + "\n")
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
	
	// Init with the unicorn logo
	c := color.New(color.FgCyan).Add(color.Bold)
	logo.InitLogo(c)
	
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
			fmt.Println("\nThanks for playing! ??")
			return
		default:
			color.Red("Invalid choice!")
			time.Sleep(1 * time.Second)
		}
	}
}

func displayMainMenu() string {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Println("\n" + strings.Repeat("=", 50))
	cyan.Println("           ?? UNICORN - MAIN MENU ??")
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
		color.Green("\n? Score saved to leaderboard!")
	}
	
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
		
		fmt.Println("\n" + strings.Repeat("=", 60))
		cyan.Println("           ?? NEW ACHIEVEMENTS UNLOCKED! ??")
		fmt.Println(strings.Repeat("=", 60))
		
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
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                       ?? LEADERBOARDS ??")
	cyan.Println(strings.Repeat("=", 70))
	
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
	
	cyan.Println("\n" + strings.Repeat("=", 50))
	cyan.Println("              PLAYER STATISTICS")
	cyan.Println(strings.Repeat("=", 50))
	
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
	
	fmt.Printf("\n?? Total Games Played: ")
	green.Printf("%d\n", stats.TotalGames)
	
	fmt.Printf("?? Best Net Worth: ")
	green.Printf("$%s\n", formatMoney(stats.BestNetWorth))
	
	fmt.Printf("?? Best ROI: ")
	green.Printf("%.2f%%\n", stats.BestROI)
	
	fmt.Printf("?? Total Successful Exits: ")
	green.Printf("%d\n", stats.TotalExits)
	
	fmt.Printf("?? Average Net Worth: ")
	green.Printf("$%.0f\n", stats.AverageNetWorth)
	
	fmt.Printf("?? Win Rate (Positive ROI): ")
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
	
	cyan.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("                    ?? ACHIEVEMENTS ??")
	cyan.Println(strings.Repeat("=", 60))
	
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
	fmt.Printf("\n?? Progress: %d/%d (%.1f%%)\n", unlockedCount, totalAchievements, progress)
	fmt.Printf("? Total Points: ")
	green.Printf("%d\n", totalPoints)
	fmt.Printf("???  Career Level: ")
	yellow.Printf("%d - %s\n", level, title)
	if level < 10 {
		fmt.Printf("?? Next Level: %d points needed\n", nextLevelPoints-totalPoints)
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

func displayHelpGuide() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    ?? HELP & INFORMATION ??")
	cyan.Println(strings.Repeat("=", 70))
	
	yellow.Println("\n?? GAME OVERVIEW")
	fmt.Println("You're a VC investor with limited capital. Invest in 20 startups")
	fmt.Println("and watch your portfolio grow (or shrink) over 10 years.")
	
	yellow.Println("\n?? HOW TO PLAY")
	fmt.Println("1. Select difficulty (Easy/Medium/Hard/Expert)")
	fmt.Println("2. Review 20 available companies with metrics")
	fmt.Println("3. Invest your capital across multiple startups")
	fmt.Println("4. Watch events unfold each turn (1 turn = 1 month)")
	fmt.Println("5. After 90-120 turns, see your final score")
	
	yellow.Println("\n?? COMPANY METRICS")
	fmt.Println("? Risk Score: Low/Medium/High - chance of failure")
	fmt.Println("? Growth Potential: Projected growth trajectory")
	fmt.Println("? Valuation: Current company worth (in millions)")
	fmt.Println("? Category: Industry sector (FinTech, BioTech, etc.)")
	
	yellow.Println("\n?? SCORING")
	fmt.Println("? Net Worth: Cash + Portfolio Value")
	fmt.Println("? ROI: Return on Investment percentage")
	fmt.Println("? Successful Exits: Companies that 5x or more")
	fmt.Println("? Rating: Based on ROI (Unicorn Hunter = 1000%+)")
	
	yellow.Println("\n?? DIFFICULTY LEVELS")
	fmt.Println("? Easy: $500k, 20% events, 3% volatility")
	fmt.Println("? Medium: $250k, 30% events, 5% volatility")
	fmt.Println("? Hard: $150k, 40% events, 7% volatility")
	fmt.Println("? Expert: $100k, 50% events, 10% volatility, 90 turns")
	
	yellow.Println("\n?? ANALYTICS")
	fmt.Println("After each game, view detailed portfolio analytics:")
	fmt.Println("? Best/Worst performers")
	fmt.Println("? Sector breakdown")
	fmt.Println("? Win/loss ratio")
	fmt.Println("? Investment distribution")
	
	yellow.Println("\n?? AVAILABLE COMPANIES")
	fmt.Println("20 diverse startups across 12+ sectors:")
	fmt.Println("FinTech ? BioTech ? CleanTech ? HealthTech ? EdTech")
	fmt.Println("Robotics ? Security ? Gaming ? LegalTech ? AgriTech")
	fmt.Println("Logistics ? IoT ? Creative ? CloudTech ? and more!")
	
	yellow.Println("\n?? RANDOM EVENTS")
	fmt.Println("60+ possible events can affect your companies:")
	fmt.Println("? Funding rounds (Series A/B, IPO)")
	fmt.Println("? Product launches (success/failure)")
	fmt.Println("? Partnerships & acquisitions")
	fmt.Println("? Scandals & regulatory issues")
	fmt.Println("? Market conditions & competition")
	
	yellow.Println("\n?? STRATEGY TIPS")
	fmt.Println("? Diversify: Don't put everything in one company")
	fmt.Println("? Balance: Mix high-risk and low-risk investments")
	fmt.Println("? Sectors: Different industries perform differently")
	fmt.Println("? Research: Read company metrics carefully")
	fmt.Println("? Patience: Some companies take time to grow")
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	fmt.Print("\nPress 'Enter' to return to menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
