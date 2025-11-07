package founder

import (
	"fmt"
	"math/rand"
)

// Company names for generating deals
var companyNames = []string{
	"Acme Corp", "TechVentures Inc", "DataCorp", "CloudSystems", "InnovateLabs",
	"GlobalTech", "FutureSoft", "MegaData", "SmartSolutions", "DigitalFirst",
	"AgileWorks", "NextGen Systems", "Velocity Corp", "Quantum Inc", "Synergy Solutions",
	"Enterprise Hub", "Prodigy Tech", "Visionary Corp", "Apex Solutions", "Pioneer Systems",
}

// InitializeSalesPipeline sets up the sales pipeline
func (fs *FounderState) InitializeSalesPipeline() {
	fs.SalesPipeline = &SalesPipeline{
		ActiveDeals:       []Deal{},
		ClosedDeals:       []Deal{},
		LeadsPerMonth:     0,
		ConversionRates:   make(map[string]float64),
		AverageDealSize:   fs.AvgDealSize,
		AverageSalesCycle: 90, // days
		WinRate:           0.20,
		TotalDealsCreated: 0,
		NextDealID:        1,
	}

	// Initialize conversion rates
	fs.SalesPipeline.ConversionRates["lead"] = 0.30       // 30% of leads qualify
	fs.SalesPipeline.ConversionRates["qualified"] = 0.40  // 40% move to demo
	fs.SalesPipeline.ConversionRates["demo"] = 0.50       // 50% move to negotiation
	fs.SalesPipeline.ConversionRates["negotiation"] = 0.40 // 40% close won
}

// GenerateNewDeals creates new sales opportunities based on team size and growth
func (fs *FounderState) GenerateNewDeals() []string {
	var messages []string

	if fs.SalesPipeline == nil {
		fs.InitializeSalesPipeline()
	}

	// Calculate lead generation based on sales team + marketing + content
	salesTeam := len(fs.Team.Sales)
	leadMultiplier := 1.0

	// Content marketing boosts inbound leads
	if fs.ContentProgram != nil && fs.ContentProgram.InboundLeads > 0 {
		leadMultiplier += 0.5
	}

	// Calculate new leads
	baseLeads := salesTeam * 3 // Each sales rep generates ~3 leads per month
	if baseLeads < 2 {
		baseLeads = 2 // Minimum 2 leads even with no sales team (founder selling)
	}

	newLeads := int(float64(baseLeads) * leadMultiplier)
	fs.SalesPipeline.LeadsPerMonth = newLeads

	// Generate new deal opportunities
	for i := 0; i < newLeads; i++ {
		// Determine segment for this deal
		segment := fs.SuggestSegmentForNewCustomer()
		if segment == "" {
			segment = "SMB" // Default
		}

		// Generate deal size appropriate for segment
		dealSize := fs.GenerateDealSizeForSegment(segment)

		// Pick random vertical if we have a focus
		vertical := ""
		if fs.SelectedVertical != "" {
			// 70% chance to get focused vertical
			if rand.Float64() < 0.70 {
				vertical = fs.SelectedVertical
			}
		}

		// Create the deal
		deal := Deal{
			ID:               fs.SalesPipeline.NextDealID,
			CompanyName:      companyNames[rand.Intn(len(companyNames))],
			DealSize:         dealSize,
			Stage:            "lead",
			CloseProbability: 0.05 + rand.Float64()*0.10, // 5-15% initial probability
			DaysInStage:      0,
			RequiredActions:  []string{"Qualify lead", "Schedule discovery call"},
			AssignedSalesRep: "",
			MonthCreated:     fs.Turn,
			Segment:          segment,
			Vertical:         vertical,
		}

		// Assign to sales rep if available
		if len(fs.Team.Sales) > 0 {
			deal.AssignedSalesRep = fs.Team.Sales[rand.Intn(len(fs.Team.Sales))].Name
		} else {
			deal.AssignedSalesRep = fs.FounderName + " (Founder)"
		}

		fs.SalesPipeline.ActiveDeals = append(fs.SalesPipeline.ActiveDeals, deal)
		fs.SalesPipeline.NextDealID++
		fs.SalesPipeline.TotalDealsCreated++
	}

	if newLeads > 0 {
		messages = append(messages, fmt.Sprintf("📞 Generated %d new sales leads this month", newLeads))
	}

	return messages
}

