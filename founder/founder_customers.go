package founder

import (
	"fmt"
	"math"
	"math/rand"
)


func (fs *FounderState) addCustomer(dealSize int64, source string) Customer {
	// Determine if contract is perpetual (80% chance) or fixed term (20% chance)
	termMonths := 0 // Default to perpetual
	if rand.Float64() < 0.2 {
		// Fixed term contracts: 6, 12, 24, or 36 months
		terms := []int{6, 12, 24, 36}
		termMonths = terms[rand.Intn(len(terms))]
	}

	// Initial health score based on product maturity and source
	healthScore := 0.5 + (fs.ProductMaturity * 0.3) // Base: 0.5-0.8
	if source == "affiliate" {
		healthScore -= 0.1 // Affiliates slightly less sticky
	}
	healthScore = math.Min(1.0, math.Max(0.1, healthScore))

	customer := Customer{
		ID:           fs.NextCustomerID,
		Source:       source,
		DealSize:     dealSize,
		TermMonths:   termMonths,
		MonthAdded:   fs.Turn,
		MonthChurned: 0,
		IsActive:     true,
		HealthScore:  healthScore,
	}

	fs.CustomerList = append(fs.CustomerList, customer)
	fs.NextCustomerID++
	fs.TotalCustomersEver++

	return customer
}

func (fs *FounderState) churnCustomer(customerID int) {
	for i := range fs.CustomerList {
		if fs.CustomerList[i].ID == customerID && fs.CustomerList[i].IsActive {
			fs.CustomerList[i].IsActive = false
			fs.CustomerList[i].MonthChurned = fs.Turn
			fs.TotalChurned++
			break
		}
	}
}

func (fs *FounderState) GetActiveCustomers() []Customer {
	var active []Customer
	for _, c := range fs.CustomerList {
		if c.IsActive {
			active = append(active, c)
		}
	}
	return active
}

func (fs *FounderState) GetChurnedCustomers() []Customer {
	var churned []Customer
	for _, c := range fs.CustomerList {
		if !c.IsActive {
			churned = append(churned, c)
		}
	}
	return churned
}

func (fs *FounderState) syncMRR() {
	var directMRR int64
	var affiliateMRR int64
	var directCount int
	var affiliateCount int

	// Recalculate from actual customer list to prevent accumulation errors
	for _, c := range fs.CustomerList {
		if !c.IsActive {
			continue
		}

		if c.Source == "affiliate" {
			affiliateMRR += c.DealSize
			affiliateCount++
		} else {
			// "direct", "partnership", "market" all count as direct
			directMRR += c.DealSize
			directCount++
		}
	}

	fs.DirectMRR = directMRR
	fs.AffiliateMRR = affiliateMRR
	fs.MRR = directMRR + affiliateMRR
	fs.DirectCustomers = directCount
	fs.AffiliateCustomers = affiliateCount
	fs.Customers = directCount + affiliateCount

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}
}

func (fs *FounderState) updateDealSizeRange(newDealSize int64) {
	if fs.MinDealSize == 0 || newDealSize < fs.MinDealSize {
		fs.MinDealSize = newDealSize
	}
	if newDealSize > fs.MaxDealSize {
		fs.MaxDealSize = newDealSize
	}
}

func (fs *FounderState) RecalculateChurnRate() {
	// Calculate base churn from product maturity
	baseChurnFromMaturity := (1.0-fs.ProductMaturity)*0.65 + 0.05
	baseChurnFromMaturity = math.Max(0.05, math.Min(0.70, baseChurnFromMaturity))
	baseChurn := baseChurnFromMaturity

	// CS team reduces churn
	// COO counts as 3x CS reps
	csImpact := 0.0
	for _, cs := range fs.Team.CustomerSuccess {
		csImpact += (cs.Impact * 0.02) // Each CS rep reduces churn by ~2%
	}
	for _, exec := range fs.Team.Executives {
		if exec.Role == RoleCOO {
			csImpact += (exec.Impact * 0.02) // COO has 3x impact already built into Impact field
		}
	}

	// Calculate effective churn rate (after CS team impact)
	actualChurn := math.Max(0.01, baseChurn-csImpact)
	fs.CustomerChurnRate = actualChurn
	fs.ChurnRate = fs.CustomerChurnRate
}

