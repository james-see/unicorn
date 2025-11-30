package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// HelpScreen shows help and info
type HelpScreen struct {
	width        int
	height       int
	scrollOffset int
}

// NewHelpScreen creates a new help screen
func NewHelpScreen(width, height int) *HelpScreen {
	return &HelpScreen{
		width:  width,
		height: height,
	}
}

// Init initializes the help screen
func (s *HelpScreen) Init() tea.Cmd {
	return nil
}

// Update handles help input
func (s *HelpScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Back):
			return s, PopScreen()
		case key.Matches(msg, keys.Global.Down):
			s.scrollOffset++
		case key.Matches(msg, keys.Global.Up):
			if s.scrollOffset > 0 {
				s.scrollOffset--
			}
		}
	}
	return s, nil
}

// View renders help
func (s *HelpScreen) View() string {
	var b strings.Builder
	
	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Yellow).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("❓ HELP & INFO ❓")))
	b.WriteString("\n\n")
	
	// Help content
	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Padding(1, 2).
		Width(65)
	
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	sectionStyle := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true)
	
	var help strings.Builder
	
	help.WriteString(titleStyle.Render("🦄 UNICORN - The VC Simulation Game"))
	help.WriteString("\n\n")
	
	help.WriteString(sectionStyle.Render("HOW TO PLAY"))
	help.WriteString("\n")
	help.WriteString(`
In VC Mode, you manage a venture capital fund:
• Start with a pool of capital
• Invest in promising startups
• Navigate market events and funding rounds
• Compete against AI investors
• Exit investments for profit

Goal: Maximize your net worth by game end!
`)
	
	help.WriteString(sectionStyle.Render("GAME MECHANICS"))
	help.WriteString("\n")
	help.WriteString(`
• Each turn = 1 month
• Startups can: grow, raise rounds, get acquired, fail
• Dilution happens when companies raise new rounds
• Pro-rata rights let you maintain ownership
• Board seats give you voting power
`)
	
	help.WriteString(sectionStyle.Render("INVESTMENT TERMS"))
	help.WriteString("\n")
	help.WriteString(`
• Common Stock: Basic ownership
• Preferred Stock: Better liquidation rights
• SAFE: Simple agreement, converts later
• Convertible Note: Debt that converts to equity
`)
	
	help.WriteString(sectionStyle.Render("KEYBOARD SHORTCUTS"))
	help.WriteString("\n")
	help.WriteString(`
• ↑/↓/←/→ or hjkl: Navigate
• Enter: Select/Confirm
• Esc: Back/Cancel
• q: Quit
• d: Dashboard (in game)
• v: Value-Add (in game)
• s: Secondary Market (in game)
`)
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(helpBox.Render(help.String())))
	b.WriteString("\n\n")
	
	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("↑/↓ scroll • esc back"))
	
	return b.String()
}
