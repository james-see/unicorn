package ui

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	achievements "github.com/jamesacampbell/unicorn/achievements"
	animations "github.com/jamesacampbell/unicorn/animations"
	ascii "github.com/jamesacampbell/unicorn/ascii"
	clear "github.com/jamesacampbell/unicorn/clear"
	db "github.com/jamesacampbell/unicorn/database"
	founder "github.com/jamesacampbell/unicorn/founder"
	leaderboard "github.com/jamesacampbell/unicorn/leaderboard"
	progression "github.com/jamesacampbell/unicorn/progression"
	upgrades "github.com/jamesacampbell/unicorn/upgrades"
)

func PlayFounderMode(username string) {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Load available startups
	startups, err := founder.LoadFounderStartups("founder/startups.json")
	if err != nil {
		color.Red("Error loading startups: %v", err)
		fmt.Print("\nPress 'Enter' to return to main menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// Step 1: Category selection
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("              SELECT STARTUP CATEGORY")
	cyan.Println(strings.Repeat("=", 70))

	// Get unique categories from startups
	categoryMap := make(map[string]int)
	for _, s := range startups {
		categoryMap[s.Type]++
	}

	// Build ordered list of categories
	categories := []string{"SaaS", "DeepTech", "FinTech", "HealthTech", "GovTech"}
	var availableCategories []string
	for _, cat := range categories {
		if count, exists := categoryMap[cat]; exists {
			availableCategories = append(availableCategories, cat)
			green.Printf("\n%d. %s", len(availableCategories), cat)
			fmt.Printf(" (%d companies)\n", count)
		}
	}

	fmt.Print("\nSelect category (1-" + fmt.Sprintf("%d", len(availableCategories)) + "): ")
	reader := bufio.NewReader(os.Stdin)
	categoryChoice, _ := reader.ReadString('\n')
	categoryChoice = strings.TrimSpace(categoryChoice)
	categoryNum, err := strconv.Atoi(categoryChoice)
	if err != nil || categoryNum < 1 || categoryNum > len(availableCategories) {
		color.Red("Invalid choice!")
		fmt.Print("\nPress 'Enter' to return to main menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	selectedCategory := availableCategories[categoryNum-1]

	// Step 2: Filter startups by category
	var filteredStartups []founder.StartupTemplate
	for _, s := range startups {
		if s.Type == selectedCategory {
			filteredStartups = append(filteredStartups, s)
		}
	}

	clear.ClearIt()

	// Display startup selection
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Printf("              CHOOSE YOUR %s STARTUP\n", strings.ToUpper(selectedCategory))
	cyan.Println(strings.Repeat("=", 70))

	// Randomize competition levels and cash for display (and later use)
	competitionLevels := []string{"low", "medium", "high", "very_high"}
	rand.Seed(time.Now().UnixNano()) // Seed randomization

	for i, s := range filteredStartups {
		// Randomize competition level for display
		randomCompetition := competitionLevels[rand.Intn(len(competitionLevels))]

		// Randomize cash (±20%)
		cashVariance := 0.20
		cashMultiplier := 1.0 + (rand.Float64()*cashVariance*2 - cashVariance) // 0.8 to 1.2
		randomizedCash := int64(float64(s.InitialCash) * cashMultiplier)

		// Randomize CAC (±25%)
		cacVariance := 0.25
		cacMultiplier := 1.0 + (rand.Float64()*cacVariance*2 - cacVariance) // 0.75 to 1.25
		randomizedCAC := int64(float64(s.BaseCAC) * cacMultiplier)

		yellow.Printf("\n%d. %s", i+1, s.Name)
		fmt.Printf(" [%s]\n", s.Type)
		fmt.Printf("   %s\n", s.Tagline)
		fmt.Printf("   Cash: $%s | MRR: $%s | Customers: %d\n",
			formatFounderCurrency(randomizedCash),
			formatFounderCurrency(s.InitialMRR),
			s.InitialCustomers)
		fmt.Printf("   CAC: $%s | Market Size: %s | Competition: %s\n",
			formatFounderCurrency(randomizedCAC),
			formatFounderNumber(int64(s.TargetMarketSize)),
			randomCompetition)

		// Store randomized values back in template for use when game starts
		filteredStartups[i].CompetitionLevel = randomCompetition
		filteredStartups[i].InitialCash = randomizedCash
		filteredStartups[i].BaseCAC = randomizedCAC
	}

	fmt.Print("\nSelect your startup (1-" + fmt.Sprintf("%d", len(filteredStartups)) + "): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(filteredStartups) {
		color.Red("Invalid choice!")
		fmt.Print("\nPress 'Enter' to return to main menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	selectedStartup := filteredStartups[choiceNum-1]
	clear.ClearIt()

	// Get player upgrades
	playerUpgrades, err := db.GetPlayerUpgrades(username)
	if err != nil {
		playerUpgrades = []string{}
	}

	// Initialize founder game
	fs := founder.NewFounderGame(username, selectedStartup, playerUpgrades)

	// Welcome message
	displayFounderWelcome(fs)

	// Main game loop
	for !fs.IsGameOver() {
		playFounderTurn(fs)
	}

	// Show final score
	displayFounderFinalScore(fs)

	// Save score to database and check achievements
	saveFounderScoreAndCheckAchievements(fs)

	// Submit to global leaderboard
	askToSubmitFounderToGlobalLeaderboard(fs)

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
	if fs.CashRunwayMonths < 0 {
		fmt.Printf("⏱️  Runway: ∞ (cash positive!)\n")
	} else {
		fmt.Printf("⏱️  Runway: %d months\n", fs.CashRunwayMonths)
	}
	fmt.Printf("📈 Product Maturity: %.0f%%\n", fs.ProductMaturity*100)
	
	// Display active upgrades (filtered for Founder mode)
	founderUpgrades := upgrades.FilterUpgradeIDsForGameMode(fs.PlayerUpgrades, "founder")
	if len(founderUpgrades) > 0 {
		green := color.New(color.FgGreen)
		fmt.Println("\n" + strings.Repeat("─", 70))
		green.Println("✨ ACTIVE UPGRADES FOR THIS GAME:")
		fmt.Println(strings.Repeat("─", 70))
		for _, upgradeID := range founderUpgrades {
			if upgrade, exists := upgrades.AllUpgrades[upgradeID]; exists {
				fmt.Printf("  %s %s\n", upgrade.Icon, upgrade.Name)
			}
		}
	}

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

	// Check for low cash warning (calculate accurate runway)
	netRevenue := int64(float64(fs.MRR) * 0.67)
	totalExpenses := fs.MonthlyTeamCost + fs.MonthlyComputeCost + fs.MonthlyODCCost
	if fs.AffiliateProgram != nil {
		totalExpenses += fs.AffiliateProgram.MonthlyPlatformFee
		totalExpenses += int64(float64(fs.AffiliateMRR) * fs.AffiliateProgram.Commission)
	}
	for _, m := range fs.GlobalMarkets {
		totalExpenses += m.MonthlyCost
	}
	netIncome := netRevenue - totalExpenses

	if fs.NeedsLowCashWarning() {
		red.Println("\n⚠️  WARNING: Cash is running low!")
		if netIncome > 0 {
			red.Printf("   Cash: $%s | Runway: ∞ (profitable)\n", formatFounderCurrency(fs.Cash))
		} else {
			burnRate := -netIncome
			runway := int(float64(fs.Cash) / float64(burnRate))
			red.Printf("   Cash: $%s | Runway: %d months\n", formatFounderCurrency(fs.Cash), runway)
		}
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

	// Calculate actual net income for accurate runway display
	netRevenue2 := int64(float64(fs.MRR) * 0.67) // After 33% deductions
	totalExpenses2 := fs.MonthlyTeamCost + fs.MonthlyComputeCost + fs.MonthlyODCCost
	if fs.AffiliateProgram != nil {
		totalExpenses2 += fs.AffiliateProgram.MonthlyPlatformFee
		totalExpenses2 += int64(float64(fs.AffiliateMRR) * fs.AffiliateProgram.Commission)
	}
	for _, m := range fs.GlobalMarkets {
		totalExpenses2 += m.MonthlyCost
	}
	netIncome2 := netRevenue2 - totalExpenses2

	fmt.Printf("💰 Cash: $%s", formatFounderCurrency(fs.Cash))
	if netIncome2 > 0 {
		green.Printf(" (∞ runway - profitable!)\n")
	} else if netIncome2 < 0 {
		burnRate := -netIncome2
		runway := int(float64(fs.Cash) / float64(burnRate))
		if runway == 0 {
			red.Printf(" (⚠️  <1 month runway)\n")
		} else if runway <= 3 {
			red.Printf(" (⚠️  %d months runway)\n", runway)
		} else if runway <= 6 {
			yellow.Printf(" (%d months runway)\n", runway)
		} else {
			fmt.Printf(" (%d months runway)\n", runway)
		}
	} else {
		fmt.Printf(" (break-even)\n")
	}

	fmt.Printf("📈 MRR: $%s", formatFounderCurrency(fs.MRR))
	// Only show growth rate if we have MRR and it's not the first month
	if fs.MRR > 0 && fs.Turn > 1 {
		if fs.MonthlyGrowthRate > 0 {
			green.Printf(" (↑%.1f%%)\n", fs.MonthlyGrowthRate*100)
		} else if fs.MonthlyGrowthRate < 0 {
			red.Printf(" (↓%.1f%%)\n", -fs.MonthlyGrowthRate*100)
		} else {
			fmt.Println(" (flat)")
		}
	} else {
		fmt.Println() // Just newline, no growth indicator
	}

	// Show MRR breakdown if affiliate program is active
	if fs.AffiliateProgram != nil && fs.AffiliateMRR > 0 {
		fmt.Printf("   Direct: $%s | Affiliate: $%s\n",
			formatFounderCurrency(fs.DirectMRR), formatFounderCurrency(fs.AffiliateMRR))
	}

	fmt.Printf("👥 Customers: %d", fs.Customers)
	if fs.AffiliateProgram != nil && fs.AffiliateCustomers > 0 {
		fmt.Printf(" (Direct: %d | Affiliate: %d)", fs.DirectCustomers, fs.AffiliateCustomers)
	}
	fmt.Println()

	// Customer section
	if fs.TotalCustomersEver > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("👥 CUSTOMERS")
		fmt.Println(strings.Repeat("─", 70))
		fmt.Printf("Active: %d | Churned: %d | Total Ever: %d\n",
			fs.Customers, fs.TotalChurned, fs.TotalCustomersEver)
		if fs.Customers > 0 {
			fmt.Printf("Active MRR: $%s | Avg Deal: $%s/mo",
				formatFounderCurrency(fs.MRR), formatFounderCurrency(fs.AvgDealSize))
			if fs.MinDealSize > 0 && fs.MaxDealSize > fs.MinDealSize {
				fmt.Printf(" | Range: $%s - $%s/mo\n",
					formatFounderCurrency(fs.MinDealSize), formatFounderCurrency(fs.MaxDealSize))
			} else {
				fmt.Println()
			}
		}
	} else {
		if fs.Customers > 0 {
			fmt.Printf("💸 Deal Size: $%s/mo avg", formatFounderCurrency(fs.AvgDealSize))
			if fs.MinDealSize > 0 && fs.MaxDealSize > fs.MinDealSize {
				fmt.Printf(" | Range: $%s - $%s/mo\n",
					formatFounderCurrency(fs.MinDealSize), formatFounderCurrency(fs.MaxDealSize))
			} else {
				fmt.Println()
			}
		} else {
			// Show expected deal size when no customers yet
			fmt.Printf("💸 Expected Deal Size: $%s/mo\n", formatFounderCurrency(fs.AvgDealSize))
		}
	}
	fmt.Printf("🔄 Churn Rate: %.1f%% | ", fs.CustomerChurnRate*100)
	fmt.Printf("📦 Product Maturity: %.0f%%\n", fs.ProductMaturity*100)
	fmt.Printf("	💼 Your Equity: %.1f%%", 100.0-fs.EquityPool-fs.EquityGivenAway)
	if fs.EquityPool > 0 {
		fmt.Printf(" | Employee Pool: %.1f%% (%.1f%% used, %.1f%% available)\n",
			fs.EquityPool, fs.EquityAllocated, fs.EquityPool-fs.EquityAllocated)
	} else {
		fmt.Println()
	}

	// Show key metrics if we have meaningful data
	if fs.Turn > 1 && fs.Customers > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("📊 KEY METRICS")
		fmt.Println(strings.Repeat("─", 70))

		ltvCac := fs.CalculateLTVToCAC()
		cacPayback := fs.CalculateCACPayback()
		ruleOf40 := fs.CalculateRuleOf40()
		burnMultiple := fs.CalculateBurnMultiple()
		magicNumber := fs.CalculateMagicNumber()

		if ltvCac > 0 {
			statusIcon := "🟢"
			if ltvCac < 1 {
				statusIcon = "🔴"
			} else if ltvCac < 3 {
				statusIcon = "🟡"
			}
			fmt.Printf("%s LTV:CAC Ratio: %.1f:1", statusIcon, ltvCac)
			if ltvCac < 1 {
				fmt.Print(" (losing money per customer!)")
			} else if ltvCac < 3 {
				fmt.Print(" (aim for 3:1+)")
			} else {
				fmt.Print(" (healthy)")
			}
			fmt.Println()
		}

		if cacPayback > 0 {
			statusIcon := "🟢"
			if cacPayback > 18 {
				statusIcon = "🔴"
			} else if cacPayback > 12 {
				statusIcon = "🟡"
			}
			fmt.Printf("%s CAC Payback: %.1f months", statusIcon, cacPayback)
			if cacPayback > 18 {
				fmt.Print(" (too long!)")
			} else if cacPayback > 12 {
				fmt.Print(" (aim for <12)")
			} else {
				fmt.Print(" (good)")
			}
			fmt.Println()
		}

		if fs.MRR > 50000 {
			statusIcon := "🟢"
			if ruleOf40 < 0 {
				statusIcon = "🔴"
			} else if ruleOf40 < 40 {
				statusIcon = "🟡"
			}
			fmt.Printf("%s Rule of 40: %.0f%%", statusIcon, ruleOf40)
			if ruleOf40 < 0 {
				fmt.Print(" (unprofitable & shrinking!)")
			} else if ruleOf40 < 40 {
				fmt.Print(" (aim for 40%+)")
			} else {
				fmt.Print(" (excellent)")
			}
			fmt.Println()
		}

		if burnMultiple > 0 && burnMultiple < 999 {
			statusIcon := "🟢"
			if burnMultiple > 2 {
				statusIcon = "🔴"
			} else if burnMultiple > 1.5 {
				statusIcon = "🟡"
			}
			fmt.Printf("%s Burn Multiple: %.1fx", statusIcon, burnMultiple)
			if burnMultiple > 2 {
				fmt.Print(" (burning too much per $ of growth!)")
			} else if burnMultiple > 1.5 {
				fmt.Print(" (room for improvement)")
			} else {
				fmt.Print(" (efficient)")
			}
			fmt.Println()
		}

		if magicNumber > 0 {
			statusIcon := "🟢"
			if magicNumber < 0.75 {
				statusIcon = "🔴"
			} else if magicNumber < 1.0 {
				statusIcon = "🟡"
			}
			fmt.Printf("%s Magic Number: %.2f", statusIcon, magicNumber)
			if magicNumber < 0.75 {
				fmt.Print(" (low sales efficiency)")
			} else if magicNumber < 1.0 {
				fmt.Print(" (aim for 1.0+)")
			} else {
				fmt.Print(" (strong sales efficiency)")
			}
			fmt.Println()
		}
	}

	// Show monthly highlights
	highlights := fs.GenerateMonthlyHighlights()
	if len(highlights) > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("⭐ MONTHLY HIGHLIGHTS")
		fmt.Println(strings.Repeat("─", 70))

		wins := []founder.MonthlyHighlight{}
		concerns := []founder.MonthlyHighlight{}

		for _, h := range highlights {
			if h.Type == "win" {
				wins = append(wins, h)
			} else {
				concerns = append(concerns, h)
			}
		}

		if len(wins) > 0 {
			green := color.New(color.FgGreen)
			green.Println("\n🎉 Wins:")
			for _, w := range wins {
				fmt.Printf("  • %s\n", w.Message)
			}
		}

		if len(concerns) > 0 {
			red := color.New(color.FgRed)
			red.Println("\n⚠️  Concerns:")
			for _, c := range concerns {
				fmt.Printf("  • %s\n", c.Message)
			}
		}
	}

	// Show customer health if we have customers
	if fs.Customers > 0 {
		healthy, atRisk, critical, atRiskMRR, criticalMRR := fs.GetCustomerHealthSegments()
		if atRisk > 0 || critical > 0 {
			fmt.Println("\n" + strings.Repeat("─", 70))
			yellow.Println("❤️  CUSTOMER HEALTH")
			fmt.Println(strings.Repeat("─", 70))

			if healthy > 0 {
				green := color.New(color.FgGreen)
				green.Printf("🟢 Healthy: %d customers (low churn risk)\n", healthy)
			}
			if atRisk > 0 {
				yellow.Printf("🟡 At Risk: %d customers ($%s MRR at risk)\n", atRisk, formatFounderCurrency(atRiskMRR))
			}
			if critical > 0 {
				red := color.New(color.FgRed)
				red.Printf("🔴 Critical: %d customers ($%s MRR likely to churn)\n", critical, formatFounderCurrency(criticalMRR))
			}

			// Advice
			if critical > 0 || atRisk > 0 {
				fmt.Println("\n💡 Tip: Hire CS team or COO to improve customer health")
			}
		}
	}

	// Show board sentiment if raised funding
	if len(fs.FundingRounds) > 0 && fs.BoardSentiment != "" {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("👔 BOARD / INVESTOR SENTIMENT")
		fmt.Println(strings.Repeat("─", 70))

		sentimentIcon := "😐"
		sentimentColor := color.New(color.FgYellow)
		switch fs.BoardSentiment {
		case "happy":
			sentimentIcon = "😄"
			sentimentColor = color.New(color.FgGreen)
		case "pleased":
			sentimentIcon = "🙂"
			sentimentColor = color.New(color.FgGreen)
		case "neutral":
			sentimentIcon = "😐"
			sentimentColor = color.New(color.FgYellow)
		case "concerned":
			sentimentIcon = "😟"
			sentimentColor = color.New(color.FgRed)
		case "angry":
			sentimentIcon = "😡"
			sentimentColor = color.New(color.FgRed)
		}

		sentimentColor.Printf("%s Sentiment: %s", sentimentIcon, fs.BoardSentiment)
		fmt.Printf(" | Pressure: %d/100\n", fs.BoardPressure)

		if fs.BoardPressure >= 75 {
			red := color.New(color.FgRed)
			red.Println("\n⚠️  Board is demanding faster progress - risk of founder replacement!")
		} else if fs.BoardPressure >= 50 {
			fmt.Println("\n💬 Board wants to see improved metrics soon")
		} else if fs.BoardPressure <= 25 {
			green := color.New(color.FgGreen)
			green.Println("\n✓ Board is very supportive - keep up the good work!")
		}
	}

	// Show strategic opportunity if one is pending
	if fs.PendingOpportunity != nil {
		opp := fs.PendingOpportunity
		fmt.Println("\n" + strings.Repeat("─", 70))
		cyan := color.New(color.FgCyan, color.Bold)
		cyan.Println("💡 STRATEGIC OPPORTUNITY")
		fmt.Println(strings.Repeat("─", 70))

		yellow.Printf("\n%s\n", opp.Title)
		fmt.Printf("\n%s\n", opp.Description)

		green := color.New(color.FgGreen)
		green.Printf("\n✓ Benefits: %s\n", opp.Benefit)

		red := color.New(color.FgRed)
		red.Printf("⚠️  Risks: %s\n", opp.Risk)

		if opp.Cost > 0 {
			fmt.Printf("💰 Cost: $%s\n", formatFounderCurrency(opp.Cost))
		}

		fmt.Printf("\n⏰ Expires in %d month(s) - decide at end of this turn!\n", opp.ExpiresIn)
	}

	// Show partnerships if any are active
	activePartnerships := []founder.Partnership{}
	for _, p := range fs.Partnerships {
		if p.Status == "active" {
			activePartnerships = append(activePartnerships, p)
		}
	}

	if len(activePartnerships) > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("🤝 ACTIVE PARTNERSHIPS")
		fmt.Println(strings.Repeat("─", 70))

		totalMRRBoost := int64(0)
		totalChurnReduction := 0.0

		for _, p := range activePartnerships {
			monthsRemaining := p.Duration - (fs.Turn - p.MonthStarted)
			fmt.Printf("• %s (%s) - %d months remaining\n", p.Partner, p.Type, monthsRemaining)

			// Show benefits based on partnership type
			switch p.Type {
			case "distribution":
				fmt.Printf("  └─ MRR Boost: +$%s/mo (customer acquisition)\n",
					formatFounderCurrency(p.MRRBoost))
				totalMRRBoost += p.MRRBoost
			case "technology":
				fmt.Printf("  └─ Churn Reduction: -%.1f%% | Product Integration Boost: +$%s/mo\n",
					p.ChurnReduction*100, formatFounderCurrency(p.MRRBoost))
				totalChurnReduction += p.ChurnReduction
				totalMRRBoost += p.MRRBoost
			case "co-marketing":
				fmt.Printf("  └─ MRR Boost: +$%s/mo (marketing reach)\n",
					formatFounderCurrency(p.MRRBoost))
				totalMRRBoost += p.MRRBoost
			case "data":
				fmt.Printf("  └─ Churn Reduction: -%.1f%% | Analytics Boost: +$%s/mo\n",
					p.ChurnReduction*100, formatFounderCurrency(p.MRRBoost))
				totalChurnReduction += p.ChurnReduction
				totalMRRBoost += p.MRRBoost
			}
		}

		// Show cumulative benefits
		if len(activePartnerships) > 1 {
			fmt.Printf("\nTotal Partnership Benefits:\n")
			if totalMRRBoost > 0 {
				fmt.Printf("  • Combined MRR Boost: +$%s/mo\n", formatFounderCurrency(totalMRRBoost))
			}
			if totalChurnReduction > 0 {
				fmt.Printf("  • Combined Churn Reduction: -%.1f%%\n", totalChurnReduction*100)
			}
		}
	}

	// Show team
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("👥 TEAM")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("Engineers: %d | Sales: %d | CS: %d | Marketing: %d | C-Suite: %d\n",
		len(fs.Team.Engineers), len(fs.Team.Sales),
		len(fs.Team.CustomerSuccess), len(fs.Team.Marketing), len(fs.Team.Executives))
	fmt.Printf("Monthly Team Cost: $%s\n", formatFounderCurrency(fs.MonthlyTeamCost))

	// Show infrastructure costs if there are customers
	if fs.Customers > 0 {
		fmt.Printf("💻 Compute Cost: $%s/mo | 📦 ODC: $%s/mo\n",
			formatFounderCurrency(fs.MonthlyComputeCost),
			formatFounderCurrency(fs.MonthlyODCCost))
	}

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
			if len(round.Investors) > 0 {
				cyan := color.New(color.FgCyan)
				cyan.Printf("   Investors: %s\n", strings.Join(round.Investors, ", "))
			}
		}
	}

	// Decision menu
	// Capture MRR before any decisions are made (for comparison)
	preDecisionMRR := fs.MRR
	handleFounderDecisions(fs)

	// Process month
	messages := fs.ProcessMonthWithBaseline(preDecisionMRR)

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

	// Loop until user makes an action (not just viewing data)
	for {
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
		} else {
			fmt.Println("8. View Affiliate Program")
		}
		if len(fs.Competitors) > 0 {
			fmt.Println("9. Handle Competitors")
		}
		fmt.Println("10. Expand to New Market")
		fmt.Println("11. Execute Pivot/Strategy Change")
		if fs.AffiliateProgram != nil {
			fmt.Println("11b. End Affiliate Program")
		}
		
		// Strategic opportunities (numbered sequentially)
		nextOption := 12
		if fs.PendingOpportunity != nil {
			green := color.New(color.FgGreen, color.Bold)
			green.Printf("%d. 💡 Respond to Strategic Opportunity\n", nextOption)
			nextOption++
		}

		cyan.Println("\n[VIEW DATA]")
		fmt.Printf("%d. View Team Roster\n", nextOption)
		nextOption++
		fmt.Printf("%d. View Customer Deals\n", nextOption)
		nextOption++
		if fs.Customers > 0 {
			fmt.Printf("%d. Solicit Customer Feedback\n", nextOption)
			nextOption++
		}
		fmt.Printf("%d. View Financials & Cash Flow\n", nextOption)
		nextOption++

		// Exit option (always last)
		exits := fs.GetAvailableExits()
		hasAvailableExit := false
		for _, exit := range exits {
			if exit.CanExit && exit.Type != "continue" {
				hasAvailableExit = true
				break
			}
		}
		if hasAvailableExit {
			green := color.New(color.FgGreen, color.Bold)
			green.Printf("%d. 💰 Consider Exit Options (IPO/Acquisition/Secondary)\n", nextOption)
		}

		fmt.Println("\n0. Skip (Do Nothing)")

		maxChoice := nextOption - 1
		if hasAvailableExit {
			maxChoice = nextOption
		}
		fmt.Printf("\nWhat would you like to do? (0-%d): ", maxChoice)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		// Check for "11b" to end affiliate program before parsing as integer
		if choice == "11b" && fs.AffiliateProgram != nil {
			shouldContinue := handleEndAffiliateProgram(fs)
			if shouldContinue {
				continue
			}
			break
		}

		// Parse choice and validate range
		choiceNum, err := strconv.Atoi(choice)
		if err != nil || choiceNum < 0 || choiceNum > maxChoice {
			fmt.Println("\nInvalid choice, skipping...")
			break
		}

		// Calculate view option start position
		viewOptionStart := 12
		if fs.PendingOpportunity != nil {
			viewOptionStart = 13
		}

		// Handle view options first (they don't advance the turn)
		if choiceNum >= viewOptionStart && choiceNum < nextOption {
			viewIndex := choiceNum - viewOptionStart
			if viewIndex == 0 {
				handleViewTeamRoster(fs)
			} else if viewIndex == 1 {
				handleViewCustomerDeals(fs)
			} else if viewIndex == 2 && fs.Customers > 0 {
				// Solicit Feedback - this is an action that can cancel
				shouldContinue := handleSolicitFeedback(fs)
				if shouldContinue {
					continue
				}
				break
			} else if viewIndex == 2 && fs.Customers <= 0 {
				handleViewFinancials(fs)
			} else if viewIndex == 3 && fs.Customers > 0 {
				handleViewFinancials(fs)
			} else {
				fmt.Println("\nInvalid choice, skipping...")
				break
			}
			continue // Loop back to menu for view options
		}

		// Handle exit option
		if choiceNum == maxChoice && hasAvailableExit {
			shouldContinue := handleExitOptions(fs)
			if shouldContinue {
				continue
			}
			break
		}

		// Action options - may return true to continue (cancelled) or false to process month
		shouldContinue := false
		switch choiceNum {
		case 1:
			shouldContinue = handleHiring(fs)
		case 2:
			shouldContinue = handleFiring(fs)
		case 3:
			shouldContinue = handleMarketing(fs)
		case 4:
			shouldContinue = handleFundraising(fs)
		case 5:
			shouldContinue = handleBuyback(fs)
		case 6:
			shouldContinue = handleBoardAndEquity(fs)
		case 7:
			shouldContinue = handlePartnership(fs)
		case 8:
			if fs.AffiliateProgram != nil {
				handleViewAffiliateProgram(fs)
				shouldContinue = true // Viewing doesn't advance
			} else {
				shouldContinue = handleAffiliateLaunch(fs)
			}
		case 9:
			shouldContinue = handleCompetitorManagement(fs)
		case 10:
			shouldContinue = handleGlobalExpansion(fs)
		case 11:
			shouldContinue = handlePivot(fs)
		case 12:
			if fs.PendingOpportunity != nil {
				shouldContinue = handleStrategicOpportunity(fs)
			} else {
				fmt.Println("\nInvalid choice, skipping...")
				break
			}
		case 0:
			fmt.Println("\n✓ Focusing on operations this month...")
			shouldContinue = false
		default:
			fmt.Println("\nInvalid choice, skipping...")
			break
		}

		// If cancelled, loop back to menu
		if shouldContinue {
			continue
		}

		// Exit loop and process month
		break
	}
}

