package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/components"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// MainMenuScreen is the main menu
type MainMenuScreen struct {
	width    int
	height   int
	menu     *components.Menu
	gameData *GameData
}

// NewMainMenuScreen creates a new main menu screen
func NewMainMenuScreen(width, height int, gameData *GameData) *MainMenuScreen {
	menuItems := []components.MenuItem{
		{
			ID:          "new_game",
			Title:       "New Game",
			Description: "Start a new VC or Founder game",
			Icon:        "🎮",
		},
		{
			ID:          "leaderboard",
			Title:       "Leaderboards",
			Description: "View top players and scores",
			Icon:        "🏆",
		},
		{
			ID:          "stats",
			Title:       "Player Statistics",
			Description: "View your gameplay statistics",
			Icon:        "📊",
		},
		{
			ID:          "achievements",
			Title:       "Achievements",
			Description: "View unlocked achievements",
			Icon:        "🎖️",
		},
		{
			ID:          "upgrades",
			Title:       "Upgrades",
			Description: "Purchase upgrades with achievement points",
			Icon:        "⬆️",
		},
		{
			ID:          "progression",
			Title:       "Progression & Levels",
			Description: "View your career progression",
			Icon:        "📈",
		},
		{
			ID:          "analytics",
			Title:       "Analytics Dashboard",
			Description: "Deep dive into your performance",
			Icon:        "📉",
		},
		{
			ID:          "help",
			Title:       "Help & Info",
			Description: "Learn how to play",
			Icon:        "❓",
		},
		{
			ID:          "quit",
			Title:       "Quit",
			Description: "Exit the game",
			Icon:        "🚪",
		},
	}

	menu := components.NewMenu("", menuItems)
	menu.SetSize(40, 20)
	menu.SetHideHelp(true)

	return &MainMenuScreen{
		width:    width,
		height:   height,
		menu:     menu,
		gameData: gameData,
	}
}

// Init initializes the main menu
func (m *MainMenuScreen) Init() tea.Cmd {
	return nil
}

// Update handles main menu input
func (m *MainMenuScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Global.Quit) {
			return m, Quit()
		}

	case components.MenuSelectedMsg:
		return m, m.handleSelection(msg.ID)
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m *MainMenuScreen) handleSelection(id string) tea.Cmd {
	switch id {
	case "new_game":
		return SwitchTo(ScreenVCSetup)
	case "leaderboard":
		return PushTo(ScreenLeaderboard)
	case "stats":
		return PushTo(ScreenStats)
	case "achievements":
		return PushTo(ScreenAchievements)
	case "upgrades":
		return PushTo(ScreenUpgrades)
	case "progression":
		return PushTo(ScreenProgression)
	case "analytics":
		return PushTo(ScreenAnalytics)
	case "help":
		return PushTo(ScreenHelp)
	case "quit":
		return Quit()
	}
	return nil
}

// View renders the main menu
func (m *MainMenuScreen) View() string {
	var b strings.Builder

	// Header with logo
	logoStyle := lipgloss.NewStyle().
		Foreground(styles.Magenta).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	logo := `
🦄 UNICORN 🦄
`
	b.WriteString(logoStyle.Render(logo))
	b.WriteString("\n")

	// Title bar
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(50).
		Align(lipgloss.Center).
		Padding(0, 2)

	title := titleStyle.Render("MAIN MENU")
	titleContainer := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)
	b.WriteString(titleContainer.Render(title))
	b.WriteString("\n\n")

	// Menu
	menuContainer := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2)

	b.WriteString(menuContainer.Render(menuBox.Render(m.menu.View())))

	// Version info
	b.WriteString("\n\n")
	versionStyle := lipgloss.NewStyle().
		Foreground(styles.Gray).
		Width(m.width).
		Align(lipgloss.Center)
	b.WriteString(versionStyle.Render("v2.0 - Bubble Tea Edition"))

	return b.String()
}
