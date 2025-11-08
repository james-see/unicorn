package founder

import (
	"fmt"
	"math"
	"math/rand"
)

// InitializeProductRoadmap creates the roadmap with available features
func (fs *FounderState) InitializeProductRoadmap() {
	fs.ProductRoadmap = &ProductRoadmap{
		Features:          []ProductFeature{},
		AvailableFeatures: getAvailableFeatures(),
		CompletedCount:    0,
		InProgressCount:   0,
		CompetitorLaunches: []CompetitorFeatureLaunch{},
	}
}

// getAvailableFeatures returns the template of features that can be built
func getAvailableFeatures() []ProductFeature {
	return []ProductFeature{
		{
			Name:              "REST API",
			Category:          "Integration",
			EngineerMonths:    4,
			Cost:              50000,
			ChurnReduction:    0.02,
			CloseRateIncrease: 0.10,
			DealSizeIncrease:  0.15,
			MarketAppealScore: 85,
			Status:            "available",
		},
		{
			Name:              "Mobile App",
			Category:          "Platform",
			EngineerMonths:    8,
			Cost:              120000,
			ChurnReduction:    0.03,
			CloseRateIncrease: 0.15,
			DealSizeIncrease:  0.20,
			MarketAppealScore: 90,
			Status:            "available",
		},
		{
			Name:              "Enterprise SSO",
			Category:          "Security",
			EngineerMonths:    3,
			Cost:              40000,
			ChurnReduction:    0.03,
			CloseRateIncrease: 0.20,
			DealSizeIncrease:  0.25,
			MarketAppealScore: 95,
			Status:            "available",
		},
		{
			Name:              "Advanced Analytics",
			Category:          "Analytics",
			EngineerMonths:    5,
			Cost:              70000,
			ChurnReduction:    0.02,
			CloseRateIncrease: 0.12,
			DealSizeIncrease:  0.18,
			MarketAppealScore: 80,
			Status:            "available",
		},
		{
			Name:              "AI/ML Capabilities",
			Category:          "AI",
			EngineerMonths:    10,
			Cost:              200000,
			ChurnReduction:    0.04,
			CloseRateIncrease: 0.25,
			DealSizeIncrease:  0.35,
			MarketAppealScore: 100,
			Status:            "available",
		},
		{
			Name:              "Integrations Hub",
			Category:          "Integration",
			EngineerMonths:    6,
			Cost:              90000,
			ChurnReduction:    0.04,
			CloseRateIncrease: 0.15,
			DealSizeIncrease:  0.22,
			MarketAppealScore: 88,
			Status:            "available",
		},
		{
			Name:              "Security Suite",
			Category:          "Security",
			EngineerMonths:    4,
			Cost:              60000,
			ChurnReduction:    0.02,
			CloseRateIncrease: 0.10,
			DealSizeIncrease:  0.15,
			MarketAppealScore: 82,
			Status:            "available",
		},
		{
			Name:              "Performance Optimization",
			Category:          "Infrastructure",
			EngineerMonths:    3,
			Cost:              50000,
			ChurnReduction:    0.03,
			CloseRateIncrease: 0.05,
			DealSizeIncrease:  0.10,
			MarketAppealScore: 70,
			Status:            "available",
		},
		{
			Name:              "White Label",
			Category:          "Enterprise",
			EngineerMonths:    7,
			Cost:              110000,
			ChurnReduction:    0.02,
			CloseRateIncrease: 0.18,
			DealSizeIncrease:  0.30,
			MarketAppealScore: 85,
			Status:            "available",
		},
		{
			Name:              "Workflow Automation",
			Category:          "Product",
			EngineerMonths:    5,
			Cost:              75000,
			ChurnReduction:    0.03,
			CloseRateIncrease: 0.12,
			DealSizeIncrease:  0.20,
			MarketAppealScore: 78,
			Status:            "available",
		},
	}
}

// StartFeature begins development on a feature
func (fs *FounderState) StartFeature(featureName string, engineers int) error {
	if fs.ProductRoadmap == nil {
		fs.InitializeProductRoadmap()
	}

	// Check if enough engineers available
	allocatedEngineers := fs.GetAllocatedEngineers()
	availableEngineers := len(fs.Team.Engineers) - allocatedEngineers
	if engineers > availableEngineers {
		return fmt.Errorf("only %d engineers available (need %d)", availableEngineers, engineers)
	}

	// Find the feature in available features
	var featureTemplate *ProductFeature
	for i := range fs.ProductRoadmap.AvailableFeatures {
		if fs.ProductRoadmap.AvailableFeatures[i].Name == featureName {
			featureTemplate = &fs.ProductRoadmap.AvailableFeatures[i]
			break
		}
	}

	if featureTemplate == nil {
		return fmt.Errorf("feature not found: %s", featureName)
	}

	// Check if already in progress or completed
	for _, f := range fs.ProductRoadmap.Features {
		if f.Name == featureName && (f.Status == "in_progress" || f.Status == "completed") {
			return fmt.Errorf("feature already %s", f.Status)
		}
	}

	// Check if have enough cash
	if fs.Cash < featureTemplate.Cost {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(featureTemplate.Cost))
	}

	// Pay the cost and start feature
	fs.Cash -= featureTemplate.Cost

	newFeature := *featureTemplate
	newFeature.Status = "in_progress"
	newFeature.MonthStarted = fs.Turn
	newFeature.DevelopmentProgress = 0
	newFeature.AllocatedEngineers = engineers

	fs.ProductRoadmap.Features = append(fs.ProductRoadmap.Features, newFeature)
	fs.ProductRoadmap.InProgressCount++

	return nil
}

