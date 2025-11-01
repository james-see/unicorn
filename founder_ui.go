package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	ascii "github.com/jamesacampbell/unicorn/ascii"
	clear "github.com/jamesacampbell/unicorn/clear"
	founder "github.com/jamesacampbell/unicorn/founder"
)

func playFounderMode(username string) {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	// Load available startups
	startups, err := founder.LoadFounderStartups("founder/startups.json")
	if err != nil {
		color.Red("Error loading startups: %v", err)
		fmt.Print("\nPress 'Enter' to return to main menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// Display startup selection
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("              CHOOSE YOUR STARTUP")
	cyan.Println(strings.Repeat("=", 70))

	for i, s := range startups {
		yellow.Printf("\n%d. %s", i+1, s.Name)
		fmt.Printf(" [%s]\n", s.Type)
		fmt.Printf("   %s\n", s.Tagline)
		fmt.Printf("   Cash: $%s | MRR: $%s | Customers: %d\n",
			formatFounderCurrency(s.InitialCash),
			formatFounderCurrency(s.InitialMRR),
			s.InitialCustomers)
		fmt.Printf("   Market Size: %d | Competition: %s\n",
			s.TargetMarketSize, s.CompetitionLevel)
	}

	fmt.Print("\nSelect your startup (1-" + fmt.Sprintf("%d", len(startups)) + "): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(startups) {
		color.Red("Invalid choice!")
		fmt.Print("\nPress 'Enter' to return to main menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	selectedStartup := startups[choiceNum-1]
	clear.ClearIt()

	// Initialize founder game
	fs := founder.NewFounderGame(username, selectedStartup)

	// Welcome message
	displayFounderWelcome(fs)

	// Main game loop
	for !fs.IsGameOver() {
		playFounderTurn(fs)
	}

	// Show final score
	displayFounderFinalScore(fs)

	fmt.Print("\nPress 'Enter' to return to main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func displayFounderWelcome(fs *founder.FounderState) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Printf("\n🚀 Welcome to %s!\n", fs.CompanyName)
	fmt.Printf("\nFounder: %s\n", fs.FounderName)
	fmt.Printf("Category: %s\n", fs.Category)
	fmt.Printf("\n💰 Starting Cash: $%s\n", formatFounderCurrency(fs.Cash))
	fmt.Printf("📊 MRR: $%s\n", formatFounderCurrency(fs.MRR))
	fmt.Printf("👥 Team Size: %d\n", fs.Team.TotalEmployees)
	fmt.Printf("⏱️  Runway: %d months\n", fs.CashRunwayMonths)
	fmt.Printf("📈 Product Maturity: %.0f%%\n", fs.ProductMaturity*100)

	yellow.Println("\n\n🎯 YOUR GOAL:")
	fmt.Println("Build a successful company over 60 months (5 years)")
	fmt.Println("• Hire the right team to grow your product and sales")
	fmt.Println("• Raise funding when needed (watch your runway!)")
	fmt.Println("• Manage churn and keep customers happy")
	fmt.Println("• Avoid running out of cash!")

	fmt.Print("\nPress 'Enter' to start your journey...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	clear.ClearIt()
}

func playFounderTurn(fs *founder.FounderState) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	// Display current status
	clear.ClearIt()
	cyan.Printf("\n[📅] MONTH %d of %d\n", fs.Turn, fs.MaxTurns)
	fmt.Println(strings.Repeat("=", 70))

	// Check for low cash warning
	if fs.NeedsLowCashWarning() {
		red.Println("\n⚠️  WARNING: Cash is running low!")
		red.Printf("   Cash: $%s | Runway: %d months\n", formatFounderCurrency(fs.Cash), fs.CashRunwayMonths)
		yellow.Println("   Consider: Raise funding, cut costs, or speed up revenue growth!")
	}

	// Check for acquisition offer
	if offer := fs.CheckForAcquisition(); offer != nil {
		displayAcquisitionOffer(fs, offer)
	}

	// Show company metrics
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("📊 COMPANY METRICS")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Printf("💰 Cash: $%s", formatFounderCurrency(fs.Cash))
	if fs.CashRunwayMonths < 6 {
		red.Printf(" (⚠️  %d months runway)\n", fs.CashRunwayMonths)
	} else {
		fmt.Printf(" (%d months runway)\n", fs.CashRunwayMonths)
	}
	
	fmt.Printf("📈 MRR: $%s", formatFounderCurrency(fs.MRR))
	if fs.MonthlyGrowthRate > 0 {
		green.Printf(" (↑%.1f%%)\n", fs.MonthlyGrowthRate*100)
	} else if fs.MonthlyGrowthRate < 0 {
		red.Printf(" (↓%.1f%%)\n", -fs.MonthlyGrowthRate*100)
	} else {
		fmt.Println(" (flat)")
	}
	
	fmt.Printf("👥 Customers: %d | ", fs.Customers)
	fmt.Printf("💸 Avg Deal: $%s/mo\n", formatFounderCurrency(fs.AvgDealSize))
	fmt.Printf("🔄 Churn Rate: %.1f%% | ", fs.CustomerChurnRate*100)
	fmt.Printf("📦 Product Maturity: %.0f%%\n", fs.ProductMaturity*100)
	fmt.Printf("💼 Your Equity: %.1f%%\n", 100.0-fs.EquityGivenAway)

	// Show team
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("👥 TEAM")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("Engineers: %d | Sales: %d | CS: %d | Marketing: %d\n",
		len(fs.Team.Engineers), len(fs.Team.Sales),
		len(fs.Team.CustomerSuccess), len(fs.Team.Marketing))
	fmt.Printf("Monthly Team Cost: $%s\n", formatFounderCurrency(fs.MonthlyTeamCost))

	// Show funding history
	if len(fs.FundingRounds) > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("💰 FUNDING HISTORY")
		fmt.Println(strings.Repeat("─", 70))
		for _, round := range fs.FundingRounds {
			fmt.Printf("%s: $%s @ $%s valuation (%.1f%% equity, %s terms) - Month %d\n",
				round.RoundName,
				formatFounderCurrency(round.Amount),
				formatFounderCurrency(round.Valuation),
				round.EquityGiven,
				round.Terms,
				round.Month)
		}
	}

	// Decision menu
	handleFounderDecisions(fs)

	// Process month
	messages := fs.ProcessMonth()

	// Display month results
	if len(messages) > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		cyan.Println("📰 MONTHLY UPDATE")
		fmt.Println(strings.Repeat("─", 70))
		for _, msg := range messages {
			fmt.Println(msg)
		}
	}

	fmt.Print("\nPress 'Enter' to continue to next month...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleFounderDecisions(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🎯 DECISIONS FOR THIS MONTH")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Println("\n1. Hire Team Member")
	fmt.Println("2. Fire Team Member")
	fmt.Println("3. Spend on Marketing")
	fmt.Println("4. Raise Funding Round")
	fmt.Println("5. Skip (Do Nothing)")

	fmt.Print("\nWhat would you like to do? (1-5): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		handleHiring(fs)
	case "2":
		handleFiring(fs)
	case "3":
		handleMarketing(fs)
	case "4":
		handleFundraising(fs)
	case "5":
		fmt.Println("\n✓ Focusing on operations this month...")
	default:
		fmt.Println("\nInvalid choice, skipping...")
	}
}

func handleHiring(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💼 HIRING")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Println("\n1. Engineer ($100k/year) - Builds product, reduces churn")
	fmt.Println("2. Sales Rep ($100k/year) - Increases customer acquisition")
	fmt.Println("3. Customer Success ($100k/year) - Reduces churn")
	fmt.Println("4. Marketing ($100k/year) - Supports customer acquisition")
	fmt.Println("5. Cancel")

	fmt.Print("\nWho would you like to hire? (1-5): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var role founder.EmployeeRole
	switch choice {
	case "1":
		role = founder.RoleEngineer
	case "2":
		role = founder.RoleSales
	case "3":
		role = founder.RoleCustomerSuccess
	case "4":
		role = founder.RoleMarketing
	case "5":
		fmt.Println("\nCanceled hiring")
		return
	default:
		color.Red("\nInvalid choice!")
		return
	}

	err := fs.HireEmployee(role)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Hired a new %s!", role)
		fmt.Printf("New runway: %d months\n", fs.CashRunwayMonths)
	}
}

func handleFiring(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	if fs.Team.TotalEmployees == 0 {
		color.Yellow("\nYou have no employees to fire.")
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("⚠️  LAYOFFS")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Printf("\n1. Engineer (current: %d)\n", len(fs.Team.Engineers))
	fmt.Printf("2. Sales Rep (current: %d)\n", len(fs.Team.Sales))
	fmt.Printf("3. Customer Success (current: %d)\n", len(fs.Team.CustomerSuccess))
	fmt.Printf("4. Marketing (current: %d)\n", len(fs.Team.Marketing))
	fmt.Println("5. Cancel")

	fmt.Print("\nWho would you like to let go? (1-5): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var role founder.EmployeeRole
	switch choice {
	case "1":
		role = founder.RoleEngineer
	case "2":
		role = founder.RoleSales
	case "3":
		role = founder.RoleCustomerSuccess
	case "4":
		role = founder.RoleMarketing
	case "5":
		fmt.Println("\nCanceled")
		return
	default:
		color.Red("\nInvalid choice!")
		return
	}

	err := fs.FireEmployee(role)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Let go one %s", role)
		fmt.Printf("New runway: %d months\n", fs.CashRunwayMonths)
	}
}

func handleMarketing(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("📣 MARKETING SPEND")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Printf("\nCurrent Cash: $%s\n", formatFounderCurrency(fs.Cash))
	fmt.Printf("Competition Level: %s\n", fs.CompetitionLevel)
	
	var cac int64
	switch fs.CompetitionLevel {
	case "very_high":
		cac = 10000
	case "high":
		cac = 7500
	case "medium":
		cac = 5000
	default:
		cac = 3000
	}
	fmt.Printf("Estimated CAC: $%s per customer\n", formatFounderCurrency(cac))

	fmt.Print("\nHow much to spend on marketing? (enter amount or 0 to cancel): $")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(strings.ReplaceAll(amountStr, ",", ""))
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount < 0 {
		color.Red("\nInvalid amount!")
		return
	}

	if amount == 0 {
		fmt.Println("\nCanceled")
		return
	}

	if amount > fs.Cash {
		color.Red("\n❌ Not enough cash!")
		return
	}

	newCustomers := fs.SpendOnMarketing(amount)
	color.Green("\n✓ Marketing campaign launched!")
	color.Green("  Acquired %d new customers!", newCustomers)
	fmt.Printf("  New MRR: $%s\n", formatFounderCurrency(fs.MRR))
}

func handleFundraising(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💰 RAISE FUNDING")
	fmt.Println(strings.Repeat("─", 70))
	
	// Determine what rounds are available
	hasSeed := false
	hasSeriesA := false
	hasSeriesB := false
	
	for _, round := range fs.FundingRounds {
		if round.RoundName == "Seed" {
			hasSeed = true
		}
		if round.RoundName == "Series A" {
			hasSeriesA = true
		}
		if round.RoundName == "Series B" {
			hasSeriesB = true
		}
	}

	fmt.Println("\nAvailable Rounds:")
	options := []string{}
	if !hasSeed {
		fmt.Println("1. Seed Round ($2-5M)")
		options = append(options, "Seed")
	}
	if hasSeed && !hasSeriesA {
		fmt.Println("2. Series A ($10-20M)")
		options = append(options, "Series A")
	}
	if hasSeriesA && !hasSeriesB {
		fmt.Println("3. Series B ($30-50M)")
		options = append(options, "Series B")
	}
	fmt.Println("0. Cancel")

	if len(options) == 0 {
		color.Yellow("\nNo more funding rounds available!")
		return
	}

	fmt.Print("\nWhich round? (0-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "0" {
		fmt.Println("\nCanceled")
		return
	}

	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(options) {
		color.Red("\nInvalid choice!")
		return
	}

	roundName := options[choiceNum-1]
	success, amount, terms, equityGiven := fs.RaiseFunding(roundName)

	if !success {
		color.Red("\n❌ Failed to raise funding!")
		return
	}

	color.Green("\n✓ Successfully raised %s!", roundName)
	color.Green("  Amount: $%s", formatFounderCurrency(amount))
	fmt.Printf("  Valuation: $%s\n", formatFounderCurrency(fs.FundingRounds[len(fs.FundingRounds)-1].Valuation))
	fmt.Printf("  Equity Given: %.1f%%\n", equityGiven)
	fmt.Printf("  Terms: %s\n", terms)
	fmt.Printf("  Your remaining equity: %.1f%%\n", 100.0-fs.EquityGivenAway)
	fmt.Printf("  New runway: %d months\n", fs.CashRunwayMonths)
}

func displayAcquisitionOffer(fs *founder.FounderState, offer *founder.AcquisitionOffer) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	fmt.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("🎉 ACQUISITION OFFER!")
	fmt.Println(strings.Repeat("=", 70))

	yellow.Printf("\n%s wants to acquire your company!\n", offer.Acquirer)
	green.Printf("\nOffer Amount: $%s\n", formatFounderCurrency(offer.OfferAmount))
	
	founderPayout := int64(float64(offer.OfferAmount) * (100.0 - fs.EquityGivenAway) / 100.0)
	green.Printf("Your Payout (%.1f%% equity): $%s\n", 100.0-fs.EquityGivenAway, formatFounderCurrency(founderPayout))
	
	fmt.Printf("\nDue Diligence: %s\n", offer.DueDiligence)
	fmt.Printf("Terms Quality: %s\n", offer.TermsQuality)

	fmt.Print("\nAccept this offer? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice == "y" || choice == "yes" {
		fs.Cash = founderPayout
		fs.Turn = fs.MaxTurns + 1 // End game
		color.Green("\n🎉 Congratulations! You've successfully exited!")
	} else {
		color.Yellow("\n✓ Declined the offer. Continuing to build...")
	}
}

func displayFounderFinalScore(fs *founder.FounderState) {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	outcome, valuation, founderEquity := fs.GetFinalScore()

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    GAME OVER - FINAL RESULTS")
	cyan.Println(strings.Repeat("=", 70))

	fmt.Printf("\n%s\n", ascii.Trophy)
	yellow.Printf("Founder: %s\n", fs.FounderName)
	yellow.Printf("Company: %s\n", fs.CompanyName)
	fmt.Printf("\nMonths Played: %d\n", fs.Turn)
	
	// Outcome
	fmt.Println("\n" + strings.Repeat("─", 70))
	if fs.Cash <= 0 {
		red.Printf("OUTCOME: %s\n", outcome)
	} else {
		green.Printf("OUTCOME: %s\n", outcome)
	}
	fmt.Println(strings.Repeat("─", 70))

	// Final metrics
	fmt.Printf("\n💰 Final Cash: $%s\n", formatFounderCurrency(fs.Cash))
	fmt.Printf("📊 Company Valuation: $%s\n", formatFounderCurrency(valuation))
	fmt.Printf("📈 MRR: $%s\n", formatFounderCurrency(fs.MRR))
	fmt.Printf("👥 Customers: %d\n", fs.Customers)
	fmt.Printf("💼 Your Equity: %.1f%%\n", founderEquity)
	
	if founderEquity > 0 && valuation > 0 {
		yourValue := int64(float64(valuation) * founderEquity / 100.0)
		green.Printf("💎 Your Equity Value: $%s\n", formatFounderCurrency(yourValue))
	}

	// Team
	fmt.Printf("\n👥 Final Team Size: %d\n", fs.Team.TotalEmployees)
	fmt.Printf("   Engineers: %d | Sales: %d | CS: %d | Marketing: %d\n",
		len(fs.Team.Engineers), len(fs.Team.Sales),
		len(fs.Team.CustomerSuccess), len(fs.Team.Marketing))

	// Funding
	if len(fs.FundingRounds) > 0 {
		fmt.Println("\n💰 Funding Rounds:")
		totalRaised := int64(0)
		for _, round := range fs.FundingRounds {
			fmt.Printf("   %s: $%s (%.1f%% equity)\n",
				round.RoundName,
				formatFounderCurrency(round.Amount),
				round.EquityGiven)
			totalRaised += round.Amount
		}
		green.Printf("   Total Raised: $%s\n", formatFounderCurrency(totalRaised))
	}
}

func formatFounderCurrency(amount int64) string {
	if amount < 0 {
		return fmt.Sprintf("-$%s", formatFounderCurrency(-amount))
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