func handleHiring(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💼 HIRING")
	fmt.Println(strings.Repeat("─", 70))

	fmt.Println("\n== INDIVIDUAL CONTRIBUTORS ==")
	fmt.Println("1. Engineer ($100k/year) - Builds product, reduces churn")
	fmt.Println("2. Sales Rep ($100k/year) - Increases customer acquisition")
	fmt.Println("3. Customer Success ($100k/year) - Reduces churn")
	fmt.Println("4. Marketing ($100k/year) - Supports customer acquisition")

	green.Println("\n== C-LEVEL EXECUTIVES (3x impact, $300k/year) ==")
	fmt.Println("5. CTO - Like hiring 3 engineers at once")
	fmt.Println("6. CGO (Chief Growth Officer) - Like hiring 3 sales reps")
	fmt.Println("7. COO (Chief Operating Officer) - Like hiring 3 CS reps")
	fmt.Println("8. CFO (Chief Financial Officer) - Reduces burn by 10%, 1 role only")

	fmt.Println("\n0. Cancel")

	fmt.Print("\nWho would you like to hire? (0-8): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var role founder.EmployeeRole
	var isExec bool
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
		role = founder.RoleCTO
		isExec = true
	case "6":
		role = founder.RoleCGO
		isExec = true
	case "7":
		role = founder.RoleCOO
		isExec = true
	case "8":
		role = founder.RoleCFO
		isExec = true
	case "0":
		fmt.Println("\nCanceled hiring")
		return true // Cancelled, go back to menu
	default:
		color.Red("\nInvalid choice!")
		return true // Invalid, go back to menu
	}

	// For sales, marketing, and CS roles, ask about market assignment if there are global markets
	var market string = "USA"
	if !isExec && (role == founder.RoleSales || role == founder.RoleMarketing || role == founder.RoleCustomerSuccess) && len(fs.GlobalMarkets) > 0 {
		fmt.Println("\n" + strings.Repeat("─", 70))
		yellow.Println("📍 MARKET ASSIGNMENT")
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("\nWhich market should this employee focus on?")
		fmt.Println("1. USA (home market)")
		
		marketOptions := []string{"USA"}
		optionNum := 2
		for _, m := range fs.GlobalMarkets {
			fmt.Printf("%d. %s (%d customers, $%s MRR)\n", optionNum, m.Region, m.CustomerCount, formatFounderCurrency(m.MRR))
			marketOptions = append(marketOptions, m.Region)
			optionNum++
		}
		fmt.Printf("%d. All Markets (works globally)\n", optionNum)
		marketOptions = append(marketOptions, "All")
		
		fmt.Print("\nSelect market (1-" + fmt.Sprintf("%d", optionNum) + "): ")
		marketChoice, _ := reader.ReadString('\n')
		marketChoice = strings.TrimSpace(marketChoice)
		
		marketNum, err := strconv.Atoi(marketChoice)
		if err != nil || marketNum < 1 || marketNum > len(marketOptions) {
			color.Red("\nInvalid choice! Defaulting to USA")
			market = "USA"
		} else {
			market = marketOptions[marketNum-1]
		}
	}

	var err error
	if !isExec && market != "USA" {
		err = fs.HireEmployeeWithMarket(role, market)
	} else {
		err = fs.HireEmployee(role)
	}
	
	if err != nil {
		color.Red("\n❌ Error: %v", err)
		return true // Error, go back to menu
	} else {
		if isExec {
			green.Printf("\n🎉 Hired a %s! (3x impact)\n", role)
			fmt.Printf("   Cost: $300k/year ($25k/month)\n")
		} else {
			color.Green("\n✓ Hired a new %s!", role)
			fmt.Printf("   Cost: $100k/year ($8.3k/month)\n")
			if market != "USA" {
				cyan := color.New(color.FgCyan)
				cyan.Printf("   Assigned to: %s\n", market)
			}
		}
		if fs.CashRunwayMonths < 0 {
			fmt.Printf("   Runway: ∞ (still profitable!)\n")
		} else {
			fmt.Printf("   New runway: %d months\n", fs.CashRunwayMonths)
		}
	}
	return false // Action taken
}