// HireEmployee adds a new team member

func (fs *FounderState) UpdateCAC() {
	// Start with business-specific base CAC
	effectiveCAC := float64(fs.BaseCAC)

	// Product maturity reduces CAC (better product = better conversion)
	// At 100% maturity, CAC is 60% of base (40% reduction)
	maturityDiscount := fs.ProductMaturity * 0.4
	effectiveCAC *= (1.0 - maturityDiscount)

	// Competition increases CAC
	switch fs.CompetitionLevel {
	case "very_high":
		effectiveCAC *= 1.5 // +50%
	case "high":
		effectiveCAC *= 1.3 // +30%
	case "medium":
		effectiveCAC *= 1.1 // +10%
		// low = no change
	}

	fs.CustomerAcquisitionCost = int64(effectiveCAC)
}


func (fs *FounderState) SpendOnMarketing(amount int64) int {
	if amount > fs.Cash {
		return 0
	}

	fs.Cash -= amount

	// Use current effective CAC (which accounts for product maturity and competition)
	fs.UpdateCAC()

	newCustomers := int(amount / fs.CustomerAcquisitionCost)

	// Calculate MRR with variable deal sizes
	// Use category-based defaults if AvgDealSize is 0 (no customers yet)
	baseDealSize := fs.AvgDealSize
	if baseDealSize == 0 {
		switch fs.Category {
		case "SaaS":
			baseDealSize = 1000 // Default $1k/month for SaaS
		case "DeepTech":
			baseDealSize = 5000 // Default $5k/month for DeepTech
		case "GovTech":
			baseDealSize = 2000 // Default $2k/month for GovTech
		case "Hardware":
			baseDealSize = 3000 // Default $3k/month for Hardware
		default:
			baseDealSize = 1000 // Default $1k/month
		}
	}

	var totalMRR int64
	var dealSizes []int64 // Store deal sizes for customer tracking
	for i := 0; i < newCustomers; i++ {
		dealSize := generateDealSize(baseDealSize, fs.Category)
		fs.updateDealSizeRange(dealSize)
		totalMRR += dealSize
		dealSizes = append(dealSizes, dealSize)
	}

	// These are direct customers (not from affiliate program)
	fs.Customers += newCustomers
	fs.DirectCustomers += newCustomers
	fs.DirectMRR += totalMRR

	// Add customers to tracking system
	for _, dealSize := range dealSizes {
		fs.addCustomer(dealSize, "direct")
	}

	// Sync MRR from DirectMRR + AffiliateMRR
	fs.syncMRR()

	// Recalculate average deal size
	if fs.Customers > 0 {
		fs.AvgDealSize = fs.MRR / int64(fs.Customers)
	}

	fs.CalculateRunway()

	return newCustomers
}

// CheckForAcquisition checks if an acquisition offer comes in

func (fs *FounderState) SolicitCustomerFeedback() error {
	if fs.Customers == 0 {
		return fmt.Errorf("no customers to solicit feedback from")
	}

	// Feedback improves product maturity by 1-5% based on customer count
	// More customers = better feedback = more improvement
	improvement := 0.01 + (float64(fs.Customers)/100.0)*0.04 // 1-5% improvement
	if improvement > 0.05 {
		improvement = 0.05 // Cap at 5%
	}

	fs.ProductMaturity = math.Min(1.0, fs.ProductMaturity+improvement)

	// Customer feedback also reduces churn by 3-10%
	churnReduction := 0.03 + rand.Float64()*0.07                               // 3-10% reduction
	fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-churnReduction) // Minimum 1% churn
	fs.ChurnRate = fs.CustomerChurnRate

	return nil
}

// ============================================================================
// KEY METRICS & ANALYTICS
// ============================================================================

// CalculateLTVToCAC calculates the lifetime value to customer acquisition cost ratio