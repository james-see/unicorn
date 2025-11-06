package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/analytics"
	"github.com/jamesacampbell/unicorn/clear"
	"github.com/jamesacampbell/unicorn/database"
)

// DisplayAnalyticsDashboard shows comprehensive analytics for a player
func DisplayAnalyticsDashboard(playerName string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                    PERFORMANCE ANALYTICS")
	cyan.Println(strings.Repeat("=", 70))
	
	// Generate trend analysis
	trendReport, err := analytics.GenerateTrendAnalysis(playerName, 30)
	if err != nil {
		red.Printf("\nError generating analytics: %v\n", err)
		return
	}
	
	// Get player stats
	stats, err := database.GetPlayerStats(playerName)
	if err != nil || stats.TotalGames == 0 {
		yellow.Println("\nNo game data available for this player.")
		yellow.Println("Play some games to see your analytics!")
		return
	}
	
	// Section 1: Performance Overview
	fmt.Println()
	yellow.Println("[1] PERFORMANCE OVERVIEW")
	
	if trendReport.Last7Days.GamesPlayed > 0 {
		fmt.Printf("    Last 7 Days:   %d games | %.0f%% win rate | Avg ROI: %.1fx\n",
			trendReport.Last7Days.GamesPlayed,
			trendReport.Last7Days.WinRate,
			trendReport.Last7Days.AvgROI)
	} else {
		fmt.Println("    Last 7 Days:   No games played")
	}
	
	if trendReport.Last30Days.GamesPlayed > 0 {
		fmt.Printf("    Last 30 Days: %d games | %.0f%% win rate | Avg ROI: %.1fx\n",
			trendReport.Last30Days.GamesPlayed,
			trendReport.Last30Days.WinRate,
			trendReport.Last30Days.AvgROI)
	} else {
		fmt.Println("    Last 30 Days:  No games played")
	}
	
	fmt.Printf("    All Time:     %d games | %.0f%% win rate | Avg ROI: %.1fx\n",
		trendReport.AllTime.GamesPlayed,
		stats.WinRate,
		trendReport.AllTime.AvgROI)
	
	// Display trend
	fmt.Println()
	fmt.Printf("    Trend: ")
	if strings.Contains(trendReport.TrendVector, "Improving") || strings.Contains(trendReport.TrendVector, "Upward") {
		green.Printf("%s\n", trendReport.TrendVector)
	} else if strings.Contains(trendReport.TrendVector, "Declining") || strings.Contains(trendReport.TrendVector, "Decline") {
		red.Printf("%s\n", trendReport.TrendVector)
	} else {
		yellow.Printf("%s\n", trendReport.TrendVector)
	}
	
	// Section 2: Difficulty Breakdown
	fmt.Println()
	yellow.Println("[2] DIFFICULTY BREAKDOWN")
	
	difficulties := []string{"easy", "medium", "hard", "expert"}
	for _, diff := range difficulties {
		diffScores, err := database.GetTopScoresByNetWorth(1000, diff)
		if err != nil {
			continue
		}
		
		// Filter for this player
		playerDiffScores := []database.GameScore{}
		for _, score := range diffScores {
			if score.PlayerName == playerName {
				playerDiffScores = append(playerDiffScores, score)
			}
		}
		
		if len(playerDiffScores) == 0 {
			continue
		}
		
		wins := 0
		totalNetWorth := int64(0)
		for _, score := range playerDiffScores {
			if score.ROI > 0 {
				wins++
			}
			totalNetWorth += score.FinalNetWorth
		}
		
		winRate := float64(wins) / float64(len(playerDiffScores)) * 100
		avgNetWorth := totalNetWorth / int64(len(playerDiffScores))
		
		fmt.Printf("    %-8s  %d games | %.0f%% wins | Avg Net Worth: $%s\n",
			strings.Title(diff)+":",
			len(playerDiffScores),
			winRate,
			FormatCurrency(avgNetWorth))
	}
	
	// Section 3: Historical Performance
	fmt.Println()
	yellow.Println("[3] HISTORICAL PERFORMANCE")
	
	// Get last 6 months of data
	now := time.Now()
	fmt.Println("    Month       Games  Wins  Avg ROI  Best Result")
	fmt.Println("    " + strings.Repeat("─", 55))
	
	for i := 0; i < 6; i++ {
		month := now.AddDate(0, -i, 0)
		monthReport, err := analytics.GetMonthlyStats(playerName, month.Year(), int(month.Month()))
		if err != nil || monthReport.GamesPlayed == 0 {
			continue
		}
		
		fmt.Printf("    %-11s %5d  %4d  %6.1fx  $%s\n",
			monthReport.Month[:3]+" "+fmt.Sprintf("%d", monthReport.Year),
			monthReport.GamesPlayed,
			monthReport.Wins,
			monthReport.AvgROI,
			FormatCurrency(monthReport.BestNetWorth))
	}
	
	// Section 4: Top Games
	fmt.Println()
	yellow.Println("[4] TOP GAMES (All Time)")
	
	topScores, err := database.GetTopScoresByPlayer(playerName, 5)
	if err == nil && len(topScores) > 0 {
		fmt.Println("    Rank  Net Worth      ROI    Difficulty  Date")
		fmt.Println("    " + strings.Repeat("─", 60))
		
		for i, score := range topScores {
			if i >= 5 {
				break
			}
			
			date := score.PlayedAt.Format("2006-01-02")
			fmt.Printf("    #%-4d $%-12s %.1fx   %-10s  %s\n",
				i+1,
				FormatCurrency(score.FinalNetWorth),
				score.ROI,
				strings.Title(score.Difficulty),
				date)
		}
	}
	
	// Section 5: Insights & Recommendations
	fmt.Println()
	yellow.Println("[5] INSIGHTS & RECOMMENDATIONS")
	
	if len(trendReport.Insights) > 0 {
		for _, insight := range trendReport.Insights {
			fmt.Printf("    %s\n", insight)
		}
	} else {
		fmt.Println("    Keep playing to generate insights!")
	}
	
	// Section 6: Comparison to Global Stats
	fmt.Println()
	yellow.Println("[6] GLOBAL COMPARISON")
	
	comparison, err := analytics.CompareToGlobal(*stats)
	if err == nil {
		fmt.Printf("    Your Avg Net Worth:  $%s\n", FormatCurrency(int64(comparison.YourAvg)))
		fmt.Printf("    Global Avg:          $%s\n", FormatCurrency(int64(comparison.GlobalAvg)))
		
		if comparison.YourAvg > comparison.GlobalAvg {
			diff := ((comparison.YourAvg - comparison.GlobalAvg) / comparison.GlobalAvg) * 100
			green.Printf("    You're %.0f%% above average! 🎉\n", diff)
		} else if comparison.YourAvg < comparison.GlobalAvg {
			diff := ((comparison.GlobalAvg - comparison.YourAvg) / comparison.GlobalAvg) * 100
			fmt.Printf("    You're %.0f%% below average\n", diff)
		}
		
		if comparison.Percentile <= 25 {
			green.Printf("    Ranking: Top %d%% of players! 🏆\n", comparison.Percentile)
		} else {
			fmt.Printf("    Ranking: Top %d%% of players\n", comparison.Percentile)
		}
	}
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayTrendChart creates an ASCII line chart for trends
func DisplayTrendChart(data []float64, title string) {
	if len(data) == 0 {
		fmt.Println("No data available for chart")
		return
	}
	
	cyan := color.New(color.FgCyan, color.Bold)
	
	cyan.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("─", 50))
	
	// Find min and max for scaling
	min, max := data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	// Create 10 rows for the chart
	rows := 10
	height := max - min
	if height == 0 {
		height = 1
	}
	
	// Draw chart from top to bottom
	for row := rows; row >= 0; row-- {
		threshold := min + (float64(row)/float64(rows))*height
		
		// Y-axis label
		fmt.Printf("%6.1f │ ", threshold)
		
		// Plot points
		for _, v := range data {
			if v >= threshold-height/float64(rows)/2 && v <= threshold+height/float64(rows)/2 {
				fmt.Print("●")
			} else if v > threshold {
				fmt.Print("│")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	
	// X-axis
	fmt.Print("       └")
	fmt.Print(strings.Repeat("─", len(data)))
	fmt.Println()
	
	fmt.Printf("        ")
	for i := range data {
		if i%5 == 0 {
			fmt.Printf("%-5d", i+1)
		}
	}
	fmt.Println()
}

// DisplayPerformanceHeatmap shows a color-coded grid of monthly performance
func DisplayPerformanceHeatmap(playerName string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                  PERFORMANCE HEATMAP")
	cyan.Println(strings.Repeat("=", 70))
	
	// Get last 12 months
	now := time.Now()
	
	fmt.Println("\nMonthly Win Rate Heatmap:")
	fmt.Println()
	
	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		monthReport, err := analytics.GetMonthlyStats(playerName, month.Year(), int(month.Month()))
		
		monthStr := month.Format("Jan 2006")
		fmt.Printf("%-12s ", monthStr)
		
		if err != nil || monthReport.GamesPlayed == 0 {
			fmt.Println("[No Data]")
			continue
		}
		
		winRate := float64(monthReport.Wins) / float64(monthReport.GamesPlayed) * 100
		
		// Create visual bar
		bars := int(winRate / 10) // 10% per bar
		if bars > 10 {
			bars = 10
		}
		
		// Color code based on performance
		if winRate >= 70 {
			green.Print(strings.Repeat("█", bars))
		} else if winRate >= 50 {
			yellow.Print(strings.Repeat("█", bars))
		} else {
			red.Print(strings.Repeat("█", bars))
		}
		
		fmt.Print(strings.Repeat("░", 10-bars))
		fmt.Printf(" %.0f%% (%d games)\n", winRate, monthReport.GamesPlayed)
	}
	
	fmt.Println()
	green.Print("█ ")
	fmt.Print("70%+ | ")
	yellow.Print("█ ")
	fmt.Print("50-69% | ")
	red.Print("█ ")
	fmt.Println("<50%")
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// DisplayComparison shows side-by-side comparison
func DisplayComparison(playerName1, playerName2 string) {
	clear.ClearIt()
	
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	cyan.Println("                  PLAYER COMPARISON")
	cyan.Println(strings.Repeat("=", 70))
	
	// Get stats for both players
	stats1, err1 := database.GetPlayerStats(playerName1)
	stats2, err2 := database.GetPlayerStats(playerName2)
	
	if err1 != nil || err2 != nil {
		color.Red("\nError loading player stats")
		return
	}
	
	fmt.Println()
	yellow.Printf("%-30s %-30s\n", playerName1, playerName2)
	fmt.Println(strings.Repeat("─", 70))
	
	fmt.Printf("Games Played:     %-12d         %-12d\n", stats1.TotalGames, stats2.TotalGames)
	fmt.Printf("Win Rate:         %-12.1f%%       %-12.1f%%\n", stats1.WinRate, stats2.WinRate)
	fmt.Printf("Best Net Worth:   $%-12s       $%-12s\n", 
		FormatCurrency(stats1.BestNetWorth), 
		FormatCurrency(stats2.BestNetWorth))
	fmt.Printf("Avg Net Worth:    $%-12s       $%-12s\n", 
		FormatCurrency(int64(stats1.AverageNetWorth)), 
		FormatCurrency(int64(stats2.AverageNetWorth)))
	fmt.Printf("Best ROI:         %-12.1fx       %-12.1fx\n", stats1.BestROI, stats2.BestROI)
	fmt.Printf("Total Exits:      %-12d         %-12d\n", stats1.TotalExits, stats2.TotalExits)
	
	cyan.Println("\n" + strings.Repeat("=", 70))
	
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// AnalyticsMenu provides options for viewing analytics
func AnalyticsMenu() {
	for {
		clear.ClearIt()
		
		cyan := color.New(color.FgCyan, color.Bold)
		yellow := color.New(color.FgYellow)
		
		cyan.Println("\n" + strings.Repeat("=", 60))
		cyan.Println("                ANALYTICS DASHBOARD")
		cyan.Println(strings.Repeat("=", 60))
		
		yellow.Println("\n1. View Performance Analytics")
		yellow.Println("2. View Performance Heatmap")
		yellow.Println("3. Compare Two Players")
		yellow.Println("4. View Trend Chart")
		yellow.Println("5. Back to Main Menu")
		
		fmt.Print("\nEnter your choice: ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		if choice == "5" {
			return
		}
		
		switch choice {
		case "1":
			fmt.Print("\nEnter your name: ")
			playerName, _ := reader.ReadString('\n')
			playerName = strings.TrimSpace(playerName)
			
			if playerName != "" {
				DisplayAnalyticsDashboard(playerName)
			}
			
		case "2":
			fmt.Print("\nEnter your name: ")
			playerName, _ := reader.ReadString('\n')
			playerName = strings.TrimSpace(playerName)
			
			if playerName != "" {
				DisplayPerformanceHeatmap(playerName)
			}
			
		case "3":
			fmt.Print("\nEnter first player name: ")
			player1, _ := reader.ReadString('\n')
			player1 = strings.TrimSpace(player1)
			
			fmt.Print("Enter second player name: ")
			player2, _ := reader.ReadString('\n')
			player2 = strings.TrimSpace(player2)
			
			if player1 != "" && player2 != "" {
				DisplayComparison(player1, player2)
			}
			
		case "4":
			fmt.Print("\nEnter your name: ")
			playerName, _ := reader.ReadString('\n')
			playerName = strings.TrimSpace(playerName)
			
			if playerName != "" {
				// Get recent games for trend
				scores, err := database.GetTopScoresByPlayer(playerName, 20)
				if err == nil && len(scores) > 0 {
					// Extract ROI values
					roiData := make([]float64, len(scores))
					for i, score := range scores {
						roiData[i] = score.ROI
					}
					
					clear.ClearIt()
					DisplayTrendChart(roiData, "ROI Trend (Last 20 Games)")
					
					fmt.Print("\nPress 'Enter' to continue...")
					bufio.NewReader(os.Stdin).ReadBytes('\n')
				} else {
					color.Red("\nNot enough game data for trend chart")
					time.Sleep(2 * time.Second)
				}
			}
			
		default:
			color.Red("Invalid choice!")
			time.Sleep(1 * time.Second)
		}
	}
}

