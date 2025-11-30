package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// FounderView represents what we're showing in founder mode
type FounderView int

const (
	FounderViewMain FounderView = iota
	FounderViewActions
	FounderViewHiring
	FounderViewFunding
	FounderViewMetrics
)

// FounderGameScreen handles the founder game
type FounderGameScreen struct {
	width    int
	height   int
	gameData *GameData
	view     FounderView

	// Turn state
	turnMessages []string

	// Menus
	actionsMenu *components.Menu
}

// NewFounderGameScreen creates a new founder game screen
func NewFounderGameScreen(width, height int, gameData *GameData) *FounderGameScreen {
	// Actions menu
	actionItems := []components.MenuItem{
		{ID: "continue", Title: "Continue to next month", Description: "Process this month's events", Icon: "⏭️"},
		{ID: "hiring", Title: "View Team", Description: "See your team status", Icon: "👥"},
		{ID: "exit", Title: "View Exit Options", Description: "See available exit strategies", Icon: "🚪"},
	}
	actionsMenu := components.NewMenu("ACTIONS", actionItems)
	actionsMenu.SetSize(50, 15)
	actionsMenu.SetHideHelp(true)

	return &FounderGameScreen{
		width:       width,
		height:      height,
		gameData:    gameData,
		view:        FounderViewMain,
		actionsMenu: actionsMenu,
	}
}

// Init initializes the founder game screen
func (s *FounderGameScreen) Init() tea.Cmd {
	return nil
}

// Update handles founder game input
func (s *FounderGameScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.view {
		case FounderViewMain:
			switch {
			case key.Matches(msg, keys.Global.Enter):
				s.view = FounderViewActions
				return s, nil
			case key.Matches(msg, keys.Global.Back), msg.String() == "q":
				return s, SwitchTo(ScreenMainMenu)
			}

		case FounderViewActions:
			if key.Matches(msg, keys.Global.Back) {
				s.view = FounderViewMain
				return s, nil
			}
		}

	case components.MenuSelectedMsg:
		if s.view == FounderViewActions {
			return s.handleAction(msg.ID)
		}
	}

	// Check for game over
	if fg != nil && fg.IsGameOver() {
		// Show final results message and return to menu
		return s, SwitchTo(ScreenMainMenu)
	}

	// Update current component
	var cmd tea.Cmd
	switch s.view {
	case FounderViewActions:
		s.actionsMenu, cmd = s.actionsMenu.Update(msg)
	}

	return s, cmd
}

func (s *FounderGameScreen) handleAction(id string) (ScreenModel, tea.Cmd) {
	fg := s.gameData.FounderState

	switch id {
	case "continue":
		// Process month using the actual FounderState method
		msgs := fg.ProcessMonth()
		s.turnMessages = msgs
		s.view = FounderViewMain

	case "hiring":
		// Show team info in messages
		teamInfo := fmt.Sprintf("Team: %d engineers, %d sales, %d CS. Monthly cost: $%s",
			len(fg.Team.Engineers), len(fg.Team.Sales), len(fg.Team.CustomerSuccess),
			formatCompactMoney(fg.MonthlyTeamCost))
		s.turnMessages = []string{teamInfo}
		s.view = FounderViewMain

	case "exit":
		// Show exit options
		exits := fg.GetAvailableExits()
		var msgs []string
		for _, exit := range exits {
			status := "✓"
			if !exit.CanExit {
				status = "🔒"
			}
			msgs = append(msgs, fmt.Sprintf("%s %s: $%s", status, exit.Type, formatCompactMoney(exit.Valuation)))
		}
		s.turnMessages = msgs
		s.view = FounderViewMain
	}

	return s, nil
}

// View renders the founder game screen
func (s *FounderGameScreen) View() string {
	switch s.view {
	case FounderViewActions:
		return s.renderActions()
	default:
		return s.renderMain()
	}
}

func (s *FounderGameScreen) renderMain() string {
	fg := s.gameData.FounderState
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(70).
		Align(lipgloss.Center)

	header := fmt.Sprintf("🚀 %s - MONTH %d", fg.CompanyName, fg.Turn)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render(header)))
	b.WriteString("\n\n")

	// Main layout
	leftPanel := s.renderCompanyPanel()
	rightPanel := s.renderMetricsPanel()

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, "  ", rightPanel)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(content))
	b.WriteString("\n")

	// News
	if len(s.turnMessages) > 0 {
		b.WriteString(s.renderNews())
		b.WriteString("\n")
	}

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(styles.White).
		Background(styles.DarkGray).
		Width(70).
		Padding(0, 2)

	// Calculate valuation (simplified ARR * multiple)
	arr := fg.MRR * 12
	multiple := 10.0
	if fg.MonthlyGrowthRate > 0.05 {
		multiple = 15.0
	}
	valuation := int64(float64(arr) * multiple)

	status := fmt.Sprintf("💰 Cash: $%s  |  📊 Valuation: $%s  |  ⏳ Runway: %d mo",
		formatCompactMoney(fg.Cash),
		formatCompactMoney(valuation),
		fg.CashRunwayMonths)
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(statusStyle.Render(status)))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("enter actions • q quit"))

	return b.String()
}

func (s *FounderGameScreen) renderCompanyPanel() string {
	fg := s.gameData.FounderState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Magenta).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
	b.WriteString(titleStyle.Render("🏢 COMPANY"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	b.WriteString(labelStyle.Render("Founder: "))
	b.WriteString(fg.FounderName)
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Category: "))
	b.WriteString(fg.Category)
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Type: "))
	b.WriteString(fg.StartupType)
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Team Size: "))
	b.WriteString(fmt.Sprintf("%d", fg.Team.TotalEmployees))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("MRR: "))
	b.WriteString(fmt.Sprintf("$%s", formatCompactMoney(fg.MRR)))
	b.WriteString("\n")

	return panelStyle.Render(b.String())
}

func (s *FounderGameScreen) renderMetricsPanel() string {
	fg := s.gameData.FounderState

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(0, 1).
		Width(35)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	b.WriteString(titleStyle.Render("📊 METRICS"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)

	b.WriteString(labelStyle.Render("Customers: "))
	b.WriteString(fmt.Sprintf("%d", fg.Customers))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Growth: "))
	growthStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if fg.MonthlyGrowthRate < 0 {
		growthStyle = lipgloss.NewStyle().Foreground(styles.Red)
	}
	b.WriteString(growthStyle.Render(fmt.Sprintf("%.1f%%", fg.MonthlyGrowthRate*100)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Churn: "))
	b.WriteString(fmt.Sprintf("%.1f%%/mo", fg.ChurnRate*100))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Product Maturity: "))
	b.WriteString(fmt.Sprintf("%.0f%%", fg.ProductMaturity*100))
	b.WriteString("\n")

	return panelStyle.Render(b.String())
}

func (s *FounderGameScreen) renderNews() string {
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(0, 1).
		Width(72)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	b.WriteString(titleStyle.Render("📰 NEWS"))
	b.WriteString("\n")

	// Show up to 5 messages
	limit := 5
	if len(s.turnMessages) < limit {
		limit = len(s.turnMessages)
	}
	for i := 0; i < limit; i++ {
		msg := s.turnMessages[i]
		if len(msg) > 65 {
			msg = msg[:62] + "..."
		}
		b.WriteString("• " + msg + "\n")
	}

	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(panelStyle.Render(b.String()))
}

func (s *FounderGameScreen) renderActions() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)

	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("⚡ CHOOSE ACTION")))
	b.WriteString("\n\n")

	// Menu
	menuContainer := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(s.actionsMenu.View())))
	b.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back • enter select"))

	return b.String()
}
