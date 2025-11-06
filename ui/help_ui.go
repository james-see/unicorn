package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/ascii"
	"github.com/jamesacampbell/unicorn/clear"
)


func DisplayInvestingFAQ() {
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

func DisplayHelpGuide() {
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
		DisplayInvestingFAQ()
		return
	} else if choice == "3" {
		return
	}

	clear.ClearIt()
	cyan.Println("\n" + strings.Repeat("?", 70))
	cyan.Println("              GAME OVERVIEW & RULES")
	cyan.Println(strings.Repeat("?", 70))

	yellow.Printf("\n%s GAME OVERVIEW\n", ascii.Lightbulb)
	fmt.Println("TWO GAME MODES:")
	fmt.Println("• VC Investor Mode: Invest in startups and build a portfolio")
	fmt.Println("• Startup Founder Mode: Build your own startup from the ground up")
	fmt.Println()
	fmt.Println("Both modes feature persistent progression, achievements, and XP rewards!")

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

	yellow.Printf("\n%s PLAYER PROGRESSION SYSTEM\n", ascii.Trophy)
	fmt.Println("Earn XP and level up after every game!")
	fmt.Printf("%s 50 Levels: From 'Novice Investor' to 'Titan of Industry'\n", ascii.Star)
	fmt.Printf("%s XP Sources:\n", ascii.Chart)
	fmt.Println("   • 100 XP - Game completion")
	fmt.Println("   • +50 XP - Positive ROI")
	fmt.Println("   • +200 XP per successful exit (5x+)")
	fmt.Println("   • +0/50/100/200 XP - Difficulty bonus")
	fmt.Println("   • +10 XP × points - New achievements")
	fmt.Println("   • +500 XP - IPO exit (Founder mode)")
	fmt.Println("   • +300 XP - Acquisition exit (Founder mode)")
	fmt.Printf("%s Level Unlocks:\n", ascii.Zap)
	fmt.Println("   • Level 5: Hard Difficulty")
	fmt.Println("   • Level 10: Expert Difficulty + Analytics Dashboard")
	fmt.Println("   • Level 15: Nightmare Mode + Achievement Chains")
	fmt.Println("   • Level 20+: Game Modifiers, Secondary Markets")

	yellow.Printf("\n%s ACHIEVEMENTS & UPGRADES\n", ascii.Medal)
	fmt.Printf("%s 45+ Achievements: Earn points across 6 categories\n", ascii.Star)
	fmt.Printf("%s Achievement Points: Spend on permanent upgrades\n", ascii.Coin)
	fmt.Printf("%s VC Upgrades: Better terms, more intel, special abilities\n", ascii.Portfolio)
	fmt.Printf("%s Founder Upgrades: Lower costs, better metrics, advantages\n", ascii.Rocket)
	fmt.Printf("%s Hidden Achievements: Secret challenges to discover\n", ascii.Crown)
	fmt.Printf("%s Progressive Achievements: Track long-term progress\n", ascii.Chart)

	yellow.Printf("\n%s ANALYTICS DASHBOARD\n", ascii.Chart)
	fmt.Println("Deep insights into your performance:")
	fmt.Printf("%s Trend Analysis: 7-day, 30-day, all-time trends\n", ascii.Star)
	fmt.Printf("%s Difficulty Breakdown: Per-difficulty statistics\n", ascii.Shield)
	fmt.Printf("%s Monthly Reports: Historical performance tracking\n", ascii.Calendar)
	fmt.Printf("%s AI Insights: Smart recommendations based on patterns\n", ascii.Lightbulb)
	fmt.Printf("%s Global Comparisons: Compare to global averages\n", ascii.Trophy)

	yellow.Printf("\n%s STRATEGY TIPS\n", ascii.Lightbulb)
	fmt.Printf("%s Diversify: Don't put everything in one company\n", ascii.Portfolio)
	fmt.Printf("%s Balance: Mix high-risk and low-risk investments\n", ascii.Shield)
	fmt.Printf("%s Sectors: Different industries perform differently\n", ascii.Building)
	fmt.Printf("%s Research: Read company metrics carefully\n", ascii.Star)
	fmt.Printf("%s Early Entry: Invest before dilution hits\n", ascii.Target)
	fmt.Printf("%s Level Up: Higher levels unlock more features\n", ascii.Trophy)
	fmt.Printf("%s Achievements: Unlock upgrades for permanent advantages\n", ascii.Medal)

	cyan.Println("\n" + strings.Repeat("=", 70))
	fmt.Print("\nPress 'Enter' to return to menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}