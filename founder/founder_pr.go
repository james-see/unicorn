package founder

import "fmt"

// LaunchPRProgram starts PR and media relations
func (fs *FounderState) LaunchPRProgram(monthlyRetainer int64) error {
	if fs.MRR < 500000 {
		return fmt.Errorf("PR programs recommended at $500k+ MRR")
	}

	fs.PRProgram = &PRProgram{
		HasPRFirm:       true,
		MonthlyRetainer: monthlyRetainer,
		Campaigns:       []PRCampaign{},
		MediaCoverage:   []MediaCoverage{},
		BrandScore:      50, // Start at 50/100
		LaunchedMonth:   fs.Turn,
	}

	return nil
}

// LaunchPRCampaign runs a specific PR initiative
func (fs *FounderState) LaunchPRCampaign(campaignType string, cost int64) error {
	if fs.PRProgram == nil {
		return fmt.Errorf("no PR program active")
	}

	if fs.Cash < cost {
		return fmt.Errorf("insufficient cash")
	}

	fs.Cash -= cost

	campaign := PRCampaign{
		Type:        campaignType,
		Cost:        cost,
		Duration:    3, // 3 months
		TargetMedia: []string{"TechCrunch", "WSJ"},
		StartMonth:  fs.Turn,
		Success:     true,
		Impact: PRImpact{
			CACReduction:   0.15, // -15% CAC
			BrandBoost:     10,
			InboundLeads:   20,
			DurationMonths: 6,
		},
	}

	fs.PRProgram.Campaigns = append(fs.PRProgram.Campaigns, campaign)
	fs.PRProgram.BrandScore += 10

	return nil
}

// UpdatePRProgram processes monthly PR effects
func (fs *FounderState) UpdatePRProgram() []string {
	var messages []string

	if fs.PRProgram == nil {
		return messages
	}

	// Pay retainer
	if fs.Cash < fs.PRProgram.MonthlyRetainer {
		return messages
	}
	fs.Cash -= fs.PRProgram.MonthlyRetainer

	// Check for active campaigns
	activeCampaigns := 0
	for _, campaign := range fs.PRProgram.Campaigns {
		if fs.Turn-campaign.StartMonth < campaign.Duration {
			activeCampaigns++
		}
	}

	if activeCampaigns > 0 {
		messages = append(messages, fmt.Sprintf("📰 %d active PR campaigns boosting brand awareness", activeCampaigns))
	}

	return messages
}

