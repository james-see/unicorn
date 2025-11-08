package founder

import (
	"fmt"
	"math/rand"
)

// InitializePartnershipIntegrations initializes enhanced partnership system
func InitializePartnershipIntegrations(fs *FounderState) {
	if fs.PartnershipIntegrations == nil {
		fs.PartnershipIntegrations = []PartnershipIntegration{}
	}
}

// CreateDeepIntegration creates a deep integration with a partner
func (fs *FounderState) CreateDeepIntegration(partnerName string, integrationDepth string) error {
	// Check if partnership exists
	partnershipExists := false
	for _, p := range fs.Partnerships {
		if p.Partner == partnerName && p.Status == "active" {
			partnershipExists = true
			break
		}
	}

	if !partnershipExists {
		return fmt.Errorf("partnership with %s not found or not active", partnerName)
	}

	// Check if integration already exists
	for _, pi := range fs.PartnershipIntegrations {
		if pi.PartnerName == partnerName {
			return fmt.Errorf("integration with %s already exists", partnerName)
		}
	}

	validDepths := map[string]bool{
		"surface": true,
		"deep":    true,
		"native":  true,
	}
	if !validDepths[integrationDepth] {
		return fmt.Errorf("invalid integration depth: %s", integrationDepth)
	}

	// Integration cost based on depth
	var cost int64
	switch integrationDepth {
	case "surface":
		cost = 20000 // $20k
	case "deep":
		cost = 50000 // $50k
	case "native":
		cost = 100000 // $100k
	}

	if cost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	fs.Cash -= cost

	// Create integration
	integration := PartnershipIntegration{
		PartnerName:      partnerName,
		IntegrationDepth: integrationDepth,
		RevenueShare:     0.0,
		CACReduction:     0.0,
		ChurnReduction:   0.0,
		MRRContribution:  0,
		CoMarketingActive: false,
		DataSharingActive: false,
	}

	// Set benefits based on depth
	switch integrationDepth {
	case "surface":
		integration.CACReduction = 0.05  // 5% CAC reduction
		integration.ChurnReduction = 0.01 // 1% churn reduction
	case "deep":
		integration.CACReduction = 0.15  // 15% CAC reduction
		integration.ChurnReduction = 0.03 // 3% churn reduction
		integration.RevenueShare = 0.10   // 10% revenue share
	case "native":
		integration.CACReduction = 0.25   // 25% CAC reduction
		integration.ChurnReduction = 0.05 // 5% churn reduction
		integration.RevenueShare = 0.20  // 20% revenue share
	}

	fs.PartnershipIntegrations = append(fs.PartnershipIntegrations, integration)

	return nil
}

// LaunchCoMarketingCampaign launches a co-marketing campaign with a partner
func (fs *FounderState) LaunchCoMarketingCampaign(partnerName string) error {
	// Find integration
	var integration *PartnershipIntegration
	for i := range fs.PartnershipIntegrations {
		if fs.PartnershipIntegrations[i].PartnerName == partnerName {
			integration = &fs.PartnershipIntegrations[i]
			break
		}
	}

	if integration == nil {
		return fmt.Errorf("integration with %s not found", partnerName)
	}

	if integration.CoMarketingActive {
		return fmt.Errorf("co-marketing campaign already active")
	}

	// Campaign cost: $10-30k
	cost := int64(10000 + rand.Int63n(20000))
	if cost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	fs.Cash -= cost
	integration.CoMarketingActive = true

	// Additional CAC reduction during campaign
	integration.CACReduction += 0.10 // +10% CAC reduction

	return nil
}

// EnableDataSharing enables data sharing partnership
func (fs *FounderState) EnableDataSharing(partnerName string) error {
	// Find integration
	var integration *PartnershipIntegration
	for i := range fs.PartnershipIntegrations {
		if fs.PartnershipIntegrations[i].PartnerName == partnerName {
			integration = &fs.PartnershipIntegrations[i]
			break
		}
	}

	if integration == nil {
		return fmt.Errorf("integration with %s not found", partnerName)
	}

	if integration.DataSharingActive {
		return fmt.Errorf("data sharing already active")
	}

	integration.DataSharingActive = true

	// Data sharing improves product (reduces churn)
	integration.ChurnReduction += 0.02 // +2% churn reduction

	return nil
}

// ProcessPartnershipIntegrations processes partnership benefits monthly
func (fs *FounderState) ProcessPartnershipIntegrations() []string {
	var messages []string
	totalCACReduction := 0.0
	totalChurnReduction := 0.0
	totalMRRContribution := int64(0)

	for i := range fs.PartnershipIntegrations {
		pi := &fs.PartnershipIntegrations[i]

		// Calculate MRR contribution from revenue share
		if pi.RevenueShare > 0 {
			// Estimate partnership deals as 5-15% of new MRR
			newMRRFromPartnership := int64(float64(fs.MRR) * pi.RevenueShare * (0.05 + rand.Float64()*0.10))
			pi.MRRContribution = newMRRFromPartnership
			totalMRRContribution += newMRRFromPartribution
		}

		totalCACReduction += pi.CACReduction
		totalChurnReduction += pi.ChurnReduction
	}

	// Apply benefits (capped)
	if totalCACReduction > 0.5 {
		totalCACReduction = 0.5 // Max 50% CAC reduction
	}
	if totalChurnReduction > 0.15 {
		totalChurnReduction = 0.15 // Max 15% churn reduction
	}

	// Apply MRR contribution
	if totalMRRContribution > 0 {
		fs.DirectMRR += totalMRRContribution
		fs.syncMRR()
		messages = append(messages, fmt.Sprintf("🤝 Partnerships contributed $%s MRR this month", formatCurrency(totalMRRContribution)))
	}

	return messages
}

// GetPartnershipIntegrationSummary returns summary of integrations
func (fs *FounderState) GetPartnershipIntegrationSummary() (totalIntegrations int, deepIntegrations int, nativeIntegrations int, totalMRR int64) {
	totalIntegrations = len(fs.PartnershipIntegrations)
	deepIntegrations = 0
	nativeIntegrations = 0
	totalMRR = 0

	for _, pi := range fs.PartnershipIntegrations {
		if pi.IntegrationDepth == "deep" {
			deepIntegrations++
		} else if pi.IntegrationDepth == "native" {
			nativeIntegrations++
		}
		totalMRR += pi.MRRContribution
	}

	return totalIntegrations, deepIntegrations, nativeIntegrations, totalMRR
}

