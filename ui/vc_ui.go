package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/animations"
	"github.com/jamesacampbell/unicorn/ascii"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/game"
	"github.com/jamesacampbell/unicorn/upgrades"
)

func FormatCurrency(amount int64) string {
	if amount < 0 {
		return fmt.Sprintf("-$%s", FormatCurrency(-amount))
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

func AskForAutomatedMode() bool {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    GAME MODE SELECTION")
	cyan.Println(strings.Repeat("=", 70))

	fmt.Println()
	green.Println("1. Manual Mode (Press Enter each turn) [RECOMMENDED]")
	fmt.Println("   → Full access to all features:")
	magenta.Println("     • Operational value-add actions")
	magenta.Println("     • Due diligence before investing")
	magenta.Println("     • Secondary market stake sales")
	magenta.Println("     • Active founder relationship management")
	fmt.Println("   → Strategic decision-making between turns")

	fmt.Println()
	yellow.Println("2. Automated Mode (1 second per turn)")
	fmt.Println("   → Simplified gameplay - no interactive features")
	fmt.Println("   → Reputation affects deal quality only")
	fmt.Println("   → Best for quick games or testing")

	fmt.Print("\nEnter your choice (1-2, default 1): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Show spinner after input
	spinner, _ := animations.StartSpinner("Processing...")
	time.Sleep(500 * time.Millisecond)
	spinner.Stop()

	return choice == "2"
}

func formatStrategyName(strategy string) string {
	switch strategy {
	case "conservative":
		return "Conservative"
	case "aggressive":
		return "Aggressive"
	case "balanced":
		return "Balanced"
	case "deep_tech":
		return "Deep Tech"
	case "consumer_focused":
		return "Consumer Focused"
	case "enterprise_focused":
		return "Enterprise Focused"
	case "early_stage":
		return "Early Stage"
	case "growth_stage":
		return "Growth Stage"
	case "mega_fund":
		return "Mega Fund"
	case "seed_focused":
		return "Seed Focused"
	default:
		// Capitalize first letter and replace underscores with spaces
		if len(strategy) == 0 {
			return strategy
		}
		result := strings.ToUpper(string(strategy[0]))
		for i := 1; i < len(strategy); i++ {
			if strategy[i] == '_' {
				result += " "
				if i+1 < len(strategy) {
					result += strings.ToUpper(string(strategy[i+1]))
					i++ // Skip the next character since we already capitalized it
				}
			} else {
				result += string(strategy[i])
			}
		}
		return result
	}
}

func DisplayWelcome(username string, difficulty game.Difficulty, playerUpgrades []string, aiPlayers []game.AIPlayer) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta, color.Bold)
	green := color.New(color.FgGreen)

	cyan.Printf("\n%s, welcome to your investment journey!\n", username)
	fmt.Printf("\nDifficulty: ")
	yellow.Printf("%s\n", difficulty.Name)
	fmt.Printf("Fund Size: $%s\n", FormatMoney(difficulty.StartingCash))
	lpCommittedCapital := difficulty.StartingCash * 2
	fmt.Printf("LP Committed Capital: $%s (capital calls quarterly)\n", FormatMoney(lpCommittedCapital))
	fmt.Printf("Follow-on Reserve: $%s ($100k base + $50k per round)\n", FormatMoney(int64(1000000)))
	fmt.Printf("Management Fee: 2%% annually ($%s/year)\n", FormatMoney(int64(float64(difficulty.StartingCash)*0.02)))
	fmt.Printf("Game Duration: %d turns (%d years)\n", difficulty.MaxTurns, difficulty.MaxTurns/12)

	// Display active upgrades (filtered for VC mode)
	vcUpgrades := upgrades.FilterUpgradeIDsForGameMode(playerUpgrades, "vc")
	if len(vcUpgrades) > 0 {
		fmt.Println()
		green.Println("✨ ACTIVE UPGRADES FOR THIS GAME:")
		for _, upgradeID := range vcUpgrades {
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
	for _, aiPlayer := range aiPlayers {
		strategyName := formatStrategyName(aiPlayer.Strategy)
		fmt.Printf("   ? %s (%s) - %s\n", aiPlayer.Name, aiPlayer.Firm, strategyName)
	}

	fmt.Print("\nPress 'Enter' to see available startups...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func DisplayStartup(s game.Startup, index int, availableCash int64, playerUpgrades []string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)

	cyan.Printf("\n[%d] %s\n", index+1, s.Name)
	fmt.Printf("    %s\n", s.Description)
	yellow.Printf("    Category: %s\n", s.Category)
	fmt.Printf("    Valuation: $%s\n", FormatMoney(s.Valuation))
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
	green.Printf("    Max Investment: $%s", FormatMoney(maxAvailable))
	if maxAvailable < maxInvestment {
		fmt.Printf(" (limited by available cash, max would be $%s)", FormatMoney(maxInvestment))
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
	fmt.Printf("    Website Visitors: %s/month\n", FormatNumber(s.MonthlyWebsiteVisitors))

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

func HandleFollowOnOpportunities(gs *game.GameState, opportunities []game.FollowOnOpportunity) {
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
		fmt.Printf("   Pre-money Valuation: $%s\n", FormatMoney(opp.PreMoneyVal))
		fmt.Printf("   Post-money Valuation: $%s\n", FormatMoney(opp.PostMoneyVal))
		yellow.Printf("   Your Current Equity: %.2f%%\n", opp.CurrentEquity)
		availableFunds := gs.Portfolio.Cash + gs.Portfolio.FollowOnReserve
		green.Printf("   Available Funds: $%s (Cash: $%s + Reserve: $%s)\n",
			FormatMoney(availableFunds),
			FormatMoney(gs.Portfolio.Cash),
			FormatMoney(gs.Portfolio.FollowOnReserve))

		fmt.Println("\n" + strings.Repeat("-", 70))
		cyan.Println("\n?? INVEST MORE TO AVOID DILUTION!")
		fmt.Println("   If you don't invest, your ownership % will decrease.")
		fmt.Println("   If you DO invest, you'll maintain or increase your stake.")
		fmt.Printf("\n   Investment Range: $%s to $%s\n",
			FormatMoney(opp.MinInvestment), FormatMoney(opp.MaxInvestment))

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
			color.Red("Amount below minimum investment of $%s", FormatMoney(opp.MinInvestment))
			continue
		}

		if amount > opp.MaxInvestment {
			color.Red("Amount exceeds maximum investment of $%s", FormatMoney(opp.MaxInvestment))
			continue
		}

		err = gs.MakeFollowOnInvestment(opp.CompanyName, amount)
		if err != nil {
			color.Red("Error: %v", err)
		} else {
			green.Printf("\n%s Follow-on investment successful! Invested $%s in %s\n",
				ascii.Check, FormatMoney(amount), opp.CompanyName)
			fmt.Printf("Follow-on Reserve Remaining: $%s\n", FormatMoney(gs.Portfolio.FollowOnReserve))
			fmt.Printf("Cash Remaining: $%s\n", FormatMoney(gs.Portfolio.Cash))
		}
	}

	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayBoardMeeting shows a visual board meeting interface
func DisplayBoardMeeting(gs *game.GameState, companyName string, voteTitle string) {
	clear.ClearIt()

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	green := color.New(color.FgGreen)

	// Get board members
	members := gs.GetBoardMembers(companyName)

	// Display board table header
	fmt.Println("\n" + strings.Repeat("=", 70))
	magenta.Println("                    🏛️  BOARD MEETING 🏛️")
	fmt.Println(strings.Repeat("=", 70))

	cyan.Printf("\nCompany: %s\n", companyName)
	yellow.Printf("Agenda: %s\n\n", voteTitle)

	// Display board members
	fmt.Println(strings.Repeat("-", 70))
	green.Println("Board Members Present:")
	fmt.Println(strings.Repeat("-", 70))

	if len(members) == 0 {
		fmt.Println("  (No board members found)")
	} else {
		for i, member := range members {
			memberLabel := fmt.Sprintf("%d. %s", i+1, member.Name)
			if member.IsPlayer {
				magenta.Printf("  %s (%s) - %d vote(s)\n", memberLabel, member.Firm, member.VoteWeight)
			} else {
				fmt.Printf("  %s (%s) - %d vote(s)\n", memberLabel, member.Firm, member.VoteWeight)
			}
		}
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Println()

	// Simple board table visualization
	fmt.Println("                    ╔═══════════════════════════════════════╗")
	fmt.Println("                    ║                                       ║")
	fmt.Println("                    ║         BOARD OF DIRECTORS            ║")
	fmt.Println("                    ║                                       ║")
	fmt.Println("                    ╚═══════════════════════════════════════╝")
	fmt.Println()

	// Show members around the table (up to 4 positions)
	positions := []string{"", "", "", ""}
	for i, member := range members {
		if i < 4 {
			if member.IsPlayer {
				positions[i] = fmt.Sprintf("👤 %s", member.Name)
			} else {
				positions[i] = fmt.Sprintf("💼 %s", member.Name)
			}
		}
	}

	if len(positions) > 0 && positions[0] != "" {
		fmt.Printf("         %-25s │                           │", positions[0])
		if len(positions) > 1 && positions[1] != "" {
			fmt.Printf(" %s\n", positions[1])
		} else {
			fmt.Println()
		}
	} else {
		fmt.Println("                    │                           │")
	}
	fmt.Println("                    │                           │")
	if len(positions) > 2 && positions[2] != "" {
		fmt.Printf("         %-25s │                           │", positions[2])
		if len(positions) > 3 && positions[3] != "" {
			fmt.Printf(" %s\n", positions[3])
		} else {
			fmt.Println()
		}
	} else {
		fmt.Println("                    │                           │")
	}
	fmt.Println("                    └───────────────────────────┘")
	fmt.Println()

	time.Sleep(800 * time.Millisecond) // Brief pause for effect
}

func HandleBoardVotes(gs *game.GameState, votes []game.BoardVote) {
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

		// Display board meeting interface
		DisplayBoardMeeting(gs, vote.CompanyName, vote.Title)

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

func SelectInvestmentTerms(gs *game.GameState, startup *game.Startup, amount int64) game.InvestmentTerms {
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
		if opt.ValuationCap > 0 {
			green.Printf("   ✓ Valuation Cap: $%s (converts at cap if company raises above)\n", FormatMoney(opt.ValuationCap))
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

func PlayTurn(gs *game.GameState, autoMode bool) {
	yellow := color.New(color.FgYellow, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)

	// Check for follow-on investment opportunities BEFORE processing turn
	// This way the player can invest before dilution happens
	// Always pause for follow-on investments, even in automated mode
	opportunities := gs.GetFollowOnOpportunities()
	if len(opportunities) > 0 {
		HandleFollowOnOpportunities(gs, opportunities)
	}

	// Print separator line instead of clearing screen
	fmt.Println(strings.Repeat("=", 70))

	// Show round transition animation for milestones (every 5 turns)
	if gs.Portfolio.Turn%5 == 0 {
		animations.ShowRoundTransition(gs.Portfolio.Turn)
	}

	yellow.Printf("\n%s MONTH %d of %d\n", ascii.Calendar, gs.Portfolio.Turn, gs.Portfolio.MaxTurns)

	// Strategic Advisor: Show preview of next board vote
	playerUpgrades, _ := database.GetPlayerUpgrades(gs.PlayerName)
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
	
	// Process active value-add actions
	valueAddMsgs := gs.ProcessActiveValueAddActions()
	messages = append(messages, valueAddMsgs...)
	
	// Generate and process relationship events
	for i := range gs.Portfolio.Investments {
		inv := &gs.Portfolio.Investments[i]
		if inv.FounderName != "" {
			event := game.GenerateRelationshipEvent(inv, gs.Portfolio.Turn)
			if event != nil {
				inv.RelationshipScore = game.ApplyRelationshipChange(
					inv.RelationshipScore, 
					event.ScoreChange)
				messages = append(messages, event.Description)
				
				// Check for board removal due to poor relationship
				if inv.Terms.HasBoardSeat && game.CanBeFiredFromBoard(inv.RelationshipScore) {
					inv.Terms.HasBoardSeat = false
					messages = append(messages, fmt.Sprintf("⚠️  Poor relationship with %s led to board removal!", inv.FounderName))
				}
			}
		}
	}
	
	// Generate secondary market offers
	newOffers := gs.GenerateSecondaryOffers()
	gs.SecondaryMarketOffers = append(gs.SecondaryMarketOffers, newOffers...)
	
	// Process offer expirations
	expiredMsgs := gs.ProcessSecondaryOfferExpirations()
	messages = append(messages, expiredMsgs...)

	// Check for pending board votes AFTER processing turn
	// Board votes are created during ProcessTurn for acquisitions/down rounds
	pendingVotes := gs.GetPendingBoardVotes()
	if len(pendingVotes) > 0 {
		HandleBoardVotes(gs, pendingVotes)
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
				inv.CompanyName, FormatMoney(inv.AmountInvested), inv.EquityPercent, dilutionInfo)
			fmt.Printf("      Current Value: $%s ", FormatMoney(value))
			profitColor.Printf("(%s$%s)\n", profitSign, FormatMoney(abs(profit)))

			totalCompanyValuation += inv.CurrentValuation
		}

		// Display total company valuation
		yellow.Printf("\n   Total Company Valuation: $%s\n", FormatMoney(totalCompanyValuation))
	}

	fmt.Printf("\n%s Net Worth: $%s", ascii.Money, FormatMoney(gs.Portfolio.NetWorth))
	fmt.Printf(" | Cash: $%s | Follow-on Reserve: $%s\n",
		FormatMoney(gs.Portfolio.Cash), FormatMoney(gs.Portfolio.FollowOnReserve))
	fmt.Printf("   Management Fees Paid: $%s\n", FormatMoney(gs.Portfolio.ManagementFeesCharged))

	// Show competitive leaderboard every quarter
	if gs.Portfolio.Turn%3 == 0 {
		DisplayMiniLeaderboard(gs)
	}

	// Portfolio Dashboard option
	if !autoMode {
		fmt.Println()
		yellow.Println("Press 'd' for Portfolio Dashboard, or Enter to continue...")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "d" || input == "dashboard" {
			DisplayPortfolioDashboard(gs)
			// After dashboard, show portfolio again and ask to continue
			fmt.Println()
			yellow.Println("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			// Re-show portfolio summary
			cyan.Print(ascii.PortfolioHeader)
			if len(gs.Portfolio.Investments) > 0 {
				for _, inv := range gs.Portfolio.Investments {
					value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
					profit := value - inv.AmountInvested
					profitColor := color.New(color.FgGreen)
					profitSign := "+"
					if profit < 0 {
						profitColor = color.New(color.FgRed)
						profitSign = ""
					}
					fmt.Printf("   %s: $%s → $%s ", inv.CompanyName, FormatMoney(inv.AmountInvested), FormatMoney(value))
					profitColor.Printf("(%s$%s)\n", profitSign, FormatMoney(abs(profit)))
				}
			}
			fmt.Printf("\n%s Net Worth: $%s\n", ascii.Money, FormatMoney(gs.Portfolio.NetWorth))
		}
	}

	// Show interactive features (Manual Mode only)
	if !autoMode {
		// Value-add opportunities
		ShowValueAddMenu(gs, autoMode)
		
		// Secondary market offers
		ShowSecondaryMarketOffers(gs, autoMode)
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

func DisplayFinalScore(gs *game.GameState) {
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
	fmt.Printf("%s Management Fees Paid: $%s\n", ascii.Money, FormatMoney(gs.Portfolio.ManagementFeesCharged))
	if gs.Portfolio.CarryInterestPaid > 0 {
		yellow := color.New(color.FgYellow)
		yellow.Printf("%s Carry Interest Paid (20%%): $%s\n", ascii.Money, FormatMoney(gs.Portfolio.CarryInterestPaid))
	}
	fmt.Println()

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("%s Final Net Worth (after carry): $%s\n", ascii.Money, FormatMoney(netWorth))

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
		fmt.Printf("$%-14s ", FormatMoney(entry.NetWorth))
		roiColorEntry.Printf("%.1f%%\n", entry.ROI)
	}

	if leaderboard[0].IsPlayer {
		magenta.Println("\n?? CONGRATULATIONS! You beat all the AI investors!")
	} else {
		magenta.Printf("\nYou finished in position #%d. Better luck next time!\n", FindPlayerRank(leaderboard))
	}

	fmt.Println("\n" + strings.Repeat("?", 50))
	fmt.Println("FINAL PORTFOLIO:")
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		fmt.Printf("   %s: $%s ? $%s\n",
			inv.CompanyName, FormatMoney(inv.AmountInvested), FormatMoney(value))
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

func DisplayMiniLeaderboard(gs *game.GameState) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n?? Current Standings:")
	leaderboard := gs.GetLeaderboard()

	for i, entry := range leaderboard {
		marker := "  "
		if entry.IsPlayer {
			marker = "? "
			yellow.Printf("%s%d. %s (%s): $%s (ROI: %.1f%%)\n",
				marker, i+1, entry.Name, entry.Firm, FormatMoney(entry.NetWorth), entry.ROI)
		} else {
			fmt.Printf("%s%d. %s (%s): $%s (ROI: %.1f%%)\n",
				marker, i+1, entry.Name, entry.Firm, FormatMoney(entry.NetWorth), entry.ROI)
		}
	}
}

func FindPlayerRank(leaderboard []game.PlayerScore) int {
	for i, entry := range leaderboard {
		if entry.IsPlayer {
			return i + 1
		}
	}
	return len(leaderboard)
}

func SelectDifficulty(username string) game.Difficulty {
	clear.ClearIt()
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	magenta := color.New(color.FgMagenta)
	gray := color.New(color.FgHiBlack)

	// Get player level
	profile, err := database.GetPlayerProfile(username)
	playerLevel := 1
	if err == nil {
		playerLevel = profile.Level
	}

	cyan.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("                 SELECT DIFFICULTY")
	cyan.Println(strings.Repeat("=", 60))

	if playerLevel > 1 {
		fmt.Printf("\n   Your Level: %d\n", playerLevel)
	}

	// Easy - always available
	green.Printf("\n1. Easy")
	fmt.Printf(" - %s\n", game.EasyDifficulty.Description)
	fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
		FormatMoney(game.EasyDifficulty.StartingCash), game.EasyDifficulty.MaxTurns)

	// Medium - always available
	yellow.Printf("\n2. Medium")
	fmt.Printf(" - %s\n", game.MediumDifficulty.Description)
	fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
		FormatMoney(game.MediumDifficulty.StartingCash), game.MediumDifficulty.MaxTurns)

	// Hard - requires level 5
	if playerLevel >= 5 {
		red.Printf("\n3. Hard")
		fmt.Printf(" - %s\n", game.HardDifficulty.Description)
		fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
			FormatMoney(game.HardDifficulty.StartingCash), game.HardDifficulty.MaxTurns)
	} else {
		gray.Printf("\n3. Hard 🔒")
		gray.Printf(" - Unlocks at Level 5\n")
		fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
			FormatMoney(game.HardDifficulty.StartingCash), game.HardDifficulty.MaxTurns)
	}

	// Expert - requires level 10
	if playerLevel >= 10 {
		magenta.Printf("\n4. Expert")
		fmt.Printf(" - %s\n", game.ExpertDifficulty.Description)
		fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
			FormatMoney(game.ExpertDifficulty.StartingCash), game.ExpertDifficulty.MaxTurns)
	} else {
		gray.Printf("\n4. Expert 🔒")
		gray.Printf(" - Unlocks at Level 10\n")
		fmt.Printf("   Starting Cash: $%s | Max Turns: %d\n",
			FormatMoney(game.ExpertDifficulty.StartingCash), game.ExpertDifficulty.MaxTurns)
	}

	fmt.Print("\nEnter your choice (1-4): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Show spinner after input
	spinner, _ := animations.StartSpinner("Loading difficulty...")
	time.Sleep(500 * time.Millisecond)
	spinner.Stop()

	switch choice {
	case "1":
		return game.EasyDifficulty
	case "3":
		if playerLevel >= 5 {
			return game.HardDifficulty
		} else {
			color.Yellow("\n🔒 Hard difficulty unlocks at Level 5. Playing Medium instead.")
			time.Sleep(1500 * time.Millisecond)
			return game.MediumDifficulty
		}
	case "4":
		if playerLevel >= 10 {
			return game.ExpertDifficulty
		} else {
			color.Yellow("\n🔒 Expert difficulty unlocks at Level 10. Playing Medium instead.")
			time.Sleep(1500 * time.Millisecond)
			return game.MediumDifficulty
		}
	default:
		return game.MediumDifficulty
	}
}

func handleSyndicateInvestment(gs *game.GameState) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta, color.Bold)

	if len(gs.SyndicateOpportunities) == 0 {
		color.Yellow("\nNo syndicate opportunities available.")
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	magenta.Println("🤝 SYNDICATE INVESTMENT")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\nAvailable Syndicate Deals:")
	for i, opp := range gs.SyndicateOpportunities {
		startup := gs.AvailableStartups[opp.StartupIndex]
		fmt.Printf("\n%d. %s\n", i+1, opp.CompanyName)
		fmt.Printf("   Lead: %s (%s)\n", opp.LeadInvestor, opp.LeadInvestorFirm)
		fmt.Printf("   Round Size: $%s | Valuation: $%s\n", FormatMoney(opp.TotalRoundSize), FormatMoney(opp.Valuation))
		green.Printf("   Investment Range: $%s - $%s\n", FormatMoney(opp.YourMinShare), FormatMoney(opp.YourMaxShare))
		fmt.Printf("   Category: %s\n", startup.Category)
	}

	fmt.Printf("\nSelect syndicate deal (1-%d, or 0 to cancel): ", len(gs.SyndicateOpportunities))
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)

	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 0 || choice > len(gs.SyndicateOpportunities) {
		color.Yellow("Cancelled.")
		return
	}

	if choice == 0 {
		return
	}

	opp := gs.SyndicateOpportunities[choice-1]

	cyan.Printf("\n💵 CO-INVESTING IN: %s\n", opp.CompanyName)
	fmt.Printf("   Lead Investor: %s (%s)\n", opp.LeadInvestor, opp.LeadInvestorFirm)
	fmt.Printf("   Total Round: $%s\n", FormatMoney(opp.TotalRoundSize))
	fmt.Printf("   Company Valuation: $%s\n", FormatMoney(opp.Valuation))
	yellow.Printf("   Your Investment Range: $%s - $%s\n", FormatMoney(opp.YourMinShare), FormatMoney(opp.YourMaxShare))
	fmt.Printf("\n   Benefits:\n")
	for _, benefit := range opp.Benefits {
		fmt.Printf("     • %s\n", benefit)
	}

	fmt.Printf("\nEnter investment amount ($%s - $%s, or 0 to cancel): $",
		FormatMoney(opp.YourMinShare), FormatMoney(opp.YourMaxShare))
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)

	if amountStr == "" || amountStr == "0" {
		color.Yellow("Cancelled.")
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		color.Red("Invalid amount!")
		return
	}

	err = gs.MakeSyndicateInvestment(choice-1, amount)
	if err != nil {
		color.Red("Error: %v", err)
	} else {
		green.Printf("\n%s Syndicate investment successful!\n", ascii.Check)
		fmt.Printf("Invested $%s in %s via %s's syndicate\n", FormatMoney(amount), opp.CompanyName, opp.LeadInvestor)
		fmt.Printf("Cash remaining: $%s\n", FormatMoney(gs.Portfolio.Cash))

		// Calculate equity
		equityPercent := (float64(amount) / float64(opp.TotalRoundSize)) * 100.0 * 1.05 // 5% bonus
		if equityPercent > 20.0 {
			equityPercent = 20.0
		}
		fmt.Printf("Equity acquired: %.2f%%\n", equityPercent)
	}
}

func investmentPhase(gs *game.GameState, playerLevel int, autoMode bool) {
	clear.ClearIt()
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Print(ascii.InvestmentHeader)
	green.Printf("\nTurn %d/%d\n", gs.Portfolio.Turn, gs.Portfolio.MaxTurns)
	fmt.Printf("Fund Size: $%s\n", FormatMoney(gs.Portfolio.InitialFundSize))
	fmt.Printf("Cash Available: $%s\n", FormatMoney(gs.Portfolio.Cash))
	fmt.Printf("Follow-on Reserve: $%s\n", FormatMoney(gs.Portfolio.FollowOnReserve))
	fmt.Printf("Portfolio Value: $%s\n", FormatMoney(gs.GetPortfolioValue()))
	fmt.Printf("Net Worth: $%s\n", FormatMoney(gs.Portfolio.NetWorth))

	// Calculate and display total company valuation
	totalValuation := int64(0)
	for _, startup := range gs.AvailableStartups {
		totalValuation += startup.Valuation
	}
	yellow := color.New(color.FgYellow)
	yellow.Printf("Total Company Valuation: $%s\n", FormatMoney(totalValuation))

	// Show available startups
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("AVAILABLE STARTUPS:")

	// Get player upgrades for display
	playerUpgrades, err := database.GetPlayerUpgrades(gs.PlayerName)
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
		DisplayStartup(startup, i, gs.Portfolio.Cash, playerUpgrades)
	}
	fmt.Println(strings.Repeat("=", 50))

	// Show syndicate opportunities if unlocked
	if playerLevel >= 2 && len(gs.SyndicateOpportunities) > 0 {
		magenta := color.New(color.FgMagenta, color.Bold)
		fmt.Println("\n" + strings.Repeat("=", 50))
		magenta.Println("🤝 SYNDICATE OPPORTUNITIES (Co-Invest with Other VCs)")
		fmt.Println(strings.Repeat("=", 50))

		for i, opp := range gs.SyndicateOpportunities {
			startup := gs.AvailableStartups[opp.StartupIndex]
			cyan.Printf("\n[SYNDICATE %d] %s\n", i+1, opp.CompanyName)
			fmt.Printf("   %s\n", startup.Description)
			yellow.Printf("   Lead Investor: %s (%s)\n", opp.LeadInvestor, opp.LeadInvestorFirm)
			fmt.Printf("   Total Round Size: $%s\n", FormatMoney(opp.TotalRoundSize))
			fmt.Printf("   Company Valuation: $%s\n", FormatMoney(opp.Valuation))
			green.Printf("   Your Share: $%s - $%s\n", FormatMoney(opp.YourMinShare), FormatMoney(opp.YourMaxShare))
			fmt.Printf("   Benefits:\n")
			for _, benefit := range opp.Benefits {
				fmt.Printf("     • %s\n", benefit)
			}
			fmt.Println()
		}
		fmt.Println(strings.Repeat("=", 50))
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		// Auto-start if out of money
		if gs.Portfolio.Cash <= 0 {
			color.Yellow("\n⚠️  Out of investment capital! Starting game...")
			gs.AIPlayerMakeInvestments()
			time.Sleep(2 * time.Second)
			break
		}

		fmt.Printf("\nEnter company number (1-%d) to invest, 's' for syndicates, or press Enter to start: ", len(gs.AvailableStartups))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "done" || input == "" {
			// Have AI players make their investments too
			gs.AIPlayerMakeInvestments()
			color.Green("\n✓ AI players have made their investments!")
			break
		}

		// Handle syndicate investments
		if (input == "s" || input == "S") && playerLevel >= 2 && len(gs.SyndicateOpportunities) > 0 {
			handleSyndicateInvestment(gs)
			continue
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
		fmt.Printf("   Valuation: $%s\n", FormatMoney(startup.Valuation))
		yellow.Printf("   Max Investment Available: $%s", FormatMoney(maxInvestmentDisplay))
		if maxInvestmentDisplay < maxInvestment {
			fmt.Printf(" (limited by cash, max would be $%s)", FormatMoney(maxInvestment))
		} else {
			if hasSuperProRata {
				fmt.Printf(" (50%% of valuation)")
			} else {
				fmt.Printf(" (20%% of valuation)")
			}
		}
		fmt.Println()
		fmt.Printf("\nEnter investment amount ($10,000 - $%s, or 0 to skip): $", FormatMoney(maxInvestmentDisplay))
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
			color.Red("Maximum investment is $%s (%s of company valuation: $%s)", FormatMoney(maxInvestmentDisplay), maxPercentText, FormatMoney(startup.Valuation))
			continue
		}

		// Perform due diligence (Manual Mode only)
		ddLevel := ShowDueDiligenceMenu(gs, &gs.AvailableStartups[companyNum-1], amount, autoMode)
		if ddLevel == "cancelled" {
			color.Yellow("Investment cancelled after due diligence")
			continue
		}
		
		// Show term options for investments $50k+
		var selectedTerms game.InvestmentTerms
		if amount >= 50000 {
			selectedTerms = SelectInvestmentTerms(gs, &gs.AvailableStartups[companyNum-1], amount)
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
			// Initialize founder relationship for new investment
			if len(gs.Portfolio.Investments) > 0 {
				lastInv := &gs.Portfolio.Investments[len(gs.Portfolio.Investments)-1]
				lastInv.FounderName = game.GenerateFounderName()
				lastInv.RelationshipScore = game.CalculateInitialRelationship(
					selectedTerms, 
					ddLevel != "none", 
					amount)
				lastInv.DDLevel = ddLevel
				lastInv.HasDueDiligence = ddLevel != "none"
				lastInv.LastInteraction = gs.Portfolio.Turn
				lastInv.ValueAddProvided = 0
				
				// Apply reputation bonus
				if gs.PlayerReputation != nil {
					bonus := game.GetReputationBonus(gs.PlayerReputation)
					game.ApplyReputationBonusToInvestment(lastInv, bonus)
				}
			}
			
			color.Green("%s Investment successful!", ascii.Check)
			fmt.Printf("Cash remaining: $%s\n", FormatMoney(gs.Portfolio.Cash))
			fmt.Printf("Terms: %s\n", selectedTerms.Type)
		}
	}
}

func AskForFirmName(username string) string {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	defaultFirmName := game.GenerateDefaultFirmName(username)

	cyan.Println("\n" + strings.Repeat("=", 60))
	cyan.Println("              VC FIRM NAME")
	cyan.Println(strings.Repeat("=", 60))

	yellow.Printf("\nEnter your VC firm name (default: %s): ", defaultFirmName)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		firmName := defaultFirmName
		green.Printf("\n✓ Using default firm name: %s\n", firmName)
		time.Sleep(1 * time.Second)
		return firmName
	}

	green.Printf("\n✓ Firm name set: %s\n", input)
	time.Sleep(1 * time.Second)
	return input
}

