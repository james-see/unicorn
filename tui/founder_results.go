package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/achievements"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/founder"
	"github.com/jamesacampbell/unicorn/leaderboard"
	"github.com/jamesacampbell/unicorn/progression"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// FounderResultsScreen shows the end-game results for founder mode
type FounderResultsScreen struct {
	width    int
	height   int
	gameData *GameData
	phase    ResultsPhase

	// Computed results
	outcome       string
	valuation     int64
	founderEquity float64
	founderPayout int64
	roi           float64
	rating        string
	ratingIcon    string

	// New achievements unlocked
	newAchievements    []string
	newAchievementObjs []achievements.Achievement

	// XP and Level tracking
	xpBreakdown   map[string]int
	totalXP       int
	leveledUp     bool
	oldLevel      int
	newLevel      int
	levelUpPoints int
	profileBefore *database.PlayerProfile
	profileAfter  *database.PlayerProfile

	// Score saved
	scoreSaved bool

	// Leaderboard submission
	leaderboardSubmitted bool
	leaderboardError     string
}

// NewFounderResultsScreen creates a new founder results screen
func NewFounderResultsScreen(width, height int, gameData *GameData) *FounderResultsScreen {
	fs := gameData.FounderState

	outcome, valuation, founderEquity := fs.GetFinalScore()

	// Calculate founder payout
	var founderPayout int64
	if fs.HasExited {
		switch fs.ExitType {
		case "ipo":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.20 / 100.0)
		case "acquisition":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity / 100.0)
		case "secondary":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.5 / 100.0)
		default:
			founderPayout = fs.Cash + int64(float64(valuation)*founderEquity/100.0)
		}
	} else {
		founderPayout = fs.Cash + int64(float64(valuation)*founderEquity/100.0)
	}

	// Calculate ROI
	initialCash := int64(500000)
	if len(fs.FundingRounds) > 0 {
		initialCash = fs.FundingRounds[0].Amount * 2
	}
	roi := 0.0
	if initialCash > 0 {
		roi = (float64(founderPayout-initialCash) / float64(initialCash)) * 100.0
	}

	rating, ratingIcon := calculateFounderRating(fs, valuation, founderPayout)

	return &FounderResultsScreen{
		width:         width,
		height:        height,
		gameData:      gameData,
		phase:         PhaseResults,
		outcome:       outcome,
		valuation:     valuation,
		founderEquity: founderEquity,
		founderPayout: founderPayout,
		roi:           roi,
		rating:        rating,
		ratingIcon:    ratingIcon,
	}
}

func calculateFounderRating(fs *founder.FounderState, valuation int64, founderPayout int64) (string, string) {
	if fs.HasExited && fs.ExitType == "ipo" && valuation >= 1000000000 {
		return "UNICORN FOUNDER - Legendary!", "👑"
	}
	if fs.HasExited && fs.ExitType == "ipo" {
		return "IPO SUCCESS - Outstanding!", "🏆"
	}
	if founderPayout >= 100000000 {
		return "Mega Exit - Excellent!", "⭐"
	}
	if fs.HasExited && founderPayout >= 10000000 {
		return "Great Exit - Well Done!", "🎯"
	}
	if fs.HasExited {
		return "Successful Exit - Good Job!", "✓"
	}
	if fs.MRR >= 1000000 {
		return "Unicorn Path - Impressive!", "🦄"
	}
	if fs.MRR >= 100000 {
		return "Real Business - Solid!", "✓"
	}
	if fs.Cash <= 0 {
		return "Ran Out of Cash - Try Again", "⚠"
	}
	if fs.MRR > 0 {
		return "Survived - Keep Going", "="
	}
	return "No Traction - Better Luck Next Time", "⚠"
}

