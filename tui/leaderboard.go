package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// LeaderboardScreen shows the leaderboards
type LeaderboardScreen struct {
	width         int
	height        int
	table         *components.GameTable
	menu          *components.Menu
	currentView   string
	currentFilter string
	currentMode   string // "all", "vc", "founder"
}

// NewLeaderboardScreen creates a new leaderboard screen
func NewLeaderboardScreen(width, height int) *LeaderboardScreen {
	// Filter menu
	menuItems := []components.MenuItem{
		{ID: "header_mode", Title: "── MODE ──", Disabled: true, Icon: ""},
		{ID: "mode_all", Title: "All Modes", Icon: "🌟"},
		{ID: "mode_vc", Title: "VC Mode", Icon: "🎩"},
		{ID: "mode_founder", Title: "Founder Mode", Icon: "🚀"},
		{ID: "header_sort", Title: "── SORT ──", Disabled: true, Icon: ""},
		{ID: "net_worth", Title: "By Net Worth", Icon: "💰"},
		{ID: "roi", Title: "By ROI", Icon: "📈"},
		{ID: "recent", Title: "Recent Games", Icon: "🕐"},
		{ID: "header_diff", Title: "── DIFFICULTY ──", Disabled: true, Icon: ""},
		{ID: "easy", Title: "Easy Difficulty", Icon: "🟢"},
		{ID: "medium", Title: "Medium Difficulty", Icon: "🟡"},
		{ID: "hard", Title: "Hard Difficulty", Icon: "🔴"},
		{ID: "expert", Title: "Expert Difficulty", Icon: "💀"},
	}
	menu := components.NewMenu("LEADERBOARD FILTERS", menuItems)
	menu.SetSize(35, 18)
	menu.SetHideHelp(true)

	s := &LeaderboardScreen{
		width:       width,
		height:      height,
		menu:        menu,
		currentView: "net_worth",
		currentMode: "all",
	}

	s.loadLeaderboard("net_worth", "all", "all")
	return s
}

func (s *LeaderboardScreen) loadLeaderboard(sortBy, difficulty, mode string) {
	var scores []database.GameScore
	var err error

	switch sortBy {
	case "recent":
		scores, err = database.GetRecentGamesAndMode(20, mode)
	case "roi":
		scores, err = database.GetTopScoresByROIAndMode(20, difficulty, mode)
	default:
		scores, err = database.GetTopScoresByNetWorthAndMode(20, difficulty, mode)
	}

	if err != nil || len(scores) == 0 {
		columns := []table.Column{
			{Title: "#", Width: 4},
			{Title: "Player", Width: 15},
			{Title: "Net Worth", Width: 12},
			{Title: "ROI", Width: 8},
			{Title: "Mode", Width: 9},
			{Title: "Difficulty", Width: 10},
		}
		s.table = components.NewGameTable("", columns, []table.Row{})
		s.table.SetSize(62, 12)
		return
	}

	rows := make([]table.Row, len(scores))
	for i, score := range scores {
		modeStr := score.Mode
		if modeStr == "" {
			modeStr = "vc"
		}
		rows[i] = table.Row{
			fmt.Sprintf("%d", i+1),
			truncate(score.PlayerName, 15),
			fmt.Sprintf("$%s", formatCompactMoney(score.FinalNetWorth)),
			fmt.Sprintf("%.1f%%", score.ROI),
			modeStr,
			score.Difficulty,
		}
	}

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Player", Width: 15},
		{Title: "Net Worth", Width: 12},
		{Title: "ROI", Width: 8},
		{Title: "Mode", Width: 9},
		{Title: "Difficulty", Width: 10},
	}

	s.table = components.NewGameTable("", columns, rows)
	s.table.SetSize(62, 12)
}

// Init initializes the leaderboard screen
func (s *LeaderboardScreen) Init() tea.Cmd {
	return nil
}

// Update handles leaderboard input
func (s *LeaderboardScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Global.Back) {
			return s, PopScreen()
		}

	case components.MenuSelectedMsg:
		switch msg.ID {
		// Mode filters
		case "mode_all":
			s.currentMode = "all"
		case "mode_vc":
			s.currentMode = "vc"
		case "mode_founder":
			s.currentMode = "founder"
		// Sort
		case "net_worth":
			s.currentView = "net_worth"
		case "roi":
			s.currentView = "roi"
		case "recent":
			s.currentView = "recent"
		// Difficulty (only applies to VC mode; founder uses "Founder" tag)
		case "easy":
			s.currentView = "net_worth"
			s.loadLeaderboard("net_worth", "Easy", s.currentMode)
			s.currentView = "easy"
			var cmd tea.Cmd
			s.menu, cmd = s.menu.Update(msg)
			return s, cmd
		case "medium":
			s.currentView = "net_worth"
			s.loadLeaderboard("net_worth", "Medium", s.currentMode)
			s.currentView = "medium"
			var cmd tea.Cmd
			s.menu, cmd = s.menu.Update(msg)
			return s, cmd
		case "hard":
			s.currentView = "net_worth"
			s.loadLeaderboard("net_worth", "Hard", s.currentMode)
			s.currentView = "hard"
			var cmd tea.Cmd
			s.menu, cmd = s.menu.Update(msg)
			return s, cmd
		case "expert":
			s.currentView = "net_worth"
			s.loadLeaderboard("net_worth", "Expert", s.currentMode)
			s.currentView = "expert"
			var cmd tea.Cmd
			s.menu, cmd = s.menu.Update(msg)
			return s, cmd
		}

		// Reload with current sort + mode + difficulty
		difficulty := "all"
		switch s.currentView {
		case "easy":
			difficulty = "Easy"
		case "medium":
			difficulty = "Medium"
		case "hard":
			difficulty = "Hard"
		case "expert":
			difficulty = "Expert"
		}
		s.loadLeaderboard(s.currentView, difficulty, s.currentMode)
	}

	var cmd tea.Cmd
	s.menu, cmd = s.menu.Update(msg)
	return s, cmd
}

// View renders the leaderboard
func (s *LeaderboardScreen) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Gold).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("🏆 LEADERBOARDS 🏆")))
	b.WriteString("\n\n")

	// Layout: menu on left, table on right
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Gold).
		Padding(0, 1)

	tableBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(0, 1)

	leftPanel := menuBox.Render(s.menu.View())
	rightPanel := tableBox.Render(s.table.View())

	// Center with margin only so bordered boxes never reflow (avoids disjointed bottom border).
	panelRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, "  ", rightPanel)
	rowWidth := lipgloss.Width(panelRow)
	margin := (s.width - rowWidth) / 2
	if margin < 0 {
		margin = 0
	}
	b.WriteString(lipgloss.NewStyle().MarginLeft(margin).Render(panelRow))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter select • esc back"))

	return b.String()
}