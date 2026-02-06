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
}

// NewLeaderboardScreen creates a new leaderboard screen
func NewLeaderboardScreen(width, height int) *LeaderboardScreen {
	// Filter menu
	menuItems := []components.MenuItem{
		{ID: "net_worth_all", Title: "By Net Worth (All)", Icon: "💰"},
		{ID: "roi_all", Title: "By ROI (All)", Icon: "📈"},
		{ID: "easy", Title: "Easy Difficulty", Icon: "🟢"},
		{ID: "medium", Title: "Medium Difficulty", Icon: "🟡"},
		{ID: "hard", Title: "Hard Difficulty", Icon: "🔴"},
		{ID: "expert", Title: "Expert Difficulty", Icon: "💀"},
		{ID: "recent", Title: "Recent Games", Icon: "🕐"},
	}
	menu := components.NewMenu("LEADERBOARD FILTERS", menuItems)
	menu.SetSize(35, 15)
	menu.SetHideHelp(true)

	s := &LeaderboardScreen{
		width:       width,
		height:      height,
		menu:        menu,
		currentView: "net_worth_all",
	}

	s.loadLeaderboard("net_worth", "all")
	return s
}

func (s *LeaderboardScreen) loadLeaderboard(sortBy, difficulty string) {
	var scores []database.GameScore
	var err error

	if sortBy == "recent" {
		scores, err = database.GetRecentGames(20)
	} else if sortBy == "roi" {
		scores, err = database.GetTopScoresByROI(20, difficulty)
	} else {
		scores, err = database.GetTopScoresByNetWorth(20, difficulty)
	}

	if err != nil || len(scores) == 0 {
		// Empty table
		columns := []table.Column{
			{Title: "#", Width: 4},
			{Title: "Player", Width: 15},
			{Title: "Net Worth", Width: 12},
			{Title: "ROI", Width: 8},
			{Title: "Difficulty", Width: 10},
		}
		s.table = components.NewGameTable("", columns, []table.Row{})
		s.table.SetSize(55, 12)
		return
	}

	rows := make([]table.Row, len(scores))
	for i, score := range scores {
		rows[i] = table.Row{
			fmt.Sprintf("%d", i+1),
			truncate(score.PlayerName, 15),
			fmt.Sprintf("$%s", formatCompactMoney(score.FinalNetWorth)),
			fmt.Sprintf("%.1f%%", score.ROI),
			score.Difficulty,
		}
	}

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Player", Width: 15},
		{Title: "Net Worth", Width: 12},
		{Title: "ROI", Width: 8},
		{Title: "Difficulty", Width: 10},
	}

	s.table = components.NewGameTable("", columns, rows)
	s.table.SetSize(55, 12)
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
		s.currentView = msg.ID
		switch msg.ID {
		case "net_worth_all":
			s.loadLeaderboard("net_worth", "all")
		case "roi_all":
			s.loadLeaderboard("roi", "all")
		case "easy":
			s.loadLeaderboard("net_worth", "Easy")
		case "medium":
			s.loadLeaderboard("net_worth", "Medium")
		case "hard":
			s.loadLeaderboard("net_worth", "Hard")
		case "expert":
			s.loadLeaderboard("net_worth", "Expert")
		case "recent":
			s.loadLeaderboard("recent", "all")
		}
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