// Init initializes the results screen
func (s *FounderResultsScreen) Init() tea.Cmd {
	fs := s.gameData.FounderState

	// Calculate total funding raised
	totalFundingRaised := int64(0)
	for _, round := range fs.FundingRounds {
		totalFundingRaised += round.Amount
	}

	// Calculate successful exits
	successfulExits := 0
	if fs.HasExited {
		successfulExits = 1
	}

	// Calculate actual months played (Turn gets inflated on exit)
	actualTurns := fs.Turn
	if fs.HasExited && fs.ExitMonth > 0 {
		actualTurns = fs.ExitMonth
	}

	// Save score to database
	score := database.GameScore{
		PlayerName:      fs.FounderName,
		FinalNetWorth:   s.founderPayout,
		ROI:             s.roi,
		SuccessfulExits: successfulExits,
		TurnsPlayed:     actualTurns,
		Difficulty:      "Founder",
		PlayedAt:        time.Now(),
	}

	err := database.SaveGameScore(score)
	if err == nil {
		s.scoreSaved = true
	}

	// Get profile before XP is added
	s.profileBefore, _ = database.GetPlayerProfile(fs.FounderName)
	if s.profileBefore != nil {
		s.oldLevel = s.profileBefore.Level
	} else {
		s.oldLevel = 1
	}

	// Check for achievements
	s.checkAchievements()

	// Calculate XP breakdown
	s.calculateXPBreakdown()

	// Add XP to player
	var levelUpErr error
	s.leveledUp, s.newLevel, s.levelUpPoints, levelUpErr = database.AddExperience(fs.FounderName, s.totalXP)
	if levelUpErr == nil {
		s.profileAfter, _ = database.GetPlayerProfile(fs.FounderName)
	}

	return nil
}

func (s *FounderResultsScreen) calculateXPBreakdown() {
	fs := s.gameData.FounderState
	s.xpBreakdown = make(map[string]int)

	// Base XP for completing a game
	s.xpBreakdown["Game Completion"] = progression.XPGameComplete
	s.totalXP = progression.XPGameComplete

	// Positive ROI bonus
	if s.roi > 0 {
		s.xpBreakdown["Positive ROI"] = progression.XPPositiveROI
		s.totalXP += progression.XPPositiveROI
	}

	// Successful exit bonus
	if fs.HasExited {
		exitXP := progression.XPSuccessfulExit
		s.xpBreakdown["Successful Exit"] = exitXP
		s.totalXP += exitXP

		// Founder-specific exit bonuses
		if fs.ExitType == "ipo" {
			s.xpBreakdown["IPO Exit Bonus"] = 500
			s.totalXP += 500
		} else if fs.ExitType == "acquisition" {
			s.xpBreakdown["Acquisition Exit Bonus"] = 300
			s.totalXP += 300
		}
	}

	// Profitability bonus
	if fs.MonthReachedProfitability > 0 && fs.MonthReachedProfitability <= 24 {
		s.xpBreakdown["Reached Profitability"] = 100
		s.totalXP += 100
	}

	// Achievement bonuses
	if len(s.newAchievementObjs) > 0 {
		achXP := 0
		for _, ach := range s.newAchievementObjs {
			achXP += ach.Points * progression.XPAchievementBase
		}
		if achXP > 0 {
			s.xpBreakdown[fmt.Sprintf("New Achievements (%d)", len(s.newAchievementObjs))] = achXP
			s.totalXP += achXP
		}
	}
}

