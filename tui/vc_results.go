package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/achievements"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/progression"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// ResultsPhase tracks what part of results we're showing
type ResultsPhase int

const (
	PhaseResults ResultsPhase = iota
	PhaseXP
	PhaseLevelUp
	PhaseLeaderboard
	PhaseAchievements
	PhaseDone
)

// VCResultsScreen shows the end-game results
type VCResultsScreen struct {
	width    int
	height   int
	gameData *GameData
	phase    ResultsPhase

	// Computed results
	netWorth        int64
	roi             float64
	successfulExits int
	playerRank      int
	won             bool
	rating          string
	ratingIcon      string

	// New achievements unlocked
	newAchievements    []string
	newAchievementObjs []achievements.Achievement

	// XP and Level tracking
	xpBreakdown    map[string]int
	totalXP        int
	leveledUp      bool
	oldLevel       int
	newLevel       int
	levelUpPoints  int
	profileBefore  *database.PlayerProfile
	profileAfter   *database.PlayerProfile

	// Leaderboard
	leaderboardTable *components.GameTable

	// Score saved
	scoreSaved bool
}

// NewVCResultsScreen creates a new results screen
func NewVCResultsScreen(width, height int, gameData *GameData) *VCResultsScreen {
	gs := gameData.GameState
	netWorth, roi, successfulExits := gs.GetFinalScore()

	// Determine if player won (doubled starting cash)
	won := netWorth >= gs.Difficulty.StartingCash*2

	// Calculate rating
	rating, ratingIcon := calculateRating(roi)

	// Find player rank
	leaderboard := gs.GetLeaderboard()
	playerRank := 1
	for i, entry := range leaderboard {
		if entry.IsPlayer {
			playerRank = i + 1
			break
		}
	}

	// Build leaderboard table
	rows := make([]table.Row, len(leaderboard))
	for i, entry := range leaderboard {
		marker := ""
		if entry.IsPlayer {
			marker = "→"
		}

		// Format ROI - use compact notation for very large numbers
		roiStr := fmt.Sprintf("%.0f%%", entry.ROI)
		if entry.ROI >= 10000 {
			roiStr = fmt.Sprintf("%.0fK%%", entry.ROI/1000)
		} else if entry.ROI >= 1000 {
			roiStr = fmt.Sprintf("%.1fK%%", entry.ROI/1000)
		}

		rows[i] = table.Row{
			fmt.Sprintf("%s%d", marker, i+1),
			truncate(entry.Name, 18),
			truncate(entry.Firm, 22),
			fmt.Sprintf("$%s", formatCompactMoney(entry.NetWorth)),
			roiStr,
		}
	}

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Investor", Width: 18},
		{Title: "Firm", Width: 22},
		{Title: "Net Worth", Width: 10},
		{Title: "ROI", Width: 10},
	}

	leaderboardTable := components.NewGameTable("", columns, rows)
	// Set height to accommodate all players (max 5 AI + 1 player = 6, plus header)
	leaderboardTable.SetSize(68, len(leaderboard)+2)

	s := &VCResultsScreen{
		width:            width,
		height:           height,
		gameData:         gameData,
		phase:            PhaseResults,
		netWorth:         netWorth,
		roi:              roi,
		successfulExits:  successfulExits,
		playerRank:       playerRank,
		won:              won,
		rating:           rating,
		ratingIcon:       ratingIcon,
		leaderboardTable: leaderboardTable,
	}

	return s
}

func calculateRating(roi float64) (string, string) {
	if roi >= 1000 {
		return "UNICORN HUNTER - Legendary!", "👑"
	} else if roi >= 500 {
		return "Elite VC - Outstanding!", "🏆"
	} else if roi >= 200 {
		return "Great Investor - Excellent!", "⭐"
	} else if roi >= 50 {
		return "Solid Performance - Good!", "✓"
	} else if roi >= 10 {
		return "Profitable - Not Bad", "✓"
	} else if roi >= -10 {
		return "Break Even - Survived", "="
	}
	return "Lost Money - Better Luck Next Time", "⚠"
}

// Init initializes the results screen
func (s *VCResultsScreen) Init() tea.Cmd {
	// Save score to database
	gs := s.gameData.GameState
	score := database.GameScore{
		PlayerName:      gs.PlayerName,
		FinalNetWorth:   s.netWorth,
		ROI:             s.roi,
		SuccessfulExits: s.successfulExits,
		TurnsPlayed:     gs.Portfolio.Turn - 1,
		Difficulty:      gs.Difficulty.Name,
		PlayedAt:        time.Now(),
	}

	err := database.SaveGameScore(score)
	if err == nil {
		s.scoreSaved = true
	}

	// Get profile before XP is added
	s.profileBefore, _ = database.GetPlayerProfile(gs.PlayerName)
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
	s.leveledUp, s.newLevel, s.levelUpPoints, levelUpErr = database.AddExperience(gs.PlayerName, s.totalXP)
	if levelUpErr == nil {
		s.profileAfter, _ = database.GetPlayerProfile(gs.PlayerName)
	}

	return nil
}

