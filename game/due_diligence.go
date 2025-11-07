package game

import (
	"fmt"
	"math/rand"
)

// DDDecision represents a due diligence opportunity before investment
type DDDecision struct {
	CompanyName      string
	StartupIndex     int
	InvestmentAmount int64
	Turn             int
}

// DDLevel represents due diligence depth
type DDLevel struct {
	ID          string
	Name        string
	Description string
	Cost        int64
	Duration    int // Days (not turns, just flavor)
	Reveals     []string
}

// DDFinding represents something discovered during due diligence
type DDFinding struct {
	Type        string // "red_flag", "green_flag", "neutral"
	Category    string // "founder", "tech", "legal", "financial", "market"
	Description string
	Impact      float64 // Impact on valuation trajectory (-0.2 to +0.2)
	RiskImpact  float64 // Impact on risk score (-0.1 to +0.1)
}

// GetDDLevels returns available due diligence options
func GetDDLevels() []DDLevel {
	return []DDLevel{
		{
			ID:          "none",
			Name:        "No Due Diligence",
			Description: "Invest immediately without additional research",
			Cost:        0,
			Duration:    0,
			Reveals:     []string{},
		},
		{
			ID:          "quick",
			Name:        "Quick DD",
			Description: "Basic review of metrics and team",
			Cost:        5000,
			Duration:    1,
			Reveals:     []string{"risk_score", "growth_potential"},
		},
		{
			ID:          "standard",
			Name:        "Standard DD",
			Description: "Thorough review including founder references",
			Cost:        15000,
			Duration:    3,
			Reveals:     []string{"risk_score", "growth_potential", "founder_quality", "2_hidden_metrics"},
		},
		{
			ID:          "deep",
			Name:        "Deep DD",
			Description: "Comprehensive analysis with technical audit",
			Cost:        30000,
			Duration:    7,
			Reveals:     []string{"full_disclosure", "risk_reduction"},
		},
	}
}

// PerformDueDiligence executes due diligence and generates findings
func PerformDueDiligence(startup *Startup, level string) []DDFinding {
	findings := []DDFinding{}

	if level == "none" {
		return findings
	}

	// Generate findings based on DD level and company characteristics

	// Quick DD: Just basic metrics
	if level == "quick" {
		findings = append(findings, DDFinding{
			Type:     "neutral",
			Category: "financial",
			Description: fmt.Sprintf("Risk Score: %.1f/1.0, Growth Potential: %.1f/1.0",
				startup.RiskScore, startup.GrowthPotential),
			Impact:     0.0,
			RiskImpact: 0.0,
		})
	}

	// Standard DD: More detailed findings
	if level == "standard" || level == "deep" {
		// Founder quality check
		founderRoll := rand.Float64()
		if founderRoll < 0.15 { // 15% chance of red flag
			findings = append(findings, DDFinding{
				Type:        "red_flag",
				Category:    "founder",
				Description: "Previous startup failure not disclosed; co-founder conflict rumors",
				Impact:      -0.10,
				RiskImpact:  0.05,
			})
		} else if founderRoll > 0.85 { // 15% chance of green flag
			findings = append(findings, DDFinding{
				Type:        "green_flag",
				Category:    "founder",
				Description: "Founder has successful exit history; strong industry reputation",
				Impact:      0.08,
				RiskImpact:  -0.03,
			})
		} else {
			findings = append(findings, DDFinding{
				Type:        "neutral",
				Category:    "founder",
				Description: "Founder has solid background with relevant industry experience",
				Impact:      0.0,
				RiskImpact:  0.0,
			})
		}

		// Financial metrics check
		financialRoll := rand.Float64()
		if financialRoll < 0.20 {
			findings = append(findings, DDFinding{
				Type:        "red_flag",
				Category:    "financial",
				Description: "Burn rate higher than disclosed; runway concerns",
				Impact:      -0.08,
				RiskImpact:  0.04,
			})
		} else if financialRoll > 0.75 {
			findings = append(findings, DDFinding{
				Type:        "green_flag",
				Category:    "financial",
				Description: "Unit economics better than expected; clear path to profitability",
				Impact:      0.10,
				RiskImpact:  -0.02,
			})
		}
	}

	// Deep DD: Additional technical and legal findings
	if level == "deep" {
		// Technical audit
		techRoll := rand.Float64()
		if techRoll < 0.18 {
			findings = append(findings, DDFinding{
				Type:        "red_flag",
				Category:    "tech",
				Description: "Significant technical debt; scalability concerns identified",
				Impact:      -0.12,
				RiskImpact:  0.06,
			})
		} else if techRoll > 0.80 {
			findings = append(findings, DDFinding{
				Type:        "green_flag",
				Category:    "tech",
				Description: "Strong technical architecture; defensible IP portfolio",
				Impact:      0.12,
				RiskImpact:  -0.05,
			})
		} else {
			findings = append(findings, DDFinding{
				Type:        "neutral",
				Category:    "tech",
				Description: "Technology stack is solid with normal levels of technical debt",
				Impact:      0.0,
				RiskImpact:  0.0,
			})
		}

		// Legal/compliance check
		legalRoll := rand.Float64()
		if legalRoll < 0.12 {
			findings = append(findings, DDFinding{
				Type:        "red_flag",
				Category:    "legal",
				Description: "Pending litigation; IP ownership disputes",
				Impact:      -0.15,
				RiskImpact:  0.08,
			})
		} else if legalRoll > 0.88 {
			findings = append(findings, DDFinding{
				Type:        "green_flag",
				Category:    "legal",
				Description: "Clean cap table; all IP properly assigned",
				Impact:      0.05,
				RiskImpact:  -0.02,
			})
		}

		// Market positioning
		marketRoll := rand.Float64()
		if marketRoll < 0.15 {
			findings = append(findings, DDFinding{
				Type:        "red_flag",
				Category:    "market",
				Description: "Strong competitor just raised large round; market share concerns",
				Impact:      -0.10,
				RiskImpact:  0.05,
			})
		} else if marketRoll > 0.82 {
			findings = append(findings, DDFinding{
				Type:        "green_flag",
				Category:    "market",
				Description: "Market timing excellent; strong product-market fit signals",
				Impact:      0.10,
				RiskImpact:  -0.04,
			})
		}
	}

	// Always add a summary finding
	redFlags := 0
	greenFlags := 0
	totalImpact := 0.0
	totalRiskImpact := 0.0

	for _, f := range findings {
		if f.Type == "red_flag" {
			redFlags++
		} else if f.Type == "green_flag" {
			greenFlags++
		}
		totalImpact += f.Impact
		totalRiskImpact += f.RiskImpact
	}

	summaryType := "neutral"
	summaryDesc := "Overall assessment: Mixed signals, normal startup risk profile"

	if redFlags > greenFlags+1 {
		summaryType = "red_flag"
		summaryDesc = "Overall assessment: Significant concerns identified, elevated risk"
	} else if greenFlags > redFlags+1 {
		summaryType = "green_flag"
		summaryDesc = "Overall assessment: Strong opportunity, better than expected"
	}

	findings = append(findings, DDFinding{
		Type:        summaryType,
		Category:    "summary",
		Description: summaryDesc,
		Impact:      totalImpact,
		RiskImpact:  totalRiskImpact,
	})

	return findings
}

