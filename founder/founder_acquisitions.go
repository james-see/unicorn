package founder

import (
	"fmt"
	"math/rand"
)

// InitializeAcquisitions initializes the acquisitions system
func InitializeAcquisitions(fs *FounderState) {
	if fs.AcquisitionTargets == nil {
		fs.AcquisitionTargets = []AcquisitionTarget{}
	}
	if fs.Acquisitions == nil {
		fs.Acquisitions = []Acquisition{}
	}
}

// CanAcquire checks if founder can acquire companies
func (fs *FounderState) CanAcquire() bool {
	// Unlock: $2M+ ARR OR Series B raised
	arr := fs.MRR * 12
	hasSeriesB := false
	for _, round := range fs.FundingRounds {
		if round.RoundName == "Series B" {
			hasSeriesB = true
			break
		}
	}
	return arr >= 2000000 || hasSeriesB
}

// GenerateAcquisitionTargets generates 1-3 acquisition targets per quarter
func (fs *FounderState) GenerateAcquisitionTargets() {
	if !fs.CanAcquire() {
		return
	}

	// Generate targets quarterly (every 3 months)
	if fs.Turn%3 != 0 {
		return
	}

	// Remove expired targets
	activeTargets := []AcquisitionTarget{}
	for _, target := range fs.AcquisitionTargets {
		if fs.Turn-target.MonthAppeared < target.ExpiresIn {
			activeTargets = append(activeTargets, target)
		}
	}
	fs.AcquisitionTargets = activeTargets

	// Generate 1-3 new targets
	numTargets := 1 + rand.Intn(3)
	for i := 0; i < numTargets; i++ {
		target := fs.generateAcquisitionTarget()
		if target != nil {
			fs.AcquisitionTargets = append(fs.AcquisitionTargets, *target)
		}
	}
}

// generateAcquisitionTarget creates a random acquisition target
func (fs *FounderState) generateAcquisitionTarget() *AcquisitionTarget {
	// Target MRR: 5-50% of your MRR
	targetMRR := int64(float64(fs.MRR) * (0.05 + rand.Float64()*0.45))
	if targetMRR < 1000 {
		targetMRR = 1000 // Minimum $1k MRR
	}

	// Customers: roughly proportional to MRR
	avgDealSize := targetMRR / 20 // Assume ~20 customers for the MRR
	if avgDealSize < 100 {
		avgDealSize = 100
	}
	customers := int(targetMRR / avgDealSize)
	if customers < 5 {
		customers = 5
	}

	// Acquisition cost: 2-10x revenue multiple
	revenueMultiple := 2.0 + rand.Float64()*8.0
	acquisitionCost := int64(float64(targetMRR*12) * revenueMultiple)
	if acquisitionCost < 500000 {
		acquisitionCost = 500000 // Minimum $500k
	}
	if acquisitionCost > 5000000 {
		acquisitionCost = 5000000 // Maximum $5M
	}

	// Integration cost: 20-50% of acquisition cost
	integrationCost := int64(float64(acquisitionCost) * (0.20 + rand.Float64()*0.30))

	// Synergy bonus: 10-50% revenue boost
	synergyBonus := 0.10 + rand.Float64()*0.40

	// Risk level
	riskRoll := rand.Float64()
	risk := "medium"
	if riskRoll < 0.3 {
		risk = "low"
	} else if riskRoll > 0.7 {
		risk = "high"
	}

	// Technology/IP gained
	technologies := []string{}
	techOptions := []string{"API Integration", "Mobile SDK", "Analytics Engine", "ML Models", "Customer Data", "Brand Assets"}
	numTech := 1 + rand.Intn(3)
	for i := 0; i < numTech && i < len(techOptions); i++ {
		technologies = append(technologies, techOptions[rand.Intn(len(techOptions))])
	}

	// Team size: 2-10 people
	teamSize := 2 + rand.Intn(9)

	// Company names
	companyNames := []string{
		"TechFlow", "DataSync", "CloudBridge", "AppForge", "DevTools Pro",
		"Analytics Co", "MobileFirst", "API Gateway", "CodeBase", "PlatformX",
	}

	return &AcquisitionTarget{
		Name:            companyNames[rand.Intn(len(companyNames))],
		Category:        fs.Category,
		MRR:             targetMRR,
		Customers:       customers,
		TeamSize:        teamSize,
		Technology:      technologies,
		AcquisitionCost: acquisitionCost,
		IntegrationCost: integrationCost,
		SynergyBonus:    synergyBonus,
		Risk:            risk,
		MonthAppeared:   fs.Turn,
		ExpiresIn:       3, // Expires in 3 months
	}
}

