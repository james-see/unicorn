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
	cyan := color.New(color.FgCyan, color.Bold)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🎯 DECISIONS FOR THIS MONTH")
	fmt.Println(strings.Repeat("─", 70))
	
	cyan.Println("\n[TEAM & OPERATIONS]")
	fmt.Println("1. Hire Team Member")
	fmt.Println("2. Fire Team Member")
	fmt.Println("3. Spend on Marketing")
	
	cyan.Println("\n[FUNDING & EQUITY]")
	fmt.Println("4. Raise Funding Round")
	if len(fs.FundingRounds) > 0 && fs.MRR > fs.MonthlyTeamCost {
		fmt.Println("5. Buy Back Equity (profitable companies only)")
	}
	fmt.Println("6. Manage Board & Equity Pool")
	
	cyan.Println("\n[STRATEGIC]")
	fmt.Println("7. Start Partnership")
	if fs.AffiliateProgram == nil {
		fmt.Println("8. Launch Affiliate Program")
	}
	if len(fs.Competitors) > 0 {
		fmt.Println("9. Handle Competitors")
	}
	fmt.Println("10. Expand to New Market")
	fmt.Println("11. Execute Pivot/Strategy Change")
	
	fmt.Println("\n0. Skip (Do Nothing)")

	fmt.Print("\nWhat would you like to do? (0-11): ")
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
		handleBuyback(fs)
	case "6":
		handleBoardAndEquity(fs)
	case "7":
		handlePartnership(fs)
	case "8":
		handleAffiliateLaunch(fs)
	case "9":
		handleCompetitorManagement(fs)
	case "10":
		handleGlobalExpansion(fs)
	case "11":
		handlePivot(fs)
	case "0":
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
	
	fs.UpdateCAC() // Recalculate current CAC
	
	fmt.Printf("\nCurrent Cash: $%s\n", formatFounderCurrency(fs.Cash))
	fmt.Printf("Base CAC (your business): $%s\n", formatFounderCurrency(fs.BaseCAC))
	fmt.Printf("Current Effective CAC: $%s\n", formatFounderCurrency(fs.CustomerAcquisitionCost))
	fmt.Printf("  Product Maturity: %.0f%% (reduces CAC up to 40%%)\n", fs.ProductMaturity*100)
	fmt.Printf("  Competition: %s (impacts CAC)\n", fs.CompetitionLevel)
	
	cacReduction := (1.0 - float64(fs.CustomerAcquisitionCost)/float64(fs.BaseCAC)) * 100
	if cacReduction > 0 {
		color.Green("  ✓ CAC reduced by %.0f%% from base!\n", cacReduction)
	}

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

func handlePartnership(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🤝 STRATEGIC PARTNERSHIPS")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Println("\n1. Distribution Partnership ($50-150k) - 10-30% MRR boost")
	fmt.Println("2. Technology Partnership ($30-100k) - Product integration, reduce churn")
	fmt.Println("3. Co-Marketing Partnership ($25-75k) - 15-40% MRR boost")
	fmt.Println("4. Data Partnership ($40-100k) - Analytics/insights, reduce churn")
	fmt.Println("0. Cancel")

	fmt.Print("\nSelect partnership type (0-4): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var partnerType string
	switch choice {
	case "1":
		partnerType = "distribution"
	case "2":
		partnerType = "technology"
	case "3":
		partnerType = "co-marketing"
	case "4":
		partnerType = "data"
	case "0":
		fmt.Println("\nCanceled")
		return
	default:
		color.Red("\nInvalid choice!")
		return
	}

	partnership, err := fs.StartPartnership(partnerType)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Partnership with %s started!", partnership.Partner)
		fmt.Printf("  Type: %s\n", partnership.Type)
		fmt.Printf("  Cost: $%s\n", formatFounderCurrency(partnership.Cost))
		fmt.Printf("  Duration: %d months\n", partnership.Duration)
		fmt.Printf("  Expected MRR Boost: $%s\n", formatFounderCurrency(partnership.MRRBoost))
	}
}

