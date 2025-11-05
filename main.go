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
	animations "github.com/jamesacampbell/unicorn/animations"
	ascii "github.com/jamesacampbell/unicorn/ascii"

	// analytics "github.com/jamesacampbell/unicorn/analytics"
	clear "github.com/jamesacampbell/unicorn/clear"
	db "github.com/jamesacampbell/unicorn/database"
	game "github.com/jamesacampbell/unicorn/game"
	leaderboard "github.com/jamesacampbell/unicorn/leaderboard"
	logo "github.com/jamesacampbell/unicorn/logo"
	upgrades "github.com/jamesacampbell/unicorn/upgrades"
	yaml "gopkg.in/yaml.v2"
)

type gameData struct {
	Pot        int64  `yaml:"starting-cash"`
	BadThings  int64  `yaml:"number-of-bad-things-per-year"`
	Foreground string `yaml:"foreground-color"`
}

func formatCurrency(amount int64) string {
	if amount < 0 {
		return fmt.Sprintf("-$%s", formatCurrency(-amount))
	}

	str := fmt.Sprintf("%d", amount)
	var result string
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}

func initMenu() (username string) {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Name: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	// Check if this player has played before
	stats, err := db.GetPlayerStats(text)
	if err == nil && stats.TotalGames > 0 {
		// Returning player - show welcome back message with stats
		fmt.Println()
		green.Printf("🎉 Welcome back, %s!\n\n", text)

		cyan.Println(strings.Repeat("=", 60))
		cyan.Println("                  YOUR PLAYER STATS")
		cyan.Println(strings.Repeat("=", 60))

		yellow.Printf("\n📊 Games Played: %d\n", stats.TotalGames)
		yellow.Printf("💰 Best Net Worth: $%s\n", formatCurrency(stats.BestNetWorth))
		yellow.Printf("📈 Best ROI: %.1f%%\n", stats.BestROI*100)
		yellow.Printf("🚀 Total Exits: %d\n", stats.TotalExits)
		yellow.Printf("📊 Average Net Worth: $%s\n", formatCurrency(int64(stats.AverageNetWorth)))
		yellow.Printf("🎯 Win Rate: %.1f%%\n", stats.WinRate)

		// Get achievement count
		achievementCount, _ := db.GetPlayerAchievementCount(text)
		if achievementCount > 0 {
			yellow.Printf("🏆 Achievements Unlocked: %d\n", achievementCount)
		}

		// Get and display active upgrades
		playerUpgrades, err := db.GetPlayerUpgrades(text)
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

func displayWelcome(username string, difficulty game.Difficulty, playerUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta, color.Bold)
	green := color.New(color.FgGreen)

	cyan.Printf("\n%s, welcome to your investment journey!\n", username)
	fmt.Printf("\nDifficulty: ")
	yellow.Printf("%s\n", difficulty.Name)
	fmt.Printf("Fund Size: $%s\n", formatMoney(difficulty.StartingCash))
	fmt.Printf("Follow-on Reserve: $%s ($100k base + $50k per round)\n", formatMoney(int64(1000000)))
	fmt.Printf("Management Fee: 2%% annually ($%s/year)\n", formatMoney(int64(float64(difficulty.StartingCash)*0.02)))
	fmt.Printf("Game Duration: %d turns (%d years)\n", difficulty.MaxTurns, difficulty.MaxTurns/12)

	// Display active upgrades
	if len(playerUpgrades) > 0 {
		fmt.Println()
		green.Println("✨ ACTIVE UPGRADES:")
		for _, upgradeID := range playerUpgrades {
			if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
				fmt.Printf("   %s %s - %s\n", upgrade.Icon, upgrade.Name, upgrade.Description)
			}
		}
	}

	fmt.Println("\nEach turn = 1 month. Choose your investments wisely!")
	fmt.Println("Random events and funding rounds will affect valuations.")
	fmt.Println("Watch out for dilution when companies raise new rounds!")
	fmt.Println("Note: Uninvested cash from your initial fund is also available for follow-on investments!")
	fmt.Println("You can invest more from your available cash + follow-on reserve when companies raise Series rounds!")

	magenta.Println("\n?? COMPETING AGAINST:")
	fmt.Println("   ? CARL (Sterling & Cooper) - Conservative")
	fmt.Println("   ? Sarah Chen (Accel Partners) - Aggressive")
	fmt.Println("   ? Marcus Williams (Sequoia Capital) - Balanced")

	fmt.Print("\nPress 'Enter' to see available startups...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func displayStartup(s game.Startup, index int, availableCash int64, playerUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)

	cyan.Printf("\n[%d] %s\n", index+1, s.Name)
	fmt.Printf("    %s\n", s.Description)
	yellow.Printf("    Category: %s\n", s.Category)
	fmt.Printf("    Valuation: $%s\n", formatMoney(s.Valuation))
	// Check for Super Pro-Rata upgrade
	maxInvestmentPercent := 0.20
	hasSuperProRata := false
	for _, upgradeID := range playerUpgrades {
		if upgradeID == "super_pro_rata" {
			maxInvestmentPercent = 0.50
			hasSuperProRata = true
			break
		}
	}
	maxInvestment := int64(float64(s.Valuation) * maxInvestmentPercent)
	maxAvailable := maxInvestment
	if maxAvailable > availableCash {
		maxAvailable = availableCash
	}
	green.Printf("    Max Investment: $%s", formatMoney(maxAvailable))
	if maxAvailable < maxInvestment {
		fmt.Printf(" (limited by available cash, max would be $%s)", formatMoney(maxInvestment))
	} else {
		if hasSuperProRata {
			fmt.Printf(" (50%% of valuation)")
		} else {
			fmt.Printf(" (20%% of valuation)")
		}
	}
	fmt.Println()
	fmt.Printf("    Monthly Sales: %d units\n", s.MonthlySales)
	fmt.Printf("    Margin: %d%%\n", s.PercentMargin)
	fmt.Printf("    Website Visitors: %s/month\n", formatNumber(s.MonthlyWebsiteVisitors))

	// Risk indicator - show numbers if Due Diligence upgrade is owned
	hasDueDiligence := false
	for _, upgradeID := range playerUpgrades {
		if upgradeID == "due_diligence" {
			hasDueDiligence = true
			break
		}
	}

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
	if hasDueDiligence {
		// Show exact risk score number
		magenta.Printf(" (%.2f)", s.RiskScore)
	}
	fmt.Println()

	// Revenue Tracker - show revenue trends if upgrade is owned
	hasRevenueTracker := false
	for _, upgradeID := range playerUpgrades {
		if upgradeID == "revenue_tracker" {
			hasRevenueTracker = true
			break
		}
	}

	if hasRevenueTracker && len(s.RevenueHistory) > 1 {
		// Show revenue trend
		fmt.Printf("    Revenue Trend: ")
		trend := s.RevenueHistory[len(s.RevenueHistory)-1] - s.RevenueHistory[0]
		if trend > 0 {
			green.Printf("↑ +%.1f%%", float64(trend)/float64(s.RevenueHistory[0])*100)
		} else if trend < 0 {
			red := color.New(color.FgRed)
			red.Printf("↓ %.1f%%", float64(trend)/float64(s.RevenueHistory[0])*100)
		} else {
			fmt.Printf("→ Stable")
		}
		fmt.Printf(" over %d months\n", len(s.RevenueHistory))
	}

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
	growthColor.Printf("    Growth Potential: %s\n", growthLabel)
}

func handleFollowOnOpportunities(gs *game.GameState, opportunities []game.FollowOnOpportunity) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)

	// Clear screen to make it obvious
	fmt.Println("\n\n")
	fmt.Println(strings.Repeat("=", 70))
	magenta.Println("            ???? FOLLOW-ON INVESTMENT OPPORTUNITY!")
	fmt.Println(strings.Repeat("=", 70))
	cyan.Println("\nOne of your portfolio companies is raising a new funding round!")
	cyan.Println("You can invest MORE money to avoid dilution and increase ownership.")

	for _, opp := range opportunities {
		fmt.Println("\n" + strings.Repeat("-", 70))
		magenta.Printf("\n?? COMPANY: %s\n", opp.CompanyName)
		fmt.Printf("   Raising: %s round\n", opp.RoundName)
		fmt.Printf("   Pre-money Valuation: $%s\n", formatMoney(opp.PreMoneyVal))
		fmt.Printf("   Post-money Valuation: $%s\n", formatMoney(opp.PostMoneyVal))
		yellow.Printf("   Your Current Equity: %.2f%%\n", opp.CurrentEquity)
		availableFunds := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve
		green.Printf("   Available Funds: $%s (Cash: $%s + Reserve: $%s)\n",
			formatMoney(availableFunds),
			formatMoney(gs.Portfolio.Cash),
			formatMoney(gs.Portfolio.FollowOnReserve))

		fmt.Println("\n" + strings.Repeat("-", 70))
		cyan.Println("\n?? INVEST MORE TO AVOID DILUTION!")
		fmt.Println("   If you don't invest, your ownership % will decrease.")
		fmt.Println("   If you DO invest, you'll maintain or increase your stake.")
		fmt.Printf("\n   Investment Range: $%s to $%s\n",
			formatMoney(opp.MinInvestment), formatMoney(opp.MaxInvestment))

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n💰 Enter amount to invest (0 or Enter to skip): $")
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)

		if amountStr == "" || amountStr == "0" {
			color.Yellow("Skipping follow-on investment.")
			continue
		}

		amount, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil || amount < 0 {
			color.Yellow("Invalid amount, skipping.")
			continue
		}

		if amount == 0 {
			color.Yellow("Skipping follow-on investment.")
			continue
		}

		if amount < opp.MinInvestment {
			color.Red("Amount below minimum investment of $%s", formatMoney(opp.MinInvestment))
			continue
		}

		if amount > opp.MaxInvestment {
			color.Red("Amount exceeds maximum investment of $%s", formatMoney(opp.MaxInvestment))
			continue
		}

		err = gs.MakeFollowOnInvestment(opp.CompanyName, amount)
		if err != nil {
			color.Red("Error: %v", err)
		} else {
			green.Printf("\n%s Follow-on investment successful! Invested $%s in %s\n",
				ascii.Check, formatMoney(amount), opp.CompanyName)
			fmt.Printf("Follow-on Reserve Remaining: $%s\n", formatMoney(gs.Portfolio.FollowOnReserve))
			fmt.Printf("Cash Remaining: $%s\n", formatMoney(gs.Portfolio.Cash))
		}
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleBoardVotes(gs *game.GameState, votes []game.BoardVote) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	red := color.New(color.FgRed, color.Bold)

	for _, vote := range votes {
		// Find the vote index in the full pending votes list
		voteIndex := -1
		for i, v := range gs.GetPendingBoardVotes() {
			if v.CompanyName == vote.CompanyName && v.VoteType == vote.VoteType && v.Turn == vote.Turn {
				voteIndex = i
				break
			}
		}

		if voteIndex == -1 {
			continue // Vote already processed
		}

		fmt.Println("\n\n")
		fmt.Println(strings.Repeat("=", 70))
		magenta.Println("            🏛️  BOARD VOTE REQUIRED!")
		fmt.Println(strings.Repeat("=", 70))

		cyan.Printf("\nCompany: %s\n", vote.CompanyName)
		yellow.Printf("\n%s\n", vote.Title)
		fmt.Println("\n" + strings.Repeat("-", 70))
		fmt.Printf("\n%s\n", vote.Description)
		fmt.Println("\n" + strings.Repeat("-", 70))

		green.Printf("\nOption A: %s\n", vote.OptionA)
		fmt.Printf("   → %s\n", vote.ConsequenceA)

		red.Printf("\nOption B: %s\n", vote.OptionB)
		fmt.Printf("   → %s\n", vote.ConsequenceB)

		fmt.Println("\n" + strings.Repeat("-", 70))
		cyan.Println("\nYour vote as a board member:")

		// Check for double board seat upgrade
		voteWeight := 1
		for _, inv := range gs.Portfolio.Investments {
			if inv.CompanyName == vote.CompanyName && inv.Terms.HasBoardSeat {
				voteWeight = inv.Terms.BoardSeatMultiplier
				if voteWeight == 0 {
					voteWeight = 1
				}
				break
			}
		}

		if voteWeight > 1 {
			magenta.Printf("Voting Power: %d votes (Double Board Seat upgrade active!)\n", voteWeight)
		} else {
			fmt.Println("Voting Power: 1 vote")
		}

		fmt.Print("Vote (A/1 for Accept/Approve, B/2 for Reject/Disapprove): ")

		reader := bufio.NewReader(os.Stdin)
		voteChoice, _ := reader.ReadString('\n')
		voteChoice = strings.TrimSpace(voteChoice)

		// Re-find vote index since list may have changed
		voteIndex = -1
		for i, v := range gs.GetPendingBoardVotes() {
			if v.CompanyName == vote.CompanyName && v.VoteType == vote.VoteType && v.Turn == vote.Turn {
				voteIndex = i
				break
			}
		}

		if voteIndex == -1 {
			color.Yellow("Vote already processed.")
			continue
		}

		result, passed, err := gs.ProcessBoardVote(voteIndex, voteChoice)
		if err != nil {
			color.Red("Error: %v", err)
			fmt.Print("\nPress 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		fmt.Println("\n" + strings.Repeat("-", 70))
		if passed {
			green.Printf("\n✅ %s\n", result)
		} else {
			yellow.Printf("\n❌ %s\n", result)
		}

		// Execute vote outcome
		outcomeMessages := gs.ExecuteBoardVoteOutcome(vote, passed)
		for _, msg := range outcomeMessages {
			fmt.Println(msg)
		}

		fmt.Print("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func selectInvestmentTerms(gs *game.GameState, startup *game.Startup, amount int64) game.InvestmentTerms {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	fmt.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                      INVESTMENT TERMS")
	fmt.Println(strings.Repeat("=", 70))

	options := gs.GenerateTermOptions(startup, amount)

	fmt.Println("\nSelect your investment structure:")
	for i, opt := range options {
		fmt.Printf("\n%d. %s\n", i+1, opt.Type)
		if opt.HasProRataRights {
			green.Println("   ✓ Pro-Rata Rights (participate in future rounds)")
		}
		if opt.HasInfoRights {
			green.Println("   ✓ Information Rights (quarterly financials)")
		}
		if opt.HasBoardSeat {
			green.Println("   ✓ Board Observer Seat")
		}
		if opt.LiquidationPref > 0 {
			green.Printf("   ✓ %dx Liquidation Preference (get paid first)\n", int(opt.LiquidationPref))
		}
		if opt.HasAntiDilution {
			green.Println("   ✓ Anti-Dilution Protection (protect from down rounds)")
		}
		if opt.ConversionDiscount > 0 {
			green.Printf("   ✓ %.0f%% Conversion Discount (bonus equity)\n", opt.ConversionDiscount*100)
		}
	}

	maxOption := len(options)
	fmt.Printf("\nSelect terms (1-%d, or Enter for Preferred): ", maxOption)
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "" {
		choice = "1"
	}

	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(options) {
		yellow.Println("Invalid choice, using Preferred Stock")
		return options[0]
	}

	return options[choiceNum-1]
}

func investmentPhase(gs *game.GameState) {
	clear.ClearIt()
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Print(ascii.InvestmentHeader)
	green.Printf("\nTurn %d/%d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	fmt.Printf("Fund Size: $%s\n", formatMoney(gs.Portfolio.InitialFundSize))
	fmt.Printf("Cash Available: $%s\n", formatMoney(gs.Portfolio.Cash))
	fmt.Printf("Follow-on Reserve: $%s\n", formatMoney(gs.Portfolio.FollowOnReserve))
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

	// Get player upgrades for display
	playerUpgrades, err := db.GetPlayerUpgrades(gs.PlayerName)
	if err != nil {
		playerUpgrades = []string{}
	}

	// Market Intelligence: Show sector trends
	hasMarketIntelligence := false
	for _, upgradeID := range playerUpgrades {
		if upgradeID == "market_intelligence" {
			hasMarketIntelligence = true
			break
		}
	}
	if hasMarketIntelligence {
		sectorTrends := gs.GetSectorTrends()
		if len(sectorTrends) > 0 {
			fmt.Println()
			magenta := color.New(color.FgMagenta, color.Bold)
			magenta.Println("📈 MARKET INTELLIGENCE - Category Trends:")
			fmt.Println("   (Trends based on average valuations across available startups)")
			for sector, trend := range sectorTrends {
				fmt.Printf("   %s: %s\n", sector, trend)
			}
			fmt.Println()
		}
	}

	for i, startup := range gs.AvailableStartups {
		displayStartup(startup, i, gs.Portfolio.Cash, playerUpgrades)
	}
	fmt.Println(strings.Repeat("=", 50))

	reader := bufio.NewReader(os.Stdin)

	for {
		// Auto-start if out of money
		if gs.Portfolio.Cash <= 0 {
			color.Yellow("\n⚠️  Out of investment capital! Starting game...")
			gs.AIPlayerMakeInvestments()
			time.Sleep(2 * time.Second)
			break
		}

		fmt.Printf("\nEnter company number (1-%d) to invest, or press Enter to start: ", len(gs.AvailableStartups))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "done" || input == "" {
			// Have AI players make their investments too
			gs.AIPlayerMakeInvestments()
			color.Green("\n✓ AI players have made their investments!")
			break
		}

		companyNum, err := strconv.Atoi(input)
		if err != nil || companyNum < 1 || companyNum > len(gs.AvailableStartups) {
			color.Red("Invalid company number!")
			continue
		}

		startup := gs.AvailableStartups[companyNum-1]
		// Check for Super Pro-Rata upgrade
		maxInvestmentPercent := 0.20
		hasSuperProRata := false
		for _, upgradeID := range playerUpgrades {
			if upgradeID == "super_pro_rata" {
				maxInvestmentPercent = 0.50
				hasSuperProRata = true
				break
			}
		}
		maxInvestment := int64(float64(startup.Valuation) * maxInvestmentPercent)
		maxInvestmentDisplay := maxInvestment
		if maxInvestmentDisplay > gs.Portfolio.Cash {
			maxInvestmentDisplay = gs.Portfolio.Cash
		}

		cyan := color.New(color.FgCyan, color.Bold)
		yellow := color.New(color.FgYellow)
		cyan.Printf("\n💵 INVESTING IN: %s\n", startup.Name)
		fmt.Printf("   Valuation: $%s\n", formatMoney(startup.Valuation))
		yellow.Printf("   Max Investment Available: $%s", formatMoney(maxInvestmentDisplay))
		if maxInvestmentDisplay < maxInvestment {
			fmt.Printf(" (limited by cash, max would be $%s)", formatMoney(maxInvestment))
		} else {
			if hasSuperProRata {
				fmt.Printf(" (50%% of valuation)")
			} else {
				fmt.Printf(" (20%% of valuation)")
			}
		}
		fmt.Println()
		fmt.Printf("\nEnter investment amount ($10,000 - $%s, or 0 to skip): $", formatMoney(maxInvestmentDisplay))
		amountStr, _ := reader.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)

		if amountStr == "" || amountStr == "0" {
			color.Yellow("Skipped investment in this company.")
			continue
		}

		amount, err := strconv.ParseInt(amountStr, 10, 64)

		if err != nil {
			color.Red("Invalid amount!")
			continue
		}

		if amount == 0 {
			color.Yellow("Skipped investment in this company.")
			continue
		}

		// Validate amount before proceeding
		if amount < 10000 {
			color.Red("Minimum investment is $10,000")
			continue
		}
		if amount > maxInvestmentDisplay {
			maxPercentText := "20%"
			if hasSuperProRata {
				maxPercentText = "50%"
			}
			color.Red("Maximum investment is $%s (%s of company valuation: $%s)", formatMoney(maxInvestmentDisplay), maxPercentText, formatMoney(startup.Valuation))
			continue
		}

		// Show term options for investments $50k+
		var selectedTerms game.InvestmentTerms
		if amount >= 50000 {
			selectedTerms = selectInvestmentTerms(gs, &gs.AvailableStartups[companyNum-1], amount)
		} else {
			// Default terms for smaller investments
			selectedTerms = game.InvestmentTerms{
				Type:                "Common Stock",
				HasProRataRights:    false,
				HasInfoRights:       false,
				HasBoardSeat:        false,
				BoardSeatMultiplier: 1,
				LiquidationPref:     0.0,
				HasAntiDilution:     false,
			}
		}

		err = gs.MakeInvestmentWithTerms(companyNum-1, amount, selectedTerms)
		if err != nil {
			color.Red("Error: %v", err)
		} else {
			color.Green("%s Investment successful!", ascii.Check)
			fmt.Printf("Cash remaining: $%s\n", formatMoney(gs.Portfolio.Cash))
			fmt.Printf("Terms: %s\n", selectedTerms.Type)
		}
	}
}

func playTurn(gs *game.GameState, autoMode bool) {
	yellow := color.New(color.FgYellow, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)

	// Check for follow-on investment opportunities BEFORE processing turn
	// This way the player can invest before dilution happens
	// Always pause for follow-on investments, even in automated mode
	opportunities := gs.GetFollowOnOpportunities()
	if len(opportunities) > 0 {
		handleFollowOnOpportunities(gs, opportunities)
	}

	// Print separator line instead of clearing screen
	fmt.Println(strings.Repeat("=", 70))

	// Show round transition animation for milestones (every 5 turns)
	if gs.Portfolio.Turn%5 == 0 {
		animations.ShowRoundTransition(gs.Portfolio.Turn)
	}

	yellow.Printf("\n%s MONTH %d of %d\n", ascii.Calendar, gs.Portfolio.Turn, gs.Portfolio.MaxTurns)

	// Strategic Advisor: Show preview of next board vote
	playerUpgrades, _ := db.GetPlayerUpgrades(gs.PlayerName)
	hasStrategicAdvisor := false
	for _, upgradeID := range playerUpgrades {
		if upgradeID == "strategic_advisor" {
			hasStrategicAdvisor = true
			break
		}
	}
	if hasStrategicAdvisor {
		nextVote := gs.GetNextBoardVotePreview()
		if nextVote != "" {
			magenta := color.New(color.FgMagenta, color.Bold)
			fmt.Println()
			magenta.Println("🔮 STRATEGIC ADVISOR PREVIEW:")
			fmt.Println(nextVote)
			fmt.Println()
		}
	}

	messages := gs.ProcessTurn()

	// Check for pending board votes AFTER processing turn
	// Board votes are created during ProcessTurn for acquisitions/down rounds
	pendingVotes := gs.GetPendingBoardVotes()
	if len(pendingVotes) > 0 {
		handleBoardVotes(gs, pendingVotes)
		// Re-execute vote outcomes (votes are already processed, just need to show outcomes)
		// Note: Vote outcomes are executed in handleBoardVotes, so we don't need to ProcessTurn again
	}

	// Separate critical messages (that need pause) from informational messages
	criticalMessages := []string{}
	infoMessages := []string{}
	hasExitEvent := false

	for _, msg := range messages {
		// Check for exit events (ACQUIRED, IPO, etc.)
		if strings.Contains(msg, "ACQUIRED") || strings.Contains(msg, "acquisition") {
			hasExitEvent = true
			criticalMessages = append(criticalMessages, msg)
			// Check for dramatic events (scandals, co-founder splits, fraud, etc.)
		} else if strings.Contains(msg, "💔") || strings.Contains(msg, "🔥") ||
			strings.Contains(msg, "⚖️") || strings.Contains(msg, "🚨") ||
			strings.Contains(msg, "🔓") || strings.Contains(msg, "👋") ||
			strings.Contains(msg, "📋") || strings.Contains(msg, "🔄") ||
			strings.Contains(msg, "⚔️") || strings.Contains(msg, "💥") {
			criticalMessages = append(criticalMessages, msg)
			// Critical messages: funding rounds, dilution, negative news, down rounds
		} else if strings.Contains(msg, "raised") ||
			strings.Contains(msg, "diluted") ||
			strings.Contains(msg, "DOWN ROUND") {
			criticalMessages = append(criticalMessages, msg)
		} else {
			infoMessages = append(infoMessages, msg)
		}
	}

	// Always show all messages
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
	fmt.Printf(" | Cash: $%s | Follow-on Reserve: $%s\n",
		formatMoney(gs.Portfolio.Cash), formatMoney(gs.Portfolio.FollowOnReserve))
	fmt.Printf("   Management Fees Paid: $%s\n", formatMoney(gs.Portfolio.ManagementFeesCharged))

	// Show competitive leaderboard every quarter
	if gs.Portfolio.Turn%3 == 0 {
		displayMiniLeaderboard(gs)
	}

	// Always pause for exit events, even in auto mode
	if hasExitEvent {
		magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Println("\n🎉 COMPANY EXIT EVENT! 🎉")
		fmt.Print("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		// Only pause for other critical messages in auto mode; informational messages don't pause
	} else if len(criticalMessages) > 0 {
		// Critical message - always pause
		fmt.Print("\nPress 'Enter' to continue to next month...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	} else if autoMode {
		// No critical messages and auto mode - don't pause
		time.Sleep(1 * time.Second)
	} else {
		// Manual mode - always pause
		fmt.Print("\nPress 'Enter' to continue to next month...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func displayFinalScore(gs *game.GameState) {
	// Show animated game over screen FIRST (before clearing)
	netWorth, roi, successfulExits := gs.GetFinalScore()
	won := netWorth >= gs.Difficulty.StartingCash*2 // Won if doubled starting cash
	animations.ShowGameOverAnimation(won, netWorth)

	// Pause to let user see the animation
	fmt.Print("\nPress 'Enter' to see detailed results...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	clear.ClearIt()

	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)

	cyan.Print(ascii.GameOverHeader)

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
	} else if roi >= 10 {
		rating = "Profitable - Not Bad"
		icon = ascii.Check
	} else if roi >= -10 {
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
	// Show animated splash screen on first launch
	animations.ShowGameStartAnimation()

	// Pause to let user enjoy the splash screen
	fmt.Print("\nPress 'Enter' to continue to main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

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
			displayUpgradeMenu()
		case "6":
			displayHelpGuide()
		case "7":
			animations.ShowInfoMessage("Thanks for playing! " + ascii.Star2)
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
	yellow.Println("5. Upgrades")
	yellow.Println("6. Help & Info")
	yellow.Println("7. Quit")

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

	// Select game mode
	gameMode := askForGameMode()
	clear.ClearIt()

	if gameMode == "founder" {
		playFounderMode(username)
	} else {
		playVCMode(username)
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

	if choice == "2" {
		return "founder"
	}
	return "vc"
}

func playVCMode(username string) {
	// Select difficulty
	difficulty := selectDifficulty()
	clear.ClearIt()

	// Ask for automated mode
	autoMode := askForAutomatedMode()
	clear.ClearIt()

	// Get player upgrades
	playerUpgrades, err := db.GetPlayerUpgrades(username)
	if err != nil {
		playerUpgrades = []string{}
	}

	// Display welcome and rules (with upgrades)
	displayWelcome(username, difficulty, playerUpgrades)

	// Initialize game
	gs := game.NewGame(username, difficulty, playerUpgrades)

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

	err = db.SaveGameScore(score)
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
			db.UnlockAchievement(gs.PlayerName, ach.ID)

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
	allUnlocked, _ := db.GetPlayerAchievements(gs.PlayerName)
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalLifetimePoints += ach.Points
		}
	}

	// Get owned upgrades to calculate available balance
	ownedUpgrades, _ := db.GetPlayerUpgrades(gs.PlayerName)
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

	// Show spinner while loading
	spinner, _ := animations.StartSpinner("Loading leaderboard...")
	defer spinner.Stop()

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

	// Show spinner while loading
	spinner, _ := animations.StartSpinner("Loading recent games...")
	defer spinner.Stop()

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

	// Show spinner while loading
	spinner, _ := animations.StartSpinner("Loading player stats...")
	defer spinner.Stop()

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
	yellow := color.New(color.FgYellow)

	// Show spinner while loading detailed stats
	spinner2, _ := animations.StartSpinner("Loading detailed stats...")
	defer spinner2.Stop()

	// Get stats for VC mode
	vcStats, err := db.GetPlayerStatsByMode(playerName, "vc")
	if err != nil {
		color.Red("Error loading VC stats: %v", err)
		time.Sleep(2 * time.Second)
		return
	}

	// Get stats for Founder mode
	founderStats, err := db.GetPlayerStatsByMode(playerName, "founder")
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
		green.Printf("$%s\n", formatMoney(vcStats.BestNetWorth))

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
		green.Printf("$%s\n", formatMoney(founderStats.BestNetWorth))

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
	playerUpgrades, err := db.GetPlayerUpgrades(playerName)
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
	yellow.Println("https://james-see.github.io/unicorn")
}

func displayInvestingFAQ() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("              STARTUP INVESTING FAQ")
	cyan.Println(strings.Repeat("=", 70))

	yellow.Println("\n💰 INVESTMENT TERMS")
	fmt.Println()
	fmt.Println("Q: What's the difference between Preferred and Common Stock?")
	green.Println("A: Preferred Stock gives you:")
	fmt.Println("   • Liquidation preference (get paid back first in exit)")
	fmt.Println("   • Anti-dilution protection (protects value in down rounds)")
	fmt.Println("   • Information rights (quarterly financial updates)")
	fmt.Println("   • Pro-rata rights (right to invest in future rounds)")
	fmt.Println("   Common Stock has none of these protections.")

	fmt.Println("\nQ: What is a SAFE?")
	green.Println("A: Simple Agreement for Future Equity:")
	fmt.Println("   • Converts to equity in the next priced round")
	fmt.Println("   • Usually includes a 15-20% discount")
	fmt.Println("   • No liquidation preference")
	fmt.Println("   • Simpler than convertible notes")
	fmt.Println("   • Popular for early-stage investing")

	fmt.Println("\nQ: What are Pro-Rata Rights?")
	green.Println("A: The right to maintain your ownership % in future rounds:")
	fmt.Println("   • When a company raises Series A, you can invest more")
	fmt.Println("   • Prevents dilution of your stake")
	fmt.Println("   • Requires additional capital from available cash or follow-on reserve")
	fmt.Println("   • Essential for successful investments")

	yellow.Println("\n📊 VALUATION & EQUITY")
	fmt.Println()
	fmt.Println("Q: How is equity calculated?")
	green.Println("A: Your ownership % = (Your Investment / Post-Money Valuation) × 100")
	fmt.Println("   Example: $100k into $1M valuation = 10% ownership")
	fmt.Println("   Post-Money = Pre-Money + Total Round Size")

	fmt.Println("\nQ: What is dilution?")
	green.Println("A: When a company raises new funding, all existing shareholders")
	fmt.Println("   get diluted unless they invest more (pro-rata rights):")
	fmt.Println("   • You own 10% after your investment")
	fmt.Println("   • Company raises Series A ($10M)")
	fmt.Println("   • Your 10% becomes ~7% (30% dilution)")
	fmt.Println("   • Your $ value may still increase if valuation grows")

	yellow.Println("\n🚀 FUNDING ROUNDS")
	fmt.Println()
	fmt.Println("Q: What are the typical funding stages?")
	green.Println("A: Pre-Seed → Seed → Series A → Series B → Series C → IPO")
	fmt.Println("   • Pre-Seed: $250k-$1M (you invest here)")
	fmt.Println("   • Seed: $2M-$5M (3-9 months)")
	fmt.Println("   • Series A: $10M-$20M (12-24 months)")
	fmt.Println("   • Series B: $30M-$50M (30-48 months)")
	fmt.Println("   • Series C+: $50M-$100M+ (48+ months)")

	fmt.Println("\nQ: What is a down round?")
	green.Println("A: When a company raises at a LOWER valuation:")
	fmt.Println("   • Bad signal to market")
	fmt.Println("   • Heavy dilution for existing investors")
	fmt.Println("   • Anti-dilution protection helps here")
	fmt.Println("   • Happens ~20% of the time")

	yellow.Println("\n💼 EXIT STRATEGIES")
	fmt.Println()
	fmt.Println("Q: How do I make money?")
	green.Println("A: Three main exits:")
	fmt.Println("   • Acquisition: Company gets bought (most common)")
	fmt.Println("   • IPO: Company goes public (rare but huge)")
	fmt.Println("   • Secondary Sale: Sell shares to another investor")

	fmt.Println("\nQ: What's a good return?")
	green.Println("A: VC benchmarks:")
	fmt.Println("   • 3x: Good return")
	fmt.Println("   • 5x: Great return")
	fmt.Println("   • 10x: Excellent return")
	fmt.Println("   • 100x: Unicorn! (1 in 1000 startups)")

	yellow.Println("\n⚠️  RISK MANAGEMENT")
	fmt.Println()
	fmt.Println("Q: How should I diversify?")
	green.Println("A: Rule of thumb:")
	fmt.Println("   • Invest in 8-12 companies minimum")
	fmt.Println("   • Mix of high-risk/high-reward and safer bets")
	fmt.Println("   • Different sectors (FinTech, BioTech, etc.)")
	fmt.Println("   • ~70% of startups will fail or break even")
	fmt.Println("   • You need 1-2 big winners to make up for losses")

	fmt.Println("\nQ: What kills startups?")
	green.Println("A: Top reasons:")
	fmt.Println("   • Running out of cash (38%)")
	fmt.Println("   • No market need (35%)")
	fmt.Println("   • Competition (20%)")
	fmt.Println("   • Bad timing (17%)")
	fmt.Println("   • Co-founder conflicts (13%)")

	yellow.Println("\n📈 KEY METRICS")
	fmt.Println()
	fmt.Println("Q: What metrics matter?")
	green.Println("A: Watch these:")
	fmt.Println("   • Monthly Recurring Revenue (MRR) - predictable income")
	fmt.Println("   • Burn Rate - how fast they spend cash")
	fmt.Println("   • Customer Acquisition Cost (CAC) - cost per customer")
	fmt.Println("   • Lifetime Value (LTV) - revenue per customer")
	fmt.Println("   • LTV:CAC Ratio - should be 3:1 or better")

	cyan.Println("\n" + strings.Repeat("=", 70))
	fmt.Print("\nPress 'Enter' to return...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func displayHelpGuide() {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("              HELP & INFORMATION")
	cyan.Println(strings.Repeat("=", 70))

	fmt.Println("\n1. Game Overview & Rules")
	fmt.Println("2. Startup Investing FAQ")
	fmt.Println("3. Back to Main Menu")

	fmt.Print("\nEnter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "2" {
		displayInvestingFAQ()
		return
	} else if choice == "3" {
		return
	}

	clear.ClearIt()
	cyan.Println("\n" + strings.Repeat("?", 70))
	cyan.Println("              GAME OVERVIEW & RULES")
	cyan.Println(strings.Repeat("?", 70))

	yellow.Printf("\n%s GAME OVERVIEW\n", ascii.Lightbulb)
	fmt.Println("You're a VC investor with limited capital. Invest in 15 randomly")
	fmt.Println("selected startups from a pool of 30 and watch your portfolio")
	fmt.Println("grow (or shrink) over 5 years.")

	yellow.Printf("\n%s HOW TO PLAY\n", ascii.Target)
	fmt.Println("1. Select difficulty (Easy/Medium/Hard/Expert)")
	fmt.Println("2. Review 15 randomly selected companies (from 30 total)")
	fmt.Println("3. Invest your capital across multiple startups")
	fmt.Println("4. Watch events unfold each turn (1 turn = 1 month)")
	fmt.Println("5. Invest follow-on capital when companies raise new rounds")
	fmt.Println("6. After 60 turns (5 years), see your final score")

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
	fmt.Printf("%s Easy: $1M fund, 20%% events, 3%% volatility, 5 years\n", ascii.Check)
	fmt.Printf("%s Medium: $750k fund, 30%% events, 5%% volatility, 5 years\n", ascii.Star)
	fmt.Printf("%s Hard: $500k fund, 40%% events, 7%% volatility, 5 years\n", ascii.Warning)
	fmt.Printf("%s Expert: $500k fund, 50%% events, 10%% volatility, 5 years\n", ascii.Zap)

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

func displayUpgradeMenu() {
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
	allUnlocked, err := db.GetPlayerAchievements(playerName)
	if err != nil {
		allUnlocked = []string{}
	}

	totalLifetimePoints := 0
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalLifetimePoints += ach.Points
		}
	}

	// Get owned upgrades
	ownedUpgrades, err := db.GetPlayerUpgrades(playerName)
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
		browseAllUpgrades(playerName, availableBalance, ownedUpgrades)
	case "2":
		viewPlayerUpgrades(playerName, ownedUpgrades)
	case "3":
		purchaseUpgrades(playerName, availableBalance, ownedUpgrades)
	case "4":
		return
	default:
		color.Red("Invalid choice!")
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	displayUpgradeMenu()
}

func browseAllUpgrades(playerName string, availableBalance int, ownedUpgrades []string) {
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

func viewPlayerUpgrades(playerName string, ownedUpgrades []string) {
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

func purchaseUpgrades(playerName string, totalPoints int, ownedUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Refresh points and owned upgrades from database
	allUnlocked, err := db.GetPlayerAchievements(playerName)
	if err != nil {
		allUnlocked = []string{}
	}

	currentPoints := 0
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			currentPoints += ach.Points
		}
	}

	// Deduct cost of owned upgrades
	currentOwnedUpgrades, err := db.GetPlayerUpgrades(playerName)
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
	err = db.PurchaseUpgrade(playerName, upgrade.ID)
	if err != nil {
		color.Red("Error purchasing upgrade: %v", err)
		return
	}

	green.Printf("\n✓ Successfully purchased: %s %s!\n", upgrade.Icon, upgrade.Name)

	// Refresh points and owned upgrades from database
	allUnlockedRefresh, errRefresh := db.GetPlayerAchievements(playerName)
	if errRefresh != nil {
		allUnlockedRefresh = []string{}
	}

	newTotalPoints := 0
	for _, id := range allUnlockedRefresh {
		if ach, exists := achievements.AllAchievements[id]; exists {
			newTotalPoints += ach.Points
		}
	}

	newOwnedUpgrades, errUpgrades := db.GetPlayerUpgrades(playerName)
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
	purchaseUpgrades(playerName, newTotalPoints, newOwnedUpgrades)
}