func (s *FounderResultsScreen) checkAchievements() {
	fs := s.gameData.FounderState

	// Calculate total funding raised
	totalFundingRaised := int64(0)
	for _, round := range fs.FundingRounds {
		totalFundingRaised += round.Amount
	}

	successfulExits := 0
	if fs.HasExited {
		successfulExits = 1
	}

	// Phase 1 stats
	featuresCompleted := 0
	enterpriseFeatures := 0
	if fs.ProductRoadmap != nil {
		featuresCompleted = fs.ProductRoadmap.CompletedCount
		for _, f := range fs.GetCompletedFeatures() {
			if f.Category == "Security" || f.Category == "Analytics" {
				enterpriseFeatures++
			}
		}
	}

	enterpriseCustomers := 0
	if len(fs.CustomerSegments) > 0 {
		for _, seg := range fs.CustomerSegments {
			if seg.Name == "Enterprise" {
				enterpriseCustomers = seg.Volume
			}
		}
	}

	verticalConcentration := fs.GetSegmentConcentration()

	pricingExperiments := 0
	if fs.PricingStrategy != nil {
		for _, change := range fs.PricingStrategy.ChangeHistory {
			if strings.Contains(change.Reason, "Applied experiment") {
				pricingExperiments++
			}
		}
	}

	premiumPricing := fs.AvgDealSize > 50000 && fs.MonthlyGrowthRate > 0.05

	lowTouchCustomers := 0
	if len(fs.CustomerSegments) > 0 {
		for _, seg := range fs.CustomerSegments {
			if seg.Name == "SMB" || seg.Name == "Startup" {
				lowTouchCustomers += seg.Volume
			}
		}
	}

	dealsClosedWon := 0
	highProbClose := false
	maxPipelineSize := 0
	if fs.SalesPipeline != nil {
		for _, deal := range fs.SalesPipeline.ClosedDeals {
			if deal.Stage == "closed_won" {
				dealsClosedWon++
				if deal.CloseProbability >= 0.90 {
					highProbClose = true
				}
			}
		}
		if len(fs.SalesPipeline.ActiveDeals) > maxPipelineSize {
			maxPipelineSize = len(fs.SalesPipeline.ActiveDeals)
		}
	}

	// Get player stats
	playerStats, _ := database.GetPlayerStats(fs.FounderName)
	winStreak, _ := database.GetWinStreak(fs.FounderName)

	// Calculate actual turns for achievements
	achievementTurns := fs.Turn
	if fs.HasExited && fs.ExitMonth > 0 {
		achievementTurns = fs.ExitMonth
	}

	// Build game stats for achievement checking
	gameStats := achievements.GameStats{
		GameMode:                    "founder",
		FinalNetWorth:               s.founderPayout,
		ROI:                         s.roi,
		SuccessfulExits:             successfulExits,
		TurnsPlayed:                 achievementTurns,
		Difficulty:                  "Founder",
		FinalMRR:                    fs.MRR,
		FinalValuation:              s.valuation,
		FinalEquity:                 s.founderEquity,
		Customers:                   fs.Customers,
		FundingRoundsRaised:         len(fs.FundingRounds),
		TotalFundingRaised:          totalFundingRaised,
		HasExited:                   fs.HasExited,
		ExitType:                    fs.ExitType,
		ExitValuation:               fs.ExitValuation,
		MonthsToProfitability:       fs.MonthReachedProfitability,
		RanOutOfCash:                fs.Cash <= 0 && !fs.HasExited,
		FeaturesCompleted:           featuresCompleted,
		InnovationLeader:            calculateInnovationLeader(fs),
		EnterpriseFeatures:          enterpriseFeatures,
		CustomerLossDuringRoadmap:   fs.CustomersLostDuringRoadmap > 0,
		EnterpriseCustomers:         enterpriseCustomers,
		VerticalConcentration:       verticalConcentration,
		PricingExperimentsCompleted: pricingExperiments,
		PremiumPricingSuccess:       premiumPricing,
		LowTouchCustomers:           lowTouchCustomers,
		DealsClosedWon:              dealsClosedWon,
		HighProbabilityClose:        highProbClose,
		MaxPipelineSize:             maxPipelineSize,
		CustomerChurnRate:           fs.CustomerChurnRate,
		TotalGames:                  playerStats.TotalGames,
		TotalWins:                   int(playerStats.WinRate * float64(playerStats.TotalGames) / 100.0),
		WinStreak:                   winStreak,
		BestNetWorth:                playerStats.BestNetWorth,
		TotalExits:                  playerStats.TotalExits,
	}

	// Get previously unlocked achievements
	previouslyUnlocked, _ := database.GetPlayerAchievements(fs.FounderName)

	// Check for new achievements
	newUnlocks := achievements.CheckAchievements(gameStats, previouslyUnlocked)

	for _, ach := range newUnlocks {
		s.newAchievements = append(s.newAchievements, ach.Name)
		s.newAchievementObjs = append(s.newAchievementObjs, ach)
		database.UnlockAchievement(fs.FounderName, ach.ID)
	}
}