// ProcessPipeline moves deals through the sales funnel
func (fs *FounderState) ProcessPipeline() []string {
	var messages []string

	if fs.SalesPipeline == nil {
		return messages
	}

	totalClosed := 0
	totalWon := 0
	totalRevenue := int64(0)

	// Process each active deal
	for i := len(fs.SalesPipeline.ActiveDeals) - 1; i >= 0; i-- {
		deal := &fs.SalesPipeline.ActiveDeals[i]
		
		// Age the deal
		deal.DaysInStage += 30 // One month = ~30 days

		// Calculate progression probability based on stage and deal attributes
		progressionChance := fs.SalesPipeline.ConversionRates[deal.Stage]

		// Apply bonuses from features, ICP match, etc.
		if fs.ProductRoadmap != nil {
			_, closeRateBonus, _ := fs.GetFeatureBonuses()
			progressionChance *= (1.0 + closeRateBonus)
		}

		if fs.SelectedICP != "" && deal.Segment == fs.SelectedICP {
			_, icpCloseBonus, _ := fs.GetICPBenefits()
			progressionChance *= (1.0 + icpCloseBonus)
		}

		// Check if deal progresses
		if rand.Float64() < progressionChance {
			switch deal.Stage {
			case "lead":
				deal.Stage = "qualified"
				deal.CloseProbability = 0.15 + rand.Float64()*0.15 // 15-30%
				deal.DaysInStage = 0
				deal.RequiredActions = []string{"Conduct discovery call", "Send proposal"}

			case "qualified":
				deal.Stage = "demo"
				deal.CloseProbability = 0.30 + rand.Float64()*0.20 // 30-50%
				deal.DaysInStage = 0
				deal.RequiredActions = []string{"Demo product", "Address objections"}

			case "demo":
				deal.Stage = "negotiation"
				deal.CloseProbability = 0.50 + rand.Float64()*0.30 // 50-80%
				deal.DaysInStage = 0
				deal.RequiredActions = []string{"Negotiate terms", "Send contract"}

			case "negotiation":
				// Close the deal!
				if rand.Float64() < deal.CloseProbability {
					// Won!
					deal.Stage = "closed_won"
					totalClosed++
					totalWon++
					totalRevenue += deal.DealSize

					// Add as customer
					customer := Customer{
						ID:           fs.NextCustomerID,
						Source:       "direct",
						DealSize:     deal.DealSize,
						IsActive:     true,
						MonthChurned: 0,
						MonthAdded:   fs.Turn,
						HealthScore:  0.80, // New customers start reasonably healthy
					}
					fs.CustomerList = append(fs.CustomerList, customer)
					fs.NextCustomerID++
					fs.Customers++
					fs.DirectCustomers++
					fs.TotalCustomersEver++
					fs.MRR += deal.DealSize
					fs.DirectMRR += deal.DealSize
				} else {
					// Lost
					deal.Stage = "closed_lost"
					deal.LostReason = getRandomLostReason()
					totalClosed++
				}

				// Move to closed deals
				fs.SalesPipeline.ClosedDeals = append(fs.SalesPipeline.ClosedDeals, *deal)
				fs.SalesPipeline.ActiveDeals = append(fs.SalesPipeline.ActiveDeals[:i], fs.SalesPipeline.ActiveDeals[i+1:]...)
			}
		} else {
			// Deal didn't progress - check if it's too old and should be marked lost
			maxDaysInStage := 120 // 4 months max per stage
			if deal.DaysInStage > maxDaysInStage {
				deal.Stage = "closed_lost"
				deal.LostReason = "No response / Deal stalled"
				fs.SalesPipeline.ClosedDeals = append(fs.SalesPipeline.ClosedDeals, *deal)
				fs.SalesPipeline.ActiveDeals = append(fs.SalesPipeline.ActiveDeals[:i], fs.SalesPipeline.ActiveDeals[i+1:]...)
				totalClosed++
			}
		}
	}

	// Update win rate
	if len(fs.SalesPipeline.ClosedDeals) > 0 {
		wonCount := 0
		for _, deal := range fs.SalesPipeline.ClosedDeals {
			if deal.Stage == "closed_won" {
				wonCount++
			}
		}
		fs.SalesPipeline.WinRate = float64(wonCount) / float64(len(fs.SalesPipeline.ClosedDeals))
	}

	// Report closed deals
	if totalWon > 0 {
		messages = append(messages, fmt.Sprintf("🎉 Closed %d deals worth $%s/month in new MRR!", totalWon, formatCurrency(totalRevenue)))
	}
	if totalClosed > totalWon {
		lostCount := totalClosed - totalWon
		messages = append(messages, fmt.Sprintf("📉 Lost %d deals this month", lostCount))
	}

	return messages
}

