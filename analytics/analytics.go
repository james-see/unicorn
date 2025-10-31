package analytics

import (
	"fmt"
	"sort"
	"strings"

	game "github.com/jamesacampbell/unicorn/game"
)

// PortfolioAnalytics tracks detailed portfolio performance
type PortfolioAnalytics struct {
	TotalInvested       int64
	CurrentValue        int64
	TotalGainLoss       int64
	PercentageChange    float64
	BestPerformer       InvestmentPerformance
	WorstPerformer      InvestmentPerformance
	SectorBreakdown     map[string]SectorPerformance
	InvestmentCount     int
	AverageReturn       float64
	PositiveInvestments int
	NegativeInvestments int
}

// InvestmentPerformance tracks individual investment metrics
type InvestmentPerformance struct {
	CompanyName      string
	Category         string
	AmountInvested   int64
	CurrentValue     int64
	GainLoss         int64
	PercentageChange float64
}

// SectorPerformance tracks performance by sector
type SectorPerformance struct {
	SectorName      string
	TotalInvested   int64
	CurrentValue    int64
	GainLoss        int64
	ROI             float64
	CompanyCount    int
	AverageROI      float64
}

// CalculateAnalytics generates comprehensive portfolio analytics
func CalculateAnalytics(gs *game.GameState) *PortfolioAnalytics {
	analytics := &PortfolioAnalytics{
		SectorBreakdown: make(map[string]SectorPerformance),
	}
	
	var performances []InvestmentPerformance
	sectorData := make(map[string]*SectorPerformance)
	
	// Analyze each investment
	for _, inv := range gs.Portfolio.Investments {
		value := int64((inv.EquityPercent / 100.0) * float64(inv.CurrentValuation))
		gainLoss := value - inv.AmountInvested
		percentChange := (float64(gainLoss) / float64(inv.AmountInvested)) * 100.0
		
		perf := InvestmentPerformance{
			CompanyName:      inv.CompanyName,
			AmountInvested:   inv.AmountInvested,
			CurrentValue:     value,
			GainLoss:         gainLoss,
			PercentageChange: percentChange,
		}
		
		// Find the company to get its category
		for _, startup := range gs.AvailableStartups {
			if startup.Name == inv.CompanyName {
				perf.Category = startup.Category
				break
			}
		}
		
		performances = append(performances, perf)
		
		analytics.TotalInvested += inv.AmountInvested
		analytics.CurrentValue += value
		analytics.InvestmentCount++
		
		if gainLoss > 0 {
			analytics.PositiveInvestments++
		} else if gainLoss < 0 {
			analytics.NegativeInvestments++
		}
		
		// Track sector performance
		if perf.Category != "" {
			if _, exists := sectorData[perf.Category]; !exists {
				sectorData[perf.Category] = &SectorPerformance{
					SectorName: perf.Category,
				}
			}
			
			sector := sectorData[perf.Category]
			sector.TotalInvested += inv.AmountInvested
			sector.CurrentValue += value
			sector.GainLoss += gainLoss
			sector.CompanyCount++
		}
	}
	
	// Calculate totals
	analytics.TotalGainLoss = analytics.CurrentValue - analytics.TotalInvested
	if analytics.TotalInvested > 0 {
		analytics.PercentageChange = (float64(analytics.TotalGainLoss) / float64(analytics.TotalInvested)) * 100.0
		analytics.AverageReturn = analytics.PercentageChange / float64(analytics.InvestmentCount)
	}
	
	// Find best and worst performers
	if len(performances) > 0 {
		sort.Slice(performances, func(i, j int) bool {
			return performances[i].PercentageChange > performances[j].PercentageChange
		})
		analytics.BestPerformer = performances[0]
		analytics.WorstPerformer = performances[len(performances)-1]
	}
	
	// Calculate sector ROI
	for sectorName, sector := range sectorData {
		if sector.TotalInvested > 0 {
			sector.ROI = (float64(sector.GainLoss) / float64(sector.TotalInvested)) * 100.0
			sector.AverageROI = sector.ROI / float64(sector.CompanyCount)
		}
		analytics.SectorBreakdown[sectorName] = *sector
	}
	
	return analytics
}

// GetTopSectors returns sectors sorted by performance
func GetTopSectors(analytics *PortfolioAnalytics) []SectorPerformance {
	var sectors []SectorPerformance
	for _, sector := range analytics.SectorBreakdown {
		sectors = append(sectors, sector)
	}
	
	sort.Slice(sectors, func(i, j int) bool {
		return sectors[i].ROI > sectors[j].ROI
	})
	
	return sectors
}