// calculateInnovationLeader checks if the founder was first to launch features
func calculateInnovationLeader(fs *founder.FounderState) bool {
	if fs.ProductRoadmap == nil || len(fs.ProductRoadmap.Features) == 0 {
		return false
	}

	featureMap := map[string][]string{
		"REST API":                  {"API", "REST API", "API Integration"},
		"Mobile App":               {"Mobile App", "Mobile", "iOS App", "Android App"},
		"Enterprise SSO":           {"SSO", "Single Sign-On", "Enterprise SSO"},
		"Advanced Analytics":       {"Analytics", "Advanced Analytics", "Reporting"},
		"AI/ML Capabilities":       {"AI", "ML", "Machine Learning", "Artificial Intelligence"},
		"Integrations Hub":         {"Integrations", "Integration Hub", "API Integration"},
		"Security Suite":           {"Security", "Security Suite", "Enterprise Security"},
		"Performance Optimization": {"Performance", "Optimization"},
		"White Label":              {"White Label", "White-Label"},
		"Workflow Automation":      {"Automation", "Workflow"},
	}

	for _, feature := range fs.ProductRoadmap.Features {
		if feature.Status != "completed" {
			continue
		}

		matchingNames, exists := featureMap[feature.Name]
		if !exists {
			matchingNames = []string{feature.Name}
		}

		wasFirst := true
		for _, launch := range fs.ProductRoadmap.CompetitorLaunches {
			for _, name := range matchingNames {
				if strings.EqualFold(launch.FeatureName, name) && launch.MonthLaunched <= feature.MonthCompleted {
					wasFirst = false
					break
				}
			}
			if !wasFirst {
				break
			}
		}

		if wasFirst {
			return true
		}
	}
	return false
}

// Update handles results screen input
func (s *FounderResultsScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Enter):
			switch s.phase {
			case PhaseResults:
				s.phase = PhaseXP
			case PhaseXP:
				if s.leveledUp {
					s.phase = PhaseLevelUp
				} else {
					s.phase = PhaseLeaderboard
				}
			case PhaseLevelUp:
				s.phase = PhaseLeaderboard
			case PhaseLeaderboard:
				if len(s.newAchievements) > 0 {
					s.phase = PhaseAchievements
				} else {
					s.phase = PhaseDone
					return s, SwitchTo(ScreenMainMenu)
				}
			case PhaseAchievements:
				s.phase = PhaseDone
				return s, SwitchTo(ScreenMainMenu)
			}

		case key.Matches(msg, keys.Global.Back):
			return s, SwitchTo(ScreenMainMenu)
		}
	}

	return s, nil
}

// View renders the results screen
func (s *FounderResultsScreen) View() string {
	switch s.phase {
	case PhaseXP:
		return s.renderXPBreakdown()
	case PhaseLevelUp:
		return s.renderLevelUp()
	case PhaseLeaderboard:
		return s.renderLeaderboard()
	case PhaseAchievements:
		return s.renderAchievements()
	default:
		return s.renderResults()
	}
}

