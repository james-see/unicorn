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

// ShowSecondaryMarketOffers displays secondary market offers (Manual Mode only)
func ShowSecondaryMarketOffers(gs *game.GameState, autoMode bool) {
	if autoMode {
		// Auto-decline all offers in automated mode
		gs.SecondaryMarketOffers = []game.SecondaryOffer{}
		return
	}

	if len(gs.SecondaryMarketOffers) == 0 {
		return // No offers
	}

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)

	fmt.Println()
	cyan.Println(strings.Repeat("=", 70))
	cyan.Println("            💰 SECONDARY MARKET OFFERS")
	cyan.Println(strings.Repeat("=", 70))

	yellow.Println("\nYou have received offers to buy stakes in your portfolio companies.")
	fmt.Println("Secondary sales allow you to realize gains early or cut losses.")
	fmt.Println()

	for i, offer := range gs.SecondaryMarketOffers {
		fmt.Printf("OFFER #%d\n", i+1)
		fmt.Println(strings.Repeat("-", 70))

		// Find the investment
		var inv *game.Investment
		for j := range gs.Portfolio.Investments {
			if gs.Portfolio.Investments[j].CompanyName == offer.CompanyName {
				inv = &gs.Portfolio.Investments[j]
				break
			}
		}

		if inv == nil {
			continue
		}

		fmt.Printf("Company:         %s\n", offer.CompanyName)
		fmt.Printf("Buyer:           %s (%s)\n", offer.BuyerName, offer.BuyerFirm)
		fmt.Printf("Your Equity:     %.2f%%\n", offer.EquityOffered)
		fmt.Printf("Current Value:   $%s\n", FormatMoney(offer.CurrentValue))
		magenta.Printf("Offer Amount:    $%s (%.1f%% of value)\n",
			FormatMoney(offer.OfferAmount), offer.OfferPercent*100)

		// Calculate ROI
		roi := game.CalculateSecondaryROI(offer, inv.AmountInvested)
		if roi > 0 {
			green.Printf("Your ROI:        +%.1f%%\n", roi)
		} else {
			color.Red("Your ROI:        %.1f%%\n", roi)
		}

		// AI recommendation
		shouldAccept, reason := game.ShouldAcceptSecondaryOffer(offer, *inv)
		fmt.Printf("\nAI Recommendation: ")
		if shouldAccept {
			green.Printf("ACCEPT\n")
		} else {
			yellow.Printf("HOLD\n")
		}
		fmt.Printf("Reason: %s\n", reason)

		fmt.Printf("\nExpires in %d turns\n", offer.ExpiresIn)
		fmt.Println()
	}

	fmt.Print("Select offer to accept (1-%d), 'd' to decline all, or Enter to decide later: ",
		len(gs.SecondaryMarketOffers))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input == "" {
		return // Decide later
	}

	if input == "d" {
		// Decline all
		gs.SecondaryMarketOffers = []game.SecondaryOffer{}
		color.Yellow("\nDeclined all offers")
		fmt.Print("\nPress Enter to continue...")
		reader.ReadBytes('\n')
		return
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(gs.SecondaryMarketOffers) {
		color.Red("Invalid selection")
		return
	}

	// Confirm acceptance
	offer := gs.SecondaryMarketOffers[choice-1]
	fmt.Printf("\nAccept $%s offer for %s from %s?\n",
		FormatMoney(offer.OfferAmount), offer.CompanyName, offer.BuyerName)
	fmt.Print("Confirm (y/n): ")
	input, _ = reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" && input != "yes" {
		fmt.Println("Offer declined")
		gs.DeclineSecondaryOffer(choice - 1)
		return
	}

	// Accept offer
	err = gs.AcceptSecondaryOffer(choice - 1)
	if err != nil {
		color.Red("Error accepting offer: %v", err)
	} else {
		green.Printf("\n✓ Successfully sold stake in %s for $%s!\n",
			offer.CompanyName, FormatMoney(offer.OfferAmount))
		fmt.Printf("Proceeds added to cash reserves.\n")
	}

	fmt.Print("\nPress Enter to continue...")
	reader.ReadBytes('\n')

	// Show remaining offers if any
	if len(gs.SecondaryMarketOffers) > 0 {
		ShowSecondaryMarketOffers(gs, autoMode)
	}
}