func handleFiring(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	if fs.Team.TotalEmployees == 0 {
		color.Yellow("\nYou have no employees to fire.")
		return true
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("⚠️  LAYOFFS")
	fmt.Println(strings.Repeat("─", 70))

	fmt.Println("\n== INDIVIDUAL CONTRIBUTORS ==")
	fmt.Printf("1. Engineer (current: %d)\n", len(fs.Team.Engineers))
	fmt.Printf("2. Sales Rep (current: %d)\n", len(fs.Team.Sales))
	fmt.Printf("3. Customer Success (current: %d)\n", len(fs.Team.CustomerSuccess))
	fmt.Printf("4. Marketing (current: %d)\n", len(fs.Team.Marketing))

	fmt.Println("\n== EXECUTIVES ==")
	execCount := make(map[founder.EmployeeRole]int)
	for _, exec := range fs.Team.Executives {
		execCount[exec.Role]++
	}
	fmt.Printf("5. CTO (current: %d)\n", execCount[founder.RoleCTO])
	fmt.Printf("6. CGO (current: %d)\n", execCount[founder.RoleCGO])
	fmt.Printf("7. COO (current: %d)\n", execCount[founder.RoleCOO])
	fmt.Printf("8. CFO (current: %d)\n", execCount[founder.RoleCFO])

	fmt.Println("\n0. Cancel")

	fmt.Print("\nWho would you like to let go? (0-8): ")
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
		role = founder.RoleCTO
	case "6":
		role = founder.RoleCGO
	case "7":
		role = founder.RoleCOO
	case "8":
		role = founder.RoleCFO
	case "0":
		fmt.Println("\nCanceled")
		return true
	default:
		color.Red("\nInvalid choice!")
		return true
	}

	err := fs.FireEmployee(role)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Let go one %s", role)
		if fs.CashRunwayMonths < 0 {
			fmt.Printf("Runway: ∞ (profitable!)\n")
		} else {
			fmt.Printf("New runway: %d months\n", fs.CashRunwayMonths)
		}
	}
	return false
}

func handleMarketing(fs *founder.FounderState) bool {
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
		return true
	}

	if amount == 0 {
		fmt.Println("\nCanceled")
		return true
	}

	if amount > fs.Cash {
		color.Red("\n❌ Not enough cash!")
		return true
	}

	newCustomers := fs.SpendOnMarketing(amount)
	color.Green("\n✓ Marketing campaign launched!")
	color.Green("  Acquired %d new customers!", newCustomers)
	fmt.Printf("  New MRR: $%s\n", formatFounderCurrency(fs.MRR))
	return false
}

func handleFundraising(fs *founder.FounderState) bool {
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
	optionNum := 1
	if !hasSeed {
		fmt.Printf("%d. Seed Round ($2-5M)\n", optionNum)
		options = append(options, "Seed")
		optionNum++
	}
	if hasSeed && !hasSeriesA {
		fmt.Printf("%d. Series A ($10-20M)\n", optionNum)
		options = append(options, "Series A")
		optionNum++
	}
	if hasSeriesA && !hasSeriesB {
		fmt.Printf("%d. Series B ($30-50M)\n", optionNum)
		options = append(options, "Series B")
		optionNum++
	}
	fmt.Println("0. Cancel")

	if len(options) == 0 {
		color.Yellow("\nNo more funding rounds available!")
		return true
	}

	fmt.Printf("\nWhich round? (0-%d): ", len(options))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "0" {
		fmt.Println("\nCanceled")
		return true
	}

	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(options) {
		color.Red("\nInvalid choice!")
		return true
	}

	roundName := options[choiceNum-1]

	// Generate term sheet options
	termSheets := fs.GenerateTermSheetOptions(roundName)
	if len(termSheets) == 0 {
		color.Red("\n❌ Unable to generate term sheets!")
		return true
	}

	// Check for fundraising advisor
	hasFoundraisingAdvisor := false
	var advisorName string
	for _, boardMember := range fs.BoardMembers {
		if boardMember.IsActive && (boardMember.Expertise == "fundraising" || boardMember.Expertise == "strategy") {
			hasFoundraisingAdvisor = true
			advisorName = boardMember.Name
			break
		}
	}

	// Display term sheet options
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("📄 TERM SHEET OPTIONS")
	fmt.Println(strings.Repeat("─", 70))
	
	currentEquity := 100.0 - fs.EquityGivenAway - fs.EquityPool
	fmt.Printf("\n💼 Your Current Equity: %.2f%%\n", currentEquity)
	if currentEquity >= 50.0 {
		green := color.New(color.FgGreen)
		green.Println("   ✓ You have majority control")
	} else {
		red := color.New(color.FgRed, color.Bold)
		red.Println("   ⚠️  You don't have majority control (investors can force exits)")
	}
	
	if hasFoundraisingAdvisor {
		green := color.New(color.FgGreen)
		green.Printf("\n💡 %s (advisor) helped improve these terms!\n", advisorName)
		fmt.Println("   → Better valuations and more favorable terms")
	}

	for i, sheet := range termSheets {
		newEquity := currentEquity - sheet.Equity
		fmt.Printf("\n%d. %s\n", i+1, sheet.Terms)
		fmt.Printf("   Amount: $%s\n", formatFounderCurrency(sheet.Amount))
		fmt.Printf("   Pre-Money Valuation: $%s\n", formatFounderCurrency(sheet.PreValuation))
		fmt.Printf("   Post-Money Valuation: $%s\n", formatFounderCurrency(sheet.PostValuation))
		fmt.Printf("   Equity Given: %.1f%%\n", sheet.Equity)
		fmt.Printf("   Your Equity After: %.2f%%", newEquity)
		if newEquity < 50.0 && currentEquity >= 50.0 {
			red := color.New(color.FgRed, color.Bold)
			red.Printf(" ⚠️  (Loses majority control!)\n")
		} else if newEquity < 50.0 {
			yellow := color.New(color.FgYellow)
			yellow.Printf(" ⚠️  (Already below 50%%)\n")
		} else {
			fmt.Println()
		}
		fmt.Printf("   → %s\n", sheet.Description)
	}
	fmt.Println("\n0. Cancel")

	fmt.Print("\nSelect term sheet (0-4): ")
	termChoice, _ := reader.ReadString('\n')
	termChoice = strings.TrimSpace(termChoice)

	if termChoice == "0" {
		fmt.Println("\nCanceled")
		return true
	}

	termNum, err := strconv.Atoi(termChoice)
	if err != nil || termNum < 1 || termNum > len(termSheets) {
		color.Red("\nInvalid choice!")
		return true
	}

	selectedSheet := termSheets[termNum-1]
	
	// Check if this will drop founder below 50% equity
	newEquity := currentEquity - selectedSheet.Equity
	
	if newEquity < 50.0 && currentEquity >= 50.0 {
		red := color.New(color.FgRed, color.Bold)
		yellow := color.New(color.FgYellow, color.Bold)
		fmt.Println()
		red.Println("⚠️  WARNING: This will drop you below 50% ownership!")
		yellow.Printf("   Current equity: %.2f%%\n", currentEquity)
		yellow.Printf("   After this round: %.2f%%\n", newEquity)
		fmt.Println("   Without majority control, investors can force exit decisions.")
		fmt.Println("   You may be forced to accept acquisition offers.")
		fmt.Print("\nContinue anyway? (y/n): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))
		if confirm != "y" && confirm != "yes" {
			fmt.Println("\nCanceled")
			return true
		}
	}
	
	success := fs.RaiseFundingWithTerms(roundName, selectedSheet)

	if !success {
		color.Red("\n❌ Failed to raise funding!")
		return true
	}

	color.Green("\n✓ Successfully raised %s!", roundName)
	color.Green("  Amount: $%s", formatFounderCurrency(selectedSheet.Amount))
	fmt.Printf("  Pre-Money Valuation: $%s\n", formatFounderCurrency(selectedSheet.PreValuation))
	fmt.Printf("  Post-Money Valuation: $%s\n", formatFounderCurrency(selectedSheet.PostValuation))
	fmt.Printf("  Equity Given: %.1f%%\n", selectedSheet.Equity)
	fmt.Printf("  Terms: %s\n", selectedSheet.Terms)
	fmt.Printf("  Your remaining equity: %.1f%%\n", 100.0-fs.EquityGivenAway-fs.EquityPool)
	if fs.CashRunwayMonths < 0 {
		fmt.Printf("  Runway: ∞ (profitable!)\n")
	} else {
		fmt.Printf("  New runway: %d months\n", fs.CashRunwayMonths)
	}
	return false
}

func displayAcquisitionOffer(fs *founder.FounderState, offer *founder.AcquisitionOffer) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)
	red := color.New(color.FgRed, color.Bold)

	fmt.Println("\n" + strings.Repeat("=", 70))
	
	// Check if this is a competitor acquisition
	if offer.IsCompetitor {
		red.Println("⚠️  COMPETITOR ACQUISITION OFFER!")
		yellow.Printf("\n%s wants to acquire your company!\n", offer.Acquirer)
		
		// Silicon Valley-specific messages
		if offer.Acquirer == "Hooli" || offer.Acquirer == "Gavin Belson's New Thing" {
			red.Println("\n🏢 Hooli Acquisition Offer!")
			fmt.Println("Gavin Belson wants to buy your company. This is typically a lowball offer")
			fmt.Println("to eliminate competition. Watch out for bad terms and due diligence issues.")
		} else if offer.Acquirer == "Nucleus" {
			yellow.Println("\n⚛️  Nucleus wants to acquire you!")
			fmt.Println("Your compression technology competitor wants to buy you out.")
		} else {
			fmt.Println("\nThis is a competitive acquisition - they're buying to eliminate competition.")
			fmt.Println("Competitor offers are typically lower than strategic acquirer offers.")
		}
	} else {
		cyan.Println("🎉 ACQUISITION OFFER!")
		yellow.Printf("\n%s wants to acquire your company!\n", offer.Acquirer)
	}
	
	fmt.Println(strings.Repeat("=", 70))
	green.Printf("\nOffer Amount: $%s\n", formatFounderCurrency(offer.OfferAmount))

	// Calculate and display payout breakdown for all cap table entities
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💰 ACQUISITION PAYOUT BREAKDOWN")
	fmt.Println(strings.Repeat("─", 70))
	
	founderEquity := 100.0 - fs.EquityGivenAway - fs.EquityPool
	founderPayout := int64(float64(offer.OfferAmount) * founderEquity / 100.0)
	
	green.Printf("\n%-40s %8.2f%%  $%s\n", "You (Founder)", founderEquity, formatFounderCurrency(founderPayout))
	
	// Calculate payouts for all cap table entries
	totalShown := founderEquity
	for _, entry := range fs.CapTable {
		payout := int64(float64(offer.OfferAmount) * entry.Equity / 100.0)
		totalShown += entry.Equity
		
		// Color code by type
		switch entry.Type {
		case "investor":
			white.Printf("%-40s %8.2f%%  $%s\n", entry.Name, entry.Equity, formatFounderCurrency(payout))
		case "executive":
			cyan.Printf("%-40s %8.2f%%  $%s\n", entry.Name+" (Executive)", entry.Equity, formatFounderCurrency(payout))
		case "employee":
			white.Printf("%-40s %8.2f%%  $%s\n", entry.Name+" (Employee)", entry.Equity, formatFounderCurrency(payout))
		case "advisor":
			white.Printf("%-40s %8.2f%%  $%s\n", entry.Name+" (Advisor)", entry.Equity, formatFounderCurrency(payout))
		default:
			white.Printf("%-40s %8.2f%%  $%s\n", entry.Name, entry.Equity, formatFounderCurrency(payout))
		}
	}
	
	// Show remaining equity pool
	if fs.EquityPool > 0 {
		yellow.Printf("%-40s %8.2f%%  (unallocated)\n", "Employee Pool", fs.EquityPool)
	}
	
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("%-40s %8.2f%%  $%s\n", "TOTAL", totalShown+fs.EquityPool, formatFounderCurrency(offer.OfferAmount))

	fmt.Printf("\nDue Diligence: %s\n", offer.DueDiligence)
	fmt.Printf("Terms Quality: %s\n", offer.TermsQuality)

	// Check if founder has majority ownership
	forcedAcceptance := founderEquity < 50.0

	if forcedAcceptance {
		red := color.New(color.FgRed, color.Bold)
		yellow := color.New(color.FgYellow, color.Bold)
		fmt.Println()
		red.Println("⚠️  WARNING: You don't have majority ownership!")
		yellow.Printf("   Your equity: %.2f%% (need 50%%+ for control)\n", founderEquity)
		fmt.Println("   Without majority control, the board can force this acquisition.")
		fmt.Println()
		yellow.Println("   The acquisition will proceed automatically...")
		fmt.Print("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		
		// Force acceptance
		fs.Cash = founderPayout
		fs.Turn = fs.MaxTurns + 1 // End game
		if offer.IsCompetitor {
			if offer.Acquirer == "Hooli" || offer.Acquirer == "Gavin Belson's New Thing" {
				red.Println("\n🏢 Hooli acquired your company!")
				fmt.Println("Gavin Belson now owns your startup. This is not ideal.")
			} else {
				red.Println("\n⚠️  Acquisition completed - Competitor took control!")
			}
		} else {
			color.Green("\n🎉 Acquisition completed - Board approved the exit!")
		}
		color.Yellow("   You received $%s (%.2f%% of $%s)", 
			formatFounderCurrency(founderPayout), founderEquity, formatFounderCurrency(offer.OfferAmount))
	} else {
		fmt.Print("\nAccept this offer? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice == "y" || choice == "yes" {
			fs.Cash = founderPayout
			fs.Turn = fs.MaxTurns + 1 // End game
			
			if offer.IsCompetitor {
				if offer.Acquirer == "Hooli" || offer.Acquirer == "Gavin Belson's New Thing" {
					yellow.Println("\n🏢 You accepted Hooli's acquisition offer.")
					fmt.Println("Gavin Belson now owns your company. At least you got paid...")
				} else {
					color.Green("\n🎉 You accepted the competitor acquisition offer!")
				}
			} else {
				color.Green("\n🎉 Congratulations! You've successfully exited!")
			}
		} else {
			color.Yellow("\n✓ Declined the offer. Continuing to build...")
		}
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
	} else if fs.HasExited {
		// Show exit-specific outcome
		switch fs.ExitType {
		case "ipo":
			green.Printf("OUTCOME: 🏛️  IPO SUCCESS - Took company public!\n")
		case "acquisition":
			green.Printf("OUTCOME: 🤝 ACQUIRED - Sold to strategic buyer!\n")
		case "secondary":
			green.Printf("OUTCOME: 💼 SECONDARY SALE - Partial liquidity event!\n")
		default:
			green.Printf("OUTCOME: %s\n", outcome)
		}
	} else {
		green.Printf("OUTCOME: %s\n", outcome)
	}
	fmt.Println(strings.Repeat("─", 70))

	// Final metrics - handle exit differently
	if fs.HasExited {
		magenta := color.New(color.FgMagenta, color.Bold)

		fmt.Printf("\n📅 Exit Month: %d of %d\n", fs.ExitMonth, fs.MaxTurns)
		fmt.Printf("📊 Exit Valuation: $%s\n", formatFounderCurrency(fs.ExitValuation))
		fmt.Printf("💼 Your Equity: %.1f%%\n", founderEquity)

		switch fs.ExitType {
		case "ipo":
			liquidation := int64(float64(fs.ExitValuation) * founderEquity * 0.20 / 100.0)
			remaining := int64(float64(fs.ExitValuation) * founderEquity * 0.80 / 100.0)
			magenta.Printf("💵 Immediate Liquidation (20%%): $%s\n", formatFounderCurrency(liquidation))
			magenta.Printf("💎 Remaining Equity Value: $%s\n", formatFounderCurrency(remaining))
			magenta.Printf("📈 Total Net Worth: $%s\n", formatFounderCurrency(liquidation+remaining))
		case "acquisition":
			payout := int64(float64(fs.ExitValuation) * founderEquity / 100.0)
			magenta.Printf("💵 Your Payout: $%s\n", formatFounderCurrency(payout))
			
			// Show complete cap table payout breakdown
			fmt.Println("\n" + strings.Repeat("─", 70))
			yellow.Println("💰 ACQUISITION PAYOUT BREAKDOWN")
			fmt.Println(strings.Repeat("─", 70))
			
			white := color.New(color.FgWhite)
			green.Printf("\n%-40s %8.2f%%  $%s\n", "You (Founder)", founderEquity, formatFounderCurrency(payout))
			
			// Calculate payouts for all cap table entries
			totalShown := founderEquity
			for _, entry := range fs.CapTable {
				entryPayout := int64(float64(fs.ExitValuation) * entry.Equity / 100.0)
				totalShown += entry.Equity
				
				// Color code by type
				switch entry.Type {
				case "investor":
					white.Printf("%-40s %8.2f%%  $%s\n", entry.Name, entry.Equity, formatFounderCurrency(entryPayout))
				case "executive":
					cyan.Printf("%-40s %8.2f%%  $%s\n", entry.Name+" (Executive)", entry.Equity, formatFounderCurrency(entryPayout))
				case "employee":
					white.Printf("%-40s %8.2f%%  $%s\n", entry.Name+" (Employee)", entry.Equity, formatFounderCurrency(entryPayout))
				case "advisor":
					white.Printf("%-40s %8.2f%%  $%s\n", entry.Name+" (Advisor)", entry.Equity, formatFounderCurrency(entryPayout))
				default:
					white.Printf("%-40s %8.2f%%  $%s\n", entry.Name, entry.Equity, formatFounderCurrency(entryPayout))
				}
			}
			
			// Show remaining equity pool
			if fs.EquityPool > 0 {
				yellow.Printf("%-40s %8.2f%%  (unallocated)\n", "Employee Pool", fs.EquityPool)
			}
			
			fmt.Println(strings.Repeat("─", 70))
			fmt.Printf("%-40s %8.2f%%  $%s\n", "TOTAL", totalShown+fs.EquityPool, formatFounderCurrency(fs.ExitValuation))
			
		case "secondary":
			sold := int64(float64(fs.ExitValuation) * founderEquity * 0.5 / 100.0)
			remaining := int64(float64(fs.ExitValuation) * founderEquity * 0.5 / 100.0)
			magenta.Printf("💵 Payout (50%% of stake): $%s\n", formatFounderCurrency(sold))
			magenta.Printf("💎 Remaining Equity Value: $%s\n", formatFounderCurrency(remaining))
			magenta.Printf("📈 Total Net Worth: $%s\n", formatFounderCurrency(sold+remaining))
		}

		fmt.Printf("\n📈 Final MRR: $%s\n", formatFounderCurrency(fs.MRR))
		fmt.Printf("👥 Customers: %d\n", fs.Customers)
	} else {
		fmt.Printf("\n💰 Final Cash: $%s\n", formatFounderCurrency(fs.Cash))
		fmt.Printf("📊 Company Valuation: $%s\n", formatFounderCurrency(valuation))
		fmt.Printf("📈 MRR: $%s\n", formatFounderCurrency(fs.MRR))
		fmt.Printf("👥 Customers: %d\n", fs.Customers)
		fmt.Printf("💼 Your Equity: %.1f%%\n", founderEquity)

		if founderEquity > 0 && valuation > 0 {
			yourValue := int64(float64(valuation) * founderEquity / 100.0)
			green.Printf("💎 Your Equity Value: $%s\n", formatFounderCurrency(yourValue))
		}
	}

	// Team
	fmt.Printf("\n👥 Final Team Size: %d\n", fs.Team.TotalEmployees)
	fmt.Printf("   Engineers: %d | Sales: %d | CS: %d | Marketing: %d | C-Suite: %d\n",
		len(fs.Team.Engineers), len(fs.Team.Sales),
		len(fs.Team.CustomerSuccess), len(fs.Team.Marketing), len(fs.Team.Executives))

	// Funding
	if len(fs.FundingRounds) > 0 {
		fmt.Println("\n💰 Funding Rounds:")
		totalRaised := int64(0)
		cyan := color.New(color.FgCyan)
		for _, round := range fs.FundingRounds {
			fmt.Printf("   %s: $%s (%.1f%% equity)\n",
				round.RoundName,
				formatFounderCurrency(round.Amount),
				round.EquityGiven)
			if len(round.Investors) > 0 {
				cyan.Printf("      Investors: %s\n", strings.Join(round.Investors, ", "))
			}
			totalRaised += round.Amount
		}
		green.Printf("   Total Raised: $%s\n", formatFounderCurrency(totalRaised))
	}
}