func (s *FounderResultsScreen) renderResults() string {
	fs := s.gameData.FounderState
	var b strings.Builder

	// Header based on outcome
	var headerStyle lipgloss.Style
	var headerText string

	if fs.HasExited && fs.ExitType == "ipo" {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Green).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "🏛️  IPO SUCCESS! 🏛️"
	} else if fs.HasExited && fs.ExitType == "acquisition" {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Green).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "🤝 ACQUISITION COMPLETE! 🤝"
	} else if fs.HasExited && fs.ExitType == "secondary" {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Cyan).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "💼 SECONDARY SALE COMPLETE 💼"
	} else if fs.Cash <= 0 {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Red).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "💸 GAME OVER - RAN OUT OF CASH 💸"
	} else {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Yellow).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "⏰ TIME'S UP - GAME OVER ⏰"
	}

	b.WriteString(headerStyle.Render(headerText))
	b.WriteString("\n\n")

	// Results box
	resultsBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 3).
		Width(60)

	var results strings.Builder

	// Player/Company info
	playerStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	results.WriteString(playerStyle.Render(fmt.Sprintf("Founder: %s", fs.FounderName)))
	results.WriteString("\n")
	results.WriteString(playerStyle.Render(fmt.Sprintf("Company: %s (%s)", fs.CompanyName, fs.Category)))
	results.WriteString("\n")
	monthsPlayed := fs.Turn
	if fs.HasExited && fs.ExitMonth > 0 {
		monthsPlayed = fs.ExitMonth
	}
	results.WriteString(fmt.Sprintf("Months Played: %d", monthsPlayed))
	results.WriteString("\n\n")

	// Outcome
	resultsHeader := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	results.WriteString(resultsHeader.Render("═══ OUTCOME ═══"))
	results.WriteString("\n\n")

	outcomeStyle := lipgloss.NewStyle().Foreground(styles.White).Bold(true)
	results.WriteString(outcomeStyle.Render(s.outcome))
	results.WriteString("\n\n")

	// Financial results
	results.WriteString(resultsHeader.Render("═══ FINANCIAL RESULTS ═══"))
	results.WriteString("\n\n")

	if fs.HasExited {
		results.WriteString(fmt.Sprintf("Exit Valuation: $%s\n", formatCompactMoney(fs.ExitValuation)))
		results.WriteString(fmt.Sprintf("Your Equity: %.1f%%\n", s.founderEquity))

		payoutStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		switch fs.ExitType {
		case "ipo":
			immediateVal := int64(float64(fs.ExitValuation) * s.founderEquity * 0.20 / 100.0)
			remainingVal := int64(float64(fs.ExitValuation) * s.founderEquity * 0.80 / 100.0)
			results.WriteString(payoutStyle.Render(fmt.Sprintf("Immediate Liquidation (20%%): $%s", formatCompactMoney(immediateVal))))
			results.WriteString("\n")
			results.WriteString(fmt.Sprintf("Remaining Equity Value: $%s", formatCompactMoney(remainingVal)))
			results.WriteString("\n")
			results.WriteString(payoutStyle.Render(fmt.Sprintf("Total Net Worth: $%s", formatCompactMoney(immediateVal+remainingVal))))
		case "acquisition":
			results.WriteString(payoutStyle.Render(fmt.Sprintf("Your Payout: $%s", formatCompactMoney(s.founderPayout))))
		case "secondary":
			soldVal := int64(float64(fs.ExitValuation) * s.founderEquity * 0.5 / 100.0)
			remainVal := int64(float64(fs.ExitValuation) * s.founderEquity * 0.5 / 100.0)
			results.WriteString(fmt.Sprintf("Sold (50%% of stake): $%s", formatCompactMoney(soldVal)))
			results.WriteString("\n")
			results.WriteString(fmt.Sprintf("Remaining Equity: $%s", formatCompactMoney(remainVal)))
			results.WriteString("\n")
			results.WriteString(payoutStyle.Render(fmt.Sprintf("Total Net Worth: $%s", formatCompactMoney(soldVal+remainVal))))
		}
		results.WriteString("\n")
		results.WriteString(fmt.Sprintf("Final MRR: $%s\n", formatCompactMoney(fs.MRR)))
		results.WriteString(fmt.Sprintf("Customers: %d", fs.Customers))
	} else {
		results.WriteString(fmt.Sprintf("Final Cash: $%s\n", formatCompactMoney(fs.Cash)))
		results.WriteString(fmt.Sprintf("Company Valuation: $%s\n", formatCompactMoney(s.valuation)))
		results.WriteString(fmt.Sprintf("MRR: $%s\n", formatCompactMoney(fs.MRR)))
		results.WriteString(fmt.Sprintf("Customers: %d\n", fs.Customers))
		results.WriteString(fmt.Sprintf("Your Equity: %.1f%%", s.founderEquity))
		if s.founderEquity > 0 && s.valuation > 0 {
			equityVal := int64(float64(s.valuation) * s.founderEquity / 100.0)
			valStyle := lipgloss.NewStyle().Foreground(styles.Green)
			results.WriteString("\n")
			results.WriteString(valStyle.Render(fmt.Sprintf("Your Equity Value: $%s", formatCompactMoney(equityVal))))
		}
	}

	results.WriteString("\n\n")

	// Team
	results.WriteString(fmt.Sprintf("Final Team Size: %d\n", fs.Team.TotalEmployees))
	results.WriteString(fmt.Sprintf("  Eng: %d | Sales: %d | CS: %d | Mkt: %d | C-Suite: %d",
		len(fs.Team.Engineers), len(fs.Team.Sales),
		len(fs.Team.CustomerSuccess), len(fs.Team.Marketing), len(fs.Team.Executives)))

	// Funding rounds
	if len(fs.FundingRounds) > 0 {
		results.WriteString("\n\n")
		results.WriteString(resultsHeader.Render("═══ FUNDING ROUNDS ═══"))
		results.WriteString("\n\n")

		totalRaised := int64(0)
		for _, round := range fs.FundingRounds {
			results.WriteString(fmt.Sprintf("%s: $%s (%.1f%% equity)",
				round.RoundName,
				formatCompactMoney(round.Amount),
				round.EquityGiven))
			if len(round.Investors) > 0 {
				investorStyle := lipgloss.NewStyle().Foreground(styles.Cyan)
				results.WriteString("\n")
				results.WriteString(investorStyle.Render(fmt.Sprintf("  Investors: %s", strings.Join(round.Investors, ", "))))
			}
			results.WriteString("\n")
			totalRaised += round.Amount
		}
		totalStyle := lipgloss.NewStyle().Foreground(styles.Green)
		results.WriteString(totalStyle.Render(fmt.Sprintf("Total Raised: $%s", formatCompactMoney(totalRaised))))
	}

	// Cap table breakdown (show for all endings)
	if len(fs.CapTable) > 0 || len(fs.FundingRounds) > 0 {
		results.WriteString("\n\n")
		if fs.HasExited {
			results.WriteString(resultsHeader.Render("═══ PAYOUT BREAKDOWN ═══"))
		} else {
			results.WriteString(resultsHeader.Render("═══ CAP TABLE ═══"))
		}
		results.WriteString("\n\n")

		// Use exit valuation if exited, otherwise computed valuation
		valForPayout := s.valuation
		if fs.HasExited && fs.ExitValuation > 0 {
			valForPayout = fs.ExitValuation
		}

		founderPayStyle := lipgloss.NewStyle().Foreground(styles.Green)
		founderVal := int64(float64(valForPayout) * s.founderEquity / 100.0)
		results.WriteString(founderPayStyle.Render(fmt.Sprintf("%-28s %6.1f%%  $%s",
			"You (Founder)", s.founderEquity, formatCompactMoney(founderVal))))
		results.WriteString("\n")

		// Investors
		investorStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
		for _, round := range fs.FundingRounds {
			investorName := round.RoundName
			if len(round.Investors) > 0 {
				investorName = round.Investors[0]
			}
			investorVal := int64(float64(valForPayout) * round.EquityGiven / 100.0)
			results.WriteString(investorStyle.Render(fmt.Sprintf("%-28s %6.1f%%  $%s",
				truncate(investorName+" ("+round.RoundName+")", 28), round.EquityGiven, formatCompactMoney(investorVal))))
			results.WriteString("\n")
		}

		// Employees, execs, advisors from cap table
		for _, entry := range fs.CapTable {
			entryVal := int64(float64(valForPayout) * entry.Equity / 100.0)
			label := entry.Name
			switch entry.Type {
			case "executive":
				label += " (Exec)"
			case "employee":
				label += " (Emp)"
			case "advisor":
				label += " (Adv)"
			}
			results.WriteString(fmt.Sprintf("%-28s %6.2f%%  $%s\n",
				truncate(label, 28), entry.Equity, formatCompactMoney(entryVal)))
		}

		// Unallocated pool
		unallocatedPool := fs.EquityPool - fs.EquityAllocated
		if unallocatedPool < 0 {
			unallocatedPool = 0
		}
		if unallocatedPool > 0 {
			dimStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			results.WriteString(dimStyle.Render(fmt.Sprintf("%-28s %6.1f%%  (cancelled at exit)",
				"Unallocated Pool", unallocatedPool)))
			results.WriteString("\n")
		}
	}

	results.WriteString("\n\n")

	// Rating
	ratingStyle := lipgloss.NewStyle().Foreground(styles.Gold).Bold(true)
	results.WriteString("Rating: ")
	results.WriteString(ratingStyle.Render(fmt.Sprintf("%s %s", s.ratingIcon, s.rating)))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(resultsBox.Render(results.String())))
	b.WriteString("\n\n")

	// Score saved indicator
	if s.scoreSaved {
		savedStyle := lipgloss.NewStyle().Foreground(styles.Green).Width(s.width).Align(lipgloss.Center)
		b.WriteString(savedStyle.Render("✓ Score saved to leaderboard"))
		b.WriteString("\n\n")
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("Press Enter to see XP earned"))

	return b.String()
}

