package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/game"
)

// DisplayPortfolioDashboard shows detailed portfolio metrics
func DisplayPortfolioDashboard(gs *game.GameState) {
	clear.ClearIt()

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	magenta := color.New(color.FgMagenta, color.Bold)

	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                  PORTFOLIO DASHBOARD")
	cyan.Println(strings.Repeat("=", 70))

	// Overall Portfolio Summary
	fmt.Println()
	magenta.Println("[PORTFOLIO SUMMARY]")
	fmt.Printf("   Total Investments: %d companies\n", len(gs.Portfolio.Investments))
	fmt.Printf("   Total Invested: $%s\n", FormatMoney(gs.GetTotalInvested()))
	fmt.Printf("   Current Portfolio Value: $%s\n", FormatMoney(gs.GetPortfolioValue()))
	
	totalProfit := gs.GetPortfolioValue() - gs.GetTotalInvested()
	profitColor := green
	profitSign := "+"
	if totalProfit < 0 {
		profitColor = red
		profitSign = ""
	}
	profitColor.Printf("   Total Profit/Loss: %s$%s\n", profitSign, FormatMoney(abs(totalProfit)))
	
	if gs.GetTotalInvested() > 0 {
		portfolioROI := (float64(totalProfit) / float64(gs.GetTotalInvested())) * 100.0
		roiColor := green
		if portfolioROI < 0 {
			roiColor = red
		}
		roiColor.Printf("   Portfolio ROI: %.2f%%\n", portfolioROI)
	}

	// Best and Worst Performers
	fmt.Println()
	magenta.Println("[TOP PERFORMERS]")
	bestPerformers := gs.GetBestPerformers(3)
	if len(bestPerformers) > 0 {
		for i, inv := range bestPerformers {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			profit := value - inv.AmountInvested
			roi := (float64(profit) / float64(inv.AmountInvested)) * 100.0
			green.Printf("   %d. %s: $%s → $%s (%.1f%% ROI)\n",
				i+1, inv.CompanyName, FormatMoney(inv.AmountInvested), FormatMoney(value), roi)
		}
	} else {
		fmt.Println("   No investments yet")
	}

	fmt.Println()
	magenta.Println("[WORST PERFORMERS]")
	worstPerformers := gs.GetWorstPerformers(3)
	if len(worstPerformers) > 0 {
		for i, inv := range worstPerformers {
			value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
			profit := value - inv.AmountInvested
			roi := (float64(profit) / float64(inv.AmountInvested)) * 100.0
			roiColor := green
			if roi < 0 {
				roiColor = red
			}
			roiColor.Printf("   %d. %s: $%s → $%s (%.1f%% ROI)\n",
				i+1, inv.CompanyName, FormatMoney(inv.AmountInvested), FormatMoney(value), roi)
		}
	} else {
		fmt.Println("   No investments yet")
	}

	// Sector Breakdown
	fmt.Println()
	magenta.Println("[SECTOR BREAKDOWN]")
	sectorBreakdown := gs.GetSectorBreakdown()
	if len(sectorBreakdown) > 0 {
		for sector, data := range sectorBreakdown {
			fmt.Printf("   %s:\n", sector)
			fmt.Printf("      Companies: %d\n", data.Count)
			fmt.Printf("      Invested: $%s\n", FormatMoney(data.TotalInvested))
			fmt.Printf("      Current Value: $%s\n", FormatMoney(data.CurrentValue))
			sectorProfit := data.CurrentValue - data.TotalInvested
			sectorROI := 0.0
			if data.TotalInvested > 0 {
				sectorROI = (float64(sectorProfit) / float64(data.TotalInvested)) * 100.0
			}
			roiColor := green
			if sectorROI < 0 {
				roiColor = red
			}
			roiColor.Printf("      ROI: %.1f%%\n", sectorROI)
		}
	} else {
		fmt.Println("   No investments yet")
	}

	// Investment Distribution
	fmt.Println()
	magenta.Println("[INVESTMENT DISTRIBUTION]")
	positiveCount := 0
	negativeCount := 0
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		if value > inv.AmountInvested {
			positiveCount++
		} else {
			negativeCount++
		}
	}
	totalCount := len(gs.Portfolio.Investments)
	if totalCount > 0 {
		winRate := (float64(positiveCount) / float64(totalCount)) * 100.0
		fmt.Printf("   Winning Investments: %d (%.1f%%)\n", positiveCount, winRate)
		fmt.Printf("   Losing Investments: %d (%.1f%%)\n", negativeCount, 100.0-winRate)
	} else {
		fmt.Println("   No investments yet")
	}

	// Financial Summary
	fmt.Println()
	magenta.Println("[FINANCIAL SUMMARY]")
	fmt.Printf("   Cash Available: $%s\n", FormatMoney(gs.Portfolio.Cash))
	fmt.Printf("   Follow-on Reserve: $%s\n", FormatMoney(gs.Portfolio.FollowOnReserve))
	fmt.Printf("   Management Fees Paid: $%s\n", FormatMoney(gs.Portfolio.ManagementFeesCharged))
	fmt.Printf("   Net Worth: $%s\n", FormatMoney(gs.Portfolio.NetWorth))
	
	// LP Commitments
	fmt.Println()
	magenta.Println("[LP COMMITMENTS]")
	fmt.Printf("   LP Committed Capital: $%s\n", FormatMoney(gs.Portfolio.LPCommittedCapital))
	fmt.Printf("   LP Called Capital: $%s\n", FormatMoney(gs.Portfolio.LPCalledCapital))
	remainingCommitment := gs.Portfolio.LPCommittedCapital - gs.Portfolio.LPCalledCapital
	fmt.Printf("   Remaining Commitment: $%s\n", FormatMoney(remainingCommitment))
	if gs.Portfolio.LastCapitalCallTurn > 0 {
		fmt.Printf("   Last Capital Call: Month %d\n", gs.Portfolio.LastCapitalCallTurn)
	}
	// Show next capital call
	nextCallTurn := -1
	for _, turn := range gs.Portfolio.CapitalCallSchedule {
		if turn > gs.Portfolio.Turn {
			nextCallTurn = turn
			break
		}
	}
	if nextCallTurn > 0 {
		fmt.Printf("   Next Capital Call: Month %d\n", nextCallTurn)
	} else if remainingCommitment > 0 {
		fmt.Println("   All capital calls completed")
	}
	
	// Carry Interest Calculation
	fmt.Println()
	magenta.Println("[CARRY INTEREST PROJECTION]")
	totalStartingCapital := gs.Portfolio.InitialFundSize + gs.Portfolio.FollowOnReserve
	projectedCarry, hurdleReturn, excessProfit, applies := gs.CalculateCarryInterest()
	
	fmt.Printf("   Starting Capital: $%s\n", FormatMoney(totalStartingCapital))
	fmt.Printf("   Hurdle Return (40%%): $%s\n", FormatMoney(int64(hurdleReturn)))
	fmt.Printf("   Current Profit: $%s\n", FormatMoney(gs.Portfolio.NetWorth - totalStartingCapital))
	
	if applies {
		yellow.Printf("   Excess Profit: $%s\n", FormatMoney(int64(excessProfit)))
		yellow.Printf("   Projected Carry (20%%): $%s\n", FormatMoney(projectedCarry))
		yellow.Printf("   Net to LPs (after carry): $%s\n", FormatMoney(gs.Portfolio.NetWorth - projectedCarry))
		
		// Show progress to hurdle
		hurdleProgress := (float64(gs.Portfolio.NetWorth - totalStartingCapital) / hurdleReturn) * 100.0
		if hurdleProgress < 100.0 {
			fmt.Printf("   Progress to Hurdle: %.1f%%\n", hurdleProgress)
		} else {
			green.Printf("   ✓ Hurdle Exceeded!\n")
		}
	} else {
		neededForHurdle := int64(hurdleReturn) - (gs.Portfolio.NetWorth - totalStartingCapital)
		if neededForHurdle > 0 {
			fmt.Printf("   Need $%s more to reach hurdle\n", FormatMoney(neededForHurdle))
			hurdleProgress := (float64(gs.Portfolio.NetWorth - totalStartingCapital) / hurdleReturn) * 100.0
			fmt.Printf("   Progress to Hurdle: %.1f%%\n", hurdleProgress)
		} else {
			fmt.Println("   No carry interest applies (below hurdle)")
		}
	}
	
	if gs.Portfolio.CarryInterestPaid > 0 {
		green.Printf("\n   Final Carry Paid: $%s\n", FormatMoney(gs.Portfolio.CarryInterestPaid))
	}

	// Board Members Section
	fmt.Println()
	magenta.Println("[BOARD SEATS]")
	companiesWithBoardSeats := []string{}
	for _, inv := range gs.Portfolio.Investments {
		if inv.Terms.HasBoardSeat {
			companiesWithBoardSeats = append(companiesWithBoardSeats, inv.CompanyName)
		}
	}
	if len(companiesWithBoardSeats) > 0 {
		fmt.Printf("   Companies where you have board seats: %d\n", len(companiesWithBoardSeats))
		for _, companyName := range companiesWithBoardSeats {
			members := gs.GetBoardMembers(companyName)
			voteWeight := 1
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == companyName && inv.Terms.HasBoardSeat {
					voteWeight = inv.Terms.BoardSeatMultiplier
					if voteWeight == 0 {
						voteWeight = 1
					}
					break
				}
			}
			fmt.Printf("   • %s: %d vote(s) (Total board members: %d)\n", companyName, voteWeight, len(members))
		}
		fmt.Println()
		yellow.Println("Press 'b' to view board members, or Enter to continue...")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "b" || input == "board" {
			ViewBoardMembers(gs)
		}
	} else {
		fmt.Println("   No board seats")
		fmt.Println()
		yellow.Println("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

// ViewBoardMembers shows board members for companies where player has board seats
func ViewBoardMembers(gs *game.GameState) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta, color.Bold)
	green := color.New(color.FgGreen)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                  BOARD MEMBERS")
	cyan.Println(strings.Repeat("=", 70))
	
	companiesWithBoardSeats := []string{}
	for _, inv := range gs.Portfolio.Investments {
		if inv.Terms.HasBoardSeat {
			companiesWithBoardSeats = append(companiesWithBoardSeats, inv.CompanyName)
		}
	}
	
	if len(companiesWithBoardSeats) == 0 {
		fmt.Println("\nYou don't have any board seats.")
		fmt.Println()
		yellow.Println("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	
	for _, companyName := range companiesWithBoardSeats {
		members := gs.GetBoardMembers(companyName)
		
		fmt.Println()
		magenta.Printf("Company: %s\n", companyName)
		fmt.Println(strings.Repeat("-", 70))
		
		if len(members) == 0 {
			fmt.Println("  No board members found")
		} else {
			for i, member := range members {
				fmt.Printf("%d. ", i+1)
				if member.IsPlayer {
					magenta.Printf("%s (%s)", member.Name, member.Firm)
					if member.VoteWeight > 1 {
						green.Printf(" - %d votes", member.VoteWeight)
					} else {
						green.Printf(" - %d vote", member.VoteWeight)
					}
				} else {
					fmt.Printf("%s (%s)", member.Name, member.Firm)
					fmt.Printf(" - %d vote", member.VoteWeight)
				}
				fmt.Println()
			}
		}
		fmt.Println(strings.Repeat("-", 70))
	}
	
	fmt.Println()
	yellow.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