func askToSubmitFounderToGlobalLeaderboard(fs *founder.FounderState) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("           🏆 GLOBAL LEADERBOARD")
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

	// Calculate final metrics
	_, valuation, founderEquity := fs.GetFinalScore()

	// Calculate founder payout based on exit type
	var founderPayout int64
	if fs.HasExited {
		switch fs.ExitType {
		case "ipo":
			// 20% immediate liquidation
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.20 / 100.0)
		case "acquisition":
			// Full payout
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity / 100.0)
		case "secondary":
			// 50% sold
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.5 / 100.0)
		}
	} else {
		// No exit, equity value only
		founderPayout = int64(float64(valuation) * founderEquity / 100.0)
	}

	// Calculate max ARR (MRR * 12)
	maxARR := fs.MRR * 12

	// Calculate total funding raised
	totalFundingRaised := int64(0)
	for _, round := range fs.FundingRounds {
		totalFundingRaised += round.Amount
	}

	// Check API availability
	fmt.Print("\nChecking global leaderboard service...")
	if !leaderboard.IsAPIAvailable("") {
		color.Yellow("\n⚠️  Global leaderboard service is not available right now.")
		color.Yellow("Your score has been saved locally.")
		return
	}
	color.Green(" ✓")

	// Submit score
	fmt.Print("Submitting your score...")
	submission := leaderboard.FounderScoreSubmission{
		PlayerName:        fs.FounderName,
		FinalValuation:    valuation,
		FounderEquity:     founderEquity,
		FounderPayout:     founderPayout,
		ExitType:          fs.ExitType,
		ExitMonth:         fs.ExitMonth,
		MaxARR:            maxARR,
		StartupTemplate:   fs.CompanyName,
		FundingRaised:     totalFundingRaised,
		CustomersAcquired: fs.TotalCustomersEver,
	}

	err := leaderboard.SubmitFounderScore(submission, "")
	if err != nil {
		color.Red("\n❌ Failed to submit score: %v", err)
		color.Yellow("Your score has been saved locally.")
		return
	}

	color.Green(" ✓")
	cyan.Println("\n🎉 Success! Your score has been submitted to the global leaderboard!")
	yellow.Println("\nView the global leaderboard at:")
	yellow.Println("https://james-see.github.io/unicorn/#founder-leaderboard")
}

func handlePartnership(fs *founder.FounderState) bool {
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
		return true
	default:
		color.Red("\nInvalid choice!")
		return true
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
	return false
}

func handleAffiliateLaunch(fs *founder.FounderState) bool {
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
		return true
	}

	err = fs.LaunchAffiliateProgram(commission)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Affiliate program launched!")
		fmt.Printf("  Commission: %.1f%%\n", commission)
		fmt.Printf("  Starting Affiliates: %d\n", fs.AffiliateProgram.Affiliates)
	}
	return false
}

func handleCompetitorManagement(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	reader := bufio.NewReader(os.Stdin)

	if len(fs.Competitors) == 0 {
		color.Yellow("\nNo active competitors at this time")
		return true
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
		fmt.Printf(" | Market Share: %.1f%% | Strategy: %s", comp.MarketShare*100, comp.Strategy)
		
		// Add Silicon Valley flavor text
		if comp.Name == "Hooli" {
			fmt.Printf(" | 🏢 Tech Giant")
		} else if comp.Name == "Nucleus" {
			fmt.Printf(" | ⚛️  Compression Competitor")
		} else if comp.Name == "Gavin Belson's New Thing" {
			fmt.Printf(" | 💼 Gavin's Latest")
		} else if comp.Name == "Pied Piper" {
			fmt.Printf(" | 🎵 Parallel Universe")
		}
		fmt.Println()
	}

	if activeCount == 0 {
		color.Yellow("\nNo active competitors")
		return true
	}

	fmt.Print("\nSelect competitor # to handle (0 to cancel): ")
	compStr, _ := reader.ReadString('\n')
	compStr = strings.TrimSpace(compStr)
	compNum, err := strconv.Atoi(compStr)
	if err != nil || compNum == 0 {
		fmt.Println("\nCanceled")
		return true
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
		return true
	}

	message, err := fs.HandleCompetitor(compIndex, strategy)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ %s", message)
	}
	return false
}

func handleGlobalExpansion(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🌍 GLOBAL EXPANSION")
	fmt.Println(strings.Repeat("─", 70))

	// Get list of active markets
	activeMarkets := make(map[string]bool)
	for _, m := range fs.GlobalMarkets {
		activeMarkets[m.Region] = true
	}

	// Define all markets with their details
	type MarketOption struct {
		Region      string
		SetupCost   string
		MonthlyCost string
		Competition string
		Number      int
	}

	allMarkets := []MarketOption{
		{"Europe", "$200k", "$30k/mo", "high", 1},
		{"Asia", "$250k", "$40k/mo", "very high", 2},
		{"LATAM", "$150k", "$20k/mo", "medium", 3},
		{"Middle East", "$180k", "$25k/mo", "low", 4},
		{"Africa", "$120k", "$15k/mo", "low", 5},
		{"Australia", "$100k", "$18k/mo", "medium", 6},
	}

	// Filter out active markets and build menu
	availableMarkets := []MarketOption{}
	for _, market := range allMarkets {
		if !activeMarkets[market.Region] {
			availableMarkets = append(availableMarkets, market)
		}
	}

	if len(availableMarkets) == 0 {
		yellow.Println("\n✓ You've expanded to all available markets!")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	fmt.Println("\nAvailable Markets:")
	for _, market := range availableMarkets {
		fmt.Printf("%d. %s - %s setup, %s, %s competition\n",
			market.Number, market.Region, market.SetupCost, market.MonthlyCost, market.Competition)
	}
	fmt.Println("0. Cancel")

	if len(fs.GlobalMarkets) > 0 {
		yellow.Println("\n✓ Active Markets:")
		for _, m := range fs.GlobalMarkets {
			fmt.Printf("   %s - %d customers, $%s MRR\n",
				m.Region, m.CustomerCount, formatFounderCurrency(m.MRR))
		}
	}

	yellow.Println("\nⓘ  Initial customers = Setup Cost ÷ Local CAC")
	fmt.Println("   Sales/Marketing teams help grow new markets faster!")
	fmt.Println("   Without CS team & immature product: ~50% monthly churn!")

	fmt.Print("\nSelect market (1-%d or 0 to cancel): ", len(availableMarkets))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "0" {
		fmt.Println("\nCanceled")
		return true
	}

	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > 6 {
		color.Red("\nInvalid choice!")
		return true
	}

	// Map choice to region using the original number
	var region string
	for _, market := range availableMarkets {
		if market.Number == choiceNum {
			region = market.Region
			break
		}
	}

	if region == "" {
		color.Red("\nMarket not available or already operating in that region!")
		return true
	}

	// Count competitors before expansion
	competitorsBefore := len(fs.Competitors)
	
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

		yellow.Println("\n⚠️  IMPORTANT:")
		fmt.Printf("  • Global churn increased to %.1f%% (operational complexity)\n", fs.CustomerChurnRate*100)
		cyan := color.New(color.FgCyan)
		cyan.Println("  • Assign CS team to this market to reduce churn!")
		cyan.Println("  • Assign Sales/Marketing to grow this market faster!")
		
		// Show new competitors
		newCompetitors := len(fs.Competitors) - competitorsBefore
		if newCompetitors > 0 {
			yellow.Printf("\n🔴 %d new competitor(s) detected in %s!\n", newCompetitors, region)
			fmt.Println("  View them in the 'Handle Competitors' menu to choose your strategy")
			for i := competitorsBefore; i < len(fs.Competitors); i++ {
				comp := fs.Competitors[i]
				fmt.Printf("    • %s (Threat: %s, Market Share: %.1f%%)\n", comp.Name, comp.Threat, comp.MarketShare*100)
			}
		}
	}
	return false
}