func (s *FounderResultsScreen) renderXPBreakdown() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("📊 EXPERIENCE EARNED 📊"))
	b.WriteString("\n\n")

	// XP breakdown box
	xpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(50)

	var xpContent strings.Builder

	sourceStyle := lipgloss.NewStyle().Foreground(styles.White)
	xpStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)

	for source, amount := range s.xpBreakdown {
		xpContent.WriteString(xpStyle.Render(fmt.Sprintf("+%d XP ", amount)))
		xpContent.WriteString(sourceStyle.Render(source))
		xpContent.WriteString("\n")
	}

	xpContent.WriteString("\n")
	xpContent.WriteString(strings.Repeat("─", 40))
	xpContent.WriteString("\n")

	totalStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	xpContent.WriteString(totalStyle.Render(fmt.Sprintf("Total XP Gained: +%d XP", s.totalXP)))

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(xpBox.Render(xpContent.String())))
	b.WriteString("\n\n")

	// Level progress
	if s.profileAfter != nil {
		levelBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Magenta).
			Padding(1, 2).
			Width(50)

		var levelContent strings.Builder
		levelStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
		levelInfo := progression.GetLevelInfo(s.profileAfter.Level)

		levelContent.WriteString(levelStyle.Render(fmt.Sprintf("Level %d - %s", s.profileAfter.Level, levelInfo.Title)))
		levelContent.WriteString("\n\n")

		progressBar := progression.FormatXPBar(s.profileAfter.ExperiencePoints, s.profileAfter.NextLevelXP, 30)
		progressStyle := lipgloss.NewStyle().Foreground(styles.Cyan)
		levelContent.WriteString(progressStyle.Render(progressBar))
		levelContent.WriteString("\n")
		levelContent.WriteString(fmt.Sprintf("%d / %d XP", s.profileAfter.ExperiencePoints, s.profileAfter.NextLevelXP))

		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(levelBox.Render(levelContent.String())))
		b.WriteString("\n\n")
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	if s.leveledUp {
		b.WriteString(helpStyle.Render("Press Enter to see level up rewards!"))
	} else {
		b.WriteString(helpStyle.Render("Press Enter to continue"))
	}

	return b.String()
}

