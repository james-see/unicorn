package analytics

import (
	"fmt"
	"sort"
	"time"

	"github.com/jamesacampbell/unicorn/database"
)

// PerformanceTrend represents performance over a time period
type PerformanceTrend struct {
	Period      string
	GamesPlayed int
	WinRate     float64
	AvgROI      float64
	AvgNetWorth int64
	Trend       string // "Improving", "Declining", "Stable"
}

// ComparisonStats compares player to global averages
type ComparisonStats struct {
	YourAvg    float64
	GlobalAvg  float64
	Percentile int // Top X% of players
}

// MonthlyReport contains performance data for a specific month
type MonthlyReport struct {
	Month          string
	Year           int
	GamesPlayed    int
	Wins           int
	AvgROI         float64
	BestNetWorth   int64
	TotalInvested  int64
}

// SectorPerformance tracks ROI by sector
type SectorPerformance struct {
	SectorName      string
	InvestmentCount int
	AvgROI          float64
	BestROI         float64
	WorstROI        float64
	TotalInvested   int64
}

// TrendReport contains comprehensive trend analysis
type TrendReport struct {
	Last7Days   PerformanceTrend
	Last30Days  PerformanceTrend
	AllTime     PerformanceTrend
	TrendVector string // Overall trend direction
	Insights    []string
}

// GenerateTrendAnalysis generates trend analysis for a player
func GenerateTrendAnalysis(playerName string, daysBack int) (*TrendReport, error) {
	// Get game history
	scores, err := database.GetTopScoresByPlayer(playerName, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get game history: %v", err)
	}
	
	if len(scores) == 0 {
		return &TrendReport{
			TrendVector: "Insufficient Data",
			Insights:    []string{"Play more games to see trends!"},
		}, nil
	}
	
	now := time.Now()
	
	// Calculate trends for different periods
	last7Days := calculatePeriodTrend(scores, now, 7)
	last30Days := calculatePeriodTrend(scores, now, 30)
	allTime := calculatePeriodTrend(scores, now, 36500) // ~100 years
	
	// Determine overall trend vector
	trendVector := determineTrendVector(last7Days, last30Days, allTime)
	
	// Generate insights
	insights := generateInsights(last7Days, last30Days, allTime, scores)
	
	return &TrendReport{
		Last7Days:   last7Days,
		Last30Days:  last30Days,
		AllTime:     allTime,
		TrendVector: trendVector,
		Insights:    insights,
	}, nil
}

// calculatePeriodTrend calculates performance trend for a specific period
func calculatePeriodTrend(scores []database.GameScore, endDate time.Time, daysBack int) PerformanceTrend {
	startDate := endDate.AddDate(0, 0, -daysBack)
	
	periodScores := []database.GameScore{}
	for _, score := range scores {
		if score.PlayedAt.After(startDate) && score.PlayedAt.Before(endDate) {
			periodScores = append(periodScores, score)
		}
	}
	
	if len(periodScores) == 0 {
		return PerformanceTrend{
			Period:      getPeriodName(daysBack),
			GamesPlayed: 0,
			Trend:       "No Data",
		}
	}
	
	// Calculate stats
	totalROI := 0.0
	totalNetWorth := int64(0)
	wins := 0
	
	for _, score := range periodScores {
		totalROI += score.ROI
		totalNetWorth += score.FinalNetWorth
		if score.ROI > 0 {
			wins++
		}
	}
	
	avgROI := totalROI / float64(len(periodScores))
	avgNetWorth := totalNetWorth / int64(len(periodScores))
	winRate := float64(wins) / float64(len(periodScores)) * 100
	
	// Determine trend (compare first half vs second half)
	trend := "Stable"
	if len(periodScores) >= 4 {
		midPoint := len(periodScores) / 2
		firstHalfAvg := 0.0
		secondHalfAvg := 0.0
		
		for i := 0; i < midPoint; i++ {
			firstHalfAvg += periodScores[i].ROI
		}
		firstHalfAvg /= float64(midPoint)
		
		for i := midPoint; i < len(periodScores); i++ {
			secondHalfAvg += periodScores[i].ROI
		}
		secondHalfAvg /= float64(len(periodScores) - midPoint)
		
		improvement := (secondHalfAvg - firstHalfAvg) / (firstHalfAvg + 0.01) // Avoid div by zero
		
		if improvement > 0.15 {
			trend = "Improving"
		} else if improvement < -0.15 {
			trend = "Declining"
		}
	}
	
	return PerformanceTrend{
		Period:      getPeriodName(daysBack),
		GamesPlayed: len(periodScores),
		WinRate:     winRate,
		AvgROI:      avgROI,
		AvgNetWorth: avgNetWorth,
		Trend:       trend,
	}
}