// GetAllocatedEngineers returns the number of engineers currently allocated to features
func (fs *FounderState) GetAllocatedEngineers() int {
	if fs.ProductRoadmap == nil {
		return 0
	}

	allocated := 0
	for _, feature := range fs.ProductRoadmap.Features {
		if feature.Status == "in_progress" {
			allocated += feature.AllocatedEngineers
		}
	}
	return allocated
}

// ReallocateEngineers changes the number of engineers on a feature
func (fs *FounderState) ReallocateEngineers(featureName string, newCount int) error {
	if fs.ProductRoadmap == nil {
		return fmt.Errorf("no product roadmap initialized")
	}

	// Find the feature
	featureIndex := -1
	for i := range fs.ProductRoadmap.Features {
		if fs.ProductRoadmap.Features[i].Name == featureName && fs.ProductRoadmap.Features[i].Status == "in_progress" {
			featureIndex = i
			break
		}
	}

	if featureIndex == -1 {
		return fmt.Errorf("feature not found or not in progress: %s", featureName)
	}

	feature := &fs.ProductRoadmap.Features[featureIndex]
	currentAllocation := feature.AllocatedEngineers

	// Calculate available engineers excluding current allocation
	allocatedEngineers := fs.GetAllocatedEngineers() - currentAllocation
	availableEngineers := len(fs.Team.Engineers) - allocatedEngineers

	if newCount > availableEngineers {
		return fmt.Errorf("only %d engineers available", availableEngineers)
	}

	feature.AllocatedEngineers = newCount
	return nil
}

// ProcessRoadmapProgress updates feature development each month
func (fs *FounderState) ProcessRoadmapProgress() []string {
	var messages []string

	if fs.ProductRoadmap == nil {
		return messages
	}

	for i := range fs.ProductRoadmap.Features {
		feature := &fs.ProductRoadmap.Features[i]
		
		if feature.Status != "in_progress" {
			continue
		}

		// Calculate progress based on engineers allocated
		// Each engineer contributes ~25% progress per month per engineer-month needed
		if feature.AllocatedEngineers > 0 {
			progressPerMonth := (float64(feature.AllocatedEngineers) / float64(feature.EngineerMonths)) * 100.0
			feature.DevelopmentProgress += int(progressPerMonth)

			if feature.DevelopmentProgress >= 100 {
				feature.DevelopmentProgress = 100
				feature.Status = "completed"
				feature.MonthCompleted = fs.Turn
				fs.ProductRoadmap.CompletedCount++
				fs.ProductRoadmap.InProgressCount--

				// Check if competitor launched this feature before player completed it
				// Map feature names to competitor feature names (normalize for comparison)
				featureMap := map[string][]string{
					"REST API": {"API", "REST API", "API Integration"},
					"Mobile App": {"Mobile App", "Mobile", "iOS App", "Android App"},
					"Enterprise SSO": {"SSO", "Single Sign-On", "Enterprise SSO"},
					"Advanced Analytics": {"Analytics", "Advanced Analytics", "Reporting"},
					"AI/ML Capabilities": {"AI", "ML", "Machine Learning", "Artificial Intelligence"},
					"Integrations Hub": {"Integrations", "Integration Hub", "API Integration"},
					"Security Suite": {"Security", "Security Suite", "Enterprise Security"},
				}
				
				// Check if any competitor launched a similar feature before completion
				competitorLaunchedFirst := false
				for _, launch := range fs.ProductRoadmap.CompetitorLaunches {
					// Check if this competitor launch matches the completed feature
					matchingNames, exists := featureMap[feature.Name]
					if !exists {
						// Try direct match
						matchingNames = []string{feature.Name}
					}
					for _, matchName := range matchingNames {
						if launch.FeatureName == matchName && launch.MonthLaunched < feature.MonthCompleted {
							competitorLaunchedFirst = true
							break
						}
					}
					if competitorLaunchedFirst {
						break
					}
				}

				// Apply benefits permanently
				fs.CustomerChurnRate = math.Max(0.01, fs.CustomerChurnRate-feature.ChurnReduction)

				if competitorLaunchedFirst {
					messages = append(messages, fmt.Sprintf("✅ FEATURE COMPLETE: %s! (Competitors launched similar features earlier) Churn -%.1f%%, Close Rate +%.0f%%, Deal Size +%.0f%%",
						feature.Name,
						feature.ChurnReduction*100,
						feature.CloseRateIncrease*100,
						feature.DealSizeIncrease*100))
				} else {
					messages = append(messages, fmt.Sprintf("✅ FEATURE COMPLETE: %s! 🏆 Innovation Leader! Churn -%.1f%%, Close Rate +%.0f%%, Deal Size +%.0f%%",
						feature.Name,
						feature.ChurnReduction*100,
						feature.CloseRateIncrease*100,
						feature.DealSizeIncrease*100))
				}
			} else {
				messages = append(messages, fmt.Sprintf("🔨 %s: %d%% complete (%d engineers working)",
					feature.Name, feature.DevelopmentProgress, feature.AllocatedEngineers))
			}
		}
	}

	// Check if competitors launch features (creates pressure)
	if rand.Float64() < 0.08 { // 8% chance per month
		competitorFeatures := []string{"API", "Mobile App", "SSO", "Analytics", "Integrations", "Security"}
		feature := competitorFeatures[rand.Intn(len(competitorFeatures))]
		
		compName := "A competitor"
		if len(fs.Competitors) > 0 {
			activeComps := []Competitor{}
			for _, c := range fs.Competitors {
				if c.Active {
					activeComps = append(activeComps, c)
				}
			}
			if len(activeComps) > 0 {
				compName = activeComps[rand.Intn(len(activeComps))].Name
			}
		}

		// Track competitor launch for innovation leader tracking
		if fs.ProductRoadmap == nil {
			fs.InitializeProductRoadmap()
		}
		launch := CompetitorFeatureLaunch{
			FeatureName:   feature,
			CompetitorName: compName,
			MonthLaunched: fs.Turn,
		}
		fs.ProductRoadmap.CompetitorLaunches = append(fs.ProductRoadmap.CompetitorLaunches, launch)

		messages = append(messages, fmt.Sprintf("⚠️  %s launched a new %s feature! Consider your product roadmap", compName, feature))
	}

	return messages
}