func (s *FounderResultsScreen) renderLevelUp() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Gold).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("🎉 LEVEL UP! 🎉"))
	b.WriteString("\n\n")

	levelBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Gold).
		Padding(1, 3).
		Width(55)

	var levelContent strings.Builder

	oldStyle := lipgloss.NewStyle().Foreground(styles.Gray)
	arrowStyle := lipgloss.NewStyle().Foreground(styles.White).Bold(true)
	newStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)

	oldLevelInfo := progression.GetLevelInfo(s.oldLevel)
	newLevelInfo := progression.GetLevelInfo(s.newLevel)

	levelContent.WriteString(oldStyle.Render(fmt.Sprintf("Level %d - %s", s.oldLevel, oldLevelInfo.Title)))
	levelContent.WriteString("\n")
	levelContent.WriteString(arrowStyle.Render("           ↓"))
	levelContent.WriteString("\n")
	levelContent.WriteString(newStyle.Render(fmt.Sprintf("Level %d - %s", s.newLevel, newLevelInfo.Title)))
	levelContent.WriteString("\n\n")

	if s.levelUpPoints > 0 {
		pointsStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		levelContent.WriteString(pointsStyle.Render(fmt.Sprintf("💰 Bonus Points: +%d", s.levelUpPoints)))
		levelContent.WriteString("\n\n")
	}

	unlocks := newLevelInfo.Unlocks
	if len(unlocks) > 0 {
		unlockHeader := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
		levelContent.WriteString(unlockHeader.Render("🔓 NEW UNLOCKS:"))
		levelContent.WriteString("\n")

		unlockStyle := lipgloss.NewStyle().Foreground(styles.White)
		for _, unlock := range unlocks {
			levelContent.WriteString(unlockStyle.Render("   • " + unlock))
			levelContent.WriteString("\n")
		}
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(levelBox.Render(levelContent.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *FounderResultsScreen) renderLeaderboard() string {
	fs := s.gameData.FounderState
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Magenta).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("🏆 GAME SUMMARY 🏆"))
	b.WriteString("\n\n")

	// Summary box
	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(60)

	var content strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	content.WriteString(titleStyle.Render("Final Statistics"))
	content.WriteString("\n\n")

	content.WriteString(fmt.Sprintf("Company:     %s\n", fs.CompanyName))
	content.WriteString(fmt.Sprintf("Category:    %s\n", fs.Category))
	durationMonths := fs.Turn
	if fs.HasExited && fs.ExitMonth > 0 {
		durationMonths = fs.ExitMonth
	}
	content.WriteString(fmt.Sprintf("Duration:    %d months\n", durationMonths))
	content.WriteString(fmt.Sprintf("Peak MRR:    $%s\n", formatCompactMoney(fs.MRR)))
	content.WriteString(fmt.Sprintf("Total Customers: %d (ever: %d)\n", fs.Customers, fs.TotalCustomersEver))

	totalRaised := int64(0)
	for _, round := range fs.FundingRounds {
		totalRaised += round.Amount
	}
	content.WriteString(fmt.Sprintf("Total Raised: $%s\n", formatCompactMoney(totalRaised)))
	content.WriteString(fmt.Sprintf("Final Team:   %d employees\n", fs.Team.TotalEmployees))

	if fs.HasExited {
		exitStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		content.WriteString(exitStyle.Render(fmt.Sprintf("Exit:        %s ($%s)", fs.ExitType, formatCompactMoney(fs.ExitValuation))))
	} else if fs.Cash <= 0 {
		failStyle := lipgloss.NewStyle().Foreground(styles.Red)
		content.WriteString(failStyle.Render("Outcome:     Ran out of cash"))
	} else {
		content.WriteString(fmt.Sprintf("Valuation:   $%s", formatCompactMoney(s.valuation)))
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(summaryBox.Render(content.String())))
	b.WriteString("\n\n")

	// Global leaderboard prompt
	if !s.leaderboardSubmitted {
		promptStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Width(s.width).Align(lipgloss.Center)
		b.WriteString(promptStyle.Render("Score saved locally. Global leaderboard submission available from main menu."))
		b.WriteString("\n\n")
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	if len(s.newAchievements) > 0 {
		b.WriteString(helpStyle.Render("Press Enter to see achievements"))
	} else {
		b.WriteString(helpStyle.Render("Press Enter to return to main menu"))
	}

	return b.String()
}

func (s *FounderResultsScreen) renderAchievements() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Gold).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("🏆 ACHIEVEMENTS UNLOCKED! 🏆"))
	b.WriteString("\n\n")

	achBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Gold).
		Padding(1, 2).
		Width(55)

	var achList strings.Builder
	for i, ach := range s.newAchievementObjs {
		achStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
		descStyle := lipgloss.NewStyle().Foreground(styles.Gray)
		achList.WriteString(achStyle.Render(fmt.Sprintf("%s %s [%s]", ach.Icon, ach.Name, ach.Rarity)))
		achList.WriteString("\n")
		achList.WriteString(descStyle.Render("   "+ach.Description))
		if i < len(s.newAchievementObjs)-1 {
			achList.WriteString("\n")
		}
		achList.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(achBox.Render(achList.String())))
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("Press Enter to return to main menu"))

	return b.String()
}