func (s *VCResultsScreen) calculateXPBreakdown() {
	gs := s.gameData.GameState
	s.xpBreakdown = make(map[string]int)

	// Base XP for completing a game
	s.xpBreakdown["Game Completion"] = progression.XPGameComplete
	s.totalXP = progression.XPGameComplete

	// Positive ROI bonus
	if s.roi > 0 {
		s.xpBreakdown["Positive ROI"] = progression.XPPositiveROI
		s.totalXP += progression.XPPositiveROI
	}

	// Successful exits bonus
	if s.successfulExits > 0 {
		exitXP := progression.XPSuccessfulExit * s.successfulExits
		s.xpBreakdown[fmt.Sprintf("Successful Exits (%d)", s.successfulExits)] = exitXP
		s.totalXP += exitXP
	}

	// Difficulty bonus
	switch strings.ToLower(gs.Difficulty.Name) {
	case "medium":
		s.xpBreakdown["Medium Difficulty"] = progression.XPDifficultyMedium
		s.totalXP += progression.XPDifficultyMedium
	case "hard":
		s.xpBreakdown["Hard Difficulty"] = progression.XPDifficultyHard
		s.totalXP += progression.XPDifficultyHard
	case "expert":
		s.xpBreakdown["Expert Difficulty"] = progression.XPDifficultyExpert
		s.totalXP += progression.XPDifficultyExpert
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

func (s *VCResultsScreen) checkAchievements() {
	gs := s.gameData.GameState

	// Build game stats for achievement checking
	stats := achievements.GameStats{
		GameMode:        "vc",
		FinalNetWorth:   s.netWorth,
		ROI:             s.roi,
		TotalInvested:   gs.GetTotalInvested(),
		TurnsPlayed:     gs.Portfolio.Turn,
		InvestmentCount: len(gs.Portfolio.Investments),
		SuccessfulExits: s.successfulExits,
		Difficulty:      gs.Difficulty.Name,
	}

	// Get previously unlocked achievements
	previouslyUnlocked, _ := database.GetPlayerAchievements(gs.PlayerName)

	// Check for new achievements
	newUnlocks := achievements.CheckAchievements(stats, previouslyUnlocked)

	for _, ach := range newUnlocks {
		s.newAchievements = append(s.newAchievements, ach.Name)
		s.newAchievementObjs = append(s.newAchievementObjs, ach)
		database.UnlockAchievement(gs.PlayerName, ach.ID)
	}
}

// Update handles results screen input
func (s *VCResultsScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
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
func (s *VCResultsScreen) View() string {
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

func (s *VCResultsScreen) renderResults() string {
	gs := s.gameData.GameState
	var b strings.Builder

	// Big header based on outcome
	var headerStyle lipgloss.Style
	var headerText string

	if s.won {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Green).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "🎉 VICTORY! 🎉"
	} else {
		headerStyle = lipgloss.NewStyle().
			Foreground(styles.Red).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		headerText = "GAME OVER"
	}

	b.WriteString(headerStyle.Render(headerText))
	b.WriteString("\n\n")

	// Results box
	resultsBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 3).
		Width(50)

	var results strings.Builder

	// Player info
	playerStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
	results.WriteString(playerStyle.Render(fmt.Sprintf("Player: %s", gs.PlayerName)))
	results.WriteString("\n")
	results.WriteString(playerStyle.Render(fmt.Sprintf("Firm: %s", s.gameData.FirmName)))
	results.WriteString("\n")
	results.WriteString(fmt.Sprintf("Difficulty: %s", gs.Difficulty.Name))
	results.WriteString("\n")
	results.WriteString(fmt.Sprintf("Turns Played: %d", gs.Portfolio.Turn-1))
	results.WriteString("\n\n")

	// Financial results
	resultsHeader := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	results.WriteString(resultsHeader.Render("═══ FINANCIAL RESULTS ═══"))
	results.WriteString("\n\n")

	// Net worth
	netWorthStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	results.WriteString("Final Net Worth: ")
	results.WriteString(netWorthStyle.Render(fmt.Sprintf("$%s", formatCompactMoney(s.netWorth))))
	results.WriteString("\n")

	// ROI
	roiStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if s.roi < 0 {
		roiStyle = lipgloss.NewStyle().Foreground(styles.Red)
	}
	results.WriteString("ROI: ")
	results.WriteString(roiStyle.Render(fmt.Sprintf("%.1f%%", s.roi)))
	results.WriteString("\n")

	// Exits
	results.WriteString(fmt.Sprintf("Successful Exits (5x+): %d", s.successfulExits))
	results.WriteString("\n")

	// Management fees
	results.WriteString(fmt.Sprintf("Management Fees Paid: $%s", formatCompactMoney(gs.Portfolio.ManagementFeesCharged)))
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
	b.WriteString(helpStyle.Render("Press Enter to see final standings"))

	return b.String()
}

func (s *VCResultsScreen) renderLeaderboard() string {
	gs := s.gameData.GameState
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Magenta).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("🏆 FINAL LEADERBOARD 🏆"))
	b.WriteString("\n\n")

	// Build a simple table without the component (to avoid scrolling issues)
	leaderboard := gs.GetLeaderboard()

	tableBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(1, 2).
		Width(78)

	var tableContent strings.Builder

	// Header row
	headerRow := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	tableContent.WriteString(headerRow.Render(fmt.Sprintf("  %-4s %-18s %-24s %-12s %s", "#", "Investor", "Firm", "Net Worth", "ROI")))
	tableContent.WriteString("\n")
	tableContent.WriteString(strings.Repeat("─", 74))
	tableContent.WriteString("\n")

	// Data rows
	for i, entry := range leaderboard {
		marker := "  "
		rowStyle := lipgloss.NewStyle()
		if entry.IsPlayer {
			marker = "→ "
			rowStyle = rowStyle.Foreground(styles.Cyan).Bold(true)
		}

		// Format ROI compactly
		roiStr := fmt.Sprintf("%.0f%%", entry.ROI)
		if entry.ROI >= 10000 {
			roiStr = fmt.Sprintf("%.0fK%%", entry.ROI/1000)
		} else if entry.ROI >= 1000 {
			roiStr = fmt.Sprintf("%.1fK%%", entry.ROI/1000)
		}

		row := fmt.Sprintf("%s%-4d %-18s %-24s %-12s %s",
			marker,
			i+1,
			truncate(entry.Name, 18),
			truncate(entry.Firm, 24),
			"$"+formatCompactMoney(entry.NetWorth),
			roiStr,
		)
		tableContent.WriteString(rowStyle.Render(row))
		tableContent.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(tableBox.Render(tableContent.String())))
	b.WriteString("\n\n")

	// Player position
	posStyle := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	if s.playerRank == 1 {
		posStyle = posStyle.Foreground(styles.Gold).Bold(true)
		b.WriteString(posStyle.Render("🎉 CONGRATULATIONS! You beat all the AI investors!"))
	} else {
		posStyle = posStyle.Foreground(styles.Yellow)
		b.WriteString(posStyle.Render(fmt.Sprintf("You finished in position #%d", s.playerRank)))
	}
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	if len(s.newAchievements) > 0 {
		b.WriteString(helpStyle.Render("Press Enter to see achievements"))
	} else {
		b.WriteString(helpStyle.Render("Press Enter to return to main menu"))
	}

	return b.String()
}

