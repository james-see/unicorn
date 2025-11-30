package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// AnalyticsScreen shows analytics dashboard
type AnalyticsScreen struct {
	width  int
	height int
}

// NewAnalyticsScreen creates a new analytics screen
func NewAnalyticsScreen(width, height int) *AnalyticsScreen {
	return &AnalyticsScreen{
		width:  width,
		height: height,
	}
}

// Init initializes the analytics screen
func (s *AnalyticsScreen) Init() tea.Cmd {
	return nil
}

// Update handles analytics input
func (s *AnalyticsScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Global.Back) {
			return s, PopScreen()
		}
	}
	return s, nil
}

// View renders analytics
func (s *AnalyticsScreen) View() string {
	var b strings.Builder
	
	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📉 ANALYTICS DASHBOARD 📉")))
	b.WriteString("\n\n")
	
	// Info box
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(2, 4).
		Width(55)
	
	info := `Analytics Dashboard

This feature provides deep insights into your
investment performance across all games.

Available Analytics:
• Investment category breakdown
• Risk vs. return analysis
• Win/loss patterns by difficulty
• Time series of portfolio values
• Comparison with AI players

Coming Soon: More detailed visualizations!`
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(info)))
	b.WriteString("\n\n")
	
	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))
	
	return b.String()
}