func handleAffiliateLaunch(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💰 LAUNCH AFFILIATE PROGRAM")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Println("\nAffiliate programs let partners sell your product for commission.")
	fmt.Println("Setup Cost: $20-50k | Monthly Platform Fees: $5-10k")
	fmt.Println("\nRecommended commission rates:")
	fmt.Println("  • 10-15% for SaaS products")
	fmt.Println("  • 15-20% for marketplaces")
	fmt.Println("  • 20-30% for high-margin products")

	fmt.Print("\nSet commission rate (5-30%): ")
	commStr, _ := reader.ReadString('\n')
	commStr = strings.TrimSpace(commStr)
	commission, err := strconv.ParseFloat(commStr, 64)
	if err != nil || commission < 5 || commission > 30 {
		color.Red("\nInvalid commission rate!")
		return
	}

	err = fs.LaunchAffiliateProgram(commission)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Affiliate program launched!")
		fmt.Printf("  Commission: %.1f%%\n", commission)
		fmt.Printf("  Starting Affiliates: %d\n", fs.AffiliateProgram.Affiliates)
	}
}

func handleCompetitorManagement(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	reader := bufio.NewReader(os.Stdin)

	if len(fs.Competitors) == 0 {
		color.Yellow("\nNo active competitors at this time")
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("⚔️  COMPETITOR MANAGEMENT")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Println("\nActive Competitors:")
	activeCount := 0
	for i, comp := range fs.Competitors {
		if !comp.Active {
			continue
		}
		activeCount++
		threatColor := color.New(color.FgYellow)
		if comp.Threat == "high" {
			threatColor = color.New(color.FgRed)
		} else if comp.Threat == "low" {
			threatColor = color.New(color.FgGreen)
		}
		
		fmt.Printf("%d. %s - ", i+1, comp.Name)
		threatColor.Printf("Threat: %s", comp.Threat)
		fmt.Printf(" | Market Share: %.1f%% | Strategy: %s\n", comp.MarketShare*100, comp.Strategy)
	}

	if activeCount == 0 {
		color.Yellow("\nNo active competitors")
		return
	}

	fmt.Print("\nSelect competitor # to handle (0 to cancel): ")
	compStr, _ := reader.ReadString('\n')
	compStr = strings.TrimSpace(compStr)
	compNum, err := strconv.Atoi(compStr)
	if err != nil || compNum == 0 {
		fmt.Println("\nCanceled")
		return
	}
	compIndex := compNum - 1

	fmt.Println("\nStrategies:")
	fmt.Println("1. Ignore - No cost, but they may take market share")
	fmt.Println("2. Compete Aggressively - $50-150k, reduce their threat")
	fmt.Println("3. Partner With Them - $100-250k, merge customer bases")

	fmt.Print("\nSelect strategy (1-3): ")
	stratChoice, _ := reader.ReadString('\n')
	stratChoice = strings.TrimSpace(stratChoice)

	var strategy string
	switch stratChoice {
	case "1":
		strategy = "ignore"
	case "2":
		strategy = "compete"
	case "3":
		strategy = "partner"
	default:
		red.Println("\nInvalid choice!")
		return
	}

	message, err := fs.HandleCompetitor(compIndex, strategy)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ %s", message)
	}
}

func handleGlobalExpansion(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🌍 GLOBAL EXPANSION")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Println("\nAvailable Markets:")
	fmt.Println("1. Europe - $200k setup (~27 initial customers), $30k/mo, high competition")
	fmt.Println("2. Asia - $250k setup (~25 initial customers), $40k/mo, very high competition")
	fmt.Println("3. LATAM - $150k setup (~30 initial customers), $20k/mo, medium competition")
	fmt.Println("4. Middle East - $180k setup (~60 initial customers), $25k/mo, low competition")
	fmt.Println("5. Africa - $120k setup (~40 initial customers), $15k/mo, low competition")
	fmt.Println("6. Australia - $100k setup (~20 initial customers), $18k/mo, medium competition")
	fmt.Println("0. Cancel")
	
	yellow.Println("\nⓘ  Initial customers = Setup Cost ÷ Local CAC")
	fmt.Println("   Without CS team & immature product: ~50% monthly churn!")
	fmt.Println("   Competitors in market will steal customers if ignored")

	// Show already launched markets
	if len(fs.GlobalMarkets) > 0 {
		fmt.Println("\nAlready Operating In:")
		for _, m := range fs.GlobalMarkets {
			fmt.Printf("  • %s (%.1f%% penetration, $%s MRR)\n", 
				m.Region, m.Penetration*100, formatFounderCurrency(m.MRR))
		}
	}

	fmt.Print("\nSelect market (0-6): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var region string
	switch choice {
	case "1":
		region = "Europe"
	case "2":
		region = "Asia"
	case "3":
		region = "LATAM"
	case "4":
		region = "Middle East"
	case "5":
		region = "Africa"
	case "6":
		region = "Australia"
	case "0":
		fmt.Println("\nCanceled")
		return
	default:
		color.Red("\nInvalid choice!")
		return
	}

	market, err := fs.ExpandToMarket(region)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Launched in %s!", market.Region)
		fmt.Printf("  Setup Cost: $%s\n", formatFounderCurrency(market.SetupCost))
		fmt.Printf("  Initial Customers: %d\n", market.CustomerCount)
		fmt.Printf("  Initial MRR: $%s\n", formatFounderCurrency(market.MRR))
		fmt.Printf("  Monthly Operating Cost: $%s\n", formatFounderCurrency(market.MonthlyCost))
		fmt.Printf("  Market Size: %d potential customers\n", market.MarketSize)
		fmt.Printf("  Initial Penetration: %.2f%%\n", market.Penetration*100)
		
		yellow.Println("\n⚠️  Your global churn rate increased due to operational complexity")
		fmt.Printf("  New churn rate: %.1f%%\n", fs.CustomerChurnRate*100)
	}
}

