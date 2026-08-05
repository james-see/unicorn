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
	
	help.WriteString(titleStyle.Render("🦄 UNICORN - Startup Adventure"))
	help.WriteString("\n\n")
	
	help.WriteString(sectionStyle.Render("TWO GAME MODES"))
	help.WriteString("\n")
	help.WriteString(`
🎩 VC Investor Mode
   Manage a VC fund, invest in startups, compete with AI investors.
   Goal: Maximize net worth over 60 turns (5 years).

🚀 Startup Founder Mode
   Build your own startup from scratch. Hire, fundraise, acquire customers.
   Goal: Grow to $20M ARR and IPO, or get acquired.
`)
	
	help.WriteString(sectionStyle.Render("VC MODE - HOW TO PLAY"))
	help.WriteString("\n")
	help.WriteString(`
• Start with a pool of capital ($500K-$1M)
• Invest in promising startups
• Navigate market events and funding rounds
• Compete against AI investors like CARL
• Exit investments for profit via secondary market

Goal: Maximize your net worth by game end!
`)
	
	help.WriteString(sectionStyle.Render("FOUNDER MODE - HOW TO PLAY"))
	help.WriteString("\n")
	help.WriteString(`
• Choose a startup template (SaaS, DeepTech, GovTech, Hardware)
• Hire team: engineers, sales, CS, marketing, C-suite
• Acquire customers via direct sales, affiliates, partnerships
• Raise funding rounds: Seed, Series A, B, C
• Manage board, advisors, equity, PR, security, tech debt
• Respond to crises, competitors, and market conditions
• Exit via IPO or acquisition

Goal: Build a unicorn and exit successfully!
`)
	
	help.WriteString(sectionStyle.Render("GAME MECHANICS"))
	help.WriteString("\n")
	help.WriteString(`
VC Mode:
• Each turn = 1 month (60 turns total)
• Startups can: grow, raise rounds, get acquired, fail
• Dilution happens when companies raise new rounds
• Pro-rata rights let you maintain ownership
• Board seats give you voting power

Founder Mode:
• Each turn = 1 month
• MRR grows with marketing, sales, product maturity
• Churn decreases as product matures
• Cash flow = MRR - team costs - infrastructure
`)
	
	help.WriteString(sectionStyle.Render("INVESTMENT TERMS (VC)"))
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
Global:
• ↑/↓/←/→ or hjkl: Navigate
• Enter: Select/Confirm
• Esc: Back/Cancel
• q: Quit (with confirmation in VC mode)

VC Mode (in game):
• d: Dashboard
• v: Value-Add actions
• s: Secondary Market (1-9 select, a accept, r reject)

Founder Mode (in game):
• Enter: Open Actions menu
• n: Quick advance to next month
• Esc/q: Quit confirmation
`)
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(helpBox.Render(help.String())))
	b.WriteString("\n\n")
	
	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("↑/↓ scroll • esc back"))
	
	return b.String()
}
