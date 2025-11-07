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

// ShowValueAddMenu displays value-add opportunities and handles selection
// Only shown in Manual Mode
func ShowValueAddMenu(gs *game.GameState, autoMode bool) {
	if autoMode {
		// Skip in automated mode
		return
	}

	opportunities := gs.GetValueAddOpportunities()

	if len(opportunities) == 0 {
		return // No opportunities
	}

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Count actions already taken this turn
	actionsThisTurn := 0
	for _, action := range gs.ActiveValueAddActions {
		if action.AppliedTurn == gs.Portfolio.Turn {
			actionsThisTurn++
		}
	}

	// Don't show menu if already at limit
	if actionsThisTurn >= 2 {
		return
	}

	fmt.Println()
	cyan.Println(strings.Repeat("=", 70))
	cyan.Println("              💼 VALUE-ADD OPPORTUNITIES")
	cyan.Println(strings.Repeat("=", 70))

	yellow.Println("\nYou can provide operational support to your portfolio companies.")
	remainingPoints := 2 - actionsThisTurn
	fmt.Printf("You have %d attention point%s remaining this turn. Each action costs 1 point.\n\n",
		remainingPoints, map[bool]string{true: "", false: "s"}[remainingPoints == 1])

	// Show available actions
	actionTypes := game.GetAvailableValueAddTypes()

	fmt.Println("AVAILABLE VALUE-ADD ACTIONS:")
	for i, actionType := range actionTypes {
		fmt.Printf("\n%d. %s - $%s\n", i+1, actionType.Name, FormatMoney(actionType.Cost))
		fmt.Printf("   %s\n", actionType.Description)
		if actionType.RequiresBoardSeat {
			fmt.Printf("   Requires: Board seat\n")
		} else {
			fmt.Printf("   Requires: %.1f%% equity\n", actionType.MinEquityPct)
		}
	}

	fmt.Println()
	green.Println("ELIGIBLE COMPANIES:")
	for i, companyName := range opportunities {
		// Find the investment
		for _, inv := range gs.Portfolio.Investments {
			if inv.CompanyName == companyName {
				emoji := game.GetRelationshipEmoji(inv.RelationshipScore)
				fmt.Printf("%d. %s (%.1f%% equity, %s %.1f/100 relationship)\n",
					i+1, companyName, inv.EquityPercent, emoji, inv.RelationshipScore)
				break
			}
		}
	}

	fmt.Println()
	fmt.Print("Select company (1-%d) or press Enter to skip: ", len(opportunities))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return // Skip
	}

	companyIdx, err := strconv.Atoi(input)
	if err != nil || companyIdx < 1 || companyIdx > len(opportunities) {
		color.Red("Invalid selection")
		return
	}

	selectedCompany := opportunities[companyIdx-1]

	// Now select action type
	fmt.Print("\nSelect action type (1-%d) or press Enter to cancel: ", len(actionTypes))
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return // Cancel
	}

	actionIdx, err := strconv.Atoi(input)
	if err != nil || actionIdx < 1 || actionIdx > len(actionTypes) {
		color.Red("Invalid selection")
		return
	}

	selectedAction := actionTypes[actionIdx-1]

	// Check if can provide
	can, reason := gs.CanProvideValueAdd(selectedCompany, selectedAction)
	if !can {
		color.Red("Cannot provide this value-add: %s", reason)
		fmt.Print("\nPress Enter to continue...")
		reader.ReadBytes('\n')
		return
	}

	// Confirm
	fmt.Printf("\nProvide %s to %s for $%s?\n",
		selectedAction.Name, selectedCompany, FormatMoney(selectedAction.Cost))
	fmt.Print("Confirm (y/n): ")
	input, _ = reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" && input != "yes" {
		fmt.Println("Cancelled")
		return
	}

	// Execute
	err = gs.ProvideValueAdd(selectedCompany, selectedAction.ID)
	if err != nil {
		color.Red("Error: %v", err)
	} else {
		green.Printf("\n✓ Successfully provided %s to %s!\n", selectedAction.Name, selectedCompany)
		fmt.Printf("Relationship improved, value-add effects will apply over next %d turns\n",
			selectedAction.Duration)
	}

	fmt.Print("\nPress Enter to continue...")
	reader.ReadBytes('\n')

	// Recount actions this turn (it increased by 1 after this action)
	actionsNow := 0
	for _, action := range gs.ActiveValueAddActions {
		if action.AppliedTurn == gs.Portfolio.Turn {
			actionsNow++
		}
	}

	// Allow another action if under the limit and opportunities remain
	if actionsNow < 2 && len(gs.GetValueAddOpportunities()) > 0 {
		fmt.Print("\nProvide another value-add action? (y/n): ")
		input, _ = reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "y" || input == "yes" {
			ShowValueAddMenu(gs, autoMode)
		}
	} else if actionsNow >= 2 {
		yellow.Println("\n✓ Maximum 2 value-add actions per turn reached")
	}
}

// DisplayValueAddHistory shows value-add actions taken
func DisplayValueAddHistory(gs *game.GameState) {
	if len(gs.ActiveValueAddActions) == 0 {
		fmt.Println("\nNo active value-add actions")
		return
	}

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	cyan.Println("ACTIVE VALUE-ADD ACTIONS:")
	fmt.Println()

	for _, action := range gs.ActiveValueAddActions {
		turnsRemaining := action.Duration - (gs.Portfolio.Turn - action.AppliedTurn)
		if turnsRemaining > 0 {
			yellow.Printf("%s - %s\n", action.CompanyName, action.Description)
			fmt.Printf("  Effects: +%.1f%% valuation boost, +%.1f relationship\n",
				action.ValuationBoost*100, action.Relationship)
			fmt.Printf("  Duration: %d turns remaining\n", turnsRemaining)
			fmt.Println()
		}
	}
}