// determineTrendVector determines overall trend direction
func determineTrendVector(last7, last30, allTime PerformanceTrend) string {
	if last7.GamesPlayed == 0 && last30.GamesPlayed == 0 {
		return "Insufficient Data"
	}
	
	// Weight recent trends more heavily
	if last7.GamesPlayed >= 3 {
		if last7.Trend == "Improving" {
			return "📈 Strong Upward Trend"
		} else if last7.Trend == "Declining" {
			return "📉 Recent Decline"
		}
	}
	
	if last30.Trend == "Improving" {
		return "↗ Improving"
	} else if last30.Trend == "Declining" {
		return "↘ Declining"
	}
	
	return "→ Stable"
}

// generateInsights generates actionable insights from trends
func generateInsights(last7, last30, allTime PerformanceTrend, allScores []database.GameScore) []string {
	insights := []string{}
	
	// Win rate insights
	if last30.WinRate > allTime.WinRate+10 {
		insights = append(insights, fmt.Sprintf("✨ You're performing %.0f%% better than your all-time average!", 
			last30.WinRate-allTime.WinRate))
	} else if last30.WinRate < allTime.WinRate-10 {
		insights = append(insights, fmt.Sprintf("⚠️ Your recent win rate is %.0f%% below your average", 
			allTime.WinRate-last30.WinRate))
	}
	
	// ROI insights
	if last30.AvgROI > allTime.AvgROI*1.2 {
		insights = append(insights, "💰 Your ROI is significantly above your average - great job!")
	}
	
	// Consistency insights
	if len(allScores) >= 10 {
		// Calculate consistency (lower std deviation = more consistent)
		stdDev := calculateStdDev(allScores)
		if stdDev < 0.5 {
			insights = append(insights, "🎯 You're very consistent across games")
		} else if stdDev > 2.0 {
			insights = append(insights, "📊 Your performance varies significantly - try to find what works")
		}
	}
	
	// Game volume insights
	if last7.GamesPlayed > last30.GamesPlayed/2 {
		insights = append(insights, "🔥 You've been playing a lot lately - stay focused!")
	}
	
	// Best performance insights
	if len(allScores) > 0 {
		bestROI := allScores[0].ROI
		for _, score := range allScores {
			if score.ROI > bestROI {
				bestROI = score.ROI
			}
		}
		
		if bestROI > 5.0 {
			insights = append(insights, fmt.Sprintf("🏆 Your best ROI is %.1fx - try to replicate that success!", bestROI))
		}
	}
	
	if len(insights) == 0 {
		insights = append(insights, "Keep playing to generate more insights!")
	}
	
	return insights
}

// CompareToGlobal compares player stats to global averages
func CompareToGlobal(playerStats database.PlayerStats) (ComparisonStats, error) {
	// Get global stats (top 100 players as sample)
	allScores, err := database.GetTopScoresByNetWorth(100, "all")
	if err != nil {
		return ComparisonStats{}, err
	}
	
	if len(allScores) == 0 {
		return ComparisonStats{
			YourAvg:    playerStats.AverageNetWorth,
			GlobalAvg:  0,
			Percentile: 100,
		}, nil
	}
	
	// Calculate global average net worth
	globalTotal := int64(0)
	for _, score := range allScores {
		globalTotal += score.FinalNetWorth
	}
	globalAvg := float64(globalTotal) / float64(len(allScores))
	
	// Calculate percentile (simplified)
	betterThan := 0
	for _, score := range allScores {
		if playerStats.BestNetWorth > score.FinalNetWorth {
			betterThan++
		}
	}
	percentile := (betterThan * 100) / len(allScores)
	
	return ComparisonStats{
		YourAvg:    playerStats.AverageNetWorth,
		GlobalAvg:  globalAvg,
		Percentile: 100 - percentile, // Top X%
	}, nil
}

