package founder

import (
	"fmt"
	"math/rand"
)


func (fs *FounderState) StartPartnership(partnerType string) (*Partnership, error) {
	// Ensure MRR is synced before calculating boost
	fs.syncMRR()

	partners := map[string][]string{
		"distribution": {"Salesforce", "HubSpot", "Oracle", "SAP", "Adobe"},
		"technology":   {"AWS", "Google Cloud", "Microsoft Azure", "IBM", "MongoDB"},
		"co-marketing": {"Shopify", "Stripe", "Zendesk", "Slack", "Atlassian"},
		"data":         {"Snowflake", "Databricks", "Tableau", "Segment", "Amplitude"},
	}

	partnerList, ok := partners[partnerType]
	if !ok {
		return nil, fmt.Errorf("unknown partnership type: %s", partnerType)
	}

	partner := partnerList[rand.Intn(len(partnerList))]

	// Calculate costs and benefits
	var cost, mrrBoost int64
	var churnReduction float64
	var duration int

	switch partnerType {
	case "distribution":
		cost = 50000 + rand.Int63n(100000)                             // $50-150k
		mrrBoost = int64(float64(fs.MRR) * (0.1 + rand.Float64()*0.2)) // 10-30% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - helps acquire first customers
			mrrBoost = 5000 + rand.Int63n(15000) // $5-20k/month minimum
		}
		churnReduction = 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
		duration = 12 + rand.Intn(12)               // 12-24 months
	case "technology":
		cost = 30000 + rand.Int63n(70000)                                // $30-100k
		mrrBoost = int64(float64(fs.MRR) * (0.05 + rand.Float64()*0.15)) // 5-20% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - product integration helps attract customers
			mrrBoost = 3000 + rand.Int63n(7000) // $3-10k/month minimum
		}
		churnReduction = 0.02 + rand.Float64()*0.03 // 2-5% churn reduction
		duration = 12 + rand.Intn(24)               // 12-36 months
	case "co-marketing":
		cost = 25000 + rand.Int63n(50000)                                // $25-75k
		mrrBoost = int64(float64(fs.MRR) * (0.15 + rand.Float64()*0.25)) // 15-40% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - marketing helps acquire customers
			mrrBoost = 8000 + rand.Int63n(12000) // $8-20k/month minimum
		}
		churnReduction = 0.005 + rand.Float64()*0.015 // 0.5-2% churn reduction
		duration = 6 + rand.Intn(12)                  // 6-18 months
	case "data":
		cost = 40000 + rand.Int63n(60000)                                // $40-100k
		mrrBoost = int64(float64(fs.MRR) * (0.08 + rand.Float64()*0.12)) // 8-20% MRR boost
		if mrrBoost == 0 && fs.MRR == 0 {
			// Minimum boost even with no MRR - analytics help attract customers
			mrrBoost = 4000 + rand.Int63n(8000) // $4-12k/month minimum
		}
		churnReduction = 0.01 + rand.Float64()*0.02 // 1-3% churn reduction
		duration = 12 + rand.Intn(24)               // 12-36 months
	}

	if cost > fs.Cash {
		return nil, fmt.Errorf("insufficient cash for partnership (need $%s)", formatCurrency(cost))
	}

	partnership := Partnership{
		Partner:        partner,
		Type:           partnerType,
		MonthStarted:   fs.Turn,
		Duration:       duration,
		Cost:           cost,
		MRRBoost:       mrrBoost,
		ChurnReduction: churnReduction,
		Status:         "active",
	}

	fs.Cash -= cost
	fs.Partnerships = append(fs.Partnerships, partnership)

	// Apply partnership benefits immediately
	fs.MRR += mrrBoost
	fs.DirectMRR += mrrBoost // Partnership boost goes to direct MRR
	fs.CustomerChurnRate -= churnReduction
	if fs.CustomerChurnRate < 0 {
		fs.CustomerChurnRate = 0
	}

	// Sync MRR to ensure consistency
	fs.syncMRR()

	return &partnership, nil
}


func (fs *FounderState) UpdatePartnerships() []string {
	var messages []string

	for i := range fs.Partnerships {
		p := &fs.Partnerships[i]
		if p.Status != "active" {
			continue
		}

		// Check if partnership has expired
		monthsActive := fs.Turn - p.MonthStarted
		if monthsActive >= p.Duration {
			p.Status = "expired"
			messages = append(messages, fmt.Sprintf("⏰ Partnership with %s has expired", p.Partner))

			// Remove benefits
			fs.MRR -= p.MRRBoost
			fs.CustomerChurnRate += p.ChurnReduction
			continue
		}

		// Apply ongoing benefits (already included in calculations)
	}

	return messages
}