// GetPipelineMetrics returns key pipeline statistics
func (fs *FounderState) GetPipelineMetrics() map[string]interface{} {
	if fs.SalesPipeline == nil {
		return nil
	}

	// Count deals by stage
	stageCounts := map[string]int{
		"lead":        0,
		"qualified":   0,
		"demo":        0,
		"negotiation": 0,
	}

	totalValue := int64(0)
	for _, deal := range fs.SalesPipeline.ActiveDeals {
		stageCounts[deal.Stage]++
		totalValue += deal.DealSize
	}

	// Calculate weighted pipeline value (probability * deal size)
	weightedValue := int64(0)
	for _, deal := range fs.SalesPipeline.ActiveDeals {
		weightedValue += int64(float64(deal.DealSize) * deal.CloseProbability)
	}

	return map[string]interface{}{
		"activeDeals":    len(fs.SalesPipeline.ActiveDeals),
		"stageCounts":    stageCounts,
		"pipelineValue":  totalValue,
		"weightedValue":  weightedValue,
		"winRate":        fs.SalesPipeline.WinRate,
		"leadsPerMonth":  fs.SalesPipeline.LeadsPerMonth,
		"avgDealSize":    fs.SalesPipeline.AverageDealSize,
		"avgSalesCycle":  fs.SalesPipeline.AverageSalesCycle,
	}
}

// GetDealsByStage returns deals filtered by stage
func (fs *FounderState) GetDealsByStage(stage string) []Deal {
	if fs.SalesPipeline == nil {
		return []Deal{}
	}

	deals := []Deal{}
	for _, deal := range fs.SalesPipeline.ActiveDeals {
		if deal.Stage == stage {
			deals = append(deals, deal)
		}
	}
	return deals
}

// GetWinLossReasons returns a summary of why deals are lost
func (fs *FounderState) GetWinLossReasons() map[string]int {
	if fs.SalesPipeline == nil {
		return map[string]int{}
	}

	reasons := make(map[string]int)
	for _, deal := range fs.SalesPipeline.ClosedDeals {
		if deal.Stage == "closed_lost" && deal.LostReason != "" {
			reasons[deal.LostReason]++
		}
	}
	return reasons
}

// AccelerateDeals allows spending money to speed up deals
func (fs *FounderState) AccelerateDeal(dealID int, action string) error {
	if fs.SalesPipeline == nil {
		return fmt.Errorf("no sales pipeline initialized")
	}

	// Find the deal
	var deal *Deal
	for i := range fs.SalesPipeline.ActiveDeals {
		if fs.SalesPipeline.ActiveDeals[i].ID == dealID {
			deal = &fs.SalesPipeline.ActiveDeals[i]
			break
		}
	}

	if deal == nil {
		return fmt.Errorf("deal not found")
	}

	// Apply action
	var cost int64
	var probabilityIncrease float64

	switch action {
	case "demo":
		cost = 5000 + rand.Int63n(5000) // $5-10k
		probabilityIncrease = 0.10      // +10% close probability
	case "poc":
		cost = 20000 + rand.Int63n(30000) // $20-50k
		probabilityIncrease = 0.20        // +20% close probability
	case "travel":
		cost = 2000 + rand.Int63n(3000) // $2-5k
		probabilityIncrease = 0.05      // +5% close probability
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	if fs.Cash < cost {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	// Pay cost and apply benefit
	fs.Cash -= cost
	deal.CloseProbability += probabilityIncrease
	if deal.CloseProbability > 0.95 {
		deal.CloseProbability = 0.95 // Cap at 95%
	}

	return nil
}

// getRandomLostReason returns a random reason for losing a deal
func getRandomLostReason() string {
	reasons := []string{
		"Price too high",
		"Chose competitor",
		"Missing features",
		"Budget constraints",
		"Timing not right",
		"Security concerns",
		"Integration issues",
		"No decision maker access",
		"Lost to status quo",
		"Product not mature enough",
	}
	return reasons[rand.Intn(len(reasons))]
}

