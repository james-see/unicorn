package founder

import (
	"fmt"
	"math"
	"math/rand"
)

// LaunchContentProgram starts a content marketing and SEO initiative
func (fs *FounderState) LaunchContentProgram(monthlyBudget int64) error {
	if monthlyBudget < 10000 {
		return fmt.Errorf("minimum budget is $10,000/month")
	}

	if monthlyBudget > 50000 {
		monthlyBudget = 50000 // Cap at $50k/month
	}

	if fs.Cash < monthlyBudget {
		return fmt.Errorf("insufficient cash for initial month")
	}

	// Pay first month
	fs.Cash -= monthlyBudget

	fs.ContentProgram = &ContentProgram{
		MonthlyBudget:   monthlyBudget,
		ContentTypes:    make(map[string]bool),
		OrganicTraffic:  100, // Start with 100 monthly visitors
		InboundLeads:    0,
		ContentQuality:  0.50, // Starts at 50% quality
		SEOScore:        10,   // Starts at 10/100
		MonthsActive:    0,
		TotalInvestment: monthlyBudget,
		CumulativeLeads: 0,
		LaunchedMonth:   fs.Turn,
	}

	// Default content types based on budget
	fs.ContentProgram.ContentTypes["blog"] = true
	fs.ContentProgram.ContentTypes["seo"] = true
	if monthlyBudget >= 30000 {
		fs.ContentProgram.ContentTypes["webinars"] = true
	}

	return nil
}

// UpdateContentProgram processes monthly content marketing growth
func (fs *FounderState) UpdateContentProgram() []string {
	var messages []string

	if fs.ContentProgram == nil {
		return messages
	}

	// Pay monthly budget
	if fs.Cash < fs.ContentProgram.MonthlyBudget {
		// Can't afford content marketing - pause it
		messages = append(messages, "⚠️  Content marketing paused due to insufficient cash")
		return messages
	}

	fs.Cash -= fs.ContentProgram.MonthlyBudget
	fs.ContentProgram.MonthsActive++
	fs.ContentProgram.TotalInvestment += fs.ContentProgram.MonthlyBudget

	// Quality improves over time (up to 0.90)
	if fs.ContentProgram.ContentQuality < 0.90 {
		qualityIncrease := 0.05
		if fs.ContentProgram.MonthlyBudget >= 30000 {
			qualityIncrease = 0.08
		}
		fs.ContentProgram.ContentQuality = math.Min(0.90, fs.ContentProgram.ContentQuality+qualityIncrease)
	}

	// SEO score improves slowly (takes 3-6 months to see results)
	if fs.ContentProgram.MonthsActive >= 3 && fs.ContentProgram.SEOScore < 90 {
		seoIncrease := 5 + rand.Intn(5) // 5-10 points per month
		if fs.ContentProgram.MonthlyBudget >= 40000 {
			seoIncrease += 5
		}
		fs.ContentProgram.SEOScore = int(math.Min(90, float64(fs.ContentProgram.SEOScore+seoIncrease)))
	}

	// Organic traffic compounds monthly
	growthRate := 0.05 + (fs.ContentProgram.ContentQuality * 0.10) // 5-15% growth
	if fs.ContentProgram.SEOScore > 50 {
		growthRate += 0.05 // Additional 5% if SEO is good
	}

	fs.ContentProgram.OrganicTraffic = int(float64(fs.ContentProgram.OrganicTraffic) * (1.0 + growthRate))

	// Convert traffic to inbound leads (1-3% conversion)
	conversionRate := 0.01 + (fs.ContentProgram.ContentQuality * 0.02)
	fs.ContentProgram.InboundLeads = int(float64(fs.ContentProgram.OrganicTraffic) * conversionRate)
	fs.ContentProgram.CumulativeLeads += fs.ContentProgram.InboundLeads

	// Add message
	if fs.ContentProgram.MonthsActive == 3 {
		messages = append(messages, "📈 Content marketing starting to show SEO results!")
	}

	if fs.ContentProgram.SEOScore >= 80 {
		messages = append(messages, fmt.Sprintf("🎯 Strong SEO! %d inbound leads this month (40%% lower CAC)", fs.ContentProgram.InboundLeads))
	} else if fs.ContentProgram.InboundLeads > 0 {
		messages = append(messages, fmt.Sprintf("📝 Content generated %d inbound leads this month", fs.ContentProgram.InboundLeads))
	}

	return messages
}

// EndContentProgram stops content marketing
func (fs *FounderState) EndContentProgram() {
	fs.ContentProgram = nil
}