// ============================================================================
// AFFILIATE PROGRAM
// ============================================================================


func (fs *FounderState) LaunchAffiliateProgram(commission float64) error {
	if fs.AffiliateProgram != nil {
		return fmt.Errorf("affiliate program already running")
	}

	setupCost := int64(20000 + rand.Int63n(30000)) // $20-50k setup
	if setupCost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(setupCost))
	}

	fs.Cash -= setupCost

	fs.AffiliateProgram = &AffiliateProgram{
		LaunchedMonth:      fs.Turn,
		Commission:         commission / 100,  // Convert % to decimal
		Affiliates:         5 + rand.Intn(10), // Start with 5-15 affiliates
		SetupCost:          setupCost,
		MonthlyPlatformFee: 5000 + rand.Int63n(5000), // $5-10k/month
		MonthlyRevenue:     0,
		CustomersAcquired:  0,
	}

	return nil
}


func (fs *FounderState) UpdateAffiliateProgram() []string {
	var messages []string

	if fs.AffiliateProgram == nil {
		return messages
	}

	prog := fs.AffiliateProgram

	// Pay platform fee
	fs.Cash -= prog.MonthlyPlatformFee

	// Calculate affiliate sales (each affiliate brings 0-2 customers/month)
	newCustomers := 0
	for i := 0; i < prog.Affiliates; i++ {
		if rand.Float64() < 0.3 { // 30% chance per affiliate
			newCustomers += 1 + rand.Intn(2)
		}
	}

	if newCustomers > 0 {
		// Calculate MRR with variable deal sizes
		// Use template AvgDealSize if current AvgDealSize is 0 (no customers yet)
		baseDealSize := fs.AvgDealSize
		if baseDealSize == 0 {
			// Fallback to category-based defaults
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

		commissionPaid := int64(float64(totalMRR) * prog.Commission)

		// These are affiliate customers
		fs.Customers += newCustomers
		fs.AffiliateCustomers += newCustomers
		fs.AffiliateMRR += totalMRR
		fs.Cash -= commissionPaid

		// Add customers to tracking system
		for _, dealSize := range dealSizes {
			fs.addCustomer(dealSize, "affiliate")
		}

		prog.CustomersAcquired += newCustomers
		prog.MonthlyRevenue += totalMRR

		// Sync MRR from DirectMRR + AffiliateMRR
		fs.syncMRR()

		// Recalculate average deal size
		if fs.Customers > 0 {
			fs.AvgDealSize = fs.MRR / int64(fs.Customers)
		}

		messages = append(messages, fmt.Sprintf("🤝 Affiliates brought %d customers ($%s MRR, $%s commission)",
			newCustomers, formatCurrency(totalMRR), formatCurrency(commissionPaid)))

		// Affiliates grow over time if successful
		if rand.Float64() < 0.2 {
			prog.Affiliates += 1 + rand.Intn(3)
		}
	}

	return messages
}

// ============================================================================
// CHAIRMAN OF THE BOARD
// ============================================================================

// GetChairman returns the current chairman, if any
func (fs *FounderState) GetChairman() *BoardMember {
	for i := range fs.BoardMembers {
		if fs.BoardMembers[i].IsActive && fs.BoardMembers[i].IsChairman {
			return &fs.BoardMembers[i]
		}
	}
	return nil
}

// SetChairman promotes an advisor to chairman of the board
// Requires additional equity (1.5-2x) and higher monthly retainer
func (fs *FounderState) SetChairman(advisorName string) error {
	// Find the advisor
	var advisorIndex = -1
	for i := range fs.BoardMembers {
		if fs.BoardMembers[i].Name == advisorName && fs.BoardMembers[i].IsActive {
			advisorIndex = i
			break
		}
	}

	if advisorIndex == -1 {
		return fmt.Errorf("advisor not found: %s", advisorName)
	}

	advisor := &fs.BoardMembers[advisorIndex]

	// Check if already chairman
	if advisor.IsChairman {
		return fmt.Errorf("%s is already the chairman", advisorName)
	}

	// Remove any existing chairman
	currentChairman := fs.GetChairman()
	if currentChairman != nil {
		currentChairman.IsChairman = false
	}

	// Calculate additional equity cost (chairman needs 1.5-2x equity)
	additionalEquity := advisor.EquityCost * (0.5 + rand.Float64()*0.5) // 0.5-1x additional
	totalEquityNeeded := advisor.EquityCost + additionalEquity

	// Check if we have enough equity pool
	availableEquity := fs.EquityPool - fs.EquityAllocated
	if availableEquity < additionalEquity {
		return fmt.Errorf("insufficient equity pool (need %.2f%%, have %.2f%% available)", additionalEquity, availableEquity)
	}

	// Update equity allocation
	fs.EquityAllocated += additionalEquity
	advisor.EquityCost = totalEquityNeeded

	// Update cap table
	for i := range fs.CapTable {
		if fs.CapTable[i].Name == advisorName {
			fs.CapTable[i].Equity = totalEquityNeeded
			break
		}
	}

	// Set as chairman
	advisor.IsChairman = true

	return nil
}

// RemoveChairman removes the chairman role (with consequences)
func (fs *FounderState) RemoveChairman() error {
	chairman := fs.GetChairman()
	if chairman == nil {
		return fmt.Errorf("no chairman to remove")
	}

	// Remove chairman role
	chairman.IsChairman = false

	// Consequences: negative PR, board sentiment drop
	// This will be handled in UpdateBoardSentiment
	// Board pressure increases by 20-30 points
	if fs.BoardPressure < 100 {
		fs.BoardPressure += 20 + rand.Intn(11) // 20-30 point increase
		if fs.BoardPressure > 100 {
			fs.BoardPressure = 100
		}
	}

	// Board sentiment becomes more negative
	if fs.BoardSentiment == "happy" {
		fs.BoardSentiment = "pleased"
	} else if fs.BoardSentiment == "pleased" {
		fs.BoardSentiment = "neutral"
	} else if fs.BoardSentiment == "neutral" {
		fs.BoardSentiment = "concerned"
	} else if fs.BoardSentiment == "concerned" {
		fs.BoardSentiment = "angry"
	}

	return nil
}

// MitigateCrisis allows chairman to reduce severity of negative events
// Returns true if chairman successfully mitigated, false otherwise
func (fs *FounderState) MitigateCrisis(event *RandomEvent) bool {
	chairman := fs.GetChairman()
	if chairman == nil || event.IsPositive {
		return false
	}

	// Chairman can mitigate: legal, press, regulation events
	mitigatableTypes := []string{"legal", "press", "regulation"}
	isMitigatable := false
	for _, t := range mitigatableTypes {
		if event.Type == t {
			isMitigatable = true
			break
		}
	}

	if !isMitigatable {
		return false
	}

	// 70% chance chairman successfully mitigates
	if rand.Float64() < 0.7 {
		// Reduce impact by 30-50%
		mitigationFactor := 0.5 + rand.Float64()*0.2 // 0.5-0.7 (30-50% reduction)

		// Apply mitigation to impact
		if event.Impact.CashCost > 0 {
			event.Impact.CashCost = int64(float64(event.Impact.CashCost) * mitigationFactor)
		}
		if event.Impact.ChurnChange > 0 {
			event.Impact.ChurnChange *= mitigationFactor
		}
		if event.Impact.CACChange > 1.0 {
			// Reduce CAC increase
			event.Impact.CACChange = 1.0 + (event.Impact.CACChange-1.0)*mitigationFactor
		}
		if event.Impact.GrowthChange < 1.0 {
			// Reduce growth decrease
			event.Impact.GrowthChange = 1.0 - (1.0-event.Impact.GrowthChange)*mitigationFactor
		}
		if event.Impact.DurationMonths > 0 {
			event.Impact.DurationMonths = int(float64(event.Impact.DurationMonths) * mitigationFactor)
		}

		// Downgrade severity
		if event.Severity == "critical" {
			event.Severity = "major"
		} else if event.Severity == "major" {
			event.Severity = "moderate"
		} else if event.Severity == "moderate" {
			event.Severity = "minor"
		}

		return true
	}

	return false
}

// RemoveAdvisor removes an advisor from the board (with equity buyback option)
func (fs *FounderState) RemoveAdvisor(advisorName string, buybackEquity bool) error {
	// Find the advisor
	var advisorIndex = -1
	for i := range fs.BoardMembers {
		if fs.BoardMembers[i].Name == advisorName && fs.BoardMembers[i].IsActive {
			advisorIndex = i
			break
		}
	}

	if advisorIndex == -1 {
		return fmt.Errorf("advisor not found: %s", advisorName)
	}

	advisor := &fs.BoardMembers[advisorIndex]

	// Cannot remove chairman directly - must remove chairman role first
	if advisor.IsChairman {
		return fmt.Errorf("cannot remove chairman directly. Remove chairman role first, then remove advisor")
	}

	// Calculate buyback cost if requested
	if buybackEquity {
		// Buyback equity at current valuation (expensive)
		// Estimate valuation: MRR * 10-20x multiple (conservative)
		estimatedValuation := fs.MRR * 15 // 15x MRR multiple
		if estimatedValuation < 1000000 {
			estimatedValuation = 1000000 // Minimum $1M valuation
		}
		
		buybackCost := int64(float64(estimatedValuation) * (advisor.EquityCost / 100.0))
		
		if buybackCost > fs.Cash {
			return fmt.Errorf("insufficient cash for equity buyback (need $%s, have $%s)", 
				formatCurrency(buybackCost), formatCurrency(fs.Cash))
		}

		fs.Cash -= buybackCost
		
		// Return equity to founder (reduce EquityGivenAway)
		fs.EquityGivenAway -= advisor.EquityCost
		if fs.EquityGivenAway < 0 {
			fs.EquityGivenAway = 0
		}

		// Remove from cap table
		for i := range fs.CapTable {
			if fs.CapTable[i].Name == advisorName {
				fs.CapTable = append(fs.CapTable[:i], fs.CapTable[i+1:]...)
				break
			}
		}
	} else {
		// No buyback - advisor keeps equity but is removed from board
		// This causes negative PR and board sentiment issues
		if fs.BoardPressure < 100 {
			fs.BoardPressure += 10 + rand.Intn(11) // 10-20 point increase
			if fs.BoardPressure > 100 {
				fs.BoardPressure = 100
			}
		}

		// Board sentiment becomes more negative
		if fs.BoardSentiment == "happy" {
			fs.BoardSentiment = "pleased"
		} else if fs.BoardSentiment == "pleased" {
			fs.BoardSentiment = "neutral"
		} else if fs.BoardSentiment == "neutral" {
			fs.BoardSentiment = "concerned"
		}
	}

	// Mark advisor as inactive
	advisor.IsActive = false

	return nil
}

// FireBoardMember removes an investor board member (as majority owner)
// This has serious consequences - investors don't like being fired
func (fs *FounderState) FireBoardMember(memberName string) error {
	// Find the board member
	var memberIndex = -1
	for i := range fs.BoardMembers {
		if fs.BoardMembers[i].Name == memberName && fs.BoardMembers[i].IsActive {
			if fs.BoardMembers[i].Type == "investor" {
				memberIndex = i
				break
			}
		}
	}

	if memberIndex == -1 {
		return fmt.Errorf("investor board member not found: %s", memberName)
	}

	member := &fs.BoardMembers[memberIndex]

	// Check if founder has majority ownership (51%+)
	founderEquity := 100.0 - fs.EquityGivenAway - fs.EquityPool
	if founderEquity < 51.0 {
		return fmt.Errorf("cannot fire board member: you need 51%%+ ownership (you have %.1f%%)", founderEquity)
	}

	// Cannot fire chairman directly - must remove chairman role first
	if member.IsChairman {
		return fmt.Errorf("cannot fire chairman directly. Remove chairman role first, then fire board member")
	}

	// Serious consequences for firing an investor board member
	// Board pressure increases significantly
	if fs.BoardPressure < 100 {
		fs.BoardPressure += 30 + rand.Intn(21) // 30-50 point increase
		if fs.BoardPressure > 100 {
			fs.BoardPressure = 100
		}
	}

	// Board sentiment becomes very negative
	if fs.BoardSentiment == "happy" || fs.BoardSentiment == "pleased" {
		fs.BoardSentiment = "angry"
	} else if fs.BoardSentiment == "neutral" {
		fs.BoardSentiment = "angry"
	} else if fs.BoardSentiment == "concerned" {
		fs.BoardSentiment = "angry"
	}

	// Reduce board seats
	if fs.BoardSeats > 1 {
		fs.BoardSeats--
	}

	// Mark member as inactive
	member.IsActive = false

	return nil
}

// ============================================================================
// COMPETITORS
// ============================================================================
