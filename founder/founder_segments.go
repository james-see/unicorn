package founder

import (
	"fmt"
	"math/rand"
)

// InitializeSegments sets up default customer segments
func (fs *FounderState) InitializeSegments() {
	fs.CustomerSegments = []CustomerSegment{
		{
			Name:                "Enterprise",
			AvgDealSize:         50000, // $50k/month
			ChurnRate:           0.05,  // 5% monthly
			SalesCycle:          9,     // 9 months
			CAC:                 25000, // $25k per customer
			FeatureRequirements: []string{"Enterprise SSO", "Advanced Analytics", "Security Suite"},
			Volume:              0,
		},
		{
			Name:                "Mid-Market",
			AvgDealSize:         10000, // $10k/month
			ChurnRate:           0.08,  // 8% monthly
			SalesCycle:          4,     // 4 months
			CAC:                 8000,  // $8k per customer
			FeatureRequirements: []string{"REST API", "Integrations Hub"},
			Volume:              0,
		},
		{
			Name:                "SMB",
			AvgDealSize:         1000, // $1k/month
			ChurnRate:           0.15, // 15% monthly
			SalesCycle:          1,    // 1 month
			CAC:                 1500, // $1.5k per customer
			FeatureRequirements: []string{"Mobile App"},
			Volume:              0,
		},
		{
			Name:                "Startup",
			AvgDealSize:         500,  // $500/month
			ChurnRate:           0.25, // 25% monthly (high churn)
			SalesCycle:          1,    // 1 month
			CAC:                 800,  // $800 per customer
			FeatureRequirements: []string{},
			Volume:              0,
		},
	}
}

// InitializeVerticals sets up industry verticals
func (fs *FounderState) InitializeVerticals() {
	fs.VerticalFocuses = []VerticalFocus{
		{
			Industry:    "FinTech",
			MarketSize:  50000,
			Competition: "very_high",
			IsActive:    false,
		},
		{
			Industry:    "HealthTech",
			MarketSize:  40000,
			Competition: "high",
			IsActive:    false,
		},
		{
			Industry:    "Retail",
			MarketSize:  60000,
			Competition: "high",
			IsActive:    false,
		},
		{
			Industry:    "Manufacturing",
			MarketSize:  35000,
			Competition: "medium",
			IsActive:    false,
		},
		{
			Industry:    "Education",
			MarketSize:  30000,
			Competition: "medium",
			IsActive:    false,
		},
		{
			Industry:    "Real Estate",
			MarketSize:  25000,
			Competition: "low",
			IsActive:    false,
		},
		{
			Industry:    "Legal",
			MarketSize:  20000,
			Competition: "low",
			IsActive:    false,
		},
		{
			Industry:    "Media",
			MarketSize:  28000,
			Competition: "medium",
			IsActive:    false,
		},
	}
}

// SelectICP sets the focused customer segment
func (fs *FounderState) SelectICP(segmentName string) error {
	if fs.Customers < 50 {
		return fmt.Errorf("need at least 50 customers to focus on a segment")
	}

	// Validate segment exists
	found := false
	for _, seg := range fs.CustomerSegments {
		if seg.Name == segmentName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid segment: %s", segmentName)
	}

	fs.SelectedICP = segmentName
	return nil
}