func handlePivot(fs *founder.FounderState) bool {
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
		return true
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
		return true
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
	return false
}

func handleBuyback(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	reader := bufio.NewReader(os.Stdin)

	if len(fs.FundingRounds) == 0 {
		color.Yellow("\nNo funding rounds to buy back from")
		return true
	}

	monthlyProfit := fs.MRR - fs.MonthlyTeamCost
	if monthlyProfit <= 0 {
		color.Red("\nMust be profitable to buy back equity")
		fmt.Printf("Current monthly profit/loss: $%s\n", formatFounderCurrency(monthlyProfit))
		return true
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
		return true
	}
	if roundNum < 1 || roundNum > len(fs.FundingRounds) {
		color.Red("\nInvalid round!")
		return true
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
		return true
	}

	costEstimate := int64(float64(currentVal) * equity / 100.0)
	fmt.Printf("\nEstimated cost: $%s\n", formatFounderCurrency(costEstimate))
	fmt.Print("Confirm? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "yes" && confirm != "y" {
		fmt.Println("\nCanceled")
		return true
	}

	buyback, err := fs.BuybackEquity(selectedRound.RoundName, equity)
	if err != nil {
		color.Red("\n❌ Error: %v", err)
	} else {
		color.Green("\n✓ Successfully bought back %.1f%% equity!", buyback.EquityBought)
		fmt.Printf("  Paid: $%s\n", formatFounderCurrency(buyback.PricePaid))
		fmt.Printf("  Your new ownership: %.1f%%\n", 100.0-fs.EquityGivenAway-fs.EquityPool)
	}
	return false
}

func handleBoardAndEquity(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("👔 BOARD & EQUITY POOL")
	fmt.Println(strings.Repeat("─", 70))

	fmt.Printf("\nCurrent Board Seats: %d\n", fs.BoardSeats)
	fmt.Printf("Employee Equity Pool: %.1f%%\n", fs.EquityPool)
	fmt.Printf("Your Equity: %.1f%%\n", 100.0-fs.EquityPool-fs.EquityGivenAway)

	// Show current advisors/board members
	if len(fs.BoardMembers) > 0 {
		cyan.Println("\n👥 CURRENT ADVISORS:")
		for _, member := range fs.BoardMembers {
			if member.IsActive {
				role := "Advisor"
				if member.IsChairman {
					role = "👔 Chairman"
				}
				fmt.Printf("  • %s (%s, %s) - %.2f%% equity - %s\n",
					member.Name, member.Type, member.Expertise, member.EquityCost, role)
			}
		}
	}

	fmt.Println("\nOptions:")
	if len(fs.BoardMembers) > 0 || fs.BoardSeats > 1 {
		magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Println("0. View Board Table (visual board members display)")
	}
	fmt.Println("1. Add Board Seat (costs ~2% from equity pool)")
	fmt.Println("2. Expand Equity Pool (dilutes you by 1-10%)")
	green.Println("3. Add Advisor (0.25-1% equity for strategic guidance)")
	if len(fs.BoardMembers) > 0 {
		chairman := fs.GetChairman()
		red := color.New(color.FgRed)
		if chairman == nil {
			yellow.Println("4. Set Advisor as Chairman (requires additional 0.5-1x equity)")
		} else {
			red.Println("4. Remove Chairman (causes negative PR & board pressure)")
		}
		red.Println("5. Remove Advisor (with equity buyback option)")
		// Check if there are investor board members
		hasInvestorMembers := false
		for _, member := range fs.BoardMembers {
			if member.IsActive && member.Type == "investor" {
				hasInvestorMembers = true
				break
			}
		}
		if hasInvestorMembers {
			red.Println("6. Fire Board Member (investor - requires 51%+ ownership)")
		}
	}
	fmt.Println("9. Cancel")

	maxOption := "3"
	if len(fs.BoardMembers) > 0 {
		maxOption = "5"
		// Check if there are investor board members
		hasInvestorMembers := false
		for _, member := range fs.BoardMembers {
			if member.IsActive && member.Type == "investor" {
				hasInvestorMembers = true
				break
			}
		}
		if hasInvestorMembers {
			maxOption = "6"
		}
	}
	fmt.Printf("\nSelect option (0-%s): ", maxOption)
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

		fs.AddBoardSeat(reason)
		color.Green("\n✓ Added board seat")
		fmt.Printf("  New board seats: %d\n", fs.BoardSeats)
		fmt.Printf("  Remaining equity pool: %.1f%%\n", fs.EquityPool)

	case "2":
		fmt.Print("\nHow much to add to equity pool? (1-10%): ")
		pctStr, _ := reader.ReadString('\n')
		pctStr = strings.TrimSpace(pctStr)
		pct, err := strconv.ParseFloat(pctStr, 64)
		if err != nil || pct < 1 || pct > 10 {
			color.Red("\nInvalid percentage!")
			return true
		}

		// Check if this will drop founder below 50% equity
		currentEquity := 100.0 - fs.EquityGivenAway - fs.EquityPool
		newEquity := currentEquity - pct
		
		if newEquity < 50.0 && currentEquity >= 50.0 {
			red := color.New(color.FgRed, color.Bold)
			yellow := color.New(color.FgYellow, color.Bold)
			fmt.Println()
			red.Println("⚠️  WARNING: This will drop you below 50% ownership!")
			yellow.Printf("   Current equity: %.2f%%\n", currentEquity)
			yellow.Printf("   After expansion: %.2f%%\n", newEquity)
			fmt.Println("   Without majority control, investors can force exit decisions.")
			fmt.Print("\nContinue anyway? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))
			if confirm != "y" && confirm != "yes" {
				fmt.Println("\nCanceled")
				return true
			}
		}

		fs.ExpandEquityPool(pct)
		color.Green("\n✓ Expanded equity pool by %.1f%%", pct)
		fmt.Printf("  New equity pool: %.1f%%\n", fs.EquityPool)
		fmt.Printf("  Your equity: %.1f%%\n", 100.0-fs.EquityPool-fs.EquityGivenAway)

	case "3":
		handleAddAdvisor(fs)

	case "4":
		if len(fs.BoardMembers) == 0 {
			color.Red("\nInvalid choice!")
			return true
		}
		chairman := fs.GetChairman()
		if chairman == nil {
			handleSetChairman(fs)
		} else {
			handleRemoveChairman(fs)
		}

	case "5":
		if len(fs.BoardMembers) == 0 {
			color.Red("\nInvalid choice!")
			return true
		}
		handleRemoveAdvisor(fs)

	case "6":
		if len(fs.BoardMembers) == 0 {
			color.Red("\nInvalid choice!")
			return true
		}
		handleFireBoardMember(fs)

	case "0":
		handleViewBoardTable(fs)
		return true
	default:
		color.Red("\nInvalid choice!")
		return true
	}
	return false
}