// ApplyDDFindings applies due diligence findings to startup
func ApplyDDFindings(startup *Startup, findings []DDFinding) {
	for _, finding := range findings {
		// Apply risk impact
		startup.RiskScore += finding.RiskImpact

		// Ensure risk stays in bounds
		if startup.RiskScore < 0.1 {
			startup.RiskScore = 0.1
		}
		if startup.RiskScore > 0.9 {
			startup.RiskScore = 0.9
		}

		// Apply growth impact
		startup.GrowthPotential += finding.Impact

		// Ensure growth stays in bounds
		if startup.GrowthPotential < 0.1 {
			startup.GrowthPotential = 0.1
		}
		if startup.GrowthPotential > 1.0 {
			startup.GrowthPotential = 1.0
		}
	}
}

// GetDDRelationshipBonus returns relationship boost from doing DD
func GetDDRelationshipBonus(level string) float64 {
	switch level {
	case "quick":
		return 3.0
	case "standard":
		return 6.0
	case "deep":
		return 10.0
	default:
		return 0.0
	}
}

// ShouldBlockInvestment checks if DD findings are so bad investment should be blocked
func ShouldBlockInvestment(findings []DDFinding) (bool, string) {
	redFlagCount := 0
	criticalIssues := []string{}

	for _, finding := range findings {
		if finding.Type == "red_flag" {
			redFlagCount++
			if finding.RiskImpact >= 0.07 { // Serious risk increase
				criticalIssues = append(criticalIssues, finding.Description)
			}
		}
	}

	// If 3+ red flags with at least one critical, warn player
	if redFlagCount >= 3 && len(criticalIssues) > 0 {
		return true, fmt.Sprintf("CRITICAL: %d red flags found including: %s. Recommend passing on this deal.",
			redFlagCount, criticalIssues[0])
	}

	return false, ""
}

// CalculateDDROI estimates if DD cost was worth it
func CalculateDDROI(ddCost int64, investmentAmount int64, findings []DDFinding) string {
	// Calculate total impact from findings
	totalImpact := 0.0
	for _, finding := range findings {
		totalImpact += finding.Impact
	}

	// Estimate value added/saved
	estimatedValue := int64(float64(investmentAmount) * totalImpact)

	if estimatedValue > ddCost*2 {
		return "Excellent - DD likely saved/added significant value"
	} else if estimatedValue > ddCost {
		return "Good - DD was worthwhile"
	} else if estimatedValue > 0 {
		return "Fair - DD provided some value"
	} else if estimatedValue < -ddCost {
		return "Poor - Should have skipped this deal"
	}

	return "Neutral - Normal due diligence outcome"
}