// GetAvailableFeaturesToStart returns features that can be started
func (fs *FounderState) GetAvailableFeaturesToStart() []ProductFeature {
	if fs.ProductRoadmap == nil {
		fs.InitializeProductRoadmap()
	}

	available := []ProductFeature{}
	
	// Create a map of features already started or completed
	started := make(map[string]bool)
	for _, f := range fs.ProductRoadmap.Features {
		if f.Status == "in_progress" || f.Status == "completed" {
			started[f.Name] = true
		}
	}

	// Return features not yet started
	for _, f := range fs.ProductRoadmap.AvailableFeatures {
		if !started[f.Name] {
			available = append(available, f)
		}
	}

	return available
}

// GetInProgressFeatures returns features currently being built
func (fs *FounderState) GetInProgressFeatures() []ProductFeature {
	if fs.ProductRoadmap == nil {
		return []ProductFeature{}
	}

	inProgress := []ProductFeature{}
	for _, f := range fs.ProductRoadmap.Features {
		if f.Status == "in_progress" {
			inProgress = append(inProgress, f)
		}
	}
	return inProgress
}

// GetCompletedFeatures returns features that are done
func (fs *FounderState) GetCompletedFeatures() []ProductFeature {
	if fs.ProductRoadmap == nil {
		return []ProductFeature{}
	}

	completed := []ProductFeature{}
	for _, f := range fs.ProductRoadmap.Features {
		if f.Status == "completed" {
			completed = append(completed, f)
		}
	}
	return completed
}

// GetFeatureBonuses calculates cumulative bonuses from completed features
func (fs *FounderState) GetFeatureBonuses() (churnReduction float64, closeRateBonus float64, dealSizeBonus float64) {
	if fs.ProductRoadmap == nil {
		return 0, 0, 0
	}

	for _, f := range fs.ProductRoadmap.Features {
		if f.Status == "completed" {
			churnReduction += f.ChurnReduction
			closeRateBonus += f.CloseRateIncrease
			dealSizeBonus += f.DealSizeIncrease
		}
	}
	return
}

// CancelFeature stops development on a feature (forfeit cost)
func (fs *FounderState) CancelFeature(featureName string) error {
	if fs.ProductRoadmap == nil {
		return fmt.Errorf("no product roadmap initialized")
	}

	featureIndex := -1
	for i := range fs.ProductRoadmap.Features {
		if fs.ProductRoadmap.Features[i].Name == featureName && fs.ProductRoadmap.Features[i].Status == "in_progress" {
			featureIndex = i
			break
		}
	}

	if featureIndex == -1 {
		return fmt.Errorf("feature not found or not in progress: %s", featureName)
	}

	// Remove the feature
	fs.ProductRoadmap.Features = append(fs.ProductRoadmap.Features[:featureIndex], fs.ProductRoadmap.Features[featureIndex+1:]...)
	fs.ProductRoadmap.InProgressCount--

	return nil
}