// GenerateASCIIChart creates a simple ASCII bar chart
func GenerateASCIIChart(values map[string]float64, width int) string {
	if len(values) == 0 {
		return "No data"
	}
	
	// Find max value for scaling
	maxVal := 0.0
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	
	if maxVal == 0 {
		maxVal = 1 // Avoid division by zero
	}
	
	// Sort keys
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var chart strings.Builder
	for _, key := range keys {
		value := values[key]
		barLength := int((value / maxVal) * float64(width))
		if barLength < 0 {
			barLength = 0
		}
		
		bar := strings.Repeat("=", barLength)
		chart.WriteString(fmt.Sprintf("%-15s %s %.1f%%\n", key, bar, value))
	}
	
	return chart.String()
}

// GetPortfolioSummary returns a formatted summary
func GetPortfolioSummary(analytics *PortfolioAnalytics) string {
	var summary strings.Builder
	
	summary.WriteString("\n?? PORTFOLIO ANALYTICS\n")
	summary.WriteString(strings.Repeat("=", 50) + "\n\n")
	
	summary.WriteString(fmt.Sprintf("Total Invested:     $%s\n", formatMoney(analytics.TotalInvested)))
	summary.WriteString(fmt.Sprintf("Current Value:      $%s\n", formatMoney(analytics.CurrentValue)))
	
	if analytics.TotalGainLoss >= 0 {
		summary.WriteString(fmt.Sprintf("Total Gain:         +$%s (%.1f%%)\n", 
			formatMoney(analytics.TotalGainLoss), analytics.PercentageChange))
	} else {
		summary.WriteString(fmt.Sprintf("Total Loss:         -$%s (%.1f%%)\n", 
			formatMoney(-analytics.TotalGainLoss), analytics.PercentageChange))
	}
	
	summary.WriteString(fmt.Sprintf("\nInvestments:        %d total\n", analytics.InvestmentCount))
	summary.WriteString(fmt.Sprintf("Positive:           %d (%.0f%%)\n", 
		analytics.PositiveInvestments, 
		float64(analytics.PositiveInvestments)/float64(analytics.InvestmentCount)*100))
	summary.WriteString(fmt.Sprintf("Negative:           %d (%.0f%%)\n", 
		analytics.NegativeInvestments,
		float64(analytics.NegativeInvestments)/float64(analytics.InvestmentCount)*100))
	
	if analytics.BestPerformer.CompanyName != "" {
		summary.WriteString(fmt.Sprintf("\n?? Best:             %s (%.1f%%)\n", 
			analytics.BestPerformer.CompanyName, analytics.BestPerformer.PercentageChange))
	}
	
	if analytics.WorstPerformer.CompanyName != "" {
		summary.WriteString(fmt.Sprintf("?? Worst:            %s (%.1f%%)\n", 
			analytics.WorstPerformer.CompanyName, analytics.WorstPerformer.PercentageChange))
	}
	
	return summary.String()
}

// GetSectorBreakdown returns formatted sector analysis
func GetSectorBreakdown(analytics *PortfolioAnalytics) string {
	var breakdown strings.Builder
	
	breakdown.WriteString("\n?? SECTOR BREAKDOWN\n")
	breakdown.WriteString(strings.Repeat("=", 50) + "\n\n")
	
	sectors := GetTopSectors(analytics)
	
	if len(sectors) == 0 {
		breakdown.WriteString("No sector data available\n")
		return breakdown.String()
	}
	
	for i, sector := range sectors {
		breakdown.WriteString(fmt.Sprintf("%d. %s\n", i+1, sector.SectorName))
		breakdown.WriteString(fmt.Sprintf("   Companies: %d | ROI: %.1f%%\n", 
			sector.CompanyCount, sector.ROI))
		breakdown.WriteString(fmt.Sprintf("   Invested: $%s ? Value: $%s\n\n", 
			formatMoney(sector.TotalInvested), formatMoney(sector.CurrentValue)))
	}
	
	return breakdown.String()
}

// Helper function to format money
func formatMoney(amount int64) string {
	abs := amount
	if abs < 0 {
		abs = -abs
	}
	
	s := fmt.Sprintf("%d", abs)
	
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