// SelectVertical sets the focused industry vertical
func (fs *FounderState) SelectVertical(industry string) error {
	if fs.Customers < 50 {
		return fmt.Errorf("need at least 50 customers to focus on a vertical")
	}

	// Find and activate the vertical
	found := false
	for i := range fs.VerticalFocuses {
		if fs.VerticalFocuses[i].Industry == industry {
			fs.VerticalFocuses[i].IsActive = true
			fs.VerticalFocuses[i].ICPMatch = 0.0 // Start at 0, builds over time
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid vertical: %s", industry)
	}

	fs.SelectedVertical = industry
	return nil
}

// GetICPBenefits calculates the benefits of having a focused ICP
func (fs *FounderState) GetICPBenefits() (cacReduction float64, closeRateBonus float64, dealSizeBonus float64) {
	if fs.SelectedICP == "" && fs.SelectedVertical == "" {
		return 0, 0, 0 // No focus = no benefits
	}

	// Base benefits for having any focus
	cacReduction = 0.20      // 20% CAC reduction
	closeRateBonus = 0.15    // 15% close rate increase
	dealSizeBonus = 0.10     // 10% deal size increase

	// Additional benefits if both segment AND vertical are selected
	if fs.SelectedICP != "" && fs.SelectedVertical != "" {
		cacReduction += 0.10    // Additional 10% CAC reduction
		closeRateBonus += 0.10  // Additional 10% close rate
		dealSizeBonus += 0.05   // Additional 5% deal size
	}

	// Bonus for specialized sales team
	if fs.SelectedVertical != "" {
		for i := range fs.VerticalFocuses {
			if fs.VerticalFocuses[i].Industry == fs.SelectedVertical && fs.VerticalFocuses[i].IsActive {
				// ICPMatch grows over time as you specialize
				matchBonus := fs.VerticalFocuses[i].ICPMatch * 0.15 // Up to 15% additional benefits
				cacReduction += matchBonus
				closeRateBonus += matchBonus
				dealSizeBonus += matchBonus * 0.5
				break
			}
		}
	}

	return cacReduction, closeRateBonus, dealSizeBonus
}

// UpdateSegmentVolumes recounts customers in each segment
func (fs *FounderState) UpdateSegmentVolumes() {
	// Reset all volumes
	for i := range fs.CustomerSegments {
		fs.CustomerSegments[i].Volume = 0
	}

	// Count customers by segment
	for _, customer := range fs.CustomerList {
		if !customer.IsActive {
			continue
		}

		// Determine segment based on deal size
		segment := fs.DetermineSegmentFromDealSize(customer.DealSize)
		for i := range fs.CustomerSegments {
			if fs.CustomerSegments[i].Name == segment {
				fs.CustomerSegments[i].Volume++
				break
			}
		}
	}
}

// DetermineSegmentFromDealSize returns the segment name based on deal size
func (fs *FounderState) DetermineSegmentFromDealSize(dealSize int64) string {
	if dealSize >= 50000 {
		return "Enterprise"
	} else if dealSize >= 10000 {
		return "Mid-Market"
	} else if dealSize >= 1000 {
		return "SMB"
	}
	return "Startup"
}

// GenerateDealSizeForSegment generates a deal size appropriate for the segment
func (fs *FounderState) GenerateDealSizeForSegment(segmentName string) int64 {
	for _, seg := range fs.CustomerSegments {
		if seg.Name == segmentName {
			// Generate deal size with ±30% variance
			variance := 0.30
			multiplier := 1.0 + (rand.Float64()*variance*2 - variance)
			dealSize := int64(float64(seg.AvgDealSize) * multiplier)
			
			// Apply ICP benefits if focused on this segment
			if fs.SelectedICP == segmentName {
				_, _, dealSizeBonus := fs.GetICPBenefits()
				dealSize = int64(float64(dealSize) * (1.0 + dealSizeBonus))
			}
			
			return dealSize
		}
	}
	return fs.AvgDealSize // Fallback
}

// ProcessSegmentFocus updates vertical focus metrics monthly
func (fs *FounderState) ProcessSegmentFocus() []string {
	var messages []string

	// Update vertical ICP match over time (improves with focus)
	if fs.SelectedVertical != "" {
		for i := range fs.VerticalFocuses {
			if fs.VerticalFocuses[i].Industry == fs.SelectedVertical && fs.VerticalFocuses[i].IsActive {
				// ICP match improves by 0.05 per month, capped at 1.0
				if fs.VerticalFocuses[i].ICPMatch < 1.0 {
					fs.VerticalFocuses[i].ICPMatch += 0.05
					if fs.VerticalFocuses[i].ICPMatch > 1.0 {
						fs.VerticalFocuses[i].ICPMatch = 1.0
					}
					
					if fs.VerticalFocuses[i].ICPMatch >= 1.0 {
						messages = append(messages, fmt.Sprintf("🎯 %s vertical focus fully optimized! Maximum benefits unlocked", fs.SelectedVertical))
					}
				}
				break
			}
		}
	}

	// Update segment volumes
	fs.UpdateSegmentVolumes()

	return messages
}

// GetSegmentConcentration returns what % of customers are in the selected ICP
func (fs *FounderState) GetSegmentConcentration() float64 {
	if fs.SelectedICP == "" || fs.Customers == 0 {
		return 0
	}

	for _, seg := range fs.CustomerSegments {
		if seg.Name == fs.SelectedICP {
			return float64(seg.Volume) / float64(fs.Customers)
		}
	}
	return 0
}

// SuggestSegmentForNewCustomer suggests which segment a new customer should be in
// based on ICP focus and current distribution
func (fs *FounderState) SuggestSegmentForNewCustomer() string {
	// If focused on a segment, 70% chance of getting that segment
	if fs.SelectedICP != "" && rand.Float64() < 0.70 {
		return fs.SelectedICP
	}

	// Otherwise, random distribution weighted by typical mix
	// Enterprise: 10%, Mid-Market: 30%, SMB: 40%, Startup: 20%
	roll := rand.Float64()
	if roll < 0.10 {
		return "Enterprise"
	} else if roll < 0.40 {
		return "Mid-Market"
	} else if roll < 0.80 {
		return "SMB"
	}
	return "Startup"
}

// ChangeICP pivots to a new segment focus (with penalties)
func (fs *FounderState) ChangeICP(newSegment string) error {
	if fs.SelectedICP == "" {
		return fmt.Errorf("no existing ICP to change from")
	}

	if fs.SelectedICP == newSegment {
		return fmt.Errorf("already focused on %s", newSegment)
	}

	// Validate new segment exists
	found := false
	for _, seg := range fs.CustomerSegments {
		if seg.Name == newSegment {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid segment: %s", newSegment)
	}

	// Lose 10-20% of customers from old segment due to pivot
	lossRate := 0.10 + rand.Float64()*0.10
	for i := range fs.CustomerSegments {
		if fs.CustomerSegments[i].Name == fs.SelectedICP {
			customersLost := int(float64(fs.CustomerSegments[i].Volume) * lossRate)
			
			// Mark customers as churned
			lostCount := 0
			for j := range fs.CustomerList {
				if fs.CustomerList[j].IsActive {
					customerSegment := fs.DetermineSegmentFromDealSize(fs.CustomerList[j].DealSize)
					if customerSegment == fs.SelectedICP && lostCount < customersLost {
						fs.churnCustomer(fs.CustomerList[j].ID)
						fs.MRR -= fs.CustomerList[j].DealSize
						fs.DirectMRR -= fs.CustomerList[j].DealSize
						lostCount++
					}
				}
			}
			
			fs.Customers -= lostCount
			fs.DirectCustomers -= lostCount
			break
		}
	}

	fs.SelectedICP = newSegment
	return nil
}

// ChangeVertical pivots to a new vertical focus (with penalties)
func (fs *FounderState) ChangeVertical(newVertical string) error {
	if fs.SelectedVertical == "" {
		return fmt.Errorf("no existing vertical to change from")
	}

	if fs.SelectedVertical == newVertical {
		return fmt.Errorf("already focused on %s", newVertical)
	}

	// Deactivate old vertical
	for i := range fs.VerticalFocuses {
		if fs.VerticalFocuses[i].Industry == fs.SelectedVertical {
			fs.VerticalFocuses[i].IsActive = false
			fs.VerticalFocuses[i].ICPMatch = 0
			break
		}
	}

	// Activate new vertical
	return fs.SelectVertical(newVertical)
}