// AcquireCompany executes an acquisition
func (fs *FounderState) AcquireCompany(targetIndex int) (*Acquisition, error) {
	if targetIndex < 0 || targetIndex >= len(fs.AcquisitionTargets) {
		return nil, fmt.Errorf("invalid target index")
	}

	target := fs.AcquisitionTargets[targetIndex]
	totalCost := target.AcquisitionCost + target.IntegrationCost

	if totalCost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash (need $%s)", formatCurrency(totalCost))
	}

	fs.Cash -= totalCost

	// Create acquisition record
	acquisition := Acquisition{
		TargetName:        target.Name,
		Month:             fs.Turn,
		Cost:              totalCost,
		CustomersGained:   target.Customers,
		MRRGained:         target.MRR,
		TeamGained:        target.TeamSize,
		Success:           false, // Will be determined after integration
		IntegrationMonths: 3 + rand.Intn(4), // 3-6 months
		IntegrationProgress: 0,
		SynergyRealized:   0.0,
	}

	// Add customers and MRR immediately (but integration affects retention)
	fs.Customers += target.Customers
	fs.DirectCustomers += target.Customers
	fs.DirectMRR += target.MRR
	fs.syncMRR()

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

	// Add to acquisitions list
	fs.Acquisitions = append(fs.Acquisitions, acquisition)

	// Remove target from available targets
	fs.AcquisitionTargets = append(fs.AcquisitionTargets[:targetIndex], fs.AcquisitionTargets[targetIndex+1:]...)

	return &acquisition, nil
}

// ProcessAcquisitionIntegration processes integration progress for active acquisitions
func (fs *FounderState) ProcessAcquisitionIntegration() []string {
	var messages []string

	for i := range fs.Acquisitions {
		acq := &fs.Acquisitions[i]
		if acq.Success || acq.IntegrationProgress >= 100 {
			continue
		}

		// Integration progress: depends on engineering team size
		engineers := len(fs.Team.Engineers)
		progressPerMonth := 20.0 // Base 20% per month
		if engineers > 0 {
			progressPerMonth += float64(engineers) * 5.0 // +5% per engineer
		}

		// Technical debt slows integration
		if fs.TechnicalDebt != nil && fs.TechnicalDebt.CurrentLevel > 50 {
			progressPerMonth *= (1.0 - float64(fs.TechnicalDebt.CurrentLevel)/200.0)
		}

		acq.IntegrationProgress += int(progressPerMonth)
		if acq.IntegrationProgress > 100 {
			acq.IntegrationProgress = 100
		}

		// Check if integration complete
		if acq.IntegrationProgress >= 100 {
			acq.Success = true

		// Determine synergy realized (based on integration quality)
		// Find the original target to get risk level
		baseSynergy := 0.5 // Start at 50% of promised synergy
		// Risk assessment based on integration speed and team size
		if acq.TeamGained > 5 {
			baseSynergy = 0.6 // Larger teams = better integration
		}

			// Integration speed bonus
			expectedMonths := acq.IntegrationMonths
			actualMonths := fs.Turn - acq.Month
			if actualMonths < expectedMonths {
				baseSynergy += 0.1 // Bonus for fast integration
			}

			acq.SynergyRealized = baseSynergy

			// Apply synergy bonus to MRR
			synergyMRR := int64(float64(acq.MRRGained) * acq.SynergyRealized)
			fs.DirectMRR += synergyMRR
			fs.syncMRR()

			// Reduce CAC due to cross-selling
			fs.BaseCAC = int64(float64(fs.BaseCAC) * 0.95) // -5% CAC

			messages = append(messages, fmt.Sprintf("✅ Acquisition integration complete: %s (+$%s MRR synergy)", acq.TargetName, formatCurrency(synergyMRR)))
		} else {
			// During integration, some customers may churn
			integrationChurnRate := 0.02 // 2% churn per month during integration
			// Higher churn for larger acquisitions
			if acq.TeamGained > 10 {
				integrationChurnRate = 0.05 // 5% for large acquisitions
			}

			customersLost := int(float64(acq.CustomersGained) * integrationChurnRate)
			if customersLost > 0 {
				mrrLost := int64(customersLost) * fs.AvgDealSize
				fs.Customers -= customersLost
				fs.DirectCustomers -= customersLost
				fs.DirectMRR -= mrrLost
				fs.syncMRR()

				acq.CustomersGained -= customersLost
				acq.MRRGained -= mrrLost

				messages = append(messages, fmt.Sprintf("⚠️  %s integration: Lost %d customers during transition", acq.TargetName, customersLost))
			}

			// Integration reduces engineering velocity
			if fs.TechnicalDebt != nil {
				fs.TechnicalDebt.CurrentLevel += 5 // +5 tech debt per month during integration
			}
		}
	}

	return messages
}

// GetAcquisitionSummary returns summary of acquisitions
func (fs *FounderState) GetAcquisitionSummary() (completed int, inProgress int, totalMRRGained int64, totalCost int64) {
	completed = 0
	inProgress = 0
	totalMRRGained = 0
	totalCost = 0

	for _, acq := range fs.Acquisitions {
		totalCost += acq.Cost
		if acq.Success {
			completed++
			totalMRRGained += acq.MRRGained + int64(float64(acq.MRRGained)*acq.SynergyRealized)
		} else {
			inProgress++
			totalMRRGained += acq.MRRGained // Count current MRR even if not fully integrated
		}
	}

	return completed, inProgress, totalMRRGained, totalCost
}