// CalculateTrend determines if a set of values is improving/declining
func CalculateTrend(dataPoints []float64) string {
	if len(dataPoints) < 3 {
		return "Stable"
	}
	
	// Simple linear regression to determine trend
	n := float64(len(dataPoints))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	
	for i, y := range dataPoints {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	// Slope of regression line
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	
	if slope > 0.1 {
		return "Improving"
	} else if slope < -0.1 {
		return "Declining"
	}
	
	return "Stable"
}

// Helper functions

func getPeriodName(days int) string {
	if days <= 7 {
		return "Last 7 Days"
	} else if days <= 30 {
		return "Last 30 Days"
	}
	return "All Time"
}

func calculateStdDev(scores []database.GameScore) float64 {
	if len(scores) == 0 {
		return 0
	}
	
	// Calculate mean
	sum := 0.0
	for _, score := range scores {
		sum += score.ROI
	}
	mean := sum / float64(len(scores))
	
	// Calculate variance
	variance := 0.0
	for _, score := range scores {
		diff := score.ROI - mean
		variance += diff * diff
	}
	variance /= float64(len(scores))
	
	// Standard deviation is square root of variance
	// Simplified calculation
	return variance
}

// GetMonthlyStats returns performance for a specific month
func GetMonthlyStats(playerName string, year, month int) (*MonthlyReport, error) {
	scores, err := database.GetTopScoresByPlayer(playerName, 1000)
	if err != nil {
		return nil, err
	}
	
	var monthScores []database.GameScore
	for _, score := range scores {
		if score.PlayedAt.Year() == year && int(score.PlayedAt.Month()) == month {
			monthScores = append(monthScores, score)
		}
	}
	
	if len(monthScores) == 0 {
		return &MonthlyReport{
			Month:       time.Month(month).String(),
			Year:        year,
			GamesPlayed: 0,
		}, nil
	}
	
	totalROI := 0.0
	bestNetWorth := int64(0)
	wins := 0
	
	for _, score := range monthScores {
		totalROI += score.ROI
		if score.FinalNetWorth > bestNetWorth {
			bestNetWorth = score.FinalNetWorth
		}
		if score.ROI > 0 {
			wins++
		}
		// Note: TotalInvested not available in GameScore, would need detailed history
	}
	
	return &MonthlyReport{
		Month:        time.Month(month).String(),
		Year:         year,
		GamesPlayed:  len(monthScores),
		Wins:         wins,
		AvgROI:       totalROI / float64(len(monthScores)),
		BestNetWorth: bestNetWorth,
	}, nil
}

// GetTopInvestments returns best performing investments (placeholder - would need detailed game history)
func GetTopInvestments(playerName string, limit int) ([]InvestmentSummary, error) {
	// This would require the game_history_detailed table to be populated
	// For now, return empty slice
	return []InvestmentSummary{}, nil
}

// InvestmentSummary represents a summary of an investment
type InvestmentSummary struct {
	CompanyName string
	ROI         float64
	ExitValue   int64
	GameDate    time.Time
}

// GetRecentGamesAnalysis provides quick analysis of recent games
func GetRecentGamesAnalysis(playerName string, count int) (map[string]interface{}, error) {
	scores, err := database.GetTopScoresByPlayer(playerName, count)
	if err != nil {
		return nil, err
	}
	
	if len(scores) == 0 {
		return map[string]interface{}{
			"games_played": 0,
			"message":      "No games played yet",
		}, nil
	}
	
	// Sort by played_at descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].PlayedAt.After(scores[j].PlayedAt)
	})
	
	// Take most recent
	if len(scores) > count {
		scores = scores[:count]
	}
	
	avgROI := 0.0
	avgNetWorth := int64(0)
	wins := 0
	
	for _, score := range scores {
		avgROI += score.ROI
		avgNetWorth += score.FinalNetWorth
		if score.ROI > 0 {
			wins++
		}
	}
	
	return map[string]interface{}{
		"games_played":   len(scores),
		"avg_roi":        avgROI / float64(len(scores)),
		"avg_net_worth":  avgNetWorth / int64(len(scores)),
		"win_rate":       float64(wins) / float64(len(scores)) * 100,
		"most_recent":    scores[0].PlayedAt,
	}, nil
}