func (s *VCResultsScreen) renderAchievements() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Gold).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("🏆 ACHIEVEMENTS UNLOCKED! 🏆"))
	b.WriteString("\n\n")

	// Achievement list
	achBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Gold).
		Padding(1, 2).
		Width(50)

	var achList strings.Builder
	for _, achName := range s.newAchievements {
		achStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
		achList.WriteString(achStyle.Render("⭐ " + achName))
		achList.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(achBox.Render(achList.String())))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("Press Enter to return to main menu"))

	return b.String()
}

func (s *VCResultsScreen) renderXPBreakdown() string {
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

	// Individual XP sources
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

	// Total XP
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

		// Progress bar
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
		b.WriteString(helpStyle.Render("Press Enter to see final standings"))
	}

	return b.String()
}

func (s *VCResultsScreen) renderLevelUp() string {
	var b strings.Builder

	// Big celebration header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Gold).
		Bold(true).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(headerStyle.Render("🎉 LEVEL UP! 🎉"))
	b.WriteString("\n\n")

	// Level up box
	levelBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Gold).
		Padding(1, 3).
		Width(55)

	var levelContent strings.Builder

	// Level change
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

	// Points earned
	if s.levelUpPoints > 0 {
		pointsStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		levelContent.WriteString(pointsStyle.Render(fmt.Sprintf("💰 Bonus Points: +%d", s.levelUpPoints)))
		levelContent.WriteString("\n\n")
	}

	// Unlocks
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

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("Press Enter to see final standings"))

	return b.String()
}