func handlePivot(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🔄 PIVOT / STRATEGY CHANGE")
	fmt.Println(strings.Repeat("─", 70))
	
	red.Println("\n⚠️  WARNING: Pivots are expensive and risky!")
	fmt.Printf("Current Strategy: %s\n", fs.StartupType)
	fmt.Println("\nExpected Costs:")
	fmt.Println("  • $100-300k in pivot costs")
	fmt.Println("  • Lose 20-50% of customers")
	fmt.Println("  • 30-70% success rate (depends on product maturity & timing)")

	fmt.Println("\nNew Strategy Options:")
	fmt.Println("1. Enterprise B2B")
	fmt.Println("2. SMB B2B")
	fmt.Println("3. B2C")
	fmt.Println("4. Marketplace")
	fmt.Println("5. Platform")
	fmt.Println("6. Vertical SaaS")
	fmt.Println("7. Horizontal SaaS")
	fmt.Println("8. Deep Tech")
	fmt.Println("9. Consumer Apps")
	fmt.Println("0. Cancel")

	fmt.Print("\nSelect new strategy (0-9): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	strategies := map[string]string{
		"1": "Enterprise B2B",
		"2": "SMB B2B",
		"3": "B2C",
		"4": "Marketplace",
		"5": "Platform",
		"6": "Vertical SaaS",
		"7": "Horizontal SaaS",
		"8": "Deep Tech",
		"9": "Consumer Apps",
	}

	toStrategy, ok := strategies[choice]
	if !ok {
		if choice == "0" {
			fmt.Println("\nCanceled")
		} else {
			red.Println("\nInvalid choice!")
		}
		return
	}

	fmt.Print("\nReason for pivot: ")
	reason, _ := reader.ReadString('\n')
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "Strategic repositioning"
	}

	fmt.Print("\nAre you sure? This is risky! (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "yes" && confirm != "y" {
		fmt.Println("\nCanceled")
		return
	}

	pivot, err := fs.ExecutePivot(toStrategy, reason)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		if pivot.Success {
			color.Green("\n🎉 Pivot SUCCESSFUL!")
			fmt.Printf("  New Strategy: %s\n", pivot.ToStrategy)
			fmt.Printf("  New Market Size: %d (expanded!)\n", fs.TargetMarketSize)
		} else {
			red.Println("\n😞 Pivot FAILED")
			fmt.Println("  The market didn't respond well to the change")
		}
		fmt.Printf("  Cost: $%s\n", formatFounderCurrency(pivot.Cost))
		fmt.Printf("  Customers Lost: %d\n", pivot.CustomersLost)
	}
}