func handleViewBoardTable(fs *founder.FounderState) {
	clear.ClearIt()
	
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta, color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	
	fmt.Println("\n" + strings.Repeat("=", 70))
	magenta.Println("                    🏛️  BOARD TABLE 🏛️")
	fmt.Println(strings.Repeat("=", 70))
	
	fmt.Printf("\nTotal Board Seats: %d\n", fs.BoardSeats)
	fmt.Printf("Your Equity: %.2f%%\n", 100.0-fs.EquityPool-fs.EquityGivenAway)
	fmt.Printf("Equity Pool: %.2f%%\n", fs.EquityPool)
	
	// Filter active board members
	activeMembers := []founder.BoardMember{}
	chairman := (*founder.BoardMember)(nil)
	for i := range fs.BoardMembers {
		if fs.BoardMembers[i].IsActive {
			activeMembers = append(activeMembers, fs.BoardMembers[i])
			if fs.BoardMembers[i].IsChairman {
				chairman = &fs.BoardMembers[i]
			}
		}
	}
	
	fmt.Println("\n" + strings.Repeat("-", 70))
	fmt.Println(ascii.BoardTable)
	fmt.Println(strings.Repeat("-", 70))
	
	if len(activeMembers) == 0 {
		yellow.Println("\n📭 No board members currently seated.")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	
	// Display board members around the table
	fmt.Println("\n" + strings.Repeat("-", 70))
	green.Println("Board Members:")
	fmt.Println(strings.Repeat("-", 70))
	
	// Show chairman first if exists
	if chairman != nil {
		magenta.Printf("\n👔 CHAIRMAN:\n")
		fmt.Printf("   %s (%s)\n", chairman.Name, chairman.Type)
		expertise := chairman.Expertise
		if len(expertise) > 0 {
			expertise = strings.ToUpper(expertise[:1]) + expertise[1:]
		}
		fmt.Printf("   Expertise: %s\n", expertise)
		fmt.Printf("   Equity: %.2f%%\n", chairman.EquityCost)
		if chairman.ContributionScore > 0 {
			green.Printf("   Contribution Score: %.0f%%\n", chairman.ContributionScore*100)
		}
		fmt.Println()
	}
	
	// Show other members
	otherMembers := []founder.BoardMember{}
	for _, member := range activeMembers {
		if !member.IsChairman {
			otherMembers = append(otherMembers, member)
		}
	}
	
	if len(otherMembers) > 0 {
		fmt.Println("Board Members:")
		for i, member := range otherMembers {
			memberType := member.Type
			if member.Type == "investor" {
				red.Printf("   %d. %s (%s)\n", i+1, member.Name, memberType)
			} else {
				fmt.Printf("   %d. %s (%s)\n", i+1, member.Name, memberType)
			}
			expertise := member.Expertise
			if len(expertise) > 0 {
				expertise = strings.ToUpper(expertise[:1]) + expertise[1:]
			}
			fmt.Printf("      Expertise: %s\n", expertise)
			fmt.Printf("      Equity: %.2f%%\n", member.EquityCost)
			if member.ContributionScore > 0 {
				green.Printf("      Contribution Score: %.0f%%\n", member.ContributionScore*100)
			}
			fmt.Println()
		}
	}
	
	// Show empty seats
	emptySeats := fs.BoardSeats - len(activeMembers)
	if emptySeats > 0 {
		yellow.Printf("\n📭 Empty Seats: %d\n", emptySeats)
	}
	
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleAddAdvisor(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💼 ADD ADVISOR")
	fmt.Println(strings.Repeat("─", 70))

	cyan.Println("\nSelect Advisor Expertise:")
	fmt.Println("1. Sales - Helps with customer acquisition & MRR growth")
	fmt.Println("2. Product - Improves product maturity faster")
	fmt.Println("3. Fundraising - Better terms on future rounds")
	fmt.Println("4. Operations - Reduces costs & burn")
	fmt.Println("5. Strategy - Reduces churn, strategic guidance")
	fmt.Println("\n💰 Cost: $10-50k setup fee + 0.25-1% equity")
	fmt.Println("   (50% chance of $2-8k/month retainer for less equity)")
	fmt.Println("\n0. Cancel")

	fmt.Print("\nSelect (0-5): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var expertise string
	switch choice {
	case "1":
		expertise = "sales"
	case "2":
		expertise = "product"
	case "3":
		expertise = "fundraising"
	case "4":
		expertise = "operations"
	case "5":
		expertise = "strategy"
	case "0":
		fmt.Println("\nCanceled")
		return
	default:
		color.Red("\nInvalid choice!")
		return
	}

	// Generate advisor - real Silicon Valley legends & tech advisors
	var names []string
	switch expertise {
	case "sales":
		names = []string{
			"Gary Vaynerchuk", // Gary V - marketing/sales legend
			"Aaron Ross",      // Predictable Revenue author
			"Jill Konrath",    // Sales strategist
			"Grant Cardone",   // Sales expert
			"Mark Roberge",    // Former HubSpot CRO
		}
	case "product":
		names = []string{
			"Marty Cagan",   // Product management guru
			"Julie Zhuo",    // Former Facebook VP of Design
			"Ken Norton",    // Google Ventures partner
			"April Dunford", // Positioning expert
			"Teresa Torres", // Continuous Discovery
		}
	case "fundraising":
		names = []string{
			"Marc Andreessen", // a16z founder
			"Naval Ravikant",  // AngelList founder
			"Jason Calacanis", // Angel investor
			"Chris Sacca",     // Lowercase Capital
			"Cyan Banister",   // Long Journey Ventures
		}
	case "operations":
		names = []string{
			"Sheryl Sandberg",       // Former Facebook COO
			"Keith Rabois",          // Founders Fund
			"Frank Slootman",        // Snowflake CEO
			"Claire Hughes Johnson", // Stripe COO
			"Gokul Rajaram",         // "Godfather of AdSense"
		}
	case "strategy":
		names = []string{
			"Reid Hoffman", // LinkedIn founder
			"Eric Ries",    // Lean Startup author
			"Steve Blank",  // Customer Development pioneer
			"Paul Graham",  // Y Combinator founder
			"Ben Horowitz", // a16z founder
		}
	}
	advisorName := names[rand.Intn(len(names))]

	// Advisor costs: equity + cash setup fee + optional monthly retainer
	equityCost := 0.25 + rand.Float64()*0.75    // 0.25-1% equity
	setupFee := int64(10000 + rand.Intn(40000)) // $10-50k one-time (recruiting, onboarding)

	// 50% chance of monthly retainer vs equity-only
	var monthlyRetainer int64
	if rand.Float64() < 0.5 {
		monthlyRetainer = int64(2000 + rand.Intn(6000)) // $2-8k/month ongoing
		equityCost = equityCost * 0.5                   // Lower equity if paying cash
	}

	// Check if we have enough equity
	availableEquity := fs.EquityPool - fs.EquityAllocated
	if availableEquity < equityCost {
		color.Red("\n❌ Not enough equity pool! Need %.2f%%, have %.2f%% available", equityCost, availableEquity)
		fmt.Println("   Expand your equity pool first.")
		return
	}

	// Check if we have enough cash for setup fee
	if fs.Cash < setupFee {
		color.Red("\n❌ Not enough cash! Need $%s for setup, have $%s",
			formatFounderCurrency(setupFee), formatFounderCurrency(fs.Cash))
		return
	}

	// Show cost breakdown and confirm
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Printf("💰 %s (%s Advisor)\n", advisorName, strings.Title(expertise))
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("Setup Fee: $%s (one-time)\n", formatFounderCurrency(setupFee))
	if monthlyRetainer > 0 {
		fmt.Printf("Monthly Retainer: $%s/month\n", formatFounderCurrency(monthlyRetainer))
	} else {
		cyan.Println("Monthly Retainer: Equity-only (no cash)")
	}
	fmt.Printf("Equity: %.2f%%\n", equityCost)
	fmt.Printf("\nAfter hiring:\n")
	fmt.Printf("  Cash: $%s → $%s\n",
		formatFounderCurrency(fs.Cash), formatFounderCurrency(fs.Cash-setupFee))
	fmt.Printf("  Employee Pool: %.1f%% → %.1f%%\n",
		fs.EquityPool-fs.EquityAllocated,
		fs.EquityPool-fs.EquityAllocated-equityCost)
	fmt.Printf("  Your Equity: %.1f%% (unchanged)\n",
		100.0-fs.EquityPool-fs.EquityGivenAway)

	fmt.Print("\nConfirm hire? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("\nCanceled")
		return
	}

	// Deduct costs
	fs.Cash -= setupFee

	advisor := founder.BoardMember{
		Name:              advisorName,
		Type:              "advisor",
		Expertise:         expertise,
		MonthAdded:        fs.Turn,
		EquityCost:        equityCost,
		IsActive:          true,
		ContributionScore: 0.5, // Start neutral
	}

	fs.BoardMembers = append(fs.BoardMembers, advisor)
	fs.EquityAllocated += equityCost
	fs.CapTable = append(fs.CapTable, founder.CapTableEntry{
		Name:         advisorName,
		Type:         "advisor",
		Equity:       equityCost,
		MonthGranted: fs.Turn,
	})

	color.Green("\n✓ Added %s as %s advisor!", advisorName, expertise)
	fmt.Printf("   Setup Cost: $%s\n", formatFounderCurrency(setupFee))
	if monthlyRetainer > 0 {
		fmt.Printf("   Monthly Retainer: $%s/month (added to team cost)\n", formatFounderCurrency(monthlyRetainer))
		fs.MonthlyTeamCost += monthlyRetainer
	}
	fmt.Printf("   Equity Cost: %.2f%% (from employee pool)\n", equityCost)
	fmt.Printf("   Remaining Pool: %.2f%%\n", fs.EquityPool-fs.EquityAllocated)
	fmt.Printf("   Your Equity: %.1f%% (unchanged)\n", 100.0-fs.EquityPool-fs.EquityGivenAway)
	fmt.Println("\n   💡 They will provide monthly guidance based on their expertise.")
}

func handleSetChairman(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("👔 SET ADVISOR AS CHAIRMAN")
	fmt.Println(strings.Repeat("─", 70))

	// List available advisors
	availableAdvisors := []founder.BoardMember{}
	for _, member := range fs.BoardMembers {
		if member.IsActive && !member.IsChairman {
			availableAdvisors = append(availableAdvisors, member)
		}
	}

	if len(availableAdvisors) == 0 {
		red.Println("\n❌ No advisors available to set as chairman!")
		return
	}

	fmt.Println("\nSelect advisor to promote to Chairman:")
	for i, advisor := range availableAdvisors {
		fmt.Printf("%d. %s (%s) - Current equity: %.2f%%\n",
			i+1, advisor.Name, advisor.Expertise, advisor.EquityCost)
	}
	fmt.Println("0. Cancel")

	fmt.Print("\nSelect (0-%d): ", len(availableAdvisors))
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 0 || choice > len(availableAdvisors) {
		color.Red("\nInvalid choice!")
		return
	}

	if choice == 0 {
		fmt.Println("\nCanceled")
		return
	}

	selectedAdvisor := availableAdvisors[choice-1]

	// Calculate additional equity needed
	currentEquity := selectedAdvisor.EquityCost
	additionalEquity := currentEquity * (0.5 + rand.Float64()*0.5) // 0.5-1x additional
	totalEquityNeeded := currentEquity + additionalEquity

	// Check available equity
	availableEquity := fs.EquityPool - fs.EquityAllocated
	if availableEquity < additionalEquity {
		red.Printf("\n❌ Insufficient equity pool! Need %.2f%%, have %.2f%% available\n", additionalEquity, availableEquity)
		fmt.Println("   Expand your equity pool first.")
		return
	}

	// Show cost breakdown
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Printf("💰 Chairman Promotion: %s\n", selectedAdvisor.Name)
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("Current Equity: %.2f%%\n", currentEquity)
	fmt.Printf("Additional Equity Required: %.2f%%\n", additionalEquity)
	fmt.Printf("Total Equity After Promotion: %.2f%%\n", totalEquityNeeded)
	fmt.Printf("\nBenefits:\n")
	green.Println("  • 60% chance of guidance (vs 30% for advisors)")
	green.Println("  • 2x impact on all benefits")
	green.Println("  • Crisis management (mitigates negative events)")
	green.Println("  • Investor relations (reduces board pressure)")
	green.Println("  • Represents company at events")
	fmt.Printf("\nTradeoffs:\n")
	red.Printf("  • Additional %.2f%% equity dilution\n", additionalEquity)
	red.Println("  • Higher monthly retainer ($5-15k vs $2-8k)")
	red.Println("  • Removing chairman causes negative PR")

	fmt.Print("\nConfirm promotion? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("\nCanceled")
		return
	}

	err = fs.SetChairman(selectedAdvisor.Name)
	if err != nil {
		red.Printf("\n❌ Error: %v\n", err)
		return
	}

	green.Printf("\n✓ %s is now Chairman of the Board!\n", selectedAdvisor.Name)
	fmt.Printf("  Additional equity: %.2f%%\n", additionalEquity)
	fmt.Printf("  Total equity: %.2f%%\n", totalEquityNeeded)
	fmt.Printf("  Remaining equity pool: %.2f%%\n", fs.EquityPool-fs.EquityAllocated)
}

func handleRemoveChairman(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	reader := bufio.NewReader(os.Stdin)

	chairman := fs.GetChairman()
	if chairman == nil {
		red.Println("\n❌ No chairman to remove!")
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	red.Println("⚠️  REMOVE CHAIRMAN")
	fmt.Println(strings.Repeat("─", 70))

	fmt.Printf("\nCurrent Chairman: %s (%s)\n", chairman.Name, chairman.Expertise)
	fmt.Printf("Current Equity: %.2f%%\n", chairman.EquityCost)

	red.Println("\n⚠️  WARNING: Removing chairman will cause:")
	red.Println("  • Negative PR and media attention")
	red.Println("  • Board pressure increase (20-30 points)")
	red.Println("  • Board sentiment deterioration")
	red.Println("  • Potential investor concerns")

	fmt.Printf("\nRemove %s as chairman? (y/n): ", chairman.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("\nCanceled")
		return
	}

	err := fs.RemoveChairman()
	if err != nil {
		red.Printf("\n❌ Error: %v\n", err)
		return
	}

	yellow.Printf("\n⚠️  %s removed as chairman\n", chairman.Name)
	red.Println("  Board pressure increased")
	red.Println("  Board sentiment deteriorated")
	fmt.Println("  (Chairman remains as advisor)")
}

func handleRemoveAdvisor(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	reader := bufio.NewReader(os.Stdin)

	// Get active advisors (not chairman, not investors)
	advisors := []founder.BoardMember{}
	for _, member := range fs.BoardMembers {
		if member.IsActive && member.Type == "advisor" && !member.IsChairman {
			advisors = append(advisors, member)
		}
	}

	if len(advisors) == 0 {
		red.Println("\n❌ No advisors to remove (or all are chairman/investors)")
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	red.Println("⚠️  REMOVE ADVISOR")
	fmt.Println(strings.Repeat("─", 70))

	cyan.Println("\nSelect advisor to remove:")
	for i, advisor := range advisors {
		fmt.Printf("%d. %s (%s) - %.2f%% equity\n", i+1, advisor.Name, advisor.Expertise, advisor.EquityCost)
	}
	fmt.Println("0. Cancel")

	fmt.Printf("\nSelect (0-%d): ", len(advisors))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 0 || choiceNum > len(advisors) {
		red.Println("\nInvalid choice!")
		return
	}

	if choiceNum == 0 {
		fmt.Println("\nCanceled")
		return
	}

	selectedAdvisor := advisors[choiceNum-1]

	// Ask about equity buyback
	fmt.Println("\n" + strings.Repeat("-", 70))
	yellow.Println("Equity Buyback Option:")
	fmt.Println("  • Buyback equity: Expensive but advisor keeps no equity")
	fmt.Println("  • No buyback: Advisor keeps equity but causes board pressure")
	
	// Calculate buyback cost
	estimatedValuation := fs.MRR * 15
	if estimatedValuation < 1000000 {
		estimatedValuation = 1000000
	}
	buybackCost := int64(float64(estimatedValuation) * (selectedAdvisor.EquityCost / 100.0))
	
	fmt.Printf("\nBuyback cost: $%s (%.2f%% of estimated $%s valuation)\n", 
		formatFounderCurrency(buybackCost), selectedAdvisor.EquityCost, formatFounderCurrency(estimatedValuation))
	fmt.Printf("Your cash: $%s\n", formatFounderCurrency(fs.Cash))
	
	fmt.Print("\nBuyback equity? (y/n): ")
	buybackChoice, _ := reader.ReadString('\n')
	buybackChoice = strings.TrimSpace(strings.ToLower(buybackChoice))
	buybackEquity := (buybackChoice == "y" || buybackChoice == "yes")

	if buybackEquity && buybackCost > fs.Cash {
		red.Printf("\n❌ Insufficient cash for buyback (need $%s)\n", formatFounderCurrency(buybackCost))
		return
	}

	// Confirm removal
	red.Printf("\n⚠️  WARNING: Removing %s will:\n", selectedAdvisor.Name)
	if buybackEquity {
		red.Printf("  • Cost $%s for equity buyback\n", formatFounderCurrency(buybackCost))
		green.Println("  • Return equity to you")
	} else {
		red.Println("  • Increase board pressure (10-20 points)")
		red.Println("  • Deteriorate board sentiment")
		yellow.Println("  • Advisor keeps their equity")
	}

	fmt.Printf("\nRemove %s? (y/n): ", selectedAdvisor.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("\nCanceled")
		return
	}

	err = fs.RemoveAdvisor(selectedAdvisor.Name, buybackEquity)
	if err != nil {
		red.Printf("\n❌ Error: %v\n", err)
		return
	}

	if buybackEquity {
		green.Printf("\n✅ %s removed (equity bought back for $%s)\n", selectedAdvisor.Name, formatFounderCurrency(buybackCost))
		green.Printf("  Equity returned: %.2f%%\n", selectedAdvisor.EquityCost)
	} else {
		yellow.Printf("\n⚠️  %s removed from board\n", selectedAdvisor.Name)
		red.Println("  Board pressure increased")
		red.Println("  Board sentiment deteriorated")
		fmt.Printf("  (%s keeps %.2f%% equity)\n", selectedAdvisor.Name, selectedAdvisor.EquityCost)
	}
}

func handleFireBoardMember(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	red := color.New(color.FgRed, color.Bold)
	reader := bufio.NewReader(os.Stdin)

	// Get investor board members
	investorMembers := []founder.BoardMember{}
	for _, member := range fs.BoardMembers {
		if member.IsActive && member.Type == "investor" {
			investorMembers = append(investorMembers, member)
		}
	}

	if len(investorMembers) == 0 {
		red.Println("\n❌ No investor board members to fire")
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	red.Println("🔥 FIRE BOARD MEMBER")
	fmt.Println(strings.Repeat("─", 70))

	// Check founder equity
	founderEquity := 100.0 - fs.EquityGivenAway - fs.EquityPool
	if founderEquity < 51.0 {
		red.Printf("\n❌ Cannot fire board members: You need 51%%+ ownership\n")
		yellow.Printf("   Your current equity: %.1f%%\n", founderEquity)
		return
	}

	cyan.Printf("\nYour equity: %.1f%% (majority owner)\n", founderEquity)
	cyan.Println("\nSelect investor board member to fire:")
	for i, member := range investorMembers {
		role := ""
		if member.IsChairman {
			role = " (Chairman)"
		}
		fmt.Printf("%d. %s - %.2f%% equity%s\n", i+1, member.Name, member.EquityCost, role)
	}
	fmt.Println("0. Cancel")

	fmt.Printf("\nSelect (0-%d): ", len(investorMembers))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 0 || choiceNum > len(investorMembers) {
		red.Println("\nInvalid choice!")
		return
	}

	if choiceNum == 0 {
		fmt.Println("\nCanceled")
		return
	}

	selectedMember := investorMembers[choiceNum-1]

	// Serious warning
	red.Printf("\n⚠️  ⚠️  ⚠️  SERIOUS CONSEQUENCES ⚠️  ⚠️  ⚠️\n")
	red.Printf("\nFiring %s will cause:\n", selectedMember.Name)
	red.Println("  • Massive board pressure increase (30-50 points)")
	red.Println("  • Board sentiment becomes ANGRY")
	red.Println("  • Negative PR and media attention")
	red.Println("  • Other investors will be very concerned")
	red.Println("  • May impact future fundraising")
	yellow.Printf("  • Board seats reduced by 1\n")
	yellow.Printf("  • %s keeps %.2f%% equity (cannot buyback)\n", selectedMember.Name, selectedMember.EquityCost)

	fmt.Printf("\n🔥 FIRE %s? (type 'FIRE' to confirm): ", selectedMember.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "FIRE" {
		fmt.Println("\nCanceled (must type 'FIRE' to confirm)")
		return
	}

	err = fs.FireBoardMember(selectedMember.Name)
	if err != nil {
		red.Printf("\n❌ Error: %v\n", err)
		return
	}

	red.Printf("\n🔥 %s FIRED from board\n", selectedMember.Name)
	red.Println("  Board pressure increased significantly")
	red.Println("  Board sentiment: ANGRY")
	red.Println("  Negative PR generated")
	yellow.Printf("  Board seats: %d\n", fs.BoardSeats)
	fmt.Printf("  (%s keeps %.2f%% equity)\n", selectedMember.Name, selectedMember.EquityCost)
}

func handleExitOptions(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💰 EXIT OPTIONS")
	fmt.Println(strings.Repeat("─", 70))

	exits := fs.GetAvailableExits()

	cyan.Println("\nAvailable Exit Paths:")
	fmt.Println()

	for i, exit := range exits {
		if exit.Type == "continue" {
			continue // Don't show "continue" as an exit option
		}

		optionNum := i + 1
		var icon string
		switch exit.Type {
		case "ipo":
			icon = "🏛️ "
		case "acquisition":
			icon = "🤝"
		case "secondary":
			icon = "💼"
		}

		fmt.Printf("%d. %s %s\n", optionNum, icon, strings.ToUpper(exit.Type))
		fmt.Printf("   %s\n", exit.Description)
		fmt.Printf("   Company Valuation: $%s\n", formatFounderCurrency(exit.Valuation))
		fmt.Printf("   Your Payout: $%s (%.1f%% equity)\n",
			formatFounderCurrency(exit.FounderPayout),
			(100.0 - fs.EquityPool - fs.EquityGivenAway))

		fmt.Println("   Requirements:")
		for _, req := range exit.Requirements {
			if strings.Contains(req, "❌") {
				red.Printf("     %s\n", req)
			} else {
				green.Printf("     %s\n", req)
			}
		}

		if exit.CanExit {
			green.Printf("   ✅ AVAILABLE\n")
		} else {
			red.Printf("   ❌ NOT AVAILABLE YET\n")
		}
		fmt.Println()
	}

	fmt.Println("0. Cancel (keep building)")
	fmt.Print("\nSelect exit option (0-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "0" {
		fmt.Println("\n✓ Continuing to build...")
		return true // Return to menu
	}

	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(exits) {
		color.Red("\nInvalid choice!")
		return true
	}

	selectedExit := exits[choiceNum-1]

	if !selectedExit.CanExit {
		color.Red("\n❌ This exit is not available yet. Keep building!")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	// Check if founder has majority ownership
	founderEquity := 100.0 - fs.EquityPool - fs.EquityGivenAway
	forcedExit := founderEquity < 50.0

	if forcedExit {
		red := color.New(color.FgRed, color.Bold)
		yellow := color.New(color.FgYellow, color.Bold)
		fmt.Println()
		red.Println("⚠️  WARNING: You don't have majority ownership!")
		yellow.Printf("   Your equity: %.2f%% (need 50%%+ for control)\n", founderEquity)
		fmt.Println("   Without majority control, investors can force exit decisions.")
		fmt.Println()
		yellow.Printf("   The %s exit will proceed automatically...\n", strings.ToUpper(selectedExit.Type))
		fmt.Print("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		
		// Execute exit automatically
		fs.ExecuteExit(selectedExit.Type)
		
		green.Printf("\n🎉 Exit completed via %s (Board approved)\n", strings.ToUpper(selectedExit.Type))
		fmt.Printf("   Company Valuation: $%s\n", formatFounderCurrency(selectedExit.Valuation))
		fmt.Printf("   Your Payout: $%s (%.2f%% equity)\n\n", 
			formatFounderCurrency(selectedExit.FounderPayout), founderEquity)
		
		fmt.Println("Press 'Enter' to see final results...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		
		return false // Game ends
	}

	// Confirm exit (only if founder has majority)
	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Printf("⚠️  CONFIRM %s\n", strings.ToUpper(selectedExit.Type))
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("\nCompany Valuation: $%s\n", formatFounderCurrency(selectedExit.Valuation))
	fmt.Printf("Your Payout: $%s\n", formatFounderCurrency(selectedExit.FounderPayout))

	if selectedExit.Type == "ipo" {
		fmt.Println("\n📊 IPO: You'll remain as CEO of a public company.")
		fmt.Println("   You can sell 20%% of your shares now for liquidity.")
	} else if selectedExit.Type == "acquisition" {
		fmt.Println("\n🤝 Acquisition: You'll likely have a 1-2 year earn-out,")
		fmt.Println("   then you're free to start your next venture.")
	} else if selectedExit.Type == "secondary" {
		fmt.Println("\n💼 Secondary Sale: Selling half your stake for liquidity,")
		fmt.Println("   but you stay on as CEO to build further.")
	}

	red.Println("\n⚠️  THIS IS IRREVERSIBLE. THE GAME WILL END.")
	fmt.Print("\nType 'EXIT' to confirm: ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToUpper(confirm))

	if confirm != "EXIT" {
		fmt.Println("\n✓ Continuing to build...")
		return true
	}

	// Execute exit
	fs.ExecuteExit(selectedExit.Type)

	green.Printf("\n🎉 Congratulations! You've exited via %s!\n", strings.ToUpper(selectedExit.Type))
	fmt.Printf("   Company Valuation: $%s\n", formatFounderCurrency(selectedExit.Valuation))
	fmt.Printf("   Your Payout: $%s\n\n", formatFounderCurrency(selectedExit.FounderPayout))
	
	fmt.Println("Press 'Enter' to see final results...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	return false // Game ends, process final month
}

func saveFounderScoreAndCheckAchievements(fs *founder.FounderState) {
	// Calculate final metrics
	_, valuation, founderEquity := fs.GetFinalScore()
	
	// Calculate founder payout (net worth)
	var founderPayout int64
	if fs.HasExited {
		switch fs.ExitType {
		case "ipo":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.20 / 100.0) // 20% immediate
		case "acquisition":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity / 100.0)
		case "secondary":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.5 / 100.0)
		default:
			founderPayout = fs.Cash + int64(float64(valuation)*founderEquity/100.0)
		}
	} else {
		founderPayout = fs.Cash + int64(float64(valuation)*founderEquity/100.0)
	}
	
	// Calculate initial investment (starting cash)
	// Use the starting cash amount from the game state
	initialCash := int64(500000) // Default estimate - founder mode starts with varying amounts
	if fs.Turn > 0 {
		// Estimate based on first funding round or use a reasonable default
		if len(fs.FundingRounds) > 0 {
			// Starting cash is roughly 2x first funding round
			initialCash = fs.FundingRounds[0].Amount * 2
		} else {
			// No funding rounds, use template default
			initialCash = 500000
		}
	}
	
	// Calculate ROI
	roi := 0.0
	if initialCash > 0 {
		roi = (float64(founderPayout-initialCash) / float64(initialCash)) * 100.0
	}
	
	// Calculate total funding raised
	totalFundingRaised := int64(0)
	for _, round := range fs.FundingRounds {
		totalFundingRaised += round.Amount
	}
	
	// Calculate successful exits
	successfulExits := 0
	if fs.HasExited {
		successfulExits = 1
	}
	
	// Save score to database (using GameScore table - founder mode counts as "Founder" difficulty)
	score := db.GameScore{
		PlayerName:      fs.FounderName,
		FinalNetWorth:   founderPayout,
		ROI:             roi,
		SuccessfulExits: successfulExits,
		TurnsPlayed:     fs.Turn,
		Difficulty:      "Founder",
		PlayedAt:        time.Now(),
	}
	
	err := db.SaveGameScore(score)
	if err != nil {
		color.Yellow("\nWarning: Could not save score: %v", err)
	} else {
		color.Green("\n%s Score saved to local leaderboard!", ascii.Check)
	}
	
	// Check for achievements
	previouslyUnlocked, err := db.GetPlayerAchievements(fs.FounderName)
	if err != nil {
		previouslyUnlocked = []string{}
	}
	
	// Get player stats
	playerStats, _ := db.GetPlayerStats(fs.FounderName)
	winStreak, _ := db.GetWinStreak(fs.FounderName)
	
	// Build game stats for achievement checking
	gameStats := achievements.GameStats{
		GameMode:             "founder",
		FinalNetWorth:        founderPayout,
		ROI:                  roi,
		SuccessfulExits:      successfulExits,
		TurnsPlayed:          fs.Turn,
		Difficulty:           "Founder",
		FinalMRR:             fs.MRR,
		FinalValuation:       valuation,
		FinalEquity:          founderEquity,
		Customers:            fs.Customers,
		FundingRoundsRaised:  len(fs.FundingRounds),
		TotalFundingRaised:   totalFundingRaised,
		HasExited:            fs.HasExited,
		ExitType:             fs.ExitType,
		ExitValuation:        fs.ExitValuation,
		MonthsToProfitability: fs.MonthReachedProfitability,
		RanOutOfCash:         fs.Cash <= 0 && !fs.HasExited, // Ran out of cash and didn't exit = lost
		TotalGames:           playerStats.TotalGames,
		TotalWins:            int(playerStats.WinRate * float64(playerStats.TotalGames) / 100.0),
		WinStreak:            winStreak,
		BestNetWorth:         playerStats.BestNetWorth,
		TotalExits:           playerStats.TotalExits,
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
			db.UnlockAchievement(fs.FounderName, ach.ID)
			
			// Display with animation
			achievementText := fmt.Sprintf("%s %s [%s]\n+%d points", ach.Icon, ach.Name, ach.Rarity, ach.Points)
			animations.ShowAchievementUnlock(achievementText, ach.Description)
		}
	} else {
		yellow.Println("\nNo new achievements unlocked this game.")
		yellow.Println("Keep playing to unlock more achievements!")
	}
	
	// Calculate and award XP
	xpEarned := progression.CalculateXPReward(&gameStats, newAchievementIDs)
	
	// Get profile before adding XP to compare levels
	profileBefore, _ := db.GetPlayerProfile(fs.FounderName)
	oldLevel := profileBefore.Level
	
	// Add XP to player profile
	leveledUp, newLevel, err := db.AddExperience(fs.FounderName, xpEarned)
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
			xpBreakdown["Successful Exit"] = progression.XPSuccessfulExit * gameStats.SuccessfulExits
		}
		
		// Founder mode specific bonuses
		if gameStats.HasExited && gameStats.ExitType == "ipo" {
			xpBreakdown["IPO Exit"] = 500
		} else if gameStats.HasExited && gameStats.ExitType == "acquisition" {
			xpBreakdown["Acquisition Exit"] = 300
		}
		
		if gameStats.MonthsToProfitability > 0 && gameStats.MonthsToProfitability <= 24 {
			xpBreakdown["Reached Profitability"] = 100
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
			DisplayLevelUp(fs.FounderName, oldLevel, newLevel, levelInfo.Unlocks)
		} else {
			// Show progress towards next level
			profileAfter, _ := db.GetPlayerProfile(fs.FounderName)
			fmt.Println()
			yellow.Printf("   Level %d Progress: ", profileAfter.Level)
			progressBar := progression.FormatXPBar(profileAfter.ExperiencePoints, profileAfter.NextLevelXP, 20)
			fmt.Printf("%s %d/%d XP\n", progressBar, profileAfter.ExperiencePoints, profileAfter.NextLevelXP)
		}
	}
	
	// Calculate and display career level and points (for old achievement system compatibility)
	totalLifetimePoints := 0
	allUnlocked, _ := db.GetPlayerAchievements(fs.FounderName)
	for _, id := range allUnlocked {
		if ach, exists := achievements.AllAchievements[id]; exists {
			totalLifetimePoints += ach.Points
		}
	}

	// Get owned upgrades to calculate available balance
	ownedUpgrades, _ := db.GetPlayerUpgrades(fs.FounderName)
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

func formatFounderNumber(amount int64) string {
	if amount < 0 {
		return fmt.Sprintf("-%s", formatFounderNumber(-amount))
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

func handleViewCustomerDeals(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("📋 CUSTOMER DEALS")
	fmt.Println(strings.Repeat("─", 70))

	activeCustomers := fs.GetActiveCustomers()
	churnedCustomers := fs.GetChurnedCustomers()

	if len(activeCustomers) == 0 && len(churnedCustomers) == 0 {
		cyan.Println("\nNo customers yet.")
		return
	}

	if len(activeCustomers) > 0 {
		green.Printf("\n✓ ACTIVE CUSTOMERS (%d):\n\n", len(activeCustomers))
		fmt.Printf("%-8s %-12s %-15s %-20s %-10s\n", "ID", "Source", "Deal Size/mo", "Term", "Added")
		fmt.Println(strings.Repeat("─", 70))
		for _, c := range activeCustomers {
			termStr := "Perpetual"
			if c.TermMonths > 0 {
				termStr = fmt.Sprintf("%d months", c.TermMonths)
			}
			fmt.Printf("%-8d %-12s $%-14s %-20s Month %-3d\n",
				c.ID, c.Source, formatFounderCurrency(c.DealSize), termStr, c.MonthAdded)
		}
	}

	if len(churnedCustomers) > 0 {
		red.Printf("\n\n✗ CHURNED CUSTOMERS (%d):\n\n", len(churnedCustomers))
		fmt.Printf("%-8s %-12s %-15s %-20s %-10s %-10s\n", "ID", "Source", "Deal Size/mo", "Term", "Added", "Churned")
		fmt.Println(strings.Repeat("─", 70))
		for _, c := range churnedCustomers {
			termStr := "Perpetual"
			if c.TermMonths > 0 {
				termStr = fmt.Sprintf("%d months", c.TermMonths)
			}
			fmt.Printf("%-8d %-12s $%-14s %-20s Month %-3d Month %-3d\n",
				c.ID, c.Source, formatFounderCurrency(c.DealSize), termStr, c.MonthAdded, c.MonthChurned)
		}
	}

	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleSolicitFeedback(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💬 SOLICIT CUSTOMER FEEDBACK")
	fmt.Println(strings.Repeat("─", 70))

	if fs.Customers == 0 {
		cyan.Println("\nYou need customers before you can solicit feedback!")
		fmt.Println("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	currentMaturity := fs.ProductMaturity
	currentChurn := fs.CustomerChurnRate
	err := fs.SolicitCustomerFeedback()
	if err != nil {
		color.Red("\nError: %s", err)
		fmt.Println("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	improvement := fs.ProductMaturity - currentMaturity
	churnReduction := currentChurn - fs.CustomerChurnRate
	green.Printf("\n✓ Customer feedback collected from %d customers!\n", fs.Customers)
	fmt.Printf("  Product maturity improved by %.1f%%\n", improvement*100)
	fmt.Printf("  New product maturity: %.0f%%\n", fs.ProductMaturity*100)
	green.Printf("  Churn rate reduced by %.1f%% (now %.1f%%)\n", churnReduction*100, fs.CustomerChurnRate*100)
	fmt.Println("\n  Note: This action takes effect this month.")
	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return false
}

func handleViewAffiliateProgram(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🤝 AFFILIATE PROGRAM")
	fmt.Println(strings.Repeat("─", 70))

	if fs.AffiliateProgram == nil {
		cyan.Println("\nNo affiliate program launched yet!")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	prog := fs.AffiliateProgram

	fmt.Printf("\n📊 Program Stats\n")
	fmt.Printf("Launched: Month %d (%d months ago)\n", prog.LaunchedMonth, fs.Turn-prog.LaunchedMonth)
	fmt.Printf("Active Affiliates: %d\n", prog.Affiliates)
	fmt.Printf("Commission Rate: %.0f%%\n", prog.Commission*100)
	fmt.Printf("\n💰 Costs\n")
	fmt.Printf("Setup Cost: $%s (one-time)\n", formatFounderCurrency(prog.SetupCost))
	fmt.Printf("Monthly Platform Fee: $%s\n", formatFounderCurrency(prog.MonthlyPlatformFee))
	fmt.Printf("\n📈 Performance\n")
	fmt.Printf("Total Customers Acquired: %d\n", prog.CustomersAcquired)
	fmt.Printf("Total Revenue Generated: $%s\n", formatFounderCurrency(prog.MonthlyRevenue))

	if prog.CustomersAcquired > 0 {
		avgPerAffiliate := float64(prog.CustomersAcquired) / float64(prog.Affiliates)
		green.Printf("\nAvg Customers per Affiliate: %.1f\n", avgPerAffiliate)
	}

	// Show active affiliate customers
	affiliateCustomers := 0
	var totalAffiliateMRR int64
	for _, c := range fs.CustomerList {
		if c.IsActive && c.Source == "affiliate" {
			affiliateCustomers++
			totalAffiliateMRR += c.DealSize
		}
	}

	if affiliateCustomers > 0 {
		fmt.Printf("\n📊 Current Active Affiliate Customers\n")
		fmt.Printf("Active: %d customers\n", affiliateCustomers)
		fmt.Printf("Monthly Recurring Revenue: $%s\n", formatFounderCurrency(totalAffiliateMRR))
		monthlyCommission := int64(float64(totalAffiliateMRR) * prog.Commission)
		fmt.Printf("Monthly Commission Paid: $%s\n", formatFounderCurrency(monthlyCommission))
	}

	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleEndAffiliateProgram(fs *founder.FounderState) bool {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)
	reader := bufio.NewReader(os.Stdin)

	if fs.AffiliateProgram == nil {
		cyan.Println("\nNo affiliate program is running!")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	prog := fs.AffiliateProgram
	
	// Count active affiliate customers
	affiliateCustomers := 0
	var totalAffiliateMRR int64
	for _, c := range fs.CustomerList {
		if c.IsActive && c.Source == "affiliate" {
			affiliateCustomers++
			totalAffiliateMRR += c.DealSize
		}
	}

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("🛑 END AFFILIATE PROGRAM")
	fmt.Println(strings.Repeat("─", 70))

	fmt.Printf("\nCurrent Program Stats:\n")
	fmt.Printf("Active Affiliates: %d\n", prog.Affiliates)
	fmt.Printf("Active Affiliate Customers: %d\n", affiliateCustomers)
	fmt.Printf("Affiliate MRR: $%s/month\n", formatFounderCurrency(totalAffiliateMRR))
	fmt.Printf("Monthly Platform Fee: $%s\n", formatFounderCurrency(prog.MonthlyPlatformFee))

	fmt.Println("\nWhat would you like to do with affiliate customers?")
	fmt.Println("1. Transition to Direct Sales (customers stay, no churn)")
	fmt.Println("2. Let Customers Churn (customers leave when program ends)")
	fmt.Println("0. Cancel")

	fmt.Print("\nYour choice (0-2): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "0" {
		fmt.Println("\n✓ Keeping affiliate program running...")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	transitionCustomers := false
	if choice == "1" {
		transitionCustomers = true
	} else if choice != "2" {
		red.Println("\nInvalid choice!")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	// End the program
	err := fs.EndAffiliateProgram(transitionCustomers)
	if err != nil {
		red.Printf("\n✗ Error: %v\n", err)
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return true
	}

	if transitionCustomers {
		green.Printf("\n✓ Affiliate program ended. %d customers transitioned to direct sales.\n", affiliateCustomers)
		green.Println("   No customer churn - they're now direct customers.")
	} else {
		yellow.Printf("\n✓ Affiliate program ended. %d affiliate customers churned.\n", affiliateCustomers)
		red.Printf("   Lost $%s/month in MRR.\n", formatFounderCurrency(totalAffiliateMRR))
	}

	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return true // Return to menu (doesn't advance turn)
}

func handleViewFinancials(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💰 FINANCIALS & CASH FLOW")
	fmt.Println(strings.Repeat("─", 70))

	// Revenue
	fmt.Println("\n📈 REVENUE")
	fmt.Printf("MRR: $%s/month\n", formatFounderCurrency(fs.MRR))
	if fs.AffiliateProgram != nil && fs.AffiliateMRR > 0 {
		fmt.Printf("  ├─ Direct MRR: $%s\n", formatFounderCurrency(fs.DirectMRR))
		fmt.Printf("  └─ Affiliate MRR: $%s\n", formatFounderCurrency(fs.AffiliateMRR))
	}

	// After 33% deductions (taxes, fees, overhead, savings)
	netRevenue := int64(float64(fs.MRR) * 0.67)
	fmt.Printf("Net Revenue (after 33%% deductions): $%s/month\n", formatFounderCurrency(netRevenue))
	fmt.Printf("  Deductions: Taxes (20%%), Processing (3%%), Overhead (5%%), Savings (5%%)\n")

	// Expenses
	fmt.Println("\n💸 EXPENSES")
	totalExpenses := int64(0)

	fmt.Printf("Team Salaries: $%s/month\n", formatFounderCurrency(fs.MonthlyTeamCost))
	totalExpenses += fs.MonthlyTeamCost

	if fs.MonthlyComputeCost > 0 || fs.MonthlyODCCost > 0 {
		infraCosts := fs.MonthlyComputeCost + fs.MonthlyODCCost
		fmt.Printf("Infrastructure (Compute + ODC): $%s/month\n", formatFounderCurrency(infraCosts))
		totalExpenses += infraCosts
	}

	if fs.AffiliateProgram != nil {
		fmt.Printf("Affiliate Program Platform Fee: $%s/month\n", formatFounderCurrency(fs.AffiliateProgram.MonthlyPlatformFee))
		totalExpenses += fs.AffiliateProgram.MonthlyPlatformFee

		// Commission is variable, estimate based on current affiliate MRR
		commission := int64(float64(fs.AffiliateMRR) * fs.AffiliateProgram.Commission)
		if commission > 0 {
			fmt.Printf("Affiliate Commissions (~%.0f%%): $%s/month\n",
				fs.AffiliateProgram.Commission*100, formatFounderCurrency(commission))
			totalExpenses += commission
		}
	}

	// Global market costs
	if len(fs.GlobalMarkets) > 0 {
		var marketCosts int64
		for _, m := range fs.GlobalMarkets {
			marketCosts += m.MonthlyCost
		}
		fmt.Printf("Global Markets (%d): $%s/month\n", len(fs.GlobalMarkets), formatFounderCurrency(marketCosts))
		totalExpenses += marketCosts
	}

	fmt.Printf("\nTotal Monthly Expenses: $%s\n", formatFounderCurrency(totalExpenses))

	// Net Income
	netIncome := netRevenue - totalExpenses
	fmt.Println("\n📊 NET INCOME")
	if netIncome > 0 {
		green.Printf("✅ Profitable: +$%s/month\n", formatFounderCurrency(netIncome))
	} else {
		red.Printf("⚠️  Burning: -$%s/month\n", formatFounderCurrency(-netIncome))
	}

	// Margins
	if fs.MRR > 0 {
		grossMargin := float64(netRevenue) / float64(fs.MRR) * 100
		netMargin := float64(netIncome) / float64(fs.MRR) * 100
		fmt.Printf("\nGross Margin: %.1f%%\n", grossMargin)
		if netMargin > 0 {
			green.Printf("Net Margin: %.1f%%\n", netMargin)
		} else {
			red.Printf("Net Margin: %.1f%%\n", netMargin)
		}
	}

	// Runway
	fmt.Println("\n⏱️  RUNWAY")
	fmt.Printf("Current Cash: $%s\n", formatFounderCurrency(fs.Cash))
	if netIncome > 0 {
		green.Println("Runway: Infinite (profitable)")
		fmt.Printf("Monthly Cash Growth: +$%s\n", formatFounderCurrency(netIncome))
	} else if netIncome < 0 {
		burnRate := -netIncome
		runway := int(float64(fs.Cash) / float64(burnRate))
		if runway <= 3 {
			red.Printf("Runway: %d months ⚠️  CRITICAL\n", runway)
		} else if runway <= 6 {
			yellow.Printf("Runway: %d months\n", runway)
		} else {
			fmt.Printf("Runway: %d months\n", runway)
		}
		red.Printf("Monthly Burn: -$%s\n", formatFounderCurrency(burnRate))
	} else {
		fmt.Println("Runway: Break-even")
	}

	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleViewTeamRoster(fs *founder.FounderState) {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("👥 TEAM ROSTER")
	fmt.Println(strings.Repeat("─", 70))

	allEmployees := []founder.Employee{}

	// Collect all employees
	for _, e := range fs.Team.Engineers {
		allEmployees = append(allEmployees, e)
	}
	for _, e := range fs.Team.Sales {
		allEmployees = append(allEmployees, e)
	}
	for _, e := range fs.Team.CustomerSuccess {
		allEmployees = append(allEmployees, e)
	}
	for _, e := range fs.Team.Marketing {
		allEmployees = append(allEmployees, e)
	}
	for _, e := range fs.Team.Executives {
		allEmployees = append(allEmployees, e)
	}

	if len(allEmployees) == 0 {
		cyan.Println("\nNo employees yet!")
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// Display each employee
	for _, e := range allEmployees {
		roleTitle := strings.ToUpper(string(e.Role))
		if e.IsExecutive {
			green.Printf("\n━━━ %s ━━━\n", roleTitle)
		} else {
			fmt.Printf("\n── %s ──\n", roleTitle)
		}

		fmt.Printf("Name: %s\n", e.Name)
		fmt.Printf("Salary: $%s/year ($%s/month)\n",
			formatFounderCurrency(e.MonthlyCost*12), formatFounderCurrency(e.MonthlyCost))
		
		// Show market assignment for non-executives
		if !e.IsExecutive && e.AssignedMarket != "" {
			if e.AssignedMarket == "All" {
				cyan.Printf("Assigned Market: %s (works globally)\n", e.AssignedMarket)
			} else {
				fmt.Printf("Assigned Market: %s\n", e.AssignedMarket)
			}
		}

		if e.Equity > 0 {
			vestedEquity := 0.0
			if e.HasCliff {
				// After cliff, equity vests monthly
				vestedEquity = e.Equity * (float64(e.VestedMonths) / float64(e.VestingMonths))
			}

			fmt.Printf("Equity: %.2f%% ", e.Equity)
			if vestedEquity > 0 {
				green.Printf("(%.2f%% vested)\n", vestedEquity)
			} else {
				fmt.Printf("(unvested)\n")
			}

			if !e.HasCliff {
				fmt.Printf("Vesting: %d/%d months until cliff (%d months remaining)\n",
					e.VestedMonths, e.CliffMonths, e.CliffMonths-e.VestedMonths)
			} else {
				fmt.Printf("Vesting: %d/%d months (%.0f%% vested)\n",
					e.VestedMonths, e.VestingMonths, (float64(e.VestedMonths)/float64(e.VestingMonths))*100)
			}
		}

		fmt.Printf("Hired: Month %d", e.MonthHired)
		if fs.Turn > e.MonthHired {
			fmt.Printf(" (%d months ago)\n", fs.Turn-e.MonthHired)
		} else {
			fmt.Println(" (this month)")
		}

		fmt.Printf("Performance: %.1fx impact\n", e.Impact)
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Printf("Total Team: %d employees\n", len(allEmployees))
	fmt.Printf("Monthly Payroll: $%s\n", formatFounderCurrency(fs.MonthlyTeamCost))
	fmt.Printf("Total Employee Equity: %.1f%% (%.1f%% available in pool)\n",
		fs.EquityAllocated, fs.EquityPool-fs.EquityAllocated)

	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func handleStrategicOpportunity(fs *founder.FounderState) bool {
	if fs.PendingOpportunity == nil {
		fmt.Println("\nNo pending opportunity")
		return true
	}

	opp := fs.PendingOpportunity
	reader := bufio.NewReader(os.Stdin)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	fmt.Println("\n" + strings.Repeat("─", 70))
	yellow.Println("💡 STRATEGIC OPPORTUNITY DECISION")
	fmt.Println(strings.Repeat("─", 70))

	fmt.Printf("\n%s\n\n", opp.Title)
	fmt.Printf("%s\n\n", opp.Description)
	green.Printf("✓ Benefits: %s\n", opp.Benefit)
	red.Printf("⚠️  Risks: %s\n", opp.Risk)
	if opp.Cost > 0 {
		fmt.Printf("\n💰 Cost: $%s", formatFounderCurrency(opp.Cost))
		fmt.Printf(" (Current cash: $%s)\n", formatFounderCurrency(fs.Cash))
	}

	// Check if chairman exists and can delegate
	chairman := fs.GetChairman()
	hasChairman := chairman != nil && chairman.IsActive && chairman.IsChairman
	
	fmt.Println("\nOptions:")
	fmt.Println("1. Accept Opportunity")
	fmt.Println("2. Decline Opportunity")
	if hasChairman && (opp.Type == "conference" || opp.Type == "press") {
		green.Println("3. Delegate to Chairman (saves founder time, chairman handles it)")
	}
	fmt.Println("0. Decide Later")

	maxOption := 2
	if hasChairman && (opp.Type == "conference" || opp.Type == "press") {
		maxOption = 3
	}
	fmt.Printf("\nYour decision (0-%d): ", maxOption)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		// Check if can afford
		if opp.Cost > fs.Cash {
			color.Red("\n✗ Insufficient cash! Need $%s, have $%s",
				formatFounderCurrency(opp.Cost), formatFounderCurrency(fs.Cash))
			fmt.Println("\nPress 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return true
		}

		// Accept opportunity - apply effects based on type
		fs.Cash -= opp.Cost

		switch opp.Type {
		case "press":
			// Add customers over next few months
			newCustomers := 5 + rand.Intn(10)
			green.Printf("\n✓ TechCrunch feature published! Expect %d new customers over next 3 months\n", newCustomers)
		case "enterprise_pilot":
			if rand.Float64() < 0.8 {
				green.Println("\n✓ Enterprise pilot successful! Major customer acquired")
			} else {
				yellow.Println("\n~ Enterprise pilot didn't convert, but gained valuable learnings")
			}
		case "bridge_round":
			green.Println("\n✓ Bridge round closed! Cash runway extended")
		case "conference":
			green.Println("\n✓ Conference presentation successful! Leads and recruiting pipeline boosted")
		case "talent":
			green.Println("\n✓ Star engineer hired! Product development accelerated")
		case "competitor_distress":
			green.Println("\n✓ Competitor customers acquired! Market position strengthened")
		}

		fs.PendingOpportunity = nil
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return false

	case "3":
		// Delegate to chairman (only for conference/press)
		if !hasChairman || (opp.Type != "conference" && opp.Type != "press") {
			color.Red("\nInvalid choice")
			return true
		}
		
		// Check if can afford
		if opp.Cost > fs.Cash {
			color.Red("\n✗ Insufficient cash! Need $%s, have $%s",
				formatFounderCurrency(opp.Cost), formatFounderCurrency(fs.Cash))
			fmt.Println("\nPress 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return true
		}
		
		// Chairman handles it - same benefits but founder saves time
		fs.Cash -= opp.Cost
		green.Printf("\n✓ Delegated to %s (Chairman). They'll handle this on your behalf.\n", chairman.Name)
		green.Println("   You save founder time and can focus on other priorities.")
		
		// Apply same effects as accepting, but with chairman bonus
		switch opp.Type {
		case "press":
			newCustomers := 5 + rand.Intn(10)
			green.Printf("✓ TechCrunch feature published! Expect %d new customers over next 3 months\n", newCustomers)
		case "conference":
			green.Println("✓ Conference presentation successful! Leads and recruiting pipeline boosted")
			green.Println("   Chairman's network connections enhanced the results.")
		}
		
		fs.PendingOpportunity = nil
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return false

	case "2":
		yellow.Printf("\n✓ Declined: %s\n", opp.Title)
		fs.PendingOpportunity = nil
		fmt.Println("\nPress 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return false

	case "0":
		fmt.Println("\n✓ Will decide later (opportunity still pending)")
		return true

	default:
		color.Red("\nInvalid choice")
		return true
	}
}