// submitToGlobalLeaderboard submits the score to the global leaderboard
func (s *FounderResultsScreen) submitToGlobalLeaderboard() {
	fs := s.gameData.FounderState

	_, valuation, founderEquity := fs.GetFinalScore()

	var founderPayout int64
	if fs.HasExited {
		switch fs.ExitType {
		case "ipo":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.20 / 100.0)
		case "acquisition":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity / 100.0)
		case "secondary":
			founderPayout = int64(float64(fs.ExitValuation) * founderEquity * 0.5 / 100.0)
		}
	} else {
		founderPayout = int64(float64(valuation) * founderEquity / 100.0)
	}

	maxARR := fs.MRR * 12

	totalFundingRaised := int64(0)
	for _, round := range fs.FundingRounds {
		totalFundingRaised += round.Amount
	}

	submission := leaderboard.FounderScoreSubmission{
		PlayerName:        fs.FounderName,
		FinalValuation:    valuation,
		FounderEquity:     founderEquity,
		FounderPayout:     founderPayout,
		ExitType:          fs.ExitType,
		ExitMonth:         fs.ExitMonth,
		MaxARR:            maxARR,
		StartupTemplate:   fs.CompanyName,
		FundingRaised:     totalFundingRaised,
		CustomersAcquired: fs.TotalCustomersEver,
	}

	err := leaderboard.SubmitFounderScore(submission, "")
	if err != nil {
		s.leaderboardError = err.Error()
	} else {
		s.leaderboardSubmitted = true
	}
}