func handleBuyback(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	if len(fs.FundingRounds) == 0 {
		color.Yellow("\nNo funding rounds to buy back from")
		return
	}

	monthlyProfit := fs.MRR - fs.MonthlyTeamCost
	if monthlyProfit <= 0 {
		color.Red("\nMust be profitable to buy back equity")
		fmt.Printf("Current monthly profit/loss: $%s\n", formatFounderCurrency(monthlyProfit))
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💎 BUY BACK EQUITY")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Printf("\nYou're profitable! Monthly profit: $%s\n", formatFounderCurrency(monthlyProfit))
	fmt.Println("\nFunding Rounds:")
	for i, round := range fs.FundingRounds {
		fmt.Printf("%d. %s - %.1f%% equity given\n", i+1, round.RoundName, round.EquityGiven)
	}

	fmt.Print("\nSelect round to buy back from (0 to cancel): ")
	roundStr, _ := reader.ReadString('\n')
	roundStr = strings.TrimSpace(roundStr)
	roundNum, err := strconv.Atoi(roundStr)
	if err != nil || roundNum == 0 {
		fmt.Println("\nCanceled")
		return
	}
	if roundNum < 1 || roundNum > len(fs.FundingRounds) {
		color.Red("\nInvalid round!")
		return
	}

	selectedRound := fs.FundingRounds[roundNum-1]
	currentVal := int64(float64(fs.MRR) * 12 * 12) // 12x ARR

	fmt.Printf("\nCurrent valuation: $%s (12x ARR)\n", formatFounderCurrency(currentVal))
	fmt.Printf("Available to buy back: %.1f%%\n", selectedRound.EquityGiven)

	fmt.Print("\nHow much equity to buy back? (%): ")
	equityStr, _ := reader.ReadString('\n')
	equityStr = strings.TrimSpace(equityStr)
	equity, err := strconv.ParseFloat(equityStr, 64)
	if err != nil || equity <= 0 {
		color.Red("\nInvalid percentage!")
		return
	}

	costEstimate := int64(float64(currentVal) * equity / 100.0)
	fmt.Printf("\nEstimated cost: $%s\n", formatFounderCurrency(costEstimate))
	fmt.Print("Confirm? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "yes" && confirm != "y" {
		fmt.Println("\nCanceled")
		return
	}

	buyback, err := fs.BuybackEquity(selectedRound.RoundName, equity)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Successfully bought back %.1f%% equity!", buyback.EquityBought)
		fmt.Printf("  Paid: $%s\n", formatFounderCurrency(buyback.PricePaid))
		fmt.Printf("  Your new ownership: %.1f%%\n", 100.0-fs.EquityGivenAway)
	}
}

func handleBoardAndEquity(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("👔 BOARD & EQUITY POOL")
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Printf("\nCurrent Board Seats: %d\n", fs.BoardSeats)
	fmt.Printf("Employee Equity Pool: %.1f%%\n", fs.EquityPool)
	fmt.Printf("Your Equity: %.1f%%\n", 100.0-fs.EquityGivenAway)

	fmt.Println("\nOptions:")
	fmt.Println("1. Add Board Seat (costs ~2% from equity pool)")
	fmt.Println("2. Expand Equity Pool (dilutes you by 1-10%)")
	fmt.Println("0. Cancel")

	fmt.Print("\nSelect option (0-2): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Print("\nReason for new board seat: ")
		reason, _ := reader.ReadString('\n')
		reason = strings.TrimSpace(reason)
		if reason == "" {
			reason = "Strategic advisor"
		}

		err := fs.AddBoardSeat(reason)
		if err != nil {
			color.Red("\n❌ Error: %v", err)
		} else {
			color.Green("\n✓ Added board seat")
			fmt.Printf("  New board seats: %d\n", fs.BoardSeats)
			fmt.Printf("  Remaining equity pool: %.1f%%\n", fs.EquityPool)
		}

	case "2":
		fmt.Print("\nHow much to add to equity pool? (1-10%): ")
		pctStr, _ := reader.ReadString('\n')
		pctStr = strings.TrimSpace(pctStr)
		pct, err := strconv.ParseFloat(pctStr, 64)
		if err != nil || pct < 1 || pct > 10 {
			color.Red("\nInvalid percentage!")
			return
		}

		err = fs.ExpandEquityPool(pct)
		if err != nil {
			color.Red("\n❌ Error: %v", err)
		} else {
			color.Green("\n✓ Expanded equity pool by %.1f%%", pct)
			fmt.Printf("  New equity pool: %.1f%%\n", fs.EquityPool)
			fmt.Printf("  Your equity: %.1f%%\n", 100.0-fs.EquityGivenAway)
		}

	case "0":
		fmt.Println("\nCanceled")
	default:
		color.Red("\nInvalid choice!")
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

