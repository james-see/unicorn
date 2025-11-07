package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/game"
)

// ShowDueDiligenceMenu shows DD options before investment (Manual Mode only)
func ShowDueDiligenceMenu(gs *game.GameState, startup *game.Startup, amount int64, autoMode bool) string {
	if autoMode {
		// Skip in automated mode - go straight to investment
		return "none"
	}

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta)

	fmt.Println()
	cyan.Println(strings.Repeat("=", 70))
	cyan.Println("              🔍 DUE DILIGENCE OPTIONS")
	cyan.Println(strings.Repeat("=", 70))

	fmt.Printf("\nYou're about to invest $%s in %s\n", FormatMoney(amount), startup.Name)
	yellow.Println("\nWould you like to perform due diligence first?")
	fmt.Println("DD can reveal red flags, hidden gems, and improve founder relationships.")

	levels := game.GetDDLevels()

	fmt.Println()
	for i, level := range levels {
		if level.ID == "none" {
			fmt.Printf("%d. %s (FREE)\n", i+1, level.Name)
			fmt.Printf("   %s\n", level.Description)
		} else {
			fmt.Printf("%d. %s ($%s, %d days)\n", i+1, level.Name, FormatMoney(level.Cost), level.Duration)
			fmt.Printf("   %s\n", level.Description)
			magenta.Printf("   Reveals: %s\n", strings.Join(level.Reveals, ", "))
		}
		fmt.Println()
	}

	fmt.Print("Select option (1-4): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(levels) {
		fmt.Println("Invalid choice, proceeding without DD")
		return "none"
	}

	selectedLevel := levels[choice-1]

	// Check if have enough cash
	if selectedLevel.Cost > gs.Portfolio.Cash {
		color.Red("Insufficient funds for DD. Proceeding without DD.")
		fmt.Print("\nPress Enter to continue...")
		reader.ReadBytes('\n')
		return "none"
	}

	if selectedLevel.ID == "none" {
		return "none"
	}

	// Perform DD
	fmt.Printf("\nPerforming %s...\n", selectedLevel.Name)

	findings := game.PerformDueDiligence(startup, selectedLevel.ID)

	// Deduct cost
	gs.Portfolio.Cash -= selectedLevel.Cost

	// Display findings
	fmt.Println()
	cyan.Println("DUE DILIGENCE FINDINGS:")
	fmt.Println()

	for _, finding := range findings {
		switch finding.Type {
		case "red_flag":
			color.Red("🚩 RED FLAG - %s: %s\n", finding.Category, finding.Description)
		case "green_flag":
			color.Green("✓ POSITIVE - %s: %s\n", finding.Category, finding.Description)
		default:
			fmt.Printf("ℹ️  INFO - %s: %s\n", finding.Category, finding.Description)
		}
	}

	// Check if should block investment
	shouldBlock, blockReason := game.ShouldBlockInvestment(findings)
	if shouldBlock {
		fmt.Println()
		color.Red("⚠️  WARNING: %s\n", blockReason)
		fmt.Print("\nStill proceed with investment? (y/n): ")
		input, _ = reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input != "y" && input != "yes" {
			fmt.Println("Investment cancelled. Returning to menu.")
			gs.Portfolio.Cash += selectedLevel.Cost // Refund DD cost
			return "cancelled"
		}
	}

	// Apply findings to startup
	game.ApplyDDFindings(startup, findings)

	fmt.Print("\nPress Enter to proceed with investment...")
	reader.ReadBytes('\n')

	return selectedLevel.ID
}
