package game

import (
	"math"
)

// VCReputation represents a player's reputation in the VC ecosystem
type VCReputation struct {
	PlayerName       string
	PerformanceScore float64 // 0-100, based on past fund returns
	FounderScore     float64 // 0-100, based on founder relationship quality
	MarketScore      float64 // 0-100, based on win streaks and achievements
	TotalGamesPlayed int
	SuccessfulExits  int
	AvgROILast5      float64 // Average ROI from last 5 games
	LastUpdated      string
}

// GetAggregateReputation calculates the overall reputation score (0-100)
// Weighted: Performance 40%, Founder 30%, Market 30%
func (r *VCReputation) GetAggregateReputation() float64 {
	return (r.PerformanceScore * 0.4) + (r.FounderScore * 0.3) + (r.MarketScore * 0.3)
}

// GetDealQualityTier returns the deal quality tier based on aggregate reputation
func (r *VCReputation) GetDealQualityTier() string {
	aggregate := r.GetAggregateReputation()
	if aggregate >= 70 {
		return "Tier 1 (Hot Deals)"
	} else if aggregate >= 40 {
		return "Tier 2 (Standard Deals)"
	}
	return "Tier 3 (Struggling Deals)"
}

// GetDealQualityDescription returns a description of what the tier means
func (r *VCReputation) GetDealQualityDescription() string {
	aggregate := r.GetAggregateReputation()
	if aggregate >= 70 {
		return "Access to high-quality startups with lower risk and higher growth potential"
	} else if aggregate >= 40 {
		return "Access to standard startup opportunities with balanced risk/reward"
	}
	return "Limited to higher-risk startups with lower growth potential"
}

// CalculatePerformanceScore calculates performance score from recent game history
// Based on average ROI from last 5 games
func CalculatePerformanceScore(avgROI float64, successfulExits int, totalGames int) float64 {
	if totalGames == 0 {
		return 50.0 // Starting score
	}

	// Base score from ROI (0-80 points)
	// ROI of 0% = 40 points, 100% = 70 points, 200%+ = 80 points
	roiScore := 40.0
	if avgROI > 0 {
		roiScore = 40.0 + math.Min(avgROI/2.5, 40.0) // Max 80 from ROI
	} else if avgROI < 0 {
		// Negative ROI reduces score
		roiScore = math.Max(0, 40.0+(avgROI/2.0)) // Can go down to 0
	}

	// Bonus from successful exits (0-20 points)
	exitRate := float64(successfulExits) / float64(totalGames)
	exitBonus := exitRate * 20.0

	score := roiScore + exitBonus

	// Cap at 100
	return math.Min(score, 100.0)
}

// CalculateMarketScore calculates market score from achievements and win streaks
func CalculateMarketScore(achievementPoints int, winStreak int) float64 {
	// Base score from achievement points (0-70 points)
	// 0 points = 30, 500 points = 50, 1000+ points = 70
	achievementScore := 30.0 + math.Min(float64(achievementPoints)/25.0, 40.0)

	// Win streak bonus (0-30 points)
	// 1 win = 5 points, 2 wins = 10, 3 wins = 15, 4+ wins = 30
	streakBonus := math.Min(float64(winStreak)*7.5, 30.0)

	score := achievementScore + streakBonus

	// Cap at 100
	return math.Min(score, 100.0)
}

// UpdateReputationAfterGame updates reputation based on game performance
func UpdateReputationAfterGame(current *VCReputation, gameROI float64, hadSuccessfulExit bool, achievementPoints int, winStreak int) *VCReputation {
	updated := &VCReputation{
		PlayerName:       current.PlayerName,
		TotalGamesPlayed: current.TotalGamesPlayed + 1,
		SuccessfulExits:  current.SuccessfulExits,
		LastUpdated:      current.LastUpdated,
	}

	if hadSuccessfulExit {
		updated.SuccessfulExits++
	}

	// Update AvgROILast5 (rolling average of last 5 games)
	// This is simplified - we'd need to track individual game ROIs for true rolling average
	// For now, we'll use weighted average where new game is 20% weight
	if current.TotalGamesPlayed == 0 {
		updated.AvgROILast5 = gameROI
	} else if current.TotalGamesPlayed < 5 {
		// Building up to 5 games
		weight := 1.0 / float64(current.TotalGamesPlayed+1)
		updated.AvgROILast5 = (current.AvgROILast5 * (1 - weight)) + (gameROI * weight)
	} else {
		// Rolling average with 20% weight on new game
		updated.AvgROILast5 = (current.AvgROILast5 * 0.8) + (gameROI * 0.2)
	}

	// Recalculate scores
	updated.PerformanceScore = CalculatePerformanceScore(updated.AvgROILast5, updated.SuccessfulExits, updated.TotalGamesPlayed)
	updated.MarketScore = CalculateMarketScore(achievementPoints, winStreak)

	// Founder score carries over and will be updated separately
	updated.FounderScore = current.FounderScore

	return updated
}

// UpdateFounderScore updates the founder reputation component
func (r *VCReputation) UpdateFounderScore(averageRelationshipScore float64) {
	// Founder score is a weighted average
	// New average has 30% weight, existing score has 70% weight (builds slowly)
	if r.TotalGamesPlayed == 0 {
		r.FounderScore = averageRelationshipScore
	} else {
		r.FounderScore = (r.FounderScore * 0.7) + (averageRelationshipScore * 0.3)
	}
}

// GetNewPlayerReputation returns default reputation for new players
func GetNewPlayerReputation(playerName string) *VCReputation {
	return &VCReputation{
		PlayerName:       playerName,
		PerformanceScore: 50.0, // Neutral starting score
		FounderScore:     50.0, // Neutral starting score
		MarketScore:      50.0, // Neutral starting score
		TotalGamesPlayed: 0,
		SuccessfulExits:  0,
		AvgROILast5:      0.0,
	}
}

// GetReputationLevel returns a descriptive level based on aggregate score
func (r *VCReputation) GetReputationLevel() string {
	aggregate := r.GetAggregateReputation()
	if aggregate >= 90 {
		return "Legendary VC"
	} else if aggregate >= 80 {
		return "Top-Tier VC"
	} else if aggregate >= 70 {
		return "Established VC"
	} else if aggregate >= 60 {
		return "Rising VC"
	} else if aggregate >= 50 {
		return "Competent VC"
	} else if aggregate >= 40 {
		return "Developing VC"
	}
	return "Emerging VC"
}