func PlayVCMode(username string) {
	// Select difficulty (passing username for level checking)
	difficulty := SelectDifficulty(username)
	clear.ClearIt()

	// Ask for firm name
	firmName := AskForFirmName(username)
	clear.ClearIt()

	// Ask for automated mode
	autoMode := AskForAutomatedMode()
	clear.ClearIt()

	// Get player upgrades
	playerUpgrades, err := database.GetPlayerUpgrades(username)
	if err != nil {
		playerUpgrades = []string{}
	}

	// Load player reputation
	dbRep, err := database.GetVCReputation(username)
	if err != nil {
		// New player - create default reputation
		dbRep = &database.VCReputation{
			PlayerName:       username,
			PerformanceScore: 50.0,
			FounderScore:     50.0,
			MarketScore:      50.0,
			TotalGamesPlayed: 0,
			SuccessfulExits:  0,
			AvgROILast5:      0.0,
		}
	}

	// Initialize game first (so we can get randomized AI players)
	gs := game.NewGame(username, firmName, difficulty, playerUpgrades)
	
	// Set player reputation in game state
	gs.PlayerReputation = &game.VCReputation{
		PlayerName:       dbRep.PlayerName,
		PerformanceScore: dbRep.PerformanceScore,
		FounderScore:     dbRep.FounderScore,
		MarketScore:      dbRep.MarketScore,
		TotalGamesPlayed: dbRep.TotalGamesPlayed,
		SuccessfulExits:  dbRep.SuccessfulExits,
		AvgROILast5:      dbRep.AvgROILast5,
	}

	// Display welcome and rules (with upgrades and randomized AI players)
	DisplayWelcome(username, difficulty, playerUpgrades, gs.AIPlayers)
	
	// Show reputation summary
	DisplayReputationSummary(gs.PlayerReputation)

	// Get player level to check for syndicate unlock
	profile, err := database.GetPlayerProfile(username)
	playerLevel := 1
	if err == nil {
		playerLevel = profile.Level
	}

	// Generate syndicate opportunities if unlocked (level 2+)
	if playerLevel >= 2 {
		gs.GenerateSyndicateOpportunities(playerLevel)
	}

	// Investment phase at start
	investmentPhase(gs, playerLevel, autoMode)

	// Main game loop
	for !gs.IsGameOver() {
		PlayTurn(gs, autoMode)
	}

	// Show final score
	DisplayFinalScore(gs)

	// Save score to database
	netWorth, roi, successfulExits := gs.GetFinalScore()
	score := database.GameScore{
		PlayerName:      gs.PlayerName,
		FinalNetWorth:   netWorth,
		ROI:             roi,
		SuccessfulExits: successfulExits,
		TurnsPlayed:     gs.Portfolio.Turn - 1,
		Difficulty:      gs.Difficulty.Name,
		PlayedAt:        time.Now(),
	}

	err = database.SaveGameScore(score)
	if err != nil {
		color.Yellow("\nWarning: Could not save score: %v", err)
	} else {
		color.Green("\n%s Score saved to local leaderboard!", ascii.Check)
	}

	// Ask if player wants to submit to global leaderboard
	AskToSubmitToGlobalLeaderboard(score)

	// Check for achievements
	CheckAndUnlockAchievements(gs)
	
	// Update reputation based on game performance
	updatePlayerReputation(gs, username, roi, successfulExits)

	fmt.Print("\nPress 'Enter' to return to main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func updatePlayerReputation(gs *game.GameState, username string, roi float64, successfulExits int) {
	if gs.PlayerReputation == nil {
		return // Safety check
	}
	
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	
	// Calculate average founder relationship
	totalRelationship := 0.0
	count := 0
	for _, inv := range gs.Portfolio.Investments {
		if inv.FounderName != "" {
			totalRelationship += inv.RelationshipScore
			count++
		}
	}
	avgFounderRelationship := 50.0
	if count > 0 {
		avgFounderRelationship = totalRelationship / float64(count)
	}
	
	// Get achievement points (approximate from career level)
	profile, _ := database.GetPlayerProfile(username)
	achievementPoints := profile.Level * 50 // Rough estimate
	
	// Calculate win streak (simple: this game was a win if ROI > 0)
	winStreak := 0
	if roi > 0 {
		winStreak = 1 // Simplified for now
	}
	
	// Save current reputation for comparison
	oldRep := *gs.PlayerReputation
	
	// Update reputation
	hadSuccessfulExit := successfulExits > 0
	updatedRep := game.UpdateReputationAfterGame(
		gs.PlayerReputation,
		roi,
		hadSuccessfulExit,
		achievementPoints,
		winStreak)
	
	// Update founder score
	updatedRep.UpdateFounderScore(avgFounderRelationship)
	
	// Save to database
	dbRep := &database.VCReputation{
		PlayerName:       updatedRep.PlayerName,
		PerformanceScore: updatedRep.PerformanceScore,
		FounderScore:     updatedRep.FounderScore,
		MarketScore:      updatedRep.MarketScore,
		TotalGamesPlayed: updatedRep.TotalGamesPlayed,
		SuccessfulExits:  updatedRep.SuccessfulExits,
		AvgROILast5:      updatedRep.AvgROILast5,
	}
	
	err := database.SaveVCReputation(dbRep)
	if err != nil {
		color.Yellow("\nWarning: Could not save reputation: %v", err)
		return
	}
	
	// Display reputation changes
	fmt.Println()
	cyan.Println(strings.Repeat("=", 70))
	cyan.Println("                    REPUTATION UPDATE")
	cyan.Println(strings.Repeat("=", 70))
	
	fmt.Println()
	fmt.Printf("Performance: %.1f → ", oldRep.PerformanceScore)
	if updatedRep.PerformanceScore > oldRep.PerformanceScore {
		green.Printf("%.1f (+%.1f)\n", updatedRep.PerformanceScore, 
			updatedRep.PerformanceScore-oldRep.PerformanceScore)
	} else if updatedRep.PerformanceScore < oldRep.PerformanceScore {
		color.Red("%.1f (%.1f)\n", updatedRep.PerformanceScore,
			updatedRep.PerformanceScore-oldRep.PerformanceScore)
	} else {
		fmt.Printf("%.1f (no change)\n", updatedRep.PerformanceScore)
	}
	
	fmt.Printf("Founder:     %.1f → ", oldRep.FounderScore)
	if updatedRep.FounderScore > oldRep.FounderScore {
		green.Printf("%.1f (+%.1f)\n", updatedRep.FounderScore,
			updatedRep.FounderScore-oldRep.FounderScore)
	} else if updatedRep.FounderScore < oldRep.FounderScore {
		color.Red("%.1f (%.1f)\n", updatedRep.FounderScore,
			updatedRep.FounderScore-oldRep.FounderScore)
	} else {
		fmt.Printf("%.1f (no change)\n", updatedRep.FounderScore)
	}
	
	fmt.Printf("Market:      %.1f → ", oldRep.MarketScore)
	if updatedRep.MarketScore > oldRep.MarketScore {
		green.Printf("%.1f (+%.1f)\n", updatedRep.MarketScore,
			updatedRep.MarketScore-oldRep.MarketScore)
	} else if updatedRep.MarketScore < oldRep.MarketScore {
		color.Red("%.1f (%.1f)\n", updatedRep.MarketScore,
			updatedRep.MarketScore-oldRep.MarketScore)
	} else {
		fmt.Printf("%.1f (no change)\n", updatedRep.MarketScore)
	}
	
	oldAggregate := oldRep.GetAggregateReputation()
	newAggregate := updatedRep.GetAggregateReputation()
	
	fmt.Println()
	fmt.Printf("Overall:     %.1f → ", oldAggregate)
	if newAggregate > oldAggregate {
		green.Printf("%.1f (%s)\n", newAggregate, updatedRep.GetReputationLevel())
	} else if newAggregate < oldAggregate {
		yellow.Printf("%.1f (%s)\n", newAggregate, updatedRep.GetReputationLevel())
	} else {
		fmt.Printf("%.1f (%s)\n", newAggregate, updatedRep.GetReputationLevel())
	}
	
	// Show tier change if any
	oldTier := oldRep.GetDealQualityTier()
	newTier := updatedRep.GetDealQualityTier()
	if oldTier != newTier {
		fmt.Println()
		green.Printf("🎉 Deal Flow Quality: %s → %s\n", oldTier, newTier)
	}
	
	cyan.Println(strings.Repeat("=", 70))
}

func FormatMoney(amount int64) string {
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

func FormatNumber(n int) string {
	return FormatMoney(int64(n))
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
