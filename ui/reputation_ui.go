package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesacampbell/unicorn/game"
)

// DisplayReputation shows player's VC reputation
func DisplayReputation(rep *game.VCReputation) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta, color.Bold)

	fmt.Println()
	cyan.Println(strings.Repeat("=", 70))
	cyan.Println("                    YOUR VC REPUTATION")
	cyan.Println(strings.Repeat("=", 70))

	// Aggregate score
	aggregate := (rep.PerformanceScore * 0.4) + (rep.FounderScore * 0.3) + (rep.MarketScore * 0.3)
	level := getReputationLevel(aggregate)

	fmt.Println()
	magenta.Printf("Overall Reputation: %.1f/100 - %s\n", aggregate, level)
	fmt.Println()

	// Component scores
	yellow.Println("REPUTATION COMPONENTS:")
	fmt.Printf("  Performance Score:  %.1f/100 ", rep.PerformanceScore)
	printScoreBar(rep.PerformanceScore)
	fmt.Printf("  Founder Score:      %.1f/100 ", rep.FounderScore)
	printScoreBar(rep.FounderScore)
	fmt.Printf("  Market Score:       %.1f/100 ", rep.MarketScore)
	printScoreBar(rep.MarketScore)

	fmt.Println()
	yellow.Println("CAREER STATS:")
	fmt.Printf("  Games Played:       %d\n", rep.TotalGamesPlayed)
	fmt.Printf("  Successful Exits:   %d\n", rep.SuccessfulExits)
	fmt.Printf("  Avg ROI (Last 5):   %.1f%%\n", rep.AvgROILast5)

	// Deal flow quality
	fmt.Println()
	tier := getDealQualityTier(aggregate)
	tierDesc := getDealQualityDescription(aggregate)

	if aggregate >= 70 {
		green.Printf("DEAL FLOW: %s\n", tier)
	} else if aggregate >= 40 {
		yellow.Printf("DEAL FLOW: %s\n", tier)
	} else {
		color.Red("DEAL FLOW: %s\n", tier)
	}
	fmt.Printf("  %s\n", tierDesc)

	fmt.Println()
	cyan.Println(strings.Repeat("=", 70))
}

func printScoreBar(score float64) {
	bars := int(score / 10)
	fmt.Print("(")
	for i := 0; i < 10; i++ {
		if i < bars {
			color.Green("█")
		} else {
			fmt.Print("░")
		}
	}
	fmt.Println(")")
}

func getReputationLevel(score float64) string {
	if score >= 90 {
		return "Legendary VC"
	} else if score >= 80 {
		return "Top-Tier VC"
	} else if score >= 70 {
		return "Established VC"
	} else if score >= 60 {
		return "Rising VC"
	} else if score >= 50 {
		return "Competent VC"
	} else if score >= 40 {
		return "Developing VC"
	}
	return "Emerging VC"
}

func getDealQualityTier(score float64) string {
	if score >= 70 {
		return "Tier 1 (Hot Deals)"
	} else if score >= 40 {
		return "Tier 2 (Standard Deals)"
	}
	return "Tier 3 (Struggling Deals)"
}

func getDealQualityDescription(score float64) string {
	if score >= 70 {
		return "Access to high-quality startups with lower risk and higher growth potential"
	} else if score >= 40 {
		return "Access to standard startup opportunities with balanced risk/reward"
	}
	return "Limited to higher-risk startups with lower growth potential"
}

// DisplayReputationSummary shows a brief reputation summary
func DisplayReputationSummary(rep *game.VCReputation) {
	aggregate := (rep.PerformanceScore * 0.4) + (rep.FounderScore * 0.3) + (rep.MarketScore * 0.3)
	level := getReputationLevel(aggregate)
	tier := getDealQualityTier(aggregate)

	cyan := color.New(color.FgCyan)
	cyan.Printf("\nReputation: %.1f/100 (%s) | Deal Flow: %s\n", aggregate, level, tier)
}
